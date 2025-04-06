package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initModifierForInNode() {
	c := &value.ModifierForInNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argThenExpression := args[1].MustReference().(ast.ExpressionNode)
			argPattern := args[2].MustReference().(ast.PatternNode)
			argInExpression := args[3].MustReference().(ast.ExpressionNode)

			var argSpan *position.Span
			if args[4].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[4].Pointer())
			}
			self := ast.NewModifierForInNode(
				argSpan,
				argThenExpression,
				argPattern,
				argInExpression,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"then_expression",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModifierForInNode)
			result := value.Ref(self.ThenExpression)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"pattern",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModifierForInNode)
			result := value.Ref(self.Pattern)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"in_expression",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModifierForInNode)
			result := value.Ref(self.InExpression)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModifierForInNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModifierForInNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModifierForInNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
