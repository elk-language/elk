package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a closure eg. `|i| -> println(i)`
type ClosureLiteralNode struct {
	TypedNodeBase
	Parameters []ParameterNode // formal parameters of the closure separated by semicolons
	ReturnType TypeNode
	ThrowType  TypeNode
	Body       []StatementNode // body of the closure
}

func (n *ClosureLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ClosureLiteralNode)
	if !ok {
		return false
	}

	if len(n.Parameters) != len(o.Parameters) ||
		len(n.Body) != len(o.Body) {
		return false
	}

	if n.ReturnType == o.ReturnType {
	} else if n.ReturnType == nil || o.ReturnType == nil {
		return false
	} else if !n.ReturnType.Equal(value.Ref(o.ReturnType)) {
		return false
	}

	if n.ThrowType == o.ThrowType {
	} else if n.ThrowType == nil || o.ThrowType == nil {
		return false
	} else if !n.ThrowType.Equal(value.Ref(o.ThrowType)) {
		return false
	}

	for i, param := range n.Parameters {
		if !param.Equal(value.Ref(o.Parameters[i])) {
			return false
		}
	}

	for i, stmt := range n.Body {
		if !stmt.Equal(value.Ref(o.Body[i])) {
			return false
		}
	}

	return n.Span().Equal(o.Span())
}

func (n *ClosureLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("|")
	for i, param := range n.Parameters {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(param.String())
	}
	buff.WriteRune('|')

	if n.ReturnType != nil {
		buff.WriteString(": ")
		buff.WriteString(n.ReturnType.String())
	}

	if n.ThrowType != nil {
		buff.WriteString(" ! ")
		buff.WriteString(n.ThrowType.String())
	}

	buff.WriteString(" ->")

	if len(n.Body) == 1 {
		buff.WriteRune(' ')
		then := n.Body[0]
		parens := ExpressionPrecedence(n) > StatementPrecedence(then)
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteString(then.String())
		if parens {
			buff.WriteRune(')')
		}
	} else {
		buff.WriteRune('\n')
		for _, stmt := range n.Body {
			indent.IndentString(&buff, stmt.String(), 1)
			buff.WriteRune('\n')
		}
		buff.WriteString("end")
	}

	return buff.String()
}

func (*ClosureLiteralNode) IsStatic() bool {
	return false
}

// Create a new closure expression node eg. `|i| -> println(i)`
func NewClosureLiteralNode(span *position.Span, params []ParameterNode, retType TypeNode, throwType TypeNode, body []StatementNode) *ClosureLiteralNode {
	return &ClosureLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Parameters:    params,
		ReturnType:    retType,
		ThrowType:     throwType,
		Body:          body,
	}
}

func (*ClosureLiteralNode) Class() *value.Class {
	return value.ClosureLiteralNodeClass
}

func (*ClosureLiteralNode) DirectClass() *value.Class {
	return value.ClosureLiteralNodeClass
}

func (n *ClosureLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ClosureLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  parameters: %[\n")
	for i, element := range n.Parameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  return_type: ")
	indent.IndentStringFromSecondLine(&buff, n.ReturnType.Inspect(), 1)

	buff.WriteString(",\n  throw_type: ")
	indent.IndentStringFromSecondLine(&buff, n.ThrowType.Inspect(), 1)

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

func (n *ClosureLiteralNode) Error() string {
	return n.Inspect()
}
