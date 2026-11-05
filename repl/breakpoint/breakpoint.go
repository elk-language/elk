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
	thread                  *vm.Thread
	elkTypechecker          *checker.Checker
	checkerContext          *checker.BreakpointContext
	compilerContext         *compiler.BytecodeBreakpointContext
	sourceMap               map[string]string
	inputIndex              int
	prompt                  *goprompt.Prompt
	underscoreVariableIndex int
	prepend                 []ast.StatementNode
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

var helpBanner = `
Elk Breakpoint Help
===================
`

var commandsHelp = `
Available Commands:
  :continue, :c           Exit the breakpoint and continue code execution
  :quit, :q, exit         Exit the program
  :help, :h, help         Print this help message
  :variables, :v          Inspect the currently available variables
  :stack, :s              Inspect the values on the VM's stack
  :disassemble, :disasm   Disassemble the bytecode and print it
`

func (e *evaluator) evaluate(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

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
		fn, dl = e.elkTypechecker.DumpVariablesForBreakpoint(sourceName, e.compilerContext.Location)
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.prompt.ClearScreen()
		return
	default:
		fn, dl = e.elkTypechecker.CheckBreakpointSource(sourceName, input, e.deepCopyPrepend())
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

	programNode, ok := e.elkTypechecker.ASTCache.Get(sourceName)
	if ok {
		using := e.elkTypechecker.ExtractTopLevelUsing(programNode)
		e.prepend = e.usingToPrepend(using)
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

func (e *evaluator) deepCopyPrepend() []ast.StatementNode {
	if len(e.prepend) < 1 {
		return nil
	}

	result := make([]ast.StatementNode, len(e.prepend))

	for i, statement := range e.prepend {
		result[i] = ast.DeepCopy(statement).(ast.StatementNode)
	}

	return result
}

func printHelp() {
	fmt.Print(helpBanner)
	fmt.Print(prompt.KeyBindsHelp)
	fmt.Println()
	fmt.Print(commandsHelp)
	fmt.Println()
}

type BreakpointHandler struct{}

func (bh BreakpointHandler) RunBreakpoint(thread *vm.Thread, context value.Value) error {
	return Run(thread, context.AsReference().(*compiler.BytecodeBreakpointContext))
}

func NewPrompt(executor goprompt.Executor) (*goprompt.Prompt, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	breakpointHistoryFilePath := filepath.Join(cwd, ".elk_breakpoint_history")
	p := prompt.New(
		executor,
		goprompt.WithPrefix("!> "),
		goprompt.WithPrefixTextColor(goprompt.Yellow),
		goprompt.WithFSHistory(breakpointHistoryFilePath),
	)
	return p, nil
}

// Start the REPL.
func Run(thread *vm.Thread, context *compiler.BytecodeBreakpointContext) error {
	eval := newEvaluator(thread, context)
	p, err := NewPrompt(eval.evaluate)
	if err != nil {
		return err
	}

	eval.prompt = p
	p.Run()
	return nil
}

func newEvaluator(thread *vm.Thread, context *compiler.BytecodeBreakpointContext) *evaluator {
	typechecker := checker.NewBreakpointChecker(context)

	printLocation(context.Location)

	return &evaluator{
		thread:          thread,
		elkTypechecker:  typechecker,
		checkerContext:  context.TypecheckerContext.(*checker.BreakpointContext),
		compilerContext: context,
		sourceMap:       make(map[string]string),
	}
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
