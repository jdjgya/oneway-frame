package process

import (
	"strings"
	"sync/atomic"

	"github.com/goinggo/mapstructure"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/sync/plugin"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/process"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	module = "dummy"
)

type DummyProcessor struct {
	strBuilder strings.Builder
	validator  *validator.Validate

	config

	log  *zap.Logger
	logf *zap.SugaredLogger
}

type config struct {
	Name string
}

func init() {
	process.Plugins[module] = &DummyProcessor{
		validator: validator.New(),
	}
}

func (d *DummyProcessor) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &d.config)

	d.log = log.GetLogger(module)
	d.logf = d.log.Sugar()
}

func (d *DummyProcessor) CheckConfig() error {
	return d.validator.Struct(d.config)
}

func (d *DummyProcessor) Execute(job *map[string]string) error {
	return nil
}

func (d *DummyProcessor) AddSuccess() {
	atomic.AddUint64(&plugin.Metrics.ProcessOK, 1)
}

func (d *DummyProcessor) AddError() {
	atomic.AddUint64(&plugin.Metrics.ProcessErr, 1)
}
