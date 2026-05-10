package breakpoint

import (
	"fmt"
	"os"

	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/repl/prompt"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	goprompt "github.com/elk-language/go-prompt"
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
	case ":v", ":variables":
		fn, dl = e.elkTypechecker.DumpVariablesForBreakpoint(sourceName, e.compilerContext.Location)
	case ":s", ":stack":
		e.thread.InspectValueStack()
		return
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

	eval := &evaluator{
		thread:          thread,
		elkTypechecker:  typechecker,
		checkerContext:  context.TypecheckerContext.(*checker.BreakpointContext),
		compilerContext: context,
		sourceMap:       make(map[string]string),
	}
	return eval.evaluate
}
