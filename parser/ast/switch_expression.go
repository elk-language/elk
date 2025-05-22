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
	Cases    []*CaseNode
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
func NewSwitchExpressionNode(loc *position.Location, val ExpressionNode, cases []*CaseNode, els []StatementNode) *SwitchExpressionNode {
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

	buff.WriteString(",\n  body: %[\n")
	for i, stmt := range n.Cases {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  else_body: %[\n")
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

func (n *CaseNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &CaseNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Pattern:  n.Pattern.splice(loc, args, unquote).(PatternNode),
		Body:     SpliceSlice(n.Body, loc, args, unquote),
	}
}

func (n *CaseNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::CaseNode", env)
}

func (n *CaseNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

	return n.loc.Equal(o.loc)
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

	fmt.Fprintf(&buff, "Std::Elk::AST::CaseNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *CaseNode) Error() string {
	return n.Inspect()
}

// Create a new `case` node
func NewCaseNode(loc *position.Location, pattern PatternNode, body []StatementNode) *CaseNode {
	return &CaseNode{
		NodeBase: NodeBase{loc: loc},
		Pattern:  pattern,
		Body:     body,
	}
}
