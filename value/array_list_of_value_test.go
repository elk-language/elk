package value_test

import (
	"regexp"
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestArrayListOfValue_Grow(t *testing.T) {
	tests := map[string]struct {
		l        *value.ArrayListOfValue
		n        int
		after    *value.ArrayListOfValue
		capAfter int
	}{
		"grow a list with spare capacity": {
			l:        value.NewArrayListOfValue(10),
			n:        5,
			after:    &value.ArrayListOfValue{},
			capAfter: 15,
		},
		"grow an empty list": {
			l:        &value.ArrayListOfValue{},
			n:        5,
			after:    &value.ArrayListOfValue{},
			capAfter: 5,
		},
		"grow a list with elements": {
			l:        &value.ArrayListOfValue{value.SmallInt(2).ToValue()},
			n:        5,
			after:    &value.ArrayListOfValue{value.SmallInt(2).ToValue()},
			capAfter: 6,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.l.Grow(tc.n)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.capAfter, tc.l.Capacity(), opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayListOfValue_Concat(t *testing.T) {
	tests := map[string]struct {
		left  *value.ArrayListOfValue
		right value.Value
		want  *value.ArrayListOfValue
		err   value.Value
	}{
		"ArrayList + ArrayList => ArrayList": {
			left:  &value.ArrayListOfValue{value.SmallInt(2).ToValue()},
			right: value.Ref(&value.ArrayListOfValue{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))}),
			want:  &value.ArrayListOfValue{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
		},
		"ArrayList + ArrayTuple => ArrayList": {
			left:  &value.ArrayListOfValue{value.SmallInt(2).ToValue()},
			right: value.Ref(&value.ArrayTupleOfValue{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))}),
			want:  &value.ArrayListOfValue{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
		},
		"ArrayList + Int => TypeError": {
			left:  &value.ArrayListOfValue{value.SmallInt(2).ToValue()},
			right: value.Int8(5).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot concat 5i8 with list [2]`)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Concat(tc.right)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
			if !err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayListOfValue_Repeat(t *testing.T) {
	tests := map[string]struct {
		left  *value.ArrayListOfValue
		right value.Value
		want  *value.ArrayListOfValue
		err   value.Value
	}{
		"ArrayList * SmallInt => ArrayList": {
			left:  &value.ArrayListOfValue{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.SmallInt(3).ToValue(),
			want:  &value.ArrayListOfValue{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue(), value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
		},
		"ArrayList * -SmallInt => OutOfRangeError": {
			left:  &value.ArrayListOfValue{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.SmallInt(-3).ToValue(),
			err:   value.Ref(value.NewError(value.OutOfRangeErrorClass, `list repeat count cannot be negative: -3`)),
		},
		"ArrayList * 0 => ArrayList": {
			left:  &value.ArrayListOfValue{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.SmallInt(0).ToValue(),
			want:  &value.ArrayListOfValue{},
		},
		"ArrayList * BigInt => OutOfRangeError": {
			left:  &value.ArrayListOfValue{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.Ref(value.NewBigInt(3)),
			err:   value.Ref(value.NewError(value.OutOfRangeErrorClass, `list repeat count is too large 3`)),
		},
		"ArrayList * Int8 => TypeError": {
			left:  &value.ArrayListOfValue{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.Int8(3).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot repeat a list using 3i8`)),
		},
		"ArrayList * String => TypeError": {
			left:  &value.ArrayListOfValue{value.SmallInt(2).ToValue(), value.SmallInt(3).ToValue()},
			right: value.Ref(value.String("bar")),
			err:   value.Ref(value.NewError(value.TypeErrorClass, `cannot repeat a list using "bar"`)),
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

func TestArrayListOfValue_Inspect(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayListOfValue
		want string
	}{
		"empty": {
			l:    &value.ArrayListOfValue{},
			want: "[]",
		},
		"with one element": {
			l:    &value.ArrayListOfValue{value.SmallInt(3).ToValue()},
			want: `[3]`,
		},
		"with elements": {
			l:    &value.ArrayListOfValue{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
			want: `[3, "foo"]`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.l.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayListOfValue_Expand(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayListOfValue
		new  int
		want *value.ArrayListOfValue
	}{
		"add 3 to an empty list": {
			l:    &value.ArrayListOfValue{},
			new:  3,
			want: &value.ArrayListOfValue{value.Nil, value.Nil, value.Nil},
		},
		"add 0 to an empty list": {
			l:    &value.ArrayListOfValue{},
			new:  0,
			want: &value.ArrayListOfValue{},
		},
		"add 2 to a filled list": {
			l:    &value.ArrayListOfValue{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue()},
			new:  2,
			want: &value.ArrayListOfValue{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue(), value.Nil, value.Nil},
		},
		"add 0 to a filled list": {
			l:    &value.ArrayListOfValue{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue()},
			new:  0,
			want: &value.ArrayListOfValue{value.SmallInt(-3).ToValue(), value.Float(10.5).ToValue()},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.l.Expand(tc.new)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, tc.l, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayListOfValue_Append(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayListOfValue
		val  value.Value
		want *value.ArrayListOfValue
	}{
		"append to an empty list": {
			l:    &value.ArrayListOfValue{},
			val:  value.SmallInt(3).ToValue(),
			want: &value.ArrayListOfValue{value.SmallInt(3).ToValue()},
		},
		"append to a filled list": {
			l:    &value.ArrayListOfValue{value.SmallInt(3).ToValue()},
			val:  value.Ref(value.String("foo")),
			want: &value.ArrayListOfValue{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.l.Append(tc.val)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, tc.l, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayListOfValue_SubscriptSet(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayListOfValue
		key  value.Value
		val  value.Value
		want *value.ArrayListOfValue
		err  value.Value
	}{
		"set index 0 in an empty list": {
			l:    &value.ArrayListOfValue{},
			key:  value.SmallInt(0).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayListOfValue{},
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			)),
		},
		"set index 0 in a populated list": {
			l:    &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(0).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayListOfValue{value.SmallInt(25).ToValue(), value.Ref(value.String("foo"))},
		},
		"set index -1 in a populated list": {
			l:    &value.ArrayListOfValue{value.Nil, value.Float(89.2).ToValue(), value.Ref(value.String("foo"))},
			key:  value.SmallInt(-1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayListOfValue{value.Nil, value.Float(89.2).ToValue(), value.SmallInt(25).ToValue()},
		},
		"set index -2 in a populated list": {
			l:    &value.ArrayListOfValue{value.Float(89.2).ToValue(), value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(-2).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayListOfValue{value.Float(89.2).ToValue(), value.SmallInt(25).ToValue(), value.Ref(value.String("foo"))},
		},
		"set index in the middle of a populated list": {
			l:    &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.SmallInt(1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayListOfValue{value.Nil, value.SmallInt(25).ToValue(), value.Float(21.37).ToValue()},
		},
		"set uint8 index": {
			l:    &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.UInt8(1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayListOfValue{value.Nil, value.SmallInt(25).ToValue(), value.Float(21.37).ToValue()},
		},
		"set string index": {
			l:    &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.Ref(value.String("lol")),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			)),
		},
		"set float index": {
			l:    &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.Float(1).ToValue(),
			val:  value.SmallInt(25).ToValue(),
			want: &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
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

func TestArrayListOfValue_Subscript(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayListOfValue
		key  value.Value
		want value.Value
		err  value.Value
	}{
		"get index 0 in an empty list": {
			l:   &value.ArrayListOfValue{},
			key: value.SmallInt(0).ToValue(),
			err: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 0 out of range: 0...0",
			)),
		},
		"get index 0 in a populated list": {
			l:    &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(0).ToValue(),
			want: value.Nil,
		},
		"get index -1 in a populated list": {
			l:    &value.ArrayListOfValue{value.Nil, value.Float(89.2).ToValue(), value.Ref(value.String("foo"))},
			key:  value.SmallInt(-1).ToValue(),
			want: value.Ref(value.String("foo")),
		},
		"get index -2 in a populated list": {
			l:    &value.ArrayListOfValue{value.Float(89.2).ToValue(), value.Nil, value.Ref(value.String("foo"))},
			key:  value.SmallInt(-2).ToValue(),
			want: value.Nil,
		},
		"get index in the middle of a populated list": {
			l:    &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.SmallInt(1).ToValue(),
			want: value.Ref(value.String("foo")),
		},
		"get uint8 index": {
			l:    &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key:  value.UInt8(1).ToValue(),
			want: value.Ref(value.String("foo")),
		},
		"get string index": {
			l:   &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
			key: value.Ref(value.String("lol")),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			)),
		},
		"get float index": {
			l:   &value.ArrayListOfValue{value.Nil, value.Ref(value.String("foo")), value.Float(21.37).ToValue()},
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
				t.Fatal(diff)
			}
			if !tc.err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.want, want, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayListOfValue_Length(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayListOfValue
		want int
	}{
		"empty list": {
			l:    &value.ArrayListOfValue{},
			want: 0,
		},
		"one element": {
			l:    &value.ArrayListOfValue{value.SmallInt(3).ToValue()},
			want: 1,
		},
		"5 elements": {
			l: &value.ArrayListOfValue{
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
				t.Fatal(diff)
			}
		})
	}
}

func TestArrayListOfValueIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		l    *value.ArrayListOfValueIterator
		want string
	}{
		"empty": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayListOfValue{},
				0,
			),
			want: `Std::ArrayListOfValue::Iterator\{&: 0x[[:xdigit:]]{4,12}, list: \[\], index: 0\}`,
		},
		"with one element": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayListOfValue{value.SmallInt(3).ToValue()},
				0,
			),
			want: `Std::ArrayListOfValue::Iterator\{&: 0x[[:xdigit:]]{4,12}, list: \[3\], index: 0\}`,
		},
		"with elements": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayListOfValue{value.SmallInt(3).ToValue(), value.Ref(value.String("foo"))},
				1,
			),
			want: `Std::ArrayListOfValue::Iterator\{&: 0x[[:xdigit:]]{4,12}, list: \[3, "foo"\], index: 1\}`,
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

func TestArrayListOfValueIterator_Next(t *testing.T) {
	tests := map[string]struct {
		l     *value.ArrayListOfValueIterator
		after *value.ArrayListOfValueIterator
		want  value.Value
		err   value.Value
	}{
		"empty": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayListOfValue{},
				0,
			),
			after: value.NewArrayListIteratorWithIndex(
				&value.ArrayListOfValue{},
				0,
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
		"with two elements index 0": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayListOfValue{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				0,
			),
			after: value.NewArrayListIteratorWithIndex(
				&value.ArrayListOfValue{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				1,
			),
			want: value.SmallInt(3).ToValue(),
		},
		"with two elements index 1": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayListOfValue{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				1,
			),
			after: value.NewArrayListIteratorWithIndex(
				&value.ArrayListOfValue{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				2,
			),
			want: value.Float(7.2).ToValue(),
		},
		"with two elements index 2": {
			l: value.NewArrayListIteratorWithIndex(
				&value.ArrayListOfValue{value.SmallInt(3).ToValue(), value.Float(7.2).ToValue()},
				2,
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.l.NextValue()
			opts := comparer.Options()
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
			if !tc.err.IsUndefined() {
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
