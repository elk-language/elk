package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
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

// Check if this node equals another node.
func (n *DoExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*DoExpressionNode)
	if !ok {
		return false
	}

	if len(n.Body) != len(o.Body) ||
		len(n.Catches) != len(o.Catches) ||
		len(n.Finally) != len(o.Finally) {
		return false
	}

	for i, stmt := range n.Body {
		if !stmt.Equal(value.Ref(o.Body[i])) {
			return false
		}
	}

	for i, catch := range n.Catches {
		if !catch.Equal(value.Ref(o.Catches[i])) {
			return false
		}
	}

	for i, stmt := range n.Finally {
		if !stmt.Equal(value.Ref(o.Finally[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *DoExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("do\n")

	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteString("\n")
	}

	for _, catch := range n.Catches {
		buff.WriteString(catch.String())
		buff.WriteString("\n")
	}

	if len(n.Finally) > 0 {
		buff.WriteString("finally\n")
		for _, stmt := range n.Finally {
			indent.IndentString(&buff, stmt.String(), 1)
			buff.WriteString("\n")
		}
	}

	buff.WriteString("end")

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::DoExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  body: %[\n")
	for i, stmt := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  catches: %[\n")
	for i, stmt := range n.Catches {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  finally: %[\n")
	for i, stmt := range n.Finally {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
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

func (n *CatchNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*CatchNode)
	if !ok {
		return false
	}

	if !n.Pattern.Equal(value.Ref(o.Pattern)) {
		return false
	}

	if n.StackTraceVar == o.StackTraceVar {
	} else if n.StackTraceVar == nil || o.StackTraceVar == nil {
		return false
	} else if !n.StackTraceVar.Equal(value.Ref(o.StackTraceVar)) {
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

	return n.Span().Equal(o.Span())
}

func (n *CatchNode) String() string {
	var buff strings.Builder

	buff.WriteString("catch ")
	buff.WriteString(n.Pattern.String())

	if n.StackTraceVar != nil {
		buff.WriteString(", ")
		buff.WriteString(n.StackTraceVar.String())
	}

	buff.WriteRune('\n')
	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::CatchNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  stack_trace_var: ")
	indent.IndentStringFromSecondLine(&buff, n.StackTraceVar.Inspect(), 1)

	buff.WriteString(",\n  body: %[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
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
