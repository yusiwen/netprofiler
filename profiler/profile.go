package profiler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/yusiwen/netprofiles/utils"
	"io"
	"log"
	"os"
	"path/filepath"
)

type ProfileManager interface {
	Save(profile string) error
	Load(profile string) error

	PostUp() error
	PostDown() error

	SetPostUpCommand(command string) error
	SetPortDownCommand(command string) error

	GetLocation() string
	SetLocation(location string)
	IsForce() bool
	SetForce(force bool)
	ListProfiles() error
	ListUnits() error
	GetCurrentProfile() *Profile
}

type DefaultProfileManager struct {
	Units    []Unit
	Location string
	Force    bool
}

var PM ProfileManager

type Profile struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func (p *DefaultProfileManager) Save(profile string) error {
	err := utils.CreateIfNotExists(os.ExpandEnv(p.Location), os.ModePerm)
	if err != nil {
		return err
	}

	// Overwriting confirmation
	if !p.Force {
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
	fmt.Printf("Profile '%s' saved\n", profile)
	return nil
}

func (p *DefaultProfileManager) Load(profile string) error {
	currentProfile := p.GetCurrentProfile()
	if currentProfile != nil && profile == currentProfile.Name {
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

	info := Profile{Name: profile, Location: os.ExpandEnv(p.Location)}
	bytes, err := json.Marshal(info)
	fmt.Println(string(bytes))
	if err != nil {
		log.Printf("error: marshal json error: %v\n", err)
		return err
	}
	w := bufio.NewWriter(file)
	_, err = w.Write(bytes)
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

func (p *DefaultProfileManager) PostUp() error {
	//TODO implement me
	panic("implement me")
}

func (p *DefaultProfileManager) PostDown() error {
	//TODO implement me
	panic("implement me")
}

func (p *DefaultProfileManager) GetCurrentProfile() *Profile {
	path := filepath.Join(os.ExpandEnv(p.Location), ".current")
	if !utils.Exists(path) {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		log.Printf("failed to get current profile: %v\n", err)
		return nil
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
		return nil
	}

	profile := Profile{}
	err = json.Unmarshal(bytes, &profile)
	if err != nil {
		log.Printf("failed to parse .current: %v\n", err)
		return nil
	}
	return &profile
}

func (p *DefaultProfileManager) ListProfiles() error {
	if !utils.Exists(os.ExpandEnv(p.Location)) {
		return nil
	}

	files, err := os.ReadDir(os.ExpandEnv(p.Location))
	if err != nil {
		log.Fatal(err)
		return err
	}

	var current = ""
	currentProfile := p.GetCurrentProfile()
	if currentProfile != nil {
		current = currentProfile.Name
	}

	for _, f := range files {
		if f.IsDir() {
			if f.Name() == current {
				fmt.Println(">> " + f.Name())
			} else {
				fmt.Println(f.Name())
			}
		}
	}

	return nil
}

func (p *DefaultProfileManager) ListUnits() error {
	for _, p := range p.Units {
		s, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s\n", s)
		}
	}
	return nil
}

func (p *DefaultProfileManager) SetPostUpCommand(command string) error {
	//TODO implement me
	panic("implement me")
}

func (p *DefaultProfileManager) SetPortDownCommand(command string) error {
	//TODO implement me
	panic("implement me")
}

func (p *DefaultProfileManager) GetLocation() string {
	return p.Location
}

func (p *DefaultProfileManager) IsForce() bool {
	return p.Force
}

func (p *DefaultProfileManager) SetForce(force bool) {
	p.Force = force
}

func (p *DefaultProfileManager) SetLocation(location string) {
	p.Location = location
}
