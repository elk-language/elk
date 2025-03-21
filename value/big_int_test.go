package value_test

import (
	"testing"

	"math"
	"math/big"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestBigInt_Add(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.NewBigInt(3),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"add SmallInt and return BigInt": {
			a:    value.ParseBigIntPanic("9223372036854775815", 10),
			b:    value.SmallInt(10).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775825", 10)),
		},
		"add SmallInt and return SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775837", 10),
			b:    value.SmallInt(-10).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775827", 10)),
		},
		"add BigInt and return BigInt": {
			a:    value.ParseBigIntPanic("9223372036854775827", 10),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775830", 10)),
		},
		"add Float and return Float": {
			a:    value.NewBigInt(3),
			b:    value.Float(2.5).ToValue(),
			want: value.Float(5.5).ToValue(),
		},
		"add Float NaN and return Float NaN": {
			a:    value.NewBigInt(3),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"add Float -Inf and return Float -Inf": {
			a:    value.NewBigInt(3),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatNegInf().ToValue(),
		},
		"add BigFloat and return BigFloat with 64bit precision": {
			a:    value.NewBigInt(56),
			b:    value.Ref(value.NewBigFloat(2.5)),
			want: value.Ref(value.NewBigFloat(58.5).SetPrecision(64)),
		},
		"add BigFloat NaN and return BigFloat NaN": {
			a:    value.NewBigInt(56),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"add BigFloat +Inf and return BigFloat +Inf": {
			a:    value.NewBigInt(56),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"add BigFloat -Inf and return BigFloat -Inf": {
			a:    value.NewBigInt(56),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatNegInf()),
		},
		"add BigFloat and return BigFloat with 80bit precision": {
			a:    value.NewBigInt(56),
			b:    value.Ref(value.NewBigFloat(2.5).SetPrecision(80)),
			want: value.Ref(value.NewBigFloat(58.5).SetPrecision(80)),
		},
		"add BigFloat and return BigFloat with 65bit precision": {
			a:    value.ParseBigIntPanic("36893488147419103228", 10),
			b:    value.Ref(value.NewBigFloat(2.5).SetPrecision(64)),
			want: value.Ref(value.ParseBigFloatPanic("36893488147419103230.5").SetPrecision(65)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.AddVal(tc.b)
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

func TestBigInt_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.ParseBigIntPanic("9223372036854775817", 10),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"subtract SmallInt and return BigInt": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.SmallInt(5).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775812", 10)),
		},
		"subtract BigInt and return BigInt": {
			a:    value.ParseBigIntPanic("27670116110564327451", 10),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775817", 10)),
			want: value.Ref(value.ParseBigIntPanic("18446744073709551634", 10)),
		},
		"subtract BigInt and return SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775812", 10)),
			want: value.SmallInt(5).ToValue(),
		},
		"subtract Float and return Float": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Float(15.5).ToValue(),
			want: value.Float(9223372036854775801.5).ToValue(),
		},

		"subtract Float NaN and return Float NaN": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"subtract Float +Inf and return Float -Inf": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.FloatInf().ToValue(),
			want: value.FloatNegInf().ToValue(),
		},
		"subtract Float -Inf and return Float +Inf": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},

		"subtract BigFloat NaN and return BigFloat NaN": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"subtract BigFloat +Inf and return BigFloat -Inf": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatNegInf()),
		},
		"subtract BigFloat -Inf and return BigFloat +Inf": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"subtract BigFloat and return BigFloat with 64bit precision": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Ref(value.NewBigFloat(854775817)),
			want: value.Ref(value.NewBigFloat(9223372036000000000).SetPrecision(64)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.SubtractVal(tc.b)
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
func TestBigInt_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"multiply by String and return an error": {
			a:   value.ParseBigIntPanic("9223372036854775817", 10),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"multiply by SmallInt and return BigInt": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.SmallInt(10).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("92233720368547758170", 10)),
		},
		"multiply by BigInt and return BigInt": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775825", 10)),
			want: value.Ref(value.ParseBigIntPanic("85070591730234616105651324816166224025", 10)),
		},
		"multiply by Float and return Float": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Float(0.00001).ToValue(),
			want: value.Float(92233720368547.77).ToValue(),
		},
		"multiply by BigFloat and return BigFloat with 64bit precision": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Ref(value.NewBigFloat(10)),
			want: value.Ref(value.ParseBigFloatPanic("92233720368547758170").SetPrecision(64)),
		},

		"multiply by Float NaN and return Float NaN": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"multiply by Float +Inf and return Float +Inf": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.FloatInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"multiply by Float +Inf and return Float -Inf": {
			a:    value.ParseBigIntPanic("-9223372036854775817", 10),
			b:    value.FloatInf().ToValue(),
			want: value.FloatNegInf().ToValue(),
		},
		"multiply by Float -Inf and return Float -Inf": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatNegInf().ToValue(),
		},
		"multiply by Float -Inf and return Float +Inf": {
			a:    value.ParseBigIntPanic("-9223372036854775817", 10),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},

		"multiply by BigFloat NaN and return BigFloat NaN": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"multiply by BigFloat +Inf and return BigFloat +Inf": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"multiply by BigFloat +Inf and return BigFloat -Inf": {
			a:    value.ParseBigIntPanic("-9223372036854775817", 10),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatNegInf()),
		},
		"multiply by BigFloat -Inf and return BigFloat -Inf": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatNegInf()),
		},
		"multiply by BigFloat -Inf and return BigFloat +Inf": {
			a:    value.ParseBigIntPanic("-9223372036854775817", 10),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatInf()),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.MultiplyVal(tc.b)
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

