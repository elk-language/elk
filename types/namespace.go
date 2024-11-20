package types

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/elk-language/elk/ds"
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
	SetSingleton(*SingletonClass)
	IsAbstract() bool
	IsSealed() bool
	IsPrimitive() bool
	IsGeneric() bool
	IsDefined() bool
	SetDefined(bool)

	TypeParameters() []*TypeParameter
	SetTypeParameters([]*TypeParameter)

	Constants() ConstantMap
	Constant(name value.Symbol) (Constant, bool)
	ConstantString(name string) (Constant, bool)
	DefineConstant(name value.Symbol, val Type)
	DefineConstantWithFullName(name value.Symbol, fullName string, val Type)

	Subtypes() ConstantMap
	Subtype(name value.Symbol) (Constant, bool)
	SubtypeString(name string) (Constant, bool)
	MustSubtype(name string) Type
	DefineSubtype(name value.Symbol, val Type)
	DefineSubtypeWithFullName(name value.Symbol, fullName string, val Type)

	Methods() MethodMap
	Method(name value.Symbol) *Method
	MethodString(name string) *Method
	DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method
	SetMethod(name value.Symbol, method *Method)

	InstanceVariables() TypeMap
	InstanceVariable(name value.Symbol) Type
	InstanceVariableString(name string) Type
	DefineInstanceVariable(name value.Symbol, val Type)

	MethodAliases() MethodAliasMap
	SetMethodAlias(name value.Symbol, method *Method)

	DefineClass(docComment string, primitive, abstract, sealed bool, name value.Symbol, parent Namespace, env *GlobalEnvironment) *Class
	DefineModule(docComment string, name value.Symbol, env *GlobalEnvironment) *Module
	DefineMixin(docComment string, abstract bool, name value.Symbol, env *GlobalEnvironment) *Mixin
	DefineInterface(docComment string, name value.Symbol, env *GlobalEnvironment) *Interface
}

func TypeParametersDeepCopyEnv(typeParameters []*TypeParameter, oldEnv, newEnv *GlobalEnvironment) []*TypeParameter {
	newTypeParameters := make([]*TypeParameter, len(typeParameters))
	for name, typeParam := range typeParameters {
		newTypeParameters[name] = typeParam.DeepCopyEnv(oldEnv, newEnv)
	}
	return newTypeParameters
}

func ConstantsDeepCopyEnv(constants ConstantMap, oldEnv, newEnv *GlobalEnvironment) ConstantMap {
	newConstants := make(ConstantMap, len(constants))
	for constName, constant := range constants {
		newConstants[constName] = Constant{
			FullName: constant.FullName,
			Type:     DeepCopyEnv(constant.Type, oldEnv, newEnv),
		}
	}
	return newConstants
}

func TypesDeepCopyEnv(types TypeMap, oldEnv, newEnv *GlobalEnvironment) TypeMap {
	newTypes := make(TypeMap, len(types))
	for typeName, typ := range types {
		newTypes[typeName] = DeepCopyEnv(typ, oldEnv, newEnv)
	}
	return newTypes
}

func MethodsDeepCopyEnv(methods MethodMap, oldEnv, newEnv *GlobalEnvironment) MethodMap {
	newMethods := make(MethodMap, len(methods))
	for methodName, method := range methods {
		newMethods[methodName] = method.DeepCopyEnv(oldEnv, newEnv)
	}
	return newMethods
}

func NamespaceHasAnyDefinableMethods(namespace Namespace) bool {
	for _, method := range namespace.Methods() {
		if method.IsDefinable() {
			return true
		}
	}

	for _, alias := range namespace.MethodAliases() {
		if alias.IsDefinable() {
			return true
		}
	}

	return false
}

func ConstructTypeArgumentsFromTypeParameterUpperBounds(typeParams []*TypeParameter) *TypeArguments {
	typeArgMap := make(TypeArgumentMap, len(typeParams))
	typeArgOrder := make([]value.Symbol, len(typeParams))

	for i, typeParam := range typeParams {
		arg := typeParam.UpperBound

		typeArg := NewTypeArgument(
			arg,
			typeParam.Variance,
		)
		typeArgMap[typeParam.Name] = typeArg
		typeArgOrder[i] = typeParam.Name
	}

	return NewTypeArguments(
		typeArgMap,
		typeArgOrder,
	)
}

