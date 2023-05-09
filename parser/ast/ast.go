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
	lexer.Positioner
}

// Check whether the token can be used as a left value
// in a variable/constant declaration.
func IsValidDeclarationTarget(node Node) bool {
	switch node.(type) {
	case *PrivateConstantNode, *ConstantNode, *PrivateIdentifierNode, *IdentifierNode:
		return true
	default:
		return false
	}
}

// Check whether the token can be used as a left value
// in an assignment expression.
func IsValidAssignmentTarget(node Node) bool {
	switch node.(type) {
	case *PrivateIdentifierNode, *IdentifierNode:
		return true
	default:
		return false
	}
}

// Check whether the node is a constant.
func IsConstant(node Node) bool {
	switch node.(type) {
	case *PrivateConstantNode, *ConstantNode:
		return true
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
func (*EmptyStatementNode) statementNode()      {}

// All expression nodes implement the Expr interface.
type ExpressionNode interface {
	Node
	expressionNode()
}

func (*InvalidNode) expressionNode()              {}
func (*ModifierNode) expressionNode()             {}
func (*ModifierIfElseNode) expressionNode()       {}
func (*AssignmentExpressionNode) expressionNode() {}
func (*BinaryExpressionNode) expressionNode()     {}
func (*LogicalExpressionNode) expressionNode()    {}
func (*UnaryExpressionNode) expressionNode()      {}
func (*TrueLiteralNode) expressionNode()          {}
func (*FalseLiteralNode) expressionNode()         {}
func (*NilLiteralNode) expressionNode()           {}
func (*RawStringLiteralNode) expressionNode()     {}
func (*IntLiteralNode) expressionNode()           {}
func (*FloatLiteralNode) expressionNode()         {}
func (*StringLiteralNode) expressionNode()        {}
func (*IdentifierNode) expressionNode()           {}
func (*PrivateIdentifierNode) expressionNode()    {}
func (*ConstantNode) expressionNode()             {}
func (*PrivateConstantNode) expressionNode()      {}
func (*SelfLiteralNode) expressionNode()          {}
func (*IfExpressionNode) expressionNode()         {}
func (*UnlessExpressionNode) expressionNode()     {}
func (*WhileExpressionNode) expressionNode()      {}
func (*UntilExpressionNode) expressionNode()      {}
func (*LoopExpressionNode) expressionNode()       {}
func (*BreakExpressionNode) expressionNode()      {}
func (*ReturnExpressionNode) expressionNode()     {}
func (*ContinueExpressionNode) expressionNode()   {}
func (*ThrowExpressionNode) expressionNode()      {}

// Nodes that implement this interface can appear
// inside of a String literal.
type StringLiteralContentNode interface {
	Node
	stringLiteralContentNode()
}

func (*InvalidNode) stringLiteralContentNode()                     {}
func (*StringInterpolationNode) stringLiteralContentNode()         {}
func (*StringLiteralContentSectionNode) stringLiteralContentNode() {}

// Represents a single Elk program (usually a single file).
type ProgramNode struct {
	*lexer.Position
	Body []StatementNode
}

// Represents an empty statement.
type EmptyStatementNode struct {
	*lexer.Position
}

// Expression optionally terminated with a newline or a semicolon.
type ExpressionStatementNode struct {
	*lexer.Position
	Expression ExpressionNode
}

// Assignment with the specified operator.
type AssignmentExpressionNode struct {
	*lexer.Position
	Left  ExpressionNode // left hand side
	Op    *lexer.Token   // operator
	Right ExpressionNode // right hand side
}

// Expression of an operator with two operands.
type BinaryExpressionNode struct {
	*lexer.Position
	Left  ExpressionNode // left hand side
	Op    *lexer.Token   // operator
	Right ExpressionNode // right hand side
}

// Expression of a logical operator with two operands.
type LogicalExpressionNode struct {
	*lexer.Position
	Left  ExpressionNode // left hand side
	Op    *lexer.Token   // operator
	Right ExpressionNode // right hand side
}

// Expression of an operator with one operand.
type UnaryExpressionNode struct {
	*lexer.Position
	Op    *lexer.Token   // operator
	Right ExpressionNode // right hand side
}

type TrueLiteralNode struct {
	*lexer.Position
}

type FalseLiteralNode struct {
	*lexer.Position
}

type SelfLiteralNode struct {
	*lexer.Position
}

type NilLiteralNode struct {
	*lexer.Position
}

type RawStringLiteralNode struct {
	*lexer.Position
	Value string // value of the string literal
}

type IntLiteralNode struct {
	*lexer.Position
	Token *lexer.Token
}

type FloatLiteralNode struct {
	*lexer.Position
	Value string
}

// Represents a syntax error.
type InvalidNode struct {
	*lexer.Position
	Token *lexer.Token
}

// Represents a single section of characters of a string literal.
type StringLiteralContentSectionNode struct {
	*lexer.Position
	Value string
}

// Represents a single interpolated section of a string literal.
type StringInterpolationNode struct {
	*lexer.Position
	Expression ExpressionNode
}

// Represents a string literal.
type StringLiteralNode struct {
	*lexer.Position
	Content []StringLiteralContentNode
}

// Represents a public identifier.
type IdentifierNode struct {
	*lexer.Position
	Value string
}

// Represents a private identifier.
type PrivateIdentifierNode struct {
	*lexer.Position
	Value string
}

// Represents a public constant.
type ConstantNode struct {
	*lexer.Position
	Value string
}

// Represents a private constant.
type PrivateConstantNode struct {
	*lexer.Position
	Value string
}

// Represents an `if`, `unless`, `while` or `until` modifier expression.
type ModifierNode struct {
	*lexer.Position
	Left     ExpressionNode // left hand side
	Modifier *lexer.Token   // modifier token
	Right    ExpressionNode // right hand side
}

// Represents an `if .. else` modifier expression.
type ModifierIfElseNode struct {
	*lexer.Position
	ThenExpression ExpressionNode // then expression body
	Condition      ExpressionNode // if condition
	ElseExpression ExpressionNode // else expression body
}

// Represents an `if` expression.
type IfExpressionNode struct {
	*lexer.Position
	Condition ExpressionNode  // if condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

// Represents an `unless` expression.
type UnlessExpressionNode struct {
	*lexer.Position
	Condition ExpressionNode  // unless condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

// Represents a `while` expression.
type WhileExpressionNode struct {
	*lexer.Position
	Condition ExpressionNode  // while condition
	ThenBody  []StatementNode // then expression body
}

// Represents a `until` expression.
type UntilExpressionNode struct {
	*lexer.Position
	Condition ExpressionNode  // until condition
	ThenBody  []StatementNode // then expression body
}

// Represents a `loop` expression.
type LoopExpressionNode struct {
	*lexer.Position
	ThenBody []StatementNode // then expression body
}

// Represents a `break` expression.
type BreakExpressionNode struct {
	*lexer.Position
}

// Represents a `return` expression.
type ReturnExpressionNode struct {
	*lexer.Position
	Value ExpressionNode
}

// Represents a `continue` expression.
type ContinueExpressionNode struct {
	*lexer.Position
	Value ExpressionNode
}

// Represents a `throw` expression.
type ThrowExpressionNode struct {
	*lexer.Position
	Value ExpressionNode
}
