package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value"
	"github.com/fatih/color"
)

const (
	padding  = 2
	maxWidth = 80
)

func main() {
	m := &RichReporter{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type RichReporter struct {
	progress          progress.Model
	spinner           spinner.Model
	events            chan *ReportEvent
	totalCaseCount    int
	finishedCaseCount int
	errorCounter      int
	failedCounter     int
	skippedCounter    int
	successCounter    int
	quit              bool
	startTime         time.Time
	interrupted       bool
	shutdown          context.CancelFunc
}

const successBarColor = "#27F57D"
const failBarColor = "#FC1249"

func NewRichReporter() *RichReporter {
	return &RichReporter{
		progress: progress.New(
			progress.WithSolidFill(successBarColor),
		),
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(successBarColor))),
		),
		startTime: time.Now(),
	}
}

func (r *RichReporter) Init() tea.Cmd {
	return tea.Batch(waitForEvent(r.events), r.spinner.Tick)
}

type finishRichReport struct{}

func finishRichReportCmd() tea.Msg {
	return finishRichReport{}
}

func (r *RichReporter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.progress.Width = min(msg.Width-padding*2-4, maxWidth)
		return r, nil
	case finishRichReport:
		r.quit = true
		return r, tea.Quit
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			r.shutdown()
			r.interrupted = true
			return r, finishRichReportCmd
		}
		return r, nil
	case *ReportEvent:
		switch msg.Type {
		case REPORT_FINISH_SUITE:
			if msg.SuiteReport.Suite == RootSuite {
				cmd := r.progress.SetPercent(1.0)
				return r, cmd
			}
		case REPORT_FINISH_CASE:
			r.finishedCaseCount++

			var printCmd tea.Cmd
			switch msg.CaseReport.status {
			case TEST_ERROR:
				printCmd = r.reportError(msg.CaseReport)
				r.setFailed()
				r.errorCounter++
			case TEST_FAILED:
				printCmd = r.reportFailure(msg.CaseReport)
				r.setFailed()
				r.failedCounter++
			case TEST_SKIPPED:
				r.skippedCounter++
			case TEST_SUCCESS:
				r.successCounter++
			}

			cmd := r.progress.SetPercent(float64(r.finishedCaseCount) / float64(r.totalCaseCount))
			return r, tea.Batch(waitForEvent(r.events), cmd, printCmd)
		}

		return r, waitForEvent(r.events)
	case spinner.TickMsg:
		var cmd tea.Cmd
		r.spinner, cmd = r.spinner.Update(msg)
		return r, cmd
	case progress.FrameMsg:
		// FrameMsg is sent when the progress bar wants to animate itself
		progressModel, cmd := r.progress.Update(msg)
		r.progress = progressModel.(progress.Model)
		if !r.progress.IsAnimating() && r.progress.Percent() >= 1.0 {
			return r, finishRichReportCmd
		}
		return r, cmd

	default:
		return r, nil
	}
}

func (r *RichReporter) setFailed() {
	r.progress.FullColor = failBarColor
	r.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(failBarColor))
}

// Cmd helper: wait for next value from channel
func waitForEvent(ch <-chan *ReportEvent) tea.Cmd {
	return func() tea.Msg {
		d, ok := <-ch
		if !ok {
			return tea.Quit
		}
		return d
	}
}

func (r *RichReporter) View() string {
	var result strings.Builder

	result.WriteString("\n ")
	result.WriteString(r.spinner.View())
	result.WriteString(" ")
	result.WriteString(r.progress.View())
	result.WriteString("\n\n")
	duration := r.Duration()
	eta := r.ETA()
	fmt.Fprintf(
		&result,
		" [%d/%d] Time: %02d:%02d:%02d, ETA: %02d:%02d:%02d",
		r.finishedCaseCount,
		r.totalCaseCount,
		int(duration.Hours()),
		int(duration.Minutes()),
		int(duration.Seconds()),
		int(eta.Hours()),
		int(eta.Minutes()),
		int(eta.Seconds()),
	)

	result.WriteByte('\n')

	if r.quit {
		var failColor *color.Color
		var successColor *color.Color

		if r.interrupted {
			failColor = color.New(color.FgRed)
			successColor = failColor

			fmt.Fprintf(
				&result,
				"\n\nInterrupted after %s\n",
				duration.String(),
			)
		} else {
			failColor = color.New(color.FgRed)
			successColor = color.New(color.FgGreen)

			fmt.Fprintf(
				&result,
				"\n\nFinished in %s\n",
				duration.String(),
			)
		}

		fmt.Fprintf(
			&result,
			"Summary: %d cases, ",
			r.totalCaseCount,
		)

		if r.failedCounter > 0 || r.errorCounter > 0 {
			fmt.Fprintf(
				&result,
				"%d passed, %d skipped, ",
				r.successCounter,
				r.skippedCounter,
			)
			failColor.Fprintf(
				&result,
				"%d failed, %d errors\n",
				r.failedCounter,
				r.errorCounter,
			)
		} else {
			successColor.Fprintf(
				&result,
				"%d passed, %d skipped\n",
				r.successCounter,
				r.skippedCounter,
			)
		}
	}

	return result.String()
}

