package value

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/indent"
)

type ImmutableNativeBox[T ValueInterface] struct {
	ptr *T
}

func (b *ImmutableNativeBox[T]) Get() T {
	return *b.ptr
}

func (b *ImmutableNativeBox[T]) GetValue() Value {
	return b.Get().ToValue()
}

func (b *ImmutableNativeBox[T]) Copy() Reference {
	return &ImmutableNativeBox[T]{
		ptr: b.ptr,
	}
}

func (*ImmutableNativeBox[T]) Class() *Class {
	return ImmutableBoxClass
}

func (*ImmutableNativeBox[T]) DirectClass() *Class {
	return BoxClass
}

func (*ImmutableNativeBox[T]) SingletonClass() *Class {
	return nil
}

func (b *ImmutableNativeBox[T]) ToValue() Value {
	return Ref(b)
}

func (b *ImmutableNativeBox[T]) Error() string {
	return b.Inspect()
}

func (b *ImmutableNativeBox[T]) Inspect() string {
	valInspect := b.Get().Inspect()
	if !strings.ContainsRune(valInspect, '\n') {
		return fmt.Sprintf("Std::ImmutableBox{&: %p, %s}", b, valInspect)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::ImmutableBox{\n  &: %p", b)

	buff.WriteString(",\n  ")
	indent.IndentStringFromSecondLine(&buff, b.Get().Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (b *ImmutableNativeBox[T]) Address() UInt {
	return UInt(uintptr(unsafe.Pointer(b.ptr)))
}

func (*ImmutableNativeBox[T]) InstanceVariables() *InstanceVariables {
	return nil
}

func NewImmutableNativeBox[T ValueInterface](v *T) *ImmutableNativeBox[T] {
	return &ImmutableNativeBox[T]{ptr: v}
}
