package types

import (
	"fmt"

	"github.com/elk-language/elk/lexer"
)

type Type interface {
	ToNonLiteral(*GlobalEnvironment) Type
	inspect() string
}

func InspectModifier(abstract, sealed bool) string {
	if abstract {
		return "abstract"
	}
	if sealed {
		return "sealed"
	}

	return "default"
}

func Inspect(typ Type) string {
	if typ == nil {
		return "void"
	}

	return typ.inspect()
}

func InspectInstanceVariable(name string) string {
	return fmt.Sprintf("@%s", name)
}

func InspectInstanceVariableWithColor(name string) string {
	return lexer.Colorize(InspectInstanceVariable(name))
}

func InspectInstanceVariableDeclaration(name string, typ Type) string {
	return fmt.Sprintf("var @%s: %s", name, Inspect(typ))
}

func InspectInstanceVariableDeclarationWithColor(name string, typ Type) string {
	return lexer.Colorize(InspectInstanceVariableDeclaration(name, typ))
}

func InspectWithColor(typ Type) string {
	return lexer.Colorize(Inspect(typ))
}

func GetMethod(typ Type, name string, env *GlobalEnvironment) *Method {
	typ = typ.ToNonLiteral(env)

	switch t := typ.(type) {
	case *Class:
		return t.MethodString(name)
	case *Module:
		return t.MethodString(name)
	}

	return nil
}
