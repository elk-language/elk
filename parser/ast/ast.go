// Package ast defines types
// used by the Elk parser.
//
// All the nodes of the Abstract Syntax Tree
// constructed by the Elk parser are defined in this package.
package ast

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

const indentUnitString = "  "

var indentUnitBytes = []byte(indentUnitString)
var newlineBytes = []byte{'\n'}

func indentString(out io.Writer, str string, indentLevel int) {
	scanner := bufio.NewScanner(strings.NewReader(str))

	firstIteration := true
	for scanner.Scan() {
		if !firstIteration {
			out.Write(newlineBytes)
		} else {
			firstIteration = false
		}

		for range indentLevel {
			out.Write(indentUnitBytes)
		}
		line := scanner.Bytes()
		out.Write(line)
	}
}

func indentStringFromSecondLine(out io.Writer, str string, indentLevel int) {
	scanner := bufio.NewScanner(strings.NewReader(str))

	firstIteration := true
	for scanner.Scan() {
		if !firstIteration {
			out.Write(newlineBytes)
			for range indentLevel {
				out.Write(indentUnitBytes)
			}
		} else {
			firstIteration = false
		}

		line := scanner.Bytes()
		out.Write(line)
	}
}

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
	return value.NodeClass
}

func (t *TypedNodeBase) DirectClass() *value.Class {
	return value.NodeClass
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
	return value.NodeClass
}

func (t *TypedNodeBaseWithLoc) DirectClass() *value.Class {
	return value.NodeClass
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
	return value.NodeClass
}

