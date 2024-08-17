package types

import (
	"strings"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type TypeArgument struct {
	Type     Type
	Variance Variance
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

func NewTypeArguments(m map[value.Symbol]*TypeArgument, order []value.Symbol) *TypeArguments {
	return &TypeArguments{
		ArgumentMap:   m,
		ArgumentOrder: order,
	}
}

type Generic struct {
	Type Type
	*TypeArguments
}

func NewGeneric(typ Type, typeArgs *TypeArguments) *Generic {
	generic := &Generic{
		Type:          typ,
		TypeArguments: typeArgs,
	}
	typeArgs.ArgumentMap[symbol.M_self] = NewTypeArgument(
		generic,
		INVARIANT,
	)

	return generic
}

func (g *Generic) ToNonLiteral(env *GlobalEnvironment) Type {
	return g
}

func (*Generic) IsLiteral() bool {
	return false
}

func (g *Generic) inspect() string {
	buffer := new(strings.Builder)

	buffer.WriteString(Inspect(g.Type))
	buffer.WriteRune('[')
	for i, argName := range g.ArgumentOrder {
		arg := g.ArgumentMap[argName]
		if i > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString(Inspect(arg.Type))
	}
	buffer.WriteRune(']')
	return buffer.String()
}