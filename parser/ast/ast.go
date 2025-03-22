// Package ast defines types
// used by the Elk parser.
//
// All the nodes of the Abstract Syntax Tree
// constructed by the Elk parser are defined in this package.
package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Checks whether all expressions in the given list are static.
func isExpressionSliceStatic(elements []ExpressionNode) bool {
	for _, element := range elements {
		if !element.IsStatic() {
			return false
		}
	}
	return true
}

// Checks whether all expressions in the given list are static.
func areExpressionsStatic(elements ...ExpressionNode) bool {
	for _, element := range elements {
		if element != nil && !element.IsStatic() {
			return false
		}
	}
	return true
}

// Checks whether all nodes in the given list are static.
func areNodesStatic(elements ...Node) bool {
	for _, element := range elements {
		if element != nil && !element.IsStatic() {
			return false
		}
	}
	return true
}

// Turn an expression to a statement
func ExpressionToStatement(expr ExpressionNode) StatementNode {
	return NewExpressionStatementNode(expr.Span(), expr)
}

// Turn an expression to a collection of statements.
func ExpressionToStatements(expr ExpressionNode) []StatementNode {
	return []StatementNode{ExpressionToStatement(expr)}
}

// Every node type implements this interface.
type Node interface {
	position.SpanInterface
	value.Reference
	IsStatic() bool // Value is known at compile-time
	Type(*types.GlobalEnvironment) types.Type
	SetType(types.Type)
	SkipTypechecking() bool
}

type DocCommentableNode interface {
	DocComment() string
	SetDocComment(string)
}

type DocCommentableNodeBase struct {
	comment string
}

func (d *DocCommentableNodeBase) DocComment() string {
	return d.comment
}

func (d *DocCommentableNodeBase) SetDocComment(comment string) {
	d.comment = comment
}

// Base typed AST node.
type TypedNodeBase struct {
	span *position.Span
	typ  types.Type
}

func (t *TypedNodeBase) Type(*types.GlobalEnvironment) types.Type {
	return t.typ
}

func (t *TypedNodeBase) SkipTypechecking() bool {
	return t.typ != nil
}

func (t *TypedNodeBase) SetType(typ types.Type) {
	t.typ = typ
}

func (t *TypedNodeBase) Span() *position.Span {
	return t.span
}

func (t *TypedNodeBase) SetSpan(span *position.Span) {
	t.span = span
}

func (t *TypedNodeBase) Class() *value.Class {
	return nil
}

func (t *TypedNodeBase) DirectClass() *value.Class {
	return nil
}

func (t *TypedNodeBase) SingletonClass() *value.Class {
	return nil
}

func (t *TypedNodeBase) InstanceVariables() value.SymbolMap {
	return nil
}

func (t *TypedNodeBase) Copy() value.Reference {
	return t
}

func (t *TypedNodeBase) Inspect() string {
	return fmt.Sprintf("Std::Node{&: %p}", t)
}

func (t *TypedNodeBase) Error() string {
	return t.Inspect()
}

// Base typed AST node.
type TypedNodeBaseWithLoc struct {
	loc *position.Location
	typ types.Type
}

func (t *TypedNodeBaseWithLoc) Type(*types.GlobalEnvironment) types.Type {
	return t.typ
}

func (t *TypedNodeBaseWithLoc) SkipTypechecking() bool {
	return t.typ != nil
}

func (t *TypedNodeBaseWithLoc) SetType(typ types.Type) {
	t.typ = typ
}

func (t *TypedNodeBaseWithLoc) Span() *position.Span {
	return &t.loc.Span
}

func (t *TypedNodeBaseWithLoc) SetSpan(span *position.Span) {
	t.loc.Span = *span
}

func (t *TypedNodeBaseWithLoc) Location() *position.Location {
	return t.loc
}

func (t *TypedNodeBaseWithLoc) SetLocation(loc *position.Location) {
	t.loc = loc
}

func (t *TypedNodeBaseWithLoc) Class() *value.Class {
	return nil
}

func (t *TypedNodeBaseWithLoc) DirectClass() *value.Class {
	return nil
}

func (t *TypedNodeBaseWithLoc) SingletonClass() *value.Class {
	return nil
}

func (t *TypedNodeBaseWithLoc) InstanceVariables() value.SymbolMap {
	return nil
}

func (t *TypedNodeBaseWithLoc) Copy() value.Reference {
	return t
}

func (t *TypedNodeBaseWithLoc) Inspect() string {
	return fmt.Sprintf("Std::Node{&: %p}", t)
}

func (t *TypedNodeBaseWithLoc) Error() string {
	return t.Inspect()
}

// Base AST node.
type NodeBase struct {
	span *position.Span
}

func (*NodeBase) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Void{}
}

func (*NodeBase) SetType(types.Type) {}

func (t *NodeBase) SkipTypechecking() bool {
	return false
}

func (n *NodeBase) Span() *position.Span {
	return n.span
}

func (n *NodeBase) SetSpan(span *position.Span) {
	n.span = span
}

func (n *NodeBase) Class() *value.Class {
	return nil
}

func (n *NodeBase) DirectClass() *value.Class {
	return nil
}

func (n *NodeBase) SingletonClass() *value.Class {
	return nil
}

func (n *NodeBase) InstanceVariables() value.SymbolMap {
	return nil
}

func (n *NodeBase) Copy() value.Reference {
	return n
}

func (n *NodeBase) Inspect() string {
	return fmt.Sprintf("Std::Node{&: %p}", n)
}

func (n *NodeBase) Error() string {
	return n.Inspect()
}

// Check whether the node can be used as a left value
// in a variable/constant declaration.
func IsValidDeclarationTarget(node Node) bool {
	switch node.(type) {
	case *PrivateIdentifierNode, *PublicIdentifierNode:
		return true
	default:
		return false
	}
}

// Check whether the node can be used as a left value
// in an assignment expression.
func IsValidAssignmentTarget(node Node) bool {
	switch node.(type) {
	case *PrivateIdentifierNode, *PublicIdentifierNode,
		*AttributeAccessNode, *InstanceVariableNode, *SubscriptExpressionNode:
		return true
	default:
		return false
	}
}

type StringOrSymbolLiteralNode interface {
	Node
	PatternExpressionNode
	stringOrSymbolLiteralNode()
}

func (*InvalidNode) stringOrSymbolLiteralNode()                   {}
func (*InterpolatedSymbolLiteralNode) stringOrSymbolLiteralNode() {}
func (*SimpleSymbolLiteralNode) stringOrSymbolLiteralNode()       {}
func (*DoubleQuotedStringLiteralNode) stringOrSymbolLiteralNode() {}
func (*RawStringLiteralNode) stringOrSymbolLiteralNode()          {}
func (*InterpolatedStringLiteralNode) stringOrSymbolLiteralNode() {}
