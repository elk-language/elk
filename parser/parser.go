// Package parser implements a parser
// used by the Elk interpreter.
//
// Parser expects a slice of bytes containing Elk source code
// parses it, registering any encountered errors, and returns an Abstract Syntax Tree.
package parser

import (
	"fmt"
	"slices"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Parsing mode.
type mode uint8

const (
	zeroMode             mode = iota // initial zero value mode
	normalMode                       // regular parsing mode
	panicMode                        // triggered after encountering a syntax error, changes to `normalMode` after synchronisation
	withoutBitwiseOrMode             // disables bitwise OR `|` from the grammar
	incompleteMode                   // the input is incomplete, parser expected more tokens but received an END_OF_FILE.
	withoutUnionTypeMode             // disables union type `|` from the grammar
)

// Holds the current state of the parsing process.
type Parser struct {
	sourceName       string       // Path to the source file or some name.
	source           string       // Elk source code
	lexer            *lexer.Lexer // lexer which outputs a stream of tokens
	lookahead        *token.Token // next token used for predicting productions
	secondLookahead  *token.Token // second next token used for predicting productions
	thirdLookahead   *token.Token // third next token used for predicting productions
	diagnostics      diagnostic.DiagnosticList
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
func Parse(sourceName string, source string) (*ast.ProgramNode, diagnostic.DiagnosticList) {
	return New(sourceName, source).Parse()
}

func (*Parser) Class() *value.Class {
	return value.ElkParserClass
}

func (*Parser) DirectClass() *value.Class {
	return value.ElkParserClass
}

func (p *Parser) Inspect() string {
	return fmt.Sprintf("Std::Elk::Parser{&: %p}", p)
}

func (p *Parser) Error() string {
	return p.Inspect()
}

func (p *Parser) SingletonClass() *value.Class {
	return nil
}

func (p *Parser) InstanceVariables() value.SymbolMap {
	return nil
}

func (p *Parser) Copy() value.Reference {
	return p
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
func (p *Parser) Parse() (*ast.ProgramNode, diagnostic.DiagnosticList) {
	p.reset()

	p.advance() // populate thirdLookahead
	p.advance() // populate secondLookahead
	p.advance() // populate lookahead
	return p.program(), p.diagnostics
}

func (p *Parser) reset() {
	p.lexer = lexer.NewWithName(p.sourceName, p.source)
	p.mode = normalMode
	p.diagnostics = nil
}

// Adds an error which tells the user that the received
// token is unexpected.
func (p *Parser) errorUnexpected(message string) {
	p.errorMessage(fmt.Sprintf("unexpected %s, %s", p.lookahead.Type.Name(), message))
}

// Adds an error which tells the user that another type of token
// was expected.
func (p *Parser) errorExpected(expected string) {
	p.errorMessage(fmt.Sprintf("unexpected %s, expected %s", p.lookahead.Type.Name(), expected))
}

// Same as [errorExpected] but lets you pass a token type.
func (p *Parser) errorExpectedToken(expected token.Type) {
	p.errorExpected(expected.Name())
}

// Adds an error with a custom message.
func (p *Parser) errorMessage(message string) {
	p.errorMessageLocation(message, p.lookahead.Location())
}

// Same as [errorMessage] but let's you pass a Span.
func (p *Parser) errorMessageLocation(message string, loc *position.Location) {
	if p.mode == panicMode {
		return
	}

	p.diagnostics.AddFailure(
		message,
		loc,
	)
}

// Add the content of an error token to the syntax error list.
func (p *Parser) errorToken(err *token.Token) {
	p.errorMessageLocation(err.Value, err.Location())
}

// Attempt to consume the specified token type.
// If the next token doesn't match an error is added and the parser
// enters panic mode.
func (p *Parser) consume(tokenType token.Type) (*token.Token, bool) {
	return p.consumeExpected(tokenType, tokenType.Name())
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
func (p *Parser) acceptSecond(tokenTypes ...token.Type) bool {
	for _, typ := range tokenTypes {
		if p.secondLookahead.Type == typ {
			return true
		}
	}
	return false
}

// Checks whether the third next token matches any the specified types.
func (p *Parser) acceptThird(tokenTypes ...token.Type) bool {
	for _, typ := range tokenTypes {
		if p.thirdLookahead.Type == typ {
			return true
		}
	}
	return false
}

// Move over to the next token.
func (p *Parser) advance() *token.Token {
	previous := p.lookahead
	previousSecond := p.secondLookahead
	previousThird := p.thirdLookahead
	if previousSecond != nil && previousSecond.Type == token.ERROR {
		p.errorToken(previousSecond)
	}
	p.thirdLookahead = p.lexer.Next()
	p.secondLookahead = previousThird
	p.lookahead = previousSecond
	return previous
}

// Discards tokens until something resembling a new statement is encountered.
// Used for recovering after error.
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
func (p *Parser) statementBlock(stopTokens ...token.Type) (*position.Location, []ast.StatementNode, bool) {
	var thenBody []ast.StatementNode
	var lastLocation *position.Location
	var multiline bool

	if p.lookahead.Type == token.END_OF_FILE {
		p.errorExpected("a statement separator or an expression")
		p.updateErrorMode(true)
		return p.lookahead.Location(), nil, true
	}

	if !p.lookahead.IsStatementSeparator() {
		expr := p.expressionWithoutModifier()
		thenBody = append(thenBody, ast.NewExpressionStatementNode(
			expr.Location(),
			expr,
		))
		lastLocation = expr.Location()
	} else {
		multiline = true
		p.advance()

		if p.accept(token.END) {
			lastLocation = p.lookahead.Location()
		} else if !containsToken(stopTokens, p.lookahead.Type) {
			thenBody = p.statements(stopTokens...)
			if len(thenBody) > 0 {
				lastLocation = position.LocationOfLastElement(thenBody)
			}
		}
	}

	return lastLocation, thenBody, multiline
}

// statementProduction = subProduction [SEPARATOR]
func statementProduction[Expression, Statement ast.Node](p *Parser, constructor statementConstructor[Expression, Statement], expressionProduction func() Expression, separators ...token.Type) Statement {
	expr := expressionProduction()
	var sep *token.Token
	if p.lookahead.IsStatementSeparator() || p.lookahead.Type == token.END_OF_FILE {
		sep = p.advance()
		return constructor(
			expr.Location().Join(sep.Location()),
			expr,
		)
	}
	for _, sepType := range separators {
		if p.lookahead.Type == sepType {
			return constructor(
				expr.Location(),
				expr,
			)
		}
	}
	if p.match(token.ERROR) {
		if p.synchronise() {
			p.advance()
		}
		return constructor(
			expr.Location(),
			expr,
		)
	}

	p.updateErrorMode(false)
	p.errorExpected(statementSeparatorMessage)
	if p.synchronise() {
		p.advance()
	}

	return constructor(
		expr.Location(),
		expr,
	)
}

type statementsProduction[Statement ast.Node] func(...token.Type) []Statement

// Represents an AST Node constructor function for a new ast.StatementNode
type statementConstructor[Expression, Statement ast.Node] func(*position.Location, Expression) Statement

// Consume a block of statements, like in `if`, `elsif` or `while` expressions,
// that terminates with `end` or can be single-line when it begins with `then`.
func genericStatementBlockWithThen[Expression, Statement ast.Node](
	p *Parser,
	statementsProduction statementsProduction[Statement],
	expressionProduction func() Expression,
	statementConstructor statementConstructor[Expression, Statement],
	stopTokens ...token.Type,
) (*position.Location, []Statement, bool) {
	var thenBody []Statement
	var lastLocation *position.Location
	var multiline bool

	if p.lookahead.Type == token.THEN {
		p.advance()
		expr := expressionProduction()
		thenBody = append(thenBody, statementConstructor(
			expr.Location(),
			expr,
		))
		lastLocation = expr.Location()
	} else {
		multiline = true
		if p.lookahead.IsStatementSeparator() {
			p.advance()
		} else {
			p.updateErrorMode(false)
			p.errorExpected(statementSeparatorMessage)
		}

		if p.accept(token.END) {
			lastLocation = p.lookahead.Location()
		} else if !containsToken(stopTokens, p.lookahead.Type) {
			thenBody = statementsProduction(stopTokens...)
			if len(thenBody) > 0 {
				lastLocation = thenBody[len(thenBody)-1].Location()
			}
		}
	}

	return lastLocation, thenBody, multiline
}

// Consume a block of statements, like in `if`, `elsif` or `while` expressions,
// that terminates with `end` or can be single-line when it begins with `then`.
func (p *Parser) statementBlockWithThen(stopTokens ...token.Type) (*position.Location, []ast.StatementNode, bool) {
	return genericStatementBlockWithThen(p, p.statements, p.expressionWithoutModifier, ast.NewExpressionStatementNodeI, stopTokens...)
}

// Consume a block of statements, like in `if`, `elsif` or `while` expressions,
// that terminates with `end` or can be single-line when it begins with `then`.
func (p *Parser) structBodyStatementBlockWithThen(stopTokens ...token.Type) (*position.Location, []ast.StructBodyStatementNode, bool) {
	return genericStatementBlockWithThen(p, p.structBodyStatements, p.attributeParameter, ast.NewParameterStatementNodeI, stopTokens...)
}

// Represents an AST Node constructor function for binary operators
type binaryConstructor[Element ast.Node] func(*position.Location, *token.Token, Element, Element) Element

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
			left.Location().Join(right.Location()),
			operator,
			left,
			right,
		)
	}

	return left
}

// binaryExpression = subProduction | binaryExpression operators subProduction
func (p *Parser) binaryExpression(subProduction func() ast.ExpressionNode, operators ...token.Type) ast.ExpressionNode {
	return binaryProduction(p, ast.NewBinaryExpressionNodeI, subProduction, operators...)
}

// binaryTypeExpression = subProduction | binaryTypeExpression operators subProduction
func (p *Parser) binaryTypeExpression(subProduction func() ast.TypeNode, operators ...token.Type) ast.TypeNode {
	return binaryProduction(p, ast.NewBinaryTypeNodeI, subProduction, operators...)
}

// binaryPattern = subProduction | binaryPattern operators subProduction
func (p *Parser) binaryPattern(subProduction func() ast.PatternNode, operators ...token.Type) ast.PatternNode {
	return binaryProduction(p, ast.NewBinaryPatternNodeI, subProduction, operators...)
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
		if slices.Contains(stopTokens, p.lookahead.Type) {
			break
		}
		if !p.match(token.COMMA) {
			break
		}
		p.swallowNewlines()
		if slices.Contains(stopTokens, p.lookahead.Type) {
			break
		}
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
		if slices.Contains(stopTokens, p.lookahead.Type) {
			break
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
		if slices.Contains(stopTokens, p.lookahead.Type) {
			return list
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
		if slices.Contains(stopTokens, p.lookahead.Type) {
			return list
		}
		element := subProduction(stopTokens...)
		list = append(list, element)
	}
}

// ==== Productions ====

// program = topLevelStatements
func (p *Parser) program() *ast.ProgramNode {
	statements := p.topLevelStatements()
	lastLocation := position.LocationOfLastElement(statements)
	return ast.NewProgramNode(
		p.newLocation(position.NewSpanFromPosition(position.New(0, 1, 1))).Join(lastLocation),
		statements,
	)
}

// topLevelStatements = topLevelStatement*
func (p *Parser) topLevelStatements(stopTokens ...token.Type) []ast.StatementNode {
	return repeatedProductionWithStop(p, p.topLevelStatement, stopTokens...)
}

// topLevelStatement = emptyStatement | expressionStatement | importStatement
func (p *Parser) topLevelStatement(separators ...token.Type) ast.StatementNode {
	if p.lookahead.IsStatementSeparator() {
		return p.emptyStatement()
	}

	if p.accept(token.IMPORT) {
		return p.importStatement()
	}

	return p.expressionStatement(separators...)
}

func (p *Parser) importStatement() ast.StatementNode {
	importTok := p.advance()

	var stringLiteral ast.StringLiteralNode
	switch p.lookahead.Type {
	case token.STRING_BEG:
		stringLiteral = p.stringLiteral(false)
	case token.RAW_STRING:
		stringLiteral = p.rawStringLiteral()
	default:
		p.errorExpected("a string literal")
		errTok := p.advance()
		return ast.NewInvalidNode(errTok.Location(), errTok)
	}

	return ast.NewImportStatementNode(
		importTok.Location().Join(stringLiteral.Location()),
		stringLiteral,
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
	return statementProduction(p, ast.NewParameterStatementNode, p.attributeParameter, separators...)
}

// emptyStatement = SEPARATOR
func (p *Parser) emptyStatement() *ast.EmptyStatementNode {
	sepTok := p.advance()
	return ast.NewEmptyStatementNode(sepTok.Location())
}

const statementSeparatorMessage = "a statement separator `\\n`, `;`"

// expressionStatement = topLevelExpression [SEPARATOR]
func (p *Parser) expressionStatement(separators ...token.Type) *ast.ExpressionStatementNode {
	return statementProduction(p, ast.NewExpressionStatementNode, p.topLevelExpression, separators...)
}

// topLevelExpression = declarationExpression
func (p *Parser) topLevelExpression() ast.ExpressionNode {
	expr := p.declarationExpression(true)
	if p.mode == panicMode {
		p.synchronise()
	}
	return expr
}

func (p *Parser) declarationExpression(allowed bool) ast.ExpressionNode {
	switch p.lookahead.Type {
	case token.MACRO:
		return p.macroDefinition(allowed)
	case token.DEF:
		return p.methodDefinition(allowed)
	case token.SIG:
		return p.methodSignatureDefinition(allowed)
	case token.SINGLETON:
		return p.singletonBlockExpression(allowed)
	case token.EXTEND:
		return p.extendWhereBlockExpression(allowed)
	case token.INIT:
		return p.initDefinition(allowed)
	case token.USING:
		return p.usingDeclaration(allowed)
	case token.CLASS:
		return p.classDeclaration(allowed)
	case token.MODULE:
		return p.moduleDeclaration(allowed)
	case token.MIXIN:
		return p.mixinDeclaration(allowed)
	case token.INTERFACE:
		return p.interfaceDeclaration(allowed)
	case token.STRUCT:
		return p.structDeclaration(allowed)
	case token.GETTER:
		return p.getterDeclaration(allowed)
	case token.SETTER:
		return p.setterDeclaration(allowed)
	case token.ATTR:
		return p.attrDeclaration(allowed)
	case token.CONST:
		return p.constantDeclaration(allowed)
	case token.VAR:
		return p.variableDeclaration(allowed)
	case token.TYPEDEF:
		return p.typeDefinition(allowed)
	case token.INCLUDE:
		return p.includeExpression(allowed)
	case token.IMPLEMENT:
		return p.implementExpression(allowed)
	case token.NOINIT:
		return p.noinitModifier(allowed)
	case token.ABSTRACT:
		return p.abstractModifier(allowed)
	case token.PRIMITIVE:
		return p.primitiveModifier(allowed)
	case token.SEALED:
		return p.sealedModifier(allowed)
	case token.ASYNC:
		return p.asyncModifier(allowed)
	case token.DOC_COMMENT:
		return p.docComment(allowed)
	case token.ALIAS:
		return p.aliasDeclaration(allowed)
	}

	return p.modifierExpression()
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
// expressionWithoutModifier "for" pattern "in" expressionWithoutModifier
func (p *Parser) modifierExpression() ast.ExpressionNode {
	left := p.expressionWithoutModifier()

	switch p.lookahead.Type {
	case token.UNLESS, token.WHILE, token.UNTIL:
		mod := p.advance()
		p.swallowNewlines()
		right := p.expressionWithoutModifier()
		return ast.NewModifierNode(
			left.Location().Join(right.Location()),
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
				left.Location().Join(elseExpr.Location()),
				left,
				cond,
				elseExpr,
			)
		}
		return ast.NewModifierNode(
			left.Location().Join(cond.Location()),
			ifTok,
			left,
			cond,
		)
	case token.FOR:
		p.advance()
		p.swallowNewlines()
		param := p.pattern()
		if !ast.PatternDeclaresVariables(param) {
			p.errorMessageLocation("patterns in for in loops should define at least one variable", param.Location())
		}
		p.swallowNewlines()
		inTok, ok := p.consume(token.IN)
		if !ok {
			return ast.NewInvalidNode(inTok.Location(), inTok)
		}
		p.swallowNewlines()
		inExpr := p.expressionWithoutModifier()
		return ast.NewModifierForInNode(
			left.Location().Join(inExpr.Location()),
			left,
			param,
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
			p.errorMessageLocation(
				fmt.Sprintf("invalid `%s` declaration target", p.lookahead.Type.Name()),
				left.Location(),
			)
		}
	} else if ast.IsConstant(left) {
		p.errorMessageLocation(
			"constants cannot be assigned, maybe you meant to declare it with `:=`",
			left.Location(),
		)
	} else if !ast.IsValidAssignmentTarget(left) {
		p.errorMessageLocation(
			fmt.Sprintf("invalid `%s` assignment target", p.lookahead.Type.Name()),
			left.Location(),
		)
	}

	operator := p.advance()
	p.swallowNewlines()

	p.indentedSection = true
	right := p.assignmentExpression()
	p.indentedSection = false

	return ast.NewAssignmentExpressionNode(
		left.Location().Join(right.Location()),
		operator,
		left,
		right,
	)
}

// formalParameter = identifier [":" typeAnnotationWithoutUnionAndVoid] ["=" expressionWithoutModifier]
func (p *Parser) formalParameter() ast.ParameterNode {
	var init ast.ExpressionNode
	var typ ast.TypeNode

	var paramName *token.Token
	var kind ast.ParameterKind
	var location *position.Location

	var restParam bool
	if starTok, ok := p.matchOk(token.STAR); ok {
		kind = ast.PositionalRestParameterKind
		location = starTok.Location()
		restParam = true
	} else if starStarTok, ok := p.matchOk(token.STAR_STAR); ok {
		kind = ast.NamedRestParameterKind
		location = starStarTok.Location()
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
			tok.Location(),
			tok,
		)
	}
	location = location.Join(paramName.Location())

	if p.match(token.COLON) {
		typ = p.typeAnnotationWithoutUnionAndVoid()
		location = location.Join(typ.Location())
	}

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		if restParam {
			p.errorMessageLocation("rest parameters cannot have default values", init.Location())
		}
		location = location.Join(init.Location())
	}

	return ast.NewFormalParameterNode(
		location,
		paramName.Value,
		typ,
		init,
		kind,
	)
}

// formalParameter = ["*" | "**"] identifier [":" typeAnnotationWithoutVoid] ["=" expressionWithoutModifier]
func (p *Parser) methodParameter() ast.ParameterNode {
	var init ast.ExpressionNode
	var typ ast.TypeNode

	var paramName *token.Token
	var setIvar bool
	var kind ast.ParameterKind
	var location *position.Location

	var restParam bool
	if starTok, ok := p.matchOk(token.STAR); ok {
		kind = ast.PositionalRestParameterKind
		location = starTok.Location()
		restParam = true
	} else if starStarTok, ok := p.matchOk(token.STAR_STAR); ok {
		kind = ast.NamedRestParameterKind
		location = starStarTok.Location()
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
			tok.Location(),
			tok,
		)
	}
	location = location.Join(paramName.Location())

	if p.match(token.COLON) {
		typ = p.typeAnnotationWithoutVoid()
		location = location.Join(typ.Location())
	}

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		location = location.Join(init.Location())
		if restParam {
			p.errorMessageLocation("rest parameters cannot have default values", location)
		}
	}

	return ast.NewMethodParameterNode(
		location,
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
	postRestSeen := false
	namedRestSeen := ast.IsNamedRestParam(element)
	elements = append(elements, element)

	for {
		if p.lookahead.Type == token.END_OF_FILE {
			break
		}
		if slices.Contains(stopTokens, p.lookahead.Type) {
			break
		}
		if !p.match(token.COMMA) {
			break
		}
		p.swallowNewlines()
		if slices.Contains(stopTokens, p.lookahead.Type) {
			break
		}
		element := parameter()
		elements = append(elements, element)

		posRest := ast.IsPositionalRestParam(element)
		if posRest {
			if posRestSeen {
				p.errorMessageLocation("there should be only a single positional rest parameter", element.Location())
				continue
			}
			if namedRestSeen {
				p.errorMessageLocation("named rest parameters should appear last", element.Location())
			}
			posRestSeen = true
			continue
		}

		namedRest := ast.IsNamedRestParam(element)
		if namedRest {
			if postRestSeen {
				p.errorMessageLocation("named rest parameters cannot appear after a post parameter", element.Location())
				continue
			}
			if namedRestSeen {
				p.errorMessageLocation("there should be only a single named rest parameter", element.Location())
				continue
			}
			namedRestSeen = true
			continue
		}

		if posRestSeen {
			postRestSeen = true
		}
		if namedRestSeen {
			p.errorMessageLocation("named rest parameters should appear last", element.Location())
			continue
		}

		opt := element.IsOptional()
		if !opt && optionalSeen {
			p.errorMessageLocation("required parameters cannot appear after optional parameters", element.Location())
		} else if opt {
			if !optionalSeen {
				optionalSeen = true
			}
			if posRestSeen {
				p.errorMessageLocation("optional parameters cannot appear after rest parameters", element.Location())
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

// signatureParameter = ["*" | "**"] identifier ["?"] [":" typeAnnotation]
func (p *Parser) signatureParameter() ast.ParameterNode {
	var typ ast.TypeNode
	var opt bool

	var paramName *token.Token
	var kind ast.ParameterKind
	var location *position.Location

	var restParam bool
	if starTok, ok := p.matchOk(token.STAR); ok {
		kind = ast.PositionalRestParameterKind
		location = starTok.Location()
		restParam = true
	} else if starStarTok, ok := p.matchOk(token.STAR_STAR); ok {
		kind = ast.NamedRestParameterKind
		location = starStarTok.Location()
		restParam = true
	}

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
			tok.Location(),
			tok,
		)
	}
	location = location.Join(paramName.Location())

	if questionTok, ok := p.matchOk(token.QUESTION); ok {
		opt = true
		location = location.Join(questionTok.Location())
		if restParam {
			p.errorMessageLocation("rest parameters cannot have default values", location)
		}
	}

	if p.match(token.COLON) {
		typ = p.typeAnnotationWithoutVoid()
		location = location.Join(typ.Location())
	}

	return ast.NewSignatureParameterNode(
		location,
		paramName.Value,
		typ,
		opt,
		kind,
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

// logicalAndExpression = pipeExpression |
// logicalAndExpression "&&" pipeExpression |
// logicalAndExpression "&!" pipeExpression
func (p *Parser) logicalAndExpression() ast.ExpressionNode {
	return p.logicalExpression(p.pipeExpression, token.AND_AND, token.AND_BANG)
}

// pipeExpression = bitwiseOrExpression |
// pipeExpression "|>" bitwiseOrExpression
func (p *Parser) pipeExpression() ast.ExpressionNode {
	left := p.bitwiseOrExpression()

	for {
		op, ok := p.matchOk(token.PIPE_OP)
		if !ok {
			break
		}
		p.swallowNewlines()

		p.indentedSection = true
		right := p.bitwiseOrExpression()
		p.indentedSection = false
		if !ast.IsValidPipeExpressionTarget(right) {
			p.errorMessageLocation(
				"invalid right hand side of a pipe expression, only method and function calls are allowed",
				right.Location(),
			)
		}

		left = ast.NewBinaryExpressionNode(
			left.Location().Join(right.Location()),
			op,
			left,
			right,
		)
	}

	return left
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

// bitwiseAndExpression = bitwiseAndNotExpression | bitwiseAndExpression "&" bitwiseAndNotExpression
func (p *Parser) bitwiseAndExpression() ast.ExpressionNode {
	return p.binaryExpression(p.bitwiseAndNotExpression, token.AND)
}

// bitwiseAndNotExpression = equalityExpression | bitwiseAndNotExpression "&~" equalityExpression
func (p *Parser) bitwiseAndNotExpression() ast.ExpressionNode {
	return p.binaryExpression(p.equalityExpression, token.AND_TILDE)
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
			left.Location().Join(right.Location()),
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
			left.Location().Join(right.Location()),
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

// multiplicativeExpression = rangeLiteral | multiplicativeExpression ("*" | "/" | "%") rangeLiteral
func (p *Parser) multiplicativeExpression() ast.ExpressionNode {
	return p.binaryExpression(p.rangeLiteral, token.STAR, token.SLASH, token.PERCENT)
}

// rangeLiteral = asExpression | asExpression ("..." | "<.." | "..<" | "<.<") [asExpression]
func (p *Parser) rangeLiteral() ast.ExpressionNode {
	left := p.asExpression()
	op, ok := p.matchOk(token.CLOSED_RANGE_OP, token.LEFT_OPEN_RANGE_OP, token.RIGHT_OPEN_RANGE_OP, token.OPEN_RANGE_OP)
	if !ok {
		return left
	}

	if !p.lookahead.IsValidAsEndInRangeLiteral() {
		return ast.NewRangeLiteralNode(
			left.Location().Join(op.Location()),
			op,
			left,
			nil,
		)
	}

	right := p.asExpression()

	return ast.NewRangeLiteralNode(
		left.Location().Join(right.Location()),
		op,
		left,
		right,
	)
}

// asExpression = unaryExpression ["as" strictConstantLookup]
func (p *Parser) asExpression() ast.ExpressionNode {
	expr := p.unaryExpression()
	if p.accept(token.AS) && p.acceptSecond(token.SCOPE_RES_OP, token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT) {
		p.advance()
		runtimeType := p.strictConstantLookup()
		return ast.NewAsExpressionNode(
			expr.Location().Join(runtimeType.Location()),
			expr,
			runtimeType,
		)
	}

	return expr
}

// unaryExpression = powerExpression | ("!" | "-" | "+" | "~" | "&") unaryExpression
func (p *Parser) unaryExpression() ast.ExpressionNode {
	if operator, ok := p.matchOk(token.BANG, token.MINUS, token.PLUS, token.TILDE, token.AND); ok {
		p.swallowNewlines()

		p.indentedSection = true
		right := p.unaryExpression()
		p.indentedSection = false

		return ast.NewUnaryExpressionNode(
			operator.Location().Join(right.Location()),
			operator,
			right,
		)
	}

	return p.powerExpression()
}

// powerExpression = postfixExpression | postfixExpression "**" powerExpression
func (p *Parser) powerExpression() ast.ExpressionNode {
	left := p.postfixExpression()

	if p.lookahead.Type != token.STAR_STAR {
		return left
	}

	operator := p.advance()
	p.swallowNewlines()

	p.indentedSection = true
	right := p.powerExpression()
	p.indentedSection = false

	return ast.NewBinaryExpressionNode(
		left.Location().Join(right.Location()),
		operator,
		left,
		right,
	)
}

// postfixExpression = methodCall ["++" | "--"]
func (p *Parser) postfixExpression() ast.ExpressionNode {
	expr := p.methodCall()

	var op *token.Token
	switch p.lookahead.Type {
	case token.PLUS_PLUS, token.MINUS_MINUS:
		op = p.advance()
	default:
		return expr
	}

	if !ast.IsValidAssignmentTarget(expr) {
		p.errorMessageLocation(
			fmt.Sprintf("invalid `%s` assignment target", op.Type.Name()),
			expr.Location(),
		)
	}

	return ast.NewPostfixExpressionNode(
		expr.Location().Join(op.Location()),
		op,
		expr,
	)
}

// The boolean value indicates whether a comma was the last consumed token.
//
// positionalArgumentList = positionalArgument ("," positionalArgument)*
func (p *Parser) positionalArgumentList(stopTokens ...token.Type) ([]ast.ExpressionNode, bool) {
	var elements []ast.ExpressionNode
	if p.accept(token.STAR_STAR) || p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) && p.acceptSecond(token.COLON) {
		return elements, false
	}
	elements = append(elements, p.positionalArgument())

	for {
		if p.accept(token.END_OF_FILE) {
			break
		}

		if slices.Contains(stopTokens, p.lookahead.Type) {
			break
		}
		if !p.match(token.COMMA) {
			break
		}
		p.swallowNewlines()
		if slices.Contains(stopTokens, p.lookahead.Type) {
			break
		}

		if p.accept(token.STAR_STAR) || p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) && p.acceptSecond(token.COLON) {
			return elements, true
		}
		elements = append(elements, p.positionalArgument())
	}

	return elements, false
}

// positionalArgument = ["*"] expressionWithoutModifier
func (p *Parser) positionalArgument() ast.ExpressionNode {
	if starTok, ok := p.matchOk(token.STAR); ok {
		expr := p.expressionWithoutModifier()
		return ast.NewSplatExpressionNode(
			starTok.Location().Join(expr.Location()),
			expr,
		)
	}

	return p.expressionWithoutModifier()
}

// namedArgument = identifier ":" expressionWithoutModifier | "**" expressionWithoutModifier
func (p *Parser) namedArgument() ast.NamedArgumentNode {
	if starTok, ok := p.matchOk(token.STAR_STAR); ok {
		expr := p.expressionWithoutModifier()
		return ast.NewDoubleSplatExpressionNode(
			starTok.Location().Join(expr.Location()),
			expr,
		)
	}

	ident, ok := p.matchOk(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER)
	if !ok {
		p.errorExpected("an identifier")
		errTok := p.advance()
		return ast.NewInvalidNode(
			errTok.Location(),
			errTok,
		)
	}
	colon, ok := p.consume(token.COLON)
	if !ok {
		return ast.NewInvalidNode(
			colon.Location(),
			colon,
		)
	}
	val := p.expressionWithoutModifier()

	return ast.NewNamedCallArgumentNode(
		ident.Location().Join(val.Location()),
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

func (p *Parser) methodCall() ast.ExpressionNode {
	// function call
	var receiver ast.ExpressionNode

	// receiverless macro
	if p.accept(token.PRIVATE_IDENTIFIER, token.PUBLIC_IDENTIFIER) && p.acceptSecond(token.BANG) {
		macroName := p.advance()
		location := macroName.Location()

		p.advance() // bang !

		var posArgs []ast.ExpressionNode
		var namedArgs []ast.NamedArgumentNode
		if p.accept(token.LPAREN) || p.lookahead.IsValidAsArgumentToNoParenFunctionCall() {
			var lastArgSpan *position.Location
			var errToken *token.Token

			lastArgSpan, posArgs, namedArgs, errToken = p.callArgumentList()
			if errToken != nil {
				return ast.NewInvalidNode(
					errToken.Location(),
					errToken,
				)
			}
			if lastArgSpan == nil {
				p.errorExpected("macro arguments")
				errToken = p.advance()
				return ast.NewInvalidNode(
					errToken.Location(),
					errToken,
				)
			}
			location = location.Join(lastArgSpan)
		}

		if p.hasTrailingClosure() {
			function := p.closureExpression()
			if len(namedArgs) > 0 {
				namedArgs = append(
					namedArgs,
					ast.NewNamedCallArgumentNode(function.Location(), "func", function),
				)
			} else {
				posArgs = append(posArgs, function)
			}
			location = location.Join(function.Location())
		}

		receiver = ast.NewReceiverlessMacroCallNode(
			location,
			macroName.Value,
			posArgs,
			namedArgs,
		)

		// receiverless method
	} else if p.accept(token.PRIVATE_IDENTIFIER, token.PUBLIC_IDENTIFIER) &&
		(p.acceptSecond(token.LPAREN, token.COLON_COLON_LBRACKET) || p.secondLookahead.IsValidAsArgumentToNoParenFunctionCall()) {
		methodName := p.advance()
		location := methodName.Location()

		var typeArgs []ast.TypeNode
		if p.match(token.COLON_COLON_LBRACKET) {
			p.swallowNewlines()
			// generic constructor call
			typeArgs = p.typeAnnotationList(token.RBRACKET)
			p.swallowNewlines()
			rbracket, ok := p.consume(token.RBRACKET)
			if !ok {
				return ast.NewInvalidNode(rbracket.Location(), rbracket)
			}
			location = location.Join(rbracket.Location())
		}
		lastArgSpan, posArgs, namedArgs, errToken := p.callArgumentList()
		if errToken != nil {
			return ast.NewInvalidNode(
				errToken.Location(),
				errToken,
			)
		}
		if lastArgSpan == nil {
			p.errorExpected("method arguments")
			errToken = p.advance()
			return ast.NewInvalidNode(
				errToken.Location(),
				errToken,
			)
		}
		location = location.Join(lastArgSpan)

		if p.hasTrailingClosure() {
			function := p.closureExpression()
			if len(namedArgs) > 0 {
				namedArgs = append(
					namedArgs,
					ast.NewNamedCallArgumentNode(function.Location(), "func", function),
				)
			} else {
				posArgs = append(posArgs, function)
			}
			location = location.Join(function.Location())
		}

		if len(typeArgs) > 0 {
			receiver = ast.NewGenericReceiverlessMethodCallNode(
				location,
				methodName.Value,
				typeArgs,
				posArgs,
				namedArgs,
			)
		} else {
			receiver = ast.NewReceiverlessMethodCallNode(
				location,
				methodName.Value,
				posArgs,
				namedArgs,
			)
		}
	}

	// method call
	if receiver == nil {
		receiver = p.constructorCall()
	}
methodCallLoop:
	for {
		var opToken *token.Token

		// subscript
		if p.accept(token.LBRACKET, token.QUESTION_LBRACKET) {
			nilSafe := p.accept(token.QUESTION_LBRACKET)
			p.advance()
			p.swallowNewlines()

			p.indentedSection = true
			key := p.expressionWithoutModifier()
			p.indentedSection = false

			p.swallowNewlines()
			tok, _ := p.consume(token.RBRACKET)

			if nilSafe {
				receiver = ast.NewNilSafeSubscriptExpressionNode(
					receiver.Location().Join(tok.Location()),
					receiver,
					key,
				)
			} else {
				receiver = ast.NewSubscriptExpressionNode(
					receiver.Location().Join(tok.Location()),
					receiver,
					key,
				)
			}
		}

		if p.accept(token.NEWLINE) && p.acceptSecond(token.DOT, token.QUESTION_DOT, token.DOT_DOT, token.QUESTION_DOT_DOT) {
			p.advance()
			opToken = p.advance()
		} else {
			switch p.lookahead.Type {
			case token.DOT, token.QUESTION_DOT, token.DOT_DOT, token.QUESTION_DOT_DOT:
				opToken = p.advance()
			case token.LBRACKET, token.QUESTION_LBRACKET:
				continue methodCallLoop
			default:
				return receiver
			}
		}

		// closure call
		if p.accept(token.LPAREN) {
			lastLocation, posArgs, namedArgs, errToken := p.callArgumentList()
			if errToken != nil {
				return ast.NewInvalidNode(
					errToken.Location(),
					errToken,
				)
			}
			if p.hasTrailingClosure() {
				function := p.closureExpression()
				if len(namedArgs) > 0 {
					namedArgs = append(
						namedArgs,
						ast.NewNamedCallArgumentNode(function.Location(), "func", function),
					)
				} else {
					posArgs = append(posArgs, function)
				}
				lastLocation = function.Location()
			}

			receiver = ast.NewCallNode(
				receiver.Location().Join(lastLocation),
				receiver,
				opToken.Type == token.QUESTION_DOT,
				posArgs,
				namedArgs,
			)
			continue
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
				errTok.Location(),
				errTok,
			)
		}

		methodNameTok := p.advance()
		methodName := methodNameTok.FetchValue()
		location := receiver.Location().Join(methodNameTok.Location())

		switch methodNameTok.Type {
		case token.AWAIT:
			if opToken.Type != token.DOT {
				p.errorMessageLocation("invalid await operator", opToken.Location())
			}
			receiver = ast.NewAwaitExpressionNode(
				location,
				receiver,
			)
			continue
		case token.MUST:
			if methodNameTok.Type == token.MUST {
				if opToken.Type != token.DOT {
					p.errorMessageLocation("invalid must operator", opToken.Location())
				}
				receiver = ast.NewMustExpressionNode(
					location,
					receiver,
				)
				continue
			}
		}

		var isMacro bool
		if tok, ok := p.matchOk(token.BANG); ok {
			isMacro = true
			location = location.Join(tok.Location())
			if opToken.Type != token.DOT {
				p.errorMessageLocation("invalid macro call operator", opToken.Location())
			}
			if methodNameTok.Type != token.PUBLIC_IDENTIFIER && !methodNameTok.IsKeyword() {
				p.errorMessageLocation("invalid macro name", methodNameTok.Location())
			}
		}

		var typeArgs []ast.TypeNode
		if op, ok := p.matchOk(token.COLON_COLON_LBRACKET); ok {
			p.swallowNewlines()
			// generic constructor call
			typeArgs = p.typeAnnotationList(token.RBRACKET)
			p.swallowNewlines()
			rbracket, ok := p.consume(token.RBRACKET)
			if !ok {
				return ast.NewInvalidNode(rbracket.Location(), rbracket)
			}
			location = location.Join(rbracket.Location())
			if isMacro {
				p.errorMessageLocation("macro calls cannot be generic", op.Location().Join(rbracket.Location()))
			}
		}

		hasParentheses := p.lookahead.Type == token.LPAREN
		lastArgSpan, posArgs, namedArgs, errToken := p.callArgumentList()
		if errToken != nil {
			return ast.NewInvalidNode(
				errToken.Location(),
				errToken,
			)
		}
		if lastArgSpan != nil {
			location = location.Join(lastArgSpan)
		}

		hasTrailingClosure := p.hasTrailingClosure()

		if !isMacro &&
			!hasTrailingClosure &&
			!hasParentheses &&
			len(posArgs) == 0 &&
			len(namedArgs) == 0 &&
			opToken.Type == token.DOT &&
			len(typeArgs) == 0 {
			receiver = ast.NewAttributeAccessNode(
				location,
				receiver,
				methodName,
			)
			continue
		}

		if hasTrailingClosure && (hasParentheses || len(posArgs) == 0 && len(namedArgs) == 0) {
			function := p.closureExpression()
			if len(namedArgs) > 0 {
				namedArgs = append(
					namedArgs,
					ast.NewNamedCallArgumentNode(function.Location(), "func", function),
				)
			} else {
				posArgs = append(posArgs, function)
			}
			location = location.Join(function.Location())
		}

		if isMacro {
			receiver = ast.NewMacroCallNode(
				location,
				receiver,
				methodName,
				posArgs,
				namedArgs,
			)
		} else if len(typeArgs) > 0 {
			receiver = ast.NewGenericMethodCallNode(
				location,
				receiver,
				opToken,
				methodName,
				typeArgs,
				posArgs,
				namedArgs,
			)
		} else {
			receiver = ast.NewMethodCallNode(
				location,
				receiver,
				opToken,
				methodName,
				posArgs,
				namedArgs,
			)
		}
	}
}

func (p *Parser) hasTrailingClosure() bool {
	return p.accept(token.THIN_ARROW) ||
		(p.accept(token.OR_OR) && p.acceptSecond(token.THIN_ARROW)) ||
		(p.accept(token.OR) && (p.acceptSecond(token.STAR, token.STAR_STAR) ||
			p.acceptSecond(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) && p.acceptThird(token.OR, token.COMMA, token.COLON)))
}

// beginlessRangeLiteral = ("..." | "<.<" | "<.." | "..<") constructorCall
func (p *Parser) beginlessRangeLiteral() ast.ExpressionNode {
	op := p.advance()
	right := p.constructorCall()
	return ast.NewRangeLiteralNode(
		op.Location().Join(right.Location()),
		op,
		nil,
		right,
	)
}

// callArgumentListInternal = (positionalArgumentList | namedArgumentList | positionalArgumentList "," namedArgumentList)
// callArgumentList = "(" callArgumentListInternal ")" | callArgumentListInternal
func (p *Parser) callArgumentList() (*position.Location, []ast.ExpressionNode, []ast.NamedArgumentNode, *token.Token) {
	if p.match(token.LPAREN) {
		p.swallowNewlines()
		if rparen, ok := p.matchOk(token.RPAREN); ok {
			return rparen.Location(),
				nil,
				nil,
				nil
		}
		posArgs, commaConsumed := p.positionalArgumentList(token.RPAREN)
		var namedArgs []ast.NamedArgumentNode
		if len(posArgs) == 0 || len(posArgs) > 0 && commaConsumed {
			namedArgs = p.namedArgumentList(token.RPAREN)
		}
		p.swallowNewlines()
		rparen, ok := p.consume(token.RPAREN)
		if !ok {
			return nil,
				nil,
				nil,
				rparen
		}

		return rparen.Location(),
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
	location := position.LocationOfLastElement(posArgs)
	var namedArgs []ast.NamedArgumentNode
	if len(posArgs) == 0 || len(posArgs) > 0 && commaConsumed {
		namedArgs = p.namedArgumentList()
		location = position.LocationOfLastElement(namedArgs)
	}

	return location,
		posArgs,
		namedArgs,
		nil
}

// constructorCall = constantLookup |
// strictConstantLookup ["::[" typeAnnotationList "]"] ( "(" argumentList ")" | argumentList )
func (p *Parser) constructorCall() ast.ExpressionNode {
	if !p.accept(token.PRIVATE_CONSTANT, token.PUBLIC_CONSTANT, token.SCOPE_RES_OP) {
		return p.constantLookup()
	}

	constant := p.strictConstantLookup()
	location := constant.Location()

	var typeArgs []ast.TypeNode
	if p.match(token.COLON_COLON_LBRACKET) {
		p.swallowNewlines()
		// generic constructor call
		typeArgs = p.typeAnnotationList(token.RBRACKET)
		p.swallowNewlines()
		rbracket, ok := p.consume(token.RBRACKET)
		if !ok {
			return ast.NewInvalidNode(rbracket.Location(), rbracket)
		}
		location = location.Join(rbracket.Location())
	}

	lastArgSpan, posArgs, namedArgs, errToken := p.callArgumentList()
	if errToken != nil {
		return ast.NewInvalidNode(
			errToken.Location(),
			errToken,
		)
	}
	if lastArgSpan == nil {
		if typeArgs != nil {
			p.errorMessageLocation("invalid generic constant", location)
		}
		return constant
	}
	location = location.Join(lastArgSpan)

	if len(typeArgs) > 0 {
		return ast.NewGenericConstructorCallNode(
			location,
			constant,
			typeArgs,
			posArgs,
			namedArgs,
		)
	}
	return ast.NewConstructorCallNode(
		location,
		constant,
		posArgs,
		namedArgs,
	)
}

const privateConstantAccessMessage = "cannot access a private constant from the outside"

// constantLookup = primaryExpression | "::" publicConstant | constantLookup "::" publicConstant
func (p *Parser) constantLookup() ast.ExpressionNode {
	var left ast.ExpressionNode
	if tok, ok := p.matchOk(token.SCOPE_RES_OP); ok {
		if p.accept(token.PRIVATE_CONSTANT) {
			p.errorUnexpected(privateConstantAccessMessage)
		}
		right := p.constant()
		left = ast.NewConstantLookupNode(
			tok.Location().Join(right.Location()),
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
			left.Location().Join(right.Location()),
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
			tok.Location().Join(right.Location()),
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
			left.Location().Join(right.Location()),
			left,
			right,
		)
	}

	return left
}

func (p *Parser) publicConstant() *ast.PublicConstantNode {
	tok, ok := p.matchOk(token.PUBLIC_CONSTANT)
	if !ok {
		panic(fmt.Sprintf("invalid public constant token: %#v", tok))
	}
	return ast.NewPublicConstantNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) privateConstant() *ast.PrivateConstantNode {
	tok, ok := p.matchOk(token.PRIVATE_CONSTANT)
	if !ok {
		panic(fmt.Sprintf("invalid private constant token: %#v", tok))
	}
	return ast.NewPrivateConstantNode(
		tok.Location(),
		tok.Value,
	)
}

// constant = privateConstant | publicConstant | unquoteConstant
func (p *Parser) constant() ast.ConstantNode {
	if p.accept(token.PRIVATE_CONSTANT) {
		return p.privateConstant()
	}

	if p.accept(token.PUBLIC_CONSTANT) {
		return p.publicConstant()
	}

	if p.accept(token.UNQUOTE) {
		return p.unquoteConstant()
	}

	p.errorExpected("a constant")
	tok := p.advance()
	p.mode = panicMode
	return ast.NewInvalidNode(
		tok.Location(),
		tok,
	)
}

// newExpression = "new" ( "(" argumentList ")" | argumentList )
func (p *Parser) newExpression() ast.ExpressionNode {
	newTok := p.advance()
	location := newTok.Location()

	lastArgSpan, posArgs, namedArgs, errToken := p.callArgumentList()
	if errToken != nil {
		return ast.NewInvalidNode(
			errToken.Location(),
			errToken,
		)
	}
	location = location.Join(lastArgSpan)

	return ast.NewNewExpressionNode(
		location,
		posArgs,
		namedArgs,
	)
}

func (p *Parser) primaryExpression() ast.ExpressionNode {
	switch p.lookahead.Type {
	case token.NEW:
		return p.newExpression()
	case token.TRUE:
		tok := p.advance()
		return ast.NewTrueLiteralNode(tok.Location())
	case token.FALSE:
		tok := p.advance()
		return ast.NewFalseLiteralNode(tok.Location())
	case token.NIL:
		tok := p.advance()
		return ast.NewNilLiteralNode(tok.Location())
	case token.THIN_ARROW:
		return p.closureAfterArrow(nil, nil, nil, nil)
	case token.SELF:
		return p.selfLiteral()
	case token.BREAK:
		return p.breakExpression()
	case token.AWAIT:
		return p.awaitExpression()
	case token.GO:
		return p.goExpression()
	case token.RETURN:
		return p.returnExpression()
	case token.YIELD:
		return p.yieldExpression()
	case token.CONTINUE:
		return p.continueExpression()
	case token.THROW:
		return p.throwExpression()
	case token.MUST:
		return p.mustExpression()
	case token.TRY:
		return p.tryExpression()
	case token.TYPEOF:
		return p.typeofExpression()
	case token.REGEX_BEG:
		return p.regexLiteral()
	case token.SPECIAL_IDENTIFIER:
		if p.acceptSecond(token.COLON) {
			label := p.advance()
			p.advance() // colon
			expr := p.expressionWithModifier()

			return ast.NewLabeledExpressionNode(
				label.Location().Join(expr.Location()),
				label.Value,
				expr,
			)
		}
	case token.LPAREN:
		p.advance()
		if p.mode == withoutBitwiseOrMode {
			p.mode = normalMode
		}
		expr := p.expressionWithModifier()
		p.consume(token.RPAREN)
		return expr
	case token.LBRACKET:
		return p.arrayListLiteral()
	case token.TUPLE_LITERAL_BEG:
		return p.arrayTupleLiteral()
	case token.HASH_SET_LITERAL_BEG:
		return p.hashSetLiteral()
	case token.WORD_ARRAY_LIST_BEG:
		return p.wordArrayListLiteral()
	case token.WORD_ARRAY_TUPLE_BEG:
		return p.wordArrayTupleLiteral()
	case token.WORD_HASH_SET_BEG:
		return p.wordHashSetLiteral()
	case token.SYMBOL_ARRAY_LIST_BEG:
		return p.symbolArrayListLiteral()
	case token.SYMBOL_ARRAY_TUPLE_BEG:
		return p.symbolArrayTupleLiteral()
	case token.SYMBOL_HASH_SET_BEG:
		return p.symbolHashSetLiteral()
	case token.HEX_ARRAY_LIST_BEG:
		return p.hexArrayListLiteral()
	case token.HEX_ARRAY_TUPLE_BEG:
		return p.hexArrayTupleLiteral()
	case token.HEX_HASH_SET_BEG:
		return p.hexHashSetLiteral()
	case token.BIN_ARRAY_LIST_BEG:
		return p.binArrayListLiteral()
	case token.BIN_ARRAY_TUPLE_BEG:
		return p.binArrayTupleLiteral()
	case token.BIN_HASH_SET_BEG:
		return p.binHashSetLiteral()
	case token.LBRACE:
		return p.hashMapLiteral()
	case token.RECORD_LITERAL_BEG:
		return p.hashRecordLiteral()
	case token.CHAR_LITERAL:
		return p.charLiteral()
	case token.RAW_CHAR_LITERAL:
		return p.rawCharLiteral()
	case token.RAW_STRING:
		return p.rawStringLiteral()
	case token.STRING_BEG:
		return p.stringLiteral(true)
	case token.COLON:
		return p.symbolLiteral(true)
	case token.OR, token.OR_OR:
		return p.closureExpression()
	case token.DOC_COMMENT:
		return p.docComment(false)
	case token.VAR:
		return p.variableDeclaration(false)
	case token.VAL:
		return p.valueDeclaration()
	case token.CONST:
		return p.constantDeclaration(false)
	case token.DEF:
		return p.methodDefinition(false)
	case token.MACRO:
		return p.macroDefinition(false)
	case token.INIT:
		return p.initDefinition(false)
	case token.SWITCH:
		return p.switchExpression()
	case token.IF:
		return p.ifExpression()
	case token.QUOTE:
		return p.quoteExpression()
	case token.UNQUOTE:
		return p.unquoteExpression()
	case token.DO:
		return p.doExpressionOrMacroBoundary()
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
	case token.FORNUM:
		return p.fornumExpression()
	case token.NOINIT:
		return p.noinitModifier(false)
	case token.ABSTRACT:
		return p.abstractModifier(false)
	case token.PRIMITIVE:
		return p.primitiveModifier(false)
	case token.SEALED:
		return p.sealedModifier(false)
	case token.ASYNC:
		return p.asyncModifier(false)
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		return p.identifierOrFunction()
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		return p.constant()
	case token.INSTANCE_VARIABLE:
		return p.instanceVariable()
	case token.INT:
		tok := p.advance()
		return ast.NewIntLiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.INT64:
		tok := p.advance()
		return ast.NewInt64LiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.UINT64:
		tok := p.advance()
		return ast.NewUInt64LiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.INT32:
		tok := p.advance()
		return ast.NewInt32LiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.UINT32:
		tok := p.advance()
		return ast.NewUInt32LiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.INT16:
		tok := p.advance()
		return ast.NewInt16LiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.UINT16:
		tok := p.advance()
		return ast.NewUInt16LiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.INT8:
		tok := p.advance()
		return ast.NewInt8LiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.UINT8:
		tok := p.advance()
		return ast.NewUInt8LiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.FLOAT:
		tok := p.advance()
		return ast.NewFloatLiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.BIG_FLOAT:
		tok := p.advance()
		return ast.NewBigFloatLiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.FLOAT64:
		tok := p.advance()
		return ast.NewFloat64LiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.FLOAT32:
		tok := p.advance()
		return ast.NewFloat32LiteralNode(
			tok.Location(),
			tok.Value,
		)
	case token.ERROR:
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Location(),
			tok,
		)
	case token.GETTER:
		return p.getterDeclaration(false)
	case token.SETTER:
		return p.setterDeclaration(false)
	case token.ATTR:
		return p.attrDeclaration(false)
	case token.USING:
		return p.usingDeclaration(false)
	case token.CLASS:
		return p.classDeclaration(false)
	case token.MODULE:
		return p.moduleDeclaration(false)
	case token.MIXIN:
		return p.mixinDeclaration(false)
	case token.INTERFACE:
		return p.interfaceDeclaration(false)
	case token.STRUCT:
		return p.structDeclaration(false)
	case token.TYPEDEF:
		return p.typeDefinition(false)
	case token.TYPE:
		return p.typeExpression()
	case token.ALIAS:
		return p.aliasDeclaration(false)
	case token.SIG:
		return p.methodSignatureDefinition(false)
	case token.EXTEND:
		return p.extendWhereBlockExpression(false)
	case token.SINGLETON:
		return p.singletonBlockExpression(false)
	case token.INCLUDE:
		return p.includeExpression(false)
	case token.IMPLEMENT:
		return p.implementExpression(false)
	case token.CLOSED_RANGE_OP, token.RIGHT_OPEN_RANGE_OP, token.LEFT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
		return p.beginlessRangeLiteral()
	}

	p.errorExpected("an expression")
	p.updateErrorMode(true)
	tok := p.advance()
	return ast.NewInvalidNode(
		tok.Location(),
		tok,
	)
}

type specialCollectionLiteralWithoutCapacityConstructor[Return, Element ast.Node] func(*position.Location, []Element) Return
type invalidNodeConstructor[Return ast.Node] func(*position.Location, *token.Token) Return

// specialCollectionLiteralWithoutCapacity = beginTokenType (elementProduction)* endTokenType
func specialCollectionLiteralWithoutCapacity[Return, Element ast.Node](p *Parser, elementProduction func() Element, constructor specialCollectionLiteralWithoutCapacityConstructor[Return, Element], invalidConstructor invalidNodeConstructor[Return], endTokenType token.Type) Return {
	begTok := p.advance()
	content := repeatedProduction(p, elementProduction, endTokenType)
	endTok, ok := p.consume(endTokenType)

	if !ok {
		return invalidConstructor(endTok.Location(), endTok)
	}

	return constructor(
		begTok.Location().Join(endTok.Location()),
		content,
	)
}

type specialCollectionLiteralWithCapacityConstructor[Element ast.ExpressionNode] func(*position.Location, []Element, ast.ExpressionNode) ast.ExpressionNode

// specialCollectionLiteralWithCapacity = beginTokenType (elementProduction)* endTokenType [":" primaryExpression]
func specialCollectionLiteralWithCapacity[Element ast.ExpressionNode](p *Parser, elementProduction func() Element, constructor specialCollectionLiteralWithCapacityConstructor[Element], endTokenType token.Type) ast.ExpressionNode {
	begTok := p.advance()
	content := repeatedProduction(p, elementProduction, endTokenType)
	endTok, ok := p.consume(endTokenType)

	if !ok {
		return ast.NewInvalidNode(endTok.Location(), endTok)
	}

	var capacity ast.ExpressionNode
	location := begTok.Location().Join(endTok.Location())
	if p.match(token.COLON) {
		p.swallowNewlines()
		capacity = p.primaryExpression()
		location = location.Join(capacity.Location())
	}

	return constructor(
		location,
		content,
		capacity,
	)
}

// wordArrayListLiteral = "\w[" (rawString)* "]" [":" expressionWithoutModifier]
func (p *Parser) wordArrayListLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithCapacity(
		p,
		p.wordCollectionElement,
		ast.NewWordArrayListLiteralExpressionNode,
		token.WORD_ARRAY_LIST_END,
	)
}

// wordArrayTupleLiteral = "%w[" (rawString)* "]"
func (p *Parser) wordArrayTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.wordCollectionElement,
		ast.NewWordArrayTupleLiteralExpressionNode,
		ast.NewInvalidExpressionNode,
		token.WORD_ARRAY_TUPLE_END,
	)
}

// wordHashSetLiteral = "^w[" (rawString)* "]" [":" expressionWithoutModifier]
func (p *Parser) wordHashSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithCapacity(
		p,
		p.wordCollectionElement,
		ast.NewWordHashSetLiteralNodeI,
		token.WORD_HASH_SET_END,
	)
}

// symbolArrayListLiteral = "\s[" (rawString)* "]" [":" expressionWithoutModifier]
func (p *Parser) symbolArrayListLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithCapacity(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolArrayListLiteralExpressionNode,
		token.SYMBOL_ARRAY_LIST_END,
	)
}

