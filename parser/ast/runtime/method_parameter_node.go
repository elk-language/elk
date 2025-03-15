package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initMethodParameterNode() {
	c := &value.MethodParameterNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argName := (string)(args[0].MustReference().(value.String))
			argTypeNode := args[1].MustReference().(ast.TypeNode)
			argInitialiser := args[2].MustReference().(ast.ExpressionNode)
			argSetInstanceVariable := value.Truthy(args[3])
			argKind := ast.ParameterKind(args[4].AsUInt8())

			var argSpan *position.Span
			if args[5].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[5].Pointer())
			}
			self := ast.NewMethodParameterNode(
				argSpan,
				argName,
				argSetInstanceVariable,
				argTypeNode,
				argInitialiser,
				argKind,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(6),
	)

	vm.Def(
		c,
		"name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			result := value.Ref(value.String(self.Name))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			result := value.Ref(self.TypeNode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"initialiser",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			result := value.Ref(self.Initialiser)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"set_instance_variable",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			result := value.ToElkBool(self.SetInstanceVariable)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"kind",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			result := value.UInt8(self.Kind).ToValue()
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_optional",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			return value.ToElkBool(self.IsOptional()), value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
