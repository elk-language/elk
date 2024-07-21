package types

import (
	"fmt"

	"github.com/elk-language/elk/value/symbol"
)

type FloatLiteral struct {
	Value string
}

func (f *FloatLiteral) StringValue() string {
	return f.Value
}

func NewFloatLiteral(value string) *FloatLiteral {
	return &FloatLiteral{
		Value: value,
	}
}

func (*FloatLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Float)
}

func (*FloatLiteral) IsLiteral() bool {
	return true
}

func (f *FloatLiteral) inspect() string {
	return fmt.Sprintf("Std::Float(%s)", f.Value)
}

type Float64Literal struct {
	Value string
}

func (f *Float64Literal) StringValue() string {
	return f.Value
}

func NewFloat64Literal(value string) *Float64Literal {
	return &Float64Literal{
		Value: value,
	}
}

func (*Float64Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Float64)
}

func (*Float64Literal) IsLiteral() bool {
	return true
}

func (f *Float64Literal) inspect() string {
	return fmt.Sprintf("%sf64", f.Value)
}

type Float32Literal struct {
	Value string
}

func (f *Float32Literal) StringValue() string {
	return f.Value
}

func NewFloat32Literal(value string) *Float32Literal {
	return &Float32Literal{
		Value: value,
	}
}

func (*Float32Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Float32)
}

func (*Float32Literal) IsLiteral() bool {
	return true
}

func (f *Float32Literal) inspect() string {
	return fmt.Sprintf("%sf32", f.Value)
}

type BigFloatLiteral struct {
	Value string
}

func (f *BigFloatLiteral) StringValue() string {
	return f.Value
}

func NewBigFloatLiteral(value string) *BigFloatLiteral {
	return &BigFloatLiteral{
		Value: value,
	}
}

func (*BigFloatLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.BigFloat)
}

func (*BigFloatLiteral) IsLiteral() bool {
	return true
}

func (f *BigFloatLiteral) inspect() string {
	return fmt.Sprintf("%sbf", f.Value)
}
