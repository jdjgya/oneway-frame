package output

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/goinggo/mapstructure"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/output"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	module = "dummy"
)

type DummyOutputer struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	output func()
	config

	log  *zap.Logger
	logf *zap.SugaredLogger
}

type config struct {
	Name  string `validate:"required"`
	Group string `validate:"required"`
}

func init() {
	output.Plugin[module] = &DummyOutputer{}
}

func (d *DummyOutputer) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &d.config)

	d.wg = &sync.WaitGroup{}

	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.output = output.WrapWithOutputLoop(d.ctx, d.wg, d.Group, d.coreFunc)

	d.log = log.GetLogger(module)
	d.logf = d.log.Sugar()
}

func (d *DummyOutputer) CheckConfig() error {
	validate := validator.New()
	return validate.Struct(d.config)
}

func (d *DummyOutputer) coreFunc(msg map[string]string) error {
	b, _ := json.Marshal(msg)
	d.log.Info(string(b))
	return nil
}

func (d *DummyOutputer) DoOutput() {
	d.output()
}

func (d *DummyOutputer) Stop() {
	d.cancel()
	d.wg.Wait()
}
