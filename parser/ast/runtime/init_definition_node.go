package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initInitDefinitionNode() {
	c := &value.InitDefinitionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			argParametersTuple := args[0].MustReference().(*value.ArrayTuple)
			argParameters := make([]ast.ParameterNode, argParametersTuple.Length())
			for i, el := range *argParametersTuple {
				argParameters[i] = el.MustReference().(ast.ParameterNode)
			}
			argThrowType := args[1].MustReference().(ast.TypeNode)

			argBodyTuple := args[2].MustReference().(*value.ArrayTuple)
			argBody := make([]ast.StatementNode, argBodyTuple.Length())
			for i, el := range *argBodyTuple {
				argBody[i] = el.MustReference().(ast.StatementNode)
			}

			var argLocation *position.Location
			if args[3].IsUndefined() {
				argLocation = position.DefaultLocation
			} else {
				argLocation = (*position.Location)(args[3].Pointer())
			}
			self := ast.NewInitDefinitionNode(
				argLocation,
				argParameters,
				argThrowType,
				argBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"parameters",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)

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
		"throw_type",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)
			result := value.Ref(self.ThrowType)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)

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
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
