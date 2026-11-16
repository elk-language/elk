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
	"unicode"

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
	executor                         goprompt.Executor
	ctx                              context.Context
	vm                               *vm.Thread
	inspectStack                     bool
	sourceMap                        map[string]string
	inputIndex                       int
	elkTypechecker                   *checker.Checker
	goTypechecker                    *compilerTypes.GoTypechecker
	prepend                          []ast.StatementNode
	prompt                           *goprompt.Prompt
	completer                        goprompt.Completer
	lastCompletionSource             string
	lastCompletionParsedProgramNode  *ast.ProgramNode
	lastCompletionCheckedProgramNode *ast.ProgramNode
}

const sourceNamePrefix = "<repl:"

func (e *evaluator) sourceName() string {
	return fmt.Sprintf("%s%d>", sourceNamePrefix, e.inputIndex)
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
		stmt := ast.DeepCopyGeneric(statement)
		stmt.SetPrepended(true)
		result[i] = stmt
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

const underscoreDocComment = "Contains the value returned by the last execution"

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
	e.elkTypechecker.DefineLocal(underscoreDocComment, "_", resultType, true, true, nil)

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
	cmp, dl := checker.CheckSourceNative(sourceName, input, nil, bitfield.BitField32{}, &buff, vm.DefaultThreadPool)

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

func isLookupNode(node ast.Node) bool {
	switch node.(type) {
	case *ast.MethodCallNode, *ast.ReceiverlessMethodCallNode,
		*ast.AttributeAccessNode, *ast.ConstantLookupNode, *ast.ScopedMacroCallNode,
		*ast.MethodLookupNode:
		return true
	default:
		return false
	}
}

const completionSourceName = "<completion>"

// Handles autocomplete
func (e *evaluator) complete(doc goprompt.Document) (suggestions []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	source := doc.Text
	var parsedProgramNode *ast.ProgramNode
	_ = parsedProgramNode
	var checkedProgramNode *ast.ProgramNode
	if e.lastCompletionSource == source && e.lastCompletionParsedProgramNode != nil {
		parsedProgramNode = e.lastCompletionParsedProgramNode
		checkedProgramNode = e.lastCompletionCheckedProgramNode
	} else {
		parsedProgramNode, checkedProgramNode = e.elkTypechecker.CheckSourceREPLAutocomplete(completionSourceName, source, e.deepCopyPrepend())
		e.lastCompletionSource = source
		e.lastCompletionParsedProgramNode = parsedProgramNode
		e.lastCompletionCheckedProgramNode = checkedProgramNode
	}
	e.elkTypechecker.ClearErrors()

	cursorByteOffset := doc.CursorByteOffset()

	var narrowestParsedNode ast.Node
	var narrowestParsedLookupNode ast.Node
	ast.Traverse(
		parsedProgramNode,
		func(node, parent ast.Node) ast.TraverseOption {
			loc := node.Location()
			if loc == nil || cursorByteOffset < pstrings.ByteNumber(loc.StartPos.ByteOffset) || cursorByteOffset-1 > pstrings.ByteNumber(loc.EndPos.ByteOffset) {
				return ast.TraverseSkip
			}

			if narrowestParsedNode == nil {
				narrowestParsedNode = node
				if isLookupNode(node) {
					narrowestParsedLookupNode = node
				}
				return ast.TraverseContinue
			}

			if node.Location().ByteCount() <= narrowestParsedNode.Location().ByteCount() {
				narrowestParsedNode = node
				if isLookupNode(node) {
					narrowestParsedLookupNode = node
				}
			}

			return ast.TraverseContinue
		},
		nil,
	)

	var narrowestCheckedNode ast.Node
	var narrowestCheckedLookupNode ast.Node
	ast.Traverse(
		checkedProgramNode,
		func(node, parent ast.Node) ast.TraverseOption {
			if stmt, ok := node.(ast.StatementNode); ok && stmt.Prepended() {
				return ast.TraverseSkip
			}

			loc := node.Location()
			if loc == nil || cursorByteOffset < pstrings.ByteNumber(loc.StartPos.ByteOffset) || cursorByteOffset-1 > pstrings.ByteNumber(loc.EndPos.ByteOffset) {
				return ast.TraverseSkip
			}

			if narrowestCheckedNode == nil {
				narrowestCheckedNode = node
				if isLookupNode(node) {
					narrowestCheckedLookupNode = node
				}
				return ast.TraverseContinue
			}

			if node.Location().ByteCount() <= narrowestCheckedNode.Location().ByteCount() {
				narrowestCheckedNode = node
				if isLookupNode(node) {
					narrowestCheckedLookupNode = node
				}
			}

			return ast.TraverseContinue
		},
		nil,
	)

	if narrowestCheckedNode == nil {
		return
	}

	switch n := narrowestParsedLookupNode.(type) {
	case *ast.ConstantLookupNode:
		return e.completeConstantLookup(doc, n)
	case *ast.MethodLookupNode:
		return e.completeMethodLookup(doc, n)
	case *ast.ScopedMacroCallNode:
		return e.completeScopedMacroCall(doc, n)
	case *ast.AttributeAccessNode:
		receiver := getReceiverOfMethodCall(narrowestCheckedLookupNode)
		if receiver == nil {
			return
		}
		return e.completeMethodCall(doc, receiver, n.AttributeName)
	case *ast.MethodCallNode:
		receiver := getReceiverOfMethodCall(narrowestCheckedLookupNode)
		if receiver == nil {
			return
		}
		return e.completeMethodCall(doc, receiver, n.MethodName)
	case *ast.ReceiverlessMethodCallNode:
		return e.completeFunction(doc, n.MethodName)
	}

	switch n := narrowestParsedNode.(type) {
	case *ast.PublicConstantNode:
		return e.completeConstant(doc, n.Value, n.Location())
	case *ast.PrivateConstantNode:
		return e.completeConstant(doc, n.Value, n.Location())
	case *ast.ReceiverlessMethodCallNode:
		return e.completeFunction(doc, n.MethodName)
	case *ast.PublicIdentifierNode:
		return e.completeLocalOrFunction(doc, n.Value, n.Location())
	case *ast.PrivateIdentifierNode:
		return e.completeLocalOrFunction(doc, n.Value, n.Location())
	default:
		return
	}
}

func getReceiverOfMethodCall(node ast.Node) ast.ExpressionNode {
	switch n := node.(type) {
	case *ast.AttributeAccessNode:
		return n.Receiver
	case *ast.MethodCallNode:
		return n.Receiver
	default:
		return nil
	}
}

func (e *evaluator) completeConstant(doc goprompt.Document, name string, loc *position.Location) (suggestions []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	cursorByteOffset := doc.CursorByteOffset()
	var currentName string
	currentName, startChar, endChar, _ = getCompletionStartAndEnd(doc, name, loc)

	node, _ := e.elkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := node.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := lspData.GetConstantCompletions(e.elkTypechecker, currentName, int(cursorByteOffset))
	appendCompletionsToSuggestions(&suggestions, completions, completionTextFieldValue)

	return suggestions, startChar, endChar
}

func constantToString(constant ast.ComplexConstantNode) string {
	switch constant := constant.(type) {
	case *ast.PublicConstantNode:
		return constant.Value
	case *ast.PrivateConstantNode:
		return constant.Value
	}

	return ""
}

func (e *evaluator) completeConstantLookup(doc goprompt.Document, node *ast.ConstantLookupNode) (suggestions []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	name := constantToString(node.Right)

	var currentName string
	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, name, node.Right.Location())
	if !cursorIsWithin {
		return
	}

	program, _ := e.elkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := program.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := e.elkTypechecker.GetConstantLookupCompletions(node.Left, currentName)
	appendCompletionsToSuggestions(&suggestions, completions, completionTextFieldValue)

	if currentName == "" {
		completions := e.elkTypechecker.GetScopedMacroCompletions(node.Left, currentName)
		appendCompletionsToSuggestions(&suggestions, completions, completionTextFieldValue)
	}

	return suggestions, startChar, endChar
}

