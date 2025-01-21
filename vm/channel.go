package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// ::Std::Channel
func initChannel() {
	// Instance methods
	c := &value.ChannelClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			nVal := args[1]
			var n int
			if nVal.IsUndefined() {

			} else {
				var ok bool
				n, ok = value.ToGoInt(nVal)
				if !ok {
					return value.Undefined, value.Ref(value.NewError(value.OutOfRangeErrorClass, "channel capacity is too large"))
				}
			}
			self := value.NewChannel(n)
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Channel)(args[0].Pointer())
			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Channel)(args[0].Pointer())
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Channel)(args[0].Pointer())
			return value.SmallInt(self.LeftCapacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.ToElkBool(args[0] == args[1]), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			return value.ToElkBool(args[0] == args[1]), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Channel)(args[0].Pointer())
			err := self.Push(args[1])
			if err.IsUndefined() {
				return value.Ref(self), value.Undefined
			}
			return value.Undefined, err
		},
		DefWithParameters(1),
	)
	Alias(c, "push", "<<")

	Def(
		c,
		"pop",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Channel)(args[0].Pointer())
			result, ok := self.Pop()
			if !ok {
				return value.Undefined, symbol.L_channel_closed.ToValue()
			}
			return result, value.Undefined
		},
	)

	Def(
		c,
		"close",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Channel)(args[0].Pointer())
			err := self.Close()
			if err.IsUndefined() {
				return value.Nil, value.Undefined
			} else {
				return value.Undefined, err
			}
		},
	)

	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Channel)(args[0].Pointer())
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)

}
