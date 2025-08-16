package test

import (
	"time"

	"github.com/elk-language/elk/value"
)

// Contains the result of running a test case
type CaseReport struct {
	Case       *Case
	Error      value.Value // Assertion failure or runtime error
	StackTrace *value.StackTrace
	Duration   time.Duration // The amount of time it took to run the test
	Status     TestStatus
}

func NewCaseReport(cas *Case) *CaseReport {
	return &CaseReport{
		Case: cas,
	}
}
