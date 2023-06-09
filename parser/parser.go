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
	"github.com/elk-language/elk/token"
)

// Parsing mode.
type mode uint8

const (
	normalMode           mode = iota // regular parsing mode
	panicMode                        // triggered after encountering a syntax error, changes to `normalMode` after synchronisation
	withoutBitwiseOrMode             // disables bitwise OR `|` from the grammar
)

// Holds the current state of the parsing process.
type Parser struct {
	sourceName    string       // Path to the source file or some name.
	source        []byte       // Elk source code
	lexer         *lexer.Lexer // lexer which outputs a stream of tokens
	lookahead     *token.Token // next token used for predicting productions
	nextLookahead *token.Token // second next token used for predicting productions
	errors        position.ErrorList
	mode          mode
}

// Instantiate a new parser.
func new(sourceName string, source []byte) *Parser {
	return &Parser{
		sourceName: sourceName,
		source:     source,
		lexer:      lexer.NewWithName(sourceName, source),
		mode:       normalMode,
	}
}

// Parse the given source code and return an Abstract Syntax Tree.
// Main entry point to the parser.
func Parse(sourceName string, source []byte) (*ast.ProgramNode, position.ErrorList) {
	return new(sourceName, source).parse()
}

// Start the parsing process from the top.
func (p *Parser) parse() (*ast.ProgramNode, position.ErrorList) {
	p.advance() // populate nextLookahead
	p.advance() // populate lookahead
	return p.program(), p.errors
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
	p.errorMessagePos(message, p.lookahead.Position)
}

// Same as [errorMessage] but let's you pass a Position.
func (p *Parser) errorMessagePos(message string, pos *position.Position) {
	if p.mode == panicMode {
		return
	}

	p.errors.Add(
		message,
		position.NewLocationWithPosition(p.sourceName, pos),
	)
}

