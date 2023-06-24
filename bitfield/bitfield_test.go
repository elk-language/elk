package bitfield

import (
	"testing"
)

func TestBitfield8HasFlag(t *testing.T) {
	tests := map[string]struct {
		bitfield Bitfield8
		flag     BitFlag8
		want     bool
	}{
		"return false when empty map": {
			flag: 0b10000000,
			want: false,
		},
		"return true when the bit is set": {
			bitfield: Bitfield8FromInt(0b11001011),
			flag:     0b10000000,
			want:     true,
		},
		"return false when the bit is not set": {
			bitfield: Bitfield8FromInt(0b11001011),
			flag:     0b00100000,
			want:     false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.bitfield.HasFlag(tc.flag)
			if tc.want != got {
				t.Fatalf("wanted: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestBitfield8SetFlag(t *testing.T) {
	tests := map[string]struct {
		in   Bitfield8
		flag BitFlag8
		want Bitfield8
	}{
		"set bit when map empty": {
			flag: 0b10000000,
			want: Bitfield8FromInt(0b10000000),
		},
		"do nothing when bit already set": {
			in:   Bitfield8FromInt(0b11001011),
			flag: 0b10000000,
			want: Bitfield8FromInt(0b11001011),
		},
		"set bit when unset": {
			in:   Bitfield8FromInt(0b11001011),
			flag: 0b00100000,
			want: Bitfield8FromInt(0b11101011),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.SetFlag(tc.flag)
			got := tc.in
			if tc.want.bitfield != got.bitfield {
				t.Fatalf("wanted: %b, got: %b", tc.want.bitfield, got.bitfield)
			}
		})
	}
}
