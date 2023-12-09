// Package parser implements a parser
// used by the Elk interpreter.
//
// Parser expects a slice of bytes containing Elk source code
// parses it, registering any encountered errors, and returns an Abstract Syntax Tree.
package parser

import (
	"fmt"
	"unicode/utf8"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
)

// Parsing mode.
type mode uint8

const (
	zeroMode             mode = iota // initial zero value mode
	normalMode                       // regular parsing mode
	panicMode                        // triggered after encountering a syntax error, changes to `normalMode` after synchronisation
	withoutBitwiseOrMode             // disables bitwise OR `|` from the grammar
	incompleteMode                   // the input is incomplete, parser expected more tokens but received an END_OF_FILE.
)

// Holds the current state of the parsing process.
type Parser struct {
	sourceName       string       // Path to the source file or some name.
	source           string       // Elk source code
	lexer            *lexer.Lexer // lexer which outputs a stream of tokens
	lookahead        *token.Token // next token used for predicting productions
	nextLookahead    *token.Token // second next token used for predicting productions
	errors           errors.ErrorList
	mode             mode
	indentedSection  bool
	incompleteIndent bool
}

// Instantiate a new parser.
func New(sourceName string, source string) *Parser {
	return &Parser{
		sourceName: sourceName,
		source:     source,
		mode:       zeroMode,
	}
}

// Parse the given source code and return an Abstract Syntax Tree.
// Main entry point to the parser.
func Parse(sourceName string, source string) (*ast.ProgramNode, errors.ErrorList) {
	return New(sourceName, source).Parse()
}

// Returns true when the parser had finished early
// because of an END_OF_FILE token.
func (p *Parser) IsIncomplete() bool {
	return p.mode == incompleteMode
}

// Returns true when the parser had finished early
// because of an END_OF_FILE token and the following
// code should be indented.
func (p *Parser) ShouldIndent() bool {
	return p.incompleteIndent
}

// Start the parsing process from the top.
func (p *Parser) Parse() (*ast.ProgramNode, errors.ErrorList) {
	p.reset()

	p.advance() // populate nextLookahead
	p.advance() // populate lookahead
	return p.program(), p.errors
}

func (p *Parser) reset() {
	p.lexer = lexer.NewWithName(p.sourceName, p.source)
	p.mode = normalMode
	p.errors = nil
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

// Same as [errorExpected] but lets you pass a token type.
func (p *Parser) errorExpectedToken(expected token.Type) {
	p.errorExpected(expected.String())
}

// Adds an error with a custom message.
func (p *Parser) errorMessage(message string) {
	p.errorMessageSpan(message, p.lookahead.Span())
}

// Same as [errorMessage] but let's you pass a Span.
func (p *Parser) errorMessageSpan(message string, span *position.Span) {
	if p.mode == panicMode {
		return
	}

	p.errors.Add(
		message,
		position.NewLocationWithSpan(p.sourceName, span),
	)
}

// Add the content of an error token to the syntax error list.
func (p *Parser) errorToken(err *token.Token) {
	p.errorMessageSpan(err.Value, err.Span())
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
		p.updateErrorMode(true)
		return p.advance(), false
	}

	return p.advance(), true
}

// Update the mode of the parser after encountering an error.
func (p *Parser) updateErrorMode(panic bool) {
	if p.lookahead.Type == token.END_OF_FILE && p.indentedSection {
		p.incompleteIndent = true
		p.mode = incompleteMode
	} else if p.lookahead.Type == token.END_OF_FILE {
		p.mode = incompleteMode
	} else if panic {
		p.mode = panicMode
	}
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

// Discards tokens until something resembling a new statement is encountered.
// Used for recovering after errors.
// Returns `true` when there is a statement separator to consume.
func (p *Parser) synchronise() bool {
	p.mode = normalMode
	for {
		switch p.lookahead.Type {
		case token.END_OF_FILE:
			return false
		case token.SEMICOLON,
			token.NEWLINE:
			return true
		}

		p.advance()
	}
}

// Accept and ignore any number of consecutive newline tokens.
func (p *Parser) swallowNewlines() {
	for {
		if !p.match(token.NEWLINE) {
			break
		}
	}
}

// Checks if the given slice of token types contains
// the given token type.
func containsToken(slice []token.Type, v token.Type) bool {
	for _, s := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// Consume a block of statements, like in `else` expressions,
// that terminates with `end`.
func (p *Parser) statementBlock(stopTokens ...token.Type) (*position.Span, []ast.StatementNode, bool) {
	var thenBody []ast.StatementNode
	var lastSpan *position.Span
	var multiline bool

	if p.lookahead.Type == token.END_OF_FILE {
		p.errorExpected("a statement separator or an expression")
		p.updateErrorMode(true)
		return p.lookahead.Span(), nil, false
	}

	if !p.lookahead.IsStatementSeparator() {
		expr := p.expressionWithoutModifier()
		thenBody = append(thenBody, ast.NewExpressionStatementNode(
			expr.Span(),
			expr,
		))
		lastSpan = expr.Span()
	} else {
		multiline = true
		p.advance()

		if p.accept(token.END) {
			lastSpan = p.lookahead.Span()
		} else if !containsToken(stopTokens, p.lookahead.Type) {
			thenBody = p.statements(stopTokens...)
			if len(thenBody) > 0 {
				lastSpan = position.SpanOfLastElement(thenBody)
			}
		}
	}

	return lastSpan, thenBody, multiline
}

// statementProduction = subProduction [SEPARATOR]
func statementProduction[Expression, Statement ast.Node](p *Parser, constructor statementConstructor[Expression, Statement], expressionProduction func() Expression, separators ...token.Type) Statement {
	expr := expressionProduction()
	var sep *token.Token
	if p.lookahead.IsStatementSeparator() || p.lookahead.Type == token.END_OF_FILE {
		sep = p.advance()
		return constructor(
			expr.Span().Join(sep.Span()),
			expr,
		)
	}
	for _, sepType := range separators {
		if p.lookahead.Type == sepType {
			return constructor(
				expr.Span(),
				expr,
			)
		}
	}
	if p.match(token.ERROR) {
		if p.synchronise() {
			p.advance()
		}
		return constructor(
			expr.Span(),
			expr,
		)
	}

	p.updateErrorMode(false)
	p.errorExpected(statementSeparatorMessage)
	if p.synchronise() {
		p.advance()
	}

	return constructor(
		expr.Span(),
		expr,
	)
}

type statementsProduction[Statement ast.Node] func(...token.Type) []Statement

// Represents an AST Node constructor function for a new ast.StatementNode
type statementConstructor[Expression, Statement ast.Node] func(*position.Span, Expression) Statement

// Consume a block of statements, like in `if`, `elsif` or `while` expressions,
// that terminates with `end` or can be single-line when it begins with `then`.
func genericStatementBlockWithThen[Expression, Statement ast.Node](
	p *Parser,
	statementsProduction statementsProduction[Statement],
	expressionProduction func() Expression,
	statementConstructor statementConstructor[Expression, Statement],
	stopTokens ...token.Type,
) (*position.Span, []Statement, bool) {
	var thenBody []Statement
	var lastSpan *position.Span
	var multiline bool

	if p.lookahead.Type == token.THEN {
		p.advance()
		expr := expressionProduction()
		thenBody = append(thenBody, statementConstructor(
			expr.Span(),
			expr,
		))
		lastSpan = expr.Span()
	} else {
		multiline = true
		if p.lookahead.IsStatementSeparator() {
			p.advance()
		} else {
			p.updateErrorMode(false)
			p.errorExpected(statementSeparatorMessage)
		}

		if p.accept(token.END) {
			lastSpan = p.lookahead.Span()
		} else if !containsToken(stopTokens, p.lookahead.Type) {
			thenBody = statementsProduction(stopTokens...)
			if len(thenBody) > 0 {
				lastSpan = thenBody[len(thenBody)-1].Span()
			}
		}
	}

	return lastSpan, thenBody, multiline
}

// Consume a block of statements, like in `if`, `elsif` or `while` expressions,
// that terminates with `end` or can be single-line when it begins with `then`.
func (p *Parser) statementBlockWithThen(stopTokens ...token.Type) (*position.Span, []ast.StatementNode, bool) {
	return genericStatementBlockWithThen(p, p.statements, p.expressionWithoutModifier, ast.NewExpressionStatementNodeI, stopTokens...)
}

// Consume a block of statements, like in `if`, `elsif` or `while` expressions,
// that terminates with `end` or can be single-line when it begins with `then`.
func (p *Parser) structBodyStatementBlockWithThen(stopTokens ...token.Type) (*position.Span, []ast.StructBodyStatementNode, bool) {
	return genericStatementBlockWithThen(p, p.structBodyStatements, p.formalParameter, ast.NewParameterStatementNodeI, stopTokens...)
}

// Represents an AST Node constructor function for binary operators
type binaryConstructor[Element ast.Node] func(*position.Span, *token.Token, Element, Element) Element

// binaryProduction = subProduction | binaryProduction operators subProduction
func binaryProduction[Element ast.Node](p *Parser, constructor binaryConstructor[Element], subProduction func() Element, operators ...token.Type) Element {
	left := subProduction()

	for {
		operator, ok := p.matchOk(operators...)
		if !ok {
			break
		}
		p.swallowNewlines()

		p.indentedSection = true
		right := subProduction()
		p.indentedSection = false

		left = constructor(
			left.Span().Join(right.Span()),
			operator,
			left,
			right,
		)
	}

	return left
}

// Represents an AST Node constructor function for an `include`- like expression
// eg. `include`, `extend`, `enhance`
type includelikeConstructor[T ast.Node] func(*position.Span, []ast.ComplexConstantNode) T

// includelikeExpression = keyword genericConstantList
func includelikeExpression[T ast.Node](p *Parser, constructor includelikeConstructor[T]) T {
	keyword := p.advance()
	consts := p.genericConstantList()
	span := position.JoinSpanOfLastElement(keyword.Span(), consts)

	return constructor(
		span,
		consts,
	)
}

// binaryExpression = subProduction | binaryExpression operators subProduction
func (p *Parser) binaryExpression(subProduction func() ast.ExpressionNode, operators ...token.Type) ast.ExpressionNode {
	return binaryProduction(p, ast.NewBinaryExpressionNodeI, subProduction, operators...)
}

// binaryTypeExpression = subProduction | binaryTypeExpression operators subProduction
func (p *Parser) binaryTypeExpression(subProduction func() ast.TypeNode, operators ...token.Type) ast.TypeNode {
	return binaryProduction(p, ast.NewBinaryTypeExpressionNodeI, subProduction, operators...)
}

// logicalExpression = subProduction | logicalExpression operators subProduction
func (p *Parser) logicalExpression(subProduction func() ast.ExpressionNode, operators ...token.Type) ast.ExpressionNode {
	return binaryProduction(p, ast.NewLogicalExpressionNodeI, subProduction, operators...)
}

// commaSeparatedList = element ("," element)*
func commaSeparatedList[Element ast.Node](p *Parser, elementProduction func() Element, stopTokens ...token.Type) []Element {
	var elements []Element
	elements = append(elements, elementProduction())

	for {
		p.swallowNewlines()
		if p.accept(token.END_OF_FILE) {
			break
		}
		for _, stopToken := range stopTokens {
			if p.lookahead.Type == stopToken {
				break
			}
		}
		if !p.match(token.COMMA) {
			break
		}
		p.swallowNewlines()
		elements = append(elements, elementProduction())
	}

	return elements
}

// commaSeparatedListWithoutTerminator = element ("," element)*
func commaSeparatedListWithoutTerminator[Element ast.Node](p *Parser, elementProduction func() Element, stopTokens ...token.Type) []Element {
	var elements []Element
	elements = append(elements, elementProduction())

	for {
		if p.accept(token.END_OF_FILE) {
			break
		}
		for _, stopToken := range stopTokens {
			if p.lookahead.Type == stopToken {
				break
			}
		}
		if !p.match(token.COMMA) {
			break
		}
		p.swallowNewlines()
		elements = append(elements, elementProduction())
	}

	return elements
}

// Consume subProductions until one of the provided token types is encountered.
//
// repeatedProduction = subProduction*
func repeatedProduction[Element ast.Node](p *Parser, subProduction func() Element, stopTokens ...token.Type) []Element {
	var list []Element

	for {
		if p.lookahead.Type == token.END_OF_FILE {
			return list
		}
		for _, stopToken := range stopTokens {
			if p.lookahead.Type == stopToken {
				return list
			}
		}
		element := subProduction()
		list = append(list, element)
	}
}

// A production that can be repeated as in `repeatableProductionWithStop*`
type repeatableProductionWithStop[Element ast.Node] func(...token.Type) Element

// Consume subProductions until one of the provided token types are encountered.
//
// repeatedProductionWithStop = subProduction*
func repeatedProductionWithStop[Element ast.Node](p *Parser, subProduction repeatableProductionWithStop[Element], stopTokens ...token.Type) []Element {
	var list []Element

	for {
		if p.lookahead.Type == token.END_OF_FILE {
			return list
		}
		for _, stopToken := range stopTokens {
			if p.lookahead.Type == stopToken {
				return list
			}
		}
		element := subProduction(stopTokens...)
		list = append(list, element)
	}
}

// ==== Productions ====

// program = statements
func (p *Parser) program() *ast.ProgramNode {
	statements := p.statements()
	lastSpan := position.SpanOfLastElement(statements)
	return ast.NewProgramNode(
		position.NewSpanFromPosition(position.New(0, 1, 1)).Join(lastSpan),
		statements,
	)
}

// statements = statement*
func (p *Parser) statements(stopTokens ...token.Type) []ast.StatementNode {
	return repeatedProductionWithStop(p, p.statement, stopTokens...)
}

// statement = emptyStatement | expressionStatement
func (p *Parser) statement(separators ...token.Type) ast.StatementNode {
	if p.lookahead.IsStatementSeparator() {
		return p.emptyStatement()
	}

	return p.expressionStatement(separators...)
}

// structBodyStatements = structBodyStatement*
func (p *Parser) structBodyStatements(stopTokens ...token.Type) []ast.StructBodyStatementNode {
	return repeatedProductionWithStop(p, p.structBodyStatement, stopTokens...)
}

// structBodyStatement = emptyStatement | parameterStatement
func (p *Parser) structBodyStatement(separators ...token.Type) ast.StructBodyStatementNode {
	if p.lookahead.IsStatementSeparator() {
		return p.emptyStatement()
	}

	return p.parameterStatement(separators...)
}

// parameterStatement = formalParameter [SEPARATOR]
func (p *Parser) parameterStatement(separators ...token.Type) *ast.ParameterStatementNode {
	return statementProduction(p, ast.NewParameterStatementNode, p.formalParameter, separators...)
}

// emptyStatement = SEPARATOR
func (p *Parser) emptyStatement() *ast.EmptyStatementNode {
	sepTok := p.advance()
	return ast.NewEmptyStatementNode(sepTok.Span())
}

const statementSeparatorMessage = "a statement separator `\\n`, `;`"

// expressionStatement = expressionWithModifier [SEPARATOR]
func (p *Parser) expressionStatement(separators ...token.Type) *ast.ExpressionStatementNode {
	return statementProduction(p, ast.NewExpressionStatementNode, p.expressionWithModifier, separators...)
}

// expressionWithModifier = modifierExpression
func (p *Parser) expressionWithModifier() ast.ExpressionNode {
	expr := p.modifierExpression()
	if p.mode == panicMode {
		p.synchronise()
	}
	return expr
}

// expressionWithoutModifier = assignmentExpression
func (p *Parser) expressionWithoutModifier() ast.ExpressionNode {
	asgmt := p.assignmentExpression()
	if p.mode == panicMode {
		p.synchronise()
	}
	return asgmt
}

// modifierExpression = expressionWithoutModifier |
// expressionWithoutModifier ("if" | "unless" | "while" | "until") expressionWithoutModifier |
// expressionWithoutModifier "if" expressionWithoutModifier "else" expressionWithoutModifier |
// expressionWithoutModifier "for" identifierList "in" expressionWithoutModifier
func (p *Parser) modifierExpression() ast.ExpressionNode {
	left := p.expressionWithoutModifier()

	switch p.lookahead.Type {
	case token.UNLESS, token.WHILE, token.UNTIL:
		mod := p.advance()
		p.swallowNewlines()
		right := p.expressionWithoutModifier()
		return ast.NewModifierNode(
			left.Span().Join(right.Span()),
			mod,
			left,
			right,
		)
	case token.IF:
		ifTok := p.advance()
		p.swallowNewlines()
		cond := p.expressionWithoutModifier()
		if p.lookahead.Type == token.ELSE {
			p.advance()
			p.swallowNewlines()
			elseExpr := p.expressionWithoutModifier()
			return ast.NewModifierIfElseNode(
				left.Span().Join(elseExpr.Span()),
				left,
				cond,
				elseExpr,
			)
		}
		return ast.NewModifierNode(
			left.Span().Join(cond.Span()),
			ifTok,
			left,
			cond,
		)
	case token.FOR:
		p.advance()
		p.swallowNewlines()
		params := p.identifierList(token.IN)
		inTok, ok := p.consume(token.IN)
		if !ok {
			return ast.NewInvalidNode(inTok.Span(), inTok)
		}
		p.swallowNewlines()
		inExpr := p.expressionWithoutModifier()
		return ast.NewModifierForInNode(
			left.Span().Join(inExpr.Span()),
			left,
			params,
			inExpr,
		)
	}

	return left
}

// assignmentExpression = logicalOrExpression | expression ASSIGN_OP assignmentExpression
func (p *Parser) assignmentExpression() ast.ExpressionNode {
	left := p.logicalOrExpression()

	if !p.lookahead.IsAssignmentOperator() {
		return left
	}

	if p.lookahead.Type == token.COLON_EQUAL {
		if !ast.IsValidDeclarationTarget(left) {
			p.errorMessageSpan(
				fmt.Sprintf("invalid `%s` declaration target", p.lookahead.Type.String()),
				left.Span(),
			)
		}
	} else if ast.IsConstant(left) {
		p.errorMessageSpan(
			"constants can't be assigned, maybe you meant to declare it with `:=`",
			left.Span(),
		)
	} else if !ast.IsValidAssignmentTarget(left) {
		p.errorMessageSpan(
			fmt.Sprintf("invalid `%s` assignment target", p.lookahead.Type.String()),
			left.Span(),
		)
	}

	operator := p.advance()
	p.swallowNewlines()

	p.indentedSection = true
	right := p.assignmentExpression()
	p.indentedSection = false

	return ast.NewAssignmentExpressionNode(
		left.Span().Join(right.Span()),
		operator,
		left,
		right,
	)
}

// formalParameter = identifier [":" typeAnnotation] ["=" expressionWithoutModifier]
func (p *Parser) formalParameter() ast.ParameterNode {
	var init ast.ExpressionNode
	var typ ast.TypeNode

	var paramName *token.Token
	var kind ast.ParameterKind
	var span *position.Span

	var restParam bool
	if starTok, ok := p.matchOk(token.STAR); ok {
		kind = ast.PositionalRestParameterKind
		span = starTok.Span()
		restParam = true
	} else if starStarTok, ok := p.matchOk(token.STAR_STAR); ok {
		kind = ast.NamedRestParameterKind
		span = starStarTok.Span()
		restParam = true
	}

	switch p.lookahead.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		paramName = p.advance()
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		p.errorExpected("a lowercase identifier as the name of the declared formalParameter")
		paramName = p.advance()
	default:
		p.errorExpected("an identifier as the name of the declared formalParameter")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Span(),
			tok,
		)
	}
	span = span.Join(paramName.Span())

	if p.match(token.COLON) {
		typ = p.intersectionType()
		span = span.Join(typ.Span())
	}

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		if restParam {
			p.errorMessageSpan("rest parameters can't have default values", init.Span())
		}
		span = span.Join(init.Span())
	}

	return ast.NewFormalParameterNode(
		span,
		paramName.Value,
		typ,
		init,
		kind,
	)
}

