package main

import (
	"net/http"
	_ "net/http/pprof"

	_ "github.com/jdjgya/service-frame/example/oneway/plugin/cronjob"
	_ "github.com/jdjgya/service-frame/example/oneway/plugin/input"
	_ "github.com/jdjgya/service-frame/example/oneway/plugin/output"
	_ "github.com/jdjgya/service-frame/example/oneway/plugin/process"
	_ "github.com/jdjgya/service-frame/example/oneway/plugin/transit"

	"github.com/jdjgya/service-frame/pkg/oneway/controller"
)

func main() {
	go http.ListenAndServe("localhost:6060", nil)

	c := controller.GetInstance()
	c.InitService()
	c.ActivateService()

	c.Start()

	c.TrapSignals()
	c.TraceStatus()
}
