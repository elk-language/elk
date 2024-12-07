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
			l:        &value.ArrayList{value.SmallInt(2).ToValue()},
			n:        5,
			after:    &value.ArrayList{value.SmallInt(2).ToValue()},
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
		err   value.Value
	}{
		"ArrayList + ArrayList => ArrayList": {
			left:  &value.ArrayList{value.SmallInt(2).ToValue()},
			right: value.Ref(&value.ArrayList{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))}),
			want:  value.Ref(&value.ArrayList{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))}),
		},
		"ArrayList + ArrayTuple => ArrayList": {
			left:  &value.ArrayList{value.SmallInt(2).ToValue()},
			right: value.Ref(&value.ArrayTuple{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))}),
			want:  value.Ref(&value.ArrayList{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))}),
		},
		"ArrayList + Int => TypeError": {
			left:  &value.ArrayList{value.SmallInt(2).ToValue()},
			right: value.Int8(5).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot concat 5i8 with list [2]`)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Concat(tc.right)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if !err.IsUndefined() {
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
		err   value.Value
	}{
		"ArrayList * SmallInt => ArrayList": {
			left:  &value.ArrayList{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.SmallInt(3).ToValue(),
			want:  &value.ArrayList{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
		},
		"ArrayList * -SmallInt => OutOfRangeError": {
			left:  &value.ArrayList{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.SmallInt(-3).ToValue(),
			err:   value.Ref(value.NewError(value.OutOfRangeErrorClass, `list repeat count cannot be negative: -3`)),
		},
		"ArrayList * 0 => ArrayList": {
			left:  &value.ArrayList{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.SmallInt(0).ToValue(),
			want:  &value.ArrayList{},
		},
		"ArrayList * BigInt => OutOfRangeError": {
			left:  &value.ArrayList{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.Ref(value.NewBigInt(3)),
			err:   value.Ref(value.NewError(value.OutOfRangeErrorClass, `list repeat count is too large 3`)),
		},
		"ArrayList * Int8 => TypeError": {
			left:  &value.ArrayList{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.Int8(3).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot repeat a list using 3i8`)),
		},
		"ArrayList * String => TypeError": {
			left:  &value.ArrayList{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.Ref(value.String("bar")),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot repeat a list using "bar"`)),
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
			l:    &value.ArrayList{value.SmallInt(3).ToValue()},
			want: `[3]`,
		},
		"with elements": {
			l:    &value.ArrayList{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
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
			l:    &value.ArrayList{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue()},
			new:  2,
			want: &value.ArrayList{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue(), value.Nil, value.Nil},
		},
		"add 0 to a filled list": {
			l:    &value.ArrayList{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue()},
			new:  0,
			want: &value.ArrayList{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue()},
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
			val:  value.SmallInt(3).ToValue(),
			want: &value.ArrayList{value.SmallInt(3).ToValue()},
		},
		"append to a filled list": {
			l:    &value.ArrayList{value.SmallInt(3).ToValue()},
			val:  value.Ref(value.String("foo")),
			want: &value.ArrayList{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
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
		err  value.Value
	}{
		"set index 0 in an empty list": {
			l:    &value.ArrayList{},
			key:  value.SmallInt(0).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayList{},
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			)),
		},
		"set index 0 in a populated list": {
			l:    &value.ArrayList{value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(0).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayList{value.SmallInt(25).ToValue(), value.Ref(value.String("foo"))},
		},
		"set index -1 in a populated list": {
			l:    &value.ArrayList{value.Nil, value.Float(89.2).ToValue(), value.Ref(value.String("foo"))},
			key:  value.SmallInt(-1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayList{value.Nil, value.Float(89.2).ToValue(), value.SmallInt(25).ToValue()},
		},
		"set index -2 in a populated list": {
			l:    &value.ArrayList{value.Float(89.2).ToValue(), value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(-2).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayList{value.Float(89.2).ToValue(), value.SmallInt(25).ToValue(), value.Ref(value.String("foo"))},
		},
		"set index in the middle of a populated list": {
			l:    &value.ArrayList{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.SmallInt(1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayList{value.Nil, value.SmallInt(25).ToValue(), value.Float(21.37).ToValue()},
		},
		"set uint8 index": {
			l:    &value.ArrayList{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.UInt8(1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayList{value.Nil, value.SmallInt(25).ToValue(), value.Float(21.37).ToValue()},
		},
		"set string index": {
			l:    &value.ArrayList{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.Ref(value.String("lol")),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayList{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			)),
		},
		"set float index": {
			l:    &value.ArrayList{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.Float(1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayList{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int`",
			)),
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
		err  value.Value
	}{
		"get index 0 in an empty list": {
			l:   &value.ArrayList{},
			key: value.SmallInt(0).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			)),
		},
		"get index 0 in a populated list": {
			l:    &value.ArrayList{value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(0).ToValue(),
			want: value.Nil,
		},
		"get index -1 in a populated list": {
			l:    &value.ArrayList{value.Nil, value.Float(89.2).ToValue(), value.Ref(value.String("foo"))},
			key:  value.SmallInt(-1).ToValue(),
			want: value.Ref(value.String("foo")),
		},
		"get index -2 in a populated list": {
			l:    &value.ArrayList{value.Float(89.2).ToValue(), value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(-2).ToValue(),
			want: value.Nil,
		},
		"get index in the middle of a populated list": {
			l:    &value.ArrayList{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.SmallInt(1).ToValue(),
			want: value.Ref(value.String("foo")),
		},
		"get uint8 index": {
			l:    &value.ArrayList{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.UInt8(1).ToValue(),
			want: value.Ref(value.String("foo")),
		},
		"get string index": {
			l:   &value.ArrayList{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key: value.Ref(value.String("lol")),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			)),
		},
		"get float index": {
			l:   &value.ArrayList{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
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
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			if !tc.err.IsUndefined() {
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
			l:    &value.ArrayList{value.SmallInt(3).ToValue()},
			want: 1,
		},
		"5 elements": {
			l: &value.ArrayList{
				value.SmallInt(3).ToValue(),
				value.Nil,
				value.Float(4.5).ToValue(),
				value.Ref(value.String("bar")),
				value.Ref(value.String("foo")),
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
				&value.ArrayList{value.SmallInt(3).ToValue()},
				0,
			),
			want: `Std::ArrayList::Iterator{list: [3], index: 0}`,
		},
		"with elements": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
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
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
		"with two elements index 0": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				0,
			),
			after: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				1,
			),
			want: value.SmallInt(3).ToValue(),
		},
		"with two elements index 1": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				1,
			),
			after: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				2,
			),
			want: value.Float(7.2).ToValue(),
		},
		"with two elements index 2": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayList{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				2,
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.l.Next()
			if diff := cmp.Diff(tc.err, err); diff != "" {
				t.Fatalf(diff)
			}
			if !tc.err.IsUndefined() {
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
