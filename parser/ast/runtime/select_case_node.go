package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initSelectCaseNode() {
	c := &value.SelectCaseNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argExpr := args[1].MustReference().(ast.ExpressionNode)

			argBodyTuple := args[2].AsReference().(value.ArrayTuple)
			argBody := value.TransformArrayTupleIntoNativeArrayTuple(argBodyTuple, func(v value.Value) ast.StatementNode {
				return v.AsReference().(ast.StatementNode)
			}).ToSlice()

			var argOk ast.IdentifierNode
			if args[3].IsNotUndefined() {
				argOk = args[3].AsReference().(ast.IdentifierNode)
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewSelectCaseNode(
				argLoc,
				argExpr,
				argOk,
				argBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"expression_node",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectCaseNode)
			result := value.Ref(self.Expression)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"ok_var",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectCaseNode)
			if self.OkVar == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.OkVar), value.Undefined
		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectCaseNode)
			entries := value.CastNativeArrayTuplePtr(&self.Body)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectCaseNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectCaseNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelectCaseNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
