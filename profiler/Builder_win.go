//go:build windows
// +build windows

package profiler

import (
	"os"
)

var DefaultProfileManager ProfileManager

func init() {
	DefaultProfileManager.Location = "$USERPROFILE\\.config\\netprofiles\\$COMPUTERNAME"
	DefaultProfileManager.IsForce = false

	DefaultProfileManager.Units = []Unit{
		&FileUnit{
			Name:  "hosts",
			Files: []File{{Path: "c:\\Windows\\System32\\drivers\\etc\\hosts", RootPrivilege: true}},
		},
		&FileUnit{
			Name:  "git",
			Files: []File{{Path: os.ExpandEnv("$USERPROFILE\\.gitconfig"), RootPrivilege: false}},
		},
	}

}
