//go:build windows
// +build windows

package profiler

import (
	"os"
)

var PM DefaultProfileManager

func init() {
	PM.Location = "$USERPROFILE\\.config\\netprofiles\\$COMPUTERNAME"
	PM.IsForce = false

	PM.Units = []Unit{
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
