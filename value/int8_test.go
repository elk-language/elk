package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestInt8Inspect(t *testing.T) {
	tests := map[string]struct {
		i    value.Int8
		want string
	}{
		"positive number": {
			i:    125,
			want: "125i8",
		},
		"negative number": {
			i:    -25,
			want: "-25i8",
		},
		"zero": {
			i:    0,
			want: "0i8",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.i.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestInt8_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Int8
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"exponentiate Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"exponentiate positive Int8 5 ** 2": {
			a:    value.Int8(5),
			b:    value.Int8(2).ToValue(),
			want: value.Int8(25),
		},
		"exponentiate positive Int8 2 ** 5": {
			a:    value.Int8(2),
			b:    value.Int8(5).ToValue(),
			want: value.Int8(32),
		},
		"exponentiate positive Int8 6 ** 1": {
			a:    value.Int8(6),
			b:    value.Int8(1).ToValue(),
			want: value.Int8(6),
		},
		"exponentiate negative Int8": {
			a:    value.Int8(4),
			b:    value.Int8(-2).ToValue(),
			want: value.Int8(1),
		},
		"exponentiate zero": {
			a:    value.Int8(25),
			b:    value.Int8(0).ToValue(),
			want: value.Int8(1),
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

func TestInt8_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Int8
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"add Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"add Int8": {
			a:    value.Int8(5),
			b:    value.Int8(2).ToValue(),
			want: value.Int8(7),
		},
		"add positive Int8": {
			a:    value.Int8(53),
			b:    value.Int8(21).ToValue(),
			want: value.Int8(74),
		},
		"add negative Int8": {
			a:    value.Int8(25),
			b:    value.Int8(-50).ToValue(),
			want: value.Int8(-25),
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

func TestInt8_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Int8
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"subtract Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"subtract positive Int8": {
			a:    value.Int8(53),
			b:    value.Int8(21).ToValue(),
			want: value.Int8(32),
		},
		"subtract negative Int8": {
			a:    value.Int8(25),
			b:    value.Int8(-50).ToValue(),
			want: value.Int8(75),
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

func TestInt8_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Int8
		err  value.Value
	}{
		"multiply String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"multiply Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"multiply positive Int8": {
			a:    value.Int8(53),
			b:    value.Int8(2).ToValue(),
			want: value.Int8(106),
		},
		"multiply negative Int8": {
			a:    value.Int8(25),
			b:    value.Int8(-2).ToValue(),
			want: value.Int8(-50),
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

func TestInt8_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"15i8 > 30i8": {
			a:    value.Int8(15),
			b:    value.Int8(30).ToValue(),
			want: value.False,
		},
		"78i8 > -80i8": {
			a:    value.Int8(78),
			b:    value.Int8(-80).ToValue(),
			want: value.True,
		},
		"15i8 > 15i8": {
			a:    value.Int8(15),
			b:    value.Int8(15).ToValue(),
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

func TestInt8_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"15i8 >= 30i8": {
			a:    value.Int8(15),
			b:    value.Int8(30).ToValue(),
			want: value.False,
		},
		"78i8 >= -80i8": {
			a:    value.Int8(78),
			b:    value.Int8(-80).ToValue(),
			want: value.True,
		},
		"15i8 >= 15i8": {
			a:    value.Int8(15),
			b:    value.Int8(15).ToValue(),
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

func TestInt8_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"15i8 < 30i8": {
			a:    value.Int8(15),
			b:    value.Int8(30).ToValue(),
			want: value.True,
		},
		"7 < -8": {
			a:    value.Int8(7),
			b:    value.Int8(-8).ToValue(),
			want: value.False,
		},
		"15i8 < 15i8": {
			a:    value.Int8(15),
			b:    value.Int8(15).ToValue(),
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

func TestInt8_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"15i8 <= 30i8": {
			a:    value.Int8(15),
			b:    value.Int8(30).ToValue(),
			want: value.True,
		},
		"7 <= -8": {
			a:    value.Int8(7),
			b:    value.Int8(-8).ToValue(),
			want: value.False,
		},
		"15i8 <= 15i8": {
			a:    value.Int8(15),
			b:    value.Int8(15).ToValue(),
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

func TestInt8_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String 5i8 == '5'": {
			a:    value.Int8(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},

		"Int64 5i8 == 5i64": {
			a:    value.Int8(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4i8 == 5i64": {
			a:    value.Int8(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 5i8 == -5i64": {
			a:    value.Int8(5),
			b:    value.Int64(-5).ToValue(),
			want: value.False,
		},

		"UInt64 5i8 == 5u64": {
			a:    value.Int8(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 -5i8 == 5u64": {
			a:    value.Int8(-5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4i8 == 5u64": {
			a:    value.Int8(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"Int32 5i8 == 5i32": {
			a:    value.Int8(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4i8 == 5i32": {
			a:    value.Int8(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"UInt32 5i8 == 5u32": {
			a:    value.Int8(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 -5i8 == 5u32": {
			a:    value.Int8(-5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 4i8 == 5u32": {
			a:    value.Int8(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"Int16 5i8 == 5i16": {
			a:    value.Int8(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4i8 == 5i16": {
			a:    value.Int8(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"UInt16 5i8 == 5u16": {
			a:    value.Int8(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 -5i8 == 5u16": {
			a:    value.Int8(-5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 4i8 == 5u16": {
			a:    value.Int8(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"Int8 5i8 == 5i8": {
			a:    value.Int8(5),
			b:    value.Int8(5).ToValue(),
			want: value.True,
		},
		"Int8 4i8 == 5i8": {
			a:    value.Int8(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt8 5i8 == 5u8": {
			a:    value.Int8(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 -5i8 == 5u8": {
			a:    value.Int8(-5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 4i8 == 5u8": {
			a:    value.Int8(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5i8 == 5.0f64": {
			a:    value.Int8(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5i8 == 5.5f64": {
			a:    value.Int8(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 5i8 == NaN": {
			a:    value.Int8(5),
			b:    value.Float64NaN().ToValue(),
			want: value.False,
		},
		"Float64 5i8 == +Inf": {
			a:    value.Int8(5),
			b:    value.Float64Inf().ToValue(),
			want: value.False,
		},
		"Float64 5i8 == -Inf": {
			a:    value.Int8(5),
			b:    value.Float64NegInf().ToValue(),
			want: value.False,
		},

		"Float32 5i8 == 5.0f32": {
			a:    value.Int8(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5i8 == 5.5f32": {
			a:    value.Int8(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 5i8 == NaN": {
			a:    value.Int8(5),
			b:    value.Float32NaN().ToValue(),
			want: value.False,
		},
		"Float32 5i8 == +Inf": {
			a:    value.Int8(5),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 5i8 == -Inf": {
			a:    value.Int8(5),
			b:    value.Float32NegInf().ToValue(),
			want: value.False,
		},

		"SmallInt 5i8 == 5": {
			a:    value.Int8(5),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},
		"SmallInt 4i8 == 5": {
			a:    value.Int8(4),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},

		"BigInt 5i8 == 5bi": {
			a:    value.Int8(5),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},
		"BigInt 4i8 == 5bi": {
			a:    value.Int8(4),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},

		"Float 5i8 == 5.0": {
			a:    value.Int8(5),
			b:    value.Float(5).ToValue(),
			want: value.False,
		},
		"Float 5i8 == 5.5": {
			a:    value.Int8(5),
			b:    value.Float(5.5).ToValue(),
			want: value.False,
		},
		"Float 5i8 == +Inf": {
			a:    value.Int8(5),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 5i8 == -Inf": {
			a:    value.Int8(5),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 5i8 == NaN": {
			a:    value.Int8(5),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 5i8 == 5.0bf": {
			a:    value.Int8(5),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.False,
		},
		"BigFloat 5i8 == 5.5bf": {
			a:    value.Int8(5),
			b:    value.Ref(value.NewBigFloat(5.5)),
			want: value.False,
		},
		"BigFloat 5i8 == +Inf": {
			a:    value.Int8(5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 5i8 == -Inf": {
			a:    value.Int8(5),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 5i8 == NaN": {
			a:    value.Int8(5),
			b:    value.Ref(value.BigFloatNaN()),
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

func TestInt8_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Int8
		err  value.Value
	}{
		"perform AND for String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"perform AND for Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"11 & 7": {
			a:    value.Int8(0b1011),
			b:    value.Int8(0b0111).ToValue(),
			want: value.Int8(0b0011),
		},
		"-14 & 23": {
			a:    value.Int8(-14),
			b:    value.Int8(23).ToValue(),
			want: value.Int8(18),
		},
		"124 & 0": {
			a:    value.Int8(124),
			b:    value.Int8(0).ToValue(),
			want: value.Int8(0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseAnd(tc.b)
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

func TestInt8_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Int8
		err  value.Value
	}{
		"perform OR for String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"perform OR for Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"21 | 13": {
			a:    value.Int8(0b10101),
			b:    value.Int8(0b01101).ToValue(),
			want: value.Int8(0b11101),
		},
		"-14 | 23": {
			a:    value.Int8(-14),
			b:    value.Int8(23).ToValue(),
			want: value.Int8(-9),
		},
		"124 | 0": {
			a:    value.Int8(124),
			b:    value.Int8(0).ToValue(),
			want: value.Int8(124),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseOr(tc.b)
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

func TestInt8_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Int8
		err  value.Value
	}{
		"perform XOR for String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"perform XOR for Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"21 ^ 13": {
			a:    value.Int8(0b10101),
			b:    value.Int8(0b01101).ToValue(),
			want: value.Int8(0b11000),
		},
		"-14 ^ 23": {
			a:    value.Int8(-14),
			b:    value.Int8(23).ToValue(),
			want: value.Int8(-27),
		},
		"124 ^ 0": {
			a:    value.Int8(124),
			b:    value.Int8(0).ToValue(),
			want: value.Int8(124),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseXor(tc.b)
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

func TestInt8_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.Int8
		b    value.Value
		want value.Int8
		err  value.Value
	}{
		"perform modulo for String and return an error": {
			a:   value.Int8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int8`")),
		},
		"perform modulo for Int64 and return an error": {
			a:   value.Int8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int8`")),
		},
		"21 % 10": {
			a:    value.Int8(21),
			b:    value.Int8(10).ToValue(),
			want: value.Int8(1),
		},
		"38 % 3": {
			a:    value.Int8(38),
			b:    value.Int8(3).ToValue(),
			want: value.Int8(2),
		},
		"122 % 39": {
			a:    value.Int8(122),
			b:    value.Int8(39).ToValue(),
			want: value.Int8(5),
		},
		"38 % 0": {
			a:   value.Int8(38),
			b:   value.Int8(0).ToValue(),
			err: value.Ref(value.NewError(value.ZeroDivisionErrorClass, "cannot divide by zero")),
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
