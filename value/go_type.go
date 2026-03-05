package value

import "strings"

// Represents a native Go type
type GoType struct {
	Name     string
	TypeArgs []*GoType
}

func (g *GoType) Equal(other *GoType) bool {
	if g == other {
		return true
	}
	if g.Name != other.Name {
		return false
	}

	if len(g.TypeArgs) != len(other.TypeArgs) {
		return false
	}

	for i := range len(g.TypeArgs) {
		if !g.TypeArgs[i].Equal(other.TypeArgs[i]) {
			return false
		}
	}

	return true
}

func NewGoType(name string) *GoType {
	return &GoType{
		Name: name,
	}
}

func NewGenericGoType(name string, typeArgs []*GoType) *GoType {
	return &GoType{
		Name:     name,
		TypeArgs: typeArgs,
	}
}

func (g *GoType) IsGeneric() bool {
	return len(g.TypeArgs) != 0
}

func (g *GoType) String() string {
	if len(g.TypeArgs) == 0 {
		return g.Name
	}

	var b strings.Builder

	b.WriteString(g.Name)
	b.WriteRune('[')
	for i, typeArg := range g.TypeArgs {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(typeArg.String())
	}
	b.WriteRune(']')

	return b.String()
}
