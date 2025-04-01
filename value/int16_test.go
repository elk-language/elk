package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestInt16Inspect(t *testing.T) {
	tests := map[string]struct {
		i    value.Int16
		want string
	}{
		"positive number": {
			i:    216,
			want: "216i16",
		},
		"negative number": {
			i:    -25,
			want: "-25i16",
		},
		"zero": {
			i:    0,
			want: "0i16",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.i.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Int16
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.Int16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"exponentiate Int64 and return an error": {
			a:   value.Int16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"exponentiate positive Int16 5 ** 2": {
			a:    value.Int16(5),
			b:    value.Int16(2).ToValue(),
			want: value.Int16(25),
		},
		"exponentiate positive Int16 2 ** 5": {
			a:    value.Int16(2),
			b:    value.Int16(5).ToValue(),
			want: value.Int16(32),
		},
		"exponentiate positive Int16 6 ** 1": {
			a:    value.Int16(6),
			b:    value.Int16(1).ToValue(),
			want: value.Int16(6),
		},
		"exponentiate negative Int16": {
			a:    value.Int16(4),
			b:    value.Int16(-2).ToValue(),
			want: value.Int16(1),
		},
		"exponentiate zero": {
			a:    value.Int16(25),
			b:    value.Int16(0).ToValue(),
			want: value.Int16(1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.ExponentiateVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Int16
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.Int16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"add Int16": {
			a:    value.Int16(5),
			b:    value.Int16(2).ToValue(),
			want: value.Int16(7),
		},
		"add positive Int16": {
			a:    value.Int16(53),
			b:    value.Int16(21).ToValue(),
			want: value.Int16(74),
		},
		"add negative Int16": {
			a:    value.Int16(25),
			b:    value.Int16(-50).ToValue(),
			want: value.Int16(-25),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Add(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Int16
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.Int16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"subtract Int64 and return an error": {
			a:   value.Int16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"subtract positive Int16": {
			a:    value.Int16(53),
			b:    value.Int16(21).ToValue(),
			want: value.Int16(32),
		},
		"subtract negative Int16": {
			a:    value.Int16(25),
			b:    value.Int16(-50).ToValue(),
			want: value.Int16(75),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Subtract(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Int16
		err  value.Value
	}{
		"multiply String and return an error": {
			a:   value.Int16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"multiply Int64 and return an error": {
			a:   value.Int16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"multiply positive Int16": {
			a:    value.Int16(53),
			b:    value.Int16(2).ToValue(),
			want: value.Int16(106),
		},
		"multiply negative Int16": {
			a:    value.Int16(25),
			b:    value.Int16(-2).ToValue(),
			want: value.Int16(-50),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Multiply(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Int16(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"Int64 and return an error": {
			a:    value.Int16(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"15i16 > 30i16": {
			a:    value.Int16(15),
			b:    value.Int16(30).ToValue(),
			want: value.False,
		},
		"780i16 > -800i16": {
			a:    value.Int16(780),
			b:    value.Int16(-800).ToValue(),
			want: value.True,
		},
		"15i16 > 15i16": {
			a:    value.Int16(15),
			b:    value.Int16(15).ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Int16(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"Int64 and return an error": {
			a:    value.Int16(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"15i16 >= 30i16": {
			a:    value.Int16(15),
			b:    value.Int16(30).ToValue(),
			want: value.False,
		},
		"780i16 >= -800i16": {
			a:    value.Int16(780),
			b:    value.Int16(-800).ToValue(),
			want: value.True,
		},
		"15i16 >= 15i16": {
			a:    value.Int16(15),
			b:    value.Int16(15).ToValue(),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqualVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Int16(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"Int64 and return an error": {
			a:    value.Int16(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"15i16 < 30i16": {
			a:    value.Int16(15),
			b:    value.Int16(30).ToValue(),
			want: value.True,
		},
		"7 < -8": {
			a:    value.Int16(7),
			b:    value.Int16(-8).ToValue(),
			want: value.False,
		},
		"15i16 < 15i16": {
			a:    value.Int16(15),
			b:    value.Int16(15).ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Int16(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"Int64 and return an error": {
			a:    value.Int16(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"15i16 <= 30i16": {
			a:    value.Int16(15),
			b:    value.Int16(30).ToValue(),
			want: value.True,
		},
		"7 <= -8": {
			a:    value.Int16(7),
			b:    value.Int16(-8).ToValue(),
			want: value.False,
		},
		"15i16 <= 15i16": {
			a:    value.Int16(15),
			b:    value.Int16(15).ToValue(),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqualVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String 5i16 == '5'": {
			a:    value.Int16(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},

		"Int64 5i16 == 5i64": {
			a:    value.Int16(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4i16 == 5i64": {
			a:    value.Int16(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 5i16 == -5i64": {
			a:    value.Int16(5),
			b:    value.Int64(-5).ToValue(),
			want: value.False,
		},

		"UInt64 5i16 == 5u64": {
			a:    value.Int16(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 -5i16 == 5u64": {
			a:    value.Int16(-5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4i16 == 5u64": {
			a:    value.Int16(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"Int32 5i16 == 5i32": {
			a:    value.Int16(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4i16 == 5i32": {
			a:    value.Int16(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"UInt32 5i16 == 5u32": {
			a:    value.Int16(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 -5i16 == 5u32": {
			a:    value.Int16(-5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 4i16 == 5u32": {
			a:    value.Int16(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"Int16 5i16 == 5i16": {
			a:    value.Int16(5),
			b:    value.Int16(5).ToValue(),
			want: value.True,
		},
		"Int16 4i16 == 5i16": {
			a:    value.Int16(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"UInt16 5i16 == 5u16": {
			a:    value.Int16(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 -5i16 == 5u16": {
			a:    value.Int16(-5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 4i16 == 5u16": {
			a:    value.Int16(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"Int8 5i16 == 5i8": {
			a:    value.Int16(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4i16 == 5i8": {
			a:    value.Int16(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt8 5i16 == 5u8": {
			a:    value.Int16(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 -5i16 == 5u8": {
			a:    value.Int16(-5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 4i16 == 5u8": {
			a:    value.Int16(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5i16 == 5.0f64": {
			a:    value.Int16(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5i16 == 5.5f64": {
			a:    value.Int16(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 5i16 == NaN": {
			a:    value.Int16(5),
			b:    value.Float64NaN().ToValue(),
			want: value.False,
		},
		"Float64 5i16 == +Inf": {
			a:    value.Int16(5),
			b:    value.Float64Inf().ToValue(),
			want: value.False,
		},
		"Float64 5i16 == -Inf": {
			a:    value.Int16(5),
			b:    value.Float64NegInf().ToValue(),
			want: value.False,
		},

		"Float32 5i16 == 5.0f32": {
			a:    value.Int16(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5i16 == 5.5f32": {
			a:    value.Int16(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 5i16 == NaN": {
			a:    value.Int16(5),
			b:    value.Float32NaN().ToValue(),
			want: value.False,
		},
		"Float32 5i16 == +Inf": {
			a:    value.Int16(5),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 5i16 == -Inf": {
			a:    value.Int16(5),
			b:    value.Float32NegInf().ToValue(),
			want: value.False,
		},

		"SmallInt 5i16 == 5": {
			a:    value.Int16(5),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},
		"SmallInt 4i16 == 5": {
			a:    value.Int16(4),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},

		"BigInt 5i16 == 5bi": {
			a:    value.Int16(5),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},
		"BigInt 4i16 == 5bi": {
			a:    value.Int16(4),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},

		"Float 5i16 == 5.0": {
			a:    value.Int16(5),
			b:    value.Float(5).ToValue(),
			want: value.False,
		},
		"Float 5i16 == 5.5": {
			a:    value.Int16(5),
			b:    value.Float(5.5).ToValue(),
			want: value.False,
		},
		"Float 5i16 == +Inf": {
			a:    value.Int16(5),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 5i16 == -Inf": {
			a:    value.Int16(5),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 5i16 == NaN": {
			a:    value.Int16(5),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 5i16 == 5.0bf": {
			a:    value.Int16(5),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.False,
		},
		"BigFloat 5i16 == 5.5bf": {
			a:    value.Int16(5),
			b:    value.Ref(value.NewBigFloat(5.5)),
			want: value.False,
		},
		"BigFloat 5i16 == +Inf": {
			a:    value.Int16(5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 5i16 == -Inf": {
			a:    value.Int16(5),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 5i16 == NaN": {
			a:    value.Int16(5),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.EqualVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Int16
		err  value.Value
	}{
		"perform AND for String and return an error": {
			a:   value.Int16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"perform AND for Int64 and return an error": {
			a:   value.Int16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"11 & 7": {
			a:    value.Int16(0b1011),
			b:    value.Int16(0b0111).ToValue(),
			want: value.Int16(0b0011),
		},
		"-14 & 23": {
			a:    value.Int16(-14),
			b:    value.Int16(23).ToValue(),
			want: value.Int16(18),
		},
		"258 & 0": {
			a:    value.Int16(258),
			b:    value.Int16(0).ToValue(),
			want: value.Int16(0),
		},
		"124 & 255": {
			a:    value.Int16(124),
			b:    value.Int16(255).ToValue(),
			want: value.Int16(124),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseAnd(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Int16
		err  value.Value
	}{
		"perform OR for String and return an error": {
			a:   value.Int16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"perform OR for Int64 and return an error": {
			a:   value.Int16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"21 | 13": {
			a:    value.Int16(0b10101),
			b:    value.Int16(0b01101).ToValue(),
			want: value.Int16(0b11101),
		},
		"-14 | 23": {
			a:    value.Int16(-14),
			b:    value.Int16(23).ToValue(),
			want: value.Int16(-9),
		},
		"258 | 0": {
			a:    value.Int16(258),
			b:    value.Int16(0).ToValue(),
			want: value.Int16(258),
		},
		"124 | 255": {
			a:    value.Int16(124),
			b:    value.Int16(255).ToValue(),
			want: value.Int16(255),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseOr(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Int16
		err  value.Value
	}{
		"perform XOR for String and return an error": {
			a:   value.Int16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"perform XOR for Int64 and return an error": {
			a:   value.Int16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"21 ^ 13": {
			a:    value.Int16(0b10101),
			b:    value.Int16(0b01101).ToValue(),
			want: value.Int16(0b11000),
		},
		"-14 ^ 23": {
			a:    value.Int16(-14),
			b:    value.Int16(23).ToValue(),
			want: value.Int16(-27),
		},
		"258 ^ 0": {
			a:    value.Int16(258),
			b:    value.Int16(0).ToValue(),
			want: value.Int16(258),
		},
		"124 ^ 255": {
			a:    value.Int16(124),
			b:    value.Int16(255).ToValue(),
			want: value.Int16(131),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseXor(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt16_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.Int16
		b    value.Value
		want value.Int16
		err  value.Value
	}{
		"perform modulo for String and return an error": {
			a:   value.Int16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int16`")),
		},
		"perform modulo for Int64 and return an error": {
			a:   value.Int16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int16`")),
		},
		"21 % 10": {
			a:    value.Int16(21),
			b:    value.Int16(10).ToValue(),
			want: value.Int16(1),
		},
		"38 % 3": {
			a:    value.Int16(38),
			b:    value.Int16(3).ToValue(),
			want: value.Int16(2),
		},
		"522 % 39": {
			a:    value.Int16(522),
			b:    value.Int16(39).ToValue(),
			want: value.Int16(15),
		},
		"38 % 0": {
			a:   value.Int16(38),
			b:   value.Int16(0).ToValue(),
			err: value.Ref(value.NewError(value.ZeroDivisionErrorClass, "cannot divide by zero")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.ModuloVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
