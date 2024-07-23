package types

import (
	"fmt"

	"github.com/elk-language/elk/value/symbol"
)

type IntLiteral struct {
	Value string
}

func (i *IntLiteral) StringValue() string {
	return i.Value
}

func NewIntLiteral(value string) *IntLiteral {
	return &IntLiteral{
		Value: value,
	}
}

func (i *IntLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int)
}

func (*IntLiteral) IsLiteral() bool {
	return true
}

func (i *IntLiteral) inspect() string {
	return i.Value
}

type Int64Literal struct {
	Value string
}

func (i *Int64Literal) StringValue() string {
	return i.Value
}

func NewInt64Literal(value string) *Int64Literal {
	return &Int64Literal{
		Value: value,
	}
}

func (i *Int64Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int64)
}

func (*Int64Literal) IsLiteral() bool {
	return true
}

func (i *Int64Literal) inspect() string {
	return fmt.Sprintf("%si64", i.Value)
}

type Int32Literal struct {
	Value string
}

func (i *Int32Literal) StringValue() string {
	return i.Value
}

func NewInt32Literal(value string) *Int32Literal {
	return &Int32Literal{
		Value: value,
	}
}

func (i *Int32Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int32)
}

func (*Int32Literal) IsLiteral() bool {
	return true
}

func (i *Int32Literal) inspect() string {
	return fmt.Sprintf("%si32", i.Value)
}

type Int16Literal struct {
	Value string
}

func (i *Int16Literal) StringValue() string {
	return i.Value
}

func NewInt16Literal(value string) *Int16Literal {
	return &Int16Literal{
		Value: value,
	}
}

func (i *Int16Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int16)
}

func (*Int16Literal) IsLiteral() bool {
	return true
}

func (i *Int16Literal) inspect() string {
	return fmt.Sprintf("%si16", i.Value)
}

type Int8Literal struct {
	Value string
}

func (i *Int8Literal) StringValue() string {
	return i.Value
}

func NewInt8Literal(value string) *Int8Literal {
	return &Int8Literal{
		Value: value,
	}
}

func (i *Int8Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int8)
}

func (*Int8Literal) IsLiteral() bool {
	return true
}

func (i *Int8Literal) inspect() string {
	return fmt.Sprintf("%si8", i.Value)
}

type UInt64Literal struct {
	Value string
}

func (i *UInt64Literal) StringValue() string {
	return i.Value
}

func NewUInt64Literal(value string) *UInt64Literal {
	return &UInt64Literal{
		Value: value,
	}
}

func (i *UInt64Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt64)
}

func (*UInt64Literal) IsLiteral() bool {
	return true
}

func (i *UInt64Literal) inspect() string {
	return fmt.Sprintf("%su64", i.Value)
}

type UInt32Literal struct {
	Value string
}

func (i *UInt32Literal) StringValue() string {
	return i.Value
}

func NewUInt32Literal(value string) *UInt32Literal {
	return &UInt32Literal{
		Value: value,
	}
}

func (i *UInt32Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt32)
}

func (*UInt32Literal) IsLiteral() bool {
	return true
}

func (i *UInt32Literal) inspect() string {
	return fmt.Sprintf("%su32", i.Value)
}

type UInt16Literal struct {
	Value string
}

func (i *UInt16Literal) StringValue() string {
	return i.Value
}

func NewUInt16Literal(value string) *UInt16Literal {
	return &UInt16Literal{
		Value: value,
	}
}

func (i *UInt16Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt16)
}

func (*UInt16Literal) IsLiteral() bool {
	return true
}

func (i *UInt16Literal) inspect() string {
	return fmt.Sprintf("%su16", i.Value)
}

type UInt8Literal struct {
	Value string
}

func (i *UInt8Literal) StringValue() string {
	return i.Value
}

func NewUInt8Literal(value string) *UInt8Literal {
	return &UInt8Literal{
		Value: value,
	}
}

func (i *UInt8Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt8)
}

func (*UInt8Literal) IsLiteral() bool {
	return true
}

func (i *UInt8Literal) inspect() string {
	return fmt.Sprintf("%su8", i.Value)
}
