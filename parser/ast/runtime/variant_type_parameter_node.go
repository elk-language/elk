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
			argName := (string)(args[1].MustReference().(value.String))

			var argLowerBound ast.TypeNode
			if !args[2].IsUndefined() {
				argLowerBound = args[2].MustReference().(ast.TypeNode)
			}

			var argUpperBound ast.TypeNode
			if !args[3].IsUndefined() {
				argUpperBound = args[3].MustReference().(ast.TypeNode)
			}

			var argDefault ast.TypeNode
			if !args[4].IsUndefined() {
				argDefault = args[4].MustReference().(ast.TypeNode)
			}

			var argVariance ast.Variance
			if !args[5].IsUndefined() {
				argVariance = ast.Variance(args[5].AsUInt8())
			}

			var argSpan *position.Span
			if args[6].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[6].Pointer())
			}
			self := ast.NewVariantTypeParameterNode(
				argSpan,
				argVariance,
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
		"default_node",
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

	vm.Def(
		c,
		"is_invariant",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			result := value.ToElkBool(self.Variance == ast.INVARIANT)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_covariant",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			result := value.ToElkBool(self.Variance == ast.COVARIANT)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_contravariant",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			result := value.ToElkBool(self.Variance == ast.CONTRAVARIANT)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariantTypeParameterNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)
}
