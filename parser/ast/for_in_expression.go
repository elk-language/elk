package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `for in` expression eg. `for i in 5..15 then println(i)`
type ForInExpressionNode struct {
	TypedNodeBase
	Pattern      PatternNode
	InExpression ExpressionNode  // expression that will be iterated through
	ThenBody     []StatementNode // then expression body
}

func (n *ForInExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ForInExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Pattern:       n.Pattern.Splice(loc, args, unquote).(PatternNode),
		InExpression:  n.InExpression.Splice(loc, args, unquote).(ExpressionNode),
		ThenBody:      SpliceSlice(n.ThenBody, loc, args, unquote),
	}
}

func (n *ForInExpressionNode) Traverse(yield func(Node) bool) bool {
	if !n.Pattern.Traverse(yield) {
		return false
	}
	if !n.InExpression.Traverse(yield) {
		return false
	}
	for _, stmt := range n.ThenBody {
		if !stmt.Traverse(yield) {
			return false
		}
	}
	return yield(n)
}

func (n *ForInExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ForInExpressionNode)
	if !ok {
		return false
	}

	if !n.loc.Equal(o.loc) ||
		!n.Pattern.Equal(value.Ref(o.Pattern)) ||
		!n.InExpression.Equal(value.Ref(o.InExpression)) ||
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

func (n *ForInExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("for ")
	buff.WriteString(n.Pattern.String())
	buff.WriteString(" in ")
	buff.WriteString(n.InExpression.String())
	buff.WriteRune('\n')

	for _, stmt := range n.ThenBody {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteString("\n")
	}
	buff.WriteString("end")

	return buff.String()
}

func (*ForInExpressionNode) IsStatic() bool {
	return false
}

// Create a new `for in` expression node eg. `for i in 5..15 then println(i)`
func NewForInExpressionNode(loc *position.Location, pattern PatternNode, inExpr ExpressionNode, then []StatementNode) *ForInExpressionNode {
	return &ForInExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Pattern:       pattern,
		InExpression:  inExpr,
		ThenBody:      then,
	}
}

func (*ForInExpressionNode) Class() *value.Class {
	return value.ForInExpressionNodeClass
}

func (*ForInExpressionNode) DirectClass() *value.Class {
	return value.ForInExpressionNodeClass
}

func (n *ForInExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ForInExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  in_expression: ")
	indent.IndentStringFromSecondLine(&buff, n.InExpression.Inspect(), 1)

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

func (n *ForInExpressionNode) Error() string {
	return n.Inspect()
}
