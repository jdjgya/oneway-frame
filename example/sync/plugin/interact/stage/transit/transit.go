package transit

import (
	"fmt"
	"sync/atomic"

	"github.com/goinggo/mapstructure"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/sync/plugin"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/transit"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	module = "dummy"
)

type DummyTransiter struct {
	validator *validator.Validate
	config

	log  *zap.Logger
	logf *zap.SugaredLogger
}

type config struct {
	Name string
	Echo string
}

func init() {
	transit.Plugins[module] = &DummyTransiter{}
}

func (d *DummyTransiter) New() transit.Transit {
	return &DummyTransiter{
		validator: validator.New(),
	}
}

func (d *DummyTransiter) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &d.config)
	fmt.Println("init echo:", d.config.Echo)

	d.log = log.GetLogger(module)
	d.logf = d.log.Sugar()
}

func (d *DummyTransiter) CheckConfig() error {
	return d.validator.Struct(d.config)
}

func (d *DummyTransiter) Execute(job *map[string]string) error {
	return nil
}

func (d *DummyTransiter) AddSuccess() {
	atomic.AddUint64(&plugin.Metrics.TransitOK, 1)
}

func (d *DummyTransiter) AddError() {
	atomic.AddUint64(&plugin.Metrics.TransitErr, 1)
}
