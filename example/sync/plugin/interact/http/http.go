package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goinggo/mapstructure"
	"github.com/jdjgya/service-frame/example/sync/plugin/interact/http/pattern"
	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/sync/plugin"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/interact"
	"github.com/jdjgya/service-frame/pkg/sync/plugin/plug"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	module = "http"
)

var (
	Router *gin.Engine
)

type Http struct {
	ctx    context.Context
	cancel context.CancelFunc

	listener plug.Listener
	router   *gin.Engine
	config

	log  *zap.Logger
	logf *zap.SugaredLogger
}

type config struct {
	Name    string `validate:"required"`
	Address string `validate:"required"`
	Port    int    `validate:"required"`
	Pattern
}

type Pattern struct {
	Name    string `validate:"required"`
	Method  string `validate:"required"`
	Path    string `validate:"required"`
	Timeout int
}

func init() {
	http := &Http{}
	interact.Plugins[module] = http
}

func (h *Http) SetConfig(conf interface{}) {
	_ = mapstructure.Decode(conf, &h.config)
	h.ctx, h.cancel = context.WithCancel(context.Background())

	h.log = log.GetLogger(module)
	h.logf = h.log.Sugar()

	gin.DefaultWriter = ioutil.Discard
	h.router = gin.New()
	h.router.Use(gin.Recovery())

	patternName := h.config.Pattern.Name
	patternMethod := h.config.Pattern.Method
	err := pattern.Plugins[patternName].RegisterRouter(h.router, patternMethod, h.Path)
	if err != nil {
		h.logf.Errorf("faild to register router. error details: %s", err.Error())
	}
	pattern.Plugins[h.config.Pattern.Name].SetRouterStage(plugin.ActivatedTransit, plugin.ActivatedProcess, plugin.ActivatedRequest)

	socket := strings.Join([]string{h.Address, strconv.Itoa(h.Port)}, ":")
	h.listener = &http.Server{
		Addr:    socket,
		Handler: h.router,
	}

}

func (h *Http) CheckConfig() error {
	validate := validator.New()
	return validate.Struct(h.config)
}

func (h *Http) DoInteract() {
	err := h.listener.ListenAndServe()
	if err != nil {
		h.logf.Errorf("failed to listen on particular protocol. error: %s", err.Error())
		os.Exit(1)
	}
}

func (h *Http) Stop() {
	if err := h.listener.Shutdown(h.ctx); err != nil {
		h.logf.Errorf("failed to stop interact pluing(%s). error: %s", module, err.Error())
	}
}
