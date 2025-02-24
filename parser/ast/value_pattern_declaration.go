package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a value pattern declaration eg. `val [foo, { bar }] = baz()`
type ValuePatternDeclarationNode struct {
	NodeBase
	Pattern     PatternNode
	Initialiser ExpressionNode // value assigned to the value
}

func (*ValuePatternDeclarationNode) IsStatic() bool {
	return false
}

func (*ValuePatternDeclarationNode) Class() *value.Class {
	return value.ValuePatternDeclarationNodeClass
}

func (*ValuePatternDeclarationNode) DirectClass() *value.Class {
	return value.ValuePatternDeclarationNodeClass
}

func (v *ValuePatternDeclarationNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::ValuePatternDeclarationNode{&: %p, pattern: %s, initialiser: %s}",
		v,
		v.Pattern.Inspect(),
		v.Initialiser.Inspect(),
	)
}

func (v *ValuePatternDeclarationNode) Error() string {
	return v.Inspect()
}

// Create a new value declaration node eg. `val foo: String`
func NewValuePatternDeclarationNode(span *position.Span, pattern PatternNode, init ExpressionNode) *ValuePatternDeclarationNode {
	return &ValuePatternDeclarationNode{
		NodeBase:    NodeBase{span: span},
		Pattern:     pattern,
		Initialiser: init,
	}
}
