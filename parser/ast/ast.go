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

// returns true if the value is known at compile-time
func IsStatic(expr ExpressionNode) bool {
	switch n := expr.(type) {
	case *TrueLiteralNode, *FalseLiteralNode, *NilLiteralNode, *RawStringLiteralNode,
		*IntLiteralNode, *FloatLiteralNode, *DoubleQuotedStringLiteralNode, *ClosureExpressionNode,
		*SimpleSymbolLiteralNode:
		return true
	case *NamedValueLiteralNode:
		return IsStatic(n.Value)
	case *KeyValueExpressionNode:
		return IsStatic(n.Key) && IsStatic(n.Value)
	case *SymbolKeyValueExpressionNode:
		return IsStatic(n.Value)
	case *ListLiteralNode:
		for _, element := range n.Elements {
			if !IsStatic(element) {
				return false
			}
		}
		return true
	case *TupleLiteralNode:
		for _, element := range n.Elements {
			if !IsStatic(element) {
				return false
			}
		}
		return true
	case *SetLiteralNode:
		for _, element := range n.Elements {
			if !IsStatic(element) {
				return false
			}
		}
		return true
	case *MapLiteralNode:
		for _, element := range n.Elements {
			if !IsStatic(element) {
				return false
			}
		}
		return true
	default:
		return false
	}
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

// All expression nodes implement this interface.
type ExpressionNode interface {
	Node
	expressionNode()
}

func (*InvalidNode) expressionNode()                   {}
func (*ModifierNode) expressionNode()                  {}
func (*ModifierIfElseNode) expressionNode()            {}
func (*ModifierForInNode) expressionNode()             {}
func (*AssignmentExpressionNode) expressionNode()      {}
func (*BinaryExpressionNode) expressionNode()          {}
func (*LogicalExpressionNode) expressionNode()         {}
func (*UnaryExpressionNode) expressionNode()           {}
func (*TrueLiteralNode) expressionNode()               {}
func (*FalseLiteralNode) expressionNode()              {}
func (*NilLiteralNode) expressionNode()                {}
func (*SimpleSymbolLiteralNode) expressionNode()       {}
func (*InterpolatedSymbolLiteral) expressionNode()     {}
func (*NamedValueLiteralNode) expressionNode()         {}
func (*IntLiteralNode) expressionNode()                {}
func (*FloatLiteralNode) expressionNode()              {}
func (*RawStringLiteralNode) expressionNode()          {}
func (*DoubleQuotedStringLiteralNode) expressionNode() {}
func (*InterpolatedStringLiteralNode) expressionNode() {}
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
func (*ForExpressionNode) expressionNode()             {}
func (*BreakExpressionNode) expressionNode()           {}
func (*ReturnExpressionNode) expressionNode()          {}
func (*ContinueExpressionNode) expressionNode()        {}
func (*ThrowExpressionNode) expressionNode()           {}
func (*VariableDeclarationNode) expressionNode()       {}
func (*ConstantDeclarationNode) expressionNode()       {}
func (*ConstantLookupNode) expressionNode()            {}
func (*ClosureExpressionNode) expressionNode()         {}
func (*ClassDeclarationNode) expressionNode()          {}
func (*ModuleDeclarationNode) expressionNode()         {}
func (*MixinDeclarationNode) expressionNode()          {}
func (*InterfaceDeclarationNode) expressionNode()      {}
func (*StructDeclarationNode) expressionNode()         {}
func (*MethodDefinitionNode) expressionNode()          {}
func (*InitDefinitionNode) expressionNode()            {}
func (*MethodSignatureDefinitionNode) expressionNode() {}
func (*GenericConstantNode) expressionNode()           {}
func (*TypeDefinitionNode) expressionNode()            {}
func (*AliasExpressionNode) expressionNode()           {}
func (*IncludeExpressionNode) expressionNode()         {}
func (*ExtendExpressionNode) expressionNode()          {}
func (*EnhanceExpressionNode) expressionNode()         {}
func (*ConstructorCallNode) expressionNode()           {}
func (*MethodCallNode) expressionNode()                {}
func (*FunctionCallNode) expressionNode()              {}
func (*KeyValueExpressionNode) expressionNode()        {}
func (*SymbolKeyValueExpressionNode) expressionNode()  {}
func (*ListLiteralNode) expressionNode()               {}
func (*TupleLiteralNode) expressionNode()              {}
func (*SetLiteralNode) expressionNode()                {}
func (*MapLiteralNode) expressionNode()                {}
func (*RangeLiteralNode) expressionNode()              {}
func (*TypeLiteralNode) expressionNode()               {}

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
}

func (*InvalidNode) parameterNode()            {}
func (*FormalParameterNode) parameterNode()    {}
func (*MethodParameterNode) parameterNode()    {}
func (*SignatureParameterNode) parameterNode() {}
func (*LoopParameterNode) parameterNode()      {}

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

