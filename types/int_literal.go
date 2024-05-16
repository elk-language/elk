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

func (i *IntLiteral) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *IntLiteral:
		return i.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Int)
	default:
		return false
	}
}

func (i *IntLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int)
}

func (i *IntLiteral) Inspect() string {
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

func (i *Int64Literal) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *Int64Literal:
		return i.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Int64)
	default:
		return false
	}
}

func (i *Int64Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int64)
}

func (i *Int64Literal) Inspect() string {
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

func (i *Int32Literal) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *Int32Literal:
		return i.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Int32)
	default:
		return false
	}
}

func (i *Int32Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int32)
}

func (i *Int32Literal) Inspect() string {
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

func (i *Int16Literal) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *Int16Literal:
		return i.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Int16)
	default:
		return false
	}
}

func (i *Int16Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int16)
}

func (i *Int16Literal) Inspect() string {
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

func (i *Int8Literal) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *Int8Literal:
		return i.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Int8)
	default:
		return false
	}
}

func (i *Int8Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int8)
}

func (i *Int8Literal) Inspect() string {
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

func (i *UInt64Literal) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *UInt64Literal:
		return i.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.UInt64)
	default:
		return false
	}
}

func (i *UInt64Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt64)
}

func (i *UInt64Literal) Inspect() string {
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

func (i *UInt32Literal) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *UInt32Literal:
		return i.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.UInt32)
	default:
		return false
	}
}

func (i *UInt32Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt32)
}

func (i *UInt32Literal) Inspect() string {
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

func (i *UInt16Literal) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *UInt16Literal:
		return i.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.UInt16)
	default:
		return false
	}
}

func (i *UInt16Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt16)
}

func (i *UInt16Literal) Inspect() string {
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

func (i *UInt8Literal) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *UInt8Literal:
		return i.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.UInt8)
	default:
		return false
	}
}

func (i *UInt8Literal) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt8)
}

func (i *UInt8Literal) Inspect() string {
	return fmt.Sprintf("%su8", i.Value)
}
