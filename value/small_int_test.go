package value_test

import (
	"math/big"
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestSmallInt_Add(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"add String and return an error": {
			a:   value.SmallInt(3),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewCoerceError(value.IntClass, value.StringClass)),
		},
		"add SmallInt and return SmallInt": {
			a:    value.SmallInt(3),
			b:    value.SmallInt(10).ToValue(),
			want: value.SmallInt(13).ToValue(),
		},
		"add Float and return Float": {
			a:    value.SmallInt(-20),
			b:    value.Float(2.5).ToValue(),
			want: value.Float(-17.5).ToValue(),
		},
		"add Float NaN and return Float NaN": {
			a:    value.SmallInt(-20),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"add Float +Inf and return Float +Inf": {
			a:    value.SmallInt(-20),
			b:    value.FloatInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"add Float -Inf and return Float -Inf": {
			a:    value.SmallInt(-20),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatNegInf().ToValue(),
		},
		"add BigFloat NaN and return BigFloat NaN": {
			a:    value.SmallInt(56),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"add BigFloat +Inf and return BigFloat +Inf": {
			a:    value.SmallInt(56),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"add BigFloat -Inf and return BigFloat -Inf": {
			a:    value.SmallInt(56),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatNegInf()),
		},
		"add BigFloat and return BigFloat with 64bit precision": {
			a:    value.SmallInt(56),
			b:    value.Ref(value.NewBigFloat(2.5)),
			want: value.Ref(value.NewBigFloat(58.5).SetPrecision(64)),
		},
		"add BigFloat and return BigFloat with 80bit precision": {
			a:    value.SmallInt(56),
			b:    value.Ref(value.NewBigFloat(2.5).SetPrecision(80)),
			want: value.Ref(value.ToElkBigFloat((&big.Float{}).SetPrec(80).Add(big.NewFloat(56), big.NewFloat(2.5)))),
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

func TestSmallInt_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"subtract String and return an error": {
			a:   value.SmallInt(3),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewCoerceError(value.IntClass, value.StringClass)),
		},
		"subtract SmallInt and return SmallInt": {
			a:    value.SmallInt(3),
			b:    value.SmallInt(10).ToValue(),
			want: value.SmallInt(-7).ToValue(),
		},

		"subtract Float and return Float": {
			a:    value.SmallInt(-20),
			b:    value.Float(2.5).ToValue(),
			want: value.Float(-22.5).ToValue(),
		},
		"subtract Float NaN and return Float NaN": {
			a:    value.SmallInt(26),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"subtract Float +Inf and return Float -Inf": {
			a:    value.SmallInt(19),
			b:    value.FloatInf().ToValue(),
			want: value.FloatNegInf().ToValue(),
		},
		"subtract Float -Inf and return Float +Inf": {
			a:    value.SmallInt(3),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},

		"subtract BigFloat and return BigFloat with 64bit precision": {
			a:    value.SmallInt(56),
			b:    value.Ref(value.NewBigFloat(2.5)),
			want: value.Ref(value.NewBigFloat(53.5).SetPrecision(64)),
		},
		"subtract BigFloat and return BigFloat with 80bit precision": {
			a:    value.SmallInt(56),
			b:    value.Ref(value.NewBigFloat(2.5).SetPrecision(80)),
			want: value.Ref(value.ToElkBigFloat((&big.Float{}).SetPrec(80).Sub(big.NewFloat(56), big.NewFloat(2.5)))),
		},

		"subtract BigFloat NaN and return BigFloat NaN": {
			a:    value.SmallInt(35),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"subtract BigFloat +Inf and return BigFloat -Inf": {
			a:    value.SmallInt(56),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatNegInf()),
		},
		"subtract BigFloat -Inf and return BigFloat +Inf": {
			a:    value.SmallInt(-12),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatInf()),
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

func TestSmallInt_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"multiply by String and return an error": {
			a:   value.SmallInt(3),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewCoerceError(value.IntClass, value.StringClass)),
		},
		"multiply by SmallInt and return SmallInt": {
			a:    value.SmallInt(3),
			b:    value.SmallInt(10).ToValue(),
			want: value.SmallInt(30).ToValue(),
		},
		"multiply by Float and return Float": {
			a:    value.SmallInt(-20),
			b:    value.Float(2.5).ToValue(),
			want: value.Float(-50).ToValue(),
		},
		"multiply by BigFloat and return BigFloat with 64bit precision": {
			a:    value.SmallInt(56),
			b:    value.Ref(value.NewBigFloat(2.5)),
			want: value.Ref(value.NewBigFloat(140).SetPrecision(64)),
		},
		"multiply by BigFloat and return BigFloat with 80bit precision": {
			a:    value.SmallInt(56),
			b:    value.Ref(value.NewBigFloat(2.5).SetPrecision(80)),
			want: value.Ref(value.ToElkBigFloat((&big.Float{}).SetPrec(80).Mul(big.NewFloat(56), big.NewFloat(2.5)))),
		},

		"multiply by Float NaN and return Float NaN": {
			a:    value.SmallInt(234),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"multiply by Float +Inf and return Float +Inf": {
			a:    value.SmallInt(234),
			b:    value.FloatInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"multiply by Float +Inf and return Float -Inf": {
			a:    value.SmallInt(-123),
			b:    value.FloatInf().ToValue(),
			want: value.FloatNegInf().ToValue(),
		},
		"multiply by Float -Inf and return Float -Inf": {
			a:    value.SmallInt(56),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatNegInf().ToValue(),
		},
		"multiply by Float -Inf and return Float +Inf": {
			a:    value.SmallInt(-5),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},

		"multiply by BigFloat NaN and return BigFloat NaN": {
			a:    value.SmallInt(75),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"multiply by BigFloat +Inf and return BigFloat +Inf": {
			a:    value.SmallInt(15),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"multiply by BigFloat +Inf and return BigFloat -Inf": {
			a:    value.SmallInt(-2),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatNegInf()),
		},
		"multiply by BigFloat -Inf and return BigFloat -Inf": {
			a:    value.SmallInt(7),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatNegInf()),
		},
		"multiply by BigFloat -Inf and return BigFloat +Inf": {
			a:    value.SmallInt(-9),
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

func TestSmallInt_Divide(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"divide by String and return an error": {
			a:   value.SmallInt(3),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewCoerceError(value.IntClass, value.StringClass)),
		},
		"divide by SmallInt and return SmallInt": {
			a:    value.SmallInt(30),
			b:    value.SmallInt(10).ToValue(),
			want: value.SmallInt(3).ToValue(),
		},

		"Float -20 / 0.5": {
			a:    value.SmallInt(-20),
			b:    value.Float(0.5).ToValue(),
			want: value.Float(-40).ToValue(),
		},
		"Float 234 / 0": {
			a:    234,
			b:    value.Float(0).ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float -234 / 0": {
			a:    -234,
			b:    value.Float(0).ToValue(),
			want: value.FloatNegInf().ToValue(),
		},
		"Float 234 / NaN": {
			a:    234,
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float 234 / +Inf": {
			a:    234,
			b:    value.FloatInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float 56 / -Inf": {
			a:    56,
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(0).ToValue(),
		},

		"BigFloat 55 / 2.5 with 64bit precision": {
			a:    value.SmallInt(55),
			b:    value.Ref(value.NewBigFloat(2.5)),
			want: value.Ref(value.NewBigFloat(22.0).SetPrecision(64)),
		},
		"BigFloat 55 / 2.5 with 80bit precision": {
			a:    value.SmallInt(55),
			b:    value.Ref(value.NewBigFloat(2.5).SetPrecision(80)),
			want: value.Ref(value.NewBigFloat(22).SetPrecision(80)),
		},
		"BigFloat 234 / 0": {
			a:    234,
			b:    value.Ref(value.NewBigFloat(0)),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat -234 / 0": {
			a:    -234,
			b:    value.Ref(value.NewBigFloat(0)),
			want: value.Ref(value.BigFloatNegInf()),
		},
		"BigFloat 234 / NaN": {
			a:    234,
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat 234 / +Inf": {
			a:    234,
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.DivideVal(tc.b)
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

func TestSmallInt_addOverflow(t *testing.T) {
	tests := map[string]struct {
		a, b value.SmallInt
		want value.SmallInt
		ok   bool
	}{
		"not overflow": {
			a:    value.SmallInt(3),
			b:    value.SmallInt(10),
			want: value.SmallInt(13),
			ok:   true,
		},
		"overflow": {
			a:    value.SmallInt(10),
			b:    value.SmallInt(value.MaxSmallInt),
			want: value.SmallInt(value.MinSmallInt + 9),
			ok:   false,
		},
		"not underflow": {
			a:    value.SmallInt(value.MinSmallInt + 20),
			b:    value.SmallInt(-18),
			want: value.SmallInt(value.MinSmallInt + 2),
			ok:   true,
		},
		"not underflow positive": {
			a:    value.SmallInt(value.MinSmallInt + 20),
			b:    value.SmallInt(18),
			want: value.SmallInt(value.MinSmallInt + 38),
			ok:   true,
		},
		"underflow": {
			a:    value.SmallInt(value.MinSmallInt),
			b:    value.SmallInt(-20),
			want: value.SmallInt(value.MaxSmallInt - 19),
			ok:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.a.AddOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSmallInt_subtractOverflow(t *testing.T) {
	tests := map[string]struct {
		a, b value.SmallInt
		want value.SmallInt
		ok   bool
	}{
		"not underflow": {
			a:    value.SmallInt(3),
			b:    value.SmallInt(10),
			want: value.SmallInt(-7),
			ok:   true,
		},
		"underflow": {
			a:    value.SmallInt(value.MinSmallInt),
			b:    value.SmallInt(3),
			want: value.SmallInt(value.MaxSmallInt - 2),
			ok:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.a.SubtractOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSmallInt_multiplyOverflow(t *testing.T) {
	tests := map[string]struct {
		a, b value.SmallInt
		want value.SmallInt
		ok   bool
	}{
		"not overflow": {
			a:    value.SmallInt(3),
			b:    value.SmallInt(10),
			want: value.SmallInt(30),
			ok:   true,
		},
		"overflow": {
			a:    value.SmallInt(10),
			b:    value.SmallInt(value.MaxSmallInt),
			want: value.SmallInt(-10),
			ok:   false,
		},
		"not underflow": {
			a:    value.SmallInt(-125),
			b:    value.SmallInt(5),
			want: value.SmallInt(-625),
			ok:   true,
		},
		"underflow": {
			a:    value.SmallInt(value.MinSmallInt),
			b:    value.SmallInt(2),
			want: value.SmallInt(0),
			ok:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.a.MultiplyOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSmallInt_divideOverflow(t *testing.T) {
	tests := map[string]struct {
		a, b value.SmallInt
		want value.SmallInt
		ok   bool
	}{
		"not overflow": {
			a:    value.SmallInt(20),
			b:    value.SmallInt(5),
			want: value.SmallInt(4),
			ok:   true,
		},
		"overflow": {
			a:    value.SmallInt(value.MinSmallInt),
			b:    value.SmallInt(-1),
			want: value.SmallInt(value.MinSmallInt),
			ok:   false,
		},
		"not underflow": {
			a:    value.SmallInt(-625),
			b:    value.SmallInt(5),
			want: value.SmallInt(-125),
			ok:   true,
		},
		"division by zero": {
			a:    value.SmallInt(value.MinSmallInt),
			b:    value.SmallInt(0),
			want: value.SmallInt(0),
			ok:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.a.DivideOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSmallInt_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.SmallInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"exponentiate Int32 and return an error": {
			a:   value.SmallInt(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 5 ** 2": {
			a:    value.SmallInt(5),
			b:    value.SmallInt(2).ToValue(),
			want: value.SmallInt(25).ToValue(),
		},
		"SmallInt 7 ** 8": {
			a:    value.SmallInt(7),
			b:    value.SmallInt(8).ToValue(),
			want: value.SmallInt(5764801).ToValue(),
		},
		"SmallInt 2 ** 5": {
			a:    value.SmallInt(2),
			b:    value.SmallInt(5).ToValue(),
			want: value.SmallInt(32).ToValue(),
		},
		"SmallInt 6 ** 1": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(1).ToValue(),
			want: value.SmallInt(6).ToValue(),
		},
		"SmallInt 2 ** 64": {
			a:    value.SmallInt(2),
			b:    value.SmallInt(64).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("18446744073709551616", 10)),
		},
		"SmallInt 4 ** -2": {
			a:    value.SmallInt(4),
			b:    value.SmallInt(-2).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"SmallInt 25 ** 0": {
			a:    value.SmallInt(25),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},

		"BigInt 5 ** 2": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.NewBigInt(2)),
			want: value.SmallInt(25).ToValue(),
		},
		"BigInt 7 ** 8": {
			a:    value.SmallInt(7),
			b:    value.Ref(value.NewBigInt(8)),
			want: value.SmallInt(5764801).ToValue(),
		},
		"BigInt 2 ** 5": {
			a:    value.SmallInt(2),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.SmallInt(32).ToValue(),
		},
		"BigInt 6 ** 1": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(1)),
			want: value.SmallInt(6).ToValue(),
		},
		"BigInt 2 ** 64": {
			a:    value.SmallInt(2),
			b:    value.Ref(value.NewBigInt(64)),
			want: value.Ref(value.ParseBigIntPanic("18446744073709551616", 10)),
		},
		"BigInt 4 ** -2": {
			a:    value.SmallInt(4),
			b:    value.Ref(value.NewBigInt(-2)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigInt 25 ** 0": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigInt(0)),
			want: value.SmallInt(1).ToValue(),
		},

		"Float 5 ** 2": {
			a:    value.SmallInt(5),
			b:    value.Float(2).ToValue(),
			want: value.Float(25).ToValue(),
		},
		"Float 7 ** 8": {
			a:    value.SmallInt(7),
			b:    value.Float(8).ToValue(),
			want: value.Float(5764801).ToValue(),
		},
		"Float 3 ** 2.5": {
			a:    value.SmallInt(3),
			b:    value.Float(2.5).ToValue(),
			want: value.Float(15.588457268119894).ToValue(),
		},
		"Float 6 ** 1": {
			a:    value.SmallInt(6),
			b:    value.Float(1).ToValue(),
			want: value.Float(6).ToValue(),
		},
		"Float 4 ** -2": {
			a:    value.SmallInt(4),
			b:    value.Float(-2).ToValue(),
			want: value.Float(0.0625).ToValue(),
		},
		"Float 25 ** 0": {
			a:    value.SmallInt(25),
			b:    value.Float(0).ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float 25 ** NaN": {
			a:    value.SmallInt(25),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    value.SmallInt(0),
			b:    value.Float(-5).ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    value.SmallInt(0),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    value.SmallInt(0),
			b:    value.FloatInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    value.SmallInt(0),
			b:    value.Float(-8).ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    value.SmallInt(0),
			b:    value.Float(7).ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    value.SmallInt(0),
			b:    value.Float(8).ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    value.SmallInt(-1),
			b:    value.FloatInf().ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    value.SmallInt(-1),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.SmallInt(2),
			b:    value.FloatInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.SmallInt(-2),
			b:    value.FloatInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.SmallInt(2),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.SmallInt(-2),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"BigFloat 5 ** 2": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.NewBigFloat(2)),
			want: value.Ref(value.NewBigFloat(25).SetPrecision(64)),
		},
		"BigFloat 7 ** 8": {
			a:    value.SmallInt(7),
			b:    value.Ref(value.NewBigFloat(8)),
			want: value.Ref(value.NewBigFloat(5764801).SetPrecision(64)),
		},
		"BigFloat 3 ** 2.5": {
			a:    value.SmallInt(3),
			b:    value.Ref(value.NewBigFloat(2.5)),
			want: value.Ref(value.ParseBigFloatPanic("15.5884572681198956415").SetPrecision(64)),
		},
		"BigFloat 6 ** 1": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(1)),
			want: value.Ref(value.NewBigFloat(6).SetPrecision(64)),
		},
		"BigFloat 4 ** -2": {
			a:    value.SmallInt(4),
			b:    value.Ref(value.NewBigFloat(-2)),
			want: value.Ref(value.NewBigFloat(0.0625).SetPrecision(64)),
		},
		"BigFloat 25 ** 0": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigFloat(0)),
			want: value.Ref(value.NewBigFloat(1).SetPrecision(64)),
		},
		"BigFloat 25 ** NaN": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    value.SmallInt(0),
			b:    value.Ref(value.NewBigFloat(-5)),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    value.SmallInt(0),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    value.SmallInt(0),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
		"BigFloat 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    value.SmallInt(0),
			b:    value.Ref(value.NewBigFloat(-8)),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    value.SmallInt(0),
			b:    value.Ref(value.NewBigFloat(7)),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
		"BigFloat 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    value.SmallInt(0),
			b:    value.Ref(value.NewBigFloat(8)),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
		"BigFloat -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    value.SmallInt(-1),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(1).SetPrecision(64)),
		},
		"BigFloat -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    value.SmallInt(-1),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(1).SetPrecision(64)),
		},
		"BigFloat 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.SmallInt(2),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.SmallInt(-2),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.SmallInt(2),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
		"BigFloat -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.SmallInt(-2),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(0).SetPrecision(64)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.ExponentiateVal(tc.b)
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

func TestSmallInt_RightBitshift(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"shift by String and return an error": {
			a:   value.SmallInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be used as a bitshift operand")),
		},
		"shift by Float and return an error": {
			a:   value.SmallInt(5),
			b:   value.Float(2.5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be used as a bitshift operand")),
		},

		"shift by SmallInt 5 >> 1": {
			a:    value.SmallInt(5),
			b:    value.SmallInt(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by SmallInt 5 >> -1": {
			a:    value.SmallInt(5),
			b:    value.SmallInt(-1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by SmallInt 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by SmallInt -32 >> 2": {
			a:    value.SmallInt(-32),
			b:    value.SmallInt(2).ToValue(),
			want: value.SmallInt(-8).ToValue(),
		},
		"shift by SmallInt 80 >> 60": {
			a:    value.SmallInt(80),
			b:    value.SmallInt(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},

		"shift by BigInt 5 >> 1": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.NewBigInt(1)),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by BigInt 5 >> -1": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.NewBigInt(-1)),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by BigInt 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.Ref(value.NewBigInt(0)),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by BigInt 80 >> 60": {
			a:    value.SmallInt(80),
			b:    value.Ref(value.NewBigInt(60)),
			want: value.SmallInt(0).ToValue(),
		},

		"shift by Int64 5 >> 1": {
			a:    value.SmallInt(5),
			b:    value.Int64(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by Int64 5 >> -1": {
			a:    value.SmallInt(5),
			b:    value.Int64(-1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by Int64 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.Int64(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by Int64 80 >> 60": {
			a:    value.SmallInt(80),
			b:    value.Int64(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},

		"shift by Int32 5 >> 1": {
			a:    value.SmallInt(5),
			b:    value.Int32(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by Int32 5 >> -1": {
			a:    value.SmallInt(5),
			b:    value.Int32(-1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by Int32 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.Int32(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by Int32 80 >> 60": {
			a:    value.SmallInt(80),
			b:    value.Int32(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by Int16 5 >> 1": {
			a:    value.SmallInt(5),
			b:    value.Int16(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by Int16 5 >> -1": {
			a:    value.SmallInt(5),
			b:    value.Int16(-1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by Int16 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.Int16(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by Int16 80 >> 60": {
			a:    value.SmallInt(80),
			b:    value.Int16(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},

		"shift by Int8 5 >> 1": {
			a:    value.SmallInt(5),
			b:    value.Int8(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by Int8 5 >> -1": {
			a:    value.SmallInt(5),
			b:    value.Int8(-1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by Int8 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.Int8(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by Int8 80 >> 60": {
			a:    value.SmallInt(80),
			b:    value.Int8(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by UInt64 5 >> 1": {
			a:    value.SmallInt(5),
			b:    value.UInt64(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by UInt64 28 >> 2": {
			a:    value.SmallInt(28),
			b:    value.UInt64(2).ToValue(),
			want: value.SmallInt(7).ToValue(),
		},
		"shift by UInt64 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.UInt64(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by UInt64 80 >> 60": {
			a:    value.SmallInt(80),
			b:    value.UInt64(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},

		"shift by UInt32 5 >> 1": {
			a:    value.SmallInt(5),
			b:    value.UInt32(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by UInt32 28 >> 2": {
			a:    value.SmallInt(28),
			b:    value.UInt32(2).ToValue(),
			want: value.SmallInt(7).ToValue(),
		},
		"shift by UInt32 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.UInt32(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by UInt32 80 >> 60": {
			a:    value.SmallInt(80),
			b:    value.UInt32(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by UInt16 5 >> 1": {
			a:    value.SmallInt(5),
			b:    value.UInt16(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by UInt16 28 >> 2": {
			a:    value.SmallInt(28),
			b:    value.UInt16(2).ToValue(),
			want: value.SmallInt(7).ToValue(),
		},
		"shift by UInt16 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.UInt16(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by UInt16 80 >> 60": {
			a:    value.SmallInt(80),
			b:    value.UInt16(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},

		"shift by UInt8 5 >> 1": {
			a:    value.SmallInt(5),
			b:    value.UInt8(1).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"shift by UInt8 28 >> 2": {
			a:    value.SmallInt(28),
			b:    value.UInt8(2).ToValue(),
			want: value.SmallInt(7).ToValue(),
		},
		"shift by UInt8 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.UInt8(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by UInt8 80 >> 60": {
			a:    value.SmallInt(80),
			b:    value.UInt8(60).ToValue(),
			want: value.SmallInt(0).ToValue(),
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

func TestSmallInt_LeftBitshift(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"shift by String and return an error": {
			a:   value.SmallInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be used as a bitshift operand")),
		},
		"shift by Float and return an error": {
			a:   value.SmallInt(5),
			b:   value.Float(2.5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be used as a bitshift operand")),
		},

		"shift by SmallInt 5 << 1": {
			a:    value.SmallInt(5),
			b:    value.SmallInt(1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by SmallInt 12 << -1": {
			a:    value.SmallInt(12),
			b:    value.SmallInt(-1).ToValue(),
			want: value.SmallInt(6).ToValue(),
		},
		"shift by SmallInt 418 << -5": {
			a:    value.SmallInt(418),
			b:    value.SmallInt(-5).ToValue(),
			want: value.SmallInt(13).ToValue(),
		},
		"shift by SmallInt 75 << 0": {
			a:    value.SmallInt(75),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by BigInt 5 << 1": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.NewBigInt(1)),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by BigInt 12 << -1": {
			a:    value.SmallInt(12),
			b:    value.Ref(value.NewBigInt(-1)),
			want: value.SmallInt(6).ToValue(),
		},
		"shift by BigInt 418 << -5": {
			a:    value.SmallInt(418),
			b:    value.Ref(value.NewBigInt(-5)),
			want: value.SmallInt(13).ToValue(),
		},
		"shift by BigInt 75 >> 0": {
			a:    value.SmallInt(75),
			b:    value.Ref(value.NewBigInt(0)),
			want: value.SmallInt(75).ToValue(),
		},

		"shift by Int64 5 << 1": {
			a:    value.SmallInt(5),
			b:    value.Int64(1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by Int64 12 << -1": {
			a:    value.SmallInt(12),
			b:    value.Int64(-1).ToValue(),
			want: value.SmallInt(6).ToValue(),
		},
		"shift by Int64 418 << -5": {
			a:    value.SmallInt(418),
			b:    value.Int64(-5).ToValue(),
			want: value.SmallInt(13).ToValue(),
		},
		"shift by Int64 75 << 0": {
			a:    value.SmallInt(75),
			b:    value.Int64(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by Int32 5 << 1": {
			a:    value.SmallInt(5),
			b:    value.Int32(1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by Int32 12 << -1": {
			a:    value.SmallInt(12),
			b:    value.Int32(-1).ToValue(),
			want: value.SmallInt(6).ToValue(),
		},
		"shift by Int32 418 << -5": {
			a:    value.SmallInt(418),
			b:    value.Int32(-5).ToValue(),
			want: value.SmallInt(13).ToValue(),
		},
		"shift by Int32 75 << 0": {
			a:    value.SmallInt(75),
			b:    value.Int32(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},

		"shift by Int16 5 << 1": {
			a:    value.SmallInt(5),
			b:    value.Int16(1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by Int16 12 << -1": {
			a:    value.SmallInt(12),
			b:    value.Int16(-1).ToValue(),
			want: value.SmallInt(6).ToValue(),
		},
		"shift by Int16 418 << -5": {
			a:    value.SmallInt(418),
			b:    value.Int16(-5).ToValue(),
			want: value.SmallInt(13).ToValue(),
		},
		"shift by Int16 75 << 0": {
			a:    value.SmallInt(75),
			b:    value.Int16(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by Int8 5 << 1": {
			a:    value.SmallInt(5),
			b:    value.Int8(1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by Int8 12 << -1": {
			a:    value.SmallInt(12),
			b:    value.Int8(-1).ToValue(),
			want: value.SmallInt(6).ToValue(),
		},
		"shift by Int8 418 << -5": {
			a:    value.SmallInt(418),
			b:    value.Int8(-5).ToValue(),
			want: value.SmallInt(13).ToValue(),
		},
		"shift by Int8 75 << 0": {
			a:    value.SmallInt(75),
			b:    value.Int8(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},

		"shift by UInt64 5 << 1": {
			a:    value.SmallInt(5),
			b:    value.UInt64(1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by UInt64 75 << 0": {
			a:    value.SmallInt(75),
			b:    value.UInt64(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},

		"shift by UInt32 5 << 1": {
			a:    value.SmallInt(5),
			b:    value.UInt32(1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by UInt32 75 << 0": {
			a:    value.SmallInt(75),
			b:    value.UInt32(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
		"shift by UInt16 5 << 1": {
			a:    value.SmallInt(5),
			b:    value.UInt16(1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by UInt16 75 << 0": {
			a:    value.SmallInt(75),
			b:    value.UInt16(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},

		"shift by UInt8 5 << 1": {
			a:    value.SmallInt(5),
			b:    value.UInt8(1).ToValue(),
			want: value.SmallInt(10).ToValue(),
		},
		"shift by UInt8 75 << 0": {
			a:    value.SmallInt(75),
			b:    value.UInt8(0).ToValue(),
			want: value.SmallInt(75).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LeftBitshiftVal(tc.b)
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

func TestSmallInt_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt & String and return an error": {
			a:   value.SmallInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"SmallInt & Int32 and return an error": {
			a:   value.SmallInt(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int`")),
		},
		"SmallInt & Float and return an error": {
			a:   value.SmallInt(5),
			b:   value.Float(2.5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Int`")),
		},

		"23 & 10": {
			a:    value.SmallInt(23),
			b:    value.SmallInt(10).ToValue(),
			want: value.SmallInt(2).ToValue(),
		},
		"11 & 7": {
			a:    value.SmallInt(11),
			b:    value.SmallInt(7).ToValue(),
			want: value.SmallInt(3).ToValue(),
		},
		"-14 & 23": {
			a:    value.SmallInt(-14),
			b:    value.SmallInt(23).ToValue(),
			want: value.SmallInt(18).ToValue(),
		},
		"258 & 0": {
			a:    value.SmallInt(258),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"124 & 255": {
			a:    value.SmallInt(124),
			b:    value.SmallInt(255).ToValue(),
			want: value.SmallInt(124).ToValue(),
		},

		"255 & 9223372036857247042": {
			a:    value.SmallInt(255),
			b:    value.Ref(value.ParseBigIntPanic("9223372036857247042", 10)),
			want: value.SmallInt(66).ToValue(),
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

func TestSmallInt_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt | String and return an error": {
			a:   value.SmallInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"SmallInt | Int32 and return an error": {
			a:   value.SmallInt(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int`")),
		},
		"SmallInt | Float and return an error": {
			a:   value.SmallInt(5),
			b:   value.Float(2.5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Int`")),
		},

		"23 | 10": {
			a:    value.SmallInt(23),
			b:    value.SmallInt(10).ToValue(),
			want: value.SmallInt(31).ToValue(),
		},
		"11 | 7": {
			a:    value.SmallInt(11),
			b:    value.SmallInt(7).ToValue(),
			want: value.SmallInt(15).ToValue(),
		},
		"-14 | 23": {
			a:    value.SmallInt(-14),
			b:    value.SmallInt(23).ToValue(),
			want: value.SmallInt(-9).ToValue(),
		},
		"258 | 0": {
			a:    value.SmallInt(258),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(258).ToValue(),
		},
		"124 | 255": {
			a:    value.SmallInt(124),
			b:    value.SmallInt(255).ToValue(),
			want: value.SmallInt(255).ToValue(),
		},

		"255 | 9223372036857247042": {
			a:    value.SmallInt(255),
			b:    value.Ref(value.ParseBigIntPanic("9223372036857247042", 10)),
			want: value.Ref(value.ParseBigIntPanic("9223372036857247231", 10)),
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

func TestSmallInt_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt ^ String and return an error": {
			a:   value.SmallInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},
		"SmallInt ^ Int32 and return an error": {
			a:   value.SmallInt(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Int`")),
		},
		"SmallInt ^ Float and return an error": {
			a:   value.SmallInt(5),
			b:   value.Float(2.5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Int`")),
		},

		"23 ^ 10": {
			a:    value.SmallInt(23),
			b:    value.SmallInt(10).ToValue(),
			want: value.SmallInt(29).ToValue(),
		},
		"11 ^ 7": {
			a:    value.SmallInt(11),
			b:    value.SmallInt(7).ToValue(),
			want: value.SmallInt(12).ToValue(),
		},
		"-14 ^ 23": {
			a:    value.SmallInt(-14),
			b:    value.SmallInt(23).ToValue(),
			want: value.SmallInt(-27).ToValue(),
		},
		"258 ^ 0": {
			a:    value.SmallInt(258),
			b:    value.SmallInt(0).ToValue(),
			want: value.SmallInt(258).ToValue(),
		},
		"124 ^ 255": {
			a:    value.SmallInt(124),
			b:    value.SmallInt(255).ToValue(),
			want: value.SmallInt(131).ToValue(),
		},

		"255 ^ 9223372036857247042": {
			a:    value.SmallInt(255),
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

func TestSmallInt_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.SmallInt(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 % 3": {
			a:    value.SmallInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"SmallInt 76 % 6": {
			a:    value.SmallInt(76),
			b:    value.SmallInt(6).ToValue(),
			want: value.SmallInt(4).ToValue(),
		},
		"SmallInt -76 % 6": {
			a:    value.SmallInt(-76),
			b:    value.SmallInt(6).ToValue(),
			want: value.SmallInt(-4).ToValue(),
		},
		"SmallInt 76 % -6": {
			a:    value.SmallInt(76),
			b:    value.SmallInt(-6).ToValue(),
			want: value.SmallInt(4).ToValue(),
		},
		"SmallInt -76 % -6": {
			a:    value.SmallInt(-76),
			b:    value.SmallInt(-6).ToValue(),
			want: value.SmallInt(-4).ToValue(),
		},
		"SmallInt 124 % 9": {
			a:    value.SmallInt(124),
			b:    value.SmallInt(9).ToValue(),
			want: value.SmallInt(7).ToValue(),
		},
		"BigInt 25 % 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigInt 76 % 6": {
			a:    value.SmallInt(76),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.SmallInt(4).ToValue(),
		},
		"BigInt -76 % 6": {
			a:    value.SmallInt(-76),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.SmallInt(-4).ToValue(),
		},
		"BigInt 76 % -6": {
			a:    value.SmallInt(76),
			b:    value.Ref(value.NewBigInt(-6)),
			want: value.SmallInt(4).ToValue(),
		},
		"BigInt -76 % -6": {
			a:    value.SmallInt(-76),
			b:    value.Ref(value.NewBigInt(-6)),
			want: value.SmallInt(-4).ToValue(),
		},
		"BigInt 124 % 9": {
			a:    value.SmallInt(124),
			b:    value.Ref(value.NewBigInt(9)),
			want: value.SmallInt(7).ToValue(),
		},
		"Float 25 % 3": {
			a:    value.SmallInt(25),
			b:    value.Float(3).ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float 76 % 6": {
			a:    value.SmallInt(76),
			b:    value.Float(6).ToValue(),
			want: value.Float(4).ToValue(),
		},
		"Float 124 % 9": {
			a:    value.SmallInt(124),
			b:    value.Float(9).ToValue(),
			want: value.Float(7).ToValue(),
		},
		"Float 124 % +Inf": {
			a:    value.SmallInt(124),
			b:    value.FloatInf().ToValue(),
			want: value.Float(124).ToValue(),
		},
		"Float 124 % -Inf": {
			a:    value.SmallInt(124),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(124).ToValue(),
		},
		"Float 74 % 6.25": {
			a:    value.SmallInt(74),
			b:    value.Float(6.25).ToValue(),
			want: value.Float(5.25).ToValue(),
		},
		"Float -74 % 6.25": {
			a:    value.SmallInt(-74),
			b:    value.Float(6.25).ToValue(),
			want: value.Float(-5.25).ToValue(),
		},
		"Float 74 % -6.25": {
			a:    value.SmallInt(74),
			b:    value.Float(-6.25).ToValue(),
			want: value.Float(5.25).ToValue(),
		},
		"Float -74 % -6.25": {
			a:    value.SmallInt(-74),
			b:    value.Float(-6.25).ToValue(),
			want: value.Float(-5.25).ToValue(),
		},
		"Float 25 % 0": { // Mod(x, 0) = NaN
			a:    value.SmallInt(25),
			b:    value.Float(0).ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    value.SmallInt(25),
			b:    value.FloatInf().ToValue(),
			want: value.Float(25).ToValue(),
		},
		"Float -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    value.SmallInt(-87),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(-87).ToValue(),
		},
		"Float 49 % NaN": { // Mod(x, NaN) = NaN
			a:    value.SmallInt(49),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"BigFloat 25 % 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.Ref(value.NewBigFloat(1).SetPrecision(64)),
		},
		"BigFloat 76 % 6": {
			a:    value.SmallInt(76),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.Ref(value.NewBigFloat(4).SetPrecision(64)),
		},
		"BigFloat 124 % 9": {
			a:    value.SmallInt(124),
			b:    value.Ref(value.NewBigFloat(9)),
			want: value.Ref(value.NewBigFloat(7).SetPrecision(64)),
		},
		"BigFloat 74 % 6.25": {
			a:    value.SmallInt(74),
			b:    value.Ref(value.NewBigFloat(6.25)),
			want: value.Ref(value.NewBigFloat(5.25).SetPrecision(64)),
		},
		"BigFloat 74 % 6.25p92": {
			a:    value.SmallInt(74),
			b:    value.Ref(value.NewBigFloat(6.25).SetPrecision(92)),
			want: value.Ref(value.NewBigFloat(5.25).SetPrecision(92)),
		},
		"BigFloat -74 % 6.25": {
			a:    value.SmallInt(-74),
			b:    value.Ref(value.NewBigFloat(6.25)),
			want: value.Ref(value.NewBigFloat(-5.25).SetPrecision(64)),
		},
		"BigFloat 74 % -6.25": {
			a:    value.SmallInt(74),
			b:    value.Ref(value.NewBigFloat(-6.25)),
			want: value.Ref(value.NewBigFloat(5.25).SetPrecision(64)),
		},
		"BigFloat -74 % -6.25": {
			a:    value.SmallInt(-74),
			b:    value.Ref(value.NewBigFloat(-6.25)),
			want: value.Ref(value.NewBigFloat(-5.25).SetPrecision(64)),
		},
		"BigFloat 25 % 0": { // Mod(x, 0) = NaN
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigFloat(0)),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    value.SmallInt(25),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(25).SetPrecision(64)),
		},
		"BigFloat -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    value.SmallInt(-87),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(-87).SetPrecision(64)),
		},
		"BigFloat 49 % NaN": { // Mod(x, NaN) = NaN
			a:    value.SmallInt(49),
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

func TestSmallInt_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 > 3": {
			a:    value.SmallInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.True,
		},
		"SmallInt 6 > 18": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6 > 6": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.False,
		},

		"BigInt 25 > 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.True,
		},
		"BigInt 6 > 18": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6 > 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.False,
		},
		"Float 25 > 3": {
			a:    value.SmallInt(25),
			b:    value.Float(3).ToValue(),
			want: value.True,
		},
		"Float 6 > 18.5": {
			a:    value.SmallInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6 > 6": {
			a:    value.SmallInt(6),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6 > Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6 > -Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.True,
		},
		"Float 6 > NaN": {
			a:    value.SmallInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25 > 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.True,
		},
		"BigFloat 6 > 18.5": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6 > 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat 6 > Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6 > -Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.True,
		},
		"BigFloat 6 > NaN": {
			a:    value.SmallInt(6),
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

func TestSmallInt_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 >= 3": {
			a:    value.SmallInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.True,
		},
		"SmallInt 6 >= 18": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6 >= 6": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25 >= 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.True,
		},
		"BigInt 6 >= 18": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6 >= 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},
		"Float 25 >= 3": {
			a:    value.SmallInt(25),
			b:    value.Float(3).ToValue(),
			want: value.True,
		},
		"Float 6 >= 18.5": {
			a:    value.SmallInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6 >= 6": {
			a:    value.SmallInt(6),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6 >= Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6 >= -Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.True,
		},
		"Float 6 >= NaN": {
			a:    value.SmallInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25 >= 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.True,
		},
		"BigFloat 6 >= 18.5": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6 >= 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat 6 >= Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6 >= -Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.True,
		},
		"BigFloat 6 >= NaN": {
			a:    value.SmallInt(6),
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

func TestSmallInt_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 < 3": {
			a:    value.SmallInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6 < 18": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.True,
		},
		"SmallInt 6 < 6": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.False,
		},

		"BigInt 25 < 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6 < 18": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.True,
		},
		"BigInt 6 < 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.False,
		},
		"Float 25 < 3": {
			a:    value.SmallInt(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6 < 18.5": {
			a:    value.SmallInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.True,
		},
		"Float 6 < 6": {
			a:    value.SmallInt(6),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6 < Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float 6 < -Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6 < NaN": {
			a:    value.SmallInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25 < 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6 < 18.5": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.True,
		},
		"BigFloat 6 < 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat 6 < Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat 6 < -Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6 < NaN": {
			a:    value.SmallInt(6),
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

func TestSmallInt_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.String("foo")),
			want: value.False,
			err:  value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Int`")),
		},

		"SmallInt 25 <= 3": {
			a:    value.SmallInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6 <= 18": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.True,
		},
		"SmallInt 6 <= 6": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25 <= 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6 <= 18": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.True,
		},
		"BigInt 6 <= 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},
		"Float 25 <= 3": {
			a:    value.SmallInt(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6 <= 18.5": {
			a:    value.SmallInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.True,
		},
		"Float 6 <= 6": {
			a:    value.SmallInt(6),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6 <= Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float 6 <= -Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6 <= NaN": {
			a:    value.SmallInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25 <= 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6 <= 18.5": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.True,
		},
		"BigFloat 6 <= 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat 6 <= Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat 6 <= -Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6 <= NaN": {
			a:    value.SmallInt(6),
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

func TestSmallInt_LaxEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
	}{
		"String 5 =~ '5'": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},
		"Char 5 =~ `5`": {
			a:    value.SmallInt(5),
			b:    value.Char('5').ToValue(),
			want: value.False,
		},

		"Int64 5 =~ 5i64": {
			a:    value.SmallInt(5),
			b:    value.Int64(5).ToValue(),
			want: value.True,
		},
		"Int64 4 =~ 5i64": {
			a:    value.SmallInt(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},

		"Int32 5 =~ 5i32": {
			a:    value.SmallInt(5),
			b:    value.Int32(5).ToValue(),
			want: value.True,
		},
		"Int32 4 =~ 5i32": {
			a:    value.SmallInt(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"Int16 5 =~ 5i16": {
			a:    value.SmallInt(5),
			b:    value.Int16(5).ToValue(),
			want: value.True,
		},
		"Int16 4 =~ 5i16": {
			a:    value.SmallInt(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"Int8 5 =~ 5i8": {
			a:    value.SmallInt(5),
			b:    value.Int8(5).ToValue(),
			want: value.True,
		},
		"Int8 4 =~ 5i8": {
			a:    value.SmallInt(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt64 5 =~ 5u64": {
			a:    value.SmallInt(5),
			b:    value.UInt64(5).ToValue(),
			want: value.True,
		},
		"UInt64 4 =~ 5u64": {
			a:    value.SmallInt(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"UInt32 5 =~ 5u32": {
			a:    value.SmallInt(5),
			b:    value.UInt32(5).ToValue(),
			want: value.True,
		},
		"UInt32 4 =~ 5u32": {
			a:    value.SmallInt(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt16 5 =~ 5u16": {
			a:    value.SmallInt(5),
			b:    value.UInt16(5).ToValue(),
			want: value.True,
		},
		"UInt16 4 =~ 5u16": {
			a:    value.SmallInt(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"UInt8 5 =~ 5u8": {
			a:    value.SmallInt(5),
			b:    value.UInt8(5).ToValue(),
			want: value.True,
		},
		"UInt8 4 =~ 5u8": {
			a:    value.SmallInt(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5 =~ 5f64": {
			a:    value.SmallInt(5),
			b:    value.Float64(5).ToValue(),
			want: value.True,
		},
		"Float64 5 =~ 5.5f64": {
			a:    value.SmallInt(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 4 =~ 5f64": {
			a:    value.SmallInt(4),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float32 5 =~ 5f32": {
			a:    value.SmallInt(5),
			b:    value.Float32(5).ToValue(),
			want: value.True,
		},
		"Float32 5 =~ 5.5f32": {
			a:    value.SmallInt(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 4 =~ 5f32": {
			a:    value.SmallInt(4),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},

		"SmallInt 25 =~ 3": {
			a:    value.SmallInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6 =~ 18": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6 =~ 6": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25 =~ 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6 =~ 18": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6 =~ 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},
		"Float 25 =~ 3": {
			a:    value.SmallInt(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6 =~ 18.5": {
			a:    value.SmallInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6 =~ 6": {
			a:    value.SmallInt(6),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6 =~ Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6 =~ -Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6 =~ NaN": {
			a:    value.SmallInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25 =~ 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6 =~ 18.5": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6 =~ 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat 6 =~ Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6 =~ -Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6 =~ NaN": {
			a:    value.SmallInt(6),
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

func TestSmallInt_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
	}{
		"String 5 == '5'": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},
		"Char 5 == `5`": {
			a:    value.SmallInt(5),
			b:    value.Char('5').ToValue(),
			want: value.False,
		},

		"Int64 5 == 5i64": {
			a:    value.SmallInt(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4 == 5i64": {
			a:    value.SmallInt(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},

		"Int32 5 == 5i32": {
			a:    value.SmallInt(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4 == 5i32": {
			a:    value.SmallInt(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"Int16 5 == 5i16": {
			a:    value.SmallInt(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4 == 5i16": {
			a:    value.SmallInt(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int8 5 == 5i8": {
			a:    value.SmallInt(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4 == 5i8": {
			a:    value.SmallInt(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt64 5 == 5u64": {
			a:    value.SmallInt(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4 == 5u64": {
			a:    value.SmallInt(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"UInt32 5 == 5u32": {
			a:    value.SmallInt(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 4 == 5u32": {
			a:    value.SmallInt(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"UInt16 5 == 5u16": {
			a:    value.SmallInt(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 4 == 5u16": {
			a:    value.SmallInt(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt8 5 == 5u8": {
			a:    value.SmallInt(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 4 == 5u8": {
			a:    value.SmallInt(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5 == 5f64": {
			a:    value.SmallInt(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5 == 5.5f64": {
			a:    value.SmallInt(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 4 == 5f64": {
			a:    value.SmallInt(4),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},

		"Float32 5 == 5f32": {
			a:    value.SmallInt(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5 == 5.5f32": {
			a:    value.SmallInt(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 4 == 5f32": {
			a:    value.SmallInt(4),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"SmallInt 25 == 3": {
			a:    value.SmallInt(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6 == 18": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6 == 6": {
			a:    value.SmallInt(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25 == 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6 == 18": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6 == 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},

		"Float 25 == 3": {
			a:    value.SmallInt(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6 == 18.5": {
			a:    value.SmallInt(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6 == 6": {
			a:    value.SmallInt(6),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6 == Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6 == -Inf": {
			a:    value.SmallInt(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6 == NaN": {
			a:    value.SmallInt(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"BigFloat 25 == 3": {
			a:    value.SmallInt(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6 == 18.5": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6 == 6": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat 6 == Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6 == -Inf": {
			a:    value.SmallInt(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6 == NaN": {
			a:    value.SmallInt(6),
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
