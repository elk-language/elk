package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initBinArrayTupleLiteralNode() {
	c := &value.BinArrayTupleLiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {

			argElementsTuple := args[1].AsReference().(value.ArrayTuple)
			argElements := value.TransformArrayTupleIntoNativeArrayTuple(argElementsTuple, func(v value.Value) ast.IntCollectionContentNode {
				return v.AsReference().(ast.IntCollectionContentNode)
			}).ToSlice()

			var argLoc *position.Location
			if args[2].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[2].Pointer())
			}
			self := ast.NewBinArrayTupleLiteralNode(
				argLoc,
				argElements,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"elements",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.BinArrayTupleLiteralNode)
			entries := value.CastNativeArrayTuplePtr(&self.Elements)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.BinArrayTupleLiteralNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.BinArrayTupleLiteralNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.BinArrayTupleLiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
