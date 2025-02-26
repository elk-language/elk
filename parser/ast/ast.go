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
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
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

// Represents a single statement of a struct body
// optionally terminated with a newline or semicolon.
type StructBodyStatementNode interface {
	Node
	structBodyStatementNode()
}

func (*InvalidNode) structBodyStatementNode()            {}
func (*EmptyStatementNode) structBodyStatementNode()     {}
func (*ParameterStatementNode) structBodyStatementNode() {}

// All nodes that should be able to appear as
// elements of word collection literals should
// implement this interface.
type WordCollectionContentNode interface {
	Node
	ExpressionNode
	wordCollectionContentNode()
}

func (*InvalidNode) wordCollectionContentNode()          {}
func (*RawStringLiteralNode) wordCollectionContentNode() {}

// All nodes that should be able to appear as
// elements of symbol collection literals should
// implement this interface.
type SymbolCollectionContentNode interface {
	Node
	ExpressionNode
	symbolCollectionContentNode()
}

func (*InvalidNode) symbolCollectionContentNode()             {}
func (*SimpleSymbolLiteralNode) symbolCollectionContentNode() {}

// All nodes that should be able to appear as
// elements of Int collection literals should
// implement this interface.
type IntCollectionContentNode interface {
	Node
	ExpressionNode
	intCollectionContentNode()
}

func (*InvalidNode) intCollectionContentNode()    {}
func (*IntLiteralNode) intCollectionContentNode() {}

// All nodes that should be valid in type annotations should
// implement this interface
type TypeNode interface {
	Node
	typeNode()
}

func (*InvalidNode) typeNode()                   {}
func (*SelfLiteralNode) typeNode()               {}
func (*AnyTypeNode) typeNode()                   {}
func (*NeverTypeNode) typeNode()                 {}
func (*VoidTypeNode) typeNode()                  {}
func (*UnionTypeNode) typeNode()                 {}
func (*IntersectionTypeNode) typeNode()          {}
func (*BinaryTypeExpressionNode) typeNode()      {}
func (*NilableTypeNode) typeNode()               {}
func (*InstanceOfTypeNode) typeNode()            {}
func (*SingletonTypeNode) typeNode()             {}
func (*ClosureTypeNode) typeNode()               {}
func (*NotTypeNode) typeNode()                   {}
func (*UnaryTypeNode) typeNode()                 {}
func (*PublicConstantNode) typeNode()            {}
func (*PrivateConstantNode) typeNode()           {}
func (*ConstantLookupNode) typeNode()            {}
func (*GenericConstantNode) typeNode()           {}
func (*NilLiteralNode) typeNode()                {}
func (*BoolLiteralNode) typeNode()               {}
func (*TrueLiteralNode) typeNode()               {}
func (*FalseLiteralNode) typeNode()              {}
func (*CharLiteralNode) typeNode()               {}
func (*RawCharLiteralNode) typeNode()            {}
func (*RawStringLiteralNode) typeNode()          {}
func (*InterpolatedStringLiteralNode) typeNode() {}
func (*DoubleQuotedStringLiteralNode) typeNode() {}
func (*SimpleSymbolLiteralNode) typeNode()       {}
func (*InterpolatedSymbolLiteralNode) typeNode() {}
func (*IntLiteralNode) typeNode()                {}
func (*Int64LiteralNode) typeNode()              {}
func (*Int32LiteralNode) typeNode()              {}
func (*Int16LiteralNode) typeNode()              {}
func (*Int8LiteralNode) typeNode()               {}
func (*UInt64LiteralNode) typeNode()             {}
func (*UInt32LiteralNode) typeNode()             {}
func (*UInt16LiteralNode) typeNode()             {}
func (*UInt8LiteralNode) typeNode()              {}
func (*FloatLiteralNode) typeNode()              {}
func (*Float64LiteralNode) typeNode()            {}
func (*Float32LiteralNode) typeNode()            {}
func (*BigFloatLiteralNode) typeNode()           {}

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

// All nodes that represent regexes should
// implement this interface.
type RegexLiteralNode interface {
	Node
	PatternExpressionNode
	regexLiteralNode()
}

func (*InvalidNode) regexLiteralNode()                    {}
func (*UninterpolatedRegexLiteralNode) regexLiteralNode() {}
func (*InterpolatedRegexLiteralNode) regexLiteralNode()   {}

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

// All nodes that represent strings should
// implement this interface.
type StringLiteralNode interface {
	Node
	PatternExpressionNode
	StringOrSymbolLiteralNode
	stringLiteralNode()
}

func (*InvalidNode) stringLiteralNode()                   {}
func (*DoubleQuotedStringLiteralNode) stringLiteralNode() {}
func (*RawStringLiteralNode) stringLiteralNode()          {}
func (*InterpolatedStringLiteralNode) stringLiteralNode() {}

// All nodes that represent simple strings (without interpolation)
// should implement this interface.
type SimpleStringLiteralNode interface {
	Node
	ExpressionNode
	StringLiteralNode
	StringOrSymbolLiteralNode
	simpleStringLiteralNode()
}

func (*InvalidNode) simpleStringLiteralNode()                   {}
func (*DoubleQuotedStringLiteralNode) simpleStringLiteralNode() {}
func (*RawStringLiteralNode) simpleStringLiteralNode()          {}

// All nodes that should be valid in parameter declaration lists
// of methods or functions should implement this interface.
type ParameterNode interface {
	Node
	parameterNode()
	IsOptional() bool
}

func (*InvalidNode) parameterNode()            {}
func (*FormalParameterNode) parameterNode()    {}
func (*MethodParameterNode) parameterNode()    {}
func (*SignatureParameterNode) parameterNode() {}
func (*AttributeParameterNode) parameterNode() {}

// checks whether the given parameter is a positional rest parameter.
func IsPositionalRestParam(p ParameterNode) bool {
	switch param := p.(type) {
	case *MethodParameterNode:
		return param.Kind == PositionalRestParameterKind
	case *FormalParameterNode:
		return param.Kind == PositionalRestParameterKind
	case *SignatureParameterNode:
		return param.Kind == PositionalRestParameterKind
	default:
		return false
	}
}

// checks whether the given parameter is a named rest parameter.
func IsNamedRestParam(p ParameterNode) bool {
	switch param := p.(type) {
	case *MethodParameterNode:
		return param.Kind == NamedRestParameterKind
	case *FormalParameterNode:
		return param.Kind == NamedRestParameterKind
	case *SignatureParameterNode:
		return param.Kind == NamedRestParameterKind
	default:
		return false
	}
}

// Represents a type variable in generics like `class Foo[+V]; end`
type TypeParameterNode interface {
	Node
	typeVariableNode()
}

func (*InvalidNode) typeVariableNode()              {}
func (*VariantTypeParameterNode) typeVariableNode() {}

// All nodes that should be valid in constant lookups
// should implement this interface.
type ComplexConstantNode interface {
	Node
	TypeNode
	ExpressionNode
	PatternNode
	PatternExpressionNode
	UsingEntryNode
	complexConstantNode()
}

func (*InvalidNode) complexConstantNode()         {}
func (*PublicConstantNode) complexConstantNode()  {}
func (*PrivateConstantNode) complexConstantNode() {}
func (*ConstantLookupNode) complexConstantNode()  {}
func (*GenericConstantNode) complexConstantNode() {}
func (*NilLiteralNode) complexConstantNode()      {}

// Represents all nodes that are valid in using declarations
type UsingEntryNode interface {
	Node
	ExpressionNode
	usingEntryNode()
}

func (*InvalidNode) usingEntryNode()                  {}
func (*PublicConstantNode) usingEntryNode()           {}
func (*PrivateConstantNode) usingEntryNode()          {}
func (*ConstantLookupNode) usingEntryNode()           {}
func (*MethodLookupNode) usingEntryNode()             {}
func (*UsingAllEntryNode) usingEntryNode()            {}
func (*UsingEntryWithSubentriesNode) usingEntryNode() {}
func (*ConstantAsNode) usingEntryNode()               {}
func (*MethodLookupAsNode) usingEntryNode()           {}
func (*GenericConstantNode) usingEntryNode()          {}
func (*NilLiteralNode) usingEntryNode()               {}

type UsingSubentryNode interface {
	Node
	ExpressionNode
	usingSubentryNode()
}

func (*InvalidNode) usingSubentryNode()            {}
func (*PublicConstantNode) usingSubentryNode()     {}
func (*PublicConstantAsNode) usingSubentryNode()   {}
func (*PublicIdentifierNode) usingSubentryNode()   {}
func (*PublicIdentifierAsNode) usingSubentryNode() {}

// All nodes that should be valid constants
// should implement this interface.
type ConstantNode interface {
	Node
	TypeNode
	ExpressionNode
	UsingEntryNode
	ComplexConstantNode
	constantNode()
}

func (*InvalidNode) constantNode()         {}
func (*PublicConstantNode) constantNode()  {}
func (*PrivateConstantNode) constantNode() {}

// All nodes that should be valid identifiers
// should implement this interface.
type IdentifierNode interface {
	Node
	PatternExpressionNode
	identifierNode()
}

func (*InvalidNode) identifierNode()           {}
func (*PublicIdentifierNode) identifierNode()  {}
func (*PrivateIdentifierNode) identifierNode() {}

