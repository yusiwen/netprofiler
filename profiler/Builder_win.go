//go:build windows
// +build windows

package profiler

import (
	"os"
)

func init() {
	DefaultLocation = "$USERPROFILE\\.config\\netprofiles\\$COMPUTERNAME"
	IsForce = false

	Profilers = []Profiler{
		&FileProfiler{
			Name:  "hosts",
			Files: []File{{Path: "c:\\Windows\\System32\\drivers\\etc\\hosts", RootPrivilege: true}},
		},
		&FileProfiler{
			Name:  "git",
			Files: []File{{Path: os.ExpandEnv("$USERPROFILE\\.gitconfig"), RootPrivilege: false}},
		},
	}

}
