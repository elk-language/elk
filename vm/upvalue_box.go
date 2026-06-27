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

var _ value.Box = &UpvalueBox{}

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
	return value.BoxClass
}

func (*UpvalueBox) DirectClass() *value.Class {
	return value.BoxClass
}

func (*UpvalueBox) SingletonClass() *value.Class {
	return nil
}

// Retrieve the value stored in the box
func (l *UpvalueBox) Get() value.Value {
	return *l.slot
}

func (l *UpvalueBox) GetValue() value.Value {
	return l.Get()
}

// Set the value in the box
func (l *UpvalueBox) Set(v value.Value) {
	*l.slot = v
}

func (l *UpvalueBox) SetValue(v value.Value) {
	l.Set(v)
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

func (l *UpvalueBox) ToImmutableBoxInterface() value.ImmutableBox {
	return l.ToImmutableBox()
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

func (l *UpvalueBox) Address() value.UInt {
	return value.UInt(uintptr(unsafe.Pointer(l)))
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

// Represents a pointer to a local variable/value
type ImmutableUpvalueBox Upvalue

var _ value.ImmutableBox = &ImmutableUpvalueBox{}

func (l *ImmutableUpvalueBox) IsClosed() bool {
	return (*Upvalue)(l).IsClosed()
}

func (l *ImmutableUpvalueBox) IsOpen() bool {
	return !l.IsClosed()
}

func (l *ImmutableUpvalueBox) Close() {
	(*Upvalue)(l).Close()
}

func (*ImmutableUpvalueBox) Class() *value.Class {
	return value.ImmutableBoxClass
}

func (*ImmutableUpvalueBox) DirectClass() *value.Class {
	return value.ImmutableBoxClass
}

func (*ImmutableUpvalueBox) SingletonClass() *value.Class {
	return nil
}

// Retrieve the value stored in the box
func (l *ImmutableUpvalueBox) Get() value.Value {
	return *l.slot
}

func (l *ImmutableUpvalueBox) GetValue() value.Value {
	return l.Get()
}

func (l *ImmutableUpvalueBox) ToImmutableBox() *value.ImmutableBoxOfValue {
	return (*value.ImmutableBoxOfValue)(l.slot)
}

func (l *ImmutableUpvalueBox) LocalAddress() uintptr {
	return uintptr(unsafe.Pointer(l.slot))
}

func (l *ImmutableUpvalueBox) ToImmutableBoxInterface() value.ImmutableBox {
	return l.ToImmutableBox()
}

// Return the box of the next value in memory
func (l *ImmutableUpvalueBox) NextImmutableBox(step int) *value.ImmutableBoxOfValue {
	ptr := unsafe.Pointer(l.slot)
	return (*value.ImmutableBoxOfValue)(unsafe.Add(ptr, step*int(value.ValueSize)))
}

// Return the box of the previous value in memory
func (l *ImmutableUpvalueBox) PrevImmutableBox(step int) *value.ImmutableBoxOfValue {
	ptr := unsafe.Pointer(l.slot)
	return (*value.ImmutableBoxOfValue)(unsafe.Add(ptr, -step*int(value.ValueSize)))
}

func (l *ImmutableUpvalueBox) Address() value.UInt {
	return value.UInt(uintptr(unsafe.Pointer(l)))
}

func (l *ImmutableUpvalueBox) Inspect() string {
	valInspect := l.Get().Inspect()
	if !strings.ContainsRune(valInspect, '\n') {
		return fmt.Sprintf("Std::ImmutableBox{&: %p, %s}", l, valInspect)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::ImmutableBox{\n  &: %p", l)

	buff.WriteString(",\n  ")
	indent.IndentStringFromSecondLine(&buff, valInspect, 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (l *ImmutableUpvalueBox) Error() string {
	return l.Inspect()
}

func (*ImmutableUpvalueBox) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (v *ImmutableUpvalueBox) Copy() value.Reference {
	return v
}

func (v *ImmutableUpvalueBox) ToValue() value.Value {
	return value.Ref(v)
}
