package vm

import (
	"github.com/elk-language/elk/value"
)

func initFloat32() {
	// Instance methods
	c := &value.Float32Class.MethodContainer
	Def(
		c,
		"hash",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
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
			self := args[0].AsFloat32()
			return (-self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.ToValueErr(self.Add(args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.ToValueErr(self.Subtract(args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.ToValueErr(self.Multiply(args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.ToValueErr(self.Divide(args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"**",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.ToValueErr(self.ExponentiateVal(args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			other := args[1]
			return self.CompareVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			other := args[1]
			return self.GreaterThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			other := args[1]
			return self.GreaterThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			other := args[1]
			return self.LessThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			other := args[1]
			return self.LessThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			other := args[1]
			return self.EqualVal(other), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			other := args[1]
			return value.StrictFloatLaxEqual(self, other), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"%",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			other := args[1]
			return value.ToValueErr(self.ModuloVal(other))
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.Ref(self.ToString()), value.Undefined
		},
	)

	Def(
		c,
		"to_float",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.Float(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.SmallInt(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_float64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.Float64(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_float64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)

	Def(
		c,
		"to_int64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.Int64(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.Int32(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.Int16(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.Int8(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.UInt(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.UInt64(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.UInt32(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.UInt16(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			return value.UInt8(self).ToValue(), value.Undefined
		},
	)
}
