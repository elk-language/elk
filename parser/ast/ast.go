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

// Represents a single statement, so for example
// a single valid "line" of Elk code.
// Usually its an expression optionally terminated with a newline or a semicolon.
type StatementNode interface {
	Node
	statementNode()
}

// All expression nodes implement the Expr interface.
type ExpressionNode interface {
	Node
	expressionNode()
}

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

func (*ExpressionStatementNode) statementNode() {}

type AssignmentExpressionNode struct {
	lexer.Position
	Lhs ExpressionNode // left hand side
	Op  *lexer.Token   // operator
	Rhs ExpressionNode // right hand side
}

func (*AssignmentExpressionNode) expressionNode() {}
