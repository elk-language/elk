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
)

type evaluator struct {
	prompt.Evaluator
	ctx           context.Context
	vm            *vm.Thread
	inspectStack  bool
	goTypechecker *compilerTypes.GoTypechecker
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

const underscoreDocComment = "Contains the value returned by the last execution"

func (e *evaluator) evaluate(input string) {
	sourceName := e.AddSource(input)
	defer e.DeleteSource(sourceName)

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
		e.Prompt.ClearScreen()
		return
	case ":v", ":variables":
		fn, dl = e.ElkTypechecker.DumpVariablesForREPL(sourceName, position.DefaultLocation)
	default:
		fn, dl = e.ElkTypechecker.CheckSourceREPLBytecode(sourceName, input, e.DeepCopyPrepend())
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

	programNode, _ := e.ElkTypechecker.ASTCache.Get(sourceName)
	using := e.ElkTypechecker.ExtractTopLevelUsing(programNode)
	e.Prepend = e.UsingToPrepend(using)

	resultType := e.ElkTypechecker.TypeOf(programNode)
	e.ElkTypechecker.DefineLocal(underscoreDocComment, "_", resultType, true, true, nil)

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
	sourceName := e.AddSource(input)
	defer e.DeleteSource(sourceName)

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.Prompt.ClearScreen()
		return
	default:
	}

	ast, dl := parser.Parse(sourceName, input)

	fmt.Printf("%s\n", lexer.Colorize(ast.Inspect()))
	// pp.Println(ast) // print Go structures

	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.SourceMap)
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
	sourceName := e.AddSource(input)
	defer e.DeleteSource(sourceName)

	var fn *vm.BytecodeFunction
	var dl diagnostic.DiagnosticList

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":v", ":variables":
		fn, dl = e.ElkTypechecker.DumpVariablesForREPL(sourceName, position.DefaultLocation)
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.Prompt.ClearScreen()
		return
	default:
		fn, dl = e.ElkTypechecker.CheckSourceREPLBytecode(sourceName, input, e.DeepCopyPrepend())
	}

	using, _ := e.ElkTypechecker.ExtractTopLevelUsingFromFile(sourceName)
	e.Prepend = e.UsingToPrepend(using)

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

	fn.Disassemble(os.Stdout)
}

// compiles the input to Go source code and dumps it to the output
func (e *evaluator) transpile(input string) {
	sourceName := e.AddSource(input)
	defer e.DeleteSource(sourceName)

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.Prompt.ClearScreen()
		return
	default:
	}

	var buff bytes.Buffer
	cmp, dl := checker.CheckSourceNative(sourceName, input, nil, bitfield.BitField32{}, &buff, vm.DefaultThreadPool)

	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.SourceMap)
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
	sourceName := e.AddSource(input)
	defer e.DeleteSource(sourceName)

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.Prompt.ClearScreen()
		return
	default:
	}

	err := compileRunNative(sourceName, input)
	if err != nil {
		if dl, ok := errors.AsType[diagnostic.DiagnosticList](err); ok {
			fmt.Println()

			str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.SourceMap)
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
	sourceName := e.AddSource(input)
	defer e.DeleteSource(sourceName)

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.Prompt.ClearScreen()
		return
	default:
	}

	_, dl := e.ElkTypechecker.CheckSourceREPL(sourceName, input, e.Prepend)
	using, _ := e.ElkTypechecker.ExtractTopLevelUsingFromFile(sourceName)
	e.Prepend = e.UsingToPrepend(using)

	for _, ast := range e.ElkTypechecker.ASTCache.Map {
		if ast.IsBuiltin {
			continue
		}
		fmt.Println(lexer.Colorize(ast.Inspect()))
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

	fmt.Println("OK")
}

// parses, typechecks, expands macros and prints the AST to the output
func (e *evaluator) expand(input string) {
	sourceName := e.AddSource(input)
	defer e.DeleteSource(sourceName)

	switch input {
	case ":q", ":quit", "exit":
		fmt.Println()
		return
	case ":h", ":help", "help":
		printHelp()
		return
	case ":clear", "clear":
		e.Prompt.ClearScreen()
		return
	default:
	}

	_, dl := e.ElkTypechecker.CheckSourceREPL(sourceName, input, e.Prepend)
	using, _ := e.ElkTypechecker.ExtractTopLevelUsingFromFile(sourceName)
	e.Prepend = e.UsingToPrepend(using)

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

	ast, ok := e.ElkTypechecker.ASTCache.GetUnsafe(sourceName)
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
	eval.Prompt = prompt.New(
		&eval.Evaluator,
		goprompt.WithFSHistory(replHistoryFilePath),
	)
	eval.Prompt.Run()

	return nil
}

var helpBanner = `
Elk REPL Help
=============
`

func printHelp() {
	prompt.PrintHelp(helpBanner, prompt.Commands)
}

func newEvaluator(ctx context.Context, opts *Options) *evaluator {
	eval := &evaluator{
		Commands:     prompt.Commands,
		ctx:          ctx,
		inspectStack: opts.InspectStack,
		SourceMap:    make(map[string]string),
		Completer:    goprompt.NoopCompleter,
	}
	eval.setExecutor(opts)
	return eval
}

func (e *evaluator) setExecutor(opts *Options) {
	if opts.Lex {
		e.Executor = e.lex
		info.CurrentMode = info.ReplLex
		return
	}
	if opts.Disassemble {
		e.ElkTypechecker = checker.New()
		e.ElkTypechecker.SetREPL(true)
		e.ElkTypechecker.DefineLocal(underscoreDocComment, "_", types.Nil{}, true, true, nil)

		e.Completer = e.Complete
		e.Executor = e.disassemble
		info.CurrentMode = info.ReplDisassemble
		return
	}
	if opts.Transpile {
		checker, err := compilerTypes.NewGoTypechecker()
		if err != nil {
			panic(fmt.Sprintf("go typechecker error: %s\n", err))
		}

		e.goTypechecker = checker
		e.Executor = e.transpile
		info.CurrentMode = info.ReplTranspile
		return
	}
	if opts.Native {
		e.Executor = e.native
		info.CurrentMode = info.ReplNative
		return
	}
	if opts.Parse {
		e.Executor = e.parse
		info.CurrentMode = info.ReplParse
		return
	}
	if opts.Typecheck {
		e.ElkTypechecker = checker.New()
		e.ElkTypechecker.SetREPL(true)

		e.Completer = e.Complete
		e.Executor = e.typecheck
		info.CurrentMode = info.ReplTypecheck
		return
	}
	if opts.Expand {
		e.ElkTypechecker = checker.New()
		e.ElkTypechecker.SetREPL(true)

		e.Completer = e.Complete
		e.Executor = e.expand
		info.CurrentMode = info.ReplExpand
		return
	}

	e.ElkTypechecker = checker.New()
	e.ElkTypechecker.SetREPL(true)
	e.ElkTypechecker.SetLSP(true)
	e.vm = vm.New()
	e.ElkTypechecker.DefineLocal(underscoreDocComment, "_", types.Nil{}, true, true, nil)

	e.Completer = e.Complete
	e.Executor = e.evaluate
	info.CurrentMode = info.ReplInterpretMode
}
