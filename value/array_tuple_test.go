package value_test

import (
	"regexp"
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestArrayTuple_Concat(t *testing.T) {
	tests := map[string]struct {
		left  *value.ArrayTuple
		right value.Value
		want  value.Value
		err   value.Value
	}{
		"ArrayTuple + ArrayTuple => ArrayTuple": {
			left:  &value.ArrayTuple{value.SmallInt(2).ToValue()},
			right: value.Ref(&value.ArrayTuple{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))}),
			want:  value.Ref(&value.ArrayTuple{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))}),
		},
		"ArrayTuple + ArrayList => ArrayList": {
			left:  &value.ArrayTuple{value.SmallInt(2).ToValue()},
			right: value.Ref(&value.ArrayList{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))}),
			want:  value.Ref(&value.ArrayList{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))}),
		},
		"ArrayTuple + Int => TypeError": {
			left:  &value.ArrayTuple{value.SmallInt(2).ToValue()},
			right: value.Int8(5).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot concat 5i8 with arrayTuple %[2]`)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.ConcatVal(tc.right)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayTuple_Repeat(t *testing.T) {
	tests := map[string]struct {
		left  *value.ArrayTuple
		right value.Value
		want  *value.ArrayTuple
		err   value.Value
	}{
		"ArrayTuple * SmallInt => ArrayTuple": {
			left:  &value.ArrayTuple{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.SmallInt(3).ToValue(),
			want:  &value.ArrayTuple{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
		},
		"ArrayTuple * -SmallInt => OutOfRangeError": {
			left:  &value.ArrayTuple{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.SmallInt(-3).ToValue(),
			err:   value.Ref(value.NewError(value.OutOfRangeErrorClass, `arrayTuple repeat count cannot be negative: -3`)),
		},
		"ArrayTuple * 0 => ArrayTuple": {
			left:  &value.ArrayTuple{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.SmallInt(0).ToValue(),
			want:  &value.ArrayTuple{},
		},
		"ArrayTuple * BigInt => OutOfRangeError": {
			left:  &value.ArrayTuple{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.Ref(value.NewBigInt(3)),
			err:   value.Ref(value.NewError(value.OutOfRangeErrorClass, `arrayTuple repeat count is too large 3`)),
		},
		"ArrayTuple * Int8 => TypeError": {
			left:  &value.ArrayTuple{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.Int8(3).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot repeat a arrayTuple using 3i8`)),
		},
		"ArrayTuple * String => TypeError": {
			left:  &value.ArrayTuple{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.Ref(value.String("bar")),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot repeat a arrayTuple using "bar"`)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Repeat(tc.right)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayTuple_Inspect(t *testing.T) {
	tests := map[string]struct {
		t    *value.ArrayTuple
		want string
	}{
		"empty": {
			t:    &value.ArrayTuple{},
			want: "%[]",
		},
		"with one element": {
			t:    &value.ArrayTuple{value.SmallInt(3).ToValue()},
			want: `%[3]`,
		},
		"with elements": {
			t:    &value.ArrayTuple{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
			want: `%[3, "foo"]`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.t.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayTuple_Expand(t *testing.T) {
	tests := map[string]struct {
		t    *value.ArrayTuple
		new  int
		want *value.ArrayTuple
	}{
		"add 3 to an empty list": {
			t:    &value.ArrayTuple{},
			new:  3,
			want: &value.ArrayTuple{value.Nil, value.Nil, value.Nil},
		},
		"add 0 to an empty list": {
			t:    &value.ArrayTuple{},
			new:  0,
			want: &value.ArrayTuple{},
		},
		"add 2 to a filled list": {
			t:    &value.ArrayTuple{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue()},
			new:  2,
			want: &value.ArrayTuple{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue(), value.Nil, value.Nil},
		},
		"add 0 to a filled list": {
			t:    &value.ArrayTuple{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue()},
			new:  0,
			want: &value.ArrayTuple{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue()},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.t.Expand(tc.new)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, tc.t, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayTuple_Append(t *testing.T) {
	tests := map[string]struct {
		t    *value.ArrayTuple
		val  value.Value
		want *value.ArrayTuple
	}{
		"append to an empty list": {
			t:    &value.ArrayTuple{},
			val:  value.SmallInt(3).ToValue(),
			want: &value.ArrayTuple{value.SmallInt(3).ToValue()},
		},
		"append to a filled list": {
			t:    &value.ArrayTuple{value.SmallInt(3).ToValue()},
			val:  value.Ref(value.String("foo")),
			want: &value.ArrayTuple{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.t.Append(tc.val)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, tc.t, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayTuple_SubscriptSet(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayTuple
		key  value.Value
		val  value.Value
		want *value.ArrayTuple
		err  value.Value
	}{
		"set index 0 in an empty arrayTuple": {
			l:    &value.ArrayTuple{},
			key:  value.SmallInt(0).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayTuple{},
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			)),
		},
		"set index 0 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(0).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayTuple{value.SmallInt(25).ToValue(), value.Ref(value.String("foo"))},
		},
		"set index -1 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.Float(89.2).ToValue(), value.Ref(value.String("foo"))},
			key:  value.SmallInt(-1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayTuple{value.Nil, value.Float(89.2).ToValue(), value.SmallInt(25).ToValue()},
		},
		"set index -2 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Float(89.2).ToValue(), value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(-2).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayTuple{value.Float(89.2).ToValue(), value.SmallInt(25).ToValue(), value.Ref(value.String("foo"))},
		},
		"set index in the middle of a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.SmallInt(1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayTuple{value.Nil, value.SmallInt(25).ToValue(), value.Float(21.37).ToValue()},
		},
		"set uint8 index": {
			l:    &value.ArrayTuple{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.UInt8(1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayTuple{value.Nil, value.SmallInt(25).ToValue(), value.Float(21.37).ToValue()},
		},
		"set string index": {
			l:    &value.ArrayTuple{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.Ref(value.String("lol")),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayTuple{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			)),
		},
		"set float index": {
			l:    &value.ArrayTuple{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.Float(1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayTuple{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int`",
			)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.l.SubscriptSet(tc.key, tc.val)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, tc.l, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayTuple_Subscript(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayTuple
		key  value.Value
		want value.Value
		err  value.Value
	}{
		"get index 0 in an empty arrayTuple": {
			l:   &value.ArrayTuple{},
			key: value.SmallInt(0).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			)),
		},
		"get index 0 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(0).ToValue(),
			want: value.Nil,
		},
		"get index out of range": {
			l:   &value.ArrayTuple{value.Nil, value.Float(89.2).ToValue(), value.Ref(value.String("foo"))},
			key: value.SmallInt(31).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 31 out of range: -3...3",
			)),
		},
		"get negative index out of range": {
			l:   &value.ArrayTuple{value.Nil, value.Float(89.2).ToValue(), value.Ref(value.String("foo"))},
			key: value.SmallInt(-31).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index -31 out of range: -3...3",
			)),
		},
		"get index -1 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.Float(89.2).ToValue(), value.Ref(value.String("foo"))},
			key:  value.SmallInt(-1).ToValue(),
			want: value.Ref(value.String("foo")),
		},
		"get index -2 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Float(89.2).ToValue(), value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(-2).ToValue(),
			want: value.Nil,
		},
		"get index in the middle of a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.SmallInt(1).ToValue(),
			want: value.Ref(value.String("foo")),
		},
		"get uint8 index": {
			l:    &value.ArrayTuple{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.UInt8(1).ToValue(),
			want: value.Ref(value.String("foo")),
		},
		"get string index": {
			l:   &value.ArrayTuple{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key: value.Ref(value.String("lol")),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			)),
		},
		"get float index": {
			l:   &value.ArrayTuple{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key: value.Float(1).ToValue(),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int`",
			)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			want, err := tc.l.Subscript(tc.key)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
			if !tc.err.IsNil() {
				return
			}
			if diff := cmp.Diff(tc.want, want, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayTupleIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayTupleIterator
		want string
	}{
		"empty": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{},
				0,
			),
			want: `Std::ArrayTuple::Iterator\{&: 0x[[:xdigit:]]{4,12}, tuple: %\[\], index: 0\}`,
		},
		"with one element": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3).ToValue()},
				0,
			),
			want: `Std::ArrayTuple::Iterator\{&: 0x[[:xdigit:]]{4,12}, tuple: %\[3\], index: 0\}`,
		},
		"with elements": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
				1,
			),
			want: `Std::ArrayTuple::Iterator\{&: 0x[[:xdigit:]]{4,12}, tuple: %\[3, "foo"\], index: 1\}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.l.Inspect()
			ok, _ := regexp.MatchString(tc.want, got)
			if !ok {
				t.Fatalf("got %q, expected to match pattern %q", got, tc.want)
			}
		})
	}
}

func TestArrayTupleIterator_Next(t *testing.T) {
	tests := map[string]struct {
		l     *value.ArrayTupleIterator
		after *value.ArrayTupleIterator
		want  value.Value
		err   value.Value
	}{
		"empty": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{},
				0,
			),
			after: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{},
				0,
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
		"with two elements index 0": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				0,
			),
			after: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				1,
			),
			want: value.SmallInt(3).ToValue(),
		},
		"with two elements index 1": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				1,
			),
			after: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				2,
			),
			want: value.Float(7.2).ToValue(),
		},
		"with two elements index 2": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				2,
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.l.Next()
			opts := comparer.Options()
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
			if !tc.err.IsNil() {
				return
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.after, tc.l, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
