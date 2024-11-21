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
		err   value.Value
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
			err:   value.NewError(value.TypeErrorClass, "cannot concat 5i8 with char `a`"),
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

func TestCharRepeat(t *testing.T) {
	tests := map[string]struct {
		left  value.Char
		right value.Value
		want  value.String
		err   value.Value
	}{
		"Char * SmallInt => String": {
			left:  value.Char('a'),
			right: value.SmallInt(3),
			want:  value.String("aaa"),
		},
		"Char * 0 => String": {
			left:  value.Char('a'),
			right: value.SmallInt(0),
			want:  value.String(""),
		},
		"Char * -SmallInt => OutOfRangeError": {
			left:  value.Char('a'),
			right: value.SmallInt(-3),
			err:   value.NewError(value.OutOfRangeErrorClass, `repeat count cannot be negative: -3`),
		},
		"Char * BigInt => OutOfRangeError": {
			left:  value.Char('a'),
			right: value.NewBigInt(3),
			err:   value.NewError(value.OutOfRangeErrorClass, `repeat count is too large 3`),
		},
		"Char * Int8 => TypeError": {
			left:  value.Char('a'),
			right: value.Int8(3),
			err:   value.NewError(value.TypeErrorClass, `cannot repeat a char using 3i8`),
		},
		"String * String => TypeError": {
			left:  value.Char('a'),
			right: value.String("bar"),
			err:   value.NewError(value.TypeErrorClass, `cannot repeat a char using "bar"`),
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

func TestCharInspect(t *testing.T) {
	tests := map[string]struct {
		c    value.Char
		want string
	}{
		"ascii letter": {
			c:    'd',
			want: "`d`",
		},
		"utf-8 character": {
			c:    'ś',
			want: "`ś`",
		},
		"newline": {
			c:    '\n',
			want: "`\\n`",
		},
		"double quote": {
			c:    '"',
			want: "`\"`",
		},
		"backtick": {
			c:    '`',
			want: "`\\``",
		},
		"backslash": {
			c:    '\\',
			want: "`\\\\`",
		},
		"hex byte": {
			c:    '\x02',
			want: "`\\x02`",
		},
		"unicode codepoint": {
			c:    '\U0010FFFF',
			want: "`\\U0010FFFF`",
		},
		"small unicode codepoint": {
			c:    '\u200d',
			want: "`\\u200d`",
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

func TestChar_Compare(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt and return an error": {
			a:   value.Char('a'),
			b:   value.SmallInt(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   value.Char('a'),
			b:   value.Float(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   value.Char('a'),
			b:   value.NewBigFloat(5),
			err: value.NewError(value.TypeErrorClass, "`Std::BigFloat` cannot be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   value.Char('a'),
			b:   value.Int64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   value.Char('a'),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   value.Char('a'),
			b:   value.Int16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int16` cannot be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   value.Char('a'),
			b:   value.Int8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int8` cannot be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt64` cannot be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt32` cannot be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt16` cannot be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt8` cannot be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   value.Char('a'),
			b:   value.Float64(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   value.Char('a'),
			b:   value.Float32(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::Char`"),
		},

		"Char `a` <=> `a`": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.SmallInt(0),
		},
		"Char `a` <=> `b`": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.SmallInt(-1),
		},
		"Char `b` <=> `a`": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.SmallInt(1),
		},
		"Char `ś` <=> `ą`": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.SmallInt(1),
		},
		"Char `ą` <=> `ś`": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.SmallInt(-1),
		},

		"String `a` <=> 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.SmallInt(0),
		},
		"String `a` <=> 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.SmallInt(-1),
		},
		"String `b` <=> 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.SmallInt(1),
		},
		"String `a` <=> 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.SmallInt(-1),
		},
		"String `b` <=> 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.SmallInt(1),
		},
		"String `ś` <=> 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.SmallInt(1),
		},
		"String `ą` <=> 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
			want: value.SmallInt(-1),
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

func TestChar_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt and return an error": {
			a:   value.Char('a'),
			b:   value.SmallInt(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   value.Char('a'),
			b:   value.Float(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   value.Char('a'),
			b:   value.NewBigFloat(5),
			err: value.NewError(value.TypeErrorClass, "`Std::BigFloat` cannot be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   value.Char('a'),
			b:   value.Int64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   value.Char('a'),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   value.Char('a'),
			b:   value.Int16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int16` cannot be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   value.Char('a'),
			b:   value.Int8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int8` cannot be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt64` cannot be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt32` cannot be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt16` cannot be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt8` cannot be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   value.Char('a'),
			b:   value.Float64(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   value.Char('a'),
			b:   value.Float32(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::Char`"),
		},

		"Char `a` > `a`": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char `a` > `b`": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.False,
		},
		"Char `b` > `a`": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char `ś` > `ą`": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.True,
		},
		"Char `ą` > `ś`": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.False,
		},

		"String `a` > 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.False,
		},
		"String `a` > 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.False,
		},
		"String `b` > 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.True,
		},
		"String `a` > 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String `b` > 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.True,
		},
		"String `ś` > 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.True,
		},
		"String `ą` > 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
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

func TestChar_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt and return an error": {
			a:   value.Char('a'),
			b:   value.SmallInt(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   value.Char('a'),
			b:   value.Float(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   value.Char('a'),
			b:   value.NewBigFloat(5),
			err: value.NewError(value.TypeErrorClass, "`Std::BigFloat` cannot be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   value.Char('a'),
			b:   value.Int64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   value.Char('a'),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   value.Char('a'),
			b:   value.Int16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int16` cannot be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   value.Char('a'),
			b:   value.Int8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int8` cannot be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt64` cannot be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt32` cannot be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt16` cannot be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt8` cannot be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   value.Char('a'),
			b:   value.Float64(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   value.Char('a'),
			b:   value.Float32(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::Char`"),
		},

		"Char `a` >= `a`": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char `a` >= `b`": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.False,
		},
		"Char `b` >= `a`": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char `ś` >= `ą`": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.True,
		},
		"Char `ą` >= `ś`": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.False,
		},

		"String `a` >= 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.True,
		},
		"String `a` >= 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.False,
		},
		"String `b` >= 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.True,
		},
		"String `a` >= 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String `b` >= 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.True,
		},
		"String `ś` >= 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.True,
		},
		"String `ą` >= 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
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

func TestChar_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt and return an error": {
			a:   value.Char('a'),
			b:   value.SmallInt(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   value.Char('a'),
			b:   value.Float(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   value.Char('a'),
			b:   value.NewBigFloat(5),
			err: value.NewError(value.TypeErrorClass, "`Std::BigFloat` cannot be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   value.Char('a'),
			b:   value.Int64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   value.Char('a'),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   value.Char('a'),
			b:   value.Int16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int16` cannot be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   value.Char('a'),
			b:   value.Int8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int8` cannot be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt64` cannot be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt32` cannot be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt16` cannot be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt8` cannot be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   value.Char('a'),
			b:   value.Float64(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   value.Char('a'),
			b:   value.Float32(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::Char`"),
		},

		"Char `a` < `a`": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char `a` < `b`": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.True,
		},
		"Char `b` < `a`": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char `ś` < `ą`": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.False,
		},
		"Char `ą` < `ś`": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.True,
		},

		"String `a` < 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.False,
		},
		"String `a` < 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.True,
		},
		"String `b` < 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.False,
		},
		"String `a` < 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.True,
		},
		"String `b` < 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String `ś` < 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.False,
		},
		"String `ą` < 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
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

