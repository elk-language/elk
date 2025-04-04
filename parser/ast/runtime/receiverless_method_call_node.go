package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initReceiverlessMethodCallNode() {
	c := &value.ReceiverlessMethodCallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argMethodName := (string)(args[0].MustReference().(value.String))

			var argPositionalArguments []ast.ExpressionNode
			if !args[1].IsUndefined() {
				argPositionalArgumentsTuple := args[1].MustReference().(*value.ArrayTuple)
				argPositionalArguments = make([]ast.ExpressionNode, argPositionalArgumentsTuple.Length())
				for i, el := range *argPositionalArgumentsTuple {
					argPositionalArguments[i] = el.MustReference().(ast.ExpressionNode)
				}
			}

			var argNamedArguments []ast.NamedArgumentNode
			if !args[2].IsUndefined() {
				argNamedArgumentsTuple := args[2].MustReference().(*value.ArrayTuple)
				argNamedArguments = make([]ast.NamedArgumentNode, argNamedArgumentsTuple.Length())
				for i, el := range *argNamedArgumentsTuple {
					argNamedArguments[i] = el.MustReference().(ast.NamedArgumentNode)
				}
			}

			var argSpan *position.Span
			if args[3].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[3].Pointer())
			}
			self := ast.NewReceiverlessMethodCallNode(
				argSpan,
				argMethodName,
				argPositionalArguments,
				argNamedArguments,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"method_name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)
			result := value.Ref(value.String(self.MethodName))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)

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
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)

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
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
