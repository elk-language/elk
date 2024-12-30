package value_test

import (
	"math"
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestFloat64_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Float64
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.Float64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float64`")),
		},
		"exponentiate Int32 and return an error": {
			a:   value.Float64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float64`")),
		},
		"exponentiate positive value.Float64": {
			a:    value.Float64(5.5),
			b:    value.Float64(3).ToValue(),
			want: value.Float64(166.375),
		},
		"exponentiate negative value.Float64": {
			a:    value.Float64(5.5),
			b:    value.Float64(-2).ToValue(),
			want: value.Float64(0.03305785123966942),
		},
		"exponentiate zero": {
			a:    value.Float64(5.5),
			b:    value.Float64(0).ToValue(),
			want: value.Float64(1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Exponentiate(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat64_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Float64
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.Float64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float64`")),
		},
		"add Int32 and return an error": {
			a:   value.Float64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float64`")),
		},
		"add positive Float64": {
			a:    value.Float64(53.5),
			b:    value.Float64(21).ToValue(),
			want: value.Float64(74.5),
		},
		"add negative Float64": {
			a:    value.Float64(25.5),
			b:    value.Float64(-50).ToValue(),
			want: value.Float64(-24.5),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Add(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat64_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Float64
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.Float64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float64`")),
		},
		"subtract Int32 and return an error": {
			a:   value.Float64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float64`")),
		},
		"subtract positive Float64": {
			a:    value.Float64(53.5),
			b:    value.Float64(21).ToValue(),
			want: value.Float64(32.5),
		},
		"subtract negative Float64": {
			a:    value.Float64(25.5),
			b:    value.Float64(-50).ToValue(),
			want: value.Float64(75.5),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Subtract(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat64_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Float64
		err  value.Value
	}{
		"multiply String and return an error": {
			a:   value.Float64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float64`")),
		},
		"multiply Int32 and return an error": {
			a:   value.Float64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float64`")),
		},
		"multiply positive Float64": {
			a:    value.Float64(53.5),
			b:    value.Float64(2).ToValue(),
			want: value.Float64(107),
		},
		"multiply negative Float64": {
			a:    value.Float64(25.5),
			b:    value.Float64(-2).ToValue(),
			want: value.Float64(-51),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Multiply(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat64_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Float64(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float64`")),
		},
		"Int32 and return an error": {
			a:    value.Float64(5),
			b:    value.Int32(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float64`")),
		},
		"15f64 > 30f64": {
			a:    value.Float64(15),
			b:    value.Float64(30).ToValue(),
			want: value.False,
		},
		"780f64 > -800f64": {
			a:    value.Float64(780),
			b:    value.Float64(-800).ToValue(),
			want: value.True,
		},
		"15f64 > 15f64": {
			a:    value.Float64(15),
			b:    value.Float64(15).ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThan(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat64_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Float64(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float64`")),
		},
		"Int32 and return an error": {
			a:    value.Float64(5),
			b:    value.Int32(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float64`")),
		},
		"15f64 >= 30f64": {
			a:    value.Float64(15),
			b:    value.Float64(30).ToValue(),
			want: value.False,
		},
		"780f64 >= -800f64": {
			a:    value.Float64(780),
			b:    value.Float64(-800).ToValue(),
			want: value.True,
		},
		"15f64 >= 15f64": {
			a:    value.Float64(15),
			b:    value.Float64(15).ToValue(),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqual(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat64_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Float64(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float64`")),
		},
		"Int32 and return an error": {
			a:    value.Float64(5),
			b:    value.Int32(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float64`")),
		},
		"15f64 < 30f64": {
			a:    value.Float64(15),
			b:    value.Float64(30).ToValue(),
			want: value.True,
		},
		"780f64 < -800f64": {
			a:    value.Float64(780),
			b:    value.Float64(-800).ToValue(),
			want: value.False,
		},
		"15f64 < 15f64": {
			a:    value.Float64(15),
			b:    value.Float64(15).ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThan(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat64_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Float64(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float64`")),
		},
		"Int32 and return an error": {
			a:    value.Float64(5),
			b:    value.Int32(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float64`")),
		},
		"15f64 <= 30f64": {
			a:    value.Float64(15),
			b:    value.Float64(30).ToValue(),
			want: value.True,
		},
		"780f64 <= -800f64": {
			a:    value.Float64(780),
			b:    value.Float64(-800).ToValue(),
			want: value.False,
		},
		"15f64 <= 15f64": {
			a:    value.Float64(15),
			b:    value.Float64(15).ToValue(),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqual(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat64_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String 5.5f64 == '5.5'": {
			a:    value.Float64(5.5),
			b:    value.Ref(value.String("5.5")),
			want: value.False,
		},

		"Int64 5.0f64 == 5i64": {
			a:    value.Float64(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 5.5f64 == 5i64": {
			a:    value.Float64(5.5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 NaN == 0i64": {
			a:    value.Float64NaN(),
			b:    value.Int64(0).ToValue(),
			want: value.False,
		},
		"Int64 +Inf == 69i64": {
			a:    value.Float64Inf(),
			b:    value.Int64(69).ToValue(),
			want: value.False,
		},
		"Int64 -Inf == -89i64": {
			a:    value.Float64NegInf(),
			b:    value.Int64(-89).ToValue(),
			want: value.False,
		},

		"UInt64 5.0f64 == 5u64": {
			a:    value.Float64(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 -5.0f64 == 5u64": {
			a:    value.Float64(-5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 5.5f64 == 5u64": {
			a:    value.Float64(5.5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 NaN == 0u64": {
			a:    value.Float64NaN(),
			b:    value.UInt64(0).ToValue(),
			want: value.False,
		},
		"UInt64 +Inf == 69u64": {
			a:    value.Float64Inf(),
			b:    value.UInt64(69).ToValue(),
			want: value.False,
		},
		"UInt64 -Inf == 89u64": {
			a:    value.Float64NegInf(),
			b:    value.UInt64(89).ToValue(),
			want: value.False,
		},

		"Int32 5.0f64 == 5i32": {
			a:    value.Float64(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 5.5f64 == 5i32": {
			a:    value.Float64(5.5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 NaN == 0i32": {
			a:    value.Float64NaN(),
			b:    value.Int32(0).ToValue(),
			want: value.False,
		},
		"Int32 +Inf == 69i32": {
			a:    value.Float64Inf(),
			b:    value.Int32(69).ToValue(),
			want: value.False,
		},
		"Int32 -Inf == -89i32": {
			a:    value.Float64NegInf(),
			b:    value.Int32(-89).ToValue(),
			want: value.False,
		},

		"UInt32 5.0f64 == 5u32": {
			a:    value.Float64(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 -5.0f64 == 5u32": {
			a:    value.Float64(-5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 5.5f64 == 5u32": {
			a:    value.Float64(5.5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 NaN == 0u32": {
			a:    value.Float64NaN(),
			b:    value.UInt32(0).ToValue(),
			want: value.False,
		},
		"UInt32 +Inf == 69u32": {
			a:    value.Float64Inf(),
			b:    value.UInt32(69).ToValue(),
			want: value.False,
		},
		"UInt32 -Inf == 89u32": {
			a:    value.Float64NegInf(),
			b:    value.UInt32(89).ToValue(),
			want: value.False,
		},

		"Int16 5.0f64 == 5i16": {
			a:    value.Float64(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 5.5f64 == 5i16": {
			a:    value.Float64(5.5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 NaN == 0i16": {
			a:    value.Float64NaN(),
			b:    value.Int16(0).ToValue(),
			want: value.False,
		},
		"Int16 +Inf == 69i16": {
			a:    value.Float64Inf(),
			b:    value.Int16(69).ToValue(),
			want: value.False,
		},
		"Int16 -Inf == -89i16": {
			a:    value.Float64NegInf(),
			b:    value.Int16(-89).ToValue(),
			want: value.False,
		},

		"UInt16 5.0f64 == 5u16": {
			a:    value.Float64(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 -5.0f64 == 5u16": {
			a:    value.Float64(-5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 5.5f64 == 5u16": {
			a:    value.Float64(5.5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 NaN == 0u16": {
			a:    value.Float64NaN(),
			b:    value.UInt16(0).ToValue(),
			want: value.False,
		},
		"UInt16 +Inf == 69u16": {
			a:    value.Float64Inf(),
			b:    value.UInt16(69).ToValue(),
			want: value.False,
		},
		"UInt16 -Inf == 89u16": {
			a:    value.Float64NegInf(),
			b:    value.UInt16(89).ToValue(),
			want: value.False,
		},

		"Int8 5.0f64 == 5i8": {
			a:    value.Float64(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 5.5f64 == 5i8": {
			a:    value.Float64(5.5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 NaN == 0i8": {
			a:    value.Float64NaN(),
			b:    value.Int8(0).ToValue(),
			want: value.False,
		},
		"Int8 +Inf == 69i8": {
			a:    value.Float64Inf(),
			b:    value.Int8(69).ToValue(),
			want: value.False,
		},
		"Int8 -Inf == -89i8": {
			a:    value.Float64NegInf(),
			b:    value.Int8(-89).ToValue(),
			want: value.False,
		},

		"UInt8 5.0f64 == 5u8": {
			a:    value.Float64(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 -5.0f64 == 5u8": {
			a:    value.Float64(-5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 5.5f64 == 5u8": {
			a:    value.Float64(5.5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 NaN == 0u8": {
			a:    value.Float64NaN(),
			b:    value.UInt8(0).ToValue(),
			want: value.False,
		},
		"UInt8 +Inf == 69u8": {
			a:    value.Float64Inf(),
			b:    value.UInt8(69).ToValue(),
			want: value.False,
		},
		"UInt8 -Inf == 89u8": {
			a:    value.Float64NegInf(),
			b:    value.UInt8(89).ToValue(),
			want: value.False,
		},

		"value.Float64 21.9f64 == 21.9f64": {
			a:    value.Float64(21.9),
			b:    value.Float64(21.9).ToValue(),
			want: value.True,
		},
		"value.Float64 21.9f64 == 38.0f64": {
			a:    value.Float64(21.9),
			b:    value.Float64(38).ToValue(),
			want: value.False,
		},
		"value.Float64 NaN == NaN": {
			a:    value.Float64NaN(),
			b:    value.Float64NaN().ToValue(),
			want: value.False,
		},
		"value.Float64 +Inf == +Inf": {
			a:    value.Float64Inf(),
			b:    value.Float64Inf().ToValue(),
			want: value.True,
		},
		"value.Float64 -Inf == -Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float64NegInf().ToValue(),
			want: value.True,
		},
		"value.Float64 +Inf == -Inf": {
			a:    value.Float64Inf(),
			b:    value.Float64NegInf().ToValue(),
			want: value.False,
		},
		"value.Float64 -Inf == +Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float64Inf().ToValue(),
			want: value.False,
		},
		"value.Float64 8.5f64 == +Inf": {
			a:    value.Float64(8.5),
			b:    value.Float64Inf().ToValue(),
			want: value.False,
		},
		"value.Float64 +Inf == 98.0f64": {
			a:    value.Float64Inf(),
			b:    value.Float64(98).ToValue(),
			want: value.False,
		},

		"Float32 21.0f64 == 21.0f32": {
			a:    value.Float64(21),
			b:    value.Float32(21).ToValue(),
			want: value.False,
		},
		"Float32 21.9f64 == 38.0f32": {
			a:    value.Float64(21.9),
			b:    value.Float32(38).ToValue(),
			want: value.False,
		},
		"Float32 NaN == NaN": {
			a:    value.Float64NaN(),
			b:    value.Float32NaN().ToValue(),
			want: value.False,
		},
		"Float32 +Inf == +Inf": {
			a:    value.Float64Inf(),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 -Inf == -Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float32NegInf().ToValue(),
			want: value.False,
		},
		"Float32 +Inf == -Inf": {
			a:    value.Float64Inf(),
			b:    value.Float32NegInf().ToValue(),
			want: value.False,
		},
		"Float32 -Inf == +Inf": {
			a:    value.Float64NegInf(),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 8.5f64 == +Inf": {
			a:    value.Float64(8.5),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 +Inf == 98.0f32": {
			a:    value.Float64Inf(),
			b:    value.Float32(98).ToValue(),
			want: value.False,
		},

		"SmallInt 16.0f64 == 16": {
			a:    value.Float64(16),
			b:    value.SmallInt(16).ToValue(),
			want: value.False,
		},
		"SmallInt 16.5f64 == 16": {
			a:    value.Float64(16.5),
			b:    value.SmallInt(16).ToValue(),
			want: value.False,
		},
		"SmallInt NaN == 0": {
			a:    value.Float64NaN(),
			b:    value.SmallInt(0).ToValue(),
			want: value.False,
		},
		"SmallInt +Inf == 69": {
			a:    value.Float64Inf(),
			b:    value.SmallInt(69).ToValue(),
			want: value.False,
		},
		"SmallInt -Inf == -89": {
			a:    value.Float64NegInf(),
			b:    value.SmallInt(-89).ToValue(),
			want: value.False,
		},

		"BigInt 16.0f64 == 16": {
			a:    value.Float64(16),
			b:    value.Ref(value.NewBigInt(16)),
			want: value.False,
		},
		"BigInt 16.5f64 == 16": {
			a:    value.Float64(16.5),
			b:    value.Ref(value.NewBigInt(16)),
			want: value.False,
		},
		"BigInt NaN == 0": {
			a:    value.Float64NaN(),
			b:    value.Ref(value.NewBigInt(0)),
			want: value.False,
		},
		"BigInt +Inf == 69": {
			a:    value.Float64Inf(),
			b:    value.Ref(value.NewBigInt(69)),
			want: value.False,
		},
		"BigInt -Inf == -89": {
			a:    value.Float64NegInf(),
			b:    value.Ref(value.NewBigInt(-89)),
			want: value.False,
		},

		"Float 21.9f64 == 21.9": {
			a:    value.Float64(21.9),
			b:    value.Float(21.9).ToValue(),
			want: value.False,
		},
		"Float 21.9f64 == 38.0": {
			a:    value.Float64(21.9),
			b:    value.Float(38).ToValue(),
			want: value.False,
		},
		"Float NaN == NaN": {
			a:    value.Float64NaN(),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"Float +Inf == +Inf": {
			a:    value.Float64Inf(),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float -Inf == -Inf": {
			a:    value.Float64NegInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float +Inf == -Inf": {
			a:    value.Float64Inf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float -Inf == +Inf": {
			a:    value.Float64NegInf(),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 8.5f64 == +Inf": {
			a:    value.Float64(8.5),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float +Inf == 98.0": {
			a:    value.Float64Inf(),
			b:    value.Float(98).ToValue(),
			want: value.False,
		},

		"BigFloat 21.9f64 == 21.9bf": {
			a:    value.Float64(21.9),
			b:    value.Ref(value.NewBigFloat(21.9)),
			want: value.False,
		},
		"BigFloat 21.9f64 == 38.0bf": {
			a:    value.Float64(21.9),
			b:    value.Ref(value.NewBigFloat(38)),
			want: value.False,
		},
		"BigFloat NaN == NaN": {
			a:    value.Float64NaN(),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
		"BigFloat +Inf == +Inf": {
			a:    value.Float64Inf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat -Inf == -Inf": {
			a:    value.Float64NegInf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat +Inf == -Inf": {
			a:    value.Float64Inf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat -Inf == +Inf": {
			a:    value.Float64NegInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 8.5f64 == +Inf": {
			a:    value.Float64(8.5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat +Inf == 98.0bf": {
			a:    value.Float64Inf(),
			b:    value.Ref(value.NewBigFloat(98)),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat64_Divide(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Float64
		err  value.Value
	}{
		"divide by String and return an error": {
			a:   value.Float64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float64`")),
		},
		"divide Int32 and return an error": {
			a:   value.Float64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float64`")),
		},
		"divide positive value.Float64": {
			a:    value.Float64(54.5),
			b:    value.Float64(2).ToValue(),
			want: value.Float64(27.25),
		},
		"divide negative value.Float64": {
			a:    value.Float64(50),
			b:    value.Float64(-2).ToValue(),
			want: value.Float64(-25),
		},
		"divide by zero": {
			a:    value.Float64(50),
			b:    value.Float64(0).ToValue(),
			want: value.Float64(math.Inf(1)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Divide(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat64_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.Float64
		b    value.Value
		want value.Float64
		err  value.Value
	}{
		"perform modulo for String and return an error": {
			a:   value.Float64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float64`")),
		},
		"perform modulo for Int32 and return an error": {
			a:   value.Float64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float64`")),
		},
		"perform modulo for Float32 and return an error": {
			a:   value.Float64(5),
			b:   value.Float32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::Float64`")),
		},
		"21 % 10": {
			a:    value.Float64(21),
			b:    value.Float64(10).ToValue(),
			want: value.Float64(1),
		},
		"38 % 3": {
			a:    value.Float64(38),
			b:    value.Float64(3).ToValue(),
			want: value.Float64(2),
		},
		"522 % 39": {
			a:    value.Float64(522),
			b:    value.Float64(39).ToValue(),
			want: value.Float64(15),
		},
		"56.87 % 3": {
			a:    value.Float64(56.87),
			b:    value.Float64(3).ToValue(),
			want: value.Float64(2.8699999999999974),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Modulo(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
