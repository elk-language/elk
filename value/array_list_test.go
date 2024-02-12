package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestArrayList_Grow(t *testing.T) {
	tests := map[string]struct {
		l        *value.ArrayList
		n        int
		after    *value.ArrayList
		capAfter int
	}{
		"grow a list with spare capacity": {
			l:        value.NewArrayList(10),
			n:        5,
			after:    &value.ArrayList{},
			capAfter: 15,
		},
		"grow an empty list": {
			l:        &value.ArrayList{},
			n:        5,
			after:    &value.ArrayList{},
			capAfter: 5,
		},
		"grow a list with elements": {
			l:        &value.ArrayList{value.SmallInt(2)},
			n:        5,
			after:    &value.ArrayList{value.SmallInt(2)},
			capAfter: 6,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.l.Grow(tc.n)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.capAfter, tc.l.Capacity(), opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestArrayList_Concat(t *testing.T) {
	tests := map[string]struct {
		left  *value.ArrayList
		right value.Value
		want  value.Value
		err   *value.Error
	}{
		"ArrayList + ArrayList => ArrayList": {
			left:  &value.ArrayList{value.SmallInt(2)},
			right: &value.ArrayList{value.SmallInt(3), value.String("foo")},
			want:  &value.ArrayList{value.SmallInt(2), value.SmallInt(3), value.String("foo")},
		},
		"ArrayList + ArrayTuple => ArrayList": {
			left:  &value.ArrayList{value.SmallInt(2)},
			right: &value.ArrayTuple{value.SmallInt(3), value.String("foo")},
			want:  &value.ArrayList{value.SmallInt(2), value.SmallInt(3), value.String("foo")},
		},
		"ArrayList + Int => TypeError": {
			left:  &value.ArrayList{value.SmallInt(2)},
			right: value.Int8(5),
			err:   value.NewError(value.TypeErrorClass, `cannot concat 5i8 with list [2]`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Concat(tc.right)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestArrayList_Repeat(t *testing.T) {
	tests := map[string]struct {
		left  *value.ArrayList
		right value.Value
		want  *value.ArrayList
		err   *value.Error
	}{
		"ArrayList * SmallInt => ArrayList": {
			left:  &value.ArrayList{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(3),
			want:  &value.ArrayList{value.SmallInt(2), value.SmallInt(3), value.SmallInt(2), value.SmallInt(3), value.SmallInt(2), value.SmallInt(3)},
		},
		"ArrayList * -SmallInt => OutOfRangeError": {
			left:  &value.ArrayList{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(-3),
			err:   value.NewError(value.OutOfRangeErrorClass, `list repeat count cannot be negative: -3`),
		},
		"ArrayList * 0 => ArrayList": {
			left:  &value.ArrayList{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(0),
			want:  &value.ArrayList{},
		},
		"ArrayList * BigInt => OutOfRangeError": {
			left:  &value.ArrayList{value.SmallInt(2), value.SmallInt(3)},
			right: value.NewBigInt(3),
			err:   value.NewError(value.OutOfRangeErrorClass, `list repeat count is too large 3`),
		},
		"ArrayList * Int8 => TypeError": {
			left:  &value.ArrayList{value.SmallInt(2), value.SmallInt(3)},
			right: value.Int8(3),
			err:   value.NewError(value.TypeErrorClass, `cannot repeat a list using 3i8`),
		},
		"ArrayList * String => TypeError": {
			left:  &value.ArrayList{value.SmallInt(2), value.SmallInt(3)},
			right: value.String("bar"),
			err:   value.NewError(value.TypeErrorClass, `cannot repeat a list using "bar"`),
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

func TestArrayList_Inspect(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayList
		want string
	}{
		"empty": {
			l:    &value.ArrayList{},
			want: "[]",
		},
		"with one element": {
			l:    &value.ArrayList{value.SmallInt(3)},
			want: `[3]`,
		},
		"with elements": {
			l:    &value.ArrayList{value.SmallInt(3), value.String("foo")},
			want: `[3, "foo"]`,
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

func TestArrayList_Expand(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayList
		new  int
		want *value.ArrayList
	}{
		"add 3 to an empty list": {
			l:    &value.ArrayList{},
			new:  3,
			want: &value.ArrayList{value.Nil, value.Nil, value.Nil},
		},
		"add 0 to an empty list": {
			l:    &value.ArrayList{},
			new:  0,
			want: &value.ArrayList{},
		},
		"add 2 to a filled list": {
			l:    &value.ArrayList{value.SmallInt(-3), value.Float(10.5)},
			new:  2,
			want: &value.ArrayList{value.SmallInt(-3), value.Float(10.5), value.Nil, value.Nil},
		},
		"add 0 to a filled list": {
			l:    &value.ArrayList{value.SmallInt(-3), value.Float(10.5)},
			new:  0,
			want: &value.ArrayList{value.SmallInt(-3), value.Float(10.5)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.l.Expand(tc.new)
			if diff := cmp.Diff(tc.want, tc.l); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestArrayList_Append(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayList
		val  value.Value
		want *value.ArrayList
	}{
		"append to an empty list": {
			l:    &value.ArrayList{},
			val:  value.SmallInt(3),
			want: &value.ArrayList{value.SmallInt(3)},
		},
		"append to a filled list": {
			l:    &value.ArrayList{value.SmallInt(3)},
			val:  value.String("foo"),
			want: &value.ArrayList{value.SmallInt(3), value.String("foo")},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.l.Append(tc.val)
			if diff := cmp.Diff(tc.want, tc.l); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestArrayList_SubscriptSet(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayList
		key  value.Value
		val  value.Value
		want *value.ArrayList
		err  *value.Error
	}{
		"set index 0 in an empty list": {
			l:    &value.ArrayList{},
			key:  value.SmallInt(0),
			val:  value.SmallInt(25),
			want: &value.ArrayList{},
			err: value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			),
		},
		"set index 0 in a populated list": {
			l:    &value.ArrayList{value.Nil, value.String("foo")},
			key:  value.SmallInt(0),
			val:  value.SmallInt(25),
			want: &value.ArrayList{value.SmallInt(25), value.String("foo")},
		},
		"set index -1 in a populated list": {
			l:    &value.ArrayList{value.Nil, value.Float(89.2), value.String("foo")},
			key:  value.SmallInt(-1),
			val:  value.SmallInt(25),
			want: &value.ArrayList{value.Nil, value.Float(89.2), value.SmallInt(25)},
		},
		"set index -2 in a populated list": {
			l:    &value.ArrayList{value.Float(89.2), value.Nil, value.String("foo")},
			key:  value.SmallInt(-2),
			val:  value.SmallInt(25),
			want: &value.ArrayList{value.Float(89.2), value.SmallInt(25), value.String("foo")},
		},
		"set index in the middle of a populated list": {
			l:    &value.ArrayList{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.SmallInt(1),
			val:  value.SmallInt(25),
			want: &value.ArrayList{value.Nil, value.SmallInt(25), value.Float(21.37)},
		},
		"set uint8 index": {
			l:    &value.ArrayList{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.UInt8(1),
			val:  value.SmallInt(25),
			want: &value.ArrayList{value.Nil, value.SmallInt(25), value.Float(21.37)},
		},
		"set string index": {
			l:    &value.ArrayList{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.String("lol"),
			val:  value.SmallInt(25),
			want: &value.ArrayList{value.Nil, value.String("foo"), value.Float(21.37)},
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			),
		},
		"set float index": {
			l:    &value.ArrayList{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.Float(1),
			val:  value.SmallInt(25),
			want: &value.ArrayList{value.Nil, value.String("foo"), value.Float(21.37)},
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int`",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.l.SubscriptSet(tc.key, tc.val)
			if diff := cmp.Diff(tc.want, tc.l, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestArrayList_Subscript(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayList
		key  value.Value
		want value.Value
		err  *value.Error
	}{
		"get index 0 in an empty list": {
			l:   &value.ArrayList{},
			key: value.SmallInt(0),
			err: value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			),
		},
		"get index 0 in a populated list": {
			l:    &value.ArrayList{value.Nil, value.String("foo")},
			key:  value.SmallInt(0),
			want: value.Nil,
		},
		"get index -1 in a populated list": {
			l:    &value.ArrayList{value.Nil, value.Float(89.2), value.String("foo")},
			key:  value.SmallInt(-1),
			want: value.String("foo"),
		},
		"get index -2 in a populated list": {
			l:    &value.ArrayList{value.Float(89.2), value.Nil, value.String("foo")},
			key:  value.SmallInt(-2),
			want: value.Nil,
		},
		"get index in the middle of a populated list": {
			l:    &value.ArrayList{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.SmallInt(1),
			want: value.String("foo"),
		},
		"get uint8 index": {
			l:    &value.ArrayList{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.UInt8(1),
			want: value.String("foo"),
		},
		"get string index": {
			l:   &value.ArrayList{value.Nil, value.String("foo"), value.Float(21.37)},
			key: value.String("lol"),
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			),
		},
		"get float index": {
			l:   &value.ArrayList{value.Nil, value.String("foo"), value.Float(21.37)},
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
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			if tc.err != nil {
				return
			}
			if diff := cmp.Diff(tc.want, want, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestArrayList_Length(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayList
		want int
	}{
		"empty list": {
			l:    &value.ArrayList{},
			want: 0,
		},
		"one element": {
			l:    &value.ArrayList{value.SmallInt(3)},
			want: 1,
		},
		"5 elements": {
			l: &value.ArrayList{
				value.SmallInt(3),
				value.Nil,
				value.Float(4.5),
				value.String("bar"),
				value.String("foo"),
			},
			want: 5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.l.Length()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestArrayListIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayListIterator
		want string
	}{
		"empty": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{},
				0,
			),
			want: "Std::ArrayList::Iterator{list: [], index: 0}",
		},
		"with one element": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3)},
				0,
			),
			want: `Std::ArrayList::Iterator{list: [3], index: 0}`,
		},
		"with elements": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3), value.String("foo")},
				1,
			),
			want: `Std::ArrayList::Iterator{list: [3, "foo"], index: 1}`,
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

func TestArrayListIterator_Next(t *testing.T) {
	tests := map[string]struct {
		l     *value.ArrayListIterator
		after *value.ArrayListIterator
		want  value.Value
		err   value.Value
	}{
		"empty": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{},
				0,
			),
			after: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{},
				0,
			),
			err: value.ToSymbol("stop_iteration"),
		},
		"with two elements index 0": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3), value.Float(7.2)},
				0,
			),
			after: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3), value.Float(7.2)},
				1,
			),
			want: value.SmallInt(3),
		},
		"with two elements index 1": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3), value.Float(7.2)},
				1,
			),
			after: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3), value.Float(7.2)},
				2,
			),
			want: value.Float(7.2),
		},
		"with two elements index 2": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3), value.Float(7.2)},
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
