package runtime

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initPatternExpressionNode() {
	c := &value.PatternExpressionNodeMixin.MethodContainer
	vm.Def(
		c,
		"to_ast_pattern_expr_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
}
