package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

type UnquoteKind uint8

const (
	UNQUOTE_EXPRESSION_KIND         UnquoteKind = iota // Default kind of unquote for expression nodes
	UNQUOTE_PATTERN_KIND                               // Unquote kind for pattern nodes
	UNQUOTE_PATTERN_EXPRESSION_KIND                    // Unquote kind for pattern expression nodes
	UNQUOTE_TYPE_KIND                                  // Unquote kind for type nodes
	UNQUOTE_CONSTANT_KIND                              // Unquote kind for constant nodes
	UNQUOTE_IDENTIFIER_KIND                            // Unquote kind for identifier nodes
)

type UnquoteOrInvalidNode interface {
	ExpressionNode
	PatternNode
	TypeNode
	ConstantNode
	IdentifierNode
	unquote()
}

func (n *UnquoteNode) unquote() {}
func (n *InvalidNode) unquote() {}

// Represents an unquoted piece of AST inside of a quote
type UnquoteNode struct {
	TypedNodeBase
	Kind       UnquoteKind
	Expression ExpressionNode
}

func (n *UnquoteNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Expression.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *UnquoteNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UnquoteNode", env)
}

func (n *UnquoteNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	if args == nil || len(*args) == 0 {
		panic("too few arguments for splicing AST nodes")
	}

	arg := (*args)[0]
	*args = (*args)[1:]

	var targetLoc *position.Location
	if loc != nil && loc != position.ZeroLocation {
		targetLoc = loc.Copy()
		targetLoc.Parent = n.loc
	}

	return arg.splice(targetLoc, nil, true)
}

// Check if this node equals another node.
func (n *UnquoteNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UnquoteNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Expression.Equal(value.Ref(o.Expression))
}

// Return a string representation of the node.
func (n *UnquoteNode) String() string {
	var buff strings.Builder

	buff.WriteString("unquote(")

	exprStr := n.Expression.String()
	if strings.ContainsRune(exprStr, '\n') {
		buff.WriteRune('\n')
		indent.IndentString(&buff, exprStr, 1)
		buff.WriteRune('\n')
	} else {
		buff.WriteString(exprStr)
	}

	buff.WriteRune(')')

	return buff.String()
}

func (*UnquoteNode) IsStatic() bool {
	return false
}

func (*UnquoteNode) Class() *value.Class {
	return value.MacroBoundaryNodeClass
}

func (*UnquoteNode) DirectClass() *value.Class {
	return value.MacroBoundaryNodeClass
}

func (n *UnquoteNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UnquoteNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  node: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UnquoteNode) Error() string {
	return n.Inspect()
}

// Create an unquote node eg.
//
//	unquote(x)
func NewUnquoteNode(loc *position.Location, kind UnquoteKind, expr ExpressionNode) *UnquoteNode {
	return &UnquoteNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Kind:          kind,
		Expression:    expr,
	}
}
