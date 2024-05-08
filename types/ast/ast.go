// Package ast defines typed AST nodes
// used by Elk
package ast

import (
	"go/token"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
)

type Node interface {
	ast.Node
	typ() types.Type
}

// Return the type of the given node.
func TypeOf(node Node, globalEnv *types.GlobalEnvironment) types.Type {
	switch node.(type) {
	case *FalseLiteralNode:
		return globalEnv.StdConst("False")
	case *TrueLiteralNode:
		return globalEnv.StdConst("True")
	case *NilLiteralNode:
		return globalEnv.StdConst("Nil")
	case *DoubleQuotedStringLiteralNode:
		return globalEnv.StdConst("String")
	case *IntLiteralNode:
		return globalEnv.StdConst("Int")
	case *Int64LiteralNode:
		return globalEnv.StdConst("Int64")
	case *Int32LiteralNode:
		return globalEnv.StdConst("Int32")
	case *Int16LiteralNode:
		return globalEnv.StdConst("Int16")
	case *Int8LiteralNode:
		return globalEnv.StdConst("Int8")
	case *UInt64LiteralNode:
		return globalEnv.StdConst("UInt64")
	case *UInt32LiteralNode:
		return globalEnv.StdConst("UInt32")
	case *UInt16LiteralNode:
		return globalEnv.StdConst("UInt16")
	case *UInt8LiteralNode:
		return globalEnv.StdConst("UInt8")
	case *FloatLiteralNode:
		return globalEnv.StdConst("Float")
	case *Float64LiteralNode:
		return globalEnv.StdConst("Float64")
	case *Float32LiteralNode:
		return globalEnv.StdConst("Float32")
	case *BigFloatLiteralNode:
		return globalEnv.StdConst("BigFloat")
	}
	return node.typ()
}

// Base struct of every AST node.
type NodeBase struct {
	span *position.Span
}

func (n *NodeBase) Span() *position.Span {
	return n.span
}

func (n *NodeBase) SetSpan(span *position.Span) {
	n.span = span
}

func (*NodeBase) typ() types.Type {
	return types.Void{}
}

// Represents a single statement, so for example
// a single valid "line" of Elk code.
// Usually its an expression optionally terminated with a newline ors semicolon.
type StatementNode interface {
	Node
	statementNode()
}

func (*InvalidNode) statementNode()             {}
func (*ExpressionStatementNode) statementNode() {}

// All expression nodes implement this interface.
type ExpressionNode interface {
	Node
	expressionNode()
}

func (*InvalidNode) expressionNode()                   {}
func (*TrueLiteralNode) expressionNode()               {}
func (*FalseLiteralNode) expressionNode()              {}
func (*NilLiteralNode) expressionNode()                {}
func (*IntLiteralNode) expressionNode()                {}
func (*Int64LiteralNode) expressionNode()              {}
func (*Int32LiteralNode) expressionNode()              {}
func (*Int16LiteralNode) expressionNode()              {}
func (*Int8LiteralNode) expressionNode()               {}
func (*UInt64LiteralNode) expressionNode()             {}
func (*UInt32LiteralNode) expressionNode()             {}
func (*UInt16LiteralNode) expressionNode()             {}
func (*UInt8LiteralNode) expressionNode()              {}
func (*FloatLiteralNode) expressionNode()              {}
func (*Float64LiteralNode) expressionNode()            {}
func (*Float32LiteralNode) expressionNode()            {}
func (*BigFloatLiteralNode) expressionNode()           {}
func (*DoubleQuotedStringLiteralNode) expressionNode() {}
func (*VariableDeclarationNode) expressionNode()       {}

// All nodes that should be valid in type annotations should
// implement this interface
type TypeNode interface {
	Node
	typeNode()
}

func (*InvalidNode) typeNode() {}

// Represents a syntax error.
type InvalidNode struct {
	NodeBase
	Token *token.Token
}

func (*InvalidNode) IsStatic() bool {
	return false
}

func (*InvalidNode) IsOptional() bool {
	return false
}

// Create a new invalid node.
func NewInvalidNode(span *position.Span, tok *token.Token) *InvalidNode {
	return &InvalidNode{
		NodeBase: NodeBase{span: span},
		Token:    tok,
	}
}

