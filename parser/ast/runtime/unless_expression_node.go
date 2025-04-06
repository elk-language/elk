package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initUnlessExpressionNode() {
	c := &value.UnlessExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argCondition := args[1].MustReference().(ast.ExpressionNode)

			var argThenBody []ast.StatementNode
			if !args[2].IsUndefined() {
				argThenBodyTuple := args[2].MustReference().(*value.ArrayTuple)
				argThenBody = make([]ast.StatementNode, argThenBodyTuple.Length())
				for i, el := range *argThenBodyTuple {
					argThenBody[i] = el.MustReference().(ast.StatementNode)
				}
			}

			var argElseBody []ast.StatementNode
			if !args[3].IsUndefined() {
				argElseBodyTuple := args[3].MustReference().(*value.ArrayTuple)
				argElseBody = make([]ast.StatementNode, argElseBodyTuple.Length())
				for i, el := range *argElseBodyTuple {
					argElseBody[i] = el.MustReference().(ast.StatementNode)
				}
			}

			var argSpan *position.Span
			if args[4].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[4].Pointer())
			}
			self := ast.NewUnlessExpressionNode(
				argSpan,
				argCondition,
				argThenBody,
				argElseBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"condition",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnlessExpressionNode)
			result := value.Ref(self.Condition)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"then_body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnlessExpressionNode)

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
		"else_body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnlessExpressionNode)

			collection := self.ElseBody
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
			self := args[0].MustReference().(*ast.UnlessExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnlessExpressionNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnlessExpressionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
