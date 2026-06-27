package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initMethodSignatureDefinitionNode() {
	c := &value.MethodSignatureDefinitionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argName := args[1].MustReference().(ast.IdentifierNode)

			var argTypeParameters []ast.TypeParameterNode
			if !args[2].IsUndefined() {
				argTypeParametersTuple := args[2].AsReference().(value.ArrayTuple)
				argTypeParameters = value.TransformArrayTupleIntoNativeArrayTuple(argTypeParametersTuple, func(v value.Value) ast.TypeParameterNode {
					return v.AsReference().(ast.TypeParameterNode)
				}).ToSlice()
			}

			var argParameters []ast.ParameterNode
			if !args[3].IsUndefined() {
				argParametersTuple := args[3].AsReference().(value.ArrayTuple)
				argParameters = value.TransformArrayTupleIntoNativeArrayTuple(argParametersTuple, func(v value.Value) ast.ParameterNode {
					return v.AsReference().(ast.ParameterNode)
				}).ToSlice()
			}

			var argReturnType ast.TypeNode
			if !args[4].IsUndefined() {
				argReturnType = args[4].MustReference().(ast.TypeNode)
			}

			var argThrowType ast.TypeNode
			if !args[5].IsUndefined() {
				argThrowType = args[5].MustReference().(ast.TypeNode)
			}

			var argDocComment string
			if !args[6].IsUndefined() {
				argDocComment = string(args[6].MustReference().(value.String))
			}

			var argLoc *position.Location
			if args[7].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[7].Pointer())
			}
			self := ast.NewMethodSignatureDefinitionNode(
				argLoc,
				argDocComment,
				argName,
				argTypeParameters,
				argParameters,
				argReturnType,
				argThrowType,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(7),
	)

	vm.Def(
		c,
		"doc_comment",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			result := value.Ref(value.String(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"name",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			result := value.Ref(self.Name)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_parameters",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			entries := value.CastNativeArrayTuplePtr(&self.TypeParameters)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"parameters",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			entries := value.CastNativeArrayTuplePtr(&self.Parameters)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"return_type",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			if self.ReturnType == nil {
				return value.Nil, value.Undefined
			}
			return value.Ref(self.ReturnType), value.Undefined

		},
	)

	vm.Def(
		c,
		"throw_type",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			if self.ThrowType == nil {
				return value.Nil, value.Undefined
			}
			return value.Ref(self.ThrowType), value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
