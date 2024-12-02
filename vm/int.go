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
			switch s := self.(type) {
			case value.SmallInt:
				return s.Increment(), nil
			case *value.BigInt:
				return s.Increment(), nil
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"--",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Decrement(), nil
			case *value.BigInt:
				return s.Decrement(), nil
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Add(other)
			case *value.BigInt:
				return s.Add(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Subtract(other)
			case *value.BigInt:
				return s.Subtract(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Multiply(other)
			case *value.BigInt:
				return s.Multiply(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Divide(other)
			case *value.BigInt:
				return s.Divide(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"**",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Exponentiate(other)
			case *value.BigInt:
				return s.Exponentiate(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Compare(other)
			case *value.BigInt:
				return s.Compare(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.GreaterThan(other)
			case *value.BigInt:
				return s.GreaterThan(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.GreaterThanEqual(other)
			case *value.BigInt:
				return s.GreaterThanEqual(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.LessThan(other)
			case *value.BigInt:
				return s.LessThan(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.LessThanEqual(other)
			case *value.BigInt:
				return s.LessThanEqual(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Equal(other), nil
			case *value.BigInt:
				return s.Equal(other), nil
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.LeftBitshift(other)
			case *value.BigInt:
				return s.LeftBitshift(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.RightBitshift(other)
			case *value.BigInt:
				return s.RightBitshift(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"&",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.BitwiseAnd(other)
			case *value.BigInt:
				return s.BitwiseAnd(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"~",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.BitwiseNot(), nil
			case *value.BigInt:
				return s.BitwiseNot(), nil
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"&~",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.BitwiseAndNot(other)
			case *value.BigInt:
				return s.BitwiseAndNot(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"|",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.BitwiseOr(other)
			case *value.BigInt:
				return s.BitwiseOr(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"^",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.BitwiseXor(other)
			case *value.BigInt:
				return s.BitwiseXor(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"%",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Modulo(other)
			case *value.BigInt:
				return s.Modulo(other)
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)
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
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Negate(), nil
			case *value.BigInt:
				return s.Negate(), nil
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return value.String(s.Inspect()), nil
			case *value.BigInt:
				return value.String(s.Inspect()), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Alias(c, "to_string", "inspect")

	Def(
		c,
		"to_int",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)
	Def(
		c,
		"to_float",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToFloat(), nil
			case *value.BigInt:
				return s.ToFloat(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"to_float64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToFloat64(), nil
			case *value.BigInt:
				return s.ToFloat64(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"to_float32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToFloat32(), nil
			case *value.BigInt:
				return s.ToFloat32(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"to_int64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToInt64(), nil
			case *value.BigInt:
				return s.ToInt64(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"to_int32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToInt32(), nil
			case *value.BigInt:
				return s.ToInt32(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"to_int16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToInt16(), nil
			case *value.BigInt:
				return s.ToInt16(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"to_int8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToInt8(), nil
			case *value.BigInt:
				return s.ToInt8(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"to_uint64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToUInt64(), nil
			case *value.BigInt:
				return s.ToUInt64(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"to_uint32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToUInt32(), nil
			case *value.BigInt:
				return s.ToUInt32(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"to_uint16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToUInt16(), nil
			case *value.BigInt:
				return s.ToUInt16(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"to_uint8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.ToUInt8(), nil
			case *value.BigInt:
				return s.ToUInt8(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Def(
		c,
		"times",
		func(v *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			fn := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				for i := range s {
					_, err := v.CallCallable(fn, value.SmallInt(i))
					if err != nil {
						return nil, err
					}
				}
				return value.Nil, nil
			case *value.BigInt:
				if s.IsSmallInt() {
					for i := range s.ToSmallInt() {
						_, err := v.CallCallable(fn, value.SmallInt(i))
						if err != nil {
							return nil, err
						}
					}
					return value.Nil, nil
				}
				for i := range value.MaxSmallInt {
					_, err := v.CallCallable(fn, value.SmallInt(i))
					if err != nil {
						return nil, err
					}
				}

				sGo := s.ToGoBigInt()
				one := big.NewInt(1)
				bigI := big.NewInt(value.MaxSmallInt)
				for ; bigI.Cmp(sGo) == -1; bigI.Add(bigI, one) {
					v.CallCallable(fn, value.ToElkBigInt(bigI))
				}
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"nanoseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Nanoseconds(), nil
			case *value.BigInt:
				return s.Nanoseconds(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Alias(c, "nanosecond", "nanoseconds")

	Def(
		c,
		"microseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Microseconds(), nil
			case *value.BigInt:
				return s.Microseconds(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Alias(c, "microsecond", "microseconds")

	Def(
		c,
		"milliseconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Milliseconds(), nil
			case *value.BigInt:
				return s.Milliseconds(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Alias(c, "millisecond", "milliseconds")

	Def(
		c,
		"seconds",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Seconds(), nil
			case *value.BigInt:
				return s.Seconds(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Alias(c, "second", "seconds")

	Def(
		c,
		"minutes",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Minutes(), nil
			case *value.BigInt:
				return s.Minutes(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Alias(c, "minute", "minutes")

	Def(
		c,
		"hours",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Hours(), nil
			case *value.BigInt:
				return s.Hours(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Alias(c, "hour", "hours")

	Def(
		c,
		"days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Days(), nil
			case *value.BigInt:
				return s.Days(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Alias(c, "day", "days")

	Def(
		c,
		"weeks",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Weeks(), nil
			case *value.BigInt:
				return s.Weeks(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Alias(c, "week", "weeks")

	Def(
		c,
		"years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Years(), nil
			case *value.BigInt:
				return s.Years(), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)
	Alias(c, "year", "years")
}
