package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestList_Concat(t *testing.T) {
	tests := map[string]struct {
		left  *value.List
		right value.Value
		want  value.Value
		err   *value.Error
	}{
		"List + List => List": {
			left:  &value.List{value.SmallInt(2)},
			right: &value.List{value.SmallInt(3), value.String("foo")},
			want:  &value.List{value.SmallInt(2), value.SmallInt(3), value.String("foo")},
		},
		"List + ArrayTuple => List": {
			left:  &value.List{value.SmallInt(2)},
			right: &value.ArrayTuple{value.SmallInt(3), value.String("foo")},
			want:  &value.List{value.SmallInt(2), value.SmallInt(3), value.String("foo")},
		},
		"List + Int => TypeError": {
			left:  &value.List{value.SmallInt(2)},
			right: value.Int8(5),
			err:   value.NewError(value.TypeErrorClass, `cannot concat 5i8 with list [2]`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Concat(tc.right)
			opts := comparer.Comparer
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

func TestList_Repeat(t *testing.T) {
	tests := map[string]struct {
		left  *value.List
		right value.Value
		want  *value.List
		err   *value.Error
	}{
		"List * SmallInt => List": {
			left:  &value.List{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(3),
			want:  &value.List{value.SmallInt(2), value.SmallInt(3), value.SmallInt(2), value.SmallInt(3), value.SmallInt(2), value.SmallInt(3)},
		},
		"List * -SmallInt => OutOfRangeError": {
			left:  &value.List{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(-3),
			err:   value.NewError(value.OutOfRangeErrorClass, `list repeat count cannot be negative: -3`),
		},
		"List * 0 => List": {
			left:  &value.List{value.SmallInt(2), value.SmallInt(3)},
			right: value.SmallInt(0),
			want:  &value.List{},
		},
		"List * BigInt => OutOfRangeError": {
			left:  &value.List{value.SmallInt(2), value.SmallInt(3)},
			right: value.NewBigInt(3),
			err:   value.NewError(value.OutOfRangeErrorClass, `list repeat count is too large 3`),
		},
		"List * Int8 => TypeError": {
			left:  &value.List{value.SmallInt(2), value.SmallInt(3)},
			right: value.Int8(3),
			err:   value.NewError(value.TypeErrorClass, `cannot repeat a list using 3i8`),
		},
		"List * String => TypeError": {
			left:  &value.List{value.SmallInt(2), value.SmallInt(3)},
			right: value.String("bar"),
			err:   value.NewError(value.TypeErrorClass, `cannot repeat a list using "bar"`),
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

func TestList_Inspect(t *testing.T) {
	tests := map[string]struct {
		l    *value.List
		want string
	}{
		"empty": {
			l:    &value.List{},
			want: "[]",
		},
		"with one element": {
			l:    &value.List{value.SmallInt(3)},
			want: `[3]`,
		},
		"with elements": {
			l:    &value.List{value.SmallInt(3), value.String("foo")},
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

func TestList_Expand(t *testing.T) {
	tests := map[string]struct {
		l    *value.List
		new  int
		want *value.List
	}{
		"add 3 to an empty list": {
			l:    &value.List{},
			new:  3,
			want: &value.List{value.Nil, value.Nil, value.Nil},
		},
		"add 0 to an empty list": {
			l:    &value.List{},
			new:  0,
			want: &value.List{},
		},
		"add 2 to a filled list": {
			l:    &value.List{value.SmallInt(-3), value.Float(10.5)},
			new:  2,
			want: &value.List{value.SmallInt(-3), value.Float(10.5), value.Nil, value.Nil},
		},
		"add 0 to a filled list": {
			l:    &value.List{value.SmallInt(-3), value.Float(10.5)},
			new:  0,
			want: &value.List{value.SmallInt(-3), value.Float(10.5)},
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

func TestList_Append(t *testing.T) {
	tests := map[string]struct {
		l    *value.List
		val  value.Value
		want *value.List
	}{
		"append to an empty list": {
			l:    &value.List{},
			val:  value.SmallInt(3),
			want: &value.List{value.SmallInt(3)},
		},
		"append to a filled list": {
			l:    &value.List{value.SmallInt(3)},
			val:  value.String("foo"),
			want: &value.List{value.SmallInt(3), value.String("foo")},
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

func TestList_SubscriptSet(t *testing.T) {
	tests := map[string]struct {
		l    *value.List
		key  value.Value
		val  value.Value
		want *value.List
		err  *value.Error
	}{
		"set index 0 in an empty list": {
			l:    &value.List{},
			key:  value.SmallInt(0),
			val:  value.SmallInt(25),
			want: &value.List{},
			err: value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			),
		},
		"set index 0 in a populated list": {
			l:    &value.List{value.Nil, value.String("foo")},
			key:  value.SmallInt(0),
			val:  value.SmallInt(25),
			want: &value.List{value.SmallInt(25), value.String("foo")},
		},
		"set index -1 in a populated list": {
			l:    &value.List{value.Nil, value.Float(89.2), value.String("foo")},
			key:  value.SmallInt(-1),
			val:  value.SmallInt(25),
			want: &value.List{value.Nil, value.Float(89.2), value.SmallInt(25)},
		},
		"set index -2 in a populated list": {
			l:    &value.List{value.Float(89.2), value.Nil, value.String("foo")},
			key:  value.SmallInt(-2),
			val:  value.SmallInt(25),
			want: &value.List{value.Float(89.2), value.SmallInt(25), value.String("foo")},
		},
		"set index in the middle of a populated list": {
			l:    &value.List{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.SmallInt(1),
			val:  value.SmallInt(25),
			want: &value.List{value.Nil, value.SmallInt(25), value.Float(21.37)},
		},
		"set uint8 index": {
			l:    &value.List{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.UInt8(1),
			val:  value.SmallInt(25),
			want: &value.List{value.Nil, value.SmallInt(25), value.Float(21.37)},
		},
		"set string index": {
			l:    &value.List{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.String("lol"),
			val:  value.SmallInt(25),
			want: &value.List{value.Nil, value.String("foo"), value.Float(21.37)},
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			),
		},
		"set float index": {
			l:    &value.List{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.Float(1),
			val:  value.SmallInt(25),
			want: &value.List{value.Nil, value.String("foo"), value.Float(21.37)},
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

func TestList_Subscript(t *testing.T) {
	tests := map[string]struct {
		l    *value.List
		key  value.Value
		want value.Value
		err  *value.Error
	}{
		"get index 0 in an empty list": {
			l:   &value.List{},
			key: value.SmallInt(0),
			err: value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			),
		},
		"get index 0 in a populated list": {
			l:    &value.List{value.Nil, value.String("foo")},
			key:  value.SmallInt(0),
			want: value.Nil,
		},
		"get index -1 in a populated list": {
			l:    &value.List{value.Nil, value.Float(89.2), value.String("foo")},
			key:  value.SmallInt(-1),
			want: value.String("foo"),
		},
		"get index -2 in a populated list": {
			l:    &value.List{value.Float(89.2), value.Nil, value.String("foo")},
			key:  value.SmallInt(-2),
			want: value.Nil,
		},
		"get index in the middle of a populated list": {
			l:    &value.List{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.SmallInt(1),
			want: value.String("foo"),
		},
		"get uint8 index": {
			l:    &value.List{value.Nil, value.String("foo"), value.Float(21.37)},
			key:  value.UInt8(1),
			want: value.String("foo"),
		},
		"get string index": {
			l:   &value.List{value.Nil, value.String("foo"), value.Float(21.37)},
			key: value.String("lol"),
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			),
		},
		"get float index": {
			l:   &value.List{value.Nil, value.String("foo"), value.Float(21.37)},
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

func TestList_Length(t *testing.T) {
	tests := map[string]struct {
		l    *value.List
		want int
	}{
		"empty list": {
			l:    &value.List{},
			want: 0,
		},
		"one element": {
			l:    &value.List{value.SmallInt(3)},
			want: 1,
		},
		"5 elements": {
			l: &value.List{
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

func TestListIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		l    *value.ListIterator
		want string
	}{
		"empty": {
			l: value.NewListIteratorWithIndex(
				&value.List{},
				0,
			),
			want: "Std::List::Iterator{list: [], index: 0}",
		},
		"with one element": {
			l: value.NewListIteratorWithIndex(
				&value.List{value.SmallInt(3)},
				0,
			),
			want: `Std::List::Iterator{list: [3], index: 0}`,
		},
		"with elements": {
			l: value.NewListIteratorWithIndex(
				&value.List{value.SmallInt(3), value.String("foo")},
				1,
			),
			want: `Std::List::Iterator{list: [3, "foo"], index: 1}`,
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

func TestListIterator_Next(t *testing.T) {
	tests := map[string]struct {
		l     *value.ListIterator
		after *value.ListIterator
		want  value.Value
		err   value.Value
	}{
		"empty": {
			l: value.NewListIteratorWithIndex(
				&value.List{},
				0,
			),
			after: value.NewListIteratorWithIndex(
				&value.List{},
				0,
			),
			err: value.ToSymbol("stop_iteration"),
		},
		"with two elements index 0": {
			l: value.NewListIteratorWithIndex(
				&value.List{value.SmallInt(3), value.Float(7.2)},
				0,
			),
			after: value.NewListIteratorWithIndex(
				&value.List{value.SmallInt(3), value.Float(7.2)},
				1,
			),
			want: value.SmallInt(3),
		},
		"with two elements index 1": {
			l: value.NewListIteratorWithIndex(
				&value.List{value.SmallInt(3), value.Float(7.2)},
				1,
			),
			after: value.NewListIteratorWithIndex(
				&value.List{value.SmallInt(3), value.Float(7.2)},
				2,
			),
			want: value.Float(7.2),
		},
		"with two elements index 2": {
			l: value.NewListIteratorWithIndex(
				&value.List{value.SmallInt(3), value.Float(7.2)},
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
