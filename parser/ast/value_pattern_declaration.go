package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a value pattern declaration eg. `val [foo, { bar }] = baz()`
type ValuePatternDeclarationNode struct {
	NodeBase
	Pattern     PatternNode
	Initialiser ExpressionNode // value assigned to the value
}

func (n *ValuePatternDeclarationNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ValuePatternDeclarationNode{
		NodeBase:    NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Pattern:     n.Pattern.splice(loc, args, unquote).(PatternNode),
		Initialiser: n.Initialiser.splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *ValuePatternDeclarationNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::ValuePatternDeclarationNode", env)
}

func (n *ValuePatternDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Pattern != nil {
		if n.Pattern.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	if n.Initialiser != nil {
		if n.Initialiser.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *ValuePatternDeclarationNode) traverseMap(parent Node, fn func(node, parent Node) Node) Node {
	var pattern PatternNode
	if n.Pattern != nil {
		pattern = n.Pattern.traverseMap(n, fn).(PatternNode)
	}

	var init ExpressionNode
	if n.Initialiser != nil {
		init = n.Initialiser.traverseMap(n, fn).(ExpressionNode)
	}

	newNode := NewValuePatternDeclarationNode(n.loc, pattern, init)

	return fn(newNode, parent)
}

func (n *ValuePatternDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ValuePatternDeclarationNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Pattern.Equal(value.Ref(o.Pattern)) &&
		n.Initialiser.Equal(value.Ref(o.Initialiser))
}

func (n *ValuePatternDeclarationNode) LaxEqual(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ValuePatternDeclarationNode)
	if !ok {
		return false
	}

	return n.Pattern.LaxEqual(value.Ref(o.Pattern)) &&
		n.Initialiser.LaxEqual(value.Ref(o.Initialiser))
}

func (n *ValuePatternDeclarationNode) String() string {
	var buff strings.Builder

	buff.WriteString("val ")
	buff.WriteString(n.Pattern.String())
	buff.WriteString(" = ")
	buff.WriteString(n.Initialiser.String())

	return buff.String()
}

func (*ValuePatternDeclarationNode) IsStatic() bool {
	return false
}

func (*ValuePatternDeclarationNode) Class() *value.Class {
	return value.ValuePatternDeclarationNodeClass
}

func (*ValuePatternDeclarationNode) DirectClass() *value.Class {
	return value.ValuePatternDeclarationNodeClass
}

func (n *ValuePatternDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ValuePatternDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	indent.IndentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ValuePatternDeclarationNode) ToValue() value.Value {
	return value.Ref(n)
}

func (v *ValuePatternDeclarationNode) Error() string {
	return v.Inspect()
}

// Create a new value declaration node eg. `val foo: String`
func NewValuePatternDeclarationNode(loc *position.Location, pattern PatternNode, init ExpressionNode) *ValuePatternDeclarationNode {
	return &ValuePatternDeclarationNode{
		NodeBase:    NodeBase{loc: loc},
		Pattern:     pattern,
		Initialiser: init,
	}
}
