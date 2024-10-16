package types

import (
	"fmt"
	"iter"
	"strings"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type TypeArgument struct {
	Type     Type
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
	return &TypeArgument{
		Type:     typ,
		Variance: variance,
	}
}

type TypeArguments struct {
	ArgumentMap   map[value.Symbol]*TypeArgument
	ArgumentOrder []value.Symbol
}

// Create a shallow copy
func (t *TypeArguments) Copy() *TypeArguments {
	return &TypeArguments{
		ArgumentMap:   t.ArgumentMap,
		ArgumentOrder: t.ArgumentOrder,
	}
}

// Create a deep copy with ArgumentMap
func (t *TypeArguments) DeepCopy() *TypeArguments {
	newMap := make(map[value.Symbol]*TypeArgument, len(t.ArgumentMap))
	for key, val := range t.ArgumentMap {
		newMap[key] = val.Copy()
	}
	return &TypeArguments{
		ArgumentMap:   newMap,
		ArgumentOrder: t.ArgumentOrder,
	}
}

// Iterates over every type argument in definition order.
func (t *TypeArguments) AllArguments() iter.Seq2[value.Symbol, *TypeArgument] {
	return func(yield func(name value.Symbol, arg *TypeArgument) bool) {
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

func NewTypeArguments(m map[value.Symbol]*TypeArgument, order []value.Symbol) *TypeArguments {
	return &TypeArguments{
		ArgumentMap:   m,
		ArgumentOrder: order,
	}
}

type Generic struct {
	Namespace
	*TypeArguments
}

func IsGeneric(typ Type) bool {
	_, ok := typ.(*Generic)
	return ok
}

func NewGeneric(typ Namespace, typeArgs *TypeArguments) *Generic {
	generic := &Generic{
		Namespace:     typ,
		TypeArguments: typeArgs,
	}
	typeArgs.ArgumentMap[symbol.L_self] = NewTypeArgument(
		generic,
		INVARIANT,
	)

	return generic
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
		panic(fmt.Sprintf("invalid type argument count in new generic, expected %d, got %d", len(namespace.TypeParameters()), len(args)))
	}

	typeArgMap := make(map[value.Symbol]*TypeArgument, len(args))
	typeArgOrder := make([]value.Symbol, len(args))

	for i, typeParam := range namespace.TypeParameters() {
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

func NewGenericWithVariance(namespace Namespace, variance Variance, args ...Type) *Generic {
	if len(namespace.TypeParameters()) != len(args) {
		panic(fmt.Sprintf("invalid type argument count in new generic, expected %d, got %d", len(namespace.TypeParameters()), len(args)))
	}

	typeArgMap := make(map[value.Symbol]*TypeArgument, len(args))
	typeArgOrder := make([]value.Symbol, len(args))

	for i, typeParam := range namespace.TypeParameters() {
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
	for _, typeParam := range g.Namespace.TypeParameters() {
		arg := g.ArgumentMap[typeParam.Name]
		arg.Variance = typeParam.Variance
	}
}

func (g *Generic) ToNonLiteral(env *GlobalEnvironment) Type {
	return g
}

func (*Generic) IsLiteral() bool {
	return false
}

func (g *Generic) inspect() string {
	buffer := new(strings.Builder)

	buffer.WriteString(Inspect(g.Namespace))
	buffer.WriteRune('[')
	first := true
	for _, arg := range g.AllArguments() {
		if !first {
			buffer.WriteString(", ")
		} else {
			first = false
		}

		buffer.WriteString(Inspect(arg.Type))
	}
	buffer.WriteRune(']')
	return buffer.String()
}
