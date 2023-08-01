package repl

import (
	"fmt"

	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/vm"
	"github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
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
func Run() {
	p := prompt.New(
		executor,
		prompt.WithLexer(&Lexer{}),
		prompt.WithExecuteOnEnterCallback(executeOnEnter),
	)
	p.Run()
}

const (
	blockEndKeyword = "end"
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
	if len(input) >= 3 && input[len(input)-3:] == blockEndKeyword {
		var indentDiff int
		indentDiff = indentSize - (baseIndent % indentSize)
		if indentDiff > baseIndent {
			indentDiff = baseIndent
		}

		pr.CursorLeftRunes(pstrings.RuneNumber(indentDiff + len(blockEndKeyword)))
		pr.InsertTextMoveCursor(blockEndKeyword, false)
		pr.DeleteRunes(pstrings.RuneNumber(indentDiff + len(blockEndKeyword)))
		baseIndent -= indentSize
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

func executor(input string) {
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
	fmt.Printf("=> %s\n", value.Inspect())
}
