package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a nilable pattern eg. `Foo(bar, baz)?`
type NilablePatternNode struct {
	TypedNodeBase
	Pattern PatternNode
}

func (n *NilablePatternNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &NilablePatternNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Pattern:       n.Pattern.splice(loc, args, unquote).(PatternNode),
	}
}

func (n *NilablePatternNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::NilablePatternNode", env)
}

func (n *NilablePatternNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Pattern.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (*NilablePatternNode) IsStatic() bool {
	return false
}

// Create a nilable pattern node eg. `Foo(foo: 5, bar: a, c)?`
func NewNilablePatternNode(loc *position.Location, pattern PatternNode) *NilablePatternNode {
	return &NilablePatternNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Pattern:       pattern,
	}
}

func (*NilablePatternNode) Class() *value.Class {
	return value.NilablePatternNodeClass
}

func (*NilablePatternNode) DirectClass() *value.Class {
	return value.NilablePatternNodeClass
}

func (n *NilablePatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NilablePatternNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NilablePatternNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *NilablePatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NilablePatternNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Pattern.Equal(value.Ref(o.Pattern))
}

func (n *NilablePatternNode) Error() string {
	return n.Inspect()
}

func (n *NilablePatternNode) String() string {
	var buff strings.Builder

	leftParen := PatternPrecedence(n) > PatternPrecedence(n.Pattern)
	if leftParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Pattern.String())
	if leftParen {
		buff.WriteRune(')')
	}

	buff.WriteRune('?')

	return buff.String()
}
