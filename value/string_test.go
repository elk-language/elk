package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStringConcat(t *testing.T) {
	tests := map[string]struct {
		left  String
		right Value
		want  String
		err   *Error
	}{
		"String + String => String": {
			left:  String("foo"),
			right: String("bar"),
			want:  String("foobar"),
		},
		"String + Char => String": {
			left:  String("foo"),
			right: Char('b'),
			want:  String("foob"),
		},
		"String + Int => TypeError": {
			left:  String("foo"),
			right: Int8(5),
			err:   NewError(TypeErrorClass, `can't concat 5i8 to string "foo"`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Concat(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmp.AllowUnexported(BigFloat{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringRepeat(t *testing.T) {
	tests := map[string]struct {
		left  String
		right Value
		want  String
		err   *Error
	}{
		"String * SmallInt => String": {
			left:  String("a"),
			right: SmallInt(3),
			want:  "aaa",
		},
		"String * BigInt => OutOfRangeError": {
			left:  String("foo"),
			right: NewBigInt(3),
			err:   NewError(OutOfRangeErrorClass, `repeat count is too large 3`),
		},
		"String * Int8 => TypeError": {
			left:  String("foo"),
			right: Int8(3),
			err:   NewError(TypeErrorClass, `can't repeat a string using 3i8`),
		},
		"String * String => TypeError": {
			left:  String("foo"),
			right: String("bar"),
			err:   NewError(TypeErrorClass, `can't repeat a string using "bar"`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Repeat(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmp.AllowUnexported(BigFloat{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_RemoveSuffix(t *testing.T) {
	tests := map[string]struct {
		str    String
		suffix Value
		want   String
		err    *Error
	}{
		"return a type error when int is given": {
			str:    String("foo bar"),
			suffix: SmallInt(3),
			err:    NewError(TypeErrorClass, "`Std::SmallInt` can't be coerced into `Std::String`"),
		},
		"return a string without the given string suffix": {
			str:    String("foo bar"),
			suffix: String("bar"),
			want:   String("foo "),
		},
		"return the same string if there is no such string suffix": {
			str:    String("foo bar"),
			suffix: String("foo"),
			want:   String("foo bar"),
		},
		"return a string without the given char suffix": {
			str:    String("foo bar"),
			suffix: Char('r'),
			want:   String("foo ba"),
		},
		"return the same string if there is no such char suffix": {
			str:    String("foo bar"),
			suffix: Char('f'),
			want:   String("foo bar"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.str.RemoveSuffix(tc.suffix)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmp.AllowUnexported(BigFloat{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringByteLength(t *testing.T) {
	tests := map[string]struct {
		str  String
		want int
	}{
		"only ascii": {
			str:  "foo123",
			want: 6,
		},
		"unicode": {
			str:  "ślązak",
			want: 8,
		},
		"grapheme clusters": {
			str:  "👩‍🚀🇵🇱",
			want: 19,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.ByteLength()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringCharLength(t *testing.T) {
	tests := map[string]struct {
		str  String
		want int
	}{
		"only ascii": {
			str:  "foo123",
			want: 6,
		},
		"unicode": {
			str:  "ślązak",
			want: 6,
		},
		"grapheme clusters": {
			str:  "👩‍🚀🇵🇱",
			want: 5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.CharLength()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringGraphemeLength(t *testing.T) {
	tests := map[string]struct {
		str  String
		want int
	}{
		"only ascii": {
			str:  "foo123",
			want: 6,
		},
		"unicode": {
			str:  "ślązak",
			want: 6,
		},
		"graphemes clusters": {
			str:  "👩‍🚀🇵🇱",
			want: 2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.GraphemeLength()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringReverseBytes(t *testing.T) {
	tests := map[string]struct {
		str  String
		want String
	}{
		"only ascii": {
			str:  "foo123",
			want: "321oof",
		},
		"unicode": {
			str:  "ślązak",
			want: "kaz\x85\xc4l\x9b\xc5",
		},
		"grapheme clusters": {
			str:  "👩‍🚀🇵🇱",
			want: "\xb1\x87\x9f\U000351df\xf0\x80\x9a\x9f\xf0\x8d\x80⩑\x9f\xf0",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.ReverseBytes()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringReverseChars(t *testing.T) {
	tests := map[string]struct {
		str  String
		want String
	}{
		"only ascii": {
			str:  "foo123",
			want: "321oof",
		},
		"unicode": {
			str:  "ślązak",
			want: "kaząlś",
		},
		"grapheme clusters": {
			str:  "👩‍🚀🇵🇱",
			want: "🇱🇵🚀\u200d👩",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.ReverseChars()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringReverseGraphemes(t *testing.T) {
	tests := map[string]struct {
		str  String
		want String
	}{
		"only ascii": {
			str:  "foo123",
			want: "321oof",
		},
		"unicode": {
			str:  "ślązak",
			want: "kaząlś",
		},
		"grapheme clusters": {
			str:  "👩‍🚀🇵🇱",
			want: "🇵🇱👩‍🚀",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.ReverseGraphemes()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
