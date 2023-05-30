// Package parser implements a parser
// used by the Elk interpreter.
//
// Parser expects a slice of bytes containing Elk source code
// parses it, registering any encountered errors, and returns an Abstract Syntax Tree.
package parser

import (
	"fmt"

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
	source        []byte       // Elk source code
	lexer         *lexer.Lexer // lexer which outputs a stream of tokens
	lookahead     *token.Token // next token used for predicting productions
	nextLookahead *token.Token // second next token used for predicting productions
	errors        ErrorList
	mode          mode
}

// Instantiate a new parser.
func new(source []byte) *Parser {
	return &Parser{
		source: source,
		lexer:  lexer.New(source),
		mode:   normalMode,
	}
}

// Parse the given source code and return an Abstract Syntax Tree.
// Main entry point to the parser.
func Parse(source []byte) (*ast.ProgramNode, ErrorList) {
	return new(source).parse()
}

// Start the parsing process from the top.
func (p *Parser) parse() (*ast.ProgramNode, ErrorList) {
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
		pos,
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
	if p.lookahead.Type == token.ERROR {
		return p.advance(), false
	}

	if p.lookahead.Type != tokenType {
		p.errorExpectedToken(tokenType)
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

// Accept and ignore any number of consecutive end-line tokens.
func (p *Parser) swallowEndLines() {
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
		p.swallowEndLines()
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
	pos := keyword.Position
	consts := p.genericConstantList()
	if len(consts) > 0 {
		pos = pos.Join(consts[len(consts)-1].Pos())
	}

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
		p.swallowEndLines()
		elements = append(elements, elementProduction())
	}

	return elements
}

// A production that can be repeated as in `repeatableProduction*`
type repeatableProduction[Element ast.Node] func(...token.Type) Element

// Consume subProductions until one of the provided token types is encountered.
//
// repeatedProduction = subProduction*
func repeatedProduction[Element ast.Node](p *Parser, subProduction repeatableProduction[Element], stopTokens ...token.Type) []Element {
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
	return repeatedProduction(p, p.statement, stopTokens...)
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
	return repeatedProduction(p, p.structBodyStatement, stopTokens...)
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
// expressionWithoutModifier "if" expressionWithoutModifier "else" expressionWithoutModifier
func (p *Parser) modifierExpression() ast.ExpressionNode {
	left := p.expressionWithoutModifier()

	switch p.lookahead.Type {
	case token.UNLESS, token.WHILE, token.UNTIL:
		mod := p.advance()
		right := p.expressionWithoutModifier()
		return ast.NewModifierNode(
			left.Pos().Join(right.Pos()),
			mod,
			left,
			right,
		)
	case token.IF:
		mod := p.advance()
		cond := p.expressionWithoutModifier()
		if p.lookahead.Type == token.ELSE {
			p.advance()
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
			mod,
			left,
			cond,
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
	p.swallowEndLines()
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
	param, _ := p.formalParameterOptional()
	return param
}

// same as [formalParameter] but returns a boolean indicating whether
// the parameter is optional
func (p *Parser) formalParameterOptional() (ast.ParameterNode, bool) {
	var init ast.ExpressionNode
	var typ ast.TypeNode

	paramName, ok := p.matchOk(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER)
	if !ok {
		p.errorExpected("an identifier as the name of the declared formalParameter")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Position,
			tok,
		), init != nil
	}
	lastPos := paramName.Position

	if p.match(token.COLON) {
		typ = p.intersectionType()
		lastPos = typ.Pos()
	}

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		lastPos = init.Pos()
	}

	return ast.NewFormalParameterNode(
		paramName.Position.Join(lastPos.Pos()),
		paramName.Value,
		typ,
		init,
	), init != nil
}

// parameterList = parameter ("," parameter)*
func (p *Parser) parameterList(parameter func() (ast.ParameterNode, bool), stopTokens ...token.Type) []ast.ParameterNode {
	var elements []ast.ParameterNode
	element, optionalSeen := parameter()
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
		p.swallowEndLines()
		element, opt := parameter()
		if !opt && optionalSeen {
			p.errorMessagePos("required parameters can't appear after optional parameters", element.Pos())
		} else if opt && !optionalSeen {
			optionalSeen = true
		}
		elements = append(elements, element)
	}

	return elements
}

// formalParameterList = formalParameter ("," formalParameter)*
func (p *Parser) formalParameterList(stopTokens ...token.Type) []ast.ParameterNode {
	return p.parameterList(p.formalParameterOptional, stopTokens...)
}

// signatureParameter = identifier ["?"] [":" typeAnnotation]
func (p *Parser) signatureParameter() (ast.ParameterNode, bool) {
	var typ ast.TypeNode
	var opt bool

	paramName, ok := p.matchOk(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER)
	if !ok {
		p.errorExpected("an identifier as the name of the declared signatureParameter")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Position,
			tok,
		), opt
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
	), opt
}

