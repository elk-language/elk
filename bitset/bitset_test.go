package bitset

import (
	"testing"
)

func TestBitset8HasFlag(t *testing.T) {
	tests := map[string]struct {
		bitset Bitset8
		flag   BitFlag8
		want   bool
	}{
		"return false when empty map": {
			flag: 0b10000000,
			want: false,
		},
		"return true when the bit is set": {
			bitset: Bitset8FromInt(0b11001011),
			flag:   0b10000000,
			want:   true,
		},
		"return false when the bit is not set": {
			bitset: Bitset8FromInt(0b11001011),
			flag:   0b00100000,
			want:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.bitset.HasFlag(tc.flag)
			if tc.want != got {
				t.Fatalf("wanted: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestBitset8SetFlag(t *testing.T) {
	tests := map[string]struct {
		in   Bitset8
		flag BitFlag8
		want Bitset8
	}{
		"set bit when map empty": {
			flag: 0b10000000,
			want: Bitset8FromInt(0b10000000),
		},
		"do nothing when bit already set": {
			in:   Bitset8FromInt(0b11001011),
			flag: 0b10000000,
			want: Bitset8FromInt(0b11001011),
		},
		"set bit when unset": {
			in:   Bitset8FromInt(0b11001011),
			flag: 0b00100000,
			want: Bitset8FromInt(0b11101011),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.SetFlag(tc.flag)
			got := tc.in
			if tc.want.bitset != got.bitset {
				t.Fatalf("wanted: %b, got: %b", tc.want.bitset, got.bitset)
			}
		})
	}
}
