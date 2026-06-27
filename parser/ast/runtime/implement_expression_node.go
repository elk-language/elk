package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initImplementExpressionNode() {
	c := &value.ImplementExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argConstantsTuple := args[1].AsReference().(value.ArrayTuple)
			argConstants := value.TransformArrayTupleIntoNativeArrayTuple(argConstantsTuple, func(v value.Value) ast.ComplexConstantNode {
				return v.AsReference().(ast.ComplexConstantNode)
			}).ToSlice()

			var argLoc *position.Location
			if args[2].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[2].Pointer())
			}
			self := ast.NewImplementExpressionNode(
				argLoc,
				argConstants,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"constants",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ImplementExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.Constants)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ImplementExpressionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ImplementExpressionNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ImplementExpressionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
