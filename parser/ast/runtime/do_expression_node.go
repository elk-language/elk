package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initDoExpressionNode() {
	c := value.DoExpressionNodeClass
	vm.Def(
		c,
		"#init",
		func(thread *vm.Thread, args []value.Value) (result value.Value, err value.Value) {
			var argBody []ast.StatementNode
			if !args[1].IsUndefined() {
				argBodyTuple, err := vm.TransformTupleValueIntoNativeArrayTuple(thread, args[1], func(v value.Value) ast.StatementNode {
					return v.AsReference().(ast.StatementNode)
				})
				if err.IsNotUndefined() {
					return value.Undefined, err
				}
				argBody = argBodyTuple.ToSlice()
			}

			var argCatches []*ast.CatchNode
			if !args[2].IsUndefined() {
				argCatchesTuple, err := vm.TransformTupleValueIntoNativeArrayTuple(thread, args[2], func(v value.Value) *ast.CatchNode {
					return v.AsReference().(*ast.CatchNode)
				})
				if err.IsNotUndefined() {
					return value.Undefined, err
				}
				argCatches = argCatchesTuple.ToSlice()
			}

			var argFinally []ast.StatementNode
			if !args[3].IsUndefined() {
				argFinallyTuple, err := vm.TransformTupleValueIntoNativeArrayTuple(thread, args[3], func(v value.Value) ast.StatementNode {
					return v.AsReference().(ast.StatementNode)
				})
				if err.IsNotUndefined() {
					return value.Undefined, err
				}
				argFinally = argFinallyTuple.ToSlice()
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewDoExpressionNode(
				argLoc,
				argBody,
				argCatches,
				argFinally,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.Body)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"catches",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.Catches)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"finally_body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.Finally)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"=~",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)
			other := args[1]
			return value.BoolVal(self.LaxEqual(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
