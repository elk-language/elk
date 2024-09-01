package types

import (
	"fmt"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value/symbol"
)

type Type interface {
	ToNonLiteral(*GlobalEnvironment) Type
	IsLiteral() bool
	inspect() string
}

func CanBeFalsy(typ Type, env *GlobalEnvironment) bool {
	switch t := typ.(type) {
	case *Nilable, Nil, False, Bool, Nothing, Void:
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
	case *NamedType:
		return CanBeFalsy(t.Type, env)
	default:
		return false
	}
}

func CanBeTruthy(typ Type, env *GlobalEnvironment) bool {
	switch t := typ.(type) {
	case *Nilable:
		return CanBeTruthy(t.Type, env)
	case Nil, False:
		return false
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
	case *NamedType:
		return CanBeTruthy(t.Type, env)
	default:
		return true
	}
}

func InspectModifier(abstract, sealed, primitive bool) string {
	if abstract {
		if primitive {
			return "abstract primitive"
		}
		return "abstract"
	}
	if sealed {
		if primitive {
			return "sealed primitive"
		}
		return "sealed"
	}

	if primitive {
		return "primitive"
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

func I(typ Type) string {
	return InspectWithColor(typ)
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
