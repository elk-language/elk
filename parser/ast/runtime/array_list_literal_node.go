package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initArrayListLiteralNode() {
	c := &value.ArrayListLiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			var argElements []ast.ExpressionNode
			if !args[1].IsUndefined() {
				arg0Tuple := args[1].AsReference().(value.ArrayTuple)
				arg0 := value.TransformArrayTupleIntoNativeArrayTuple(arg0Tuple, func(v value.Value) ast.ExpressionNode {
					return v.AsReference().(ast.ExpressionNode)
				})
				argElements = arg0.ToSlice()
			}

			var argCapacity ast.ExpressionNode
			if !args[2].IsUndefined() {
				argCapacity = args[2].MustReference().(ast.ExpressionNode)
			}

			var argLoc *position.Location
			if args[3].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[3].Pointer())
			}
			self := ast.NewArrayListLiteralNode(
				argLoc,
				argElements,
				argCapacity,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"elements",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ArrayListLiteralNode)
			entries := value.CastNativeArrayTuplePtr(&self.Elements)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"capacity",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ArrayListLiteralNode)
			if self.Capacity == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.Capacity), value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ArrayListLiteralNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ArrayListLiteralNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ArrayListLiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
