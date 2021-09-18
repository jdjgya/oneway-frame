package monitoring

type MetricsReporter interface {
	DoReport()
}

type Monitor struct {
	MetricsReporter
}

func getReporter() MetricsReporter {
	return GetMetricPuller()
}

func (m *Monitor) SetMetricReporter() {
	m.MetricsReporter = GetMetricPuller()
}

func (m *Monitor) TraceMetric() {
	go collectMetric()
	go m.MetricsReporter.DoReport()
}
