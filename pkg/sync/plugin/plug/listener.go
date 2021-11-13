package plug

import (
	"context"
)

type Listener interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
