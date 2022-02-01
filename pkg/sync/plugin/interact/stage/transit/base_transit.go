package transit

import "github.com/jdjgya/service-frame/pkg/sync/plugin/plug"

var (
	Plugins = make(map[string]Transit)
)

type Transit interface {
	New() Transit
	plug.Stager
}