// symbolArrayTupleLiteral = "%s[" (rawString)* "]"
func (p *Parser) symbolArrayTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolArrayTupleLiteralExpressionNode,
		ast.NewInvalidExpressionNode,
		token.SYMBOL_ARRAY_TUPLE_END,
	)
}

// symbolHashSetLiteral = "^s[" (rawString)* "]" [":" expressionWithoutModifier]
func (p *Parser) symbolHashSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithCapacity(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolHashSetLiteralNodeI,
		token.SYMBOL_HASH_SET_END,
	)
}

// hexArrayListLiteral = "\x[" (HEX_INT)* "]" [":" expressionWithoutModifier]
func (p *Parser) hexArrayListLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithCapacity(
		p,
		p.intCollectionElement,
		ast.NewHexArrayListLiteralExpressionNode,
		token.HEX_ARRAY_LIST_END,
	)
}

// hexArrayTupleLiteral = "%x[" (HEX_INT)* "]"
func (p *Parser) hexArrayTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.intCollectionElement,
		ast.NewHexArrayTupleLiteralExpressionNode,
		ast.NewInvalidExpressionNode,
		token.HEX_ARRAY_TUPLE_END,
	)
}

// hexHashSetLiteral = "^x[" (HEX_INT)* "]" [":" expressionWithoutModifier]
func (p *Parser) hexHashSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithCapacity(
		p,
		p.intCollectionElement,
		ast.NewHexHashSetLiteralNodeI,
		token.HEX_HASH_SET_END,
	)
}

