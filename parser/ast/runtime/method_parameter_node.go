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
			argName := (string)(args[1].MustReference().(value.String))
			argTypeNode := args[2].MustReference().(ast.TypeNode)

			var argInitialiser ast.ExpressionNode
			if !args[3].IsUndefined() {
				argInitialiser = args[3].MustReference().(ast.ExpressionNode)
			}

			var argSetInstanceVariable bool
			if !args[4].IsUndefined() {
				argSetInstanceVariable = value.Truthy(args[4])
			}

			var argKind ast.ParameterKind
			if !args[5].IsUndefined() {
				argKind = ast.ParameterKind(args[5].AsUInt8())
			}

			var argLoc *position.Location
			if args[6].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[6].Pointer())
			}
			self := ast.NewMethodParameterNode(
				argLoc,
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
			if self.TypeNode == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.TypeNode), value.Undefined
		},
	)

	vm.Def(
		c,
		"initialiser",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			if self.Initialiser == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.Initialiser), value.Undefined
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
		"is_normal",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			result := value.ToElkBool(self.Kind == ast.NormalParameterKind)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_positional_rest",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			result := value.ToElkBool(self.Kind == ast.PositionalRestParameterKind)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_named_rest",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			result := value.ToElkBool(self.Kind == ast.NamedRestParameterKind)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodParameterNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