// Add the content of an error token to the syntax error list.
func (p *Parser) errorToken(err *token.Token) {
	p.errorMessagePos(err.Value, err.Position)
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
		p.mode = panicMode
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
func (p *Parser) statementBlock(stopTokens ...token.Type) (*position.Position, []ast.StatementNode, bool) {
	var thenBody []ast.StatementNode
	var lastPos *position.Position
	var multiline bool

	if !p.lookahead.IsStatementSeparator() {
		expr := p.expressionWithoutModifier()
		thenBody = append(thenBody, ast.NewExpressionStatementNode(
			expr.Pos(),
			expr,
		))
		lastPos = expr.Pos()
	} else {
		multiline = true
		p.advance()

		if p.accept(token.END) {
			lastPos = p.lookahead.Position
		} else if !containsToken(stopTokens, p.lookahead.Type) {
			thenBody = p.statements(stopTokens...)
			if len(thenBody) > 0 {
				lastPos = thenBody[len(thenBody)-1].Pos()
			}
		}
	}

	return lastPos, thenBody, multiline
}

// statementProduction = subProduction [SEPARATOR]
func statementProduction[Expression, Statement ast.Node](p *Parser, constructor statementConstructor[Expression, Statement], expressionProduction func() Expression, separators ...token.Type) Statement {
	expr := expressionProduction()
	var sep *token.Token
	if p.lookahead.IsStatementSeparator() || p.lookahead.Type == token.END_OF_FILE {
		sep = p.advance()
		return constructor(
			expr.Pos().Join(sep.Pos()),
			expr,
		)
	}
	for _, sepType := range separators {
		if p.lookahead.Type == sepType {
			return constructor(
				expr.Pos(),
				expr,
			)
		}
	}
	if p.match(token.ERROR) {
		if p.synchronise() {
			p.advance()
		}
		return constructor(
			expr.Pos(),
			expr,
		)
	}

	p.errorExpected(statementSeparatorMessage)
	if p.synchronise() {
		p.advance()
	}

	return constructor(
		expr.Pos(),
		expr,
	)
}

type statementsProduction[Statement ast.Node] func(...token.Type) []Statement

// Represents an AST Node constructor function for a new ast.StatementNode
type statementConstructor[Expression, Statement ast.Node] func(*position.Position, Expression) Statement

// Consume a block of statements, like in `if`, `elsif` or `while` expressions,
// that terminates with `end` or can be single-line when it begins with `then`.
func genericStatementBlockWithThen[Expression, Statement ast.Node](
	p *Parser,
	statementsProduction statementsProduction[Statement],
	expressionProduction func() Expression,
	statementConstructor statementConstructor[Expression, Statement],
	stopTokens ...token.Type,
) (*position.Position, []Statement, bool) {
	var thenBody []Statement
	var lastPos *position.Position
	var multiline bool

	if p.lookahead.Type == token.THEN {
		p.advance()
		expr := expressionProduction()
		thenBody = append(thenBody, statementConstructor(
			expr.Pos(),
			expr,
		))
		lastPos = expr.Pos()
	} else {
		multiline = true
		if p.lookahead.IsStatementSeparator() {
			p.advance()
		} else {
			p.errorExpected(statementSeparatorMessage)
		}

		if p.accept(token.END) {
			lastPos = p.lookahead.Position
		} else if !containsToken(stopTokens, p.lookahead.Type) {
			thenBody = statementsProduction(stopTokens...)
			if len(thenBody) > 0 {
				lastPos = thenBody[len(thenBody)-1].Pos()
			}
		}
	}

	return lastPos, thenBody, multiline
}

// Consume a block of statements, like in `if`, `elsif` or `while` expressions,
// that terminates with `end` or can be single-line when it begins with `then`.
func (p *Parser) statementBlockWithThen(stopTokens ...token.Type) (*position.Position, []ast.StatementNode, bool) {
	return genericStatementBlockWithThen(p, p.statements, p.expressionWithoutModifier, ast.NewExpressionStatementNodeI, stopTokens...)
}

// Consume a block of statements, like in `if`, `elsif` or `while` expressions,
// that terminates with `end` or can be single-line when it begins with `then`.
func (p *Parser) structBodyStatementBlockWithThen(stopTokens ...token.Type) (*position.Position, []ast.StructBodyStatementNode, bool) {
	return genericStatementBlockWithThen(p, p.structBodyStatements, p.formalParameter, ast.NewParameterStatementNodeI, stopTokens...)
}

// Represents an AST Node constructor function for binary operators
type binaryConstructor[Element ast.Node] func(*position.Position, *token.Token, Element, Element) Element

// binaryProduction = subProduction | binaryProduction operators subProduction
func binaryProduction[Element ast.Node](p *Parser, constructor binaryConstructor[Element], subProduction func() Element, operators ...token.Type) Element {
	left := subProduction()

	for {
		operator, ok := p.matchOk(operators...)
		if !ok {
			break
		}
		p.swallowNewlines()
		right := subProduction()
		left = constructor(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
}

// Represents an AST Node constructor function for an `include`- like expression
// eg. `include`, `extend`, `enhance`
type includelikeConstrutor[T ast.Node] func(*position.Position, []ast.ComplexConstantNode) T

// includelikeExpression = keyword genericConstantList
func includelikeExpression[T ast.Node](p *Parser, constructor includelikeConstrutor[T]) T {
	keyword := p.advance()
	consts := p.genericConstantList()
	pos := position.JoinLastElement(keyword.Position, consts)

	return constructor(
		pos,
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
	return ast.NewProgramNode(
		position.New(0, len(p.source), 1, 1),
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
	return ast.NewEmptyStatementNode(sepTok.Position)
}

const statementSeparatorMessage = "a statement separator `\\n`, `;`"

// expressionStatement = expressionWithModifier [SEPARATOR]
func (p *Parser) expressionStatement(separators ...token.Type) *ast.ExpressionStatementNode {
	return statementProduction(p, ast.NewExpressionStatementNode, p.expressionWithModifier, separators...)
}

// expressionWithModifier = modifierExpression
func (p *Parser) expressionWithModifier() ast.ExpressionNode {
	asgmt := p.modifierExpression()
	if p.mode == panicMode {
		p.synchronise()
	}
	return asgmt
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
// expressionWithoutModifier "for" loopParameterList "in" expressionWithoutModifier
func (p *Parser) modifierExpression() ast.ExpressionNode {
	left := p.expressionWithoutModifier()

	switch p.lookahead.Type {
	case token.UNLESS, token.WHILE, token.UNTIL:
		mod := p.advance()
		p.swallowNewlines()
		right := p.expressionWithoutModifier()
		return ast.NewModifierNode(
			left.Pos().Join(right.Pos()),
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
				left.Pos().Join(elseExpr.Pos()),
				left,
				cond,
				elseExpr,
			)
		}
		return ast.NewModifierNode(
			left.Pos().Join(cond.Pos()),
			ifTok,
			left,
			cond,
		)
	case token.FOR:
		p.advance()
		p.swallowNewlines()
		params := p.loopParameterList(token.IN)
		inTok, ok := p.consume(token.IN)
		if !ok {
			return ast.NewInvalidNode(inTok.Position, inTok)
		}
		p.swallowNewlines()
		inExpr := p.expressionWithoutModifier()
		return ast.NewModifierForInNode(
			left.Pos().Join(inExpr.Pos()),
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
	if p.lookahead.Type == token.COLON_EQUAL {
		if !ast.IsValidDeclarationTarget(left) {
			p.errorMessagePos(
				fmt.Sprintf("invalid `%s` declaration target", p.lookahead.Type.String()),
				left.Pos(),
			)
		}
	}

	if !p.lookahead.IsAssignmentOperator() {
		return left
	}

	if ast.IsConstant(left) {
		p.errorMessagePos(
			"constants can't be assigned, maybe you meant to declare it with `:=`",
			left.Pos(),
		)
	} else if !ast.IsValidAssignmentTarget(left) {
		p.errorMessagePos(
			fmt.Sprintf("invalid `%s` assignment target", p.lookahead.Type.String()),
			left.Pos(),
		)
	}

	operator := p.advance()
	p.swallowNewlines()
	right := p.assignmentExpression()

	return ast.NewAssignmentExpressionNode(
		left.Pos().Join(right.Pos()),
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
	var pos *position.Position

	if starTok, ok := p.matchOk(token.STAR); ok {
		kind = ast.PositionalRestParameterKind
		pos = starTok.Position
	} else if starStarTok, ok := p.matchOk(token.STAR_STAR); ok {
		kind = ast.NamedRestParameterKind
		pos = starStarTok.Position
	}

	switch p.lookahead.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		paramName = p.advance()
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		p.errorExpected("a lowercased identifier as the name of the declared formalParameter")
		paramName = p.advance()
	default:
		p.errorExpected("an identifier as the name of the declared formalParameter")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Position,
			tok,
		)
	}
	pos = pos.Join(paramName.Position)

	if p.match(token.COLON) {
		typ = p.intersectionType()
		pos = pos.Join(typ.Pos())
	}

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		pos = pos.Join(init.Pos())
	}

	return ast.NewFormalParameterNode(
		pos,
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
	var pos *position.Position

	if starTok, ok := p.matchOk(token.STAR); ok {
		kind = ast.PositionalRestParameterKind
		pos = starTok.Position
	} else if starStarTok, ok := p.matchOk(token.STAR_STAR); ok {
		kind = ast.NamedRestParameterKind
		pos = starStarTok.Position
	}

	switch p.lookahead.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		paramName = p.advance()
	case token.INSTANCE_VARIABLE:
		paramName = p.advance()
		setIvar = true
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		p.errorExpected("a lowercased identifier as the name of the declared formalParameter")
		paramName = p.advance()
	default:
		p.errorExpected("an identifier as the name of the declared formalParameter")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Position,
			tok,
		)
	}
	pos = pos.Join(paramName.Position)

	if p.match(token.COLON) {
		typ = p.intersectionType()
		pos = pos.Join(typ.Pos())
	}

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		pos = pos.Join(init.Pos())
	}

	return ast.NewMethodParameterNode(
		pos,
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
			p.errorMessagePos("there should be only a single positional rest parameter", element.Pos())
			continue
		}

		if posRest {
			if namedRestSeen {
				p.errorMessagePos("named rest parameters should appear last", element.Pos())
			}
			posRestSeen = true
			continue
		}

		namedRest := ast.IsNamedRestParam(element)
		if namedRest && namedRestSeen {
			p.errorMessagePos("there should be only a single named rest parameter", element.Pos())
			continue
		}

		if namedRestSeen {
			p.errorMessagePos("named rest parameters should appear last", element.Pos())
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
			p.errorMessagePos("required parameters can't appear after optional parameters", element.Pos())
		} else if opt && !optionalSeen {
			optionalSeen = true
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
			tok.Position,
			tok,
		)
	}
	lastPos := paramName.Position

	if questionTok, ok := p.matchOk(token.QUESTION); ok {
		opt = true
		lastPos = questionTok.Pos()
	}

	if p.match(token.COLON) {
		typ = p.intersectionType()
		lastPos = typ.Pos()
	}

	return ast.NewSignatureParameterNode(
		paramName.Position.Join(lastPos.Pos()),
		paramName.Value,
		typ,
		opt,
	)
}

// signatureParameterList = signatureParameter ("," signatureParameter)*
func (p *Parser) signatureParameterList(stopTokens ...token.Type) []ast.ParameterNode {
	return p.parameterList(p.signatureParameter, stopTokens...)
}

// loopParameter = identifier [":" typeAnnotation]
func (p *Parser) loopParameter() ast.ParameterNode {
	var typ ast.TypeNode

	var paramName *token.Token

	switch p.lookahead.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		paramName = p.advance()
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		p.errorExpected("a lowercase identifier as the name of the declared loopParameter")
		paramName = p.advance()
	default:
		p.errorExpected("an identifier as the name of the declared loopParameter")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Position,
			tok,
		)
	}
	lastPos := paramName.Position

	if p.match(token.COLON) {
		typ = p.intersectionType()
		lastPos = typ.Pos()
	}

	return ast.NewLoopParameterNode(
		paramName.Position.Join(lastPos.Pos()),
		paramName.Value,
		typ,
	)
}

