package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/value"
)

func TestVMSource_GreaterThan(t *testing.T) {
	tests := sourceTestTable{
		// String
		"'25' > '25'": {
			source:       "'25' > '25'",
			wantStackTop: value.False.ToValue(),
		},
		"'7' > '10'": {
			source:       "'7' > '10'",
			wantStackTop: value.True.ToValue(),
		},
		"'10' > '7'": {
			source:       "'10' > '7'",
			wantStackTop: value.False.ToValue(),
		},
		"'25' > '22'": {
			source:       "'25' > '22'",
			wantStackTop: value.True.ToValue(),
		},
		"'22' > '25'": {
			source:       "'22' > '25'",
			wantStackTop: value.False.ToValue(),
		},
		"'foo' > 'foo'": {
			source:       "'foo' > 'foo'",
			wantStackTop: value.False.ToValue(),
		},
		"'foo' > 'foa'": {
			source:       "'foo' > 'foa'",
			wantStackTop: value.True.ToValue(),
		},
		"'foa' > 'foo'": {
			source:       "'foa' > 'foo'",
			wantStackTop: value.False.ToValue(),
		},
		"'foo' > 'foo bar'": {
			source:       "'foo' > 'foo bar'",
			wantStackTop: value.False.ToValue(),
		},
		"'foo bar' > 'foo'": {
			source:       "'foo bar' > 'foo'",
			wantStackTop: value.True.ToValue(),
		},

		"'2' > `2`": {
			source:       "'2' > `2`",
			wantStackTop: value.False.ToValue(),
		},
		"'72' > `7`": {
			source:       "'72' > `7`",
			wantStackTop: value.True.ToValue(),
		},
		"'8' > `7`": {
			source:       "'8' > `7`",
			wantStackTop: value.True.ToValue(),
		},
		"'7' > `8`": {
			source:       "'7' > `8`",
			wantStackTop: value.False.ToValue(),
		},
		"'ba' > `b`": {
			source:       "'ba' > `b`",
			wantStackTop: value.True.ToValue(),
		},
		"'b' > `a`": {
			source:       "'b' > `a`",
			wantStackTop: value.True.ToValue(),
		},
		"'a' > `b`": {
			source:       "'a' > `b`",
			wantStackTop: value.False.ToValue(),
		},

		"'2' > 2.0": {
			source: "'2' > 2.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(8, 1, 9)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `2.0`"),
			},
		},

		"'28' > 25.2bf": {
			source: "'28' > 25.2bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `25.2bf`"),
			},
		},

		"'28.8' > 12.9f64": {
			source: "'28.8' > 12.9f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(9, 1, 10), P(15, 1, 16)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `12.9f64`"),
			},
		},

		"'28.8' > 12.9f32": {
			source: "'28.8' > 12.9f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(9, 1, 10), P(15, 1, 16)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `12.9f32`"),
			},
		},

		"'93' > 19i64": {
			source: "'93' > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `19i64`"),
			},
		},

		"'93' > 19i32": {
			source: "'93' > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `19i32`"),
			},
		},

		"'93' > 19i16": {
			source: "'93' > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `19i16`"),
			},
		},

		"'93' > 19i8": {
			source: "'93' > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `19i8`"),
			},
		},

		"'93' > 19u64": {
			source: "'93' > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `19u64`"),
			},
		},

		"'93' > 19u32": {
			source: "'93' > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `19u32`"),
			},
		},

		"'93' > 19u16": {
			source: "'93' > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `19u16`"),
			},
		},

		"'93' > 19u8": {
			source: "'93' > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>`, got type `19u8`"),
			},
		},

		// Char
		"`2` > `2`": {
			source:       "`2` > `2`",
			wantStackTop: value.False.ToValue(),
		},
		"`8` > `7`": {
			source:       "`8` > `7`",
			wantStackTop: value.True.ToValue(),
		},
		"`7` > `8`": {
			source:       "`7` > `8`",
			wantStackTop: value.False.ToValue(),
		},
		"`b` > `a`": {
			source:       "`b` > `a`",
			wantStackTop: value.True.ToValue(),
		},
		"`a` > `b`": {
			source:       "`a` > `b`",
			wantStackTop: value.False.ToValue(),
		},

		"`2` > '2'": {
			source:       "`2` > '2'",
			wantStackTop: value.False.ToValue(),
		},
		"`7` > '72'": {
			source:       "`7` > '72'",
			wantStackTop: value.False.ToValue(),
		},
		"`8` > '7'": {
			source:       "`8` > '7'",
			wantStackTop: value.True.ToValue(),
		},
		"`7` > '8'": {
			source:       "`7` > '8'",
			wantStackTop: value.False.ToValue(),
		},
		"`b` > 'a'": {
			source:       "`b` > 'a'",
			wantStackTop: value.True.ToValue(),
		},
		"`b` > 'ba'": {
			source:       "`b` > 'ba'",
			wantStackTop: value.False.ToValue(),
		},
		"`a` > 'b'": {
			source:       "`a` > 'b'",
			wantStackTop: value.False.ToValue(),
		},

		"`2` > 2.0": {
			source: "`2` > 2.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(8, 1, 9)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `2.0`"),
			},
		},
		"`i` > 25.2bf": {
			source: "`i` > 25.2bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `25.2bf`"),
			},
		},
		"`f` > 12.9f64": {
			source: "`f` > 12.9f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `12.9f64`"),
			},
		},
		"`0` > 12.9f32": {
			source: "`0` > 12.9f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `12.9f32`"),
			},
		},
		"`9` > 19i64": {
			source: "`9` > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `19i64`"),
			},
		},
		"`u` > 19i32": {
			source: "`u` > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `19i32`"),
			},
		},
		"`4` > 19i16": {
			source: "`4` > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `19i16`"),
			},
		},
		"`6` > 19i8": {
			source: "`6` > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `19i8`"),
			},
		},
		"`9` > 19u64": {
			source: "`9` > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `19u64`"),
			},
		},
		"`u` > 19u32": {
			source: "`u` > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `19u32`"),
			},
		},
		"`4` > 19u16": {
			source: "`4` > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `19u16`"),
			},
		},
		"`6` > 19u8": {
			source: "`6` > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>`, got type `19u8`"),
			},
		},

		// Int
		"25 > 25": {
			source:       "25 > 25",
			wantStackTop: value.False.ToValue(),
		},
		"25 > -25": {
			source:       "25 > -25",
			wantStackTop: value.True.ToValue(),
		},
		"-25 > 25": {
			source:       "-25 > 25",
			wantStackTop: value.False.ToValue(),
		},
		"13 > 7": {
			source:       "13 > 7",
			wantStackTop: value.True.ToValue(),
		},
		"7 > 13": {
			source:       "7 > 13",
			wantStackTop: value.False.ToValue(),
		},

		"25 > 25.0": {
			source:       "25 > 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"25 > -25.0": {
			source:       "25 > -25.0",
			wantStackTop: value.True.ToValue(),
		},
		"-25 > 25.0": {
			source:       "-25 > 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"13 > 7.0": {
			source:       "13 > 7.0",
			wantStackTop: value.True.ToValue(),
		},
		"7 > 13.0": {
			source:       "7 > 13.0",
			wantStackTop: value.False.ToValue(),
		},
		"7 > 7.5": {
			source:       "7 > 7.5",
			wantStackTop: value.False.ToValue(),
		},
		"7 > 6.9": {
			source:       "7 > 6.9",
			wantStackTop: value.True.ToValue(),
		},

		"25 > 25bf": {
			source:       "25 > 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"25 > -25bf": {
			source:       "25 > -25bf",
			wantStackTop: value.True.ToValue(),
		},
		"-25 > 25bf": {
			source:       "-25 > 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"13 > 7bf": {
			source:       "13 > 7bf",
			wantStackTop: value.True.ToValue(),
		},
		"7 > 13bf": {
			source:       "7 > 13bf",
			wantStackTop: value.False.ToValue(),
		},
		"7 > 7.5bf": {
			source:       "7 > 7.5bf",
			wantStackTop: value.False.ToValue(),
		},
		"7 > 6.9bf": {
			source:       "7 > 6.9bf",
			wantStackTop: value.True.ToValue(),
		},

		"6 > 19f64": {
			source: "6 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>`, got type `19f64`"),
			},
		},
		"6 > 19f32": {
			source: "6 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>`, got type `19f32`"),
			},
		},
		"6 > 19i64": {
			source: "6 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>`, got type `19i64`"),
			},
		},
		"6 > 19i32": {
			source: "6 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>`, got type `19i32`"),
			},
		},
		"6 > 19i16": {
			source: "6 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>`, got type `19i16`"),
			},
		},
		"6 > 19i8": {
			source: "6 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(7, 1, 8)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>`, got type `19i8`"),
			},
		},
		"6 > 19u64": {
			source: "6 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>`, got type `19u64`"),
			},
		},
		"6 > 19u32": {
			source: "6 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>`, got type `19u32`"),
			},
		},
		"6 > 19u16": {
			source: "6 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>`, got type `19u16`"),
			},
		},
		"6 > 19u8": {
			source: "6 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(7, 1, 8)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>`, got type `19u8`"),
			},
		},

		// Float
		"25.0 > 25.0": {
			source:       "25.0 > 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"25.0 > -25.0": {
			source:       "25.0 > -25.0",
			wantStackTop: value.True.ToValue(),
		},
		"-25.0 > 25.0": {
			source:       "-25.0 > 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"13.0 > 7.0": {
			source:       "13.0 > 7.0",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 > 13.0": {
			source:       "7.0 > 13.0",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 > 7.5": {
			source:       "7.0 > 7.5",
			wantStackTop: value.False.ToValue(),
		},
		"7.5 > 7.0": {
			source:       "7.5 > 7.0",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 > 6.9": {
			source:       "7.0 > 6.9",
			wantStackTop: value.True.ToValue(),
		},

		"25.0 > 25": {
			source:       "25.0 > 25",
			wantStackTop: value.False.ToValue(),
		},
		"25.0 > -25": {
			source:       "25.0 > -25",
			wantStackTop: value.True.ToValue(),
		},
		"-25.0 > 25": {
			source:       "-25.0 > 25",
			wantStackTop: value.False.ToValue(),
		},
		"13.0 > 7": {
			source:       "13.0 > 7",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 > 13": {
			source:       "7.0 > 13",
			wantStackTop: value.False.ToValue(),
		},
		"7.5 > 7": {
			source:       "7.5 > 7",
			wantStackTop: value.True.ToValue(),
		},

		"25.0 > 25bf": {
			source:       "25.0 > 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"25.0 > -25bf": {
			source:       "25.0 > -25bf",
			wantStackTop: value.True.ToValue(),
		},
		"-25.0 > 25bf": {
			source:       "-25.0 > 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"13.0 > 7bf": {
			source:       "13.0 > 7bf",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 > 13bf": {
			source:       "7.0 > 13bf",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 > 7.5bf": {
			source:       "7.0 > 7.5bf",
			wantStackTop: value.False.ToValue(),
		},
		"7.5 > 7bf": {
			source:       "7.5 > 7bf",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 > 6.9bf": {
			source:       "7.0 > 6.9bf",
			wantStackTop: value.True.ToValue(),
		},

		"6.0 > 19f64": {
			source: "6.0 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>`, got type `19f64`"),
			},
		},
		"6.0 > 19f32": {
			source: "6.0 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>`, got type `19f32`"),
			},
		},
		"6.0 > 19i64": {
			source: "6.0 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>`, got type `19i64`"),
			},
		},
		"6.0 > 19i32": {
			source: "6.0 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>`, got type `19i32`"),
			},
		},
		"6.0 > 19i16": {
			source: "6.0 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>`, got type `19i16`"),
			},
		},
		"6.0 > 19i8": {
			source: "6.0 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>`, got type `19i8`"),
			},
		},
		"6.0 > 19u64": {
			source: "6.0 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>`, got type `19u64`"),
			},
		},
		"6.0 > 19u32": {
			source: "6.0 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>`, got type `19u32`"),
			},
		},
		"6.0 > 19u16": {
			source: "6.0 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>`, got type `19u16`"),
			},
		},
		"6.0 > 19u8": {
			source: "6.0 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>`, got type `19u8`"),
			},
		},

		// BigFloat
		"25bf > 25.0": {
			source:       "25bf > 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"25bf > -25.0": {
			source:       "25bf > -25.0",
			wantStackTop: value.True.ToValue(),
		},
		"-25bf > 25.0": {
			source:       "-25bf > 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"13bf > 7.0": {
			source:       "13bf > 7.0",
			wantStackTop: value.True.ToValue(),
		},
		"7bf > 13.0": {
			source:       "7bf > 13.0",
			wantStackTop: value.False.ToValue(),
		},
		"7bf > 7.5": {
			source:       "7bf > 7.5",
			wantStackTop: value.False.ToValue(),
		},
		"7.5bf > 7.0": {
			source:       "7.5bf > 7.0",
			wantStackTop: value.True.ToValue(),
		},
		"7bf > 6.9": {
			source:       "7bf > 6.9",
			wantStackTop: value.True.ToValue(),
		},

		"25bf > 25": {
			source:       "25bf > 25",
			wantStackTop: value.False.ToValue(),
		},
		"25bf > -25": {
			source:       "25bf > -25",
			wantStackTop: value.True.ToValue(),
		},
		"-25bf > 25": {
			source:       "-25bf > 25",
			wantStackTop: value.False.ToValue(),
		},
		"13bf > 7": {
			source:       "13bf > 7",
			wantStackTop: value.True.ToValue(),
		},
		"7bf > 13": {
			source:       "7bf > 13",
			wantStackTop: value.False.ToValue(),
		},
		"7.5bf > 7": {
			source:       "7.5bf > 7",
			wantStackTop: value.True.ToValue(),
		},

		"25bf > 25bf": {
			source:       "25bf > 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"25bf > -25bf": {
			source:       "25bf > -25bf",
			wantStackTop: value.True.ToValue(),
		},
		"-25bf > 25bf": {
			source:       "-25bf > 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"13bf > 7bf": {
			source:       "13bf > 7bf",
			wantStackTop: value.True.ToValue(),
		},
		"7bf > 13bf": {
			source:       "7bf > 13bf",
			wantStackTop: value.False.ToValue(),
		},
		"7bf > 7.5bf": {
			source:       "7bf > 7.5bf",
			wantStackTop: value.False.ToValue(),
		},
		"7.5bf > 7bf": {
			source:       "7.5bf > 7bf",
			wantStackTop: value.True.ToValue(),
		},
		"7bf > 6.9bf": {
			source:       "7bf > 6.9bf",
			wantStackTop: value.True.ToValue(),
		},

		"6bf > 19f64": {
			source: "6bf > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>`, got type `19f64`"),
			},
		},
		"6bf > 19f32": {
			source: "6bf > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>`, got type `19f32`"),
			},
		},
		"6bf > 19i64": {
			source: "6bf > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>`, got type `19i64`"),
			},
		},
		"6bf > 19i32": {
			source: "6bf > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>`, got type `19i32`"),
			},
		},
		"6bf > 19i16": {
			source: "6bf > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>`, got type `19i16`"),
			},
		},
		"6bf > 19i8": {
			source: "6bf > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>`, got type `19i8`"),
			},
		},
		"6bf > 19u64": {
			source: "6bf > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>`, got type `19u64`"),
			},
		},
		"6bf > 19u32": {
			source: "6bf > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>`, got type `19u32`"),
			},
		},
		"6bf > 19u16": {
			source: "6bf > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>`, got type `19u16`"),
			},
		},
		"6bf > 19u8": {
			source: "6bf > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>`, got type `19u8`"),
			},
		},

		// Float64
		"25f64 > 25f64": {
			source:       "25f64 > 25f64",
			wantStackTop: value.False.ToValue(),
		},
		"25f64 > -25f64": {
			source:       "25f64 > -25f64",
			wantStackTop: value.True.ToValue(),
		},
		"-25f64 > 25f64": {
			source:       "-25f64 > 25f64",
			wantStackTop: value.False.ToValue(),
		},
		"13f64 > 7f64": {
			source:       "13f64 > 7f64",
			wantStackTop: value.True.ToValue(),
		},
		"7f64 > 13f64": {
			source:       "7f64 > 13f64",
			wantStackTop: value.False.ToValue(),
		},
		"7f64 > 7.5f64": {
			source:       "7f64 > 7.5f64",
			wantStackTop: value.False.ToValue(),
		},
		"7.5f64 > 7f64": {
			source:       "7.5f64 > 7f64",
			wantStackTop: value.True.ToValue(),
		},
		"7f64 > 6.9f64": {
			source:       "7f64 > 6.9f64",
			wantStackTop: value.True.ToValue(),
		},

		"6f64 > 19.0": {
			source: "6f64 > 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19.0`"),
			},
		},

		"6f64 > 19": {
			source: "6f64 > 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19`"),
			},
		},
		"6f64 > 19bf": {
			source: "6f64 > 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19bf`"),
			},
		},
		"6f64 > 19f32": {
			source: "6f64 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19f32`"),
			},
		},
		"6f64 > 19i64": {
			source: "6f64 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19i64`"),
			},
		},
		"6f64 > 19i32": {
			source: "6f64 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19i32`"),
			},
		},
		"6f64 > 19i16": {
			source: "6f64 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19i16`"),
			},
		},
		"6f64 > 19i8": {
			source: "6f64 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19i8`"),
			},
		},
		"6f64 > 19u64": {
			source: "6f64 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19u64`"),
			},
		},
		"6f64 > 19u32": {
			source: "6f64 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19u32`"),
			},
		},
		"6f64 > 19u16": {
			source: "6f64 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19u16`"),
			},
		},
		"6f64 > 19u8": {
			source: "6f64 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>`, got type `19u8`"),
			},
		},

		// Float32
		"25f32 > 25f32": {
			source:       "25f32 > 25f32",
			wantStackTop: value.False.ToValue(),
		},
		"25f32 > -25f32": {
			source:       "25f32 > -25f32",
			wantStackTop: value.True.ToValue(),
		},
		"-25f32 > 25f32": {
			source:       "-25f32 > 25f32",
			wantStackTop: value.False.ToValue(),
		},
		"13f32 > 7f32": {
			source:       "13f32 > 7f32",
			wantStackTop: value.True.ToValue(),
		},
		"7f32 > 13f32": {
			source:       "7f32 > 13f32",
			wantStackTop: value.False.ToValue(),
		},
		"7f32 > 7.5f32": {
			source:       "7f32 > 7.5f32",
			wantStackTop: value.False.ToValue(),
		},
		"7.5f32 > 7f32": {
			source:       "7.5f32 > 7f32",
			wantStackTop: value.True.ToValue(),
		},
		"7f32 > 6.9f32": {
			source:       "7f32 > 6.9f32",
			wantStackTop: value.True.ToValue(),
		},

		"6f32 > 19.0": {
			source: "6f32 > 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19.0`"),
			},
		},

		"6f32 > 19": {
			source: "6f32 > 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19`"),
			},
		},
		"6f32 > 19bf": {
			source: "6f32 > 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19bf`"),
			},
		},
		"6f32 > 19f64": {
			source: "6f32 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19f64`"),
			},
		},
		"6f32 > 19i64": {
			source: "6f32 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19i64`"),
			},
		},
		"6f32 > 19i32": {
			source: "6f32 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19i32`"),
			},
		},
		"6f32 > 19i16": {
			source: "6f32 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19i16`"),
			},
		},
		"6f32 > 19i8": {
			source: "6f32 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19i8`"),
			},
		},
		"6f32 > 19u64": {
			source: "6f32 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19u64`"),
			},
		},
		"6f32 > 19u32": {
			source: "6f32 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19u32`"),
			},
		},
		"6f32 > 19u16": {
			source: "6f32 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19u16`"),
			},
		},
		"6f32 > 19u8": {
			source: "6f32 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>`, got type `19u8`"),
			},
		},

		// Int64
		"25i64 > 25i64": {
			source:       "25i64 > 25i64",
			wantStackTop: value.False.ToValue(),
		},
		"25i64 > -25i64": {
			source:       "25i64 > -25i64",
			wantStackTop: value.True.ToValue(),
		},
		"-25i64 > 25i64": {
			source:       "-25i64 > 25i64",
			wantStackTop: value.False.ToValue(),
		},
		"13i64 > 7i64": {
			source:       "13i64 > 7i64",
			wantStackTop: value.True.ToValue(),
		},
		"7i64 > 13i64": {
			source:       "7i64 > 13i64",
			wantStackTop: value.False.ToValue(),
		},

		"6i64 > 19": {
			source: "6i64 > 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19`"),
			},
		},
		"6i64 > 19.0": {
			source: "6i64 > 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19.0`"),
			},
		},
		"6i64 > 19bf": {
			source: "6i64 > 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19bf`"),
			},
		},
		"6i64 > 19f64": {
			source: "6i64 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19f64`"),
			},
		},
		"6i64 > 19f32": {
			source: "6i64 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19f32`"),
			},
		},
		"6i64 > 19i32": {
			source: "6i64 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19i32`"),
			},
		},
		"6i64 > 19i16": {
			source: "6i64 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19i16`"),
			},
		},
		"6i64 > 19i8": {
			source: "6i64 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19i8`"),
			},
		},
		"6i64 > 19u64": {
			source: "6i64 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19u64`"),
			},
		},
		"6i64 > 19u32": {
			source: "6i64 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19u32`"),
			},
		},
		"6i64 > 19u16": {
			source: "6i64 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19u16`"),
			},
		},
		"6i64 > 19u8": {
			source: "6i64 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>`, got type `19u8`"),
			},
		},

		// Int32
		"25i32 > 25i32": {
			source:       "25i32 > 25i32",
			wantStackTop: value.False.ToValue(),
		},
		"25i32 > -25i32": {
			source:       "25i32 > -25i32",
			wantStackTop: value.True.ToValue(),
		},
		"-25i32 > 25i32": {
			source:       "-25i32 > 25i32",
			wantStackTop: value.False.ToValue(),
		},
		"13i32 > 7i32": {
			source:       "13i32 > 7i32",
			wantStackTop: value.True.ToValue(),
		},
		"7i32 > 13i32": {
			source:       "7i32 > 13i32",
			wantStackTop: value.False.ToValue(),
		},

		"6i32 > 19": {
			source: "6i32 > 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19`"),
			},
		},
		"6i32 > 19.0": {
			source: "6i32 > 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19.0`"),
			},
		},
		"6i32 > 19bf": {
			source: "6i32 > 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19bf`"),
			},
		},
		"6i32 > 19f64": {
			source: "6i32 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19f64`"),
			},
		},
		"6i32 > 19f32": {
			source: "6i32 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19f32`"),
			},
		},
		"6i32 > 19i64": {
			source: "6i32 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19i64`"),
			},
		},
		"6i32 > 19i16": {
			source: "6i32 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19i16`"),
			},
		},
		"6i32 > 19i8": {
			source: "6i32 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19i8`"),
			},
		},
		"6i32 > 19u64": {
			source: "6i32 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19u64`"),
			},
		},
		"6i32 > 19u32": {
			source: "6i32 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19u32`"),
			},
		},
		"6i32 > 19u16": {
			source: "6i32 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19u16`"),
			},
		},
		"6i32 > 19u8": {
			source: "6i32 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>`, got type `19u8`"),
			},
		},

		// Int16
		"25i16 > 25i16": {
			source:       "25i16 > 25i16",
			wantStackTop: value.False.ToValue(),
		},
		"25i16 > -25i16": {
			source:       "25i16 > -25i16",
			wantStackTop: value.True.ToValue(),
		},
		"-25i16 > 25i16": {
			source:       "-25i16 > 25i16",
			wantStackTop: value.False.ToValue(),
		},
		"13i16 > 7i16": {
			source:       "13i16 > 7i16",
			wantStackTop: value.True.ToValue(),
		},
		"7i16 > 13i16": {
			source:       "7i16 > 13i16",
			wantStackTop: value.False.ToValue(),
		},
		"6i16 > 19": {
			source: "6i16 > 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19`"),
			},
		},
		"6i16 > 19.0": {
			source: "6i16 > 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19.0`"),
			},
		},
		"6i16 > 19bf": {
			source: "6i16 > 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19bf`"),
			},
		},
		"6i16 > 19f64": {
			source: "6i16 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19f64`"),
			},
		},
		"6i16 > 19f32": {
			source: "6i16 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19f32`"),
			},
		},
		"6i16 > 19i64": {
			source: "6i16 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19i64`"),
			},
		},
		"6i16 > 19i32": {
			source: "6i16 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19i32`"),
			},
		},
		"6i16 > 19i8": {
			source: "6i16 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19i8`"),
			},
		},
		"6i16 > 19u64": {
			source: "6i16 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19u64`"),
			},
		},
		"6i16 > 19u32": {
			source: "6i16 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19u32`"),
			},
		},
		"6i16 > 19u16": {
			source: "6i16 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19u16`"),
			},
		},
		"6i16 > 19u8": {
			source: "6i16 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>`, got type `19u8`"),
			},
		},

		// Int8
		"25i8 > 25i8": {
			source:       "25i8 > 25i8",
			wantStackTop: value.False.ToValue(),
		},
		"25i8 > -25i8": {
			source:       "25i8 > -25i8",
			wantStackTop: value.True.ToValue(),
		},
		"-25i8 > 25i8": {
			source:       "-25i8 > 25i8",
			wantStackTop: value.False.ToValue(),
		},
		"13i8 > 7i8": {
			source:       "13i8 > 7i8",
			wantStackTop: value.True.ToValue(),
		},
		"7i8 > 13i8": {
			source:       "7i8 > 13i8",
			wantStackTop: value.False.ToValue(),
		},
		"6i8 > 19": {
			source: "6i8 > 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(7, 1, 8)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19`"),
			},
		},
		"6i8 > 19.0": {
			source: "6i8 > 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19.0`"),
			},
		},
		"6i8 > 19bf": {
			source: "6i8 > 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19bf`"),
			},
		},
		"6i8 > 19f64": {
			source: "6i8 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19f64`"),
			},
		},
		"6i8 > 19f32": {
			source: "6i8 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19f32`"),
			},
		},
		"6i8 > 19i64": {
			source: "6i8 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19i64`"),
			},
		},
		"6i8 > 19i32": {
			source: "6i8 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19i32`"),
			},
		},
		"6i8 > 19i16": {
			source: "6i8 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19i16`"),
			},
		},
		"6i8 > 19u64": {
			source: "6i8 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19u64`"),
			},
		},
		"6i8 > 19u32": {
			source: "6i8 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19u32`"),
			},
		},
		"6i8 > 19u16": {
			source: "6i8 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19u16`"),
			},
		},
		"6i8 > 19u8": {
			source: "6i8 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>`, got type `19u8`"),
			},
		},

		// UInt64
		"25u64 > 25u64": {
			source:       "25u64 > 25u64",
			wantStackTop: value.False.ToValue(),
		},
		"13u64 > 7u64": {
			source:       "13u64 > 7u64",
			wantStackTop: value.True.ToValue(),
		},
		"7u64 > 13u64": {
			source:       "7u64 > 13u64",
			wantStackTop: value.False.ToValue(),
		},
		"6u64 > 19": {
			source: "6u64 > 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19`"),
			},
		},
		"6u64 > 19.0": {
			source: "6u64 > 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19.0`"),
			},
		},
		"6u64 > 19bf": {
			source: "6u64 > 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19bf`"),
			},
		},
		"6u64 > 19f64": {
			source: "6u64 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19f64`"),
			},
		},
		"6u64 > 19f32": {
			source: "6u64 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19f32`"),
			},
		},
		"6u64 > 19i64": {
			source: "6u64 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19i64`"),
			},
		},
		"6u64 > 19i32": {
			source: "6u64 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19i32`"),
			},
		},
		"6u64 > 19i16": {
			source: "6u64 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19i16`"),
			},
		},
		"6u64 > 19i8": {
			source: "6u64 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19i8`"),
			},
		},
		"6u64 > 19u32": {
			source: "6u64 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19u32`"),
			},
		},
		"6u64 > 19u16": {
			source: "6u64 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19u16`"),
			},
		},
		"6u64 > 19u8": {
			source: "6u64 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>`, got type `19u8`"),
			},
		},

		// UInt32
		"25u32 > 25u32": {
			source:       "25u32 > 25u32",
			wantStackTop: value.False.ToValue(),
		},
		"13u32 > 7u32": {
			source:       "13u32 > 7u32",
			wantStackTop: value.True.ToValue(),
		},
		"7u32 > 13u32": {
			source:       "7u32 > 13u32",
			wantStackTop: value.False.ToValue(),
		},
		"6u32 > 19": {
			source: "6u32 > 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19`"),
			},
		},
		"6u32 > 19.0": {
			source: "6u32 > 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19.0`"),
			},
		},
		"6u32 > 19bf": {
			source: "6u32 > 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19bf`"),
			},
		},
		"6u32 > 19f64": {
			source: "6u32 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19f64`"),
			},
		},
		"6u32 > 19f32": {
			source: "6u32 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19f32`"),
			},
		},
		"6u32 > 19i64": {
			source: "6u32 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19i64`"),
			},
		},
		"6u32 > 19i32": {
			source: "6u32 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19i32`"),
			},
		},
		"6u32 > 19i16": {
			source: "6u32 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19i16`"),
			},
		},
		"6u32 > 19i8": {
			source: "6u32 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19i8`"),
			},
		},
		"6u32 > 19u64": {
			source: "6u32 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19u64`"),
			},
		},
		"6u32 > 19u16": {
			source: "6u32 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19u16`"),
			},
		},
		"6u32 > 19u8": {
			source: "6u32 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>`, got type `19u8`"),
			},
		},

		// UInt16
		"25u16 > 25u16": {
			source:       "25u16 > 25u16",
			wantStackTop: value.False.ToValue(),
		},
		"13u16 > 7u16": {
			source:       "13u16 > 7u16",
			wantStackTop: value.True.ToValue(),
		},
		"7u16 > 13u16": {
			source:       "7u16 > 13u16",
			wantStackTop: value.False.ToValue(),
		},
		"6u16 > 19": {
			source: "6u16 > 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19`"),
			},
		},
		"6u16 > 19.0": {
			source: "6u16 > 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19.0`"),
			},
		},
		"6u16 > 19bf": {
			source: "6u16 > 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19bf`"),
			},
		},
		"6u16 > 19f64": {
			source: "6u16 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19f64`"),
			},
		},
		"6u16 > 19f32": {
			source: "6u16 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19f32`"),
			},
		},
		"6u16 > 19i64": {
			source: "6u16 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19i64`"),
			},
		},
		"6u16 > 19i32": {
			source: "6u16 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19i32`"),
			},
		},
		"6u16 > 19i16": {
			source: "6u16 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19i16`"),
			},
		},
		"6u16 > 19i8": {
			source: "6u16 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19i8`"),
			},
		},
		"6u16 > 19u64": {
			source: "6u16 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19u64`"),
			},
		},
		"6u16 > 19u32": {
			source: "6u16 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19u32`"),
			},
		},
		"6u16 > 19u8": {
			source: "6u16 > 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>`, got type `19u8`"),
			},
		},

		// Int8
		"25u8 > 25u8": {
			source:       "25u8 > 25u8",
			wantStackTop: value.False.ToValue(),
		},
		"13u8 > 7u8": {
			source:       "13u8 > 7u8",
			wantStackTop: value.True.ToValue(),
		},
		"7u8 > 13u8": {
			source:       "7u8 > 13u8",
			wantStackTop: value.False.ToValue(),
		},

		"6u8 > 19": {
			source: "6u8 > 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(7, 1, 8)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19`"),
			},
		},
		"6u8 > 19.0": {
			source: "6u8 > 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19.0`"),
			},
		},
		"6u8 > 19bf": {
			source: "6u8 > 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19bf`"),
			},
		},
		"6u8 > 19f64": {
			source: "6u8 > 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19f64`"),
			},
		},
		"6u8 > 19f32": {
			source: "6u8 > 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19f32`"),
			},
		},
		"6u8 > 19i64": {
			source: "6u8 > 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19i64`"),
			},
		},
		"6u8 > 19i32": {
			source: "6u8 > 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19i32`"),
			},
		},
		"6u8 > 19i16": {
			source: "6u8 > 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19i16`"),
			},
		},
		"6u8 > 19i8": {
			source: "6u8 > 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19i8`"),
			},
		},
		"6u8 > 19u64": {
			source: "6u8 > 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19u64`"),
			},
		},
		"6u8 > 19u32": {
			source: "6u8 > 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19u32`"),
			},
		},
		"6u8 > 19u16": {
			source: "6u8 > 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>`, got type `19u16`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_GreaterThanEqual(t *testing.T) {
	tests := sourceTestTable{
		// String
		"'25' >= '25'": {
			source:       "'25' >= '25'",
			wantStackTop: value.True.ToValue(),
		},
		"'7' >= '10'": {
			source:       "'7' >= '10'",
			wantStackTop: value.True.ToValue(),
		},
		"'10' >= '7'": {
			source:       "'10' >= '7'",
			wantStackTop: value.False.ToValue(),
		},
		"'25' >= '22'": {
			source:       "'25' >= '22'",
			wantStackTop: value.True.ToValue(),
		},
		"'22' >= '25'": {
			source:       "'22' >= '25'",
			wantStackTop: value.False.ToValue(),
		},
		"'foo' >= 'foo'": {
			source:       "'foo' >= 'foo'",
			wantStackTop: value.True.ToValue(),
		},
		"'foo' >= 'foa'": {
			source:       "'foo' >= 'foa'",
			wantStackTop: value.True.ToValue(),
		},
		"'foa' >= 'foo'": {
			source:       "'foa' >= 'foo'",
			wantStackTop: value.False.ToValue(),
		},
		"'foo' >= 'foo bar'": {
			source:       "'foo' >= 'foo bar'",
			wantStackTop: value.False.ToValue(),
		},
		"'foo bar' >= 'foo'": {
			source:       "'foo bar' >= 'foo'",
			wantStackTop: value.True.ToValue(),
		},

		"'2' >= `2`": {
			source:       "'2' >= `2`",
			wantStackTop: value.True.ToValue(),
		},
		"'72' >= `7`": {
			source:       "'72' >= `7`",
			wantStackTop: value.True.ToValue(),
		},
		"'8' >= `7`": {
			source:       "'8' >= `7`",
			wantStackTop: value.True.ToValue(),
		},
		"'7' >= `8`": {
			source:       "'7' >= `8`",
			wantStackTop: value.False.ToValue(),
		},
		"'ba' >= `b`": {
			source:       "'ba' >= `b`",
			wantStackTop: value.True.ToValue(),
		},
		"'b' >= `a`": {
			source:       "'b' >= `a`",
			wantStackTop: value.True.ToValue(),
		},
		"'a' >= `b`": {
			source:       "'a' >= `b`",
			wantStackTop: value.False.ToValue(),
		},
		"'2' >= 2.0": {
			source: "'2' >= 2.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(9, 1, 10)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `2.0`"),
			},
		},

		"'28' >= 25.2bf": {
			source: "'28' >= 25.2bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(13, 1, 14)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `25.2bf`"),
			},
		},

		"'28.8' >= 12.9f64": {
			source: "'28.8' >= 12.9f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(10, 1, 11), P(16, 1, 17)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `12.9f64`"),
			},
		},

		"'28.8' >= 12.9f32": {
			source: "'28.8' >= 12.9f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(10, 1, 11), P(16, 1, 17)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `12.9f32`"),
			},
		},

		"'93' >= 19i64": {
			source: "'93' >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `19i64`"),
			},
		},

		"'93' >= 19i32": {
			source: "'93' >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `19i32`"),
			},
		},

		"'93' >= 19i16": {
			source: "'93' >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `19i16`"),
			},
		},

		"'93' >= 19i8": {
			source: "'93' >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `19i8`"),
			},
		},

		"'93' >= 19u64": {
			source: "'93' >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `19u64`"),
			},
		},

		"'93' >= 19u32": {
			source: "'93' >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `19u32`"),
			},
		},

		"'93' >= 19u16": {
			source: "'93' >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `19u16`"),
			},
		},

		"'93' >= 19u8": {
			source: "'93' >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:>=`, got type `19u8`"),
			},
		},

		// Char
		"`2` >= `2`": {
			source:       "`2` >= `2`",
			wantStackTop: value.True.ToValue(),
		},
		"`8` >= `7`": {
			source:       "`8` >= `7`",
			wantStackTop: value.True.ToValue(),
		},
		"`7` >= `8`": {
			source:       "`7` >= `8`",
			wantStackTop: value.False.ToValue(),
		},
		"`b` >= `a`": {
			source:       "`b` >= `a`",
			wantStackTop: value.True.ToValue(),
		},
		"`a` >= `b`": {
			source:       "`a` >= `b`",
			wantStackTop: value.False.ToValue(),
		},

		"`2` >= '2'": {
			source:       "`2` >= '2'",
			wantStackTop: value.True.ToValue(),
		},
		"`7` >= '72'": {
			source:       "`7` >= '72'",
			wantStackTop: value.False.ToValue(),
		},
		"`8` >= '7'": {
			source:       "`8` >= '7'",
			wantStackTop: value.True.ToValue(),
		},
		"`7` >= '8'": {
			source:       "`7` >= '8'",
			wantStackTop: value.False.ToValue(),
		},
		"`b` >= 'a'": {
			source:       "`b` >= 'a'",
			wantStackTop: value.True.ToValue(),
		},
		"`b` >= 'ba'": {
			source:       "`b` >= 'ba'",
			wantStackTop: value.False.ToValue(),
		},
		"`a` >= 'b'": {
			source:       "`a` >= 'b'",
			wantStackTop: value.False.ToValue(),
		},
		"`2` >= 2.0": {
			source: "`2` >= 2.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(9, 1, 10)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `2.0`"),
			},
		},
		"`i` >= 25.2bf": {
			source: "`i` >= 25.2bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `25.2bf`"),
			},
		},
		"`f` >= 12.9f64": {
			source: "`f` >= 12.9f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(13, 1, 14)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `12.9f64`"),
			},
		},
		"`0` >= 12.9f32": {
			source: "`0` >= 12.9f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(13, 1, 14)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `12.9f32`"),
			},
		},
		"`9` >= 19i64": {
			source: "`9` >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `19i64`"),
			},
		},
		"`u` >= 19i32": {
			source: "`u` >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `19i32`"),
			},
		},
		"`4` >= 19i16": {
			source: "`4` >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `19i16`"),
			},
		},
		"`6` >= 19i8": {
			source: "`6` >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `19i8`"),
			},
		},
		"`9` >= 19u64": {
			source: "`9` >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `19u64`"),
			},
		},
		"`u` >= 19u32": {
			source: "`u` >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `19u32`"),
			},
		},
		"`4` >= 19u16": {
			source: "`4` >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `19u16`"),
			},
		},
		"`6` >= 19u8": {
			source: "`6` >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:>=`, got type `19u8`"),
			},
		},

		// Int
		"25 >= 25": {
			source:       "25 >= 25",
			wantStackTop: value.True.ToValue(),
		},
		"25 >= -25": {
			source:       "25 >= -25",
			wantStackTop: value.True.ToValue(),
		},
		"-25 >= 25": {
			source:       "-25 >= 25",
			wantStackTop: value.False.ToValue(),
		},
		"13 >= 7": {
			source:       "13 >= 7",
			wantStackTop: value.True.ToValue(),
		},
		"7 >= 13": {
			source:       "7 >= 13",
			wantStackTop: value.False.ToValue(),
		},

		"25 >= 25.0": {
			source:       "25 >= 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"25 >= -25.0": {
			source:       "25 >= -25.0",
			wantStackTop: value.True.ToValue(),
		},
		"-25 >= 25.0": {
			source:       "-25 >= 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"13 >= 7.0": {
			source:       "13 >= 7.0",
			wantStackTop: value.True.ToValue(),
		},
		"7 >= 13.0": {
			source:       "7 >= 13.0",
			wantStackTop: value.False.ToValue(),
		},
		"7 >= 7.5": {
			source:       "7 >= 7.5",
			wantStackTop: value.False.ToValue(),
		},
		"7 >= 6.9": {
			source:       "7 >= 6.9",
			wantStackTop: value.True.ToValue(),
		},

		"25 >= 25bf": {
			source:       "25 >= 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"25 >= -25bf": {
			source:       "25 >= -25bf",
			wantStackTop: value.True.ToValue(),
		},
		"-25 >= 25bf": {
			source:       "-25 >= 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"13 >= 7bf": {
			source:       "13 >= 7bf",
			wantStackTop: value.True.ToValue(),
		},
		"7 >= 13bf": {
			source:       "7 >= 13bf",
			wantStackTop: value.False.ToValue(),
		},
		"7 >= 7.5bf": {
			source:       "7 >= 7.5bf",
			wantStackTop: value.False.ToValue(),
		},
		"7 >= 6.9bf": {
			source:       "7 >= 6.9bf",
			wantStackTop: value.True.ToValue(),
		},
		"6 >= 19f64": {
			source: "6 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>=`, got type `19f64`"),
			},
		},
		"6 >= 19f32": {
			source: "6 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>=`, got type `19f32`"),
			},
		},
		"6 >= 19i64": {
			source: "6 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>=`, got type `19i64`"),
			},
		},
		"6 >= 19i32": {
			source: "6 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>=`, got type `19i32`"),
			},
		},
		"6 >= 19i16": {
			source: "6 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>=`, got type `19i16`"),
			},
		},
		"6 >= 19i8": {
			source: "6 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>=`, got type `19i8`"),
			},
		},
		"6 >= 19u64": {
			source: "6 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>=`, got type `19u64`"),
			},
		},
		"6 >= 19u32": {
			source: "6 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>=`, got type `19u32`"),
			},
		},
		"6 >= 19u16": {
			source: "6 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>=`, got type `19u16`"),
			},
		},
		"6 >= 19u8": {
			source: "6 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:>=`, got type `19u8`"),
			},
		},

		// Float
		"25.0 >= 25.0": {
			source:       "25.0 >= 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"25.0 >= -25.0": {
			source:       "25.0 >= -25.0",
			wantStackTop: value.True.ToValue(),
		},
		"-25.0 >= 25.0": {
			source:       "-25.0 >= 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"13.0 >= 7.0": {
			source:       "13.0 >= 7.0",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 >= 13.0": {
			source:       "7.0 >= 13.0",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 >= 7.5": {
			source:       "7.0 >= 7.5",
			wantStackTop: value.False.ToValue(),
		},
		"7.5 >= 7.0": {
			source:       "7.5 >= 7.0",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 >= 6.9": {
			source:       "7.0 >= 6.9",
			wantStackTop: value.True.ToValue(),
		},

		"25.0 >= 25": {
			source:       "25.0 >= 25",
			wantStackTop: value.True.ToValue(),
		},
		"25.0 >= -25": {
			source:       "25.0 >= -25",
			wantStackTop: value.True.ToValue(),
		},
		"-25.0 >= 25": {
			source:       "-25.0 >= 25",
			wantStackTop: value.False.ToValue(),
		},
		"13.0 >= 7": {
			source:       "13.0 >= 7",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 >= 13": {
			source:       "7.0 >= 13",
			wantStackTop: value.False.ToValue(),
		},
		"7.5 >= 7": {
			source:       "7.5 >= 7",
			wantStackTop: value.True.ToValue(),
		},

		"25.0 >= 25bf": {
			source:       "25.0 >= 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"25.0 >= -25bf": {
			source:       "25.0 >= -25bf",
			wantStackTop: value.True.ToValue(),
		},
		"-25.0 >= 25bf": {
			source:       "-25.0 >= 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"13.0 >= 7bf": {
			source:       "13.0 >= 7bf",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 >= 13bf": {
			source:       "7.0 >= 13bf",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 >= 7.5bf": {
			source:       "7.0 >= 7.5bf",
			wantStackTop: value.False.ToValue(),
		},
		"7.5 >= 7bf": {
			source:       "7.5 >= 7bf",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 >= 6.9bf": {
			source:       "7.0 >= 6.9bf",
			wantStackTop: value.True.ToValue(),
		},
		"6.0 >= 19f64": {
			source: "6.0 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>=`, got type `19f64`"),
			},
		},
		"6.0 >= 19f32": {
			source: "6.0 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>=`, got type `19f32`"),
			},
		},
		"6.0 >= 19i64": {
			source: "6.0 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>=`, got type `19i64`"),
			},
		},
		"6.0 >= 19i32": {
			source: "6.0 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>=`, got type `19i32`"),
			},
		},
		"6.0 >= 19i16": {
			source: "6.0 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>=`, got type `19i16`"),
			},
		},
		"6.0 >= 19i8": {
			source: "6.0 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>=`, got type `19i8`"),
			},
		},
		"6.0 >= 19u64": {
			source: "6.0 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>=`, got type `19u64`"),
			},
		},
		"6.0 >= 19u32": {
			source: "6.0 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>=`, got type `19u32`"),
			},
		},
		"6.0 >= 19u16": {
			source: "6.0 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>=`, got type `19u16`"),
			},
		},
		"6.0 >= 19u8": {
			source: "6.0 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:>=`, got type `19u8`"),
			},
		},

		// BigFloat
		"25bf >= 25.0": {
			source:       "25bf >= 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"25bf >= -25.0": {
			source:       "25bf >= -25.0",
			wantStackTop: value.True.ToValue(),
		},
		"-25bf >= 25.0": {
			source:       "-25bf >= 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"13bf >= 7.0": {
			source:       "13bf >= 7.0",
			wantStackTop: value.True.ToValue(),
		},
		"7bf >= 13.0": {
			source:       "7bf >= 13.0",
			wantStackTop: value.False.ToValue(),
		},
		"7bf >= 7.5": {
			source:       "7bf >= 7.5",
			wantStackTop: value.False.ToValue(),
		},
		"7.5bf >= 7.0": {
			source:       "7.5bf >= 7.0",
			wantStackTop: value.True.ToValue(),
		},
		"7bf >= 6.9": {
			source:       "7bf >= 6.9",
			wantStackTop: value.True.ToValue(),
		},

		"25bf >= 25": {
			source:       "25bf >= 25",
			wantStackTop: value.True.ToValue(),
		},
		"25bf >= -25": {
			source:       "25bf >= -25",
			wantStackTop: value.True.ToValue(),
		},
		"-25bf >= 25": {
			source:       "-25bf >= 25",
			wantStackTop: value.False.ToValue(),
		},
		"13bf >= 7": {
			source:       "13bf >= 7",
			wantStackTop: value.True.ToValue(),
		},
		"7bf >= 13": {
			source:       "7bf >= 13",
			wantStackTop: value.False.ToValue(),
		},
		"7.5bf >= 7": {
			source:       "7.5bf >= 7",
			wantStackTop: value.True.ToValue(),
		},

		"25bf >= 25bf": {
			source:       "25bf >= 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"25bf >= -25bf": {
			source:       "25bf >= -25bf",
			wantStackTop: value.True.ToValue(),
		},
		"-25bf >= 25bf": {
			source:       "-25bf >= 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"13bf >= 7bf": {
			source:       "13bf >= 7bf",
			wantStackTop: value.True.ToValue(),
		},
		"7bf >= 13bf": {
			source:       "7bf >= 13bf",
			wantStackTop: value.False.ToValue(),
		},
		"7bf >= 7.5bf": {
			source:       "7bf >= 7.5bf",
			wantStackTop: value.False.ToValue(),
		},
		"7.5bf >= 7bf": {
			source:       "7.5bf >= 7bf",
			wantStackTop: value.True.ToValue(),
		},
		"7bf >= 6.9bf": {
			source:       "7bf >= 6.9bf",
			wantStackTop: value.True.ToValue(),
		},
		"6bf >= 19f64": {
			source: "6bf >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>=`, got type `19f64`"),
			},
		},
		"6bf >= 19f32": {
			source: "6bf >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>=`, got type `19f32`"),
			},
		},
		"6bf >= 19i64": {
			source: "6bf >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>=`, got type `19i64`"),
			},
		},
		"6bf >= 19i32": {
			source: "6bf >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>=`, got type `19i32`"),
			},
		},
		"6bf >= 19i16": {
			source: "6bf >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>=`, got type `19i16`"),
			},
		},
		"6bf >= 19i8": {
			source: "6bf >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>=`, got type `19i8`"),
			},
		},
		"6bf >= 19u64": {
			source: "6bf >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>=`, got type `19u64`"),
			},
		},
		"6bf >= 19u32": {
			source: "6bf >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>=`, got type `19u32`"),
			},
		},
		"6bf >= 19u16": {
			source: "6bf >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>=`, got type `19u16`"),
			},
		},
		"6bf >= 19u8": {
			source: "6bf >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:>=`, got type `19u8`"),
			},
		},

		// Float64
		"25f64 >= 25f64": {
			source:       "25f64 >= 25f64",
			wantStackTop: value.True.ToValue(),
		},
		"25f64 >= -25f64": {
			source:       "25f64 >= -25f64",
			wantStackTop: value.True.ToValue(),
		},
		"-25f64 >= 25f64": {
			source:       "-25f64 >= 25f64",
			wantStackTop: value.False.ToValue(),
		},
		"13f64 >= 7f64": {
			source:       "13f64 >= 7f64",
			wantStackTop: value.True.ToValue(),
		},
		"7f64 >= 13f64": {
			source:       "7f64 >= 13f64",
			wantStackTop: value.False.ToValue(),
		},
		"7f64 >= 7.5f64": {
			source:       "7f64 >= 7.5f64",
			wantStackTop: value.False.ToValue(),
		},
		"7.5f64 >= 7f64": {
			source:       "7.5f64 >= 7f64",
			wantStackTop: value.True.ToValue(),
		},
		"7f64 >= 6.9f64": {
			source:       "7f64 >= 6.9f64",
			wantStackTop: value.True.ToValue(),
		},
		"6f64 >= 19.0": {
			source: "6f64 >= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19.0`"),
			},
		},
		"6f64 >= 19": {
			source: "6f64 >= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19`"),
			},
		},
		"6f64 >= 19bf": {
			source: "6f64 >= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19bf`"),
			},
		},
		"6f64 >= 19f32": {
			source: "6f64 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19f32`"),
			},
		},
		"6f64 >= 19i64": {
			source: "6f64 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19i64`"),
			},
		},
		"6f64 >= 19i32": {
			source: "6f64 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19i32`"),
			},
		},
		"6f64 >= 19i16": {
			source: "6f64 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19i16`"),
			},
		},
		"6f64 >= 19i8": {
			source: "6f64 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19i8`"),
			},
		},
		"6f64 >= 19u64": {
			source: "6f64 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19u64`"),
			},
		},
		"6f64 >= 19u32": {
			source: "6f64 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19u32`"),
			},
		},
		"6f64 >= 19u16": {
			source: "6f64 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19u16`"),
			},
		},
		"6f64 >= 19u8": {
			source: "6f64 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:>=`, got type `19u8`"),
			},
		},

		// Float32
		"25f32 >= 25f32": {
			source:       "25f32 >= 25f32",
			wantStackTop: value.True.ToValue(),
		},
		"25f32 >= -25f32": {
			source:       "25f32 >= -25f32",
			wantStackTop: value.True.ToValue(),
		},
		"-25f32 >= 25f32": {
			source:       "-25f32 >= 25f32",
			wantStackTop: value.False.ToValue(),
		},
		"13f32 >= 7f32": {
			source:       "13f32 >= 7f32",
			wantStackTop: value.True.ToValue(),
		},
		"7f32 >= 13f32": {
			source:       "7f32 >= 13f32",
			wantStackTop: value.False.ToValue(),
		},
		"7f32 >= 7.5f32": {
			source:       "7f32 >= 7.5f32",
			wantStackTop: value.False.ToValue(),
		},
		"7.5f32 >= 7f32": {
			source:       "7.5f32 >= 7f32",
			wantStackTop: value.True.ToValue(),
		},
		"7f32 >= 6.9f32": {
			source:       "7f32 >= 6.9f32",
			wantStackTop: value.True.ToValue(),
		},
		"6f32 >= 19.0": {
			source: "6f32 >= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19.0`"),
			},
		},
		"6f32 >= 19": {
			source: "6f32 >= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19`"),
			},
		},
		"6f32 >= 19bf": {
			source: "6f32 >= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19bf`"),
			},
		},
		"6f32 >= 19f64": {
			source: "6f32 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19f64`"),
			},
		},
		"6f32 >= 19i64": {
			source: "6f32 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19i64`"),
			},
		},
		"6f32 >= 19i32": {
			source: "6f32 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19i32`"),
			},
		},
		"6f32 >= 19i16": {
			source: "6f32 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19i16`"),
			},
		},
		"6f32 >= 19i8": {
			source: "6f32 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19i8`"),
			},
		},
		"6f32 >= 19u64": {
			source: "6f32 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19u64`"),
			},
		},
		"6f32 >= 19u32": {
			source: "6f32 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19u32`"),
			},
		},
		"6f32 >= 19u16": {
			source: "6f32 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19u16`"),
			},
		},
		"6f32 >= 19u8": {
			source: "6f32 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:>=`, got type `19u8`"),
			},
		},

		// Int64
		"25i64 >= 25i64": {
			source:       "25i64 >= 25i64",
			wantStackTop: value.True.ToValue(),
		},
		"25i64 >= -25i64": {
			source:       "25i64 >= -25i64",
			wantStackTop: value.True.ToValue(),
		},
		"-25i64 >= 25i64": {
			source:       "-25i64 >= 25i64",
			wantStackTop: value.False.ToValue(),
		},
		"13i64 >= 7i64": {
			source:       "13i64 >= 7i64",
			wantStackTop: value.True.ToValue(),
		},
		"7i64 >= 13i64": {
			source:       "7i64 >= 13i64",
			wantStackTop: value.False.ToValue(),
		},
		"6i64 >= 19": {
			source: "6i64 >= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19`"),
			},
		},
		"6i64 >= 19.0": {
			source: "6i64 >= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19.0`"),
			},
		},
		"6i64 >= 19bf": {
			source: "6i64 >= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19bf`"),
			},
		},
		"6i64 >= 19f64": {
			source: "6i64 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19f64`"),
			},
		},
		"6i64 >= 19f32": {
			source: "6i64 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19f32`"),
			},
		},
		"6i64 >= 19i32": {
			source: "6i64 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19i32`"),
			},
		},
		"6i64 >= 19i16": {
			source: "6i64 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19i16`"),
			},
		},
		"6i64 >= 19i8": {
			source: "6i64 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19i8`"),
			},
		},
		"6i64 >= 19u64": {
			source: "6i64 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19u64`"),
			},
		},
		"6i64 >= 19u32": {
			source: "6i64 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19u32`"),
			},
		},
		"6i64 >= 19u16": {
			source: "6i64 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19u16`"),
			},
		},
		"6i64 >= 19u8": {
			source: "6i64 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:>=`, got type `19u8`"),
			},
		},

		// Int32
		"25i32 >= 25i32": {
			source:       "25i32 >= 25i32",
			wantStackTop: value.True.ToValue(),
		},
		"25i32 >= -25i32": {
			source:       "25i32 >= -25i32",
			wantStackTop: value.True.ToValue(),
		},
		"-25i32 >= 25i32": {
			source:       "-25i32 >= 25i32",
			wantStackTop: value.False.ToValue(),
		},
		"13i32 >= 7i32": {
			source:       "13i32 >= 7i32",
			wantStackTop: value.True.ToValue(),
		},
		"7i32 >= 13i32": {
			source:       "7i32 >= 13i32",
			wantStackTop: value.False.ToValue(),
		},
		"6i32 >= 19": {
			source: "6i32 >= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19`"),
			},
		},
		"6i32 >= 19.0": {
			source: "6i32 >= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19.0`"),
			},
		},
		"6i32 >= 19bf": {
			source: "6i32 >= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19bf`"),
			},
		},
		"6i32 >= 19f64": {
			source: "6i32 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19f64`"),
			},
		},
		"6i32 >= 19f32": {
			source: "6i32 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19f32`"),
			},
		},
		"6i32 >= 19i64": {
			source: "6i32 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19i64`"),
			},
		},
		"6i32 >= 19i16": {
			source: "6i32 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19i16`"),
			},
		},
		"6i32 >= 19i8": {
			source: "6i32 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19i8`"),
			},
		},
		"6i32 >= 19u64": {
			source: "6i32 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19u64`"),
			},
		},
		"6i32 >= 19u32": {
			source: "6i32 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19u32`"),
			},
		},
		"6i32 >= 19u16": {
			source: "6i32 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19u16`"),
			},
		},
		"6i32 >= 19u8": {
			source: "6i32 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:>=`, got type `19u8`"),
			},
		},

		// Int16
		"25i16 >= 25i16": {
			source:       "25i16 >= 25i16",
			wantStackTop: value.True.ToValue(),
		},
		"25i16 >= -25i16": {
			source:       "25i16 >= -25i16",
			wantStackTop: value.True.ToValue(),
		},
		"-25i16 >= 25i16": {
			source:       "-25i16 >= 25i16",
			wantStackTop: value.False.ToValue(),
		},
		"13i16 >= 7i16": {
			source:       "13i16 >= 7i16",
			wantStackTop: value.True.ToValue(),
		},
		"7i16 >= 13i16": {
			source:       "7i16 >= 13i16",
			wantStackTop: value.False.ToValue(),
		},

		"6i16 >= 19": {
			source: "6i16 >= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19`"),
			},
		},
		"6i16 >= 19.0": {
			source: "6i16 >= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19.0`"),
			},
		},
		"6i16 >= 19bf": {
			source: "6i16 >= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19bf`"),
			},
		},
		"6i16 >= 19f64": {
			source: "6i16 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19f64`"),
			},
		},
		"6i16 >= 19f32": {
			source: "6i16 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19f32`"),
			},
		},
		"6i16 >= 19i64": {
			source: "6i16 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19i64`"),
			},
		},
		"6i16 >= 19i32": {
			source: "6i16 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19i32`"),
			},
		},
		"6i16 >= 19i8": {
			source: "6i16 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19i8`"),
			},
		},
		"6i16 >= 19u64": {
			source: "6i16 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19u64`"),
			},
		},
		"6i16 >= 19u32": {
			source: "6i16 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19u32`"),
			},
		},
		"6i16 >= 19u16": {
			source: "6i16 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19u16`"),
			},
		},
		"6i16 >= 19u8": {
			source: "6i16 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:>=`, got type `19u8`"),
			},
		},

		// Int8
		"25i8 >= 25i8": {
			source:       "25i8 >= 25i8",
			wantStackTop: value.True.ToValue(),
		},
		"25i8 >= -25i8": {
			source:       "25i8 >= -25i8",
			wantStackTop: value.True.ToValue(),
		},
		"-25i8 >= 25i8": {
			source:       "-25i8 >= 25i8",
			wantStackTop: value.False.ToValue(),
		},
		"13i8 >= 7i8": {
			source:       "13i8 >= 7i8",
			wantStackTop: value.True.ToValue(),
		},
		"7i8 >= 13i8": {
			source:       "7i8 >= 13i8",
			wantStackTop: value.False.ToValue(),
		},
		"6i8 >= 19": {
			source: "6i8 >= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19`"),
			},
		},
		"6i8 >= 19.0": {
			source: "6i8 >= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19.0`"),
			},
		},
		"6i8 >= 19bf": {
			source: "6i8 >= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19bf`"),
			},
		},
		"6i8 >= 19f64": {
			source: "6i8 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19f64`"),
			},
		},
		"6i8 >= 19f32": {
			source: "6i8 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19f32`"),
			},
		},
		"6i8 >= 19i64": {
			source: "6i8 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19i64`"),
			},
		},
		"6i8 >= 19i32": {
			source: "6i8 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19i32`"),
			},
		},
		"6i8 >= 19i16": {
			source: "6i8 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19i16`"),
			},
		},
		"6i8 >= 19u64": {
			source: "6i8 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19u64`"),
			},
		},
		"6i8 >= 19u32": {
			source: "6i8 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19u32`"),
			},
		},
		"6i8 >= 19u16": {
			source: "6i8 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19u16`"),
			},
		},
		"6i8 >= 19u8": {
			source: "6i8 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:>=`, got type `19u8`"),
			},
		},

		// UInt64
		"25u64 >= 25u64": {
			source:       "25u64 >= 25u64",
			wantStackTop: value.True.ToValue(),
		},
		"13u64 >= 7u64": {
			source:       "13u64 >= 7u64",
			wantStackTop: value.True.ToValue(),
		},
		"7u64 >= 13u64": {
			source:       "7u64 >= 13u64",
			wantStackTop: value.False.ToValue(),
		},
		"6u64 >= 19": {
			source: "6u64 >= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19`"),
			},
		},
		"6u64 >= 19.0": {
			source: "6u64 >= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19.0`"),
			},
		},
		"6u64 >= 19bf": {
			source: "6u64 >= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19bf`"),
			},
		},
		"6u64 >= 19f64": {
			source: "6u64 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19f64`"),
			},
		},
		"6u64 >= 19f32": {
			source: "6u64 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19f32`"),
			},
		},
		"6u64 >= 19i64": {
			source: "6u64 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19i64`"),
			},
		},
		"6u64 >= 19i32": {
			source: "6u64 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19i32`"),
			},
		},
		"6u64 >= 19i16": {
			source: "6u64 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19i16`"),
			},
		},
		"6u64 >= 19i8": {
			source: "6u64 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19i8`"),
			},
		},
		"6u64 >= 19u32": {
			source: "6u64 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19u32`"),
			},
		},
		"6u64 >= 19u16": {
			source: "6u64 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19u16`"),
			},
		},
		"6u64 >= 19u8": {
			source: "6u64 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:>=`, got type `19u8`"),
			},
		},

		// UInt32
		"25u32 >= 25u32": {
			source:       "25u32 >= 25u32",
			wantStackTop: value.True.ToValue(),
		},
		"13u32 >= 7u32": {
			source:       "13u32 >= 7u32",
			wantStackTop: value.True.ToValue(),
		},
		"7u32 >= 13u32": {
			source:       "7u32 >= 13u32",
			wantStackTop: value.False.ToValue(),
		},
		"6u32 >= 19": {
			source: "6u32 >= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19`"),
			},
		},
		"6u32 >= 19.0": {
			source: "6u32 >= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19.0`"),
			},
		},
		"6u32 >= 19bf": {
			source: "6u32 >= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19bf`"),
			},
		},
		"6u32 >= 19f64": {
			source: "6u32 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19f64`"),
			},
		},
		"6u32 >= 19f32": {
			source: "6u32 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19f32`"),
			},
		},
		"6u32 >= 19i64": {
			source: "6u32 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19i64`"),
			},
		},
		"6u32 >= 19i32": {
			source: "6u32 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19i32`"),
			},
		},
		"6u32 >= 19i16": {
			source: "6u32 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19i16`"),
			},
		},
		"6u32 >= 19i8": {
			source: "6u32 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19i8`"),
			},
		},
		"6u32 >= 19u64": {
			source: "6u32 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19u64`"),
			},
		},
		"6u32 >= 19u16": {
			source: "6u32 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19u16`"),
			},
		},
		"6u32 >= 19u8": {
			source: "6u32 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:>=`, got type `19u8`"),
			},
		},

		// Int16
		"25u16 >= 25u16": {
			source:       "25u16 >= 25u16",
			wantStackTop: value.True.ToValue(),
		},
		"13u16 >= 7u16": {
			source:       "13u16 >= 7u16",
			wantStackTop: value.True.ToValue(),
		},
		"7u16 >= 13u16": {
			source:       "7u16 >= 13u16",
			wantStackTop: value.False.ToValue(),
		},
		"6u16 >= 19": {
			source: "6u16 >= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19`"),
			},
		},
		"6u16 >= 19.0": {
			source: "6u16 >= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19.0`"),
			},
		},
		"6u16 >= 19bf": {
			source: "6u16 >= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19bf`"),
			},
		},
		"6u16 >= 19f64": {
			source: "6u16 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19f64`"),
			},
		},
		"6u16 >= 19f32": {
			source: "6u16 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19f32`"),
			},
		},
		"6u16 >= 19i64": {
			source: "6u16 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19i64`"),
			},
		},
		"6u16 >= 19i32": {
			source: "6u16 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19i32`"),
			},
		},
		"6u16 >= 19i16": {
			source: "6u16 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19i16`"),
			},
		},
		"6u16 >= 19i8": {
			source: "6u16 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19i8`"),
			},
		},
		"6u16 >= 19u64": {
			source: "6u16 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19u64`"),
			},
		},
		"6u16 >= 19u32": {
			source: "6u16 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19u32`"),
			},
		},
		"6u16 >= 19u8": {
			source: "6u16 >= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:>=`, got type `19u8`"),
			},
		},

		// UInt8
		"25u8 >= 25u8": {
			source:       "25u8 >= 25u8",
			wantStackTop: value.True.ToValue(),
		},
		"13u8 >= 7u8": {
			source:       "13u8 >= 7u8",
			wantStackTop: value.True.ToValue(),
		},
		"7u8 >= 13u8": {
			source:       "7u8 >= 13u8",
			wantStackTop: value.False.ToValue(),
		},
		"6u8 >= 19": {
			source: "6u8 >= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19`"),
			},
		},
		"6u8 >= 19.0": {
			source: "6u8 >= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19.0`"),
			},
		},
		"6u8 >= 19bf": {
			source: "6u8 >= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19bf`"),
			},
		},
		"6u8 >= 19f64": {
			source: "6u8 >= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19f64`"),
			},
		},
		"6u8 >= 19f32": {
			source: "6u8 >= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19f32`"),
			},
		},
		"6u8 >= 19i64": {
			source: "6u8 >= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19i64`"),
			},
		},
		"6u8 >= 19i32": {
			source: "6u8 >= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19i32`"),
			},
		},
		"6u8 >= 19i16": {
			source: "6u8 >= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19i16`"),
			},
		},
		"6u8 >= 19i8": {
			source: "6u8 >= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19i8`"),
			},
		},
		"6u8 >= 19u64": {
			source: "6u8 >= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19u64`"),
			},
		},
		"6u8 >= 19u32": {
			source: "6u8 >= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19u32`"),
			},
		},
		"6u8 >= 19u16": {
			source: "6u8 >= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:>=`, got type `19u16`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LessThan(t *testing.T) {
	tests := sourceTestTable{
		// String
		"'25' < '25'": {
			source:       "'25' < '25'",
			wantStackTop: value.False.ToValue(),
		},
		"'7' < '10'": {
			source:       "'7' < '10'",
			wantStackTop: value.False.ToValue(),
		},
		"'10' < '7'": {
			source:       "'10' < '7'",
			wantStackTop: value.True.ToValue(),
		},
		"'25' < '22'": {
			source:       "'25' < '22'",
			wantStackTop: value.False.ToValue(),
		},
		"'22' < '25'": {
			source:       "'22' < '25'",
			wantStackTop: value.True.ToValue(),
		},
		"'foo' < 'foo'": {
			source:       "'foo' < 'foo'",
			wantStackTop: value.False.ToValue(),
		},
		"'foo' < 'foa'": {
			source:       "'foo' < 'foa'",
			wantStackTop: value.False.ToValue(),
		},
		"'foa' < 'foo'": {
			source:       "'foa' < 'foo'",
			wantStackTop: value.True.ToValue(),
		},
		"'foo' < 'foo bar'": {
			source:       "'foo' < 'foo bar'",
			wantStackTop: value.True.ToValue(),
		},
		"'foo bar' < 'foo'": {
			source:       "'foo bar' < 'foo'",
			wantStackTop: value.False.ToValue(),
		},

		"'2' < `2`": {
			source:       "'2' < `2`",
			wantStackTop: value.False.ToValue(),
		},
		"'72' < `7`": {
			source:       "'72' < `7`",
			wantStackTop: value.False.ToValue(),
		},
		"'8' < `7`": {
			source:       "'8' < `7`",
			wantStackTop: value.False.ToValue(),
		},
		"'7' < `8`": {
			source:       "'7' < `8`",
			wantStackTop: value.True.ToValue(),
		},
		"'ba' < `b`": {
			source:       "'ba' < `b`",
			wantStackTop: value.False.ToValue(),
		},
		"'b' < `a`": {
			source:       "'b' < `a`",
			wantStackTop: value.False.ToValue(),
		},
		"'a' < `b`": {
			source:       "'a' < `b`",
			wantStackTop: value.True.ToValue(),
		},
		"'2' < 2.0": {
			source: "'2' < 2.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(8, 1, 9)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `2.0`"),
			},
		},

		"'28' < 25.2bf": {
			source: "'28' < 25.2bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `25.2bf`"),
			},
		},

		"'28.8' < 12.9f64": {
			source: "'28.8' < 12.9f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(9, 1, 10), P(15, 1, 16)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `12.9f64`"),
			},
		},

		"'28.8' < 12.9f32": {
			source: "'28.8' < 12.9f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(9, 1, 10), P(15, 1, 16)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `12.9f32`"),
			},
		},

		"'93' < 19i64": {
			source: "'93' < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `19i64`"),
			},
		},

		"'93' < 19i32": {
			source: "'93' < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `19i32`"),
			},
		},

		"'93' < 19i16": {
			source: "'93' < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `19i16`"),
			},
		},

		"'93' < 19i8": {
			source: "'93' < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `19i8`"),
			},
		},

		"'93' < 19u64": {
			source: "'93' < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `19u64`"),
			},
		},

		"'93' < 19u32": {
			source: "'93' < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `19u32`"),
			},
		},

		"'93' < 19u16": {
			source: "'93' < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `19u16`"),
			},
		},

		"'93' < 19u8": {
			source: "'93' < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<`, got type `19u8`"),
			},
		},

		// Char
		"`2` < `2`": {
			source:       "`2` < `2`",
			wantStackTop: value.False.ToValue(),
		},
		"`8` < `7`": {
			source:       "`8` < `7`",
			wantStackTop: value.False.ToValue(),
		},
		"`7` < `8`": {
			source:       "`7` < `8`",
			wantStackTop: value.True.ToValue(),
		},
		"`b` < `a`": {
			source:       "`b` < `a`",
			wantStackTop: value.False.ToValue(),
		},
		"`a` < `b`": {
			source:       "`a` < `b`",
			wantStackTop: value.True.ToValue(),
		},

		"`2` < '2'": {
			source:       "`2` < '2'",
			wantStackTop: value.False.ToValue(),
		},
		"`7` < '72'": {
			source:       "`7` < '72'",
			wantStackTop: value.True.ToValue(),
		},
		"`8` < '7'": {
			source:       "`8` < '7'",
			wantStackTop: value.False.ToValue(),
		},
		"`7` < '8'": {
			source:       "`7` < '8'",
			wantStackTop: value.True.ToValue(),
		},
		"`b` < 'a'": {
			source:       "`b` < 'a'",
			wantStackTop: value.False.ToValue(),
		},
		"`b` < 'ba'": {
			source:       "`b` < 'ba'",
			wantStackTop: value.True.ToValue(),
		},
		"`a` < 'b'": {
			source:       "`a` < 'b'",
			wantStackTop: value.True.ToValue(),
		},
		"`2` < 2.0": {
			source: "`2` < 2.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(8, 1, 9)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `2.0`"),
			},
		},
		"`i` < 25.2bf": {
			source: "`i` < 25.2bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `25.2bf`"),
			},
		},
		"`f` < 12.9f64": {
			source: "`f` < 12.9f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `12.9f64`"),
			},
		},
		"`0` < 12.9f32": {
			source: "`0` < 12.9f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `12.9f32`"),
			},
		},
		"`9` < 19i64": {
			source: "`9` < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `19i64`"),
			},
		},
		"`u` < 19i32": {
			source: "`u` < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `19i32`"),
			},
		},
		"`4` < 19i16": {
			source: "`4` < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `19i16`"),
			},
		},
		"`6` < 19i8": {
			source: "`6` < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `19i8`"),
			},
		},
		"`9` < 19u64": {
			source: "`9` < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `19u64`"),
			},
		},
		"`u` < 19u32": {
			source: "`u` < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `19u32`"),
			},
		},
		"`4` < 19u16": {
			source: "`4` < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `19u16`"),
			},
		},
		"`6` < 19u8": {
			source: "`6` < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<`, got type `19u8`"),
			},
		},

		// Int
		"25 < 25": {
			source:       "25 < 25",
			wantStackTop: value.False.ToValue(),
		},
		"25 < -25": {
			source:       "25 < -25",
			wantStackTop: value.False.ToValue(),
		},
		"-25 < 25": {
			source:       "-25 < 25",
			wantStackTop: value.True.ToValue(),
		},
		"13 < 7": {
			source:       "13 < 7",
			wantStackTop: value.False.ToValue(),
		},
		"7 < 13": {
			source:       "7 < 13",
			wantStackTop: value.True.ToValue(),
		},

		"25 < 25.0": {
			source:       "25 < 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"25 < -25.0": {
			source:       "25 < -25.0",
			wantStackTop: value.False.ToValue(),
		},
		"-25 < 25.0": {
			source:       "-25 < 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"13 < 7.0": {
			source:       "13 < 7.0",
			wantStackTop: value.False.ToValue(),
		},
		"7 < 13.0": {
			source:       "7 < 13.0",
			wantStackTop: value.True.ToValue(),
		},
		"7 < 7.5": {
			source:       "7 < 7.5",
			wantStackTop: value.True.ToValue(),
		},
		"7 < 6.9": {
			source:       "7 < 6.9",
			wantStackTop: value.False.ToValue(),
		},

		"25 < 25bf": {
			source:       "25 < 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"25 < -25bf": {
			source:       "25 < -25bf",
			wantStackTop: value.False.ToValue(),
		},
		"-25 < 25bf": {
			source:       "-25 < 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"13 < 7bf": {
			source:       "13 < 7bf",
			wantStackTop: value.False.ToValue(),
		},
		"7 < 13bf": {
			source:       "7 < 13bf",
			wantStackTop: value.True.ToValue(),
		},
		"7 < 7.5bf": {
			source:       "7 < 7.5bf",
			wantStackTop: value.True.ToValue(),
		},
		"7 < 6.9bf": {
			source:       "7 < 6.9bf",
			wantStackTop: value.False.ToValue(),
		},
		"6 < 19f64": {
			source: "6 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<`, got type `19f64`"),
			},
		},
		"6 < 19f32": {
			source: "6 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<`, got type `19f32`"),
			},
		},
		"6 < 19i64": {
			source: "6 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<`, got type `19i64`"),
			},
		},
		"6 < 19i32": {
			source: "6 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<`, got type `19i32`"),
			},
		},
		"6 < 19i16": {
			source: "6 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<`, got type `19i16`"),
			},
		},
		"6 < 19i8": {
			source: "6 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(7, 1, 8)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<`, got type `19i8`"),
			},
		},
		"6 < 19u64": {
			source: "6 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<`, got type `19u64`"),
			},
		},
		"6 < 19u32": {
			source: "6 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<`, got type `19u32`"),
			},
		},
		"6 < 19u16": {
			source: "6 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<`, got type `19u16`"),
			},
		},
		"6 < 19u8": {
			source: "6 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(4, 1, 5), P(7, 1, 8)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<`, got type `19u8`"),
			},
		},

		// Float
		"25.0 < 25.0": {
			source:       "25.0 < 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"25.0 < -25.0": {
			source:       "25.0 < -25.0",
			wantStackTop: value.False.ToValue(),
		},
		"-25.0 < 25.0": {
			source:       "-25.0 < 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"13.0 < 7.0": {
			source:       "13.0 < 7.0",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 < 13.0": {
			source:       "7.0 < 13.0",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 < 7.5": {
			source:       "7.0 < 7.5",
			wantStackTop: value.True.ToValue(),
		},
		"7.5 < 7.0": {
			source:       "7.5 < 7.0",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 < 6.9": {
			source:       "7.0 < 6.9",
			wantStackTop: value.False.ToValue(),
		},

		"25.0 < 25": {
			source:       "25.0 < 25",
			wantStackTop: value.False.ToValue(),
		},
		"25.0 < -25": {
			source:       "25.0 < -25",
			wantStackTop: value.False.ToValue(),
		},
		"-25.0 < 25": {
			source:       "-25.0 < 25",
			wantStackTop: value.True.ToValue(),
		},
		"13.0 < 7": {
			source:       "13.0 < 7",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 < 13": {
			source:       "7.0 < 13",
			wantStackTop: value.True.ToValue(),
		},
		"7.5 < 7": {
			source:       "7.5 < 7",
			wantStackTop: value.False.ToValue(),
		},

		"25.0 < 25bf": {
			source:       "25.0 < 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"25.0 < -25bf": {
			source:       "25.0 < -25bf",
			wantStackTop: value.False.ToValue(),
		},
		"-25.0 < 25bf": {
			source:       "-25.0 < 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"13.0 < 7bf": {
			source:       "13.0 < 7bf",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 < 13bf": {
			source:       "7.0 < 13bf",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 < 7.5bf": {
			source:       "7.0 < 7.5bf",
			wantStackTop: value.True.ToValue(),
		},
		"7.5 < 7bf": {
			source:       "7.5 < 7bf",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 < 6.9bf": {
			source:       "7.0 < 6.9bf",
			wantStackTop: value.False.ToValue(),
		},
		"6.0 < 19f64": {
			source: "6.0 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<`, got type `19f64`"),
			},
		},
		"6.0 < 19f32": {
			source: "6.0 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<`, got type `19f32`"),
			},
		},
		"6.0 < 19i64": {
			source: "6.0 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<`, got type `19i64`"),
			},
		},
		"6.0 < 19i32": {
			source: "6.0 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<`, got type `19i32`"),
			},
		},
		"6.0 < 19i16": {
			source: "6.0 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<`, got type `19i16`"),
			},
		},
		"6.0 < 19i8": {
			source: "6.0 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<`, got type `19i8`"),
			},
		},
		"6.0 < 19u64": {
			source: "6.0 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<`, got type `19u64`"),
			},
		},
		"6.0 < 19u32": {
			source: "6.0 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<`, got type `19u32`"),
			},
		},
		"6.0 < 19u16": {
			source: "6.0 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<`, got type `19u16`"),
			},
		},
		"6.0 < 19u8": {
			source: "6.0 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<`, got type `19u8`"),
			},
		},

		// BigFloat
		"25bf < 25.0": {
			source:       "25bf < 25.0",
			wantStackTop: value.False.ToValue(),
		},
		"25bf < -25.0": {
			source:       "25bf < -25.0",
			wantStackTop: value.False.ToValue(),
		},
		"-25bf < 25.0": {
			source:       "-25bf < 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"13bf < 7.0": {
			source:       "13bf < 7.0",
			wantStackTop: value.False.ToValue(),
		},
		"7bf < 13.0": {
			source:       "7bf < 13.0",
			wantStackTop: value.True.ToValue(),
		},
		"7bf < 7.5": {
			source:       "7bf < 7.5",
			wantStackTop: value.True.ToValue(),
		},
		"7.5bf < 7.0": {
			source:       "7.5bf < 7.0",
			wantStackTop: value.False.ToValue(),
		},
		"7bf < 6.9": {
			source:       "7bf < 6.9",
			wantStackTop: value.False.ToValue(),
		},

		"25bf < 25": {
			source:       "25bf < 25",
			wantStackTop: value.False.ToValue(),
		},
		"25bf < -25": {
			source:       "25bf < -25",
			wantStackTop: value.False.ToValue(),
		},
		"-25bf < 25": {
			source:       "-25bf < 25",
			wantStackTop: value.True.ToValue(),
		},
		"13bf < 7": {
			source:       "13bf < 7",
			wantStackTop: value.False.ToValue(),
		},
		"7bf < 13": {
			source:       "7bf < 13",
			wantStackTop: value.True.ToValue(),
		},
		"7.5bf < 7": {
			source:       "7.5bf < 7",
			wantStackTop: value.False.ToValue(),
		},

		"25bf < 25bf": {
			source:       "25bf < 25bf",
			wantStackTop: value.False.ToValue(),
		},
		"25bf < -25bf": {
			source:       "25bf < -25bf",
			wantStackTop: value.False.ToValue(),
		},
		"-25bf < 25bf": {
			source:       "-25bf < 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"13bf < 7bf": {
			source:       "13bf < 7bf",
			wantStackTop: value.False.ToValue(),
		},
		"7bf < 13bf": {
			source:       "7bf < 13bf",
			wantStackTop: value.True.ToValue(),
		},
		"7bf < 7.5bf": {
			source:       "7bf < 7.5bf",
			wantStackTop: value.True.ToValue(),
		},
		"7.5bf < 7bf": {
			source:       "7.5bf < 7bf",
			wantStackTop: value.False.ToValue(),
		},
		"7bf < 6.9bf": {
			source:       "7bf < 6.9bf",
			wantStackTop: value.False.ToValue(),
		},
		"6bf < 19f64": {
			source: "6bf < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<`, got type `19f64`"),
			},
		},
		"6bf < 19f32": {
			source: "6bf < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<`, got type `19f32`"),
			},
		},
		"6bf < 19i64": {
			source: "6bf < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<`, got type `19i64`"),
			},
		},
		"6bf < 19i32": {
			source: "6bf < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<`, got type `19i32`"),
			},
		},
		"6bf < 19i16": {
			source: "6bf < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<`, got type `19i16`"),
			},
		},
		"6bf < 19i8": {
			source: "6bf < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<`, got type `19i8`"),
			},
		},
		"6bf < 19u64": {
			source: "6bf < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<`, got type `19u64`"),
			},
		},
		"6bf < 19u32": {
			source: "6bf < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<`, got type `19u32`"),
			},
		},
		"6bf < 19u16": {
			source: "6bf < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<`, got type `19u16`"),
			},
		},
		"6bf < 19u8": {
			source: "6bf < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<`, got type `19u8`"),
			},
		},

		// Float64
		"25f64 < 25f64": {
			source:       "25f64 < 25f64",
			wantStackTop: value.False.ToValue(),
		},
		"25f64 < -25f64": {
			source:       "25f64 < -25f64",
			wantStackTop: value.False.ToValue(),
		},
		"-25f64 < 25f64": {
			source:       "-25f64 < 25f64",
			wantStackTop: value.True.ToValue(),
		},
		"13f64 < 7f64": {
			source:       "13f64 < 7f64",
			wantStackTop: value.False.ToValue(),
		},
		"7f64 < 13f64": {
			source:       "7f64 < 13f64",
			wantStackTop: value.True.ToValue(),
		},
		"7f64 < 7.5f64": {
			source:       "7f64 < 7.5f64",
			wantStackTop: value.True.ToValue(),
		},
		"7.5f64 < 7f64": {
			source:       "7.5f64 < 7f64",
			wantStackTop: value.False.ToValue(),
		},
		"7f64 < 6.9f64": {
			source:       "7f64 < 6.9f64",
			wantStackTop: value.False.ToValue(),
		},
		"6f64 < 19.0": {
			source: "6f64 < 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19.0`"),
			},
		},

		"6f64 < 19": {
			source: "6f64 < 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19`"),
			},
		},
		"6f64 < 19bf": {
			source: "6f64 < 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19bf`"),
			},
		},
		"6f64 < 19f32": {
			source: "6f64 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19f32`"),
			},
		},
		"6f64 < 19i64": {
			source: "6f64 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19i64`"),
			},
		},
		"6f64 < 19i32": {
			source: "6f64 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19i32`"),
			},
		},
		"6f64 < 19i16": {
			source: "6f64 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19i16`"),
			},
		},
		"6f64 < 19i8": {
			source: "6f64 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19i8`"),
			},
		},
		"6f64 < 19u64": {
			source: "6f64 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19u64`"),
			},
		},
		"6f64 < 19u32": {
			source: "6f64 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19u32`"),
			},
		},
		"6f64 < 19u16": {
			source: "6f64 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19u16`"),
			},
		},
		"6f64 < 19u8": {
			source: "6f64 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<`, got type `19u8`"),
			},
		},

		// Float32
		"25f32 < 25f32": {
			source:       "25f32 < 25f32",
			wantStackTop: value.False.ToValue(),
		},
		"25f32 < -25f32": {
			source:       "25f32 < -25f32",
			wantStackTop: value.False.ToValue(),
		},
		"-25f32 < 25f32": {
			source:       "-25f32 < 25f32",
			wantStackTop: value.True.ToValue(),
		},
		"13f32 < 7f32": {
			source:       "13f32 < 7f32",
			wantStackTop: value.False.ToValue(),
		},
		"7f32 < 13f32": {
			source:       "7f32 < 13f32",
			wantStackTop: value.True.ToValue(),
		},
		"7f32 < 7.5f32": {
			source:       "7f32 < 7.5f32",
			wantStackTop: value.True.ToValue(),
		},
		"7.5f32 < 7f32": {
			source:       "7.5f32 < 7f32",
			wantStackTop: value.False.ToValue(),
		},
		"7f32 < 6.9f32": {
			source:       "7f32 < 6.9f32",
			wantStackTop: value.False.ToValue(),
		},

		"6f32 < 19.0": {
			source: "6f32 < 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19.0`"),
			},
		},
		"6f32 < 19": {
			source: "6f32 < 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19`"),
			},
		},
		"6f32 < 19bf": {
			source: "6f32 < 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19bf`"),
			},
		},
		"6f32 < 19f64": {
			source: "6f32 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19f64`"),
			},
		},
		"6f32 < 19i64": {
			source: "6f32 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19i64`"),
			},
		},
		"6f32 < 19i32": {
			source: "6f32 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19i32`"),
			},
		},
		"6f32 < 19i16": {
			source: "6f32 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19i16`"),
			},
		},
		"6f32 < 19i8": {
			source: "6f32 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19i8`"),
			},
		},
		"6f32 < 19u64": {
			source: "6f32 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19u64`"),
			},
		},
		"6f32 < 19u32": {
			source: "6f32 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19u32`"),
			},
		},
		"6f32 < 19u16": {
			source: "6f32 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19u16`"),
			},
		},
		"6f32 < 19u8": {
			source: "6f32 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<`, got type `19u8`"),
			},
		},

		// Int64
		"25i64 < 25i64": {
			source:       "25i64 < 25i64",
			wantStackTop: value.False.ToValue(),
		},
		"25i64 < -25i64": {
			source:       "25i64 < -25i64",
			wantStackTop: value.False.ToValue(),
		},
		"-25i64 < 25i64": {
			source:       "-25i64 < 25i64",
			wantStackTop: value.True.ToValue(),
		},
		"13i64 < 7i64": {
			source:       "13i64 < 7i64",
			wantStackTop: value.False.ToValue(),
		},
		"7i64 < 13i64": {
			source:       "7i64 < 13i64",
			wantStackTop: value.True.ToValue(),
		},
		"6i64 < 19": {
			source: "6i64 < 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19`"),
			},
		},
		"6i64 < 19.0": {
			source: "6i64 < 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19.0`"),
			},
		},
		"6i64 < 19bf": {
			source: "6i64 < 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19bf`"),
			},
		},
		"6i64 < 19f64": {
			source: "6i64 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19f64`"),
			},
		},
		"6i64 < 19f32": {
			source: "6i64 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19f32`"),
			},
		},
		"6i64 < 19i32": {
			source: "6i64 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19i32`"),
			},
		},
		"6i64 < 19i16": {
			source: "6i64 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19i16`"),
			},
		},
		"6i64 < 19i8": {
			source: "6i64 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19i8`"),
			},
		},
		"6i64 < 19u64": {
			source: "6i64 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19u64`"),
			},
		},
		"6i64 < 19u32": {
			source: "6i64 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19u32`"),
			},
		},
		"6i64 < 19u16": {
			source: "6i64 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19u16`"),
			},
		},
		"6i64 < 19u8": {
			source: "6i64 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<`, got type `19u8`"),
			},
		},

		// Int32
		"25i32 < 25i32": {
			source:       "25i32 < 25i32",
			wantStackTop: value.False.ToValue(),
		},
		"25i32 < -25i32": {
			source:       "25i32 < -25i32",
			wantStackTop: value.False.ToValue(),
		},
		"-25i32 < 25i32": {
			source:       "-25i32 < 25i32",
			wantStackTop: value.True.ToValue(),
		},
		"13i32 < 7i32": {
			source:       "13i32 < 7i32",
			wantStackTop: value.False.ToValue(),
		},
		"7i32 < 13i32": {
			source:       "7i32 < 13i32",
			wantStackTop: value.True.ToValue(),
		},
		"6i32 < 19": {
			source: "6i32 < 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19`"),
			},
		},
		"6i32 < 19.0": {
			source: "6i32 < 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19.0`"),
			},
		},
		"6i32 < 19bf": {
			source: "6i32 < 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19bf`"),
			},
		},
		"6i32 < 19f64": {
			source: "6i32 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19f64`"),
			},
		},
		"6i32 < 19f32": {
			source: "6i32 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19f32`"),
			},
		},
		"6i32 < 19i64": {
			source: "6i32 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19i64`"),
			},
		},
		"6i32 < 19i16": {
			source: "6i32 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19i16`"),
			},
		},
		"6i32 < 19i8": {
			source: "6i32 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19i8`"),
			},
		},
		"6i32 < 19u64": {
			source: "6i32 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19u64`"),
			},
		},
		"6i32 < 19u32": {
			source: "6i32 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19u32`"),
			},
		},
		"6i32 < 19u16": {
			source: "6i32 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19u16`"),
			},
		},
		"6i32 < 19u8": {
			source: "6i32 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<`, got type `19u8`"),
			},
		},

		// Int16
		"25i16 < 25i16": {
			source:       "25i16 < 25i16",
			wantStackTop: value.False.ToValue(),
		},
		"25i16 < -25i16": {
			source:       "25i16 < -25i16",
			wantStackTop: value.False.ToValue(),
		},
		"-25i16 < 25i16": {
			source:       "-25i16 < 25i16",
			wantStackTop: value.True.ToValue(),
		},
		"13i16 < 7i16": {
			source:       "13i16 < 7i16",
			wantStackTop: value.False.ToValue(),
		},
		"7i16 < 13i16": {
			source:       "7i16 < 13i16",
			wantStackTop: value.True.ToValue(),
		},

		"6i16 < 19": {
			source: "6i16 < 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19`"),
			},
		},
		"6i16 < 19.0": {
			source: "6i16 < 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19.0`"),
			},
		},
		"6i16 < 19bf": {
			source: "6i16 < 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19bf`"),
			},
		},
		"6i16 < 19f64": {
			source: "6i16 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19f64`"),
			},
		},
		"6i16 < 19f32": {
			source: "6i16 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19f32`"),
			},
		},
		"6i16 < 19i64": {
			source: "6i16 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19i64`"),
			},
		},
		"6i16 < 19i32": {
			source: "6i16 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19i32`"),
			},
		},
		"6i16 < 19i8": {
			source: "6i16 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19i8`"),
			},
		},
		"6i16 < 19u64": {
			source: "6i16 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19u64`"),
			},
		},
		"6i16 < 19u32": {
			source: "6i16 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19u32`"),
			},
		},
		"6i16 < 19u16": {
			source: "6i16 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19u16`"),
			},
		},
		"6i16 < 19u8": {
			source: "6i16 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<`, got type `19u8`"),
			},
		},

		// Int8
		"25i8 < 25i8": {
			source:       "25i8 < 25i8",
			wantStackTop: value.False.ToValue(),
		},
		"25i8 < -25i8": {
			source:       "25i8 < -25i8",
			wantStackTop: value.False.ToValue(),
		},
		"-25i8 < 25i8": {
			source:       "-25i8 < 25i8",
			wantStackTop: value.True.ToValue(),
		},
		"13i8 < 7i8": {
			source:       "13i8 < 7i8",
			wantStackTop: value.False.ToValue(),
		},
		"7i8 < 13i8": {
			source:       "7i8 < 13i8",
			wantStackTop: value.True.ToValue(),
		},
		"6i8 < 19": {
			source: "6i8 < 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(7, 1, 8)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19`"),
			},
		},
		"6i8 < 19.0": {
			source: "6i8 < 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19.0`"),
			},
		},
		"6i8 < 19bf": {
			source: "6i8 < 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19bf`"),
			},
		},
		"6i8 < 19f64": {
			source: "6i8 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19f64`"),
			},
		},
		"6i8 < 19f32": {
			source: "6i8 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19f32`"),
			},
		},
		"6i8 < 19i64": {
			source: "6i8 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19i64`"),
			},
		},
		"6i8 < 19i32": {
			source: "6i8 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19i32`"),
			},
		},
		"6i8 < 19i16": {
			source: "6i8 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19i16`"),
			},
		},
		"6i8 < 19u64": {
			source: "6i8 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19u64`"),
			},
		},
		"6i8 < 19u32": {
			source: "6i8 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19u32`"),
			},
		},
		"6i8 < 19u16": {
			source: "6i8 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19u16`"),
			},
		},
		"6i8 < 19u8": {
			source: "6i8 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<`, got type `19u8`"),
			},
		},

		// UInt64
		"25u64 < 25u64": {
			source:       "25u64 < 25u64",
			wantStackTop: value.False.ToValue(),
		},
		"13u64 < 7u64": {
			source:       "13u64 < 7u64",
			wantStackTop: value.False.ToValue(),
		},
		"7u64 < 13u64": {
			source:       "7u64 < 13u64",
			wantStackTop: value.True.ToValue(),
		},
		"6u64 < 19": {
			source: "6u64 < 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19`"),
			},
		},
		"6u64 < 19.0": {
			source: "6u64 < 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19.0`"),
			},
		},
		"6u64 < 19bf": {
			source: "6u64 < 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19bf`"),
			},
		},
		"6u64 < 19f64": {
			source: "6u64 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19f64`"),
			},
		},
		"6u64 < 19f32": {
			source: "6u64 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19f32`"),
			},
		},
		"6u64 < 19i64": {
			source: "6u64 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19i64`"),
			},
		},
		"6u64 < 19i32": {
			source: "6u64 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19i32`"),
			},
		},
		"6u64 < 19i16": {
			source: "6u64 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19i16`"),
			},
		},
		"6u64 < 19i8": {
			source: "6u64 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19i8`"),
			},
		},
		"6u64 < 19u32": {
			source: "6u64 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19u32`"),
			},
		},
		"6u64 < 19u16": {
			source: "6u64 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19u16`"),
			},
		},
		"6u64 < 19u8": {
			source: "6u64 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<`, got type `19u8`"),
			},
		},

		// UInt32
		"25u32 < 25u32": {
			source:       "25u32 < 25u32",
			wantStackTop: value.False.ToValue(),
		},
		"13u32 < 7u32": {
			source:       "13u32 < 7u32",
			wantStackTop: value.False.ToValue(),
		},
		"7u32 < 13u32": {
			source:       "7u32 < 13u32",
			wantStackTop: value.True.ToValue(),
		},
		"6u32 < 19": {
			source: "6u32 < 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19`"),
			},
		},
		"6u32 < 19.0": {
			source: "6u32 < 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19.0`"),
			},
		},
		"6u32 < 19bf": {
			source: "6u32 < 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19bf`"),
			},
		},
		"6u32 < 19f64": {
			source: "6u32 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19f64`"),
			},
		},
		"6u32 < 19f32": {
			source: "6u32 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19f32`"),
			},
		},
		"6u32 < 19i64": {
			source: "6u32 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19i64`"),
			},
		},
		"6u32 < 19i32": {
			source: "6u32 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19i32`"),
			},
		},
		"6u32 < 19i16": {
			source: "6u32 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19i16`"),
			},
		},
		"6u32 < 19i8": {
			source: "6u32 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19i8`"),
			},
		},
		"6u32 < 19u64": {
			source: "6u32 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19u64`"),
			},
		},
		"6u32 < 19u16": {
			source: "6u32 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19u16`"),
			},
		},
		"6u32 < 19u8": {
			source: "6u32 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<`, got type `19u8`"),
			},
		},

		// UInt16
		"25u16 < 25u16": {
			source:       "25u16 < 25u16",
			wantStackTop: value.False.ToValue(),
		},
		"13u16 < 7u16": {
			source:       "13u16 < 7u16",
			wantStackTop: value.False.ToValue(),
		},
		"7u16 < 13u16": {
			source:       "7u16 < 13u16",
			wantStackTop: value.True.ToValue(),
		},
		"6u16 < 19": {
			source: "6u16 < 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19`"),
			},
		},
		"6u16 < 19.0": {
			source: "6u16 < 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19.0`"),
			},
		},
		"6u16 < 19bf": {
			source: "6u16 < 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19bf`"),
			},
		},
		"6u16 < 19f64": {
			source: "6u16 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19f64`"),
			},
		},
		"6u16 < 19f32": {
			source: "6u16 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19f32`"),
			},
		},
		"6u16 < 19i64": {
			source: "6u16 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19i64`"),
			},
		},
		"6u16 < 19i32": {
			source: "6u16 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19i32`"),
			},
		},
		"6u16 < 19i16": {
			source: "6u16 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19i16`"),
			},
		},
		"6u16 < 19i8": {
			source: "6u16 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19i8`"),
			},
		},
		"6u16 < 19u64": {
			source: "6u16 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19u64`"),
			},
		},
		"6u16 < 19u32": {
			source: "6u16 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19u32`"),
			},
		},
		"6u16 < 19u8": {
			source: "6u16 < 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<`, got type `19u8`"),
			},
		},

		// UInt8
		"25u8 < 25u8": {
			source:       "25u8 < 25u8",
			wantStackTop: value.False.ToValue(),
		},
		"13u8 < 7u8": {
			source:       "13u8 < 7u8",
			wantStackTop: value.False.ToValue(),
		},
		"7u8 < 13u8": {
			source:       "7u8 < 13u8",
			wantStackTop: value.True.ToValue(),
		},
		"6u8 < 19": {
			source: "6u8 < 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(7, 1, 8)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19`"),
			},
		},
		"6u8 < 19.0": {
			source: "6u8 < 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19.0`"),
			},
		},
		"6u8 < 19bf": {
			source: "6u8 < 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19bf`"),
			},
		},
		"6u8 < 19f64": {
			source: "6u8 < 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19f64`"),
			},
		},
		"6u8 < 19f32": {
			source: "6u8 < 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19f32`"),
			},
		},
		"6u8 < 19i64": {
			source: "6u8 < 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19i64`"),
			},
		},
		"6u8 < 19i32": {
			source: "6u8 < 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19i32`"),
			},
		},
		"6u8 < 19i16": {
			source: "6u8 < 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19i16`"),
			},
		},
		"6u8 < 19i8": {
			source: "6u8 < 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(9, 1, 10)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19i8`"),
			},
		},
		"6u8 < 19u64": {
			source: "6u8 < 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19u64`"),
			},
		},
		"6u8 < 19u32": {
			source: "6u8 < 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19u32`"),
			},
		},
		"6u8 < 19u16": {
			source: "6u8 < 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(6, 1, 7), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<`, got type `19u16`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LessThanEqual(t *testing.T) {
	tests := sourceTestTable{
		// String
		"'25' <= '25'": {
			source:       "'25' <= '25'",
			wantStackTop: value.True.ToValue(),
		},
		"'7' <= '10'": {
			source:       "'7' <= '10'",
			wantStackTop: value.False.ToValue(),
		},
		"'10' <= '7'": {
			source:       "'10' <= '7'",
			wantStackTop: value.True.ToValue(),
		},
		"'25' <= '22'": {
			source:       "'25' <= '22'",
			wantStackTop: value.False.ToValue(),
		},
		"'22' <= '25'": {
			source:       "'22' <= '25'",
			wantStackTop: value.True.ToValue(),
		},
		"'foo' <= 'foo'": {
			source:       "'foo' <= 'foo'",
			wantStackTop: value.True.ToValue(),
		},
		"'foo' <= 'foa'": {
			source:       "'foo' <= 'foa'",
			wantStackTop: value.False.ToValue(),
		},
		"'foa' <= 'foo'": {
			source:       "'foa' <= 'foo'",
			wantStackTop: value.True.ToValue(),
		},
		"'foo' <= 'foo bar'": {
			source:       "'foo' <= 'foo bar'",
			wantStackTop: value.True.ToValue(),
		},
		"'foo bar' <= 'foo'": {
			source:       "'foo bar' <= 'foo'",
			wantStackTop: value.False.ToValue(),
		},

		"'2' <= `2`": {
			source:       "'2' <= `2`",
			wantStackTop: value.True.ToValue(),
		},
		"'72' <= `7`": {
			source:       "'72' <= `7`",
			wantStackTop: value.False.ToValue(),
		},
		"'8' <= `7`": {
			source:       "'8' <= `7`",
			wantStackTop: value.False.ToValue(),
		},
		"'7' <= `8`": {
			source:       "'7' <= `8`",
			wantStackTop: value.True.ToValue(),
		},
		"'ba' <= `b`": {
			source:       "'ba' <= `b`",
			wantStackTop: value.False.ToValue(),
		},
		"'b' <= `a`": {
			source:       "'b' <= `a`",
			wantStackTop: value.False.ToValue(),
		},
		"'a' <= `b`": {
			source:       "'a' <= `b`",
			wantStackTop: value.True.ToValue(),
		},
		"'2' <= 2.0": {
			source: "'2' <= 2.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(9, 1, 10)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `2.0`"),
			},
		},

		"'28' <= 25.2bf": {
			source: "'28' <= 25.2bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(13, 1, 14)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `25.2bf`"),
			},
		},

		"'28.8' <= 12.9f64": {
			source: "'28.8' <= 12.9f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(10, 1, 11), P(16, 1, 17)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `12.9f64`"),
			},
		},

		"'28.8' <= 12.9f32": {
			source: "'28.8' <= 12.9f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(10, 1, 11), P(16, 1, 17)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `12.9f32`"),
			},
		},

		"'93' <= 19i64": {
			source: "'93' <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `19i64`"),
			},
		},

		"'93' <= 19i32": {
			source: "'93' <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `19i32`"),
			},
		},

		"'93' <= 19i16": {
			source: "'93' <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `19i16`"),
			},
		},

		"'93' <= 19i8": {
			source: "'93' <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `19i8`"),
			},
		},

		"'93' <= 19u64": {
			source: "'93' <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `19u64`"),
			},
		},

		"'93' <= 19u32": {
			source: "'93' <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `19u32`"),
			},
		},

		"'93' <= 19u16": {
			source: "'93' <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `19u16`"),
			},
		},

		"'93' <= 19u8": {
			source: "'93' <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::String.:<=`, got type `19u8`"),
			},
		},

		// Char
		"`2` <= `2`": {
			source:       "`2` <= `2`",
			wantStackTop: value.True.ToValue(),
		},
		"`8` <= `7`": {
			source:       "`8` <= `7`",
			wantStackTop: value.False.ToValue(),
		},
		"`7` <= `8`": {
			source:       "`7` <= `8`",
			wantStackTop: value.True.ToValue(),
		},
		"`b` <= `a`": {
			source:       "`b` <= `a`",
			wantStackTop: value.False.ToValue(),
		},
		"`a` <= `b`": {
			source:       "`a` <= `b`",
			wantStackTop: value.True.ToValue(),
		},

		"`2` <= '2'": {
			source:       "`2` <= '2'",
			wantStackTop: value.True.ToValue(),
		},
		"`7` <= '72'": {
			source:       "`7` <= '72'",
			wantStackTop: value.True.ToValue(),
		},
		"`8` <= '7'": {
			source:       "`8` <= '7'",
			wantStackTop: value.False.ToValue(),
		},
		"`7` <= '8'": {
			source:       "`7` <= '8'",
			wantStackTop: value.True.ToValue(),
		},
		"`b` <= 'a'": {
			source:       "`b` <= 'a'",
			wantStackTop: value.False.ToValue(),
		},
		"`b` <= 'ba'": {
			source:       "`b` <= 'ba'",
			wantStackTop: value.True.ToValue(),
		},
		"`a` <= 'b'": {
			source:       "`a` <= 'b'",
			wantStackTop: value.True.ToValue(),
		},
		"`2` <= 2.0": {
			source: "`2` <= 2.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(9, 1, 10)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `2.0`"),
			},
		},
		"`i` <= 25.2bf": {
			source: "`i` <= 25.2bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(12, 1, 13)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `25.2bf`"),
			},
		},
		"`f` <= 12.9f64": {
			source: "`f` <= 12.9f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(13, 1, 14)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `12.9f64`"),
			},
		},
		"`0` <= 12.9f32": {
			source: "`0` <= 12.9f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(13, 1, 14)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `12.9f32`"),
			},
		},
		"`9` <= 19i64": {
			source: "`9` <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `19i64`"),
			},
		},
		"`u` <= 19i32": {
			source: "`u` <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `19i32`"),
			},
		},
		"`4` <= 19i16": {
			source: "`4` <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `19i16`"),
			},
		},
		"`6` <= 19i8": {
			source: "`6` <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `19i8`"),
			},
		},
		"`9` <= 19u64": {
			source: "`9` <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `19u64`"),
			},
		},
		"`u` <= 19u32": {
			source: "`u` <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `19u32`"),
			},
		},
		"`4` <= 19u16": {
			source: "`4` <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `19u16`"),
			},
		},
		"`6` <= 19u8": {
			source: "`6` <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::String | Std::Char` for parameter `other` in call to `Std::Char.:<=`, got type `19u8`"),
			},
		},

		// Int
		"25 <= 25": {
			source:       "25 <= 25",
			wantStackTop: value.True.ToValue(),
		},
		"25 <= -25": {
			source:       "25 <= -25",
			wantStackTop: value.False.ToValue(),
		},
		"-25 <= 25": {
			source:       "-25 <= 25",
			wantStackTop: value.True.ToValue(),
		},
		"13 <= 7": {
			source:       "13 <= 7",
			wantStackTop: value.False.ToValue(),
		},
		"7 <= 13": {
			source:       "7 <= 13",
			wantStackTop: value.True.ToValue(),
		},

		"25 <= 25.0": {
			source:       "25 <= 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"25 <= -25.0": {
			source:       "25 <= -25.0",
			wantStackTop: value.False.ToValue(),
		},
		"-25 <= 25.0": {
			source:       "-25 <= 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"13 <= 7.0": {
			source:       "13 <= 7.0",
			wantStackTop: value.False.ToValue(),
		},
		"7 <= 13.0": {
			source:       "7 <= 13.0",
			wantStackTop: value.True.ToValue(),
		},
		"7 <= 7.5": {
			source:       "7 <= 7.5",
			wantStackTop: value.True.ToValue(),
		},
		"7 <= 6.9": {
			source:       "7 <= 6.9",
			wantStackTop: value.False.ToValue(),
		},

		"25 <= 25bf": {
			source:       "25 <= 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"25 <= -25bf": {
			source:       "25 <= -25bf",
			wantStackTop: value.False.ToValue(),
		},
		"-25 <= 25bf": {
			source:       "-25 <= 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"13 <= 7bf": {
			source:       "13 <= 7bf",
			wantStackTop: value.False.ToValue(),
		},
		"7 <= 13bf": {
			source:       "7 <= 13bf",
			wantStackTop: value.True.ToValue(),
		},
		"7 <= 7.5bf": {
			source:       "7 <= 7.5bf",
			wantStackTop: value.True.ToValue(),
		},
		"7 <= 6.9bf": {
			source:       "7 <= 6.9bf",
			wantStackTop: value.False.ToValue(),
		},
		"6 <= 19f64": {
			source: "6 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<=`, got type `19f64`"),
			},
		},
		"6 <= 19f32": {
			source: "6 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<=`, got type `19f32`"),
			},
		},
		"6 <= 19i64": {
			source: "6 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<=`, got type `19i64`"),
			},
		},
		"6 <= 19i32": {
			source: "6 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<=`, got type `19i32`"),
			},
		},
		"6 <= 19i16": {
			source: "6 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<=`, got type `19i16`"),
			},
		},
		"6 <= 19i8": {
			source: "6 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<=`, got type `19i8`"),
			},
		},
		"6 <= 19u64": {
			source: "6 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<=`, got type `19u64`"),
			},
		},
		"6 <= 19u32": {
			source: "6 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<=`, got type `19u32`"),
			},
		},
		"6 <= 19u16": {
			source: "6 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<=`, got type `19u16`"),
			},
		},
		"6 <= 19u8": {
			source: "6 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 1, 6), P(8, 1, 9)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Int.:<=`, got type `19u8`"),
			},
		},

		// Float
		"25.0 <= 25.0": {
			source:       "25.0 <= 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"25.0 <= -25.0": {
			source:       "25.0 <= -25.0",
			wantStackTop: value.False.ToValue(),
		},
		"-25.0 <= 25.0": {
			source:       "-25.0 <= 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"13.0 <= 7.0": {
			source:       "13.0 <= 7.0",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 <= 13.0": {
			source:       "7.0 <= 13.0",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 <= 7.5": {
			source:       "7.0 <= 7.5",
			wantStackTop: value.True.ToValue(),
		},
		"7.5 <= 7.0": {
			source:       "7.5 <= 7.0",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 <= 6.9": {
			source:       "7.0 <= 6.9",
			wantStackTop: value.False.ToValue(),
		},

		"25.0 <= 25": {
			source:       "25.0 <= 25",
			wantStackTop: value.True.ToValue(),
		},
		"25.0 <= -25": {
			source:       "25.0 <= -25",
			wantStackTop: value.False.ToValue(),
		},
		"-25.0 <= 25": {
			source:       "-25.0 <= 25",
			wantStackTop: value.True.ToValue(),
		},
		"13.0 <= 7": {
			source:       "13.0 <= 7",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 <= 13": {
			source:       "7.0 <= 13",
			wantStackTop: value.True.ToValue(),
		},
		"7.5 <= 7": {
			source:       "7.5 <= 7",
			wantStackTop: value.False.ToValue(),
		},

		"25.0 <= 25bf": {
			source:       "25.0 <= 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"25.0 <= -25bf": {
			source:       "25.0 <= -25bf",
			wantStackTop: value.False.ToValue(),
		},
		"-25.0 <= 25bf": {
			source:       "-25.0 <= 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"13.0 <= 7bf": {
			source:       "13.0 <= 7bf",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 <= 13bf": {
			source:       "7.0 <= 13bf",
			wantStackTop: value.True.ToValue(),
		},
		"7.0 <= 7.5bf": {
			source:       "7.0 <= 7.5bf",
			wantStackTop: value.True.ToValue(),
		},
		"7.5 <= 7bf": {
			source:       "7.5 <= 7bf",
			wantStackTop: value.False.ToValue(),
		},
		"7.0 <= 6.9bf": {
			source:       "7.0 <= 6.9bf",
			wantStackTop: value.False.ToValue(),
		},
		"6.0 <= 19f64": {
			source: "6.0 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<=`, got type `19f64`"),
			},
		},
		"6.0 <= 19f32": {
			source: "6.0 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<=`, got type `19f32`"),
			},
		},
		"6.0 <= 19i64": {
			source: "6.0 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<=`, got type `19i64`"),
			},
		},
		"6.0 <= 19i32": {
			source: "6.0 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<=`, got type `19i32`"),
			},
		},
		"6.0 <= 19i16": {
			source: "6.0 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<=`, got type `19i16`"),
			},
		},
		"6.0 <= 19i8": {
			source: "6.0 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<=`, got type `19i8`"),
			},
		},
		"6.0 <= 19u64": {
			source: "6.0 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<=`, got type `19u64`"),
			},
		},
		"6.0 <= 19u32": {
			source: "6.0 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<=`, got type `19u32`"),
			},
		},
		"6.0 <= 19u16": {
			source: "6.0 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<=`, got type `19u16`"),
			},
		},
		"6.0 <= 19u8": {
			source: "6.0 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::Float.:<=`, got type `19u8`"),
			},
		},

		// BigFloat
		"25bf <= 25.0": {
			source:       "25bf <= 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"25bf <= -25.0": {
			source:       "25bf <= -25.0",
			wantStackTop: value.False.ToValue(),
		},
		"-25bf <= 25.0": {
			source:       "-25bf <= 25.0",
			wantStackTop: value.True.ToValue(),
		},
		"13bf <= 7.0": {
			source:       "13bf <= 7.0",
			wantStackTop: value.False.ToValue(),
		},
		"7bf <= 13.0": {
			source:       "7bf <= 13.0",
			wantStackTop: value.True.ToValue(),
		},
		"7bf <= 7.5": {
			source:       "7bf <= 7.5",
			wantStackTop: value.True.ToValue(),
		},
		"7.5bf <= 7.0": {
			source:       "7.5bf <= 7.0",
			wantStackTop: value.False.ToValue(),
		},
		"7bf <= 6.9": {
			source:       "7bf <= 6.9",
			wantStackTop: value.False.ToValue(),
		},

		"25bf <= 25": {
			source:       "25bf <= 25",
			wantStackTop: value.True.ToValue(),
		},
		"25bf <= -25": {
			source:       "25bf <= -25",
			wantStackTop: value.False.ToValue(),
		},
		"-25bf <= 25": {
			source:       "-25bf <= 25",
			wantStackTop: value.True.ToValue(),
		},
		"13bf <= 7": {
			source:       "13bf <= 7",
			wantStackTop: value.False.ToValue(),
		},
		"7bf <= 13": {
			source:       "7bf <= 13",
			wantStackTop: value.True.ToValue(),
		},
		"7.5bf <= 7": {
			source:       "7.5bf <= 7",
			wantStackTop: value.False.ToValue(),
		},

		"25bf <= 25bf": {
			source:       "25bf <= 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"25bf <= -25bf": {
			source:       "25bf <= -25bf",
			wantStackTop: value.False.ToValue(),
		},
		"-25bf <= 25bf": {
			source:       "-25bf <= 25bf",
			wantStackTop: value.True.ToValue(),
		},
		"13bf <= 7bf": {
			source:       "13bf <= 7bf",
			wantStackTop: value.False.ToValue(),
		},
		"7bf <= 13bf": {
			source:       "7bf <= 13bf",
			wantStackTop: value.True.ToValue(),
		},
		"7bf <= 7.5bf": {
			source:       "7bf <= 7.5bf",
			wantStackTop: value.True.ToValue(),
		},
		"7.5bf <= 7bf": {
			source:       "7.5bf <= 7bf",
			wantStackTop: value.False.ToValue(),
		},
		"7bf <= 6.9bf": {
			source:       "7bf <= 6.9bf",
			wantStackTop: value.False.ToValue(),
		},
		"6bf <= 19f64": {
			source: "6bf <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<=`, got type `19f64`"),
			},
		},
		"6bf <= 19f32": {
			source: "6bf <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<=`, got type `19f32`"),
			},
		},
		"6bf <= 19i64": {
			source: "6bf <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<=`, got type `19i64`"),
			},
		},
		"6bf <= 19i32": {
			source: "6bf <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<=`, got type `19i32`"),
			},
		},
		"6bf <= 19i16": {
			source: "6bf <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<=`, got type `19i16`"),
			},
		},
		"6bf <= 19i8": {
			source: "6bf <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<=`, got type `19i8`"),
			},
		},
		"6bf <= 19u64": {
			source: "6bf <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<=`, got type `19u64`"),
			},
		},
		"6bf <= 19u32": {
			source: "6bf <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<=`, got type `19u32`"),
			},
		},
		"6bf <= 19u16": {
			source: "6bf <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<=`, got type `19u16`"),
			},
		},
		"6bf <= 19u8": {
			source: "6bf <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::CoercibleNumeric` for parameter `other` in call to `Std::BigFloat.:<=`, got type `19u8`"),
			},
		},

		// Float64
		"25f64 <= 25f64": {
			source:       "25f64 <= 25f64",
			wantStackTop: value.True.ToValue(),
		},
		"25f64 <= -25f64": {
			source:       "25f64 <= -25f64",
			wantStackTop: value.False.ToValue(),
		},
		"-25f64 <= 25f64": {
			source:       "-25f64 <= 25f64",
			wantStackTop: value.True.ToValue(),
		},
		"13f64 <= 7f64": {
			source:       "13f64 <= 7f64",
			wantStackTop: value.False.ToValue(),
		},
		"7f64 <= 13f64": {
			source:       "7f64 <= 13f64",
			wantStackTop: value.True.ToValue(),
		},
		"7f64 <= 7.5f64": {
			source:       "7f64 <= 7.5f64",
			wantStackTop: value.True.ToValue(),
		},
		"7.5f64 <= 7f64": {
			source:       "7.5f64 <= 7f64",
			wantStackTop: value.False.ToValue(),
		},
		"7f64 <= 6.9f64": {
			source:       "7f64 <= 6.9f64",
			wantStackTop: value.False.ToValue(),
		},
		"6f64 <= 19.0": {
			source: "6f64 <= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19.0`"),
			},
		},
		"6f64 <= 19": {
			source: "6f64 <= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19`"),
			},
		},
		"6f64 <= 19bf": {
			source: "6f64 <= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19bf`"),
			},
		},
		"6f64 <= 19f32": {
			source: "6f64 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19f32`"),
			},
		},
		"6f64 <= 19i64": {
			source: "6f64 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19i64`"),
			},
		},
		"6f64 <= 19i32": {
			source: "6f64 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19i32`"),
			},
		},
		"6f64 <= 19i16": {
			source: "6f64 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19i16`"),
			},
		},
		"6f64 <= 19i8": {
			source: "6f64 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19i8`"),
			},
		},
		"6f64 <= 19u64": {
			source: "6f64 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19u64`"),
			},
		},
		"6f64 <= 19u32": {
			source: "6f64 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19u32`"),
			},
		},
		"6f64 <= 19u16": {
			source: "6f64 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19u16`"),
			},
		},
		"6f64 <= 19u8": {
			source: "6f64 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float64` for parameter `other` in call to `Std::Float64.:<=`, got type `19u8`"),
			},
		},

		// Float32
		"25f32 <= 25f32": {
			source:       "25f32 <= 25f32",
			wantStackTop: value.True.ToValue(),
		},
		"25f32 <= -25f32": {
			source:       "25f32 <= -25f32",
			wantStackTop: value.False.ToValue(),
		},
		"-25f32 <= 25f32": {
			source:       "-25f32 <= 25f32",
			wantStackTop: value.True.ToValue(),
		},
		"13f32 <= 7f32": {
			source:       "13f32 <= 7f32",
			wantStackTop: value.False.ToValue(),
		},
		"7f32 <= 13f32": {
			source:       "7f32 <= 13f32",
			wantStackTop: value.True.ToValue(),
		},
		"7f32 <= 7.5f32": {
			source:       "7f32 <= 7.5f32",
			wantStackTop: value.True.ToValue(),
		},
		"7.5f32 <= 7f32": {
			source:       "7.5f32 <= 7f32",
			wantStackTop: value.False.ToValue(),
		},
		"7f32 <= 6.9f32": {
			source:       "7f32 <= 6.9f32",
			wantStackTop: value.False.ToValue(),
		},
		"6f32 <= 19.0": {
			source: "6f32 <= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19.0`"),
			},
		},
		"6f32 <= 19": {
			source: "6f32 <= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19`"),
			},
		},
		"6f32 <= 19bf": {
			source: "6f32 <= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19bf`"),
			},
		},
		"6f32 <= 19f64": {
			source: "6f32 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19f64`"),
			},
		},
		"6f32 <= 19i64": {
			source: "6f32 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19i64`"),
			},
		},
		"6f32 <= 19i32": {
			source: "6f32 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19i32`"),
			},
		},
		"6f32 <= 19i16": {
			source: "6f32 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19i16`"),
			},
		},
		"6f32 <= 19i8": {
			source: "6f32 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19i8`"),
			},
		},
		"6f32 <= 19u64": {
			source: "6f32 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19u64`"),
			},
		},
		"6f32 <= 19u32": {
			source: "6f32 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19u32`"),
			},
		},
		"6f32 <= 19u16": {
			source: "6f32 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19u16`"),
			},
		},
		"6f32 <= 19u8": {
			source: "6f32 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Float32` for parameter `other` in call to `Std::Float32.:<=`, got type `19u8`"),
			},
		},

		// Int64
		"25i64 <= 25i64": {
			source:       "25i64 <= 25i64",
			wantStackTop: value.True.ToValue(),
		},
		"25i64 <= -25i64": {
			source:       "25i64 <= -25i64",
			wantStackTop: value.False.ToValue(),
		},
		"-25i64 <= 25i64": {
			source:       "-25i64 <= 25i64",
			wantStackTop: value.True.ToValue(),
		},
		"13i64 <= 7i64": {
			source:       "13i64 <= 7i64",
			wantStackTop: value.False.ToValue(),
		},
		"7i64 <= 13i64": {
			source:       "7i64 <= 13i64",
			wantStackTop: value.True.ToValue(),
		},
		"6i64 <= 19": {
			source: "6i64 <= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19`"),
			},
		},
		"6i64 <= 19.0": {
			source: "6i64 <= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19.0`"),
			},
		},
		"6i64 <= 19bf": {
			source: "6i64 <= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19bf`"),
			},
		},
		"6i64 <= 19f64": {
			source: "6i64 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19f64`"),
			},
		},
		"6i64 <= 19f32": {
			source: "6i64 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19f32`"),
			},
		},
		"6i64 <= 19i32": {
			source: "6i64 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19i32`"),
			},
		},
		"6i64 <= 19i16": {
			source: "6i64 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19i16`"),
			},
		},
		"6i64 <= 19i8": {
			source: "6i64 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19i8`"),
			},
		},
		"6i64 <= 19u64": {
			source: "6i64 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19u64`"),
			},
		},
		"6i64 <= 19u32": {
			source: "6i64 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19u32`"),
			},
		},
		"6i64 <= 19u16": {
			source: "6i64 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19u16`"),
			},
		},
		"6i64 <= 19u8": {
			source: "6i64 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `Std::Int64.:<=`, got type `19u8`"),
			},
		},

		// Int32
		"25i32 <= 25i32": {
			source:       "25i32 <= 25i32",
			wantStackTop: value.True.ToValue(),
		},
		"25i32 <= -25i32": {
			source:       "25i32 <= -25i32",
			wantStackTop: value.False.ToValue(),
		},
		"-25i32 <= 25i32": {
			source:       "-25i32 <= 25i32",
			wantStackTop: value.True.ToValue(),
		},
		"13i32 <= 7i32": {
			source:       "13i32 <= 7i32",
			wantStackTop: value.False.ToValue(),
		},
		"7i32 <= 13i32": {
			source:       "7i32 <= 13i32",
			wantStackTop: value.True.ToValue(),
		},
		"6i32 <= 19": {
			source: "6i32 <= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19`"),
			},
		},
		"6i32 <= 19.0": {
			source: "6i32 <= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19.0`"),
			},
		},
		"6i32 <= 19bf": {
			source: "6i32 <= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19bf`"),
			},
		},
		"6i32 <= 19f64": {
			source: "6i32 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19f64`"),
			},
		},
		"6i32 <= 19f32": {
			source: "6i32 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19f32`"),
			},
		},
		"6i32 <= 19i64": {
			source: "6i32 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19i64`"),
			},
		},
		"6i32 <= 19i16": {
			source: "6i32 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19i16`"),
			},
		},
		"6i32 <= 19i8": {
			source: "6i32 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19i8`"),
			},
		},
		"6i32 <= 19u64": {
			source: "6i32 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19u64`"),
			},
		},
		"6i32 <= 19u32": {
			source: "6i32 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19u32`"),
			},
		},
		"6i32 <= 19u16": {
			source: "6i32 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19u16`"),
			},
		},
		"6i32 <= 19u8": {
			source: "6i32 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int32` for parameter `other` in call to `Std::Int32.:<=`, got type `19u8`"),
			},
		},

		// Int16
		"25i16 <= 25i16": {
			source:       "25i16 <= 25i16",
			wantStackTop: value.True.ToValue(),
		},
		"25i16 <= -25i16": {
			source:       "25i16 <= -25i16",
			wantStackTop: value.False.ToValue(),
		},
		"-25i16 <= 25i16": {
			source:       "-25i16 <= 25i16",
			wantStackTop: value.True.ToValue(),
		},
		"13i16 <= 7i16": {
			source:       "13i16 <= 7i16",
			wantStackTop: value.False.ToValue(),
		},
		"7i16 <= 13i16": {
			source:       "7i16 <= 13i16",
			wantStackTop: value.True.ToValue(),
		},
		"6i16 <= 19": {
			source: "6i16 <= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19`"),
			},
		},
		"6i16 <= 19.0": {
			source: "6i16 <= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19.0`"),
			},
		},
		"6i16 <= 19bf": {
			source: "6i16 <= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19bf`"),
			},
		},
		"6i16 <= 19f64": {
			source: "6i16 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19f64`"),
			},
		},
		"6i16 <= 19f32": {
			source: "6i16 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19f32`"),
			},
		},
		"6i16 <= 19i64": {
			source: "6i16 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19i64`"),
			},
		},
		"6i16 <= 19i32": {
			source: "6i16 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19i32`"),
			},
		},
		"6i16 <= 19i8": {
			source: "6i16 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19i8`"),
			},
		},
		"6i16 <= 19u64": {
			source: "6i16 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19u64`"),
			},
		},
		"6i16 <= 19u32": {
			source: "6i16 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19u32`"),
			},
		},
		"6i16 <= 19u16": {
			source: "6i16 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19u16`"),
			},
		},
		"6i16 <= 19u8": {
			source: "6i16 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int16` for parameter `other` in call to `Std::Int16.:<=`, got type `19u8`"),
			},
		},

		// Int8
		"25i8 <= 25i8": {
			source:       "25i8 <= 25i8",
			wantStackTop: value.True.ToValue(),
		},
		"25i8 <= -25i8": {
			source:       "25i8 <= -25i8",
			wantStackTop: value.False.ToValue(),
		},
		"-25i8 <= 25i8": {
			source:       "-25i8 <= 25i8",
			wantStackTop: value.True.ToValue(),
		},
		"13i8 <= 7i8": {
			source:       "13i8 <= 7i8",
			wantStackTop: value.False.ToValue(),
		},
		"7i8 <= 13i8": {
			source:       "7i8 <= 13i8",
			wantStackTop: value.True.ToValue(),
		},
		"6i8 <= 19": {
			source: "6i8 <= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19`"),
			},
		},
		"6i8 <= 19.0": {
			source: "6i8 <= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19.0`"),
			},
		},
		"6i8 <= 19bf": {
			source: "6i8 <= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19bf`"),
			},
		},
		"6i8 <= 19f64": {
			source: "6i8 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19f64`"),
			},
		},
		"6i8 <= 19f32": {
			source: "6i8 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19f32`"),
			},
		},
		"6i8 <= 19i64": {
			source: "6i8 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19i64`"),
			},
		},
		"6i8 <= 19i32": {
			source: "6i8 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19i32`"),
			},
		},
		"6i8 <= 19i16": {
			source: "6i8 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19i16`"),
			},
		},
		"6i8 <= 19u64": {
			source: "6i8 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19u64`"),
			},
		},
		"6i8 <= 19u32": {
			source: "6i8 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19u32`"),
			},
		},
		"6i8 <= 19u16": {
			source: "6i8 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19u16`"),
			},
		},
		"6i8 <= 19u8": {
			source: "6i8 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::Int8` for parameter `other` in call to `Std::Int8.:<=`, got type `19u8`"),
			},
		},

		// UInt64
		"25u64 <= 25u64": {
			source:       "25u64 <= 25u64",
			wantStackTop: value.True.ToValue(),
		},
		"13u64 <= 7u64": {
			source:       "13u64 <= 7u64",
			wantStackTop: value.False.ToValue(),
		},
		"7u64 <= 13u64": {
			source:       "7u64 <= 13u64",
			wantStackTop: value.True.ToValue(),
		},
		"6u64 <= 19": {
			source: "6u64 <= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19`"),
			},
		},
		"6u64 <= 19.0": {
			source: "6u64 <= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19.0`"),
			},
		},
		"6u64 <= 19bf": {
			source: "6u64 <= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19bf`"),
			},
		},
		"6u64 <= 19f64": {
			source: "6u64 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19f64`"),
			},
		},
		"6u64 <= 19f32": {
			source: "6u64 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19f32`"),
			},
		},
		"6u64 <= 19i64": {
			source: "6u64 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19i64`"),
			},
		},
		"6u64 <= 19i32": {
			source: "6u64 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19i32`"),
			},
		},
		"6u64 <= 19i16": {
			source: "6u64 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19i16`"),
			},
		},
		"6u64 <= 19i8": {
			source: "6u64 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19i8`"),
			},
		},
		"6u64 <= 19u32": {
			source: "6u64 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19u32`"),
			},
		},
		"6u64 <= 19u16": {
			source: "6u64 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19u16`"),
			},
		},
		"6u64 <= 19u8": {
			source: "6u64 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt64` for parameter `other` in call to `Std::UInt64.:<=`, got type `19u8`"),
			},
		},

		// UInt32
		"25u32 <= 25u32": {
			source:       "25u32 <= 25u32",
			wantStackTop: value.True.ToValue(),
		},
		"13u32 <= 7u32": {
			source:       "13u32 <= 7u32",
			wantStackTop: value.False.ToValue(),
		},
		"7u32 <= 13u32": {
			source:       "7u32 <= 13u32",
			wantStackTop: value.True.ToValue(),
		},
		"6u32 <= 19": {
			source: "6u32 <= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19`"),
			},
		},
		"6u32 <= 19.0": {
			source: "6u32 <= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19.0`"),
			},
		},
		"6u32 <= 19bf": {
			source: "6u32 <= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19bf`"),
			},
		},
		"6u32 <= 19f64": {
			source: "6u32 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19f64`"),
			},
		},
		"6u32 <= 19f32": {
			source: "6u32 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19f32`"),
			},
		},
		"6u32 <= 19i64": {
			source: "6u32 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19i64`"),
			},
		},
		"6u32 <= 19i32": {
			source: "6u32 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19i32`"),
			},
		},
		"6u32 <= 19i16": {
			source: "6u32 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19i16`"),
			},
		},
		"6u32 <= 19i8": {
			source: "6u32 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19i8`"),
			},
		},
		"6u32 <= 19u64": {
			source: "6u32 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19u64`"),
			},
		},
		"6u32 <= 19u16": {
			source: "6u32 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19u16`"),
			},
		},
		"6u32 <= 19u8": {
			source: "6u32 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt32` for parameter `other` in call to `Std::UInt32.:<=`, got type `19u8`"),
			},
		},

		// UInt16
		"25u16 <= 25u16": {
			source:       "25u16 <= 25u16",
			wantStackTop: value.True.ToValue(),
		},
		"13u16 <= 7u16": {
			source:       "13u16 <= 7u16",
			wantStackTop: value.False.ToValue(),
		},
		"7u16 <= 13u16": {
			source:       "7u16 <= 13u16",
			wantStackTop: value.True.ToValue(),
		},
		"6u16 <= 19": {
			source: "6u16 <= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(9, 1, 10)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19`"),
			},
		},
		"6u16 <= 19.0": {
			source: "6u16 <= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19.0`"),
			},
		},
		"6u16 <= 19bf": {
			source: "6u16 <= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19bf`"),
			},
		},
		"6u16 <= 19f64": {
			source: "6u16 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19f64`"),
			},
		},
		"6u16 <= 19f32": {
			source: "6u16 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19f32`"),
			},
		},
		"6u16 <= 19i64": {
			source: "6u16 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19i64`"),
			},
		},
		"6u16 <= 19i32": {
			source: "6u16 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19i32`"),
			},
		},
		"6u16 <= 19i16": {
			source: "6u16 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19i16`"),
			},
		},
		"6u16 <= 19i8": {
			source: "6u16 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19i8`"),
			},
		},
		"6u16 <= 19u64": {
			source: "6u16 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19u64`"),
			},
		},
		"6u16 <= 19u32": {
			source: "6u16 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19u32`"),
			},
		},
		"6u16 <= 19u8": {
			source: "6u16 <= 19u8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::UInt16` for parameter `other` in call to `Std::UInt16.:<=`, got type `19u8`"),
			},
		},

		// UInt8
		"25u8 <= 25u8": {
			source:       "25u8 <= 25u8",
			wantStackTop: value.True.ToValue(),
		},
		"13u8 <= 7u8": {
			source:       "13u8 <= 7u8",
			wantStackTop: value.False.ToValue(),
		},
		"7u8 <= 13u8": {
			source:       "7u8 <= 13u8",
			wantStackTop: value.True.ToValue(),
		},
		"6u8 <= 19": {
			source: "6u8 <= 19",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(8, 1, 9)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19`"),
			},
		},
		"6u8 <= 19.0": {
			source: "6u8 <= 19.0",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19.0`"),
			},
		},
		"6u8 <= 19bf": {
			source: "6u8 <= 19bf",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19bf`"),
			},
		},
		"6u8 <= 19f64": {
			source: "6u8 <= 19f64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19f64`"),
			},
		},
		"6u8 <= 19f32": {
			source: "6u8 <= 19f32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19f32`"),
			},
		},
		"6u8 <= 19i64": {
			source: "6u8 <= 19i64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19i64`"),
			},
		},
		"6u8 <= 19i32": {
			source: "6u8 <= 19i32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19i32`"),
			},
		},
		"6u8 <= 19i16": {
			source: "6u8 <= 19i16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19i16`"),
			},
		},
		"6u8 <= 19i8": {
			source: "6u8 <= 19i8",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(10, 1, 11)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19i8`"),
			},
		},
		"6u8 <= 19u64": {
			source: "6u8 <= 19u64",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19u64`"),
			},
		},
		"6u8 <= 19u32": {
			source: "6u8 <= 19u32",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19u32`"),
			},
		},
		"6u8 <= 19u16": {
			source: "6u8 <= 19u16",
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::UInt8` for parameter `other` in call to `Std::UInt8.:<=`, got type `19u16`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LaxEqual(t *testing.T) {
	tests := simpleSourceTestTable{
		// String
		"'25' =~ '25'":   value.True.ToValue(),
		"'25' =~ '25.0'": value.False.ToValue(),
		"'25' =~ '7'":    value.False.ToValue(),

		"'7' =~ `7`":  value.True.ToValue(),
		"'a' =~ `a`":  value.True.ToValue(),
		"'7' =~ `5`":  value.False.ToValue(),
		"'ab' =~ `a`": value.False.ToValue(),

		"'25' =~ 25.0":   value.False.ToValue(),
		"'13.3' =~ 13.3": value.False.ToValue(),

		"'25' =~ 25bf":     value.False.ToValue(),
		"'13.3' =~ 13.3bf": value.False.ToValue(),

		"'25' =~ 25f64": value.False.ToValue(),

		"'25' =~ 25f32": value.False.ToValue(),

		"'1' =~ 1i64": value.False.ToValue(),

		"'5' =~ 5i32": value.False.ToValue(),

		"'5' =~ 5i16": value.False.ToValue(),

		"'5' =~ 5i8": value.False.ToValue(),

		"'1' =~ 1u64": value.False.ToValue(),

		"'5' =~ 5u32": value.False.ToValue(),

		"'5' =~ 5u16": value.False.ToValue(),

		"'5' =~ 5u8": value.False.ToValue(),

		// Char
		"`2` =~ '2'":   value.True.ToValue(),
		"`a` =~ 'a'":   value.True.ToValue(),
		"`a` =~ 'ab'":  value.False.ToValue(),
		"`2` =~ '2.0'": value.False.ToValue(),

		"`7` =~ `7`": value.True.ToValue(),
		"`a` =~ `a`": value.True.ToValue(),
		"`7` =~ `5`": value.False.ToValue(),
		"`a` =~ `b`": value.False.ToValue(),

		"`2` =~ 2.0": value.False.ToValue(),

		"`9` =~ 9bf": value.False.ToValue(),

		"`3` =~ 3f64": value.False.ToValue(),

		"`7` =~ 7f32": value.False.ToValue(),

		"`1` =~ 1i64": value.False.ToValue(),

		"`5` =~ 5i32": value.False.ToValue(),

		"`5` =~ 5i16": value.False.ToValue(),

		"`5` =~ 5i8": value.False.ToValue(),

		"`1` =~ 1u64": value.False.ToValue(),

		"`5` =~ 5u32": value.False.ToValue(),

		"`5` =~ 5u16": value.False.ToValue(),

		"`5` =~ 5u8": value.False.ToValue(),

		// Int
		"25 =~ 25":  value.True.ToValue(),
		"-25 =~ 25": value.False.ToValue(),
		"25 =~ -25": value.False.ToValue(),
		"25 =~ 28":  value.False.ToValue(),
		"28 =~ 25":  value.False.ToValue(),

		"25 =~ '25'": value.False.ToValue(),

		"7 =~ `7`": value.False.ToValue(),

		"-73 =~ 73.0": value.False.ToValue(),
		"73 =~ -73.0": value.False.ToValue(),
		"25 =~ 25.0":  value.True.ToValue(),
		"1 =~ 1.2":    value.False.ToValue(),

		"-73 =~ 73bf": value.False.ToValue(),
		"73 =~ -73bf": value.False.ToValue(),
		"25 =~ 25bf":  value.True.ToValue(),
		"1 =~ 1.2bf":  value.False.ToValue(),

		"-73 =~ 73f64": value.False.ToValue(),
		"73 =~ -73f64": value.False.ToValue(),
		"25 =~ 25f64":  value.True.ToValue(),
		"1 =~ 1.2f64":  value.False.ToValue(),

		"-73 =~ 73f32": value.False.ToValue(),
		"73 =~ -73f32": value.False.ToValue(),
		"25 =~ 25f32":  value.True.ToValue(),
		"1 =~ 1.2f32":  value.False.ToValue(),

		"1 =~ 1i64":   value.True.ToValue(),
		"4 =~ -4i64":  value.False.ToValue(),
		"-8 =~ 8i64":  value.False.ToValue(),
		"-8 =~ -8i64": value.True.ToValue(),
		"91 =~ 27i64": value.False.ToValue(),

		"5 =~ 5i32":  value.True.ToValue(),
		"4 =~ -4i32": value.False.ToValue(),
		"-8 =~ 8i32": value.False.ToValue(),
		"3 =~ 71i32": value.False.ToValue(),

		"5 =~ 5i16":  value.True.ToValue(),
		"4 =~ -4i16": value.False.ToValue(),
		"-8 =~ 8i16": value.False.ToValue(),
		"3 =~ 71i16": value.False.ToValue(),

		"5 =~ 5i8":  value.True.ToValue(),
		"4 =~ -4i8": value.False.ToValue(),
		"-8 =~ 8i8": value.False.ToValue(),
		"3 =~ 71i8": value.False.ToValue(),

		"1 =~ 1u64":   value.True.ToValue(),
		"-8 =~ 8u64":  value.False.ToValue(),
		"91 =~ 27u64": value.False.ToValue(),

		"5 =~ 5u32":  value.True.ToValue(),
		"-8 =~ 8u32": value.False.ToValue(),
		"3 =~ 71u32": value.False.ToValue(),

		"53000 =~ 32767u16": value.False.ToValue(),
		"5 =~ 5u16":         value.True.ToValue(),
		"-8 =~ 8u16":        value.False.ToValue(),
		"3 =~ 71u16":        value.False.ToValue(),

		"256 =~ 127u8": value.False.ToValue(),
		"5 =~ 5u8":     value.True.ToValue(),
		"-8 =~ 8u8":    value.False.ToValue(),
		"3 =~ 71u8":    value.False.ToValue(),

		// Int64
		"25i64 =~ 25":  value.True.ToValue(),
		"-25i64 =~ 25": value.False.ToValue(),
		"25i64 =~ -25": value.False.ToValue(),
		"25i64 =~ 28":  value.False.ToValue(),
		"28i64 =~ 25":  value.False.ToValue(),

		"25i64 =~ '25'": value.False.ToValue(),

		"7i64 =~ `7`": value.False.ToValue(),

		"-73i64 =~ 73.0": value.False.ToValue(),
		"73i64 =~ -73.0": value.False.ToValue(),
		"25i64 =~ 25.0":  value.True.ToValue(),
		"1i64 =~ 1.2":    value.False.ToValue(),

		"-73i64 =~ 73bf": value.False.ToValue(),
		"73i64 =~ -73bf": value.False.ToValue(),
		"25i64 =~ 25bf":  value.True.ToValue(),
		"1i64 =~ 1.2bf":  value.False.ToValue(),

		"-73i64 =~ 73f64": value.False.ToValue(),
		"73i64 =~ -73f64": value.False.ToValue(),
		"25i64 =~ 25f64":  value.True.ToValue(),
		"1i64 =~ 1.2f64":  value.False.ToValue(),

		"-73i64 =~ 73f32": value.False.ToValue(),
		"73i64 =~ -73f32": value.False.ToValue(),
		"25i64 =~ 25f32":  value.True.ToValue(),
		"1i64 =~ 1.2f32":  value.False.ToValue(),

		"1i64 =~ 1i64":   value.True.ToValue(),
		"4i64 =~ -4i64":  value.False.ToValue(),
		"-8i64 =~ 8i64":  value.False.ToValue(),
		"-8i64 =~ -8i64": value.True.ToValue(),
		"91i64 =~ 27i64": value.False.ToValue(),

		"5i64 =~ 5i32":  value.True.ToValue(),
		"4i64 =~ -4i32": value.False.ToValue(),
		"-8i64 =~ 8i32": value.False.ToValue(),
		"3i64 =~ 71i32": value.False.ToValue(),

		"5i64 =~ 5i16":  value.True.ToValue(),
		"4i64 =~ -4i16": value.False.ToValue(),
		"-8i64 =~ 8i16": value.False.ToValue(),
		"3i64 =~ 71i16": value.False.ToValue(),

		"5i64 =~ 5i8":  value.True.ToValue(),
		"4i64 =~ -4i8": value.False.ToValue(),
		"-8i64 =~ 8i8": value.False.ToValue(),
		"3i64 =~ 71i8": value.False.ToValue(),

		"1i64 =~ 1u64":   value.True.ToValue(),
		"-8i64 =~ 8u64":  value.False.ToValue(),
		"91i64 =~ 27u64": value.False.ToValue(),

		"5i64 =~ 5u32":  value.True.ToValue(),
		"-8i64 =~ 8u32": value.False.ToValue(),
		"3i64 =~ 71u32": value.False.ToValue(),

		"53000i64 =~ 32767u16": value.False.ToValue(),
		"5i64 =~ 5u16":         value.True.ToValue(),
		"-8i64 =~ 8u16":        value.False.ToValue(),
		"3i64 =~ 71u16":        value.False.ToValue(),

		"256i64 =~ 127u8": value.False.ToValue(),
		"5i64 =~ 5u8":     value.True.ToValue(),
		"-8i64 =~ 8u8":    value.False.ToValue(),
		"3i64 =~ 71u8":    value.False.ToValue(),

		// Int32
		"25i32 =~ 25":  value.True.ToValue(),
		"-25i32 =~ 25": value.False.ToValue(),
		"25i32 =~ -25": value.False.ToValue(),
		"25i32 =~ 28":  value.False.ToValue(),
		"28i32 =~ 25":  value.False.ToValue(),

		"25i32 =~ '25'": value.False.ToValue(),

		"7i32 =~ `7`": value.False.ToValue(),

		"-73i32 =~ 73.0": value.False.ToValue(),
		"73i32 =~ -73.0": value.False.ToValue(),
		"25i32 =~ 25.0":  value.True.ToValue(),
		"1i32 =~ 1.2":    value.False.ToValue(),

		"-73i32 =~ 73bf": value.False.ToValue(),
		"73i32 =~ -73bf": value.False.ToValue(),
		"25i32 =~ 25bf":  value.True.ToValue(),
		"1i32 =~ 1.2bf":  value.False.ToValue(),

		"-73i32 =~ 73f64": value.False.ToValue(),
		"73i32 =~ -73f64": value.False.ToValue(),
		"25i32 =~ 25f64":  value.True.ToValue(),
		"1i32 =~ 1.2f64":  value.False.ToValue(),

		"-73i32 =~ 73f32": value.False.ToValue(),
		"73i32 =~ -73f32": value.False.ToValue(),
		"25i32 =~ 25f32":  value.True.ToValue(),
		"1i32 =~ 1.2f32":  value.False.ToValue(),

		"1i32 =~ 1i64":   value.True.ToValue(),
		"4i32 =~ -4i64":  value.False.ToValue(),
		"-8i32 =~ 8i64":  value.False.ToValue(),
		"-8i32 =~ -8i64": value.True.ToValue(),
		"91i32 =~ 27i64": value.False.ToValue(),

		"5i32 =~ 5i32":  value.True.ToValue(),
		"4i32 =~ -4i32": value.False.ToValue(),
		"-8i32 =~ 8i32": value.False.ToValue(),
		"3i32 =~ 71i32": value.False.ToValue(),

		"5i32 =~ 5i16":  value.True.ToValue(),
		"4i32 =~ -4i16": value.False.ToValue(),
		"-8i32 =~ 8i16": value.False.ToValue(),
		"3i32 =~ 71i16": value.False.ToValue(),

		"5i32 =~ 5i8":  value.True.ToValue(),
		"4i32 =~ -4i8": value.False.ToValue(),
		"-8i32 =~ 8i8": value.False.ToValue(),
		"3i32 =~ 71i8": value.False.ToValue(),

		"1i32 =~ 1u64":   value.True.ToValue(),
		"-8i32 =~ 8u64":  value.False.ToValue(),
		"91i32 =~ 27u64": value.False.ToValue(),

		"5i32 =~ 5u32":  value.True.ToValue(),
		"-8i32 =~ 8u32": value.False.ToValue(),
		"3i32 =~ 71u32": value.False.ToValue(),

		"53000i32 =~ 32767u16": value.False.ToValue(),
		"5i32 =~ 5u16":         value.True.ToValue(),
		"-8i32 =~ 8u16":        value.False.ToValue(),
		"3i32 =~ 71u16":        value.False.ToValue(),

		"256i32 =~ 127u8": value.False.ToValue(),
		"5i32 =~ 5u8":     value.True.ToValue(),
		"-8i32 =~ 8u8":    value.False.ToValue(),
		"3i32 =~ 71u8":    value.False.ToValue(),

		// Int16
		"25i16 =~ 25":  value.True.ToValue(),
		"-25i16 =~ 25": value.False.ToValue(),
		"25i16 =~ -25": value.False.ToValue(),
		"25i16 =~ 28":  value.False.ToValue(),
		"28i16 =~ 25":  value.False.ToValue(),

		"25i16 =~ '25'": value.False.ToValue(),

		"7i16 =~ `7`": value.False.ToValue(),

		"-73i16 =~ 73.0": value.False.ToValue(),
		"73i16 =~ -73.0": value.False.ToValue(),
		"25i16 =~ 25.0":  value.True.ToValue(),
		"1i16 =~ 1.2":    value.False.ToValue(),

		"-73i16 =~ 73bf": value.False.ToValue(),
		"73i16 =~ -73bf": value.False.ToValue(),
		"25i16 =~ 25bf":  value.True.ToValue(),
		"1i16 =~ 1.2bf":  value.False.ToValue(),

		"-73i16 =~ 73f64": value.False.ToValue(),
		"73i16 =~ -73f64": value.False.ToValue(),
		"25i16 =~ 25f64":  value.True.ToValue(),
		"1i16 =~ 1.2f64":  value.False.ToValue(),

		"-73i16 =~ 73f32": value.False.ToValue(),
		"73i16 =~ -73f32": value.False.ToValue(),
		"25i16 =~ 25f32":  value.True.ToValue(),
		"1i16 =~ 1.2f32":  value.False.ToValue(),

		"1i16 =~ 1i64":   value.True.ToValue(),
		"4i16 =~ -4i64":  value.False.ToValue(),
		"-8i16 =~ 8i64":  value.False.ToValue(),
		"-8i16 =~ -8i64": value.True.ToValue(),
		"91i16 =~ 27i64": value.False.ToValue(),

		"5i16 =~ 5i32":  value.True.ToValue(),
		"4i16 =~ -4i32": value.False.ToValue(),
		"-8i16 =~ 8i32": value.False.ToValue(),
		"3i16 =~ 71i32": value.False.ToValue(),

		"5i16 =~ 5i16":  value.True.ToValue(),
		"4i16 =~ -4i16": value.False.ToValue(),
		"-8i16 =~ 8i16": value.False.ToValue(),
		"3i16 =~ 71i16": value.False.ToValue(),

		"5i16 =~ 5i8":  value.True.ToValue(),
		"4i16 =~ -4i8": value.False.ToValue(),
		"-8i16 =~ 8i8": value.False.ToValue(),
		"3i16 =~ 71i8": value.False.ToValue(),

		"1i16 =~ 1u64":   value.True.ToValue(),
		"-8i16 =~ 8u64":  value.False.ToValue(),
		"91i16 =~ 27u64": value.False.ToValue(),

		"5i16 =~ 5u32":  value.True.ToValue(),
		"-8i16 =~ 8u32": value.False.ToValue(),
		"3i16 =~ 71u32": value.False.ToValue(),

		"5i16 =~ 5u16":  value.True.ToValue(),
		"-8i16 =~ 8u16": value.False.ToValue(),
		"3i16 =~ 71u16": value.False.ToValue(),

		"256i16 =~ 127u8": value.False.ToValue(),
		"5i16 =~ 5u8":     value.True.ToValue(),
		"-8i16 =~ 8u8":    value.False.ToValue(),
		"3i16 =~ 71u8":    value.False.ToValue(),

		// Int8
		"25i8 =~ 25":  value.True.ToValue(),
		"-25i8 =~ 25": value.False.ToValue(),
		"25i8 =~ -25": value.False.ToValue(),
		"25i8 =~ 28":  value.False.ToValue(),
		"28i8 =~ 25":  value.False.ToValue(),

		"25i8 =~ '25'": value.False.ToValue(),

		"7i8 =~ `7`": value.False.ToValue(),

		"-73i8 =~ 73.0": value.False.ToValue(),
		"73i8 =~ -73.0": value.False.ToValue(),
		"25i8 =~ 25.0":  value.True.ToValue(),
		"1i8 =~ 1.2":    value.False.ToValue(),

		"-73i8 =~ 73bf": value.False.ToValue(),
		"73i8 =~ -73bf": value.False.ToValue(),
		"25i8 =~ 25bf":  value.True.ToValue(),
		"1i8 =~ 1.2bf":  value.False.ToValue(),

		"-73i8 =~ 73f64": value.False.ToValue(),
		"73i8 =~ -73f64": value.False.ToValue(),
		"25i8 =~ 25f64":  value.True.ToValue(),
		"1i8 =~ 1.2f64":  value.False.ToValue(),

		"-73i8 =~ 73f32": value.False.ToValue(),
		"73i8 =~ -73f32": value.False.ToValue(),
		"25i8 =~ 25f32":  value.True.ToValue(),
		"1i8 =~ 1.2f32":  value.False.ToValue(),

		"1i8 =~ 1i64":   value.True.ToValue(),
		"4i8 =~ -4i64":  value.False.ToValue(),
		"-8i8 =~ 8i64":  value.False.ToValue(),
		"-8i8 =~ -8i64": value.True.ToValue(),
		"91i8 =~ 27i64": value.False.ToValue(),

		"5i8 =~ 5i32":  value.True.ToValue(),
		"4i8 =~ -4i32": value.False.ToValue(),
		"-8i8 =~ 8i32": value.False.ToValue(),
		"3i8 =~ 71i32": value.False.ToValue(),

		"5i8 =~ 5i16":  value.True.ToValue(),
		"4i8 =~ -4i16": value.False.ToValue(),
		"-8i8 =~ 8i16": value.False.ToValue(),
		"3i8 =~ 71i16": value.False.ToValue(),

		"5i8 =~ 5i8":  value.True.ToValue(),
		"4i8 =~ -4i8": value.False.ToValue(),
		"-8i8 =~ 8i8": value.False.ToValue(),
		"3i8 =~ 71i8": value.False.ToValue(),

		"1i8 =~ 1u64":   value.True.ToValue(),
		"-8i8 =~ 8u64":  value.False.ToValue(),
		"91i8 =~ 27u64": value.False.ToValue(),

		"5i8 =~ 5u32":  value.True.ToValue(),
		"-8i8 =~ 8u32": value.False.ToValue(),
		"3i8 =~ 71u32": value.False.ToValue(),

		"5i8 =~ 5u16":  value.True.ToValue(),
		"-8i8 =~ 8u16": value.False.ToValue(),
		"3i8 =~ 71u16": value.False.ToValue(),

		"5i8 =~ 5u8":  value.True.ToValue(),
		"-8i8 =~ 8u8": value.False.ToValue(),
		"3i8 =~ 71u8": value.False.ToValue(),

		// UInt64
		"25u64 =~ 25":  value.True.ToValue(),
		"25u64 =~ -25": value.False.ToValue(),
		"25u64 =~ 28":  value.False.ToValue(),
		"28u64 =~ 25":  value.False.ToValue(),

		"25u64 =~ '25'": value.False.ToValue(),

		"7u64 =~ `7`": value.False.ToValue(),

		"73u64 =~ -73.0": value.False.ToValue(),
		"25u64 =~ 25.0":  value.True.ToValue(),
		"1u64 =~ 1.2":    value.False.ToValue(),

		"73u64 =~ -73bf": value.False.ToValue(),
		"25u64 =~ 25bf":  value.True.ToValue(),
		"1u64 =~ 1.2bf":  value.False.ToValue(),

		"73u64 =~ -73f64": value.False.ToValue(),
		"25u64 =~ 25f64":  value.True.ToValue(),
		"1u64 =~ 1.2f64":  value.False.ToValue(),

		"73u64 =~ -73f32": value.False.ToValue(),
		"25u64 =~ 25f32":  value.True.ToValue(),
		"1u64 =~ 1.2f32":  value.False.ToValue(),

		"1u64 =~ 1i64":   value.True.ToValue(),
		"4u64 =~ -4i64":  value.False.ToValue(),
		"91u64 =~ 27i64": value.False.ToValue(),

		"5u64 =~ 5i32":  value.True.ToValue(),
		"4u64 =~ -4i32": value.False.ToValue(),
		"3u64 =~ 71i32": value.False.ToValue(),

		"5u64 =~ 5i16":  value.True.ToValue(),
		"4u64 =~ -4i16": value.False.ToValue(),
		"3u64 =~ 71i16": value.False.ToValue(),

		"5u64 =~ 5i8":  value.True.ToValue(),
		"4u64 =~ -4i8": value.False.ToValue(),
		"3u64 =~ 71i8": value.False.ToValue(),

		"1u64 =~ 1u64":   value.True.ToValue(),
		"91u64 =~ 27u64": value.False.ToValue(),

		"5u64 =~ 5u32":  value.True.ToValue(),
		"3u64 =~ 71u32": value.False.ToValue(),

		"53000u64 =~ 32767u16": value.False.ToValue(),
		"5u64 =~ 5u16":         value.True.ToValue(),
		"3u64 =~ 71u16":        value.False.ToValue(),

		"256u64 =~ 127u8": value.False.ToValue(),
		"5u64 =~ 5u8":     value.True.ToValue(),
		"3u64 =~ 71u8":    value.False.ToValue(),

		// UInt32
		"25u32 =~ 25":  value.True.ToValue(),
		"25u32 =~ -25": value.False.ToValue(),
		"25u32 =~ 28":  value.False.ToValue(),
		"28u32 =~ 25":  value.False.ToValue(),

		"25u32 =~ '25'": value.False.ToValue(),

		"7u32 =~ `7`": value.False.ToValue(),

		"73u32 =~ -73.0": value.False.ToValue(),
		"25u32 =~ 25.0":  value.True.ToValue(),
		"1u32 =~ 1.2":    value.False.ToValue(),

		"73u32 =~ -73bf": value.False.ToValue(),
		"25u32 =~ 25bf":  value.True.ToValue(),
		"1u32 =~ 1.2bf":  value.False.ToValue(),

		"73u32 =~ -73f64": value.False.ToValue(),
		"25u32 =~ 25f64":  value.True.ToValue(),
		"1u32 =~ 1.2f64":  value.False.ToValue(),

		"73u32 =~ -73f32": value.False.ToValue(),
		"25u32 =~ 25f32":  value.True.ToValue(),
		"1u32 =~ 1.2f32":  value.False.ToValue(),

		"1u32 =~ 1i64":   value.True.ToValue(),
		"4u32 =~ -4i64":  value.False.ToValue(),
		"91u32 =~ 27i64": value.False.ToValue(),

		"5u32 =~ 5i32":  value.True.ToValue(),
		"4u32 =~ -4i32": value.False.ToValue(),
		"3u32 =~ 71i32": value.False.ToValue(),

		"5u32 =~ 5i16":  value.True.ToValue(),
		"4u32 =~ -4i16": value.False.ToValue(),
		"3u32 =~ 71i16": value.False.ToValue(),

		"5u32 =~ 5i8":  value.True.ToValue(),
		"4u32 =~ -4i8": value.False.ToValue(),
		"3u32 =~ 71i8": value.False.ToValue(),

		"1u32 =~ 1u64":   value.True.ToValue(),
		"91u32 =~ 27u64": value.False.ToValue(),

		"5u32 =~ 5u32":  value.True.ToValue(),
		"3u32 =~ 71u32": value.False.ToValue(),

		"53000u32 =~ 32767u16": value.False.ToValue(),
		"5u32 =~ 5u16":         value.True.ToValue(),
		"3u32 =~ 71u16":        value.False.ToValue(),

		"256u32 =~ 127u8": value.False.ToValue(),
		"5u32 =~ 5u8":     value.True.ToValue(),
		"3u32 =~ 71u8":    value.False.ToValue(),

		// UInt16
		"25u16 =~ 25":  value.True.ToValue(),
		"25u16 =~ -25": value.False.ToValue(),
		"25u16 =~ 28":  value.False.ToValue(),
		"28u16 =~ 25":  value.False.ToValue(),

		"25u16 =~ '25'": value.False.ToValue(),

		"7u16 =~ `7`": value.False.ToValue(),

		"73u16 =~ -73.0": value.False.ToValue(),
		"25u16 =~ 25.0":  value.True.ToValue(),
		"1u16 =~ 1.2":    value.False.ToValue(),

		"73u16 =~ -73bf": value.False.ToValue(),
		"25u16 =~ 25bf":  value.True.ToValue(),
		"1u16 =~ 1.2bf":  value.False.ToValue(),

		"73u16 =~ -73f64": value.False.ToValue(),
		"25u16 =~ 25f64":  value.True.ToValue(),
		"1u16 =~ 1.2f64":  value.False.ToValue(),

		"73u16 =~ -73f32": value.False.ToValue(),
		"25u16 =~ 25f32":  value.True.ToValue(),
		"1u16 =~ 1.2f32":  value.False.ToValue(),

		"1u16 =~ 1i64":   value.True.ToValue(),
		"4u16 =~ -4i64":  value.False.ToValue(),
		"91u16 =~ 27i64": value.False.ToValue(),

		"5u16 =~ 5i32":  value.True.ToValue(),
		"4u16 =~ -4i32": value.False.ToValue(),
		"3u16 =~ 71i32": value.False.ToValue(),

		"5u16 =~ 5i16":  value.True.ToValue(),
		"4u16 =~ -4i16": value.False.ToValue(),
		"3u16 =~ 71i16": value.False.ToValue(),

		"5u16 =~ 5i8":  value.True.ToValue(),
		"4u16 =~ -4i8": value.False.ToValue(),
		"3u16 =~ 71i8": value.False.ToValue(),

		"1u16 =~ 1u64":   value.True.ToValue(),
		"91u16 =~ 27u64": value.False.ToValue(),

		"5u16 =~ 5u32":  value.True.ToValue(),
		"3u16 =~ 71u32": value.False.ToValue(),

		"53000u16 =~ 32767u16": value.False.ToValue(),
		"5u16 =~ 5u16":         value.True.ToValue(),
		"3u16 =~ 71u16":        value.False.ToValue(),

		"256u16 =~ 127u8": value.False.ToValue(),
		"5u16 =~ 5u8":     value.True.ToValue(),
		"3u16 =~ 71u8":    value.False.ToValue(),

		// UInt8
		"25u8 =~ 25":  value.True.ToValue(),
		"25u8 =~ -25": value.False.ToValue(),
		"25u8 =~ 28":  value.False.ToValue(),
		"28u8 =~ 25":  value.False.ToValue(),

		"25u8 =~ '25'": value.False.ToValue(),

		"7u8 =~ `7`": value.False.ToValue(),

		"73u8 =~ -73.0": value.False.ToValue(),
		"25u8 =~ 25.0":  value.True.ToValue(),
		"1u8 =~ 1.2":    value.False.ToValue(),

		"73u8 =~ -73bf": value.False.ToValue(),
		"25u8 =~ 25bf":  value.True.ToValue(),
		"1u8 =~ 1.2bf":  value.False.ToValue(),

		"73u8 =~ -73f64": value.False.ToValue(),
		"25u8 =~ 25f64":  value.True.ToValue(),
		"1u8 =~ 1.2f64":  value.False.ToValue(),

		"73u8 =~ -73f32": value.False.ToValue(),
		"25u8 =~ 25f32":  value.True.ToValue(),
		"1u8 =~ 1.2f32":  value.False.ToValue(),

		"1u8 =~ 1i64":   value.True.ToValue(),
		"4u8 =~ -4i64":  value.False.ToValue(),
		"91u8 =~ 27i64": value.False.ToValue(),

		"5u8 =~ 5i32":  value.True.ToValue(),
		"4u8 =~ -4i32": value.False.ToValue(),
		"3u8 =~ 71i32": value.False.ToValue(),

		"5u8 =~ 5i16":  value.True.ToValue(),
		"4u8 =~ -4i16": value.False.ToValue(),
		"3u8 =~ 71i16": value.False.ToValue(),

		"5u8 =~ 5i8":  value.True.ToValue(),
		"4u8 =~ -4i8": value.False.ToValue(),
		"3u8 =~ 71i8": value.False.ToValue(),

		"1u8 =~ 1u64":   value.True.ToValue(),
		"91u8 =~ 27u64": value.False.ToValue(),

		"5u8 =~ 5u32":  value.True.ToValue(),
		"3u8 =~ 71u32": value.False.ToValue(),

		"5u8 =~ 5u16":  value.True.ToValue(),
		"3u8 =~ 71u16": value.False.ToValue(),

		"5u8 =~ 5u8":  value.True.ToValue(),
		"3u8 =~ 71u8": value.False.ToValue(),

		// Float
		"-73.0 =~ 73.0": value.False.ToValue(),
		"73.0 =~ -73.0": value.False.ToValue(),
		"25.0 =~ 25.0":  value.True.ToValue(),
		"1.0 =~ 1.2":    value.False.ToValue(),
		"1.2 =~ 1.0":    value.False.ToValue(),
		"78.5 =~ 78.5":  value.True.ToValue(),

		"8.25 =~ '8.25'": value.False.ToValue(),

		"4.0 =~ `4`": value.False.ToValue(),

		"25.0 =~ 25":  value.True.ToValue(),
		"32.3 =~ 32":  value.False.ToValue(),
		"-25.0 =~ 25": value.False.ToValue(),
		"25.0 =~ -25": value.False.ToValue(),
		"25.0 =~ 28":  value.False.ToValue(),
		"28.0 =~ 25":  value.False.ToValue(),

		"-73.0 =~ 73bf":  value.False.ToValue(),
		"73.0 =~ -73bf":  value.False.ToValue(),
		"25.0 =~ 25bf":   value.True.ToValue(),
		"1.0 =~ 1.2bf":   value.False.ToValue(),
		"15.5 =~ 15.5bf": value.True.ToValue(),

		"-73.0 =~ 73f64":    value.False.ToValue(),
		"73.0 =~ -73f64":    value.False.ToValue(),
		"25.0 =~ 25f64":     value.True.ToValue(),
		"1.0 =~ 1.2f64":     value.False.ToValue(),
		"15.26 =~ 15.26f64": value.True.ToValue(),

		"-73.0 =~ 73f32":  value.False.ToValue(),
		"73.0 =~ -73f32":  value.False.ToValue(),
		"25.0 =~ 25f32":   value.True.ToValue(),
		"1.0 =~ 1.2f32":   value.False.ToValue(),
		"15.5 =~ 15.5f32": value.True.ToValue(),

		"1.0 =~ 1i64":   value.True.ToValue(),
		"1.5 =~ 1i64":   value.False.ToValue(),
		"4.0 =~ -4i64":  value.False.ToValue(),
		"-8.0 =~ 8i64":  value.False.ToValue(),
		"-8.0 =~ -8i64": value.True.ToValue(),
		"91.0 =~ 27i64": value.False.ToValue(),

		"1.0 =~ 1i32":   value.True.ToValue(),
		"1.5 =~ 1i32":   value.False.ToValue(),
		"4.0 =~ -4i32":  value.False.ToValue(),
		"-8.0 =~ 8i32":  value.False.ToValue(),
		"-8.0 =~ -8i32": value.True.ToValue(),
		"91.0 =~ 27i32": value.False.ToValue(),

		"1.0 =~ 1i16":   value.True.ToValue(),
		"1.5 =~ 1i16":   value.False.ToValue(),
		"4.0 =~ -4i16":  value.False.ToValue(),
		"-8.0 =~ 8i16":  value.False.ToValue(),
		"-8.0 =~ -8i16": value.True.ToValue(),
		"91.0 =~ 27i16": value.False.ToValue(),

		"1.0 =~ 1i8":   value.True.ToValue(),
		"1.5 =~ 1i8":   value.False.ToValue(),
		"4.0 =~ -4i8":  value.False.ToValue(),
		"-8.0 =~ 8i8":  value.False.ToValue(),
		"-8.0 =~ -8i8": value.True.ToValue(),
		"91.0 =~ 27i8": value.False.ToValue(),

		"1.0 =~ 1u64":   value.True.ToValue(),
		"1.5 =~ 1u64":   value.False.ToValue(),
		"-8.0 =~ 8u64":  value.False.ToValue(),
		"91.0 =~ 27u64": value.False.ToValue(),

		"1.0 =~ 1u32":   value.True.ToValue(),
		"1.5 =~ 1u32":   value.False.ToValue(),
		"-8.0 =~ 8u32":  value.False.ToValue(),
		"91.0 =~ 27u32": value.False.ToValue(),

		"53000.0 =~ 32767u16": value.False.ToValue(),
		"1.0 =~ 1u16":         value.True.ToValue(),
		"1.5 =~ 1u16":         value.False.ToValue(),
		"-8.0 =~ 8u16":        value.False.ToValue(),
		"91.0 =~ 27u16":       value.False.ToValue(),

		"256.0 =~ 127u8": value.False.ToValue(),
		"1.0 =~ 1u8":     value.True.ToValue(),
		"1.5 =~ 1u8":     value.False.ToValue(),
		"-8.0 =~ 8u8":    value.False.ToValue(),
		"91.0 =~ 27u8":   value.False.ToValue(),

		// Float64
		"-73f64 =~ 73.0":  value.False.ToValue(),
		"73f64 =~ -73.0":  value.False.ToValue(),
		"25f64 =~ 25.0":   value.True.ToValue(),
		"1f64 =~ 1.2":     value.False.ToValue(),
		"1.2f64 =~ 1.0":   value.False.ToValue(),
		"78.5f64 =~ 78.5": value.True.ToValue(),

		"8.25f64 =~ '8.25'": value.False.ToValue(),

		"4f64 =~ `4`": value.False.ToValue(),

		"25f64 =~ 25":   value.True.ToValue(),
		"32.3f64 =~ 32": value.False.ToValue(),
		"-25f64 =~ 25":  value.False.ToValue(),
		"25f64 =~ -25":  value.False.ToValue(),
		"25f64 =~ 28":   value.False.ToValue(),
		"28f64 =~ 25":   value.False.ToValue(),

		"-73f64 =~ 73bf":    value.False.ToValue(),
		"73f64 =~ -73bf":    value.False.ToValue(),
		"25f64 =~ 25bf":     value.True.ToValue(),
		"1f64 =~ 1.2bf":     value.False.ToValue(),
		"15.5f64 =~ 15.5bf": value.True.ToValue(),

		"-73f64 =~ 73f64":      value.False.ToValue(),
		"73f64 =~ -73f64":      value.False.ToValue(),
		"25f64 =~ 25f64":       value.True.ToValue(),
		"1f64 =~ 1.2f64":       value.False.ToValue(),
		"15.26f64 =~ 15.26f64": value.True.ToValue(),

		"-73f64 =~ 73f32":    value.False.ToValue(),
		"73f64 =~ -73f32":    value.False.ToValue(),
		"25f64 =~ 25f32":     value.True.ToValue(),
		"1f64 =~ 1.2f32":     value.False.ToValue(),
		"15.5f64 =~ 15.5f32": value.True.ToValue(),

		"1f64 =~ 1i64":   value.True.ToValue(),
		"1.5f64 =~ 1i64": value.False.ToValue(),
		"4f64 =~ -4i64":  value.False.ToValue(),
		"-8f64 =~ 8i64":  value.False.ToValue(),
		"-8f64 =~ -8i64": value.True.ToValue(),
		"91f64 =~ 27i64": value.False.ToValue(),

		"1f64 =~ 1i32":   value.True.ToValue(),
		"1.5f64 =~ 1i32": value.False.ToValue(),
		"4f64 =~ -4i32":  value.False.ToValue(),
		"-8f64 =~ 8i32":  value.False.ToValue(),
		"-8f64 =~ -8i32": value.True.ToValue(),
		"91f64 =~ 27i32": value.False.ToValue(),

		"1f64 =~ 1i16":   value.True.ToValue(),
		"1.5f64 =~ 1i16": value.False.ToValue(),
		"4f64 =~ -4i16":  value.False.ToValue(),
		"-8f64 =~ 8i16":  value.False.ToValue(),
		"-8f64 =~ -8i16": value.True.ToValue(),
		"91f64 =~ 27i16": value.False.ToValue(),

		"1f64 =~ 1i8":   value.True.ToValue(),
		"1.5f64 =~ 1i8": value.False.ToValue(),
		"4f64 =~ -4i8":  value.False.ToValue(),
		"-8f64 =~ 8i8":  value.False.ToValue(),
		"-8f64 =~ -8i8": value.True.ToValue(),
		"91f64 =~ 27i8": value.False.ToValue(),

		"1f64 =~ 1u64":   value.True.ToValue(),
		"1.5f64 =~ 1u64": value.False.ToValue(),
		"-8f64 =~ 8u64":  value.False.ToValue(),
		"91f64 =~ 27u64": value.False.ToValue(),

		"1f64 =~ 1u32":   value.True.ToValue(),
		"1.5f64 =~ 1u32": value.False.ToValue(),
		"-8f64 =~ 8u32":  value.False.ToValue(),
		"91f64 =~ 27u32": value.False.ToValue(),

		"53000f64 =~ 32767u16": value.False.ToValue(),
		"1f64 =~ 1u16":         value.True.ToValue(),
		"1.5f64 =~ 1u16":       value.False.ToValue(),
		"-8f64 =~ 8u16":        value.False.ToValue(),
		"91f64 =~ 27u16":       value.False.ToValue(),

		"256f64 =~ 127u8": value.False.ToValue(),
		"1f64 =~ 1u8":     value.True.ToValue(),
		"1.5f64 =~ 1u8":   value.False.ToValue(),
		"-8f64 =~ 8u8":    value.False.ToValue(),
		"91f64 =~ 27u8":   value.False.ToValue(),

		// Float32
		"-73f32 =~ 73.0":  value.False.ToValue(),
		"73f32 =~ -73.0":  value.False.ToValue(),
		"25f32 =~ 25.0":   value.True.ToValue(),
		"1f32 =~ 1.2":     value.False.ToValue(),
		"1.2f32 =~ 1.0":   value.False.ToValue(),
		"78.5f32 =~ 78.5": value.True.ToValue(),

		"8.25f32 =~ '8.25'": value.False.ToValue(),

		"4f32 =~ `4`": value.False.ToValue(),

		"25f32 =~ 25":   value.True.ToValue(),
		"32.3f32 =~ 32": value.False.ToValue(),
		"-25f32 =~ 25":  value.False.ToValue(),
		"25f32 =~ -25":  value.False.ToValue(),
		"25f32 =~ 28":   value.False.ToValue(),
		"28f32 =~ 25":   value.False.ToValue(),

		"-73f32 =~ 73bf":    value.False.ToValue(),
		"73f32 =~ -73bf":    value.False.ToValue(),
		"25f32 =~ 25bf":     value.True.ToValue(),
		"1f32 =~ 1.2bf":     value.False.ToValue(),
		"15.5f32 =~ 15.5bf": value.True.ToValue(),

		"-73f32 =~ 73f64":    value.False.ToValue(),
		"73f32 =~ -73f64":    value.False.ToValue(),
		"25f32 =~ 25f64":     value.True.ToValue(),
		"1f32 =~ 1.2f64":     value.False.ToValue(),
		"15.5f32 =~ 15.5f64": value.True.ToValue(),

		"-73f32 =~ 73f32":    value.False.ToValue(),
		"73f32 =~ -73f32":    value.False.ToValue(),
		"25f32 =~ 25f32":     value.True.ToValue(),
		"1f32 =~ 1.2f32":     value.False.ToValue(),
		"15.5f32 =~ 15.5f32": value.True.ToValue(),

		"1f32 =~ 1i64":   value.True.ToValue(),
		"1.5f32 =~ 1i64": value.False.ToValue(),
		"4f32 =~ -4i64":  value.False.ToValue(),
		"-8f32 =~ 8i64":  value.False.ToValue(),
		"-8f32 =~ -8i64": value.True.ToValue(),
		"91f32 =~ 27i64": value.False.ToValue(),

		"1f32 =~ 1i32":   value.True.ToValue(),
		"1.5f32 =~ 1i32": value.False.ToValue(),
		"4f32 =~ -4i32":  value.False.ToValue(),
		"-8f32 =~ 8i32":  value.False.ToValue(),
		"-8f32 =~ -8i32": value.True.ToValue(),
		"91f32 =~ 27i32": value.False.ToValue(),

		"1f32 =~ 1i16":   value.True.ToValue(),
		"1.5f32 =~ 1i16": value.False.ToValue(),
		"4f32 =~ -4i16":  value.False.ToValue(),
		"-8f32 =~ 8i16":  value.False.ToValue(),
		"-8f32 =~ -8i16": value.True.ToValue(),
		"91f32 =~ 27i16": value.False.ToValue(),

		"1f32 =~ 1i8":   value.True.ToValue(),
		"1.5f32 =~ 1i8": value.False.ToValue(),
		"4f32 =~ -4i8":  value.False.ToValue(),
		"-8f32 =~ 8i8":  value.False.ToValue(),
		"-8f32 =~ -8i8": value.True.ToValue(),
		"91f32 =~ 27i8": value.False.ToValue(),

		"1f32 =~ 1u64":   value.True.ToValue(),
		"1.5f32 =~ 1u64": value.False.ToValue(),
		"-8f32 =~ 8u64":  value.False.ToValue(),
		"91f32 =~ 27u64": value.False.ToValue(),

		"1f32 =~ 1u32":   value.True.ToValue(),
		"1.5f32 =~ 1u32": value.False.ToValue(),
		"-8f32 =~ 8u32":  value.False.ToValue(),
		"91f32 =~ 27u32": value.False.ToValue(),

		"53000f32 =~ 32767u16": value.False.ToValue(),
		"1f32 =~ 1u16":         value.True.ToValue(),
		"1.5f32 =~ 1u16":       value.False.ToValue(),
		"-8f32 =~ 8u16":        value.False.ToValue(),
		"91f32 =~ 27u16":       value.False.ToValue(),

		"256f32 =~ 127u8": value.False.ToValue(),
		"1f32 =~ 1u8":     value.True.ToValue(),
		"1.5f32 =~ 1u8":   value.False.ToValue(),
		"-8f32 =~ 8u8":    value.False.ToValue(),
		"91f32 =~ 27u8":   value.False.ToValue(),
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}

func TestVMSource_LaxNotEqual(t *testing.T) {
	tests := simpleSourceTestTable{
		// String
		"'25' !~ '25'":   value.False.ToValue(),
		"'25' !~ '25.0'": value.True.ToValue(),
		"'25' !~ '7'":    value.True.ToValue(),

		"'7' !~ `7`":  value.False.ToValue(),
		"'a' !~ `a`":  value.False.ToValue(),
		"'7' !~ `5`":  value.True.ToValue(),
		"'ab' !~ `a`": value.True.ToValue(),

		"'25' !~ 25.0":   value.True.ToValue(),
		"'13.3' !~ 13.3": value.True.ToValue(),

		"'25' !~ 25bf":     value.True.ToValue(),
		"'13.3' !~ 13.3bf": value.True.ToValue(),

		"'25' !~ 25f64": value.True.ToValue(),

		"'25' !~ 25f32": value.True.ToValue(),

		"'1' !~ 1i64": value.True.ToValue(),

		"'5' !~ 5i32": value.True.ToValue(),

		"'5' !~ 5i16": value.True.ToValue(),

		"'5' !~ 5i8": value.True.ToValue(),

		"'1' !~ 1u64": value.True.ToValue(),

		"'5' !~ 5u32": value.True.ToValue(),

		"'5' !~ 5u16": value.True.ToValue(),

		"'5' !~ 5u8": value.True.ToValue(),

		// Char
		"`2` !~ '2'":   value.False.ToValue(),
		"`a` !~ 'a'":   value.False.ToValue(),
		"`a` !~ 'ab'":  value.True.ToValue(),
		"`2` !~ '2.0'": value.True.ToValue(),

		"`7` !~ `7`": value.False.ToValue(),
		"`a` !~ `a`": value.False.ToValue(),
		"`7` !~ `5`": value.True.ToValue(),
		"`a` !~ `b`": value.True.ToValue(),

		"`2` !~ 2.0": value.True.ToValue(),

		"`9` !~ 9bf": value.True.ToValue(),

		"`3` !~ 3f64": value.True.ToValue(),

		"`7` !~ 7f32": value.True.ToValue(),

		"`1` !~ 1i64": value.True.ToValue(),

		"`5` !~ 5i32": value.True.ToValue(),

		"`5` !~ 5i16": value.True.ToValue(),

		"`5` !~ 5i8": value.True.ToValue(),

		"`1` !~ 1u64": value.True.ToValue(),

		"`5` !~ 5u32": value.True.ToValue(),

		"`5` !~ 5u16": value.True.ToValue(),

		"`5` !~ 5u8": value.True.ToValue(),

		// Int
		"25 !~ 25":  value.False.ToValue(),
		"-25 !~ 25": value.True.ToValue(),
		"25 !~ -25": value.True.ToValue(),
		"25 !~ 28":  value.True.ToValue(),
		"28 !~ 25":  value.True.ToValue(),

		"25 !~ '25'": value.True.ToValue(),

		"7 !~ `7`": value.True.ToValue(),

		"-73 !~ 73.0": value.True.ToValue(),
		"73 !~ -73.0": value.True.ToValue(),
		"25 !~ 25.0":  value.False.ToValue(),
		"1 !~ 1.2":    value.True.ToValue(),

		"-73 !~ 73bf": value.True.ToValue(),
		"73 !~ -73bf": value.True.ToValue(),
		"25 !~ 25bf":  value.False.ToValue(),
		"1 !~ 1.2bf":  value.True.ToValue(),

		"-73 !~ 73f64": value.True.ToValue(),
		"73 !~ -73f64": value.True.ToValue(),
		"25 !~ 25f64":  value.False.ToValue(),
		"1 !~ 1.2f64":  value.True.ToValue(),

		"-73 !~ 73f32": value.True.ToValue(),
		"73 !~ -73f32": value.True.ToValue(),
		"25 !~ 25f32":  value.False.ToValue(),
		"1 !~ 1.2f32":  value.True.ToValue(),

		"1 !~ 1i64":   value.False.ToValue(),
		"4 !~ -4i64":  value.True.ToValue(),
		"-8 !~ 8i64":  value.True.ToValue(),
		"-8 !~ -8i64": value.False.ToValue(),
		"91 !~ 27i64": value.True.ToValue(),

		"5 !~ 5i32":  value.False.ToValue(),
		"4 !~ -4i32": value.True.ToValue(),
		"-8 !~ 8i32": value.True.ToValue(),
		"3 !~ 71i32": value.True.ToValue(),

		"5 !~ 5i16":  value.False.ToValue(),
		"4 !~ -4i16": value.True.ToValue(),
		"-8 !~ 8i16": value.True.ToValue(),
		"3 !~ 71i16": value.True.ToValue(),

		"5 !~ 5i8":  value.False.ToValue(),
		"4 !~ -4i8": value.True.ToValue(),
		"-8 !~ 8i8": value.True.ToValue(),
		"3 !~ 71i8": value.True.ToValue(),

		"1 !~ 1u64":   value.False.ToValue(),
		"-8 !~ 8u64":  value.True.ToValue(),
		"91 !~ 27u64": value.True.ToValue(),

		"5 !~ 5u32":  value.False.ToValue(),
		"-8 !~ 8u32": value.True.ToValue(),
		"3 !~ 71u32": value.True.ToValue(),

		"53000 !~ 32767u16": value.True.ToValue(),
		"5 !~ 5u16":         value.False.ToValue(),
		"-8 !~ 8u16":        value.True.ToValue(),
		"3 !~ 71u16":        value.True.ToValue(),

		"256 !~ 127u8": value.True.ToValue(),
		"5 !~ 5u8":     value.False.ToValue(),
		"-8 !~ 8u8":    value.True.ToValue(),
		"3 !~ 71u8":    value.True.ToValue(),

		// Int64
		"25i64 !~ 25":  value.False.ToValue(),
		"-25i64 !~ 25": value.True.ToValue(),
		"25i64 !~ -25": value.True.ToValue(),
		"25i64 !~ 28":  value.True.ToValue(),
		"28i64 !~ 25":  value.True.ToValue(),

		"25i64 !~ '25'": value.True.ToValue(),

		"7i64 !~ `7`": value.True.ToValue(),

		"-73i64 !~ 73.0": value.True.ToValue(),
		"73i64 !~ -73.0": value.True.ToValue(),
		"25i64 !~ 25.0":  value.False.ToValue(),
		"1i64 !~ 1.2":    value.True.ToValue(),

		"-73i64 !~ 73bf": value.True.ToValue(),
		"73i64 !~ -73bf": value.True.ToValue(),
		"25i64 !~ 25bf":  value.False.ToValue(),
		"1i64 !~ 1.2bf":  value.True.ToValue(),

		"-73i64 !~ 73f64": value.True.ToValue(),
		"73i64 !~ -73f64": value.True.ToValue(),
		"25i64 !~ 25f64":  value.False.ToValue(),
		"1i64 !~ 1.2f64":  value.True.ToValue(),

		"-73i64 !~ 73f32": value.True.ToValue(),
		"73i64 !~ -73f32": value.True.ToValue(),
		"25i64 !~ 25f32":  value.False.ToValue(),
		"1i64 !~ 1.2f32":  value.True.ToValue(),

		"1i64 !~ 1i64":   value.False.ToValue(),
		"4i64 !~ -4i64":  value.True.ToValue(),
		"-8i64 !~ 8i64":  value.True.ToValue(),
		"-8i64 !~ -8i64": value.False.ToValue(),
		"91i64 !~ 27i64": value.True.ToValue(),

		"5i64 !~ 5i32":  value.False.ToValue(),
		"4i64 !~ -4i32": value.True.ToValue(),
		"-8i64 !~ 8i32": value.True.ToValue(),
		"3i64 !~ 71i32": value.True.ToValue(),

		"5i64 !~ 5i16":  value.False.ToValue(),
		"4i64 !~ -4i16": value.True.ToValue(),
		"-8i64 !~ 8i16": value.True.ToValue(),
		"3i64 !~ 71i16": value.True.ToValue(),

		"5i64 !~ 5i8":  value.False.ToValue(),
		"4i64 !~ -4i8": value.True.ToValue(),
		"-8i64 !~ 8i8": value.True.ToValue(),
		"3i64 !~ 71i8": value.True.ToValue(),

		"1i64 !~ 1u64":   value.False.ToValue(),
		"-8i64 !~ 8u64":  value.True.ToValue(),
		"91i64 !~ 27u64": value.True.ToValue(),

		"5i64 !~ 5u32":  value.False.ToValue(),
		"-8i64 !~ 8u32": value.True.ToValue(),
		"3i64 !~ 71u32": value.True.ToValue(),

		"53000i64 !~ 32767u16": value.True.ToValue(),
		"5i64 !~ 5u16":         value.False.ToValue(),
		"-8i64 !~ 8u16":        value.True.ToValue(),
		"3i64 !~ 71u16":        value.True.ToValue(),

		"256i64 !~ 127u8": value.True.ToValue(),
		"5i64 !~ 5u8":     value.False.ToValue(),
		"-8i64 !~ 8u8":    value.True.ToValue(),
		"3i64 !~ 71u8":    value.True.ToValue(),

		// Int32
		"25i32 !~ 25":  value.False.ToValue(),
		"-25i32 !~ 25": value.True.ToValue(),
		"25i32 !~ -25": value.True.ToValue(),
		"25i32 !~ 28":  value.True.ToValue(),
		"28i32 !~ 25":  value.True.ToValue(),

		"25i32 !~ '25'": value.True.ToValue(),

		"7i32 !~ `7`": value.True.ToValue(),

		"-73i32 !~ 73.0": value.True.ToValue(),
		"73i32 !~ -73.0": value.True.ToValue(),
		"25i32 !~ 25.0":  value.False.ToValue(),
		"1i32 !~ 1.2":    value.True.ToValue(),

		"-73i32 !~ 73bf": value.True.ToValue(),
		"73i32 !~ -73bf": value.True.ToValue(),
		"25i32 !~ 25bf":  value.False.ToValue(),
		"1i32 !~ 1.2bf":  value.True.ToValue(),

		"-73i32 !~ 73f64": value.True.ToValue(),
		"73i32 !~ -73f64": value.True.ToValue(),
		"25i32 !~ 25f64":  value.False.ToValue(),
		"1i32 !~ 1.2f64":  value.True.ToValue(),

		"-73i32 !~ 73f32": value.True.ToValue(),
		"73i32 !~ -73f32": value.True.ToValue(),
		"25i32 !~ 25f32":  value.False.ToValue(),
		"1i32 !~ 1.2f32":  value.True.ToValue(),

		"1i32 !~ 1i64":   value.False.ToValue(),
		"4i32 !~ -4i64":  value.True.ToValue(),
		"-8i32 !~ 8i64":  value.True.ToValue(),
		"-8i32 !~ -8i64": value.False.ToValue(),
		"91i32 !~ 27i64": value.True.ToValue(),

		"5i32 !~ 5i32":  value.False.ToValue(),
		"4i32 !~ -4i32": value.True.ToValue(),
		"-8i32 !~ 8i32": value.True.ToValue(),
		"3i32 !~ 71i32": value.True.ToValue(),

		"5i32 !~ 5i16":  value.False.ToValue(),
		"4i32 !~ -4i16": value.True.ToValue(),
		"-8i32 !~ 8i16": value.True.ToValue(),
		"3i32 !~ 71i16": value.True.ToValue(),

		"5i32 !~ 5i8":  value.False.ToValue(),
		"4i32 !~ -4i8": value.True.ToValue(),
		"-8i32 !~ 8i8": value.True.ToValue(),
		"3i32 !~ 71i8": value.True.ToValue(),

		"1i32 !~ 1u64":   value.False.ToValue(),
		"-8i32 !~ 8u64":  value.True.ToValue(),
		"91i32 !~ 27u64": value.True.ToValue(),

		"5i32 !~ 5u32":  value.False.ToValue(),
		"-8i32 !~ 8u32": value.True.ToValue(),
		"3i32 !~ 71u32": value.True.ToValue(),

		"53000i32 !~ 32767u16": value.True.ToValue(),
		"5i32 !~ 5u16":         value.False.ToValue(),
		"-8i32 !~ 8u16":        value.True.ToValue(),
		"3i32 !~ 71u16":        value.True.ToValue(),

		"256i32 !~ 127u8": value.True.ToValue(),
		"5i32 !~ 5u8":     value.False.ToValue(),
		"-8i32 !~ 8u8":    value.True.ToValue(),
		"3i32 !~ 71u8":    value.True.ToValue(),

		// Int16
		"25i16 !~ 25":  value.False.ToValue(),
		"-25i16 !~ 25": value.True.ToValue(),
		"25i16 !~ -25": value.True.ToValue(),
		"25i16 !~ 28":  value.True.ToValue(),
		"28i16 !~ 25":  value.True.ToValue(),

		"25i16 !~ '25'": value.True.ToValue(),

		"7i16 !~ `7`": value.True.ToValue(),

		"-73i16 !~ 73.0": value.True.ToValue(),
		"73i16 !~ -73.0": value.True.ToValue(),
		"25i16 !~ 25.0":  value.False.ToValue(),
		"1i16 !~ 1.2":    value.True.ToValue(),

		"-73i16 !~ 73bf": value.True.ToValue(),
		"73i16 !~ -73bf": value.True.ToValue(),
		"25i16 !~ 25bf":  value.False.ToValue(),
		"1i16 !~ 1.2bf":  value.True.ToValue(),

		"-73i16 !~ 73f64": value.True.ToValue(),
		"73i16 !~ -73f64": value.True.ToValue(),
		"25i16 !~ 25f64":  value.False.ToValue(),
		"1i16 !~ 1.2f64":  value.True.ToValue(),

		"-73i16 !~ 73f32": value.True.ToValue(),
		"73i16 !~ -73f32": value.True.ToValue(),
		"25i16 !~ 25f32":  value.False.ToValue(),
		"1i16 !~ 1.2f32":  value.True.ToValue(),

		"1i16 !~ 1i64":   value.False.ToValue(),
		"4i16 !~ -4i64":  value.True.ToValue(),
		"-8i16 !~ 8i64":  value.True.ToValue(),
		"-8i16 !~ -8i64": value.False.ToValue(),
		"91i16 !~ 27i64": value.True.ToValue(),

		"5i16 !~ 5i32":  value.False.ToValue(),
		"4i16 !~ -4i32": value.True.ToValue(),
		"-8i16 !~ 8i32": value.True.ToValue(),
		"3i16 !~ 71i32": value.True.ToValue(),

		"5i16 !~ 5i16":  value.False.ToValue(),
		"4i16 !~ -4i16": value.True.ToValue(),
		"-8i16 !~ 8i16": value.True.ToValue(),
		"3i16 !~ 71i16": value.True.ToValue(),

		"5i16 !~ 5i8":  value.False.ToValue(),
		"4i16 !~ -4i8": value.True.ToValue(),
		"-8i16 !~ 8i8": value.True.ToValue(),
		"3i16 !~ 71i8": value.True.ToValue(),

		"1i16 !~ 1u64":   value.False.ToValue(),
		"-8i16 !~ 8u64":  value.True.ToValue(),
		"91i16 !~ 27u64": value.True.ToValue(),

		"5i16 !~ 5u32":  value.False.ToValue(),
		"-8i16 !~ 8u32": value.True.ToValue(),
		"3i16 !~ 71u32": value.True.ToValue(),

		"5i16 !~ 5u16":  value.False.ToValue(),
		"-8i16 !~ 8u16": value.True.ToValue(),
		"3i16 !~ 71u16": value.True.ToValue(),

		"256i16 !~ 127u8": value.True.ToValue(),
		"5i16 !~ 5u8":     value.False.ToValue(),
		"-8i16 !~ 8u8":    value.True.ToValue(),
		"3i16 !~ 71u8":    value.True.ToValue(),

		// Int8
		"25i8 !~ 25":  value.False.ToValue(),
		"-25i8 !~ 25": value.True.ToValue(),
		"25i8 !~ -25": value.True.ToValue(),
		"25i8 !~ 28":  value.True.ToValue(),
		"28i8 !~ 25":  value.True.ToValue(),

		"25i8 !~ '25'": value.True.ToValue(),

		"7i8 !~ `7`": value.True.ToValue(),

		"-73i8 !~ 73.0": value.True.ToValue(),
		"73i8 !~ -73.0": value.True.ToValue(),
		"25i8 !~ 25.0":  value.False.ToValue(),
		"1i8 !~ 1.2":    value.True.ToValue(),

		"-73i8 !~ 73bf": value.True.ToValue(),
		"73i8 !~ -73bf": value.True.ToValue(),
		"25i8 !~ 25bf":  value.False.ToValue(),
		"1i8 !~ 1.2bf":  value.True.ToValue(),

		"-73i8 !~ 73f64": value.True.ToValue(),
		"73i8 !~ -73f64": value.True.ToValue(),
		"25i8 !~ 25f64":  value.False.ToValue(),
		"1i8 !~ 1.2f64":  value.True.ToValue(),

		"-73i8 !~ 73f32": value.True.ToValue(),
		"73i8 !~ -73f32": value.True.ToValue(),
		"25i8 !~ 25f32":  value.False.ToValue(),
		"1i8 !~ 1.2f32":  value.True.ToValue(),

		"1i8 !~ 1i64":   value.False.ToValue(),
		"4i8 !~ -4i64":  value.True.ToValue(),
		"-8i8 !~ 8i64":  value.True.ToValue(),
		"-8i8 !~ -8i64": value.False.ToValue(),
		"91i8 !~ 27i64": value.True.ToValue(),

		"5i8 !~ 5i32":  value.False.ToValue(),
		"4i8 !~ -4i32": value.True.ToValue(),
		"-8i8 !~ 8i32": value.True.ToValue(),
		"3i8 !~ 71i32": value.True.ToValue(),

		"5i8 !~ 5i16":  value.False.ToValue(),
		"4i8 !~ -4i16": value.True.ToValue(),
		"-8i8 !~ 8i16": value.True.ToValue(),
		"3i8 !~ 71i16": value.True.ToValue(),

		"5i8 !~ 5i8":  value.False.ToValue(),
		"4i8 !~ -4i8": value.True.ToValue(),
		"-8i8 !~ 8i8": value.True.ToValue(),
		"3i8 !~ 71i8": value.True.ToValue(),

		"1i8 !~ 1u64":   value.False.ToValue(),
		"-8i8 !~ 8u64":  value.True.ToValue(),
		"91i8 !~ 27u64": value.True.ToValue(),

		"5i8 !~ 5u32":  value.False.ToValue(),
		"-8i8 !~ 8u32": value.True.ToValue(),
		"3i8 !~ 71u32": value.True.ToValue(),

		"5i8 !~ 5u16":  value.False.ToValue(),
		"-8i8 !~ 8u16": value.True.ToValue(),
		"3i8 !~ 71u16": value.True.ToValue(),

		"5i8 !~ 5u8":  value.False.ToValue(),
		"-8i8 !~ 8u8": value.True.ToValue(),
		"3i8 !~ 71u8": value.True.ToValue(),

		// UInt64
		"25u64 !~ 25":  value.False.ToValue(),
		"25u64 !~ -25": value.True.ToValue(),
		"25u64 !~ 28":  value.True.ToValue(),
		"28u64 !~ 25":  value.True.ToValue(),

		"25u64 !~ '25'": value.True.ToValue(),

		"7u64 !~ `7`": value.True.ToValue(),

		"73u64 !~ -73.0": value.True.ToValue(),
		"25u64 !~ 25.0":  value.False.ToValue(),
		"1u64 !~ 1.2":    value.True.ToValue(),

		"73u64 !~ -73bf": value.True.ToValue(),
		"25u64 !~ 25bf":  value.False.ToValue(),
		"1u64 !~ 1.2bf":  value.True.ToValue(),

		"73u64 !~ -73f64": value.True.ToValue(),
		"25u64 !~ 25f64":  value.False.ToValue(),
		"1u64 !~ 1.2f64":  value.True.ToValue(),

		"73u64 !~ -73f32": value.True.ToValue(),
		"25u64 !~ 25f32":  value.False.ToValue(),
		"1u64 !~ 1.2f32":  value.True.ToValue(),

		"1u64 !~ 1i64":   value.False.ToValue(),
		"4u64 !~ -4i64":  value.True.ToValue(),
		"91u64 !~ 27i64": value.True.ToValue(),

		"5u64 !~ 5i32":  value.False.ToValue(),
		"4u64 !~ -4i32": value.True.ToValue(),
		"3u64 !~ 71i32": value.True.ToValue(),

		"5u64 !~ 5i16":  value.False.ToValue(),
		"4u64 !~ -4i16": value.True.ToValue(),
		"3u64 !~ 71i16": value.True.ToValue(),

		"5u64 !~ 5i8":  value.False.ToValue(),
		"4u64 !~ -4i8": value.True.ToValue(),
		"3u64 !~ 71i8": value.True.ToValue(),

		"1u64 !~ 1u64":   value.False.ToValue(),
		"91u64 !~ 27u64": value.True.ToValue(),

		"5u64 !~ 5u32":  value.False.ToValue(),
		"3u64 !~ 71u32": value.True.ToValue(),

		"53000u64 !~ 32767u16": value.True.ToValue(),
		"5u64 !~ 5u16":         value.False.ToValue(),
		"3u64 !~ 71u16":        value.True.ToValue(),

		"256u64 !~ 127u8": value.True.ToValue(),
		"5u64 !~ 5u8":     value.False.ToValue(),
		"3u64 !~ 71u8":    value.True.ToValue(),

		// UInt32
		"25u32 !~ 25":  value.False.ToValue(),
		"25u32 !~ -25": value.True.ToValue(),
		"25u32 !~ 28":  value.True.ToValue(),
		"28u32 !~ 25":  value.True.ToValue(),

		"25u32 !~ '25'": value.True.ToValue(),

		"7u32 !~ `7`": value.True.ToValue(),

		"73u32 !~ -73.0": value.True.ToValue(),
		"25u32 !~ 25.0":  value.False.ToValue(),
		"1u32 !~ 1.2":    value.True.ToValue(),

		"73u32 !~ -73bf": value.True.ToValue(),
		"25u32 !~ 25bf":  value.False.ToValue(),
		"1u32 !~ 1.2bf":  value.True.ToValue(),

		"73u32 !~ -73f64": value.True.ToValue(),
		"25u32 !~ 25f64":  value.False.ToValue(),
		"1u32 !~ 1.2f64":  value.True.ToValue(),

		"73u32 !~ -73f32": value.True.ToValue(),
		"25u32 !~ 25f32":  value.False.ToValue(),
		"1u32 !~ 1.2f32":  value.True.ToValue(),

		"1u32 !~ 1i64":   value.False.ToValue(),
		"4u32 !~ -4i64":  value.True.ToValue(),
		"91u32 !~ 27i64": value.True.ToValue(),

		"5u32 !~ 5i32":  value.False.ToValue(),
		"4u32 !~ -4i32": value.True.ToValue(),
		"3u32 !~ 71i32": value.True.ToValue(),

		"5u32 !~ 5i16":  value.False.ToValue(),
		"4u32 !~ -4i16": value.True.ToValue(),
		"3u32 !~ 71i16": value.True.ToValue(),

		"5u32 !~ 5i8":  value.False.ToValue(),
		"4u32 !~ -4i8": value.True.ToValue(),
		"3u32 !~ 71i8": value.True.ToValue(),

		"1u32 !~ 1u64":   value.False.ToValue(),
		"91u32 !~ 27u64": value.True.ToValue(),

		"5u32 !~ 5u32":  value.False.ToValue(),
		"3u32 !~ 71u32": value.True.ToValue(),

		"53000u32 !~ 32767u16": value.True.ToValue(),
		"5u32 !~ 5u16":         value.False.ToValue(),
		"3u32 !~ 71u16":        value.True.ToValue(),

		"256u32 !~ 127u8": value.True.ToValue(),
		"5u32 !~ 5u8":     value.False.ToValue(),
		"3u32 !~ 71u8":    value.True.ToValue(),

		// UInt16
		"25u16 !~ 25":  value.False.ToValue(),
		"25u16 !~ -25": value.True.ToValue(),
		"25u16 !~ 28":  value.True.ToValue(),
		"28u16 !~ 25":  value.True.ToValue(),

		"25u16 !~ '25'": value.True.ToValue(),

		"7u16 !~ `7`": value.True.ToValue(),

		"73u16 !~ -73.0": value.True.ToValue(),
		"25u16 !~ 25.0":  value.False.ToValue(),
		"1u16 !~ 1.2":    value.True.ToValue(),

		"73u16 !~ -73bf": value.True.ToValue(),
		"25u16 !~ 25bf":  value.False.ToValue(),
		"1u16 !~ 1.2bf":  value.True.ToValue(),

		"73u16 !~ -73f64": value.True.ToValue(),
		"25u16 !~ 25f64":  value.False.ToValue(),
		"1u16 !~ 1.2f64":  value.True.ToValue(),

		"73u16 !~ -73f32": value.True.ToValue(),
		"25u16 !~ 25f32":  value.False.ToValue(),
		"1u16 !~ 1.2f32":  value.True.ToValue(),

		"1u16 !~ 1i64":   value.False.ToValue(),
		"4u16 !~ -4i64":  value.True.ToValue(),
		"91u16 !~ 27i64": value.True.ToValue(),

		"5u16 !~ 5i32":  value.False.ToValue(),
		"4u16 !~ -4i32": value.True.ToValue(),
		"3u16 !~ 71i32": value.True.ToValue(),

		"5u16 !~ 5i16":  value.False.ToValue(),
		"4u16 !~ -4i16": value.True.ToValue(),
		"3u16 !~ 71i16": value.True.ToValue(),

		"5u16 !~ 5i8":  value.False.ToValue(),
		"4u16 !~ -4i8": value.True.ToValue(),
		"3u16 !~ 71i8": value.True.ToValue(),

		"1u16 !~ 1u64":   value.False.ToValue(),
		"91u16 !~ 27u64": value.True.ToValue(),

		"5u16 !~ 5u32":  value.False.ToValue(),
		"3u16 !~ 71u32": value.True.ToValue(),

		"53000u16 !~ 32767u16": value.True.ToValue(),
		"5u16 !~ 5u16":         value.False.ToValue(),
		"3u16 !~ 71u16":        value.True.ToValue(),

		"256u16 !~ 127u8": value.True.ToValue(),
		"5u16 !~ 5u8":     value.False.ToValue(),
		"3u16 !~ 71u8":    value.True.ToValue(),

		// UInt8
		"25u8 !~ 25":  value.False.ToValue(),
		"25u8 !~ -25": value.True.ToValue(),
		"25u8 !~ 28":  value.True.ToValue(),
		"28u8 !~ 25":  value.True.ToValue(),

		"25u8 !~ '25'": value.True.ToValue(),

		"7u8 !~ `7`": value.True.ToValue(),

		"73u8 !~ -73.0": value.True.ToValue(),
		"25u8 !~ 25.0":  value.False.ToValue(),
		"1u8 !~ 1.2":    value.True.ToValue(),

		"73u8 !~ -73bf": value.True.ToValue(),
		"25u8 !~ 25bf":  value.False.ToValue(),
		"1u8 !~ 1.2bf":  value.True.ToValue(),

		"73u8 !~ -73f64": value.True.ToValue(),
		"25u8 !~ 25f64":  value.False.ToValue(),
		"1u8 !~ 1.2f64":  value.True.ToValue(),

		"73u8 !~ -73f32": value.True.ToValue(),
		"25u8 !~ 25f32":  value.False.ToValue(),
		"1u8 !~ 1.2f32":  value.True.ToValue(),

		"1u8 !~ 1i64":   value.False.ToValue(),
		"4u8 !~ -4i64":  value.True.ToValue(),
		"91u8 !~ 27i64": value.True.ToValue(),

		"5u8 !~ 5i32":  value.False.ToValue(),
		"4u8 !~ -4i32": value.True.ToValue(),
		"3u8 !~ 71i32": value.True.ToValue(),

		"5u8 !~ 5i16":  value.False.ToValue(),
		"4u8 !~ -4i16": value.True.ToValue(),
		"3u8 !~ 71i16": value.True.ToValue(),

		"5u8 !~ 5i8":  value.False.ToValue(),
		"4u8 !~ -4i8": value.True.ToValue(),
		"3u8 !~ 71i8": value.True.ToValue(),

		"1u8 !~ 1u64":   value.False.ToValue(),
		"91u8 !~ 27u64": value.True.ToValue(),

		"5u8 !~ 5u32":  value.False.ToValue(),
		"3u8 !~ 71u32": value.True.ToValue(),

		"5u8 !~ 5u16":  value.False.ToValue(),
		"3u8 !~ 71u16": value.True.ToValue(),

		"5u8 !~ 5u8":  value.False.ToValue(),
		"3u8 !~ 71u8": value.True.ToValue(),

		// Float
		"-73.0 !~ 73.0": value.True.ToValue(),
		"73.0 !~ -73.0": value.True.ToValue(),
		"25.0 !~ 25.0":  value.False.ToValue(),
		"1.0 !~ 1.2":    value.True.ToValue(),
		"1.2 !~ 1.0":    value.True.ToValue(),
		"78.5 !~ 78.5":  value.False.ToValue(),

		"8.25 !~ '8.25'": value.True.ToValue(),

		"4.0 !~ `4`": value.True.ToValue(),

		"25.0 !~ 25":  value.False.ToValue(),
		"32.3 !~ 32":  value.True.ToValue(),
		"-25.0 !~ 25": value.True.ToValue(),
		"25.0 !~ -25": value.True.ToValue(),
		"25.0 !~ 28":  value.True.ToValue(),
		"28.0 !~ 25":  value.True.ToValue(),

		"-73.0 !~ 73bf":  value.True.ToValue(),
		"73.0 !~ -73bf":  value.True.ToValue(),
		"25.0 !~ 25bf":   value.False.ToValue(),
		"1.0 !~ 1.2bf":   value.True.ToValue(),
		"15.5 !~ 15.5bf": value.False.ToValue(),

		"-73.0 !~ 73f64":    value.True.ToValue(),
		"73.0 !~ -73f64":    value.True.ToValue(),
		"25.0 !~ 25f64":     value.False.ToValue(),
		"1.0 !~ 1.2f64":     value.True.ToValue(),
		"15.26 !~ 15.26f64": value.False.ToValue(),

		"-73.0 !~ 73f32":  value.True.ToValue(),
		"73.0 !~ -73f32":  value.True.ToValue(),
		"25.0 !~ 25f32":   value.False.ToValue(),
		"1.0 !~ 1.2f32":   value.True.ToValue(),
		"15.5 !~ 15.5f32": value.False.ToValue(),

		"1.0 !~ 1i64":   value.False.ToValue(),
		"1.5 !~ 1i64":   value.True.ToValue(),
		"4.0 !~ -4i64":  value.True.ToValue(),
		"-8.0 !~ 8i64":  value.True.ToValue(),
		"-8.0 !~ -8i64": value.False.ToValue(),
		"91.0 !~ 27i64": value.True.ToValue(),

		"1.0 !~ 1i32":   value.False.ToValue(),
		"1.5 !~ 1i32":   value.True.ToValue(),
		"4.0 !~ -4i32":  value.True.ToValue(),
		"-8.0 !~ 8i32":  value.True.ToValue(),
		"-8.0 !~ -8i32": value.False.ToValue(),
		"91.0 !~ 27i32": value.True.ToValue(),

		"1.0 !~ 1i16":   value.False.ToValue(),
		"1.5 !~ 1i16":   value.True.ToValue(),
		"4.0 !~ -4i16":  value.True.ToValue(),
		"-8.0 !~ 8i16":  value.True.ToValue(),
		"-8.0 !~ -8i16": value.False.ToValue(),
		"91.0 !~ 27i16": value.True.ToValue(),

		"1.0 !~ 1i8":   value.False.ToValue(),
		"1.5 !~ 1i8":   value.True.ToValue(),
		"4.0 !~ -4i8":  value.True.ToValue(),
		"-8.0 !~ 8i8":  value.True.ToValue(),
		"-8.0 !~ -8i8": value.False.ToValue(),
		"91.0 !~ 27i8": value.True.ToValue(),

		"1.0 !~ 1u64":   value.False.ToValue(),
		"1.5 !~ 1u64":   value.True.ToValue(),
		"-8.0 !~ 8u64":  value.True.ToValue(),
		"91.0 !~ 27u64": value.True.ToValue(),

		"1.0 !~ 1u32":   value.False.ToValue(),
		"1.5 !~ 1u32":   value.True.ToValue(),
		"-8.0 !~ 8u32":  value.True.ToValue(),
		"91.0 !~ 27u32": value.True.ToValue(),

		"53000.0 !~ 32767u16": value.True.ToValue(),
		"1.0 !~ 1u16":         value.False.ToValue(),
		"1.5 !~ 1u16":         value.True.ToValue(),
		"-8.0 !~ 8u16":        value.True.ToValue(),
		"91.0 !~ 27u16":       value.True.ToValue(),

		"256.0 !~ 127u8": value.True.ToValue(),
		"1.0 !~ 1u8":     value.False.ToValue(),
		"1.5 !~ 1u8":     value.True.ToValue(),
		"-8.0 !~ 8u8":    value.True.ToValue(),
		"91.0 !~ 27u8":   value.True.ToValue(),

		// Float64
		"-73f64 !~ 73.0":  value.True.ToValue(),
		"73f64 !~ -73.0":  value.True.ToValue(),
		"25f64 !~ 25.0":   value.False.ToValue(),
		"1f64 !~ 1.2":     value.True.ToValue(),
		"1.2f64 !~ 1.0":   value.True.ToValue(),
		"78.5f64 !~ 78.5": value.False.ToValue(),

		"8.25f64 !~ '8.25'": value.True.ToValue(),

		"4f64 !~ `4`": value.True.ToValue(),

		"25f64 !~ 25":   value.False.ToValue(),
		"32.3f64 !~ 32": value.True.ToValue(),
		"-25f64 !~ 25":  value.True.ToValue(),
		"25f64 !~ -25":  value.True.ToValue(),
		"25f64 !~ 28":   value.True.ToValue(),
		"28f64 !~ 25":   value.True.ToValue(),

		"-73f64 !~ 73bf":    value.True.ToValue(),
		"73f64 !~ -73bf":    value.True.ToValue(),
		"25f64 !~ 25bf":     value.False.ToValue(),
		"1f64 !~ 1.2bf":     value.True.ToValue(),
		"15.5f64 !~ 15.5bf": value.False.ToValue(),

		"-73f64 !~ 73f64":      value.True.ToValue(),
		"73f64 !~ -73f64":      value.True.ToValue(),
		"25f64 !~ 25f64":       value.False.ToValue(),
		"1f64 !~ 1.2f64":       value.True.ToValue(),
		"15.26f64 !~ 15.26f64": value.False.ToValue(),

		"-73f64 !~ 73f32":    value.True.ToValue(),
		"73f64 !~ -73f32":    value.True.ToValue(),
		"25f64 !~ 25f32":     value.False.ToValue(),
		"1f64 !~ 1.2f32":     value.True.ToValue(),
		"15.5f64 !~ 15.5f32": value.False.ToValue(),

		"1f64 !~ 1i64":   value.False.ToValue(),
		"1.5f64 !~ 1i64": value.True.ToValue(),
		"4f64 !~ -4i64":  value.True.ToValue(),
		"-8f64 !~ 8i64":  value.True.ToValue(),
		"-8f64 !~ -8i64": value.False.ToValue(),
		"91f64 !~ 27i64": value.True.ToValue(),

		"1f64 !~ 1i32":   value.False.ToValue(),
		"1.5f64 !~ 1i32": value.True.ToValue(),
		"4f64 !~ -4i32":  value.True.ToValue(),
		"-8f64 !~ 8i32":  value.True.ToValue(),
		"-8f64 !~ -8i32": value.False.ToValue(),
		"91f64 !~ 27i32": value.True.ToValue(),

		"1f64 !~ 1i16":   value.False.ToValue(),
		"1.5f64 !~ 1i16": value.True.ToValue(),
		"4f64 !~ -4i16":  value.True.ToValue(),
		"-8f64 !~ 8i16":  value.True.ToValue(),
		"-8f64 !~ -8i16": value.False.ToValue(),
		"91f64 !~ 27i16": value.True.ToValue(),

		"1f64 !~ 1i8":   value.False.ToValue(),
		"1.5f64 !~ 1i8": value.True.ToValue(),
		"4f64 !~ -4i8":  value.True.ToValue(),
		"-8f64 !~ 8i8":  value.True.ToValue(),
		"-8f64 !~ -8i8": value.False.ToValue(),
		"91f64 !~ 27i8": value.True.ToValue(),

		"1f64 !~ 1u64":   value.False.ToValue(),
		"1.5f64 !~ 1u64": value.True.ToValue(),
		"-8f64 !~ 8u64":  value.True.ToValue(),
		"91f64 !~ 27u64": value.True.ToValue(),

		"1f64 !~ 1u32":   value.False.ToValue(),
		"1.5f64 !~ 1u32": value.True.ToValue(),
		"-8f64 !~ 8u32":  value.True.ToValue(),
		"91f64 !~ 27u32": value.True.ToValue(),

		"53000f64 !~ 32767u16": value.True.ToValue(),
		"1f64 !~ 1u16":         value.False.ToValue(),
		"1.5f64 !~ 1u16":       value.True.ToValue(),
		"-8f64 !~ 8u16":        value.True.ToValue(),
		"91f64 !~ 27u16":       value.True.ToValue(),

		"256f64 !~ 127u8": value.True.ToValue(),
		"1f64 !~ 1u8":     value.False.ToValue(),
		"1.5f64 !~ 1u8":   value.True.ToValue(),
		"-8f64 !~ 8u8":    value.True.ToValue(),
		"91f64 !~ 27u8":   value.True.ToValue(),

		// Float32
		"-73f32 !~ 73.0":  value.True.ToValue(),
		"73f32 !~ -73.0":  value.True.ToValue(),
		"25f32 !~ 25.0":   value.False.ToValue(),
		"1f32 !~ 1.2":     value.True.ToValue(),
		"1.2f32 !~ 1.0":   value.True.ToValue(),
		"78.5f32 !~ 78.5": value.False.ToValue(),

		"8.25f32 !~ '8.25'": value.True.ToValue(),

		"4f32 !~ `4`": value.True.ToValue(),

		"25f32 !~ 25":   value.False.ToValue(),
		"32.3f32 !~ 32": value.True.ToValue(),
		"-25f32 !~ 25":  value.True.ToValue(),
		"25f32 !~ -25":  value.True.ToValue(),
		"25f32 !~ 28":   value.True.ToValue(),
		"28f32 !~ 25":   value.True.ToValue(),

		"-73f32 !~ 73bf":    value.True.ToValue(),
		"73f32 !~ -73bf":    value.True.ToValue(),
		"25f32 !~ 25bf":     value.False.ToValue(),
		"1f32 !~ 1.2bf":     value.True.ToValue(),
		"15.5f32 !~ 15.5bf": value.False.ToValue(),

		"-73f32 !~ 73f64":    value.True.ToValue(),
		"73f32 !~ -73f64":    value.True.ToValue(),
		"25f32 !~ 25f64":     value.False.ToValue(),
		"1f32 !~ 1.2f64":     value.True.ToValue(),
		"15.5f32 !~ 15.5f64": value.False.ToValue(),

		"-73f32 !~ 73f32":    value.True.ToValue(),
		"73f32 !~ -73f32":    value.True.ToValue(),
		"25f32 !~ 25f32":     value.False.ToValue(),
		"1f32 !~ 1.2f32":     value.True.ToValue(),
		"15.5f32 !~ 15.5f32": value.False.ToValue(),

		"1f32 !~ 1i64":   value.False.ToValue(),
		"1.5f32 !~ 1i64": value.True.ToValue(),
		"4f32 !~ -4i64":  value.True.ToValue(),
		"-8f32 !~ 8i64":  value.True.ToValue(),
		"-8f32 !~ -8i64": value.False.ToValue(),
		"91f32 !~ 27i64": value.True.ToValue(),

		"1f32 !~ 1i32":   value.False.ToValue(),
		"1.5f32 !~ 1i32": value.True.ToValue(),
		"4f32 !~ -4i32":  value.True.ToValue(),
		"-8f32 !~ 8i32":  value.True.ToValue(),
		"-8f32 !~ -8i32": value.False.ToValue(),
		"91f32 !~ 27i32": value.True.ToValue(),

		"1f32 !~ 1i16":   value.False.ToValue(),
		"1.5f32 !~ 1i16": value.True.ToValue(),
		"4f32 !~ -4i16":  value.True.ToValue(),
		"-8f32 !~ 8i16":  value.True.ToValue(),
		"-8f32 !~ -8i16": value.False.ToValue(),
		"91f32 !~ 27i16": value.True.ToValue(),

		"1f32 !~ 1i8":   value.False.ToValue(),
		"1.5f32 !~ 1i8": value.True.ToValue(),
		"4f32 !~ -4i8":  value.True.ToValue(),
		"-8f32 !~ 8i8":  value.True.ToValue(),
		"-8f32 !~ -8i8": value.False.ToValue(),
		"91f32 !~ 27i8": value.True.ToValue(),

		"1f32 !~ 1u64":   value.False.ToValue(),
		"1.5f32 !~ 1u64": value.True.ToValue(),
		"-8f32 !~ 8u64":  value.True.ToValue(),
		"91f32 !~ 27u64": value.True.ToValue(),

		"1f32 !~ 1u32":   value.False.ToValue(),
		"1.5f32 !~ 1u32": value.True.ToValue(),
		"-8f32 !~ 8u32":  value.True.ToValue(),
		"91f32 !~ 27u32": value.True.ToValue(),

		"53000f32 !~ 32767u16": value.True.ToValue(),
		"1f32 !~ 1u16":         value.False.ToValue(),
		"1.5f32 !~ 1u16":       value.True.ToValue(),
		"-8f32 !~ 8u16":        value.True.ToValue(),
		"91f32 !~ 27u16":       value.True.ToValue(),

		"256f32 !~ 127u8": value.True.ToValue(),
		"1f32 !~ 1u8":     value.False.ToValue(),
		"1.5f32 !~ 1u8":   value.True.ToValue(),
		"-8f32 !~ 8u8":    value.True.ToValue(),
		"91f32 !~ 27u8":   value.True.ToValue(),
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}

func TestVMSource_Equal(t *testing.T) {
	tests := simpleSourceTestTable{
		// String
		"'25' == '25'":   value.True.ToValue(),
		"'25' == '25.0'": value.False.ToValue(),
		"'25' == '7'":    value.False.ToValue(),

		"'7' == `7`": value.False.ToValue(),

		"'25' == 25.0": value.False.ToValue(),

		"'25' == 25bf": value.False.ToValue(),

		"'25' == 25f64": value.False.ToValue(),

		"'25' == 25f32": value.False.ToValue(),

		"'1' == 1i64": value.False.ToValue(),

		"'5' == 5i32": value.False.ToValue(),

		"'5' == 5i16": value.False.ToValue(),

		"'5' == 5i8": value.False.ToValue(),

		"'1' == 1u64": value.False.ToValue(),

		"'5' == 5u32": value.False.ToValue(),

		"'5' == 5u16": value.False.ToValue(),

		"'5' == 5u8": value.False.ToValue(),

		// Char
		"`2` == 25": value.False.ToValue(),

		"`2` == '2'": value.False.ToValue(),

		"`7` == `7`": value.True.ToValue(),
		"`b` == `b`": value.True.ToValue(),
		"`c` == `g`": value.False.ToValue(),
		"`7` == `8`": value.False.ToValue(),

		"`2` == 2.0": value.False.ToValue(),

		"`3` == 3bf": value.False.ToValue(),

		"`9` == 9f64": value.False.ToValue(),

		"`1` == 1f32": value.False.ToValue(),

		"`1` == 1i64": value.False.ToValue(),

		"`5` == 5i32": value.False.ToValue(),

		"`5` == 5i16": value.False.ToValue(),

		"`5` == 5i8": value.False.ToValue(),

		"`1` == 1u64": value.False.ToValue(),

		"`5` == 5u32": value.False.ToValue(),

		"`5` == 5u16": value.False.ToValue(),

		"`5` == 5u8": value.False.ToValue(),

		// Int
		"25 == 25":  value.True.ToValue(),
		"-25 == 25": value.False.ToValue(),
		"25 == -25": value.False.ToValue(),
		"25 == 28":  value.False.ToValue(),
		"28 == 25":  value.False.ToValue(),

		"25 == '25'": value.False.ToValue(),

		"7 == `7`": value.False.ToValue(),

		"25 == 25.0": value.False.ToValue(),

		"25 == 25bf": value.False.ToValue(),

		"25 == 25f64": value.False.ToValue(),

		"25 == 25f32": value.False.ToValue(),

		"1 == 1i64": value.False.ToValue(),

		"5 == 5i32": value.False.ToValue(),

		"5 == 5i16": value.False.ToValue(),

		"5 == 5i8": value.False.ToValue(),

		"1 == 1u64": value.False.ToValue(),

		"5 == 5u32": value.False.ToValue(),

		"5 == 5u16": value.False.ToValue(),

		"5 == 5u8": value.False.ToValue(),

		// Int64
		"25i64 == 25": value.False.ToValue(),

		"25i64 == '25'": value.False.ToValue(),

		"7i64 == `7`": value.False.ToValue(),

		"25i64 == 25.0": value.False.ToValue(),

		"25i64 == 25bf": value.False.ToValue(),

		"25i64 == 25f64": value.False.ToValue(),

		"25i64 == 25f32": value.False.ToValue(),

		"1i64 == 1i64":   value.True.ToValue(),
		"4i64 == -4i64":  value.False.ToValue(),
		"-8i64 == 8i64":  value.False.ToValue(),
		"-8i64 == -8i64": value.True.ToValue(),
		"91i64 == 27i64": value.False.ToValue(),

		"5i64 == 5i32": value.False.ToValue(),

		"5i64 == 5i16": value.False.ToValue(),

		"5i64 == 5i8": value.False.ToValue(),

		"1i64 == 1u64": value.False.ToValue(),

		"5i64 == 5u32": value.False.ToValue(),

		"5i64 == 5u16": value.False.ToValue(),

		"5i64 == 5u8": value.False.ToValue(),

		// Int32
		"25i32 == 25": value.False.ToValue(),

		"25i32 == '25'": value.False.ToValue(),

		"7i32 == `7`": value.False.ToValue(),

		"25i32 == 25.0": value.False.ToValue(),

		"25i32 == 25bf": value.False.ToValue(),

		"25i32 == 25f64": value.False.ToValue(),

		"25i32 == 25f32": value.False.ToValue(),

		"1i32 == 1i64": value.False.ToValue(),

		"5i32 == 5i32":  value.True.ToValue(),
		"4i32 == -4i32": value.False.ToValue(),
		"-8i32 == 8i32": value.False.ToValue(),
		"3i32 == 71i32": value.False.ToValue(),

		"5i32 == 5i16": value.False.ToValue(),

		"5i32 == 5i8": value.False.ToValue(),

		"1i32 == 1u64": value.False.ToValue(),

		"5i32 == 5u32": value.False.ToValue(),

		"5i32 == 5u16": value.False.ToValue(),

		"5i32 == 5u8": value.False.ToValue(),

		// Int16
		"25i16 == 25": value.False.ToValue(),

		"25i16 == '25'": value.False.ToValue(),

		"7i16 == `7`": value.False.ToValue(),

		"25i16 == 25.0": value.False.ToValue(),

		"25i16 == 25bf": value.False.ToValue(),

		"25i16 == 25f64": value.False.ToValue(),

		"25i16 == 25f32": value.False.ToValue(),

		"1i16 == 1i64": value.False.ToValue(),

		"5i16 == 5i32": value.False.ToValue(),

		"5i16 == 5i16":  value.True.ToValue(),
		"4i16 == -4i16": value.False.ToValue(),
		"-8i16 == 8i16": value.False.ToValue(),
		"3i16 == 71i16": value.False.ToValue(),

		"5i16 == 5i8": value.False.ToValue(),

		"1i16 == 1u64": value.False.ToValue(),

		"5i16 == 5u32": value.False.ToValue(),

		"5i16 == 5u16": value.False.ToValue(),

		"5i16 == 5u8": value.False.ToValue(),

		// Int8
		"25i8 == 25": value.False.ToValue(),

		"25i8 == '25'": value.False.ToValue(),

		"7i8 == `7`": value.False.ToValue(),

		"25i8 == 25.0": value.False.ToValue(),

		"25i8 == 25bf": value.False.ToValue(),

		"25i8 == 25f64": value.False.ToValue(),

		"25i8 == 25f32": value.False.ToValue(),

		"1i8 == 1i64": value.False.ToValue(),

		"5i8 == 5i32": value.False.ToValue(),

		"5i8 == 5i16": value.False.ToValue(),

		"5i8 == 5i8":  value.True.ToValue(),
		"4i8 == -4i8": value.False.ToValue(),
		"-8i8 == 8i8": value.False.ToValue(),
		"3i8 == 71i8": value.False.ToValue(),

		"1i8 == 1u64": value.False.ToValue(),

		"5i8 == 5u32": value.False.ToValue(),

		"5i8 == 5u16": value.False.ToValue(),

		"5i8 == 5u8": value.False.ToValue(),

		// UInt64
		"25u64 == 25": value.False.ToValue(),

		"25u64 == '25'": value.False.ToValue(),

		"7u64 == `7`": value.False.ToValue(),

		"25u64 == 25.0": value.False.ToValue(),

		"25u64 == 25bf": value.False.ToValue(),

		"25u64 == 25f64": value.False.ToValue(),

		"25u64 == 25f32": value.False.ToValue(),

		"1u64 == 1i64": value.False.ToValue(),

		"5u64 == 5i32": value.False.ToValue(),

		"5u64 == 5i16": value.False.ToValue(),

		"5u64 == 5i8": value.False.ToValue(),

		"1u64 == 1u64":   value.True.ToValue(),
		"91u64 == 27u64": value.False.ToValue(),

		"5u64 == 5u32": value.False.ToValue(),

		"5u64 == 5u16": value.False.ToValue(),

		"5u64 == 5u8": value.False.ToValue(),

		// UInt32
		"25u32 == 25": value.False.ToValue(),

		"25u32 == '25'": value.False.ToValue(),

		"7u32 == `7`": value.False.ToValue(),

		"25u32 == 25.0": value.False.ToValue(),

		"25u32 == 25bf": value.False.ToValue(),

		"25u32 == 25f64": value.False.ToValue(),

		"25u32 == 25f32": value.False.ToValue(),

		"1u32 == 1i64": value.False.ToValue(),

		"5u32 == 5i32": value.False.ToValue(),

		"5u32 == 5i16": value.False.ToValue(),

		"5u32 == 5i8": value.False.ToValue(),

		"1u32 == 1u64": value.False.ToValue(),

		"5u32 == 5u32":  value.True.ToValue(),
		"3u32 == 71u32": value.False.ToValue(),

		"5u32 == 5u16": value.False.ToValue(),

		"5u32 == 5u8": value.False.ToValue(),

		// UInt16
		"25u16 == 25": value.False.ToValue(),

		"25u16 == '25'": value.False.ToValue(),

		"7u16 == `7`": value.False.ToValue(),

		"25u16 == 25.0": value.False.ToValue(),

		"25u16 == 25bf": value.False.ToValue(),

		"25u16 == 25f64": value.False.ToValue(),

		"25u16 == 25f32": value.False.ToValue(),

		"1u16 == 1i64": value.False.ToValue(),

		"5u16 == 5i32": value.False.ToValue(),

		"5u16 == 5i16": value.False.ToValue(),

		"5u16 == 5i8": value.False.ToValue(),

		"1u16 == 1u64": value.False.ToValue(),

		"5u16 == 5u32": value.False.ToValue(),

		"53000u16 == 32767u16": value.False.ToValue(),
		"5u16 == 5u16":         value.True.ToValue(),
		"3u16 == 71u16":        value.False.ToValue(),

		"5u16 == 5u8": value.False.ToValue(),

		// UInt8
		"25u8 == 25": value.False.ToValue(),

		"25u8 == '25'": value.False.ToValue(),

		"7u8 == `7`": value.False.ToValue(),

		"25u8 == 25.0": value.False.ToValue(),

		"25u8 == 25bf": value.False.ToValue(),

		"25u8 == 25f64": value.False.ToValue(),

		"25u8 == 25f32": value.False.ToValue(),

		"1u8 == 1i64": value.False.ToValue(),

		"5u8 == 5i32": value.False.ToValue(),

		"5u8 == 5i16": value.False.ToValue(),

		"5u8 == 5i8": value.False.ToValue(),

		"1u8 == 1u64": value.False.ToValue(),

		"5u8 == 5u32": value.False.ToValue(),

		"5u8 == 5u16": value.False.ToValue(),

		"5u8 == 5u8":  value.True.ToValue(),
		"3u8 == 71u8": value.False.ToValue(),

		// Float
		"-73.0 == 73.0": value.False.ToValue(),
		"73.0 == -73.0": value.False.ToValue(),
		"25.0 == 25.0":  value.True.ToValue(),
		"1.0 == 1.2":    value.False.ToValue(),
		"1.2 == 1.0":    value.False.ToValue(),
		"78.5 == 78.5":  value.True.ToValue(),

		"8.25 == '8.25'": value.False.ToValue(),

		"4.0 == `4`": value.False.ToValue(),

		"25.0 == 25": value.False.ToValue(),

		"25.0 == 25bf":   value.False.ToValue(),
		"15.5 == 15.5bf": value.False.ToValue(),

		"25.0 == 25f64":     value.False.ToValue(),
		"15.26 == 15.26f64": value.False.ToValue(),

		"25.0 == 25f32":   value.False.ToValue(),
		"15.5 == 15.5f32": value.False.ToValue(),

		"1.0 == 1i64":   value.False.ToValue(),
		"-8.0 == -8i64": value.False.ToValue(),

		"1.0 == 1i32":   value.False.ToValue(),
		"-8.0 == -8i32": value.False.ToValue(),

		"1.0 == 1i16":   value.False.ToValue(),
		"-8.0 == -8i16": value.False.ToValue(),

		"1.0 == 1i8":   value.False.ToValue(),
		"-8.0 == -8i8": value.False.ToValue(),

		"1.0 == 1u64": value.False.ToValue(),

		"1.0 == 1u32": value.False.ToValue(),

		"1.0 == 1u16": value.False.ToValue(),

		"1.0 == 1u8": value.False.ToValue(),

		// Float64
		"25f64 == 25.0":   value.False.ToValue(),
		"78.5f64 == 78.5": value.False.ToValue(),

		"8.25f64 == '8.25'": value.False.ToValue(),

		"4f64 == `4`": value.False.ToValue(),

		"25f64 == 25": value.False.ToValue(),

		"25f64 == 25bf":     value.False.ToValue(),
		"15.5f64 == 15.5bf": value.False.ToValue(),

		"-73f64 == 73f64":      value.False.ToValue(),
		"73f64 == -73f64":      value.False.ToValue(),
		"25f64 == 25f64":       value.True.ToValue(),
		"1f64 == 1.2f64":       value.False.ToValue(),
		"15.26f64 == 15.26f64": value.True.ToValue(),

		"25f64 == 25f32":     value.False.ToValue(),
		"15.5f64 == 15.5f32": value.False.ToValue(),

		"1f64 == 1i64":   value.False.ToValue(),
		"-8f64 == -8i64": value.False.ToValue(),

		"1f64 == 1i32":   value.False.ToValue(),
		"-8f64 == -8i32": value.False.ToValue(),

		"1f64 == 1i16":   value.False.ToValue(),
		"-8f64 == -8i16": value.False.ToValue(),

		"1f64 == 1i8":   value.False.ToValue(),
		"-8f64 == -8i8": value.False.ToValue(),

		"1f64 == 1u64": value.False.ToValue(),

		"1f64 == 1u32": value.False.ToValue(),

		"1f64 == 1u16": value.False.ToValue(),

		"1f64 == 1u8": value.False.ToValue(),

		// Float32
		"25f32 == 25.0":   value.False.ToValue(),
		"78.5f32 == 78.5": value.False.ToValue(),

		"8.25f32 == '8.25'": value.False.ToValue(),

		"4f32 == `4`": value.False.ToValue(),

		"25f32 == 25": value.False.ToValue(),

		"25f32 == 25bf":     value.False.ToValue(),
		"15.5f32 == 15.5bf": value.False.ToValue(),

		"25f32 == 25f64":     value.False.ToValue(),
		"15.5f32 == 15.5f64": value.False.ToValue(),

		"-73f32 == 73f32":    value.False.ToValue(),
		"73f32 == -73f32":    value.False.ToValue(),
		"25f32 == 25f32":     value.True.ToValue(),
		"1f32 == 1.2f32":     value.False.ToValue(),
		"15.5f32 == 15.5f32": value.True.ToValue(),

		"1f32 == 1i64":   value.False.ToValue(),
		"-8f32 == -8i64": value.False.ToValue(),

		"1f32 == 1i32":   value.False.ToValue(),
		"-8f32 == -8i32": value.False.ToValue(),

		"1f32 == 1i16":   value.False.ToValue(),
		"-8f32 == -8i16": value.False.ToValue(),

		"1f32 == 1i8":   value.False.ToValue(),
		"-8f32 == -8i8": value.False.ToValue(),

		"1f32 == 1u64": value.False.ToValue(),

		"1f32 == 1u32": value.False.ToValue(),

		"1f32 == 1u16": value.False.ToValue(),

		"1f32 == 1u8": value.False.ToValue(),
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}

func TestVMSource_NotEqual(t *testing.T) {
	tests := simpleSourceTestTable{
		// String
		"'25' != '25'":   value.False.ToValue(),
		"'25' != '25.0'": value.True.ToValue(),
		"'25' != '7'":    value.True.ToValue(),

		"'7' != `7`": value.True.ToValue(),

		"'25' != 25.0": value.True.ToValue(),

		"'25' != 25bf": value.True.ToValue(),

		"'25' != 25f64": value.True.ToValue(),

		"'25' != 25f32": value.True.ToValue(),

		"'1' != 1i64": value.True.ToValue(),

		"'5' != 5i32": value.True.ToValue(),

		"'5' != 5i16": value.True.ToValue(),

		"'5' != 5i8": value.True.ToValue(),

		"'1' != 1u64": value.True.ToValue(),

		"'5' != 5u32": value.True.ToValue(),

		"'5' != 5u16": value.True.ToValue(),

		"'5' != 5u8": value.True.ToValue(),

		// Char
		"`2` != 25": value.True.ToValue(),

		"`2` != '2'": value.True.ToValue(),

		"`7` != `7`": value.False.ToValue(),
		"`b` != `b`": value.False.ToValue(),
		"`c` != `g`": value.True.ToValue(),
		"`7` != `8`": value.True.ToValue(),

		"`2` != 2.0": value.True.ToValue(),

		"`3` != 3bf": value.True.ToValue(),

		"`9` != 9f64": value.True.ToValue(),

		"`1` != 1f32": value.True.ToValue(),

		"`1` != 1i64": value.True.ToValue(),

		"`5` != 5i32": value.True.ToValue(),

		"`5` != 5i16": value.True.ToValue(),

		"`5` != 5i8": value.True.ToValue(),

		"`1` != 1u64": value.True.ToValue(),

		"`5` != 5u32": value.True.ToValue(),

		"`5` != 5u16": value.True.ToValue(),

		"`5` != 5u8": value.True.ToValue(),

		// Int
		"25 != 25":  value.False.ToValue(),
		"-25 != 25": value.True.ToValue(),
		"25 != -25": value.True.ToValue(),
		"25 != 28":  value.True.ToValue(),
		"28 != 25":  value.True.ToValue(),

		"25 != '25'": value.True.ToValue(),

		"7 != `7`": value.True.ToValue(),

		"25 != 25.0": value.True.ToValue(),

		"25 != 25bf": value.True.ToValue(),

		"25 != 25f64": value.True.ToValue(),

		"25 != 25f32": value.True.ToValue(),

		"1 != 1i64": value.True.ToValue(),

		"5 != 5i32": value.True.ToValue(),

		"5 != 5i16": value.True.ToValue(),

		"5 != 5i8": value.True.ToValue(),

		"1 != 1u64": value.True.ToValue(),

		"5 != 5u32": value.True.ToValue(),

		"5 != 5u16": value.True.ToValue(),

		"5 != 5u8": value.True.ToValue(),

		// Int64
		"25i64 != 25": value.True.ToValue(),

		"25i64 != '25'": value.True.ToValue(),

		"7i64 != `7`": value.True.ToValue(),

		"25i64 != 25.0": value.True.ToValue(),

		"25i64 != 25bf": value.True.ToValue(),

		"25i64 != 25f64": value.True.ToValue(),

		"25i64 != 25f32": value.True.ToValue(),

		"1i64 != 1i64":   value.False.ToValue(),
		"4i64 != -4i64":  value.True.ToValue(),
		"-8i64 != 8i64":  value.True.ToValue(),
		"-8i64 != -8i64": value.False.ToValue(),
		"91i64 != 27i64": value.True.ToValue(),

		"5i64 != 5i32": value.True.ToValue(),

		"5i64 != 5i16": value.True.ToValue(),

		"5i64 != 5i8": value.True.ToValue(),

		"1i64 != 1u64": value.True.ToValue(),

		"5i64 != 5u32": value.True.ToValue(),

		"5i64 != 5u16": value.True.ToValue(),

		"5i64 != 5u8": value.True.ToValue(),

		// Int32
		"25i32 != 25": value.True.ToValue(),

		"25i32 != '25'": value.True.ToValue(),

		"7i32 != `7`": value.True.ToValue(),

		"25i32 != 25.0": value.True.ToValue(),

		"25i32 != 25bf": value.True.ToValue(),

		"25i32 != 25f64": value.True.ToValue(),

		"25i32 != 25f32": value.True.ToValue(),

		"1i32 != 1i64": value.True.ToValue(),

		"5i32 != 5i32":  value.False.ToValue(),
		"4i32 != -4i32": value.True.ToValue(),
		"-8i32 != 8i32": value.True.ToValue(),
		"3i32 != 71i32": value.True.ToValue(),

		"5i32 != 5i16": value.True.ToValue(),

		"5i32 != 5i8": value.True.ToValue(),

		"1i32 != 1u64": value.True.ToValue(),

		"5i32 != 5u32": value.True.ToValue(),

		"5i32 != 5u16": value.True.ToValue(),

		"5i32 != 5u8": value.True.ToValue(),

		// Int16
		"25i16 != 25": value.True.ToValue(),

		"25i16 != '25'": value.True.ToValue(),

		"7i16 != `7`": value.True.ToValue(),

		"25i16 != 25.0": value.True.ToValue(),

		"25i16 != 25bf": value.True.ToValue(),

		"25i16 != 25f64": value.True.ToValue(),

		"25i16 != 25f32": value.True.ToValue(),

		"1i16 != 1i64": value.True.ToValue(),

		"5i16 != 5i32": value.True.ToValue(),

		"5i16 != 5i16":  value.False.ToValue(),
		"4i16 != -4i16": value.True.ToValue(),
		"-8i16 != 8i16": value.True.ToValue(),
		"3i16 != 71i16": value.True.ToValue(),

		"5i16 != 5i8": value.True.ToValue(),

		"1i16 != 1u64": value.True.ToValue(),

		"5i16 != 5u32": value.True.ToValue(),

		"5i16 != 5u16": value.True.ToValue(),

		"5i16 != 5u8": value.True.ToValue(),

		// Int8
		"25i8 != 25": value.True.ToValue(),

		"25i8 != '25'": value.True.ToValue(),

		"7i8 != `7`": value.True.ToValue(),

		"25i8 != 25.0": value.True.ToValue(),

		"25i8 != 25bf": value.True.ToValue(),

		"25i8 != 25f64": value.True.ToValue(),

		"25i8 != 25f32": value.True.ToValue(),

		"1i8 != 1i64": value.True.ToValue(),

		"5i8 != 5i32": value.True.ToValue(),

		"5i8 != 5i16": value.True.ToValue(),

		"5i8 != 5i8":  value.False.ToValue(),
		"4i8 != -4i8": value.True.ToValue(),
		"-8i8 != 8i8": value.True.ToValue(),
		"3i8 != 71i8": value.True.ToValue(),

		"1i8 != 1u64": value.True.ToValue(),

		"5i8 != 5u32": value.True.ToValue(),

		"5i8 != 5u16": value.True.ToValue(),

		"5i8 != 5u8": value.True.ToValue(),

		// UInt64
		"25u64 != 25": value.True.ToValue(),

		"25u64 != '25'": value.True.ToValue(),

		"7u64 != `7`": value.True.ToValue(),

		"25u64 != 25.0": value.True.ToValue(),

		"25u64 != 25bf": value.True.ToValue(),

		"25u64 != 25f64": value.True.ToValue(),

		"25u64 != 25f32": value.True.ToValue(),

		"1u64 != 1i64": value.True.ToValue(),

		"5u64 != 5i32": value.True.ToValue(),

		"5u64 != 5i16": value.True.ToValue(),

		"5u64 != 5i8": value.True.ToValue(),

		"1u64 != 1u64":   value.False.ToValue(),
		"91u64 != 27u64": value.True.ToValue(),

		"5u64 != 5u32": value.True.ToValue(),

		"5u64 != 5u16": value.True.ToValue(),

		"5u64 != 5u8": value.True.ToValue(),

		// UInt32
		"25u32 != 25": value.True.ToValue(),

		"25u32 != '25'": value.True.ToValue(),

		"7u32 != `7`": value.True.ToValue(),

		"25u32 != 25.0": value.True.ToValue(),

		"25u32 != 25bf": value.True.ToValue(),

		"25u32 != 25f64": value.True.ToValue(),

		"25u32 != 25f32": value.True.ToValue(),

		"1u32 != 1i64": value.True.ToValue(),

		"5u32 != 5i32": value.True.ToValue(),

		"5u32 != 5i16": value.True.ToValue(),

		"5u32 != 5i8": value.True.ToValue(),

		"1u32 != 1u64": value.True.ToValue(),

		"5u32 != 5u32":  value.False.ToValue(),
		"3u32 != 71u32": value.True.ToValue(),

		"5u32 != 5u16": value.True.ToValue(),

		"5u32 != 5u8": value.True.ToValue(),

		// UInt16
		"25u16 != 25": value.True.ToValue(),

		"25u16 != '25'": value.True.ToValue(),

		"7u16 != `7`": value.True.ToValue(),

		"25u16 != 25.0": value.True.ToValue(),

		"25u16 != 25bf": value.True.ToValue(),

		"25u16 != 25f64": value.True.ToValue(),

		"25u16 != 25f32": value.True.ToValue(),

		"1u16 != 1i64": value.True.ToValue(),

		"5u16 != 5i32": value.True.ToValue(),

		"5u16 != 5i16": value.True.ToValue(),

		"5u16 != 5i8": value.True.ToValue(),

		"1u16 != 1u64": value.True.ToValue(),

		"5u16 != 5u32": value.True.ToValue(),

		"53000u16 != 32767u16": value.True.ToValue(),
		"5u16 != 5u16":         value.False.ToValue(),
		"3u16 != 71u16":        value.True.ToValue(),

		"5u16 != 5u8": value.True.ToValue(),

		// UInt8
		"25u8 != 25": value.True.ToValue(),

		"25u8 != '25'": value.True.ToValue(),

		"7u8 != `7`": value.True.ToValue(),

		"25u8 != 25.0": value.True.ToValue(),

		"25u8 != 25bf": value.True.ToValue(),

		"25u8 != 25f64": value.True.ToValue(),

		"25u8 != 25f32": value.True.ToValue(),

		"1u8 != 1i64": value.True.ToValue(),

		"5u8 != 5i32": value.True.ToValue(),

		"5u8 != 5i16": value.True.ToValue(),

		"5u8 != 5i8": value.True.ToValue(),

		"1u8 != 1u64": value.True.ToValue(),

		"5u8 != 5u32": value.True.ToValue(),

		"5u8 != 5u16": value.True.ToValue(),

		"5u8 != 5u8":  value.False.ToValue(),
		"3u8 != 71u8": value.True.ToValue(),

		// Float
		"-73.0 != 73.0": value.True.ToValue(),
		"73.0 != -73.0": value.True.ToValue(),
		"25.0 != 25.0":  value.False.ToValue(),
		"1.0 != 1.2":    value.True.ToValue(),
		"1.2 != 1.0":    value.True.ToValue(),
		"78.5 != 78.5":  value.False.ToValue(),

		"8.25 != '8.25'": value.True.ToValue(),

		"4.0 != `4`": value.True.ToValue(),

		"25.0 != 25": value.True.ToValue(),

		"25.0 != 25bf":   value.True.ToValue(),
		"15.5 != 15.5bf": value.True.ToValue(),

		"25.0 != 25f64":     value.True.ToValue(),
		"15.26 != 15.26f64": value.True.ToValue(),

		"25.0 != 25f32":   value.True.ToValue(),
		"15.5 != 15.5f32": value.True.ToValue(),

		"1.0 != 1i64":   value.True.ToValue(),
		"-8.0 != -8i64": value.True.ToValue(),

		"1.0 != 1i32":   value.True.ToValue(),
		"-8.0 != -8i32": value.True.ToValue(),

		"1.0 != 1i16":   value.True.ToValue(),
		"-8.0 != -8i16": value.True.ToValue(),

		"1.0 != 1i8":   value.True.ToValue(),
		"-8.0 != -8i8": value.True.ToValue(),

		"1.0 != 1u64": value.True.ToValue(),

		"1.0 != 1u32": value.True.ToValue(),

		"1.0 != 1u16": value.True.ToValue(),

		"1.0 != 1u8": value.True.ToValue(),

		// Float64
		"25f64 != 25.0":   value.True.ToValue(),
		"78.5f64 != 78.5": value.True.ToValue(),

		"8.25f64 != '8.25'": value.True.ToValue(),

		"4f64 != `4`": value.True.ToValue(),

		"25f64 != 25": value.True.ToValue(),

		"25f64 != 25bf":     value.True.ToValue(),
		"15.5f64 != 15.5bf": value.True.ToValue(),

		"-73f64 != 73f64":      value.True.ToValue(),
		"73f64 != -73f64":      value.True.ToValue(),
		"25f64 != 25f64":       value.False.ToValue(),
		"1f64 != 1.2f64":       value.True.ToValue(),
		"15.26f64 != 15.26f64": value.False.ToValue(),

		"25f64 != 25f32":     value.True.ToValue(),
		"15.5f64 != 15.5f32": value.True.ToValue(),

		"1f64 != 1i64":   value.True.ToValue(),
		"-8f64 != -8i64": value.True.ToValue(),

		"1f64 != 1i32":   value.True.ToValue(),
		"-8f64 != -8i32": value.True.ToValue(),

		"1f64 != 1i16":   value.True.ToValue(),
		"-8f64 != -8i16": value.True.ToValue(),

		"1f64 != 1i8":   value.True.ToValue(),
		"-8f64 != -8i8": value.True.ToValue(),

		"1f64 != 1u64": value.True.ToValue(),

		"1f64 != 1u32": value.True.ToValue(),

		"1f64 != 1u16": value.True.ToValue(),

		"1f64 != 1u8": value.True.ToValue(),

		// Float32
		"25f32 != 25.0":   value.True.ToValue(),
		"78.5f32 != 78.5": value.True.ToValue(),

		"8.25f32 != '8.25'": value.True.ToValue(),

		"4f32 != `4`": value.True.ToValue(),

		"25f32 != 25": value.True.ToValue(),

		"25f32 != 25bf":     value.True.ToValue(),
		"15.5f32 != 15.5bf": value.True.ToValue(),

		"25f32 != 25f64":     value.True.ToValue(),
		"15.5f32 != 15.5f64": value.True.ToValue(),

		"-73f32 != 73f32":    value.True.ToValue(),
		"73f32 != -73f32":    value.True.ToValue(),
		"25f32 != 25f32":     value.False.ToValue(),
		"1f32 != 1.2f32":     value.True.ToValue(),
		"15.5f32 != 15.5f32": value.False.ToValue(),

		"1f32 != 1i64":   value.True.ToValue(),
		"-8f32 != -8i64": value.True.ToValue(),

		"1f32 != 1i32":   value.True.ToValue(),
		"-8f32 != -8i32": value.True.ToValue(),

		"1f32 != 1i16":   value.True.ToValue(),
		"-8f32 != -8i16": value.True.ToValue(),

		"1f32 != 1i8":   value.True.ToValue(),
		"-8f32 != -8i8": value.True.ToValue(),

		"1f32 != 1u64": value.True.ToValue(),

		"1f32 != 1u32": value.True.ToValue(),

		"1f32 != 1u16": value.True.ToValue(),

		"1f32 != 1u8": value.True.ToValue(),
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}
