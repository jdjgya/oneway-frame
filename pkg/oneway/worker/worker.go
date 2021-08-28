package worker

import "github.com/jdjgya/service-frame/pkg/oneway/plugin/plug"

type Worker interface {
	plug.PartUser
	plug.CronUser
	plug.Statuser
}
