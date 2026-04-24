package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initForInExpressionNode() {
	c := &value.ForInExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argPattern := args[1].MustReference().(ast.PatternNode)
			argInExpression := args[2].MustReference().(ast.ExpressionNode)

			argThenTuple := args[3].AsReference().(value.ArrayTuple)
			argThen := value.TransformArrayTupleIntoNativeArrayTuple(argThenTuple, func(v value.Value) ast.StatementNode {
				return v.AsReference().(ast.StatementNode)
			}).ToSlice()

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewForInExpressionNode(
				argLoc,
				argPattern,
				argInExpression,
				argThen,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"pattern_node",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			result := value.Ref(self.Pattern)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"in_expression",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			result := value.Ref(self.InExpression)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"then_body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			entries := value.CastNativeArrayTuplePtr(&self.ThenBody)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
