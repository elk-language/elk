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

// Same as [consume] but let's you specify a custom expected error message.
func (p *Parser) consumeExpected(tokenType token.Type, expectedMsg string) (*token.Token, bool) {
	if p.lookahead.Type == token.ERROR {
		return p.advance(), false
	}

	if p.lookahead.Type != tokenType {
		p.errorExpected(expectedMsg)
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

// Consume statements until the provided token type is encountered.
func (p *Parser) statementsWithStop(stopTokens ...token.Type) []ast.StatementNode {
	var statementList []ast.StatementNode

	for {
		if p.lookahead.Type == token.END_OF_FILE {
			return statementList
		}
		for _, stopToken := range stopTokens {
			if p.lookahead.Type == stopToken {
				return statementList
			}
		}
		statement := p.statement(stopTokens...)
		statementList = append(statementList, statement)
	}
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
func (p *Parser) statementBlockBody(stopTokens ...token.Type) (*position.Position, []ast.StatementNode, bool) {
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
			thenBody = p.statementsWithStop(stopTokens...)
			if len(thenBody) > 0 {
				lastPos = thenBody[len(thenBody)-1].Pos()
			}
		}
	}

	return lastPos, thenBody, multiline
}

// Consume a block of statements, like in `if`, `elsif` or `while` expressions,
// that terminates with `end` or can be single-line when it begins with `then`.
func (p *Parser) statementBlockBodyWithThen(stopTokens ...token.Type) (*position.Position, []ast.StatementNode, bool) {
	var thenBody []ast.StatementNode
	var lastPos *position.Position
	var multiline bool

	if p.lookahead.Type == token.THEN {
		p.advance()
		expr := p.expressionWithoutModifier()
		thenBody = append(thenBody, ast.NewExpressionStatementNode(
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
			thenBody = p.statementsWithStop(stopTokens...)
			if len(thenBody) > 0 {
				lastPos = thenBody[len(thenBody)-1].Pos()
			}
		}
	}

	return lastPos, thenBody, multiline
}

// binaryExpression = subProduction | binaryExpression operators subProduction
func (p *Parser) binaryExpression(subProduction func() ast.ExpressionNode, operators ...token.Type) ast.ExpressionNode {
	left := subProduction()

	for {
		operator, ok := p.matchOk(operators...)
		if !ok {
			break
		}
		p.swallowEndLines()
		right := subProduction()
		left = ast.NewBinaryExpressionNode(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
}

// logicalExpression = subProduction | logicalExpression operators subProduction
func (p *Parser) logicalExpression(subProduction func() ast.ExpressionNode, operators ...token.Type) ast.ExpressionNode {
	left := subProduction()

	for {
		operator, ok := p.matchOk(operators...)
		if !ok {
			break
		}
		p.swallowEndLines()
		right := subProduction()
		left = ast.NewLogicalExpressionNode(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
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
func (p *Parser) statements() []ast.StatementNode {
	return p.statementsWithStop()
}

// statement = emptyStatement | expressionStatement
func (p *Parser) statement(separators ...token.Type) ast.StatementNode {
	if p.lookahead.IsStatementSeparator() {
		return p.emptyStatement()
	}

	return p.expressionStatement(separators...)
}

// emptyStatement = SEPARATOR
func (p *Parser) emptyStatement() *ast.EmptyStatementNode {
	sepTok := p.advance()
	return ast.NewEmptyStatementNode(sepTok.Position)
}

const statementSeparatorMessage = "a statement separator `\\n`, `;`"

// expressionStatement = expressionWithModifier [SEPARATOR]
func (p *Parser) expressionStatement(separators ...token.Type) *ast.ExpressionStatementNode {
	expr := p.expressionWithModifier()
	var sep *token.Token
	if p.lookahead.IsStatementSeparator() {
		sep = p.advance()
		return ast.NewExpressionStatementNode(
			expr.Pos().Join(sep.Pos()),
			expr,
		)
	}
	for _, sepType := range separators {
		if p.lookahead.Type == sepType {
			return ast.NewExpressionStatementNode(
				expr.Pos(),
				expr,
			)
		}
	}

	if p.lookahead.Type == token.END_OF_FILE {
		return ast.NewExpressionStatementNode(
			position.New(
				expr.Pos().StartByte,
				p.lookahead.StartByte-expr.Pos().StartByte,
				expr.Pos().Line,
				expr.Pos().Column,
			),
			expr,
		)
	}

	p.errorExpected(statementSeparatorMessage)
	if p.synchronise() {
		p.advance()
	}

	return ast.NewExpressionStatementNode(
		expr.Pos(),
		expr,
	)
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

// parameter = identifier [":" typeAnnotation] ["=" expressionWithoutModifier]
func (p *Parser) parameter() ast.ParameterNode {
	var init ast.ExpressionNode
	var typ ast.TypeNode

	paramName, ok := p.matchOk(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER)
	if !ok {
		p.errorExpected("an identifier as the name of the declared parameter")
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

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		lastPos = init.Pos()
	}

	return ast.NewFormalParameterNode(
		paramName.Position.Join(lastPos.Pos()),
		paramName.Value,
		typ,
		init,
	)
}

// parameterList = parameter ("," parameter)* [","]
func (p *Parser) parameterList(stopTokens ...token.Type) []ast.ParameterNode {
	var args []ast.ParameterNode
	args = append(args, p.parameter())

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
		if !p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) {
			break
		}
		param := p.parameter()
		args = append(args, param)
	}

	return args
}

// positionalArguments = [expressionWithoutModifier ("," expressionWithoutModifier)* [","]]
func (p *Parser) positionalArguments(stopTokens ...token.Type) []ast.ExpressionNode {
	var args []ast.ExpressionNode

	for {
		if p.lookahead.Type == token.END_OF_FILE {
			return args
		}
		for _, stopToken := range stopTokens {
			if p.lookahead.Type == stopToken {
				return args
			}
		}
		expr := p.expressionWithoutModifier()
		args = append(args, expr)
	}
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

// powerExpression = constantLookup | constantLookup "**" powerExpression
func (p *Parser) powerExpression() ast.ExpressionNode {
	left := p.constantLookup()

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
	case token.PUBLIC_CONSTANT:
		tok := p.advance()
		return ast.NewPublicConstantNode(
			tok.Position,
			tok.Value,
		)
	case token.PRIVATE_CONSTANT:
		tok := p.advance()
		return ast.NewPrivateConstantNode(
			tok.Position,
			tok.Value,
		)
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
	case token.TYPEDEF:
		return p.typeDefinition()
	case token.ALIAS:
		return p.aliasExpression()
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

// methodDefinition = "def" METHOD_NAME ["(" parameterList ")"] [":" typeAnnotation] ["!" typeAnnotation] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
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

	// parameterList
	if p.match(token.LPAREN) {
		if !p.match(token.RPAREN) {
			params = p.parameterList(token.RPAREN)

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

	lastPos, body, multiline := p.statementBlockBodyWithThen(token.END)
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

// typeVariableList = typeVariable ("," typeVariable)* [","]
func (p *Parser) typeVariableList(stopTokens ...token.Type) []ast.TypeVariableNode {
	var vars []ast.TypeVariableNode
	vars = append(vars, p.typeVariable())

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
		if !p.accept(token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT, token.PLUS, token.MINUS) {
			break
		}
		typeVar := p.typeVariable()
		vars = append(vars, typeVar)
	}

	return vars
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

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(token.END)
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

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(token.END)
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

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(token.END)
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
	left := p.intersectionType()

	for p.lookahead.Type == token.OR {
		operator := p.advance()

		p.swallowEndLines()
		right := p.intersectionType()

		left = ast.NewBinaryTypeExpressionNode(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
}

// intersectionType = nilableType | intersectionType "&" nilableType
func (p *Parser) intersectionType() ast.TypeNode {
	left := p.nilableType()

	for p.lookahead.Type == token.AND {
		operator := p.advance()

		p.swallowEndLines()
		right := p.nilableType()

		left = ast.NewBinaryTypeExpressionNode(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
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

// genericConstant = strictConstantLookup | strictConstantLookup "[" [constantList] "]"
func (p *Parser) genericConstant() ast.ComplexConstantNode {
	constant := p.strictConstantLookup()
	if !p.match(token.LBRACKET) {
		return constant
	}

	if p.match(token.RBRACKET) {
		p.errorExpected("a constant")
		return constant
	}

	constList := p.constantList(token.RBRACKET)
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

// constantList = genericConstant ("," genericConstant)* [","]
func (p *Parser) constantList(stopTokens ...token.Type) []ast.ComplexConstantNode {
	var consts []ast.ComplexConstantNode
	consts = append(consts, p.genericConstant())

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
		consts = append(consts, p.genericConstant())
	}

	return consts
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

	lastPos, thenBody, multiline := p.statementBlockBody(token.END)
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

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(token.END)
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

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(token.END)
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

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(token.END, token.ELSE)
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
		lastPos, thenBody, multiline = p.statementBlockBody(token.END)
		currentExpr.ElseBody = thenBody
		if lastPos != nil {
			*currentExpr.Position = *currentExpr.Position.Join(lastPos)
		}
	} else if p.match(token.ELSE) {
		lastPos, thenBody, multiline = p.statementBlockBody(token.END)
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

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(token.END, token.ELSE, token.ELSIF)
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

		lastPos, thenBody, multiline = p.statementBlockBodyWithThen(token.END, token.ELSE, token.ELSIF)
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
		lastPos, thenBody, multiline = p.statementBlockBody(token.END)
		currentExpr.ElseBody = thenBody
		if lastPos != nil {
			*currentExpr.Position = *currentExpr.Position.Join(lastPos)
		}
	} else if p.match(token.ELSE) {
		lastPos, thenBody, multiline = p.statementBlockBody(token.END)
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
		body := p.statementsWithStop(token.RBRACE)
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

	lastPos, body, multiline := p.statementBlockBody(token.END)
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

// closureExpression = [(("|" parameterList "|") | "||") [: typeAnnotation]] closureAfterArrow
func (p *Parser) closureExpression() ast.ExpressionNode {
	var params []ast.ParameterNode
	var firstPos *position.Position
	var returnType ast.TypeNode

	if p.accept(token.OR) {
		firstPos = p.advance().Position
		if !p.accept(token.OR) {
			p.mode = withoutBitwiseOrMode
			params = p.parameterList(token.OR)
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
