package monitoring

import (
	"time"

	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	module = "metricsPusher"
)

var (
	pushLogger = log.GetLogger(module)
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
		err := m.Pusher.Add()
		if err != nil {
			pushLogger.Error("failed to add metric")
		}

		time.Sleep(60 * time.Second)
	}
}
