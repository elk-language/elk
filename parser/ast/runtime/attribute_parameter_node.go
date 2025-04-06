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
			argName := (string)(args[1].MustReference().(value.String))

			var argTypeNode ast.TypeNode
			if !args[2].IsUndefined() {
				argTypeNode = args[2].MustReference().(ast.TypeNode)
			}

			var argInit ast.ExpressionNode
			if !args[3].IsUndefined() {
				argInit = args[3].MustReference().(ast.ExpressionNode)
			}

			var argSpan *position.Span
			if args[4].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[4].Pointer())
			}
			self := ast.NewAttributeParameterNode(
				argSpan,
				argName,
				argTypeNode,
				argInit,
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
			if self.TypeNode == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.TypeNode), value.Undefined

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
		"is_normal",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			return value.True, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_positional_rest",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			return value.False, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_named_rest",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			return value.False, value.Undefined
		},
	)

	vm.Def(
		c,
		"initialiser",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttributeParameterNode)
			if self.Initialiser == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.Initialiser), value.Undefined

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
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttributeParameterNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttributeParameterNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
