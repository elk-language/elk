package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCharConcat(t *testing.T) {
	tests := map[string]struct {
		left  Char
		right Value
		want  String
		err   *Error
	}{
		"Char + String => String": {
			left:  Char('f'),
			right: String("oo"),
			want:  String("foo"),
		},
		"Char + Char => String": {
			left:  Char('a'),
			right: Char('b'),
			want:  String("ab"),
		},
		"Char + Int => TypeError": {
			left:  Char('a'),
			right: Int8(5),
			err:   NewError(TypeErrorClass, `can't concat 5i8 with char c"a"`),
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

func TestCharRepeat(t *testing.T) {
	tests := map[string]struct {
		left  Char
		right Value
		want  String
		err   *Error
	}{
		"Char * SmallInt => String": {
			left:  Char('a'),
			right: SmallInt(3),
			want:  String("aaa"),
		},
		"Char * BigInt => OutOfRangeError": {
			left:  Char('a'),
			right: NewBigInt(3),
			err:   NewError(OutOfRangeErrorClass, `repeat count is too large 3`),
		},
		"Char * Int8 => TypeError": {
			left:  Char('a'),
			right: Int8(3),
			err:   NewError(TypeErrorClass, `can't repeat a char using 3i8`),
		},
		"String * String => TypeError": {
			left:  Char('a'),
			right: String("bar"),
			err:   NewError(TypeErrorClass, `can't repeat a char using "bar"`),
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

func TestCharInspect(t *testing.T) {
	tests := map[string]struct {
		c    Char
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
