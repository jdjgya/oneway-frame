package cronjob

import (
	"context"
	"sync"

	"github.com/jdjgya/service-frame/pkg/oneway/plugin/plug"
	"github.com/robfig/cron"
)

var (
	Plugin = make(map[string]Cronjob)
)

type Cronjob interface {
	plug.Cronner
	DoSchedule()
}

func WrapWithCron(ctx context.Context, wg *sync.WaitGroup, schedule string, coreFunc func()) func() {
	return func() {
		wg.Add(1)
		defer wg.Done()

		c := cron.New()
		c.AddFunc(schedule, coreFunc)

		c.Start()

		for {
			select {
			case <-ctx.Done():
				coreFunc()
				c.Stop()
				return
			}
		}
	}
}
