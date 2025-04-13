package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a module declaration eg. `module Foo; end`
type ModuleDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant ExpressionNode  // The constant that will hold the module value
	Body     []StatementNode // body of the module
	Bytecode value.Method
}

func (n *ModuleDeclarationNode) Splice(loc *position.Location, args *[]Node) Node {
	var constant ExpressionNode
	if n.Constant != nil {
		constant = n.Constant.Splice(loc, args).(ExpressionNode)
	}

	body := SpliceSlice(n.Body, loc, args)

	return &ModuleDeclarationNode{
		TypedNodeBase:          n.TypedNodeBase,
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Constant:               constant,
		Body:                   body,
		Bytecode:               n.Bytecode,
	}
}

func (n *ModuleDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ModuleDeclarationNode)
	if !ok {
		return false
	}

	if !n.loc.Equal(o.loc) ||
		n.comment != o.comment ||
		len(n.Body) != len(o.Body) {
		return false
	}

	if n.Constant == o.Constant {
	} else if n.Constant == nil || o.Constant == nil {
		return false
	} else if !n.Constant.Equal(value.Ref(o.Constant)) {
		return false
	}

	for i, stmt := range n.Body {
		if !stmt.Equal(value.Ref(o.Body[i])) {
			return false
		}
	}

	return true
}

func (n *ModuleDeclarationNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	buff.WriteString("module")
	if n.Constant != nil {
		buff.WriteRune(' ')
		buff.WriteString(n.Constant.String())
	}

	if len(n.Body) > 0 {
		buff.WriteString("\n")
		for _, stmt := range n.Body {
			indent.IndentString(&buff, stmt.String(), 1)
			buff.WriteString("\n")
		}
		buff.WriteString("end")
	} else {
		buff.WriteString("; end")
	}

	return buff.String()
}

func (*ModuleDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*ModuleDeclarationNode) IsStatic() bool {
	return false
}

// Create a new module declaration node eg. `module Foo; end`
func NewModuleDeclarationNode(
	loc *position.Location,
	docComment string,
	constant ExpressionNode,
	body []StatementNode,
) *ModuleDeclarationNode {

	return &ModuleDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant: constant,
		Body:     body,
	}
}

func (*ModuleDeclarationNode) Class() *value.Class {
	return value.ClassDeclarationNodeClass
}

func (*ModuleDeclarationNode) DirectClass() *value.Class {
	return value.ClassDeclarationNodeClass
}

func (n *ModuleDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ModuleDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  constant: ")
	if n.Constant == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)
	}

	buff.WriteString(",\n  body: %[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ModuleDeclarationNode) Error() string {
	return n.Inspect()
}
