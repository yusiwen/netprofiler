package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"git.yusiwen.cn/yusiwen/netprofiler/constant"
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

func process(opts *AppOptions) error {
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
					fmt.Println("added task: ", c.Args().First())
					return nil
				},
			},
			{
				Name:    "load",
				Aliases: []string{"L"},
				Usage:   "load a profile to system",
				Action: func(c *R.Context) error {
					fmt.Println("added task: ", c.Args().First())
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all profiles",
				Action: func(c *R.Context) error {
					fmt.Println("added task: ", c.Args().First())
					return nil
				},
			},
		}}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
