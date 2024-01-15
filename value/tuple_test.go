package value_test

import (
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
		err   *value.Error
	}{
		"ArrayTuple + ArrayTuple => ArrayTuple": {
			left:  &value.ArrayTuple{value.SmallInt(2)},
			right: &value.ArrayTuple{value.SmallInt(3), value.String("foo")},
			want:  &value.ArrayTuple{value.SmallInt(2), value.SmallInt(3), value.String("foo")},
		},
		"ArrayTuple + List => List": {
			left:  &value.ArrayTuple{value.SmallInt(2)},
			right: &value.List{value.SmallInt(3), value.String("foo")},
			want:  &value.List{value.SmallInt(2), value.SmallInt(3), value.String("foo")},
		},
		"ArrayTuple + Int => TypeError": {
			left:  &value.ArrayTuple{value.SmallInt(2)},
			right: value.Int8(5),
			err:   value.NewError(value.TypeErrorClass, `cannot concat 5i8 with arrayTuple %[2]`),
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

func TestArrayTuple_Repeat(t *testing.T) {
	tests := map[string]struct {
		left  *value.ArrayTuple
		right value.Value
		want  *value.ArrayTuple
		err   *value.Error
	}{
		"ArrayTuple * SmallInt => ArrayTuple": {
			left:  &value.ArrayTuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(3),
			want:  &value.ArrayTuple{value.SmallInt(2), value.SmallInt(3), value.SmallInt(2), value.SmallInt(3), value.SmallInt(2), value.SmallInt(3)},
		},
		"ArrayTuple * -SmallInt => OutOfRangeError": {
			left:  &value.ArrayTuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(-3),
			err:   value.NewError(value.OutOfRangeErrorClass, `arrayTuple repeat count cannot be negative: -3`),
		},
		"ArrayTuple * 0 => ArrayTuple": {
			left:  &value.ArrayTuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(0),
			want:  &value.ArrayTuple{},
		},
		"ArrayTuple * BigInt => OutOfRangeError": {
			left:  &value.ArrayTuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.NewBigInt(3),
			err:   value.NewError(value.OutOfRangeErrorClass, `arrayTuple repeat count is too large 3`),
		},
		"ArrayTuple * Int8 => TypeError": {
			left:  &value.ArrayTuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.Int8(3),
			err:   value.NewError(value.TypeErrorClass, `cannot repeat a arrayTuple using 3i8`),
		},
		"ArrayTuple * String => TypeError": {
			left:  &value.ArrayTuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.String("bar"),
			err:   value.NewError(value.TypeErrorClass, `cannot repeat a arrayTuple using "bar"`),
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
			t:    &value.ArrayTuple{value.SmallInt(3)},
			want: `%[3]`,
		},
		"with elements": {
			t:    &value.ArrayTuple{value.SmallInt(3), value.String("foo")},
			want: `%[3, "foo"]`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.t.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
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
			t:    &value.ArrayTuple{value.SmallInt(-3), value.Float(10.5)},
			new:  2,
			want: &value.ArrayTuple{value.SmallInt(-3), value.Float(10.5), value.Nil, value.Nil},
		},
		"add 0 to a filled list": {
			t:    &value.ArrayTuple{value.SmallInt(-3), value.Float(10.5)},
			new:  0,
			want: &value.ArrayTuple{value.SmallInt(-3), value.Float(10.5)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.t.Expand(tc.new)
			if diff := cmp.Diff(tc.want, tc.t); diff != "" {
				t.Fatalf(diff)
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
			val:  value.SmallInt(3),
			want: &value.ArrayTuple{value.SmallInt(3)},
		},
		"append to a filled list": {
			t:    &value.ArrayTuple{value.SmallInt(3)},
			val:  value.String("foo"),
			want: &value.ArrayTuple{value.SmallInt(3), value.String("foo")},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.t.Append(tc.val)
			if diff := cmp.Diff(tc.want, tc.t); diff != "" {
				t.Fatalf(diff)
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
		err  *value.Error
	}{
		"set index 0 in an empty arrayTuple": {
			l:    &value.ArrayTuple{},
			key:  value.SmallInt(0),
			val:  value.SmallInt(25),
			want: &value.ArrayTuple{},
			err: value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			),
		},
		"set index 0 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.String("foo")},
			key:  value.SmallInt(0),
			val:  value.SmallInt(25),
			want: &value.ArrayTuple{value.SmallInt(25), value.String("foo")},
		},
		"set index -1 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.Float(89.2), value.String("foo")},
			key:  value.SmallInt(-1),
			val:  value.SmallInt(25),
			want: &value.ArrayTuple{value.Nil, value.Float(89.2), value.SmallInt(25)},
		},
		"set index -2 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Float(89.2), value.Nil, value.String("foo")},
			key:  value.SmallInt(-2),
			val:  value.SmallInt(25),
			want: &value.ArrayTuple{value.Float(89.2), value.SmallInt(25), value.String("foo")},
		},
		"set index in the middle of a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.SmallInt(1),
			val:  value.SmallInt(25),
			want: &value.ArrayTuple{value.Nil, value.SmallInt(25), value.Float(21.37)},
		},
		"set uint8 index": {
			l:    &value.ArrayTuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.UInt8(1),
			val:  value.SmallInt(25),
			want: &value.ArrayTuple{value.Nil, value.SmallInt(25), value.Float(21.37)},
		},
		"set string index": {
			l:    &value.ArrayTuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.String("lol"),
			val:  value.SmallInt(25),
			want: &value.ArrayTuple{value.Nil, value.String("foo"), value.Float(21.37)},
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			),
		},
		"set float index": {
			l:    &value.ArrayTuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.Float(1),
			val:  value.SmallInt(25),
			want: &value.ArrayTuple{value.Nil, value.String("foo"), value.Float(21.37)},
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int`",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.l.SubscriptSet(tc.key, tc.val)
			if diff := cmp.Diff(tc.want, tc.l, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestArrayTuple_Subscript(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayTuple
		key  value.Value
		want value.Value
		err  *value.Error
	}{
		"get index 0 in an empty arrayTuple": {
			l:   &value.ArrayTuple{},
			key: value.SmallInt(0),
			err: value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			),
		},
		"get index 0 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.String("foo")},
			key:  value.SmallInt(0),
			want: value.Nil,
		},
		"get index out of range": {
			l:   &value.ArrayTuple{value.Nil, value.Float(89.2), value.String("foo")},
			key: value.SmallInt(31),
			err: value.NewError(
				value.IndexErrorClass,
				"index 31 out of range: -3...3",
			),
		},
		"get negative index out of range": {
			l:   &value.ArrayTuple{value.Nil, value.Float(89.2), value.String("foo")},
			key: value.SmallInt(-31),
			err: value.NewError(
				value.IndexErrorClass,
				"index -31 out of range: -3...3",
			),
		},
		"get index -1 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.Float(89.2), value.String("foo")},
			key:  value.SmallInt(-1),
			want: value.String("foo"),
		},
		"get index -2 in a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Float(89.2), value.Nil, value.String("foo")},
			key:  value.SmallInt(-2),
			want: value.Nil,
		},
		"get index in the middle of a populated arrayTuple": {
			l:    &value.ArrayTuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.SmallInt(1),
			want: value.String("foo"),
		},
		"get uint8 index": {
			l:    &value.ArrayTuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.UInt8(1),
			want: value.String("foo"),
		},
		"get string index": {
			l:   &value.ArrayTuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key: value.String("lol"),
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			),
		},
		"get float index": {
			l:   &value.ArrayTuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key: value.Float(1),
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int`",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			want, err := tc.l.Subscript(tc.key)
			if diff := cmp.Diff(tc.err, err, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
			if tc.err != nil {
				return
			}
			if diff := cmp.Diff(tc.want, want, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
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
			want: "Std::ArrayTuple::Iterator{arrayTuple: %[], index: 0}",
		},
		"with one element": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3)},
				0,
			),
			want: `Std::ArrayTuple::Iterator{arrayTuple: %[3], index: 0}`,
		},
		"with elements": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3), value.String("foo")},
				1,
			),
			want: `Std::ArrayTuple::Iterator{arrayTuple: %[3, "foo"], index: 1}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.l.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
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
			err: value.ToSymbol("stop_iteration"),
		},
		"with two elements index 0": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3), value.Float(7.2)},
				0,
			),
			after: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3), value.Float(7.2)},
				1,
			),
			want: value.SmallInt(3),
		},
		"with two elements index 1": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3), value.Float(7.2)},
				1,
			),
			after: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3), value.Float(7.2)},
				2,
			),
			want: value.Float(7.2),
		},
		"with two elements index 2": {
			l: value.NewArrayTupleIteratorWithIndex(
				&value.ArrayTuple{value.SmallInt(3), value.Float(7.2)},
				2,
			),
			err: value.ToSymbol("stop_iteration"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.l.Next()
			if diff := cmp.Diff(tc.err, err); diff != "" {
				t.Fatalf(diff)
			}
			if tc.err != nil {
				return
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.after, tc.l); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
