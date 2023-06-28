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
