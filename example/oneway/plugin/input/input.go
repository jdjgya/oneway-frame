package input

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/goinggo/mapstructure"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/input"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	module = "dummy"
)

type DummyInputter struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	input func()
	config

	log  *zap.Logger
	logf *zap.SugaredLogger
}

type config struct {
	Name          string `validate:"required"`
	Group         string `validate:"required"`
	FetchInterval int    `validate:"required"`
}

func init() {
	input.Plugin[module] = &DummyInputter{}
}

func (d *DummyInputter) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &d.config)

	d.wg = &sync.WaitGroup{}
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.input = input.WrapWithInputLoop(d.ctx, d.wg, d.Group, d.coreFunc, time.Duration(d.FetchInterval))

	d.log = log.GetLogger(module)
	d.logf = d.log.Sugar()

	fmt.Println("after setconfig")
}

func (d *DummyInputter) CheckConfig() error {
	validate := validator.New()
	return validate.Struct(d.config)
}

func (d *DummyInputter) coreFunc() ([]byte, error) {
	fmt.Println("dummy input demo")
	return []byte("dummy input demo"), nil
}

func (d *DummyInputter) DoInput() {
	d.input()
}

func (d *DummyInputter) Stop() {
	d.cancel()
	d.wg.Wait()
}
