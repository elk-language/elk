package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Union type eg. `String | Int | Float`
type UnionTypeNode struct {
	TypedNodeBase
	Elements []TypeNode
}

func (n *UnionTypeNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UnionTypeNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Elements:      SpliceSlice(n.Elements, loc, args, unquote),
	}
}

func (n *UnionTypeNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UnionTypeNode", env)
}

func (n *UnionTypeNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

func (n *UnionTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UnionTypeNode)
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

func (n *UnionTypeNode) String() string {
	var buff strings.Builder

	for i, element := range n.Elements {
		if i > 0 {
			buff.WriteString(" | ")
		}
		buff.WriteString(element.String())
	}

	return buff.String()
}

func (*UnionTypeNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String | Int`
func NewUnionTypeNode(loc *position.Location, elements []TypeNode) *UnionTypeNode {
	return &UnionTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Elements:      elements,
	}
}

func (*UnionTypeNode) Class() *value.Class {
	return value.UnionTypeNodeClass
}

func (*UnionTypeNode) DirectClass() *value.Class {
	return value.UnionTypeNodeClass
}

func (n *UnionTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UnionTypeNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *UnionTypeNode) Error() string {
	return n.Inspect()
}
