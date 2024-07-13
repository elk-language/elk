package types

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
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
	DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, params []*Parameter, returnType, throwType Type) *Method
	SetMethod(name value.Symbol, method *Method)

	InstanceVariables() *TypeMap
	InstanceVariable(name value.Symbol) Type
	InstanceVariableString(name string) Type
	DefineInstanceVariable(name value.Symbol, val Type)

	DefineClass(docComment string, primitive, abstract, sealed bool, name string, parent Namespace, env *GlobalEnvironment) *Class
	DefineModule(docComment string, name string) *Module
	DefineMixin(docComment string, abstract bool, name string, env *GlobalEnvironment) *Mixin
	DefineInterface(docComment string, name string, env *GlobalEnvironment) *Interface
}

func GetMethodInNamespace(namespace Namespace, name value.Symbol) *Method {
	currentNamespace := namespace
	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		method := currentNamespace.Method(name)
		if method != nil {
			return method
		}
	}

	return nil
}

func GetInstanceVariableInNamespace(namespace Namespace, name value.Symbol) (Type, Namespace) {
	currentNamespace := namespace
	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		ivar := currentNamespace.InstanceVariable(name)
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
	return strings.Split(fullConstantPath, "::")
}

func GetConstantName(fullConstantPath string) string {
	constantPath := GetConstantPath(fullConstantPath)
	return constantPath[len(constantPath)-1]
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

// Iterate over every subtype
func ForeachSubtype(namespace Namespace, f func(name value.Symbol, typ Type)) {
	for name, typ := range namespace.Subtypes().Map {
		f(name, typ)
	}
}

// Iterate over every subtype, sorted by name
func ForeachSubtypeSorted(namespace Namespace, f func(name value.Symbol, typ Type)) {
	subtypes := namespace.Subtypes().Map
	names := symbol.SortKeys(subtypes)

	for _, name := range names {
		typ := subtypes[name]
		f(name, typ)
	}
}

// Iterate over every constant that is not a subtype
func ForeachConstant(namespace Namespace, f func(name value.Symbol, typ Type)) {
	for name, typ := range namespace.Constants().Map {
		if namespace.Subtype(name) != nil {
			continue
		}

		f(name, typ)
	}
}

// Iterate over every constant that is not a subtype, sorted by name
func ForeachConstantSorted(namespace Namespace, f func(name value.Symbol, typ Type)) {
	constants := namespace.Constants().Map
	names := symbol.SortKeys(constants)
	for _, name := range names {
		typ := constants[name]
		if namespace.Subtype(name) != nil {
			continue
		}

		f(name, typ)
	}
}

// Iterate over every method defined in the given namespace including the inherited ones
func ForeachMethod(namespace Namespace, f func(name value.Symbol, method *Method)) {
	currentNamespace := namespace
	seenMethods := make(map[value.Symbol]bool)

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		for name, method := range currentNamespace.Methods().Map {
			if seenMethods[name] {
				continue
			}

			f(name, method)
			seenMethods[name] = true
		}
	}
}

// Iterate over every method defined in the given namespace including the inherited ones, sorted by name
func ForeachMethodSorted(namespace Namespace, f func(name value.Symbol, method *Method)) {
	currentNamespace := namespace
	seenMethods := make(map[value.Symbol]bool)

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		methods := currentNamespace.Methods().Map
		names := symbol.SortKeys(methods)
		for _, name := range names {
			method := methods[name]
			if seenMethods[name] {
				continue
			}

			f(name, method)
			seenMethods[name] = true
		}
	}
}

// Iterate over every instance variable defined in the given namespace including the inherited ones
func ForeachInstanceVariable(namespace Namespace, f func(name value.Symbol, typ Type, namespace Namespace)) {
	currentNamespace := namespace
	seenIvars := make(map[value.Symbol]bool)

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		for name, typ := range currentNamespace.InstanceVariables().Map {
			if seenIvars[name] {
				continue
			}

			f(name, typ, currentNamespace)
			seenIvars[name] = true
		}
	}
}

// Iterate over every instance variable defined in the given namespace including the inherited ones
func ForeachInstanceVariableSorted(namespace Namespace, f func(name value.Symbol, typ Type, namespace Namespace)) {
	currentNamespace := namespace
	seenIvars := make(map[value.Symbol]bool)

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		ivars := currentNamespace.InstanceVariables().Map
		names := symbol.SortKeys(ivars)
		for _, name := range names {
			typ := ivars[name]
			if seenIvars[name] {
				continue
			}

			f(name, typ, currentNamespace)
			seenIvars[name] = true
		}
	}
}
