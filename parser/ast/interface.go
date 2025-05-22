package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an interface declaration eg. `interface Foo; end`
type InterfaceDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant       ExpressionNode      // The constant that will hold the interface value
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Body           []StatementNode     // body of the interface
	Implements     []*ImplementExpressionNode
	Bytecode       value.Method
}

func (n *InterfaceDeclarationNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	var constant ExpressionNode
	if n.Constant != nil {
		constant = n.Constant.splice(loc, args, unquote).(ExpressionNode)
	}

	typeParams := SpliceSlice(n.TypeParameters, loc, args, unquote)
	body := SpliceSlice(n.Body, loc, args, unquote)
	implements := SpliceSlice(n.Implements, loc, args, unquote)

	return &InterfaceDeclarationNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Constant:               constant,
		TypeParameters:         typeParams,
		Body:                   body,
		Implements:             implements,
		Bytecode:               n.Bytecode,
	}
}

func (n *InterfaceDeclarationNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::InterfaceDeclarationNode", env)
}

func (n *InterfaceDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Constant != nil {
		if n.Constant.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, param := range n.TypeParameters {
		if param.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, stmt := range n.Body {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *InterfaceDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InterfaceDeclarationNode)
	if !ok {
		return false
	}

	if len(n.TypeParameters) != len(o.TypeParameters) ||
		len(n.Body) != len(o.Body) {
		return false
	}

	if n.Constant == o.Constant {
	} else if n.Constant == nil || o.Constant == nil {
		return false
	} else if !n.Constant.Equal(value.Ref(o.Constant)) {
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

	return n.comment == o.comment &&
		n.loc.Equal(o.loc)
}

func (n *InterfaceDeclarationNode) String() string {
	var buff strings.Builder
	buff.WriteString("interface")
	if n.Constant != nil {
		buff.WriteRune(' ')
		buff.WriteString(n.Constant.String())
	}

	if len(n.TypeParameters) > 0 {
		buff.WriteString("[")
		for i, param := range n.TypeParameters {
			if i > 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(param.String())
		}
		buff.WriteString("]")
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

func (*InterfaceDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*InterfaceDeclarationNode) IsStatic() bool {
	return false
}

// Create a new interface declaration node eg. `interface Foo; end`
func NewInterfaceDeclarationNode(
	loc *position.Location,
	docComment string,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StatementNode,
) *InterfaceDeclarationNode {

	return &InterfaceDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

func (*InterfaceDeclarationNode) Class() *value.Class {
	return value.InterfaceDeclarationNodeClass
}

func (*InterfaceDeclarationNode) DirectClass() *value.Class {
	return value.InterfaceDeclarationNodeClass
}

func (n *InterfaceDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::InterfaceDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *InterfaceDeclarationNode) Error() string {
	return n.Inspect()
}
