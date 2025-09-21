package vm

import (
	"unsafe"

	"github.com/elk-language/elk/value"
)

// Std::ImmutableBox
func initImmutableBox() {
	// Class methods
	c := &value.ImmutableBoxClass.SingletonClass().MethodContainer
	Def(
		c,
		"at",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			address := args[1].AsUInt()
			b := (*value.ImmutableBox)(unsafe.Pointer(uintptr(address)))

			return value.Ref(b), value.Undefined
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.ImmutableBoxClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.ImmutableBox)(args[0].Pointer())
			v := args[1]
			*self = value.ImmutableBox(v)

			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"get",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.ImmutableBox)(args[0].Pointer())
			return self.Get(), value.Undefined
		},
	)
	Def(
		c,
		"address",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.ImmutableBox)(args[0].Pointer())
			return value.UInt(uintptr(unsafe.Pointer(self))).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_immutable_box",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)

	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.ImmutableBox)(args[0].Pointer())

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

	Def(
		c,
		"prev",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.ImmutableBox)(args[0].Pointer())

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
}
