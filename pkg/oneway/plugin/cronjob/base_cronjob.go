package cronjob

import (
	"context"
	"os"
	"sync"

	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/plug"
	"github.com/robfig/cron"
)

const (
	module = "baseCronJob"
)

var (
	Plugin     = make(map[string]Cronjob)
	cronLogger = log.GetLogger(module)
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
		err := c.AddFunc(schedule, coreFunc)
		if err != nil {
			cronLogger.Error("failed to init the cronjob. please check the cronjob conf")
			os.Exit(1)
		}

		c.Start()

		for {
			<-ctx.Done()
			coreFunc()
			c.Stop()
			return
		}
	}
}
