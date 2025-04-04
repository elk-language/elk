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
			argStart := args[0].MustReference().(ast.ExpressionNode)
			argEnd := args[1].MustReference().(ast.ExpressionNode)

			var argOp *token.Token
			if !args[2].IsUndefined() {
				argOp = args[2].MustReference().(*token.Token)
			} else {
				argOp = token.New(position.DefaultSpan, token.CLOSED_RANGE_OP)
			}

			var argSpan *position.Span
			if args[3].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[3].Pointer())
			}
			self := ast.NewRangeLiteralNode(
				argSpan,
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
