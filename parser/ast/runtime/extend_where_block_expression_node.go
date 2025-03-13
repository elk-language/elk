package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initExtendWhereBlockExpressionNode() {
	c := &value.ExtendWhereBlockExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			arg0Tuple := args[0].MustReference().(*value.ArrayTuple)
			arg0 := make([]ast.StatementNode, arg0Tuple.Length())
			for i, el := range *arg0Tuple {
				arg0[i] = el.MustReference().(ast.StatementNode)
			}

			arg1Tuple := args[1].MustReference().(*value.ArrayTuple)
			arg1 := make([]ast.TypeParameterNode, arg1Tuple.Length())
			for i, el := range *arg1Tuple {
				arg1[i] = el.MustReference().(ast.TypeParameterNode)
			}

			var argSpan *position.Span
			if args[2].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[2].Pointer())
			}
			self := ast.NewExtendWhereBlockExpressionNode(
				argSpan,
				arg0,
				arg1,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ExtendWhereBlockExpressionNode)

			collection := self.Body
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
		"where",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ExtendWhereBlockExpressionNode)

			collection := self.Where
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
			self := args[0].MustReference().(*ast.ExtendWhereBlockExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
