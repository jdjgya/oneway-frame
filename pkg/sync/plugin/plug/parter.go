package plug

var Parters = make(map[string]Parter)

type PartConfigSetter interface {
	SetConfig(interface{})
}

type Parter interface {
	PartConfigSetter
	ConfigChecker
	Stopper
}

type PartUser interface {
	StageUser
	SetParter(string)
	StartParters()
	StopParters()
}
