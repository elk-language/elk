package object

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
