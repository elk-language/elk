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
