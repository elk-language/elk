package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFloatAdd(t *testing.T) {
	tests := map[string]struct {
		left  Float
		right Value
		want  Value
		err   *Error
	}{
		"Float + Float => Float": {
			left:  2.5,
			right: Float(10.2),
			want:  Float(12.7),
		},
		"Float + BigFloat => BigFloat": {
			left:  2.5,
			right: NewBigFloat(10.2),
			want:  NewBigFloat(12.7),
		},
		"Float + SmallInt => Float": {
			left:  2.5,
			right: SmallInt(120),
			want:  Float(122.5),
		},
		"Float + BigInt => Float": {
			left:  2.5,
			right: NewBigInt(120),
			want:  Float(122.5),
		},
		"Float + Int64 => TypeError": {
			left:  2.5,
			right: Int64(20),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Float`"),
		},
		"Float + String => TypeError": {
			left:  2.5,
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
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

func TestFloatSubtract(t *testing.T) {
	tests := map[string]struct {
		left  Float
		right Value
		want  Value
		err   *Error
	}{
		"Float - Float => Float": {
			left:  10.0,
			right: Float(5.5),
			want:  Float(4.5),
		},
		"Float - BigFloat => BigFloat": {
			left:  12.5,
			right: NewBigFloat(2.5),
			want:  NewBigFloat(10.0),
		},
		"Float - SmallInt => Float": {
			left:  12.5,
			right: SmallInt(2),
			want:  Float(10.5),
		},
		"Float - BigInt => Float": {
			left:  2.5,
			right: NewBigInt(2),
			want:  Float(.5),
		},
		"Float - Int64 => TypeError": {
			left:  2.5,
			right: Int64(2),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Float`"),
		},
		"Float - String => TypeError": {
			left:  2.5,
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
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
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloatMultiply(t *testing.T) {
	tests := map[string]struct {
		left  Float
		right Value
		want  Value
		err   *Error
	}{
		"Float * Float => Float": {
			left:  2.55,
			right: Float(10.0),
			want:  Float(25.5),
		},
		"Float * BigFloat => BigFloat": {
			left:  2.55,
			right: NewBigFloat(10.0),
			want:  NewBigFloat(25.5),
		},
		"Float * SmallInt => Float": {
			left:  2.55,
			right: SmallInt(20),
			want:  Float(51),
		},
		"Float * BigInt => Float": {
			left:  2.55,
			right: NewBigInt(20),
			want:  Float(51),
		},
		"Float * Int64 => TypeError": {
			left:  2.5,
			right: Int64(20),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Float`"),
		},
		"Float * String => TypeError": {
			left:  2.5,
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
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

func TestFloatDivide(t *testing.T) {
	tests := map[string]struct {
		left  Float
		right Value
		want  Value
		err   *Error
	}{
		"Float / Float => Float": {
			left:  2.68,
			right: Float(2.0),
			want:  Float(1.34),
		},
		"Float / BigFloat => BigFloat": {
			left:  2.68,
			right: NewBigFloat(2.0),
			want:  NewBigFloat(1.34),
		},
		"Float / SmallInt => Float": {
			left:  2.68,
			right: SmallInt(2),
			want:  Float(1.34),
		},
		"Float / BigInt => Float": {
			left:  2.68,
			right: NewBigInt(2),
			want:  Float(1.34),
		},
		"Float / Int64 => TypeError": {
			left:  2.5,
			right: Int64(20),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Float`"),
		},
		"Float / String => TypeError": {
			left:  2.5,
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
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

func TestFloat_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    Float
		b    Value
		want Value
		err  *Error
	}{
		"exponentiate String and return an error": {
			a:   Float(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
		},
		"exponentiate Int32 and return an error": {
			a:   Float(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Float`"),
		},
		"exponentiate positive SmallInt 5 ** 2": {
			a:    Float(5),
			b:    SmallInt(2),
			want: Float(25),
		},
		"exponentiate positive SmallInt 7 ** 8": {
			a:    Float(7),
			b:    SmallInt(8),
			want: Float(5764801),
		},
		"exponentiate positive SmallInt 2.5 ** 5": {
			a:    Float(2.5),
			b:    SmallInt(5),
			want: Float(97.65625),
		},
		"exponentiate positive SmallInt 7.12 ** 1": {
			a:    Float(7.12),
			b:    SmallInt(1),
			want: Float(7.12),
		},
		"exponentiate negative SmallInt": {
			a:    Float(4),
			b:    SmallInt(-2),
			want: Float(0.0625),
		},
		"exponentiate SmallInt zero": {
			a:    Float(25),
			b:    SmallInt(0),
			want: Float(1),
		},

		"exponentiate positive BigInt 5 ** 2": {
			a:    Float(5),
			b:    NewBigInt(2),
			want: Float(25),
		},
		"exponentiate positive BigInt 7 ** 8": {
			a:    Float(7),
			b:    NewBigInt(8),
			want: Float(5764801),
		},
		"exponentiate positive BigInt 2.5 ** 5": {
			a:    Float(2.5),
			b:    NewBigInt(5),
			want: Float(97.65625),
		},
		"exponentiate positive BigInt 7.12 ** 1": {
			a:    Float(7.12),
			b:    NewBigInt(1),
			want: Float(7.12),
		},
		"exponentiate negative BigInt": {
			a:    Float(4),
			b:    NewBigInt(-2),
			want: Float(0.0625),
		},
		"exponentiate BigInt zero": {
			a:    Float(25),
			b:    NewBigInt(0),
			want: Float(1),
		},

		"exponentiate positive Float 5 ** 2": {
			a:    Float(5),
			b:    Float(2),
			want: Float(25),
		},
		"exponentiate positive Float 7 ** 8": {
			a:    Float(7),
			b:    Float(8),
			want: Float(5764801),
		},
		"exponentiate positive Float 2.5 ** 2.5": {
			a:    Float(2.5),
			b:    Float(2.5),
			want: Float(9.882117688026186),
		},
		"exponentiate positive Float 3 ** 2.5": {
			a:    Float(3),
			b:    Float(2.5),
			want: Float(15.588457268119894),
		},
		"exponentiate positive Float 6 ** 1": {
			a:    Float(6),
			b:    Float(1),
			want: Float(6),
		},
		"exponentiate negative Float": {
			a:    Float(4),
			b:    Float(-2),
			want: Float(0.0625),
		},
		"exponentiate Float zero": {
			a:    Float(25),
			b:    Float(0),
			want: Float(1),
		},

		"exponentiate positive BigFloat 5 ** 2": {
			a:    Float(5),
			b:    NewBigFloat(2),
			want: NewBigFloat(25).SetPrecision(53),
		},
		"exponentiate positive BigFloat 7 ** 8": {
			a:    Float(7),
			b:    NewBigFloat(8),
			want: NewBigFloat(5764801).SetPrecision(53),
		},
		"exponentiate positive BigFloat 2.5 ** 2.5": {
			a:    Float(2.5),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("9.882117688026186").SetPrecision(53),
		},
		"exponentiate positive BigFloat 3 ** 2.5": {
			a:    Float(3),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("15.5884572681198956415").SetPrecision(53),
		},
		"exponentiate positive BigFloat 6 ** 1": {
			a:    Float(6),
			b:    NewBigFloat(1),
			want: NewBigFloat(6).SetPrecision(53),
		},
		"exponentiate negative BigFloat": {
			a:    Float(4),
			b:    NewBigFloat(-2),
			want: NewBigFloat(0.0625).SetPrecision(53),
		},
		"exponentiate BigFloat zero": {
			a:    Float(25),
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
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    Float
		b    Value
		want Value
		err  *Error
	}{
		"mod by String and return an error": {
			a:   Float(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
		},

		"mod by SmallInt 25 % 3": {
			a:    Float(25),
			b:    SmallInt(3),
			want: Float(1),
		},
		"mod by SmallInt 25.6 % 3": {
			a:    Float(25.6),
			b:    SmallInt(3),
			want: Float(1.6000000000000014),
		},
		"mod by SmallInt 76 % 6": {
			a:    Float(76),
			b:    SmallInt(6),
			want: Float(4),
		},
		"mod by SmallInt -76 % 6": {
			a:    Float(-76),
			b:    SmallInt(6),
			want: Float(-4),
		},
		"mod by SmallInt 76 % -6": {
			a:    Float(76),
			b:    SmallInt(-6),
			want: Float(4),
		},
		"mod by SmallInt -76 % -6": {
			a:    Float(-76),
			b:    SmallInt(-6),
			want: Float(-4),
		},
		"mod by SmallInt 124 % 9": {
			a:    Float(124),
			b:    SmallInt(9),
			want: Float(7),
		},

		"mod by BigInt 25 % 3": {
			a:    Float(25),
			b:    NewBigInt(3),
			want: Float(1),
		},
		"mod by BigInt 76 % 6": {
			a:    Float(76),
			b:    NewBigInt(6),
			want: Float(4),
		},
		"mod by BigInt 76.5 % 6": {
			a:    Float(76.5),
			b:    NewBigInt(6),
			want: Float(4.5),
		},
		"mod by BigInt -76 % 6": {
			a:    Float(-76),
			b:    NewBigInt(6),
			want: Float(-4),
		},
		"mod by BigInt 76 % -6": {
			a:    Float(76),
			b:    NewBigInt(-6),
			want: Float(4),
		},
		"mod by BigInt -76 % -6": {
			a:    Float(-76),
			b:    NewBigInt(-6),
			want: Float(-4),
		},
		"mod by BigInt 124 % 9": {
			a:    Float(124),
			b:    NewBigInt(9),
			want: Float(7),
		},
		"mod by BigInt 9765 % 9223372036854775808": {
			a:    Float(9765),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: Float(9765),
		},

		"mod by Float 25 % 3": {
			a:    Float(25),
			b:    Float(3),
			want: Float(1),
		},
		"mod by Float 76 % 6": {
			a:    Float(76),
			b:    Float(6),
			want: Float(4),
		},
		"mod by Float 124 % 9": {
			a:    Float(124),
			b:    Float(9),
			want: Float(7),
		},
		"mod by Float 124 % +Inf": {
			a:    Float(124),
			b:    FloatInf(),
			want: Float(124),
		},
		"mod by Float 124 % -Inf": {
			a:    Float(124),
			b:    FloatNegInf(),
			want: Float(124),
		},
		"mod by Float 74.5 % 6.25": {
			a:    Float(74.5),
			b:    Float(6.25),
			want: Float(5.75),
		},
		"mod by Float 74 % 6.25": {
			a:    Float(74),
			b:    Float(6.25),
			want: Float(5.25),
		},
		"mod by Float -74 % 6.25": {
			a:    Float(-74),
			b:    Float(6.25),
			want: Float(-5.25),
		},
		"mod by Float 74 % -6.25": {
			a:    Float(74),
			b:    Float(-6.25),
			want: Float(5.25),
		},
		"mod by Float -74 % -6.25": {
			a:    Float(-74),
			b:    Float(-6.25),
			want: Float(-5.25),
		},

		"mod by BigFloat 25 % 3": {
			a:    Float(25),
			b:    NewBigFloat(3),
			want: NewBigFloat(1),
		},
		"mod by BigFloat 76 % 6": {
			a:    Float(76),
			b:    NewBigFloat(6),
			want: NewBigFloat(4),
		},
		"mod by BigFloat 124 % 9": {
			a:    Float(124),
			b:    NewBigFloat(9),
			want: NewBigFloat(7),
		},
		"mod by BigFloat 74 % 6.25": {
			a:    Float(74),
			b:    NewBigFloat(6.25),
			want: NewBigFloat(5.25),
		},
		"mod by BigFloat 74 % 6.25 with higher precision": {
			a:    Float(74),
			b:    NewBigFloat(6.25).SetPrecision(64),
			want: NewBigFloat(5.25).SetPrecision(64),
		},
		"mod by BigFloat -74 % 6.25": {
			a:    Float(-74),
			b:    NewBigFloat(6.25),
			want: NewBigFloat(-5.25),
		},
		"mod by BigFloat 74 % -6.25": {
			a:    Float(74),
			b:    NewBigFloat(-6.25),
			want: NewBigFloat(5.25),
		},
		"mod by BigFloat -74 % -6.25": {
			a:    Float(-74),
			b:    NewBigFloat(-6.25),
			want: NewBigFloat(-5.25),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Modulo(tc.b)
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
