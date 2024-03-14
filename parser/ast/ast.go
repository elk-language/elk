// Package ast defines types
// used by the Elk parser.
//
// All the nodes of the Abstract Syntax Tree
// constructed by the Elk parser are defined in this package.
package ast

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/token"
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
	IsStatic() bool // Value is known at compile-time
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

// Check whether the token can be used as a left value
// in a variable/constant declaration.
func IsValidDeclarationTarget(node Node) bool {
	switch node.(type) {
	case *PrivateConstantNode, *PublicConstantNode,
		*ConstantLookupNode, *PrivateIdentifierNode, *PublicIdentifierNode:
		return true
	default:
		return false
	}
}

// Check whether the token can be used as a left value
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

// All expression nodes implement this interface.
type ExpressionNode interface {
	Node
	expressionNode()
}

func (*InvalidNode) expressionNode()                    {}
func (*PostfixExpressionNode) expressionNode()          {}
func (*ModifierNode) expressionNode()                   {}
func (*ModifierIfElseNode) expressionNode()             {}
func (*ModifierForInNode) expressionNode()              {}
func (*AssignmentExpressionNode) expressionNode()       {}
func (*BinaryExpressionNode) expressionNode()           {}
func (*LogicalExpressionNode) expressionNode()          {}
func (*UnaryExpressionNode) expressionNode()            {}
func (*TrueLiteralNode) expressionNode()                {}
func (*FalseLiteralNode) expressionNode()               {}
func (*NilLiteralNode) expressionNode()                 {}
func (*InstanceVariableNode) expressionNode()           {}
func (*SimpleSymbolLiteralNode) expressionNode()        {}
func (*InterpolatedSymbolLiteral) expressionNode()      {}
func (*IntLiteralNode) expressionNode()                 {}
func (*Int64LiteralNode) expressionNode()               {}
func (*UInt64LiteralNode) expressionNode()              {}
func (*Int32LiteralNode) expressionNode()               {}
func (*UInt32LiteralNode) expressionNode()              {}
func (*Int16LiteralNode) expressionNode()               {}
func (*UInt16LiteralNode) expressionNode()              {}
func (*Int8LiteralNode) expressionNode()                {}
func (*UInt8LiteralNode) expressionNode()               {}
func (*FloatLiteralNode) expressionNode()               {}
func (*BigFloatLiteralNode) expressionNode()            {}
func (*Float64LiteralNode) expressionNode()             {}
func (*Float32LiteralNode) expressionNode()             {}
func (*UninterpolatedRegexLiteralNode) expressionNode() {}
func (*InterpolatedRegexLiteralNode) expressionNode()   {}
func (*RawStringLiteralNode) expressionNode()           {}
func (*CharLiteralNode) expressionNode()                {}
func (*RawCharLiteralNode) expressionNode()             {}
func (*DoubleQuotedStringLiteralNode) expressionNode()  {}
func (*InterpolatedStringLiteralNode) expressionNode()  {}
func (*VariableDeclarationNode) expressionNode()        {}
func (*ValueDeclarationNode) expressionNode()           {}
func (*PublicIdentifierNode) expressionNode()           {}
func (*PrivateIdentifierNode) expressionNode()          {}
func (*PublicConstantNode) expressionNode()             {}
func (*PrivateConstantNode) expressionNode()            {}
func (*SelfLiteralNode) expressionNode()                {}
func (*DoExpressionNode) expressionNode()               {}
func (*SingletonBlockExpressionNode) expressionNode()   {}
func (*IfExpressionNode) expressionNode()               {}
func (*UnlessExpressionNode) expressionNode()           {}
func (*WhileExpressionNode) expressionNode()            {}
func (*UntilExpressionNode) expressionNode()            {}
func (*LoopExpressionNode) expressionNode()             {}
func (*NumericForExpressionNode) expressionNode()       {}
func (*ForInExpressionNode) expressionNode()            {}
func (*BreakExpressionNode) expressionNode()            {}
func (*LabeledExpressionNode) expressionNode()          {}
func (*ReturnExpressionNode) expressionNode()           {}
func (*ContinueExpressionNode) expressionNode()         {}
func (*ThrowExpressionNode) expressionNode()            {}
func (*ConstantDeclarationNode) expressionNode()        {}
func (*ConstantLookupNode) expressionNode()             {}
func (*ClosureLiteralNode) expressionNode()             {}
func (*ClassDeclarationNode) expressionNode()           {}
func (*ModuleDeclarationNode) expressionNode()          {}
func (*MixinDeclarationNode) expressionNode()           {}
func (*InterfaceDeclarationNode) expressionNode()       {}
func (*StructDeclarationNode) expressionNode()          {}
func (*MethodDefinitionNode) expressionNode()           {}
func (*InitDefinitionNode) expressionNode()             {}
func (*MethodSignatureDefinitionNode) expressionNode()  {}
func (*GenericConstantNode) expressionNode()            {}
func (*TypeDefinitionNode) expressionNode()             {}
func (*AliasDeclarationNode) expressionNode()           {}
func (*GetterDeclarationNode) expressionNode()          {}
func (*SetterDeclarationNode) expressionNode()          {}
func (*AccessorDeclarationNode) expressionNode()        {}
func (*IncludeExpressionNode) expressionNode()          {}
func (*ExtendExpressionNode) expressionNode()           {}
func (*EnhanceExpressionNode) expressionNode()          {}
func (*ConstructorCallNode) expressionNode()            {}
func (*SubscriptExpressionNode) expressionNode()        {}
func (*NilSafeSubscriptExpressionNode) expressionNode() {}
func (*MethodCallNode) expressionNode()                 {}
func (*FunctionCallNode) expressionNode()               {}
func (*AttributeAccessNode) expressionNode()            {}
func (*KeyValueExpressionNode) expressionNode()         {}
func (*SymbolKeyValueExpressionNode) expressionNode()   {}
func (*ArrayListLiteralNode) expressionNode()           {}
func (*WordArrayListLiteralNode) expressionNode()       {}
func (*WordArrayTupleLiteralNode) expressionNode()      {}
func (*WordHashSetLiteralNode) expressionNode()         {}
func (*SymbolArrayListLiteralNode) expressionNode()     {}
func (*SymbolArrayTupleLiteralNode) expressionNode()    {}
func (*SymbolHashSetLiteralNode) expressionNode()       {}
func (*HexArrayListLiteralNode) expressionNode()        {}
func (*HexArrayTupleLiteralNode) expressionNode()       {}
func (*HexHashSetLiteralNode) expressionNode()          {}
func (*BinArrayListLiteralNode) expressionNode()        {}
func (*BinArrayTupleLiteralNode) expressionNode()       {}
func (*BinHashSetLiteralNode) expressionNode()          {}
func (*ArrayTupleLiteralNode) expressionNode()          {}
func (*HashSetLiteralNode) expressionNode()             {}
func (*HashMapLiteralNode) expressionNode()             {}
func (*HashRecordLiteralNode) expressionNode()          {}
func (*RangeLiteralNode) expressionNode()               {}
func (*ArithmeticSequenceLiteralNode) expressionNode()  {}
func (*DocCommentNode) expressionNode()                 {}

