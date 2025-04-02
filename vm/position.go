package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::String::Position
func initPosition() {
	// Instance methods
	c := &value.PositionClass.MethodContainer

	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Position)(args[0].Pointer())

			byteOffset, ok := value.IntToGoInt(args[1])
			if !ok {
				return value.Undefined, value.Ref(value.NewError(value.OutOfRangeErrorClass, "too large integer value in position"))
			}
			line, ok := value.IntToGoInt(args[2])
			if !ok {
				return value.Undefined, value.Ref(value.NewError(value.OutOfRangeErrorClass, "too large integer value in position"))
			}
			column, ok := value.IntToGoInt(args[3])
			if !ok {
				return value.Undefined, value.Ref(value.NewError(value.OutOfRangeErrorClass, "too large integer value in position"))
			}

			self.ByteOffset = byteOffset
			self.Line = line
			self.Column = column
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(3),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Position)(args[0].Pointer())
			other := (*value.Position)(args[1].Pointer())
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"byte_offset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Position)(args[0].Pointer())
			return value.SmallInt(self.ByteOffset).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"line",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Position)(args[0].Pointer())
			return value.SmallInt(self.Line).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"column",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Position)(args[0].Pointer())
			return value.SmallInt(self.Column).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Position)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

}
