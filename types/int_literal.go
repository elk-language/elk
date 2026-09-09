package types

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/value/symbol"
)

type IntLiteral struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *IntLiteral) ToRef() Ref[*IntLiteral] {
	return Ref[*IntLiteral](c.id)
}

func (f *IntLiteral) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("int:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *IntLiteral) EqualAny(other any) bool {
	o, ok := other.(*IntLiteral)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *IntLiteral) ID() ID {
	return f.id
}

func (f *IntLiteral) SetID(id ID) {
	f.id = id
}

func (i *IntLiteral) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
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
	t := &IntLiteral{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (i *IntLiteral) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Int)
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
		id:         i.id,
	}
}

type Int64Literal struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *Int64Literal) ToRef() Ref[*Int64Literal] {
	return Ref[*Int64Literal](c.id)
}

func (f *Int64Literal) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("int64:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *Int64Literal) EqualAny(other any) bool {
	o, ok := other.(*Int64Literal)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *Int64Literal) ID() ID {
	return f.id
}

func (f *Int64Literal) SetID(id ID) {
	f.id = id
}

func (i *Int64Literal) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
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
	t := &Int64Literal{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (i *Int64Literal) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Int64)
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
		id:         i.id,
	}
}

type Int32Literal struct {
	Value      string
	isNegative bool
	id         ID
}

func (i *Int32Literal) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
}

func (c *Int32Literal) ToRef() Ref[*Int32Literal] {
	return Ref[*Int32Literal](c.id)
}

func (f *Int32Literal) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("int32:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *Int32Literal) EqualAny(other any) bool {
	o, ok := other.(*Int32Literal)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *Int32Literal) ID() ID {
	return f.id
}

func (f *Int32Literal) SetID(id ID) {
	f.id = id
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
	t := &Int32Literal{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (i *Int32Literal) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Int32)
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
		id:         i.id,
	}
}

type Int16Literal struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *Int16Literal) ToRef() Ref[*Int16Literal] {
	return Ref[*Int16Literal](c.id)
}

func (f *Int16Literal) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("int16:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *Int16Literal) EqualAny(other any) bool {
	o, ok := other.(*Int16Literal)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *Int16Literal) ID() ID {
	return f.id
}

func (f *Int16Literal) SetID(id ID) {
	f.id = id
}

func (i *Int16Literal) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
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
	t := &Int16Literal{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (i *Int16Literal) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Int16)
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
		id:         i.id,
	}
}

type Int8Literal struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *Int8Literal) ToRef() Ref[*Int8Literal] {
	return Ref[*Int8Literal](c.id)
}

func (f *Int8Literal) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("int8:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *Int8Literal) EqualAny(other any) bool {
	o, ok := other.(*Int8Literal)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *Int8Literal) ID() ID {
	return f.id
}

func (f *Int8Literal) SetID(id ID) {
	f.id = id
}

func (i *Int8Literal) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
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
	t := &Int8Literal{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (i *Int8Literal) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Int8)
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
		id:         i.id,
	}
}

type UIntLiteral struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *UIntLiteral) ToRef() Ref[*UIntLiteral] {
	return Ref[*UIntLiteral](c.id)
}

func (f *UIntLiteral) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("uint:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *UIntLiteral) EqualAny(other any) bool {
	o, ok := other.(*UIntLiteral)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *UIntLiteral) ID() ID {
	return f.id
}

func (f *UIntLiteral) SetID(id ID) {
	f.id = id
}

func (i *UIntLiteral) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
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
	t := &UIntLiteral{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (i *UIntLiteral) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_UInt)
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
		id:         i.id,
	}
}

type UInt64Literal struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *UInt64Literal) ToRef() Ref[*UInt64Literal] {
	return Ref[*UInt64Literal](c.id)
}

func (f *UInt64Literal) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("uint64:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *UInt64Literal) EqualAny(other any) bool {
	o, ok := other.(*UInt64Literal)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *UInt64Literal) ID() ID {
	return f.id
}

func (f *UInt64Literal) SetID(id ID) {
	f.id = id
}

func (i *UInt64Literal) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
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
	t := &UInt64Literal{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (i *UInt64Literal) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_UInt64)
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
		id:         i.id,
	}
}

type UInt32Literal struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *UInt32Literal) ToRef() Ref[*UInt32Literal] {
	return Ref[*UInt32Literal](c.id)
}

func (f *UInt32Literal) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("uint32:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *UInt32Literal) EqualAny(other any) bool {
	o, ok := other.(*UInt32Literal)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *UInt32Literal) ID() ID {
	return f.id
}

func (f *UInt32Literal) SetID(id ID) {
	f.id = id
}

func (i *UInt32Literal) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
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
	t := &UInt32Literal{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (i *UInt32Literal) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_UInt32)
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
		id:         i.id,
	}
}

type UInt16Literal struct {
	Value      string
	isNegative bool
	id         ID
}

func (c *UInt16Literal) ToRef() Ref[*UInt16Literal] {
	return Ref[*UInt16Literal](c.id)
}

func (f *UInt16Literal) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("uint16:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *UInt16Literal) EqualAny(other any) bool {
	o, ok := other.(*UInt16Literal)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *UInt16Literal) ID() ID {
	return f.id
}

func (f *UInt16Literal) SetID(id ID) {
	f.id = id
}

func (i *UInt16Literal) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
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
	t := &UInt16Literal{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (i *UInt16Literal) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_UInt16)
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
	id         ID
}

func (c *UInt8Literal) ToRef() Ref[*UInt8Literal] {
	return Ref[*UInt8Literal](c.id)
}

func (f *UInt8Literal) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("uint8:")
	d.WriteString(f.Value)

	var isNegativeByte byte
	if f.isNegative {
		isNegativeByte = 1
	}
	d.Write([]byte{isNegativeByte})

	return d.Sum64()
}

func (f *UInt8Literal) EqualAny(other any) bool {
	o, ok := other.(*UInt8Literal)
	if !ok {
		return false
	}

	if f.id > 0 {
		return f.id == o.id
	}

	return f.Value == o.Value && f.isNegative == o.isNegative
}

func (f *UInt8Literal) ID() ID {
	return f.id
}

func (f *UInt8Literal) SetID(id ID) {
	f.id = id
}

func (i *UInt8Literal) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
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
	t := &UInt8Literal{
		Value: value,
	}
	Env.RegisterType(t)
	return t
}

func (i *UInt8Literal) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_UInt8)
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
		id:         i.id,
	}
}
