package types

import (
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

	DefineClass(docComment string, name string, parent Namespace, env *GlobalEnvironment) *Class
	DefineModule(docComment string, name string) *Module
	DefineMixin(docComment string, name string, env *GlobalEnvironment) *Mixin
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
