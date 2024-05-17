// Package ast defines typed AST nodes
// used by Elk
package ast

import (
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value/symbol"
)

type Node interface {
	ast.Node
	typ() types.Type
}

// Return the type of the given node.
func TypeOf(node Node, globalEnv *types.GlobalEnvironment) types.Type {
	switch node.(type) {
	case *FalseLiteralNode:
		return globalEnv.StdSubtype(symbol.False)
	case *TrueLiteralNode:
		return globalEnv.StdSubtype(symbol.True)
	case *NilLiteralNode:
		return globalEnv.StdSubtype(symbol.Nil)
	case *MethodDefinitionNode, *InitDefinitionNode:
		return globalEnv.StdSubtype(symbol.Method)
	case *MethodSignatureDefinitionNode:
		return types.Void{}
	case *InterpolatedSymbolLiteralNode:
		return globalEnv.StdSubtype(symbol.Symbol)
	case *InterpolatedStringLiteralNode:
		return globalEnv.StdSubtype(symbol.String)
	}
	return node.typ()
}

// Return the type of the given node.
func IsLiteral(node Node) bool {
	switch node.(type) {
	case *FalseLiteralNode, *TrueLiteralNode, *NilLiteralNode,
		*InterpolatedSymbolLiteralNode, *SimpleSymbolLiteralNode,
		*InterpolatedStringLiteralNode, *DoubleQuotedStringLiteralNode, *RawStringLiteralNode,
		*CharLiteralNode, *RawCharLiteralNode, *IntLiteralNode,
		*Int64LiteralNode, *Int32LiteralNode, *Int16LiteralNode, *Int8LiteralNode,
		*UInt64LiteralNode, *UInt32LiteralNode, *UInt16LiteralNode, *UInt8LiteralNode,
		*FloatLiteralNode, *Float64LiteralNode, *Float32LiteralNode, *BigFloatLiteralNode:
		return true
	}
	return false
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
func (*EmptyStatementNode) statementNode()      {}

// All expression nodes implement this interface.
type ExpressionNode interface {
	Node
	expressionNode()
}

func (*InvalidNode) expressionNode()        {}
func (*TypeExpressionNode) expressionNode() {}

// func (*VariablePatternDeclarationNode) expressionNode() {}
// func (*ValuePatternDeclarationNode) expressionNode()    {}
// func (*PostfixExpressionNode) expressionNode()          {}
// func (*ModifierNode) expressionNode()                   {}
// func (*ModifierIfElseNode) expressionNode()             {}
// func (*ModifierForInNode) expressionNode()              {}
// func (*AssignmentExpressionNode) expressionNode()       {}
// func (*BinaryExpressionNode) expressionNode()           {}
// func (*LogicalExpressionNode) expressionNode()          {}
// func (*UnaryExpressionNode) expressionNode()            {}
func (*TrueLiteralNode) expressionNode()  {}
func (*FalseLiteralNode) expressionNode() {}
func (*NilLiteralNode) expressionNode()   {}

// func (*InstanceVariableNode) expressionNode()           {}
func (*SimpleSymbolLiteralNode) expressionNode()       {}
func (*InterpolatedSymbolLiteralNode) expressionNode() {}
func (*IntLiteralNode) expressionNode()                {}
func (*Int64LiteralNode) expressionNode()              {}
func (*UInt64LiteralNode) expressionNode()             {}
func (*Int32LiteralNode) expressionNode()              {}
func (*UInt32LiteralNode) expressionNode()             {}
func (*Int16LiteralNode) expressionNode()              {}
func (*UInt16LiteralNode) expressionNode()             {}
func (*Int8LiteralNode) expressionNode()               {}
func (*UInt8LiteralNode) expressionNode()              {}
func (*FloatLiteralNode) expressionNode()              {}
func (*BigFloatLiteralNode) expressionNode()           {}
func (*Float64LiteralNode) expressionNode()            {}
func (*Float32LiteralNode) expressionNode()            {}

// func (*UninterpolatedRegexLiteralNode) expressionNode() {}
// func (*InterpolatedRegexLiteralNode) expressionNode()   {}
func (*RawStringLiteralNode) expressionNode()          {}
func (*CharLiteralNode) expressionNode()               {}
func (*RawCharLiteralNode) expressionNode()            {}
func (*DoubleQuotedStringLiteralNode) expressionNode() {}
func (*InterpolatedStringLiteralNode) expressionNode() {}
func (*VariableDeclarationNode) expressionNode()       {}
func (*ValueDeclarationNode) expressionNode()          {}
func (*PublicIdentifierNode) expressionNode()          {}
func (*PrivateIdentifierNode) expressionNode()         {}
func (*PublicConstantNode) expressionNode()            {}
func (*PrivateConstantNode) expressionNode()           {}

// func (*SelfLiteralNode) expressionNode()                {}
// func (*DoExpressionNode) expressionNode()               {}
// func (*SingletonBlockExpressionNode) expressionNode()   {}
// func (*SwitchExpressionNode) expressionNode()           {}
// func (*IfExpressionNode) expressionNode()               {}
// func (*UnlessExpressionNode) expressionNode()           {}
// func (*WhileExpressionNode) expressionNode()            {}
// func (*UntilExpressionNode) expressionNode()            {}
// func (*LoopExpressionNode) expressionNode()             {}
// func (*NumericForExpressionNode) expressionNode()       {}
// func (*ForInExpressionNode) expressionNode()            {}
// func (*BreakExpressionNode) expressionNode()            {}
// func (*LabeledExpressionNode) expressionNode()          {}
// func (*ReturnExpressionNode) expressionNode()           {}
// func (*ContinueExpressionNode) expressionNode()         {}
// func (*ThrowExpressionNode) expressionNode()            {}
func (*ConstantDeclarationNode) expressionNode() {}

// func (*FunctionLiteralNode) expressionNode()            {}
func (*ClassDeclarationNode) expressionNode()  {}
func (*ModuleDeclarationNode) expressionNode() {}
func (*MixinDeclarationNode) expressionNode()  {}

// func (*InterfaceDeclarationNode) expressionNode()       {}
// func (*StructDeclarationNode) expressionNode()          {}
func (*MethodDefinitionNode) expressionNode()          {}
func (*InitDefinitionNode) expressionNode()            {}
func (*MethodSignatureDefinitionNode) expressionNode() {}

// func (*GenericConstantNode) expressionNode()            {}
// func (*TypeDefinitionNode) expressionNode()             {}
// func (*AliasDeclarationNode) expressionNode()           {}
// func (*GetterDeclarationNode) expressionNode()          {}
// func (*SetterDeclarationNode) expressionNode()          {}
// func (*AccessorDeclarationNode) expressionNode()        {}
// func (*IncludeExpressionNode) expressionNode()          {}
// func (*ExtendExpressionNode) expressionNode()           {}
// func (*EnhanceExpressionNode) expressionNode()          {}
// func (*ConstructorCallNode) expressionNode()            {}
// func (*SubscriptExpressionNode) expressionNode()        {}
// func (*NilSafeSubscriptExpressionNode) expressionNode() {}
// func (*CallNode) expressionNode()                       {}
// func (*MethodCallNode) expressionNode()                 {}
// func (*ReceiverlessMethodCallNode) expressionNode()     {}
// func (*AttributeAccessNode) expressionNode()            {}
// func (*KeyValueExpressionNode) expressionNode()         {}
// func (*SymbolKeyValueExpressionNode) expressionNode()   {}
// func (*ArrayListLiteralNode) expressionNode()           {}
// func (*WordArrayListLiteralNode) expressionNode()       {}
// func (*WordArrayTupleLiteralNode) expressionNode()      {}
// func (*WordHashSetLiteralNode) expressionNode()         {}
// func (*SymbolArrayListLiteralNode) expressionNode()     {}
// func (*SymbolArrayTupleLiteralNode) expressionNode()    {}
// func (*SymbolHashSetLiteralNode) expressionNode()       {}
// func (*HexArrayListLiteralNode) expressionNode()        {}
// func (*HexArrayTupleLiteralNode) expressionNode()       {}
// func (*HexHashSetLiteralNode) expressionNode()          {}
// func (*BinArrayListLiteralNode) expressionNode()        {}
// func (*BinArrayTupleLiteralNode) expressionNode()       {}
// func (*BinHashSetLiteralNode) expressionNode()          {}
// func (*ArrayTupleLiteralNode) expressionNode()          {}
// func (*HashSetLiteralNode) expressionNode()             {}
// func (*HashMapLiteralNode) expressionNode()             {}
// func (*HashRecordLiteralNode) expressionNode()          {}
// func (*RangeLiteralNode) expressionNode()               {}
// func (*DocCommentNode) expressionNode()                 {}

// Represents a type variable in generics like `class Foo[+V]; end`
type TypeVariableNode interface {
	Node
	typeVariableNode()
}

func (*InvalidNode) typeVariableNode()             {}
func (*VariantTypeVariableNode) typeVariableNode() {}

// All nodes that should be valid in type annotations should
// implement this interface
type TypeNode interface {
	Node
	typeNode()
}

func (*InvalidNode) typeNode() {}

func (*UnionTypeNode) typeNode()        {}
func (*IntersectionTypeNode) typeNode() {}

// func (*NilableTypeNode) typeNode()          {}
// func (*SingletonTypeNode) typeNode()        {}
func (*PublicConstantNode) typeNode()            {}
func (*PrivateConstantNode) typeNode()           {}
func (*NilLiteralNode) typeNode()                {}
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

// func (*GenericConstantNode) typeNode()      {}

// All nodes that should be valid in constant lookups
// should implement this interface.
type ComplexConstantNode interface {
	Node
	TypeNode
	ExpressionNode
	// PatternNode
	// PatternExpressionNode
	complexConstantNode()
}

func (*InvalidNode) complexConstantNode()         {}
func (*PublicConstantNode) complexConstantNode()  {}
func (*PrivateConstantNode) complexConstantNode() {}

// func (*GenericConstantNode) complexConstantNode() {}

// All nodes that should be valid constants
// should implement this interface.
type ConstantNode interface {
	Node
	TypeNode
	ExpressionNode
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
	// PatternExpressionNode
	identifierNode()
}

func (*InvalidNode) identifierNode()           {}
func (*PublicIdentifierNode) identifierNode()  {}
func (*PrivateIdentifierNode) identifierNode() {}

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
	default:
		return false
	}
}

