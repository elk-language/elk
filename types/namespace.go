package types

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/value"
)

type Namespace interface {
	Type
	Name() string
	DocComment() string
	SetDocComment(string)
	AppendDocComment(string)
	Parent() Namespace
	SetParent(Namespace)
	Singleton() *SingletonClass
	IsAbstract() bool
	IsSealed() bool
	IsPrimitive() bool

	Constants() *TypeMap
	Constant(name value.Symbol) Type
	ConstantString(name string) Type
	DefineConstant(name string, val Type)

	Subtypes() *TypeMap
	Subtype(name value.Symbol) Type
	SubtypeString(name string) Type
	DefineSubtype(name string, val Type)

	Methods() *MethodMap
	Method(name value.Symbol) *Method
	MethodString(name string) *Method
	DefineMethod(docComment string, name string, params []*Parameter, returnType, throwType Type) *Method
	SetMethod(name string, method *Method)

	InstanceVariables() *TypeMap
	InstanceVariable(name value.Symbol) Type
	InstanceVariableString(name string) Type
	DefineInstanceVariable(name string, val Type)

	DefineClass(docComment string, primitive, abstract, sealed bool, name string, parent Namespace, env *GlobalEnvironment) *Class
	DefineModule(docComment string, name string) *Module
	DefineMixin(docComment string, abstract bool, name string, env *GlobalEnvironment) *Mixin
	DefineInterface(docComment string, name string, env *GlobalEnvironment) *Interface
}

func GetMethodInNamespace(namespace Namespace, name string) *Method {
	currentNamespace := namespace
	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		method := currentNamespace.MethodString(name)
		if method != nil {
			return method
		}
	}

	return nil
}

func GetInstanceVariableInNamespace(namespace Namespace, name string) (Type, Namespace) {
	currentNamespace := namespace
	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		ivar := currentNamespace.InstanceVariableString(name)
		if ivar != nil {
			return ivar, currentNamespace
		}
	}

	return nil, nil
}

func NamespacesAreEqual(left, right Namespace) bool {
	if left == right {
		return true
	}

	switch l := left.(type) {
	case *Mixin:
		switch r := right.(type) {
		case *MixinProxy:
			return l == r.Mixin
		}
	case *MixinProxy:
		switch r := right.(type) {
		case *Mixin:
			return l.Mixin == r
		case *MixinProxy:
			return l.Mixin == r.Mixin
		}
	case *Interface:
		switch r := right.(type) {
		case *InterfaceProxy:
			return l == r.Interface
		}
	case *InterfaceProxy:
		switch r := right.(type) {
		case *Interface:
			return l.Interface == r
		case *InterfaceProxy:
			return l.Interface == r.Interface
		}
	}

	return false
}

func GetConstantPath(fullConstantPath string) []string {
	return strings.Split(fullConstantPath, ":")
}

func GetConstantName(fullConstantPath string) string {
	constantPath := GetConstantPath(fullConstantPath)
	return constantPath[len(constantPath)-1]
}

// Serialise the type to Go code
func TypeToCode(typ Type) string {
	switch t := typ.(type) {
	case Any:
		return "types.Any{}"
	case Void, nil:
		return "types.Void{}"
	case Never:
		return "types.Never{}"
	case *Class:
		return fmt.Sprintf(
			"types.NameToType(%q, env)",
			t.name,
		)
	case *Mixin:
		return fmt.Sprintf(
			"types.NameToType(%q, env)",
			t.name,
		)
	case *Module:
		return fmt.Sprintf(
			"types.NameToType(%q, env)",
			t.name,
		)
	case *Interface:
		return fmt.Sprintf(
			"types.NameToType(%q, env)",
			t.name,
		)
	case *Nilable:
		return fmt.Sprintf(
			"types.NewNilable(%s)",
			TypeToCode(t.Type),
		)
	case *Union:
		buff := new(strings.Builder)
		buff.WriteString("types.NewUnion(")
		for _, element := range t.Elements {
			fmt.Fprintf(
				buff,
				"%s, ",
				TypeToCode(element),
			)
		}
		buff.WriteRune(')')
		return buff.String()
	case *Intersection:
		buff := new(strings.Builder)
		buff.WriteString("types.NewIntersection(")
		for _, element := range t.Elements {
			fmt.Fprintf(
				buff,
				"%s, ",
				TypeToCode(element),
			)
		}
		buff.WriteRune(')')
		return buff.String()
	case *SymbolLiteral:
		return fmt.Sprintf("types.NewSymbolLiteral(%q)", t.Value)
	case *StringLiteral:
		return fmt.Sprintf("types.NewStringLiteral(%q)", t.Value)
	case *CharLiteral:
		return fmt.Sprintf("types.NewCharLiteral(%q)", t.Value)
	case *FloatLiteral:
		return fmt.Sprintf("types.NewFloatLiteral(%q)", t.Value)
	case *Float32Literal:
		return fmt.Sprintf("types.NewFloat32Literal(%q)", t.Value)
	case *Float64Literal:
		return fmt.Sprintf("types.NewFloat64Literal(%q)", t.Value)
	case *IntLiteral:
		return fmt.Sprintf("types.NewIntLiteral(%q)", t.Value)
	case *Int64Literal:
		return fmt.Sprintf("types.NewInt64Literal(%q)", t.Value)
	case *Int32Literal:
		return fmt.Sprintf("types.NewInt32Literal(%q)", t.Value)
	case *Int16Literal:
		return fmt.Sprintf("types.NewInt16Literal(%q)", t.Value)
	case *Int8Literal:
		return fmt.Sprintf("types.NewInt8Literal(%q)", t.Value)
	case *UInt64Literal:
		return fmt.Sprintf("types.NewUInt64Literal(%q)", t.Value)
	case *UInt32Literal:
		return fmt.Sprintf("types.NewUInt32Literal(%q)", t.Value)
	case *UInt16Literal:
		return fmt.Sprintf("types.NewUInt16Literal(%q)", t.Value)
	case *UInt8Literal:
		return fmt.Sprintf("types.NewUInt8Literal(%q)", t.Value)
	default:
		panic(
			fmt.Sprintf("invalid type: %T", typ),
		)
	}
}

