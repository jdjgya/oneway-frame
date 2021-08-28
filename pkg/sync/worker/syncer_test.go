package worker

import (
	"github.com/jdjgya/service-frame/pkg/sync/plugin"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/cronjob"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/process"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/request"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/transit"

	"errors"
	"io"
	"testing"
	"time"

	"github.com/goinggo/mapstructure"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

const (
	dummy                 = "dummy"
	doInteract            = "doInteract"
	doTransit             = "doTransit"
	doProcess             = "doProcess"
	doOutput              = "doOutput"
	doCron                = "doCron"
	nonExistentPluginType = "nonExistentPluginType"
	onStop                = "onStop"
)

// worker for testing
var tester Worker

// configer for testing
type testConfiger struct{}

var testConfer *testConfiger

func (tcfgr *testConfiger) SetConfigType(string) {}

func (tcfgr *testConfiger) ReadConfig(ir io.Reader) error { return errors.New("use for error cases") }

func (tcfgr *testConfiger) Get(pluginType string) interface{} {
	switch pluginType {
	case plugin.Interact:
		return map[string]interface{}{
			"name": dummy,
			"transit": map[string]interface{}{
				"name": dummy,
			},
			"process": map[string]interface{}{
				"name": dummy,
			},
			"request": map[string]interface{}{
				"name": dummy,
			},
		}
	case plugin.StageTransit:
		return map[string]interface{}{"name": dummy}
	case plugin.StageProcess:
		return map[string]interface{}{"name": dummy}
	case plugin.StageRequest:
		return map[string]interface{}{"name": dummy}
	case plugin.CronJob:
		return []map[interface{}]interface{}{{"name": dummy}}
	default:
		return map[string]interface{}{}
	}
}

func (tcfgr *testConfiger) GetInt32(confKey string) int32 { return 0 }

// general conf struct for testing and for all plugins
type Conf struct {
	Name string `validate:"required"`
}

// interact plugin for testing
type testInteract struct {
	Status string
	Conf
}

func (ti *testInteract) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &ti.Status)
}

func (ti *testInteract) CheckConfig() error {
	validate := validator.New()
	err := validate.Struct(ti.Conf)
	if err != nil {
		return err
	}

	return nil
}

func (ti *testInteract) DoInteract() {
	ti.Status = doInteract
}

func (ti *testInteract) Stop() {
	ti.Status = onStop
}

// transit plugin for testing
type testTransit struct {
	Status string
	Conf
}

func (tt *testTransit) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &tt.Conf)
}

func (tt *testTransit) CheckConfig() error {
	validate := validator.New()
	err := validate.Struct(tt.Conf)
	if err != nil {
		return err
	}

	return nil
}

func (tt *testTransit) Execute(job *map[string]string) error {
	tt.Status = doTransit
	return nil
}

func (tt *testTransit) AddSuccess() {}

func (tt *testTransit) AddError() {}

// process plugin for testing
type testProcess struct {
	Status string
	Conf
}

func (tp *testProcess) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &tp.Conf)
}

func (tp *testProcess) CheckConfig() error {
	validate := validator.New()
	err := validate.Struct(tp.Conf)
	if err != nil {
		return err
	}

	return nil
}

func (tp *testProcess) Execute(job *map[string]string) error {
	tp.Status = doProcess
	return nil
}

func (tp *testProcess) AddSuccess() {}

func (tp *testProcess) AddError() {}

// Request plugin for testing
type testRequest struct {
	Status string
	Conf
}

func (tr *testRequest) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &tr.Conf)
}

func (tr *testRequest) CheckConfig() error {
	validate := validator.New()
	err := validate.Struct(tr.Conf)
	if err != nil {
		return err
	}

	return nil
}

func (tr *testRequest) Execute(job *map[string]string) error {
	tr.Status = doOutput
	return nil
}

func (tr *testRequest) AddSuccess() {}

func (tr *testRequest) AddError() {}

