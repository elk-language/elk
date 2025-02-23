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
	return fmt.Sprintf("Std::AST::InvalidNode{&: %p, token: %s}", n, n.Token.Inspect())
}

func (p *InvalidNode) Error() string {
	return p.Inspect()
}

// Create a new invalid node.
func NewInvalidNode(span *position.Span, tok *token.Token) *InvalidNode {
	return &InvalidNode{
		NodeBase: NodeBase{span: span},
		Token:    tok,
	}
}

func NewInvalidExpressionNode(span *position.Span, tok *token.Token) ExpressionNode {
	return NewInvalidNode(span, tok)
}

func NewInvalidPatternNode(span *position.Span, tok *token.Token) PatternNode {
	return NewInvalidNode(span, tok)
}

func NewInvalidPatternExpressionNode(span *position.Span, tok *token.Token) PatternExpressionNode {
	return NewInvalidNode(span, tok)
}