// All nodes that should be valid in type annotations should
// implement this interface
type TypeNode interface {
	Node
	typeNode()
}

func (*InvalidNode) typeNode()              {}
func (*BinaryTypeExpressionNode) typeNode() {}
func (*NilableTypeNode) typeNode()          {}
func (*SingletonTypeNode) typeNode()        {}
func (*PublicConstantNode) typeNode()       {}
func (*PrivateConstantNode) typeNode()      {}
func (*ConstantLookupNode) typeNode()       {}
func (*GenericConstantNode) typeNode()      {}

// All nodes that represent regexes should
// implement this interface.
type RegexLiteralNode interface {
	Node
	ExpressionNode
	regexLiteralNode()
}

func (*InvalidNode) regexLiteralNode()                    {}
func (*UninterpolatedRegexLiteralNode) regexLiteralNode() {}
func (*InterpolatedRegexLiteralNode) regexLiteralNode()   {}

// All nodes that represent strings should
// implement this interface.
type StringLiteralNode interface {
	Node
	ExpressionNode
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
	simpleStringLiteralNode()
}

func (*InvalidNode) simpleStringLiteralNode()                   {}
func (*DoubleQuotedStringLiteralNode) simpleStringLiteralNode() {}
func (*RawStringLiteralNode) simpleStringLiteralNode()          {}

// All nodes that should be valid in parameter declaration lists
// of methods or closures should implement this interface.
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

// Represents a type variable in generics like `class Foo[+V]; end`
type TypeVariableNode interface {
	Node
	typeVariableNode()
}

func (*InvalidNode) typeVariableNode()             {}
func (*VariantTypeVariableNode) typeVariableNode() {}

// All nodes that should be valid in constant lookups
// should implement this interface.
type ComplexConstantNode interface {
	Node
	TypeNode
	ExpressionNode
	complexConstantNode()
}

func (*InvalidNode) complexConstantNode()         {}
func (*PublicConstantNode) complexConstantNode()  {}
func (*PrivateConstantNode) complexConstantNode() {}
func (*ConstantLookupNode) complexConstantNode()  {}
func (*GenericConstantNode) complexConstantNode() {}

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
	ExpressionNode
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
// symbol literals.
type SymbolLiteralNode interface {
	Node
	ExpressionNode
	symbolLiteralNode()
}

func (*InvalidNode) symbolLiteralNode()               {}
func (*SimpleSymbolLiteralNode) symbolLiteralNode()   {}
func (*InterpolatedSymbolLiteral) symbolLiteralNode() {}

// Nodes that implement this interface represent
// named arguments in method calls.
type NamedArgumentNode interface {
	Node
	namedArgumentNode()
}

func (*InvalidNode) namedArgumentNode()           {}
func (*NamedCallArgumentNode) namedArgumentNode() {}

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

