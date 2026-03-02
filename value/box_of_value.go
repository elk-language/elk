package value

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/indent"
)

func BoxOfValueConstructor(class *Class) Value {
	return Ref(NewBoxOfValue(Undefined))
}

// BoxOfValue wraps another value, it's a pointer to another `Value`.
type BoxOfValue Value

func NewBoxOfValue(v Value) *BoxOfValue {
	b := BoxOfValue(v)
	return &b
}

// Retrieve the value stored in the box
func (b *BoxOfValue) Get() Value {
	return Value(*b)
}

func (b *BoxOfValue) GetValue() Value {
	return b.Get()
}

func (b *BoxOfValue) Set(v Value) {
	*b = (BoxOfValue)(v)
}

func (b *BoxOfValue) SetValue(v Value) {
	b.Set(v)
}

// Return the box of the next value in memory
func (b *BoxOfValue) Next(step int) *BoxOfValue {
	ptr := unsafe.Pointer(b)
	return (*BoxOfValue)(unsafe.Add(ptr, step*int(ValueSize)))
}

// Return the box of the previous value in memory
func (b *BoxOfValue) Prev(step int) *BoxOfValue {
	ptr := unsafe.Pointer(b)
	return (*BoxOfValue)(unsafe.Add(ptr, -step*int(ValueSize)))
}

func (*BoxOfValue) Class() *Class {
	return BoxClass
}

func (*BoxOfValue) DirectClass() *Class {
	return BoxClass
}

func (*BoxOfValue) SingletonClass() *Class {
	return nil
}

func (b *BoxOfValue) Copy() Reference {
	return b
}

func (b *BoxOfValue) ToValue() Value {
	return Ref(b)
}

func (*BoxOfValue) InstanceVariables() *InstanceVariables {
	return nil
}

func (b *BoxOfValue) Address() UInt {
	return UInt(uintptr(unsafe.Pointer(b)))
}

func (b *BoxOfValue) ToImmutableBox() *ImmutableBoxOfValue {
	return (*ImmutableBoxOfValue)(unsafe.Pointer(b))
}

func (b *BoxOfValue) ToImmutableBoxInterface() ImmutableBox {
	return b.ToImmutableBox()
}

func (b *BoxOfValue) Inspect() string {
	valInspect := b.Get().Inspect()
	if !strings.ContainsRune(valInspect, '\n') {
		return fmt.Sprintf("Std::Box{&: %p, %s}", b, valInspect)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Box{\n  &: %p", b)

	buff.WriteString(",\n  ")
	indent.IndentStringFromSecondLine(&buff, b.Get().Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (b *BoxOfValue) Error() string {
	return b.Inspect()
}
