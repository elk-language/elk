package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Module
func initModule() {
	// Instance methods
	c := &value.ModuleClass.MethodContainer
	Accessor(c, "doc")

	Def(
		c,
		"name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Module)
			return value.String(self.Name), nil
		},
	)
}
