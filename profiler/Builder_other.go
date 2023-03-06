//go:build !windows
// +build !windows

package profiler

import (
	"log"
	"os"
	"os/exec"
)

func init() {
	DefaultLocation = "$HOME/.config/netprofiles/$HOSTNAME"
	IsForce = false

	Profilers = []Profiler{
		&FileProfiler{
			Name: "netplan",
			Files: []File{
				{Path: "/etc/netplan/00-installer-config.yaml", RootPrivilege: true},
				{Path: "/etc/netplan/01-network-manager-all.yaml", RootPrivilege: true},
				{Path: "/etc/netplan/99-custom.yaml", RootPrivilege: true},
			},
			PostLoad: func() error {
				log.Println("Running 'netplan generate --debug'")
				_, err := exec.Command("sudo", "netplan", "generate", "--debug").Output()
				if err != nil {
					log.Fatalf("'sudo netplan generate --debug' failed: %v", err)
					return err
				}

				log.Println("Running 'netplan apply'")
				_, err = exec.Command("sudo", "netplan", "apply").Output()
				if err != nil {
					log.Fatal(err)
					return err
				}
				return nil
			},
		},
		&FileProfiler{
			Name:  "hosts",
			Files: []File{{Path: "/etc/hosts", RootPrivilege: true}},
		},
		&FileProfiler{
			Name:  "apt",
			Files: []File{{Path: "/etc/apt/apt.conf.d/02proxy", RootPrivilege: true}},
		},
		&FileProfiler{
			Name: "docker",
			Files: []File{
				// https://docs.docker.com/config/daemon/systemd/#httphttps-proxy
				{Path: "/etc/systemd/system/docker.service.d/proxy.conf", RootPrivilege: true},
				{Path: os.ExpandEnv("$HOME/.docker/config.json"), RootPrivilege: false},
			},
			PostLoad: func() error {
				log.Println("Running 'systemctl daemon-reload'")
				_, err := exec.Command("sudo", "systemctl", "daemon-reload").Output()
				if err != nil {
					log.Fatalf("'systemctl daemon-reload' failed: %v", err)
					return err
				}

				log.Println("Running 'systemctl restart docker'")
				_, err = exec.Command("sudo", "systemctl", "restart", "docker").Output()
				if err != nil {
					log.Fatal(err)
					return err
				}
				return nil
			},
		},
		&FileProfiler{
			Name: "containerd",
			Files: []File{
				{Path: "/etc/systemd/system/containerd.service.d/proxy.conf", RootPrivilege: true},
				{Path: "/etc/containerd/config.toml", RootPrivilege: true},
			},
			PostLoad: func() error {
				log.Println("Running 'systemctl daemon-reload'")
				_, err := exec.Command("sudo", "systemctl", "daemon-reload").Output()
				if err != nil {
					log.Fatalf("'systemctl daemon-reload' failed: %v", err)
					return err
				}

				log.Println("Running 'systemctl restart containerd'")
				_, err = exec.Command("sudo", "systemctl", "restart", "containerd").Output()
				if err != nil {
					log.Fatal(err)
					return err
				}
				return nil
			},
		},
		&FileProfiler{
			Name:  "git",
			Files: []File{{Path: os.ExpandEnv("$HOME/.gitconfig"), RootPrivilege: false}},
		},
	}
}