// formalParameter = ["*" | "**"] identifier [":" typeAnnotation] ["=" expressionWithoutModifier]
func (p *Parser) methodParameter() ast.ParameterNode {
	var init ast.ExpressionNode
	var typ ast.TypeNode

	var paramName *token.Token
	var setIvar bool
	var kind ast.ParameterKind
	var span *position.Span

	var restParam bool
	if starTok, ok := p.matchOk(token.STAR); ok {
		kind = ast.PositionalRestParameterKind
		span = starTok.Span()
		restParam = true
	} else if starStarTok, ok := p.matchOk(token.STAR_STAR); ok {
		kind = ast.NamedRestParameterKind
		span = starStarTok.Span()
		restParam = true
	}

	switch p.lookahead.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		paramName = p.advance()
	case token.INSTANCE_VARIABLE:
		paramName = p.advance()
		setIvar = true
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		p.errorExpected("a lowercase identifier as the name of the declared formalParameter")
		paramName = p.advance()
	default:
		p.errorExpected("an identifier as the name of the declared formalParameter")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Span(),
			tok,
		)
	}
	span = span.Join(paramName.Span())

	if p.match(token.COLON) {
		typ = p.intersectionType()
		span = span.Join(typ.Span())
	}

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		span = span.Join(init.Span())
		if restParam {
			p.errorMessageSpan("rest parameters can't have default values", span)
		}
	}

	return ast.NewMethodParameterNode(
		span,
		paramName.Value,
		setIvar,
		typ,
		init,
		kind,
	)
}

// parameterList = parameter ("," parameter)*
func (p *Parser) parameterList(parameter func() ast.ParameterNode, stopTokens ...token.Type) []ast.ParameterNode {
	var elements []ast.ParameterNode
	element := parameter()
	optionalSeen := element.IsOptional()
	posRestSeen := ast.IsPositionalRestParam(element)
	namedRestSeen := ast.IsNamedRestParam(element)
	elements = append(elements, element)

	for {
		if p.lookahead.Type == token.END_OF_FILE {
			break
		}
		for _, stopToken := range stopTokens {
			if p.lookahead.Type == stopToken {
				break
			}
		}
		if !p.match(token.COMMA) {
			break
		}
		p.swallowNewlines()
		element := parameter()
		elements = append(elements, element)

		posRest := ast.IsPositionalRestParam(element)
		if posRest && posRestSeen {
			p.errorMessageSpan("there should be only a single positional rest parameter", element.Span())
			continue
		}

		if posRest {
			if namedRestSeen {
				p.errorMessageSpan("named rest parameters should appear last", element.Span())
			}
			posRestSeen = true
			continue
		}

		namedRest := ast.IsNamedRestParam(element)
		if namedRest && namedRestSeen {
			p.errorMessageSpan("there should be only a single named rest parameter", element.Span())
			continue
		}

		if namedRestSeen {
			p.errorMessageSpan("named rest parameters should appear last", element.Span())
			continue
		}

		if namedRest {
			namedRestSeen = true
			continue
		}

		if posRest {
			continue
		}

		opt := element.IsOptional()
		if !opt && optionalSeen {
			p.errorMessageSpan("required parameters can't appear after optional parameters", element.Span())
		} else if opt {
			if !optionalSeen {
				optionalSeen = true
			}
			if posRestSeen {
				p.errorMessageSpan("optional parameters can't appear after rest parameters", element.Span())
			}
		}
	}

	return elements
}

// formalParameterList = formalParameter ("," formalParameter)*
func (p *Parser) formalParameterList(stopTokens ...token.Type) []ast.ParameterNode {
	return p.parameterList(p.formalParameter, stopTokens...)
}

// methodParameterList = methodParameter ("," methodParameter)*
func (p *Parser) methodParameterList(stopTokens ...token.Type) []ast.ParameterNode {
	return p.parameterList(p.methodParameter, stopTokens...)
}

