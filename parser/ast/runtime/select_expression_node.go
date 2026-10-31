package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initSelectExpressionNode() {
	c := value.SelectExpressionNodeClass
	vm.Def(
		c,
		"#init",
		func(thread *vm.Thread, args []value.Value) (result value.Value, err value.Value) {
			argCasesTuple, err := vm.TransformTupleValueIntoNativeArrayTuple(thread, args[1], func(v value.Value) *ast.SelectCaseNode {
				return v.AsReference().(*ast.SelectCaseNode)
			})
			if err.IsNotUndefined() {
				return value.Undefined, err
			}
			argCases := argCasesTuple.ToSlice()

			var argElseBody []ast.StatementNode
			if !args[2].IsUndefined() {
				argElseBodyTuple, err := vm.TransformTupleValueIntoNativeArrayTuple(thread, args[2], func(v value.Value) ast.StatementNode {
					return v.AsReference().(ast.StatementNode)
				})
				if err.IsNotUndefined() {
					return value.Undefined, err
				}
				argElseBody = argElseBodyTuple.ToSlice()
			}

			var argLoc *position.Location
			if args[3].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[3].Pointer())
			}
			self := ast.NewSelectExpressionNode(
				argLoc,
				argCases,
				argElseBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"cases",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.Cases)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"else_body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.ElseBody)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectExpressionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectExpressionNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"=~",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectExpressionNode)
			other := args[1]
			return value.BoolVal(self.LaxEqual(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectExpressionNode)
			result := value.Ref(value.String(self.String()))
			return result, value.Undefined
		},
	)

}
