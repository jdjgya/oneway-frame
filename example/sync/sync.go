package main

import (
	_ "github.com/jdjgya/service-frame/example/sync/plugin/cronjob"
	_ "github.com/jdjgya/service-frame/example/sync/plugin/interact/http"
	_ "github.com/jdjgya/service-frame/example/sync/plugin/interact/http/pattern/dummy0"
	_ "github.com/jdjgya/service-frame/example/sync/plugin/interact/http/pattern/dummy1"
	_ "github.com/jdjgya/service-frame/example/sync/plugin/interact/stage/process"
	_ "github.com/jdjgya/service-frame/example/sync/plugin/interact/stage/request"
	_ "github.com/jdjgya/service-frame/example/sync/plugin/interact/stage/transit"

	"github.com/jdjgya/service-frame/pkg/sync/controller"
)

func main() {
	c := controller.GetInstance()
	c.InitService()
	c.ActivateService()

	c.Start()

	c.TrapSignals()
	c.TraceStatus()
}
