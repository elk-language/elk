// Package parser implements a regex parser.
package parser

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/regex/flag"
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
	errors        diagnostic.DiagnosticList
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
func Parse(source string) (ast.Node, diagnostic.DiagnosticList) {
	return New(source).Parse()
}

// Start the parsing process from the top.
func (p *Parser) Parse() (ast.Node, diagnostic.DiagnosticList) {
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
	p.errors.AddFailure(
		message,
		position.NewLocationWithSpan("regex", span),
	)
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
func (p *Parser) union(stopTokens ...token.Type) ast.Node {
	stopTokens = append(stopTokens, token.PIPE)
	left := p.concatenation(stopTokens...)

	for {
		_, ok := p.matchOk(token.PIPE)
		if !ok {
			break
		}

		right := p.concatenation(stopTokens...)

		left = ast.NewUnionNode(
			left.Span().Join(right.Span()),
			left,
			right,
		)
	}

	return left
}

// concatenation = quantifier*
func (p *Parser) concatenation(stopTokens ...token.Type) ast.Node {
	var list []ast.ConcatenationElementNode

	for {
		if p.lookahead.Type == token.END_OF_FILE || slices.Contains(stopTokens, p.lookahead.Type) {
			if len(list) == 1 {
				return list[0]
			}
			return ast.NewConcatenationNode(
				position.JoinSpanOfCollection(list),
				list,
			)
		}
		element := p.quantifier()
		list = append(list, element)
	}
}

func (p *Parser) consumeDigits(stopTokens ...token.Type) (string, *position.Span) {
	var buffer strings.Builder
	var span *position.Span
	for {
		if p.accept(token.END_OF_FILE) || p.accept(stopTokens...) {
			break
		}

		tok, ok := p.consumeExpected(token.CHAR, "a decimal digit")
		if !ok {
			continue
		}
		if !charIsDigit(tok.Char()) {
			p.errorMessageSpan(
				fmt.Sprintf("unexpected %s, expected a decimal digit", tok.StringValue()),
				tok.Span(),
			)
		}

		span = span.Join(tok.Span())
		buffer.WriteString(tok.StringValue())
	}

	return buffer.String(), span
}

func (p *Parser) consumeLetters(stopTokens ...token.Type) (string, *position.Span) {
	var buffer strings.Builder
	var span *position.Span
	for {
		if p.accept(token.END_OF_FILE) || p.accept(stopTokens...) {
			break
		}

		tok, ok := p.consumeExpected(token.CHAR, "an ASCII letter")
		if !ok {
			continue
		}
		if !charIsAsciiLetter(tok.Char()) {
			p.errorMessageSpan(
				fmt.Sprintf("unexpected %s, expected an ASCII letter", tok.StringValue()),
				tok.Span(),
			)
		}

		span = span.Join(tok.Span())
		buffer.WriteString(tok.StringValue())
	}

	return buffer.String(), span
}

func (p *Parser) consumeFlags(stopTokens ...token.Type) (setFlags bitfield.BitField8, unsetFlags bitfield.BitField8, span *position.Span) {
	var disable bool

	for {
		if p.accept(token.END_OF_FILE) || p.accept(stopTokens...) {
			break
		}

		var tok *token.Token
		var ok bool
		if p.match(token.DASH) {
			disable = true
			continue
		} else {
			tok, ok = p.consumeExpected(token.CHAR, "a regex flag")
			if !ok {
				continue
			}
		}

		var currentFlag bitfield.BitFlag8
		switch tok.Char() {
		case 'm':
			currentFlag = flag.MultilineFlag
		case 'i':
			currentFlag = flag.CaseInsensitiveFlag
		case 's':
			currentFlag = flag.DotAllFlag
		case 'U':
			currentFlag = flag.UngreedyFlag
		case 'x':
			currentFlag = flag.ExtendedFlag
		case 'a':
			currentFlag = flag.ASCIIFlag
		default:
			p.errorMessageSpan(
				fmt.Sprintf("unexpected %s, expected a regex flag", tok.StringValue()),
				tok.Span(),
			)
		}
		if disable {
			unsetFlags.SetFlag(currentFlag)
		} else {
			setFlags.SetFlag(currentFlag)
		}

		span = span.Join(tok.Span())
	}

	return setFlags, unsetFlags, span
}

// quantifier = primaryRegex ["+" | "+?" | "*" | "*?" | "?" | "??" | "{" DIGIT+ ["," DIGIT*] "}"]
func (p *Parser) quantifier() ast.ConcatenationElementNode {
	r := p.primaryRegex()
	switch p.lookahead.Type {
	case token.PLUS:
		lastTok := p.advance()
		var alt bool
		if q, ok := p.matchOk(token.QUESTION); ok {
			alt = true
			lastTok = q
		}
		return ast.NewOneOrMoreQuantifierNode(
			r.Span().Join(lastTok.Span()),
			r,
			alt,
		)
	case token.STAR:
		lastTok := p.advance()
		var alt bool
		if q, ok := p.matchOk(token.QUESTION); ok {
			alt = true
			lastTok = q
		}
		return ast.NewZeroOrMoreQuantifierNode(
			r.Span().Join(lastTok.Span()),
			r,
			alt,
		)
	case token.QUESTION:
		lastTok := p.advance()
		var alt bool
		if q, ok := p.matchOk(token.QUESTION); ok {
			alt = true
			lastTok = q
		}
		return ast.NewZeroOrOneQuantifierNode(
			r.Span().Join(lastTok.Span()),
			r,
			alt,
		)
	case token.LBRACE:
		lbrace := p.advance()
		var max, min string
		var commaPresent bool
		var lastSpan *position.Span

		if p.match(token.COMMA) {
			commaPresent = true
			max, lastSpan = p.consumeDigits(token.RBRACE)
		} else {
			min, lastSpan = p.consumeDigits(token.RBRACE, token.COMMA)
			if comma, ok := p.matchOk(token.COMMA); ok {
				commaPresent = true
				lastSpan = comma.Span()
				if !p.accept(token.RBRACE) {
					max, lastSpan = p.consumeDigits(token.RBRACE)
				}
			}
		}
		rbrace, ok := p.consume(token.RBRACE)
		if ok {
			lastSpan = rbrace.Span()
		}
		var alt bool
		if q, ok := p.matchOk(token.QUESTION); ok {
			lastSpan = q.Span()
			alt = true
		}
		if commaPresent {
			return ast.NewNMQuantifierNode(
				r.Span().Join(lastSpan),
				r,
				min,
				max,
				alt,
			)
		}
		if len(min) == 0 {
			if lastSpan == nil {
				lastSpan = lbrace.Span()
			}
			p.errorMessageSpan("expected decimal digits", lastSpan)
		}
		return ast.NewNQuantifierNode(
			r.Span().Join(lastSpan),
			r,
			min,
			alt,
		)
	default:
		return r
	}
}

// primaryRegex = char | escapes | anchors | group
func (p *Parser) primaryRegex() ast.PrimaryRegexNode {
	switch p.lookahead.Type {
	case token.CHAR, token.COMMA, token.RBRACE, token.RBRACKET,
		token.DASH, token.COLON, token.LANGLE, token.RANGLE, token.SINGLE_QUOTE:
		return p.char()
	case token.META_CHAR_ESCAPE:
		tok := p.advance()
		return ast.NewMetaCharEscapeNode(
			tok.Span(),
			tok.Char(),
		)
	case token.QUOTED_TEXT:
		tok := p.advance()
		return ast.NewQuotedTextNode(
			tok.Span(),
			tok.Value,
		)
	case token.LBRACKET:
		return p.charClass()
	case token.LPAREN:
		return p.group()
	case token.BELL_ESCAPE:
		return p.bellEscape()
	case token.FORM_FEED_ESCAPE:
		return p.formFeedEscape()
	case token.TAB_ESCAPE:
		return p.tabEscape()
	case token.NEWLINE_ESCAPE:
		return p.newlineEscape()
	case token.CARRIAGE_RETURN_ESCAPE:
		return p.carriageReturnEscape()
	case token.ABSOLUTE_START_OF_STRING_ANCHOR:
		tok := p.advance()
		return ast.NewAbsoluteStartOfStringAnchorNode(tok.Span())
	case token.ABSOLUTE_END_OF_STRING_ANCHOR:
		tok := p.advance()
		return ast.NewAbsoluteEndOfStringAnchorNode(tok.Span())
	case token.CARET:
		tok := p.advance()
		return ast.NewStartOfStringAnchorNode(tok.Span())
	case token.DOLLAR:
		tok := p.advance()
		return ast.NewEndOfStringAnchorNode(tok.Span())
	case token.DOT:
		tok := p.advance()
		return ast.NewAnyCharClassNode(tok.Span())
	case token.WORD_BOUNDARY_ANCHOR:
		tok := p.advance()
		return ast.NewWordBoundaryAnchorNode(tok.Span())
	case token.NOT_WORD_BOUNDARY_ANCHOR:
		tok := p.advance()
		return ast.NewNotWordBoundaryAnchorNode(tok.Span())
	case token.WORD_CHAR_CLASS:
		tok := p.advance()
		return ast.NewWordCharClassNode(tok.Span())
	case token.NOT_WORD_CHAR_CLASS:
		tok := p.advance()
		return ast.NewNotWordCharClassNode(tok.Span())
	case token.DIGIT_CHAR_CLASS:
		tok := p.advance()
		return ast.NewDigitCharClassNode(tok.Span())
	case token.NOT_DIGIT_CHAR_CLASS:
		tok := p.advance()
		return ast.NewNotDigitCharClassNode(tok.Span())
	case token.WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewWhitespaceCharClassNode(tok.Span())
	case token.NOT_WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewNotWhitespaceCharClassNode(tok.Span())
	case token.HORIZONTAL_WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewHorizontalWhitespaceCharClassNode(tok.Span())
	case token.NOT_HORIZONTAL_WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewNotHorizontalWhitespaceCharClassNode(tok.Span())
	case token.VERTICAL_WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewVerticalWhitespaceCharClassNode(tok.Span())
	case token.NOT_VERTICAL_WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewNotVerticalWhitespaceCharClassNode(tok.Span())
	case token.CARET_ESCAPE:
		return p.caretEscape()
	case token.LONG_UNICODE_ESCAPE:
		return p.longUnicodeEscape()
	case token.UNICODE_ESCAPE:
		return p.unicodeEscape()
	case token.HEX_ESCAPE:
		return p.hexEscape()
	case token.OCTAL_ESCAPE:
		return p.octalEscape()
	case token.SIMPLE_OCTAL_ESCAPE:
		return p.simpleOctalEscape()
	case token.UNICODE_CHAR_CLASS:
		return p.unicodeCharClass()
	case token.NEGATED_UNICODE_CHAR_CLASS:
		return p.negatedUnicodeCharClass()
	}
	if p.lookahead.Type != token.ERROR {
		p.errorExpected("a primary expression")
	}
	t := p.advance()
	return ast.NewInvalidNode(t.Span(), t)
}

func (p *Parser) bellEscape() *ast.BellEscapeNode {
	tok := p.advance()
	return ast.NewBellEscapeNode(tok.Span())
}

func (p *Parser) formFeedEscape() *ast.FormFeedEscapeNode {
	tok := p.advance()
	return ast.NewFormFeedEscapeNode(tok.Span())
}

func (p *Parser) tabEscape() *ast.TabEscapeNode {
	tok := p.advance()
	return ast.NewTabEscapeNode(tok.Span())
}

func (p *Parser) newlineEscape() *ast.NewlineEscapeNode {
	tok := p.advance()
	return ast.NewNewlineEscapeNode(tok.Span())
}

func (p *Parser) carriageReturnEscape() *ast.CarriageReturnEscapeNode {
	tok := p.advance()
	return ast.NewCarriageReturnEscapeNode(tok.Span())
}

func (p *Parser) simpleOctalEscape() *ast.OctalEscapeNode {
	tok := p.advance()
	return ast.NewOctalEscapeNode(tok.Span(), tok.Value)
}

// charClass = "[" ["^"] namedCharClass* "]"
func (p *Parser) charClass() ast.PrimaryRegexNode {
	var list []ast.CharClassElementNode
	var negated bool
	lbracket := p.advance()

	if p.match(token.CARET) {
		negated = true
	}

	for {
		if p.accept(token.END_OF_FILE) {
			p.errorMessage("unterminated character class, missing ]")
			return ast.NewCharClassNode(
				position.JoinSpanOfLastElement(lbracket.Span(), list),
				list,
				negated,
			)
		}
		if p.accept(token.RBRACKET) {
			rbracket := p.advance()
			return ast.NewCharClassNode(
				lbracket.Span().Join(rbracket.Span()),
				list,
				negated,
			)
		}
		element := p.namedCharClass()
		list = append(list, element)
	}
}

// namedCharClass = "[" ":"  ["^"] LETTER* ":" "]" | charRange
func (p *Parser) namedCharClass() ast.CharClassElementNode {
	if !p.accept(token.LBRACKET) {
		return p.charRange()
	}

	var negated bool

	lbracket := p.advance()
	if t, ok := p.consume(token.COLON); !ok {
		return ast.NewInvalidNode(t.Span(), t)
	}
	if p.match(token.CARET) {
		negated = true
	}

	var nameBuffer strings.Builder
	for {
		if p.accept(token.END_OF_FILE) {
			p.errorMessage("unterminated named char class")
			t := p.advance()
			return ast.NewInvalidNode(t.Span(), t)
		}
		if p.match(token.COLON) {
			break
		}

		chTok, ok := p.consume(token.CHAR)
		if !ok {
			continue
		}
		char := chTok.Char()
		if !charIsAsciiLetter(char) {
			p.errorMessageSpan(fmt.Sprintf("unexpected %s, expected an ASCII letter", chTok.StringValue()), chTok.Span())
		}
		nameBuffer.WriteRune(char)
	}

	rbracket, ok := p.consume(token.RBRACKET)
	if !ok {
		return ast.NewInvalidNode(rbracket.Span(), rbracket)
	}

	if nameBuffer.Len() == 0 {
		p.errorMessageSpan(
			"empty named char class",
			lbracket.Span().Join(rbracket.Span()),
		)
	}

	return ast.NewNamedCharClassNode(
		lbracket.Span().Join(rbracket.Span()),
		nameBuffer.String(),
		negated,
	)
}

// charRange = primaryCharClassElement "-" primaryCharClassElement | primaryCharClassElement
func (p *Parser) charRange() ast.CharClassElementNode {
	left := p.primaryCharClassElement()
	if !ast.IsValidCharRangeElement(left) {
		return left
	}

	if !p.match(token.DASH) {
		return left
	}

	right := p.primaryCharClassElement()
	return ast.NewCharRangeNode(
		left.Span().Join(right.Span()),
		left,
		right,
	)
}

// primaryCharClassElement = char | escape
func (p *Parser) primaryCharClassElement() ast.CharClassElementNode {
	switch p.lookahead.Type {
	case token.CHAR, token.COMMA, token.LBRACE, token.RBRACE, token.RBRACKET,
		token.COLON, token.LANGLE, token.RANGLE, token.SINGLE_QUOTE,
		token.CARET, token.DOLLAR, token.DOT, token.PLUS, token.STAR,
		token.QUESTION, token.PIPE, token.LPAREN, token.RPAREN:
		return p.char()
	case token.META_CHAR_ESCAPE:
		tok := p.advance()
		return ast.NewMetaCharEscapeNode(
			tok.Span(),
			tok.Char(),
		)
	case token.BELL_ESCAPE:
		return p.bellEscape()
	case token.FORM_FEED_ESCAPE:
		return p.formFeedEscape()
	case token.TAB_ESCAPE:
		return p.tabEscape()
	case token.NEWLINE_ESCAPE:
		return p.newlineEscape()
	case token.CARRIAGE_RETURN_ESCAPE:
		return p.carriageReturnEscape()
	case token.WORD_CHAR_CLASS:
		tok := p.advance()
		return ast.NewWordCharClassNode(tok.Span())
	case token.NOT_WORD_CHAR_CLASS:
		tok := p.advance()
		return ast.NewNotWordCharClassNode(tok.Span())
	case token.DIGIT_CHAR_CLASS:
		tok := p.advance()
		return ast.NewDigitCharClassNode(tok.Span())
	case token.NOT_DIGIT_CHAR_CLASS:
		tok := p.advance()
		return ast.NewNotDigitCharClassNode(tok.Span())
	case token.WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewWhitespaceCharClassNode(tok.Span())
	case token.NOT_WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewNotWhitespaceCharClassNode(tok.Span())
	case token.HORIZONTAL_WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewHorizontalWhitespaceCharClassNode(tok.Span())
	case token.NOT_HORIZONTAL_WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewNotHorizontalWhitespaceCharClassNode(tok.Span())
	case token.VERTICAL_WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewVerticalWhitespaceCharClassNode(tok.Span())
	case token.NOT_VERTICAL_WHITESPACE_CHAR_CLASS:
		tok := p.advance()
		return ast.NewNotVerticalWhitespaceCharClassNode(tok.Span())
	case token.LONG_UNICODE_ESCAPE:
		return p.longUnicodeEscape()
	case token.UNICODE_ESCAPE:
		return p.unicodeEscape()
	case token.HEX_ESCAPE:
		return p.hexEscape()
	case token.OCTAL_ESCAPE:
		return p.octalEscape()
	case token.SIMPLE_OCTAL_ESCAPE:
		return p.simpleOctalEscape()
	case token.UNICODE_CHAR_CLASS:
		return p.unicodeCharClass()
	case token.NEGATED_UNICODE_CHAR_CLASS:
		return p.negatedUnicodeCharClass()
	}
	if p.lookahead.Type != token.ERROR {
		p.errorExpected("a char class element")
	}
	t := p.advance()
	return ast.NewInvalidNode(t.Span(), t)
}

// group = "(" ["?" (":" | (["P"] "<" ALPHA_CHAR* ">")] union ")"
func (p *Parser) group() ast.PrimaryRegexNode {
	lparen := p.advance()
	var nonCapturing, onlyFlags bool
	var name string
	var setFlags, unsetFlags bitfield.BitField8
	var lastSpan *position.Span
	var content ast.Node

	if _, ok := p.matchOk(token.QUESTION); ok {
		if _, ok := p.matchOk(token.COLON); ok {
			nonCapturing = true
		} else if p.match(token.LANGLE) {
			name, _ = p.consumeLetters(token.RANGLE, token.RPAREN)
			rangle, _ := p.consume(token.RANGLE)
			if len(name) == 0 {
				p.errorMessageSpan("expected a group name", rangle.Span())
			}
		} else if p.match(token.SINGLE_QUOTE) {
			name, _ = p.consumeLetters(token.SINGLE_QUOTE, token.RPAREN)
			apostrophe, _ := p.consume(token.SINGLE_QUOTE)
			if len(name) == 0 {
				p.errorMessageSpan("expected a group name", apostrophe.Span())
			}
		} else if p.lookahead.Type == token.CHAR && p.lookahead.Value == "P" {
			p.advance()
			p.consume(token.LANGLE)
			name, _ = p.consumeLetters(token.RANGLE, token.RPAREN)
			rangle, _ := p.consume(token.RANGLE)
			if len(name) == 0 {
				p.errorMessageSpan("expected a group name", rangle.Span())
			}
		} else {
			setFlags, unsetFlags, _ = p.consumeFlags(token.RPAREN, token.COLON)
			if !p.match(token.COLON) {
				onlyFlags = true
			}
		}
	}
	if !onlyFlags {
		content = p.union(token.RPAREN)
	}
	rparen, ok := p.consume(token.RPAREN)
	if !ok {
		return ast.NewInvalidNode(rparen.Span(), rparen)
	}
	lastSpan = rparen.Span()

	return ast.NewGroupNode(
		lparen.Span().Join(lastSpan),
		content,
		name,
		setFlags,
		unsetFlags,
		nonCapturing,
	)
}

const hexLiteralChars = "0123456789abcdefABCDEF"

func charIsHex(char rune) bool {
	return strings.ContainsRune(hexLiteralChars, char)
}

func charIsOctal(char rune) bool {
	return char >= '0' && char <= '7'
}

func charIsDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func charIsAsciiLetter(char rune) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'
}

