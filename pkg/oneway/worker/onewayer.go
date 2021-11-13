package worker

import (
	"os"
	"reflect"
	"strings"

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
	module = "onewayer"
)

var (
	configer = config.GetConfiger()
	osExit   = os.Exit
)

type Onewayer struct {
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

	return &Onewayer{
		log:  logger,
		logf: logger.Sugar(),
	}
}

func (o *Onewayer) getParterConfig(pluginType string) (string, interface{}) {
	rawConfig := configer.Get(pluginType)
	pluginName := (rawConfig.(map[string]interface{}))[name].(string)
	parterName := strings.Join([]string{pluginType, pluginName}, "-")

	switch pluginType {
	case plugin.Input:
		o.Input = pluginName
		plug.Parters[parterName] = input.Plugin[pluginName]
	case plugin.Transit:
		o.Transit = pluginName
		plug.Parters[parterName] = transit.Plugin[pluginName]
	case plugin.Process:
		o.Process = pluginName
		plug.Parters[parterName] = process.Plugin[pluginName]
	case plugin.Output:
		o.Output = pluginName
		plug.Parters[parterName] = output.Plugin[pluginName]
	}

	return parterName, rawConfig
}

func (o *Onewayer) setParterConfig(parterName string, parterConfig interface{}) {
	plug.Parters[parterName].SetConfig(parterConfig)
	err := plug.Parters[parterName].CheckConfig()
	if err != nil {
		o.logf.Errorf("failed to set parter config. error details: %s", err.Error())
		osExit(1)
	}
}

func (o *Onewayer) SetParter(pluginType string) {
	parterName, parterConfig := o.getParterConfig(pluginType)
	o.setParterConfig(parterName, parterConfig)
}

func (o *Onewayer) StartParters() {
	o.logf.Infof("start input plugin(%s)", o.Input)
	go input.Plugin[o.Input].DoInput()

	o.logf.Infof("start transit plugin(%s)", o.Transit)
	go transit.Plugin[o.Transit].DoTransit()

	o.logf.Infof("start process plugin(%s)", o.Process)
	go process.Plugin[o.Process].DoProcess()

	o.logf.Infof("start output plugin(%s)", o.Output)
	go output.Plugin[o.Output].DoOutput()
}

func (o *Onewayer) StopParters() {
	input.Plugin[o.Input].Stop()
	o.logf.Infof("stop input plugin(%s)", o.Input)

	transit.Plugin[o.Transit].Stop()
	o.logf.Infof("stop transit plugin(%s)", o.Transit)

	process.Plugin[o.Process].Stop()
	o.logf.Infof("stop process plugin(%s)", o.Process)

	output.Plugin[o.Output].Stop()
	o.logf.Infof("stop output plugin(%s)", o.Output)
}

func (o *Onewayer) getCronnerConfig(pluginType string) map[string]map[interface{}]interface{} {
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

func (o *Onewayer) setCronnerConfig(pluginType string, cronConfigs map[string]map[interface{}]interface{}) {
	var cronName string
	defer func() {
		if err := recover(); err != nil {
			o.logf.Errorf("%s cronjob plugin(%s) was not defined", pluginType, cronName)
			osExit(1)
		}
	}()

	for cronName, cronConfig := range cronConfigs {
		cronner := strings.Join([]string{pluginType, cronName}, "-")
		plug.Cronners[cronner].SetConfig(cronConfig)
		err := plug.Cronners[cronner].CheckConfig()
		if err != nil {
			o.logf.Errorf("failed to set cronjob plugin(%s). error: %s", cronName, err.Error())
			osExit(1)
		}

		o.CronJobs = append(o.CronJobs, cronName)
	}
}

func (o *Onewayer) SetCronner(pluginType string) {
	cronConfigs := o.getCronnerConfig(pluginType)
	o.setCronnerConfig(pluginType, cronConfigs)
}

func (o *Onewayer) StartCronners() {
	for _, cronName := range o.CronJobs {
		o.logf.Infof("start cronjob plugin(%s)", cronName)
		go cronjob.Plugin[cronName].DoSchedule()
	}
}

func (o *Onewayer) StopCronners() {
	for _, cronName := range o.CronJobs {
		cronjob.Plugin[cronName].Stop()
		o.logf.Infof("stop cronjob plugin(%s)", cronName)
	}
}

func (o *Onewayer) GetStatus() bool {
	areAllPluginsCompleted := plugin.InputStatus.Completed && plugin.TransitStatus.Completed && plugin.ProcessStatus.Completed && plugin.OutputStatus.Completed
	return areAllPluginsCompleted
}
