package monitoring

type MetricsReporter interface {
	DoReport()
}

type Monitor struct {
	MetricsReporter
}

func getReporter(isOneTimeExec bool) MetricsReporter {
	if isOneTimeExec {
		return GetActiveReporter()
	}

	return GetPassiveReporter()
}

func (m *Monitor) SetRunMode(isOneTimeExec bool) {
	m.MetricsReporter = getReporter(isOneTimeExec)
}

func (m *Monitor) TraceMetric() {
	go m.MetricsReporter.DoReport()
}
