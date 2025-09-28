package vm

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
)

const (
	UpvalueLongIndexFlag bitfield.BitFlag8 = 1 << iota
	UpvalueLocalFlag
)

// Represents a captured variable from an outer context
type Upvalue struct {
	// Points to the region where the closed variable lives.
	// In an open upvalue slot points to a slot on the value stack.
	// In a closed upvalue slot points to the `closed` field
	// within the upvalue.
	slot *value.Value
	// Undefined in open upvalues, contains the variable's value in closed upvalues.
	closed value.Value
	// Points to the next upvalue on the stack creating a linked list
	next *Upvalue
}

func (u *Upvalue) IsClosed() bool {
	return u.slot == &u.closed
}

func (u *Upvalue) IsOpen() bool {
	return !u.IsClosed()
}

func (u *Upvalue) Close() {
	if u.IsClosed() {
		return
	}

	u.unsafeClose()
}

func (u *Upvalue) unsafeClose() {
	// move the variable from the stack to the heap
	// inside of the upvalue
	u.closed = *u.slot
	// the location pointer now points to the `closed` field
	// within the upvalue
	u.slot = &u.closed
}

func NewUpvalue(slot *value.Value) *Upvalue {
	return &Upvalue{
		slot: slot,
	}
}

// Implementation of the value.Value interface

func (*Upvalue) Class() *value.Class {
	return nil
}

func (*Upvalue) DirectClass() *value.Class {
	return nil
}

func (*Upvalue) SingletonClass() *value.Class {
	return nil
}

func (*Upvalue) Inspect() string {
	return "upvalue"
}

func (u *Upvalue) Error() string {
	return u.Inspect()
}

func (*Upvalue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (v *Upvalue) Copy() value.Reference {
	return v
}