// signatureParameter = identifier ["?"] [":" typeAnnotation]
func (p *Parser) signatureParameter() ast.ParameterNode {
	var typ ast.TypeNode
	var opt bool

	var paramName *token.Token

	switch p.lookahead.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		paramName = p.advance()
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		p.errorExpected("a lowercase identifier as the name of the declared signatureParameter")
		paramName = p.advance()
	default:
		p.errorExpected("an identifier as the name of the declared signatureParameter")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Span(),
			tok,
		)
	}
	lastSpan := paramName.Span()

	if questionTok, ok := p.matchOk(token.QUESTION); ok {
		opt = true
		lastSpan = questionTok.Span()
	}

	if p.match(token.COLON) {
		typ = p.intersectionType()
		lastSpan = typ.Span()
	}

	return ast.NewSignatureParameterNode(
		paramName.Span().Join(lastSpan.Span()),
		paramName.Value,
		typ,
		opt,
	)
}

// signatureParameterList = signatureParameter ("," signatureParameter)*
func (p *Parser) signatureParameterList(stopTokens ...token.Type) []ast.ParameterNode {
	return p.parameterList(p.signatureParameter, stopTokens...)
}

// identifierList = identifier ("," identifier)*
func (p *Parser) identifierList(stopTokens ...token.Type) []ast.IdentifierNode {
	return commaSeparatedList(p, p.identifier, stopTokens...)
}

// logicalOrExpression = logicalAndExpression |
// logicalOrExpression "||" logicalAndExpression |
// logicalOrExpression "??" logicalAndExpression |
// logicalOrExpression "|!" logicalAndExpression
func (p *Parser) logicalOrExpression() ast.ExpressionNode {
	return p.logicalExpression(p.logicalAndExpression, token.OR_OR, token.QUESTION_QUESTION, token.OR_BANG)
}

// logicalAndExpression = bitwiseOrExpression |
// logicalAndExpression "&&" bitwiseOrExpression |
// logicalAndExpression "&!" bitwiseOrExpression
func (p *Parser) logicalAndExpression() ast.ExpressionNode {
	return p.logicalExpression(p.bitwiseOrExpression, token.AND_AND, token.AND_BANG)
}

// bitwiseOrExpression = bitwiseXorExpression | bitwiseOrExpression "|" bitwiseXorExpression
func (p *Parser) bitwiseOrExpression() ast.ExpressionNode {
	if p.mode == withoutBitwiseOrMode {
		return p.bitwiseXorExpression()
	}
	return p.binaryExpression(p.bitwiseXorExpression, token.OR)
}

// bitwiseXorExpression = bitwiseAndExpression | bitwiseXorExpression "^" bitwiseAndExpression
func (p *Parser) bitwiseXorExpression() ast.ExpressionNode {
	return p.binaryExpression(p.bitwiseAndExpression, token.XOR)
}

// bitwiseAndExpression = equalityExpression | bitwiseAndExpression "&" equalityExpression
func (p *Parser) bitwiseAndExpression() ast.ExpressionNode {
	return p.binaryExpression(p.equalityExpression, token.AND)
}

// equalityExpression = comparisonExpression | equalityExpression EQUALITY_OP comparisonExpression
func (p *Parser) equalityExpression() ast.ExpressionNode {
	left := p.comparisonExpression()

	for p.lookahead.IsEqualityOperator() {
		operator := p.advance()

		p.swallowNewlines()

		p.indentedSection = true
		right := p.comparisonExpression()
		p.indentedSection = false

		left = ast.NewBinaryExpressionNode(
			left.Span().Join(right.Span()),
			operator,
			left,
			right,
		)
	}

	return left
}

// comparisonExpression = bitwiseShiftExpression | comparison COMP_OP bitwiseShiftExpression
func (p *Parser) comparisonExpression() ast.ExpressionNode {
	left := p.bitwiseShiftExpression()

	for p.lookahead.IsComparisonOperator() {
		operator := p.advance()

		p.swallowNewlines()

		p.indentedSection = true
		right := p.bitwiseShiftExpression()
		p.indentedSection = false

		left = ast.NewBinaryExpressionNode(
			left.Span().Join(right.Span()),
			operator,
			left,
			right,
		)
	}

	return left
}

// bitwiseShiftExpression = additiveExpression | bitwiseShiftExpression ("<<" | "<<<" | ">>" | ">>>") additiveExpression
func (p *Parser) bitwiseShiftExpression() ast.ExpressionNode {
	return p.binaryExpression(p.additiveExpression, token.LBITSHIFT, token.LTRIPLE_BITSHIFT, token.RBITSHIFT, token.RTRIPLE_BITSHIFT)
}

// additiveExpression = multiplicativeExpression | additiveExpression ("+" | "-") multiplicativeExpression
func (p *Parser) additiveExpression() ast.ExpressionNode {
	return p.binaryExpression(p.multiplicativeExpression, token.PLUS, token.MINUS)
}

// multiplicativeExpression = unaryExpression | multiplicativeExpression ("*" | "/" | "%") unaryExpression
func (p *Parser) multiplicativeExpression() ast.ExpressionNode {
	return p.binaryExpression(p.unaryExpression, token.STAR, token.SLASH, token.PERCENT)
}

// unaryExpression = powerExpression | ("!" | "-" | "+" | "~") unaryExpression
func (p *Parser) unaryExpression() ast.ExpressionNode {
	if operator, ok := p.matchOk(token.BANG, token.MINUS, token.PLUS, token.TILDE); ok {
		p.swallowNewlines()

		p.indentedSection = true
		right := p.unaryExpression()
		p.indentedSection = false

		return ast.NewUnaryExpressionNode(
			operator.Span().Join(right.Span()),
			operator,
			right,
		)
	}

	return p.powerExpression()
}

// powerExpression = methodCall | methodCall "**" powerExpression
func (p *Parser) powerExpression() ast.ExpressionNode {
	left := p.methodCall()

	if p.lookahead.Type != token.STAR_STAR {
		return left
	}

	operator := p.advance()
	p.swallowNewlines()

	p.indentedSection = true
	right := p.powerExpression()
	p.indentedSection = false

	return ast.NewBinaryExpressionNode(
		left.Span().Join(right.Span()),
		operator,
		left,
		right,
	)
}

// The boolean value indicates whether a comma was the last consumed token.
//
// positionalArgumentList = expressionWithoutModifier ("," expressionWithoutModifier)*
func (p *Parser) positionalArgumentList(stopTokens ...token.Type) ([]ast.ExpressionNode, bool) {
	var elements []ast.ExpressionNode
	if p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) && p.nextLookahead.Type == token.COLON {
		return elements, false
	}
	elements = append(elements, p.expressionWithoutModifier())

	for {
		if p.accept(token.END_OF_FILE) {
			break
		}

		for _, stopToken := range stopTokens {
			if p.lookahead.Type == stopToken {
				break
			}
		}
		if !p.match(token.COMMA) {
			break
		}
		p.swallowNewlines()
		if p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) && p.nextLookahead.Type == token.COLON {
			return elements, true
		}
		elements = append(elements, p.expressionWithoutModifier())
	}

	return elements, false
}

// namedArgument = identifier ":" expressionWithoutModifier
func (p *Parser) namedArgument() ast.NamedArgumentNode {
	ident, ok := p.matchOk(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER)
	if !ok {
		p.errorExpected("an identifier")
		errTok := p.advance()
		return ast.NewInvalidNode(
			errTok.Span(),
			errTok,
		)
	}
	colon, ok := p.consume(token.COLON)
	if !ok {
		return ast.NewInvalidNode(
			colon.Span(),
			colon,
		)
	}
	val := p.expressionWithoutModifier()

	return ast.NewNamedCallArgumentNode(
		ident.Span().Join(val.Span()),
		ident.Value,
		val,
	)
}

// namedArgumentList = namedArgument ("," namedArgument)*
func (p *Parser) namedArgumentList(stopTokens ...token.Type) []ast.NamedArgumentNode {
	return commaSeparatedList(p, p.namedArgument, stopTokens...)
}

const (
	expectedPublicMethodMessage = "a public method name (public identifier, keyword or overridable operator)"
	expectedMethodMessage       = "a method name (identifier, keyword or overridable operator)"
)

// methodCall = rangeLiteral |
// identifier ( "(" argumentList ")" | argumentList) |
// "self" ("."| "?.") (identifier | keyword | overridableOperator) ( "(" argumentList ")" | argumentList) |
// methodCall ("."| "?.") (publicIdentifier | keyword | overridableOperator) ( "(" argumentList ")" | argumentList)
func (p *Parser) methodCall() ast.ExpressionNode {
	// function call
	var receiver ast.ExpressionNode

	if p.accept(token.PRIVATE_IDENTIFIER, token.PUBLIC_IDENTIFIER) &&
		(p.nextLookahead.Type == token.LPAREN || p.nextLookahead.IsValidAsArgumentToNoParenFunctionCall()) {
		methodName := p.advance()
		lastSpan, posArgs, namedArgs, errToken := p.callArgumentList()
		if errToken != nil {
			return ast.NewInvalidNode(
				errToken.Span(),
				errToken,
			)
		}
		if lastSpan == nil {
			p.errorExpected("method arguments")
			errToken = p.advance()
			return ast.NewInvalidNode(
				errToken.Span(),
				errToken,
			)
		}

		receiver = ast.NewFunctionCallNode(
			methodName.Span().Join(lastSpan),
			methodName.Value,
			posArgs,
			namedArgs,
		)
	}

	// method call
	if receiver == nil {
		receiver = p.rangeOrArithmeticSequenceLiteral()
	}
	for {
		var opToken *token.Token

		if p.accept(token.NEWLINE) && p.acceptNext(token.DOT, token.QUESTION_DOT) {
			p.advance()
			opToken = p.advance()
		} else {
			t, ok := p.matchOk(token.DOT, token.QUESTION_DOT)
			if !ok {
				return receiver
			}
			opToken = t
		}

		_, selfReceiver := receiver.(*ast.SelfLiteralNode)
		p.swallowNewlines()

		if (!selfReceiver && p.accept(token.PRIVATE_IDENTIFIER)) || p.lookahead.IsNonOverridableOperator() {
			p.errorExpected(expectedPublicMethodMessage)
		} else if !p.lookahead.IsValidMethodName() {
			p.errorExpected(expectedPublicMethodMessage)
			p.indentedSection = true
			p.updateErrorMode(true)
			p.indentedSection = false
			errTok := p.advance()
			return ast.NewInvalidNode(
				errTok.Span(),
				errTok,
			)
		}

		methodNameTok := p.advance()
		methodName := methodNameTok.StringValue()

		lastSpan, posArgs, namedArgs, errToken := p.callArgumentList()
		if errToken != nil {
			return ast.NewInvalidNode(
				errToken.Span(),
				errToken,
			)
		}
		if lastSpan == nil {
			lastSpan = methodNameTok.Span()
		}

		receiver = ast.NewMethodCallNode(
			receiver.Span().Join(lastSpan),
			receiver,
			opToken.Type == token.QUESTION_DOT,
			methodName,
			posArgs,
			namedArgs,
		)
	}
}

// rangeOrArithmeticSequenceLiteral = constructorCall (".." | "...") [constructorCall] [":" constructorCall]
func (p *Parser) rangeOrArithmeticSequenceLiteral() ast.ExpressionNode {
	left := p.constructorCall()
	op, ok := p.matchOk(token.RANGE_OP, token.EXCLUSIVE_RANGE_OP)
	if !ok {
		return left
	}

	// endless arithmetic sequence
	if p.match(token.COLON) {
		step := p.constructorCall()
		return ast.NewArithmeticSequenceLiteralNode(
			left.Span().Join(step.Span()),
			op.Type == token.EXCLUSIVE_RANGE_OP,
			left,
			nil,
			step,
		)
	}

	if !p.lookahead.IsValidAsEndInRangeLiteral() {
		return ast.NewRangeLiteralNode(
			left.Span().Join(op.Span()),
			op.Type == token.EXCLUSIVE_RANGE_OP,
			left,
			nil,
		)
	}

	right := p.constructorCall()

	if p.match(token.COLON) {
		step := p.constructorCall()
		return ast.NewArithmeticSequenceLiteralNode(
			left.Span().Join(step.Span()),
			op.Type == token.EXCLUSIVE_RANGE_OP,
			left,
			right,
			step,
		)
	}

	return ast.NewRangeLiteralNode(
		left.Span().Join(right.Span()),
		op.Type == token.EXCLUSIVE_RANGE_OP,
		left,
		right,
	)
}

