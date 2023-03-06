package profiler

var Profilers []Profiler
var DefaultLocation string
var IsForce bool

type Profiler interface {
	Save(profile, location string) error
	Load(profile, location string) error
}
