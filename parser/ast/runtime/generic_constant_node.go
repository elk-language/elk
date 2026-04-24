package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initGenericConstantNode() {
	c := &value.GenericConstantNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argConst := args[1].MustReference().(ast.ComplexConstantNode)

			argArgsTuple := args[2].AsReference().(value.ArrayTuple)
			argArgs := value.TransformArrayTupleIntoNativeArrayTuple(argArgsTuple, func(v value.Value) ast.TypeNode {
				return v.AsReference().(ast.TypeNode)
			}).ToSlice()

			var argLoc *position.Location
			if args[3].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[3].Pointer())
			}
			self := ast.NewGenericConstantNode(
				argLoc,
				argConst,
				argArgs,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"constant",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstantNode)
			result := value.Ref(self.Constant)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstantNode)
			entries := value.CastNativeArrayTuplePtr(&self.TypeArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstantNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstantNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstantNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)
}
