package value

import (
	"fmt"
	"unsafe"
	"weak"
)

var WeakClass *Class // ::Std::Weak

type Weak weak.Pointer[Box]

type internalWeak struct {
	ptr unsafe.Pointer
}

// Make a weak pointer
func MakeWeak(box *Box) Weak {
	ptr := weak.Make(box)
	return Weak(ptr)
}

// Convert to a box (a strong pointer)
func (w Weak) ToBox() *Box {
	ptr := weak.Pointer[Box](w)
	return ptr.Value()
}

// Convert to an immutable box (a strong pointer)
func (w Weak) ToImmutableBox() *ImmutableBox {
	return w.ToBox().ToImmutableBox()
}

func (w Weak) ToBoxValue() Value {
	box := w.ToBox()
	if box == nil {
		return Nil
	}

	return Ref(box)
}

func (w Weak) ToImmutableBoxValue() Value {
	box := w.ToImmutableBox()
	if box == nil {
		return Nil
	}

	return Ref(box)
}

func (w Weak) internal() internalWeak {
	return *(*internalWeak)(unsafe.Pointer(&w))
}

func (w Weak) ToValue() Value {
	internal := w.internal()
	return Value{
		flag: WEAK_FLAG,
		ptr:  internal.ptr,
	}
}

func (Weak) Class() *Class {
	return WeakClass
}

func (Weak) DirectClass() *Class {
	return WeakClass
}

func (Weak) SingletonClass() *Class {
	return nil
}

func (Weak) InstanceVariables() *InstanceVariables {
	return nil
}

func (w Weak) Inspect() string {
	return fmt.Sprintf("Std::Weak{&: %p}", w.internal().ptr)
}

func (w Weak) Error() string {
	return w.Inspect()
}

func initWeak() {
	WeakClass = NewClass()
	StdModule.AddConstantString("Weak", Ref(WeakClass))
}
