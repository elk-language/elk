package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

func init() {
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
		DefWithSealed(),
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
		DefWithSealed(),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.Add(other))
			case *value.BigInt:
				return value.ToValueErr(s.Add(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.Subtract(other))
			case *value.BigInt:
				return value.ToValueErr(s.Subtract(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.Multiply(other))
			case *value.BigInt:
				return value.ToValueErr(s.Multiply(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.Divide(other))
			case *value.BigInt:
				return value.ToValueErr(s.Divide(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"**",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.Exponentiate(other))
			case *value.BigInt:
				return value.ToValueErr(s.Exponentiate(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.Compare(other))
			case *value.BigInt:
				return value.ToValueErr(s.Compare(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.GreaterThan(other))
			case *value.BigInt:
				return value.ToValueErr(s.GreaterThan(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.GreaterThanEqual(other))
			case *value.BigInt:
				return value.ToValueErr(s.GreaterThanEqual(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.LessThan(other))
			case *value.BigInt:
				return value.ToValueErr(s.LessThan(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.LessThanEqual(other))
			case *value.BigInt:
				return value.ToValueErr(s.LessThanEqual(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
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
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.LeftBitshift(other))
			case *value.BigInt:
				return value.ToValueErr(s.LeftBitshift(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		">>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.RightBitshift(other))
			case *value.BigInt:
				return value.ToValueErr(s.RightBitshift(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"&",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.BitwiseAnd(other))
			case *value.BigInt:
				return value.ToValueErr(s.BitwiseAnd(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
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
		DefWithSealed(),
	)
	Def(
		c,
		"&~",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.BitwiseAndNot(other))
			case *value.BigInt:
				return value.ToValueErr(s.BitwiseAndNot(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"|",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.BitwiseOr(other))
			case *value.BigInt:
				return value.ToValueErr(s.BitwiseOr(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"^",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.BitwiseXor(other))
			case *value.BigInt:
				return value.ToValueErr(s.BitwiseXor(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"%",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.Modulo(other))
			case *value.BigInt:
				return value.ToValueErr(s.Modulo(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithParameters("other"),
		DefWithSealed(),
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
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return s.Negate(), nil
			case *value.BigInt:
				return s.Negate(), nil
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		DefWithSealed(),
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
}
