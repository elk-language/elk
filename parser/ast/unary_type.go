package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Type of an operator with one operand eg. `-2`, `+3`
type UnaryTypeNode struct {
	TypedNodeBase
	Op       *token.Token // operator
	TypeNode TypeNode     // right hand side
}

func (n *UnaryTypeNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UnaryTypeNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Op:            n.Op.Splice(loc, unquote),
		TypeNode:      n.TypeNode.splice(loc, args, unquote).(TypeNode),
	}
}

func (n *UnaryTypeNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UnaryTypeNode", env)
}

func (n *UnaryTypeNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.TypeNode.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *UnaryTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UnaryTypeNode)
	if !ok {
		return false
	}

	return n.Op.Equal(o.Op) &&
		n.TypeNode.Equal(value.Ref(o.TypeNode)) &&
		n.loc.Equal(o.loc)
}

func (n *UnaryTypeNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Op.FetchValue())

	parens := TypePrecedence(n) > TypePrecedence(n.TypeNode)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.TypeNode.String())
	if parens {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (u *UnaryTypeNode) IsStatic() bool {
	return false
}

// Create a new unary expression node.
func NewUnaryTypeNode(loc *position.Location, op *token.Token, typeNode TypeNode) *UnaryTypeNode {
	return &UnaryTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Op:            op,
		TypeNode:      typeNode,
	}
}

func (*UnaryTypeNode) Class() *value.Class {
	return value.UnaryTypeNodeClass
}

func (*UnaryTypeNode) DirectClass() *value.Class {
	return value.UnaryTypeNodeClass
}

func (n *UnaryTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UnaryTypeNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UnaryTypeNode) Error() string {
	return n.Inspect()
}
