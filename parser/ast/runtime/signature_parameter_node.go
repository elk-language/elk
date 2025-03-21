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
			argName := (string)(args[0].MustReference().(value.String))

			var argTypeNode ast.TypeNode
			if !args[1].IsUndefined() {
				argTypeNode = args[1].MustReference().(ast.TypeNode)
			}

			var argOptional bool
			if !args[2].IsUndefined() {
				argOptional = value.Truthy(args[2])
			}

			var argKind ast.ParameterKind
			if !args[3].IsUndefined() {
				argKind = ast.ParameterKind(args[3].AsUInt8())
			}

			var argSpan *position.Span
			if args[4].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[4].Pointer())
			}
			self := ast.NewSignatureParameterNode(
				argSpan,
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
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SignatureParameterNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
