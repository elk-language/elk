package vm

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
)

// Represents a pointer to a local variable/value
type LocalBox Upvalue

func (l *LocalBox) IsClosed() bool {
	return (*Upvalue)(l).IsClosed()
}

func (l *LocalBox) IsOpen() bool {
	return !l.IsClosed()
}

func (l *LocalBox) Close() {
	(*Upvalue)(l).Close()
}

func (*LocalBox) Class() *value.Class {
	return value.LocalBoxClass
}

func (*LocalBox) DirectClass() *value.Class {
	return value.LocalBoxClass
}

func (*LocalBox) SingletonClass() *value.Class {
	return nil
}

// Retrieve the value stored in the box
func (l *LocalBox) Get() value.Value {
	return *l.slot
}

// Set the value in the box
func (l *LocalBox) Set(v value.Value) {
	*l.slot = v
}

func (l *LocalBox) ToBox() *value.Box {
	return (*value.Box)(l.slot)
}

func (l *LocalBox) LocalAddress() uintptr {
	return uintptr(unsafe.Pointer(l.slot))
}

func (l *LocalBox) ToImmutableBox() *value.ImmutableBox {
	return (*value.ImmutableBox)(l.slot)
}

// Return the box of the next value in memory
func (l *LocalBox) Next(step int) *value.Box {
	ptr := unsafe.Pointer(l.slot)
	return (*value.Box)(unsafe.Add(ptr, step*int(value.ValueSize)))
}

func (l *LocalBox) NextImmutableBox(step int) *value.ImmutableBox {
	return (*value.ImmutableBox)(l.Next(step))
}

// Return the box of the previous value in memory
func (l *LocalBox) Prev(step int) *value.Box {
	ptr := unsafe.Pointer(l.slot)
	return (*value.Box)(unsafe.Add(ptr, -step*int(value.ValueSize)))
}

func (l *LocalBox) PrevImmutableBox(step int) *value.ImmutableBox {
	return (*value.ImmutableBox)(l.Prev(step))
}

func (l *LocalBox) Inspect() string {
	valInspect := l.Get().Inspect()
	if !strings.ContainsRune(valInspect, '\n') {
		return fmt.Sprintf("Std::LocalBox{&: %p, %s}", l, valInspect)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::LocalBox{\n  &: %p", l)

	buff.WriteString(",\n  ")
	indent.IndentStringFromSecondLine(&buff, valInspect, 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (l *LocalBox) Error() string {
	return l.Inspect()
}

func (*LocalBox) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (v *LocalBox) Copy() value.Reference {
	return v
}

// Std::LocalBox
func initLocalBox() {
	// Class methods
	c := &value.BoxClass.SingletonClass().MethodContainer
	Def(
		c,
		"at",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			address := args[1].AsUInt()
			b := (*LocalBox)(unsafe.Pointer(uintptr(address)))

			return value.Ref(b), value.Undefined
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.LocalBoxClass.MethodContainer
	Def(
		c,
		"get",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())
			return self.Get(), value.Undefined
		},
	)
	Def(
		c,
		"set",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())
			v := args[1]
			self.Set(v)

			return v, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"address",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())
			return value.UInt(uintptr(unsafe.Pointer(self))).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"close",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())
			self.Close()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"is_closed",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())
			return value.ToElkBool(self.IsClosed()), value.Undefined
		},
	)
	Def(
		c,
		"is_open",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())
			return value.ToElkBool(self.IsOpen()), value.Undefined
		},
	)
	Def(
		c,
		"local_address",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())
			return value.UInt(self.LocalAddress()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_immutable_box",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())
			return value.Ref(self.ToImmutableBox()), value.Undefined
		},
	)
	Def(
		c,
		"to_box",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())
			return value.Ref(self.ToBox()), value.Undefined
		},
	)

	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())

			var step int
			if args[1].IsUndefined() {
				step = 1
			} else {
				step = args[1].AsInt()
			}

			return value.Ref(self.Next(step)), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "next_box", "next")

	Def(
		c,
		"next_immutable_box",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())

			var step int
			if args[1].IsUndefined() {
				step = 1
			} else {
				step = args[1].AsInt()
			}

			return value.Ref(self.NextImmutableBox(step)), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"prev",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())

			var step int
			if args[1].IsUndefined() {
				step = 1
			} else {
				step = args[1].AsInt()
			}

			return value.Ref(self.Prev(step)), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "prev_box", "prev")

	Def(
		c,
		"prev_immutable_box",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*LocalBox)(args[0].Pointer())

			var step int
			if args[1].IsUndefined() {
				step = 1
			} else {
				step = args[1].AsInt()
			}

			return value.Ref(self.PrevImmutableBox(step)), value.Undefined
		},
		DefWithParameters(1),
	)
}
