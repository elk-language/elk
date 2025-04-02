package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an `if` expression eg. `if foo then println("bar")`
type IfExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // if condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

// Check if this node equals another node.
func (n *IfExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*IfExpressionNode)
	if !ok {
		return false
	}

	if !n.Condition.Equal(value.Ref(o.Condition)) ||
		!n.span.Equal(o.span) {
		return false
	}

	if len(n.ThenBody) != len(o.ThenBody) ||
		len(n.ElseBody) != len(o.ElseBody) {
		return false
	}

	for i, element := range n.ThenBody {
		if !element.Equal(value.Ref(o.ThenBody[i])) {
			return false
		}
	}

	for i, element := range n.ElseBody {
		if !element.Equal(value.Ref(o.ElseBody[i])) {
			return false
		}
	}

	return true
}

// Return a string representation of the node.
func (n *IfExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("if ")
	buff.WriteString(n.Condition.String())
	buff.WriteRune('\n')

	for _, stmt := range n.ThenBody {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}

	if len(n.ElseBody) > 0 {
		then := n.ElseBody[0]
		parens := ExpressionPrecedence(n) > StatementPrecedence(then)
		if len(n.ElseBody) == 1 && !parens {
			buff.WriteString("else ")
			buff.WriteString(then.String())
		} else {
			buff.WriteString("else\n")
			for _, stmt := range n.ElseBody {
				indent.IndentString(&buff, stmt.String(), 1)
				buff.WriteRune('\n')
			}
			buff.WriteString("end")
		}
	} else {
		buff.WriteString("end")
	}

	return buff.String()
}

func (*IfExpressionNode) IsStatic() bool {
	return false
}

// Create a new `if` expression node eg. `if foo then println("bar")`
func NewIfExpressionNode(span *position.Span, cond ExpressionNode, then, els []StatementNode) *IfExpressionNode {
	return &IfExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ThenBody:      then,
		Condition:     cond,
		ElseBody:      els,
	}
}

func (*IfExpressionNode) Class() *value.Class {
	return value.IfExpressionNodeClass
}

func (*IfExpressionNode) DirectClass() *value.Class {
	return value.IfExpressionNodeClass
}

func (n *IfExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::IfExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  condition: ")
	indent.IndentStringFromSecondLine(&buff, n.Condition.Inspect(), 1)

	buff.WriteString(",\n  then_body: %[\n")
	for i, element := range n.ThenBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  else_body: %[\n")
	for i, element := range n.ElseBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *IfExpressionNode) Error() string {
	return n.Inspect()
}
