package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Mixin
func initMixin() {
	// Instance methods
	c := &value.MixinClass.MethodContainer
	Accessor(c, "doc")

	Def(
		c,
		"name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Mixin)
			return value.String(self.Name), nil
		},
	)
}