// binArrayListLiteral = "\b[" (BIN_INT)* "]" [":" expressionWithoutModifier]
func (p *Parser) binArrayListLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithCapacity(
		p,
		p.intCollectionElement,
		ast.NewBinArrayListLiteralExpressionNode,
		token.BIN_ARRAY_LIST_END,
	)
}

// binArrayTupleLiteral = "%b[" (BIN_INT)* "]"
func (p *Parser) binArrayTupleLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.intCollectionElement,
		ast.NewBinArrayTupleLiteralExpressionNode,
		ast.NewInvalidExpressionNode,
		token.BIN_ARRAY_TUPLE_END,
	)
}

// binHashSetLiteral = "^b[" (BIN_INT)* "]" [":" expressionWithoutModifier]
func (p *Parser) binHashSetLiteral() ast.ExpressionNode {
	return specialCollectionLiteralWithCapacity(
		p,
		p.intCollectionElement,
		ast.NewBinHashSetLiteralNodeI,
		token.BIN_HASH_SET_END,
	)
}

type collectionWithoutCapacityConstructor[N ast.Node] func(*position.Location, []N) N
type collectionElementsProduction[N ast.Node] func(...token.Type) []N

// collectionLiteralWithoutCapacity = startTok [elementsProduction] endTok
func collectionLiteralWithoutCapacity[N ast.Node](p *Parser, endTokType token.Type, elementsProduction collectionElementsProduction[N], constructor collectionWithoutCapacityConstructor[N]) N {
	startTok := p.advance()
	p.swallowNewlines()

	if endTok, ok := p.matchOk(endTokType); ok {
		return constructor(
			startTok.Location().Join(endTok.Location()),
			nil,
		)
	}

	elements := elementsProduction(endTokType)
	p.swallowNewlines()
	endTok, ok := p.consume(endTokType)
	location := startTok.Location()
	if ok {
		location = location.Join(endTok.Location())
	}

	return constructor(
		location,
		elements,
	)
}

type collectionWithCapacityConstructor func(*position.Location, []ast.ExpressionNode, ast.ExpressionNode) ast.ExpressionNode

// collectionLiteralWithCapacity = startTok [elementsProduction] endTok [":" primaryExpression]
func (p *Parser) collectionLiteralWithCapacity(endTokType token.Type, elementsProduction collectionElementsProduction[ast.ExpressionNode], constructor collectionWithCapacityConstructor) ast.ExpressionNode {
	startTok := p.advance()
	p.swallowNewlines()
	var capacity ast.ExpressionNode

	if endTok, ok := p.matchOk(endTokType); ok {
		location := startTok.Location().Join(endTok.Location())
		if p.match(token.COLON) {
			p.swallowNewlines()
			capacity = p.primaryExpression()
			location = location.Join(capacity.Location())
		}
		return constructor(
			location,
			nil,
			capacity,
		)
	}

	elements := elementsProduction(endTokType)
	p.swallowNewlines()
	endTok, ok := p.consume(endTokType)
	if !ok {
		return ast.NewInvalidNode(
			endTok.Location(),
			endTok,
		)
	}

	location := startTok.Location().Join(endTok.Location())
	if p.match(token.COLON) {
		p.swallowNewlines()
		capacity = p.primaryExpression()
		location = location.Join(capacity.Location())
	}

	return constructor(
		location,
		elements,
		capacity,
	)
}

// collectionElementModifier = subProduction |
// subProduction ("if" | "unless") expressionWithoutModifier |
// subProduction "if" expressionWithoutModifier "else" subProduction |
// subProduction "for" identifierList "in" expressionWithoutModifier
func (p *Parser) collectionElementModifier(subProduction func() ast.ExpressionNode) ast.ExpressionNode {
	left := subProduction()

	switch p.lookahead.Type {
	case token.UNLESS:
		mod := p.advance()
		p.swallowNewlines()
		right := p.expressionWithoutModifier()
		return ast.NewModifierNode(
			left.Location().Join(right.Location()),
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
			elseExpr := subProduction()
			return ast.NewModifierIfElseNode(
				left.Location().Join(elseExpr.Location()),
				left,
				cond,
				elseExpr,
			)
		}
		return ast.NewModifierNode(
			left.Location().Join(cond.Location()),
			ifTok,
			left,
			cond,
		)
	case token.FOR:
		p.advance()
		p.swallowNewlines()
		param := p.pattern()
		p.swallowNewlines()
		inTok, ok := p.consume(token.IN)
		if !ok {
			return ast.NewInvalidNode(inTok.Location(), inTok)
		}
		p.swallowNewlines()
		inExpr := p.expressionWithoutModifier()
		return ast.NewModifierForInNode(
			left.Location().Join(inExpr.Location()),
			left,
			param,
			inExpr,
		)
	}

	return left
}

// "{" [recordLiteralElements] "}" [":" primaryExpression]
func (p *Parser) hashMapLiteral() ast.ExpressionNode {
	return p.collectionLiteralWithCapacity(token.RBRACE, p.recordLiteralElements, ast.NewHashMapLiteralNodeI)
}

// "%{" [recordLiteralElements] "}"
func (p *Parser) hashRecordLiteral() ast.ExpressionNode {
	return collectionLiteralWithoutCapacity(p, token.RBRACE, p.recordLiteralElements, ast.NewHashRecordLiteralNodeI)
}

// arrayListLiteral = "[" [listLikeLiteralElements] "]" [":" primaryExpression]
func (p *Parser) arrayListLiteral() ast.ExpressionNode {
	return p.collectionLiteralWithCapacity(token.RBRACKET, p.listLikeLiteralElements, ast.NewArrayListLiteralNodeI)
}

// arrayTupleLiteral = "%[" [listLikeLiteralElements] "]"
func (p *Parser) arrayTupleLiteral() ast.ExpressionNode {
	return collectionLiteralWithoutCapacity(p, token.RBRACKET, p.listLikeLiteralElements, ast.NewArrayTupleLiteralNodeI)
}

// listLikeLiteralElements = listLikeLiteralElement ("," listLikeLiteralElement)*
func (p *Parser) listLikeLiteralElements(stopTokens ...token.Type) []ast.ExpressionNode {
	return commaSeparatedList(p, p.listLikeLiteralElement, stopTokens...)
}

// listLikeLiteralElement = "*" expressionWithoutModifier |
// keyValueExpression |
// keyValueExpression ("if" | "unless") expressionWithoutModifier |
// keyValueExpression "if" expressionWithoutModifier "else" expressionWithoutModifier |
// keyValueExpression "for" identifierList "in" expressionWithoutModifier
func (p *Parser) listLikeLiteralElement() ast.ExpressionNode {
	if p.accept(token.STAR) {
		starTok := p.advance()
		expr := p.expressionWithoutModifier()
		return ast.NewSplatExpressionNode(
			starTok.Location().Join(expr.Location()),
			expr,
		)
	}
	if p.accept(token.STAR_STAR) {
		p.errorMessage("double splats cannot appear in list, tuple nor set literals")
		starTok := p.advance()
		expr := p.expressionWithoutModifier()
		return ast.NewDoubleSplatExpressionNode(
			starTok.Location().Join(expr.Location()),
			expr,
		)
	}

	return p.collectionElementModifier(p.keyValueExpression)
}

// recordLiteralElements = recordLiteralElement ("," recordLiteralElement)*
func (p *Parser) recordLiteralElements(stopTokens ...token.Type) []ast.ExpressionNode {
	return commaSeparatedList(p, p.recordLiteralElement, stopTokens...)
}

// recordLiteralElement = "**" expressionWithoutModifier |
// keyValueMapExpression |
// keyValueMapExpression ("if" | "unless") expressionWithoutModifier |
// keyValueMapExpression "if" expressionWithoutModifier "else" expressionWithoutModifier |
// keyValueMapExpression "for" identifierList "in" expressionWithoutModifier
func (p *Parser) recordLiteralElement() ast.ExpressionNode {
	if p.accept(token.STAR_STAR) {
		starTok := p.advance()
		expr := p.expressionWithoutModifier()
		return ast.NewDoubleSplatExpressionNode(
			starTok.Location().Join(expr.Location()),
			expr,
		)
	}
	if p.accept(token.STAR) {
		p.errorMessage("splats cannot appear in record nor map literals")
		starTok := p.advance()
		expr := p.expressionWithoutModifier()

		return ast.NewSplatExpressionNode(
			starTok.Location().Join(expr.Location()),
			expr,
		)
	}

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
		p.acceptSecond(token.COLON) {
		key := p.advance()
		p.advance()
		p.swallowNewlines()
		val := p.expressionWithoutModifier()
		return ast.NewSymbolKeyValueExpressionNode(
			key.Location().Join(val.Location()),
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
			p.errorMessageLocation("expected a key-value pair, map literals should consist of key-value pairs", key.Location())
			return key
		}
	}

	p.swallowNewlines()
	val := p.expressionWithoutModifier()

	return ast.NewKeyValueExpressionNode(
		key.Location().Join(val.Location()),
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
			key.Location().Join(value.Location()),
			key,
			value,
		)
	}

	return key
}

