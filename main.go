package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"git.yusiwen.cn/yusiwen/netprofiler/constant"
	P "git.yusiwen.cn/yusiwen/netprofiler/profiler"
	"git.yusiwen.cn/yusiwen/netprofiler/utils"
	R "github.com/urfave/cli/v2"
)

var default_location string = "$HOME/.config/netprofiles"

var profilers = []P.Profiler{
	&P.FileProfiler{
		Name:  "netplan",
		Files: []P.File{{Path: "/etc/netplan/99-custom.yaml", RootPrivilege: true}},
		PostLoad: func() error {
			log.Println("Running 'netplan generate --debug'")
			_, err := exec.Command("sudo", "netplan", "generate", "--debug").Output()
			if err != nil {
				log.Fatalf("'sudo netplan generate --debug' failed: %v", err)
				return err
			}

			log.Println("Running 'netplan apply'")
			_, err = exec.Command("sudo", "netplan", "apply").Output()
			if err != nil {
				log.Fatal(err)
				return err
			}
			return nil
		},
	},
	&P.FileProfiler{
		Name:  "hosts",
		Files: []P.File{{Path: "/etc/hosts", RootPrivilege: true}},
	},
	&P.FileProfiler{
		Name:  "apt",
		Files: []P.File{{Path: "/etc/apt/apt.conf.d/02proxy", RootPrivilege: true}},
	},
	&P.FileProfiler{
		Name: "docker",
		Files: []P.File{
			{Path: "/etc/systemd/system/docker.service.d/proxy.conf", RootPrivilege: true},
			{Path: os.ExpandEnv("$HOME/.docker/config.json"), RootPrivilege: false},
		},
		PostLoad: func() error {
			log.Println("Running 'systemctl daemon-reload'")
			_, err := exec.Command("sudo", "systemctl", "daemon-reload").Output()
			if err != nil {
				log.Fatalf("'systemctl daemon-reload' failed: %v", err)
				return err
			}

			log.Println("Running 'systemctl restart docker'")
			_, err = exec.Command("sudo", "systemctl", "restart", "docker").Output()
			if err != nil {
				log.Fatal(err)
				return err
			}
			return nil
		},
	},
	&P.FileProfiler{
		Name:  "git",
		Files: []P.File{{Path: os.ExpandEnv("$HOME/.gitconfig"), RootPrivilege: false}},
	},
}

func save(profile string) error {
	for _, p := range profilers {
		err := p.Save(profile, os.ExpandEnv(default_location))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Profile '%s' saved\n", profile)
	return nil
}

func load(profile string) error {
	for _, p := range profilers {
		err := p.Load(profile, os.ExpandEnv(default_location))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Profile '%s' loaded\n", profile)
	return nil
}

func list() error {
	files, err := ioutil.ReadDir(os.ExpandEnv(default_location))

	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	return nil
}

func copy(srcProfile string, dstProfile string) error {
	src := filepath.Join(os.ExpandEnv(default_location), srcProfile)
	if !utils.Exists(src) {
		return fmt.Errorf("profile '%s' not exists", srcProfile)
	}
	dst := filepath.Join(os.ExpandEnv(default_location), dstProfile)
	if utils.Exists(dst) {
		return fmt.Errorf("profile '%s' already exists", dstProfile)
	}

	err := utils.CopyDirectory(src, dst)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Copy '%s' to '%s'\n", srcProfile, dstProfile)
	return nil
}

func delete(profile string) error {
	path := filepath.Join(os.ExpandEnv(default_location), profile)
	if !utils.Exists(path) {
		return fmt.Errorf("profile '%s' not exists", profile)
	}
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	fmt.Printf("Profile '%s' deleted\n", profile)
	return nil
}

func main() {
	app := &R.App{
		Name:    "netprofiler",
		Usage:   "My network profiles switcher for working at home, office and business travels",
		Version: strings.Join([]string{constant.Version, " (", constant.BuildTime, ")"}, ""),
		Flags: []R.Flag{
			&R.StringFlag{
				Name:    "location",
				Aliases: []string{"l"},
				Value:   "$HOME/.config/netprofiles",
				Usage:   "Set location to save profiles",
			},
		},
		Commands: []*R.Command{
			{
				Name:    "save",
				Aliases: []string{"S"},
				Usage:   "save current environment to a profile",
				Action: func(c *R.Context) error {
					default_location = c.String("location")
					profile := c.Args().First()
					if len(profile) == 0 {
						return errors.New("profile name must not be null")
					}
					return save(profile)
				},
			},
			{
				Name:    "load",
				Aliases: []string{"L"},
				Usage:   "load a profile to system",
				Action: func(c *R.Context) error {
					default_location = c.String("location")
					profile := c.Args().First()
					if len(profile) == 0 {
						return errors.New("profile name must not be null")
					}
					return load(profile)
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all profiles",
				Action: func(c *R.Context) error {
					default_location = c.String("location")
					return list()
				},
			},
			{
				Name:    "copy",
				Aliases: []string{"C"},
				Usage:   "copy profile to another profile",
				Action: func(c *R.Context) error {
					default_location = c.String("location")
					if c.Args().Len() != 2 {
						return errors.New("wrong parameter")
					}
					return copy(c.Args().First(), c.Args().Get(1))
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"D"},
				Usage:   "delete a profile",
				Action: func(c *R.Context) error {
					default_location = c.String("location")
					profile := c.Args().First()
					if len(profile) == 0 {
						return errors.New("profile name must not be null")
					}
					return delete(profile)
				},
			},
		}}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
