package test

import (
	"time"

	"github.com/elk-language/elk/value"
)

// Contains the result of running a test suite
type SuiteReport struct {
	Suite           *Suite
	SubSuiteReports []*SuiteReport
	CaseReports     []*CaseReport
	Error           value.Value
	StackTrace      *value.StackTrace
	Duration        time.Duration
	Status          TestStatus
}

func NewSuiteReport(suite *Suite) *SuiteReport {
	return &SuiteReport{
		Suite: suite,
	}
}

func (s *SuiteReport) RegisterSubSuiteReport(subSuiteReport *SuiteReport) {
	s.SubSuiteReports = append(s.SubSuiteReports, subSuiteReport)
}

func (s *SuiteReport) RegisterCaseReport(caseReport *CaseReport) {
	s.CaseReports = append(s.CaseReports, caseReport)
}
