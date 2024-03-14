package vm_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
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
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `[]` is not available to value of class `Std::Nil`: nil",
			),
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
				list := ["foo", 2, 7.8]
				list?[0]
			`,
			wantStackTop: value.String("foo"),
		},
		"get index -1 of a list": {
			source: `
				list := ["foo", 2, 7.8]
				list?[-1]
			`,
			wantStackTop: value.Float(7.8),
		},
		"get too big index": {
			source: `
				list := ["foo", 2, 7.8]
				list?[50]
			`,
			wantRuntimeErr: value.NewError(
				value.IndexErrorClass,
				"index 50 out of range: -3...3",
			),
		},
		"get too small index": {
			source: `
				list := ["foo", 2, 7.8]
				list?[-10]
			`,
			wantRuntimeErr: value.NewError(
				value.IndexErrorClass,
				"index -10 out of range: -3...3",
			),
		},
		"get from nil": {
			source: `
				list := nil
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
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
		},
		"instantiate a class without an initialiser with arguments": {
			source: `
				class Foo; end

				::Foo(2)
			`,
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"wrong number of arguments, given: 1, expected: 0",
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
		},
		"instantiate a class with an initialiser without arguments": {
			source: `
				class Foo
					init(a); end
				end

				::Foo()
			`,
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"wrong number of arguments, given: 0, expected: 1..1",
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
		},
		"instantiate a class with an initialiser with arguments": {
			source: `
				class Foo
					init(a)
						println("a: " + a)
					end
				end

				::Foo("bar")
			`,
			wantStdout: "a: bar\n",
			wantStackTop: value.NewObject(
				value.ObjectWithClass(
					value.NewClassWithOptions(
						value.ClassWithName("Foo"),
						value.ClassWithMethods(
							value.MethodMap{
								value.ToSymbol("#init"): vm.NewBytecodeMethod(
									value.ToSymbol("#init"),
									[]byte{
										byte(bytecode.LOAD_VALUE8), 0,
										byte(bytecode.GET_LOCAL8), 1,
										byte(bytecode.ADD),
										byte(bytecode.CALL_FUNCTION8), 1,
										byte(bytecode.POP),
										byte(bytecode.RETURN_SELF),
									},
									L(P(20, 3, 6), P(60, 5, 8)),
									bytecode.LineInfoList{
										bytecode.NewLineInfo(4, 4),
										bytecode.NewLineInfo(5, 2),
									},
									[]value.Symbol{
										value.ToSymbol("a"),
									},
									0,
									-1,
									false,
									false,
									[]value.Value{
										value.String("a: "),
										value.NewCallSiteInfo(
											value.ToSymbol("println"),
											1,
											nil,
										),
									},
								),
							},
						),
					),
				),
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
		},
		"instantiate a class with an initialiser with ivar parameters": {
			source: `
				class Foo
					init(@a)
						println("a: " + a)
					end
				end

				::Foo("bar")
			`,
			wantStdout: "a: bar\n",
			wantStackTop: value.NewObject(
				value.ObjectWithInstanceVariables(
					value.SymbolMap{
						value.ToSymbol("a"): value.String("bar"),
					},
				),
				value.ObjectWithClass(
					value.NewClassWithOptions(
						value.ClassWithName("Foo"),
						value.ClassWithMethods(
							value.MethodMap{
								value.ToSymbol("#init"): vm.NewBytecodeMethod(
									value.ToSymbol("#init"),
									[]byte{
										byte(bytecode.GET_LOCAL8), 1,
										byte(bytecode.SET_IVAR8), 0,
										byte(bytecode.POP),
										byte(bytecode.LOAD_VALUE8), 1,
										byte(bytecode.GET_LOCAL8), 1,
										byte(bytecode.ADD),
										byte(bytecode.CALL_FUNCTION8), 2,
										byte(bytecode.POP),
										byte(bytecode.RETURN_SELF),
									},
									L(P(20, 3, 6), P(61, 5, 8)),
									bytecode.LineInfoList{
										bytecode.NewLineInfo(3, 3),
										bytecode.NewLineInfo(4, 4),
										bytecode.NewLineInfo(5, 2),
									},
									[]value.Symbol{
										value.ToSymbol("a"),
									},
									0,
									-1,
									false,
									false,
									[]value.Value{
										value.ToSymbol("a"),
										value.String("a: "),
										value.NewCallSiteInfo(
											value.ToSymbol("println"),
											1,
											nil,
										),
									},
								),
							},
						),
					),
				),
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
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
				sealed class ::Std::Int
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
			wantStackTop: value.Nil,
			teardown:     func() { delete(value.GlobalObject.SingletonClass().Methods, value.ToSymbol("klass")) },
		},
		"add an alias to a nonexistent method": {
			source: `
				alias foo blabla
			`,
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"cannot create an alias for a nonexistent method: blabla",
			),
		},
		"add an alias overriding a sealed method": {
			source: `
				sealed class ::Std::Int
				  alias + class
				end
			`,
			wantRuntimeErr: value.NewError(
				value.SealedMethodErrorClass,
				"cannot override a sealed method: +",
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
				value.ToSymbol("foo"),
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
					value.ToSymbol("bar"),
				},
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
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
				value.ToSymbol("foo"),
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
					value.ToSymbol("a"),
					value.ToSymbol("b"),
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
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("bar"))
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
						value.ToSymbol("foo"): vm.NewBytecodeMethod(
							value.ToSymbol("foo"),
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
								value.ToSymbol("a"),
								value.ToSymbol("b"),
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
						value.ClassWithName("&Bar"),
						value.ClassWithMethods(
							value.MethodMap{
								value.ToSymbol("foo"): vm.NewBytecodeMethod(
									value.ToSymbol("foo"),
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
										value.ToSymbol("a"),
										value.ToSymbol("b"),
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

func TestVMSource_OverrideSealedMethod(t *testing.T) {
	tests := sourceTestTable{
		"override a sealed builtin method": {
			source: `
				sealed class ::Std::String
				  def +(other)
						"lol"
					end
				end
			`,
			wantRuntimeErr: value.NewError(
				value.SealedMethodErrorClass,
				"cannot override a sealed method: +",
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
		"nil safe call on nil": {
			source: `
				a := nil

				a?.foo(3, 4)?.bar
			`,
			wantStackTop: value.Nil,
		},
		"nil safe call on not nil": {
			source: `
				a := 5

				a?.inspect
			`,
			wantStackTop: value.String("5"),
		},
		"call a global method without arguments": {
			source: `
				def foo: Symbol
					:bar
				end

				self.foo
			`,
			wantStackTop: value.ToSymbol("bar"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
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
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("add"))
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
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("add"))
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
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("add"))
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
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("add"))
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
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("add"))
			},
		},
		"call a method with only named arguments": {
			source: `
				def foo(a: String, b: String, c: String, d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				self.foo(b: "b", a: "a", c: "c", e: "e")
			`,
			wantStackTop: value.String("a: a, b: b, c: c, d: default d, e: e"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with all required arguments and named arguments": {
			source: `
				def foo(a: String, b: String, c: String, d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				self.foo("a", c: "c", b: "b")
			`,
			wantStackTop: value.String("a: a, b: b, c: c, d: default d, e: default e"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with optional arguments and named arguments": {
			source: `
				def foo(a: String, b: String, c: String = "default c", d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				self.foo("a", "b", "c", e: "e")
			`,
			wantStackTop: value.String("a: a, b: b, c: c, d: default d, e: e"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with a named rest param and no args": {
			source: `
				def foo(**a: String): String
					"a: ${a.inspect}"
				end

				self.foo()
			`,
			wantStackTop: value.String("a: {}"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with a named rest param and a few named args": {
			source: `
				def foo(**a: String): HashMap[Symbol, String]
					a
				end

				self.foo(d: "foo", a: "bar")
			`,
			wantStackTop: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.ToSymbol("a"), Value: value.String("bar")},
				value.Pair{Key: value.ToSymbol("d"), Value: value.String("foo")},
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with regular params, named rest param and a few named args": {
			source: `
				def foo(a, **b: String): List
					[a, b]
				end

				self.foo("foo", c: "bar", d: "baz")
			`,
			wantStackTop: value.NewArrayListWithElements(
				2,
				value.String("foo"),
				vm.MustNewHashMapWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("c"), Value: value.String("bar")},
					value.Pair{Key: value.ToSymbol("d"), Value: value.String("baz")},
				),
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with regular params, named rest param and only required args": {
			source: `
				def foo(a, **b: String): List
					[a, b]
				end

				self.foo("foo")
			`,
			wantStackTop: value.NewArrayListWithElements(
				2,
				value.String("foo"),
				value.NewHashMap(0),
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with regular params, optional params, named rest param and a few named args": {
			source: `
				def foo(a, b = 5, **c: String): List
					[a, b, c]
				end

				self.foo("foo", c: "bar", d: "baz")
			`,
			wantStackTop: value.NewArrayListWithElements(
				3,
				value.String("foo"),
				value.SmallInt(5),
				vm.MustNewHashMapWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("c"), Value: value.String("bar")},
					value.Pair{Key: value.ToSymbol("d"), Value: value.String("baz")},
				),
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with regular params, optional params, named rest param and all args": {
			source: `
				def foo(a, b = 5, **c): List
					[a, b, c]
				end

				self.foo("foo", 9, c: "bar", d: "baz")
			`,
			wantStackTop: value.NewArrayListWithElements(
				3,
				value.String("foo"),
				value.SmallInt(9),
				vm.MustNewHashMapWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("c"), Value: value.String("bar")},
					value.Pair{Key: value.ToSymbol("d"), Value: value.String("baz")},
				),
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with regular params, optional params, named rest param and optional named arg": {
			source: `
				def foo(a, b = 5, **c): ArrayList[Value]
					[a, b, c]
				end

				self.foo("foo", c: "bar", d: "baz", b: 9)
			`,
			wantStackTop: value.NewArrayListWithElements(
				2,
				value.String("foo"),
				value.SmallInt(9),
				vm.MustNewHashMapWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("d"), Value: value.String("baz")},
					value.Pair{Key: value.ToSymbol("c"), Value: value.String("bar")},
				),
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with positional rest params and named rest params and no args": {
			source: `
				def foo(*a, **b): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				self.foo()
			`,
			wantStackTop: value.String(`a: [], b: {}`),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with positional rest params and named rest params and positional args": {
			source: `
				def foo(*a, **b): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				self.foo(1, 5, 7)
			`,
			wantStackTop: value.String(`a: [1, 5, 7], b: {}`),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with positional rest params and named rest params and named args": {
			source: `
				def foo(*a, **b): ArrayList[Value]
					[a, b]
				end

				self.foo(foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.NewArrayListWithElements(
				2,
				&value.NilArrayList,
				vm.MustNewHashMapWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo"), Value: value.SmallInt(5)},
					value.Pair{Key: value.ToSymbol("bar"), Value: value.SmallInt(2)},
					value.Pair{Key: value.ToSymbol("baz"), Value: value.SmallInt(8)},
				),
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with positional rest params and named rest params and both types of args": {
			source: `
				def foo(*a, **b): ArrayList[Value]
					[a, b]
				end

				self.foo(10, 20, 30, foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.NewArrayListWithElements(
				2,
				value.NewArrayListWithElements(
					3,
					value.SmallInt(10),
					value.SmallInt(20),
					value.SmallInt(30),
				),
				vm.MustNewHashMapWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo"), Value: value.SmallInt(5)},
					value.Pair{Key: value.ToSymbol("bar"), Value: value.SmallInt(2)},
					value.Pair{Key: value.ToSymbol("baz"), Value: value.SmallInt(8)},
				),
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},

		"call a method with regular, positional rest params and named rest params and no args": {
			source: `
				def foo(a, *b, **c): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				self.foo(5)
			`,
			wantStackTop: value.String(`a: 5, b: [], c: {}`),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with regular, positional rest params and named rest params and positional args": {
			source: `
				def foo(a, *b, **c): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				self.foo(1, 5, 7)
			`,
			wantStackTop: value.String(`a: 1, b: [5, 7], c: {}`),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with regular, positional rest params and named rest params and named args": {
			source: `
				def foo(a, *b, **c): ArrayList[Value]
					[a, b, c]
				end

				self.foo(1, foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.NewArrayListWithElements(
				3,
				value.SmallInt(1),
				&value.NilArrayList,
				vm.MustNewHashMapWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo"), Value: value.SmallInt(5)},
					value.Pair{Key: value.ToSymbol("bar"), Value: value.SmallInt(2)},
					value.Pair{Key: value.ToSymbol("baz"), Value: value.SmallInt(8)},
				),
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with regular, positional rest params and named rest params and both types of args": {
			source: `
				def foo(a, *b, **c): ArrayList[Value]
					[a, b, c]
				end

				self.foo(10, 20, 30, foo: 5, bar: 2, baz: 8)
			`,
			wantStackTop: value.NewArrayListWithElements(
				3,
				value.SmallInt(10),
				value.NewArrayListWithElements(
					2,
					value.SmallInt(20),
					value.SmallInt(30),
				),
				vm.MustNewHashMapWithElements(
					nil,
					value.Pair{Key: value.ToSymbol("foo"), Value: value.SmallInt(5)},
					value.Pair{Key: value.ToSymbol("bar"), Value: value.SmallInt(2)},
					value.Pair{Key: value.ToSymbol("baz"), Value: value.SmallInt(8)},
				),
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
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
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
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
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with rest parameters and required arguments": {
			source: `
				def foo(a, *b): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				self.foo(1)
			`,
			wantStackTop: value.String("a: 1, b: []"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with rest parameters and all arguments": {
			source: `
				def foo(a, *b): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				self.foo(1, 2, 3, 4)
			`,
			wantStackTop: value.String("a: 1, b: [2, 3, 4]"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with rest parameters and named args": {
			source: `
				def foo(a, b, *c): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				self.foo(b: 1, a: 2)
			`,
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"wrong number of positional arguments, given: 0, expected: 2..",
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with rest parameters and no optional arguments": {
			source: `
				def foo(a = 3, *b): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				self.foo()
			`,
			wantStackTop: value.String("a: 3, b: []"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with rest parameters and optional arguments": {
			source: `
				def foo(a = 3, *b): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				self.foo(1)
			`,
			wantStackTop: value.String("a: 1, b: []"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with rest parameters and all optional arguments": {
			source: `
				def foo(a = 3, *b): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				self.foo(1, 2, 3)
			`,
			wantStackTop: value.String("a: 1, b: [2, 3]"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with post parameters": {
			source: `
				def foo(*a, b): String
					"a: ${a.inspect}, b: ${b.inspect}"
				end

				self.foo(1, 2, 3)
			`,
			wantStackTop: value.String("a: [1, 2], b: 3"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with multiple post arguments": {
			source: `
				def foo(*a, b, c): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				self.foo(1, 2, 3, 4)
			`,
			wantStackTop: value.String("a: [1, 2], b: 3, c: 4"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with pre and post arguments": {
			source: `
				def foo(a, b, *c, d, e): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}, d: ${d.inspect}, e: ${e.inspect}"
				end

				self.foo(1, 2, 3, 4, 5, 6)
			`,
			wantStackTop: value.String("a: 1, b: 2, c: [3, 4], d: 5, e: 6"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with named post arguments": {
			source: `
				def foo(*a, b, c): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}"
				end

				self.foo(1, 2, c: 3, b: 4)
			`,
			wantStackTop: value.String("a: [1, 2], b: 4, c: 3"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with pre and named post arguments": {
			source: `
				def foo(a, *b, c, d): String
					"a: ${a.inspect}, b: ${b.inspect}, c: ${c.inspect}, d: ${d.inspect}"
				end

				self.foo(1, 2, 3, d: 4, c: 5)
			`,
			wantStackTop: value.String("a: 1, b: [2, 3], c: 5, d: 4"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with named arguments and missing required arguments": {
			source: `
				def foo(a: String, b: String, c: String = "default c", d: String = "default d", e: String = "default e"): String
					"a: ${a}, b: ${b}, c: ${c}, d: ${d}, e: ${e}"
				end

				self.foo("a", e: "e")
			`,
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"missing required argument `b` in call to `foo`",
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with duplicated arguments": {
			source: `
				def foo(a: String, b: String): String
					"a: ${a}, b: ${b}"
				end

				self.foo("a", b: "b", a: "a2")
			`,
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"duplicated argument `a` in call to `foo`",
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"call a method with unknown named arguments": {
			source: `
				def foo(a: String, b: String): String
					"a: ${a}, b: ${b}"
				end

				self.foo("a", unknown: "lala", moo: "meow", b: "b")
			`,
			wantRuntimeErr: value.NewError(
				value.ArgumentErrorClass,
				"unknown arguments: [:unknown, :moo]",
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
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
			wantStackTop: value.ToSymbol("baz"),
			teardown: func() {
				delete(value.ObjectClass.Methods, value.ToSymbol("bar"))
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
				delete(value.ObjectClass.Methods, value.ToSymbol("add"))
			},
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
				def foo=(v)
					:bar
				end

				self.foo = 3
			`,
			wantStackTop: value.SmallInt(3),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo="))
			},
		},
		"call subscript set": {
			source: `
				list := [5, 8, 20]
				list[0] = :foo
				list
			`,
			wantStackTop: &value.ArrayList{value.ToSymbol("foo"), value.SmallInt(8), value.SmallInt(20)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
