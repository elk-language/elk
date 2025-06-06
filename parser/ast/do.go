package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
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

func (n *DoExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &DoExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Body:          SpliceSlice(n.Body, loc, args, unquote),
		Catches:       SpliceSlice(n.Catches, loc, args, unquote),
		Finally:       SpliceSlice(n.Finally, loc, args, unquote),
	}
}

func (n *DoExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::DoExpressionNode", env)
}

func (n *DoExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, stmt := range n.Body {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, catch := range n.Catches {
		if catch.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, stmt := range n.Finally {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
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

	return n.loc.Equal(o.loc)
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

	fmt.Fprintf(&buff, "Std::Elk::AST::DoExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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
func NewDoExpressionNode(loc *position.Location, body []StatementNode, catches []*CatchNode, finally []StatementNode) *DoExpressionNode {
	return &DoExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
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

func (n *CatchNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	pattern := n.Pattern.splice(loc, args, unquote).(PatternNode)

	var stackTraceVar IdentifierNode
	if n.StackTraceVar != nil {
		stackTraceVar = n.StackTraceVar.splice(loc, args, unquote).(IdentifierNode)
	}

	body := SpliceSlice(n.Body, loc, args, unquote)

	return &CatchNode{
		NodeBase:      NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Pattern:       pattern,
		StackTraceVar: stackTraceVar,
		Body:          body,
	}
}

func (n *CatchNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::CatchNode", env)
}

func (n *CatchNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, stmt := range n.Body {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	if n.Pattern.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.StackTraceVar != nil {
		if n.StackTraceVar.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
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

	return n.loc.Equal(o.loc)
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

	fmt.Fprintf(&buff, "Std::Elk::AST::CatchNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  stack_trace_var: ")
	if n.StackTraceVar == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.StackTraceVar.Inspect(), 1)
	}

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
func NewCatchNode(loc *position.Location, pattern PatternNode, stackTraceVar IdentifierNode, body []StatementNode) *CatchNode {
	return &CatchNode{
		NodeBase:      NodeBase{loc: loc},
		Pattern:       pattern,
		StackTraceVar: stackTraceVar,
		Body:          body,
	}
}
