package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Assignment with the specified operator.
type AssignmentExpressionNode struct {
	TypedNodeBase
	Op    *token.Token   // operator
	Left  ExpressionNode // left hand side
	Right ExpressionNode // right hand side
}

func (n *AssignmentExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &AssignmentExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Op:            n.Op,
		Left:          n.Left.splice(loc, args, unquote).(ExpressionNode),
		Right:         n.Right.splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *AssignmentExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::AssignmentExpressionNode", env)
}

func (n *AssignmentExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Right.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}
	if n.Left.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (*AssignmentExpressionNode) IsStatic() bool {
	return false
}

// Create a new assignment expression node eg. `foo = 3`
func NewAssignmentExpressionNode(loc *position.Location, op *token.Token, left, right ExpressionNode) *AssignmentExpressionNode {
	return &AssignmentExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

func (*AssignmentExpressionNode) Class() *value.Class {
	return value.AssignmentExpressionNodeClass
}

func (*AssignmentExpressionNode) DirectClass() *value.Class {
	return value.AssignmentExpressionNodeClass
}

func (n *AssignmentExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AssignmentExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  left: ")
	indent.IndentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indent.IndentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AssignmentExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AssignmentExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Op.Equal(o.Op) &&
		n.Left.Equal(value.Ref(o.Left)) &&
		n.Right.Equal(value.Ref(o.Right))
}

func (n *AssignmentExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Left.String())
	buff.WriteRune(' ')
	buff.WriteString(n.Op.FetchValue())
	buff.WriteRune(' ')

	rightStr := n.Right.String()
	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Right)
	if strings.ContainsRune(rightStr, '\n') {
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteRune('\n')
		indent.IndentString(&buff, rightStr, 1)
		if parens {
			buff.WriteString("\n)")
		}
	} else {
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteString(rightStr)
		if parens {
			buff.WriteRune(')')
		}
	}

	return buff.String()
}

func (p *AssignmentExpressionNode) Error() string {
	return p.Inspect()
}
