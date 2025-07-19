package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a variable declaration eg. `var foo: String`
type VariableDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name        IdentifierNode // name of the variable
	TypeNode    TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

func (n *VariableDeclarationNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	name := n.Name.splice(loc, args, unquote).(IdentifierNode)

	var typeNode TypeNode
	if n.TypeNode != nil {
		typeNode = n.TypeNode.splice(loc, args, unquote).(TypeNode)
	}

	var init ExpressionNode
	if n.Initialiser != nil {
		init = n.Initialiser.splice(loc, args, unquote).(ExpressionNode)
	}

	return &VariableDeclarationNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Name:                   name,
		TypeNode:               typeNode,
		Initialiser:            init,
	}
}

func (n *VariableDeclarationNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::VariableDeclarationNode", env)
}

func (n *VariableDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Name.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
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

func (n *VariableDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*VariableDeclarationNode)
	if !ok {
		return false
	}

	if !n.Name.Equal(value.Ref(o.Name)) ||
		n.comment != o.comment ||
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

func (n *VariableDeclarationNode) String() string {
	var buff strings.Builder

	buff.WriteString("var ")
	buff.WriteString(n.Name.String())

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

func (*VariableDeclarationNode) IsStatic() bool {
	return false
}

func (*VariableDeclarationNode) Class() *value.Class {
	return value.VariableDeclarationNodeClass
}

func (*VariableDeclarationNode) DirectClass() *value.Class {
	return value.VariableDeclarationNodeClass
}

func (n *VariableDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::VariableDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  name: ")
	indent.IndentStringFromSecondLine(&buff, n.Name.Inspect(), 1)

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

func (v *VariableDeclarationNode) Error() string {
	return v.Inspect()
}

// Create a new variable declaration node eg. `var foo: String`
func NewVariableDeclarationNode(loc *position.Location, docComment string, name IdentifierNode, typ TypeNode, init ExpressionNode) *VariableDeclarationNode {
	return &VariableDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Name:        name,
		TypeNode:    typ,
		Initialiser: init,
	}
}
