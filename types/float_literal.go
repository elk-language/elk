package types

import (
	"fmt"

	"github.com/elk-language/elk/value/symbol"
)

type FloatLiteral struct {
	Value string
}

func NewFloatLiteral(value string) *FloatLiteral {
	return &FloatLiteral{
		Value: value,
	}
}

func (f *FloatLiteral) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *FloatLiteral:
		return f.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Float)
	default:
		return false
	}
}

func (*FloatLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Float)
}

func (f *FloatLiteral) Inspect() string {
	return fmt.Sprintf("Std::Float(%s)", f.Value)
}

type Float64Literal struct {
	Value string
}

func NewFloat64Literal(value string) *Float64Literal {
	return &Float64Literal{
		Value: value,
	}
}

func (f *Float64Literal) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *Float64Literal:
		return f.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Float64)
	default:
		return false
	}
}

func (*Float64Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Float64)
}

func (f *Float64Literal) Inspect() string {
	return fmt.Sprintf("%sf64", f.Value)
}

type Float32Literal struct {
	Value string
}

func NewFloat32Literal(value string) *Float32Literal {
	return &Float32Literal{
		Value: value,
	}
}

func (f *Float32Literal) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *Float32Literal:
		return f.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Float32)
	default:
		return false
	}
}

func (*Float32Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Float32)
}

func (f *Float32Literal) Inspect() string {
	return fmt.Sprintf("%sf32", f.Value)
}

type BigFloatLiteral struct {
	Value string
}

func NewBigFloatLiteral(value string) *BigFloatLiteral {
	return &BigFloatLiteral{
		Value: value,
	}
}

func (f *BigFloatLiteral) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *BigFloatLiteral:
		return f.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.BigFloat)
	default:
		return false
	}
}

func (*BigFloatLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.BigFloat)
}

func (f *BigFloatLiteral) Inspect() string {
	return fmt.Sprintf("%sbf", f.Value)
}
