package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `while` expression eg. `while i < 5 then i += 5`
type WhileExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // while condition
	ThenBody  []StatementNode // then expression body
}

func (n *WhileExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*WhileExpressionNode)
	if !ok {
		return false
	}

	if !n.Condition.Equal(value.Ref(o.Condition)) ||
		!n.span.Equal(o.span) ||
		len(n.ThenBody) != len(o.ThenBody) {
		return false
	}

	for i, stmt := range n.ThenBody {
		if !stmt.Equal(value.Ref(o.ThenBody[i])) {
			return false
		}
	}

	return true
}

func (n *WhileExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("while ")
	buff.WriteString(n.Condition.String())

	buff.WriteRune('\n')
	for _, stmt := range n.ThenBody {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}
	buff.WriteString("end")

	return buff.String()
}

func (*WhileExpressionNode) IsStatic() bool {
	return false
}

// Create a new `while` expression node eg. `while i < 5 then i += 5`
func NewWhileExpressionNode(span *position.Span, cond ExpressionNode, then []StatementNode) *WhileExpressionNode {
	return &WhileExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Condition:     cond,
		ThenBody:      then,
	}
}

func (*WhileExpressionNode) Class() *value.Class {
	return value.WhileExpressionNodeClass
}

func (*WhileExpressionNode) DirectClass() *value.Class {
	return value.WhileExpressionNodeClass
}

func (n *WhileExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::WhileExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  condition: ")
	indent.IndentStringFromSecondLine(&buff, n.Condition.Inspect(), 1)

	buff.WriteString(",\n  then_body: %%[\n")
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

func (n *WhileExpressionNode) Error() string {
	return n.Inspect()
}