func (n *NodeBase) DirectClass() *value.Class {
	return value.NodeClass
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

// Check whether the node can be used as a range pattern element.
func IsValidRangePatternElement(node Node) bool {
	switch node.(type) {
	case *TrueLiteralNode, *FalseLiteralNode, *NilLiteralNode, *CharLiteralNode,
		*RawCharLiteralNode, *RawStringLiteralNode, *DoubleQuotedStringLiteralNode,
		*InterpolatedStringLiteralNode, *SimpleSymbolLiteralNode, *InterpolatedSymbolLiteralNode,
		*FloatLiteralNode, *Float64LiteralNode, *Float32LiteralNode, *BigFloatLiteralNode,
		*IntLiteralNode, *Int64LiteralNode, *UInt64LiteralNode, *Int32LiteralNode, *UInt32LiteralNode,
		*Int16LiteralNode, *UInt16LiteralNode, *Int8LiteralNode, *UInt8LiteralNode,
		*PublicConstantNode, *PrivateConstantNode, *ConstantLookupNode, *UnaryExpressionNode:
		return true
	default:
		return false
	}
}

// Check whether the node is a constant.
func IsConstant(node Node) bool {
	switch node.(type) {
	case *PrivateConstantNode, *PublicConstantNode:
		return true
	default:
		return false
	}
}

// Check whether the node is a complex constant.
func IsComplexConstant(node Node) bool {
	switch node.(type) {
	case *PrivateConstantNode, *PublicConstantNode, *ConstantLookupNode:
		return true
	default:
		return false
	}
}

// Check whether the node is a valid right operand of the pipe operator `|>`.
func IsValidPipeExpressionTarget(node Node) bool {
	switch node.(type) {
	case *MethodCallNode, *GenericMethodCallNode, *ReceiverlessMethodCallNode,
		*GenericReceiverlessMethodCallNode, *AttributeAccessNode, *ConstructorCallNode, *CallNode:
		return true
	default:
		return false
	}
}

type StringTypeNode interface {
	Node
	TypeNode
	StringLiteralNode
}

type StringOrSymbolTypeNode interface {
	Node
	TypeNode
	StringOrSymbolLiteralNode
}

// All nodes that should be valid in pattern matching should
// implement this interface
type PatternNode interface {
	Node
	patternNode()
}

func (*InvalidNode) patternNode()                    {}
func (*AsPatternNode) patternNode()                  {}
func (*BinHashSetLiteralNode) patternNode()          {}
func (*BinArrayTupleLiteralNode) patternNode()       {}
func (*BinArrayListLiteralNode) patternNode()        {}
func (*HexHashSetLiteralNode) patternNode()          {}
func (*HexArrayTupleLiteralNode) patternNode()       {}
func (*HexArrayListLiteralNode) patternNode()        {}
func (*SymbolHashSetLiteralNode) patternNode()       {}
func (*SymbolArrayTupleLiteralNode) patternNode()    {}
func (*SymbolArrayListLiteralNode) patternNode()     {}
func (*WordHashSetLiteralNode) patternNode()         {}
func (*WordArrayTupleLiteralNode) patternNode()      {}
func (*WordArrayListLiteralNode) patternNode()       {}
func (*SymbolKeyValuePatternNode) patternNode()      {}
func (*KeyValuePatternNode) patternNode()            {}
func (*ObjectPatternNode) patternNode()              {}
func (*RecordPatternNode) patternNode()              {}
func (*MapPatternNode) patternNode()                 {}
func (*RestPatternNode) patternNode()                {}
func (*SetPatternNode) patternNode()                 {}
func (*ListPatternNode) patternNode()                {}
func (*TuplePatternNode) patternNode()               {}
func (*ConstantLookupNode) patternNode()             {}
func (*PublicConstantNode) patternNode()             {}
func (*PrivateConstantNode) patternNode()            {}
func (*GenericConstantNode) patternNode()            {}
func (*PublicIdentifierNode) patternNode()           {}
func (*PrivateIdentifierNode) patternNode()          {}
func (*RangeLiteralNode) patternNode()               {}
func (*BinaryPatternNode) patternNode()              {}
func (*UnaryExpressionNode) patternNode()            {}
func (*TrueLiteralNode) patternNode()                {}
func (*FalseLiteralNode) patternNode()               {}
func (*NilLiteralNode) patternNode()                 {}
func (*CharLiteralNode) patternNode()                {}
func (*RawCharLiteralNode) patternNode()             {}
func (*DoubleQuotedStringLiteralNode) patternNode()  {}
func (*InterpolatedStringLiteralNode) patternNode()  {}
func (*RawStringLiteralNode) patternNode()           {}
func (*SimpleSymbolLiteralNode) patternNode()        {}
func (*InterpolatedSymbolLiteralNode) patternNode()  {}
func (*IntLiteralNode) patternNode()                 {}
func (*Int64LiteralNode) patternNode()               {}
func (*UInt64LiteralNode) patternNode()              {}
func (*Int32LiteralNode) patternNode()               {}
func (*UInt32LiteralNode) patternNode()              {}
func (*Int16LiteralNode) patternNode()               {}
func (*UInt16LiteralNode) patternNode()              {}
func (*Int8LiteralNode) patternNode()                {}
func (*UInt8LiteralNode) patternNode()               {}
func (*FloatLiteralNode) patternNode()               {}
func (*Float32LiteralNode) patternNode()             {}
func (*Float64LiteralNode) patternNode()             {}
func (*BigFloatLiteralNode) patternNode()            {}
func (*UninterpolatedRegexLiteralNode) patternNode() {}
func (*InterpolatedRegexLiteralNode) patternNode()   {}

func anyPatternDeclaresVariables(patterns []PatternNode) bool {
	for _, pat := range patterns {
		if PatternDeclaresVariables(pat) {
			return true
		}
	}
	return false
}

func PatternDeclaresVariables(pattern PatternNode) bool {
	switch pat := pattern.(type) {
	case *PublicIdentifierNode, *PrivateIdentifierNode, *AsPatternNode:
		return true
	case *BinaryPatternNode:
		return PatternDeclaresVariables(pat.Left) ||
			PatternDeclaresVariables(pat.Right)
	case *ObjectPatternNode:
		return anyPatternDeclaresVariables(pat.Attributes)
	case *SymbolKeyValuePatternNode:
		return PatternDeclaresVariables(pat.Value)
	case *KeyValuePatternNode:
		return PatternDeclaresVariables(pat.Value)
	case *MapPatternNode:
		return anyPatternDeclaresVariables(pat.Elements)
	case *RecordPatternNode:
		return anyPatternDeclaresVariables(pat.Elements)
	case *ListPatternNode:
		return anyPatternDeclaresVariables(pat.Elements)
	case *TuplePatternNode:
		return anyPatternDeclaresVariables(pat.Elements)
	case *RestPatternNode:
		switch pat.Identifier.(type) {
		case *PrivateIdentifierNode, *PublicIdentifierNode:
			return true
		}
		return false
	default:
		return false
	}
}

type PatternExpressionNode interface {
	Node
	ExpressionNode
	PatternNode
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

// Represents a `catch` eg.
//
//	catch SomeError(message)
//		print("awesome!")
//	end
type CatchNode struct {
	NodeBase
	Pattern       PatternNode
	StackTraceVar IdentifierNode
	Body          []StatementNode // do expression body
}

func (*CatchNode) IsStatic() bool {
	return false
}

// Create a new `catch` node eg.
//
//	catch SomeError(message)
//		print("awesome!")
//	end
func NewCatchNode(span *position.Span, pattern PatternNode, stackTraceVar IdentifierNode, body []StatementNode) *CatchNode {
	return &CatchNode{
		NodeBase:      NodeBase{span: span},
		Pattern:       pattern,
		StackTraceVar: stackTraceVar,
		Body:          body,
	}
}

// Pattern with two operands eg. `> 10 && < 50`
type BinaryPatternNode struct {
	TypedNodeBase
	Op    *token.Token // operator
	Left  PatternNode  // left hand side
	Right PatternNode  // right hand side
}

func (*BinaryPatternNode) IsStatic() bool {
	return false
}

// Create a new binary pattern node eg. `> 10 && < 50`
func NewBinaryPatternNode(span *position.Span, op *token.Token, left, right PatternNode) *BinaryPatternNode {
	return &BinaryPatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Same as [NewBinaryPatternNode] but returns an interface
func NewBinaryPatternNodeI(span *position.Span, op *token.Token, left, right PatternNode) PatternNode {
	return NewBinaryPatternNode(span, op, left, right)
}

// Represents an as pattern eg. `> 5 && < 20 as foo`
type AsPatternNode struct {
	NodeBase
	Pattern PatternNode
	Name    IdentifierNode
}

func (*AsPatternNode) IsStatic() bool {
	return false
}

// Create an Object pattern node eg. `Foo(foo: 5, bar: a, c)`
func NewAsPatternNode(span *position.Span, pattern PatternNode, name IdentifierNode) *AsPatternNode {
	return &AsPatternNode{
		NodeBase: NodeBase{span: span},
		Pattern:  pattern,
		Name:     name,
	}
}

// Represents a symbol value pattern eg. `foo: bar`
type SymbolKeyValuePatternNode struct {
	NodeBase
	Key   string
	Value PatternNode
}

func (s *SymbolKeyValuePatternNode) IsStatic() bool {
	return false
}

// Create a symbol key value node eg. `foo: bar`
func NewSymbolKeyValuePatternNode(span *position.Span, key string, val PatternNode) *SymbolKeyValuePatternNode {
	return &SymbolKeyValuePatternNode{
		NodeBase: NodeBase{span: span},
		Key:      key,
		Value:    val,
	}
}

// Represents a key value pattern eg. `foo => bar`
type KeyValuePatternNode struct {
	NodeBase
	Key   PatternExpressionNode
	Value PatternNode
}

func (k *KeyValuePatternNode) IsStatic() bool {
	return false
}

// Create a key value pattern node eg. `foo => bar`
func NewKeyValuePatternNode(span *position.Span, key PatternExpressionNode, val PatternNode) *KeyValuePatternNode {
	return &KeyValuePatternNode{
		NodeBase: NodeBase{span: span},
		Key:      key,
		Value:    val,
	}
}

// Represents an Object pattern eg. `Foo(foo: 5, bar: a, c)`
type ObjectPatternNode struct {
	TypedNodeBase
	ObjectType ComplexConstantNode
	Attributes []PatternNode
}

func (m *ObjectPatternNode) IsStatic() bool {
	return false
}

// Create an Object pattern node eg. `Foo(foo: 5, bar: a, c)`
func NewObjectPatternNode(span *position.Span, objectType ComplexConstantNode, attrs []PatternNode) *ObjectPatternNode {
	return &ObjectPatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ObjectType:    objectType,
		Attributes:    attrs,
	}
}

// Represents a Record pattern eg. `%{ foo: 5, bar: a, 5 => >= 10 }`
type RecordPatternNode struct {
	NodeBase
	Elements []PatternNode
}

func (m *RecordPatternNode) IsStatic() bool {
	return false
}

// Create a Record pattern node eg. `%{ foo: 5, bar: a, 5 => >= 10 }`
func NewRecordPatternNode(span *position.Span, elements []PatternNode) *RecordPatternNode {
	return &RecordPatternNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewRecordPatternNode] but returns an interface
func NewRecordPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewRecordPatternNode(span, elements)
}

// Represents a Map pattern eg. `{ foo: 5, bar: a, 5 => >= 10 }`
type MapPatternNode struct {
	TypedNodeBase
	Elements []PatternNode
}

func (m *MapPatternNode) IsStatic() bool {
	return false
}

// Create a Map pattern node eg. `{ foo: 5, bar: a, 5 => >= 10 }`
func NewMapPatternNode(span *position.Span, elements []PatternNode) *MapPatternNode {
	return &MapPatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewMapPatternNode] but returns an interface
func NewMapPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewMapPatternNode(span, elements)
}

