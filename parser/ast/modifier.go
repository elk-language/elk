package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Represents an `if`, `unless`, `while` or `until` modifier expression eg. `return true if foo`.
type ModifierNode struct {
	TypedNodeBase
	Modifier *token.Token   // modifier token
	Left     ExpressionNode // left hand side
	Right    ExpressionNode // right hand side
}

func (*ModifierNode) IsStatic() bool {
	return false
}

func (*ModifierNode) Class() *value.Class {
	return value.ModifierNodeClass
}

func (*ModifierNode) DirectClass() *value.Class {
	return value.ModifierNodeClass
}

func (n *ModifierNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::ModifierNode{\n  &: %p", n)

	buff.WriteString(",\n  modifier: ")
	indentStringFromSecondLine(&buff, n.Modifier.Inspect(), 1)

	buff.WriteString(",\n  left: ")
	indentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (m *ModifierNode) Error() string {
	return m.Inspect()
}

// Create a new modifier node eg. `return true if foo`.
func NewModifierNode(span *position.Span, mod *token.Token, left, right ExpressionNode) *ModifierNode {
	return &ModifierNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Modifier:      mod,
		Left:          left,
		Right:         right,
	}
}

// Represents an `if .. else` modifier expression eg. `foo = 1 if bar else foo = 2`
type ModifierIfElseNode struct {
	TypedNodeBase
	ThenExpression ExpressionNode // then expression body
	Condition      ExpressionNode // if condition
	ElseExpression ExpressionNode // else expression body
}

func (*ModifierIfElseNode) IsStatic() bool {
	return false
}

func (*ModifierIfElseNode) Class() *value.Class {
	return value.ModifierIfElseNodeClass
}

func (*ModifierIfElseNode) DirectClass() *value.Class {
	return value.ModifierIfElseNodeClass
}

func (n *ModifierIfElseNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::ModifierIfElseNode{\n  &: %p", n)

	buff.WriteString(",\n  then_expression: ")
	indentStringFromSecondLine(&buff, n.ThenExpression.Inspect(), 1)

	buff.WriteString(",\n  condition: ")
	indentStringFromSecondLine(&buff, n.Condition.Inspect(), 1)

	buff.WriteString(",\n  else_expression: ")
	indentStringFromSecondLine(&buff, n.ElseExpression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (m *ModifierIfElseNode) Error() string {
	return m.Inspect()
}

// Create a new modifier `if` .. `else` node eg. `foo = 1 if bar else foo = 2“.
func NewModifierIfElseNode(span *position.Span, then, cond, els ExpressionNode) *ModifierIfElseNode {
	return &ModifierIfElseNode{
		TypedNodeBase:  TypedNodeBase{span: span},
		ThenExpression: then,
		Condition:      cond,
		ElseExpression: els,
	}
}

// Represents an `for .. in` modifier expression eg. `println(i) for i in 10..30`
type ModifierForInNode struct {
	TypedNodeBase
	ThenExpression ExpressionNode // then expression body
	Pattern        PatternNode
	InExpression   ExpressionNode // expression that will be iterated through
}

func (*ModifierForInNode) IsStatic() bool {
	return false
}

func (*ModifierForInNode) Class() *value.Class {
	return value.ModifierForInNodeClass
}

func (*ModifierForInNode) DirectClass() *value.Class {
	return value.ModifierForInNodeClass
}

func (n *ModifierForInNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::ModifierForInNode{\n  &: %p", n)

	buff.WriteString(",\n  then_expression: ")
	indentStringFromSecondLine(&buff, n.ThenExpression.Inspect(), 1)

	buff.WriteString(",\n  pattern: ")
	indentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  in_expression: ")
	indentStringFromSecondLine(&buff, n.InExpression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (m *ModifierForInNode) Error() string {
	return m.Inspect()
}

// Create a new modifier `for` .. `in` node eg. `println(i) for i in 10..30`
func NewModifierForInNode(span *position.Span, then ExpressionNode, pattern PatternNode, in ExpressionNode) *ModifierForInNode {
	return &ModifierForInNode{
		TypedNodeBase:  TypedNodeBase{span: span},
		ThenExpression: then,
		Pattern:        pattern,
		InExpression:   in,
	}
}