func ConstructTypeArgumentsFromTypeParameterUpperBoundsAndVariance(typeParams []*TypeParameter, variance Variance) *TypeArguments {
	typeArgMap := make(TypeArgumentMap, len(typeParams))
	typeArgOrder := make([]value.Symbol, len(typeParams))

	for i, typeParam := range typeParams {
		arg := typeParam.UpperBound

		typeArg := NewTypeArgument(
			arg,
			variance,
		)
		typeArgMap[typeParam.Name] = typeArg
		typeArgOrder[i] = typeParam.Name
	}

	return NewTypeArguments(
		typeArgMap,
		typeArgOrder,
	)
}

func implementInterface(target Namespace, iface *Interface) {
	proxy := NewInterfaceProxy(iface, target.Parent())
	target.SetParent(proxy)
}

func ImplementInterface(target, interfaceNamespace Namespace) {
	switch implemented := interfaceNamespace.(type) {
	case *Interface:
		implementInterface(target, implemented)
	case *Generic:
		iface := implemented.Namespace.(*Interface)
		proxy := NewInterfaceProxy(iface, target.Parent())
		generic := NewGeneric(proxy, implemented.TypeArguments)
		target.SetParent(generic)
	default:
		panic(fmt.Sprintf("wrong interface type: %T", interfaceNamespace))
	}
}

func includeMixin(target Namespace, mixin *Mixin) {
	proxy := NewMixinProxy(mixin, target.Parent())
	target.SetParent(proxy)
}

