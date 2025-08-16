package test

// Contains the result of running a test suite
type Reporter interface {
	Report(events chan *ReportEvent)
}

type SimpleReporter struct{}

func (SimpleReporter) Report(events chan *ReportEvent) {
	for event := range events {
		switch event.Type {
		case REPORT_FINISH_SUITE:
		case REPORT_FINISH_CASE:
			// TODO print `.`
		}
	}
}
