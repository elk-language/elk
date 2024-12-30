package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestUInt32Inspect(t *testing.T) {
	tests := map[string]struct {
		i    value.UInt32
		want string
	}{
		"positive number": {
			i:    216,
			want: "216u32",
		},
		"zero": {
			i:    0,
			want: "0u32",
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

func TestUInt32_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.UInt32
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.UInt32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"exponentiate Int64 and return an error": {
			a:   value.UInt32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt32`")),
		},
		"exponentiate positive UInt32 5 ** 2": {
			a:    value.UInt32(5),
			b:    value.UInt32(2).ToValue(),
			want: value.UInt32(25),
		},
		"exponentiate positive UInt32 7 ** 8": {
			a:    value.UInt32(7),
			b:    value.UInt32(8).ToValue(),
			want: value.UInt32(5764801),
		},
		"exponentiate positive UInt32 2 ** 5": {
			a:    value.UInt32(2),
			b:    value.UInt32(5).ToValue(),
			want: value.UInt32(32),
		},
		"exponentiate positive UInt32 6 ** 1": {
			a:    value.UInt32(6),
			b:    value.UInt32(1).ToValue(),
			want: value.UInt32(6),
		},
		"exponentiate zero": {
			a:    value.UInt32(25),
			b:    value.UInt32(0).ToValue(),
			want: value.UInt32(1),
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

func TestUInt32_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.UInt32
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.UInt32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"add Int32": {
			a:   value.UInt32(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt32`")),
		},
		"add UInt32": {
			a:    value.UInt32(53),
			b:    value.UInt32(21).ToValue(),
			want: value.UInt32(74),
		},
		"add zero": {
			a:    value.UInt32(25),
			b:    value.UInt32(0).ToValue(),
			want: value.UInt32(25),
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

func TestUInt32_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.UInt32
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.UInt32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"subtract Int32 and return an error": {
			a:   value.UInt32(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt32`")),
		},
		"subtract positive UInt32": {
			a:    value.UInt32(53),
			b:    value.UInt32(21).ToValue(),
			want: value.UInt32(32),
		},
		"subtract zero": {
			a:    value.UInt32(25),
			b:    value.UInt32(0).ToValue(),
			want: value.UInt32(25),
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

func TestUInt32_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.UInt32
		err  value.Value
	}{
		"multiply String and return an error": {
			a:   value.UInt32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"multiply Int64 and return an error": {
			a:   value.UInt32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt32`")),
		},
		"multiply positive UInt32": {
			a:    value.UInt32(53),
			b:    value.UInt32(2).ToValue(),
			want: value.UInt32(106),
		},
		"multiply zero": {
			a:    value.UInt32(25),
			b:    value.UInt32(0).ToValue(),
			want: value.UInt32(0),
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

func TestUInt32_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.UInt32(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"Int64 and return an error": {
			a:    value.UInt32(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt32`")),
		},
		"15u32 > 30u32": {
			a:    value.UInt32(15),
			b:    value.UInt32(30).ToValue(),
			want: value.False,
		},
		"780u32 > 800u32": {
			a:    value.UInt32(780),
			b:    value.UInt32(800).ToValue(),
			want: value.False,
		},
		"15u32 > 15u32": {
			a:    value.UInt32(15),
			b:    value.UInt32(15).ToValue(),
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
func TestUInt32_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.UInt32(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"Int64 and return an error": {
			a:    value.UInt32(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt32`")),
		},
		"15u32 >= 30u32": {
			a:    value.UInt32(15),
			b:    value.UInt32(30).ToValue(),
			want: value.False,
		},
		"780u32 >= 800u32": {
			a:    value.UInt32(780),
			b:    value.UInt32(800).ToValue(),
			want: value.False,
		},
		"15u32 >= 15u32": {
			a:    value.UInt32(15),
			b:    value.UInt32(15).ToValue(),
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

func TestUInt32_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.UInt32(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"Int64 and return an error": {
			a:    value.UInt32(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt32`")),
		},
		"15u32 < 30u32": {
			a:    value.UInt32(15),
			b:    value.UInt32(30).ToValue(),
			want: value.True,
		},
		"780u32 < 800u32": {
			a:    value.UInt32(780),
			b:    value.UInt32(800).ToValue(),
			want: value.True,
		},
		"15u32 < 15u32": {
			a:    value.UInt32(15),
			b:    value.UInt32(15).ToValue(),
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

func TestUInt32_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.UInt32(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"Int64 and return an error": {
			a:    value.UInt32(5),
			b:    value.Int64(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt32`")),
		},
		"15u32 <= 30u32": {
			a:    value.UInt32(15),
			b:    value.UInt32(30).ToValue(),
			want: value.True,
		},
		"780u32 <= 800u32": {
			a:    value.UInt32(780),
			b:    value.UInt32(800).ToValue(),
			want: value.True,
		},
		"15u32 <= 15u32": {
			a:    value.UInt32(15),
			b:    value.UInt32(15).ToValue(),
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

func TestUInt32_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String 5u32 == '5'": {
			a:    value.UInt32(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},

		"Int64 5u32 == 5i64": {
			a:    value.UInt32(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4u32 == 5i64": {
			a:    value.UInt32(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 5u32 == -5i64": {
			a:    value.UInt32(5),
			b:    value.Int64(-5).ToValue(),
			want: value.False,
		},

		"UInt64 5u32 == 5u64": {
			a:    value.UInt32(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4u32 == 5u64": {
			a:    value.UInt32(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"Int32 5u32 == 5i32": {
			a:    value.UInt32(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4u32 == 5i32": {
			a:    value.UInt32(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"UInt32 5u32 == 5u32": {
			a:    value.UInt32(5),
			b:    value.UInt32(5).ToValue(),
			want: value.True,
		},
		"UInt32 4u32 == 5u32": {
			a:    value.UInt32(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"Int16 5u32 == 5i16": {
			a:    value.UInt32(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4u32 == 5i16": {
			a:    value.UInt32(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"UInt16 5u32 == 5u16": {
			a:    value.UInt32(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 4u32 == 5u16": {
			a:    value.UInt32(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"Int8 5u32 == 5i8": {
			a:    value.UInt32(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4u32 == 5i8": {
			a:    value.UInt32(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt8 5u32 == 5u8": {
			a:    value.UInt32(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 4u32 == 5u8": {
			a:    value.UInt32(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5u32 == 5.0f64": {
			a:    value.UInt32(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5u32 == 5.5f64": {
			a:    value.UInt32(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 5u32 == NaN": {
			a:    value.UInt32(5),
			b:    value.Float64NaN().ToValue(),
			want: value.False,
		},
		"Float64 5u32 == +Inf": {
			a:    value.UInt32(5),
			b:    value.Float64Inf().ToValue(),
			want: value.False,
		},
		"Float64 5u32 == -Inf": {
			a:    value.UInt32(5),
			b:    value.Float64NegInf().ToValue(),
			want: value.False,
		},

		"Float32 5u32 == 5.0f32": {
			a:    value.UInt32(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5u32 == 5.5f32": {
			a:    value.UInt32(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 5u32 == NaN": {
			a:    value.UInt32(5),
			b:    value.Float32NaN().ToValue(),
			want: value.False,
		},
		"Float32 5u32 == +Inf": {
			a:    value.UInt32(5),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 5u32 == -Inf": {
			a:    value.UInt32(5),
			b:    value.Float32NegInf().ToValue(),
			want: value.False,
		},

		"SmallInt 5u32 == 5": {
			a:    value.UInt32(5),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},
		"SmallInt 4u32 == 5": {
			a:    value.UInt32(4),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},

		"BigInt 5u32 == 5bi": {
			a:    value.UInt32(5),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},
		"BigInt 4u32 == 5bi": {
			a:    value.UInt32(4),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},

		"Float 5u32 == 5.0": {
			a:    value.UInt32(5),
			b:    value.Float(5).ToValue(),
			want: value.False,
		},
		"Float 5u32 == 5.5": {
			a:    value.UInt32(5),
			b:    value.Float(5.5).ToValue(),
			want: value.False,
		},
		"Float 5u32 == +Inf": {
			a:    value.UInt32(5),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 5u32 == -Inf": {
			a:    value.UInt32(5),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 5u32 == NaN": {
			a:    value.UInt32(5),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 5u32 == 5.0bf": {
			a:    value.UInt32(5),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.False,
		},
		"BigFloat 5u32 == 5.5bf": {
			a:    value.UInt32(5),
			b:    value.Ref(value.NewBigFloat(5.5)),
			want: value.False,
		},
		"BigFloat 5u32 == +Inf": {
			a:    value.UInt32(5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 5u32 == -Inf": {
			a:    value.UInt32(5),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 5u32 == NaN": {
			a:    value.UInt32(5),
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

func TestUInt32_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.UInt32
		err  value.Value
	}{
		"perform AND for String and return an error": {
			a:   value.UInt32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"perform AND for Int64 and return an error": {
			a:   value.UInt32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt32`")),
		},
		"11 & 7": {
			a:    value.UInt32(0b1011),
			b:    value.UInt32(0b0111).ToValue(),
			want: value.UInt32(0b0011),
		},
		"124 & 23": {
			a:    value.UInt32(124),
			b:    value.UInt32(23).ToValue(),
			want: value.UInt32(20),
		},
		"258 & 0": {
			a:    value.UInt32(258),
			b:    value.UInt32(0).ToValue(),
			want: value.UInt32(0),
		},
		"124 & 255": {
			a:    value.UInt32(124),
			b:    value.UInt32(255).ToValue(),
			want: value.UInt32(124),
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
func TestUInt32_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.UInt32
		err  value.Value
	}{
		"perform OR for String and return an error": {
			a:   value.UInt32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"perform OR for Int64 and return an error": {
			a:   value.UInt32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt32`")),
		},
		"21 | 13": {
			a:    value.UInt32(0b10101),
			b:    value.UInt32(0b01101).ToValue(),
			want: value.UInt32(0b11101),
		},
		"124 | 23": {
			a:    value.UInt32(124),
			b:    value.UInt32(23).ToValue(),
			want: value.UInt32(127),
		},
		"258 | 0": {
			a:    value.UInt32(258),
			b:    value.UInt32(0).ToValue(),
			want: value.UInt32(258),
		},
		"124 | 255": {
			a:    value.UInt32(124),
			b:    value.UInt32(255).ToValue(),
			want: value.UInt32(255),
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

func TestUInt32_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.UInt32
		err  value.Value
	}{
		"perform XOR for String and return an error": {
			a:   value.UInt32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"perform XOR for Int64 and return an error": {
			a:   value.UInt32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt32`")),
		},
		"21 ^ 13": {
			a:    value.UInt32(0b10101),
			b:    value.UInt32(0b01101).ToValue(),
			want: value.UInt32(0b11000),
		},
		"124 ^ 23": {
			a:    value.UInt32(124),
			b:    value.UInt32(23).ToValue(),
			want: value.UInt32(107),
		},
		"258 ^ 0": {
			a:    value.UInt32(258),
			b:    value.UInt32(0).ToValue(),
			want: value.UInt32(258),
		},
		"124 ^ 255": {
			a:    value.UInt32(124),
			b:    value.UInt32(255).ToValue(),
			want: value.UInt32(131),
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

func TestUInt32_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt32
		b    value.Value
		want value.UInt32
		err  value.Value
	}{
		"perform modulo for String and return an error": {
			a:   value.UInt32(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt32`")),
		},
		"perform modulo for Int64 and return an error": {
			a:   value.UInt32(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt32`")),
		},
		"21 % 10": {
			a:    value.UInt32(21),
			b:    value.UInt32(10).ToValue(),
			want: value.UInt32(1),
		},
		"38 % 3": {
			a:    value.UInt32(38),
			b:    value.UInt32(3).ToValue(),
			want: value.UInt32(2),
		},
		"522 % 39": {
			a:    value.UInt32(522),
			b:    value.UInt32(39).ToValue(),
			want: value.UInt32(15),
		},
		"38 % 0": {
			a:   value.UInt32(38),
			b:   value.UInt32(0).ToValue(),
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
