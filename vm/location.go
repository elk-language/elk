package vm

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Std::FS::Location
func initLocation() {
	// Instance methods
	c := &value.SpanClass.MethodContainer

	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Location)(args[0].Pointer())
			path := (*value.Path)(args[1].Pointer())
			span := (*value.Span)(args[2].Pointer())

			self.FilePath = path.String()
			self.Span = *(*position.Span)(span)
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Location)(args[0].Pointer())
			other := (*value.Location)(args[1].Pointer())
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"span",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Location)(args[0].Pointer())
			return value.Ref(self.SpanValue()), value.Undefined
		},
	)
	Def(
		c,
		"file_path",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Location)(args[0].Pointer())
			return value.Ref(value.NewPath(self.FilePath)), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Location)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

}
