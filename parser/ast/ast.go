// Package ast defines types
// used by the Elk parser.
//
// All the nodes of the Abstract Syntax Tree
// constructed by the Elk parser are defined in this package.
package ast

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
)

// Every node type implements this interface.
type Node interface {
	position.Interface
}

// Check whether the token can be used as a left value
// in a variable/constant declaration.
func IsValidDeclarationTarget(node Node) bool {
	switch node.(type) {
	case *PrivateConstantNode, *PublicConstantNode, *PrivateIdentifierNode, *PublicIdentifierNode:
		return true
	default:
		return false
	}
}

// Check whether the token can be used as a left value
// in an assignment expression.
func IsValidAssignmentTarget(node Node) bool {
	switch node.(type) {
	case *PrivateIdentifierNode, *PublicIdentifierNode:
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

// Represents a single statement, so for example
// a single valid "line" of Elk code.
// Usually its an expression optionally terminated with a newline or a semicolon.
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

func (*InvalidNode) expressionNode()                   {}
func (*ModifierNode) expressionNode()                  {}
func (*ModifierIfElseNode) expressionNode()            {}
func (*AssignmentExpressionNode) expressionNode()      {}
func (*BinaryExpressionNode) expressionNode()          {}
func (*LogicalExpressionNode) expressionNode()         {}
func (*UnaryExpressionNode) expressionNode()           {}
func (*TrueLiteralNode) expressionNode()               {}
func (*FalseLiteralNode) expressionNode()              {}
func (*NilLiteralNode) expressionNode()                {}
func (*RawStringLiteralNode) expressionNode()          {}
func (*SimpleSymbolLiteralNode) expressionNode()       {}
func (*ComplexSymbolLiteralNode) expressionNode()      {}
func (*NamedValueLiteralNode) expressionNode()         {}
func (*IntLiteralNode) expressionNode()                {}
func (*FloatLiteralNode) expressionNode()              {}
func (*StringLiteralNode) expressionNode()             {}
func (*PublicIdentifierNode) expressionNode()          {}
func (*PrivateIdentifierNode) expressionNode()         {}
func (*PublicConstantNode) expressionNode()            {}
func (*PrivateConstantNode) expressionNode()           {}
func (*SelfLiteralNode) expressionNode()               {}
func (*IfExpressionNode) expressionNode()              {}
func (*UnlessExpressionNode) expressionNode()          {}
func (*WhileExpressionNode) expressionNode()           {}
func (*UntilExpressionNode) expressionNode()           {}
func (*LoopExpressionNode) expressionNode()            {}
func (*BreakExpressionNode) expressionNode()           {}
func (*ReturnExpressionNode) expressionNode()          {}
func (*ContinueExpressionNode) expressionNode()        {}
func (*ThrowExpressionNode) expressionNode()           {}
func (*VariableDeclarationNode) expressionNode()       {}
func (*ConstantDeclarationNode) expressionNode()       {}
func (*ConstantLookupNode) expressionNode()            {}
func (*FormalParameterNode) expressionNode()           {}
func (*ClosureExpressionNode) expressionNode()         {}
func (*ClassDeclarationNode) expressionNode()          {}
func (*ModuleDeclarationNode) expressionNode()         {}
func (*MixinDeclarationNode) expressionNode()          {}
func (*InterfaceDeclarationNode) expressionNode()      {}
func (*MethodDefinitionNode) expressionNode()          {}
func (*MethodSignatureDefinitionNode) expressionNode() {}
func (*GenericConstantNode) expressionNode()           {}
func (*TypeDefinitionNode) expressionNode()            {}
func (*AliasExpressionNode) expressionNode()           {}
func (*IncludeExpressionNode) expressionNode()         {}
func (*ExtendExpressionNode) expressionNode()          {}
func (*EnhanceExpressionNode) expressionNode()         {}

// All nodes that should be valid in type annotations should
// implement this interface
type TypeNode interface {
	Node
	typeNode()
}

func (*InvalidNode) typeNode()              {}
func (*BinaryTypeExpressionNode) typeNode() {}
func (*NilableTypeNode) typeNode()          {}
func (*PublicConstantNode) typeNode()       {}
func (*PrivateConstantNode) typeNode()      {}
func (*ConstantLookupNode) typeNode()       {}
func (*GenericConstantNode) typeNode()      {}

// All nodes that should be valid in parameter declaration lists
// of methods or closures should implement this interface.
type ParameterNode interface {
	Node
	parameterNode()
}

func (*InvalidNode) parameterNode()         {}
func (*FormalParameterNode) parameterNode() {}

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

// Nodes that implement this interface represent
// symbol literals.
type SymbolLiteralNode interface {
	Node
	ExpressionNode
	symbolLiteralNode()
}

func (*InvalidNode) symbolLiteralNode()              {}
func (*SimpleSymbolLiteralNode) symbolLiteralNode()  {}
func (*ComplexSymbolLiteralNode) symbolLiteralNode() {}

// Represents a single Elk program (usually a single file).
type ProgramNode struct {
	*position.Position
	Body []StatementNode
}

// Create a new program node.
func NewProgramNode(pos *position.Position, body []StatementNode) *ProgramNode {
	return &ProgramNode{
		Position: pos,
		Body:     body,
	}
}

// Represents an empty statement eg. a statement with only a semicolon or a newline.
type EmptyStatementNode struct {
	*position.Position
}

// Create a new empty statement node.
func NewEmptyStatementNode(pos *position.Position) *EmptyStatementNode {
	return &EmptyStatementNode{
		Position: pos,
	}
}

// Expression optionally terminated with a newline or a semicolon.
type ExpressionStatementNode struct {
	*position.Position
	Expression ExpressionNode
}

// Create a new expression statement node eg. `5 * 2\n`
func NewExpressionStatementNode(pos *position.Position, expr ExpressionNode) *ExpressionStatementNode {
	return &ExpressionStatementNode{
		Position:   pos,
		Expression: expr,
	}
}

// Assignment with the specified operator.
type AssignmentExpressionNode struct {
	*position.Position
	Op    *token.Token   // operator
	Left  ExpressionNode // left hand side
	Right ExpressionNode // right hand side
}

// Create a new assignment expression node eg. `foo = 3`
func NewAssignmentExpressionNode(pos *position.Position, op *token.Token, left, right ExpressionNode) *AssignmentExpressionNode {
	return &AssignmentExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Expression of an operator with two operands eg. `2 + 5`, `foo > bar`
type BinaryExpressionNode struct {
	*position.Position
	Op    *token.Token   // operator
	Left  ExpressionNode // left hand side
	Right ExpressionNode // right hand side
}

// Create a new binary expression node.
func NewBinaryExpressionNode(pos *position.Position, op *token.Token, left, right ExpressionNode) *BinaryExpressionNode {
	return &BinaryExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Create a new binary expression node and wrap it in the ExpressionNode interface
func NewBinaryExpressionNodeI(pos *position.Position, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return &BinaryExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Expression of a logical operator with two operands eg. `foo &&  bar`
type LogicalExpressionNode struct {
	*position.Position
	Op    *token.Token   // operator
	Left  ExpressionNode // left hand side
	Right ExpressionNode // right hand side
}

// Create a new logical expression node.
func NewLogicalExpressionNode(pos *position.Position, op *token.Token, left, right ExpressionNode) *LogicalExpressionNode {
	return &LogicalExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Create a new logical expression node and wrap it in the ExpressionNode interface
func NewLogicalExpressionNodeI(pos *position.Position, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return &LogicalExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Expression of an operator with one operand eg. `!foo`, `-bar`
type UnaryExpressionNode struct {
	*position.Position
	Op    *token.Token   // operator
	Right ExpressionNode // right hand side
}

// Create a new unary expression node.
func NewUnaryExpressionNode(pos *position.Position, op *token.Token, right ExpressionNode) *UnaryExpressionNode {
	return &UnaryExpressionNode{
		Position: pos,
		Op:       op,
		Right:    right,
	}
}

// `true` literal.
type TrueLiteralNode struct {
	*position.Position
}

// Create a new `true` literal node.
func NewTrueLiteralNode(pos *position.Position) *TrueLiteralNode {
	return &TrueLiteralNode{
		Position: pos,
	}
}

// `self` literal.
type FalseLiteralNode struct {
	*position.Position
}

// Create a new `false` literal node.
func NewFalseLiteralNode(pos *position.Position) *FalseLiteralNode {
	return &FalseLiteralNode{
		Position: pos,
	}
}

// `self` literal.
type SelfLiteralNode struct {
	*position.Position
}

// Create a new `self` literal node.
func NewSelfLiteralNode(pos *position.Position) *SelfLiteralNode {
	return &SelfLiteralNode{
		Position: pos,
	}
}

// `nil` literal.
type NilLiteralNode struct {
	*position.Position
}

// Create a new `nil` literal node.
func NewNilLiteralNode(pos *position.Position) *NilLiteralNode {
	return &NilLiteralNode{
		Position: pos,
	}
}

// Raw string literal enclosed with single quotes eg. `'foo'`.
type RawStringLiteralNode struct {
	*position.Position
	Value string // value of the string literal
}

// Create a new raw string literal node eg. `'foo'`.
func NewRawStringLiteralNode(pos *position.Position, val string) *RawStringLiteralNode {
	return &RawStringLiteralNode{
		Position: pos,
		Value:    val,
	}
}

// Int literal eg. `5`, `125_355`, `0xff`
type IntLiteralNode struct {
	*position.Position
	Token *token.Token
}

// Create a new raw string literal node eg. `5`, `125_355`, `0xff`
func NewIntLiteralNode(pos *position.Position, tok *token.Token) *IntLiteralNode {
	return &IntLiteralNode{
		Position: pos,
		Token:    tok,
	}
}

// Float literal eg. `5.2`, `.5`, `45e20`
type FloatLiteralNode struct {
	*position.Position
	Value string
}

// Create a new float literal node eg. `5.2`, `.5`, `45e20`
func NewFloatLiteralNode(pos *position.Position, val string) *FloatLiteralNode {
	return &FloatLiteralNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a syntax error.
type InvalidNode struct {
	*position.Position
	Token *token.Token
}

// Create a new invalid node.
func NewInvalidNode(pos *position.Position, tok *token.Token) *InvalidNode {
	return &InvalidNode{
		Position: pos,
		Token:    tok,
	}
}

// Represents a single section of characters of a string literal eg. `foo` in `"foo${bar}"`.
type StringLiteralContentSectionNode struct {
	*position.Position
	Value string
}

// Create a new string literal content section node eg. `foo` in `"foo${bar}"`.
func NewStringLiteralContentSectionNode(pos *position.Position, val string) *StringLiteralContentSectionNode {
	return &StringLiteralContentSectionNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a single interpolated section of a string literal eg. `bar + 2` in `"foo${bar + 2}"`
type StringInterpolationNode struct {
	*position.Position
	Expression ExpressionNode
}

// Create a new string interpolation node eg. `bar + 2` in `"foo${bar + 2}"`
func NewStringInterpolationNode(pos *position.Position, expr ExpressionNode) *StringInterpolationNode {
	return &StringInterpolationNode{
		Position:   pos,
		Expression: expr,
	}
}

// Represents a string literal eg. `"foo ${bar} baz"`
type StringLiteralNode struct {
	*position.Position
	Content []StringLiteralContentNode
}

// Create a new string literal node eg. `"foo ${bar} baz"`
func NewStringLiteralNode(pos *position.Position, cont []StringLiteralContentNode) *StringLiteralNode {
	return &StringLiteralNode{
		Position: pos,
		Content:  cont,
	}
}

// Represents a public identifier eg. `foo`.
type PublicIdentifierNode struct {
	*position.Position
	Value string
}

// Create a new public identifier node eg. `foo`.
func NewPublicIdentifierNode(pos *position.Position, val string) *PublicIdentifierNode {
	return &PublicIdentifierNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a private identifier eg. `_foo`
type PrivateIdentifierNode struct {
	*position.Position
	Value string
}

// Create a new private identifier node eg. `_foo`.
func NewPrivateIdentifierNode(pos *position.Position, val string) *PrivateIdentifierNode {
	return &PrivateIdentifierNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a public constant eg. `Foo`.
type PublicConstantNode struct {
	*position.Position
	Value string
}

// Create a new public constant node eg. `Foo`.
func NewPublicConstantNode(pos *position.Position, val string) *PublicConstantNode {
	return &PublicConstantNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a private constant eg. `_Foo`
type PrivateConstantNode struct {
	*position.Position
	Value string
}

// Create a new private constant node eg. `_Foo`.
func NewPrivateConstantNode(pos *position.Position, val string) *PrivateConstantNode {
	return &PrivateConstantNode{
		Position: pos,
		Value:    val,
	}
}

// Represents an `if`, `unless`, `while` or `until` modifier expression eg. `return true if foo`.
type ModifierNode struct {
	*position.Position
	Modifier *token.Token   // modifier token
	Left     ExpressionNode // left hand side
	Right    ExpressionNode // right hand side
}

// Create a new modifier node eg. `return true if foo`.
func NewModifierNode(pos *position.Position, mod *token.Token, left, right ExpressionNode) *ModifierNode {
	return &ModifierNode{
		Position: pos,
		Modifier: mod,
		Left:     left,
		Right:    right,
	}
}

// Represents an `if .. else` modifier expression eg. `foo = 1 if bar else foo = 2`
type ModifierIfElseNode struct {
	*position.Position
	ThenExpression ExpressionNode // then expression body
	Condition      ExpressionNode // if condition
	ElseExpression ExpressionNode // else expression body
}

// Create a new modifier `if` .. `else` node eg. `foo = 1 if bar else foo = 2“.
func NewModifierIfElseNode(pos *position.Position, then, cond, els ExpressionNode) *ModifierIfElseNode {
	return &ModifierIfElseNode{
		Position:       pos,
		ThenExpression: then,
		Condition:      cond,
		ElseExpression: els,
	}
}

// Represents an `if` expression eg. `if foo then println("bar")`
type IfExpressionNode struct {
	*position.Position
	Condition ExpressionNode  // if condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

// Create a new `if` expression node eg. `if foo then println("bar")`
func NewIfExpressionNode(pos *position.Position, cond ExpressionNode, then, els []StatementNode) *IfExpressionNode {
	return &IfExpressionNode{
		Position:  pos,
		ThenBody:  then,
		Condition: cond,
		ElseBody:  els,
	}
}

// Represents an `unless` expression eg. `unless foo then println("bar")`
type UnlessExpressionNode struct {
	*position.Position
	Condition ExpressionNode  // unless condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

// Create a new `unless` expression node eg. `unless foo then println("bar")`
func NewUnlessExpressionNode(pos *position.Position, cond ExpressionNode, then, els []StatementNode) *UnlessExpressionNode {
	return &UnlessExpressionNode{
		Position:  pos,
		ThenBody:  then,
		Condition: cond,
		ElseBody:  els,
	}
}

// Represents a `while` expression eg. `while i < 5 then i += 5`
type WhileExpressionNode struct {
	*position.Position
	Condition ExpressionNode  // while condition
	ThenBody  []StatementNode // then expression body
}

// Create a new `while` expression node eg. `while i < 5 then i += 5`
func NewWhileExpressionNode(pos *position.Position, cond ExpressionNode, then []StatementNode) *WhileExpressionNode {
	return &WhileExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  then,
	}
}

// Represents a `until` expression eg. `until i >= 5 then i += 5`
type UntilExpressionNode struct {
	*position.Position
	Condition ExpressionNode  // until condition
	ThenBody  []StatementNode // then expression body
}

// Create a new `until` expression node eg. `until i >= 5 then i += 5`
func NewUntilExpressionNode(pos *position.Position, cond ExpressionNode, then []StatementNode) *UntilExpressionNode {
	return &UntilExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  then,
	}
}

// Represents a `loop` expression.
type LoopExpressionNode struct {
	*position.Position
	ThenBody []StatementNode // then expression body
}

// Create a new `loop` expression node eg. `loop println('elk is awesome')`
func NewLoopExpressionNode(pos *position.Position, then []StatementNode) *LoopExpressionNode {
	return &LoopExpressionNode{
		Position: pos,
		ThenBody: then,
	}
}

// Represents a `break` expression eg. `break`
type BreakExpressionNode struct {
	*position.Position
}

// Create a new `break` expression node eg. `break`
func NewBreakExpressionNode(pos *position.Position) *BreakExpressionNode {
	return &BreakExpressionNode{
		Position: pos,
	}
}

// Represents a `return` expression eg. `return`, `return true`
type ReturnExpressionNode struct {
	*position.Position
	Value ExpressionNode
}

// Create a new `return` expression node eg. `return`, `return true`
func NewReturnExpressionNode(pos *position.Position, val ExpressionNode) *ReturnExpressionNode {
	return &ReturnExpressionNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a `continue` expression eg. `continue`, `continue "foo"`
type ContinueExpressionNode struct {
	*position.Position
	Value ExpressionNode
}

// Create a new `continue` expression node eg. `continue`, `continue "foo"`
func NewContinueExpressionNode(pos *position.Position, val ExpressionNode) *ContinueExpressionNode {
	return &ContinueExpressionNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a `throw` expression eg. `throw ArgumentError.new("foo")`
type ThrowExpressionNode struct {
	*position.Position
	Value ExpressionNode
}

// Create a new `throw` expression node eg. `throw ArgumentError.new("foo")`
func NewThrowExpressionNode(pos *position.Position, val ExpressionNode) *ThrowExpressionNode {
	return &ThrowExpressionNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a variable declaration eg. `var foo: String`
type VariableDeclarationNode struct {
	*position.Position
	Name        *token.Token   // name of the variable
	Type        TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

// Create a new variable declaration node eg. `var foo: String`
func NewVariableDeclarationNode(pos *position.Position, name *token.Token, typ TypeNode, init ExpressionNode) *VariableDeclarationNode {
	return &VariableDeclarationNode{
		Position:    pos,
		Name:        name,
		Type:        typ,
		Initialiser: init,
	}
}

// Represents a constant declaration eg. `const Foo: List[String] = ["foo", "bar"]`
type ConstantDeclarationNode struct {
	*position.Position
	Name        *token.Token   // name of the constant
	Type        TypeNode       // type of the constant
	Initialiser ExpressionNode // value assigned to the constant
}

// Create a new constant declaration node eg. `const Foo: List[String] = ["foo", "bar"]`
func NewConstantDeclarationNode(pos *position.Position, name *token.Token, typ TypeNode, init ExpressionNode) *ConstantDeclarationNode {
	return &ConstantDeclarationNode{
		Position:    pos,
		Name:        name,
		Type:        typ,
		Initialiser: init,
	}
}

// Type expression of an operator with two operands eg. `String | Int`
type BinaryTypeExpressionNode struct {
	*position.Position
	Op    *token.Token // operator
	Left  TypeNode     // left hand side
	Right TypeNode     // right hand side
}

// Create a new binary type expression node eg. `String | Int`
func NewBinaryTypeExpressionNode(pos *position.Position, op *token.Token, left, right TypeNode) *BinaryTypeExpressionNode {
	return &BinaryTypeExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Create a new binary type expression node eg. `String | Int` and wrap it in the TypeNode interface
func NewBinaryTypeExpressionNodeI(pos *position.Position, op *token.Token, left, right TypeNode) TypeNode {
	return &BinaryTypeExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Represents an optional or nilable type eg. `String?`
type NilableTypeNode struct {
	*position.Position
	Type TypeNode // right hand side
}

// Create a new nilable type node eg. `String?`
func NewNilableTypeNode(pos *position.Position, typ TypeNode) *NilableTypeNode {
	return &NilableTypeNode{
		Position: pos,
		Type:     typ,
	}
}

// Represents a constant lookup expressions eg. `Foo::Bar`
type ConstantLookupNode struct {
	*position.Position
	Left  ExpressionNode      // left hand side
	Right ComplexConstantNode // right hand side
}

// Create a new constant lookup expression node eg. `Foo::Bar`
func NewConstantLookupNode(pos *position.Position, left ExpressionNode, right ComplexConstantNode) *ConstantLookupNode {
	return &ConstantLookupNode{
		Position: pos,
		Left:     left,
		Right:    right,
	}
}

// Represents a formal parameter in method and closure declarations eg. `foo: String`
type FormalParameterNode struct {
	*position.Position
	Name        string         // name of the variable
	Type        TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

// Create a new formal parameter node eg. `foo: String`
func NewFormalParameterNode(pos *position.Position, name string, typ TypeNode, init ExpressionNode) *FormalParameterNode {
	return &FormalParameterNode{
		Position:    pos,
		Name:        name,
		Type:        typ,
		Initialiser: init,
	}
}

// Represents a closure eg. `|i| -> println(i)`
type ClosureExpressionNode struct {
	*position.Position
	Parameters []ParameterNode // formal parameters of the closure separated by semicolons
	ReturnType TypeNode
	Body       []StatementNode // body of the closure
}

// Create a new closure expression node eg. `|i| -> println(i)`
func NewClosureExpressionNode(pos *position.Position, params []ParameterNode, retType TypeNode, body []StatementNode) *ClosureExpressionNode {
	return &ClosureExpressionNode{
		Position:   pos,
		Parameters: params,
		ReturnType: retType,
		Body:       body,
	}
}

// Represents a class declaration eg. `class Foo; end`
type ClassDeclarationNode struct {
	*position.Position
	Constant      ExpressionNode     // The constant that will hold the class object
	TypeVariables []TypeVariableNode // Generic type variable definitions
	Superclass    ExpressionNode     // the super/parent class of this class
	Body          []StatementNode    // body of the class
}

// Create a new class declaration node eg. `class Foo; end`
func NewClassDeclarationNode(
	pos *position.Position,
	constant ExpressionNode,
	typeVars []TypeVariableNode,
	superclass ExpressionNode,
	body []StatementNode,
) *ClassDeclarationNode {

	return &ClassDeclarationNode{
		Position:      pos,
		Constant:      constant,
		TypeVariables: typeVars,
		Superclass:    superclass,
		Body:          body,
	}
}

// Represents a module declaration eg. `module Foo; end`
type ModuleDeclarationNode struct {
	*position.Position
	Constant ExpressionNode  // The constant that will hold the module object
	Body     []StatementNode // body of the module
}

// Create a new module declaration node eg. `module Foo; end`
func NewModuleDeclarationNode(
	pos *position.Position,
	constant ExpressionNode,
	body []StatementNode,
) *ModuleDeclarationNode {

	return &ModuleDeclarationNode{
		Position: pos,
		Constant: constant,
		Body:     body,
	}
}

// Represents a mixin declaration eg. `mixin Foo; end`
type MixinDeclarationNode struct {
	*position.Position
	Constant      ExpressionNode     // The constant that will hold the mixin object
	TypeVariables []TypeVariableNode // Generic type variable definitions
	Body          []StatementNode    // body of the mixin
}

// Create a new mixin declaration node eg. `mixin Foo; end`
func NewMixinDeclarationNode(
	pos *position.Position,
	constant ExpressionNode,
	typeVars []TypeVariableNode,
	body []StatementNode,
) *MixinDeclarationNode {

	return &MixinDeclarationNode{
		Position:      pos,
		Constant:      constant,
		TypeVariables: typeVars,
		Body:          body,
	}
}

// Represents a interface declaration eg. `interface Foo; end`
type InterfaceDeclarationNode struct {
	*position.Position
	Constant      ExpressionNode     // The constant that will hold the interface object
	TypeVariables []TypeVariableNode // Generic type variable definitions
	Body          []StatementNode    // body of the interface
}

// Create a new interface declaration node eg. `interface Foo; end`
func NewInterfaceDeclarationNode(
	pos *position.Position,
	constant ExpressionNode,
	typeVars []TypeVariableNode,
	body []StatementNode,
) *InterfaceDeclarationNode {

	return &InterfaceDeclarationNode{
		Position:      pos,
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
	*position.Position
	Variance   Variance // Variance level of this type variable
	Name       string   // Name of the type variable eg. `T`
	UpperBound ComplexConstantNode
}

// Create a new type variable node eg. `+V`
func NewVariantTypeVariableNode(pos *position.Position, variance Variance, name string, upper ComplexConstantNode) *VariantTypeVariableNode {
	return &VariantTypeVariableNode{
		Position:   pos,
		Variance:   variance,
		Name:       name,
		UpperBound: upper,
	}
}

// Represents a symbol literal with simple content eg. `:foo`, `:'foo bar`
type SimpleSymbolLiteralNode struct {
	*position.Position
	Content string
}

// Create a simple symbol literal node eg. `:foo`, `:'foo bar`
func NewSimpleSymbolLiteralNode(pos *position.Position, cont string) *SimpleSymbolLiteralNode {
	return &SimpleSymbolLiteralNode{
		Position: pos,
		Content:  cont,
	}
}

// Represents a symbol literal with complex content eg. `:"foo\n"`, `:"foo ${bar + 2}"`
type ComplexSymbolLiteralNode struct {
	*position.Position
	Content *StringLiteralNode
}

// Create a simple symbol literal node eg. `:"foo\n"`, `:"foo ${bar + 2}"`
func NewComplexSymbolLiteralNode(pos *position.Position, cont *StringLiteralNode) *ComplexSymbolLiteralNode {
	return &ComplexSymbolLiteralNode{
		Position: pos,
		Content:  cont,
	}
}

// Represents a named value literal eg. `:foo{2}`
type NamedValueLiteralNode struct {
	*position.Position
	Name  string
	Value ExpressionNode
}

// Create a named value node eg. `:foo{2}`
func NewNamedValueLiteralNode(pos *position.Position, name string, value ExpressionNode) *NamedValueLiteralNode {
	return &NamedValueLiteralNode{
		Position: pos,
		Name:     name,
		Value:    value,
	}
}

// Represents a method definition eg. `def foo: String then 'hello world'`
type MethodDefinitionNode struct {
	*position.Position
	Name       string
	Parameters []ParameterNode // formal parameters
	ReturnType TypeNode
	ThrowType  TypeNode
	Body       []StatementNode // body of the method
}

// Create a method definition node eg. `def foo: String then 'hello world'`
func NewMethodDefinitionNode(pos *position.Position, name string, params []ParameterNode, returnType, throwType TypeNode, body []StatementNode) *MethodDefinitionNode {
	return &MethodDefinitionNode{
		Position:   pos,
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		ThrowType:  throwType,
		Body:       body,
	}
}

// Represents a method signature definition eg. `sig to_string(val: Int): String`
type MethodSignatureDefinitionNode struct {
	*position.Position
	Name       string
	Parameters []ParameterNode // formal parameters
	ReturnType TypeNode
	ThrowType  TypeNode
}

// Create a method signature node eg. `sig to_string(val: Int): String`
func NewMethodSignatureDefinitionNode(pos *position.Position, name string, params []ParameterNode, returnType, throwType TypeNode) *MethodSignatureDefinitionNode {
	return &MethodSignatureDefinitionNode{
		Position:   pos,
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		ThrowType:  throwType,
	}
}

// Represents a generic constant in type annotations eg. `List[String]`
type GenericConstantNode struct {
	*position.Position
	Constant         ComplexConstantNode
	GenericArguments []ComplexConstantNode
}

// Create a generic constant node eg. `List[String]`
func NewGenericConstantNode(pos *position.Position, constant ComplexConstantNode, args []ComplexConstantNode) *GenericConstantNode {
	return &GenericConstantNode{
		Position:         pos,
		Constant:         constant,
		GenericArguments: args,
	}
}

// Represents a new type definition eg. `typedef StringList = List[String]`
type TypeDefinitionNode struct {
	*position.Position
	Constant ComplexConstantNode // new name of the type
	Type     TypeNode            // the type
}

// Create a type definition node eg. `typedef StringList = List[String]`
func NewTypeDefinitionNode(pos *position.Position, constant ComplexConstantNode, typ TypeNode) *TypeDefinitionNode {
	return &TypeDefinitionNode{
		Position: pos,
		Constant: constant,
		Type:     typ,
	}
}

// Represents a new alias expression eg. `alias push = append`
type AliasExpressionNode struct {
	*position.Position
	NewName IdentifierNode
	OldName IdentifierNode
}

// Create a alias expression node eg. `alias push = append`
func NewAliasExpressionNode(pos *position.Position, newName, oldName IdentifierNode) *AliasExpressionNode {
	return &AliasExpressionNode{
		Position: pos,
		NewName:  newName,
		OldName:  oldName,
	}
}

// Represents an include expression eg. `include Enumerable[V]`
type IncludeExpressionNode struct {
	*position.Position
	Constants []ComplexConstantNode
}

// Create an include expression node eg. `include Enumerable[V]`
func NewIncludeExpressionNode(pos *position.Position, consts []ComplexConstantNode) *IncludeExpressionNode {
	return &IncludeExpressionNode{
		Position:  pos,
		Constants: consts,
	}
}

// Represents an extend expression eg. `extend Enumerable[V]`
type ExtendExpressionNode struct {
	*position.Position
	Constants []ComplexConstantNode
}

// Create an extend expression node eg. `extend Enumerable[V]`
func NewExtendExpressionNode(pos *position.Position, consts []ComplexConstantNode) *ExtendExpressionNode {
	return &ExtendExpressionNode{
		Position:  pos,
		Constants: consts,
	}
}

// Represents an enhance expression eg. `enhance Enumerable[V]`
type EnhanceExpressionNode struct {
	*position.Position
	Constants []ComplexConstantNode
}

// Create an enhance expression node eg. `enhance Enumerable[V]`
func NewEnhanceExpressionNode(pos *position.Position, consts []ComplexConstantNode) *EnhanceExpressionNode {
	return &EnhanceExpressionNode{
		Position:  pos,
		Constants: consts,
	}
}
