package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an as pattern eg. `> 5 && < 20 as foo`
type AsPatternNode struct {
	NodeBase
	Pattern PatternNode
	Name    IdentifierNode
}

func (*AsPatternNode) IsStatic() bool {
	return false
}

// Create an Object pattern node eg. `Foo(foo: 5, bar: a, c)`
func NewAsPatternNode(span *position.Span, pattern PatternNode, name IdentifierNode) *AsPatternNode {
	return &AsPatternNode{
		NodeBase: NodeBase{span: span},
		Pattern:  pattern,
		Name:     name,
	}
}

func (*AsPatternNode) Class() *value.Class {
	return value.AsPatternNodeClass
}

func (*AsPatternNode) DirectClass() *value.Class {
	return value.AsPatternNodeClass
}

func (n *AsPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AsPatternNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  name: ")
	indent.IndentStringFromSecondLine(&buff, n.Name.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AsPatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AsPatternNode)
	if !ok {
		return false
	}

	if !n.Span().Equal(o.Span()) {
		return false
	}

	if !n.Pattern.Equal(value.Ref(o.Pattern)) {
		return false
	}

	if !n.Name.Equal(value.Ref(o.Name)) {
		return false
	}

	return true
}

func (n *AsPatternNode) Error() string {
	return n.Inspect()
}

func (n *AsPatternNode) String() string {
	var buff strings.Builder

	leftParen := PatternPrecedence(n) > PatternPrecedence(n.Pattern)
	if leftParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Pattern.String())
	if leftParen {
		buff.WriteRune(')')
	}

	buff.WriteString(" as ")

	buff.WriteString(n.Name.String())

	return buff.String()
}
