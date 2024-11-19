package vm

import (
	"github.com/elk-language/elk/value"
)

func initClass() {
	// Instance methods
	c := &value.ClassClass.MethodContainer
	Accessor(c, "doc")

	Def(
		c,
		"superclass",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Class)
			superclass := self.Superclass()
			if superclass == nil {
				return value.Nil, nil
			}
			return self.Superclass(), nil
		},
	)
	Def(
		c,
		"name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Class)
			return value.String(self.Name), nil
		},
	)
}
