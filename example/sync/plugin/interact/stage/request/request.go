package request

import (
	"sync/atomic"

	"github.com/goinggo/mapstructure"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/sync/plugin"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/request"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	module = "dummy"
)

type DummyRequester struct {
	config

	log  *zap.Logger
	logf *zap.SugaredLogger
}

type config struct {
	Name string
}

func init() {
	request.Plugins[module] = &DummyRequester{}
}

func (d *DummyRequester) New() request.Request {
	return &DummyRequester{}
}

func (d *DummyRequester) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &d.config)

	d.log = log.GetLogger(module)
	d.logf = d.log.Sugar()
}

func (d *DummyRequester) CheckConfig() error {
	validate := validator.New()
	return validate.Struct(d.config)
}

func (d *DummyRequester) Execute(job *map[string]string) (bool, error) {
	return true, nil
}

func (d *DummyRequester) AddSuccess() {
	atomic.AddUint64(&plugin.Metrics.RequestOK, 1)
}

func (d *DummyRequester) AddError() {
	atomic.AddUint64(&plugin.Metrics.RequestErr, 1)
}
