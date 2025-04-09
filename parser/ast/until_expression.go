package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `until` expression eg. `until i >= 5 then i += 5`
type UntilExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // until condition
	ThenBody  []StatementNode // then expression body
}

func (n *UntilExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UntilExpressionNode)
	if !ok {
		return false
	}

	if len(n.ThenBody) != len(o.ThenBody) ||
		!n.Condition.Equal(value.Ref(o.Condition)) ||
		!n.loc.Equal(o.loc) {
		return false
	}

	for i, stmt := range n.ThenBody {
		if !stmt.Equal(value.Ref(o.ThenBody[i])) {
			return false
		}
	}

	return true
}

func (n *UntilExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("until ")
	buff.WriteString(n.Condition.String())

	buff.WriteRune('\n')
	for _, stmt := range n.ThenBody {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}

	buff.WriteString("end")

	return buff.String()
}

func (*UntilExpressionNode) IsStatic() bool {
	return false
}

// Create a new `until` expression node eg. `until i >= 5 then i += 5`
func NewUntilExpressionNode(loc *position.Location, cond ExpressionNode, then []StatementNode) *UntilExpressionNode {
	return &UntilExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Condition:     cond,
		ThenBody:      then,
	}
}

func (*UntilExpressionNode) Class() *value.Class {
	return value.UntilExpressionNodeClass
}

func (*UntilExpressionNode) DirectClass() *value.Class {
	return value.UntilExpressionNodeClass
}

func (n *UntilExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UntilExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  condition: ")
	indent.IndentStringFromSecondLine(&buff, n.Condition.Inspect(), 1)

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

func (n *UntilExpressionNode) Error() string {
	return n.Inspect()
}
