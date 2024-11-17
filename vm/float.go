package vm

import (
	"github.com/elk-language/elk/value"
)

func initFloat() {
	// Instance methods
	c := &value.FloatClass.MethodContainer
	Def(
		c,
		"+@",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)
	Def(
		c,
		"-@",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return -self, nil
		},
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Add(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Subtract(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Multiply(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Divide(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"**",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Exponentiate(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Compare(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.GreaterThan(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.GreaterThanEqual(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.LessThan(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.LessThanEqual(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return self.Equal(other), nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"%",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Modulo(other))
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.String(self.Inspect()), nil
		},
	)
	Alias(c, "to_string", "inspect")

	Def(
		c,
		"to_float",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self, nil
		},
	)
	Def(
		c,
		"to_int",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToInt(), nil
		},
	)
	Def(
		c,
		"to_float64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToFloat64(), nil
		},
	)
	Def(
		c,
		"to_float32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToFloat32(), nil
		},
	)

	Def(
		c,
		"to_int64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToInt64(), nil
		},
	)
	Def(
		c,
		"to_int32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToInt32(), nil
		},
	)
	Def(
		c,
		"to_int16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToInt16(), nil
		},
	)
	Def(
		c,
		"to_int8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToInt8(), nil
		},
	)
	Def(
		c,
		"to_uint64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToUInt64(), nil
		},
	)
	Def(
		c,
		"to_uint32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToUInt32(), nil
		},
	)
	Def(
		c,
		"to_uint16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToUInt16(), nil
		},
	)
	Def(
		c,
		"to_uint8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToUInt8(), nil
		},
	)
}
