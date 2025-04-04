package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a rest element in a list pattern eg. `*a`
type RestPatternNode struct {
	NodeBase
	Identifier IdentifierNode
}

func (n *RestPatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*RestPatternNode)
	if !ok {
		return false
	}

	return n.Identifier.Equal(value.Ref(o.Identifier)) &&
		n.span.Equal(o.span)
}

func (n *RestPatternNode) String() string {
	var buff strings.Builder

	buff.WriteRune('*')
	buff.WriteString(n.Identifier.String())

	return buff.String()
}

func (r *RestPatternNode) IsStatic() bool {
	return false
}

// Create a rest pattern node eg. `*a`
func NewRestPatternNode(span *position.Span, ident IdentifierNode) *RestPatternNode {
	return &RestPatternNode{
		NodeBase:   NodeBase{span: span},
		Identifier: ident,
	}
}

func (*RestPatternNode) Class() *value.Class {
	return value.RestPatternNodeClass
}

func (*RestPatternNode) DirectClass() *value.Class {
	return value.RestPatternNodeClass
}

func (n *RestPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::RestPatternNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  identifier: ")
	indent.IndentStringFromSecondLine(&buff, n.Identifier.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *RestPatternNode) Error() string {
	return n.Inspect()
}
