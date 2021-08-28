package pattern

import (
	"github.com/gin-gonic/gin"
)

var (
	Plugins = make(map[string]Pattern)
)

type Pattern interface {
	SetConfig()
	RegisterRouter(*gin.Engine, string, string) error
	SetRouterStage(string, string, string)
}
