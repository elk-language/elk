package vm_test

import (
	"testing"

	"github.com/elk-language/elk/value"
)

func TestVMSource_GreaterThan(t *testing.T) {
	tests := sourceTestTable{
		// String
		"'25' > '25'": {
			source:       "'25' > '25'",
			wantStackTop: value.False,
		},
		"'7' > '10'": {
			source:       "'7' > '10'",
			wantStackTop: value.True,
		},
		"'10' > '7'": {
			source:       "'10' > '7'",
			wantStackTop: value.False,
		},
		"'25' > '22'": {
			source:       "'25' > '22'",
			wantStackTop: value.True,
		},
		"'22' > '25'": {
			source:       "'22' > '25'",
			wantStackTop: value.False,
		},
		"'foo' > 'foo'": {
			source:       "'foo' > 'foo'",
			wantStackTop: value.False,
		},
		"'foo' > 'foa'": {
			source:       "'foo' > 'foa'",
			wantStackTop: value.True,
		},
		"'foa' > 'foo'": {
			source:       "'foa' > 'foo'",
			wantStackTop: value.False,
		},
		"'foo' > 'foo bar'": {
			source:       "'foo' > 'foo bar'",
			wantStackTop: value.False,
		},
		"'foo bar' > 'foo'": {
			source:       "'foo bar' > 'foo'",
			wantStackTop: value.True,
		},

		"'2' > `2`": {
			source:       "'2' > `2`",
			wantStackTop: value.False,
		},
		"'72' > `7`": {
			source:       "'72' > `7`",
			wantStackTop: value.True,
		},
		"'8' > `7`": {
			source:       "'8' > `7`",
			wantStackTop: value.True,
		},
		"'7' > `8`": {
			source:       "'7' > `8`",
			wantStackTop: value.False,
		},
		"'ba' > `b`": {
			source:       "'ba' > `b`",
			wantStackTop: value.True,
		},
		"'b' > `a`": {
			source:       "'b' > `a`",
			wantStackTop: value.True,
		},
		"'a' > `b`": {
			source:       "'a' > `b`",
			wantStackTop: value.False,
		},

		"'2' > 2.0": {
			source: "'2' > 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::String`",
			),
		},

		"'28' > 25.2bf": {
			source: "'28' > 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::String`",
			),
		},

		"'28.8' > 12.9f64": {
			source: "'28.8' > 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::String`",
			),
		},

		"'28.8' > 12.9f32": {
			source: "'28.8' > 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::String`",
			),
		},

		"'93' > 19i64": {
			source: "'93' > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::String`",
			),
		},

		"'93' > 19i32": {
			source: "'93' > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::String`",
			),
		},

		"'93' > 19i16": {
			source: "'93' > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::String`",
			),
		},

		"'93' > 19i8": {
			source: "'93' > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::String`",
			),
		},

		"'93' > 19u64": {
			source: "'93' > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::String`",
			),
		},

		"'93' > 19u32": {
			source: "'93' > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::String`",
			),
		},

		"'93' > 19u16": {
			source: "'93' > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::String`",
			),
		},

		"'93' > 19u8": {
			source: "'93' > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::String`",
			),
		},

		// Char
		"`2` > `2`": {
			source:       "`2` > `2`",
			wantStackTop: value.False,
		},
		"`8` > `7`": {
			source:       "`8` > `7`",
			wantStackTop: value.True,
		},
		"`7` > `8`": {
			source:       "`7` > `8`",
			wantStackTop: value.False,
		},
		"`b` > `a`": {
			source:       "`b` > `a`",
			wantStackTop: value.True,
		},
		"`a` > `b`": {
			source:       "`a` > `b`",
			wantStackTop: value.False,
		},

		"`2` > '2'": {
			source:       "`2` > '2'",
			wantStackTop: value.False,
		},
		"`7` > '72'": {
			source:       "`7` > '72'",
			wantStackTop: value.False,
		},
		"`8` > '7'": {
			source:       "`8` > '7'",
			wantStackTop: value.True,
		},
		"`7` > '8'": {
			source:       "`7` > '8'",
			wantStackTop: value.False,
		},
		"`b` > 'a'": {
			source:       "`b` > 'a'",
			wantStackTop: value.True,
		},
		"`b` > 'ba'": {
			source:       "`b` > 'ba'",
			wantStackTop: value.False,
		},
		"`a` > 'b'": {
			source:       "`a` > 'b'",
			wantStackTop: value.False,
		},

		"`2` > 2.0": {
			source: "`2` > 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Char`",
			),
		},
		"`i` > 25.2bf": {
			source: "`i` > 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Char`",
			),
		},
		"`f` > 12.9f64": {
			source: "`f` > 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Char`",
			),
		},
		"`0` > 12.9f32": {
			source: "`0` > 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Char`",
			),
		},
		"`9` > 19i64": {
			source: "`9` > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Char`",
			),
		},
		"`u` > 19i32": {
			source: "`u` > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Char`",
			),
		},
		"`4` > 19i16": {
			source: "`4` > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Char`",
			),
		},
		"`6` > 19i8": {
			source: "`6` > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Char`",
			),
		},
		"`9` > 19u64": {
			source: "`9` > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Char`",
			),
		},
		"`u` > 19u32": {
			source: "`u` > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Char`",
			),
		},
		"`4` > 19u16": {
			source: "`4` > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Char`",
			),
		},
		"`6` > 19u8": {
			source: "`6` > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Char`",
			),
		},

		// Int
		"25 > 25": {
			source:       "25 > 25",
			wantStackTop: value.False,
		},
		"25 > -25": {
			source:       "25 > -25",
			wantStackTop: value.True,
		},
		"-25 > 25": {
			source:       "-25 > 25",
			wantStackTop: value.False,
		},
		"13 > 7": {
			source:       "13 > 7",
			wantStackTop: value.True,
		},
		"7 > 13": {
			source:       "7 > 13",
			wantStackTop: value.False,
		},

		"25 > 25.0": {
			source:       "25 > 25.0",
			wantStackTop: value.False,
		},
		"25 > -25.0": {
			source:       "25 > -25.0",
			wantStackTop: value.True,
		},
		"-25 > 25.0": {
			source:       "-25 > 25.0",
			wantStackTop: value.False,
		},
		"13 > 7.0": {
			source:       "13 > 7.0",
			wantStackTop: value.True,
		},
		"7 > 13.0": {
			source:       "7 > 13.0",
			wantStackTop: value.False,
		},
		"7 > 7.5": {
			source:       "7 > 7.5",
			wantStackTop: value.False,
		},
		"7 > 6.9": {
			source:       "7 > 6.9",
			wantStackTop: value.True,
		},

		"25 > 25bf": {
			source:       "25 > 25bf",
			wantStackTop: value.False,
		},
		"25 > -25bf": {
			source:       "25 > -25bf",
			wantStackTop: value.True,
		},
		"-25 > 25bf": {
			source:       "-25 > 25bf",
			wantStackTop: value.False,
		},
		"13 > 7bf": {
			source:       "13 > 7bf",
			wantStackTop: value.True,
		},
		"7 > 13bf": {
			source:       "7 > 13bf",
			wantStackTop: value.False,
		},
		"7 > 7.5bf": {
			source:       "7 > 7.5bf",
			wantStackTop: value.False,
		},
		"7 > 6.9bf": {
			source:       "7 > 6.9bf",
			wantStackTop: value.True,
		},

		"6 > 19f64": {
			source: "6 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int`",
			),
		},
		"6 > 19f32": {
			source: "6 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int`",
			),
		},
		"6 > 19i64": {
			source: "6 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int`",
			),
		},
		"6 > 19i32": {
			source: "6 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int`",
			),
		},
		"6 > 19i16": {
			source: "6 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int`",
			),
		},
		"6 > 19i8": {
			source: "6 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int`",
			),
		},
		"6 > 19u64": {
			source: "6 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int`",
			),
		},
		"6 > 19u32": {
			source: "6 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int`",
			),
		},
		"6 > 19u16": {
			source: "6 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int`",
			),
		},
		"6 > 19u8": {
			source: "6 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int`",
			),
		},

		// Float
		"25.0 > 25.0": {
			source:       "25.0 > 25.0",
			wantStackTop: value.False,
		},
		"25.0 > -25.0": {
			source:       "25.0 > -25.0",
			wantStackTop: value.True,
		},
		"-25.0 > 25.0": {
			source:       "-25.0 > 25.0",
			wantStackTop: value.False,
		},
		"13.0 > 7.0": {
			source:       "13.0 > 7.0",
			wantStackTop: value.True,
		},
		"7.0 > 13.0": {
			source:       "7.0 > 13.0",
			wantStackTop: value.False,
		},
		"7.0 > 7.5": {
			source:       "7.0 > 7.5",
			wantStackTop: value.False,
		},
		"7.5 > 7.0": {
			source:       "7.5 > 7.0",
			wantStackTop: value.True,
		},
		"7.0 > 6.9": {
			source:       "7.0 > 6.9",
			wantStackTop: value.True,
		},

		"25.0 > 25": {
			source:       "25.0 > 25",
			wantStackTop: value.False,
		},
		"25.0 > -25": {
			source:       "25.0 > -25",
			wantStackTop: value.True,
		},
		"-25.0 > 25": {
			source:       "-25.0 > 25",
			wantStackTop: value.False,
		},
		"13.0 > 7": {
			source:       "13.0 > 7",
			wantStackTop: value.True,
		},
		"7.0 > 13": {
			source:       "7.0 > 13",
			wantStackTop: value.False,
		},
		"7.5 > 7": {
			source:       "7.5 > 7",
			wantStackTop: value.True,
		},

		"25.0 > 25bf": {
			source:       "25.0 > 25bf",
			wantStackTop: value.False,
		},
		"25.0 > -25bf": {
			source:       "25.0 > -25bf",
			wantStackTop: value.True,
		},
		"-25.0 > 25bf": {
			source:       "-25.0 > 25bf",
			wantStackTop: value.False,
		},
		"13.0 > 7bf": {
			source:       "13.0 > 7bf",
			wantStackTop: value.True,
		},
		"7.0 > 13bf": {
			source:       "7.0 > 13bf",
			wantStackTop: value.False,
		},
		"7.0 > 7.5bf": {
			source:       "7.0 > 7.5bf",
			wantStackTop: value.False,
		},
		"7.5 > 7bf": {
			source:       "7.5 > 7bf",
			wantStackTop: value.True,
		},
		"7.0 > 6.9bf": {
			source:       "7.0 > 6.9bf",
			wantStackTop: value.True,
		},

		"6.0 > 19f64": {
			source: "6.0 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 > 19f32": {
			source: "6.0 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 > 19i64": {
			source: "6.0 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 > 19i32": {
			source: "6.0 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 > 19i16": {
			source: "6.0 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 > 19i8": {
			source: "6.0 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 > 19u64": {
			source: "6.0 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 > 19u32": {
			source: "6.0 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 > 19u16": {
			source: "6.0 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 > 19u8": {
			source: "6.0 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float`",
			),
		},

		// BigFloat
		"25bf > 25.0": {
			source:       "25bf > 25.0",
			wantStackTop: value.False,
		},
		"25bf > -25.0": {
			source:       "25bf > -25.0",
			wantStackTop: value.True,
		},
		"-25bf > 25.0": {
			source:       "-25bf > 25.0",
			wantStackTop: value.False,
		},
		"13bf > 7.0": {
			source:       "13bf > 7.0",
			wantStackTop: value.True,
		},
		"7bf > 13.0": {
			source:       "7bf > 13.0",
			wantStackTop: value.False,
		},
		"7bf > 7.5": {
			source:       "7bf > 7.5",
			wantStackTop: value.False,
		},
		"7.5bf > 7.0": {
			source:       "7.5bf > 7.0",
			wantStackTop: value.True,
		},
		"7bf > 6.9": {
			source:       "7bf > 6.9",
			wantStackTop: value.True,
		},

		"25bf > 25": {
			source:       "25bf > 25",
			wantStackTop: value.False,
		},
		"25bf > -25": {
			source:       "25bf > -25",
			wantStackTop: value.True,
		},
		"-25bf > 25": {
			source:       "-25bf > 25",
			wantStackTop: value.False,
		},
		"13bf > 7": {
			source:       "13bf > 7",
			wantStackTop: value.True,
		},
		"7bf > 13": {
			source:       "7bf > 13",
			wantStackTop: value.False,
		},
		"7.5bf > 7": {
			source:       "7.5bf > 7",
			wantStackTop: value.True,
		},

		"25bf > 25bf": {
			source:       "25bf > 25bf",
			wantStackTop: value.False,
		},
		"25bf > -25bf": {
			source:       "25bf > -25bf",
			wantStackTop: value.True,
		},
		"-25bf > 25bf": {
			source:       "-25bf > 25bf",
			wantStackTop: value.False,
		},
		"13bf > 7bf": {
			source:       "13bf > 7bf",
			wantStackTop: value.True,
		},
		"7bf > 13bf": {
			source:       "7bf > 13bf",
			wantStackTop: value.False,
		},
		"7bf > 7.5bf": {
			source:       "7bf > 7.5bf",
			wantStackTop: value.False,
		},
		"7.5bf > 7bf": {
			source:       "7.5bf > 7bf",
			wantStackTop: value.True,
		},
		"7bf > 6.9bf": {
			source:       "7bf > 6.9bf",
			wantStackTop: value.True,
		},

		"6bf > 19f64": {
			source: "6bf > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19f32": {
			source: "6bf > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19i64": {
			source: "6bf > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19i32": {
			source: "6bf > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19i16": {
			source: "6bf > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19i8": {
			source: "6bf > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19u64": {
			source: "6bf > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19u32": {
			source: "6bf > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19u16": {
			source: "6bf > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19u8": {
			source: "6bf > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::BigFloat`",
			),
		},

		// Float64
		"25f64 > 25f64": {
			source:       "25f64 > 25f64",
			wantStackTop: value.False,
		},
		"25f64 > -25f64": {
			source:       "25f64 > -25f64",
			wantStackTop: value.True,
		},
		"-25f64 > 25f64": {
			source:       "-25f64 > 25f64",
			wantStackTop: value.False,
		},
		"13f64 > 7f64": {
			source:       "13f64 > 7f64",
			wantStackTop: value.True,
		},
		"7f64 > 13f64": {
			source:       "7f64 > 13f64",
			wantStackTop: value.False,
		},
		"7f64 > 7.5f64": {
			source:       "7f64 > 7.5f64",
			wantStackTop: value.False,
		},
		"7.5f64 > 7f64": {
			source:       "7.5f64 > 7f64",
			wantStackTop: value.True,
		},
		"7f64 > 6.9f64": {
			source:       "7f64 > 6.9f64",
			wantStackTop: value.True,
		},

		"6f64 > 19.0": {
			source: "6f64 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Float64`",
			),
		},

		"6f64 > 19": {
			source: "6f64 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19bf": {
			source: "6f64 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19f32": {
			source: "6f64 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19i64": {
			source: "6f64 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19i32": {
			source: "6f64 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19i16": {
			source: "6f64 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19i8": {
			source: "6f64 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19u64": {
			source: "6f64 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19u32": {
			source: "6f64 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19u16": {
			source: "6f64 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19u8": {
			source: "6f64 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float64`",
			),
		},

		// Float32
		"25f32 > 25f32": {
			source:       "25f32 > 25f32",
			wantStackTop: value.False,
		},
		"25f32 > -25f32": {
			source:       "25f32 > -25f32",
			wantStackTop: value.True,
		},
		"-25f32 > 25f32": {
			source:       "-25f32 > 25f32",
			wantStackTop: value.False,
		},
		"13f32 > 7f32": {
			source:       "13f32 > 7f32",
			wantStackTop: value.True,
		},
		"7f32 > 13f32": {
			source:       "7f32 > 13f32",
			wantStackTop: value.False,
		},
		"7f32 > 7.5f32": {
			source:       "7f32 > 7.5f32",
			wantStackTop: value.False,
		},
		"7.5f32 > 7f32": {
			source:       "7.5f32 > 7f32",
			wantStackTop: value.True,
		},
		"7f32 > 6.9f32": {
			source:       "7f32 > 6.9f32",
			wantStackTop: value.True,
		},

		"6f32 > 19.0": {
			source: "6f32 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Float32`",
			),
		},

		"6f32 > 19": {
			source: "6f32 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19bf": {
			source: "6f32 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19f64": {
			source: "6f32 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19i64": {
			source: "6f32 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19i32": {
			source: "6f32 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19i16": {
			source: "6f32 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19i8": {
			source: "6f32 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19u64": {
			source: "6f32 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19u32": {
			source: "6f32 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19u16": {
			source: "6f32 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19u8": {
			source: "6f32 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float32`",
			),
		},

		// Int64
		"25i64 > 25i64": {
			source:       "25i64 > 25i64",
			wantStackTop: value.False,
		},
		"25i64 > -25i64": {
			source:       "25i64 > -25i64",
			wantStackTop: value.True,
		},
		"-25i64 > 25i64": {
			source:       "-25i64 > 25i64",
			wantStackTop: value.False,
		},
		"13i64 > 7i64": {
			source:       "13i64 > 7i64",
			wantStackTop: value.True,
		},
		"7i64 > 13i64": {
			source:       "7i64 > 13i64",
			wantStackTop: value.False,
		},

		"6i64 > 19": {
			source: "6i64 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19.0": {
			source: "6i64 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19bf": {
			source: "6i64 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19f64": {
			source: "6i64 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19f32": {
			source: "6i64 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19i32": {
			source: "6i64 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19i16": {
			source: "6i64 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19i8": {
			source: "6i64 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19u64": {
			source: "6i64 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19u32": {
			source: "6i64 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19u16": {
			source: "6i64 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19u8": {
			source: "6i64 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int64`",
			),
		},

		// Int32
		"25i32 > 25i32": {
			source:       "25i32 > 25i32",
			wantStackTop: value.False,
		},
		"25i32 > -25i32": {
			source:       "25i32 > -25i32",
			wantStackTop: value.True,
		},
		"-25i32 > 25i32": {
			source:       "-25i32 > 25i32",
			wantStackTop: value.False,
		},
		"13i32 > 7i32": {
			source:       "13i32 > 7i32",
			wantStackTop: value.True,
		},
		"7i32 > 13i32": {
			source:       "7i32 > 13i32",
			wantStackTop: value.False,
		},

		"6i32 > 19": {
			source: "6i32 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19.0": {
			source: "6i32 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19bf": {
			source: "6i32 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19f64": {
			source: "6i32 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19f32": {
			source: "6i32 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19i64": {
			source: "6i32 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19i16": {
			source: "6i32 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19i8": {
			source: "6i32 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19u64": {
			source: "6i32 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19u32": {
			source: "6i32 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19u16": {
			source: "6i32 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19u8": {
			source: "6i32 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int32`",
			),
		},

		// Int16
		"25i16 > 25i16": {
			source:       "25i16 > 25i16",
			wantStackTop: value.False,
		},
		"25i16 > -25i16": {
			source:       "25i16 > -25i16",
			wantStackTop: value.True,
		},
		"-25i16 > 25i16": {
			source:       "-25i16 > 25i16",
			wantStackTop: value.False,
		},
		"13i16 > 7i16": {
			source:       "13i16 > 7i16",
			wantStackTop: value.True,
		},
		"7i16 > 13i16": {
			source:       "7i16 > 13i16",
			wantStackTop: value.False,
		},

		"6i16 > 19": {
			source: "6i16 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19.0": {
			source: "6i16 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19bf": {
			source: "6i16 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19f64": {
			source: "6i16 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19f32": {
			source: "6i16 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19i64": {
			source: "6i16 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19i32": {
			source: "6i16 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19i8": {
			source: "6i16 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19u64": {
			source: "6i16 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19u32": {
			source: "6i16 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19u16": {
			source: "6i16 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19u8": {
			source: "6i16 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int16`",
			),
		},

		// Int8
		"25i8 > 25i8": {
			source:       "25i8 > 25i8",
			wantStackTop: value.False,
		},
		"25i8 > -25i8": {
			source:       "25i8 > -25i8",
			wantStackTop: value.True,
		},
		"-25i8 > 25i8": {
			source:       "-25i8 > 25i8",
			wantStackTop: value.False,
		},
		"13i8 > 7i8": {
			source:       "13i8 > 7i8",
			wantStackTop: value.True,
		},
		"7i8 > 13i8": {
			source:       "7i8 > 13i8",
			wantStackTop: value.False,
		},

		"6i8 > 19": {
			source: "6i8 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19.0": {
			source: "6i8 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19bf": {
			source: "6i8 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19f64": {
			source: "6i8 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19f32": {
			source: "6i8 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19i64": {
			source: "6i8 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19i32": {
			source: "6i8 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19i16": {
			source: "6i8 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19u64": {
			source: "6i8 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19u32": {
			source: "6i8 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19u16": {
			source: "6i8 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19u8": {
			source: "6i8 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int8`",
			),
		},

		// UInt64
		"25u64 > 25u64": {
			source:       "25u64 > 25u64",
			wantStackTop: value.False,
		},
		"13u64 > 7u64": {
			source:       "13u64 > 7u64",
			wantStackTop: value.True,
		},
		"7u64 > 13u64": {
			source:       "7u64 > 13u64",
			wantStackTop: value.False,
		},

		"6u64 > 19": {
			source: "6u64 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19.0": {
			source: "6u64 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19bf": {
			source: "6u64 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19f64": {
			source: "6u64 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19f32": {
			source: "6u64 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19i64": {
			source: "6u64 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19i32": {
			source: "6u64 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19i16": {
			source: "6u64 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19i8": {
			source: "6u64 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19u32": {
			source: "6u64 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19u16": {
			source: "6u64 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19u8": {
			source: "6u64 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt64`",
			),
		},

		// UInt32
		"25u32 > 25u32": {
			source:       "25u32 > 25u32",
			wantStackTop: value.False,
		},
		"13u32 > 7u32": {
			source:       "13u32 > 7u32",
			wantStackTop: value.True,
		},
		"7u32 > 13u32": {
			source:       "7u32 > 13u32",
			wantStackTop: value.False,
		},

		"6u32 > 19": {
			source: "6u32 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19.0": {
			source: "6u32 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19bf": {
			source: "6u32 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19f64": {
			source: "6u32 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19f32": {
			source: "6u32 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19i64": {
			source: "6u32 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19i32": {
			source: "6u32 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19i16": {
			source: "6u32 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19i8": {
			source: "6u32 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19u64": {
			source: "6u32 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19u16": {
			source: "6u32 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19u8": {
			source: "6u32 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt32`",
			),
		},

		// Int16
		"25u16 > 25u16": {
			source:       "25u16 > 25u16",
			wantStackTop: value.False,
		},
		"13u16 > 7u16": {
			source:       "13u16 > 7u16",
			wantStackTop: value.True,
		},
		"7u16 > 13u16": {
			source:       "7u16 > 13u16",
			wantStackTop: value.False,
		},

		"6u16 > 19": {
			source: "6u16 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19.0": {
			source: "6u16 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19bf": {
			source: "6u16 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19f64": {
			source: "6u16 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19f32": {
			source: "6u16 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19i64": {
			source: "6u16 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19i32": {
			source: "6u16 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19i16": {
			source: "6u16 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19i8": {
			source: "6u16 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19u64": {
			source: "6u16 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19u32": {
			source: "6u16 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19u8": {
			source: "6u16 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt16`",
			),
		},

		// Int8
		"25u8 > 25u8": {
			source:       "25u8 > 25u8",
			wantStackTop: value.False,
		},
		"13u8 > 7u8": {
			source:       "13u8 > 7u8",
			wantStackTop: value.True,
		},
		"7u8 > 13u8": {
			source:       "7u8 > 13u8",
			wantStackTop: value.False,
		},

		"6u8 > 19": {
			source: "6u8 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19.0": {
			source: "6u8 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19bf": {
			source: "6u8 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19f64": {
			source: "6u8 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19f32": {
			source: "6u8 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19i64": {
			source: "6u8 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19i32": {
			source: "6u8 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19i16": {
			source: "6u8 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19i8": {
			source: "6u8 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19u64": {
			source: "6u8 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19u32": {
			source: "6u8 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19u16": {
			source: "6u8 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt8`",
			),
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
			wantStackTop: value.True,
		},
		"'7' >= '10'": {
			source:       "'7' >= '10'",
			wantStackTop: value.True,
		},
		"'10' >= '7'": {
			source:       "'10' >= '7'",
			wantStackTop: value.False,
		},
		"'25' >= '22'": {
			source:       "'25' >= '22'",
			wantStackTop: value.True,
		},
		"'22' >= '25'": {
			source:       "'22' >= '25'",
			wantStackTop: value.False,
		},
		"'foo' >= 'foo'": {
			source:       "'foo' >= 'foo'",
			wantStackTop: value.True,
		},
		"'foo' >= 'foa'": {
			source:       "'foo' >= 'foa'",
			wantStackTop: value.True,
		},
		"'foa' >= 'foo'": {
			source:       "'foa' >= 'foo'",
			wantStackTop: value.False,
		},
		"'foo' >= 'foo bar'": {
			source:       "'foo' >= 'foo bar'",
			wantStackTop: value.False,
		},
		"'foo bar' >= 'foo'": {
			source:       "'foo bar' >= 'foo'",
			wantStackTop: value.True,
		},

		"'2' >= `2`": {
			source:       "'2' >= `2`",
			wantStackTop: value.True,
		},
		"'72' >= `7`": {
			source:       "'72' >= `7`",
			wantStackTop: value.True,
		},
		"'8' >= `7`": {
			source:       "'8' >= `7`",
			wantStackTop: value.True,
		},
		"'7' >= `8`": {
			source:       "'7' >= `8`",
			wantStackTop: value.False,
		},
		"'ba' >= `b`": {
			source:       "'ba' >= `b`",
			wantStackTop: value.True,
		},
		"'b' >= `a`": {
			source:       "'b' >= `a`",
			wantStackTop: value.True,
		},
		"'a' >= `b`": {
			source:       "'a' >= `b`",
			wantStackTop: value.False,
		},

		"'2' >= 2.0": {
			source: "'2' >= 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::String`",
			),
		},

		"'28' >= 25.2bf": {
			source: "'28' >= 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::String`",
			),
		},

		"'28.8' >= 12.9f64": {
			source: "'28.8' >= 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::String`",
			),
		},

		"'28.8' >= 12.9f32": {
			source: "'28.8' >= 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::String`",
			),
		},

		"'93' >= 19i64": {
			source: "'93' >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::String`",
			),
		},

		"'93' >= 19i32": {
			source: "'93' >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::String`",
			),
		},

		"'93' >= 19i16": {
			source: "'93' >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::String`",
			),
		},

		"'93' >= 19i8": {
			source: "'93' >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::String`",
			),
		},

		"'93' >= 19u64": {
			source: "'93' >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::String`",
			),
		},

		"'93' >= 19u32": {
			source: "'93' >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::String`",
			),
		},

		"'93' >= 19u16": {
			source: "'93' >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::String`",
			),
		},

		"'93' >= 19u8": {
			source: "'93' >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::String`",
			),
		},

		// Char
		"`2` >= `2`": {
			source:       "`2` >= `2`",
			wantStackTop: value.True,
		},
		"`8` >= `7`": {
			source:       "`8` >= `7`",
			wantStackTop: value.True,
		},
		"`7` >= `8`": {
			source:       "`7` >= `8`",
			wantStackTop: value.False,
		},
		"`b` >= `a`": {
			source:       "`b` >= `a`",
			wantStackTop: value.True,
		},
		"`a` >= `b`": {
			source:       "`a` >= `b`",
			wantStackTop: value.False,
		},

		"`2` >= '2'": {
			source:       "`2` >= '2'",
			wantStackTop: value.True,
		},
		"`7` >= '72'": {
			source:       "`7` >= '72'",
			wantStackTop: value.False,
		},
		"`8` >= '7'": {
			source:       "`8` >= '7'",
			wantStackTop: value.True,
		},
		"`7` >= '8'": {
			source:       "`7` >= '8'",
			wantStackTop: value.False,
		},
		"`b` >= 'a'": {
			source:       "`b` >= 'a'",
			wantStackTop: value.True,
		},
		"`b` >= 'ba'": {
			source:       "`b` >= 'ba'",
			wantStackTop: value.False,
		},
		"`a` >= 'b'": {
			source:       "`a` >= 'b'",
			wantStackTop: value.False,
		},

		"`2` >= 2.0": {
			source: "`2` >= 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Char`",
			),
		},
		"`i` >= 25.2bf": {
			source: "`i` >= 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Char`",
			),
		},
		"`f` >= 12.9f64": {
			source: "`f` >= 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Char`",
			),
		},
		"`0` >= 12.9f32": {
			source: "`0` >= 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Char`",
			),
		},
		"`9` >= 19i64": {
			source: "`9` >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Char`",
			),
		},
		"`u` >= 19i32": {
			source: "`u` >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Char`",
			),
		},
		"`4` >= 19i16": {
			source: "`4` >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Char`",
			),
		},
		"`6` >= 19i8": {
			source: "`6` >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Char`",
			),
		},
		"`9` >= 19u64": {
			source: "`9` >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Char`",
			),
		},
		"`u` >= 19u32": {
			source: "`u` >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Char`",
			),
		},
		"`4` >= 19u16": {
			source: "`4` >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Char`",
			),
		},
		"`6` >= 19u8": {
			source: "`6` >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Char`",
			),
		},

		// Int
		"25 >= 25": {
			source:       "25 >= 25",
			wantStackTop: value.True,
		},
		"25 >= -25": {
			source:       "25 >= -25",
			wantStackTop: value.True,
		},
		"-25 >= 25": {
			source:       "-25 >= 25",
			wantStackTop: value.False,
		},
		"13 >= 7": {
			source:       "13 >= 7",
			wantStackTop: value.True,
		},
		"7 >= 13": {
			source:       "7 >= 13",
			wantStackTop: value.False,
		},

		"25 >= 25.0": {
			source:       "25 >= 25.0",
			wantStackTop: value.True,
		},
		"25 >= -25.0": {
			source:       "25 >= -25.0",
			wantStackTop: value.True,
		},
		"-25 >= 25.0": {
			source:       "-25 >= 25.0",
			wantStackTop: value.False,
		},
		"13 >= 7.0": {
			source:       "13 >= 7.0",
			wantStackTop: value.True,
		},
		"7 >= 13.0": {
			source:       "7 >= 13.0",
			wantStackTop: value.False,
		},
		"7 >= 7.5": {
			source:       "7 >= 7.5",
			wantStackTop: value.False,
		},
		"7 >= 6.9": {
			source:       "7 >= 6.9",
			wantStackTop: value.True,
		},

		"25 >= 25bf": {
			source:       "25 >= 25bf",
			wantStackTop: value.True,
		},
		"25 >= -25bf": {
			source:       "25 >= -25bf",
			wantStackTop: value.True,
		},
		"-25 >= 25bf": {
			source:       "-25 >= 25bf",
			wantStackTop: value.False,
		},
		"13 >= 7bf": {
			source:       "13 >= 7bf",
			wantStackTop: value.True,
		},
		"7 >= 13bf": {
			source:       "7 >= 13bf",
			wantStackTop: value.False,
		},
		"7 >= 7.5bf": {
			source:       "7 >= 7.5bf",
			wantStackTop: value.False,
		},
		"7 >= 6.9bf": {
			source:       "7 >= 6.9bf",
			wantStackTop: value.True,
		},

		"6 >= 19f64": {
			source: "6 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int`",
			),
		},
		"6 >= 19f32": {
			source: "6 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int`",
			),
		},
		"6 >= 19i64": {
			source: "6 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int`",
			),
		},
		"6 >= 19i32": {
			source: "6 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int`",
			),
		},
		"6 >= 19i16": {
			source: "6 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int`",
			),
		},
		"6 >= 19i8": {
			source: "6 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int`",
			),
		},
		"6 >= 19u64": {
			source: "6 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int`",
			),
		},
		"6 >= 19u32": {
			source: "6 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int`",
			),
		},
		"6 >= 19u16": {
			source: "6 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int`",
			),
		},
		"6 >= 19u8": {
			source: "6 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int`",
			),
		},

		// Float
		"25.0 >= 25.0": {
			source:       "25.0 >= 25.0",
			wantStackTop: value.True,
		},
		"25.0 >= -25.0": {
			source:       "25.0 >= -25.0",
			wantStackTop: value.True,
		},
		"-25.0 >= 25.0": {
			source:       "-25.0 >= 25.0",
			wantStackTop: value.False,
		},
		"13.0 >= 7.0": {
			source:       "13.0 >= 7.0",
			wantStackTop: value.True,
		},
		"7.0 >= 13.0": {
			source:       "7.0 >= 13.0",
			wantStackTop: value.False,
		},
		"7.0 >= 7.5": {
			source:       "7.0 >= 7.5",
			wantStackTop: value.False,
		},
		"7.5 >= 7.0": {
			source:       "7.5 >= 7.0",
			wantStackTop: value.True,
		},
		"7.0 >= 6.9": {
			source:       "7.0 >= 6.9",
			wantStackTop: value.True,
		},

		"25.0 >= 25": {
			source:       "25.0 >= 25",
			wantStackTop: value.True,
		},
		"25.0 >= -25": {
			source:       "25.0 >= -25",
			wantStackTop: value.True,
		},
		"-25.0 >= 25": {
			source:       "-25.0 >= 25",
			wantStackTop: value.False,
		},
		"13.0 >= 7": {
			source:       "13.0 >= 7",
			wantStackTop: value.True,
		},
		"7.0 >= 13": {
			source:       "7.0 >= 13",
			wantStackTop: value.False,
		},
		"7.5 >= 7": {
			source:       "7.5 >= 7",
			wantStackTop: value.True,
		},

		"25.0 >= 25bf": {
			source:       "25.0 >= 25bf",
			wantStackTop: value.True,
		},
		"25.0 >= -25bf": {
			source:       "25.0 >= -25bf",
			wantStackTop: value.True,
		},
		"-25.0 >= 25bf": {
			source:       "-25.0 >= 25bf",
			wantStackTop: value.False,
		},
		"13.0 >= 7bf": {
			source:       "13.0 >= 7bf",
			wantStackTop: value.True,
		},
		"7.0 >= 13bf": {
			source:       "7.0 >= 13bf",
			wantStackTop: value.False,
		},
		"7.0 >= 7.5bf": {
			source:       "7.0 >= 7.5bf",
			wantStackTop: value.False,
		},
		"7.5 >= 7bf": {
			source:       "7.5 >= 7bf",
			wantStackTop: value.True,
		},
		"7.0 >= 6.9bf": {
			source:       "7.0 >= 6.9bf",
			wantStackTop: value.True,
		},

		"6.0 >= 19f64": {
			source: "6.0 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19f32": {
			source: "6.0 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19i64": {
			source: "6.0 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19i32": {
			source: "6.0 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19i16": {
			source: "6.0 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19i8": {
			source: "6.0 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19u64": {
			source: "6.0 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19u32": {
			source: "6.0 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19u16": {
			source: "6.0 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19u8": {
			source: "6.0 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float`",
			),
		},

		// BigFloat
		"25bf >= 25.0": {
			source:       "25bf >= 25.0",
			wantStackTop: value.True,
		},
		"25bf >= -25.0": {
			source:       "25bf >= -25.0",
			wantStackTop: value.True,
		},
		"-25bf >= 25.0": {
			source:       "-25bf >= 25.0",
			wantStackTop: value.False,
		},
		"13bf >= 7.0": {
			source:       "13bf >= 7.0",
			wantStackTop: value.True,
		},
		"7bf >= 13.0": {
			source:       "7bf >= 13.0",
			wantStackTop: value.False,
		},
		"7bf >= 7.5": {
			source:       "7bf >= 7.5",
			wantStackTop: value.False,
		},
		"7.5bf >= 7.0": {
			source:       "7.5bf >= 7.0",
			wantStackTop: value.True,
		},
		"7bf >= 6.9": {
			source:       "7bf >= 6.9",
			wantStackTop: value.True,
		},

		"25bf >= 25": {
			source:       "25bf >= 25",
			wantStackTop: value.True,
		},
		"25bf >= -25": {
			source:       "25bf >= -25",
			wantStackTop: value.True,
		},
		"-25bf >= 25": {
			source:       "-25bf >= 25",
			wantStackTop: value.False,
		},
		"13bf >= 7": {
			source:       "13bf >= 7",
			wantStackTop: value.True,
		},
		"7bf >= 13": {
			source:       "7bf >= 13",
			wantStackTop: value.False,
		},
		"7.5bf >= 7": {
			source:       "7.5bf >= 7",
			wantStackTop: value.True,
		},

		"25bf >= 25bf": {
			source:       "25bf >= 25bf",
			wantStackTop: value.True,
		},
		"25bf >= -25bf": {
			source:       "25bf >= -25bf",
			wantStackTop: value.True,
		},
		"-25bf >= 25bf": {
			source:       "-25bf >= 25bf",
			wantStackTop: value.False,
		},
		"13bf >= 7bf": {
			source:       "13bf >= 7bf",
			wantStackTop: value.True,
		},
		"7bf >= 13bf": {
			source:       "7bf >= 13bf",
			wantStackTop: value.False,
		},
		"7bf >= 7.5bf": {
			source:       "7bf >= 7.5bf",
			wantStackTop: value.False,
		},
		"7.5bf >= 7bf": {
			source:       "7.5bf >= 7bf",
			wantStackTop: value.True,
		},
		"7bf >= 6.9bf": {
			source:       "7bf >= 6.9bf",
			wantStackTop: value.True,
		},

		"6bf >= 19f64": {
			source: "6bf >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19f32": {
			source: "6bf >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19i64": {
			source: "6bf >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19i32": {
			source: "6bf >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19i16": {
			source: "6bf >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19i8": {
			source: "6bf >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19u64": {
			source: "6bf >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19u32": {
			source: "6bf >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19u16": {
			source: "6bf >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19u8": {
			source: "6bf >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::BigFloat`",
			),
		},

		// Float64
		"25f64 >= 25f64": {
			source:       "25f64 >= 25f64",
			wantStackTop: value.True,
		},
		"25f64 >= -25f64": {
			source:       "25f64 >= -25f64",
			wantStackTop: value.True,
		},
		"-25f64 >= 25f64": {
			source:       "-25f64 >= 25f64",
			wantStackTop: value.False,
		},
		"13f64 >= 7f64": {
			source:       "13f64 >= 7f64",
			wantStackTop: value.True,
		},
		"7f64 >= 13f64": {
			source:       "7f64 >= 13f64",
			wantStackTop: value.False,
		},
		"7f64 >= 7.5f64": {
			source:       "7f64 >= 7.5f64",
			wantStackTop: value.False,
		},
		"7.5f64 >= 7f64": {
			source:       "7.5f64 >= 7f64",
			wantStackTop: value.True,
		},
		"7f64 >= 6.9f64": {
			source:       "7f64 >= 6.9f64",
			wantStackTop: value.True,
		},

		"6f64 >= 19.0": {
			source: "6f64 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Float64`",
			),
		},

		"6f64 >= 19": {
			source: "6f64 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19bf": {
			source: "6f64 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19f32": {
			source: "6f64 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19i64": {
			source: "6f64 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19i32": {
			source: "6f64 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19i16": {
			source: "6f64 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19i8": {
			source: "6f64 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19u64": {
			source: "6f64 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19u32": {
			source: "6f64 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19u16": {
			source: "6f64 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19u8": {
			source: "6f64 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float64`",
			),
		},

		// Float32
		"25f32 >= 25f32": {
			source:       "25f32 >= 25f32",
			wantStackTop: value.True,
		},
		"25f32 >= -25f32": {
			source:       "25f32 >= -25f32",
			wantStackTop: value.True,
		},
		"-25f32 >= 25f32": {
			source:       "-25f32 >= 25f32",
			wantStackTop: value.False,
		},
		"13f32 >= 7f32": {
			source:       "13f32 >= 7f32",
			wantStackTop: value.True,
		},
		"7f32 >= 13f32": {
			source:       "7f32 >= 13f32",
			wantStackTop: value.False,
		},
		"7f32 >= 7.5f32": {
			source:       "7f32 >= 7.5f32",
			wantStackTop: value.False,
		},
		"7.5f32 >= 7f32": {
			source:       "7.5f32 >= 7f32",
			wantStackTop: value.True,
		},
		"7f32 >= 6.9f32": {
			source:       "7f32 >= 6.9f32",
			wantStackTop: value.True,
		},

		"6f32 >= 19.0": {
			source: "6f32 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Float32`",
			),
		},

		"6f32 >= 19": {
			source: "6f32 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19bf": {
			source: "6f32 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19f64": {
			source: "6f32 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19i64": {
			source: "6f32 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19i32": {
			source: "6f32 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19i16": {
			source: "6f32 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19i8": {
			source: "6f32 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19u64": {
			source: "6f32 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19u32": {
			source: "6f32 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19u16": {
			source: "6f32 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19u8": {
			source: "6f32 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float32`",
			),
		},

		// Int64
		"25i64 >= 25i64": {
			source:       "25i64 >= 25i64",
			wantStackTop: value.True,
		},
		"25i64 >= -25i64": {
			source:       "25i64 >= -25i64",
			wantStackTop: value.True,
		},
		"-25i64 >= 25i64": {
			source:       "-25i64 >= 25i64",
			wantStackTop: value.False,
		},
		"13i64 >= 7i64": {
			source:       "13i64 >= 7i64",
			wantStackTop: value.True,
		},
		"7i64 >= 13i64": {
			source:       "7i64 >= 13i64",
			wantStackTop: value.False,
		},

		"6i64 >= 19": {
			source: "6i64 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19.0": {
			source: "6i64 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19bf": {
			source: "6i64 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19f64": {
			source: "6i64 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19f32": {
			source: "6i64 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19i32": {
			source: "6i64 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19i16": {
			source: "6i64 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19i8": {
			source: "6i64 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19u64": {
			source: "6i64 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19u32": {
			source: "6i64 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19u16": {
			source: "6i64 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19u8": {
			source: "6i64 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int64`",
			),
		},

		// Int32
		"25i32 >= 25i32": {
			source:       "25i32 >= 25i32",
			wantStackTop: value.True,
		},
		"25i32 >= -25i32": {
			source:       "25i32 >= -25i32",
			wantStackTop: value.True,
		},
		"-25i32 >= 25i32": {
			source:       "-25i32 >= 25i32",
			wantStackTop: value.False,
		},
		"13i32 >= 7i32": {
			source:       "13i32 >= 7i32",
			wantStackTop: value.True,
		},
		"7i32 >= 13i32": {
			source:       "7i32 >= 13i32",
			wantStackTop: value.False,
		},

		"6i32 >= 19": {
			source: "6i32 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19.0": {
			source: "6i32 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19bf": {
			source: "6i32 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19f64": {
			source: "6i32 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19f32": {
			source: "6i32 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19i64": {
			source: "6i32 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19i16": {
			source: "6i32 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19i8": {
			source: "6i32 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19u64": {
			source: "6i32 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19u32": {
			source: "6i32 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19u16": {
			source: "6i32 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19u8": {
			source: "6i32 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int32`",
			),
		},

		// Int16
		"25i16 >= 25i16": {
			source:       "25i16 >= 25i16",
			wantStackTop: value.True,
		},
		"25i16 >= -25i16": {
			source:       "25i16 >= -25i16",
			wantStackTop: value.True,
		},
		"-25i16 >= 25i16": {
			source:       "-25i16 >= 25i16",
			wantStackTop: value.False,
		},
		"13i16 >= 7i16": {
			source:       "13i16 >= 7i16",
			wantStackTop: value.True,
		},
		"7i16 >= 13i16": {
			source:       "7i16 >= 13i16",
			wantStackTop: value.False,
		},

		"6i16 >= 19": {
			source: "6i16 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19.0": {
			source: "6i16 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19bf": {
			source: "6i16 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19f64": {
			source: "6i16 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19f32": {
			source: "6i16 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19i64": {
			source: "6i16 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19i32": {
			source: "6i16 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19i8": {
			source: "6i16 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19u64": {
			source: "6i16 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19u32": {
			source: "6i16 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19u16": {
			source: "6i16 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19u8": {
			source: "6i16 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int16`",
			),
		},

		// Int8
		"25i8 >= 25i8": {
			source:       "25i8 >= 25i8",
			wantStackTop: value.True,
		},
		"25i8 >= -25i8": {
			source:       "25i8 >= -25i8",
			wantStackTop: value.True,
		},
		"-25i8 >= 25i8": {
			source:       "-25i8 >= 25i8",
			wantStackTop: value.False,
		},
		"13i8 >= 7i8": {
			source:       "13i8 >= 7i8",
			wantStackTop: value.True,
		},
		"7i8 >= 13i8": {
			source:       "7i8 >= 13i8",
			wantStackTop: value.False,
		},

		"6i8 >= 19": {
			source: "6i8 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19.0": {
			source: "6i8 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19bf": {
			source: "6i8 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19f64": {
			source: "6i8 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19f32": {
			source: "6i8 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19i64": {
			source: "6i8 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19i32": {
			source: "6i8 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19i16": {
			source: "6i8 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19u64": {
			source: "6i8 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19u32": {
			source: "6i8 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19u16": {
			source: "6i8 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19u8": {
			source: "6i8 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int8`",
			),
		},

		// UInt64
		"25u64 >= 25u64": {
			source:       "25u64 >= 25u64",
			wantStackTop: value.True,
		},
		"13u64 >= 7u64": {
			source:       "13u64 >= 7u64",
			wantStackTop: value.True,
		},
		"7u64 >= 13u64": {
			source:       "7u64 >= 13u64",
			wantStackTop: value.False,
		},

		"6u64 >= 19": {
			source: "6u64 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19.0": {
			source: "6u64 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19bf": {
			source: "6u64 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19f64": {
			source: "6u64 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19f32": {
			source: "6u64 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19i64": {
			source: "6u64 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19i32": {
			source: "6u64 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19i16": {
			source: "6u64 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19i8": {
			source: "6u64 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19u32": {
			source: "6u64 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19u16": {
			source: "6u64 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19u8": {
			source: "6u64 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt64`",
			),
		},

		// UInt32
		"25u32 >= 25u32": {
			source:       "25u32 >= 25u32",
			wantStackTop: value.True,
		},
		"13u32 >= 7u32": {
			source:       "13u32 >= 7u32",
			wantStackTop: value.True,
		},
		"7u32 >= 13u32": {
			source:       "7u32 >= 13u32",
			wantStackTop: value.False,
		},

		"6u32 >= 19": {
			source: "6u32 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19.0": {
			source: "6u32 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19bf": {
			source: "6u32 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19f64": {
			source: "6u32 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19f32": {
			source: "6u32 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19i64": {
			source: "6u32 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19i32": {
			source: "6u32 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19i16": {
			source: "6u32 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19i8": {
			source: "6u32 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19u64": {
			source: "6u32 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19u16": {
			source: "6u32 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19u8": {
			source: "6u32 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt32`",
			),
		},

		// Int16
		"25u16 >= 25u16": {
			source:       "25u16 >= 25u16",
			wantStackTop: value.True,
		},
		"13u16 >= 7u16": {
			source:       "13u16 >= 7u16",
			wantStackTop: value.True,
		},
		"7u16 >= 13u16": {
			source:       "7u16 >= 13u16",
			wantStackTop: value.False,
		},

		"6u16 >= 19": {
			source: "6u16 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19.0": {
			source: "6u16 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19bf": {
			source: "6u16 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19f64": {
			source: "6u16 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19f32": {
			source: "6u16 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19i64": {
			source: "6u16 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19i32": {
			source: "6u16 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19i16": {
			source: "6u16 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19i8": {
			source: "6u16 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19u64": {
			source: "6u16 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19u32": {
			source: "6u16 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19u8": {
			source: "6u16 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt16`",
			),
		},

		// Int8
		"25u8 >= 25u8": {
			source:       "25u8 >= 25u8",
			wantStackTop: value.True,
		},
		"13u8 >= 7u8": {
			source:       "13u8 >= 7u8",
			wantStackTop: value.True,
		},
		"7u8 >= 13u8": {
			source:       "7u8 >= 13u8",
			wantStackTop: value.False,
		},

		"6u8 >= 19": {
			source: "6u8 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19.0": {
			source: "6u8 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19bf": {
			source: "6u8 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19f64": {
			source: "6u8 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19f32": {
			source: "6u8 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19i64": {
			source: "6u8 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19i32": {
			source: "6u8 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19i16": {
			source: "6u8 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19i8": {
			source: "6u8 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19u64": {
			source: "6u8 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19u32": {
			source: "6u8 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19u16": {
			source: "6u8 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt8`",
			),
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
			wantStackTop: value.False,
		},
		"'7' < '10'": {
			source:       "'7' < '10'",
			wantStackTop: value.False,
		},
		"'10' < '7'": {
			source:       "'10' < '7'",
			wantStackTop: value.True,
		},
		"'25' < '22'": {
			source:       "'25' < '22'",
			wantStackTop: value.False,
		},
		"'22' < '25'": {
			source:       "'22' < '25'",
			wantStackTop: value.True,
		},
		"'foo' < 'foo'": {
			source:       "'foo' < 'foo'",
			wantStackTop: value.False,
		},
		"'foo' < 'foa'": {
			source:       "'foo' < 'foa'",
			wantStackTop: value.False,
		},
		"'foa' < 'foo'": {
			source:       "'foa' < 'foo'",
			wantStackTop: value.True,
		},
		"'foo' < 'foo bar'": {
			source:       "'foo' < 'foo bar'",
			wantStackTop: value.True,
		},
		"'foo bar' < 'foo'": {
			source:       "'foo bar' < 'foo'",
			wantStackTop: value.False,
		},

		"'2' < `2`": {
			source:       "'2' < `2`",
			wantStackTop: value.False,
		},
		"'72' < `7`": {
			source:       "'72' < `7`",
			wantStackTop: value.False,
		},
		"'8' < `7`": {
			source:       "'8' < `7`",
			wantStackTop: value.False,
		},
		"'7' < `8`": {
			source:       "'7' < `8`",
			wantStackTop: value.True,
		},
		"'ba' < `b`": {
			source:       "'ba' < `b`",
			wantStackTop: value.False,
		},
		"'b' < `a`": {
			source:       "'b' < `a`",
			wantStackTop: value.False,
		},
		"'a' < `b`": {
			source:       "'a' < `b`",
			wantStackTop: value.True,
		},

		"'2' < 2.0": {
			source: "'2' < 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::String`",
			),
		},

		"'28' < 25.2bf": {
			source: "'28' < 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::String`",
			),
		},

		"'28.8' < 12.9f64": {
			source: "'28.8' < 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::String`",
			),
		},

		"'28.8' < 12.9f32": {
			source: "'28.8' < 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::String`",
			),
		},

		"'93' < 19i64": {
			source: "'93' < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::String`",
			),
		},

		"'93' < 19i32": {
			source: "'93' < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::String`",
			),
		},

		"'93' < 19i16": {
			source: "'93' < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::String`",
			),
		},

		"'93' < 19i8": {
			source: "'93' < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::String`",
			),
		},

		"'93' < 19u64": {
			source: "'93' < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::String`",
			),
		},

		"'93' < 19u32": {
			source: "'93' < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::String`",
			),
		},

		"'93' < 19u16": {
			source: "'93' < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::String`",
			),
		},

		"'93' < 19u8": {
			source: "'93' < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::String`",
			),
		},

		// Char
		"`2` < `2`": {
			source:       "`2` < `2`",
			wantStackTop: value.False,
		},
		"`8` < `7`": {
			source:       "`8` < `7`",
			wantStackTop: value.False,
		},
		"`7` < `8`": {
			source:       "`7` < `8`",
			wantStackTop: value.True,
		},
		"`b` < `a`": {
			source:       "`b` < `a`",
			wantStackTop: value.False,
		},
		"`a` < `b`": {
			source:       "`a` < `b`",
			wantStackTop: value.True,
		},

		"`2` < '2'": {
			source:       "`2` < '2'",
			wantStackTop: value.False,
		},
		"`7` < '72'": {
			source:       "`7` < '72'",
			wantStackTop: value.True,
		},
		"`8` < '7'": {
			source:       "`8` < '7'",
			wantStackTop: value.False,
		},
		"`7` < '8'": {
			source:       "`7` < '8'",
			wantStackTop: value.True,
		},
		"`b` < 'a'": {
			source:       "`b` < 'a'",
			wantStackTop: value.False,
		},
		"`b` < 'ba'": {
			source:       "`b` < 'ba'",
			wantStackTop: value.True,
		},
		"`a` < 'b'": {
			source:       "`a` < 'b'",
			wantStackTop: value.True,
		},

		"`2` < 2.0": {
			source: "`2` < 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Char`",
			),
		},
		"`i` < 25.2bf": {
			source: "`i` < 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Char`",
			),
		},
		"`f` < 12.9f64": {
			source: "`f` < 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Char`",
			),
		},
		"`0` < 12.9f32": {
			source: "`0` < 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Char`",
			),
		},
		"`9` < 19i64": {
			source: "`9` < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Char`",
			),
		},
		"`u` < 19i32": {
			source: "`u` < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Char`",
			),
		},
		"`4` < 19i16": {
			source: "`4` < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Char`",
			),
		},
		"`6` < 19i8": {
			source: "`6` < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Char`",
			),
		},
		"`9` < 19u64": {
			source: "`9` < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Char`",
			),
		},
		"`u` < 19u32": {
			source: "`u` < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Char`",
			),
		},
		"`4` < 19u16": {
			source: "`4` < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Char`",
			),
		},
		"`6` < 19u8": {
			source: "`6` < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Char`",
			),
		},

		// Int
		"25 < 25": {
			source:       "25 < 25",
			wantStackTop: value.False,
		},
		"25 < -25": {
			source:       "25 < -25",
			wantStackTop: value.False,
		},
		"-25 < 25": {
			source:       "-25 < 25",
			wantStackTop: value.True,
		},
		"13 < 7": {
			source:       "13 < 7",
			wantStackTop: value.False,
		},
		"7 < 13": {
			source:       "7 < 13",
			wantStackTop: value.True,
		},

		"25 < 25.0": {
			source:       "25 < 25.0",
			wantStackTop: value.False,
		},
		"25 < -25.0": {
			source:       "25 < -25.0",
			wantStackTop: value.False,
		},
		"-25 < 25.0": {
			source:       "-25 < 25.0",
			wantStackTop: value.True,
		},
		"13 < 7.0": {
			source:       "13 < 7.0",
			wantStackTop: value.False,
		},
		"7 < 13.0": {
			source:       "7 < 13.0",
			wantStackTop: value.True,
		},
		"7 < 7.5": {
			source:       "7 < 7.5",
			wantStackTop: value.True,
		},
		"7 < 6.9": {
			source:       "7 < 6.9",
			wantStackTop: value.False,
		},

		"25 < 25bf": {
			source:       "25 < 25bf",
			wantStackTop: value.False,
		},
		"25 < -25bf": {
			source:       "25 < -25bf",
			wantStackTop: value.False,
		},
		"-25 < 25bf": {
			source:       "-25 < 25bf",
			wantStackTop: value.True,
		},
		"13 < 7bf": {
			source:       "13 < 7bf",
			wantStackTop: value.False,
		},
		"7 < 13bf": {
			source:       "7 < 13bf",
			wantStackTop: value.True,
		},
		"7 < 7.5bf": {
			source:       "7 < 7.5bf",
			wantStackTop: value.True,
		},
		"7 < 6.9bf": {
			source:       "7 < 6.9bf",
			wantStackTop: value.False,
		},

		"6 < 19f64": {
			source: "6 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int`",
			),
		},
		"6 < 19f32": {
			source: "6 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int`",
			),
		},
		"6 < 19i64": {
			source: "6 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int`",
			),
		},
		"6 < 19i32": {
			source: "6 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int`",
			),
		},
		"6 < 19i16": {
			source: "6 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int`",
			),
		},
		"6 < 19i8": {
			source: "6 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int`",
			),
		},
		"6 < 19u64": {
			source: "6 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int`",
			),
		},
		"6 < 19u32": {
			source: "6 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int`",
			),
		},
		"6 < 19u16": {
			source: "6 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int`",
			),
		},
		"6 < 19u8": {
			source: "6 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int`",
			),
		},

		// Float
		"25.0 < 25.0": {
			source:       "25.0 < 25.0",
			wantStackTop: value.False,
		},
		"25.0 < -25.0": {
			source:       "25.0 < -25.0",
			wantStackTop: value.False,
		},
		"-25.0 < 25.0": {
			source:       "-25.0 < 25.0",
			wantStackTop: value.True,
		},
		"13.0 < 7.0": {
			source:       "13.0 < 7.0",
			wantStackTop: value.False,
		},
		"7.0 < 13.0": {
			source:       "7.0 < 13.0",
			wantStackTop: value.True,
		},
		"7.0 < 7.5": {
			source:       "7.0 < 7.5",
			wantStackTop: value.True,
		},
		"7.5 < 7.0": {
			source:       "7.5 < 7.0",
			wantStackTop: value.False,
		},
		"7.0 < 6.9": {
			source:       "7.0 < 6.9",
			wantStackTop: value.False,
		},

		"25.0 < 25": {
			source:       "25.0 < 25",
			wantStackTop: value.False,
		},
		"25.0 < -25": {
			source:       "25.0 < -25",
			wantStackTop: value.False,
		},
		"-25.0 < 25": {
			source:       "-25.0 < 25",
			wantStackTop: value.True,
		},
		"13.0 < 7": {
			source:       "13.0 < 7",
			wantStackTop: value.False,
		},
		"7.0 < 13": {
			source:       "7.0 < 13",
			wantStackTop: value.True,
		},
		"7.5 < 7": {
			source:       "7.5 < 7",
			wantStackTop: value.False,
		},

		"25.0 < 25bf": {
			source:       "25.0 < 25bf",
			wantStackTop: value.False,
		},
		"25.0 < -25bf": {
			source:       "25.0 < -25bf",
			wantStackTop: value.False,
		},
		"-25.0 < 25bf": {
			source:       "-25.0 < 25bf",
			wantStackTop: value.True,
		},
		"13.0 < 7bf": {
			source:       "13.0 < 7bf",
			wantStackTop: value.False,
		},
		"7.0 < 13bf": {
			source:       "7.0 < 13bf",
			wantStackTop: value.True,
		},
		"7.0 < 7.5bf": {
			source:       "7.0 < 7.5bf",
			wantStackTop: value.True,
		},
		"7.5 < 7bf": {
			source:       "7.5 < 7bf",
			wantStackTop: value.False,
		},
		"7.0 < 6.9bf": {
			source:       "7.0 < 6.9bf",
			wantStackTop: value.False,
		},

		"6.0 < 19f64": {
			source: "6.0 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 < 19f32": {
			source: "6.0 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 < 19i64": {
			source: "6.0 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 < 19i32": {
			source: "6.0 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 < 19i16": {
			source: "6.0 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 < 19i8": {
			source: "6.0 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 < 19u64": {
			source: "6.0 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 < 19u32": {
			source: "6.0 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 < 19u16": {
			source: "6.0 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 < 19u8": {
			source: "6.0 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float`",
			),
		},

		// BigFloat
		"25bf < 25.0": {
			source:       "25bf < 25.0",
			wantStackTop: value.False,
		},
		"25bf < -25.0": {
			source:       "25bf < -25.0",
			wantStackTop: value.False,
		},
		"-25bf < 25.0": {
			source:       "-25bf < 25.0",
			wantStackTop: value.True,
		},
		"13bf < 7.0": {
			source:       "13bf < 7.0",
			wantStackTop: value.False,
		},
		"7bf < 13.0": {
			source:       "7bf < 13.0",
			wantStackTop: value.True,
		},
		"7bf < 7.5": {
			source:       "7bf < 7.5",
			wantStackTop: value.True,
		},
		"7.5bf < 7.0": {
			source:       "7.5bf < 7.0",
			wantStackTop: value.False,
		},
		"7bf < 6.9": {
			source:       "7bf < 6.9",
			wantStackTop: value.False,
		},

		"25bf < 25": {
			source:       "25bf < 25",
			wantStackTop: value.False,
		},
		"25bf < -25": {
			source:       "25bf < -25",
			wantStackTop: value.False,
		},
		"-25bf < 25": {
			source:       "-25bf < 25",
			wantStackTop: value.True,
		},
		"13bf < 7": {
			source:       "13bf < 7",
			wantStackTop: value.False,
		},
		"7bf < 13": {
			source:       "7bf < 13",
			wantStackTop: value.True,
		},
		"7.5bf < 7": {
			source:       "7.5bf < 7",
			wantStackTop: value.False,
		},

		"25bf < 25bf": {
			source:       "25bf < 25bf",
			wantStackTop: value.False,
		},
		"25bf < -25bf": {
			source:       "25bf < -25bf",
			wantStackTop: value.False,
		},
		"-25bf < 25bf": {
			source:       "-25bf < 25bf",
			wantStackTop: value.True,
		},
		"13bf < 7bf": {
			source:       "13bf < 7bf",
			wantStackTop: value.False,
		},
		"7bf < 13bf": {
			source:       "7bf < 13bf",
			wantStackTop: value.True,
		},
		"7bf < 7.5bf": {
			source:       "7bf < 7.5bf",
			wantStackTop: value.True,
		},
		"7.5bf < 7bf": {
			source:       "7.5bf < 7bf",
			wantStackTop: value.False,
		},
		"7bf < 6.9bf": {
			source:       "7bf < 6.9bf",
			wantStackTop: value.False,
		},

		"6bf < 19f64": {
			source: "6bf < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19f32": {
			source: "6bf < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19i64": {
			source: "6bf < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19i32": {
			source: "6bf < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19i16": {
			source: "6bf < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19i8": {
			source: "6bf < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19u64": {
			source: "6bf < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19u32": {
			source: "6bf < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19u16": {
			source: "6bf < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19u8": {
			source: "6bf < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::BigFloat`",
			),
		},

		// Float64
		"25f64 < 25f64": {
			source:       "25f64 < 25f64",
			wantStackTop: value.False,
		},
		"25f64 < -25f64": {
			source:       "25f64 < -25f64",
			wantStackTop: value.False,
		},
		"-25f64 < 25f64": {
			source:       "-25f64 < 25f64",
			wantStackTop: value.True,
		},
		"13f64 < 7f64": {
			source:       "13f64 < 7f64",
			wantStackTop: value.False,
		},
		"7f64 < 13f64": {
			source:       "7f64 < 13f64",
			wantStackTop: value.True,
		},
		"7f64 < 7.5f64": {
			source:       "7f64 < 7.5f64",
			wantStackTop: value.True,
		},
		"7.5f64 < 7f64": {
			source:       "7.5f64 < 7f64",
			wantStackTop: value.False,
		},
		"7f64 < 6.9f64": {
			source:       "7f64 < 6.9f64",
			wantStackTop: value.False,
		},

		"6f64 < 19.0": {
			source: "6f64 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Float64`",
			),
		},

		"6f64 < 19": {
			source: "6f64 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19bf": {
			source: "6f64 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19f32": {
			source: "6f64 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19i64": {
			source: "6f64 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19i32": {
			source: "6f64 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19i16": {
			source: "6f64 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19i8": {
			source: "6f64 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19u64": {
			source: "6f64 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19u32": {
			source: "6f64 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19u16": {
			source: "6f64 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19u8": {
			source: "6f64 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float64`",
			),
		},

		// Float32
		"25f32 < 25f32": {
			source:       "25f32 < 25f32",
			wantStackTop: value.False,
		},
		"25f32 < -25f32": {
			source:       "25f32 < -25f32",
			wantStackTop: value.False,
		},
		"-25f32 < 25f32": {
			source:       "-25f32 < 25f32",
			wantStackTop: value.True,
		},
		"13f32 < 7f32": {
			source:       "13f32 < 7f32",
			wantStackTop: value.False,
		},
		"7f32 < 13f32": {
			source:       "7f32 < 13f32",
			wantStackTop: value.True,
		},
		"7f32 < 7.5f32": {
			source:       "7f32 < 7.5f32",
			wantStackTop: value.True,
		},
		"7.5f32 < 7f32": {
			source:       "7.5f32 < 7f32",
			wantStackTop: value.False,
		},
		"7f32 < 6.9f32": {
			source:       "7f32 < 6.9f32",
			wantStackTop: value.False,
		},

		"6f32 < 19.0": {
			source: "6f32 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Float32`",
			),
		},

		"6f32 < 19": {
			source: "6f32 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19bf": {
			source: "6f32 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19f64": {
			source: "6f32 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19i64": {
			source: "6f32 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19i32": {
			source: "6f32 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19i16": {
			source: "6f32 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19i8": {
			source: "6f32 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19u64": {
			source: "6f32 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19u32": {
			source: "6f32 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19u16": {
			source: "6f32 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19u8": {
			source: "6f32 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float32`",
			),
		},

		// Int64
		"25i64 < 25i64": {
			source:       "25i64 < 25i64",
			wantStackTop: value.False,
		},
		"25i64 < -25i64": {
			source:       "25i64 < -25i64",
			wantStackTop: value.False,
		},
		"-25i64 < 25i64": {
			source:       "-25i64 < 25i64",
			wantStackTop: value.True,
		},
		"13i64 < 7i64": {
			source:       "13i64 < 7i64",
			wantStackTop: value.False,
		},
		"7i64 < 13i64": {
			source:       "7i64 < 13i64",
			wantStackTop: value.True,
		},

		"6i64 < 19": {
			source: "6i64 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19.0": {
			source: "6i64 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19bf": {
			source: "6i64 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19f64": {
			source: "6i64 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19f32": {
			source: "6i64 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19i32": {
			source: "6i64 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19i16": {
			source: "6i64 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19i8": {
			source: "6i64 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19u64": {
			source: "6i64 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19u32": {
			source: "6i64 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19u16": {
			source: "6i64 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19u8": {
			source: "6i64 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int64`",
			),
		},

		// Int32
		"25i32 < 25i32": {
			source:       "25i32 < 25i32",
			wantStackTop: value.False,
		},
		"25i32 < -25i32": {
			source:       "25i32 < -25i32",
			wantStackTop: value.False,
		},
		"-25i32 < 25i32": {
			source:       "-25i32 < 25i32",
			wantStackTop: value.True,
		},
		"13i32 < 7i32": {
			source:       "13i32 < 7i32",
			wantStackTop: value.False,
		},
		"7i32 < 13i32": {
			source:       "7i32 < 13i32",
			wantStackTop: value.True,
		},

		"6i32 < 19": {
			source: "6i32 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19.0": {
			source: "6i32 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19bf": {
			source: "6i32 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19f64": {
			source: "6i32 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19f32": {
			source: "6i32 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19i64": {
			source: "6i32 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19i16": {
			source: "6i32 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19i8": {
			source: "6i32 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19u64": {
			source: "6i32 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19u32": {
			source: "6i32 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19u16": {
			source: "6i32 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19u8": {
			source: "6i32 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int32`",
			),
		},

		// Int16
		"25i16 < 25i16": {
			source:       "25i16 < 25i16",
			wantStackTop: value.False,
		},
		"25i16 < -25i16": {
			source:       "25i16 < -25i16",
			wantStackTop: value.False,
		},
		"-25i16 < 25i16": {
			source:       "-25i16 < 25i16",
			wantStackTop: value.True,
		},
		"13i16 < 7i16": {
			source:       "13i16 < 7i16",
			wantStackTop: value.False,
		},
		"7i16 < 13i16": {
			source:       "7i16 < 13i16",
			wantStackTop: value.True,
		},

		"6i16 < 19": {
			source: "6i16 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19.0": {
			source: "6i16 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19bf": {
			source: "6i16 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19f64": {
			source: "6i16 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19f32": {
			source: "6i16 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19i64": {
			source: "6i16 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19i32": {
			source: "6i16 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19i8": {
			source: "6i16 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19u64": {
			source: "6i16 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19u32": {
			source: "6i16 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19u16": {
			source: "6i16 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19u8": {
			source: "6i16 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int16`",
			),
		},

		// Int8
		"25i8 < 25i8": {
			source:       "25i8 < 25i8",
			wantStackTop: value.False,
		},
		"25i8 < -25i8": {
			source:       "25i8 < -25i8",
			wantStackTop: value.False,
		},
		"-25i8 < 25i8": {
			source:       "-25i8 < 25i8",
			wantStackTop: value.True,
		},
		"13i8 < 7i8": {
			source:       "13i8 < 7i8",
			wantStackTop: value.False,
		},
		"7i8 < 13i8": {
			source:       "7i8 < 13i8",
			wantStackTop: value.True,
		},

		"6i8 < 19": {
			source: "6i8 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19.0": {
			source: "6i8 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19bf": {
			source: "6i8 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19f64": {
			source: "6i8 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19f32": {
			source: "6i8 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19i64": {
			source: "6i8 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19i32": {
			source: "6i8 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19i16": {
			source: "6i8 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19u64": {
			source: "6i8 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19u32": {
			source: "6i8 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19u16": {
			source: "6i8 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19u8": {
			source: "6i8 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int8`",
			),
		},

		// UInt64
		"25u64 < 25u64": {
			source:       "25u64 < 25u64",
			wantStackTop: value.False,
		},
		"13u64 < 7u64": {
			source:       "13u64 < 7u64",
			wantStackTop: value.False,
		},
		"7u64 < 13u64": {
			source:       "7u64 < 13u64",
			wantStackTop: value.True,
		},

		"6u64 < 19": {
			source: "6u64 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19.0": {
			source: "6u64 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19bf": {
			source: "6u64 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19f64": {
			source: "6u64 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19f32": {
			source: "6u64 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19i64": {
			source: "6u64 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19i32": {
			source: "6u64 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19i16": {
			source: "6u64 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19i8": {
			source: "6u64 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19u32": {
			source: "6u64 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19u16": {
			source: "6u64 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19u8": {
			source: "6u64 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt64`",
			),
		},

		// UInt32
		"25u32 < 25u32": {
			source:       "25u32 < 25u32",
			wantStackTop: value.False,
		},
		"13u32 < 7u32": {
			source:       "13u32 < 7u32",
			wantStackTop: value.False,
		},
		"7u32 < 13u32": {
			source:       "7u32 < 13u32",
			wantStackTop: value.True,
		},

		"6u32 < 19": {
			source: "6u32 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19.0": {
			source: "6u32 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19bf": {
			source: "6u32 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19f64": {
			source: "6u32 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19f32": {
			source: "6u32 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19i64": {
			source: "6u32 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19i32": {
			source: "6u32 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19i16": {
			source: "6u32 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19i8": {
			source: "6u32 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19u64": {
			source: "6u32 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19u16": {
			source: "6u32 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19u8": {
			source: "6u32 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt32`",
			),
		},

		// Int16
		"25u16 < 25u16": {
			source:       "25u16 < 25u16",
			wantStackTop: value.False,
		},
		"13u16 < 7u16": {
			source:       "13u16 < 7u16",
			wantStackTop: value.False,
		},
		"7u16 < 13u16": {
			source:       "7u16 < 13u16",
			wantStackTop: value.True,
		},

		"6u16 < 19": {
			source: "6u16 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19.0": {
			source: "6u16 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19bf": {
			source: "6u16 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19f64": {
			source: "6u16 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19f32": {
			source: "6u16 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19i64": {
			source: "6u16 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19i32": {
			source: "6u16 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19i16": {
			source: "6u16 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19i8": {
			source: "6u16 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19u64": {
			source: "6u16 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19u32": {
			source: "6u16 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19u8": {
			source: "6u16 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt16`",
			),
		},

		// Int8
		"25u8 < 25u8": {
			source:       "25u8 < 25u8",
			wantStackTop: value.False,
		},
		"13u8 < 7u8": {
			source:       "13u8 < 7u8",
			wantStackTop: value.False,
		},
		"7u8 < 13u8": {
			source:       "7u8 < 13u8",
			wantStackTop: value.True,
		},

		"6u8 < 19": {
			source: "6u8 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19.0": {
			source: "6u8 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19bf": {
			source: "6u8 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19f64": {
			source: "6u8 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19f32": {
			source: "6u8 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19i64": {
			source: "6u8 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19i32": {
			source: "6u8 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19i16": {
			source: "6u8 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19i8": {
			source: "6u8 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19u64": {
			source: "6u8 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19u32": {
			source: "6u8 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19u16": {
			source: "6u8 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt8`",
			),
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
			wantStackTop: value.True,
		},
		"'7' <= '10'": {
			source:       "'7' <= '10'",
			wantStackTop: value.False,
		},
		"'10' <= '7'": {
			source:       "'10' <= '7'",
			wantStackTop: value.True,
		},
		"'25' <= '22'": {
			source:       "'25' <= '22'",
			wantStackTop: value.False,
		},
		"'22' <= '25'": {
			source:       "'22' <= '25'",
			wantStackTop: value.True,
		},
		"'foo' <= 'foo'": {
			source:       "'foo' <= 'foo'",
			wantStackTop: value.True,
		},
		"'foo' <= 'foa'": {
			source:       "'foo' <= 'foa'",
			wantStackTop: value.False,
		},
		"'foa' <= 'foo'": {
			source:       "'foa' <= 'foo'",
			wantStackTop: value.True,
		},
		"'foo' <= 'foo bar'": {
			source:       "'foo' <= 'foo bar'",
			wantStackTop: value.True,
		},
		"'foo bar' <= 'foo'": {
			source:       "'foo bar' <= 'foo'",
			wantStackTop: value.False,
		},

		"'2' <= `2`": {
			source:       "'2' <= `2`",
			wantStackTop: value.True,
		},
		"'72' <= `7`": {
			source:       "'72' <= `7`",
			wantStackTop: value.False,
		},
		"'8' <= `7`": {
			source:       "'8' <= `7`",
			wantStackTop: value.False,
		},
		"'7' <= `8`": {
			source:       "'7' <= `8`",
			wantStackTop: value.True,
		},
		"'ba' <= `b`": {
			source:       "'ba' <= `b`",
			wantStackTop: value.False,
		},
		"'b' <= `a`": {
			source:       "'b' <= `a`",
			wantStackTop: value.False,
		},
		"'a' <= `b`": {
			source:       "'a' <= `b`",
			wantStackTop: value.True,
		},

		"'2' <= 2.0": {
			source: "'2' <= 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::String`",
			),
		},

		"'28' <= 25.2bf": {
			source: "'28' <= 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::String`",
			),
		},

		"'28.8' <= 12.9f64": {
			source: "'28.8' <= 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::String`",
			),
		},

		"'28.8' <= 12.9f32": {
			source: "'28.8' <= 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::String`",
			),
		},

		"'93' <= 19i64": {
			source: "'93' <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::String`",
			),
		},

		"'93' <= 19i32": {
			source: "'93' <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::String`",
			),
		},

		"'93' <= 19i16": {
			source: "'93' <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::String`",
			),
		},

		"'93' <= 19i8": {
			source: "'93' <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::String`",
			),
		},

		"'93' <= 19u64": {
			source: "'93' <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::String`",
			),
		},

		"'93' <= 19u32": {
			source: "'93' <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::String`",
			),
		},

		"'93' <= 19u16": {
			source: "'93' <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::String`",
			),
		},

		"'93' <= 19u8": {
			source: "'93' <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::String`",
			),
		},

		// Char
		"`2` <= `2`": {
			source:       "`2` <= `2`",
			wantStackTop: value.True,
		},
		"`8` <= `7`": {
			source:       "`8` <= `7`",
			wantStackTop: value.False,
		},
		"`7` <= `8`": {
			source:       "`7` <= `8`",
			wantStackTop: value.True,
		},
		"`b` <= `a`": {
			source:       "`b` <= `a`",
			wantStackTop: value.False,
		},
		"`a` <= `b`": {
			source:       "`a` <= `b`",
			wantStackTop: value.True,
		},

		"`2` <= '2'": {
			source:       "`2` <= '2'",
			wantStackTop: value.True,
		},
		"`7` <= '72'": {
			source:       "`7` <= '72'",
			wantStackTop: value.True,
		},
		"`8` <= '7'": {
			source:       "`8` <= '7'",
			wantStackTop: value.False,
		},
		"`7` <= '8'": {
			source:       "`7` <= '8'",
			wantStackTop: value.True,
		},
		"`b` <= 'a'": {
			source:       "`b` <= 'a'",
			wantStackTop: value.False,
		},
		"`b` <= 'ba'": {
			source:       "`b` <= 'ba'",
			wantStackTop: value.True,
		},
		"`a` <= 'b'": {
			source:       "`a` <= 'b'",
			wantStackTop: value.True,
		},

		"`2` <= 2.0": {
			source: "`2` <= 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Char`",
			),
		},
		"`i` <= 25.2bf": {
			source: "`i` <= 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Char`",
			),
		},
		"`f` <= 12.9f64": {
			source: "`f` <= 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Char`",
			),
		},
		"`0` <= 12.9f32": {
			source: "`0` <= 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Char`",
			),
		},
		"`9` <= 19i64": {
			source: "`9` <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Char`",
			),
		},
		"`u` <= 19i32": {
			source: "`u` <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Char`",
			),
		},
		"`4` <= 19i16": {
			source: "`4` <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Char`",
			),
		},
		"`6` <= 19i8": {
			source: "`6` <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Char`",
			),
		},
		"`9` <= 19u64": {
			source: "`9` <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Char`",
			),
		},
		"`u` <= 19u32": {
			source: "`u` <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Char`",
			),
		},
		"`4` <= 19u16": {
			source: "`4` <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Char`",
			),
		},
		"`6` <= 19u8": {
			source: "`6` <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Char`",
			),
		},

		// Int
		"25 <= 25": {
			source:       "25 <= 25",
			wantStackTop: value.True,
		},
		"25 <= -25": {
			source:       "25 <= -25",
			wantStackTop: value.False,
		},
		"-25 <= 25": {
			source:       "-25 <= 25",
			wantStackTop: value.True,
		},
		"13 <= 7": {
			source:       "13 <= 7",
			wantStackTop: value.False,
		},
		"7 <= 13": {
			source:       "7 <= 13",
			wantStackTop: value.True,
		},

		"25 <= 25.0": {
			source:       "25 <= 25.0",
			wantStackTop: value.True,
		},
		"25 <= -25.0": {
			source:       "25 <= -25.0",
			wantStackTop: value.False,
		},
		"-25 <= 25.0": {
			source:       "-25 <= 25.0",
			wantStackTop: value.True,
		},
		"13 <= 7.0": {
			source:       "13 <= 7.0",
			wantStackTop: value.False,
		},
		"7 <= 13.0": {
			source:       "7 <= 13.0",
			wantStackTop: value.True,
		},
		"7 <= 7.5": {
			source:       "7 <= 7.5",
			wantStackTop: value.True,
		},
		"7 <= 6.9": {
			source:       "7 <= 6.9",
			wantStackTop: value.False,
		},

		"25 <= 25bf": {
			source:       "25 <= 25bf",
			wantStackTop: value.True,
		},
		"25 <= -25bf": {
			source:       "25 <= -25bf",
			wantStackTop: value.False,
		},
		"-25 <= 25bf": {
			source:       "-25 <= 25bf",
			wantStackTop: value.True,
		},
		"13 <= 7bf": {
			source:       "13 <= 7bf",
			wantStackTop: value.False,
		},
		"7 <= 13bf": {
			source:       "7 <= 13bf",
			wantStackTop: value.True,
		},
		"7 <= 7.5bf": {
			source:       "7 <= 7.5bf",
			wantStackTop: value.True,
		},
		"7 <= 6.9bf": {
			source:       "7 <= 6.9bf",
			wantStackTop: value.False,
		},

		"6 <= 19f64": {
			source: "6 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int`",
			),
		},
		"6 <= 19f32": {
			source: "6 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int`",
			),
		},
		"6 <= 19i64": {
			source: "6 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int`",
			),
		},
		"6 <= 19i32": {
			source: "6 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int`",
			),
		},
		"6 <= 19i16": {
			source: "6 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int`",
			),
		},
		"6 <= 19i8": {
			source: "6 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int`",
			),
		},
		"6 <= 19u64": {
			source: "6 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int`",
			),
		},
		"6 <= 19u32": {
			source: "6 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int`",
			),
		},
		"6 <= 19u16": {
			source: "6 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int`",
			),
		},
		"6 <= 19u8": {
			source: "6 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int`",
			),
		},

		// Float
		"25.0 <= 25.0": {
			source:       "25.0 <= 25.0",
			wantStackTop: value.True,
		},
		"25.0 <= -25.0": {
			source:       "25.0 <= -25.0",
			wantStackTop: value.False,
		},
		"-25.0 <= 25.0": {
			source:       "-25.0 <= 25.0",
			wantStackTop: value.True,
		},
		"13.0 <= 7.0": {
			source:       "13.0 <= 7.0",
			wantStackTop: value.False,
		},
		"7.0 <= 13.0": {
			source:       "7.0 <= 13.0",
			wantStackTop: value.True,
		},
		"7.0 <= 7.5": {
			source:       "7.0 <= 7.5",
			wantStackTop: value.True,
		},
		"7.5 <= 7.0": {
			source:       "7.5 <= 7.0",
			wantStackTop: value.False,
		},
		"7.0 <= 6.9": {
			source:       "7.0 <= 6.9",
			wantStackTop: value.False,
		},

		"25.0 <= 25": {
			source:       "25.0 <= 25",
			wantStackTop: value.True,
		},
		"25.0 <= -25": {
			source:       "25.0 <= -25",
			wantStackTop: value.False,
		},
		"-25.0 <= 25": {
			source:       "-25.0 <= 25",
			wantStackTop: value.True,
		},
		"13.0 <= 7": {
			source:       "13.0 <= 7",
			wantStackTop: value.False,
		},
		"7.0 <= 13": {
			source:       "7.0 <= 13",
			wantStackTop: value.True,
		},
		"7.5 <= 7": {
			source:       "7.5 <= 7",
			wantStackTop: value.False,
		},

		"25.0 <= 25bf": {
			source:       "25.0 <= 25bf",
			wantStackTop: value.True,
		},
		"25.0 <= -25bf": {
			source:       "25.0 <= -25bf",
			wantStackTop: value.False,
		},
		"-25.0 <= 25bf": {
			source:       "-25.0 <= 25bf",
			wantStackTop: value.True,
		},
		"13.0 <= 7bf": {
			source:       "13.0 <= 7bf",
			wantStackTop: value.False,
		},
		"7.0 <= 13bf": {
			source:       "7.0 <= 13bf",
			wantStackTop: value.True,
		},
		"7.0 <= 7.5bf": {
			source:       "7.0 <= 7.5bf",
			wantStackTop: value.True,
		},
		"7.5 <= 7bf": {
			source:       "7.5 <= 7bf",
			wantStackTop: value.False,
		},
		"7.0 <= 6.9bf": {
			source:       "7.0 <= 6.9bf",
			wantStackTop: value.False,
		},

		"6.0 <= 19f64": {
			source: "6.0 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19f32": {
			source: "6.0 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19i64": {
			source: "6.0 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19i32": {
			source: "6.0 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19i16": {
			source: "6.0 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19i8": {
			source: "6.0 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19u64": {
			source: "6.0 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19u32": {
			source: "6.0 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19u16": {
			source: "6.0 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19u8": {
			source: "6.0 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float`",
			),
		},

		// BigFloat
		"25bf <= 25.0": {
			source:       "25bf <= 25.0",
			wantStackTop: value.True,
		},
		"25bf <= -25.0": {
			source:       "25bf <= -25.0",
			wantStackTop: value.False,
		},
		"-25bf <= 25.0": {
			source:       "-25bf <= 25.0",
			wantStackTop: value.True,
		},
		"13bf <= 7.0": {
			source:       "13bf <= 7.0",
			wantStackTop: value.False,
		},
		"7bf <= 13.0": {
			source:       "7bf <= 13.0",
			wantStackTop: value.True,
		},
		"7bf <= 7.5": {
			source:       "7bf <= 7.5",
			wantStackTop: value.True,
		},
		"7.5bf <= 7.0": {
			source:       "7.5bf <= 7.0",
			wantStackTop: value.False,
		},
		"7bf <= 6.9": {
			source:       "7bf <= 6.9",
			wantStackTop: value.False,
		},

		"25bf <= 25": {
			source:       "25bf <= 25",
			wantStackTop: value.True,
		},
		"25bf <= -25": {
			source:       "25bf <= -25",
			wantStackTop: value.False,
		},
		"-25bf <= 25": {
			source:       "-25bf <= 25",
			wantStackTop: value.True,
		},
		"13bf <= 7": {
			source:       "13bf <= 7",
			wantStackTop: value.False,
		},
		"7bf <= 13": {
			source:       "7bf <= 13",
			wantStackTop: value.True,
		},
		"7.5bf <= 7": {
			source:       "7.5bf <= 7",
			wantStackTop: value.False,
		},

		"25bf <= 25bf": {
			source:       "25bf <= 25bf",
			wantStackTop: value.True,
		},
		"25bf <= -25bf": {
			source:       "25bf <= -25bf",
			wantStackTop: value.False,
		},
		"-25bf <= 25bf": {
			source:       "-25bf <= 25bf",
			wantStackTop: value.True,
		},
		"13bf <= 7bf": {
			source:       "13bf <= 7bf",
			wantStackTop: value.False,
		},
		"7bf <= 13bf": {
			source:       "7bf <= 13bf",
			wantStackTop: value.True,
		},
		"7bf <= 7.5bf": {
			source:       "7bf <= 7.5bf",
			wantStackTop: value.True,
		},
		"7.5bf <= 7bf": {
			source:       "7.5bf <= 7bf",
			wantStackTop: value.False,
		},
		"7bf <= 6.9bf": {
			source:       "7bf <= 6.9bf",
			wantStackTop: value.False,
		},

		"6bf <= 19f64": {
			source: "6bf <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19f32": {
			source: "6bf <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19i64": {
			source: "6bf <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19i32": {
			source: "6bf <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19i16": {
			source: "6bf <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19i8": {
			source: "6bf <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19u64": {
			source: "6bf <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19u32": {
			source: "6bf <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19u16": {
			source: "6bf <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19u8": {
			source: "6bf <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::BigFloat`",
			),
		},

		// Float64
		"25f64 <= 25f64": {
			source:       "25f64 <= 25f64",
			wantStackTop: value.True,
		},
		"25f64 <= -25f64": {
			source:       "25f64 <= -25f64",
			wantStackTop: value.False,
		},
		"-25f64 <= 25f64": {
			source:       "-25f64 <= 25f64",
			wantStackTop: value.True,
		},
		"13f64 <= 7f64": {
			source:       "13f64 <= 7f64",
			wantStackTop: value.False,
		},
		"7f64 <= 13f64": {
			source:       "7f64 <= 13f64",
			wantStackTop: value.True,
		},
		"7f64 <= 7.5f64": {
			source:       "7f64 <= 7.5f64",
			wantStackTop: value.True,
		},
		"7.5f64 <= 7f64": {
			source:       "7.5f64 <= 7f64",
			wantStackTop: value.False,
		},
		"7f64 <= 6.9f64": {
			source:       "7f64 <= 6.9f64",
			wantStackTop: value.False,
		},

		"6f64 <= 19.0": {
			source: "6f64 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Float64`",
			),
		},

		"6f64 <= 19": {
			source: "6f64 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19bf": {
			source: "6f64 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19f32": {
			source: "6f64 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19i64": {
			source: "6f64 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19i32": {
			source: "6f64 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19i16": {
			source: "6f64 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19i8": {
			source: "6f64 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19u64": {
			source: "6f64 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19u32": {
			source: "6f64 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19u16": {
			source: "6f64 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19u8": {
			source: "6f64 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float64`",
			),
		},

		// Float32
		"25f32 <= 25f32": {
			source:       "25f32 <= 25f32",
			wantStackTop: value.True,
		},
		"25f32 <= -25f32": {
			source:       "25f32 <= -25f32",
			wantStackTop: value.False,
		},
		"-25f32 <= 25f32": {
			source:       "-25f32 <= 25f32",
			wantStackTop: value.True,
		},
		"13f32 <= 7f32": {
			source:       "13f32 <= 7f32",
			wantStackTop: value.False,
		},
		"7f32 <= 13f32": {
			source:       "7f32 <= 13f32",
			wantStackTop: value.True,
		},
		"7f32 <= 7.5f32": {
			source:       "7f32 <= 7.5f32",
			wantStackTop: value.True,
		},
		"7.5f32 <= 7f32": {
			source:       "7.5f32 <= 7f32",
			wantStackTop: value.False,
		},
		"7f32 <= 6.9f32": {
			source:       "7f32 <= 6.9f32",
			wantStackTop: value.False,
		},

		"6f32 <= 19.0": {
			source: "6f32 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Float32`",
			),
		},

		"6f32 <= 19": {
			source: "6f32 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19bf": {
			source: "6f32 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19f64": {
			source: "6f32 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19i64": {
			source: "6f32 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19i32": {
			source: "6f32 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19i16": {
			source: "6f32 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19i8": {
			source: "6f32 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19u64": {
			source: "6f32 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19u32": {
			source: "6f32 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19u16": {
			source: "6f32 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19u8": {
			source: "6f32 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Float32`",
			),
		},

		// Int64
		"25i64 <= 25i64": {
			source:       "25i64 <= 25i64",
			wantStackTop: value.True,
		},
		"25i64 <= -25i64": {
			source:       "25i64 <= -25i64",
			wantStackTop: value.False,
		},
		"-25i64 <= 25i64": {
			source:       "-25i64 <= 25i64",
			wantStackTop: value.True,
		},
		"13i64 <= 7i64": {
			source:       "13i64 <= 7i64",
			wantStackTop: value.False,
		},
		"7i64 <= 13i64": {
			source:       "7i64 <= 13i64",
			wantStackTop: value.True,
		},

		"6i64 <= 19": {
			source: "6i64 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19.0": {
			source: "6i64 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19bf": {
			source: "6i64 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19f64": {
			source: "6i64 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19f32": {
			source: "6i64 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19i32": {
			source: "6i64 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19i16": {
			source: "6i64 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19i8": {
			source: "6i64 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19u64": {
			source: "6i64 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19u32": {
			source: "6i64 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19u16": {
			source: "6i64 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19u8": {
			source: "6i64 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int64`",
			),
		},

		// Int32
		"25i32 <= 25i32": {
			source:       "25i32 <= 25i32",
			wantStackTop: value.True,
		},
		"25i32 <= -25i32": {
			source:       "25i32 <= -25i32",
			wantStackTop: value.False,
		},
		"-25i32 <= 25i32": {
			source:       "-25i32 <= 25i32",
			wantStackTop: value.True,
		},
		"13i32 <= 7i32": {
			source:       "13i32 <= 7i32",
			wantStackTop: value.False,
		},
		"7i32 <= 13i32": {
			source:       "7i32 <= 13i32",
			wantStackTop: value.True,
		},

		"6i32 <= 19": {
			source: "6i32 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19.0": {
			source: "6i32 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19bf": {
			source: "6i32 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19f64": {
			source: "6i32 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19f32": {
			source: "6i32 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19i64": {
			source: "6i32 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19i16": {
			source: "6i32 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19i8": {
			source: "6i32 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19u64": {
			source: "6i32 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19u32": {
			source: "6i32 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19u16": {
			source: "6i32 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19u8": {
			source: "6i32 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int32`",
			),
		},

		// Int16
		"25i16 <= 25i16": {
			source:       "25i16 <= 25i16",
			wantStackTop: value.True,
		},
		"25i16 <= -25i16": {
			source:       "25i16 <= -25i16",
			wantStackTop: value.False,
		},
		"-25i16 <= 25i16": {
			source:       "-25i16 <= 25i16",
			wantStackTop: value.True,
		},
		"13i16 <= 7i16": {
			source:       "13i16 <= 7i16",
			wantStackTop: value.False,
		},
		"7i16 <= 13i16": {
			source:       "7i16 <= 13i16",
			wantStackTop: value.True,
		},

		"6i16 <= 19": {
			source: "6i16 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19.0": {
			source: "6i16 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19bf": {
			source: "6i16 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19f64": {
			source: "6i16 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19f32": {
			source: "6i16 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19i64": {
			source: "6i16 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19i32": {
			source: "6i16 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19i8": {
			source: "6i16 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19u64": {
			source: "6i16 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19u32": {
			source: "6i16 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19u16": {
			source: "6i16 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19u8": {
			source: "6i16 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int16`",
			),
		},

		// Int8
		"25i8 <= 25i8": {
			source:       "25i8 <= 25i8",
			wantStackTop: value.True,
		},
		"25i8 <= -25i8": {
			source:       "25i8 <= -25i8",
			wantStackTop: value.False,
		},
		"-25i8 <= 25i8": {
			source:       "-25i8 <= 25i8",
			wantStackTop: value.True,
		},
		"13i8 <= 7i8": {
			source:       "13i8 <= 7i8",
			wantStackTop: value.False,
		},
		"7i8 <= 13i8": {
			source:       "7i8 <= 13i8",
			wantStackTop: value.True,
		},

		"6i8 <= 19": {
			source: "6i8 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19.0": {
			source: "6i8 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19bf": {
			source: "6i8 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19f64": {
			source: "6i8 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19f32": {
			source: "6i8 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19i64": {
			source: "6i8 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19i32": {
			source: "6i8 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19i16": {
			source: "6i8 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19u64": {
			source: "6i8 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19u32": {
			source: "6i8 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19u16": {
			source: "6i8 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19u8": {
			source: "6i8 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::Int8`",
			),
		},

		// UInt64
		"25u64 <= 25u64": {
			source:       "25u64 <= 25u64",
			wantStackTop: value.True,
		},
		"13u64 <= 7u64": {
			source:       "13u64 <= 7u64",
			wantStackTop: value.False,
		},
		"7u64 <= 13u64": {
			source:       "7u64 <= 13u64",
			wantStackTop: value.True,
		},

		"6u64 <= 19": {
			source: "6u64 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19.0": {
			source: "6u64 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19bf": {
			source: "6u64 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19f64": {
			source: "6u64 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19f32": {
			source: "6u64 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19i64": {
			source: "6u64 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19i32": {
			source: "6u64 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19i16": {
			source: "6u64 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19i8": {
			source: "6u64 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19u32": {
			source: "6u64 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19u16": {
			source: "6u64 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19u8": {
			source: "6u64 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt64`",
			),
		},

		// UInt32
		"25u32 <= 25u32": {
			source:       "25u32 <= 25u32",
			wantStackTop: value.True,
		},
		"13u32 <= 7u32": {
			source:       "13u32 <= 7u32",
			wantStackTop: value.False,
		},
		"7u32 <= 13u32": {
			source:       "7u32 <= 13u32",
			wantStackTop: value.True,
		},

		"6u32 <= 19": {
			source: "6u32 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19.0": {
			source: "6u32 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19bf": {
			source: "6u32 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19f64": {
			source: "6u32 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19f32": {
			source: "6u32 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19i64": {
			source: "6u32 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19i32": {
			source: "6u32 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19i16": {
			source: "6u32 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19i8": {
			source: "6u32 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19u64": {
			source: "6u32 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19u16": {
			source: "6u32 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19u8": {
			source: "6u32 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt32`",
			),
		},

		// Int16
		"25u16 <= 25u16": {
			source:       "25u16 <= 25u16",
			wantStackTop: value.True,
		},
		"13u16 <= 7u16": {
			source:       "13u16 <= 7u16",
			wantStackTop: value.False,
		},
		"7u16 <= 13u16": {
			source:       "7u16 <= 13u16",
			wantStackTop: value.True,
		},

		"6u16 <= 19": {
			source: "6u16 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19.0": {
			source: "6u16 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19bf": {
			source: "6u16 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19f64": {
			source: "6u16 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19f32": {
			source: "6u16 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19i64": {
			source: "6u16 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19i32": {
			source: "6u16 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19i16": {
			source: "6u16 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19i8": {
			source: "6u16 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19u64": {
			source: "6u16 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19u32": {
			source: "6u16 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19u8": {
			source: "6u16 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` cannot be coerced into `Std::UInt16`",
			),
		},

		// Int8
		"25u8 <= 25u8": {
			source:       "25u8 <= 25u8",
			wantStackTop: value.True,
		},
		"13u8 <= 7u8": {
			source:       "13u8 <= 7u8",
			wantStackTop: value.False,
		},
		"7u8 <= 13u8": {
			source:       "7u8 <= 13u8",
			wantStackTop: value.True,
		},

		"6u8 <= 19": {
			source: "6u8 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19.0": {
			source: "6u8 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19bf": {
			source: "6u8 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19f64": {
			source: "6u8 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19f32": {
			source: "6u8 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19i64": {
			source: "6u8 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19i32": {
			source: "6u8 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19i16": {
			source: "6u8 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19i8": {
			source: "6u8 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19u64": {
			source: "6u8 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19u32": {
			source: "6u8 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` cannot be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19u16": {
			source: "6u8 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` cannot be coerced into `Std::UInt8`",
			),
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
		"'25' =~ '25'":   value.True,
		"'25' =~ '25.0'": value.False,
		"'25' =~ '7'":    value.False,

		"'7' =~ `7`":  value.True,
		"'a' =~ `a`":  value.True,
		"'7' =~ `5`":  value.False,
		"'ab' =~ `a`": value.False,

		"'25' =~ 25.0":   value.False,
		"'13.3' =~ 13.3": value.False,

		"'25' =~ 25bf":     value.False,
		"'13.3' =~ 13.3bf": value.False,

		"'25' =~ 25f64": value.False,

		"'25' =~ 25f32": value.False,

		"'1' =~ 1i64": value.False,

		"'5' =~ 5i32": value.False,

		"'5' =~ 5i16": value.False,

		"'5' =~ 5i8": value.False,

		"'1' =~ 1u64": value.False,

		"'5' =~ 5u32": value.False,

		"'5' =~ 5u16": value.False,

		"'5' =~ 5u8": value.False,

		// Char
		"`2` =~ '2'":   value.True,
		"`a` =~ 'a'":   value.True,
		"`a` =~ 'ab'":  value.False,
		"`2` =~ '2.0'": value.False,

		"`7` =~ `7`": value.True,
		"`a` =~ `a`": value.True,
		"`7` =~ `5`": value.False,
		"`a` =~ `b`": value.False,

		"`2` =~ 2.0": value.False,

		"`9` =~ 9bf": value.False,

		"`3` =~ 3f64": value.False,

		"`7` =~ 7f32": value.False,

		"`1` =~ 1i64": value.False,

		"`5` =~ 5i32": value.False,

		"`5` =~ 5i16": value.False,

		"`5` =~ 5i8": value.False,

		"`1` =~ 1u64": value.False,

		"`5` =~ 5u32": value.False,

		"`5` =~ 5u16": value.False,

		"`5` =~ 5u8": value.False,

		// Int
		"25 =~ 25":  value.True,
		"-25 =~ 25": value.False,
		"25 =~ -25": value.False,
		"25 =~ 28":  value.False,
		"28 =~ 25":  value.False,

		"25 =~ '25'": value.False,

		"7 =~ `7`": value.False,

		"-73 =~ 73.0": value.False,
		"73 =~ -73.0": value.False,
		"25 =~ 25.0":  value.True,
		"1 =~ 1.2":    value.False,

		"-73 =~ 73bf": value.False,
		"73 =~ -73bf": value.False,
		"25 =~ 25bf":  value.True,
		"1 =~ 1.2bf":  value.False,

		"-73 =~ 73f64": value.False,
		"73 =~ -73f64": value.False,
		"25 =~ 25f64":  value.True,
		"1 =~ 1.2f64":  value.False,

		"-73 =~ 73f32": value.False,
		"73 =~ -73f32": value.False,
		"25 =~ 25f32":  value.True,
		"1 =~ 1.2f32":  value.False,

		"1 =~ 1i64":   value.True,
		"4 =~ -4i64":  value.False,
		"-8 =~ 8i64":  value.False,
		"-8 =~ -8i64": value.True,
		"91 =~ 27i64": value.False,

		"5 =~ 5i32":  value.True,
		"4 =~ -4i32": value.False,
		"-8 =~ 8i32": value.False,
		"3 =~ 71i32": value.False,

		"5 =~ 5i16":  value.True,
		"4 =~ -4i16": value.False,
		"-8 =~ 8i16": value.False,
		"3 =~ 71i16": value.False,

		"5 =~ 5i8":  value.True,
		"4 =~ -4i8": value.False,
		"-8 =~ 8i8": value.False,
		"3 =~ 71i8": value.False,

		"1 =~ 1u64":   value.True,
		"-8 =~ 8u64":  value.False,
		"91 =~ 27u64": value.False,

		"5 =~ 5u32":  value.True,
		"-8 =~ 8u32": value.False,
		"3 =~ 71u32": value.False,

		"53000 =~ 32767u16": value.False,
		"5 =~ 5u16":         value.True,
		"-8 =~ 8u16":        value.False,
		"3 =~ 71u16":        value.False,

		"256 =~ 127u8": value.False,
		"5 =~ 5u8":     value.True,
		"-8 =~ 8u8":    value.False,
		"3 =~ 71u8":    value.False,

		// Int64
		"25i64 =~ 25":  value.True,
		"-25i64 =~ 25": value.False,
		"25i64 =~ -25": value.False,
		"25i64 =~ 28":  value.False,
		"28i64 =~ 25":  value.False,

		"25i64 =~ '25'": value.False,

		"7i64 =~ `7`": value.False,

		"-73i64 =~ 73.0": value.False,
		"73i64 =~ -73.0": value.False,
		"25i64 =~ 25.0":  value.True,
		"1i64 =~ 1.2":    value.False,

		"-73i64 =~ 73bf": value.False,
		"73i64 =~ -73bf": value.False,
		"25i64 =~ 25bf":  value.True,
		"1i64 =~ 1.2bf":  value.False,

		"-73i64 =~ 73f64": value.False,
		"73i64 =~ -73f64": value.False,
		"25i64 =~ 25f64":  value.True,
		"1i64 =~ 1.2f64":  value.False,

		"-73i64 =~ 73f32": value.False,
		"73i64 =~ -73f32": value.False,
		"25i64 =~ 25f32":  value.True,
		"1i64 =~ 1.2f32":  value.False,

		"1i64 =~ 1i64":   value.True,
		"4i64 =~ -4i64":  value.False,
		"-8i64 =~ 8i64":  value.False,
		"-8i64 =~ -8i64": value.True,
		"91i64 =~ 27i64": value.False,

		"5i64 =~ 5i32":  value.True,
		"4i64 =~ -4i32": value.False,
		"-8i64 =~ 8i32": value.False,
		"3i64 =~ 71i32": value.False,

		"5i64 =~ 5i16":  value.True,
		"4i64 =~ -4i16": value.False,
		"-8i64 =~ 8i16": value.False,
		"3i64 =~ 71i16": value.False,

		"5i64 =~ 5i8":  value.True,
		"4i64 =~ -4i8": value.False,
		"-8i64 =~ 8i8": value.False,
		"3i64 =~ 71i8": value.False,

		"1i64 =~ 1u64":   value.True,
		"-8i64 =~ 8u64":  value.False,
		"91i64 =~ 27u64": value.False,

		"5i64 =~ 5u32":  value.True,
		"-8i64 =~ 8u32": value.False,
		"3i64 =~ 71u32": value.False,

		"53000i64 =~ 32767u16": value.False,
		"5i64 =~ 5u16":         value.True,
		"-8i64 =~ 8u16":        value.False,
		"3i64 =~ 71u16":        value.False,

		"256i64 =~ 127u8": value.False,
		"5i64 =~ 5u8":     value.True,
		"-8i64 =~ 8u8":    value.False,
		"3i64 =~ 71u8":    value.False,

		// Int32
		"25i32 =~ 25":  value.True,
		"-25i32 =~ 25": value.False,
		"25i32 =~ -25": value.False,
		"25i32 =~ 28":  value.False,
		"28i32 =~ 25":  value.False,

		"25i32 =~ '25'": value.False,

		"7i32 =~ `7`": value.False,

		"-73i32 =~ 73.0": value.False,
		"73i32 =~ -73.0": value.False,
		"25i32 =~ 25.0":  value.True,
		"1i32 =~ 1.2":    value.False,

		"-73i32 =~ 73bf": value.False,
		"73i32 =~ -73bf": value.False,
		"25i32 =~ 25bf":  value.True,
		"1i32 =~ 1.2bf":  value.False,

		"-73i32 =~ 73f64": value.False,
		"73i32 =~ -73f64": value.False,
		"25i32 =~ 25f64":  value.True,
		"1i32 =~ 1.2f64":  value.False,

		"-73i32 =~ 73f32": value.False,
		"73i32 =~ -73f32": value.False,
		"25i32 =~ 25f32":  value.True,
		"1i32 =~ 1.2f32":  value.False,

		"1i32 =~ 1i64":   value.True,
		"4i32 =~ -4i64":  value.False,
		"-8i32 =~ 8i64":  value.False,
		"-8i32 =~ -8i64": value.True,
		"91i32 =~ 27i64": value.False,

		"5i32 =~ 5i32":  value.True,
		"4i32 =~ -4i32": value.False,
		"-8i32 =~ 8i32": value.False,
		"3i32 =~ 71i32": value.False,

		"5i32 =~ 5i16":  value.True,
		"4i32 =~ -4i16": value.False,
		"-8i32 =~ 8i16": value.False,
		"3i32 =~ 71i16": value.False,

		"5i32 =~ 5i8":  value.True,
		"4i32 =~ -4i8": value.False,
		"-8i32 =~ 8i8": value.False,
		"3i32 =~ 71i8": value.False,

		"1i32 =~ 1u64":   value.True,
		"-8i32 =~ 8u64":  value.False,
		"91i32 =~ 27u64": value.False,

		"5i32 =~ 5u32":  value.True,
		"-8i32 =~ 8u32": value.False,
		"3i32 =~ 71u32": value.False,

		"53000i32 =~ 32767u16": value.False,
		"5i32 =~ 5u16":         value.True,
		"-8i32 =~ 8u16":        value.False,
		"3i32 =~ 71u16":        value.False,

		"256i32 =~ 127u8": value.False,
		"5i32 =~ 5u8":     value.True,
		"-8i32 =~ 8u8":    value.False,
		"3i32 =~ 71u8":    value.False,

		// Int16
		"25i16 =~ 25":  value.True,
		"-25i16 =~ 25": value.False,
		"25i16 =~ -25": value.False,
		"25i16 =~ 28":  value.False,
		"28i16 =~ 25":  value.False,

		"25i16 =~ '25'": value.False,

		"7i16 =~ `7`": value.False,

		"-73i16 =~ 73.0": value.False,
		"73i16 =~ -73.0": value.False,
		"25i16 =~ 25.0":  value.True,
		"1i16 =~ 1.2":    value.False,

		"-73i16 =~ 73bf": value.False,
		"73i16 =~ -73bf": value.False,
		"25i16 =~ 25bf":  value.True,
		"1i16 =~ 1.2bf":  value.False,

		"-73i16 =~ 73f64": value.False,
		"73i16 =~ -73f64": value.False,
		"25i16 =~ 25f64":  value.True,
		"1i16 =~ 1.2f64":  value.False,

		"-73i16 =~ 73f32": value.False,
		"73i16 =~ -73f32": value.False,
		"25i16 =~ 25f32":  value.True,
		"1i16 =~ 1.2f32":  value.False,

		"1i16 =~ 1i64":   value.True,
		"4i16 =~ -4i64":  value.False,
		"-8i16 =~ 8i64":  value.False,
		"-8i16 =~ -8i64": value.True,
		"91i16 =~ 27i64": value.False,

		"5i16 =~ 5i32":  value.True,
		"4i16 =~ -4i32": value.False,
		"-8i16 =~ 8i32": value.False,
		"3i16 =~ 71i32": value.False,

		"5i16 =~ 5i16":  value.True,
		"4i16 =~ -4i16": value.False,
		"-8i16 =~ 8i16": value.False,
		"3i16 =~ 71i16": value.False,

		"5i16 =~ 5i8":  value.True,
		"4i16 =~ -4i8": value.False,
		"-8i16 =~ 8i8": value.False,
		"3i16 =~ 71i8": value.False,

		"1i16 =~ 1u64":   value.True,
		"-8i16 =~ 8u64":  value.False,
		"91i16 =~ 27u64": value.False,

		"5i16 =~ 5u32":  value.True,
		"-8i16 =~ 8u32": value.False,
		"3i16 =~ 71u32": value.False,

		"5i16 =~ 5u16":  value.True,
		"-8i16 =~ 8u16": value.False,
		"3i16 =~ 71u16": value.False,

		"256i16 =~ 127u8": value.False,
		"5i16 =~ 5u8":     value.True,
		"-8i16 =~ 8u8":    value.False,
		"3i16 =~ 71u8":    value.False,

		// Int8
		"25i8 =~ 25":  value.True,
		"-25i8 =~ 25": value.False,
		"25i8 =~ -25": value.False,
		"25i8 =~ 28":  value.False,
		"28i8 =~ 25":  value.False,

		"25i8 =~ '25'": value.False,

		"7i8 =~ `7`": value.False,

		"-73i8 =~ 73.0": value.False,
		"73i8 =~ -73.0": value.False,
		"25i8 =~ 25.0":  value.True,
		"1i8 =~ 1.2":    value.False,

		"-73i8 =~ 73bf": value.False,
		"73i8 =~ -73bf": value.False,
		"25i8 =~ 25bf":  value.True,
		"1i8 =~ 1.2bf":  value.False,

		"-73i8 =~ 73f64": value.False,
		"73i8 =~ -73f64": value.False,
		"25i8 =~ 25f64":  value.True,
		"1i8 =~ 1.2f64":  value.False,

		"-73i8 =~ 73f32": value.False,
		"73i8 =~ -73f32": value.False,
		"25i8 =~ 25f32":  value.True,
		"1i8 =~ 1.2f32":  value.False,

		"1i8 =~ 1i64":   value.True,
		"4i8 =~ -4i64":  value.False,
		"-8i8 =~ 8i64":  value.False,
		"-8i8 =~ -8i64": value.True,
		"91i8 =~ 27i64": value.False,

		"5i8 =~ 5i32":  value.True,
		"4i8 =~ -4i32": value.False,
		"-8i8 =~ 8i32": value.False,
		"3i8 =~ 71i32": value.False,

		"5i8 =~ 5i16":  value.True,
		"4i8 =~ -4i16": value.False,
		"-8i8 =~ 8i16": value.False,
		"3i8 =~ 71i16": value.False,

		"5i8 =~ 5i8":  value.True,
		"4i8 =~ -4i8": value.False,
		"-8i8 =~ 8i8": value.False,
		"3i8 =~ 71i8": value.False,

		"1i8 =~ 1u64":   value.True,
		"-8i8 =~ 8u64":  value.False,
		"91i8 =~ 27u64": value.False,

		"5i8 =~ 5u32":  value.True,
		"-8i8 =~ 8u32": value.False,
		"3i8 =~ 71u32": value.False,

		"5i8 =~ 5u16":  value.True,
		"-8i8 =~ 8u16": value.False,
		"3i8 =~ 71u16": value.False,

		"5i8 =~ 5u8":  value.True,
		"-8i8 =~ 8u8": value.False,
		"3i8 =~ 71u8": value.False,

		// UInt64
		"25u64 =~ 25":  value.True,
		"25u64 =~ -25": value.False,
		"25u64 =~ 28":  value.False,
		"28u64 =~ 25":  value.False,

		"25u64 =~ '25'": value.False,

		"7u64 =~ `7`": value.False,

		"73u64 =~ -73.0": value.False,
		"25u64 =~ 25.0":  value.True,
		"1u64 =~ 1.2":    value.False,

		"73u64 =~ -73bf": value.False,
		"25u64 =~ 25bf":  value.True,
		"1u64 =~ 1.2bf":  value.False,

		"73u64 =~ -73f64": value.False,
		"25u64 =~ 25f64":  value.True,
		"1u64 =~ 1.2f64":  value.False,

		"73u64 =~ -73f32": value.False,
		"25u64 =~ 25f32":  value.True,
		"1u64 =~ 1.2f32":  value.False,

		"1u64 =~ 1i64":   value.True,
		"4u64 =~ -4i64":  value.False,
		"91u64 =~ 27i64": value.False,

		"5u64 =~ 5i32":  value.True,
		"4u64 =~ -4i32": value.False,
		"3u64 =~ 71i32": value.False,

		"5u64 =~ 5i16":  value.True,
		"4u64 =~ -4i16": value.False,
		"3u64 =~ 71i16": value.False,

		"5u64 =~ 5i8":  value.True,
		"4u64 =~ -4i8": value.False,
		"3u64 =~ 71i8": value.False,

		"1u64 =~ 1u64":   value.True,
		"91u64 =~ 27u64": value.False,

		"5u64 =~ 5u32":  value.True,
		"3u64 =~ 71u32": value.False,

		"53000u64 =~ 32767u16": value.False,
		"5u64 =~ 5u16":         value.True,
		"3u64 =~ 71u16":        value.False,

		"256u64 =~ 127u8": value.False,
		"5u64 =~ 5u8":     value.True,
		"3u64 =~ 71u8":    value.False,

		// UInt32
		"25u32 =~ 25":  value.True,
		"25u32 =~ -25": value.False,
		"25u32 =~ 28":  value.False,
		"28u32 =~ 25":  value.False,

		"25u32 =~ '25'": value.False,

		"7u32 =~ `7`": value.False,

		"73u32 =~ -73.0": value.False,
		"25u32 =~ 25.0":  value.True,
		"1u32 =~ 1.2":    value.False,

		"73u32 =~ -73bf": value.False,
		"25u32 =~ 25bf":  value.True,
		"1u32 =~ 1.2bf":  value.False,

		"73u32 =~ -73f64": value.False,
		"25u32 =~ 25f64":  value.True,
		"1u32 =~ 1.2f64":  value.False,

		"73u32 =~ -73f32": value.False,
		"25u32 =~ 25f32":  value.True,
		"1u32 =~ 1.2f32":  value.False,

		"1u32 =~ 1i64":   value.True,
		"4u32 =~ -4i64":  value.False,
		"91u32 =~ 27i64": value.False,

		"5u32 =~ 5i32":  value.True,
		"4u32 =~ -4i32": value.False,
		"3u32 =~ 71i32": value.False,

		"5u32 =~ 5i16":  value.True,
		"4u32 =~ -4i16": value.False,
		"3u32 =~ 71i16": value.False,

		"5u32 =~ 5i8":  value.True,
		"4u32 =~ -4i8": value.False,
		"3u32 =~ 71i8": value.False,

		"1u32 =~ 1u64":   value.True,
		"91u32 =~ 27u64": value.False,

		"5u32 =~ 5u32":  value.True,
		"3u32 =~ 71u32": value.False,

		"53000u32 =~ 32767u16": value.False,
		"5u32 =~ 5u16":         value.True,
		"3u32 =~ 71u16":        value.False,

		"256u32 =~ 127u8": value.False,
		"5u32 =~ 5u8":     value.True,
		"3u32 =~ 71u8":    value.False,

		// UInt16
		"25u16 =~ 25":  value.True,
		"25u16 =~ -25": value.False,
		"25u16 =~ 28":  value.False,
		"28u16 =~ 25":  value.False,

		"25u16 =~ '25'": value.False,

		"7u16 =~ `7`": value.False,

		"73u16 =~ -73.0": value.False,
		"25u16 =~ 25.0":  value.True,
		"1u16 =~ 1.2":    value.False,

		"73u16 =~ -73bf": value.False,
		"25u16 =~ 25bf":  value.True,
		"1u16 =~ 1.2bf":  value.False,

		"73u16 =~ -73f64": value.False,
		"25u16 =~ 25f64":  value.True,
		"1u16 =~ 1.2f64":  value.False,

		"73u16 =~ -73f32": value.False,
		"25u16 =~ 25f32":  value.True,
		"1u16 =~ 1.2f32":  value.False,

		"1u16 =~ 1i64":   value.True,
		"4u16 =~ -4i64":  value.False,
		"91u16 =~ 27i64": value.False,

		"5u16 =~ 5i32":  value.True,
		"4u16 =~ -4i32": value.False,
		"3u16 =~ 71i32": value.False,

		"5u16 =~ 5i16":  value.True,
		"4u16 =~ -4i16": value.False,
		"3u16 =~ 71i16": value.False,

		"5u16 =~ 5i8":  value.True,
		"4u16 =~ -4i8": value.False,
		"3u16 =~ 71i8": value.False,

		"1u16 =~ 1u64":   value.True,
		"91u16 =~ 27u64": value.False,

		"5u16 =~ 5u32":  value.True,
		"3u16 =~ 71u32": value.False,

		"53000u16 =~ 32767u16": value.False,
		"5u16 =~ 5u16":         value.True,
		"3u16 =~ 71u16":        value.False,

		"256u16 =~ 127u8": value.False,
		"5u16 =~ 5u8":     value.True,
		"3u16 =~ 71u8":    value.False,

		// UInt8
		"25u8 =~ 25":  value.True,
		"25u8 =~ -25": value.False,
		"25u8 =~ 28":  value.False,
		"28u8 =~ 25":  value.False,

		"25u8 =~ '25'": value.False,

		"7u8 =~ `7`": value.False,

		"73u8 =~ -73.0": value.False,
		"25u8 =~ 25.0":  value.True,
		"1u8 =~ 1.2":    value.False,

		"73u8 =~ -73bf": value.False,
		"25u8 =~ 25bf":  value.True,
		"1u8 =~ 1.2bf":  value.False,

		"73u8 =~ -73f64": value.False,
		"25u8 =~ 25f64":  value.True,
		"1u8 =~ 1.2f64":  value.False,

		"73u8 =~ -73f32": value.False,
		"25u8 =~ 25f32":  value.True,
		"1u8 =~ 1.2f32":  value.False,

		"1u8 =~ 1i64":   value.True,
		"4u8 =~ -4i64":  value.False,
		"91u8 =~ 27i64": value.False,

		"5u8 =~ 5i32":  value.True,
		"4u8 =~ -4i32": value.False,
		"3u8 =~ 71i32": value.False,

		"5u8 =~ 5i16":  value.True,
		"4u8 =~ -4i16": value.False,
		"3u8 =~ 71i16": value.False,

		"5u8 =~ 5i8":  value.True,
		"4u8 =~ -4i8": value.False,
		"3u8 =~ 71i8": value.False,

		"1u8 =~ 1u64":   value.True,
		"91u8 =~ 27u64": value.False,

		"5u8 =~ 5u32":  value.True,
		"3u8 =~ 71u32": value.False,

		"5u8 =~ 5u16":  value.True,
		"3u8 =~ 71u16": value.False,

		"5u8 =~ 5u8":  value.True,
		"3u8 =~ 71u8": value.False,

		// Float
		"-73.0 =~ 73.0": value.False,
		"73.0 =~ -73.0": value.False,
		"25.0 =~ 25.0":  value.True,
		"1.0 =~ 1.2":    value.False,
		"1.2 =~ 1.0":    value.False,
		"78.5 =~ 78.5":  value.True,

		"8.25 =~ '8.25'": value.False,

		"4.0 =~ `4`": value.False,

		"25.0 =~ 25":  value.True,
		"32.3 =~ 32":  value.False,
		"-25.0 =~ 25": value.False,
		"25.0 =~ -25": value.False,
		"25.0 =~ 28":  value.False,
		"28.0 =~ 25":  value.False,

		"-73.0 =~ 73bf":  value.False,
		"73.0 =~ -73bf":  value.False,
		"25.0 =~ 25bf":   value.True,
		"1.0 =~ 1.2bf":   value.False,
		"15.5 =~ 15.5bf": value.True,

		"-73.0 =~ 73f64":    value.False,
		"73.0 =~ -73f64":    value.False,
		"25.0 =~ 25f64":     value.True,
		"1.0 =~ 1.2f64":     value.False,
		"15.26 =~ 15.26f64": value.True,

		"-73.0 =~ 73f32":  value.False,
		"73.0 =~ -73f32":  value.False,
		"25.0 =~ 25f32":   value.True,
		"1.0 =~ 1.2f32":   value.False,
		"15.5 =~ 15.5f32": value.True,

		"1.0 =~ 1i64":   value.True,
		"1.5 =~ 1i64":   value.False,
		"4.0 =~ -4i64":  value.False,
		"-8.0 =~ 8i64":  value.False,
		"-8.0 =~ -8i64": value.True,
		"91.0 =~ 27i64": value.False,

		"1.0 =~ 1i32":   value.True,
		"1.5 =~ 1i32":   value.False,
		"4.0 =~ -4i32":  value.False,
		"-8.0 =~ 8i32":  value.False,
		"-8.0 =~ -8i32": value.True,
		"91.0 =~ 27i32": value.False,

		"1.0 =~ 1i16":   value.True,
		"1.5 =~ 1i16":   value.False,
		"4.0 =~ -4i16":  value.False,
		"-8.0 =~ 8i16":  value.False,
		"-8.0 =~ -8i16": value.True,
		"91.0 =~ 27i16": value.False,

		"1.0 =~ 1i8":   value.True,
		"1.5 =~ 1i8":   value.False,
		"4.0 =~ -4i8":  value.False,
		"-8.0 =~ 8i8":  value.False,
		"-8.0 =~ -8i8": value.True,
		"91.0 =~ 27i8": value.False,

		"1.0 =~ 1u64":   value.True,
		"1.5 =~ 1u64":   value.False,
		"-8.0 =~ 8u64":  value.False,
		"91.0 =~ 27u64": value.False,

		"1.0 =~ 1u32":   value.True,
		"1.5 =~ 1u32":   value.False,
		"-8.0 =~ 8u32":  value.False,
		"91.0 =~ 27u32": value.False,

		"53000.0 =~ 32767u16": value.False,
		"1.0 =~ 1u16":         value.True,
		"1.5 =~ 1u16":         value.False,
		"-8.0 =~ 8u16":        value.False,
		"91.0 =~ 27u16":       value.False,

		"256.0 =~ 127u8": value.False,
		"1.0 =~ 1u8":     value.True,
		"1.5 =~ 1u8":     value.False,
		"-8.0 =~ 8u8":    value.False,
		"91.0 =~ 27u8":   value.False,

		// Float64
		"-73f64 =~ 73.0":  value.False,
		"73f64 =~ -73.0":  value.False,
		"25f64 =~ 25.0":   value.True,
		"1f64 =~ 1.2":     value.False,
		"1.2f64 =~ 1.0":   value.False,
		"78.5f64 =~ 78.5": value.True,

		"8.25f64 =~ '8.25'": value.False,

		"4f64 =~ `4`": value.False,

		"25f64 =~ 25":   value.True,
		"32.3f64 =~ 32": value.False,
		"-25f64 =~ 25":  value.False,
		"25f64 =~ -25":  value.False,
		"25f64 =~ 28":   value.False,
		"28f64 =~ 25":   value.False,

		"-73f64 =~ 73bf":    value.False,
		"73f64 =~ -73bf":    value.False,
		"25f64 =~ 25bf":     value.True,
		"1f64 =~ 1.2bf":     value.False,
		"15.5f64 =~ 15.5bf": value.True,

		"-73f64 =~ 73f64":      value.False,
		"73f64 =~ -73f64":      value.False,
		"25f64 =~ 25f64":       value.True,
		"1f64 =~ 1.2f64":       value.False,
		"15.26f64 =~ 15.26f64": value.True,

		"-73f64 =~ 73f32":    value.False,
		"73f64 =~ -73f32":    value.False,
		"25f64 =~ 25f32":     value.True,
		"1f64 =~ 1.2f32":     value.False,
		"15.5f64 =~ 15.5f32": value.True,

		"1f64 =~ 1i64":   value.True,
		"1.5f64 =~ 1i64": value.False,
		"4f64 =~ -4i64":  value.False,
		"-8f64 =~ 8i64":  value.False,
		"-8f64 =~ -8i64": value.True,
		"91f64 =~ 27i64": value.False,

		"1f64 =~ 1i32":   value.True,
		"1.5f64 =~ 1i32": value.False,
		"4f64 =~ -4i32":  value.False,
		"-8f64 =~ 8i32":  value.False,
		"-8f64 =~ -8i32": value.True,
		"91f64 =~ 27i32": value.False,

		"1f64 =~ 1i16":   value.True,
		"1.5f64 =~ 1i16": value.False,
		"4f64 =~ -4i16":  value.False,
		"-8f64 =~ 8i16":  value.False,
		"-8f64 =~ -8i16": value.True,
		"91f64 =~ 27i16": value.False,

		"1f64 =~ 1i8":   value.True,
		"1.5f64 =~ 1i8": value.False,
		"4f64 =~ -4i8":  value.False,
		"-8f64 =~ 8i8":  value.False,
		"-8f64 =~ -8i8": value.True,
		"91f64 =~ 27i8": value.False,

		"1f64 =~ 1u64":   value.True,
		"1.5f64 =~ 1u64": value.False,
		"-8f64 =~ 8u64":  value.False,
		"91f64 =~ 27u64": value.False,

		"1f64 =~ 1u32":   value.True,
		"1.5f64 =~ 1u32": value.False,
		"-8f64 =~ 8u32":  value.False,
		"91f64 =~ 27u32": value.False,

		"53000f64 =~ 32767u16": value.False,
		"1f64 =~ 1u16":         value.True,
		"1.5f64 =~ 1u16":       value.False,
		"-8f64 =~ 8u16":        value.False,
		"91f64 =~ 27u16":       value.False,

		"256f64 =~ 127u8": value.False,
		"1f64 =~ 1u8":     value.True,
		"1.5f64 =~ 1u8":   value.False,
		"-8f64 =~ 8u8":    value.False,
		"91f64 =~ 27u8":   value.False,

		// Float32
		"-73f32 =~ 73.0":  value.False,
		"73f32 =~ -73.0":  value.False,
		"25f32 =~ 25.0":   value.True,
		"1f32 =~ 1.2":     value.False,
		"1.2f32 =~ 1.0":   value.False,
		"78.5f32 =~ 78.5": value.True,

		"8.25f32 =~ '8.25'": value.False,

		"4f32 =~ `4`": value.False,

		"25f32 =~ 25":   value.True,
		"32.3f32 =~ 32": value.False,
		"-25f32 =~ 25":  value.False,
		"25f32 =~ -25":  value.False,
		"25f32 =~ 28":   value.False,
		"28f32 =~ 25":   value.False,

		"-73f32 =~ 73bf":    value.False,
		"73f32 =~ -73bf":    value.False,
		"25f32 =~ 25bf":     value.True,
		"1f32 =~ 1.2bf":     value.False,
		"15.5f32 =~ 15.5bf": value.True,

		"-73f32 =~ 73f64":    value.False,
		"73f32 =~ -73f64":    value.False,
		"25f32 =~ 25f64":     value.True,
		"1f32 =~ 1.2f64":     value.False,
		"15.5f32 =~ 15.5f64": value.True,

		"-73f32 =~ 73f32":    value.False,
		"73f32 =~ -73f32":    value.False,
		"25f32 =~ 25f32":     value.True,
		"1f32 =~ 1.2f32":     value.False,
		"15.5f32 =~ 15.5f32": value.True,

		"1f32 =~ 1i64":   value.True,
		"1.5f32 =~ 1i64": value.False,
		"4f32 =~ -4i64":  value.False,
		"-8f32 =~ 8i64":  value.False,
		"-8f32 =~ -8i64": value.True,
		"91f32 =~ 27i64": value.False,

		"1f32 =~ 1i32":   value.True,
		"1.5f32 =~ 1i32": value.False,
		"4f32 =~ -4i32":  value.False,
		"-8f32 =~ 8i32":  value.False,
		"-8f32 =~ -8i32": value.True,
		"91f32 =~ 27i32": value.False,

		"1f32 =~ 1i16":   value.True,
		"1.5f32 =~ 1i16": value.False,
		"4f32 =~ -4i16":  value.False,
		"-8f32 =~ 8i16":  value.False,
		"-8f32 =~ -8i16": value.True,
		"91f32 =~ 27i16": value.False,

		"1f32 =~ 1i8":   value.True,
		"1.5f32 =~ 1i8": value.False,
		"4f32 =~ -4i8":  value.False,
		"-8f32 =~ 8i8":  value.False,
		"-8f32 =~ -8i8": value.True,
		"91f32 =~ 27i8": value.False,

		"1f32 =~ 1u64":   value.True,
		"1.5f32 =~ 1u64": value.False,
		"-8f32 =~ 8u64":  value.False,
		"91f32 =~ 27u64": value.False,

		"1f32 =~ 1u32":   value.True,
		"1.5f32 =~ 1u32": value.False,
		"-8f32 =~ 8u32":  value.False,
		"91f32 =~ 27u32": value.False,

		"53000f32 =~ 32767u16": value.False,
		"1f32 =~ 1u16":         value.True,
		"1.5f32 =~ 1u16":       value.False,
		"-8f32 =~ 8u16":        value.False,
		"91f32 =~ 27u16":       value.False,

		"256f32 =~ 127u8": value.False,
		"1f32 =~ 1u8":     value.True,
		"1.5f32 =~ 1u8":   value.False,
		"-8f32 =~ 8u8":    value.False,
		"91f32 =~ 27u8":   value.False,
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
		"'25' !~ '25'":   value.False,
		"'25' !~ '25.0'": value.True,
		"'25' !~ '7'":    value.True,

		"'7' !~ `7`":  value.False,
		"'a' !~ `a`":  value.False,
		"'7' !~ `5`":  value.True,
		"'ab' !~ `a`": value.True,

		"'25' !~ 25.0":   value.True,
		"'13.3' !~ 13.3": value.True,

		"'25' !~ 25bf":     value.True,
		"'13.3' !~ 13.3bf": value.True,

		"'25' !~ 25f64": value.True,

		"'25' !~ 25f32": value.True,

		"'1' !~ 1i64": value.True,

		"'5' !~ 5i32": value.True,

		"'5' !~ 5i16": value.True,

		"'5' !~ 5i8": value.True,

		"'1' !~ 1u64": value.True,

		"'5' !~ 5u32": value.True,

		"'5' !~ 5u16": value.True,

		"'5' !~ 5u8": value.True,

		// Char
		"`2` !~ '2'":   value.False,
		"`a` !~ 'a'":   value.False,
		"`a` !~ 'ab'":  value.True,
		"`2` !~ '2.0'": value.True,

		"`7` !~ `7`": value.False,
		"`a` !~ `a`": value.False,
		"`7` !~ `5`": value.True,
		"`a` !~ `b`": value.True,

		"`2` !~ 2.0": value.True,

		"`9` !~ 9bf": value.True,

		"`3` !~ 3f64": value.True,

		"`7` !~ 7f32": value.True,

		"`1` !~ 1i64": value.True,

		"`5` !~ 5i32": value.True,

		"`5` !~ 5i16": value.True,

		"`5` !~ 5i8": value.True,

		"`1` !~ 1u64": value.True,

		"`5` !~ 5u32": value.True,

		"`5` !~ 5u16": value.True,

		"`5` !~ 5u8": value.True,

		// Int
		"25 !~ 25":  value.False,
		"-25 !~ 25": value.True,
		"25 !~ -25": value.True,
		"25 !~ 28":  value.True,
		"28 !~ 25":  value.True,

		"25 !~ '25'": value.True,

		"7 !~ `7`": value.True,

		"-73 !~ 73.0": value.True,
		"73 !~ -73.0": value.True,
		"25 !~ 25.0":  value.False,
		"1 !~ 1.2":    value.True,

		"-73 !~ 73bf": value.True,
		"73 !~ -73bf": value.True,
		"25 !~ 25bf":  value.False,
		"1 !~ 1.2bf":  value.True,

		"-73 !~ 73f64": value.True,
		"73 !~ -73f64": value.True,
		"25 !~ 25f64":  value.False,
		"1 !~ 1.2f64":  value.True,

		"-73 !~ 73f32": value.True,
		"73 !~ -73f32": value.True,
		"25 !~ 25f32":  value.False,
		"1 !~ 1.2f32":  value.True,

		"1 !~ 1i64":   value.False,
		"4 !~ -4i64":  value.True,
		"-8 !~ 8i64":  value.True,
		"-8 !~ -8i64": value.False,
		"91 !~ 27i64": value.True,

		"5 !~ 5i32":  value.False,
		"4 !~ -4i32": value.True,
		"-8 !~ 8i32": value.True,
		"3 !~ 71i32": value.True,

		"5 !~ 5i16":  value.False,
		"4 !~ -4i16": value.True,
		"-8 !~ 8i16": value.True,
		"3 !~ 71i16": value.True,

		"5 !~ 5i8":  value.False,
		"4 !~ -4i8": value.True,
		"-8 !~ 8i8": value.True,
		"3 !~ 71i8": value.True,

		"1 !~ 1u64":   value.False,
		"-8 !~ 8u64":  value.True,
		"91 !~ 27u64": value.True,

		"5 !~ 5u32":  value.False,
		"-8 !~ 8u32": value.True,
		"3 !~ 71u32": value.True,

		"53000 !~ 32767u16": value.True,
		"5 !~ 5u16":         value.False,
		"-8 !~ 8u16":        value.True,
		"3 !~ 71u16":        value.True,

		"256 !~ 127u8": value.True,
		"5 !~ 5u8":     value.False,
		"-8 !~ 8u8":    value.True,
		"3 !~ 71u8":    value.True,

		// Int64
		"25i64 !~ 25":  value.False,
		"-25i64 !~ 25": value.True,
		"25i64 !~ -25": value.True,
		"25i64 !~ 28":  value.True,
		"28i64 !~ 25":  value.True,

		"25i64 !~ '25'": value.True,

		"7i64 !~ `7`": value.True,

		"-73i64 !~ 73.0": value.True,
		"73i64 !~ -73.0": value.True,
		"25i64 !~ 25.0":  value.False,
		"1i64 !~ 1.2":    value.True,

		"-73i64 !~ 73bf": value.True,
		"73i64 !~ -73bf": value.True,
		"25i64 !~ 25bf":  value.False,
		"1i64 !~ 1.2bf":  value.True,

		"-73i64 !~ 73f64": value.True,
		"73i64 !~ -73f64": value.True,
		"25i64 !~ 25f64":  value.False,
		"1i64 !~ 1.2f64":  value.True,

		"-73i64 !~ 73f32": value.True,
		"73i64 !~ -73f32": value.True,
		"25i64 !~ 25f32":  value.False,
		"1i64 !~ 1.2f32":  value.True,

		"1i64 !~ 1i64":   value.False,
		"4i64 !~ -4i64":  value.True,
		"-8i64 !~ 8i64":  value.True,
		"-8i64 !~ -8i64": value.False,
		"91i64 !~ 27i64": value.True,

		"5i64 !~ 5i32":  value.False,
		"4i64 !~ -4i32": value.True,
		"-8i64 !~ 8i32": value.True,
		"3i64 !~ 71i32": value.True,

		"5i64 !~ 5i16":  value.False,
		"4i64 !~ -4i16": value.True,
		"-8i64 !~ 8i16": value.True,
		"3i64 !~ 71i16": value.True,

		"5i64 !~ 5i8":  value.False,
		"4i64 !~ -4i8": value.True,
		"-8i64 !~ 8i8": value.True,
		"3i64 !~ 71i8": value.True,

		"1i64 !~ 1u64":   value.False,
		"-8i64 !~ 8u64":  value.True,
		"91i64 !~ 27u64": value.True,

		"5i64 !~ 5u32":  value.False,
		"-8i64 !~ 8u32": value.True,
		"3i64 !~ 71u32": value.True,

		"53000i64 !~ 32767u16": value.True,
		"5i64 !~ 5u16":         value.False,
		"-8i64 !~ 8u16":        value.True,
		"3i64 !~ 71u16":        value.True,

		"256i64 !~ 127u8": value.True,
		"5i64 !~ 5u8":     value.False,
		"-8i64 !~ 8u8":    value.True,
		"3i64 !~ 71u8":    value.True,

		// Int32
		"25i32 !~ 25":  value.False,
		"-25i32 !~ 25": value.True,
		"25i32 !~ -25": value.True,
		"25i32 !~ 28":  value.True,
		"28i32 !~ 25":  value.True,

		"25i32 !~ '25'": value.True,

		"7i32 !~ `7`": value.True,

		"-73i32 !~ 73.0": value.True,
		"73i32 !~ -73.0": value.True,
		"25i32 !~ 25.0":  value.False,
		"1i32 !~ 1.2":    value.True,

		"-73i32 !~ 73bf": value.True,
		"73i32 !~ -73bf": value.True,
		"25i32 !~ 25bf":  value.False,
		"1i32 !~ 1.2bf":  value.True,

		"-73i32 !~ 73f64": value.True,
		"73i32 !~ -73f64": value.True,
		"25i32 !~ 25f64":  value.False,
		"1i32 !~ 1.2f64":  value.True,

		"-73i32 !~ 73f32": value.True,
		"73i32 !~ -73f32": value.True,
		"25i32 !~ 25f32":  value.False,
		"1i32 !~ 1.2f32":  value.True,

		"1i32 !~ 1i64":   value.False,
		"4i32 !~ -4i64":  value.True,
		"-8i32 !~ 8i64":  value.True,
		"-8i32 !~ -8i64": value.False,
		"91i32 !~ 27i64": value.True,

		"5i32 !~ 5i32":  value.False,
		"4i32 !~ -4i32": value.True,
		"-8i32 !~ 8i32": value.True,
		"3i32 !~ 71i32": value.True,

		"5i32 !~ 5i16":  value.False,
		"4i32 !~ -4i16": value.True,
		"-8i32 !~ 8i16": value.True,
		"3i32 !~ 71i16": value.True,

		"5i32 !~ 5i8":  value.False,
		"4i32 !~ -4i8": value.True,
		"-8i32 !~ 8i8": value.True,
		"3i32 !~ 71i8": value.True,

		"1i32 !~ 1u64":   value.False,
		"-8i32 !~ 8u64":  value.True,
		"91i32 !~ 27u64": value.True,

		"5i32 !~ 5u32":  value.False,
		"-8i32 !~ 8u32": value.True,
		"3i32 !~ 71u32": value.True,

		"53000i32 !~ 32767u16": value.True,
		"5i32 !~ 5u16":         value.False,
		"-8i32 !~ 8u16":        value.True,
		"3i32 !~ 71u16":        value.True,

		"256i32 !~ 127u8": value.True,
		"5i32 !~ 5u8":     value.False,
		"-8i32 !~ 8u8":    value.True,
		"3i32 !~ 71u8":    value.True,

		// Int16
		"25i16 !~ 25":  value.False,
		"-25i16 !~ 25": value.True,
		"25i16 !~ -25": value.True,
		"25i16 !~ 28":  value.True,
		"28i16 !~ 25":  value.True,

		"25i16 !~ '25'": value.True,

		"7i16 !~ `7`": value.True,

		"-73i16 !~ 73.0": value.True,
		"73i16 !~ -73.0": value.True,
		"25i16 !~ 25.0":  value.False,
		"1i16 !~ 1.2":    value.True,

		"-73i16 !~ 73bf": value.True,
		"73i16 !~ -73bf": value.True,
		"25i16 !~ 25bf":  value.False,
		"1i16 !~ 1.2bf":  value.True,

		"-73i16 !~ 73f64": value.True,
		"73i16 !~ -73f64": value.True,
		"25i16 !~ 25f64":  value.False,
		"1i16 !~ 1.2f64":  value.True,

		"-73i16 !~ 73f32": value.True,
		"73i16 !~ -73f32": value.True,
		"25i16 !~ 25f32":  value.False,
		"1i16 !~ 1.2f32":  value.True,

		"1i16 !~ 1i64":   value.False,
		"4i16 !~ -4i64":  value.True,
		"-8i16 !~ 8i64":  value.True,
		"-8i16 !~ -8i64": value.False,
		"91i16 !~ 27i64": value.True,

		"5i16 !~ 5i32":  value.False,
		"4i16 !~ -4i32": value.True,
		"-8i16 !~ 8i32": value.True,
		"3i16 !~ 71i32": value.True,

		"5i16 !~ 5i16":  value.False,
		"4i16 !~ -4i16": value.True,
		"-8i16 !~ 8i16": value.True,
		"3i16 !~ 71i16": value.True,

		"5i16 !~ 5i8":  value.False,
		"4i16 !~ -4i8": value.True,
		"-8i16 !~ 8i8": value.True,
		"3i16 !~ 71i8": value.True,

		"1i16 !~ 1u64":   value.False,
		"-8i16 !~ 8u64":  value.True,
		"91i16 !~ 27u64": value.True,

		"5i16 !~ 5u32":  value.False,
		"-8i16 !~ 8u32": value.True,
		"3i16 !~ 71u32": value.True,

		"5i16 !~ 5u16":  value.False,
		"-8i16 !~ 8u16": value.True,
		"3i16 !~ 71u16": value.True,

		"256i16 !~ 127u8": value.True,
		"5i16 !~ 5u8":     value.False,
		"-8i16 !~ 8u8":    value.True,
		"3i16 !~ 71u8":    value.True,

		// Int8
		"25i8 !~ 25":  value.False,
		"-25i8 !~ 25": value.True,
		"25i8 !~ -25": value.True,
		"25i8 !~ 28":  value.True,
		"28i8 !~ 25":  value.True,

		"25i8 !~ '25'": value.True,

		"7i8 !~ `7`": value.True,

		"-73i8 !~ 73.0": value.True,
		"73i8 !~ -73.0": value.True,
		"25i8 !~ 25.0":  value.False,
		"1i8 !~ 1.2":    value.True,

		"-73i8 !~ 73bf": value.True,
		"73i8 !~ -73bf": value.True,
		"25i8 !~ 25bf":  value.False,
		"1i8 !~ 1.2bf":  value.True,

		"-73i8 !~ 73f64": value.True,
		"73i8 !~ -73f64": value.True,
		"25i8 !~ 25f64":  value.False,
		"1i8 !~ 1.2f64":  value.True,

		"-73i8 !~ 73f32": value.True,
		"73i8 !~ -73f32": value.True,
		"25i8 !~ 25f32":  value.False,
		"1i8 !~ 1.2f32":  value.True,

		"1i8 !~ 1i64":   value.False,
		"4i8 !~ -4i64":  value.True,
		"-8i8 !~ 8i64":  value.True,
		"-8i8 !~ -8i64": value.False,
		"91i8 !~ 27i64": value.True,

		"5i8 !~ 5i32":  value.False,
		"4i8 !~ -4i32": value.True,
		"-8i8 !~ 8i32": value.True,
		"3i8 !~ 71i32": value.True,

		"5i8 !~ 5i16":  value.False,
		"4i8 !~ -4i16": value.True,
		"-8i8 !~ 8i16": value.True,
		"3i8 !~ 71i16": value.True,

		"5i8 !~ 5i8":  value.False,
		"4i8 !~ -4i8": value.True,
		"-8i8 !~ 8i8": value.True,
		"3i8 !~ 71i8": value.True,

		"1i8 !~ 1u64":   value.False,
		"-8i8 !~ 8u64":  value.True,
		"91i8 !~ 27u64": value.True,

		"5i8 !~ 5u32":  value.False,
		"-8i8 !~ 8u32": value.True,
		"3i8 !~ 71u32": value.True,

		"5i8 !~ 5u16":  value.False,
		"-8i8 !~ 8u16": value.True,
		"3i8 !~ 71u16": value.True,

		"5i8 !~ 5u8":  value.False,
		"-8i8 !~ 8u8": value.True,
		"3i8 !~ 71u8": value.True,

		// UInt64
		"25u64 !~ 25":  value.False,
		"25u64 !~ -25": value.True,
		"25u64 !~ 28":  value.True,
		"28u64 !~ 25":  value.True,

		"25u64 !~ '25'": value.True,

		"7u64 !~ `7`": value.True,

		"73u64 !~ -73.0": value.True,
		"25u64 !~ 25.0":  value.False,
		"1u64 !~ 1.2":    value.True,

		"73u64 !~ -73bf": value.True,
		"25u64 !~ 25bf":  value.False,
		"1u64 !~ 1.2bf":  value.True,

		"73u64 !~ -73f64": value.True,
		"25u64 !~ 25f64":  value.False,
		"1u64 !~ 1.2f64":  value.True,

		"73u64 !~ -73f32": value.True,
		"25u64 !~ 25f32":  value.False,
		"1u64 !~ 1.2f32":  value.True,

		"1u64 !~ 1i64":   value.False,
		"4u64 !~ -4i64":  value.True,
		"91u64 !~ 27i64": value.True,

		"5u64 !~ 5i32":  value.False,
		"4u64 !~ -4i32": value.True,
		"3u64 !~ 71i32": value.True,

		"5u64 !~ 5i16":  value.False,
		"4u64 !~ -4i16": value.True,
		"3u64 !~ 71i16": value.True,

		"5u64 !~ 5i8":  value.False,
		"4u64 !~ -4i8": value.True,
		"3u64 !~ 71i8": value.True,

		"1u64 !~ 1u64":   value.False,
		"91u64 !~ 27u64": value.True,

		"5u64 !~ 5u32":  value.False,
		"3u64 !~ 71u32": value.True,

		"53000u64 !~ 32767u16": value.True,
		"5u64 !~ 5u16":         value.False,
		"3u64 !~ 71u16":        value.True,

		"256u64 !~ 127u8": value.True,
		"5u64 !~ 5u8":     value.False,
		"3u64 !~ 71u8":    value.True,

		// UInt32
		"25u32 !~ 25":  value.False,
		"25u32 !~ -25": value.True,
		"25u32 !~ 28":  value.True,
		"28u32 !~ 25":  value.True,

		"25u32 !~ '25'": value.True,

		"7u32 !~ `7`": value.True,

		"73u32 !~ -73.0": value.True,
		"25u32 !~ 25.0":  value.False,
		"1u32 !~ 1.2":    value.True,

		"73u32 !~ -73bf": value.True,
		"25u32 !~ 25bf":  value.False,
		"1u32 !~ 1.2bf":  value.True,

		"73u32 !~ -73f64": value.True,
		"25u32 !~ 25f64":  value.False,
		"1u32 !~ 1.2f64":  value.True,

		"73u32 !~ -73f32": value.True,
		"25u32 !~ 25f32":  value.False,
		"1u32 !~ 1.2f32":  value.True,

		"1u32 !~ 1i64":   value.False,
		"4u32 !~ -4i64":  value.True,
		"91u32 !~ 27i64": value.True,

		"5u32 !~ 5i32":  value.False,
		"4u32 !~ -4i32": value.True,
		"3u32 !~ 71i32": value.True,

		"5u32 !~ 5i16":  value.False,
		"4u32 !~ -4i16": value.True,
		"3u32 !~ 71i16": value.True,

		"5u32 !~ 5i8":  value.False,
		"4u32 !~ -4i8": value.True,
		"3u32 !~ 71i8": value.True,

		"1u32 !~ 1u64":   value.False,
		"91u32 !~ 27u64": value.True,

		"5u32 !~ 5u32":  value.False,
		"3u32 !~ 71u32": value.True,

		"53000u32 !~ 32767u16": value.True,
		"5u32 !~ 5u16":         value.False,
		"3u32 !~ 71u16":        value.True,

		"256u32 !~ 127u8": value.True,
		"5u32 !~ 5u8":     value.False,
		"3u32 !~ 71u8":    value.True,

		// UInt16
		"25u16 !~ 25":  value.False,
		"25u16 !~ -25": value.True,
		"25u16 !~ 28":  value.True,
		"28u16 !~ 25":  value.True,

		"25u16 !~ '25'": value.True,

		"7u16 !~ `7`": value.True,

		"73u16 !~ -73.0": value.True,
		"25u16 !~ 25.0":  value.False,
		"1u16 !~ 1.2":    value.True,

		"73u16 !~ -73bf": value.True,
		"25u16 !~ 25bf":  value.False,
		"1u16 !~ 1.2bf":  value.True,

		"73u16 !~ -73f64": value.True,
		"25u16 !~ 25f64":  value.False,
		"1u16 !~ 1.2f64":  value.True,

		"73u16 !~ -73f32": value.True,
		"25u16 !~ 25f32":  value.False,
		"1u16 !~ 1.2f32":  value.True,

		"1u16 !~ 1i64":   value.False,
		"4u16 !~ -4i64":  value.True,
		"91u16 !~ 27i64": value.True,

		"5u16 !~ 5i32":  value.False,
		"4u16 !~ -4i32": value.True,
		"3u16 !~ 71i32": value.True,

		"5u16 !~ 5i16":  value.False,
		"4u16 !~ -4i16": value.True,
		"3u16 !~ 71i16": value.True,

		"5u16 !~ 5i8":  value.False,
		"4u16 !~ -4i8": value.True,
		"3u16 !~ 71i8": value.True,

		"1u16 !~ 1u64":   value.False,
		"91u16 !~ 27u64": value.True,

		"5u16 !~ 5u32":  value.False,
		"3u16 !~ 71u32": value.True,

		"53000u16 !~ 32767u16": value.True,
		"5u16 !~ 5u16":         value.False,
		"3u16 !~ 71u16":        value.True,

		"256u16 !~ 127u8": value.True,
		"5u16 !~ 5u8":     value.False,
		"3u16 !~ 71u8":    value.True,

		// UInt8
		"25u8 !~ 25":  value.False,
		"25u8 !~ -25": value.True,
		"25u8 !~ 28":  value.True,
		"28u8 !~ 25":  value.True,

		"25u8 !~ '25'": value.True,

		"7u8 !~ `7`": value.True,

		"73u8 !~ -73.0": value.True,
		"25u8 !~ 25.0":  value.False,
		"1u8 !~ 1.2":    value.True,

		"73u8 !~ -73bf": value.True,
		"25u8 !~ 25bf":  value.False,
		"1u8 !~ 1.2bf":  value.True,

		"73u8 !~ -73f64": value.True,
		"25u8 !~ 25f64":  value.False,
		"1u8 !~ 1.2f64":  value.True,

		"73u8 !~ -73f32": value.True,
		"25u8 !~ 25f32":  value.False,
		"1u8 !~ 1.2f32":  value.True,

		"1u8 !~ 1i64":   value.False,
		"4u8 !~ -4i64":  value.True,
		"91u8 !~ 27i64": value.True,

		"5u8 !~ 5i32":  value.False,
		"4u8 !~ -4i32": value.True,
		"3u8 !~ 71i32": value.True,

		"5u8 !~ 5i16":  value.False,
		"4u8 !~ -4i16": value.True,
		"3u8 !~ 71i16": value.True,

		"5u8 !~ 5i8":  value.False,
		"4u8 !~ -4i8": value.True,
		"3u8 !~ 71i8": value.True,

		"1u8 !~ 1u64":   value.False,
		"91u8 !~ 27u64": value.True,

		"5u8 !~ 5u32":  value.False,
		"3u8 !~ 71u32": value.True,

		"5u8 !~ 5u16":  value.False,
		"3u8 !~ 71u16": value.True,

		"5u8 !~ 5u8":  value.False,
		"3u8 !~ 71u8": value.True,

		// Float
		"-73.0 !~ 73.0": value.True,
		"73.0 !~ -73.0": value.True,
		"25.0 !~ 25.0":  value.False,
		"1.0 !~ 1.2":    value.True,
		"1.2 !~ 1.0":    value.True,
		"78.5 !~ 78.5":  value.False,

		"8.25 !~ '8.25'": value.True,

		"4.0 !~ `4`": value.True,

		"25.0 !~ 25":  value.False,
		"32.3 !~ 32":  value.True,
		"-25.0 !~ 25": value.True,
		"25.0 !~ -25": value.True,
		"25.0 !~ 28":  value.True,
		"28.0 !~ 25":  value.True,

		"-73.0 !~ 73bf":  value.True,
		"73.0 !~ -73bf":  value.True,
		"25.0 !~ 25bf":   value.False,
		"1.0 !~ 1.2bf":   value.True,
		"15.5 !~ 15.5bf": value.False,

		"-73.0 !~ 73f64":    value.True,
		"73.0 !~ -73f64":    value.True,
		"25.0 !~ 25f64":     value.False,
		"1.0 !~ 1.2f64":     value.True,
		"15.26 !~ 15.26f64": value.False,

		"-73.0 !~ 73f32":  value.True,
		"73.0 !~ -73f32":  value.True,
		"25.0 !~ 25f32":   value.False,
		"1.0 !~ 1.2f32":   value.True,
		"15.5 !~ 15.5f32": value.False,

		"1.0 !~ 1i64":   value.False,
		"1.5 !~ 1i64":   value.True,
		"4.0 !~ -4i64":  value.True,
		"-8.0 !~ 8i64":  value.True,
		"-8.0 !~ -8i64": value.False,
		"91.0 !~ 27i64": value.True,

		"1.0 !~ 1i32":   value.False,
		"1.5 !~ 1i32":   value.True,
		"4.0 !~ -4i32":  value.True,
		"-8.0 !~ 8i32":  value.True,
		"-8.0 !~ -8i32": value.False,
		"91.0 !~ 27i32": value.True,

		"1.0 !~ 1i16":   value.False,
		"1.5 !~ 1i16":   value.True,
		"4.0 !~ -4i16":  value.True,
		"-8.0 !~ 8i16":  value.True,
		"-8.0 !~ -8i16": value.False,
		"91.0 !~ 27i16": value.True,

		"1.0 !~ 1i8":   value.False,
		"1.5 !~ 1i8":   value.True,
		"4.0 !~ -4i8":  value.True,
		"-8.0 !~ 8i8":  value.True,
		"-8.0 !~ -8i8": value.False,
		"91.0 !~ 27i8": value.True,

		"1.0 !~ 1u64":   value.False,
		"1.5 !~ 1u64":   value.True,
		"-8.0 !~ 8u64":  value.True,
		"91.0 !~ 27u64": value.True,

		"1.0 !~ 1u32":   value.False,
		"1.5 !~ 1u32":   value.True,
		"-8.0 !~ 8u32":  value.True,
		"91.0 !~ 27u32": value.True,

		"53000.0 !~ 32767u16": value.True,
		"1.0 !~ 1u16":         value.False,
		"1.5 !~ 1u16":         value.True,
		"-8.0 !~ 8u16":        value.True,
		"91.0 !~ 27u16":       value.True,

		"256.0 !~ 127u8": value.True,
		"1.0 !~ 1u8":     value.False,
		"1.5 !~ 1u8":     value.True,
		"-8.0 !~ 8u8":    value.True,
		"91.0 !~ 27u8":   value.True,

		// Float64
		"-73f64 !~ 73.0":  value.True,
		"73f64 !~ -73.0":  value.True,
		"25f64 !~ 25.0":   value.False,
		"1f64 !~ 1.2":     value.True,
		"1.2f64 !~ 1.0":   value.True,
		"78.5f64 !~ 78.5": value.False,

		"8.25f64 !~ '8.25'": value.True,

		"4f64 !~ `4`": value.True,

		"25f64 !~ 25":   value.False,
		"32.3f64 !~ 32": value.True,
		"-25f64 !~ 25":  value.True,
		"25f64 !~ -25":  value.True,
		"25f64 !~ 28":   value.True,
		"28f64 !~ 25":   value.True,

		"-73f64 !~ 73bf":    value.True,
		"73f64 !~ -73bf":    value.True,
		"25f64 !~ 25bf":     value.False,
		"1f64 !~ 1.2bf":     value.True,
		"15.5f64 !~ 15.5bf": value.False,

		"-73f64 !~ 73f64":      value.True,
		"73f64 !~ -73f64":      value.True,
		"25f64 !~ 25f64":       value.False,
		"1f64 !~ 1.2f64":       value.True,
		"15.26f64 !~ 15.26f64": value.False,

		"-73f64 !~ 73f32":    value.True,
		"73f64 !~ -73f32":    value.True,
		"25f64 !~ 25f32":     value.False,
		"1f64 !~ 1.2f32":     value.True,
		"15.5f64 !~ 15.5f32": value.False,

		"1f64 !~ 1i64":   value.False,
		"1.5f64 !~ 1i64": value.True,
		"4f64 !~ -4i64":  value.True,
		"-8f64 !~ 8i64":  value.True,
		"-8f64 !~ -8i64": value.False,
		"91f64 !~ 27i64": value.True,

		"1f64 !~ 1i32":   value.False,
		"1.5f64 !~ 1i32": value.True,
		"4f64 !~ -4i32":  value.True,
		"-8f64 !~ 8i32":  value.True,
		"-8f64 !~ -8i32": value.False,
		"91f64 !~ 27i32": value.True,

		"1f64 !~ 1i16":   value.False,
		"1.5f64 !~ 1i16": value.True,
		"4f64 !~ -4i16":  value.True,
		"-8f64 !~ 8i16":  value.True,
		"-8f64 !~ -8i16": value.False,
		"91f64 !~ 27i16": value.True,

		"1f64 !~ 1i8":   value.False,
		"1.5f64 !~ 1i8": value.True,
		"4f64 !~ -4i8":  value.True,
		"-8f64 !~ 8i8":  value.True,
		"-8f64 !~ -8i8": value.False,
		"91f64 !~ 27i8": value.True,

		"1f64 !~ 1u64":   value.False,
		"1.5f64 !~ 1u64": value.True,
		"-8f64 !~ 8u64":  value.True,
		"91f64 !~ 27u64": value.True,

		"1f64 !~ 1u32":   value.False,
		"1.5f64 !~ 1u32": value.True,
		"-8f64 !~ 8u32":  value.True,
		"91f64 !~ 27u32": value.True,

		"53000f64 !~ 32767u16": value.True,
		"1f64 !~ 1u16":         value.False,
		"1.5f64 !~ 1u16":       value.True,
		"-8f64 !~ 8u16":        value.True,
		"91f64 !~ 27u16":       value.True,

		"256f64 !~ 127u8": value.True,
		"1f64 !~ 1u8":     value.False,
		"1.5f64 !~ 1u8":   value.True,
		"-8f64 !~ 8u8":    value.True,
		"91f64 !~ 27u8":   value.True,

		// Float32
		"-73f32 !~ 73.0":  value.True,
		"73f32 !~ -73.0":  value.True,
		"25f32 !~ 25.0":   value.False,
		"1f32 !~ 1.2":     value.True,
		"1.2f32 !~ 1.0":   value.True,
		"78.5f32 !~ 78.5": value.False,

		"8.25f32 !~ '8.25'": value.True,

		"4f32 !~ `4`": value.True,

		"25f32 !~ 25":   value.False,
		"32.3f32 !~ 32": value.True,
		"-25f32 !~ 25":  value.True,
		"25f32 !~ -25":  value.True,
		"25f32 !~ 28":   value.True,
		"28f32 !~ 25":   value.True,

		"-73f32 !~ 73bf":    value.True,
		"73f32 !~ -73bf":    value.True,
		"25f32 !~ 25bf":     value.False,
		"1f32 !~ 1.2bf":     value.True,
		"15.5f32 !~ 15.5bf": value.False,

		"-73f32 !~ 73f64":    value.True,
		"73f32 !~ -73f64":    value.True,
		"25f32 !~ 25f64":     value.False,
		"1f32 !~ 1.2f64":     value.True,
		"15.5f32 !~ 15.5f64": value.False,

		"-73f32 !~ 73f32":    value.True,
		"73f32 !~ -73f32":    value.True,
		"25f32 !~ 25f32":     value.False,
		"1f32 !~ 1.2f32":     value.True,
		"15.5f32 !~ 15.5f32": value.False,

		"1f32 !~ 1i64":   value.False,
		"1.5f32 !~ 1i64": value.True,
		"4f32 !~ -4i64":  value.True,
		"-8f32 !~ 8i64":  value.True,
		"-8f32 !~ -8i64": value.False,
		"91f32 !~ 27i64": value.True,

		"1f32 !~ 1i32":   value.False,
		"1.5f32 !~ 1i32": value.True,
		"4f32 !~ -4i32":  value.True,
		"-8f32 !~ 8i32":  value.True,
		"-8f32 !~ -8i32": value.False,
		"91f32 !~ 27i32": value.True,

		"1f32 !~ 1i16":   value.False,
		"1.5f32 !~ 1i16": value.True,
		"4f32 !~ -4i16":  value.True,
		"-8f32 !~ 8i16":  value.True,
		"-8f32 !~ -8i16": value.False,
		"91f32 !~ 27i16": value.True,

		"1f32 !~ 1i8":   value.False,
		"1.5f32 !~ 1i8": value.True,
		"4f32 !~ -4i8":  value.True,
		"-8f32 !~ 8i8":  value.True,
		"-8f32 !~ -8i8": value.False,
		"91f32 !~ 27i8": value.True,

		"1f32 !~ 1u64":   value.False,
		"1.5f32 !~ 1u64": value.True,
		"-8f32 !~ 8u64":  value.True,
		"91f32 !~ 27u64": value.True,

		"1f32 !~ 1u32":   value.False,
		"1.5f32 !~ 1u32": value.True,
		"-8f32 !~ 8u32":  value.True,
		"91f32 !~ 27u32": value.True,

		"53000f32 !~ 32767u16": value.True,
		"1f32 !~ 1u16":         value.False,
		"1.5f32 !~ 1u16":       value.True,
		"-8f32 !~ 8u16":        value.True,
		"91f32 !~ 27u16":       value.True,

		"256f32 !~ 127u8": value.True,
		"1f32 !~ 1u8":     value.False,
		"1.5f32 !~ 1u8":   value.True,
		"-8f32 !~ 8u8":    value.True,
		"91f32 !~ 27u8":   value.True,
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
		"'25' == '25'":   value.True,
		"'25' == '25.0'": value.False,
		"'25' == '7'":    value.False,

		"'7' == `7`": value.False,

		"'25' == 25.0": value.False,

		"'25' == 25bf": value.False,

		"'25' == 25f64": value.False,

		"'25' == 25f32": value.False,

		"'1' == 1i64": value.False,

		"'5' == 5i32": value.False,

		"'5' == 5i16": value.False,

		"'5' == 5i8": value.False,

		"'1' == 1u64": value.False,

		"'5' == 5u32": value.False,

		"'5' == 5u16": value.False,

		"'5' == 5u8": value.False,

		// Char
		"`2` == 25": value.False,

		"`2` == '2'": value.False,

		"`7` == `7`": value.True,
		"`b` == `b`": value.True,
		"`c` == `g`": value.False,
		"`7` == `8`": value.False,

		"`2` == 2.0": value.False,

		"`3` == 3bf": value.False,

		"`9` == 9f64": value.False,

		"`1` == 1f32": value.False,

		"`1` == 1i64": value.False,

		"`5` == 5i32": value.False,

		"`5` == 5i16": value.False,

		"`5` == 5i8": value.False,

		"`1` == 1u64": value.False,

		"`5` == 5u32": value.False,

		"`5` == 5u16": value.False,

		"`5` == 5u8": value.False,

		// Int
		"25 == 25":  value.True,
		"-25 == 25": value.False,
		"25 == -25": value.False,
		"25 == 28":  value.False,
		"28 == 25":  value.False,

		"25 == '25'": value.False,

		"7 == `7`": value.False,

		"25 == 25.0": value.False,

		"25 == 25bf": value.False,

		"25 == 25f64": value.False,

		"25 == 25f32": value.False,

		"1 == 1i64": value.False,

		"5 == 5i32": value.False,

		"5 == 5i16": value.False,

		"5 == 5i8": value.False,

		"1 == 1u64": value.False,

		"5 == 5u32": value.False,

		"5 == 5u16": value.False,

		"5 == 5u8": value.False,

		// Int64
		"25i64 == 25": value.False,

		"25i64 == '25'": value.False,

		"7i64 == `7`": value.False,

		"25i64 == 25.0": value.False,

		"25i64 == 25bf": value.False,

		"25i64 == 25f64": value.False,

		"25i64 == 25f32": value.False,

		"1i64 == 1i64":   value.True,
		"4i64 == -4i64":  value.False,
		"-8i64 == 8i64":  value.False,
		"-8i64 == -8i64": value.True,
		"91i64 == 27i64": value.False,

		"5i64 == 5i32": value.False,

		"5i64 == 5i16": value.False,

		"5i64 == 5i8": value.False,

		"1i64 == 1u64": value.False,

		"5i64 == 5u32": value.False,

		"5i64 == 5u16": value.False,

		"5i64 == 5u8": value.False,

		// Int32
		"25i32 == 25": value.False,

		"25i32 == '25'": value.False,

		"7i32 == `7`": value.False,

		"25i32 == 25.0": value.False,

		"25i32 == 25bf": value.False,

		"25i32 == 25f64": value.False,

		"25i32 == 25f32": value.False,

		"1i32 == 1i64": value.False,

		"5i32 == 5i32":  value.True,
		"4i32 == -4i32": value.False,
		"-8i32 == 8i32": value.False,
		"3i32 == 71i32": value.False,

		"5i32 == 5i16": value.False,

		"5i32 == 5i8": value.False,

		"1i32 == 1u64": value.False,

		"5i32 == 5u32": value.False,

		"5i32 == 5u16": value.False,

		"5i32 == 5u8": value.False,

		// Int16
		"25i16 == 25": value.False,

		"25i16 == '25'": value.False,

		"7i16 == `7`": value.False,

		"25i16 == 25.0": value.False,

		"25i16 == 25bf": value.False,

		"25i16 == 25f64": value.False,

		"25i16 == 25f32": value.False,

		"1i16 == 1i64": value.False,

		"5i16 == 5i32": value.False,

		"5i16 == 5i16":  value.True,
		"4i16 == -4i16": value.False,
		"-8i16 == 8i16": value.False,
		"3i16 == 71i16": value.False,

		"5i16 == 5i8": value.False,

		"1i16 == 1u64": value.False,

		"5i16 == 5u32": value.False,

		"5i16 == 5u16": value.False,

		"5i16 == 5u8": value.False,

		// Int8
		"25i8 == 25": value.False,

		"25i8 == '25'": value.False,

		"7i8 == `7`": value.False,

		"25i8 == 25.0": value.False,

		"25i8 == 25bf": value.False,

		"25i8 == 25f64": value.False,

		"25i8 == 25f32": value.False,

		"1i8 == 1i64": value.False,

		"5i8 == 5i32": value.False,

		"5i8 == 5i16": value.False,

		"5i8 == 5i8":  value.True,
		"4i8 == -4i8": value.False,
		"-8i8 == 8i8": value.False,
		"3i8 == 71i8": value.False,

		"1i8 == 1u64": value.False,

		"5i8 == 5u32": value.False,

		"5i8 == 5u16": value.False,

		"5i8 == 5u8": value.False,

		// UInt64
		"25u64 == 25": value.False,

		"25u64 == '25'": value.False,

		"7u64 == `7`": value.False,

		"25u64 == 25.0": value.False,

		"25u64 == 25bf": value.False,

		"25u64 == 25f64": value.False,

		"25u64 == 25f32": value.False,

		"1u64 == 1i64": value.False,

		"5u64 == 5i32": value.False,

		"5u64 == 5i16": value.False,

		"5u64 == 5i8": value.False,

		"1u64 == 1u64":   value.True,
		"91u64 == 27u64": value.False,

		"5u64 == 5u32": value.False,

		"5u64 == 5u16": value.False,

		"5u64 == 5u8": value.False,

		// UInt32
		"25u32 == 25": value.False,

		"25u32 == '25'": value.False,

		"7u32 == `7`": value.False,

		"25u32 == 25.0": value.False,

		"25u32 == 25bf": value.False,

		"25u32 == 25f64": value.False,

		"25u32 == 25f32": value.False,

		"1u32 == 1i64": value.False,

		"5u32 == 5i32": value.False,

		"5u32 == 5i16": value.False,

		"5u32 == 5i8": value.False,

		"1u32 == 1u64": value.False,

		"5u32 == 5u32":  value.True,
		"3u32 == 71u32": value.False,

		"5u32 == 5u16": value.False,

		"5u32 == 5u8": value.False,

		// UInt16
		"25u16 == 25": value.False,

		"25u16 == '25'": value.False,

		"7u16 == `7`": value.False,

		"25u16 == 25.0": value.False,

		"25u16 == 25bf": value.False,

		"25u16 == 25f64": value.False,

		"25u16 == 25f32": value.False,

		"1u16 == 1i64": value.False,

		"5u16 == 5i32": value.False,

		"5u16 == 5i16": value.False,

		"5u16 == 5i8": value.False,

		"1u16 == 1u64": value.False,

		"5u16 == 5u32": value.False,

		"53000u16 == 32767u16": value.False,
		"5u16 == 5u16":         value.True,
		"3u16 == 71u16":        value.False,

		"5u16 == 5u8": value.False,

		// UInt8
		"25u8 == 25": value.False,

		"25u8 == '25'": value.False,

		"7u8 == `7`": value.False,

		"25u8 == 25.0": value.False,

		"25u8 == 25bf": value.False,

		"25u8 == 25f64": value.False,

		"25u8 == 25f32": value.False,

		"1u8 == 1i64": value.False,

		"5u8 == 5i32": value.False,

		"5u8 == 5i16": value.False,

		"5u8 == 5i8": value.False,

		"1u8 == 1u64": value.False,

		"5u8 == 5u32": value.False,

		"5u8 == 5u16": value.False,

		"5u8 == 5u8":  value.True,
		"3u8 == 71u8": value.False,

		// Float
		"-73.0 == 73.0": value.False,
		"73.0 == -73.0": value.False,
		"25.0 == 25.0":  value.True,
		"1.0 == 1.2":    value.False,
		"1.2 == 1.0":    value.False,
		"78.5 == 78.5":  value.True,

		"8.25 == '8.25'": value.False,

		"4.0 == `4`": value.False,

		"25.0 == 25": value.False,

		"25.0 == 25bf":   value.False,
		"15.5 == 15.5bf": value.False,

		"25.0 == 25f64":     value.False,
		"15.26 == 15.26f64": value.False,

		"25.0 == 25f32":   value.False,
		"15.5 == 15.5f32": value.False,

		"1.0 == 1i64":   value.False,
		"-8.0 == -8i64": value.False,

		"1.0 == 1i32":   value.False,
		"-8.0 == -8i32": value.False,

		"1.0 == 1i16":   value.False,
		"-8.0 == -8i16": value.False,

		"1.0 == 1i8":   value.False,
		"-8.0 == -8i8": value.False,

		"1.0 == 1u64": value.False,

		"1.0 == 1u32": value.False,

		"1.0 == 1u16": value.False,

		"1.0 == 1u8": value.False,

		// Float64
		"25f64 == 25.0":   value.False,
		"78.5f64 == 78.5": value.False,

		"8.25f64 == '8.25'": value.False,

		"4f64 == `4`": value.False,

		"25f64 == 25": value.False,

		"25f64 == 25bf":     value.False,
		"15.5f64 == 15.5bf": value.False,

		"-73f64 == 73f64":      value.False,
		"73f64 == -73f64":      value.False,
		"25f64 == 25f64":       value.True,
		"1f64 == 1.2f64":       value.False,
		"15.26f64 == 15.26f64": value.True,

		"25f64 == 25f32":     value.False,
		"15.5f64 == 15.5f32": value.False,

		"1f64 == 1i64":   value.False,
		"-8f64 == -8i64": value.False,

		"1f64 == 1i32":   value.False,
		"-8f64 == -8i32": value.False,

		"1f64 == 1i16":   value.False,
		"-8f64 == -8i16": value.False,

		"1f64 == 1i8":   value.False,
		"-8f64 == -8i8": value.False,

		"1f64 == 1u64": value.False,

		"1f64 == 1u32": value.False,

		"1f64 == 1u16": value.False,

		"1f64 == 1u8": value.False,

		// Float32
		"25f32 == 25.0":   value.False,
		"78.5f32 == 78.5": value.False,

		"8.25f32 == '8.25'": value.False,

		"4f32 == `4`": value.False,

		"25f32 == 25": value.False,

		"25f32 == 25bf":     value.False,
		"15.5f32 == 15.5bf": value.False,

		"25f32 == 25f64":     value.False,
		"15.5f32 == 15.5f64": value.False,

		"-73f32 == 73f32":    value.False,
		"73f32 == -73f32":    value.False,
		"25f32 == 25f32":     value.True,
		"1f32 == 1.2f32":     value.False,
		"15.5f32 == 15.5f32": value.True,

		"1f32 == 1i64":   value.False,
		"-8f32 == -8i64": value.False,

		"1f32 == 1i32":   value.False,
		"-8f32 == -8i32": value.False,

		"1f32 == 1i16":   value.False,
		"-8f32 == -8i16": value.False,

		"1f32 == 1i8":   value.False,
		"-8f32 == -8i8": value.False,

		"1f32 == 1u64": value.False,

		"1f32 == 1u32": value.False,

		"1f32 == 1u16": value.False,

		"1f32 == 1u8": value.False,
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
		"'25' != '25'":   value.False,
		"'25' != '25.0'": value.True,
		"'25' != '7'":    value.True,

		"'7' != `7`": value.True,

		"'25' != 25.0": value.True,

		"'25' != 25bf": value.True,

		"'25' != 25f64": value.True,

		"'25' != 25f32": value.True,

		"'1' != 1i64": value.True,

		"'5' != 5i32": value.True,

		"'5' != 5i16": value.True,

		"'5' != 5i8": value.True,

		"'1' != 1u64": value.True,

		"'5' != 5u32": value.True,

		"'5' != 5u16": value.True,

		"'5' != 5u8": value.True,

		// Char
		"`2` != 25": value.True,

		"`2` != '2'": value.True,

		"`7` != `7`": value.False,
		"`b` != `b`": value.False,
		"`c` != `g`": value.True,
		"`7` != `8`": value.True,

		"`2` != 2.0": value.True,

		"`3` != 3bf": value.True,

		"`9` != 9f64": value.True,

		"`1` != 1f32": value.True,

		"`1` != 1i64": value.True,

		"`5` != 5i32": value.True,

		"`5` != 5i16": value.True,

		"`5` != 5i8": value.True,

		"`1` != 1u64": value.True,

		"`5` != 5u32": value.True,

		"`5` != 5u16": value.True,

		"`5` != 5u8": value.True,

		// Int
		"25 != 25":  value.False,
		"-25 != 25": value.True,
		"25 != -25": value.True,
		"25 != 28":  value.True,
		"28 != 25":  value.True,

		"25 != '25'": value.True,

		"7 != `7`": value.True,

		"25 != 25.0": value.True,

		"25 != 25bf": value.True,

		"25 != 25f64": value.True,

		"25 != 25f32": value.True,

		"1 != 1i64": value.True,

		"5 != 5i32": value.True,

		"5 != 5i16": value.True,

		"5 != 5i8": value.True,

		"1 != 1u64": value.True,

		"5 != 5u32": value.True,

		"5 != 5u16": value.True,

		"5 != 5u8": value.True,

		// Int64
		"25i64 != 25": value.True,

		"25i64 != '25'": value.True,

		"7i64 != `7`": value.True,

		"25i64 != 25.0": value.True,

		"25i64 != 25bf": value.True,

		"25i64 != 25f64": value.True,

		"25i64 != 25f32": value.True,

		"1i64 != 1i64":   value.False,
		"4i64 != -4i64":  value.True,
		"-8i64 != 8i64":  value.True,
		"-8i64 != -8i64": value.False,
		"91i64 != 27i64": value.True,

		"5i64 != 5i32": value.True,

		"5i64 != 5i16": value.True,

		"5i64 != 5i8": value.True,

		"1i64 != 1u64": value.True,

		"5i64 != 5u32": value.True,

		"5i64 != 5u16": value.True,

		"5i64 != 5u8": value.True,

		// Int32
		"25i32 != 25": value.True,

		"25i32 != '25'": value.True,

		"7i32 != `7`": value.True,

		"25i32 != 25.0": value.True,

		"25i32 != 25bf": value.True,

		"25i32 != 25f64": value.True,

		"25i32 != 25f32": value.True,

		"1i32 != 1i64": value.True,

		"5i32 != 5i32":  value.False,
		"4i32 != -4i32": value.True,
		"-8i32 != 8i32": value.True,
		"3i32 != 71i32": value.True,

		"5i32 != 5i16": value.True,

		"5i32 != 5i8": value.True,

		"1i32 != 1u64": value.True,

		"5i32 != 5u32": value.True,

		"5i32 != 5u16": value.True,

		"5i32 != 5u8": value.True,

		// Int16
		"25i16 != 25": value.True,

		"25i16 != '25'": value.True,

		"7i16 != `7`": value.True,

		"25i16 != 25.0": value.True,

		"25i16 != 25bf": value.True,

		"25i16 != 25f64": value.True,

		"25i16 != 25f32": value.True,

		"1i16 != 1i64": value.True,

		"5i16 != 5i32": value.True,

		"5i16 != 5i16":  value.False,
		"4i16 != -4i16": value.True,
		"-8i16 != 8i16": value.True,
		"3i16 != 71i16": value.True,

		"5i16 != 5i8": value.True,

		"1i16 != 1u64": value.True,

		"5i16 != 5u32": value.True,

		"5i16 != 5u16": value.True,

		"5i16 != 5u8": value.True,

		// Int8
		"25i8 != 25": value.True,

		"25i8 != '25'": value.True,

		"7i8 != `7`": value.True,

		"25i8 != 25.0": value.True,

		"25i8 != 25bf": value.True,

		"25i8 != 25f64": value.True,

		"25i8 != 25f32": value.True,

		"1i8 != 1i64": value.True,

		"5i8 != 5i32": value.True,

		"5i8 != 5i16": value.True,

		"5i8 != 5i8":  value.False,
		"4i8 != -4i8": value.True,
		"-8i8 != 8i8": value.True,
		"3i8 != 71i8": value.True,

		"1i8 != 1u64": value.True,

		"5i8 != 5u32": value.True,

		"5i8 != 5u16": value.True,

		"5i8 != 5u8": value.True,

		// UInt64
		"25u64 != 25": value.True,

		"25u64 != '25'": value.True,

		"7u64 != `7`": value.True,

		"25u64 != 25.0": value.True,

		"25u64 != 25bf": value.True,

		"25u64 != 25f64": value.True,

		"25u64 != 25f32": value.True,

		"1u64 != 1i64": value.True,

		"5u64 != 5i32": value.True,

		"5u64 != 5i16": value.True,

		"5u64 != 5i8": value.True,

		"1u64 != 1u64":   value.False,
		"91u64 != 27u64": value.True,

		"5u64 != 5u32": value.True,

		"5u64 != 5u16": value.True,

		"5u64 != 5u8": value.True,

		// UInt32
		"25u32 != 25": value.True,

		"25u32 != '25'": value.True,

		"7u32 != `7`": value.True,

		"25u32 != 25.0": value.True,

		"25u32 != 25bf": value.True,

		"25u32 != 25f64": value.True,

		"25u32 != 25f32": value.True,

		"1u32 != 1i64": value.True,

		"5u32 != 5i32": value.True,

		"5u32 != 5i16": value.True,

		"5u32 != 5i8": value.True,

		"1u32 != 1u64": value.True,

		"5u32 != 5u32":  value.False,
		"3u32 != 71u32": value.True,

		"5u32 != 5u16": value.True,

		"5u32 != 5u8": value.True,

		// UInt16
		"25u16 != 25": value.True,

		"25u16 != '25'": value.True,

		"7u16 != `7`": value.True,

		"25u16 != 25.0": value.True,

		"25u16 != 25bf": value.True,

		"25u16 != 25f64": value.True,

		"25u16 != 25f32": value.True,

		"1u16 != 1i64": value.True,

		"5u16 != 5i32": value.True,

		"5u16 != 5i16": value.True,

		"5u16 != 5i8": value.True,

		"1u16 != 1u64": value.True,

		"5u16 != 5u32": value.True,

		"53000u16 != 32767u16": value.True,
		"5u16 != 5u16":         value.False,
		"3u16 != 71u16":        value.True,

		"5u16 != 5u8": value.True,

		// UInt8
		"25u8 != 25": value.True,

		"25u8 != '25'": value.True,

		"7u8 != `7`": value.True,

		"25u8 != 25.0": value.True,

		"25u8 != 25bf": value.True,

		"25u8 != 25f64": value.True,

		"25u8 != 25f32": value.True,

		"1u8 != 1i64": value.True,

		"5u8 != 5i32": value.True,

		"5u8 != 5i16": value.True,

		"5u8 != 5i8": value.True,

		"1u8 != 1u64": value.True,

		"5u8 != 5u32": value.True,

		"5u8 != 5u16": value.True,

		"5u8 != 5u8":  value.False,
		"3u8 != 71u8": value.True,

		// Float
		"-73.0 != 73.0": value.True,
		"73.0 != -73.0": value.True,
		"25.0 != 25.0":  value.False,
		"1.0 != 1.2":    value.True,
		"1.2 != 1.0":    value.True,
		"78.5 != 78.5":  value.False,

		"8.25 != '8.25'": value.True,

		"4.0 != `4`": value.True,

		"25.0 != 25": value.True,

		"25.0 != 25bf":   value.True,
		"15.5 != 15.5bf": value.True,

		"25.0 != 25f64":     value.True,
		"15.26 != 15.26f64": value.True,

		"25.0 != 25f32":   value.True,
		"15.5 != 15.5f32": value.True,

		"1.0 != 1i64":   value.True,
		"-8.0 != -8i64": value.True,

		"1.0 != 1i32":   value.True,
		"-8.0 != -8i32": value.True,

		"1.0 != 1i16":   value.True,
		"-8.0 != -8i16": value.True,

		"1.0 != 1i8":   value.True,
		"-8.0 != -8i8": value.True,

		"1.0 != 1u64": value.True,

		"1.0 != 1u32": value.True,

		"1.0 != 1u16": value.True,

		"1.0 != 1u8": value.True,

		// Float64
		"25f64 != 25.0":   value.True,
		"78.5f64 != 78.5": value.True,

		"8.25f64 != '8.25'": value.True,

		"4f64 != `4`": value.True,

		"25f64 != 25": value.True,

		"25f64 != 25bf":     value.True,
		"15.5f64 != 15.5bf": value.True,

		"-73f64 != 73f64":      value.True,
		"73f64 != -73f64":      value.True,
		"25f64 != 25f64":       value.False,
		"1f64 != 1.2f64":       value.True,
		"15.26f64 != 15.26f64": value.False,

		"25f64 != 25f32":     value.True,
		"15.5f64 != 15.5f32": value.True,

		"1f64 != 1i64":   value.True,
		"-8f64 != -8i64": value.True,

		"1f64 != 1i32":   value.True,
		"-8f64 != -8i32": value.True,

		"1f64 != 1i16":   value.True,
		"-8f64 != -8i16": value.True,

		"1f64 != 1i8":   value.True,
		"-8f64 != -8i8": value.True,

		"1f64 != 1u64": value.True,

		"1f64 != 1u32": value.True,

		"1f64 != 1u16": value.True,

		"1f64 != 1u8": value.True,

		// Float32
		"25f32 != 25.0":   value.True,
		"78.5f32 != 78.5": value.True,

		"8.25f32 != '8.25'": value.True,

		"4f32 != `4`": value.True,

		"25f32 != 25": value.True,

		"25f32 != 25bf":     value.True,
		"15.5f32 != 15.5bf": value.True,

		"25f32 != 25f64":     value.True,
		"15.5f32 != 15.5f64": value.True,

		"-73f32 != 73f32":    value.True,
		"73f32 != -73f32":    value.True,
		"25f32 != 25f32":     value.False,
		"1f32 != 1.2f32":     value.True,
		"15.5f32 != 15.5f32": value.False,

		"1f32 != 1i64":   value.True,
		"-8f32 != -8i64": value.True,

		"1f32 != 1i32":   value.True,
		"-8f32 != -8i32": value.True,

		"1f32 != 1i16":   value.True,
		"-8f32 != -8i16": value.True,

		"1f32 != 1i8":   value.True,
		"-8f32 != -8i8": value.True,

		"1f32 != 1u64": value.True,

		"1f32 != 1u32": value.True,

		"1f32 != 1u16": value.True,

		"1f32 != 1u8": value.True,
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}
