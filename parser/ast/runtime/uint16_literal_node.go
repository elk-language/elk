package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initUInt16LiteralNode() {
	c := &value.UInt16LiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argValue := (string)(args[1].MustReference().(value.String))

			var argLoc *position.Location
			if args[2].IsUndefined() {
				argLoc = position.DefaultLocation
			} else {
				argLoc = (*position.Location)(args[2].Pointer())
			}
			self := ast.NewUInt16LiteralNode(
				argLoc,
				argValue,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"value",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UInt16LiteralNode)
			result := value.Ref(value.String(self.Value))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UInt16LiteralNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UInt16LiteralNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UInt16LiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
