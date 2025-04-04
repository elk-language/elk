package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `loop` expression.
type LoopExpressionNode struct {
	TypedNodeBase
	ThenBody []StatementNode // then expression body
}

func (n *LoopExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*LoopExpressionNode)
	if !ok {
		return false
	}

	if len(n.ThenBody) != len(o.ThenBody) {
		return false
	}

	for i, stmt := range n.ThenBody {
		if !stmt.Equal(value.Ref(o.ThenBody[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

func (n *LoopExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("loop\n")

	for _, stmt := range n.ThenBody {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteString("\n")
	}

	buff.WriteString("end")

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::LoopExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  then_body: %[\n")
	for i, stmt := range n.ThenBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *LoopExpressionNode) Error() string {
	return n.Inspect()
}
