package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initRangeLiteralNode() {
	c := &value.RangeLiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argStart := args[1].MustReference().(ast.ExpressionNode)
			argEnd := args[2].MustReference().(ast.ExpressionNode)

			var argOp *token.Token
			if !args[3].IsUndefined() {
				argOp = args[3].MustReference().(*token.Token)
			} else {
				argOp = token.New(position.DefaultLocation, token.CLOSED_RANGE_OP)
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.DefaultLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewRangeLiteralNode(
				argLoc,
				argOp,
				argStart,
				argEnd,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"start",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.RangeLiteralNode)
			result := value.Ref(self.Start)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"end",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.RangeLiteralNode)
			result := value.Ref(self.End)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"op",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.RangeLiteralNode)
			result := value.Ref(self.Op)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.RangeLiteralNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.RangeLiteralNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.RangeLiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)
}
