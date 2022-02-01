package dummy1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jdjgya/service-frame/example/sync/plugin/interact/http/pattern"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/process"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/request"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact/stage/transit"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/plug"
	"go.uber.org/zap"
)

const (
	module   = "dummy1"
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
	fmt.Println(trans, proc, req)
	m.stages = append(m.stages, transit.Plugins[trans], process.Plugins[proc], request.Plugins[req])
}

func (m *DummyHandler) post(g *gin.Context) {
	dummyMsg := map[string]string{}
	for _, stage := range m.stages {
		err := stage.Execute(&dummyMsg)
		if err != nil {
			stage.AddError()
		}

		stage.AddSuccess()
	}

	g.JSON(http.StatusOK, response)
}
