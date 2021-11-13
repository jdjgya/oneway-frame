package worker

import (
	"github.com/jdjgya/service-frame/pkg/config"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/sync/plugin"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/cronjob"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/process"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/request"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/transit"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/plug"

	"os"
	"reflect"
	"strings"

	"go.uber.org/zap"
)

const (
	name   = "name"
	module = "interactor"
)

var (
	configer = config.GetConfiger()
	osExit   = os.Exit
)

type Syncer struct {
	Interact     string
	StageTransit string
	StageProcess string
	StageRequest string
	CronJobs     []string

	log  *zap.Logger
	logf *zap.SugaredLogger
}

func InitWorker() Worker {
	logger := log.GetLogger(module)

	return &Syncer{
		log:  logger,
		logf: logger.Sugar(),
	}
}

func (r *Syncer) setStagerConfig(stager string, stageConfig interface{}) {
	plug.Stagers[stager].SetConfig(stageConfig)
	err := plug.Stagers[stager].CheckConfig()
	if err != nil {
		r.logf.Error(err)
		osExit(1)
	}
}

func (r *Syncer) getStagerConfig(stageType string, rawConf interface{}) (string, interface{}) {
	var stager string
	stageConfig := (rawConf.(map[string]interface{}))[stageType]

	switch stageType {
	case plugin.StageTransit:
		stageName := (stageConfig.(map[string]interface{}))[name].(string)
		stager = strings.Join([]string{stageType, stageName}, "-")

		r.StageTransit = stageName
		plugin.ActivatedTransit = stageName
		plug.Stagers[stager] = transit.Plugins[r.StageTransit]

	case plugin.StageProcess:
		stageName := (stageConfig.(map[string]interface{}))[name].(string)
		stager = strings.Join([]string{stageType, stageName}, "-")

		r.StageProcess = stageName
		plugin.ActivatedProcess = stageName
		plug.Stagers[stager] = process.Plugins[r.StageProcess]

	case plugin.StageRequest:
		stageName := (stageConfig.(map[string]interface{}))[name].(string)
		stager = strings.Join([]string{stageType, stageName}, "-")

		r.StageRequest = stageName
		plugin.ActivatedRequest = stageName
		plug.Stagers[stager] = request.Plugins[r.StageRequest]
	}

	return stager, stageConfig
}

func (r *Syncer) SetStager(stageType string, stageConfig interface{}) {
	stager, stageConfig := r.getStagerConfig(stageType, stageConfig)
	r.setStagerConfig(stager, stageConfig)
}

func (r *Syncer) setParterConfig(parter string, pluginConf interface{}) {
	plug.Parters[parter].SetConfig(pluginConf)
	err := plug.Parters[parter].CheckConfig()
	if err != nil {
		r.logf.Error(err)
		osExit(1)
	}
}

func (r *Syncer) getParterConfig(pluginType string) (string, interface{}) {
	rawConf := configer.Get(pluginType)
	pluginName := (rawConf.(map[string]interface{}))[name].(string)
	r.Interact = pluginName

	parter := strings.Join([]string{pluginType, pluginName}, "-")
	plug.Parters[parter] = interact.Plugins[r.Interact]

	return parter, rawConf
}

func (r *Syncer) SetParter(pluginType string) {
	parter, pluginConf := r.getParterConfig(pluginType)
	r.SetStager(plugin.StageTransit, pluginConf)
	r.SetStager(plugin.StageProcess, pluginConf)
	r.SetStager(plugin.StageRequest, pluginConf)
	r.setParterConfig(parter, pluginConf)
}

func (r *Syncer) StartParters() {
	r.logf.Infof("start interact plugin(%s)", r.Interact)
	go interact.Plugins[r.Interact].DoInteract()
}

func (r *Syncer) StopParters() {
	interact.Plugins[r.Interact].Stop()
	r.logf.Infof("stop input plugin(%s)", r.Interact)
}

func (r *Syncer) getCronnerConfig(pluginType string) map[string]map[interface{}]interface{} {
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

func (r *Syncer) setCronnerConfig(pluginType string, cronConfigs map[string]map[interface{}]interface{}) {
	for cronName, cronConfig := range cronConfigs {
		cronner := strings.Join([]string{pluginType, cronName}, "-")
		plug.Cronners[cronner].SetConfig(cronConfig)
		err := plug.Cronners[cronner].CheckConfig()
		if err != nil {
			r.logf.Errorf("failed to set cronjob plugin(%s). error: %s", cronName, err.Error())
			osExit(1)
		}

		r.CronJobs = append(r.CronJobs, cronName)
	}
}

func (r *Syncer) SetCronner(pluginType string) {
	cronConfigs := r.getCronnerConfig(pluginType)
	r.setCronnerConfig(pluginType, cronConfigs)
}

func (r *Syncer) StartCronners() {
	for _, cronName := range r.CronJobs {
		r.logf.Infof("start cronjob plugin(%s)", cronName)
		go cronjob.Plugin[cronName].DoSchedule()
	}
}

func (r *Syncer) StopCronners() {
	for _, cronName := range r.CronJobs {
		cronjob.Plugin[cronName].Stop()
		r.logf.Infof("stop cronjob plugin(%s)", cronName)
	}
}
