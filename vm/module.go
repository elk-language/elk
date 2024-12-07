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
			self := args[0].MustReference().(*value.Module)
			return value.Ref(value.String(self.Name)), value.Undefined
		},
	)
}
