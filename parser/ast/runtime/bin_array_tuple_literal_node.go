package ast

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
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			arg0Tuple := args[0].MustReference().(*value.ArrayTuple)
			arg0 := make([]ast.IntCollectionContentNode, arg0Tuple.Length())
			for i, el := range *arg0Tuple {
				arg0[i] = el.MustReference().(ast.IntCollectionContentNode)
			}

			var argSpan *position.Span
			if args[1].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[1].Pointer())
			}
			self := ast.NewBinArrayTupleLiteralNode(
				argSpan,
				arg0,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"elements",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.BinArrayTupleLiteralNode)

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
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.BinArrayTupleLiteralNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
