package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::ImmutableBox
func initImmutableBox() {
	// Instance methods
	c := &value.ImmutableBoxClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ImmutableBoxOfValue)(args[0].Pointer())
			v := args[1]
			*self = value.ImmutableBoxOfValue(v)

			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"get",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ImmutableBox)
			return self.GetValue(), value.Undefined
		},
	)
	Def(
		c,
		"address",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ImmutableBox)
			return self.Address().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_immutable_box",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
}
