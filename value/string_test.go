package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestString_GraphemeAt(t *testing.T) {
	tests := map[string]struct {
		s    value.String
		i    value.Value
		want value.String
		err  value.Value
	}{
		"get index 0 in an empty string": {
			s: "",
			i: value.SmallInt(0).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			)),
		},
		"get index 0 in an ascii string": {
			s:    "foo",
			i:    value.SmallInt(0).ToValue(),
			want: "f",
		},
		"get index 0 in a binary string": {
			s:    "\x86foo",
			i:    value.SmallInt(0).ToValue(),
			want: "\x86",
		},
		"get index 1 in a binary string": {
			s:    "\x86foo",
			i:    value.SmallInt(1).ToValue(),
			want: "f",
		},
		"get index 1 in an ascii string": {
			s:    "foo",
			i:    value.SmallInt(1).ToValue(),
			want: "o",
		},
		"get index 1 in a unicode string": {
			s:    "żółć",
			i:    value.SmallInt(1).ToValue(),
			want: "ó",
		},
		"get index 0 in grapheme cluster string": {
			s:    "🇵🇱🥟👨🏻‍💻",
			i:    value.SmallInt(0).ToValue(),
			want: "🇵🇱",
		},
		"get index 1 in grapheme cluster string": {
			s:    "🇵🇱🥟👨🏻‍💻",
			i:    value.SmallInt(1).ToValue(),
			want: "🥟",
		},
		"get index -1 in grapheme cluster string": {
			s:    "🇵🇱🥟👨🏻‍💻",
			i:    value.SmallInt(-1).ToValue(),
			want: "👨🏻‍💻",
		},
		"get index -1 in a unicode string": {
			s:    "żółć",
			i:    value.SmallInt(-1).ToValue(),
			want: "ć",
		},
		"get index -2 in a unicode string": {
			s:    "żółć",
			i:    value.SmallInt(-2).ToValue(),
			want: "ł",
		},
		"get positive index out of range": {
			s: "żółć",
			i: value.SmallInt(25).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 25 out of range: -4...4",
			)),
		},
		"get negative index out of range": {
			s: "żółć",
			i: value.SmallInt(-25).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index -25 out of range: -4...4",
			)),
		},
		"get uint8 index": {
			s:    "żółć",
			i:    value.UInt8(1).ToValue(),
			want: "ó",
		},
		"get int16 index": {
			s:    "żółć",
			i:    value.Int16(1).ToValue(),
			want: "ó",
		},
		"get string index": {
			s: "żółć",
			i: value.Ref(value.String("lol")),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			)),
		},
		"get float index": {
			s: "żółć",
			i: value.Float(3).ToValue(),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int`",
			)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.s.GraphemeAt(tc.i)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if !tc.err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_Subscript(t *testing.T) {
	tests := map[string]struct {
		s    value.String
		i    value.Value
		want value.Char
		err  value.Value
	}{
		"get index 0 in an empty string": {
			s: "",
			i: value.SmallInt(0).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			)),
		},
		"get index 0 in an ascii string": {
			s:    "foo",
			i:    value.SmallInt(0).ToValue(),
			want: 'f',
		},
		"get index 0 in a binary string": {
			s:    "\x86foo",
			i:    value.SmallInt(0).ToValue(),
			want: '\x86',
		},
		"get index 1 in a binary string": {
			s:    "\x86foo",
			i:    value.SmallInt(1).ToValue(),
			want: 'f',
		},
		"get index 1 in an ascii string": {
			s:    "foo",
			i:    value.SmallInt(1).ToValue(),
			want: 'o',
		},
		"get index 1 in a unicode string": {
			s:    "żółć",
			i:    value.SmallInt(1).ToValue(),
			want: 'ó',
		},
		"get index -1 in a unicode string": {
			s:    "żółć",
			i:    value.SmallInt(-1).ToValue(),
			want: 'ć',
		},
		"get index -2 in a unicode string": {
			s:    "żółć",
			i:    value.SmallInt(-2).ToValue(),
			want: 'ł',
		},
		"get positive index out of range": {
			s: "żółć",
			i: value.SmallInt(25).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 25 out of range: -4...4",
			)),
		},
		"get negative index out of range": {
			s: "żółć",
			i: value.SmallInt(-25).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index -25 out of range: -4...4",
			)),
		},
		"get uint8 index": {
			s:    "żółć",
			i:    value.UInt8(1).ToValue(),
			want: 'ó',
		},
		"get int16 index": {
			s:    "żółć",
			i:    value.Int16(1).ToValue(),
			want: 'ó',
		},
		"get string index": {
			s: "żółć",
			i: value.Ref(value.String("lol")),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			)),
		},
		"get float index": {
			s: "żółć",
			i: value.Float(3).ToValue(),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int`",
			)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.s.Subscript(tc.i)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if !tc.err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_ByteAt(t *testing.T) {
	tests := map[string]struct {
		s    value.String
		i    value.Value
		want value.UInt8
		err  value.Value
	}{
		"get index 0 in an empty string": {
			s: "",
			i: value.SmallInt(0).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			)),
		},
		"get index 0 in an ascii string": {
			s:    "foo",
			i:    value.SmallInt(0).ToValue(),
			want: 'f',
		},
		"get index 1 in an ascii string": {
			s:    "foo",
			i:    value.SmallInt(1).ToValue(),
			want: 'o',
		},
		"get index 1 in a unicode string": {
			s:    "żółć",
			i:    value.SmallInt(1).ToValue(),
			want: '\xbc',
		},
		"get index -1 in a unicode string": {
			s:    "żółć",
			i:    value.SmallInt(-1).ToValue(),
			want: '\x87',
		},
		"get index -2 in a unicode string": {
			s:    "żółć",
			i:    value.SmallInt(-2).ToValue(),
			want: '\xc4',
		},
		"get positive index out of range": {
			s: "żółć",
			i: value.SmallInt(25).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 25 out of range: -8...8",
			)),
		},
		"get negative index out of range": {
			s: "żółć",
			i: value.SmallInt(-25).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index -25 out of range: -8...8",
			)),
		},
		"get uint8 index": {
			s:    "foo",
			i:    value.UInt8(1).ToValue(),
			want: 'o',
		},
		"get int16 index": {
			s:    "foo",
			i:    value.Int16(1).ToValue(),
			want: 'o',
		},
		"get string index": {
			s: "foo",
			i: value.Ref(value.String("lol")),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			)),
		},
		"get float index": {
			s: "foo",
			i: value.Float(3).ToValue(),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int`",
			)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.s.ByteAt(tc.i)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if !tc.err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_Inspect(t *testing.T) {
	tests := map[string]struct {
		s    value.String
		want string
	}{
		"ascii letter": {
			s:    "d",
			want: `"d"`,
		},
		"utf-8 character": {
			s:    "ślężak",
			want: `"ślężak"`,
		},
		"newline": {
			s:    "\n",
			want: `"\n"`,
		},
		"double quote": {
			s:    `"`,
			want: `"\""`,
		},
		"dollar": {
			s:    `$foo`,
			want: `"\$foo"`,
		},
		"pound": {
			s:    `#foo`,
			want: `"\#foo"`,
		},
		"backslash": {
			s:    `\`,
			want: `"\\"`,
		},
		"hex byte": {
			s:    "\x02",
			want: `"\x02"`,
		},
		"unicode codepoint": {
			s:    "\U0010FFFF",
			want: `"\U0010FFFF"`,
		},
		"small unicode codepoint": {
			s:    "\u200d",
			want: `"\u200d"`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringConcat(t *testing.T) {
	tests := map[string]struct {
		left  value.String
		right value.Value
		want  value.String
		err   value.Value
	}{
		"String + String => String": {
			left:  value.String("foo"),
			right: value.Ref(value.String("bar")),
			want:  value.String("foobar"),
		},
		"String + Char => String": {
			left:  value.String("foo"),
			right: value.Char('b').ToValue(),
			want:  value.String("foob"),
		},
		"String + Int => TypeError": {
			left:  value.String("foo"),
			right: value.Int8(5).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot concat 5i8 to string "foo"`)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Concat(tc.right)
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

