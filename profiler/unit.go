package profiler

type Unit interface {
	Save(profile, location string) error
	Load(profile, location string) error
}
