package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initHexArrayListLiteralNode() {
	c := &value.HexArrayListLiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			var argElements []ast.IntCollectionContentNode
			if !args[1].IsUndefined() {
				argElementsTuple := args[1].AsReference().(value.ArrayTuple)
				argElements = value.TransformArrayTupleIntoNativeArrayTuple(argElementsTuple, func(v value.Value) ast.IntCollectionContentNode {
					return v.AsReference().(ast.IntCollectionContentNode)
				}).ToSlice()
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
			self := ast.NewHexArrayListLiteralNode(
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
			self := args[0].MustReference().(*ast.HexArrayListLiteralNode)
			entries := value.CastNativeArrayTuplePtr(&self.Elements)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"capacity",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.HexArrayListLiteralNode)
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
			self := args[0].MustReference().(*ast.HexArrayListLiteralNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.HexArrayListLiteralNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.HexArrayListLiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
