package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initYieldExpressionNode() {
	c := &value.YieldExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			var argValue ast.ExpressionNode
			if !args[0].IsUndefined() {
				argValue = args[0].MustReference().(ast.ExpressionNode)
			}

			var argForward bool
			if !args[1].IsUndefined() {
				argForward = value.Truthy(args[1])
			}

			var argSpan *position.Span
			if args[2].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[2].Pointer())
			}
			self := ast.NewYieldExpressionNode(
				argSpan,
				argForward,
				argValue,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"value",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.YieldExpressionNode)
			if self.Value == nil {
				return value.Nil, value.Undefined
			}
			result := value.Ref(self.Value)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"forward",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.YieldExpressionNode)
			result := value.ToElkBool(self.Forward)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.YieldExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.YieldExpressionNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.YieldExpressionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