const flags = "imsUxa-"

func charIsFlag(char rune) bool {
	return strings.ContainsRune(flags, char)
}

// unicodeCharClass = "\p" ((CHAR) | "{" ["^"] CHAR+ "}")
func (p *Parser) unicodeCharClass() *ast.UnicodeCharClassNode {
	begTok := p.advance()
	span := begTok.Span()
	var buffer strings.Builder
	var negated bool

	if p.match(token.LBRACE) {
		if p.match(token.CARET) {
			negated = true
		}
		for {
			if p.match(token.RBRACE) {
				break
			}
			if p.accept(token.END_OF_FILE) {
				p.errorExpected("an alphabetic character")
				p.advance()
				break
			}
			chTok, ok := p.consumeExpected(token.CHAR, "an alphabetic character")
			if !ok {
				continue
			}
			char := chTok.Char()
			if !unicode.IsLetter(char) {
				p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected an alphabetic character", char), chTok.Span())
			}
			buffer.WriteRune(char)
			span = span.Join(chTok.Span())
		}

		if buffer.Len() == 0 {
			p.errorExpected("a hex digit")
		}
	} else {
		chTok, _ := p.consumeExpected(token.CHAR, "an alphabetic character")
		char := chTok.Char()
		if !unicode.IsLetter(char) {
			p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected an alphabetic character", char), chTok.Span())
		}
		buffer.WriteRune(char)
		span = span.Join(chTok.Span())
	}

	return ast.NewUnicodeCharClassNode(
		span,
		buffer.String(),
		negated,
	)
}

