package controller

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/jdjgya/service-frame/pkg/config"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/sync/plugin"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	interact = "testInteract"
	transit  = "testTransit"
	process  = "testProcess"
	request  = "testOutput"
	cronJob  = "cronJob"

	successConf            = "../../../example/sync/sync-conf.yaml"
	failureNotFoundConf    = "can_not_find_this_file.yaml"
	failureReadConfContent = "../../../example/sync/sync-conf.yaml"

	start = "start"
	stop  = "stop"
)

var (
	isTestWorkerCompleted bool
	_                     = func() bool {
		testing.Init()
		return true
	}()
)

type testWorker struct {
	Interact         string
	Transit          string
	Process          string
	Request          string
	CronJobs         []string
	log              *zap.Logger
	testWokrerStatus string
	testCronJbStatus string
}

var tester *testWorker

func (t *testWorker) SetStager(stageType string, stageConfig interface{}) {}

func (t *testWorker) SetParter(pluginType string) {
	switch pluginType {
	case plugin.Interact:
		t.Interact = interact
	case plugin.StageTransit:
		t.Transit = transit
	case plugin.StageProcess:
		t.Process = process
	case plugin.StageRequest:
		t.Request = request
	}
}

func (t *testWorker) StartParters() { t.testWokrerStatus = start }

func (t *testWorker) StopParters() { t.testWokrerStatus = stop }

func (t *testWorker) SetCronner(pluginType string) {
	t.CronJobs = append(t.CronJobs, cronJob)
}

func (t *testWorker) StartCronners() { t.testCronJbStatus = start }

func (t *testWorker) StopCronners() { t.testCronJbStatus = stop }

type testConfiger struct{}

var testConfer *testConfiger

func (tcfgr *testConfiger) SetConfigType(string) {}

func (tcfgr *testConfiger) ReadConfig(ir io.Reader) error { return errors.New("use for error cases") }

func (tcfgr *testConfiger) Get(confKey string) interface{} { return confKey }

func (tcfgr *testConfiger) GetInt32(confKey string) int32 { return 0 }

func TestController(t *testing.T) {
	GetInstance()

	assert.NotEqual(t, nil, instance, "failed to get controller instance")
	assert.NotEqual(t, nil, instance.Worker, "failed to get worker instance")
}

func TestInitService(t *testing.T) {
	var testReturnCode int
	osExit = func(code int) {
		testReturnCode = code
	}

	defer func() {
		osExit = os.Exit
		conf = successConf
	}()

	conf = successConf
	instance.InitService()
	assert.NotEmpty(t, instance.log, "failed to init logger")
	assert.NotEqual(t, nil, instance.ctx, "failed to init ctx")
	assert.NotEqual(t, nil, instance.wg, "failed to init wg")
	assert.NotEqual(t, nil, plugin.Metrics, "failed to init plugin metrics")

	conf = ""
	instance.InitService()
	assert.Equal(t, 1, testReturnCode, "failed to get error return code")
}

func TestLoadConfFile(t *testing.T) {
	var err error
	configer = &testConfiger{}
	defer func() { configer = config.GetConfiger() }()

	err = instance.loadConfig("")
	assert.Equal(t, "conf file is required, please specify the path of conf file", err.Error(), "failed to get conf error")

	err = instance.loadConfig(failureNotFoundConf)
	assert.NotEqual(t, nil, err, "failed to get conf error")

	err = instance.loadConfig(failureReadConfContent)
	assert.Equal(t, "use for error cases", err.Error(), "failed to get conf error")
}

func TestActivateService(t *testing.T) {
	tester = &testWorker{}
	instance.Worker = tester

	instance.ActivateService()
	assert.Equal(t, interact, tester.Interact, "failed to set interact plugin")
}

func TestStartService(t *testing.T) {
	tester = &testWorker{}
	instance.Worker = tester

	instance.Start()
	assert.Equal(t, start, tester.testWokrerStatus, "failed to start worker")
	assert.Equal(t, start, tester.testCronJbStatus, "failed to start cronjob")
	instance.wg.Done()
}

func TestRestartService(t *testing.T) {
	tester = &testWorker{}
	instance.Worker = tester
	instance.wg.Add(1)

	instance.log = log.GetLogger(module)
	instance.ctx, instance.cancel = context.WithCancel(context.Background())
	instance.wg = &wg
	instance.signalChannel = make(chan os.Signal, 1)
	tester.testWokrerStatus = stop
	tester.testCronJbStatus = stop

	instance.Restart()
	assert.Equal(t, start, tester.testWokrerStatus, "failed to restart worker")
	assert.Equal(t, start, tester.testCronJbStatus, "failed to restart cronjob")
}

func TestStopService(t *testing.T) {
	tester = &testWorker{}
	instance.Worker = tester
	tester.testWokrerStatus = start
	tester.testCronJbStatus = start
	instance.wg.Add(1)

	instance.Stop()
	assert.Equal(t, stop, tester.testWokrerStatus, "failed to stop worker")
	assert.Equal(t, stop, tester.testCronJbStatus, "failed to stop cronjob")
}

func TestTrapSignalsSIGTERM(t *testing.T) {
	tester = &testWorker{}
	instance.Worker = tester
	tester.testWokrerStatus = start
	tester.testCronJbStatus = start

	instance.TrapSignals()
	instance.signalChannel <- syscall.SIGTERM

	time.Sleep(2 * time.Second)
	assert.Equal(t, stop, tester.testWokrerStatus, "failed to stop worker by SIGTERM")
	assert.Equal(t, stop, tester.testCronJbStatus, "failed to stop cronjob by SIGTERM")
}

func TestTrapSignalsSIGHUP(t *testing.T) {
	tester = &testWorker{}
	instance.Worker = tester
	tester.testWokrerStatus = stop
	tester.testCronJbStatus = stop

	instance.TrapSignals()
	instance.signalChannel <- syscall.SIGHUP

	time.Sleep(2 * time.Second)
	assert.Equal(t, start, tester.testWokrerStatus, "failed to stop worker by SIGHUP")
	assert.Equal(t, start, tester.testCronJbStatus, "failed to stop cronjob by SIGHUP")
}

func TestTraceStatus(t *testing.T) {
	tester = &testWorker{}
	instance.Worker = tester

	instance.wg = &sync.WaitGroup{}
	instance.wg.Add(1)

	go func() {
		time.Sleep(2 * time.Second)
		instance.wg.Done()
	}()

	instance.TraceStatus()
}
