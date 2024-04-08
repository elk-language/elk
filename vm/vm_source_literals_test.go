package vm_test

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestVMSource_RangeLiteral(t *testing.T) {
	tests := sourceTestTable{
		"closed range": {
			source:       `2...5`,
			wantStackTop: value.NewClosedRange(value.SmallInt(2), value.SmallInt(5)),
		},
		"dynamic closed range": {
			source: `
				a := 2.5
				a...23.8
			`,
			wantStackTop: value.NewClosedRange(value.Float(2.5), value.Float(23.8)),
		},

		"open range": {
			source:       `2<.<5`,
			wantStackTop: value.NewOpenRange(value.SmallInt(2), value.SmallInt(5)),
		},
		"dynamic open range": {
			source: `
				a := 2.5
				a<.<23.8
			`,
			wantStackTop: value.NewOpenRange(value.Float(2.5), value.Float(23.8)),
		},

		"left open range": {
			source:       `2<..5`,
			wantStackTop: value.NewLeftOpenRange(value.SmallInt(2), value.SmallInt(5)),
		},
		"dynamic left open range": {
			source: `
				a := 2.5
				a<..23.8
			`,
			wantStackTop: value.NewLeftOpenRange(value.Float(2.5), value.Float(23.8)),
		},

		"right open range": {
			source:       `2..<5`,
			wantStackTop: value.NewRightOpenRange(value.SmallInt(2), value.SmallInt(5)),
		},
		"dynamic right open range": {
			source: `
				a := 2.5
				a..<23.8
			`,
			wantStackTop: value.NewRightOpenRange(value.Float(2.5), value.Float(23.8)),
		},

		"beginless closed range": {
			source:       `...5`,
			wantStackTop: value.NewBeginlessClosedRange(value.SmallInt(5)),
		},
		"dynamic beginless closed range": {
			source: `
				a := 2.5
				...a
			`,
			wantStackTop: value.NewBeginlessClosedRange(value.Float(2.5)),
		},

		"beginless open range": {
			source:       `..<5`,
			wantStackTop: value.NewBeginlessOpenRange(value.SmallInt(5)),
		},
		"dynamic beginless open range": {
			source: `
				a := 2.5
				..<a
			`,
			wantStackTop: value.NewBeginlessOpenRange(value.Float(2.5)),
		},

		"endless closed range": {
			source:       `5...`,
			wantStackTop: value.NewEndlessClosedRange(value.SmallInt(5)),
		},
		"dynamic endless closed range": {
			source: `
				a := 2.5
				a...
			`,
			wantStackTop: value.NewEndlessClosedRange(value.Float(2.5)),
		},

		"endless open range": {
			source:       `5<..`,
			wantStackTop: value.NewEndlessOpenRange(value.SmallInt(5)),
		},
		"dynamic endless open range": {
			source: `
				a := 2.5
				a<..
			`,
			wantStackTop: value.NewEndlessOpenRange(value.Float(2.5)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

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
		"with static elements and for in loops": {
			source: `
			  arr := [5, 6, 7]
				%[1, i * 2 for i in arr, %[:foo]]
			`,
			wantStackTop: &value.ArrayTuple{
				value.SmallInt(1),
				value.SmallInt(10),
				value.SmallInt(12),
				value.SmallInt(14),
				&value.ArrayTuple{
					value.ToSymbol("foo"),
				},
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
		"word arrayTuple literal with static capacity": {
			source: `
				print(\w[foo bar baz]:20.inspect)
			`,
			wantStdout:   `["foo", "bar", "baz"]:20`,
			wantStackTop: value.Nil,
		},
		"symbol arrayTuple literal with static capacity": {
			source: `
				print(\s[foo bar baz]:20.inspect)
			`,
			wantStdout:   `[:foo, :bar, :baz]:20`,
			wantStackTop: value.Nil,
		},
		"bin arrayTuple literal with static capacity": {
			source: `
				print(\b[101 10 11]:20.inspect)
			`,
			wantStdout:   `[5, 2, 3]:20`,
			wantStackTop: value.Nil,
		},
		"hex arrayTuple literal with static capacity": {
			source: `
				print(\x[ff de4 5]:20.inspect)
			`,
			wantStdout:   `[255, 3556, 5]:20`,
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
		"with static elements and for in loops": {
			source: `
			  arr := [5, 6, 7]
				[1, i * 2 for i in arr, %[:foo]]
			`,
			wantStackTop: &value.ArrayList{
				value.SmallInt(1),
				value.SmallInt(10),
				value.SmallInt(12),
				value.SmallInt(14),
				&value.ArrayTuple{
					value.ToSymbol("foo"),
				},
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

func TestVMSource_HashSetLiteral(t *testing.T) {
	tests := sourceTestTable{
		"empty hashSet literal": {
			source:       `^[]`,
			wantStackTop: &value.HashSet{},
		},
		"static hashSet literal": {
			source: `^[1, 2.5, :foo]`,
			wantStackTop: vm.MustNewHashSetWithElements(
				nil,
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("foo"),
			),
		},
		"static hashSet literal with static capacity": {
			source: `
				^[1, 2.5, :foo]:20
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				23,
				1,
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("foo"),
			),
		},
		"word hashSet literal with static capacity": {
			source: `
				^w[foo bar baz]:20
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				23,
				1,
				value.String("foo"),
				value.String("bar"),
				value.String("baz"),
			),
		},
		"symbol hashSet literal with static capacity": {
			source: `
				^s[foo bar baz]:20
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				23,
				1,
				value.ToSymbol("foo"),
				value.ToSymbol("bar"),
				value.ToSymbol("baz"),
			),
		},
		"bin hashSet literal with static capacity": {
			source: `
				^b[101 10 11]:20
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				23,
				1,
				value.SmallInt(5),
				value.SmallInt(2),
				value.SmallInt(3),
			),
		},
		"hex arrayTuple literal with static capacity": {
			source: `
				^x[ff de4 5]:20
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				23,
				1,
				value.SmallInt(255),
				value.SmallInt(3556),
				value.SmallInt(5),
			),
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				^[1, 2.5, foo, :bar]
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				4,
				1,
				value.SmallInt(1),
				value.Float(2.5),
				value.String("foo var"),
				value.ToSymbol("bar"),
			),
		},
		"starts with dynamic elements": {
			source: `
				foo := "foo var"
				^[foo, 1, 2.5, :bar]
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				4,
				1,
				value.String("foo var"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			),
		},
		"starts with dynamic elements and has capacity": {
			source: `
			  cap := 5
				foo := "foo var"
				^[foo, 1, 2.5, :bar]:(cap + 2)
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElements(
				nil,
				11,
				value.String("foo var"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			),
		},
		"with falsy if": {
			source: `
				foo := nil
				^["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElements(
				nil,
				4,
				value.String("awesome"),
				value.Float(2.5),
				value.ToSymbol("bar"),
			),
		},
		"with truthy if": {
			source: `
				foo := 57
				^["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				6,
				1,
				value.String("awesome"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			),
		},
		"with falsy unless": {
			source: `
				foo := nil
				^["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				6,
				1,
				value.String("awesome"),
				value.SmallInt(1),
				value.Float(2.5),
				value.ToSymbol("bar"),
			),
		},
		"with truthy unless": {
			source: `
				foo := true
				^["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElements(
				nil,
				4,
				value.String("awesome"),
				value.Float(2.5),
				value.ToSymbol("bar"),
			),
		},
		"with static elements and for in loops": {
			source: `
			  arr := [5, 6, 7]
				^[1, i * 2 for i in arr, 2]
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElements(
				nil,
				6,
				value.SmallInt(1),
				value.SmallInt(10),
				value.SmallInt(12),
				value.SmallInt(14),
				value.SmallInt(2),
			),
		},
		"with initial modifier": {
			source: `
			  foo := true
				^[3 if foo]
			`,
			wantStackTop: vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				1,
				1,
				value.SmallInt(3),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_HashMapLiteral(t *testing.T) {
	tests := sourceTestTable{
		"empty": {
			source:       `{}`,
			wantStackTop: value.NewHashMap(0),
		},
		"static literal": {
			source: `{ 1 => 2.5, "bar" => :foo }`,
			wantStackTop: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.String("bar"),
					Value: value.ToSymbol("foo"),
				},
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.Float(2.5),
				},
			),
		},
		"static elements with for loops": {
			source: `{ 1 => 'foo', i => i ** 2 for i in [1, 2, 3, 4], 2 => 5.6 }`,
			wantStackTop: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.SmallInt(1),
				},
				value.Pair{
					Key:   value.SmallInt(2),
					Value: value.Float(5.6),
				},
				value.Pair{
					Key:   value.SmallInt(3),
					Value: value.SmallInt(9),
				},
				value.Pair{
					Key:   value.SmallInt(4),
					Value: value.SmallInt(16),
				},
			),
		},
		"static literal with mutable elements": {
			source: `{ 1 => 2.5, 0 => [1, 2], "bar" => :foo }`,
			wantStackTop: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.String("bar"),
					Value: value.ToSymbol("foo"),
				},
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.Float(2.5),
				},
				value.Pair{
					Key: value.SmallInt(0),
					Value: &value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
					},
				},
			),
		},
		"static arrayTuple literal with static capacity": {
			source: `
				{ 1 => 2.5, "bar" => :foo }:20
			`,
			wantStackTop: vm.MustNewHashMapWithCapacityAndElements(
				nil,
				22,
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.Float(2.5),
				},
				value.Pair{
					Key:   value.String("bar"),
					Value: value.ToSymbol("foo"),
				},
			),
		},
		"nested static": {
			source: `{ 1 => 2.5, foo: { "bar" => [] }, "baz" => [4] }`,
			wantStackTop: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.Float(2.5),
				},
				value.Pair{
					Key: value.ToSymbol("foo"),
					Value: vm.MustNewHashMapWithElements(
						nil,
						value.Pair{
							Key:   value.String("bar"),
							Value: &value.ArrayList{},
						},
					),
				},
				value.Pair{
					Key: value.String("baz"),
					Value: &value.ArrayList{
						value.SmallInt(4),
					},
				},
			),
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				{ 1 => 2.5, foo => :bar, "baz" => 9 }
			`,
			wantStackTop: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.Float(2.5),
				},
				value.Pair{
					Key:   value.String("foo var"),
					Value: value.ToSymbol("bar"),
				},
				value.Pair{
					Key:   value.String("baz"),
					Value: value.SmallInt(9),
				},
			),
		},
		"starts with dynamic elements": {
			source: `
				foo := "foo var"
				{ foo => 1, 2.5 => :bar }
			`,
			wantStackTop: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.String("foo var"),
					Value: value.SmallInt(1),
				},
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"starts with dynamic elements and has capacity": {
			source: `
			  cap := 5
				foo := "foo var"
				{ foo => 1, 2.5 => :bar }:(cap + 2)
			`,
			wantStackTop: vm.MustNewHashMapWithCapacityAndElements(
				nil,
				9,
				value.Pair{
					Key:   value.String("foo var"),
					Value: value.SmallInt(1),
				},
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"with falsy if": {
			source: `
				foo := nil
				{ "awesome" => 1 if foo, 2.5 => :bar }
			`,
			wantStackTop: vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"with truthy if": {
			source: `
				foo := 57
				{ "awesome" => 1 if foo, 2.5 => :bar }
			`,
			wantStackTop: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.String("awesome"),
					Value: value.SmallInt(1),
				},
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"with falsy unless": {
			source: `
				foo := nil
				{ "awesome" => 1 unless foo, 2.5 => :bar }
			`,
			wantStackTop: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.String("awesome"),
					Value: value.SmallInt(1),
				},
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"with truthy unless": {
			source: `
				foo := true
				{ "awesome" => 1 unless foo, 2.5 => :bar }
			`,
			wantStackTop: vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"with initial modifier": {
			source: `
			  foo := true
				{ 3 => 2 if foo }
			`,
			wantStackTop: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(3),
					Value: value.SmallInt(2),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_HashRecordLiteral(t *testing.T) {
	tests := sourceTestTable{
		"empty": {
			source:       `%{}`,
			wantStackTop: value.NewHashRecord(0),
		},
		"static literal": {
			source: `%{ 1 => 2.5, "bar" => :foo }`,
			wantStackTop: vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{
					Key:   value.String("bar"),
					Value: value.ToSymbol("foo"),
				},
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.Float(2.5),
				},
			),
		},
		"static elements with for loops": {
			source: `%{ 1 => 'foo', i => i ** 2 for i in [1, 2, 3, 4], 2 => 5.6 }`,
			wantStackTop: vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.SmallInt(1),
				},
				value.Pair{
					Key:   value.SmallInt(2),
					Value: value.Float(5.6),
				},
				value.Pair{
					Key:   value.SmallInt(3),
					Value: value.SmallInt(9),
				},
				value.Pair{
					Key:   value.SmallInt(4),
					Value: value.SmallInt(16),
				},
			),
		},
		"static literal with mutable elements": {
			source: `%{ 1 => 2.5, 0 => [1, 2], "bar" => :foo }`,
			wantStackTop: vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.Float(2.5),
				},
				value.Pair{
					Key:   value.String("bar"),
					Value: value.ToSymbol("foo"),
				},
				value.Pair{
					Key: value.SmallInt(0),
					Value: &value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
					},
				},
			),
		},
		"nested static": {
			source: `%{ 1 => 2.5, foo: %{ "bar" => [] }, "baz" => [4] }`,
			wantStackTop: vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key: value.String("baz"),
					Value: &value.ArrayList{
						value.SmallInt(4),
					},
				},
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.Float(2.5),
				},
				value.Pair{
					Key: value.ToSymbol("foo"),
					Value: vm.MustNewHashRecordWithCapacityAndElements(
						nil,
						5,
						value.Pair{
							Key:   value.String("bar"),
							Value: &value.ArrayList{},
						},
					),
				},
			),
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				%{ 1 => 2.5, foo => :bar, "baz" => 9 }
			`,
			wantStackTop: vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.SmallInt(1),
					Value: value.Float(2.5),
				},
				value.Pair{
					Key:   value.String("foo var"),
					Value: value.ToSymbol("bar"),
				},
				value.Pair{
					Key:   value.String("baz"),
					Value: value.SmallInt(9),
				},
			),
		},
		"starts with dynamic elements": {
			source: `
				foo := "foo var"
				%{ foo => 1, 2.5 => :bar }
			`,
			wantStackTop: vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{
					Key:   value.String("foo var"),
					Value: value.SmallInt(1),
				},
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"with falsy if": {
			source: `
				foo := nil
				%{ "awesome" => 1 if foo, 2.5 => :bar }
			`,
			wantStackTop: vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"with truthy if": {
			source: `
				foo := 57
				%{ "awesome" => 1 if foo, 2.5 => :bar }
			`,
			wantStackTop: vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.String("awesome"),
					Value: value.SmallInt(1),
				},
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"with falsy unless": {
			source: `
				foo := nil
				%{ "awesome" => 1 unless foo, 2.5 => :bar }
			`,
			wantStackTop: vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.String("awesome"),
					Value: value.SmallInt(1),
				},
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"with truthy unless": {
			source: `
				foo := true
				%{ "awesome" => 1 unless foo, 2.5 => :bar }
			`,
			wantStackTop: vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.Float(2.5),
					Value: value.ToSymbol("bar"),
				},
			),
		},
		"with initial modifier": {
			source: `
			  foo := true
				%{ 3 => 2 if foo }
			`,
			wantStackTop: vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.SmallInt(3),
					Value: value.SmallInt(2),
				},
			),
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

func TestVMSource_RegexLiteral(t *testing.T) {
	tests := sourceTestTable{
		"static regex": {
			source:       `%/foo \w+/im`,
			wantStackTop: value.MustCompileRegex(`foo \w+`, bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag|flag.MultilineFlag)),
		},
		"interpolated regex": {
			source: `
				bar := 15.2
				foo := %/foo/sa
				x := "x"

				%/foo: ${foo}, bar: ${bar + 2}, baz: ${nil}, ${x}/xi
			`,
			wantStackTop: value.MustCompileRegex(
				"foo: (?sa-imUx:foo), bar: 17.2, baz: , x",
				bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag|flag.ExtendedFlag),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