// negatedUnicodeCharClass = "\P" ((CHAR) | "{" ["^"] CHAR+ "}")
func (p *Parser) negatedUnicodeCharClass() *ast.UnicodeCharClassNode {
	begTok := p.advance()
	span := begTok.Span()
	var buffer strings.Builder
	var negated bool

	if p.match(token.LBRACE) {
		if p.match(token.CARET) {
			negated = true
		}
		for {
			if p.match(token.RBRACE) {
				break
			}
			if p.accept(token.END_OF_FILE) {
				p.errorExpected("an alphabetic character")
				p.advance()
				break
			}
			chTok, ok := p.consumeExpected(token.CHAR, "an alphabetic character")
			if !ok {
				continue
			}
			char := chTok.Char()
			if !unicode.IsLetter(char) {
				p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected an alphabetic character", char), chTok.Span())
			}
			buffer.WriteRune(char)
			span = span.Join(chTok.Span())
		}

		if buffer.Len() == 0 {
			p.errorExpected("a hex digit")
		}
	} else {
		chTok, _ := p.consumeExpected(token.CHAR, "an alphabetic character")
		char := chTok.Char()
		if !unicode.IsLetter(char) {
			p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected an alphabetic character", char), chTok.Span())
		}
		buffer.WriteRune(char)
		span = span.Join(chTok.Span())
	}

	return ast.NewUnicodeCharClassNode(
		span,
		buffer.String(),
		!negated,
	)
}

