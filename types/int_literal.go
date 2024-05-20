package types

import (
	"fmt"

	"github.com/elk-language/elk/value/symbol"
)

type IntLiteral struct {
	Value string
}

func NewIntLiteral(value string) *IntLiteral {
	return &IntLiteral{
		Value: value,
	}
}

func (i *IntLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int)
}

func (i *IntLiteral) inspect() string {
	return fmt.Sprintf("Std::Int(%s)", i.Value)
}

type Int64Literal struct {
	Value string
}

func NewInt64Literal(value string) *Int64Literal {
	return &Int64Literal{
		Value: value,
	}
}

func (i *Int64Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int64)
}

func (i *Int64Literal) inspect() string {
	return fmt.Sprintf("%si64", i.Value)
}

type Int32Literal struct {
	Value string
}

func NewInt32Literal(value string) *Int32Literal {
	return &Int32Literal{
		Value: value,
	}
}

func (i *Int32Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int32)
}

func (i *Int32Literal) inspect() string {
	return fmt.Sprintf("%si32", i.Value)
}

type Int16Literal struct {
	Value string
}

func NewInt16Literal(value string) *Int16Literal {
	return &Int16Literal{
		Value: value,
	}
}

func (i *Int16Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int16)
}

func (i *Int16Literal) inspect() string {
	return fmt.Sprintf("%si16", i.Value)
}

type Int8Literal struct {
	Value string
}

func NewInt8Literal(value string) *Int8Literal {
	return &Int8Literal{
		Value: value,
	}
}

func (i *Int8Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int8)
}

func (i *Int8Literal) inspect() string {
	return fmt.Sprintf("%si8", i.Value)
}

type UInt64Literal struct {
	Value string
}

func NewUInt64Literal(value string) *UInt64Literal {
	return &UInt64Literal{
		Value: value,
	}
}

func (i *UInt64Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt64)
}

func (i *UInt64Literal) inspect() string {
	return fmt.Sprintf("%su64", i.Value)
}

type UInt32Literal struct {
	Value string
}

func NewUInt32Literal(value string) *UInt32Literal {
	return &UInt32Literal{
		Value: value,
	}
}

func (i *UInt32Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt32)
}

func (i *UInt32Literal) inspect() string {
	return fmt.Sprintf("%su32", i.Value)
}

type UInt16Literal struct {
	Value string
}

func NewUInt16Literal(value string) *UInt16Literal {
	return &UInt16Literal{
		Value: value,
	}
}

func (i *UInt16Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt16)
}

func (i *UInt16Literal) inspect() string {
	return fmt.Sprintf("%su16", i.Value)
}

type UInt8Literal struct {
	Value string
}

func NewUInt8Literal(value string) *UInt8Literal {
	return &UInt8Literal{
		Value: value,
	}
}

func (i *UInt8Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt8)
}

func (i *UInt8Literal) inspect() string {
	return fmt.Sprintf("%su8", i.Value)
}
