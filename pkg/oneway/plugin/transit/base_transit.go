package transit

import (
	"context"

	"sync"

	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/plug"
)

var (
	Plugin = make(map[string]Transit)
)

type Transit interface {
	plug.Parter
	DoTransit()
}

func WrapWithTransitLoop(ctx context.Context, wg *sync.WaitGroup, group string, coreFunc func([]byte) ([]map[string]string, error)) func() {
	return func() {
		wg.Add(1)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, isChnOpen := <-plugin.I2TChan[group]:
				switch isChnOpen {
				case true:
					jobs, err := coreFunc(msg)
					if err != nil || len(jobs) == 0 {
						plugin.Metrics.TransitErr++
						continue
					}

					for _, message := range jobs {
						plugin.T2PChan[group] <- message
						plugin.Metrics.TransitOK++
					}
				case false:
					close(plugin.T2PChan[group])
					plugin.TransitStatus.Completed = true
					return
				}
			}
		}
	}
}