// octalEscape = "\o" ((OCTAL_CHAR OCTAL_CHAR) | "{" OCTAL_CHAR+ "}")
func (p *Parser) octalEscape() *ast.OctalEscapeNode {
	begTok := p.advance()
	span := begTok.Span()
	var buffer strings.Builder

	if p.match(token.LBRACE) {
		for {
			if p.match(token.RBRACE) {
				break
			}
			if p.accept(token.END_OF_FILE) {
				p.errorExpected("an octal digit")
				p.advance()
				break
			}
			chTok, ok := p.consumeExpected(token.CHAR, "an octal digit")
			if !ok {
				continue
			}
			char := chTok.Char()
			if !charIsOctal(char) {
				p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected an octal digit", char), chTok.Span())
			}
			buffer.WriteRune(char)
			span = span.Join(chTok.Span())
		}

		if buffer.Len() == 0 {
			p.errorExpected("an octal digit")
		}
		if buffer.Len() > 3 {
			p.errorMessageSpan("too many octal digits", span)
		}
	} else {
		for range 3 {
			chTok, ok := p.consumeExpected(token.CHAR, "an octal digit")
			if !ok {
				continue
			}
			char := chTok.Char()
			if !charIsOctal(char) {
				p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected an octal digit", char), chTok.Span())
			}
			buffer.WriteRune(char)
			span = span.Join(chTok.Span())
		}
	}

	return ast.NewOctalEscapeNode(
		span,
		buffer.String(),
	)
}