func TestChar_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"SmallInt and return an error": {
			a:   value.Char('a'),
			b:   value.SmallInt(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int` cannot be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   value.Char('a'),
			b:   value.Float(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float` cannot be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   value.Char('a'),
			b:   value.NewBigFloat(5),
			err: value.NewError(value.TypeErrorClass, "`Std::BigFloat` cannot be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   value.Char('a'),
			b:   value.Int64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   value.Char('a'),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   value.Char('a'),
			b:   value.Int16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int16` cannot be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   value.Char('a'),
			b:   value.Int8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int8` cannot be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt64(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt64` cannot be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt32` cannot be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt16(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt16` cannot be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   value.Char('a'),
			b:   value.UInt8(2),
			err: value.NewError(value.TypeErrorClass, "`Std::UInt8` cannot be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   value.Char('a'),
			b:   value.Float64(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   value.Char('a'),
			b:   value.Float32(5),
			err: value.NewError(value.TypeErrorClass, "`Std::Float32` cannot be coerced into `Std::Char`"),
		},

		"Char `a` <= `a`": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char `a` <= `b`": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.True,
		},
		"Char `b` <= `a`": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char `ś` <= `ą`": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.False,
		},
		"Char `ą` <= `ś`": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
			want: value.True,
		},

		"String `a` <= 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.True,
		},
		"String `a` <= 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.True,
		},
		"String `b` <= 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.False,
		},
		"String `a` <= 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.True,
		},
		"String `b` <= 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String `ś` <= 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.False,
		},
		"String `ą` <= 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
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

func TestChar_LaxEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
	}{
		"SmallInt `2` =~ 2": {
			a:    value.Char('2'),
			b:    value.SmallInt(2),
			want: value.False,
		},
		"Float `5` =~ 5.0": {
			a:    value.Char('5'),
			b:    value.Float(5),
			want: value.False,
		},
		"BigFloat `4` =~ 4bf": {
			a:    value.Char('4'),
			b:    value.NewBigFloat(4),
			want: value.False,
		},
		"Int64 `5` =~ 5i64": {
			a:    value.Char('5'),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int32 `2` =~ 2i32": {
			a:    value.Char('2'),
			b:    value.Int32(2),
			want: value.False,
		},
		"Int16 `8` =~ 8i16": {
			a:    value.Char('8'),
			b:    value.Int16(8),
			want: value.False,
		},
		"Int8 `8` =~ 8i8": {
			a:    value.Char('8'),
			b:    value.Int8(8),
			want: value.False,
		},
		"UInt64 `3` =~ 3u64": {
			a:    value.Char('3'),
			b:    value.UInt64(3),
			want: value.False,
		},
		"UInt32 `9` =~ 9u32": {
			a:    value.Char('9'),
			b:    value.UInt32(9),
			want: value.False,
		},
		"UInt16 `7` =~ 7u16": {
			a:    value.Char('7'),
			b:    value.UInt16(7),
			want: value.False,
		},
		"UInt8 '1' =~ 1u8": {
			a:    value.Char('1'),
			b:    value.Int8(12),
			want: value.False,
		},
		"Float64 `0` =~ 0f64": {
			a:    value.Char('0'),
			b:    value.Float64(0),
			want: value.False,
		},
		"Float32 `5` =~ 5f32": {
			a:    value.Char('5'),
			b:    value.Float32(5),
			want: value.False,
		},

		"String `a` =~ 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.True,
		},
		"String `a` =~ 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.False,
		},
		"String `b` =~ 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.False,
		},
		"String `a` =~ 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String `b` =~ 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String `ś` =~ 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.False,
		},
		"String `ś` =~ 'ś'": {
			a:    value.Char('ś'),
			b:    value.String("ś"),
			want: value.True,
		},
		"String `ą` =~ 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
			want: value.False,
		},

		"Char `a` =~ `a`": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char `a` =~ `b`": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.False,
		},
		"Char `b` =~ `a`": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char `ś` =~ `ą`": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.False,
		},
		"Char `ś` =~ `ś`": {
			a:    value.Char('ś'),
			b:    value.Char('ś'),
			want: value.True,
		},
		"Char `ą` =~ `ś`": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
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

func TestChar_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.Char
		b    value.Value
		want value.Value
	}{
		"SmallInt `2` == 2": {
			a:    value.Char('2'),
			b:    value.SmallInt(2),
			want: value.False,
		},
		"Float `5` == 5.0": {
			a:    value.Char('5'),
			b:    value.Float(5),
			want: value.False,
		},
		"BigFloat `4` == 4bf": {
			a:    value.Char('4'),
			b:    value.NewBigFloat(4),
			want: value.False,
		},
		"Int64 `5` == 5i64": {
			a:    value.Char('5'),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int32 `2` == 2i32": {
			a:    value.Char('2'),
			b:    value.Int32(2),
			want: value.False,
		},
		"Int16 `8` == 8i16": {
			a:    value.Char('8'),
			b:    value.Int16(8),
			want: value.False,
		},
		"Int8 `8` == 8i8": {
			a:    value.Char('8'),
			b:    value.Int8(8),
			want: value.False,
		},
		"UInt64 `3` == 3u64": {
			a:    value.Char('3'),
			b:    value.UInt64(3),
			want: value.False,
		},
		"UInt32 `9` == 9u32": {
			a:    value.Char('9'),
			b:    value.UInt32(9),
			want: value.False,
		},
		"UInt16 `7` == 7u16": {
			a:    value.Char('7'),
			b:    value.UInt16(7),
			want: value.False,
		},
		"UInt8 '1' == 1u8": {
			a:    value.Char('1'),
			b:    value.Int8(12),
			want: value.False,
		},
		"Float64 `0` == 0f64": {
			a:    value.Char('0'),
			b:    value.Float64(0),
			want: value.False,
		},
		"Float32 `5` == 5f32": {
			a:    value.Char('5'),
			b:    value.Float32(5),
			want: value.False,
		},

		"String `a` == 'a'": {
			a:    value.Char('a'),
			b:    value.String("a"),
			want: value.False,
		},
		"String `a` == 'b'": {
			a:    value.Char('a'),
			b:    value.String("b"),
			want: value.False,
		},
		"String `b` == 'a'": {
			a:    value.Char('b'),
			b:    value.String("a"),
			want: value.False,
		},
		"String `a` == 'aa'": {
			a:    value.Char('a'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String `b` == 'aa'": {
			a:    value.Char('b'),
			b:    value.String("aa"),
			want: value.False,
		},
		"String `ś` == 'ą'": {
			a:    value.Char('ś'),
			b:    value.String("ą"),
			want: value.False,
		},
		"String `ś` == 'ś'": {
			a:    value.Char('ś'),
			b:    value.String("ś"),
			want: value.False,
		},
		"String `ą` == 'ś'": {
			a:    value.Char('ą'),
			b:    value.String("ś"),
			want: value.False,
		},

		"Char `a` == `a`": {
			a:    value.Char('a'),
			b:    value.Char('a'),
			want: value.True,
		},
		"Char `a` == `b`": {
			a:    value.Char('a'),
			b:    value.Char('b'),
			want: value.False,
		},
		"Char `b` == `a`": {
			a:    value.Char('b'),
			b:    value.Char('a'),
			want: value.False,
		},
		"Char `ś` == `ą`": {
			a:    value.Char('ś'),
			b:    value.Char('ą'),
			want: value.False,
		},
		"Char `ś` == `ś`": {
			a:    value.Char('ś'),
			b:    value.Char('ś'),
			want: value.True,
		},
		"Char `ą` == `ś`": {
			a:    value.Char('ą'),
			b:    value.Char('ś'),
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
