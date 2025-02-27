package ast

import (
	"fmt"
	"strings"

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

	fmt.Fprintf(&buff, "Std::AST::ClosureLiteralNode{\n  &: %p", n)

	buff.WriteString(",\n  parameters: %%[\n")
	for i, element := range n.Parameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  return_type: ")
	indentStringFromSecondLine(&buff, n.ReturnType.Inspect(), 1)

	buff.WriteString(",\n  throw_type: ")
	indentStringFromSecondLine(&buff, n.ThrowType.Inspect(), 1)

	buff.WriteString(",\n  body: %%[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ClosureLiteralNode) Error() string {
	return n.Inspect()
}
