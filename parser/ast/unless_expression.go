package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an `unless` expression eg. `unless foo then println("bar")`
type UnlessExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // unless condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

func (n *UnlessExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UnlessExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Condition:     n.Condition.splice(loc, args, unquote).(ExpressionNode),
		ThenBody:      SpliceSlice(n.ThenBody, loc, args, unquote),
		ElseBody:      SpliceSlice(n.ElseBody, loc, args, unquote),
	}
}

func (n *UnlessExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UnlessExpressionNode", env)
}

func (n *UnlessExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Condition.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, stmt := range n.ThenBody {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
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

func (n *UnlessExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UnlessExpressionNode)
	if !ok {
		return false
	}

	if len(n.ThenBody) != len(o.ThenBody) ||
		len(n.ElseBody) != len(o.ElseBody) ||
		!n.Condition.Equal(value.Ref(o.Condition)) ||
		!n.loc.Equal(o.loc) {
		return false
	}

	for i, stmt := range n.ThenBody {
		if !stmt.Equal(value.Ref(o.ThenBody[i])) {
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

func (n *UnlessExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("unless ")
	buff.WriteString(n.Condition.String())

	buff.WriteRune('\n')
	for _, stmt := range n.ThenBody {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}

	if len(n.ElseBody) > 0 {
		buff.WriteString("else\n")
		for _, stmt := range n.ElseBody {
			indent.IndentString(&buff, stmt.String(), 1)
			buff.WriteRune('\n')
		}
	}

	buff.WriteString("end")

	return buff.String()
}

func (*UnlessExpressionNode) IsStatic() bool {
	return false
}

// Create a new `unless` expression node eg. `unless foo then println("bar")`
func NewUnlessExpressionNode(loc *position.Location, cond ExpressionNode, then, els []StatementNode) *UnlessExpressionNode {
	return &UnlessExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		ThenBody:      then,
		Condition:     cond,
		ElseBody:      els,
	}
}

func (*UnlessExpressionNode) Class() *value.Class {
	return value.UnlessExpressionNodeClass
}

func (*UnlessExpressionNode) DirectClass() *value.Class {
	return value.UnlessExpressionNodeClass
}

func (n *UnlessExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UnlessExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  condition: ")
	indent.IndentStringFromSecondLine(&buff, n.Condition.Inspect(), 1)

	buff.WriteString(",\n  then_body: %[\n")
	for i, stmt := range n.ThenBody {
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

func (n *UnlessExpressionNode) Error() string {
	return n.Inspect()
}
