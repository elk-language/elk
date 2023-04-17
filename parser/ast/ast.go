// Package ast defines types
// used by the Elk parser.
//
// All the nodes of the Abstract Syntax Tree
// constructed by the Elk parser are defined in this package.
package ast

import (
	"github.com/elk-language/elk/lexer"
)

// Every node type implements this interface.
type Node interface {
	Pos() lexer.Position
}

// Check whether the token can be used as a left value,
// as in the left side of an assignment etc.
func IsValidLeftValue(node Node) bool {
	switch node.(type) {
	// case IdentifierNode:
	// 	return true
	default:
		return false
	}
}

// Represents a single statement, so for example
// a single valid "line" of Elk code.
// Usually its an expression optionally terminated with a newline or a semicolon.
type StatementNode interface {
	Node
	statementNode()
}

func (*ExpressionStatementNode) statementNode() {}

// All expression nodes implement the Expr interface.
type ExpressionNode interface {
	Node
	expressionNode()
}

func (*AssignmentExpressionNode) expressionNode() {}
func (*BinaryExpressionNode) expressionNode()     {}
func (*UnaryExpressionNode) expressionNode()      {}
func (*TrueLiteralNode) expressionNode()          {}
func (*FalseLiteralNode) expressionNode()         {}
func (*NilLiteralNode) expressionNode()           {}
func (*RawStringLiteralNode) expressionNode()     {}
func (*IntLiteralNode) expressionNode()           {}
func (*FloatLiteralNode) expressionNode()         {}
func (*InvalidNode) expressionNode()              {}

// Represents a single Elk program (usually a single file).
type ProgramNode struct {
	lexer.Position
	Body []StatementNode
}

// Expression optionally terminated with a newline or a semicolon.
type ExpressionStatementNode struct {
	lexer.Position
	Expression ExpressionNode
}

// Assignment with the specified operator.
type AssignmentExpressionNode struct {
	lexer.Position
	Left  ExpressionNode // left hand side
	Op    *lexer.Token   // operator
	Right ExpressionNode // right hand side
}

// Expression of an operator with two operands.
type BinaryExpressionNode struct {
	lexer.Position
	Left  ExpressionNode // left hand side
	Op    *lexer.Token   // operator
	Right ExpressionNode // right hand side
}

// Expression of an operator with one operand.
type UnaryExpressionNode struct {
	lexer.Position
	Op    *lexer.Token   // operator
	Right ExpressionNode // right hand side
}

type TrueLiteralNode struct {
	lexer.Position
}

type FalseLiteralNode struct {
	lexer.Position
}

type NilLiteralNode struct {
	lexer.Position
}

type RawStringLiteralNode struct {
	lexer.Position
	Value string // value of the string literal
}

type IntLiteralNode struct {
	lexer.Position
	Token *lexer.Token
}

type FloatLiteralNode struct {
	lexer.Position
	Value string
}

// Represents a syntax error.
type InvalidNode struct {
	lexer.Position
	Token *lexer.Token
}
