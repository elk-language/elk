package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initWordArrayListLiteralNode() {
	c := &value.WordArrayListLiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			var argElements []ast.WordCollectionContentNode
			if !args[1].IsUndefined() {
				argElementsTuple := args[1].MustReference().(*value.ArrayTuple)
				argElements = make([]ast.WordCollectionContentNode, argElementsTuple.Length())
				for i, el := range *argElementsTuple {
					argElements[i] = el.MustReference().(ast.WordCollectionContentNode)
				}
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
			self := ast.NewWordArrayListLiteralNode(
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
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.WordArrayListLiteralNode)

			collection := self.Elements
			arrayTuple := value.NewArrayTupleWithLength(len(collection))
			for i, el := range collection {
				arrayTuple.SetAt(i, value.Ref(el))
			}
			result := value.Ref(arrayTuple)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"capacity",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.WordArrayListLiteralNode)
			if self.Capacity == nil {
				return value.Nil, value.Undefined
			}
			result := value.Ref(self.Capacity)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.WordArrayListLiteralNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.WordArrayListLiteralNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.WordArrayListLiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
