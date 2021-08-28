package main

import (
	_ "exampleoneway/plugin/cronjob"
	_ "exampleoneway/plugin/input"
	_ "exampleoneway/plugin/output"
	_ "exampleoneway/plugin/process"
	_ "exampleoneway/plugin/transit"

	"github.com/jdjgya/service-frame/pkg/oneway/controller"
)

func main() {
	c := controller.GetInstance()
	c.InitService()

	c.Start()

	c.TrapSignals()
	c.TraceStatus()
}
