package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
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
	Cases    []*SwitchCaseNode
	ElseBody []StatementNode
}

func (n *SwitchExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SwitchExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value.splice(loc, args, unquote).(ExpressionNode),
		Cases:         SpliceSlice(n.Cases, loc, args, unquote),
		ElseBody:      SpliceSlice(n.ElseBody, loc, args, unquote),
	}
}

func (n *SwitchExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SwitchExpressionNode", env)
}

func (n *SwitchExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Value.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
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

func (n *SwitchExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SwitchExpressionNode)
	if !ok {
		return false
	}

	if len(n.Cases) != len(o.Cases) ||
		len(n.ElseBody) != len(o.ElseBody) ||
		!n.loc.Equal(o.loc) ||
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
func NewSwitchExpressionNode(loc *position.Location, val ExpressionNode, cases []*SwitchCaseNode, els []StatementNode) *SwitchExpressionNode {
	return &SwitchExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
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

	fmt.Fprintf(&buff, "Std::Elk::AST::SwitchExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

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

func (n *SwitchExpressionNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *SwitchExpressionNode) Error() string {
	return n.Inspect()
}

// Represents a `case` node eg. `case 3 then println("eureka!")`
type SwitchCaseNode struct {
	NodeBase
	Pattern PatternNode
	Body    []StatementNode
}

func (n *SwitchCaseNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SwitchCaseNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Pattern:  n.Pattern.splice(loc, args, unquote).(PatternNode),
		Body:     SpliceSlice(n.Body, loc, args, unquote),
	}
}

func (n *SwitchCaseNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SwitchCaseNode", env)
}

func (n *SwitchCaseNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Pattern.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, stmt := range n.Body {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *SwitchCaseNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SwitchCaseNode)
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

	return n.loc.Equal(o.loc)
}

func (n *SwitchCaseNode) String() string {
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

func (*SwitchCaseNode) IsStatic() bool {
	return false
}

func (*SwitchCaseNode) Class() *value.Class {
	return value.SwitchCaseNodeClass
}

func (*SwitchCaseNode) DirectClass() *value.Class {
	return value.SwitchCaseNodeClass
}

func (n *SwitchCaseNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SwitchCaseNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

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

func (n *SwitchCaseNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *SwitchCaseNode) Error() string {
	return n.Inspect()
}

// Create a new switch `case` node
func NewSwitchCaseNode(loc *position.Location, pattern PatternNode, body []StatementNode) *SwitchCaseNode {
	return &SwitchCaseNode{
		NodeBase: NodeBase{loc: loc},
		Pattern:  pattern,
		Body:     body,
	}
}