// caretEscape = "\c" ASCII_CHAR
func (p *Parser) caretEscape() *ast.CaretEscapeNode {
	begTok := p.advance()
	span := begTok.Span()

	chTok, _ := p.consumeExpected(token.CHAR, "an ASCII character")
	char := chTok.Char()
	if !charIsAsciiLetter(char) {
		p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected an ASCII letter", char), chTok.Span())
	}
	span = span.Join(chTok.Span())

	return ast.NewCaretEscapeNode(
		span,
		char,
	)
}

// unicodeEscape = "\U" ((HEX_CHAR{8}) | "{" HEX_CHAR+ "}")
func (p *Parser) longUnicodeEscape() *ast.UnicodeEscapeNode {
	begTok := p.advance()
	span := begTok.Span()
	var buffer strings.Builder

	if p.match(token.LBRACE) {
		for {
			if p.match(token.RBRACE) {
				break
			}
			if p.accept(token.END_OF_FILE) {
				p.errorExpected("a hex digit")
				p.advance()
				break
			}
			chTok, ok := p.consumeExpected(token.CHAR, "a hex digit")
			if !ok {
				continue
			}
			char := chTok.Char()
			if !charIsHex(char) {
				p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected a hex digit", char), chTok.Span())
			}
			buffer.WriteRune(char)
			span = span.Join(chTok.Span())
		}

		if buffer.Len() == 0 {
			p.errorExpected("a hex digit")
		}
	} else {
		for range 8 {
			chTok, ok := p.consumeExpected(token.CHAR, "a hex digit")
			if !ok {
				continue
			}
			char := chTok.Char()
			if !charIsHex(char) {
				p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected a hex digit", char), chTok.Span())
			}
			buffer.WriteRune(char)
			span = span.Join(chTok.Span())
		}
	}

	return ast.NewUnicodeEscapeNode(
		span,
		buffer.String(),
	)
}

