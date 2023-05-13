// Package ast defines types
// used by the Elk parser.
//
// All the nodes of the Abstract Syntax Tree
// constructed by the Elk parser are defined in this package.
package ast

import (
	"github.com/elk-language/elk/lexer"
)

// Every node type implements this interface.
type Node interface {
	lexer.Positioner
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

func (*InvalidNode) expressionNode()              {}
func (*ModifierNode) expressionNode()             {}
func (*ModifierIfElseNode) expressionNode()       {}
func (*AssignmentExpressionNode) expressionNode() {}
func (*BinaryExpressionNode) expressionNode()     {}
func (*LogicalExpressionNode) expressionNode()    {}
func (*UnaryExpressionNode) expressionNode()      {}
func (*TrueLiteralNode) expressionNode()          {}
func (*FalseLiteralNode) expressionNode()         {}
func (*NilLiteralNode) expressionNode()           {}
func (*RawStringLiteralNode) expressionNode()     {}
func (*IntLiteralNode) expressionNode()           {}
func (*FloatLiteralNode) expressionNode()         {}
func (*StringLiteralNode) expressionNode()        {}
func (*PublicIdentifierNode) expressionNode()     {}
func (*PrivateIdentifierNode) expressionNode()    {}
func (*PublicConstantNode) expressionNode()       {}
func (*PrivateConstantNode) expressionNode()      {}
func (*SelfLiteralNode) expressionNode()          {}
func (*IfExpressionNode) expressionNode()         {}
func (*UnlessExpressionNode) expressionNode()     {}
func (*WhileExpressionNode) expressionNode()      {}
func (*UntilExpressionNode) expressionNode()      {}
func (*LoopExpressionNode) expressionNode()       {}
func (*BreakExpressionNode) expressionNode()      {}
func (*ReturnExpressionNode) expressionNode()     {}
func (*ContinueExpressionNode) expressionNode()   {}
func (*ThrowExpressionNode) expressionNode()      {}
func (*VariableDeclarationNode) expressionNode()  {}
func (*ConstantLookupNode) expressionNode()       {}

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

// All nodes that should be valid in constant lookups
// should implement this interface.
type ConstantNode interface {
	Node
	TypeNode
	ExpressionNode
	constantNode()
}

func (*InvalidNode) constantNode()         {}
func (*PublicConstantNode) constantNode()  {}
func (*PrivateConstantNode) constantNode() {}
func (*ConstantLookupNode) constantNode()  {}

// Nodes that implement this interface can appear
// inside of a String literal.
type StringLiteralContentNode interface {
	Node
	stringLiteralContentNode()
}

func (*InvalidNode) stringLiteralContentNode()                     {}
func (*StringInterpolationNode) stringLiteralContentNode()         {}
func (*StringLiteralContentSectionNode) stringLiteralContentNode() {}

// Represents a single Elk program (usually a single file).
type ProgramNode struct {
	*lexer.Position
	Body []StatementNode
}

// Create a new program node.
func NewProgramNode(pos *lexer.Position, body []StatementNode) *ProgramNode {
	return &ProgramNode{
		Position: pos,
		Body:     body,
	}
}

// Represents an empty statement eg. a statement with only a semicolon or a newline.
type EmptyStatementNode struct {
	*lexer.Position
}

// Create a new empty statement node.
func NewEmptyStatementNode(pos *lexer.Position) *EmptyStatementNode {
	return &EmptyStatementNode{
		Position: pos,
	}
}

// Expression optionally terminated with a newline or a semicolon.
type ExpressionStatementNode struct {
	*lexer.Position
	Expression ExpressionNode
}

// Create a new expression statement node eg. `5 * 2\n`
func NewExpressionStatementNode(pos *lexer.Position, expr ExpressionNode) *ExpressionStatementNode {
	return &ExpressionStatementNode{
		Position:   pos,
		Expression: expr,
	}
}

// Assignment with the specified operator.
type AssignmentExpressionNode struct {
	*lexer.Position
	Op    *lexer.Token   // operator
	Left  ExpressionNode // left hand side
	Right ExpressionNode // right hand side
}

// Create a new assignment expression node eg. `foo = 3`
func NewAssignmentExpressionNode(pos *lexer.Position, op *lexer.Token, left, right ExpressionNode) *AssignmentExpressionNode {
	return &AssignmentExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Expression of an operator with two operands eg. `2 + 5`, `foo > bar`
type BinaryExpressionNode struct {
	*lexer.Position
	Op    *lexer.Token   // operator
	Left  ExpressionNode // left hand side
	Right ExpressionNode // right hand side
}

// Create a new binary expression node.
func NewBinaryExpressionNode(pos *lexer.Position, op *lexer.Token, left, right ExpressionNode) *BinaryExpressionNode {
	return &BinaryExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Expression of a logical operator with two operands eg. `foo &&  bar`
type LogicalExpressionNode struct {
	*lexer.Position
	Op    *lexer.Token   // operator
	Left  ExpressionNode // left hand side
	Right ExpressionNode // right hand side
}

// Create a new logical expression node.
func NewLogicalExpressionNode(pos *lexer.Position, op *lexer.Token, left, right ExpressionNode) *LogicalExpressionNode {
	return &LogicalExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Expression of an operator with one operand eg. `!foo`, `-bar`
type UnaryExpressionNode struct {
	*lexer.Position
	Op    *lexer.Token   // operator
	Right ExpressionNode // right hand side
}

// Create a new unary expression node.
func NewUnaryExpressionNode(pos *lexer.Position, op *lexer.Token, right ExpressionNode) *UnaryExpressionNode {
	return &UnaryExpressionNode{
		Position: pos,
		Op:       op,
		Right:    right,
	}
}

// `true` literal.
type TrueLiteralNode struct {
	*lexer.Position
}

// Create a new `true` literal node.
func NewTrueLiteralNode(pos *lexer.Position) *TrueLiteralNode {
	return &TrueLiteralNode{
		Position: pos,
	}
}

// `self` literal.
type FalseLiteralNode struct {
	*lexer.Position
}

// Create a new `false` literal node.
func NewFalseLiteralNode(pos *lexer.Position) *FalseLiteralNode {
	return &FalseLiteralNode{
		Position: pos,
	}
}

// `self` literal.
type SelfLiteralNode struct {
	*lexer.Position
}

// Create a new `self` literal node.
func NewSelfLiteralNode(pos *lexer.Position) *SelfLiteralNode {
	return &SelfLiteralNode{
		Position: pos,
	}
}

// `nil` literal.
type NilLiteralNode struct {
	*lexer.Position
}

// Create a new `nil` literal node.
func NewNilLiteralNode(pos *lexer.Position) *NilLiteralNode {
	return &NilLiteralNode{
		Position: pos,
	}
}

// Raw string literal enclosed with single quotes eg. `'foo'`.
type RawStringLiteralNode struct {
	*lexer.Position
	Value string // value of the string literal
}

// Create a new raw string literal node eg. `'foo'`.
func NewRawStringLiteralNode(pos *lexer.Position, val string) *RawStringLiteralNode {
	return &RawStringLiteralNode{
		Position: pos,
		Value:    val,
	}
}

// Int literal eg. `5`, `125_355`, `0xff`
type IntLiteralNode struct {
	*lexer.Position
	Token *lexer.Token
}

// Create a new raw string literal node eg. `5`, `125_355`, `0xff`
func NewIntLiteralNode(pos *lexer.Position, tok *lexer.Token) *IntLiteralNode {
	return &IntLiteralNode{
		Position: pos,
		Token:    tok,
	}
}

// Float literal eg. `5.2`, `.5`, `45e20`
type FloatLiteralNode struct {
	*lexer.Position
	Value string
}

// Create a new float literal node eg. `5.2`, `.5`, `45e20`
func NewFloatLiteralNode(pos *lexer.Position, val string) *FloatLiteralNode {
	return &FloatLiteralNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a syntax error.
type InvalidNode struct {
	*lexer.Position
	Token *lexer.Token
}

// Create a new invalid node.
func NewInvalidNode(pos *lexer.Position, tok *lexer.Token) *InvalidNode {
	return &InvalidNode{
		Position: pos,
		Token:    tok,
	}
}

// Represents a single section of characters of a string literal eg. `foo` in `"foo${bar}"`.
type StringLiteralContentSectionNode struct {
	*lexer.Position
	Value string
}

// Create a new string literal content section node eg. `foo` in `"foo${bar}"`.
func NewStringLiteralContentSectionNode(pos *lexer.Position, val string) *StringLiteralContentSectionNode {
	return &StringLiteralContentSectionNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a single interpolated section of a string literal eg. `bar + 2` in `"foo${bar + 2}"`
type StringInterpolationNode struct {
	*lexer.Position
	Expression ExpressionNode
}

// Create a new string interpolation node eg. `bar + 2` in `"foo${bar + 2}"`
func NewStringInterpolationNode(pos *lexer.Position, expr ExpressionNode) *StringInterpolationNode {
	return &StringInterpolationNode{
		Position:   pos,
		Expression: expr,
	}
}

// Represents a string literal eg. `"foo ${bar} baz"`
type StringLiteralNode struct {
	*lexer.Position
	Content []StringLiteralContentNode
}

// Create a new string literal node eg. `"foo ${bar} baz"`
func NewStringLiteralNode(pos *lexer.Position, cont []StringLiteralContentNode) *StringLiteralNode {
	return &StringLiteralNode{
		Position: pos,
		Content:  cont,
	}
}

// Represents a public identifier eg. `foo`.
type PublicIdentifierNode struct {
	*lexer.Position
	Value string
}

// Create a new public identifier node eg. `foo`.
func NewPublicIdentifierNode(pos *lexer.Position, val string) *PublicIdentifierNode {
	return &PublicIdentifierNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a private identifier eg. `_foo`
type PrivateIdentifierNode struct {
	*lexer.Position
	Value string
}

// Create a new private identifier node eg. `_foo`.
func NewPrivateIdentifierNode(pos *lexer.Position, val string) *PrivateIdentifierNode {
	return &PrivateIdentifierNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a public constant eg. `Foo`.
type PublicConstantNode struct {
	*lexer.Position
	Value string
}

// Create a new public constant node eg. `Foo`.
func NewPublicConstantNode(pos *lexer.Position, val string) *PublicConstantNode {
	return &PublicConstantNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a private constant eg. `_Foo`
type PrivateConstantNode struct {
	*lexer.Position
	Value string
}

// Create a new private constant node eg. `_Foo`.
func NewPrivateConstantNode(pos *lexer.Position, val string) *PrivateConstantNode {
	return &PrivateConstantNode{
		Position: pos,
		Value:    val,
	}
}

// Represents an `if`, `unless`, `while` or `until` modifier expression eg. `return true if foo`.
type ModifierNode struct {
	*lexer.Position
	Modifier *lexer.Token   // modifier token
	Left     ExpressionNode // left hand side
	Right    ExpressionNode // right hand side
}

// Create a new modifier node eg. `return true if foo`.
func NewModifierNode(pos *lexer.Position, mod *lexer.Token, left, right ExpressionNode) *ModifierNode {
	return &ModifierNode{
		Position: pos,
		Modifier: mod,
		Left:     left,
		Right:    right,
	}
}

// Represents an `if .. else` modifier expression eg. `foo = 1 if bar else foo = 2`
type ModifierIfElseNode struct {
	*lexer.Position
	ThenExpression ExpressionNode // then expression body
	Condition      ExpressionNode // if condition
	ElseExpression ExpressionNode // else expression body
}

// Create a new modifier `if` .. `else` node eg. `foo = 1 if bar else foo = 2“.
func NewModifierIfElseNode(pos *lexer.Position, then, cond, els ExpressionNode) *ModifierIfElseNode {
	return &ModifierIfElseNode{
		Position:       pos,
		ThenExpression: then,
		Condition:      cond,
		ElseExpression: els,
	}
}

// Represents an `if` expression eg. `if foo then println("bar")`
type IfExpressionNode struct {
	*lexer.Position
	Condition ExpressionNode  // if condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

// Create a new `if` expression node eg. `if foo then println("bar")`
func NewIfExpressionNode(pos *lexer.Position, cond ExpressionNode, then, els []StatementNode) *IfExpressionNode {
	return &IfExpressionNode{
		Position:  pos,
		ThenBody:  then,
		Condition: cond,
		ElseBody:  els,
	}
}

// Represents an `unless` expression eg. `unless foo then println("bar")`
type UnlessExpressionNode struct {
	*lexer.Position
	Condition ExpressionNode  // unless condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

// Create a new `unless` expression node eg. `unless foo then println("bar")`
func NewUnlessExpressionNode(pos *lexer.Position, cond ExpressionNode, then, els []StatementNode) *UnlessExpressionNode {
	return &UnlessExpressionNode{
		Position:  pos,
		ThenBody:  then,
		Condition: cond,
		ElseBody:  els,
	}
}

// Represents a `while` expression eg. `while i < 5 then i += 5`
type WhileExpressionNode struct {
	*lexer.Position
	Condition ExpressionNode  // while condition
	ThenBody  []StatementNode // then expression body
}

// Create a new `while` expression node eg. `while i < 5 then i += 5`
func NewWhileExpressionNode(pos *lexer.Position, cond ExpressionNode, then []StatementNode) *WhileExpressionNode {
	return &WhileExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  then,
	}
}

// Represents a `until` expression eg. `until i >= 5 then i += 5`
type UntilExpressionNode struct {
	*lexer.Position
	Condition ExpressionNode  // until condition
	ThenBody  []StatementNode // then expression body
}

// Create a new `until` expression node eg. `until i >= 5 then i += 5`
func NewUntilExpressionNode(pos *lexer.Position, cond ExpressionNode, then []StatementNode) *UntilExpressionNode {
	return &UntilExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  then,
	}
}

// Represents a `loop` expression.
type LoopExpressionNode struct {
	*lexer.Position
	ThenBody []StatementNode // then expression body
}

// Create a new `loop` expression node eg. `loop println('elk is awesome')`
func NewLoopExpressionNode(pos *lexer.Position, then []StatementNode) *LoopExpressionNode {
	return &LoopExpressionNode{
		Position: pos,
		ThenBody: then,
	}
}

// Represents a `break` expression eg. `break`
type BreakExpressionNode struct {
	*lexer.Position
}

// Create a new `break` expression node eg. `break`
func NewBreakExpressionNode(pos *lexer.Position) *BreakExpressionNode {
	return &BreakExpressionNode{
		Position: pos,
	}
}

// Represents a `return` expression eg. `return`, `return true`
type ReturnExpressionNode struct {
	*lexer.Position
	Value ExpressionNode
}

// Create a new `return` expression node eg. `return`, `return true`
func NewReturnExpressionNode(pos *lexer.Position, val ExpressionNode) *ReturnExpressionNode {
	return &ReturnExpressionNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a `continue` expression eg. `continue`, `continue "foo"`
type ContinueExpressionNode struct {
	*lexer.Position
	Value ExpressionNode
}

// Create a new `continue` expression node eg. `continue`, `continue "foo"`
func NewContinueExpressionNode(pos *lexer.Position, val ExpressionNode) *ContinueExpressionNode {
	return &ContinueExpressionNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a `throw` expression eg. `throw ArgumentError.new("foo")`
type ThrowExpressionNode struct {
	*lexer.Position
	Value ExpressionNode
}

// Create a new `throw` expression node eg. `throw ArgumentError.new("foo")`
func NewThrowExpressionNode(pos *lexer.Position, val ExpressionNode) *ThrowExpressionNode {
	return &ThrowExpressionNode{
		Position: pos,
		Value:    val,
	}
}

// Represents a variable declaration eg. `var foo: String`
type VariableDeclarationNode struct {
	*lexer.Position
	Name        *lexer.Token   // name of the variable
	Type        TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

// Create a new variable declaration node eg. `var foo: String`
func NewVariableDeclarationNode(pos *lexer.Position, name *lexer.Token, typ TypeNode, init ExpressionNode) *VariableDeclarationNode {
	return &VariableDeclarationNode{
		Position:    pos,
		Name:        name,
		Type:        typ,
		Initialiser: init,
	}
}

// Type expression of an operator with two operands eg. `String | Int`
type BinaryTypeExpressionNode struct {
	*lexer.Position
	Op    *lexer.Token // operator
	Left  TypeNode     // left hand side
	Right TypeNode     // right hand side
}

// Create a new binary type expression node eg. `String | Int`
func NewBinaryTypeExpressionNode(pos *lexer.Position, op *lexer.Token, left, right TypeNode) *BinaryTypeExpressionNode {
	return &BinaryTypeExpressionNode{
		Position: pos,
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Represents an optional or nilable type eg. `String?`
type NilableTypeNode struct {
	*lexer.Position
	Type TypeNode // right hand side
}

// Create a new nilable type node eg. `String?`
func NewNilableTypeNode(pos *lexer.Position, typ TypeNode) *NilableTypeNode {
	return &NilableTypeNode{
		Position: pos,
		Type:     typ,
	}
}

// Represents a constant lookup expressions eg. `Foo::Bar`
type ConstantLookupNode struct {
	*lexer.Position
	Left  ExpressionNode // left hand side
	Right ConstantNode   // right hand side
}

// Create a new constant lookup expression node eg. `Foo::Bar`
func NewConstantLookupNode(pos *lexer.Position, left ExpressionNode, right ConstantNode) *ConstantLookupNode {
	return &ConstantLookupNode{
		Position: pos,
		Left:     left,
		Right:    right,
	}
}
