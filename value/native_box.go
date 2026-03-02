package value

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/indent"
)

type NativeBox[T ValueInterface] struct {
	ptr *T
}

func (b *NativeBox[T]) Get() T {
	return *b.ptr
}

func (b *NativeBox[T]) Set(v T) {
	*b.ptr = v
}

func (b *NativeBox[T]) GetValue() Value {
	return b.Get().ToValue()
}

func (b *NativeBox[T]) SetValue(v Value) {
	*b.ptr = v.ToInterface().(T)
}

func (b *NativeBox[T]) Copy() Reference {
	return &NativeBox[T]{
		ptr: b.ptr,
	}
}

func (*NativeBox[T]) Class() *Class {
	return BoxClass
}

func (*NativeBox[T]) DirectClass() *Class {
	return BoxClass
}

func (*NativeBox[T]) SingletonClass() *Class {
	return nil
}

func (b *NativeBox[T]) ToValue() Value {
	return Ref(b)
}

func (b *NativeBox[T]) Error() string {
	return b.Inspect()
}

func (b *NativeBox[T]) Address() UInt {
	return UInt(uintptr(unsafe.Pointer(b.ptr)))
}

func (b *NativeBox[T]) ToImmutableBox() *ImmutableNativeBox[T] {
	return NewImmutableNativeBox(b.ptr)
}

func (b *NativeBox[T]) ToImmutableBoxInterface() ImmutableBox {
	return b.ToImmutableBox()
}

func (b *NativeBox[T]) Inspect() string {
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

func (*NativeBox[T]) InstanceVariables() *InstanceVariables {
	return nil
}

func NewNativeBox[T ValueInterface](v *T) *NativeBox[T] {
	return &NativeBox[T]{ptr: v}
}