// unicodeEscape = "\u" ((HEX_CHAR{4}) | "{" HEX_CHAR+ "}")
func (p *Parser) unicodeEscape() *ast.UnicodeEscapeNode {
	begTok := p.advance()
	span := begTok.Span()
	var buffer strings.Builder

	if p.match(token.LBRACE) {
		for {
			if p.match(token.RBRACE) {
				break
			}
			if p.accept(token.END_OF_FILE) {
				p.errorExpected("a hex digit")
				p.advance()
				break
			}
			chTok, ok := p.consumeExpected(token.CHAR, "a hex digit")
			if !ok {
				continue
			}
			char := chTok.Char()
			if !charIsHex(char) {
				p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected a hex digit", char), chTok.Span())
			}
			buffer.WriteRune(char)
			span = span.Join(chTok.Span())
		}

		if buffer.Len() == 0 {
			p.errorExpected("a hex digit")
		}
	} else {
		for range 4 {
			chTok, ok := p.consumeExpected(token.CHAR, "a hex digit")
			if !ok {
				continue
			}
			char := chTok.Char()
			if !charIsHex(char) {
				p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected a hex digit", char), chTok.Span())
			}
			buffer.WriteRune(char)
			span = span.Join(chTok.Span())
		}
	}

	return ast.NewUnicodeEscapeNode(
		span,
		buffer.String(),
	)
}

