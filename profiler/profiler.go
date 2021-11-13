package profiler

type Profiler interface {
	Save(profile, location string) error
	Load(profile, location string) error
}
