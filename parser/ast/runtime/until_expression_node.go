package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initUntilExpressionNode() {
	c := &value.UntilExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argCondition := args[1].MustReference().(ast.ExpressionNode)

			var argThenBody []ast.StatementNode
			if !args[2].IsUndefined() {
				argThenBodyTuple := args[2].AsReference().(value.ArrayTuple)
				argThenBody = value.TransformArrayTupleIntoNativeArrayTuple(argThenBodyTuple, func(v value.Value) ast.StatementNode {
					return v.AsReference().(ast.StatementNode)
				}).ToSlice()
			}

			var argLoc *position.Location
			if args[3].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[3].Pointer())
			}
			self := ast.NewUntilExpressionNode(
				argLoc,
				argCondition,
				argThenBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"condition",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UntilExpressionNode)
			result := value.Ref(self.Condition)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"then_body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UntilExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.ThenBody)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UntilExpressionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UntilExpressionNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UntilExpressionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
