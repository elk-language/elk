package types

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/value/symbol"
)

type FloatLiteral struct {
	Value      string
	id         ID
	isNegative bool
}

func (c *FloatLiteral) ToRef() Ref[*FloatLiteral] {
	return Ref[*FloatLiteral](c.id)
}

func (f *FloatLiteral) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("float:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *FloatLiteral) EqualAny(other any) bool {
	o, ok := other.(*FloatLiteral)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *FloatLiteral) ID() ID {
	return f.id
}

func (f *FloatLiteral) SetID(id ID) {
	f.id = id
}

func (f *FloatLiteral) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(f, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(f, parent)
	}
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
	t := &FloatLiteral{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (*FloatLiteral) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Float)
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
		id:         f.id,
	}
}

type Float64Literal struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *Float64Literal) ToRef() Ref[*Float64Literal] {
	return Ref[*Float64Literal](c.id)
}

func (f *Float64Literal) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("float64:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *Float64Literal) EqualAny(other any) bool {
	o, ok := other.(*Float64Literal)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *Float64Literal) ID() ID {
	return f.id
}

func (f *Float64Literal) SetID(id ID) {
	f.id = id
}

func (f *Float64Literal) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(f, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(f, parent)
	}
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
	t := &Float64Literal{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (*Float64Literal) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Float64)
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
		id:         f.id,
	}
}

type Float32Literal struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *Float32Literal) ToRef() Ref[*Float32Literal] {
	return Ref[*Float32Literal](c.id)
}

func (f *Float32Literal) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("float32:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *Float32Literal) EqualAny(other any) bool {
	o, ok := other.(*Float32Literal)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *Float32Literal) ID() ID {
	return f.id
}

func (f *Float32Literal) SetID(id ID) {
	f.id = id
}

func (f *Float32Literal) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(f, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(f, parent)
	}
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
	t := &Float32Literal{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (*Float32Literal) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Float32)
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
		id:         f.id,
	}
}

type BigFloatLiteral struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *BigFloatLiteral) ToRef() Ref[*BigFloatLiteral] {
	return Ref[*BigFloatLiteral](c.id)
}

func (f *BigFloatLiteral) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("bigfloat:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *BigFloatLiteral) EqualAny(other any) bool {
	o, ok := other.(*BigFloatLiteral)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *BigFloatLiteral) ID() ID {
	return f.id
}

func (f *BigFloatLiteral) SetID(id ID) {
	f.id = id
}

func (f *BigFloatLiteral) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(f, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(f, parent)
	}
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
	t := &BigFloatLiteral{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (*BigFloatLiteral) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_BigFloat)
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
		id:         f.id,
	}
}
