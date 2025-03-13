package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initExpressionStatementNode() {
	c := &value.ExpressionStatementNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			arg0 := args[0].MustReference().(ast.ExpressionNode)

			var argSpan *position.Span
			if args[1].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[1].Pointer())
			}
			self := ast.NewExpressionStatementNode(
				argSpan,
				arg0,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"expression",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ExpressionStatementNode)
			result := value.Ref(self.Expression)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ExpressionStatementNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
