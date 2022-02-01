package plug

var Stagers = make(map[string]Stager)

type StageNewer interface {
	New() Stager
}

type StageConfigSetter interface {
	SetConfig(interface{})
}

type Stager interface {
	StageConfigSetter
	ConfigChecker
	Execute(*map[string]string) error
	Statuser
}

type StageUser interface {
	SetStager(string, int, interface{})
}
