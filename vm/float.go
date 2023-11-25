package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Add(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Subtract(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Multiply(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"/",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Divide(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"**",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Exponentiate(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.GreaterThan(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.GreaterThanEqual(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.LessThan(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.LessThanEqual(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"%",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return value.ToValueErr(self.Modulo(other))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return self.Equal(other), nil
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"===",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			other := args[1]
			return self.StrictEqual(other), nil
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)

	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.String(self.Inspect()), nil
		},
	)
	value.FloatClass.DefineAliasString("to_string", "inspect")

	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_float",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self, nil
		},
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_int",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToInt(), nil
		},
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_float64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToFloat64(), nil
		},
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_float32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToFloat32(), nil
		},
	)

	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_int64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToInt64(), nil
		},
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_int32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToInt32(), nil
		},
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_int16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToInt16(), nil
		},
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_int8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToInt8(), nil
		},
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_uint64",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToUInt64(), nil
		},
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_uint32",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToUInt32(), nil
		},
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_uint16",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToUInt16(), nil
		},
	)
	DefineMethodWithOptions(
		value.FloatClass.Methods,
		"to_uint8",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Float)
			return self.ToUInt8(), nil
		},
	)
}
