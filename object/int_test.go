package object

import (
	"testing"

	"math"
	"math/big"

	"github.com/google/go-cmp/cmp"
)

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
