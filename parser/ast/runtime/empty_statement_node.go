package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initEmptyStatementNode() {
	c := &value.EmptyStatementNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			var argLoc *position.Location
			if args[1].IsUndefined() {
				argLoc = position.DefaultLocation
			} else {
				argLoc = (*position.Location)(args[1].Pointer())
			}
			self := ast.NewEmptyStatementNode(
				argLoc,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.EmptyStatementNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.EmptyStatementNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.EmptyStatementNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