func NameToType(fullSubtypePath string, env *GlobalEnvironment) Type {
	subtypePath := GetConstantPath(fullSubtypePath)
	var namespace Namespace = env.Root
	var currentType Type = namespace
	for _, subtypeName := range subtypePath {
		if namespace == nil {
			panic(
				fmt.Sprintf(
					"`%s` is not a namespace",
					InspectWithColor(currentType),
				),
			)
		}
		currentType = namespace.SubtypeString(subtypeName)
		if currentType == nil {
			panic(
				fmt.Sprintf(
					"Undefined subtype `%s` in namespace `%s`",
					subtypeName,
					InspectWithColor(namespace),
				),
			)
		}

		namespace, _ = currentType.(Namespace)
	}

	return currentType
}

func NameToNamespace(fullSubtypePath string, env *GlobalEnvironment) Namespace {
	return NameToType(fullSubtypePath, env).(Namespace)
}

// iterate over every mixin that is included in the given namespace
func ForeachIncludedMixin(namespace Namespace, f func(*Mixin)) {
	currentNamespace := namespace.Parent()
	seenMixins := make(map[string]bool)

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		var mixin *Mixin
		switch n := currentNamespace.(type) {
		case *MixinProxy:
			mixin = n.Mixin
		default:
			continue
		}

		if seenMixins[mixin.name] {
			continue
		}

		f(mixin)

		seenMixins[mixin.name] = true
	}
}

// iterate over every interface that is implemented in the given namespace
func ForeachImplementedInterface(namespace Namespace, f func(*Interface)) {
	currentNamespace := namespace.Parent()
	seenInterfaces := make(map[string]bool)

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		var iface *Interface
		switch n := currentNamespace.(type) {
		case *InterfaceProxy:
			iface = n.Interface
		default:
			continue
		}

		if seenInterfaces[iface.name] {
			continue
		}

		f(iface)

		seenInterfaces[iface.name] = true
	}
}

func FindRootParent(namespace Namespace) Namespace {
	currentNamespace := namespace
	for {
		parent := currentNamespace.Parent()
		if parent == nil {
			return currentNamespace
		}
		currentNamespace = parent
	}
}

// Iterate over every method defined in the given namespace including the inherited ones
func ForeachMethod(namespace Namespace, f func(*Method)) {
	currentNamespace := namespace
	seenMethods := make(map[string]bool)

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		for _, method := range currentNamespace.Methods().Map {
			if seenMethods[method.Name] {
				continue
			}

			f(method)
			seenMethods[method.Name] = true
		}
	}
}

// Iterate over every instance variable defined in the given namespace including the inherited ones
func ForeachInstanceVariable(namespace Namespace, f func(name string, typ Type, namespace Namespace)) {
	currentNamespace := namespace
	seenIvars := make(map[value.Symbol]bool)

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		for name, typ := range currentNamespace.InstanceVariables().Map {
			if seenIvars[name] {
				continue
			}

			f(name.String(), typ, currentNamespace)
			seenIvars[name] = true
		}
	}
}
