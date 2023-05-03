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
	lookahead     *lexer.Token // next token used for predicting productions
	nextLookahead *lexer.Token // second next token used for predicting productions
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

// Adds an error which tells the user that another type of token
// was expected.
func (p *Parser) errorExpected(expected string) {
	p.errorMessage(fmt.Sprintf("unexpected %s, expected %s", p.lookahead.TokenType.String(), expected))
}

// Same as [errorExpected] but lets you pass a token type.
func (p *Parser) errorExpectedToken(expected lexer.TokenType) {
	p.errorExpected(expected.String())
}

// Adds an error with a custom message.
func (p *Parser) errorMessage(message string) {
	p.errorMessagePos(message, p.lookahead.Position)
}

// Same as [errorMessage] but let's you pass a Position.
func (p *Parser) errorMessagePos(message string, pos *lexer.Position) {
	if p.mode == panicMode {
		return
	}

	p.errors.Add(
		message,
		pos,
	)
}

// Add the content of an error token to the syntax error list.
func (p *Parser) errorToken(err *lexer.Token) {
	p.errorMessagePos(err.Value, err.Position)
}

// Attempt to consume the specified token type.
// If the next token doesn't match an error is added and the parser
// enters panic mode.
func (p *Parser) consume(tokenType lexer.TokenType) (*lexer.Token, bool) {
	if p.lookahead.TokenType == lexer.ErrorToken {
		return p.advance(), false
	}

	if p.lookahead.TokenType != tokenType {
		p.errorExpectedToken(tokenType)
		p.mode = panicMode
		return p.advance(), false
	}

	return p.advance(), true
}

// Checks if the next token matches any of the given types,
// if so it gets consumed.
func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, typ := range types {
		if p.accept(typ) {
			p.advance()
			return true
		}
	}

	return false
}

// Same as [match] but returns the consumed token.
func (p *Parser) matchOk(types ...lexer.TokenType) (*lexer.Token, bool) {
	for _, typ := range types {
		if p.accept(typ) {
			return p.advance(), true
		}
	}

	return nil, false
}

// Checks whether there are any more tokens to be consumed.
func (p *Parser) isAtEnd() bool {
	return p.lookahead.TokenType == lexer.EndOfFileToken
}

// Checks whether the next token matches the specified type.
func (p *Parser) accept(tokenType lexer.TokenType) bool {
	return p.lookahead.TokenType == tokenType
}

// Move over to the next token.
func (p *Parser) advance() *lexer.Token {
	previous := p.lookahead
	previousNext := p.nextLookahead
	if previousNext != nil && previousNext.TokenType == lexer.ErrorToken {
		p.errorToken(previousNext)
	}
	p.nextLookahead = p.lexer.Next()
	p.lookahead = previousNext
	return previous
}

