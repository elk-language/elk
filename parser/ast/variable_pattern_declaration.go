package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a variable declaration with patterns eg. `var [foo, { bar }] = baz()`
type VariablePatternDeclarationNode struct {
	NodeBase
	Pattern     PatternNode
	Initialiser ExpressionNode // value assigned to the variable
}

func (*VariablePatternDeclarationNode) IsStatic() bool {
	return false
}

func (*VariablePatternDeclarationNode) Class() *value.Class {
	return value.VariablePatternDeclarationNodeClass
}

func (*VariablePatternDeclarationNode) DirectClass() *value.Class {
	return value.VariablePatternDeclarationNodeClass
}

func (v *VariablePatternDeclarationNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::VariablePatternDeclarationNode{&: %p, pattern: %s, initialiser: %s}",
		v,
		v.Pattern.Inspect(),
		v.Initialiser.Inspect(),
	)
}

func (v *VariablePatternDeclarationNode) Error() string {
	return v.Inspect()
}

// Create a new variable declaration node with patterns eg. `var [foo, { bar }] = baz()`
func NewVariablePatternDeclarationNode(span *position.Span, pattern PatternNode, init ExpressionNode) *VariablePatternDeclarationNode {
	return &VariablePatternDeclarationNode{
		NodeBase:    NodeBase{span: span},
		Pattern:     pattern,
		Initialiser: init,
	}
}
