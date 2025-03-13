package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initNumericForExpressionNode() {
	c := &value.NumericForExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argInitialiser := args[0].MustReference().(ast.ExpressionNode)
			argCondition := args[1].MustReference().(ast.ExpressionNode)
			argIncrement := args[2].MustReference().(ast.ExpressionNode)

			argThenBodyTuple := args[3].MustReference().(*value.ArrayTuple)
			argThenBody := make([]ast.StatementNode, argThenBodyTuple.Length())
			for i, el := range *argThenBodyTuple {
				argThenBody[i] = el.MustReference().(ast.StatementNode)
			}

			var argSpan *position.Span
			if args[4].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[4].Pointer())
			}
			self := ast.NewNumericForExpressionNode(
				argSpan,
				argInitialiser,
				argCondition,
				argIncrement,
				argThenBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"initialiser",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.NumericForExpressionNode)
			result := value.Ref(self.Initialiser)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"condition",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.NumericForExpressionNode)
			result := value.Ref(self.Condition)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"increment",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.NumericForExpressionNode)
			result := value.Ref(self.Increment)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"then_body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.NumericForExpressionNode)

			collection := self.ThenBody
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
			self := args[0].MustReference().(*ast.NumericForExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
