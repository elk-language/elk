package runtime

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initInterpolatedRegexLiteralNode() {
	c := &value.InterpolatedRegexLiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {

			argContentTuple := args[1].AsReference().(value.ArrayTuple)
			argContent := value.TransformArrayTupleIntoNativeArrayTuple(argContentTuple, func(v value.Value) ast.RegexLiteralContentNode {
				return v.AsReference().(ast.RegexLiteralContentNode)
			}).ToSlice()

			var argFlags value.UInt8
			if !args[2].IsUndefined() {
				argFlags = args[2].AsUInt8()
			}

			var argLoc *position.Location
			if args[3].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[3].Pointer())
			}
			self := ast.NewInterpolatedRegexLiteralNode(
				argLoc,
				argContent,
				bitfield.BitField8FromInt(argFlags),
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"content",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			entries := value.CastNativeArrayTuplePtr(&self.Content)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"flags",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.UInt8(self.Flags.Byte()).ToValue()
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_case_insensitive",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.BoolVal(self.IsCaseInsensitive())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_multiline",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.BoolVal(self.IsMultiline())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_dot_all",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.BoolVal(self.IsDotAll())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_ungreedy",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.BoolVal(self.IsUngreedy())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_ascii",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.BoolVal(self.IsASCII())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_extended",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.BoolVal(self.IsExtended())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined
		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
