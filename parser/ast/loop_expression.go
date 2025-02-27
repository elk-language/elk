package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `loop` expression.
type LoopExpressionNode struct {
	TypedNodeBase
	ThenBody []StatementNode // then expression body
}

func (*LoopExpressionNode) IsStatic() bool {
	return false
}

// Create a new `loop` expression node eg. `loop println('elk is awesome')`
func NewLoopExpressionNode(span *position.Span, then []StatementNode) *LoopExpressionNode {
	return &LoopExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ThenBody:      then,
	}
}

func (*LoopExpressionNode) Class() *value.Class {
	return value.LoopExpressionNodeClass
}

func (*LoopExpressionNode) DirectClass() *value.Class {
	return value.LoopExpressionNodeClass
}

func (n *LoopExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::LoopExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  then_body: %%[\n")
	for i, stmt := range n.ThenBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *LoopExpressionNode) Error() string {
	return n.Inspect()
}
