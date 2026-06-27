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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argName := args[1].MustReference().(ast.IdentifierNode)

			argTypeArgsTuple := args[2].AsReference().(value.ArrayTuple)
			argTypeArgs := value.TransformArrayTupleIntoNativeArrayTuple(argTypeArgsTuple, func(v value.Value) ast.TypeNode {
				return v.AsReference().(ast.TypeNode)
			}).ToSlice()

			var argPosArgs []ast.ExpressionNode
			if !args[3].IsUndefined() {
				argPosArgsTuple := args[3].AsReference().(value.ArrayTuple)
				argPosArgs = value.TransformArrayTupleIntoNativeArrayTuple(argPosArgsTuple, func(v value.Value) ast.ExpressionNode {
					return v.AsReference().(ast.ExpressionNode)
				}).ToSlice()
			}

			var argNamedArgs []ast.NamedArgumentNode
			if !args[4].IsUndefined() {
				argNamedArgsTuple := args[4].AsReference().(value.ArrayTuple)
				argNamedArgs = value.TransformArrayTupleIntoNativeArrayTuple(argNamedArgsTuple, func(v value.Value) ast.NamedArgumentNode {
					return v.AsReference().(ast.NamedArgumentNode)
				}).ToSlice()
			}

			var argLoc *position.Location
			if args[5].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[5].Pointer())
			}
			self := ast.NewGenericReceiverlessMethodCallNode(
				argLoc,
				argName,
				argTypeArgs,
				argPosArgs,
				argNamedArgs,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"method_name",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)
			result := value.Ref(self.MethodName)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)
			entries := value.CastNativeArrayTuplePtr(&self.TypeArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)
			entries := value.CastNativeArrayTuplePtr(&self.PositionalArguments)
			return entries.ToValue(), value.Undefined

		},
	)

	vm.Def(
		c,
		"named_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)
			entries := value.CastNativeArrayTuplePtr(&self.NamedArguments)
			return entries.ToValue(), value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericReceiverlessMethodCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)
}