// formalParameterList = signatureParameter ("," signatureParameter)*
func (p *Parser) signatureParameterList(stopTokens ...token.Type) []ast.ParameterNode {
	return p.parameterList(p.signatureParameter, stopTokens...)
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

		p.swallowEndLines()
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

		p.swallowEndLines()
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
		p.swallowEndLines()
		right := p.unaryExpression()
		return ast.NewUnaryExpressionNode(
			operator.Pos().Join(right.Pos()),
			operator,
			right,
		)
	}

	return p.powerExpression()
}

// powerExpression = constructorCall | constructorCall "**" powerExpression
func (p *Parser) powerExpression() ast.ExpressionNode {
	left := p.constructorCall()

	if p.lookahead.Type != token.STAR_STAR {
		return left
	}

	operator := p.advance()
	p.swallowEndLines()
	right := p.powerExpression()

	return ast.NewBinaryExpressionNode(
		left.Pos().Join(right.Pos()),
		operator,
		left,
		right,
	)
}

// positionalArgumentList = expressionWithoutModifier ("," expressionWithoutModifier)*
func (p *Parser) positionalArgumentList(stopTokens ...token.Type) []ast.ExpressionNode {
	var elements []ast.ExpressionNode
	if p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) && p.nextLookahead.Type == token.COLON {
		return elements
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
		p.swallowEndLines()
		if p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) && p.nextLookahead.Type == token.COLON {
			return elements
		}
		elements = append(elements, p.expressionWithoutModifier())
	}

	return elements
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

