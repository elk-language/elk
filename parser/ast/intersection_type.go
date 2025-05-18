package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Union type eg. `String & Int & Float`
type IntersectionTypeNode struct {
	TypedNodeBase
	Elements []TypeNode
}

func (n *IntersectionTypeNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &IntersectionTypeNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Elements:      SpliceSlice(n.Elements, loc, args, unquote),
	}
}

func (n *IntersectionTypeNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, elem := range n.Elements {
		if elem.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *IntersectionTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*IntersectionTypeNode)
	if !ok {
		return false
	}

	if len(n.Elements) != len(o.Elements) {
		return false
	}

	for i, element := range n.Elements {
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

func (n *IntersectionTypeNode) String() string {
	var buff strings.Builder

	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(" & ")
		}
		buff.WriteString(element.String())
	}

	return buff.String()
}

func (*IntersectionTypeNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String & Int`
func NewIntersectionTypeNode(loc *position.Location, elements []TypeNode) *IntersectionTypeNode {
	return &IntersectionTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Elements:      elements,
	}
}

func (*IntersectionTypeNode) Class() *value.Class {
	return value.IntersectionTypeNodeClass
}

func (*IntersectionTypeNode) DirectClass() *value.Class {
	return value.IntersectionTypeNodeClass
}

func (n *IntersectionTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::IntersectionTypeNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  elements: %[\n")
	for i, stmt := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *IntersectionTypeNode) Error() string {
	return n.Inspect()
}
