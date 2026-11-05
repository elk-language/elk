package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestSize_InBytes(t *testing.T) {
	tests := map[string]struct {
		s    value.Size
		want value.Float
	}{
		"returns self in bytes as float": {
			s:    value.Byte,
			want: value.Float(1),
		},
		"returns self in kilobytes as float": {
			s:    value.Kilobyte,
			want: value.Float(1024),
		},
		"returns self in megabytes as float": {
			s:    value.Megabyte,
			want: value.Float(1048576),
		},
		"returns self in gigabytes as float": {
			s:    value.Gigabyte,
			want: value.Float(1073741824),
		},
		"returns self in terabytes as float": {
			s:    value.Terabyte,
			want: value.Float(1099511627776),
		},
		"returns self in petabytes as float": {
			s:    value.Petabyte,
			want: value.Float(1125899906842624),
		},
		"handles complex values": {
			s:    8*value.Byte + 4*value.Kilobyte + 8*value.Megabyte + 16*value.Gigabyte + 32*value.Terabyte + 64*value.Petabyte,
			want: value.Float(72092795598278656), // ? what
		},
		"what1": {
			s:    4*value.Kilobyte + 8*value.Megabyte + 16*value.Gigabyte + 32*value.Terabyte + 64*value.Petabyte,
			want: value.Float(72092795598278656), // ? what
		},
		"what2": {
			s:    -8*value.Byte + 4*value.Kilobyte + 8*value.Megabyte + 16*value.Gigabyte + 32*value.Terabyte + 64*value.Petabyte,
			want: value.Float(72092795598278656), // ? what
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.InBytes()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSize_InKilobytes(t *testing.T) {
	tests := map[string]struct {
		s    value.Size
		want value.Float
	}{
		"returns self in bytes as float": {
			s:    value.Byte,
			want: value.Float(0.0009765625),
		},
		"returns self in kilobytes as float": {
			s:    value.Kilobyte,
			want: value.Float(1),
		},
		"returns self in megabytes as float": {
			s:    value.Megabyte,
			want: value.Float(1024),
		},
		"returns self in gigabytes as float": {
			s:    value.Gigabyte,
			want: value.Float(1048576),
		},
		"returns self in terabytes as float": {
			s:    value.Terabyte,
			want: value.Float(1073741824),
		},
		"returns self in petabytes as float": {
			s:    value.Petabyte,
			want: value.Float(1099511627776),
		},
		"handles complex values": {
			s:    64*value.Byte + 32*value.Kilobyte + 16*value.Megabyte + 8*value.Gigabyte + 4*value.Terabyte + 2*value.Petabyte,
			want: value.Float(2.2033266278720625e+12),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.InKilobytes()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSize_InMegabytes(t *testing.T) {
	tests := map[string]struct {
		s    value.Size
		want value.Float
	}{
		"returns self in bytes as float": {
			s:    value.Byte,
			want: value.Float(9.5367431640625e-07),
		},
		"returns self in kilobytes as float": {
			s:    value.Kilobyte,
			want: value.Float(0.0009765625),
		},
		"returns self in megabytes as float": {
			s:    value.Megabyte,
			want: value.Float(1),
		},
		"returns self in gigabytes as float": {
			s:    value.Gigabyte,
			want: value.Float(1024),
		},
		"returns self in terabytes as float": {
			s:    value.Terabyte,
			want: value.Float(1048576),
		},
		"returns self in petabytes as float": {
			s:    value.Petabyte,
			want: value.Float(1073741824),
		},
		"handles complex values": {
			s:    64*value.Byte + 32*value.Kilobyte + 16*value.Megabyte + 8*value.Gigabyte + 4*value.Terabyte + 2*value.Petabyte,
			want: value.Float(2.151686160031311e+09),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.InMegabytes()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSize_InGigabytes(t *testing.T) {
	tests := map[string]struct {
		s    value.Size
		want value.Float
	}{
		"returns self in bytes as float": {
			s:    value.Byte,
			want: value.Float(9.313225746154785e-10),
		},
		"returns self in kilobytes as float": {
			s:    value.Kilobyte,
			want: value.Float(9.5367431640625e-07),
		},
		"returns self in megabytes as float": {
			s:    value.Megabyte,
			want: value.Float(0.0009765625),
		},
		"returns self in gigabytes as float": {
			s:    value.Gigabyte,
			want: value.Float(1),
		},
		"returns self in terabytes as float": {
			s:    value.Terabyte,
			want: value.Float(1024),
		},
		"returns self in petabytes as float": {
			s:    value.Petabyte,
			want: value.Float(1048576),
		},
		"handles complex values": {
			s:    64*value.Byte + 32*value.Kilobyte + 16*value.Megabyte + 8*value.Gigabyte + 4*value.Terabyte + 2*value.Petabyte,
			want: value.Float(2.101256015655577e+06),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.InGigabytes()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSize_InTerabytes(t *testing.T) {
	tests := map[string]struct {
		s    value.Size
		want value.Float
	}{
		"returns self in bytes as float": {
			s:    value.Byte,
			want: value.Float(9.094947017729282e-13),
		},
		"returns self in kilobytes as float": {
			s:    value.Kilobyte,
			want: value.Float(9.313225746154785e-10),
		},
		"returns self in megabytes as float": {
			s:    value.Megabyte,
			want: value.Float(9.5367431640625e-07),
		},
		"returns self in gigabytes as float": {
			s:    value.Gigabyte,
			want: value.Float(0.0009765625),
		},
		"returns self in terabytes as float": {
			s:    value.Terabyte,
			want: value.Float(1),
		},
		"returns self in petabytes as float": {
			s:    value.Petabyte,
			want: value.Float(1024),
		},
		"handles complex values": {
			s:    64*value.Byte + 32*value.Kilobyte + 16*value.Megabyte + 8*value.Gigabyte + 4*value.Terabyte + 2*value.Petabyte,
			want: value.Float(2052.0078277886496),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.InTerabytes()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSize_InPetabytes(t *testing.T) {
	tests := map[string]struct {
		s    value.Size
		want value.Float
	}{
		"returns self in bytes as float": {
			s:    value.Byte,
			want: value.Float(8.881784197001252e-16),
		},
		"returns self in kilobytes as float": {
			s:    value.Kilobyte,
			want: value.Float(9.094947017729282e-13),
		},
		"returns self in megabytes as float": {
			s:    value.Megabyte,
			want: value.Float(9.313225746154785e-10),
		},
		"returns self in gigabytes as float": {
			s:    value.Gigabyte,
			want: value.Float(9.5367431640625e-07),
		},
		"returns self in terabytes as float": {
			s:    value.Terabyte,
			want: value.Float(0.0009765625),
		},
		"returns self in petabytes as float": {
			s:    value.Petabyte,
			want: value.Float(1),
		},
		"handles complex values": {
			s:    64*value.Byte + 32*value.Kilobyte + 16*value.Megabyte + 8*value.Gigabyte + 4*value.Terabyte + 2*value.Petabyte,
			want: value.Float(2.003913894324853),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.InPetabytes()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSize_Bytes(t *testing.T) {
	tests := map[string]struct {
		s    value.Size
		want value.Value
	}{
		"returns self as total number of bytes as int": {
			s:    value.Byte,
			want: value.ToElkInt(1),
		},
		"returns self as total number of kilobytes as int": {
			s:    value.Kilobyte,
			want: value.ToElkInt(1024),
		},
		"returns self as total number of megabytes as int": {
			s:    value.Megabyte,
			want: value.ToElkInt(1048576),
		},
		"returns self as total number of gigabytes as int": {
			s:    value.Gigabyte,
			want: value.ToElkInt(1073741824),
		},
		"returns self as total number of terabytes as int": {
			s:    value.Terabyte,
			want: value.ToElkInt(1099511627776),
		},
		"returns self as total number of petabytes as int": {
			s:    value.Petabyte,
			want: value.ToElkInt(1125899906842624),
		},
		"handles complex values": {
			s:    64*value.Byte + 32*value.Kilobyte + 16*value.Megabyte + 8*value.Gigabyte + 4*value.Terabyte + 2*value.Petabyte,
			want: value.ToElkInt(2256206466940992),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.Bytes()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSize_ToString(t *testing.T) {
	tests := map[string]struct {
		s    value.Size
		want value.String
	}{
		"returns string for 0 bytes size": {
			s:    0 * value.Byte,
			want: "0.bytes",
		},
		"returns string with only bytes": {
			s:    512 * value.Byte,
			want: "512.bytes",
		},
		"returns string with only kilobytes": {
			s:    512 * value.Kilobyte,
			want: "512.kilobytes",
		},
		"returns string with only megabytes": {
			s:    512 * value.Megabyte,
			want: "512.megabytes",
		},
		"returns string with only gigabytes": {
			s:    512 * value.Gigabyte,
			want: "512.gigabytes",
		},
		"returns string with only terabytes": {
			s:    512 * value.Terabyte,
			want: "512.terabytes",
		},
		"returns string with only petabytes": {
			s:    512 * value.Petabyte,
			want: "512.petabytes",
		},
		"returns string with overflowing sizes": {
			s:    2020*value.Byte + 2020*value.Kilobyte + 2020*value.Megabyte + 2020*value.Gigabyte + 2020*value.Terabyte,
			want: "1.petabytes + 997.terabytes + 997.gigabytes + 997.megabytes + 997.kilobytes + 996.bytes",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.ToString()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
