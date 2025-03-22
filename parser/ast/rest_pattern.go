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

	fmt.Fprintf(&buff, "Std::Elk::AST::RestPatternNode{\n  &: %p", n)

	buff.WriteString(",\n  identifier: ")
	indent.IndentStringFromSecondLine(&buff, n.Identifier.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *RestPatternNode) Error() string {
	return n.Inspect()
}