// loopParameterList = loopParameter ("," loopParameter)*
func (p *Parser) loopParameterList(stopTokens ...token.Type) []ast.ParameterNode {
	return commaSeparatedList(p, p.loopParameter, stopTokens...)
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
		right := p.comparisonExpression()

		left = ast.NewBinaryExpressionNode(
			left.Pos().Join(right.Pos()),
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
		right := p.bitwiseShiftExpression()

		left = ast.NewBinaryExpressionNode(
			left.Pos().Join(right.Pos()),
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

// multiplicativeExpression = unaryExpression | multiplicativeExpression ("*" | "/") unaryExpression
func (p *Parser) multiplicativeExpression() ast.ExpressionNode {
	return p.binaryExpression(p.unaryExpression, token.STAR, token.SLASH)
}

// unaryExpression = powerExpression | ("!" | "-" | "+" | "~") unaryExpression
func (p *Parser) unaryExpression() ast.ExpressionNode {
	if operator, ok := p.matchOk(token.BANG, token.MINUS, token.PLUS, token.TILDE); ok {
		p.swallowNewlines()
		right := p.unaryExpression()
		return ast.NewUnaryExpressionNode(
			operator.Pos().Join(right.Pos()),
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
	right := p.powerExpression()

	return ast.NewBinaryExpressionNode(
		left.Pos().Join(right.Pos()),
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
			errTok.Position,
			errTok,
		)
	}
	colon, ok := p.consume(token.COLON)
	if !ok {
		return ast.NewInvalidNode(
			colon.Position,
			colon,
		)
	}
	val := p.expressionWithoutModifier()

	return ast.NewNamedCallArgumentNode(
		ident.Pos().Join(val.Pos()),
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
		lastPos, posArgs, namedArgs, errToken := p.callArgumentList()
		if errToken != nil {
			return ast.NewInvalidNode(
				errToken.Position,
				errToken,
			)
		}
		if lastPos == nil {
			p.errorExpected("method arguments")
			errToken = p.advance()
			return ast.NewInvalidNode(
				errToken.Position,
				errToken,
			)
		}

		receiver = ast.NewFunctionCallNode(
			methodName.Position.Join(lastPos),
			methodName.Value,
			posArgs,
			namedArgs,
		)
	}

	// method call
	if receiver == nil {
		receiver = p.rangeLiteral()
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
			p.mode = panicMode
			errTok := p.advance()
			return ast.NewInvalidNode(
				errTok.Position,
				errTok,
			)
		}

		methodNameTok := p.advance()
		methodName := methodNameTok.StringValue()

		lastPos, posArgs, namedArgs, errToken := p.callArgumentList()
		if errToken != nil {
			return ast.NewInvalidNode(
				errToken.Position,
				errToken,
			)
		}
		if lastPos == nil {
			lastPos = methodNameTok.Position
		}

		receiver = ast.NewMethodCallNode(
			receiver.Pos().Join(lastPos),
			receiver,
			opToken.Type == token.QUESTION_DOT,
			methodName,
			posArgs,
			namedArgs,
		)
	}
}

// rangeLiteral = constructorCall (".." | "...") [constructorCall]
func (p *Parser) rangeLiteral() ast.ExpressionNode {
	left := p.constructorCall()
	op, ok := p.matchOk(token.RANGE_OP, token.EXCLUSIVE_RANGE_OP)
	if !ok {
		return left
	}

	if !p.lookahead.IsValidAsEndInRangeLiteral() {
		return ast.NewRangeLiteralNode(
			left.Pos().Join(op.Position),
			op.Type == token.EXCLUSIVE_RANGE_OP,
			left,
			nil,
		)
	}

	right := p.constructorCall()

	return ast.NewRangeLiteralNode(
		left.Pos().Join(right.Pos()),
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
		op.Position.Join(right.Pos()),
		op.Type == token.EXCLUSIVE_RANGE_OP,
		nil,
		right,
	)
}

// callArgumentListInternal = (positionalArgumentList | namedArgumentList | positionalArgumentList "," namedArgumentList)
// callArgumentList = "(" callArgumentList ")" | callArgumentList
func (p *Parser) callArgumentList() (*position.Position, []ast.ExpressionNode, []ast.NamedArgumentNode, *token.Token) {
	if p.match(token.LPAREN) {
		p.swallowNewlines()
		if rparen, ok := p.matchOk(token.RPAREN); ok {
			return rparen.Position,
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

		return rparen.Position,
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
	pos := position.OfLastElement(posArgs)
	var namedArgs []ast.NamedArgumentNode
	if len(posArgs) == 0 || len(posArgs) > 0 && commaConsumed {
		namedArgs = p.namedArgumentList()
		pos = position.OfLastElement(namedArgs)
	}

	return pos,
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

	lastPos, posArgs, namedArgs, errToken := p.callArgumentList()
	if errToken != nil {
		return ast.NewInvalidNode(
			errToken.Position,
			errToken,
		)
	}
	if lastPos == nil {
		return constant
	}

	return ast.NewConstructorCallNode(
		constant.Pos().Join(lastPos),
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
			tok.Pos().Join(right.Pos()),
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
			left.Pos().Join(right.Pos()),
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
			tok.Pos().Join(right.Pos()),
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
			left.Pos().Join(right.Pos()),
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
			tok.Position,
			tok.Value,
		)
	}

	if tok, ok := p.matchOk(token.PUBLIC_CONSTANT); ok {
		return ast.NewPublicConstantNode(
			tok.Position,
			tok.Value,
		)
	}

	p.errorExpected("a constant")
	tok := p.advance()
	p.mode = panicMode
	return ast.NewInvalidNode(
		tok.Position,
		tok,
	)
}

func (p *Parser) primaryExpression() ast.ExpressionNode {
	switch p.lookahead.Type {
	case token.TRUE:
		tok := p.advance()
		return ast.NewTrueLiteralNode(tok.Position)
	case token.FALSE:
		tok := p.advance()
		return ast.NewFalseLiteralNode(tok.Position)
	case token.NIL:
		tok := p.advance()
		return ast.NewNilLiteralNode(tok.Position)
	case token.SELF:
		return p.selfLiteral()
	case token.BREAK:
		tok := p.advance()
		return ast.NewBreakExpressionNode(tok.Position)
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
	case token.CHAR_LITERAL:
		return p.charLiteral()
	case token.RAW_CHAR_LITERAL:
		return p.rawCharLiteral()
	case token.RAW_STRING:
		return p.rawStringLiteral()
	case token.STRING_BEG:
		return p.stringLiteral()
	case token.COLON:
		return p.symbolOrNamedValueLiteral()
	case token.OR, token.OR_OR:
		return p.closureExpression()
	case token.VAR:
		return p.variableDeclaration()
	case token.CONST:
		return p.constantDeclaration()
	case token.DEF:
		return p.methodDefinition()
	case token.INIT:
		return p.initDefinition()
	case token.IF:
		return p.ifExpression()
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
			tok.Position,
			tok.Value,
		)
	case token.INT64:
		tok := p.advance()
		return ast.NewInt64LiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.UINT64:
		tok := p.advance()
		return ast.NewUInt64LiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.INT32:
		tok := p.advance()
		return ast.NewInt32LiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.UINT32:
		tok := p.advance()
		return ast.NewUInt32LiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.INT16:
		tok := p.advance()
		return ast.NewInt16LiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.UINT16:
		tok := p.advance()
		return ast.NewUInt16LiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.INT8:
		tok := p.advance()
		return ast.NewInt8LiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.UINT8:
		tok := p.advance()
		return ast.NewUInt8LiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.FLOAT:
		tok := p.advance()
		return ast.NewFloatLiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.BIG_FLOAT:
		tok := p.advance()
		return ast.NewBigFloatLiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.FLOAT64:
		tok := p.advance()
		return ast.NewFloat64LiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.FLOAT32:
		tok := p.advance()
		return ast.NewFloat32LiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.ERROR:
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Position,
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
	case token.TYPE:
		return p.typeLiteral()
	default:
		p.errorExpected("an expression")
		p.mode = panicMode
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Position,
			tok,
		)
	}
}

type specialCollectionLiteralConstructor[Element ast.ExpressionNode] func(*position.Position, []Element) ast.ExpressionNode

// specialCollectionLiteral = beginTokenType (elementProduction)* endTokenType
func specialCollectionLiteral[Element ast.ExpressionNode](p *Parser, elementProduction func() Element, constructor specialCollectionLiteralConstructor[Element], endTokenType token.Type) ast.ExpressionNode {
	begTok := p.advance()
	content := repeatedProduction(p, elementProduction, endTokenType)
	endTok, ok := p.consume(endTokenType)

	if !ok {
		return ast.NewInvalidNode(endTok.Position, endTok)
	}

	return constructor(
		begTok.Position.Join(endTok.Position),
		content,
	)
}

// wordListLiteral = "%w[" (rawString)* "]"
func (p *Parser) wordListLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.wordCollectionElement,
		ast.NewWordListLiteralNodeI,
		token.WORD_LIST_END,
	)
}

