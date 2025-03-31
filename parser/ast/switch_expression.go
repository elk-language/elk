package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
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

func (n *SwitchExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SwitchExpressionNode)
	if !ok {
		return false
	}

	if len(n.Cases) != len(o.Cases) ||
		len(n.ElseBody) != len(o.ElseBody) ||
		!n.span.Equal(o.span) ||
		!n.Value.Equal(value.Ref(o.Value)) {
		return false
	}

	for i, c := range n.Cases {
		if !c.Equal(value.Ref(o.Cases[i])) {
			return false
		}
	}

	for i, stmt := range n.ElseBody {
		if !stmt.Equal(value.Ref(o.ElseBody[i])) {
			return false
		}
	}

	return true
}

func (n *SwitchExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("switch ")
	buff.WriteString(n.Value.String())
	buff.WriteString("\n")

	for _, c := range n.Cases {
		buff.WriteString(c.String())
		buff.WriteString("\n")
	}

	if len(n.ElseBody) > 0 {
		buff.WriteString("else\n")
		for _, stmt := range n.ElseBody {
			indent.IndentString(&buff, stmt.String(), 1)
			buff.WriteString("\n")
		}
	}

	buff.WriteString("end")

	return buff.String()
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
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString(",\n  body: %%[\n")
	for i, stmt := range n.Cases {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  else_body: %%[\n")
	for i, stmt := range n.ElseBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
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

func (n *CaseNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*CaseNode)
	if !ok {
		return false
	}

	if !n.Pattern.Equal(value.Ref(o.Pattern)) {
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

func (n *CaseNode) String() string {
	var buff strings.Builder

	buff.WriteString("case ")
	buff.WriteString(n.Pattern.String())
	buff.WriteRune('\n')

	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}

	return buff.String()
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
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  body: %%[\n")
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
