package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestTuple_Concat(t *testing.T) {
	tests := map[string]struct {
		left  *value.Tuple
		right value.Value
		want  value.Value
		err   *value.Error
	}{
		"Tuple + Tuple => Tuple": {
			left:  &value.Tuple{value.SmallInt(2)},
			right: &value.Tuple{value.SmallInt(3), value.String("foo")},
			want:  &value.Tuple{value.SmallInt(2), value.SmallInt(3), value.String("foo")},
		},
		"Tuple + List => List": {
			left:  &value.Tuple{value.SmallInt(2)},
			right: &value.List{value.SmallInt(3), value.String("foo")},
			want:  &value.List{value.SmallInt(2), value.SmallInt(3), value.String("foo")},
		},
		"Tuple + Int => TypeError": {
			left:  &value.Tuple{value.SmallInt(2)},
			right: value.Int8(5),
			err:   value.NewError(value.TypeErrorClass, `cannot concat 5i8 with tuple %[2]`),
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

func TestTuple_Repeat(t *testing.T) {
	tests := map[string]struct {
		left  *value.Tuple
		right value.Value
		want  *value.Tuple
		err   *value.Error
	}{
		"Tuple * SmallInt => Tuple": {
			left:  &value.Tuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(3),
			want:  &value.Tuple{value.SmallInt(2), value.SmallInt(3), value.SmallInt(2), value.SmallInt(3), value.SmallInt(2), value.SmallInt(3)},
		},
		"Tuple * -SmallInt => OutOfRangeError": {
			left:  &value.Tuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(-3),
			err:   value.NewError(value.OutOfRangeErrorClass, `tuple repeat count cannot be negative: -3`),
		},
		"Tuple * 0 => Tuple": {
			left:  &value.Tuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(0),
			want:  &value.Tuple{},
		},
		"Tuple * BigInt => OutOfRangeError": {
			left:  &value.Tuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.NewBigInt(3),
			err:   value.NewError(value.OutOfRangeErrorClass, `tuple repeat count is too large 3`),
		},
		"Tuple * Int8 => TypeError": {
			left:  &value.Tuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.Int8(3),
			err:   value.NewError(value.TypeErrorClass, `cannot repeat a tuple using 3i8`),
		},
		"Tuple * String => TypeError": {
			left:  &value.Tuple{value.SmallInt(2), value.SmallInt(3)},
			right: value.String("bar"),
			err:   value.NewError(value.TypeErrorClass, `cannot repeat a tuple using "bar"`),
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

func TestTuple_Inspect(t *testing.T) {
	tests := map[string]struct {
		t    *value.Tuple
		want string
	}{
		"empty": {
			t:    &value.Tuple{},
			want: "%[]",
		},
		"with one element": {
			t:    &value.Tuple{value.SmallInt(3)},
			want: `%[3]`,
		},
		"with elements": {
			t:    &value.Tuple{value.SmallInt(3), value.String("foo")},
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

func TestTuple_Expand(t *testing.T) {
	tests := map[string]struct {
		t    *value.Tuple
		new  int
		want *value.Tuple
	}{
		"add 3 to an empty list": {
			t:    &value.Tuple{},
			new:  3,
			want: &value.Tuple{value.Nil, value.Nil, value.Nil},
		},
		"add 0 to an empty list": {
			t:    &value.Tuple{},
			new:  0,
			want: &value.Tuple{},
		},
		"add 2 to a filled list": {
			t:    &value.Tuple{value.SmallInt(-3), value.Float(10.5)},
			new:  2,
			want: &value.Tuple{value.SmallInt(-3), value.Float(10.5), value.Nil, value.Nil},
		},
		"add 0 to a filled list": {
			t:    &value.Tuple{value.SmallInt(-3), value.Float(10.5)},
			new:  0,
			want: &value.Tuple{value.SmallInt(-3), value.Float(10.5)},
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

func TestTuple_Append(t *testing.T) {
	tests := map[string]struct {
		t    *value.Tuple
		val  value.Value
		want *value.Tuple
	}{
		"append to an empty list": {
			t:    &value.Tuple{},
			val:  value.SmallInt(3),
			want: &value.Tuple{value.SmallInt(3)},
		},
		"append to a filled list": {
			t:    &value.Tuple{value.SmallInt(3)},
			val:  value.String("foo"),
			want: &value.Tuple{value.SmallInt(3), value.String("foo")},
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

func TestTuple_SubscriptSet(t *testing.T) {
	tests := map[string]struct {
		l    *value.Tuple
		key  value.Value
		val  value.Value
		want *value.Tuple
		err  *value.Error
	}{
		"set index 0 in an empty tuple": {
			l:    &value.Tuple{},
			key:  value.SmallInt(0),
			val:  value.SmallInt(25),
			want: &value.Tuple{},
			err: value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			),
		},
		"set index 0 in a populated tuple": {
			l:    &value.Tuple{value.Nil, value.String("foo")},
			key:  value.SmallInt(0),
			val:  value.SmallInt(25),
			want: &value.Tuple{value.SmallInt(25), value.String("foo")},
		},
		"set index -1 in a populated tuple": {
			l:    &value.Tuple{value.Nil, value.Float(89.2), value.String("foo")},
			key:  value.SmallInt(-1),
			val:  value.SmallInt(25),
			want: &value.Tuple{value.Nil, value.Float(89.2), value.SmallInt(25)},
		},
		"set index -2 in a populated tuple": {
			l:    &value.Tuple{value.Float(89.2), value.Nil, value.String("foo")},
			key:  value.SmallInt(-2),
			val:  value.SmallInt(25),
			want: &value.Tuple{value.Float(89.2), value.SmallInt(25), value.String("foo")},
		},
		"set index in the middle of a populated tuple": {
			l:    &value.Tuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.SmallInt(1),
			val:  value.SmallInt(25),
			want: &value.Tuple{value.Nil, value.SmallInt(25), value.Float(21.37)},
		},
		"set uint8 index": {
			l:    &value.Tuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.UInt8(1),
			val:  value.SmallInt(25),
			want: &value.Tuple{value.Nil, value.SmallInt(25), value.Float(21.37)},
		},
		"set string index": {
			l:    &value.Tuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.String("lol"),
			val:  value.SmallInt(25),
			want: &value.Tuple{value.Nil, value.String("foo"), value.Float(21.37)},
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			),
		},
		"set float index": {
			l:    &value.Tuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.Float(1),
			val:  value.SmallInt(25),
			want: &value.Tuple{value.Nil, value.String("foo"), value.Float(21.37)},
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

func TestTuple_Subscript(t *testing.T) {
	tests := map[string]struct {
		l    *value.Tuple
		key  value.Value
		want value.Value
		err  *value.Error
	}{
		"get index 0 in an empty tuple": {
			l:   &value.Tuple{},
			key: value.SmallInt(0),
			err: value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			),
		},
		"get index 0 in a populated tuple": {
			l:    &value.Tuple{value.Nil, value.String("foo")},
			key:  value.SmallInt(0),
			want: value.Nil,
		},
		"get index -1 in a populated tuple": {
			l:    &value.Tuple{value.Nil, value.Float(89.2), value.String("foo")},
			key:  value.SmallInt(-1),
			want: value.String("foo"),
		},
		"get index -2 in a populated tuple": {
			l:    &value.Tuple{value.Float(89.2), value.Nil, value.String("foo")},
			key:  value.SmallInt(-2),
			want: value.Nil,
		},
		"get index in the middle of a populated tuple": {
			l:    &value.Tuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.SmallInt(1),
			want: value.String("foo"),
		},
		"get uint8 index": {
			l:    &value.Tuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.UInt8(1),
			want: value.String("foo"),
		},
		"get string index": {
			l:   &value.Tuple{value.Nil, value.String("foo"), value.Float(21.37)},
			key: value.String("lol"),
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			),
		},
		"get float index": {
			l:   &value.Tuple{value.Nil, value.String("foo"), value.Float(21.37)},
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

func TestTupleIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		l    *value.TupleIterator
		want string
	}{
		"empty": {
			l: value.NewTupleIteratorWithIndex(
				&value.Tuple{},
				0,
			),
			want: "Std::Tuple::Iterator{tuple: %[], index: 0}",
		},
		"with one element": {
			l: value.NewTupleIteratorWithIndex(
				&value.Tuple{value.SmallInt(3)},
				0,
			),
			want: `Std::Tuple::Iterator{tuple: %[3], index: 0}`,
		},
		"with elements": {
			l: value.NewTupleIteratorWithIndex(
				&value.Tuple{value.SmallInt(3), value.String("foo")},
				1,
			),
			want: `Std::Tuple::Iterator{tuple: %[3, "foo"], index: 1}`,
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

func TestTupleIterator_Next(t *testing.T) {
	tests := map[string]struct {
		l     *value.TupleIterator
		after *value.TupleIterator
		want  value.Value
		err   value.Value
	}{
		"empty": {
			l: value.NewTupleIteratorWithIndex(
				&value.Tuple{},
				0,
			),
			after: value.NewTupleIteratorWithIndex(
				&value.Tuple{},
				0,
			),
			err: value.ToSymbol("stop_iteration"),
		},
		"with two elements index 0": {
			l: value.NewTupleIteratorWithIndex(
				&value.Tuple{value.SmallInt(3), value.Float(7.2)},
				0,
			),
			after: value.NewTupleIteratorWithIndex(
				&value.Tuple{value.SmallInt(3), value.Float(7.2)},
				1,
			),
			want: value.SmallInt(3),
		},
		"with two elements index 1": {
			l: value.NewTupleIteratorWithIndex(
				&value.Tuple{value.SmallInt(3), value.Float(7.2)},
				1,
			),
			after: value.NewTupleIteratorWithIndex(
				&value.Tuple{value.SmallInt(3), value.Float(7.2)},
				2,
			),
			want: value.Float(7.2),
		},
		"with two elements index 2": {
			l: value.NewTupleIteratorWithIndex(
				&value.Tuple{value.SmallInt(3), value.Float(7.2)},
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
