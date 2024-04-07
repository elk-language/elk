package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestArrayListContains(t *testing.T) {
	tests := map[string]struct {
		list     *value.ArrayList
		val      value.Value
		contains bool
		err      value.Value
	}{
		"empty list": {
			list:     &value.ArrayList{},
			val:      value.SmallInt(5),
			contains: false,
		},
		"coercible elements": {
			list:     &value.ArrayList{value.String("foo"), value.Float(5)},
			val:      value.SmallInt(5),
			contains: false,
		},
		"has the value": {
			list:     &value.ArrayList{value.String("foo"), value.SmallInt(5), value.Float(9.3)},
			val:      value.SmallInt(5),
			contains: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			contains, err := vm.ArrayListContains(v, tc.list, tc.val)
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tc.contains, contains, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestArrayListEqual(t *testing.T) {
	tests := map[string]struct {
		list  *value.ArrayList
		other *value.ArrayList
		equal bool
		err   value.Value
	}{
		"two identical lists": {
			list:  &value.ArrayList{value.String("foo"), value.Float(5)},
			other: &value.ArrayList{value.String("foo"), value.Float(5)},
			equal: true,
		},
		"different length": {
			list:  &value.ArrayList{value.String("foo"), value.Float(5)},
			other: &value.ArrayList{value.String("foo"), value.Float(5), value.Nil},
			equal: false,
		},
		"the same values of different types": {
			list:  &value.ArrayList{value.String("foo"), value.SmallInt(5)},
			other: &value.ArrayList{value.String("foo"), value.Float(5)},
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.ArrayListEqual(v, tc.list, tc.other)
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tc.equal, equal, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
