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

func TestString_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    String
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt and return an error": {
			a:   String("a"),
			b:   SmallInt(2),
			err: NewError(TypeErrorClass, "`Std::SmallInt` can't be coerced into `Std::String`"),
		},
		"Float and return an error": {
			a:   String("a"),
			b:   Float(5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::String`"),
		},
		"BigFloat and return an error": {
			a:   String("a"),
			b:   NewBigFloat(5),
			err: NewError(TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::String`"),
		},
		"Int64 and return an error": {
			a:   String("a"),
			b:   Int64(2),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::String`"),
		},
		"Int32 and return an error": {
			a:   String("a"),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::String`"),
		},
		"Int16 and return an error": {
			a:   String("a"),
			b:   Int16(2),
			err: NewError(TypeErrorClass, "`Std::Int16` can't be coerced into `Std::String`"),
		},
		"Int8 and return an error": {
			a:   String("a"),
			b:   Int8(2),
			err: NewError(TypeErrorClass, "`Std::Int8` can't be coerced into `Std::String`"),
		},
		"UInt64 and return an error": {
			a:   String("a"),
			b:   UInt64(2),
			err: NewError(TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::String`"),
		},
		"UInt32 and return an error": {
			a:   String("a"),
			b:   UInt32(2),
			err: NewError(TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::String`"),
		},
		"UInt16 and return an error": {
			a:   String("a"),
			b:   UInt16(2),
			err: NewError(TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::String`"),
		},
		"UInt8 and return an error": {
			a:   String("a"),
			b:   UInt8(2),
			err: NewError(TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::String`"),
		},
		"Float64 and return an error": {
			a:   String("a"),
			b:   Float64(5),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::String`"),
		},
		"Float32 and return an error": {
			a:   String("a"),
			b:   Float32(5),
			err: NewError(TypeErrorClass, "`Std::Float32` can't be coerced into `Std::String`"),
		},

		"String 'a' > 'a'": {
			a:    String("a"),
			b:    String("a"),
			want: False,
		},
		"String 'a' > 'b'": {
			a:    String("a"),
			b:    String("b"),
			want: False,
		},
		"String 'b' > 'a'": {
			a:    String("b"),
			b:    String("a"),
			want: True,
		},
		"String 'aa' > 'a'": {
			a:    String("aa"),
			b:    String("a"),
			want: True,
		},
		"String 'a' > 'aa'": {
			a:    String("a"),
			b:    String("aa"),
			want: False,
		},
		"String 'aa' > 'b'": {
			a:    String("aa"),
			b:    String("b"),
			want: False,
		},
		"String 'b' > 'aa'": {
			a:    String("b"),
			b:    String("aa"),
			want: True,
		},
		"String 'abdf' > 'abcf'": {
			a:    String("abdf"),
			b:    String("abcf"),
			want: True,
		},
		"String 'abcf' > 'abdf'": {
			a:    String("abcf"),
			b:    String("abdf"),
			want: False,
		},
		"String 'ś' > 'ą'": {
			a:    String("ś"),
			b:    String("ą"),
			want: True,
		},
		"String 'ą' > 'ś'": {
			a:    String("ą"),
			b:    String("ś"),
			want: False,
		},

		"Char 'a' > c'a'": {
			a:    String("a"),
			b:    Char('a'),
			want: False,
		},
		"Char 'a' > c'b'": {
			a:    String("a"),
			b:    Char('b'),
			want: False,
		},
		"Char 'b' > c'a'": {
			a:    String("b"),
			b:    Char('a'),
			want: True,
		},
		"Char 'aa' > c'a'": {
			a:    String("aa"),
			b:    Char('a'),
			want: True,
		},
		"Char 'aa' > c'b'": {
			a:    String("aa"),
			b:    Char('b'),
			want: False,
		},
		"Char 'ś' > c'ą'": {
			a:    String("ś"),
			b:    Char('ą'),
			want: True,
		},
		"Char 'ą' > c'ś'": {
			a:    String("ą"),
			b:    Char('ś'),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThan(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
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
		a    String
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt and return an error": {
			a:   String("a"),
			b:   SmallInt(2),
			err: NewError(TypeErrorClass, "`Std::SmallInt` can't be coerced into `Std::String`"),
		},
		"Float and return an error": {
			a:   String("a"),
			b:   Float(5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::String`"),
		},
		"BigFloat and return an error": {
			a:   String("a"),
			b:   NewBigFloat(5),
			err: NewError(TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::String`"),
		},
		"Int64 and return an error": {
			a:   String("a"),
			b:   Int64(2),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::String`"),
		},
		"Int32 and return an error": {
			a:   String("a"),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::String`"),
		},
		"Int16 and return an error": {
			a:   String("a"),
			b:   Int16(2),
			err: NewError(TypeErrorClass, "`Std::Int16` can't be coerced into `Std::String`"),
		},
		"Int8 and return an error": {
			a:   String("a"),
			b:   Int8(2),
			err: NewError(TypeErrorClass, "`Std::Int8` can't be coerced into `Std::String`"),
		},
		"UInt64 and return an error": {
			a:   String("a"),
			b:   UInt64(2),
			err: NewError(TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::String`"),
		},
		"UInt32 and return an error": {
			a:   String("a"),
			b:   UInt32(2),
			err: NewError(TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::String`"),
		},
		"UInt16 and return an error": {
			a:   String("a"),
			b:   UInt16(2),
			err: NewError(TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::String`"),
		},
		"UInt8 and return an error": {
			a:   String("a"),
			b:   UInt8(2),
			err: NewError(TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::String`"),
		},
		"Float64 and return an error": {
			a:   String("a"),
			b:   Float64(5),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::String`"),
		},
		"Float32 and return an error": {
			a:   String("a"),
			b:   Float32(5),
			err: NewError(TypeErrorClass, "`Std::Float32` can't be coerced into `Std::String`"),
		},

		"String 'a' >= 'a'": {
			a:    String("a"),
			b:    String("a"),
			want: True,
		},
		"String 'foo' >= 'foo'": {
			a:    String("foo"),
			b:    String("foo"),
			want: True,
		},
		"String 'a' >= 'b'": {
			a:    String("a"),
			b:    String("b"),
			want: False,
		},
		"String 'b' >= 'a'": {
			a:    String("b"),
			b:    String("a"),
			want: True,
		},
		"String 'aa' >= 'a'": {
			a:    String("aa"),
			b:    String("a"),
			want: True,
		},
		"String 'a' >= 'aa'": {
			a:    String("a"),
			b:    String("aa"),
			want: False,
		},
		"String 'aa' >= 'b'": {
			a:    String("aa"),
			b:    String("b"),
			want: False,
		},
		"String 'b' >= 'aa'": {
			a:    String("b"),
			b:    String("aa"),
			want: True,
		},
		"String 'abdf' >= 'abcf'": {
			a:    String("abdf"),
			b:    String("abcf"),
			want: True,
		},
		"String 'abcf' >= 'abdf'": {
			a:    String("abcf"),
			b:    String("abdf"),
			want: False,
		},
		"String 'ś' >= 'ą'": {
			a:    String("ś"),
			b:    String("ą"),
			want: True,
		},
		"String 'ą' >= 'ś'": {
			a:    String("ą"),
			b:    String("ś"),
			want: False,
		},

		"Char 'a' >= c'a'": {
			a:    String("a"),
			b:    Char('a'),
			want: True,
		},
		"Char 'a' >= c'b'": {
			a:    String("a"),
			b:    Char('b'),
			want: False,
		},
		"Char 'b' >= c'a'": {
			a:    String("b"),
			b:    Char('a'),
			want: True,
		},
		"Char 'aa' >= c'a'": {
			a:    String("aa"),
			b:    Char('a'),
			want: True,
		},
		"Char 'aa' >= c'b'": {
			a:    String("aa"),
			b:    Char('b'),
			want: False,
		},
		"Char 'ś' >= c'ą'": {
			a:    String("ś"),
			b:    Char('ą'),
			want: True,
		},
		"Char 'ą' >= c'ś'": {
			a:    String("ą"),
			b:    Char('ś'),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqual(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
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
		a    String
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt and return an error": {
			a:   String("a"),
			b:   SmallInt(2),
			err: NewError(TypeErrorClass, "`Std::SmallInt` can't be coerced into `Std::String`"),
		},
		"Float and return an error": {
			a:   String("a"),
			b:   Float(5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::String`"),
		},
		"BigFloat and return an error": {
			a:   String("a"),
			b:   NewBigFloat(5),
			err: NewError(TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::String`"),
		},
		"Int64 and return an error": {
			a:   String("a"),
			b:   Int64(2),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::String`"),
		},
		"Int32 and return an error": {
			a:   String("a"),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::String`"),
		},
		"Int16 and return an error": {
			a:   String("a"),
			b:   Int16(2),
			err: NewError(TypeErrorClass, "`Std::Int16` can't be coerced into `Std::String`"),
		},
		"Int8 and return an error": {
			a:   String("a"),
			b:   Int8(2),
			err: NewError(TypeErrorClass, "`Std::Int8` can't be coerced into `Std::String`"),
		},
		"UInt64 and return an error": {
			a:   String("a"),
			b:   UInt64(2),
			err: NewError(TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::String`"),
		},
		"UInt32 and return an error": {
			a:   String("a"),
			b:   UInt32(2),
			err: NewError(TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::String`"),
		},
		"UInt16 and return an error": {
			a:   String("a"),
			b:   UInt16(2),
			err: NewError(TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::String`"),
		},
		"UInt8 and return an error": {
			a:   String("a"),
			b:   UInt8(2),
			err: NewError(TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::String`"),
		},
		"Float64 and return an error": {
			a:   String("a"),
			b:   Float64(5),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::String`"),
		},
		"Float32 and return an error": {
			a:   String("a"),
			b:   Float32(5),
			err: NewError(TypeErrorClass, "`Std::Float32` can't be coerced into `Std::String`"),
		},

		"String 'a' < 'a'": {
			a:    String("a"),
			b:    String("a"),
			want: False,
		},
		"String 'foo' < 'foo'": {
			a:    String("foo"),
			b:    String("foo"),
			want: False,
		},
		"String 'a' < 'b'": {
			a:    String("a"),
			b:    String("b"),
			want: True,
		},
		"String 'b' < 'a'": {
			a:    String("b"),
			b:    String("a"),
			want: False,
		},
		"String 'aa' < 'a'": {
			a:    String("aa"),
			b:    String("a"),
			want: False,
		},
		"String 'a' < 'aa'": {
			a:    String("a"),
			b:    String("aa"),
			want: True,
		},
		"String 'aa' < 'b'": {
			a:    String("aa"),
			b:    String("b"),
			want: True,
		},
		"String 'b' < 'aa'": {
			a:    String("b"),
			b:    String("aa"),
			want: False,
		},
		"String 'abdf' < 'abcf'": {
			a:    String("abdf"),
			b:    String("abcf"),
			want: False,
		},
		"String 'abcf' < 'abdf'": {
			a:    String("abcf"),
			b:    String("abdf"),
			want: True,
		},
		"String 'ś' < 'ą'": {
			a:    String("ś"),
			b:    String("ą"),
			want: False,
		},
		"String 'ą' < 'ś'": {
			a:    String("ą"),
			b:    String("ś"),
			want: True,
		},

		"Char 'a' < c'a'": {
			a:    String("a"),
			b:    Char('a'),
			want: False,
		},
		"Char 'a' < c'b'": {
			a:    String("a"),
			b:    Char('b'),
			want: True,
		},
		"Char 'b' < c'a'": {
			a:    String("b"),
			b:    Char('a'),
			want: False,
		},
		"Char 'aa' < c'a'": {
			a:    String("aa"),
			b:    Char('a'),
			want: False,
		},
		"Char 'aa' < c'b'": {
			a:    String("aa"),
			b:    Char('b'),
			want: True,
		},
		"Char 'ś' < c'ą'": {
			a:    String("ś"),
			b:    Char('ą'),
			want: False,
		},
		"Char 'ą' < c'ś'": {
			a:    String("ą"),
			b:    Char('ś'),
			want: True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThan(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
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
		a    String
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt and return an error": {
			a:   String("a"),
			b:   SmallInt(2),
			err: NewError(TypeErrorClass, "`Std::SmallInt` can't be coerced into `Std::String`"),
		},
		"Float and return an error": {
			a:   String("a"),
			b:   Float(5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::String`"),
		},
		"BigFloat and return an error": {
			a:   String("a"),
			b:   NewBigFloat(5),
			err: NewError(TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::String`"),
		},
		"Int64 and return an error": {
			a:   String("a"),
			b:   Int64(2),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::String`"),
		},
		"Int32 and return an error": {
			a:   String("a"),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::String`"),
		},
		"Int16 and return an error": {
			a:   String("a"),
			b:   Int16(2),
			err: NewError(TypeErrorClass, "`Std::Int16` can't be coerced into `Std::String`"),
		},
		"Int8 and return an error": {
			a:   String("a"),
			b:   Int8(2),
			err: NewError(TypeErrorClass, "`Std::Int8` can't be coerced into `Std::String`"),
		},
		"UInt64 and return an error": {
			a:   String("a"),
			b:   UInt64(2),
			err: NewError(TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::String`"),
		},
		"UInt32 and return an error": {
			a:   String("a"),
			b:   UInt32(2),
			err: NewError(TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::String`"),
		},
		"UInt16 and return an error": {
			a:   String("a"),
			b:   UInt16(2),
			err: NewError(TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::String`"),
		},
		"UInt8 and return an error": {
			a:   String("a"),
			b:   UInt8(2),
			err: NewError(TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::String`"),
		},
		"Float64 and return an error": {
			a:   String("a"),
			b:   Float64(5),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::String`"),
		},
		"Float32 and return an error": {
			a:   String("a"),
			b:   Float32(5),
			err: NewError(TypeErrorClass, "`Std::Float32` can't be coerced into `Std::String`"),
		},

		"String 'a' <= 'a'": {
			a:    String("a"),
			b:    String("a"),
			want: True,
		},
		"String 'foo' <= 'foo'": {
			a:    String("foo"),
			b:    String("foo"),
			want: True,
		},
		"String 'a' <= 'b'": {
			a:    String("a"),
			b:    String("b"),
			want: True,
		},
		"String 'b' <= 'a'": {
			a:    String("b"),
			b:    String("a"),
			want: False,
		},
		"String 'aa' <= 'a'": {
			a:    String("aa"),
			b:    String("a"),
			want: False,
		},
		"String 'a' <= 'aa'": {
			a:    String("a"),
			b:    String("aa"),
			want: True,
		},
		"String 'aa' <= 'b'": {
			a:    String("aa"),
			b:    String("b"),
			want: True,
		},
		"String 'b' <= 'aa'": {
			a:    String("b"),
			b:    String("aa"),
			want: False,
		},
		"String 'abdf' <= 'abcf'": {
			a:    String("abdf"),
			b:    String("abcf"),
			want: False,
		},
		"String 'abcf' <= 'abdf'": {
			a:    String("abcf"),
			b:    String("abdf"),
			want: True,
		},
		"String 'ś' <= 'ą'": {
			a:    String("ś"),
			b:    String("ą"),
			want: False,
		},
		"String 'ą' <= 'ś'": {
			a:    String("ą"),
			b:    String("ś"),
			want: True,
		},

		"Char 'a' <= c'a'": {
			a:    String("a"),
			b:    Char('a'),
			want: True,
		},
		"Char 'a' <= c'b'": {
			a:    String("a"),
			b:    Char('b'),
			want: True,
		},
		"Char 'b' <= c'a'": {
			a:    String("b"),
			b:    Char('a'),
			want: False,
		},
		"Char 'aa' <= c'a'": {
			a:    String("aa"),
			b:    Char('a'),
			want: False,
		},
		"Char 'aa' <= c'b'": {
			a:    String("aa"),
			b:    Char('b'),
			want: True,
		},
		"Char 'ś' <= c'ą'": {
			a:    String("ś"),
			b:    Char('ą'),
			want: False,
		},
		"Char 'ą' <= c'ś'": {
			a:    String("ą"),
			b:    Char('ś'),
			want: True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqual(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
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

func TestString_Equal(t *testing.T) {
	tests := map[string]struct {
		a    String
		b    Value
		want Value
	}{
		"SmallInt '2' == 2": {
			a:    String("2"),
			b:    SmallInt(2),
			want: False,
		},
		"Float '5.2' == 5.2": {
			a:    String("5.2"),
			b:    Float(5.2),
			want: False,
		},
		"BigFloat '4.5' == 4.5bf": {
			a:    String("4.5"),
			b:    NewBigFloat(4.5),
			want: False,
		},
		"Int64 '5' == 5i64": {
			a:    String("5"),
			b:    Int64(5),
			want: False,
		},
		"Int32 '2' == 2i32": {
			a:    String("2"),
			b:    Int32(2),
			want: False,
		},
		"Int16 '25' == 25i16": {
			a:    String("25"),
			b:    Int16(25),
			want: False,
		},
		"Int8 '8' == 8i8": {
			a:    String("8"),
			b:    Int8(8),
			want: False,
		},
		"UInt64 '31' == 31u64": {
			a:    String("31"),
			b:    UInt64(31),
			want: False,
		},
		"UInt32 '9' == 9u32": {
			a:    String("9"),
			b:    UInt32(9),
			want: False,
		},
		"UInt16 '74' == 74u16": {
			a:    String("74"),
			b:    UInt16(74),
			want: False,
		},
		"UInt8 '12' == 12u8": {
			a:    String("12"),
			b:    Int8(12),
			want: False,
		},
		"Float64 '49.2' == 49.2f64": {
			a:    String("49.2"),
			b:    Float64(49.2),
			want: False,
		},
		"Float32 '57.9' == 57.9f32": {
			a:    String("57.9"),
			b:    Float32(57.9),
			want: False,
		},

		"String 'a' == 'a'": {
			a:    String("a"),
			b:    String("a"),
			want: True,
		},
		"String 'foo' == 'foo'": {
			a:    String("foo"),
			b:    String("foo"),
			want: True,
		},
		"String 'a' == 'b'": {
			a:    String("a"),
			b:    String("b"),
			want: False,
		},
		"String 'b' == 'a'": {
			a:    String("b"),
			b:    String("a"),
			want: False,
		},
		"String 'aa' == 'a'": {
			a:    String("aa"),
			b:    String("a"),
			want: False,
		},
		"String 'a' == 'aa'": {
			a:    String("a"),
			b:    String("aa"),
			want: False,
		},
		"String 'aa' == 'b'": {
			a:    String("aa"),
			b:    String("b"),
			want: False,
		},
		"String 'b' == 'aa'": {
			a:    String("b"),
			b:    String("aa"),
			want: False,
		},
		"String 'abdf' == 'abcf'": {
			a:    String("abdf"),
			b:    String("abcf"),
			want: False,
		},
		"String 'abcf' == 'abdf'": {
			a:    String("abcf"),
			b:    String("abdf"),
			want: False,
		},
		"String 'ś' == 'ą'": {
			a:    String("ś"),
			b:    String("ą"),
			want: False,
		},
		"String 'ą' == 'ś'": {
			a:    String("ą"),
			b:    String("ś"),
			want: False,
		},

		"Char 'a' == c'a'": {
			a:    String("a"),
			b:    Char('a'),
			want: True,
		},
		"Char 'a' == c'b'": {
			a:    String("a"),
			b:    Char('b'),
			want: False,
		},
		"Char 'b' == c'a'": {
			a:    String("b"),
			b:    Char('a'),
			want: False,
		},
		"Char 'aa' == c'a'": {
			a:    String("aa"),
			b:    Char('a'),
			want: False,
		},
		"Char 'aa' == c'b'": {
			a:    String("aa"),
			b:    Char('b'),
			want: False,
		},
		"Char 'ś' == c'ą'": {
			a:    String("ś"),
			b:    Char('ą'),
			want: False,
		},
		"Char 'ą' == c'ś'": {
			a:    String("ą"),
			b:    Char('ś'),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestString_StrictEqual(t *testing.T) {
	tests := map[string]struct {
		a    String
		b    Value
		want Value
	}{
		"SmallInt '2' === 2": {
			a:    String("2"),
			b:    SmallInt(2),
			want: False,
		},
		"Float '5.2' === 5.2": {
			a:    String("5.2"),
			b:    Float(5.2),
			want: False,
		},
		"BigFloat '4.5' === 4.5bf": {
			a:    String("4.5"),
			b:    NewBigFloat(4.5),
			want: False,
		},
		"Int64 '5' === 5i64": {
			a:    String("5"),
			b:    Int64(5),
			want: False,
		},
		"Int32 '2' === 2i32": {
			a:    String("2"),
			b:    Int32(2),
			want: False,
		},
		"Int16 '25' === 25i16": {
			a:    String("25"),
			b:    Int16(25),
			want: False,
		},
		"Int8 '8' === 8i8": {
			a:    String("8"),
			b:    Int8(8),
			want: False,
		},
		"UInt64 '31' === 31u64": {
			a:    String("31"),
			b:    UInt64(31),
			want: False,
		},
		"UInt32 '9' === 9u32": {
			a:    String("9"),
			b:    UInt32(9),
			want: False,
		},
		"UInt16 '74' === 74u16": {
			a:    String("74"),
			b:    UInt16(74),
			want: False,
		},
		"UInt8 '12' === 12u8": {
			a:    String("12"),
			b:    Int8(12),
			want: False,
		},
		"Float64 '49.2' === 49.2f64": {
			a:    String("49.2"),
			b:    Float64(49.2),
			want: False,
		},
		"Float32 '57.9' === 57.9f32": {
			a:    String("57.9"),
			b:    Float32(57.9),
			want: False,
		},

		"String 'a' === 'a'": {
			a:    String("a"),
			b:    String("a"),
			want: True,
		},
		"String 'foo' === 'foo'": {
			a:    String("foo"),
			b:    String("foo"),
			want: True,
		},
		"String 'a' === 'b'": {
			a:    String("a"),
			b:    String("b"),
			want: False,
		},
		"String 'b' === 'a'": {
			a:    String("b"),
			b:    String("a"),
			want: False,
		},
		"String 'aa' === 'a'": {
			a:    String("aa"),
			b:    String("a"),
			want: False,
		},
		"String 'a' === 'aa'": {
			a:    String("a"),
			b:    String("aa"),
			want: False,
		},
		"String 'aa' === 'b'": {
			a:    String("aa"),
			b:    String("b"),
			want: False,
		},
		"String 'b' === 'aa'": {
			a:    String("b"),
			b:    String("aa"),
			want: False,
		},
		"String 'abdf' === 'abcf'": {
			a:    String("abdf"),
			b:    String("abcf"),
			want: False,
		},
		"String 'abcf' === 'abdf'": {
			a:    String("abcf"),
			b:    String("abdf"),
			want: False,
		},
		"String 'ś' === 'ą'": {
			a:    String("ś"),
			b:    String("ą"),
			want: False,
		},
		"String 'ą' === 'ś'": {
			a:    String("ą"),
			b:    String("ś"),
			want: False,
		},

		"Char 'a' === c'a'": {
			a:    String("a"),
			b:    Char('a'),
			want: False,
		},
		"Char 'a' === c'b'": {
			a:    String("a"),
			b:    Char('b'),
			want: False,
		},
		"Char 'b' === c'a'": {
			a:    String("b"),
			b:    Char('a'),
			want: False,
		},
		"Char 'aa' === c'a'": {
			a:    String("aa"),
			b:    Char('a'),
			want: False,
		},
		"Char 'aa' === c'b'": {
			a:    String("aa"),
			b:    Char('b'),
			want: False,
		},
		"Char 'ś' === c'ą'": {
			a:    String("ś"),
			b:    Char('ą'),
			want: False,
		},
		"Char 'ą' === c'ś'": {
			a:    String("ą"),
			b:    Char('ś'),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.StrictEqual(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}
