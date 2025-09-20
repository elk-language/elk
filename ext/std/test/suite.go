package test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

type SuiteOrCase interface {
	traverse(enter func(test SuiteOrCase) TraverseOption, leave func(test SuiteOrCase) TraverseOption) TraverseOption
}

// Represents a test suite, a group of tests like `describe` or `context`
type Suite struct {
	Name       string
	Location   *position.Location
	Parent     *Suite
	SubSuites  []*Suite
	Cases      []*Case
	BeforeEach []*vm.Closure
	AfterEach  []*vm.Closure
	BeforeAll  []*vm.Closure
	AfterAll   []*vm.Closure
	FullMatch  bool
	caseCount  int
}

func (s *Suite) countCases() int {
	var counter int

	counter += len(s.Cases)
	for _, subSuite := range s.SubSuites {
		counter += subSuite.CaseCount()
	}

	return counter
}

func (s *Suite) CaseCount() int {
	if s.caseCount >= 0 {
		return s.caseCount
	}

	s.caseCount = s.countCases()
	return s.caseCount
}

func (s *Suite) traverse(enter func(test SuiteOrCase) TraverseOption, leave func(test SuiteOrCase) TraverseOption) TraverseOption {
	switch enter(s) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(s)
	}

	for _, testCase := range s.Cases {
		if testCase.traverse(enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, subSuite := range s.SubSuites {
		if subSuite.traverse(enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(s)
}

// Create a new tests suite
func NewSuite(name string, parent *Suite, loc *position.Location) *Suite {
	return &Suite{
		Name:      name,
		Parent:    parent,
		Location:  loc,
		caseCount: -1,
	}
}

func (s *Suite) NewSubSuite(name string, loc *position.Location) *Suite {
	subSuite := NewSuite(name, s, loc)
	subSuite.FullMatch = s.FullMatch
	return subSuite
}

func (s *Suite) RegisterSubSuite(subSuite *Suite) {
	if !subSuite.FullMatch {
		subSuite.FullMatch = s.FullMatch
	}
	s.SubSuites = append(s.SubSuites, subSuite)
}

func (s *Suite) NewCase(name string, fn *vm.Closure) *Case {
	return NewCase(name, fn, s)
}

func (s *Suite) RegisterCase(testCase *Case) {
	s.Cases = append(s.Cases, testCase)
}

func (s *Suite) FullName() string {
	if s.Parent == nil {
		return s.Name
	}
	parentFullName := s.Parent.FullName()
	if parentFullName == "" {
		return s.Name
	}

	return fmt.Sprintf("%s %s", parentFullName, s.Name)
}

func (s *Suite) FullNameWithSeparator() string {
	if s.Parent == nil {
		return s.Name
	}
	parentFullName := s.Parent.FullName()
	if parentFullName == "" {
		return s.Name
	}

	return fmt.Sprintf("%s > %s", parentFullName, s.Name)
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

func (s *Suite) Run(v *vm.VM, events chan<- *ReportEvent, rng *rand.Rand, ctx context.Context) *SuiteReport {
	if isDone(ctx) {
		return nil
	}

	var ok bool
	startTime := time.Now()

	suiteReport := NewSuiteReport(s)
	suiteReport.status = TEST_RUNNING
	events <- NewSuiteReportEvent(suiteReport, REPORT_START_SUITE)

	if s.CaseCount() == 0 {
		suiteReport.status = TEST_SKIPPED
		events <- NewSuiteReportEvent(suiteReport, REPORT_FINISH_SUITE)
		return suiteReport
	}

	suiteReport, ok = s.runBeforeAll(startTime, suiteReport, v, events, ctx)
	if !ok {
		return suiteReport
	}

	for _, testCase := range shuffleCases(s.Cases, rng) {
		caseReport := testCase.Run(v, events, ctx)
		if caseReport == nil {
			s.runAfterAll(startTime, suiteReport, v)
			return nil
		}
		suiteReport.RegisterCaseReport(caseReport)
	}

	for _, subSuite := range s.SubSuites {
		subSuiteReport := subSuite.Run(v, events, rng, ctx)
		if subSuiteReport == nil {
			return nil
		}
		suiteReport.RegisterSubSuiteReport(subSuiteReport)
	}

	s.runAfterAll(startTime, suiteReport, v)
	if isDone(ctx) {
		return nil
	}

	suiteReport.UpdateStatus(TEST_SUCCESS)
	suiteReport.duration = time.Since(startTime)
	events <- NewSuiteReportEvent(suiteReport, REPORT_FINISH_SUITE)
	return suiteReport
}

func (s *Suite) runBeforeAll(startTime time.Time, report *SuiteReport, v *vm.VM, events chan<- *ReportEvent, ctx context.Context) (*SuiteReport, bool) {
	for _, hook := range s.BeforeAll {
		if isDone(ctx) {
			return nil, false
		}
		_, err := v.CallClosure(hook)
		if !err.IsUndefined() {
			var status TestStatus
			if value.IsA(err, AssertionErrorClass) {
				status = TEST_FAILED
			} else {
				status = TEST_ERROR
			}
			report.status = status
			report.RegisterErr(
				Err{
					Typ:        ErrBeforeAll,
					Err:        err,
					StackTrace: v.GetStackTrace(),
				},
			)
			report.duration = time.Since(startTime)
			v.ResetError()
			events <- NewSuiteReportEvent(report, REPORT_FINISH_SUITE)
			return report, false
		}
	}

	return report, true
}

func (s *Suite) runAfterAll(startTime time.Time, report *SuiteReport, v *vm.VM) {
	for _, hook := range s.AfterAll {
		_, err := v.CallClosure(hook)
		if !err.IsUndefined() {
			var status TestStatus
			if value.IsA(err, AssertionErrorClass) {
				status = TEST_FAILED
			} else {
				status = TEST_ERROR
			}
			report.status = status
			report.RegisterErr(
				Err{
					Typ:        ErrAfterAll,
					Err:        err,
					StackTrace: v.GetStackTrace(),
				},
			)
			report.duration = time.Since(startTime)
			v.ResetError()
		}
	}
}

func noopTraverseSuite(test SuiteOrCase) TraverseOption { return TraverseContinue }

func TraverseSuite(test SuiteOrCase, enter func(test SuiteOrCase) TraverseOption, leave func(test SuiteOrCase) TraverseOption) {
	if enter == nil {
		enter = noopTraverseSuite
	}
	if leave == nil {
		leave = noopTraverseSuite
	}
	test.traverse(enter, leave)
}
