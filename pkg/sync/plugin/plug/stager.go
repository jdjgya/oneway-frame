package plug

var Stagers = make(map[string]Stager)

type StageConfigSetter interface {
	SetConfig(interface{})
}

type Stager interface {
	StageConfigSetter
	ConfigChecker
	Execute(*map[string]string) (bool, error)
	Statuser
}

type StageUser interface {
	SetStager(string, int, interface{})
}
