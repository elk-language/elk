package vm

import (
	"fmt"

	"math/big"

	"github.com/elk-language/elk/value"
)

// Std::Int
func initInt() {
	// Instance methods
	c := &value.IntClass.MethodContainer
	Def(
		c,
		"++",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Increment(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return value.Ref(s.Increment()), value.Undefined
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"--",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Decrement(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.DecrementVal(), value.Undefined
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().AddVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.AddVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().SubtractVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.SubtractVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().MultiplyVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.MultiplyVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().DivideVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.DivideVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"**",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().ExponentiateVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ExponentiateVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().CompareVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.CompareVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().GreaterThanVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.GreaterThanVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().GreaterThanEqualVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.GreaterThanEqualVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().LessThanVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.LessThanVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().LessThanEqualVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.LessThanEqualVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().EqualVal(other), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.EqualVal(other), value.Undefined
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().LeftBitshiftVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.LeftBitshiftVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().RightBitshiftVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.RightBitshiftVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"&",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().BitwiseAndVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.BitwiseAndVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"~",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().BitwiseNot().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return value.Ref(s.BitwiseNot()), value.Undefined
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"&~",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().BitwiseAndNotVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.BitwiseAndNotVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"|",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().BitwiseOrVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.BitwiseOrVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"^",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().BitwiseXorVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.BitwiseXorVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"%",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			if self.IsSmallInt() {
				return self.AsSmallInt().ModuloVal(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ModuloVal(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
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
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().NegateVal(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return value.Ref(s.Negate()), value.Undefined
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return value.Ref(value.String(self.AsSmallInt().Inspect())), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return value.Ref(value.String(s.Inspect())), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "to_string", "inspect")

	Def(
		c,
		"to_int",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"to_float",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToFloat().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToFloat().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"to_float64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToFloat64().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToFloat64().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"to_float32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToFloat32().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToFloat32().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"to_int64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToInt64().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToInt64().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"to_int32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToInt32().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToInt32().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"to_int16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToInt16().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToInt16().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"to_int8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToInt8().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToInt8().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"to_uint64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToUInt64().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToUInt64().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"to_uint32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToUInt32().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToUInt32().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"to_uint16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToUInt16().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToUInt16().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"to_uint8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().ToUInt8().ToValue(), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.ToUInt8().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Def(
		c,
		"times",
		func(v *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			fn := args[1]

			if self.IsSmallInt() {
				for i := range self.AsSmallInt() {
					_, err := v.CallCallable(fn, value.SmallInt(i).ToValue())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Nil, value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				if s.IsSmallInt() {
					for i := range s.ToSmallInt() {
						_, err := v.CallCallable(fn, value.SmallInt(i).ToValue())
						if !err.IsUndefined() {
							return value.Undefined, err
						}
					}
					return value.Nil, value.Undefined
				}
				for i := range value.MaxSmallInt {
					_, err := v.CallCallable(fn, value.SmallInt(i).ToValue())
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}

				sGo := s.ToGoBigInt()
				one := big.NewInt(1)
				bigI := big.NewInt(value.MaxSmallInt)
				for ; bigI.Cmp(sGo) == -1; bigI.Add(bigI, one) {
					v.CallCallable(fn, value.Ref(value.ToElkBigInt(bigI)))
				}
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Nanoseconds().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Nanoseconds().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "nanosecond", "nanoseconds")

	Def(
		c,
		"microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Microseconds().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Microseconds().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "microsecond", "microseconds")

	Def(
		c,
		"milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Milliseconds().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Milliseconds().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "millisecond", "milliseconds")

	Def(
		c,
		"seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Seconds().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Seconds().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "second", "seconds")

	Def(
		c,
		"minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Minutes().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Minutes().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "minute", "minutes")

	Def(
		c,
		"hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Hours().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Hours().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "hour", "hours")

	Def(
		c,
		"days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Days().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Days().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "day", "days")

	Def(
		c,
		"weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Weeks().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Weeks().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "week", "weeks")
	Def(
		c,
		"months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Months().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Months().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "month", "months")
	Def(
		c,
		"years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Years().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Years().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "year", "years")
	Def(
		c,
		"centuries",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Centuries().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Centuries().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "century", "centuries")
	Def(
		c,
		"millenia",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return self.AsSmallInt().Millenia().ToValue(), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return s.Millenia().ToValue(), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
	Alias(c, "millenium", "millenia")

	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			if self.IsSmallInt() {
				return value.Ref(value.NewSmallIntIterator(self.AsSmallInt())), value.Undefined
			}
			switch s := self.AsReference().(type) {
			case *value.BigInt:
				return value.Ref(value.NewBigIntIterator(s)), value.Undefined
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %s", self.Inspect()))
		},
	)
}

// ::Std::Int::Iterator
func initIntIterator() {
	// Instance methods
	c := &value.IntIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch self := self.SafeAsReference().(type) {
			case *value.SmallIntIterator:
				return self.Next()
			case *value.BigIntIterator:
				return self.Next()
			default:
				panic(fmt.Sprintf("expected SmallIntIterator or BigIntIterator, got: %s", self.Inspect()))
			}
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch self := self.SafeAsReference().(type) {
			case *value.SmallIntIterator:
				self.Reset()
				return args[0], value.Undefined
			case *value.BigIntIterator:
				self.Reset()
				return args[0], value.Undefined
			default:
				panic(fmt.Sprintf("expected SmallIntIterator or BigIntIterator, got: %s", self.Inspect()))
			}
		},
	)
}
