package repl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/format"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/elk-language/elk"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/compiler/colorize"
	compilerTypes "github.com/elk-language/elk/compiler/types"
	"github.com/elk-language/elk/info"
	"github.com/elk-language/elk/info/banner"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/repl/prompt"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	goprompt "github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
)

type evaluator struct {
	executor                  goprompt.Executor
	ctx                       context.Context
	vm                        *vm.Thread
	inspectStack              bool
	sourceMap                 map[string]string
	inputIndex                int
	elkTypechecker            *checker.Checker
	goTypechecker             *compilerTypes.GoTypechecker
	prepend                   []ast.StatementNode
	prompt                    *goprompt.Prompt
	lastCompletionSource      string
	lastCompletionProgramNode *ast.ProgramNode
}

func (e *evaluator) sourceName() string {
	return fmt.Sprintf("<repl:%d>", e.inputIndex)
}

func (e *evaluator) deleteSource(sourceName string) {
	if e.elkTypechecker != nil && !e.elkTypechecker.DefinedMacros() {
		delete(e.sourceMap, sourceName)
	}
}

func (e *evaluator) addSource(input string) string {
	sourceName := e.sourceName()
	e.inputIndex++
	e.sourceMap[sourceName] = input
	return sourceName
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

func (e *evaluator) programContainsBreakpoint(program *ast.ProgramNode) bool {
	var result bool
	ast.Traverse(
		program,
		func(node, parent ast.Node) ast.TraverseOption {
			if _, ok := node.(*ast.BreakpointNode); ok {
				result = true
				return ast.TraverseBreak
			}

			return ast.TraverseContinue
		},
		nil,
	)

	return result
}

func (e *evaluator) evaluate(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	var fn *vm.BytecodeFunction
	var dl diagnostic.DiagnosticList

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":s", ":stack":
		e.vm.InspectValueStack()
		return
	case ":disassemble", ":disasm":
		e.vm.Bytecode().DisassembleStdout()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.prompt.ClearScreen()
		return
	case ":v", ":variables":
		fn, dl = e.elkTypechecker.DumpVariablesForREPL(sourceName, position.DefaultLocation)
	default:
		fn, dl = e.elkTypechecker.CheckSourceREPLBytecode(sourceName, input, e.deepCopyPrepend())
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

	programNode, _ := e.elkTypechecker.ASTCache.Get(sourceName)
	using := e.elkTypechecker.ExtractTopLevelUsing(programNode)
	e.prepend = e.usingToPrepend(using)

	resultType := e.elkTypechecker.TypeOf(programNode)
	e.elkTypechecker.DefineLocal("_", resultType, true, true)

	eval := func(spinner *prompt.Spinner, cd prompt.SpinnerCaptureData) error {
		executionFinishedCtx, markExecutionFinished := context.WithCancel(e.ctx)
		defer markExecutionFinished()

		vmCtx, abortExecution := context.WithCancel(e.ctx)
		defer abortExecution()

		cancelSignalCtx, cancelSignal := context.WithCancel(e.ctx)
		defer cancelSignal()

		go func() {
			signalCtx, stop := signal.NotifyContext(
				cancelSignalCtx,
				syscall.SIGINT,
				syscall.SIGTERM,
			)
			defer stop()

			<-signalCtx.Done()

			abortExecution()

			select {
			case <-executionFinishedCtx.Done():
			case <-time.After(5 * time.Second):
				fmt.Fprintf(e.vm.Stderr, "\ntimed out waiting for execution to finish\n")
				os.Exit(1)
			}
		}()

		startTime := value.TimeNow()
		e.vm.Aborter = value.NewAborter(vmCtx, abortExecution)

		if spinner != nil {
			spinner.SetMessage("running")

			e.vm.Stdout = cd.NewStdout
			e.vm.Stderr = cd.NewStderr
		}

		returnVal, runtimeErr := e.vm.InterpretREPL(fn)
		duration := value.TimeSince(startTime)
		cancelSignal()
		markExecutionFinished()
		if !runtimeErr.IsUndefined() {
			e.vm.PrintError()
			e.vm.ResetError()
			return nil
		}
		fmt.Printf("=> %s\n", lexer.Colorize(returnVal.Inspect()))
		fmt.Printf("?: %s\n\n", duration.String())

		if e.inspectStack {
			e.vm.InspectValueStack()
		}
		return nil
	}

	if e.programContainsBreakpoint(programNode) {
		err := eval(nil, prompt.SpinnerCaptureData{})
		if err != nil {
			panic(err)
		}
	} else {
		spinner := prompt.NewSpinner(os.Stdout, "compiling")
		spinner.RunWithRestore(
			func(cd prompt.SpinnerCaptureData) error {
				return eval(spinner, cd)
			},
			func(cd prompt.SpinnerCaptureData) {
				e.vm.Stdout = cd.OriginalStdout
				e.vm.Stderr = cd.OriginalStderr
				cd.Restore()
			},
		)
	}
}

// parses the input and prints it to the output
func (e *evaluator) parse(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.prompt.ClearScreen()
		return
	default:
	}

	ast, dl := parser.Parse(sourceName, input)

	fmt.Printf("%s\n", lexer.Colorize(ast.Inspect()))
	// pp.Println(ast) // print Go structures

	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
		if dl.IsFailure() {
			return
		}
	}

}