// hashSetLiteral = "^[" [hashSetLiteralElements] "]" [":" primaryExpression]
func (p *Parser) hashSetLiteral() ast.ExpressionNode {
	return p.collectionLiteralWithCapacity(token.RBRACKET, p.hashSetLiteralElements, ast.NewHashSetLiteralNodeI)
}

// hashSetLiteralElements = hashSetLiteralElement ("," hashSetLiteralElement)*
func (p *Parser) hashSetLiteralElements(stopTokens ...token.Type) []ast.ExpressionNode {
	return commaSeparatedList(p, p.hashSetLiteralElement, stopTokens...)
}

// hashSetLiteralElement = ["*"] expressionWithoutModifier |
// expressionWithoutModifier ("if" | "unless") expressionWithoutModifier |
// expressionWithoutModifier "if" expressionWithoutModifier "else" expressionWithoutModifier |
// expressionWithoutModifier "for" identifierList "in" expressionWithoutModifier
func (p *Parser) hashSetLiteralElement() ast.ExpressionNode {
	if p.accept(token.STAR) {
		starTok := p.advance()
		expr := p.expressionWithoutModifier()
		return ast.NewSplatExpressionNode(
			starTok.Location().Join(expr.Location()),
			expr,
		)
	}
	if p.accept(token.STAR_STAR) {
		p.errorMessage("double splats cannot appear in list, tuple nor set literals")
		starTok := p.advance()
		expr := p.expressionWithoutModifier()
		return ast.NewDoubleSplatExpressionNode(
			starTok.Location().Join(expr.Location()),
			expr,
		)
	}

	return p.collectionElementModifier(p.expressionWithoutModifier)
}

// selfLiteral = "self"
func (p *Parser) selfLiteral() *ast.SelfLiteralNode {
	tok := p.advance()
	return ast.NewSelfLiteralNode(tok.Location())
}

// genericConstantList = genericConstant ("," genericConstant)*
func (p *Parser) genericConstantList(stopTokens ...token.Type) []ast.ComplexConstantNode {
	return commaSeparatedListWithoutTerminator(p, p.genericConstant, stopTokens...)
}

// constantList = constant ("," constant)*
func (p *Parser) constantList(stopTokens ...token.Type) []ast.ComplexConstantNode {
	return commaSeparatedListWithoutTerminator(p, p.strictConstantLookup, stopTokens...)
}

// usingEntryList = usingEntry ("," usingEntry)*
func (p *Parser) usingEntryList(stopTokens ...token.Type) []ast.UsingEntryNode {
	return commaSeparatedListWithoutTerminator(p, p.usingEntry, stopTokens...)
}

// usingSubentryList = usingSubentry ("," usingSubentry)*
func (p *Parser) usingSubentryList(stopTokens ...token.Type) []ast.UsingSubentryNode {
	return commaSeparatedList(p, p.usingSubentry, stopTokens...)
}

// typeAnnotationList = typeAnnotation ("," typeAnnotation)*
func (p *Parser) typeAnnotationList(stopTokens ...token.Type) []ast.TypeNode {
	return commaSeparatedList(p, p.typeAnnotation, stopTokens...)
}

// includeExpression = "include" genericConstantList ["where" typeParameterListWithoutTerminator]
func (p *Parser) includeExpression(allowed bool) *ast.IncludeExpressionNode {
	keyword := p.advance()
	consts := p.genericConstantList()
	location := position.JoinLocationOfLastElement(keyword.Location(), consts)

	if !allowed {
		p.errorMessageLocation(
			"this definition cannot appear in expressions",
			location,
		)
	}
	return ast.NewIncludeExpressionNode(
		location,
		consts,
	)
}

// implementExpression = "implement" genericConstantList
func (p *Parser) implementExpression(allowed bool) *ast.ImplementExpressionNode {
	keyword := p.advance()
	consts := p.genericConstantList()
	location := position.JoinLocationOfLastElement(keyword.Location(), consts)

	if !allowed {
		p.errorMessageLocation(
			"this definition cannot appear in expressions",
			location,
		)
	}
	return ast.NewImplementExpressionNode(
		location,
		consts,
	)
}

// methodDefinition = "sig" methodName ["(" signatureParameterList ")"] [":" typeAnnotation] ["!" typeAnnotation]
func (p *Parser) methodSignatureDefinition(allowed bool) ast.ExpressionNode {
	var params []ast.ParameterNode
	var returnType ast.TypeNode
	var throwType ast.TypeNode
	var location *position.Location

	sigTok := p.advance()
	location = sigTok.Location()

	methodName, mSpan := p.methodName()
	location = location.Join(mSpan)

	var typeParams []ast.TypeParameterNode
	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeParams = p.typeParameterList(token.RBRACKET)
			rbracket, ok := p.consume(token.RBRACKET)
			if !ok {
				return ast.NewInvalidNode(
					rbracket.Location(),
					rbracket,
				)
			}
			location = location.Join(rbracket.Location())
		}
	}

	if p.match(token.LPAREN) {
		if rparen, ok := p.matchOk(token.RPAREN); ok {
			location = location.Join(rparen.Location())
		} else {
			params = p.signatureParameterList(token.RPAREN)

			rparen, ok := p.consume(token.RPAREN)
			location = location.Join(rparen.Location())
			if !ok {
				return ast.NewInvalidNode(
					rparen.Location(),
					rparen,
				)
			}
		}
	}

	// return type
	if p.match(token.COLON) {
		returnType = p.typeAnnotation()
		location = location.Join(returnType.Location())
	}

	// throw type
	if p.match(token.BANG) {
		throwType = p.typeAnnotation()
		location = location.Join(throwType.Location())
	}

	if !allowed {
		p.errorMessageLocation(
			"signature definitions cannot appear in expressions",
			location,
		)
	}

	return ast.NewMethodSignatureDefinitionNode(
		location,
		"",
		methodName,
		typeParams,
		params,
		returnType,
		throwType,
	)
}

// aliasEntry = identifier identifier
func (p *Parser) aliasEntry() *ast.AliasDeclarationEntry {
	newName, newNameSpan := p.methodName()
	p.swallowNewlines()
	oldName, oldNameSpan := p.methodName()

	return ast.NewAliasDeclarationEntry(
		newNameSpan.Join(oldNameSpan),
		newName,
		oldName,
	)
}

// aliasEntryList = aliasEntry ("," aliasEntry)*
func (p *Parser) aliasEntryList(stopTokens ...token.Type) []*ast.AliasDeclarationEntry {
	return commaSeparatedListWithoutTerminator(p, p.aliasEntry, stopTokens...)
}

// aliasDeclaration = "alias" methodName methodName
func (p *Parser) aliasDeclaration(allowed bool) ast.ExpressionNode {
	aliasTok := p.advance()
	p.swallowNewlines()

	entries := p.aliasEntryList()

	location := position.JoinLocationOfLastElement(aliasTok.Location(), entries)
	if !allowed {
		p.errorMessageLocation(
			"alias definitions cannot appear in expressions",
			location,
		)
	}
	return ast.NewAliasDeclarationNode(
		location,
		entries,
	)
}

// typeExpression = "type" typeAnnotation
func (p *Parser) typeExpression() ast.ExpressionNode {
	typeTok := p.advance()

	p.swallowNewlines()

	typ := p.typeAnnotation()
	return ast.NewTypeExpressionNode(
		typeTok.Location().Join(typ.Location()),
		typ,
	)
}

// typeDeclaration = "typedef" strictConstantLookup "=" typeAnnotation
func (p *Parser) typeDefinition(allowed bool) ast.ExpressionNode {
	typedefTok := p.advance()

	name := p.strictConstantLookup()
	var typeVars []ast.TypeParameterNode

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeParameterList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Location(),
					errTok,
				)
			}
		}
	}
	equalTok, ok := p.consume(token.EQUAL_OP)
	if !ok {
		return ast.NewInvalidNode(equalTok.Location(), equalTok)
	}
	p.swallowNewlines()

	typ := p.typeAnnotation()
	location := typedefTok.Location().Join(typ.Location())
	if !allowed {
		p.errorMessageLocation(
			"type definitions cannot appear in expressions",
			location,
		)
	}
	if len(typeVars) > 0 {
		return ast.NewGenericTypeDefinitionNode(
			location,
			"",
			name,
			typeVars,
			typ,
		)
	}
	return ast.NewTypeDefinitionNode(
		location,
		"",
		name,
		typ,
	)
}

func (p *Parser) macroName() (*token.Token, bool) {
	if p.lookahead.IsValidMacroName() {
		macroNameTok := p.advance()
		return macroNameTok, true
	} else {
		p.errorExpected("a macro name (public identifier, keyword)")
		p.updateErrorMode(true)
		tok := p.advance()
		return tok, false
	}
}

func (p *Parser) methodName() (string, *position.Location) {
	var methodName string
	var location *position.Location

	if p.lookahead.IsValidRegularMethodName() {
		methodNameTok := p.advance()
		methodName = methodNameTok.FetchValue()
		location = methodNameTok.Location()
		if tok, ok := p.matchOk(token.EQUAL_OP); ok {
			methodName += "="
			location = location.Join(tok.Location())
		}
	} else if p.accept(token.LBRACKET) && p.acceptSecond(token.RBRACKET) {
		// [
		tok := p.advance()
		location = tok.Location()

		// ]
		tok = p.advance()
		location = location.Join(tok.Location())
		methodName = "[]"

		if tok, ok := p.matchOk(token.EQUAL_OP); ok {
			methodName += "="
			location = location.Join(tok.Location())
		}
	} else {
		if !p.lookahead.IsOverridableOperator() {
			p.errorExpected("a method name (identifier, overridable operator)")
		}
		tok := p.advance()
		methodName = tok.FetchValue()
		location = tok.Location()
	}

	return methodName, location
}

// macroDefinition = "macro" macroName ["(" formalParameterList ")"] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) macroDefinition(allowed bool) ast.ExpressionNode {
	var params []ast.ParameterNode
	var body []ast.StatementNode
	var location *position.Location

	macroTok := p.advance()
	location = macroTok.Location()

	p.swallowNewlines()
	macroNameTok, ok := p.macroName()
	if !ok {
		return ast.NewInvalidNode(macroNameTok.Location(), macroNameTok)
	}

	if p.match(token.LPAREN) {
		p.swallowNewlines()
		if !p.match(token.RPAREN) {
			params = p.formalParameterList(token.RPAREN)

			p.swallowNewlines()
			rparen, ok := p.consume(token.RPAREN)
			if !ok {
				return ast.NewInvalidNode(
					rparen.Location(),
					rparen,
				)
			}
			location = location.Join(rparen.Location())
		}
	}

	lastLocation, body, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = location.Join(lastLocation)
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
			location = location.Join(endTok.Location())
		}
	}

	if !allowed {
		p.errorMessageLocation(
			"macro definitions cannot appear in expressions",
			location,
		)
	}

	return ast.NewMacroDefinitionNode(
		location,
		"",
		false,
		macroNameTok.FetchValue(),
		params,
		body,
	)
}

// methodDefinition = "def" ["*"] methodName ["(" methodParameterList ")"] [":" typeAnnotation] ["!" typeAnnotation] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) methodDefinition(allowed bool) ast.ExpressionNode {
	var params []ast.ParameterNode
	var returnType ast.TypeNode
	var throwType ast.TypeNode
	var body []ast.StatementNode
	var location *position.Location

	defTok := p.advance()
	location = defTok.Location()
	var isGenerator bool
	if p.accept(token.STAR) && p.secondLookahead.IsValidMethodName() {
		p.advance() // swallow "*"
		isGenerator = true
	}
	p.swallowNewlines()
	methodName, methodNameSpan := p.methodName()
	var isSetter bool
	var isSubscriptSetter bool

	if methodName == "[]=" {
		isSubscriptSetter = true
	} else if len(methodName) > 0 {
		firstChar, _ := utf8.DecodeRuneInString(methodName)
		lastChar := methodName[len(methodName)-1]
		if (unicode.IsLetter(firstChar) || firstChar == '_') && lastChar == '=' {
			isSetter = true
		}
	}

	var typeParams []ast.TypeParameterNode
	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeParams = p.typeParameterList(token.RBRACKET)
			rbracket, ok := p.consume(token.RBRACKET)
			if !ok {
				return ast.NewInvalidNode(
					rbracket.Location(),
					rbracket,
				)
			}
			location = location.Join(rbracket.Location())
		}
	}

	if p.match(token.LPAREN) {
		p.swallowNewlines()
		if !p.match(token.RPAREN) {
			params = p.methodParameterList(token.RPAREN)

			p.swallowNewlines()
			rparen, ok := p.consume(token.RPAREN)
			if !ok {
				return ast.NewInvalidNode(
					rparen.Location(),
					rparen,
				)
			}
			location = location.Join(rparen.Location())
		}
	}

	if isSubscriptSetter {
		var location *position.Location
		if len(params) == 0 {
			location = methodNameSpan
		} else {
			location = position.JoinLocationOfCollection(params)
		}

		if len(params) != 2 {
			p.errorMessageLocation(fmt.Sprintf("subscript setter methods must have two parameters, got: %d", len(params)), location)
		}
	} else if isSetter {
		if len(params) == 0 {
			p.errorMessageLocation("setter methods must have a single parameter, got: 0", methodNameSpan)
		} else if len(params) > 1 {
			p.errorMessageLocation(fmt.Sprintf("setter methods must have a single parameter, got: %d", len(params)), position.JoinLocationOfCollection(params[1:]))
		}
	}

	// return type
	if p.match(token.COLON) {
		returnType = p.typeAnnotation()
		location = location.Join(returnType.Location())
		if isSetter || isSubscriptSetter {
			p.errorMessageLocation("setter methods cannot be defined with custom return types", returnType.Location())
		}
	}

	// throw type
	if p.match(token.BANG) {
		throwType = p.typeAnnotation()
		location = location.Join(throwType.Location())
	}

	lastLocation, body, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = location.Join(lastLocation)
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
			location = location.Join(endTok.Location())
		}
	}

	if !allowed {
		p.errorMessageLocation(
			"method definitions cannot appear in expressions",
			location,
		)
	}

	var flags bitfield.BitFlag8
	if isGenerator {
		flags |= ast.METHOD_GENERATOR_FLAG
	}
	return ast.NewMethodDefinitionNode(
		location,
		"",
		flags,
		methodName,
		typeParams,
		params,
		returnType,
		throwType,
		body,
	)
}

func (p *Parser) newLocation(location *position.Span) *position.Location {
	return position.NewLocation(p.sourceName, location)
}

// initDefinition = "init" ["(" methodParameterList ")"] ["!" typeAnnotation] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) initDefinition(allowed bool) ast.ExpressionNode {
	var params []ast.ParameterNode
	var throwType ast.TypeNode
	var body []ast.StatementNode
	var location *position.Location

	initTok := p.advance()

	// methodParameterList
	if p.match(token.LPAREN) {
		p.swallowNewlines()
		if !p.match(token.RPAREN) {
			params = p.methodParameterList(token.RPAREN)

			p.swallowNewlines()
			if tok, ok := p.consume(token.RPAREN); !ok {
				return ast.NewInvalidNode(
					tok.Location(),
					tok,
				)
			}
		}
	}

	// throw type
	if p.match(token.BANG) {
		throwType = p.typeAnnotation()
	}

	lastLocation, body, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = initTok.Location().Join(lastLocation)
	} else {
		location = initTok.Location()
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
			location = location.Join(endTok.Location())
		}
	}

	if !allowed {
		p.errorMessageLocation(
			"method definitions cannot appear in expressions",
			location,
		)
	}

	return ast.NewInitDefinitionNode(
		location,
		params,
		throwType,
		body,
	)
}

// typeParameter = ["+" | "-"] constant (":=" TypeNode | [">" TypeNode] ["<" TypeNode] ["=" TypeNode])
func (p *Parser) typeParameter() ast.TypeParameterNode {
	variance := ast.INVARIANT
	var firstSpan *position.Location
	var lastLocation *position.Location
	var lowerBound ast.TypeNode
	var upperBound ast.TypeNode
	var def ast.TypeNode

	switch p.lookahead.Type {
	case token.PLUS:
		plusTok := p.advance()
		firstSpan = plusTok.Location()
		variance = ast.COVARIANT
	case token.MINUS:
		minusTok := p.advance()
		firstSpan = minusTok.Location()
		variance = ast.CONTRAVARIANT
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
	default:
		errTok := p.advance()
		p.errorExpected("a type variable")
		return ast.NewInvalidNode(
			errTok.Location(),
			errTok,
		)
	}

	if !p.accept(token.PRIVATE_CONSTANT, token.PUBLIC_CONSTANT) {
		errTok := p.advance()
		return ast.NewInvalidNode(
			errTok.Location(),
			errTok,
		)
	}
	nameTok := p.advance()
	if firstSpan == nil {
		firstSpan = nameTok.Location()
	}
	lastLocation = nameTok.Location()

	if p.match(token.COLON_EQUAL) {
		lowerBound = p.typeAnnotation()
		upperBound = lowerBound
		def = lowerBound
		lastLocation = lowerBound.Location()
	} else {
		if p.match(token.GREATER) {
			lowerBound = p.typeAnnotation()
			lastLocation = lowerBound.Location()
		}

		if p.match(token.LESS) {
			upperBound = p.typeAnnotation()
			lastLocation = upperBound.Location()
		}

		if p.match(token.EQUAL_OP) {
			def = p.typeAnnotation()
			lastLocation = def.Location()
		}
	}

	return ast.NewVariantTypeParameterNode(
		firstSpan.Join(lastLocation),
		variance,
		nameTok.Value,
		lowerBound,
		upperBound,
		def,
	)
}

// typeParameterList = typeParameter ("," typeParameter)*
func (p *Parser) typeParameterList(stopTokens ...token.Type) []ast.TypeParameterNode {
	return commaSeparatedList(p, p.typeParameter, stopTokens...)
}

// typeParameterListWithoutTerminator = typeParameter ("," typeParameter)*
func (p *Parser) typeParameterListWithoutTerminator(stopTokens ...token.Type) []ast.TypeParameterNode {
	return commaSeparatedListWithoutTerminator(p, p.typeParameter, stopTokens...)
}

// attributeParameter = identifier [":" typeAnnotation] ["=" expressionWithoutModifier]
func (p *Parser) attributeParameter() ast.ParameterNode {
	var init ast.ExpressionNode
	var typ ast.TypeNode

	var paramName *token.Token
	var location *position.Location

	switch p.lookahead.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		paramName = p.advance()
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		p.errorExpected("a lowercase identifier as the name of the declared attribute")
		paramName = p.advance()
	default:
		p.errorExpected("an identifier as the name of the declared attribute")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Location(),
			tok,
		)
	}
	location = location.Join(paramName.Location())

	if p.match(token.COLON) {
		typ = p.typeAnnotation()
		location = location.Join(typ.Location())
	}

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		location = location.Join(init.Location())
	}

	return ast.NewAttributeParameterNode(
		location,
		paramName.Value,
		typ,
		init,
	)
}

// setterParameter = identifier [":" typeAnnotation]
func (p *Parser) setterParameter() ast.ParameterNode {
	var init ast.ExpressionNode
	var typ ast.TypeNode

	var paramName *token.Token
	var location *position.Location

	switch p.lookahead.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		paramName = p.advance()
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT:
		p.errorExpected("a lowercase identifier as the name of the declared attribute")
		paramName = p.advance()
	default:
		p.errorExpected("an identifier as the name of the declared attribute")
		tok := p.advance()
		return ast.NewInvalidNode(
			tok.Location(),
			tok,
		)
	}
	location = location.Join(paramName.Location())

	if p.match(token.COLON) {
		typ = p.typeAnnotation()
		location = location.Join(typ.Location())
	}

	if p.match(token.EQUAL_OP) {
		init = p.expressionWithoutModifier()
		p.errorMessageLocation("setter declarations cannot have initialisers", init.Location())
		location = location.Join(init.Location())
	}

	return ast.NewAttributeParameterNode(
		location,
		paramName.Value,
		typ,
		init,
	)
}

// attributeParameterList = attributeParameter ("," attributeParameter)*
func (p *Parser) attributeParameterList(stopTokens ...token.Type) []ast.ParameterNode {
	return commaSeparatedListWithoutTerminator(p, p.attributeParameter, stopTokens...)
}

// setterParameterList = setterParameter ("," setterParameter)*
func (p *Parser) setterParameterList(stopTokens ...token.Type) []ast.ParameterNode {
	return commaSeparatedListWithoutTerminator(p, p.setterParameter, stopTokens...)
}

// getterDeclaration = "getter" attributeParameterList
func (p *Parser) getterDeclaration(allowed bool) ast.ExpressionNode {
	getterTok := p.advance()
	p.swallowNewlines()
	attrList := p.attributeParameterList()

	if !allowed {
		p.errorMessageLocation(
			"getter declarations cannot appear in expressions",
			getterTok.Location(),
		)
	}

	return ast.NewGetterDeclarationNode(
		position.JoinLocationOfLastElement(getterTok.Location(), attrList),
		"",
		attrList,
	)
}

// setterDeclaration = "setter" attributeParameterList
func (p *Parser) setterDeclaration(allowed bool) ast.ExpressionNode {
	setterTok := p.advance()
	p.swallowNewlines()
	attrList := p.setterParameterList()

	if !allowed {
		p.errorMessageLocation(
			"setter declarations cannot appear in expressions",
			setterTok.Location(),
		)
	}

	return ast.NewSetterDeclarationNode(
		position.JoinLocationOfLastElement(setterTok.Location(), attrList),
		"",
		attrList,
	)
}

