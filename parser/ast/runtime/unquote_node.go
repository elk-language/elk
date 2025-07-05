package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initUnquoteNode() {
	c := &value.UnquoteNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argExprNode := args[1].MustReference().(ast.ExpressionNode)

			var argKind ast.UnquoteKind
			if !args[2].IsUndefined() {
				argKind = ast.UnquoteKind(args[2].AsUInt8())
			}

			var argLoc *position.Location
			if args[3].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[3].Pointer())
			}

			self := ast.NewUnquoteNode(
				argLoc,
				argKind,
				argExprNode,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"expression",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnquoteNode)
			result := value.Ref(self.Expression)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"kind",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnquoteNode)
			result := value.UInt8(self.Kind)
			return result.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnquoteNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnquoteNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UnquoteNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

	vm.Def(
		c,
		"value",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			return value.Ref(value.String("")), value.Undefined
		},
	)

}
