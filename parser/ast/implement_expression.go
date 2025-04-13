package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an enhance expression eg. `implement Enumerable[V]`
type ImplementExpressionNode struct {
	TypedNodeBase
	Constants []ComplexConstantNode
}

func (n *ImplementExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ImplementExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Constants:     SpliceSlice(n.Constants, loc, args, unquote),
	}
}

// Check if this node equals another node.
func (n *ImplementExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ImplementExpressionNode)
	if !ok {
		return false
	}

	if len(n.Constants) != len(o.Constants) {
		return false
	}

	for i, constant := range n.Constants {
		if !constant.Equal(value.Ref(o.Constants[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *ImplementExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("implement ")

	for i, constant := range n.Constants {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(constant.String())
	}

	return buff.String()
}

func (*ImplementExpressionNode) SkipTypechecking() bool {
	return false
}

func (*ImplementExpressionNode) IsStatic() bool {
	return false
}

// Create an enhance expression node eg. `implement Enumerable[V]`
func NewImplementExpressionNode(loc *position.Location, consts []ComplexConstantNode) *ImplementExpressionNode {
	return &ImplementExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Constants:     consts,
	}
}

func (*ImplementExpressionNode) Class() *value.Class {
	return value.ImplementExpressionNodeClass
}

func (*ImplementExpressionNode) DirectClass() *value.Class {
	return value.ImplementExpressionNodeClass
}

func (n *ImplementExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ImplementExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  constants: %[\n")
	for i, element := range n.Constants {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ImplementExpressionNode) Error() string {
	return n.Inspect()
}
