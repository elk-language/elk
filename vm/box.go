package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Box
func initBox() {
	// Instance methods
	c := &value.BoxClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BoxOfValue)(args[0].Pointer())
			v := args[1]
			self.Set(v)

			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"get",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Box)
			return self.GetValue(), value.Undefined
		},
	)
	Def(
		c,
		"set",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Box)
			v := args[1]
			self.SetValue(v)

			return v, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"address",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Box)
			return self.Address().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_immutable_box",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.Box)
			return self.ToImmutableBoxInterface().ToValue(), value.Undefined
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
