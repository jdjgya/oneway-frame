package transit

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/goinggo/mapstructure"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/transit"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	module = "dummy"
)

type DummyTransitter struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	transit func()
	config

	log  *zap.Logger
	logf *zap.SugaredLogger
}

type config struct {
	Name  string `validate:"required"`
	Group string `validate:"required"`
}

func init() {
	transit.Plugin[module] = &DummyTransitter{}
}

func (d *DummyTransitter) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &d.config)

	d.wg = &sync.WaitGroup{}
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.transit = transit.WrapWithTransitLoop(d.ctx, d.wg, d.Group, d.coreFunc)

	d.log = log.GetLogger(module)
	d.logf = d.log.Sugar()
}

func (d *DummyTransitter) CheckConfig() error {
	validator := validator.New()
	return validator.Struct(d.config)
}

func (d *DummyTransitter) coreFunc(msgs []byte) ([]map[string]string, error) {
	msg := make(map[string]string)
	msg["dummy"] = string(msgs)

	var msgList []map[string]string
	msgList = append(msgList, msg)

	b, _ := json.Marshal(msgList)
	d.log.Info(string(b))

	return msgList, nil
}

func (d *DummyTransitter) DoTransit() {
	d.transit()
}

func (d *DummyTransitter) Stop() {
	d.cancel()
	d.wg.Wait()
}
