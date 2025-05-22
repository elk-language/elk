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

// Expression of a logical operator with two operands eg. `foo && bar`
type LogicalExpressionNode struct {
	TypedNodeBase
	Op     *token.Token   // operator
	Left   ExpressionNode // left hand side
	Right  ExpressionNode // right hand side
	static bool
}

func (n *LogicalExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &LogicalExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Op:            n.Op.Splice(loc, unquote),
		Left:          n.Left.splice(loc, args, unquote).(ComplexConstantNode),
		Right:         n.Right.splice(loc, args, unquote).(ComplexConstantNode),
	}
}

func (n *LogicalExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::LogicalExpressionNode", env)
}

func (n *LogicalExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Left.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.Right.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *LogicalExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*LogicalExpressionNode)
	if !ok {
		return false
	}

	return n.Op.Equal(o.Op) &&
		n.Left.Equal(value.Ref(o.Left)) &&
		n.Right.Equal(value.Ref(o.Right)) &&
		n.loc.Equal(o.loc)
}

func (n *LogicalExpressionNode) String() string {
	var buff strings.Builder

	associativity := ExpressionAssociativity(n)

	var leftParen bool
	var rightParen bool
	if associativity == LEFT_ASSOCIATIVE {
		leftParen = ExpressionPrecedence(n) > ExpressionPrecedence(n.Left)
		rightParen = ExpressionPrecedence(n) >= ExpressionPrecedence(n.Right)
	} else {
		leftParen = ExpressionPrecedence(n) >= ExpressionPrecedence(n.Left)
		rightParen = ExpressionPrecedence(n) > ExpressionPrecedence(n.Right)
	}

	if leftParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Left.String())
	if leftParen {
		buff.WriteRune(')')
	}

	buff.WriteString(" ")
	buff.WriteString(n.Op.String())
	buff.WriteString(" ")

	if rightParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Right.String())
	if rightParen {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (l *LogicalExpressionNode) IsStatic() bool {
	return l.static
}

// Create a new logical expression node.
func NewLogicalExpressionNode(loc *position.Location, op *token.Token, left, right ExpressionNode) *LogicalExpressionNode {
	return &LogicalExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Op:            op,
		Left:          left,
		Right:         right,
		static:        areExpressionsStatic(left, right),
	}
}

// Same as [NewLogicalExpressionNode] but returns an interface
func NewLogicalExpressionNodeI(loc *position.Location, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return &LogicalExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

func (*LogicalExpressionNode) Class() *value.Class {
	return value.LogicalExpressionNodeClass
}

func (*LogicalExpressionNode) DirectClass() *value.Class {
	return value.LogicalExpressionNodeClass
}

func (n *LogicalExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::LogicalExpressionNode{\n  loc: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  left: ")
	indent.IndentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indent.IndentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *LogicalExpressionNode) Error() string {
	return n.Inspect()
}
