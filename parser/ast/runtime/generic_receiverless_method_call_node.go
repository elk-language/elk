package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initGenericReceiverlessMethodCallNode() {
	c := &value.GenericReceiverlessMethodCallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argName := (string)(args[0].MustReference().(value.String))

			argTypeArgsTuple := args[1].MustReference().(*value.ArrayTuple)
			argTypeArgs := make([]ast.TypeNode, argTypeArgsTuple.Length())
			for i, el := range *argTypeArgsTuple {
				argTypeArgs[i] = el.MustReference().(ast.TypeNode)
			}

			var argPosArgs []ast.ExpressionNode
			if !args[2].IsUndefined() {
				argPosArgsTuple := args[2].MustReference().(*value.ArrayTuple)
				argPosArgs = make([]ast.ExpressionNode, argPosArgsTuple.Length())
				for i, el := range *argPosArgsTuple {
					argPosArgs[i] = el.MustReference().(ast.ExpressionNode)
				}
			}

			var argNamedArgs []ast.NamedArgumentNode
			if !args[3].IsUndefined() {
				argNamedArgsTuple := args[3].MustReference().(*value.ArrayTuple)
				argNamedArgs = make([]ast.NamedArgumentNode, argNamedArgsTuple.Length())
				for i, el := range *argNamedArgsTuple {
					argNamedArgs[i] = el.MustReference().(ast.NamedArgumentNode)
				}
			}

			var argSpan *position.Span
			if args[4].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[4].Pointer())
			}
			self := ast.NewGenericReceiverlessMethodCallNode(
				argSpan,
				argName,
				argTypeArgs,
				argPosArgs,
				argNamedArgs,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(6),
	)

	vm.Def(
		c,
		"method_name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)
			result := value.Ref(value.String(self.MethodName))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)

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
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)

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
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)

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
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