// wordTupleLiteral = "%w(" (rawString)* ")"
func (p *Parser) wordTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.wordCollectionElement,
		ast.NewWordTupleLiteralNodeI,
		token.WORD_TUPLE_END,
	)
}

// wordSetLiteral = "%w{" (rawString)* "}"
func (p *Parser) wordSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.wordCollectionElement,
		ast.NewWordSetLiteralNodeI,
		token.WORD_SET_END,
	)
}

// symbolListLiteral = "%s[" (rawString)* "]"
func (p *Parser) symbolListLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolListLiteralNodeI,
		token.SYMBOL_LIST_END,
	)
}

// symbolTupleLiteral = "%s(" (rawString)* ")"
func (p *Parser) symbolTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolTupleLiteralNodeI,
		token.SYMBOL_TUPLE_END,
	)
}

// symbolSetLiteral = "%s{" (rawString)* "}"
func (p *Parser) symbolSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolSetLiteralNodeI,
		token.SYMBOL_SET_END,
	)
}

// hexListLiteral = "%x[" (HEX_INT)* "]"
func (p *Parser) hexListLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewHexListLiteralNodeI,
		token.HEX_LIST_END,
	)
}

// hexTupleLiteral = "%x(" (HEX_INT)* ")"
func (p *Parser) hexTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewHexTupleLiteralNodeI,
		token.HEX_TUPLE_END,
	)
}

// hexSetLiteral = "%x{" (HEX_INT)* "}"
func (p *Parser) hexSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewHexSetLiteralNodeI,
		token.HEX_SET_END,
	)
}

// binListLiteral = "%b[" (BIN_INT)* "]"
func (p *Parser) binListLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewBinListLiteralNodeI,
		token.BIN_LIST_END,
	)
}

// binTupleLiteral = "%b(" (BIN_INT)* ")"
func (p *Parser) binTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewBinTupleLiteralNodeI,
		token.BIN_TUPLE_END,
	)
}

// binTupleLiteral = "%b{" (BIN_INT)* "}"
func (p *Parser) binSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteral(
		p,
		p.intCollectionElement,
		ast.NewBinSetLiteralNodeI,
		token.BIN_SET_END,
	)
}

// typeLiteral = "type" typeAnnotation
func (p *Parser) typeLiteral() ast.ExpressionNode {
	typeTok := p.advance()
	typeExpr := p.typeAnnotation()

	return ast.NewTypeLiteralNode(
		typeTok.Position.Join(typeExpr.Pos()),
		typeExpr,
	)
}

type listLikeConstructor func(*position.Position, []ast.ExpressionNode) ast.ExpressionNode
type collectionElementsProduction func(...token.Type) []ast.ExpressionNode

// collectionLiteral = startTok [elementsProduction] endTok
func (p *Parser) collectionLiteral(endTokType token.Type, elementsProduction collectionElementsProduction, constructor listLikeConstructor) ast.ExpressionNode {
	startTok := p.advance()
	p.swallowNewlines()

	if endTok, ok := p.matchOk(endTokType); ok {
		return constructor(
			startTok.Position.Join(endTok.Position),
			nil,
		)
	}

	elements := elementsProduction(endTokType)
	p.swallowNewlines()
	endTok, ok := p.consume(endTokType)
	if !ok {
		return ast.NewInvalidNode(
			endTok.Position,
			endTok,
		)
	}

	return constructor(
		startTok.Position.Join(endTok.Position),
		elements,
	)
}

