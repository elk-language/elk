package types

import (
	"fmt"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value/symbol"
)

type Type interface {
	ToNonLiteral(*GlobalEnvironment) Type
	inspect() string
}

func CanBeFalsy(typ Type, env *GlobalEnvironment) bool {
	switch t := typ.(type) {
	case *Nilable:
		return true
	case *Class:
		if t == env.StdSubtype(symbol.Bool) || t == env.StdSubtype(symbol.False) || t == env.StdSubtype(symbol.Nil) {
			return true
		}
		return false
	case *Union:
		for _, element := range t.Elements {
			if CanBeFalsy(element, env) {
				return true
			}
		}
		return false
	case *Intersection:
		for _, element := range t.Elements {
			if CanBeFalsy(element, env) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func CanBeTruthy(typ Type, env *GlobalEnvironment) bool {
	switch t := typ.(type) {
	case *Nilable:
		return CanBeTruthy(t.Type, env)
	case *Class:
		if t == env.StdSubtype(symbol.False) || t == env.StdSubtype(symbol.Nil) {
			return false
		}
		return true
	case *Union:
		for _, element := range t.Elements {
			if CanBeTruthy(element, env) {
				return true
			}
		}
		return false
	case *Intersection:
		for _, element := range t.Elements {
			if CanBeTruthy(element, env) {
				return true
			}
		}
		return false
	default:
		return false
	}
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
