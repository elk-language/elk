package runtime

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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			var argParameters []ast.ParameterNode
			if !args[1].IsUndefined() {
				argParametersTuple := args[1].AsReference().(value.ArrayTuple)
				argParameters = value.TransformArrayTupleIntoNativeArrayTuple(argParametersTuple, func(v value.Value) ast.ParameterNode {
					return v.AsReference().(ast.ParameterNode)
				}).ToSlice()
			}

			var argReturnType ast.TypeNode
			if !args[2].IsUndefined() {
				argReturnType = args[2].MustReference().(ast.TypeNode)
			}

			var argThrowType ast.TypeNode
			if !args[3].IsUndefined() {
				argThrowType = args[3].MustReference().(ast.TypeNode)
			}

			var argBody []ast.StatementNode
			if !args[4].IsUndefined() {
				argBodyTuple := args[4].AsReference().(value.ArrayTuple)
				argBody = value.TransformArrayTupleIntoNativeArrayTuple(argBodyTuple, func(v value.Value) ast.StatementNode {
					return v.AsReference().(ast.StatementNode)
				}).ToSlice()
			}

			var argLambda bool
			if !args[5].IsUndefined() {
				argThrowType = args[5].MustReference().(ast.TypeNode)
			}

			var argLoc *position.Location
			if args[6].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[6].Pointer())
			}
			self := ast.NewClosureLiteralNode(
				argLoc,
				argParameters,
				argReturnType,
				argThrowType,
				argBody,
				argLambda,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(6),
	)

	vm.Def(
		c,
		"parameters",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			entries := value.CastNativeArrayTuplePtr(&self.Parameters)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"lambda",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			return value.BoolVal(self.Lambda), value.Undefined
		},
	)

	vm.Def(
		c,
		"return_type",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			if self.ReturnType == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.ReturnType), value.Undefined

		},
	)

	vm.Def(
		c,
		"throw_type",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			if self.ReturnType == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.ThrowType), value.Undefined
		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			entries := value.CastNativeArrayTuplePtr(&self.Body)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureLiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
