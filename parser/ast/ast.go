// Package ast defines types
// used by the Elk parser.
//
// All the nodes of the Abstract Syntax Tree
// constructed by the Elk parser are defined in this package.
package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
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

type Associativity uint8

const (
	NON_ASSOCIATIVE Associativity = iota
	LEFT_ASSOCIATIVE
	RIGHT_ASSOCIATIVE
)

func ExpressionAssociativity(expr ExpressionNode) Associativity {
	switch e := expr.(type) {
	case *ReturnExpressionNode, *BreakExpressionNode,
		*ContinueExpressionNode, *UnaryExpressionNode,
		*AwaitExpressionNode, *YieldExpressionNode,
		*ThrowExpressionNode, *MustExpressionNode, *TryExpressionNode,
		*TypeofExpressionNode, *LoopExpressionNode, *GoExpressionNode,
		*DoExpressionNode, *IfExpressionNode, *UnlessExpressionNode,
		*WhileExpressionNode, *UntilExpressionNode, *ForInExpressionNode,
		*NumericForExpressionNode, *TypeExpressionNode, *ClosureLiteralNode:
		return RIGHT_ASSOCIATIVE
	case *BinaryExpressionNode:
		switch e.Op.Type {
		case token.STAR_STAR:
			return RIGHT_ASSOCIATIVE
		default:
			return LEFT_ASSOCIATIVE
		}
	case *LogicalExpressionNode:
		return LEFT_ASSOCIATIVE
	}

	return NON_ASSOCIATIVE
}

func StatementPrecedence(stmt StatementNode) uint8 {
	switch stmt := stmt.(type) {
	case *ExpressionStatementNode:
		return ExpressionPrecedence(stmt.Expression)
	default:
		return 0
	}
}

func ExpressionPrecedence(expr ExpressionNode) uint8 {
	switch e := expr.(type) {
	case *ModifierNode, *ModifierForInNode, *ModifierIfElseNode:
		return 10
	case *ReturnExpressionNode, *BreakExpressionNode,
		*ContinueExpressionNode, *YieldExpressionNode, *ThrowExpressionNode,
		*MustExpressionNode, *TryExpressionNode, *TypeofExpressionNode,
		*LoopExpressionNode, *GoExpressionNode, *DoExpressionNode,
		*IfExpressionNode, *UnlessExpressionNode, *WhileExpressionNode,
		*UntilExpressionNode, *ForInExpressionNode, *NumericForExpressionNode,
		*TypeExpressionNode, *ClosureLiteralNode, *ConstantDeclarationNode,
		*DoubleSplatExpressionNode, *SplatExpressionNode:
		return 20
	case *AssignmentExpressionNode:
		return 30
	case *LogicalExpressionNode:
		switch e.Op.Type {
		case token.OR_OR, token.QUESTION_QUESTION, token.OR_BANG:
			return 40
		case token.AND_AND, token.AND_BANG:
			return 50
		}
	case *BinaryExpressionNode:
		switch e.Op.Type {
		case token.PIPE_OP:
			return 60
		case token.OR:
			return 70
		case token.XOR:
			return 80
		case token.AND:
			return 90
		case token.AND_TILDE:
			return 100
		case token.EQUAL_EQUAL, token.NOT_EQUAL, token.STRICT_EQUAL,
			token.STRICT_NOT_EQUAL, token.LAX_EQUAL, token.LAX_NOT_EQUAL:
			return 110
		case token.LESS, token.LESS_EQUAL, token.GREATER,
			token.GREATER_EQUAL, token.ISA_OP, token.REVERSE_ISA_OP,
			token.INSTANCE_OF_OP, token.REVERSE_INSTANCE_OF_OP, token.SPACESHIP_OP:
			return 120
		case token.LBITSHIFT, token.LTRIPLE_BITSHIFT, token.RBITSHIFT, token.RTRIPLE_BITSHIFT:
			return 130
		case token.PLUS, token.MINUS:
			return 140
		case token.STAR, token.SLASH, token.PERCENT:
			return 150
		case token.STAR_STAR:
			return 190
		}
	case *RangeLiteralNode:
		return 160
	case *AsExpressionNode:
		return 170
	case *UnaryExpressionNode:
		return 180
	case *PostfixExpressionNode:
		return 200
	case *GenericReceiverlessMethodCallNode,
		*ReceiverlessMethodCallNode, *NilSafeSubscriptExpressionNode,
		*SubscriptExpressionNode, *CallNode, *AttributeAccessNode,
		*GenericMethodCallNode, *MethodCallNode, *AwaitExpressionNode:
		return 210
	case *ConstructorCallNode, *GenericConstructorCallNode:
		return 220
	case *ConstantLookupNode:
		return 230
	}

	return 255
}

func TypeAssociativity(expr TypeNode) Associativity {
	switch expr.(type) {
	case *BinaryTypeNode, *NilableTypeNode:
		return LEFT_ASSOCIATIVE
	case *NotTypeNode, *SingletonTypeNode, *InstanceOfTypeNode,
		*UnaryTypeNode, *ClosureTypeNode:
		return RIGHT_ASSOCIATIVE
	}

	return NON_ASSOCIATIVE
}

func TypePrecedence(expr TypeNode) uint8 {
	switch e := expr.(type) {
	case *ClosureTypeNode:
		return 10
	case *BinaryTypeNode:
		switch e.Op.Type {
		case token.OR:
			return 20
		case token.AND:
			return 30
		case token.SLASH:
			return 40
		}
	case *NotTypeNode, *SingletonTypeNode, *InstanceOfTypeNode:
		return 50
	case *NilableTypeNode:
		return 60
	case *UnaryTypeNode:
		return 70
	}

	return 255
}

func PatternAssociativity(expr PatternNode) Associativity {
	switch expr.(type) {
	case *BinaryTypeNode, *NilableTypeNode:
		return LEFT_ASSOCIATIVE
	case *NotTypeNode, *SingletonTypeNode, *InstanceOfTypeNode,
		*UnaryTypeNode:
		return RIGHT_ASSOCIATIVE
	}

	return NON_ASSOCIATIVE
}

func PatternPrecedence(expr PatternNode) uint8 {
	switch e := expr.(type) {
	case *AsPatternNode:
		return 10
	case *BinaryPatternNode:
		switch e.Op.Type {
		case token.OR_OR:
			return 20
		case token.AND_AND:
			return 30
		}
	case *UnaryExpressionNode:
		switch e.Op.Type {
		case token.GREATER, token.GREATER_EQUAL,
			token.LESS, token.LESS_EQUAL,
			token.EQUAL_EQUAL, token.NOT_EQUAL,
			token.STRICT_EQUAL, token.STRICT_NOT_EQUAL,
			token.LAX_EQUAL, token.LAX_NOT_EQUAL:
			return 40
		case token.MINUS, token.PLUS:
			return 60
		}
	case *RangeLiteralNode:
		return 50
	}

	return 255
}

// Every node type implements this interface.
type Node interface {
	position.SpanInterface
	value.Reference
	IsStatic() bool // Value is known at compile-time
	Type(*types.GlobalEnvironment) types.Type
	SetType(types.Type)
	SkipTypechecking() bool
	Equal(value.Value) bool
	String() string
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
