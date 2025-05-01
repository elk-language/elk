package repl

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
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
	typechecker  *checker.Checker
	vm           *vm.VM
	inspectStack bool
}

const replName = "<repl>"

func (e *evaluator) evaluate(input string) {
	if e.typechecker == nil {
		e.typechecker = checker.New()
		e.vm = vm.New()
	}
	fn, dl := e.typechecker.CheckSource(replName, input)

	if dl != nil {
		fmt.Println()

		sourceMap := map[string]string{replName: input}
		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, sourceMap)
		if err != nil {
			panic(err)
		}
		fmt.Println(str)
		isFailure := dl.IsFailure()
		e.typechecker.ClearErrors()
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
	ast, dl := parser.Parse(replName, input)

	if dl != nil {
		fmt.Println()

		sourceMap := map[string]string{replName: input}
		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, sourceMap)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
		if dl.IsFailure() {
			return
		}
	}
	pp.Println(ast)
}

// compiles the input to bytecode and dumps it to the output
func (e *evaluator) disassemble(input string) {
	if e.typechecker == nil {
		e.typechecker = checker.New()
	}
	fn, dl := e.typechecker.CheckSource(replName, input)

	if dl != nil {
		fmt.Println()

		sourceMap := map[string]string{replName: input}
		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, sourceMap)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
		isFailure := dl.IsFailure()
		e.typechecker.ClearErrors()
		if isFailure {
			return
		}
	}

	fn.Disassemble(os.Stdout)
}

// parsers, typechecks the input and prints it to the output
func (e *evaluator) typecheck(input string) {
	if e.typechecker == nil {
		e.typechecker = checker.New()
	}
	_, dl := e.typechecker.CheckSource(replName, input)

	if dl != nil {
		fmt.Println()

		sourceMap := map[string]string{replName: input}
		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, sourceMap)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
		isFailure := dl.IsFailure()
		e.typechecker.ClearErrors()
		if isFailure {
			return
		}
	}

	fmt.Println("OK")
}

// lexes the input and prints it to the output
func (e *evaluator) lex(input string) {
	tokens := lexer.Lex(input)
	pp.Println(tokens)
}

// Start the REPL.
func Run(disassemble, inspectStack, parse, lex, typecheck bool) {
	p := prompt.New(
		executor(disassemble, inspectStack, parse, lex, typecheck),
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

		indentDiff = indentSize - (baseIndent % indentSize)
		if indentDiff > baseIndent {
			indentDiff = baseIndent
		}
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

func executor(disassemble, inspectStack, parse, lex, typecheck bool) prompt.Executor {
	eval := &evaluator{
		inspectStack: inspectStack,
	}
	if lex {
		return eval.lex
	}
	if disassemble {
		return eval.disassemble
	}
	if parse {
		return eval.parse
	}
	if typecheck {
		return eval.typecheck
	}

	return eval.evaluate
}
