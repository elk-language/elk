package types

import (
	"strings"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/value/symbol"
)

type GenericNamedType struct {
	docComment      string
	Name            string
	declarationFile string
	Type            Type
	TypeParameters  []*TypeParameter
	useFiles        concurrent.Set[string]
}

func NewGenericNamedType(docComment string, name string, typ Type, typeVars []*TypeParameter, decl string) *GenericNamedType {
	return &GenericNamedType{
		docComment:      docComment,
		Name:            name,
		Type:            typ,
		TypeParameters:  typeVars,
		declarationFile: decl,
		useFiles:        concurrent.MakeSet[string](),
	}
}

func (g *GenericNamedType) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(g, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(g, parent)
	}

	if g.Type.traverse(g, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, typeParam := range g.TypeParameters {
		if typeParam.traverse(g, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(g, parent)
}

func (g *GenericNamedType) ToNonLiteral(env *GlobalEnvironment) Type {
	return g
}

func (*GenericNamedType) IsLiteral() bool {
	return false
}

func (n *GenericNamedType) DocComment() string {
	return n.docComment
}

func (n *GenericNamedType) DeclarationFile() string {
	return n.declarationFile
}

func (n *GenericNamedType) UseFiles() *concurrent.Set[string] {
	return &n.useFiles
}

func (c *GenericNamedType) AddUseFile(decl string) {
	c.useFiles.Add(decl)
}

func (g *GenericNamedType) inspect() string {
	if len(g.TypeParameters) < 1 {
		return g.Name
	}

	buffer := new(strings.Builder)

	buffer.WriteString(g.Name)
	buffer.WriteRune('[')
	for i, typeVar := range g.TypeParameters {
		if i > 0 {
			buffer.WriteString(", ")
		}
		switch typeVar.Variance {
		case COVARIANT:
			buffer.WriteRune('+')
		case CONTRAVARIANT:
			buffer.WriteRune('-')
		case BIVARIANT:
			buffer.WriteString("+-")
		}

		buffer.WriteString(typeVar.Name.String())
		if !IsNever(typeVar.LowerBound) {
			buffer.WriteString(" > ")
			buffer.WriteString(Inspect(typeVar.LowerBound))
		}
		if !IsAny(typeVar.UpperBound) {
			buffer.WriteString(" < ")
			buffer.WriteString(Inspect(typeVar.UpperBound))
		}
	}
	buffer.WriteRune(']')
	return buffer.String()
}

func (g *GenericNamedType) Copy() *GenericNamedType {
	return &GenericNamedType{
		docComment:      g.docComment,
		declarationFile: g.declarationFile,
		useFiles:        g.useFiles,
		Name:            g.Name,
		Type:            g.Type,
		TypeParameters:  g.TypeParameters,
	}
}

func (g *GenericNamedType) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *GenericNamedType {
	if newType, ok := NameToTypeOk(g.Name, newEnv); ok {
		return newType.(*GenericNamedType)
	}

	newType := &GenericNamedType{
		Name:            g.Name,
		docComment:      g.docComment,
		declarationFile: g.declarationFile,
		useFiles:        g.useFiles.CloneValue(),
	}

	classConstantPath := GetConstantPath(g.Name)
	parentNamespace := DeepCopyNamespacePath(classConstantPath[:len(classConstantPath)-1], oldEnv, newEnv)
	classConstantName := classConstantPath[len(classConstantPath)-1]
	parentNamespace.DefineSubtype(symbol.ToSymbol(classConstantName), newType)

	newType.Type = DeepCopyEnv(g.Type, oldEnv, newEnv)
	newType.TypeParameters = TypeParametersDeepCopyEnv(g.TypeParameters, oldEnv, newEnv)

	return newType
}