func TestBigInt_Divide(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"divide by String and return an error": {
			a:   value.ParseBigIntPanic("9223372036854775817", 10),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"divide by SmallInt and return BigInt": {
			a:    value.ParseBigIntPanic("27670116110564327454", 10),
			b:    value.SmallInt(2).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("13835058055282163727", 10)),
		},
		"divide by BigInt and return SmallInt": {
			a:    value.ParseBigIntPanic("27670116110564327454", 10),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775818", 10)),
			want: value.SmallInt(3).ToValue(),
		},
		"divide by Float and return Float": {
			a:    value.ParseBigIntPanic("9223372036854775818", 10),
			b:    value.Float(2).ToValue(),
			want: value.Float(4611686018427387909).ToValue(),
		},
		"divide by BigFloat and return BigFloat with 64bit precision": {
			a:    value.ParseBigIntPanic("1000000000000000000", 10),
			b:    value.Ref(value.NewBigFloat(20000)),
			want: value.Ref(value.NewBigFloat(50000000000000).SetPrecision(64)),
		},

		"divide by Float 0 and return Float +Inf": {
			a:    value.NewBigInt(234),
			b:    value.Float(0).ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"divide by Float 0 and return Float -Inf": {
			a:    value.NewBigInt(-234),
			b:    value.Float(0).ToValue(),
			want: value.FloatNegInf().ToValue(),
		},
		"divide by Float NaN and return Float NaN": {
			a:    value.NewBigInt(234),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"divide by Float +Inf and return Float 0": {
			a:    value.NewBigInt(234),
			b:    value.FloatInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"divide by Float -Inf and return Float -Inf": {
			a:    value.NewBigInt(234),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(0).ToValue(),
		},

		"divide by BigFloat 0 and return BigFloat +Inf": {
			a:    value.NewBigInt(234),
			b:    value.Ref(value.NewBigFloat(0)),
			want: value.Ref(value.BigFloatInf()),
		},
		"divide by BigFloat 0 and return BigFloat -Inf": {
			a:    value.NewBigInt(-234),
			b:    value.Ref(value.NewBigFloat(0)),
			want: value.Ref(value.BigFloatNegInf()),
		},
		"divide by BigFloat NaN and return BigFloat NaN": {
			a:    value.NewBigInt(234),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"divide by BigFloat +Inf and return BigFloat 0": {
			a:    value.NewBigInt(234),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.DivideVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Log(got.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBigInt_IsSmallInt(t *testing.T) {
	tests := map[string]struct {
		i    *value.BigInt
		want bool
	}{
		"fits in SmallInt": {
			i:    value.NewBigInt(math.MaxInt64 - 1),
			want: true,
		},
		"does not fit in SmallInt": {
			i:    value.ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5))),
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.i.IsSmallInt()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBigInt_ToSmallInt(t *testing.T) {
	tests := map[string]struct {
		i    *value.BigInt
		want value.SmallInt
	}{
		"fits in SmallInt": {
			i:    value.NewBigInt(value.MaxSmallInt - 1),
			want: value.SmallInt(value.MaxSmallInt - 1),
		},
		"overflows SmallInt": {
			i:    value.ToElkBigInt((&big.Int{}).Add(big.NewInt(value.MaxSmallInt), big.NewInt(5))),
			want: value.MinSmallInt + 4,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.i.ToSmallInt()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
func TestBigInt_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"exponentiate Int32 and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 5 ** 2": {
			a:    value.NewBigInt(5),
			b:    value.SmallInt(2).ToValue(),
			want: value.SmallInt(25).ToValue(),
		},
		"SmallInt 7 ** 8": {
			a:    value.NewBigInt(7),
			b:    value.SmallInt(8).ToValue(),
			want: value.SmallInt(5764801).ToValue(),
		},
		"SmallInt 2 ** 5": {
			a:    value.NewBigInt(2),
			b:    value.SmallInt(5).ToValue(),
			want: value.SmallInt(32).ToValue(),
		},
		"SmallInt 6 ** 1": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(1).ToValue(),
			want: value.SmallInt(6).ToValue(),
		},
		"SmallInt 2 ** 64": {
			a:    value.NewBigInt(2),
			b:    value.SmallInt(64).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("18446744073709551616", 10)),
		},
		"SmallInt 9223372036854775808 ** 2": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.SmallInt(2).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("85070591730234615865843651857942052864", 10)),
		},
		"SmallInt 4 ** -2": {
			a:    value.NewBigInt(4),
			b:    value.SmallInt(-2).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"SmallInt 25 ** 0": {
			a:    value.NewBigInt(25),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},

		"BigInt 5 ** 2": {
			a:    value.NewBigInt(5),
			b:    value.Ref(value.NewBigInt(2)),
			want: value.SmallInt(25).ToValue(),
		},
		"BigInt 7 ** 8": {
			a:    value.NewBigInt(7),
			b:    value.Ref(value.NewBigInt(8)),
			want: value.SmallInt(5764801).ToValue(),
		},
		"BigInt 2 ** 5": {
			a:    value.NewBigInt(2),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.SmallInt(32).ToValue(),
		},
		"BigInt 6 ** 1": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(1)),
			want: value.SmallInt(6).ToValue(),
		},
		"BigInt 9223372036854775808 ** 2": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Ref(value.NewBigInt(2)),
			want: value.Ref(value.ParseBigIntPanic("85070591730234615865843651857942052864", 10)),
		},
		"BigInt 1 ** 9223372036854775808": {
			a:    value.NewBigInt(1),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigInt 2 ** 64": {
			a:    value.NewBigInt(2),
			b:    value.Ref(value.NewBigInt(64)),
			want: value.Ref(value.ParseBigIntPanic("18446744073709551616", 10)),
		},
		"BigInt 4 ** -2": {
			a:    value.NewBigInt(4),
			b:    value.Ref(value.NewBigInt(-2)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigInt 25 ** 0": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigInt(0)),
			want: value.SmallInt(1).ToValue(),
		},

		"Float 5 ** 2": {
			a:    value.NewBigInt(5),
			b:    value.Float(2).ToValue(),
			want: value.Float(25).ToValue(),
		},
		"Float 7 ** 8": {
			a:    value.NewBigInt(7),
			b:    value.Float(8).ToValue(),
			want: value.Float(5764801).ToValue(),
		},
		"Float 3 ** 2.5": {
			a:    value.NewBigInt(3),
			b:    value.Float(2.5).ToValue(),
			want: value.Float(15.588457268119894).ToValue(),
		},
		"Float 6 ** 1": {
			a:    value.NewBigInt(6),
			b:    value.Float(1).ToValue(),
			want: value.Float(6).ToValue(),
		},
		"Float 9223372036854775808 ** 2": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Float(2).ToValue(),
			want: value.Float(8.507059173023462e+37).ToValue(),
		},
		"Float 4 ** -2": {
			a:    value.NewBigInt(4),
			b:    value.Float(-2).ToValue(),
			want: value.Float(0.0625).ToValue(),
		},
		"Float 25 ** 0": {
			a:    value.NewBigInt(25),
			b:    value.Float(0).ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float 25 ** NaN": {
			a:    value.NewBigInt(25),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    value.NewBigInt(0),
			b:    value.Float(-5).ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    value.NewBigInt(0),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    value.NewBigInt(0),
			b:    value.FloatInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    value.NewBigInt(0),
			b:    value.Float(-8).ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    value.NewBigInt(0),
			b:    value.Float(7).ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    value.NewBigInt(0),
			b:    value.Float(8).ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    value.NewBigInt(-1),
			b:    value.FloatInf().ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    value.NewBigInt(-1),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.NewBigInt(2),
			b:    value.FloatInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.NewBigInt(-2),
			b:    value.FloatInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.NewBigInt(2),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.NewBigInt(-2),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(0).ToValue(),
		},

		"BigFloat 5 ** 2": {
			a:    value.NewBigInt(5),
			b:    value.Ref(value.NewBigFloat(2)),
			want: value.Ref(value.NewBigFloat(25).SetPrecision(64)),
		},
		"BigFloat 7 ** 8": {
			a:    value.NewBigInt(7),
			b:    value.Ref(value.NewBigFloat(8)),
			want: value.Ref(value.NewBigFloat(5764801).SetPrecision(64)),
		},
		"BigFloat 3 ** 2.5": {
			a:    value.NewBigInt(3),
			b:    value.Ref(value.NewBigFloat(2.5)),
			want: value.Ref(value.ParseBigFloatPanic("15.5884572681198956415").SetPrecision(64)),
		},
		"BigFloat 6 ** 1": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(1)),
			want: value.Ref(value.NewBigFloat(6).SetPrecision(64)),
		},
		"BigFloat 4 ** -2": {
			a:    value.NewBigInt(4),
			b:    value.Ref(value.NewBigFloat(-2)),
			want: value.Ref(value.NewBigFloat(0.0625).SetPrecision(64)),
		},
		"BigFloat 25 ** 0": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigFloat(0)),
			want: value.Ref(value.NewBigFloat(1).SetPrecision(64)),
		},
		"BigFloat 25 ** NaN": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    value.NewBigInt(0),
			b:    value.Ref(value.NewBigFloat(-5)),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    value.NewBigInt(0),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    value.NewBigInt(0),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
		"BigFloat 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    value.NewBigInt(0),
			b:    value.Ref(value.NewBigFloat(-8)),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    value.NewBigInt(0),
			b:    value.Ref(value.NewBigFloat(7)),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
		"BigFloat 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    value.NewBigInt(0),
			b:    value.Ref(value.NewBigFloat(8)),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
		"BigFloat -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    value.NewBigInt(-1),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(1).SetPrecision(64)),
		},
		"BigFloat -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    value.NewBigInt(-1),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(1).SetPrecision(64)),
		},
		"BigFloat 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.NewBigInt(2),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.NewBigInt(-2),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.NewBigInt(2),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
		"BigFloat -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.NewBigInt(-2),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.ExponentiateVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Log(got.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBigInt_Compare(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"Int64 and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Int64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 <=> 3": {
			a:    value.NewBigInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"SmallInt 6 <=> 18": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"SmallInt 6 <=> 6": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},

		"BigInt 25 <=> 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigInt 6 <=> 18": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.SmallInt(-1).ToValue(),
		},
		"BigInt 6 <=> 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.SmallInt(0).ToValue(),
		},
		"Float 25 <=> 3": {
			a:    value.NewBigInt(25),
			b:    value.Float(3).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"Float 6 <=> 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"Float 6 <=> 6": {
			a:    value.NewBigInt(6),
			b:    value.Float(6).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"Float 6 <=> Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"Float 6 <=> -Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"Float 6 <=> NaN": {
			a:    value.NewBigInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.Nil,
		},

		"BigFloat 25 <=> 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigFloat 6 <=> 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.SmallInt(-1).ToValue(),
		},
		"BigFloat 6 <=> 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.SmallInt(0).ToValue(),
		},
		"BigFloat 6 <=> Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.SmallInt(-1).ToValue(),
		},
		"BigFloat 6 <=> -Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.SmallInt(1).ToValue(),
		},
		"BigFloat 6 <=> NaN": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.CompareVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBigInt_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.NewBigInt(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"Int64 and return an error": {
			a:    value.NewBigInt(5),
			b:    value.Int64(7).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 > 3": {
			a:    value.NewBigInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.True,
		},
		"SmallInt 6 > 18": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6 > 6": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.False,
		},

		"BigInt 25 > 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.True,
		},
		"BigInt 6 > 18": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6 > 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.False,
		},

		"Float 25 > 3": {
			a:    value.NewBigInt(25),
			b:    value.Float(3).ToValue(),
			want: value.True,
		},
		"Float 6 > 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6 > 6": {
			a:    value.NewBigInt(6),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6 > Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6 > -Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.True,
		},
		"Float 6 > NaN": {
			a:    value.NewBigInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"BigFloat 25 > 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.True,
		},
		"BigFloat 6 > 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6 > 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat 6 > Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6 > -Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.True,
		},
		"BigFloat 6 > NaN": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBigInt_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.NewBigInt(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"Int64 and return an error": {
			a:    value.NewBigInt(5),
			b:    value.Int64(7).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 >= 3": {
			a:    value.NewBigInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.True,
		},
		"SmallInt 6 >= 18": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6 >= 6": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25 >= 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.True,
		},
		"BigInt 6 >= 18": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6 >= 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},

		"Float 25 >= 3": {
			a:    value.NewBigInt(25),
			b:    value.Float(3).ToValue(),
			want: value.True,
		},
		"Float 6 >= 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6 >= 6": {
			a:    value.NewBigInt(6),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6 >= Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6 >= -Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.True,
		},
		"Float 6 >= NaN": {
			a:    value.NewBigInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"BigFloat 25 >= 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.True,
		},
		"BigFloat 6 >= 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6 >= 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat 6 >= Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6 >= -Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.True,
		},
		"BigFloat 6 >= NaN": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqualVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBigInt_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.NewBigInt(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"Int64 and return an error": {
			a:    value.NewBigInt(5),
			b:    value.Int64(7).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 < 3": {
			a:    value.NewBigInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6 < 18": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.True,
		},
		"SmallInt 6 < 6": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.False,
		},

		"BigInt 25 < 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6 < 18": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.True,
		},
		"BigInt 6 < 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.False,
		},

		"Float 25 < 3": {
			a:    value.NewBigInt(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6 < 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.True,
		},
		"Float 6 < 6": {
			a:    value.NewBigInt(6),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6 < Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float 6 < -Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6 < NaN": {
			a:    value.NewBigInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"BigFloat 25 < 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6 < 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.True,
		},
		"BigFloat 6 < 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat 6 < Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat 6 < -Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6 < NaN": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBigInt_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.NewBigInt(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"Int64 and return an error": {
			a:    value.NewBigInt(5),
			b:    value.Int64(7).ToValue(),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 <= 3": {
			a:    value.NewBigInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6 <= 18": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.True,
		},
		"SmallInt 6 <= 6": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25 <= 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6 <= 18": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.True,
		},
		"BigInt 6 <= 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},

		"Float 25 <= 3": {
			a:    value.NewBigInt(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6 <= 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.True,
		},
		"Float 6 <= 6": {
			a:    value.NewBigInt(6),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6 <= Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float 6 <= -Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6 <= NaN": {
			a:    value.NewBigInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"BigFloat 25 <= 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6 <= 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.True,
		},
		"BigFloat 6 <= 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat 6 <= Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat 6 <= -Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6 <= NaN": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqualVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBigInt_LaxEqual(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
	}{
		"String 5 =~ '5'": {
			a:    value.NewBigInt(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},
		"Char 5 =~ `5`": {
			a:    value.NewBigInt(5),
			b:    value.Char('5').ToValue(),
			want: value.False,
		},

		"Int64 5 =~ 5i64": {
			a:    value.NewBigInt(5),
			b:    value.Int64(5).ToValue(),
			want: value.True,
		},
		"Int64 4 =~ 5i64": {
			a:    value.NewBigInt(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},

		"Int32 5 =~ 5i32": {
			a:    value.NewBigInt(5),
			b:    value.Int32(5).ToValue(),
			want: value.True,
		},
		"Int32 4 =~ 5i32": {
			a:    value.NewBigInt(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"Int16 5 =~ 5i16": {
			a:    value.NewBigInt(5),
			b:    value.Int16(5).ToValue(),
			want: value.True,
		},
		"Int16 4 =~ 5i16": {
			a:    value.NewBigInt(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"Int8 5 =~ 5i8": {
			a:    value.NewBigInt(5),
			b:    value.Int8(5).ToValue(),
			want: value.True,
		},
		"Int8 4 =~ 5i8": {
			a:    value.NewBigInt(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt64 5 =~ 5u64": {
			a:    value.NewBigInt(5),
			b:    value.UInt64(5).ToValue(),
			want: value.True,
		},
		"UInt64 4 =~ 5u64": {
			a:    value.NewBigInt(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"UInt32 5 =~ 5u32": {
			a:    value.NewBigInt(5),
			b:    value.UInt32(5).ToValue(),
			want: value.True,
		},
		"UInt32 4 =~ 5u32": {
			a:    value.NewBigInt(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"UInt16 5 =~ 5u16": {
			a:    value.NewBigInt(5),
			b:    value.UInt16(5).ToValue(),
			want: value.True,
		},
		"UInt16 4 =~ 5u16": {
			a:    value.NewBigInt(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"UInt8 5 =~ 5u8": {
			a:    value.NewBigInt(5),
			b:    value.UInt8(5).ToValue(),
			want: value.True,
		},
		"UInt8 4 =~ 5u8": {
			a:    value.NewBigInt(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5 =~ 5f64": {
			a:    value.NewBigInt(5),
			b:    value.Float64(5).ToValue(),
			want: value.True,
		},
		"Float64 5 =~ 5.5f64": {
			a:    value.NewBigInt(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 4 =~ 5f64": {
			a:    value.NewBigInt(4),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float32 5 =~ 5f32": {
			a:    value.NewBigInt(5),
			b:    value.Float32(5).ToValue(),
			want: value.True,
		},
		"Float32 5 =~ 5.5f32": {
			a:    value.NewBigInt(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 4 =~ 5f32": {
			a:    value.NewBigInt(4),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},

		"SmallInt 25 =~ 3": {
			a:    value.NewBigInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6 =~ 18": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6 =~ 6": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25 =~ 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6 =~ 18": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6 =~ 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},

		"Float 25 =~ 3": {
			a:    value.NewBigInt(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6 =~ 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6 =~ 6": {
			a:    value.NewBigInt(6),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6 =~ Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6 =~ -Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6 =~ NaN": {
			a:    value.NewBigInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25 =~ 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6 =~ 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6 =~ 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat 6 =~ Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6 =~ -Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6 =~ NaN": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.LaxEqualVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
		})
	}
}

func TestBigInt_Equal(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
	}{
		"String 5 == '5'": {
			a:    value.NewBigInt(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},
		"Char 5 == `5`": {
			a:    value.NewBigInt(5),
			b:    value.Char('5').ToValue(),
			want: value.False,
		},

		"Int64 5 == 5i64": {
			a:    value.NewBigInt(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4 == 5i64": {
			a:    value.NewBigInt(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},

		"Int32 5 == 5i32": {
			a:    value.NewBigInt(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4 == 5i32": {
			a:    value.NewBigInt(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"Int16 5 == 5i16": {
			a:    value.NewBigInt(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4 == 5i16": {
			a:    value.NewBigInt(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"Int8 5 == 5i8": {
			a:    value.NewBigInt(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4 == 5i8": {
			a:    value.NewBigInt(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt64 5 == 5u64": {
			a:    value.NewBigInt(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4 == 5u64": {
			a:    value.NewBigInt(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"UInt32 5 == 5u32": {
			a:    value.NewBigInt(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 4 == 5u32": {
			a:    value.NewBigInt(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"UInt16 5 == 5u16": {
			a:    value.NewBigInt(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 4 == 5u16": {
			a:    value.NewBigInt(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"UInt8 5 == 5u8": {
			a:    value.NewBigInt(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 4 == 5u8": {
			a:    value.NewBigInt(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5 == 5f64": {
			a:    value.NewBigInt(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5 == 5.5f64": {
			a:    value.NewBigInt(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 4 == 5f64": {
			a:    value.NewBigInt(4),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float32 5 == 5f32": {
			a:    value.NewBigInt(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5 == 5.5f32": {
			a:    value.NewBigInt(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 4 == 5f32": {
			a:    value.NewBigInt(4),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},

		"SmallInt 25 == 3": {
			a:    value.NewBigInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6 == 18": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6 == 6": {
			a:    value.NewBigInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25 == 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6 == 18": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6 == 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},

		"Float 25 == 3": {
			a:    value.NewBigInt(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6 == 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6 == 6": {
			a:    value.NewBigInt(6),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6 == Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6 == -Inf": {
			a:    value.NewBigInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6 == NaN": {
			a:    value.NewBigInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25 == 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6 == 18.5": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6 == 6": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat 6 == Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6 == -Inf": {
			a:    value.NewBigInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6 == NaN": {
			a:    value.NewBigInt(6),
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

func TestBigInt_RightBitshift(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"shift by String and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be used as a bitshift operand")),
		},
		"shift by Float and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Float(2.5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be used as a bitshift operand")),
		},

		"shift by SmallInt 73786976294838206464 >> 3": {
			a:    value.ParseBigIntPanic("73786976294838206464", 10),
			b:    value.SmallInt(3).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"shift by SmallInt 5 >> 1": {
			a:    value.NewBigInt(5),
			b:    value.SmallInt(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by SmallInt 5 >> -1": {
			a:    value.NewBigInt(5),
			b:    value.SmallInt(-1).ToValue(),
			want: value.Ref(value.NewBigInt(10)),
		},
		"shift by SmallInt 75 >> 0": {
			a:    value.NewBigInt(75),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by SmallInt -32 >> 2": {
			a:    value.NewBigInt(-32),
			b:    value.SmallInt(2).ToValue(),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by SmallInt 80 >> 60": {
			a:    value.NewBigInt(80),
			b:    value.SmallInt(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by SmallInt fall down to SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.SmallInt(40).ToValue(),
			want: value.SmallInt(8388608).ToValue(),
		},
		"shift by SmallInt huge result": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.SmallInt(-40).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("10141204801825835211973625643008", 10)),
		},

		"shift by BigInt 73786976294838206464 >> 3": {
			a:    value.ParseBigIntPanic("73786976294838206464", 10),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"shift by BigInt 5 >> 1": {
			a:    value.NewBigInt(5),
			b:    value.Ref(value.NewBigInt(1)),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by BigInt 5 >> -1": {
			a:    value.NewBigInt(5),
			b:    value.Ref(value.NewBigInt(-1)),
			want: value.Ref(value.NewBigInt(10)),
		},
		"shift by BigInt 75 >> 0": {
			a:    value.NewBigInt(75),
			b:    value.Ref(value.NewBigInt(0)),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by BigInt -32 >> 2": {
			a:    value.NewBigInt(-32),
			b:    value.Ref(value.NewBigInt(2)),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by BigInt 80 >> 60": {
			a:    value.NewBigInt(80),
			b:    value.Ref(value.NewBigInt(60)),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by BigInt 80 >> -9223372036854775808": {
			a:    value.NewBigInt(80),
			b:    value.Ref(value.NewBigInt(-9223372036854775808)),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by BigInt fall down to SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Ref(value.NewBigInt(40)),
			want: value.SmallInt(8388608).ToValue(),
		},
		"shift by BigInt huge result": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Ref(value.NewBigInt(-40)),
			want: value.Ref(value.ParseBigIntPanic("10141204801825835211973625643008", 10)),
		},
		"shift by huge BigInt": {
			a:    value.NewBigInt(10),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
			want: value.SmallInt(0).ToValue(),
		},

		"shift by Int64 73786976294838206464 >> 3": {
			a:    value.ParseBigIntPanic("73786976294838206464", 10),
			b:    value.Int64(3).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"shift by Int64 5 >> 1": {
			a:    value.NewBigInt(5),
			b:    value.Int64(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by Int64 5 >> -1": {
			a:    value.NewBigInt(5),
			b:    value.Int64(-1).ToValue(),
			want: value.Ref(value.NewBigInt(10)),
		},
		"shift by Int64 75 >> 0": {
			a:    value.NewBigInt(75),
			b:    value.Int64(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by Int64 -32 >> 2": {
			a:    value.NewBigInt(-32),
			b:    value.Int64(2).ToValue(),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by Int64 80 >> 60": {
			a:    value.NewBigInt(80),
			b:    value.Int64(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by Int64 80 >> -9223372036854775808": {
			a:    value.NewBigInt(80),
			b:    value.Int64(-9223372036854775808).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by Int64 fall down to SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Int64(40).ToValue(),
			want: value.SmallInt(8388608).ToValue(),
		},
		"shift by Int64 huge result": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Int64(-40).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("10141204801825835211973625643008", 10)),
		},

		"shift by Int32 73786976294838206464 >> 3": {
			a:    value.ParseBigIntPanic("73786976294838206464", 10),
			b:    value.Int32(3).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"shift by Int32 5 >> 1": {
			a:    value.NewBigInt(5),
			b:    value.Int32(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by Int32 5 >> -1": {
			a:    value.NewBigInt(5),
			b:    value.Int32(-1).ToValue(),
			want: value.Ref(value.NewBigInt(10)),
		},
		"shift by Int32 75 >> 0": {
			a:    value.NewBigInt(75),
			b:    value.Int32(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by Int32 -32 >> 2": {
			a:    value.NewBigInt(-32),
			b:    value.Int32(2).ToValue(),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by Int32 80 >> 60": {
			a:    value.NewBigInt(80),
			b:    value.Int32(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by Int32 fall down to SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Int32(40).ToValue(),
			want: value.SmallInt(8388608).ToValue(),
		},
		"shift by Int32 huge result": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Int32(-40).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("10141204801825835211973625643008", 10)),
		},

		"shift by Int16 73786976294838206464 >> 3": {
			a:    value.ParseBigIntPanic("73786976294838206464", 10),
			b:    value.Int16(3).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"shift by Int16 5 >> 1": {
			a:    value.NewBigInt(5),
			b:    value.Int16(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by Int16 5 >> -1": {
			a:    value.NewBigInt(5),
			b:    value.Int16(-1).ToValue(),
			want: value.Ref(value.NewBigInt(10)),
		},
		"shift by Int16 75 >> 0": {
			a:    value.NewBigInt(75),
			b:    value.Int16(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by Int16 -32 >> 2": {
			a:    value.NewBigInt(-32),
			b:    value.Int16(2).ToValue(),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by Int16 80 >> 60": {
			a:    value.NewBigInt(80),
			b:    value.Int16(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by Int16 fall down to SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Int16(40).ToValue(),
			want: value.SmallInt(8388608).ToValue(),
		},
		"shift by Int16 huge result": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Int16(-40).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("10141204801825835211973625643008", 10)),
		},

		"shift by Int8 73786976294838206464 >> 3": {
			a:    value.ParseBigIntPanic("73786976294838206464", 10),
			b:    value.Int8(3).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"shift by Int8 5 >> 1": {
			a:    value.NewBigInt(5),
			b:    value.Int8(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by Int8 5 >> -1": {
			a:    value.NewBigInt(5),
			b:    value.Int8(-1).ToValue(),
			want: value.Ref(value.NewBigInt(10)),
		},
		"shift by Int8 75 >> 0": {
			a:    value.NewBigInt(75),
			b:    value.Int8(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by Int8 -32 >> 2": {
			a:    value.NewBigInt(-32),
			b:    value.Int8(2).ToValue(),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by Int8 80 >> 60": {
			a:    value.NewBigInt(80),
			b:    value.Int8(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by Int8 fall down to SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Int8(40).ToValue(),
			want: value.SmallInt(8388608).ToValue(),
		},
		"shift by Int8 huge result": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Int8(-40).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("10141204801825835211973625643008", 10)),
		},

		"shift by UInt64 73786976294838206464 >> 3": {
			a:    value.ParseBigIntPanic("73786976294838206464", 10),
			b:    value.UInt64(3).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"shift by UInt64 5 >> 1": {
			a:    value.NewBigInt(5),
			b:    value.UInt64(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by UInt64 28 >> 2": {
			a:    value.NewBigInt(28),
			b:    value.UInt64(2).ToValue(),
			want: value.SmallInt(7).ToValue(),
		},
		"shift by UInt64 75 >> 0": {
			a:    value.NewBigInt(75),
			b:    value.UInt64(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by UInt64 -32 >> 2": {
			a:    value.NewBigInt(-32),
			b:    value.UInt64(2).ToValue(),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by UInt64 80 >> 60": {
			a:    value.NewBigInt(80),
			b:    value.UInt64(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by UInt64 fall down to SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.UInt64(40).ToValue(),
			want: value.SmallInt(8388608).ToValue(),
		},

		"shift by UInt32 73786976294838206464 >> 3": {
			a:    value.ParseBigIntPanic("73786976294838206464", 10),
			b:    value.UInt32(3).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"shift by UInt32 5 >> 1": {
			a:    value.NewBigInt(5),
			b:    value.UInt32(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by UInt32 28 >> 2": {
			a:    value.NewBigInt(28),
			b:    value.UInt32(2).ToValue(),
			want: value.SmallInt(7).ToValue(),
		},
		"shift by UInt32 75 >> 0": {
			a:    value.NewBigInt(75),
			b:    value.UInt32(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by UInt32 -32 >> 2": {
			a:    value.NewBigInt(-32),
			b:    value.UInt32(2).ToValue(),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by UInt32 80 >> 60": {
			a:    value.NewBigInt(80),
			b:    value.UInt32(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by UInt32 fall down to SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.UInt32(40).ToValue(),
			want: value.SmallInt(8388608).ToValue(),
		},
		"shift by UInt16 73786976294838206464 >> 3": {
			a:    value.ParseBigIntPanic("73786976294838206464", 10),
			b:    value.UInt16(3).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"shift by UInt16 5 >> 1": {
			a:    value.NewBigInt(5),
			b:    value.UInt16(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by UInt16 28 >> 2": {
			a:    value.NewBigInt(28),
			b:    value.UInt16(2).ToValue(),
			want: value.SmallInt(7).ToValue(),
		},
		"shift by UInt16 75 >> 0": {
			a:    value.NewBigInt(75),
			b:    value.UInt16(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by UInt16 -32 >> 2": {
			a:    value.NewBigInt(-32),
			b:    value.UInt16(2).ToValue(),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by UInt16 80 >> 60": {
			a:    value.NewBigInt(80),
			b:    value.UInt16(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by UInt16 fall down to SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.UInt16(40).ToValue(),
			want: value.SmallInt(8388608).ToValue(),
		},

		"shift by UInt8 73786976294838206464 >> 3": {
			a:    value.ParseBigIntPanic("73786976294838206464", 10),
			b:    value.UInt8(3).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"shift by UInt8 5 >> 1": {
			a:    value.NewBigInt(5),
			b:    value.UInt8(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by UInt8 28 >> 2": {
			a:    value.NewBigInt(28),
			b:    value.UInt8(2).ToValue(),
			want: value.SmallInt(7).ToValue(),
		},
		"shift by UInt8 75 >> 0": {
			a:    value.NewBigInt(75),
			b:    value.UInt8(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by UInt8 -32 >> 2": {
			a:    value.NewBigInt(-32),
			b:    value.UInt8(2).ToValue(),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by UInt8 80 >> 60": {
			a:    value.NewBigInt(80),
			b:    value.UInt8(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by UInt8 fall down to SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.UInt8(40).ToValue(),
			want: value.SmallInt(8388608).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.RightBitshiftVal(tc.b)
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

func TestBigInt_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"BigInt & String and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"BigInt & Int32 and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int`")),
		},
		"BigInt & Float and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Float(2.5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Int`")),
		},

		"23 & 10": {
			a:    value.NewBigInt(23),
			b:    value.SmallInt(10).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"11 & 7": {
			a:    value.NewBigInt(11),
			b:    value.SmallInt(7).ToValue(),
			want: value.SmallInt(3).ToValue(),
		},
		"-14 & 23": {
			a:    value.NewBigInt(-14),
			b:    value.SmallInt(23).ToValue(),
			want: value.SmallInt(18).ToValue(),
		},
		"258 & 0": {
			a:    value.NewBigInt(258),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"124 & 255": {
			a:    value.NewBigInt(124),
			b:    value.SmallInt(255).ToValue(),
			want: value.SmallInt(124).ToValue(),
		},

		"255 & 9223372036857247042": {
			a:    value.NewBigInt(255),
			b:    value.Ref(value.ParseBigIntPanic("9223372036857247042", 10)),
			want: value.SmallInt(66).ToValue(),
		},
		"9223372036857247042 & 10223372099998981329": {
			a:    value.ParseBigIntPanic("9223372036857247042", 10),
			b:    value.Ref(value.ParseBigIntPanic("10223372099998981329", 10)),
			want: value.Ref(value.ParseBigIntPanic("9223372036855043136", 10)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseAndVal(tc.b)
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

func TestBigInt_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"BigInt | String and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"BigInt | Int32 and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int`")),
		},
		"BigInt | Float and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Float(2.5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Int`")),
		},

		"23 | 10": {
			a:    value.NewBigInt(23),
			b:    value.SmallInt(10).ToValue(),
			want: value.SmallInt(31).ToValue(),
		},
		"11 | 7": {
			a:    value.NewBigInt(11),
			b:    value.SmallInt(7).ToValue(),
			want: value.SmallInt(15).ToValue(),
		},
		"-14 | 23": {
			a:    value.NewBigInt(-14),
			b:    value.SmallInt(23).ToValue(),
			want: value.SmallInt(-9).ToValue(),
		},
		"258 | 0": {
			a:    value.NewBigInt(258),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(258).ToValue(),
		},
		"124 | 255": {
			a:    value.NewBigInt(124),
			b:    value.SmallInt(255).ToValue(),
			want: value.SmallInt(255).ToValue(),
		},

		"255 | 9223372036857247042": {
			a:    value.NewBigInt(255),
			b:    value.Ref(value.ParseBigIntPanic("9223372036857247042", 10)),
			want: value.Ref(value.ParseBigIntPanic("9223372036857247231", 10)),
		},
		"9223372036857247042 | 10223372099998981329": {
			a:    value.ParseBigIntPanic("9223372036857247042", 10),
			b:    value.Ref(value.ParseBigIntPanic("10223372099998981329", 10)),
			want: value.Ref(value.ParseBigIntPanic("10223372100001185235", 10)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseOrVal(tc.b)
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

func TestBigInt_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"BigInt ^ String and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"BigInt ^ Int32 and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int`")),
		},
		"BigInt ^ Float and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Float(2.5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Int`")),
		},

		"23 ^ 10": {
			a:    value.NewBigInt(23),
			b:    value.SmallInt(10).ToValue(),
			want: value.SmallInt(29).ToValue(),
		},
		"11 ^ 7": {
			a:    value.NewBigInt(11),
			b:    value.SmallInt(7).ToValue(),
			want: value.SmallInt(12).ToValue(),
		},
		"-14 ^ 23": {
			a:    value.NewBigInt(-14),
			b:    value.SmallInt(23).ToValue(),
			want: value.SmallInt(-27).ToValue(),
		},
		"258 ^ 0": {
			a:    value.NewBigInt(258),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(258).ToValue(),
		},
		"124 ^ 255": {
			a:    value.NewBigInt(124),
			b:    value.SmallInt(255).ToValue(),
			want: value.SmallInt(131).ToValue(),
		},

		"255 ^ 9223372036857247042": {
			a:    value.NewBigInt(255),
			b:    value.Ref(value.ParseBigIntPanic("9223372036857247042", 10)),
			want: value.Ref(value.ParseBigIntPanic("9223372036857247165", 10)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseXorVal(tc.b)
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

func TestBigInt_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.NewBigInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 % 3": {
			a:    value.NewBigInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"SmallInt 76 % 6": {
			a:    value.NewBigInt(76),
			b:    value.SmallInt(6).ToValue(),
			want: value.SmallInt(4).ToValue(),
		},
		"SmallInt -76 % 6": {
			a:    value.NewBigInt(-76),
			b:    value.SmallInt(6).ToValue(),
			want: value.SmallInt(-4).ToValue(),
		},
		"SmallInt 76 % -6": {
			a:    value.NewBigInt(76),
			b:    value.SmallInt(-6).ToValue(),
			want: value.SmallInt(4).ToValue(),
		},
		"SmallInt -76 % -6": {
			a:    value.NewBigInt(-76),
			b:    value.SmallInt(-6).ToValue(),
			want: value.SmallInt(-4).ToValue(),
		},
		"SmallInt 124 % 9": {
			a:    value.NewBigInt(124),
			b:    value.SmallInt(9).ToValue(),
			want: value.SmallInt(7).ToValue(),
		},
		"SmallInt 124 % 0": {
			a:   value.NewBigInt(124),
			b:   value.SmallInt(0).ToValue(),
			err: value.Ref(value.NewError(value.ZeroDivisionErrorClass, "cannot divide by zero")),
		},
		"SmallInt 9223372036854775808 % 9": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.SmallInt(9).ToValue(),
			want: value.SmallInt(8).ToValue(),
		},

		"BigInt 25 % 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigInt 76 % 6": {
			a:    value.NewBigInt(76),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.SmallInt(4).ToValue(),
		},
		"BigInt -76 % 6": {
			a:    value.NewBigInt(-76),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.SmallInt(-4).ToValue(),
		},
		"BigInt 76 % -6": {
			a:    value.NewBigInt(76),
			b:    value.Ref(value.NewBigInt(-6)),
			want: value.SmallInt(4).ToValue(),
		},
		"BigInt -76 % -6": {
			a:    value.NewBigInt(-76),
			b:    value.Ref(value.NewBigInt(-6)),
			want: value.SmallInt(-4).ToValue(),
		},
		"BigInt 124 % 9": {
			a:    value.NewBigInt(124),
			b:    value.Ref(value.NewBigInt(9)),
			want: value.SmallInt(7).ToValue(),
		},
		"BigIntInt 124 % 0": {
			a:   value.NewBigInt(124),
			b:   value.Ref(value.NewBigInt(0)),
			err: value.Ref(value.NewError(value.ZeroDivisionErrorClass, "cannot divide by zero")),
		},
		"BigInt 36893488147419103230 % 18446744073709551616": {
			a:    value.ParseBigIntPanic("36893488147419103230", 10),
			b:    value.Ref(value.ParseBigIntPanic("18446744073709551616", 10)),
			want: value.Ref(value.ParseBigIntPanic("18446744073709551614", 10)),
		},
		"BigInt 9765 % 9223372036854775808": {
			a:    value.NewBigInt(9765),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
			want: value.SmallInt(9765).ToValue(),
		},

		"Float 25 % 3": {
			a:    value.NewBigInt(25),
			b:    value.Float(3).ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float 76 % 6": {
			a:    value.NewBigInt(76),
			b:    value.Float(6).ToValue(),
			want: value.Float(4).ToValue(),
		},
		"Float 124 % 9": {
			a:    value.NewBigInt(124),
			b:    value.Float(9).ToValue(),
			want: value.Float(7).ToValue(),
		},
		"Float 124 % +Inf": {
			a:    value.NewBigInt(124),
			b:    value.FloatInf().ToValue(),
			want: value.Float(124).ToValue(),
		},
		"Float 124 % -Inf": {
			a:    value.NewBigInt(124),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(124).ToValue(),
		},
		"Float 74 % 6.25": {
			a:    value.NewBigInt(74),
			b:    value.Float(6.25).ToValue(),
			want: value.Float(5.25).ToValue(),
		},
		"Float -74 % 6.25": {
			a:    value.NewBigInt(-74),
			b:    value.Float(6.25).ToValue(),
			want: value.Float(-5.25).ToValue(),
		},
		"Float 74 % -6.25": {
			a:    value.NewBigInt(74),
			b:    value.Float(-6.25).ToValue(),
			want: value.Float(5.25).ToValue(),
		},
		"Float -74 % -6.25": {
			a:    value.NewBigInt(-74),
			b:    value.Float(-6.25).ToValue(),
			want: value.Float(-5.25).ToValue(),
		},
		"Float 9223372036854775808 % 9.5": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.Float(9.5).ToValue(),
			want: value.Float(8.5).ToValue(),
		},
		"Float 25 % 0": { // Mod(x, 0) = NaN
			a:    value.NewBigInt(25),
			b:    value.Float(0).ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    value.NewBigInt(25),
			b:    value.FloatInf().ToValue(),
			want: value.Float(25).ToValue(),
		},
		"Float -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    value.NewBigInt(-87),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(-87).ToValue(),
		},
		"Float 49 % NaN": { // Mod(x, NaN) = NaN
			a:    value.NewBigInt(49),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},

		"BigFloat 25 % 3": {
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.Ref(value.NewBigFloat(1).SetPrecision(64)),
		},
		"BigFloat 76 % 6": {
			a:    value.NewBigInt(76),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.Ref(value.NewBigFloat(4).SetPrecision(64)),
		},
		"BigFloat 124 % 9": {
			a:    value.NewBigInt(124),
			b:    value.Ref(value.NewBigFloat(9)),
			want: value.Ref(value.NewBigFloat(7).SetPrecision(64)),
		},
		"BigFloat 74 % 6.25": {
			a:    value.NewBigInt(74),
			b:    value.Ref(value.NewBigFloat(6.25)),
			want: value.Ref(value.NewBigFloat(5.25).SetPrecision(64)),
		},
		"BigFloat 74 % 6.25p86": {
			a:    value.NewBigInt(74),
			b:    value.Ref(value.NewBigFloat(6.25).SetPrecision(86)),
			want: value.Ref(value.NewBigFloat(5.25).SetPrecision(86)),
		},
		"BigFloat -74 % 6.25": {
			a:    value.NewBigInt(-74),
			b:    value.Ref(value.NewBigFloat(6.25)),
			want: value.Ref(value.NewBigFloat(-5.25).SetPrecision(64)),
		},
		"BigFloat 74 % -6.25": {
			a:    value.NewBigInt(74),
			b:    value.Ref(value.NewBigFloat(-6.25)),
			want: value.Ref(value.NewBigFloat(5.25).SetPrecision(64)),
		},
		"BigFloat -74 % -6.25": {
			a:    value.NewBigInt(-74),
			b:    value.Ref(value.NewBigFloat(-6.25)),
			want: value.Ref(value.NewBigFloat(-5.25).SetPrecision(64)),
		},
		"BigFloat 25 % 0": { // Mod(x, 0) = NaN
			a:    value.NewBigInt(25),
			b:    value.Ref(value.NewBigFloat(0)),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    value.NewBigInt(25),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(25).SetPrecision(64)),
		},
		"BigFloat -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    value.NewBigInt(-87),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(-87).SetPrecision(64)),
		},
		"BigFloat 49 % NaN": { // Mod(x, NaN) = NaN
			a:    value.NewBigInt(49),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.ModuloVal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