func TestStringRepeat(t *testing.T) {
	tests := map[string]struct {
		left  value.String
		right value.Value
		want  value.String
		err   value.Value
	}{
		"String * SmallInt => String": {
			left:  value.String("a"),
			right: value.SmallInt(3).ToValue(),
			want:  "aaa",
		},
		"String * 0 => String": {
			left:  value.String("a"),
			right: value.SmallInt(0).ToValue(),
			want:  "",
		},
		"String * -SmallInt => OutOfRangeError": {
			left:  value.String("a"),
			right: value.SmallInt(-3).ToValue(),
			err:   value.Ref(value.NewError(value.OutOfRangeErrorClass, `repeat count cannot be negative: -3`)),
		},
		"String * BigInt => OutOfRangeError": {
			left:  value.String("foo"),
			right: value.Ref(value.NewBigInt(3)),
			err:   value.Ref(value.NewError(value.OutOfRangeErrorClass, `repeat count is too large 3`)),
		},
		"String * Int8 => TypeError": {
			left:  value.String("foo"),
			right: value.Int8(3).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot repeat a string using 3i8`)),
		},
		"String * String => TypeError": {
			left:  value.String("foo"),
			right: value.Ref(value.String("bar")),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot repeat a string using "bar"`)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Repeat(tc.right)
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

func TestString_RemoveSuffix(t *testing.T) {
	tests := map[string]struct {
		str    value.String
		suffix value.Value
		want   value.String
		err    value.Value
	}{
		"return a type error when int is given": {
			str:    value.String("foo bar"),
			suffix: value.SmallInt(3).ToValue(),
			err:    value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::String`")),
		},
		"return a string without the given string suffix": {
			str:    value.String("foo bar"),
			suffix: value.Ref(value.String("bar")),
			want:   value.String("foo "),
		},
		"return the same string if there is no such string suffix": {
			str:    value.String("foo bar"),
			suffix: value.Ref(value.String("foo")),
			want:   value.String("foo bar"),
		},
		"return a string without the given char suffix": {
			str:    value.String("foo bar"),
			suffix: value.Char('r').ToValue(),
			want:   value.String("foo ba"),
		},
		"return the same string if there is no such char suffix": {
			str:    value.String("foo bar"),
			suffix: value.Char('f').ToValue(),
			want:   value.String("foo bar"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.str.RemoveSuffix(tc.suffix)
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

func TestStringByteCount(t *testing.T) {
	tests := map[string]struct {
		str  value.String
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
			got := tc.str.ByteCount()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringCharCount(t *testing.T) {
	tests := map[string]struct {
		str  value.String
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
			got := tc.str.CharCount()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringGraphemeCount(t *testing.T) {
	tests := map[string]struct {
		str  value.String
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
			got := tc.str.GraphemeCount()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringReverseBytes(t *testing.T) {
	tests := map[string]struct {
		str  value.String
		want value.String
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
		str  value.String
		want value.String
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
		str  value.String
		want value.String
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

func TestString_Compare(t *testing.T) {
	tests := map[string]struct {
		a    value.String
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt and return an error": {
			a:   value.String("a"),
			b:   value.SmallInt(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::String`")),
		},
		"Float and return an error": {
			a:   value.String("a"),
			b:   value.Float(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::String`")),
		},
		"BigFloat and return an error": {
			a:   value.String("a"),
			b:   value.Ref(value.NewBigFloat(5)),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::BigFloat` cannot be coerced into `Std::String`")),
		},
		"Int64 and return an error": {
			a:   value.String("a"),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::String`")),
		},
		"Int32 and return an error": {
			a:   value.String("a"),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::String`")),
		},
		"Int16 and return an error": {
			a:   value.String("a"),
			b:   value.Int16(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int16` cannot be coerced into `Std::String`")),
		},
		"Int8 and return an error": {
			a:   value.String("a"),
			b:   value.Int8(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int8` cannot be coerced into `Std::String`")),
		},
		"UInt64 and return an error": {
			a:   value.String("a"),
			b:   value.UInt64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt64` cannot be coerced into `Std::String`")),
		},
		"UInt32 and return an error": {
			a:   value.String("a"),
			b:   value.UInt32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt32` cannot be coerced into `Std::String`")),
		},
		"UInt16 and return an error": {
			a:   value.String("a"),
			b:   value.UInt16(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt16` cannot be coerced into `Std::String`")),
		},
		"UInt8 and return an error": {
			a:   value.String("a"),
			b:   value.UInt8(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt8` cannot be coerced into `Std::String`")),
		},
		"Float64 and return an error": {
			a:   value.String("a"),
			b:   value.Float64(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::String`")),
		},
		"Float32 and return an error": {
			a:   value.String("a"),
			b:   value.Float32(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::String`")),
		},
		"String 'a' <=> 'a'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("a")),
			want: value.SmallInt(0).ToValue(),
		},
		"String 'a' <=> 'b'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("b")),
			want: value.SmallInt(-1).ToValue(),
		},
		"String 'b' <=> 'a'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("a")),
			want: value.SmallInt(1).ToValue(),
		},
		"String 'aa' <=> 'a'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("a")),
			want: value.SmallInt(1).ToValue(),
		},
		"String 'a' <=> 'aa'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("aa")),
			want: value.SmallInt(-1).ToValue(),
		},
		"String 'aa' <=> 'b'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("b")),
			want: value.SmallInt(-1).ToValue(),
		},
		"String 'b' <=> 'aa'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("aa")),
			want: value.SmallInt(1).ToValue(),
		},
		"String 'abdf' <=> 'abcf'": {
			a:    value.String("abdf"),
			b:    value.Ref(value.String("abcf")),
			want: value.SmallInt(1).ToValue(),
		},
		"String 'abcf' <=> 'abdf'": {
			a:    value.String("abcf"),
			b:    value.Ref(value.String("abdf")),
			want: value.SmallInt(-1).ToValue(),
		},
		"String 'ś' <=> 'ą'": {
			a:    value.String("ś"),
			b:    value.Ref(value.String("ą")),
			want: value.SmallInt(1).ToValue(),
		},
		"String 'ą' <=> 'ś'": {
			a:    value.String("ą"),
			b:    value.Ref(value.String("ś")),
			want: value.SmallInt(-1).ToValue(),
		},

		"Char 'a' <=> `a`": {
			a:    value.String("a"),
			b:    value.Char('a').ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"Char 'a' <=> `b`": {
			a:    value.String("a"),
			b:    value.Char('b').ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"Char 'b' <=> `a`": {
			a:    value.String("b"),
			b:    value.Char('a').ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"Char 'aa' <=> `a`": {
			a:    value.String("aa"),
			b:    value.Char('a').ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"Char 'aa' <=> `b`": {
			a:    value.String("aa"),
			b:    value.Char('b').ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"Char 'ś' <=> `ą`": {
			a:    value.String("ś"),
			b:    value.Char('ą').ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"Char 'ą' <=> `ś`": {
			a:    value.String("ą"),
			b:    value.Char('ś').ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Compare(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.String
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt and return an error": {
			a:   value.String("a"),
			b:   value.SmallInt(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::String`")),
		},
		"Float and return an error": {
			a:   value.String("a"),
			b:   value.Float(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::String`")),
		},
		"BigFloat and return an error": {
			a:   value.String("a"),
			b:   value.Ref(value.NewBigFloat(5)),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::BigFloat` cannot be coerced into `Std::String`")),
		},
		"Int64 and return an error": {
			a:   value.String("a"),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::String`")),
		},
		"Int32 and return an error": {
			a:   value.String("a"),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::String`")),
		},
		"Int16 and return an error": {
			a:   value.String("a"),
			b:   value.Int16(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int16` cannot be coerced into `Std::String`")),
		},
		"Int8 and return an error": {
			a:   value.String("a"),
			b:   value.Int8(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int8` cannot be coerced into `Std::String`")),
		},
		"UInt64 and return an error": {
			a:   value.String("a"),
			b:   value.UInt64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt64` cannot be coerced into `Std::String`")),
		},
		"UInt32 and return an error": {
			a:   value.String("a"),
			b:   value.UInt32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt32` cannot be coerced into `Std::String`")),
		},
		"UInt16 and return an error": {
			a:   value.String("a"),
			b:   value.UInt16(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt16` cannot be coerced into `Std::String`")),
		},
		"UInt8 and return an error": {
			a:   value.String("a"),
			b:   value.UInt8(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt8` cannot be coerced into `Std::String`")),
		},
		"Float64 and return an error": {
			a:   value.String("a"),
			b:   value.Float64(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::String`")),
		},
		"Float32 and return an error": {
			a:   value.String("a"),
			b:   value.Float32(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::String`")),
		},

		"String 'a' > 'a'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("a")),
			want: value.False,
		},
		"String 'a' > 'b'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("b")),
			want: value.False,
		},
		"String 'b' > 'a'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("a")),
			want: value.True,
		},
		"String 'aa' > 'a'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("a")),
			want: value.True,
		},
		"String 'a' > 'aa'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("aa")),
			want: value.False,
		},
		"String 'aa' > 'b'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("b")),
			want: value.False,
		},
		"String 'b' > 'aa'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("aa")),
			want: value.True,
		},
		"String 'abdf' > 'abcf'": {
			a:    value.String("abdf"),
			b:    value.Ref(value.String("abcf")),
			want: value.True,
		},
		"String 'abcf' > 'abdf'": {
			a:    value.String("abcf"),
			b:    value.Ref(value.String("abdf")),
			want: value.False,
		},
		"String 'ś' > 'ą'": {
			a:    value.String("ś"),
			b:    value.Ref(value.String("ą")),
			want: value.True,
		},
		"String 'ą' > 'ś'": {
			a:    value.String("ą"),
			b:    value.Ref(value.String("ś")),
			want: value.False,
		},

		"Char 'a' > `a`": {
			a:    value.String("a"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'a' > `b`": {
			a:    value.String("a"),
			b:    value.Char('b').ToValue(),
			want: value.False,
		},
		"Char 'b' > `a`": {
			a:    value.String("b"),
			b:    value.Char('a').ToValue(),
			want: value.True,
		},
		"Char 'aa' > `a`": {
			a:    value.String("aa"),
			b:    value.Char('a').ToValue(),
			want: value.True,
		},
		"Char 'aa' > `b`": {
			a:    value.String("aa"),
			b:    value.Char('b').ToValue(),
			want: value.False,
		},
		"Char 'ś' > `ą`": {
			a:    value.String("ś"),
			b:    value.Char('ą').ToValue(),
			want: value.True,
		},
		"Char 'ą' > `ś`": {
			a:    value.String("ą"),
			b:    value.Char('ś').ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThan(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.String
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt and return an error": {
			a:   value.String("a"),
			b:   value.SmallInt(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::String`")),
		},
		"Float and return an error": {
			a:   value.String("a"),
			b:   value.Float(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::String`")),
		},
		"BigFloat and return an error": {
			a:   value.String("a"),
			b:   value.Ref(value.NewBigFloat(5)),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::BigFloat` cannot be coerced into `Std::String`")),
		},
		"Int64 and return an error": {
			a:   value.String("a"),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::String`")),
		},
		"Int32 and return an error": {
			a:   value.String("a"),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::String`")),
		},
		"Int16 and return an error": {
			a:   value.String("a"),
			b:   value.Int16(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int16` cannot be coerced into `Std::String`")),
		},
		"Int8 and return an error": {
			a:   value.String("a"),
			b:   value.Int8(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int8` cannot be coerced into `Std::String`")),
		},
		"UInt64 and return an error": {
			a:   value.String("a"),
			b:   value.UInt64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt64` cannot be coerced into `Std::String`")),
		},
		"UInt32 and return an error": {
			a:   value.String("a"),
			b:   value.UInt32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt32` cannot be coerced into `Std::String`")),
		},
		"UInt16 and return an error": {
			a:   value.String("a"),
			b:   value.UInt16(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt16` cannot be coerced into `Std::String`")),
		},
		"UInt8 and return an error": {
			a:   value.String("a"),
			b:   value.UInt8(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt8` cannot be coerced into `Std::String`")),
		},
		"Float64 and return an error": {
			a:   value.String("a"),
			b:   value.Float64(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::String`")),
		},
		"Float32 and return an error": {
			a:   value.String("a"),
			b:   value.Float32(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::String`")),
		},
		"String 'a' >= 'a'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("a")),
			want: value.True,
		},
		"String 'foo' >= 'foo'": {
			a:    value.String("foo"),
			b:    value.Ref(value.String("foo")),
			want: value.True,
		},
		"String 'a' >= 'b'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("b")),
			want: value.False,
		},
		"String 'b' >= 'a'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("a")),
			want: value.True,
		},
		"String 'aa' >= 'a'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("a")),
			want: value.True,
		},
		"String 'a' >= 'aa'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("aa")),
			want: value.False,
		},
		"String 'aa' >= 'b'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("b")),
			want: value.False,
		},
		"String 'b' >= 'aa'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("aa")),
			want: value.True,
		},
		"String 'abdf' >= 'abcf'": {
			a:    value.String("abdf"),
			b:    value.Ref(value.String("abcf")),
			want: value.True,
		},
		"String 'abcf' >= 'abdf'": {
			a:    value.String("abcf"),
			b:    value.Ref(value.String("abdf")),
			want: value.False,
		},
		"String 'ś' >= 'ą'": {
			a:    value.String("ś"),
			b:    value.Ref(value.String("ą")),
			want: value.True,
		},
		"String 'ą' >= 'ś'": {
			a:    value.String("ą"),
			b:    value.Ref(value.String("ś")),
			want: value.False,
		},

		"Char 'a' >= `a`": {
			a:    value.String("a"),
			b:    value.Char('a').ToValue(),
			want: value.True,
		},
		"Char 'a' >= `b`": {
			a:    value.String("a"),
			b:    value.Char('b').ToValue(),
			want: value.False,
		},
		"Char 'b' >= `a`": {
			a:    value.String("b"),
			b:    value.Char('a').ToValue(),
			want: value.True,
		},
		"Char 'aa' >= `a`": {
			a:    value.String("aa"),
			b:    value.Char('a').ToValue(),
			want: value.True,
		},
		"Char 'aa' >= `b`": {
			a:    value.String("aa"),
			b:    value.Char('b').ToValue(),
			want: value.False,
		},
		"Char 'ś' >= `ą`": {
			a:    value.String("ś"),
			b:    value.Char('ą').ToValue(),
			want: value.True,
		},
		"Char 'ą' >= `ś`": {
			a:    value.String("ą"),
			b:    value.Char('ś').ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqual(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.String
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt and return an error": {
			a:   value.String("a"),
			b:   value.SmallInt(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::String`")),
		},
		"Float and return an error": {
			a:   value.String("a"),
			b:   value.Float(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::String`")),
		},
		"BigFloat and return an error": {
			a:   value.String("a"),
			b:   value.Ref(value.NewBigFloat(5)),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::BigFloat` cannot be coerced into `Std::String`")),
		},
		"Int64 and return an error": {
			a:   value.String("a"),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::String`")),
		},
		"Int32 and return an error": {
			a:   value.String("a"),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::String`")),
		},
		"Int16 and return an error": {
			a:   value.String("a"),
			b:   value.Int16(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int16` cannot be coerced into `Std::String`")),
		},
		"Int8 and return an error": {
			a:   value.String("a"),
			b:   value.Int8(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int8` cannot be coerced into `Std::String`")),
		},
		"UInt64 and return an error": {
			a:   value.String("a"),
			b:   value.UInt64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt64` cannot be coerced into `Std::String`")),
		},
		"UInt32 and return an error": {
			a:   value.String("a"),
			b:   value.UInt32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt32` cannot be coerced into `Std::String`")),
		},
		"UInt16 and return an error": {
			a:   value.String("a"),
			b:   value.UInt16(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt16` cannot be coerced into `Std::String`")),
		},
		"UInt8 and return an error": {
			a:   value.String("a"),
			b:   value.UInt8(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt8` cannot be coerced into `Std::String`")),
		},
		"Float64 and return an error": {
			a:   value.String("a"),
			b:   value.Float64(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::String`")),
		},
		"Float32 and return an error": {
			a:   value.String("a"),
			b:   value.Float32(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::String`")),
		},
		"String 'a' < 'a'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("a")),
			want: value.False,
		},
		"String 'foo' < 'foo'": {
			a:    value.String("foo"),
			b:    value.Ref(value.String("foo")),
			want: value.False,
		},
		"String 'a' < 'b'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("b")),
			want: value.True,
		},
		"String 'b' < 'a'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("a")),
			want: value.False,
		},
		"String 'aa' < 'a'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("a")),
			want: value.False,
		},
		"String 'a' < 'aa'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("aa")),
			want: value.True,
		},
		"String 'aa' < 'b'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("b")),
			want: value.True,
		},
		"String 'b' < 'aa'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("aa")),
			want: value.False,
		},
		"String 'abdf' < 'abcf'": {
			a:    value.String("abdf"),
			b:    value.Ref(value.String("abcf")),
			want: value.False,
		},
		"String 'abcf' < 'abdf'": {
			a:    value.String("abcf"),
			b:    value.Ref(value.String("abdf")),
			want: value.True,
		},
		"String 'ś' < 'ą'": {
			a:    value.String("ś"),
			b:    value.Ref(value.String("ą")),
			want: value.False,
		},
		"String 'ą' < 'ś'": {
			a:    value.String("ą"),
			b:    value.Ref(value.String("ś")),
			want: value.True,
		},

		"Char 'a' < `a`": {
			a:    value.String("a"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'a' < `b`": {
			a:    value.String("a"),
			b:    value.Char('b').ToValue(),
			want: value.True,
		},
		"Char 'b' < `a`": {
			a:    value.String("b"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'aa' < `a`": {
			a:    value.String("aa"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'aa' < `b`": {
			a:    value.String("aa"),
			b:    value.Char('b').ToValue(),
			want: value.True,
		},
		"Char 'ś' < `ą`": {
			a:    value.String("ś"),
			b:    value.Char('ą').ToValue(),
			want: value.False,
		},
		"Char 'ą' < `ś`": {
			a:    value.String("ą"),
			b:    value.Char('ś').ToValue(),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThan(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.String
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt and return an error": {
			a:   value.String("a"),
			b:   value.SmallInt(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::String`")),
		},
		"Float and return an error": {
			a:   value.String("a"),
			b:   value.Float(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::String`")),
		},
		"BigFloat and return an error": {
			a:   value.String("a"),
			b:   value.Ref(value.NewBigFloat(5)),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::BigFloat` cannot be coerced into `Std::String`")),
		},
		"Int64 and return an error": {
			a:   value.String("a"),
			b:   value.Int64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::String`")),
		},
		"Int32 and return an error": {
			a:   value.String("a"),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::String`")),
		},
		"Int16 and return an error": {
			a:   value.String("a"),
			b:   value.Int16(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int16` cannot be coerced into `Std::String`")),
		},
		"Int8 and return an error": {
			a:   value.String("a"),
			b:   value.Int8(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int8` cannot be coerced into `Std::String`")),
		},
		"UInt64 and return an error": {
			a:   value.String("a"),
			b:   value.UInt64(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt64` cannot be coerced into `Std::String`")),
		},
		"UInt32 and return an error": {
			a:   value.String("a"),
			b:   value.UInt32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt32` cannot be coerced into `Std::String`")),
		},
		"UInt16 and return an error": {
			a:   value.String("a"),
			b:   value.UInt16(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt16` cannot be coerced into `Std::String`")),
		},
		"UInt8 and return an error": {
			a:   value.String("a"),
			b:   value.UInt8(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::UInt8` cannot be coerced into `Std::String`")),
		},
		"Float64 and return an error": {
			a:   value.String("a"),
			b:   value.Float64(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::String`")),
		},
		"Float32 and return an error": {
			a:   value.String("a"),
			b:   value.Float32(5).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::String`")),
		},
		"String 'a' <= 'a'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("a")),
			want: value.True,
		},
		"String 'foo' <= 'foo'": {
			a:    value.String("foo"),
			b:    value.Ref(value.String("foo")),
			want: value.True,
		},
		"String 'a' <= 'b'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("b")),
			want: value.True,
		},
		"String 'b' <= 'a'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("a")),
			want: value.False,
		},
		"String 'aa' <= 'a'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("a")),
			want: value.False,
		},
		"String 'a' <= 'aa'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("aa")),
			want: value.True,
		},
		"String 'aa' <= 'b'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("b")),
			want: value.True,
		},
		"String 'b' <= 'aa'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("aa")),
			want: value.False,
		},
		"String 'abdf' <= 'abcf'": {
			a:    value.String("abdf"),
			b:    value.Ref(value.String("abcf")),
			want: value.False,
		},
		"String 'abcf' <= 'abdf'": {
			a:    value.String("abcf"),
			b:    value.Ref(value.String("abdf")),
			want: value.True,
		},
		"String 'ś' <= 'ą'": {
			a:    value.String("ś"),
			b:    value.Ref(value.String("ą")),
			want: value.False,
		},
		"String 'ą' <= 'ś'": {
			a:    value.String("ą"),
			b:    value.Ref(value.String("ś")),
			want: value.True,
		},

		"Char 'a' <= `a`": {
			a:    value.String("a"),
			b:    value.Char('a').ToValue(),
			want: value.True,
		},
		"Char 'a' <= `b`": {
			a:    value.String("a"),
			b:    value.Char('b').ToValue(),
			want: value.True,
		},
		"Char 'b' <= `a`": {
			a:    value.String("b"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'aa' <= `a`": {
			a:    value.String("aa"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'aa' <= `b`": {
			a:    value.String("aa"),
			b:    value.Char('b').ToValue(),
			want: value.True,
		},
		"Char 'ś' <= `ą`": {
			a:    value.String("ś"),
			b:    value.Char('ą').ToValue(),
			want: value.False,
		},
		"Char 'ą' <= `ś`": {
			a:    value.String("ą"),
			b:    value.Char('ś').ToValue(),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqual(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_LaxEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.String
		b    value.Value
		want value.Value
	}{
		"SmallInt '2' =~ 2": {
			a:    value.String("2"),
			b:    value.SmallInt(2).ToValue(),
			want: value.False,
		},
		"Float '5.2' =~ 5.2": {
			a:    value.String("5.2"),
			b:    value.Float(5.2).ToValue(),
			want: value.False,
		},
		"BigFloat '4.5' =~ 4.5bf": {
			a:    value.String("4.5"),
			b:    value.Ref(value.NewBigFloat(4.5)),
			want: value.False,
		},
		"Int64 '5' =~ 5i64": {
			a:    value.String("5"),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int32 '2' =~ 2i32": {
			a:    value.String("2"),
			b:    value.Int32(2).ToValue(),
			want: value.False,
		},
		"Int16 '25' =~ 25i16": {
			a:    value.String("25"),
			b:    value.Int16(25).ToValue(),
			want: value.False,
		},
		"Int8 '8' =~ 8i8": {
			a:    value.String("8"),
			b:    value.Int8(8).ToValue(),
			want: value.False,
		},
		"UInt64 '31' =~ 31u64": {
			a:    value.String("31"),
			b:    value.UInt64(31).ToValue(),
			want: value.False,
		},
		"UInt32 '9' =~ 9u32": {
			a:    value.String("9"),
			b:    value.UInt32(9).ToValue(),
			want: value.False,
		},
		"UInt16 '74' =~ 74u16": {
			a:    value.String("74"),
			b:    value.UInt16(74).ToValue(),
			want: value.False,
		},
		"UInt8 '12' =~ 12u8": {
			a:    value.String("12"),
			b:    value.UInt8(12).ToValue(),
			want: value.False,
		},
		"Float64 '49.2' =~ 49.2f64": {
			a:    value.String("49.2"),
			b:    value.Float64(49.2).ToValue(),
			want: value.False,
		},
		"Float32 '57.9' =~ 57.9f32": {
			a:    value.String("57.9"),
			b:    value.Float32(57.9).ToValue(),
			want: value.False,
		},
		"String 'a' =~ 'a'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("a")),
			want: value.True,
		},
		"String 'foo' =~ 'foo'": {
			a:    value.String("foo"),
			b:    value.Ref(value.String("foo")),
			want: value.True,
		},
		"String 'a' =~ 'b'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("b")),
			want: value.False,
		},
		"String 'b' =~ 'a'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("a")),
			want: value.False,
		},
		"String 'aa' =~ 'a'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("a")),
			want: value.False,
		},
		"String 'a' =~ 'aa'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("aa")),
			want: value.False,
		},
		"String 'aa' =~ 'b'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("b")),
			want: value.False,
		},
		"String 'b' =~ 'aa'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("aa")),
			want: value.False,
		},
		"String 'abdf' =~ 'abcf'": {
			a:    value.String("abdf"),
			b:    value.Ref(value.String("abcf")),
			want: value.False,
		},
		"String 'abcf' =~ 'abdf'": {
			a:    value.String("abcf"),
			b:    value.Ref(value.String("abdf")),
			want: value.False,
		},
		"String 'ś' =~ 'ą'": {
			a:    value.String("ś"),
			b:    value.Ref(value.String("ą")),
			want: value.False,
		},
		"String 'ą' =~ 'ś'": {
			a:    value.String("ą"),
			b:    value.Ref(value.String("ś")),
			want: value.False,
		},

		"Char 'a' =~ `a`": {
			a:    value.String("a"),
			b:    value.Char('a').ToValue(),
			want: value.True,
		},
		"Char 'a' =~ `b`": {
			a:    value.String("a"),
			b:    value.Char('b').ToValue(),
			want: value.False,
		},
		"Char 'b' =~ `a`": {
			a:    value.String("b"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'aa' =~ `a`": {
			a:    value.String("aa"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'aa' =~ `b`": {
			a:    value.String("aa"),
			b:    value.Char('b').ToValue(),
			want: value.False,
		},
		"Char 'ś' =~ `ą`": {
			a:    value.String("ś"),
			b:    value.Char('ą').ToValue(),
			want: value.False,
		},
		"Char 'ą' =~ `ś`": {
			a:    value.String("ą"),
			b:    value.Char('ś').ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.LaxEqual(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.String
		b    value.Value
		want value.Value
	}{
		"SmallInt '2' == 2": {
			a:    value.String("2"),
			b:    value.SmallInt(2).ToValue(),
			want: value.False,
		},
		"Float '5.2' == 5.2": {
			a:    value.String("5.2"),
			b:    value.Float(5.2).ToValue(),
			want: value.False,
		},
		"BigFloat '4.5' == 4.5bf": {
			a:    value.String("4.5"),
			b:    value.Ref(value.NewBigFloat(4.5)),
			want: value.False,
		},
		"Int64 '5' == 5i64": {
			a:    value.String("5"),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int32 '2' == 2i32": {
			a:    value.String("2"),
			b:    value.Int32(2).ToValue(),
			want: value.False,
		},
		"Int16 '25' == 25i16": {
			a:    value.String("25"),
			b:    value.Int16(25).ToValue(),
			want: value.False,
		},
		"Int8 '8' == 8i8": {
			a:    value.String("8"),
			b:    value.Int8(8).ToValue(),
			want: value.False,
		},
		"UInt64 '31' == 31u64": {
			a:    value.String("31"),
			b:    value.UInt64(31).ToValue(),
			want: value.False,
		},
		"UInt32 '9' == 9u32": {
			a:    value.String("9"),
			b:    value.UInt32(9).ToValue(),
			want: value.False,
		},
		"UInt16 '74' == 74u16": {
			a:    value.String("74"),
			b:    value.UInt16(74).ToValue(),
			want: value.False,
		},
		"UInt8 '12' == 12u8": {
			a:    value.String("12"),
			b:    value.UInt8(12).ToValue(),
			want: value.False,
		},
		"Float64 '49.2' == 49.2f64": {
			a:    value.String("49.2"),
			b:    value.Float64(49.2).ToValue(),
			want: value.False,
		},
		"Float32 '57.9' == 57.9f32": {
			a:    value.String("57.9"),
			b:    value.Float32(57.9).ToValue(),
			want: value.False,
		},

		"String 'a' == 'a'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("a")),
			want: value.True,
		},
		"String 'foo' == 'foo'": {
			a:    value.String("foo"),
			b:    value.Ref(value.String("foo")),
			want: value.True,
		},
		"String 'a' == 'b'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("b")),
			want: value.False,
		},
		"String 'b' == 'a'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("a")),
			want: value.False,
		},
		"String 'aa' == 'a'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("a")),
			want: value.False,
		},
		"String 'a' == 'aa'": {
			a:    value.String("a"),
			b:    value.Ref(value.String("aa")),
			want: value.False,
		},
		"String 'aa' == 'b'": {
			a:    value.String("aa"),
			b:    value.Ref(value.String("b")),
			want: value.False,
		},
		"String 'b' == 'aa'": {
			a:    value.String("b"),
			b:    value.Ref(value.String("aa")),
			want: value.False,
		},
		"String 'abdf' == 'abcf'": {
			a:    value.String("abdf"),
			b:    value.Ref(value.String("abcf")),
			want: value.False,
		},
		"String 'abcf' == 'abdf'": {
			a:    value.String("abcf"),
			b:    value.Ref(value.String("abdf")),
			want: value.False,
		},
		"String 'ś' == 'ą'": {
			a:    value.String("ś"),
			b:    value.Ref(value.String("ą")),
			want: value.False,
		},
		"String 'ą' == 'ś'": {
			a:    value.String("ą"),
			b:    value.Ref(value.String("ś")),
			want: value.False,
		},

		"Char 'a' == `a`": {
			a:    value.String("a"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'a' == `b`": {
			a:    value.String("a"),
			b:    value.Char('b').ToValue(),
			want: value.False,
		},
		"Char 'b' == `a`": {
			a:    value.String("b"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'aa' == `a`": {
			a:    value.String("aa"),
			b:    value.Char('a').ToValue(),
			want: value.False,
		},
		"Char 'aa' == `b`": {
			a:    value.String("aa"),
			b:    value.Char('b').ToValue(),
			want: value.False,
		},
		"Char 'ś' == `ą`": {
			a:    value.String("ś"),
			b:    value.Char('ą').ToValue(),
			want: value.False,
		},
		"Char 'ą' == `ś`": {
			a:    value.String("ą"),
			b:    value.Char('ś').ToValue(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringCharIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		i    *value.StringCharIterator
		want string
	}{
		"empty": {
			i: value.NewStringCharIteratorWithByteOffset(
				"",
				0,
			),
			want: `Std::String::CharIterator{string: "", byte_offset: 0}`,
		},
		"not empty with offset 1": {
			i: value.NewStringCharIteratorWithByteOffset(
				value.String("ab"),
				1,
			),
			want: `Std::String::CharIterator{string: "ab", byte_offset: 1}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.i.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringCharIterator_Next(t *testing.T) {
	tests := map[string]struct {
		s     *value.StringCharIterator
		after *value.StringCharIterator
		want  value.Value
		err   value.Value
	}{
		"empty": {
			s: value.NewStringCharIteratorWithByteOffset(
				value.String(""),
				0,
			),
			after: value.NewStringCharIteratorWithByteOffset(
				value.String(""),
				0,
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
		"with two chars offset 0": {
			s: value.NewStringCharIteratorWithByteOffset(
				value.String("ab"),
				0,
			),
			after: value.NewStringCharIteratorWithByteOffset(
				value.String("ab"),
				1,
			),
			want: value.Char('a').ToValue(),
		},
		"with two-byte unicode chars offset 0": {
			s: value.NewStringCharIteratorWithByteOffset(
				value.String("śę"),
				0,
			),
			after: value.NewStringCharIteratorWithByteOffset(
				value.String("śę"),
				2,
			),
			want: value.Char('ś').ToValue(),
		},
		"with three-byte unicode chars offset 0": {
			s: value.NewStringCharIteratorWithByteOffset(
				value.String("≈∫"),
				0,
			),
			after: value.NewStringCharIteratorWithByteOffset(
				value.String("≈∫"),
				3,
			),
			want: value.Char('≈').ToValue(),
		},
		"with four-byte unicode chars offset 0": {
			s: value.NewStringCharIteratorWithByteOffset(
				value.String("😀🔥"),
				0,
			),
			after: value.NewStringCharIteratorWithByteOffset(
				value.String("😀🔥"),
				4,
			),
			want: value.Char('😀').ToValue(),
		},
		"with two chars offset 1": {
			s: value.NewStringCharIteratorWithByteOffset(
				value.String("ab"),
				1,
			),
			after: value.NewStringCharIteratorWithByteOffset(
				value.String("ab"),
				2,
			),
			want: value.Char('b').ToValue(),
		},
		"with two chars offset 2": {
			s: value.NewStringCharIteratorWithByteOffset(
				value.String("ab"),
				2,
			),
			after: value.NewStringCharIteratorWithByteOffset(
				value.String("ab"),
				2,
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.s.Next()
			if diff := cmp.Diff(tc.err, err); diff != "" {
				t.Fatalf(diff)
			}
			if !tc.err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.after, tc.s); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringByteIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		i    *value.StringByteIterator
		want string
	}{
		"empty": {
			i: value.NewStringByteIteratorWithByteOffset(
				"",
				0,
			),
			want: `Std::String::ByteIterator{string: "", byte_offset: 0}`,
		},
		"not empty with offset 1": {
			i: value.NewStringByteIteratorWithByteOffset(
				value.String("ab"),
				1,
			),
			want: `Std::String::ByteIterator{string: "ab", byte_offset: 1}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.i.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringByteIterator_Next(t *testing.T) {
	tests := map[string]struct {
		s     *value.StringByteIterator
		after *value.StringByteIterator
		want  value.Value
		err   value.Value
	}{
		"empty": {
			s: value.NewStringByteIteratorWithByteOffset(
				value.String(""),
				0,
			),
			after: value.NewStringByteIteratorWithByteOffset(
				value.String(""),
				0,
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
		"with two chars offset 0": {
			s: value.NewStringByteIteratorWithByteOffset(
				value.String("ab"),
				0,
			),
			after: value.NewStringByteIteratorWithByteOffset(
				value.String("ab"),
				1,
			),
			want: value.UInt8('a').ToValue(),
		},
		"with two-byte unicode chars offset 0": {
			s: value.NewStringByteIteratorWithByteOffset(
				value.String("śę"),
				0,
			),
			after: value.NewStringByteIteratorWithByteOffset(
				value.String("śę"),
				1,
			),
			want: value.UInt8('\xc5').ToValue(),
		},
		"with three-byte unicode chars offset 1": {
			s: value.NewStringByteIteratorWithByteOffset(
				value.String("≈∫"),
				1,
			),
			after: value.NewStringByteIteratorWithByteOffset(
				value.String("≈∫"),
				2,
			),
			want: value.UInt8('\x89').ToValue(),
		},
		"with four-byte unicode chars offset 3": {
			s: value.NewStringByteIteratorWithByteOffset(
				value.String("😀🔥"),
				3,
			),
			after: value.NewStringByteIteratorWithByteOffset(
				value.String("😀🔥"),
				4,
			),
			want: value.UInt8('\x80').ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.s.Next()
			if diff := cmp.Diff(tc.err, err); diff != "" {
				t.Fatalf(diff)
			}
			if !tc.err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.after, tc.s); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
