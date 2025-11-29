package test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value"
	"github.com/fatih/color"
)

type PlainReporter struct {
	caseCounter    int
	errorCounter   int
	failedCounter  int
	skippedCounter int
	successCounter int
}

func NewPlainReporter() *PlainReporter {
	return &PlainReporter{}
}

func ListenForInterrupt(
	shutdown context.CancelFunc,
) {
	// Create context that listens for the interrupt signal from the OS.
	signalCtx, signalCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer signalCancel()

	// Wait for the interrupt signal.
	<-signalCtx.Done()

	signalCancel() // Stop listening for interrupt signals, allow Ctrl+C to force shutdown
	shutdown()
}

func (s *PlainReporter) Report(events chan *ReportEvent, shutdown context.CancelFunc) {
	go ListenForInterrupt(shutdown)

	for event := range events {
		switch event.Type {
		case REPORT_FINISH_SUITE:
			hasBeforeAllErr := slices.ContainsFunc(
				event.SuiteReport.err,
				func(e Err) bool {
					return e.Typ == ErrBeforeAll
				},
			)

			if hasBeforeAllErr {
				suite := event.SuiteReport.Suite
				suiteCases := suite.CaseCount()
				s.caseCounter += suiteCases
				switch event.SuiteReport.Status() {
				case TEST_ERROR:
					s.errorCounter += suiteCases
					fmt.Print(color.RedString(strings.Repeat("E", suiteCases)))
				case TEST_FAILED:
					s.failedCounter += suiteCases
					fmt.Print(color.RedString(strings.Repeat("F", suiteCases)))
				}
			}

			if event.SuiteReport.Suite == RootSuite {
				s.finishReport(event.SuiteReport)
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

func (s *PlainReporter) finishReport(report *SuiteReport) {
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

				s.reportFailure(report)

				return TraverseContinue
			},
			nil,
		)
	}

	fmt.Printf(
		"\nFinished in %s\n",
		report.duration.String(),
	)
	fmt.Printf(
		"Summary: %d cases, %d passed, %d skipped, %d failed, %d errors\n",
		s.caseCounter,
		s.successCounter,
		s.skippedCounter,
		s.failedCounter,
		s.errorCounter,
	)
}

func (s *PlainReporter) reportFailure(report Report) {
	if len(report.Err()) == 0 {
		return
	}

	fmt.Printf("%s:\n", report.FullNameWithSeparator())
	switch report.Status() {
	case TEST_FAILED, TEST_ERROR:
		for i, testErr := range report.Err() {
			if i != 0 {
				fmt.Println()
			}

			fmt.Printf("   %s:\n", testErr.Typ.String())
			if value.IsA(testErr.Err, AssertionErrorClass) {
				assertionErr := testErr.Err.AsReference().(*value.Object)
				frame, err := testErr.StackTrace.Get(-1)
				if !err.IsUndefined() {
					panic(err)
				}

				fmt.Printf(
					"    failure: %s\n    at: %s:%d",
					lexer.ColorizeEmbellishedText(assertionErr.Message().AsString().String()),
					frame.FileName,
					frame.LineNumber,
				)
			} else if value.IsA(testErr.Err, value.ErrorClass) {
				err := testErr.Err.AsReference().(*value.Object)
				fmt.Printf(
					"    error: %s,\n    message: %s\n",
					lexer.Colorize(testErr.Err.Class().Name),
					lexer.ColorizeEmbellishedText(err.Message().AsString().String()),
				)
				indent.IndentString(os.Stdout, testErr.StackTrace.String(), 2)
			} else {
				fmt.Printf(
					"    error: %s\n",
					lexer.Colorize(testErr.Err.Inspect()),
				)
				indent.IndentString(os.Stdout, testErr.StackTrace.String(), 2)
			}
		}

		stdout := report.Stdout()
		if stdout.Len() > 0 {
			fmt.Println("\n\n    --- stdout ---")
			indent.IndentString(os.Stdout, stdout.String(), 2)
		}

		stderr := report.Stdout()
		if stderr.Len() > 0 {
			fmt.Println("\n\n    --- stderr ---")
			indent.IndentString(os.Stdout, stderr.String(), 2)
		}
		fmt.Print("\n\n")
	}
}
