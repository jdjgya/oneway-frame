package process

import (
	"context"
	"sync"

	"github.com/goinggo/mapstructure"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/process"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	module = "dummy"
)

type DummyProcessor struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	process func()
	config

	log  *zap.Logger
	logf *zap.SugaredLogger
}

type config struct {
	Name string `validate:"required"`
}

func init() {
	process.Plugin[module] = &DummyProcessor{
		wg: &sync.WaitGroup{},
	}
}

func (d *DummyProcessor) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &d.config)
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.process = process.WrapWithProcessLoop(d.ctx, d.wg, d.coreFunc)

	d.log = log.GetLogger(module)
	d.logf = d.log.Sugar()
}

func (d *DummyProcessor) CheckConfig() error {
	validate := validator.New()
	return validate.Struct(d.config)
}

func (d *DummyProcessor) coreFunc(msg map[string]string, isChnOpen bool) (map[string]string, error) {
	d.log.Info("dummy process demo")
	return msg, nil
}

func (d *DummyProcessor) DoProcess() {
	d.process()
}

func (d *DummyProcessor) Stop() {
	d.cancel()
	d.wg.Wait()
}
