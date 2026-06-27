package runtime

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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {

			argBodyTuple := args[1].AsReference().(value.ArrayTuple)
			argBody := value.TransformArrayTupleIntoNativeArrayTuple(argBodyTuple, func(v value.Value) ast.StatementNode {
				return v.AsReference().(ast.StatementNode)
			}).ToSlice()

			whereTuple := args[2].AsReference().(value.ArrayTuple)
			argWhere := value.TransformArrayTupleIntoNativeArrayTuple(whereTuple, func(v value.Value) ast.TypeParameterNode {
				return v.AsReference().(ast.TypeParameterNode)
			}).ToSlice()

			var argLoc *position.Location
			if args[3].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[3].Pointer())
			}
			self := ast.NewExtendWhereBlockExpressionNode(
				argLoc,
				argBody,
				argWhere,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ExtendWhereBlockExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.Body)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"where_params",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ExtendWhereBlockExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.Where)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ExtendWhereBlockExpressionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ExtendWhereBlockExpressionNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ExtendWhereBlockExpressionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
