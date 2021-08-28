package interact

import (
	"github.com/jdjgya/service-frame/pkg/sync/plugin/plug"
)

var (
	Plugins = make(map[string]Interactor)
)

type Interactor interface {
	plug.Parter
	DoInteract()
}
