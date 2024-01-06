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
		"with falsy unless": {
			source: `
				foo := nil
				%["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: &value.Tuple{
				value.String("awesome"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"with truthy unless": {
			source: `
				foo := true
				%["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: &value.Tuple{
				value.String("awesome"),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"static with indices": {
			source: `%["awesome", 5 => :foo, 2 => 8.3]`,
			wantStackTop: &value.Tuple{
				value.String("awesome"),
				value.Nil,
				value.Float(8.3),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo"),
			},
		},
		"static indices with dynamic elements": {
			source: `
			  foo := 3
				%["awesome", 5 => :foo, 2 => 8.3, foo]
			`,
			wantStackTop: &value.Tuple{
				value.String("awesome"),
				value.Nil,
				value.Float(8.3),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo"),
				value.SmallInt(3),
			},
		},
		"with dynamic elements and indices": {
			source: `
			  foo := 3
				%[foo, "awesome", 5 => :foo, 2 => 8.3]
			`,
			wantStackTop: &value.Tuple{
				value.SmallInt(3),
				value.String("awesome"),
				value.Float(8.3),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo"),
			},
		},
		"with dynamic indices": {
			source: `
			  foo := 3
				%[foo => :bar, "awesome"]
			`,
			wantStackTop: &value.Tuple{
				value.Nil,
				value.Nil,
				value.Nil,
				value.ToSymbol("bar"),
				value.String("awesome"),
			},
		},
		"with initial modifier": {
			source: `
			  foo := true
				%[3 if foo]
			`,
			wantStackTop: &value.Tuple{
				value.SmallInt(3),
			},
		},
		"with string index": {
			source: `
			  foo := "3"
				%[foo => :bar]
			`,
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			),
		},
		"with indices and if modifiers": {
			source: `
			  foo := "3"
				%[3 => :bar if foo]
			`,
			wantStackTop: &value.Tuple{
				value.Nil,
				value.Nil,
				value.Nil,
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

func TestVMSource_ListLiteral(t *testing.T) {
	tests := sourceTestTable{
		"empty tuple literal": {
			source:       `[]`,
			wantStackTop: &value.List{},
		},
		"static tuple literal": {
			source: `[1, 2.5, :foo]`,
			wantStackTop: &value.List{
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("foo"),
			},
		},
		"nested static tuple literal": {
			source: `[1, 2.5, ["bar", []], [:foo]]`,
			wantStackTop: &value.List{
				value.SmallInt(1),
				value.Float(2.5),
				&value.List{
					value.String("bar"),
					&value.List{},
				},
				&value.List{
					value.ToSymbol("foo"),
				},
			},
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				[1, 2.5, foo, :bar]
			`,
			wantStackTop: &value.List{
				value.SmallInt(1),
				value.Float(2.5),
				value.String("foo var"),
				value.ToSymbol("bar"),
			},
		},
		"starts with dynamic elements": {
			source: `
				foo := "foo var"
				[foo, 1, 2.5, :bar]
			`,
			wantStackTop: &value.List{
				value.String("foo var"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"with falsy if": {
			source: `
				foo := nil
				["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: &value.List{
				value.String("awesome"),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"with truthy if": {
			source: `
				foo := 57
				["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: &value.List{
				value.String("awesome"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"with falsy unless": {
			source: `
				foo := nil
				["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: &value.List{
				value.String("awesome"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"with truthy unless": {
			source: `
				foo := true
				["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: &value.List{
				value.String("awesome"),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"static with indices": {
			source: `["awesome", 5 => :foo, 2 => 8.3]`,
			wantStackTop: &value.List{
				value.String("awesome"),
				value.Nil,
				value.Float(8.3),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo"),
			},
		},
		"static indices with dynamic elements": {
			source: `
			  foo := 3
				["awesome", 5 => :foo, 2 => 8.3, foo]
			`,
			wantStackTop: &value.List{
				value.String("awesome"),
				value.Nil,
				value.Float(8.3),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo"),
				value.SmallInt(3),
			},
		},
		"with dynamic elements and indices": {
			source: `
			  foo := 3
				[foo, "awesome", 5 => :foo, 2 => 8.3]
			`,
			wantStackTop: &value.List{
				value.SmallInt(3),
				value.String("awesome"),
				value.Float(8.3),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo"),
			},
		},
		"with dynamic indices": {
			source: `
			  foo := 3
				[foo => :bar, "awesome"]
			`,
			wantStackTop: &value.List{
				value.Nil,
				value.Nil,
				value.Nil,
				value.ToSymbol("bar"),
				value.String("awesome"),
			},
		},
		"with initial modifier": {
			source: `
			  foo := true
				[3 if foo]
			`,
			wantStackTop: &value.List{
				value.SmallInt(3),
			},
		},
		"with string index": {
			source: `
			  foo := "3"
				[foo => :bar]
			`,
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::String` cannot be coerced into `Std::Int`",
			),
		},
		"with indices and if modifiers": {
			source: `
			  foo := "3"
				[3 => :bar if foo]
			`,
			wantStackTop: &value.List{
				value.Nil,
				value.Nil,
				value.Nil,
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
