package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initKernel() {
	// Std::Kernel
	c := &value.KernelModule.SingletonClass().MethodContainer
	vm.Def(
		c,
		"#splice",
		func(vm *vm.Thread, args []value.Value) (value.Value, value.Value) {
			baseNode := args[1].AsReference().(ast.Node)

			var replacementNodes *[]ast.Node
			if !args[2].IsUndefined() {
				tuple := args[2].AsReference().(value.ArrayTuple)
				replacementNodes = (*[]ast.Node)(
					value.TransformArrayTupleIntoNativeArrayTuple(tuple, func(v value.Value) ast.Node {
						return v.AsReference().(ast.Node)
					}),
				)
			}

			result := ast.Splice(baseNode, nil, replacementNodes)
			return value.Ref(result), value.Undefined
		},
		vm.DefWithParameters(2),
	)
}