// Nodes that implement this interface can appear
// inside of a String literal.
type StringLiteralContentNode interface {
	Node
	stringLiteralContentNode()
}

func (*InvalidNode) stringLiteralContentNode()                     {}
func (*StringInspectInterpolationNode) stringLiteralContentNode()  {}
func (*StringInterpolationNode) stringLiteralContentNode()         {}
func (*StringLiteralContentSectionNode) stringLiteralContentNode() {}

// Nodes that implement this interface can appear
// inside of a Regex literal.
type RegexLiteralContentNode interface {
	Node
	regexLiteralContentNode()
}

func (*InvalidNode) regexLiteralContentNode()                    {}
func (*RegexInterpolationNode) regexLiteralContentNode()         {}
func (*RegexLiteralContentSectionNode) regexLiteralContentNode() {}

// Nodes that implement this interface represent
// named arguments in method calls.
type NamedArgumentNode interface {
	Node
	namedArgumentNode()
}

func (*InvalidNode) namedArgumentNode()               {}
func (*NamedCallArgumentNode) namedArgumentNode()     {}
func (*DoubleSplatExpressionNode) namedArgumentNode() {}

type ProgramState uint8

const (
	UNCHECKED ProgramState = iota
	CHECKING_NAMESPACES
	CHECKED_NAMESPACES

	CHECKING_METHODS
	CHECKED_METHODS

	CHECKING_EXPRESSIONS
	CHECKED_EXPRESSIONS
)

// Formal parameter optionally terminated with a newline or a semicolon.
type ParameterStatementNode struct {
	NodeBase
	Parameter ParameterNode
}

func (*ParameterStatementNode) IsStatic() bool {
	return false
}

// Create a new formal parameter statement node eg. `foo: Bar\n`
func NewParameterStatementNode(span *position.Span, param ParameterNode) *ParameterStatementNode {
	return &ParameterStatementNode{
		NodeBase:  NodeBase{span: span},
		Parameter: param,
	}
}

// Same as [NewParameterStatementNode] but returns an interface
func NewParameterStatementNodeI(span *position.Span, param ParameterNode) StructBodyStatementNode {
	return &ParameterStatementNode{
		NodeBase:  NodeBase{span: span},
		Parameter: param,
	}
}

// `bool` literal.
type BoolLiteralNode struct {
	NodeBase
}

func (*BoolLiteralNode) IsStatic() bool {
	return true
}

func (*BoolLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Bool{}
}

// Create a new `bool` literal node.
func NewBoolLiteralNode(span *position.Span) *BoolLiteralNode {
	return &BoolLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `void` type.
type VoidTypeNode struct {
	NodeBase
}

func (*VoidTypeNode) IsStatic() bool {
	return true
}

// Create a new `void` type node.
func NewVoidTypeNode(span *position.Span) *VoidTypeNode {
	return &VoidTypeNode{
		NodeBase: NodeBase{span: span},
	}
}

// `never` type.
type NeverTypeNode struct {
	NodeBase
}

func (*NeverTypeNode) IsStatic() bool {
	return true
}

func (*NeverTypeNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Never{}
}

// Create a new `never` type node.
func NewNeverTypeNode(span *position.Span) *NeverTypeNode {
	return &NeverTypeNode{
		NodeBase: NodeBase{span: span},
	}
}

// `any` type.
type AnyTypeNode struct {
	NodeBase
}

func (*AnyTypeNode) IsStatic() bool {
	return true
}

func (*AnyTypeNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Any{}
}

// Create a new `any` type node.
func NewAnyTypeNode(span *position.Span) *AnyTypeNode {
	return &AnyTypeNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a public identifier eg. `foo`.
type PublicIdentifierNode struct {
	TypedNodeBase
	Value string
}

func (*PublicIdentifierNode) IsStatic() bool {
	return false
}

// Create a new public identifier node eg. `foo`.
func NewPublicIdentifierNode(span *position.Span, val string) *PublicIdentifierNode {
	return &PublicIdentifierNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a private identifier eg. `_foo`
type PrivateIdentifierNode struct {
	TypedNodeBase
	Value string
}

func (*PrivateIdentifierNode) IsStatic() bool {
	return false
}

// Create a new private identifier node eg. `_foo`.
func NewPrivateIdentifierNode(span *position.Span, val string) *PrivateIdentifierNode {
	return &PrivateIdentifierNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a public constant eg. `Foo`.
type PublicConstantNode struct {
	TypedNodeBase
	Value string
}

func (*PublicConstantNode) IsStatic() bool {
	return false
}

// Create a new public constant node eg. `Foo`.
func NewPublicConstantNode(span *position.Span, val string) *PublicConstantNode {
	return &PublicConstantNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a private constant eg. `_Foo`
type PrivateConstantNode struct {
	TypedNodeBase
	Value string
}

func (*PrivateConstantNode) IsStatic() bool {
	return false
}

// Create a new private constant node eg. `_Foo`.
func NewPrivateConstantNode(span *position.Span, val string) *PrivateConstantNode {
	return &PrivateConstantNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents an as type downcast eg. `foo as String`
type AsExpressionNode struct {
	TypedNodeBase
	Value       ExpressionNode
	RuntimeType ComplexConstantNode
}

func (*AsExpressionNode) IsStatic() bool {
	return false
}

// Create a new private constant node eg. `_Foo`.
func NewAsExpressionNode(span *position.Span, val ExpressionNode, runtimeType ComplexConstantNode) *AsExpressionNode {
	return &AsExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
		RuntimeType:   runtimeType,
	}
}

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

// Represents a `do` expression eg.
//
//	do
//		print("awesome!")
//	end
type DoExpressionNode struct {
	TypedNodeBase
	Body    []StatementNode // do expression body
	Catches []*CatchNode
	Finally []StatementNode
}

func (*DoExpressionNode) IsStatic() bool {
	return false
}

// Create a new `do` expression node eg.
//
//	do
//		print("awesome!")
//	end
func NewDoExpressionNode(span *position.Span, body []StatementNode, catches []*CatchNode, finally []StatementNode) *DoExpressionNode {
	return &DoExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
		Catches:       catches,
		Finally:       finally,
	}
}

// Represents an `extend where` block expression eg.
//
//	extend where T < Foo
//		def hello then println("awesome!")
//	end
type ExtendWhereBlockExpressionNode struct {
	TypedNodeBase
	Body  []StatementNode
	Where []TypeParameterNode
}

func (*ExtendWhereBlockExpressionNode) SkipTypechecking() bool {
	return false
}

func (*ExtendWhereBlockExpressionNode) IsStatic() bool {
	return false
}

// Create a new `singleton` block expression node eg.
//
//	singleton
//		def hello then println("awesome!")
//	end
func NewExtendWhereBlockExpressionNode(span *position.Span, body []StatementNode, where []TypeParameterNode) *ExtendWhereBlockExpressionNode {
	return &ExtendWhereBlockExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
		Where:         where,
	}
}

// Represents a `singleton` block expression eg.
//
//	singleton
//		def hello then println("awesome!")
//	end
type SingletonBlockExpressionNode struct {
	TypedNodeBase
	Body     []StatementNode // do expression body
	Bytecode *vm.BytecodeFunction
}

func (*SingletonBlockExpressionNode) SkipTypechecking() bool {
	return false
}

func (*SingletonBlockExpressionNode) IsStatic() bool {
	return false
}

// Create a new `singleton` block expression node eg.
//
//	singleton
//		def hello then println("awesome!")
//	end
func NewSingletonBlockExpressionNode(span *position.Span, body []StatementNode) *SingletonBlockExpressionNode {
	return &SingletonBlockExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
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

// Represents a `switch` expression eg.
//
//	switch a
//	case 3
//	  println("eureka!")
//	case nil
//	  println("boo")
//	else
//	  println("nothing")
//	end
type SwitchExpressionNode struct {
	TypedNodeBase
	Value    ExpressionNode
	Cases    []*CaseNode
	ElseBody []StatementNode
}

func (*SwitchExpressionNode) IsStatic() bool {
	return false
}

// Create a new `switch` expression node
func NewSwitchExpressionNode(span *position.Span, val ExpressionNode, cases []*CaseNode, els []StatementNode) *SwitchExpressionNode {
	return &SwitchExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
		Cases:         cases,
		ElseBody:      els,
	}
}

// Represents an `if` expression eg. `if foo then println("bar")`
type IfExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // if condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

func (*IfExpressionNode) IsStatic() bool {
	return false
}

// Create a new `if` expression node eg. `if foo then println("bar")`
func NewIfExpressionNode(span *position.Span, cond ExpressionNode, then, els []StatementNode) *IfExpressionNode {
	return &IfExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ThenBody:      then,
		Condition:     cond,
		ElseBody:      els,
	}
}

// Represents an `unless` expression eg. `unless foo then println("bar")`
type UnlessExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // unless condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

func (*UnlessExpressionNode) IsStatic() bool {
	return false
}

// Create a new `unless` expression node eg. `unless foo then println("bar")`
func NewUnlessExpressionNode(span *position.Span, cond ExpressionNode, then, els []StatementNode) *UnlessExpressionNode {
	return &UnlessExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ThenBody:      then,
		Condition:     cond,
		ElseBody:      els,
	}
}

// Represents a `while` expression eg. `while i < 5 then i += 5`
type WhileExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // while condition
	ThenBody  []StatementNode // then expression body
}

