package repl

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"unicode/utf8"

	"github.com/elk-language/elk"
	"github.com/elk-language/elk/compiler/colorize"
	"github.com/elk-language/elk/compiler/types"
	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/vm"
	"github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
	"github.com/k0kubun/pp/v3"
)

// Adapter for `lexer.Lexer` that
// provides an interface compatible with go-prompt.
type Lexer struct {
	lexer.Lexer
}

func (l *Lexer) Init(input string) {
	l.Lexer = *lexer.New(input)
}

func (l *Lexer) Next() (prompt.Token, bool) {
	t := l.Lexer.Next()
	if t.Type == token.END_OF_FILE {
		return nil, false
	}

	return t, true
}

type evaluator struct {
	vm             *vm.Thread
	inspectStack   bool
	sourceMap      map[string]string
	inputIndex     int
	elkTypechecker *checker.Checker
	goTypechecker  *types.GoTypechecker
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

func (e *evaluator) evaluate(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	if e.elkTypechecker == nil {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetIncremental(true)
		e.vm = vm.New()
	}

	fn, dl := e.elkTypechecker.CheckSourceBytecode(sourceName, input)

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

	value, runtimeErr := e.vm.InterpretREPL(fn)
	if !runtimeErr.IsUndefined() {
		e.vm.PrintError()
		e.vm.ResetError()
		return
	}
	fmt.Printf("=> %s\n\n", lexer.Colorize(value.Inspect()))

	if e.inspectStack {
		e.vm.InspectValueStack()
	}
}

// parses the input and prints it to the output
func (e *evaluator) parse(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	ast, dl := parser.Parse(sourceName, input)

	pp.Println(ast)
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
		e.elkTypechecker.SetIncremental(true)
	}
	fn, dl := e.elkTypechecker.CheckSourceBytecode(sourceName, input)

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

	var buff bytes.Buffer
	cmp, dl := checker.CheckSourceNative(sourceName, input, nil, &buff, vm.DefaultThreadPool)

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

	err := elk.CompileRunSource(sourceName, input)
	if err != nil {
		var dl diagnostic.DiagnosticList
		if errors.As(err, &dl) {
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

// parsers, typechecks the input and prints it to the output
func (e *evaluator) typecheck(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	if e.elkTypechecker == nil {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetIncremental(true)
	}
	_, dl := e.elkTypechecker.CheckSource(sourceName, input)

	for _, ast := range e.elkTypechecker.ASTCache.Map {
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

// parsers, typechecks, expands macros and prints the AST to the output
func (e *evaluator) expand(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	if e.elkTypechecker == nil {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetIncremental(true)
	}
	_, dl := e.elkTypechecker.CheckSource(sourceName, input)

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
	pp.Println(tokens)
}

// Start the REPL.
func Run(disassemble, transpile, native, inspectStack, parse, lex, typecheck, expand bool) {
	p := prompt.New(
		executor(disassemble, transpile, native, inspectStack, parse, lex, typecheck, expand),
		prompt.WithLexer(&Lexer{}),
		prompt.WithExecuteOnEnterCallback(executeOnEnter),
		prompt.WithPrefix(">> "),
	)
	p.Run()
}

// A Set of keywords that end a block of code
var blockEndKeywords = ds.MakeSet(
	"end",
)

// A Set of keywords that separate multiple blocks of code
var blockSeparatorKeywords = ds.MakeSet(
	"else",
	"elsif",
	"case",
	"catch",
	"finally",
)

// Callback triggered when the Enter key is pressed.
// Decides whether the input is complete and should be executed
// or whether a newline with indentation should be added to the buffer.
func executeOnEnter(pr *prompt.Prompt, indentSize int) (indent int, execute bool) {
	doc := pr.Buffer().Document()
	var input string
	if doc.OnLastLine() {
		input = doc.Text
	} else {
		input = doc.TextBeforeCursor()
	}

	p := parser.New(sourceName, input)
	p.Parse()

	baseIndent := doc.CurrentLineIndentSpaces()
	currentLine := doc.CurrentLine()
	lex := lexer.New(currentLine)
	firstToken := lex.Next()
	firstWord := firstToken.FetchValue()
	isBlockEnd := blockEndKeywords.Contains(firstWord)
	isBlockSeparator := blockSeparatorKeywords.Contains(firstWord)

	var movedBack bool
	if isBlockEnd || isBlockSeparator {
		var indentDiff int
		var nextIndentDiff int

		indentDiff = min(indentSize-(baseIndent%indentSize), baseIndent)
		if isBlockEnd {
			nextIndentDiff = indentDiff
		}

		if indentDiff != 0 {
			movedBack = true
		}
		toLeft := pstrings.RuneNumber(utf8.RuneCountInString(currentLine) - baseIndent + indentDiff)
		pr.CursorLeftRunes(toLeft)
		pr.InsertTextMoveCursor(currentLine[baseIndent:], false)
		pr.DeleteRunes(toLeft)
		baseIndent -= nextIndentDiff
	}

	indent = baseIndent / indentSize
	if doc.OnLastLine() {
		if p.ShouldIndent() && !movedBack {
			return indent + 1, false
		}
		if p.IsIncomplete() {
			return indent, false
		}

		return 0, true
	}

	if p.ShouldIndent() && !movedBack {
		return indent + 1, false
	}

	return indent, false
}

const (
	sourceName = "REPL"
)

func executor(disassemble, transpile, native, inspectStack, parse, lex, typecheck, expand bool) prompt.Executor {
	eval := &evaluator{
		inspectStack: inspectStack,
		sourceMap:    make(map[string]string),
	}
	if lex {
		return eval.lex
	}
	if disassemble {
		return eval.disassemble
	}
	if transpile {
		checker, err := types.NewGoTypechecker()
		if err != nil {
			panic(fmt.Sprintf("go typechecker error: %s\n", err))
		}

		eval.goTypechecker = checker
		return eval.transpile
	}
	if native {
		return eval.native
	}
	if parse {
		return eval.parse
	}
	if typecheck {
		return eval.typecheck
	}
	if expand {
		return eval.expand
	}

	return eval.evaluate
}
