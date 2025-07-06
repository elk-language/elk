package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
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
			wantStackTop: value.Ref(value.String("foo")),
		},
		"get index -1 of a list": {
			source: `
				list := ["foo", 2, 7.8]
				list[-1]
			`,
			wantStackTop: value.Float(7.8).ToValue(),
		},
		"get too big index": {
			source: `
				list := ["foo", 2, 7.8]
				list[50]
			`,
			wantRuntimeErr: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 50 out of range: -3...3",
			)),
		},
		"get too small index": {
			source: `
				list := ["foo", 2, 7.8]
				list[-10]
			`,
			wantRuntimeErr: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index -10 out of range: -3...3",
			)),
		},
		"get from nil": {
			source: `
				list := nil
				list[-10]
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(21, 3, 5), P(29, 3, 13)), "method `[]` is not defined on type `Std::Nil`"),
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
			wantStackTop: value.Ref(value.String("foo")),
		},
		"get index -1 of a list": {
			source: `
				var list: List[String | Int | Float]? = ["foo", 2, 7.8]
				list?[-1]
			`,
			wantStackTop: value.Float(7.8).ToValue(),
		},
		"get too big index": {
			source: `
				var list: List[String | Int | Float]? = ["foo", 2, 7.8]
				list?[50]
			`,
			wantRuntimeErr: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index 50 out of range: -3...3",
			)),
		},
		"get too small index": {
			source: `
				var list: List[String | Int | Float]? = ["foo", 2, 7.8]
				list?[-10]
			`,
			wantRuntimeErr: value.Ref(value.NewError(
				value.IndexErrorClass,
				"index -10 out of range: -3...3",
			)),
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
			wantStackTop: value.Ref(value.NewObject(
				value.ObjectWithClass(
					value.NewClassWithOptions(
						value.ClassWithName("Foo"),
					),
				),
			)),
		},
		"instantiate a class without an initialiser with arguments": {
			source: `
				class Foo; end

				::Foo(2)
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(25, 4, 5), P(32, 4, 12)), "expected 0 arguments in call to `Foo.:#init`, got 1"),
			},
		},
		"instantiate a class with an initialiser without arguments": {
			source: `
				class Foo
					init(a: String); end
				end

				::Foo()
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(54, 6, 5), P(60, 6, 11)), "argument `a` is missing in call to `Foo.:#init`"),
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
			wantStackTop: value.Ref(value.String("bar")),
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
				sealed primitive noinit class ::Std::Int < Value
					alias add +
				end

				3.add(4)
			`,
			wantStackTop: value.SmallInt(7).ToValue(),
		},
		"add an alias to a nonexistent method": {
			source: `
				alias foo blabla
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(11, 2, 11), P(20, 2, 20)), "method `blabla` is not defined on type `Std::Kernel`"),
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
			wantStackTop: value.ToSymbol("bar").ToValue(),
		},
		"define a method with positional arguments in top level": {
			source: `
				def foo(a: Int, b: Int): Int
					c := 5
					a + b + c
				end
				foo(1, 2)
			`,
			wantStackTop: value.SmallInt(8).ToValue(),
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
			wantStackTop: value.SmallInt(8).ToValue(),
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
			wantStackTop: value.SmallInt(8).ToValue(),
		},
		"define a generator": {
			source: `
				module Bar
					def *foo(a: Int, b: Int): Int
						c := 5
						a + b + c
					end
				end

				println Bar.foo(1, 2).inspect
			`,
			wantStdout:   "Std::Generator{location: sourceName:3:6}\n",
			wantStackTop: value.Nil,
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
			wantStackTop: value.SmallInt(25).ToValue(),
		},
		"invalid args": {
			source: `
				pow2 := |a: Int| -> a ** 2
				pow2.call(5, 8)
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(36, 3, 5), P(50, 3, 19)), "expected 1 arguments in call to `call`, got 2"),
			},
		},
		"call closure in a native method": {
			source: `
				5.times |i| ->
					println i
				end
			`,
			wantStackTop: value.Nil,
			wantStdout:   "0\n1\n2\n3\n4\n",
		},
		"call closure that throws in a native method": {
			source: `
				5.times |i| ->
					throw unchecked i
				end
			`,
			wantRuntimeErr: value.SmallInt(0).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Async(t *testing.T) {
	tests := sourceTestTable{
		"start 8 promises and await them": {
			source: `
				async def foo: Int
					println "START"
					await timeout(10.millisecond)
					println "STOP"

					10
				end

				var promises: List[Promise[Int]] = []
				for i in 1...8 then promises << foo()
				for p in promises then await p
			`,
			wantStdout:   "START\nSTART\nSTART\nSTART\nSTART\nSTART\nSTART\nSTART\nSTOP\nSTOP\nSTOP\nSTOP\nSTOP\nSTOP\nSTOP\nSTOP\n",
			wantStackTop: value.Nil,
		},
		"await a promise that throws": {
			source: `
				def lol: String
					throw unchecked 5
				end
				async def foo: String
					lol() + "u"
				end
				async def bar: String
				  await foo()
				end
				async def baz: String
					await bar()
				end

				baz().await_sync
			`,
			wantRuntimeErr: value.SmallInt(5).ToValue(),
			wantStackTrace: &value.StackTrace{
				{FuncName: "sourceName", FileName: "sourceName", LineNumber: 15},
				{FuncName: "baz", FileName: "sourceName", LineNumber: 12},
				{FuncName: "bar", FileName: "sourceName", LineNumber: 9},
				{FuncName: "foo", FileName: "sourceName", LineNumber: 6},
				{FuncName: "lol", FileName: "sourceName", LineNumber: 3},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Generator(t *testing.T) {
	tests := sourceTestTable{
		"call a generator": {
			source: `
				def *foo(a: Int, b: Int = 9): Int
					yield a + b
					yield 6

					0
				end

				g := foo(3)
				println(try g.next)
				println(try g.next)
				println(try g.next)
			`,
			wantStdout:   "12\n6\n0\n",
			wantStackTop: value.Nil,
		},
		"throws :stop_iteration after the last yield": {
			source: `
				def *foo(): Int
					5
				end

				g := foo()
				println(try g.next)
				for i in 1...3
					do
						g.next
					catch :stop_iteration
						println "caught"
					end
				end
			`,
			wantStdout:   "5\ncaught\ncaught\ncaught\n",
			wantStackTop: value.Nil,
		},
		"can throw custom errors": {
			source: `
				def *foo(): Int ! String
					yield 5
					throw "bar"
				end

				g := foo()
				println(try g.next)
				println(try g.next)
			`,
			wantStdout:     "5\n",
			wantRuntimeErr: value.Ref(value.String("bar")),
		},
		"throws stop_iteration after a custom error had been thrown": {
			source: `
				def *foo(): Int ! String
					throw "bar"
				end

				g := foo()
				do
					try g.next
				catch String() as str
					println "caught custom"
				end

				for i in 1...3
					do
						try g.next
					catch :stop_iteration
						println "caught stop"
					end
				end
			`,
			wantStdout:   "caught custom\ncaught stop\ncaught stop\ncaught stop\n",
			wantStackTop: value.Nil,
		},
		"can be used in for loops": {
			source: `
				def *foo(a: Int, b: Int = 9): Int
					yield a + b
					yield 6

					0
				end

				for i in foo(3)
					println i
				end
			`,
			wantStdout:   "12\n6\n0\n",
			wantStackTop: value.Nil,
		},
		"can be reset": {
			source: `
				def *foo(a: Int, b: Int = 9): Int
					yield a + b
					yield 6

					0
				end

				g := foo(3)
				println(try g.next)
				println(try g.next)
				g.reset
				println(try g.next)
			`,
			wantStdout:   "12\n6\n12\n",
			wantStackTop: value.Nil,
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
			wantStackTop: value.Ref(value.String("5")),
		},
		"call a variable": {
			source: `
				module Foo
					def call: Symbol
						:bar
					end
				end

				a := Foo
				a()
			`,
			wantStackTop: value.ToSymbol("bar").ToValue(),
		},
		"call method from using": {
			source: `
				using Foo::bar

				module Foo
					def bar: Symbol
						:bar
					end
				end

				bar()
			`,
			wantStackTop: value.ToSymbol("bar").ToValue(),
		},
		"call method from using all": {
			source: `
				using Foo::*

				module Foo
					def bar: Symbol
						:bar
					end
				end

				bar()
			`,
			wantStackTop: value.ToSymbol("bar").ToValue(),
		},
		"call variable": {
			source: `
				foo := |n: Int| -> n ** 2
				foo(5)
			`,
			wantStackTop: value.SmallInt(25).ToValue(),
		},
		"call a global method without arguments": {
			source: `
				def foo: Symbol
					:bar
				end

				foo()
			`,
			wantStackTop: value.ToSymbol("bar").ToValue(),
		},
		"call a global method with positional arguments": {
			source: `
				def add(a: Int, b: Int): Int
					a + b
				end

				add(5, 9)
			`,
			wantStackTop: value.SmallInt(14).ToValue(),
		},
		"call a method with missing required arguments": {
			source: `
				def add(a: Int, b: Int): Int
					a + b
				end

				add(5)
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(58, 6, 5), P(63, 6, 10)), "argument `b` is missing in call to `Std::Kernel::add`"),
			},
		},
		"call a method without optional arguments": {
			source: `
				def add(a: Int, b: Int = 3, c: Float = 20.5): Float
					a + b + c
				end

				add(5)
			`,
			wantStackTop: value.Float(28.5).ToValue(),
		},
		"call a method with some optional arguments": {
			source: `
				def add(a: Int, b: Int = 3, c: Float = 20.5): Float
					a + b + c
				end

				add(5, 0)
			`,
			wantStackTop: value.Float(25.5).ToValue(),
		},
		"call a method with all optional arguments": {
			source: `
				def add(a: Int, b: Int = 3, c: Float = 20.5): Float
					a + b + c
				end

				add(3, 2, 3.5)
			`,
			wantStackTop: value.Float(8.5).ToValue(),
		},
		"call a method with only named arguments": {
			source: `
				def foo(a: String, b: String, c: String, d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				foo(b: "b", a: "a", c: "c", e: "e")
			`,
			wantStackTop: value.Ref(value.String("a: a, b: b, c: c, d: default d, e: e")),
		},
		"call a method with all required arguments and named arguments": {
			source: `
				def foo(a: String, b: String, c: String, d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				foo("a", c: "c", b: "b")
			`,
			wantStackTop: value.Ref(value.String("a: a, b: b, c: c, d: default d, e: default e")),
		},
		"call a method with optional arguments and named arguments": {
			source: `
				def foo(a: String, b: String, c: String = "default c", d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				foo("a", "b", "c", e: "e")
			`,
			wantStackTop: value.Ref(value.String("a: a, b: b, c: c, d: default d, e: e")),
		},
		"call a method with a named rest param and no args": {
			source: `
				def foo(**a: String): String
					"a: #{a}"
				end

				foo()
			`,
			wantStackTop: value.Ref(value.String("a: %{}")),
		},
		"call a method with a named rest param and a few named args": {
			source: `
				def foo(**a: String): Record[Symbol, String]
					a
				end

				foo(d: "foo", a: "bar")
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{Key: value.ToSymbol("a").ToValue(), Value: value.Ref(value.String("bar"))},
				value.Pair{Key: value.ToSymbol("d").ToValue(), Value: value.Ref(value.String("foo"))},
			)),
		},
		"call a method with a named rest param and a double splat argument": {
			source: `
				def foo(**a: Int): Record[Symbol, Int]
					a
				end

				map := { foo: 1, bar: 2 }
				foo(**map)
			`,
			wantStackTop: value.Ref(vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.ToSymbol("foo").ToValue(), Value: value.SmallInt(1).ToValue()},
				value.Pair{Key: value.ToSymbol("bar").ToValue(), Value: value.SmallInt(2).ToValue()},
			)),
		},
		"call a method with a named rest param, named and double splat arguments": {
			source: `
				def foo(**a: Int): Record[Symbol, Int]
					a
				end

				map := { foo: 1, bar: 2 }
				foo(a: 20, **map, b: 9)
			`,
			wantStackTop: value.Ref(vm.MustNewHashRecordWithElements(
				nil,
				value.Pair{Key: value.ToSymbol("a").ToValue(), Value: value.SmallInt(20).ToValue()},
				value.Pair{Key: value.ToSymbol("b").ToValue(), Value: value.SmallInt(9).ToValue()},
				value.Pair{Key: value.ToSymbol("foo").ToValue(), Value: value.SmallInt(1).ToValue()},
				value.Pair{Key: value.ToSymbol("bar").ToValue(), Value: value.SmallInt(2).ToValue()},
			)),
		},
		"call a method with regular params, named rest param and a few named args": {
			source: `
				def foo(a: String, **b: String): Tuple[any]
					[a, b]
				end

				foo("foo", c: "bar", d: "baz")
			`,
			wantStackTop: value.Ref(value.NewArrayListWithElements(
				2,
				value.Ref(value.String("foo")),
				value.Ref(vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("c").ToValue(), Value: value.Ref(value.String("bar"))},
					value.Pair{Key: value.ToSymbol("d").ToValue(), Value: value.Ref(value.String("baz"))},
				)),
			)),
		},
		"call a method with regular params, named rest param and only required args": {
			source: `
				def foo(a: String, **b: String): Tuple[any]
					[a, b]
				end

				foo("foo")
			`,
			wantStackTop: value.Ref(value.NewArrayListWithElements(
				2,
				value.Ref(value.String("foo")),
				value.Ref(value.NewHashRecord(0)),
			)),
		},
		"call a method with regular params, optional params, named rest param and a few named args": {
			source: `
				def foo(a: String, b: Int = 5, **c: String): Tuple[any]
					[a, b, c]
				end

				foo("foo", c: "bar", d: "baz")
			`,
			wantStackTop: value.Ref(value.NewArrayListWithElements(
				3,
				value.Ref(value.String("foo")),
				value.SmallInt(5).ToValue(),
				value.Ref(vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("c").ToValue(), Value: value.Ref(value.String("bar"))},
					value.Pair{Key: value.ToSymbol("d").ToValue(), Value: value.Ref(value.String("baz"))},
				)),
			)),
		},
		"call a method with regular params, optional params, named rest param and all args": {
			source: `
				def foo(a: String, b: Int = 5, **c: String): Tuple[any]
					[a, b, c]
				end

				foo("foo", 9, c: "bar", d: "baz")
			`,
			wantStackTop: value.Ref(value.NewArrayListWithElements(
				3,
				value.Ref(value.String("foo")),
				value.SmallInt(9).ToValue(),
				value.Ref(vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("c").ToValue(), Value: value.Ref(value.String("bar"))},
					value.Pair{Key: value.ToSymbol("d").ToValue(), Value: value.Ref(value.String("baz"))},
				)),
			)),
		},
		"call a method with regular params, optional params, named rest param and optional named arg": {
			source: `
				def foo(a: String, b: String | Int = "b", **c: String): Tuple[any]
					[a, b, c]
				end

				foo("foo", c: "bar", d: "baz", b: 9)
			`,
			wantStackTop: value.Ref(value.NewArrayListWithElements(
				3,
				value.Ref(value.String("foo")),
				value.SmallInt(9).ToValue(),
				value.Ref(vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("d").ToValue(), Value: value.Ref(value.String("baz"))},
					value.Pair{Key: value.ToSymbol("c").ToValue(), Value: value.Ref(value.String("bar"))},
				)),
			)),
		},
		"call a method with positional rest params and named rest params and no args": {
			source: `
				def foo(*a: Int, **b: Int): String
					"a: #{a}, b: #{b}"
				end

				foo()
			`,
			wantStackTop: value.Ref(value.String(`a: %[], b: %{}`)),
		},
		"call a method with positional rest params and named rest params and positional args": {
			source: `
				def foo(*a: Int, **b: Int): String
					"a: #{a}, b: #{b}"
				end

				foo(1, 5, 7)
			`,
			wantStackTop: value.Ref(value.String(`a: %[1, 5, 7], b: %{}`)),
		},
		"call a method with positional rest params and named rest params and named args": {
			source: `
				def foo(*a: Int, **b: Int): Tuple[any]
					[a, b]
				end

				foo(foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.Ref(value.NewArrayListWithElements(
				2,
				value.Ref(&value.ArrayTuple{}),
				value.Ref(vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo").ToValue(), Value: value.SmallInt(5).ToValue()},
					value.Pair{Key: value.ToSymbol("bar").ToValue(), Value: value.SmallInt(2).ToValue()},
					value.Pair{Key: value.ToSymbol("baz").ToValue(), Value: value.SmallInt(8).ToValue()},
				)),
			)),
		},
		"call a method with positional rest params and named rest params and both types of args": {
			source: `
				def foo(*a: Int, **b: Int): Tuple[any]
					[a, b]
				end

				foo(10, 20, 30, foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.Ref(value.NewArrayListWithElements(
				2,
				value.Ref(value.NewArrayTupleWithElements(
					3,
					value.SmallInt(10).ToValue(),
					value.SmallInt(20).ToValue(),
					value.SmallInt(30).ToValue(),
				)),
				value.Ref(vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo").ToValue(), Value: value.SmallInt(5).ToValue()},
					value.Pair{Key: value.ToSymbol("bar").ToValue(), Value: value.SmallInt(2).ToValue()},
					value.Pair{Key: value.ToSymbol("baz").ToValue(), Value: value.SmallInt(8).ToValue()},
				)),
			)),
		},

		"call a method with regular, positional rest params and named rest params and no args": {
			source: `
				def foo(a: Int, *b: Int, **c: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				foo(5)
			`,
			wantStackTop: value.Ref(value.String(`a: 5, b: %[], c: %{}`)),
		},
		"call a method with regular, positional rest params and named rest params and positional args": {
			source: `
				def foo(a: Int, *b: Int, **c: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				foo(1, 5, 7)
			`,
			wantStackTop: value.Ref(value.String(`a: 1, b: %[5, 7], c: %{}`)),
		},
		"call a method with regular, positional rest params and named rest params and named args": {
			source: `
				def foo(a: Int, *b: Int, **c: Int): Tuple[any]
					[a, b, c]
				end

				foo(1, foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.Ref(value.NewArrayListWithElements(
				3,
				value.SmallInt(1).ToValue(),
				value.Ref(&value.ArrayTuple{}),
				value.Ref(vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo").ToValue(), Value: value.SmallInt(5).ToValue()},
					value.Pair{Key: value.ToSymbol("bar").ToValue(), Value: value.SmallInt(2).ToValue()},
					value.Pair{Key: value.ToSymbol("baz").ToValue(), Value: value.SmallInt(8).ToValue()},
				)),
			)),
		},
		"call a method with regular, positional rest params and named rest params and both types of args": {
			source: `
				def foo(a: Int, *b: Int, **c: Int): Tuple[any]
					[a, b, c]
				end

				foo(10, 20, 30, foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.Ref(value.NewArrayListWithElements(
				3,
				value.SmallInt(10).ToValue(),
				value.Ref(value.NewArrayTupleWithElements(
					2,
					value.SmallInt(20).ToValue(),
					value.SmallInt(30).ToValue(),
				)),
				value.Ref(vm.MustNewHashRecordWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo").ToValue(), Value: value.SmallInt(5).ToValue()},
					value.Pair{Key: value.ToSymbol("bar").ToValue(), Value: value.SmallInt(2).ToValue()},
					value.Pair{Key: value.ToSymbol("baz").ToValue(), Value: value.SmallInt(8).ToValue()},
				)),
			)),
		},
		"call a method with rest parameters and no arguments": {
			source: `
				def foo(*b: Int): String
					b.inspect
				end

				foo()
			`,
			wantStackTop: value.Ref(value.String("%[]")),
		},
		"call a method with rest parameter and list splat argument": {
			source: `
				def foo(*b: Int): String
					b.inspect
				end

				list := [5, 9, 2]
				foo(*list)
			`,
			wantStackTop: value.Ref(value.String("[5, 9, 2]")),
		},
		"call a method with rest parameter and splat argument": {
			source: `
				def foo(*b: Int): String
					b.inspect
				end

				foo(*3)
			`,
			wantStackTop: value.Ref(value.String("%[0, 1, 2]")),
		},
		"call a method with rest parameter and rest + splat arguments": {
			source: `
				def foo(*b: Int): String
					b.inspect
				end

				list := [85, 22]
				foo(20, *3, 9, *list, 17)
			`,
			wantStackTop: value.Ref(value.String("%[20, 0, 1, 2, 9, 85, 22, 17]")),
		},
		"call a method with rest parameters and arguments": {
			source: `
				def foo(*b: Int): String
					b.inspect
				end

				foo(1, 2, 3)
			`,
			wantStackTop: value.Ref(value.String("%[1, 2, 3]")),
		},
		"call a method with rest parameters and required arguments": {
			source: `
				def foo(a: Int, *b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(1)
			`,
			wantStackTop: value.Ref(value.String("a: 1, b: %[]")),
		},
		"call a method with rest parameters and rest + splat arguments": {
			source: `
				def foo(a: Int, *b: Int): String
					"a: #a, b: #b"
				end

				list := [85, 22]
				foo(20, *3, 9, *list, 17)
			`,
			wantStackTop: value.Ref(value.String("a: 20, b: %[0, 1, 2, 9, 85, 22, 17]")),
		},
		"call a method with rest parameters and all arguments": {
			source: `
				def foo(a: Int, *b: Int): String
					"a: #a, b: #b"
				end

				foo(1, 2, 3, 4)
			`,
			wantStackTop: value.Ref(value.String("a: 1, b: %[2, 3, 4]")),
		},
		"call a method with rest parameters and named args": {
			source: `
				def foo(a: Int, b: Int, *c: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(b: 1, a: 2)
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(99, 6, 5), P(113, 6, 19)), "expected 2... positional arguments in call to `Std::Kernel::foo`, got 0"),
			},
		},
		"call a method with rest parameters and no optional arguments": {
			source: `
				def foo(a: Int = 3, *b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo()
			`,
			wantStackTop: value.Ref(value.String("a: 3, b: %[]")),
		},
		"call a method with rest parameters and optional arguments": {
			source: `
				def foo(a: Int = 3, *b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(1)
			`,
			wantStackTop: value.Ref(value.String("a: 1, b: %[]")),
		},
		"call a method with rest parameters and all optional arguments": {
			source: `
				def foo(a: Int = 3, *b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(1, 2, 3)
			`,
			wantStackTop: value.Ref(value.String("a: 1, b: %[2, 3]")),
		},
		"call a method with post parameters": {
			source: `
				def foo(*a: Int, b: Int): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				foo(1, 2, 3)
			`,
			wantStackTop: value.Ref(value.String("a: %[1, 2], b: 3")),
		},
		"call a method with post parameters and splat arguments": {
			source: `
				def foo(*a: Int, b: Int): String
					"a: #a, b: #b"
				end

				foo(*3, 2, 3)
			`,
			wantStackTop: value.Ref(value.String("a: %[0, 1, 2, 2], b: 3")),
		},
		"call a method with multiple post arguments": {
			source: `
				def foo(*a: Int, b: Int, c: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				foo(1, 2, 3, 4)
			`,
			wantStackTop: value.Ref(value.String("a: %[1, 2], b: 3, c: 4")),
		},
		"call a method with pre and post arguments": {
			source: `
				def foo(a: Int, b: Int, *c: Int, d: Int, e: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}, d: ${d.inspect}, e: ${e.inspect}"
				end

				foo(1, 2, 3, 4, 5, 6)
			`,
			wantStackTop: value.Ref(value.String("a: 1, b: 2, c: %[3, 4], d: 5, e: 6")),
		},
		"call a method with named post arguments": {
			source: `
				def foo(*a: Int, b: Int, c: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				foo(1, 2, c: 3, b: 4)
			`,
			wantStackTop: value.Ref(value.String("a: %[1, 2], b: 4, c: 3")),
		},
		"call a method with pre and named post arguments": {
			source: `
				def foo(a: Int, *b: Int, c: Int, d: Int): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}, d: ${d.inspect}"
				end

				foo(1, 2, 3, d: 4, c: 5)
			`,
			wantStackTop: value.Ref(value.String("a: 1, b: %[2, 3], c: 5, d: 4")),
		},
		"call a method with named arguments and missing required arguments": {
			source: `
				def foo(a: String, b: String, c: String = "default c", d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				foo("a", e: "e")
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(182, 6, 5), P(197, 6, 20)), "argument `b` is missing in call to `Std::Kernel::foo`"),
			},
		},
		"call a method with duplicated arguments": {
			source: `
				def foo(a: String, b: String): String
					"a: ${a}, b: ${b}"
				end

				foo("a", b: "b", a: "a2")
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(97, 6, 22), P(103, 6, 28)), "duplicated argument `a` in call to `Std::Kernel::foo`"),
			},
		},
		"call a method with unknown named arguments": {
			source: `
				def foo(a: String, b: String): String
					"a: ${a}, b: ${b}"
				end

				foo("a", unknown: "lala", moo: "meow", b: "b")
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(89, 6, 14), P(103, 6, 28)), "nonexistent parameter `unknown` given in call to `Std::Kernel::foo`"),
				diagnostic.NewFailure(L(P(106, 6, 31), P(116, 6, 41)), "nonexistent parameter `moo` given in call to `Std::Kernel::foo`"),
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
			wantStackTop: value.ToSymbol("baz").ToValue(),
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
			wantStackTop: value.SmallInt(16).ToValue(),
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
			wantStackTop: value.ToSymbol("baz").ToValue(),
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
			wantStackTop: value.SmallInt(9).ToValue(),
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
			wantStackTop: value.SmallInt(3).ToValue(),
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
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(96, 8, 5), P(104, 8, 13)), "method `++` is not defined on type `Std::Nil`"),
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
			wantStackTop: value.SmallInt(2).ToValue(),
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
			wantStackTop: value.SmallInt(0).ToValue(),
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
			wantStackTop: value.SmallInt(3).ToValue(),
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
			wantStackTop: value.SmallInt(-1).ToValue(),
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
			wantStackTop: value.SmallInt(6).ToValue(),
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
			wantStackTop: value.SmallInt(3).ToValue(),
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
			wantStackTop: value.SmallInt(144).ToValue(),
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
			wantStackTop: value.SmallInt(2).ToValue(),
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
			wantStackTop: value.SmallInt(20).ToValue(),
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
			wantStackTop: value.Int8(20).ToValue(),
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
			wantStackTop: value.SmallInt(2).ToValue(),
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
			wantStackTop: value.Int8(2).ToValue(),
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
			wantStackTop: value.SmallInt(4).ToValue(),
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
			wantStackTop: value.SmallInt(7).ToValue(),
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
			wantStackTop: value.SmallInt(3).ToValue(),
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
			wantStackTop: value.SmallInt(5).ToValue(),
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
			wantStackTop: value.SmallInt(2).ToValue(),
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
			wantStackTop: value.SmallInt(5).ToValue(),
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
			wantStackTop: value.SmallInt(5).ToValue(),
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
			wantStackTop: value.SmallInt(2).ToValue(),
		},
		"call subscript set": {
			source: `
				var list: List[Int | Symbol] = [5, 8, 20]
				list[0] = :foo
				list
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.ToSymbol("foo").ToValue(),
				value.SmallInt(8).ToValue(),
				value.SmallInt(20).ToValue(),
			}),
		},
		"subscript return value": {
			source: `
				list := ["foo", 2, 7.8]
				list[0] = 8
			`,
			wantStackTop: value.SmallInt(8).ToValue(),
		},
		"set index 0 of a list": {
			source: `
				list := ["foo", 2, 7.8]
				list[0] = 8
				list
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(8).ToValue(),
				value.SmallInt(2).ToValue(),
				value.Float(7.8).ToValue(),
			}),
		},
		"subscript setter increment": {
			source: `
				list := [5, 2, 7]
				list[0]++
			`,
			wantStackTop: value.SmallInt(6).ToValue(),
		},
		"subscript setter decrement": {
			source: `
				list := [5, 2, 7]
				list[0]--
			`,
			wantStackTop: value.SmallInt(4).ToValue(),
		},
		"subscript setter add": {
			source: `
				list := [5, 2, 7]
				list[0] += 8
			`,
			wantStackTop: value.SmallInt(13).ToValue(),
		},
		"subscript setter subtract": {
			source: `
				list := [5, 2, 7]
				list[0] -= 8
			`,
			wantStackTop: value.SmallInt(-3).ToValue(),
		},
		"subscript setter multiply": {
			source: `
				list := [5, 2, 7]
				list[1] *= 3
			`,
			wantStackTop: value.SmallInt(6).ToValue(),
		},
		"subscript setter divide": {
			source: `
				list := [5, 8, 7]
				list[1] /= 2
			`,
			wantStackTop: value.SmallInt(4).ToValue(),
		},
		"subscript setter exponentiate": {
			source: `
				list := [5, 8, 7]
				list[1] **= 2
			`,
			wantStackTop: value.SmallInt(64).ToValue(),
		},
		"subscript setter modulo type error": {
			source: `
				list := [5, 8, 7]
				list[0] %= 2
			`,
			wantStackTop: value.SmallInt(1).ToValue(),
		},
		"subscript setter modulo": {
			source: `
				list := [5, 8, 7]
				list[0] %= 2
			`,
			wantStackTop: value.SmallInt(1).ToValue(),
		},
		"subscript setter left bitshift": {
			source: `
				list := [5, 8, 7]
				list[0] <<= 2
			`,
			wantStackTop: value.SmallInt(20).ToValue(),
		},
		"subscript setter logic left bitshift": {
			source: `
				list := [5i8, 8i8, 7i8]
				list[0] <<<= 2
			`,
			wantStackTop: value.Int8(20).ToValue(),
		},
		"subscript setter right bitshift": {
			source: `
				list := [10, 8, 7]
				list[0] >>= 2
			`,
			wantStackTop: value.SmallInt(2).ToValue(),
		},
		"subscript setter logic right bitshift": {
			source: `
				list := [10i8, 8i8, 7i8]
				list[0] >>>= 2
			`,
			wantStackTop: value.Int8(2).ToValue(),
		},
		"subscript setter bitwise and": {
			source: `
				list := [6, 8, 7]
				list[0] &= 5
			`,
			wantStackTop: value.SmallInt(4).ToValue(),
		},
		"subscript setter bitwise or": {
			source: `
				list := [6, 8, 7]
				list[0] |= 5
			`,
			wantStackTop: value.SmallInt(7).ToValue(),
		},
		"subscript setter bitwise xor": {
			source: `
				list := [6, 8, 7]
				list[0] ^= 5
			`,
			wantStackTop: value.SmallInt(3).ToValue(),
		},
		"subscript setter logic or falsy": {
			source: `
				list := [nil, 8, 7.8]
				list[0] ||= 5
			`,
			wantStackTop: value.SmallInt(5).ToValue(),
		},
		"subscript setter logic or truthy": {
			source: `
				list := [1, 8, 7.8]
				list[0] ||= 5
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(29, 3, 5), P(35, 3, 11)), "this condition will always have the same result since type `Std::Int | Std::Float` is truthy"),
				diagnostic.NewWarning(L(P(41, 3, 17), P(41, 3, 17)), "unreachable code"),
			},
			wantStackTop: value.SmallInt(1).ToValue(),
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
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(29, 3, 5), P(35, 3, 11)), "this condition will always have the same result since type `Std::Int | Std::Float` is truthy"),
			},
			wantStackTop: value.SmallInt(5).ToValue(),
		},
		"subscript setter nil coalesce nil": {
			source: `
				list := [nil, 8, 7.8]
				list[0] ??= 5
			`,
			wantStackTop: value.SmallInt(5).ToValue(),
		},
		"subscript setter nil coalesce false": {
			source: `
				list := [false, 8, 7.8]
				list[0] ??= 5
			`,
			wantStackTop: value.False,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(33, 3, 5), P(39, 3, 11)), "this condition will always have the same result since type `bool | Std::Int | Std::Float` can never be nil"),
				diagnostic.NewWarning(L(P(45, 3, 17), P(45, 3, 17)), "unreachable code"),
			},
		},
		"subscript setter nil coalesce truthy": {
			source: `
				list := [1, 8, 7.8]
				list[0] ??= 5
			`,
			wantStackTop: value.SmallInt(1).ToValue(),
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(29, 3, 5), P(35, 3, 11)), "this condition will always have the same result since type `Std::Int | Std::Float` can never be nil"),
				diagnostic.NewWarning(L(P(41, 3, 17), P(41, 3, 17)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
