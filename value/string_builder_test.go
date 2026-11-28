package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestInspectStringBuilder(t *testing.T) {
	tests := map[string]struct {
		builder *value.StringBuilder
		want    string
	}{
		"inspect empty": {
			builder: value.NewStringBuilder(),
			want:    `Std::StringBuilder("")`,
		},

		"inspect after writing string": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.WriteString(value.String("foo"))
				return buf
			}(),
			want: `Std::StringBuilder("foo")`,
		},
		"inspect after writing string passed as value": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.Write(value.String("foo").ToValue())
				return buf
			}(),
			want: `Std::StringBuilder("foo")`,
		},

		"inspect after writing char": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.WriteChar(value.Char('f'))
				return buf
			}(),
			want: `Std::StringBuilder("f")`,
		},
		"inspect after writing char passed as value": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.Write(value.Char('f').ToValue())
				return buf
			}(),
			want: `Std::StringBuilder("f")`,
		},

		"inspect after writing byte": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.WriteByte(value.UInt8('\x02'))
				return buf
			}(),
			want: `Std::StringBuilder("\x02")`,
		},
		"inspect after writing byte passed as value": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.Write(value.UInt8('\x02').ToValue())
				return buf
			}(),
			want: `Std::StringBuilder("\x02")`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.builder.Inspect()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestStringifyStringBuilder(t *testing.T) {
	tests := map[string]struct {
		builder *value.StringBuilder
		want    value.String
	}{
		"build empty string value": {
			builder: value.NewStringBuilder(),
			want:    value.String(""),
		},

		"build string value after writing string": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.WriteString(value.String("foo"))
				return buf
			}(),
			want: value.String("foo"),
		},
		"build string value after writing string passed as value": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.Write(value.String("foo").ToValue())
				return buf
			}(),
			want: value.String("foo"),
		},

		"build string value after writing char": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.WriteChar(value.Char('f'))
				return buf
			}(),
			want: value.String("f"),
		},
		"build string value after writing char passed as value": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.Write(value.Char('f').ToValue())
				return buf
			}(),
			want: value.String("f"),
		},

		"build string value after writing byte": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.WriteByte(value.UInt8('\x02'))
				return buf
			}(),
			want: value.String("\x02"),
		},
		"build string value after writing byte passed as value": {
			builder: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.Write(value.UInt8('\x02').ToValue())
				return buf
			}(),
			want: value.String("\x02"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.builder.ToString()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestStringBuilderGoMethods(t *testing.T) {
	tests := map[string]struct {
		initial  *value.StringBuilder
		change   func(*value.StringBuilder)
		length   int
		capacity int
	}{
		"check length and capacity": {
			initial:  value.NewStringBuilder(),
			change:   func(buf *value.StringBuilder) {},
			length:   0,
			capacity: 0,
		},
		"check length and capacity after writing string": {
			initial: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.WriteString(value.String("foo"))
				return buf
			}(),
			change: func(buf *value.StringBuilder) {
				buf.WriteString(value.String("bar"))
			},
			length:   6,
			capacity: 8,
		},
		"check length and capacity after writing char": {
			initial: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.WriteString(value.String("foo"))
				return buf
			}(),
			change: func(buf *value.StringBuilder) {
				buf.WriteChar(value.Char('b'))
			},
			length:   4,
			capacity: 8,
		},
		"check length and capacity after writing byte": {
			initial: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.WriteString(value.String("foo"))
				return buf
			}(),
			change: func(buf *value.StringBuilder) {
				buf.WriteByte(value.UInt8('\x02'))
			},
			length:   4,
			capacity: 8,
		},
		"check length and capacity after exceeding original capacity": {
			initial: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.WriteString(value.String("foo"))
				return buf
			}(),
			change: func(buf *value.StringBuilder) {
				buf.WriteString(value.String("_bar!?"))
			},
			length:   9,
			capacity: 16,
		},
		"check length and capacity after manual capacity increase": {
			initial: value.NewStringBuilder(),
			change: func(buf *value.StringBuilder) {
				buf.Grow(9)
			},
			length:   0,
			capacity: 16,
		},
		"check length and capacity after reset": {
			initial: func() *value.StringBuilder {
				buf := value.NewStringBuilder()
				buf.Grow(17)
				buf.WriteString(value.String("foo"))
				return buf
			}(),
			change: func(buf *value.StringBuilder) {
				buf.Reset()
			},
			length:   0,
			capacity: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.initial
			tc.change(got)

			length := got.Length()
			capacity := got.Capacity()

			lenOpts := comparer.Options()
			if diff := cmp.Diff(tc.length, length, lenOpts...); diff != "" {
				t.Fatalf("Length diff:\n%s", diff)
			}

			CapOpts := comparer.Options()
			if diff := cmp.Diff(tc.capacity, capacity, CapOpts...); diff != "" {
				t.Fatalf("Capacity diff:\n%s", diff)
			}
		})
	}
}
