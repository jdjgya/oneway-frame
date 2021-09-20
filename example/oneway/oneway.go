package main

import (
	_ "github.com/jdjgya/service-frame/example/oneway/plugin/cronjob"
	_ "github.com/jdjgya/service-frame/example/oneway/plugin/input"
	_ "github.com/jdjgya/service-frame/example/oneway/plugin/output"
	_ "github.com/jdjgya/service-frame/example/oneway/plugin/process"
	_ "github.com/jdjgya/service-frame/example/oneway/plugin/transit"

	"github.com/jdjgya/service-frame/pkg/oneway/controller"
)

func main() {
	c := controller.GetInstance()
	c.InitService()
	c.ActivateService()

	c.Start()

	c.TrapSignals()
	c.TraceStatus()
}
