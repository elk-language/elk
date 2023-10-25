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
		"Float + Float NaN => Float NaN": {
			left:  2.5,
			right: FloatNaN(),
			want:  FloatNaN(),
		},
		"Float NaN + Float => Float NaN": {
			left:  FloatNaN(),
			right: Float(2.5),
			want:  FloatNaN(),
		},
		"Float NaN + Float NaN => Float NaN": {
			left:  FloatNaN(),
			right: FloatNaN(),
			want:  FloatNaN(),
		},
		"Float +Inf + Float +Inf => Float +Inf": {
			left:  FloatInf(),
			right: FloatInf(),
			want:  FloatInf(),
		},
		"Float -Inf + Float -Inf => Float -Inf": {
			left:  FloatNegInf(),
			right: FloatNegInf(),
			want:  FloatNegInf(),
		},
		"Float +Inf + Float -Inf => Float NaN": {
			left:  FloatInf(),
			right: FloatNegInf(),
			want:  FloatNaN(),
		},
		"Float + BigFloat => BigFloat": {
			left:  2.5,
			right: NewBigFloat(10.2),
			want:  NewBigFloat(12.7),
		},
		"Float NaN + BigFloat => BigFloat NaN": {
			left:  FloatNaN(),
			right: NewBigFloat(10.2),
			want:  BigFloatNaN(),
		},
		"Float + BigFloat NaN => BigFloat NaN": {
			left:  2.5,
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"Float NaN + BigFloat NaN => BigFloat NaN": {
			left:  FloatNaN(),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"Float +Inf + BigFloat -Inf => BigFloat NaN": {
			left:  FloatInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNaN(),
		},
		"Float +Inf + BigFloat +Inf => BigFloat +Inf": {
			left:  FloatInf(),
			right: BigFloatInf(),
			want:  BigFloatInf(),
		},
		"Float -Inf + BigFloat -Inf => BigFloat -Inf": {
			left:  FloatNegInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNegInf(),
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
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				BigFloatComparer,
				FloatComparer,
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

		"Float - Float NaN => Float NaN": {
			left:  2.5,
			right: FloatNaN(),
			want:  FloatNaN(),
		},
		"Float NaN - Float => Float NaN": {
			left:  FloatNaN(),
			right: Float(2.5),
			want:  FloatNaN(),
		},
		"Float NaN - Float NaN => Float NaN": {
			left:  FloatNaN(),
			right: FloatNaN(),
			want:  FloatNaN(),
		},
		"Float +Inf - Float +Inf => Float NaN": {
			left:  FloatInf(),
			right: FloatInf(),
			want:  FloatNaN(),
		},
		"Float -Inf - Float -Inf => Float NaN": {
			left:  FloatNegInf(),
			right: FloatNegInf(),
			want:  FloatNaN(),
		},
		"Float +Inf - Float -Inf => Float +Inf": {
			left:  FloatInf(),
			right: FloatNegInf(),
			want:  FloatInf(),
		},

		"Float - BigFloat NaN => BigFloat NaN": {
			left:  2.5,
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"Float NaN - BigFloat => BigFloat NaN": {
			left:  FloatNaN(),
			right: NewBigFloat(2.5),
			want:  BigFloatNaN(),
		},
		"Float NaN - BigFloat NaN => BigFloat NaN": {
			left:  FloatNaN(),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"Float +Inf - BigFloat +Inf => BigFloat NaN": {
			left:  FloatInf(),
			right: BigFloatInf(),
			want:  BigFloatNaN(),
		},
		"Float -Inf - BigFloat -Inf => BigFloat NaN": {
			left:  FloatNegInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNaN(),
		},
		"Float +Inf - BigFloat -Inf => BigFloat +Inf": {
			left:  FloatInf(),
			right: BigFloatNegInf(),
			want:  BigFloatInf(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Subtract(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				FloatComparer,
				BigFloatComparer,
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

		"Float * BigFloat NaN => BigFloat NaN": {
			left:  Float(2.5),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"Float NaN * BigFloat => BigFloat NaN": {
			left:  FloatNaN(),
			right: NewBigFloat(10.2),
			want:  BigFloatNaN(),
		},
		"Float NaN * BigFloat NaN => BigFloat NaN": {
			left:  FloatNaN(),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"Float +Inf * BigFloat => BigFloat +Inf": {
			left:  FloatInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatInf(),
		},
		"Float * BigFloat +Inf => BigFloat +Inf": {
			left:  Float(10.2),
			right: BigFloatInf(),
			want:  BigFloatInf(),
		},
		"Float +Inf * BigFloat +Inf => BigFloat +Inf": {
			left:  FloatInf(),
			right: BigFloatInf(),
			want:  BigFloatInf(),
		},
		"Float -Inf * +BigFloat => BigFloat -Inf": {
			left:  FloatNegInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatNegInf(),
		},
		"Float -Inf * -BigFloat => BigFloat +Inf": {
			left:  FloatNegInf(),
			right: NewBigFloat(-10.2),
			want:  BigFloatInf(),
		},
		"+Float * BigFloat -Inf => BigFloat -Inf": {
			left:  Float(10.2),
			right: BigFloatNegInf(),
			want:  BigFloatNegInf(),
		},
		"-Float * BigFloat -Inf => BigFloat +Inf": {
			left:  Float(-10.2),
			right: BigFloatNegInf(),
			want:  BigFloatInf(),
		},
		"Float -Inf * BigFloat -Inf => BigFloat +Inf": {
			left:  FloatNegInf(),
			right: BigFloatNegInf(),
			want:  BigFloatInf(),
		},
		"Float +Inf * BigFloat -Inf => BigFloat -Inf": {
			left:  FloatInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNegInf(),
		},
		"Float -Inf * BigFloat +Inf => BigFloat -Inf": {
			left:  FloatNegInf(),
			right: BigFloatInf(),
			want:  BigFloatNegInf(),
		},
		"Float -Inf * BigFloat 0 => BigFloat NaN": {
			left:  FloatNegInf(),
			right: NewBigFloat(0),
			want:  BigFloatNaN(),
		},
		"Float 0 * BigFloat +Inf => BigFloat NaN": {
			left:  Float(0),
			right: BigFloatInf(),
			want:  BigFloatNaN(),
		},

		"Float * Float NaN => Float NaN": {
			left:  Float(2.5),
			right: FloatNaN(),
			want:  FloatNaN(),
		},
		"Float NaN * Float => Float NaN": {
			left:  FloatNaN(),
			right: Float(10.2),
			want:  FloatNaN(),
		},
		"Float NaN * Float NaN => Float NaN": {
			left:  FloatNaN(),
			right: FloatNaN(),
			want:  FloatNaN(),
		},
		"Float +Inf * Float => Float +Inf": {
			left:  FloatInf(),
			right: Float(10.2),
			want:  FloatInf(),
		},
		"Float * Float +Inf => Float +Inf": {
			left:  Float(10.2),
			right: FloatInf(),
			want:  FloatInf(),
		},
		"Float +Inf * Float +Inf => Float +Inf": {
			left:  FloatInf(),
			right: FloatInf(),
			want:  FloatInf(),
		},
		"Float -Inf * +Float => Float -Inf": {
			left:  FloatNegInf(),
			right: Float(10.2),
			want:  FloatNegInf(),
		},
		"Float -Inf * -Float => Float +Inf": {
			left:  FloatNegInf(),
			right: Float(-10.2),
			want:  FloatInf(),
		},
		"+Float * Float -Inf => Float -Inf": {
			left:  Float(10.2),
			right: FloatNegInf(),
			want:  FloatNegInf(),
		},
		"-Float * Float -Inf => Float +Inf": {
			left:  Float(-10.2),
			right: FloatNegInf(),
			want:  FloatInf(),
		},
		"Float -Inf * Float -Inf => Float +Inf": {
			left:  FloatNegInf(),
			right: FloatNegInf(),
			want:  FloatInf(),
		},
		"Float +Inf * Float -Inf => Float -Inf": {
			left:  FloatInf(),
			right: FloatNegInf(),
			want:  FloatNegInf(),
		},
		"Float -Inf * Float +Inf => Float -Inf": {
			left:  FloatNegInf(),
			right: FloatInf(),
			want:  FloatNegInf(),
		},
		"Float -Inf * Float 0 => Float NaN": {
			left:  FloatNegInf(),
			right: Float(0),
			want:  FloatNaN(),
		},
		"Float 0 * Float +Inf => Float NaN": {
			left:  Float(0),
			right: FloatInf(),
			want:  FloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Multiply(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				BigFloatComparer,
				FloatComparer,
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

		"Float / BigFloat NaN => BigFloat NaN": {
			left:  Float(2.5),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"Float NaN / BigFloat => BigFloat NaN": {
			left:  FloatNaN(),
			right: NewBigFloat(10.2),
			want:  BigFloatNaN(),
		},
		"Float NaN / BigFloat NaN => BigFloat NaN": {
			left:  FloatNaN(),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"Float +Inf / BigFloat => BigFloat +Inf": {
			left:  FloatInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatInf(),
		},
		"Float / BigFloat +Inf => BigFloat 0": {
			left:  Float(10.2),
			right: BigFloatInf(),
			want:  NewBigFloat(0),
		},
		"Float +Inf / BigFloat +Inf => BigFloat NaN": {
			left:  FloatInf(),
			right: BigFloatInf(),
			want:  BigFloatNaN(),
		},
		"Float -Inf / +BigFloat => BigFloat -Inf": {
			left:  FloatNegInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatNegInf(),
		},
		"Float -Inf / -BigFloat => BigFloat +Inf": {
			left:  FloatNegInf(),
			right: NewBigFloat(-10.2),
			want:  BigFloatInf(),
		},
		"+Float / BigFloat -Inf => BigFloat -0": {
			left:  Float(10.2),
			right: BigFloatNegInf(),
			want:  NewBigFloat(-0),
		},
		"-Float / BigFloat -Inf => BigFloat +0": {
			left:  Float(-10.2),
			right: BigFloatNegInf(),
			want:  NewBigFloat(0),
		},
		"Float -Inf / BigFloat -Inf => BigFloat NaN": {
			left:  FloatNegInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNaN(),
		},
		"Float +Inf / BigFloat -Inf => BigFloat NaN": {
			left:  FloatInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNaN(),
		},
		"Float -Inf / BigFloat +Inf => BigFloat NaN": {
			left:  FloatNegInf(),
			right: BigFloatInf(),
			want:  BigFloatNaN(),
		},
		"Float -Inf / BigFloat 0 => BigFloat -Inf": {
			left:  FloatNegInf(),
			right: NewBigFloat(0),
			want:  BigFloatNegInf(),
		},
		"Float 0 / BigFloat +Inf => BigFloat 0": {
			left:  Float(0),
			right: BigFloatInf(),
			want:  NewBigFloat(0),
		},

		"Float / Float NaN => Float NaN": {
			left:  Float(2.5),
			right: FloatNaN(),
			want:  FloatNaN(),
		},
		"Float NaN / Float => Float NaN": {
			left:  FloatNaN(),
			right: Float(10.2),
			want:  FloatNaN(),
		},
		"Float NaN / Float NaN => Float NaN": {
			left:  FloatNaN(),
			right: FloatNaN(),
			want:  FloatNaN(),
		},
		"Float +Inf / Float => Float +Inf": {
			left:  FloatInf(),
			right: Float(10.2),
			want:  FloatInf(),
		},
		"Float / Float +Inf => Float 0": {
			left:  Float(10.2),
			right: FloatInf(),
			want:  Float(0),
		},
		"Float +Inf / Float +Inf => Float NaN": {
			left:  FloatInf(),
			right: FloatInf(),
			want:  FloatNaN(),
		},
		"Float -Inf / +Float => Float -Inf": {
			left:  FloatNegInf(),
			right: Float(10.2),
			want:  FloatNegInf(),
		},
		"Float -Inf / -Float => Float +Inf": {
			left:  FloatNegInf(),
			right: Float(-10.2),
			want:  FloatInf(),
		},
		"+Float / Float -Inf => Float -0": {
			left:  Float(10.2),
			right: FloatNegInf(),
			want:  Float(-0),
		},
		"-Float / Float -Inf => Float +0": {
			left:  Float(-10.2),
			right: FloatNegInf(),
			want:  Float(0),
		},
		"Float -Inf / Float -Inf => Float NaN": {
			left:  FloatNegInf(),
			right: FloatNegInf(),
			want:  FloatNaN(),
		},
		"Float +Inf / Float -Inf => Float NaN": {
			left:  FloatInf(),
			right: FloatNegInf(),
			want:  FloatNaN(),
		},
		"Float -Inf / Float +Inf => Float NaN": {
			left:  FloatNegInf(),
			right: FloatInf(),
			want:  FloatNaN(),
		},
		"Float -Inf / Float 0 => Float -Inf": {
			left:  FloatNegInf(),
			right: Float(0),
			want:  FloatNegInf(),
		},
		"Float 0 / Float +Inf => Float 0": {
			left:  Float(0),
			right: FloatInf(),
			want:  Float(0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Divide(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				BigFloatComparer,
				FloatComparer,
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
		"SmallInt 5 ** 2": {
			a:    Float(5),
			b:    SmallInt(2),
			want: Float(25),
		},
		"SmallInt 7 ** 8": {
			a:    Float(7),
			b:    SmallInt(8),
			want: Float(5764801),
		},
		"SmallInt 2.5 ** 5": {
			a:    Float(2.5),
			b:    SmallInt(5),
			want: Float(97.65625),
		},
		"SmallInt 7.12 ** 1": {
			a:    Float(7.12),
			b:    SmallInt(1),
			want: Float(7.12),
		},
		"SmallInt 4 ** -2": {
			a:    Float(4),
			b:    SmallInt(-2),
			want: Float(0.0625),
		},
		"SmallInt 25 ** 0": {
			a:    Float(25),
			b:    SmallInt(0),
			want: Float(1),
		},

		"BigInt 5 ** 2": {
			a:    Float(5),
			b:    NewBigInt(2),
			want: Float(25),
		},
		"BigInt 7 ** 8": {
			a:    Float(7),
			b:    NewBigInt(8),
			want: Float(5764801),
		},
		"BigInt 2.5 ** 5": {
			a:    Float(2.5),
			b:    NewBigInt(5),
			want: Float(97.65625),
		},
		"BigInt 7.12 ** 1": {
			a:    Float(7.12),
			b:    NewBigInt(1),
			want: Float(7.12),
		},
		"BigInt 4 ** -2": {
			a:    Float(4),
			b:    NewBigInt(-2),
			want: Float(0.0625),
		},
		"BigInt 25 ** 0": {
			a:    Float(25),
			b:    NewBigInt(0),
			want: Float(1),
		},

		"Float 5 ** 2": {
			a:    Float(5),
			b:    Float(2),
			want: Float(25),
		},
		"Float 7 ** 8": {
			a:    Float(7),
			b:    Float(8),
			want: Float(5764801),
		},
		"Float 2.5 ** 2.5": {
			a:    Float(2.5),
			b:    Float(2.5),
			want: Float(9.882117688026186),
		},
		"Float 3 ** 2.5": {
			a:    Float(3),
			b:    Float(2.5),
			want: Float(15.588457268119894),
		},
		"Float 6 ** 1": {
			a:    Float(6),
			b:    Float(1),
			want: Float(6),
		},
		"Float 4 ** -2": {
			a:    Float(4),
			b:    Float(-2),
			want: Float(0.0625),
		},
		"Float 25 ** 0": {
			a:    Float(25),
			b:    Float(0),
			want: Float(1),
		},
		"Float 25 ** NaN": {
			a:    Float(25),
			b:    FloatNaN(),
			want: FloatNaN(),
		},
		"Float NaN ** 25": {
			a:    FloatNaN(),
			b:    Float(25),
			want: FloatNaN(),
		},
		"Float NaN ** NaN": {
			a:    FloatNaN(),
			b:    FloatNaN(),
			want: FloatNaN(),
		},
		"Float 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    Float(0),
			b:    Float(-5),
			want: FloatInf(),
		},
		"Float 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    Float(0),
			b:    FloatNegInf(),
			want: FloatInf(),
		},
		"Float 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    Float(0),
			b:    FloatInf(),
			want: Float(0),
		},
		"Float 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    Float(0),
			b:    Float(-8),
			want: FloatInf(),
		},
		"Float 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    Float(0),
			b:    Float(7),
			want: Float(0),
		},
		"Float 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    Float(0),
			b:    Float(8),
			want: Float(0),
		},
		"Float -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    Float(-1),
			b:    FloatInf(),
			want: Float(1),
		},
		"Float -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    Float(-1),
			b:    FloatNegInf(),
			want: Float(1),
		},
		"Float 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    Float(2),
			b:    FloatInf(),
			want: FloatInf(),
		},
		"Float -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    Float(-2),
			b:    FloatInf(),
			want: FloatInf(),
		},
		"Float 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    Float(2),
			b:    FloatNegInf(),
			want: Float(0),
		},
		"Float -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    Float(-2),
			b:    FloatNegInf(),
			want: Float(0),
		},
		"Float 0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    Float(0.5),
			b:    FloatInf(),
			want: Float(0),
		},
		"Float -0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    Float(-0.5),
			b:    FloatInf(),
			want: Float(0),
		},
		"Float 0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    Float(0.5),
			b:    FloatNegInf(),
			want: FloatInf(),
		},
		"Float -0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    Float(-0.5),
			b:    FloatNegInf(),
			want: FloatInf(),
		},
		"Float +Inf ** 5": { // Pow(+Inf, y) = +Inf for y > 0
			a:    FloatInf(),
			b:    Float(5),
			want: FloatInf(),
		},
		"Float +Inf ** -7": { // Pow(+Inf, y) = +0 for y < 0
			a:    FloatInf(),
			b:    Float(-7),
			want: Float(0),
		},
		"Float -Inf ** -7": {
			a:    FloatNegInf(),
			b:    Float(-7),
			want: Float(0),
		},
		"Float -5.5 ** 3.8": { // Pow(x, y) = NaN for finite x < 0 and finite non-integer y
			a:    Float(-5.5),
			b:    Float(3.8),
			want: FloatNaN(),
		},

		"BigFloat 5 ** 2": {
			a:    Float(5),
			b:    NewBigFloat(2),
			want: NewBigFloat(25).SetPrecision(53),
		},
		"BigFloat 7 ** 8": {
			a:    Float(7),
			b:    NewBigFloat(8),
			want: NewBigFloat(5764801).SetPrecision(53),
		},
		"BigFloat 2.5 ** 2.5": {
			a:    Float(2.5),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("9.882117688026186").SetPrecision(53),
		},
		"BigFloat 3 ** 2.5": {
			a:    Float(3),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("15.5884572681198956415").SetPrecision(53),
		},
		"BigFloat 6 ** 1": {
			a:    Float(6),
			b:    NewBigFloat(1),
			want: NewBigFloat(6).SetPrecision(53),
		},
		"BigFloat 4 ** -2": {
			a:    Float(4),
			b:    NewBigFloat(-2),
			want: NewBigFloat(0.0625).SetPrecision(53),
		},
		"BigFloat 25 ** 0": {
			a:    Float(25),
			b:    NewBigFloat(0),
			want: NewBigFloat(1).SetPrecision(53),
		},
		"BigFloat 25 ** NaN": {
			a:    Float(25),
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
		"BigFloat NaN ** 25": {
			a:    FloatNaN(),
			b:    NewBigFloat(25),
			want: BigFloatNaN(),
		},
		"BigFloat NaN ** NaN": {
			a:    FloatNaN(),
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
		"BigFloat 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    Float(0),
			b:    NewBigFloat(-5),
			want: BigFloatInf(),
		},
		"BigFloat 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    Float(0),
			b:    BigFloatNegInf(),
			want: BigFloatInf(),
		},
		"BigFloat 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    Float(0),
			b:    BigFloatInf(),
			want: NewBigFloat(0),
		},
		"BigFloat 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    Float(0),
			b:    NewBigFloat(-8),
			want: BigFloatInf(),
		},
		"BigFloat 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    Float(0),
			b:    NewBigFloat(7),
			want: NewBigFloat(0),
		},
		"BigFloat 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    Float(0),
			b:    NewBigFloat(8),
			want: NewBigFloat(0),
		},
		"BigFloat -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    Float(-1),
			b:    BigFloatInf(),
			want: NewBigFloat(1),
		},
		"BigFloat -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    Float(-1),
			b:    BigFloatNegInf(),
			want: NewBigFloat(1),
		},
		"BigFloat 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    Float(2),
			b:    BigFloatInf(),
			want: BigFloatInf(),
		},
		"BigFloat -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    Float(-2),
			b:    BigFloatInf(),
			want: BigFloatInf(),
		},
		"BigFloat 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    Float(2),
			b:    BigFloatNegInf(),
			want: NewBigFloat(0),
		},
		"BigFloat -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    Float(-2),
			b:    BigFloatNegInf(),
			want: NewBigFloat(0),
		},
		"BigFloat 0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    Float(0.5),
			b:    BigFloatInf(),
			want: NewBigFloat(0),
		},
		"BigFloat -0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    Float(-0.5),
			b:    BigFloatInf(),
			want: NewBigFloat(0),
		},
		"BigFloat 0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    Float(0.5),
			b:    BigFloatNegInf(),
			want: BigFloatInf(),
		},
		"BigFloat -0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    Float(-0.5),
			b:    BigFloatNegInf(),
			want: BigFloatInf(),
		},
		"BigFloat +Inf ** 5": { // Pow(+Inf, y) = +Inf for y > 0
			a:    FloatInf(),
			b:    NewBigFloat(5),
			want: BigFloatInf(),
		},
		"BigFloat +Inf ** -7": { // Pow(+Inf, y) = +0 for y < 0
			a:    FloatInf(),
			b:    NewBigFloat(-7),
			want: NewBigFloat(0),
		},
		"BigFloat -Inf ** -7": {
			a:    FloatNegInf(),
			b:    NewBigFloat(-7),
			want: NewBigFloat(0),
		},
		"BigFloat -5.5 ** 3.8": { // Pow(x, y) = NaN for finite x < 0 and finite non-integer y
			a:    Float(-5.5),
			b:    NewBigFloat(3.8),
			want: BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Exponentiate(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				BigFloatComparer,
				FloatComparer,
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
		"String and return an error": {
			a:   Float(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
		},

		"SmallInt 25 % 3": {
			a:    Float(25),
			b:    SmallInt(3),
			want: Float(1),
		},
		"SmallInt 25.6 % 3": {
			a:    Float(25.6),
			b:    SmallInt(3),
			want: Float(1.6000000000000014),
		},
		"SmallInt 76 % 6": {
			a:    Float(76),
			b:    SmallInt(6),
			want: Float(4),
		},
		"SmallInt -76 % 6": {
			a:    Float(-76),
			b:    SmallInt(6),
			want: Float(-4),
		},
		"SmallInt 76 % -6": {
			a:    Float(76),
			b:    SmallInt(-6),
			want: Float(4),
		},
		"SmallInt -76 % -6": {
			a:    Float(-76),
			b:    SmallInt(-6),
			want: Float(-4),
		},
		"SmallInt 124 % 9": {
			a:    Float(124),
			b:    SmallInt(9),
			want: Float(7),
		},

		"BigInt 25 % 3": {
			a:    Float(25),
			b:    NewBigInt(3),
			want: Float(1),
		},
		"BigInt 76 % 6": {
			a:    Float(76),
			b:    NewBigInt(6),
			want: Float(4),
		},
		"BigInt 76.5 % 6": {
			a:    Float(76.5),
			b:    NewBigInt(6),
			want: Float(4.5),
		},
		"BigInt -76 % 6": {
			a:    Float(-76),
			b:    NewBigInt(6),
			want: Float(-4),
		},
		"BigInt 76 % -6": {
			a:    Float(76),
			b:    NewBigInt(-6),
			want: Float(4),
		},
		"BigInt -76 % -6": {
			a:    Float(-76),
			b:    NewBigInt(-6),
			want: Float(-4),
		},
		"BigInt 124 % 9": {
			a:    Float(124),
			b:    NewBigInt(9),
			want: Float(7),
		},
		"BigInt 9765 % 9223372036854775808": {
			a:    Float(9765),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: Float(9765),
		},

		"Float 25 % 3": {
			a:    Float(25),
			b:    Float(3),
			want: Float(1),
		},
		"Float 76 % 6": {
			a:    Float(76),
			b:    Float(6),
			want: Float(4),
		},
		"Float 124 % 9": {
			a:    Float(124),
			b:    Float(9),
			want: Float(7),
		},
		"Float 124 % +Inf": {
			a:    Float(124),
			b:    FloatInf(),
			want: Float(124),
		},
		"Float 124 % -Inf": {
			a:    Float(124),
			b:    FloatNegInf(),
			want: Float(124),
		},
		"Float 74.5 % 6.25": {
			a:    Float(74.5),
			b:    Float(6.25),
			want: Float(5.75),
		},
		"Float 74 % 6.25": {
			a:    Float(74),
			b:    Float(6.25),
			want: Float(5.25),
		},
		"Float -74 % 6.25": {
			a:    Float(-74),
			b:    Float(6.25),
			want: Float(-5.25),
		},
		"Float 74 % -6.25": {
			a:    Float(74),
			b:    Float(-6.25),
			want: Float(5.25),
		},
		"Float -74 % -6.25": {
			a:    Float(-74),
			b:    Float(-6.25),
			want: Float(-5.25),
		},
		"Float +Inf % 5": { // Mod(±Inf, y) = NaN
			a:    FloatInf(),
			b:    Float(5),
			want: FloatNaN(),
		},
		"Float -Inf % 5": { // Mod(±Inf, y) = NaN
			a:    FloatNegInf(),
			b:    Float(5),
			want: FloatNaN(),
		},
		"Float NaN % 625": { // Mod(NaN, y) = NaN
			a:    FloatNaN(),
			b:    Float(625),
			want: FloatNaN(),
		},
		"Float 25 % 0": { // Mod(x, 0) = NaN
			a:    Float(25),
			b:    Float(0),
			want: FloatNaN(),
		},
		"Float 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    Float(25),
			b:    FloatInf(),
			want: Float(25),
		},
		"Float -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    Float(-87),
			b:    FloatNegInf(),
			want: Float(-87),
		},
		"Float 49 % NaN": { // Mod(x, NaN) = NaN
			a:    Float(49),
			b:    FloatNaN(),
			want: FloatNaN(),
		},

		"BigFloat 25 % 3": {
			a:    Float(25),
			b:    NewBigFloat(3),
			want: NewBigFloat(1),
		},
		"BigFloat 76 % 6": {
			a:    Float(76),
			b:    NewBigFloat(6),
			want: NewBigFloat(4),
		},
		"BigFloat 124 % 9": {
			a:    Float(124),
			b:    NewBigFloat(9),
			want: NewBigFloat(7),
		},
		"BigFloat 74 % 6.25": {
			a:    Float(74),
			b:    NewBigFloat(6.25),
			want: NewBigFloat(5.25),
		},
		"BigFloat 74 % 6.25 with higher precision": {
			a:    Float(74),
			b:    NewBigFloat(6.25).SetPrecision(64),
			want: NewBigFloat(5.25).SetPrecision(64),
		},
		"BigFloat -74 % 6.25": {
			a:    Float(-74),
			b:    NewBigFloat(6.25),
			want: NewBigFloat(-5.25),
		},
		"BigFloat 74 % -6.25": {
			a:    Float(74),
			b:    NewBigFloat(-6.25),
			want: NewBigFloat(5.25),
		},
		"BigFloat -74 % -6.25": {
			a:    Float(-74),
			b:    NewBigFloat(-6.25),
			want: NewBigFloat(-5.25),
		},
		"BigFloat +Inf % 5": { // Mod(±Inf, y) = NaN
			a:    FloatInf(),
			b:    NewBigFloat(5),
			want: BigFloatNaN(),
		},
		"BigFloat -Inf % 5": { // Mod(±Inf, y) = NaN
			a:    FloatNegInf(),
			b:    NewBigFloat(5),
			want: BigFloatNaN(),
		},
		"BigFloat NaN % 625": { // Mod(NaN, y) = NaN
			a:    FloatNaN(),
			b:    NewBigFloat(625),
			want: BigFloatNaN(),
		},
		"BigFloat 25 % 0": { // Mod(x, 0) = NaN
			a:    Float(25),
			b:    NewBigFloat(0),
			want: BigFloatNaN(),
		},
		"BigFloat 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    Float(25),
			b:    BigFloatInf(),
			want: NewBigFloat(25),
		},
		"BigFloat -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    Float(-87),
			b:    BigFloatNegInf(),
			want: NewBigFloat(-87),
		},
		"BigFloat 49 % NaN": { // Mod(x, NaN) = NaN
			a:    Float(49),
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
				BigFloatComparer,
				FloatComparer,
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

func TestFloat_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    Float
		b    Value
		want Value
		err  *Error
	}{
		"String and return an error": {
			a:   Float(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
		},
		"Char and return an error": {
			a:   Float(5),
			b:   Char('f'),
			err: NewError(TypeErrorClass, "`Std::Char` can't be coerced into `Std::Float`"),
		},
		"Int64 and return an error": {
			a:   Float(5),
			b:   Int64(7),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Float`"),
		},
		"Float64 and return an error": {
			a:   Float(5),
			b:   Float64(7),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Float`"),
		},

		"SmallInt 25.0 > 3": {
			a:    Float(25),
			b:    SmallInt(3),
			want: True,
		},
		"SmallInt 6.0 > 18": {
			a:    Float(6),
			b:    SmallInt(18),
			want: False,
		},
		"SmallInt 6.0 > 6": {
			a:    Float(6),
			b:    SmallInt(6),
			want: False,
		},
		"SmallInt 6.5 > 6": {
			a:    Float(6.5),
			b:    SmallInt(6),
			want: True,
		},

		"BigInt 25.0 > 3": {
			a:    Float(25),
			b:    NewBigInt(3),
			want: True,
		},
		"BigInt 6.0 > 18": {
			a:    Float(6),
			b:    NewBigInt(18),
			want: False,
		},
		"BigInt 6.0 > 6": {
			a:    Float(6),
			b:    NewBigInt(6),
			want: False,
		},
		"BigInt 6.5 > 6": {
			a:    Float(6.5),
			b:    NewBigInt(6),
			want: True,
		},

		"Float 25.0 > 3.0": {
			a:    Float(25),
			b:    Float(3),
			want: True,
		},
		"Float 6.0 > 18.5": {
			a:    Float(6),
			b:    Float(18.5),
			want: False,
		},
		"Float 6.0 > 6.0": {
			a:    Float(6),
			b:    Float(6),
			want: False,
		},
		"Float 6.0 > -6.0": {
			a:    Float(6),
			b:    Float(-6),
			want: True,
		},
		"Float -6.0 > 6.0": {
			a:    Float(-6),
			b:    Float(6),
			want: False,
		},
		"Float 6.5 > 6.0": {
			a:    Float(6.5),
			b:    Float(6),
			want: True,
		},
		"Float 6.0 > 6.5": {
			a:    Float(6),
			b:    Float(6.5),
			want: False,
		},
		"Float 6.0 > +Inf": {
			a:    Float(6),
			b:    FloatInf(),
			want: False,
		},
		"Float 6.0 > -Inf": {
			a:    Float(6),
			b:    FloatNegInf(),
			want: True,
		},
		"Float +Inf > +Inf": {
			a:    FloatInf(),
			b:    FloatInf(),
			want: False,
		},
		"Float +Inf > -Inf": {
			a:    FloatInf(),
			b:    FloatNegInf(),
			want: True,
		},
		"Float -Inf > +Inf": {
			a:    FloatNegInf(),
			b:    FloatInf(),
			want: False,
		},
		"Float -Inf > -Inf": {
			a:    FloatNegInf(),
			b:    FloatNegInf(),
			want: False,
		},
		"Float 6.0 > NaN": {
			a:    Float(6),
			b:    FloatNaN(),
			want: False,
		},
		"Float NaN > 6.0": {
			a:    FloatNaN(),
			b:    Float(6),
			want: False,
		},
		"Float NaN > NaN": {
			a:    FloatNaN(),
			b:    FloatNaN(),
			want: False,
		},

		"BigFloat 25.0 > 3.0bf": {
			a:    Float(25),
			b:    NewBigFloat(3),
			want: True,
		},
		"BigFloat 6.0 > 18.5bf": {
			a:    Float(6),
			b:    NewBigFloat(18.5),
			want: False,
		},
		"BigFloat 6.0 > 6.0bf": {
			a:    Float(6),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat -6.0 > 6.0bf": {
			a:    Float(-6),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat 6.0 > -6.0bf": {
			a:    Float(6),
			b:    NewBigFloat(-6),
			want: True,
		},
		"BigFloat -6.0 > -6.0bf": {
			a:    Float(-6),
			b:    NewBigFloat(-6),
			want: False,
		},
		"BigFloat 6.5 > 6.0bf": {
			a:    Float(6.5),
			b:    NewBigFloat(6),
			want: True,
		},
		"BigFloat 6.0 > 6.5bf": {
			a:    Float(6),
			b:    NewBigFloat(6.5),
			want: False,
		},
		"BigFloat 6.0 > +Inf": {
			a:    Float(6),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat 6.0 > -Inf": {
			a:    Float(6),
			b:    BigFloatNegInf(),
			want: True,
		},
		"BigFloat +Inf > 6.0": {
			a:    FloatInf(),
			b:    NewBigFloat(6),
			want: True,
		},
		"BigFloat -Inf > 6.0": {
			a:    FloatNegInf(),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat +Inf > +Inf": {
			a:    FloatInf(),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat +Inf > -Inf": {
			a:    FloatInf(),
			b:    BigFloatNegInf(),
			want: True,
		},
		"BigFloat -Inf > +Inf": {
			a:    FloatNegInf(),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat -Inf > -Inf": {
			a:    FloatNegInf(),
			b:    BigFloatNegInf(),
			want: False,
		},
		"BigFloat 6.0 > NaN": {
			a:    Float(6),
			b:    BigFloatNaN(),
			want: False,
		},
		"BigFloat NaN > 6.0bf": {
			a:    FloatNaN(),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat NaN > NaN": {
			a:    FloatNaN(),
			b:    BigFloatNaN(),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThan(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				FloatComparer,
				BigFloatComparer,
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
func TestFloat_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    Float
		b    Value
		want Value
		err  *Error
	}{
		"String and return an error": {
			a:   Float(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
		},
		"Char and return an error": {
			a:   Float(5),
			b:   Char('f'),
			err: NewError(TypeErrorClass, "`Std::Char` can't be coerced into `Std::Float`"),
		},
		"Int64 and return an error": {
			a:   Float(5),
			b:   Int64(7),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Float`"),
		},
		"Float64 and return an error": {
			a:   Float(5),
			b:   Float64(7),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Float`"),
		},

		"SmallInt 25.0 >= 3": {
			a:    Float(25),
			b:    SmallInt(3),
			want: True,
		},
		"SmallInt 6.0 >= 18": {
			a:    Float(6),
			b:    SmallInt(18),
			want: False,
		},
		"SmallInt 6.0 >= 6": {
			a:    Float(6),
			b:    SmallInt(6),
			want: True,
		},
		"SmallInt 6.5 >= 6": {
			a:    Float(6.5),
			b:    SmallInt(6),
			want: True,
		},

		"BigInt 25.0 >= 3": {
			a:    Float(25),
			b:    NewBigInt(3),
			want: True,
		},
		"BigInt 6.0 >= 18": {
			a:    Float(6),
			b:    NewBigInt(18),
			want: False,
		},
		"BigInt 6.0 >= 6": {
			a:    Float(6),
			b:    NewBigInt(6),
			want: True,
		},
		"BigInt 6.5 >= 6": {
			a:    Float(6.5),
			b:    NewBigInt(6),
			want: True,
		},

		"Float 25.0 >= 3.0": {
			a:    Float(25),
			b:    Float(3),
			want: True,
		},
		"Float 6.0 >= 18.5": {
			a:    Float(6),
			b:    Float(18.5),
			want: False,
		},
		"Float 6.0 >= 6.0": {
			a:    Float(6),
			b:    Float(6),
			want: True,
		},
		"Float 6.0 >= -6.0": {
			a:    Float(6),
			b:    Float(-6),
			want: True,
		},
		"Float -6.0 >= 6.0": {
			a:    Float(-6),
			b:    Float(6),
			want: False,
		},
		"Float 6.5 >= 6.0": {
			a:    Float(6.5),
			b:    Float(6),
			want: True,
		},
		"Float 6.0 >= 6.5": {
			a:    Float(6),
			b:    Float(6.5),
			want: False,
		},
		"Float 6.0 >= +Inf": {
			a:    Float(6),
			b:    FloatInf(),
			want: False,
		},
		"Float 6.0 >= -Inf": {
			a:    Float(6),
			b:    FloatNegInf(),
			want: True,
		},
		"Float +Inf >= +Inf": {
			a:    FloatInf(),
			b:    FloatInf(),
			want: True,
		},
		"Float +Inf >= -Inf": {
			a:    FloatInf(),
			b:    FloatNegInf(),
			want: True,
		},
		"Float -Inf >= +Inf": {
			a:    FloatNegInf(),
			b:    FloatInf(),
			want: False,
		},
		"Float -Inf >= -Inf": {
			a:    FloatNegInf(),
			b:    FloatNegInf(),
			want: True,
		},
		"Float 6.0 >= NaN": {
			a:    Float(6),
			b:    FloatNaN(),
			want: False,
		},
		"Float NaN >= 6.0": {
			a:    FloatNaN(),
			b:    Float(6),
			want: False,
		},
		"Float NaN >= NaN": {
			a:    FloatNaN(),
			b:    FloatNaN(),
			want: False,
		},

		"BigFloat 25.0 >= 3.0bf": {
			a:    Float(25),
			b:    NewBigFloat(3),
			want: True,
		},
		"BigFloat 6.0 >= 18.5bf": {
			a:    Float(6),
			b:    NewBigFloat(18.5),
			want: False,
		},
		"BigFloat 6.0 >= 6.0bf": {
			a:    Float(6),
			b:    NewBigFloat(6),
			want: True,
		},
		"BigFloat -6.0 >= 6.0bf": {
			a:    Float(-6),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat 6.0 >= -6.0bf": {
			a:    Float(6),
			b:    NewBigFloat(-6),
			want: True,
		},
		"BigFloat -6.0 >= -6.0bf": {
			a:    Float(-6),
			b:    NewBigFloat(-6),
			want: True,
		},
		"BigFloat 6.5 >= 6.0bf": {
			a:    Float(6.5),
			b:    NewBigFloat(6),
			want: True,
		},
		"BigFloat 6.0 >= 6.5bf": {
			a:    Float(6),
			b:    NewBigFloat(6.5),
			want: False,
		},
		"BigFloat 6.0 >= +Inf": {
			a:    Float(6),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat 6.0 >= -Inf": {
			a:    Float(6),
			b:    BigFloatNegInf(),
			want: True,
		},
		"BigFloat +Inf >= 6.0": {
			a:    FloatInf(),
			b:    NewBigFloat(6),
			want: True,
		},
		"BigFloat -Inf >= 6.0": {
			a:    FloatNegInf(),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat +Inf >= +Inf": {
			a:    FloatInf(),
			b:    BigFloatInf(),
			want: True,
		},
		"BigFloat +Inf >= -Inf": {
			a:    FloatInf(),
			b:    BigFloatNegInf(),
			want: True,
		},
		"BigFloat -Inf >= +Inf": {
			a:    FloatNegInf(),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat -Inf >= -Inf": {
			a:    FloatNegInf(),
			b:    BigFloatNegInf(),
			want: True,
		},
		"BigFloat 6.0 >= NaN": {
			a:    Float(6),
			b:    BigFloatNaN(),
			want: False,
		},
		"BigFloat NaN >= 6.0bf": {
			a:    FloatNaN(),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat NaN >= NaN": {
			a:    FloatNaN(),
			b:    BigFloatNaN(),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqual(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				FloatComparer,
				BigFloatComparer,
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

func TestFloat_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    Float
		b    Value
		want Value
		err  *Error
	}{
		"String and return an error": {
			a:   Float(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
		},
		"Char and return an error": {
			a:   Float(5),
			b:   Char('f'),
			err: NewError(TypeErrorClass, "`Std::Char` can't be coerced into `Std::Float`"),
		},
		"Int64 and return an error": {
			a:   Float(5),
			b:   Int64(7),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Float`"),
		},
		"Float64 and return an error": {
			a:   Float(5),
			b:   Float64(7),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Float`"),
		},

		"SmallInt 25.0 < 3": {
			a:    Float(25),
			b:    SmallInt(3),
			want: False,
		},
		"SmallInt 6.0 < 18": {
			a:    Float(6),
			b:    SmallInt(18),
			want: True,
		},
		"SmallInt 6.0 < 6": {
			a:    Float(6),
			b:    SmallInt(6),
			want: False,
		},
		"SmallInt 5.5 < 6": {
			a:    Float(5.5),
			b:    SmallInt(6),
			want: True,
		},

		"BigInt 25.0 < 3": {
			a:    Float(25),
			b:    NewBigInt(3),
			want: False,
		},
		"BigInt 6.0 < 18": {
			a:    Float(6),
			b:    NewBigInt(18),
			want: True,
		},
		"BigInt 6.0 < 6": {
			a:    Float(6),
			b:    NewBigInt(6),
			want: False,
		},
		"BigInt 5.5 < 6": {
			a:    Float(5.5),
			b:    NewBigInt(6),
			want: True,
		},

		"Float 25.0 < 3.0": {
			a:    Float(25),
			b:    Float(3),
			want: False,
		},
		"Float 6.0 < 18.5": {
			a:    Float(6),
			b:    Float(18.5),
			want: True,
		},
		"Float 6.0 < 6.0": {
			a:    Float(6),
			b:    Float(6),
			want: False,
		},
		"Float 5.5 < 6.0": {
			a:    Float(5.5),
			b:    Float(6),
			want: True,
		},
		"Float 6.0 < 6.5": {
			a:    Float(6),
			b:    Float(6.5),
			want: True,
		},
		"Float 6.3 < 6.0": {
			a:    Float(6.3),
			b:    Float(6),
			want: False,
		},
		"Float 6.0 < +Inf": {
			a:    Float(6),
			b:    FloatInf(),
			want: True,
		},
		"Float 6.0 < -Inf": {
			a:    Float(6),
			b:    FloatNegInf(),
			want: False,
		},
		"Float +Inf < 6.0": {
			a:    FloatInf(),
			b:    Float(6),
			want: False,
		},
		"Float -Inf < 6.0": {
			a:    FloatNegInf(),
			b:    Float(6),
			want: True,
		},
		"Float +Inf < +Inf": {
			a:    FloatInf(),
			b:    FloatInf(),
			want: False,
		},
		"Float -Inf < +Inf": {
			a:    FloatNegInf(),
			b:    FloatInf(),
			want: True,
		},
		"Float +Inf < -Inf": {
			a:    FloatInf(),
			b:    FloatNegInf(),
			want: False,
		},
		"Float -Inf < -Inf": {
			a:    FloatNegInf(),
			b:    FloatNegInf(),
			want: False,
		},
		"Float 6.0 < NaN": {
			a:    Float(6),
			b:    FloatNaN(),
			want: False,
		},
		"Float NaN < 6.0": {
			a:    FloatNaN(),
			b:    Float(6),
			want: False,
		},
		"Float NaN < NaN": {
			a:    FloatNaN(),
			b:    FloatNaN(),
			want: False,
		},

		"BigFloat 25.0 < 3.0bf": {
			a:    Float(25),
			b:    NewBigFloat(3),
			want: False,
		},
		"BigFloat 6.0 < 18.5bf": {
			a:    Float(6),
			b:    NewBigFloat(18.5),
			want: True,
		},
		"BigFloat 6.0 < 6bf": {
			a:    Float(6),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat 6.0 < +Inf": {
			a:    Float(6),
			b:    BigFloatInf(),
			want: True,
		},
		"BigFloat 6.0 < -Inf": {
			a:    Float(6),
			b:    BigFloatNegInf(),
			want: False,
		},
		"BigFloat +Inf < +Inf": {
			a:    FloatInf(),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat -Inf < +Inf": {
			a:    FloatNegInf(),
			b:    BigFloatInf(),
			want: True,
		},
		"BigFloat -Inf < -Inf": {
			a:    FloatNegInf(),
			b:    BigFloatNegInf(),
			want: False,
		},
		"BigFloat 6.0 < NaN": {
			a:    Float(6),
			b:    BigFloatNaN(),
			want: False,
		},
		"BigFloat NaN < 6.0bf": {
			a:    FloatNaN(),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat NaN < NaN": {
			a:    FloatNaN(),
			b:    BigFloatNaN(),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThan(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				FloatComparer,
				BigFloatComparer,
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
func TestFloat_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    Float
		b    Value
		want Value
		err  *Error
	}{
		"String and return an error": {
			a:   Float(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float`"),
		},
		"Char and return an error": {
			a:   Float(5),
			b:   Char('f'),
			err: NewError(TypeErrorClass, "`Std::Char` can't be coerced into `Std::Float`"),
		},
		"Int64 and return an error": {
			a:   Float(5),
			b:   Int64(7),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Float`"),
		},
		"Float64 and return an error": {
			a:   Float(5),
			b:   Float64(7),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Float`"),
		},

		"SmallInt 25.0 <= 3": {
			a:    Float(25),
			b:    SmallInt(3),
			want: False,
		},
		"SmallInt 6.0 <= 18": {
			a:    Float(6),
			b:    SmallInt(18),
			want: True,
		},
		"SmallInt 6.0 <= 6": {
			a:    Float(6),
			b:    SmallInt(6),
			want: True,
		},
		"SmallInt 6.5 <= 6": {
			a:    Float(6.5),
			b:    SmallInt(6),
			want: False,
		},
		"SmallInt 5.5 <= 6": {
			a:    Float(5.5),
			b:    SmallInt(6),
			want: True,
		},

		"BigInt 25.0 <= 3": {
			a:    Float(25),
			b:    NewBigInt(3),
			want: False,
		},
		"BigInt 6.0 <= 18": {
			a:    Float(6),
			b:    NewBigInt(18),
			want: True,
		},
		"BigInt 6.0 <= 6": {
			a:    Float(6),
			b:    NewBigInt(6),
			want: True,
		},
		"BigInt 6.5 <= 6": {
			a:    Float(6.5),
			b:    NewBigInt(6),
			want: False,
		},
		"BigInt 5.5 <= 6": {
			a:    Float(5.5),
			b:    NewBigInt(6),
			want: True,
		},

		"Float 25.0 <= 3.0": {
			a:    Float(25),
			b:    Float(3),
			want: False,
		},
		"Float 6.0 <= 18.5": {
			a:    Float(6),
			b:    Float(18.5),
			want: True,
		},
		"Float 6.0 <= 6.0": {
			a:    Float(6),
			b:    Float(6),
			want: True,
		},
		"Float 5.5 <= 6.0": {
			a:    Float(5.5),
			b:    Float(6),
			want: True,
		},
		"Float 6.0 <= 6.5": {
			a:    Float(6),
			b:    Float(6.5),
			want: True,
		},
		"Float 6.3 <= 6.0": {
			a:    Float(6.3),
			b:    Float(6),
			want: False,
		},
		"Float 6.0 <= +Inf": {
			a:    Float(6),
			b:    FloatInf(),
			want: True,
		},
		"Float 6.0 <= -Inf": {
			a:    Float(6),
			b:    FloatNegInf(),
			want: False,
		},
		"Float +Inf <= 6.0": {
			a:    FloatInf(),
			b:    Float(6),
			want: False,
		},
		"Float -Inf <= 6.0": {
			a:    FloatNegInf(),
			b:    Float(6),
			want: True,
		},
		"Float +Inf <= +Inf": {
			a:    FloatInf(),
			b:    FloatInf(),
			want: True,
		},
		"Float -Inf <= +Inf": {
			a:    FloatNegInf(),
			b:    FloatInf(),
			want: True,
		},
		"Float +Inf <= -Inf": {
			a:    FloatInf(),
			b:    FloatNegInf(),
			want: False,
		},
		"Float -Inf <= -Inf": {
			a:    FloatNegInf(),
			b:    FloatNegInf(),
			want: True,
		},
		"Float 6.0 <= NaN": {
			a:    Float(6),
			b:    FloatNaN(),
			want: False,
		},
		"Float NaN <= 6.0": {
			a:    FloatNaN(),
			b:    Float(6),
			want: False,
		},
		"Float NaN <= NaN": {
			a:    FloatNaN(),
			b:    FloatNaN(),
			want: False,
		},

		"BigFloat 25.0 <= 3.0bf": {
			a:    Float(25),
			b:    NewBigFloat(3),
			want: False,
		},
		"BigFloat 6.0 <= 18.5bf": {
			a:    Float(6),
			b:    NewBigFloat(18.5),
			want: True,
		},
		"BigFloat 6.0 <= 6bf": {
			a:    Float(6),
			b:    NewBigFloat(6),
			want: True,
		},
		"BigFloat 6.0 <= +Inf": {
			a:    Float(6),
			b:    BigFloatInf(),
			want: True,
		},
		"BigFloat 6.0 <= -Inf": {
			a:    Float(6),
			b:    BigFloatNegInf(),
			want: False,
		},
		"BigFloat +Inf <= +Inf": {
			a:    FloatInf(),
			b:    BigFloatInf(),
			want: True,
		},
		"BigFloat -Inf <= +Inf": {
			a:    FloatNegInf(),
			b:    BigFloatInf(),
			want: True,
		},
		"BigFloat -Inf <= -Inf": {
			a:    FloatNegInf(),
			b:    BigFloatNegInf(),
			want: True,
		},
		"BigFloat 6.0 <= NaN": {
			a:    Float(6),
			b:    BigFloatNaN(),
			want: False,
		},
		"BigFloat NaN <= 6.0bf": {
			a:    FloatNaN(),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat NaN <= NaN": {
			a:    FloatNaN(),
			b:    BigFloatNaN(),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqual(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				FloatComparer,
				BigFloatComparer,
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

func TestFloat_Equal(t *testing.T) {
	tests := map[string]struct {
		a    Float
		b    Value
		want Value
	}{
		"String 5.0 == '5'": {
			a:    Float(5),
			b:    String("5"),
			want: False,
		},
		"Char 5.0 == c'5'": {
			a:    Float(5),
			b:    Char('5'),
			want: False,
		},

		"Int64 5.0 == 5i64": {
			a:    Float(5),
			b:    Int64(5),
			want: True,
		},
		"Int64 5.5 == 5i64": {
			a:    Float(5.5),
			b:    Int64(5),
			want: False,
		},
		"Int64 4.0 == 5i64": {
			a:    Float(4),
			b:    Int64(5),
			want: False,
		},

		"Int32 5.0 == 5i32": {
			a:    Float(5),
			b:    Int32(5),
			want: True,
		},
		"Int32 5.5 == 5i32": {
			a:    Float(5.5),
			b:    Int32(5),
			want: False,
		},
		"Int32 4.0 == 5i32": {
			a:    Float(4),
			b:    Int32(5),
			want: False,
		},

		"Int16 5.0 == 5i16": {
			a:    Float(5),
			b:    Int16(5),
			want: True,
		},
		"Int16 5.5 == 5i16": {
			a:    Float(5.5),
			b:    Int16(5),
			want: False,
		},
		"Int16 4.0 == 5i16": {
			a:    Float(4),
			b:    Int16(5),
			want: False,
		},

		"Int8 5.0 == 5i8": {
			a:    Float(5),
			b:    Int8(5),
			want: True,
		},
		"Int8 5.5 == 5i8": {
			a:    Float(5.5),
			b:    Int8(5),
			want: False,
		},
		"Int8 4.0 == 5i8": {
			a:    Float(4),
			b:    Int8(5),
			want: False,
		},

		"UInt64 5.0 == 5u64": {
			a:    Float(5),
			b:    UInt64(5),
			want: True,
		},
		"UInt64 5.5 == 5u64": {
			a:    Float(5.5),
			b:    UInt64(5),
			want: False,
		},
		"UInt64 4.0 == 5u64": {
			a:    Float(4),
			b:    UInt64(5),
			want: False,
		},

		"UInt32 5.0 == 5u32": {
			a:    Float(5),
			b:    UInt32(5),
			want: True,
		},
		"UInt32 5.5 == 5u32": {
			a:    Float(5.5),
			b:    UInt32(5),
			want: False,
		},
		"UInt32 4.0 == 5u32": {
			a:    Float(4),
			b:    UInt32(5),
			want: False,
		},

		"UInt16 5.0 == 5u16": {
			a:    Float(5),
			b:    UInt16(5),
			want: True,
		},
		"UInt16 5.5 == 5u16": {
			a:    Float(5.5),
			b:    UInt16(5),
			want: False,
		},
		"UInt16 4.0 == 5u16": {
			a:    Float(4),
			b:    UInt16(5),
			want: False,
		},

		"UInt8 5.0 == 5u8": {
			a:    Float(5),
			b:    UInt8(5),
			want: True,
		},
		"UInt8 5.5 == 5u8": {
			a:    Float(5.5),
			b:    UInt8(5),
			want: False,
		},
		"UInt8 4.0 == 5u8": {
			a:    Float(4),
			b:    UInt8(5),
			want: False,
		},

		"Float64 5.0 == 5f64": {
			a:    Float(5),
			b:    Float64(5),
			want: True,
		},
		"Float64 5.5 == 5f64": {
			a:    Float(5.5),
			b:    Float64(5),
			want: False,
		},
		"Float64 5.0 == 5.5f64": {
			a:    Float(5),
			b:    Float64(5.5),
			want: False,
		},
		"Float64 5.5 == 5.5f64": {
			a:    Float(5.5),
			b:    Float64(5.5),
			want: True,
		},
		"Float64 4.0 == 5f64": {
			a:    Float(4),
			b:    Float64(5),
			want: False,
		},

		"Float32 5.0 == 5f32": {
			a:    Float(5),
			b:    Float32(5),
			want: True,
		},
		"Float32 5.5 == 5f32": {
			a:    Float(5.5),
			b:    Float32(5),
			want: False,
		},
		"Float32 5.0 == 5.5f32": {
			a:    Float(5),
			b:    Float32(5.5),
			want: False,
		},
		"Float32 5.5 == 5.5f32": {
			a:    Float(5.5),
			b:    Float32(5.5),
			want: True,
		},
		"Float32 4.0 == 5f32": {
			a:    Float(4),
			b:    Float32(5),
			want: False,
		},

		"SmallInt 25.0 == 3": {
			a:    Float(25),
			b:    SmallInt(3),
			want: False,
		},
		"SmallInt 6.0 == 18": {
			a:    Float(6),
			b:    SmallInt(18),
			want: False,
		},
		"SmallInt 6.0 == 6": {
			a:    Float(6),
			b:    SmallInt(6),
			want: True,
		},
		"SmallInt 6.5 == 6": {
			a:    Float(6.5),
			b:    SmallInt(6),
			want: False,
		},

		"BigInt 25.0 == 3": {
			a:    Float(25),
			b:    NewBigInt(3),
			want: False,
		},
		"BigInt 6.0 == 18": {
			a:    Float(6),
			b:    NewBigInt(18),
			want: False,
		},
		"BigInt 6.0 == 6": {
			a:    Float(6),
			b:    NewBigInt(6),
			want: True,
		},
		"BigInt 6.5 == 6": {
			a:    Float(6.5),
			b:    NewBigInt(6),
			want: False,
		},

		"Float 25.0 == 3.0": {
			a:    Float(25),
			b:    Float(3),
			want: False,
		},
		"Float 6.0 == 18.5": {
			a:    Float(6),
			b:    Float(18.5),
			want: False,
		},
		"Float 6.0 == 6": {
			a:    Float(6),
			b:    Float(6),
			want: True,
		},
		"Float 6.0 == +Inf": {
			a:    Float(6),
			b:    FloatInf(),
			want: False,
		},
		"Float 6.0 == -Inf": {
			a:    Float(6),
			b:    FloatNegInf(),
			want: False,
		},
		"Float +Inf == 6.0": {
			a:    FloatInf(),
			b:    Float(6),
			want: False,
		},
		"Float -Inf == 6.0": {
			a:    FloatNegInf(),
			b:    Float(6),
			want: False,
		},
		"Float +Inf == +Inf": {
			a:    FloatInf(),
			b:    FloatInf(),
			want: True,
		},
		"Float +Inf == -Inf": {
			a:    FloatInf(),
			b:    FloatNegInf(),
			want: False,
		},
		"Float -Inf == +Inf": {
			a:    FloatNegInf(),
			b:    FloatInf(),
			want: False,
		},
		"Float 6.0 == NaN": {
			a:    Float(6),
			b:    FloatNaN(),
			want: False,
		},
		"Float NaN == 6.0": {
			a:    FloatNaN(),
			b:    Float(6),
			want: False,
		},
		"Float NaN == NaN": {
			a:    FloatNaN(),
			b:    FloatNaN(),
			want: False,
		},

		"BigFloat 25.0 == 3.0bf": {
			a:    Float(25),
			b:    NewBigFloat(3),
			want: False,
		},
		"BigFloat 6.0 == 18.5bf": {
			a:    Float(6),
			b:    NewBigFloat(18.5),
			want: False,
		},
		"BigFloat 6.0 == 6bf": {
			a:    Float(6),
			b:    NewBigFloat(6),
			want: True,
		},
		"BigFloat 6.0 == +Inf": {
			a:    Float(6),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat 6.0 == -Inf": {
			a:    Float(6),
			b:    BigFloatNegInf(),
			want: False,
		},
		"BigFloat +Inf == 6.0bf": {
			a:    FloatInf(),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat -Inf == 6.0bf": {
			a:    FloatNegInf(),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat +Inf == +Inf": {
			a:    FloatInf(),
			b:    BigFloatInf(),
			want: True,
		},
		"BigFloat +Inf == -Inf": {
			a:    FloatInf(),
			b:    BigFloatNegInf(),
			want: False,
		},
		"BigFloat -Inf == +Inf": {
			a:    FloatNegInf(),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat 6.0 == NaN": {
			a:    Float(6),
			b:    BigFloatNaN(),
			want: False,
		},
		"BigFloat NaN == 6.0bf": {
			a:    FloatNaN(),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat NaN == NaN": {
			a:    FloatNaN(),
			b:    BigFloatNaN(),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				FloatComparer,
				BigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat_StrictEqual(t *testing.T) {
	tests := map[string]struct {
		a    Float
		b    Value
		want Value
	}{
		"String 5.0 === '5'": {
			a:    Float(5),
			b:    String("5"),
			want: False,
		},
		"Char 5.0 === c'5'": {
			a:    Float(5),
			b:    Char('5'),
			want: False,
		},

		"Int64 5.0 === 5i64": {
			a:    Float(5),
			b:    Int64(5),
			want: False,
		},
		"Int64 5.3 === 5i64": {
			a:    Float(5.3),
			b:    Int64(5),
			want: False,
		},
		"Int64 4.0 === 5i64": {
			a:    Float(4),
			b:    Int64(5),
			want: False,
		},

		"Int32 5.0 === 5i32": {
			a:    Float(5),
			b:    Int32(5),
			want: False,
		},
		"Int32 5.2 === 5i32": {
			a:    Float(5.2),
			b:    Int32(5),
			want: False,
		},
		"Int32 4.0 === 5i32": {
			a:    Float(4),
			b:    Int32(5),
			want: False,
		},

		"Int16 5.0 === 5i16": {
			a:    Float(5),
			b:    Int16(5),
			want: False,
		},
		"Int16 5.8 === 5i16": {
			a:    Float(5.8),
			b:    Int16(5),
			want: False,
		},
		"Int16 4.0 === 5i16": {
			a:    Float(4),
			b:    Int16(5),
			want: False,
		},

		"Int8 5.0 === 5i8": {
			a:    Float(5),
			b:    Int8(5),
			want: False,
		},
		"Int8 4.0 === 5i8": {
			a:    Float(4),
			b:    Int8(5),
			want: False,
		},

		"UInt64 5.0 === 5u64": {
			a:    Float(5),
			b:    UInt64(5),
			want: False,
		},
		"UInt64 5.7 === 5u64": {
			a:    Float(5.7),
			b:    UInt64(5),
			want: False,
		},
		"UInt64 4.0 === 5u64": {
			a:    Float(4),
			b:    UInt64(5),
			want: False,
		},

		"UInt32 5.0 === 5u32": {
			a:    Float(5),
			b:    UInt32(5),
			want: False,
		},
		"UInt32 5.3 === 5u32": {
			a:    Float(5.3),
			b:    UInt32(5),
			want: False,
		},
		"UInt32 4.0 === 5u32": {
			a:    Float(4),
			b:    UInt32(5),
			want: False,
		},

		"UInt16 5.0 === 5u16": {
			a:    Float(5),
			b:    UInt16(5),
			want: False,
		},
		"UInt16 5.65 === 5u16": {
			a:    Float(5.65),
			b:    UInt16(5),
			want: False,
		},
		"UInt16 4.0 === 5u16": {
			a:    Float(4),
			b:    UInt16(5),
			want: False,
		},

		"UInt8 5.0 === 5u8": {
			a:    Float(5),
			b:    UInt8(5),
			want: False,
		},
		"UInt8 5.12 === 5u8": {
			a:    Float(5.12),
			b:    UInt8(5),
			want: False,
		},
		"UInt8 4.0 === 5u8": {
			a:    Float(4),
			b:    UInt8(5),
			want: False,
		},

		"Float64 5.0 === 5f64": {
			a:    Float(5),
			b:    Float64(5),
			want: False,
		},
		"Float64 5.0 === 5.5f64": {
			a:    Float(5),
			b:    Float64(5.5),
			want: False,
		},
		"Float64 5.5 === 5.5f64": {
			a:    Float(5),
			b:    Float64(5.5),
			want: False,
		},
		"Float64 4.0 === 5f64": {
			a:    Float(4),
			b:    Float64(5),
			want: False,
		},

		"Float32 5.0 === 5f32": {
			a:    Float(5),
			b:    Float32(5),
			want: False,
		},
		"Float32 5.0 === 5.5f32": {
			a:    Float(5),
			b:    Float32(5.5),
			want: False,
		},
		"Float32 5.5 === 5.5f32": {
			a:    Float(5.5),
			b:    Float32(5.5),
			want: False,
		},
		"Float32 4.0 === 5f32": {
			a:    Float(4),
			b:    Float32(5),
			want: False,
		},

		"SmallInt 25.0 === 3": {
			a:    Float(25),
			b:    SmallInt(3),
			want: False,
		},
		"SmallInt 6.0 === 18": {
			a:    Float(6),
			b:    SmallInt(18),
			want: False,
		},
		"SmallInt 6.0 === 6": {
			a:    Float(6),
			b:    SmallInt(6),
			want: False,
		},
		"SmallInt 6.5 === 6": {
			a:    Float(6.5),
			b:    SmallInt(6),
			want: False,
		},

		"BigInt 25.0 === 3": {
			a:    Float(25),
			b:    NewBigInt(3),
			want: False,
		},
		"BigInt 6.0 === 18": {
			a:    Float(6),
			b:    NewBigInt(18),
			want: False,
		},
		"BigInt 6.0 === 6": {
			a:    Float(6),
			b:    NewBigInt(6),
			want: False,
		},
		"BigInt 6.5 === 6": {
			a:    Float(6.5),
			b:    NewBigInt(6),
			want: False,
		},

		"Float 25.0 === 3.0": {
			a:    Float(25),
			b:    Float(3),
			want: False,
		},
		"Float 6.0 === 18.5": {
			a:    Float(6),
			b:    Float(18.5),
			want: False,
		},
		"Float 6.0 === 6.0": {
			a:    Float(6),
			b:    Float(6),
			want: True,
		},
		"Float 27.5 === 27.5": {
			a:    Float(27.5),
			b:    Float(27.5),
			want: True,
		},
		"Float 6.5 === 6.0": {
			a:    Float(6.5),
			b:    Float(6),
			want: False,
		},
		"Float 6.0 === Inf": {
			a:    Float(6),
			b:    FloatInf(),
			want: False,
		},
		"Float 6.0 === -Inf": {
			a:    Float(6),
			b:    FloatNegInf(),
			want: False,
		},
		"Float 6.0 === NaN": {
			a:    Float(6),
			b:    FloatNaN(),
			want: False,
		},

		"BigFloat 25.0 === 3bf": {
			a:    Float(25),
			b:    NewBigFloat(3),
			want: False,
		},
		"BigFloat 6.0 === 18.5bf": {
			a:    Float(6),
			b:    NewBigFloat(18.5),
			want: False,
		},
		"BigFloat 6.0 === 6bf": {
			a:    Float(6),
			b:    NewBigFloat(6),
			want: False,
		},
		"BigFloat 6.5 === 6.5bf": {
			a:    Float(6.5),
			b:    NewBigFloat(6.5),
			want: False,
		},
		"BigFloat 6.0 === Inf": {
			a:    Float(6),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat 6.0 === -Inf": {
			a:    Float(6),
			b:    BigFloatNegInf(),
			want: False,
		},
		"BigFloat 6.0 === NaN": {
			a:    Float(6),
			b:    BigFloatNaN(),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.StrictEqual(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				FloatComparer,
				BigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}
