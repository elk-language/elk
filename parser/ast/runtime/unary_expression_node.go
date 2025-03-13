package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initUnaryExpressionNode() {
	c := &value.UnaryExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argOp := args[0].MustReference().(*token.Token)
			argRight := args[1].MustReference().(ast.ExpressionNode)

			var argSpan *position.Span
			if args[2].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[2].Pointer())
			}
			self := ast.NewUnaryExpressionNode(
				argSpan,
				argOp,
				argRight,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"op",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnaryExpressionNode)
			result := value.Ref(self.Op)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"right",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnaryExpressionNode)
			result := value.Ref(self.Right)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnaryExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