// Same as [NewExpressionStatementNode] but returns an interface
func NewExpressionStatementNodeI(pos *position.Position, expr ExpressionNode) StatementNode {
	return &ExpressionStatementNode{
		Position:   pos,
		Expression: expr,
	}
}

// Formal parameter optionally terminated with a newline or a semicolon.
type ParameterStatementNode struct {
	*position.Position
	Parameter ParameterNode
}

// Create a new formal parameter statement node eg. `foo: Bar\n`
func NewParameterStatementNode(pos *position.Position, param ParameterNode) *ParameterStatementNode {
	return &ParameterStatementNode{
		Position:  pos,
		Parameter: param,
	}
}

// Same as [NewParameterStatementNode] but returns an interface
func NewParameterStatementNodeI(pos *position.Position, param ParameterNode) StructBodyStatementNode {
	return &ParameterStatementNode{
		Position:  pos,
		Parameter: param,
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

// Same as [NewBinaryExpressionNode] but returns an interface
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

// Same as [NewLogicalExpressionNode] but returns an interface
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

// Represents an interpolated string literal eg. `"foo ${bar} baz"`
type InterpolatedStringLiteralNode struct {
	*position.Position
	Content []StringLiteralContentNode
}

// Create a new interpolated string literal node eg. `"foo ${bar} baz"`
func NewInterpolatedStringLiteralNode(pos *position.Position, cont []StringLiteralContentNode) *InterpolatedStringLiteralNode {
	return &InterpolatedStringLiteralNode{
		Position: pos,
		Content:  cont,
	}
}

// Represents a simple double quoted string literal eg. `"foo baz"`
type DoubleQuotedStringLiteralNode struct {
	*position.Position
	Value string
}

// Create a new double quoted string literal node eg. `"foo baz"`
func NewDoubleQuotedStringLiteralNode(pos *position.Position, val string) *DoubleQuotedStringLiteralNode {
	return &DoubleQuotedStringLiteralNode{
		Position: pos,
		Value:    val,
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

// Represents an `for .. in` modifier expression eg. `println(i) for i in 10..30`
type ModifierForInNode struct {
	*position.Position
	ThenExpression ExpressionNode  // then expression body
	Parameters     []ParameterNode // list of parameters
	InExpression   ExpressionNode  // expression that will be iterated through
}

// Create a new modifier `for` .. `in` node eg. `println(i) for i in 10..30`
func NewModifierForInNode(pos *position.Position, then ExpressionNode, params []ParameterNode, in ExpressionNode) *ModifierForInNode {
	return &ModifierForInNode{
		Position:       pos,
		ThenExpression: then,
		Parameters:     params,
		InExpression:   in,
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

// Represents a `for` expression eg. `for i in 5..15 then println(i)`
type ForExpressionNode struct {
	*position.Position
	Parameters   []ParameterNode // list of parameters
	InExpression ExpressionNode  // expression that will be iterated through
	ThenBody     []StatementNode // then expression body
}

// Create a new `for` expression node eg. `for i in 5..15 then println(i)`
func NewForExpressionNode(pos *position.Position, params []ParameterNode, inExpr ExpressionNode, then []StatementNode) *ForExpressionNode {
	return &ForExpressionNode{
		Position:     pos,
		Parameters:   params,
		InExpression: inExpr,
		ThenBody:     then,
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

// Same as [NewBinaryTypeExpressionNode] but returns an interface
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

// Represents a formal parameter in closure or struct declarations eg. `foo: String = 'bar'`
type FormalParameterNode struct {
	*position.Position
	Name        string         // name of the variable
	Type        TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

// Create a new formal parameter node eg. `foo: String = 'bar'`
func NewFormalParameterNode(pos *position.Position, name string, typ TypeNode, init ExpressionNode) *FormalParameterNode {
	return &FormalParameterNode{
		Position:    pos,
		Name:        name,
		Type:        typ,
		Initialiser: init,
	}
}

// Represents a formal parameter in method declarations eg. `foo: String = 'bar'`
type MethodParameterNode struct {
	*position.Position
	Name                string         // name of the variable
	SetInstanceVariable bool           // whether an instance variable with this name gets automatically assigned
	Type                TypeNode       // type of the variable
	Initialiser         ExpressionNode // value assigned to the variable
}

// Create a new formal parameter node eg. `foo: String = 'bar'`
func NewMethodParameterNode(pos *position.Position, name string, setIvar bool, typ TypeNode, init ExpressionNode) *MethodParameterNode {
	return &MethodParameterNode{
		Position:            pos,
		SetInstanceVariable: setIvar,
		Name:                name,
		Type:                typ,
		Initialiser:         init,
	}
}

// Represents a signature parameter in method and closure signatures eg. `foo?: String`
type SignatureParameterNode struct {
	*position.Position
	Name     string   // name of the variable
	Type     TypeNode // type of the variable
	Optional bool     // whether this parameter is optional
}

// Create a new signature parameter node eg. `foo?: String`
func NewSignatureParameterNode(pos *position.Position, name string, typ TypeNode, opt bool) *SignatureParameterNode {
	return &SignatureParameterNode{
		Position: pos,
		Name:     name,
		Type:     typ,
		Optional: opt,
	}
}

// Represents a parameter in loop expressions eg. `foo: String`
type LoopParameterNode struct {
	*position.Position
	Name string   // name of the variable
	Type TypeNode // type of the variable
}

// Create a new loop parameter node eg. `foo: String`
func NewLoopParameterNode(pos *position.Position, name string, typ TypeNode) *LoopParameterNode {
	return &LoopParameterNode{
		Position: pos,
		Name:     name,
		Type:     typ,
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

// Represents an interface declaration eg. `interface Foo; end`
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

// Represents a struct declaration eg. `struct Foo; end`
type StructDeclarationNode struct {
	*position.Position
	Constant      ExpressionNode            // The constant that will hold the struct object
	TypeVariables []TypeVariableNode        // Generic type variable definitions
	Body          []StructBodyStatementNode // body of the struct
}

// Create a new struct declaration node eg. `struct Foo; end`
func NewStructDeclarationNode(
	pos *position.Position,
	constant ExpressionNode,
	typeVars []TypeVariableNode,
	body []StructBodyStatementNode,
) *StructDeclarationNode {

	return &StructDeclarationNode{
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

// Represents a symbol literal with simple content eg. `:foo`, `:'foo bar`, `:"lol"`
type SimpleSymbolLiteralNode struct {
	*position.Position
	Content string
}

// Create a simple symbol literal node eg. `:foo`, `:'foo bar`, `:"lol"`
func NewSimpleSymbolLiteralNode(pos *position.Position, cont string) *SimpleSymbolLiteralNode {
	return &SimpleSymbolLiteralNode{
		Position: pos,
		Content:  cont,
	}
}

// Represents an interpolated symbol eg. `:"foo ${bar + 2}"`
type InterpolatedSymbolLiteral struct {
	*position.Position
	Content *InterpolatedStringLiteralNode
}

// Create an interpolated symbol literal node eg. `:"foo ${bar + 2}"`
func NewInterpolatedSymbolLiteral(pos *position.Position, cont *InterpolatedStringLiteralNode) *InterpolatedSymbolLiteral {
	return &InterpolatedSymbolLiteral{
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

// Represents a constructor definition eg. `init then 'hello world'`
type InitDefinitionNode struct {
	*position.Position
	Parameters []ParameterNode // formal parameters
	ThrowType  TypeNode
	Body       []StatementNode // body of the method
}

// Create a constructor definition node eg. `init then 'hello world'`
func NewInitDefinitionNode(pos *position.Position, params []ParameterNode, throwType TypeNode, body []StatementNode) *InitDefinitionNode {
	return &InitDefinitionNode{
		Position:   pos,
		Parameters: params,
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

// Represents a new alias expression eg. `alias push append`
type AliasExpressionNode struct {
	*position.Position
	NewName string
	OldName string
}

// Create a alias expression node eg. `alias push append`
func NewAliasExpressionNode(pos *position.Position, newName, oldName string) *AliasExpressionNode {
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

// Represents a named argument in a function call eg. `foo: 123`
type NamedCallArgumentNode struct {
	*position.Position
	Name  string
	Value ExpressionNode
}

// Create a named argument node eg. `foo: 123`
func NewNamedCallArgumentNode(pos *position.Position, name string, val ExpressionNode) *NamedCallArgumentNode {
	return &NamedCallArgumentNode{
		Position: pos,
		Name:     name,
		Value:    val,
	}
}

// Represents a constructor call eg. `String(123)`
type ConstructorCallNode struct {
	*position.Position
	Class               ComplexConstantNode // class that is being instantiated
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

// Create a constructor call node eg. `String(123)`
func NewConstructorCallNode(pos *position.Position, class ComplexConstantNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ConstructorCallNode {
	return &ConstructorCallNode{
		Position:            pos,
		Class:               class,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a method call eg. `'123'.to_int()`
type MethodCallNode struct {
	*position.Position
	Receiver            ExpressionNode
	NilSafe             bool
	MethodName          string
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

// Create a method call node eg. `'123'.to_int()`
func NewMethodCallNode(pos *position.Position, recv ExpressionNode, nilSafe bool, methodName string, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *MethodCallNode {
	return &MethodCallNode{
		Position:            pos,
		Receiver:            recv,
		NilSafe:             nilSafe,
		MethodName:          methodName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a function-like call eg. `to_string(123)`
type FunctionCallNode struct {
	*position.Position
	MethodName          string
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

// Create a function call node eg. `to_string(123)`
func NewFunctionCallNode(pos *position.Position, methodName string, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *FunctionCallNode {
	return &FunctionCallNode{
		Position:            pos,
		MethodName:          methodName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a symbol value expression eg. `foo: bar`
type SymbolKeyValueExpressionNode struct {
	*position.Position
	Key   string
	Value ExpressionNode
}

// Create a symbol key value node eg. `foo: bar`
func NewSymbolKeyValueExpressionNode(pos *position.Position, key string, val ExpressionNode) *SymbolKeyValueExpressionNode {
	return &SymbolKeyValueExpressionNode{
		Position: pos,
		Key:      key,
		Value:    val,
	}
}

// Represents a key value expression eg. `foo => bar`
type KeyValueExpressionNode struct {
	*position.Position
	Key   ExpressionNode
	Value ExpressionNode
}

// Create a key value expression node eg. `foo => bar`
func NewKeyValueExpressionNode(pos *position.Position, key, val ExpressionNode) *KeyValueExpressionNode {
	return &KeyValueExpressionNode{
		Position: pos,
		Key:      key,
		Value:    val,
	}
}

// Represents a List literal eg. `[1, 5, -6]`
type ListLiteralNode struct {
	*position.Position
	Elements []ExpressionNode
}

// Create a List literal node eg. `[1, 5, -6]`
func NewListLiteralNode(pos *position.Position, elements []ExpressionNode) *ListLiteralNode {
	return &ListLiteralNode{
		Position: pos,
		Elements: elements,
	}
}

// Same as [NewListLiteralNode] but returns an interface
func NewListLiteralNodeI(pos *position.Position, elements []ExpressionNode) ExpressionNode {
	return &ListLiteralNode{
		Position: pos,
		Elements: elements,
	}
}

// Represents a Tuple literal eg. `%(1, 5, -6)`
type TupleLiteralNode struct {
	*position.Position
	Elements []ExpressionNode
}

// Create a Tuple literal node eg. `%(1, 5, -6)`
func NewTupleLiteralNode(pos *position.Position, elements []ExpressionNode) *TupleLiteralNode {
	return &TupleLiteralNode{
		Position: pos,
		Elements: elements,
	}
}

// Same as [NewTupleLiteralNode] but returns an interface
func NewTupleLiteralNodeI(pos *position.Position, elements []ExpressionNode) ExpressionNode {
	return &TupleLiteralNode{
		Position: pos,
		Elements: elements,
	}
}

// Represents a Set literal eg. `%{1, 5, -6}`
type SetLiteralNode struct {
	*position.Position
	Elements []ExpressionNode
}

// Create a Set literal node eg. `%{1, 5, -6}`
func NewSetLiteralNode(pos *position.Position, elements []ExpressionNode) *SetLiteralNode {
	return &SetLiteralNode{
		Position: pos,
		Elements: elements,
	}
}

// Same as [NewSetLiteralNode] but returns an interface
func NewSetLiteralNodeI(pos *position.Position, elements []ExpressionNode) ExpressionNode {
	return &SetLiteralNode{
		Position: pos,
		Elements: elements,
	}
}

// Represents a Map literal eg. `{ foo: 1, 'bar' => 5, baz }`
type MapLiteralNode struct {
	*position.Position
	Elements []ExpressionNode
}

// Create a Map literal node eg. `{ foo: 1, 'bar' => 5, baz }`
func NewMapLiteralNode(pos *position.Position, elements []ExpressionNode) *MapLiteralNode {
	return &MapLiteralNode{
		Position: pos,
		Elements: elements,
	}
}

// Same as [NewMapLiteralNode] but returns an interface
func NewMapLiteralNodeI(pos *position.Position, elements []ExpressionNode) ExpressionNode {
	return &MapLiteralNode{
		Position: pos,
		Elements: elements,
	}
}

// Represents a Range literal eg. `1..5`
type RangeLiteralNode struct {
	*position.Position
	Exclusive bool
	Left      ExpressionNode
	Right     ExpressionNode
}

// Create a Range literal node eg. `1..5`
func NewRangeLiteralNode(pos *position.Position, exclusive bool, left, right ExpressionNode) *RangeLiteralNode {
	return &RangeLiteralNode{
		Position:  pos,
		Exclusive: exclusive,
		Left:      left,
		Right:     right,
	}
}

// Represents a Type literal eg. `type String?`
type TypeLiteralNode struct {
	*position.Position
	TypeExpression TypeNode
}

// Create a Type literal node eg. `type String?`
func NewTypeLiteralNode(pos *position.Position, texpr TypeNode) *TypeLiteralNode {
	return &TypeLiteralNode{
		Position:       pos,
		TypeExpression: texpr,
	}
}
