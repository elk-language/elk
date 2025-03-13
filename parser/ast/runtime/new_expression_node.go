package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initNewExpressionNode() {
	c := &value.NewExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			argPositionalArgumentsTuple := args[0].MustReference().(*value.ArrayTuple)
			argPositionalArguments := make([]ast.ExpressionNode, argPositionalArgumentsTuple.Length())
			for i, el := range *argPositionalArgumentsTuple {
				argPositionalArguments[i] = el.MustReference().(ast.ExpressionNode)
			}

			argNamedArgumentsTuple := args[1].MustReference().(*value.ArrayTuple)
			argNamedArguments := make([]ast.NamedArgumentNode, argNamedArgumentsTuple.Length())
			for i, el := range *argNamedArgumentsTuple {
				argNamedArguments[i] = el.MustReference().(ast.NamedArgumentNode)
			}

			var argSpan *position.Span
			if args[2].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[2].Pointer())
			}
			self := ast.NewNewExpressionNode(
				argSpan,
				argPositionalArguments,
				argNamedArguments,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.NewExpressionNode)

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
			self := args[0].MustReference().(*ast.NewExpressionNode)

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
			self := args[0].MustReference().(*ast.NewExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
