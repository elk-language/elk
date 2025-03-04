package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Represents a `singleton` block expression eg.
//
//	singleton
//		def hello then println("awesome!")
//	end
type SingletonBlockExpressionNode struct {
	TypedNodeBase
	Body     []StatementNode // do expression body
	Bytecode *vm.BytecodeFunction
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

	fmt.Fprintf(&buff, "Std::Elk::AST::SingletonBlockExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  body: %%[\n")
	for i, stmt := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
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
func NewSingletonBlockExpressionNode(span *position.Span, body []StatementNode) *SingletonBlockExpressionNode {
	return &SingletonBlockExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
	}
}
