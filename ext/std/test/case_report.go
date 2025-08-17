package test

import (
	"bytes"
	"time"

	"github.com/elk-language/elk/value"
)

// Contains the result of running a test case
type CaseReport struct {
	Case       *Case
	err        value.Value // Assertion failure or runtime error
	stackTrace *value.StackTrace
	duration   time.Duration // The amount of time it took to run the test
	status     TestStatus
	stdout     bytes.Buffer
	stderr     bytes.Buffer
}

func (c *CaseReport) traverse(enter func(report Report) TraverseOption, leave func(report Report) TraverseOption) TraverseOption {
	switch enter(c) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(c)
	}

	return leave(c)
}

func (c *CaseReport) Error() value.Value {
	return c.err
}

func (c *CaseReport) StackTrace() *value.StackTrace {
	return c.stackTrace
}

func (c *CaseReport) Duration() time.Duration {
	return c.duration
}

func (c *CaseReport) Status() TestStatus {
	return c.status
}

func NewCaseReport(cas *Case) *CaseReport {
	return &CaseReport{
		Case: cas,
	}
}
