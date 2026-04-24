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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argMethodName := args[1].MustReference().(ast.IdentifierNode)

			var argPositionalArguments []ast.ExpressionNode
			if !args[2].IsUndefined() {
				argPositionalArgumentsTuple := args[2].AsReference().(value.ArrayTuple)
				argPositionalArguments = value.TransformArrayTupleIntoNativeArrayTuple(argPositionalArgumentsTuple, func(v value.Value) ast.ExpressionNode {
					return v.AsReference().(ast.ExpressionNode)
				}).ToSlice()
			}

			var argNamedArguments []ast.NamedArgumentNode
			if !args[3].IsUndefined() {
				argNamedArgumentsTuple := args[3].AsReference().(value.ArrayTuple)
				argNamedArguments = value.TransformArrayTupleIntoNativeArrayTuple(argNamedArgumentsTuple, func(v value.Value) ast.NamedArgumentNode {
					return v.AsReference().(ast.NamedArgumentNode)
				}).ToSlice()
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewReceiverlessMethodCallNode(
				argLoc,
				argMethodName,
				argPositionalArguments,
				argNamedArguments,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"method_name",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)
			result := value.Ref(self.MethodName)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)
			entries := value.CastNativeArrayTuplePtr(&self.PositionalArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"named_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)

			entries := value.CastNativeArrayTuplePtr(&self.NamedArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMethodCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
