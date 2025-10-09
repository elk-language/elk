// Package ast defines types
// used by the Elk parser.
//
// All the nodes of the Abstract Syntax Tree
// constructed by the Elk parser are defined in this package.
package ast

import (
	"fmt"
	"iter"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

type static uint8

const (
	staticUnset static = iota
	staticFalse
	staticTrue
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
		if !isExpressionStatic(element) {
			return false
		}
	}

	return true
}

func isExpressionStatic(expr ExpressionNode) bool {
	return expr == nil || expr.IsStatic()
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

// Turn a type to a statement
func TypeToStatement(typ TypeNode) StatementNode {
	return NewTypeStatementNode(typ.Location(), typ)
}

// Turn a type to a collection of statements.
func TypeToStatements(typ TypeNode) []StatementNode {
	return []StatementNode{TypeToStatement(typ)}
}

// Turn a pattern to a statement
func PatternToStatement(pattern PatternNode) StatementNode {
	return NewPatternStatementNode(pattern.Location(), pattern)
}

// Turn a pattern to a collection of statements.
func PatternToStatements(pattern PatternNode) []StatementNode {
	return []StatementNode{PatternToStatement(pattern)}
}

// Turn an expression to a statement
func ExpressionToStatement(expr ExpressionNode) StatementNode {
	return NewExpressionStatementNode(expr.Location(), expr)
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
		*NumericForExpressionNode, *TypeExpressionNode, *ClosureLiteralNode,
		*LabeledExpressionNode, *QuoteExpressionNode:
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
	case *LabeledExpressionNode:
		return 5
	case *ModifierNode, *ModifierForInNode, *ModifierIfElseNode:
		return 10
	case *ReturnExpressionNode, *BreakExpressionNode,
		*ContinueExpressionNode, *YieldExpressionNode, *ThrowExpressionNode,
		*MustExpressionNode, *TryExpressionNode, *TypeofExpressionNode,
		*LoopExpressionNode, *GoExpressionNode, *DoExpressionNode,
		*IfExpressionNode, *UnlessExpressionNode, *WhileExpressionNode,
		*UntilExpressionNode, *ForInExpressionNode, *NumericForExpressionNode,
		*TypeExpressionNode, *ClosureLiteralNode, *ConstantDeclarationNode,
		*DoubleSplatExpressionNode, *SplatExpressionNode, *QuoteExpressionNode:
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
		*GenericMethodCallNode, *MethodCallNode, *AwaitExpressionNode,
		*MacroCallNode, *ReceiverlessMacroCallNode:
		return 210
	case *ConstructorCallNode, *GenericConstructorCallNode:
		return 220
	case *ConstantLookupNode, *MethodLookupNode, *InstanceMethodLookupNode:
		return 230
	}

	return 255
}

func TypeAssociativity(expr TypeNode) Associativity {
	switch expr.(type) {
	case *BinaryTypeNode, *NilableTypeNode,
		*IntersectionTypeNode, *UnionTypeNode:
		return LEFT_ASSOCIATIVE
	case *NotTypeNode, *SingletonTypeNode, *InstanceOfTypeNode,
		*UnaryTypeNode, *CallableTypeNode:
		return RIGHT_ASSOCIATIVE
	}

	return NON_ASSOCIATIVE
}

