package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an `if` expression eg. `if foo then println("bar")`
type IfExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // if condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

func (n *IfExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &IfExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Condition:     n.Condition.splice(loc, args, unquote).(ComplexConstantNode),
		ThenBody:      SpliceSlice(n.ThenBody, loc, args, unquote),
		ElseBody:      SpliceSlice(n.ElseBody, loc, args, unquote),
	}
}

func (n *IfExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::IfExpressionNode", env)
}

func (n *IfExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

// Check if this node equals another node.
func (n *IfExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*IfExpressionNode)
	if !ok {
		return false
	}

	if !n.Condition.Equal(value.Ref(o.Condition)) ||
		!n.loc.Equal(o.loc) {
		return false
	}

	if len(n.ThenBody) != len(o.ThenBody) ||
		len(n.ElseBody) != len(o.ElseBody) {
		return false
	}

	for i, element := range n.ThenBody {
		if !element.Equal(value.Ref(o.ThenBody[i])) {
			return false
		}
	}

	for i, element := range n.ElseBody {
		if !element.Equal(value.Ref(o.ElseBody[i])) {
			return false
		}
	}

	return true
}

// Return a string representation of the node.
func (n *IfExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("if ")
	buff.WriteString(n.Condition.String())
	buff.WriteRune('\n')

	for _, stmt := range n.ThenBody {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}

	if len(n.ElseBody) > 0 {
		then := n.ElseBody[0]
		parens := ExpressionPrecedence(n) > StatementPrecedence(then)
		if len(n.ElseBody) == 1 && !parens {
			buff.WriteString("else ")
			buff.WriteString(then.String())
		} else {
			buff.WriteString("else\n")
			for _, stmt := range n.ElseBody {
				indent.IndentString(&buff, stmt.String(), 1)
				buff.WriteRune('\n')
			}
			buff.WriteString("end")
		}
	} else {
		buff.WriteString("end")
	}

	return buff.String()
}

func (*IfExpressionNode) IsStatic() bool {
	return false
}

// Create a new `if` expression node eg. `if foo then println("bar")`
func NewIfExpressionNode(loc *position.Location, cond ExpressionNode, then, els []StatementNode) *IfExpressionNode {
	return &IfExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		ThenBody:      then,
		Condition:     cond,
		ElseBody:      els,
	}
}

func (*IfExpressionNode) Class() *value.Class {
	return value.IfExpressionNodeClass
}

func (*IfExpressionNode) DirectClass() *value.Class {
	return value.IfExpressionNodeClass
}

func (n *IfExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::IfExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  condition: ")
	indent.IndentStringFromSecondLine(&buff, n.Condition.Inspect(), 1)

	buff.WriteString(",\n  then_body: %[\n")
	for i, element := range n.ThenBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  else_body: %[\n")
	for i, element := range n.ElseBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *IfExpressionNode) Error() string {
	return n.Inspect()
}
