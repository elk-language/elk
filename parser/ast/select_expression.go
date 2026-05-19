package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a `select` expression eg.
//
//	select
//	case v := <<ch
//	  println(v)
//	case ch2 << "foo"
//	  println("boo")
//	else
//	  println("nothing")
//	end
type SelectExpressionNode struct {
	TypedNodeBase
	Cases    []*SelectCaseNode
	ElseBody []StatementNode
}

func (n *SelectExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SelectExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Cases:         SpliceSlice(n.Cases, loc, args, unquote),
		ElseBody:      SpliceSlice(n.ElseBody, loc, args, unquote),
	}
}

func (n *SelectExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SelectExpressionNode", env)
}

func (n *SelectExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, caseNode := range n.Cases {
		if caseNode.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, stmt := range n.ElseBody {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *SelectExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SelectExpressionNode)
	if !ok {
		return false
	}

	if len(n.Cases) != len(o.Cases) ||
		len(n.ElseBody) != len(o.ElseBody) ||
		!n.loc.Equal(o.loc) {
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

func (n *SelectExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("select\n")

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

func (*SelectExpressionNode) IsStatic() bool {
	return false
}

// Create a new `select` expression node
func NewSelectExpressionNode(loc *position.Location, cases []*SelectCaseNode, els []StatementNode) *SelectExpressionNode {
	return &SelectExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Cases:         cases,
		ElseBody:      els,
	}
}

func (*SelectExpressionNode) Class() *value.Class {
	return value.SelectExpressionNodeClass
}

func (*SelectExpressionNode) DirectClass() *value.Class {
	return value.SelectExpressionNodeClass
}

func (n *SelectExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SelectExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  cases: %[")
	if len(n.Cases) > 0 {
		buff.WriteRune('\n')
		for i, element := range n.Cases {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, element.Inspect(), 2)
		}
		buff.WriteString("\n  ")
	}
	buff.WriteRune(']')

	buff.WriteString(",\n  else_body: %[")
	if len(n.ElseBody) > 0 {
		buff.WriteRune('\n')
		for i, element := range n.ElseBody {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, element.Inspect(), 2)
		}
		buff.WriteString("\n  ")
	}
	buff.WriteRune(']')

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SelectExpressionNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *SelectExpressionNode) Error() string {
	return n.Inspect()
}

// Represents a select `case` node eg. `case v := <<ch then println("eureka!")`
type SelectCaseNode struct {
	NodeBase
	Expression ExpressionNode
	Body       []StatementNode
}

func (n *SelectCaseNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	expr := n.Expression.splice(loc, args, unquote).(ExpressionNode)
	body := SpliceSlice(n.Body, loc, args, unquote)

	return &SelectCaseNode{
		NodeBase:   NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Expression: expr,
		Body:       body,
	}
}

func (n *SelectCaseNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SelectCaseNode", env)
}

func (n *SelectCaseNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Expression.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, stmt := range n.Body {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *SelectCaseNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SelectCaseNode)
	if !ok {
		return false
	}

	if !n.Expression.Equal(value.Ref(o.Expression)) {
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

func (n *SelectCaseNode) String() string {
	var buff strings.Builder

	buff.WriteString("case ")
	buff.WriteString(n.Expression.String())
	buff.WriteRune('\n')

	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}

	return buff.String()
}

func (*SelectCaseNode) IsStatic() bool {
	return false
}

func (*SelectCaseNode) Class() *value.Class {
	return value.SelectCaseNodeClass
}

func (*SelectCaseNode) DirectClass() *value.Class {
	return value.SelectCaseNodeClass
}

func (n *SelectCaseNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SelectCaseNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

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

func (n *SelectCaseNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *SelectCaseNode) Error() string {
	return n.Inspect()
}

// Create a new select `case` node
func NewSelectCaseNode(loc *position.Location, expr ExpressionNode, body []StatementNode) *SelectCaseNode {
	return &SelectCaseNode{
		NodeBase:   NodeBase{loc: loc},
		Expression: expr,
		Body:       body,
	}
}
