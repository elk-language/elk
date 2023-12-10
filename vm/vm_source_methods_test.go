package vm_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestVMSource_Alias(t *testing.T) {
	tests := sourceTestTable{
		"add an alias to a builtin method in Std::Int": {
			source: `
				class ::Std::Int
					alias add +
				end

				3.add(4)
			`,
			wantStackTop: value.SmallInt(7),
			teardown:     func() { delete(value.IntClass.Methods, value.ToSymbol("add")) },
		},
		"add an alias to a builtin method": {
			source: `
				alias klass class
			`,
			wantStackTop: vm.NewNativeMethodWithOptions(
				vm.NativeMethodWithStringName("class"),
			),
			teardown: func() { delete(value.GlobalObject.SingletonClass().Methods, value.ToSymbol("klass")) },
		},
		"add an alias to a nonexistent method": {
			source: `
				alias foo blabla
			`,
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"can't create an alias for a nonexistent method: blabla",
			),
		},
		"add an alias overriding a frozen method": {
			source: `
				class ::Std::Int
				  alias + class
				end
			`,
			wantRuntimeErr: value.NewError(
				value.FrozenMethodErrorClass,
				"can't override a frozen method: +",
			),
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
			`,
			wantStackTop: vm.NewBytecodeMethod(
				value.SymbolTable.Add("foo"),
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(5, 2, 5), P(29, 2, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 2),
				},
				nil,
				0,
				-1,
				false,
				false,
				[]value.Value{
					value.SymbolTable.Add("bar"),
				},
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"define a method with positional arguments in top level": {
			source: `
				def foo(a: Int, b: Int): Int
					c := 5
					a + b + c
				end
			`,
			wantStackTop: vm.NewBytecodeMethod(
				value.SymbolTable.Add("foo"),
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.ADD),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(5, 2, 5), P(67, 5, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Symbol{
					value.SymbolTable.Add("a"),
					value.SymbolTable.Add("b"),
				},
				0,
				-1,
				false,
				false,
				[]value.Value{
					value.SmallInt(5),
				},
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("bar"))
			},
		},
		"define a method with positional arguments in a class": {
			source: `
				class Bar
					def foo(a: Int, b: Int): Int
						c := 5
						a + b + c
					end
				end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Bar"),
				value.ClassWithMethods(
					value.MethodMap{
						value.SymbolTable.Add("foo"): vm.NewBytecodeMethod(
							value.SymbolTable.Add("foo"),
							[]byte{
								byte(bytecode.PREP_LOCALS8), 1,
								byte(bytecode.LOAD_VALUE8), 0,
								byte(bytecode.SET_LOCAL8), 3,
								byte(bytecode.POP),
								byte(bytecode.GET_LOCAL8), 1,
								byte(bytecode.GET_LOCAL8), 2,
								byte(bytecode.ADD),
								byte(bytecode.GET_LOCAL8), 3,
								byte(bytecode.ADD),
								byte(bytecode.RETURN),
							},
							L(P(20, 3, 6), P(85, 6, 8)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(4, 4),
								bytecode.NewLineInfo(5, 5),
								bytecode.NewLineInfo(6, 1),
							},
							[]value.Symbol{
								value.SymbolTable.Add("a"),
								value.SymbolTable.Add("b"),
							},
							0,
							-1,
							false,
							false,
							[]value.Value{
								value.SmallInt(5),
							},
						),
					},
				),
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Bar")
			},
		},
		"define a method with positional arguments in a module": {
			source: `
				module Bar
					def foo(a: Int, b: Int): Int
						c := 5
						a + b + c
					end
				end
			`,
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithName("Bar"),
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithParent(value.ModuleClass),
						value.ClassWithMethods(
							value.MethodMap{
								value.SymbolTable.Add("foo"): vm.NewBytecodeMethod(
									value.SymbolTable.Add("foo"),
									[]byte{
										byte(bytecode.PREP_LOCALS8), 1,
										byte(bytecode.LOAD_VALUE8), 0,
										byte(bytecode.SET_LOCAL8), 3,
										byte(bytecode.POP),
										byte(bytecode.GET_LOCAL8), 1,
										byte(bytecode.GET_LOCAL8), 2,
										byte(bytecode.ADD),
										byte(bytecode.GET_LOCAL8), 3,
										byte(bytecode.ADD),
										byte(bytecode.RETURN),
									},
									L(P(21, 3, 6), P(86, 6, 8)),
									bytecode.LineInfoList{
										bytecode.NewLineInfo(4, 4),
										bytecode.NewLineInfo(5, 5),
										bytecode.NewLineInfo(6, 1),
									},
									[]value.Symbol{
										value.SymbolTable.Add("a"),
										value.SymbolTable.Add("b"),
									},
									0,
									-1,
									false,
									false,
									[]value.Value{
										value.SmallInt(5),
									},
								),
							},
						),
					),
				),
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Bar")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_OverrideFrozenMethod(t *testing.T) {
	tests := sourceTestTable{
		"override a frozen builtin method": {
			source: `
				class ::Std::String
				  def +(other)
						"lol"
					end
				end
			`,
			wantRuntimeErr: value.NewError(
				value.FrozenMethodErrorClass,
				"can't override a frozen method: +",
			),
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
		"call a global method without arguments": {
			source: `
				def foo: Symbol
					:bar
				end

				self.foo
			`,
			wantStackTop: value.SymbolTable.Add("bar"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a global method with positional arguments": {
			source: `
				def add(a: Int, b: Int): Int
					a + b
				end

				self.add(5, 9)
			`,
			wantStackTop: value.SmallInt(14),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("add"))
			},
		},
		"call a method with missing required arguments": {
			source: `
				def add(a: Int, b: Int): Int
					a + b
				end

				self.add(5)
			`,
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"wrong number of arguments, given: 1, expected: 2..2",
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("add"))
			},
		},
		"call a method without optional arguments": {
			source: `
				def add(a: Int, b: Int = 3, c: Float = 20.5): Int
					a + b + c
				end

				self.add(5)
			`,
			wantStackTop: value.Float(28.5),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("add"))
			},
		},
		"call a method with some optional arguments": {
			source: `
				def add(a: Int, b: Int = 3, c: Float = 20.5): Int
					a + b + c
				end

				self.add(5, 0)
			`,
			wantStackTop: value.Float(25.5),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("add"))
			},
		},
		"call a method with all optional arguments": {
			source: `
				def add(a: Int, b: Int = 3, c: Float = 20.5): Int
					a + b + c
				end

				self.add(3, 2, 3.5)
			`,
			wantStackTop: value.Float(8.5),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("add"))
			},
		},
		"call a method with only named arguments": {
			source: `
				def foo(a: String, b: String, c: String, d: String = "default d", e: String = "default e"): String
					"a: " + a + ", b: " + b + ", c: " + c + ", d: " + d + ", e: " + e
				end

				self.foo(b: "b", a: "a", c: "c", e: "e")
			`,
			wantStackTop: value.String("a: a, b: b, c: c, d: default d, e: e"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with all required arguments and named arguments": {
			source: `
				def foo(a: String, b: String, c: String, d: String = "default d", e: String = "default e"): String
					"a: " + a + ", b: " + b + ", c: " + c + ", d: " + d + ", e: " + e
				end

				self.foo("a", c: "c", b: "b")
			`,
			wantStackTop: value.String("a: a, b: b, c: c, d: default d, e: default e"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with optional arguments and named arguments": {
			source: `
				def foo(a: String, b: String, c: String = "default c", d: String = "default d", e: String = "default e"): String
					"a: " + a + ", b: " + b + ", c: " + c + ", d: " + d + ", e: " + e
				end

				self.foo("a", "b", "c", e: "e")
			`,
			wantStackTop: value.String("a: a, b: b, c: c, d: default d, e: e"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with rest parameters and no arguments": {
			source: `
				def foo(*b): String
					b.inspect
				end

				self.foo()
			`,
			wantStackTop: value.String("[]"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with rest parameters and arguments": {
			source: `
				def foo(*b): String
					b.inspect
				end

				self.foo(1, 2, 3)
			`,
			wantStackTop: value.String("[1, 2, 3]"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with rest parameters and required arguments": {
			source: `
				def foo(a, *b): String
					"a: " + a.inspect + ", b: " + b.inspect
				end

				self.foo(1)
			`,
			wantStackTop: value.String("a: 1, b: []"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with rest parameters and all arguments": {
			source: `
				def foo(a, *b): String
					"a: " + a.inspect + ", b: " + b.inspect
				end

				self.foo(1, 2, 3, 4)
			`,
			wantStackTop: value.String("a: 1, b: [2, 3, 4]"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with rest parameters and no optional arguments": {
			source: `
				def foo(a = 3, *b): String
					"a: " + a.inspect + ", b: " + b.inspect
				end

				self.foo()
			`,
			wantStackTop: value.String("a: 3, b: []"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with rest parameters and optional arguments": {
			source: `
				def foo(a = 3, *b): String
					"a: " + a.inspect + ", b: " + b.inspect
				end

				self.foo(1)
			`,
			wantStackTop: value.String("a: 1, b: []"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with rest parameters and all optional arguments": {
			source: `
				def foo(a = 3, *b): String
					"a: " + a.inspect + ", b: " + b.inspect
				end

				self.foo(1, 2, 3)
			`,
			wantStackTop: value.String("a: 1, b: [2, 3]"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with post parameters": {
			source: `
				def foo(*a, b): String
					"a: " + a.inspect + ", b: " + b.inspect
				end

				self.foo(1, 2, 3)
			`,
			wantStackTop: value.String("a: [1, 2], b: 3"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with multiple post arguments": {
			source: `
				def foo(*a, b, c): String
					"a: " + a.inspect + ", b: " + b.inspect + ", c: " + c.inspect
				end

				self.foo(1, 2, 3, 4)
			`,
			wantStackTop: value.String("a: [1, 2], b: 3, c: 4"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with pre and post arguments": {
			source: `
				def foo(a, b, *c, d, e): String
					"a: " + a.inspect + ", b: " + b.inspect + ", c: " + c.inspect + ", d: " + d.inspect + ", e: " + e.inspect
				end

				self.foo(1, 2, 3, 4, 5, 6)
			`,
			wantStackTop: value.String("a: 1, b: 2, c: [3, 4], d: 5, e: 6"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with named post arguments": {
			source: `
				def foo(*a, b, c): String
					"a: " + a.inspect + ", b: " + b.inspect + ", c: " + c.inspect
				end

				self.foo(1, 2, c: 3, b: 4)
			`,
			wantStackTop: value.String("a: [1, 2], b: 4, c: 3"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with pre and named post arguments": {
			source: `
				def foo(a, *b, c, d): String
					"a: " + a.inspect + ", b: " + b.inspect + ", c: " + c.inspect + ", d: " + d.inspect
				end

				self.foo(1, 2, 3, d: 4, c: 5)
			`,
			wantStackTop: value.String("a: 1, b: [2, 3], c: 5, d: 4"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with named arguments and missing required arguments": {
			source: `
				def foo(a: String, b: String, c: String = "default c", d: String = "default d", e: String = "default e"): String
					"a: " + a + ", b: " + b + ", c: " + c + ", d: " + d + ", e: " + e
				end

				self.foo("a", e: "e")
			`,
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"missing required argument `b` in call to `foo`",
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with duplicated arguments": {
			source: `
				def foo(a: String, b: String): String
					"a: " + a + ", b: " + b
				end

				self.foo("a", b: "b", a: "a2")
			`,
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"duplicated argument `a` in call to `foo`",
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a method with unknown named arguments": {
			source: `
				def foo(a: String, b: String): String
					"a: " + a + ", b: " + b
				end

				self.foo("a", unknown: "lala", moo: "meow", b: "b")
			`,
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"unknown arguments: [:unknown, :moo]",
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
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
			wantStackTop: value.SymbolTable.Add("baz"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
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
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
		},
		"call an instance method without arguments": {
			source: `
				class ::Std::Object
					def bar: Symbol
						:baz
					end
				end

				self.bar
			`,
			wantStackTop: value.SymbolTable.Add("baz"),
			teardown: func() {
				delete(value.ObjectClass.Methods, value.SymbolTable.Add("bar"))
			},
		},
		"call an instance method with positional arguments": {
			source: `
				class ::Std::Object
					def add(a: Int, b: Int): Int
						a + b
					end
				end

				self.add 1, 8
			`,
			wantStackTop: value.SmallInt(9),
			teardown: func() {
				delete(value.ObjectClass.Methods, value.SymbolTable.Add("add"))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
