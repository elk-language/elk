package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `do` expression eg.
//
//	do
//		print("awesome!")
//	end
type DoExpressionNode struct {
	TypedNodeBase
	Body    []StatementNode // do expression body
	Catches []*CatchNode
	Finally []StatementNode
}

func (*DoExpressionNode) IsStatic() bool {
	return false
}

func (*DoExpressionNode) Class() *value.Class {
	return value.DoExpressionNodeClass
}

func (*DoExpressionNode) DirectClass() *value.Class {
	return value.DoExpressionNodeClass
}

func (n *DoExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::DoExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  body: %%[\n")
	for i, stmt := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  catches: %%[\n")
	for i, stmt := range n.Catches {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  finally: %%[\n")
	for i, stmt := range n.Finally {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *DoExpressionNode) Error() string {
	return n.Inspect()
}

// Create a new `do` expression node eg.
//
//	do
//		print("awesome!")
//	end
func NewDoExpressionNode(span *position.Span, body []StatementNode, catches []*CatchNode, finally []StatementNode) *DoExpressionNode {
	return &DoExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
		Catches:       catches,
		Finally:       finally,
	}
}

// Represents a `catch` eg.
//
//	catch SomeError(message)
//		print("awesome!")
//	end
type CatchNode struct {
	NodeBase
	Pattern       PatternNode
	StackTraceVar IdentifierNode
	Body          []StatementNode // do expression body
}

func (*CatchNode) IsStatic() bool {
	return false
}

func (*CatchNode) Class() *value.Class {
	return value.CatchNodeClass
}

func (*CatchNode) DirectClass() *value.Class {
	return value.CatchNodeClass
}

func (n *CatchNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::CatchNode{\n  &: %p", n)

	buff.WriteString(",\n  pattern: ")
	indentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  stack_trace_var: ")
	indentStringFromSecondLine(&buff, n.StackTraceVar.Inspect(), 1)

	buff.WriteString(",\n  body: %%[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *CatchNode) Error() string {
	return n.Inspect()
}

// Create a new `catch` node eg.
//
//	catch SomeError(message)
//		print("awesome!")
//	end
func NewCatchNode(span *position.Span, pattern PatternNode, stackTraceVar IdentifierNode, body []StatementNode) *CatchNode {
	return &CatchNode{
		NodeBase:      NodeBase{span: span},
		Pattern:       pattern,
		StackTraceVar: stackTraceVar,
		Body:          body,
	}
}
