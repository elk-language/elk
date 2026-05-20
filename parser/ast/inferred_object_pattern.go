package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an inferred object pattern eg. `@{foo: 5, bar: a, c}`
type InferredObjectPatternNode struct {
	TypedNodeBase
	Attributes []PatternNode
}

func (n *InferredObjectPatternNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &InferredObjectPatternNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Attributes:    SpliceSlice(n.Attributes, loc, args, unquote),
	}
}

func (n *InferredObjectPatternNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::InferredObjectPatternNode", env)
}

func (n *InferredObjectPatternNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, attr := range n.Attributes {
		if attr.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *InferredObjectPatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ObjectPatternNode)
	if !ok {
		return false
	}

	if len(n.Attributes) != len(o.Attributes) {
		return false
	}

	for i, attr := range n.Attributes {
		if !attr.Equal(value.Ref(o.Attributes[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

func (n *InferredObjectPatternNode) String() string {
	var buff strings.Builder

	buff.WriteString("@{")

	for i, attr := range n.Attributes {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(attr.String())
	}

	buff.WriteRune('}')

	return buff.String()
}

func (m *InferredObjectPatternNode) IsStatic() bool {
	return false
}

// Create an inferred object pattern node eg. `@{foo: 5, bar: a, c}`
func NewInferredObjectPatternNode(loc *position.Location, attrs []PatternNode) *InferredObjectPatternNode {
	return &InferredObjectPatternNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Attributes:    attrs,
	}
}

func (*InferredObjectPatternNode) Class() *value.Class {
	return value.InferredObjectPatternNodeClass
}

func (*InferredObjectPatternNode) DirectClass() *value.Class {
	return value.InferredObjectPatternNodeClass
}

func (n *InferredObjectPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::InferredObjectPatternNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  attributes: %[")
	if len(n.Attributes) > 0 {
		buff.WriteRune('\n')
		for i, element := range n.Attributes {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, element.Inspect(), 2)
		}
		buff.WriteString("\n  ")
	}
	buff.WriteRune(']')

	buff.WriteString("\n}")

	return buff.String()
}

func (n *InferredObjectPatternNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *InferredObjectPatternNode) Error() string {
	return n.Inspect()
}
