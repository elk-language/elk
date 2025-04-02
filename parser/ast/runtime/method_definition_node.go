package runtime

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initMethodDefinitionNode() {
	c := &value.MethodDefinitionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argName := (string)(args[0].MustReference().(value.String))

			var argTypeParameters []ast.TypeParameterNode
			if !args[1].IsUndefined() {
				argTypeParametersTuple := args[1].MustReference().(*value.ArrayTuple)
				argTypeParameters = make([]ast.TypeParameterNode, argTypeParametersTuple.Length())
				for i, el := range *argTypeParametersTuple {
					argTypeParameters[i] = el.MustReference().(ast.TypeParameterNode)
				}
			}

			var argParameters []ast.ParameterNode
			if !args[2].IsUndefined() {
				argParametersTuple := args[2].MustReference().(*value.ArrayTuple)
				argParameters = make([]ast.ParameterNode, argParametersTuple.Length())
				for i, el := range *argParametersTuple {
					argParameters[i] = el.MustReference().(ast.ParameterNode)
				}
			}

			var argReturnType ast.TypeNode
			if !args[3].IsUndefined() {
				argReturnType = args[3].MustReference().(ast.TypeNode)
			}

			var argThrowType ast.TypeNode
			if !args[4].IsUndefined() {
				argThrowType = args[4].MustReference().(ast.TypeNode)
			}

			var argBody []ast.StatementNode
			if !args[5].IsUndefined() {
				argBodyTuple := args[5].MustReference().(*value.ArrayTuple)
				argBody = make([]ast.StatementNode, argBodyTuple.Length())
				for i, el := range *argBodyTuple {
					argBody[i] = el.MustReference().(ast.StatementNode)
				}
			}

			var argFlags bitfield.BitFlag8
			if !args[6].IsUndefined() {
				argFlags = bitfield.BitFlag8(args[6].AsUInt8())
			}

			var argDocComment string
			if !args[7].IsUndefined() {
				argDocComment = string(args[7].MustReference().(value.String))
			}

			var argLocation *position.Location
			if args[8].IsUndefined() {
				argLocation = position.DefaultLocation
			} else {
				argLocation = (*position.Location)(args[8].Pointer())
			}
			self := ast.NewMethodDefinitionNode(
				argLocation,
				argDocComment,
				argFlags,
				argName,
				argTypeParameters,
				argParameters,
				argReturnType,
				argThrowType,
				argBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(9),
	)

	vm.Def(
		c,
		"doc_comment",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.Ref(value.String(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.Ref(value.String(self.Name))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_parameters",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)

			collection := self.TypeParameters
			arrayTuple := value.NewArrayTupleWithLength(len(collection))
			for i, el := range collection {
				arrayTuple.SetAt(i, value.Ref(el))
			}
			result := value.Ref(arrayTuple)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"parameters",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)

			collection := self.Parameters
			arrayTuple := value.NewArrayTupleWithLength(len(collection))
			for i, el := range collection {
				arrayTuple.SetAt(i, value.Ref(el))
			}
			result := value.Ref(arrayTuple)
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"return_type",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			if self.ReturnType == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.ReturnType), value.Undefined
		},
	)

	vm.Def(
		c,
		"throw_type",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			if self.ThrowType == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.ThrowType), value.Undefined
		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)

			collection := self.Body
			arrayTuple := value.NewArrayTupleWithLength(len(collection))
			for i, el := range collection {
				arrayTuple.SetAt(i, value.Ref(el))
			}
			result := value.Ref(arrayTuple)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"flags",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.UInt8(self.Flags.Byte()).ToValue()
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_abstract",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.ToElkBool(self.IsAbstract())
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_sealed",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.ToElkBool(self.IsSealed())
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_generator",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.ToElkBool(self.IsGenerator())
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_async",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.ToElkBool(self.IsAsync())
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
