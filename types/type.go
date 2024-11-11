package types

import (
	"fmt"
	"reflect"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value/symbol"
)

type Type interface {
	ToNonLiteral(*GlobalEnvironment) Type
	IsLiteral() bool
	inspect() string
}

func DeepCopyNamespacePath(constantPath []string, oldEnv, newEnv *GlobalEnvironment) Namespace {
	var newNamespace Namespace = ToNamespaceInterface(newEnv.Root)
	var oldNamespace Namespace = ToNamespaceInterface(oldEnv.Root)
	var newCurrentType Type = ToTypeInterface(newEnv.Root)
	var oldCurrentType Type = ToTypeInterface(oldEnv.Root)

	for _, subtypeName := range constantPath {
		oldSubtype, _ := oldNamespace.SubtypeString(subtypeName)
		oldCurrentType = oldSubtype.Type
		oldNamespace = ToNamespaceInterface(oldCurrentType.(Namespace))

		newSubtype, ok := newNamespace.SubtypeString(subtypeName)
		if !ok {
			newCurrentType = DeepCopyEnv(oldNamespace, oldEnv, newEnv)
		} else {
			newCurrentType = newSubtype.Type
		}
		newNamespace = ToNamespaceInterface(newCurrentType.(Namespace))
	}

	return newNamespace
}

func IsPointerNil(val any) bool {
	if val == nil {
		return true
	}

	value := reflect.ValueOf(val)
	kind := value.Kind()
	return kind == reflect.Pointer && value.IsNil()
}

func ToTypeInterface[T Type](typ T) Type {
	if IsPointerNil(typ) {
		return nil
	}

	return typ
}

func ToNamespaceInterface[T Namespace](typ T) Namespace {
	if IsPointerNil(typ) {
		return nil
	}

	return typ
}

var counter uint64

func DeepCopyEnv(t Type, oldEnv, newEnv *GlobalEnvironment) Type {
	counter++
	if counter > 1000 {
		panic("dupa")
	}
	fmt.Printf("deep copy: %s\n", I(t))
	defer fmt.Printf("returned\n")
	switch t := t.(type) {
	case *Module:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *Class:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *Mixin:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *Method:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *Interface:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *SingletonClass:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *MixinProxy:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *InterfaceProxy:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *MixinWithWhere:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *Nilable:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *InstanceOf:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *SingletonOf:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *Generic:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *Not:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *NamedType:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *GenericNamedType:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *Union:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *Intersection:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *ConstantPlaceholder:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *ModulePlaceholder:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *NamespacePlaceholder:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *UsingBufferNamespace:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *TypeParamNamespace:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	case *Closure:
		return ToTypeInterface(t.DeepCopyEnv(oldEnv, newEnv))
	default:
		return ToTypeInterface(t)
	}
}

func CanBeFalsy(typ Type, env *GlobalEnvironment) bool {
	switch t := typ.(type) {
	case *Nilable, Nil, False, Bool, Untyped, Void:
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
