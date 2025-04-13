package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `singleton` block expression eg.
//
//	singleton
//		def hello then println("awesome!")
//	end
type SingletonBlockExpressionNode struct {
	TypedNodeBase
	Body     []StatementNode // do expression body
	Bytecode value.Method
}

func (n *SingletonBlockExpressionNode) Splice(loc *position.Location, args *[]Node) Node {
	return &SingletonBlockExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: getLoc(loc, n.loc), typ: n.typ},
		Body:          SpliceSlice(n.Body, loc, args),
		Bytecode:      n.Bytecode,
	}
}

func (n *SingletonBlockExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SingletonBlockExpressionNode)
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

func (n *SingletonBlockExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("singleton\n")

	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteString("\n")
	}

	buff.WriteString("end")

	return buff.String()
}

func (*SingletonBlockExpressionNode) SkipTypechecking() bool {
	return false
}

func (*SingletonBlockExpressionNode) IsStatic() bool {
	return false
}

func (*SingletonBlockExpressionNode) Class() *value.Class {
	return value.DoExpressionNodeClass
}

func (*SingletonBlockExpressionNode) DirectClass() *value.Class {
	return value.DoExpressionNodeClass
}

func (n *SingletonBlockExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SingletonBlockExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *SingletonBlockExpressionNode) Error() string {
	return n.Inspect()
}

// Create a new `singleton` block expression node eg.
//
//	singleton
//		def hello then println("awesome!")
//	end
func NewSingletonBlockExpressionNode(loc *position.Location, body []StatementNode) *SingletonBlockExpressionNode {
	return &SingletonBlockExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Body:          body,
	}
}
