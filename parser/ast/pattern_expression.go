package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a pattern expression `pattern String?`
type PatternExpressionNode struct {
	NodeBase
	PatternNode PatternNode
}

func (n *PatternExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &PatternExpressionNode{
		NodeBase:    NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		PatternNode: n.PatternNode.splice(loc, args, unquote).(PatternNode),
	}
}

func (n *PatternExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::PatternExpressionNode", env)
}

func (n *PatternExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.PatternNode.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *PatternExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PatternExpressionNode)
	if !ok {
		return false
	}

	return n.PatternNode.Equal(value.Ref(o.PatternNode)) &&
		n.loc.Equal(o.loc)
}

func (n *PatternExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("pattern ")
	buff.WriteString(n.PatternNode.String())

	return buff.String()
}

func (*PatternExpressionNode) IsStatic() bool {
	return false
}

func (*PatternExpressionNode) Class() *value.Class {
	return value.PatternExpressionNodeClass
}

func (*PatternExpressionNode) DirectClass() *value.Class {
	return value.PatternExpressionNodeClass
}

func (n *PatternExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::PatternExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  pattern_node: ")
	indent.IndentStringFromSecondLine(&buff, n.PatternNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *PatternExpressionNode) Error() string {
	return n.Inspect()
}

// Create a new pattern expression `pattern Foo(a: String, b: > 2)`
func NewPatternExpressionNode(loc *position.Location, patternNode PatternNode) *PatternExpressionNode {
	return &PatternExpressionNode{
		NodeBase:    NodeBase{loc: loc},
		PatternNode: patternNode,
	}
}
