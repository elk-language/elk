package vm

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Std::String::Span
func initSpan() {
	// Instance methods
	c := &value.SpanClass.MethodContainer

	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Span)(args[0].Pointer())
			startPos := (*value.Position)(args[1].Pointer())
			endPos := (*value.Position)(args[2].Pointer())

			self.StartPos = (*position.Position)(startPos)
			self.EndPos = (*position.Position)(endPos)
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Span)(args[0].Pointer())
			other := (*value.Span)(args[1].Pointer())
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"start_pos",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Span)(args[0].Pointer())
			return value.Ref(self.StartPosition()), value.Undefined
		},
	)
	Def(
		c,
		"end_pos",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Span)(args[0].Pointer())
			return value.Ref(self.EndPosition()), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Span)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

}
