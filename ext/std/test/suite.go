package test

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Represents a test suite, a group of tests like `describe` or `context`
type Suite struct {
	Name       string
	Parent     *Suite
	SubSuites  []*Suite
	Cases      []*Case
	BeforeEach []*vm.Closure
	AfterEach  []*vm.Closure
	BeforeAll  []*vm.Closure
	AfterAll   []*vm.Closure
}

// Create a new tests suite
func NewSuite(name string, parent *Suite) *Suite {
	return &Suite{
		Name:   name,
		Parent: parent,
	}
}

func (s *Suite) NewSubSuite(name string) *Suite {
	newSuite := NewSuite(name, s)
	s.SubSuites = append(s.SubSuites, newSuite)
	return newSuite
}

func (s *Suite) NewCase(name string, fn *vm.Closure) *Case {
	newCase := NewCase(name, fn, s)
	s.Cases = append(s.Cases, newCase)
	return newCase
}

func (s *Suite) FullName() string {
	if s.Parent == nil {
		return s.Name
	}

	return fmt.Sprintf("%s %s", s.Parent.FullName(), s.Name)
}

func (s *Suite) RegisterBeforeEach(fn *vm.Closure) {
	s.BeforeEach = append(s.BeforeEach, fn)
}

func (s *Suite) RegisterBeforeAll(fn *vm.Closure) {
	s.BeforeAll = append(s.BeforeAll, fn)
}

func (s *Suite) RegisterAfterEach(fn *vm.Closure) {
	s.AfterEach = append(s.AfterEach, fn)
}

func (s *Suite) RegisterAfterAll(fn *vm.Closure) {
	s.AfterAll = append(s.AfterAll, fn)
}

func shuffleCases(cases []*Case, rng *rand.Rand) []*Case {
	shuffled := slices.Clone(cases)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

func (s *Suite) Run(v *vm.VM, events chan *ReportEvent, rng *rand.Rand) *SuiteReport {
	var err value.Value
	startTime := time.Now()

	suiteReport := NewSuiteReport(s)
	suiteReport.Status = TEST_RUNNING
	events <- NewSuiteReportEvent(suiteReport, REPORT_START_SUITE)

	for _, hook := range s.BeforeAll {
		_, err = v.CallClosure(hook)
		if !err.IsUndefined() {
			suiteReport.Status = TEST_ERROR
			suiteReport.Error = err
			suiteReport.StackTrace = v.ErrStackTrace()
			suiteReport.Duration = time.Since(startTime)
			events <- NewSuiteReportEvent(suiteReport, REPORT_FINISH_SUITE)
			return suiteReport
		}
	}

	for _, testCase := range shuffleCases(s.Cases, rng) {
		caseReport := testCase.Run(v, events)
		suiteReport.RegisterCaseReport(caseReport)
	}

	for _, subSuite := range s.SubSuites {
		subSuiteReport := subSuite.Run(v, events, rng)
		suiteReport.RegisterSubSuiteReport(subSuiteReport)
	}

	for _, hook := range s.AfterAll {
		_, err = v.CallClosure(hook)
		if !err.IsUndefined() {
			suiteReport.Status = TEST_ERROR
			suiteReport.Error = err
			suiteReport.StackTrace = v.ErrStackTrace()
			suiteReport.Duration = time.Since(startTime)
			events <- NewSuiteReportEvent(suiteReport, REPORT_FINISH_SUITE)
			return suiteReport
		}
	}

	suiteReport.Status = TEST_SUCCESS
	suiteReport.Duration = time.Since(startTime)
	events <- NewSuiteReportEvent(suiteReport, REPORT_FINISH_SUITE)
	return suiteReport
}
