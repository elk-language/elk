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
	return fmt.Sprintf("Std::AST::InstanceVariableNode{&: %p, value: %s}", i, value.String(i.Value).Inspect())
}

func (p *InstanceVariableNode) Error() string {
	return p.Inspect()
}
