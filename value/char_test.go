package value

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

func TestChar_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    Char
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt and return an error": {
			a:   Char('a'),
			b:   SmallInt(2),
			err: NewError(TypeErrorClass, "`Std::SmallInt` can't be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   Char('a'),
			b:   Float(5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   Char('a'),
			b:   NewBigFloat(5),
			err: NewError(TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   Char('a'),
			b:   Int64(2),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   Char('a'),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   Char('a'),
			b:   Int16(2),
			err: NewError(TypeErrorClass, "`Std::Int16` can't be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   Char('a'),
			b:   Int8(2),
			err: NewError(TypeErrorClass, "`Std::Int8` can't be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   Char('a'),
			b:   UInt64(2),
			err: NewError(TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   Char('a'),
			b:   UInt32(2),
			err: NewError(TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   Char('a'),
			b:   UInt16(2),
			err: NewError(TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   Char('a'),
			b:   UInt8(2),
			err: NewError(TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   Char('a'),
			b:   Float64(5),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   Char('a'),
			b:   Float32(5),
			err: NewError(TypeErrorClass, "`Std::Float32` can't be coerced into `Std::Char`"),
		},

		"Char c'a' > c'a'": {
			a:    Char('a'),
			b:    Char('a'),
			want: False,
		},
		"Char c'a' > c'b'": {
			a:    Char('a'),
			b:    Char('b'),
			want: False,
		},
		"Char c'b' > c'a'": {
			a:    Char('b'),
			b:    Char('a'),
			want: True,
		},
		"Char c'ś' > c'ą'": {
			a:    Char('ś'),
			b:    Char('ą'),
			want: True,
		},
		"Char c'ą' > c'ś'": {
			a:    Char('ą'),
			b:    Char('ś'),
			want: False,
		},

		"String c'a' > 'a'": {
			a:    Char('a'),
			b:    String("a"),
			want: False,
		},
		"String c'a' > 'b'": {
			a:    Char('a'),
			b:    String("b"),
			want: False,
		},
		"String c'b' > 'a'": {
			a:    Char('b'),
			b:    String("a"),
			want: True,
		},
		"String c'a' > 'aa'": {
			a:    Char('a'),
			b:    String("aa"),
			want: False,
		},
		"String c'b' > 'aa'": {
			a:    Char('b'),
			b:    String("aa"),
			want: True,
		},
		"String c'ś' > 'ą'": {
			a:    Char('ś'),
			b:    String("ą"),
			want: True,
		},
		"String c'ą' > 'ś'": {
			a:    Char('ą'),
			b:    String("ś"),
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

func TestChar_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    Char
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt and return an error": {
			a:   Char('a'),
			b:   SmallInt(2),
			err: NewError(TypeErrorClass, "`Std::SmallInt` can't be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   Char('a'),
			b:   Float(5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   Char('a'),
			b:   NewBigFloat(5),
			err: NewError(TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   Char('a'),
			b:   Int64(2),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   Char('a'),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   Char('a'),
			b:   Int16(2),
			err: NewError(TypeErrorClass, "`Std::Int16` can't be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   Char('a'),
			b:   Int8(2),
			err: NewError(TypeErrorClass, "`Std::Int8` can't be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   Char('a'),
			b:   UInt64(2),
			err: NewError(TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   Char('a'),
			b:   UInt32(2),
			err: NewError(TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   Char('a'),
			b:   UInt16(2),
			err: NewError(TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   Char('a'),
			b:   UInt8(2),
			err: NewError(TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   Char('a'),
			b:   Float64(5),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   Char('a'),
			b:   Float32(5),
			err: NewError(TypeErrorClass, "`Std::Float32` can't be coerced into `Std::Char`"),
		},

		"Char c'a' >= c'a'": {
			a:    Char('a'),
			b:    Char('a'),
			want: True,
		},
		"Char c'a' >= c'b'": {
			a:    Char('a'),
			b:    Char('b'),
			want: False,
		},
		"Char c'b' >= c'a'": {
			a:    Char('b'),
			b:    Char('a'),
			want: True,
		},
		"Char c'ś' >= c'ą'": {
			a:    Char('ś'),
			b:    Char('ą'),
			want: True,
		},
		"Char c'ą' >= c'ś'": {
			a:    Char('ą'),
			b:    Char('ś'),
			want: False,
		},

		"String c'a' >= 'a'": {
			a:    Char('a'),
			b:    String("a"),
			want: True,
		},
		"String c'a' >= 'b'": {
			a:    Char('a'),
			b:    String("b"),
			want: False,
		},
		"String c'b' >= 'a'": {
			a:    Char('b'),
			b:    String("a"),
			want: True,
		},
		"String c'a' >= 'aa'": {
			a:    Char('a'),
			b:    String("aa"),
			want: False,
		},
		"String c'b' >= 'aa'": {
			a:    Char('b'),
			b:    String("aa"),
			want: True,
		},
		"String c'ś' >= 'ą'": {
			a:    Char('ś'),
			b:    String("ą"),
			want: True,
		},
		"String c'ą' >= 'ś'": {
			a:    Char('ą'),
			b:    String("ś"),
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

func TestChar_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    Char
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt and return an error": {
			a:   Char('a'),
			b:   SmallInt(2),
			err: NewError(TypeErrorClass, "`Std::SmallInt` can't be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   Char('a'),
			b:   Float(5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   Char('a'),
			b:   NewBigFloat(5),
			err: NewError(TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   Char('a'),
			b:   Int64(2),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   Char('a'),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   Char('a'),
			b:   Int16(2),
			err: NewError(TypeErrorClass, "`Std::Int16` can't be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   Char('a'),
			b:   Int8(2),
			err: NewError(TypeErrorClass, "`Std::Int8` can't be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   Char('a'),
			b:   UInt64(2),
			err: NewError(TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   Char('a'),
			b:   UInt32(2),
			err: NewError(TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   Char('a'),
			b:   UInt16(2),
			err: NewError(TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   Char('a'),
			b:   UInt8(2),
			err: NewError(TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   Char('a'),
			b:   Float64(5),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   Char('a'),
			b:   Float32(5),
			err: NewError(TypeErrorClass, "`Std::Float32` can't be coerced into `Std::Char`"),
		},

		"Char c'a' < c'a'": {
			a:    Char('a'),
			b:    Char('a'),
			want: False,
		},
		"Char c'a' < c'b'": {
			a:    Char('a'),
			b:    Char('b'),
			want: True,
		},
		"Char c'b' < c'a'": {
			a:    Char('b'),
			b:    Char('a'),
			want: False,
		},
		"Char c'ś' < c'ą'": {
			a:    Char('ś'),
			b:    Char('ą'),
			want: False,
		},
		"Char c'ą' < c'ś'": {
			a:    Char('ą'),
			b:    Char('ś'),
			want: True,
		},

		"String c'a' < 'a'": {
			a:    Char('a'),
			b:    String("a"),
			want: False,
		},
		"String c'a' < 'b'": {
			a:    Char('a'),
			b:    String("b"),
			want: True,
		},
		"String c'b' < 'a'": {
			a:    Char('b'),
			b:    String("a"),
			want: False,
		},
		"String c'a' < 'aa'": {
			a:    Char('a'),
			b:    String("aa"),
			want: True,
		},
		"String c'b' < 'aa'": {
			a:    Char('b'),
			b:    String("aa"),
			want: False,
		},
		"String c'ś' < 'ą'": {
			a:    Char('ś'),
			b:    String("ą"),
			want: False,
		},
		"String c'ą' < 'ś'": {
			a:    Char('ą'),
			b:    String("ś"),
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

func TestChar_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    Char
		b    Value
		want Value
		err  *Error
	}{
		"SmallInt and return an error": {
			a:   Char('a'),
			b:   SmallInt(2),
			err: NewError(TypeErrorClass, "`Std::SmallInt` can't be coerced into `Std::Char`"),
		},
		"Float and return an error": {
			a:   Char('a'),
			b:   Float(5),
			err: NewError(TypeErrorClass, "`Std::Float` can't be coerced into `Std::Char`"),
		},
		"BigFloat and return an error": {
			a:   Char('a'),
			b:   NewBigFloat(5),
			err: NewError(TypeErrorClass, "`Std::BigFloat` can't be coerced into `Std::Char`"),
		},
		"Int64 and return an error": {
			a:   Char('a'),
			b:   Int64(2),
			err: NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::Char`"),
		},
		"Int32 and return an error": {
			a:   Char('a'),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Char`"),
		},
		"Int16 and return an error": {
			a:   Char('a'),
			b:   Int16(2),
			err: NewError(TypeErrorClass, "`Std::Int16` can't be coerced into `Std::Char`"),
		},
		"Int8 and return an error": {
			a:   Char('a'),
			b:   Int8(2),
			err: NewError(TypeErrorClass, "`Std::Int8` can't be coerced into `Std::Char`"),
		},
		"UInt64 and return an error": {
			a:   Char('a'),
			b:   UInt64(2),
			err: NewError(TypeErrorClass, "`Std::UInt64` can't be coerced into `Std::Char`"),
		},
		"UInt32 and return an error": {
			a:   Char('a'),
			b:   UInt32(2),
			err: NewError(TypeErrorClass, "`Std::UInt32` can't be coerced into `Std::Char`"),
		},
		"UInt16 and return an error": {
			a:   Char('a'),
			b:   UInt16(2),
			err: NewError(TypeErrorClass, "`Std::UInt16` can't be coerced into `Std::Char`"),
		},
		"UInt8 and return an error": {
			a:   Char('a'),
			b:   UInt8(2),
			err: NewError(TypeErrorClass, "`Std::UInt8` can't be coerced into `Std::Char`"),
		},
		"Float64 and return an error": {
			a:   Char('a'),
			b:   Float64(5),
			err: NewError(TypeErrorClass, "`Std::Float64` can't be coerced into `Std::Char`"),
		},
		"Float32 and return an error": {
			a:   Char('a'),
			b:   Float32(5),
			err: NewError(TypeErrorClass, "`Std::Float32` can't be coerced into `Std::Char`"),
		},

		"Char c'a' <= c'a'": {
			a:    Char('a'),
			b:    Char('a'),
			want: True,
		},
		"Char c'a' <= c'b'": {
			a:    Char('a'),
			b:    Char('b'),
			want: True,
		},
		"Char c'b' <= c'a'": {
			a:    Char('b'),
			b:    Char('a'),
			want: False,
		},
		"Char c'ś' <= c'ą'": {
			a:    Char('ś'),
			b:    Char('ą'),
			want: False,
		},
		"Char c'ą' <= c'ś'": {
			a:    Char('ą'),
			b:    Char('ś'),
			want: True,
		},

		"String c'a' <= 'a'": {
			a:    Char('a'),
			b:    String("a"),
			want: True,
		},
		"String c'a' <= 'b'": {
			a:    Char('a'),
			b:    String("b"),
			want: True,
		},
		"String c'b' <= 'a'": {
			a:    Char('b'),
			b:    String("a"),
			want: False,
		},
		"String c'a' <= 'aa'": {
			a:    Char('a'),
			b:    String("aa"),
			want: True,
		},
		"String c'b' <= 'aa'": {
			a:    Char('b'),
			b:    String("aa"),
			want: False,
		},
		"String c'ś' <= 'ą'": {
			a:    Char('ś'),
			b:    String("ą"),
			want: False,
		},
		"String c'ą' <= 'ś'": {
			a:    Char('ą'),
			b:    String("ś"),
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

func TestChar_Equal(t *testing.T) {
	tests := map[string]struct {
		a    Char
		b    Value
		want Value
	}{
		"SmallInt c'2' == 2": {
			a:    Char('2'),
			b:    SmallInt(2),
			want: False,
		},
		"Float c'5' == 5.0": {
			a:    Char('5'),
			b:    Float(5),
			want: False,
		},
		"BigFloat c'4' == 4bf": {
			a:    Char('4'),
			b:    NewBigFloat(4),
			want: False,
		},
		"Int64 c'5' == 5i64": {
			a:    Char('5'),
			b:    Int64(5),
			want: False,
		},
		"Int32 c'2' == 2i32": {
			a:    Char('2'),
			b:    Int32(2),
			want: False,
		},
		"Int16 c'8' == 8i16": {
			a:    Char('8'),
			b:    Int16(8),
			want: False,
		},
		"Int8 c'8' == 8i8": {
			a:    Char('8'),
			b:    Int8(8),
			want: False,
		},
		"UInt64 c'3' == 3u64": {
			a:    Char('3'),
			b:    UInt64(3),
			want: False,
		},
		"UInt32 c'9' == 9u32": {
			a:    Char('9'),
			b:    UInt32(9),
			want: False,
		},
		"UInt16 c'7' == 7u16": {
			a:    Char('7'),
			b:    UInt16(7),
			want: False,
		},
		"UInt8 '1' == 1u8": {
			a:    Char('1'),
			b:    Int8(12),
			want: False,
		},
		"Float64 c'0' == 0f64": {
			a:    Char('0'),
			b:    Float64(0),
			want: False,
		},
		"Float32 c'5' == 5f32": {
			a:    Char('5'),
			b:    Float32(5),
			want: False,
		},

		"String c'a' == 'a'": {
			a:    Char('a'),
			b:    String("a"),
			want: True,
		},
		"String c'a' == 'b'": {
			a:    Char('a'),
			b:    String("b"),
			want: False,
		},
		"String c'b' == 'a'": {
			a:    Char('b'),
			b:    String("a"),
			want: False,
		},
		"String c'a' == 'aa'": {
			a:    Char('a'),
			b:    String("aa"),
			want: False,
		},
		"String c'b' == 'aa'": {
			a:    Char('b'),
			b:    String("aa"),
			want: False,
		},
		"String c'ś' == 'ą'": {
			a:    Char('ś'),
			b:    String("ą"),
			want: False,
		},
		"String c'ś' == 'ś'": {
			a:    Char('ś'),
			b:    String("ś"),
			want: True,
		},
		"String c'ą' == 'ś'": {
			a:    Char('ą'),
			b:    String("ś"),
			want: False,
		},

		"Char c'a' == c'a'": {
			a:    Char('a'),
			b:    Char('a'),
			want: True,
		},
		"Char c'a' == c'b'": {
			a:    Char('a'),
			b:    Char('b'),
			want: False,
		},
		"Char c'b' == c'a'": {
			a:    Char('b'),
			b:    Char('a'),
			want: False,
		},
		"Char c'ś' == c'ą'": {
			a:    Char('ś'),
			b:    Char('ą'),
			want: False,
		},
		"Char c'ś' == c'ś'": {
			a:    Char('ś'),
			b:    Char('ś'),
			want: True,
		},
		"Char c'ą' == c'ś'": {
			a:    Char('ą'),
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

func TestChar_StrictEqual(t *testing.T) {
	tests := map[string]struct {
		a    Char
		b    Value
		want Value
	}{
		"SmallInt c'2' === 2": {
			a:    Char('2'),
			b:    SmallInt(2),
			want: False,
		},
		"Float c'5' === 5.0": {
			a:    Char('5'),
			b:    Float(5),
			want: False,
		},
		"BigFloat c'4' === 4bf": {
			a:    Char('4'),
			b:    NewBigFloat(4),
			want: False,
		},
		"Int64 c'5' === 5i64": {
			a:    Char('5'),
			b:    Int64(5),
			want: False,
		},
		"Int32 c'2' === 2i32": {
			a:    Char('2'),
			b:    Int32(2),
			want: False,
		},
		"Int16 c'8' === 8i16": {
			a:    Char('8'),
			b:    Int16(8),
			want: False,
		},
		"Int8 c'8' === 8i8": {
			a:    Char('8'),
			b:    Int8(8),
			want: False,
		},
		"UInt64 c'3' === 3u64": {
			a:    Char('3'),
			b:    UInt64(3),
			want: False,
		},
		"UInt32 c'9' === 9u32": {
			a:    Char('9'),
			b:    UInt32(9),
			want: False,
		},
		"UInt16 c'7' === 7u16": {
			a:    Char('7'),
			b:    UInt16(7),
			want: False,
		},
		"UInt8 '1' === 1u8": {
			a:    Char('1'),
			b:    Int8(12),
			want: False,
		},
		"Float64 c'0' === 0f64": {
			a:    Char('0'),
			b:    Float64(0),
			want: False,
		},
		"Float32 c'5' === 5f32": {
			a:    Char('5'),
			b:    Float32(5),
			want: False,
		},

		"String c'a' === 'a'": {
			a:    Char('a'),
			b:    String("a"),
			want: False,
		},
		"String c'a' === 'b'": {
			a:    Char('a'),
			b:    String("b"),
			want: False,
		},
		"String c'b' === 'a'": {
			a:    Char('b'),
			b:    String("a"),
			want: False,
		},
		"String c'a' === 'aa'": {
			a:    Char('a'),
			b:    String("aa"),
			want: False,
		},
		"String c'b' === 'aa'": {
			a:    Char('b'),
			b:    String("aa"),
			want: False,
		},
		"String c'ś' === 'ą'": {
			a:    Char('ś'),
			b:    String("ą"),
			want: False,
		},
		"String c'ś' === 'ś'": {
			a:    Char('ś'),
			b:    String("ś"),
			want: False,
		},
		"String c'ą' === 'ś'": {
			a:    Char('ą'),
			b:    String("ś"),
			want: False,
		},

		"Char c'a' === c'a'": {
			a:    Char('a'),
			b:    Char('a'),
			want: True,
		},
		"Char c'a' === c'b'": {
			a:    Char('a'),
			b:    Char('b'),
			want: False,
		},
		"Char c'b' === c'a'": {
			a:    Char('b'),
			b:    Char('a'),
			want: False,
		},
		"Char c'ś' === c'ą'": {
			a:    Char('ś'),
			b:    Char('ą'),
			want: False,
		},
		"Char c'ś' === c'ś'": {
			a:    Char('ś'),
			b:    Char('ś'),
			want: True,
		},
		"Char c'ą' === c'ś'": {
			a:    Char('ą'),
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
