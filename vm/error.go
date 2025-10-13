package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Error
func initError() {
	// Instance methods
	c := &value.ErrorClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(*value.Object)
			message := args[1]
			self.SetMessage(message)
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"message",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0].AsReference().(*value.Object)
			return self.Message(), value.Undefined
		},
	)

}
