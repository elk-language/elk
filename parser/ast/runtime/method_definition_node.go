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
			argDocComment := (string)(args[0].MustReference().(value.String))
			argName := (string)(args[1].MustReference().(value.String))

			argTypeParametersTuple := args[2].MustReference().(*value.ArrayTuple)
			argTypeParameters := make([]ast.TypeParameterNode, argTypeParametersTuple.Length())
			for i, el := range *argTypeParametersTuple {
				argTypeParameters[i] = el.MustReference().(ast.TypeParameterNode)
			}

			argParametersTuple := args[3].MustReference().(*value.ArrayTuple)
			argParameters := make([]ast.ParameterNode, argParametersTuple.Length())
			for i, el := range *argParametersTuple {
				argParameters[i] = el.MustReference().(ast.ParameterNode)
			}
			argReturnType := args[4].MustReference().(ast.TypeNode)
			argThrowType := args[5].MustReference().(ast.TypeNode)

			argBodyTuple := args[6].MustReference().(*value.ArrayTuple)
			argBody := make([]ast.StatementNode, argBodyTuple.Length())
			for i, el := range *argBodyTuple {
				argBody[i] = el.MustReference().(ast.StatementNode)
			}
			argFlags := bitfield.BitFlag8(args[7].AsUInt8())

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
			result := value.Ref(self.ReturnType)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"throw_type",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.Ref(self.ThrowType)
			return result, value.Undefined

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

}
