package vm_test

import (
	"testing"

	"github.com/elk-language/elk/value"
)

func TestVMSource_TupleLiteral(t *testing.T) {
	tests := sourceTestTable{
		"empty tuple literal": {
			source:       `%[]`,
			wantStackTop: &value.Tuple{},
		},
		"static tuple literal": {
			source: `%[1, 2.5, :foo]`,
			wantStackTop: &value.Tuple{
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("foo"),
			},
		},
		"nested static tuple literal": {
			source: `%[1, 2.5, %["bar", %[]], %[:foo]]`,
			wantStackTop: &value.Tuple{
				value.SmallInt(1),
				value.Float(2.5),
				&value.Tuple{
					value.String("bar"),
					&value.Tuple{},
				},
				&value.Tuple{
					value.ToSymbol("foo"),
				},
			},
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				%[1, 2.5, foo, :bar]
			`,
			wantStackTop: &value.Tuple{
				value.SmallInt(1),
				value.Float(2.5),
				value.String("foo var"),
				value.ToSymbol("bar"),
			},
		},
		"starts with dynamic elements": {
			source: `
				foo := "foo var"
				%[foo, 1, 2.5, :bar]
			`,
			wantStackTop: &value.Tuple{
				value.String("foo var"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"with falsy if": {
			source: `
				foo := nil
				%["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: &value.Tuple{
				value.String("awesome"),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"with truthy if": {
			source: `
				foo := 57
				%["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: &value.Tuple{
				value.String("awesome"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