func TypePrecedence(expr TypeNode) uint8 {
	switch e := expr.(type) {
	case *CallableTypeNode:
		return 10
	case *UnionTypeNode:
		return 20
	case *IntersectionTypeNode:
		return 30
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
	case *BinaryPatternNode:
		return LEFT_ASSOCIATIVE
	case *UnaryExpressionNode:
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

// Value used to decide what whether to skip the children of the node,
// break the traversal or continue in the AST Traverse method.
// The zero value continues the traversal.
type TraverseOption uint8

const (
	TraverseContinue TraverseOption = iota
	TraverseSkip
	TraverseBreak
)

// Every node type implements this interface.
type Node interface {
	position.LocationInterface
	value.Reference
	IsStatic() bool // Value is known at compile-time
	// Return the static type of the value represented
	// by the AST node
	Type(*types.GlobalEnvironment) types.Type
	SetType(types.Type)
	SkipTypechecking() bool
	Equal(value.Value) bool
	String() string
	// Return the type of the AST node object
	// for use in macros
	MacroType(*types.GlobalEnvironment) types.Type
	splice(loc *position.Location, args *[]Node, unquote bool) Node
	traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption
}

// Create a deep copy of the AST node
func DeepCopy(node Node) Node {
	return Splice(node, nil, nil)
}

// Create a copy of AST replacing consecutive unquote nodes with the given arguments
func Splice(node Node, loc *position.Location, args *[]Node) Node {
	return node.splice(loc, args, false)
}

func noopTraverse(node, parent Node) TraverseOption { return TraverseContinue }

func Traverse(node Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) {
	if enter == nil {
		enter = noopTraverse
	}
	if leave == nil {
		leave = noopTraverse
	}
	node.traverse(nil, enter, leave)
}

func Iter(node Node) iter.Seq[Node] {
	return func(yield func(Node) bool) {
		Traverse(
			node,
			func(node, parent Node) TraverseOption {
				if !yield(node) {
					return TraverseBreak
				}
				return TraverseContinue
			},
			nil,
		)
	}
}

func NewNodeIterator(node Node) *value.ArrayTupleIterator {
	var tuple value.ArrayTuple
	for n := range Iter(node) {
		tuple = append(tuple, value.Ref(n))
	}

	return value.NewArrayTupleIterator(&tuple)
}

func SpliceSlice[N Node](slice []N, loc *position.Location, args *[]Node, unquote bool) []N {
	result := make([]N, len(slice))
	for i, n := range slice {
		result[i] = n.splice(loc, args, unquote).(N)
	}

	return result
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
	loc *position.Location
	typ types.Type
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
	return t.loc.Span
}

func (t *TypedNodeBase) SetSpan(span *position.Span) {
	t.loc.Span = span
}

func (t *TypedNodeBase) Location() *position.Location {
	return t.loc
}

func (t *TypedNodeBase) SetLocation(loc *position.Location) {
	t.loc = loc
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

func (t *TypedNodeBase) InstanceVariables() *value.InstanceVariables {
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

// Base AST node.
type NodeBase struct {
	loc *position.Location
}

func (*NodeBase) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Void{}
}

func (*NodeBase) SetType(types.Type) {}

func (t *NodeBase) SkipTypechecking() bool {
	return false
}

func (t *NodeBase) Span() *position.Span {
	return t.loc.Span
}

func (t *NodeBase) SetSpan(span *position.Span) {
	t.loc.Span = span
}

func (n *NodeBase) Location() *position.Location {
	return n.loc
}

func (n *NodeBase) SetLocation(loc *position.Location) {
	n.loc = loc
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

func (n *NodeBase) InstanceVariables() *value.InstanceVariables {
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
	switch node := node.(type) {
	case *PrivateIdentifierNode, *PublicIdentifierNode:
		return true
	case *UnquoteNode:
		switch node.Kind {
		case UNQUOTE_IDENTIFIER_KIND:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// Check whether the node can be used as a left value
// in an assignment expression.
func IsValidAssignmentTarget(node Node) bool {
	switch node.(type) {
	case *PrivateIdentifierNode, *PublicIdentifierNode,
		*AttributeAccessNode, *PublicInstanceVariableNode, *SubscriptExpressionNode:
		return true
	default:
		return false
	}
}

type StringOrSymbolLiteralNode interface {
	Node
	LiteralPatternNode
	stringOrSymbolLiteralNode()
}

func (*InvalidNode) stringOrSymbolLiteralNode()                   {}
func (*InterpolatedSymbolLiteralNode) stringOrSymbolLiteralNode() {}
func (*SimpleSymbolLiteralNode) stringOrSymbolLiteralNode()       {}
func (*DoubleQuotedStringLiteralNode) stringOrSymbolLiteralNode() {}
func (*RawStringLiteralNode) stringOrSymbolLiteralNode()          {}
func (*InterpolatedStringLiteralNode) stringOrSymbolLiteralNode() {}
