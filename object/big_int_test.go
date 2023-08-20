package object

import (
	"testing"

	"math"
	"math/big"

	"github.com/google/go-cmp/cmp"
)

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
