package monitoring

import (
	"net/http"

	"github.com/jdjgya/service-frame/pkg/sync/plugin/plug"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricPuller struct {
	listener plug.Listener
}

func GetMetricPuller() MetricsReporter {
	http.Handle("/metrics", promhttp.Handler())

	return &MetricPuller{
		listener: &http.Server{
			Addr: "0.0.0.0:2112",
		},
	}
}

func (m *MetricPuller) DoReport() {
	go m.listener.ListenAndServe()
}
