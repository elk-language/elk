package value

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/indent"
)

var ImmutableBoxClass *Class // ::Std::ImmutableBox

func initImmutableBox() {
	ImmutableBoxClass = NewClassWithOptions(ClassWithConstructor(ImmutableBoxConstructor))
	StdModule.AddConstantString("ImmutableBox", Ref(ImmutableBoxClass))
}

// Box wraps another value, it's a read only pointer to another `Value`.
type ImmutableBox Value

func ImmutableBoxConstructor(class *Class) Value {
	return Ref(NewImmutableBox(Undefined))
}

func NewImmutableBox(v Value) *ImmutableBox {
	b := ImmutableBox(v)
	return &b
}

func (b *ImmutableBox) ToBox() *Box {
	return (*Box)(b)
}

// Retrieve the value stored in the box
func (b *ImmutableBox) Get() Value {
	return Value(*b)
}

// Return the box of the next value in memory
func (b *ImmutableBox) Next(step int) *ImmutableBox {
	ptr := unsafe.Pointer(b)
	return (*ImmutableBox)(unsafe.Add(ptr, step*int(ValueSize)))
}

// Return the box of the previous value in memory
func (b *ImmutableBox) Prev(step int) *ImmutableBox {
	ptr := unsafe.Pointer(b)
	return (*ImmutableBox)(unsafe.Add(ptr, -step*int(ValueSize)))
}

func (*ImmutableBox) Class() *Class {
	return ImmutableBoxClass
}

func (*ImmutableBox) DirectClass() *Class {
	return ImmutableBoxClass
}

func (*ImmutableBox) SingletonClass() *Class {
	return nil
}

func (b *ImmutableBox) Copy() Reference {
	return b
}

func (*ImmutableBox) InstanceVariables() *InstanceVariables {
	return nil
}

func (b *ImmutableBox) Inspect() string {
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

func (b *ImmutableBox) Error() string {
	return b.Inspect()
}
