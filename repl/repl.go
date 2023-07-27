package repl

import (
	"strings"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/go-prompt"
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

func executeOnEnter(input string, indentSize int) (indent int, execute bool) {
	p := parser.New("(eval)", []byte(input))
	p.Parse()

	var spaces int
	lastNewline := strings.LastIndexByte(input, '\n')
	for i := lastNewline + 1; i < len(input); i++ {
		b := input[i]
		if b != ' ' {
			break
		}

		spaces++
	}
	baseIndent := spaces / indentSize
	if len(input) > 3 && baseIndent > 0 && input[len(input)-3:] == "end" {
		baseIndent--
	}

	if p.ShouldIndent() {
		return baseIndent + 1, false
	}
	if p.IsIncomplete() {
		return baseIndent, false
	}

	return 0, true
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
