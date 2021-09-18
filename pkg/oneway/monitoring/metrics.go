package monitoring

import (
	"time"

	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	inputOK = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "input_ok",
			Help: "",
		},
	)

	inputErr = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "input_err",
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

	outputOK = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "output_ok",
			Help: "",
		},
	)

	outputErr = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "output_err",
			Help: "",
		},
	)
)

func collectMetric() {
	for {
		inputOK.Set(float64(plugin.Metrics.InputOK))
		inputErr.Set(float64(plugin.Metrics.InputErr))

		transitOK.Set(float64(plugin.Metrics.TransitOK))
		transitErr.Set(float64(plugin.Metrics.TransitErr))

		processOK.Set(float64(plugin.Metrics.ProcessOK))
		processErr.Set(float64(plugin.Metrics.ProcessErr))

		outputOK.Set(float64(plugin.Metrics.OutputOK))
		outputErr.Set(float64(plugin.Metrics.OutputErr))

		time.Sleep(60 * time.Second)
	}
}
