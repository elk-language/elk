package test

import (
	"time"

	"github.com/elk-language/elk/value"
)

// Contains the result of running a test case
type CaseReport struct {
	Case     *Case
	Failure  value.Value   // Failed assertion, test executed successfully bute expectations were not met
	Error    value.Value   // Runtime error, when executing the test an unexpected runtime error occurred
	Duration time.Duration // The amount of time it took to run the test
	Status   TestStatus
}
