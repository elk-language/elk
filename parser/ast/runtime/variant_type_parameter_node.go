package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initVariantTypeParameterNode() {
	c := &value.VariantTypeParameterNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argVariance := args[0].AsUInt8()
			argName := (string)(args[1].MustReference().(value.String))
			argLowerBound := args[2].MustReference().(ast.TypeNode)
			argUpperBound := args[3].MustReference().(ast.TypeNode)
			argDefault := args[4].MustReference().(ast.TypeNode)

			var argSpan *position.Span
			if args[5].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[5].Pointer())
			}
			self := ast.NewVariantTypeParameterNode(
				argSpan,
				ast.Variance(argVariance),
				argName,
				argLowerBound,
				argUpperBound,
				argDefault,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(6),
	)

	vm.Def(
		c,
		"variance",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			result := value.UInt8(self.Variance).ToValue()
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			result := value.Ref(value.String(self.Name))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"lower_bound",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			result := value.Ref(self.LowerBound)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"upper_bound",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			result := value.Ref(self.UpperBound)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_default",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			result := value.Ref(self.Default)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
