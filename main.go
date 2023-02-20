package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.yusiwen.cn/yusiwen/netprofiler/constant"
	P "git.yusiwen.cn/yusiwen/netprofiler/profiler"
	"git.yusiwen.cn/yusiwen/netprofiler/utils"
	R "github.com/urfave/cli/v2"
)

func getCurrentProfile() string {
	path := filepath.Join(os.ExpandEnv(P.DefaultLocation), ".current")
	if !utils.Exists(path) {
		return ""
	}
	file, err := os.Open(path)
	if err != nil {
		log.Printf("failed to get current profile: %v\n", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("warning: file close error: %v\n", err)
		}
	}(file)

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("failed to get current profile: %v\n", err)
		return ""
	}
	return string(bytes)
}

func save(profile string) error {
	err := utils.CreateIfNotExists(os.ExpandEnv(P.DefaultLocation), os.ModePerm)
	if err != nil {
		return err
	}

	// Overwriting confirmation
	if !P.IsForce {
		if utils.Exists(filepath.Join(os.ExpandEnv(P.DefaultLocation), profile)) {
			if !utils.AskForConfirmation(fmt.Sprintf("Profile '%s' already exists, overwrite it?", profile)) {
				return nil
			}
		}
	}

	for _, p := range P.Profilers {
		err := p.Save(profile, os.ExpandEnv(P.DefaultLocation))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Profile '%s' saved\n", profile)
	return nil
}

func load(profile string) error {
	for _, p := range P.Profilers {
		err := p.Load(profile, os.ExpandEnv(P.DefaultLocation))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Profile '%s' loaded\n", profile)
	file, err := os.Create(filepath.Join(os.ExpandEnv(P.DefaultLocation), ".current"))
	if err != nil {
		log.Printf("save current profile failed: %v\n", err)
		return nil
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("warning: file close error: %v\n", err)
		}
	}(file)
	w := bufio.NewWriter(file)
	_, err = w.WriteString(profile)
	if err != nil {
		log.Printf("error: file write error: %v\n", err)
		return err
	}

	err = w.Flush()
	if err != nil {
		log.Printf("error: file flush error: %v\n", err)
		return err
	}

	return nil
}

func list() error {
	if !utils.Exists(os.ExpandEnv(P.DefaultLocation)) {
		return nil
	}

	files, err := ioutil.ReadDir(os.ExpandEnv(P.DefaultLocation))
	currentProfile := getCurrentProfile()

	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if f.Name() == currentProfile {
				fmt.Println(">> " + f.Name())
			} else {
				fmt.Println(f.Name())
			}
		}
	}

	return nil
}

func processCopyCommand(srcProfile string, dstProfile string) error {
	src := filepath.Join(os.ExpandEnv(P.DefaultLocation), srcProfile)
	if !utils.Exists(src) {
		return fmt.Errorf("profile '%s' not exists", srcProfile)
	}
	dst := filepath.Join(os.ExpandEnv(P.DefaultLocation), dstProfile)
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

func processDeleteCommand(profile string) error {
	path := filepath.Join(os.ExpandEnv(P.DefaultLocation), profile)
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
				Value:   P.DefaultLocation,
				Usage:   "Set location to save profiles",
			},
		},
		Commands: []*R.Command{
			{
				Name:    "save",
				Aliases: []string{"S"},
				Usage:   "save current environment to a profile",
				Action: func(c *R.Context) error {
					P.DefaultLocation = c.String("location")
					if c.Bool("force") {
						P.IsForce = true
					}
					profile := c.Args().First()
					if len(profile) == 0 {
						return errors.New("profile name must not be null")
					}
					return save(profile)
				},
				Flags: []R.Flag{
					&R.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "Save without confirmation",
					},
				},
			},
			{
				Name:    "load",
				Aliases: []string{"L"},
				Usage:   "load a profile to system",
				Action: func(c *R.Context) error {
					P.DefaultLocation = c.String("location")
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
					P.DefaultLocation = c.String("location")
					return list()
				},
			},
			{
				Name:    "copy",
				Aliases: []string{"C"},
				Usage:   "copy profile to another profile",
				Action: func(c *R.Context) error {
					P.DefaultLocation = c.String("location")
					if c.Args().Len() != 2 {
						return errors.New("wrong parameter")
					}
					return processCopyCommand(c.Args().First(), c.Args().Get(1))
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"D"},
				Usage:   "delete a profile",
				Action: func(c *R.Context) error {
					P.DefaultLocation = c.String("location")
					profile := c.Args().First()
					if len(profile) == 0 {
						return errors.New("profile name must not be null")
					}
					return processDeleteCommand(profile)
				},
			},
			{
				Name:    "list-profilers",
				Aliases: []string{"p"},
				Usage:   "list all profilers",
				Action: func(c *R.Context) error {
					for _, p := range P.Profilers {
						s, err := json.Marshal(p)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Printf("%s\n", s)
						}
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
