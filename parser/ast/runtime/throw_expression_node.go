package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initThrowExpressionNode() {
	c := &value.ThrowExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			var argValue ast.ExpressionNode
			if !args[0].IsUndefined() {
				argValue = args[0].MustReference().(ast.ExpressionNode)
			}

			var argUnchecked bool
			if !args[1].IsUndefined() {
				argUnchecked = value.Truthy(args[1])
			}

			var argSpan *position.Span
			if args[2].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[2].Pointer())
			}

			self := ast.NewThrowExpressionNode(
				argSpan,
				argUnchecked,
				argValue,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"is_unchecked",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ThrowExpressionNode)
			result := value.ToElkBool(self.Unchecked)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"value",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ThrowExpressionNode)
			if self.Value == nil {
				return value.Nil, value.Undefined
			}
			result := value.Ref(self.Value)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ThrowExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ThrowExpressionNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ThrowExpressionNode)
			result := value.Ref(value.String(self.String()))
			return result, value.Undefined
		},
	)

}
