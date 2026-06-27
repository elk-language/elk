package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initAliasDeclarationNode() {
	c := &value.AliasDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {

			arg0Tuple := args[1].AsReference().(value.ArrayTuple)
			arg0 := value.TransformArrayTupleIntoNativeArrayTuple(arg0Tuple, func(v value.Value) *ast.AliasDeclarationEntry {
				return (*ast.AliasDeclarationEntry)(v.Pointer())
			})

			var argLoc *position.Location
			if args[2].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[2].Pointer())
			}
			self := ast.NewAliasDeclarationNode(
				argLoc,
				arg0.ToSlice(),
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"entries",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AliasDeclarationNode)
			entries := value.CastNativeArrayTuplePtr(&self.Entries)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AliasDeclarationNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AliasDeclarationNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AliasDeclarationNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
