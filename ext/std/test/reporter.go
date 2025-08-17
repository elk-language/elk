package test

import (
	"fmt"
	"os"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value"
	"github.com/fatih/color"
)

// Contains the result of running a test suite
type Reporter interface {
	Report(events chan *ReportEvent)
}

type SimpleReporter struct {
	caseCounter    int
	errorCounter   int
	failedCounter  int
	skippedCounter int
	successCounter int
}

func NewSimpleReporter() *SimpleReporter {
	return &SimpleReporter{}
}

func (s *SimpleReporter) Report(events chan *ReportEvent) {
	for event := range events {
		switch event.Type {
		case REPORT_FINISH_SUITE:
			if event.SuiteReport.Suite == RootSuite {
				s.reportFinish(event.SuiteReport)
			}
		case REPORT_FINISH_CASE:
			s.caseCounter++
			switch event.CaseReport.status {
			case TEST_ERROR:
				s.errorCounter++
				fmt.Print(color.RedString("E"))
			case TEST_FAILED:
				s.failedCounter++
				fmt.Print(color.RedString("F"))
			case TEST_SKIPPED:
				s.skippedCounter++
				fmt.Print(color.BlueString("S"))
			case TEST_SUCCESS:
				s.successCounter++
				fmt.Print(color.GreenString("."))
			}
		}
	}
}

func (s *SimpleReporter) reportFinish(report *SuiteReport) {
	switch report.status {
	case TEST_ERROR, TEST_FAILED:
		color.Red("\n\nFailures:\n\n")
		TraverseReport(
			report,
			func(report Report) TraverseOption {
				switch report.Status() {
				case TEST_SKIPPED, TEST_SUCCESS:
					return TraverseSkip
				}

				r, ok := report.(*CaseReport)
				if !ok {
					return TraverseContinue
				}

				fmt.Printf("  %s:\n", r.Case.FullNameWithSeparator())
				switch r.status {
				case TEST_FAILED:
					assertionErr := r.err.AsReference().(*value.Object)
					frame, err := r.stackTrace.Get(-1)
					if !err.IsUndefined() {
						panic(err)
					}

					fmt.Printf(
						"    failure: %s\n    at: %s:%d\n",
						assertionErr.Message().AsString().String(),
						frame.FileName,
						frame.LineNumber,
					)

					if r.stdout.Len() > 0 {
						fmt.Println("\n    --- stdout ---")
						indent.IndentString(os.Stdout, r.stdout.String(), 2)
					}

					if r.stderr.Len() > 0 {
						fmt.Println("\n\n    --- stderr ---")
						indent.IndentString(os.Stdout, r.stderr.String(), 2)
					}
					fmt.Println()
				case TEST_ERROR:
					fmt.Printf(
						"    error: %s\n\n",
						lexer.Colorize(r.err.Inspect()),
					)

					indent.IndentString(os.Stdout, r.stackTrace.String(), 2)

					if r.stdout.Len() > 0 {
						fmt.Println("\n\n    --- stdout ---")
						indent.IndentString(os.Stdout, r.stdout.String(), 2)
					}

					if r.stderr.Len() > 0 {
						fmt.Println("\n\n    --- stderr ---")
						indent.IndentString(os.Stdout, r.stderr.String(), 2)
					}
					fmt.Println()
				}

				return TraverseContinue
			},
			nil,
		)
	}

	fmt.Printf(
		"\n\nSummary: %d cases, %d passed, %d skipped, %d failed, %d errors\n",
		s.caseCounter,
		s.successCounter,
		s.skippedCounter,
		s.failedCounter,
		s.errorCounter,
	)
	fmt.Printf(
		"Finished in %s\n",
		report.duration.String(),
	)
}
