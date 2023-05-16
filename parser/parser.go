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
	normalMode mode = iota // regular parsing mode
	panicMode              // triggered after encountering a syntax error, changes to `normalMode` after synchronisation
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
		statement := p.statement()
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
func (p *Parser) statement() ast.StatementNode {
	if p.lookahead.IsStatementSeparator() {
		return p.emptyStatement()
	}

	return p.expressionStatement()
}

// emptyStatement = SEPARATOR
func (p *Parser) emptyStatement() *ast.EmptyStatementNode {
	sepTok := p.advance()
	return ast.NewEmptyStatementNode(sepTok.Position)
}

const statementSeparatorMessage = "a statement separator `\\n`, `;` or end of file"

// expressionStatement = expressionWithModifier [SEPARATOR]
func (p *Parser) expressionStatement() *ast.ExpressionStatementNode {
	expr := p.expressionWithModifier()
	var sep *token.Token
	if p.lookahead.IsStatementSeparator() {
		sep = p.advance()
		return ast.NewExpressionStatementNode(
			expr.Pos().Join(sep.Pos()),
			expr,
		)
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
func (p *Parser) parameter(stopTokens ...token.Type) ast.ParameterNode {
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
		paramName,
		typ,
		init,
	)
}

// parameters = parameter ("," parameter)* [","]
func (p *Parser) parameters(stopTokens ...token.Type) []ast.ParameterNode {
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
	left := p.logicalAndExpression()

FOR:
	for {
		switch p.lookahead.Type {
		case token.OR_OR, token.QUESTION_QUESTION, token.OR_BANG:
		default:
			break FOR
		}
		operator := p.advance()

		p.swallowEndLines()
		right := p.logicalAndExpression()

		left = ast.NewLogicalExpressionNode(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
}

// logicalAndExpression = equalityExpression |
// logicalAndExpression "&&" equalityExpression |
// logicalAndExpression "&!" equalityExpression
func (p *Parser) logicalAndExpression() ast.ExpressionNode {
	left := p.equalityExpression()

FOR:
	for {
		switch p.lookahead.Type {
		case token.AND_AND, token.AND_BANG:
		default:
			break FOR
		}
		operator := p.advance()

		p.swallowEndLines()
		right := p.equalityExpression()

		left = ast.NewLogicalExpressionNode(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
}

// equalityExpression = comparison | equalityExpression EQUALITY_OP comparison
func (p *Parser) equalityExpression() ast.ExpressionNode {
	left := p.comparison()

	for p.lookahead.IsEqualityOperator() {
		operator := p.advance()

		p.swallowEndLines()
		right := p.comparison()

		left = ast.NewBinaryExpressionNode(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
}

// comparison = additiveExpression | comparison COMP_OP additiveExpression
func (p *Parser) comparison() ast.ExpressionNode {
	left := p.additiveExpression()

	for p.lookahead.IsComparisonOperator() {
		operator := p.advance()

		p.swallowEndLines()
		right := p.additiveExpression()

		left = ast.NewBinaryExpressionNode(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
}

// additiveExpression = multiplicativeExpression | additiveExpression ("+" | "-") multiplicativeExpression
func (p *Parser) additiveExpression() ast.ExpressionNode {
	left := p.multiplicativeExpression()

	for {
		operator, ok := p.matchOk(token.MINUS, token.PLUS)
		if !ok {
			break
		}
		p.swallowEndLines()
		right := p.multiplicativeExpression()
		left = ast.NewBinaryExpressionNode(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
}

// multiplicativeExpression = unaryExpression | multiplicativeExpression ("*" | "/") unaryExpression
func (p *Parser) multiplicativeExpression() ast.ExpressionNode {
	left := p.unaryExpression()

	for {
		operator, ok := p.matchOk(token.STAR, token.SLASH)
		if !ok {
			break
		}
		p.swallowEndLines()
		right := p.unaryExpression()
		left = ast.NewBinaryExpressionNode(
			left.Pos().Join(right.Pos()),
			operator,
			left,
			right,
		)
	}

	return left
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

// constantLookup = primaryExpression | constantLookup "::" publicConstant
func (p *Parser) constantLookup() ast.ExpressionNode {
	left := p.primaryExpression()

	for p.lookahead.Type == token.SCOPE_RES_OP {
		p.advance()

		p.swallowEndLines()
		if p.accept(token.PRIVATE_CONSTANT) {
			p.errorUnexpected("can't access a private constant from the outside")
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
		expr := p.expressionWithModifier()
		p.consume(token.RPAREN)
		return expr
	case token.RAW_STRING:
		tok := p.advance()
		return ast.NewRawStringLiteralNode(
			tok.Position,
			tok.Value,
		)
	case token.OR, token.THIN_ARROW:
		return p.closureExpression()
	case token.VAR:
		return p.variableDeclaration()
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
	case token.STRING_BEG:
		return p.stringLiteral()
	case token.PUBLIC_IDENTIFIER:
		tok := p.advance()
		return ast.NewPublicIdentifierNode(
			tok.Position,
			tok.Value,
		)
	case token.PRIVATE_IDENTIFIER:
		tok := p.advance()
		return ast.NewPrivateIdentifierNode(
			tok.Position,
			tok.Value,
		)
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

// variableDeclaration = "var" identifier [":" typeAnnotation] ["=" expressionWithoutModifier]
func (p *Parser) variableDeclaration() ast.ExpressionNode {
	varTok := p.advance()
	var init ast.ExpressionNode
	var typ ast.TypeNode

	varName, ok := p.matchOk(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER)
	lastPos := varName.Position
	if !ok {
		p.errorExpected("an identifier as the name of the declared variable")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Position,
			tok,
		)
	}

	if p.match(token.COLON) {
		typ = p.typeAnnotation()
		lastPos = typ.Pos()
	}

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		lastPos = init.Pos()
	}

	return ast.NewVariableDeclarationNode(
		varTok.Position.Join(lastPos.Pos()),
		varName,
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

// namedType = typeConstantLookup
func (p *Parser) namedType() ast.TypeNode {
	return p.typeConstantLookup()
}

// typeConstantLookup = constant | typeConstantLookup "::" publicConstant
func (p *Parser) typeConstantLookup() ast.ConstantNode {
	var left ast.ConstantNode
	left = p.constant()

	for p.lookahead.Type == token.SCOPE_RES_OP {
		p.advance()

		p.swallowEndLines()
		if p.accept(token.PRIVATE_CONSTANT) {
			p.errorUnexpected("can't access a private constant from the outside")
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
	return ast.NewInvalidNode(
		tok.Position,
		tok,
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

// closureExpression = logicalOrExpression |
// [(("|" closureArguments "|") | "||") [: typeAnnotation]] "->" (expressionWithoutModifier | SEPARATOR [statements] "end")
func (p *Parser) closureExpression() ast.ExpressionNode {
	var params []ast.ParameterNode
	var firstPos *position.Position
	var pos *position.Position
	var returnType ast.TypeNode

	if p.accept(token.OR) {
		firstPos = p.advance().Position
		if !p.accept(token.OR) {
			params = p.parameters(token.OR)
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
	}

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
