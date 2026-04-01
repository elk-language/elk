package value

import (
	"strings"

	"github.com/elk-language/elk/concurrent"
)

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

var goTypeMap = concurrent.NewMap[string, *GoType]()

func NewGoType(name string) *GoType {
	return &GoType{
		Name: name,
	}
}

func FetchGoType(name string) *GoType {
	goTypeMap.Lock()
	defer goTypeMap.Unlock()

	typ, ok := goTypeMap.GetUnsafe(name)
	if ok {
		return typ
	}

	gt := NewGoType(name)
	goTypeMap.SetUnsafe(name, gt)
	return gt
}

func NewGenericGoType(name string, typeArgs []*GoType) *GoType {
	return &GoType{
		Name:     name,
		TypeArgs: typeArgs,
	}
}

func FetchGenericGoType(name string, typeArgs []*GoType) *GoType {
	goTypeMap.Lock()
	defer goTypeMap.Unlock()

	key := GoTypeKey(
		name,
		typeArgs,
	)
	typ, ok := goTypeMap.GetUnsafe(key)
	if ok {
		return typ
	}

	gt := NewGenericGoType(name, typeArgs)
	goTypeMap.SetUnsafe(name, gt)
	return gt
}

func (g *GoType) IsGeneric() bool {
	return len(g.TypeArgs) != 0
}

func (g *GoType) String() string {
	return GoTypeKey(g.Name, g.TypeArgs)
}

func GoTypeKey(name string, typeArgs []*GoType) string {
	if len(typeArgs) == 0 {
		return name
	}

	var b strings.Builder

	b.WriteString(name)
	b.WriteRune('[')
	for i, typeArg := range typeArgs {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(typeArg.String())
	}
	b.WriteRune(']')

	return b.String()
}
