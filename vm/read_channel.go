package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::ReadChannel
func initReadChannel() {
	// Singleton methods
	c := &value.ReadChannelClass.SingletonClass().MethodContainer
	Def(
		c,
		"closed",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			ch := value.NewChannelOfValue(0)
			ch.Close()
			return ch.ToReadChannelOfValue().ToValue(), value.Undefined
		},
	)

	// Instance methods
	c = &value.ReadChannelClass.MethodContainer
	Def(
		c,
		"capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ReadChannel)
			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ReadChannel)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ReadChannel)
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
		"pop",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ReadChannel)
			result, err := self.PopCtx(vm.Aborter.Context())
			if err.IsNotUndefined() {
				return value.Undefined, err
			}
			return result, value.Undefined
		},
	)
	Alias(c, "<<@", "pop")

	Def(
		c,
		"next",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ReadChannel)
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
