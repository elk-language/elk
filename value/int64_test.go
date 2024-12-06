package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestInt64Inspect(t *testing.T) {
	tests := map[string]struct {
		i    value.Int64
		want string
	}{
		"positive number": {
			i:    216,
			want: "216i64",
		},
		"negative number": {
			i:    -25,
			want: "-25i64",
		},
		"zero": {
			i:    0,
			want: "0i64",
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

func TestInt64_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.Int64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int64`")),
		},
		"exponentiate Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int64`")),
		},
		"exponentiate positive Int64 5 ** 2": {
			a:    value.Int64(5),
			b:    value.Int64(2).ToValue(),
			want: value.Int64(25),
		},
		"exponentiate positive Int64 7 ** 8": {
			a:    value.Int64(7),
			b:    value.Int64(8).ToValue(),
			want: value.Int64(5764801),
		},
		"exponentiate positive Int64 2 ** 5": {
			a:    value.Int64(2),
			b:    value.Int64(5).ToValue(),
			want: value.Int64(32),
		},
		"exponentiate positive Int64 6 ** 1": {
			a:    value.Int64(6),
			b:    value.Int64(1).ToValue(),
			want: value.Int64(6),
		},
		"exponentiate negative Int64": {
			a:    value.Int64(4),
			b:    value.Int64(-2).ToValue(),
			want: value.Int64(1),
		},
		"exponentiate zero": {
			a:    value.Int64(25),
			b:    value.Int64(0).ToValue(),
			want: value.Int64(1),
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

func TestInt64_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.Int64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int64`")),
		},
		"add Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int64`")),
		},
		"add positive Int64": {
			a:    value.Int64(53),
			b:    value.Int64(21).ToValue(),
			want: value.Int64(74),
		},
		"add negative Int64": {
			a:    value.Int64(25),
			b:    value.Int64(-50).ToValue(),
			want: value.Int64(-25),
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

func TestInt64_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.Int64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int64`")),
		},
		"subtract Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int64`")),
		},
		"subtract positive Int64": {
			a:    value.Int64(53),
			b:    value.Int64(21).ToValue(),
			want: value.Int64(32),
		},
		"subtract negative Int64": {
			a:    value.Int64(25),
			b:    value.Int64(-50).ToValue(),
			want: value.Int64(75),
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

func TestInt64_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Int64
		err  value.Value
	}{
		"multiply String and return an error": {
			a:   value.Int64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int64`")),
		},
		"multiply Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int64`")),
		},
		"multiply positive Int64": {
			a:    value.Int64(53),
			b:    value.Int64(2).ToValue(),
			want: value.Int64(106),
		},
		"multiply negative Int64": {
			a:    value.Int64(25),
			b:    value.Int64(-2).ToValue(),
			want: value.Int64(-50),
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

func TestInt64_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Int64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int64`")),
		},
		"Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int64`")),
		},
		"15i64 > 30i64": {
			a:    value.Int64(15),
			b:    value.Int64(30).ToValue(),
			want: value.False,
		},
		"780i64 > -800i64": {
			a:    value.Int64(780),
			b:    value.Int64(-800).ToValue(),
			want: value.True,
		},
		"15i64 > 15i64": {
			a:    value.Int64(15),
			b:    value.Int64(15).ToValue(),
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

func TestInt64_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Int64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int64`")),
		},
		"Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int64`")),
		},
		"15i64 >= 30i64": {
			a:    value.Int64(15),
			b:    value.Int64(30).ToValue(),
			want: value.False,
		},
		"780i64 >= -800i64": {
			a:    value.Int64(780),
			b:    value.Int64(-800).ToValue(),
			want: value.True,
		},
		"15i64 >= 15i64": {
			a:    value.Int64(15),
			b:    value.Int64(15).ToValue(),
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

func TestInt64_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Int64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int64`")),
		},
		"Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int64`")),
		},
		"15i64 < 30i64": {
			a:    value.Int64(15),
			b:    value.Int64(30).ToValue(),
			want: value.True,
		},
		"780i64 < -800i64": {
			a:    value.Int64(780),
			b:    value.Int64(-800).ToValue(),
			want: value.False,
		},
		"15i64 < 15i64": {
			a:    value.Int64(15),
			b:    value.Int64(15).ToValue(),
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

func TestInt64_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Int64(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int64`")),
		},
		"Int32 and return an error": {
			a:   value.Int64(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int64`")),
		},
		"15i64 <= 30i64": {
			a:    value.Int64(15),
			b:    value.Int64(30).ToValue(),
			want: value.True,
		},
		"780i64 <= -800i64": {
			a:    value.Int64(780),
			b:    value.Int64(-800).ToValue(),
			want: value.False,
		},
		"15i64 <= 15i64": {
			a:    value.Int64(15),
			b:    value.Int64(15).ToValue(),
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

func TestInt64_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.Int64
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String 5i64 == '5'": {
			a:    value.Int64(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},

		"Int64 5i64 == 5i64": {
			a:    value.Int64(5),
			b:    value.Int64(5).ToValue(),
			want: value.True,
		},
		"Int64 4i64 == 5i64": {
			a:    value.Int64(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 5i64 == -5i64": {
			a:    value.Int64(5),
			b:    value.Int64(-5).ToValue(),
			want: value.False,
		},

		"UInt64 5i64 == 5u64": {
			a:    value.Int64(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 -5i64 == 5u64": {
			a:    value.Int64(-5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4i64 == 5u64": {
			a:    value.Int64(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"Int32 5i64 == 5i32": {
			a:    value.Int64(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4i64 == 5i32": {
			a:    value.Int64(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"UInt32 5i64 == 5u32": {
			a:    value.Int64(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 -5i64 == 5u32": {
			a:    value.Int64(-5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 4i64 == 5u32": {
			a:    value.Int64(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"Int16 5i64 == 5i16": {
			a:    value.Int64(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4i64 == 5i16": {
			a:    value.Int64(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"UInt16 5i64 == 5u16": {
			a:    value.Int64(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 -5i64 == 5u16": {
			a:    value.Int64(-5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 4i64 == 5u16": {
			a:    value.Int64(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"Int8 5i64 == 5i8": {
			a:    value.Int64(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4i64 == 5i8": {
			a:    value.Int64(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt8 5i64 == 5u8": {
			a:    value.Int64(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 -5i64 == 5u8": {
			a:    value.Int64(-5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 4i64 == 5u8": {
			a:    value.Int64(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5i64 == 5.0f64": {
			a:    value.Int64(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5i64 == 5.5f64": {
			a:    value.Int64(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 5i64 == NaN": {
			a:    value.Int64(5),
			b:    value.Float64NaN().ToValue(),
			want: value.False,
		},
		"Float64 5i64 == +Inf": {
			a:    value.Int64(5),
			b:    value.Float64Inf().ToValue(),
			want: value.False,
		},
		"Float64 5i64 == -Inf": {
			a:    value.Int64(5),
			b:    value.Float64NegInf().ToValue(),
			want: value.False,
		},

		"Float32 5i64 == 5.0f32": {
			a:    value.Int64(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5i64 == 5.5f32": {
			a:    value.Int64(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 5i64 == NaN": {
			a:    value.Int64(5),
			b:    value.Float32NaN().ToValue(),
			want: value.False,
		},
		"Float32 5i64 == +Inf": {
			a:    value.Int64(5),
			b:    value.Float32Inf().ToValue(),
			want: value.False,
		},
		"Float32 5i64 == -Inf": {
			a:    value.Int64(5),
			b:    value.Float32NegInf().ToValue(),
			want: value.False,
		},

		"SmallInt 5i64 == 5": {
			a:    value.Int64(5),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},
		"SmallInt 4i64 == 5": {
			a:    value.Int64(4),
			b:    value.SmallInt(5).ToValue(),
			want: value.False,
		},

		"BigInt 5i64 == 5bi": {
			a:    value.Int64(5),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},
		"BigInt 4i64 == 5bi": {
			a:    value.Int64(4),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.False,
		},

		"Float 5i64 == 5.0": {
			a:    value.Int64(5),
			b:    value.Float(5).ToValue(),
			want: value.False,
		},
		"Float 5i64 == 5.5": {
			a:    value.Int64(5),
			b:    value.Float(5.5).ToValue(),
			want: value.False,
		},
		"Float 5i64 == +Inf": {
			a:    value.Int64(5),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 5i64 == -Inf": {
			a:    value.Int64(5),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 5i64 == NaN": {
			a:    value.Int64(5),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 5i64 == 5.0bf": {
			a:    value.Int64(5),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.False,
		},
		"BigFloat 5i64 == 5.5bf": {
			a:    value.Int64(5),
			b:    value.Ref(value.NewBigFloat(5.5)),
			want: value.False,
		},
		"BigFloat 5i64 == +Inf": {
			a:    value.Int64(5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 5i64 == -Inf": {
			a:    value.Int64(5),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 5i64 == NaN": {
			a:    value.Int64(5),
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