// beginlessRangeLiteral = (".." | "...") constructorCall
func (p *Parser) beginlessRangeLiteral() ast.ExpressionNode {
	op := p.advance()
	right := p.constructorCall()
	return ast.NewRangeLiteralNode(
		op.Span().Join(right.Span()),
		op.Type == token.EXCLUSIVE_RANGE_OP,
		nil,
		right,
	)
}

// callArgumentListInternal = (positionalArgumentList | namedArgumentList | positionalArgumentList "," namedArgumentList)
// callArgumentList = "(" callArgumentList ")" | callArgumentList
func (p *Parser) callArgumentList() (*position.Span, []ast.ExpressionNode, []ast.NamedArgumentNode, *token.Token) {
	if p.match(token.LPAREN) {
		p.swallowNewlines()
		if rparen, ok := p.matchOk(token.RPAREN); ok {
			return rparen.Span(),
				nil,
				nil,
				nil
		}
		posArgs, commaConsumed := p.positionalArgumentList()
		var namedArgs []ast.NamedArgumentNode
		if len(posArgs) == 0 || len(posArgs) > 0 && commaConsumed {
			namedArgs = p.namedArgumentList()
		}
		p.swallowNewlines()
		rparen, ok := p.consume(token.RPAREN)
		if !ok {
			return nil,
				nil,
				nil,
				rparen
		}

		return rparen.Span(),
			posArgs,
			namedArgs,
			nil
	}

	// no parentheses
	if !p.lookahead.IsValidAsArgumentToNoParenFunctionCall() {
		return nil,
			nil,
			nil,
			nil
	}

	posArgs, commaConsumed := p.positionalArgumentList()
	span := position.SpanOfLastElement(posArgs)
	var namedArgs []ast.NamedArgumentNode
	if len(posArgs) == 0 || len(posArgs) > 0 && commaConsumed {
		namedArgs = p.namedArgumentList()
		span = position.SpanOfLastElement(namedArgs)
	}

	return span,
		posArgs,
		namedArgs,
		nil
}

// constructorCall = constantLookup |
// strictConstantLookup "(" argumentList ")" |
// strictConstantLookup argumentList
func (p *Parser) constructorCall() ast.ExpressionNode {
	if !p.accept(token.PRIVATE_CONSTANT, token.PUBLIC_CONSTANT, token.SCOPE_RES_OP) {
		return p.constantLookup()
	}

	constant := p.strictConstantLookup()

	lastSpan, posArgs, namedArgs, errToken := p.callArgumentList()
	if errToken != nil {
		return ast.NewInvalidNode(
			errToken.Span(),
			errToken,
		)
	}
	if lastSpan == nil {
		return constant
	}

	return ast.NewConstructorCallNode(
		constant.Span().Join(lastSpan),
		constant,
		posArgs,
		namedArgs,
	)
}

const privateConstantAccessMessage = "can't access a private constant from the outside"

// constantLookup = primaryExpression | "::" publicConstant | constantLookup "::" publicConstant
func (p *Parser) constantLookup() ast.ExpressionNode {
	var left ast.ExpressionNode
	if tok, ok := p.matchOk(token.SCOPE_RES_OP); ok {
		if p.accept(token.PRIVATE_CONSTANT) {
			p.errorUnexpected(privateConstantAccessMessage)
		}
		right := p.constant()
		left = ast.NewConstantLookupNode(
			tok.Span().Join(right.Span()),
			nil,
			right,
		)
	} else {
		left = p.primaryExpression()
	}

	for p.lookahead.Type == token.SCOPE_RES_OP {
		p.advance()

		p.swallowNewlines()
		if p.accept(token.PRIVATE_CONSTANT) {
			p.errorUnexpected(privateConstantAccessMessage)
		}
		right := p.constant()

		left = ast.NewConstantLookupNode(
			left.Span().Join(right.Span()),
			left,
			right,
		)
	}

	return left
}

// strictConstantLookup = constant | "::" publicConstant | strictConstantLookup "::" publicConstant
func (p *Parser) strictConstantLookup() ast.ComplexConstantNode {
	var left ast.ComplexConstantNode
	if tok, ok := p.matchOk(token.SCOPE_RES_OP); ok {
		if p.accept(token.PRIVATE_CONSTANT) {
			p.errorUnexpected(privateConstantAccessMessage)
		}
		right := p.constant()
		left = ast.NewConstantLookupNode(
			tok.Span().Join(right.Span()),
			nil,
			right,
		)
	} else {
		left = p.constant()
	}

	for p.lookahead.Type == token.SCOPE_RES_OP {
		p.advance()

		p.swallowNewlines()
		if p.accept(token.PRIVATE_CONSTANT) {
			p.errorUnexpected(privateConstantAccessMessage)
		}
		right := p.constant()

		left = ast.NewConstantLookupNode(
			left.Span().Join(right.Span()),
			left,
			right,
		)
	}

	return left
}

// constant = privateConstant | publicConstant
func (p *Parser) constant() ast.ConstantNode {
	if tok, ok := p.matchOk(token.PRIVATE_CONSTANT); ok {
		return ast.NewPrivateConstantNode(
			tok.Span(),
			tok.Value,
		)
	}

	if tok, ok := p.matchOk(token.PUBLIC_CONSTANT); ok {
		return ast.NewPublicConstantNode(
			tok.Span(),
			tok.Value,
		)
	}

	p.errorExpected("a constant")
	tok := p.advance()
	p.mode = panicMode
	return ast.NewInvalidNode(
		tok.Span(),
		tok,
	)
}

func (p *Parser) primaryExpression() ast.ExpressionNode {
	switch p.lookahead.Type {
	case token.TRUE:
		tok := p.advance()
		return ast.NewTrueLiteralNode(tok.Span())
	case token.FALSE:
		tok := p.advance()
		return ast.NewFalseLiteralNode(tok.Span())
	case token.NIL:
		tok := p.advance()
		return ast.NewNilLiteralNode(tok.Span())
	case token.SELF:
		return p.selfLiteral()
	case token.BREAK:
		tok := p.advance()
		return ast.NewBreakExpressionNode(tok.Span())
	case token.RETURN:
		return p.returnExpression()
	case token.CONTINUE:
		return p.continueExpression()
	case token.THROW:
		return p.throwExpression()
	case token.LPAREN:
		p.advance()
		if p.mode == withoutBitwiseOrMode {
			p.mode = normalMode
		}
		expr := p.expressionWithModifier()
		p.consume(token.RPAREN)
		return expr
	case token.LBRACKET:
		return p.listLiteral()
	case token.TUPLE_LITERAL_BEG:
		return p.tupleLiteral()
	case token.SET_LITERAL_BEG:
		return p.setLiteral()
	case token.WORD_LIST_BEG:
		return p.wordListLiteral()
	case token.WORD_TUPLE_BEG:
		return p.wordTupleLiteral()
	case token.WORD_SET_BEG:
		return p.wordSetLiteral()
	case token.SYMBOL_LIST_BEG:
		return p.symbolListLiteral()
	case token.SYMBOL_TUPLE_BEG:
		return p.symbolTupleLiteral()
	case token.SYMBOL_SET_BEG:
		return p.symbolSetLiteral()
	case token.HEX_LIST_BEG:
		return p.hexListLiteral()
	case token.HEX_TUPLE_BEG:
		return p.hexTupleLiteral()
	case token.HEX_SET_BEG:
		return p.hexSetLiteral()
	case token.BIN_LIST_BEG:
		return p.binListLiteral()
	case token.BIN_TUPLE_BEG:
		return p.binTupleLiteral()
	case token.BIN_SET_BEG:
		return p.binSetLiteral()
	case token.LBRACE:
		return p.mapLiteral()
	case token.RECORD_LITERAL_BEG:
		return p.recordLiteral()
	case token.CHAR_LITERAL:
		return p.charLiteral()
	case token.RAW_CHAR_LITERAL:
		return p.rawCharLiteral()
	case token.RAW_STRING:
		return p.rawStringLiteral()
	case token.STRING_BEG:
		return p.stringLiteral()
	case token.COLON:
		return p.symbolLiteral()
	case token.OR, token.OR_OR:
		return p.closureExpression()
	case token.DOC_COMMENT:
		return p.docComment()
	case token.VAR:
		return p.variableDeclaration()
	case token.VAL:
		return p.valueDeclaration()
	case token.CONST:
		return p.constantDeclaration()
	case token.DEF:
		return p.methodDefinition()
	case token.INIT:
		return p.initDefinition()
	case token.IF:
		return p.ifExpression()
	case token.DO:
		return p.doExpression()
	case token.UNLESS:
		return p.unlessExpression()
	case token.WHILE:
		return p.whileExpression()
	case token.UNTIL:
		return p.untilExpression()
	case token.LOOP:
		return p.loopExpression()
	case token.FOR:
		return p.forExpression()
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		return p.identifierOrClosure()
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		return p.constant()
	case token.INT:
		tok := p.advance()
		return ast.NewIntLiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.INT64:
		tok := p.advance()
		return ast.NewInt64LiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.UINT64:
		tok := p.advance()
		return ast.NewUInt64LiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.INT32:
		tok := p.advance()
		return ast.NewInt32LiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.UINT32:
		tok := p.advance()
		return ast.NewUInt32LiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.INT16:
		tok := p.advance()
		return ast.NewInt16LiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.UINT16:
		tok := p.advance()
		return ast.NewUInt16LiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.INT8:
		tok := p.advance()
		return ast.NewInt8LiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.UINT8:
		tok := p.advance()
		return ast.NewUInt8LiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.FLOAT:
		tok := p.advance()
		return ast.NewFloatLiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.BIG_FLOAT:
		tok := p.advance()
		return ast.NewBigFloatLiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.FLOAT64:
		tok := p.advance()
		return ast.NewFloat64LiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.FLOAT32:
		tok := p.advance()
		return ast.NewFloat32LiteralNode(
			tok.Span(),
			tok.Value,
		)
	case token.ERROR:
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Span(),
			tok,
		)
	case token.CLASS:
		return p.classDeclaration()
	case token.MODULE:
		return p.moduleDeclaration()
	case token.MIXIN:
		return p.mixinDeclaration()
	case token.INTERFACE:
		return p.interfaceDeclaration()
	case token.STRUCT:
		return p.structDeclaration()
	case token.TYPEDEF:
		return p.typeDefinition()
	case token.ALIAS:
		return p.aliasExpression()
	case token.SIG:
		return p.methodSignatureDefinition()
	case token.INCLUDE:
		return p.includeExpression()
	case token.EXTEND:
		return p.extendExpression()
	case token.ENHANCE:
		return p.enhanceExpression()
	case token.RANGE_OP, token.EXCLUSIVE_RANGE_OP:
		return p.beginlessRangeLiteral()
	default:
		p.errorExpected("an expression")
		p.updateErrorMode(true)
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Span(),
			tok,
		)
	}
}

type specialCollectionLiteralConstructor[Element ast.ExpressionNode] func(*position.Span, []Element) ast.ExpressionNode

// specialCollectionLiteral = beginTokenType (elementProduction)* endTokenType
func specialCollectionLiteral[Element ast.ExpressionNode](p *Parser, elementProduction func() Element, constructor specialCollectionLiteralConstructor[Element], endTokenType token.Type) ast.ExpressionNode {
	begTok := p.advance()
	content := repeatedProduction(p, elementProduction, endTokenType)
	endTok, ok := p.consume(endTokenType)

	if !ok {
		return ast.NewInvalidNode(endTok.Span(), endTok)
	}

	return constructor(
		begTok.Span().Join(endTok.Span()),
		content,
	)
}