// collectionElementModifier = subProduction |
// subProduction ("if" | "unless") expressionWithoutModifier |
// subProduction "if" expressionWithoutModifier "else" expressionWithoutModifier |
// subProduction "for" loopParameterList "in" expressionWithoutModifier
func (p *Parser) collectionElementModifier(subProduction func() ast.ExpressionNode) ast.ExpressionNode {
	left := subProduction()

	switch p.lookahead.Type {
	case token.UNLESS:
		mod := p.advance()
		p.swallowNewlines()
		right := p.expressionWithoutModifier()
		return ast.NewModifierNode(
			left.Pos().Join(right.Pos()),
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
				left.Pos().Join(elseExpr.Pos()),
				left,
				cond,
				elseExpr,
			)
		}
		return ast.NewModifierNode(
			left.Pos().Join(cond.Pos()),
			ifTok,
			left,
			cond,
		)
	case token.FOR:
		p.advance()
		p.swallowNewlines()
		params := p.loopParameterList(token.IN)
		inTok, ok := p.consume(token.IN)
		if !ok {
			return ast.NewInvalidNode(inTok.Position, inTok)
		}
		p.swallowNewlines()
		inExpr := p.expressionWithoutModifier()
		return ast.NewModifierForInNode(
			left.Pos().Join(inExpr.Pos()),
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

// listLiteral = "[" [listLikeLiteralElements] "]"
func (p *Parser) listLiteral() ast.ExpressionNode {
	return p.collectionLiteral(token.RBRACKET, p.listLikeLiteralElements, ast.NewListLiteralNodeI)
}

// tupleLiteral = "%(" [listLikeLiteralElements] ")"
func (p *Parser) tupleLiteral() ast.ExpressionNode {
	return p.collectionLiteral(token.RPAREN, p.listLikeLiteralElements, ast.NewTupleLiteralNodeI)
}

// listLikeLiteralElements = listLikeLiteralElement ("," listLikeLiteralElement)*
func (p *Parser) listLikeLiteralElements(stopTokens ...token.Type) []ast.ExpressionNode {
	return commaSeparatedList(p, p.listLikeLiteralElement, stopTokens...)
}

// listLikeLiteralElement = keyValueExpression |
// keyValueExpression ("if" | "unless") expressionWithoutModifier |
// keyValueExpression "if" expressionWithoutModifier "else" expressionWithoutModifier |
// keyValueExpression "for" loopParameterList "in" expressionWithoutModifier
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
// keyValueMapExpression "for" loopParameterList "in" expressionWithoutModifier
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
			key.Pos().Join(val.Pos()),
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
			p.errorMessagePos("expected a key-value pair, map literals should consist of key-value pairs", key.Pos())
			return key
		}
	}

	p.swallowNewlines()
	val := p.expressionWithoutModifier()

	return ast.NewKeyValueExpressionNode(
		key.Pos().Join(val.Pos()),
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
			key.Pos().Join(value.Pos()),
			key,
			value,
		)
	}

	return key
}

// setLiteral = "%{" [setLiteralElements] "}"
func (p *Parser) setLiteral() ast.ExpressionNode {
	return p.collectionLiteral(token.RBRACE, p.setLiteralElements, ast.NewSetLiteralNodeI)
}

// setLiteralElements = setLiteralElement ("," setLiteralElement)*
func (p *Parser) setLiteralElements(stopTokens ...token.Type) []ast.ExpressionNode {
	return commaSeparatedList(p, p.setLiteralElement, stopTokens...)
}

// setLiteralElement = expressionWithoutModifier |
// expressionWithoutModifier ("if" | "unless") expressionWithoutModifier |
// expressionWithoutModifier "if" expressionWithoutModifier "else" expressionWithoutModifier |
// expressionWithoutModifier "for" loopParameterList "in" expressionWithoutModifier
func (p *Parser) setLiteralElement() ast.ExpressionNode {
	return p.collectionElementModifier(p.expressionWithoutModifier)
}

// selfLiteral = "self"
func (p *Parser) selfLiteral() *ast.SelfLiteralNode {
	tok := p.advance()
	return ast.NewSelfLiteralNode(tok.Position)
}

// genericConstantList = genericConstant ("," genericConstant)*
func (p *Parser) genericConstantList(stopTokens ...token.Type) []ast.ComplexConstantNode {
	return commaSeparatedList(p, p.genericConstant, stopTokens...)
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

// methodDefinition = "sig" METHOD_NAME ["(" signatureParameterList ")"] [":" typeAnnotation] ["!" typeAnnotation]
func (p *Parser) methodSignatureDefinition() ast.ExpressionNode {
	var params []ast.ParameterNode
	var returnType ast.TypeNode
	var throwType ast.TypeNode
	var pos *position.Position

	sigTok := p.advance()
	if !p.lookahead.IsValidMethodName() {
		p.errorExpected("a method name (identifier, overridable operator)")
	}

	methodName := p.advance()
	pos = sigTok.Position.Join(methodName.Position)

	if p.match(token.LPAREN) {
		if rparen, ok := p.matchOk(token.RPAREN); ok {
			pos = pos.Join(rparen.Position)
		} else {
			params = p.signatureParameterList(token.RPAREN)

			rparen, ok := p.consume(token.RPAREN)
			pos = pos.Join(rparen.Position)
			if !ok {
				return ast.NewInvalidNode(
					rparen.Position,
					rparen,
				)
			}
		}
	}

	// return type
	if p.match(token.COLON) {
		returnType = p.typeAnnotation()
		pos = pos.Join(returnType.Pos())
	}

	// throw type
	if p.match(token.BANG) {
		throwType = p.typeAnnotation()
		pos = pos.Join(throwType.Pos())
	}

	return ast.NewMethodSignatureDefinitionNode(
		pos,
		methodName.StringValue(),
		params,
		returnType,
		throwType,
	)
}

// typeDeclaration = "alias" identifier identifier
func (p *Parser) aliasExpression() ast.ExpressionNode {
	aliasTok := p.advance()
	p.swallowNewlines()

	var (
		lastPos *position.Position
		newName string
		oldName string
	)
	if !p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) {
		p.errorExpected("an identifier")
		errTok := p.advance()
		return ast.NewInvalidNode(errTok.Position, errTok)
	}
	newNameTok := p.advance()
	newName = newNameTok.Value
	if p.match(token.EQUAL_OP) {
		newName = newName + "="
	}
	p.swallowNewlines()

	if !p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) {
		p.errorExpected("an identifier")
		errTok := p.advance()
		return ast.NewInvalidNode(errTok.Position, errTok)
	}
	oldNameTok := p.advance()
	oldName = oldNameTok.Value
	lastPos = oldNameTok.Position

	if eq, ok := p.matchOk(token.EQUAL_OP); ok {
		oldName = oldName + "="
		lastPos = eq.Position
	}

	return ast.NewAliasExpressionNode(
		aliasTok.Position.Join(lastPos),
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
		return ast.NewInvalidNode(equalTok.Position, equalTok)
	}
	p.swallowNewlines()

	typ := p.typeAnnotation()
	return ast.NewTypeDefinitionNode(
		typedefTok.Position.Join(typ.Pos()),
		name,
		typ,
	)
}