func (e *evaluator) completeMethodLookup(doc goprompt.Document, node *ast.MethodLookupNode) (suggestions []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	name := ast.IdentifierToString(node.Name)

	var currentName string
	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, name, node.Name.Location())
	if !cursorIsWithin {
		return
	}

	program, _ := e.elkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := program.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := e.elkTypechecker.GetScopedMacroCompletions(node.Receiver, currentName)
	appendCompletionsToSuggestions(&suggestions, completions, completionTextFieldValue)
	return suggestions, startChar, endChar
}

func (e *evaluator) completeScopedMacroCall(doc goprompt.Document, node *ast.ScopedMacroCallNode) (suggestions []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	name := ast.IdentifierToString(node.MacroName)

	var currentName string
	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, name, node.MacroName.Location())
	if !cursorIsWithin {
		return
	}

	program, _ := e.elkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := program.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := e.elkTypechecker.GetScopedMacroCompletions(node.Receiver, currentName)
	appendCompletionsToSuggestions(&suggestions, completions, completionTextFieldValue)

	return suggestions, startChar, endChar
}

func (e *evaluator) completeFunction(doc goprompt.Document, nameNode ast.IdentifierNode) (suggestions []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	name := ast.IdentifierToString(nameNode)
	cursorByteOffset := doc.CursorByteOffset()

	var currentName string
	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, name, nameNode.Location())
	if !cursorIsWithin {
		return
	}

	node, _ := e.elkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := node.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := lspData.GetFunctionCompletions(e.elkTypechecker, currentName, int(cursorByteOffset))
	appendCompletionsToSuggestions(&suggestions, completions, completionTextFieldLabel)

	return suggestions, startChar, endChar
}

