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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			var argLoc *position.Location
			if args[7].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[7].Pointer())
			}

			argReceiver := args[1].MustReference().(ast.ExpressionNode)
			argName := args[2].MustReference().(ast.IdentifierNode)

			argTypeArgsTuple := args[3].AsReference().(value.ArrayTuple)
			argTypeArgs := value.TransformArrayTupleIntoNativeArrayTuple(argTypeArgsTuple, func(v value.Value) ast.TypeNode {
				return v.AsReference().(ast.TypeNode)
			}).ToSlice()

			var argOp *token.Token
			if args[4].IsUndefined() {
				argOp = token.New(argLoc, token.DOT)
			} else {
				argOp = args[4].MustReference().(*token.Token)
			}

			var argPosArgs []ast.ExpressionNode
			if !args[5].IsUndefined() {
				argPosArgsTuple := args[5].AsReference().(value.ArrayTuple)
				argPosArgs = value.TransformArrayTupleIntoNativeArrayTuple(argPosArgsTuple, func(v value.Value) ast.ExpressionNode {
					return v.AsReference().(ast.ExpressionNode)
				}).ToSlice()
			}

			var argNamedArgs []ast.NamedArgumentNode
			if !args[6].IsUndefined() {
				argNamedArgsTuple := args[6].AsReference().(value.ArrayTuple)
				argNamedArgs = value.TransformArrayTupleIntoNativeArrayTuple(argNamedArgsTuple, func(v value.Value) ast.NamedArgumentNode {
					return v.AsReference().(ast.NamedArgumentNode)
				}).ToSlice()
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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			result := value.Ref(self.Receiver)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"op",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			result := value.Ref(self.Op)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"method_name",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			result := value.Ref(self.MethodName)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			entries := value.CastNativeArrayTuplePtr(&self.TypeArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			entries := value.CastNativeArrayTuplePtr(&self.PositionalArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"named_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			entries := value.CastNativeArrayTuplePtr(&self.NamedArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericMethodCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
