package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::WriteChannel
func initWriteChannel() {
	// Singleton methods
	c := &value.WriteChannelClass.SingletonClass().MethodContainer
	Def(
		c,
		"closed",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			ch := value.NewChannelOfValue(0)
			ch.Close()
			return ch.ToWriteChannelOfValue().ToValue(), value.Undefined
		},
	)

	// Instance methods
	c = &value.WriteChannelClass.MethodContainer
	Def(
		c,
		"capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.WriteChannel)
			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.WriteChannel)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.WriteChannel)
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
			self := args[0].AsReference().(value.WriteChannel)
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
		"close",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.WriteChannel)
			err := self.Close()
			if err.IsUndefined() {
				return value.Nil, value.Undefined
			} else {
				return value.Undefined, err
			}
		},
	)

}
