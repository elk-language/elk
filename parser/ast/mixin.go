package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Represents a mixin declaration eg. `mixin Foo; end`
type MixinDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Abstract              bool
	Constant              ExpressionNode      // The constant that will hold the mixin value
	TypeParameters        []TypeParameterNode // Generic type variable definitions
	Body                  []StatementNode     // body of the mixin
	IncludesAndImplements []ExpressionNode
	Bytecode              *vm.BytecodeFunction
}

func (n *MixinDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MixinDeclarationNode)
	if !ok {
		return false
	}

	if n.Abstract != o.Abstract ||
		!n.span.Equal(o.span) ||
		n.comment != o.comment {
		return false
	}

	if n.Constant == o.Constant {
	} else if n.Constant == nil || o.Constant == nil {
		return false
	} else if !n.Constant.Equal(value.Ref(o.Constant)) {
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

	return true
}

func (n *MixinDeclarationNode) String() string {
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

	buff.WriteString("mixin ")
	buff.WriteString(n.Constant.String())

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

	if len(n.Body) == 0 {
		buff.WriteString("; end")
		return buff.String()
	}

	buff.WriteString("\n")
	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteString("\n")
	}
	buff.WriteString("end")

	return buff.String()
}

func (*MixinDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*MixinDeclarationNode) IsStatic() bool {
	return false
}

// Create a new mixin declaration node eg. `mixin Foo; end`
func NewMixinDeclarationNode(
	span *position.Span,
	docComment string,
	abstract bool,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StatementNode,
) *MixinDeclarationNode {

	return &MixinDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Abstract:       abstract,
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

func (*MixinDeclarationNode) Class() *value.Class {
	return value.MixinDeclarationNodeClass
}

func (*MixinDeclarationNode) DirectClass() *value.Class {
	return value.MixinDeclarationNodeClass
}

func (n *MixinDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MixinDeclarationNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	fmt.Fprintf(&buff, ",\n  abstract: %t", n.Abstract)

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  constant: ")
	indent.IndentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  type_parameters: %%[\n")
	for i, element := range n.TypeParameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  body: %%[\n")
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

func (n *MixinDeclarationNode) Error() string {
	return n.Inspect()
}
