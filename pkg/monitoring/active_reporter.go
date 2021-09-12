package monitoring

import (
	"time"

	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type ActiveReporter struct {
	Pusher
}

type Pusher interface {
	Add() error
	Gatherer(prometheus.Gatherer) *push.Pusher
}

func GetActiveReporter() MetricsReporter {
	registry := prometheus.NewRegistry()
	registry.MustRegister(inputOK)

	return &ActiveReporter{
		Pusher: push.New("http://127.0.0.1:9091", plugin.Service).Gatherer(registry),
	}
}

func (a *ActiveReporter) DoReport() {
	for {
		inputOK.Set(float64(plugin.Metrics.InputOK))

		a.Pusher.Add()
		time.Sleep(60 * time.Second)
	}
}
