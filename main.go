package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"git.yusiwen.cn/yusiwen/netprofiler/constant"
	P "git.yusiwen.cn/yusiwen/netprofiler/profiler"
	R "github.com/urfave/cli/v2"
)

type AppOptions struct {
	EnableDNS bool
	Template  string
	Output    string
	Secret    string
	Port      int
	RedirPort int
	LogLevel  string
}

var profilers = []P.Profiler{
	&P.FileProfiler{
		Name:  "netplan",
		Files: []string{"/etc/netplan/99-custom.yaml"},
		PostLoad: func() error {
			log.Println("Running 'netplan generate --debug'")
			out, err := exec.Command("netplan", "generate", "--debug").Output()
			if err != nil {
				log.Fatal(err)
				return err
			}
			fmt.Println(out)

			log.Println("Running 'netplan apply'")
			out, err = exec.Command("netplan", "apply").Output()
			if err != nil {
				log.Fatal(err)
				return err
			}
			fmt.Println(out)
			return nil
		},
	},
	&P.FileProfiler{
		Name:  "apt",
		Files: []string{"/etc/apt/apt.conf.d/02proxy"},
	},
	&P.FileProfiler{
		Name: "docker",
		Files: []string{
			"/etc/systemd/system/docker.service.d/proxy.conf",
			"$HOME/.docker/config.json",
		},
	},
	&P.FileProfiler{
		Name:  "git",
		Files: []string{"$HOME/.gitconfig"},
	},
}

func save(profile string) error {
	log.Printf("Saving to profile '%s'\n", profile)
	for _, p := range profilers {
		err := p.Save(profile, "$HOME/.config/netprofiles")
		if err != nil {
			return err
		}
	}
	return nil
}

func load(profile string) error {
	log.Printf("Loading profile '%s'\n", profile)
	for _, p := range profilers {
		err := p.Load(profile, "$HOME/.config/netprofiles")
		if err != nil {
			return err
		}
	}
	return nil
}

func list() error {
	return nil
}

func main() {
	app := &R.App{
		Name:    "netprofiler",
		Usage:   "My network profiles switcher for working at home, office and business travels",
		Version: strings.Join([]string{constant.Version, " (", constant.BuildTime, ")"}, ""),
		Commands: []*R.Command{
			{
				Name:    "save",
				Aliases: []string{"S"},
				Usage:   "save current environment to a profile",
				Action: func(c *R.Context) error {
					profile := c.Args().First()
					if len(profile) == 0 {
						return errors.New("Profile name must not be null")
					}
					return save(profile)
				},
			},
			{
				Name:    "load",
				Aliases: []string{"L"},
				Usage:   "load a profile to system",
				Action: func(c *R.Context) error {
					profile := c.Args().First()
					if len(profile) == 0 {
						return errors.New("Profile name must not be null")
					}
					return load(profile)
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all profiles",
				Action: func(c *R.Context) error {
					fmt.Println("added task: ", c.Args().First())
					return list()
				},
			},
		}}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
