package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `go` expression eg. `go foo()`, `go; foo(); end`
type GoExpressionNode struct {
	TypedNodeBase
	Body []StatementNode
}

func (n *GoExpressionNode) Splice(loc *position.Location, args *[]Node) Node {
	return &GoExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: getLoc(loc, n.loc), typ: n.typ},
		Body:          SpliceSlice(n.Body, loc, args),
	}
}

// Check if this node equals another node.
func (n *GoExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*GoExpressionNode)
	if !ok {
		return false
	}

	if len(n.Body) != len(o.Body) {
		return false
	}

	for i, stmt := range n.Body {
		if !stmt.Equal(value.Ref(o.Body[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *GoExpressionNode) String() string {
	var buff strings.Builder

	if len(n.Body) == 1 {
		buff.WriteString("go ")

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
		buff.WriteString("go\n")
		for _, stmt := range n.Body {
			indent.IndentString(&buff, stmt.String(), 1)
			buff.WriteString("\n")
		}
		buff.WriteString("end")
	}

	return buff.String()
}

func (*GoExpressionNode) IsStatic() bool {
	return false
}

// Create a new `go` expression node eg. `go foo()`, `go; foo(); end`
func NewGoExpressionNode(loc *position.Location, body []StatementNode) *GoExpressionNode {
	return &GoExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Body:          body,
	}
}

func (*GoExpressionNode) Class() *value.Class {
	return value.GoExpressionNodeClass
}

func (*GoExpressionNode) DirectClass() *value.Class {
	return value.GoExpressionNodeClass
}

func (n *GoExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::GoExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  body: %[\n")
	for i, stmt := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *GoExpressionNode) Error() string {
	return n.Inspect()
}
