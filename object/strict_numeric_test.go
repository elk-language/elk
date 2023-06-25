package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStrictNumericParseUint(t *testing.T) {
	tests := map[string]struct {
		str     string
		base    int
		bitSize int
		want    uint64
		err     *Error
	}{
		"explicit decimal": {
			str:     "123",
			base:    10,
			bitSize: 8,
			want:    123,
		},
		"implicit decimal": {
			str:     "123",
			base:    0,
			bitSize: 8,
			want:    123,
		},
		"implicit decimal with underscores": {
			str:     "2_500",
			base:    0,
			bitSize: 16,
			want:    2500,
		},
		"8bit implicit decimal out of range": {
			str:     "300",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"explicit hex": {
			str:     "ff",
			base:    16,
			bitSize: 8,
			want:    255,
		},
		"implicit hex": {
			str:     "0xff",
			base:    0,
			bitSize: 8,
			want:    255,
		},
		"implicit hex with underscores": {
			str:     "0x12_34",
			base:    0,
			bitSize: 16,
			want:    4660,
		},
		"8bit implicit hex out of range": {
			str:     "0xfff",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"explicit duodecimal": {
			str:     "1a",
			base:    12,
			bitSize: 8,
			want:    22,
		},
		"implicit duodecimal": {
			str:     "0d1a",
			base:    0,
			bitSize: 8,
			want:    22,
		},
		"implicit duodecimal with underscores": {
			str:     "0d12_34",
			base:    0,
			bitSize: 16,
			want:    2056,
		},
		"8bit implicit duodecimal out of range": {
			str:     "0d194",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"explicit octal": {
			str:     "67",
			base:    8,
			bitSize: 8,
			want:    55,
		},
		"implicit octal": {
			str:     "0o67",
			base:    0,
			bitSize: 8,
			want:    55,
		},
		"implicit octal with underscores": {
			str:     "0o12_34",
			base:    0,
			bitSize: 16,
			want:    668,
		},
		"8bit implicit octal out of range": {
			str:     "0o400",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"explicit quaternary": {
			str:     "33",
			base:    4,
			bitSize: 8,
			want:    15,
		},
		"implicit quaternary": {
			str:     "0q33",
			base:    0,
			bitSize: 8,
			want:    15,
		},
		"implicit quaternary with underscores": {
			str:     "0q12_33",
			base:    0,
			bitSize: 8,
			want:    111,
		},
		"8bit implicit quaternary out of range": {
			str:     "0q10000",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"explicit binary": {
			str:     "101",
			base:    2,
			bitSize: 8,
			want:    5,
		},
		"implicit binary": {
			str:     "0b101",
			base:    0,
			bitSize: 8,
			want:    5,
		},
		"implicit binary with underscores": {
			str:     "0b100_111",
			base:    0,
			bitSize: 16,
			want:    39,
		},
		"8bit implicit binary out of range": {
			str:     "0b100000000",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"64bit decimal": {
			str:     "123",
			base:    10,
			bitSize: 64,
			want:    123,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictParseUint(tc.str, tc.base, tc.bitSize)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}),
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
