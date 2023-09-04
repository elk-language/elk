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
