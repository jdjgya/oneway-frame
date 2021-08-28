package input

import (
	"context"
	"onewayframe/pkg/plugin"
	"onewayframe/pkg/plugin/plug"
	"sync"
	"time"
)

var (
	Plugin = make(map[string]Input)
)

type Input interface {
	plug.Parter
	DoInput()
}

func WrapWithInputLoop(ctx context.Context, wg *sync.WaitGroup, coreFunc func() ([]byte, error), interval time.Duration) func() {
	return func() {
		wg.Add(1)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := coreFunc()
				if err != nil {
					plugin.Metrics.InputErr++
					continue
				}

				plugin.I2TChan <- msg
				plugin.Metrics.InputOK++

				if plugin.IsOneTimeExec {
					close(plugin.I2TChan)
					plugin.InputStatus.Completed = true
					return
				}

				time.Sleep(interval * time.Second)
			}
		}
	}
}
