package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a syntax error.
type InvalidNode struct {
	NodeBase
	Token *token.Token
}

func (n *InvalidNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &InvalidNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Token:    n.Token.Splice(loc, unquote),
	}
}

func (n *InvalidNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *InvalidNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InvalidNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Token.Equal(o.Token)
}

func (n *InvalidNode) String() string {
	return "<invalid>"
}

func (*InvalidNode) SetType(types.Type) {}

func (*InvalidNode) IsStatic() bool {
	return false
}

func (*InvalidNode) IsOptional() bool {
	return false
}

func (*InvalidNode) Class() *value.Class {
	return value.InvalidNodeClass
}

func (*InvalidNode) DirectClass() *value.Class {
	return value.InvalidNodeClass
}

func (n *InvalidNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::InvalidNode{span: %s, token: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Token.Inspect(),
	)
}

func (p *InvalidNode) Error() string {
	return p.Inspect()
}

// Create a new invalid node.
func NewInvalidNode(loc *position.Location, tok *token.Token) *InvalidNode {
	return &InvalidNode{
		NodeBase: NodeBase{loc: loc},
		Token:    tok,
	}
}

func NewInvalidExpressionNode(loc *position.Location, tok *token.Token) ExpressionNode {
	return NewInvalidNode(loc, tok)
}

func NewInvalidPatternNode(loc *position.Location, tok *token.Token) PatternNode {
	return NewInvalidNode(loc, tok)
}

func NewInvalidPatternExpressionNode(loc *position.Location, tok *token.Token) PatternExpressionNode {
	return NewInvalidNode(loc, tok)
}
