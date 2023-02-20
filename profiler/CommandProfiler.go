package profiler

type CommandProfiler struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func (fp *CommandProfiler) Save(profile, location string) error {
	return nil
}

func (fp *CommandProfiler) Load(profile, location string) error {
	return nil
}