// compiles the input to bytecode and dumps it to the output
func (e *evaluator) disassemble(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	if e.elkTypechecker == nil {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetREPL(true)
	}

	var fn *vm.BytecodeFunction
	var dl diagnostic.DiagnosticList

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":v", ":variables":
		fn, dl = e.elkTypechecker.DumpVariablesForREPL(sourceName, position.DefaultLocation)
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.prompt.ClearScreen()
		return
	default:
		fn, dl = e.elkTypechecker.CheckSourceREPLBytecode(sourceName, input, e.deepCopyPrepend())
	}

	using, _ := e.elkTypechecker.ExtractTopLevelUsingFromFile(sourceName)
	e.prepend = e.usingToPrepend(using)

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

	fn.Disassemble(os.Stdout)
}

// compiles the input to Go source code and dumps it to the output
func (e *evaluator) transpile(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.prompt.ClearScreen()
		return
	default:
	}

	var buff bytes.Buffer
	cmp, dl := checker.CheckSourceNative(sourceName, input, nil, bitfield.BitField16{}, &buff, vm.DefaultThreadPool)

	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
		if dl.IsFailure() {
			return
		}
	}

	cmp.Flush()
	result, err := format.Source(buff.Bytes())
	if err != nil {
		fmt.Println(buff.String())
		fmt.Fprintf(os.Stderr, "cannot format target go file: %s\n", err)
		return
	}

	_, err = os.Stdout.Write(colorize.ColorizeGo(result))
	if err != nil {
		panic(err)
	}
	fmt.Println()

	err = e.goTypechecker.CheckBytes(result)
	if err != nil {
		panic(err)
	}
}

// compiles the input to Go source code and executes it
func (e *evaluator) native(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.prompt.ClearScreen()
		return
	default:
	}

	err := compileRunNative(sourceName, input)
	if err != nil {
		if dl, ok := errors.AsType[diagnostic.DiagnosticList](err); ok {
			fmt.Println()

			str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
			if err != nil {
				panic(err)
			}

			fmt.Println(str)
		} else {
			fmt.Println(err.Error())
		}
		return
	}
}

// Compile the input to a native binary and execute it.
// A spinner stays on the last line until both phases finish.
func compileRunNative(sourceName, input string) error {
	spinner := prompt.NewSpinner(os.Stdout, "compiling")
	return spinner.Run(func(prompt.SpinnerCaptureData) error {
		binPath, err := elk.CompileSource(sourceName, input)
		if err != nil {
			return err
		}

		spinner.SetMessage("running")
		return elk.RunBinary(binPath)
	})
}

