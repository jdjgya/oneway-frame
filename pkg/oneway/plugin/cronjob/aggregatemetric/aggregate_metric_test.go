package aggregatemetric

import (
	"context"

	"sync"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/cronjob"
)

var (
	schedule    = "* * * * * *"
	isStopped   = false
	didSchedule = false
)

// wrapper for testing
func testLoopWrapper(ctx context.Context, wg *sync.WaitGroup, schedule string, coreFunc func()) func() {
	return func() {
		for {
			select {
			case <-ctx.Done():
				isStopped = true
				return
			default:
				didSchedule = true
			}
		}
	}
}

func TestInitRestAPI(t *testing.T) {
	assert.NotEqual(t, nil, cronjob.Plugin[module], "failed to init aggregate-metric into cronjob.Plugin")
}

func TestSetConfigSuccess(t *testing.T) {
	tester := &MetricAggregator{
		wg:  &sync.WaitGroup{},
		log: log.GetLogger(module),
	}

	testConf := map[interface{}]interface{}{
		"name":     module,
		"schedule": schedule,
	}

	tester.SetConfig(testConf)
	err := tester.CheckConfig()
	assert.Equal(t, nil, err, "failed to set valid config")
	assert.Equal(t, module, tester.Name, "failed to set conf name")
	assert.Equal(t, schedule, tester.Schedule, "failed to set schedule")
}

func TestSetConfigFailed(t *testing.T) {
	tester := &MetricAggregator{
		wg:  &sync.WaitGroup{},
		log: log.GetLogger(module),
	}

	testConf := map[interface{}]interface{}{
		"name": module,
	}

	tester.SetConfig(testConf)
	err := tester.CheckConfig()
	assert.Equal(t, "Key: 'config.Schedule' Error:Field validation for 'Schedule' failed on the 'required' tag", err.Error(), "failed to capture invalid config")
}

func TestCoreFunc(t *testing.T) {
	tester := &MetricAggregator{
		wg:  &sync.WaitGroup{},
		log: log.GetLogger(module),
	}

	plugin.Metrics = &plugin.Metric{}
	plugin.Metrics.InputOK++
	tester.coreFunc()
	assert.Equal(t, int32(0), plugin.Metrics.InputOK, "failed to swap new and old metrics pointer")
}

func TestDoSchedule(t *testing.T) {
	tester := &MetricAggregator{
		wg:  &sync.WaitGroup{},
		log: log.GetLogger(module),
	}
	tester.ctx, tester.cancel = context.WithCancel(context.Background())

	go func() {
		time.Sleep(2 * time.Second)
		tester.cancel()
	}()

	didSchedule = false
	tester.schedule = testLoopWrapper(tester.ctx, tester.wg, schedule, tester.coreFunc)
	tester.DoSchedule()
	assert.Equal(t, true, didSchedule, "failed to do schedule")
}

func TestStop(t *testing.T) {
	tester := &MetricAggregator{
		wg:  &sync.WaitGroup{},
		log: log.GetLogger(module),
	}
	tester.ctx, tester.cancel = context.WithCancel(context.Background())

	isStopped = false
	tester.schedule = testLoopWrapper(tester.ctx, tester.wg, schedule, tester.coreFunc)
	go func() {
		time.Sleep(2 * time.Second)
		tester.Stop()
	}()

	tester.DoSchedule()
	assert.Equal(t, true, isStopped, "failed to do stop")
}
