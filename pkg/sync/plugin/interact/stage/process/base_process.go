package process

import "github.com/jdjgya/service-frame/pkg/sync/plugin/plug"

var (
	Plugins = make(map[string]Process)
)

type Process interface {
	New() Process
	plug.Stager
}
