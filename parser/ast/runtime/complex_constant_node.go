package runtime

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initComplexConstantNode() {
	c := &value.ComplexConstantNodeMixin.MethodContainer
	vm.Def(
		c,
		"to_ast_complex_const_node",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
}