// Nodes that implement this interface represent
// symbol literals.
type SymbolLiteralNode interface {
	Node
	ExpressionNode
	symbolLiteralNode()
}

func (*InvalidNode) symbolLiteralNode()                   {}
func (*SimpleSymbolLiteralNode) symbolLiteralNode()       {}
func (*InterpolatedSymbolLiteralNode) symbolLiteralNode() {}

type StringOrSymbolLiteralNode interface {
	Node
	// PatternExpressionNode
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
	// PatternExpressionNode
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

// Represents an empty statement eg. a statement with only a semicolon or a newline.
type EmptyStatementNode struct {
	NodeBase
}

func (*EmptyStatementNode) IsStatic() bool {
	return false
}

// Create a new empty statement node.
func NewEmptyStatementNode(span *position.Span) *EmptyStatementNode {
	return &EmptyStatementNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a type expression `type String?`
type TypeExpressionNode struct {
	NodeBase
	TypeNode TypeNode
}

func (*TypeExpressionNode) IsStatic() bool {
	return false
}

// Create a new type expression `type String?`
func NewTypeExpressionNode(span *position.Span, typeNode TypeNode) *TypeExpressionNode {
	return &TypeExpressionNode{
		NodeBase: NodeBase{span: span},
		TypeNode: typeNode,
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
	_typ  types.Type
}

func (*IntLiteralNode) IsStatic() bool {
	return true
}

func (i *IntLiteralNode) typ() types.Type {
	return i._typ
}

// Create a new int literal node eg. `5`, `125_355`, `0xff`
func NewIntLiteralNode(span *position.Span, val string, typ types.Type) *IntLiteralNode {
	return &IntLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Int64 literal eg. `5i64`, `125_355i64`, `0xffi64`
type Int64LiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*Int64LiteralNode) IsStatic() bool {
	return true
}

func (i *Int64LiteralNode) typ() types.Type {
	return i._typ
}

// Create a new Int64 literal node eg. `5i64`, `125_355i64`, `0xffi64`
func NewInt64LiteralNode(span *position.Span, val string, typ types.Type) *Int64LiteralNode {
	return &Int64LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// UInt64 literal eg. `5u64`, `125_355u64`, `0xffu64`
type UInt64LiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*UInt64LiteralNode) IsStatic() bool {
	return true
}

func (i *UInt64LiteralNode) typ() types.Type {
	return i._typ
}

// Create a new UInt64 literal node eg. `5u64`, `125_355u64`, `0xffu64`
func NewUInt64LiteralNode(span *position.Span, val string, typ types.Type) *UInt64LiteralNode {
	return &UInt64LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Int32 literal eg. `5i32`, `1_20i32`, `0xffi32`
type Int32LiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*Int32LiteralNode) IsStatic() bool {
	return true
}

func (i *Int32LiteralNode) typ() types.Type {
	return i._typ
}

// Create a new Int32 literal node eg. `5i32`, `1_20i32`, `0xffi32`
func NewInt32LiteralNode(span *position.Span, val string, typ types.Type) *Int32LiteralNode {
	return &Int32LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// UInt32 literal eg. `5u32`, `1_20u32`, `0xffu32`
type UInt32LiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*UInt32LiteralNode) IsStatic() bool {
	return true
}

func (i *UInt32LiteralNode) typ() types.Type {
	return i._typ
}

// Create a new UInt32 literal node eg. `5u32`, `1_20u32`, `0xffu32`
func NewUInt32LiteralNode(span *position.Span, val string, typ types.Type) *UInt32LiteralNode {
	return &UInt32LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Int16 literal eg. `5i16`, `1_20i16`, `0xffi16`
type Int16LiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*Int16LiteralNode) IsStatic() bool {
	return true
}

func (i *Int16LiteralNode) typ() types.Type {
	return i._typ
}

// Create a new Int16 literal node eg. `5i16`, `1_20i16`, `0xffi16`
func NewInt16LiteralNode(span *position.Span, val string, typ types.Type) *Int16LiteralNode {
	return &Int16LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// UInt16 literal eg. `5u16`, `1_20u16`, `0xffu16`
type UInt16LiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*UInt16LiteralNode) IsStatic() bool {
	return true
}

func (i *UInt16LiteralNode) typ() types.Type {
	return i._typ
}

// Create a new UInt16 literal node eg. `5u16`, `1_20u16`, `0xffu16`
func NewUInt16LiteralNode(span *position.Span, val string, typ types.Type) *UInt16LiteralNode {
	return &UInt16LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Int8 literal eg. `5i8`, `1_20i8`, `0xffi8`
type Int8LiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*Int8LiteralNode) IsStatic() bool {
	return true
}

func (i *Int8LiteralNode) typ() types.Type {
	return i._typ
}

// Create a new Int8 literal node eg. `5i8`, `1_20i8`, `0xffi8`
func NewInt8LiteralNode(span *position.Span, val string, typ types.Type) *Int8LiteralNode {
	return &Int8LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// UInt8 literal eg. `5u8`, `1_20u8`, `0xffu8`
type UInt8LiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*UInt8LiteralNode) IsStatic() bool {
	return true
}

func (i *UInt8LiteralNode) typ() types.Type {
	return i._typ
}

// Create a new UInt8 literal node eg. `5u8`, `1_20u8`, `0xffu8`
func NewUInt8LiteralNode(span *position.Span, val string, typ types.Type) *UInt8LiteralNode {
	return &UInt8LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Float literal eg. `5.2`, `.5`, `45e20`
type FloatLiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*FloatLiteralNode) IsStatic() bool {
	return true
}

func (i *FloatLiteralNode) typ() types.Type {
	return i._typ
}

// Create a new float literal node eg. `5.2`, `.5`, `45e20`
func NewFloatLiteralNode(span *position.Span, val string, typ types.Type) *FloatLiteralNode {
	return &FloatLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// BigFloat literal eg. `5.2bf`, `.5bf`, `45e20bf`
type BigFloatLiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*BigFloatLiteralNode) IsStatic() bool {
	return true
}

func (i *BigFloatLiteralNode) typ() types.Type {
	return i._typ
}

// Create a new BigFloat literal node eg. `5.2bf`, `.5bf`, `45e20bf`
func NewBigFloatLiteralNode(span *position.Span, val string, typ types.Type) *BigFloatLiteralNode {
	return &BigFloatLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Float64 literal eg. `5.2f64`, `.5f64`, `45e20f64`
type Float64LiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*Float64LiteralNode) IsStatic() bool {
	return true
}

func (i *Float64LiteralNode) typ() types.Type {
	return i._typ
}

// Create a new Float64 literal node eg. `5.2f64`, `.5f64`, `45e20f64`
func NewFloat64LiteralNode(span *position.Span, val string, typ types.Type) *Float64LiteralNode {
	return &Float64LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Float32 literal eg. `5.2f32`, `.5f32`, `45e20f32`
type Float32LiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*Float32LiteralNode) IsStatic() bool {
	return true
}

func (i *Float32LiteralNode) typ() types.Type {
	return i._typ
}

// Create a new Float32 literal node eg. `5.2f32`, `.5f32`, `45e20f32`
func NewFloat32LiteralNode(span *position.Span, val string, typ types.Type) *Float32LiteralNode {
	return &Float32LiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Raw string literal enclosed with single quotes eg. `'foo'`.
type RawStringLiteralNode struct {
	NodeBase
	Value string // value of the string literal
	_typ  types.Type
}

func (*RawStringLiteralNode) IsStatic() bool {
	return true
}

func (r *RawStringLiteralNode) typ() types.Type {
	return r._typ
}

// Create a new raw string literal node eg. `'foo'`.
func NewRawStringLiteralNode(span *position.Span, val string, typ types.Type) *RawStringLiteralNode {
	return &RawStringLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Char literal eg. `c"a"`
type CharLiteralNode struct {
	NodeBase
	Value rune // value of the string literal
	_typ  types.Type
}

func (*CharLiteralNode) IsStatic() bool {
	return true
}

func (c *CharLiteralNode) typ() types.Type {
	return c._typ
}

// Create a new char literal node eg. `c"a"`
func NewCharLiteralNode(span *position.Span, val rune, typ types.Type) *CharLiteralNode {
	return &CharLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Raw Char literal eg. `a`
type RawCharLiteralNode struct {
	NodeBase
	Value rune // value of the char literal
	_typ  types.Type
}

func (*RawCharLiteralNode) IsStatic() bool {
	return true
}

func (r *RawCharLiteralNode) typ() types.Type {
	return r._typ
}

// Create a new raw char literal node eg. r`a`
func NewRawCharLiteralNode(span *position.Span, val rune, typ types.Type) *RawCharLiteralNode {
	return &RawCharLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Represents a single section of characters of a string literal eg. `foo` in `"foo${bar}"`.
type StringLiteralContentSectionNode struct {
	NodeBase
	Value string
}

func (*StringLiteralContentSectionNode) IsStatic() bool {
	return true
}

// Create a new string literal content section node eg. `foo` in `"foo${bar}"`.
func NewStringLiteralContentSectionNode(span *position.Span, val string) *StringLiteralContentSectionNode {
	return &StringLiteralContentSectionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a single inspect interpolated section of a string literal eg. `bar + 2` in `"foo#{bar + 2}"`
type StringInspectInterpolationNode struct {
	NodeBase
	Expression ExpressionNode
}

func (*StringInspectInterpolationNode) IsStatic() bool {
	return false
}

// Create a new string inspect interpolation node eg. `bar + 2` in `"foo#{bar + 2}"`
func NewStringInspectInterpolationNode(span *position.Span, expr ExpressionNode) *StringInspectInterpolationNode {
	return &StringInspectInterpolationNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Represents a single interpolated section of a string literal eg. `bar + 2` in `"foo${bar + 2}"`
type StringInterpolationNode struct {
	NodeBase
	Expression ExpressionNode
}

func (*StringInterpolationNode) IsStatic() bool {
	return false
}

// Create a new string interpolation node eg. `bar + 2` in `"foo${bar + 2}"`
func NewStringInterpolationNode(span *position.Span, expr ExpressionNode) *StringInterpolationNode {
	return &StringInterpolationNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Represents an interpolated string literal eg. `"foo ${bar} baz"`
type InterpolatedStringLiteralNode struct {
	NodeBase
	Content []StringLiteralContentNode
}

func (*InterpolatedStringLiteralNode) IsStatic() bool {
	return false
}

// Create a new interpolated string literal node eg. `"foo ${bar} baz"`
func NewInterpolatedStringLiteralNode(span *position.Span, cont []StringLiteralContentNode) *InterpolatedStringLiteralNode {
	return &InterpolatedStringLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  cont,
	}
}

// Represents a simple double quoted string literal eg. `"foo baz"`
type DoubleQuotedStringLiteralNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*DoubleQuotedStringLiteralNode) IsStatic() bool {
	return true
}

func (d *DoubleQuotedStringLiteralNode) typ() types.Type {
	return d._typ
}

// Create a new double quoted string literal node eg. `"foo baz"`
func NewDoubleQuotedStringLiteralNode(span *position.Span, val string, typ types.Type) *DoubleQuotedStringLiteralNode {
	return &DoubleQuotedStringLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Represents a symbol literal with simple content eg. `:foo`, `:'foo bar`, `:"lol"`
type SimpleSymbolLiteralNode struct {
	NodeBase
	Content string
	_typ    types.Type
}

func (*SimpleSymbolLiteralNode) IsStatic() bool {
	return true
}

func (s *SimpleSymbolLiteralNode) typ() types.Type {
	return s._typ
}

// Create a simple symbol literal node eg. `:foo`, `:'foo bar`, `:"lol"`
func NewSimpleSymbolLiteralNode(span *position.Span, cont string, typ types.Type) *SimpleSymbolLiteralNode {
	return &SimpleSymbolLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  cont,
		_typ:     typ,
	}
}

// Represents an interpolated symbol eg. `:"foo ${bar + 2}"`
type InterpolatedSymbolLiteralNode struct {
	NodeBase
	Content *InterpolatedStringLiteralNode
}

func (*InterpolatedSymbolLiteralNode) IsStatic() bool {
	return false
}

// Create an interpolated symbol literal node eg. `:"foo ${bar + 2}"`
func NewInterpolatedSymbolLiteralNode(span *position.Span, cont *InterpolatedStringLiteralNode) *InterpolatedSymbolLiteralNode {
	return &InterpolatedSymbolLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  cont,
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

// Represents a value declaration eg. `val foo: String`
type ValueDeclarationNode struct {
	NodeBase
	Name        *token.Token   // name of the value
	TypeNode    TypeNode       // type of the value
	Initialiser ExpressionNode // value assigned to the value
	_typ        types.Type
}

func (*ValueDeclarationNode) IsStatic() bool {
	return false
}

func (v *ValueDeclarationNode) typ() types.Type {
	return v._typ
}

// Create a new value declaration node eg. `val foo: String`
func NewValueDeclarationNode(span *position.Span, name *token.Token, typeNode TypeNode, init ExpressionNode, typ types.Type) *ValueDeclarationNode {
	return &ValueDeclarationNode{
		NodeBase:    NodeBase{span: span},
		Name:        name,
		TypeNode:    typeNode,
		Initialiser: init,
		_typ:        typ,
	}
}

// Represents a constant declaration eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
type ConstantDeclarationNode struct {
	NodeBase
	Name        *token.Token   // name of the constant
	TypeNode    TypeNode       // type of the constant
	Initialiser ExpressionNode // value assigned to the constant
	_typ        types.Type
}

func (*ConstantDeclarationNode) IsStatic() bool {
	return false
}

func (c *ConstantDeclarationNode) typ() types.Type {
	return c._typ
}

// Create a new constant declaration node eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
func NewConstantDeclarationNode(span *position.Span, name *token.Token, typeNode TypeNode, init ExpressionNode, typ types.Type) *ConstantDeclarationNode {
	return &ConstantDeclarationNode{
		NodeBase:    NodeBase{span: span},
		Name:        name,
		TypeNode:    typeNode,
		Initialiser: init,
		_typ:        typ,
	}
}

// Union type eg. `String & Int & Float`
type IntersectionTypeNode struct {
	NodeBase
	Elements []TypeNode
	_typ     types.Type
}

func (*IntersectionTypeNode) IsStatic() bool {
	return false
}

func (b *IntersectionTypeNode) typ() types.Type {
	return b._typ
}

// Create a new binary type expression node eg. `String & Int`
func NewIntersectionTypeNode(span *position.Span, elements []TypeNode, typ types.Type) *IntersectionTypeNode {
	return &IntersectionTypeNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
		_typ:     typ,
	}
}

// Union type eg. `String | Int | Float`
type UnionTypeNode struct {
	NodeBase
	Elements []TypeNode
	_typ     types.Type
}

func (*UnionTypeNode) IsStatic() bool {
	return false
}

func (b *UnionTypeNode) typ() types.Type {
	return b._typ
}

// Create a new binary type expression node eg. `String | Int`
func NewUnionTypeNode(span *position.Span, elements []TypeNode, typ types.Type) *UnionTypeNode {
	return &UnionTypeNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
		_typ:     typ,
	}
}

// Represents a public identifier eg. `foo`.
type PublicIdentifierNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*PublicIdentifierNode) IsStatic() bool {
	return false
}

func (p *PublicIdentifierNode) typ() types.Type {
	return p._typ
}

// Create a new public identifier node eg. `foo`.
func NewPublicIdentifierNode(span *position.Span, val string, typ types.Type) *PublicIdentifierNode {
	return &PublicIdentifierNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Represents a private identifier eg. `_foo`
type PrivateIdentifierNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*PrivateIdentifierNode) IsStatic() bool {
	return false
}

func (p *PrivateIdentifierNode) typ() types.Type {
	return p._typ
}

// Create a new private identifier node eg. `_foo`.
func NewPrivateIdentifierNode(span *position.Span, val string, typ types.Type) *PrivateIdentifierNode {
	return &PrivateIdentifierNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Represents a public constant eg. `Foo`.
type PublicConstantNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*PublicConstantNode) IsStatic() bool {
	return false
}

func (p *PublicConstantNode) typ() types.Type {
	return p._typ
}

// Create a new public constant node eg. `Foo`.
func NewPublicConstantNode(span *position.Span, val string, typ types.Type) *PublicConstantNode {
	return &PublicConstantNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
	}
}

// Represents a private constant eg. `_Foo`
type PrivateConstantNode struct {
	NodeBase
	Value string
	_typ  types.Type
}

func (*PrivateConstantNode) IsStatic() bool {
	return false
}

func (p *PrivateConstantNode) typ() types.Type {
	return p._typ
}

// Create a new private constant node eg. `_Foo`.
func NewPrivateConstantNode(span *position.Span, val string, typ types.Type) *PrivateConstantNode {
	return &PrivateConstantNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		_typ:     typ,
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
	NodeBase
	Name        string         // name of the variable
	Type        TypeNode       // type of the variable
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
		NodeBase:    NodeBase{span: span},
		Name:        name,
		Type:        typ,
		Initialiser: init,
		Kind:        kind,
	}
}

// Represents a formal parameter in method declarations eg. `foo: String = 'bar'`
type MethodParameterNode struct {
	NodeBase
	Name                string         // name of the variable
	Type                TypeNode       // type of the variable
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
		NodeBase:            NodeBase{span: span},
		SetInstanceVariable: setIvar,
		Name:                name,
		Type:                typ,
		Initialiser:         init,
		Kind:                kind,
	}
}

// Represents a signature parameter in method and function signatures eg. `foo?: String`
type SignatureParameterNode struct {
	NodeBase
	Name     string   // name of the variable
	Type     TypeNode // type of the variable
	Optional bool     // whether this parameter is optional
}

func (*SignatureParameterNode) IsStatic() bool {
	return false
}

func (f *SignatureParameterNode) IsOptional() bool {
	return f.Optional
}

// Create a new signature parameter node eg. `foo?: String`
func NewSignatureParameterNode(span *position.Span, name string, typ TypeNode, opt bool) *SignatureParameterNode {
	return &SignatureParameterNode{
		NodeBase: NodeBase{span: span},
		Name:     name,
		Type:     typ,
		Optional: opt,
	}
}

// Represents an attribute declaration in getters, setters and accessors eg. `foo: String`
type AttributeParameterNode struct {
	NodeBase
	Name string   // name of the variable
	Type TypeNode // type of the variable
}

func (*AttributeParameterNode) IsStatic() bool {
	return false
}

func (a *AttributeParameterNode) IsOptional() bool {
	return false
}

// Create a new attribute declaration in getters, setters and accessors eg. `foo: String`
func NewAttributeParameterNode(span *position.Span, name string, typ TypeNode) *AttributeParameterNode {
	return &AttributeParameterNode{
		NodeBase: NodeBase{span: span},
		Name:     name,
		Type:     typ,
	}
}

// Represents the variance of a type variable.
type Variance uint8

const (
	INVARIANT Variance = iota
	COVARIANT
	CONTRAVARIANT
)

// Represents a type variable eg. `+V`
type VariantTypeVariableNode struct {
	NodeBase
	Variance   Variance // Variance level of this type variable
	Name       string   // Name of the type variable eg. `T`
	UpperBound ComplexConstantNode
}

func (*VariantTypeVariableNode) IsStatic() bool {
	return false
}

// Create a new type variable node eg. `+V`
func NewVariantTypeVariableNode(span *position.Span, variance Variance, name string, upper ComplexConstantNode) *VariantTypeVariableNode {
	return &VariantTypeVariableNode{
		NodeBase:   NodeBase{span: span},
		Variance:   variance,
		Name:       name,
		UpperBound: upper,
	}
}

// Represents a class declaration eg. `class Foo; end`
type ClassDeclarationNode struct {
	NodeBase
	Abstract      bool
	Sealed        bool
	Constant      ExpressionNode     // The constant that will hold the class value
	TypeVariables []TypeVariableNode // Generic type variable definitions
	Superclass    ExpressionNode     // the super/parent class of this class
	Body          []StatementNode    // body of the class
	_typ          types.Type
}

func (*ClassDeclarationNode) IsStatic() bool {
	return false
}

func (c *ClassDeclarationNode) typ() types.Type {
	return c._typ
}

// Create a new class declaration node eg. `class Foo; end`
func NewClassDeclarationNode(
	span *position.Span,
	abstract bool,
	sealed bool,
	constant ExpressionNode,
	typeVars []TypeVariableNode,
	superclass ExpressionNode,
	body []StatementNode,
	typ types.Type,
) *ClassDeclarationNode {

	return &ClassDeclarationNode{
		NodeBase:      NodeBase{span: span},
		Abstract:      abstract,
		Sealed:        sealed,
		Constant:      constant,
		TypeVariables: typeVars,
		Superclass:    superclass,
		Body:          body,
		_typ:          typ,
	}
}

// Represents a module declaration eg. `module Foo; end`
type ModuleDeclarationNode struct {
	NodeBase
	Constant ExpressionNode  // The constant that will hold the module value
	Body     []StatementNode // body of the module
	_typ     types.Type
}

func (*ModuleDeclarationNode) IsStatic() bool {
	return false
}

func (m *ModuleDeclarationNode) typ() types.Type {
	return m._typ
}

// Create a new module declaration node eg. `module Foo; end`
func NewModuleDeclarationNode(
	span *position.Span,
	constant ExpressionNode,
	body []StatementNode,
	typ types.Type,
) *ModuleDeclarationNode {

	return &ModuleDeclarationNode{
		NodeBase: NodeBase{span: span},
		Constant: constant,
		Body:     body,
		_typ:     typ,
	}
}

// Represents a mixin declaration eg. `mixin Foo; end`
type MixinDeclarationNode struct {
	NodeBase
	Constant      ExpressionNode     // The constant that will hold the mixin value
	TypeVariables []TypeVariableNode // Generic type variable definitions
	Body          []StatementNode    // body of the mixin
	_typ          types.Type
}

func (*MixinDeclarationNode) IsStatic() bool {
	return false
}

func (m *MixinDeclarationNode) typ() types.Type {
	return m._typ
}

// Create a new mixin declaration node eg. `mixin Foo; end`
func NewMixinDeclarationNode(
	span *position.Span,
	constant ExpressionNode,
	typeVars []TypeVariableNode,
	body []StatementNode,
	typ types.Type,
) *MixinDeclarationNode {

	return &MixinDeclarationNode{
		NodeBase:      NodeBase{span: span},
		Constant:      constant,
		TypeVariables: typeVars,
		Body:          body,
		_typ:          typ,
	}
}

// Represents a method definition eg. `def foo: String then 'hello world'`
type MethodDefinitionNode struct {
	NodeBase
	Name       string
	Parameters []ParameterNode // formal parameters
	ReturnType TypeNode
	ThrowType  TypeNode
	Body       []StatementNode // body of the method
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

// Create a method definition node eg. `def foo: String then 'hello world'`
func NewMethodDefinitionNode(span *position.Span, name string, params []ParameterNode, returnType, throwType TypeNode, body []StatementNode) *MethodDefinitionNode {
	return &MethodDefinitionNode{
		NodeBase:   NodeBase{span: span},
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		ThrowType:  throwType,
		Body:       body,
	}
}

// Represents a constructor definition eg. `init then 'hello world'`
type InitDefinitionNode struct {
	NodeBase
	Parameters []ParameterNode // formal parameters
	ThrowType  TypeNode
	Body       []StatementNode // body of the method
}

func (*InitDefinitionNode) IsStatic() bool {
	return false
}

// Create a constructor definition node eg. `init then 'hello world'`
func NewInitDefinitionNode(span *position.Span, params []ParameterNode, throwType TypeNode, body []StatementNode) *InitDefinitionNode {
	return &InitDefinitionNode{
		NodeBase:   NodeBase{span: span},
		Parameters: params,
		ThrowType:  throwType,
		Body:       body,
	}
}

// Represents a method signature definition eg. `sig to_string(val: Int): String`
type MethodSignatureDefinitionNode struct {
	NodeBase
	Name       string
	Parameters []ParameterNode // formal parameters
	ReturnType TypeNode
	ThrowType  TypeNode
}

func (*MethodSignatureDefinitionNode) IsStatic() bool {
	return false
}

// Create a method signature node eg. `sig to_string(val: Int): String`
func NewMethodSignatureDefinitionNode(span *position.Span, name string, params []ParameterNode, returnType, throwType TypeNode) *MethodSignatureDefinitionNode {
	return &MethodSignatureDefinitionNode{
		NodeBase:   NodeBase{span: span},
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		ThrowType:  throwType,
	}
}
