package checker_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestGoExpression(t *testing.T) {
	tests := testTable{
		"has its own scope": {
			input: `
				a := 5
				go
					a + "foo"
					b := 5
				end
				b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(28, 4, 10), P(32, 4, 14)), "expected type `Std::Int` for parameter `other` in call to `Std::Int.:+`, got type `\"foo\"`"),
				diagnostic.NewFailure(L("<main>", P(58, 7, 5), P(58, 7, 5)), "undefined local `b`"),
			},
		},
		"returns a Thread": {
			input: `
				var a: nil = go println("foo")
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(18, 2, 18), P(34, 2, 34)), "type `Std::Thread` cannot be assigned to type `nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestDoCatchExpression(t *testing.T) {
	tests := testTable{
		"do and finally have their own scopes": {
			input: `
				a := 5
				do
					a + "foo"
					b := 5
				finally
					a - "bar"
					b := 2.5
				end
				b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(67, 7, 10), P(71, 7, 14)), "expected type `Std::Int` for parameter `other` in call to `Std::Int.:-`, got type `\"bar\"`"),
				diagnostic.NewFailure(L("<main>", P(28, 4, 10), P(32, 4, 14)), "expected type `Std::Int` for parameter `other` in call to `Std::Int.:+`, got type `\"foo\"`"),
				diagnostic.NewFailure(L("<main>", P(99, 10, 5), P(99, 10, 5)), "undefined local `b`"),
			},
		},
		"catches have their own scopes": {
			input: `
				a := 5
				do
					a + "foo"
					b := 5
				catch String()
					a - "bar"
					b := 2.5
				catch Char()
					a * "bar"
					b := 2u8
				end
				b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(74, 7, 10), P(78, 7, 14)), "expected type `Std::Int` for parameter `other` in call to `Std::Int.:-`, got type `\"bar\"`"),
				diagnostic.NewFailure(L("<main>", P(120, 10, 10), P(124, 10, 14)), "expected type `Std::Int` for parameter `other` in call to `Std::Int.:*`, got type `\"bar\"`"),
				diagnostic.NewFailure(L("<main>", P(28, 4, 10), P(32, 4, 14)), "expected type `Std::Int` for parameter `other` in call to `Std::Int.:+`, got type `\"foo\"`"),
				diagnostic.NewFailure(L("<main>", P(152, 13, 5), P(152, 13, 5)), "undefined local `b`"),
			},
		},
		"checks invalid patterns": {
			input: `
				do
					println(5)
				catch String(length) && Int()
					println(length)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(34, 4, 11), P(56, 4, 33)), "this pattern is impossible to satisfy"),
			},
		},
		"throw with wider catch": {
			input: `
				do
					throw :foo
				catch Symbol()
					println "lol"
				end
			`,
		},
		"catch with stack trace": {
			input: `
				do
					throw :foo
				catch Symbol(), st
					println "lol"
					var a: nil = st
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(84, 6, 19), P(85, 6, 20)), "type `Std::StackTrace` cannot be assigned to type `nil`"),
			},
		},
		"method call with wider catch": {
			input: `
				def foo! :foo
					throw :foo
				end

				do
					foo()
				catch Symbol()
					println "lol"
				end
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestTryExpression(t *testing.T) {
	tests := testTable{
		"try with a value": {
			input: `
				a := 6
				var b: -1 = try a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(28, 3, 17), P(32, 3, 21)), "unnecessary `try`, the expression does not throw a checked error"),
				diagnostic.NewFailure(L("<main>", P(28, 3, 17), P(32, 3, 21)), "type `Std::Int` cannot be assigned to type `-1`"),
			},
		},
		"try with a function call": {
			input: `
				def foo: Int
					3
				end

				var b: -1 = try foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(50, 6, 17), P(58, 6, 25)), "unnecessary `try`, the expression does not throw a checked error"),
				diagnostic.NewFailure(L("<main>", P(50, 6, 17), P(58, 6, 25)), "type `Std::Int` cannot be assigned to type `-1`"),
			},
		},
		"try with a function call that throws a checked error": {
			input: `
				def foo: Int ! Symbol
					throw :lol
				end

				var b: -1 = try foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(68, 6, 17), P(76, 6, 25)), "type `Std::Int` cannot be assigned to type `-1`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestMustExpression(t *testing.T) {
	tests := testTable{
		"must with a non-nilable value": {
			input: `
				a := 6
				var b: -1 = must a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(28, 3, 17), P(33, 3, 22)), "unnecessary `must`, type `Std::Int` is not nilable"),
				diagnostic.NewFailure(L("<main>", P(28, 3, 17), P(33, 3, 22)), "type `Std::Int` cannot be assigned to type `-1`"),
			},
		},
		"must with a nilable value": {
			input: `
				var a: Int? = 6
				var b: -1 = must a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(37, 3, 17), P(42, 3, 22)), "type `Std::Int` cannot be assigned to type `-1`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestAsExpression(t *testing.T) {
	tests := testTable{
		"valid type downcast": {
			input: `
				var a: String | Float = .2
				var b: 9 = a as String
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(47, 3, 16), P(57, 3, 26)), "type `Std::String` cannot be assigned to type `9`"),
			},
		},
		"valid type upcast": {
			input: `
				var a: String = "foo"
				var b: 9 = a as Value
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(42, 3, 16), P(51, 3, 25)), "type `Std::Value` cannot be assigned to type `9`"),
			},
		},
		"invalid cast": {
			input: `
				var a: String | Float = .2
				var b: 9 = a as Int
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(47, 3, 16), P(54, 3, 23)), "cannot cast type `Std::String | Std::Float` to type `Std::Int`"),
			},
		},
		"invalid type in cast": {
			input: `
				typedef Foo = "foo"
				var a = 5
				var b: 9 = a as Foo
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(59, 4, 21), P(61, 4, 23)), "only classes and mixins are allowed in `as` type casts"),
				diagnostic.NewFailure(L("<main>", P(54, 4, 16), P(61, 4, 23)), "cannot cast type `Std::Int` to type `Foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestThrowExpression(t *testing.T) {
	tests := testTable{
		"throw without value": {
			input: `
				throw
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(5, 2, 5), P(9, 2, 9)), "thrown value of type `Std::Error` must be caught"),
			},
		},
		"throw without catch": {
			input: `
				throw :foo
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(5, 2, 5), P(14, 2, 14)), "thrown value of type `:foo` must be caught"),
			},
		},
		"throw with different catch": {
			input: `
				do
					throw :foo
				catch String() as str
					println str
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(13, 3, 6), P(22, 3, 15)), "thrown value of type `:foo` must be caught"),
			},
		},
		"throw with matching catch": {
			input: `
				do
					throw :foo
				catch :foo
					println "lol"
				end
			`,
		},
		"throw with non-exhaustive catch": {
			input: `
				do
					var a: Symbol = :foo
					throw a
				catch :foo
					println "lol"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(39, 4, 6), P(45, 4, 12)), "thrown value of type `Std::Symbol` must be caught"),
			},
		},
		"throw with wider catch": {
			input: `
				do
					throw :foo
				catch Symbol()
					println "lol"
				end
			`,
		},
		"call receiverless method that throws": {
			input: `
				def foo! Symbol
					throw :foo
				end

				foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(50, 6, 5), P(52, 6, 7)), "thrown value of type `Std::Symbol` must be caught"),
			},
		},
		"call receiverless method that throws and catch": {
			input: `
				def foo! Symbol
					throw :foo
				end

				do
					foo()
				catch Symbol() as sym
					println sym
				end
			`,
		},
		"call method that throws": {
			input: `
				module Foo
					def foo! Symbol
						throw :foo
					end
				end

				Foo.foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(80, 8, 9), P(82, 8, 11)), "thrown value of type `Std::Symbol` must be caught"),
			},
		},
		"call method that throws and catch": {
			input: `
				module Foo
					def foo! Symbol
						throw :foo
					end
				end

				do
					Foo.foo()
				catch Symbol() as sym
					println sym
				end
			`,
		},

		"call constructor that throws": {
			input: `
				class Foo
					init! Symbol
						throw :foo
					end
				end

				Foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(72, 8, 5), P(76, 8, 9)), "thrown value of type `Std::Symbol` must be caught"),
			},
		},
		"call constructor that throws and catch": {
			input: `
				class Foo
					init! Symbol
						throw :foo
					end
				end

				do
					Foo()
				catch Symbol() as sym
					println sym
				end
			`,
		},
		"call implicit generic constructor that throws": {
			input: `
				class Foo[V]
					init(a: V)! V
						throw a
					end
				end

				Foo("lol")
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(73, 8, 5), P(82, 8, 14)), "thrown value of type `Std::String` must be caught"),
			},
		},
		"call explicit generic constructor that throws": {
			input: `
				class Foo[V]
					init(a: V)! V
						throw a
					end
				end

				Foo::[String?]("lol")
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(73, 8, 5), P(93, 8, 25)), "thrown value of type `Std::String?` must be caught"),
			},
		},

		"call closure that throws": {
			input: `
				foo := ||! Symbol -> throw :foo
				foo.()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(41, 3, 5), P(46, 3, 10)), "thrown value of type `Std::Symbol` must be caught"),
			},
		},
		"call closure that throws and catch": {
			input: `
				foo := ||! Symbol -> throw :foo

				do
					foo.()
				catch Symbol() as sym
					println sym
				end
			`,
		},

		"method body - throw without catch": {
			input: `
				def bar
					throw :foo
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(18, 3, 6), P(27, 3, 15)), "thrown value of type `:foo` must be caught or added to the signature of the function `! :foo`"),
			},
		},
		"method body - throw with different catch": {
			input: `
				def bar
					do
						throw :foo
					catch String() as str
						println str
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(27, 4, 7), P(36, 4, 16)), "thrown value of type `:foo` must be caught or added to the signature of the function `! :foo`"),
			},
		},
		"method body - throw with matching catch": {
			input: `
				def bar
					do
						throw :foo
					catch :foo
						println "lol"
					end
				end
			`,
		},
		"method body - throw with non-exhaustive catch": {
			input: `
				def bar
					do
						var a: Symbol = :foo
						throw a
					catch :foo
						println "lol"
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(54, 5, 7), P(60, 5, 13)), "thrown value of type `Std::Symbol` must be caught or added to the signature of the function `! Std::Symbol`"),
			},
		},
		"method body - throw with wider catch": {
			input: `
				def bar
					do
						throw :foo
					catch Symbol()
						println "lol"
					end
				end
			`,
		},
		"method body - throw with different method throw type": {
			input: `
				def bar! String
					throw :foo
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(26, 3, 6), P(35, 3, 15)), "thrown value of type `:foo` must be caught or added to the signature of the function `! Std::String | :foo`"),
			},
		},
		"method body - throw with matching method throw type": {
			input: `
				def bar! :foo
					throw :foo
				end
			`,
		},
		"method body - throw with non-exhaustive method throw type": {
			input: `
				def bar! :foo
					var a: Symbol = :foo
					throw a
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(50, 4, 6), P(56, 4, 12)), "thrown value of type `Std::Symbol` must be caught or added to the signature of the function `! Std::Symbol`"),
			},
		},
		"method body - throw with wider method throw type": {
			input: `
				def bar! Symbol
					throw :foo
				end
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSwitchExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				a := 5
				var b: String? = nil
				switch b
				case String()
					var c: 9 = a
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(84, 6, 17), P(84, 6, 17)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"narrows variables": {
			input: `
				var a: String? = nil
				switch a
				case String()
					var c: 1 = a
				case nil
					var d: 2 = a
				case false
					a = 3
				else
					var f: 4 = a
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(73, 5, 17), P(73, 5, 17)), "type `Std::String` cannot be assigned to type `1`"),
				diagnostic.NewFailure(L("<main>", P(104, 7, 17), P(104, 7, 17)), "type `nil` cannot be assigned to type `2`"),
				diagnostic.NewFailure(L("<main>", P(115, 8, 10), P(119, 8, 14)), "type `Std::String?` cannot ever match type `false`"),
				diagnostic.NewFailure(L("<main>", P(130, 9, 10), P(130, 9, 10)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(157, 11, 17), P(157, 11, 17)), "type `Std::String | nil` cannot be assigned to type `4`"),
			},
		},
		"narrows variable declarations": {
			input: `
				var a: String? = nil
				switch var b: String | Float | nil = a
				case String()
					var c: 1 = b
				case nil
					var d: 2 = b
				case false
					b = 3
				else
					var f: 4 = b
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(103, 5, 17), P(103, 5, 17)), "type `Std::String` cannot be assigned to type `1`"),
				diagnostic.NewFailure(L("<main>", P(134, 7, 17), P(134, 7, 17)), "type `nil` cannot be assigned to type `2`"),
				diagnostic.NewFailure(L("<main>", P(145, 8, 10), P(149, 8, 14)), "type `Std::String?` cannot ever match type `false`"),
				diagnostic.NewFailure(L("<main>", P(160, 9, 10), P(160, 9, 10)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(187, 11, 17), P(187, 11, 17)), "type `Std::String | nil` cannot be assigned to type `4`"),
			},
		},
		"narrows value declarations": {
			input: `
				var a: String? = nil
				switch val b: String | Float | nil = a
				case String()
					var c: 1 = b
				case nil
					var d: 2 = b
				case false
					b = 3
				else
					var f: 4 = b
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(103, 5, 17), P(103, 5, 17)), "type `Std::String` cannot be assigned to type `1`"),
				diagnostic.NewFailure(L("<main>", P(134, 7, 17), P(134, 7, 17)), "type `nil` cannot be assigned to type `2`"),
				diagnostic.NewFailure(L("<main>", P(145, 8, 10), P(149, 8, 14)), "type `Std::String?` cannot ever match type `false`"),
				diagnostic.NewFailure(L("<main>", P(156, 9, 6), P(156, 9, 6)), "local value `b` cannot be reassigned"),
				diagnostic.NewFailure(L("<main>", P(160, 9, 10), P(160, 9, 10)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(187, 11, 17), P(187, 11, 17)), "type `Std::String | nil` cannot be assigned to type `4`"),
			},
		},
		"narrows variable assignments": {
			input: `
				var a: String | Float | nil = nil
				def b: String? then nil

				switch a = b()
				case String()
					var c: 1 = a
				case nil
					var d: 2 = a
				case Float()
					a = 3
				else
					var f: 4 = a
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(121, 7, 17), P(121, 7, 17)), "type `Std::String` cannot be assigned to type `1`"),
				diagnostic.NewFailure(L("<main>", P(152, 9, 17), P(152, 9, 17)), "type `nil` cannot be assigned to type `2`"),
				diagnostic.NewFailure(L("<main>", P(163, 10, 10), P(169, 10, 16)), "type `Std::String?` cannot ever match type `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(180, 11, 10), P(180, 11, 10)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(207, 13, 17), P(207, 13, 17)), "type `Std::String | nil` cannot be assigned to type `4`"),
			},
		},
		"returns a union of last expression types in each block": {
			input: `
				var a: String? = nil
				var c: 0 =
				  switch a = b()
				  case String()
					  5
				  case nil
					  2.5
				  else
					  "else"
				  end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(58, 4, 18), P(58, 4, 18)), "method `b` is not defined on type `Std::Object`"),
				diagnostic.NewFailure(L("<main>", P(47, 4, 7), P(150, 11, 9)), "type `5 | 2.5 | \"else\"` cannot be assigned to type `0`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestBreakExpression(t *testing.T) {
	tests := testTable{
		"the return type is never": {
			input: `
				loop
					a := break
					a = 4
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(35, 4, 10), P(35, 4, 10)), "type `4` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(31, 4, 6), P(35, 4, 10)), "unreachable code"),
			},
		},
		"outside of a loop": {
			input: `break`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(4, 1, 5)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"outside of a loop with a nonexistent label": {
			input: `break[foo]`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(8, 1, 9)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"with a nonexistent label": {
			input: `
				loop
					break[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(15, 3, 6), P(23, 3, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"with a displaced label": {
			input: `
				$foo: loop; end
				loop
					break[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(35, 4, 6), P(43, 4, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
				diagnostic.NewWarning(L("<main>", P(25, 3, 5), P(52, 5, 7)), "unreachable code"),
			},
		},
		"with a valid label": {
			input: `
				$foo: loop
					loop
						break[foo]
					end
				end
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestContinueExpression(t *testing.T) {
	tests := testTable{
		"the return type is never": {
			input: `
				loop
					a := continue
					a = 4
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(38, 4, 10), P(38, 4, 10)), "type `4` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(34, 4, 6), P(38, 4, 10)), "unreachable code"),
			},
		},
		"outside of a loop": {
			input: `continue`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(7, 1, 8)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"outside of a loop with a nonexistent label": {
			input: `continue[foo]`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(11, 1, 12)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"with a nonexistent label": {
			input: `
				loop
					continue[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(15, 3, 6), P(26, 3, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"with a displaced label": {
			input: `
				$foo: loop; end
				loop
					continue[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(35, 4, 6), P(46, 4, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
				diagnostic.NewWarning(L("<main>", P(25, 3, 5), P(55, 5, 7)), "unreachable code"),
			},
		},
		"with a valid label": {
			input: `
				$foo: loop
					loop
						continue[foo]
					end
				end
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestReturnExpression(t *testing.T) {
	tests := testTable{
		"the return type is never": {
			input: `
				a := return
				a = 4
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(25, 3, 9), P(25, 3, 9)), "type `4` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(21, 3, 5), P(25, 3, 9)), "unreachable code"),
			},
		},
		"warn about values returned in the top level": {
			input: `
				return 4
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(12, 2, 12), P(12, 2, 12)), "values returned in void context will be ignored"),
			},
		},
		"warn about values returned in void methods": {
			input: `
				def foo
					return 4
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(25, 3, 13), P(25, 3, 13)), "values returned in void context will be ignored"),
			},
		},
		"accept matching return type": {
			input: `
				def foo: String
					return "foo"
				end
			`,
		},
		"invalid return type": {
			input: `
				def foo: String
					return 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(26, 3, 6), P(33, 3, 13)), "type `2` cannot be assigned to type `Std::String`"),
			},
		},
		"accept matching return type in a generator": {
			input: `
				def *foo: String
					return "foo"
				end
			`,
		},
		"invalid return type in a generator": {
			input: `
				def *foo: String
					return 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(27, 3, 6), P(34, 3, 13)), "type `2` cannot be assigned to type `Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestYieldExpression(t *testing.T) {
	tests := testTable{
		"cannot yield in the top level": {
			input: `
				yield 4
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(5, 2, 5), P(11, 2, 11)), "yield cannot be used outside of generators"),
				diagnostic.NewWarning(L("<main>", P(11, 2, 11), P(11, 2, 11)), "values yielded in void context will be ignored"),
			},
		},
		"warn about values yielded in void methods": {
			input: `
				def *foo
					yield 4
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(25, 3, 12), P(25, 3, 12)), "values yielded in void context will be ignored"),
			},
		},
		"cannot yield in closures in generators": {
			input: `
				def* iter: String
					|str: String| -> yield "foo"

					"bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(45, 3, 23), P(55, 3, 33)), "yield cannot be used outside of generators"),
			},
		},
		"accept matching yield type": {
			input: `
				def *foo: String
					yield "foo"
					"bar"
				end
			`,
		},
		"invalid yield type": {
			input: `
				def *foo: String
					yield 2

					"bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(27, 3, 6), P(33, 3, 12)), "type `2` cannot be assigned to type `Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestAwaitExpression(t *testing.T) {
	tests := testTable{
		"cannot await non-promise": {
			input: `
				await 4
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(11, 2, 11), P(11, 2, 11)), "only promises can be awaited"),
			},
		},
		"await returns the result type of the promise": {
			input: `
				var a: Promise[Int] = loop; end
				var b: nil = a.await
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(54, 3, 18), P(60, 3, 24)), "type `Std::Int` cannot be assigned to type `nil`"),
				diagnostic.NewWarning(L("<main>", P(41, 3, 5), P(60, 3, 24)), "unreachable code"),
			},
		},
		"await throws if the promise has a throw type": {
			input: `
				var a: Promise[Int, String] = loop; end
				a.await
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(49, 3, 5), P(55, 3, 11)), "thrown value of type `Std::String` must be caught"),
				diagnostic.NewWarning(L("<main>", P(49, 3, 5), P(55, 3, 11)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestDoExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				a := 5
				do
					var b: Int = a
				end
			`,
		},
		"returns the last expression": {
			input: `
				a := 2
				var b: Int = do
					"foo" + "bar"
					a + 2
				end
			`,
		},
		"returns nil when empty": {
			input: `
				var b: nil = do; end
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestNumericForExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				fornum ;;
					var b: Int? = a
				end
			`,
		},
		"use variables defined in the header": {
			input: `
				fornum i := 0; i < 8; i++
					var b: Int = i
				end
			`,
		},
		"typecheck the header and body": {
			input: `
				fornum a; b; c
					d
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(12, 2, 12), P(12, 2, 12)), "undefined local `a`"),
				diagnostic.NewFailure(L("<main>", P(15, 2, 15), P(15, 2, 15)), "undefined local `b`"),
				diagnostic.NewFailure(L("<main>", P(18, 2, 18), P(18, 2, 18)), "undefined local `c`"),
				diagnostic.NewFailure(L("<main>", P(25, 3, 6), P(25, 3, 6)), "undefined local `d`"),
			},
		},
		"returns never if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(81, 7, 9), P(81, 7, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(77, 7, 5), P(81, 7, 9)), "unreachable code"),
			},
		},
		"returns never when there is no condition": {
			input: `
				a := 2
				b := fornum ;;
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(77, 7, 9), P(77, 7, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(73, 7, 5), P(77, 7, 9)), "unreachable code"),
			},
		},
		"returns nil if condition is falsy": {
			input: `
				a := 2
				var b: nil = fornum ;false;
					"foo" + "bar"
					a + 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(37, 3, 26), P(41, 3, 30)), "this loop will never execute since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(49, 4, 6), P(62, 4, 19)), "unreachable code"),
			},
		},
		"returns a nilable body type if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(72, 5, 7)), "type `Std::String?` cannot be assigned to type `8`"),
			},
		},
		"cannot use void in the condition": {
			input: `
				def foo; end
				fornum foo(); foo(); foo()
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 19), P(40, 3, 23)), "cannot use type `void` as a value in this context"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				fornum ;a;
					var b: true = a
				end
			`,
		},

		"returns nil with a naked break if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					c := false
					if c
						break
					end
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(109, 10, 9), P(109, 10, 9)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"returns the value given to break if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					c := false
					if c
						break "foo" + "bar"
					end
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(123, 10, 9), P(123, 10, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"returns nil with a break if condition is falsy": {
			input: `
				a := 2
				var b: nil = fornum ;false;
					c := false
					if c
						break "foo" + "bar"
					end
					a + 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(37, 3, 26), P(41, 3, 30)), "this loop will never execute since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(49, 4, 6), P(59, 4, 16)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					c := false
					if c
						break
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(119, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					c := false
					if c
						break 2.5
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(123, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"break from a nested labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: fornum ;a;
					var b: Int? = 9
					fornum ;b;
						break[foo] 2.5
						"foo" + "bar"
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(117, 7, 7), P(129, 7, 19)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(146, 9, 7)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},

		"returns never with a naked continue if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					c := false
					if c
						continue
					end
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(112, 10, 9), P(112, 10, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(108, 10, 5), P(112, 10, 9)), "unreachable code"),
			},
		},
		"returns never with continue if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					continue "foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewWarning(L("<main>", P(68, 5, 6), P(72, 5, 10)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(90, 7, 9), P(90, 7, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(86, 7, 5), P(90, 7, 9)), "unreachable code"),
			},
		},
		"returns nil with a continue if condition is falsy": {
			input: `
				a := 2
				var b: nil = fornum ;false;
					c := false
					if c
						continue "foo" + "bar"
					end
					a + 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(37, 3, 26), P(41, 3, 30)), "this loop will never execute since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(49, 4, 6), P(59, 4, 16)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					var c: Int? = 1
					if c
						continue
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(127, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					c := false
					if c
						continue 2.5
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(126, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"continue a parent labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: fornum ;a;
					var b: Int? = 9
					fornum ;b;
						continue[foo] 2.5
						"foo" + "bar"
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(120, 7, 7), P(132, 7, 19)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(149, 9, 7)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestForInExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				for i in [1, 2, 3]
					var b: Int? = a
				end
			`,
		},
		"use variables defined in the header": {
			input: `
				for i in [1, 2, 3]
					var b: Int = i
				end
			`,
		},
		"typecheck the header and body": {
			input: `
				for a in b
					c
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(14, 2, 14), P(14, 2, 14)), "undefined local `b`"),
				diagnostic.NewFailure(L("<main>", P(21, 3, 6), P(21, 3, 6)), "undefined local `c`"),
			},
		},
		"cannot use void in the condition": {
			input: `
				def foo; end
				for i in foo()
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 14), P(35, 3, 18)), "type `void` cannot be iterated over, it does not implement `Std::PrimitiveIterable[any, any]`"),
			},
		},

		"returns nil": {
			input: `
				var a = [1, 2, 3]
				b := for i in a
					i
					5
				end
				var c: 9 = b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(80, 7, 16), P(80, 7, 16)), "type `nil` cannot be assigned to type `9`"),
			},
		},

		"returns nil with a naked break": {
			input: `
				a := 2
				b := for i in [1, 2, 3]
					c := false
					if c
						break
					end
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(114, 10, 9), P(114, 10, 9)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"returns the nilable value given to conditional break": {
			input: `
				a := 2
				b := for i in [1, 2, 3]
					c := false
					if c
						break "foo" + "bar"
					end
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(128, 10, 9), P(128, 10, 9)), "type `3` cannot be assigned to type `Std::String?`"),
			},
		},
		"returns the nilable value given to break": {
			input: `
				a := 2
				b := for i in [1, 2, 3]
					break "foo" + "bar"
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(81, 6, 9), P(81, 6, 9)), "type `3` cannot be assigned to type `Std::String?`"),
			},
		},

		"break from a nested labeled loop": {
			input: `
				var b: 8 = $foo: for i in 2...10
					for j in [9, 2, 6]
						break[foo] 2.5
						"foo" + "bar"
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(89, 5, 7), P(101, 5, 19)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(16, 2, 16), P(118, 7, 7)), "type `2.5?` cannot be assigned to type `8`"),
			},
		},

		"returns nil with naked continue": {
			input: `
				var b: 8 = for i in [1, 2, 3]
					var c: Int? = 1
					if c
						continue
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(16, 2, 16), P(115, 8, 7)), "type `nil` cannot be assigned to type `8`"),
			},
		},
		"returns nil with continue": {
			input: `
				var b: 8 = for c in 2...10
					d := false
					if d
						continue 2.5
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(16, 2, 16), P(111, 8, 7)), "type `nil` cannot be assigned to type `8`"),
			},
		},
		"continue a parent labeled loop": {
			input: `
				var b: 8 = $foo: for i in 1...20
					for j in 7...30
						continue[foo] 2.5
						"foo" + "bar"
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(89, 5, 7), P(101, 5, 19)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(16, 2, 16), P(118, 7, 7)), "type `nil` cannot be assigned to type `8`"),
			},
		},

		"unmatchable pattern": {
			input: `
				for [1, i] in 1...20
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(9, 2, 9), P(14, 2, 14)), "type `Std::Int` cannot ever match type `Std::List[any]`"),
			},
		},
		"valid object pattern": {
			input: `
				for Pair(key, value) in { foo: "bar", baz: "lol" }
					var a: 9 = key
					var b: nil = value
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(72, 3, 17), P(74, 3, 19)), "type `Std::Symbol` cannot be assigned to type `9`"),
				diagnostic.NewFailure(L("<main>", P(94, 4, 19), P(98, 4, 23)), "type `Std::String` cannot be assigned to type `nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestModifierForInExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				var a: Int = 5
				println(a) for i in [1, 2, 3]
			`,
		},
		"use variables defined in the header": {
			input: `
				println(i) for i in [1, 2, 3]
			`,
		},
		"typecheck the header and body": {
			input: `
				c for a in b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(16, 2, 16), P(16, 2, 16)), "undefined local `b`"),
				diagnostic.NewFailure(L("<main>", P(5, 2, 5), P(5, 2, 5)), "undefined local `c`"),
			},
		},
		"cannot use void in the condition": {
			input: `
				def foo; end
				println(i) for i in foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(42, 3, 25), P(46, 3, 29)), "type `void` cannot be iterated over, it does not implement `Std::PrimitiveIterable[any, any]`"),
			},
		},

		"returns nil": {
			input: `
				var a = [1, 2, 3]
				b := (i for i in a)
				var c: 9 = b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(62, 4, 16), P(62, 4, 16)), "type `nil` cannot be assigned to type `9`"),
			},
		},

		"returns the nilable value given to break": {
			input: `
				a := 2
				b := (break "foo" + "bar" for i in [1, 2, 3])
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(70, 4, 9), P(70, 4, 9)), "type `3` cannot be assigned to type `Std::String?`"),
			},
		},

		"break from a nested labeled loop": {
			input: `
				var b: 8 = $foo: for i in 2...10
					break[foo] 2.5 for j in [9, 2, 6]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(16, 2, 16), P(83, 4, 7)), "type `2.5?` cannot be assigned to type `8`"),
			},
		},

		"returns nil with continue": {
			input: `
				var b: 8 = (continue 2.5 for c in 2...10)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(17, 2, 17), P(44, 2, 44)), "type `nil` cannot be assigned to type `8`"),
			},
		},
		"continue a parent labeled loop": {
			input: `
				var b: 8 = $foo: for i in 1...20
					continue[foo] 2.5 for j in 7...30
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(16, 2, 16), P(83, 4, 7)), "type `nil` cannot be assigned to type `8`"),
			},
		},

		"unmatchable pattern": {
			input: `
				i for [1, i] in 1...20
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(11, 2, 11), P(16, 2, 16)), "type `Std::Int` cannot ever match type `Std::List[any]`"),
			},
		},
		"valid object pattern": {
			input: `
				key + value for Pair(key, value) in { foo: "bar", baz: "lol" }
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(5, 2, 5), P(15, 2, 15)), "method `+` is not defined on type `Std::Symbol`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestLoopExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				loop
					var b: Int? = a
				end
			`,
		},
		"returns never": {
			input: `
				a := 2
				b := loop
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(72, 7, 9), P(72, 7, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(68, 7, 5), P(72, 7, 9)), "unreachable code"),
			},
		},
		"returns nil when a naked break is present": {
			input: `
				var a: Int? = 2
				b := loop
					if a
						break
					end
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(82, 8, 9), P(82, 8, 9)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"returns the value given to break": {
			input: `
				var a: Int? = 2
				var b = loop
					if a
						break ""
					end
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(88, 8, 9), P(88, 8, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"returns the union of values given to break": {
			input: `
				var a: Int? = 2
				var b = loop
					if a
						break ""
					else
						break 2.5
					end
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(114, 10, 9), P(114, 10, 9)), "type `3` cannot be assigned to type `Std::String | Std::Float`"),
			},
		},
		"break nested labeled loop": {
			input: `
				var a: Int? = 2
				var b = $foo: loop
					loop
						break[foo] ""
					end
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(99, 8, 9), P(99, 8, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},

		"returns never when a naked continue is present": {
			input: `
				var a: Int? = 2
				b := loop
					if a
						continue
					end
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(85, 8, 9), P(85, 8, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(81, 8, 5), P(85, 8, 9)), "unreachable code"),
			},
		},
		"does not return the value given to continue": {
			input: `
				var a: Int? = 2
				var b = loop
					if a
						continue ""
					end
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(91, 8, 9), P(91, 8, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(87, 8, 5), P(91, 8, 9)), "unreachable code"),
			},
		},
		"continue in nested labeled loop": {
			input: `
				var a: Int? = 2
				var b = $foo: loop
					loop
						continue[foo] ""
					end
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(102, 8, 9), P(102, 8, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(98, 8, 5), P(102, 8, 9)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestWhileExpression(t *testing.T) {
	tests := testTable{
		"cannot use void in the condition": {
			input: `
				def foo; end
				while foo()
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(28, 3, 11), P(32, 3, 15)), "cannot use type `void` as a value in this context"),
			},
		},
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				while a
					var b: Int? = a
				end
			`,
		},
		"returns never if condition is truthy": {
			input: `
				a := 2
				b := while true
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(78, 7, 9), P(78, 7, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(74, 7, 5), P(78, 7, 9)), "unreachable code"),
			},
		},
		"returns nil if condition is falsy": {
			input: `
				a := 2
				var b: nil = while false
					"foo" + "bar"
					a + 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(35, 3, 24), P(39, 3, 28)), "this loop will never execute since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(46, 4, 6), P(59, 4, 19)), "unreachable code"),
			},
		},
		"returns a nilable body type if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(69, 5, 7)), "type `Std::String?` cannot be assigned to type `8`"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				while a
					var b: true = a
				end
			`,
		},

		"returns nil with a naked break if condition is truthy": {
			input: `
				a := 2
				b := while true
					c := false
					if c
						break
					end
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(106, 10, 9), P(106, 10, 9)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"returns the value given to break if condition is truthy": {
			input: `
				a := 2
				b := while true
					c := false
					if c
						break "foo" + "bar"
					end
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(120, 10, 9), P(120, 10, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"returns nil with a break if condition is falsy": {
			input: `
				a := 2
				var b: nil = while false
					c := false
					if c
						break "foo" + "bar"
					end
					a + 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(35, 3, 24), P(39, 3, 28)), "this loop will never execute since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(46, 4, 6), P(56, 4, 16)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					c := false
					if c
						break
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(116, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					c := false
					if c
						break 2.5
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(120, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"break from a nested labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: while a
					var b: Int? = 9
					while b
						break[foo] 2.5
						"foo" + "bar"
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(111, 7, 7), P(123, 7, 19)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(140, 9, 7)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},

		"returns never with a naked continue if condition is truthy": {
			input: `
				a := 2
				b := while true
					continue
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewWarning(L("<main>", P(51, 5, 6), P(55, 5, 10)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(73, 7, 9), P(73, 7, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(69, 7, 5), P(73, 7, 9)), "unreachable code"),
			},
		},
		"returns never with continue if condition is truthy": {
			input: `
				a := 2
				b := while true
					continue "foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewWarning(L("<main>", P(65, 5, 6), P(69, 5, 10)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(87, 7, 9), P(87, 7, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(83, 7, 5), P(87, 7, 9)), "unreachable code"),
			},
		},
		"returns nil with a continue if condition is falsy": {
			input: `
				a := 2
				var b: nil = while false
					continue "foo" + "bar"
					a + 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(35, 3, 24), P(39, 3, 28)), "this loop will never execute since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(46, 4, 6), P(68, 4, 28)), "unreachable code"),
				diagnostic.NewWarning(L("<main>", P(74, 5, 6), P(78, 5, 10)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					c := false
					if c
						continue
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(119, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					c := false
					if c
						continue 2.5
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(123, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"continue a parent labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: while a
					var b: Int? = 9
					while b
						continue[foo] 2.5
						"foo" + "bar"
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(114, 7, 7), P(126, 7, 19)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(143, 9, 7)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestWhileModifier(t *testing.T) {
	tests := testTable{
		"cannot use void in the condition": {
			input: `
				def foo; end
				nil while foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(32, 3, 15), P(36, 3, 19)), "cannot use type `void` as a value in this context"),
			},
		},
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				var b: Int? = (a while a)
			`,
		},
		"returns never if condition is truthy": {
			input: `
				a := 2
				b := (a + 2 while true)
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(34, 3, 23), P(37, 3, 26)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(48, 4, 9), P(48, 4, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(44, 4, 5), P(48, 4, 9)), "unreachable code"),
			},
		},
		"returns a nilable body type if condition is falsy": {
			input: `
				a := 2
				var b: nil = (a + 2 while false)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(42, 3, 31), P(46, 3, 35)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewFailure(L("<main>", P(30, 3, 19), P(46, 3, 35)), "type `Std::Int?` cannot be assigned to type `nil`"),
			},
		},
		"returns a nilable body type if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = ("foo" + "bar" while a)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(37, 3, 17), P(57, 3, 37)), "type `Std::String?` cannot be assigned to type `8`"),
			},
		},

		"do not narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				(var b: true = a while a)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(38, 3, 20), P(38, 3, 20)), "type `bool` cannot be assigned to type `true`"),
			},
		},

		"returns the value given to break if condition is truthy": {
			input: `
				a := 2
				b := (break "foo" + "bar" while true)
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(48, 3, 37), P(51, 3, 40)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(62, 4, 9), P(62, 4, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"returns nil with a break if condition is falsy": {
			input: `
				a := 2
				var b: nil = (break "foo" + "bar" while false)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(56, 3, 45), P(60, 3, 49)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewFailure(L("<main>", P(30, 3, 19), P(60, 3, 49)), "type `nil | Std::String` cannot be assigned to type `nil`"),
			},
		},
		"break from a nested labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: do
					var b: Int? = 9
					break[foo] 2.5 while b
				end while a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(108, 6, 15)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},

		"returns never with continue if condition is truthy": {
			input: `
				a := 2
				b := (continue "foo" + "bar" while true)
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(51, 3, 40), P(54, 3, 43)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(65, 4, 9), P(65, 4, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(61, 4, 5), P(65, 4, 9)), "unreachable code"),
			},
		},
		"returns a nilable body type with a continue if condition is falsy": {
			input: `
				a := 2
				var b: nil = (continue "foo" + "bar" while false)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(59, 3, 48), P(63, 3, 52)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewFailure(L("<main>", P(30, 3, 19), P(63, 3, 52)), "type `nil | Std::String` cannot be assigned to type `nil`"),
			},
		},
		"continue a parent labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: while a
					var b: Int? = 9
					continue[foo] 2.5 while b
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(108, 6, 7)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestUntilExpression(t *testing.T) {
	tests := testTable{
		"cannot use void in the condition": {
			input: `
				def foo; end
				until foo()
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(28, 3, 11), P(32, 3, 15)), "cannot use type `void` as a value in this context"),
			},
		},
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				until a
					var b: Int? = a
				end
			`,
		},
		"returns never if condition is falsy": {
			input: `
				a := 2
				b := until false
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(27, 3, 16), P(31, 3, 20)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewFailure(L("<main>", P(79, 7, 9), P(79, 7, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(75, 7, 5), P(79, 7, 9)), "unreachable code"),
			},
		},
		"returns nil if condition is truthy": {
			input: `
				a := 2
				var b: nil = until true
					"foo" + "bar"
					a + 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(35, 3, 24), P(38, 3, 27)), "this loop will never execute since type `true` is truthy"),
				diagnostic.NewWarning(L("<main>", P(45, 4, 6), P(58, 4, 19)), "unreachable code"),
			},
		},
		"returns a nilable body type if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = until a
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(69, 5, 7)), "type `Std::String?` cannot be assigned to type `8`"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				until a
					var b: false = a
				end
			`,
		},

		"returns nil with a naked break if condition is falsy": {
			input: `
				a := 2
				b := until false
					break
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(27, 3, 16), P(31, 3, 20)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(49, 5, 6), P(53, 5, 10)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(71, 7, 9), P(71, 7, 9)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"returns the value given to break if condition is falsy": {
			input: `
				a := 2
				b := until false
					break "foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(27, 3, 16), P(31, 3, 20)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(63, 5, 6), P(67, 5, 10)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(85, 7, 9), P(85, 7, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"returns nil with a break if condition is truthy": {
			input: `
				a := 2
				var b: nil = until true
					break "foo" + "bar"
					a + 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(35, 3, 24), P(38, 3, 27)), "this loop will never execute since type `true` is truthy"),
				diagnostic.NewWarning(L("<main>", P(45, 4, 6), P(64, 4, 25)), "unreachable code"),
				diagnostic.NewWarning(L("<main>", P(70, 5, 6), P(74, 5, 10)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = until a
					var c: Int? = 3
					if c
						break
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(121, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = until a
					var c: Int? = 3
					if c
						break 2.5
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(125, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"break from a nested labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: until a
					var b: Int? = 9
					until b
						var c: Int? = 3
						if c
							break[foo] 2.5
						end
						"foo" + "bar"
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(184, 12, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},

		"returns never with a naked continue if condition is falsy": {
			input: `
				a := 2
				b := until false
					continue
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(27, 3, 16), P(31, 3, 20)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(52, 5, 6), P(56, 5, 10)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(74, 7, 9), P(74, 7, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(70, 7, 5), P(74, 7, 9)), "unreachable code"),
			},
		},
		"returns never with continue if condition is falsy": {
			input: `
				a := 2
				b := until false
					continue "foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(27, 3, 16), P(31, 3, 20)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(66, 5, 6), P(70, 5, 10)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(88, 7, 9), P(88, 7, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(84, 7, 5), P(88, 7, 9)), "unreachable code"),
			},
		},
		"returns nil with a continue if condition is truthy": {
			input: `
				a := 2
				var b: nil = until true
					continue "foo" + "bar"
					a + 2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(35, 3, 24), P(38, 3, 27)), "this loop will never execute since type `true` is truthy"),
				diagnostic.NewWarning(L("<main>", P(45, 4, 6), P(67, 4, 28)), "unreachable code"),
				diagnostic.NewWarning(L("<main>", P(73, 5, 6), P(77, 5, 10)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = until a
					var c: Int? = 3
					if c
						continue
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(124, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = until a
					var c: Int? = 3
					if c
						continue 2.5
					end
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(128, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"continue a parent labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: until a
					var b: Int? = 9
					until b
						var c: Int? = 3
						if c
							continue[foo] 2.5
						end
						"foo" + "bar"
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(187, 12, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestUntilModifier(t *testing.T) {
	tests := testTable{
		"cannot use void in the condition": {
			input: `
				def foo; end
				nil until foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(32, 3, 15), P(36, 3, 19)), "cannot use type `void` as a value in this context"),
			},
		},
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				(var b: Int? = a) until a
			`,
		},
		"returns never if condition is falsy": {
			input: `
				a := 2
				b := (a + 2 until false)
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(34, 3, 23), P(38, 3, 27)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewFailure(L("<main>", P(49, 4, 9), P(49, 4, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(45, 4, 5), P(49, 4, 9)), "unreachable code"),
			},
		},
		"returns nil if condition is truthy": {
			input: `
				a := 2
				var b: nil = (a + 2 until true)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(42, 3, 31), P(45, 3, 34)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(30, 3, 19), P(45, 3, 34)), "type `Std::Int?` cannot be assigned to type `nil`"),
			},
		},
		"returns a nilable body type if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = ("foo" + "bar" until a)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(37, 3, 17), P(57, 3, 37)), "type `Std::String?` cannot be assigned to type `8`"),
			},
		},

		"do not narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				(var b: false = a) until a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(39, 3, 21), P(39, 3, 21)), "type `bool` cannot be assigned to type `false`"),
			},
		},

		"returns the value given to break if condition is falsy": {
			input: `
				a := 2
				b := (break "foo" + "bar" until false)
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(48, 3, 37), P(52, 3, 41)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewFailure(L("<main>", P(63, 4, 9), P(63, 4, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"returns a nilable body type with a break if condition is truthy": {
			input: `
				a := 2
				var b: nil = (break "foo" + "bar" until true)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(56, 3, 45), P(59, 3, 48)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(30, 3, 19), P(59, 3, 48)), "type `nil | Std::String` cannot be assigned to type `nil`"),
			},
		},
		"break from a nested labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: do
					var b: Int? = 9
					break[foo] 2.5 until b
				end until a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(108, 6, 15)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},

		"returns never with continue if condition is falsy": {
			input: `
				a := 2
				b := (continue "foo" + "bar" until false)
				b = 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(51, 3, 40), P(55, 3, 44)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewFailure(L("<main>", P(66, 4, 9), P(66, 4, 9)), "type `3` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(62, 4, 5), P(66, 4, 9)), "unreachable code"),
			},
		},
		"returns nil with a continue if condition is truthy": {
			input: `
				a := 2
				var b: nil = (continue "foo" + "bar" until true)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(59, 3, 48), P(62, 3, 51)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewFailure(L("<main>", P(30, 3, 19), P(62, 3, 51)), "type `nil | Std::String` cannot be assigned to type `nil`"),
			},
		},
		"continue a parent labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: until a
					var b: Int? = 9
					continue[foo] 2.5 until b
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(108, 6, 7)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestUnlessExpression(t *testing.T) {
	tests := testTable{
		"cannot use void in the condition": {
			input: `
				def foo; end
				unless foo()
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(29, 3, 12), P(33, 3, 16)), "cannot use type `void` as a value in this context"),
			},
		},
		"checks modifier version": {
			input: `
				def foo; end
				nil unless foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(33, 3, 16), P(37, 3, 20)), "cannot use type `void` as a value in this context"),
			},
		},
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				unless a
					var b: Int? = a
				else
					a
				end
			`,
		},
		"returns the last else expression if condition is truthy": {
			input: `
				a := 2
				var b: Float = unless true
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(38, 3, 27), P(41, 3, 30)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewWarning(L("<main>", P(48, 4, 6), P(61, 4, 19)), "unreachable code"),
			},
		},
		"returns the last then expression if condition is falsy": {
			input: `
				a := 2
				var b: Int = unless false
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(36, 3, 25), P(40, 3, 29)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(86, 7, 6), P(89, 7, 9)), "unreachable code"),
			},
		},
		"returns a union of both branches if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = unless a
					"foo" + "bar"
				else
					2.2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(88, 7, 7)), "type `Std::String | 2.2` cannot be assigned to type `8`"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				unless a
					var b: false = a
				else
					var b: true = a
				end
			`,
		},
		"narrow nilable variable and return": {
			input: `
				def foo(a: Int?)
					return unless a

					var b: nil = a
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(62, 5, 19), P(62, 5, 19)), "type `Std::Int` cannot be assigned to type `nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestIfExpression(t *testing.T) {
	tests := testTable{
		"cannot use void in the condition": {
			input: `
				def foo; end
				if foo()
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(25, 3, 8), P(29, 3, 12)), "cannot use type `void` as a value in this context"),
			},
		},
		"checks modifier version": {
			input: `
				def foo; end
				nil if foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(29, 3, 12), P(33, 3, 16)), "cannot use type `void` as a value in this context"),
			},
		},
		"checks modifier version with else": {
			input: `
				a := 2
				var b: Int = (a + 2 if true else 2.2)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(39, 3, 28), P(42, 3, 31)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewWarning(L("<main>", P(49, 3, 38), P(51, 3, 40)), "unreachable code"),
			},
		},
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				if a
					var b: Int? = a
				else
					a
				end
			`,
		},
		"returns the last then expression if condition is truthy": {
			input: `
				a := 2
				var b: Int = if true
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(32, 3, 21), P(35, 3, 24)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewWarning(L("<main>", P(81, 7, 6), P(84, 7, 9)), "unreachable code"),
			},
		},
		"returns the last else expression if condition is truthy": {
			input: `
				a := 2
				var b: Float = if false
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(34, 3, 23), P(38, 3, 27)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewWarning(L("<main>", P(45, 4, 6), P(58, 4, 19)), "unreachable code"),
			},
		},
		"returns a union of both branches if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = if a
					"foo" + "bar"
				else
					2.2
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(84, 7, 7)), "type `Std::String | 2.2` cannot be assigned to type `8`"),
			},
		},
		"returns a union of all branches if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: Float? = nil
				var c: 8 = if a
					"foo" + "bar"
				else if b
					2.5
				else
					%/foo/
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(60, 4, 16), P(134, 10, 7)), "type `Std::String | 2.5 | Std::Regex` cannot be assigned to type `8`"),
			},
		},
		"returns a union of then and nil if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = if a
					"foo" + "bar"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 3, 16), P(66, 5, 7)), "type `Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"returns nil when empty": {
			input: `
				var a: Int? = nil
				var b: nil = if a; end
			`,
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				if a
					var b: true = a
				else
					var b: false = a
				end
			`,
		},
		"narrow nilable variable type by using truthiness": {
			input: `
				var a: Int? = nil
				if a
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow nilable variable declaration by using truthiness": {
			input: `
				var a: Int? = nil
				if var b: Int | Float | nil = a
					var c: Int = b
				else
					var c: nil = b
				end
			`,
		},
		"narrow nilable variable assignment by using truthiness": {
			input: `
				var a: Int? = nil
				var b: Int | Float | nil
				if b = a
					var c: Int = b
				else
					var c: nil = b
				end
			`,
		},
		"narrow nilable variable by using truthiness and assign wider type": {
			input: `
				var a: Int? = nil
				if a
					var c: Int = a
				else
					var c: nil = a
					a = 8
					var d: Int = a
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(110, 8, 19), P(110, 8, 19)), "type `Std::Int?` cannot be assigned to type `Std::Int`"),
			},
		},
		"narrow nilable variable and return": {
			input: `
				def foo(a: Int?)
					return if !a

					var b: nil = a
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(59, 5, 19), P(59, 5, 19)), "type `Std::Int` cannot be assigned to type `nil`"),
			},
		},
		"widen narrowed type if a wider value is assigned": {
			input: `
				var a: Int | Float | nil = nil
				switch a
				case Int() || nil
					if a
						var c: Int = a
					else
						var c: nil = a
						a = 2.5
						var d: 2 = a
					end
					var e: 3 = a
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(164, 10, 18), P(164, 10, 18)), "type `Std::Int | Std::Float | nil` cannot be assigned to type `2`"),
				diagnostic.NewFailure(L("<main>", P(191, 12, 17), P(191, 12, 17)), "type `Std::Int | Std::Float | nil` cannot be assigned to type `3`"),
			},
		},
		"narrow named nilable variable type by using truthiness": {
			input: `
				typedef Foo = Int?
				var a: Foo = nil
				if a
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow nilable variable type by using negated truthiness": {
			input: `
				var a: Int? = nil
				if !a
					var b: nil = a
				else
					var b: Int = a
				end
			`,
		},
		"narrow named nilable variable type by using negated truthiness": {
			input: `
				typedef Foo = Int?
				var a: Foo = nil
				if !a
					var b: nil = a
				else
					var b: Int = a
				end
			`,
		},
		"narrow union type by using <:": {
			input: `
				var a: Int | String = "foo"
				if a <: Int
					var b: Int = a
				else
					var b: String = a
				end
			`,
		},
		"narrow variable type by using <:": {
			input: `
				var a: Int? = nil
				if a <: Int
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow variable type by using <<:": {
			input: `
				var a: Int? = nil
				if a <<: Int
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow a few variables with &&": {
			input: `
				var a: Int? = nil
				var b = false
				if a && b
					var c: Int = a
					var d: true = b
				end
			`,
		},
		"narrow with an impossible && branch": {
			input: `
				var a: Int? = 3
				if a && nil
					a = :foo
				else
					a = :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(46, 4, 10), P(49, 4, 13)), "type `:foo` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(69, 6, 10), P(72, 6, 13)), "type `:bar` cannot be assigned to type `Std::Int?`"),
				diagnostic.NewWarning(L("<main>", P(28, 3, 8), P(35, 3, 15)), "this condition will always have the same result since type `nil` is falsy"),
				diagnostic.NewWarning(L("<main>", P(42, 4, 6), P(50, 4, 14)), "unreachable code"),
			},
		},
		"narrow a few variables with ||": {
			input: `
				var a: Int? = nil
				var b = false
				if a || b
				else
					var c: nil = a
					var d: false = b
				end
			`,
		},
		"narrow with an impossible || branch": {
			input: `
				var a: Int? = nil
				if a || 3
					a = :foo
				else
					a = :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(46, 4, 10), P(49, 4, 13)), "type `:foo` cannot be assigned to type `Std::Int?`"),
				diagnostic.NewFailure(L("<main>", P(69, 6, 10), P(72, 6, 13)), "type `:bar` cannot be assigned to type `never`"),
				diagnostic.NewWarning(L("<main>", P(30, 3, 8), P(35, 3, 13)), "this condition will always have the same result since type `Std::Int` is truthy"),
				diagnostic.NewWarning(L("<main>", P(65, 6, 6), P(73, 6, 14)), "unreachable code"),
			},
		},

		"narrow with ===": {
			input: `
				var a: Int | Float = 1
				var b: Float? = .2
				if a === b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(75, 5, 10), P(78, 5, 13)), "type `:foo` cannot be assigned to type `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(89, 6, 10), P(92, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(112, 8, 10), P(115, 8, 13)), "type `:baz` cannot be assigned to type `Std::Int | Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(126, 9, 10), P(130, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float?`"),
			},
		},
		"narrow with an impossible ===": {
			input: `
				var a: Int = 1
				var b: Float = .2
				if a === b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(49, 4, 8), P(49, 4, 8)), "this strict equality check is impossible, `Std::Int` cannot ever be equal to `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(66, 5, 10), P(69, 5, 13)), "type `:foo` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(80, 6, 10), P(83, 6, 13)), "type `:bar` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(103, 8, 10), P(106, 8, 13)), "type `:baz` cannot be assigned to type `Std::Int`"),
				diagnostic.NewFailure(L("<main>", P(117, 9, 10), P(121, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float`"),
			},
		},
		"narrow with !==": {
			input: `
				var a: Int | Float = 1
				var b: Float? = .2
				if a !== b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(75, 5, 10), P(78, 5, 13)), "type `:foo` cannot be assigned to type `Std::Int | Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(89, 6, 10), P(92, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float?`"),
				diagnostic.NewFailure(L("<main>", P(112, 8, 10), P(115, 8, 13)), "type `:baz` cannot be assigned to type `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(126, 9, 10), P(130, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float`"),
			},
		},
		"narrow with an impossible !==": {
			input: `
				var a: Int = 1
				var b: Float = .2
				if a !== b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(49, 4, 8), P(49, 4, 8)), "this strict equality check is impossible, `Std::Int` cannot ever be equal to `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(66, 5, 10), P(69, 5, 13)), "type `:foo` cannot be assigned to type `Std::Int`"),
				diagnostic.NewFailure(L("<main>", P(80, 6, 10), P(83, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(103, 8, 10), P(106, 8, 13)), "type `:baz` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(117, 9, 10), P(121, 9, 14)), "type `:fizz` cannot be assigned to type `never`"),
			},
		},

		"narrow with ==": {
			input: `
				var a: Int | Float = 1
				var b: Float? = .2
				if a == b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(74, 5, 10), P(77, 5, 13)), "type `:foo` cannot be assigned to type `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(88, 6, 10), P(91, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(111, 8, 10), P(114, 8, 13)), "type `:baz` cannot be assigned to type `Std::Int | Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(125, 9, 10), P(129, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float?`"),
			},
		},
		"narrow with an impossible ==": {
			input: `
				var a: Int = 1
				var b: Float = .2
				if a == b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(49, 4, 8), P(49, 4, 8)), "this equality check is impossible, `Std::Int` cannot ever be equal to `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(65, 5, 10), P(68, 5, 13)), "type `:foo` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(79, 6, 10), P(82, 6, 13)), "type `:bar` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(102, 8, 10), P(105, 8, 13)), "type `:baz` cannot be assigned to type `Std::Int`"),
				diagnostic.NewFailure(L("<main>", P(116, 9, 10), P(120, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float`"),
			},
		},
		"narrow with !=": {
			input: `
				var a: Int | Float = 1
				var b: Float? = .2
				if a != b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(74, 5, 10), P(77, 5, 13)), "type `:foo` cannot be assigned to type `Std::Int | Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(88, 6, 10), P(91, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float?`"),
				diagnostic.NewFailure(L("<main>", P(111, 8, 10), P(114, 8, 13)), "type `:baz` cannot be assigned to type `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(125, 9, 10), P(129, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float`"),
			},
		},
		"narrow with an impossible !=": {
			input: `
				var a: Int = 1
				var b: Float = .2
				if a != b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(49, 4, 8), P(49, 4, 8)), "this equality check is impossible, `Std::Int` cannot ever be equal to `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(65, 5, 10), P(68, 5, 13)), "type `:foo` cannot be assigned to type `Std::Int`"),
				diagnostic.NewFailure(L("<main>", P(79, 6, 10), P(82, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float`"),
				diagnostic.NewFailure(L("<main>", P(102, 8, 10), P(105, 8, 13)), "type `:baz` cannot be assigned to type `never`"),
				diagnostic.NewFailure(L("<main>", P(116, 9, 10), P(120, 9, 14)), "type `:fizz` cannot be assigned to type `never`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestLogicalAnd(t *testing.T) {
	tests := testTable{
		"cannot use void on the left hand side": {
			input: `
				def foo; end
				foo() && foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 5), P(26, 3, 9)), "cannot use type `void` as a value in this context"),
			},
		},
		"returns the right type when the left type is truthy": {
			input: `
				var a = "foo"
				var b = 2
				var c: Int = a && b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(50, 4, 18), P(50, 4, 18)), "this condition will always have the same result since type `Std::String` is truthy"),
			},
		},
		"returns the left type when the left type is falsy": {
			input: `
				var a = nil
				var b = 2
				var c: nil = a && b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result since type `nil` is falsy"),
				diagnostic.NewWarning(L("<main>", P(53, 4, 23), P(53, 4, 23)), "unreachable code"),
			},
		},
		"returns a union of both types with only nil when the left can be both truthy and falsy": {
			input: `
				var a: String? = "foo"
				var b = 2
				var c: 9 = a && b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(57, 4, 16), P(62, 4, 21)), "type `nil | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types with only false when the left can be both truthy and falsy": {
			input: `
				var a = false
				var b = 2
				var c: 9 = a && b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(48, 4, 16), P(53, 4, 21)), "type `false | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without duplication": {
			input: `
				var a: false | nil | Int = nil
				var b: Float | Int | nil = 2.2
				var c: 9 = a && b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(86, 4, 16), P(91, 4, 21)), "type `false | nil | Std::Float | Std::Int` cannot be assigned to type `9`"),
			},
		},

		"narrow left variable to non falsy": {
			input: `
				var a: false | nil | Int = nil
				a && a + 2
			`,
		},
		"narrow a few variables to non falsy": {
			input: `
				var a: Int? = nil
				var b: Int? = nil
				a && b && a + b
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestLogicalOr(t *testing.T) {
	tests := testTable{
		"cannot use void on the left hand side": {
			input: `
				def foo; end
				foo() || foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 5), P(26, 3, 9)), "cannot use type `void` as a value in this context"),
			},
		},
		"returns the left type when it is truthy": {
			input: `
				var a = "foo"
				var b = 2
				var c: String = a || b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(53, 4, 21), P(53, 4, 21)), "this condition will always have the same result since type `Std::String` is truthy"),
				diagnostic.NewWarning(L("<main>", P(58, 4, 26), P(58, 4, 26)), "unreachable code"),
			},
		},
		"returns the right type when the left type is falsy": {
			input: `
				var a = nil
				var b = 2
				var c: Int = a || b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result since type `nil` is falsy"),
			},
		},
		"returns a union of both types without nil when the left can be both truthy and falsy": {
			input: `
				var a: String? = "foo"
				var b = 2
				var c: 9 = a || b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(57, 4, 16), P(62, 4, 21)), "type `Std::String | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without false when the left can be both truthy and falsy": {
			input: `
				var a = false
				var b = 2
				var c: 9 = a || b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(48, 4, 16), P(53, 4, 21)), "type `true | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without duplication": {
			input: `
				var a: String | Int | nil = nil
				var b: Float | Int | nil = 2.2
				var c: 9 = a || b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(87, 4, 16), P(92, 4, 21)), "type `Std::String | Std::Int | Std::Float | nil` cannot be assigned to type `9`"),
			},
		},

		"narrow left variable to falsy": {
			input: `
				var a: false | nil | Int = nil
				a || var b: 9 = a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(56, 3, 21), P(56, 3, 21)), "type `false | nil` cannot be assigned to type `9`"),
			},
		},
		"narrow a few variables to non falsy": {
			input: `
				var a: Int? = nil
				var b: Int? = nil
				a || b || var c: 9 = a && b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(70, 4, 26), P(70, 4, 26)), "this condition will always have the same result since type `nil` is falsy"),
				diagnostic.NewWarning(L("<main>", P(75, 4, 31), P(75, 4, 31)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(70, 4, 26), P(75, 4, 31)), "type `nil` cannot be assigned to type `9`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestNilCoalescing(t *testing.T) {
	tests := testTable{
		"cannot use void on the left hand side": {
			input: `
				def foo; end
				foo() ?? foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 5), P(26, 3, 9)), "cannot use type `void` as a value in this context"),
			},
		},
		"returns the left type when it is not nilable": {
			input: `
				var a = "foo"
				var b = 2
				var c: String = a ?? b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(53, 4, 21), P(53, 4, 21)), "this condition will always have the same result since type `Std::String` can never be nil"),
				diagnostic.NewWarning(L("<main>", P(58, 4, 26), P(58, 4, 26)), "unreachable code"),
			},
		},
		"returns the right type when the left type is nil": {
			input: `
				var a = nil
				var b = 2
				var c: Int = a ?? b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result"),
			},
		},
		"returns a union of both types without nil when the left can be both nil and not nil": {
			input: `
				var a: String? = "foo"
				var b = 2
				var c: 9 = a ?? b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(57, 4, 16), P(62, 4, 21)), "type `Std::String | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without duplication": {
			input: `
				var a: String | Int | nil = nil
				var b: Float | Int | nil = 2.2
				var c: 9 = a ?? b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(87, 4, 16), P(92, 4, 21)), "type `Std::String | Std::Int | Std::Float | nil` cannot be assigned to type `9`"),
			},
		},

		"narrow left variable to non nilable": {
			input: `
				var a: false | nil | Int = nil
				a ?? var b: 9 = a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(56, 3, 21), P(56, 3, 21)), "type `nil` cannot be assigned to type `9`"),
			},
		},
		"narrow a few variables to non nilable": {
			input: `
				var a: Int? = nil
				var b: Int? = nil
				a ?? b ?? var c: 9 = a && b
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L("<main>", P(70, 4, 26), P(70, 4, 26)), "this condition will always have the same result since type `nil` is falsy"),
				diagnostic.NewWarning(L("<main>", P(75, 4, 31), P(75, 4, 31)), "unreachable code"),
				diagnostic.NewFailure(L("<main>", P(70, 4, 26), P(75, 4, 31)), "type `nil` cannot be assigned to type `9`"),
			},
		},
		"narrow nested ||": {
			input: `
				var a: bool? = false
				var b: bool? = false
				(a || b) ?? do
					a = :foo
					b = :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(79, 5, 10), P(82, 5, 13)), "type `:foo` cannot be assigned to type `nil`"),
				diagnostic.NewFailure(L("<main>", P(93, 6, 10), P(96, 6, 13)), "type `:bar` cannot be assigned to type `nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
