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
			self := args[0].MustReference().(*value.Class)
			superclass := self.Superclass()
			if superclass == nil {
				return value.Nil, value.Nil
			}
			return value.Ref(superclass), value.Nil
		},
	)
	Def(
		c,
		"name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Class)
			return value.Ref(value.String(self.Name)), value.Nil
		},
	)
}