// attrDeclaration = "attr" attributeParameterList
func (p *Parser) attrDeclaration(allowed bool) ast.ExpressionNode {
	attrTok := p.advance()
	p.swallowNewlines()
	attrList := p.attributeParameterList()

	if !allowed {
		p.errorMessageLocation(
			"attr declarations cannot appear in expressions",
			attrTok.Location(),
		)
	}

	return ast.NewAttrDeclarationNode(
		position.JoinLocationOfLastElement(attrTok.Location(), attrList),
		"",
		attrList,
	)
}

// usingSubentry = publicConstant ["as" publicConstant] | publicIdentifier ["as" publicIdentifier]
func (p *Parser) usingSubentry() ast.UsingSubentryNode {
	switch p.lookahead.Type {
	case token.PUBLIC_IDENTIFIER:
		identTok := p.advance()
		if !p.accept(token.AS) {
			return ast.NewPublicIdentifierNode(
				identTok.Location(),
				identTok.Value,
			)
		}

		p.advance() // as
		asIdentTok, ok := p.consume(token.PUBLIC_IDENTIFIER)
		if !ok {
			return ast.NewInvalidNode(asIdentTok.Location(), asIdentTok)
		}
		return ast.NewPublicIdentifierAsNode(
			identTok.Location().Join(asIdentTok.Location()),
			ast.NewPublicIdentifierNode(identTok.Location(), identTok.Value),
			asIdentTok.Value,
		)
	case token.PUBLIC_CONSTANT:
		constTok := p.advance()
		if !p.accept(token.AS) {
			return ast.NewPublicConstantNode(
				constTok.Location(),
				constTok.Value,
			)
		}

		p.advance() // as
		asConstTok, ok := p.consume(token.PUBLIC_CONSTANT)
		if !ok {
			return ast.NewInvalidNode(asConstTok.Location(), asConstTok)
		}
		return ast.NewPublicConstantAsNode(
			constTok.Location().Join(asConstTok.Location()),
			ast.NewPublicConstantNode(constTok.Location(), constTok.Value),
			asConstTok.Value,
		)
	default:
		p.errorExpected("a public identifier or public constant")
		tok := p.advance()
		return ast.NewInvalidNode(tok.Location(), tok)
	}
}

// usingEntry = strictConstantLookup ["::" (publicIdentifier | "*" | "{" usingSubentryList "}")]
func (p *Parser) usingEntry() ast.UsingEntryNode {
	var left ast.ComplexConstantNode
	if tok, ok := p.matchOk(token.SCOPE_RES_OP); ok {
		if p.accept(token.PRIVATE_CONSTANT) {
			p.errorUnexpected(privateConstantAccessMessage)
		}
		right := p.constant()
		left = ast.NewConstantLookupNode(
			tok.Location().Join(right.Location()),
			nil,
			right,
		)
	} else {
		left = p.constant()
	}

	for p.lookahead.Type == token.SCOPE_RES_OP {
		p.advance()

		p.swallowNewlines()
		if p.accept(token.PUBLIC_IDENTIFIER) {
			nameTok := p.advance()
			methodLookup := ast.NewMethodLookupNode(
				left.Location().Join(nameTok.Location()),
				left,
				nameTok.Value,
			)
			if p.match(token.AS) {
				asNameTok, ok := p.consume(token.PUBLIC_IDENTIFIER)
				if !ok {
					return ast.NewInvalidNode(
						asNameTok.Location(),
						asNameTok,
					)
				}
				return ast.NewMethodLookupAsNode(
					left.Location(),
					methodLookup,
					asNameTok.Value,
				)
			}
			return methodLookup
		}
		if p.accept(token.STAR) {
			starTok := p.advance()
			return ast.NewUsingAllEntryNode(
				left.Location().Join(starTok.Location()),
				left,
			)
		}
		if p.accept(token.LBRACE) {
			p.advance()
			subentryList := p.usingSubentryList(token.RBRACE)
			rbrace, ok := p.consume(token.RBRACE)
			if !ok {
				ast.NewInvalidNode(rbrace.Location(), rbrace)
			}

			return ast.NewUsingEntryWithSubentriesNode(
				left.Location().Join(rbrace.Location()),
				left,
				subentryList,
			)
		}

		if p.accept(token.PRIVATE_CONSTANT) {
			p.errorUnexpected(privateConstantAccessMessage)
		}
		right := p.constant()

		left = ast.NewConstantLookupNode(
			left.Location().Join(right.Location()),
			left,
			right,
		)
	}

	if p.match(token.AS) {
		asNameTok, ok := p.consume(token.PUBLIC_CONSTANT)
		if !ok {
			return ast.NewInvalidNode(
				asNameTok.Location(),
				asNameTok,
			)
		}
		return ast.NewConstantAsNode(
			left.Location(),
			left,
			asNameTok.Value,
		)
	}

	return left
}

// usingDeclaration = "using" usingEntryList
func (p *Parser) usingDeclaration(allowed bool) ast.ExpressionNode {
	usingTok := p.advance()
	p.swallowNewlines()
	constList := p.usingEntryList()

	if !allowed {
		p.errorMessageLocation(
			"using declarations cannot appear in expressions",
			usingTok.Location(),
		)
	}

	return ast.NewUsingExpressionNode(
		position.JoinLocationOfLastElement(usingTok.Location(), constList),
		constList,
	)
}

// classDeclaration = "class" [constantLookup] ["[" typeVariableList "]"] ["<" genericConstantOrNil] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) classDeclaration(allowed bool) ast.ExpressionNode {
	classTok := p.advance()
	var superclass ast.ExpressionNode
	var constant ast.ExpressionNode
	var typeVars []ast.TypeParameterNode
	if !p.accept(token.LESS, token.LBRACKET, token.SEMICOLON, token.NEWLINE) {
		constant = p.constantLookup()

		switch c := constant.(type) {
		case *ast.PublicConstantNode,
			*ast.PrivateConstantNode,
			*ast.ConstantLookupNode:
		case *ast.UnquoteNode:
			c.Kind = ast.UNQUOTE_CONSTANT_KIND
		default:
			p.errorMessageLocation("invalid class name, expected a constant", constant.Location())
		}
	}
	var location *position.Location

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeParameterList(token.RBRACKET)
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Location(),
					errTok,
				)
			}
		}
	}

	if p.match(token.LESS) {
		superclass = p.genericConstantOrNil()
	}

	lastLocation, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = classTok.Location().Join(lastLocation)
	} else {
		location = classTok.Location()
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
			location = location.Join(endTok.Location())
		}
	}

	if !allowed {
		p.errorMessageLocation(
			"class declarations cannot appear in expressions",
			location,
		)
	}

	if constant == nil {
		p.errorMessageLocation(
			"anonymous classes are not supported",
			location,
		)
	}

	return ast.NewClassDeclarationNode(
		location,
		"",
		false,
		false,
		false,
		false,
		constant,
		typeVars,
		superclass,
		thenBody,
	)
}

// moduleDeclaration = "module" [constantLookup] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) moduleDeclaration(allowed bool) ast.ExpressionNode {
	moduleTok := p.advance()
	var constant ast.ExpressionNode
	if !p.accept(token.LBRACKET, token.SEMICOLON, token.NEWLINE) {
		constant = p.constantLookup()

		switch c := constant.(type) {
		case *ast.PublicConstantNode,
			*ast.PrivateConstantNode,
			*ast.ConstantLookupNode:
		case *ast.UnquoteNode:
			c.Kind = ast.UNQUOTE_CONSTANT_KIND
		default:
			p.errorMessageLocation("invalid module name, expected a constant", constant.Location())
		}
	}
	var location *position.Location

	if lbracket, ok := p.matchOk(token.LBRACKET); ok {
		errPos := lbracket.Location()
		if p.accept(token.RBRACKET) {
			rbracket := p.advance()
			errPos = errPos.Join(rbracket.Location())
		} else {
			p.typeParameterList()
			rbracket, ok := p.consume(token.RBRACKET)
			if !ok {
				return ast.NewInvalidNode(
					rbracket.Location(),
					rbracket,
				)
			}
			errPos = errPos.Join(rbracket.Location())
		}
		p.errorMessageLocation("modules cannot be generic", errPos)
	}

	lastLocation, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = moduleTok.Location().Join(lastLocation)
	} else {
		location = moduleTok.Location()
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
			location = location.Join(endTok.Location())
		}
	}

	if !allowed {
		p.errorMessageLocation(
			"module declarations cannot appear in expressions",
			location,
		)
	}

	if constant == nil {
		p.errorMessageLocation(
			"anonymous modules are not supported",
			location,
		)
	}

	return ast.NewModuleDeclarationNode(
		location,
		"",
		constant,
		thenBody,
	)
}

// mixinDeclaration = "mixin" [constantLookup] ["[" typeVariableList "]"] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) mixinDeclaration(allowed bool) ast.ExpressionNode {
	mixinTok := p.advance()
	var constant ast.ExpressionNode
	var typeVars []ast.TypeParameterNode
	if !p.accept(token.LBRACKET, token.SEMICOLON, token.NEWLINE) {
		constant = p.constantLookup()

		switch c := constant.(type) {
		case *ast.PublicConstantNode,
			*ast.PrivateConstantNode,
			*ast.ConstantLookupNode:
		case *ast.UnquoteNode:
			c.Kind = ast.UNQUOTE_CONSTANT_KIND
		default:
			p.errorMessageLocation("invalid mixin name, expected a constant", constant.Location())
		}
	}
	var location *position.Location

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeParameterList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Location(),
					errTok,
				)
			}
		}
	}

	lastLocation, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = mixinTok.Location().Join(lastLocation)
	} else {
		location = mixinTok.Location()
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
			location = location.Join(endTok.Location())
		}
	}

	if !allowed {
		p.errorMessageLocation(
			"mixin declarations cannot appear in expressions",
			location,
		)
	}

	if constant == nil {
		p.errorMessageLocation(
			"anonymous mixins are not supported",
			location,
		)
	}

	return ast.NewMixinDeclarationNode(
		location,
		"",
		false,
		constant,
		typeVars,
		thenBody,
	)
}

// interfaceDeclaration = "interface" [constantLookup] ["[" typeVariableList "]"] ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) interfaceDeclaration(allowed bool) ast.ExpressionNode {
	interfaceTok := p.advance()
	var constant ast.ExpressionNode
	var typeVars []ast.TypeParameterNode
	if !p.accept(token.LBRACKET, token.SEMICOLON, token.NEWLINE) {
		constant = p.constantLookup()

		switch c := constant.(type) {
		case *ast.PublicConstantNode,
			*ast.PrivateConstantNode,
			*ast.ConstantLookupNode:
		case *ast.UnquoteNode:
			c.Kind = ast.UNQUOTE_CONSTANT_KIND
		default:
			p.errorMessageLocation("invalid interface name, expected a constant", constant.Location())
		}
	}
	var location *position.Location

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeParameterList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Location(),
					errTok,
				)
			}
		}
	}

	lastLocation, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = interfaceTok.Location().Join(lastLocation)
	} else {
		location = interfaceTok.Location()
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
			location = location.Join(endTok.Location())
		}
	}

	if !allowed {
		p.errorMessageLocation(
			"interface declarations cannot appear in expressions",
			location,
		)
	}

	if constant == nil {
		p.errorMessageLocation(
			"anonymous interfaces are not supported",
			location,
		)
	}

	return ast.NewInterfaceDeclarationNode(
		location,
		"",
		constant,
		typeVars,
		thenBody,
	)
}

// structDeclaration = "struct" [constantLookup] ["[" typeVariableList "]"] ((SEPARATOR [structBodyStatements] "end") | ("then" formalParameter))
func (p *Parser) structDeclaration(allowed bool) ast.ExpressionNode {
	structTok := p.advance()
	var constant ast.ExpressionNode
	var typeVars []ast.TypeParameterNode
	if !p.accept(token.LBRACKET, token.SEMICOLON, token.NEWLINE) {
		constant = p.constantLookup()

		switch c := constant.(type) {
		case *ast.PublicConstantNode,
			*ast.PrivateConstantNode,
			*ast.ConstantLookupNode:
		case *ast.UnquoteNode:
			c.Kind = ast.UNQUOTE_CONSTANT_KIND
		default:
			p.errorMessageLocation("invalid struct name, expected a constant", constant.Location())
		}
	}
	var location *position.Location

	if p.match(token.LBRACKET) {
		if p.accept(token.RBRACKET) {
			p.errorExpected("a list of type variables")
			p.advance()
		} else {
			typeVars = p.typeParameterList()
			if errTok, ok := p.consume(token.RBRACKET); !ok {
				return ast.NewInvalidNode(
					errTok.Location(),
					errTok,
				)
			}
		}
	}

	lastLocation, thenBody, multiline := p.structBodyStatementBlockWithThen(token.END)
	if lastLocation != nil {
		location = structTok.Location().Join(lastLocation)
	} else {
		location = structTok.Location()
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
			location = location.Join(endTok.Location())
		}
	}

	if !allowed {
		p.errorMessageLocation(
			"struct declarations cannot appear in expressions",
			location,
		)
	}

	if constant == nil {
		p.errorMessageLocation(
			"anonymous structs are not supported",
			location,
		)
	}

	return ast.NewStructDeclarationNode(
		location,
		"",
		constant,
		typeVars,
		thenBody,
	)
}

// variableDeclaration = "var" identifier [":" typeAnnotationWithoutVoid] ["=" expressionWithoutModifier] |
// "var" pattern "=" expressionWithoutModifier
func (p *Parser) variableDeclaration(instanceVariableAllowed bool) ast.ExpressionNode {
	varTok := p.advance()
	var init ast.ExpressionNode

	if varName, ok := p.matchOk(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER,
		token.INSTANCE_VARIABLE); ok {
		var typ ast.TypeNode
		lastLocation := varName.Location()

		if p.match(token.COLON) {
			typ = p.typeAnnotationWithoutVoid()
			lastLocation = typ.Location()
		}

		if p.match(token.EQUAL_OP) {
			p.swallowNewlines()
			init = p.expressionWithoutModifier()
			lastLocation = init.Location()
			if varName.Type == token.INSTANCE_VARIABLE {
				p.errorMessageLocation("instance variables cannot be initialised when declared", lastLocation)
			}
		}

		location := varTok.Location().Join(lastLocation)
		if varName.Type == token.INSTANCE_VARIABLE {
			if !instanceVariableAllowed {
				p.errorMessageLocation(
					"instance variable declarations cannot appear in expressions",
					location,
				)
			}
			if typ == nil {
				p.errorMessageLocation(
					"instance variable declarations must have an explicit type",
					location,
				)
			}
			return ast.NewInstanceVariableDeclarationNode(
				location,
				"",
				varName.Value,
				typ,
			)
		}

		return ast.NewVariableDeclarationNode(
			location,
			"",
			varName.Value,
			typ,
			init,
		)
	}

	pattern := p.pattern()
	if tok, ok := p.consume(token.EQUAL_OP); !ok {
		return ast.NewInvalidNode(
			tok.Location(),
			tok,
		)
	}
	p.swallowNewlines()
	init = p.expressionWithoutModifier()

	if !ast.PatternDeclaresVariables(pattern) {
		p.errorMessageLocation("patterns in variable declarations should define at least one variable", pattern.Location())
	}

	return ast.NewVariablePatternDeclarationNode(
		varTok.Location().Join(init.Location()),
		pattern,
		init,
	)
}

// valueDeclaration = "val" identifier [":" typeAnnotationWithoutVoid] ["=" expressionWithoutModifier] |
// "val" pattern "=" expressionWithoutModifier
func (p *Parser) valueDeclaration() ast.ExpressionNode {
	valTok := p.advance()
	var init ast.ExpressionNode

	if valName, ok := p.matchOk(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER,
		token.INSTANCE_VARIABLE); ok {
		var typ ast.TypeNode
		lastLocation := valName.Location()

		switch valName.Type {
		case token.INSTANCE_VARIABLE:
			p.errorMessageLocation("instance variables cannot be declared using `val`", valName.Location())
		}

		if p.match(token.COLON) {
			typ = p.typeAnnotationWithoutVoid()
			lastLocation = typ.Location()
		}

		if p.match(token.EQUAL_OP) {
			p.swallowNewlines()
			init = p.expressionWithoutModifier()
			lastLocation = init.Location()
		}

		return ast.NewValueDeclarationNode(
			valTok.Location().Join(lastLocation),
			valName.Value,
			typ,
			init,
		)
	}

	pattern := p.pattern()
	if tok, ok := p.consume(token.EQUAL_OP); !ok {
		return ast.NewInvalidNode(
			tok.Location(),
			tok,
		)
	}
	p.swallowNewlines()
	init = p.expressionWithoutModifier()

	if !ast.PatternDeclaresVariables(pattern) {
		p.errorMessageLocation("patterns in value declarations should define at least one value", pattern.Location())
	}

	return ast.NewValuePatternDeclarationNode(
		valTok.Location().Join(init.Location()),
		pattern,
		init,
	)
}

// constantDeclaration = "const" complexConstant [":" typeAnnotationWithoutVoid] ["=" expressionWithoutModifier]
func (p *Parser) constantDeclaration(allowed bool) ast.ExpressionNode {
	constTok := p.advance()
	var init ast.ExpressionNode
	var typ ast.TypeNode

	constant := p.constantLookup()
	switch c := constant.(type) {
	case *ast.PublicConstantNode,
		*ast.PrivateConstantNode,
		*ast.ConstantLookupNode:
	case *ast.UnquoteNode:
		c.Kind = ast.UNQUOTE_CONSTANT_KIND
	default:
		p.errorMessageLocation("invalid constant name", constant.Location())
	}
	lastLocation := constant.Location()

	if p.match(token.COLON) {
		typ = p.typeAnnotationWithoutVoid()
		lastLocation = typ.Location()
	}

	if p.match(token.EQUAL_OP) {
		p.swallowNewlines()
		init = p.expressionWithoutModifier()
		lastLocation = init.Location()
	}

	location := constTok.Location().Join(lastLocation)

	if !allowed {
		p.errorMessageLocation(
			"constant declarations cannot appear in expressions",
			location,
		)
	}

	return ast.NewConstantDeclarationNode(
		location,
		"",
		constant,
		typ,
		init,
	)
}

// typeAnnotation = "void" | "never" | "any" | unionType
func (p *Parser) typeAnnotation() ast.TypeNode {
	switch p.lookahead.Type {
	case token.VOID:
		tok := p.advance()
		return ast.NewVoidTypeNode(tok.Location())
	}
	return p.unionType()
}

// typeAnnotationWithoutUnionAndVoid = "any" | intersectionType
func (p *Parser) typeAnnotationWithoutUnionAndVoid() ast.TypeNode {
	switch p.lookahead.Type {
	case token.VOID:
		p.errorMessage("type `void` cannot be used in this context")
		tok := p.advance()
		return ast.NewVoidTypeNode(tok.Location())
	}
	return p.intersectionType()
}

// typeAnnotationWithoutVoid = "any" | unionType
func (p *Parser) typeAnnotationWithoutVoid() ast.TypeNode {
	switch p.lookahead.Type {
	case token.VOID:
		p.errorMessage("type `void` cannot be used in this context")
		tok := p.advance()
		return ast.NewVoidTypeNode(tok.Location())
	}
	return p.unionType()
}

// unionType = intersectionType | unionType "|" intersectionType
func (p *Parser) unionType() ast.TypeNode {
	if p.mode == withoutUnionTypeMode {
		return p.intersectionType()
	}
	return p.binaryTypeExpression(p.intersectionType, token.OR)
}

// intersectionType = differenceType | intersectionType "&" differenceType
func (p *Parser) intersectionType() ast.TypeNode {
	return p.binaryTypeExpression(p.differenceType, token.AND)
}

// differenceType = unaryType | differenceType "/" unaryType
func (p *Parser) differenceType() ast.TypeNode {
	return p.binaryTypeExpression(p.unaryType, token.SLASH)
}

// unaryType = ("~" | "&" | "%") unaryType | nilableType
func (p *Parser) unaryType() ast.TypeNode {
	switch p.lookahead.Type {
	case token.TILDE:
		opTok := p.advance()
		typ := p.unaryType()
		return ast.NewNotTypeNode(
			opTok.Location().Join(typ.Location()),
			typ,
		)
	case token.AND:
		opTok := p.advance()
		typ := p.unaryType()
		return ast.NewSingletonTypeNode(
			opTok.Location().Join(typ.Location()),
			typ,
		)
	case token.PERCENT:
		opTok := p.advance()
		typ := p.unaryType()
		return ast.NewInstanceOfTypeNode(
			opTok.Location().Join(typ.Location()),
			typ,
		)
	default:
		return p.nilableType()
	}
}

// nilableType = unaryLiteralType ["?"]
func (p *Parser) nilableType() ast.TypeNode {
	typ := p.unaryLiteralType()

	if questTok, ok := p.matchOk(token.QUESTION); ok {
		return ast.NewNilableTypeNode(
			typ.Location().Join(questTok.Location()),
			typ,
		)
	}

	return typ
}

// unaryLiteralType = ("-" | "+") unaryLiteralType | primaryType
func (p *Parser) unaryLiteralType() ast.TypeNode {
	switch p.lookahead.Type {
	case token.PLUS, token.MINUS:
		opTok := p.advance()
		typ := p.unaryLiteralType()
		return ast.NewUnaryTypeNode(
			opTok.Location().Join(typ.Location()),
			opTok,
			typ,
		)
	default:
		return p.primaryType()
	}
}

