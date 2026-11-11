package types

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value/symbol"
)

type NamedType struct {
	docComment  string
	Name        string
	Type        Type
	declaration *position.Location
}

func NewNamedType(docComment string, name string, typ Type, decl *position.Location) *NamedType {
	return &NamedType{
		docComment:  docComment,
		Name:        name,
		Type:        typ,
		declaration: decl,
	}
}

func (n *NamedType) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(n, parent)
	}
}

func (n *NamedType) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func (n *NamedType) DocComment() string {
	return n.docComment
}

func (n *NamedType) Declaration() *position.Location {
	return n.declaration
}

func (*NamedType) IsLiteral() bool {
	return false
}

func (n *NamedType) inspect() string {
	return n.Name
}

func (n *NamedType) Copy() *NamedType {
	return &NamedType{
		docComment:  n.docComment,
		Name:        n.Name,
		Type:        n.Type,
		declaration: n.declaration,
	}
}

func (n *NamedType) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *NamedType {
	if newType, ok := NameToTypeOk(n.Name, newEnv); ok {
		return newType.(*NamedType)
	}

	newType := &NamedType{
		docComment:  n.docComment,
		Name:        n.Name,
		declaration: n.declaration,
	}

	classConstantPath := GetConstantPath(n.Name)
	parentNamespace := DeepCopyNamespacePath(classConstantPath[:len(classConstantPath)-1], oldEnv, newEnv)
	classConstantName := classConstantPath[len(classConstantPath)-1]
	parentNamespace.DefineSubtype(symbol.ToSymbol(classConstantName), newType)

	newType.Type = DeepCopyEnv(n.Type, oldEnv, newEnv)

	return newType
}
