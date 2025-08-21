package test

import (
	"bytes"
	"time"
)

type Report interface {
	Duration() time.Duration
	Err() []Err
	Status() TestStatus
	Stdout() *bytes.Buffer
	Stderr() *bytes.Buffer
	FullNameWithSeparator() string
	traverse(enter func(report Report) TraverseOption, leave func(report Report) TraverseOption) TraverseOption
}

// Value used to decide what whether to skip the children of the report,
// break the traversal or continue in the Report Traverse method.
// The zero value continues the traversal.
type TraverseOption uint8

const (
	TraverseContinue TraverseOption = iota
	TraverseSkip
	TraverseBreak
)

func noopTraverseReport(report Report) TraverseOption { return TraverseContinue }

func TraverseReport(report Report, enter func(report Report) TraverseOption, leave func(report Report) TraverseOption) {
	if enter == nil {
		enter = noopTraverseReport
	}
	if leave == nil {
		leave = noopTraverseReport
	}
	report.traverse(enter, leave)
}
