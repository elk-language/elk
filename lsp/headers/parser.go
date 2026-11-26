package headers

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lsp/headers/token"
)

// Error is the list of errors found while parsing headers.
type Error []string

// Error concatenates every message, one per line.
func (e Error) Error() string {
	return strings.Join([]string(e), "; ")
}

// Parsing mode.
type parserMode uint8

const (
	zeroParserMode   parserMode = iota // initial zero value mode
	normalParserMode                   // regular parsing mode
	errorParserMode                    // error parsing mode
)

// Holds the current state of the header parsing process.
type Parser struct {
	source    []byte       // Headers source
	lexer     *Lexer       // lexer which outputs a stream of tokens
	lookahead *token.Token // next token used for predicting productions
	errors    []string
	headerMap *Headers
	mode      parserMode
}

// Instantiate a new header parser.
func NewParser(source []byte) *Parser {
	return &Parser{
		source:    source,
		headerMap: ds.NewOrderedMap[string, string](),
	}
}

// Parse the given source code and return an Abstract Syntax Tree.
// Main entry point to the parser.
func Parse(source []byte) (*Headers, int, error) {
	return NewParser(source).Parse()
}

// Start the parsing process from the top.
func (p *Parser) Parse() (*Headers, int, error) {
	p.reset()

	p.advance() // populate lookahead
	p.headers()

	if p.mode == errorParserMode {
		return nil, 0, Error(p.errors)
	}
	return p.headerMap, p.lexer.cursor, nil
}

func (p *Parser) reset() {
	p.lexer = NewLexer(p.source)
	p.mode = normalParserMode
	p.errors = nil
}

// Adds an error with a custom message.
func (p *Parser) errorMessage(message string) {
	p.mode = errorParserMode
	p.errors = append(p.errors, message)
}

// Adds an error which tells the user that another type of token
// was expected.
func (p *Parser) errorExpected(expected string) {
	p.errorMessage(fmt.Sprintf("unexpected %s, expected %s", p.lookahead.Type.String(), expected))
}

// Add the content of an error token to the syntax error list.
func (p *Parser) errorToken(err *token.Token) {
	p.errorMessage(err.Value)
}

// Attempt to consume the specified token type.
// If the next token doesn't match an error is added.
func (p *Parser) consume(tokenType token.Type) (*token.Token, bool) {
	return p.consumeExpected(tokenType, tokenType.String())
}

// Same as [consume] but lets you specify a custom expected error message.
func (p *Parser) consumeExpected(tokenType token.Type, expected string) (*token.Token, bool) {
	if p.lookahead.Type == token.ERROR {
		return p.advance(), false
	}

	if p.lookahead.Type != tokenType {
		p.errorExpected(expected)
		return p.advance(), false
	}

	return p.advance(), true
}

// Move over to the next token.
func (p *Parser) advance() *token.Token {
	previous := p.lookahead
	if previous != nil && previous.Type == token.ERROR {
		p.errorToken(previous)
	}

	p.lookahead = p.lexer.Next()
	return previous
}

// ==== Productions ====

func (p *Parser) headers() {
	for {
		switch p.lookahead.Type {
		case token.END_OF_FILE, token.NEWLINE:
			return
		case token.STRING:
			p.header()
		default:
			p.errorToken(p.lookahead)
			p.advance()
		}

		if p.mode == errorParserMode {
			return
		}
	}
}

func (p *Parser) header() {
	name := p.advance()
	_, ok := p.consume(token.COLON)
	if !ok {
		return
	}

	value, ok := p.consume(token.STRING)
	if !ok {
		return
	}

	p.headerMap.Set(name.Value, value.Value)
}
