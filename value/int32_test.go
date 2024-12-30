package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestInt32Inspect(t *testing.T) {
	tests := map[string]struct {
		i    value.Int32
		want string
	}{
		"positive number": {
			i:    216,
			want: "216i32",
		},
		"negative number": {
			i:    -25,
			want: "-25i32",
		},
		"zero": {
			i:    0,
			want: "0i32",
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

func TestInt32_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Int32
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.Int32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"exponentiate Int64 and return an error": {
			a:   value.Int32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"exponentiate positive Int32 5 ** 2": {
			a:    value.Int32(5),
			b:    value.Int32(2).ToValue(),
			want: value.Int32(25),
		},
		"exponentiate positive Int32 7 ** 8": {
			a:    value.Int32(7),
			b:    value.Int32(8).ToValue(),
			want: value.Int32(5764801),
		},
		"exponentiate positive Int32 2 ** 5": {
			a:    value.Int32(2),
			b:    value.Int32(5).ToValue(),
			want: value.Int32(32),
		},
		"exponentiate positive Int32 6 ** 1": {
			a:    value.Int32(6),
			b:    value.Int32(1).ToValue(),
			want: value.Int32(6),
		},
		"exponentiate negative Int32": {
			a:    value.Int32(4),
			b:    value.Int32(-2).ToValue(),
			want: value.Int32(1),
		},
		"exponentiate zero": {
			a:    value.Int32(25),
			b:    value.Int32(0).ToValue(),
			want: value.Int32(1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Exponentiate(tc.b)
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

func TestInt32_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Int32
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.Int32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"add Int32": {
			a:    value.Int32(5),
			b:    value.Int32(2).ToValue(),
			want: value.Int32(7),
		},
		"add positive Int32": {
			a:    value.Int32(53),
			b:    value.Int32(21).ToValue(),
			want: value.Int32(74),
		},
		"add negative Int32": {
			a:    value.Int32(25),
			b:    value.Int32(-50).ToValue(),
			want: value.Int32(-25),
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

func TestInt32_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Int32
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.Int32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"subtract Int64 and return an error": {
			a:   value.Int32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"subtract positive Int32": {
			a:    value.Int32(53),
			b:    value.Int32(21).ToValue(),
			want: value.Int32(32),
		},
		"subtract negative Int32": {
			a:    value.Int32(25),
			b:    value.Int32(-50).ToValue(),
			want: value.Int32(75),
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

func TestInt32_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Int32
		err  value.Value
	}{
		"multiply String and return an error": {
			a:   value.Int32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"multiply Int64 and return an error": {
			a:   value.Int32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"multiply positive Int32": {
			a:    value.Int32(53),
			b:    value.Int32(2).ToValue(),
			want: value.Int32(106),
		},
		"multiply negative Int32": {
			a:    value.Int32(25),
			b:    value.Int32(-2).ToValue(),
			want: value.Int32(-50),
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

func TestInt32_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Int32(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"Int64 and return an error": {
			a:    value.Int32(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"15i32 > 30i32": {
			a:    value.Int32(15),
			b:    value.Int32(30).ToValue(),
			want: value.False,
		},
		"780i32 > -800i32": {
			a:    value.Int32(780),
			b:    value.Int32(-800).ToValue(),
			want: value.True,
		},
		"15i32 > 15i32": {
			a:    value.Int32(15),
			b:    value.Int32(15).ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThan(tc.b)
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

func TestInt32_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Int32(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"Int64 and return an error": {
			a:    value.Int32(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"15i32 >= 30i32": {
			a:    value.Int32(15),
			b:    value.Int32(30).ToValue(),
			want: value.False,
		},
		"780i32 >= -800i32": {
			a:    value.Int32(780),
			b:    value.Int32(-800).ToValue(),
			want: value.True,
		},
		"15i32 >= 15i32": {
			a:    value.Int32(15),
			b:    value.Int32(15).ToValue(),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqual(tc.b)
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

func TestInt32_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Int32(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"Int64 and return an error": {
			a:    value.Int32(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"15i32 < 30i32": {
			a:    value.Int32(15),
			b:    value.Int32(30).ToValue(),
			want: value.True,
		},
		"780i32 < -800i32": {
			a:    value.Int32(780),
			b:    value.Int32(-800).ToValue(),
			want: value.False,
		},
		"15i32 < 15i32": {
			a:    value.Int32(15),
			b:    value.Int32(15).ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThan(tc.b)
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

func TestInt32_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.Int32(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"Int64 and return an error": {
			a:    value.Int32(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"15i32 <= 30i32": {
			a:    value.Int32(15),
			b:    value.Int32(30).ToValue(),
			want: value.True,
		},
		"780i32 <= -800i32": {
			a:    value.Int32(780),
			b:    value.Int32(-800).ToValue(),
			want: value.False,
		},
		"15i32 <= 15i32": {
			a:    value.Int32(15),
			b:    value.Int32(15).ToValue(),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqual(tc.b)
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

func TestInt32_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String 5i32 == '5'": {
			a:    value.Int32(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},

		"Int64 5i32 == 5i64": {
			a:    value.Int32(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4i32 == 5i64": {
			a:    value.Int32(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 5i32 == -5i64": {
			a:    value.Int32(5),
			b:    value.Int64(-5).ToValue(),
			want: value.False,
		},

		"UInt64 5i32 == 5u64": {
			a:    value.Int32(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 -5i32 == 5u64": {
			a:    value.Int32(-5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4i32 == 5u64": {
			a:    value.Int32(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"Int32 5i32 == 5i32": {
			a:    value.Int32(5),
			b:    value.Int32(5).ToValue(),
			want: value.True,
		},
		"Int32 4i32 == 5i32": {
			a:    value.Int32(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"UInt32 5i32 == 5u32": {
			a:    value.Int32(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 -5i32 == 5u32": {
			a:    value.Int32(-5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 4i32 == 5u32": {
			a:    value.Int32(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"Int16 5i32 == 5i16": {
			a:    value.Int32(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4i32 == 5i16": {
			a:    value.Int32(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"UInt16 5i32 == 5u16": {
			a:    value.Int32(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 -5i32 == 5u16": {
			a:    value.Int32(-5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 4i32 == 5u16": {
			a:    value.Int32(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"Int8 5i32 == 5i8": {
			a:    value.Int32(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4i32 == 5i8": {
			a:    value.Int32(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt8 5i32 == 5u8": {
			a:    value.Int32(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 -5i32 == 5u8": {
			a:    value.Int32(-5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 4i32 == 5u8": {
			a:    value.Int32(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5i32 == 5.0f64": {
			a:    value.Int32(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5i32 == 5.5f64": {
			a:    value.Int32(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 5i32 == NaN": {
			a:    value.Int32(5),
			b:    value.Float64NaN().ToValue(),
			want: value.False,
		},
		"Float64 5i32 == +Inf": {
			a:    value.Int32(5),
			b:    value.Float64Inf().ToValue(),
			want: value.False,
		},
		"Float64 5i32 == -Inf": {
			a:    value.Int32(5),
			b:    value.Float64NegInf().ToValue(),
			want: value.False,
		},

		"Float32 5i32 == 5.0f32": {
			a:    value.Int32(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5i32 == 5.5f32": {
			a:    value.Int32(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 5i32 == NaN": {
			a:    value.Int32(5),
			b:    value.Float32NaN().ToValue(),
			want: value.False,
		},
		"Float32 5i32 == +Inf": {
			a:    value.Int32(5),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 5i32 == -Inf": {
			a:    value.Int32(5),
			b:    value.Float32NegInf().ToValue(),
			want: value.False,
		},

		"SmallInt 5i32 == 5": {
			a:    value.Int32(5),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},
		"SmallInt 4i32 == 5": {
			a:    value.Int32(4),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},

		"BigInt 5i32 == 5bi": {
			a:    value.Int32(5),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},
		"BigInt 4i32 == 5bi": {
			a:    value.Int32(4),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},

		"Float 5i32 == 5.0": {
			a:    value.Int32(5),
			b:    value.Float(5).ToValue(),
			want: value.False,
		},
		"Float 5i32 == 5.5": {
			a:    value.Int32(5),
			b:    value.Float(5.5).ToValue(),
			want: value.False,
		},
		"Float 5i32 == +Inf": {
			a:    value.Int32(5),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 5i32 == -Inf": {
			a:    value.Int32(5),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 5i32 == NaN": {
			a:    value.Int32(5),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 5i32 == 5.0bf": {
			a:    value.Int32(5),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.False,
		},
		"BigFloat 5i32 == 5.5bf": {
			a:    value.Int32(5),
			b:    value.Ref(value.NewBigFloat(5.5)),
			want: value.False,
		},
		"BigFloat 5i32 == +Inf": {
			a:    value.Int32(5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 5i32 == -Inf": {
			a:    value.Int32(5),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 5i32 == NaN": {
			a:    value.Int32(5),
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
				t.Fatal(diff)
			}
		})
	}
}

func TestInt32_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Int32
		err  value.Value
	}{
		"perform AND for String and return an error": {
			a:   value.Int32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"perform AND for Int64 and return an error": {
			a:   value.Int32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"11 & 7": {
			a:    value.Int32(0b1011),
			b:    value.Int32(0b0111).ToValue(),
			want: value.Int32(0b0011),
		},
		"-14 & 23": {
			a:    value.Int32(-14),
			b:    value.Int32(23).ToValue(),
			want: value.Int32(18),
		},
		"258 & 0": {
			a:    value.Int32(258),
			b:    value.Int32(0).ToValue(),
			want: value.Int32(0),
		},
		"124 & 255": {
			a:    value.Int32(124),
			b:    value.Int32(255).ToValue(),
			want: value.Int32(124),
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

func TestInt32_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Int32
		err  value.Value
	}{
		"perform OR for String and return an error": {
			a:   value.Int32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"perform OR for Int64 and return an error": {
			a:   value.Int32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"21 | 13": {
			a:    value.Int32(0b10101),
			b:    value.Int32(0b01101).ToValue(),
			want: value.Int32(0b11101),
		},
		"-14 | 23": {
			a:    value.Int32(-14),
			b:    value.Int32(23).ToValue(),
			want: value.Int32(-9),
		},
		"258 | 0": {
			a:    value.Int32(258),
			b:    value.Int32(0).ToValue(),
			want: value.Int32(258),
		},
		"124 | 255": {
			a:    value.Int32(124),
			b:    value.Int32(255).ToValue(),
			want: value.Int32(255),
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

func TestInt32_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Int32
		err  value.Value
	}{
		"perform XOR for String and return an error": {
			a:   value.Int32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"perform XOR for Int64 and return an error": {
			a:   value.Int32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"21 ^ 13": {
			a:    value.Int32(0b10101),
			b:    value.Int32(0b01101).ToValue(),
			want: value.Int32(0b11000),
		},
		"-14 ^ 23": {
			a:    value.Int32(-14),
			b:    value.Int32(23).ToValue(),
			want: value.Int32(-27),
		},
		"258 ^ 0": {
			a:    value.Int32(258),
			b:    value.Int32(0).ToValue(),
			want: value.Int32(258),
		},
		"124 ^ 255": {
			a:    value.Int32(124),
			b:    value.Int32(255).ToValue(),
			want: value.Int32(131),
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

func TestInt32_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.Int32
		b    value.Value
		want value.Int32
		err  value.Value
	}{
		"perform modulo for String and return an error": {
			a:   value.Int32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int32`")),
		},
		"perform modulo for Int64 and return an error": {
			a:   value.Int32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int32`")),
		},
		"21 % 10": {
			a:    value.Int32(21),
			b:    value.Int32(10).ToValue(),
			want: value.Int32(1),
		},
		"38 % 3": {
			a:    value.Int32(38),
			b:    value.Int32(3).ToValue(),
			want: value.Int32(2),
		},
		"522 % 39": {
			a:    value.Int32(522),
			b:    value.Int32(39).ToValue(),
			want: value.Int32(15),
		},
		"38 % 0": {
			a:   value.Int32(38),
			b:   value.Int32(0).ToValue(),
			err: value.Ref(value.NewError(value.ZeroDivisionErrorClass, "cannot divide by zero")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Modulo(tc.b)
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
