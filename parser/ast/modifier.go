package ast

import (
	"fmt"

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

func (m *ModifierNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::ModifierNode{&: %p, modifier: %s, left: %s, right: %s}",
		m,
		m.Modifier.Inspect(),
		m.Left.Inspect(),
		m.Right.Inspect(),
	)
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

func (m *ModifierIfElseNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::ModifierIfElseNode{&: %p, then_expression: %s, condition: %s, else_expression: %s}",
		m,
		m.ThenExpression.Inspect(),
		m.Condition.Inspect(),
		m.ElseExpression.Inspect(),
	)
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

func (m *ModifierForInNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::ModifierForInNode{&: %p, then_expression: %s, pattern: %s, in_expression: %s}",
		m,
		m.ThenExpression.Inspect(),
		m.Pattern.Inspect(),
		m.InExpression.Inspect(),
	)
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
