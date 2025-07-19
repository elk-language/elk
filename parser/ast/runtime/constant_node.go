package runtime

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initConstantNode() {
	c := &value.ConstantNodeMixin.MethodContainer
	vm.Def(
		c,
		"to_ast_const_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
}
