package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initConstructorCallNode() {
	c := &value.ConstructorCallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argClassNode := args[1].MustReference().(ast.ComplexConstantNode)

			var argPosArgs []ast.ExpressionNode
			if !args[2].IsUndefined() {
				argPosArgsTuple := args[2].AsReference().(value.ArrayTuple)
				argPosArgs = value.TransformArrayTupleIntoNativeArrayTuple(argPosArgsTuple, func(v value.Value) ast.ExpressionNode {
					return v.AsReference().(ast.ExpressionNode)
				}).ToSlice()
			}

			var argNamedArgs []ast.NamedArgumentNode
			if !args[3].IsUndefined() {
				argNamedArgsTuple := args[3].AsReference().(value.ArrayTuple)
				argNamedArgs = value.TransformArrayTupleIntoNativeArrayTuple(argNamedArgsTuple, func(v value.Value) ast.NamedArgumentNode {
					return v.AsReference().(ast.NamedArgumentNode)
				}).ToSlice()
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewConstructorCallNode(
				argLoc,
				argClassNode,
				argPosArgs,
				argNamedArgs,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"class_node",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstructorCallNode)
			result := value.Ref(self.ClassNode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstructorCallNode)
			entries := value.CastNativeArrayTuplePtr(&self.PositionalArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"named_arguments",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstructorCallNode)
			entries := value.CastNativeArrayTuplePtr(&self.NamedArguments)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstructorCallNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstructorCallNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstructorCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