// hexEscape = "\x" ((HEX_CHAR HEX_CHAR) | "{" HEX_CHAR+ "}")
func (p *Parser) hexEscape() *ast.HexEscapeNode {
	begTok := p.advance()
	span := begTok.Span()
	var buffer strings.Builder

	if p.match(token.LBRACE) {
		for {
			if p.match(token.RBRACE) {
				break
			}
			if p.accept(token.END_OF_FILE) {
				p.errorExpected("a hex digit")
				p.advance()
				break
			}
			chTok, ok := p.consumeExpected(token.CHAR, "a hex digit")
			if !ok {
				continue
			}
			char := chTok.Char()
			if !charIsHex(char) {
				p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected a hex digit", char), chTok.Span())
			}
			buffer.WriteRune(char)
			span = span.Join(chTok.Span())
		}

		if buffer.Len() == 0 {
			p.errorExpected("a hex digit")
		}
	} else {
		for range 2 {
			chTok, ok := p.consumeExpected(token.CHAR, "a hex digit")
			if !ok {
				continue
			}
			char := chTok.Char()
			if !charIsHex(char) {
				p.errorMessageSpan(fmt.Sprintf("unexpected %c, expected a hex digit", char), chTok.Span())
			}
			buffer.WriteRune(char)
			span = span.Join(chTok.Span())
		}
	}

	return ast.NewHexEscapeNode(
		span,
		buffer.String(),
	)
}

func (p *Parser) char() *ast.CharNode {
	charTok := p.advance()
	return ast.NewCharNode(
		charTok.Span(),
		charTok.Char(),
	)
}
