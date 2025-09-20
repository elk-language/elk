package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Int32
func initInt32() {
	// Instance methods
	c := &value.Int32Class.MethodContainer
	Def(
		c,
		"hash",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.Hash().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"++",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return (self + 1).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"--",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return (self - 1).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.ToValueErr(self.Add(args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.ToValueErr(self.Subtract(args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.ToValueErr(self.Multiply(args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.ToValueErr(self.Divide(args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"**",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.ToValueErr(self.ExponentiateVal(args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.CompareVal(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.GreaterThanVal(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.GreaterThanEqualVal(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.LessThanVal(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.LessThanEqualVal(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.EqualVal(args[1]), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.ToValueErr(value.StrictIntLeftBitshift(self, args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.ToValueErr(value.StrictIntLogicalLeftBitshift(self, args[1], value.LogicalRightShift32))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.ToValueErr(value.StrictIntRightBitshift(self, args[1]))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">>>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.ToValueErr(value.StrictIntLogicalRightBitshift(self, args[1], value.LogicalRightShift32))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"&",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			other := args[1].AsInt32()

			return (self & other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"~",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return (^self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"&~",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			other := args[1].AsInt32()

			return (self &^ other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"|",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			other := args[1].AsInt32()

			return (self | other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"^",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			other := args[1].AsInt32()

			return (self ^ other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"%",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			other := args[1].AsInt32()

			return (self % other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
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
			self := args[0].AsInt32()

			return (-self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return value.Ref(self.ToString()), value.Undefined
		},
	)

	Def(
		c,
		"to_int",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToSmallInt().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_float",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToFloat().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_float64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToFloat64().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_float32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToFloat32().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToInt64().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"to_int16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToInt16().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToInt8().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToUInt().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToUInt64().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToUInt32().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToUInt16().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt32()
			return self.ToUInt8().ToValue(), value.Undefined
		},
	)
}
