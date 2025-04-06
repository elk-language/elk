package runtime

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
			var argInitialiser ast.ExpressionNode
			if !args[1].IsUndefined() {
				argInitialiser = args[1].MustReference().(ast.ExpressionNode)
			}

			var argCondition ast.ExpressionNode
			if !args[2].IsUndefined() {
				argCondition = args[2].MustReference().(ast.ExpressionNode)
			}

			var argIncrement ast.ExpressionNode
			if !args[3].IsUndefined() {
				argIncrement = args[3].MustReference().(ast.ExpressionNode)
			}

			var argThenBody []ast.StatementNode
			if !args[4].IsUndefined() {
				argThenBodyTuple := args[4].MustReference().(*value.ArrayTuple)
				argThenBody = make([]ast.StatementNode, argThenBodyTuple.Length())
				for i, el := range *argThenBodyTuple {
					argThenBody[i] = el.MustReference().(ast.StatementNode)
				}
			}

			var argSpan *position.Span
			if args[5].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[5].Pointer())
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
			if self.Initialiser == nil {
				return value.Nil, value.Undefined
			}
			result := value.Ref(self.Initialiser)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"condition",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.NumericForExpressionNode)
			if self.Condition == nil {
				return value.Nil, value.Undefined
			}
			result := value.Ref(self.Condition)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"increment",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.NumericForExpressionNode)
			if self.Increment == nil {
				return value.Nil, value.Undefined
			}
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

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.NumericForExpressionNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.NumericForExpressionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