func (e *evaluator) completeLocalOrFunction(doc goprompt.Document, name string, loc *position.Location) (suggestions []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	cursorByteOffset := doc.CursorByteOffset()
	var currentName string
	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, name, loc)
	if !cursorIsWithin {
		return
	}

	node, _ := e.elkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := node.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := lspData.GetLocalCompletions(currentName, int(cursorByteOffset))
	appendCompletionsToSuggestions(&suggestions, completions, completionTextFieldValue)

	completions = lspData.GetFunctionCompletions(e.elkTypechecker, currentName, int(cursorByteOffset))
	appendCompletionsToSuggestions(&suggestions, completions, completionTextFieldLabel)

	return suggestions, startChar, endChar
}

type completionTextField uint8

const (
	completionTextFieldValue completionTextField = iota
	completionTextFieldLabel
	completionTextFieldIdentifier
)

func appendCompletionsToSuggestions(target *[]goprompt.Suggest, source []*checker.Completion, completionTextField completionTextField) {
	for _, completion := range source {
		var text string
		switch completionTextField {
		case completionTextFieldValue:
			text = completion.Value
		case completionTextFieldLabel:
			text = completion.Label
		case completionTextFieldIdentifier:
			text = completion.Identifier
		}

		*target = append(*target, goprompt.Suggest{
			Label:       completion.Label,
			Text:        text,
			Description: normaliseSuggestionDescription(completion.Description),
		})
	}
}

func normaliseSuggestionDescription(desc string) string {
	var buf bytes.Buffer

	var lastChar rune
	for _, char := range desc {
		if unicode.IsSpace(char) {
			goto cont
		}
		if unicode.IsSpace(lastChar) {
			buf.WriteRune(' ')
		}

		buf.WriteRune(char)
	cont:
		lastChar = char
	}

	return buf.String()
}

func (e *evaluator) completeMethodCall(doc goprompt.Document, receiverNode, nameNode ast.Node) (suggestions []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	receiverType := e.elkTypechecker.TypeOf(receiverNode)

	var currentName string
	switch n := nameNode.(type) {
	case *ast.PublicIdentifierNode:
		currentName = n.Value
	case *ast.PrivateIdentifierNode:
		currentName = n.Value
	}

	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, currentName, nameNode.Location())
	if !cursorIsWithin {
		return
	}

	completions := e.elkTypechecker.GetMethodCompletions(receiverType, currentName)
	appendCompletionsToSuggestions(&suggestions, completions, completionTextFieldValue)

	return suggestions, startChar, endChar
}

func getCompletionStartAndEnd(doc goprompt.Document, name string, loc *position.Location) (currentName string, startChar, endChar pstrings.RuneNumber, cursorIsWithin bool) {
	currentName = name
	source := doc.Text
	cursorRuneOffset := doc.CursorRuneOffset()
	cursorByteOffset := doc.CursorByteOffset()

	if name == "" {
		cursorIsWithin = true
	}

	if loc.ByteIsWithin(int(cursorByteOffset - 1)) {
		startChar = pstrings.RuneNumber(loc.StartPos.RuneOffset(source))
		endChar = pstrings.RuneNumber(loc.EndPos.RuneOffset(source)) + 1
		truncate := loc.EndPos.ByteOffset - int(cursorByteOffset) + 1
		if truncate > 0 {
			currentName = currentName[:len(currentName)-truncate]
		}
		cursorIsWithin = true
	} else if cursorByteOffset > pstrings.ByteNumber(loc.EndPos.ByteOffset) {
		startChar = pstrings.RuneNumber(loc.StartPos.RuneOffset(source))
		endChar = pstrings.RuneNumber(loc.EndPos.RuneOffset(source)) + 1
	} else {
		startChar = cursorRuneOffset
		endChar = cursorRuneOffset
	}

	return
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
		completer:    goprompt.NoopCompleter,
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
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetREPL(true)
		e.elkTypechecker.DefineLocal(underscoreDocComment, "_", types.Nil{}, true, true, nil)

		e.completer = e.complete
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
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetREPL(true)

		e.completer = e.complete
		e.executor = e.typecheck
		info.CurrentMode = info.ReplTypecheck
		return
	}
	if opts.Expand {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetREPL(true)

		e.completer = e.complete
		e.executor = e.expand
		info.CurrentMode = info.ReplExpand
		return
	}

	e.elkTypechecker = checker.New()
	e.elkTypechecker.SetREPL(true)
	e.elkTypechecker.SetLSP(true)
	e.vm = vm.New()
	e.elkTypechecker.DefineLocal(underscoreDocComment, "_", types.Nil{}, true, true, nil)

	e.completer = e.complete
	e.executor = e.evaluate
	info.CurrentMode = info.ReplInterpretMode
}
