package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a value declaration eg. `val foo: String`
type ValueDeclarationNode struct {
	TypedNodeBase
	Name        string         // name of the value
	TypeNode    TypeNode       // type of the value
	Initialiser ExpressionNode // value assigned to the value
}

func (n *ValueDeclarationNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	var typeNode TypeNode
	if n.TypeNode != nil {
		typeNode = n.TypeNode.Splice(loc, args, unquote).(TypeNode)
	}

	var init ExpressionNode
	if n.Initialiser != nil {
		init = n.Initialiser.Splice(loc, args, unquote).(ExpressionNode)
	}

	return &ValueDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Name:          n.Name,
		TypeNode:      typeNode,
		Initialiser:   init,
	}
}

func (n *ValueDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.TypeNode != nil {
		if n.TypeNode.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	if n.Initialiser != nil {
		if n.Initialiser.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *ValueDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ValueDeclarationNode)
	if !ok {
		return false
	}

	if n.Name != o.Name ||
		!n.loc.Equal(o.loc) {
		return false
	}

	if n.TypeNode == o.TypeNode {
	} else if n.TypeNode == nil || o.TypeNode == nil {
		return false
	} else if !n.TypeNode.Equal(value.Ref(o.TypeNode)) {
		return false
	}

	if n.Initialiser == o.Initialiser {
	} else if n.Initialiser == nil || o.Initialiser == nil {
		return false
	} else if !n.Initialiser.Equal(value.Ref(o.Initialiser)) {
		return false
	}

	return true
}

func (n *ValueDeclarationNode) String() string {
	var buff strings.Builder

	buff.WriteString("val ")
	buff.WriteString(n.Name)

	if n.TypeNode != nil {
		buff.WriteString(": ")
		buff.WriteString(n.TypeNode.String())
	}

	if n.Initialiser != nil {
		buff.WriteString(" = ")
		buff.WriteString(n.Initialiser.String())
	}

	return buff.String()
}

func (*ValueDeclarationNode) IsStatic() bool {
	return false
}

// Create a new value declaration node eg. `val foo: String`
func NewValueDeclarationNode(loc *position.Location, name string, typ TypeNode, init ExpressionNode) *ValueDeclarationNode {
	return &ValueDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Name:          name,
		TypeNode:      typ,
		Initialiser:   init,
	}
}

func (*ValueDeclarationNode) Class() *value.Class {
	return value.ValueDeclarationNodeClass
}

func (*ValueDeclarationNode) DirectClass() *value.Class {
	return value.ValueDeclarationNodeClass
}

func (n *ValueDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ValueDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  name: ")
	buff.WriteString(n.Name)

	buff.WriteString(",\n  type_node: ")
	if n.TypeNode == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)
	}

	buff.WriteString(",\n  initialiser: ")
	if n.Initialiser == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)
	}

	buff.WriteString("\n}")

	return buff.String()
}

func (v *ValueDeclarationNode) Error() string {
	return v.Inspect()
}
