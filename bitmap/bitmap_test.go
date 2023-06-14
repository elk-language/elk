package bitmap

import (
	"testing"
)

func TestBitmap8HasFlag(t *testing.T) {
	tests := map[string]struct {
		bitmap Bitmap8
		flag   BitFlag8
		want   bool
	}{
		"return false when empty map": {
			flag: 0b10000000,
			want: false,
		},
		"return true when the bit is set": {
			bitmap: Bitmap8FromInt(0b11001011),
			flag:   0b10000000,
			want:   true,
		},
		"return false when the bit is not set": {
			bitmap: Bitmap8FromInt(0b11001011),
			flag:   0b00100000,
			want:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.bitmap.HasFlag(tc.flag)
			if tc.want != got {
				t.Fatalf("wanted: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestBitmap8SetFlag(t *testing.T) {
	tests := map[string]struct {
		in   Bitmap8
		flag BitFlag8
		want Bitmap8
	}{
		"set bit when map empty": {
			flag: 0b10000000,
			want: Bitmap8FromInt(0b10000000),
		},
		"do nothing when bit already set": {
			in:   Bitmap8FromInt(0b11001011),
			flag: 0b10000000,
			want: Bitmap8FromInt(0b11001011),
		},
		"set bit when unset": {
			in:   Bitmap8FromInt(0b11001011),
			flag: 0b00100000,
			want: Bitmap8FromInt(0b11101011),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.SetFlag(tc.flag)
			got := tc.in
			if tc.want.bitmap != got.bitmap {
				t.Fatalf("wanted: %b, got: %b", tc.want.bitmap, got.bitmap)
			}
		})
	}
}
