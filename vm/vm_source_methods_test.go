package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestVMSource_Subscript(t *testing.T) {
	tests := sourceTestTable{
		"get index 0 of a list": {
			source: `
				list := ["foo", 2, 7.8]
				list[0]
			`,
			wantStackTop: value.String("foo"),
		},
		"get index -1 of a list": {
			source: `
				list := ["foo", 2, 7.8]
				list[-1]
			`,
			wantStackTop: value.Float(7.8),
		},
		"get too big index": {
			source: `
				list := ["foo", 2, 7.8]
				list[50]
			`,
			wantRuntimeErr: value.NewError(
				value.IndexErrorClass,
				"index 50 out of range: -3...3",
			),
		},
		"get too small index": {
			source: `
				list := ["foo", 2, 7.8]
				list[-10]
			`,
			wantRuntimeErr: value.NewError(
				value.IndexErrorClass,
				"index -10 out of range: -3...3",
			),
		},
		"get from nil": {
			source: `
				list := nil
				list[-10]
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(21, 3, 5), P(29, 3, 13)), "method `[]` is not defined on type `Std::Nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_NilSafeSubscript(t *testing.T) {
	tests := sourceTestTable{
		"get index 0 of a list": {
			source: `
				var list: List[String | Int | Float]? = ["foo", 2, 7.8]
				list?[0]
			`,
			wantStackTop: value.String("foo"),
		},
		"get index -1 of a list": {
			source: `
				var list: List[String | Int | Float]? = ["foo", 2, 7.8]
				list?[-1]
			`,
			wantStackTop: value.Float(7.8),
		},
		"get too big index": {
			source: `
				var list: List[String | Int | Float]? = ["foo", 2, 7.8]
				list?[50]
			`,
			wantRuntimeErr: value.NewError(
				value.IndexErrorClass,
				"index 50 out of range: -3...3",
			),
		},
		"get too small index": {
			source: `
				var list: List[String | Int | Float]? = ["foo", 2, 7.8]
				list?[-10]
			`,
			wantRuntimeErr: value.NewError(
				value.IndexErrorClass,
				"index -10 out of range: -3...3",
			),
		},
		"get from nil": {
			source: `
				var list: List[Int]? =  nil
				list?[-10]
			`,
			wantStackTop: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Instantiate(t *testing.T) {
	tests := sourceTestTable{
		"instantiate a class without an initialiser without arguments": {
			source: `
				class Foo; end

				::Foo()
			`,
			wantStackTop: value.NewObject(
				value.ObjectWithClass(
					value.NewClassWithOptions(
						value.ClassWithName("Foo"),
					),
				),
			),
		},
		"instantiate a class without an initialiser with arguments": {
			source: `
				class Foo; end

				::Foo(2)
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(25, 4, 5), P(32, 4, 12)), "expected 0 arguments in call to `#init`, got 1"),
			},
		},
		"instantiate a class with an initialiser without arguments": {
			source: `
				class Foo
					init(a: String); end
				end

				::Foo()
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(54, 6, 5), P(60, 6, 11)), "argument `a` is missing in call to `#init`"),
			},
		},
		"instantiate a class with an initialiser with arguments": {
			source: `
				class Foo
					init(a: String)
						println("a: " + a)
					end
				end

				::Foo("bar")
				nil
			`,
			wantStdout:   "a: bar\n",
			wantStackTop: value.Nil,
		},
		"instantiate a class with an initialiser with ivar parameters": {
			source: `
				class Foo
					init(@a: String)
						println("a: " + a)
					end

					getter a: String
				end

				f := Foo("bar")
				f.a
			`,
			wantStdout:   "a: bar\n",
			wantStackTop: value.String("bar"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Alias(t *testing.T) {
	tests := sourceTestTable{
		"add an alias to a builtin method in Std::Int": {
			source: `
				sealed primitive class ::Std::Int < Value
					alias add +
				end

				3.add(4)
			`,
			wantStackTop: value.SmallInt(7),
		},
		"add an alias to a nonexistent method": {
			source: `
				alias foo blabla
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(11, 2, 11), P(20, 2, 20)), "method `blabla` is not defined on type `Std::Kernel`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_DefineMethod(t *testing.T) {
	tests := sourceTestTable{
		"define a method in top level": {
			source: `
				def foo: Symbol then :bar
				foo()
			`,
			wantStackTop: value.ToSymbol("bar"),
		},
		"define a method with positional arguments in top level": {
			source: `
				def foo(a: Int, b: Int): Int
					c := 5
					a + b + c
				end
				foo(1, 2)
			`,
			wantStackTop: value.SmallInt(8),
		},
		"define a method with positional arguments in a class": {
			source: `
				class Bar
					def foo(a: Int, b: Int): Int
						c := 5
						a + b + c
					end
				end

				Bar().foo(1, 2)
			`,
			wantStackTop: value.SmallInt(8),
		},
		"define a method with positional arguments in a module": {
			source: `
				module Bar
					def foo(a: Int, b: Int): Int
						c := 5
						a + b + c
					end
				end

				Bar.foo(1, 2)
			`,
			wantStackTop: value.SmallInt(8),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_CallClosure(t *testing.T) {
	tests := sourceTestTable{
		"call": {
			source: `
				pow2 := |a: Int| -> a ** 2
				pow2.call(5)
			`,
			wantStackTop: value.SmallInt(25),
		},
		"invalid args": {
			source: `
				pow2 := |a: Int| -> a ** 2
				pow2.call(5, 8)
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(36, 3, 5), P(50, 3, 19)), "expected 1 arguments in call to `call`, got 2"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_CallMethod(t *testing.T) {
	tests := sourceTestTable{
		"nil safe call on nil": {
			source: `
				var a: Int? = nil

				try a?.to_string?.to_int
			`,
			wantStackTop: value.Nil,
		},
		"nil safe call on not nil": {
			source: `
				var a: Int? = 5

				a?.inspect
			`,
			wantStackTop: value.String("5"),
		},
		"call a global method without arguments": {
			source: `
				def foo: Symbol
					:bar
				end

				foo()
			`,
			wantStackTop: value.ToSymbol("bar"),
		},
		"call a global method with positional arguments": {
			source: `
				def add(a: Int, b: Int): Int
					a + b
				end

				add(5, 9)
			`,
			wantStackTop: value.SmallInt(14),
		},
		"call a method with missing required arguments": {
			source: `
				def add(a: Int, b: Int): Int
					a + b
				end

				add(5)
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(58, 6, 5), P(63, 6, 10)), "argument `b` is missing in call to `add`"),
			},
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"`add` wrong number of arguments, given: 1, expected: 2..2",
			),
		},
		"call a method without optional arguments": {
			source: `
				def add(a: Int, b: Int = 3, c: Float = 20.5): Float
					a + b + c
				end

				add(5)
			`,
			wantStackTop: value.Float(28.5),
		},
		"call a method with some optional arguments": {
			source: `
				def add(a: Int, b: Int = 3, c: Float = 20.5): Float
					a + b + c
				end

				add(5, 0)
			`,
			wantStackTop: value.Float(25.5),
		},
		"call a method with all optional arguments": {
			source: `
				def add(a: Int, b: Int = 3, c: Float = 20.5): Float
					a + b + c
				end

				add(3, 2, 3.5)
			`,
			wantStackTop: value.Float(8.5),
		},
		"call a method with only named arguments": {
			source: `
				def foo(a: String, b: String, c: String, d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				foo(b: "b", a: "a", c: "c", e: "e")
			`,
			wantStackTop: value.String("a: a, b: b, c: c, d: default d, e: e"),
		},
		"call a method with all required arguments and named arguments": {
			source: `
				def foo(a: String, b: String, c: String, d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				foo("a", c: "c", b: "b")
			`,
			wantStackTop: value.String("a: a, b: b, c: c, d: default d, e: default e"),
		},
		"call a method with optional arguments and named arguments": {
			source: `
				def foo(a: String, b: String, c: String = "default c", d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				foo("a", "b", "c", e: "e")
			`,
			wantStackTop: value.String("a: a, b: b, c: c, d: default d, e: e"),
		},
		"call a method with a named rest param and no args": {
			source: `
				def foo(**a: String): String
					"a: #{a}"
				end

				foo()
			`,
			wantStackTop: value.String("a: %{}"),
		},
		"call a method with a named rest param and a few named args": {
			source: `
				def foo(**a: String): Record[Symbol, String]
					a
				end

				foo(d: "foo", a: "bar")
			`,
			wantStackTop: vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{Key: value.ToSymbol("a"), Value: value.String("bar")},
				value.Pair{Key: value.ToSymbol("d"), Value: value.String("foo")},
			),
		},
		"call a method with regular params, named rest param and a few named args": {
			source: `
				def foo(a: String, **b: String): Tuple[any]
					[a, b]
				end

				foo("foo", c: "bar", d: "baz")
			`,
			wantStackTop: value.NewArrayListWithElements(
				2,
				value.String("foo"),
				vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("c"), Value: value.String("bar")},
					value.Pair{Key: value.ToSymbol("d"), Value: value.String("baz")},
				),
			),
		},
		"call a method with regular params, named rest param and only required args": {
			source: `
				def foo(a: String, **b: String): Tuple[any]
					[a, b]
				end

				foo("foo")
			`,
			wantStackTop: value.NewArrayListWithElements(
				2,
				value.String("foo"),
				value.NewHashRecord(0),
			),
		},
		"call a method with regular params, optional params, named rest param and a few named args": {
			source: `
				def foo(a: String, b: Int = 5, **c: String): Tuple[any]
					[a, b, c]
				end

				foo("foo", c: "bar", d: "baz")
			`,
			wantStackTop: value.NewArrayListWithElements(
				3,
				value.String("foo"),
				value.SmallInt(5),
				vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("c"), Value: value.String("bar")},
					value.Pair{Key: value.ToSymbol("d"), Value: value.String("baz")},
				),
			),
		},
		"call a method with regular params, optional params, named rest param and all args": {
			source: `
				def foo(a: String, b: Int = 5, **c: String): Tuple[any]
					[a, b, c]
				end

				foo("foo", 9, c: "bar", d: "baz")
			`,
			wantStackTop: value.NewArrayListWithElements(
				3,
				value.String("foo"),
				value.SmallInt(9),
				vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("c"), Value: value.String("bar")},
					value.Pair{Key: value.ToSymbol("d"), Value: value.String("baz")},
				),
			),
		},
		"call a method with regular params, optional params, named rest param and optional named arg": {
			source: `
				def foo(a: String, b: String | Int = "b", **c: String): Tuple[any]
					[a, b, c]
				end

				foo("foo", c: "bar", d: "baz", b: 9)
			`,
			wantStackTop: value.NewArrayListWithElements(
				2,
				value.String("foo"),
				value.SmallInt(9),
				vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("d"), Value: value.String("baz")},
					value.Pair{Key: value.ToSymbol("c"), Value: value.String("bar")},
				),
			),
		},
		"call a method with positional rest params and named rest params and no args": {
			source: `
				def foo(*a: Int, **b: Int): String
					"a: #{a}, b: #{b}"
				end

				foo()
			`,
			wantStackTop: value.String(`a: %[], b: %{}`),
		},
		"call a method with positional rest params and named rest params and positional args": {
			source: `
				def foo(*a: Int, **b: Int): String
					"a: #{a}, b: #{b}"
				end

				foo(1, 5, 7)
			`,
			wantStackTop: value.String(`a: %[1, 5, 7], b: %{}`),
		},
		"call a method with positional rest params and named rest params and named args": {
			source: `
				def foo(*a: Int, **b: Int): Tuple[any]
					[a, b]
				end

				foo(foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.NewArrayListWithElements(
				2,
				&value.ArrayTuple{},
				vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo"), Value: value.SmallInt(5)},
					value.Pair{Key: value.ToSymbol("bar"), Value: value.SmallInt(2)},
					value.Pair{Key: value.ToSymbol("baz"), Value: value.SmallInt(8)},
				),
			),
		},
		"call a method with positional rest params and named rest params and both types of args": {
			source: `
				def foo(*a: Int, **b: Int): Tuple[any]
					[a, b]
				end

				foo(10, 20, 30, foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.NewArrayListWithElements(
				2,
				value.NewArrayTupleWithElements(
					3,
					value.SmallInt(10),
					value.SmallInt(20),
					value.SmallInt(30),
				),
				vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo"), Value: value.SmallInt(5)},
					value.Pair{Key: value.ToSymbol("bar"), Value: value.SmallInt(2)},
					value.Pair{Key: value.ToSymbol("baz"), Value: value.SmallInt(8)},
				),
			),
		},

		"call a method with regular, positional rest params and named rest params and no args": {
			source: `
				def foo(a: Int, *b: Int, **c: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				foo(5)
			`,
			wantStackTop: value.String(`a: 5, b: %[], c: %{}`),
		},
		"call a method with regular, positional rest params and named rest params and positional args": {
			source: `
				def foo(a: Int, *b: Int, **c: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				foo(1, 5, 7)
			`,
			wantStackTop: value.String(`a: 1, b: %[5, 7], c: %{}`),
		},
		"call a method with regular, positional rest params and named rest params and named args": {
			source: `
				def foo(a: Int, *b: Int, **c: Int): Tuple[any]
					[a, b, c]
				end

				foo(1, foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.NewArrayListWithElements(
				3,
				value.SmallInt(1),
				&value.ArrayTuple{},
				vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo"), Value: value.SmallInt(5)},
					value.Pair{Key: value.ToSymbol("bar"), Value: value.SmallInt(2)},
					value.Pair{Key: value.ToSymbol("baz"), Value: value.SmallInt(8)},
				),
			),
		},
		"call a method with regular, positional rest params and named rest params and both types of args": {
			source: `
				def foo(a: Int, *b: Int, **c: Int): Tuple[any]
					[a, b, c]
				end

				foo(10, 20, 30, foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.NewArrayListWithElements(
				3,
				value.SmallInt(10),
				value.NewArrayTupleWithElements(
					2,
					value.SmallInt(20),
					value.SmallInt(30),
				),
				vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo"), Value: value.SmallInt(5)},
					value.Pair{Key: value.ToSymbol("bar"), Value: value.SmallInt(2)},
					value.Pair{Key: value.ToSymbol("baz"), Value: value.SmallInt(8)},
				),
			),
		},
		"call a method with rest parameters and no arguments": {
			source: `
				def foo(*b: Int): String
					b.inspect
				end

				foo()
			`,
			wantStackTop: value.String("%[]"),
		},
		"call a method with rest parameters and arguments": {
			source: `
				def foo(*b: Int): String
					b.inspect
				end

				foo(1, 2, 3)
			`,
			wantStackTop: value.String("%[1, 2, 3]"),
		},
		"call a method with rest parameters and required arguments": {
			source: `
				def foo(a: Int, *b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(1)
			`,
			wantStackTop: value.String("a: 1, b: %[]"),
		},
		"call a method with rest parameters and all arguments": {
			source: `
				def foo(a: Int, *b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(1, 2, 3, 4)
			`,
			wantStackTop: value.String("a: 1, b: %[2, 3, 4]"),
		},
		"call a method with rest parameters and named args": {
			source: `
				def foo(a: Int, b: Int, *c: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(b: 1, a: 2)
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(99, 6, 5), P(113, 6, 19)), "expected 2... positional arguments in call to `foo`, got 0"),
			},
		},
		"call a method with rest parameters and no optional arguments": {
			source: `
				def foo(a: Int = 3, *b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo()
			`,
			wantStackTop: value.String("a: 3, b: %[]"),
		},
		"call a method with rest parameters and optional arguments": {
			source: `
				def foo(a: Int = 3, *b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(1)
			`,
			wantStackTop: value.String("a: 1, b: %[]"),
		},
		"call a method with rest parameters and all optional arguments": {
			source: `
				def foo(a: Int = 3, *b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(1, 2, 3)
			`,
			wantStackTop: value.String("a: 1, b: %[2, 3]"),
		},
		"call a method with post parameters": {
			source: `
				def foo(*a: Int, b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(1, 2, 3)
			`,
			wantStackTop: value.String("a: %[1, 2], b: 3"),
		},
		"call a method with multiple post arguments": {
			source: `
				def foo(*a: Int, b: Int, c: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				foo(1, 2, 3, 4)
			`,
			wantStackTop: value.String("a: %[1, 2], b: 3, c: 4"),
		},
		"call a method with pre and post arguments": {
			source: `
				def foo(a: Int, b: Int, *c: Int, d: Int, e: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}, d: ${d.inspect}, e: ${e.inspect}"
				end

				foo(1, 2, 3, 4, 5, 6)
			`,
			wantStackTop: value.String("a: 1, b: 2, c: %[3, 4], d: 5, e: 6"),
		},
		"call a method with named post arguments": {
			source: `
				def foo(*a: Int, b: Int, c: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				foo(1, 2, c: 3, b: 4)
			`,
			wantStackTop: value.String("a: %[1, 2], b: 4, c: 3"),
		},
		"call a method with pre and named post arguments": {
			source: `
				def foo(a: Int, *b: Int, c: Int, d: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}, d: ${d.inspect}"
				end

				foo(1, 2, 3, d: 4, c: 5)
			`,
			wantStackTop: value.String("a: 1, b: %[2, 3], c: 5, d: 4"),
		},
		"call a method with named arguments and missing required arguments": {
			source: `
				def foo(a: String, b: String, c: String = "default c", d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				foo("a", e: "e")
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(182, 6, 5), P(197, 6, 20)), "argument `b` is missing in call to `foo`"),
			},
		},
		"call a method with duplicated arguments": {
			source: `
				def foo(a: String, b: String): String
					"a: ${a}, b: ${b}"
				end

				foo("a", b: "b", a: "a2")
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(97, 6, 22), P(103, 6, 28)), "duplicated argument `a` in call to `foo`"),
			},
		},
		"call a method with unknown named arguments": {
			source: `
				def foo(a: String, b: String): String
					"a: ${a}, b: ${b}"
				end

				foo("a", unknown: "lala", moo: "meow", b: "b")
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(89, 6, 14), P(103, 6, 28)), "nonexistent parameter `unknown` given in call to `foo`"),
				error.NewFailure(L(P(106, 6, 31), P(116, 6, 41)), "nonexistent parameter `moo` given in call to `foo`"),
			},
		},
		"call a module method without arguments": {
			source: `
				module Foo
					def bar: Symbol
						:baz
					end
				end

				::Foo.bar
			`,
			wantStackTop: value.ToSymbol("baz"),
		},
		"call a module method with positional arguments": {
			source: `
				module Foo
					def add(a: Int, b: Int): Int
						a + b
					end
				end

				::Foo.add 4, 12
			`,
			wantStackTop: value.SmallInt(16),
		},
		"call an instance method without arguments": {
			source: `
				class ::Std::Object < Value
					def bar: Symbol
						:baz
					end
				end

				self.bar
			`,
			wantStackTop: value.ToSymbol("baz"),
		},
		"call an instance method with positional arguments": {
			source: `
				class Std::Object < Std::Value
					def add(a: Int, b: Int): Int
						a + b
					end
				end

				self.add 1, 8
			`,
			wantStackTop: value.SmallInt(9),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Setters(t *testing.T) {
	tests := sourceTestTable{
		"call a manually defined setter": {
			source: `
				def foo=(v: Int)
					:bar
				end

				Kernel.foo = 3
			`,
			wantStackTop: value.SmallInt(3),
		},
		"setter increment type error": {
			source: `
				class Foo
				  attr bar: Int?
					init(@bar: Int?); end
				end

				foo := ::Foo(1)
				foo.bar++
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(96, 8, 5), P(104, 8, 13)), "method `++` is not defined on type `Std::Nil`"),
			},
		},
		"setter increment": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(1)
				foo.bar++
			`,
			wantStackTop: value.SmallInt(2),
		},
		"setter decrement": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(1)
				foo.bar--
			`,
			wantStackTop: value.SmallInt(0),
		},
		"setter add": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(1)
				foo.bar += 2
			`,
			wantStackTop: value.SmallInt(3),
		},
		"setter subtract": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(1)
				foo.bar -= 2
			`,
			wantStackTop: value.SmallInt(-1),
		},
		"setter multiply": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(3)
				foo.bar *= 2
			`,
			wantStackTop: value.SmallInt(6),
		},
		"setter divide": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(12)
				foo.bar /= 4
			`,
			wantStackTop: value.SmallInt(3),
		},
		"setter exponentiate": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(12)
				foo.bar **= 2
			`,
			wantStackTop: value.SmallInt(144),
		},
		"setter modulo": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(12)
				foo.bar %= 5
			`,
			wantStackTop: value.SmallInt(2),
		},
		"setter left bitshift": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(5)
				foo.bar <<= 2
			`,
			wantStackTop: value.SmallInt(20),
		},
		"setter logic left bitshift": {
			source: `
				class Foo
				  attr bar: Int8
					init(@bar: Int8); end
				end

				foo := ::Foo(5i8)
				foo.bar <<<= 2i8
			`,
			wantStackTop: value.Int8(20),
		},
		"setter right bitshift": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(10)
				foo.bar >>= 2
			`,
			wantStackTop: value.SmallInt(2),
		},
		"setter logic right bitshift": {
			source: `
				class Foo
				  attr bar: Int8
					init(@bar: Int8); end
				end

				foo := ::Foo(10i8)
				foo.bar >>>= 2i8
			`,
			wantStackTop: value.Int8(2),
		},
		"setter bitwise and": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(6)
				foo.bar &= 5
			`,
			wantStackTop: value.SmallInt(4),
		},
		"setter bitwise or": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(6)
				foo.bar |= 5
			`,
			wantStackTop: value.SmallInt(7),
		},
		"setter bitwise xor": {
			source: `
				class Foo
				  attr bar: Int
					init(@bar: Int); end
				end

				foo := ::Foo(6)
				foo.bar ^= 5
			`,
			wantStackTop: value.SmallInt(3),
		},
		"setter logic or falsy": {
			source: `
				class Foo
				  attr bar: Int?
					init(@bar: Int?); end
				end

				foo := ::Foo(nil)
				foo.bar ||= 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"setter logic or truthy": {
			source: `
				class Foo
				  attr bar: Int?
					init(@bar: Int?); end
				end

				foo := ::Foo(2)
				foo.bar ||= 5
			`,
			wantStackTop: value.SmallInt(2),
		},
		"setter logic and nil": {
			source: `
				class Foo
				  attr bar: Int?
					init(@bar: Int?); end
				end

				foo := ::Foo(nil)
				foo.bar &&= 5
			`,
			wantStackTop: value.Nil,
		},
		"setter logic and truthy": {
			source: `
				class Foo
				  attr bar: Int?
					init(@bar: Int?); end
				end

				foo := ::Foo(2)
				foo.bar &&= 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"setter nil coalesce falsy": {
			source: `
				class Foo
				  attr bar: Int?
					init(@bar: Int?); end
				end

				foo := ::Foo(nil)
				foo.bar ??= 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"setter nil coalesce truthy": {
			source: `
				class Foo
				  attr bar: Int?
					init(@bar: Int?); end
				end

				foo := ::Foo(2)
				foo.bar ??= 5
			`,
			wantStackTop: value.SmallInt(2),
		},
		"call subscript set": {
			source: `
				var list: List[Int | Symbol] = [5, 8, 20]
				list[0] = :foo
				list
			`,
			wantStackTop: &value.ArrayList{value.ToSymbol("foo"), value.SmallInt(8), value.SmallInt(20)},
		},
		"subscript return value": {
			source: `
				list := ["foo", 2, 7.8]
				list[0] = 8
			`,
			wantStackTop: value.SmallInt(8),
		},
		"set index 0 of a list": {
			source: `
				list := ["foo", 2, 7.8]
				list[0] = 8
				list
			`,
			wantStackTop: &value.ArrayList{
				value.SmallInt(8),
				value.SmallInt(2),
				value.Float(7.8),
			},
		},
		"subscript setter increment": {
			source: `
				list := [5, 2, 7]
				list[0]++
			`,
			wantStackTop: value.SmallInt(6),
		},
		"subscript setter decrement": {
			source: `
				list := [5, 2, 7]
				list[0]--
			`,
			wantStackTop: value.SmallInt(4),
		},
		"subscript setter add": {
			source: `
				list := [5, 2, 7]
				list[0] += 8
			`,
			wantStackTop: value.SmallInt(13),
		},
		"subscript setter subtract": {
			source: `
				list := [5, 2, 7]
				list[0] -= 8
			`,
			wantStackTop: value.SmallInt(-3),
		},
		"subscript setter multiply": {
			source: `
				list := [5, 2, 7]
				list[1] *= 3
			`,
			wantStackTop: value.SmallInt(6),
		},
		"subscript setter divide": {
			source: `
				list := [5, 8, 7]
				list[1] /= 2
			`,
			wantStackTop: value.SmallInt(4),
		},
		"subscript setter exponentiate": {
			source: `
				list := [5, 8, 7]
				list[1] **= 2
			`,
			wantStackTop: value.SmallInt(64),
		},
		"subscript setter modulo type error": {
			source: `
				list := [5, 8, 7]
				list[0] %= 2
			`,
			wantStackTop: value.SmallInt(1),
		},
		"subscript setter modulo": {
			source: `
				list := [5, 8, 7]
				list[0] %= 2
			`,
			wantStackTop: value.SmallInt(1),
		},
		"subscript setter left bitshift": {
			source: `
				list := [5, 8, 7]
				list[0] <<= 2
			`,
			wantStackTop: value.SmallInt(20),
		},
		"subscript setter logic left bitshift": {
			source: `
				list := [5i8, 8i8, 7i8]
				list[0] <<<= 2
			`,
			wantStackTop: value.Int8(20),
		},
		"subscript setter right bitshift": {
			source: `
				list := [10, 8, 7]
				list[0] >>= 2
			`,
			wantStackTop: value.SmallInt(2),
		},
		"subscript setter logic right bitshift": {
			source: `
				list := [10i8, 8i8, 7i8]
				list[0] >>>= 2
			`,
			wantStackTop: value.Int8(2),
		},
		"subscript setter bitwise and": {
			source: `
				list := [6, 8, 7]
				list[0] &= 5
			`,
			wantStackTop: value.SmallInt(4),
		},
		"subscript setter bitwise or": {
			source: `
				list := [6, 8, 7]
				list[0] |= 5
			`,
			wantStackTop: value.SmallInt(7),
		},
		"subscript setter bitwise xor": {
			source: `
				list := [6, 8, 7]
				list[0] ^= 5
			`,
			wantStackTop: value.SmallInt(3),
		},
		"subscript setter logic or falsy": {
			source: `
				list := [nil, 8, 7.8]
				list[0] ||= 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"subscript setter logic or truthy": {
			source: `
				list := [1, 8, 7.8]
				list[0] ||= 5
			`,
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(29, 3, 5), P(35, 3, 11)), "this condition will always have the same result since type `Std::Int | Std::Float` is truthy"),
				error.NewWarning(L(P(41, 3, 17), P(41, 3, 17)), "unreachable code"),
			},
			wantStackTop: value.SmallInt(1),
		},
		"subscript setter logic and nil": {
			source: `
				list := [nil, 8, 7.8]
				list[0] &&= 5
			`,
			wantStackTop: value.Nil,
		},
		"subscript setter logic and false": {
			source: `
				list := [false, 8, 7.8]
				list[0] &&= 5
			`,
			wantStackTop: value.False,
		},
		"subscript setter logic and truthy": {
			source: `
				list := [1, 8, 7.8]
				list[0] &&= 5
			`,
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(29, 3, 5), P(35, 3, 11)), "this condition will always have the same result since type `Std::Int | Std::Float` is truthy"),
			},
			wantStackTop: value.SmallInt(5),
		},
		"subscript setter nil coalesce nil": {
			source: `
				list := [nil, 8, 7.8]
				list[0] ??= 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"subscript setter nil coalesce false": {
			source: `
				list := [false, 8, 7.8]
				list[0] ??= 5
			`,
			wantStackTop: value.False,
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(33, 3, 5), P(39, 3, 11)), "this condition will always have the same result since type `bool | Std::Int | Std::Float` can never be nil"),
				error.NewWarning(L(P(45, 3, 17), P(45, 3, 17)), "unreachable code"),
			},
		},
		"subscript setter nil coalesce truthy": {
			source: `
				list := [1, 8, 7.8]
				list[0] ??= 5
			`,
			wantStackTop: value.SmallInt(1),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(29, 3, 5), P(35, 3, 11)), "this condition will always have the same result since type `Std::Int | Std::Float` can never be nil"),
				error.NewWarning(L(P(41, 3, 17), P(41, 3, 17)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