// wordListLiteral = "\w[" (rawString)* "]"
func (p *Parser) wordListLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.wordCollectionElement,
		ast.NewWordListLiteralNodeI,
		token.WORD_LIST_END,
	)
}

// wordTupleLiteral = "%w[" (rawString)* "]"
func (p *Parser) wordTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.wordCollectionElement,
		ast.NewWordTupleLiteralNodeI,
		token.WORD_TUPLE_END,
	)
}

// wordSetLiteral = "^w[" (rawString)* "]"
func (p *Parser) wordSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.wordCollectionElement,
		ast.NewWordSetLiteralNodeI,
		token.WORD_SET_END,
	)
}

// symbolListLiteral = "\s[" (rawString)* "]"
func (p *Parser) symbolListLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolListLiteralNodeI,
		token.SYMBOL_LIST_END,
	)
}

// symbolTupleLiteral = "%s[" (rawString)* "]"
func (p *Parser) symbolTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolTupleLiteralNodeI,
		token.SYMBOL_TUPLE_END,
	)
}

// symbolSetLiteral = "^s[" (rawString)* "]"
func (p *Parser) symbolSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolSetLiteralNodeI,
		token.SYMBOL_SET_END,
	)
}

// hexListLiteral = "\x[" (HEX_INT)* "]"
func (p *Parser) hexListLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewHexListLiteralNodeI,
		token.HEX_LIST_END,
	)
}

// hexTupleLiteral = "%x[" (HEX_INT)* "]"
func (p *Parser) hexTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewHexTupleLiteralNodeI,
		token.HEX_TUPLE_END,
	)
}

// hexSetLiteral = "^x[" (HEX_INT)* "]"
func (p *Parser) hexSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewHexSetLiteralNodeI,
		token.HEX_SET_END,
	)
}

// binListLiteral = "\b[" (BIN_INT)* "]"
func (p *Parser) binListLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewBinListLiteralNodeI,
		token.BIN_LIST_END,
	)
}

// binTupleLiteral = "%b[" (BIN_INT)* "]"
func (p *Parser) binTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewBinTupleLiteralNodeI,
		token.BIN_TUPLE_END,
	)
}

// binSetLiteral = "^b[" (BIN_INT)* "]"
func (p *Parser) binSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewBinSetLiteralNodeI,
		token.BIN_SET_END,
	)
}

type listLikeConstructor func(*position.Span, []ast.ExpressionNode) ast.ExpressionNode
type collectionElementsProduction func(...token.Type) []ast.ExpressionNode

// collectionLiteral = startTok [elementsProduction] endTok
func (p *Parser) collectionLiteral(endTokType token.Type, elementsProduction collectionElementsProduction, constructor listLikeConstructor) ast.ExpressionNode {
	startTok := p.advance()
	p.swallowNewlines()

	if endTok, ok := p.matchOk(endTokType); ok {
		return constructor(
			startTok.Span().Join(endTok.Span()),
			nil,
		)
	}

	elements := elementsProduction(endTokType)
	p.swallowNewlines()
	endTok, ok := p.consume(endTokType)
	if !ok {
		return ast.NewInvalidNode(
			endTok.Span(),
			endTok,
		)
	}

	return constructor(
		startTok.Span().Join(endTok.Span()),
		elements,
	)
}

// collectionElementModifier = subProduction |
// subProduction ("if" | "unless") expressionWithoutModifier |
// subProduction "if" expressionWithoutModifier "else" expressionWithoutModifier |
// subProduction "for" identifierList "in" expressionWithoutModifier
func (p *Parser) collectionElementModifier(subProduction func() ast.ExpressionNode) ast.ExpressionNode {
	left := subProduction()

	switch p.lookahead.Type {
	case token.UNLESS:
		mod := p.advance()
		p.swallowNewlines()
		right := p.expressionWithoutModifier()
		return ast.NewModifierNode(
			left.Span().Join(right.Span()),
			mod,
			left,
			right,
		)
	case token.IF:
		ifTok := p.advance()
		p.swallowNewlines()
		cond := p.expressionWithoutModifier()
		if p.lookahead.Type == token.ELSE {
			p.advance()
			p.swallowNewlines()
			elseExpr := p.expressionWithoutModifier()
			return ast.NewModifierIfElseNode(
				left.Span().Join(elseExpr.Span()),
				left,
				cond,
				elseExpr,
			)
		}
		return ast.NewModifierNode(
			left.Span().Join(cond.Span()),
			ifTok,
			left,
			cond,
		)
	case token.FOR:
		p.advance()
		p.swallowNewlines()
		params := p.identifierList(token.IN)
		inTok, ok := p.consume(token.IN)
		if !ok {
			return ast.NewInvalidNode(inTok.Span(), inTok)
		}
		p.swallowNewlines()
		inExpr := p.expressionWithoutModifier()
		return ast.NewModifierForInNode(
			left.Span().Join(inExpr.Span()),
			left,
			params,
			inExpr,
		)
	}

	return left
}

// "{" [mapLiteralElements] "}"
func (p *Parser) mapLiteral() ast.ExpressionNode {
	return p.collectionLiteral(token.RBRACE, p.mapLiteralElements, ast.NewMapLiteralNodeI)
}

// "%{" [mapLiteralElements] "}"
func (p *Parser) recordLiteral() ast.ExpressionNode {
	return p.collectionLiteral(token.RBRACE, p.mapLiteralElements, ast.NewRecordLiteralNodeI)
}

// listLiteral = "[" [listLikeLiteralElements] "]"
func (p *Parser) listLiteral() ast.ExpressionNode {
	return p.collectionLiteral(token.RBRACKET, p.listLikeLiteralElements, ast.NewListLiteralNodeI)
}

// tupleLiteral = "%[" [listLikeLiteralElements] "]"
func (p *Parser) tupleLiteral() ast.ExpressionNode {
	return p.collectionLiteral(token.RBRACKET, p.listLikeLiteralElements, ast.NewTupleLiteralNodeI)
}

// listLikeLiteralElements = listLikeLiteralElement ("," listLikeLiteralElement)*
func (p *Parser) listLikeLiteralElements(stopTokens ...token.Type) []ast.ExpressionNode {
	return commaSeparatedList(p, p.listLikeLiteralElement, stopTokens...)
}

// listLikeLiteralElement = keyValueExpression |
// keyValueExpression ("if" | "unless") expressionWithoutModifier |
// keyValueExpression "if" expressionWithoutModifier "else" expressionWithoutModifier |
// keyValueExpression "for" identifierList "in" expressionWithoutModifier
func (p *Parser) listLikeLiteralElement() ast.ExpressionNode {
	return p.collectionElementModifier(p.keyValueExpression)
}

// mapLiteralElements = mapLiteralElement ("," mapLiteralElement)*
func (p *Parser) mapLiteralElements(stopTokens ...token.Type) []ast.ExpressionNode {
	return commaSeparatedList(p, p.mapLiteralElement, stopTokens...)
}

// mapLiteralElement = keyValueMapExpression |
// keyValueMapExpression ("if" | "unless") expressionWithoutModifier |
// keyValueMapExpression "if" expressionWithoutModifier "else" expressionWithoutModifier |
// keyValueMapExpression "for" identifierList "in" expressionWithoutModifier
func (p *Parser) mapLiteralElement() ast.ExpressionNode {
	return p.collectionElementModifier(p.keyValueMapExpression)
}

// keyValueMapExpression = (identifier | constant) |
// (identifier | constant) ":" expressionWithoutModifier |
// expressionWithoutModifier "=>" expressionWithoutModifier
func (p *Parser) keyValueMapExpression() ast.ExpressionNode {
	if p.accept(
		token.PUBLIC_IDENTIFIER,
		token.PRIVATE_IDENTIFIER,
		token.PUBLIC_CONSTANT,
		token.PRIVATE_CONSTANT,
	) &&
		p.acceptNext(token.COLON) {
		key := p.advance()
		p.advance()
		p.swallowNewlines()
		val := p.expressionWithoutModifier()
		return ast.NewSymbolKeyValueExpressionNode(
			key.Span().Join(val.Span()),
			key.Value,
			val,
		)
	}
	key := p.expressionWithoutModifier()
	if !p.match(token.THICK_ARROW) {
		switch key.(type) {
		case *ast.PublicIdentifierNode, *ast.PrivateIdentifierNode,
			*ast.PublicConstantNode, *ast.PrivateConstantNode:
			return key
		default:
			p.errorMessageSpan("expected a key-value pair, map literals should consist of key-value pairs", key.Span())
			return key
		}
	}

	p.swallowNewlines()
	val := p.expressionWithoutModifier()

	return ast.NewKeyValueExpressionNode(
		key.Span().Join(val.Span()),
		key,
		val,
	)
}

// keyValueExpression = expressionWithoutModifier |
// expressionWithoutModifier "=>" expressionWithoutModifier
func (p *Parser) keyValueExpression() ast.ExpressionNode {
	key := p.expressionWithoutModifier()
	if p.match(token.THICK_ARROW) {
		p.swallowNewlines()
		value := p.expressionWithoutModifier()
		return ast.NewKeyValueExpressionNode(
			key.Span().Join(value.Span()),
			key,
			value,
		)
	}

	return key
}

// setLiteral = "^[" [setLiteralElements] "]"
func (p *Parser) setLiteral() ast.ExpressionNode {
	return p.collectionLiteral(token.RBRACKET, p.setLiteralElements, ast.NewSetLiteralNodeI)
}

// setLiteralElements = setLiteralElement ("," setLiteralElement)*
func (p *Parser) setLiteralElements(stopTokens ...token.Type) []ast.ExpressionNode {
	return commaSeparatedList(p, p.setLiteralElement, stopTokens...)
}

// setLiteralElement = expressionWithoutModifier |
// expressionWithoutModifier ("if" | "unless") expressionWithoutModifier |
// expressionWithoutModifier "if" expressionWithoutModifier "else" expressionWithoutModifier |
// expressionWithoutModifier "for" identifierList "in" expressionWithoutModifier
func (p *Parser) setLiteralElement() ast.ExpressionNode {
	return p.collectionElementModifier(p.expressionWithoutModifier)
}

// selfLiteral = "self"
func (p *Parser) selfLiteral() *ast.SelfLiteralNode {
	tok := p.advance()
	return ast.NewSelfLiteralNode(tok.Span())
}

// genericConstantList = genericConstant ("," genericConstant)*
func (p *Parser) genericConstantList(stopTokens ...token.Type) []ast.ComplexConstantNode {
	return commaSeparatedListWithoutTerminator(p, p.genericConstant, stopTokens...)
}

// includeExpression = "include" genericConstantList
func (p *Parser) includeExpression() *ast.IncludeExpressionNode {
	return includelikeExpression(p, ast.NewIncludeExpressionNode)
}

// extendExpression = "extend" genericConstantList
func (p *Parser) extendExpression() *ast.ExtendExpressionNode {
	return includelikeExpression(p, ast.NewExtendExpressionNode)
}

// enhanceExpression = "enhance" genericConstantList
func (p *Parser) enhanceExpression() *ast.EnhanceExpressionNode {
	return includelikeExpression(p, ast.NewEnhanceExpressionNode)
}

// methodDefinition = "sig" methodName ["(" signatureParameterList ")"] [":" typeAnnotation] ["!" typeAnnotation]
func (p *Parser) methodSignatureDefinition() ast.ExpressionNode {
	var params []ast.ParameterNode
	var returnType ast.TypeNode
	var throwType ast.TypeNode
	var span *position.Span

	sigTok := p.advance()
	span = sigTok.Span()

	methodName, mSpan := p.methodName()
	span = span.Join(mSpan)

	if p.match(token.LPAREN) {
		if rparen, ok := p.matchOk(token.RPAREN); ok {
			span = span.Join(rparen.Span())
		} else {
			params = p.signatureParameterList(token.RPAREN)

			rparen, ok := p.consume(token.RPAREN)
			span = span.Join(rparen.Span())
			if !ok {
				return ast.NewInvalidNode(
					rparen.Span(),
					rparen,
				)
			}
		}
	}

	// return type
	if p.match(token.COLON) {
		returnType = p.typeAnnotation()
		span = span.Join(returnType.Span())
	}

	// throw type
	if p.match(token.BANG) {
		throwType = p.typeAnnotation()
		span = span.Join(throwType.Span())
	}

	return ast.NewMethodSignatureDefinitionNode(
		span,
		methodName,
		params,
		returnType,
		throwType,
	)
}

