package test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

type SuiteOrCase interface {
	traverse(enter func(test SuiteOrCase) TraverseOption, leave func(test SuiteOrCase) TraverseOption) TraverseOption
}

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

func (s *Suite) FullNameWithSeparator() string {
	if s.Parent == nil {
		return s.Name
	}

	return fmt.Sprintf("%s › %s", s.Parent.FullName(), s.Name)
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

	var err value.Value
	startTime := time.Now()

	suiteReport := NewSuiteReport(s)
	suiteReport.status = TEST_RUNNING
	events <- NewSuiteReportEvent(suiteReport, REPORT_START_SUITE)

	for _, hook := range s.BeforeAll {
		if isDone(ctx) {
			return nil
		}
		_, err = v.CallClosure(hook)
		if !err.IsUndefined() {
			suiteReport.status = TEST_ERROR
			suiteReport.err = err
			suiteReport.stackTrace = v.GetStackTrace()
			suiteReport.duration = time.Since(startTime)
			v.ResetError()
			events <- NewSuiteReportEvent(suiteReport, REPORT_FINISH_SUITE)
			return suiteReport
		}
	}

	for _, testCase := range shuffleCases(s.Cases, rng) {
		caseReport := testCase.Run(v, events, ctx)
		if caseReport == nil {
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

	for _, hook := range s.AfterAll {
		if isDone(ctx) {
			return nil
		}
		_, err = v.CallClosure(hook)
		if !err.IsUndefined() {
			suiteReport.status = TEST_ERROR
			suiteReport.err = err
			suiteReport.stackTrace = v.GetStackTrace()
			suiteReport.duration = time.Since(startTime)
			v.ResetError()
			events <- NewSuiteReportEvent(suiteReport, REPORT_FINISH_SUITE)
			return suiteReport
		}
	}

	if isDone(ctx) {
		return nil
	}

	suiteReport.UpdateStatus(TEST_SUCCESS)
	suiteReport.duration = time.Since(startTime)
	events <- NewSuiteReportEvent(suiteReport, REPORT_FINISH_SUITE)
	return suiteReport
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
