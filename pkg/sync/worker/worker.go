package worker

import "github.com/jdjgya/service-frame/pkg/sync/plugin/plug"

type Worker interface {
	plug.PartUser
	plug.CronUser
}
