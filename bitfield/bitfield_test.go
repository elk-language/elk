package bitfield

import (
	"testing"
)

func TestBitfield8_HasFlag(t *testing.T) {
	tests := map[string]struct {
		bitfield BitField8
		flag     BitFlag8
		want     bool
	}{
		"return false when empty map": {
			flag: 0b10000000,
			want: false,
		},
		"return true when the bit is set": {
			bitfield: BitField8FromInt(0b11001011),
			flag:     0b10000000,
			want:     true,
		},
		"return false when the bit is not set": {
			bitfield: BitField8FromInt(0b11001011),
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

func TestBitfield8_SetFlag(t *testing.T) {
	tests := map[string]struct {
		in   BitField8
		flag BitFlag8
		want BitField8
	}{
		"set bit when map empty": {
			flag: 0b10000000,
			want: BitField8FromInt(0b10000000),
		},
		"do nothing when bit already set": {
			in:   BitField8FromInt(0b11001011),
			flag: 0b10000000,
			want: BitField8FromInt(0b11001011),
		},
		"set bit when unset": {
			in:   BitField8FromInt(0b11001011),
			flag: 0b00100000,
			want: BitField8FromInt(0b11101011),
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

func TestBitfield8_UnsetFlag(t *testing.T) {
	tests := map[string]struct {
		in   BitField8
		flag BitFlag8
		want BitField8
	}{
		"unset bit when map empty": {
			flag: 0b10000000,
			want: BitField8FromInt(0b00000000),
		},
		"unset existing flag": {
			in:   BitField8FromInt(0b11001011),
			flag: 0b10000000,
			want: BitField8FromInt(0b01001011),
		},
		"do nothing when already unset": {
			in:   BitField8FromInt(0b11001011),
			flag: 0b00100000,
			want: BitField8FromInt(0b11001011),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.UnsetFlag(tc.flag)
			got := tc.in
			if tc.want.bitfield != got.bitfield {
				t.Fatalf("wanted: %b, got: %b", tc.want.bitfield, got.bitfield)
			}
		})
	}
}
