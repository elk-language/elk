package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBigFloatAdd(t *testing.T) {
	tests := map[string]struct {
		left  *BigFloat
		right Value
		want  Value
		err   *Error
	}{
		"BigFloat + BigFloat => BigFloat": {
			left:  NewBigFloat(2.5),
			right: NewBigFloat(10.2),
			want:  NewBigFloat(12.7),
		},
		"result takes the max precision from its operands": {
			left:  NewBigFloat(2.5).SetPrecision(31),
			right: NewBigFloat(10.2).SetPrecision(54),
			want:  NewBigFloat(12.7).SetPrecision(54),
		},
		"result takes the max precision from its operands (left)": {
			left:  NewBigFloat(2.5).SetPrecision(54),
			right: NewBigFloat(10.2).SetPrecision(52),
			want:  NewBigFloat(12.7).SetPrecision(54),
		},
		"BigFloat + SmallInt => BigFloat": {
			left:  NewBigFloat(2.5),
			right: SmallInt(120),
			want:  NewBigFloat(122.5).SetPrecision(64),
		},
		"BigFloat + BigInt => BigFloat": {
			left:  NewBigFloat(2.5),
			right: NewBigInt(120),
			want:  NewBigFloat(122.5).SetPrecision(64),
		},
		"BigFloat + Int64 => TypeError": {
			left:  NewBigFloat(2.5),
			right: Int64(20),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat + String => TypeError": {
			left:  NewBigFloat(2.5),
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Add(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmp.AllowUnexported(BigFloat{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
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

func TestCountFloatDigits(t *testing.T) {
	tests := map[string]struct {
		str  string
		want int
	}{
		"int": {
			str:  "35",
			want: 2,
		},
		"float": {
			str:  "254.671",
			want: 6,
		},
		"int with exponent": {
			str:  "257e20",
			want: 3,
		},
		"float with exponent": {
			str:  "257.1223e91",
			want: 7,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := CountFloatDigits(tc.str)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloatSubtract(t *testing.T) {
	tests := map[string]struct {
		left  *BigFloat
		right Value
		want  Value
		err   *Error
	}{
		"BigFloat - BigFloat => BigFloat": {
			left:  NewBigFloat(10.0),
			right: NewBigFloat(2.5),
			want:  NewBigFloat(7.5),
		},
		"result takes the max precision from its operands": {
			left:  NewBigFloat(10.0).SetPrecision(54),
			right: NewBigFloat(2.5).SetPrecision(31),
			want:  NewBigFloat(7.5).SetPrecision(54),
		},
		"BigFloat - SmallInt => BigFloat": {
			left:  NewBigFloat(120.5),
			right: SmallInt(2),
			want:  NewBigFloat(118.5).SetPrecision(64),
		},
		"BigFloat - BigInt => BigFloat": {
			left:  NewBigFloat(120.5),
			right: NewBigInt(2),
			want:  NewBigFloat(118.5).SetPrecision(64),
		},
		"BigFloat - Int64 => TypeError": {
			left:  NewBigFloat(20.5),
			right: Int64(2),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat - String => TypeError": {
			left:  NewBigFloat(2.5),
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Subtract(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmp.AllowUnexported(BigFloat{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf("want: %s, got: %s\n%s", tc.want.Inspect(), got.Inspect(), diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloatMultiply(t *testing.T) {
	tests := map[string]struct {
		left  *BigFloat
		right Value
		want  Value
		err   *Error
	}{
		"BigFloat * BigFloat => BigFloat": {
			left:  NewBigFloat(2.55),
			right: NewBigFloat(10.0),
			want:  NewBigFloat(25.5),
		},
		"result takes the max precision from its operands": {
			left:  NewBigFloat(2.5).SetPrecision(31),
			right: NewBigFloat(10.0).SetPrecision(54),
			want:  NewBigFloat(25.0).SetPrecision(54),
		},
		"BigFloat * SmallInt => BigFloat": {
			left:  NewBigFloat(2.5),
			right: SmallInt(10),
			want:  NewBigFloat(25.0).SetPrecision(64),
		},
		"BigFloat * BigInt => BigFloat": {
			left:  NewBigFloat(2.5),
			right: NewBigInt(10),
			want:  NewBigFloat(25.0).SetPrecision(64),
		},
		"BigFloat * Int64 => TypeError": {
			left:  NewBigFloat(2.55),
			right: Int64(20),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat * String => TypeError": {
			left:  NewBigFloat(2.5),
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Multiply(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmp.AllowUnexported(BigFloat{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmpopts.IgnoreFields(BigFloat{}, "acc"),
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

func TestBigFloatDivide(t *testing.T) {
	tests := map[string]struct {
		left  *BigFloat
		right Value
		want  Value
		err   *Error
	}{
		"BigFloat / BigFloat => BigFloat": {
			left:  NewBigFloat(2.68),
			right: NewBigFloat(2),
			want:  NewBigFloat(1.34),
		},
		"result takes the max precision from its operands": {
			left:  NewBigFloat(2).SetPrecision(31),
			right: NewBigFloat(2).SetPrecision(54),
			want:  NewBigFloat(1).SetPrecision(54),
		},
		"BigFloat / SmallInt => BigFloat": {
			left:  NewBigFloat(2.68),
			right: SmallInt(2),
			want:  NewBigFloat(1.34).SetPrecision(64),
		},
		"BigFloat / BigInt => BigFloat": {
			left:  NewBigFloat(2.68),
			right: NewBigInt(2),
			want:  NewBigFloat(1.34).SetPrecision(64),
		},
		"BigFloat / Int64 => TypeError": {
			left:  NewBigFloat(2.68),
			right: Int64(2),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat / String => TypeError": {
			left:  NewBigFloat(2.5),
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Divide(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmp.AllowUnexported(BigFloat{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmpopts.IgnoreFields(BigFloat{}, "acc"),
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

func TestBigFloat_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    *BigFloat
		b    Value
		want Value
		err  *Error
	}{
		"exponentiate String and return an error": {
			a:   NewBigFloat(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
		"exponentiate Int32 and return an error": {
			a:   NewBigFloat(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::BigFloat`"),
		},
		"exponentiate positive SmallInt 5 ** 2": {
			a:    NewBigFloat(5),
			b:    SmallInt(2),
			want: NewBigFloat(25),
		},
		"exponentiate positive SmallInt 7 ** 8": {
			a:    NewBigFloat(7),
			b:    SmallInt(8),
			want: NewBigFloat(5764801),
		},
		"exponentiate positive SmallInt 2.5 ** 5": {
			a:    NewBigFloat(2.5),
			b:    SmallInt(5),
			want: NewBigFloat(97.65625),
		},
		"exponentiate positive SmallInt 7.12 ** 1": {
			a:    NewBigFloat(7.12),
			b:    SmallInt(1),
			want: NewBigFloat(7.12),
		},
		"exponentiate negative SmallInt": {
			a:    NewBigFloat(4),
			b:    SmallInt(-2),
			want: NewBigFloat(0.0625),
		},
		"exponentiate SmallInt zero": {
			a:    NewBigFloat(25),
			b:    SmallInt(0),
			want: NewBigFloat(1),
		},

		"exponentiate positive BigInt 5 ** 2": {
			a:    NewBigFloat(5),
			b:    NewBigInt(2),
			want: NewBigFloat(25),
		},
		"exponentiate positive BigInt 7 ** 8": {
			a:    NewBigFloat(7),
			b:    NewBigInt(8),
			want: NewBigFloat(5764801),
		},
		"exponentiate positive BigInt 2.5 ** 5": {
			a:    NewBigFloat(2.5),
			b:    NewBigInt(5),
			want: NewBigFloat(97.65625),
		},
		"exponentiate positive BigInt 7.12 ** 1": {
			a:    NewBigFloat(7.12),
			b:    NewBigInt(1),
			want: NewBigFloat(7.12),
		},
		"exponentiate negative BigInt": {
			a:    NewBigFloat(4),
			b:    NewBigInt(-2),
			want: NewBigFloat(0.0625),
		},
		"exponentiate BigInt zero": {
			a:    NewBigFloat(25),
			b:    NewBigInt(0),
			want: NewBigFloat(1),
		},

		"exponentiate positive Float 5 ** 2": {
			a:    NewBigFloat(5),
			b:    Float(2),
			want: NewBigFloat(25),
		},
		"exponentiate positive Float 7 ** 8": {
			a:    NewBigFloat(7),
			b:    Float(8),
			want: NewBigFloat(5764801),
		},
		"exponentiate positive Float 2.5 ** 2.5": {
			a:    NewBigFloat(2.5),
			b:    Float(2.5),
			want: NewBigFloat(9.882117688026186),
		},
		"exponentiate positive Float 3 ** 2.5": {
			a:    NewBigFloat(3),
			b:    Float(2.5),
			want: NewBigFloat(15.588457268119896),
		},
		"exponentiate positive Float 6 ** 1": {
			a:    NewBigFloat(6),
			b:    Float(1),
			want: NewBigFloat(6),
		},
		"exponentiate negative Float": {
			a:    NewBigFloat(4),
			b:    Float(-2),
			want: NewBigFloat(0.0625),
		},
		"exponentiate Float zero": {
			a:    NewBigFloat(25),
			b:    Float(0),
			want: NewBigFloat(1),
		},

		"exponentiate positive BigFloat 5 ** 2": {
			a:    NewBigFloat(5),
			b:    NewBigFloat(2),
			want: NewBigFloat(25).SetPrecision(53),
		},
		"exponentiate positive BigFloat 7 ** 8": {
			a:    NewBigFloat(7),
			b:    NewBigFloat(8),
			want: NewBigFloat(5764801).SetPrecision(53),
		},
		"exponentiate positive BigFloat 2.5 ** 2.5": {
			a:    NewBigFloat(2.5),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("9.882117688026186").SetPrecision(53),
		},
		"exponentiate positive BigFloat 3 ** 2.5": {
			a:    NewBigFloat(3),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("15.588457268119896").SetPrecision(53),
		},
		"exponentiate positive BigFloat 6 ** 1": {
			a:    NewBigFloat(6),
			b:    NewBigFloat(1),
			want: NewBigFloat(6).SetPrecision(53),
		},
		"exponentiate negative BigFloat": {
			a:    NewBigFloat(4),
			b:    NewBigFloat(-2),
			want: NewBigFloat(0.0625).SetPrecision(53),
		},
		"exponentiate BigFloat zero": {
			a:    NewBigFloat(25),
			b:    NewBigFloat(0),
			want: NewBigFloat(1).SetPrecision(53),
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
