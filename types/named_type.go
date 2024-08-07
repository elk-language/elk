package types

import "github.com/elk-language/elk/position"

type AstNode interface {
	position.SpanInterface
}

type NamedType struct {
	Name string
	Type Type
	Node AstNode
}

func NewNamedType(name string, typ Type) *NamedType {
	return &NamedType{
		Name: name,
		Type: typ,
	}
}

func (n *NamedType) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func (*NamedType) IsLiteral() bool {
	return false
}

func (n *NamedType) inspect() string {
	return n.Name
}
