package parser

import (
	"github.com/elk-language/elk/lexer"
)

// Holds the current state of the parsing process.
type Parser struct {
	source    []byte       // Elk source code
	lexer     *lexer.Lexer // lexer which outputs a stream of tokens
	lookahead *lexer.Token // next token used for predicting productions
}

// Instantiate a new parser.
func new(source []byte) *Parser {
	return &Parser{
		source: source,
		lexer:  lexer.New(source),
	}
}

// Parse the given source code and return an Abstract Syntax Tree.
// Main entry point to the parser.
// func Parse(source []byte) (*ast.ProgramNode, error) {
// 	return new(source).parse()
// }

// func (*Parser) parse() (*ast.ProgramNode, error) {

// }
