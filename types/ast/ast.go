// Package ast defines typed AST nodes
// used by Elk
package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
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
		return globalEnv.StdSubtype("False")
	case *TrueLiteralNode:
		return globalEnv.StdSubtype("True")
	case *NilLiteralNode:
		return globalEnv.StdSubtype("Nil")
	case *DoubleQuotedStringLiteralNode:
		return globalEnv.StdSubtype("String")
	case *IntLiteralNode:
		return globalEnv.StdSubtype("Int")
	case *Int64LiteralNode:
		return globalEnv.StdSubtype("Int64")
	case *Int32LiteralNode:
		return globalEnv.StdSubtype("Int32")
	case *Int16LiteralNode:
		return globalEnv.StdSubtype("Int16")
	case *Int8LiteralNode:
		return globalEnv.StdSubtype("Int8")
	case *UInt64LiteralNode:
		return globalEnv.StdSubtype("UInt64")
	case *UInt32LiteralNode:
		return globalEnv.StdSubtype("UInt32")
	case *UInt16LiteralNode:
		return globalEnv.StdSubtype("UInt16")
	case *UInt8LiteralNode:
		return globalEnv.StdSubtype("UInt8")
	case *FloatLiteralNode:
		return globalEnv.StdSubtype("Float")
	case *Float64LiteralNode:
		return globalEnv.StdSubtype("Float64")
	case *Float32LiteralNode:
		return globalEnv.StdSubtype("Float32")
	case *BigFloatLiteralNode:
		return globalEnv.StdSubtype("BigFloat")
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

// func (*EmptyStatementNode) statementNode()      {}

// All expression nodes implement this interface.
type ExpressionNode interface {
	Node
	expressionNode()
}

func (*InvalidNode) expressionNode() {}

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
// func (*SimpleSymbolLiteralNode) expressionNode()        {}
// func (*InterpolatedSymbolLiteralNode) expressionNode()  {}
func (*IntLiteralNode) expressionNode()      {}
func (*Int64LiteralNode) expressionNode()    {}
func (*UInt64LiteralNode) expressionNode()   {}
func (*Int32LiteralNode) expressionNode()    {}
func (*UInt32LiteralNode) expressionNode()   {}
func (*Int16LiteralNode) expressionNode()    {}
func (*UInt16LiteralNode) expressionNode()   {}
func (*Int8LiteralNode) expressionNode()     {}
func (*UInt8LiteralNode) expressionNode()    {}
func (*FloatLiteralNode) expressionNode()    {}
func (*BigFloatLiteralNode) expressionNode() {}
func (*Float64LiteralNode) expressionNode()  {}
func (*Float32LiteralNode) expressionNode()  {}

// func (*UninterpolatedRegexLiteralNode) expressionNode() {}
// func (*InterpolatedRegexLiteralNode) expressionNode()   {}
// func (*RawStringLiteralNode) expressionNode()           {}
// func (*CharLiteralNode) expressionNode()                {}
// func (*RawCharLiteralNode) expressionNode()             {}
func (*DoubleQuotedStringLiteralNode) expressionNode() {}

// func (*InterpolatedStringLiteralNode) expressionNode()  {}
func (*VariableDeclarationNode) expressionNode() {}

func (*ValueDeclarationNode) expressionNode() {}

func (*PublicIdentifierNode) expressionNode()  {}
func (*PrivateIdentifierNode) expressionNode() {}
func (*PublicConstantNode) expressionNode()    {}
func (*PrivateConstantNode) expressionNode()   {}

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
// func (*ConstantDeclarationNode) expressionNode()        {}

// func (*FunctionLiteralNode) expressionNode()            {}
func (*ClassDeclarationNode) expressionNode()  {}
func (*ModuleDeclarationNode) expressionNode() {}
func (*MixinDeclarationNode) expressionNode()  {}

// func (*InterfaceDeclarationNode) expressionNode()       {}
// func (*StructDeclarationNode) expressionNode()          {}
// func (*MethodDefinitionNode) expressionNode()           {}
// func (*InitDefinitionNode) expressionNode()             {}
// func (*MethodSignatureDefinitionNode) expressionNode()  {}
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

// func (*BinaryTypeExpressionNode) typeNode() {}
// func (*NilableTypeNode) typeNode()          {}
// func (*SingletonTypeNode) typeNode()        {}
func (*PublicConstantNode) typeNode()  {}
func (*PrivateConstantNode) typeNode() {}

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