// Represents a single Elk program (usually a single file).
type ProgramNode struct {
	NodeBase
	Body []StatementNode
}

func (*ProgramNode) IsStatic() bool {
	return false
}

// Create a new program node.
func NewProgramNode(span *position.Span, body []StatementNode) *ProgramNode {
	return &ProgramNode{
		NodeBase: NodeBase{span: span},
		Body:     body,
	}
}

// Expression optionally terminated with a newline or a semicolon.
type ExpressionStatementNode struct {
	NodeBase
	Expression ExpressionNode
}

func (e *ExpressionStatementNode) IsStatic() bool {
	return e.Expression.IsStatic()
}

// Create a new expression statement node eg. `5 * 2\n`
func NewExpressionStatementNode(span *position.Span, expr ExpressionNode) *ExpressionStatementNode {
	return &ExpressionStatementNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// `true` literal.
type TrueLiteralNode struct {
	NodeBase
}

func (*TrueLiteralNode) IsStatic() bool {
	return true
}

// Create a new `true` literal node.
func NewTrueLiteralNode(span *position.Span) *TrueLiteralNode {
	return &TrueLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `self` literal.
type FalseLiteralNode struct {
	NodeBase
}

func (*FalseLiteralNode) IsStatic() bool {
	return true
}

// Create a new `false` literal node.
func NewFalseLiteralNode(span *position.Span) *FalseLiteralNode {
	return &FalseLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `nil` literal.
type NilLiteralNode struct {
	NodeBase
}

func (*NilLiteralNode) IsStatic() bool {
	return true
}

// Create a new `nil` literal node.
func NewNilLiteralNode(span *position.Span) *NilLiteralNode {
	return &NilLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// Int literal eg. `5`, `125_355`, `0xff`
type IntLiteralNode struct {
	NodeBase
	Value string
}

func (*IntLiteralNode) IsStatic() bool {
	return true
}

// Create a new int literal node eg. `5`, `125_355`, `0xff`
func NewIntLiteralNode(span *position.Span, val string) *IntLiteralNode {
	return &IntLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Int64 literal eg. `5i64`, `125_355i64`, `0xffi64`
type Int64LiteralNode struct {
	NodeBase
	Value string
}

func (*Int64LiteralNode) IsStatic() bool {
	return true
}

// Create a new Int64 literal node eg. `5i64`, `125_355i64`, `0xffi64`
func NewInt64LiteralNode(span *position.Span, val string) *Int64LiteralNode {
	return &Int64LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// UInt64 literal eg. `5u64`, `125_355u64`, `0xffu64`
type UInt64LiteralNode struct {
	NodeBase
	Value string
}

func (*UInt64LiteralNode) IsStatic() bool {
	return true
}

// Create a new UInt64 literal node eg. `5u64`, `125_355u64`, `0xffu64`
func NewUInt64LiteralNode(span *position.Span, val string) *UInt64LiteralNode {
	return &UInt64LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Int32 literal eg. `5i32`, `1_20i32`, `0xffi32`
type Int32LiteralNode struct {
	NodeBase
	Value string
}

func (*Int32LiteralNode) IsStatic() bool {
	return true
}

// Create a new Int32 literal node eg. `5i32`, `1_20i32`, `0xffi32`
func NewInt32LiteralNode(span *position.Span, val string) *Int32LiteralNode {
	return &Int32LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// UInt32 literal eg. `5u32`, `1_20u32`, `0xffu32`
type UInt32LiteralNode struct {
	NodeBase
	Value string
}

func (*UInt32LiteralNode) IsStatic() bool {
	return true
}

// Create a new UInt32 literal node eg. `5u32`, `1_20u32`, `0xffu32`
func NewUInt32LiteralNode(span *position.Span, val string) *UInt32LiteralNode {
	return &UInt32LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Int16 literal eg. `5i16`, `1_20i16`, `0xffi16`
type Int16LiteralNode struct {
	NodeBase
	Value string
}

func (*Int16LiteralNode) IsStatic() bool {
	return true
}

// Create a new Int16 literal node eg. `5i16`, `1_20i16`, `0xffi16`
func NewInt16LiteralNode(span *position.Span, val string) *Int16LiteralNode {
	return &Int16LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// UInt16 literal eg. `5u16`, `1_20u16`, `0xffu16`
type UInt16LiteralNode struct {
	NodeBase
	Value string
}

func (*UInt16LiteralNode) IsStatic() bool {
	return true
}

// Create a new UInt16 literal node eg. `5u16`, `1_20u16`, `0xffu16`
func NewUInt16LiteralNode(span *position.Span, val string) *UInt16LiteralNode {
	return &UInt16LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Int8 literal eg. `5i8`, `1_20i8`, `0xffi8`
type Int8LiteralNode struct {
	NodeBase
	Value string
}

func (*Int8LiteralNode) IsStatic() bool {
	return true
}

// Create a new Int8 literal node eg. `5i8`, `1_20i8`, `0xffi8`
func NewInt8LiteralNode(span *position.Span, val string) *Int8LiteralNode {
	return &Int8LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// UInt8 literal eg. `5u8`, `1_20u8`, `0xffu8`
type UInt8LiteralNode struct {
	NodeBase
	Value string
}

func (*UInt8LiteralNode) IsStatic() bool {
	return true
}

// Create a new UInt8 literal node eg. `5u8`, `1_20u8`, `0xffu8`
func NewUInt8LiteralNode(span *position.Span, val string) *UInt8LiteralNode {
	return &UInt8LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Float literal eg. `5.2`, `.5`, `45e20`
type FloatLiteralNode struct {
	NodeBase
	Value string
}

func (*FloatLiteralNode) IsStatic() bool {
	return true
}

// Create a new float literal node eg. `5.2`, `.5`, `45e20`
func NewFloatLiteralNode(span *position.Span, val string) *FloatLiteralNode {
	return &FloatLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// BigFloat literal eg. `5.2bf`, `.5bf`, `45e20bf`
type BigFloatLiteralNode struct {
	NodeBase
	Value string
}

func (*BigFloatLiteralNode) IsStatic() bool {
	return true
}

// Create a new BigFloat literal node eg. `5.2bf`, `.5bf`, `45e20bf`
func NewBigFloatLiteralNode(span *position.Span, val string) *BigFloatLiteralNode {
	return &BigFloatLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Float64 literal eg. `5.2f64`, `.5f64`, `45e20f64`
type Float64LiteralNode struct {
	NodeBase
	Value string
}

func (*Float64LiteralNode) IsStatic() bool {
	return true
}

// Create a new Float64 literal node eg. `5.2f64`, `.5f64`, `45e20f64`
func NewFloat64LiteralNode(span *position.Span, val string) *Float64LiteralNode {
	return &Float64LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Float32 literal eg. `5.2f32`, `.5f32`, `45e20f32`
type Float32LiteralNode struct {
	NodeBase
	Value string
}

func (*Float32LiteralNode) IsStatic() bool {
	return true
}

// Create a new Float32 literal node eg. `5.2f32`, `.5f32`, `45e20f32`
func NewFloat32LiteralNode(span *position.Span, val string) *Float32LiteralNode {
	return &Float32LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a simple double quoted string literal eg. `"foo baz"`
type DoubleQuotedStringLiteralNode struct {
	NodeBase
	Value string
}

func (*DoubleQuotedStringLiteralNode) IsStatic() bool {
	return true
}

// Create a new double quoted string literal node eg. `"foo baz"`
func NewDoubleQuotedStringLiteralNode(span *position.Span, val string) *DoubleQuotedStringLiteralNode {
	return &DoubleQuotedStringLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a variable declaration eg. `var foo: String`
type VariableDeclarationNode struct {
	NodeBase
	Name        *token.Token   // name of the variable
	TypeNode    TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
	_typ        types.Type
}

func (*VariableDeclarationNode) IsStatic() bool {
	return false
}

func (v *VariableDeclarationNode) typ() types.Type {
	return v._typ
}

// Create a new variable declaration node eg. `var foo: String`
func NewVariableDeclarationNode(span *position.Span, name *token.Token, typeNode TypeNode, init ExpressionNode, typ types.Type) *VariableDeclarationNode {
	return &VariableDeclarationNode{
		NodeBase:    NodeBase{span: span},
		Name:        name,
		TypeNode:    typeNode,
		Initialiser: init,
		_typ:        typ,
	}
}