// primaryType = namedType | "(" typeAnnotation ")"
func (p *Parser) primaryType() ast.TypeNode {
	if p.match(token.LPAREN) {
		if p.mode == withoutUnionTypeMode {
			p.mode = normalMode
		}
		t := p.typeAnnotation()
		p.consume(token.RPAREN)
		return t
	}

	switch p.lookahead.Type {
	case token.UNQUOTE:
		return p.unquoteType()
	case token.BOOL:
		tok := p.advance()
		return ast.NewBoolLiteralNode(tok.Location())
	case token.SELF:
		return p.selfLiteral()
	case token.TRUE:
		return p.trueLiteral()
	case token.FALSE:
		return p.falseLiteral()
	case token.NIL:
		return p.nilLiteral()
	case token.CHAR_LITERAL:
		return p.charLiteral()
	case token.RAW_CHAR_LITERAL:
		return p.rawCharLiteral()
	case token.RAW_STRING:
		return p.rawStringLiteral()
	case token.STRING_BEG:
		return p.stringLiteral(false)
	case token.COLON:
		return p.symbolLiteral(false)
	case token.INT:
		return p.int()
	case token.INT64:
		return p.int64()
	case token.UINT64:
		return p.uint64()
	case token.INT32:
		return p.int32()
	case token.UINT32:
		return p.uint32()
	case token.INT16:
		return p.int16()
	case token.UINT16:
		return p.uint16()
	case token.INT8:
		return p.int8()
	case token.UINT8:
		return p.uint8()
	case token.FLOAT:
		return p.float()
	case token.BIG_FLOAT:
		return p.bigFloat()
	case token.FLOAT64:
		return p.float64()
	case token.FLOAT32:
		return p.float32()
	case token.NEVER:
		tok := p.advance()
		return ast.NewNeverTypeNode(tok.Location())
	case token.ANY:
		tok := p.advance()
		return ast.NewAnyTypeNode(tok.Location())
	case token.OR, token.OR_OR:
		return p.closureType()
	default:
		return p.namedType()
	}
}

// closureType = (("|" signatureParameterList "|") | "||") [: typeAnnotation] ["!" typeAnnotation]
func (p *Parser) closureType() ast.TypeNode {
	var params []ast.ParameterNode
	var location *position.Location
	var returnType ast.TypeNode
	var throwType ast.TypeNode

	if p.accept(token.OR) {
		location = p.advance().Location()
		if !p.accept(token.OR) {
			p.mode = withoutUnionTypeMode
			params = p.signatureParameterList(token.OR)
			p.mode = normalMode
			location = position.JoinLocationOfLastElement(location, params)
		}
		if tok, ok := p.consume(token.OR); !ok {
			return ast.NewInvalidNode(
				tok.Location(),
				tok,
			)
		}
	} else {
		orOr, ok := p.consume(token.OR_OR)
		location = orOr.Location()
		if !ok {
			return ast.NewInvalidNode(
				orOr.Location(),
				orOr,
			)
		}
	}

	// return type
	if p.match(token.COLON) {
		returnType = p.typeAnnotation()
		location = location.Join(returnType.Location())
	}

	// throw type
	if p.match(token.BANG) {
		throwType = p.typeAnnotation()
		location = location.Join(throwType.Location())
	}

	return ast.NewClosureTypeNode(
		location,
		params,
		returnType,
		throwType,
	)
}

// namedType = genericConstant
func (p *Parser) namedType() ast.TypeNode {
	return p.genericConstant()
}

// genericConstantOrNil = genericConstant | "nil"
func (p *Parser) genericConstantOrNil() ast.ComplexConstantNode {
	if p.accept(token.NIL) {
		return p.nilLiteral()
	}
	return p.genericConstant()
}

// genericConstant = strictConstantLookup | strictConstantLookup "[" [typeAnnotationList] "]"
func (p *Parser) genericConstant() ast.ComplexConstantNode {
	constant := p.strictConstantLookup()
	if !p.match(token.LBRACKET) {
		return constant
	}

	if p.match(token.RBRACKET) {
		p.errorExpected("a constant")
		return constant
	}

	typeList := p.typeAnnotationList(token.RBRACKET)
	rbracket, ok := p.consume(token.RBRACKET)
	if !ok {
		return ast.NewInvalidNode(rbracket.Location(), rbracket)
	}

	return ast.NewGenericConstantNode(
		constant.Location().Join(rbracket.Location()),
		constant,
		typeList,
	)
}

// throwExpression = "%/" (REGEX_CONTENT | "%{" expressionWithoutModifier "}")* "/" REGEX_FLAG*
func (p *Parser) regexLiteral() ast.RegexLiteralNode {
	begTok := p.advance()
	var endTok *token.Token

	var reContent []ast.RegexLiteralContentNode
	for {
		if tok, ok := p.matchOk(token.REGEX_CONTENT); ok {
			reContent = append(reContent, ast.NewRegexLiteralContentSectionNode(
				tok.Location(),
				tok.Value,
			))
			continue
		}

		if beg, ok := p.matchOk(token.REGEX_INTERP_BEG); ok {
			expr := p.expressionWithoutModifier()
			end, _ := p.consume(token.REGEX_INTERP_END)
			reContent = append(reContent, ast.NewRegexInterpolationNode(
				beg.Location().Join(end.Location()),
				expr,
			))
			continue
		}

		tok, ok := p.consume(token.REGEX_END)
		endTok = tok
		if tok.Type == token.END_OF_FILE {
			break
		}
		if !ok {
			reContent = append(reContent, ast.NewInvalidNode(
				tok.Location(),
				tok,
			))
			continue
		}
		break
	}

	var flags bitfield.BitField8
tokenLoop:
	for {
		switch p.lookahead.Type {
		case token.REGEX_FLAG_i:
			flags.SetFlag(flag.CaseInsensitiveFlag)
		case token.REGEX_FLAG_m:
			flags.SetFlag(flag.MultilineFlag)
		case token.REGEX_FLAG_s:
			flags.SetFlag(flag.DotAllFlag)
		case token.REGEX_FLAG_x:
			flags.SetFlag(flag.ExtendedFlag)
		case token.REGEX_FLAG_U:
			flags.SetFlag(flag.UngreedyFlag)
		case token.REGEX_FLAG_a:
			flags.SetFlag(flag.ASCIIFlag)
		default:
			break tokenLoop
		}
		endTok = p.advance()
	}

	if len(reContent) == 0 {
		return ast.NewUninterpolatedRegexLiteralNode(begTok.Location().Join(endTok.Location()), "", flags)
	}
	reVal, ok := reContent[0].(*ast.RegexLiteralContentSectionNode)
	if len(reContent) == 1 && ok {
		return ast.NewUninterpolatedRegexLiteralNode(
			begTok.Location().Join(endTok.Location()),
			reVal.Value,
			flags,
		)
	}

	return ast.NewInterpolatedRegexLiteralNode(
		begTok.Location().Join(endTok.Location()),
		reContent,
		flags,
	)
}

// throwExpression = "throw" ["unchecked"] [expressionWithoutModifier]
func (p *Parser) throwExpression() *ast.ThrowExpressionNode {
	throwTok := p.advance()
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() {
		return ast.NewThrowExpressionNode(
			throwTok.Location(),
			false,
			nil,
		)
	}

	var unchecked bool
	if p.match(token.UNCHECKED) {
		unchecked = true
	}

	expr := p.expressionWithoutModifier()

	return ast.NewThrowExpressionNode(
		throwTok.Location().Join(expr.Location()),
		unchecked,
		expr,
	)
}

// mustExpression = "must" [expressionWithoutModifier]
func (p *Parser) mustExpression() *ast.MustExpressionNode {
	mustTok := p.advance()
	expr := p.expressionWithoutModifier()

	return ast.NewMustExpressionNode(
		mustTok.Location().Join(expr.Location()),
		expr,
	)
}

// unquoteExpression = "unquote" "(" expressionWithoutModifier ")"
func (p *Parser) unquoteExpression() ast.ExpressionNode {
	return p.unquote(ast.UNQUOTE_EXPRESSION_KIND)
}

// unquotePatternExpression = "unquote" "(" expressionWithoutModifier ")"
func (p *Parser) unquotePattern() ast.PatternNode {
	return p.unquote(ast.UNQUOTE_PATTERN_KIND)
}

// unquotePattern = "unquote" "(" expressionWithoutModifier ")"
func (p *Parser) unquotePatternExpression() ast.PatternExpressionNode {
	return p.unquote(ast.UNQUOTE_PATTERN_EXPRESSION_KIND)
}

// unquoteConstant = "unquote" "(" expressionWithoutModifier ")"
func (p *Parser) unquoteConstant() ast.ConstantNode {
	return p.unquote(ast.UNQUOTE_CONSTANT_KIND)
}

// unquoteType = "unquote" "(" expressionWithoutModifier ")"
func (p *Parser) unquoteType() ast.TypeNode {
	return p.unquote(ast.UNQUOTE_TYPE_KIND)
}

// unquoteIdentifier = "unquote" "(" expressionWithoutModifier ")"
func (p *Parser) unquoteIdentifier() ast.IdentifierNode {
	return p.unquote(ast.UNQUOTE_IDENTIFIER_KIND)
}

// unquote = "unquote" "(" expressionWithoutModifier ")"
func (p *Parser) unquote(kind ast.UnquoteKind) ast.UnquoteOrInvalidNode {
	unquoteTok := p.advance()

	lparen, ok := p.consume(token.LPAREN)
	if !ok {
		return ast.NewInvalidNode(lparen.Location(), lparen)
	}

	expr := p.expressionWithoutModifier()

	rparen, ok := p.consume(token.RPAREN)
	if !ok {
		return ast.NewInvalidNode(rparen.Location(), rparen)
	}

	location := unquoteTok.Location().Join(rparen.Location())

	return ast.NewUnquoteNode(
		location,
		kind,
		expr,
	)
}

// quoteExpression = "quote" ((SEPARATOR [statements]) "end" | (expressionWithoutModifier))
func (p *Parser) quoteExpression() ast.ExpressionNode {
	quoteTok := p.advance()

	lastLocation, body, multiline := p.statementBlock(token.END)

	var location *position.Location
	if lastLocation != nil {
		location = quoteTok.Location().Join(lastLocation)
	} else {
		location = quoteTok.Location()
	}

	quoteExpression := ast.NewQuoteExpressionNode(
		location,
		body,
	)

	if multiline {
		if len(body) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(body) == 0 {
			p.indentedSection = false
		}
		if ok {
			quoteExpression.SetLocation(quoteExpression.Location().Join(endTok.Location()))
		}
	}

	return quoteExpression
}

// macroBoundary = "do" "macro" [RAW_STRING] ((SEPARATOR [statements]) "end" | (expressionWithoutModifier))
func (p *Parser) macroBoundary() ast.ExpressionNode {
	doTok := p.advance()
	p.advance() // macro

	var name string
	if p.accept(token.RAW_STRING) {
		nameTok := p.advance()
		name = nameTok.Value
	}

	lastLocation, body, multiline := p.statementBlock(token.END)

	var location *position.Location
	if lastLocation != nil {
		location = doTok.Location().Join(lastLocation)
	} else {
		location = doTok.Location()
	}

	macroBoundary := ast.NewMacroBoundaryNode(
		location,
		body,
		name,
	)

	if multiline {
		if len(body) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(body) == 0 {
			p.indentedSection = false
		}
		if ok {
			macroBoundary.SetLocation(macroBoundary.Location().Join(endTok.Location()))
		}
	}

	return macroBoundary
}

// tryExpression = "try" expressionWithoutModifier
func (p *Parser) tryExpression() *ast.TryExpressionNode {
	tryTok := p.advance()
	expr := p.expressionWithoutModifier()

	return ast.NewTryExpressionNode(
		tryTok.Location().Join(expr.Location()),
		expr,
	)
}

// typeofExpression = "typeof" expressionWithoutModifier
func (p *Parser) typeofExpression() *ast.TypeofExpressionNode {
	mustTok := p.advance()
	expr := p.expressionWithoutModifier()

	return ast.NewTypeofExpressionNode(
		mustTok.Location().Join(expr.Location()),
		expr,
	)
}

// breakExpression = "break" [SPECIAL_IDENTIFIER] [expressionWithoutModifier]
func (p *Parser) breakExpression() *ast.BreakExpressionNode {
	breakTok := p.advance()
	location := breakTok.Location()
	var label string
	if p.lookahead.Type == token.SPECIAL_IDENTIFIER {
		labelTok := p.advance()
		label = labelTok.Value
		location = location.Join(labelTok.Location())
	}
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() || p.accept(token.IF, token.UNLESS) {
		return ast.NewBreakExpressionNode(
			location,
			label,
			nil,
		)
	}

	expr := p.expressionWithoutModifier()

	return ast.NewBreakExpressionNode(
		location.Join(expr.Location()),
		label,
		expr,
	)
}

// continueExpression = "continue" [SPECIAL_IDENTIFIER] [expressionWithoutModifier]
func (p *Parser) continueExpression() *ast.ContinueExpressionNode {
	continueTok := p.advance()
	location := continueTok.Location()
	var label string
	if p.lookahead.Type == token.SPECIAL_IDENTIFIER {
		labelTok := p.advance()
		label = labelTok.Value
		location = location.Join(labelTok.Location())
	}
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() || p.accept(token.IF, token.UNLESS) {
		return ast.NewContinueExpressionNode(
			location,
			label,
			nil,
		)
	}

	expr := p.expressionWithoutModifier()

	return ast.NewContinueExpressionNode(
		location.Join(expr.Location()),
		label,
		expr,
	)
}

// goExpression = "go" expressionWithoutModifier
func (p *Parser) goExpression() *ast.GoExpressionNode {
	goTok := p.advance()
	lastLocation, body, multiline := p.statementBlock(token.END)

	var location *position.Location
	if lastLocation != nil {
		location = goTok.Location().Join(lastLocation)
	} else {
		location = goTok.Location()
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
			location.Join(endTok.Location())
		}
	}

	return ast.NewGoExpressionNode(
		location,
		body,
	)
}

// awaitExpression = "await" expressionWithoutModifier
func (p *Parser) awaitExpression() *ast.AwaitExpressionNode {
	awaitTok := p.advance()
	expr := p.expressionWithoutModifier()

	return ast.NewAwaitExpressionNode(
		awaitTok.Location().Join(expr.Location()),
		expr,
	)
}

// returnExpression = "return" [expressionWithoutModifier]
func (p *Parser) returnExpression() *ast.ReturnExpressionNode {
	returnTok := p.advance()
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() || p.accept(token.IF, token.UNLESS) {
		return ast.NewReturnExpressionNode(
			returnTok.Location(),
			nil,
		)
	}

	expr := p.expressionWithoutModifier()

	return ast.NewReturnExpressionNode(
		returnTok.Location().Join(expr.Location()),
		expr,
	)
}

// yieldExpression = "yield" ["*"] [expressionWithoutModifier]
func (p *Parser) yieldExpression() *ast.YieldExpressionNode {
	yieldTok := p.advance()
	if p.lookahead.IsStatementSeparator() || p.lookahead.IsEndOfFile() || p.accept(token.IF, token.UNLESS) {
		return ast.NewYieldExpressionNode(
			yieldTok.Location(),
			false,
			nil,
		)
	}

	var forward bool
	if p.match(token.STAR) {
		forward = true
	}

	expr := p.expressionWithoutModifier()

	return ast.NewYieldExpressionNode(
		yieldTok.Location().Join(expr.Location()),
		forward,
		expr,
	)
}

// loopExpression = "loop" ((SEPARATOR [statements] "end") | expressionWithoutModifier)
func (p *Parser) loopExpression() *ast.LoopExpressionNode {
	loopTok := p.advance()
	var location *position.Location

	lastLocation, thenBody, multiline := p.statementBlock(token.END)
	if lastLocation != nil {
		location = loopTok.Location().Join(lastLocation)
	} else {
		location = loopTok.Location()
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
			location = location.Join(endTok.Location())
		}
	}

	return ast.NewLoopExpressionNode(
		location,
		thenBody,
	)
}

// asyncModifier = "async" declarationExpression
func (p *Parser) asyncModifier(allowed bool) ast.ExpressionNode {
	asyncTok := p.advance()

	p.swallowNewlines()
	node := p.declarationExpression(allowed)
	switch n := node.(type) {
	case *ast.MethodDefinitionNode:
		if n.IsAsync() {
			p.errorMessageLocation("the async modifier can only be attached once", asyncTok.Location())
		}
		n.SetAsync()
		n.SetLocation(asyncTok.Location().Join(n.Location()))
	default:
		p.errorMessageLocation("the async modifier can only be attached to methods", node.Location())
	}

	return node
}

// sealedModifier = "sealed" declarationExpression
func (p *Parser) sealedModifier(allowed bool) ast.ExpressionNode {
	sealedTok := p.advance()

	p.swallowNewlines()
	node := p.declarationExpression(allowed)
	switch n := node.(type) {
	case *ast.ClassDeclarationNode:
		if n.Sealed {
			p.errorMessageLocation("the sealed modifier can only be attached once", sealedTok.Location())
		}
		if n.Abstract {
			p.errorMessageLocation("the sealed modifier cannot be attached to abstract classes", sealedTok.Location())
		}
		n.Sealed = true
		n.SetLocation(sealedTok.Location().Join(n.Location()))
	case *ast.MacroDefinitionNode:
		if n.IsSealed() {
			p.errorMessageLocation("the sealed modifier can only be attached once", sealedTok.Location())
		}
		n.SetSealed()
		n.SetLocation(sealedTok.Location().Join(n.Location()))
	case *ast.MethodDefinitionNode:
		if n.IsSealed() {
			p.errorMessageLocation("the sealed modifier can only be attached once", sealedTok.Location())
		}
		if n.IsAbstract() {
			p.errorMessageLocation("the sealed modifier cannot be attached to abstract methods", sealedTok.Location())
		}
		n.SetSealed()
		n.SetLocation(sealedTok.Location().Join(n.Location()))
	default:
		p.errorMessageLocation("the sealed modifier can only be attached to classes and methods", node.Location())
	}

	return node
}

// abstractModifier = "abstract" declarationExpression
func (p *Parser) abstractModifier(allowed bool) ast.ExpressionNode {
	abstractTok := p.advance()

	p.swallowNewlines()
	node := p.declarationExpression(allowed)
	switch n := node.(type) {
	case *ast.ClassDeclarationNode:
		if n.Abstract {
			p.errorMessageLocation("the abstract modifier can only be attached once", abstractTok.Location())
		}
		if n.Sealed {
			p.errorMessageLocation("the abstract modifier cannot be attached to sealed classes", abstractTok.Location())
		}
		n.Abstract = true
		n.SetLocation(abstractTok.Location().Join(n.Location()))
	case *ast.MethodDefinitionNode:
		if n.IsAbstract() {
			p.errorMessageLocation("the abstract modifier can only be attached once", abstractTok.Location())
		}
		if n.IsSealed() {
			p.errorMessageLocation("the abstract modifier cannot be attached to sealed methods", abstractTok.Location())
		}
		n.SetAbstract()
		n.SetLocation(abstractTok.Location().Join(n.Location()))
	case *ast.MixinDeclarationNode:
		if n.Abstract {
			p.errorMessageLocation("the abstract modifier can only be attached once", abstractTok.Location())
		}
		n.Abstract = true
		n.SetLocation(abstractTok.Location().Join(n.Location()))
	default:
		p.errorMessageLocation("the abstract modifier can only be attached to classes, mixins and methods", node.Location())
	}

	return node
}

// noinitModifier = "noinit" declarationExpression
func (p *Parser) noinitModifier(allowed bool) ast.ExpressionNode {
	abstractTok := p.advance()

	p.swallowNewlines()
	node := p.declarationExpression(allowed)
	switch n := node.(type) {
	case *ast.ClassDeclarationNode:
		n.NoInit = true
		n.SetLocation(abstractTok.Location().Join(n.Location()))
	default:
		p.errorMessageLocation("the noinit modifier can only be attached to classes", node.Location())
	}

	return node
}

// primitiveModifier = "primitive" declarationExpression
func (p *Parser) primitiveModifier(allowed bool) ast.ExpressionNode {
	primitiveTok := p.advance()

	p.swallowNewlines()
	node := p.declarationExpression(allowed)
	switch n := node.(type) {
	case *ast.ClassDeclarationNode:
		if n.Primitive {
			p.errorMessageLocation("the primitive modifier can only be attached once", primitiveTok.Location())
		}
		n.Primitive = true
		n.SetLocation(primitiveTok.Location().Join(n.Location()))
	default:
		p.errorMessageLocation("the primitive modifier can only be attached to classes", node.Location())
	}

	return node
}

// fornumExpression = ("fornum" [expressionWithoutModifier] ";" [expressionWithoutModifier] ";" [expressionWithoutModifier])
// ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) fornumExpression() ast.ExpressionNode {
	forTok := p.advance()
	p.swallowNewlines()
	location := forTok.Location()
	var init, cond, incr ast.ExpressionNode
	if !p.accept(token.SEMICOLON) {
		init = p.expressionWithoutModifier()
	}
	semiTok, ok := p.consume(token.SEMICOLON)
	if !ok {
		return ast.NewInvalidNode(semiTok.Location(), semiTok)
	}

	if !p.accept(token.SEMICOLON) {
		cond = p.expressionWithoutModifier()
	}
	semiTok, ok = p.consume(token.SEMICOLON)
	if !ok {
		return ast.NewInvalidNode(semiTok.Location(), semiTok)
	}
	location = location.Join(semiTok.Location())

	if !p.accept(token.NEWLINE, token.SEMICOLON, token.THEN) {
		incr = p.expressionWithoutModifier()
	}

	lastLocation, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = location.Join(lastLocation)
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
			location = location.Join(endTok.Location())
		}
	}

	return ast.NewNumericForExpressionNode(
		location,
		init,
		cond,
		incr,
		thenBody,
	)
}

// forExpression = ("for" pattern "in" expressionWithoutModifier)
// ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) forExpression() ast.ExpressionNode {
	forTok := p.advance()
	p.swallowNewlines()
	parameter := p.pattern()
	if !ast.PatternDeclaresVariables(parameter) {
		p.errorMessageLocation("patterns in for in loops should define at least one variable", parameter.Location())
	}

	p.swallowNewlines()
	inTok, ok := p.consume(token.IN)
	if !ok {
		return ast.NewInvalidNode(inTok.Location(), inTok)
	}
	p.swallowNewlines()
	inExpr := p.expressionWithoutModifier()
	var location *position.Location

	lastLocation, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = forTok.Location().Join(lastLocation)
	} else {
		location = forTok.Location()
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
			location = location.Join(endTok.Location())
		}
	}

	return ast.NewForInExpressionNode(
		location,
		parameter,
		inExpr,
		thenBody,
	)
}

