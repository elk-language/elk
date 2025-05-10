package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initSignatureParameterNode() {
	c := &value.SignatureParameterNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argName := (string)(args[1].MustReference().(value.String))

			var argTypeNode ast.TypeNode
			if !args[2].IsUndefined() {
				argTypeNode = args[2].MustReference().(ast.TypeNode)
			}

			var argOptional bool
			if !args[3].IsUndefined() {
				argOptional = value.Truthy(args[3])
			}

			var argKind ast.ParameterKind
			if !args[4].IsUndefined() {
				argKind = ast.ParameterKind(args[4].AsUInt8())
			}

			var argLoc *position.Location
			if args[5].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[5].Pointer())
			}
			self := ast.NewSignatureParameterNode(
				argLoc,
				argName,
				argTypeNode,
				argOptional,
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
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			result := value.Ref(value.String(self.Name))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			if self.TypeNode == nil {
				return value.Nil, value.Undefined
			}
			result := value.Ref(self.TypeNode)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_optional",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			result := value.ToElkBool(self.Optional)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_normal",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			result := value.ToElkBool(self.Kind == ast.NormalParameterKind)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_positional_rest",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			result := value.ToElkBool(self.Kind == ast.PositionalRestParameterKind)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_named_rest",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			result := value.ToElkBool(self.Kind == ast.NamedRestParameterKind)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"kind",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			result := value.UInt8(self.Kind)
			return result.ToValue(), value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			result := value.Ref(value.String(self.String()))
			return result, value.Undefined
		},
	)

}
