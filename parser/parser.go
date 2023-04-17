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
	panicMode              // triggered after encountering a syntax error, changes to `errorMode` after synchronisation
	errorMode              // like normalMode, but errors have been detected
)

// Holds the current state of the parsing process.
type Parser struct {
	source    []byte       // Elk source code
	lexer     *lexer.Lexer // lexer which outputs a stream of tokens
	lookahead *lexer.Token // next token used for predicting productions
	errors    ErrorList
	mode      mode
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
	p.advance()
	return p.program(), p.errors
}

// Adds an error which tells the user that another type of token
// was expected.
func (p *Parser) errorExpected(expected lexer.TokenType) {
	p.mode = panicMode
	p.errors.Add(
		fmt.Sprintf("expected `%s`, got `%s`", expected.String(), p.lookahead.TokenType.String()),
		p.lookahead.Position,
	)
}

// Adds an error with a custom message.
func (p *Parser) errorMessage(message string) {
	p.mode = panicMode
	p.errors.Add(
		message,
		p.lookahead.Position,
	)
}

// Add the content of an error token to the syntax error list/
func (p *Parser) errorToken(err *lexer.Token) {
	p.mode = panicMode
	p.errors.Add(
		err.Value,
		err.Position,
	)
}

// Attempt to consume the specified token type.
// If the next token doesn't match an error is added an the parser
// enters panic mode.
func (p *Parser) consume(tokenType lexer.TokenType) (*lexer.Token, bool) {
	if p.lookahead.TokenType == lexer.ErrorToken {
		return p.advance(), false
	}

	if p.lookahead.TokenType != tokenType {
		p.errorExpected(tokenType)
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
	return p.lookahead.TokenType == lexer.EOFToken
}

// Checks whether the next token matches the specified type.
func (p *Parser) accept(tokenType lexer.TokenType) bool {
	return p.lookahead.TokenType == tokenType
}

// Move over to the next token.
func (p *Parser) advance() *lexer.Token {
	previous := p.lookahead
	if previous != nil && previous.TokenType == lexer.ErrorToken {
		p.errorToken(previous)
	}
	p.lookahead = p.lexer.Next()
	return previous
}

// Consume statements until the provided token type is encountered.
func (p *Parser) statementsWithStop(stopLookahead lexer.TokenType) []ast.StatementNode {
	statementList := []ast.StatementNode{p.statement()}

	for p.lookahead.TokenType != lexer.EOFToken && p.lookahead.TokenType != stopLookahead {
		statement := p.statement()
		statementList = append(statementList, statement)
	}

	return statementList
}

// Discards tokens until something resembling a new statement is encountered.
// Used for recovering after errors.
func (p *Parser) synchronise() {
	p.mode = errorMode

	for {
		switch p.lookahead.TokenType {
		case lexer.EOFToken:
			return
		case lexer.SemicolonToken:
			p.advance()
			return
		case lexer.EndLineToken:
			p.advance()
			return
		case lexer.UsingToken:
			return
		}

		p.advance()
	}
}

// ==== Productions ====

// program = statements
func (p *Parser) program() *ast.ProgramNode {
	statements := p.statements()
	return &ast.ProgramNode{
		Position: lexer.Position{
			StartByte:  0,
			ByteLength: len(p.source),
			Line:       1,
			Column:     1,
		},
		Body: statements,
	}
}

// statements = statement+
func (p *Parser) statements() []ast.StatementNode {
	return p.statementsWithStop(lexer.ZeroToken)
}

// statement = expressionStatement
func (p *Parser) statement() ast.StatementNode {
	stmt := p.expressionStatement()
	if p.mode == panicMode {
		p.synchronise()
	}
	return stmt
}

// expressionStatement = expression [SEPARATOR]
func (p *Parser) expressionStatement() *ast.ExpressionStatementNode {
	expr := p.expression()
	pos := expr.Pos()
	var sep *lexer.Token
	if p.lookahead.IsStatementSeparator() {
		sep = p.advance()
		pos = pos.Join(sep.Pos())
	}

	return &ast.ExpressionStatementNode{
		Expression: expr,
		Position:   pos,
	}
}

// expression = assignmentExpression
func (p *Parser) expression() ast.ExpressionNode {
	return p.assignmentExpression()
}

// assignmentExpression = logicalOrExpression | expression ASSIGN_OP assignmentExpression
func (p *Parser) assignmentExpression() ast.ExpressionNode {
	left := p.logicalOrExpression()
	if !p.lookahead.IsAssignmentOperator() {
		return left
	}

	if !ast.IsValidLeftValue(left) {
		p.errorMessage(fmt.Sprintf("invalid left value in assignment `%s`", p.lookahead.TokenType.String()))
	}

	operator := p.advance()
	right := p.assignmentExpression()

	return &ast.AssignmentExpressionNode{
		Left:     left,
		Op:       operator,
		Right:    right,
		Position: left.Pos().Join(right.Pos()),
	}
}

// logicalOrExpression = logicalAndExpression | logicalOrExpression "||" logicalAndExpression
func (p *Parser) logicalOrExpression() ast.ExpressionNode {
	left := p.logicalAndExpression()

	for p.lookahead.TokenType == lexer.OrOrToken {
		operator := p.advance()

		right := p.logicalAndExpression()

		left = &ast.BinaryExpressionNode{
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

		right := p.equalityExpression()

		left = &ast.BinaryExpressionNode{
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

// multiplicativeExpression = powerExpression | multiplicativeExpression ("*" | "/") powerExpression
func (p *Parser) multiplicativeExpression() ast.ExpressionNode {
	left := p.powerExpression()

	for {
		operator, ok := p.matchOk(lexer.StarToken, lexer.SlashToken)
		if !ok {
			break
		}
		right := p.powerExpression()
		left = &ast.BinaryExpressionNode{
			Op:       operator,
			Left:     left,
			Right:    right,
			Position: left.Pos().Join(right.Pos()),
		}
	}

	return left
}

// powerExpression = unaryExpression | unaryExpression "**" powerExpression
func (p *Parser) powerExpression() ast.ExpressionNode {
	left := p.unaryExpression()

	if p.lookahead.TokenType != lexer.PowerToken {
		return left
	}

	operator := p.advance()
	right := p.powerExpression()

	return &ast.BinaryExpressionNode{
		Op:       operator,
		Left:     left,
		Right:    right,
		Position: left.Pos().Join(right.Pos()),
	}
}

// unaryExpression = primaryExpression | ("!" | "-" | "+" | "~") unaryExpression
func (p *Parser) unaryExpression() ast.ExpressionNode {
	if operator, ok := p.matchOk(lexer.BangToken, lexer.MinusToken, lexer.PlusToken, lexer.TildeToken); ok {
		right := p.unaryExpression()
		return &ast.UnaryExpressionNode{
			Op:       operator,
			Right:    right,
			Position: operator.Pos().Join(right.Pos()),
		}
	}

	return p.primaryExpression()
}

// primaryExpression = "true" | "false" | "nil" | INT | FLOAT | STRING | "(" expression ")"
func (p *Parser) primaryExpression() ast.ExpressionNode {
	if tok, ok := p.matchOk(lexer.TrueToken); ok {
		return &ast.TrueLiteralNode{Position: tok.Position}
	}
	if tok, ok := p.matchOk(lexer.FalseToken); ok {
		return &ast.FalseLiteralNode{Position: tok.Position}
	}
	if tok, ok := p.matchOk(lexer.NilToken); ok {
		return &ast.NilLiteralNode{Position: tok.Position}
	}

	if p.match(lexer.LParenToken) {
		expr := p.expression()
		p.consume(lexer.RParenToken)
		return expr
	}

	if tok, ok := p.matchOk(lexer.RawStringToken); ok {
		return &ast.RawStringLiteralNode{
			Value:    tok.Value,
			Position: tok.Position,
		}
	}

	if p.accept(lexer.StringBegToken) {
		return p.stringLiteral()
	}

	if p.lookahead.IsIntLiteral() {
		tok := p.advance()
		return &ast.IntLiteralNode{
			Token:    tok,
			Position: tok.Position,
		}
	}

	if tok, ok := p.matchOk(lexer.FloatToken); ok {
		return &ast.FloatLiteralNode{
			Value:    tok.Value,
			Position: tok.Position,
		}
	}

	p.errorMessage(fmt.Sprintf("expected an expression, got `%s`", p.lookahead.TokenType.String()))
	tok := p.advance()
	return &ast.InvalidNode{
		Token:    tok,
		Position: tok.Position,
	}
}

// stringLiteral = "\"" (STRING_CONTENT | "${" expression "}")* "\""
func (p *Parser) stringLiteral() ast.ExpressionNode {
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
			expr := p.expression()
			end, _ := p.consume(lexer.StringInterpEndToken)
			strContent = append(strContent, &ast.StringInterpolationNode{
				Expression: expr,
				Position:   beg.Position.Join(end.Position),
			})
			continue
		}

		tok, ok := p.consume(lexer.StringEndToken)
		quoteEnd = tok
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
