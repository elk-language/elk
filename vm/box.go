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
}
