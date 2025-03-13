package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initClosureLiteralNode() {
	c := &value.ClosureLiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			arg0Tuple := args[0].MustReference().(*value.ArrayTuple)
			arg0 := make([]ast.ParameterNode, arg0Tuple.Length())
			for i, el := range *arg0Tuple {
				arg0[i] = el.MustReference().(ast.ParameterNode)
			}
			arg1 := args[1].MustReference().(ast.TypeNode)
			arg2 := args[2].MustReference().(ast.TypeNode)

			arg3Tuple := args[3].MustReference().(*value.ArrayTuple)
			arg3 := make([]ast.StatementNode, arg3Tuple.Length())
			for i, el := range *arg3Tuple {
				arg3[i] = el.MustReference().(ast.StatementNode)
			}

			var argSpan *position.Span
			if args[4].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[4].Pointer())
			}
			self := ast.NewClosureLiteralNode(
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
		"parameters",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)

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
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			result := value.Ref(self.ReturnType)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"throw_type",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			result := value.Ref(self.ThrowType)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)

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
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
