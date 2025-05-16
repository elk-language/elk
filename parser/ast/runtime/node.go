package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
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
	vm.Def(
		c,
		"traverse",
		func(v *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(ast.Node)
			fn := args[1]

			switch f := fn.SafeAsReference().(type) {
			case *vm.Closure:
				for node := range ast.Traverse(self) {
					ok, err := v.CallClosure(f, value.Ref(node))
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					if value.Falsy(ok) {
						return value.False, value.Undefined
					}
				}
			default:
				for node := range ast.Traverse(self) {
					ok, err := v.CallMethodByName(symbol.L_call, fn, value.Ref(node))
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					if value.Falsy(ok) {
						return value.False, value.Undefined
					}
				}
			}

			return value.True, value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"iter",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(ast.Node)
			iterator := ast.NewNodeIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
}
