package types

import (
	"encoding/binary"
	"fmt"
	"iter"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value/symbol"
)

type TypeArgument struct {
	Type     Ref[Type]
	Variance Variance
}

// Create a shallow copy
func (t *TypeArgument) Copy() *TypeArgument {
	return &TypeArgument{
		Type:     t.Type,
		Variance: t.Variance,
	}
}

func NewTypeArgument(typ Type, variance Variance) *TypeArgument {
	t := &TypeArgument{
		Type:     ToRef(typ),
		Variance: variance,
	}
	return t
}

func NewTypeArgumentWithRef(typ Ref[Type], variance Variance) *TypeArgument {
	t := &TypeArgument{
		Type:     typ,
		Variance: variance,
	}
	return t
}

type TypeArgumentMap map[symbol.Symbol]*TypeArgument

func (t TypeArgumentMap) HasAllTypeParams(typeParams []Ref[*TypeParameter]) bool {
	for _, typeParam := range typeParams {
		_, ok := t[typeParam.Get().Name]
		if !ok {
			return false
		}
	}

	return true
}

type TypeArguments struct {
	ArgumentMap   TypeArgumentMap
	ArgumentOrder []symbol.Symbol
}

func CreateTypeArgumentOrderFromTypeParams(typeParams []Ref[*TypeParameter]) []symbol.Symbol {
	order := make([]symbol.Symbol, len(typeParams))
	for i, typeParam := range typeParams {
		order[i] = typeParam.Get().Name
	}
	return order
}

// Create a shallow copy
func (t *TypeArguments) Copy() *TypeArguments {
	return &TypeArguments{
		ArgumentMap:   t.ArgumentMap,
		ArgumentOrder: t.ArgumentOrder,
	}
}

func (t *TypeArguments) DeleteUnnecessaryArgs() {
	newMap := make(TypeArgumentMap, len(t.ArgumentOrder))
	for name, arg := range t.AllArguments() {
		newMap[name] = arg
	}
	t.ArgumentMap = newMap
}

// Create a deep copy with ArgumentMap
func (t *TypeArguments) DeepCopy() *TypeArguments {
	newMap := make(TypeArgumentMap, len(t.ArgumentMap))
	for key, val := range t.ArgumentMap {
		newMap[key] = val.Copy()
	}
	return &TypeArguments{
		ArgumentMap:   newMap,
		ArgumentOrder: t.ArgumentOrder,
	}
}

// Iterates over every type argument in definition order.
func (t *TypeArguments) AllArguments() iter.Seq2[symbol.Symbol, *TypeArgument] {
	return func(yield func(name symbol.Symbol, arg *TypeArgument) bool) {
		for _, name := range t.ArgumentOrder {
			arg := t.ArgumentMap[name]
			if !yield(name, arg) {
				break
			}
		}
	}
}

// Get the type argument under the given index
func (t *TypeArguments) Get(i int) *TypeArgument {
	return t.ArgumentMap[t.ArgumentOrder[i]]
}

func (t *TypeArguments) Len() int {
	return len(t.ArgumentOrder)
}

func NewTypeArguments(m TypeArgumentMap, order []symbol.Symbol) *TypeArguments {
	return &TypeArguments{
		ArgumentMap:   m,
		ArgumentOrder: order,
	}
}

type Generic struct {
	Namespace Ref[Namespace]
	*TypeArguments
	id ID
}

var _ Namespace = &Generic{}

func IsGeneric(typ Type) bool {
	_, ok := typ.(*Generic)
	return ok
}

func (g *Generic) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(g, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(g, parent)
	}
}

func (g *Generic) Name() string {
	return g.Namespace.Get().Name()
}

func (g *Generic) DocComment() string {
	return g.Namespace.Get().DocComment()
}

func (g *Generic) SetDocComment(comment string) {
	g.Namespace.Get().SetDocComment(comment)
}

func (g *Generic) AppendDocComment(comment string) {
	g.Namespace.Get().AppendDocComment(comment)
}

func (g *Generic) Parent() Namespace {
	return g.Namespace.Get().Parent()
}

func (g *Generic) SetParent(parent Namespace) {
	g.Namespace.Get().SetParent(parent)
}

func (g *Generic) Singleton() *SingletonClass {
	return g.Namespace.Get().Singleton()
}

func (g *Generic) SetSingleton(singleton *SingletonClass) {
	g.Namespace.Get().SetSingleton(singleton)
}

func (g *Generic) IsAbstract() bool {
	return g.Namespace.Get().IsAbstract()
}

func (g *Generic) IsSealed() bool {
	return g.Namespace.Get().IsSealed()
}

func (g *Generic) IsPrimitive() bool {
	return g.Namespace.Get().IsPrimitive()
}

func (g *Generic) IsImmutable() bool {
	return g.Namespace.Get().IsImmutable()
}

func (g *Generic) IsGeneric() bool {
	return g.Namespace.Get().IsGeneric()
}

func (g *Generic) IsDefined() bool {
	return g.Namespace.Get().IsDefined()
}

func (g *Generic) SetDefined(defined bool) {
	g.Namespace.Get().SetDefined(defined)
}