func (r *RichReporter) Percent() float64 {
	return float64(r.finishedCaseCount) / float64(r.totalCaseCount)
}

func (r *RichReporter) Duration() time.Duration {
	return time.Since(r.startTime)
}

func (r *RichReporter) ETA() time.Duration {
	dur := r.Duration()
	percent := r.Percent()
	totalEstimated := 1 / percent * dur.Seconds()
	eta := time.Duration(totalEstimated)*time.Second - dur
	if eta < 0 {
		return 0
	}
	return eta
}

func (r *RichReporter) Report(events chan *ReportEvent, shutdown context.CancelFunc) {
	r.events = events
	r.shutdown = shutdown

	TraverseSuite(
		RootSuite,
		func(test SuiteOrCase) TraverseOption {
			switch test.(type) {
			case *Case:
				r.totalCaseCount++
			}

			return TraverseContinue
		},
		nil,
	)

	if _, err := tea.NewProgram(r).Run(); err != nil {
		os.Exit(1)
	}
}

func (r *RichReporter) reportFailure(report *CaseReport) tea.Cmd {
	var result strings.Builder
	result.WriteByte('\n')
	fmt.Fprintf(&result, " %s%s:\n", color.RedString("ERR"), report.Case.FullNameWithSeparator())

	assertionErr := report.err.AsReference().(*value.Object)
	frame, err := report.stackTrace.Get(-1)
	if !err.IsUndefined() {
		panic(err)
	}

	fmt.Fprintf(
		&result,
		"    failure: %s\n    took: %s\n    at: %s:%d\n",
		lexer.ColorizeEmbellishedText(assertionErr.Message().AsString().String()),
		report.duration,
		frame.FileName,
		frame.LineNumber,
	)

	if report.stdout.Len() > 0 {
		fmt.Fprintln(&result, "\n\n    --- stdout ---")
		indent.IndentString(&result, report.stdout.String(), 2)
	}

	if report.stderr.Len() > 0 {
		fmt.Fprintln(&result, "\n\n    --- stderr ---")
		indent.IndentString(&result, report.stderr.String(), 2)
	}

	return tea.Println(result.String())
}

func (r *RichReporter) reportError(report *CaseReport) tea.Cmd {
	var result strings.Builder
	result.WriteByte('\n')
	fmt.Fprintf(&result, " %s%s:\n", color.RedString("FAIL"), report.Case.FullNameWithSeparator())

	if value.IsA(report.err, value.ErrorClass) {
		err := report.err.AsReference().(*value.Object)
		fmt.Fprintf(
			&result,
			"    error: %s,\n    message: %s\n    took: %s\n\n",
			lexer.Colorize(report.err.Class().Name),
			lexer.ColorizeEmbellishedText(err.Message().AsString().String()),
			report.duration,
		)
	} else {
		fmt.Fprintf(
			&result,
			"    error: %s\n\n",
			lexer.Colorize(report.err.Inspect()),
		)
	}

	indent.IndentString(&result, report.stackTrace.String(), 2)

	if report.stdout.Len() > 0 {
		fmt.Fprintln(&result, "\n\n    --- stdout ---")
		indent.IndentString(&result, report.stdout.String(), 2)
	}

	if report.stderr.Len() > 0 {
		fmt.Fprintln(&result, "\n\n    --- stderr ---")
		indent.IndentString(&result, report.stderr.String(), 2)
	}

	return tea.Println(result.String())
}