// parses, typechecks the input and prints it to the output
func (e *evaluator) typecheck(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.prompt.ClearScreen()
		return
	default:
	}

	if e.elkTypechecker == nil {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetREPL(true)
	}
	_, dl := e.elkTypechecker.CheckSourceREPL(sourceName, input, e.prepend)
	using, _ := e.elkTypechecker.ExtractTopLevelUsingFromFile(sourceName)
	e.prepend = e.usingToPrepend(using)

	for _, ast := range e.elkTypechecker.ASTCache.Map {
		if ast.IsBuiltin {
			continue
		}
		fmt.Println(lexer.Colorize(ast.Inspect()))
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

	fmt.Println("OK")
}

// parses, typechecks, expands macros and prints the AST to the output
func (e *evaluator) expand(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.prompt.ClearScreen()
		return
	default:
	}

	if e.elkTypechecker == nil {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetREPL(true)
	}
	_, dl := e.elkTypechecker.CheckSourceREPL(sourceName, input, e.prepend)
	using, _ := e.elkTypechecker.ExtractTopLevelUsingFromFile(sourceName)
	e.prepend = e.usingToPrepend(using)

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

	ast, ok := e.elkTypechecker.ASTCache.GetUnsafe(sourceName)
	if !ok {
		panic(fmt.Sprintf("cannot get AST of %s in REPL", sourceName))
	}

	fmt.Println(lexer.Colorize(ast.String()))
}

// lexes the input and prints it to the output
func (e *evaluator) lex(input string) {
	tokens := lexer.Lex(input)
	for _, token := range tokens {
		t := (*value.Token)(token)
		fmt.Printf("%s\n", lexer.Colorize(t.Inspect()))
	}
	fmt.Println()
}

// Handles autocomplete
func (e *evaluator) completer(doc goprompt.Document) (suggestions []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	source := doc.Text
	var programNode *ast.ProgramNode
	if e.lastCompletionSource == source && e.lastCompletionProgramNode != nil {
		programNode = e.lastCompletionProgramNode
	} else {
		programNode = e.elkTypechecker.CheckSourceREPLAutocomplete("<repl>", source, e.deepCopyPrepend())
	}
	e.elkTypechecker.ClearErrors()

	cursorOffset := int(doc.CursorByteOffset())

	var narrowestNode ast.Node
	var narrowestMethodCallNode ast.Node
	ast.Traverse(
		programNode,
		func(node, parent ast.Node) ast.TraverseOption {
			loc := node.Location()
			if cursorOffset < loc.StartPos.ByteOffset || cursorOffset-1 > loc.EndPos.ByteOffset {
				return ast.TraverseSkip
			}

			if narrowestNode == nil {
				narrowestNode = node
				if ast.IsMethodCallNode(node) {
					narrowestMethodCallNode = node
				}
				return ast.TraverseContinue
			}

			if node.Location().ByteCount() <= narrowestNode.Location().ByteCount() {
				narrowestNode = node
				if ast.IsMethodCallNode(node) {
					narrowestMethodCallNode = node
				}
			}

			return ast.TraverseContinue
		},
		nil,
	)

	if narrowestNode == nil {
		return
	}

	log("narrowestNode: %s", narrowestNode.Inspect())
	if narrowestMethodCallNode != nil {
		log(", narrowestMethodNode: %s", narrowestMethodCallNode.Inspect())
	}
	log("\n")

	var receiverNode ast.ExpressionNode
	var nameNode ast.IdentifierNode
	switch n := narrowestMethodCallNode.(type) {
	case *ast.AttributeAccessNode:
		receiverNode = n.Receiver
		nameNode = n.AttributeName
	case *ast.MethodCallNode:
		receiverNode = n.Receiver
		nameNode = n.MethodName
	default:
		return
	}

	receiverType := e.elkTypechecker.ToNonLiteral(e.elkTypechecker.TypeOf(receiverNode), true)
	receiverNamespace, ok := receiverType.(types.Namespace)
	if !ok {
		return
	}

	var currentName string
	switch n := nameNode.(type) {
	case *ast.PublicIdentifierNode:
		currentName = n.Value
	case *ast.PrivateIdentifierNode:
		currentName = n.Value
	}

	nameLocation := nameNode.Location()
	if nameLocation.ByteIsWithin(cursorOffset - 1) {
		startChar = pstrings.RuneNumber(nameLocation.StartPos.RuneOffset(source))
		endChar = pstrings.RuneNumber(nameLocation.EndPos.RuneOffset(source)) + 1
	} else {
		startChar = pstrings.RuneNumber(cursorOffset)
		endChar = pstrings.RuneNumber(cursorOffset)
	}

	log("startChar: %d, endChar: %d\n", startChar, endChar)

	for name, method := range e.elkTypechecker.MethodsInNamespace(receiverNamespace) {
		if strings.Contains(name.String(), currentName) {
			suggestions = append(
				suggestions,
				goprompt.Suggest{
					Text:        name.String(),
					Description: method.DocComment,
				},
			)
		}
	}

	return
}

