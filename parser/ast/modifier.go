package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Represents an `if`, `unless`, `while` or `until` modifier expression eg. `return true if foo`.
type ModifierNode struct {
	TypedNodeBase
	Modifier *token.Token   // modifier token
	Left     ExpressionNode // left hand side
	Right    ExpressionNode // right hand side
}

func (n *ModifierNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ModifierNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Modifier:      n.Modifier.Splice(loc, unquote),
		Left:          n.Left.Splice(loc, args, unquote).(ExpressionNode),
		Right:         n.Right.Splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *ModifierNode) Traverse(yield func(Node) bool) bool {
	if n.Left.Traverse(yield) {
		return false
	}
	if n.Right.Traverse(yield) {
		return false
	}
	return yield(n)
}

func (n *ModifierNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ModifierNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Modifier.Equal(o.Modifier) &&
		n.Left.Equal(value.Ref(o.Left)) &&
		n.Right.Equal(value.Ref(o.Right))
}

func (n *ModifierNode) String() string {
	var buff strings.Builder

	leftParen := ExpressionPrecedence(n) > ExpressionPrecedence(n.Left)
	rightParen := ExpressionPrecedence(n) >= ExpressionPrecedence(n.Right)

	if leftParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Left.String())
	if leftParen {
		buff.WriteRune(')')
	}

	buff.WriteRune(' ')
	buff.WriteString(n.Modifier.FetchValue())
	buff.WriteRune(' ')

	if rightParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Right.String())
	if rightParen {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (*ModifierNode) IsStatic() bool {
	return false
}

func (*ModifierNode) Class() *value.Class {
	return value.ModifierNodeClass
}

func (*ModifierNode) DirectClass() *value.Class {
	return value.ModifierNodeClass
}

func (n *ModifierNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ModifierNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  modifier: ")
	indent.IndentStringFromSecondLine(&buff, n.Modifier.Inspect(), 1)

	buff.WriteString(",\n  left: ")
	indent.IndentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indent.IndentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (m *ModifierNode) Error() string {
	return m.Inspect()
}

// Create a new modifier node eg. `return true if foo`.
func NewModifierNode(loc *position.Location, mod *token.Token, left, right ExpressionNode) *ModifierNode {
	return &ModifierNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Modifier:      mod,
		Left:          left,
		Right:         right,
	}
}

// Represents an `if .. else` modifier expression eg. `foo = 1 if bar else foo = 2`
type ModifierIfElseNode struct {
	TypedNodeBase
	ThenExpression ExpressionNode // then expression body
	Condition      ExpressionNode // if condition
	ElseExpression ExpressionNode // else expression body
}

func (n *ModifierIfElseNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ModifierIfElseNode{
		TypedNodeBase:  TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		ThenExpression: n.ThenExpression.Splice(loc, args, unquote).(ExpressionNode),
		Condition:      n.Condition.Splice(loc, args, unquote).(ExpressionNode),
		ElseExpression: n.ElseExpression.Splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *ModifierIfElseNode) Traverse(yield func(Node) bool) bool {
	if n.ThenExpression.Traverse(yield) {
		return false
	}
	if n.Condition.Traverse(yield) {
		return false
	}
	if n.ElseExpression.Traverse(yield) {
		return false
	}
	return yield(n)
}

func (n *ModifierIfElseNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ModifierIfElseNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.ThenExpression.Equal(value.Ref(o.ThenExpression)) &&
		n.Condition.Equal(value.Ref(o.Condition)) &&
		n.ElseExpression.Equal(value.Ref(o.ElseExpression))
}

func (n *ModifierIfElseNode) String() string {
	var buff strings.Builder

	thenParens := ExpressionPrecedence(n) > ExpressionPrecedence(n.ThenExpression)
	if thenParens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.ThenExpression.String())
	if thenParens {
		buff.WriteRune(')')
	}

	buff.WriteString(" if ")
	buff.WriteString(n.Condition.String())
	buff.WriteString(" else ")

	elseParens := ExpressionPrecedence(n) > ExpressionPrecedence(n.ElseExpression)
	if elseParens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.ElseExpression.String())
	if elseParens {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (*ModifierIfElseNode) IsStatic() bool {
	return false
}

func (*ModifierIfElseNode) Class() *value.Class {
	return value.ModifierIfElseNodeClass
}

func (*ModifierIfElseNode) DirectClass() *value.Class {
	return value.ModifierIfElseNodeClass
}

func (n *ModifierIfElseNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ModifierIfElseNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  then_expression: ")
	indent.IndentStringFromSecondLine(&buff, n.ThenExpression.Inspect(), 1)

	buff.WriteString(",\n  condition: ")
	indent.IndentStringFromSecondLine(&buff, n.Condition.Inspect(), 1)

	buff.WriteString(",\n  else_expression: ")
	indent.IndentStringFromSecondLine(&buff, n.ElseExpression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (m *ModifierIfElseNode) Error() string {
	return m.Inspect()
}

// Create a new modifier `if` .. `else` node eg. `foo = 1 if bar else foo = 2“.
func NewModifierIfElseNode(loc *position.Location, then, cond, els ExpressionNode) *ModifierIfElseNode {
	return &ModifierIfElseNode{
		TypedNodeBase:  TypedNodeBase{loc: loc},
		ThenExpression: then,
		Condition:      cond,
		ElseExpression: els,
	}
}

// Represents an `for .. in` modifier expression eg. `println(i) for i in 10..30`
type ModifierForInNode struct {
	TypedNodeBase
	ThenExpression ExpressionNode // then expression body
	Pattern        PatternNode
	InExpression   ExpressionNode // expression that will be iterated through
}

func (n *ModifierForInNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ModifierForInNode{
		TypedNodeBase:  TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		ThenExpression: n.ThenExpression.Splice(loc, args, unquote).(ExpressionNode),
		Pattern:        n.Pattern.Splice(loc, args, unquote).(PatternNode),
		InExpression:   n.InExpression.Splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *ModifierForInNode) Traverse(yield func(Node) bool) bool {
	if n.ThenExpression.Traverse(yield) {
		return false
	}
	if n.Pattern.Traverse(yield) {
		return false
	}
	if n.InExpression.Traverse(yield) {
		return false
	}
	return yield(n)
}

func (n *ModifierForInNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ModifierForInNode)
	if !ok {
		return false
	}

	return n.ThenExpression.Equal(value.Ref(o.ThenExpression)) &&
		n.Pattern.Equal(value.Ref(o.Pattern)) &&
		n.InExpression.Equal(value.Ref(o.InExpression)) &&
		n.loc.Equal(o.loc)
}

func (n *ModifierForInNode) String() string {
	var buff strings.Builder

	thenParens := ExpressionPrecedence(n) > ExpressionPrecedence(n.ThenExpression)
	if thenParens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.ThenExpression.String())
	if thenParens {
		buff.WriteRune(')')
	}

	buff.WriteString(" for ")
	buff.WriteString(n.Pattern.String())
	buff.WriteString(" in ")
	buff.WriteString(n.InExpression.String())
	return buff.String()
}

func (*ModifierForInNode) IsStatic() bool {
	return false
}

func (*ModifierForInNode) Class() *value.Class {
	return value.ModifierForInNodeClass
}

func (*ModifierForInNode) DirectClass() *value.Class {
	return value.ModifierForInNodeClass
}

func (n *ModifierForInNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ModifierForInNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  then_expression: ")
	indent.IndentStringFromSecondLine(&buff, n.ThenExpression.Inspect(), 1)

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  in_expression: ")
	indent.IndentStringFromSecondLine(&buff, n.InExpression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (m *ModifierForInNode) Error() string {
	return m.Inspect()
}

// Create a new modifier `for` .. `in` node eg. `println(i) for i in 10..30`
func NewModifierForInNode(loc *position.Location, then ExpressionNode, pattern PatternNode, in ExpressionNode) *ModifierForInNode {
	return &ModifierForInNode{
		TypedNodeBase:  TypedNodeBase{loc: loc},
		ThenExpression: then,
		Pattern:        pattern,
		InExpression:   in,
	}
}
