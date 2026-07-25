package runtime

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initAttrDeclarationNode() {
	c := &value.AttrDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			entriesTuple := args[1].AsReference().(value.ArrayTuple)
			entries := value.TransformArrayTupleIntoNativeArrayTuple(entriesTuple, func(v value.Value) ast.ParameterNode {
				return v.AsReference().(ast.ParameterNode)
			}).ToSlice()

			var argFlags bitfield.BitFlag8
			if !args[2].IsUndefined() {
				argFlags = bitfield.BitFlag8(args[2].AsUInt8())
			}

			var docComment string
			if !args[3].IsUndefined() {
				docComment = (string)(args[3].MustReference().(value.String))
			}
			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewAttrDeclarationNode(
				argLoc,
				docComment,
				argFlags,
				entries,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"entries",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttrDeclarationNode)
			entries := value.CastNativeArrayTuplePtr(&self.Entries)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"doc_comment",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttrDeclarationNode)
			result := value.Ref((value.String)(self.DocComment()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttrDeclarationNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"flags",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttrDeclarationNode)
			result := value.UInt8(self.Flags.Byte()).ToValue()
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_pure",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttrDeclarationNode)
			result := value.BoolVal(self.IsPure())
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttrDeclarationNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.AttrDeclarationNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
