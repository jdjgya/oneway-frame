package monitoring

import (
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
)
