package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	// Instance methods
	c := &value.BigFloatClass.MethodContainer
	Def(
		c,
		"set_precision",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			arg := args[1]
			p, ok := value.ToGoUInt(arg)
			if !ok {
				return nil, value.NewBigFloatPrecisionError(arg.Inspect())
			}
			return self.SetPrecision(p), nil
		},
		DefWithParameters("precision"),
	)
	Alias(c, "p", "set_precision")
	Def(
		c,
		"precision",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return value.UInt64(self.Precision()), nil
		},
	)

	Def(
		c,
		"+@",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
		DefWithSealed(),
	)
	Def(
		c,
		"-@",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.Negate(), nil
		},
		DefWithSealed(),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.Add(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.Subtract(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.Multiply(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.Divide(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"**",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.Exponentiate(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.Compare(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.GreaterThan(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.GreaterThanEqual(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.LessThan(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.LessThanEqual(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return self.Equal(other), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"%",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			other := args[1]
			return value.ToValueErr(self.Modulo(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
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
		"to_big_float",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self, nil
		},
	)
	Def(
		c,
		"to_float",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToFloat(), nil
		},
	)
	Def(
		c,
		"to_int",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToInt(), nil
		},
	)
	Def(
		c,
		"to_float64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToFloat64(), nil
		},
	)
	Def(
		c,
		"to_float32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToFloat32(), nil
		},
	)

	Def(
		c,
		"to_int64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToInt64(), nil
		},
	)
	Def(
		c,
		"to_int32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToInt32(), nil
		},
	)
	Def(
		c,
		"to_int16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToInt16(), nil
		},
	)
	Def(
		c,
		"to_int8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToInt8(), nil
		},
	)
	Def(
		c,
		"to_uint64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToUInt64(), nil
		},
	)
	Def(
		c,
		"to_uint32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToUInt32(), nil
		},
	)
	Def(
		c,
		"to_uint16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToUInt16(), nil
		},
	)
	Def(
		c,
		"to_uint8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.BigFloat)
			return self.ToUInt8(), nil
		},
	)
}