// methodDefinition = "def" METHOD_NAME ["(" methodParameterList ")"] [":" typeAnnotation] ["!" typeAnnotation] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) methodDefinition() ast.ExpressionNode {
	var params []ast.ParameterNode
	var returnType ast.TypeNode
	var throwType ast.TypeNode
	var body []ast.StatementNode
	var pos *position.Position
	var methodName string

	defTok := p.advance()
	p.swallowNewlines()
	if p.lookahead.IsValidRegularMethodName() {
		methodNameTok := p.advance()
		methodName = methodNameTok.StringValue()
		if p.match(token.EQUAL_OP) {
			methodName += "="
		}
	} else {
		if !p.lookahead.IsOverridableOperator() {
			p.errorExpected("a method name (identifier, overridable operator)")
		}
		methodName = p.advance().StringValue()
	}

	if p.match(token.LPAREN) {
		p.swallowNewlines()
		if !p.match(token.RPAREN) {
			params = p.methodParameterList(token.RPAREN, token.STAR)

			p.swallowNewlines()
			if tok, ok := p.consume(token.RPAREN); !ok {
				return ast.NewInvalidNode(
					tok.Position,
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

	lastPos, body, multiline := p.statementBlockWithThen(token.END)
	if lastPos != nil {
		pos = defTok.Position.Join(lastPos)
	} else {
		pos = defTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewMethodDefinitionNode(
		pos,
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
	var pos *position.Position

	initTok := p.advance()

	// methodParameterList
	if p.match(token.LPAREN) {
		p.swallowNewlines()
		if !p.match(token.RPAREN) {
			params = p.methodParameterList(token.RPAREN)

			p.swallowNewlines()
			if tok, ok := p.consume(token.RPAREN); !ok {
				return ast.NewInvalidNode(
					tok.Position,
					tok,
				)
			}
		}
	}

	// throw type
	if p.match(token.BANG) {
		throwType = p.typeAnnotation()
	}

	lastPos, body, multiline := p.statementBlockWithThen(token.END)
	if lastPos != nil {
		pos = initTok.Position.Join(lastPos)
	} else {
		pos = initTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewInitDefinitionNode(
		pos,
		params,
		throwType,
		body,
	)
}

// typeVariable = ["+" | "-"] constant ["<" strictConstantLookup]
func (p *Parser) typeVariable() ast.TypeVariableNode {
	variance := ast.INVARIANT
	var firstPos *position.Position
	var lastPos *position.Position
	var upperBound ast.ComplexConstantNode

	switch p.lookahead.Type {
	case token.PLUS:
		plusTok := p.advance()
		firstPos = plusTok.Position
		variance = ast.COVARIANT
	case token.MINUS:
		minusTok := p.advance()
		firstPos = minusTok.Position
		variance = ast.CONTRAVARIANT
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
	default:
		errTok := p.advance()
		p.errorExpected("a type variable")
		return ast.NewInvalidNode(
			errTok.Position,
			errTok,
		)
	}

	if !p.accept(token.PRIVATE_CONSTANT, token.PUBLIC_CONSTANT) {
		errTok := p.advance()
		return ast.NewInvalidNode(
			errTok.Position,
			errTok,
		)
	}
	nameTok := p.advance()
	if firstPos == nil {
		firstPos = nameTok.Position
	}
	lastPos = nameTok.Position

	if p.match(token.LESS) {
		upperBound = p.strictConstantLookup()
		lastPos = upperBound.Pos()
	}

	return ast.NewVariantTypeVariableNode(
		firstPos.Join(lastPos),
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
			p.errorMessagePos("invalid class name, expected a constant", constant.Pos())
		}
	}
	var pos *position.Position

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeVariableList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Position,
					errTok,
				)
			}
		}
	}

	if p.match(token.LESS) {
		superclass = p.genericConstant()
	}

	lastPos, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastPos != nil {
		pos = classTok.Position.Join(lastPos)
	} else {
		pos = classTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewClassDeclarationNode(
		pos,
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
			p.errorMessagePos("invalid module name, expected a constant", constant.Pos())
		}
	}
	var pos *position.Position

	if lbracket, ok := p.matchOk(token.LBRACKET); ok {
		errPos := lbracket.Pos()
		if p.accept(token.RBRACKET) {
			rbracket := p.advance()
			errPos = errPos.Join(rbracket.Position)
		} else {
			p.typeVariableList()
			rbracket, ok := p.consume(token.RBRACKET)
			if !ok {
				return ast.NewInvalidNode(
					rbracket.Position,
					rbracket,
				)
			}
			errPos = errPos.Join(rbracket.Position)
		}
		p.errorMessagePos("modules can't be generic", errPos)
	}

	lastPos, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastPos != nil {
		pos = moduleTok.Position.Join(lastPos)
	} else {
		pos = moduleTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewModuleDeclarationNode(
		pos,
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
			p.errorMessagePos("invalid mixin name, expected a constant", constant.Pos())
		}
	}
	var pos *position.Position

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeVariableList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Position,
					errTok,
				)
			}
		}
	}

	lastPos, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastPos != nil {
		pos = mixinTok.Position.Join(lastPos)
	} else {
		pos = mixinTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewMixinDeclarationNode(
		pos,
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
			p.errorMessagePos("invalid interface name, expected a constant", constant.Pos())
		}
	}
	var pos *position.Position

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeVariableList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Position,
					errTok,
				)
			}
		}
	}

	lastPos, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastPos != nil {
		pos = interfaceTok.Position.Join(lastPos)
	} else {
		pos = interfaceTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewInterfaceDeclarationNode(
		pos,
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
			p.errorMessagePos("invalid struct name, expected a constant", constant.Pos())
		}
	}
	var pos *position.Position

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeVariableList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Position,
					errTok,
				)
			}
		}
	}

	lastPos, thenBody, multiline := p.structBodyStatementBlockWithThen(token.END)
	if lastPos != nil {
		pos = structTok.Position.Join(lastPos)
	} else {
		pos = structTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewStructDeclarationNode(
		pos,
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
			tok.Position,
			tok,
		)
	}
	lastPos := varName.Position

	if p.match(token.COLON) {
		typ = p.typeAnnotation()
		lastPos = typ.Pos()
	}

	if p.match(token.EQUAL_OP) {
		p.swallowNewlines()
		init = p.expressionWithoutModifier()
		lastPos = init.Pos()
	}

	return ast.NewVariableDeclarationNode(
		varTok.Position.Join(lastPos),
		varName,
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
		p.errorExpected("an uppercased identifier as the name of the declared constant")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Position,
			tok,
		)
	}
	lastPos := constName.Position

	if p.match(token.COLON) {
		typ = p.typeAnnotation()
		lastPos = typ.Pos()
	}

	if p.match(token.EQUAL_OP) {
		p.swallowNewlines()
		init = p.expressionWithoutModifier()
		lastPos = init.Pos()
	} else {
		p.errorMessagePos("constants must be initialised", constTok.Position.Join(lastPos))
	}

	return ast.NewConstantDeclarationNode(
		constTok.Position.Join(lastPos),
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
			primType.Pos().Join(questTok.Position),
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
		return ast.NewInvalidNode(rbracket.Position, rbracket)
	}

	return ast.NewGenericConstantNode(
		constant.Pos().Join(rbracket.Position),
		constant,
		constList,
	)
}

// throwExpression = "throw" [expressionWithoutModifier]
func (p *Parser) throwExpression() *ast.ThrowExpressionNode {
	throwTok := p.advance()
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() {
		return ast.NewThrowExpressionNode(
			throwTok.Position,
			nil,
		)
	}

	expr := p.expressionWithoutModifier()

	return ast.NewThrowExpressionNode(
		throwTok.Position.Join(expr.Pos()),
		expr,
	)
}

// continueExpression = "continue" [expressionWithoutModifier]
func (p *Parser) continueExpression() *ast.ContinueExpressionNode {
	continueTok := p.advance()
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() {
		return ast.NewContinueExpressionNode(
			continueTok.Position,
			nil,
		)
	}

	expr := p.expressionWithoutModifier()

	return ast.NewContinueExpressionNode(
		continueTok.Position.Join(expr.Pos()),
		expr,
	)
}

