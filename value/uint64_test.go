package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestUInt64Inspect(t *testing.T) {
	tests := map[string]struct {
		i    value.UInt64
		want string
	}{
		"positive number": {
			i:    216,
			want: "216u64",
		},
		"zero": {
			i:    0,
			want: "0u64",
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

func TestUInt64_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.UInt64
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.UInt64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"exponentiate Int32 and return an error": {
			a:   value.UInt64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"exponentiate positive UInt64 5 ** 2": {
			a:    value.UInt64(5),
			b:    value.UInt64(2).ToValue(),
			want: value.UInt64(25),
		},
		"exponentiate positive UInt64 7 ** 8": {
			a:    value.UInt64(7),
			b:    value.UInt64(8).ToValue(),
			want: value.UInt64(5764801),
		},
		"exponentiate positive UInt64 2 ** 5": {
			a:    value.UInt64(2),
			b:    value.UInt64(5).ToValue(),
			want: value.UInt64(32),
		},
		"exponentiate positive UInt64 6 ** 1": {
			a:    value.UInt64(6),
			b:    value.UInt64(1).ToValue(),
			want: value.UInt64(6),
		},
		"exponentiate zero": {
			a:    value.UInt64(25),
			b:    value.UInt64(0).ToValue(),
			want: value.UInt64(1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.ExponentiateVal(tc.b)
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
func TestUInt64_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.UInt64
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.UInt64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"add Int32 and return an error": {
			a:   value.UInt64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"add positive UInt64": {
			a:    value.UInt64(53),
			b:    value.UInt64(21).ToValue(),
			want: value.UInt64(74),
		},
		"add UInt64": {
			a:    value.UInt64(25),
			b:    value.UInt64(50).ToValue(),
			want: value.UInt64(75),
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

func TestUInt64_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.UInt64
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.UInt64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"subtract Int32 and return an error": {
			a:   value.UInt64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"subtract positive UInt64": {
			a:    value.UInt64(53),
			b:    value.UInt64(21).ToValue(),
			want: value.UInt64(32),
		},
		"subtract UInt64": {
			a:    value.UInt64(75),
			b:    value.UInt64(50).ToValue(),
			want: value.UInt64(25),
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
func TestUInt64_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.UInt64
		err  value.Value
	}{
		"multiply String and return an error": {
			a:   value.UInt64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"multiply Int32 and return an error": {
			a:   value.UInt64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"multiply positive UInt64": {
			a:    value.UInt64(53),
			b:    value.UInt64(2).ToValue(),
			want: value.UInt64(106),
		},
		"multiply UInt64": {
			a:    value.UInt64(25),
			b:    value.UInt64(2).ToValue(),
			want: value.UInt64(50),
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

func TestUInt64_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.UInt64(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"Int32 and return an error": {
			a:    value.UInt64(5),
			b:    value.Int32(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"15u64 > 30u64": {
			a:    value.UInt64(15),
			b:    value.UInt64(30).ToValue(),
			want: value.False,
		},
		"780u64 > 800u64": {
			a:    value.UInt64(780),
			b:    value.UInt64(800).ToValue(),
			want: value.False,
		},
		"15u64 > 15u64": {
			a:    value.UInt64(15),
			b:    value.UInt64(15).ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanVal(tc.b)
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

func TestUInt64_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.UInt64(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"Int32 and return an error": {
			a:    value.UInt64(5),
			b:    value.Int32(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"15u64 >= 30u64": {
			a:    value.UInt64(15),
			b:    value.UInt64(30).ToValue(),
			want: value.False,
		},
		"780u64 >= 800u64": {
			a:    value.UInt64(780),
			b:    value.UInt64(800).ToValue(),
			want: value.False,
		},
		"15u64 >= 15u64": {
			a:    value.UInt64(15),
			b:    value.UInt64(15).ToValue(),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqualVal(tc.b)
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

func TestUInt64_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.UInt64(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"Int32 and return an error": {
			a:    value.UInt64(5),
			b:    value.Int32(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"15u64 < 30u64": {
			a:    value.UInt64(15),
			b:    value.UInt64(30).ToValue(),
			want: value.True,
		},
		"780u64 < 800u64": {
			a:    value.UInt64(780),
			b:    value.UInt64(800).ToValue(),
			want: value.True,
		},
		"15u64 < 15u64": {
			a:    value.UInt64(15),
			b:    value.UInt64(15).ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanVal(tc.b)
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

func TestUInt64_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.UInt64(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"Int32 and return an error": {
			a:    value.UInt64(5),
			b:    value.Int32(2).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"15u64 <= 30u64": {
			a:    value.UInt64(15),
			b:    value.UInt64(30).ToValue(),
			want: value.True,
		},
		"780u64 <= 800u64": {
			a:    value.UInt64(780),
			b:    value.UInt64(800).ToValue(),
			want: value.True,
		},
		"15u64 <= 15u64": {
			a:    value.UInt64(15),
			b:    value.UInt64(15).ToValue(),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqualVal(tc.b)
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

func TestUInt64_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String 5u64 == '5'": {
			a:    value.UInt64(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},

		"Int64 5u64 == 5i64": {
			a:    value.UInt64(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4u64 == 5i64": {
			a:    value.UInt64(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 5u64 == -5i64": {
			a:    value.UInt64(5),
			b:    value.Int64(-5).ToValue(),
			want: value.False,
		},

		"UInt64 5u64 == 5u64": {
			a:    value.UInt64(5),
			b:    value.UInt64(5).ToValue(),
			want: value.True,
		},
		"UInt64 5u64 == 4u64": {
			a:    value.UInt64(5),
			b:    value.UInt64(4).ToValue(),
			want: value.False,
		},
		"UInt64 4u64 == 5u64": {
			a:    value.UInt64(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"Int32 5u64 == 5i32": {
			a:    value.UInt64(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4u64 == 5i32": {
			a:    value.UInt64(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"UInt32 5u64 == 5u32": {
			a:    value.UInt64(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 5u64 == 4u32": {
			a:    value.UInt64(5),
			b:    value.UInt32(4).ToValue(),
			want: value.False,
		},
		"UInt32 4u64 == 5u32": {
			a:    value.UInt64(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"Int16 5u64 == 5i16": {
			a:    value.UInt64(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4u64 == 5i16": {
			a:    value.UInt64(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"UInt16 5u64 == 5u16": {
			a:    value.UInt64(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 5u64 == 4u16": {
			a:    value.UInt64(5),
			b:    value.UInt16(4).ToValue(),
			want: value.False,
		},
		"UInt16 4u64 == 5u16": {
			a:    value.UInt64(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"Int8 5u64 == 5i8": {
			a:    value.UInt64(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4u64 == 5i8": {
			a:    value.UInt64(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt8 5u64 == 5u8": {
			a:    value.UInt64(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 5u64 == 4u8": {
			a:    value.UInt64(5),
			b:    value.UInt8(4).ToValue(),
			want: value.False,
		},
		"UInt8 4u64 == 5u8": {
			a:    value.UInt64(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5u64 == 5.0f64": {
			a:    value.UInt64(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5u64 == 5.5f64": {
			a:    value.UInt64(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 5u64 == NaN": {
			a:    value.UInt64(5),
			b:    value.Float64NaN().ToValue(),
			want: value.False,
		},
		"Float64 5u64 == +Inf": {
			a:    value.UInt64(5),
			b:    value.Float64Inf().ToValue(),
			want: value.False,
		},
		"Float64 5u64 == -Inf": {
			a:    value.UInt64(5),
			b:    value.Float64NegInf().ToValue(),
			want: value.False,
		},

		"Float32 5u64 == 5.0f32": {
			a:    value.UInt64(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5u64 == 5.5f32": {
			a:    value.UInt64(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 5u64 == NaN": {
			a:    value.UInt64(5),
			b:    value.Float32NaN().ToValue(),
			want: value.False,
		},
		"Float32 5u64 == +Inf": {
			a:    value.UInt64(5),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 5u64 == -Inf": {
			a:    value.UInt64(5),
			b:    value.Float32NegInf().ToValue(),
			want: value.False,
		},

		"SmallInt 5u64 == 5": {
			a:    value.UInt64(5),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},
		"SmallInt 4u64 == 5": {
			a:    value.UInt64(4),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},

		"BigInt 5u64 == 5bi": {
			a:    value.UInt64(5),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},
		"BigInt 4u64 == 5bi": {
			a:    value.UInt64(4),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},

		"Float 5u64 == 5.0": {
			a:    value.UInt64(5),
			b:    value.Float(5).ToValue(),
			want: value.False,
		},
		"Float 5u64 == 5.5": {
			a:    value.UInt64(5),
			b:    value.Float(5.5).ToValue(),
			want: value.False,
		},
		"Float 5u64 == +Inf": {
			a:    value.UInt64(5),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 5u64 == -Inf": {
			a:    value.UInt64(5),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 5u64 == NaN": {
			a:    value.UInt64(5),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 5u64 == 5.0bf": {
			a:    value.UInt64(5),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.False,
		},
		"BigFloat 5u64 == 5.5bf": {
			a:    value.UInt64(5),
			b:    value.Ref(value.NewBigFloat(5.5)),
			want: value.False,
		},
		"BigFloat 5u64 == +Inf": {
			a:    value.UInt64(5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 5u64 == -Inf": {
			a:    value.UInt64(5),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 5u64 == NaN": {
			a:    value.UInt64(5),
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
				t.Fatalf(diff)
			}
		})
	}
}

func TestUInt64_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.UInt64
		err  value.Value
	}{
		"perform AND for String and return an error": {
			a:   value.UInt64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"perform AND for Int32 and return an error": {
			a:   value.UInt64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"11 & 7": {
			a:    value.UInt64(0b1011),
			b:    value.UInt64(0b0111).ToValue(),
			want: value.UInt64(0b0011),
		},
		"124 & 23": {
			a:    value.UInt64(124),
			b:    value.UInt64(23).ToValue(),
			want: value.UInt64(20),
		},
		"258 & 0": {
			a:    value.UInt64(258),
			b:    value.UInt64(0).ToValue(),
			want: value.UInt64(0),
		},
		"124 & 255": {
			a:    value.UInt64(124),
			b:    value.UInt64(255).ToValue(),
			want: value.UInt64(124),
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

func TestUInt64_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.UInt64
		err  value.Value
	}{
		"perform OR for String and return an error": {
			a:   value.UInt64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"perform OR for Int32 and return an error": {
			a:   value.UInt64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"21 | 13": {
			a:    value.UInt64(0b10101),
			b:    value.UInt64(0b01101).ToValue(),
			want: value.UInt64(0b11101),
		},
		"124 | 23": {
			a:    value.UInt64(124),
			b:    value.UInt64(23).ToValue(),
			want: value.UInt64(127),
		},
		"258 | 0": {
			a:    value.UInt64(258),
			b:    value.UInt64(0).ToValue(),
			want: value.UInt64(258),
		},
		"124 | 255": {
			a:    value.UInt64(124),
			b:    value.UInt64(255).ToValue(),
			want: value.UInt64(255),
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

func TestUInt64_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.UInt64
		err  value.Value
	}{
		"perform XOR for String and return an error": {
			a:   value.UInt64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"perform XOR for Int32 and return an error": {
			a:   value.UInt64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"21 ^ 13": {
			a:    value.UInt64(0b10101),
			b:    value.UInt64(0b01101).ToValue(),
			want: value.UInt64(0b11000),
		},
		"124 ^ 23": {
			a:    value.UInt64(124),
			b:    value.UInt64(23).ToValue(),
			want: value.UInt64(107),
		},
		"258 ^ 0": {
			a:    value.UInt64(258),
			b:    value.UInt64(0).ToValue(),
			want: value.UInt64(258),
		},
		"124 ^ 255": {
			a:    value.UInt64(124),
			b:    value.UInt64(255).ToValue(),
			want: value.UInt64(131),
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

func TestUInt64_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt64
		b    value.Value
		want value.UInt64
		err  value.Value
	}{
		"perform modulo for String and return an error": {
			a:   value.UInt64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt64`")),
		},
		"perform modulo for Int32 and return an error": {
			a:   value.UInt64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::UInt64`")),
		},
		"21 % 10": {
			a:    value.UInt64(21),
			b:    value.UInt64(10).ToValue(),
			want: value.UInt64(1),
		},
		"38 % 3": {
			a:    value.UInt64(38),
			b:    value.UInt64(3).ToValue(),
			want: value.UInt64(2),
		},
		"522 % 39": {
			a:    value.UInt64(522),
			b:    value.UInt64(39).ToValue(),
			want: value.UInt64(15),
		},
		"38 % 0": {
			a:   value.UInt64(38),
			b:   value.UInt64(0).ToValue(),
			err: value.Ref(value.NewError(value.ZeroDivisionErrorClass, "cannot divide by zero")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.ModuloVal(tc.b)
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
