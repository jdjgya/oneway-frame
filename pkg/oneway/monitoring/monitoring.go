package monitoring

type MetricsReporter interface {
	DoReport()
}

type Monitor struct {
	MetricsReporter
}

func getReporter(isOneTimeExec bool) MetricsReporter {
	if isOneTimeExec {
		return GetMetricPusher()
	}

	return GetMetricPuller()
}

func (m *Monitor) SetRunMode(isOneTimeExec bool) {
	m.MetricsReporter = getReporter(isOneTimeExec)
}

func (m *Monitor) TraceMetric() {
	go collectMetric()
	go m.MetricsReporter.DoReport()
}
