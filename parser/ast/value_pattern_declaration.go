package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a value pattern declaration eg. `val [foo, { bar }] = baz()`
type ValuePatternDeclarationNode struct {
	NodeBase
	Pattern     PatternNode
	Initialiser ExpressionNode // value assigned to the value
}

func (n *ValuePatternDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ValuePatternDeclarationNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
		n.Pattern.Equal(value.Ref(o.Pattern)) &&
		n.Initialiser.Equal(value.Ref(o.Initialiser))
}

func (n *ValuePatternDeclarationNode) String() string {
	var buff strings.Builder

	buff.WriteString("val ")
	buff.WriteString(n.Pattern.String())
	buff.WriteString(" = ")
	buff.WriteString(n.Initialiser.String())

	return buff.String()
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

func (n *ValuePatternDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ValuePatternDeclarationNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	indent.IndentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
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
