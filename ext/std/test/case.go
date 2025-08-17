package test

import (
	"fmt"
	"time"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Represents a single test case
type Case struct {
	Name   string
	Fn     *vm.Closure
	Parent *Suite
}

func NewCase(name string, fn *vm.Closure, parent *Suite) *Case {
	return &Case{
		Name:   name,
		Fn:     fn,
		Parent: parent,
	}
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

	return fmt.Sprintf("%s › %s", c.Parent.FullNameWithSeparator(), c.Name)
}

func callCaseClosure(v *vm.VM, caseReport *CaseReport, events chan *ReportEvent, startTime time.Time, closure *vm.Closure) bool {
	var err value.Value
	_, err = v.CallClosure(closure)
	if !err.IsUndefined() {
		var status TestStatus
		if err.Class() == AssertionErrorClass {
			status = TEST_FAILED
		} else {
			status = TEST_ERROR
		}

		caseReport.status = status
		caseReport.err = err
		caseReport.stackTrace = v.GetStackTrace()
		caseReport.duration = time.Since(startTime)
		events <- NewCaseReportEvent(caseReport, REPORT_FINISH_CASE)
		return false
	}

	return true
}

func (c *Case) Run(v *vm.VM, events chan *ReportEvent) *CaseReport {
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

	if c.Parent != nil {
		for _, hook := range c.Parent.BeforeEach {
			if !callCaseClosure(v, caseReport, events, startTime, hook) {
				return caseReport
			}
		}
	}

	if !callCaseClosure(v, caseReport, events, startTime, c.Fn) {
		return caseReport
	}

	if c.Parent != nil {
		for _, hook := range c.Parent.AfterEach {
			if !callCaseClosure(v, caseReport, events, startTime, hook) {
				return caseReport
			}
		}
	}

	caseReport.status = TEST_SUCCESS
	caseReport.duration = time.Since(startTime)
	events <- NewCaseReportEvent(caseReport, REPORT_FINISH_CASE)
	return caseReport
}
