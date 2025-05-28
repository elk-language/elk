package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initFormalParameterNode() {
	c := &value.FormalParameterNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argName := args[1].MustReference().(ast.IdentifierNode)
			argKind := ast.ParameterKind(args[2].AsUInt8())

			var argType ast.TypeNode
			if !args[3].IsUndefined() {
				argType = args[3].MustReference().(ast.TypeNode)
			}

			var argInit ast.ExpressionNode
			if !args[4].IsUndefined() {
				argInit = args[4].MustReference().(ast.ExpressionNode)
			}

			var argLoc *position.Location
			if args[5].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[5].Pointer())
			}
			self := ast.NewFormalParameterNode(
				argLoc,
				argName,
				argType,
				argInit,
				argKind,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			result := value.Ref(self.Name)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"type_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			result := value.Ref(self.TypeNode)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"initialiser",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			result := value.Ref(self.Initialiser)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"kind",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			result := value.UInt8(self.Kind).ToValue()
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_optional",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			result := value.ToElkBool(self.IsOptional())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_normal",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			result := value.ToElkBool(self.Kind == ast.NormalParameterKind)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_positional_rest",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			result := value.ToElkBool(self.Kind == ast.PositionalRestParameterKind)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_named_rest",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			result := value.ToElkBool(self.Kind == ast.NamedRestParameterKind)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.FormalParameterNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)
}
