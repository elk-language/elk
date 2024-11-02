package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestIdentifierPattern(t *testing.T) {
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
		"declares a variable in a variable declaration": {
			input: `
				var [a] = ["", 8]
				a = 3
			`,
		},
		"declares a value in a value declaration": {
			input: `
				val [a] = ["", 8]
				a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 5), P(27, 3, 5)), "local value `a` cannot be reassigned"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSetPattern(t *testing.T) {
	tests := testTable{
		"set pattern with rest and wide type": {
			input: `
				var d: HashSet[String] | Set[Int] | nil = nil
				var ^[1, "foo"] as b = d
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(95, 4, 16), P(95, 4, 16)), "type `Std::HashSet[Std::String] | Std::Set[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"set pattern with invalid value": {
			input: `
				var ^[1] as a = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `\"foo\"` cannot ever match type `Std::Set[any]`"),
			},
		},
		"set pattern with a wider type declares a variable": {
			input: `
				var a: Set[Int] | nil = nil
				var ^[1] as b = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(70, 4, 16), P(70, 4, 16)), "type `Std::Set[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"set pattern with a wider type and incompatible values": {
			input: `
				var a: Set[Int] | nil = nil
				var ^["foo", 2.5] as b = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 3, 11), P(47, 3, 15)), "type `Std::Int` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(50, 3, 18), P(52, 3, 20)), "type `Std::Int` cannot ever match type `2.5`"),
			},
		},
		"set pattern with a wider type with HashSet declares a variable": {
			input: `
				var a: HashSet[Int] | nil = nil
				var ^[9, 7] as b = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(77, 4, 16), P(77, 4, 16)), "type `Std::HashSet[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"set pattern with a wider incompatible element type with HashSet": {
			input: `
				var a: HashSet[Int] | nil = nil
				var ^["foo", 2.5] as b = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(47, 3, 11), P(51, 3, 15)), "type `Std::Int` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(54, 3, 18), P(56, 3, 20)), "type `Std::Int` cannot ever match type `2.5`"),
			},
		},
		"set pattern with a wider type HashSet | Set declares a variable with element type": {
			input: `
				var a: HashSet[Int] | Set[Float] | nil = nil
				var ^[1, 2, 9.8, .1] as b = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(99, 4, 16), P(99, 4, 16)), "type `Std::HashSet[Std::Int] | Std::Set[Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"set pattern with a wider incompatible type HashSet | Set": {
			input: `
				var a: HashSet[Int] | Set[Float] | nil = nil
				var ^["foo", 9i8, 9] as b = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 3, 11), P(64, 3, 15)), "type `Std::Int | Std::Float` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(67, 3, 18), P(69, 3, 20)), "type `Std::Int | Std::Float` cannot ever match type `9i8`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestListPattern(t *testing.T) {
	tests := testTable{
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
				var [b] as c = a
				var d: 9 = b
				var e: nil = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(89, 4, 16), P(89, 4, 16)), "type `Std::Int | Std::Float` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(108, 5, 18), P(108, 5, 18)), "type `Std::ArrayList[Std::Int] | Std::List[Std::Float]` cannot be assigned to type `nil`"),
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestTuplePattern(t *testing.T) {
	tests := testTable{
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
				var %[b] as c = a
				var d: 9 = b
				var e: nil = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(91, 4, 16), P(91, 4, 16)), "type `Std::Int | Std::Float` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(110, 5, 18), P(110, 5, 18)), "type `Std::ArrayTuple[Std::Int] | Std::List[Std::Float]` cannot be assigned to type `nil`"),
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestAsPattern(t *testing.T) {
	tests := testTable{
		"pattern with as and new variable": {
			input: `
				var %[a as b] = [8]
				a = b
			`,
		},
		"declares a variable in a variable declaration": {
			input: `
				var 1 as a = 1
				a = 5
			`,
		},
		"declares a value in a value declaration": {
			input: `
				val 1 as a = 1
				a = 5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 3, 5), P(24, 3, 5)), "local value `a` cannot be reassigned"),
			},
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestIntPattern(t *testing.T) {
	tests := testTable{
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
		"pattern with negative int literal and wrong literal type": {
			input: `
				var -1 as a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(10, 2, 10)), "type `3` cannot ever match type `-1`"),
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestFloatPattern(t *testing.T) {
	tests := testTable{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestCharPattern(t *testing.T) {
	tests := testTable{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSimpleLiteralPattern(t *testing.T) {
	tests := testTable{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestStringLiteralPattern(t *testing.T) {
	tests := testTable{

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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSymbolLiteralPattern(t *testing.T) {
	tests := testTable{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSpecialListPattern(t *testing.T) {
	tests := testTable{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSpecialTuplePattern(t *testing.T) {
	tests := testTable{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSpecialSetPattern(t *testing.T) {
	tests := testTable{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestRangePattern(t *testing.T) {
	tests := testTable{
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

func TestMapPattern(t *testing.T) {
	tests := testTable{
		"declares variables in variable declaration": {
			input: `
				var { a } = { a: 3 }
				a = 5
			`,
		},
		"declares values in value declaration": {
			input: `
				val { a } = { a: 3 }
				a = 5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 3, 5), P(30, 3, 5)), "local value `a` cannot be reassigned"),
			},
		},
		"map pattern with invalid value": {
			input: `
				var { a } = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(13, 2, 13)), "type `8` cannot ever match type `Std::Map[any, any]`"),
			},
		},
		"map pattern with a wider type declares a variable with value type": {
			input: `
				var a: Map[Symbol, Int] | nil = nil
				var { b } = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(74, 4, 16), P(74, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"map pattern with a wider type declares a variable with map type": {
			input: `
				var a: Map[Symbol, Int] | nil = nil
				var { g: 8 } as b = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(82, 4, 16), P(82, 4, 16)), "type `Std::Map[Std::Symbol, Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"map pattern with a wider type and incompatible values": {
			input: `
				var a: Map[Symbol, Int] | nil = nil
				var { foo: "bar", baz: 2.5 } as b = a
				var c: nil = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(56, 3, 16), P(60, 3, 20)), "type `Std::Int` cannot ever match type `\"bar\"`"),
				error.NewFailure(L("<main>", P(68, 3, 28), P(70, 3, 30)), "type `Std::Int` cannot ever match type `2.5`"),
				error.NewFailure(L("<main>", P(100, 4, 18), P(100, 4, 18)), "type `Std::Map[Std::Symbol, Std::Int]` cannot be assigned to type `nil`"),
			},
		},
		"map pattern with a wider type and incompatible keys": {
			input: `
				var a: Map[Symbol, Int] | nil = nil
				var { "foo" => 29, 3.5 => 2 } as b = a
				var c: nil = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 3, 11), P(61, 3, 21)), "type `Std::Symbol` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(64, 3, 24), P(71, 3, 31)), "type `Std::Symbol` cannot ever match type `3.5`"),
				error.NewFailure(L("<main>", P(101, 4, 18), P(101, 4, 18)), "type `Std::Map[Std::Symbol, Std::Int]` cannot be assigned to type `nil`"),
			},
		},
		"map pattern with a wider type with HashMap declares a variable with value type": {
			input: `
				var a: HashMap[String, Int] | nil = nil
				var { "bar" => 2 as b } = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(92, 4, 16), P(92, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"map pattern with a wider incompatible value type with HashMap": {
			input: `
				var a: HashMap[Symbol, Int] | nil = nil
				var { bar: "foo", zed: 2.5, gamma: b } = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 3, 16), P(64, 3, 20)), "type `Std::Int` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(72, 3, 28), P(74, 3, 30)), "type `Std::Int` cannot ever match type `2.5`"),
			},
		},
		"map pattern with a wider incompatible key type with HashMap": {
			input: `
				var a: HashMap[Symbol, String | Int] | nil = nil
				var { bar: "foo", "zed" => 2, 1 => b } = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(76, 3, 23), P(85, 3, 32)), "type `Std::Symbol` cannot ever match type `\"zed\"`"),
				error.NewFailure(L("<main>", P(88, 3, 35), P(93, 3, 40)), "type `Std::Symbol` cannot ever match type `1`"),
			},
		},
		"map pattern with a wider type HashMap | Map declares a variable with value type": {
			input: `
				var a: HashMap[Symbol, Int] | Map[Symbol, Float] | nil = nil
				var { lol: b } as c = a
				var d: 9 = b
				var e: nil = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(109, 4, 16), P(109, 4, 16)), "type `Std::Int | Std::Float` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(128, 5, 18), P(128, 5, 18)), "type `Std::HashMap[Std::Symbol, Std::Int] | Std::Map[Std::Symbol, Std::Float]` cannot be assigned to type `nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestRecordPattern(t *testing.T) {
	tests := testTable{
		"declares variables in variable declaration": {
			input: `
				var %{ a } = %{ a: 3 }
				a = 5
			`,
		},
		"declares values in value declaration": {
			input: `
				val %{ a } = %{ a: 3 }
				a = 5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(32, 3, 5), P(32, 3, 5)), "local value `a` cannot be reassigned"),
			},
		},
		"record pattern with invalid value": {
			input: `
				var %{ a } = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(14, 2, 14)), "type `8` cannot ever match type `Std::Record[any, any]`"),
			},
		},
		"record pattern with a wider type declares a variable with value type": {
			input: `
				var a: Record[Symbol, Int] | nil = nil
				var %{ b } = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(78, 4, 16), P(78, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"record pattern with a wider type declares a variable with map type": {
			input: `
				var a: Record[Symbol, Int] | nil = nil
				var %{ g: 8 } as b = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(86, 4, 16), P(86, 4, 16)), "cannot use type `void` as a value in this context"),
			},
		},
		"record pattern with a wider type and incompatible values": {
			input: `
				var a: Record[Symbol, Int] | nil = nil
				var %{ foo: "bar", baz: 2.5 } as b = a
				var c: nil = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 3, 17), P(64, 3, 21)), "type `Std::Int` cannot ever match type `\"bar\"`"),
				error.NewFailure(L("<main>", P(72, 3, 29), P(74, 3, 31)), "type `Std::Int` cannot ever match type `2.5`"),
				error.NewFailure(L("<main>", P(104, 4, 18), P(104, 4, 18)), "cannot use type `void` as a value in this context"),
			},
		},
		"record pattern with a wider type and incompatible keys": {
			input: `
				var a: Record[Symbol, Int] | nil = nil
				var %{ "foo" => 29, 3.5 => 2 } as b = a
				var c: nil = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(55, 3, 12), P(65, 3, 22)), "type `Std::Symbol` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(68, 3, 25), P(75, 3, 32)), "type `Std::Symbol` cannot ever match type `3.5`"),
				error.NewFailure(L("<main>", P(105, 4, 18), P(105, 4, 18)), "cannot use type `void` as a value in this context"),
			},
		},
		"record pattern with a wider type with HashRecord declares a variable with value type": {
			input: `
				var a: HashRecord[String, Int] | nil = nil
				var %{ "bar" => 2 as b } = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(96, 4, 16), P(96, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"record pattern with a wider type with HashMap declares a variable with value type": {
			input: `
				var a: HashMap[String, Int] | nil = nil
				var %{ "bar" => 2 as b } = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(93, 4, 16), P(93, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"record pattern with a wider incompatible value type with HashRecord": {
			input: `
				var a: HashRecord[Symbol, Int] | nil = nil
				var %{ bar: "foo", zed: 2.5, gamma: b } = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(64, 3, 17), P(68, 3, 21)), "type `Std::Int` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(76, 3, 29), P(78, 3, 31)), "type `Std::Int` cannot ever match type `2.5`"),
			},
		},
		"record pattern with a wider incompatible key type with HashRecord": {
			input: `
				var a: HashRecord[Symbol, String | Int] | nil = nil
				var %{ bar: "foo", "zed" => 2, 1 => b } = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(80, 3, 24), P(89, 3, 33)), "type `Std::Symbol` cannot ever match type `\"zed\"`"),
				error.NewFailure(L("<main>", P(92, 3, 36), P(97, 3, 41)), "type `Std::Symbol` cannot ever match type `1`"),
			},
		},
		"record pattern with a wider type HashMap | Record declares a variable with value type": {
			input: `
				var a: HashMap[Symbol, Int] | Record[Symbol, Float] | nil = nil
				var %{ lol: b } as c = a
				var d: 9 = b
				var e: nil = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(113, 4, 16), P(113, 4, 16)), "type `Std::Int | Std::Float` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(132, 5, 18), P(132, 5, 18)), "cannot use type `void` as a value in this context"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestObjectPattern(t *testing.T) {
	tests := testTable{
		"declares variables in variable declaration": {
			input: `
				var String(length) = "foo"
				length = 5
			`,
		},
		"declares values in value declaration": {
			input: `
				val String(length) = "foo"
				length = 5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 5), P(41, 3, 10)), "local value `length` cannot be reassigned"),
			},
		},
		"identifier - invalid value": {
			input: `
				var String(length) = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(22, 2, 22)), "type `3` cannot ever match type `Std::String`"),
			},
		},
		"identifier - nonexistent method": {
			input: `
				var String(lol) = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 2, 16), P(18, 2, 18)), "method `lol` is not defined on type `Std::String`"),
			},
		},
		"identifier - valid getter": {
			input: `
				var String(length) as s = "foo"
				var a: 9 = length
				var b: 7 = s
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(52, 3, 16), P(57, 3, 21)), "type `Std::Int` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(74, 4, 16), P(74, 4, 16)), "type `Std::String` cannot be assigned to type `7`"),
			},
		},
		"identifier - valid getter and wider type": {
			input: `
				var a: Int | String | nil = nil
				var String(length) as s = a
				var b: 9 = length
				var c: 7 = s
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(84, 4, 16), P(89, 4, 21)), "type `Std::Int` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(106, 5, 16), P(106, 5, 16)), "type `Std::String` cannot be assigned to type `7`"),
			},
		},
		"identifier - method with required arguments": {
			input: `
				class Foo
					def bar(a: Int): Int then a
				end

				var Foo(bar) as f = Foo()
				var a: 9 = bar
				var b: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(69, 6, 13), P(71, 6, 15)), "argument `a` is missing in call to `bar`"),
				error.NewFailure(L("<main>", P(102, 7, 16), P(104, 7, 18)), "type `Std::Int` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(121, 8, 16), P(121, 8, 16)), "type `Foo` cannot be assigned to type `7`"),
			},
		},
		"identifier - void method": {
			input: `
				class Foo
					def bar; end
				end

				var Foo(bar) as f = Foo()
				var a: 9 = bar
				var b: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(54, 6, 13), P(56, 6, 15)), "cannot use type `void` as a value in this context"),
				error.NewFailure(L("<main>", P(106, 8, 16), P(106, 8, 16)), "type `Foo` cannot be assigned to type `7`"),
			},
		},
		"identifier - generic getter with bounds": {
			input: `
				class Foo
					def bar[T < CoercibleNumeric]: T then loop; end
				end

				var Foo(bar) as f = Foo()
				var a: 9 = bar
				var b: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(122, 7, 16), P(124, 7, 18)), "type `Std::CoercibleNumeric` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(141, 8, 16), P(141, 8, 16)), "type `Foo` cannot be assigned to type `7`"),
			},
		},
		"identifier - generic getter without bounds": {
			input: `
				class Foo
					def bar[T]: T then loop; end
				end

				var Foo(bar) as f = Foo()
				var a: 9 = bar
				var b: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(103, 7, 16), P(105, 7, 18)), "type `any` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(122, 8, 16), P(122, 8, 16)), "type `Foo` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class and simple type": {
			input: `
				class Foo[T]
					def bar: T then loop; end
				end

				var Foo(bar) as f = Foo::[String]()
				var a: 9 = bar
				var b: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(113, 7, 16), P(115, 7, 18)), "type `Std::String` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(132, 8, 16), P(132, 8, 16)), "type `Foo[Std::String]` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class and wide type": {
			input: `
				class Foo[T]
					def bar: T then loop; end
				end

				class Bar < Foo[Int]
				end

				var a: Float | Foo[String] | Bar | nil = nil
				var Foo(bar) as f = a
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(182, 11, 16), P(184, 11, 18)), "type `Std::String | Std::Int` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(201, 12, 16), P(201, 12, 16)), "type `Foo[Std::String] | Bar` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class and wide type without the class": {
			input: `
				class Foo[T < Value]
					def bar: T then loop; end
				end

				var a: any = nil
				var Foo(bar) as f = a
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(128, 8, 16), P(130, 8, 18)), "type `Std::Value` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(147, 9, 16), P(147, 9, 16)), "type `Foo[Std::Value]` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class and wide interface type that intersects with it": {
			input: `
				class Foo[T < Value]
					def bar: T then loop; end
				end

				interface Bar
					sig bar: String
				end

				var a: Bar | nil = Foo::[String]()
				var Foo(bar) as f = a
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(194, 12, 16), P(196, 12, 18)), "type `Std::String` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(213, 13, 16), P(213, 13, 16)), "type `Foo[Std::String]` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class and interface type that intersects with it": {
			input: `
				class Foo[T < Value]
					def bar: T then loop; end
				end

				interface Bar
					sig bar: String
				end

				var a: Bar = Foo::[String]()
				var Foo(bar) as f = a
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(188, 12, 16), P(190, 12, 18)), "type `Std::String` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(207, 13, 16), P(207, 13, 16)), "type `Foo[Std::String]` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class and interface type with union types": {
			input: `
				class Foo[T < Value]
					def bar: T | nil then loop; end
					def baz: T | Int then loop; end
				end

				interface Bar
					sig bar: String | Float | nil
					sig baz: String | Float | Int
				end

				def a: Bar then loop; end
				var Foo(bar) as f = a()
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(279, 14, 16), P(281, 14, 18)), "type `nil | Std::String | Std::Float` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(298, 15, 16), P(298, 15, 16)), "type `Foo[Std::String | Std::Float]` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class and interface type with exclusive return types": {
			input: `
				class Foo[T < Value]
					def bar: T then loop; end
					def baz: T then loop; end
				end

				interface Bar
					sig bar: String
					sig baz: Int
				end

				def a: Bar then loop; end
				var Foo(bar) as f = a()
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(201, 13, 9), P(208, 13, 16)), "type `Bar` cannot ever match type `Foo`"),
				error.NewFailure(L("<main>", P(236, 14, 16), P(238, 14, 18)), "type `Std::Value` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(255, 15, 16), P(255, 15, 16)), "type `Foo[Std::Value]` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class with invalid method and interface type": {
			input: `
				class Foo[T < Value]
					def bar: Int then 3
					def baz: T then loop; end
				end

				interface Bar
					sig bar: String
				end

				class Baz
					def bar: String then "bar"
				end

				var a: Bar = Baz()
				var Foo(bar) as f = a
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(225, 16, 9), P(232, 16, 16)), "type `Bar` cannot ever match type `Foo`"),
				error.NewFailure(L("<main>", P(258, 17, 16), P(260, 17, 18)), "type `Std::Int` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(277, 18, 16), P(277, 18, 16)), "type `Foo[Std::Value]` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class with missing methods and interface type": {
			input: `
				class Foo[T < Value]
					def baz: T then loop; end
				end

				interface Bar
					sig bar: String
				end

				class Baz
					def bar: String then "bar"
				end

				var a: Bar = Baz()
				var Foo(bar) as f = a
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(200, 15, 9), P(207, 15, 16)), "type `Bar` cannot ever match type `Foo`"),
				error.NewFailure(L("<main>", P(204, 15, 13), P(206, 15, 15)), "method `bar` is not defined on type `Foo[Std::Value]`"),
				error.NewFailure(L("<main>", P(252, 17, 16), P(252, 17, 16)), "type `Foo[Std::Value]` cannot be assigned to type `7`"),
			},
		},

		"identifier - getter on a generic class and generic interface type that intersects with it": {
			input: `
				class Foo[T < Value]
					def bar: T then loop; end
				end

				interface Bar[T]
					sig bar: T
				end

				var a: Bar[String] = Foo::[String]()
				var Foo(bar) as f = a
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(194, 12, 16), P(196, 12, 18)), "type `Std::String` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(213, 13, 16), P(213, 13, 16)), "type `Foo[Std::String]` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class with invalid method and generic interface type": {
			input: `
				class Foo[T < Value]
					def bar: Int then 3
					def baz: T then loop; end
				end

				interface Bar[T]
					sig bar: T
				end

				class Baz
					def bar: String then "bar"
				end

				var a: Bar[String] = Baz()
				var Foo(bar) as f = a
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(231, 16, 9), P(238, 16, 16)), "type `Bar[Std::String]` cannot ever match type `Foo`"),
				error.NewFailure(L("<main>", P(264, 17, 16), P(266, 17, 18)), "type `Std::Int` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(283, 18, 16), P(283, 18, 16)), "type `Foo[Std::Value]` cannot be assigned to type `7`"),
			},
		},
		"identifier - getter on a generic class with missing methods and generic interface type": {
			input: `
				class Foo[T < Value]
					def baz: T then loop; end
				end

				interface Bar[T]
					sig bar: T
				end

				class Baz
					def bar: String then "bar"
				end

				var a: Bar[String] = Baz()
				var Foo(bar) as f = a
				var b: 9 = bar
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(206, 15, 9), P(213, 15, 16)), "type `Bar[Std::String]` cannot ever match type `Foo`"),
				error.NewFailure(L("<main>", P(210, 15, 13), P(212, 15, 15)), "method `bar` is not defined on type `Foo[Std::Value]`"),
				error.NewFailure(L("<main>", P(258, 17, 16), P(258, 17, 16)), "type `Foo[Std::Value]` cannot be assigned to type `7`"),
			},
		},

		"key value - invalid value": {
			input: `
				var String(length: l) = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(25, 2, 25)), "type `3` cannot ever match type `Std::String`"),
			},
		},
		"key value - nonexistent method": {
			input: `
				var String(lol: l) = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 2, 16), P(21, 2, 21)), "method `lol` is not defined on type `Std::String`"),
			},
		},
		"key value - valid getter": {
			input: `
				var String(length: l) as s = "foo"
				var a: 9 = l
				var b: 7 = s
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(55, 3, 16), P(55, 3, 16)), "type `Std::Int` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(72, 4, 16), P(72, 4, 16)), "type `Std::String` cannot be assigned to type `7`"),
			},
		},
		"key value - invalid value pattern": {
			input: `
				var String(length: "lol") as s = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 2, 24), P(28, 2, 28)), "type `Std::Int` cannot ever match type `\"lol\"`"),
			},
		},
		"key value - valid getter and wider type": {
			input: `
				var a: Int | String | nil = nil
				var String(length: l) as s = a
				var b: 9 = l
				var c: 7 = s
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(87, 4, 16), P(87, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(104, 5, 16), P(104, 5, 16)), "type `Std::String` cannot be assigned to type `7`"),
			},
		},
		"key value - method with required arguments": {
			input: `
				class Foo
					def bar(a: Int): Int then a
				end

				var Foo(bar: b) as f = Foo()
				var a: 9 = b
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(69, 6, 13), P(74, 6, 18)), "argument `a` is missing in call to `bar`"),
				error.NewFailure(L("<main>", P(105, 7, 16), P(105, 7, 16)), "type `Std::Int` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(122, 8, 16), P(122, 8, 16)), "type `Foo` cannot be assigned to type `7`"),
			},
		},
		"key value - void method": {
			input: `
				class Foo
					def bar; end
				end

				var Foo(bar: b) as f = Foo()
				var a: 9 = b
				var c: 7 = f
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(54, 6, 13), P(59, 6, 18)), "cannot use type `void` as a value in this context"),
				error.NewFailure(L("<main>", P(107, 8, 16), P(107, 8, 16)), "type `Foo` cannot be assigned to type `7`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestConstantPattern(t *testing.T) {
	tests := testTable{
		"public constant - type without value": {
			input: `
				typedef Foo = 2
				var Foo as a = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(29, 3, 9), P(31, 3, 11)), "`Foo` cannot be used as a value in expressions"),
			},
		},
		"public constant - nonexistent constant": {
			input: `
				var Foo as a = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "undefined constant `Foo`"),
			},
		},
		"public constant - invalid literal value": {
			input: `
				const Foo = 3
				var Foo as a = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 9), P(29, 3, 11)), "type `8` cannot ever match type `3`"),
			},
		},
		"public constant - valid literal value": {
			input: `
				const Foo = 3
				var Foo as a = 3
				var b: nil = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 18), P(57, 4, 18)), "type `Std::Int` cannot be assigned to type `nil`"),
			},
		},

		"private constant - type without value": {
			input: `
				typedef _Foo = 2
				var _Foo as a = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 3, 9), P(33, 3, 12)), "`_Foo` cannot be used as a value in expressions"),
			},
		},
		"private constant - nonexistent constant": {
			input: `
				var _Foo as a = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "undefined constant `_Foo`"),
			},
		},
		"private constant - invalid literal value": {
			input: `
				const _Foo = 3
				var _Foo as a = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 9), P(31, 3, 12)), "type `8` cannot ever match type `3`"),
			},
		},
		"private constant - valid literal value": {
			input: `
				const _Foo = 3
				var _Foo as a = 3
				var b: nil = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(59, 4, 18), P(59, 4, 18)), "type `Std::Int` cannot be assigned to type `nil`"),
			},
		},

		"constant lookup - type without value": {
			input: `
				module Foo
					typedef Bar = 2
				end
				var Foo::Bar as a = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(58, 5, 14), P(60, 5, 16)), "`Foo::Bar` cannot be used as a value in expressions"),
			},
		},
		"constant lookup - nonexistent module": {
			input: `
				var Foo::Bar as a = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "undefined constant `Foo`"),
			},
		},
		"constant lookup - nonexistent constant": {
			input: `
				module Foo; end
				var Foo::Bar as a = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 14), P(36, 3, 16)), "undefined constant `Foo::Bar`"),
			},
		},
		"constant lookup - invalid literal value": {
			input: `
				module Foo
					const Bar = 3
				end
				var Foo::Bar as a = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 5, 9), P(58, 5, 16)), "type `8` cannot ever match type `3`"),
			},
		},
		"constant lookup - valid literal value": {
			input: `
				module Foo
					const Bar = 3
				end
				var Foo::Bar as a = 3
				var b: nil = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(86, 6, 18), P(86, 6, 18)), "type `Std::Int` cannot be assigned to type `nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestUnaryPattern(t *testing.T) {
	tests := testTable{
		"== - invalid type": {
			input: `
				var == 10 as a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 2, 12), P(13, 2, 13)), "type `3` cannot ever match type `10`"),
			},
		},
		"== - valid type": {
			input: `
				var a: String | Int | nil = nil
				var b: String? = nil
				var == b as c = a
				var d: 1 = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(99, 5, 16), P(99, 5, 16)), "type `Std::String?` cannot be assigned to type `1`"),
			},
		},
		"!= - invalid type": {
			input: `
				var != 10 as a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 2, 12), P(13, 2, 13)), "type `3` cannot ever match type `10`"),
			},
		},
		"!= - valid type": {
			input: `
				var a: String | Int | nil = nil
				var b: String? = nil
				var != b as c = a
				var d: 1 = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(99, 5, 16), P(99, 5, 16)), "type `Std::String | Std::Int | nil` cannot be assigned to type `1`"),
			},
		},

		"=== - invalid type": {
			input: `
				var === 10 as a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(13, 2, 13), P(14, 2, 14)), "type `3` cannot ever match type `10`"),
			},
		},
		"=== - valid type": {
			input: `
				var a: String | Int | nil = nil
				var b: String? = nil
				var === b as c = a
				var d: 1 = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(100, 5, 16), P(100, 5, 16)), "type `Std::String?` cannot be assigned to type `1`"),
			},
		},
		"!== - invalid type": {
			input: `
				var !== 10 as a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(13, 2, 13), P(14, 2, 14)), "type `3` cannot ever match type `10`"),
			},
		},
		"!== - valid type": {
			input: `
				var a: String | Int | nil = nil
				var b: String? = nil
				var !== b as c = a
				var d: 1 = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(100, 5, 16), P(100, 5, 16)), "type `Std::String | Std::Int | nil` cannot be assigned to type `1`"),
			},
		},

		"=~ - different type": {
			input: `
				var =~ 10 as a = 3
			`,
		},
		"=~ - intersecting type": {
			input: `
				var a: String | Int | nil = nil
				var b: String? = nil
				var =~ b as c = a
				var d: 1 = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(99, 5, 16), P(99, 5, 16)), "type `Std::String | Std::Int | nil` cannot be assigned to type `1`"),
			},
		},
		"!~ - different type": {
			input: `
				var !~ 10 as a = 3
			`,
		},
		"!~ - intersecting type": {
			input: `
				var a: String | Int | nil = nil
				var b: String? = nil
				var !~ b as c = a
				var d: 1 = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(99, 5, 16), P(99, 5, 16)), "type `Std::String | Std::Int | nil` cannot be assigned to type `1`"),
			},
		},

		"< - invalid type": {
			input: `
				var < .1 as a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 2, 11), P(12, 2, 12)), "type `3` cannot ever match type `Std::Float`"),
			},
		},
		"< - valid type": {
			input: `
				var a: Float | Int | nil = nil
				var b: BigFloat | Float = 2.2
				var < b as c = a
				var d: 1 = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(106, 5, 16), P(106, 5, 16)), "type `Std::Float` cannot be assigned to type `1`"),
			},
		},
		"< - type without method": {
			input: `
				var a: String? = nil
				var < nil as c = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 9), P(34, 3, 9)), "method `<` is not defined on type `Std::Nil`"),
			},
		},
		"<= - invalid type": {
			input: `
				var <= .1 as a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 2, 12), P(13, 2, 13)), "type `3` cannot ever match type `Std::Float`"),
			},
		},
		"<= - valid type": {
			input: `
				var a: Float | Int | nil = nil
				var b: BigFloat | Float = 2.2
				var <= b as c = a
				var d: 1 = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(107, 5, 16), P(107, 5, 16)), "type `Std::Float` cannot be assigned to type `1`"),
			},
		},
		"<= - type without method": {
			input: `
				var a: String? = nil
				var <= nil as c = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 9), P(35, 3, 10)), "method `<=` is not defined on type `Std::Nil`"),
			},
		},
		"> - invalid type": {
			input: `
				var > 1.2 as a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 2, 11), P(13, 2, 13)), "type `3` cannot ever match type `Std::Float`"),
			},
		},
		"> - valid type": {
			input: `
				var a: Float | Int | nil = nil
				var b: BigFloat | Float = 2.2
				var > b as c = a
				var d: 1 = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(106, 5, 16), P(106, 5, 16)), "type `Std::Float` cannot be assigned to type `1`"),
			},
		},
		"> - type without method": {
			input: `
				var a: String? = nil
				var > nil as c = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 9), P(34, 3, 9)), "method `>` is not defined on type `Std::Nil`"),
			},
		},
		">= - invalid type": {
			input: `
				var >= .1 as a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 2, 12), P(13, 2, 13)), "type `3` cannot ever match type `Std::Float`"),
			},
		},
		">= - valid type": {
			input: `
				var a: Float | Int | nil = nil
				var b: BigFloat | Float = 2.2
				var >= b as c = a
				var d: 1 = c
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(107, 5, 16), P(107, 5, 16)), "type `Std::Float` cannot be assigned to type `1`"),
			},
		},
		">= - type without method": {
			input: `
				var a: String? = nil
				var >= nil as c = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 9), P(35, 3, 10)), "method `>=` is not defined on type `Std::Nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestBinaryPattern(t *testing.T) {
	tests := testTable{
		"|| - matching patterns": {
			input: `
				var a: String | Float | Int = 1
				var 1 || "foo" as b = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(80, 4, 16), P(80, 4, 16)), "type `Std::Int | Std::String` cannot be assigned to type `9`"),
			},
		},
		"|| - invalid patterns": {
			input: `
				var a: Char | Float | nil = nil
				var 1 || "foo" as b = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 9), P(45, 3, 9)), "type `Std::Char | Std::Float | nil` cannot ever match type `1`"),
				error.NewFailure(L("<main>", P(50, 3, 14), P(54, 3, 18)), "type `Std::Char | Std::Float | nil` cannot ever match type `\"foo\"`"),
				error.NewFailure(L("<main>", P(80, 4, 16), P(80, 4, 16)), "type `Std::Int | Std::String` cannot be assigned to type `9`"),
			},
		},
		"|| - conditional variables": {
			input: `
				var a: String | Float | Int = 1
				var 1 || ("foo" as b) = a
				var c: 9 = b
				b = "lol"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(82, 4, 16), P(82, 4, 16)), "type `Std::String?` cannot be assigned to type `9`"),
			},
		},
		"|| - conditional values": {
			input: `
				val a: String | Float | Int = 1
				val 1 || ("foo" as b) = a
				val c: 9 = b
				b = "lol"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(82, 4, 16), P(82, 4, 16)), "type `Std::String?` cannot be assigned to type `9`"),
				error.NewFailure(L("<main>", P(88, 5, 5), P(88, 5, 5)), "local value `b` cannot be reassigned"),
			},
		},

		"&& - matching patterns": {
			input: `
				var a: String | Float | Int = 1
				var < 5 && > 10 as b = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(81, 4, 16), P(81, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"&& - non matching type": {
			input: `
				var < 5 && > 10 as b = 2.2
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 2, 11), P(11, 2, 11)), "type `2.2` cannot ever match type `Std::Int`"),
				error.NewFailure(L("<main>", P(18, 2, 18), P(19, 2, 19)), "type `2.2` cannot ever match type `Std::Int`"),
			},
		},
		"&& - incompatible patterns": {
			input: `
				var a: String | Float | Int = 1
				var < 5 && > 2.5 as b = a
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(45, 3, 9), P(56, 3, 20)), "this pattern is impossible to satisfy"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
