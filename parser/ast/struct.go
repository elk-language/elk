package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a single statement of a struct body
// optionally terminated with a newline or semicolon.
type StructBodyStatementNode interface {
	Node
	structBodyStatementNode()
}

func (*InvalidNode) structBodyStatementNode()            {}
func (*EmptyStatementNode) structBodyStatementNode()     {}
func (*ParameterStatementNode) structBodyStatementNode() {}

// Represents a struct declaration eg. `struct Foo; end`
type StructDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant       ExpressionNode            // The constant that will hold the struct value
	TypeParameters []TypeParameterNode       // Generic type variable definitions
	Body           []StructBodyStatementNode // body of the struct
}

func (n *StructDeclarationNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	var constant ExpressionNode
	if n.Constant != nil {
		constant = n.Constant.splice(loc, args, unquote).(ExpressionNode)
	}

	typeParams := SpliceSlice(n.TypeParameters, loc, args, unquote)
	body := SpliceSlice(n.Body, loc, args, unquote)

	return &StructDeclarationNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Constant:               constant,
		TypeParameters:         typeParams,
		Body:                   body,
	}
}

func (n *StructDeclarationNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::StructDeclarationNode", env)
}

func (n *StructDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

func (n *StructDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*StructDeclarationNode)
	if !ok {
		return false
	}

	if !n.Constant.Equal(value.Ref(o.Constant)) ||
		!n.loc.Equal(o.loc) ||
		n.comment != o.comment ||
		len(n.TypeParameters) != len(o.TypeParameters) ||
		len(n.Body) != len(o.Body) {
		return false
	}

	for i, typeParam := range n.TypeParameters {
		if !typeParam.Equal(value.Ref(o.TypeParameters[i])) {
			return false
		}
	}

	for i, bodyStmt := range n.Body {
		if !bodyStmt.Equal(value.Ref(o.Body[i])) {
			return false
		}
	}

	return true
}

func (n *StructDeclarationNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	buff.WriteString("struct")
	if n.Constant != nil {
		buff.WriteRune(' ')
		buff.WriteString(n.Constant.String())
	}

	if len(n.TypeParameters) > 0 {
		buff.WriteRune('[')
		for i, typeParam := range n.TypeParameters {
			if i > 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(typeParam.String())
		}
		buff.WriteRune(']')
	}

	if len(n.Body) > 0 {
		buff.WriteRune('\n')
		for _, stmt := range n.Body {
			indent.IndentString(&buff, stmt.String(), 1)
			buff.WriteRune('\n')
		}
		buff.WriteString("end")
	} else {
		buff.WriteString("; end")
	}

	return buff.String()
}

func (*StructDeclarationNode) IsStatic() bool {
	return false
}

// Create a new struct declaration node eg. `struct Foo; end`
func NewStructDeclarationNode(
	loc *position.Location,
	docComment string,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StructBodyStatementNode,
) *StructDeclarationNode {

	return &StructDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

func (*StructDeclarationNode) Class() *value.Class {
	return value.StructDeclarationNodeClass
}

func (*StructDeclarationNode) DirectClass() *value.Class {
	return value.StructDeclarationNodeClass
}

func (n *StructDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::StructDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *StructDeclarationNode) Error() string {
	return n.Inspect()
}
