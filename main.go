package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	R "github.com/urfave/cli/v2"
	"github.com/yusiwen/netprofiles/constant"
	P "github.com/yusiwen/netprofiles/profiler"
	"github.com/yusiwen/netprofiles/utils"
)

func processCopyCommand(srcProfile string, dstProfile string) error {
	src := filepath.Join(os.ExpandEnv(P.DefaultProfileManager.Location), srcProfile)
	if !utils.Exists(src) {
		return fmt.Errorf("profile '%s' not exists", srcProfile)
	}
	dst := filepath.Join(os.ExpandEnv(P.DefaultProfileManager.Location), dstProfile)
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
	path := filepath.Join(os.ExpandEnv(P.DefaultProfileManager.Location), profile)
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
				Value:   P.DefaultProfileManager.Location,
				Usage:   "Set location to save profiles",
			},
		},
		Commands: []*R.Command{
			{
				Name:    "save",
				Aliases: []string{"S"},
				Usage:   "save current environment to a profile",
				Action: func(c *R.Context) error {
					P.DefaultProfileManager.Location = c.String("location")
					if c.Bool("force") {
						P.DefaultProfileManager.IsForce = true
					}
					profile := c.Args().First()
					if len(profile) == 0 {
						return errors.New("profile name must not be null")
					}
					return P.DefaultProfileManager.Save(profile)
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
					P.DefaultProfileManager.Location = c.String("location")
					profile := c.Args().First()
					if len(profile) == 0 {
						return errors.New("profile name must not be null")
					}
					return P.DefaultProfileManager.Load(profile)
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all profiles",
				Action: func(c *R.Context) error {
					P.DefaultProfileManager.Location = c.String("location")
					return P.DefaultProfileManager.ListProfiles()
				},
			},
			{
				Name:    "copy",
				Aliases: []string{"C"},
				Usage:   "copy profile to another profile",
				Action: func(c *R.Context) error {
					P.DefaultProfileManager.Location = c.String("location")
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
					P.DefaultProfileManager.Location = c.String("location")
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
					for _, p := range P.DefaultProfileManager.Units {
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
			{
				Name:    "current",
				Aliases: []string{"c"},
				Usage:   "get current profile",
				Action: func(c *R.Context) error {
					fmt.Println(P.DefaultProfileManager.GetCurrentProfile())
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
