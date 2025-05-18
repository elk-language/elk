package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a constant declaration eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
type ConstantDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant    ExpressionNode // name of the constant
	TypeNode    TypeNode       // type of the constant
	Initialiser ExpressionNode // value assigned to the constant
}

func (n *ConstantDeclarationNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	constant := n.Constant.Splice(loc, args, unquote).(ExpressionNode)

	var typeNode TypeNode
	if n.TypeNode != nil {
		typeNode = n.TypeNode.Splice(loc, args, unquote).(TypeNode)
	}

	var init ExpressionNode
	if n.Initialiser != nil {
		init = n.Initialiser.Splice(loc, args, unquote).(ExpressionNode)
	}

	return &ConstantDeclarationNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Constant:               constant,
		TypeNode:               typeNode,
		Initialiser:            init,
	}
}

func (n *ConstantDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Constant.traverse(n, enter, leave) == TraverseBreak {
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

// Check if this node equals another node.
func (n *ConstantDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ConstantDeclarationNode)
	if !ok {
		return false
	}

	if n.DocComment() != o.DocComment() {
		return false
	}

	if !n.Constant.Equal(value.Ref(o.Constant)) {
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

	return n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *ConstantDeclarationNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	buff.WriteString("const ")
	buff.WriteString(n.Constant.String())

	if n.TypeNode != nil {
		buff.WriteString(": ")
		buff.WriteString(n.TypeNode.String())
	}

	if n.Initialiser != nil {
		parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Initialiser)
		initStr := n.Initialiser.String()
		if strings.ContainsRune(initStr, '\n') {
			if parens {
				buff.WriteRune('(')
			}
			buff.WriteRune('\n')
			indent.IndentString(&buff, initStr, 1)
			if parens {
				buff.WriteString("\n)")
			}
		} else {
			if parens {
				buff.WriteRune('(')
			}
			buff.WriteString(initStr)
			if parens {
				buff.WriteRune(')')
			}
		}
	}

	return buff.String()
}

func (*ConstantDeclarationNode) IsStatic() bool {
	return false
}

// Create a new constant declaration node eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
func NewConstantDeclarationNode(loc *position.Location, docComment string, constant ExpressionNode, typ TypeNode, init ExpressionNode) *ConstantDeclarationNode {
	return &ConstantDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:    constant,
		TypeNode:    typ,
		Initialiser: init,
	}
}

func (*ConstantDeclarationNode) Class() *value.Class {
	return value.ConstantDeclarationNodeClass
}

func (*ConstantDeclarationNode) DirectClass() *value.Class {
	return value.ConstantDeclarationNodeClass
}

func (n *ConstantDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ConstantDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  constant: ")
	if n.Constant == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)
	}

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	if n.Initialiser == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)
	}

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ConstantDeclarationNode) Error() string {
	return n.Inspect()
}
