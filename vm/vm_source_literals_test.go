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
		"dynamic tuple literal": {
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
