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
				bigFloatComparer,
				floatComparer,
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