// returnExpression = "return" [expressionWithoutModifier]
func (p *Parser) returnExpression() *ast.ReturnExpressionNode {
	returnTok := p.advance()
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() {
		return ast.NewReturnExpressionNode(
			returnTok.Position,
			nil,
		)
	}

	expr := p.expressionWithoutModifier()

	return ast.NewReturnExpressionNode(
		returnTok.Position.Join(expr.Pos()),
		expr,
	)
}

// loopExpression = "loop" ((SEPARATOR [statements] "end") | expressionWithoutModifier)
func (p *Parser) loopExpression() *ast.LoopExpressionNode {
	loopTok := p.advance()
	var pos *position.Position

	lastPos, thenBody, multiline := p.statementBlock(token.END)
	if lastPos != nil {
		pos = loopTok.Position.Join(lastPos)
	} else {
		pos = loopTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewLoopExpressionNode(
		pos,
		thenBody,
	)
}

// forExpression = "for" loopParameterList "in" expressionWithoutModifier ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) forExpression() ast.ExpressionNode {
	forTok := p.advance()
	p.swallowNewlines()
	loopParameters := p.loopParameterList(token.IN)

	inTok, ok := p.consume(token.IN)
	if !ok {
		return ast.NewInvalidNode(inTok.Position, inTok)
	}
	p.swallowNewlines()
	inExpr := p.expressionWithoutModifier()
	var pos *position.Position

	lastPos, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastPos != nil {
		pos = forTok.Position.Join(lastPos)
	} else {
		pos = forTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewForExpressionNode(
		pos,
		loopParameters,
		inExpr,
		thenBody,
	)
}

// whileExpression = "while" expressionWithoutModifier ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) whileExpression() *ast.WhileExpressionNode {
	whileTok := p.advance()
	cond := p.expressionWithoutModifier()
	var pos *position.Position

	lastPos, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastPos != nil {
		pos = whileTok.Position.Join(lastPos)
	} else {
		pos = whileTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewWhileExpressionNode(
		pos,
		cond,
		thenBody,
	)
}

// untilExpression = "until" expressionWithoutModifier ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) untilExpression() *ast.UntilExpressionNode {
	untilTok := p.advance()
	cond := p.expressionWithoutModifier()
	var pos *position.Position

	lastPos, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastPos != nil {
		pos = untilTok.Position.Join(lastPos)
	} else {
		pos = untilTok.Position
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewUntilExpressionNode(
		pos,
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
	var pos *position.Position

	lastPos, thenBody, multiline := p.statementBlockWithThen(token.END, token.ELSE)
	if lastPos != nil {
		pos = unlessTok.Position.Join(lastPos)
	} else {
		pos = unlessTok.Position
	}

	unlessExpr := ast.NewUnlessExpressionNode(
		pos,
		cond,
		thenBody,
		nil,
	)
	currentExpr := unlessExpr

	if p.lookahead.IsStatementSeparator() && p.nextLookahead.Type == token.ELSE {
		p.advance()
		p.advance()
		lastPos, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastPos != nil {
			*currentExpr.Position = *currentExpr.Position.Join(lastPos)
		}
	} else if p.match(token.ELSE) {
		lastPos, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastPos != nil {
			*currentExpr.Position = *currentExpr.Position.Join(lastPos)
		}
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			*currentExpr.Position = *currentExpr.Position.Join(endTok.Position)
		}
	}
	unlessExpr.Position = unlessExpr.Position.Join(currentExpr.Position)

	return unlessExpr
}

// ifExpression = "if" expressionWithoutModifier ((SEPARATOR [statements]) | ("then" expressionWithoutModifier))
// ("elsif" expressionWithoutModifier ((SEPARATOR [statements]) | ("then" expressionWithoutModifier)) )*
// ["else" ((SEPARATOR [statements]) | expressionWithoutModifier)]
// "end"
func (p *Parser) ifExpression() *ast.IfExpressionNode {
	ifTok := p.advance()
	cond := p.expressionWithoutModifier()
	var pos *position.Position

	lastPos, thenBody, multiline := p.statementBlockWithThen(token.END, token.ELSE, token.ELSIF)
	if lastPos != nil {
		pos = ifTok.Position.Join(lastPos)
	} else {
		pos = ifTok.Position
	}

	ifExpr := ast.NewIfExpressionNode(
		pos,
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

		lastPos, thenBody, multiline = p.statementBlockWithThen(token.END, token.ELSE, token.ELSIF)
		if lastPos != nil {
			pos = elsifTok.Position.Join(lastPos)
		} else {
			pos = elsifTok.Position
		}

		elsifExpr := ast.NewIfExpressionNode(
			pos,
			cond,
			thenBody,
			nil,
		)

		currentExpr.ElseBody = []ast.StatementNode{
			ast.NewExpressionStatementNode(
				elsifExpr.Position,
				elsifExpr,
			),
		}
		currentExpr = elsifExpr
	}

	if p.lookahead.IsStatementSeparator() && p.nextLookahead.Type == token.ELSE {
		p.advance()
		p.advance()
		lastPos, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastPos != nil {
			*currentExpr.Position = *currentExpr.Position.Join(lastPos)
		}
	} else if p.match(token.ELSE) {
		lastPos, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastPos != nil {
			*currentExpr.Position = *currentExpr.Position.Join(lastPos)
		}
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			*currentExpr.Position = *currentExpr.Position.Join(endTok.Position)
		}
	}
	ifExpr.Position = ifExpr.Position.Join(currentExpr.Position)

	return ifExpr
}

