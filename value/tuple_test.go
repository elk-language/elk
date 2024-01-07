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
