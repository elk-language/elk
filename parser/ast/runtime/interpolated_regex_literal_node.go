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
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			argContentTuple := args[0].MustReference().(*value.ArrayTuple)
			argContent := make([]ast.RegexLiteralContentNode, argContentTuple.Length())
			for i, el := range *argContentTuple {
				argContent[i] = el.MustReference().(ast.RegexLiteralContentNode)
			}
			var argFlags value.UInt8
			if !args[1].IsUndefined() {
				argFlags = args[1].AsUInt8()
			}

			var argSpan *position.Span
			if args[2].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[2].Pointer())
			}
			self := ast.NewInterpolatedRegexLiteralNode(
				argSpan,
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
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)

			collection := self.Content
			arrayTuple := value.NewArrayTupleWithLength(len(collection))
			for i, el := range collection {
				arrayTuple.SetAt(i, value.Ref(el))
			}
			result := value.Ref(arrayTuple)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"flags",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.UInt8(self.Flags.Byte()).ToValue()
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_case_insensitive",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.ToElkBool(self.IsCaseInsensitive())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_multiline",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.ToElkBool(self.IsMultiline())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_dot_all",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.ToElkBool(self.IsDotAll())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_ungreedy",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.ToElkBool(self.IsUngreedy())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_ascii",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.ToElkBool(self.IsASCII())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"is_extended",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.ToElkBool(self.IsExtended())
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InterpolatedRegexLiteralNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined
		},
	)

}
