package worker

import (
	"strconv"

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

type (
	rootConfType = map[string]interface{}
	leafConfType = map[interface{}]interface{}
)

type Syncer struct {
	Interact string
	CronJobs []string

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

func (s *Syncer) setStagerConfig(stager string, stageConfig interface{}) {
	plug.Stagers[stager].SetConfig(stageConfig)
	err := plug.Stagers[stager].CheckConfig()
	if err != nil {
		s.logf.Error(err)
		osExit(1)
	}
}

func (s *Syncer) getStagerConfig(stageType string, stageIndex int, rawConf interface{}) (string, interface{}) {
	var stager string
	stageConfig := rawConf.(leafConfType)[stageType]
	stageName := stageConfig.(leafConfType)[name].(string)
	stager = strings.Join([]string{stageType, stageName, strconv.Itoa(stageIndex)}, "-")

	switch stageType {
	case plugin.StageTransit:
		plugin.ActivatedTransit[stageIndex] = stager
		transit.Plugins[stager] = transit.Plugins[stageName].New()
		plug.Stagers[stager] = transit.Plugins[stager]

	case plugin.StageProcess:
		plugin.ActivatedProcess[stageIndex] = stager
		process.Plugins[stager] = process.Plugins[stageName].New()
		plug.Stagers[stager] = process.Plugins[stager]

	case plugin.StageRequest:
		plugin.ActivatedRequest[stageIndex] = stager
		request.Plugins[stager] = request.Plugins[stageName].New()
		plug.Stagers[stager] = request.Plugins[stager]
	}

	return stager, stageConfig
}

func (s *Syncer) SetStager(stageType string, stageIndex int, stageConfig interface{}) {
	stager, stageConfig := s.getStagerConfig(stageType, stageIndex, stageConfig)
	s.setStagerConfig(stager, stageConfig)
}

func (s *Syncer) setParterConfig(parter string, pluginConf interface{}) {
	plug.Parters[parter].SetConfig(pluginConf)
	err := plug.Parters[parter].CheckConfig()
	if err != nil {
		s.logf.Error(err)
		osExit(1)
	}
}

func (s *Syncer) getPatternConfs(pluginConf interface{}) []interface{} {
	patterns := pluginConf.(rootConfType)["patterns"]
	return patterns.([]interface{})
}

func (s *Syncer) getParterConfig(pluginType string) (string, interface{}) {
	rawConf := configer.Get(pluginType)
	pluginName := rawConf.(rootConfType)[name].(string)
	s.Interact = pluginName

	parter := strings.Join([]string{pluginType, pluginName}, "-")
	plug.Parters[parter] = interact.Plugins[s.Interact]

	return parter, rawConf
}

func (s *Syncer) SetParter(pluginType string) {
	parter, pluginConf := s.getParterConfig(pluginType)
	patternConfs := s.getPatternConfs(pluginConf)

	plugin.ActivatedTransit = make([]string, len(patternConfs))
	plugin.ActivatedProcess = make([]string, len(patternConfs))
	plugin.ActivatedRequest = make([]string, len(patternConfs))

	for patternIndex, patternConf := range patternConfs {
		s.SetStager(plugin.StageTransit, patternIndex, patternConf)
		s.SetStager(plugin.StageProcess, patternIndex, patternConf)
		s.SetStager(plugin.StageRequest, patternIndex, patternConf)
	}

	s.setParterConfig(parter, pluginConf)
}

func (s *Syncer) StartParters() {
	s.logf.Infof("start interact plugin(%s)", s.Interact)
	go interact.Plugins[s.Interact].DoInteract()
}

func (s *Syncer) StopParters() {
	interact.Plugins[s.Interact].Stop()
	s.logf.Infof("stop input plugin(%s)", s.Interact)
}

func (s *Syncer) getCronnerConfig(pluginType string) map[string]leafConfType {
	cronConfigs := make(map[string]leafConfType)
	rawConfig := reflect.ValueOf(configer.Get(pluginType))

	for i := 0; i < rawConfig.Len(); i++ {
		pluginConfig := rawConfig.Index(i).Interface().(leafConfType)
		pluginName := pluginConfig[name].(string)
		cronName := strings.Join([]string{pluginType, pluginName}, "-")

		cronConfigs[pluginName] = make(leafConfType)
		cronConfigs[pluginName] = pluginConfig
		plug.Cronners[cronName] = cronjob.Plugin[pluginName]
	}

	return cronConfigs
}

func (s *Syncer) setCronnerConfig(pluginType string, cronConfigs map[string]leafConfType) {
	for cronName, cronConfig := range cronConfigs {
		cronner := strings.Join([]string{pluginType, cronName}, "-")
		plug.Cronners[cronner].SetConfig(cronConfig)
		err := plug.Cronners[cronner].CheckConfig()
		if err != nil {
			s.logf.Errorf("failed to set cronjob plugin(%s). error: %s", cronName, err.Error())
			osExit(1)
		}

		s.CronJobs = append(s.CronJobs, cronName)
	}
}

func (s *Syncer) SetCronner(pluginType string) {
	cronConfigs := s.getCronnerConfig(pluginType)
	s.setCronnerConfig(pluginType, cronConfigs)
}

func (s *Syncer) StartCronners() {
	for _, cronName := range s.CronJobs {
		s.logf.Infof("start cronjob plugin(%s)", cronName)
		go cronjob.Plugin[cronName].DoSchedule()
	}
}

func (s *Syncer) StopCronners() {
	for _, cronName := range s.CronJobs {
		cronjob.Plugin[cronName].Stop()
		s.logf.Infof("stop cronjob plugin(%s)", cronName)
	}
}
