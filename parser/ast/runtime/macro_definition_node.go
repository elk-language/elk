package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initMacroDefinitionNode() {
	c := &value.MacroDefinitionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argName := args[1].MustReference().(ast.IdentifierNode)

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

			var argBody []ast.StatementNode
			if !args[4].IsUndefined() {
				argBodyTuple := args[4].MustReference().(*value.ArrayTuple)
				argBody = make([]ast.StatementNode, argBodyTuple.Length())
				for i, el := range *argBodyTuple {
					argBody[i] = el.MustReference().(ast.StatementNode)
				}
			}

			var argSealed bool
			if !args[5].IsUndefined() {
				argSealed = value.Truthy(args[5])
			}

			var argDocComment string
			if !args[6].IsUndefined() {
				argDocComment = string(args[6].MustReference().(value.String))
			}

			var argLocation *position.Location
			if args[7].IsUndefined() {
				argLocation = position.ZeroLocation
			} else {
				argLocation = (*position.Location)(args[7].Pointer())
			}
			self := ast.NewMacroDefinitionNode(
				argLocation,
				argDocComment,
				argSealed,
				argName,
				argParameters,
				argReturnType,
				argBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(7),
	)

	vm.Def(
		c,
		"doc_comment",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MacroDefinitionNode)
			result := value.Ref(value.String(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MacroDefinitionNode)
			result := value.Ref(self.Name)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"parameters",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MacroDefinitionNode)

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
			self := args[0].MustReference().(*ast.MacroDefinitionNode)
			if self.ReturnType == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.ReturnType), value.Undefined
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
		"is_sealed",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodDefinitionNode)
			result := value.ToElkBool(self.IsSealed())
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
