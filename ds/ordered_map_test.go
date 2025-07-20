package ds

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrderedMap_GetTest(t *testing.T) {
	tests := map[string]struct {
		m    *OrderedMap[string, int]
		key  string
		want int
	}{
		"get existing value": {
			m: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
			),
			key:  "bar",
			want: -25,
		},
		"get nonexisting value": {
			m: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
			),
			key:  "baz",
			want: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.m.Get(tc.key)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestOrderedMap_GetOkTest(t *testing.T) {
	tests := map[string]struct {
		m         *OrderedMap[string, int]
		key       string
		wantValue int
		wantOk    bool
	}{
		"get existing value": {
			m: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
			),
			key:       "bar",
			wantValue: -25,
			wantOk:    true,
		},
		"get nonexisting value": {
			m: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
			),
			key:       "baz",
			wantValue: 0,
			wantOk:    false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotValue, gotOk := tc.m.GetOk(tc.key)
			if diff := cmp.Diff(tc.wantValue, gotValue); diff != "" {
				t.Error(diff)
			}
			if diff := cmp.Diff(tc.wantOk, gotOk); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestOrderedMap_Set(t *testing.T) {
	tests := map[string]struct {
		m       *OrderedMap[string, int]
		key     string
		value   int
		wantMap *OrderedMap[string, int]
	}{
		"set existing value": {
			m: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
			),
			key:   "bar",
			value: 2,
			wantMap: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", 2),
			),
		},
		"set new value": {
			m: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
			),
			key:   "lol",
			value: 10,
			wantMap: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
				MakePair("lol", 10),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.m.Set(tc.key, tc.value)

			opts := cmp.Options{
				cmp.AllowUnexported(OrderedMap[string, int]{}),
			}
			if diff := cmp.Diff(tc.wantMap, tc.m, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestOrderedMap_Delete(t *testing.T) {
	tests := map[string]struct {
		m       *OrderedMap[string, int]
		key     string
		want    bool
		wantMap *OrderedMap[string, int]
	}{
		"delete existing value": {
			m: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
				MakePair("baz", 5),
			),
			key:  "bar",
			want: true,
			wantMap: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("baz", 5),
			),
		},
		"delete nonexisting value": {
			m: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
			),
			key:  "lol",
			want: false,
			wantMap: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.m.Delete(tc.key)

			opts := cmp.Options{
				cmp.AllowUnexported(OrderedMap[string, int]{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Error(diff)
			}
			if diff := cmp.Diff(tc.wantMap, tc.m, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestOrderedMap_All(t *testing.T) {
	tests := map[string]struct {
		m       *OrderedMap[string, int]
		wantAll []Pair[string, int]
	}{
		"iterate in insertion order": {
			m: NewOrderedMapWithPairs(
				MakePair("foo", 1),
				MakePair("bar", -25),
				MakePair("baz", 5),
			),
			wantAll: []Pair[string, int]{
				MakePair("foo", 1),
				MakePair("bar", -25),
				MakePair("baz", 5),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var gotAll []Pair[string, int]
			for key, value := range tc.m.All() {
				gotAll = append(gotAll, MakePair(key, value))
			}

			opts := cmp.Options{
				cmp.AllowUnexported(OrderedMap[string, int]{}),
			}
			if diff := cmp.Diff(tc.wantAll, gotAll, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}
