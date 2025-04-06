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
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argName := (string)(args[1].MustReference().(value.String))

			var argTypeParameters []ast.TypeParameterNode
			if !args[2].IsUndefined() {
				argTypeParametersTuple := args[2].MustReference().(*value.ArrayTuple)
				argTypeParameters = make([]ast.TypeParameterNode, argTypeParametersTuple.Length())
				for i, el := range *argTypeParametersTuple {
					argTypeParameters[i] = el.MustReference().(ast.TypeParameterNode)
				}
			}

			var argParameters []ast.ParameterNode
			if !args[3].IsUndefined() {
				argParametersTuple := args[3].MustReference().(*value.ArrayTuple)
				argParameters = make([]ast.ParameterNode, argParametersTuple.Length())
				for i, el := range *argParametersTuple {
					argParameters[i] = el.MustReference().(ast.ParameterNode)
				}
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

			var argSpan *position.Span
			if args[7].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[7].Pointer())
			}
			self := ast.NewMethodSignatureDefinitionNode(
				argSpan,
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
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			result := value.Ref(value.String(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			result := value.Ref(value.String(self.Name))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_parameters",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)

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
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)

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
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			if self.ThrowType == nil {
				return value.Nil, value.Undefined
			}
			return value.Ref(self.ThrowType), value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodSignatureDefinitionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
