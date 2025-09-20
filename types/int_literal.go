package types

import (
	"fmt"

	"github.com/elk-language/elk/value/symbol"
)

type IntLiteral struct {
	Value      string
	isNegative bool
}

func (i *IntLiteral) StringValue() string {
	return i.Value
}

func (i *IntLiteral) IsNegative() bool {
	return i.isNegative
}

func (i *IntLiteral) SetNegative(val bool) {
	i.isNegative = val
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
	if i.isNegative {
		return fmt.Sprintf("-%s", i.Value)
	}
	return i.Value
}

func (i *IntLiteral) CopyNumeric() NumericLiteral {
	return &IntLiteral{
		Value:      i.Value,
		isNegative: i.isNegative,
	}
}

type Int64Literal struct {
	Value      string
	isNegative bool
}

func (i *Int64Literal) StringValue() string {
	return i.Value
}

func (i *Int64Literal) IsNegative() bool {
	return i.isNegative
}

func (i *Int64Literal) SetNegative(val bool) {
	i.isNegative = val
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
	if i.isNegative {
		return fmt.Sprintf("-%si64", i.Value)
	}
	return fmt.Sprintf("%si64", i.Value)
}

func (i *Int64Literal) CopyNumeric() NumericLiteral {
	return &Int64Literal{
		Value:      i.Value,
		isNegative: i.isNegative,
	}
}

type Int32Literal struct {
	Value      string
	isNegative bool
}

func (i *Int32Literal) StringValue() string {
	return i.Value
}

func (i *Int32Literal) IsNegative() bool {
	return i.isNegative
}

func (i *Int32Literal) SetNegative(val bool) {
	i.isNegative = val
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
	if i.isNegative {
		return fmt.Sprintf("-%si32", i.Value)
	}
	return fmt.Sprintf("%si32", i.Value)
}

func (i *Int32Literal) CopyNumeric() NumericLiteral {
	return &Int32Literal{
		Value:      i.Value,
		isNegative: i.isNegative,
	}
}

type Int16Literal struct {
	Value      string
	isNegative bool
}

func (i *Int16Literal) StringValue() string {
	return i.Value
}

func (i *Int16Literal) IsNegative() bool {
	return i.isNegative
}

func (i *Int16Literal) SetNegative(val bool) {
	i.isNegative = val
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
	if i.isNegative {
		return fmt.Sprintf("-%si16", i.Value)
	}
	return fmt.Sprintf("%si16", i.Value)
}

func (i *Int16Literal) CopyNumeric() NumericLiteral {
	return &Int16Literal{
		Value:      i.Value,
		isNegative: i.isNegative,
	}
}

type Int8Literal struct {
	Value      string
	isNegative bool
}

func (i *Int8Literal) StringValue() string {
	return i.Value
}

func (i *Int8Literal) IsNegative() bool {
	return i.isNegative
}

func (i *Int8Literal) SetNegative(val bool) {
	i.isNegative = val
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
	if i.isNegative {
		return fmt.Sprintf("-%si8", i.Value)
	}
	return fmt.Sprintf("%si8", i.Value)
}

func (i *Int8Literal) CopyNumeric() NumericLiteral {
	return &Int8Literal{
		Value:      i.Value,
		isNegative: i.isNegative,
	}
}

type UIntLiteral struct {
	Value      string
	isNegative bool
}

func (i *UIntLiteral) StringValue() string {
	return i.Value
}

func (i *UIntLiteral) IsNegative() bool {
	return i.isNegative
}

func (i *UIntLiteral) SetNegative(val bool) {
	i.isNegative = val
}

func NewUIntLiteral(value string) *UIntLiteral {
	return &UIntLiteral{
		Value: value,
	}
}

func (i *UIntLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.UInt64)
}

func (*UIntLiteral) IsLiteral() bool {
	return true
}

func (i *UIntLiteral) inspect() string {
	if i.isNegative {
		return fmt.Sprintf("-%su", i.Value)
	}
	return fmt.Sprintf("%su", i.Value)
}

func (i *UIntLiteral) CopyNumeric() NumericLiteral {
	return &UIntLiteral{
		Value:      i.Value,
		isNegative: i.isNegative,
	}
}

type UInt64Literal struct {
	Value      string
	isNegative bool
}

func (i *UInt64Literal) StringValue() string {
	return i.Value
}

func (i *UInt64Literal) IsNegative() bool {
	return i.isNegative
}

func (i *UInt64Literal) SetNegative(val bool) {
	i.isNegative = val
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
	if i.isNegative {
		return fmt.Sprintf("-%su64", i.Value)
	}
	return fmt.Sprintf("%su64", i.Value)
}

func (i *UInt64Literal) CopyNumeric() NumericLiteral {
	return &UInt64Literal{
		Value:      i.Value,
		isNegative: i.isNegative,
	}
}

type UInt32Literal struct {
	Value      string
	isNegative bool
}

func (i *UInt32Literal) StringValue() string {
	return i.Value
}

func (i *UInt32Literal) IsNegative() bool {
	return i.isNegative
}

func (i *UInt32Literal) SetNegative(val bool) {
	i.isNegative = val
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
	if i.isNegative {
		return fmt.Sprintf("-%su32", i.Value)
	}
	return fmt.Sprintf("%su32", i.Value)
}

func (i *UInt32Literal) CopyNumeric() NumericLiteral {
	return &UInt32Literal{
		Value:      i.Value,
		isNegative: i.isNegative,
	}
}

type UInt16Literal struct {
	Value      string
	isNegative bool
}

func (i *UInt16Literal) StringValue() string {
	return i.Value
}

func (i *UInt16Literal) IsNegative() bool {
	return i.isNegative
}

func (i *UInt16Literal) SetNegative(val bool) {
	i.isNegative = val
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
	if i.isNegative {
		return fmt.Sprintf("-%su16", i.Value)
	}
	return fmt.Sprintf("%su16", i.Value)
}

func (i *UInt16Literal) CopyNumeric() NumericLiteral {
	return &UInt16Literal{
		Value:      i.Value,
		isNegative: i.isNegative,
	}
}

type UInt8Literal struct {
	Value      string
	isNegative bool
}

func (i *UInt8Literal) StringValue() string {
	return i.Value
}

func (i *UInt8Literal) IsNegative() bool {
	return i.isNegative
}

func (i *UInt8Literal) SetNegative(val bool) {
	i.isNegative = val
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
	if i.isNegative {
		return fmt.Sprintf("-%su8", i.Value)
	}
	return fmt.Sprintf("%su8", i.Value)
}

func (i *UInt8Literal) CopyNumeric() NumericLiteral {
	return &UInt8Literal{
		Value:      i.Value,
		isNegative: i.isNegative,
	}
}
