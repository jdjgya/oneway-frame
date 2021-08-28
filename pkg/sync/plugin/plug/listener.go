package plug

import (
	"context"

	"github.com/jdjgya/service-frame/pkg/log"
)

var listenerLog = log.GetLogger("listener-helper")

type Listener interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
