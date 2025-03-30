package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an instance variable eg. `@foo`
type InstanceVariableNode struct {
	TypedNodeBase
	Value string
}

func (n *InstanceVariableNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InstanceVariableNode)
	if !ok {
		return false
	}

	return n.Value != o.Value ||
		!n.span.Equal(o.span)
}

func (n *InstanceVariableNode) String() string {
	return fmt.Sprintf("@%s", n.Value)
}

func (*InstanceVariableNode) IsStatic() bool {
	return false
}

// Create an instance variable node eg. `@foo`.
func NewInstanceVariableNode(span *position.Span, val string) *InstanceVariableNode {
	return &InstanceVariableNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

func (*InstanceVariableNode) Class() *value.Class {
	return value.InstanceVariableNodeClass
}

func (*InstanceVariableNode) DirectClass() *value.Class {
	return value.InstanceVariableNodeClass
}

func (i *InstanceVariableNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::InstanceVariableNode{&: %p, value: %s}", i, value.String(i.Value).Inspect())
}

func (p *InstanceVariableNode) Error() string {
	return p.Inspect()
}
