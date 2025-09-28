package vm

import (
	"unsafe"

	"github.com/elk-language/elk/value"
)

// Std::Box
func initBox() {
	// Class methods
	c := &value.BoxClass.SingletonClass().MethodContainer
	Def(
		c,
		"at",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			address := args[1].AsUInt()
			b := (*value.Box)(unsafe.Pointer(uintptr(address)))

			return value.Ref(b), value.Undefined
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.BoxClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Box)(args[0].Pointer())
			v := args[1]
			self.Set(v)

			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"get",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Box)(args[0].Pointer())
			return self.Get(), value.Undefined
		},
	)
	Def(
		c,
		"set",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Box)(args[0].Pointer())
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
			self := (*value.Box)(args[0].Pointer())
			return value.UInt(uintptr(unsafe.Pointer(self))).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_immutable_box",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Box)(args[0].Pointer())
			return value.Ref(self.ToImmutableBox()), value.Undefined
		},
	)
	Def(
		c,
		"to_box",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)

	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Box)(args[0].Pointer())

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
			self := (*value.Box)(args[0].Pointer())

			var step int
			if args[1].IsUndefined() {
				step = 1
			} else {
				step = args[1].AsInt()
			}

			next := self.Next(step)
			return value.Ref(next.ToImmutableBox()), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"prev",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Box)(args[0].Pointer())

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
			self := (*value.Box)(args[0].Pointer())

			var step int
			if args[1].IsUndefined() {
				step = 1
			} else {
				step = args[1].AsInt()
			}

			prev := self.Prev(step)
			return value.Ref(prev.ToImmutableBox()), value.Undefined
		},
		DefWithParameters(1),
	)
}
