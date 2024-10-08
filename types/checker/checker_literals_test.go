package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestArrayTupleLiteral(t *testing.T) {
	tests := testTable{
		"modifier if else": {
			input: `
				var a: bool = false
				var foo = %[1, 2.5 if a else "bar"]
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(80, 4, 16), P(82, 4, 18)), "type `Std::ArrayTuple[1 | 2.5 | \"bar\"]` cannot be assigned to type `9`"),
			},
		},
		"modifier if": {
			input: `
				var a: bool = false
				var foo = %[1, 2.5 if a]
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(69, 4, 16), P(71, 4, 18)), "type `Std::ArrayTuple[1 | 2.5]` cannot be assigned to type `9`"),
			},
		},
		"modifier if with narrowing": {
			input: `
				var a: String? = nil
				var foo = %[1, a.lol if a]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 20), P(49, 3, 24)), "method `lol` is not defined on type `Std::String`"),
			},
		},
		"truthy modifier if": {
			input: `
				var foo = %[1, 2.5 if  5]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(28, 2, 28), P(28, 2, 28)), "this condition will always have the same result since type `5` is truthy"),
				error.NewFailure(L("<main>", P(46, 3, 16), P(48, 3, 18)), "type `Std::ArrayTuple[1 | 2.5]` cannot be assigned to type `9`"),
			},
		},
		"falsy modifier if": {
			input: `
				var foo = %[1, 2.5 if nil]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 2, 27), P(29, 2, 29)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L("<main>", P(20, 2, 20), P(22, 2, 22)), "unreachable code"),
				error.NewFailure(L("<main>", P(47, 3, 16), P(49, 3, 18)), "type `Std::ArrayTuple[1]` cannot be assigned to type `9`"),
			},
		},

		"modifier unless": {
			input: `
				var a: bool = false
				var foo = %[1, 2.5 unless a]
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(73, 4, 16), P(75, 4, 18)), "type `Std::ArrayTuple[1 | 2.5]` cannot be assigned to type `9`"),
			},
		},
		"modifier unless with narrowing": {
			input: `
				var a: String? = nil
				var foo = %[1, a.lol unless a]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 20), P(49, 3, 24)), "method `lol` is not defined on type `Std::Nil`"),
			},
		},
		"truthy modifier unless": {
			input: `
				var foo = %[1, 2.5 unless 5]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(31, 2, 31), P(31, 2, 31)), "this condition will always have the same result since type `5` is truthy"),
				error.NewWarning(L("<main>", P(20, 2, 20), P(22, 2, 22)), "unreachable code"),
				error.NewFailure(L("<main>", P(49, 3, 16), P(51, 3, 18)), "type `Std::ArrayTuple[1]` cannot be assigned to type `9`"),
			},
		},
		"falsy modifier unless": {
			input: `
				var foo = %[1, 2.5 unless nil]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(31, 2, 31), P(33, 2, 33)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewFailure(L("<main>", P(51, 3, 16), P(53, 3, 18)), "type `Std::ArrayTuple[1 | 2.5]` cannot be assigned to type `9`"),
			},
		},

		"infer array tuple": {
			input: `
				var foo = %[1, 2]
				var a: ArrayTuple[Int] = foo
			`,
		},
		"infer array tuple with different argument types": {
			input: `
				var a = %[1, 2.2, "foo"]
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 16), P(45, 3, 16)), "type `Std::ArrayTuple[1 | 2.2 | \"foo\"]` cannot be assigned to type `9`"),
			},
		},
		"infer empty array tuple": {
			input: `
				var foo = %[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 16), P(36, 3, 18)), "type `Std::ArrayTuple[any]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define array tuple type": {
			input: `
				var foo: ArrayTuple[Int] = %[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 3, 16), P(53, 3, 18)), "type `Std::ArrayTuple[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define tuple type": {
			input: `
				var foo: Tuple[Int] = %[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 3, 16), P(48, 3, 18)), "type `Std::Tuple[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define tuple incompatible type": {
			input: `
				var foo: Tuple[Int] = %["", 2.2]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 2, 27), P(36, 2, 36)), "type `\"\" | 2.2` cannot be assigned to type `Std::Int`"),
			},
		},

		"infer word array tuple": {
			input: `
				var foo = %w[foo bar]
				var a: ArrayTuple[String] = foo
			`,
		},
		"infer empty word array tuple": {
			input: `
				var foo = %w[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayTuple[Std::String]` cannot be assigned to type `9`"),
			},
		},

		"infer symbol array tuple": {
			input: `
				var foo = %s[foo bar]
				var a: ArrayTuple[Symbol] = foo
			`,
		},
		"infer empty symbol array tuple": {
			input: `
				var foo = %s[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayTuple[Std::Symbol]` cannot be assigned to type `9`"),
			},
		},

		"infer hex array tuple": {
			input: `
				var foo = %x[1f4 fff]
				var a: ArrayTuple[Int] = foo
			`,
		},
		"infer empty hex array tuple": {
			input: `
				var foo = %x[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayTuple[Std::Int]` cannot be assigned to type `9`"),
			},
		},

		"infer bin array tuple": {
			input: `
				var foo = %b[111 100]
				var a: ArrayTuple[Int] = foo
			`,
		},
		"infer empty bin array tuple": {
			input: `
				var foo = %b[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayTuple[Std::Int]` cannot be assigned to type `9`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestHashSetLiteral(t *testing.T) {
	tests := testTable{
		"modifier if else": {
			input: `
				var a: bool = false
				var foo = ^[1, 2.5 if a else "bar"]
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(80, 4, 16), P(82, 4, 18)), "type `Std::HashSet[Std::Int | Std::Float | Std::String]` cannot be assigned to type `9`"),
			},
		},
		"modifier if": {
			input: `
				var a: bool = false
				var foo = ^[1, 2.5 if a]
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(69, 4, 16), P(71, 4, 18)), "type `Std::HashSet[Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"modifier if with narrowing": {
			input: `
				var a: String? = nil
				var foo = ^[1, a.lol if a]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 20), P(49, 3, 24)), "method `lol` is not defined on type `Std::String`"),
			},
		},
		"truthy modifier if": {
			input: `
				var foo = ^[1, 2.5 if  5]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(28, 2, 28), P(28, 2, 28)), "this condition will always have the same result since type `5` is truthy"),
				error.NewFailure(L("<main>", P(46, 3, 16), P(48, 3, 18)), "type `Std::HashSet[Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"falsy modifier if": {
			input: `
				var foo = ^[1, 2.5 if nil]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 2, 27), P(29, 2, 29)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L("<main>", P(20, 2, 20), P(22, 2, 22)), "unreachable code"),
				error.NewFailure(L("<main>", P(47, 3, 16), P(49, 3, 18)), "type `Std::HashSet[Std::Int]` cannot be assigned to type `9`"),
			},
		},

		"modifier unless": {
			input: `
				var a: bool = false
				var foo = ^[1, 2.5 unless a]
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(73, 4, 16), P(75, 4, 18)), "type `Std::HashSet[Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"modifier unless with narrowing": {
			input: `
				var a: String? = nil
				var foo = ^[1, a.lol unless a]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 20), P(49, 3, 24)), "method `lol` is not defined on type `Std::Nil`"),
			},
		},
		"truthy modifier unless": {
			input: `
				var foo = ^[1, 2.5 unless 5]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(31, 2, 31), P(31, 2, 31)), "this condition will always have the same result since type `5` is truthy"),
				error.NewWarning(L("<main>", P(20, 2, 20), P(22, 2, 22)), "unreachable code"),
				error.NewFailure(L("<main>", P(49, 3, 16), P(51, 3, 18)), "type `Std::HashSet[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"falsy modifier unless": {
			input: `
				var foo = ^[1, 2.5 unless nil]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(31, 2, 31), P(33, 2, 33)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewFailure(L("<main>", P(51, 3, 16), P(53, 3, 18)), "type `Std::HashSet[Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},

		"infer array list": {
			input: `
				var foo = ^[1, 2]
				var a: HashSet[Int] = foo
			`,
		},
		"infer array list with different argument types": {
			input: `
				var a = ^[1, 2.2, "foo"]
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 16), P(45, 3, 16)), "type `Std::HashSet[Std::Int | Std::Float | Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer empty array list": {
			input: `
				var foo = ^[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 16), P(36, 3, 18)), "type `Std::HashSet[any]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define hash set type": {
			input: `
				var foo: HashSet[Int] = ^[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 3, 16), P(50, 3, 18)), "type `Std::HashSet[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define Set type": {
			input: `
				var foo: Set[Int] = ^[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 16), P(46, 3, 18)), "type `Std::Set[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define Set incompatible type": {
			input: `
				var foo: Set[Int] = ^["", 2.2]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(25, 2, 25), P(34, 2, 34)), "type `Std::String | Std::Float` cannot be assigned to type `Std::Int`"),
			},
		},
		"int capacity": {
			input: `
				var foo: HashSet[Float] = ^[1.2]:9
			`,
		},
		"uint8 capacity": {
			input: `
				var foo: HashSet[Float] = ^[1.2]:9u8
			`,
		},
		"invalid capacity": {
			input: `
				var foo: HashSet[Float] = ^[1.2]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 2, 31), P(40, 2, 40)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer word hash set": {
			input: `
				var foo = ^w[foo bar]
				var a: HashSet[String] = foo
			`,
		},
		"infer empty word hash set": {
			input: `
				var foo = ^w[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::HashSet[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"word hash set int capacity": {
			input: `
				var foo: HashSet[String] = ^w[1.2]:9
			`,
		},
		"word hash set uint8 capacity": {
			input: `
				var foo: HashSet[String] = ^w[1.2]:9u8
			`,
		},
		"word hash set invalid capacity": {
			input: `
				var foo: HashSet[String] = ^w[1.2]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(40, 2, 40), P(42, 2, 42)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer symbol hash set": {
			input: `
				var foo = ^s[foo bar]
				var a: HashSet[Symbol] = foo
			`,
		},
		"infer empty symbol hash set": {
			input: `
				var foo = ^s[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::HashSet[Std::Symbol]` cannot be assigned to type `9`"),
			},
		},
		"symbol hash set int capacity": {
			input: `
				var foo: HashSet[Symbol] = ^s[1.2]:9
			`,
		},
		"symbol hash set uint8 capacity": {
			input: `
				var foo: HashSet[Symbol] = ^s[1.2]:9u8
			`,
		},
		"symbol hash set invalid capacity": {
			input: `
				var foo: HashSet[Symbol] = ^s[1.2]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(40, 2, 40), P(42, 2, 42)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer hex hash set": {
			input: `
				var foo = ^x[1f4 fff]
				var a: HashSet[Int] = foo
			`,
		},
		"infer empty hex hash set": {
			input: `
				var foo = ^x[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::HashSet[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"hex hash set int capacity": {
			input: `
				var foo: HashSet[Int] = ^x[1f9]:9
			`,
		},
		"hex hash set uint8 capacity": {
			input: `
				var foo: HashSet[Int] = ^x[fff]:9u8
			`,
		},
		"hex hash set invalid capacity": {
			input: `
				var foo: HashSet[Int] = ^x[1ef]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(37, 2, 37), P(39, 2, 39)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer bin hash set": {
			input: `
				var foo = ^b[111 100]
				var a: HashSet[Int] = foo
			`,
		},
		"infer empty bin hash set": {
			input: `
				var foo = ^b[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::HashSet[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"bin hash set int capacity": {
			input: `
				var foo: HashSet[Int] = ^b[100]:9
			`,
		},
		"bin hash set uint8 capacity": {
			input: `
				var foo: HashSet[Int] = ^x[111]:9u8
			`,
		},
		"bin hash set invalid capacity": {
			input: `
				var foo: HashSet[Int] = ^b[101]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(37, 2, 37), P(39, 2, 39)), "capacity must be an integer, got `9.2`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestArrayListLiteral(t *testing.T) {
	tests := testTable{
		"modifier if else": {
			input: `
				var a: bool = false
				var foo = [1, 2.5 if a else "bar"]
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(79, 4, 16), P(81, 4, 18)), "type `Std::ArrayList[Std::Int | Std::Float | Std::String]` cannot be assigned to type `9`"),
			},
		},
		"modifier if": {
			input: `
				var a: bool = false
				var foo = [1, 2.5 if a]
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(68, 4, 16), P(70, 4, 18)), "type `Std::ArrayList[Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"modifier if with narrowing": {
			input: `
				var a: String? = nil
				var foo = [1, a.lol if a]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 19), P(48, 3, 23)), "method `lol` is not defined on type `Std::String`"),
			},
		},
		"truthy modifier if": {
			input: `
				var foo = [1, 2.5 if  5]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 2, 27), P(27, 2, 27)), "this condition will always have the same result since type `5` is truthy"),
				error.NewFailure(L("<main>", P(45, 3, 16), P(47, 3, 18)), "type `Std::ArrayList[Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"falsy modifier if": {
			input: `
				var foo = [1, 2.5 if nil]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(26, 2, 26), P(28, 2, 28)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L("<main>", P(19, 2, 19), P(21, 2, 21)), "unreachable code"),
				error.NewFailure(L("<main>", P(46, 3, 16), P(48, 3, 18)), "type `Std::ArrayList[Std::Int]` cannot be assigned to type `9`"),
			},
		},

		"modifier unless": {
			input: `
				var a: bool = false
				var foo = [1, 2.5 unless a]
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(72, 4, 16), P(74, 4, 18)), "type `Std::ArrayList[Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"modifier unless with narrowing": {
			input: `
				var a: String? = nil
				var foo = [1, a.lol unless a]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 19), P(48, 3, 23)), "method `lol` is not defined on type `Std::Nil`"),
			},
		},
		"truthy modifier unless": {
			input: `
				var foo = [1, 2.5 unless 5]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(30, 2, 30), P(30, 2, 30)), "this condition will always have the same result since type `5` is truthy"),
				error.NewWarning(L("<main>", P(19, 2, 19), P(21, 2, 21)), "unreachable code"),
				error.NewFailure(L("<main>", P(48, 3, 16), P(50, 3, 18)), "type `Std::ArrayList[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"falsy modifier unless": {
			input: `
				var foo = [1, 2.5 unless nil]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(30, 2, 30), P(32, 2, 32)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewFailure(L("<main>", P(50, 3, 16), P(52, 3, 18)), "type `Std::ArrayList[Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},

		"infer array list": {
			input: `
				var foo = [1, 2]
				var a: ArrayList[Int] = foo
			`,
		},
		"infer array list with different argument types": {
			input: `
				var a = [1, 2.2, "foo"]
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 16), P(44, 3, 16)), "type `Std::ArrayList[Std::Int | Std::Float | Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer empty array list": {
			input: `
				var foo = []
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 3, 16), P(35, 3, 18)), "type `Std::ArrayList[any]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define array list type": {
			input: `
				var foo: ArrayList[Int] = []
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(49, 3, 16), P(51, 3, 18)), "type `Std::ArrayList[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define List type": {
			input: `
				var foo: List[Int] = []
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 16), P(46, 3, 18)), "type `Std::List[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define List incompatible type": {
			input: `
				var foo: List[Int] = ["", 2.2]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(26, 2, 26), P(34, 2, 34)), "type `Std::String | Std::Float` cannot be assigned to type `Std::Int`"),
			},
		},
		"int capacity": {
			input: `
				var foo: ArrayList[Float] = [1.2]:9
			`,
		},
		"uint8 capacity": {
			input: `
				var foo: ArrayList[Float] = [1.2]:9u8
			`,
		},
		"invalid capacity": {
			input: `
				var foo: ArrayList[Float] = [1.2]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 2, 33), P(41, 2, 41)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer word array list": {
			input: `
				var foo = \w[foo bar]
				var a: ArrayList[String] = foo
			`,
		},
		"infer empty word array list": {
			input: `
				var foo = \w[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayList[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"word array list int capacity": {
			input: `
				var foo: ArrayList[String] = \w[1.2]:9
			`,
		},
		"word array list uint8 capacity": {
			input: `
				var foo: ArrayList[String] = \w[1.2]:9u8
			`,
		},
		"word array list invalid capacity": {
			input: `
				var foo: ArrayList[String] = \w[1.2]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 2, 42), P(44, 2, 44)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer symbol array list": {
			input: `
				var foo = \s[foo bar]
				var a: ArrayList[Symbol] = foo
			`,
		},
		"infer empty symbol array list": {
			input: `
				var foo = \s[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayList[Std::Symbol]` cannot be assigned to type `9`"),
			},
		},
		"symbol array list int capacity": {
			input: `
				var foo: ArrayList[Symbol] = \s[1.2]:9
			`,
		},
		"symbol array list uint8 capacity": {
			input: `
				var foo: ArrayList[Symbol] = \s[1.2]:9u8
			`,
		},
		"symbol array list invalid capacity": {
			input: `
				var foo: ArrayList[Symbol] = \s[1.2]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 2, 42), P(44, 2, 44)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer hex array list": {
			input: `
				var foo = \x[1f4 fff]
				var a: ArrayList[Int] = foo
			`,
		},
		"infer empty hex array list": {
			input: `
				var foo = \x[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayList[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"hex array list int capacity": {
			input: `
				var foo: ArrayList[Int] = \x[1f9]:9
			`,
		},
		"hex array list uint8 capacity": {
			input: `
				var foo: ArrayList[Int] = \x[fff]:9u8
			`,
		},
		"hex array list invalid capacity": {
			input: `
				var foo: ArrayList[Int] = \x[1ef]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 2, 39), P(41, 2, 41)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer bin array list": {
			input: `
				var foo = \b[111 100]
				var a: ArrayList[Int] = foo
			`,
		},
		"infer empty bin array list": {
			input: `
				var foo = \b[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayList[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"bin array list int capacity": {
			input: `
				var foo: ArrayList[Int] = \b[100]:9
			`,
		},
		"bin array list uint8 capacity": {
			input: `
				var foo: ArrayList[Int] = \x[111]:9u8
			`,
		},
		"bin array list invalid capacity": {
			input: `
				var foo: ArrayList[Int] = \b[101]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 2, 39), P(41, 2, 41)), "capacity must be an integer, got `9.2`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestHashMapLiteral(t *testing.T) {
	tests := testTable{
		"modifier if else": {
			input: `
				var a: bool = false
				var foo = { foo: 1, "bar" => 2.5 if a else 1 => "lol" }
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(100, 4, 16), P(102, 4, 18)), "type `Std::HashMap[Std::Symbol | Std::String | Std::Int, Std::Int | Std::Float | Std::String]` cannot be assigned to type `9`"),
			},
		},
		"modifier if": {
			input: `
				var a: bool = false
				var foo = { foo: 1, "bar" => 2.5 if a }
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(84, 4, 16), P(86, 4, 18)), "type `Std::HashMap[Std::Symbol | Std::String, Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"modifier if with narrowing": {
			input: `
				var a: String? = nil
				var foo = { a: 1, b: a.lol if a }
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 3, 26), P(55, 3, 30)), "method `lol` is not defined on type `Std::String`"),
			},
		},
		"truthy modifier if": {
			input: `
				var foo = { foo: 1, "bar" => 2.5 if 5 }
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(41, 2, 41), P(41, 2, 41)), "this condition will always have the same result since type `5` is truthy"),
				error.NewFailure(L("<main>", P(60, 3, 16), P(62, 3, 18)), "type `Std::HashMap[Std::Symbol | Std::String, Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"falsy modifier if": {
			input: `
				var foo = { foo: 1, "bar" => 2.5 if nil }
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(41, 2, 41), P(43, 2, 43)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L("<main>", P(25, 2, 25), P(36, 2, 36)), "unreachable code"),
				error.NewFailure(L("<main>", P(62, 3, 16), P(64, 3, 18)), "type `Std::HashMap[Std::Symbol, Std::Int]` cannot be assigned to type `9`"),
			},
		},

		"explicitly define hash map type": {
			input: `
				var foo: HashMap[Int, String] = {}
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(55, 3, 16), P(57, 3, 18)), "type `Std::HashMap[Std::Int, Std::String]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define map type": {
			input: `
				var foo: Map[Int, String] = {}
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 3, 16), P(53, 3, 18)), "type `Std::Map[Std::Int, Std::String]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define Map incompatible type": {
			input: `
				var foo: Map[Int, String] = { "" => 2.2 }
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 2, 33), P(45, 2, 45)), "type `Std::String` cannot be assigned to type `Std::Int`"),
				error.NewFailure(L("<main>", P(33, 2, 33), P(45, 2, 45)), "type `Std::Float` cannot be assigned to type `Std::String`"),
			},
		},
		"infer hash map": {
			input: `
				var foo = { foo: 1, bar: 2 }
				var a: HashMap[Symbol, Int] = foo
			`,
		},
		"infer hash map with different key and value types": {
			input: `
				var a = { foo: 1, "bar" => 2.2, 1 => "foo" }
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(65, 3, 16), P(65, 3, 16)), "type `Std::HashMap[Std::Symbol | Std::String | Std::Int, Std::Int | Std::Float | Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer empty hash map": {
			input: `
				var foo = {}
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 3, 16), P(35, 3, 18)), "type `Std::HashMap[any, any]` cannot be assigned to type `9`"),
			},
		},
		"int capacity": {
			input: `
				var foo: HashMap[Symbol, Float] = { foo: 1.2 }:9
			`,
		},
		"uint8 capacity": {
			input: `
				var foo: HashMap[Symbol, Float] = { foo: 1.2 }:9u8
			`,
		},
		"invalid capacity": {
			input: `
				var foo: HashMap[Symbol, Float] = { foo: 1.2 }:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 2, 39), P(54, 2, 54)), "capacity must be an integer, got `9.2`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestHashRecordLiteral(t *testing.T) {
	tests := testTable{
		"modifier if else": {
			input: `
				var a: bool = false
				var foo = %{ foo: 1, "bar" => 2.5 if a else 1 => "lol" }
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(101, 4, 16), P(103, 4, 18)), "type `Std::HashRecord[:foo | \"bar\" | 1, 1 | 2.5 | \"lol\"]` cannot be assigned to type `9`"),
			},
		},
		"modifier if": {
			input: `
				var a: bool = false
				var foo = %{ foo: 1, "bar" => 2.5 if a }
				var b: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 4, 16), P(87, 4, 18)), "type `Std::HashRecord[:foo | \"bar\", 1 | 2.5]` cannot be assigned to type `9`"),
			},
		},
		"modifier if with narrowing": {
			input: `
				var a: String? = nil
				var foo = %{ a: 1, b: a.lol if a }
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(52, 3, 27), P(56, 3, 31)), "method `lol` is not defined on type `Std::String`"),
			},
		},
		"truthy modifier if": {
			input: `
				var foo = %{ foo: 1, "bar" => 2.5 if 5 }
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(42, 2, 42), P(42, 2, 42)), "this condition will always have the same result since type `5` is truthy"),
				error.NewFailure(L("<main>", P(61, 3, 16), P(63, 3, 18)), "type `Std::HashRecord[:foo | \"bar\", 1 | 2.5]` cannot be assigned to type `9`"),
			},
		},
		"falsy modifier if": {
			input: `
				var foo = %{ foo: 1, "bar" => 2.5 if nil }
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(42, 2, 42), P(44, 2, 44)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L("<main>", P(26, 2, 26), P(37, 2, 37)), "unreachable code"),
				error.NewFailure(L("<main>", P(63, 3, 16), P(65, 3, 18)), "type `Std::HashRecord[:foo, 1]` cannot be assigned to type `9`"),
			},
		},

		"explicitly define hash record type": {
			input: `
				var foo: HashRecord[Int, String] = %{}
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(59, 3, 16), P(61, 3, 18)), "type `Std::HashRecord[Std::Int, Std::String]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define Record type": {
			input: `
				var foo: Record[Int, String] = %{}
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(55, 3, 16), P(57, 3, 18)), "type `Std::Record[Std::Int, Std::String]` cannot be assigned to type `9`"),
			},
		},
		"explicitly define Record incompatible type": {
			input: `
				var foo: Record[Int, String] = %{ "" => 2.2 }
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 2, 36), P(49, 2, 49)), "type `\"\"` cannot be assigned to type `Std::Int`"),
				error.NewFailure(L("<main>", P(36, 2, 36), P(49, 2, 49)), "type `2.2` cannot be assigned to type `Std::String`"),
			},
		},
		"infer hash record": {
			input: `
				var foo = %{ foo: 1, bar: 2 }
				var a: HashRecord[:foo | :bar, 1 | 2] = foo
			`,
		},
		"infer hash record with different key and value types": {
			input: `
				var a = %{ foo: 1, "bar" => 2.2, 1 => "foo" }
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(66, 3, 16), P(66, 3, 16)), "type `Std::HashRecord[:foo | \"bar\" | 1, 1 | 2.2 | \"foo\"]` cannot be assigned to type `9`"),
			},
		},
		"infer empty hash record": {
			input: `
				var foo = %{}
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 16), P(36, 3, 18)), "type `Std::HashRecord[any, any]` cannot be assigned to type `9`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestRegexLiteral(t *testing.T) {
	tests := testTable{
		"infer regex": {
			input: `
				var foo = %/str/
				var a: Regex = foo
			`,
		},
		"infer interpolated regex": {
			input: `
				var foo = %/${1} str/
				var b: Regex = foo
			`,
		},
		"interpolate unconvertible value": {
			input: `
				class Foo; end
				var foo = %/${Foo()} str/
				var b: Regex = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(38, 3, 19), P(42, 3, 23)), "type `Foo` does not implement interface `Std::StringConvertible`:\n\n  - missing method `Std::StringConvertible.:to_string` with signature: `def to_string(): Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestStringLiteral(t *testing.T) {
	tests := testTable{
		"infer raw string": {
			input: `
				var foo = 'str'
				var a: String = foo
			`,
		},
		"assign string literal to String": {
			input: "var foo: String = 'str'",
		},
		"assign string literal to matching literal type": {
			input: "var foo: 'str' = 'str'",
		},
		"assign string literal to non matching literal type": {
			input: "var foo: 'str' = 'foo'",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `\"foo\"` cannot be assigned to type `\"str\"`"),
			},
		},
		"infer double quoted string": {
			input: `
				var foo = "str"
				var b: String = foo
			`,
		},
		"infer interpolated string": {
			input: `
				var foo = "${1} str #{5.2}"
				var b: String = foo
			`,
		},
		"interpolate unconvertible value": {
			input: `
				class Foo; end
				var foo = "${Foo()} str"
				var b: String = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(37, 3, 18), P(41, 3, 22)), "type `Foo` does not implement interface `Std::StringConvertible`:\n\n  - missing method `Std::StringConvertible.:to_string` with signature: `def to_string(): Std::String`"),
			},
		},
		"interpolate uninspectable value": {
			input: `
				class Foo < nil; end
				var foo = "#{Foo()} str"
				var b: String = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 3, 18), P(47, 3, 22)), "type `Foo` does not implement interface `Std::Inspectable`:\n\n  - missing method `Std::Inspectable.:inspect` with signature: `def inspect(): Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSymbolLiteral(t *testing.T) {
	tests := testTable{
		"infer simple symbol": {
			input: "var foo = :str",
		},
		"infer double quoted symbol": {
			input: `var foo = :"str"`,
		},
		"infer interpolated symbol": {
			input: `var foo = :"${1} str #{5.2}"`,
		},
		"assign symbol literal to Symbol": {
			input: "var foo: Symbol = :symb",
		},
		"assign symbol literal to matching literal type": {
			input: "var foo: :symb = :symb",
		},
		"assign symbol literal to non matching literal type": {
			input: "var foo: :symb = :foob",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `:foob` cannot be assigned to type `:symb`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestCharLiteral(t *testing.T) {
	tests := testTable{
		"infer raw char": {
			input: "var foo = r`s`",
		},
		"infer char": {
			input: "var foo = `\\n`",
		},
		"assign char literal to Char": {
			input: "var foo: Char = `f`",
		},
		"assign char literal to matching literal type": {
			input: "var foo: `f` = `f`",
		},
		"assign char literal to non matching literal type": {
			input: "var foo: `b` = `f`",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type ``f`` cannot be assigned to type ``b``"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestIntLiteral(t *testing.T) {
	tests := testTable{
		"infer int": {
			input: "var foo = 1234",
		},
		"assign int literal to Int": {
			input: "var foo: Int = 12345678",
		},
		"assign int literal to matching literal type": {
			input: "var foo: 12345 = 12345",
		},
		"assign int literal to non matching literal type": {
			input: "var foo: 23456 = 12345",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `12345` cannot be assigned to type `23456`"),
			},
		},
		"infer int64": {
			input: "var foo = 1i64",
		},
		"infer int32": {
			input: "var foo = 1i32",
		},
		"infer int16": {
			input: "var foo = 1i16",
		},
		"infer int8": {
			input: "var foo = 12i8",
		},
		"infer uint64": {
			input: "var foo = 1u64",
		},
		"infer uint32": {
			input: "var foo = 1u32",
		},
		"infer uint16": {
			input: "var foo = 1u16",
		},
		"infer uint8": {
			input: "var foo = 12u8",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestRangeLiteral(t *testing.T) {
	tests := testTable{
		// closed range
		"infer type with the same argument": {
			input: "var a: 9 = 1...5",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 1, 12), P(15, 1, 16)), "type `Std::ClosedRange[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"infer type with different argument types": {
			input: "var a: 9 = 1...2.5",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 1, 12), P(17, 1, 18)), "type `Std::ClosedRange[Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"infer type with incompatible argument types": {
			input: "var a: 9 = 1...'c'",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 1, 12), P(17, 1, 18)), "type `Std::Int` does not implement interface `Std::Comparable[Std::Int | Std::String]`:\n\n  - incorrect implementation of `Std::Comparable.:<`\n      is:        `def <(other: Std::CoercibleNumeric): Std::Bool`\n      should be: `def <(other: Std::Int | Std::String): bool`\n  - incorrect implementation of `Std::Comparable.:<=`\n      is:        `def <=(other: Std::CoercibleNumeric): Std::Bool`\n      should be: `def <=(other: Std::Int | Std::String): bool`\n  - incorrect implementation of `Std::Comparable.:<=>`\n      is:        `def <=>(other: Std::CoercibleNumeric): Std::Int`\n      should be: `def <=>(other: Std::Int | Std::String): Std::Int?`\n  - incorrect implementation of `Std::Comparable.:>`\n      is:        `def >(other: Std::CoercibleNumeric): Std::Bool`\n      should be: `def >(other: Std::Int | Std::String): bool`\n  - incorrect implementation of `Std::Comparable.:>=`\n      is:        `def >=(other: Std::CoercibleNumeric): Std::Bool`\n      should be: `def >=(other: Std::Int | Std::String): bool`"),
				error.NewFailure(L("<main>", P(11, 1, 12), P(17, 1, 18)), "type Std::Int | Std::String is not comparable and cannot be used in range literals"),
				error.NewFailure(L("<main>", P(11, 1, 12), P(17, 1, 18)), "type `Std::ClosedRange[untyped]` cannot be assigned to type `9`"),
			},
		},
		"call iter on iterable range": {
			input: "var a: 9 = (1...5).iter",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 1, 13), P(22, 1, 23)), "type `Std::ClosedRange::Iterator[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"call iter on not iterable range": {
			input: "(1.2...2.5).iter",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(1, 1, 2), P(15, 1, 16)), "method `iter` is not defined on type `Std::ClosedRange[Std::Float]`"),
			},
		},
		"call iter on mixed range": {
			input: "(1...2.5).iter",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(1, 1, 2), P(13, 1, 14)), "method `iter` is not defined on type `Std::ClosedRange[Std::Int | Std::Float]`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestFloatLiteral(t *testing.T) {
	tests := testTable{
		"assign float to comparable": {
			input: "var a: Comparable[Float] = 2.5",
		},
		"infer float": {
			input: "var foo = 12.5",
		},
		"assign float literal to Float": {
			input: "var foo: Float = 1234.6",
		},
		"assign float literal to matching literal type": {
			input: "var foo: 12.45 = 12.45",
		},
		"assign Float literal to non matching literal type": {
			input: "var foo: 23.56 = 12.45",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `12.45` cannot be assigned to type `23.56`"),
			},
		},
		"infer float64": {
			input: "var foo = 1f64",
		},
		"infer float32": {
			input: "var foo = 1f32",
		},
		"infer big float": {
			input: "var foo = 12bf",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestBoolLiteral(t *testing.T) {
	tests := testTable{
		"infer true": {
			input: "var foo = true",
		},
		"assign true literal to True": {
			input: "var foo: True = true",
		},
		"assign true literal to Bool": {
			input: "var foo: Bool = true",
		},
		"assign true literal to matching literal type": {
			input: "var foo: true = true",
		},
		"assign true literal to non matching literal type": {
			input: "var foo: false = true",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 1, 18), P(20, 1, 21)), "type `true` cannot be assigned to type `false`"),
			},
		},
		"infer false": {
			input: "var foo = false",
		},
		"assign false literal to False": {
			input: "var foo: False = false",
		},
		"assign false literal to Bool": {
			input: "var foo: Bool = false",
		},
		"assign false literal to matching literal type": {
			input: "var foo: false = false",
		},
		"assign false literal to non matching literal type": {
			input: "var foo: true = false",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 1, 17), P(20, 1, 21)), "type `false` cannot be assigned to type `true`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestNilLiteral(t *testing.T) {
	tests := testTable{
		"infer nil": {
			input: "var foo = nil",
		},
		"assign nil literal to Nil": {
			input: "var foo: Nil = nil",
		},
		"assign nil literal to matching literal type": {
			input: "var foo: nil = nil",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
