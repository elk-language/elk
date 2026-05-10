package prompt

import (
	"unicode/utf8"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/token"
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
func Run(executor prompt.Executor, opts ...prompt.Option) {
	opts = append(
		[]prompt.Option{
			prompt.WithLexer(&Lexer{}),
			prompt.WithExecuteOnEnterCallback(executeOnEnter),
			prompt.WithPrefix(">> "),
		},
		opts...,
	)

	p := prompt.New(
		executor,
		opts...,
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
