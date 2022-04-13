package worker

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
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
	"github.com/mohae/deepcopy"
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
	// Input    string
	// Transit  string
	// Process  string
	// Output   string

	Inputs    []string
	Transits  []string
	Processes []string
	Outputs   []string

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

func (o *Onewayer) getParterConfig(pluginType string) ([]string, []interface{}) {
	parterNames := []string{}
	rawConfigs := configer.Get(pluginType)

	for index, rawConfig := range rawConfigs.([]interface{}) {
		subConfig := (rawConfig.(map[interface{}]interface{}))

		pluginName := subConfig[name].(string)
		pluginGroup := subConfig["group"].(string)

		parterName := strings.Join([]string{pluginType, pluginName, strconv.Itoa(index)}, "-")

		switch pluginType {
		case plugin.Input:
			// o.Input = pluginName
			o.Inputs = append(o.Inputs, parterName)
			obj := deepcopy.Copy(input.Plugin[pluginName]).(input.Input)
			plug.Parters[parterName] = obj
			input.Plugin[parterName] = obj

			plugin.I2TChan[pluginGroup] = make(chan []byte, plugin.ChanSize)
		case plugin.Transit:
			// o.Transit = pluginName
			o.Transits = append(o.Transits, parterName)
			plug.Parters[parterName] = deepcopy.Copy(transit.Plugin[pluginName]).(transit.Transit)
			transit.Plugin[parterName] = plug.Parters[parterName].(transit.Transit)

			plugin.T2PChan[pluginGroup] = make(chan map[string]string, plugin.ChanSize)

		case plugin.Process:
			// o.Process = pluginName
			o.Processes = append(o.Processes, parterName)
			plug.Parters[parterName] = deepcopy.Copy(process.Plugin[pluginName]).(process.Process)
			process.Plugin[parterName] = plug.Parters[parterName].(process.Process)

			plugin.P2OChan[pluginGroup] = make(chan map[string]string, plugin.ChanSize)

		case plugin.Output:
			// o.Output = pluginName
			o.Outputs = append(o.Outputs, parterName)
			plug.Parters[parterName] = deepcopy.Copy(output.Plugin[pluginName]).(output.Output)
			output.Plugin[parterName] = plug.Parters[parterName].(output.Output)

		}

		parterNames = append(parterNames, parterName)
	}

	// os.Exit(1)

	// parterName := strings.Join([]string{pluginType, pluginName}, "-")

	// switch pluginType {
	// case plugin.Input:
	// 	o.Input = pluginName
	// 	plug.Parters[parterName] = input.Plugin[pluginName]
	// case plugin.Transit:
	// 	o.Transit = pluginName
	// 	plug.Parters[parterName] = transit.Plugin[pluginName]
	// case plugin.Process:
	// 	o.Process = pluginName
	// 	plug.Parters[parterName] = process.Plugin[pluginName]
	// case plugin.Output:
	// 	o.Output = pluginName
	// 	plug.Parters[parterName] = output.Plugin[pluginName]
	// }

	// return parterName, rawConfig

	return parterNames, rawConfigs.([]interface{})
}

func (o *Onewayer) setParterConfig(parterNames []string, parterConfigs []interface{}) {
	for i := 0; i < len(parterNames); i++ {
		pName := parterNames[i]
		fmt.Println("pName", pName)
		plug.Parters[pName].SetConfig(parterConfigs[i])
		err := plug.Parters[pName].CheckConfig()
		if err != nil {
			o.logf.Errorf("failed to set parter config. error details: %s", err.Error())
			osExit(1)
		}
	}
}

func (o *Onewayer) SetParter(pluginType string) {
	parterName, parterConfig := o.getParterConfig(pluginType)
	o.setParterConfig(parterName, parterConfig)
}

func (o *Onewayer) StartParters() {
	for _, inputname := range o.Inputs {
		o.logf.Infof("start input plugin(%s)", inputname)
		go input.Plugin[inputname].DoInput()
	}

	for _, transitname := range o.Transits {
		o.logf.Infof("start transit plugin(%s)", transitname)
		go transit.Plugin[transitname].DoTransit()
	}

	for _, processname := range o.Processes {
		o.logf.Infof("start process plugin(%s)", processname)
		go process.Plugin[processname].DoProcess()
	}

	for _, outputname := range o.Outputs {
		o.logf.Infof("start output plugin(%s)", outputname)
		go output.Plugin[outputname].DoOutput()
	}
}

func (o *Onewayer) StopParters() {
	o.logf.Infof("start input plugin(%s)", o.Inputs)
	for _, inputname := range o.Inputs {
		input.Plugin[inputname].Stop()
		o.logf.Infof("stop input plugin(%s)", inputname)
	}

	for _, transitname := range o.Transits {
		transit.Plugin[transitname].Stop()
		o.logf.Infof("stop transit plugin(%s)", transitname)
	}

	for _, processname := range o.Processes {
		process.Plugin[processname].Stop()
		o.logf.Infof("stop process plugin(%s)", processname)
	}

	for _, outputname := range o.Outputs {
		output.Plugin[outputname].Stop()
		o.logf.Infof("stop output plugin(%s)", outputname)
	}
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
