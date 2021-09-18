package monitoring

import (
	"time"

	"github.com/jdjgya/service-frame/pkg/sync/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	interactOK = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "interact_ok",
			Help: "",
		},
	)

	interactErr = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "interact_err",
			Help: "",
		},
	)

	transitOK = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "transit_ok",
			Help: "",
		},
	)

	transitErr = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "transit_err",
			Help: "",
		},
	)

	processOK = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "process_ok",
			Help: "",
		},
	)

	processErr = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "process_err",
			Help: "",
		},
	)

	requestOK = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "request_ok",
			Help: "",
		},
	)

	requestErr = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "request_err",
			Help: "",
		},
	)
)

func collectMetric() {
	for {
		interactOK.Set(float64(plugin.Metrics.InteractOK))
		interactErr.Set(float64(plugin.Metrics.InteractErr))

		transitOK.Set(float64(plugin.Metrics.TransitOK))
		transitErr.Set(float64(plugin.Metrics.TransitErr))

		processOK.Set(float64(plugin.Metrics.ProcessOK))
		processErr.Set(float64(plugin.Metrics.ProcessErr))

		requestOK.Set(float64(plugin.Metrics.RequestOK))
		requestErr.Set(float64(plugin.Metrics.RequestErr))

		time.Sleep(60 * time.Second)
	}
}
