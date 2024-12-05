//go:build amd64 || amd64p32 || arm64 || arm64be || ppc64 || ppc64le || mips64 || mips64le || mips64p32 || mips64p32le || s390x || sparc64

package value_test

import (
	"math"
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestSmallInt_Add_64sys(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"add SmallInt overflow and return BigInt": {
			a:    value.SmallInt(math.MaxInt64),
			b:    value.SmallInt(10).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775817", 10)),
		},
		"add BigInt and return BigInt": {
			a:    value.SmallInt(20),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775817", 10)),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775837", 10)),
		},
		"add BigInt and return SmallInt": {
			a:    value.SmallInt(-20),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775817", 10)),
			want: value.SmallInt(9223372036854775797).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Add(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSmallInt_Subtract_64sys(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"subtract SmallInt underflow and return BigInt": {
			a:    value.SmallInt(math.MinInt64),
			b:    value.SmallInt(10).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("-9223372036854775818", 10)),
		},
		"subtract BigInt and return BigInt": {
			a:    value.SmallInt(5),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775817", 10)),
			want: value.Ref(value.ParseBigIntPanic("-9223372036854775812", 10)),
		},
		"subtract BigInt and return SmallInt": {
			a:    value.SmallInt(20),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775817", 10)),
			want: value.SmallInt(-9223372036854775797).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Subtract(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSmallInt_Multiply_64sys(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"multiply by SmallInt overflow and return BigInt": {
			a:    value.SmallInt(math.MaxInt64),
			b:    value.SmallInt(10).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("92233720368547758070", 10)),
		},
		"multiply by BigInt and return BigInt": {
			a:    value.SmallInt(20),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775817", 10)),
			want: value.Ref(value.ParseBigIntPanic("184467440737095516340", 10)),
		},
		"multiply BigInt and return SmallInt": {
			a:    value.SmallInt(-1),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
			want: value.SmallInt(math.MinInt64).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Multiply(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSmallInt_Divide_64sys(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"divide by SmallInt overflow and return BigInt": {
			a:    value.SmallInt(math.MinInt64),
			b:    value.SmallInt(-1).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
		},
		"divide by BigInt and return SmallInt": {
			a:    value.SmallInt(20),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775817", 10)),
			want: value.SmallInt(0).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Divide(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSmallInt_RightBitshift_64sys(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"shift by SmallInt 80 >> -9223372036854775808": {
			a:    value.SmallInt(80),
			b:    value.SmallInt(-9223372036854775808).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by SmallInt overflow": {
			a:    value.SmallInt(10),
			b:    value.SmallInt(-60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by SmallInt close to overflow": {
			a:    value.SmallInt(10),
			b:    value.SmallInt(-59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by BigInt overflow": {
			a:    value.SmallInt(10),
			b:    value.Ref(value.NewBigInt(-60)),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by BigInt close to overflow": {
			a:    value.SmallInt(10),
			b:    value.Ref(value.NewBigInt(-59)),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by huge BigInt": {
			a:    value.SmallInt(10),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by Int64 overflow": {
			a:    value.SmallInt(10),
			b:    value.Int64(-60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by Int64 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.Int64(-59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by Int32 overflow": {
			a:    value.SmallInt(10),
			b:    value.Int32(-60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by Int32 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.Int32(-59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by Int16 overflow": {
			a:    value.SmallInt(10),
			b:    value.Int16(-60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by Int16 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.Int16(-59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by Int8 overflow": {
			a:    value.SmallInt(10),
			b:    value.Int8(-60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by Int8 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.Int64(-59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.RightBitshift(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSmallInt_LeftBitshift_64sys(t *testing.T) {
	tests := map[string]struct {
		a    value.SmallInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"shift by SmallInt 80 << 56": {
			a:    value.SmallInt(80),
			b:    value.SmallInt(56).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by SmallInt 80 >> -9223372036854775808": {
			a:    value.SmallInt(80),
			b:    value.SmallInt(-9223372036854775808).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by SmallInt overflow": {
			a:    value.SmallInt(10),
			b:    value.SmallInt(60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by SmallInt close to overflow": {
			a:    value.SmallInt(10),
			b:    value.SmallInt(59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by BigInt 80 << 56": {
			a:    value.SmallInt(80),
			b:    value.Ref(value.NewBigInt(56)),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by BigInt 80 >> -9223372036854775808": {
			a:    value.SmallInt(80),
			b:    value.Ref(value.NewBigInt(-9223372036854775808)),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by BigInt overflow": {
			a:    value.SmallInt(10),
			b:    value.Ref(value.NewBigInt(60)),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by BigInt close to overflow": {
			a:    value.SmallInt(10),
			b:    value.Ref(value.NewBigInt(59)),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by huge BigInt": {
			a:    value.SmallInt(10),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by huge negative BigInt": {
			a:    value.SmallInt(10),
			b:    value.Ref(value.ParseBigIntPanic("-9223372036854775809", 10)),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by Int64 80 << 56": {
			a:    value.SmallInt(80),
			b:    value.Int64(56).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by Int64 80 >> -9223372036854775808": {
			a:    value.SmallInt(80),
			b:    value.Int64(-9223372036854775808).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by Int64 overflow": {
			a:    value.SmallInt(10),
			b:    value.Int64(60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by Int64 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.Int64(59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},

		"shift by Int32 80 << 56": {
			a:    value.SmallInt(80),
			b:    value.Int32(56).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by Int32 overflow": {
			a:    value.SmallInt(10),
			b:    value.Int32(60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by Int32 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.Int32(59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},

		"shift by Int16 80 << 56": {
			a:    value.SmallInt(80),
			b:    value.Int16(56).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by Int16 overflow": {
			a:    value.SmallInt(10),
			b:    value.Int16(60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by Int16 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.Int16(59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by Int8 80 << 56": {
			a:    value.SmallInt(80),
			b:    value.Int8(56).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by Int8 overflow": {
			a:    value.SmallInt(10),
			b:    value.Int8(60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by Int8 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.Int8(59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},

		"shift by UInt64 80 << 56": {
			a:    value.SmallInt(80),
			b:    value.UInt64(56).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by UIn64 overflow": {
			a:    value.SmallInt(10),
			b:    value.UInt64(60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by UInt64 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.UInt64(59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by UInt32 80 << 56": {
			a:    value.SmallInt(80),
			b:    value.UInt32(56).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by UIn32 overflow": {
			a:    value.SmallInt(10),
			b:    value.UInt32(60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by UInt32 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.UInt32(59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by UInt16 80 << 56": {
			a:    value.SmallInt(80),
			b:    value.UInt16(56).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by UIn16 overflow": {
			a:    value.SmallInt(10),
			b:    value.UInt16(60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by UInt16 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.UInt16(59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by UInt8 80 << 56": {
			a:    value.SmallInt(80),
			b:    value.UInt8(56).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
		"shift by UIn8 overflow": {
			a:    value.SmallInt(10),
			b:    value.UInt8(60).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("11529215046068469760", 10)),
		},
		"shift by UInt8 close to overflow": {
			a:    value.SmallInt(10),
			b:    value.UInt8(59).ToValue(),
			want: value.SmallInt(5764607523034234880).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LeftBitshift(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