// whileExpression = "while" expressionWithoutModifier ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) whileExpression() *ast.WhileExpressionNode {
	whileTok := p.advance()
	cond := p.expressionWithoutModifier()
	var location *position.Location

	lastLocation, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = whileTok.Location().Join(lastLocation)
	} else {
		location = whileTok.Location()
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
			location = location.Join(endTok.Location())
		}
	}

	return ast.NewWhileExpressionNode(
		location,
		cond,
		thenBody,
	)
}

// untilExpression = "until" expressionWithoutModifier ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) untilExpression() *ast.UntilExpressionNode {
	untilTok := p.advance()
	cond := p.expressionWithoutModifier()
	var location *position.Location

	lastLocation, thenBody, multiline := p.statementBlockWithThen(token.END)
	if lastLocation != nil {
		location = untilTok.Location().Join(lastLocation)
	} else {
		location = untilTok.Location()
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
			location = location.Join(endTok.Location())
		}
	}

	return ast.NewUntilExpressionNode(
		location,
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
	var location *position.Location

	lastLocation, thenBody, multiline := p.statementBlockWithThen(token.END, token.ELSE)
	if lastLocation != nil {
		location = unlessTok.Location().Join(lastLocation)
	} else {
		location = unlessTok.Location()
	}

	unlessExpr := ast.NewUnlessExpressionNode(
		location,
		cond,
		thenBody,
		nil,
	)
	currentExpr := unlessExpr

	if p.lookahead.IsStatementSeparator() && p.secondLookahead.Type == token.ELSE {
		p.advance()
		p.advance()
		lastLocation, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastLocation != nil {
			currentExpr.SetLocation(currentExpr.Location().Join(lastLocation))
		}
	} else if p.match(token.ELSE) {
		lastLocation, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastLocation != nil {
			currentExpr.SetLocation(currentExpr.Location().Join(lastLocation))
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
			currentExpr.SetLocation(currentExpr.Location().Join(endTok.Location()))
		}
	}
	unlessExpr.SetLocation(unlessExpr.Location().Join(currentExpr.Location()))

	return unlessExpr
}

// doExpressionOrMacroBoundary = doExpression | macroBoundary
func (p *Parser) doExpressionOrMacroBoundary() ast.ExpressionNode {
	if p.acceptSecond(token.MACRO) {
		return p.macroBoundary()
	}

	return p.doExpression()
}

// doExpression = "do" ((SEPARATOR [statements]) | (expressionWithoutModifier))
// ("catch" pattern ((SEPARATOR [statements]) | ("then" expressionWithoutModifier)) )*
// ["finally" ((SEPARATOR [statements]) | expressionWithoutModifier)]
// "end"
func (p *Parser) doExpression() *ast.DoExpressionNode {
	doTok := p.advance()
	lastLocation, body, multiline := p.statementBlock(token.END, token.CATCH, token.FINALLY)

	var location *position.Location
	if lastLocation != nil {
		location = doTok.Location().Join(lastLocation)
	} else {
		location = doTok.Location()
	}

	doExpr := ast.NewDoExpressionNode(
		location,
		body,
		nil,
		nil,
	)

	for {
		var catchTok *token.Token

		if p.lookahead.Type == token.CATCH {
			catchTok = p.advance()
		} else if p.lookahead.IsStatementSeparator() && p.secondLookahead.Type == token.CATCH {
			p.advance()
			catchTok = p.advance()
		} else {
			break
		}
		pattern := p.pattern()
		var stackTraceVar ast.IdentifierNode
		if p.match(token.COMMA) {
			stackTraceVar = p.identifier()
		}

		lastLocation, body, multiline = p.statementBlockWithThen(token.END, token.CATCH, token.FINALLY)
		if lastLocation != nil {
			location = catchTok.Location().Join(lastLocation)
		} else {
			location = catchTok.Location()
		}

		catch := ast.NewCatchNode(
			location,
			pattern,
			stackTraceVar,
			body,
		)

		doExpr.Catches = append(
			doExpr.Catches,
			catch,
		)
	}

	if p.lookahead.IsStatementSeparator() && p.secondLookahead.Type == token.FINALLY {
		p.advance()
		p.advance()
		lastLocation, body, multiline = p.statementBlock(token.END)
		doExpr.Finally = body
		if lastLocation != nil {
			doExpr.SetLocation(doExpr.Location().Join(lastLocation))
		}
	} else if p.match(token.FINALLY) {
		lastLocation, body, multiline = p.statementBlock(token.END)
		doExpr.Finally = body
		if lastLocation != nil {
			doExpr.SetLocation(doExpr.Location().Join(lastLocation))
		}
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
			doExpr.SetLocation(doExpr.Location().Join(endTok.Location()))
		}
	}

	return doExpr
}

// singletonBlockExpression = "singleton" (expressionWithoutModifier | SEPARATOR statements "end")
func (p *Parser) singletonBlockExpression(allowed bool) *ast.SingletonBlockExpressionNode {
	singletonTok := p.advance()
	lastLocation, body, multiline := p.statementBlock(token.END)

	var location *position.Location
	if lastLocation != nil {
		location = singletonTok.Location().Join(lastLocation)
	} else {
		location = singletonTok.Location()
	}

	singletonBlockExpr := ast.NewSingletonBlockExpressionNode(
		location,
		body,
	)

	if multiline {
		if len(body) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(body) == 0 {
			p.indentedSection = false
		}
		if ok {
			singletonBlockExpr.SetLocation(singletonBlockExpr.Location().Join(endTok.Location()))
		}
	}

	if !allowed {
		p.errorMessageLocation(
			"singleton definitions cannot appear in expressions",
			location,
		)
	}

	return singletonBlockExpr
}

// singletonBlockExpression = "extend" "where" typeParameterListWithoutTerminator ((SEPARATOR [statements] "end") | ("then" expressionWithoutModifier))
func (p *Parser) extendWhereBlockExpression(allowed bool) ast.ExpressionNode {
	extendTok := p.advance()

	whereTok, ok := p.consume(token.WHERE)
	if !ok {
		return ast.NewInvalidNode(whereTok.Location(), whereTok)
	}
	typeParams := p.typeParameterListWithoutTerminator()
	location := position.JoinLocationOfLastElement(extendTok.Location(), typeParams)

	lastLocation, body, multiline := p.statementBlockWithThen(token.END)

	if lastLocation != nil {
		location = location.Join(lastLocation)
	}

	extendWhereBlock := ast.NewExtendWhereBlockExpressionNode(
		location,
		body,
		typeParams,
	)

	if multiline {
		if len(body) == 0 {
			p.indentedSection = true
		}
		endTok, ok := p.consume(token.END)
		if len(body) == 0 {
			p.indentedSection = false
		}
		if ok {
			extendWhereBlock.SetLocation(extendWhereBlock.Location().Join(endTok.Location()))
		}
	}

	if !allowed {
		p.errorMessageLocation(
			"extend where definitions cannot appear in expressions",
			location,
		)
	}

	return extendWhereBlock
}

// listElementPattern = ("*" [identifier]) | pattern
func (p *Parser) listElementPattern() ast.PatternNode {
	if star, ok := p.matchOk(token.STAR); ok {
		if p.accept(token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER) {
			ident := p.identifier()
			return ast.NewRestPatternNode(
				star.Location().Join(ident.Location()),
				ident,
			)
		}

		return ast.NewRestPatternNode(
			star.Location(),
			nil,
		)
	}
	return p.pattern()
}

// pattern = asPattern
func (p *Parser) pattern() ast.PatternNode {
	return p.asPattern()
}

// asPattern = orPattern ["as" identifier]
func (p *Parser) asPattern() ast.PatternNode {
	pattern := p.orPattern()
	if !p.match(token.AS) {
		return pattern
	}

	ident := p.identifier()

	return ast.NewAsPatternNode(
		pattern.Location().Join(ident.Location()),
		pattern,
		ident,
	)
}

// orPattern = andPattern | orPattern "||" andPattern
func (p *Parser) orPattern() ast.PatternNode {
	return p.binaryPattern(p.andPattern, token.OR_OR)
}

// andPattern = unaryPattern | andPattern "&&" unaryPattern
func (p *Parser) andPattern() ast.PatternNode {
	return p.binaryPattern(p.unaryPattern, token.AND_AND)
}

func (p *Parser) objectAttributePatternList(stopTokens ...token.Type) []ast.PatternNode {
	return commaSeparatedList(p, p.objectAttributePattern, stopTokens...)
}

// objectAttributePattern = (identifier | constant) |
// (identifier | constant) ":" pattern
func (p *Parser) objectAttributePattern() ast.PatternNode {
	if p.accept(
		token.PUBLIC_IDENTIFIER,
		token.PRIVATE_IDENTIFIER,
		token.PUBLIC_CONSTANT,
		token.PRIVATE_CONSTANT,
	) &&
		p.acceptSecond(token.COLON) {
		key := p.advance()
		p.advance()
		p.swallowNewlines()
		val := p.pattern()
		return ast.NewSymbolKeyValuePatternNode(
			key.Location().Join(val.Location()),
			key.Value,
			val,
		)
	}
	switch p.lookahead.Type {
	case token.PRIVATE_IDENTIFIER, token.PUBLIC_IDENTIFIER:
		return p.identifier()
	default:
		p.errorExpected("an object pattern attribute")
		tok := p.advance()
		return ast.NewInvalidNode(tok.Location(), tok)
	}
}

// strictConstantLookupOrObjectPattern = strictConstantLookup ["(" [objectPatternAttributes] ")"]
func (p *Parser) strictConstantLookupOrObjectPattern() ast.PatternNode {
	constant := p.strictConstantLookup()
	if !p.accept(token.LPAREN) {
		return constant
	}

	p.advance()
	p.swallowNewlines()

	if rparen, ok := p.matchOk(token.RPAREN); ok {
		return ast.NewObjectPatternNode(
			constant.Location().Join(rparen.Location()),
			constant,
			nil,
		)
	}

	elements := p.objectAttributePatternList(token.RPAREN)
	p.swallowNewlines()
	rparen, ok := p.consume(token.RPAREN)
	location := constant.Location()
	if ok {
		location = location.Join(rparen.Location())
	}

	return ast.NewObjectPatternNode(
		location,
		constant,
		elements,
	)
}

// unaryPattern = rangePattern |
// collectionPattern |
// ["<" | "<=" | ">" | ">=" | "==" | "!=" | "===" | "!==" | "=~" | "!~"] bitwiseOrExpression
func (p *Parser) unaryPattern() ast.PatternNode {
	if p.lookahead.IsCollectionLiteralBeg() {
		return p.collectionPattern()
	}

	if operator, ok := p.matchOk(token.LESS, token.LESS_EQUAL, token.GREATER,
		token.GREATER_EQUAL, token.EQUAL_EQUAL, token.NOT_EQUAL,
		token.STRICT_EQUAL, token.STRICT_NOT_EQUAL,
		token.LAX_EQUAL, token.LAX_NOT_EQUAL); ok {
		p.swallowNewlines()

		p.indentedSection = true
		right := p.bitwiseOrExpression()
		p.indentedSection = false

		return ast.NewUnaryExpressionNode(
			operator.Location().Join(right.Location()),
			operator,
			right,
		)
	}

	return p.rangePattern()
}

// collectionPattern = listPattern | tuplePattern | mapPattern | recordPattern | setPattern
func (p *Parser) collectionPattern() ast.PatternNode {
	switch p.lookahead.Type {
	case token.LBRACE:
		return p.mapPattern()
	case token.RECORD_LITERAL_BEG:
		return p.recordPattern()
	case token.LBRACKET:
		return p.listPattern()
	case token.WORD_ARRAY_LIST_BEG:
		return p.wordArrayListPattern()
	case token.BIN_ARRAY_LIST_BEG:
		return p.binArrayListPattern()
	case token.SYMBOL_ARRAY_LIST_BEG:
		return p.symbolArrayListPattern()
	case token.HEX_ARRAY_LIST_BEG:
		return p.hexArrayListPattern()
	case token.TUPLE_LITERAL_BEG:
		return p.tuplePattern()
	case token.WORD_ARRAY_TUPLE_BEG:
		return p.wordArrayTuplePattern()
	case token.SYMBOL_ARRAY_TUPLE_BEG:
		return p.symbolArrayTuplePattern()
	case token.HEX_ARRAY_TUPLE_BEG:
		return p.hexArrayTuplePattern()
	case token.BIN_ARRAY_TUPLE_BEG:
		return p.binArrayTuplePattern()
	case token.HASH_SET_LITERAL_BEG:
		return p.setPattern()
	case token.WORD_HASH_SET_BEG:
		return p.wordHashSetPattern()
	case token.SYMBOL_HASH_SET_BEG:
		return p.symbolHashSetPattern()
	case token.HEX_HASH_SET_BEG:
		return p.hexHashSetPattern()
	case token.BIN_HASH_SET_BEG:
		return p.binHashSetPattern()
	default:
		p.errorExpected("a collection pattern")
		tok := p.advance()
		return ast.NewInvalidNode(tok.Location(), tok)
	}
}

// mapPattern = "{" [mapLikePatternElements] "}"
func (p *Parser) mapPattern() ast.PatternNode {
	return collectionLiteralWithoutCapacity(p, token.RBRACE, p.mapLikePatternElements, ast.NewMapPatternNodeI)
}

// recordPattern = "%{" [mapLikePatternElements] "}"
func (p *Parser) recordPattern() ast.PatternNode {
	return collectionLiteralWithoutCapacity(p, token.RBRACE, p.mapLikePatternElements, ast.NewRecordPatternNodeI)
}

func (p *Parser) mapLikePatternElements(stopTokens ...token.Type) []ast.PatternNode {
	return commaSeparatedList(p, p.mapElementPattern, stopTokens...)
}

// mapElementPattern = (identifier) |
// (identifier | constant) ":" pattern |
// simplePattern "=>" pattern
func (p *Parser) mapElementPattern() ast.PatternNode {
	if p.accept(
		token.PUBLIC_IDENTIFIER,
		token.PRIVATE_IDENTIFIER,
		token.PUBLIC_CONSTANT,
		token.PRIVATE_CONSTANT,
	) &&
		p.acceptSecond(token.COLON) {
		key := p.advance()
		p.advance()
		p.swallowNewlines()
		val := p.pattern()
		return ast.NewSymbolKeyValuePatternNode(
			key.Location().Join(val.Location()),
			key.Value,
			val,
		)
	}
	key := p.simplePattern()
	if !p.match(token.THICK_ARROW) {
		switch key.(type) {
		case *ast.PublicIdentifierNode, *ast.PrivateIdentifierNode:
			return key
		default:
			p.errorMessageLocation("expected a key-value pair, map patterns should consist of key-value pairs", key.Location())
			return key
		}
	}

	p.swallowNewlines()
	val := p.pattern()

	return ast.NewKeyValuePatternNode(
		key.Location().Join(val.Location()),
		key,
		val,
	)
}

// setPattern = "^[" [patternElements] "]"
func (p *Parser) setPattern() ast.PatternNode {
	return collectionLiteralWithoutCapacity(p, token.RBRACKET, p.setPatternElements, ast.NewSetPatternNodeI)
}

// listPattern = "[" [listLikePatternElements] "]"
func (p *Parser) listPattern() ast.PatternNode {
	return collectionLiteralWithoutCapacity(p, token.RBRACKET, p.listLikePatternElements, ast.NewListPatternNodeI)
}

// tuplePattern = "%[" [listLikePatternElements] "]"
func (p *Parser) tuplePattern() ast.PatternNode {
	return collectionLiteralWithoutCapacity(p, token.RBRACKET, p.listLikePatternElements, ast.NewTuplePatternNodeI)
}

func (p *Parser) setPatternElements(stopTokens ...token.Type) []ast.PatternNode {
	return commaSeparatedList(p, p.setPatternElement, stopTokens...)
}

func (p *Parser) listLikePatternElements(stopTokens ...token.Type) []ast.PatternNode {
	return commaSeparatedList(p, p.listElementPattern, stopTokens...)
}

// rangePattern = primaryPattern |
// primaryPattern ("..." | "..<" | "<.." | "<.<") [primaryPattern] |
// ("..." | "..<" | "<.." | "<.<") primaryPattern
func (p *Parser) rangePattern() ast.PatternNode {
	if operator, ok := p.matchOk(token.CLOSED_RANGE_OP, token.OPEN_RANGE_OP,
		token.LEFT_OPEN_RANGE_OP, token.RIGHT_OPEN_RANGE_OP); ok {
		to := p.unaryPatternArgument()

		return ast.NewRangeLiteralNode(
			operator.Location().Join(to.Location()),
			operator,
			nil,
			to,
		)
	}

	var from ast.PatternNode
	if p.accept(token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT, token.SCOPE_RES_OP) {
		from = p.strictConstantLookupOrObjectPattern()
	} else {
		from = p.primaryPattern()
	}

	operator, ok := p.matchOk(token.CLOSED_RANGE_OP, token.OPEN_RANGE_OP, token.LEFT_OPEN_RANGE_OP, token.RIGHT_OPEN_RANGE_OP)
	if !ok {
		return from
	}

	fromExpr, ok := from.(ast.PatternExpressionNode)
	if !ok || !ast.IsValidRangePatternElement(from) {
		p.errorMessageLocation("invalid range pattern element", from.Location())
	}

	if !p.lookahead.IsValidAsEndInRangePattern() {
		return ast.NewRangeLiteralNode(
			from.Location().Join(operator.Location()),
			operator,
			fromExpr,
			nil,
		)
	}

	var to ast.PatternNode
	if p.accept(token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT, token.SCOPE_RES_OP) {
		to = p.strictConstantLookupOrObjectPattern()
	} else {
		to = p.unaryPatternArgument()
	}

	toExpr, ok := to.(ast.PatternExpressionNode)
	if !ok || !ast.IsValidRangePatternElement(to) {
		p.errorMessageLocation("invalid range pattern element", to.Location())
	}

	return ast.NewRangeLiteralNode(
		from.Location().Join(to.Location()),
		operator,
		fromExpr,
		toExpr,
	)
}

// unaryPatternArgument = ["-" | "+"] simplePattern
func (p *Parser) unaryPatternArgument() ast.PatternExpressionNode {
	operator, ok := p.matchOk(token.MINUS, token.PLUS)
	val := p.simplePattern()
	if !ok {
		return val
	}

	return ast.NewUnaryExpressionNode(
		operator.Location().Join(val.Location()),
		operator,
		val,
	)
}

// setPatternElement = literalPattern | "*" | "_"
func (p *Parser) setPatternElement() ast.PatternNode {
	switch p.lookahead.Type {
	case token.STAR:
		star := p.advance()
		var ident ast.IdentifierNode
		if p.accept(token.PRIVATE_IDENTIFIER, token.PUBLIC_IDENTIFIER) {
			p.errorMessage("set patterns do not support named rest elements")
			ident = p.identifier()
		}
		return ast.NewRestPatternNode(
			star.Location(),
			ident,
		)
	case token.PRIVATE_IDENTIFIER:
		ident := p.advance()
		if ident.Value != "_" {
			p.errorMessageLocation("set patterns cannot contain identifiers other than _", ident.Location())
		}
		return ast.NewPrivateIdentifierNode(
			ident.Location(),
			ident.Value,
		)
	case token.PUBLIC_IDENTIFIER:
		ident := p.advance()
		p.errorMessageLocation("set patterns cannot contain identifiers other than _", ident.Location())
		return ast.NewPublicIdentifierNode(
			ident.Location(),
			ident.Value,
		)
	default:
		return p.literalPattern()
	}
}

// literalPattern = ["-" | "+"] innerLiteralPattern
func (p *Parser) literalPattern() ast.PatternNode {
	operator, ok := p.matchOk(token.MINUS, token.PLUS)
	val := p.innerLiteralPattern()
	if !ok {
		return val
	}

	return ast.NewUnaryExpressionNode(
		operator.Location().Join(val.Location()),
		operator,
		val,
	)
}

func (p *Parser) innerLiteralPattern() ast.PatternExpressionNode {
	switch p.lookahead.Type {
	case token.PUBLIC_CONSTANT, token.PRIVATE_CONSTANT, token.SCOPE_RES_OP:
		return p.strictConstantLookup()
	case token.UNQUOTE:
		return p.unquotePatternExpression()
	case token.TRUE:
		return p.trueLiteral()
	case token.FALSE:
		return p.falseLiteral()
	case token.NIL:
		return p.nilLiteral()
	case token.CHAR_LITERAL:
		return p.charLiteral()
	case token.RAW_CHAR_LITERAL:
		return p.rawCharLiteral()
	case token.RAW_STRING:
		return p.rawStringLiteral()
	case token.STRING_BEG:
		return p.stringLiteral(true)
	case token.COLON:
		return p.symbolLiteral(true)
	case token.INT:
		return p.int()
	case token.INT64:
		return p.int64()
	case token.UINT64:
		return p.uint64()
	case token.INT32:
		return p.int32()
	case token.UINT32:
		return p.uint32()
	case token.INT16:
		return p.int16()
	case token.UINT16:
		return p.uint16()
	case token.INT8:
		return p.int8()
	case token.UINT8:
		return p.uint8()
	case token.FLOAT:
		return p.float()
	case token.BIG_FLOAT:
		return p.bigFloat()
	case token.FLOAT64:
		return p.float64()
	case token.FLOAT32:
		return p.float32()
	case token.ERROR:
		return p.invalidNode()
	}

	p.errorExpected("a pattern")
	p.updateErrorMode(true)
	tok := p.advance()
	return ast.NewInvalidNode(
		tok.Location(),
		tok,
	)
}