// typeDeclaration = "alias" methodName methodName
func (p *Parser) aliasExpression() ast.ExpressionNode {
	aliasTok := p.advance()
	p.swallowNewlines()

	var (
		lastSpan *position.Span
		oldName  string
		newName  string
	)
	newName, _ = p.methodName()
	p.swallowNewlines()

	oldName, lastSpan = p.methodName()

	return ast.NewAliasExpressionNode(
		aliasTok.Span().Join(lastSpan),
		newName,
		oldName,
	)
}

// typeDeclaration = "typedef" strictConstantLookup "=" typeAnnotation
func (p *Parser) typeDefinition() ast.ExpressionNode {
	typedefTok := p.advance()

	name := p.strictConstantLookup()
	equalTok, ok := p.consume(token.EQUAL_OP)
	if !ok {
		return ast.NewInvalidNode(equalTok.Span(), equalTok)
	}
	p.swallowNewlines()

	typ := p.typeAnnotation()
	return ast.NewTypeDefinitionNode(
		typedefTok.Span().Join(typ.Span()),
		name,
		typ,
	)
}

func (p *Parser) methodName() (string, *position.Span) {
	var methodName string
	var span *position.Span

	if p.lookahead.IsValidRegularMethodName() {
		methodNameTok := p.advance()
		methodName = methodNameTok.StringValue()
		span = methodNameTok.Span()
		if tok, ok := p.matchOk(token.EQUAL_OP); ok {
			methodName += "="
			span = span.Join(tok.Span())
		}
	} else if p.accept(token.LBRACKET) && p.acceptNext(token.RBRACKET) {
		// [
		tok := p.advance()
		span = tok.Span()

		// ]
		tok = p.advance()
		span = span.Join(tok.Span())
		methodName = "[]"

		if tok, ok := p.matchOk(token.EQUAL_OP); ok {
			methodName += "="
			span = span.Join(tok.Span())
		}
	} else {
		if !p.lookahead.IsOverridableOperator() {
			p.errorExpected("a method name (identifier, overridable operator)")
		}
		tok := p.advance()
		methodName = tok.StringValue()
		span = tok.Span()
	}

	return methodName, span
}

// methodDefinition = "def" methodName ["(" methodParameterList ")"] [":" typeAnnotation] ["!" typeAnnotation] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) methodDefinition() ast.ExpressionNode {
	var params []ast.ParameterNode
	var returnType ast.TypeNode
	var throwType ast.TypeNode
	var body []ast.StatementNode
	var span *position.Span

	defTok := p.advance()
	p.swallowNewlines()
	methodName, _ := p.methodName()

	if p.match(token.LPAREN) {
		p.swallowNewlines()
		if !p.match(token.RPAREN) {
			params = p.methodParameterList(token.RPAREN, token.STAR)

			p.swallowNewlines()
			if tok, ok := p.consume(token.RPAREN); !ok {
				return ast.NewInvalidNode(
					tok.Span(),
					tok,
				)
			}
		}
	}

	// return type
	if p.match(token.COLON) {
		returnType = p.typeAnnotation()
	}

	// throw type
	if p.match(token.BANG) {
		throwType = p.typeAnnotation()
	}

	lastSpan, body, multiline := p.statementBlockWithThen(token.END)
	if lastSpan != nil {
		span = defTok.Span().Join(lastSpan)
	} else {
		span = defTok.Span()
	}

	if multiline {
		if len(body) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(body) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewMethodDefinitionNode(
		span,
		methodName,
		params,
		returnType,
		throwType,
		body,
	)
}

// initDefinition = "init" ["(" methodParameterList ")"] ["!" typeAnnotation] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) initDefinition() ast.ExpressionNode {
	var params []ast.ParameterNode
	var throwType ast.TypeNode
	var body []ast.StatementNode
	var span *position.Span

	initTok := p.advance()

	// methodParameterList
	if p.match(token.LPAREN) {
		p.swallowNewlines()
		if !p.match(token.RPAREN) {
			params = p.methodParameterList(token.RPAREN)

			p.swallowNewlines()
			if tok, ok := p.consume(token.RPAREN); !ok {
				return ast.NewInvalidNode(
					tok.Span(),
					tok,
				)
			}
		}
	}

	// throw type
	if p.match(token.BANG) {
		throwType = p.typeAnnotation()
	}

	lastSpan, body, multiline := p.statementBlockWithThen(token.END)
	if lastSpan != nil {
		span = initTok.Span().Join(lastSpan)
	} else {
		span = initTok.Span()
	}

	if multiline {
		if len(body) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(body) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewInitDefinitionNode(
		span,
		params,
		throwType,
		body,
	)
}

// typeVariable = ["+" | "-"] constant ["<" strictConstantLookup]
func (p *Parser) typeVariable() ast.TypeVariableNode {
	variance := ast.INVARIANT
	var firstSpan *position.Span
	var lastSpan *position.Span
	var upperBound ast.ComplexConstantNode

	switch p.lookahead.Type {
	case token.PLUS:
		plusTok := p.advance()
		firstSpan = plusTok.Span()
		variance = ast.COVARIANT
	case token.MINUS:
		minusTok := p.advance()
		firstSpan = minusTok.Span()
		variance = ast.CONTRAVARIANT
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
	default:
		errTok := p.advance()
		p.errorExpected("a type variable")
		return ast.NewInvalidNode(
			errTok.Span(),
			errTok,
		)
	}

	if !p.accept(token.PRIVATE_CONSTANT, token.PUBLIC_CONSTANT) {
		errTok := p.advance()
		return ast.NewInvalidNode(
			errTok.Span(),
			errTok,
		)
	}
	nameTok := p.advance()
	if firstSpan == nil {
		firstSpan = nameTok.Span()
	}
	lastSpan = nameTok.Span()

	if p.match(token.LESS) {
		upperBound = p.strictConstantLookup()
		lastSpan = upperBound.Span()
	}

	return ast.NewVariantTypeVariableNode(
		firstSpan.Join(lastSpan),
		variance,
		nameTok.Value,
		upperBound,
	)
}

// typeVariableList = typeVariable ("," typeVariable)*
func (p *Parser) typeVariableList(stopTokens ...token.Type) []ast.TypeVariableNode {
	return commaSeparatedList(p, p.typeVariable, stopTokens...)
}

// classDeclaration = "class" [constantLookup] ["[" typeVariableList "]"] ["<" genericConstant] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) classDeclaration() ast.ExpressionNode {
	classTok := p.advance()
	var superclass ast.ExpressionNode
	var constant ast.ExpressionNode
	var typeVars []ast.TypeVariableNode
	if !p.accept(token.LESS, token.LBRACKET, token.SEMICOLON, token.NEWLINE) {
		constant = p.constantLookup()

		switch constant.(type) {
		case *ast.PublicConstantNode,
			*ast.PrivateConstantNode,
			*ast.ConstantLookupNode:
		default:
			p.errorMessageSpan("invalid class name, expected a constant", constant.Span())
		}
	}
	var span *position.Span

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeVariableList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Span(),
					errTok,
				)
			}
		}
	}

	if p.match(token.LESS) {
		superclass = p.genericConstant()
	}

	lastSpan, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastSpan != nil {
		span = classTok.Span().Join(lastSpan)
	} else {
		span = classTok.Span()
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewClassDeclarationNode(
		span,
		constant,
		typeVars,
		superclass,
		thenBody,
	)
}

// moduleDeclaration = "module" [constantLookup] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) moduleDeclaration() ast.ExpressionNode {
	moduleTok := p.advance()
	var constant ast.ExpressionNode
	if !p.accept(token.LBRACKET, token.SEMICOLON, token.NEWLINE) {
		constant = p.constantLookup()

		switch constant.(type) {
		case *ast.PublicConstantNode,
			*ast.PrivateConstantNode,
			*ast.ConstantLookupNode:
		default:
			p.errorMessageSpan("invalid module name, expected a constant", constant.Span())
		}
	}
	var span *position.Span

	if lbracket, ok := p.matchOk(token.LBRACKET); ok {
		errPos := lbracket.Span()
		if p.accept(token.RBRACKET) {
			rbracket := p.advance()
			errPos = errPos.Join(rbracket.Span())
		} else {
			p.typeVariableList()
			rbracket, ok := p.consume(token.RBRACKET)
			if !ok {
				return ast.NewInvalidNode(
					rbracket.Span(),
					rbracket,
				)
			}
			errPos = errPos.Join(rbracket.Span())
		}
		p.errorMessageSpan("modules can't be generic", errPos)
	}

	lastSpan, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastSpan != nil {
		span = moduleTok.Span().Join(lastSpan)
	} else {
		span = moduleTok.Span()
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewModuleDeclarationNode(
		span,
		constant,
		thenBody,
	)
}

// mixinDeclaration = "mixin" [constantLookup] ["[" typeVariableList "]"] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) mixinDeclaration() ast.ExpressionNode {
	mixinTok := p.advance()
	var constant ast.ExpressionNode
	var typeVars []ast.TypeVariableNode
	if !p.accept(token.LBRACKET, token.SEMICOLON, token.NEWLINE) {
		constant = p.constantLookup()

		switch constant.(type) {
		case *ast.PublicConstantNode,
			*ast.PrivateConstantNode,
			*ast.ConstantLookupNode:
		default:
			p.errorMessageSpan("invalid mixin name, expected a constant", constant.Span())
		}
	}
	var span *position.Span

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeVariableList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Span(),
					errTok,
				)
			}
		}
	}

	lastSpan, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastSpan != nil {
		span = mixinTok.Span().Join(lastSpan)
	} else {
		span = mixinTok.Span()
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewMixinDeclarationNode(
		span,
		constant,
		typeVars,
		thenBody,
	)
}

// interfaceDeclaration = "interface" [constantLookup] ["[" typeVariableList "]"] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) interfaceDeclaration() ast.ExpressionNode {
	interfaceTok := p.advance()
	var constant ast.ExpressionNode
	var typeVars []ast.TypeVariableNode
	if !p.accept(token.LBRACKET, token.SEMICOLON, token.NEWLINE) {
		constant = p.constantLookup()

		switch constant.(type) {
		case *ast.PublicConstantNode,
			*ast.PrivateConstantNode,
			*ast.ConstantLookupNode:
		default:
			p.errorMessageSpan("invalid interface name, expected a constant", constant.Span())
		}
	}
	var span *position.Span

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeVariableList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Span(),
					errTok,
				)
			}
		}
	}

	lastSpan, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastSpan != nil {
		span = interfaceTok.Span().Join(lastSpan)
	} else {
		span = interfaceTok.Span()
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewInterfaceDeclarationNode(
		span,
		constant,
		typeVars,
		thenBody,
	)
}

// structDeclaration = "struct" [constantLookup] ["[" typeVariableList "]"] ((SEPARATOR [structBodyStatements] "end") | ("then" formalParameter))
func (p *Parser) structDeclaration() ast.ExpressionNode {
	structTok := p.advance()
	var constant ast.ExpressionNode
	var typeVars []ast.TypeVariableNode
	if !p.accept(token.LBRACKET, token.SEMICOLON, token.NEWLINE) {
		constant = p.constantLookup()

		switch constant.(type) {
		case *ast.PublicConstantNode,
			*ast.PrivateConstantNode,
			*ast.ConstantLookupNode:
		default:
			p.errorMessageSpan("invalid struct name, expected a constant", constant.Span())
		}
	}
	var span *position.Span

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeVariableList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Span(),
					errTok,
				)
			}
		}
	}

	lastSpan, thenBody, multiline := p.structBodyStatementBlockWithThen(token.END)
	if lastSpan != nil {
		span = structTok.Span().Join(lastSpan)
	} else {
		span = structTok.Span()
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewStructDeclarationNode(
		span,
		constant,
		typeVars,
		thenBody,
	)
}

