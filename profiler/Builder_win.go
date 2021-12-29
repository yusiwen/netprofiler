//go:build windows
// +build windows

package profiler

import (
	"os"
)

var Profilers []Profiler
var DefaultLocation string

func init() {
	DefaultLocation = "$USERPROFILE\\.config\\netprofiles"

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
