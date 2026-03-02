package value

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/indent"
)

// Box wraps another value, it's a read only pointer to another `Value`.
type ImmutableBoxOfValue Value

func ImmutableBoxOfValueConstructor(class *Class) Value {
	return Ref(NewImmutableBoxOfValue(Undefined))
}

func NewImmutableBoxOfValue(v Value) *ImmutableBoxOfValue {
	b := ImmutableBoxOfValue(v)
	return &b
}

func (b *ImmutableBoxOfValue) ToBox() *BoxOfValue {
	return (*BoxOfValue)(b)
}

// Retrieve the value stored in the box
func (b *ImmutableBoxOfValue) Get() Value {
	return Value(*b)
}

// Retrieve the value stored in the box
func (b *ImmutableBoxOfValue) GetValue() Value {
	return b.Get()
}

// Return the box of the next value in memory
func (b *ImmutableBoxOfValue) Next(step int) *ImmutableBoxOfValue {
	ptr := unsafe.Pointer(b)
	return (*ImmutableBoxOfValue)(unsafe.Add(ptr, step*int(ValueSize)))
}

// Return the box of the previous value in memory
func (b *ImmutableBoxOfValue) Prev(step int) *ImmutableBoxOfValue {
	ptr := unsafe.Pointer(b)
	return (*ImmutableBoxOfValue)(unsafe.Add(ptr, -step*int(ValueSize)))
}

func (*ImmutableBoxOfValue) Class() *Class {
	return ImmutableBoxClass
}

func (*ImmutableBoxOfValue) DirectClass() *Class {
	return ImmutableBoxClass
}

func (*ImmutableBoxOfValue) SingletonClass() *Class {
	return nil
}

func (b *ImmutableBoxOfValue) Copy() Reference {
	return b
}

func (b *ImmutableBoxOfValue) ToValue() Value {
	return Ref(b)
}

func (*ImmutableBoxOfValue) InstanceVariables() *InstanceVariables {
	return nil
}

func (b *ImmutableBoxOfValue) Inspect() string {
	valInspect := b.Get().Inspect()
	if !strings.ContainsRune(valInspect, '\n') {
		return fmt.Sprintf("Std::ImmutableBox{&: %p, %s}", b, valInspect)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::ImmutableBox{\n  &: %p", b)

	buff.WriteString(",\n  ")
	indent.IndentStringFromSecondLine(&buff, valInspect, 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (b *ImmutableBoxOfValue) Address() UInt {
	return UInt(uintptr(unsafe.Pointer(b)))
}

func (b *ImmutableBoxOfValue) Error() string {
	return b.Inspect()
}
