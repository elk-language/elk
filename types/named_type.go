package types

import (
	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/value/symbol"
)

type NamedType struct {
	docComment      string
	Name            string
	Type            Type
	declarationFile string
	useFiles        concurrent.Set[string]
}

var _ TypeWithUseFiles = &NamedType{}

func NewNamedType(docComment string, name string, typ Type, decl string) *NamedType {
	return &NamedType{
		docComment:      docComment,
		Name:            name,
		Type:            typ,
		declarationFile: decl,
		useFiles:        concurrent.MakeSet[string](),
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

func (n *NamedType) DeclarationFile() string {
	return n.declarationFile
}

func (n *NamedType) UseFiles() *concurrent.Set[string] {
	return &n.useFiles
}

func (c *NamedType) AddUseFile(decl string) {
	c.useFiles.Add(decl)
}

func (*NamedType) IsLiteral() bool {
	return false
}

func (n *NamedType) inspect() string {
	return n.Name
}

func (n *NamedType) Copy() *NamedType {
	return &NamedType{
		docComment:      n.docComment,
		Name:            n.Name,
		Type:            n.Type,
		declarationFile: n.declarationFile,
		useFiles:        n.useFiles,
	}
}

func (n *NamedType) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *NamedType {
	if newType, ok := NameToTypeOk(n.Name, newEnv); ok {
		return newType.(*NamedType)
	}

	newType := &NamedType{
		docComment:      n.docComment,
		Name:            n.Name,
		declarationFile: n.declarationFile,
		useFiles:        n.useFiles.CloneValue(),
	}

	classConstantPath := GetConstantPath(n.Name)
	parentNamespace := DeepCopyNamespacePath(classConstantPath[:len(classConstantPath)-1], oldEnv, newEnv)
	classConstantName := classConstantPath[len(classConstantPath)-1]
	parentNamespace.DefineSubtype(symbol.ToSymbol(classConstantName), newType)

	newType.Type = DeepCopyEnv(n.Type, oldEnv, newEnv)

	return newType
}
