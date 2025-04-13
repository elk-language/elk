package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a class declaration eg. `class Foo; end`
type ClassDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Abstract       bool
	Sealed         bool
	Primitive      bool
	NoInit         bool
	Constant       ExpressionNode      // The constant that will hold the class value
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Superclass     ExpressionNode      // the super/parent class of this class
	Body           []StatementNode     // body of the class
	Bytecode       value.Method
}

func (n *ClassDeclarationNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	var constant ExpressionNode
	if n.Constant != nil {
		constant = n.Constant.Splice(loc, args, unquote).(ExpressionNode)
	}

	typeParams := SpliceSlice(n.TypeParameters, loc, args, unquote)

	var superclass ExpressionNode
	if n.Superclass != nil {
		superclass = n.Superclass.Splice(loc, args, unquote).(ExpressionNode)
	}

	body := SpliceSlice(n.Body, loc, args, unquote)

	return &ClassDeclarationNode{
		TypedNodeBase:  TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Abstract:       n.Abstract,
		Sealed:         n.Sealed,
		Primitive:      n.Primitive,
		NoInit:         n.NoInit,
		Constant:       constant,
		TypeParameters: typeParams,
		Superclass:     superclass,
		Body:           body,
		Bytecode:       n.Bytecode,
	}
}

func (n *ClassDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ClassDeclarationNode)
	if !ok {
		return false
	}

	if n.Abstract != o.Abstract ||
		n.Sealed != o.Sealed ||
		n.Primitive != o.Primitive ||
		n.NoInit != o.NoInit {
		return false
	}

	if n.Constant == o.Constant {
	} else if n.Constant == nil || o.Constant == nil {
		return false
	} else if !n.Constant.Equal(value.Ref(o.Constant)) {
		return false
	}

	if n.Superclass == o.Superclass {
	} else if n.Superclass == nil || o.Superclass == nil {
		return false
	} else if !n.Superclass.Equal(value.Ref(o.Superclass)) {
		return false
	}

	if len(n.TypeParameters) != len(o.TypeParameters) ||
		len(n.Body) != len(o.Body) {
		return false
	}

	for i, param := range n.TypeParameters {
		if !param.Equal(value.Ref(o.TypeParameters[i])) {
			return false
		}
	}

	for i, stmt := range n.Body {
		if !stmt.Equal(value.Ref(o.Body[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

func (n *ClassDeclarationNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	if n.Abstract {
		buff.WriteString("abstract ")
	}

	if n.Sealed {
		buff.WriteString("sealed ")
	}

	if n.Primitive {
		buff.WriteString("primitive ")
	}

	if n.NoInit {
		buff.WriteString("noinit ")
	}

	buff.WriteString("class")
	if n.Constant != nil {
		buff.WriteRune(' ')
		buff.WriteString(n.Constant.String())
	}

	if len(n.TypeParameters) > 0 {
		buff.WriteRune('[')
		for i, param := range n.TypeParameters {
			if i > 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(param.String())
		}
		buff.WriteRune(']')
	}

	if n.Superclass != nil {
		buff.WriteString(" < ")
		buff.WriteString(n.Superclass.String())
	}

	if len(n.Body) == 0 {
		buff.WriteString("; end")
		return buff.String()
	}

	buff.WriteRune('\n')
	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}

	buff.WriteString("end")

	return buff.String()
}

func (*ClassDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*ClassDeclarationNode) IsStatic() bool {
	return false
}

// Create a new class declaration node eg. `class Foo; end`
func NewClassDeclarationNode(
	loc *position.Location,
	docComment string,
	abstract bool,
	sealed bool,
	primitive bool,
	noinit bool,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	superclass ExpressionNode,
	body []StatementNode,
) *ClassDeclarationNode {

	return &ClassDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Abstract:       abstract,
		Sealed:         sealed,
		Primitive:      primitive,
		NoInit:         noinit,
		Constant:       constant,
		TypeParameters: typeParams,
		Superclass:     superclass,
		Body:           body,
	}
}

func (*ClassDeclarationNode) Class() *value.Class {
	return value.ClassDeclarationNodeClass
}

func (*ClassDeclarationNode) DirectClass() *value.Class {
	return value.ClassDeclarationNodeClass
}

func (n *ClassDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ClassDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	fmt.Fprintf(&buff, ",\n  abstract: %t", n.Abstract)
	fmt.Fprintf(&buff, ",\n  sealed: %t", n.Sealed)
	fmt.Fprintf(&buff, ",\n  primitive: %t", n.Primitive)
	fmt.Fprintf(&buff, ",\n  noinit: %t", n.NoInit)

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  constant: ")
	if n.Constant == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)
	}

	buff.WriteString(",\n  type_parameters: %[\n")
	for i, element := range n.TypeParameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  superclass: ")
	if n.Superclass == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Superclass.Inspect(), 1)
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

func (n *ClassDeclarationNode) Error() string {
	return n.Inspect()
}
