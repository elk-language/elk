package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initBreakExpressionNode() {
	c := &value.BreakExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			var argLabel string
			if !args[0].IsUndefined() {
				argLabel = (string)(args[0].MustReference().(value.String))
			}

			var argValue ast.ExpressionNode
			if !args[1].IsUndefined() {
				argValue = args[1].MustReference().(ast.ExpressionNode)
			}

			var argSpan *position.Span
			if args[2].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[2].Pointer())
			}
			self := ast.NewBreakExpressionNode(
				argSpan,
				argLabel,
				argValue,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"label",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.BreakExpressionNode)
			result := value.Ref(value.String(self.Label))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"value",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.BreakExpressionNode)
			if self.Value == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.Value), value.Undefined
		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.BreakExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
