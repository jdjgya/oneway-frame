package worker

import "onewayframe/pkg/plugin/plug"

type Worker interface {
	plug.PartUser
	plug.CronUser
	plug.Statuser
}
