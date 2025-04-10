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
			var argLoc *position.Location
			if args[7].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[7].Pointer())
			}

			argReceiver := args[1].MustReference().(ast.ExpressionNode)
			argName := (string)(args[2].MustReference().(value.String))

			argTypeArgsTuple := args[3].MustReference().(*value.ArrayTuple)
			argTypeArgs := make([]ast.TypeNode, argTypeArgsTuple.Length())
			for i, el := range *argTypeArgsTuple {
				argTypeArgs[i] = el.MustReference().(ast.TypeNode)
			}

			var argOp *token.Token
			if args[4].IsUndefined() {
				argOp = token.New(argLoc, token.DOT)
			} else {
				argOp = args[4].MustReference().(*token.Token)
			}

			var argPosArgs []ast.ExpressionNode
			if !args[5].IsUndefined() {
				argPosArgsTuple := args[5].MustReference().(*value.ArrayTuple)
				argPosArgs = make([]ast.ExpressionNode, argPosArgsTuple.Length())
				for i, el := range *argPosArgsTuple {
					argPosArgs[i] = el.MustReference().(ast.ExpressionNode)
				}
			}

			var argNamedArgs []ast.NamedArgumentNode
			if !args[6].IsUndefined() {
				argNamedArgsTuple := args[6].MustReference().(*value.ArrayTuple)
				argNamedArgs = make([]ast.NamedArgumentNode, argNamedArgsTuple.Length())
				for i, el := range *argNamedArgsTuple {
					argNamedArgs[i] = el.MustReference().(ast.NamedArgumentNode)
				}
			}

			self := ast.NewGenericMethodCallNode(
				argLoc,
				argReceiver,
				argOp,
				argName,
				argTypeArgs,
				argPosArgs,
				argNamedArgs,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(7),
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

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
