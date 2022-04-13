package output

import (
	"context"

	"sync"

	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/plug"
)

var (
	Plugin = make(map[string]Output)
)

type Output interface {
	plug.Parter
	DoOutput()
}

func WrapWithOutputLoop(ctx context.Context, wg *sync.WaitGroup, group string, coreFunc func(map[string]string) error) func() {
	return func() {
		wg.Add(1)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case message, isChnOpen := <-plugin.P2OChan[group]:
				switch isChnOpen {
				case true:
					err := coreFunc(message)
					if err != nil {
						plugin.Metrics.OutputErr++
						continue
					}

					plugin.Metrics.OutputOK++
				case false:
					plugin.OutputStatus.Completed = true
					return
				}
			}
		}
	}
}
