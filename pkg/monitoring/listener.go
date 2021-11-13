package monitoring

import (
	"context"
	"net/http"

	"github.com/jdjgya/service-frame/pkg/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	socket         = "0.0.0.0:2112"
	path           = "/metrics"
	listenerModule = "metricPuller"
)

var (
	listenerLogger = log.GetLogger(listenerModule).Sugar()
)

type listener interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type MetricPuller struct {
	listener
}

func GetMetricPuller() Metricer {
	http.Handle(path, promhttp.Handler())

	return &MetricPuller{
		listener: &http.Server{
			Addr: socket,
		},
	}
}

func (m *MetricPuller) Report() {
	go func() {
		err := m.listener.ListenAndServe()
		if err != nil {
			listenerLogger.Errorf("failed to listen on metric tunnel. error details: %s", err.Error())
		}
	}()
}
