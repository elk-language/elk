package vm

import (
	"github.com/elk-language/elk/value"
)

func initBigFloat() {
	// Instance methods
	c := &value.BigFloatClass.MethodContainer
	Def(
		c,
		"set_precision",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			arg := args[1]
			p, ok := value.ToGoUInt(arg)
			if !ok {
				return value.Undefined, value.Ref(value.NewBigFloatPrecisionError(arg.Inspect()))
			}
			return value.Ref(self.SetPrecision(p)), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "p", "set_precision")
	Def(
		c,
		"precision",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return value.UInt64(self.Precision()).ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"hash",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.Hash().ToValue(), value.Undefined
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
		"-@",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return value.Ref(self.Negate()), value.Undefined
		},
	)

	Def(
		c,
		"+",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.AddVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@1",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.AddInt(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@2",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1].AsFloat()
			return self.AddFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+@3",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := (*value.BigFloat)(args[1].Pointer())
			return self.AddBigFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"-",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.SubtractVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@1",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.SubtractInt(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@2",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1].AsFloat()
			return self.SubtractFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-@3",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := (*value.BigFloat)(args[1].Pointer())
			return self.SubtractBigFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"*",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.MultiplyVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*@1",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.MultiplyInt(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*@2",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1].AsFloat()
			return self.MultiplyFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*@3",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := (*value.BigFloat)(args[1].Pointer())
			return self.MultiplyBigFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"/",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.DivideVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/@1",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.DivideInt(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/@2",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1].AsFloat()
			return self.DivideFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/@3",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := (*value.BigFloat)(args[1].Pointer())
			return self.DivideBigFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"**",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.ExponentiateVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"**@1",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.ExponentiateInt(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"**@2",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1].AsFloat()
			return self.ExponentiateFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"**@3",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := (*value.BigFloat)(args[1].Pointer())
			return self.ExponentiateBigFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"%",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.ModuloVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"%@1",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.ModuloInt(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"%@2",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1].AsFloat()
			return self.ModuloFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"%@3",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := (*value.BigFloat)(args[1].Pointer())
			return self.ModuloBigFloat(other).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"<=>",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.CompareVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.GreaterThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.GreaterThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.LessThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.LessThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.EqualVal(other), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			other := args[1]
			return self.LaxEqualVal(other), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"inspect",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)
	Alias(c, "to_string", "inspect")

	Def(
		c,
		"to_big_float",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self, value.Undefined
		},
	)
	Def(
		c,
		"to_float",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToFloat().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToInt(), value.Undefined
		},
	)
	Def(
		c,
		"to_float64",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToFloat64().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_float32",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToFloat32().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"to_int64",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToInt64().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int32",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToInt32().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int16",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToInt16().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_int8",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToInt8().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToUInt().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint64",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToUInt64().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint32",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToUInt32().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint16",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToUInt16().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_uint8",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.BigFloat)(args[0].Pointer())
			return self.ToUInt8().ToValue(), value.Undefined
		},
	)
}
