package types

import (
	"fmt"
	"iter"
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
	IsGeneric() bool

	TypeParameters() []*TypeParameter
	SetTypeParameters([]*TypeParameter)

	Constants() *TypeMap
	Constant(name value.Symbol) Type
	ConstantString(name string) Type
	DefineConstant(name value.Symbol, val Type)

	Subtypes() *TypeMap
	Subtype(name value.Symbol) Type
	SubtypeString(name string) Type
	DefineSubtype(name value.Symbol, val Type)

	Methods() *MethodMap
	Method(name value.Symbol) *Method
	MethodString(name string) *Method
	DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method
	SetMethod(name value.Symbol, method *Method)

	InstanceVariables() *TypeMap
	InstanceVariable(name value.Symbol) Type
	InstanceVariableString(name string) Type
	DefineInstanceVariable(name value.Symbol, val Type)

	DefineClass(docComment string, primitive, abstract, sealed bool, name value.Symbol, parent Namespace, env *GlobalEnvironment) *Class
	DefineModule(docComment string, name value.Symbol, env *GlobalEnvironment) *Module
	DefineMixin(docComment string, abstract bool, name value.Symbol, env *GlobalEnvironment) *Mixin
	DefineInterface(docComment string, name value.Symbol, env *GlobalEnvironment) *Interface
}

func implementInterface(target Namespace, iface *Interface) {
	headProxy, tailProxy := iface.CreateProxy()
	tailProxy.SetParent(target.Parent())
	target.SetParent(headProxy)
}

func ImplementInterface(target, interfaceNamespace Namespace) {
	switch implemented := interfaceNamespace.(type) {
	case *Interface:
		implementInterface(target, implemented)
	case *Generic:
		iface := implemented.Namespace.(*Interface)
		headProxy, tailProxy := iface.CreateProxy()
		head := NewGeneric(headProxy, implemented.TypeArguments)
		tailProxy.SetParent(target.Parent())
		target.SetParent(head)
	default:
		panic(fmt.Sprintf("wrong interface type: %T", interfaceNamespace))
	}
}

func includeMixin(target Namespace, mixin *Mixin) {
	proxy := NewMixinProxy(mixin, target.Parent())
	target.SetParent(proxy)
}

func IncludeMixin(target, includedNamespace Namespace) {
	switch included := includedNamespace.(type) {
	case *Mixin:
		includeMixin(target, included)
	case *Generic:
		includedMixin := included.Namespace.(*Mixin)
		proxy := NewMixinProxy(includedMixin, target.Parent())
		generic := NewGeneric(proxy, included.TypeArguments)
		target.SetParent(generic)
	default:
		panic(fmt.Sprintf("wrong mixin type: %T", includedNamespace))
	}
}

func NamespaceDeclaresInstanceVariables(namespace Namespace) bool {
	for parent := range Parents(namespace) {
		if parent.InstanceVariables().Len() > 0 {
			return true
		}
	}

	return false
}

