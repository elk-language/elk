package runtime

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initPatternNode() {
	c := &value.PatternNodeMixin.MethodContainer
	vm.Def(
		c,
		"to_ast_pattern_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
}
