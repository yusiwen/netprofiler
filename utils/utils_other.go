//go:build !windows
// +build !windows

package utils

import (
	"os/exec"
)

func CopySudo(src, dst string) error {
	cmd := exec.Command("sudo", "cp", "-v", src, dst)
	return cmd.Run()
}
