package dummy

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/jdjgya/service-frame/example/sync/plugin/interact/http/pattern"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/sync/plugin"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/process"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/request"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/transit"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/plug"
	"go.uber.org/zap"
)

const (
	module   = "dummy"
	response = "dummy response"
)

var (
	post = "POST"
)

type DummyHandler struct {
	stages []plug.Stager

	log  *zap.Logger
	logf *zap.SugaredLogger
}

func init() {
	pattern.Plugins[module] = &DummyHandler{}
}

func (m *DummyHandler) SetConfig() {
	m.log = log.GetLogger(module)
	m.logf = m.log.Sugar()
}

func (m *DummyHandler) RegisterRouter(router *gin.Engine, method string, path string) error {
	switch method {
	case post:
		router.POST(path, m.post)
	default:
		m.logf.Errorf("unsupported method(%s) detected in pattern(%s)", method, path)
		return errors.New("register unsupported REST methog")
	}

	return nil
}

func (m *DummyHandler) SetRouterStage(trans, proc, req string) {
	m.stages = append(m.stages, transit.Plugins[trans], process.Plugins[proc], request.Plugins[req])
}

func (m *DummyHandler) unmarshalJobFromRequest(requestBody io.ReadCloser, msg *map[string]string) error {
	bodyByte, err := ioutil.ReadAll(requestBody)
	if err != nil {
		atomic.AddUint64(&plugin.Metrics.InteractErr, 1)
		return err
	}

	err = json.Unmarshal(bodyByte, &msg)
	if err != nil {
		atomic.AddUint64(&plugin.Metrics.InteractErr, 1)
		return err
	}

	atomic.AddUint64(&plugin.Metrics.InteractOK, 1)
	return nil
}

func (m *DummyHandler) handleJobByStages(msg *map[string]string) error {
	for _, stage := range m.stages {
		err := stage.Execute(msg)
		if err != nil {
			stage.AddError()
			return err
		}

		stage.AddSuccess()
	}

	return nil
}

func (m *DummyHandler) post(g *gin.Context) {
	g.JSON(http.StatusOK, response)
}