func (p *Parser) int() *ast.IntLiteralNode {
	tok := p.advance()
	return ast.NewIntLiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) int64() *ast.Int64LiteralNode {
	tok := p.advance()
	return ast.NewInt64LiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) uint64() *ast.UInt64LiteralNode {
	tok := p.advance()
	return ast.NewUInt64LiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) int32() *ast.Int32LiteralNode {
	tok := p.advance()
	return ast.NewInt32LiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) uint32() *ast.UInt32LiteralNode {
	tok := p.advance()
	return ast.NewUInt32LiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) int16() *ast.Int16LiteralNode {
	tok := p.advance()
	return ast.NewInt16LiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) uint16() *ast.UInt16LiteralNode {
	tok := p.advance()
	return ast.NewUInt16LiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) int8() *ast.Int8LiteralNode {
	tok := p.advance()
	return ast.NewInt8LiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) uint8() *ast.UInt8LiteralNode {
	tok := p.advance()
	return ast.NewUInt8LiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) float() *ast.FloatLiteralNode {
	tok := p.advance()
	return ast.NewFloatLiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) bigFloat() *ast.BigFloatLiteralNode {
	tok := p.advance()
	return ast.NewBigFloatLiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) float64() *ast.Float64LiteralNode {
	tok := p.advance()
	return ast.NewFloat64LiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) float32() *ast.Float32LiteralNode {
	tok := p.advance()
	return ast.NewFloat32LiteralNode(
		tok.Location(),
		tok.Value,
	)
}

func (p *Parser) invalidNode() *ast.InvalidNode {
	tok := p.advance()
	return ast.NewInvalidNode(
		tok.Location(),
		tok,
	)
}

func (p *Parser) trueLiteral() *ast.TrueLiteralNode {
	tok := p.advance()
	return ast.NewTrueLiteralNode(tok.Location())
}

func (p *Parser) falseLiteral() *ast.FalseLiteralNode {
	tok := p.advance()
	return ast.NewFalseLiteralNode(tok.Location())
}

func (p *Parser) nilLiteral() *ast.NilLiteralNode {
	tok := p.advance()
	return ast.NewNilLiteralNode(tok.Location())
}

func (p *Parser) simplePattern() ast.PatternExpressionNode {
	switch p.lookahead.Type {
	case token.REGEX_BEG:
		return p.regexLiteral()
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		return p.identifier()
	}

	return p.innerLiteralPattern()
}

// wordHashSetPattern = "^w[" (rawString)* "]"
func (p *Parser) wordHashSetPattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.wordCollectionElement,
		ast.NewWordHashSetLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.WORD_HASH_SET_END,
	)
}

// symbolHashSetPattern = "^s[" (rawString)* "]"
func (p *Parser) symbolHashSetPattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolHashSetLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.SYMBOL_HASH_SET_END,
	)
}

// hexHashSetPattern = "^x[" (HEX_INT)* "]"
func (p *Parser) hexHashSetPattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.intCollectionElement,
		ast.NewHexHashSetLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.HEX_HASH_SET_END,
	)
}

// binHashSetPattern = "^b[" (BIN_INT)* "]"
func (p *Parser) binHashSetPattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.intCollectionElement,
		ast.NewBinHashSetLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.BIN_HASH_SET_END,
	)
}

// wordArrayListPattern = "\w[" (rawString)* "]"
func (p *Parser) wordArrayListPattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.wordCollectionElement,
		ast.NewWordArrayListLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.WORD_ARRAY_LIST_END,
	)
}

// symbolArrayListPattern = "\s[" (rawString)* "]"
func (p *Parser) symbolArrayListPattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolArrayListLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.SYMBOL_ARRAY_LIST_END,
	)
}

// hexArrayListPattern = "\x[" (HEX_INT)* "]"
func (p *Parser) hexArrayListPattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.intCollectionElement,
		ast.NewHexArrayListLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.HEX_ARRAY_LIST_END,
	)
}

// binArrayListPattern = "\b[" (BIN_INT)* "]"
func (p *Parser) binArrayListPattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.intCollectionElement,
		ast.NewBinArrayListLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.BIN_ARRAY_LIST_END,
	)
}

// wordArrayTuplePattern = "%w[" (rawString)* "]"
func (p *Parser) wordArrayTuplePattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.wordCollectionElement,
		ast.NewWordArrayTupleLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.WORD_ARRAY_TUPLE_END,
	)
}

// symbolArrayTuplePattern = "%s[" (rawString)* "]"
func (p *Parser) symbolArrayTuplePattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.symbolCollectionElement,
		ast.NewSymbolArrayTupleLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.SYMBOL_ARRAY_TUPLE_END,
	)
}

// hexArrayTuplePattern = "%x[" (HEX_INT)* "]"
func (p *Parser) hexArrayTuplePattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.intCollectionElement,
		ast.NewHexArrayTupleLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.HEX_ARRAY_TUPLE_END,
	)
}

// binArrayTuplePattern = "%b[" (BIN_INT)* "]"
func (p *Parser) binArrayTuplePattern() ast.PatternExpressionNode {
	return specialCollectionLiteralWithoutCapacity(
		p,
		p.intCollectionElement,
		ast.NewBinArrayTupleLiteralPatternExpressionNode,
		ast.NewInvalidPatternExpressionNode,
		token.BIN_ARRAY_TUPLE_END,
	)
}

// primaryPattern = innerPrimaryPattern | ["-" | "+"] simplePattern
func (p *Parser) primaryPattern() ast.PatternNode {
	operator, ok := p.matchOk(token.MINUS, token.PLUS)
	if ok {
		val := p.simplePattern()
		return ast.NewUnaryExpressionNode(
			operator.Location().Join(val.Location()),
			operator,
			val,
		)
	}

	return p.innerPrimaryPattern()
}

func (p *Parser) innerPrimaryPattern() ast.PatternNode {
	switch p.lookahead.Type {
	case token.LPAREN:
		p.advance()
		pattern := p.pattern()
		p.consume(token.RPAREN)
		return pattern
	case token.UNQUOTE:
		return p.unquotePattern()
	default:
		return p.simplePattern()
	}
}

// switchExpression = "switch" expressionWithoutModifier SEPARATOR
// ("case" pattern ((SEPARATOR [statements]) | ("then" expressionWithoutModifier)) )*
// ["else" ((SEPARATOR [statements]) | expressionWithoutModifier)]
func (p *Parser) switchExpression() ast.ExpressionNode {
	switchTok := p.advance()
	val := p.expressionWithoutModifier()
	p.swallowNewlines()

	var lastLocation *position.Location
	var cases []*ast.CaseNode
	var els []ast.StatementNode
	var elsePresent bool
	withoutContent := true

	for {
		if p.match(token.ELSE) {
			lastLocation, els, _ = p.statementBlock(token.END)
			withoutContent = false
			elsePresent = true
			break
		} else if caseTok, ok := p.matchOk(token.CASE); ok {
			pattern := p.pattern()
			var caseBody []ast.StatementNode
			lastLocation, caseBody, _ = p.statementBlockWithThen(token.END, token.CASE, token.ELSE)
			withoutContent = false
			cases = append(cases, ast.NewCaseNode(
				caseTok.Location().Join(lastLocation),
				pattern,
				caseBody,
			))
			p.swallowNewlines()
		} else {
			break
		}
	}

	if withoutContent {
		p.indentedSection = true
	}
	p.swallowNewlines()
	endTok, ok := p.consume(token.END)
	if withoutContent {
		p.indentedSection = false
	}
	if ok {
		lastLocation = endTok.Location()
	}
	location := switchTok.Location().Join(lastLocation)
	if len(cases) == 0 && !elsePresent {
		p.errorMessageLocation("switch cannot be empty", location)
	} else if len(cases) == 0 && elsePresent {
		p.errorMessageLocation("switch cannot only consist of else", location)
	}
	return ast.NewSwitchExpressionNode(
		location,
		val,
		cases,
		els,
	)
}

// ifExpression = "if" expressionWithoutModifier ((SEPARATOR [statements]) | ("then" expressionWithoutModifier))
// ("elsif" expressionWithoutModifier ((SEPARATOR [statements]) | ("then" expressionWithoutModifier)) )*
// ["else" ((SEPARATOR [statements]) | expressionWithoutModifier)]
// "end"
func (p *Parser) ifExpression() *ast.IfExpressionNode {
	ifTok := p.advance()
	cond := p.expressionWithoutModifier()
	var location *position.Location

	lastLocation, thenBody, multiline := p.statementBlockWithThen(token.END, token.ELSE, token.ELSIF)
	if lastLocation != nil {
		location = ifTok.Location().Join(lastLocation)
	} else {
		location = ifTok.Location()
	}

	ifExpr := ast.NewIfExpressionNode(
		location,
		cond,
		thenBody,
		nil,
	)
	currentExpr := ifExpr

	for {
		var elsifTok *token.Token

		if p.lookahead.Type == token.ELSIF {
			elsifTok = p.advance()
		} else if p.lookahead.IsStatementSeparator() && p.secondLookahead.Type == token.ELSIF {
			p.advance()
			elsifTok = p.advance()
		} else {
			break
		}
		cond = p.expressionWithoutModifier()

		lastLocation, thenBody, multiline = p.statementBlockWithThen(token.END, token.ELSE, token.ELSIF)
		if lastLocation != nil {
			location = elsifTok.Location().Join(lastLocation)
		} else {
			location = elsifTok.Location()
		}

		elsifExpr := ast.NewIfExpressionNode(
			location,
			cond,
			thenBody,
			nil,
		)

		currentExpr.ElseBody = []ast.StatementNode{
			ast.NewExpressionStatementNode(
				elsifExpr.Location(),
				elsifExpr,
			),
		}
		currentExpr = elsifExpr
	}

	if p.lookahead.IsStatementSeparator() && p.secondLookahead.Type == token.ELSE {
		p.advance()
		p.advance()
		lastLocation, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastLocation != nil {
			currentExpr.SetLocation(currentExpr.Location().Join(lastLocation))
		}
	} else if p.match(token.ELSE) {
		lastLocation, thenBody, multiline = p.statementBlock(token.END)
		currentExpr.ElseBody = thenBody
		if lastLocation != nil {
			currentExpr.SetLocation(currentExpr.Location().Join(lastLocation))
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
			currentExpr.SetLocation(currentExpr.Location().Join(endTok.Location()))
		}
	}
	ifExpr.SetLocation(ifExpr.Location().Join(currentExpr.Location()))

	return ifExpr
}

// symbolLiteral = ":" (identifier | constant | rawStringLiteral)
func (p *Parser) symbolLiteral(withInterpolation bool) ast.StringOrSymbolTypeNode {
	symbolBegTok := p.advance()
	if p.lookahead.IsValidSimpleSymbolContent() {
		contTok := p.advance()
		return ast.NewSimpleSymbolLiteralNode(
			symbolBegTok.Location().Join(contTok.Location()),
			contTok.FetchValue(),
		)
	}

	if !p.accept(token.STRING_BEG) {
		p.errorExpected("an identifier, overridable operator or string literal")
		p.mode = panicMode
		tok := p.advance()
		return ast.NewInvalidNode(
			symbolBegTok.Location().Join(tok.Location()),
			tok,
		)
	}

	str := p.stringLiteral(withInterpolation)
	switch s := str.(type) {
	case *ast.DoubleQuotedStringLiteralNode:
		return ast.NewSimpleSymbolLiteralNode(
			symbolBegTok.Location().Join(s.Location()),
			s.Value,
		)
	case *ast.InterpolatedStringLiteralNode:
		return ast.NewInterpolatedSymbolLiteralNode(
			symbolBegTok.Location().Join(s.Location()),
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
		tok.Location(),
		char,
	)
}

// charLiteral = RAW_CHAR_LITERAL
func (p *Parser) rawCharLiteral() *ast.RawCharLiteralNode {
	tok := p.advance()
	char, _ := utf8.DecodeRuneInString(tok.Value)
	return ast.NewRawCharLiteralNode(
		tok.Location(),
		char,
	)
}

// rawStringLiteral = RAW_STRING
func (p *Parser) rawStringLiteral() *ast.RawStringLiteralNode {
	tok := p.advance()
	return ast.NewRawStringLiteralNode(
		tok.Location(),
		tok.Value,
	)
}

// wordCollectionElement = RAW_STRING
func (p *Parser) wordCollectionElement() ast.WordCollectionContentNode {
	tok, ok := p.consume(token.RAW_STRING)
	if !ok {
		return ast.NewInvalidNode(tok.Location(), tok)
	}
	return ast.NewRawStringLiteralNode(
		tok.Location(),
		tok.Value,
	)
}

// symbolCollectionElement = RAW_STRING
func (p *Parser) symbolCollectionElement() ast.SymbolCollectionContentNode {
	tok, ok := p.consume(token.RAW_STRING)
	if !ok {
		return ast.NewInvalidNode(tok.Location(), tok)
	}
	return ast.NewSimpleSymbolLiteralNode(
		tok.Location(),
		tok.Value,
	)
}

// intCollectionElement = INT
func (p *Parser) intCollectionElement() ast.IntCollectionContentNode {
	tok, ok := p.consume(token.INT)
	if !ok {
		return ast.NewInvalidNode(tok.Location(), tok)
	}
	return ast.NewIntLiteralNode(tok.Location(), tok.Value)
}

// stringLiteral = "\"" (STRING_CONTENT | "${" expressionWithoutModifier "}" | "#{" expressionWithoutModifier "}")* "\""
func (p *Parser) stringLiteral(withInterpolation bool) ast.StringTypeNode {
	quoteBeg := p.advance() // consume the opening quote
	var quoteEnd *token.Token

	var strContent []ast.StringLiteralContentNode
strContentLoop:
	for {
		switch p.lookahead.Type {
		case token.STRING_CONTENT:
			tok := p.advance()
			strContent = append(strContent, ast.NewStringLiteralContentSectionNode(
				tok.Location(),
				tok.Value,
			))
			continue strContentLoop
		case token.STRING_INTERP_BEG:
			beg := p.advance()
			expr := p.expressionWithoutModifier()
			end, _ := p.consume(token.STRING_INTERP_END)
			strContent = append(strContent, ast.NewStringInterpolationNode(
				beg.Location().Join(end.Location()),
				expr,
			))
			continue strContentLoop
		case token.STRING_INTERP_LOCAL:
			tok := p.advance()
			strContent = append(strContent, ast.NewStringInterpolationNode(
				tok.Location(),
				ast.NewPublicIdentifierNode(tok.Location(), tok.Value),
			))
			continue strContentLoop
		case token.STRING_INTERP_CONSTANT:
			tok := p.advance()
			strContent = append(strContent, ast.NewStringInterpolationNode(
				tok.Location(),
				ast.NewPublicConstantNode(tok.Location(), tok.Value),
			))
			continue strContentLoop
		case token.STRING_INSPECT_INTERP_BEG:
			beg := p.advance()
			expr := p.expressionWithoutModifier()
			end, _ := p.consume(token.STRING_INTERP_END)
			strContent = append(strContent, ast.NewStringInspectInterpolationNode(
				beg.Location().Join(end.Location()),
				expr,
			))
			continue strContentLoop
		case token.STRING_INSPECT_INTERP_LOCAL:
			tok := p.advance()
			strContent = append(strContent, ast.NewStringInspectInterpolationNode(
				tok.Location(),
				ast.NewPublicIdentifierNode(tok.Location(), tok.Value),
			))
			continue strContentLoop
		case token.STRING_INSPECT_INTERP_CONSTANT:
			tok := p.advance()
			strContent = append(strContent, ast.NewStringInspectInterpolationNode(
				tok.Location(),
				ast.NewPublicConstantNode(tok.Location(), tok.Value),
			))
			continue strContentLoop
		}

		tok, ok := p.consume(token.STRING_END)
		quoteEnd = tok
		if tok.Type == token.END_OF_FILE {
			break strContentLoop
		}
		if !ok {
			strContent = append(strContent, ast.NewInvalidNode(
				tok.Location(),
				tok,
			))
			continue strContentLoop
		}
		break strContentLoop
	}
	if len(strContent) == 0 {
		return ast.NewDoubleQuotedStringLiteralNode(
			quoteBeg.Location().Join(quoteEnd.Location()),
			"",
		)
	}
	strVal, ok := strContent[0].(*ast.StringLiteralContentSectionNode)
	if len(strContent) == 1 && ok {
		return ast.NewDoubleQuotedStringLiteralNode(
			quoteBeg.Location().Join(quoteEnd.Location()),
			strVal.Value,
		)
	}

	location := quoteBeg.Location().Join(quoteEnd.Location())
	if !withInterpolation {
		p.errorMessageLocation("cannot interpolate strings in this context", location)
	}

	return ast.NewInterpolatedStringLiteralNode(
		location,
		strContent,
	)
}

// closureAfterArrow = "->" (expressionWithoutModifier | SEPARATOR [statements] "end" | "{" [statements] "}")
func (p *Parser) closureAfterArrow(firstSpan *position.Location, params []ast.ParameterNode, returnType ast.TypeNode, throwType ast.TypeNode) ast.ExpressionNode {
	var location *position.Location
	arrowTok, ok := p.consume(token.THIN_ARROW)
	if !ok {
		return ast.NewInvalidNode(
			arrowTok.Location(),
			arrowTok,
		)
	}
	if firstSpan == nil {
		firstSpan = arrowTok.Location()
	}

	// Body with curly braces
	if p.match(token.LBRACE) {
		p.swallowNewlines()
		body := p.statements(token.RBRACE)
		if tok, ok := p.consume(token.RBRACE); ok {
			location = firstSpan.Join(tok.Location())
		} else {
			location = firstSpan
		}
		return ast.NewClosureLiteralNode(
			location,
			params,
			returnType,
			throwType,
			body,
		)
	}

	lastLocation, body, multiline := p.statementBlock(token.END)
	if lastLocation != nil {
		location = firstSpan.Join(lastLocation)
	} else {
		location = firstSpan
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
			location = location.Join(endTok.Location())
		}
	}

	return ast.NewClosureLiteralNode(
		location,
		params,
		returnType,
		throwType,
		body,
	)
}

// closureExpression = (("|" formalParameterList "|") | "||") [: typeAnnotation] ["!" typeAnnotation] closureAfterArrow
func (p *Parser) closureExpression() ast.ExpressionNode {
	var params []ast.ParameterNode
	var firstSpan *position.Location
	var returnType ast.TypeNode
	var throwType ast.TypeNode

	if p.accept(token.OR) || p.accept(token.OR_OR) {
		if p.accept(token.OR) {
			firstSpan = p.advance().Location()
			if !p.accept(token.OR) {
				p.mode = withoutBitwiseOrMode
				params = p.formalParameterList(token.OR)
				p.mode = normalMode
			}
			if tok, ok := p.consume(token.OR); !ok {
				return ast.NewInvalidNode(
					tok.Location(),
					tok,
				)
			}
		} else {
			orOr, ok := p.consume(token.OR_OR)
			firstSpan = orOr.Location()
			if !ok {
				return ast.NewInvalidNode(
					orOr.Location(),
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
	}

	return p.closureAfterArrow(firstSpan, params, returnType, throwType)
}

// identifierOrFunction = identifier | identifier closureAfterArrow
func (p *Parser) identifierOrFunction() ast.ExpressionNode {
	if p.secondLookahead.Type == token.THIN_ARROW {
		ident := p.advance()
		return p.closureAfterArrow(
			ident.Location(),
			[]ast.ParameterNode{
				ast.NewFormalParameterNode(
					ident.Location(),
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

// instanceVariable = INSTANCE_VARIABLE
func (p *Parser) instanceVariable() ast.ExpressionNode {
	token, ok := p.consume(token.INSTANCE_VARIABLE)
	if !ok {
		return ast.NewInvalidNode(token.Location(), token)
	}

	return ast.NewInstanceVariableNode(
		token.Location(),
		token.Value,
	)
}

// docComment = DOC_COMMENT declarationExpression
func (p *Parser) docComment(allowed bool) ast.ExpressionNode {
	docComment := p.advance()
	var nested bool
	if p.lookahead.Type == token.DOC_COMMENT {
		nested = true
		p.errorMessage("doc comments cannot document one another")
	}
	p.swallowNewlines()
	expr := p.declarationExpression(allowed)

	docCommentableExpr, ok := expr.(ast.DocCommentableNode)
	if !ok && !nested {
		p.errorMessageLocation("doc comments cannot be attached to this expression", expr.Location())
	}
	if ok {
		docCommentableExpr.SetDocComment(docComment.Value)
	}
	return expr
}

func (p *Parser) publicIdentifier() *ast.PublicIdentifierNode {
	ident, ok := p.matchOk(token.PUBLIC_IDENTIFIER)
	if !ok {
		panic(fmt.Sprintf("invalid public identifier token: %#v", ident))
	}
	return ast.NewPublicIdentifierNode(
		ident.Location(),
		ident.Value,
	)
}

func (p *Parser) privateIdentifier() *ast.PrivateIdentifierNode {
	ident, ok := p.matchOk(token.PRIVATE_IDENTIFIER)
	if !ok {
		panic(fmt.Sprintf("invalid private identifier token: %#v", ident))
	}
	return ast.NewPrivateIdentifierNode(
		ident.Location(),
		ident.Value,
	)
}

// identifier = PUBLIC_IDENTIFIER | PRIVATE_IDENTIFIER | unquote
func (p *Parser) identifier() ast.IdentifierNode {
	if p.accept(token.PUBLIC_IDENTIFIER) {
		return p.publicIdentifier()
	}
	if p.accept(token.PRIVATE_IDENTIFIER) {
		return p.privateIdentifier()
	}
	if p.accept(token.UNQUOTE) {
		return p.unquoteIdentifier()
	}

	p.errorExpected("an identifier")
	errTok := p.advance()
	return ast.NewInvalidNode(errTok.Location(), errTok)
}
