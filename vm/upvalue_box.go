package vm

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
)

// Represents a pointer to a local variable/value
type UpvalueBox Upvalue

func (l *UpvalueBox) IsClosed() bool {
	return (*Upvalue)(l).IsClosed()
}

func (l *UpvalueBox) IsOpen() bool {
	return !l.IsClosed()
}

func (l *UpvalueBox) Close() {
	(*Upvalue)(l).Close()
}

func (*UpvalueBox) Class() *value.Class {
	return value.UpvalueBoxClass
}

func (*UpvalueBox) DirectClass() *value.Class {
	return value.UpvalueBoxClass
}

func (*UpvalueBox) SingletonClass() *value.Class {
	return nil
}

// Retrieve the value stored in the box
func (l *UpvalueBox) Get() value.Value {
	return *l.slot
}

// Set the value in the box
func (l *UpvalueBox) Set(v value.Value) {
	*l.slot = v
}

func (l *UpvalueBox) ToBox() *value.BoxOfValue {
	return (*value.BoxOfValue)(l.slot)
}

func (l *UpvalueBox) LocalAddress() uintptr {
	return uintptr(unsafe.Pointer(l.slot))
}

func (l *UpvalueBox) ToImmutableBox() *value.ImmutableBoxOfValue {
	return (*value.ImmutableBoxOfValue)(l.slot)
}

// Return the box of the next value in memory
func (l *UpvalueBox) Next(step int) *value.BoxOfValue {
	ptr := unsafe.Pointer(l.slot)
	return (*value.BoxOfValue)(unsafe.Add(ptr, step*int(value.ValueSize)))
}

func (l *UpvalueBox) NextImmutableBox(step int) *value.ImmutableBoxOfValue {
	return (*value.ImmutableBoxOfValue)(l.Next(step))
}

// Return the box of the previous value in memory
func (l *UpvalueBox) Prev(step int) *value.BoxOfValue {
	ptr := unsafe.Pointer(l.slot)
	return (*value.BoxOfValue)(unsafe.Add(ptr, -step*int(value.ValueSize)))
}

func (l *UpvalueBox) PrevImmutableBox(step int) *value.ImmutableBoxOfValue {
	return (*value.ImmutableBoxOfValue)(l.Prev(step))
}

func (l *UpvalueBox) Inspect() string {
	valInspect := l.Get().Inspect()
	if !strings.ContainsRune(valInspect, '\n') {
		return fmt.Sprintf("Std::Box{&: %p, %s}", l, valInspect)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Box{\n  &: %p", l)

	buff.WriteString(",\n  ")
	indent.IndentStringFromSecondLine(&buff, valInspect, 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (l *UpvalueBox) Error() string {
	return l.Inspect()
}

func (*UpvalueBox) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (v *UpvalueBox) Copy() value.Reference {
	return v
}

func (v *UpvalueBox) ToValue() value.Value {
	return value.Ref(v)
}

// Std::UpvalueBox
func initLocalBox() {
	// Instance methods
	c := &value.UpvalueBoxClass.MethodContainer
	Def(
		c,
		"get",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*UpvalueBox)(args[0].Pointer())
			return self.Get(), value.Undefined
		},
	)
	Def(
		c,
		"set",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*UpvalueBox)(args[0].Pointer())
			v := args[1]
			self.Set(v)

			return v, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"address",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*UpvalueBox)(args[0].Pointer())
			return value.UInt(uintptr(unsafe.Pointer(self))).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_immutable_box",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*UpvalueBox)(args[0].Pointer())
			return value.Ref(self.ToImmutableBox()), value.Undefined
		},
	)
	Def(
		c,
		"to_box",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
}
