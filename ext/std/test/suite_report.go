package test

import (
	"time"

	"github.com/elk-language/elk/value"
)

// Contains the result of running a test suite
type SuiteReport struct {
	Suite           *Suite
	CaseReports     []*CaseReport
	SubSuiteReports []*SuiteReport
	err             value.Value
	stackTrace      *value.StackTrace
	duration        time.Duration
	status          TestStatus
}

func NewSuiteReport(suite *Suite) *SuiteReport {
	return &SuiteReport{
		Suite: suite,
	}
}

func (s *SuiteReport) traverse(enter func(report Report) TraverseOption, leave func(report Report) TraverseOption) TraverseOption {
	switch enter(s) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(s)
	}

	for _, caseReport := range s.CaseReports {
		if caseReport.traverse(enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, subSuiteReport := range s.SubSuiteReports {
		if subSuiteReport.traverse(enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(s)
}

func (s *SuiteReport) Error() value.Value {
	return s.err
}

func (s *SuiteReport) StackTrace() *value.StackTrace {
	return s.stackTrace
}

func (s *SuiteReport) Duration() time.Duration {
	return s.duration
}

func (s *SuiteReport) Status() TestStatus {
	return s.status
}

func (s *SuiteReport) RegisterSubSuiteReport(subSuiteReport *SuiteReport) {
	s.SubSuiteReports = append(s.SubSuiteReports, subSuiteReport)
	s.UpdateStatus(subSuiteReport.status)
}

func (s *SuiteReport) UpdateStatus(newStatus TestStatus) {
	switch newStatus {
	case TEST_ERROR:
		s.status = TEST_ERROR
	case TEST_FAILED:
		if s.status != TEST_ERROR {
			s.status = TEST_FAILED
		}
	case TEST_SUCCESS:
		if s.status == TEST_RUNNING {
			s.status = TEST_SUCCESS
		}
	}
}

func (s *SuiteReport) RegisterCaseReport(caseReport *CaseReport) {
	s.CaseReports = append(s.CaseReports, caseReport)
	s.UpdateStatus(caseReport.status)
}
