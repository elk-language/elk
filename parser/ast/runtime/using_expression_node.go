package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initUsingExpressionNode() {
	c := &value.UsingExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			argEntriesTuple := args[0].MustReference().(*value.ArrayTuple)
			argEntries := make([]ast.UsingEntryNode, argEntriesTuple.Length())
			for i, el := range *argEntriesTuple {
				argEntries[i] = el.MustReference().(ast.UsingEntryNode)
			}

			var argSpan *position.Span
			if args[1].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[1].Pointer())
			}
			self := ast.NewUsingExpressionNode(
				argSpan,
				argEntries,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"entries",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UsingExpressionNode)

			collection := self.Entries
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
			self := args[0].MustReference().(*ast.UsingExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