func GetInstanceVariableInNamespace(namespace Namespace, name value.Symbol) (Type, Namespace) {
	for parent := range Parents(namespace) {
		ivar := parent.InstanceVariable(name)
		if ivar != nil {
			return ivar, parent
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

// Iterate over every parent of the given namespace (including itself).
func Parents(namespace Namespace) iter.Seq[Namespace] {
	return func(yield func(parent Namespace) bool) {
		if namespace == nil {
			return
		}
		namespaces := []Namespace{namespace}

		for i := 0; i < len(namespaces); i++ {
			currentParent := namespaces[i]
		parentLoop:
			for {
				if currentParent == nil {
					break
				}

				switch cn := currentParent.(type) {
				case *MixinProxy:
					mixinParent := cn.Parent()
					if mixinParent != nil {
						namespaces = append(namespaces, mixinParent)
					}
					if !yield(currentParent) {
						return
					}

					currentParent = cn.Mixin.parent
					continue parentLoop
				case *Generic:
					switch g := cn.Namespace.(type) {
					case *MixinProxy:
						genericParent := cn.Parent()
						if genericParent != nil {
							namespaces = append(namespaces, genericParent)
						}
						if !yield(currentParent) {
							return
						}

						currentParent = g.Mixin.parent
						continue parentLoop
					}
				}

				if !yield(currentParent) {
					return
				}

				currentParent = currentParent.Parent()
			}
		}
	}
}

// iterate over every mixin that is included in the given namespace
func IncludedMixins(namespace Namespace) iter.Seq[Namespace] {
	return func(yield func(mixin Namespace) bool) {
		seenMixins := make(map[string]bool)

		for parent := range Parents(namespace.Parent()) {
			switch n := parent.(type) {
			case *MixinProxy:
			case *Generic:
				if _, ok := n.Namespace.(*MixinProxy); !ok {
					continue
				}
			default:
				continue
			}

			name := parent.Name()
			if seenMixins[name] {
				continue
			}

			if !yield(parent) {
				return
			}

			seenMixins[name] = true
		}
	}
}

// iterate over every interface that is implemented in the given namespace
func ImplementedInterfaces(namespace Namespace) iter.Seq[Namespace] {
	return func(yield func(iface Namespace) bool) {
		seenInterfaces := make(map[string]bool)

		for parent := range Parents(namespace.Parent()) {
			switch n := parent.(type) {
			case *InterfaceProxy:
			case *Generic:
				if _, ok := n.Namespace.(*InterfaceProxy); !ok {
					continue
				}
			default:
				continue
			}

			name := parent.Name()
			if seenInterfaces[name] {
				continue
			}

			if !yield(parent) {
				return
			}

			seenInterfaces[name] = true
		}
	}
}

// Iterate over every subtype
func AllSubtypes(namespace Namespace) iter.Seq2[value.Symbol, Type] {
	return func(yield func(name value.Symbol, typ Type) bool) {
		for name, typ := range namespace.Subtypes().Map {
			if !yield(name, typ) {
				break
			}
		}
	}
}

// Iterate over every subtype, sorted by name
func SortedSubtypes(namespace Namespace) iter.Seq2[value.Symbol, Type] {
	return func(yield func(name value.Symbol, typ Type) bool) {
		subtypes := namespace.Subtypes().Map
		names := symbol.SortKeys(subtypes)

		for _, name := range names {
			typ := subtypes[name]
			if !yield(name, typ) {
				break
			}
		}
	}
}

// Iterate over every constant that is not a subtype
func AllConstants(namespace Namespace) iter.Seq2[value.Symbol, Type] {
	return func(yield func(name value.Symbol, typ Type) bool) {
		for name, typ := range namespace.Constants().Map {
			if namespace.Subtype(name) != nil {
				continue
			}

			if !yield(name, typ) {
				break
			}
		}
	}
}

// Iterate over every constant that is not a subtype, sorted by name
func SortedConstants(namespace Namespace) iter.Seq2[value.Symbol, Type] {
	return func(yield func(name value.Symbol, typ Type) bool) {
		constants := namespace.Constants().Map
		names := symbol.SortKeys(constants)
		for _, name := range names {
			typ := constants[name]
			if namespace.Subtype(name) != nil {
				continue
			}

			if !yield(name, typ) {
				break
			}
		}
	}
}

// Iterate over every method defined in the given namespace including the inherited ones
func AllMethods(namespace Namespace) iter.Seq2[value.Symbol, *Method] {
	return func(yield func(name value.Symbol, method *Method) bool) {
		seenMethods := make(map[value.Symbol]bool)

		for parent := range Parents(namespace) {
			for name, method := range parent.Methods().Map {
				if seenMethods[name] {
					continue
				}

				if !yield(name, method) {
					return
				}
				seenMethods[name] = true
			}
		}
	}
}

// Iterate over every method defined in the given namespace including the inherited ones, sorted by name
func SortedMethods(namespace Namespace) iter.Seq2[value.Symbol, *Method] {
	return func(yield func(name value.Symbol, method *Method) bool) {
		seenMethods := make(map[value.Symbol]bool)

		for parent := range Parents(namespace) {
			methods := parent.Methods().Map
			names := symbol.SortKeys(methods)
			for _, name := range names {
				method := methods[name]
				if seenMethods[name] {
					continue
				}

				if !yield(name, method) {
					return
				}
				seenMethods[name] = true
			}
		}
	}
}

// Iterate over every method defined directly under the given namespace
func OwnMethods(namespace Namespace) iter.Seq2[value.Symbol, *Method] {
	return func(yield func(name value.Symbol, method *Method) bool) {
		for name, method := range namespace.Methods().Map {
			if !yield(name, method) {
				break
			}
		}
	}
}

// Iterate over every method defined directly under the given namespace, sorted by name
func SortedOwnMethods(namespace Namespace) iter.Seq2[value.Symbol, *Method] {
	return func(yield func(name value.Symbol, method *Method) bool) {
		methods := namespace.Methods().Map
		names := symbol.SortKeys(methods)

		for _, name := range names {
			method := methods[name]
			if !yield(name, method) {
				break
			}
		}
	}
}

type InstanceVariable struct {
	Type      Type
	Namespace Namespace
}

// Iterate over every instance variable defined in the given namespace including the inherited ones
func AllInstanceVariables(namespace Namespace) iter.Seq2[value.Symbol, InstanceVariable] {
	return func(yield func(name value.Symbol, ivar InstanceVariable) bool) {
		seenIvars := make(map[value.Symbol]bool)

		for parent := range Parents(namespace) {
			for name, typ := range parent.InstanceVariables().Map {
				if seenIvars[name] {
					continue
				}

				if !yield(name, InstanceVariable{typ, parent}) {
					return
				}
				seenIvars[name] = true
			}
		}
	}
}

// Iterate over every instance variable defined in the given namespace including the inherited ones
func SortedInstanceVariables(namespace Namespace) iter.Seq2[value.Symbol, InstanceVariable] {
	return func(yield func(name value.Symbol, ivar InstanceVariable) bool) {
		seenIvars := make(map[value.Symbol]bool)

		for parent := range Parents(namespace) {
			ivars := parent.InstanceVariables().Map
			names := symbol.SortKeys(ivars)
			for _, name := range names {
				typ := ivars[name]
				if seenIvars[name] {
					continue
				}

				if !yield(name, InstanceVariable{typ, parent}) {
					return
				}
				seenIvars[name] = true
			}
		}
	}
}

// Iterate over every instance variable defined directly under the given namespace
func OwnInstanceVariables(namespace Namespace) iter.Seq2[value.Symbol, Type] {
	return func(yield func(name value.Symbol, typ Type) bool) {
		for name, typ := range namespace.InstanceVariables().Map {
			if !yield(name, typ) {
				break
			}
		}
	}
}

// Iterate over every instance variable defined directly under the given namespace, sorted by name
func SortedOwnInstanceVariables(namespace Namespace) iter.Seq2[value.Symbol, Type] {
	return func(yield func(name value.Symbol, typ Type) bool) {
		ivars := namespace.InstanceVariables().Map
		names := symbol.SortKeys(ivars)

		for _, name := range names {
			typ := ivars[name]
			if !yield(name, typ) {
				break
			}
		}
	}
}