// symbolOrNamedValueLiteral = ":" (identifier | constant | rawStringLiteral) "{" [expressionWithoutModifier] "}" | ":" (identifier | constant | rawStringLiteral | stringLiteral)
func (p *Parser) symbolOrNamedValueLiteral() ast.ExpressionNode {
	symbolBegTok := p.advance()
	if p.lookahead.IsValidSimpleSymbolContent() {
		contTok := p.advance()
		if p.match(token.LBRACE) {
			if p.accept(token.RBRACE) {
				rbraceTok := p.advance()
				return ast.NewNamedValueLiteralNode(
					symbolBegTok.Position.Join(rbraceTok.Position),
					contTok.StringValue(),
					nil,
				)
			}
			expr := p.expressionWithoutModifier()
			rbraceTok, ok := p.consume(token.RBRACE)
			if !ok {
				return ast.NewInvalidNode(
					symbolBegTok.Position.Join(rbraceTok.Position),
					rbraceTok,
				)
			}

			return ast.NewNamedValueLiteralNode(
				symbolBegTok.Position.Join(rbraceTok.Position),
				contTok.StringValue(),
				expr,
			)
		}
		return ast.NewSimpleSymbolLiteralNode(
			symbolBegTok.Position.Join(contTok.Position),
			contTok.StringValue(),
		)
	}

	if !p.accept(token.STRING_BEG) {
		p.errorExpected("an identifier, overridable operator or string literal")
		p.mode = panicMode
		tok := p.advance()
		return ast.NewInvalidNode(
			symbolBegTok.Position.Join(tok.Position),
			tok,
		)
	}

	str := p.stringLiteral()
	switch s := str.(type) {
	case *ast.DoubleQuotedStringLiteralNode:
		return ast.NewSimpleSymbolLiteralNode(
			symbolBegTok.Position.Join(s.Position),
			s.Value,
		)
	case *ast.InterpolatedStringLiteralNode:
		return ast.NewInterpolatedSymbolLiteral(
			symbolBegTok.Position.Join(s.Position),
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
		tok.Position,
		char,
	)
}

// charLiteral = RAW_CHAR_LITERAL
func (p *Parser) rawCharLiteral() *ast.RawCharLiteralNode {
	tok := p.advance()
	char, _ := utf8.DecodeRuneInString(tok.Value)
	return ast.NewRawCharLiteralNode(
		tok.Position,
		char,
	)
}

// rawStringLiteral = RAW_STRING
func (p *Parser) rawStringLiteral() *ast.RawStringLiteralNode {
	tok := p.advance()
	return ast.NewRawStringLiteralNode(
		tok.Position,
		tok.Value,
	)
}

// wordCollectionElement = RAW_STRING
func (p *Parser) wordCollectionElement() ast.WordCollectionContentNode {
	tok, ok := p.consume(token.RAW_STRING)
	if !ok {
		return ast.NewInvalidNode(tok.Position, tok)
	}
	return ast.NewRawStringLiteralNode(
		tok.Position,
		tok.Value,
	)
}

// symbolCollectionElement = RAW_STRING
func (p *Parser) symbolCollectionElement() ast.SymbolCollectionContentNode {
	tok, ok := p.consume(token.RAW_STRING)
	if !ok {
		return ast.NewInvalidNode(tok.Position, tok)
	}
	return ast.NewSimpleSymbolLiteralNode(
		tok.Position,
		tok.Value,
	)
}

// intCollectionElement = INT
func (p *Parser) intCollectionElement() ast.IntCollectionContentNode {
	tok, ok := p.consume(token.INT)
	if !ok {
		return ast.NewInvalidNode(tok.Position, tok)
	}
	return ast.NewIntLiteralNode(tok.Position, tok.Value)
}

// stringLiteral = "\"" (STRING_CONTENT | "${" expressionWithoutModifier "}")* "\""
func (p *Parser) stringLiteral() ast.StringLiteralNode {
	quoteBeg := p.advance() // consume the opening quote
	var quoteEnd *token.Token

	var strContent []ast.StringLiteralContentNode
	for {
		if tok, ok := p.matchOk(token.STRING_CONTENT); ok {
			strContent = append(strContent, ast.NewStringLiteralContentSectionNode(
				tok.Position,
				tok.Value,
			))
			continue
		}

		if beg, ok := p.matchOk(token.STRING_INTERP_BEG); ok {
			expr := p.expressionWithoutModifier()
			end, _ := p.consume(token.STRING_INTERP_END)
			strContent = append(strContent, ast.NewStringInterpolationNode(
				beg.Position.Join(end.Position),
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
				tok.Position,
				tok,
			))
			continue
		}
		break
	}

	strVal, ok := strContent[0].(*ast.StringLiteralContentSectionNode)
	if len(strContent) == 1 && ok {
		return ast.NewDoubleQuotedStringLiteralNode(
			quoteBeg.Position.Join(quoteEnd.Position),
			strVal.Value,
		)
	}

	return ast.NewInterpolatedStringLiteralNode(
		quoteBeg.Position.Join(quoteEnd.Position),
		strContent,
	)
}

// closureAfterArrow = "->" (expressionWithoutModifier | SEPARATOR [statements] "end" | "{" [statements] "}")
func (p *Parser) closureAfterArrow(firstPos *position.Position, params []ast.ParameterNode, returnType ast.TypeNode, throwType ast.TypeNode) ast.ExpressionNode {
	var pos *position.Position
	arrowTok, ok := p.consume(token.THIN_ARROW)
	if !ok {
		return ast.NewInvalidNode(
			arrowTok.Position,
			arrowTok,
		)
	}
	if firstPos == nil {
		firstPos = arrowTok.Position
	}

	// Body with curly braces
	if p.match(token.LBRACE) {
		p.swallowNewlines()
		body := p.statements(token.RBRACE)
		if tok, ok := p.consume(token.RBRACE); ok {
			pos = firstPos.Join(tok.Position)
		} else {
			pos = firstPos
		}
		return ast.NewClosureLiteralNode(
			pos,
			params,
			returnType,
			throwType,
			body,
		)
	}

	lastPos, body, multiline := p.statementBlock(token.END)
	if lastPos != nil {
		pos = firstPos.Join(lastPos)
	} else {
		pos = firstPos
	}

	if multiline {
		endTok, ok := p.consume(token.END)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return ast.NewClosureLiteralNode(
		pos,
		params,
		returnType,
		throwType,
		body,
	)
}

// closureExpression = (("|" formalParameterList "|") | "||") [: typeAnnotation] ["!" typeAnnotation] closureAfterArrow
func (p *Parser) closureExpression() ast.ExpressionNode {
	var params []ast.ParameterNode
	var firstPos *position.Position
	var returnType ast.TypeNode
	var throwType ast.TypeNode

	if p.accept(token.OR) {
		firstPos = p.advance().Position
		if !p.accept(token.OR) {
			p.mode = withoutBitwiseOrMode
			params = p.formalParameterList(token.OR)
			p.mode = normalMode
		}
		if tok, ok := p.consume(token.OR); !ok {
			return ast.NewInvalidNode(
				tok.Position,
				tok,
			)
		}
	} else {
		orOr, ok := p.consume(token.OR_OR)
		firstPos = orOr.Position
		if !ok {
			return ast.NewInvalidNode(
				orOr.Position,
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

	return p.closureAfterArrow(firstPos, params, returnType, throwType)
}

// identifierOrClosure = identifier | identifier closureAfterArrow
func (p *Parser) identifierOrClosure() ast.ExpressionNode {
	if p.nextLookahead.Type == token.THIN_ARROW {
		ident := p.advance()
		return p.closureAfterArrow(
			ident.Position,
			[]ast.ParameterNode{
				ast.NewFormalParameterNode(
					ident.Position,
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

// identifier = PUBLIC_IDENTIFIER | PRIVATE_IDENTIFIER
func (p *Parser) identifier() ast.IdentifierNode {
	if p.accept(token.PUBLIC_IDENTIFIER) {
		ident := p.advance()
		return ast.NewPublicIdentifierNode(
			ident.Position,
			ident.Value,
		)
	}
	if p.accept(token.PRIVATE_IDENTIFIER) {
		ident := p.advance()
		return ast.NewPrivateIdentifierNode(
			ident.Position,
			ident.Value,
		)
	}

	p.errorExpected("an identifier")
	errTok := p.advance()
	return ast.NewInvalidNode(errTok.Position, errTok)
}