// For debugging
func log(format string, a ...any) {
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, format+"\n", a...); err != nil {
		panic(err)
	}
}

var _ = log

type Options struct {
	Disassemble  bool
	Transpile    bool
	Native       bool
	InspectStack bool
	Parse        bool
	Lex          bool
	Typecheck    bool
	Expand       bool
}

// Start the REPL.
func Run(ctx context.Context, opts *Options) error {
	eval := newEvaluator(ctx, opts)
	banner.Display()
	elkDir, err := info.MustDotElkDir()
	if err != nil {
		return err
	}

	replHistoryFilePath := filepath.Join(elkDir, ".repl_history")
	eval.prompt = prompt.New(
		eval.executor,
		goprompt.WithFSHistory(replHistoryFilePath),
		goprompt.WithCompleter(eval.completer),
	)
	eval.prompt.Run()

	return nil
}

var helpBanner = `
Elk REPL Help
=============
`

var commandsHelp = `
Available Commands:
  :quit, :q, exit         Exit the program
  :help, :h, help         Print this help message
  :variables, :v          Inspect the currently available variables
  :stack, :s              Inspect the values on the VM's stack
  :disassemble, :disasm   Disassemble the bytecode of the previous input and print it
  :clear, clear           Clear the screen
`

func printHelp() {
	fmt.Print(helpBanner)
	fmt.Print(prompt.KeyBindsHelp)
	fmt.Println()
	fmt.Print(commandsHelp)
	fmt.Println()
}

func newEvaluator(ctx context.Context, opts *Options) *evaluator {
	eval := &evaluator{
		ctx:          ctx,
		inspectStack: opts.InspectStack,
		sourceMap:    make(map[string]string),
	}
	eval.setExecutor(opts)
	return eval
}

func (e *evaluator) setExecutor(opts *Options) {
	if opts.Lex {
		e.executor = e.lex
		info.CurrentMode = info.ReplLex
		return
	}
	if opts.Disassemble {
		e.executor = e.disassemble
		info.CurrentMode = info.ReplDisassemble
		return
	}
	if opts.Transpile {
		checker, err := compilerTypes.NewGoTypechecker()
		if err != nil {
			panic(fmt.Sprintf("go typechecker error: %s\n", err))
		}

		e.goTypechecker = checker
		e.executor = e.transpile
		info.CurrentMode = info.ReplTranspile
		return
	}
	if opts.Native {
		e.executor = e.native
		info.CurrentMode = info.ReplNative
		return
	}
	if opts.Parse {
		e.executor = e.parse
		info.CurrentMode = info.ReplParse
		return
	}
	if opts.Typecheck {
		e.executor = e.typecheck
		info.CurrentMode = info.ReplTypecheck
		return
	}
	if opts.Expand {
		e.executor = e.expand
		info.CurrentMode = info.ReplExpand
		return
	}

	e.elkTypechecker = checker.New()
	e.elkTypechecker.SetREPL(true)
	e.vm = vm.New()
	e.elkTypechecker.DefineLocal("_", types.Nil{}, true, true)

	e.executor = e.evaluate
	info.CurrentMode = info.ReplInterpretMode
}
