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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Subtract(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Multiply(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
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
		"divide by Float and return Float": {
			a:    SmallInt(-20),
			b:    Float(0.5),
			want: Float(-40),
		},
		"divide by BigFloat and return BigFloat with 64bit precision": {
			a:    SmallInt(55),
			b:    NewBigFloat(2.5),
			want: NewBigFloat(22.0).SetPrecision(64),
		},
		"divide by BigFloat and return BigFloat with 80bit precision": {
			a:    SmallInt(55),
			b:    NewBigFloat(2.5).SetPrecision(80),
			want: ToElkBigFloat((&big.Float{}).SetPrec(80).Quo(big.NewFloat(55), big.NewFloat(2.5))),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Divide(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
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
		"exponentiate positive SmallInt 5 ** 2": {
			a:    SmallInt(5),
			b:    SmallInt(2),
			want: SmallInt(25),
		},
		"exponentiate positive SmallInt 7 ** 8": {
			a:    SmallInt(7),
			b:    SmallInt(8),
			want: SmallInt(5764801),
		},
		"exponentiate positive SmallInt 2 ** 5": {
			a:    SmallInt(2),
			b:    SmallInt(5),
			want: SmallInt(32),
		},
		"exponentiate positive SmallInt 6 ** 1": {
			a:    SmallInt(6),
			b:    SmallInt(1),
			want: SmallInt(6),
		},
		"exponentiate positive SmallInt and overflow": {
			a:    SmallInt(2),
			b:    SmallInt(64),
			want: ParseBigIntPanic("18446744073709551616", 10),
		},
		"exponentiate negative SmallInt": {
			a:    SmallInt(4),
			b:    SmallInt(-2),
			want: SmallInt(1),
		},
		"exponentiate SmallInt zero": {
			a:    SmallInt(25),
			b:    SmallInt(0),
			want: SmallInt(1),
		},

		"exponentiate positive BigInt 5 ** 2": {
			a:    SmallInt(5),
			b:    NewBigInt(2),
			want: SmallInt(25),
		},
		"exponentiate positive BigInt 7 ** 8": {
			a:    SmallInt(7),
			b:    NewBigInt(8),
			want: SmallInt(5764801),
		},
		"exponentiate positive BigInt 2 ** 5": {
			a:    SmallInt(2),
			b:    NewBigInt(5),
			want: SmallInt(32),
		},
		"exponentiate positive BigInt 6 ** 1": {
			a:    SmallInt(6),
			b:    NewBigInt(1),
			want: SmallInt(6),
		},
		"exponentiate positive BigInt and return BigInt": {
			a:    SmallInt(2),
			b:    NewBigInt(64),
			want: ParseBigIntPanic("18446744073709551616", 10),
		},
		"exponentiate negative BigInt": {
			a:    SmallInt(4),
			b:    NewBigInt(-2),
			want: SmallInt(1),
		},
		"exponentiate BigInt zero": {
			a:    SmallInt(25),
			b:    NewBigInt(0),
			want: SmallInt(1),
		},

		"exponentiate positive Float 5 ** 2": {
			a:    SmallInt(5),
			b:    Float(2),
			want: Float(25),
		},
		"exponentiate positive Float 7 ** 8": {
			a:    SmallInt(7),
			b:    Float(8),
			want: Float(5764801),
		},
		"exponentiate positive Float 3 ** 2.5": {
			a:    SmallInt(3),
			b:    Float(2.5),
			want: Float(15.588457268119894),
		},
		"exponentiate positive Float 6 ** 1": {
			a:    SmallInt(6),
			b:    Float(1),
			want: Float(6),
		},
		"exponentiate negative Float": {
			a:    SmallInt(4),
			b:    Float(-2),
			want: Float(0.0625),
		},
		"exponentiate Float zero": {
			a:    SmallInt(25),
			b:    Float(0),
			want: Float(1),
		},

		"exponentiate positive BigFloat 5 ** 2": {
			a:    SmallInt(5),
			b:    NewBigFloat(2),
			want: NewBigFloat(25).SetPrecision(64),
		},
		"exponentiate positive BigFloat 7 ** 8": {
			a:    SmallInt(7),
			b:    NewBigFloat(8),
			want: NewBigFloat(5764801).SetPrecision(64),
		},
		"exponentiate positive BigFloat 3 ** 2.5": {
			a:    SmallInt(3),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("15.5884572681198956415").SetPrecision(64),
		},
		"exponentiate positive BigFloat 6 ** 1": {
			a:    SmallInt(6),
			b:    NewBigFloat(1),
			want: NewBigFloat(6).SetPrecision(64),
		},
		"exponentiate negative BigFloat": {
			a:    SmallInt(4),
			b:    NewBigFloat(-2),
			want: NewBigFloat(0.0625).SetPrecision(64),
		},
		"exponentiate BigFloat zero": {
			a:    SmallInt(25),
			b:    NewBigFloat(0),
			want: NewBigFloat(1).SetPrecision(64),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Exponentiate(tc.b)
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
