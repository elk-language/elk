package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Interface
func initInterface() {
	// Instance methods
	c := &value.InterfaceClass.MethodContainer

	Def(
		c,
		"name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Interface)
			return value.String(self.Name), nil
		},
	)
}
