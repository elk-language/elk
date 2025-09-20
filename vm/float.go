package vm

import (
	"github.com/elk-language/elk/value"
)

func initFloat() {
	// Instance methods
	c := &value.FloatClass.MethodContainer
	Def(
		c,
		"hash",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.Hash().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"+@",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"-@",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return (-self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.AddVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.SubtractVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.MultiplyVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.DivideVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"**",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.ExponentiateVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.CompareVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.GreaterThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.GreaterThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.LessThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.LessThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.EqualVal(other), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.LaxEqualVal(other), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"%",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			other := args[1]
			return self.ModuloVal(other)
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)
	Alias(c, "to_string", "inspect")

	Def(
		c,
		"to_float",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self, value.Undefined
		},
	)
	Def(
		c,
		"to_int",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToInt(), value.Undefined
		},
	)
	Def(
		c,
		"to_float64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToFloat64().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_float32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToFloat32().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"to_int64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToInt64().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToInt32().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToInt16().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToInt8().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToUInt().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToUInt64().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToUInt32().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToUInt16().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.ToUInt8().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.Nanoseconds().ToValue(), value.Undefined
		},
	)
	Alias(c, "nanosecond", "nanoseconds")

	Def(
		c,
		"microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.Microseconds().ToValue(), value.Undefined
		},
	)
	Alias(c, "microsecond", "microseconds")

	Def(
		c,
		"milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.Milliseconds().ToValue(), value.Undefined
		},
	)
	Alias(c, "millisecond", "milliseconds")

	Def(
		c,
		"seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.Seconds().ToValue(), value.Undefined
		},
	)
	Alias(c, "second", "seconds")

	Def(
		c,
		"minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.Minutes().ToValue(), value.Undefined
		},
	)
	Alias(c, "minute", "minutes")

	Def(
		c,
		"hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return self.Hours().ToValue(), value.Undefined
		},
	)
	Alias(c, "hour", "hours")

	Def(
		c,
		"days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return value.Ref(self.Days()), value.Undefined
		},
	)
	Alias(c, "day", "days")

	Def(
		c,
		"weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return value.Ref(self.Weeks()), value.Undefined
		},
	)
	Alias(c, "week", "weeks")

	Def(
		c,
		"months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return value.Ref(self.Months()), value.Undefined
		},
	)
	Alias(c, "month", "months")

	Def(
		c,
		"years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return value.Ref(self.Years()), value.Undefined
		},
	)
	Alias(c, "year", "years")

	Def(
		c,
		"centuries",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return value.Ref(self.Centuries()), value.Undefined
		},
	)
	Alias(c, "century", "centuries")

	Def(
		c,
		"millenia",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat()
			return value.Ref(self.Millenia()), value.Undefined
		},
	)
	Alias(c, "millenium", "millenia")
}
