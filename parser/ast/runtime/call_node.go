package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initCallNode() {
	c := &value.CallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argReceiver := args[1].MustReference().(ast.ExpressionNode)
			argNilSafe := value.Truthy(args[2])

			argPosArgsTuple := args[3].AsReference().(value.ArrayTuple)
			argPosArgs := value.TransformArrayTupleIntoNativeArrayTuple(argPosArgsTuple, func(v value.Value) ast.ExpressionNode {
				return v.AsReference().(ast.ExpressionNode)
			}).ToSlice()

			argNamedArgsTuple := args[4].AsReference().(value.ArrayTuple)
			argNamedArgs := value.TransformArrayTupleIntoNativeArrayTuple(argNamedArgsTuple, func(v value.Value) ast.NamedArgumentNode {
				return v.AsReference().(ast.NamedArgumentNode)
			}).ToSlice()

			var argLoc *position.Location
			if args[5].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[5].Pointer())
			}
			self := ast.NewCallNode(
				argLoc,
				argReceiver,
				argNilSafe,
				argPosArgs,
				argNamedArgs,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"receiver",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallNode)
			result := value.Ref(self.Receiver)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"nil_safe",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallNode)
			result := value.BoolVal(self.NilSafe)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallNode)
			entries := value.CastNativeArrayTuplePtr(&self.PositionalArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"named_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallNode)
			entries := value.CastNativeArrayTuplePtr(&self.NamedArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
