package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initAttributeParameterNode() {
	c := &value.AttributeParameterNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			arg0 := (string)(args[0].MustReference().(value.String))
			arg1 := args[1].MustReference().(ast.TypeNode)
			arg2 := args[2].MustReference().(ast.ExpressionNode)

			var argSpan *position.Span
			if args[3].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[3].Pointer())
			}
			self := ast.NewAttributeParameterNode(
				argSpan,
				arg0,
				arg1,
				arg2,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttributeParameterNode)
			result := value.Ref(value.String(self.Name))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttributeParameterNode)
			result := value.Ref(self.TypeNode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_optional",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttributeParameterNode)
			result := value.ToElkBool(self.IsOptional())
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"initialiser",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttributeParameterNode)
			result := value.Ref(self.Initialiser)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttributeParameterNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
