package test

import "time"

// Contains the result of running a test suite
type SuiteReport struct {
	Suite           *Suite
	SubSuiteReports []*SuiteReport
	CaseReports     []*CaseReport
	Duration        time.Duration
	Status          TestStatus
}
