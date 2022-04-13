package process

import (
	"context"
	"sync"

	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/plug"
)

var (
	Plugin = make(map[string]Process)
)

type Process interface {
	plug.Parter
	DoProcess()
}

func WrapWithProcessLoop(ctx context.Context, wg *sync.WaitGroup, group string, coreFunc func(map[string]string, bool) (map[string]string, error)) func() {
	return func() {
		wg.Add(1)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, isChnOpen := <-plugin.T2PChan[group]:
				switch isChnOpen {
				case true:
					msg, err := coreFunc(msg, isChnOpen)
					if err != nil {
						plugin.Metrics.ProcessErr++
						continue
					}

					plugin.P2OChan[group] <- msg
					plugin.Metrics.ProcessOK++
				case false:
					close(plugin.P2OChan[group])
					plugin.ProcessStatus.Completed = true
					return
				}
			}
		}
	}
}
