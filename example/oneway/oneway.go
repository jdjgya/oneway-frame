package exampleoneway

import (
	"github.com/jdjgya/service-frame/pkg/frame/oneway/controller"
)

func main() {
	c := controller.GetInstance()
	c.Init
	// c := controller.GetInstance()
	// c.InitService()

	// c.Start()

	// c.TrapSignals()
	// c.TraceStatus()
}
