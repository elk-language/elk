package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Size
func initSize() {
	// Instance methods
	c := &value.SizeClass.MethodContainer

	Def(
		c,
		"-@",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return (-self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"+@",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"+",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return (self + other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return (self - other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"*",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustSize()
			other := args[1]
			return value.ToValueErr(self.Multiply(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*@1",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1]
			return value.ToValueErr(self.MultiplyInt(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*@2",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsFloat()
			return value.ToValueErr(self.MultiplyFloat(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*@3",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := (*value.BigFloat)(args[1].Pointer())
			return value.ToValueErr(self.MultiplyBigFloat(other))
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"/",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustSize()
			other := args[1]
			return value.ToValueErr(self.Divide(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/@1",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1]
			return value.ToValueErr(self.DivideInt(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/@2",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsFloat()
			return value.ToValueErr(self.DivideFloat(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/@3",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := (*value.BigFloat)(args[1].Pointer())
			return value.ToValueErr(self.DivideBigFloat(other))
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"<=>",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return value.SmallInt(self.Cmp(other)).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return value.BoolVal(self > other), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return value.BoolVal(self >= other), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return value.BoolVal(self < other), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1].AsSize()
			return value.BoolVal(self <= other), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "=~", "==")

	Def(
		c,
		"to_string",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return value.Ref(self.ToString()), value.Undefined
		},
	)
	Alias(c, "to_string", "inspect")

	Def(
		c,
		"bytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Bytes(), value.Undefined
		},
	)
	Alias(c, "byte", "bytes")
	Def(
		c,
		"in_bytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InBytes().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"kilobytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Kilobytes(), value.Undefined
		},
	)
	Alias(c, "kilobyte", "kilobytes")
	Def(
		c,
		"in_kilobytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InKilobytes().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"megabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Megabytes(), value.Undefined
		},
	)
	Alias(c, "megabyte", "megabytes")
	Def(
		c,
		"in_megabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InMegabytes().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"gigabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Gigabytes(), value.Undefined
		},
	)
	Alias(c, "gigabyte", "gigabytes")
	Def(
		c,
		"in_gigabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InGigabytes().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"terabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Terabytes(), value.Undefined
		},
	)
	Alias(c, "terabyte", "terabytes")
	Def(
		c,
		"in_terabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InTerabytes().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"petabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.Petabytes(), value.Undefined
		},
	)
	Alias(c, "petabyte", "petabytes")
	Def(
		c,
		"in_petabytes",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsSize()
			return self.InPetabytes().ToValue(), value.Undefined
		},
	)

}