type testCronjob struct {
	Status string
	Conf
}

func (tc *testCronjob) SetConfig(conf map[interface{}]interface{}) {
	_ = mapstructure.Decode(conf, &tc.Conf)
}

func (tc *testCronjob) CheckConfig() error {
	validate := validator.New()
	err := validate.Struct(tc.Conf)
	if err != nil {
		return err
	}

	return nil
}

func (tc *testCronjob) DoSchedule() {
	tc.Status = doCron
}

func (tc *testCronjob) Stop() {
	tc.Status = onStop
}

func TestInitWorker(t *testing.T) {
	tester = InitWorker()
	assert.NotEqual(t, nil, tester, "failed to init worker instance")
}

func TestSetParterSuccess(t *testing.T) {
	configer = &testConfiger{}

	ti := &testInteract{
		Conf: Conf{
			Name: dummy,
		},
	}

	interact.Plugins[dummy] = ti
	transit.Plugins[dummy] = &testTransit{}
	process.Plugins[dummy] = &testProcess{}
	request.Plugins[dummy] = &testRequest{}

	tester.SetParter(plugin.Interact)
	assert.Equal(t, dummy, ti.Name, "failed to set interact plugin")
	assert.Equal(t, dummy, plugin.ActivatedTransit, "failed to set transit plugin")
	assert.Equal(t, dummy, plugin.ActivatedProcess, "failed to set process plugin")
	assert.Equal(t, dummy, plugin.ActivatedRequest, "failed to set request plugin")
}

func TestStartParters(t *testing.T) {
	configer = &testConfiger{}

	ti := &testInteract{
		Conf: Conf{
			Name: dummy,
		},
	}

	interact.Plugins[dummy] = ti
	transit.Plugins[dummy] = &testTransit{}
	process.Plugins[dummy] = &testProcess{}
	request.Plugins[dummy] = &testRequest{}

	tester.SetParter(plugin.Interact)
	tester.StartParters()
	time.Sleep(1 * time.Second)
	assert.Equal(t, doInteract, ti.Status, "failed to start parter")
}

func TestStopParters(t *testing.T) {
	configer = &testConfiger{}

	ti := &testInteract{
		Conf: Conf{
			Name: dummy,
		},
	}

	interact.Plugins[dummy] = ti
	transit.Plugins[dummy] = &testTransit{}
	process.Plugins[dummy] = &testProcess{}
	request.Plugins[dummy] = &testRequest{}

	tester.SetParter(plugin.Interact)
	tester.StopParters()
	time.Sleep(1 * time.Second)
	assert.Equal(t, onStop, ti.Status, "failed to stop parter")
}

func TestSetCronnerSuccess(t *testing.T) {
	configer = &testConfiger{}

	tc := &testCronjob{
		Conf: Conf{
			Name: dummy,
		},
	}

	cronjob.Plugin[dummy] = tc
	tester.SetCronner(plugin.CronJob)
	assert.Equal(t, dummy, tc.Name, "failed to set cronjob plugin")
}

func TestStartCronnerSuccess(t *testing.T) {
	configer = &testConfiger{}

	tc := &testCronjob{
		Conf: Conf{
			Name: dummy,
		},
	}

	cronjob.Plugin[dummy] = tc
	tester.SetCronner(plugin.CronJob)
	tester.StartCronners()
	time.Sleep(1 * time.Second)
	assert.Equal(t, doCron, tc.Status, "failed to start cronjob plugin")
}

func TestStopCronnerSuccess(t *testing.T) {
	configer = &testConfiger{}

	tc := &testCronjob{
		Conf: Conf{
			Name: dummy,
		},
	}

	cronjob.Plugin[dummy] = tc
	tester.SetCronner(plugin.CronJob)
	tester.StopCronners()
	time.Sleep(1 * time.Second)
	assert.Equal(t, onStop, tc.Status, "failed to stop cronjob plugin")
}
