package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestCharConcat(t *testing.T) {
	tests := map[string]struct {
		left  value.Char
		right value.Value
		want  value.String
		err   *value.Error
	}{
		"Char + String => String": {
			left:  value.Char('f'),
			right: value.String("oo"),
			want:  value.String("foo"),
		},
		"Char + Char => String": {
			left:  value.Char('a'),
			right: value.Char('b'),
			want:  value.String("ab"),
		},
		"Char + Int => TypeError": {
			left:  value.Char('a'),
			right: value.Int8(5),
			err:   value.NewError(value.TypeErrorClass, `can't concat 5i8 with char c"a"`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Concat(tc.right)
			opts := comparer.Comparer
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestCharRepeat(t *testing.T) {
	tests := map[string]struct {
		left  value.Char
		right value.Value
		want  value.String
		err   *value.Error
	}{
		"Char * SmallInt => String": {
			left:  value.Char('a'),
			right: value.SmallInt(3),
			want:  value.String("aaa"),
		},
		"Char * BigInt => OutOfRangeError": {
			left:  value.Char('a'),
			right: value.NewBigInt(3),
			err:   value.NewError(value.OutOfRangeErrorClass, `repeat count is too large 3`),
		},
		"Char * Int8 => TypeError": {
			left:  value.Char('a'),
			right: value.Int8(3),
			err:   value.NewError(value.TypeErrorClass, `can't repeat a char using 3i8`),
		},
		"String * String => TypeError": {
			left:  value.Char('a'),
			right: value.String("bar"),
			err:   value.NewError(value.TypeErrorClass, `can't repeat a char using "bar"`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Repeat(tc.right)
			opts := comparer.Comparer
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestCharInspect(t *testing.T) {
	tests := map[string]struct {
		c    value.Char
		want string
	}{
		"ascii letter": {
			c:    'd',
			want: `c"d"`,
		},
		"utf-8 character": {
			c:    'ś',
			want: `c"ś"`,
		},
		"newline": {
			c:    '\n',
			want: `c"\n"`,
		},
		"double quote": {
			c:    '"',
			want: `c"\""`,
		},
		"backslash": {
			c:    '\\',
			want: `c"\\"`,
		},
		"hex byte": {
			c:    '\x02',
			want: `c"\x02"`,
		},
		"unicode codepoint": {
			c:    '\U0010FFFF',
			want: `c"\U0010FFFF"`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.c.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestChar_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
		err  *value.Error
	}{
		"SmallInt and return an error": {
			a:   value.Char('a'),
			b:   value.SmallInt(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int` can't be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   value.Char('a'),
			b:   value.Float(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` can't be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   value.Char('a'),
			b:   value.NewBigFloat(5),
			err: value.NewError(value.TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   value.Char('a'),
			b:   value.Int64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   value.Char('a'),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   value.Char('a'),
			b:   value.Int16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int16` can't be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   value.Char('a'),
			b:   value.Int8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int8` can't be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   value.Char('a'),
			b:   value.Float64(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   value.Char('a'),
			b:   value.Float32(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float32` can't be coerced into `Std::Char`"),
		},

		"Char c'a' > c'a'": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char c'a' > c'b'": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.False,
		},
		"Char c'b' > c'a'": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char c'ś' > c'ą'": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.True,
		},
		"Char c'ą' > c'ś'": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.False,
		},

		"String c'a' > 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.False,
		},
		"String c'a' > 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.False,
		},
		"String c'b' > 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.True,
		},
		"String c'a' > 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String c'b' > 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.True,
		},
		"String c'ś' > 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.True,
		},
		"String c'ą' > 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThan(tc.b)
			opts := comparer.Comparer
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

func TestChar_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
		err  *value.Error
	}{
		"SmallInt and return an error": {
			a:   value.Char('a'),
			b:   value.SmallInt(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int` can't be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   value.Char('a'),
			b:   value.Float(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` can't be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   value.Char('a'),
			b:   value.NewBigFloat(5),
			err: value.NewError(value.TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   value.Char('a'),
			b:   value.Int64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   value.Char('a'),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   value.Char('a'),
			b:   value.Int16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int16` can't be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   value.Char('a'),
			b:   value.Int8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int8` can't be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   value.Char('a'),
			b:   value.Float64(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   value.Char('a'),
			b:   value.Float32(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float32` can't be coerced into `Std::Char`"),
		},

		"Char c'a' >= c'a'": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char c'a' >= c'b'": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.False,
		},
		"Char c'b' >= c'a'": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char c'ś' >= c'ą'": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.True,
		},
		"Char c'ą' >= c'ś'": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.False,
		},

		"String c'a' >= 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.True,
		},
		"String c'a' >= 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.False,
		},
		"String c'b' >= 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.True,
		},
		"String c'a' >= 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String c'b' >= 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.True,
		},
		"String c'ś' >= 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.True,
		},
		"String c'ą' >= 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqual(tc.b)
			opts := comparer.Comparer
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

func TestChar_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
		err  *value.Error
	}{
		"SmallInt and return an error": {
			a:   value.Char('a'),
			b:   value.SmallInt(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int` can't be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   value.Char('a'),
			b:   value.Float(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` can't be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   value.Char('a'),
			b:   value.NewBigFloat(5),
			err: value.NewError(value.TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   value.Char('a'),
			b:   value.Int64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   value.Char('a'),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   value.Char('a'),
			b:   value.Int16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int16` can't be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   value.Char('a'),
			b:   value.Int8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int8` can't be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   value.Char('a'),
			b:   value.Float64(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   value.Char('a'),
			b:   value.Float32(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float32` can't be coerced into `Std::Char`"),
		},

		"Char c'a' < c'a'": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char c'a' < c'b'": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.True,
		},
		"Char c'b' < c'a'": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char c'ś' < c'ą'": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.False,
		},
		"Char c'ą' < c'ś'": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.True,
		},

		"String c'a' < 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.False,
		},
		"String c'a' < 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.True,
		},
		"String c'b' < 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.False,
		},
		"String c'a' < 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.True,
		},
		"String c'b' < 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String c'ś' < 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.False,
		},
		"String c'ą' < 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThan(tc.b)
			opts := comparer.Comparer
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

func TestChar_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
		err  *value.Error
	}{
		"SmallInt and return an error": {
			a:   value.Char('a'),
			b:   value.SmallInt(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int` can't be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   value.Char('a'),
			b:   value.Float(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` can't be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   value.Char('a'),
			b:   value.NewBigFloat(5),
			err: value.NewError(value.TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   value.Char('a'),
			b:   value.Int64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   value.Char('a'),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   value.Char('a'),
			b:   value.Int16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int16` can't be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   value.Char('a'),
			b:   value.Int8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int8` can't be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   value.Char('a'),
			b:   value.Float64(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   value.Char('a'),
			b:   value.Float32(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float32` can't be coerced into `Std::Char`"),
		},

		"Char c'a' <= c'a'": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char c'a' <= c'b'": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.True,
		},
		"Char c'b' <= c'a'": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char c'ś' <= c'ą'": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.False,
		},
		"Char c'ą' <= c'ś'": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.True,
		},

		"String c'a' <= 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.True,
		},
		"String c'a' <= 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.True,
		},
		"String c'b' <= 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.False,
		},
		"String c'a' <= 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.True,
		},
		"String c'b' <= 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String c'ś' <= 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.False,
		},
		"String c'ą' <= 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqual(tc.b)
			opts := comparer.Comparer
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

func TestChar_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
	}{
		"SmallInt c'2' == 2": {
			a:    value.Char('2'),
			b:    value.SmallInt(2),
			want: value.False,
		},
		"Float c'5' == 5.0": {
			a:    value.Char('5'),
			b:    value.Float(5),
			want: value.False,
		},
		"BigFloat c'4' == 4bf": {
			a:    value.Char('4'),
			b:    value.NewBigFloat(4),
			want: value.False,
		},
		"Int64 c'5' == 5i64": {
			a:    value.Char('5'),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int32 c'2' == 2i32": {
			a:    value.Char('2'),
			b:    value.Int32(2),
			want: value.False,
		},
		"Int16 c'8' == 8i16": {
			a:    value.Char('8'),
			b:    value.Int16(8),
			want: value.False,
		},
		"Int8 c'8' == 8i8": {
			a:    value.Char('8'),
			b:    value.Int8(8),
			want: value.False,
		},
		"UInt64 c'3' == 3u64": {
			a:    value.Char('3'),
			b:    value.UInt64(3),
			want: value.False,
		},
		"UInt32 c'9' == 9u32": {
			a:    value.Char('9'),
			b:    value.UInt32(9),
			want: value.False,
		},
		"UInt16 c'7' == 7u16": {
			a:    value.Char('7'),
			b:    value.UInt16(7),
			want: value.False,
		},
		"UInt8 '1' == 1u8": {
			a:    value.Char('1'),
			b:    value.Int8(12),
			want: value.False,
		},
		"Float64 c'0' == 0f64": {
			a:    value.Char('0'),
			b:    value.Float64(0),
			want: value.False,
		},
		"Float32 c'5' == 5f32": {
			a:    value.Char('5'),
			b:    value.Float32(5),
			want: value.False,
		},

		"String c'a' == 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.True,
		},
		"String c'a' == 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.False,
		},
		"String c'b' == 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.False,
		},
		"String c'a' == 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String c'b' == 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String c'ś' == 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.False,
		},
		"String c'ś' == 'ś'": {
			a:    value.Char('ś'),
			b:    value.String("ś"),
			want: value.True,
		},
		"String c'ą' == 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
			want: value.False,
		},

		"Char c'a' == c'a'": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char c'a' == c'b'": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.False,
		},
		"Char c'b' == c'a'": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char c'ś' == c'ą'": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.False,
		},
		"Char c'ś' == c'ś'": {
			a:    value.Char('ś'),
			b:    value.Char('ś'),
			want: value.True,
		},
		"Char c'ą' == c'ś'": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			opts := comparer.Comparer
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestChar_StrictEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
	}{
		"SmallInt c'2' === 2": {
			a:    value.Char('2'),
			b:    value.SmallInt(2),
			want: value.False,
		},
		"Float c'5' === 5.0": {
			a:    value.Char('5'),
			b:    value.Float(5),
			want: value.False,
		},
		"BigFloat c'4' === 4bf": {
			a:    value.Char('4'),
			b:    value.NewBigFloat(4),
			want: value.False,
		},
		"Int64 c'5' === 5i64": {
			a:    value.Char('5'),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int32 c'2' === 2i32": {
			a:    value.Char('2'),
			b:    value.Int32(2),
			want: value.False,
		},
		"Int16 c'8' === 8i16": {
			a:    value.Char('8'),
			b:    value.Int16(8),
			want: value.False,
		},
		"Int8 c'8' === 8i8": {
			a:    value.Char('8'),
			b:    value.Int8(8),
			want: value.False,
		},
		"UInt64 c'3' === 3u64": {
			a:    value.Char('3'),
			b:    value.UInt64(3),
			want: value.False,
		},
		"UInt32 c'9' === 9u32": {
			a:    value.Char('9'),
			b:    value.UInt32(9),
			want: value.False,
		},
		"UInt16 c'7' === 7u16": {
			a:    value.Char('7'),
			b:    value.UInt16(7),
			want: value.False,
		},
		"UInt8 '1' === 1u8": {
			a:    value.Char('1'),
			b:    value.Int8(12),
			want: value.False,
		},
		"Float64 c'0' === 0f64": {
			a:    value.Char('0'),
			b:    value.Float64(0),
			want: value.False,
		},
		"Float32 c'5' === 5f32": {
			a:    value.Char('5'),
			b:    value.Float32(5),
			want: value.False,
		},

		"String c'a' === 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.False,
		},
		"String c'a' === 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.False,
		},
		"String c'b' === 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.False,
		},
		"String c'a' === 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String c'b' === 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String c'ś' === 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.False,
		},
		"String c'ś' === 'ś'": {
			a:    value.Char('ś'),
			b:    value.String("ś"),
			want: value.False,
		},
		"String c'ą' === 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
			want: value.False,
		},

		"Char c'a' === c'a'": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char c'a' === c'b'": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.False,
		},
		"Char c'b' === c'a'": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char c'ś' === c'ą'": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.False,
		},
		"Char c'ś' === c'ś'": {
			a:    value.Char('ś'),
			b:    value.Char('ś'),
			want: value.True,
		},
		"Char c'ą' === c'ś'": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.StrictEqual(tc.b)
			opts := comparer.Comparer
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}
