package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initSelfLiteralNode() {
	c := &value.SelfLiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			var argSpan *position.Span
			if args[0].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[0].Pointer())
			}
			self := ast.NewSelfLiteralNode(
				argSpan,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SelfLiteralNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