func IncludeMixinWithWhere(target Namespace, mixin *Mixin, where []*TypeParameter) *MixinWithWhere {
	proxy := NewMixinProxy(mixin, target.Parent())
	mixinWithWhere := NewMixinWithWhere(proxy, target, where)
	target.SetParent(mixinWithWhere)
	return mixinWithWhere
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
		if len(parent.InstanceVariables()) > 0 {
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

func NameToTypeOk(fullSubtypePath string, env *GlobalEnvironment) (Type, bool) {
	if env.Root == nil {
		return nil, false
	}

	subtypePath := GetConstantPath(fullSubtypePath)
	var namespace Namespace = env.Root
	var currentType Type = namespace
	for _, subtypeName := range subtypePath {
		if namespace == nil {
			return nil, false
		}
		constant, ok := namespace.SubtypeString(subtypeName)
		if !ok {
			return nil, false
		}
		currentType = constant.Type

		namespace, _ = currentType.(Namespace)
	}

	return currentType, true
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
		constant, ok := namespace.SubtypeString(subtypeName)
		if !ok {
			panic(
				fmt.Sprintf(
					"Undefined subtype `%s` in namespace `%s`",
					subtypeName,
					InspectWithColor(namespace),
				),
			)
		}
		currentType = constant.Type

		namespace, _ = currentType.(Namespace)
	}

	return currentType
}

func NameToNamespace(fullSubtypePath string, env *GlobalEnvironment) Namespace {
	return NameToType(fullSubtypePath, env).(Namespace)
}

// Iterate over direct parents of the given namespace (including itself).
func SimpleParents(namespace Namespace) iter.Seq[Namespace] {
	return func(yield func(parent Namespace) bool) {
		for parent := namespace; parent != nil; parent = parent.Parent() {
			if !yield(parent) {
				return
			}
		}
	}
}

// Iterate over every parent of the given namespace (including itself).
func Parents(namespace Namespace) iter.Seq[Namespace] {
	return func(yield func(parent Namespace) bool) {
		if namespace == nil {
			return
		}
		namespaces := []Namespace{namespace}
		seenParents := make(ds.Set[string])

		for len(namespaces) > 0 {
			// Pop an element from the stack
			currentParent := namespaces[len(namespaces)-1]
			namespaces = namespaces[0 : len(namespaces)-1]

		parentLoop:
			for {
				if currentParent == nil {
					break
				}

				name := currentParent.Name()
				if len(name) > 0 && seenParents.Contains(name) {
					currentParent = currentParent.Parent()
					continue
				}
				seenParents.Add(name)

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
				case *InterfaceProxy:
					ifaceParent := cn.Parent()
					if ifaceParent != nil {
						namespaces = append(namespaces, ifaceParent)
					}
					if !yield(currentParent) {
						return
					}

					currentParent = cn.Interface.parent
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
					case *InterfaceProxy:
						genericParent := cn.Parent()
						if genericParent != nil {
							namespaces = append(namespaces, genericParent)
						}
						if !yield(currentParent) {
							return
						}

						currentParent = g.Interface.parent
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
		seenMixins := make(ds.Set[string])

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
			if seenMixins.Contains(name) {
				continue
			}

			if !yield(parent) {
				return
			}

			seenMixins.Add(name)
		}
	}
}

// iterate backward over an iterator
func Backward[T any](iterator iter.Seq[T]) iter.Seq[T] {
	return func(yield func(element T) bool) {
		for _, element := range slices.Backward(slices.Collect(iterator)) {
			if !yield(element) {
				return
			}
		}
	}
}

// iterate over every mixin that is directly included in the given namespace
func DirectlyIncludedMixins(namespace Namespace) iter.Seq[Namespace] {
	return func(yield func(mixin Namespace) bool) {
		seenMixins := make(ds.Set[string])

		for parent := namespace.Parent(); parent != nil; parent = parent.Parent() {
			switch n := parent.(type) {
			case *MixinProxy:
			case *Generic:
				if _, ok := n.Namespace.(*MixinProxy); !ok {
					continue
				}
			case *Class:
				return
			default:
				continue
			}

			name := parent.Name()
			if seenMixins.Contains(name) {
				continue
			}

			if !yield(parent) {
				return
			}

			seenMixins.Add(name)
		}
	}
}

// iterate over every mixin that is directly included and
// every interface that is directly implemented in the given namespace
func DirectlyIncludedAndImplemented(namespace Namespace) iter.Seq[Namespace] {
	return func(yield func(mixin Namespace) bool) {
		seenMixins := make(ds.Set[string])

		for parent := namespace.Parent(); parent != nil; parent = parent.Parent() {
			switch n := parent.(type) {
			case *MixinProxy, *InterfaceProxy, *MixinWithWhere:
			case *Generic:
				switch n.Namespace.(type) {
				case *MixinProxy, *InterfaceProxy:
				default:
					continue
				}
			case *Class:
				return
			default:
				continue
			}

			name := parent.Name()
			if seenMixins.Contains(name) {
				continue
			}

			if !yield(parent) {
				return
			}

			seenMixins.Add(name)
		}
	}
}

// iterate over every interface that is implemented in the given namespace
func ImplementedInterfaces(namespace Namespace) iter.Seq[Namespace] {
	return func(yield func(iface Namespace) bool) {
		seenInterfaces := make(ds.Set[string])

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
			if seenInterfaces.Contains(name) {
				continue
			}

			if !yield(parent) {
				return
			}

			seenInterfaces.Add(name)
		}
	}
}

// iterate over every interface that is directly implemented in the given namespace
func DirectlyImplementedInterfaces(namespace Namespace) iter.Seq[Namespace] {
	return func(yield func(iface Namespace) bool) {
		seenInterfaces := make(ds.Set[string])

		for parent := namespace.Parent(); parent != nil; parent = parent.Parent() {
			switch n := parent.(type) {
			case *InterfaceProxy:
			case *Generic:
				if _, ok := n.Namespace.(*InterfaceProxy); !ok {
					continue
				}
			case *Class:
				return
			default:
				continue
			}

			name := parent.Name()
			if seenInterfaces.Contains(name) {
				continue
			}

			if !yield(parent) {
				return
			}

			seenInterfaces.Add(name)
		}
	}
}

// Iterate over every subtype
func AllSubtypes(namespace Namespace) iter.Seq2[value.Symbol, Constant] {
	return func(yield func(name value.Symbol, constant Constant) bool) {
		for name, typ := range namespace.Subtypes() {
			if !yield(name, typ) {
				break
			}
		}
	}
}

// Iterate over every subtype, sorted by name
func SortedSubtypes(namespace Namespace) iter.Seq2[value.Symbol, Constant] {
	return func(yield func(name value.Symbol, constant Constant) bool) {
		subtypes := namespace.Subtypes()
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
func AllConstants(namespace Namespace) iter.Seq2[value.Symbol, Constant] {
	return func(yield func(name value.Symbol, constant Constant) bool) {
		for name, typ := range namespace.Constants() {
			if _, ok := namespace.Subtype(name); ok {
				continue
			}

			if !yield(name, typ) {
				break
			}
		}
	}
}

// Iterate over every constant that is not a subtype, sorted by name
func SortedConstants(namespace Namespace) iter.Seq2[value.Symbol, Constant] {
	return func(yield func(name value.Symbol, constant Constant) bool) {
		constants := namespace.Constants()
		names := symbol.SortKeys(constants)
		for _, name := range names {
			constant := constants[name]
			if _, ok := namespace.Subtype(name); ok {
				continue
			}

			if !yield(name, constant) {
				break
			}
		}
	}
}

// Iterate over every method defined in the given namespace including the inherited ones
func AllMethods(namespace Namespace) iter.Seq2[value.Symbol, *Method] {
	return func(yield func(name value.Symbol, method *Method) bool) {
		seenMethods := make(ds.Set[value.Symbol])

		for parent := range Parents(namespace) {
			for name, method := range parent.Methods() {
				if seenMethods.Contains(name) {
					continue
				}

				if !yield(name, method) {
					return
				}
				seenMethods.Add(name)
			}
		}
	}
}

// Iterate over every method defined in the given namespace including the inherited ones, sorted by name
func SortedMethods(namespace Namespace) iter.Seq2[value.Symbol, *Method] {
	return func(yield func(name value.Symbol, method *Method) bool) {
		seenMethods := make(ds.Set[value.Symbol])

		for parent := range Parents(namespace) {
			methods := parent.Methods()
			names := symbol.SortKeys(methods)
			for _, name := range names {
				method := methods[name]
				if seenMethods.Contains(name) {
					continue
				}

				if !yield(name, method) {
					return
				}
				seenMethods.Add(name)
			}
		}
	}
}

// Iterate over every method defined directly under the given namespace
func OwnMethods(namespace Namespace) iter.Seq2[value.Symbol, *Method] {
	return func(yield func(name value.Symbol, method *Method) bool) {
		for name, method := range namespace.Methods() {
			if !yield(name, method) {
				break
			}
		}
	}
}

// Iterate over every method defined directly under the given namespace, sorted by name
func SortedOwnMethods(namespace Namespace) iter.Seq2[value.Symbol, *Method] {
	return func(yield func(name value.Symbol, method *Method) bool) {
		methods := namespace.Methods()
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
		seenIvars := make(ds.Set[value.Symbol])

		for parent := range Parents(namespace) {
			for name, typ := range parent.InstanceVariables() {
				if seenIvars.Contains(name) {
					continue
				}

				if !yield(name, InstanceVariable{typ, parent}) {
					return
				}
				seenIvars.Add(name)
			}
		}
	}
}

// Iterate over every instance variable defined in the given namespace including the inherited ones
func SortedInstanceVariables(namespace Namespace) iter.Seq2[value.Symbol, InstanceVariable] {
	return func(yield func(name value.Symbol, ivar InstanceVariable) bool) {
		seenIvars := make(ds.Set[value.Symbol])

		for parent := range Parents(namespace) {
			ivars := parent.InstanceVariables()
			names := symbol.SortKeys(ivars)
			for _, name := range names {
				typ := ivars[name]
				if seenIvars.Contains(name) {
					continue
				}

				if !yield(name, InstanceVariable{typ, parent}) {
					return
				}
				seenIvars.Add(name)
			}
		}
	}
}

// Iterate over every instance variable defined directly under the given namespace
func OwnInstanceVariables(namespace Namespace) iter.Seq2[value.Symbol, Type] {
	return func(yield func(name value.Symbol, typ Type) bool) {
		for name, typ := range namespace.InstanceVariables() {
			if !yield(name, typ) {
				break
			}
		}
	}
}

// Iterate over every instance variable defined directly under the given namespace, sorted by name
func SortedOwnInstanceVariables(namespace Namespace) iter.Seq2[value.Symbol, Type] {
	return func(yield func(name value.Symbol, typ Type) bool) {
		ivars := namespace.InstanceVariables()
		names := symbol.SortKeys(ivars)

		for _, name := range names {
			typ := ivars[name]
			if !yield(name, typ) {
				break
			}
		}
	}
}
