package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a variable declaration with patterns eg. `var [foo, { bar }] = baz()`
type VariablePatternDeclarationNode struct {
	NodeBase
	Pattern     PatternNode
	Initialiser ExpressionNode // value assigned to the variable
}

func (n *VariablePatternDeclarationNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &VariablePatternDeclarationNode{
		NodeBase:    NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Pattern:     n.Pattern.Splice(loc, args, unquote).(PatternNode),
		Initialiser: n.Initialiser.Splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *VariablePatternDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

func (n *VariablePatternDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*VariablePatternDeclarationNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Pattern.Equal(value.Ref(o.Pattern)) &&
		n.Initialiser.Equal(value.Ref(o.Initialiser))
}

func (n *VariablePatternDeclarationNode) String() string {
	var buff strings.Builder

	buff.WriteString("var ")
	buff.WriteString(n.Pattern.String())
	buff.WriteString(" = ")
	buff.WriteString(n.Initialiser.String())

	return buff.String()
}

func (*VariablePatternDeclarationNode) IsStatic() bool {
	return false
}

func (*VariablePatternDeclarationNode) Class() *value.Class {
	return value.VariablePatternDeclarationNodeClass
}

func (*VariablePatternDeclarationNode) DirectClass() *value.Class {
	return value.VariablePatternDeclarationNodeClass
}

func (n *VariablePatternDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::VariablePatternDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	indent.IndentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (v *VariablePatternDeclarationNode) Error() string {
	return v.Inspect()
}

// Create a new variable declaration node with patterns eg. `var [foo, { bar }] = baz()`
func NewVariablePatternDeclarationNode(loc *position.Location, pattern PatternNode, init ExpressionNode) *VariablePatternDeclarationNode {
	return &VariablePatternDeclarationNode{
		NodeBase:    NodeBase{loc: loc},
		Pattern:     pattern,
		Initialiser: init,
	}
}
