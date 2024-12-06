package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestUInt8Inspect(t *testing.T) {
	tests := map[string]struct {
		i    value.UInt8
		want string
	}{
		"positive number": {
			i:    216,
			want: "216u8",
		},
		"zero": {
			i:    0,
			want: "0u8",
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

func TestUInt8_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt8
		b    value.Value
		want value.UInt8
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.UInt8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt8`")),
		},
		"exponentiate Int64 and return an error": {
			a:   value.UInt8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt8`")),
		},
		"exponentiate positive Int8 5 ** 2": {
			a:    value.UInt8(5),
			b:    value.UInt8(2).ToValue(),
			want: value.UInt8(25),
		},
		"exponentiate positive Int8 2 ** 5": {
			a:    value.UInt8(2),
			b:    value.UInt8(5).ToValue(),
			want: value.UInt8(32),
		},
		"exponentiate positive Int8 6 ** 1": {
			a:    value.UInt8(6),
			b:    value.UInt8(1).ToValue(),
			want: value.UInt8(6),
		},
		"exponentiate negative Int8": {
			a:    value.UInt8(4),
			b:    value.UInt8(0).ToValue(),
			want: value.UInt8(1),
		},
		"exponentiate zero": {
			a:    value.UInt8(25),
			b:    value.UInt8(0).ToValue(),
			want: value.UInt8(1),
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

func TestUInt8_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt8
		b    value.Value
		want value.UInt8
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.UInt8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt8`")),
		},
		"add Int64 and return an error": {
			a:   value.UInt8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt8`")),
		},
		"add UInt8": {
			a:    value.UInt8(5),
			b:    value.UInt8(2).ToValue(),
			want: value.UInt8(7),
		},
		"add positive UInt8": {
			a:    value.UInt8(53),
			b:    value.UInt8(21).ToValue(),
			want: value.UInt8(74),
		},
		"add UInt8 overflow": {
			a:    value.UInt8(250),
			b:    value.UInt8(50).ToValue(),
			want: value.UInt8(44),
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

func TestUInt8_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt8
		b    value.Value
		want value.UInt8
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.UInt8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt8`")),
		},
		"subtract Int64 and return an error": {
			a:   value.UInt8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt8`")),
		},
		"subtract UInt8": {
			a:    value.UInt8(53),
			b:    value.UInt8(21).ToValue(),
			want: value.UInt8(32),
		},
		"subtract UInt8 underflow": {
			a:    value.UInt8(25),
			b:    value.UInt8(50).ToValue(),
			want: value.UInt8(231),
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

func TestUInt8_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt8
		b    value.Value
		want value.UInt8
		err  value.Value
	}{
		"multiply String and return an error": {
			a:   value.UInt8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt8`")),
		},
		"multiply Int64 and return an error": {
			a:   value.UInt8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt8`")),
		},
		"multiply UInt8": {
			a:    value.UInt8(53),
			b:    value.UInt8(2).ToValue(),
			want: value.UInt8(106),
		},
		"multiply UInt8 overflow": {
			a:    value.UInt8(128),
			b:    value.UInt8(2).ToValue(),
			want: value.UInt8(0),
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

func TestUInt8_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt8
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.UInt8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt8`")),
		},
		"Int64 and return an error": {
			a:   value.UInt8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt8`")),
		},
		"15u8 > 30u8": {
			a:    value.UInt8(15),
			b:    value.UInt8(30).ToValue(),
			want: value.False,
		},
		"78u8 > 80u8": {
			a:    value.UInt8(78),
			b:    value.UInt8(80).ToValue(),
			want: value.False,
		},
		"15u8 > 15u8": {
			a:    value.UInt8(15),
			b:    value.UInt8(15).ToValue(),
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

func TestUInt8_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt8
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.UInt8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt8`")),
		},
		"Int64 and return an error": {
			a:   value.UInt8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt8`")),
		},
		"15u8 >= 30u8": {
			a:    value.UInt8(15),
			b:    value.UInt8(30).ToValue(),
			want: value.False,
		},
		"78u8 >= 80u8": {
			a:    value.UInt8(78),
			b:    value.UInt8(80).ToValue(),
			want: value.False,
		},
		"15u8 >= 15u8": {
			a:    value.UInt8(15),
			b:    value.UInt8(15).ToValue(),
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

func TestUInt8_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt8
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.UInt8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt8`")),
		},
		"Int64 and return an error": {
			a:   value.UInt8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt8`")),
		},
		"15u8 < 30u8": {
			a:    value.UInt8(15),
			b:    value.UInt8(30).ToValue(),
			want: value.True,
		},
		"78u8 < 80u8": {
			a:    value.UInt8(78),
			b:    value.UInt8(80).ToValue(),
			want: value.True,
		},
		"15u8 < 15u8": {
			a:    value.UInt8(15),
			b:    value.UInt8(15).ToValue(),
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

func TestUInt8_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt8
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.UInt8(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::UInt8`")),
		},
		"Int64 and return an error": {
			a:   value.UInt8(5),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::UInt8`")),
		},
		"15u8 <= 30u8": {
			a:    value.UInt8(15),
			b:    value.UInt8(30).ToValue(),
			want: value.True,
		},
		"78u8 <= 80u8": {
			a:    value.UInt8(78),
			b:    value.UInt8(80).ToValue(),
			want: value.True,
		},
		"15u8 <= 15u8": {
			a:    value.UInt8(15),
			b:    value.UInt8(15).ToValue(),
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

func TestUInt8_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.UInt8
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String 5u8 == '5'": {
			a:    value.UInt8(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},

		"Int64 5u8 == 5i64": {
			a:    value.UInt8(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4u8 == 5i64": {
			a:    value.UInt8(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 5u8 == -5i64": {
			a:    value.UInt8(5),
			b:    value.Int64(-5).ToValue(),
			want: value.False,
		},

		"UInt64 5u8 == 5u64": {
			a:    value.UInt8(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4u8 == 5u64": {
			a:    value.UInt8(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"Int32 5u8 == 5i32": {
			a:    value.UInt8(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4u8 == 5i32": {
			a:    value.UInt8(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"UInt32 5u8 == 5u32": {
			a:    value.UInt8(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 4u8 == 5u32": {
			a:    value.UInt8(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"Int16 5u8 == 5i16": {
			a:    value.UInt8(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4u8 == 5i16": {
			a:    value.UInt8(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"UInt16 5u8 == 5u16": {
			a:    value.UInt8(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 4u8 == 5u16": {
			a:    value.UInt8(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"Int8 5u8 == 5i8": {
			a:    value.UInt8(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4u8 == 5i8": {
			a:    value.UInt8(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt8 5u8 == 5u8": {
			a:    value.UInt8(5),
			b:    value.UInt8(5).ToValue(),
			want: value.True,
		},
		"UInt8 4u8 == 5u8": {
			a:    value.UInt8(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5u8 == 5.0f64": {
			a:    value.UInt8(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5u8 == 5.5f64": {
			a:    value.UInt8(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 5u8 == NaN": {
			a:    value.UInt8(5),
			b:    value.Float64NaN().ToValue(),
			want: value.False,
		},
		"Float64 5u8 == +Inf": {
			a:    value.UInt8(5),
			b:    value.Float64Inf().ToValue(),
			want: value.False,
		},
		"Float64 5u8 == -Inf": {
			a:    value.UInt8(5),
			b:    value.Float64NegInf().ToValue(),
			want: value.False,
		},

		"Float32 5u8 == 5.0f32": {
			a:    value.UInt8(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5u8 == 5.5f32": {
			a:    value.UInt8(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 5u8 == NaN": {
			a:    value.UInt8(5),
			b:    value.Float32NaN().ToValue(),
			want: value.False,
		},
		"Float32 5u8 == +Inf": {
			a:    value.UInt8(5),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 5u8 == -Inf": {
			a:    value.UInt8(5),
			b:    value.Float32NegInf().ToValue(),
			want: value.False,
		},

		"SmallInt 5u8 == 5": {
			a:    value.UInt8(5),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},
		"SmallInt 4u8 == 5": {
			a:    value.UInt8(4),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},

		"BigInt 5u8 == 5bi": {
			a:    value.UInt8(5),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},
		"BigInt 4u8 == 5bi": {
			a:    value.UInt8(4),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},

		"Float 5u8 == 5.0": {
			a:    value.UInt8(5),
			b:    value.Float(5).ToValue(),
			want: value.False,
		},
		"Float 5u8 == 5.5": {
			a:    value.UInt8(5),
			b:    value.Float(5.5).ToValue(),
			want: value.False,
		},
		"Float 5u8 == +Inf": {
			a:    value.UInt8(5),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 5u8 == -Inf": {
			a:    value.UInt8(5),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 5u8 == NaN": {
			a:    value.UInt8(5),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 5u8 == 5.0bf": {
			a:    value.UInt8(5),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.False,
		},
		"BigFloat 5u8 == 5.5bf": {
			a:    value.UInt8(5),
			b:    value.Ref(value.NewBigFloat(5.5)),
			want: value.False,
		},
		"BigFloat 5u8 == +Inf": {
			a:    value.UInt8(5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 5u8 == -Inf": {
			a:    value.UInt8(5),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 5u8 == NaN": {
			a:    value.UInt8(5),
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
