//go:build windows
// +build windows

package profiler

import (
	"os"
)

var Profilers []Profiler

func init() {
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