// variableDeclaration = "var" identifier [":" typeAnnotation] ["=" expressionWithoutModifier]
func (p *Parser) variableDeclaration() ast.ExpressionNode {
	varTok := p.advance()
	var init ast.ExpressionNode
	var typ ast.TypeNode

	varName, ok := p.matchOk(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER, token.INSTANCE_VARIABLE)
	if !ok {
		p.errorExpected("an identifier as the name of the declared variable")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Span(),
			tok,
		)
	}
	lastSpan := varName.Span()

	if p.match(token.COLON) {
		typ = p.typeAnnotation()
		lastSpan = typ.Span()
	}

	if p.match(token.EQUAL_OP) {
		p.swallowNewlines()
		init = p.expressionWithoutModifier()
		lastSpan = init.Span()
	}

	return ast.NewVariableDeclarationNode(
		varTok.Span().Join(lastSpan),
		varName,
		typ,
		init,
	)
}

// valueDeclaration = "val" identifier [":" typeAnnotation] ["=" expressionWithoutModifier]
func (p *Parser) valueDeclaration() ast.ExpressionNode {
	varTok := p.advance()
	var init ast.ExpressionNode
	var typ ast.TypeNode
	var lastSpan *position.Span
	var valName *token.Token

	if v, ok := p.matchOk(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER); ok {
		valName = v
		lastSpan = v.Span()
	} else if v, ok := p.matchOk(token.INSTANCE_VARIABLE); ok {
		p.errorMessageSpan("instance variables can't be declared using `val`", v.Span())
		lastSpan = v.Span()
		valName = v
	} else {
		p.errorExpected("an identifier as the name of the declared value")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Span(),
			tok,
		)
	}

	if p.match(token.COLON) {
		typ = p.typeAnnotation()
		lastSpan = typ.Span()
	}

	if p.match(token.EQUAL_OP) {
		p.swallowNewlines()
		init = p.expressionWithoutModifier()
		lastSpan = init.Span()
	}

	return ast.NewValueDeclarationNode(
		varTok.Span().Join(lastSpan),
		valName,
		typ,
		init,
	)
}

// constantDeclaration = "const" identifier [":" typeAnnotation] "=" expressionWithoutModifier
func (p *Parser) constantDeclaration() ast.ExpressionNode {
	constTok := p.advance()
	var init ast.ExpressionNode
	var typ ast.TypeNode

	constName, ok := p.matchOk(token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT)
	if !ok {
		p.errorExpected("an uppercase identifier as the name of the declared constant")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Span(),
			tok,
		)
	}
	lastSpan := constName.Span()

	if p.match(token.COLON) {
		typ = p.typeAnnotation()
		lastSpan = typ.Span()
	}

	if p.match(token.EQUAL_OP) {
		p.swallowNewlines()
		init = p.expressionWithoutModifier()
		lastSpan = init.Span()
	} else {
		p.errorMessageSpan("constants must be initialised", constTok.Span().Join(lastSpan))
	}

	return ast.NewConstantDeclarationNode(
		constTok.Span().Join(lastSpan),
		constName,
		typ,
		init,
	)
}

// typeAnnotation = unionType
func (p *Parser) typeAnnotation() ast.TypeNode {
	return p.unionType()
}

// unionType = intersectionType | unionType "|" intersectionType
func (p *Parser) unionType() ast.TypeNode {
	return p.binaryTypeExpression(p.intersectionType, token.OR)
}

// intersectionType = nilableType | intersectionType "&" nilableType
func (p *Parser) intersectionType() ast.TypeNode {
	return p.binaryTypeExpression(p.nilableType, token.AND)
}

// nilableType = primaryType | primaryType "?"
func (p *Parser) nilableType() ast.TypeNode {
	primType := p.primaryType()

	if questTok, ok := p.matchOk(token.QUESTION); ok {
		return ast.NewNilableTypeNode(
			primType.Span().Join(questTok.Span()),
			primType,
		)
	}

	return primType
}

// primaryType = namedType | "(" typeAnnotation ")"
func (p *Parser) primaryType() ast.TypeNode {
	if p.match(token.LPAREN) {
		t := p.typeAnnotation()
		p.consume(token.RPAREN)
		return t
	}

	return p.namedType()
}

// namedType = genericConstant
func (p *Parser) namedType() ast.TypeNode {
	return p.genericConstant()
}

// genericConstant = strictConstantLookup | strictConstantLookup "[" [genericConstantList] "]"
func (p *Parser) genericConstant() ast.ComplexConstantNode {
	constant := p.strictConstantLookup()
	if !p.match(token.LBRACKET) {
		return constant
	}

	if p.match(token.RBRACKET) {
		p.errorExpected("a constant")
		return constant
	}

	constList := p.genericConstantList(token.RBRACKET)
	rbracket, ok := p.consume(token.RBRACKET)
	if !ok {
		return ast.NewInvalidNode(rbracket.Span(), rbracket)
	}

	return ast.NewGenericConstantNode(
		constant.Span().Join(rbracket.Span()),
		constant,
		constList,
	)
}

// throwExpression = "throw" [expressionWithoutModifier]
func (p *Parser) throwExpression() *ast.ThrowExpressionNode {
	throwTok := p.advance()
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() {
		return ast.NewThrowExpressionNode(
			throwTok.Span(),
			nil,
		)
	}

	expr := p.expressionWithoutModifier()

	return ast.NewThrowExpressionNode(
		throwTok.Span().Join(expr.Span()),
		expr,
	)
}

// continueExpression = "continue" [expressionWithoutModifier]
func (p *Parser) continueExpression() *ast.ContinueExpressionNode {
	continueTok := p.advance()
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() {
		return ast.NewContinueExpressionNode(
			continueTok.Span(),
			nil,
		)
	}

	expr := p.expressionWithoutModifier()

	return ast.NewContinueExpressionNode(
		continueTok.Span().Join(expr.Span()),
		expr,
	)
}

// returnExpression = "return" [expressionWithoutModifier]
func (p *Parser) returnExpression() *ast.ReturnExpressionNode {
	returnTok := p.advance()
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() {
		return ast.NewReturnExpressionNode(
			returnTok.Span(),
			nil,
		)
	}

	expr := p.expressionWithoutModifier()

	return ast.NewReturnExpressionNode(
		returnTok.Span().Join(expr.Span()),
		expr,
	)
}

// loopExpression = "loop" ((SEPARATOR [statements] "end") | expressionWithoutModifier)
func (p *Parser) loopExpression() *ast.LoopExpressionNode {
	loopTok := p.advance()
	var span *position.Span

	lastSpan, thenBody, multiline := p.statementBlock(token.END)
	if lastSpan != nil {
		span = loopTok.Span().Join(lastSpan)
	} else {
		span = loopTok.Span()
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewLoopExpressionNode(
		span,
		thenBody,
	)
}

// forExpression = ("for" identifierList "in" expressionWithoutModifier) |
// ("for" [expressionWithoutModifier] ";" [expressionWithoutModifier] ";" [expressionWithoutModifier])
// ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) forExpression() ast.ExpressionNode {
	forTok := p.advance()

	p.swallowNewlines()
	if p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) && p.acceptNext(token.COMMA, token.IN) {
		return p.forInExpression(forTok)
	}

	// numeric for
	return p.numericForExpression(forTok)
}

