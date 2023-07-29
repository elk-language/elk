package repl

import (
	"strings"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
)

type Lexer struct {
	lexer.Lexer
}

func (l *Lexer) Init(input string) {
	l.Lexer = *lexer.New([]byte(input))
}

func (l *Lexer) Next() (prompt.Token, bool) {
	t := l.Lexer.Next()
	if t.Type == token.END_OF_FILE {
		return nil, false
	}

	return t, true
}

func executeOnEnter(pr *prompt.Prompt, indentSize int) (indent int, execute bool) {
	doc := pr.Buffer().Document()
	if doc.OnLastLine() {
		input := doc.Text
		p := parser.New("(eval)", []byte(input))
		p.Parse()

		prevIndent := doc.PreviousLineIndentSpaces()
		baseIndent := doc.LastLineIndentSpaces()
		if len(input) >= 3 && input[len(input)-3:] == "end" {
			if baseIndent >= prevIndent {
				var indentDiff int
				if baseIndent != prevIndent {
					indentDiff = baseIndent - prevIndent + indentSize
				} else {
					indentDiff = indentSize
				}
				pr.CursorLeftRunes(pstrings.RuneNumber(indentDiff + 3))
				pr.InsertTextMoveCursor("end", true)
				pr.DeleteRunes(pstrings.RuneNumber(indentDiff))
				baseIndent -= indentSize
			} else if prevIndent > baseIndent {
				indentDiff := prevIndent - baseIndent - indentSize
				if indentDiff < 0 {
					indentDiff = 0
				}
				pr.CursorLeftRunes(3)
				pr.InsertTextMoveCursor(strings.Repeat(" ", indentDiff), false)
				pr.CursorRightRunes(3)
				baseIndent = prevIndent - 1
			}
		}

		if p.ShouldIndent() {
			return baseIndent/indentSize + 1, false
		}
		if p.IsIncomplete() {
			return baseIndent / indentSize, false
		}

		return 0, true
	}

	input := pr.Buffer().Document().TextBeforeCursor()
	p := parser.New("(eval)", []byte(input))
	p.Parse()

	baseIndent := pr.Buffer().Document().PreviousLineIndentLevel(indentSize)
	if len(input) > 3 && baseIndent > 0 && input[len(input)-3:] == "end" {
		pr.CursorLeftRunes(pstrings.RuneNumber(indentSize + 3))
		pr.InsertTextMoveCursor("end", true)
		pr.DeleteRunes(pstrings.RuneNumber(indentSize))
		baseIndent--
	}

	if p.ShouldIndent() {
		return baseIndent + 1, false
	}

	return baseIndent, false
}

func Run() {
	p := prompt.New(
		executor,
		prompt.WithLexer(&Lexer{}),
		prompt.WithExecuteOnEnterCallback(executeOnEnter),
	)
	p.Run()
}

func executor(input string) {
}
