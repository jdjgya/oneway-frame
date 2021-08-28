package worker

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/goinggo/mapstructure"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/cronjob"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/input"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/output"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/plug"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/process"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/transit"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

const (
	dummy                 = "dummy"
	doInput               = "doInput"
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
	case plugin.Input:
		return map[string]interface{}{"name": dummy}
	case plugin.Transit:
		return map[string]interface{}{"name": dummy}
	case plugin.Process:
		return map[string]interface{}{"name": dummy}
	case plugin.Output:
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

// input plugin for testing
type testInput struct {
	Status string
	Conf
}

func (ti *testInput) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &ti.Status)
}

func (ti *testInput) CheckConfig() error {
	validate := validator.New()
	err := validate.Struct(ti.Conf)
	if err != nil {
		return err
	}

	return nil
}

func (ti *testInput) DoInput() {
	ti.Status = doInput
}

func (ti *testInput) Stop() {
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

func (tt *testTransit) DoTransit() {
	tt.Status = doTransit
}

func (tt *testTransit) Stop() {
	tt.Status = onStop
}

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

func (tp *testProcess) DoProcess() {
	tp.Status = doProcess
}

func (tp *testProcess) Stop() {
	tp.Status = onStop
}

// output plugin for testing
type testOutput struct {
	Status string
	Conf
}

func (to *testOutput) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &to.Conf)
}

func (to *testOutput) CheckConfig() error {
	validate := validator.New()
	err := validate.Struct(to.Conf)
	if err != nil {
		return err
	}

	return nil
}

func (to *testOutput) DoOutput() {
	to.Status = doOutput
}

func (to *testOutput) Stop() {
	to.Status = onStop
}

// cronjob plugin for testing
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

func TestSetWorkerSuccess(t *testing.T) {
	configer = &testConfiger{}

	ti := &testInput{
		Conf: Conf{
			Name: dummy,
		},
	}
	input.Plugin[dummy] = ti
	tester.SetParter(plugin.Input)
	assert.Equal(t, dummy, ti.Name, "failed to set input plugin")

	tt := &testTransit{
		Conf: Conf{
			Name: dummy,
		},
	}

	transit.Plugin[dummy] = tt
	tester.SetParter(plugin.Transit)
	assert.Equal(t, dummy, tt.Name, "failed to set transit plugin")

	tp := &testProcess{
		Conf: Conf{
			Name: dummy,
		},
	}
	process.Plugin[dummy] = tp
	tester.SetParter(plugin.Process)
	assert.Equal(t, dummy, tp.Name, "failed to set process plugin")

	to := &testOutput{
		Conf: Conf{
			Name: dummy,
		},
	}
	output.Plugin[dummy] = to
	tester.SetParter(plugin.Output)
	assert.Equal(t, dummy, to.Name, "failed to set output plugin")
}

func TestSetCronjobSuccess(t *testing.T) {
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

func TestSetCronjobFailed(t *testing.T) {
	configer = &testConfiger{}

	parter := "cronjobs-dummy"
	plug.Cronners[parter] = &testCronjob{}
	testConf := make(map[string]map[interface{}]interface{})
	testConf[dummy] = make(map[interface{}]interface{})
	testConf[dummy] = map[interface{}]interface{}{}

	defer func() {
		osExit = os.Exit
	}()

	var testReturnCode int
	osExit = func(code int) {
		testReturnCode = code
	}

	logger := log.GetLogger(module)
	testOnewayer := &Onewayer{
		log:  logger,
		logf: logger.Sugar(),
	}
	testOnewayer.setCronnerConfig(plugin.CronJob, testConf)
	assert.Equal(t, 1, testReturnCode, "failed to prevent invalid cronjob conf")

	delete(plug.Cronners, parter)
	tester.SetCronner(plugin.CronJob)
	assert.Equal(t, 1, testReturnCode, "failed to prevent invalid cronjob plugin")
}

func TestDoWork(t *testing.T) {
	ti := &testInput{}
	input.Plugin[dummy] = ti

	tt := &testTransit{}
	transit.Plugin[dummy] = tt

	tp := &testProcess{}
	process.Plugin[dummy] = tp

	to := &testOutput{}
	output.Plugin[dummy] = to

	tester.StartParters()
	time.Sleep(1 * time.Second)

	assert.Equal(t, doInput, ti.Status, "failed to start input plugin")
	assert.Equal(t, doTransit, tt.Status, "failed to start transit plugin")
	assert.Equal(t, doProcess, tp.Status, "failed to start process plugin")
	assert.Equal(t, doOutput, to.Status, "failed to start output plugin")
}

func TestDoCron(t *testing.T) {
	tc := &testCronjob{}
	cronjob.Plugin[dummy] = tc

	tester.StartCronners()
	time.Sleep(1 * time.Second)

	assert.Equal(t, doCron, tc.Status, "failed to start cronjob plugin")
}

func TestStopWork(t *testing.T) {
	ti := &testInput{}
	input.Plugin[dummy] = ti

	tt := &testTransit{}
	transit.Plugin[dummy] = tt

	tp := &testProcess{}
	process.Plugin[dummy] = tp

	to := &testOutput{}
	output.Plugin[dummy] = to

	tester.StopParters()
	time.Sleep(1 * time.Second)

	assert.Equal(t, onStop, ti.Status, "failed to stop input plugin")
	assert.Equal(t, onStop, tt.Status, "failed to stop transit plugin")
	assert.Equal(t, onStop, tp.Status, "failed to stop process plugin")
	assert.Equal(t, onStop, to.Status, "failed to stop output plugin")
}

func TestStopCron(t *testing.T) {
	tc := &testCronjob{}
	cronjob.Plugin[dummy] = tc

	tester.StopCronners()
	time.Sleep(1 * time.Second)

	assert.Equal(t, onStop, tc.Status, "failed to stop cronjob plugin")
}

func TestGetStatusTrue(t *testing.T) {
	plugin.InputStatus.Completed = true
	plugin.TransitStatus.Completed = true
	plugin.ProcessStatus.Completed = true
	plugin.OutputStatus.Completed = true

	assert.Equal(t, true, tester.GetStatus(), "failed to get the false status")
}

func TestGetStatusFalse(t *testing.T) {
	plugin.InputStatus.Completed = false
	plugin.TransitStatus.Completed = true
	plugin.ProcessStatus.Completed = true
	plugin.OutputStatus.Completed = true

	assert.Equal(t, false, tester.GetStatus(), "failed to get the true status")
}
