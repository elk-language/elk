package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initInitDefinitionNode() {
	c := &value.InitDefinitionNodeClass.MethodContainer
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

			var argBody []ast.StatementNode
			if !args[2].IsUndefined() {
				argBodyTuple := args[2].AsReference().(value.ArrayTuple)
				argBody = value.TransformArrayTupleIntoNativeArrayTuple(argBodyTuple, func(v value.Value) ast.StatementNode {
					return v.AsReference().(ast.StatementNode)
				}).ToSlice()
			}

			var argThrowType ast.TypeNode
			if !args[3].IsUndefined() {
				argThrowType = args[3].MustReference().(ast.TypeNode)
			}

			var argLocation *position.Location
			if args[4].IsUndefined() {
				argLocation = position.ZeroLocation
			} else {
				argLocation = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewInitDefinitionNode(
				argLocation,
				argParameters,
				argThrowType,
				argBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"parameters",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)

			entries := value.CastNativeArrayTuplePtr(&self.Parameters)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"throw_type",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)
			if self.ThrowType == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.ThrowType), value.Undefined

		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)
			entries := value.CastNativeArrayTuplePtr(&self.Body)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InitDefinitionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
