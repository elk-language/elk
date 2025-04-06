package runtime

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initNode() {
	c := &value.NodeMixin.MethodContainer
	vm.Def(
		c,
		"to_ast",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
}
