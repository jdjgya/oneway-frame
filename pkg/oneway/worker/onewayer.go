package worker

import (
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/jdjgya/service-frame/pkg/config"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/cronjob"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/input"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/output"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/plug"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/process"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin/transit"
	"go.uber.org/zap"
)

const (
	name   = "name"
	module = "forwarder"
)

var (
	once     sync.Once
	configer = config.GetConfiger()
	osExit   = os.Exit
)

type Forwarder struct {
	Input    string
	Transit  string
	Process  string
	Output   string
	CronJobs []string
	log      *zap.Logger
	logf     *zap.SugaredLogger
}

func InitWorker() Worker {
	logger := log.GetLogger(module)

	return &Forwarder{
		log:  logger,
		logf: logger.Sugar(),
	}
}

func (f *Forwarder) getParterConfig(pluginType string) (string, interface{}) {
	rawConfig := configer.Get(pluginType)
	pluginName := (rawConfig.(map[string]interface{}))[name].(string)
	parterName := strings.Join([]string{pluginType, pluginName}, "-")

	switch pluginType {
	case plugin.Input:
		f.Input = pluginName
		plug.Parters[parterName] = input.Plugin[pluginName]
	case plugin.Transit:
		f.Transit = pluginName
		plug.Parters[parterName] = transit.Plugin[pluginName]
	case plugin.Process:
		f.Process = pluginName
		plug.Parters[parterName] = process.Plugin[pluginName]
	case plugin.Output:
		f.Output = pluginName
		plug.Parters[parterName] = output.Plugin[pluginName]
	}

	return parterName, rawConfig
}

func (f *Forwarder) setParterConfig(parterName string, parterConfig interface{}) {
	plug.Parters[parterName].SetConfig(parterConfig)
	err := plug.Parters[parterName].CheckConfig()
	if err != nil {
		f.logf.Errorf("failed to set parter config. error details: %s", err.Error())
		osExit(1)
	}
}

func (f *Forwarder) SetParter(pluginType string) {
	parterName, parterConfig := f.getParterConfig(pluginType)
	f.setParterConfig(parterName, parterConfig)
}

func (f *Forwarder) StartParters() {
	f.logf.Infof("start input plugin(%s)", f.Input)
	go input.Plugin[f.Input].DoInput()

	f.logf.Infof("start transit plugin(%s)", f.Transit)
	go transit.Plugin[f.Transit].DoTransit()

	f.logf.Infof("start process plugin(%s)", f.Process)
	go process.Plugin[f.Process].DoProcess()

	f.logf.Infof("start output plugin(%s)", f.Output)
	go output.Plugin[f.Output].DoOutput()
}

func (f *Forwarder) StopParters() {
	input.Plugin[f.Input].Stop()
	f.logf.Infof("stop input plugin(%s)", f.Input)

	transit.Plugin[f.Transit].Stop()
	f.logf.Infof("stop transit plugin(%s)", f.Transit)

	process.Plugin[f.Process].Stop()
	f.logf.Infof("stop process plugin(%s)", f.Process)

	output.Plugin[f.Output].Stop()
	f.logf.Infof("stop output plugin(%s)", f.Output)
}

func (f *Forwarder) getCronnerConfig(pluginType string) map[string]map[interface{}]interface{} {
	cronConfigs := make(map[string]map[interface{}]interface{})
	rawConfig := reflect.ValueOf(configer.Get(pluginType))

	for i := 0; i < rawConfig.Len(); i++ {
		pluginConfig := rawConfig.Index(i).Interface().(map[interface{}]interface{})
		pluginName := pluginConfig[name].(string)
		cronName := strings.Join([]string{pluginType, pluginName}, "-")

		cronConfigs[pluginName] = make(map[interface{}]interface{})
		cronConfigs[pluginName] = pluginConfig
		plug.Cronners[cronName] = cronjob.Plugin[pluginName]
	}

	return cronConfigs
}

func (f *Forwarder) setCronnerConfig(pluginType string, cronConfigs map[string]map[interface{}]interface{}) {
	var cronName string
	defer func() {
		if err := recover(); err != nil {
			f.logf.Errorf("%s cronjob plugin(%s) was not defined", pluginType, cronName)
			osExit(1)
		}
	}()

	for cronName, cronConfig := range cronConfigs {
		cronner := strings.Join([]string{pluginType, cronName}, "-")
		plug.Cronners[cronner].SetConfig(cronConfig)
		err := plug.Cronners[cronner].CheckConfig()
		if err != nil {
			f.logf.Errorf("failed to set cronjob plugin(%s). error: %s", cronName, err.Error())
			osExit(1)
		}

		f.CronJobs = append(f.CronJobs, cronName)
	}
}

func (f *Forwarder) SetCronner(pluginType string) {
	cronConfigs := f.getCronnerConfig(pluginType)
	f.setCronnerConfig(pluginType, cronConfigs)
}

func (f *Forwarder) StartCronners() {
	for _, cronName := range f.CronJobs {
		f.logf.Infof("start cronjob plugin(%s)", cronName)
		go cronjob.Plugin[cronName].DoSchedule()
	}
}

func (f *Forwarder) StopCronners() {
	for _, cronName := range f.CronJobs {
		cronjob.Plugin[cronName].Stop()
		f.logf.Infof("stop cronjob plugin(%s)", cronName)
	}
}

func (f *Forwarder) GetStatus() bool {
	areAllPluginsCompleted := plugin.InputStatus.Completed && plugin.TransitStatus.Completed && plugin.ProcessStatus.Completed && plugin.OutputStatus.Completed
	return areAllPluginsCompleted
}
