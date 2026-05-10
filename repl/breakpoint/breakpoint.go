package breakpoint

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"

	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/repl/prompt"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	goprompt "github.com/elk-language/go-prompt"
	"github.com/fatih/color"
)

func init() {
	vm.BREAKPOINT_HANDLER = BreakpointHandler{}
}

type evaluator struct {
	thread          *vm.Thread
	elkTypechecker  *checker.Checker
	checkerContext  *checker.BreakpointContext
	compilerContext *compiler.BytecodeBreakpointContext
	sourceMap       map[string]string
	inputIndex      int
}

func (e *evaluator) sourceName() string {
	return fmt.Sprintf("<repl:%d>", e.inputIndex)
}

func (e *evaluator) deleteSource(sourceName string) {
	if e.elkTypechecker != nil {
		delete(e.sourceMap, sourceName)
	}
}

func (e *evaluator) addSource(input string) string {
	sourceName := e.sourceName()
	e.inputIndex++
	e.sourceMap[sourceName] = input
	return sourceName
}

func (e *evaluator) evaluate(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	var fn *vm.BytecodeFunction
	var dl diagnostic.DiagnosticList

	switch input {
	case ":c", ":continue", ":q", ":quit":
		fmt.Println()
		return
	case ":s", ":stack":
		e.thread.InspectValueStack()
		return
	case ":v", ":variables":
		fn, dl = e.elkTypechecker.DumpVariablesForBreakpoint(sourceName, e.compilerContext.Location)
	default:
		fn, dl = e.elkTypechecker.CheckBreakpointSource(sourceName, input)
	}

	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
		if err != nil {
			panic(err)
		}
		fmt.Println(str)
		isFailure := dl.IsFailure()
		e.elkTypechecker.ClearErrors()
		if isFailure {
			return
		}
	}

	value, runtimeErr := e.thread.InterpretBreakpoint(fn)
	if !runtimeErr.IsUndefined() {
		e.thread.PrintError()
		e.thread.ResetError()
		return
	}
	fmt.Printf("=> %s\n\n", lexer.Colorize(value.Inspect()))
}

func exitChecker(input string, breakline bool) bool {
	if !breakline {
		return false
	}

	switch input {
	case ":c", ":continue":
		return true
	case ":q", ":quit":
		os.Exit(0)
		return true
	default:
		return false
	}
}

type BreakpointHandler struct{}

func (bh BreakpointHandler) RunBreakpoint(thread *vm.Thread, context value.Value) {
	Run(thread, context.AsReference().(*compiler.BytecodeBreakpointContext))
}

// Start the REPL.
func Run(thread *vm.Thread, context *compiler.BytecodeBreakpointContext) {
	prompt.Run(
		executor(thread, context),
		goprompt.WithPrefix("!> "),
		goprompt.WithPrefixTextColor(goprompt.Yellow),
		goprompt.WithExitChecker(exitChecker),
	)
}

func executor(thread *vm.Thread, context *compiler.BytecodeBreakpointContext) goprompt.Executor {
	typechecker := checker.NewBreakpointChecker(context)

	printLocation(context.Location)

	eval := &evaluator{
		thread:          thread,
		elkTypechecker:  typechecker,
		checkerContext:  context.TypecheckerContext.(*checker.BreakpointContext),
		compilerContext: context,
		sourceMap:       make(map[string]string),
	}
	return eval.evaluate
}

func printLocation(loc *position.Location) {
	if loc == nil {
		return
	}

	headerColor := color.New(color.Bold, color.FgYellow, color.Underline)
	currentLineColor := color.New(color.Bold, color.FgMagenta)
	faintColor := color.New(color.Faint)
	fmt.Println()
	headerColor.Printf("Breakpoint at %s\n", loc.String())
	source, err := os.ReadFile(loc.FilePath)
	if err != nil {
		return
	}

	locStartLine := loc.StartPos.Line
	startLine := max(locStartLine-5, 1)
	endLine := loc.EndPos.Line + 4
	currentLine := loc.StartPos.Line
	startCursor := loc.StartPos.ByteOffset
	endCursor := loc.StartPos.ByteOffset
	lines := make(map[int]string)

	for {
		if startCursor < 0 || source[startCursor] == '\n' {
			lines[currentLine] = string(source[startCursor+1 : endCursor+1])

			currentLine--
			endCursor = startCursor
			startCursor--
			if currentLine < startLine {
				break
			}
			continue
		}

		startCursor--
	}

	startCursor = loc.StartPos.ByteOffset + 1
	endCursor = loc.StartPos.ByteOffset + 1
	currentLine = loc.StartPos.Line

	for {
		if endCursor == len(source)-1 || source[endCursor] == '\n' {
			lines[currentLine] += string(source[startCursor : endCursor+1])

			currentLine++
			endCursor++
			startCursor = endCursor
			if currentLine > endLine {
				break
			}
			continue
		}

		endCursor++
	}

	lineNumbers := slices.Collect(maps.Keys(lines))
	slices.Sort(lineNumbers)
	startLineStr := strconv.Itoa(startLine)
	for _, lineNumber := range lineNumbers {
		line := lines[lineNumber]

		format := fmt.Sprintf(" %%%dd | ", len(startLineStr)+1)
		if lineNumber == locStartLine {
			currentLineColor.Printf(format, lineNumber)
		} else {
			faintColor.Printf(format, lineNumber)
		}
		fmt.Print(lexer.Colorize(line))
	}
}
