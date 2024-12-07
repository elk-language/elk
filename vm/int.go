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
				return s.Decrement(), value.Undefined
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
				return self.AsSmallInt().Add(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.Add(other)
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
				return self.AsSmallInt().Subtract(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.Subtract(other)
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
				return self.AsSmallInt().Multiply(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.Multiply(other)
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
				return self.AsSmallInt().Divide(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.Divide(other)
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
				return self.AsSmallInt().Exponentiate(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.Exponentiate(other)
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
				return self.AsSmallInt().Compare(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.Compare(other)
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
				return self.AsSmallInt().GreaterThan(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.GreaterThan(other)
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
				return self.AsSmallInt().GreaterThanEqual(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.GreaterThanEqual(other)
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
				return self.AsSmallInt().LessThan(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.LessThan(other)
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
				return self.AsSmallInt().LessThanEqual(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.LessThanEqual(other)
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
				return self.AsSmallInt().Equal(other), value.Undefined
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.Equal(other), value.Undefined
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
				return self.AsSmallInt().LeftBitshift(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.LeftBitshift(other)
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
				return self.AsSmallInt().RightBitshift(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.RightBitshift(other)
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
				return self.AsSmallInt().BitwiseAnd(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.BitwiseAnd(other)
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
				return self.AsSmallInt().BitwiseAndNot(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.BitwiseAndNot(other)
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
				return self.AsSmallInt().BitwiseOr(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.BitwiseOr(other)
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
				return self.AsSmallInt().BitwiseXor(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.BitwiseXor(other)
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
				return self.AsSmallInt().Modulo(other)
			}

			switch s := self.SafeAsReference().(type) {
			case *value.BigInt:
				return s.Modulo(other)
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
				return self.AsSmallInt().Negate(), value.Undefined
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
}
