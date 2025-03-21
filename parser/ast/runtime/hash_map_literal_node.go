package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initHashMapLiteralNode() {
	c := &value.HashMapLiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			var argElements []ast.ExpressionNode
			if !args[0].IsUndefined() {
				argElementsTuple := args[0].MustReference().(*value.ArrayTuple)
				argElements = make([]ast.ExpressionNode, argElementsTuple.Length())
				for i, el := range *argElementsTuple {
					argElements[i] = el.MustReference().(ast.ExpressionNode)
				}
			}

			var argCapacity ast.ExpressionNode
			if !args[1].IsUndefined() {
				argCapacity = args[1].MustReference().(ast.ExpressionNode)
			}

			var argSpan *position.Span
			if args[2].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[2].Pointer())
			}
			self := ast.NewHashMapLiteralNode(
				argSpan,
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
			self := args[0].MustReference().(*ast.HashMapLiteralNode)

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
			self := args[0].MustReference().(*ast.HashMapLiteralNode)
			if self.Capacity == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.Capacity), value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.HashMapLiteralNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
