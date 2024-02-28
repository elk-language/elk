// Package parser implements a regex parser.
package parser

import (
	"fmt"
	"slices"
	"unicode/utf8"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/regex/lexer"
	"github.com/elk-language/elk/regex/parser/ast"
	"github.com/elk-language/elk/regex/token"
)

// Parsing mode.
type mode uint8

const (
	zeroMode   mode = iota // initial zero value mode
	normalMode             // regular parsing mode
)

// Holds the current state of the parsing process.
type Parser struct {
	source        string       // Regex source
	lexer         *lexer.Lexer // lexer which outputs a stream of tokens
	lookahead     *token.Token // next token used for predicting productions
	nextLookahead *token.Token // second next token used for predicting productions
	errors        errors.ErrorList
	mode          mode
}

// Instantiate a new parser.
func New(source string) *Parser {
	return &Parser{
		source: source,
		mode:   zeroMode,
	}
}

// Parse the given source code and return an Abstract Syntax Tree.
// Main entry point to the parser.
func Parse(source string) (ast.Node, errors.ErrorList) {
	return New(source).Parse()
}

// Start the parsing process from the top.
func (p *Parser) Parse() (ast.Node, errors.ErrorList) {
	p.reset()

	p.advance() // populate nextLookahead
	p.advance() // populate lookahead
	return p.program(), p.errors
}

func (p *Parser) reset() {
	p.lexer = lexer.New(p.source)
	p.mode = normalMode
	p.errors = nil
}

// Same as [errorExpected] but lets you pass a token type.
func (p *Parser) errorExpectedToken(expected token.Type) {
	p.errorExpected(expected.String())
}

// Adds an error with a custom message.
func (p *Parser) errorMessage(message string) {
	p.errorMessageSpan(message, p.lookahead.Span())
}

// Adds an error which tells the user that the received
// token is unexpected.
func (p *Parser) errorUnexpected(message string) {
	p.errorMessage(fmt.Sprintf("unexpected %s, %s", p.lookahead.Type.String(), message))
}

// Adds an error which tells the user that another type of token
// was expected.
func (p *Parser) errorExpected(expected string) {
	p.errorMessage(fmt.Sprintf("unexpected %s, expected %s", p.lookahead.Type.String(), expected))
}

// Add the content of an error token to the syntax error list.
func (p *Parser) errorToken(err *token.Token) {
	p.errorMessageSpan(err.Value, err.Span())
}

// Same as [errorMessage] but let's you pass a Span.
func (p *Parser) errorMessageSpan(message string, span *position.Span) {
	p.errors.Add(
		message,
		position.NewLocationWithSpan("regex", span),
	)
}

// Attempt to consume the specified token type.
// If the next token doesn't match an error is added and the parser
// enters panic mode.
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

// Checks if the next token matches any of the given types,
// if so it gets consumed.
func (p *Parser) match(types ...token.Type) bool {
	for _, typ := range types {
		if p.accept(typ) {
			p.advance()
			return true
		}
	}

	return false
}

// Same as [match] but returns the consumed token.
func (p *Parser) matchOk(types ...token.Type) (*token.Token, bool) {
	for _, typ := range types {
		if p.accept(typ) {
			return p.advance(), true
		}
	}

	return nil, false
}

// Checks whether there are any more tokens to be consumed.
func (p *Parser) isAtEnd() bool {
	return p.lookahead.Type == token.END_OF_FILE
}

// Checks whether the next token matches any the specified types.
func (p *Parser) accept(tokenTypes ...token.Type) bool {
	for _, typ := range tokenTypes {
		if p.lookahead.Type == typ {
			return true
		}
	}
	return false
}

// Checks whether the second next token matches any the specified types.
func (p *Parser) acceptNext(tokenTypes ...token.Type) bool {
	for _, typ := range tokenTypes {
		if p.nextLookahead.Type == typ {
			return true
		}
	}
	return false
}

// Move over to the next token.
func (p *Parser) advance() *token.Token {
	previous := p.lookahead
	previousNext := p.nextLookahead
	if previousNext != nil && previousNext.Type == token.ERROR {
		p.errorToken(previousNext)
	}
	p.nextLookahead = p.lexer.Next()
	p.lookahead = previousNext
	return previous
}

// ==== Productions ====

// program = union
func (p *Parser) program() ast.Node {
	return p.union()
}

// union = concatenation | union "|" concatenation
func (p *Parser) union() ast.Node {
	left := p.concatenation(token.PIPE)

	for {
		_, ok := p.matchOk(token.PIPE)
		if !ok {
			break
		}

		right := p.concatenation(token.PIPE)

		left = ast.NewUnionNode(
			left.Span().Join(right.Span()),
			left,
			right,
		)
	}

	return left
}

// concatenation = primaryRegex*
func (p *Parser) concatenation(stopTokens ...token.Type) ast.Node {
	var list []ast.PrimaryRegexNode

	for {
		if p.lookahead.Type == token.END_OF_FILE {
			if len(list) == 1 {
				return list[0]
			}
			return ast.NewConcatenationNode(
				position.JoinSpanOfCollection(list),
				list,
			)
		}
		if slices.Contains(stopTokens, p.lookahead.Type) {
			if len(list) == 1 {
				return list[0]
			}
			return ast.NewConcatenationNode(
				position.JoinSpanOfCollection(list),
				list,
			)
		}
		element := p.primaryRegex()
		list = append(list, element)
	}
}

// primaryRegex = char | escapes | anchors
func (p *Parser) primaryRegex() ast.PrimaryRegexNode {
	switch p.lookahead.Type {
	case token.CHAR:
		return p.char()
	case token.BELL_ESCAPE:
		tok := p.advance()
		return ast.NewBellEscapeNode(tok.Span())
	}
	t := p.advance()
	return ast.NewInvalidNode(t.Span(), t)
}

func (p *Parser) char() *ast.CharNode {
	charTok := p.advance()
	char, _ := utf8.DecodeRuneInString(charTok.Value)
	return ast.NewCharNode(
		charTok.Span(),
		char,
	)
}
