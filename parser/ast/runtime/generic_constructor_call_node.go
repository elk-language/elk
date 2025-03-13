package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initGenericConstructorCallNode() {
	c := &value.GenericConstructorCallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			arg0 := args[0].MustReference().(ast.ComplexConstantNode)

			arg1Tuple := args[1].MustReference().(*value.ArrayTuple)
			arg1 := make([]ast.TypeNode, arg1Tuple.Length())
			for i, el := range *arg1Tuple {
				arg1[i] = el.MustReference().(ast.TypeNode)
			}

			arg2Tuple := args[2].MustReference().(*value.ArrayTuple)
			arg2 := make([]ast.ExpressionNode, arg2Tuple.Length())
			for i, el := range *arg2Tuple {
				arg2[i] = el.MustReference().(ast.ExpressionNode)
			}

			arg3Tuple := args[3].MustReference().(*value.ArrayTuple)
			arg3 := make([]ast.NamedArgumentNode, arg3Tuple.Length())
			for i, el := range *arg3Tuple {
				arg3[i] = el.MustReference().(ast.NamedArgumentNode)
			}

			var argSpan *position.Span
			if args[4].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[4].Pointer())
			}
			self := ast.NewGenericConstructorCallNode(
				argSpan,
				arg0,
				arg1,
				arg2,
				arg3,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"class_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)
			result := value.Ref(self.ClassNode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)

			collection := self.TypeArguments
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
		"positional_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)

			collection := self.PositionalArguments
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
		"named_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)

			collection := self.NamedArguments
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
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
