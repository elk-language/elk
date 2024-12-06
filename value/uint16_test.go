package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestUInt16Inspect(t *testing.T) {
	tests := map[string]struct {
		i    value.UInt16
		want string
	}{
		"positive number": {
			i:    216,
			want: "216u16",
		},
		"zero": {
			i:    0,
			want: "0u16",
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
func TestUInt16_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.UInt16
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"exponentiate Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"exponentiate positive UInt16 5 ** 2": {
			a:    value.UInt16(5),
			b:    value.UInt16(2).ToValue(),
			want: value.UInt16(25),
		},
		"exponentiate positive UInt16 2 ** 5": {
			a:    value.UInt16(2),
			b:    value.UInt16(5).ToValue(),
			want: value.UInt16(32),
		},
		"exponentiate positive UInt16 6 ** 1": {
			a:    value.UInt16(6),
			b:    value.UInt16(1).ToValue(),
			want: value.UInt16(6),
		},
		"exponentiate zero": {
			a:    value.UInt16(25),
			b:    value.UInt16(0).ToValue(),
			want: value.UInt16(1),
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

func TestUInt16_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.UInt16
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"add Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"add UInt16": {
			a:    value.UInt16(5),
			b:    value.UInt16(2).ToValue(),
			want: value.UInt16(7),
		},
		"add positive UInt16": {
			a:    value.UInt16(53),
			b:    value.UInt16(21).ToValue(),
			want: value.UInt16(74),
		},
		"add zero": {
			a:    value.UInt16(25),
			b:    value.UInt16(0).ToValue(),
			want: value.UInt16(25),
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

func TestUInt16_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.UInt16
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"subtract Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"subtract positive UInt16": {
			a:    value.UInt16(53),
			b:    value.UInt16(21).ToValue(),
			want: value.UInt16(32),
		},
		"subtract zero": {
			a:    value.UInt16(25),
			b:    value.UInt16(0).ToValue(),
			want: value.UInt16(25),
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

func TestUInt16_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.UInt16
		err  value.Value
	}{
		"multiply String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"multiply Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"multiply positive UInt16": {
			a:    value.UInt16(53),
			b:    value.UInt16(2).ToValue(),
			want: value.UInt16(106),
		},
		"multiply by zero": {
			a:    value.UInt16(25),
			b:    value.UInt16(0).ToValue(),
			want: value.UInt16(0),
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
func TestUInt16_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"15u16 > 30u16": {
			a:    value.UInt16(15),
			b:    value.UInt16(30).ToValue(),
			want: value.False,
		},
		"780u16 > 800u16": {
			a:    value.UInt16(780),
			b:    value.UInt16(800).ToValue(),
			want: value.False,
		},
		"15u16 > 15u16": {
			a:    value.UInt16(15),
			b:    value.UInt16(15).ToValue(),
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

func TestUInt16_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"15u16 >= 30u16": {
			a:    value.UInt16(15),
			b:    value.UInt16(30).ToValue(),
			want: value.False,
		},
		"780u16 >= 800u16": {
			a:    value.UInt16(780),
			b:    value.UInt16(800).ToValue(),
			want: value.False,
		},
		"15u16 >= 15u16": {
			a:    value.UInt16(15),
			b:    value.UInt16(15).ToValue(),
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

func TestUInt16_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"15u16 < 30u16": {
			a:    value.UInt16(15),
			b:    value.UInt16(30).ToValue(),
			want: value.True,
		},
		"780u16 < 800u16": {
			a:    value.UInt16(780),
			b:    value.UInt16(800).ToValue(),
			want: value.True,
		},
		"15u16 < 15u16": {
			a:    value.UInt16(15),
			b:    value.UInt16(15).ToValue(),
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
func TestUInt16_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"15u16 <= 30u16": {
			a:    value.UInt16(15),
			b:    value.UInt16(30).ToValue(),
			want: value.True,
		},
		"780u16 <= 800u16": {
			a:    value.UInt16(780),
			b:    value.UInt16(800).ToValue(),
			want: value.True,
		},
		"15u16 <= 15u16": {
			a:    value.UInt16(15),
			b:    value.UInt16(15).ToValue(),
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
func TestUInt16_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String 5u16 == '5'": {
			a:    value.UInt16(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},

		"Int64 5u16 == 5i64": {
			a:    value.UInt16(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4u16 == 5i64": {
			a:    value.UInt16(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 5u16 == -5i64": {
			a:    value.UInt16(5),
			b:    value.Int64(-5).ToValue(),
			want: value.False,
		},

		"UInt64 5u16 == 5u64": {
			a:    value.UInt16(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4u16 == 5u64": {
			a:    value.UInt16(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"Int32 5u16 == 5i32": {
			a:    value.UInt16(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4u16 == 5i32": {
			a:    value.UInt16(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"UInt32 5u16 == 5u32": {
			a:    value.UInt16(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 4u16 == 5u32": {
			a:    value.UInt16(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"Int16 5u16 == 5i16": {
			a:    value.UInt16(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4u16 == 5i16": {
			a:    value.UInt16(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"UInt16 5u16 == 5u16": {
			a:    value.UInt16(5),
			b:    value.UInt16(5).ToValue(),
			want: value.True,
		},
		"UInt16 4u16 == 5u16": {
			a:    value.UInt16(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"Int8 5u16 == 5i8": {
			a:    value.UInt16(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4u16 == 5i8": {
			a:    value.UInt16(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt8 5u16 == 5u8": {
			a:    value.UInt16(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 4u16 == 5u8": {
			a:    value.UInt16(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5u16 == 5.0f64": {
			a:    value.UInt16(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5u16 == 5.5f64": {
			a:    value.UInt16(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 5u16 == NaN": {
			a:    value.UInt16(5),
			b:    value.Float64NaN().ToValue(),
			want: value.False,
		},
		"Float64 5u16 == +Inf": {
			a:    value.UInt16(5),
			b:    value.Float64Inf().ToValue(),
			want: value.False,
		},
		"Float64 5u16 == -Inf": {
			a:    value.UInt16(5),
			b:    value.Float64NegInf().ToValue(),
			want: value.False,
		},

		"Float32 5u16 == 5.0f32": {
			a:    value.UInt16(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5u16 == 5.5f32": {
			a:    value.UInt16(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 5u16 == NaN": {
			a:    value.UInt16(5),
			b:    value.Float32NaN().ToValue(),
			want: value.False,
		},
		"Float32 5u16 == +Inf": {
			a:    value.UInt16(5),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 5u16 == -Inf": {
			a:    value.UInt16(5),
			b:    value.Float32NegInf().ToValue(),
			want: value.False,
		},

		"SmallInt 5u16 == 5": {
			a:    value.UInt16(5),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},
		"SmallInt 4u16 == 5": {
			a:    value.UInt16(4),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},

		"BigInt 5u16 == 5bi": {
			a:    value.UInt16(5),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},
		"BigInt 4u16 == 5bi": {
			a:    value.UInt16(4),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},

		"Float 5u16 == 5.0": {
			a:    value.UInt16(5),
			b:    value.Float(5).ToValue(),
			want: value.False,
		},
		"Float 5u16 == 5.5": {
			a:    value.UInt16(5),
			b:    value.Float(5.5).ToValue(),
			want: value.False,
		},
		"Float 5u16 == +Inf": {
			a:    value.UInt16(5),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 5u16 == -Inf": {
			a:    value.UInt16(5),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 5u16 == NaN": {
			a:    value.UInt16(5),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 5u16 == 5.0bf": {
			a:    value.UInt16(5),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.False,
		},
		"BigFloat 5u16 == 5.5bf": {
			a:    value.UInt16(5),
			b:    value.Ref(value.NewBigFloat(5.5)),
			want: value.False,
		},
		"BigFloat 5u16 == +Inf": {
			a:    value.UInt16(5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 5u16 == -Inf": {
			a:    value.UInt16(5),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 5u16 == NaN": {
			a:    value.UInt16(5),
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

func TestUInt16_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.UInt16
		err  value.Value
	}{
		"perform AND for String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"perform AND for Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"11 & 7": {
			a:    value.UInt16(0b1011),
			b:    value.UInt16(0b0111).ToValue(),
			want: value.UInt16(0b0011),
		},
		"124 & 23": {
			a:    value.UInt16(124),
			b:    value.UInt16(23).ToValue(),
			want: value.UInt16(20),
		},
		"258 & 0": {
			a:    value.UInt16(258),
			b:    value.UInt16(0).ToValue(),
			want: value.UInt16(0),
		},
		"124 & 255": {
			a:    value.UInt16(124),
			b:    value.UInt16(255).ToValue(),
			want: value.UInt16(124),
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
func TestUInt16_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.UInt16
		err  value.Value
	}{
		"perform OR for String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"perform OR for Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"21 | 13": {
			a:    value.UInt16(0b10101),
			b:    value.UInt16(0b01101).ToValue(),
			want: value.UInt16(0b11101),
		},
		"124 | 23": {
			a:    value.UInt16(124),
			b:    value.UInt16(23).ToValue(),
			want: value.UInt16(127),
		},
		"258 | 0": {
			a:    value.UInt16(258),
			b:    value.UInt16(0).ToValue(),
			want: value.UInt16(258),
		},
		"124 | 255": {
			a:    value.UInt16(124),
			b:    value.UInt16(255).ToValue(),
			want: value.UInt16(255),
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

func TestUInt16_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.UInt16
		err  value.Value
	}{
		"perform XOR for String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"perform XOR for Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"21 ^ 13": {
			a:    value.UInt16(0b10101),
			b:    value.UInt16(0b01101).ToValue(),
			want: value.UInt16(0b11000),
		},
		"124 ^ 23": {
			a:    value.UInt16(124),
			b:    value.UInt16(23).ToValue(),
			want: value.UInt16(107),
		},
		"258 ^ 0": {
			a:    value.UInt16(258),
			b:    value.UInt16(0).ToValue(),
			want: value.UInt16(258),
		},
		"124 ^ 255": {
			a:    value.UInt16(124),
			b:    value.UInt16(255).ToValue(),
			want: value.UInt16(131),
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
func TestUInt16_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt16
		b    value.Value
		want value.UInt16
		err  value.Value
	}{
		"perform modulo for String and return an error": {
			a:   value.UInt16(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt16`")),
		},
		"perform modulo for Int64 and return an error": {
			a:   value.UInt16(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt16`")),
		},
		"21 % 10": {
			a:    value.UInt16(21),
			b:    value.UInt16(10).ToValue(),
			want: value.UInt16(1),
		},
		"38 % 3": {
			a:    value.UInt16(38),
			b:    value.UInt16(3).ToValue(),
			want: value.UInt16(2),
		},
		"522 % 39": {
			a:    value.UInt16(522),
			b:    value.UInt16(39).ToValue(),
			want: value.UInt16(15),
		},
		"38 % 0": {
			a:   value.UInt16(38),
			b:   value.UInt16(0).ToValue(),
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