// constructorCall = constantLookup |
// strictConstantLookup "(" argumentList ")" |
// strictConstantLookup argumentList
func (p *Parser) constructorCall() ast.ExpressionNode {
	if !p.accept(token.PRIVATE_CONSTANT, token.PUBLIC_CONSTANT, token.SCOPE_RES_OP) {
		return p.constantLookup()
	}

	constant := p.strictConstantLookup()

	if p.match(token.LPAREN) {
		p.swallowEndLines()
		if rparen, ok := p.matchOk(token.RPAREN); ok {
			return ast.NewConstructorCallNode(
				constant.Pos().Join(rparen.Position),
				constant,
				nil,
				nil,
			)
		}
		posArgs := p.positionalArgumentList(token.RPAREN)
		var namedArgs []ast.NamedArgumentNode
		p.swallowEndLines()
		rparen, ok := p.matchOk(token.RPAREN)
		if !ok {
			if len(posArgs) > 0 {
				if comma, ok := p.consume(token.COMMA); !ok {
					return ast.NewInvalidNode(
						comma.Position,
						comma,
					)
				}
			}
			namedArgs = p.namedArgumentList(token.RPAREN)
			p.swallowEndLines()
			rparen, ok = p.consume(token.RPAREN)
			if !ok {
				return ast.NewInvalidNode(
					rparen.Position,
					rparen,
				)
			}
		}

		return ast.NewConstructorCallNode(
			constant.Pos().Join(rparen.Position),
			constant,
			posArgs,
			namedArgs,
		)
	}

	return constant
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

		p.swallowEndLines()
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

		p.swallowEndLines()
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
		tok := p.advance()
		return ast.NewSelfLiteralNode(tok.Position)
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
	case token.RAW_STRING:
		return p.rawStringLiteral()
	case token.STRING_BEG:
		return p.stringLiteral()
	case token.SYMBOL_BEG:
		return p.symbolOrNamedValueLiteral()
	case token.OR, token.OR_OR, token.THIN_ARROW:
		return p.closureExpression()
	case token.VAR:
		return p.variableDeclaration()
	case token.CONST:
		return p.constantDeclaration()
	case token.DEF:
		return p.methodDefinition()
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
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		return p.identifierOrClosure()
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		return p.constant()
	case token.HEX_INT, token.DUO_INT, token.DEC_INT,
		token.OCT_INT, token.QUAT_INT, token.BIN_INT:
		tok := p.advance()
		return ast.NewIntLiteralNode(
			tok.Position,
			tok,
		)
	case token.FLOAT:
		tok := p.advance()
		return ast.NewFloatLiteralNode(
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

// typeDeclaration = "alias" identifier "=" identifier
func (p *Parser) aliasExpression() ast.ExpressionNode {
	aliasTok := p.advance()

	newName := p.identifier()
	equalTok, ok := p.consume(token.EQUAL_OP)
	if !ok {
		return ast.NewInvalidNode(equalTok.Position, equalTok)
	}
	p.swallowEndLines()

	oldName := p.identifier()
	return ast.NewAliasExpressionNode(
		aliasTok.Position.Join(oldName.Pos()),
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
	p.swallowEndLines()

	typ := p.typeAnnotation()
	return ast.NewTypeDefinitionNode(
		typedefTok.Position.Join(typ.Pos()),
		name,
		typ,
	)
}

// methodDefinition = "def" METHOD_NAME ["(" formalParameterList ")"] [":" typeAnnotation] ["!" typeAnnotation] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) methodDefinition() ast.ExpressionNode {
	var params []ast.ParameterNode
	var returnType ast.TypeNode
	var throwType ast.TypeNode
	var body []ast.StatementNode
	var pos *position.Position

	defTok := p.advance()
	if !p.lookahead.IsValidMethodName() {
		p.errorExpected("a method name (identifier, overridable operator)")
	}

	methodName := p.advance()

	// formalParameterList
	if p.match(token.LPAREN) {
		if !p.match(token.RPAREN) {
			params = p.formalParameterList(token.RPAREN)

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
		methodName.StringValue(),
		params,
		returnType,
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
		p.swallowEndLines()
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
		p.swallowEndLines()
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
	return ast.NewComplexSymbolLiteralNode(
		symbolBegTok.Position.Join(str.Position),
		str,
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

// stringLiteral = "\"" (STRING_CONTENT | "${" expressionWithoutModifier "}")* "\""
func (p *Parser) stringLiteral() *ast.StringLiteralNode {
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
			strContent = append(strContent, &ast.StringInterpolationNode{
				Expression: expr,
				Position:   beg.Position.Join(end.Position),
			})
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

	return ast.NewStringLiteralNode(
		quoteBeg.Position.Join(quoteEnd.Position),
		strContent,
	)
}

// closureAfterArrow = "->" (expressionWithoutModifier | SEPARATOR [statements] "end" | "{" [statements] "}")
func (p *Parser) closureAfterArrow(firstPos *position.Position, params []ast.ParameterNode, returnType ast.TypeNode) ast.ExpressionNode {
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
		p.swallowEndLines()
		body := p.statements(token.RBRACE)
		if tok, ok := p.consume(token.RBRACE); ok {
			pos = firstPos.Join(tok.Position)
		} else {
			pos = firstPos
		}
		return ast.NewClosureExpressionNode(
			pos,
			params,
			returnType,
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

	return ast.NewClosureExpressionNode(
		pos,
		params,
		returnType,
		body,
	)
}

// closureExpression = [(("|" formalParameterList "|") | "||") [: typeAnnotation]] closureAfterArrow
func (p *Parser) closureExpression() ast.ExpressionNode {
	var params []ast.ParameterNode
	var firstPos *position.Position
	var returnType ast.TypeNode

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

		// return type
		if p.match(token.COLON) {
			returnType = p.typeAnnotation()
		}
	} else if p.accept(token.OR_OR) {
		firstPos = p.advance().Position
		if p.match(token.COLON) {
			returnType = p.typeAnnotation()
		}
	}

	return p.closureAfterArrow(firstPos, params, returnType)
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
				),
			},
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
