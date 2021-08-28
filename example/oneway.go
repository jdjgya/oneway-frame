package example

import (
	_ "cronjob/c"
	_ "input/i"
	_ "output/o"
	_ "process/p"
	_ "transit/t"
)

func main() {
	c := controller.GetInstance()
	c.InitService()

	c.Start()

	c.TrapSignals()
	c.TraceStatus()
}
