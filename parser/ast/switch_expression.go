package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `switch` expression eg.
//
//	switch a
//	case 3
//	  println("eureka!")
//	case nil
//	  println("boo")
//	else
//	  println("nothing")
//	end
type SwitchExpressionNode struct {
	TypedNodeBase
	Value    ExpressionNode
	Cases    []*CaseNode
	ElseBody []StatementNode
}

func (*SwitchExpressionNode) IsStatic() bool {
	return false
}

// Create a new `switch` expression node
func NewSwitchExpressionNode(span *position.Span, val ExpressionNode, cases []*CaseNode, els []StatementNode) *SwitchExpressionNode {
	return &SwitchExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
		Cases:         cases,
		ElseBody:      els,
	}
}

func (*SwitchExpressionNode) Class() *value.Class {
	return value.SwitchExpressionNodeClass
}

func (*SwitchExpressionNode) DirectClass() *value.Class {
	return value.SwitchExpressionNodeClass
}

func (n *SwitchExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SwitchExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  value: ")
	indentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString(",\n  body: %%[\n")
	for i, stmt := range n.Cases {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  else_body: %%[\n")
	for i, stmt := range n.ElseBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SwitchExpressionNode) Error() string {
	return n.Inspect()
}

// Represents a `case` node eg. `case 3 then println("eureka!")`
type CaseNode struct {
	NodeBase
	Pattern PatternNode
	Body    []StatementNode
}

func (*CaseNode) IsStatic() bool {
	return false
}

func (*CaseNode) Class() *value.Class {
	return value.CaseNodeClass
}

func (*CaseNode) DirectClass() *value.Class {
	return value.CaseNodeClass
}

func (n *CaseNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::CaseNode{\n  &: %p", n)

	buff.WriteString(",\n  pattern: ")
	indentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

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

func (n *CaseNode) Error() string {
	return n.Inspect()
}

// Create a new `case` node
func NewCaseNode(span *position.Span, pattern PatternNode, body []StatementNode) *CaseNode {
	return &CaseNode{
		NodeBase: NodeBase{span: span},
		Pattern:  pattern,
		Body:     body,
	}
}
