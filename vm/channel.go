package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// ::Std::Channel
func initChannel() {
	// Singleton methods
	c := &value.ChannelClass.SingletonClass().MethodContainer
	Def(
		c,
		"closed",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			ch := value.NewChannelOfValue(0)
			ch.Close()
			return value.Ref(ch), value.Undefined
		},
	)

	// Instance methods
	c = &value.ChannelClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
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
			self := value.NewChannelOfValue(n)
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Channel)
			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Channel)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Channel)
			return value.SmallInt(self.LeftCapacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.BoolVal(args[0] == args[1]), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return value.BoolVal(args[0] == args[1]), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Channel)
			err := self.PushCtx(vm.Aborter.Context(), args[1])
			if err.IsUndefined() {
				return self.ToValue(), value.Undefined
			}
			return value.Undefined, err
		},
		DefWithParameters(1),
	)
	Alias(c, "push", "<<")

	Def(
		c,
		"pop",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Channel)
			result, ok, err := self.PopCtx(vm.Aborter.Context())
			if err.IsNotUndefined() {
				return value.Undefined, err
			}
			if !ok {
				return value.Undefined, symbol.L_channel_closed.ToValue()
			}
			return result, value.Undefined
		},
	)
	Alias(c, "<<@", "pop")

	Def(
		c,
		"close",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Channel)
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
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Channel)
			return self.NextValueCtx(vm.Aborter.Context())
		},
	)
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)

}
