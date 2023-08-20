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
			want: ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(10))),
		},
		"add BigInt and return BigInt": {
			a:    SmallInt(20),
			b:    ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(10))),
			want: ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(30))),
		},
		"add BigInt and return SmallInt": {
			a:    SmallInt(-20),
			b:    ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(10))),
			want: SmallInt(math.MaxInt64 - 10),
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
			want: ToElkBigInt((&big.Int{}).Sub(big.NewInt(math.MinInt64), big.NewInt(10))),
		},
		"subtract BigInt and return BigInt": {
			a:    SmallInt(5),
			b:    ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(10))),
			want: ToElkBigInt((&big.Int{}).Sub(big.NewInt(math.MinInt64), big.NewInt(4))),
		},
		"subtract BigInt and return SmallInt": {
			a:    SmallInt(20),
			b:    ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(10))),
			want: SmallInt(math.MinInt64 + 11),
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
			want: ToElkBigInt((&big.Int{}).Mul(big.NewInt(math.MaxInt64), big.NewInt(10))),
		},
		"multiply by BigInt and return BigInt": {
			a:    SmallInt(20),
			b:    ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(10))),
			want: ToElkBigInt((&big.Int{}).Mul(big.NewInt(20), (&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(10)))),
		},
		"multiply BigInt and return SmallInt": {
			a:    SmallInt(-1),
			b:    ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(1))),
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
			want: ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(1))),
		},
		"divide by BigInt and return SmallInt": {
			a:    SmallInt(20),
			b:    ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(10))),
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
