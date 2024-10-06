package types

import (
	"fmt"

	"github.com/elk-language/elk/value/symbol"
)

type FloatLiteral struct {
	Value      string
	isNegative bool
}

func (f *FloatLiteral) StringValue() string {
	return f.Value
}

func (f *FloatLiteral) IsNegative() bool {
	return f.isNegative
}

func (f *FloatLiteral) SetNegative(val bool) {
	f.isNegative = val
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
	if f.isNegative {
		return fmt.Sprintf("-%s", f.Value)
	}
	return f.Value
}

func (f *FloatLiteral) CopyNumeric() NumericLiteral {
	return &FloatLiteral{
		Value:      f.Value,
		isNegative: f.isNegative,
	}
}

type Float64Literal struct {
	Value      string
	isNegative bool
}

func (f *Float64Literal) StringValue() string {
	return f.Value
}

func (f *Float64Literal) IsNegative() bool {
	return f.isNegative
}

func (f *Float64Literal) SetNegative(val bool) {
	f.isNegative = val
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
	if f.isNegative {
		return fmt.Sprintf("-%sf64", f.Value)
	}
	return fmt.Sprintf("%sf64", f.Value)
}

func (f *Float64Literal) CopyNumeric() NumericLiteral {
	return &Float64Literal{
		Value:      f.Value,
		isNegative: f.isNegative,
	}
}

type Float32Literal struct {
	Value      string
	isNegative bool
}

func (f *Float32Literal) StringValue() string {
	return f.Value
}

func (f *Float32Literal) IsNegative() bool {
	return f.isNegative
}

func (f *Float32Literal) SetNegative(val bool) {
	f.isNegative = val
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
	if f.isNegative {
		return fmt.Sprintf("-%sf32", f.Value)
	}
	return fmt.Sprintf("%sf32", f.Value)
}

func (f *Float32Literal) CopyNumeric() NumericLiteral {
	return &Float32Literal{
		Value:      f.Value,
		isNegative: f.isNegative,
	}
}

type BigFloatLiteral struct {
	Value      string
	isNegative bool
}

func (f *BigFloatLiteral) StringValue() string {
	return f.Value
}

func (f *BigFloatLiteral) IsNegative() bool {
	return f.isNegative
}

func (f *BigFloatLiteral) SetNegative(val bool) {
	f.isNegative = val
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
	if f.isNegative {
		return fmt.Sprintf("-%sbf", f.Value)
	}
	return fmt.Sprintf("%sbf", f.Value)
}

func (f *BigFloatLiteral) CopyNumeric() NumericLiteral {
	return &BigFloatLiteral{
		Value:      f.Value,
		isNegative: f.isNegative,
	}
}
