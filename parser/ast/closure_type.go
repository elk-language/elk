package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a closure type eg. `|i: Int|: String`
type ClosureTypeNode struct {
	TypedNodeBase
	Parameters []ParameterNode // formal parameters of the closure separated by semicolons
	ReturnType TypeNode
	ThrowType  TypeNode
}

func (*ClosureTypeNode) Class() *value.Class {
	return value.ClosureTypeNodeClass
}

func (*ClosureTypeNode) DirectClass() *value.Class {
	return value.ClosureTypeNodeClass
}

func (n *ClosureTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ClosureTypeNode{\n  &: %p", n)

	buff.WriteString(",\n  return_type: ")
	indent.IndentStringFromSecondLine(&buff, n.ReturnType.Inspect(), 1)

	buff.WriteString(",\n  throw_type: ")
	indent.IndentStringFromSecondLine(&buff, n.ThrowType.Inspect(), 1)

	buff.WriteString(",\n  parameters: %%[\n")
	for i, stmt := range n.Parameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ClosureTypeNode) Error() string {
	return n.Inspect()
}

func (*ClosureTypeNode) IsStatic() bool {
	return false
}

// Create a new closure type node eg. `|i: Int|: String`
func NewClosureTypeNode(span *position.Span, params []ParameterNode, retType TypeNode, throwType TypeNode) *ClosureTypeNode {
	return &ClosureTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Parameters:    params,
		ReturnType:    retType,
		ThrowType:     throwType,
	}
}
