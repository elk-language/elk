package object

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

func TestSmallInt_AddOverflow(t *testing.T) {
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
			got, ok := tc.a.AddOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_SubtractOverflow(t *testing.T) {
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
			got, ok := tc.a.SubtractOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_MultiplyOverflow(t *testing.T) {
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
			got, ok := tc.a.MultiplyOverflow(tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSmallInt_DivideOverflow(t *testing.T) {
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
			got, ok := tc.a.DivideOverflow(tc.b)
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