// Represents a Set pattern eg. `^[1, "foo"]`
type SetPatternNode struct {
	TypedNodeBase
	Elements []PatternNode
}

func (s *SetPatternNode) IsStatic() bool {
	return false
}

// Create a Set pattern node eg. `^[1, "foo"]`
func NewSetPatternNode(span *position.Span, elements []PatternNode) *SetPatternNode {
	return &SetPatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewSetPatternNode] but returns an interface
func NewSetPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewSetPatternNode(span, elements)
}

// Represents a rest element in a list pattern eg. `*a`
type RestPatternNode struct {
	NodeBase
	Identifier IdentifierNode
}

func (r *RestPatternNode) IsStatic() bool {
	return false
}

// Create a rest pattern node eg. `*a`
func NewRestPatternNode(span *position.Span, ident IdentifierNode) *RestPatternNode {
	return &RestPatternNode{
		NodeBase:   NodeBase{span: span},
		Identifier: ident,
	}
}

// Represents a Tuple pattern eg. `%[1, a, >= 10]`
type TuplePatternNode struct {
	TypedNodeBase
	Elements []PatternNode
}

func (l *TuplePatternNode) IsStatic() bool {
	return false
}

// Create a Tuple pattern node eg. `%[1, a, >= 10]`
func NewTuplePatternNode(span *position.Span, elements []PatternNode) *TuplePatternNode {
	return &TuplePatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewTuplePatternNode] but returns an interface
func NewTuplePatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewTuplePatternNode(span, elements)
}

// Represents a List pattern eg. `[1, a, >= 10]`
type ListPatternNode struct {
	TypedNodeBase
	Elements []PatternNode
}

func (l *ListPatternNode) IsStatic() bool {
	return false
}

// Create a List pattern node eg. `[1, a, >= 10]`
func NewListPatternNode(span *position.Span, elements []PatternNode) *ListPatternNode {
	return &ListPatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewListPatternNode] but returns an interface
func NewListPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewListPatternNode(span, elements)
}

// Represents a `case` node eg. `case 3 then println("eureka!")`
type CaseNode struct {
	NodeBase
	Pattern PatternNode
	Body    []StatementNode
}

func (*CaseNode) IsStatic() bool {
	return false
}

// Create a new `case` node
func NewCaseNode(span *position.Span, pattern PatternNode, body []StatementNode) *CaseNode {
	return &CaseNode{
		NodeBase: NodeBase{span: span},
		Pattern:  pattern,
		Body:     body,
	}
}

// Type of an operator with one operand eg. `-2`, `+3`
type UnaryTypeNode struct {
	TypedNodeBase
	Op       *token.Token // operator
	TypeNode TypeNode     // right hand side
}

func (u *UnaryTypeNode) IsStatic() bool {
	return false
}

// Create a new unary expression node.
func NewUnaryTypeNode(span *position.Span, op *token.Token, typeNode TypeNode) *UnaryTypeNode {
	return &UnaryTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		TypeNode:      typeNode,
	}
}
