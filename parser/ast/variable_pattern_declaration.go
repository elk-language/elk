package ast

import (
	"fmt"
	"strings"

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

func (n *VariablePatternDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::VariablePatternDeclarationNode{\n  &: %p", n)

	buff.WriteString(",\n  pattern: ")
	indentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	indentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
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
