package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initSwitchCaseNode() {
	c := &value.SwitchCaseNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argPattern := args[1].MustReference().(ast.PatternNode)

			argBodyTuple := args[2].AsReference().(value.ArrayTuple)
			argBody := value.TransformArrayTupleIntoNativeArrayTuple(argBodyTuple, func(v value.Value) ast.StatementNode {
				return v.AsReference().(ast.StatementNode)
			}).ToSlice()

			var argLoc *position.Location
			if args[3].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[3].Pointer())
			}
			self := ast.NewSwitchCaseNode(
				argLoc,
				argPattern,
				argBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"pattern_node",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchCaseNode)
			result := value.Ref(self.Pattern)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchCaseNode)
			entries := value.CastNativeArrayTuplePtr(&self.Body)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchCaseNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchCaseNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchCaseNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
