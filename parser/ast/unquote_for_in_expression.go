package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a `%for in` expression eg. `%for i in 5..15 then println(i)`.
// It's a for loop expression used in AST templates (quote blocks) to build ASTs.
type UnquoteForInExpressionNode struct {
	TypedNodeBase
	Pattern      PatternNode
	InExpression ExpressionNode  // expression that will be iterated through
	ThenBody     []StatementNode // then expression body
}

func (n *UnquoteForInExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UnquoteForInExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Pattern:       n.Pattern.splice(loc, args, unquote).(PatternNode),
		InExpression:  n.InExpression.splice(loc, args, unquote).(ExpressionNode),
		ThenBody:      SpliceSlice(n.ThenBody, loc, args, unquote),
	}
}

func (n *UnquoteForInExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UnquoteForInExpressionNode", env)
}

func (n *UnquoteForInExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Pattern.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.InExpression.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, stmt := range n.ThenBody {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *UnquoteForInExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UnquoteForInExpressionNode)
	if !ok {
		return false
	}

	if !n.loc.Equal(o.loc) ||
		!n.Pattern.Equal(value.Ref(o.Pattern)) ||
		!n.InExpression.Equal(value.Ref(o.InExpression)) ||
		len(n.ThenBody) != len(o.ThenBody) {
		return false
	}

	for i, stmt := range n.ThenBody {
		if !stmt.Equal(value.Ref(o.ThenBody[i])) {
			return false
		}
	}

	return true
}

func (n *UnquoteForInExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("%for ")
	buff.WriteString(n.Pattern.String())
	buff.WriteString(" in ")
	buff.WriteString(n.InExpression.String())
	buff.WriteRune('\n')

	for _, stmt := range n.ThenBody {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteString("\n")
	}
	buff.WriteString("end")

	return buff.String()
}

func (*UnquoteForInExpressionNode) IsStatic() bool {
	return false
}

// Create a new `%for in` expression node eg. `%for i in 5..15 then println(i)`
func NewUnquoteForInExpressionNode(loc *position.Location, pattern PatternNode, inExpr ExpressionNode, then []StatementNode) *UnquoteForInExpressionNode {
	return &UnquoteForInExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Pattern:       pattern,
		InExpression:  inExpr,
		ThenBody:      then,
	}
}

func NewUnquoteForInExpressionNodeI(loc *position.Location, pattern PatternNode, inExpr ExpressionNode, then []StatementNode) ExpressionNode {
	return NewUnquoteForInExpressionNode(loc, pattern, inExpr, then)
}

func (*UnquoteForInExpressionNode) Class() *value.Class {
	return value.UnquoteForInExpressionNodeClass
}

func (*UnquoteForInExpressionNode) DirectClass() *value.Class {
	return value.UnquoteForInExpressionNodeClass
}

func (n *UnquoteForInExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UnquoteForInExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  in_expression: ")
	indent.IndentStringFromSecondLine(&buff, n.InExpression.Inspect(), 1)

	buff.WriteString(",\n  then_body: %[")
	if len(n.ThenBody) > 0 {
		buff.WriteRune('\n')
		for i, element := range n.ThenBody {
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

func (n *UnquoteForInExpressionNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *UnquoteForInExpressionNode) Error() string {
	return n.Inspect()
}
