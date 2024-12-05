//go:build amd64 || amd64p32 || arm64 || arm64be || ppc64 || ppc64le || mips64 || mips64le || mips64p32 || mips64p32le || s390x || sparc64

package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestBigInt_Add_64sys(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"add SmallInt and return BigInt": {
			a:    value.ParseBigIntPanic("9223372036854775815", 10),
			b:    value.SmallInt(10).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775825", 10)),
		},
		"add SmallInt and return SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775837", 10),
			b:    value.SmallInt(-10).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775827", 10)),
		},
		"add BigInt and return BigInt": {
			a:    value.ParseBigIntPanic("9223372036854775827", 10),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775830", 10)),
		},
		"add BigInt and return SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775827", 10),
			b:    value.Ref(value.NewBigInt(-27)),
			want: value.SmallInt(9223372036854775800).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Add(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigInt_Subtract_64sys(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"subtract SmallInt and return BigInt": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.SmallInt(5).ToValue(),
			want: value.Ref(value.ParseBigIntPanic("9223372036854775812", 10)),
		},
		"subtract SmallInt and return SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775817", 10),
			b:    value.SmallInt(11).ToValue(),
			want: value.SmallInt(9223372036854775806).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Subtract(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigInt_Multiply_64sys(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"multiply by SmallInt and return SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775808", 10),
			b:    value.SmallInt(-1).ToValue(),
			want: value.SmallInt(-9223372036854775808).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Multiply(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigInt_Divide_64sys(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"divide by SmallInt and return SmallInt": {
			a:    value.ParseBigIntPanic("9223372036854775818", 10),
			b:    value.SmallInt(2).ToValue(),
			want: value.SmallInt(4611686018427387909).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Divide(tc.b)
			opts := comparer.Options()
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

func TestBigInt_RightBitshift_64sys(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"shift by SmallInt 80 >> -9223372036854775808": {
			a:    value.NewBigInt(80),
			b:    value.SmallInt(-9223372036854775808).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by BigInt 80 >> -9223372036854775808": {
			a:    value.NewBigInt(80),
			b:    value.Ref(value.NewBigInt(-9223372036854775808)),
			want: value.SmallInt(0).ToValue(),
		},
		"shift by Int64 80 >> -9223372036854775808": {
			a:    value.NewBigInt(80),
			b:    value.Int64(-9223372036854775808).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.RightBitshift(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigInt_BitwiseXor_64sys(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigInt
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"9223372036857247042 ^ 10223372099998981329": {
			a:    value.ParseBigIntPanic("9223372036857247042", 10),
			b:    value.Ref(value.ParseBigIntPanic("10223372099998981329", 10)),
			want: value.SmallInt(1000000063146142099).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.BitwiseXor(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