func (*WhileExpressionNode) IsStatic() bool {
	return false
}

// Create a new `while` expression node eg. `while i < 5 then i += 5`
func NewWhileExpressionNode(span *position.Span, cond ExpressionNode, then []StatementNode) *WhileExpressionNode {
	return &WhileExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Condition:     cond,
		ThenBody:      then,
	}
}

// Represents a `until` expression eg. `until i >= 5 then i += 5`
type UntilExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // until condition
	ThenBody  []StatementNode // then expression body
}

func (*UntilExpressionNode) IsStatic() bool {
	return false
}

// Create a new `until` expression node eg. `until i >= 5 then i += 5`
func NewUntilExpressionNode(span *position.Span, cond ExpressionNode, then []StatementNode) *UntilExpressionNode {
	return &UntilExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Condition:     cond,
		ThenBody:      then,
	}
}

// Represents a `loop` expression.
type LoopExpressionNode struct {
	TypedNodeBase
	ThenBody []StatementNode // then expression body
}

func (*LoopExpressionNode) IsStatic() bool {
	return false
}

// Create a new `loop` expression node eg. `loop println('elk is awesome')`
func NewLoopExpressionNode(span *position.Span, then []StatementNode) *LoopExpressionNode {
	return &LoopExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ThenBody:      then,
	}
}

// Represents a numeric `for` expression eg. `fornum i := 0; i < 10; i += 1 then println(i)`
type NumericForExpressionNode struct {
	TypedNodeBase
	Initialiser ExpressionNode  // i := 0
	Condition   ExpressionNode  // i < 10
	Increment   ExpressionNode  // i += 1
	ThenBody    []StatementNode // then expression body
}

func (*NumericForExpressionNode) IsStatic() bool {
	return false
}

// Create a new numeric `fornum` expression eg. `for i := 0; i < 10; i += 1 then println(i)`
func NewNumericForExpressionNode(span *position.Span, init, cond, incr ExpressionNode, then []StatementNode) *NumericForExpressionNode {
	return &NumericForExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Initialiser:   init,
		Condition:     cond,
		Increment:     incr,
		ThenBody:      then,
	}
}

// Represents a `for in` expression eg. `for i in 5..15 then println(i)`
type ForInExpressionNode struct {
	TypedNodeBase
	Pattern      PatternNode
	InExpression ExpressionNode  // expression that will be iterated through
	ThenBody     []StatementNode // then expression body
}

func (*ForInExpressionNode) IsStatic() bool {
	return false
}

// Create a new `for in` expression node eg. `for i in 5..15 then println(i)`
func NewForInExpressionNode(span *position.Span, pattern PatternNode, inExpr ExpressionNode, then []StatementNode) *ForInExpressionNode {
	return &ForInExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Pattern:       pattern,
		InExpression:  inExpr,
		ThenBody:      then,
	}
}

// Represents a labeled expression eg. `$foo: 1 + 2`
type LabeledExpressionNode struct {
	NodeBase
	Label      string
	Expression ExpressionNode
}

func (l *LabeledExpressionNode) Type(env *types.GlobalEnvironment) types.Type {
	return l.Expression.Type(env)
}

func (l *LabeledExpressionNode) IsStatic() bool {
	return l.Expression.IsStatic()
}

// Create a new labeled expression node eg. `$foo: 1 + 2`
func NewLabeledExpressionNode(span *position.Span, label string, expr ExpressionNode) *LabeledExpressionNode {
	return &LabeledExpressionNode{
		NodeBase:   NodeBase{span: span},
		Label:      label,
		Expression: expr,
	}
}

// Represents a `break` expression eg. `break`, `break false`
type BreakExpressionNode struct {
	NodeBase
	Label string
	Value ExpressionNode
}

func (*BreakExpressionNode) IsStatic() bool {
	return false
}

func (*BreakExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

// Create a new `break` expression node eg. `break`
func NewBreakExpressionNode(span *position.Span, label string, val ExpressionNode) *BreakExpressionNode {
	return &BreakExpressionNode{
		NodeBase: NodeBase{span: span},
		Label:    label,
		Value:    val,
	}
}

// Represents a `go` expression eg. `go foo()`, `go; foo(); end`
type GoExpressionNode struct {
	TypedNodeBase
	Body []StatementNode
}

func (*GoExpressionNode) IsStatic() bool {
	return false
}

// Create a new `go` expression node eg. `go foo()`, `go; foo(); end`
func NewGoExpressionNode(span *position.Span, body []StatementNode) *GoExpressionNode {
	return &GoExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
	}
}

// Represents an `await` expression eg. `await foo()`
type AwaitExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*AwaitExpressionNode) IsStatic() bool {
	return false
}

