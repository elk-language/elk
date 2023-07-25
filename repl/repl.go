package repl

import (
	"github.com/elk-language/elk/lexer"
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

func Run() {
	p := prompt.New(
		executor,
		prompt.WithLexer(&Lexer{}),
	)
	p.Run()
}

func executor(input string) {
}
