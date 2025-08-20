package test

import "context"

// Contains the result of running a test suite
type Reporter interface {
	Report(events chan *ReportEvent, shutdown context.CancelFunc)
}
