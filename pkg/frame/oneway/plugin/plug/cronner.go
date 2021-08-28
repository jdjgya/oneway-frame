package plug

var Cronners = make(map[string]Cronner)

type CronConfigSetter interface {
	SetConfig(map[interface{}]interface{})
}

type Cronner interface {
	CronConfigSetter
	ConfigChecker
	Stopper
}

type CronUser interface {
	SetCronner(string)
	StartCronners()
	StopCronners()
}