func (g *Generic) IsNative() bool {
	return g.Namespace.Get().IsNative()
}

func (g *Generic) TypeParameters() []Ref[*TypeParameter] {
	return g.Namespace.Get().TypeParameters()
}

func (g *Generic) SetTypeParameters(typeParams []Ref[*TypeParameter]) {
	g.Namespace.Get().SetTypeParameters(typeParams)
}

func (g *Generic) Constants() ConstantMap {
	return g.Namespace.Get().Constants()
}

func (g *Generic) Constant(name symbol.Symbol) (Constant, bool) {
	return g.Namespace.Get().Constant(name)
}

func (g *Generic) ConstantString(name string) (Constant, bool) {
	return g.Namespace.Get().ConstantString(name)
}

func (g *Generic) DefineConstant(name symbol.Symbol, val Type) {
	g.Namespace.Get().DefineConstant(name, val)
}

func (g *Generic) DefineConstantWithFullName(name symbol.Symbol, fullName string, val Type) {
	g.Namespace.Get().DefineConstantWithFullName(name, fullName, val)
}

func (g *Generic) Subtypes() ConstantMap {
	return g.Namespace.Get().Subtypes()
}

func (g *Generic) Subtype(name symbol.Symbol) (Constant, bool) {
	return g.Namespace.Get().Subtype(name)
}

func (g *Generic) SubtypeString(name string) (Constant, bool) {
	return g.Namespace.Get().SubtypeString(name)
}

func (g *Generic) MustSubtype(name symbol.Symbol) Type {
	return g.Namespace.Get().MustSubtype(name)
}

func (g *Generic) MustSubtypeString(name string) Type {
	return g.Namespace.Get().MustSubtypeString(name)
}

func (g *Generic) DefineSubtype(name symbol.Symbol, val Type) {
	g.Namespace.Get().DefineSubtype(name, val)
}

func (g *Generic) DefineSubtypeWithFullName(name symbol.Symbol, fullName string, val Type) {
	g.Namespace.Get().DefineSubtypeWithFullName(name, fullName, val)
}

func (g *Generic) Methods() MethodMap {
	return g.Namespace.Get().Methods()
}

func (g *Generic) Method(name symbol.Symbol) *Method {
	return g.Namespace.Get().Method(name)
}

func (g *Generic) MethodString(name string) *Method {
	return g.Namespace.Get().MethodString(name)
}

func (g *Generic) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	return g.Namespace.Get().DefineMethod(docComment, flags, name, typeParams, params, returnType, throwType)
}

func (g *Generic) SetMethod(name symbol.Symbol, method *Method) {
	g.Namespace.Get().SetMethod(name, method)
}

func (g *Generic) InstanceVariables() InstanceVariableMap {
	return g.Namespace.Get().InstanceVariables()
}

func (g *Generic) InstanceVariable(name symbol.Symbol) *InstanceVariable {
	return g.Namespace.Get().InstanceVariable(name)
}

func (g *Generic) InstanceVariableString(name string) *InstanceVariable {
	return g.Namespace.Get().InstanceVariableString(name)
}

func (g *Generic) DefineInstanceVariable(name symbol.Symbol, ivar *InstanceVariable) {
	g.Namespace.Get().DefineInstanceVariable(name, ivar)
}

func (g *Generic) DefineClass(docComment string, primitive, abstract, sealed, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	return g.Namespace.Get().DefineClass(docComment, primitive, abstract, sealed, noinit, immutable, name, parent)
}

func (g *Generic) DefineModule(docComment string, name symbol.Symbol) *Module {
	return g.Namespace.Get().DefineModule(docComment, name)
}

func (g *Generic) DefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	return g.Namespace.Get().DefineMixin(docComment, abstract, name)
}

func (g *Generic) DefineInterface(docComment string, name symbol.Symbol) *Interface {
	return g.Namespace.Get().DefineInterface(docComment, name)
}

func (c *Generic) ToRef() Ref[*Generic] {
	return ToRef(c)
}

func (c *Generic) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("generic:")
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(c.Namespace))
	d.Write(b)

	for _, typeArg := range c.TypeArguments.AllArguments() {
		binary.LittleEndian.PutUint32(b, uint32(typeArg.Type))
		d.Write(b)
	}

	return d.Sum64()
}

func (c *Generic) EqualAny(other any) bool {
	o, ok := other.(*Generic)
	if !ok {
		return false
	}

	if c.id > 0 {
		return c.id == o.ID()
	}

	if c.Namespace != o.Namespace ||
		len(c.ArgumentMap) != len(o.ArgumentMap) {
		return false
	}

	for i := range len(c.ArgumentOrder) {
		cArg := c.ArgumentMap[c.ArgumentOrder[i]]
		oArg := o.ArgumentMap[o.ArgumentOrder[i]]
		if cArg.Type != oArg.Type {
			return false
		}
	}

	return true
}

func (c *Generic) ID() ID {
	return c.id
}

func (c *Generic) SetID(id ID) {
	c.id = id
}