func (p *Parser) numericForExpression(forTok *token.Token) ast.ExpressionNode {
	span := forTok.Span()
	var init, cond, incr ast.ExpressionNode
	if !p.accept(token.SEMICOLON) {
		init = p.expressionWithoutModifier()
	}
	semiTok, ok := p.consume(token.SEMICOLON)
	if !ok {
		return ast.NewInvalidNode(semiTok.Span(), semiTok)
	}

	if !p.accept(token.SEMICOLON) {
		cond = p.expressionWithoutModifier()
	}
	semiTok, ok = p.consume(token.SEMICOLON)
	if !ok {
		return ast.NewInvalidNode(semiTok.Span(), semiTok)
	}
	span = span.Join(semiTok.Span())

	if !p.accept(token.NEWLINE, token.SEMICOLON, token.THEN) {
		incr = p.expressionWithoutModifier()
	}

	lastSpan, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastSpan != nil {
		span = span.Join(lastSpan)
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewNumericForExpressionNode(
		span,
		init,
		cond,
		incr,
		thenBody,
	)
}

func (p *Parser) forInExpression(forTok *token.Token) ast.ExpressionNode {
	loopParameters := p.identifierList(token.IN)

	inTok, ok := p.consume(token.IN)
	if !ok {
		return ast.NewInvalidNode(inTok.Span(), inTok)
	}
	p.swallowNewlines()
	inExpr := p.expressionWithoutModifier()
	var span *position.Span

	lastSpan, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastSpan != nil {
		span = forTok.Span().Join(lastSpan)
	} else {
		span = forTok.Span()
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewForInExpressionNode(
		span,
		loopParameters,
		inExpr,
		thenBody,
	)
}

// whileExpression = "while" expressionWithoutModifier ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) whileExpression() *ast.WhileExpressionNode {
	whileTok := p.advance()
	cond := p.expressionWithoutModifier()
	var span *position.Span

	lastSpan, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastSpan != nil {
		span = whileTok.Span().Join(lastSpan)
	} else {
		span = whileTok.Span()
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewWhileExpressionNode(
		span,
		cond,
		thenBody,
	)
}

// untilExpression = "until" expressionWithoutModifier ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) untilExpression() *ast.UntilExpressionNode {
	untilTok := p.advance()
	cond := p.expressionWithoutModifier()
	var span *position.Span

	lastSpan, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastSpan != nil {
		span = untilTok.Span().Join(lastSpan)
	} else {
		span = untilTok.Span()
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewUntilExpressionNode(
		span,
		cond,
		thenBody,
	)
}

// unlessExpression = "unless" expressionWithoutModifier ((SEPARATOR [statements]) | ("then" expressionWithoutModifier))
// ["else" ((SEPARATOR [statements]) | expressionWithoutModifier)]
// "end"
func (p *Parser) unlessExpression() *ast.UnlessExpressionNode {
	unlessTok := p.advance()
	cond := p.expressionWithoutModifier()
	var span *position.Span

	lastSpan, thenBody, multiline := p.statementBlockWithThen(token.END, token.ELSE)
	if lastSpan != nil {
		span = unlessTok.Span().Join(lastSpan)
	} else {
		span = unlessTok.Span()
	}

	unlessExpr := ast.NewUnlessExpressionNode(
		span,
		cond,
		thenBody,
		nil,
	)
	currentExpr := unlessExpr

	if p.lookahead.IsStatementSeparator() && p.nextLookahead.Type == token.ELSE {
		p.advance()
		p.advance()
		lastSpan, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastSpan != nil {
			currentExpr.SetSpan(currentExpr.Span().Join(lastSpan))
		}
	} else if p.match(token.ELSE) {
		lastSpan, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastSpan != nil {
			currentExpr.SetSpan(currentExpr.Span().Join(lastSpan))
		}
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			currentExpr.SetSpan(currentExpr.Span().Join(endTok.Span()))
		}
	}
	unlessExpr.SetSpan(unlessExpr.Span().Join(currentExpr.Span()))

	return unlessExpr
}

// doExpression = "do" [SEPARATOR] [statements] "end"
func (p *Parser) doExpression() *ast.DoExpressionNode {
	doTok := p.advance()
	lastSpan, body, _ := p.statementBlock(token.END)

	var span *position.Span
	if lastSpan != nil {
		span = doTok.Span().Join(lastSpan)
	} else {
		span = doTok.Span()
	}

	doExpr := ast.NewDoExpressionNode(
		span,
		body,
	)

	if len(body) == 0 {
		p.indentedSection = true
	}
	endTok, ok := p.consume(token.END)
	if len(body) == 0 {
		p.indentedSection = false
	}
	if ok {
		doExpr.SetSpan(doExpr.Span().Join(endTok.Span()))
	}

	return doExpr
}

// ifExpression = "if" expressionWithoutModifier ((SEPARATOR [statements]) | ("then" expressionWithoutModifier))
// ("elsif" expressionWithoutModifier ((SEPARATOR [statements]) | ("then" expressionWithoutModifier)) )*
// ["else" ((SEPARATOR [statements]) | expressionWithoutModifier)]
// "end"
func (p *Parser) ifExpression() *ast.IfExpressionNode {
	ifTok := p.advance()
	cond := p.expressionWithoutModifier()
	var span *position.Span

	lastSpan, thenBody, multiline := p.statementBlockWithThen(token.END, token.ELSE, token.ELSIF)
	if lastSpan != nil {
		span = ifTok.Span().Join(lastSpan)
	} else {
		span = ifTok.Span()
	}

	ifExpr := ast.NewIfExpressionNode(
		span,
		cond,
		thenBody,
		nil,
	)
	currentExpr := ifExpr

	for {
		var elsifTok *token.Token

		if p.lookahead.Type == token.ELSIF {
			elsifTok = p.advance()
		} else if p.lookahead.IsStatementSeparator() && p.nextLookahead.Type == token.ELSIF {
			p.advance()
			elsifTok = p.advance()
		} else {
			break
		}
		cond = p.expressionWithoutModifier()

		lastSpan, thenBody, multiline = p.statementBlockWithThen(token.END, token.ELSE, token.ELSIF)
		if lastSpan != nil {
			span = elsifTok.Span().Join(lastSpan)
		} else {
			span = elsifTok.Span()
		}

		elsifExpr := ast.NewIfExpressionNode(
			span,
			cond,
			thenBody,
			nil,
		)

		currentExpr.ElseBody = []ast.StatementNode{
			ast.NewExpressionStatementNode(
				elsifExpr.Span(),
				elsifExpr,
			),
		}
		currentExpr = elsifExpr
	}

	if p.lookahead.IsStatementSeparator() && p.nextLookahead.Type == token.ELSE {
		p.advance()
		p.advance()
		lastSpan, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastSpan != nil {
			currentExpr.SetSpan(currentExpr.Span().Join(lastSpan))
		}
	} else if p.match(token.ELSE) {
		lastSpan, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastSpan != nil {
			currentExpr.SetSpan(currentExpr.Span().Join(lastSpan))
		}
	}

	if multiline {
		if len(thenBody) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(thenBody) == 0 {
			p.indentedSection = false
		}
		if ok {
			currentExpr.SetSpan(currentExpr.Span().Join(endTok.Span()))
		}
	}
	ifExpr.SetSpan(ifExpr.Span().Join(currentExpr.Span()))

	return ifExpr
}

// symbolLiteral = ":" (identifier | constant | rawStringLiteral)
func (p *Parser) symbolLiteral() ast.ExpressionNode {
	symbolBegTok := p.advance()
	if p.lookahead.IsValidSimpleSymbolContent() {
		contTok := p.advance()
		return ast.NewSimpleSymbolLiteralNode(
			symbolBegTok.Span().Join(contTok.Span()),
			contTok.StringValue(),
		)
	}

	if !p.accept(token.STRING_BEG) {
		p.errorExpected("an identifier, overridable operator or string literal")
		p.mode = panicMode
		tok := p.advance()
		return ast.NewInvalidNode(
			symbolBegTok.Span().Join(tok.Span()),
			tok,
		)
	}

	str := p.stringLiteral()
	switch s := str.(type) {
	case *ast.DoubleQuotedStringLiteralNode:
		return ast.NewSimpleSymbolLiteralNode(
			symbolBegTok.Span().Join(s.Span()),
			s.Value,
		)
	case *ast.InterpolatedStringLiteralNode:
		return ast.NewInterpolatedSymbolLiteral(
			symbolBegTok.Span().Join(s.Span()),
			s,
		)
	default:
		return s
	}
}

// charLiteral = CHAR_LITERAL
func (p *Parser) charLiteral() *ast.CharLiteralNode {
	tok := p.advance()
	char, _ := utf8.DecodeRuneInString(tok.Value)
	return ast.NewCharLiteralNode(
		tok.Span(),
		char,
	)
}

// charLiteral = RAW_CHAR_LITERAL
func (p *Parser) rawCharLiteral() *ast.RawCharLiteralNode {
	tok := p.advance()
	char, _ := utf8.DecodeRuneInString(tok.Value)
	return ast.NewRawCharLiteralNode(
		tok.Span(),
		char,
	)
}

// rawStringLiteral = RAW_STRING
func (p *Parser) rawStringLiteral() *ast.RawStringLiteralNode {
	tok := p.advance()
	return ast.NewRawStringLiteralNode(
		tok.Span(),
		tok.Value,
	)
}

// wordCollectionElement = RAW_STRING
func (p *Parser) wordCollectionElement() ast.WordCollectionContentNode {
	tok, ok := p.consume(token.RAW_STRING)
	if !ok {
		return ast.NewInvalidNode(tok.Span(), tok)
	}
	return ast.NewRawStringLiteralNode(
		tok.Span(),
		tok.Value,
	)
}

// symbolCollectionElement = RAW_STRING
func (p *Parser) symbolCollectionElement() ast.SymbolCollectionContentNode {
	tok, ok := p.consume(token.RAW_STRING)
	if !ok {
		return ast.NewInvalidNode(tok.Span(), tok)
	}
	return ast.NewSimpleSymbolLiteralNode(
		tok.Span(),
		tok.Value,
	)
}

// intCollectionElement = INT
func (p *Parser) intCollectionElement() ast.IntCollectionContentNode {
	tok, ok := p.consume(token.INT)
	if !ok {
		return ast.NewInvalidNode(tok.Span(), tok)
	}
	return ast.NewIntLiteralNode(tok.Span(), tok.Value)
}

// stringLiteral = "\"" (STRING_CONTENT | "${" expressionWithoutModifier "}")* "\""
func (p *Parser) stringLiteral() ast.StringLiteralNode {
	quoteBeg := p.advance() // consume the opening quote
	var quoteEnd *token.Token

	var strContent []ast.StringLiteralContentNode
	for {
		if tok, ok := p.matchOk(token.STRING_CONTENT); ok {
			strContent = append(strContent, ast.NewStringLiteralContentSectionNode(
				tok.Span(),
				tok.Value,
			))
			continue
		}

		if beg, ok := p.matchOk(token.STRING_INTERP_BEG); ok {
			expr := p.expressionWithoutModifier()
			end, _ := p.consume(token.STRING_INTERP_END)
			strContent = append(strContent, ast.NewStringInterpolationNode(
				beg.Span().Join(end.Span()),
				expr,
			))
			continue
		}

		tok, ok := p.consume(token.STRING_END)
		quoteEnd = tok
		if tok.Type == token.END_OF_FILE {
			break
		}
		if !ok {
			strContent = append(strContent, ast.NewInvalidNode(
				tok.Span(),
				tok,
			))
			continue
		}
		break
	}
	if len(strContent) == 0 {
		return ast.NewDoubleQuotedStringLiteralNode(
			quoteBeg.Span().Join(quoteEnd.Span()),
			"",
		)
	}
	strVal, ok := strContent[0].(*ast.StringLiteralContentSectionNode)
	if len(strContent) == 1 && ok {
		return ast.NewDoubleQuotedStringLiteralNode(
			quoteBeg.Span().Join(quoteEnd.Span()),
			strVal.Value,
		)
	}

	return ast.NewInterpolatedStringLiteralNode(
		quoteBeg.Span().Join(quoteEnd.Span()),
		strContent,
	)
}

// closureAfterArrow = "->" (expressionWithoutModifier | SEPARATOR [statements] "end" | "{" [statements] "}")
func (p *Parser) closureAfterArrow(firstSpan *position.Span, params []ast.ParameterNode, returnType ast.TypeNode, throwType ast.TypeNode) ast.ExpressionNode {
	var span *position.Span
	arrowTok, ok := p.consume(token.THIN_ARROW)
	if !ok {
		return ast.NewInvalidNode(
			arrowTok.Span(),
			arrowTok,
		)
	}
	if firstSpan == nil {
		firstSpan = arrowTok.Span()
	}

	// Body with curly braces
	if p.match(token.LBRACE) {
		p.swallowNewlines()
		body := p.statements(token.RBRACE)
		if tok, ok := p.consume(token.RBRACE); ok {
			span = firstSpan.Join(tok.Span())
		} else {
			span = firstSpan
		}
		return ast.NewClosureLiteralNode(
			span,
			params,
			returnType,
			throwType,
			body,
		)
	}

	lastSpan, body, multiline := p.statementBlock(token.END)
	if lastSpan != nil {
		span = firstSpan.Join(lastSpan)
	} else {
		span = firstSpan
	}

	if multiline {
		if len(body) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(body) == 0 {
			p.indentedSection = false
		}
		if ok {
			span = span.Join(endTok.Span())
		}
	}

	return ast.NewClosureLiteralNode(
		span,
		params,
		returnType,
		throwType,
		body,
	)
}

// closureExpression = (("|" formalParameterList "|") | "||") [: typeAnnotation] ["!" typeAnnotation] closureAfterArrow
func (p *Parser) closureExpression() ast.ExpressionNode {
	var params []ast.ParameterNode
	var firstSpan *position.Span
	var returnType ast.TypeNode
	var throwType ast.TypeNode

	if p.accept(token.OR) {
		firstSpan = p.advance().Span()
		if !p.accept(token.OR) {
			p.mode = withoutBitwiseOrMode
			params = p.formalParameterList(token.OR)
			p.mode = normalMode
		}
		if tok, ok := p.consume(token.OR); !ok {
			return ast.NewInvalidNode(
				tok.Span(),
				tok,
			)
		}
	} else {
		orOr, ok := p.consume(token.OR_OR)
		firstSpan = orOr.Span()
		if !ok {
			return ast.NewInvalidNode(
				orOr.Span(),
				orOr,
			)
		}
	}

	// return type
	if p.match(token.COLON) {
		returnType = p.typeAnnotation()
	}

	// throw type
	if p.match(token.BANG) {
		throwType = p.typeAnnotation()
	}

	return p.closureAfterArrow(firstSpan, params, returnType, throwType)
}

// identifierOrClosure = identifier | identifier closureAfterArrow
func (p *Parser) identifierOrClosure() ast.ExpressionNode {
	if p.nextLookahead.Type == token.THIN_ARROW {
		ident := p.advance()
		return p.closureAfterArrow(
			ident.Span(),
			[]ast.ParameterNode{
				ast.NewFormalParameterNode(
					ident.Span(),
					ident.Value,
					nil,
					nil,
					ast.NormalParameterKind,
				),
			},
			nil,
			nil,
		)
	}

	return p.identifier()
}

// docComment = DOC_COMMENT expressionWithModifier
func (p *Parser) docComment() *ast.DocCommentNode {
	docComment := p.advance()
	if p.lookahead.Type == token.DOC_COMMENT {
		p.errorMessage("doc comments can't document one another")
	}
	p.swallowNewlines()
	expr := p.expressionWithModifier()

	return ast.NewDocCommentNode(
		docComment.Span().Join(expr.Span()),
		docComment.Value,
		expr,
	)
}

// identifier = PUBLIC_IDENTIFIER | PRIVATE_IDENTIFIER
func (p *Parser) identifier() ast.IdentifierNode {
	if p.accept(token.PUBLIC_IDENTIFIER) {
		ident := p.advance()
		return ast.NewPublicIdentifierNode(
			ident.Span(),
			ident.Value,
		)
	}
	if p.accept(token.PRIVATE_IDENTIFIER) {
		ident := p.advance()
		return ast.NewPrivateIdentifierNode(
			ident.Span(),
			ident.Value,
		)
	}

	p.errorExpected("an identifier")
	errTok := p.advance()
	return ast.NewInvalidNode(errTok.Span(), errTok)
}
