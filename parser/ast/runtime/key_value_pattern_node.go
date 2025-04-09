package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initKeyValuePatternNode() {
	c := &value.KeyValuePatternNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argKey := args[1].MustReference().(ast.PatternExpressionNode)
			argValue := args[2].MustReference().(ast.PatternNode)

			var argLoc *position.Location
			if args[3].IsUndefined() {
				argLoc = position.DefaultLocation
			} else {
				argLoc = (*position.Location)(args[3].Pointer())
			}
			self := ast.NewKeyValuePatternNode(
				argLoc,
				argKey,
				argValue,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"key",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.KeyValuePatternNode)
			result := value.Ref(self.Key)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"value",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.KeyValuePatternNode)
			result := value.Ref(self.Value)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.KeyValuePatternNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.KeyValuePatternNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.KeyValuePatternNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
