package object

import (
	"testing"

	"math"
	"math/big"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBigInt_Add(t *testing.T) {
	tests := map[string]struct {
		a    *BigInt
		b    Value
		want Value
		err  *Error
	}{
		"add String and return an error": {
			a:   NewBigInt(3),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigInt`"),
		},
		"add SmallInt and return BigInt": {
			a:    ParseBigIntPanic("9223372036854775815", 10),
			b:    SmallInt(10),
			want: ParseBigIntPanic("9223372036854775825", 10),
		},
		"add SmallInt and return SmallInt": {
			a:    ParseBigIntPanic("9223372036854775837", 10),
			b:    SmallInt(-10),
			want: ParseBigIntPanic("9223372036854775827", 10),
		},
		"add BigInt and return BigInt": {
			a:    ParseBigIntPanic("9223372036854775827", 10),
			b:    NewBigInt(3),
			want: ParseBigIntPanic("9223372036854775830", 10),
		},
		"add BigInt and return SmallInt": {
			a:    ParseBigIntPanic("9223372036854775827", 10),
			b:    NewBigInt(-27),
			want: SmallInt(9223372036854775800),
		},
		"add Float and return Float": {
			a:    NewBigInt(3),
			b:    Float(2.5),
			want: Float(5.5),
		},
		"add BigFloat and return BigFloat with 64bit precision": {
			a:    NewBigInt(56),
			b:    NewBigFloat(2.5),
			want: NewBigFloat(58.5).SetPrecision(64),
		},
		"add BigFloat and return BigFloat with 80bit precision": {
			a:    NewBigInt(56),
			b:    NewBigFloat(2.5).SetPrecision(80),
			want: ToElkBigFloat((&big.Float{}).SetPrec(80).Add(big.NewFloat(56), big.NewFloat(2.5))),
		},
		"add BigFloat and return BigFloat with 65bit precision": {
			a:    ParseBigIntPanic("36893488147419103228", 10),
			b:    NewBigFloat(2.5).SetPrecision(64),
			want: ParseBigFloatPanic("36893488147419103230.5").SetPrecision(65),
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

func TestBigInt_Subtract(t *testing.T) {
	tests := map[string]struct {
		a    *BigInt
		b    Value
		want Value
		err  *Error
	}{
		"subtract String and return an error": {
			a:   ParseBigIntPanic("9223372036854775817", 10),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigInt`"),
		},
		"subtract SmallInt and return BigInt": {
			a:    ParseBigIntPanic("9223372036854775817", 10),
			b:    SmallInt(5),
			want: ParseBigIntPanic("9223372036854775812", 10),
		},
		"subtract SmallInt and return SmallInt": {
			a:    ParseBigIntPanic("9223372036854775817", 10),
			b:    SmallInt(11),
			want: SmallInt(9223372036854775806),
		},
		"subtract BigInt and return BigInt": {
			a:    ParseBigIntPanic("27670116110564327451", 10),
			b:    ParseBigIntPanic("9223372036854775817", 10),
			want: ParseBigIntPanic("18446744073709551634", 10),
		},
		"subtract BigInt and return SmallInt": {
			a:    ParseBigIntPanic("9223372036854775817", 10),
			b:    ParseBigIntPanic("9223372036854775812", 10),
			want: SmallInt(5),
		},
		"subtract Float and return Float": {
			a:    ParseBigIntPanic("9223372036854775817", 10),
			b:    Float(15.5),
			want: Float(9223372036854775801.5),
		},
		"subtract BigFloat and return BigFloat with 64bit precision": {
			a:    ParseBigIntPanic("9223372036854775817", 10),
			b:    NewBigFloat(854775817),
			want: NewBigFloat(9223372036000000000).SetPrecision(64),
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

func TestBigInt_Multiply(t *testing.T) {
	tests := map[string]struct {
		a    *BigInt
		b    Value
		want Value
		err  *Error
	}{
		"multiply by String and return an error": {
			a:   ParseBigIntPanic("9223372036854775817", 10),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigInt`"),
		},
		"multiply by SmallInt and return BigInt": {
			a:    ParseBigIntPanic("9223372036854775817", 10),
			b:    SmallInt(10),
			want: ParseBigIntPanic("92233720368547758170", 10),
		},
		"multiply by SmallInt and return SmallInt": {
			a:    ParseBigIntPanic("9223372036854775808", 10),
			b:    SmallInt(-1),
			want: SmallInt(-9223372036854775808),
		},
		"multiply by BigInt and return BigInt": {
			a:    ParseBigIntPanic("9223372036854775817", 10),
			b:    ParseBigIntPanic("9223372036854775825", 10),
			want: ParseBigIntPanic("85070591730234616105651324816166224025", 10),
		},
		"multiply by Float and return Float": {
			a:    ParseBigIntPanic("9223372036854775817", 10),
			b:    Float(0.00001),
			want: Float(92233720368547.77),
		},
		"multiply by BigFloat and return BigFloat with 64bit precision": {
			a:    ParseBigIntPanic("9223372036854775817", 10),
			b:    NewBigFloat(10),
			want: ParseBigFloatPanic("92233720368547758170").SetPrecision(64),
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

func TestBigInt_Divide(t *testing.T) {
	tests := map[string]struct {
		a    *BigInt
		b    Value
		want Value
		err  *Error
	}{
		"divide by String and return an error": {
			a:   ParseBigIntPanic("9223372036854775817", 10),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigInt`"),
		},
		"divide by SmallInt and return SmallInt": {
			a:    ParseBigIntPanic("9223372036854775818", 10),
			b:    SmallInt(2),
			want: SmallInt(4611686018427387909),
		},
		"divide by SmallInt and return BigInt": {
			a:    ParseBigIntPanic("27670116110564327454", 10),
			b:    SmallInt(2),
			want: ParseBigIntPanic("13835058055282163727", 10),
		},
		"divide by BigInt and return SmallInt": {
			a:    ParseBigIntPanic("27670116110564327454", 10),
			b:    ParseBigIntPanic("9223372036854775818", 10),
			want: SmallInt(3),
		},
		"divide by Float and return Float": {
			a:    ParseBigIntPanic("9223372036854775818", 10),
			b:    Float(2),
			want: Float(4611686018427387909),
		},
		"divide by BigFloat and return BigFloat with 64bit precision": {
			a:    ParseBigIntPanic("1000000000000000000", 10),
			b:    NewBigFloat(20000),
			want: NewBigFloat(50000000000000).SetPrecision(64),
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

func TestBigInt_IsSmallInt(t *testing.T) {
	tests := map[string]struct {
		i    *BigInt
		want bool
	}{
		"fits in SmallInt": {
			i:    NewBigInt(math.MaxInt64 - 1),
			want: true,
		},
		"does not fit in SmallInt": {
			i:    ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5))),
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.i.IsSmallInt()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigInt_ToSmallInt(t *testing.T) {
	tests := map[string]struct {
		i    *BigInt
		want SmallInt
	}{
		"fits in SmallInt": {
			i:    NewBigInt(math.MaxInt64 - 1),
			want: SmallInt(math.MaxInt64 - 1),
		},
		"overflows SmallInt": {
			i:    ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5))),
			want: math.MinInt64 + 4,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.i.ToSmallInt()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigInt_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    *BigInt
		b    Value
		want Value
		err  *Error
	}{
		"exponentiate String and return an error": {
			a:   NewBigInt(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigInt`"),
		},
		"exponentiate Int32 and return an error": {
			a:   NewBigInt(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::BigInt`"),
		},
		"exponentiate positive SmallInt 5 ** 2": {
			a:    NewBigInt(5),
			b:    SmallInt(2),
			want: SmallInt(25),
		},
		"exponentiate positive SmallInt 7 ** 8": {
			a:    NewBigInt(7),
			b:    SmallInt(8),
			want: SmallInt(5764801),
		},
		"exponentiate positive SmallInt 2 ** 5": {
			a:    NewBigInt(2),
			b:    SmallInt(5),
			want: SmallInt(32),
		},
		"exponentiate positive SmallInt 6 ** 1": {
			a:    NewBigInt(6),
			b:    SmallInt(1),
			want: SmallInt(6),
		},
		"exponentiate positive SmallInt and overflow": {
			a:    NewBigInt(2),
			b:    SmallInt(64),
			want: ParseBigIntPanic("18446744073709551616", 10),
		},
		"exponentiate a huge value by a positive SmallInt": {
			a:    ParseBigIntPanic("9223372036854775808", 10),
			b:    SmallInt(2),
			want: ParseBigIntPanic("85070591730234615865843651857942052864", 10),
		},
		"exponentiate negative SmallInt": {
			a:    NewBigInt(4),
			b:    SmallInt(-2),
			want: SmallInt(1),
		},
		"exponentiate SmallInt zero": {
			a:    NewBigInt(25),
			b:    SmallInt(0),
			want: SmallInt(1),
		},

		"exponentiate positive BigInt 5 ** 2": {
			a:    NewBigInt(5),
			b:    NewBigInt(2),
			want: SmallInt(25),
		},
		"exponentiate positive BigInt 7 ** 8": {
			a:    NewBigInt(7),
			b:    NewBigInt(8),
			want: SmallInt(5764801),
		},
		"exponentiate positive BigInt 2 ** 5": {
			a:    NewBigInt(2),
			b:    NewBigInt(5),
			want: SmallInt(32),
		},
		"exponentiate positive BigInt 6 ** 1": {
			a:    NewBigInt(6),
			b:    NewBigInt(1),
			want: SmallInt(6),
		},
		"exponentiate a huge value by a positive BigInt": {
			a:    ParseBigIntPanic("9223372036854775808", 10),
			b:    NewBigInt(2),
			want: ParseBigIntPanic("85070591730234615865843651857942052864", 10),
		},
		"exponentiate by a huge positive BigInt": {
			a:    NewBigInt(1),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: SmallInt(1),
		},
		"exponentiate positive BigInt and return BigInt": {
			a:    NewBigInt(2),
			b:    NewBigInt(64),
			want: ParseBigIntPanic("18446744073709551616", 10),
		},
		"exponentiate negative BigInt": {
			a:    NewBigInt(4),
			b:    NewBigInt(-2),
			want: SmallInt(1),
		},
		"exponentiate BigInt zero": {
			a:    NewBigInt(25),
			b:    NewBigInt(0),
			want: SmallInt(1),
		},

		"exponentiate positive Float 5 ** 2": {
			a:    NewBigInt(5),
			b:    Float(2),
			want: Float(25),
		},
		"exponentiate positive Float 7 ** 8": {
			a:    NewBigInt(7),
			b:    Float(8),
			want: Float(5764801),
		},
		"exponentiate positive Float 3 ** 2.5": {
			a:    NewBigInt(3),
			b:    Float(2.5),
			want: Float(15.588457268119894),
		},
		"exponentiate positive Float 6 ** 1": {
			a:    NewBigInt(6),
			b:    Float(1),
			want: Float(6),
		},
		"exponentiate a huge value by a positive Float": {
			a:    ParseBigIntPanic("9223372036854775808", 10),
			b:    Float(2),
			want: Float(8.507059173023462e+37),
		},
		"exponentiate negative Float": {
			a:    NewBigInt(4),
			b:    Float(-2),
			want: Float(0.0625),
		},
		"exponentiate Float zero": {
			a:    NewBigInt(25),
			b:    Float(0),
			want: Float(1),
		},

		"exponentiate positive BigFloat 5 ** 2": {
			a:    NewBigInt(5),
			b:    NewBigFloat(2),
			want: NewBigFloat(25).SetPrecision(64),
		},
		"exponentiate positive BigFloat 7 ** 8": {
			a:    NewBigInt(7),
			b:    NewBigFloat(8),
			want: NewBigFloat(5764801).SetPrecision(64),
		},
		"exponentiate positive BigFloat 3 ** 2.5": {
			a:    NewBigInt(3),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("15.5884572681198956415").SetPrecision(64),
		},
		"exponentiate positive BigFloat 6 ** 1": {
			a:    NewBigInt(6),
			b:    NewBigFloat(1),
			want: NewBigFloat(6).SetPrecision(64),
		},
		"exponentiate negative BigFloat": {
			a:    NewBigInt(4),
			b:    NewBigFloat(-2),
			want: NewBigFloat(0.0625).SetPrecision(64),
		},
		"exponentiate BigFloat zero": {
			a:    NewBigInt(25),
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
				t.Log(got.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
