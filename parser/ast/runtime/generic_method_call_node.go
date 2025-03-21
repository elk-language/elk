package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initGenericMethodCallNode() {
	c := &value.GenericMethodCallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			var argSpan *position.Span
			if args[6].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[6].Pointer())
			}

			argReceiver := args[0].MustReference().(ast.ExpressionNode)
			argName := (string)(args[1].MustReference().(value.String))

			argTypeArgsTuple := args[2].MustReference().(*value.ArrayTuple)
			argTypeArgs := make([]ast.TypeNode, argTypeArgsTuple.Length())
			for i, el := range *argTypeArgsTuple {
				argTypeArgs[i] = el.MustReference().(ast.TypeNode)
			}

			var argOp *token.Token
			if args[3].IsUndefined() {
				argOp = token.New(argSpan, token.DOT)
			} else {
				argOp = args[3].MustReference().(*token.Token)
			}

			var argPosArgs []ast.ExpressionNode
			if !args[4].IsUndefined() {
				argPosArgsTuple := args[4].MustReference().(*value.ArrayTuple)
				argPosArgs = make([]ast.ExpressionNode, argPosArgsTuple.Length())
				for i, el := range *argPosArgsTuple {
					argPosArgs[i] = el.MustReference().(ast.ExpressionNode)
				}
			}

			var argNamedArgs []ast.NamedArgumentNode
			if !args[5].IsUndefined() {
				argNamedArgsTuple := args[5].MustReference().(*value.ArrayTuple)
				argNamedArgs = make([]ast.NamedArgumentNode, argNamedArgsTuple.Length())
				for i, el := range *argNamedArgsTuple {
					argNamedArgs[i] = el.MustReference().(ast.NamedArgumentNode)
				}
			}

			self := ast.NewGenericMethodCallNode(
				argSpan,
				argReceiver,
				argOp,
				argName,
				argTypeArgs,
				argPosArgs,
				argNamedArgs,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(8),
	)

	vm.Def(
		c,
		"receiver",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			result := value.Ref(self.Receiver)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"op",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			result := value.Ref(self.Op)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"method_name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			result := value.Ref(value.String(self.MethodName))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)

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
			self := args[0].MustReference().(*ast.GenericMethodCallNode)

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
			self := args[0].MustReference().(*ast.GenericMethodCallNode)

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
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
