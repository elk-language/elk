package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initSwitchExpressionNode() {
	c := &value.SwitchExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argValue := args[1].MustReference().(ast.ExpressionNode)

			argCasesTuple := args[2].AsReference().(value.ArrayTuple)
			argCases := value.TransformArrayTupleIntoNativeArrayTuple(argCasesTuple, func(v value.Value) *ast.SwitchCaseNode {
				return v.AsReference().(*ast.SwitchCaseNode)
			}).ToSlice()

			var argElseBody []ast.StatementNode
			if !args[3].IsUndefined() {
				argElseBodyTuple := args[3].AsReference().(value.ArrayTuple)
				argElseBody = value.TransformArrayTupleIntoNativeArrayTuple(argElseBodyTuple, func(v value.Value) ast.StatementNode {
					return v.AsReference().(ast.StatementNode)
				}).ToSlice()
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewSwitchExpressionNode(
				argLoc,
				argValue,
				argCases,
				argElseBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"value",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)
			result := value.Ref(self.Value)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"cases",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.Cases)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"else_body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.ElseBody)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)
			result := value.Ref(value.String(self.String()))
			return result, value.Undefined
		},
	)

}
