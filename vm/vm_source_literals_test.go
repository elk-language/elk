package vm_test

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestVMSource_RangeLiteral(t *testing.T) {
	tests := sourceTestTable{
		"closed range": {
			source:       `2...5`,
			wantStackTop: value.Ref(value.NewClosedRange(value.SmallInt(2).ToValue(), value.SmallInt(5).ToValue())),
		},
		"dynamic closed range": {
			source: `
				a := 2.5
				a...23.8
			`,
			wantStackTop: value.Ref(value.NewClosedRange(value.Float(2.5).ToValue(), value.Float(23.8).ToValue())),
		},

		"open range": {
			source:       `2<.<5`,
			wantStackTop: value.Ref(value.NewOpenRange(value.SmallInt(2).ToValue(), value.SmallInt(5).ToValue())),
		},
		"dynamic open range": {
			source: `
				a := 2.5
				a<.<23.8
			`,
			wantStackTop: value.Ref(value.NewOpenRange(value.Float(2.5).ToValue(), value.Float(23.8).ToValue())),
		},

		"left open range": {
			source:       `2<..5`,
			wantStackTop: value.Ref(value.NewLeftOpenRange(value.SmallInt(2).ToValue(), value.SmallInt(5).ToValue())),
		},
		"dynamic left open range": {
			source: `
				a := 2.5
				a<..23.8
			`,
			wantStackTop: value.Ref(value.NewLeftOpenRange(value.Float(2.5).ToValue(), value.Float(23.8).ToValue())),
		},

		"right open range": {
			source:       `2..<5`,
			wantStackTop: value.Ref(value.NewRightOpenRange(value.SmallInt(2).ToValue(), value.SmallInt(5).ToValue())),
		},
		"dynamic right open range": {
			source: `
				a := 2.5
				a..<23.8
			`,
			wantStackTop: value.Ref(value.NewRightOpenRange(value.Float(2.5).ToValue(), value.Float(23.8).ToValue())),
		},

		"beginless closed range": {
			source:       `...5`,
			wantStackTop: value.Ref(value.NewBeginlessClosedRange(value.SmallInt(5).ToValue())),
		},
		"dynamic beginless closed range": {
			source: `
				a := 2.5
				...a
			`,
			wantStackTop: value.Ref(value.NewBeginlessClosedRange(value.Float(2.5).ToValue())),
		},

		"beginless open range": {
			source:       `..<5`,
			wantStackTop: value.Ref(value.NewBeginlessOpenRange(value.SmallInt(5).ToValue())),
		},
		"dynamic beginless open range": {
			source: `
				a := 2.5
				..<a
			`,
			wantStackTop: value.Ref(value.NewBeginlessOpenRange(value.Float(2.5).ToValue())),
		},

		"endless closed range": {
			source:       `5...`,
			wantStackTop: value.Ref(value.NewEndlessClosedRange(value.SmallInt(5).ToValue())),
		},
		"dynamic endless closed range": {
			source: `
				a := 2.5
				a...
			`,
			wantStackTop: value.Ref(value.NewEndlessClosedRange(value.Float(2.5).ToValue())),
		},

		"endless open range": {
			source:       `5<..`,
			wantStackTop: value.Ref(value.NewEndlessOpenRange(value.SmallInt(5).ToValue())),
		},
		"dynamic endless open range": {
			source: `
				a := 2.5
				a<..
			`,
			wantStackTop: value.Ref(value.NewEndlessOpenRange(value.Float(2.5).ToValue())),
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
			wantStackTop: value.Ref(&value.ArrayTuple{}),
		},
		"static arrayTuple literal": {
			source: `%[1, 2.5, :foo]`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("foo").ToValue(),
			}),
		},
		"nested static arrayTuple literal": {
			source: `%[1, 2.5, %["bar", %[]], %[:foo]]`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.Ref(&value.ArrayTuple{
					value.Ref(value.String("bar")),
					value.Ref(&value.ArrayTuple{}),
				}),
				value.Ref(&value.ArrayTuple{
					value.ToSymbol("foo").ToValue(),
				}),
			}),
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				%[1, 2.5, foo, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.Ref(value.String("foo var")),
				value.ToSymbol("bar").ToValue(),
			}),
		},
		"starts with dynamic elements": {
			source: `
				foo := "foo var"
				%[foo, 1, 2.5, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.Ref(value.String("foo var")),
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			}),
		},
		"with falsy if": {
			source: `
				foo := nil
				%["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.Ref(value.String("awesome")),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			}),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(38, 3, 23), P(40, 3, 25)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L(P(33, 3, 18), P(33, 3, 18)), "unreachable code"),
			},
		},
		"with truthy if": {
			source: `
				foo := 57
				%["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.Ref(value.String("awesome")),
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			}),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(37, 3, 23), P(39, 3, 25)), "this condition will always have the same result since type `Std::Int` is truthy"),
			},
		},
		"with falsy unless": {
			source: `
				foo := nil
				%["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.Ref(value.String("awesome")),
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			}),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(42, 3, 27), P(44, 3, 29)), "this condition will always have the same result since type `nil` is falsy"),
			},
		},
		"with truthy unless": {
			source: `
				foo := true
				%["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.Ref(value.String("awesome")),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			}),
		},
		"static with indices": {
			source: `%["awesome", 5 => :foo, 2 => 8.3]`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.Ref(value.String("awesome")),
				value.Nil,
				value.Float(8.3).ToValue(),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo").ToValue(),
			}),
		},
		"static indices with dynamic elements": {
			source: `
			  foo := 3
				%["awesome", 5 => :foo, 2 => 8.3, foo]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.Ref(value.String("awesome")),
				value.Nil,
				value.Float(8.3).ToValue(),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo").ToValue(),
				value.SmallInt(3).ToValue(),
			}),
		},
		"with dynamic elements and indices": {
			source: `
			  foo := 3
				%[foo, "awesome", 5 => :foo, 2 => 8.3]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.SmallInt(3).ToValue(),
				value.Ref(value.String("awesome")),
				value.Float(8.3).ToValue(),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo").ToValue(),
			}),
		},
		"with static elements and for in loops": {
			source: `
			  arr := [5, 6, 7]
				%[1, i * 2 for i in arr, %[:foo]]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.SmallInt(1).ToValue(),
				value.SmallInt(10).ToValue(),
				value.SmallInt(12).ToValue(),
				value.SmallInt(14).ToValue(),
				value.Ref(&value.ArrayTuple{
					value.ToSymbol("foo").ToValue(),
				}),
			}),
		},
		"with splats": {
			source: `
			  arr := [5, 6, 7]
				%[1, *arr, %[:foo]]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.SmallInt(1).ToValue(),
				value.SmallInt(5).ToValue(),
				value.SmallInt(6).ToValue(),
				value.SmallInt(7).ToValue(),
				value.Ref(&value.ArrayTuple{
					value.ToSymbol("foo").ToValue(),
				}),
			}),
		},
		"splat generator": {
			source: `
				def *gen: Int
					yield 5
					yield 6

					7
				end

				%[1, *gen(), %[:foo]]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.SmallInt(1).ToValue(),
				value.SmallInt(5).ToValue(),
				value.SmallInt(6).ToValue(),
				value.SmallInt(7).ToValue(),
				value.Ref(&value.ArrayTuple{
					value.ToSymbol("foo").ToValue(),
				}),
			}),
		},
		"splat custom iterable": {
			source: `
				class Foo
					include Iterator::Base[Int]

					var @counter: Int
					init
						@counter = 4
					end

					def next: Int ! :stop_iteration
						throw :stop_iteration if @counter >= 7

						@counter++
					end
				end

				f := Foo()
				%[1, *f, %[:foo]]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.SmallInt(1).ToValue(),
				value.SmallInt(5).ToValue(),
				value.SmallInt(6).ToValue(),
				value.SmallInt(7).ToValue(),
				value.Ref(&value.ArrayTuple{
					value.ToSymbol("foo").ToValue(),
				}),
			}),
		},
		"with dynamic indices": {
			source: `
			  foo := 3
				%[foo => :bar, "awesome"]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.Nil,
				value.Nil,
				value.Nil,
				value.ToSymbol("bar").ToValue(),
				value.Ref(value.String("awesome")),
			}),
		},
		"with initial modifier": {
			source: `
			  foo := true
				%[3 if foo]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.SmallInt(3).ToValue(),
			}),
		},
		"with string index": {
			source: `
			  foo := "3"
				%[foo => :bar]
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(23, 3, 7), P(25, 3, 9)), "index must be an integer, got type `Std::String`"),
			},
		},
		"with indices and if modifiers": {
			source: `
			  foo := "3"
				%[3 => :bar if foo]
			`,
			wantStackTop: value.Ref(&value.ArrayTuple{
				value.Nil,
				value.Nil,
				value.Nil,
				value.ToSymbol("bar").ToValue(),
			}),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(36, 3, 20), P(38, 3, 22)), "this condition will always have the same result since type `Std::String` is truthy"),
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
			wantStackTop: value.Ref(&value.ArrayList{}),
		},
		"static arrayTuple literal": {
			source: `[1, 2.5, :foo]`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("foo").ToValue(),
			}),
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
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.Ref(&value.ArrayList{
					value.Ref(value.String("bar")),
					value.Ref(&value.ArrayList{}),
				}),
				value.Ref(&value.ArrayList{
					value.ToSymbol("foo").ToValue(),
				}),
			}),
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				[1, 2.5, foo, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.Ref(value.String("foo var")),
				value.ToSymbol("bar").ToValue(),
			}),
		},
		"starts with dynamic elements": {
			source: `
				foo := "foo var"
				[foo, 1, 2.5, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.Ref(value.String("foo var")),
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			}),
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
			wantStackTop: value.Ref(&value.ArrayList{
				value.Ref(value.String("awesome")),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			}),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(37, 3, 22), P(39, 3, 24)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L(P(32, 3, 17), P(32, 3, 17)), "unreachable code"),
			},
		},
		"with truthy if": {
			source: `
				foo := 57
				["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.Ref(value.String("awesome")),
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			}),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(36, 3, 22), P(38, 3, 24)), "this condition will always have the same result since type `Std::Int` is truthy"),
			},
		},
		"with falsy unless": {
			source: `
				foo := nil
				["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.Ref(value.String("awesome")),
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			}),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(41, 3, 26), P(43, 3, 28)), "this condition will always have the same result since type `nil` is falsy"),
			},
		},
		"with truthy unless": {
			source: `
				foo := true
				["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.Ref(value.String("awesome")),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			}),
		},
		"static with indices": {
			source: `["awesome", 5 => :foo, 2 => 8.3]`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.Ref(value.String("awesome")),
				value.Nil,
				value.Float(8.3).ToValue(),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo").ToValue(),
			}),
		},
		"static indices with dynamic elements": {
			source: `
			  foo := 3
				["awesome", 5 => :foo, 2 => 8.3, foo]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.Ref(value.String("awesome")),
				value.Nil,
				value.Float(8.3).ToValue(),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo").ToValue(),
				value.SmallInt(3).ToValue(),
			}),
		},
		"with dynamic elements and indices": {
			source: `
			  foo := 3
				[foo, "awesome", 5 => :foo, 2 => 8.3]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(3).ToValue(),
				value.Ref(value.String("awesome")),
				value.Float(8.3).ToValue(),
				value.Nil,
				value.Nil,
				value.ToSymbol("foo").ToValue(),
			}),
		},
		"with static elements and for in loops": {
			source: `
			  arr := [5, 6, 7]
				[1, i * 2 for i in arr, %[:foo]]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(1).ToValue(),
				value.SmallInt(10).ToValue(),
				value.SmallInt(12).ToValue(),
				value.SmallInt(14).ToValue(),
				value.Ref(&value.ArrayTuple{
					value.ToSymbol("foo").ToValue(),
				}),
			}),
		},
		"with splats": {
			source: `
			  arr := [5, 6, 7]
				[1, *arr, %[:foo]]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(1).ToValue(),
				value.SmallInt(5).ToValue(),
				value.SmallInt(6).ToValue(),
				value.SmallInt(7).ToValue(),
				value.Ref(&value.ArrayTuple{
					value.ToSymbol("foo").ToValue(),
				}),
			}),
		},
		"splat generator": {
			source: `
				def *gen: Int
					yield 5
					yield 6

					7
				end

				[1, *gen(), %[:foo]]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(1).ToValue(),
				value.SmallInt(5).ToValue(),
				value.SmallInt(6).ToValue(),
				value.SmallInt(7).ToValue(),
				value.Ref(&value.ArrayTuple{
					value.ToSymbol("foo").ToValue(),
				}),
			}),
		},
		"splat custom iterable": {
			source: `
				class Foo
					include Iterator::Base[Int]

					var @counter: Int
					init
						@counter = 4
					end

					def next: Int ! :stop_iteration
						throw :stop_iteration if @counter >= 7

						@counter++
					end
				end

				f := Foo()
				[1, *f, %[:foo]]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(1).ToValue(),
				value.SmallInt(5).ToValue(),
				value.SmallInt(6).ToValue(),
				value.SmallInt(7).ToValue(),
				value.Ref(&value.ArrayTuple{
					value.ToSymbol("foo").ToValue(),
				}),
			}),
		},
		"with dynamic indices": {
			source: `
			  foo := 3
				[foo => :bar, "awesome"]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.Nil,
				value.Nil,
				value.Nil,
				value.ToSymbol("bar").ToValue(),
				value.Ref(value.String("awesome")),
			}),
		},
		"with initial modifier": {
			source: `
			  foo := true
				[3 if foo]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(3).ToValue(),
			}),
		},
		"with string index": {
			source: `
			  foo := "3"
				[foo => :bar]
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(22, 3, 6), P(24, 3, 8)), "index must be an integer, got type `Std::String`"),
			},
		},
		"with indices and if modifiers": {
			source: `
			  foo := "3"
				[3 => :bar if foo]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.Nil,
				value.Nil,
				value.Nil,
				value.ToSymbol("bar").ToValue(),
			}),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(35, 3, 19), P(37, 3, 21)), "this condition will always have the same result since type `Std::String` is truthy"),
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
			wantStackTop: value.Ref(&value.HashSet{}),
		},
		"static hashSet literal": {
			source: `^[1, 2.5, :foo]`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithElements(
				nil,
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("foo").ToValue(),
			)),
		},
		"static hashSet literal with static capacity": {
			source: `
				^[1, 2.5, :foo]:20
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				23,
				1,
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("foo").ToValue(),
			)),
		},
		"word hashSet literal with static capacity": {
			source: `
				^w[foo bar baz]:20
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				23,
				1,
				value.Ref(value.String("foo")),
				value.Ref(value.String("bar")),
				value.Ref(value.String("baz")),
			)),
		},
		"symbol hashSet literal with static capacity": {
			source: `
				^s[foo bar baz]:20
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				23,
				1,
				value.ToSymbol("foo").ToValue(),
				value.ToSymbol("bar").ToValue(),
				value.ToSymbol("baz").ToValue(),
			)),
		},
		"bin hashSet literal with static capacity": {
			source: `
				^b[101 10 11]:20
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				23,
				1,
				value.SmallInt(5).ToValue(),
				value.SmallInt(2).ToValue(),
				value.SmallInt(3).ToValue(),
			)),
		},
		"hex arrayTuple literal with static capacity": {
			source: `
				^x[ff de4 5]:20
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				23,
				1,
				value.SmallInt(255).ToValue(),
				value.SmallInt(3556).ToValue(),
				value.SmallInt(5).ToValue(),
			)),
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				^[1, 2.5, foo, :bar]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				4,
				1,
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.Ref(value.String("foo var")),
				value.ToSymbol("bar").ToValue(),
			)),
		},
		"starts with dynamic elements": {
			source: `
				foo := "foo var"
				^[foo, 1, 2.5, :bar]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				4,
				1,
				value.Ref(value.String("foo var")),
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			)),
		},
		"starts with dynamic elements and has capacity": {
			source: `
			  cap := 5
				foo := "foo var"
				^[foo, 1, 2.5, :bar]:(cap + 2)
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElements(
				nil,
				11,
				value.Ref(value.String("foo var")),
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			)),
		},
		"with falsy if": {
			source: `
				foo := nil
				^["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElements(
				nil,
				4,
				value.Ref(value.String("awesome")),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			)),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(38, 3, 23), P(40, 3, 25)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L(P(33, 3, 18), P(33, 3, 18)), "unreachable code"),
			},
		},
		"with truthy if": {
			source: `
				foo := 57
				^["awesome", 1 if foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				6,
				1,
				value.Ref(value.String("awesome")),
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			)),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(37, 3, 23), P(39, 3, 25)), "this condition will always have the same result since type `Std::Int` is truthy"),
			},
		},
		"with falsy unless": {
			source: `
				foo := nil
				^["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				6,
				1,
				value.Ref(value.String("awesome")),
				value.SmallInt(1).ToValue(),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			)),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(42, 3, 27), P(44, 3, 29)), "this condition will always have the same result since type `nil` is falsy"),
			},
		},
		"with truthy unless": {
			source: `
				foo := true
				^["awesome", 1 unless foo, 2.5, :bar]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElements(
				nil,
				4,
				value.Ref(value.String("awesome")),
				value.Float(2.5).ToValue(),
				value.ToSymbol("bar").ToValue(),
			)),
		},
		"with static elements and for in loops": {
			source: `
			  arr := [5, 6, 7]
				^[1, i * 2 for i in arr, 2]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElements(
				nil,
				6,
				value.SmallInt(1).ToValue(),
				value.SmallInt(10).ToValue(),
				value.SmallInt(12).ToValue(),
				value.SmallInt(14).ToValue(),
				value.SmallInt(2).ToValue(),
			)),
		},
		"with splats": {
			source: `
			  arr := [5, 6, 7]
				^[1, *arr, 2]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElements(
				nil,
				6,
				value.SmallInt(1).ToValue(),
				value.SmallInt(5).ToValue(),
				value.SmallInt(6).ToValue(),
				value.SmallInt(7).ToValue(),
				value.SmallInt(2).ToValue(),
			)),
		},
		"splat generator": {
			source: `
				def *gen: Int
					yield 5
					yield 6

					7
				end

				^[1, *gen(), 2]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElements(
				nil,
				6,
				value.SmallInt(1).ToValue(),
				value.SmallInt(5).ToValue(),
				value.SmallInt(6).ToValue(),
				value.SmallInt(7).ToValue(),
				value.SmallInt(2).ToValue(),
			)),
		},
		"splat custom iterable": {
			source: `
				class Foo
					include Iterator::Base[Int]

					var @counter: Int
					init
						@counter = 4
					end

					def next: Int ! :stop_iteration
						throw :stop_iteration if @counter >= 7

						@counter++
					end
				end

				f := Foo()
				^[1, *f, 2]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElements(
				nil,
				6,
				value.SmallInt(1).ToValue(),
				value.SmallInt(5).ToValue(),
				value.SmallInt(6).ToValue(),
				value.SmallInt(7).ToValue(),
				value.SmallInt(2).ToValue(),
			)),
		},
		"with initial modifier": {
			source: `
			  foo := true
				^[3 if foo]
			`,
			wantStackTop: value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
				nil,
				1,
				1,
				value.SmallInt(3).ToValue(),
			)),
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
			wantStackTop: value.Ref(value.NewHashMap(0)),
		},
		"static literal": {
			source: `{ 1 => 2.5, "bar" => :foo }`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.Ref(value.String("bar")),
					Value: value.ToSymbol("foo").ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.Float(2.5).ToValue(),
				},
			)),
		},
		"splat": {
			source: `
				a := { 1 => 1, 2 => 4, 3 => 9, 4 => 16 }
				{ 1 => 'foo', **a, 2 => 5.6 }
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(2).ToValue(),
					Value: value.Float(5.6).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(3).ToValue(),
					Value: value.SmallInt(9).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(4).ToValue(),
					Value: value.SmallInt(16).ToValue(),
				},
			)),
		},
		"splat generator": {
			source: `
				def *gen: Pair[Int, Int]
					yield Pair(1, 1)
					yield Pair(2, 4)
					yield Pair(3, 9)

					Pair(4, 16)
				end

				{ 1 => 'foo', **gen(), 2 => 5.6 }
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(2).ToValue(),
					Value: value.Float(5.6).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(3).ToValue(),
					Value: value.SmallInt(9).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(4).ToValue(),
					Value: value.SmallInt(16).ToValue(),
				},
			)),
		},
		"splat custom iterable": {
			source: `
				class Foo
					include Iterator::Base[Pair[Int, Int]]

					var @counter: Int
					init
						@counter = 0
					end

					def next: Pair[Int, Int] ! :stop_iteration
						throw :stop_iteration if @counter >= 4

						@counter++
						Pair(@counter, @counter ** 2)
					end
				end

				f := Foo()
				{ 1 => 'foo', **f, 2 => 5.6 }
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(2).ToValue(),
					Value: value.Float(5.6).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(3).ToValue(),
					Value: value.SmallInt(9).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(4).ToValue(),
					Value: value.SmallInt(16).ToValue(),
				},
			)),
		},
		"static elements with for loops": {
			source: `{ 1 => 'foo', i => i ** 2 for i in [1, 2, 3, 4], 2 => 5.6 }`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(2).ToValue(),
					Value: value.Float(5.6).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(3).ToValue(),
					Value: value.SmallInt(9).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(4).ToValue(),
					Value: value.SmallInt(16).ToValue(),
				},
			)),
		},
		"static literal with mutable elements": {
			source: `{ 1 => 2.5, 0 => [1, 2], "bar" => :foo }`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.Ref(value.String("bar")),
					Value: value.ToSymbol("foo").ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.Float(2.5).ToValue(),
				},
				value.Pair{
					Key: value.SmallInt(0).ToValue(),
					Value: value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
					}),
				},
			)),
		},
		"static arrayTuple literal with static capacity": {
			source: `
				{ 1 => 2.5, "bar" => :foo }:20
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithCapacityAndElements(
				nil,
				22,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.Float(2.5).ToValue(),
				},
				value.Pair{
					Key:   value.Ref(value.String("bar")),
					Value: value.ToSymbol("foo").ToValue(),
				},
			)),
		},
		"nested static": {
			source: `{ 1 => 2.5, foo: { "bar" => [] }, "baz" => [4] }`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.Float(2.5).ToValue(),
				},
				value.Pair{
					Key: value.ToSymbol("foo").ToValue(),
					Value: value.Ref(vm.MustNewHashMapWithElements(
						nil,
						value.Pair{
							Key:   value.Ref(value.String("bar")),
							Value: value.Ref(&value.ArrayList{}),
						},
					)),
				},
				value.Pair{
					Key: value.Ref(value.String("baz")),
					Value: value.Ref(&value.ArrayList{
						value.SmallInt(4).ToValue(),
					}),
				},
			)),
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				{ 1 => 2.5, foo => :bar, "baz" => 9 }
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.Float(2.5).ToValue(),
				},
				value.Pair{
					Key:   value.Ref(value.String("foo var")),
					Value: value.ToSymbol("bar").ToValue(),
				},
				value.Pair{
					Key:   value.Ref(value.String("baz")),
					Value: value.SmallInt(9).ToValue(),
				},
			)),
		},
		"starts with dynamic elements": {
			source: `
				foo := "foo var"
				{ foo => 1, 2.5 => :bar }
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.Ref(value.String("foo var")),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
		},
		"starts with dynamic elements and has capacity": {
			source: `
			  cap := 5
				foo := "foo var"
				{ foo => 1, 2.5 => :bar }:(cap + 2)
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithCapacityAndElements(
				nil,
				9,
				value.Pair{
					Key:   value.Ref(value.String("foo var")),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
		},
		"with falsy if": {
			source: `
				foo := nil
				{ "awesome" => 1 if foo, 2.5 => :bar }
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(40, 3, 25), P(42, 3, 27)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L(P(22, 3, 7), P(35, 3, 20)), "unreachable code"),
			},
		},
		"with truthy if": {
			source: `
				foo := 57
				{ "awesome" => 1 if foo, 2.5 => :bar }
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.Ref(value.String("awesome")),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(39, 3, 25), P(41, 3, 27)), "this condition will always have the same result since type `Std::Int` is truthy"),
			},
		},
		"with falsy unless": {
			source: `
				foo := nil
				{ "awesome" => 1 unless foo, 2.5 => :bar }
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.Ref(value.String("awesome")),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(44, 3, 29), P(46, 3, 31)), "this condition will always have the same result since type `nil` is falsy"),
			},
		},
		"with truthy unless": {
			source: `
				foo := true
				{ "awesome" => 1 unless foo, 2.5 => :bar }
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
		},
		"with initial modifier": {
			source: `
			  foo := true
				{ 3 => 2 if foo }
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(3).ToValue(),
					Value: value.SmallInt(2).ToValue(),
				},
			)),
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
			wantStackTop: value.Ref(value.NewHashRecord(0)),
		},
		"static literal": {
			source: `%{ 1 => 2.5, "bar" => :foo }`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{
					Key:   value.Ref(value.String("bar")),
					Value: value.ToSymbol("foo").ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.Float(2.5).ToValue(),
				},
			)),
		},
		"static elements with for loops": {
			source: `%{ 1 => 'foo', i => i ** 2 for i in [1, 2, 3, 4], 2 => 5.6 }`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(2).ToValue(),
					Value: value.Float(5.6).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(3).ToValue(),
					Value: value.SmallInt(9).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(4).ToValue(),
					Value: value.SmallInt(16).ToValue(),
				},
			)),
		},
		"splat": {
			source: `
				a := { 1 => 1, 2 => 4, 3 => 9, 4 => 16 }
				%{ 1 => 'foo', **a, 2 => 5.6 }
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(2).ToValue(),
					Value: value.Float(5.6).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(3).ToValue(),
					Value: value.SmallInt(9).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(4).ToValue(),
					Value: value.SmallInt(16).ToValue(),
				},
			)),
		},
		"splat generator": {
			source: `
				def *gen: Pair[Int, Int]
					yield Pair(1, 1)
					yield Pair(2, 4)
					yield Pair(3, 9)

					Pair(4, 16)
				end

				%{ 1 => 'foo', **gen(), 2 => 5.6 }
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(2).ToValue(),
					Value: value.Float(5.6).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(3).ToValue(),
					Value: value.SmallInt(9).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(4).ToValue(),
					Value: value.SmallInt(16).ToValue(),
				},
			)),
		},
		"splat custom iterable": {
			source: `
				class Foo
					include Iterator::Base[Pair[Int, Int]]

					var @counter: Int
					init
						@counter = 0
					end

					def next: Pair[Int, Int] ! :stop_iteration
						throw :stop_iteration if @counter >= 4

						@counter++
						Pair(@counter, @counter ** 2)
					end
				end

				f := Foo()
				%{ 1 => 'foo', **f, 2 => 5.6 }
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(2).ToValue(),
					Value: value.Float(5.6).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(3).ToValue(),
					Value: value.SmallInt(9).ToValue(),
				},
				value.Pair{
					Key:   value.SmallInt(4).ToValue(),
					Value: value.SmallInt(16).ToValue(),
				},
			)),
		},
		"static literal with mutable elements": {
			source: `%{ 1 => 2.5, 0 => [1, 2], "bar" => :foo }`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.Float(2.5).ToValue(),
				},
				value.Pair{
					Key:   value.Ref(value.String("bar")),
					Value: value.ToSymbol("foo").ToValue(),
				},
				value.Pair{
					Key: value.SmallInt(0).ToValue(),
					Value: value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
					}),
				},
			)),
		},
		"nested static": {
			source: `%{ 1 => 2.5, foo: %{ "bar" => [] }, "baz" => [4] }`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key: value.Ref(value.String("baz")),
					Value: value.Ref(&value.ArrayList{
						value.SmallInt(4).ToValue(),
					}),
				},
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.Float(2.5).ToValue(),
				},
				value.Pair{
					Key: value.ToSymbol("foo").ToValue(),
					Value: value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
						nil,
						5,
						value.Pair{
							Key:   value.Ref(value.String("bar")),
							Value: value.Ref(&value.ArrayList{}),
						},
					)),
				},
			)),
		},
		"starts with static elements": {
			source: `
				foo := "foo var"
				%{ 1 => 2.5, foo => :bar, "baz" => 9 }
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.SmallInt(1).ToValue(),
					Value: value.Float(2.5).ToValue(),
				},
				value.Pair{
					Key:   value.Ref(value.String("foo var")),
					Value: value.ToSymbol("bar").ToValue(),
				},
				value.Pair{
					Key:   value.Ref(value.String("baz")),
					Value: value.SmallInt(9).ToValue(),
				},
			)),
		},
		"starts with dynamic elements": {
			source: `
				foo := "foo var"
				%{ foo => 1, 2.5 => :bar }
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{
					Key:   value.Ref(value.String("foo var")),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
		},
		"with falsy if": {
			source: `
				foo := nil
				%{ "awesome" => 1 if foo, 2.5 => :bar }
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(41, 3, 26), P(43, 3, 28)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L(P(23, 3, 8), P(36, 3, 21)), "unreachable code"),
			},
		},
		"with truthy if": {
			source: `
				foo := 57
				%{ "awesome" => 1 if foo, 2.5 => :bar }
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.Ref(value.String("awesome")),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(40, 3, 26), P(42, 3, 28)), "this condition will always have the same result since type `Std::Int` is truthy"),
			},
		},
		"with falsy unless": {
			source: `
				foo := nil
				%{ "awesome" => 1 unless foo, 2.5 => :bar }
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.Ref(value.String("awesome")),
					Value: value.SmallInt(1).ToValue(),
				},
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(45, 3, 30), P(47, 3, 32)), "this condition will always have the same result since type `nil` is falsy"),
			},
		},
		"with truthy unless": {
			source: `
				foo := true
				%{ "awesome" => 1 unless foo, 2.5 => :bar }
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.Float(2.5).ToValue(),
					Value: value.ToSymbol("bar").ToValue(),
				},
			)),
		},
		"with initial modifier": {
			source: `
			  foo := true
				%{ 3 => 2 if foo }
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.SmallInt(3).ToValue(),
					Value: value.SmallInt(2).ToValue(),
				},
			)),
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
			wantStackTop: value.Ref(value.String("foo")),
		},
		"interpolated string": {
			source: `
				bar := 15.2
				foo := 1
				x := "x"

				"foo: ${foo + 2}, bar: ${bar}, baz: ${nil}, ${x}"
			`,
			wantStackTop: value.Ref(value.String("foo: 3, bar: 15.2, baz: , x")),
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
			wantStackTop: value.Ref(value.MustCompileRegex(`foo \w+`, bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag|flag.MultilineFlag))),
		},
		"interpolated regex": {
			source: `
				bar := 15.2
				foo := %/foo/sa
				x := "x"

				%/foo: ${foo}, bar: ${bar + 2}, baz: ${nil}, ${x}/xi
			`,
			wantStackTop: value.Ref(value.MustCompileRegex(
				"foo: (?sa-imUx:foo), bar: 17.2, baz: , x",
				bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag|flag.ExtendedFlag),
			)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
