package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initInferredObjectPatternNode() {
	c := &value.InferredObjectPatternNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			var argAttributes []ast.PatternNode
			if !args[1].IsUndefined() {
				argAttributesTuple := args[1].AsReference().(value.ArrayTuple)
				argAttributes = value.TransformArrayTupleIntoNativeArrayTuple(argAttributesTuple, func(v value.Value) ast.PatternNode {
					return v.AsReference().(ast.PatternNode)
				}).ToSlice()
			}

			var argLoc *position.Location
			if args[2].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[2].Pointer())
			}
			self := ast.NewInferredObjectPatternNode(
				argLoc,
				argAttributes,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"attributes",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InferredObjectPatternNode)
			entries := value.CastNativeArrayTuplePtr(&self.Attributes)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InferredObjectPatternNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InferredObjectPatternNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InferredObjectPatternNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
