package test

import (
	"bytes"
	"time"
)

// Contains the result of running a test case
type CaseReport struct {
	Case     *Case
	err      []Err
	duration time.Duration // The amount of time it took to run the test
	status   TestStatus
	stdout   bytes.Buffer
	stderr   bytes.Buffer
}

func (c *CaseReport) traverse(enter func(report Report) TraverseOption, leave func(report Report) TraverseOption) TraverseOption {
	switch enter(c) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(c)
	}

	return leave(c)
}

func (c *CaseReport) FullNameWithSeparator() string {
	return c.Case.FullNameWithSeparator()
}

func (c *CaseReport) Err() []Err {
	return c.err
}

func (c *CaseReport) Stdout() *bytes.Buffer {
	return &c.stdout
}

func (c *CaseReport) Stderr() *bytes.Buffer {
	return &c.stderr
}

func (c *CaseReport) RegisterErr(err Err) {
	c.err = append(c.err, err)
}

func (c *CaseReport) Duration() time.Duration {
	return c.duration
}

func (c *CaseReport) Status() TestStatus {
	return c.status
}

func NewCaseReport(cas *Case) *CaseReport {
	return &CaseReport{
		Case: cas,
	}
}

func (c *CaseReport) UpdateStatus(newStatus TestStatus) {
	switch newStatus {
	case TEST_ERROR:
		c.status = TEST_ERROR
	case TEST_FAILED:
		if c.status != TEST_ERROR {
			c.status = TEST_FAILED
		}
	case TEST_SUCCESS:
		if c.status == TEST_RUNNING {
			c.status = TEST_SUCCESS
		}
	}
}