// Create a new `await` expression node eg. `await foo()`
func NewAwaitExpressionNode(span *position.Span, val ExpressionNode) *AwaitExpressionNode {
	return &AwaitExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a `return` expression eg. `return`, `return true`
type ReturnExpressionNode struct {
	NodeBase
	Value ExpressionNode
}

func (*ReturnExpressionNode) IsStatic() bool {
	return false
}

func (*ReturnExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

// Create a new `return` expression node eg. `return`, `return true`
func NewReturnExpressionNode(span *position.Span, val ExpressionNode) *ReturnExpressionNode {
	return &ReturnExpressionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a `yield` expression eg. `yield`, `yield true`, `yield* foo()`
type YieldExpressionNode struct {
	NodeBase
	Value   ExpressionNode
	Forward bool
}

func (*YieldExpressionNode) IsStatic() bool {
	return false
}

// Create a new `yield` expression node eg. `yield`, `yield true`, `yield* foo()`
func NewYieldExpressionNode(span *position.Span, forward bool, val ExpressionNode) *YieldExpressionNode {
	return &YieldExpressionNode{
		NodeBase: NodeBase{span: span},
		Forward:  forward,
		Value:    val,
	}
}

// Represents a `continue` expression eg. `continue`, `continue "foo"`
type ContinueExpressionNode struct {
	NodeBase
	Label string
	Value ExpressionNode
}

func (*ContinueExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

func (*ContinueExpressionNode) IsStatic() bool {
	return false
}

// Create a new `continue` expression node eg. `continue`, `continue "foo"`
func NewContinueExpressionNode(span *position.Span, label string, val ExpressionNode) *ContinueExpressionNode {
	return &ContinueExpressionNode{
		NodeBase: NodeBase{span: span},
		Label:    label,
		Value:    val,
	}
}

// Represents a `throw` expression eg. `throw ArgumentError("foo")`
type ThrowExpressionNode struct {
	NodeBase
	Unchecked bool
	Value     ExpressionNode
}

func (*ThrowExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

func (*ThrowExpressionNode) IsStatic() bool {
	return false
}

// Create a new `throw` expression node eg. `throw ArgumentError("foo")`
func NewThrowExpressionNode(span *position.Span, unchecked bool, val ExpressionNode) *ThrowExpressionNode {
	return &ThrowExpressionNode{
		NodeBase:  NodeBase{span: span},
		Unchecked: unchecked,
		Value:     val,
	}
}

// Represents a `must` expression eg. `must foo()`
type MustExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*MustExpressionNode) IsStatic() bool {
	return false
}

// Create a new `must` expression node eg. `must foo()`
func NewMustExpressionNode(span *position.Span, val ExpressionNode) *MustExpressionNode {
	return &MustExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a `try` expression eg. `try foo()`
type TryExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*TryExpressionNode) IsStatic() bool {
	return false
}

// Create a new `try` expression node eg. `try foo()`
func NewTryExpressionNode(span *position.Span, val ExpressionNode) *TryExpressionNode {
	return &TryExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a `typeof` expression eg. `typeof foo()`
type TypeofExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*TypeofExpressionNode) IsStatic() bool {
	return false
}

// Create a new `typeof` expression node eg. `typeof foo()`
func NewTypeofExpressionNode(span *position.Span, val ExpressionNode) *TypeofExpressionNode {
	return &TypeofExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a constant declaration eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
type ConstantDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant    ExpressionNode // name of the constant
	TypeNode    TypeNode       // type of the constant
	Initialiser ExpressionNode // value assigned to the constant
}

func (*ConstantDeclarationNode) IsStatic() bool {
	return false
}

// Create a new constant declaration node eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
func NewConstantDeclarationNode(span *position.Span, docComment string, constant ExpressionNode, typ TypeNode, init ExpressionNode) *ConstantDeclarationNode {
	return &ConstantDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:    constant,
		TypeNode:    typ,
		Initialiser: init,
	}
}

// Type expression of an operator with two operands eg. `String | Int`
type BinaryTypeExpressionNode struct {
	TypedNodeBase
	Op    *token.Token // operator
	Left  TypeNode     // left hand side
	Right TypeNode     // right hand side
}

func (*BinaryTypeExpressionNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String | Int`
func NewBinaryTypeExpressionNode(span *position.Span, op *token.Token, left, right TypeNode) *BinaryTypeExpressionNode {
	return &BinaryTypeExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Same as [NewBinaryTypeExpressionNode] but returns an interface
func NewBinaryTypeExpressionNodeI(span *position.Span, op *token.Token, left, right TypeNode) TypeNode {
	return &BinaryTypeExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Union type eg. `String & Int & Float`
type IntersectionTypeNode struct {
	TypedNodeBase
	Elements []TypeNode
}

func (*IntersectionTypeNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String & Int`
func NewIntersectionTypeNode(span *position.Span, elements []TypeNode) *IntersectionTypeNode {
	return &IntersectionTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Union type eg. `String | Int | Float`
type UnionTypeNode struct {
	TypedNodeBase
	Elements []TypeNode
}

func (*UnionTypeNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String | Int`
func NewUnionTypeNode(span *position.Span, elements []TypeNode) *UnionTypeNode {
	return &UnionTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Represents an optional or nilable type eg. `String?`
type NilableTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (*NilableTypeNode) IsStatic() bool {
	return false
}

// Create a new nilable type node eg. `String?`
func NewNilableTypeNode(span *position.Span, typ TypeNode) *NilableTypeNode {
	return &NilableTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		TypeNode:      typ,
	}
}

// Represents an instance type eg. `^self`
type InstanceOfTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (*InstanceOfTypeNode) IsStatic() bool {
	return false
}

// Create a new instance of type node eg. `^self`
func NewInstanceOfTypeNode(span *position.Span, typ TypeNode) *InstanceOfTypeNode {
	return &InstanceOfTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		TypeNode:      typ,
	}
}

// Represents a singleton type eg. `&String`
type SingletonTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (*SingletonTypeNode) IsStatic() bool {
	return false
}

// Create a new singleton type node eg. `&String`
func NewSingletonTypeNode(span *position.Span, typ TypeNode) *SingletonTypeNode {
	return &SingletonTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		TypeNode:      typ,
	}
}

// Represents a not type eg. `~String`
type NotTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (*NotTypeNode) IsStatic() bool {
	return false
}

// Create a new not type node eg. `~String`
func NewNotTypeNode(span *position.Span, typ TypeNode) *NotTypeNode {
	return &NotTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		TypeNode:      typ,
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

// Represents a constant lookup expressions eg. `Foo::Bar`
type ConstantLookupNode struct {
	TypedNodeBase
	Left  ExpressionNode      // left hand side
	Right ComplexConstantNode // right hand side
}

func (*ConstantLookupNode) IsStatic() bool {
	return false
}

// Create a new constant lookup expression node eg. `Foo::Bar`
func NewConstantLookupNode(span *position.Span, left ExpressionNode, right ComplexConstantNode) *ConstantLookupNode {
	return &ConstantLookupNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Left:          left,
		Right:         right,
	}
}

// Represents a method lookup expression eg. `Foo::bar`, `a::c`
type MethodLookupNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Name     string
}

func (*MethodLookupNode) IsStatic() bool {
	return false
}

// Create a new method lookup expression node eg. `Foo::bar`, `a::c`
func NewMethodLookupNode(span *position.Span, receiver ExpressionNode, name string) *MethodLookupNode {
	return &MethodLookupNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Receiver:      receiver,
		Name:          name,
	}
}

// Represents a using all entry node eg. `Foo::*`, `A::B::C::*`
type UsingAllEntryNode struct {
	TypedNodeBase
	Namespace ExpressionNode
}

func (*UsingAllEntryNode) IsStatic() bool {
	return false
}

// Create a new using all entry node eg. `Foo::*`, `A::B::C::*`
func NewUsingAllEntryNode(span *position.Span, namespace UsingEntryNode) *UsingAllEntryNode {
	return &UsingAllEntryNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Namespace:     namespace,
	}
}

// Represents a using entry node with subentries eg. `Foo::{Bar, baz}`, `A::B::C::{lol, foo as epic, Gro as Moe}`
type UsingEntryWithSubentriesNode struct {
	NodeBase
	Namespace  ExpressionNode
	Subentries []UsingSubentryNode
}

func (*UsingEntryWithSubentriesNode) IsStatic() bool {
	return false
}

// Create a new using all entry node eg. `Foo::*`, `A::B::C::*`
func NewUsingEntryWithSubentriesNode(span *position.Span, namespace UsingEntryNode, subentries []UsingSubentryNode) *UsingEntryWithSubentriesNode {
	return &UsingEntryWithSubentriesNode{
		NodeBase:   NodeBase{span: span},
		Namespace:  namespace,
		Subentries: subentries,
	}
}

// Represents an identifier with as in using declarations
// eg. `foo as bar`.
type PublicIdentifierAsNode struct {
	NodeBase
	Target *PublicIdentifierNode
	AsName string
}

func (*PublicIdentifierAsNode) IsStatic() bool {
	return false
}

// Create a new identifier with as eg. `foo as bar`.
func NewPublicIdentifierAsNode(span *position.Span, target *PublicIdentifierNode, as string) *PublicIdentifierAsNode {
	return &PublicIdentifierAsNode{
		NodeBase: NodeBase{span: span},
		Target:   target,
		AsName:   as,
	}
}

// Represents a constant with as in using declarations
// eg. `Foo as Bar`.
type PublicConstantAsNode struct {
	NodeBase
	Target *PublicConstantNode
	AsName string
}

func (*PublicConstantAsNode) IsStatic() bool {
	return false
}

// Create a new identifier with as eg. `Foo as Bar`.
func NewPublicConstantAsNode(span *position.Span, target *PublicConstantNode, as string) *PublicConstantAsNode {
	return &PublicConstantAsNode{
		NodeBase: NodeBase{span: span},
		Target:   target,
		AsName:   as,
	}
}

// Represents a constant with as in using declarations
// eg. `Foo::Bar as Bar`.
type ConstantAsNode struct {
	NodeBase
	Constant ComplexConstantNode
	AsName   string
}

func (*ConstantAsNode) IsStatic() bool {
	return false
}

// Create a new identifier with as eg. `Foo::Bar as Bar`.
func NewConstantAsNode(span *position.Span, constant ComplexConstantNode, as string) *ConstantAsNode {
	return &ConstantAsNode{
		NodeBase: NodeBase{span: span},
		Constant: constant,
		AsName:   as,
	}
}

// Represents a method lookup with as in using declarations
// eg. `Foo::bar as Bar`.
type MethodLookupAsNode struct {
	NodeBase
	MethodLookup *MethodLookupNode
	AsName       string
}

func (*MethodLookupAsNode) IsStatic() bool {
	return false
}

// Create a new identifier with as eg. `Foo::bar as Bar`.
func NewMethodLookupAsNode(span *position.Span, methodLookup *MethodLookupNode, as string) *MethodLookupAsNode {
	return &MethodLookupAsNode{
		NodeBase:     NodeBase{span: span},
		MethodLookup: methodLookup,
		AsName:       as,
	}
}

// Indicates whether the parameter is a rest param
type ParameterKind uint8

const (
	NormalParameterKind ParameterKind = iota
	PositionalRestParameterKind
	NamedRestParameterKind
)

// Represents a formal parameter in function or struct declarations eg. `foo: String = 'bar'`
type FormalParameterNode struct {
	TypedNodeBase
	Name        string         // name of the variable
	TypeNode    TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
	Kind        ParameterKind
}

func (*FormalParameterNode) IsStatic() bool {
	return false
}

func (f *FormalParameterNode) IsOptional() bool {
	return f.Initialiser != nil
}

// Create a new formal parameter node eg. `foo: String = 'bar'`
func NewFormalParameterNode(span *position.Span, name string, typ TypeNode, init ExpressionNode, kind ParameterKind) *FormalParameterNode {
	return &FormalParameterNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Name:          name,
		TypeNode:      typ,
		Initialiser:   init,
		Kind:          kind,
	}
}

// Represents a formal parameter in method declarations eg. `foo: String = 'bar'`
type MethodParameterNode struct {
	TypedNodeBase
	Name                string         // name of the variable
	TypeNode            TypeNode       // type of the variable
	Initialiser         ExpressionNode // value assigned to the variable
	SetInstanceVariable bool           // whether an instance variable with this name gets automatically assigned
	Kind                ParameterKind
}

func (*MethodParameterNode) IsStatic() bool {
	return false
}

func (f *MethodParameterNode) IsOptional() bool {
	return f.Initialiser != nil
}

// Create a new formal parameter node eg. `foo: String = 'bar'`
func NewMethodParameterNode(span *position.Span, name string, setIvar bool, typ TypeNode, init ExpressionNode, kind ParameterKind) *MethodParameterNode {
	return &MethodParameterNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		SetInstanceVariable: setIvar,
		Name:                name,
		TypeNode:            typ,
		Initialiser:         init,
		Kind:                kind,
	}
}

// Represents a signature parameter in method and function signatures eg. `foo?: String`
type SignatureParameterNode struct {
	TypedNodeBase
	Name     string   // name of the variable
	TypeNode TypeNode // type of the variable
	Optional bool     // whether this parameter is optional
	Kind     ParameterKind
}

func (*SignatureParameterNode) IsStatic() bool {
	return false
}

func (f *SignatureParameterNode) IsOptional() bool {
	return f.Optional
}

// Create a new signature parameter node eg. `foo?: String`
func NewSignatureParameterNode(span *position.Span, name string, typ TypeNode, opt bool, kind ParameterKind) *SignatureParameterNode {
	return &SignatureParameterNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Name:          name,
		TypeNode:      typ,
		Optional:      opt,
		Kind:          kind,
	}
}

// Represents an attribute declaration in getters, setters and accessors eg. `foo: String`
type AttributeParameterNode struct {
	TypedNodeBase
	Name        string         // name of the variable
	TypeNode    TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

func (*AttributeParameterNode) IsStatic() bool {
	return false
}

func (a *AttributeParameterNode) IsOptional() bool {
	return a.Initialiser != nil
}

// Create a new attribute declaration in getters, setters and accessors eg. `foo: String`
func NewAttributeParameterNode(span *position.Span, name string, typ TypeNode, init ExpressionNode) *AttributeParameterNode {
	return &AttributeParameterNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Name:          name,
		TypeNode:      typ,
		Initialiser:   init,
	}
}

// Represents a closure type eg. `|i: Int|: String`
type ClosureTypeNode struct {
	TypedNodeBase
	Parameters []ParameterNode // formal parameters of the closure separated by semicolons
	ReturnType TypeNode
	ThrowType  TypeNode
}

func (*ClosureTypeNode) IsStatic() bool {
	return false
}

// Create a new closure type node eg. `|i: Int|: String`
func NewClosureTypeNode(span *position.Span, params []ParameterNode, retType TypeNode, throwType TypeNode) *ClosureTypeNode {
	return &ClosureTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Parameters:    params,
		ReturnType:    retType,
		ThrowType:     throwType,
	}
}

// Represents a closure eg. `|i| -> println(i)`
type ClosureLiteralNode struct {
	TypedNodeBase
	Parameters []ParameterNode // formal parameters of the closure separated by semicolons
	ReturnType TypeNode
	ThrowType  TypeNode
	Body       []StatementNode // body of the closure
}

func (*ClosureLiteralNode) IsStatic() bool {
	return false
}

// Create a new closure expression node eg. `|i| -> println(i)`
func NewClosureLiteralNode(span *position.Span, params []ParameterNode, retType TypeNode, throwType TypeNode, body []StatementNode) *ClosureLiteralNode {
	return &ClosureLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Parameters:    params,
		ReturnType:    retType,
		ThrowType:     throwType,
		Body:          body,
	}
}

// Represents a class declaration eg. `class Foo; end`
type ClassDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Abstract       bool
	Sealed         bool
	Primitive      bool
	NoInit         bool
	Constant       ExpressionNode      // The constant that will hold the class value
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Superclass     ExpressionNode      // the super/parent class of this class
	Body           []StatementNode     // body of the class
	Bytecode       *vm.BytecodeFunction
}

func (*ClassDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*ClassDeclarationNode) IsStatic() bool {
	return false
}

// Create a new class declaration node eg. `class Foo; end`
func NewClassDeclarationNode(
	span *position.Span,
	docComment string,
	abstract bool,
	sealed bool,
	primitive bool,
	noinit bool,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	superclass ExpressionNode,
	body []StatementNode,
) *ClassDeclarationNode {

	return &ClassDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Abstract:       abstract,
		Sealed:         sealed,
		Primitive:      primitive,
		NoInit:         noinit,
		Constant:       constant,
		TypeParameters: typeParams,
		Superclass:     superclass,
		Body:           body,
	}
}

// Represents a module declaration eg. `module Foo; end`
type ModuleDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant ExpressionNode  // The constant that will hold the module value
	Body     []StatementNode // body of the module
	Bytecode *vm.BytecodeFunction
}

func (*ModuleDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*ModuleDeclarationNode) IsStatic() bool {
	return false
}

// Create a new module declaration node eg. `module Foo; end`
func NewModuleDeclarationNode(
	span *position.Span,
	docComment string,
	constant ExpressionNode,
	body []StatementNode,
) *ModuleDeclarationNode {

	return &ModuleDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant: constant,
		Body:     body,
	}
}

// Represents a mixin declaration eg. `mixin Foo; end`
type MixinDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Abstract              bool
	Constant              ExpressionNode      // The constant that will hold the mixin value
	TypeParameters        []TypeParameterNode // Generic type variable definitions
	Body                  []StatementNode     // body of the mixin
	IncludesAndImplements []ExpressionNode
	Bytecode              *vm.BytecodeFunction
}

func (*MixinDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*MixinDeclarationNode) IsStatic() bool {
	return false
}

// Create a new mixin declaration node eg. `mixin Foo; end`
func NewMixinDeclarationNode(
	span *position.Span,
	docComment string,
	abstract bool,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StatementNode,
) *MixinDeclarationNode {

	return &MixinDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Abstract:       abstract,
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

// Represents an interface declaration eg. `interface Foo; end`
type InterfaceDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant       ExpressionNode      // The constant that will hold the interface value
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Body           []StatementNode     // body of the interface
	Implements     []*ImplementExpressionNode
	Bytecode       *vm.BytecodeFunction
}

func (*InterfaceDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*InterfaceDeclarationNode) IsStatic() bool {
	return false
}

// Create a new interface declaration node eg. `interface Foo; end`
func NewInterfaceDeclarationNode(
	span *position.Span,
	docComment string,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StatementNode,
) *InterfaceDeclarationNode {

	return &InterfaceDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

// Represents a struct declaration eg. `struct Foo; end`
type StructDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant       ExpressionNode            // The constant that will hold the struct value
	TypeParameters []TypeParameterNode       // Generic type variable definitions
	Body           []StructBodyStatementNode // body of the struct
}

func (*StructDeclarationNode) IsStatic() bool {
	return false
}

// Create a new struct declaration node eg. `struct Foo; end`
func NewStructDeclarationNode(
	span *position.Span,
	docComment string,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StructBodyStatementNode,
) *StructDeclarationNode {

	return &StructDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

// Represents the variance of a type parameter.
type Variance uint8

const (
	INVARIANT Variance = iota
	COVARIANT
	CONTRAVARIANT
)

// Represents a type parameter eg. `+V`
type VariantTypeParameterNode struct {
	TypedNodeBase
	Variance   Variance // Variance level of this type parameter
	Name       string   // Name of the type parameter eg. `T`
	LowerBound TypeNode
	UpperBound TypeNode
	Default    TypeNode
}

func (*VariantTypeParameterNode) IsStatic() bool {
	return false
}

// Create a new type variable node eg. `+V`
func NewVariantTypeParameterNode(span *position.Span, variance Variance, name string, lower, upper, def TypeNode) *VariantTypeParameterNode {
	return &VariantTypeParameterNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Variance:      variance,
		Name:          name,
		LowerBound:    lower,
		UpperBound:    upper,
		Default:       def,
	}
}

// Represents a method definition eg. `def foo: String then 'hello world'`
type MethodDefinitionNode struct {
	TypedNodeBaseWithLoc
	DocCommentableNodeBase
	Name           string
	TypeParameters []TypeParameterNode
	Parameters     []ParameterNode // formal parameters
	ReturnType     TypeNode
	ThrowType      TypeNode
	Body           []StatementNode // body of the method
	Flags          bitfield.BitField8
}

func (*MethodDefinitionNode) IsStatic() bool {
	return false
}

// Whether the method is a setter.
func (m *MethodDefinitionNode) IsSetter() bool {
	if len(m.Name) > 0 {
		firstChar, _ := utf8.DecodeRuneInString(m.Name)
		lastChar := m.Name[len(m.Name)-1]
		if (unicode.IsLetter(firstChar) || firstChar == '_') && lastChar == '=' {
			return true
		}
	}

	return false
}

func (m *MethodDefinitionNode) IsAbstract() bool {
	return m.Flags.HasFlag(METHOD_ABSTRACT_FLAG)
}

func (m *MethodDefinitionNode) SetAbstract() {
	m.Flags.SetFlag(METHOD_ABSTRACT_FLAG)
}

func (m *MethodDefinitionNode) IsSealed() bool {
	return m.Flags.HasFlag(METHOD_SEALED_FLAG)
}

func (m *MethodDefinitionNode) SetSealed() {
	m.Flags.SetFlag(METHOD_SEALED_FLAG)
}

func (m *MethodDefinitionNode) IsGenerator() bool {
	return m.Flags.HasFlag(METHOD_GENERATOR_FLAG)
}

func (m *MethodDefinitionNode) SetGenerator() {
	m.Flags.SetFlag(METHOD_GENERATOR_FLAG)
}

func (m *MethodDefinitionNode) IsAsync() bool {
	return m.Flags.HasFlag(METHOD_ASYNC_FLAG)
}

func (m *MethodDefinitionNode) SetAsync() {
	m.Flags.SetFlag(METHOD_ASYNC_FLAG)
}

const (
	METHOD_ABSTRACT_FLAG bitfield.BitFlag8 = 1 << iota
	METHOD_SEALED_FLAG
	METHOD_GENERATOR_FLAG
	METHOD_ASYNC_FLAG
)

// Create a method definition node eg. `def foo: String then 'hello world'`
func NewMethodDefinitionNode(
	loc *position.Location,
	docComment string,
	flags bitfield.BitFlag8,
	name string,
	typeParams []TypeParameterNode,
	params []ParameterNode,
	returnType,
	throwType TypeNode,
	body []StatementNode,
) *MethodDefinitionNode {
	return &MethodDefinitionNode{
		TypedNodeBaseWithLoc: TypedNodeBaseWithLoc{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Flags:          bitfield.BitField8FromBitFlag(flags),
		Name:           name,
		TypeParameters: typeParams,
		Parameters:     params,
		ReturnType:     returnType,
		ThrowType:      throwType,
		Body:           body,
	}
}

// Represents a constructor definition eg. `init then 'hello world'`
type InitDefinitionNode struct {
	TypedNodeBaseWithLoc
	DocCommentableNodeBase
	Parameters []ParameterNode // formal parameters
	ThrowType  TypeNode
	Body       []StatementNode // body of the method
}

func (*InitDefinitionNode) IsStatic() bool {
	return false
}

// Create a constructor definition node eg. `init then 'hello world'`
func NewInitDefinitionNode(loc *position.Location, params []ParameterNode, throwType TypeNode, body []StatementNode) *InitDefinitionNode {
	return &InitDefinitionNode{
		TypedNodeBaseWithLoc: TypedNodeBaseWithLoc{loc: loc},
		Parameters:           params,
		ThrowType:            throwType,
		Body:                 body,
	}
}

// Represents a method signature definition eg. `sig to_string(val: Int): String`
type MethodSignatureDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name           string
	TypeParameters []TypeParameterNode
	Parameters     []ParameterNode // formal parameters
	ReturnType     TypeNode
	ThrowType      TypeNode
}

func (*MethodSignatureDefinitionNode) IsStatic() bool {
	return false
}

// Create a method signature node eg. `sig to_string(val: Int): String`
func NewMethodSignatureDefinitionNode(span *position.Span, docComment, name string, typeParams []TypeParameterNode, params []ParameterNode, returnType, throwType TypeNode) *MethodSignatureDefinitionNode {
	return &MethodSignatureDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Name:           name,
		TypeParameters: typeParams,
		Parameters:     params,
		ReturnType:     returnType,
		ThrowType:      throwType,
	}
}

// Represents a generic constant in type annotations eg. `ArrayList[String]`
type GenericConstantNode struct {
	TypedNodeBase
	Constant      ComplexConstantNode
	TypeArguments []TypeNode
}

func (*GenericConstantNode) IsStatic() bool {
	return true
}

// Create a generic constant node eg. `ArrayList[String]`
func NewGenericConstantNode(span *position.Span, constant ComplexConstantNode, args []TypeNode) *GenericConstantNode {
	return &GenericConstantNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Constant:      constant,
		TypeArguments: args,
	}
}

// Represents a new generic type definition eg. `typedef Nilable[T] = T | nil`
type GenericTypeDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Constant       ComplexConstantNode // new name of the type
	TypeNode       TypeNode            // the type
}

func (*GenericTypeDefinitionNode) IsStatic() bool {
	return false
}

// Create a generic type definition node eg. `typedef Nilable[T] = T | nil`
func NewGenericTypeDefinitionNode(span *position.Span, docComment string, constant ComplexConstantNode, typeVars []TypeParameterNode, typ TypeNode) *GenericTypeDefinitionNode {
	return &GenericTypeDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:       constant,
		TypeParameters: typeVars,
		TypeNode:       typ,
	}
}

// Represents a new type definition eg. `typedef StringList = ArrayList[String]`
type TypeDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant ComplexConstantNode // new name of the type
	TypeNode TypeNode            // the type
}

func (*TypeDefinitionNode) IsStatic() bool {
	return false
}

// Create a type definition node eg. `typedef StringList = ArrayList[String]`
func NewTypeDefinitionNode(span *position.Span, docComment string, constant ComplexConstantNode, typ TypeNode) *TypeDefinitionNode {
	return &TypeDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant: constant,
		TypeNode: typ,
	}
}

// A single alias entry eg. `new_name old_name`
type AliasDeclarationEntry struct {
	NodeBase
	NewName string
	OldName string
}

func (*AliasDeclarationEntry) IsStatic() bool {
	return false
}

// Create an alias alias entry eg. `new_name old_name`
func NewAliasDeclarationEntry(span *position.Span, newName, oldName string) *AliasDeclarationEntry {
	return &AliasDeclarationEntry{
		NodeBase: NodeBase{span: span},
		NewName:  newName,
		OldName:  oldName,
	}
}

// Represents a new alias declaration eg. `alias push append, add plus`
type AliasDeclarationNode struct {
	TypedNodeBase
	Entries []*AliasDeclarationEntry
}

func (*AliasDeclarationNode) IsStatic() bool {
	return false
}

// Create an alias declaration node eg. `alias push append, add plus`
func NewAliasDeclarationNode(span *position.Span, entries []*AliasDeclarationEntry) *AliasDeclarationNode {
	return &AliasDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Entries:       entries,
	}
}

// Represents a new getter declaration eg. `getter foo: String`
type GetterDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Entries []ParameterNode
}

func (*GetterDeclarationNode) IsStatic() bool {
	return false
}

// Create a getter declaration node eg. `getter foo: String`
func NewGetterDeclarationNode(span *position.Span, docComment string, entries []ParameterNode) *GetterDeclarationNode {
	return &GetterDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Entries: entries,
	}
}

// Represents a new setter declaration eg. `setter foo: String`
type SetterDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Entries []ParameterNode
}

func (*SetterDeclarationNode) IsStatic() bool {
	return false
}

// Create a setter declaration node eg. `setter foo: String`
func NewSetterDeclarationNode(span *position.Span, docComment string, entries []ParameterNode) *SetterDeclarationNode {
	return &SetterDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Entries: entries,
	}
}

// Represents a new setter declaration eg. `attr foo: String`
type AttrDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Entries []ParameterNode
}

func (*AttrDeclarationNode) IsStatic() bool {
	return false
}

// Create an attribute declaration node eg. `attr foo: String`
func NewAttrDeclarationNode(span *position.Span, docComment string, entries []ParameterNode) *AttrDeclarationNode {
	return &AttrDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Entries: entries,
	}
}

// Represents a using expression eg. `using Foo`
type UsingExpressionNode struct {
	TypedNodeBase
	Entries []UsingEntryNode
}

func (*UsingExpressionNode) SkipTypechecking() bool {
	return false
}

func (*UsingExpressionNode) IsStatic() bool {
	return false
}

// Create a using expression node eg. `using Foo`
func NewUsingExpressionNode(span *position.Span, consts []UsingEntryNode) *UsingExpressionNode {
	return &UsingExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Entries:       consts,
	}
}

// Represents an include expression eg. `include Enumerable[V]`
type IncludeExpressionNode struct {
	TypedNodeBase
	Constants []ComplexConstantNode
}

func (*IncludeExpressionNode) SkipTypechecking() bool {
	return false
}

func (*IncludeExpressionNode) IsStatic() bool {
	return false
}

// Create an include expression node eg. `include Enumerable[V]`
func NewIncludeExpressionNode(span *position.Span, consts []ComplexConstantNode) *IncludeExpressionNode {
	return &IncludeExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Constants:     consts,
	}
}

// Represents an enhance expression eg. `implement Enumerable[V]`
type ImplementExpressionNode struct {
	TypedNodeBase
	Constants []ComplexConstantNode
}

func (*ImplementExpressionNode) SkipTypechecking() bool {
	return false
}

func (*ImplementExpressionNode) IsStatic() bool {
	return false
}

// Create an enhance expression node eg. `implement Enumerable[V]`
func NewImplementExpressionNode(span *position.Span, consts []ComplexConstantNode) *ImplementExpressionNode {
	return &ImplementExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Constants:     consts,
	}
}

// Represents a named argument in a function call eg. `foo: 123`
type NamedCallArgumentNode struct {
	NodeBase
	Name  string
	Value ExpressionNode
}

func (*NamedCallArgumentNode) IsStatic() bool {
	return false
}

// Create a named argument node eg. `foo: 123`
func NewNamedCallArgumentNode(span *position.Span, name string, val ExpressionNode) *NamedCallArgumentNode {
	return &NamedCallArgumentNode{
		NodeBase: NodeBase{span: span},
		Name:     name,
		Value:    val,
	}
}

// Represents a new expression eg. `new(123)`
type NewExpressionNode struct {
	TypedNodeBase
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*NewExpressionNode) IsStatic() bool {
	return false
}

// Create a new expression node eg. `new(123)`
func NewNewExpressionNode(span *position.Span, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *NewExpressionNode {
	return &NewExpressionNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a constructor call eg. `String(123)`
type ConstructorCallNode struct {
	TypedNodeBase
	ClassNode           ComplexConstantNode // class that is being instantiated
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*ConstructorCallNode) IsStatic() bool {
	return false
}

// Create a constructor call node eg. `String(123)`
func NewConstructorCallNode(span *position.Span, class ComplexConstantNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ConstructorCallNode {
	return &ConstructorCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		ClassNode:           class,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a constructor call eg. `ArrayList::[Int](1, 2, 3)`
type GenericConstructorCallNode struct {
	TypedNodeBase
	ClassNode           ComplexConstantNode // class that is being instantiated
	TypeArguments       []TypeNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*GenericConstructorCallNode) IsStatic() bool {
	return false
}

// Create a constructor call node eg. `ArrayList::[Int](1, 2, 3)`
func NewGenericConstructorCallNode(span *position.Span, class ComplexConstantNode, typeArgs []TypeNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *GenericConstructorCallNode {
	return &GenericConstructorCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		ClassNode:           class,
		TypeArguments:       typeArgs,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents attribute access eg. `foo.bar`
type AttributeAccessNode struct {
	TypedNodeBase
	Receiver      ExpressionNode
	AttributeName string
}

func (*AttributeAccessNode) IsStatic() bool {
	return false
}

// Create an attribute access node eg. `foo.bar`
func NewAttributeAccessNode(span *position.Span, recv ExpressionNode, attrName string) *AttributeAccessNode {
	return &AttributeAccessNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Receiver:      recv,
		AttributeName: attrName,
	}
}

// Represents subscript access eg. `arr[5]`
type SubscriptExpressionNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Key      ExpressionNode
	static   bool
}

func (s *SubscriptExpressionNode) IsStatic() bool {
	return s.static
}

// Create a subscript expression node eg. `arr[5]`
func NewSubscriptExpressionNode(span *position.Span, recv, key ExpressionNode) *SubscriptExpressionNode {
	return &SubscriptExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Receiver:      recv,
		Key:           key,
		static:        recv.IsStatic() && key.IsStatic(),
	}
}

// Represents nil-safe subscript access eg. `arr?[5]`
type NilSafeSubscriptExpressionNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Key      ExpressionNode
	static   bool
}

func (s *NilSafeSubscriptExpressionNode) IsStatic() bool {
	return s.static
}

// Create a nil-safe subscript expression node eg. `arr?[5]`
func NewNilSafeSubscriptExpressionNode(span *position.Span, recv, key ExpressionNode) *NilSafeSubscriptExpressionNode {
	return &NilSafeSubscriptExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Receiver:      recv,
		Key:           key,
		static:        recv.IsStatic() && key.IsStatic(),
	}
}

// Represents a method call eg. `'123'.()`
type CallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	NilSafe             bool
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*CallNode) IsStatic() bool {
	return false
}

// Create a method call node eg. `'123'.to_int()`
func NewCallNode(span *position.Span, recv ExpressionNode, nilSafe bool, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *CallNode {
	return &CallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		Receiver:            recv,
		NilSafe:             nilSafe,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a method call eg. `foo.bar::[String](a)`
type GenericMethodCallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	Op                  *token.Token
	MethodName          string
	TypeArguments       []TypeNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
	TailCall            bool
}

func (*GenericMethodCallNode) IsStatic() bool {
	return false
}

// Create a method call node eg. `foo.bar::[String](a)`
func NewGenericMethodCallNode(span *position.Span, recv ExpressionNode, op *token.Token, methodName string, typeArgs []TypeNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *GenericMethodCallNode {
	return &GenericMethodCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		Receiver:            recv,
		Op:                  op,
		MethodName:          methodName,
		TypeArguments:       typeArgs,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a method call eg. `'123'.to_int()`
type MethodCallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	Op                  *token.Token
	MethodName          string
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
	TailCall            bool
}

func (*MethodCallNode) IsStatic() bool {
	return false
}

// Create a method call node eg. `'123'.to_int()`
func NewMethodCallNode(span *position.Span, recv ExpressionNode, op *token.Token, methodName string, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *MethodCallNode {
	return &MethodCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		Receiver:            recv,
		Op:                  op,
		MethodName:          methodName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a function-like call eg. `to_string(123)`
type ReceiverlessMethodCallNode struct {
	TypedNodeBase
	MethodName          string
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
	TailCall            bool
}

func (*ReceiverlessMethodCallNode) IsStatic() bool {
	return false
}

// Create a function call node eg. `to_string(123)`
func NewReceiverlessMethodCallNode(span *position.Span, methodName string, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ReceiverlessMethodCallNode {
	return &ReceiverlessMethodCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		MethodName:          methodName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a generic function-like call eg. `foo::[Int](123)`
type GenericReceiverlessMethodCallNode struct {
	TypedNodeBase
	MethodName          string
	TypeArguments       []TypeNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
	TailCall            bool
}

func (*GenericReceiverlessMethodCallNode) IsStatic() bool {
	return false
}

// Create a generic function call node eg. `foo::[Int](123)`
func NewGenericReceiverlessMethodCallNode(span *position.Span, methodName string, typeArgs []TypeNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *GenericReceiverlessMethodCallNode {
	return &GenericReceiverlessMethodCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		MethodName:          methodName,
		TypeArguments:       typeArgs,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a symbol value expression eg. `foo: bar`
type SymbolKeyValueExpressionNode struct {
	NodeBase
	Key   string
	Value ExpressionNode
}

func (s *SymbolKeyValueExpressionNode) IsStatic() bool {
	return s.Value.IsStatic()
}

// Create a symbol key value node eg. `foo: bar`
func NewSymbolKeyValueExpressionNode(span *position.Span, key string, val ExpressionNode) *SymbolKeyValueExpressionNode {
	return &SymbolKeyValueExpressionNode{
		NodeBase: NodeBase{span: span},
		Key:      key,
		Value:    val,
	}
}

// Represents a key value expression eg. `foo => bar`
type KeyValueExpressionNode struct {
	TypedNodeBase
	Key    ExpressionNode
	Value  ExpressionNode
	static bool
}

func (k *KeyValueExpressionNode) IsStatic() bool {
	return k.static
}

// Create a key value expression node eg. `foo => bar`
func NewKeyValueExpressionNode(span *position.Span, key, val ExpressionNode) *KeyValueExpressionNode {
	return &KeyValueExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Key:           key,
		Value:         val,
		static:        areExpressionsStatic(key, val),
	}
}

// Represents a splat expression eg. `*foo`
type SplatExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*SplatExpressionNode) IsStatic() bool {
	return false
}

// Create a splat expression node eg. `*foo`
func NewSplatExpressionNode(span *position.Span, val ExpressionNode) *SplatExpressionNode {
	return &SplatExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a double splat expression eg. `**foo`
type DoubleSplatExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*DoubleSplatExpressionNode) IsStatic() bool {
	return false
}

// Create a double splat expression node eg. `**foo`
func NewDoubleSplatExpressionNode(span *position.Span, val ExpressionNode) *DoubleSplatExpressionNode {
	return &DoubleSplatExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a ArrayList literal eg. `[1, 5, -6]`
type ArrayListLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	Capacity ExpressionNode
	static   bool
}

func (l *ArrayListLiteralNode) IsStatic() bool {
	return l.static
}

// Create a ArrayList literal node eg. `[1, 5, -6]`
func NewArrayListLiteralNode(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) *ArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = isExpressionSliceStatic(elements) && capacity.IsStatic()
	} else {
		static = isExpressionSliceStatic(elements)
	}
	return &ArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewArrayListLiteralNode] but returns an interface
func NewArrayListLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewArrayListLiteralNode(span, elements, capacity)
}

// Represents a word ArrayList literal eg. `\w[foo bar]`
type WordArrayListLiteralNode struct {
	TypedNodeBase
	Elements []WordCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (w *WordArrayListLiteralNode) IsStatic() bool {
	return w.static
}

// Create a word ArrayList literal node eg. `\w[foo bar]`
func NewWordArrayListLiteralNode(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) *WordArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &WordArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewWordArrayListLiteralNode] but returns an interface.
func NewWordArrayListLiteralExpressionNode(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewWordArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewWordArrayListLiteralNode] but returns an interface.
func NewWordArrayListLiteralPatternExpressionNode(span *position.Span, elements []WordCollectionContentNode) PatternExpressionNode {
	return NewWordArrayListLiteralNode(span, elements, nil)
}

// Represents a word ArrayTuple literal eg. `%w[foo bar]`
type WordArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []WordCollectionContentNode
}

func (*WordArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a word ArrayTuple literal node eg. `%w[foo bar]`
func NewWordArrayTupleLiteralNode(span *position.Span, elements []WordCollectionContentNode) *WordArrayTupleLiteralNode {
	return &WordArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewWordArrayTupleLiteralNode] but returns an interface.
func NewWordArrayTupleLiteralExpressionNode(span *position.Span, elements []WordCollectionContentNode) ExpressionNode {
	return NewWordArrayTupleLiteralNode(span, elements)
}

// Same as [NewWordArrayTupleLiteralNode] but returns an interface.
func NewWordArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []WordCollectionContentNode) PatternExpressionNode {
	return NewWordArrayTupleLiteralNode(span, elements)
}

// Represents a word HashSet literal eg. `^w[foo bar]`
type WordHashSetLiteralNode struct {
	TypedNodeBase
	Elements []WordCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (w *WordHashSetLiteralNode) IsStatic() bool {
	return w.static
}

// Create a word HashSet literal node eg. `^w[foo bar]`
func NewWordHashSetLiteralNode(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) *WordHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &WordHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewWordHashSetLiteralNode] but returns an interface.
func NewWordHashSetLiteralNodeI(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewWordHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewWordHashSetLiteralNode] but returns an interface.
func NewWordHashSetLiteralPatternExpressionNode(span *position.Span, elements []WordCollectionContentNode) PatternExpressionNode {
	return NewWordHashSetLiteralNode(span, elements, nil)
}

// Represents a symbol ArrayList literal eg. `\s[foo bar]`
type SymbolArrayListLiteralNode struct {
	TypedNodeBase
	Elements []SymbolCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (s *SymbolArrayListLiteralNode) IsStatic() bool {
	return s.static
}

// Create a symbol ArrayList literal node eg. `\s[foo bar]`
func NewSymbolArrayListLiteralNode(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) *SymbolArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &SymbolArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewSymbolArrayListLiteralNode] but returns an interface.
func NewSymbolArrayListLiteralExpressionNode(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewSymbolArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewSymbolArrayListLiteralNode] but returns an interface.
func NewSymbolArrayListLiteralPatternExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) PatternExpressionNode {
	return NewSymbolArrayListLiteralNode(span, elements, nil)
}

// Represents a symbol ArrayTuple literal eg. `%s[foo bar]`
type SymbolArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []SymbolCollectionContentNode
}

func (*SymbolArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a symbol arrayTuple literal node eg. `%s[foo bar]`
func NewSymbolArrayTupleLiteralNode(span *position.Span, elements []SymbolCollectionContentNode) *SymbolArrayTupleLiteralNode {
	return &SymbolArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewSymbolArrayTupleLiteralNode] but returns an interface.
func NewSymbolArrayTupleLiteralExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) ExpressionNode {
	return NewSymbolArrayTupleLiteralNode(span, elements)
}

// Same as [NewSymbolArrayTupleLiteralNode] but returns an interface.
func NewSymbolArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) PatternExpressionNode {
	return NewSymbolArrayTupleLiteralNode(span, elements)
}

// Represents a symbol HashSet literal eg. `^s[foo bar]`
type SymbolHashSetLiteralNode struct {
	TypedNodeBase
	Elements []SymbolCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (s *SymbolHashSetLiteralNode) IsStatic() bool {
	return s.static
}

// Create a symbol HashSet literal node eg. `^s[foo bar]`
func NewSymbolHashSetLiteralNode(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) *SymbolHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &SymbolHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewSymbolHashSetLiteralNode] but returns an interface.
func NewSymbolHashSetLiteralNodeI(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewSymbolHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewSymbolHashSetLiteralNode] but returns an interface.
func NewSymbolHashSetLiteralPatternExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) PatternExpressionNode {
	return NewSymbolHashSetLiteralNode(span, elements, nil)
}

// Represents a hex ArrayList literal eg. `\x[ff ee]`
type HexArrayListLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (h *HexArrayListLiteralNode) IsStatic() bool {
	return h.static
}

// Create a hex ArrayList literal node eg. `\x[ff ee]`
func NewHexArrayListLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *HexArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &HexArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHexArrayListLiteralNode] but returns an interface.
func NewHexArrayListLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewHexArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewHexArrayListLiteralNode] but returns an interface.
func NewHexArrayListLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewHexArrayListLiteralNode(span, elements, nil)
}

// Represents a hex ArrayTuple literal eg. `%x[ff ee]`
type HexArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
}

func (*HexArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a hex ArrayTuple literal node eg. `%x[ff ee]`
func NewHexArrayTupleLiteralNode(span *position.Span, elements []IntCollectionContentNode) *HexArrayTupleLiteralNode {
	return &HexArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewHexArrayTupleLiteralNode] but returns an interface.
func NewHexArrayTupleLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode) ExpressionNode {
	return NewHexArrayTupleLiteralNode(span, elements)
}

// Same as [NewHexArrayTupleLiteralNode] but returns an interface.
func NewHexArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewHexArrayTupleLiteralNode(span, elements)
}

// Represents a hex HashSet literal eg. `^x[ff ee}]`
type HexHashSetLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (h *HexHashSetLiteralNode) IsStatic() bool {
	return h.static
}

// Create a hex HashSet literal node eg. `^x[ff ee]`
func NewHexHashSetLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *HexHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &HexHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHexHashSetLiteralNode] but returns an interface.
func NewHexHashSetLiteralNodeI(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewHexHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewHexHashSetLiteralNode] but returns an interface.
func NewHexHashSetLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewHexHashSetLiteralNode(span, elements, nil)
}

// Represents a bin ArrayList literal eg. `\b[11 10]`
type BinArrayListLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (b *BinArrayListLiteralNode) IsStatic() bool {
	return b.static
}

// Create a bin ArrayList literal node eg. `\b[11 10]`
func NewBinArrayListLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *BinArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &BinArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewBinArrayListLiteralNode] but returns an interface.
func NewBinArrayListLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewBinArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewBinArrayListLiteralNode] but returns an interface.
func NewBinArrayListLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewBinArrayListLiteralNode(span, elements, nil)
}

// Represents a bin ArrayTuple literal eg. `%b[11 10]`
type BinArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
}

func (*BinArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a bin ArrayList literal node eg. `%b[11 10]`
func NewBinArrayTupleLiteralNode(span *position.Span, elements []IntCollectionContentNode) *BinArrayTupleLiteralNode {
	return &BinArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewBinArrayTupleLiteralNode] but returns an interface.
func NewBinArrayTupleLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode) ExpressionNode {
	return NewBinArrayTupleLiteralNode(span, elements)
}

// Same as [NewBinArrayTupleLiteralNode] but returns an interface.
func NewBinArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewBinArrayTupleLiteralNode(span, elements)
}

// Represents a bin HashSet literal eg. `^b[11 10]`
type BinHashSetLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (b *BinHashSetLiteralNode) IsStatic() bool {
	return b.static
}

// Create a bin HashSet literal node eg. `^b[11 10]`
func NewBinHashSetLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *BinHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &BinHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewBinHashSetLiteralNode] but returns an interface.
func NewBinHashSetLiteralNodeI(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewBinHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewBinHashSetLiteralNode] but returns an interface.
func NewBinHashSetLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewBinHashSetLiteralNode(span, elements, nil)
}

// Represents a ArrayTuple literal eg. `%[1, 5, -6]`
type ArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	static   bool
}

func (t *ArrayTupleLiteralNode) IsStatic() bool {
	return t.static
}

// Create a ArrayTuple literal node eg. `%[1, 5, -6]`
func NewArrayTupleLiteralNode(span *position.Span, elements []ExpressionNode) *ArrayTupleLiteralNode {
	return &ArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		static:        isExpressionSliceStatic(elements),
	}
}

// Same as [NewArrayTupleLiteralNode] but returns an interface
func NewArrayTupleLiteralNodeI(span *position.Span, elements []ExpressionNode) ExpressionNode {
	return NewArrayTupleLiteralNode(span, elements)
}

// Represents a HashSet literal eg. `^[1, 5, -6]`
type HashSetLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	Capacity ExpressionNode
	static   bool
}

func (s *HashSetLiteralNode) IsStatic() bool {
	return s.static
}

// Create a HashSet literal node eg. `^[1, 5, -6]`
func NewHashSetLiteralNode(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) *HashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = isExpressionSliceStatic(elements) && capacity.IsStatic()
	} else {
		static = isExpressionSliceStatic(elements)
	}
	return &HashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHashSetLiteralNode] but returns an interface
func NewHashSetLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewHashSetLiteralNode(span, elements, capacity)
}

// Represents a HashMap literal eg. `{ foo: 1, 'bar' => 5, baz }`
type HashMapLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	Capacity ExpressionNode
	static   bool
}

func (m *HashMapLiteralNode) IsStatic() bool {
	return m.static
}

// Create a HashMap literal node eg. `{ foo: 1, 'bar' => 5, baz }`
func NewHashMapLiteralNode(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) *HashMapLiteralNode {
	var static bool
	if capacity != nil {
		static = isExpressionSliceStatic(elements) && capacity.IsStatic()
	} else {
		static = isExpressionSliceStatic(elements)
	}
	return &HashMapLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHashMapLiteralNode] but returns an interface
func NewHashMapLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewHashMapLiteralNode(span, elements, capacity)
}

// Represents a Record literal eg. `%{ foo: 1, 'bar' => 5, baz }`
type HashRecordLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	static   bool
}

func (r *HashRecordLiteralNode) IsStatic() bool {
	return r.static
}

// Create a Record literal node eg. `%{ foo: 1, 'bar' => 5, baz }`
func NewHashRecordLiteralNode(span *position.Span, elements []ExpressionNode) *HashRecordLiteralNode {
	return &HashRecordLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		static:        isExpressionSliceStatic(elements),
	}
}

// Same as [NewHashRecordLiteralNode] but returns an interface
func NewHashRecordLiteralNodeI(span *position.Span, elements []ExpressionNode) ExpressionNode {
	return NewHashRecordLiteralNode(span, elements)
}

// Represents a Range literal eg. `1...5`
type RangeLiteralNode struct {
	TypedNodeBase
	Start  ExpressionNode
	End    ExpressionNode
	Op     *token.Token
	static bool
}

func (r *RangeLiteralNode) IsStatic() bool {
	return r.static
}

// Create a Range literal node eg. `1...5`
func NewRangeLiteralNode(span *position.Span, op *token.Token, start, end ExpressionNode) *RangeLiteralNode {
	return &RangeLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Start:         start,
		End:           end,
		static:        areExpressionsStatic(start, end),
	}
}