// Same as [NewExpressionStatementNode] but returns an interface
func NewExpressionStatementNodeI(span *position.Span, expr ExpressionNode) StatementNode {
	return &ExpressionStatementNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

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

// Represents a variable declaration eg. `var foo: String`
type VariableDeclarationNode struct {
	NodeBase
	Name        *token.Token   // name of the variable
	Type        TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

func (*VariableDeclarationNode) IsStatic() bool {
	return false
}

// Create a new variable declaration node eg. `var foo: String`
func NewVariableDeclarationNode(span *position.Span, name *token.Token, typ TypeNode, init ExpressionNode) *VariableDeclarationNode {
	return &VariableDeclarationNode{
		NodeBase:    NodeBase{span: span},
		Name:        name,
		Type:        typ,
		Initialiser: init,
	}
}

// Represents a value declaration eg. `val foo: String`
type ValueDeclarationNode struct {
	NodeBase
	Name        *token.Token   // name of the value
	Type        TypeNode       // type of the value
	Initialiser ExpressionNode // value assigned to the value
}

func (*ValueDeclarationNode) IsStatic() bool {
	return false
}

// Create a new value declaration node eg. `val foo: String`
func NewValueDeclarationNode(span *position.Span, name *token.Token, typ TypeNode, init ExpressionNode) *ValueDeclarationNode {
	return &ValueDeclarationNode{
		NodeBase:    NodeBase{span: span},
		Name:        name,
		Type:        typ,
		Initialiser: init,
	}
}

// Assignment with the specified operator.
type AssignmentExpressionNode struct {
	NodeBase
	Op    *token.Token   // operator
	Left  ExpressionNode // left hand side
	Right ExpressionNode // right hand side
}

func (*AssignmentExpressionNode) IsStatic() bool {
	return false
}

// Create a new assignment expression node eg. `foo = 3`
func NewAssignmentExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *AssignmentExpressionNode {
	return &AssignmentExpressionNode{
		NodeBase: NodeBase{span: span},
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Expression of an operator with two operands eg. `2 + 5`, `foo > bar`
type BinaryExpressionNode struct {
	NodeBase
	Op     *token.Token   // operator
	Left   ExpressionNode // left hand side
	Right  ExpressionNode // right hand side
	static bool
}

func (b *BinaryExpressionNode) IsStatic() bool {
	return b.static
}

// Create a new binary expression node.
func NewBinaryExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *BinaryExpressionNode {
	return &BinaryExpressionNode{
		NodeBase: NodeBase{span: span},
		Op:       op,
		Left:     left,
		Right:    right,
		static:   areExpressionsStatic(left, right),
	}
}

// Same as [NewBinaryExpressionNode] but returns an interface
func NewBinaryExpressionNodeI(span *position.Span, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return NewBinaryExpressionNode(span, op, left, right)
}

// Expression of a logical operator with two operands eg. `foo && bar`
type LogicalExpressionNode struct {
	NodeBase
	Op     *token.Token   // operator
	Left   ExpressionNode // left hand side
	Right  ExpressionNode // right hand side
	static bool
}

func (l *LogicalExpressionNode) IsStatic() bool {
	return l.static
}

// Create a new logical expression node.
func NewLogicalExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *LogicalExpressionNode {
	return &LogicalExpressionNode{
		NodeBase: NodeBase{span: span},
		Op:       op,
		Left:     left,
		Right:    right,
		static:   areExpressionsStatic(left, right),
	}
}

// Same as [NewLogicalExpressionNode] but returns an interface
func NewLogicalExpressionNodeI(span *position.Span, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return &LogicalExpressionNode{
		NodeBase: NodeBase{span: span},
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Expression of an operator with one operand eg. `!foo`, `-bar`
type UnaryExpressionNode struct {
	NodeBase
	Op    *token.Token   // operator
	Right ExpressionNode // right hand side
}

func (u *UnaryExpressionNode) IsStatic() bool {
	return u.Right.IsStatic()
}

// Create a new unary expression node.
func NewUnaryExpressionNode(span *position.Span, op *token.Token, right ExpressionNode) *UnaryExpressionNode {
	return &UnaryExpressionNode{
		NodeBase: NodeBase{span: span},
		Op:       op,
		Right:    right,
	}
}

// Postfix expression eg. `foo++`, `bar--`
type PostfixExpressionNode struct {
	NodeBase
	Op         *token.Token // operator
	Expression ExpressionNode
}

func (i *PostfixExpressionNode) IsStatic() bool {
	return false
}

// Create a new postfix expression node.
func NewPostfixExpressionNode(span *position.Span, op *token.Token, expr ExpressionNode) *PostfixExpressionNode {
	return &PostfixExpressionNode{
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

// `self` literal.
type SelfLiteralNode struct {
	NodeBase
}

func (*SelfLiteralNode) IsStatic() bool {
	return false
}

// Create a new `self` literal node.
func NewSelfLiteralNode(span *position.Span) *SelfLiteralNode {
	return &SelfLiteralNode{
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

// Raw string literal enclosed with single quotes eg. `'foo'`.
type RawStringLiteralNode struct {
	NodeBase
	Value string // value of the string literal
}

func (*RawStringLiteralNode) IsStatic() bool {
	return true
}

// Create a new raw string literal node eg. `'foo'`.
func NewRawStringLiteralNode(span *position.Span, val string) *RawStringLiteralNode {
	return &RawStringLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Char literal eg. `c"a"`
type CharLiteralNode struct {
	NodeBase
	Value rune // value of the string literal
}

func (*CharLiteralNode) IsStatic() bool {
	return true
}

// Create a new char literal node eg. `c"a"`
func NewCharLiteralNode(span *position.Span, val rune) *CharLiteralNode {
	return &CharLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Raw Char literal eg. `a`
type RawCharLiteralNode struct {
	NodeBase
	Value rune // value of the char literal
}

func (*RawCharLiteralNode) IsStatic() bool {
	return true
}

// Create a new raw char literal node eg. r`a`
func NewRawCharLiteralNode(span *position.Span, val rune) *RawCharLiteralNode {
	return &RawCharLiteralNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
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

// Represents a public identifier eg. `foo`.
type PublicIdentifierNode struct {
	NodeBase
	Value string
}

func (*PublicIdentifierNode) IsStatic() bool {
	return false
}

// Create a new public identifier node eg. `foo`.
func NewPublicIdentifierNode(span *position.Span, val string) *PublicIdentifierNode {
	return &PublicIdentifierNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents an uninterpolated regex literal eg. `%/foo/`
type UninterpolatedRegexLiteralNode struct {
	NodeBase
	Content string
	Flags   bitfield.BitField8
}

func (*UninterpolatedRegexLiteralNode) IsStatic() bool {
	return true
}

func (r *UninterpolatedRegexLiteralNode) IsCaseInsensitive() bool {
	return r.Flags.HasFlag(flag.CaseInsensitiveFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetCaseInsensitive() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.CaseInsensitiveFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsMultiline() bool {
	return r.Flags.HasFlag(flag.MultilineFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetMultiline() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.MultilineFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsDotAll() bool {
	return r.Flags.HasFlag(flag.DotAllFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetDotAll() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.DotAllFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsUngreedy() bool {
	return r.Flags.HasFlag(flag.UngreedyFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetUngreedy() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.UngreedyFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsASCII() bool {
	return r.Flags.HasFlag(flag.ASCIIFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetASCII() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ASCIIFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsExtended() bool {
	return r.Flags.HasFlag(flag.ExtendedFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetExtended() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ExtendedFlag)
	return r
}

// Create a new uninterpolated regex literal node eg. `%/foo/`.
func NewUninterpolatedRegexLiteralNode(span *position.Span, content string, flags bitfield.BitField8) *UninterpolatedRegexLiteralNode {
	return &UninterpolatedRegexLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  content,
		Flags:    flags,
	}
}

// Represents a single section of characters of a regex literal eg. `foo` in `%/foo${bar}/`.
type RegexLiteralContentSectionNode struct {
	NodeBase
	Value string
}

func (*RegexLiteralContentSectionNode) IsStatic() bool {
	return true
}

// Create a new regex literal content section node eg. `foo` in `%/foo${bar}/`.
func NewRegexLiteralContentSectionNode(span *position.Span, val string) *RegexLiteralContentSectionNode {
	return &RegexLiteralContentSectionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a single interpolated section of a regex literal eg. `bar + 2` in `%/foo${bar + 2}/`
type RegexInterpolationNode struct {
	NodeBase
	Expression ExpressionNode
}

func (*RegexInterpolationNode) IsStatic() bool {
	return false
}

// Create a new regex interpolation node eg. `bar + 2` in `%/foo${bar + 2}/`
func NewRegexInterpolationNode(span *position.Span, expr ExpressionNode) *RegexInterpolationNode {
	return &RegexInterpolationNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Represents an Interpolated regex literal eg. `%/foo${1 + 2}bar/`
type InterpolatedRegexLiteralNode struct {
	NodeBase
	Content []RegexLiteralContentNode
	Flags   bitfield.BitField8
}

func (*InterpolatedRegexLiteralNode) IsStatic() bool {
	return false
}

func (r *InterpolatedRegexLiteralNode) IsCaseInsensitive() bool {
	return r.Flags.HasFlag(flag.CaseInsensitiveFlag)
}

func (r *InterpolatedRegexLiteralNode) SetCaseInsensitive() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.CaseInsensitiveFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsMultiline() bool {
	return r.Flags.HasFlag(flag.MultilineFlag)
}

func (r *InterpolatedRegexLiteralNode) SetMultiline() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.MultilineFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsDotAll() bool {
	return r.Flags.HasFlag(flag.DotAllFlag)
}

func (r *InterpolatedRegexLiteralNode) SetDotAll() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.DotAllFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsUngreedy() bool {
	return r.Flags.HasFlag(flag.UngreedyFlag)
}

func (r *InterpolatedRegexLiteralNode) SetUngreedy() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.UngreedyFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsASCII() bool {
	return r.Flags.HasFlag(flag.ASCIIFlag)
}

func (r *InterpolatedRegexLiteralNode) SetASCII() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ASCIIFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsExtended() bool {
	return r.Flags.HasFlag(flag.ExtendedFlag)
}

func (r *InterpolatedRegexLiteralNode) SetExtended() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ExtendedFlag)
	return r
}

// Create a new interpolated regex literal node eg. `%/foo${1 + 3}bar/`.
func NewInterpolatedRegexLiteralNode(span *position.Span, content []RegexLiteralContentNode, flags bitfield.BitField8) *InterpolatedRegexLiteralNode {
	return &InterpolatedRegexLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  content,
		Flags:    flags,
	}
}

// Represents a private identifier eg. `_foo`
type PrivateIdentifierNode struct {
	NodeBase
	Value string
}

func (*PrivateIdentifierNode) IsStatic() bool {
	return false
}

// Create a new private identifier node eg. `_foo`.
func NewPrivateIdentifierNode(span *position.Span, val string) *PrivateIdentifierNode {
	return &PrivateIdentifierNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a public constant eg. `Foo`.
type PublicConstantNode struct {
	NodeBase
	Value string
}

func (*PublicConstantNode) IsStatic() bool {
	return false
}

// Create a new public constant node eg. `Foo`.
func NewPublicConstantNode(span *position.Span, val string) *PublicConstantNode {
	return &PublicConstantNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents an instance variable eg. `@foo`
type InstanceVariableNode struct {
	NodeBase
	Value string
}

func (*InstanceVariableNode) IsStatic() bool {
	return false
}

// Create an instance variable node eg. `@foo`.
func NewInstanceVariableNode(span *position.Span, val string) *InstanceVariableNode {
	return &InstanceVariableNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a private constant eg. `_Foo`
type PrivateConstantNode struct {
	NodeBase
	Value string
}

func (*PrivateConstantNode) IsStatic() bool {
	return false
}

// Create a new private constant node eg. `_Foo`.
func NewPrivateConstantNode(span *position.Span, val string) *PrivateConstantNode {
	return &PrivateConstantNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a `do` expression eg.
//
//	do
//		print("awesome!")
//	end
type DoExpressionNode struct {
	NodeBase
	Body []StatementNode // do expression body
}

func (*DoExpressionNode) IsStatic() bool {
	return false
}

// Create a new `do` expression node eg.
//
//	do
//		print("awesome!")
//	end
func NewDoExpressionNode(span *position.Span, body []StatementNode) *DoExpressionNode {
	return &DoExpressionNode{
		NodeBase: NodeBase{span: span},
		Body:     body,
	}
}

// Represents a `singleton` block expression eg.
//
//	singleton
//		def hello then println("awesome!")
//	end
type SingletonBlockExpressionNode struct {
	NodeBase
	Body []StatementNode // do expression body
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
		NodeBase: NodeBase{span: span},
		Body:     body,
	}
}

// Represents an `if`, `unless`, `while` or `until` modifier expression eg. `return true if foo`.
type ModifierNode struct {
	NodeBase
	Modifier *token.Token   // modifier token
	Left     ExpressionNode // left hand side
	Right    ExpressionNode // right hand side
}

func (*ModifierNode) IsStatic() bool {
	return false
}

// Create a new modifier node eg. `return true if foo`.
func NewModifierNode(span *position.Span, mod *token.Token, left, right ExpressionNode) *ModifierNode {
	return &ModifierNode{
		NodeBase: NodeBase{span: span},
		Modifier: mod,
		Left:     left,
		Right:    right,
	}
}

// Represents an `if .. else` modifier expression eg. `foo = 1 if bar else foo = 2`
type ModifierIfElseNode struct {
	NodeBase
	ThenExpression ExpressionNode // then expression body
	Condition      ExpressionNode // if condition
	ElseExpression ExpressionNode // else expression body
}

func (*ModifierIfElseNode) IsStatic() bool {
	return false
}

// Create a new modifier `if` .. `else` node eg. `foo = 1 if bar else foo = 2“.
func NewModifierIfElseNode(span *position.Span, then, cond, els ExpressionNode) *ModifierIfElseNode {
	return &ModifierIfElseNode{
		NodeBase:       NodeBase{span: span},
		ThenExpression: then,
		Condition:      cond,
		ElseExpression: els,
	}
}

// Represents an `for .. in` modifier expression eg. `println(i) for i in 10..30`
type ModifierForInNode struct {
	NodeBase
	ThenExpression ExpressionNode // then expression body
	Parameter      IdentifierNode
	InExpression   ExpressionNode // expression that will be iterated through
}

func (*ModifierForInNode) IsStatic() bool {
	return false
}

// Create a new modifier `for` .. `in` node eg. `println(i) for i in 10..30`
func NewModifierForInNode(span *position.Span, then ExpressionNode, param IdentifierNode, in ExpressionNode) *ModifierForInNode {
	return &ModifierForInNode{
		NodeBase:       NodeBase{span: span},
		ThenExpression: then,
		Parameter:      param,
		InExpression:   in,
	}
}

// Represents an `if` expression eg. `if foo then println("bar")`
type IfExpressionNode struct {
	NodeBase
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
		NodeBase:  NodeBase{span: span},
		ThenBody:  then,
		Condition: cond,
		ElseBody:  els,
	}
}

// Represents an `unless` expression eg. `unless foo then println("bar")`
type UnlessExpressionNode struct {
	NodeBase
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
		NodeBase:  NodeBase{span: span},
		ThenBody:  then,
		Condition: cond,
		ElseBody:  els,
	}
}

// Represents a `while` expression eg. `while i < 5 then i += 5`
type WhileExpressionNode struct {
	NodeBase
	Condition ExpressionNode  // while condition
	ThenBody  []StatementNode // then expression body
}

func (*WhileExpressionNode) IsStatic() bool {
	return false
}

// Create a new `while` expression node eg. `while i < 5 then i += 5`
func NewWhileExpressionNode(span *position.Span, cond ExpressionNode, then []StatementNode) *WhileExpressionNode {
	return &WhileExpressionNode{
		NodeBase:  NodeBase{span: span},
		Condition: cond,
		ThenBody:  then,
	}
}

// Represents a `until` expression eg. `until i >= 5 then i += 5`
type UntilExpressionNode struct {
	NodeBase
	Condition ExpressionNode  // until condition
	ThenBody  []StatementNode // then expression body
}

func (*UntilExpressionNode) IsStatic() bool {
	return false
}

// Create a new `until` expression node eg. `until i >= 5 then i += 5`
func NewUntilExpressionNode(span *position.Span, cond ExpressionNode, then []StatementNode) *UntilExpressionNode {
	return &UntilExpressionNode{
		NodeBase:  NodeBase{span: span},
		Condition: cond,
		ThenBody:  then,
	}
}

// Represents a `loop` expression.
type LoopExpressionNode struct {
	NodeBase
	ThenBody []StatementNode // then expression body
}

func (*LoopExpressionNode) IsStatic() bool {
	return false
}

// Create a new `loop` expression node eg. `loop println('elk is awesome')`
func NewLoopExpressionNode(span *position.Span, then []StatementNode) *LoopExpressionNode {
	return &LoopExpressionNode{
		NodeBase: NodeBase{span: span},
		ThenBody: then,
	}
}

// Represents a numeric `for` expression eg. `fornum i := 0; i < 10; i += 1 then println(i)`
type NumericForExpressionNode struct {
	NodeBase
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
		NodeBase:    NodeBase{span: span},
		Initialiser: init,
		Condition:   cond,
		Increment:   incr,
		ThenBody:    then,
	}
}

// Represents a `for in` expression eg. `for i in 5..15 then println(i)`
type ForInExpressionNode struct {
	NodeBase
	Parameter    IdentifierNode  // parameter
	InExpression ExpressionNode  // expression that will be iterated through
	ThenBody     []StatementNode // then expression body
}

func (*ForInExpressionNode) IsStatic() bool {
	return false
}

// Create a new `for in` expression node eg. `for i in 5..15 then println(i)`
func NewForInExpressionNode(span *position.Span, param IdentifierNode, inExpr ExpressionNode, then []StatementNode) *ForInExpressionNode {
	return &ForInExpressionNode{
		NodeBase:     NodeBase{span: span},
		Parameter:    param,
		InExpression: inExpr,
		ThenBody:     then,
	}
}

// Represents a labeled expression eg. `$foo: 1 + 2`
type LabeledExpressionNode struct {
	NodeBase
	Label      string
	Expression ExpressionNode
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

// Create a new `break` expression node eg. `break`
func NewBreakExpressionNode(span *position.Span, label string, val ExpressionNode) *BreakExpressionNode {
	return &BreakExpressionNode{
		NodeBase: NodeBase{span: span},
		Label:    label,
		Value:    val,
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

// Create a new `return` expression node eg. `return`, `return true`
func NewReturnExpressionNode(span *position.Span, val ExpressionNode) *ReturnExpressionNode {
	return &ReturnExpressionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a `continue` expression eg. `continue`, `continue "foo"`
type ContinueExpressionNode struct {
	NodeBase
	Label string
	Value ExpressionNode
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

// Represents a `throw` expression eg. `throw ArgumentError.new("foo")`
type ThrowExpressionNode struct {
	NodeBase
	Value ExpressionNode
}

func (*ThrowExpressionNode) IsStatic() bool {
	return false
}

// Create a new `throw` expression node eg. `throw ArgumentError.new("foo")`
func NewThrowExpressionNode(span *position.Span, val ExpressionNode) *ThrowExpressionNode {
	return &ThrowExpressionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a constant declaration eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
type ConstantDeclarationNode struct {
	NodeBase
	Name        *token.Token   // name of the constant
	Type        TypeNode       // type of the constant
	Initialiser ExpressionNode // value assigned to the constant
}

func (*ConstantDeclarationNode) IsStatic() bool {
	return false
}

// Create a new constant declaration node eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
func NewConstantDeclarationNode(span *position.Span, name *token.Token, typ TypeNode, init ExpressionNode) *ConstantDeclarationNode {
	return &ConstantDeclarationNode{
		NodeBase:    NodeBase{span: span},
		Name:        name,
		Type:        typ,
		Initialiser: init,
	}
}

// Type expression of an operator with two operands eg. `String | Int`
type BinaryTypeExpressionNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Same as [NewBinaryTypeExpressionNode] but returns an interface
func NewBinaryTypeExpressionNodeI(span *position.Span, op *token.Token, left, right TypeNode) TypeNode {
	return &BinaryTypeExpressionNode{
		NodeBase: NodeBase{span: span},
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Represents an optional or nilable type eg. `String?`
type NilableTypeNode struct {
	NodeBase
	Type TypeNode // right hand side
}

func (*NilableTypeNode) IsStatic() bool {
	return false
}

// Create a new nilable type node eg. `String?`
func NewNilableTypeNode(span *position.Span, typ TypeNode) *NilableTypeNode {
	return &NilableTypeNode{
		NodeBase: NodeBase{span: span},
		Type:     typ,
	}
}

// Represents a singleton type eg. `&String`
type SingletonTypeNode struct {
	NodeBase
	Type TypeNode // right hand side
}

func (*SingletonTypeNode) IsStatic() bool {
	return false
}

// Create a new singleton type node eg. `&String`
func NewSingletonTypeNode(span *position.Span, typ TypeNode) *SingletonTypeNode {
	return &SingletonTypeNode{
		NodeBase: NodeBase{span: span},
		Type:     typ,
	}
}

// Represents a constant lookup expressions eg. `Foo::Bar`
type ConstantLookupNode struct {
	NodeBase
	Left  ExpressionNode      // left hand side
	Right ComplexConstantNode // right hand side
}

func (*ConstantLookupNode) IsStatic() bool {
	return false
}

// Create a new constant lookup expression node eg. `Foo::Bar`
func NewConstantLookupNode(span *position.Span, left ExpressionNode, right ComplexConstantNode) *ConstantLookupNode {
	return &ConstantLookupNode{
		NodeBase: NodeBase{span: span},
		Left:     left,
		Right:    right,
	}
}

// Indicates whether the parameter is a rest param
type ParameterKind uint8

const (
	NormalParameterKind ParameterKind = iota
	PositionalRestParameterKind
	NamedRestParameterKind
)

// Represents a formal parameter in closure or struct declarations eg. `foo: String = 'bar'`
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

// Represents a signature parameter in method and closure signatures eg. `foo?: String`
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

// Represents a closure eg. `|i| -> println(i)`
type ClosureLiteralNode struct {
	NodeBase
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
		NodeBase:   NodeBase{span: span},
		Parameters: params,
		ReturnType: retType,
		ThrowType:  throwType,
		Body:       body,
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
}

func (*ClassDeclarationNode) IsStatic() bool {
	return false
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
) *ClassDeclarationNode {

	return &ClassDeclarationNode{
		NodeBase:      NodeBase{span: span},
		Abstract:      abstract,
		Sealed:        sealed,
		Constant:      constant,
		TypeVariables: typeVars,
		Superclass:    superclass,
		Body:          body,
	}
}

// Represents a module declaration eg. `module Foo; end`
type ModuleDeclarationNode struct {
	NodeBase
	Constant ExpressionNode  // The constant that will hold the module value
	Body     []StatementNode // body of the module
}

func (*ModuleDeclarationNode) IsStatic() bool {
	return false
}

// Create a new module declaration node eg. `module Foo; end`
func NewModuleDeclarationNode(
	span *position.Span,
	constant ExpressionNode,
	body []StatementNode,
) *ModuleDeclarationNode {

	return &ModuleDeclarationNode{
		NodeBase: NodeBase{span: span},
		Constant: constant,
		Body:     body,
	}
}

// Represents a mixin declaration eg. `mixin Foo; end`
type MixinDeclarationNode struct {
	NodeBase
	Constant      ExpressionNode     // The constant that will hold the mixin value
	TypeVariables []TypeVariableNode // Generic type variable definitions
	Body          []StatementNode    // body of the mixin
}

func (*MixinDeclarationNode) IsStatic() bool {
	return false
}

// Create a new mixin declaration node eg. `mixin Foo; end`
func NewMixinDeclarationNode(
	span *position.Span,
	constant ExpressionNode,
	typeVars []TypeVariableNode,
	body []StatementNode,
) *MixinDeclarationNode {

	return &MixinDeclarationNode{
		NodeBase:      NodeBase{span: span},
		Constant:      constant,
		TypeVariables: typeVars,
		Body:          body,
	}
}

// Represents an interface declaration eg. `interface Foo; end`
type InterfaceDeclarationNode struct {
	NodeBase
	Constant      ExpressionNode     // The constant that will hold the interface value
	TypeVariables []TypeVariableNode // Generic type variable definitions
	Body          []StatementNode    // body of the interface
}

func (*InterfaceDeclarationNode) IsStatic() bool {
	return false
}

// Create a new interface declaration node eg. `interface Foo; end`
func NewInterfaceDeclarationNode(
	span *position.Span,
	constant ExpressionNode,
	typeVars []TypeVariableNode,
	body []StatementNode,
) *InterfaceDeclarationNode {

	return &InterfaceDeclarationNode{
		NodeBase:      NodeBase{span: span},
		Constant:      constant,
		TypeVariables: typeVars,
		Body:          body,
	}
}

// Represents a struct declaration eg. `struct Foo; end`
type StructDeclarationNode struct {
	NodeBase
	Constant      ExpressionNode            // The constant that will hold the struct value
	TypeVariables []TypeVariableNode        // Generic type variable definitions
	Body          []StructBodyStatementNode // body of the struct
}

func (*StructDeclarationNode) IsStatic() bool {
	return false
}

// Create a new struct declaration node eg. `struct Foo; end`
func NewStructDeclarationNode(
	span *position.Span,
	constant ExpressionNode,
	typeVars []TypeVariableNode,
	body []StructBodyStatementNode,
) *StructDeclarationNode {

	return &StructDeclarationNode{
		NodeBase:      NodeBase{span: span},
		Constant:      constant,
		TypeVariables: typeVars,
		Body:          body,
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

// Represents a symbol literal with simple content eg. `:foo`, `:'foo bar`, `:"lol"`
type SimpleSymbolLiteralNode struct {
	NodeBase
	Content string
}

func (*SimpleSymbolLiteralNode) IsStatic() bool {
	return true
}

// Create a simple symbol literal node eg. `:foo`, `:'foo bar`, `:"lol"`
func NewSimpleSymbolLiteralNode(span *position.Span, cont string) *SimpleSymbolLiteralNode {
	return &SimpleSymbolLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  cont,
	}
}

// Represents an interpolated symbol eg. `:"foo ${bar + 2}"`
type InterpolatedSymbolLiteral struct {
	NodeBase
	Content *InterpolatedStringLiteralNode
}

func (*InterpolatedSymbolLiteral) IsStatic() bool {
	return false
}

// Create an interpolated symbol literal node eg. `:"foo ${bar + 2}"`
func NewInterpolatedSymbolLiteral(span *position.Span, cont *InterpolatedStringLiteralNode) *InterpolatedSymbolLiteral {
	return &InterpolatedSymbolLiteral{
		NodeBase: NodeBase{span: span},
		Content:  cont,
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
	return len(m.Name) > 0 && m.Name[len(m.Name)-1] == '='
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

// Represents a generic constant in type annotations eg. `ArrayList[String]`
type GenericConstantNode struct {
	NodeBase
	Constant         ComplexConstantNode
	GenericArguments []ComplexConstantNode
}

func (*GenericConstantNode) IsStatic() bool {
	return true
}

// Create a generic constant node eg. `ArrayList[String]`
func NewGenericConstantNode(span *position.Span, constant ComplexConstantNode, args []ComplexConstantNode) *GenericConstantNode {
	return &GenericConstantNode{
		NodeBase:         NodeBase{span: span},
		Constant:         constant,
		GenericArguments: args,
	}
}

// Represents a new type definition eg. `typedef StringList = ArrayList[String]`
type TypeDefinitionNode struct {
	NodeBase
	Constant ComplexConstantNode // new name of the type
	Type     TypeNode            // the type
}

func (*TypeDefinitionNode) IsStatic() bool {
	return false
}

// Create a type definition node eg. `typedef StringList = ArrayList[String]`
func NewTypeDefinitionNode(span *position.Span, constant ComplexConstantNode, typ TypeNode) *TypeDefinitionNode {
	return &TypeDefinitionNode{
		NodeBase: NodeBase{span: span},
		Constant: constant,
		Type:     typ,
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
	NodeBase
	Entries []*AliasDeclarationEntry
}

func (*AliasDeclarationNode) IsStatic() bool {
	return false
}

// Create an alias declaration node eg. `alias push append, add plus`
func NewAliasDeclarationNode(span *position.Span, entries []*AliasDeclarationEntry) *AliasDeclarationNode {
	return &AliasDeclarationNode{
		NodeBase: NodeBase{span: span},
		Entries:  entries,
	}
}

// Represents a new getter declaration eg. `getter foo: String`
type GetterDeclarationNode struct {
	NodeBase
	Entries []ParameterNode
}

func (*GetterDeclarationNode) IsStatic() bool {
	return false
}

// Create a getter declaration node eg. `getter foo: String`
func NewGetterDeclarationNode(span *position.Span, entries []ParameterNode) *GetterDeclarationNode {
	return &GetterDeclarationNode{
		NodeBase: NodeBase{span: span},
		Entries:  entries,
	}
}

// Represents a new setter declaration eg. `setter foo: String`
type SetterDeclarationNode struct {
	NodeBase
	Entries []ParameterNode
}

func (*SetterDeclarationNode) IsStatic() bool {
	return false
}

// Create a setter declaration node eg. `setter foo: String`
func NewSetterDeclarationNode(span *position.Span, entries []ParameterNode) *SetterDeclarationNode {
	return &SetterDeclarationNode{
		NodeBase: NodeBase{span: span},
		Entries:  entries,
	}
}

// Represents a new setter declaration eg. `accessor foo: String`
type AccessorDeclarationNode struct {
	NodeBase
	Entries []ParameterNode
}

func (*AccessorDeclarationNode) IsStatic() bool {
	return false
}

// Create an accessor declaration node eg. `accessor foo: String`
func NewAccessorDeclarationNode(span *position.Span, entries []ParameterNode) *AccessorDeclarationNode {
	return &AccessorDeclarationNode{
		NodeBase: NodeBase{span: span},
		Entries:  entries,
	}
}

// Represents an include expression eg. `include Enumerable[V]`
type IncludeExpressionNode struct {
	NodeBase
	Constants []ComplexConstantNode
}

func (*IncludeExpressionNode) IsStatic() bool {
	return false
}

// Create an include expression node eg. `include Enumerable[V]`
func NewIncludeExpressionNode(span *position.Span, consts []ComplexConstantNode) *IncludeExpressionNode {
	return &IncludeExpressionNode{
		NodeBase:  NodeBase{span: span},
		Constants: consts,
	}
}

// Represents an extend expression eg. `extend Enumerable[V]`
type ExtendExpressionNode struct {
	NodeBase
	Constants []ComplexConstantNode
}

func (*ExtendExpressionNode) IsStatic() bool {
	return false
}

// Create an extend expression node eg. `extend Enumerable[V]`
func NewExtendExpressionNode(span *position.Span, consts []ComplexConstantNode) *ExtendExpressionNode {
	return &ExtendExpressionNode{
		NodeBase:  NodeBase{span: span},
		Constants: consts,
	}
}

// Represents an enhance expression eg. `enhance Enumerable[V]`
type EnhanceExpressionNode struct {
	NodeBase
	Constants []ComplexConstantNode
}

func (*EnhanceExpressionNode) IsStatic() bool {
	return false
}

// Create an enhance expression node eg. `enhance Enumerable[V]`
func NewEnhanceExpressionNode(span *position.Span, consts []ComplexConstantNode) *EnhanceExpressionNode {
	return &EnhanceExpressionNode{
		NodeBase:  NodeBase{span: span},
		Constants: consts,
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

// Represents a constructor call eg. `String(123)`
type ConstructorCallNode struct {
	NodeBase
	Class               ComplexConstantNode // class that is being instantiated
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*ConstructorCallNode) IsStatic() bool {
	return false
}

// Create a constructor call node eg. `String(123)`
func NewConstructorCallNode(span *position.Span, class ComplexConstantNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ConstructorCallNode {
	return &ConstructorCallNode{
		NodeBase:            NodeBase{span: span},
		Class:               class,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents attribute access eg. `foo.bar`
type AttributeAccessNode struct {
	NodeBase
	Receiver      ExpressionNode
	AttributeName string
}

func (*AttributeAccessNode) IsStatic() bool {
	return false
}

// Create an attribute access node eg. `foo.bar`
func NewAttributeAccessNode(span *position.Span, recv ExpressionNode, attrName string) *AttributeAccessNode {
	return &AttributeAccessNode{
		NodeBase:      NodeBase{span: span},
		Receiver:      recv,
		AttributeName: attrName,
	}
}

// Represents subscript access eg. `arr[5]`
type SubscriptExpressionNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Receiver: recv,
		Key:      key,
		static:   recv.IsStatic() && key.IsStatic(),
	}
}

// Represents nil-safe subscript access eg. `arr?[5]`
type NilSafeSubscriptExpressionNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Receiver: recv,
		Key:      key,
		static:   recv.IsStatic() && key.IsStatic(),
	}
}

// Represents a method call eg. `'123'.to_int()`
type MethodCallNode struct {
	NodeBase
	Receiver            ExpressionNode
	NilSafe             bool
	MethodName          string
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*MethodCallNode) IsStatic() bool {
	return false
}

// Create a method call node eg. `'123'.to_int()`
func NewMethodCallNode(span *position.Span, recv ExpressionNode, nilSafe bool, methodName string, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *MethodCallNode {
	return &MethodCallNode{
		NodeBase:            NodeBase{span: span},
		Receiver:            recv,
		NilSafe:             nilSafe,
		MethodName:          methodName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a function-like call eg. `to_string(123)`
type FunctionCallNode struct {
	NodeBase
	MethodName          string
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*FunctionCallNode) IsStatic() bool {
	return false
}

// Create a function call node eg. `to_string(123)`
func NewFunctionCallNode(span *position.Span, methodName string, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *FunctionCallNode {
	return &FunctionCallNode{
		NodeBase:            NodeBase{span: span},
		MethodName:          methodName,
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
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Key:      key,
		Value:    val,
		static:   areExpressionsStatic(key, val),
	}
}

// Represents a ArrayList literal eg. `[1, 5, -6]`
type ArrayListLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewArrayListLiteralNode] but returns an interface
func NewArrayListLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewArrayListLiteralNode(span, elements, capacity)
}

// Represents a word ArrayList literal eg. `\w[foo bar]`
type WordArrayListLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewWordArrayListLiteralNode] but returns an interface.
func NewWordArrayListLiteralNodeI(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewWordArrayListLiteralNode(span, elements, capacity)
}

// Represents a word ArrayTuple literal eg. `%w[foo bar]`
type WordArrayTupleLiteralNode struct {
	NodeBase
	Elements []WordCollectionContentNode
}

func (*WordArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a word ArrayTuple literal node eg. `%w[foo bar]`
func NewWordArrayTupleLiteralNode(span *position.Span, elements []WordCollectionContentNode) *WordArrayTupleLiteralNode {
	return &WordArrayTupleLiteralNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewWordArrayTupleLiteralNode] but returns an interface.
func NewWordArrayTupleLiteralNodeI(span *position.Span, elements []WordCollectionContentNode) ExpressionNode {
	return &WordArrayTupleLiteralNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Represents a word HashSet literal eg. `^w[foo bar]`
type WordHashSetLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewWordHashSetLiteralNode] but returns an interface.
func NewWordHashSetLiteralNodeI(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewWordHashSetLiteralNode(span, elements, capacity)
}

// Represents a symbol ArrayList literal eg. `\s[foo bar]`
type SymbolArrayListLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewSymbolArrayListLiteralNode] but returns an interface.
func NewSymbolArrayListLiteralNodeI(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewSymbolArrayListLiteralNode(span, elements, capacity)
}

// Represents a symbol ArrayTuple literal eg. `%s[foo bar]`
type SymbolArrayTupleLiteralNode struct {
	NodeBase
	Elements []SymbolCollectionContentNode
}

func (*SymbolArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a symbol arrayTuple literal node eg. `%s[foo bar]`
func NewSymbolArrayTupleLiteralNode(span *position.Span, elements []SymbolCollectionContentNode) *SymbolArrayTupleLiteralNode {
	return &SymbolArrayTupleLiteralNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewSymbolArrayTupleLiteralNode] but returns an interface.
func NewSymbolArrayTupleLiteralNodeI(span *position.Span, elements []SymbolCollectionContentNode) ExpressionNode {
	return &SymbolArrayTupleLiteralNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Represents a symbol HashSet literal eg. `^s[foo bar]`
type SymbolHashSetLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewSymbolHashSetLiteralNode] but returns an interface.
func NewSymbolHashSetLiteralNodeI(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewSymbolHashSetLiteralNode(span, elements, capacity)
}

// Represents a hex ArrayList literal eg. `\x[ff ee]`
type HexArrayListLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewHexArrayListLiteralNode] but returns an interface.
func NewHexArrayListLiteralNodeI(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewHexArrayListLiteralNode(span, elements, capacity)
}

// Represents a hex ArrayTuple literal eg. `%x[ff ee]`
type HexArrayTupleLiteralNode struct {
	NodeBase
	Elements []IntCollectionContentNode
}

func (*HexArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a hex ArrayTuple literal node eg. `%x[ff ee]`
func NewHexArrayTupleLiteralNode(span *position.Span, elements []IntCollectionContentNode) *HexArrayTupleLiteralNode {
	return &HexArrayTupleLiteralNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewHexArrayTupleLiteralNode] but returns an interface.
func NewHexArrayTupleLiteralNodeI(span *position.Span, elements []IntCollectionContentNode) ExpressionNode {
	return &HexArrayTupleLiteralNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Represents a hex HashSet literal eg. `^x[ff ee}]`
type HexHashSetLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewHexHashSetLiteralNode] but returns an interface.
func NewHexHashSetLiteralNodeI(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewHexHashSetLiteralNode(span, elements, capacity)
}

// Represents a bin ArrayList literal eg. `\b[11 10]`
type BinArrayListLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewBinArrayListLiteralNode] but returns an interface.
func NewBinArrayListLiteralNodeI(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewBinArrayListLiteralNode(span, elements, capacity)
}

// Represents a bin ArrayTuple literal eg. `%b[11 10]`
type BinArrayTupleLiteralNode struct {
	NodeBase
	Elements []IntCollectionContentNode
}

func (*BinArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a bin ArrayList literal node eg. `%b[11 10]`
func NewBinArrayTupleLiteralNode(span *position.Span, elements []IntCollectionContentNode) *BinArrayTupleLiteralNode {
	return &BinArrayTupleLiteralNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewBinArrayTupleLiteralNode] but returns an interface.
func NewBinArrayTupleLiteralNodeI(span *position.Span, elements []IntCollectionContentNode) ExpressionNode {
	return &BinArrayTupleLiteralNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Represents a bin HashSet literal eg. `^b[11 10]`
type BinHashSetLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewBinHashSetLiteralNode] but returns an interface.
func NewBinHashSetLiteralNodeI(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewBinHashSetLiteralNode(span, elements, capacity)
}

// Represents a ArrayTuple literal eg. `%[1, 5, -6]`
type ArrayTupleLiteralNode struct {
	NodeBase
	Elements []ExpressionNode
	static   bool
}

func (t *ArrayTupleLiteralNode) IsStatic() bool {
	return t.static
}

// Create a ArrayTuple literal node eg. `%[1, 5, -6]`
func NewArrayTupleLiteralNode(span *position.Span, elements []ExpressionNode) *ArrayTupleLiteralNode {
	return &ArrayTupleLiteralNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
		static:   isExpressionSliceStatic(elements),
	}
}

// Same as [NewArrayTupleLiteralNode] but returns an interface
func NewArrayTupleLiteralNodeI(span *position.Span, elements []ExpressionNode) ExpressionNode {
	return NewArrayTupleLiteralNode(span, elements)
}

// Represents a HashSet literal eg. `^[1, 5, -6]`
type HashSetLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewHashSetLiteralNode] but returns an interface
func NewHashSetLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewHashSetLiteralNode(span, elements, capacity)
}

// Represents a HashMap literal eg. `{ foo: 1, 'bar' => 5, baz }`
type HashMapLiteralNode struct {
	NodeBase
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
		NodeBase: NodeBase{span: span},
		Elements: elements,
		Capacity: capacity,
		static:   static,
	}
}

// Same as [NewHashMapLiteralNode] but returns an interface
func NewHashMapLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewHashMapLiteralNode(span, elements, capacity)
}

// Represents a Record literal eg. `%{ foo: 1, 'bar' => 5, baz }`
type HashRecordLiteralNode struct {
	NodeBase
	Elements []ExpressionNode
	static   bool
}

func (r *HashRecordLiteralNode) IsStatic() bool {
	return r.static
}

// Create a Record literal node eg. `%{ foo: 1, 'bar' => 5, baz }`
func NewHashRecordLiteralNode(span *position.Span, elements []ExpressionNode) *HashRecordLiteralNode {
	return &HashRecordLiteralNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
		static:   isExpressionSliceStatic(elements),
	}
}

// Same as [NewHashRecordLiteralNode] but returns an interface
func NewHashRecordLiteralNodeI(span *position.Span, elements []ExpressionNode) ExpressionNode {
	return NewHashRecordLiteralNode(span, elements)
}

// Represents a Range literal eg. `1..5`
type RangeLiteralNode struct {
	NodeBase
	From      ExpressionNode
	To        ExpressionNode
	Exclusive bool
	static    bool
}

func (r *RangeLiteralNode) IsStatic() bool {
	return r.static
}

// Create a Range literal node eg. `1..5`
func NewRangeLiteralNode(span *position.Span, exclusive bool, from, to ExpressionNode) *RangeLiteralNode {
	return &RangeLiteralNode{
		NodeBase:  NodeBase{span: span},
		Exclusive: exclusive,
		From:      from,
		To:        to,
		static:    areExpressionsStatic(from, to),
	}
}

// Represents an ArithmeticSequence literal eg. `1..5:2`
type ArithmeticSequenceLiteralNode struct {
	NodeBase
	From      ExpressionNode
	To        ExpressionNode
	Step      ExpressionNode
	Exclusive bool
	static    bool
}

func (a *ArithmeticSequenceLiteralNode) IsStatic() bool {
	return a.static
}

// Create an ArithmeticSequence literal eg. `1..5:2`
func NewArithmeticSequenceLiteralNode(span *position.Span, exclusive bool, from, to, step ExpressionNode) *ArithmeticSequenceLiteralNode {
	return &ArithmeticSequenceLiteralNode{
		NodeBase:  NodeBase{span: span},
		Exclusive: exclusive,
		From:      from,
		To:        to,
		Step:      step,
		static:    areExpressionsStatic(from, to, step),
	}
}

// Represents a doc comment eg.
//
//	##[foo bar]##
//	def foo; end
type DocCommentNode struct {
	NodeBase
	Comment    string
	Expression ExpressionNode
}

func (*DocCommentNode) IsStatic() bool {
	return false
}

// Create a doc comment.
func NewDocCommentNode(span *position.Span, comment string, expr ExpressionNode) *DocCommentNode {
	return &DocCommentNode{
		NodeBase:   NodeBase{span: span},
		Comment:    comment,
		Expression: expr,
	}
}
