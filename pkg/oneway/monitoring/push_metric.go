package monitoring

import (
	"time"

	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type MetricPusher struct {
	Pusher
}

type Pusher interface {
	Add() error
	Gatherer(prometheus.Gatherer) *push.Pusher
}

func GetMetricPusher() MetricsReporter {
	registry := prometheus.NewRegistry()
	registry.MustRegister(inputOK)

	return &MetricPusher{
		Pusher: push.New("http://127.0.0.1:9091", plugin.Service).Gatherer(registry),
	}
}

func (m *MetricPusher) DoReport() {
	for {
		m.Pusher.Add()
		time.Sleep(60 * time.Second)
	}
}
