package vm_test

import (
	"testing"

	"github.com/elk-language/elk/value"
)

func TestVMSource_ArrayTupleLiteral(t *testing.T) {
	tests := sourceTestTable{
		"empty arrayTuple literal": {
			source:       `%[]`,
			wantStackTop: &value.ArrayTuple{},
		},
		"static arrayTuple literal": {
			source: `%[1, 2.5, :foo]`,
			wantStackTop: &value.ArrayTuple{
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("foo"),
			},
		},
		"nested static arrayTuple literal": {
			source: `%[1, 2.5, %["bar", %[]], %[:foo]]`,
			wantStackTop: &value.ArrayTuple{
				value.SmallInt(1),
				value.Float(2.5),
				&value.ArrayTuple{
					value.String("bar"),
					&value.ArrayTuple{},
				},
				&value.ArrayTuple{
					value.ToSymbol("foo"),
				},
			},
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				%[1, 2.5, foo, :bar]
			`,
			wantStackTop: &value.ArrayTuple{
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
			wantStackTop: &value.ArrayTuple{
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
			wantStackTop: &value.ArrayTuple{
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
			wantStackTop: &value.ArrayTuple{
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
			wantStackTop: &value.ArrayTuple{
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
			wantStackTop: &value.ArrayTuple{
				value.String("awesome"),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"static with indices": {
			source: `%["awesome", 5 => :foo, 2 => 8.3]`,
			wantStackTop: &value.ArrayTuple{
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
			wantStackTop: &value.ArrayTuple{
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
			wantStackTop: &value.ArrayTuple{
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
			wantStackTop: &value.ArrayTuple{
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
			wantStackTop: &value.ArrayTuple{
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
			wantStackTop: &value.ArrayTuple{
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

func TestVMSource_ArrayListLiteral(t *testing.T) {
	tests := sourceTestTable{
		"empty arrayTuple literal": {
			source:       `[]`,
			wantStackTop: &value.ArrayList{},
		},
		"static arrayTuple literal": {
			source: `[1, 2.5, :foo]`,
			wantStackTop: &value.ArrayList{
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("foo"),
			},
		},
		"static arrayTuple literal with static capacity": {
			source: `
				print([1, 2.5, :foo]:20.inspect)
			`,
			wantStdout:   "[1, 2.5, :foo]:20",
			wantStackTop: value.Nil,
		},
		"nested static arrayTuple literal": {
			source: `[1, 2.5, ["bar", []], [:foo]]`,
			wantStackTop: &value.ArrayList{
				value.SmallInt(1),
				value.Float(2.5),
				&value.ArrayList{
					value.String("bar"),
					&value.ArrayList{},
				},
				&value.ArrayList{
					value.ToSymbol("foo"),
				},
			},
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				[1, 2.5, foo, :bar]
			`,
			wantStackTop: &value.ArrayList{
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
			wantStackTop: &value.ArrayList{
				value.String("foo var"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"starts with dynamic elements and has capacity": {
			source: `
			  cap := 5
				foo := "foo var"
				arr := [foo, 1, 2.5, :bar]:(cap + 2)
				println arr.inspect
			`,
			wantStdout:   "[\"foo var\", 1, 2.5, :bar]:7\n",
			wantStackTop: value.Nil,
		},
		"with falsy if": {
			source: `
				foo := nil
				["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: &value.ArrayList{
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
			wantStackTop: &value.ArrayList{
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
			wantStackTop: &value.ArrayList{
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
			wantStackTop: &value.ArrayList{
				value.String("awesome"),
				value.Float(2.5),
				value.ToSymbol("bar"),
			},
		},
		"static with indices": {
			source: `["awesome", 5 => :foo, 2 => 8.3]`,
			wantStackTop: &value.ArrayList{
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
			wantStackTop: &value.ArrayList{
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
			wantStackTop: &value.ArrayList{
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
			wantStackTop: &value.ArrayList{
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
			wantStackTop: &value.ArrayList{
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
			wantStackTop: &value.ArrayList{
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

func TestVMSource_StringLiteral(t *testing.T) {
	tests := sourceTestTable{
		"static string": {
			source:       `"foo"`,
			wantStackTop: value.String("foo"),
		},
		"interpolated string": {
			source: `
				bar := 15.2
				foo := 1
				x := "x"

				"foo: ${foo + 2}, bar: ${bar}, baz: ${nil}, ${x}"
			`,
			wantStackTop: value.String("foo: 3, bar: 15.2, baz: , x"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