func NewGeneric(typ Namespace, typeArgs *TypeArguments) *Generic {
	generic := &Generic{
		Namespace:     ToRef(typ),
		TypeArguments: typeArgs,
	}

	return generic
}

// If the namespace is not generic returns itself.
// If the namespace is generic returns a new generic with the namespace's type parameters as type arguments.
func GetDefaultNamespaceGenericInstance(namespace Namespace) Namespace {
	if !namespace.IsGeneric() {
		return namespace
	}

	return NewGenericWithTypeParameters(namespace)
}

func NewGenericWithTypeParameters(namespace Namespace) *Generic {
	return NewGeneric(
		namespace,
		ConstructTypeArgumentsFromTypeParameters(
			namespace.TypeParameters(),
		),
	)
}

func NewGenericWithUpperBoundTypeArgs(namespace Namespace) *Generic {
	return NewGeneric(
		namespace,
		ConstructTypeArgumentsFromTypeParameterUpperBounds(
			namespace.TypeParameters(),
		),
	)
}

func NewGenericWithUpperBoundTypeArgsAndVariance(namespace Namespace, variance Variance) *Generic {
	return NewGeneric(
		namespace,
		ConstructTypeArgumentsFromTypeParameterUpperBoundsAndVariance(
			namespace.TypeParameters(),
			variance,
		),
	)
}

func NewGenericWithTypeArgs(namespace Namespace, args ...Type) *Generic {
	if len(namespace.TypeParameters()) != len(args) {
		panic(fmt.Sprintf("invalid type argument count in new generic for %s, expected %d, got %d", Inspect(namespace), len(namespace.TypeParameters()), len(args)))
	}

	typeArgMap := make(TypeArgumentMap, len(args))
	typeArgOrder := make([]symbol.Symbol, len(args))

	for i, typeParamRef := range namespace.TypeParameters() {
		typeParam := typeParamRef.Get()
		arg := args[i]

		typeArg := NewTypeArgument(
			arg,
			typeParam.Variance,
		)
		typeArgMap[typeParam.Name] = typeArg
		typeArgOrder[i] = typeParam.Name
	}

	return NewGeneric(
		namespace,
		NewTypeArguments(
			typeArgMap,
			typeArgOrder,
		),
	)
}

func NewGenericWithTypeArgRefs(namespace Namespace, args ...Ref[Type]) *Generic {
	if len(namespace.TypeParameters()) != len(args) {
		panic(fmt.Sprintf("invalid type argument count in new generic for %s, expected %d, got %d", Inspect(namespace), len(namespace.TypeParameters()), len(args)))
	}

	typeArgMap := make(TypeArgumentMap, len(args))
	typeArgOrder := make([]symbol.Symbol, len(args))

	for i, typeParamRef := range namespace.TypeParameters() {
		typeParam := typeParamRef.Get()
		arg := args[i]

		typeArg := NewTypeArgumentWithRef(
			arg,
			typeParam.Variance,
		)
		typeArgMap[typeParam.Name] = typeArg
		typeArgOrder[i] = typeParam.Name
	}

	return NewGeneric(
		namespace,
		NewTypeArguments(
			typeArgMap,
			typeArgOrder,
		),
	)
}

func NewGenericWithVariance(namespace Namespace, variance Variance, args ...Type) *Generic {
	if len(namespace.TypeParameters()) != len(args) {
		panic(fmt.Sprintf("invalid type argument count in new generic, expected %d, got %d", len(namespace.TypeParameters()), len(args)))
	}

	typeArgMap := make(TypeArgumentMap, len(args))
	typeArgOrder := make([]symbol.Symbol, len(args))

	for i, typeParamRef := range namespace.TypeParameters() {
		typeParam := typeParamRef.Get()
		arg := args[i]

		typeArg := NewTypeArgument(
			arg,
			variance,
		)
		typeArgMap[typeParam.Name] = typeArg
		typeArgOrder[i] = typeParam.Name
	}

	return NewGeneric(
		namespace,
		NewTypeArguments(
			typeArgMap,
			typeArgOrder,
		),
	)
}

func (g *Generic) FixVariance() {
	if g == nil {
		return
	}
	for _, typeParamRef := range g.Namespace.Get().TypeParameters() {
		typeParam := typeParamRef.Get()
		arg := g.ArgumentMap[typeParam.Name]
		arg.Variance = typeParam.Variance
	}
}

func (g *Generic) ToNonLiteral() Type {
	return g
}

func (*Generic) IsLiteral() bool {
	return false
}

func (g *Generic) inspect() string {
	buffer := new(strings.Builder)

	buffer.WriteString(Inspect(g.Namespace.Get()))
	buffer.WriteRune('[')
	first := true
	for _, arg := range g.AllArguments() {
		if !first {
			buffer.WriteString(", ")
		} else {
			first = false
		}

		buffer.WriteString(Inspect(arg.Type.Get()))
	}
	buffer.WriteRune(']')
	return buffer.String()
}

func (g *Generic) Copy() *Generic {
	return &Generic{
		Namespace:     g.Namespace,
		TypeArguments: g.TypeArguments,
	}
}

func (g *Generic) CopyType() Type {
	return g.Copy()
}
