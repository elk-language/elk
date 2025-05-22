package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an as pattern eg. `> 5 && < 20 as foo`
type AsPatternNode struct {
	NodeBase
	Pattern PatternNode
	Name    IdentifierNode
}

func (n *AsPatternNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &AsPatternNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Pattern:  n.Pattern.splice(loc, args, unquote).(PatternNode),
		Name:     n.Name.splice(loc, args, unquote).(IdentifierNode),
	}
}

func (n *AsPatternNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::AsPatternNode", env)
}

func (n *AsPatternNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Name.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}
	if n.Pattern.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (*AsPatternNode) IsStatic() bool {
	return false
}

// Create an Object pattern node eg. `Foo(foo: 5, bar: a, c)`
func NewAsPatternNode(loc *position.Location, pattern PatternNode, name IdentifierNode) *AsPatternNode {
	return &AsPatternNode{
		NodeBase: NodeBase{loc: loc},
		Pattern:  pattern,
		Name:     name,
	}
}

func (*AsPatternNode) Class() *value.Class {
	return value.AsPatternNodeClass
}

func (*AsPatternNode) DirectClass() *value.Class {
	return value.AsPatternNodeClass
}

func (n *AsPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AsPatternNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  name: ")
	indent.IndentStringFromSecondLine(&buff, n.Name.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AsPatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AsPatternNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Pattern.Equal(value.Ref(o.Pattern)) &&
		n.Name.Equal(value.Ref(o.Name))
}

func (n *AsPatternNode) Error() string {
	return n.Inspect()
}

func (n *AsPatternNode) String() string {
	var buff strings.Builder

	leftParen := PatternPrecedence(n) > PatternPrecedence(n.Pattern)
	if leftParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Pattern.String())
	if leftParen {
		buff.WriteRune(')')
	}

	buff.WriteString(" as ")

	buff.WriteString(n.Name.String())

	return buff.String()
}
