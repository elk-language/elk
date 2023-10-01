package value

import (
	"math"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSmallInt_Add(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"add String and return an error": {
			a:   SmallInt(3),
			b:   String("foo"),
			err: NewCoerceError(SmallInt(3), String("foo")),
		},
		"add SmallInt and return SmallInt": {
			a:    SmallInt(3),
			b:    SmallInt(10),
			want: SmallInt(13),
		},
		"add SmallInt overflow and return BigInt": {
			a:    SmallInt(math.MaxInt64),
			b:    SmallInt(10),
			want: ParseBigIntPanic("9223372036854775817", 10),
		},
		"add BigInt and return BigInt": {
			a:    SmallInt(20),
			b:    ParseBigIntPanic("9223372036854775817", 10),
			want: ParseBigIntPanic("9223372036854775837", 10),
		},
		"add BigInt and return SmallInt": {
			a:    SmallInt(-20),
			b:    ParseBigIntPanic("9223372036854775817", 10),
			want: SmallInt(9223372036854775797),
		},
		"add Float and return Float": {
			a:    SmallInt(-20),
			b:    Float(2.5),
			want: Float(-17.5),
		},
		"add Float NaN and return Float NaN": {
			a:    SmallInt(-20),
			b:    FloatNaN(),
			want: FloatNaN(),
		},
		"add Float +Inf and return Float +Inf": {
			a:    SmallInt(-20),
			b:    FloatInf(),
			want: FloatInf(),
		},
		"add Float -Inf and return Float -Inf": {
			a:    SmallInt(-20),
			b:    FloatNegInf(),
			want: FloatNegInf(),
		},
		"add BigFloat NaN and return BigFloat NaN": {
			a:    SmallInt(56),
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
		"add BigFloat +Inf and return BigFloat +Inf": {
			a:    SmallInt(56),
			b:    BigFloatInf(),
			want: BigFloatInf(),
		},
		"add BigFloat -Inf and return BigFloat -Inf": {
			a:    SmallInt(56),
			b:    BigFloatNegInf(),
			want: BigFloatNegInf(),
		},
		"add BigFloat and return BigFloat with 64bit precision": {
			a:    SmallInt(56),
			b:    NewBigFloat(2.5),
			want: NewBigFloat(58.5).SetPrecision(64),
		},
		"add BigFloat and return BigFloat with 80bit precision": {
			a:    SmallInt(56),
			b:    NewBigFloat(2.5).SetPrecision(80),
			want: ToElkBigFloat((&big.Float{}).SetPrec(80).Add(big.NewFloat(56), big.NewFloat(2.5))),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Add(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				bigFloatComparer,
				floatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"subtract String and return an error": {
			a:   SmallInt(3),
			b:   String("foo"),
			err: NewCoerceError(SmallInt(3), String("foo")),
		},
		"subtract SmallInt and return SmallInt": {
			a:    SmallInt(3),
			b:    SmallInt(10),
			want: SmallInt(-7),
		},
		"subtract SmallInt underflow and return BigInt": {
			a:    SmallInt(math.MinInt64),
			b:    SmallInt(10),
			want: ParseBigIntPanic("-9223372036854775818", 10),
		},
		"subtract BigInt and return BigInt": {
			a:    SmallInt(5),
			b:    ParseBigIntPanic("9223372036854775817", 10),
			want: ParseBigIntPanic("-9223372036854775812", 10),
		},
		"subtract BigInt and return SmallInt": {
			a:    SmallInt(20),
			b:    ParseBigIntPanic("9223372036854775817", 10),
			want: SmallInt(-9223372036854775797),
		},

		"subtract Float and return Float": {
			a:    SmallInt(-20),
			b:    Float(2.5),
			want: Float(-22.5),
		},
		"subtract Float NaN and return Float NaN": {
			a:    26,
			b:    FloatNaN(),
			want: FloatNaN(),
		},
		"subtract Float +Inf and return Float -Inf": {
			a:    19,
			b:    FloatInf(),
			want: FloatNegInf(),
		},
		"subtract Float -Inf and return Float +Inf": {
			a:    3,
			b:    FloatNegInf(),
			want: FloatInf(),
		},

		"subtract BigFloat and return BigFloat with 64bit precision": {
			a:    SmallInt(56),
			b:    NewBigFloat(2.5),
			want: NewBigFloat(53.5).SetPrecision(64),
		},
		"subtract BigFloat and return BigFloat with 80bit precision": {
			a:    SmallInt(56),
			b:    NewBigFloat(2.5).SetPrecision(80),
			want: ToElkBigFloat((&big.Float{}).SetPrec(80).Sub(big.NewFloat(56), big.NewFloat(2.5))),
		},

		"subtract BigFloat NaN and return BigFloat NaN": {
			a:    35,
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
		"subtract BigFloat +Inf and return BigFloat -Inf": {
			a:    56,
			b:    BigFloatInf(),
			want: BigFloatNegInf(),
		},
		"subtract BigFloat -Inf and return BigFloat +Inf": {
			a:    -12,
			b:    BigFloatNegInf(),
			want: BigFloatInf(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Subtract(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"multiply by String and return an error": {
			a:   SmallInt(3),
			b:   String("foo"),
			err: NewCoerceError(SmallInt(3), String("foo")),
		},
		"multiply by SmallInt and return SmallInt": {
			a:    SmallInt(3),
			b:    SmallInt(10),
			want: SmallInt(30),
		},
		"multiply by SmallInt overflow and return BigInt": {
			a:    SmallInt(math.MaxInt64),
			b:    SmallInt(10),
			want: ParseBigIntPanic("92233720368547758070", 10),
		},
		"multiply by BigInt and return BigInt": {
			a:    SmallInt(20),
			b:    ParseBigIntPanic("9223372036854775817", 10),
			want: ParseBigIntPanic("184467440737095516340", 10),
		},
		"multiply BigInt and return SmallInt": {
			a:    SmallInt(-1),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: SmallInt(math.MinInt64),
		},
		"multiply by Float and return Float": {
			a:    SmallInt(-20),
			b:    Float(2.5),
			want: Float(-50),
		},
		"multiply by BigFloat and return BigFloat with 64bit precision": {
			a:    SmallInt(56),
			b:    NewBigFloat(2.5),
			want: NewBigFloat(140).SetPrecision(64),
		},
		"multiply by BigFloat and return BigFloat with 80bit precision": {
			a:    SmallInt(56),
			b:    NewBigFloat(2.5).SetPrecision(80),
			want: ToElkBigFloat((&big.Float{}).SetPrec(80).Mul(big.NewFloat(56), big.NewFloat(2.5))),
		},

		"multiply by Float NaN and return Float NaN": {
			a:    234,
			b:    FloatNaN(),
			want: FloatNaN(),
		},
		"multiply by Float +Inf and return Float +Inf": {
			a:    234,
			b:    FloatInf(),
			want: FloatInf(),
		},
		"multiply by Float +Inf and return Float -Inf": {
			a:    -123,
			b:    FloatInf(),
			want: FloatNegInf(),
		},
		"multiply by Float -Inf and return Float -Inf": {
			a:    56,
			b:    FloatNegInf(),
			want: FloatNegInf(),
		},
		"multiply by Float -Inf and return Float +Inf": {
			a:    -5,
			b:    FloatNegInf(),
			want: FloatInf(),
		},

		"multiply by BigFloat NaN and return BigFloat NaN": {
			a:    75,
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
		"multiply by BigFloat +Inf and return BigFloat +Inf": {
			a:    15,
			b:    BigFloatInf(),
			want: BigFloatInf(),
		},
		"multiply by BigFloat +Inf and return BigFloat -Inf": {
			a:    -2,
			b:    BigFloatInf(),
			want: BigFloatNegInf(),
		},
		"multiply by BigFloat -Inf and return BigFloat -Inf": {
			a:    7,
			b:    BigFloatNegInf(),
			want: BigFloatNegInf(),
		},
		"multiply by BigFloat -Inf and return BigFloat +Inf": {
			a:    -9,
			b:    BigFloatNegInf(),
			want: BigFloatInf(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Multiply(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_Divide(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"divide by String and return an error": {
			a:   SmallInt(3),
			b:   String("foo"),
			err: NewCoerceError(SmallInt(3), String("foo")),
		},
		"divide by SmallInt and return SmallInt": {
			a:    SmallInt(30),
			b:    SmallInt(10),
			want: SmallInt(3),
		},
		"divide by SmallInt overflow and return BigInt": {
			a:    SmallInt(math.MinInt64),
			b:    SmallInt(-1),
			want: ParseBigIntPanic("9223372036854775808", 10),
		},
		"divide by BigInt and return SmallInt": {
			a:    SmallInt(20),
			b:    ParseBigIntPanic("9223372036854775817", 10),
			want: SmallInt(0),
		},

		"Float -20 / 0.5": {
			a:    SmallInt(-20),
			b:    Float(0.5),
			want: Float(-40),
		},
		"Float 234 / 0": {
			a:    234,
			b:    Float(0),
			want: FloatInf(),
		},
		"Float -234 / 0": {
			a:    -234,
			b:    Float(0),
			want: FloatNegInf(),
		},
		"Float 234 / NaN": {
			a:    234,
			b:    FloatNaN(),
			want: FloatNaN(),
		},
		"Float 234 / +Inf": {
			a:    234,
			b:    FloatInf(),
			want: Float(0),
		},
		"Float 56 / -Inf": {
			a:    56,
			b:    FloatNegInf(),
			want: Float(0),
		},

		"BigFloat 55 / 2.5 with 64bit precision": {
			a:    SmallInt(55),
			b:    NewBigFloat(2.5),
			want: NewBigFloat(22.0).SetPrecision(64),
		},
		"BigFloat 55 / 2.5 with 80bit precision": {
			a:    SmallInt(55),
			b:    NewBigFloat(2.5).SetPrecision(80),
			want: NewBigFloat(22).SetPrecision(80),
		},
		"BigFloat 234 / 0": {
			a:    234,
			b:    NewBigFloat(0),
			want: BigFloatInf(),
		},
		"BigFloat -234 / 0": {
			a:    -234,
			b:    NewBigFloat(0),
			want: BigFloatNegInf(),
		},
		"BigFloat 234 / NaN": {
			a:    234,
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
		"BigFloat 234 / +Inf": {
			a:    234,
			b:    BigFloatInf(),
			want: NewBigFloat(0).SetPrecision(64),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Divide(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_addOverflow(t *testing.T) {
	tests := map[string]struct {
		a, b SmallInt
		want SmallInt
		ok   bool
	}{
		"not overflow": {
			a:    SmallInt(3),
			b:    SmallInt(10),
			want: SmallInt(13),
			ok:   true,
		},
		"overflow": {
			a:    SmallInt(10),
			b:    SmallInt(math.MaxInt64),
			want: SmallInt(math.MinInt64 + 9),
			ok:   false,
		},
		"not underflow": {
			a:    SmallInt(math.MinInt64 + 20),
			b:    SmallInt(-18),
			want: SmallInt(math.MinInt64 + 2),
			ok:   true,
		},
		"not underflow positive": {
			a:    SmallInt(math.MinInt64 + 20),
			b:    SmallInt(18),
			want: SmallInt(math.MinInt64 + 38),
			ok:   true,
		},
		"underflow": {
			a:    SmallInt(math.MinInt64),
			b:    SmallInt(-20),
			want: SmallInt(math.MaxInt64 - 19),
			ok:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.a.addOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_subtractOverflow(t *testing.T) {
	tests := map[string]struct {
		a, b SmallInt
		want SmallInt
		ok   bool
	}{
		"not underflow": {
			a:    SmallInt(3),
			b:    SmallInt(10),
			want: SmallInt(-7),
			ok:   true,
		},
		"underflow": {
			a:    SmallInt(math.MinInt64),
			b:    SmallInt(3),
			want: SmallInt(math.MaxInt64 - 2),
			ok:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.a.subtractOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_multiplyOverflow(t *testing.T) {
	tests := map[string]struct {
		a, b SmallInt
		want SmallInt
		ok   bool
	}{
		"not overflow": {
			a:    SmallInt(3),
			b:    SmallInt(10),
			want: SmallInt(30),
			ok:   true,
		},
		"overflow": {
			a:    SmallInt(10),
			b:    SmallInt(math.MaxInt64),
			want: SmallInt(-10),
			ok:   false,
		},
		"not underflow": {
			a:    SmallInt(-125),
			b:    SmallInt(5),
			want: SmallInt(-625),
			ok:   true,
		},
		"underflow": {
			a:    SmallInt(math.MinInt64),
			b:    SmallInt(2),
			want: SmallInt(0),
			ok:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.a.multiplyOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_divideOverflow(t *testing.T) {
	tests := map[string]struct {
		a, b SmallInt
		want SmallInt
		ok   bool
	}{
		"not overflow": {
			a:    SmallInt(20),
			b:    SmallInt(5),
			want: SmallInt(4),
			ok:   true,
		},
		"overflow": {
			a:    SmallInt(math.MinInt64),
			b:    SmallInt(-1),
			want: SmallInt(math.MinInt64),
			ok:   false,
		},
		"not underflow": {
			a:    SmallInt(-625),
			b:    SmallInt(5),
			want: SmallInt(-125),
			ok:   true,
		},
		"division by zero": {
			a:    SmallInt(math.MinInt64),
			b:    SmallInt(0),
			want: SmallInt(0),
			ok:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.a.divideOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"exponentiate String and return an error": {
			a:   SmallInt(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::SmallInt`"),
		},
		"exponentiate Int32 and return an error": {
			a:   SmallInt(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::SmallInt`"),
		},

		"SmallInt 5 ** 2": {
			a:    SmallInt(5),
			b:    SmallInt(2),
			want: SmallInt(25),
		},
		"SmallInt 7 ** 8": {
			a:    SmallInt(7),
			b:    SmallInt(8),
			want: SmallInt(5764801),
		},
		"SmallInt 2 ** 5": {
			a:    SmallInt(2),
			b:    SmallInt(5),
			want: SmallInt(32),
		},
		"SmallInt 6 ** 1": {
			a:    SmallInt(6),
			b:    SmallInt(1),
			want: SmallInt(6),
		},
		"SmallInt 2 ** 64": {
			a:    SmallInt(2),
			b:    SmallInt(64),
			want: ParseBigIntPanic("18446744073709551616", 10),
		},
		"SmallInt 4 ** -2": {
			a:    SmallInt(4),
			b:    SmallInt(-2),
			want: SmallInt(1),
		},
		"SmallInt 25 ** 0": {
			a:    SmallInt(25),
			b:    SmallInt(0),
			want: SmallInt(1),
		},

		"BigInt 5 ** 2": {
			a:    SmallInt(5),
			b:    NewBigInt(2),
			want: SmallInt(25),
		},
		"BigInt 7 ** 8": {
			a:    SmallInt(7),
			b:    NewBigInt(8),
			want: SmallInt(5764801),
		},
		"BigInt 2 ** 5": {
			a:    SmallInt(2),
			b:    NewBigInt(5),
			want: SmallInt(32),
		},
		"BigInt 6 ** 1": {
			a:    SmallInt(6),
			b:    NewBigInt(1),
			want: SmallInt(6),
		},
		"BigInt 2 ** 64": {
			a:    SmallInt(2),
			b:    NewBigInt(64),
			want: ParseBigIntPanic("18446744073709551616", 10),
		},
		"BigInt 4 ** -2": {
			a:    SmallInt(4),
			b:    NewBigInt(-2),
			want: SmallInt(1),
		},
		"BigInt 25 ** 0": {
			a:    SmallInt(25),
			b:    NewBigInt(0),
			want: SmallInt(1),
		},

		"Float 5 ** 2": {
			a:    SmallInt(5),
			b:    Float(2),
			want: Float(25),
		},
		"Float 7 ** 8": {
			a:    SmallInt(7),
			b:    Float(8),
			want: Float(5764801),
		},
		"Float 3 ** 2.5": {
			a:    SmallInt(3),
			b:    Float(2.5),
			want: Float(15.588457268119894),
		},
		"Float 6 ** 1": {
			a:    SmallInt(6),
			b:    Float(1),
			want: Float(6),
		},
		"Float 4 ** -2": {
			a:    SmallInt(4),
			b:    Float(-2),
			want: Float(0.0625),
		},
		"Float 25 ** 0": {
			a:    SmallInt(25),
			b:    Float(0),
			want: Float(1),
		},
		"Float 25 ** NaN": {
			a:    25,
			b:    FloatNaN(),
			want: FloatNaN(),
		},
		"Float 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    0,
			b:    Float(-5),
			want: FloatInf(),
		},
		"Float 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    0,
			b:    FloatNegInf(),
			want: FloatInf(),
		},
		"Float 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    0,
			b:    FloatInf(),
			want: Float(0),
		},
		"Float 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    0,
			b:    Float(-8),
			want: FloatInf(),
		},
		"Float 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    0,
			b:    Float(7),
			want: Float(0),
		},
		"Float 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    0,
			b:    Float(8),
			want: Float(0),
		},
		"Float -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    -1,
			b:    FloatInf(),
			want: Float(1),
		},
		"Float -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    -1,
			b:    FloatNegInf(),
			want: Float(1),
		},
		"Float 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    2,
			b:    FloatInf(),
			want: FloatInf(),
		},
		"Float -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    -2,
			b:    FloatInf(),
			want: FloatInf(),
		},
		"Float 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    2,
			b:    FloatNegInf(),
			want: Float(0),
		},
		"Float -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    -2,
			b:    FloatNegInf(),
			want: Float(0),
		},

		"BigFloat 5 ** 2": {
			a:    SmallInt(5),
			b:    NewBigFloat(2),
			want: NewBigFloat(25).SetPrecision(64),
		},
		"BigFloat 7 ** 8": {
			a:    SmallInt(7),
			b:    NewBigFloat(8),
			want: NewBigFloat(5764801).SetPrecision(64),
		},
		"BigFloat 3 ** 2.5": {
			a:    SmallInt(3),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("15.5884572681198956415").SetPrecision(64),
		},
		"BigFloat 6 ** 1": {
			a:    SmallInt(6),
			b:    NewBigFloat(1),
			want: NewBigFloat(6).SetPrecision(64),
		},
		"BigFloat 4 ** -2": {
			a:    SmallInt(4),
			b:    NewBigFloat(-2),
			want: NewBigFloat(0.0625).SetPrecision(64),
		},
		"BigFloat 25 ** 0": {
			a:    SmallInt(25),
			b:    NewBigFloat(0),
			want: NewBigFloat(1).SetPrecision(64),
		},
		"BigFloat 25 ** NaN": {
			a:    25,
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
		"BigFloat 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    0,
			b:    NewBigFloat(-5),
			want: BigFloatInf(),
		},
		"BigFloat 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    0,
			b:    BigFloatNegInf(),
			want: BigFloatInf(),
		},
		"BigFloat 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    0,
			b:    BigFloatInf(),
			want: NewBigFloat(0).SetPrecision(64),
		},
		"BigFloat 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    0,
			b:    NewBigFloat(-8),
			want: BigFloatInf(),
		},
		"BigFloat 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    0,
			b:    NewBigFloat(7),
			want: NewBigFloat(0).SetPrecision(64),
		},
		"BigFloat 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    0,
			b:    NewBigFloat(8),
			want: NewBigFloat(0).SetPrecision(64),
		},
		"BigFloat -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    -1,
			b:    BigFloatInf(),
			want: NewBigFloat(1).SetPrecision(64),
		},
		"BigFloat -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    -1,
			b:    BigFloatNegInf(),
			want: NewBigFloat(1).SetPrecision(64),
		},
		"BigFloat 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    2,
			b:    BigFloatInf(),
			want: BigFloatInf(),
		},
		"BigFloat -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    -2,
			b:    BigFloatInf(),
			want: BigFloatInf(),
		},
		"BigFloat 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    2,
			b:    BigFloatNegInf(),
			want: NewBigFloat(0).SetPrecision(64),
		},
		"BigFloat -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    -2,
			b:    BigFloatNegInf(),
			want: NewBigFloat(0).SetPrecision(64),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Exponentiate(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				bigFloatComparer,
				floatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_RightBitshift(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"shift by String and return an error": {
			a:   SmallInt(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be used as a bitshift operand"),
		},
		"shift by Float and return an error": {
			a:   SmallInt(5),
			b:   Float(2.5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be used as a bitshift operand"),
		},

		"shift by SmallInt 5 >> 1": {
			a:    SmallInt(5),
			b:    SmallInt(1),
			want: SmallInt(2),
		},
		"shift by SmallInt 5 >> -1": {
			a:    SmallInt(5),
			b:    SmallInt(-1),
			want: SmallInt(10),
		},
		"shift by SmallInt 75 >> 0": {
			a:    SmallInt(75),
			b:    SmallInt(0),
			want: SmallInt(75),
		},
		"shift by SmallInt -32 >> 2": {
			a:    SmallInt(-32),
			b:    SmallInt(2),
			want: SmallInt(-8),
		},
		"shift by SmallInt 80 >> 60": {
			a:    SmallInt(80),
			b:    SmallInt(60),
			want: SmallInt(0),
		},
		"shift by SmallInt 80 >> -9223372036854775808": {
			a:    SmallInt(80),
			b:    SmallInt(-9223372036854775808),
			want: SmallInt(0),
		},
		"shift by SmallInt overflow": {
			a:    SmallInt(10),
			b:    SmallInt(-60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by SmallInt close to overflow": {
			a:    SmallInt(10),
			b:    SmallInt(-59),
			want: SmallInt(5764607523034234880),
		},

		"shift by BigInt 5 >> 1": {
			a:    SmallInt(5),
			b:    NewBigInt(1),
			want: SmallInt(2),
		},
		"shift by BigInt 5 >> -1": {
			a:    SmallInt(5),
			b:    NewBigInt(-1),
			want: SmallInt(10),
		},
		"shift by BigInt 75 >> 0": {
			a:    SmallInt(75),
			b:    NewBigInt(0),
			want: SmallInt(75),
		},
		"shift by BigInt 80 >> 60": {
			a:    SmallInt(80),
			b:    NewBigInt(60),
			want: SmallInt(0),
		},
		"shift by BigInt overflow": {
			a:    SmallInt(10),
			b:    NewBigInt(-60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by BigInt close to overflow": {
			a:    SmallInt(10),
			b:    NewBigInt(-59),
			want: SmallInt(5764607523034234880),
		},
		"shift by huge BigInt": {
			a:    SmallInt(10),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: SmallInt(0),
		},

		"shift by Int64 5 >> 1": {
			a:    SmallInt(5),
			b:    Int64(1),
			want: SmallInt(2),
		},
		"shift by Int64 5 >> -1": {
			a:    SmallInt(5),
			b:    Int64(-1),
			want: SmallInt(10),
		},
		"shift by Int64 75 >> 0": {
			a:    SmallInt(75),
			b:    Int64(0),
			want: SmallInt(75),
		},
		"shift by Int64 80 >> 60": {
			a:    SmallInt(80),
			b:    Int64(60),
			want: SmallInt(0),
		},
		"shift by Int64 overflow": {
			a:    SmallInt(10),
			b:    Int64(-60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by Int64 close to overflow": {
			a:    SmallInt(10),
			b:    Int64(-59),
			want: SmallInt(5764607523034234880),
		},

		"shift by Int32 5 >> 1": {
			a:    SmallInt(5),
			b:    Int32(1),
			want: SmallInt(2),
		},
		"shift by Int32 5 >> -1": {
			a:    SmallInt(5),
			b:    Int32(-1),
			want: SmallInt(10),
		},
		"shift by Int32 75 >> 0": {
			a:    SmallInt(75),
			b:    Int32(0),
			want: SmallInt(75),
		},
		"shift by Int32 80 >> 60": {
			a:    SmallInt(80),
			b:    Int32(60),
			want: SmallInt(0),
		},
		"shift by Int32 overflow": {
			a:    SmallInt(10),
			b:    Int32(-60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by Int32 close to overflow": {
			a:    SmallInt(10),
			b:    Int32(-59),
			want: SmallInt(5764607523034234880),
		},

		"shift by Int16 5 >> 1": {
			a:    SmallInt(5),
			b:    Int16(1),
			want: SmallInt(2),
		},
		"shift by Int16 5 >> -1": {
			a:    SmallInt(5),
			b:    Int16(-1),
			want: SmallInt(10),
		},
		"shift by Int16 75 >> 0": {
			a:    SmallInt(75),
			b:    Int16(0),
			want: SmallInt(75),
		},
		"shift by Int16 80 >> 60": {
			a:    SmallInt(80),
			b:    Int16(60),
			want: SmallInt(0),
		},
		"shift by Int16 overflow": {
			a:    SmallInt(10),
			b:    Int16(-60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by Int16 close to overflow": {
			a:    SmallInt(10),
			b:    Int16(-59),
			want: SmallInt(5764607523034234880),
		},

		"shift by Int8 5 >> 1": {
			a:    SmallInt(5),
			b:    Int8(1),
			want: SmallInt(2),
		},
		"shift by Int8 5 >> -1": {
			a:    SmallInt(5),
			b:    Int8(-1),
			want: SmallInt(10),
		},
		"shift by Int8 75 >> 0": {
			a:    SmallInt(75),
			b:    Int8(0),
			want: SmallInt(75),
		},
		"shift by Int8 80 >> 60": {
			a:    SmallInt(80),
			b:    Int8(60),
			want: SmallInt(0),
		},
		"shift by Int8 overflow": {
			a:    SmallInt(10),
			b:    Int8(-60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by Int8 close to overflow": {
			a:    SmallInt(10),
			b:    Int64(-59),
			want: SmallInt(5764607523034234880),
		},

		"shift by UInt64 5 >> 1": {
			a:    SmallInt(5),
			b:    UInt64(1),
			want: SmallInt(2),
		},
		"shift by UInt64 28 >> 2": {
			a:    SmallInt(28),
			b:    UInt64(2),
			want: SmallInt(7),
		},
		"shift by UInt64 75 >> 0": {
			a:    SmallInt(75),
			b:    UInt64(0),
			want: SmallInt(75),
		},
		"shift by UInt64 80 >> 60": {
			a:    SmallInt(80),
			b:    UInt64(60),
			want: SmallInt(0),
		},

		"shift by UInt32 5 >> 1": {
			a:    SmallInt(5),
			b:    UInt32(1),
			want: SmallInt(2),
		},
		"shift by UInt32 28 >> 2": {
			a:    SmallInt(28),
			b:    UInt32(2),
			want: SmallInt(7),
		},
		"shift by UInt32 75 >> 0": {
			a:    SmallInt(75),
			b:    UInt32(0),
			want: SmallInt(75),
		},
		"shift by UInt32 80 >> 60": {
			a:    SmallInt(80),
			b:    UInt32(60),
			want: SmallInt(0),
		},

		"shift by UInt16 5 >> 1": {
			a:    SmallInt(5),
			b:    UInt16(1),
			want: SmallInt(2),
		},
		"shift by UInt16 28 >> 2": {
			a:    SmallInt(28),
			b:    UInt16(2),
			want: SmallInt(7),
		},
		"shift by UInt16 75 >> 0": {
			a:    SmallInt(75),
			b:    UInt16(0),
			want: SmallInt(75),
		},
		"shift by UInt16 80 >> 60": {
			a:    SmallInt(80),
			b:    UInt16(60),
			want: SmallInt(0),
		},

		"shift by UInt8 5 >> 1": {
			a:    SmallInt(5),
			b:    UInt8(1),
			want: SmallInt(2),
		},
		"shift by UInt8 28 >> 2": {
			a:    SmallInt(28),
			b:    UInt8(2),
			want: SmallInt(7),
		},
		"shift by UInt8 75 >> 0": {
			a:    SmallInt(75),
			b:    UInt8(0),
			want: SmallInt(75),
		},
		"shift by UInt8 80 >> 60": {
			a:    SmallInt(80),
			b:    UInt8(60),
			want: SmallInt(0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.RightBitshift(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmpopts.IgnoreFields(BigFloat{}, "acc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_LeftBitshift(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"shift by String and return an error": {
			a:   SmallInt(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be used as a bitshift operand"),
		},
		"shift by Float and return an error": {
			a:   SmallInt(5),
			b:   Float(2.5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be used as a bitshift operand"),
		},

		"shift by SmallInt 5 << 1": {
			a:    SmallInt(5),
			b:    SmallInt(1),
			want: SmallInt(10),
		},
		"shift by SmallInt 12 << -1": {
			a:    SmallInt(12),
			b:    SmallInt(-1),
			want: SmallInt(6),
		},
		"shift by SmallInt 418 << -5": {
			a:    SmallInt(418),
			b:    SmallInt(-5),
			want: SmallInt(13),
		},
		"shift by SmallInt 75 << 0": {
			a:    SmallInt(75),
			b:    SmallInt(0),
			want: SmallInt(75),
		},
		"shift by SmallInt 80 << 56": {
			a:    SmallInt(80),
			b:    SmallInt(56),
			want: SmallInt(5764607523034234880),
		},
		"shift by SmallInt 80 >> -9223372036854775808": {
			a:    SmallInt(80),
			b:    SmallInt(-9223372036854775808),
			want: SmallInt(0),
		},
		"shift by SmallInt overflow": {
			a:    SmallInt(10),
			b:    SmallInt(60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by SmallInt close to overflow": {
			a:    SmallInt(10),
			b:    SmallInt(59),
			want: SmallInt(5764607523034234880),
		},

		"shift by BigInt 5 << 1": {
			a:    SmallInt(5),
			b:    NewBigInt(1),
			want: SmallInt(10),
		},
		"shift by BigInt 12 << -1": {
			a:    SmallInt(12),
			b:    NewBigInt(-1),
			want: SmallInt(6),
		},
		"shift by BigInt 418 << -5": {
			a:    SmallInt(418),
			b:    NewBigInt(-5),
			want: SmallInt(13),
		},
		"shift by BigInt 75 >> 0": {
			a:    SmallInt(75),
			b:    NewBigInt(0),
			want: SmallInt(75),
		},
		"shift by BigInt 80 << 56": {
			a:    SmallInt(80),
			b:    NewBigInt(56),
			want: SmallInt(5764607523034234880),
		},
		"shift by BigInt 80 >> -9223372036854775808": {
			a:    SmallInt(80),
			b:    NewBigInt(-9223372036854775808),
			want: SmallInt(0),
		},
		"shift by BigInt overflow": {
			a:    SmallInt(10),
			b:    NewBigInt(60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by BigInt close to overflow": {
			a:    SmallInt(10),
			b:    NewBigInt(59),
			want: SmallInt(5764607523034234880),
		},
		"shift by huge BigInt": {
			a:    SmallInt(10),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: SmallInt(0),
		},
		"shift by huge negative BigInt": {
			a:    SmallInt(10),
			b:    ParseBigIntPanic("-9223372036854775809", 10),
			want: SmallInt(0),
		},

		"shift by Int64 5 << 1": {
			a:    SmallInt(5),
			b:    Int64(1),
			want: SmallInt(10),
		},
		"shift by Int64 12 << -1": {
			a:    SmallInt(12),
			b:    Int64(-1),
			want: SmallInt(6),
		},
		"shift by Int64 418 << -5": {
			a:    SmallInt(418),
			b:    Int64(-5),
			want: SmallInt(13),
		},
		"shift by Int64 75 << 0": {
			a:    SmallInt(75),
			b:    Int64(0),
			want: SmallInt(75),
		},
		"shift by Int64 80 << 56": {
			a:    SmallInt(80),
			b:    Int64(56),
			want: SmallInt(5764607523034234880),
		},
		"shift by Int64 80 >> -9223372036854775808": {
			a:    SmallInt(80),
			b:    Int64(-9223372036854775808),
			want: SmallInt(0),
		},
		"shift by Int64 overflow": {
			a:    SmallInt(10),
			b:    Int64(60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by Int64 close to overflow": {
			a:    SmallInt(10),
			b:    Int64(59),
			want: SmallInt(5764607523034234880),
		},

		"shift by Int32 5 << 1": {
			a:    SmallInt(5),
			b:    Int32(1),
			want: SmallInt(10),
		},
		"shift by Int32 12 << -1": {
			a:    SmallInt(12),
			b:    Int32(-1),
			want: SmallInt(6),
		},
		"shift by Int32 418 << -5": {
			a:    SmallInt(418),
			b:    Int32(-5),
			want: SmallInt(13),
		},
		"shift by Int32 75 << 0": {
			a:    SmallInt(75),
			b:    Int32(0),
			want: SmallInt(75),
		},
		"shift by Int32 80 << 56": {
			a:    SmallInt(80),
			b:    Int32(56),
			want: SmallInt(5764607523034234880),
		},
		"shift by Int32 overflow": {
			a:    SmallInt(10),
			b:    Int32(60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by Int32 close to overflow": {
			a:    SmallInt(10),
			b:    Int32(59),
			want: SmallInt(5764607523034234880),
		},

		"shift by Int16 5 << 1": {
			a:    SmallInt(5),
			b:    Int16(1),
			want: SmallInt(10),
		},
		"shift by Int16 12 << -1": {
			a:    SmallInt(12),
			b:    Int16(-1),
			want: SmallInt(6),
		},
		"shift by Int16 418 << -5": {
			a:    SmallInt(418),
			b:    Int16(-5),
			want: SmallInt(13),
		},
		"shift by Int16 75 << 0": {
			a:    SmallInt(75),
			b:    Int16(0),
			want: SmallInt(75),
		},
		"shift by Int16 80 << 56": {
			a:    SmallInt(80),
			b:    Int16(56),
			want: SmallInt(5764607523034234880),
		},
		"shift by Int16 overflow": {
			a:    SmallInt(10),
			b:    Int16(60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by Int16 close to overflow": {
			a:    SmallInt(10),
			b:    Int16(59),
			want: SmallInt(5764607523034234880),
		},

		"shift by Int8 5 << 1": {
			a:    SmallInt(5),
			b:    Int8(1),
			want: SmallInt(10),
		},
		"shift by Int8 12 << -1": {
			a:    SmallInt(12),
			b:    Int8(-1),
			want: SmallInt(6),
		},
		"shift by Int8 418 << -5": {
			a:    SmallInt(418),
			b:    Int8(-5),
			want: SmallInt(13),
		},
		"shift by Int8 75 << 0": {
			a:    SmallInt(75),
			b:    Int8(0),
			want: SmallInt(75),
		},
		"shift by Int8 80 << 56": {
			a:    SmallInt(80),
			b:    Int8(56),
			want: SmallInt(5764607523034234880),
		},
		"shift by Int8 overflow": {
			a:    SmallInt(10),
			b:    Int8(60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by Int8 close to overflow": {
			a:    SmallInt(10),
			b:    Int8(59),
			want: SmallInt(5764607523034234880),
		},

		"shift by UInt64 5 << 1": {
			a:    SmallInt(5),
			b:    UInt64(1),
			want: SmallInt(10),
		},
		"shift by UInt64 75 << 0": {
			a:    SmallInt(75),
			b:    UInt64(0),
			want: SmallInt(75),
		},
		"shift by UInt64 80 << 56": {
			a:    SmallInt(80),
			b:    UInt64(56),
			want: SmallInt(5764607523034234880),
		},
		"shift by UIn64 overflow": {
			a:    SmallInt(10),
			b:    UInt64(60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by UInt64 close to overflow": {
			a:    SmallInt(10),
			b:    UInt64(59),
			want: SmallInt(5764607523034234880),
		},

		"shift by UInt32 5 << 1": {
			a:    SmallInt(5),
			b:    UInt32(1),
			want: SmallInt(10),
		},
		"shift by UInt32 75 << 0": {
			a:    SmallInt(75),
			b:    UInt32(0),
			want: SmallInt(75),
		},
		"shift by UInt32 80 << 56": {
			a:    SmallInt(80),
			b:    UInt32(56),
			want: SmallInt(5764607523034234880),
		},
		"shift by UIn32 overflow": {
			a:    SmallInt(10),
			b:    UInt32(60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by UInt32 close to overflow": {
			a:    SmallInt(10),
			b:    UInt32(59),
			want: SmallInt(5764607523034234880),
		},

		"shift by UInt16 5 << 1": {
			a:    SmallInt(5),
			b:    UInt16(1),
			want: SmallInt(10),
		},
		"shift by UInt16 75 << 0": {
			a:    SmallInt(75),
			b:    UInt16(0),
			want: SmallInt(75),
		},
		"shift by UInt16 80 << 56": {
			a:    SmallInt(80),
			b:    UInt16(56),
			want: SmallInt(5764607523034234880),
		},
		"shift by UIn16 overflow": {
			a:    SmallInt(10),
			b:    UInt16(60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by UInt16 close to overflow": {
			a:    SmallInt(10),
			b:    UInt16(59),
			want: SmallInt(5764607523034234880),
		},

		"shift by UInt8 5 << 1": {
			a:    SmallInt(5),
			b:    UInt8(1),
			want: SmallInt(10),
		},
		"shift by UInt8 75 << 0": {
			a:    SmallInt(75),
			b:    UInt8(0),
			want: SmallInt(75),
		},
		"shift by UInt8 80 << 56": {
			a:    SmallInt(80),
			b:    UInt8(56),
			want: SmallInt(5764607523034234880),
		},
		"shift by UIn8 overflow": {
			a:    SmallInt(10),
			b:    UInt8(60),
			want: ParseBigIntPanic("11529215046068469760", 10),
		},
		"shift by UInt8 close to overflow": {
			a:    SmallInt(10),
			b:    UInt8(59),
			want: SmallInt(5764607523034234880),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LeftBitshift(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmpopts.IgnoreFields(BigFloat{}, "acc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt & String and return an error": {
			a:   SmallInt(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::SmallInt`"),
		},
		"SmallInt & Int32 and return an error": {
			a:   SmallInt(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::SmallInt`"),
		},
		"SmallInt & Float and return an error": {
			a:   SmallInt(5),
			b:   Float(2.5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::SmallInt`"),
		},

		"23 & 10": {
			a:    SmallInt(23),
			b:    SmallInt(10),
			want: SmallInt(2),
		},
		"11 & 7": {
			a:    SmallInt(11),
			b:    SmallInt(7),
			want: SmallInt(3),
		},
		"-14 & 23": {
			a:    SmallInt(-14),
			b:    SmallInt(23),
			want: SmallInt(18),
		},
		"258 & 0": {
			a:    SmallInt(258),
			b:    SmallInt(0),
			want: SmallInt(0),
		},
		"124 & 255": {
			a:    SmallInt(124),
			b:    SmallInt(255),
			want: SmallInt(124),
		},

		"255 & 9223372036857247042": {
			a:    SmallInt(255),
			b:    ParseBigIntPanic("9223372036857247042", 10),
			want: SmallInt(66),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseAnd(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmpopts.IgnoreFields(BigFloat{}, "acc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt | String and return an error": {
			a:   SmallInt(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::SmallInt`"),
		},
		"SmallInt | Int32 and return an error": {
			a:   SmallInt(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::SmallInt`"),
		},
		"SmallInt | Float and return an error": {
			a:   SmallInt(5),
			b:   Float(2.5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::SmallInt`"),
		},

		"23 | 10": {
			a:    SmallInt(23),
			b:    SmallInt(10),
			want: SmallInt(31),
		},
		"11 | 7": {
			a:    SmallInt(11),
			b:    SmallInt(7),
			want: SmallInt(15),
		},
		"-14 | 23": {
			a:    SmallInt(-14),
			b:    SmallInt(23),
			want: SmallInt(-9),
		},
		"258 | 0": {
			a:    SmallInt(258),
			b:    SmallInt(0),
			want: SmallInt(258),
		},
		"124 | 255": {
			a:    SmallInt(124),
			b:    SmallInt(255),
			want: SmallInt(255),
		},

		"255 | 9223372036857247042": {
			a:    SmallInt(255),
			b:    ParseBigIntPanic("9223372036857247042", 10),
			want: ParseBigIntPanic("9223372036857247231", 10),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseOr(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmpopts.IgnoreFields(BigFloat{}, "acc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt ^ String and return an error": {
			a:   SmallInt(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::SmallInt`"),
		},
		"SmallInt ^ Int32 and return an error": {
			a:   SmallInt(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::SmallInt`"),
		},
		"SmallInt ^ Float and return an error": {
			a:   SmallInt(5),
			b:   Float(2.5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::SmallInt`"),
		},

		"23 ^ 10": {
			a:    SmallInt(23),
			b:    SmallInt(10),
			want: SmallInt(29),
		},
		"11 ^ 7": {
			a:    SmallInt(11),
			b:    SmallInt(7),
			want: SmallInt(12),
		},
		"-14 ^ 23": {
			a:    SmallInt(-14),
			b:    SmallInt(23),
			want: SmallInt(-27),
		},
		"258 ^ 0": {
			a:    SmallInt(258),
			b:    SmallInt(0),
			want: SmallInt(258),
		},
		"124 ^ 255": {
			a:    SmallInt(124),
			b:    SmallInt(255),
			want: SmallInt(131),
		},

		"255 ^ 9223372036857247042": {
			a:    SmallInt(255),
			b:    ParseBigIntPanic("9223372036857247042", 10),
			want: ParseBigIntPanic("9223372036857247165", 10),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseXor(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmpopts.IgnoreFields(BigFloat{}, "acc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    SmallInt
		b    Value
		want Value
		err  *Error
	}{
		"String and return an error": {
			a:   SmallInt(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::SmallInt`"),
		},

		"SmallInt 25 % 3": {
			a:    SmallInt(25),
			b:    SmallInt(3),
			want: SmallInt(1),
		},
		"SmallInt 76 % 6": {
			a:    SmallInt(76),
			b:    SmallInt(6),
			want: SmallInt(4),
		},
		"SmallInt -76 % 6": {
			a:    SmallInt(-76),
			b:    SmallInt(6),
			want: SmallInt(-4),
		},
		"SmallInt 76 % -6": {
			a:    SmallInt(76),
			b:    SmallInt(-6),
			want: SmallInt(4),
		},
		"SmallInt -76 % -6": {
			a:    SmallInt(-76),
			b:    SmallInt(-6),
			want: SmallInt(-4),
		},
		"SmallInt 124 % 9": {
			a:    SmallInt(124),
			b:    SmallInt(9),
			want: SmallInt(7),
		},

		"BigInt 25 % 3": {
			a:    SmallInt(25),
			b:    NewBigInt(3),
			want: SmallInt(1),
		},
		"BigInt 76 % 6": {
			a:    SmallInt(76),
			b:    NewBigInt(6),
			want: SmallInt(4),
		},
		"BigInt -76 % 6": {
			a:    SmallInt(-76),
			b:    NewBigInt(6),
			want: SmallInt(-4),
		},
		"BigInt 76 % -6": {
			a:    SmallInt(76),
			b:    NewBigInt(-6),
			want: SmallInt(4),
		},
		"BigInt -76 % -6": {
			a:    SmallInt(-76),
			b:    NewBigInt(-6),
			want: SmallInt(-4),
		},
		"BigInt 124 % 9": {
			a:    SmallInt(124),
			b:    NewBigInt(9),
			want: SmallInt(7),
		},
		"BigInt 9765 % 9223372036854775808": {
			a:    SmallInt(9765),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: SmallInt(9765),
		},

		"Float 25 % 3": {
			a:    SmallInt(25),
			b:    Float(3),
			want: Float(1),
		},
		"Float 76 % 6": {
			a:    SmallInt(76),
			b:    Float(6),
			want: Float(4),
		},
		"Float 124 % 9": {
			a:    SmallInt(124),
			b:    Float(9),
			want: Float(7),
		},
		"Float 124 % +Inf": {
			a:    SmallInt(124),
			b:    FloatInf(),
			want: Float(124),
		},
		"Float 124 % -Inf": {
			a:    SmallInt(124),
			b:    FloatNegInf(),
			want: Float(124),
		},
		"Float 74 % 6.25": {
			a:    SmallInt(74),
			b:    Float(6.25),
			want: Float(5.25),
		},
		"Float -74 % 6.25": {
			a:    SmallInt(-74),
			b:    Float(6.25),
			want: Float(-5.25),
		},
		"Float 74 % -6.25": {
			a:    SmallInt(74),
			b:    Float(-6.25),
			want: Float(5.25),
		},
		"Float -74 % -6.25": {
			a:    SmallInt(-74),
			b:    Float(-6.25),
			want: Float(-5.25),
		},
		"Float 25 % 0": { // Mod(x, 0) = NaN
			a:    SmallInt(25),
			b:    Float(0),
			want: FloatNaN(),
		},
		"Float 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    SmallInt(25),
			b:    FloatInf(),
			want: Float(25),
		},
		"Float -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    SmallInt(-87),
			b:    FloatNegInf(),
			want: Float(-87),
		},
		"Float 49 % NaN": { // Mod(x, NaN) = NaN
			a:    SmallInt(49),
			b:    FloatNaN(),
			want: FloatNaN(),
		},

		"BigFloat 25 % 3": {
			a:    SmallInt(25),
			b:    NewBigFloat(3),
			want: NewBigFloat(1).SetPrecision(64),
		},
		"BigFloat 76 % 6": {
			a:    SmallInt(76),
			b:    NewBigFloat(6),
			want: NewBigFloat(4).SetPrecision(64),
		},
		"BigFloat 124 % 9": {
			a:    SmallInt(124),
			b:    NewBigFloat(9),
			want: NewBigFloat(7).SetPrecision(64),
		},
		"BigFloat 74 % 6.25": {
			a:    SmallInt(74),
			b:    NewBigFloat(6.25),
			want: NewBigFloat(5.25).SetPrecision(64),
		},
		"BigFloat 74 % 6.25p92": {
			a:    SmallInt(74),
			b:    NewBigFloat(6.25).SetPrecision(92),
			want: NewBigFloat(5.25).SetPrecision(92),
		},
		"BigFloat -74 % 6.25": {
			a:    SmallInt(-74),
			b:    NewBigFloat(6.25),
			want: NewBigFloat(-5.25).SetPrecision(64),
		},
		"BigFloat 74 % -6.25": {
			a:    SmallInt(74),
			b:    NewBigFloat(-6.25),
			want: NewBigFloat(5.25).SetPrecision(64),
		},
		"BigFloat -74 % -6.25": {
			a:    SmallInt(-74),
			b:    NewBigFloat(-6.25),
			want: NewBigFloat(-5.25).SetPrecision(64),
		},
		"BigFloat 25 % 0": { // Mod(x, 0) = NaN
			a:    SmallInt(25),
			b:    NewBigFloat(0),
			want: BigFloatNaN(),
		},
		"BigFloat 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    SmallInt(25),
			b:    BigFloatInf(),
			want: NewBigFloat(25).SetPrecision(64),
		},
		"BigFloat -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    SmallInt(-87),
			b:    BigFloatNegInf(),
			want: NewBigFloat(-87).SetPrecision(64),
		},
		"BigFloat 49 % NaN": { // Mod(x, NaN) = NaN
			a:    SmallInt(49),
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Modulo(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
