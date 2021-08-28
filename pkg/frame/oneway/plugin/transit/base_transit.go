package transit

import (
	"context"
	"onewayframe/pkg/plugin"
	"onewayframe/pkg/plugin/plug"
	"sync"
)

var (
	Plugin = make(map[string]Transit)
)

type Transit interface {
	plug.Parter
	DoTransit()
}

func WrapWithTransitLoop(ctx context.Context, wg *sync.WaitGroup, coreFunc func([]byte) ([]map[string]string, error)) func() {
	return func() {
		wg.Add(1)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, isChnOpen := <-plugin.I2TChan:
				switch isChnOpen {
				case true:
					jobs, err := coreFunc(msg)
					if err != nil || len(jobs) == 0 {
						plugin.Metrics.TransitErr++
						continue
					}

					for _, message := range jobs {
						plugin.T2PChan <- message
						plugin.Metrics.TransitOK++
					}
				case false:
					close(plugin.T2PChan)
					plugin.TransitStatus.Completed = true
					return
				}
			default:
			}
		}
	}
}