// Consume statements until the provided token type is encountered.
func (p *Parser) statementsWithStop(stopTokens ...lexer.TokenType) []ast.StatementNode {
	var statementList []ast.StatementNode

	for {
		if p.lookahead.TokenType == lexer.EndOfFileToken {
			return statementList
		}
		for _, stopToken := range stopTokens {
			if p.lookahead.TokenType == stopToken {
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
		switch p.lookahead.TokenType {
		case lexer.EndOfFileToken:
			return false
		case lexer.SemicolonToken,
			lexer.EndLineToken:
			return true
		}

		p.advance()
	}
}

// Accept and ignore any number of consecutive end-line tokens.
func (p *Parser) swallowEndLines() {
	for {
		if !p.match(lexer.EndLineToken) {
			break
		}
	}
}

// Checks if the given slice of token types contains
// the given token type.
func containsToken(slice []lexer.TokenType, v lexer.TokenType) bool {
	for _, s := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// Consume a block of statements, like in `else` expressions,
// that terminates with `end`.
func (p *Parser) statementBlockBody(stopTokens ...lexer.TokenType) (*lexer.Position, []ast.StatementNode, bool) {
	var thenBody []ast.StatementNode
	var lastPos *lexer.Position
	var multiline bool

	if !p.lookahead.IsStatementSeparator() {
		expr := p.expressionWithoutModifier()
		thenBody = append(thenBody, &ast.ExpressionStatementNode{
			Position:   expr.Pos(),
			Expression: expr,
		})
		lastPos = expr.Pos()
	} else {
		multiline = true
		p.advance()

		if p.accept(lexer.EndToken) {
			lastPos = p.lookahead.Position
		} else if !containsToken(stopTokens, p.lookahead.TokenType) {
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
func (p *Parser) statementBlockBodyWithThen(stopTokens ...lexer.TokenType) (*lexer.Position, []ast.StatementNode, bool) {
	var thenBody []ast.StatementNode
	var lastPos *lexer.Position
	var multiline bool

	if p.lookahead.TokenType == lexer.ThenToken {
		p.advance()
		expr := p.expressionWithoutModifier()
		thenBody = append(thenBody, &ast.ExpressionStatementNode{
			Position:   expr.Pos(),
			Expression: expr,
		})
		lastPos = expr.Pos()
	} else {
		multiline = true
		if p.lookahead.IsStatementSeparator() {
			p.advance()
		} else {
			p.errorExpected(statementSeparatorMessage)
		}

		if p.accept(lexer.EndToken) {
			lastPos = p.lookahead.Position
		} else if !containsToken(stopTokens, p.lookahead.TokenType) {
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
	return &ast.ProgramNode{
		Position: &lexer.Position{
			StartByte:  0,
			ByteLength: len(p.source),
			Line:       1,
			Column:     1,
		},
		Body: statements,
	}
}

// statements = statement*
func (p *Parser) statements() []ast.StatementNode {
	return p.statementsWithStop(lexer.ZeroToken)
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
	return &ast.EmptyStatementNode{
		Position: sepTok.Position,
	}
}

const statementSeparatorMessage = "a statement separator `\\n`, `;` or end of file"

// expressionStatement = expressionWithModifier [SEPARATOR]
func (p *Parser) expressionStatement() *ast.ExpressionStatementNode {
	expr := p.expressionWithModifier()
	var sep *lexer.Token
	if p.lookahead.IsStatementSeparator() {
		sep = p.advance()
		return &ast.ExpressionStatementNode{
			Expression: expr,
			Position:   expr.Pos().Join(sep.Pos()),
		}
	}

	if p.lookahead.TokenType == lexer.EndOfFileToken {
		return &ast.ExpressionStatementNode{
			Expression: expr,
			Position: &lexer.Position{
				StartByte:  expr.Pos().StartByte,
				ByteLength: p.lookahead.StartByte - expr.Pos().StartByte,
				Line:       expr.Pos().Line,
				Column:     expr.Pos().Column,
			},
		}
	}

	p.errorExpected(statementSeparatorMessage)
	if p.synchronise() {
		p.advance()
	}

	return &ast.ExpressionStatementNode{
		Expression: expr,
		Position:   expr.Pos(),
	}
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

	switch p.lookahead.TokenType {
	case lexer.UnlessToken, lexer.WhileToken, lexer.UntilToken:
		mod := p.advance()
		right := p.expressionWithoutModifier()
		return &ast.ModifierNode{
			Position: left.Pos().Join(right.Pos()),
			Left:     left,
			Modifier: mod,
			Right:    right,
		}
	case lexer.IfToken:
		mod := p.advance()
		cond := p.expressionWithoutModifier()
		if p.lookahead.TokenType == lexer.ElseToken {
			p.advance()
			elseExpr := p.expressionWithoutModifier()
			return &ast.ModifierIfElseNode{
				Position:       left.Pos().Join(elseExpr.Pos()),
				ThenExpression: left,
				Condition:      cond,
				ElseExpression: elseExpr,
			}
		}
		return &ast.ModifierNode{
			Position: left.Pos().Join(cond.Pos()),
			Left:     left,
			Modifier: mod,
			Right:    cond,
		}
	}

	return left
}

// assignmentExpression = logicalOrNilCoalescingExpression | expression ASSIGN_OP assignmentExpression
func (p *Parser) assignmentExpression() ast.ExpressionNode {
	left := p.logicalOrNilCoalescingExpression()
	if p.lookahead.TokenType == lexer.ColonEqualToken {
		if !ast.IsValidDeclarationTarget(left) {
			p.errorMessagePos(
				fmt.Sprintf("invalid `%s` declaration target", p.lookahead.TokenType.String()),
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
			fmt.Sprintf("invalid `%s` assignment target", p.lookahead.TokenType.String()),
			left.Pos(),
		)
	}

	operator := p.advance()
	p.swallowEndLines()
	right := p.assignmentExpression()

	return &ast.AssignmentExpressionNode{
		Left:     left,
		Op:       operator,
		Right:    right,
		Position: left.Pos().Join(right.Pos()),
	}
}

// logicalOrNilCoalescingExpression = logicalAndExpression |
// logicalOrNilCoalescingExpression "||" logicalAndExpression |
// logicalOrNilCoalescingExpression "??" logicalAndExpression
func (p *Parser) logicalOrNilCoalescingExpression() ast.ExpressionNode {
	left := p.logicalAndExpression()

	for p.lookahead.TokenType == lexer.OrOrToken || p.lookahead.TokenType == lexer.QuestionQuestionToken {
		operator := p.advance()

		p.swallowEndLines()
		right := p.logicalAndExpression()

		left = &ast.LogicalExpressionNode{
			Op:       operator,
			Left:     left,
			Right:    right,
			Position: left.Pos().Join(right.Pos()),
		}
	}

	return left
}

// logicalAndExpression = equalityExpression | logicalAndExpression "&&" equalityExpression
func (p *Parser) logicalAndExpression() ast.ExpressionNode {
	left := p.equalityExpression()

	for p.lookahead.TokenType == lexer.AndAndToken {
		operator := p.advance()

		p.swallowEndLines()
		right := p.equalityExpression()

		left = &ast.LogicalExpressionNode{
			Op:       operator,
			Left:     left,
			Right:    right,
			Position: left.Pos().Join(right.Pos()),
		}
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

		left = &ast.BinaryExpressionNode{
			Op:       operator,
			Left:     left,
			Right:    right,
			Position: left.Pos().Join(right.Pos()),
		}
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

		left = &ast.BinaryExpressionNode{
			Op:       operator,
			Left:     left,
			Right:    right,
			Position: left.Pos().Join(right.Pos()),
		}
	}

	return left
}

// additiveExpression = multiplicativeExpression | additiveExpression ("+" | "-") multiplicativeExpression
func (p *Parser) additiveExpression() ast.ExpressionNode {
	left := p.multiplicativeExpression()

	for {
		operator, ok := p.matchOk(lexer.MinusToken, lexer.PlusToken)
		if !ok {
			break
		}
		p.swallowEndLines()
		right := p.multiplicativeExpression()
		left = &ast.BinaryExpressionNode{
			Op:       operator,
			Left:     left,
			Right:    right,
			Position: left.Pos().Join(right.Pos()),
		}
	}

	return left
}

// multiplicativeExpression = unaryExpression | multiplicativeExpression ("*" | "/") unaryExpression
func (p *Parser) multiplicativeExpression() ast.ExpressionNode {
	left := p.unaryExpression()

	for {
		operator, ok := p.matchOk(lexer.StarToken, lexer.SlashToken)
		if !ok {
			break
		}
		p.swallowEndLines()
		right := p.unaryExpression()
		left = &ast.BinaryExpressionNode{
			Op:       operator,
			Left:     left,
			Right:    right,
			Position: left.Pos().Join(right.Pos()),
		}
	}

	return left
}

// unaryExpression = powerExpression | ("!" | "-" | "+" | "~") unaryExpression
func (p *Parser) unaryExpression() ast.ExpressionNode {
	if operator, ok := p.matchOk(lexer.BangToken, lexer.MinusToken, lexer.PlusToken, lexer.TildeToken); ok {
		p.swallowEndLines()
		right := p.unaryExpression()
		return &ast.UnaryExpressionNode{
			Op:       operator,
			Right:    right,
			Position: operator.Pos().Join(right.Pos()),
		}
	}

	return p.powerExpression()
}

// powerExpression = primaryExpression | primaryExpression "**" powerExpression
func (p *Parser) powerExpression() ast.ExpressionNode {
	left := p.primaryExpression()

	if p.lookahead.TokenType != lexer.StarStarToken {
		return left
	}

	operator := p.advance()
	p.swallowEndLines()
	right := p.powerExpression()

	return &ast.BinaryExpressionNode{
		Op:       operator,
		Left:     left,
		Right:    right,
		Position: left.Pos().Join(right.Pos()),
	}
}

// primaryExpression = "true" |
// "false" |
// "nil" |
// "self" |
// INT |
// FLOAT |
// STRING |
// IDENT |
// PRIV_IDENT |
// CONST |
// PRIV_CONST |
// "(" expressionWithModifier ")"
func (p *Parser) primaryExpression() ast.ExpressionNode {
	switch p.lookahead.TokenType {
	case lexer.TrueToken:
		tok := p.advance()
		return &ast.TrueLiteralNode{Position: tok.Position}
	case lexer.FalseToken:
		tok := p.advance()
		return &ast.FalseLiteralNode{Position: tok.Position}
	case lexer.NilToken:
		tok := p.advance()
		return &ast.NilLiteralNode{Position: tok.Position}
	case lexer.SelfToken:
		tok := p.advance()
		return &ast.SelfLiteralNode{Position: tok.Position}
	case lexer.LParenToken:
		p.advance()
		expr := p.expressionWithModifier()
		p.consume(lexer.RParenToken)
		return expr
	case lexer.RawStringToken:
		tok := p.advance()
		return &ast.RawStringLiteralNode{
			Value:    tok.Value,
			Position: tok.Position,
		}
	case lexer.IfToken:
		return p.ifExpression()
	case lexer.UnlessToken:
		return p.unlessExpression()
	case lexer.WhileToken:
		return p.whileExpression()
	case lexer.UntilToken:
		return p.untilExpression()
	case lexer.StringBegToken:
		return p.stringLiteral()
	case lexer.IdentifierToken:
		tok := p.advance()
		return &ast.IdentifierNode{
			Position: tok.Position,
			Value:    tok.Value,
		}
	case lexer.PrivateIdentifierToken:
		tok := p.advance()
		return &ast.PrivateIdentifierNode{
			Position: tok.Position,
			Value:    tok.Value,
		}
	case lexer.ConstantToken:
		tok := p.advance()
		return &ast.ConstantNode{
			Position: tok.Position,
			Value:    tok.Value,
		}
	case lexer.PrivateConstantToken:
		tok := p.advance()
		return &ast.PrivateConstantNode{
			Position: tok.Position,
			Value:    tok.Value,
		}
	case lexer.HexIntToken, lexer.DuoIntToken, lexer.DecIntToken,
		lexer.OctIntToken, lexer.QuatIntToken, lexer.BinIntToken:
		tok := p.advance()
		return &ast.IntLiteralNode{
			Token:    tok,
			Position: tok.Position,
		}
	case lexer.FloatToken:
		tok := p.advance()
		return &ast.FloatLiteralNode{
			Value:    tok.Value,
			Position: tok.Position,
		}
	case lexer.ErrorToken:
		tok := p.advance()
		return &ast.InvalidNode{
			Token:    tok,
			Position: tok.Position,
		}
	default:
		p.errorExpected("an expression")
		p.mode = panicMode
		tok := p.advance()
		return &ast.InvalidNode{
			Token:    tok,
			Position: tok.Position,
		}
	}
}

// whileExpression = "while" expressionWithoutModifier ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) whileExpression() *ast.WhileExpressionNode {
	whileTok := p.advance()
	cond := p.expressionWithoutModifier()
	var pos *lexer.Position

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(lexer.EndToken)
	if lastPos != nil {
		pos = whileTok.Position.Join(lastPos)
	} else {
		pos = whileTok.Position
	}

	if multiline {
		endTok, ok := p.consume(lexer.EndToken)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return &ast.WhileExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  thenBody,
	}
}

// untilExpression = "until" expressionWithoutModifier ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) untilExpression() *ast.UntilExpressionNode {
	untilTok := p.advance()
	cond := p.expressionWithoutModifier()
	var pos *lexer.Position

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(lexer.EndToken)
	if lastPos != nil {
		pos = untilTok.Position.Join(lastPos)
	} else {
		pos = untilTok.Position
	}

	if multiline {
		endTok, ok := p.consume(lexer.EndToken)
		if ok {
			pos = pos.Join(endTok.Position)
		}
	}

	return &ast.UntilExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  thenBody,
	}
}

// unlessExpression = "unless" expressionWithoutModifier ((SEPARATOR [statements]) | ("then" expressionWithoutModifier))
// ["else" ((SEPARATOR [statements]) | expressionWithoutModifier)]
// "end"
func (p *Parser) unlessExpression() *ast.UnlessExpressionNode {
	unlessTok := p.advance()
	cond := p.expressionWithoutModifier()
	var pos *lexer.Position

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(lexer.EndToken, lexer.ElseToken)
	if lastPos != nil {
		pos = unlessTok.Position.Join(lastPos)
	} else {
		pos = unlessTok.Position
	}

	unlessExpr := &ast.UnlessExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  thenBody,
	}
	currentExpr := unlessExpr

	if p.lookahead.IsStatementSeparator() && p.nextLookahead.TokenType == lexer.ElseToken {
		p.advance()
		p.advance()
		lastPos, thenBody, multiline = p.statementBlockBody(lexer.EndToken)
		currentExpr.ElseBody = thenBody
		if lastPos != nil {
			*currentExpr.Position = *currentExpr.Position.Join(lastPos)
		}
	} else if p.match(lexer.ElseToken) {
		lastPos, thenBody, multiline = p.statementBlockBody(lexer.EndToken)
		currentExpr.ElseBody = thenBody
		if lastPos != nil {
			*currentExpr.Position = *currentExpr.Position.Join(lastPos)
		}
	}

	if multiline {
		endTok, ok := p.consume(lexer.EndToken)
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
	var pos *lexer.Position

	lastPos, thenBody, multiline := p.statementBlockBodyWithThen(lexer.EndToken, lexer.ElseToken, lexer.ElsifToken)
	if lastPos != nil {
		pos = ifTok.Position.Join(lastPos)
	} else {
		pos = ifTok.Position
	}

	ifExpr := &ast.IfExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  thenBody,
	}
	currentExpr := ifExpr

	for {
		var elsifTok *lexer.Token

		if p.lookahead.TokenType == lexer.ElsifToken {
			elsifTok = p.advance()
		} else if p.lookahead.IsStatementSeparator() && p.nextLookahead.TokenType == lexer.ElsifToken {
			p.advance()
			elsifTok = p.advance()
		} else {
			break
		}
		cond = p.expressionWithoutModifier()

		lastPos, thenBody, multiline = p.statementBlockBodyWithThen(lexer.EndToken, lexer.ElseToken, lexer.ElsifToken)
		if lastPos != nil {
			pos = elsifTok.Position.Join(lastPos)
		} else {
			pos = elsifTok.Position
		}

		elsifExpr := &ast.IfExpressionNode{
			Position:  pos,
			Condition: cond,
			ThenBody:  thenBody,
		}

		currentExpr.ElseBody = []ast.StatementNode{
			&ast.ExpressionStatementNode{
				Position:   elsifExpr.Position,
				Expression: elsifExpr,
			},
		}
		currentExpr = elsifExpr
	}

	if p.lookahead.IsStatementSeparator() && p.nextLookahead.TokenType == lexer.ElseToken {
		p.advance()
		p.advance()
		lastPos, thenBody, multiline = p.statementBlockBody(lexer.EndToken)
		currentExpr.ElseBody = thenBody
		if lastPos != nil {
			*currentExpr.Position = *currentExpr.Position.Join(lastPos)
		}
	} else if p.match(lexer.ElseToken) {
		lastPos, thenBody, multiline = p.statementBlockBody(lexer.EndToken)
		currentExpr.ElseBody = thenBody
		if lastPos != nil {
			*currentExpr.Position = *currentExpr.Position.Join(lastPos)
		}
	}

	if multiline {
		endTok, ok := p.consume(lexer.EndToken)
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
	var quoteEnd *lexer.Token

	var strContent []ast.StringLiteralContentNode
	for {
		if tok, ok := p.matchOk(lexer.StringContentToken); ok {
			strContent = append(strContent, &ast.StringLiteralContentSectionNode{
				Value:    tok.Value,
				Position: tok.Position,
			})
			continue
		}

		if beg, ok := p.matchOk(lexer.StringInterpBegToken); ok {
			expr := p.expressionWithoutModifier()
			end, _ := p.consume(lexer.StringInterpEndToken)
			strContent = append(strContent, &ast.StringInterpolationNode{
				Expression: expr,
				Position:   beg.Position.Join(end.Position),
			})
			continue
		}

		tok, ok := p.consume(lexer.StringEndToken)
		quoteEnd = tok
		if tok.TokenType == lexer.EndOfFileToken {
			break
		}
		if !ok {
			strContent = append(strContent, &ast.InvalidNode{
				Token:    tok,
				Position: tok.Position,
			})
			continue
		}
		break
	}

	return &ast.StringLiteralNode{
		Content:  strContent,
		Position: quoteBeg.Position.Join(quoteEnd.Position),
	}
}
