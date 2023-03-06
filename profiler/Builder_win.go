//go:build windows
// +build windows

package profiler

import (
	"os"
)

func init() {
	PM = &DefaultProfileManager{
		Location: "$USERPROFILE\\.config\\netprofiles\\$COMPUTERNAME",
		Force:    false,
		Units: []Unit{
			&FileUnit{
				Name:  "hosts",
				Files: []File{{Path: "c:\\Windows\\System32\\drivers\\etc\\hosts", RootPrivilege: true}},
			},
			&FileUnit{
				Name:  "git",
				Files: []File{{Path: os.ExpandEnv("$USERPROFILE\\.gitconfig"), RootPrivilege: false}},
			},
		},
	}
}
