package aggregatemetric

import (
	"context"
	"encoding/json"
	"onewayframe/pkg/log"
	"onewayframe/pkg/plugin"
	"onewayframe/pkg/plugin/cronjob"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/goinggo/mapstructure"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	module = "aggregate-metric"
)

type MetricAggregator struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	schedule func()
	config

	log  *zap.Logger
	logf *zap.SugaredLogger
}

type config struct {
	Name     string `validate:"required"`
	Schedule string `validate:"required"`
}

func init() {
	aggregateMetric := &MetricAggregator{
		wg: &sync.WaitGroup{},
	}

	cronjob.Plugin[module] = aggregateMetric
}

func (m *MetricAggregator) SetConfig(conf map[interface{}]interface{}) {
	_ = mapstructure.Decode(conf, &m.config)
	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.schedule = cronjob.WrapWithCron(m.ctx, m.wg, m.Schedule, m.coreFunc)

	m.log = log.GetLogger(module)
	m.logf = m.log.Sugar()
}

func (m *MetricAggregator) CheckConfig() error {
	validate := validator.New()
	return validate.Struct(m.config)
}

func (m *MetricAggregator) swapRecords(metrics **plugin.Metric) *plugin.Metric {
	oldMetrics := &plugin.Metric{}
	newMetrics := &plugin.Metric{}

	*(*unsafe.Pointer)(unsafe.Pointer(&oldMetrics)) = atomic.SwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(metrics)),
		*(*unsafe.Pointer)(unsafe.Pointer(&newMetrics)),
	)

	return oldMetrics
}

func (m *MetricAggregator) coreFunc() {
	metrics := m.swapRecords(&plugin.Metrics)
	jsonMetrics, _ := json.Marshal(metrics)
	raw := json.RawMessage(jsonMetrics)

	m.log.Info("periodical metrics dump", zap.Any("metrics", &raw))
}

func (m *MetricAggregator) DoSchedule() {
	m.schedule()
}

func (m *MetricAggregator) Stop() {
	m.cancel()
	m.wg.Wait()
	m.ctx, m.cancel = context.WithCancel(context.Background())
}
