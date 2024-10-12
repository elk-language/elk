package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestPatterns(t *testing.T) {
	tests := testTable{
		"public identifier pattern": {
			input: `
				var [a] = [1]
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 16), P(34, 3, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"redeclare public identifier pattern": {
			input: `
				var a: Int
				var [a] = [1]
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(49, 4, 16), P(49, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"redeclare public identifier with a different type": {
			input: `
				var a: String
				var [a] = [1]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 10), P(28, 3, 10)), "type `Std::Int` cannot be assigned to type `Std::String`"),
			},
		},
		"private identifier pattern": {
			input: `
				var [_a] = ["", 8]
				var b: 9 = _a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 3, 16), P(40, 3, 17)), "type `Std::String | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"list pattern with rest and literal": {
			input: `
				var [a, *b] = ["", 8, 1]
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 16), P(45, 3, 16)), "type `Std::ArrayList[Std::String | Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"list pattern with rest and wide type": {
			input: `
				var d: ArrayList[String] | List[Int] | nil = nil
				var [a, *b] = d
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(89, 4, 16), P(89, 4, 16)), "type `Std::ArrayList[Std::String | Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"list pattern with invalid value": {
			input: `
				var [a] = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "type `8` cannot ever match type `Std::List[any]`"),
			},
		},
		"list pattern with a wider type declares a variable with element type": {
			input: `
				var a: List[Int] | nil = nil
				var [b] = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(65, 4, 16), P(65, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"list pattern with a wider type and incompatible values": {
			input: `
				var a: List[Int] | nil = nil
				var ["foo", b] = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 3, 10), P(47, 3, 14)), "type `Std::Int` cannot ever match type `\"foo\"`"),
			},
		},
		"list pattern with a wider type with ArrayList declares a variable with element type": {
			input: `
				var a: ArrayList[Int] | nil = nil
				var [b] = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(70, 4, 16), P(70, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"list pattern with a wider incompatible element type with ArrayList": {
			input: `
				var a: ArrayList[Int] | nil = nil
				var ["foo", 2.5, b] = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 3, 10), P(52, 3, 14)), "type `Std::Int` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(55, 3, 17), P(57, 3, 19)), "type `Std::Int` cannot ever match type `2.5`"),
			},
		},
		"list pattern with a wider type ArrayList | List declares a variable with element type": {
			input: `
				var a: ArrayList[Int] | List[Float] | nil = nil
				var [b] = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(84, 4, 16), P(84, 4, 16)), "type `Std::Int | Std::Float` cannot be assigned to type `9`"),
			},
		},
		"list pattern with a wider incompatible type ArrayList | List": {
			input: `
				var a: ArrayList[Int] | List[Float] | nil = nil
				var ["foo", 9i8, b] = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(62, 3, 10), P(66, 3, 14)), "type `Std::Int | Std::Float` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(69, 3, 17), P(71, 3, 19)), "type `Std::Int | Std::Float` cannot ever match type `9i8`"),
			},
		},

		"tuple pattern with a wider type declares a variable with element type": {
			input: `
				var a: Tuple[Int] | nil = nil
				var %[b] = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(67, 4, 16), P(67, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"tuple pattern with a wider type and incompatible values": {
			input: `
				var a: Tuple[Int] | nil = nil
				var %["foo", b] = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 11), P(49, 3, 15)), "type `Std::Int` cannot ever match type `\"foo\"`"),
			},
		},
		"tuple pattern with a wider type with ArrayTuple declares a variable with element type": {
			input: `
				var a: ArrayTuple[Int] | nil = nil
				var %[b] = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(72, 4, 16), P(72, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"tuple pattern with a wider type with ArrayList declares a variable with element type": {
			input: `
				var a: ArrayList[Int] | nil = nil
				var %[b] = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(71, 4, 16), P(71, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"tuple pattern with a wider incompatible element type with ArrayTuple": {
			input: `
				var a: ArrayTuple[Int] | nil = nil
				var %["foo", 2.5, b] = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(50, 3, 11), P(54, 3, 15)), "type `Std::Int` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(57, 3, 18), P(59, 3, 20)), "type `Std::Int` cannot ever match type `2.5`"),
			},
		},
		"tuple pattern with a wider type ArrayTuple | List declares a variable with element type": {
			input: `
				var a: ArrayTuple[Int] | List[Float] | nil = nil
				var %[b] = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(86, 4, 16), P(86, 4, 16)), "type `Std::Int | Std::Float` cannot be assigned to type `9`"),
			},
		},
		"tuple pattern with a wider incompatible type ArrayList | Tuple": {
			input: `
				var a: ArrayList[Int] | Tuple[Float] | nil = nil
				var %["foo", 9i8, b] = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(64, 3, 11), P(68, 3, 15)), "type `Std::Int | Std::Float` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(71, 3, 18), P(73, 3, 20)), "type `Std::Int | Std::Float` cannot ever match type `9i8`"),
			},
		},
		"tuple pattern with rest": {
			input: `
				var %[a, *b] = ["", 8, 1]
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 3, 16), P(46, 3, 16)), "type `Std::ArrayList[Std::String | Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"tuple pattern with invalid value": {
			input: `
				var %[a] = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `8` cannot ever match type `Std::Tuple[any]`"),
			},
		},
		"pattern with as and new variable": {
			input: `
				var %[a as b] = [8]
				a = b
			`,
		},
		"pattern with as and existing variable": {
			input: `
				var b: String
				var %[a as b] = [8]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 16), P(34, 3, 16)), "type `Std::Int` cannot be assigned to type `Std::String`"),
			},
		},

		"pattern with as has the patterns type": {
			input: `
				var b: Int | Float | nil = 3
				var 1 as a = b
				var c: nil = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(70, 4, 18), P(70, 4, 18)), "type `Std::Int` cannot be assigned to type `nil`"),
			},
		},

		"pattern with int literal and Int type": {
			input: `
				b := 3
				var 1 as a = b
			`,
		},
		"pattern with int literal and wrong literal type": {
			input: `
				var 1 as a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(9, 2, 9)), "type `3` cannot ever match type `1`"),
			},
		},
		"pattern with int literal and wider type": {
			input: `
				var b: String | Int = 3
				var 1 as a = b
			`,
		},
		"pattern with int literal and wrong type": {
			input: `
				var 1 as a = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(9, 2, 9)), "type `1.2` cannot ever match type `1`"),
			},
		},

		"pattern with int64 literal and Int64 type": {
			input: `
				b := 3i64
				var 1i64 as a = b
			`,
		},
		"pattern with int64 literal and wrong literal type": {
			input: `
				var 1i64 as a = 3i64
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `3i64` cannot ever match type `1i64`"),
			},
		},
		"pattern with int64 literal and wider type": {
			input: `
				var b: String | Int64 = 3i64
				var 1i64 as a = b
			`,
		},
		"pattern with int64 literal and wrong type": {
			input: `
				var 1i64 as a = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `1.2` cannot ever match type `1i64`"),
			},
		},
		"pattern with int32 literal and Int32 type": {
			input: `
				b := 3i32
				var 1i32 as a = b
			`,
		},
		"pattern with int32 literal and wrong literal type": {
			input: `
				var 1i32 as a = 3i32
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `3i32` cannot ever match type `1i32`"),
			},
		},
		"pattern with int32 literal and wider type": {
			input: `
				var b: String | Int32 = 3i32
				var 1i32 as a = b
			`,
		},
		"pattern with int32 literal and wrong type": {
			input: `
				var 1i32 as a = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `1.2` cannot ever match type `1i32`"),
			},
		},
		"pattern with int16 literal and Int16 type": {
			input: `
				b := 3i16
				var 1i16 as a = b
			`,
		},
		"pattern with int16 literal and wrong literal type": {
			input: `
				var 1i16 as a = 3i16
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `3i16` cannot ever match type `1i16`"),
			},
		},
		"pattern with int16 literal and wider type": {
			input: `
				var b: String | Int16 = 3i16
				var 1i16 as a = b
			`,
		},
		"pattern with int16 literal and wrong type": {
			input: `
				var 1i16 as a = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `1.2` cannot ever match type `1i16`"),
			},
		},
		"pattern with int8 literal and Int8 type": {
			input: `
				b := 3i8
				var 1i8 as a = b
			`,
		},
		"pattern with int8 literal and wrong literal type": {
			input: `
				var 1i8 as a = 3i8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "type `3i8` cannot ever match type `1i8`"),
			},
		},
		"pattern with int8 literal and wider type": {
			input: `
				var b: String | Int8 = 3i8
				var 1i8 as a = b
			`,
		},
		"pattern with int8 literal and wrong type": {
			input: `
				var 1i8 as a = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "type `1.2` cannot ever match type `1i8`"),
			},
		},
		"pattern with uint64 literal and UInt64 type": {
			input: `
				b := 3u64
				var 1u64 as a = b
			`,
		},
		"pattern with uint64 literal and wrong literal type": {
			input: `
				var 1u64 as a = 3u64
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `3u64` cannot ever match type `1u64`"),
			},
		},
		"pattern with uint64 literal and wider type": {
			input: `
				var b: String | UInt64 = 3u64
				var 1u64 as a = b
			`,
		},
		"pattern with uint64 literal and wrong type": {
			input: `
				var 1u64 as a = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `1.2` cannot ever match type `1u64`"),
			},
		},
		"pattern with uint32 literal and UInt32 type": {
			input: `
				b := 3u32
				var 1u32 as a = b
			`,
		},
		"pattern with uint32 literal and wrong literal type": {
			input: `
				var 1u32 as a = 3u32
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `3u32` cannot ever match type `1u32`"),
			},
		},
		"pattern with uint32 literal and wider type": {
			input: `
				var b: String | UInt32 = 3u32
				var 1u32 as a = b
			`,
		},
		"pattern with uint32 literal and wrong type": {
			input: `
				var 1u32 as a = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `1.2` cannot ever match type `1u32`"),
			},
		},
		"pattern with uint16 literal and UInt16 type": {
			input: `
				b := 3u16
				var 1u16 as a = b
			`,
		},
		"pattern with uint16 literal and wrong literal type": {
			input: `
				var 1u16 as a = 3u16
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `3u16` cannot ever match type `1u16`"),
			},
		},
		"pattern with uint16 literal and wider type": {
			input: `
				var b: String | UInt16 = 3u16
				var 1u16 as a = b
			`,
		},
		"pattern with uint16 literal and wrong type": {
			input: `
				var 1u16 as a = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `1.2` cannot ever match type `1u16`"),
			},
		},
		"pattern with uint8 literal and UInt8 type": {
			input: `
				b := 3u8
				var 1u8 as a = b
			`,
		},
		"pattern with uint8 literal and wrong literal type": {
			input: `
				var 1u8 as a = 3u8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "type `3u8` cannot ever match type `1u8`"),
			},
		},
		"pattern with uint8 literal and wider type": {
			input: `
				var b: String | UInt8 = 3u8
				var 1u8 as a = b
			`,
		},
		"pattern with uint8 literal and wrong type": {
			input: `
				var 1u8 as a = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "type `1.2` cannot ever match type `1u8`"),
			},
		},

		"pattern with float literal and Float type": {
			input: `
				b := 3.14
				var 1.0 as a = b
			`,
		},
		"pattern with float literal and wrong literal type": {
			input: `
				var 1.0 as a = 3.14
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "type `3.14` cannot ever match type `1.0`"),
			},
		},
		"pattern with float literal and wider type": {
			input: `
				var b: String | Float = 3.14
				var 1.0 as a = b
			`,
		},
		"pattern with float literal and wrong type": {
			input: `
				var 1.0 as a = "1.2"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "type `\"1.2\"` cannot ever match type `1.0`"),
			},
		},
		"pattern with float64 literal and Float64 type": {
			input: `
				b := 3.14f64
				var 1.0f64 as a = b
			`,
		},
		"pattern with float64 literal and wrong literal type": {
			input: `
				var 1.0f64 as a = 3.14f64
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(14, 2, 14)), "type `3.14f64` cannot ever match type `1.0f64`"),
			},
		},
		"pattern with float64 literal and wider type": {
			input: `
				var b: String | Float64 = 3.14f64
				var 1.0f64 as a = b
			`,
		},
		"pattern with float64 literal and wrong type": {
			input: `
				var 1.0f64 as a = "1.2"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(14, 2, 14)), "type `\"1.2\"` cannot ever match type `1.0f64`"),
			},
		},
		"pattern with float32 literal and Float32 type": {
			input: `
				b := 3.14f32
				var 1.0f32 as a = b
			`,
		},
		"pattern with float32 literal and wrong literal type": {
			input: `
				var 1.0f32 as a = 3.14f32
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(14, 2, 14)), "type `3.14f32` cannot ever match type `1.0f32`"),
			},
		},
		"pattern with float32 literal and wider type": {
			input: `
				var b: String | Float32 = 3.14f32
				var 1.0f32 as a = b
			`,
		},
		"pattern with float32 literal and wrong type": {
			input: `
				var 1.0f32 as a = "1.2"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(14, 2, 14)), "type `\"1.2\"` cannot ever match type `1.0f32`"),
			},
		},
		"pattern with bigfloat literal and BigFloat type": {
			input: `
				b := 3.14bf
				var 1.0bf as a = b
			`,
		},
		"pattern with bigfloat literal and wrong literal type": {
			input: `
				var 1.0bf as a = 3.14bf
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(13, 2, 13)), "type `3.14bf` cannot ever match type `1.0bf`"),
			},
		},
		"pattern with bigfloat literal and wider type": {
			input: `
				var b: String | BigFloat = 3.14bf
				var 1.0bf as a = b
			`,
		},
		"pattern with bigfloat literal and wrong type": {
			input: `
				var 1.0bf as a = "1.2"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(13, 2, 13)), "type `\"1.2\"` cannot ever match type `1.0bf`"),
			},
		},

		"pattern with char literal and Char type": {
			input: "b := `a`\nvar `a` as a = b",
		},
		"pattern with char literal and wrong literal type": {
			input: "var `a` as a = `b`",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(6, 1, 7)), "type ``b`` cannot ever match type ``a``"),
			},
		},
		"pattern with char literal and wider type": {
			input: "var b: String | Char = `a`\nvar `a` as a = b",
		},
		"pattern with char literal and wrong type": {
			input: "var `a` as a = \"a\"",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(6, 1, 7)), "type `\"a\"` cannot ever match type ``a``"),
			},
		},
		"pattern with raw char literal and Char type": {
			input: "b := r`f`\nvar r`f` as a = b",
		},
		"pattern with raw char literal and wrong literal type": {
			input: "var r`f` as a = r`g`",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(7, 1, 8)), "type ``g`` cannot ever match type ``f``"),
			},
		},
		"pattern with raw char literal and wider type": {
			input: "var b: String | Char = r`f`\nvar r`f` as a = b",
		},
		"pattern with raw char literal and wrong type": {
			input: "var r`f` as a = \"f\"",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(7, 1, 8)), "type `\"f\"` cannot ever match type ``f``"),
			},
		},

		"pattern with true literal and Bool type": {
			input: `
				b := true
				var true as a = b
			`,
		},
		"pattern with true literal and wrong literal type": {
			input: `
				var true as a = false
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `false` cannot ever match type `true`"),
			},
		},
		"pattern with true literal and wider type": {
			input: `
				var b: String | Bool = true
				var true as a = b
			`,
		},
		"pattern with true literal and wrong type": {
			input: `
				var true as a = "true"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `\"true\"` cannot ever match type `true`"),
			},
		},
		"pattern with false literal and Bool type": {
			input: `
				b := false
				var false as a = b
			`,
		},
		"pattern with false literal and wrong literal type": {
			input: `
				var false as a = true
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(13, 2, 13)), "type `true` cannot ever match type `false`"),
			},
		},
		"pattern with false literal and wider type": {
			input: `
				var b: String | Bool = false
				var false as a = b
			`,
		},
		"pattern with false literal and wrong type": {
			input: `
				var false as a = "false"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(13, 2, 13)), "type `\"false\"` cannot ever match type `false`"),
			},
		},
		"pattern with nil literal and Nil type": {
			input: `
				b := nil
				var nil as a = b
			`,
		},
		"pattern with nil literal and wrong literal type": {
			input: `
				var nil as a = 42
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "type `42` cannot ever match type `nil`"),
			},
		},
		"pattern with nil literal and wider type": {
			input: `
				var b: String | Nil = nil
				var nil as a = b
			`,
		},
		"pattern with nil literal and wrong type": {
			input: `
				var nil as a = "nil"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "type `\"nil\"` cannot ever match type `nil`"),
			},
		},

		"pattern with interpolated string literal": {
			input: `
				b := "hello"
				var "hello${1 + 2i8}" as a = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(38, 3, 21), P(40, 3, 23)), "expected type `Std::Int` for parameter `other` in call to `+`, got type `2i8`"),
			},
		},
		"pattern with raw string literal": {
			input: `
				b := "hello"
				var 'foo' as a = b
			`,
		},
		"pattern with string literal and String type": {
			input: `
				b := "hello"
				var "hello" as a = b
			`,
		},
		"pattern with string literal and wrong literal type": {
			input: `
				var "hello" as a = "world"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(15, 2, 15)), "type `\"world\"` cannot ever match type `\"hello\"`"),
			},
		},
		"pattern with string literal and wider type": {
			input: `
				var b: String | Int = "hello"
				var "hello" as a = b
			`,
		},
		"pattern with string literal and wrong type": {
			input: `
				var "hello" as a = 42
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(15, 2, 15)), "type `42` cannot ever match type `\"hello\"`"),
			},
		},

		"pattern with interpolated symbol literal": {
			input: `
				b := :hello
				var :"hello${1 + 2i8}" as a = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(38, 3, 22), P(40, 3, 24)), "expected type `Std::Int` for parameter `other` in call to `+`, got type `2i8`"),
			},
		},
		"pattern with simple symbol literal": {
			input: `
				b := :hello
				var :foo as a = b
			`,
		},
		"pattern with symbol literal and Symbol type": {
			input: `
				b := :hello
				var :hello as a = b
			`,
		},
		"pattern with symbol literal and wrong literal type": {
			input: `
				var :hello as a = :world
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(14, 2, 14)), "type `:world` cannot ever match type `:hello`"),
			},
		},
		"pattern with symbol literal and wider type": {
			input: `
				var b: Symbol | Int = :hello
				var :hello as a = b
			`,
		},
		"pattern with symbol literal and wrong type": {
			input: `
				var :hello as a = 42
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(14, 2, 14)), "type `42` cannot ever match type `:hello`"),
			},
		},

		"pattern with word list literal and array list type": {
			input: `
				b := [1, "foo"]
				var \w[foo bar] as a = b
			`,
		},
		"pattern with word list literal and wrong literal type": {
			input: `
				var \w[foo bar] as a = [1]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(19, 2, 19)), "type `Std::ArrayList[Std::Int]` cannot ever match type `Std::List[Std::String]`"),
			},
		},
		"pattern with word list literal and wider type": {
			input: `
				var b: String | List[String] = ""
				var \w[foo bar] as a = b
			`,
		},
		"pattern with word list literal and wrong type": {
			input: `
				var \w[foo bar] as a = 1
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(19, 2, 19)), "type `1` cannot ever match type `Std::List[Std::String]`"),
			},
		},
		"pattern with symbol list literal and array list type": {
			input: `
				b := [:foo, :bar]
				var \s[foo bar] as a = b
			`,
		},
		"pattern with symbol list literal and wrong literal type": {
			input: `
				var \s[foo bar] as a = ["foo"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(19, 2, 19)), "type `Std::ArrayList[Std::String]` cannot ever match type `Std::List[Std::Symbol]`"),
			},
		},
		"pattern with symbol list literal and wider type": {
			input: `
				var b: Symbol | List[Symbol] = :foo
				var \s[foo bar] as a = b
			`,
		},
		"pattern with symbol list literal and wrong type": {
			input: `
				var \s[foo bar] as a = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(19, 2, 19)), "type `\"foo\"` cannot ever match type `Std::List[Std::Symbol]`"),
			},
		},
		"pattern with binary list literal and array list type": {
			input: `
				b := [0b1010, 0b1100]
				var \b[1010 1100] as a = b
			`,
		},
		"pattern with binary list literal and wrong literal type": {
			input: `
				var \b[1010 1100] as a = ["foo"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(21, 2, 21)), "type `Std::ArrayList[Std::String]` cannot ever match type `Std::List[Std::Int]`"),
			},
		},
		"pattern with binary list literal and wider type": {
			input: `
				var b: Int | List[Int] = 0b1010
				var \b[1010 1100] as a = b
			`,
		},
		"pattern with binary list literal and wrong type": {
			input: `
				var \b[1010 1100] as a = "1010"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(21, 2, 21)), "type `\"1010\"` cannot ever match type `Std::List[Std::Int]`"),
			},
		},
		"pattern with hex list literal and array list type": {
			input: `
				b := [0xA, 0xB]
				var \x[A B] as a = b
			`,
		},
		"pattern with hex list literal and wrong literal type": {
			input: `
				var \x[A B] as a = ["foo"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(15, 2, 15)), "type `Std::ArrayList[Std::String]` cannot ever match type `Std::List[Std::Int]`"),
			},
		},
		"pattern with hex list literal and wider type": {
			input: `
				var b: Int | List[Int] = 0xA
				var \x[A B] as a = b
			`,
		},
		"pattern with hex list literal and wrong type": {
			input: `
				var \x[A B] as a = "A"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(15, 2, 15)), "type `\"A\"` cannot ever match type `Std::List[Std::Int]`"),
			},
		},

		"pattern with word tuple literal and tuple type": {
			input: `
				b := %[1, "foo"]
				var %w[foo bar] as a = b
			`,
		},
		"pattern with word tuple literal and wrong literal type": {
			input: `
				var %w[foo bar] as a = %[1, 2]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(19, 2, 19)), "type `Std::ArrayTuple[1 | 2]` cannot ever match type `Std::Tuple[Std::String]`"),
			},
		},
		"pattern with binary tuple literal and tuple type": {
			input: `
				b := %[0b1010, 0b1100]
				var %b[1010 1100] as a = b
			`,
		},
		"pattern with binary tuple literal and wrong literal type": {
			input: `
				var %b[1010 1100] as a = %["foo", "bar"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(21, 2, 21)), "type `Std::ArrayTuple[\"foo\" | \"bar\"]` cannot ever match type `Std::Tuple[Std::Int]`"),
			},
		},
		"pattern with hex tuple literal and tuple type": {
			input: `
				b := %[0xA, 0xB]
				var %x[A B] as a = b
			`,
		},
		"pattern with hex tuple literal and wrong literal type": {
			input: `
				var %x[A B] as a = %["foo", "bar"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(15, 2, 15)), "type `Std::ArrayTuple[\"foo\" | \"bar\"]` cannot ever match type `Std::Tuple[Std::Int]`"),
			},
		},
		"pattern with symbol tuple literal and tuple type": {
			input: `
				b := %[:foo, :bar]
				var %s[foo bar] as a = b
			`,
		},
		"pattern with symbol tuple literal and wrong literal type": {
			input: `
				var %s[foo bar] as a = %["foo", "bar"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(19, 2, 19)), "type `Std::ArrayTuple[\"foo\" | \"bar\"]` cannot ever match type `Std::Tuple[Std::Symbol]`"),
			},
		},
		"pattern with word set literal and set type": {
			input: `
				b := ^["foo", "bar"]
				var ^w[foo bar] as a = b
			`,
		},
		"pattern with word set literal and wrong literal type": {
			input: `
				var ^w[foo bar] as a = ^[1, 2]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(19, 2, 19)), "type `Std::HashSet[Std::Int]` cannot ever match type `Std::Set[Std::String]`"),
			},
		},
		"pattern with binary set literal and set type": {
			input: `
				b := ^[0b1010, 0b1100]
				var ^b[1010 1100] as a = b
			`,
		},
		"pattern with binary set literal and wrong literal type": {
			input: `
				var ^b[1010 1100] as a = ^["foo", "bar"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(21, 2, 21)), "type `Std::HashSet[Std::String]` cannot ever match type `Std::Set[Std::Int]`"),
			},
		},
		"pattern with hex set literal and set type": {
			input: `
				b := ^[0xA, 0xB]
				var ^x[A B] as a = b
			`,
		},
		"pattern with hex set literal and wrong literal type": {
			input: `
				var ^x[A B] as a = ^["foo", "bar"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(15, 2, 15)), "type `Std::HashSet[Std::String]` cannot ever match type `Std::Set[Std::Int]`"),
			},
		},
		"pattern with symbol set literal and set type": {
			input: `
				b := ^[:foo, :bar]
				var ^s[foo bar] as a = b
			`,
		},
		"pattern with symbol set literal and wrong literal type": {
			input: `
				var ^s[foo bar] as a = ^["foo", "bar"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(19, 2, 19)), "type `Std::HashSet[Std::String]` cannot ever match type `Std::Set[Std::Symbol]`"),
			},
		},

		"range pattern and wider type": {
			input: `
				var a: Int | Float | String | nil = nil
				var 1...9 as b = a
				var c: nil = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 4, 18), P(85, 4, 18)), "type `Std::Int` cannot be assigned to type `nil`"),
			},
		},
		"range pattern with two different literal types": {
			input: `
				var 1...5.9 as a = 9
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(15, 2, 15)), "range pattern start and end must be of the same type, got `Std::Int` and `Std::Float`"),
			},
		},
		"range pattern with wider start type": {
			input: `
				const A: Int | Float = 25.9
				var 2.5...A as b = 9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(47, 3, 15), P(47, 3, 15)), "type `Std::Int | Std::Float` cannot be used in a range pattern, only class instance types are permitted"),
				error.NewFailure(L("<main>", P(41, 3, 9), P(47, 3, 15)), "range pattern start and end must be of the same type, got `Std::Float` and `Std::Int | Std::Float`"),
			},
		},
		"range pattern and correct type": {
			input: `
				var 1...15 as a = 9
				var b: nil = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 3, 18), P(42, 3, 18)), "type `Std::Int` cannot be assigned to type `nil`"),
			},
		},
		"range pattern and wrong literal type": {
			input: `
				var 1...9 as a = 5i8
				a = nil
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(13, 2, 13)), "type `5i8` cannot ever match type `Std::Int`"),
				error.NewFailure(L("<main>", P(34, 3, 9), P(36, 3, 11)), "type `nil` cannot be assigned to type `Std::Int`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
