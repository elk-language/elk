package value_test

import (
	"math"
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestStrictFloat_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Float64
		err  *value.Error
	}{
		"exponentiate String and return an error": {
			a:   value.Float64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Float64`"),
		},
		"exponentiate Int32 and return an error": {
			a:   value.Float64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Float64`"),
		},
		"exponentiate positive value.Float64": {
			a:    value.Float64(5.5),
			b:    value.Float64(3),
			want: value.Float64(166.375),
		},
		"exponentiate negative value.Float64": {
			a:    value.Float64(5.5),
			b:    value.Float64(-2),
			want: value.Float64(0.03305785123966942),
		},
		"exponentiate zero": {
			a:    value.Float64(5.5),
			b:    value.Float64(0),
			want: value.Float64(1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictFloatExponentiate(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictInt_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"exponentiate String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"exponentiate Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"exponentiate positive Int64 5 ** 2": {
			a:    value.Int64(5),
			b:    value.Int64(2),
			want: value.Int64(25),
		},
		"exponentiate positive Int64 7 ** 8": {
			a:    value.Int64(7),
			b:    value.Int64(8),
			want: value.Int64(5764801),
		},
		"exponentiate positive Int64 2 ** 5": {
			a:    value.Int64(2),
			b:    value.Int64(5),
			want: value.Int64(32),
		},
		"exponentiate positive Int64 6 ** 1": {
			a:    value.Int64(6),
			b:    value.Int64(1),
			want: value.Int64(6),
		},
		"exponentiate negative Int64": {
			a:    value.Int64(4),
			b:    value.Int64(-2),
			want: value.Int64(1),
		},
		"exponentiate zero": {
			a:    value.Int64(25),
			b:    value.Int64(0),
			want: value.Int64(1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictIntExponentiate(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumeric_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"add String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"add Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"add positive Int64": {
			a:    value.Int64(53),
			b:    value.Int64(21),
			want: value.Int64(74),
		},
		"add negative Int64": {
			a:    value.Int64(25),
			b:    value.Int64(-50),
			want: value.Int64(-25),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictNumericAdd(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumeric_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"subtract String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"subtract Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"subtract positive Int64": {
			a:    value.Int64(53),
			b:    value.Int64(21),
			want: value.Int64(32),
		},
		"subtract negative Int64": {
			a:    value.Int64(25),
			b:    value.Int64(-50),
			want: value.Int64(75),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictNumericSubtract(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumeric_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"multiply String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"multiply Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"multiply positive Int64": {
			a:    value.Int64(53),
			b:    value.Int64(2),
			want: value.Int64(106),
		},
		"multiply negative Int64": {
			a:    value.Int64(25),
			b:    value.Int64(-2),
			want: value.Int64(-50),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictNumericMultiply(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumeric_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Bool
		err  *value.Error
	}{
		"String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"15i64 > 30i64": {
			a:    value.Int64(15),
			b:    value.Int64(30),
			want: value.False,
		},
		"780i64 > -800i64": {
			a:    value.Int64(780),
			b:    value.Int64(-800),
			want: value.True,
		},
		"15i64 > 15i64": {
			a:    value.Int64(15),
			b:    value.Int64(15),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictNumericGreaterThan(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumeric_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Bool
		err  *value.Error
	}{
		"String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"15i64 >= 30i64": {
			a:    value.Int64(15),
			b:    value.Int64(30),
			want: value.False,
		},
		"780i64 >= -800i64": {
			a:    value.Int64(780),
			b:    value.Int64(-800),
			want: value.True,
		},
		"15i64 >= 15i64": {
			a:    value.Int64(15),
			b:    value.Int64(15),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictNumericGreaterThanEqual(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumeric_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Bool
		err  *value.Error
	}{
		"String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"15i64 < 30i64": {
			a:    value.Int64(15),
			b:    value.Int64(30),
			want: value.True,
		},
		"780i64 < -800i64": {
			a:    value.Int64(780),
			b:    value.Int64(-800),
			want: value.False,
		},
		"15i64 < 15i64": {
			a:    value.Int64(15),
			b:    value.Int64(15),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictNumericLessThan(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumeric_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Bool
		err  *value.Error
	}{
		"String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"15i64 <= 30i64": {
			a:    value.Int64(15),
			b:    value.Int64(30),
			want: value.True,
		},
		"780i64 <= -800i64": {
			a:    value.Int64(780),
			b:    value.Int64(-800),
			want: value.False,
		},
		"15i64 <= 15i64": {
			a:    value.Int64(15),
			b:    value.Int64(15),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictNumericLessThanEqual(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumeric_StrictEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Bool
		err  *value.Error
	}{
		"String 5.5f64 === '5.5'": {
			a:    value.Float64(5.5),
			b:    value.String("5.5"),
			want: value.False,
		},

		"Int64 5.0f64 === 5i64": {
			a:    value.Float64(5),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int64 5.5f64 === 5i64": {
			a:    value.Float64(5.5),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int64 NaN === 0i64": {
			a:    value.Float64NaN(),
			b:    value.Int64(0),
			want: value.False,
		},
		"Int64 +Inf === 69i64": {
			a:    value.Float64Inf(),
			b:    value.Int64(69),
			want: value.False,
		},
		"Int64 -Inf === -89i64": {
			a:    value.Float64NegInf(),
			b:    value.Int64(-89),
			want: value.False,
		},

		"UInt64 5.0f64 === 5u64": {
			a:    value.Float64(5),
			b:    value.UInt64(5),
			want: value.False,
		},
		"UInt64 -5.0f64 === 5u64": {
			a:    value.Float64(-5),
			b:    value.UInt64(5),
			want: value.False,
		},
		"UInt64 5.5f64 === 5u64": {
			a:    value.Float64(5.5),
			b:    value.UInt64(5),
			want: value.False,
		},
		"UInt64 NaN === 0u64": {
			a:    value.Float64NaN(),
			b:    value.UInt64(0),
			want: value.False,
		},
		"UInt64 +Inf === 69u64": {
			a:    value.Float64Inf(),
			b:    value.UInt64(69),
			want: value.False,
		},
		"UInt64 -Inf === 89u64": {
			a:    value.Float64NegInf(),
			b:    value.UInt64(89),
			want: value.False,
		},

		"Int32 5.0f64 === 5i32": {
			a:    value.Float64(5),
			b:    value.Int32(5),
			want: value.False,
		},
		"Int32 5.5f64 === 5i32": {
			a:    value.Float64(5.5),
			b:    value.Int32(5),
			want: value.False,
		},
		"Int32 NaN === 0i32": {
			a:    value.Float64NaN(),
			b:    value.Int32(0),
			want: value.False,
		},
		"Int32 +Inf === 69i32": {
			a:    value.Float64Inf(),
			b:    value.Int32(69),
			want: value.False,
		},
		"Int32 -Inf === -89i32": {
			a:    value.Float64NegInf(),
			b:    value.Int32(-89),
			want: value.False,
		},

		"UInt32 5.0f64 === 5u32": {
			a:    value.Float64(5),
			b:    value.UInt32(5),
			want: value.False,
		},
		"UInt32 -5.0f64 === 5u32": {
			a:    value.Float64(-5),
			b:    value.UInt32(5),
			want: value.False,
		},
		"UInt32 5.5f64 === 5u32": {
			a:    value.Float64(5.5),
			b:    value.UInt32(5),
			want: value.False,
		},
		"UInt32 NaN === 0u32": {
			a:    value.Float64NaN(),
			b:    value.UInt32(0),
			want: value.False,
		},
		"UInt32 +Inf === 69u32": {
			a:    value.Float64Inf(),
			b:    value.UInt32(69),
			want: value.False,
		},
		"UInt32 -Inf === 89u32": {
			a:    value.Float64NegInf(),
			b:    value.UInt32(89),
			want: value.False,
		},

		"Int16 5.0f64 === 5i16": {
			a:    value.Float64(5),
			b:    value.Int16(5),
			want: value.False,
		},
		"Int16 5.5f64 === 5i16": {
			a:    value.Float64(5.5),
			b:    value.Int16(5),
			want: value.False,
		},
		"Int16 NaN === 0i16": {
			a:    value.Float64NaN(),
			b:    value.Int16(0),
			want: value.False,
		},
		"Int16 +Inf === 69i16": {
			a:    value.Float64Inf(),
			b:    value.Int16(69),
			want: value.False,
		},
		"Int16 -Inf === -89i16": {
			a:    value.Float64NegInf(),
			b:    value.Int16(-89),
			want: value.False,
		},

		"UInt16 5.0f64 === 5u16": {
			a:    value.Float64(5),
			b:    value.UInt16(5),
			want: value.False,
		},
		"UInt16 -5.0f64 === 5u16": {
			a:    value.Float64(-5),
			b:    value.UInt16(5),
			want: value.False,
		},
		"UInt16 5.5f64 === 5u16": {
			a:    value.Float64(5.5),
			b:    value.UInt16(5),
			want: value.False,
		},
		"UInt16 NaN === 0u16": {
			a:    value.Float64NaN(),
			b:    value.UInt16(0),
			want: value.False,
		},
		"UInt16 +Inf === 69u16": {
			a:    value.Float64Inf(),
			b:    value.UInt16(69),
			want: value.False,
		},
		"UInt16 -Inf === 89u16": {
			a:    value.Float64NegInf(),
			b:    value.UInt16(89),
			want: value.False,
		},

		"Int8 5.0f64 === 5i8": {
			a:    value.Float64(5),
			b:    value.Int8(5),
			want: value.False,
		},
		"Int8 5.5f64 === 5i8": {
			a:    value.Float64(5.5),
			b:    value.Int8(5),
			want: value.False,
		},
		"Int8 NaN === 0i8": {
			a:    value.Float64NaN(),
			b:    value.Int8(0),
			want: value.False,
		},
		"Int8 +Inf === 69i8": {
			a:    value.Float64Inf(),
			b:    value.Int8(69),
			want: value.False,
		},
		"Int8 -Inf === -89i8": {
			a:    value.Float64NegInf(),
			b:    value.Int8(-89),
			want: value.False,
		},

		"UInt8 5.0f64 === 5u8": {
			a:    value.Float64(5),
			b:    value.UInt8(5),
			want: value.False,
		},
		"UInt8 -5.0f64 === 5u8": {
			a:    value.Float64(-5),
			b:    value.UInt8(5),
			want: value.False,
		},
		"UInt8 5.5f64 === 5u8": {
			a:    value.Float64(5.5),
			b:    value.UInt8(5),
			want: value.False,
		},
		"UInt8 NaN === 0u8": {
			a:    value.Float64NaN(),
			b:    value.UInt8(0),
			want: value.False,
		},
		"UInt8 +Inf === 69u8": {
			a:    value.Float64Inf(),
			b:    value.UInt8(69),
			want: value.False,
		},
		"UInt8 -Inf === 89u8": {
			a:    value.Float64NegInf(),
			b:    value.UInt8(89),
			want: value.False,
		},

		"value.Float64 21.9f64 === 21.9f64": {
			a:    value.Float64(21.9),
			b:    value.Float64(21.9),
			want: value.True,
		},
		"value.Float64 21.9f64 === 38.0f64": {
			a:    value.Float64(21.9),
			b:    value.Float64(38),
			want: value.False,
		},
		"value.Float64 NaN === NaN": {
			a:    value.Float64NaN(),
			b:    value.Float64NaN(),
			want: value.False,
		},
		"value.Float64 +Inf === +Inf": {
			a:    value.Float64Inf(),
			b:    value.Float64Inf(),
			want: value.True,
		},
		"value.Float64 -Inf === -Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float64NegInf(),
			want: value.True,
		},
		"value.Float64 +Inf === -Inf": {
			a:    value.Float64Inf(),
			b:    value.Float64NegInf(),
			want: value.False,
		},
		"value.Float64 -Inf === +Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float64Inf(),
			want: value.False,
		},
		"value.Float64 8.5f64 === +Inf": {
			a:    value.Float64(8.5),
			b:    value.Float64Inf(),
			want: value.False,
		},
		"value.Float64 +Inf === 98.0f64": {
			a:    value.Float64Inf(),
			b:    value.Float64(98),
			want: value.False,
		},

		"Float32 21.0f64 === 21.0f32": {
			a:    value.Float64(21),
			b:    value.Float32(21),
			want: value.False,
		},
		"Float32 21.9f64 === 38.0f32": {
			a:    value.Float64(21.9),
			b:    value.Float32(38),
			want: value.False,
		},
		"Float32 NaN === NaN": {
			a:    value.Float64NaN(),
			b:    value.Float32NaN(),
			want: value.False,
		},
		"Float32 +Inf === +Inf": {
			a:    value.Float64Inf(),
			b:    value.Float32Inf(),
			want: value.False,
		},
		"Float32 -Inf === -Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float32NegInf(),
			want: value.False,
		},
		"Float32 +Inf === -Inf": {
			a:    value.Float64Inf(),
			b:    value.Float32NegInf(),
			want: value.False,
		},
		"Float32 -Inf === +Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float32Inf(),
			want: value.False,
		},
		"Float32 8.5f64 === +Inf": {
			a:    value.Float64(8.5),
			b:    value.Float32Inf(),
			want: value.False,
		},
		"Float32 +Inf === 98.0f32": {
			a:    value.Float64Inf(),
			b:    value.Float32(98),
			want: value.False,
		},

		"SmallInt 16.0f64 === 16": {
			a:    value.Float64(16),
			b:    value.SmallInt(16),
			want: value.False,
		},
		"SmallInt 16.5f64 === 16": {
			a:    value.Float64(16.5),
			b:    value.SmallInt(16),
			want: value.False,
		},
		"SmallInt NaN === 0": {
			a:    value.Float64NaN(),
			b:    value.SmallInt(0),
			want: value.False,
		},
		"SmallInt +Inf === 69": {
			a:    value.Float64Inf(),
			b:    value.SmallInt(69),
			want: value.False,
		},
		"SmallInt -Inf === -89": {
			a:    value.Float64NegInf(),
			b:    value.SmallInt(-89),
			want: value.False,
		},

		"BigInt 16.0f64 === 16": {
			a:    value.Float64(16),
			b:    value.NewBigInt(16),
			want: value.False,
		},
		"BigInt 16.5f64 === 16": {
			a:    value.Float64(16.5),
			b:    value.NewBigInt(16),
			want: value.False,
		},
		"BigInt NaN === 0": {
			a:    value.Float64NaN(),
			b:    value.NewBigInt(0),
			want: value.False,
		},
		"BigInt +Inf === 69": {
			a:    value.Float64Inf(),
			b:    value.NewBigInt(69),
			want: value.False,
		},
		"BigInt -Inf === -89": {
			a:    value.Float64NegInf(),
			b:    value.NewBigInt(-89),
			want: value.False,
		},

		"Float 21.9f64 === 21.9": {
			a:    value.Float64(21.9),
			b:    value.Float(21.9),
			want: value.False,
		},
		"Float 21.9f64 === 38.0": {
			a:    value.Float64(21.9),
			b:    value.Float(38),
			want: value.False,
		},
		"Float NaN === NaN": {
			a:    value.Float64NaN(),
			b:    value.FloatNaN(),
			want: value.False,
		},
		"Float +Inf === +Inf": {
			a:    value.Float64Inf(),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float -Inf === -Inf": {
			a:    value.Float64NegInf(),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float +Inf === -Inf": {
			a:    value.Float64Inf(),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float -Inf === +Inf": {
			a:    value.Float64NegInf(),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float 8.5f64 === +Inf": {
			a:    value.Float64(8.5),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float +Inf === 98.0": {
			a:    value.Float64Inf(),
			b:    value.Float(98),
			want: value.False,
		},

		"BigFloat 21.9f64 === 21.9bf": {
			a:    value.Float64(21.9),
			b:    value.NewBigFloat(21.9),
			want: value.False,
		},
		"BigFloat 21.9f64 === 38.0bf": {
			a:    value.Float64(21.9),
			b:    value.NewBigFloat(38),
			want: value.False,
		},
		"BigFloat NaN === NaN": {
			a:    value.Float64NaN(),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
		"BigFloat +Inf === +Inf": {
			a:    value.Float64Inf(),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat -Inf === -Inf": {
			a:    value.Float64NegInf(),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
		"BigFloat +Inf === -Inf": {
			a:    value.Float64Inf(),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
		"BigFloat -Inf === +Inf": {
			a:    value.Float64NegInf(),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat 8.5f64 === +Inf": {
			a:    value.Float64(8.5),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat +Inf === 98.0bf": {
			a:    value.Float64Inf(),
			b:    value.NewBigFloat(98),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.StrictNumericStrictEqual(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictFloat_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Bool
		err  *value.Error
	}{
		"String 5.5f64 == '5.5'": {
			a:    value.Float64(5.5),
			b:    value.String("5.5"),
			want: value.False,
		},

		"Int64 5.0f64 == 5i64": {
			a:    value.Float64(5),
			b:    value.Int64(5),
			want: value.True,
		},
		"Int64 5.5f64 == 5i64": {
			a:    value.Float64(5.5),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int64 NaN == 0i64": {
			a:    value.Float64NaN(),
			b:    value.Int64(0),
			want: value.False,
		},
		"Int64 +Inf == 69i64": {
			a:    value.Float64Inf(),
			b:    value.Int64(69),
			want: value.False,
		},
		"Int64 -Inf == -89i64": {
			a:    value.Float64NegInf(),
			b:    value.Int64(-89),
			want: value.False,
		},

		"UInt64 5.0f64 == 5u64": {
			a:    value.Float64(5),
			b:    value.UInt64(5),
			want: value.True,
		},
		"UInt64 -5.0f64 == 5u64": {
			a:    value.Float64(-5),
			b:    value.UInt64(5),
			want: value.False,
		},
		"UInt64 5.5f64 == 5u64": {
			a:    value.Float64(5.5),
			b:    value.UInt64(5),
			want: value.False,
		},
		"UInt64 NaN == 0u64": {
			a:    value.Float64NaN(),
			b:    value.UInt64(0),
			want: value.False,
		},
		"UInt64 +Inf == 69u64": {
			a:    value.Float64Inf(),
			b:    value.UInt64(69),
			want: value.False,
		},
		"UInt64 -Inf == 89u64": {
			a:    value.Float64NegInf(),
			b:    value.UInt64(89),
			want: value.False,
		},

		"Int32 5.0f64 == 5i32": {
			a:    value.Float64(5),
			b:    value.Int32(5),
			want: value.True,
		},
		"Int32 5.5f64 == 5i32": {
			a:    value.Float64(5.5),
			b:    value.Int32(5),
			want: value.False,
		},
		"Int32 NaN == 0i32": {
			a:    value.Float64NaN(),
			b:    value.Int32(0),
			want: value.False,
		},
		"Int32 +Inf == 69i32": {
			a:    value.Float64Inf(),
			b:    value.Int32(69),
			want: value.False,
		},
		"Int32 -Inf == -89i32": {
			a:    value.Float64NegInf(),
			b:    value.Int32(-89),
			want: value.False,
		},

		"UInt32 5.0f64 == 5u32": {
			a:    value.Float64(5),
			b:    value.UInt32(5),
			want: value.True,
		},
		"UInt32 -5.0f64 == 5u32": {
			a:    value.Float64(-5),
			b:    value.UInt32(5),
			want: value.False,
		},
		"UInt32 5.5f64 == 5u32": {
			a:    value.Float64(5.5),
			b:    value.UInt32(5),
			want: value.False,
		},
		"UInt32 NaN == 0u32": {
			a:    value.Float64NaN(),
			b:    value.UInt32(0),
			want: value.False,
		},
		"UInt32 +Inf == 69u32": {
			a:    value.Float64Inf(),
			b:    value.UInt32(69),
			want: value.False,
		},
		"UInt32 -Inf == 89u32": {
			a:    value.Float64NegInf(),
			b:    value.UInt32(89),
			want: value.False,
		},

		"Int16 5.0f64 == 5i16": {
			a:    value.Float64(5),
			b:    value.Int16(5),
			want: value.True,
		},
		"Int16 5.5f64 == 5i16": {
			a:    value.Float64(5.5),
			b:    value.Int16(5),
			want: value.False,
		},
		"Int16 NaN == 0i16": {
			a:    value.Float64NaN(),
			b:    value.Int16(0),
			want: value.False,
		},
		"Int16 +Inf == 69i16": {
			a:    value.Float64Inf(),
			b:    value.Int16(69),
			want: value.False,
		},
		"Int16 -Inf == -89i16": {
			a:    value.Float64NegInf(),
			b:    value.Int16(-89),
			want: value.False,
		},

		"UInt16 5.0f64 == 5u16": {
			a:    value.Float64(5),
			b:    value.UInt16(5),
			want: value.True,
		},
		"UInt16 -5.0f64 == 5u16": {
			a:    value.Float64(-5),
			b:    value.UInt16(5),
			want: value.False,
		},
		"UInt16 5.5f64 == 5u16": {
			a:    value.Float64(5.5),
			b:    value.UInt16(5),
			want: value.False,
		},
		"UInt16 NaN == 0u16": {
			a:    value.Float64NaN(),
			b:    value.UInt16(0),
			want: value.False,
		},
		"UInt16 +Inf == 69u16": {
			a:    value.Float64Inf(),
			b:    value.UInt16(69),
			want: value.False,
		},
		"UInt16 -Inf == 89u16": {
			a:    value.Float64NegInf(),
			b:    value.UInt16(89),
			want: value.False,
		},

		"Int8 5.0f64 == 5i8": {
			a:    value.Float64(5),
			b:    value.Int8(5),
			want: value.True,
		},
		"Int8 5.5f64 == 5i8": {
			a:    value.Float64(5.5),
			b:    value.Int8(5),
			want: value.False,
		},
		"Int8 NaN == 0i8": {
			a:    value.Float64NaN(),
			b:    value.Int8(0),
			want: value.False,
		},
		"Int8 +Inf == 69i8": {
			a:    value.Float64Inf(),
			b:    value.Int8(69),
			want: value.False,
		},
		"Int8 -Inf == -89i8": {
			a:    value.Float64NegInf(),
			b:    value.Int8(-89),
			want: value.False,
		},

		"UInt8 5.0f64 == 5u8": {
			a:    value.Float64(5),
			b:    value.UInt8(5),
			want: value.True,
		},
		"UInt8 -5.0f64 == 5u8": {
			a:    value.Float64(-5),
			b:    value.UInt8(5),
			want: value.False,
		},
		"UInt8 5.5f64 == 5u8": {
			a:    value.Float64(5.5),
			b:    value.UInt8(5),
			want: value.False,
		},
		"UInt8 NaN == 0u8": {
			a:    value.Float64NaN(),
			b:    value.UInt8(0),
			want: value.False,
		},
		"UInt8 +Inf == 69u8": {
			a:    value.Float64Inf(),
			b:    value.UInt8(69),
			want: value.False,
		},
		"UInt8 -Inf == 89u8": {
			a:    value.Float64NegInf(),
			b:    value.UInt8(89),
			want: value.False,
		},

		"value.Float64 21.9f64 == 21.9f64": {
			a:    value.Float64(21.9),
			b:    value.Float64(21.9),
			want: value.True,
		},
		"value.Float64 21.9f64 == 38.0f64": {
			a:    value.Float64(21.9),
			b:    value.Float64(38),
			want: value.False,
		},
		"value.Float64 NaN == NaN": {
			a:    value.Float64NaN(),
			b:    value.Float64NaN(),
			want: value.False,
		},
		"value.Float64 +Inf == +Inf": {
			a:    value.Float64Inf(),
			b:    value.Float64Inf(),
			want: value.True,
		},
		"value.Float64 -Inf == -Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float64NegInf(),
			want: value.True,
		},
		"value.Float64 +Inf == -Inf": {
			a:    value.Float64Inf(),
			b:    value.Float64NegInf(),
			want: value.False,
		},
		"value.Float64 -Inf == +Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float64Inf(),
			want: value.False,
		},
		"value.Float64 8.5f64 == +Inf": {
			a:    value.Float64(8.5),
			b:    value.Float64Inf(),
			want: value.False,
		},
		"value.Float64 +Inf == 98.0f64": {
			a:    value.Float64Inf(),
			b:    value.Float64(98),
			want: value.False,
		},

		"Float32 21.0f64 == 21.0f32": {
			a:    value.Float64(21),
			b:    value.Float32(21),
			want: value.True,
		},
		"Float32 21.9f64 == 38.0f32": {
			a:    value.Float64(21.9),
			b:    value.Float32(38),
			want: value.False,
		},
		"Float32 NaN == NaN": {
			a:    value.Float64NaN(),
			b:    value.Float32NaN(),
			want: value.False,
		},
		"Float32 +Inf == +Inf": {
			a:    value.Float64Inf(),
			b:    value.Float32Inf(),
			want: value.True,
		},
		"Float32 -Inf == -Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float32NegInf(),
			want: value.True,
		},
		"Float32 +Inf == -Inf": {
			a:    value.Float64Inf(),
			b:    value.Float32NegInf(),
			want: value.False,
		},
		"Float32 -Inf == +Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float32Inf(),
			want: value.False,
		},
		"Float32 8.5f64 == +Inf": {
			a:    value.Float64(8.5),
			b:    value.Float32Inf(),
			want: value.False,
		},
		"Float32 +Inf == 98.0f32": {
			a:    value.Float64Inf(),
			b:    value.Float32(98),
			want: value.False,
		},

		"SmallInt 16.0f64 == 16": {
			a:    value.Float64(16),
			b:    value.SmallInt(16),
			want: value.True,
		},
		"SmallInt 16.5f64 == 16": {
			a:    value.Float64(16.5),
			b:    value.SmallInt(16),
			want: value.False,
		},
		"SmallInt NaN == 0": {
			a:    value.Float64NaN(),
			b:    value.SmallInt(0),
			want: value.False,
		},
		"SmallInt +Inf == 69": {
			a:    value.Float64Inf(),
			b:    value.SmallInt(69),
			want: value.False,
		},
		"SmallInt -Inf == -89": {
			a:    value.Float64NegInf(),
			b:    value.SmallInt(-89),
			want: value.False,
		},

		"BigInt 16.0f64 == 16": {
			a:    value.Float64(16),
			b:    value.NewBigInt(16),
			want: value.True,
		},
		"BigInt 16.5f64 == 16": {
			a:    value.Float64(16.5),
			b:    value.NewBigInt(16),
			want: value.False,
		},
		"BigInt NaN == 0": {
			a:    value.Float64NaN(),
			b:    value.NewBigInt(0),
			want: value.False,
		},
		"BigInt +Inf == 69": {
			a:    value.Float64Inf(),
			b:    value.NewBigInt(69),
			want: value.False,
		},
		"BigInt -Inf == -89": {
			a:    value.Float64NegInf(),
			b:    value.NewBigInt(-89),
			want: value.False,
		},

		"Float 21.9f64 == 21.9": {
			a:    value.Float64(21.9),
			b:    value.Float(21.9),
			want: value.True,
		},
		"Float 21.9f64 == 38.0": {
			a:    value.Float64(21.9),
			b:    value.Float(38),
			want: value.False,
		},
		"Float NaN == NaN": {
			a:    value.Float64NaN(),
			b:    value.FloatNaN(),
			want: value.False,
		},
		"Float +Inf == +Inf": {
			a:    value.Float64Inf(),
			b:    value.FloatInf(),
			want: value.True,
		},
		"Float -Inf == -Inf": {
			a:    value.Float64NegInf(),
			b:    value.FloatNegInf(),
			want: value.True,
		},
		"Float +Inf == -Inf": {
			a:    value.Float64Inf(),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float -Inf == +Inf": {
			a:    value.Float64NegInf(),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float 8.5f64 == +Inf": {
			a:    value.Float64(8.5),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float +Inf == 98.0": {
			a:    value.Float64Inf(),
			b:    value.Float(98),
			want: value.False,
		},

		"BigFloat 21.9f64 == 21.9bf": {
			a:    value.Float64(21.9),
			b:    value.NewBigFloat(21.9),
			want: value.True,
		},
		"BigFloat 21.9f64 == 38.0bf": {
			a:    value.Float64(21.9),
			b:    value.NewBigFloat(38),
			want: value.False,
		},
		"BigFloat NaN == NaN": {
			a:    value.Float64NaN(),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
		"BigFloat +Inf == +Inf": {
			a:    value.Float64Inf(),
			b:    value.BigFloatInf(),
			want: value.True,
		},
		"BigFloat -Inf == -Inf": {
			a:    value.Float64NegInf(),
			b:    value.BigFloatNegInf(),
			want: value.True,
		},
		"BigFloat +Inf == -Inf": {
			a:    value.Float64Inf(),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
		"BigFloat -Inf == +Inf": {
			a:    value.Float64NegInf(),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat 8.5f64 == +Inf": {
			a:    value.Float64(8.5),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat +Inf == 98.0bf": {
			a:    value.Float64Inf(),
			b:    value.NewBigFloat(98),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.StrictFloatEqual(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictSignedInt_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Bool
		err  *value.Error
	}{
		"String 5i64 == '5'": {
			a:    value.Int64(5),
			b:    value.String("5.5"),
			want: value.False,
		},

		"Int64 5i64 == 5i64": {
			a:    value.Int64(5),
			b:    value.Int64(5),
			want: value.True,
		},
		"Int64 4i64 == 5i64": {
			a:    value.Int64(4),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int64 5i64 == -5i64": {
			a:    value.Int64(5),
			b:    value.Int64(-5),
			want: value.False,
		},

		"UInt64 5i64 == 5u64": {
			a:    value.Int64(5),
			b:    value.UInt64(5),
			want: value.True,
		},
		"UInt64 -5i64 == 5u64": {
			a:    value.Int64(-5),
			b:    value.UInt64(5),
			want: value.False,
		},

		"Int32 5i64 == 5i32": {
			a:    value.Int64(5),
			b:    value.Int32(5),
			want: value.True,
		},
		"Int32 -5i64 == 5i32": {
			a:    value.Int64(-5),
			b:    value.Int32(5),
			want: value.False,
		},
		"Int32 5i64 == -5i32": {
			a:    value.Int64(5),
			b:    value.Int32(-5),
			want: value.False,
		},
		"Int32 -5i64 == -5i32": {
			a:    value.Int64(-5),
			b:    value.Int32(-5),
			want: value.True,
		},

		"UInt32 5i64 == 5u32": {
			a:    value.Int64(5),
			b:    value.UInt32(5),
			want: value.True,
		},
		"UInt32 -5i64 == 5u32": {
			a:    value.Int64(-5),
			b:    value.UInt32(5),
			want: value.False,
		},
		"UInt32 4i64 == 5u32": {
			a:    value.Int64(4),
			b:    value.UInt32(5),
			want: value.False,
		},

		"Int16 5i64 == 5i16": {
			a:    value.Int64(5),
			b:    value.Int16(5),
			want: value.True,
		},
		"Int16 -5i64 == 5i16": {
			a:    value.Int64(-5),
			b:    value.Int16(5),
			want: value.False,
		},
		"Int16 5i64 == -5i16": {
			a:    value.Int64(5),
			b:    value.Int16(-5),
			want: value.False,
		},
		"Int16 -5i64 == -5i16": {
			a:    value.Int64(-5),
			b:    value.Int16(-5),
			want: value.True,
		},
		"Int16 4i64 == 5i16": {
			a:    value.Int64(4),
			b:    value.Int16(5),
			want: value.False,
		},

		"UInt16 5i64 == 5u16": {
			a:    value.Int64(5),
			b:    value.UInt16(5),
			want: value.True,
		},
		"UInt16 -5i64 == 5u16": {
			a:    value.Int64(-5),
			b:    value.UInt16(5),
			want: value.False,
		},
		"UInt16 4i64 == 5u16": {
			a:    value.Int64(4),
			b:    value.UInt16(5),
			want: value.False,
		},

		"Int8 5i64 == 5i8": {
			a:    value.Int64(5),
			b:    value.Int8(5),
			want: value.True,
		},
		"Int8 4i64 == 5i8": {
			a:    value.Int64(4),
			b:    value.Int8(5),
			want: value.False,
		},
		"Int8 -5i64 == 5i8": {
			a:    value.Int64(-5),
			b:    value.Int8(5),
			want: value.False,
		},
		"Int8 5i64 == -5i8": {
			a:    value.Int64(5),
			b:    value.Int8(-5),
			want: value.False,
		},
		"Int8 -5i64 == -5i8": {
			a:    value.Int64(-5),
			b:    value.Int8(-5),
			want: value.True,
		},

		"UInt8 5i64 == 5u8": {
			a:    value.Int64(5),
			b:    value.UInt8(5),
			want: value.True,
		},
		"UInt8 4i64 == 5u8": {
			a:    value.Int64(4),
			b:    value.UInt8(5),
			want: value.False,
		},
		"UInt8 -5i64 == 5u8": {
			a:    value.Int64(-5),
			b:    value.UInt8(5),
			want: value.False,
		},

		"value.Float64 21i64 == 21.0f64": {
			a:    value.Int64(21),
			b:    value.Float64(21),
			want: value.True,
		},
		"value.Float64 21i64 == 21.5f64": {
			a:    value.Int64(21),
			b:    value.Float64(21.5),
			want: value.False,
		},
		"value.Float64 21i64 == 38.0f64": {
			a:    value.Int64(21),
			b:    value.Float64(38),
			want: value.False,
		},
		"value.Float64 0i64 == NaN": {
			a:    value.Int64(0),
			b:    value.Float64NaN(),
			want: value.False,
		},
		"value.Float64 8i64 == +Inf": {
			a:    value.Int64(8),
			b:    value.Float64Inf(),
			want: value.False,
		},
		"value.Float64 8i64 == -Inf": {
			a:    value.Int64(8),
			b:    value.Float64NegInf(),
			want: value.False,
		},

		"Float32 21i64 == 21.0f32": {
			a:    value.Int64(21),
			b:    value.Float32(21),
			want: value.True,
		},
		"Float32 21i64 == 21.5f32": {
			a:    value.Int64(21),
			b:    value.Float32(21.5),
			want: value.False,
		},
		"Float32 21i64 == 38.0f32": {
			a:    value.Int64(21),
			b:    value.Float32(38),
			want: value.False,
		},
		"Float32 0i64 == NaN": {
			a:    value.Int64(0),
			b:    value.Float32NaN(),
			want: value.False,
		},
		"Float32 8i64 == +Inf": {
			a:    value.Int64(8),
			b:    value.Float32Inf(),
			want: value.False,
		},
		"Float32 8i64 == -Inf": {
			a:    value.Int64(8),
			b:    value.Float32NegInf(),
			want: value.False,
		},

		"SmallInt 16i64 == 16": {
			a:    value.Int64(16),
			b:    value.SmallInt(16),
			want: value.True,
		},
		"SmallInt 97i64 == -97": {
			a:    value.Int64(97),
			b:    value.SmallInt(-97),
			want: value.False,
		},
		"SmallInt -6i64 == 6": {
			a:    value.Int64(-6),
			b:    value.SmallInt(6),
			want: value.False,
		},
		"SmallInt -120i64 == -120": {
			a:    value.Int64(-120),
			b:    value.SmallInt(-120),
			want: value.True,
		},

		"BigInt 16i64 == 16": {
			a:    value.Int64(16),
			b:    value.NewBigInt(16),
			want: value.True,
		},
		"BigInt 97i64 == -97": {
			a:    value.Int64(97),
			b:    value.NewBigInt(-97),
			want: value.False,
		},
		"BigInt -6i64 == 6": {
			a:    value.Int64(-6),
			b:    value.NewBigInt(6),
			want: value.False,
		},
		"BigInt -120i64 == -120": {
			a:    value.Int64(-120),
			b:    value.NewBigInt(-120),
			want: value.True,
		},

		"Float 21i64 == 21.0": {
			a:    value.Int64(21),
			b:    value.Float(21),
			want: value.True,
		},
		"Float 21i64 == 21.5": {
			a:    value.Int64(21),
			b:    value.Float(21.5),
			want: value.False,
		},
		"Float 21i64 == 38.0": {
			a:    value.Int64(21),
			b:    value.Float(38),
			want: value.False,
		},
		"Float 0i64 == NaN": {
			a:    value.Int64(0),
			b:    value.FloatNaN(),
			want: value.False,
		},
		"Float 8i64 == +Inf": {
			a:    value.Int64(8),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float 8i64 == -Inf": {
			a:    value.Int64(8),
			b:    value.FloatNegInf(),
			want: value.False,
		},

		"BigFloat 21i64 == 21.0bf": {
			a:    value.Int64(21),
			b:    value.NewBigFloat(21),
			want: value.True,
		},
		"BigFloat 21i64 == 21.5bf": {
			a:    value.Int64(21),
			b:    value.NewBigFloat(21.5),
			want: value.False,
		},
		"BigFloat 21i64 == 38.0bf": {
			a:    value.Int64(21),
			b:    value.NewBigFloat(38),
			want: value.False,
		},
		"BigFloat 0i64 == NaN": {
			a:    value.Int64(0),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
		"BigFloat 8i64 == +Inf": {
			a:    value.Int64(8),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat 8i64 == -Inf": {
			a:    value.Int64(8),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.StrictSignedIntEqual(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictUnsignedInt_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.Bool
		err  *value.Error
	}{
		"String 5u64 == '5'": {
			a:    value.UInt64(5),
			b:    value.String("5.5"),
			want: value.False,
		},

		"Int64 5u64 == 5i64": {
			a:    value.UInt64(5),
			b:    value.Int64(5),
			want: value.True,
		},
		"Int64 4u64 == 5i64": {
			a:    value.UInt64(4),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int64 5u64 == -5i64": {
			a:    value.UInt64(5),
			b:    value.Int64(-5),
			want: value.False,
		},

		"UInt64 5u64 == 5u64": {
			a:    value.UInt64(5),
			b:    value.UInt64(5),
			want: value.True,
		},
		"UInt64 5u64 == 7u64": {
			a:    value.UInt64(5),
			b:    value.UInt64(7),
			want: value.False,
		},

		"Int32 5u64 == 5i32": {
			a:    value.UInt64(5),
			b:    value.Int32(5),
			want: value.True,
		},
		"Int32 5u64 == 7i32": {
			a:    value.UInt64(5),
			b:    value.Int32(7),
			want: value.False,
		},
		"Int32 5u64 == -5i32": {
			a:    value.UInt64(5),
			b:    value.Int32(-5),
			want: value.False,
		},

		"UInt32 5u64 == 5u32": {
			a:    value.UInt64(5),
			b:    value.UInt32(5),
			want: value.True,
		},
		"UInt32 4u64 == 5u32": {
			a:    value.UInt64(4),
			b:    value.UInt32(5),
			want: value.False,
		},

		"Int16 5u64 == 5i16": {
			a:    value.UInt64(5),
			b:    value.Int16(5),
			want: value.True,
		},
		"Int16 5u64 == -5i16": {
			a:    value.UInt64(5),
			b:    value.Int16(-5),
			want: value.False,
		},
		"Int16 4u64 == 5i16": {
			a:    value.UInt64(4),
			b:    value.Int16(5),
			want: value.False,
		},

		"UInt16 5u64 == 5u16": {
			a:    value.UInt64(5),
			b:    value.UInt16(5),
			want: value.True,
		},
		"UInt16 4u64 == 5u16": {
			a:    value.UInt64(4),
			b:    value.UInt16(5),
			want: value.False,
		},

		"Int8 5u64 == 5i8": {
			a:    value.UInt64(5),
			b:    value.Int8(5),
			want: value.True,
		},
		"Int8 4u64 == 5i8": {
			a:    value.UInt64(4),
			b:    value.Int8(5),
			want: value.False,
		},
		"Int8 5u64 == -5i8": {
			a:    value.UInt64(5),
			b:    value.Int8(-5),
			want: value.False,
		},

		"UInt8 5u64 == 5u8": {
			a:    value.UInt64(5),
			b:    value.UInt8(5),
			want: value.True,
		},
		"UInt8 4u64 == 5u8": {
			a:    value.UInt64(4),
			b:    value.UInt8(5),
			want: value.False,
		},

		"value.Float64 21u64 == 21.0f64": {
			a:    value.UInt64(21),
			b:    value.Float64(21),
			want: value.True,
		},
		"value.Float64 21u64 == 21.5f64": {
			a:    value.UInt64(21),
			b:    value.Float64(21.5),
			want: value.False,
		},
		"value.Float64 21u64 == 38.0f64": {
			a:    value.UInt64(21),
			b:    value.Float64(38),
			want: value.False,
		},
		"value.Float64 0u64 == NaN": {
			a:    value.UInt64(0),
			b:    value.Float64NaN(),
			want: value.False,
		},
		"value.Float64 8u64 == +Inf": {
			a:    value.UInt64(8),
			b:    value.Float64Inf(),
			want: value.False,
		},
		"value.Float64 8u64 == -Inf": {
			a:    value.UInt64(8),
			b:    value.Float64NegInf(),
			want: value.False,
		},

		"Float32 21u64 == 21.0f32": {
			a:    value.UInt64(21),
			b:    value.Float32(21),
			want: value.True,
		},
		"Float32 21u64 == 21.5f32": {
			a:    value.UInt64(21),
			b:    value.Float32(21.5),
			want: value.False,
		},
		"Float32 21u64 == 38.0f32": {
			a:    value.UInt64(21),
			b:    value.Float32(38),
			want: value.False,
		},
		"Float32 0u64 == NaN": {
			a:    value.UInt64(0),
			b:    value.Float32NaN(),
			want: value.False,
		},
		"Float32 8u64 == +Inf": {
			a:    value.UInt64(8),
			b:    value.Float32Inf(),
			want: value.False,
		},
		"Float32 8u64 == -Inf": {
			a:    value.UInt64(8),
			b:    value.Float32NegInf(),
			want: value.False,
		},

		"SmallInt 16u64 == 16": {
			a:    value.UInt64(16),
			b:    value.SmallInt(16),
			want: value.True,
		},
		"SmallInt 97u64 == -97": {
			a:    value.UInt64(97),
			b:    value.SmallInt(-97),
			want: value.False,
		},

		"BigInt 16u64 == 16": {
			a:    value.UInt64(16),
			b:    value.NewBigInt(16),
			want: value.True,
		},
		"BigInt 97u64 == -97": {
			a:    value.UInt64(97),
			b:    value.NewBigInt(-97),
			want: value.False,
		},

		"Float 21u64 == 21.0": {
			a:    value.UInt64(21),
			b:    value.Float(21),
			want: value.True,
		},
		"Float 21u64 == 21.5": {
			a:    value.UInt64(21),
			b:    value.Float(21.5),
			want: value.False,
		},
		"Float 21u64 == 38.0": {
			a:    value.UInt64(21),
			b:    value.Float(38),
			want: value.False,
		},
		"Float 0u64 == NaN": {
			a:    value.UInt64(0),
			b:    value.FloatNaN(),
			want: value.False,
		},
		"Float 8u64 == +Inf": {
			a:    value.UInt64(8),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float 8u64 == -Inf": {
			a:    value.UInt64(8),
			b:    value.FloatNegInf(),
			want: value.False,
		},

		"BigFloat 21u64 == 21.0bf": {
			a:    value.UInt64(21),
			b:    value.NewBigFloat(21),
			want: value.True,
		},
		"BigFloat 21u64 == 21.5bf": {
			a:    value.UInt64(21),
			b:    value.NewBigFloat(21.5),
			want: value.False,
		},
		"BigFloat 21u64 == 38.0bf": {
			a:    value.UInt64(21),
			b:    value.NewBigFloat(38),
			want: value.False,
		},
		"BigFloat 0u64 == NaN": {
			a:    value.UInt64(0),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
		"BigFloat 8u64 == +Inf": {
			a:    value.UInt64(8),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat 8u64 == -Inf": {
			a:    value.UInt64(8),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.StrictUnsignedIntEqual(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictInt_Divide(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"divide by String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"divide Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"divide positive Int64": {
			a:    value.Int64(54),
			b:    value.Int64(2),
			want: value.Int64(27),
		},
		"divide negative Int64": {
			a:    value.Int64(50),
			b:    value.Int64(-2),
			want: value.Int64(-25),
		},
		"divide by zero": {
			a:   value.Int64(50),
			b:   value.Int64(0),
			err: value.NewError(value.ZeroDivisionErrorClass, "can't divide by zero"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictIntDivide(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictFloat_Divide(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Float64
		err  *value.Error
	}{
		"divide by String and return an error": {
			a:   value.Float64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Float64`"),
		},
		"divide Int32 and return an error": {
			a:   value.Float64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Float64`"),
		},
		"divide positive value.Float64": {
			a:    value.Float64(54.5),
			b:    value.Float64(2),
			want: value.Float64(27.25),
		},
		"divide negative value.Float64": {
			a:    value.Float64(50),
			b:    value.Float64(-2),
			want: value.Float64(-25),
		},
		"divide by zero": {
			a:    value.Float64(50),
			b:    value.Float64(0),
			want: value.Float64(math.Inf(1)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictFloatDivide(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumeric_ParseUint(t *testing.T) {
	tests := map[string]struct {
		str     string
		base    int
		bitSize int
		want    uint64
		err     *value.Error
	}{
		"explicit decimal": {
			str:     "123",
			base:    10,
			bitSize: 8,
			want:    123,
		},
		"implicit decimal": {
			str:     "123",
			base:    0,
			bitSize: 8,
			want:    123,
		},
		"implicit decimal with underscores": {
			str:     "2_500",
			base:    0,
			bitSize: 16,
			want:    2500,
		},
		"8bit implicit decimal out of range": {
			str:     "300",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     value.Errorf(value.FormatErrorClass, "value overflows"),
		},
		"explicit hex": {
			str:     "ff",
			base:    16,
			bitSize: 8,
			want:    255,
		},
		"implicit hex": {
			str:     "0xff",
			base:    0,
			bitSize: 8,
			want:    255,
		},
		"implicit hex with underscores": {
			str:     "0x12_34",
			base:    0,
			bitSize: 16,
			want:    4660,
		},
		"8bit implicit hex out of range": {
			str:     "0xfff",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     value.Errorf(value.FormatErrorClass, "value overflows"),
		},
		"explicit duodecimal": {
			str:     "1a",
			base:    12,
			bitSize: 8,
			want:    22,
		},
		"implicit duodecimal": {
			str:     "0d1a",
			base:    0,
			bitSize: 8,
			want:    22,
		},
		"implicit duodecimal with underscores": {
			str:     "0d12_34",
			base:    0,
			bitSize: 16,
			want:    2056,
		},
		"8bit implicit duodecimal out of range": {
			str:     "0d194",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     value.Errorf(value.FormatErrorClass, "value overflows"),
		},
		"explicit octal": {
			str:     "67",
			base:    8,
			bitSize: 8,
			want:    55,
		},
		"implicit octal": {
			str:     "0o67",
			base:    0,
			bitSize: 8,
			want:    55,
		},
		"implicit octal with underscores": {
			str:     "0o12_34",
			base:    0,
			bitSize: 16,
			want:    668,
		},
		"8bit implicit octal out of range": {
			str:     "0o400",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     value.Errorf(value.FormatErrorClass, "value overflows"),
		},
		"explicit quaternary": {
			str:     "33",
			base:    4,
			bitSize: 8,
			want:    15,
		},
		"implicit quaternary": {
			str:     "0q33",
			base:    0,
			bitSize: 8,
			want:    15,
		},
		"implicit quaternary with underscores": {
			str:     "0q12_33",
			base:    0,
			bitSize: 8,
			want:    111,
		},
		"8bit implicit quaternary out of range": {
			str:     "0q10000",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     value.Errorf(value.FormatErrorClass, "value overflows"),
		},
		"explicit binary": {
			str:     "101",
			base:    2,
			bitSize: 8,
			want:    5,
		},
		"implicit binary": {
			str:     "0b101",
			base:    0,
			bitSize: 8,
			want:    5,
		},
		"implicit binary with underscores": {
			str:     "0b100_111",
			base:    0,
			bitSize: 16,
			want:    39,
		},
		"8bit implicit binary out of range": {
			str:     "0b100000000",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     value.Errorf(value.FormatErrorClass, "value overflows"),
		},
		"64bit decimal": {
			str:     "123",
			base:    10,
			bitSize: 64,
			want:    123,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictParseUint(tc.str, tc.base, tc.bitSize)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictInt_RightBitshift(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"shift by String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be used as a bitshift operand"),
		},
		"shift by Float and return an error": {
			a:   value.Int64(5),
			b:   value.Float(3.2),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` can't be used as a bitshift operand"),
		},
		"shift by Int32": {
			a:    value.Int64(234),
			b:    value.Int32(2),
			want: value.Int64(58),
		},
		"shift by UInt8": {
			a:    value.Int64(234),
			b:    value.UInt8(2),
			want: value.Int64(58),
		},
		"shift by SmallInt": {
			a:    value.Int64(234),
			b:    value.SmallInt(2),
			want: value.Int64(58),
		},
		"shift by BigInt": {
			a:    value.Int64(234),
			b:    value.NewBigInt(2),
			want: value.Int64(58),
		},
		"shift by large BigInt": {
			a:    value.Int64(234),
			b:    value.ParseBigIntPanic("9223372036854775808", 10),
			want: value.Int64(0),
		},
		"shift by 10 >> 1": {
			a:    value.Int64(10),
			b:    value.Int64(1),
			want: value.Int64(5),
		},
		"shift by 10 >> 255": {
			a:    value.Int64(10),
			b:    value.Int64(255),
			want: value.Int64(0),
		},
		"shift by 25 >> 2": {
			a:    value.Int64(25),
			b:    value.Int64(2),
			want: value.Int64(6),
		},
		"shift by 25 >> -2": {
			a:    value.Int64(25),
			b:    value.Int64(-2),
			want: value.Int64(100),
		},
		"shift by -6 >> 1": {
			a:    value.Int64(-6),
			b:    value.Int64(1),
			want: value.Int64(-3),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictIntRightBitshift(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictInt_LogicalRightBitshift(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"shift by String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be used as a bitshift operand"),
		},
		"shift by Float and return an error": {
			a:   value.Int64(5),
			b:   value.Float(3.2),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` can't be used as a bitshift operand"),
		},
		"shift by Int32": {
			a:    value.Int64(234),
			b:    value.Int32(2),
			want: value.Int64(58),
		},
		"shift by UInt8": {
			a:    value.Int64(234),
			b:    value.UInt8(2),
			want: value.Int64(58),
		},
		"shift by SmallInt": {
			a:    value.Int64(234),
			b:    value.SmallInt(2),
			want: value.Int64(58),
		},
		"shift by BigInt": {
			a:    value.Int64(234),
			b:    value.NewBigInt(2),
			want: value.Int64(58),
		},
		"shift by large BigInt": {
			a:    value.Int64(234),
			b:    value.ParseBigIntPanic("9223372036854775808", 10),
			want: value.Int64(0),
		},
		"shift by 10 >>> 1": {
			a:    value.Int64(10),
			b:    value.Int64(1),
			want: value.Int64(5),
		},
		"shift by 10 >>> 255": {
			a:    value.Int64(10),
			b:    value.Int64(255),
			want: value.Int64(0),
		},
		"shift by 25 >>> 2": {
			a:    value.Int64(25),
			b:    value.Int64(2),
			want: value.Int64(6),
		},
		"shift by 25 >>> -2": {
			a:    value.Int64(25),
			b:    value.Int64(-2),
			want: value.Int64(100),
		},
		"shift by -6 >>> 1": {
			a:    value.Int64(-6),
			b:    value.Int64(1),
			want: value.Int64(9223372036854775805),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictIntLogicalRightBitshift(tc.a, tc.b, value.LogicalRightShift64)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictInt_LeftBitshift(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"shift by String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be used as a bitshift operand"),
		},
		"shift by Float and return an error": {
			a:   value.Int64(5),
			b:   value.Float(3.2),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` can't be used as a bitshift operand"),
		},
		"shift by Int32": {
			a:    value.Int64(234),
			b:    value.Int32(2),
			want: value.Int64(936),
		},
		"shift by UInt8": {
			a:    value.Int64(234),
			b:    value.UInt8(2),
			want: value.Int64(936),
		},
		"shift by SmallInt": {
			a:    value.Int64(234),
			b:    value.SmallInt(2),
			want: value.Int64(936),
		},
		"shift by BigInt": {
			a:    value.Int64(234),
			b:    value.NewBigInt(2),
			want: value.Int64(936),
		},
		"shift by large BigInt": {
			a:    value.Int64(234),
			b:    value.ParseBigIntPanic("9223372036854775808", 10),
			want: value.Int64(0),
		},
		"shift by 10 << 1": {
			a:    value.Int64(10),
			b:    value.Int64(1),
			want: value.Int64(20),
		},
		"shift by 10 << 255": {
			a:    value.Int64(10),
			b:    value.Int64(255),
			want: value.Int64(0),
		},
		"shift by 25 << 2": {
			a:    value.Int64(25),
			b:    value.Int64(2),
			want: value.Int64(100),
		},
		"shift by 25 << -2": {
			a:    value.Int64(25),
			b:    value.Int64(-2),
			want: value.Int64(6),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictIntLeftBitshift(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictInt_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"perform AND for String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"perform AND for Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"11 & 7": {
			a:    value.Int64(0b1011),
			b:    value.Int64(0b0111),
			want: value.Int64(0b0011),
		},
		"-14 & 23": {
			a:    value.Int64(-14),
			b:    value.Int64(23),
			want: value.Int64(18),
		},
		"258 & 0": {
			a:    value.Int64(258),
			b:    value.Int64(0),
			want: value.Int64(0),
		},
		"124 & 255": {
			a:    value.Int64(124),
			b:    value.Int64(255),
			want: value.Int64(124),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictIntBitwiseAnd(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictInt_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"perform OR for String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"perform OR for Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"21 | 13": {
			a:    value.Int64(0b10101),
			b:    value.Int64(0b01101),
			want: value.Int64(0b11101),
		},
		"-14 | 23": {
			a:    value.Int64(-14),
			b:    value.Int64(23),
			want: value.Int64(-9),
		},
		"258 | 0": {
			a:    value.Int64(258),
			b:    value.Int64(0),
			want: value.Int64(258),
		},
		"124 | 255": {
			a:    value.Int64(124),
			b:    value.Int64(255),
			want: value.Int64(255),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictIntBitwiseOr(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictInt_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"perform XOR for String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"perform XOR for Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"21 ^ 13": {
			a:    value.Int64(0b10101),
			b:    value.Int64(0b01101),
			want: value.Int64(0b11000),
		},
		"-14 ^ 23": {
			a:    value.Int64(-14),
			b:    value.Int64(23),
			want: value.Int64(-27),
		},
		"258 ^ 0": {
			a:    value.Int64(258),
			b:    value.Int64(0),
			want: value.Int64(258),
		},
		"124 ^ 255": {
			a:    value.Int64(124),
			b:    value.Int64(255),
			want: value.Int64(131),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictIntBitwiseXor(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictInt_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  *value.Error
	}{
		"perform modulo for String and return an error": {
			a:   value.Int64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"perform modulo for Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"21 % 10": {
			a:    value.Int64(21),
			b:    value.Int64(10),
			want: value.Int64(1),
		},
		"38 % 3": {
			a:    value.Int64(38),
			b:    value.Int64(3),
			want: value.Int64(2),
		},
		"522 % 39": {
			a:    value.Int64(522),
			b:    value.Int64(39),
			want: value.Int64(15),
		},
		"38 % 0": {
			a:   value.Int64(38),
			b:   value.Int64(0),
			err: value.NewError(value.ZeroDivisionErrorClass, "can't divide by zero"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictIntModulo(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictFloat_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Float64
		err  *value.Error
	}{
		"perform modulo for String and return an error": {
			a:   value.Float64(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::Float64`"),
		},
		"perform modulo for Int32 and return an error": {
			a:   value.Float64(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Float64`"),
		},
		"perform modulo for Float32 and return an error": {
			a:   value.Float64(5),
			b:   value.Float32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Float32` can't be coerced into `Std::Float64`"),
		},
		"21 % 10": {
			a:    value.Float64(21),
			b:    value.Float64(10),
			want: value.Float64(1),
		},
		"38 % 3": {
			a:    value.Float64(38),
			b:    value.Float64(3),
			want: value.Float64(2),
		},
		"522 % 39": {
			a:    value.Float64(522),
			b:    value.Float64(39),
			want: value.Float64(15),
		},
		"56.87 % 3": {
			a:    value.Float64(56.87),
			b:    value.Float64(3),
			want: value.Float64(2.8699999999999974),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := value.StrictFloatModulo(tc.a, tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
