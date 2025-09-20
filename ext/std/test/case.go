package test

import (
	"context"
	"fmt"
	"iter"
	"time"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Represents a single test case
type Case struct {
	Name   string
	Fn     *vm.Closure
	Parent *Suite
}

func (c *Case) FullMatch() bool {
	return c.Parent.FullMatch
}

func (c *Case) traverse(enter func(test SuiteOrCase) TraverseOption, leave func(test SuiteOrCase) TraverseOption) TraverseOption {
	switch enter(c) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(c)
	}

	return leave(c)
}

func NewCase(name string, fn *vm.Closure, parent *Suite) *Case {
	return &Case{
		Name:   name,
		Fn:     fn,
		Parent: parent,
	}
}

func (c *Case) Location() *position.Location {
	return c.Fn.Bytecode.Location
}

func (c *Case) FullName() string {
	if c.Parent == nil {
		return c.Name
	}

	return fmt.Sprintf("%s %s", c.Parent.FullName(), c.Name)
}

func (c *Case) FullNameWithSeparator() string {
	if c.Parent == nil {
		return c.Name
	}

	return fmt.Sprintf("%s > %s", c.Parent.FullNameWithSeparator(), c.Name)
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func (c *Case) Parents() iter.Seq[*Suite] {
	return func(yield func(*Suite) bool) {
		currentParent := c.Parent

		for currentParent != nil {
			if !yield(currentParent) {
				return
			}

			currentParent = currentParent.Parent
		}
	}
}

func callCaseClosure(
	v *vm.VM,
	caseReport *CaseReport,
	startTime time.Time,
	closure *vm.Closure,
	typ ErrTyp,
) bool {
	var err value.Value
	_, err = v.CallClosure(closure)
	if !err.IsUndefined() {
		var status TestStatus
		if value.IsA(err, AssertionErrorClass) {
			status = TEST_FAILED
		} else {
			status = TEST_ERROR
		}

		caseReport.status = status
		caseReport.RegisterErr(
			Err{
				Typ:        typ,
				Err:        err,
				StackTrace: v.GetStackTrace(),
			},
		)
		caseReport.duration = time.Since(startTime)
		v.ResetError()
		return false
	}

	return true
}

func (c *Case) Run(v *vm.VM, events chan<- *ReportEvent, ctx context.Context) *CaseReport {
	if isDone(ctx) {
		return nil
	}

	var ok bool
	startTime := time.Now()

	caseReport := NewCaseReport(c)
	caseReport.status = TEST_RUNNING
	events <- NewCaseReportEvent(caseReport, REPORT_START_CASE)

	prevStdout := v.Stdout
	prevStderr := v.Stderr

	v.Stdout = &caseReport.stdout
	v.Stderr = &caseReport.stderr

	defer func() {
		v.Stdout = prevStdout
		v.Stderr = prevStderr
	}()

	caseReport, ok = c.runBeforeEach(startTime, caseReport, v, events, ctx)
	if !ok {
		c.runAfterEach(startTime, caseReport, v)
		return caseReport
	}

	if isDone(ctx) {
		return nil
	}
	callCaseClosure(v, caseReport, startTime, c.Fn, ErrCase)

	c.runAfterEach(startTime, caseReport, v)

	if isDone(ctx) {
		return nil
	}
	caseReport.UpdateStatus(TEST_SUCCESS)
	caseReport.duration = time.Since(startTime)
	events <- NewCaseReportEvent(caseReport, REPORT_FINISH_CASE)
	return caseReport
}

func (c *Case) runBeforeEach(startTime time.Time, report *CaseReport, v *vm.VM, events chan<- *ReportEvent, ctx context.Context) (*CaseReport, bool) {
	for parent := range c.Parents() {
		for _, hook := range parent.BeforeEach {
			if isDone(ctx) {
				return nil, false
			}
			if !callCaseClosure(v, report, startTime, hook, ErrBeforeEach) {
				return report, false
			}
		}
	}

	return report, true
}

func (c *Case) runAfterEach(startTime time.Time, report *CaseReport, v *vm.VM) {
	for parent := range c.Parents() {
		for _, hook := range parent.AfterEach {
			callCaseClosure(v, report, startTime, hook, ErrAfterEach)
		}
	}
}
