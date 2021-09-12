package monitoring

import (
	"net/http"

	"github.com/jdjgya/service-frame/pkg/sync/plugin/plug"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PassiveReporter struct {
	listener plug.Listener
}

func GetPassiveReporter() MetricsReporter {
	http.Handle("/metrics", promhttp.Handler())

	return &PassiveReporter{
		listener: &http.Server{
			Addr: "0.0.0.0:2112",
		},
	}
}

func (p *PassiveReporter) DoReport() {
	go p.listener.ListenAndServe()
}
