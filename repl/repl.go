package repl

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/vm"
	"github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
	"github.com/k0kubun/pp"
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

// Start the REPL.
func Run(disassemble bool) {
	p := prompt.New(
		executor(disassemble),
		prompt.WithLexer(&Lexer{}),
		prompt.WithExecuteOnEnterCallback(executeOnEnter),
		prompt.WithPrefix(">> "),
	)
	p.Run()
}

// A Set of keywords that end a block of code
var blockEndKeywords = map[string]bool{
	"end": true,
}

// A Set of keywords that separate multiple blocks of code
var blockSeparatorKeywords = map[string]bool{
	"else":  true,
	"elsif": true,
}

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
	firstWord := firstToken.StringValue()
	blockEnd := blockEndKeywords[firstWord]
	blockSeparator := blockSeparatorKeywords[firstWord]

	if blockEnd || blockSeparator {
		var indentDiff int
		var nextIndentDiff int

		indentDiff = indentSize - (baseIndent % indentSize)
		if indentDiff > baseIndent {
			indentDiff = baseIndent
		}
		if blockEnd {
			nextIndentDiff = indentDiff
		}

		toLeft := pstrings.RuneNumber(utf8.RuneCountInString(currentLine) - baseIndent + indentDiff)
		pr.CursorLeftRunes(toLeft)
		pr.InsertTextMoveCursor(currentLine[baseIndent:], false)
		pr.DeleteRunes(toLeft)
		baseIndent -= nextIndentDiff
	}

	if doc.OnLastLine() {
		if p.ShouldIndent() {
			return baseIndent/indentSize + 1, false
		}
		if p.IsIncomplete() {
			return baseIndent / indentSize, false
		}

		return 0, true
	}

	if p.ShouldIndent() {
		return baseIndent/indentSize + 1, false
	}

	return baseIndent / indentSize, false
}

const (
	sourceName = "REPL"
)

func executor(disassemble bool) prompt.Executor {
	if disassemble {
		return disassembler
	}

	return evaluator
}

func evaluator(input string) {
	chunk, compileErr := compiler.CompileSource(sourceName, input)
	if compileErr != nil {
		fmt.Println()
		fmt.Println(compileErr.HumanStringWithSource(input, true))
		return
	}
	vm := vm.New()
	value, runtimeErr := vm.InterpretBytecode(chunk)
	if runtimeErr != nil {
		panic(runtimeErr)
	}
	fmt.Printf("=> %s\n\n", value.Inspect())
}

// compiles the input to bytecode and dumps it to the output
func disassembler(input string) {
	ast, _ := parser.Parse("", input)
	pp.Println(ast)
	chunk, compileErr := compiler.CompileSource(sourceName, input)
	if compileErr != nil {
		fmt.Println()
		fmt.Println(compileErr.HumanStringWithSource(input, true))
		return
	}

	chunk.Disassemble(os.Stdout)
}
