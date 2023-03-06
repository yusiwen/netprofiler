package profiler

import (
	"bufio"
	"fmt"
	"github.com/yusiwen/netprofiles/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Profile interface {
	Save(profile string) error
	Load(profile string) error
	PostUp() error
	PostDown() error
	ListProfiles() error
	GetCurrentProfile() string
}

type ProfileManager struct {
	Units    []Unit
	Location string
	IsForce  bool
}

func (p *ProfileManager) Save(profile string) error {
	err := utils.CreateIfNotExists(os.ExpandEnv(p.Location), os.ModePerm)
	if err != nil {
		return err
	}

	// Overwriting confirmation
	if !p.IsForce {
		if utils.Exists(filepath.Join(os.ExpandEnv(p.Location), profile)) {
			if !utils.AskForConfirmation(fmt.Sprintf("Profile '%s' already exists, overwrite it?", profile)) {
				return nil
			}
		}
	}

	for _, u := range p.Units {
		err := u.Save(profile, os.ExpandEnv(p.Location))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Unit '%s' saved\n", profile)
	return nil
}

func (p *ProfileManager) Load(profile string) error {
	currentProfile := p.GetCurrentProfile()
	if profile == currentProfile {
		if !utils.AskForConfirmation(fmt.Sprintf("Profile '%s' is currently loaded, force reloading?", profile)) {
			return nil
		}
	}

	for _, u := range p.Units {
		err := u.Load(profile, os.ExpandEnv(p.Location))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Profile '%s' loaded\n", profile)
	file, err := os.Create(filepath.Join(os.ExpandEnv(p.Location), ".current"))
	if err != nil {
		log.Printf("save current profile failed: %v\n", err)
		return nil
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("warning: file close error: %v\n", err)
		}
	}(file)
	w := bufio.NewWriter(file)
	_, err = w.WriteString(profile)
	if err != nil {
		log.Printf("error: file write error: %v\n", err)
		return err
	}

	err = w.Flush()
	if err != nil {
		log.Printf("error: file flush error: %v\n", err)
		return err
	}

	return nil
}

func (p *ProfileManager) PostUp() error {
	//TODO implement me
	panic("implement me")
}

func (p *ProfileManager) PostDown() error {
	//TODO implement me
	panic("implement me")
}

func (p *ProfileManager) GetCurrentProfile() string {
	path := filepath.Join(os.ExpandEnv(p.Location), ".current")
	if !utils.Exists(path) {
		return ""
	}
	file, err := os.Open(path)
	if err != nil {
		log.Printf("failed to get current profile: %v\n", err)
		return ""
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("warning: file close error: %v\n", err)
		}
	}(file)

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("failed to get current profile: %v\n", err)
		return ""
	}

	return string(bytes)
}

func (p *ProfileManager) ListProfiles() error {
	if !utils.Exists(os.ExpandEnv(p.Location)) {
		return nil
	}

	files, err := ioutil.ReadDir(os.ExpandEnv(p.Location))
	currentProfile := p.GetCurrentProfile()

	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if f.Name() == currentProfile {
				fmt.Println(">> " + f.Name())
			} else {
				fmt.Println(f.Name())
			}
		}
	}

	return nil
}
