package request

import (
	"github.com/jdjgya/service-frame/pkg/sync/plugin/plug"
)

var (
	Plugins = make(map[string]Request)
)

type Request interface {
	New() Request
	plug.Stager
}
