package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Std::StringBuilder
func initStringBuilder() {
	// Instance methods
	c := value.StringBuilderClass

	Def(
		c,
		symbol.ReservedIdentifierName("init"),
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := value.NewStringBuilder()
			other := args[1]
			self.Write(other)

			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"<<",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()
			other := args[1]
			self.Write(other)

			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"<<@1",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()
			other := args[1].AsString()
			self.WriteString(other)

			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"<<@2",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()
			other := args[1].AsChar()
			self.WriteChar(other)

			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"<<@3",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()
			other := args[1].AsUInt8()
			self.WriteByte(other)

			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Alias(c, "push", "<<")
	Alias(c, "push@1", "<<@1")
	Alias(c, "push@2", "<<@2")
	Alias(c, "push@3", "<<@3")

	Def(
		c,
		"write",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()
			other := args[1]

			return self.Write(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"write@1",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()
			other := args[1].AsString()

			return self.WriteString(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"write@2",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()
			other := args[1].AsChar()

			return self.WriteChar(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"write@3",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()
			other := args[1].AsUInt8()
			self.WriteByte(other)
			return value.SmallInt(1).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"to_string",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()

			return self.ToString().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"inspect",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()

			return value.String(self.Inspect()).ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()

			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()

			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"grow",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()
			other := args[1].AsAnyInt()
			self.Grow(other)

			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"reset",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsStringBuilder()
			self.Reset()

			return self.ToValue(), value.Undefined
		},
	)

}
