package breakpoint

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser/ast"
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
	prompt.Evaluator
	thread                  *vm.Thread
	checkerContext          *checker.BreakpointContext
	compilerContext         *compiler.BytecodeBreakpointContext
	underscoreVariableIndex int
}

var helpBanner = `
Elk Breakpoint Help
===================
`

var Commands = append(
	[]prompt.Command{
		{Names: []string{":continue", ":c"}, Description: " Exit the breakpoint and continue code execution"},
	},
	prompt.Commands...,
)

func (e *evaluator) evaluate(input string) {
	sourceName := e.AddSource(input)
	defer e.DeleteSource(sourceName)

	var fn *vm.BytecodeFunction
	var dl diagnostic.DiagnosticList

	switch input {
	case ":c", ":continue", ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":s", ":stack":
		e.thread.InspectValueStack()
		return
	case ":disassemble", ":disasm":
		e.thread.Bytecode().DisassembleStdout()
		return
	case ":v", ":variables":
		fn, dl = e.ElkTypechecker.DumpVariablesForBreakpoint(sourceName, e.compilerContext.Location)
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.Prompt.ClearScreen()
		return
	default:
		fn, dl = e.ElkTypechecker.CheckBreakpointSource(sourceName, input, e.DeepCopyPrepend())
	}

	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.SourceMap)
		if err != nil {
			panic(err)
		}
		fmt.Println(str)
		isFailure := dl.IsFailure()
		e.ElkTypechecker.ClearErrors()
		if isFailure {
			return
		}
	}

	programNode, ok := e.ElkTypechecker.ASTCache.Get(sourceName)
	if ok {
		using := e.ElkTypechecker.ExtractTopLevelUsing(programNode)
		e.Prepend = e.usingToPrepend(using)
	}

	startTime := value.TimeNow()
	val, runtimeErr := e.thread.InterpretBreakpoint(fn)
	if !runtimeErr.IsUndefined() {
		e.thread.PrintError()
		e.thread.ResetError()
		return
	}
	duration := value.TimeSince(startTime)
	fmt.Printf("=> %s\n", lexer.Colorize(val.Inspect()))
	fmt.Printf("?: %s\n\n", duration.String())
}

func (e *evaluator) usingToPrepend(using []*ast.UsingExpressionNode) []ast.StatementNode {
	if len(using) < 1 {
		return nil
	}

	result := make([]ast.StatementNode, len(using))

	for i, usingElement := range using {
		copy := ast.DeepCopy(usingElement).(*ast.UsingExpressionNode)
		copy.SetType(nil)

		result[i] = ast.NewExpressionStatementNode(
			copy.Location(),
			copy,
		)
	}

	return result
}

func printHelp() {
	prompt.PrintHelp(helpBanner, Commands)
}

type BreakpointHandler struct{}

func (bh BreakpointHandler) RunBreakpoint(thread *vm.Thread, context value.Value) error {
	return Run(thread, context.AsReference().(*compiler.BytecodeBreakpointContext))
}

func NewPrompt(evaluator *evaluator) (*goprompt.Prompt, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	breakpointHistoryFilePath := filepath.Join(cwd, ".elk_breakpoint_history")
	p := prompt.New(
		&evaluator.Evaluator,
		goprompt.WithPrefix("!> "),
		goprompt.WithPrefixTextColor(goprompt.Yellow),
		goprompt.WithFSHistory(breakpointHistoryFilePath),
	)
	return p, nil
}

// Start the REPL.
func Run(thread *vm.Thread, context *compiler.BytecodeBreakpointContext) error {
	eval := newEvaluator(thread, context)
	p, err := NewPrompt(eval)
	if err != nil {
		return err
	}

	eval.Prompt = p
	p.Run()
	return nil
}

func newEvaluator(thread *vm.Thread, context *compiler.BytecodeBreakpointContext) *evaluator {
	typechecker := checker.NewBreakpointChecker(context)

	printLocation(context.Location)

	eval := &evaluator{
		Commands:        Commands,
		thread:          thread,
		ElkTypechecker:  typechecker,
		checkerContext:  context.TypecheckerContext.(*checker.BreakpointContext),
		compilerContext: context,
		SourceMap:       make(map[string]string),
	}
	eval.Executor = eval.evaluate
	eval.Completer = eval.Complete
	return eval
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
		if endCursor >= len(source) {
			break
		}
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
