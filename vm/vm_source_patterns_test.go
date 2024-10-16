package vm_test

import (
	"testing"

	"github.com/elk-language/elk/value"
)

func TestVMSource_Switch(t *testing.T) {
	tests := sourceTestTable{
		"match no value": {
			source: `
				switch 20
		    case 0 then :a
				case 5 then :b
				case 10 then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.Nil,
		},
		"match with variable": {
			source: `
				switch 20
		    case 0 then :a
				case n then n + 2
				case 10 then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.SmallInt(22),
		},
		"match no value with else": {
			source: `
				switch 20
		    case 0 then :a
				case 5 then :b
				case 10 then :c
				case 15 then :d
				else :e
				end
			`,
			wantStackTop: value.ToSymbol("e"),
		},
		"match boolean": {
			source: `
				switch true
		    case 0 then :a
				case false then :b
				case true then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match nil": {
			source: `
				switch nil
		    case 0 then :a
				case nil then :b
				case true then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("b"),
		},
		"match string": {
			source: `
				switch "some string"
		    case 0 then :a
				case nil then :b
				case "some other string" then :b
				case "some string" then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match interpolated string": {
			source: `
				switch "some 7 string"
		    case 0 then :a
				case nil then :b
				case "some other string" then :b
				case "some ${5 + 2} string" then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match string with regex": {
			source: `
				switch "some string"
		    case 0 then :a
				case nil then :b
				case %/^some other string$/ then :b
				case %/some str\w+/ then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match string with interpolated regex": {
			source: `
				switch "some 7 string"
		    case 0 then :a
				case nil then :b
				case %/^some other string$/ then :b
				case %/some ${2 + 5} str\w+/ then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match symbol": {
			source: `
				switch :foo
		    case 0 then :a
				case :bar then :b
				case "foo" then :b
				case :foo then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match interpolated symbol": {
			source: `
				switch :foo7
		    case 0 then :a
				case :bar then :b
				case "foo7" then :b
				case :"foo${2 + 5}" then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match with method call": {
			source: `
				switch "10"
		    case > 20 then :a
				case > 5 && < 8 then :b
				case > 9 && < 15 then :c
				case == 10.to_string then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match comparison with and": {
			source: `
				switch 10
		    case > 20 then :a
				case > 5 && < 8 then :b
				case > 9 && < 15 then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match or": {
			source: `
				switch "foo"
		    case "bar" || "foo" then :a
				case > 5 && < 8 then :b
				case > 9 && < 15 then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("a"),
		},
		"match non-list with list patterns": {
			source: `
				switch :foo
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match tuple with list patterns": {
			source: `
				switch %[1, 6, 9, 20]
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},

		"match set": {
			source: `
				switch ^[1, 6, 9, 20]
		    case < 9 then :a
				case ^[1, 6, 10] then :b
				case ^[1, 6, 9, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match word set": {
			source: `
				switch ^['foo', 'bar']
		    case < 9 then :a
				case ^['foo', 'ba'] then :b
				case ^w[foo bar] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match symbol set": {
			source: `
				switch ^[:foo, :bar]
		    case < 9 then :a
				case ^[:foo, :ba] then :b
				case ^s[foo bar] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match hex set": {
			source: `
				switch ^[0xfe, 0x4]
		    case < 9 then :a
				case ^[0xfe, 0x5] then :b
				case ^x[fe 4] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match bin set": {
			source: `
				switch ^[0b11, 0b10]
		    case < 9 then :a
				case ^[0b11, 0b01] then :b
				case ^b[11 10] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match empty set": {
			source: `
				switch ^[]
		    case < 9 then :a
				case ^[1, 6, 10] then :b
				case ^[] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match set with rest elements": {
			source: `
				switch ^[1, 6, 9, 20]
		    case < 9 then :a
				case ^[1, 6, 10] then :b
				case ^[6, 20, 1, *] then :c
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match set with skip elements": {
			source: `
				switch ^[1, 6, 9, 20]
		    case < 9 then :a
				case ^[1, 6, _] then :b
				case ^[6, 20, 1, _] then :c
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},

		"match list": {
			source: `
				switch [1, 6, 9, 20]
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match word list": {
			source: `
				switch ['foo', 'bar']
		    case < 9 then :a
				case ['foo', 'ba'] then :b
				case \w[foo bar] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match tuple with word list": {
			source: `
				switch %['foo', 'bar']
		    case < 9 then :a
				case ['foo', 'ba'] then :b
				case \w[foo bar] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match symbol list": {
			source: `
				switch [:foo, :bar]
		    case < 9 then :a
				case [:foo, :ba] then :b
				case \s[foo bar] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match tuple with symbol list": {
			source: `
				switch %[:foo, :bar]
		    case < 9 then :a
				case [:foo, :ba] then :b
				case \s[foo bar] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match hex list": {
			source: `
				switch [0xfe, 0x4]
		    case < 9 then :a
				case [0xfe, 0x5] then :b
				case \x[fe 4] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match tuple with hex list": {
			source: `
				switch %[0xfe, 0x4]
		    case < 9 then :a
				case [0xfe, 0x5] then :b
				case \x[fe 4] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match bin list": {
			source: `
				switch [0b11, 0b10]
		    case < 9 then :a
				case [0b11, 0b01] then :b
				case \b[11 10] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match tuple with bin list": {
			source: `
				switch %[0b11, 0b10]
		    case < 9 then :a
				case [0b11, 0b01] then :b
				case \b[11 10] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match empty list": {
			source: `
				switch []
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match list with rest elements": {
			source: `
				switch [1, 6, 9, 20]
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, > 5, 2] then :c
				case [1, *a, > 15] then a
				case 15 then :e
				end
			`,
			wantStackTop: &value.ArrayList{
				value.SmallInt(6),
				value.SmallInt(9),
			},
		},
		"match list with unnamed rest elements": {
			source: `
				switch [1, 6, 9, 20]
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, > 5, 2] then :c
				case [1, *, < 15] then :d
				case [1, *, > 15] then :e
				case 15 then :f
				end
			`,
			wantStackTop: value.ToSymbol("e"),
		},
		"match nested lists": {
			source: `
				switch [1, 6, [17, 43, [71, 28]], 20]
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, [17, 43, [42, 28]], 20] then :c
				case [1, 6, [17, > 40, [71, 28]], > 15] then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},

		"match non-list with tuple patterns": {
			source: `
				switch :foo
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match list with tuple patterns": {
			source: `
				switch [1, 6, 9, 20]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match tuple": {
			source: `
				switch %[1, 6, 9, 20]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match word tuple": {
			source: `
				switch %['foo', 'bar']
		    case < 9 then :a
				case %['foo', 'ba'] then :b
				case %w[foo bar] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match list with word tuple": {
			source: `
				switch ['foo', 'bar']
		    case < 9 then :a
				case %['foo', 'ba'] then :b
				case %w[foo bar] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match symbol tuple": {
			source: `
				switch %[:foo, :bar]
		    case < 9 then :a
				case %[:foo, :ba] then :b
				case %s[foo bar] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match list with symbol tuple": {
			source: `
				switch [:foo, :bar]
		    case < 9 then :a
				case %[:foo, :ba] then :b
				case %s[foo bar] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match hex tuple": {
			source: `
				switch %[0xfe, 0x4]
		    case < 9 then :a
				case %[0xfe, 0x5] then :b
				case %x[fe 4] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match list with hex tuple": {
			source: `
				switch [0xfe, 0x4]
		    case < 9 then :a
				case %[0xfe, 0x5] then :b
				case %x[fe 4] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match bin tuple": {
			source: `
				switch %[0b11, 0b10]
		    case < 9 then :a
				case %[0b11, 0b01] then :b
				case %b[11 10] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match list with bin tuple": {
			source: `
				switch [0b11, 0b10]
		    case < 9 then :a
				case %[0b11, 0b01] then :b
				case %b[11 10] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match empty tuple": {
			source: `
				switch %[]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match tuple with rest elements": {
			source: `
				switch %[1, 6, 9, 20]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, > 5, 2] then :c
				case %[1, *a, > 15] then a
				case 15 then :e
				end
			`,
			wantStackTop: &value.ArrayList{
				value.SmallInt(6),
				value.SmallInt(9),
			},
		},
		"match tuple with unnamed rest elements": {
			source: `
				switch %[1, 6, 9, 20]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, > 5, 2] then :c
				case %[1, *, < 15] then :d
				case %[1, *, > 15] then :e
				case 15 then :f
				end
			`,
			wantStackTop: value.ToSymbol("e"),
		},
		"match nested tuples": {
			source: `
				switch %[1, 6, %[17, 43, %[71, 28]], 20]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, %[17, 43, %[42, 28]], 20] then :c
				case %[1, 6, %[17, > 40, %[71, 28]], > 15] then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},

		"match empty map": {
			source: `
				switch {}
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case { 1 => > 2, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case {} then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match map": {
			source: `
				switch { 1 => 5.5, foo: "bar", "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case { 1 => > 2.0, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case { 1 => > 2.0, foo: "baz" || "bar", "baz" => < 13.0 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match map with record": {
			source: `
				switch { 1 => 5.5, foo: "bar", "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %{ 1 => > 2.0, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case %{ 1 => > 2.0, foo: "baz" || "bar", "baz" => < 13.0 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match nested maps": {
			source: `
				switch { 1 => 5.5, foo: ["bar", 5, 4, { elo: "mordo" }], "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case { 1 => > 2.0, foo: ["baz" || "bar", 5, 4, { eli: "mordo" }], "baz" => < 13.0 } then :c
				case { 1 => > 2.0, foo: ["baz" || "bar", 5, 4, { elo: %/^mord\w+$/ }], "baz" => < 13.0 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},

		"match empty record": {
			source: `
				switch %{}
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %{ 1 => > 2, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case %{} then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match record": {
			source: `
				switch %{ 1 => 5.5, foo: "bar", "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %{ 1 => > 2.0, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case %{ 1 => > 2.0, foo: "baz" || "bar", "baz" => < 13.0 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match record with map": {
			source: `
				switch %{ 1 => 5.5, foo: "bar", "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case { 1 => > 2, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case { 1 => > 2, foo: "baz" || "bar", "baz" => < 13 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match nested records": {
			source: `
				switch %{ 1 => 5.5, foo: ["bar", 5, 4, %{ elo: "mordo" }], "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %{ 1 => > 2.0, foo: ["baz" || "bar", 5, 4, %{ eli: "mordo" }], "baz" => < 13.0 } then :c
				case %{ 1 => > 2.0, foo: ["baz" || "bar", 5, 4, %{ elo: %/^mord\w+$/ }], "baz" => < 13.0 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},

		"match object": {
			source: `
				switch "foo bar"
		    case ::Std::Int() then :a
				case ::Std::String(length: < 4) then :b
				case ::Std::String(uppercase, length: > 3 as l) then [uppercase, l]
				case {} then :d
				case 15 then :e
				end
			`,
			wantStackTop: &value.ArrayList{
				value.String("FOO BAR"),
				value.SmallInt(7),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
