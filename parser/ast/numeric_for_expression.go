package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a numeric `for` expression eg. `fornum i := 0; i < 10; i += 1 then println(i)`
type NumericForExpressionNode struct {
	TypedNodeBase
	Initialiser ExpressionNode  // i := 0
	Condition   ExpressionNode  // i < 10
	Increment   ExpressionNode  // i += 1
	ThenBody    []StatementNode // then expression body
}

func (n *NumericForExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NumericForExpressionNode)
	if !ok {
		return false
	}

	if n.Initialiser == o.Initialiser {
	} else if n.Initialiser == nil || o.Initialiser == nil {
		return false
	} else if !n.Initialiser.Equal(value.Ref(o.Initialiser)) {
		return false
	}

	if n.Condition == o.Condition {
	} else if n.Condition == nil || o.Condition == nil {
		return false
	} else if !n.Condition.Equal(value.Ref(o.Condition)) {
		return false
	}

	if n.Increment == o.Increment {
	} else if n.Increment == nil || o.Increment == nil {
		return false
	} else if !n.Increment.Equal(value.Ref(o.Increment)) {
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

	return n.loc.Equal(o.loc)
}

func (n *NumericForExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("fornum ")

	if n.Initialiser != nil {
		buff.WriteString(n.Initialiser.String())
	}

	buff.WriteRune(';')

	if n.Condition != nil {
		buff.WriteRune(' ')
		buff.WriteString(n.Condition.String())
	}

	buff.WriteRune(';')

	if n.Increment != nil {
		buff.WriteRune(' ')
		buff.WriteString(n.Increment.String())
	}

	buff.WriteRune('\n')
	for _, stmt := range n.ThenBody {
		buff.WriteString(stmt.String())
		buff.WriteRune('\n')
	}
	buff.WriteString("end")

	return buff.String()
}

func (*NumericForExpressionNode) IsStatic() bool {
	return false
}

// Create a new numeric `fornum` expression eg. `for i := 0; i < 10; i += 1 then println(i)`
func NewNumericForExpressionNode(loc *position.Location, init, cond, incr ExpressionNode, then []StatementNode) *NumericForExpressionNode {
	return &NumericForExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Initialiser:   init,
		Condition:     cond,
		Increment:     incr,
		ThenBody:      then,
	}
}

func (*NumericForExpressionNode) Class() *value.Class {
	return value.NumericForExpressionNodeClass
}

func (*NumericForExpressionNode) DirectClass() *value.Class {
	return value.NumericForExpressionNodeClass
}

func (n *NumericForExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NumericForExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  initialiser: ")
	indent.IndentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString(",\n  condition: ")
	indent.IndentStringFromSecondLine(&buff, n.Condition.Inspect(), 1)

	buff.WriteString(",\n  increment: ")
	indent.IndentStringFromSecondLine(&buff, n.Increment.Inspect(), 1)

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

func (n *NumericForExpressionNode) Error() string {
	return n.Inspect()
}
