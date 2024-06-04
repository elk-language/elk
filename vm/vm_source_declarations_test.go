package vm_test

import (
	"testing"

	"github.com/elk-language/elk/value"
)

func TestVMSource_DefineSingleton(t *testing.T) {
	tests := sourceTestTable{
		"define singleton methods on a class": {
			source: `
				class Foo
					singleton
						def bar then :boo
					end
				end

				::Foo.bar
			`,
			wantStackTop: value.ToSymbol("boo"),
			teardown:     func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"define singleton methods on a mixin": {
			source: `
				mixin Foo
					singleton
						def bar then :boo
					end
				end

				::Foo.bar
			`,
			wantStackTop: value.ToSymbol("boo"),
			teardown:     func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"define singleton methods on a module": {
			source: `
				module Foo
					singleton
						def bar then :boo
					end
				end

				::Foo.bar
			`,
			wantStackTop: value.ToSymbol("boo"),
			teardown:     func() { value.RootModule.Constants.DeleteString("Foo") },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_DefineMixin(t *testing.T) {
	tests := sourceTestTable{
		"mixin without a body with a relative name": {
			source: "mixin Foo; end",
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"mixin without a body with an absolute name": {
			source: "mixin ::Foo; end",
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"mixin with a body": {
			source: `
				mixin Foo
					a := 5
					Bar := a - 2
				end
			`,
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithConstants(
					value.SymbolMap{
						value.ToSymbol("Bar"): value.SmallInt(3),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"nested mixins": {
			source: `
				mixin Gdańsk
					mixin Gdynia
						mixin Sopot
							Trójmiasto := "jest super"
							::Gdańsk::Warszawa := "to stolica"
						end
					end
				end
			`,
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Gdańsk"),
				value.MixinWithConstants(
					value.SymbolMap{
						value.ToSymbol("Gdynia"): value.NewMixinWithOptions(
							value.MixinWithName("Gdańsk::Gdynia"),
							value.MixinWithConstants(
								value.SymbolMap{
									value.ToSymbol("Sopot"): value.NewMixinWithOptions(
										value.MixinWithName("Gdańsk::Gdynia::Sopot"),
										value.MixinWithConstants(
											value.SymbolMap{
												value.ToSymbol("Trójmiasto"): value.String("jest super"),
											},
										),
									),
								},
							),
						),
						value.ToSymbol("Warszawa"): value.String("to stolica"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Gdańsk") },
		},
		"open an existing mixin": {
			source: `
				mixin Foo
					FIRST_CONSTANT := "oguem"
				end

				mixin Foo
					SECOND_CONSTANT := "całe te"
				end
			`,
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithConstants(
					value.SymbolMap{
						value.ToSymbol("FIRST_CONSTANT"):  value.String("oguem"),
						value.ToSymbol("SECOND_CONSTANT"): value.String("całe te"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"redefined constant": {
			source: `
				Foo := 3
				mixin Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.RedefinedConstantErrorClass,
				"module Root already has a constant named `:Foo`",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"redefined class as mixin": {
			source: `
				class Foo; end
				mixin Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.RedefinedConstantErrorClass,
				"module Root already has a constant named `:Foo`",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Include(t *testing.T) {
	tests := sourceTestTable{
		"include a mixin to a class": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				class ::Std::Object
					include ::Foo
				end

				self.foo
			`,
			wantStackTop: value.String("hey, it's foo"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.ObjectClass.Parent = value.ValueClass
			},
		},
		"include two mixins to a class": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				mixin Bar
					def bar: String
						"hey, it's bar"
					end
				end

				class ::Std::Object
					include ::Foo, ::Bar
				end

				self.foo + "; " + self.bar
			`,
			wantStackTop: value.String("hey, it's foo; hey, it's bar"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("Bar")
				value.ObjectClass.Parent = value.ValueClass
			},
		},
		"include a complex mixin in a class": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				mixin Bar
					include ::Foo

					def bar: String
						"hey, it's bar"
					end
				end

				sealed class ::Std::Int
					include ::Bar
				end

				1.foo + "; " + 1.bar
			`,
			wantStackTop: value.String("hey, it's foo; hey, it's bar"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("Bar")
				value.ObjectClass.Parent = value.ValueClass
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Extend(t *testing.T) {
	tests := sourceTestTable{
		"extend a class with a mixin": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				sealed class ::Std::String
					extend ::Foo
				end

				::Std::String.foo
			`,
			wantStackTop: value.String("hey, it's foo"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.StringClass.SetDirectClass(value.ObjectClass)
			},
		},
		"extend a module with a mixin": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				module Std
					extend ::Foo
				end

				::Std.foo
			`,
			wantStackTop: value.String("hey, it's foo"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.StdModule.SetDirectClass(value.ModuleClass)
			},
		},
		"extend a mixin with a mixin": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				mixin Bar
					extend ::Foo
				end

				::Bar.foo
			`,
			wantStackTop: value.String("hey, it's foo"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("Bar")
			},
		},
		"extend a class with two mixins": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				mixin Bar
					def bar: String
						"hey, it's bar"
					end
				end

				sealed class ::Std::String
					extend ::Foo, ::Bar
				end

				::Std::String.foo + "; " + ::Std::String.bar
			`,
			wantStackTop: value.String("hey, it's foo; hey, it's bar"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("Bar")
				value.StringClass.SetDirectClass(value.ObjectClass)
			},
		},
		"extend a class with a complex mixin": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				mixin Bar
					include ::Foo

					def bar: String
						"hey, it's bar"
					end
				end

				class Baz
					extend ::Bar
				end

				::Baz.foo + "; " + ::Baz.bar
			`,
			wantStackTop: value.String("hey, it's foo; hey, it's bar"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("Bar")
				value.RootModule.Constants.DeleteString("Baz")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_DefineClass(t *testing.T) {
	tests := sourceTestTable{
		"class without a body with a relative name": {
			source: "class Foo; end",
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"abstract class": {
			source: "abstract class Foo; end",
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithAbstract(),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"reopen a class with the abstract modifier": {
			source: `
				class Foo; end
				abstract class Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.ModifierMismatchErrorClass,
				"class Foo < Std::Object should be reopened without the `abstract` modifier",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"reopen an abstract class": {
			source: `
				abstract class Foo; end
				abstract class Foo; end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithAbstract(),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"reopen an abstract class without the modifier": {
			source: `
				abstract class Foo; end
				class Foo; end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithAbstract(),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"reopen an abstract class with the sealed modifier": {
			source: `
				abstract class Foo; end
				sealed class Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.ModifierMismatchErrorClass,
				"abstract class Foo < Std::Object should be reopened without the `sealed` modifier",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"sealed class": {
			source: "sealed class Foo; end",
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithSealed(),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"reopen a class with the sealed modifier": {
			source: `
				class Foo; end
				sealed class Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.ModifierMismatchErrorClass,
				"class Foo < Std::Object should be reopened without the `sealed` modifier",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"reopen a sealed class": {
			source: `
				sealed class Foo; end
				sealed class Foo; end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithSealed(),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"reopen a sealed class without the modifier": {
			source: `
				sealed class Foo; end
				class Foo; end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithSealed(),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"reopen a sealed class with the abstract modifier": {
			source: `
				sealed class Foo; end
				abstract class Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.ModifierMismatchErrorClass,
				"sealed class Foo < Std::Object should be reopened without the `abstract` modifier",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"class without a body with an absolute name": {
			source: "class ::Foo; end",
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"class without a body with a parent": {
			source: "class Foo < ::Std::Error; end",
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithParent(value.ErrorClass),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"inherit from a sealed class": {
			source: `
				sealed class Foo; end
				class Bar < ::Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.SealedClassErrorClass,
				"Bar cannot inherit from sealed class Foo < Std::Object",
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("Bar")
			},
		},
		"class with a body": {
			source: `
				class Foo
					a := 5
					Bar := a - 2
				end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithConstants(
					value.SymbolMap{
						value.ToSymbol("Bar"): value.SmallInt(3),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"nested classes": {
			source: `
				class Gdańsk
					class Gdynia
						class Sopot
							Trójmiasto := "jest super"
							::Gdańsk::Warszawa := "to stolica"
						end
					end
				end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Gdańsk"),
				value.ClassWithConstants(
					value.SymbolMap{
						value.ToSymbol("Gdynia"): value.NewClassWithOptions(
							value.ClassWithName("Gdańsk::Gdynia"),
							value.ClassWithConstants(
								value.SymbolMap{
									value.ToSymbol("Sopot"): value.NewClassWithOptions(
										value.ClassWithName("Gdańsk::Gdynia::Sopot"),
										value.ClassWithConstants(
											value.SymbolMap{
												value.ToSymbol("Trójmiasto"): value.String("jest super"),
											},
										),
									),
								},
							),
						),
						value.ToSymbol("Warszawa"): value.String("to stolica"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Gdańsk") },
		},
		"open an existing class": {
			source: `
				class Foo
					FIRST_CONSTANT := "oguem"
				end

				class Foo
					SECOND_CONSTANT := "całe te"
				end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithConstants(
					value.SymbolMap{
						value.ToSymbol("FIRST_CONSTANT"):  value.String("oguem"),
						value.ToSymbol("SECOND_CONSTANT"): value.String("całe te"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"superclass mismatch": {
			source: `
				class Foo; end

				class Bar < ::Foo
					FIRST_CONSTANT := "oguem"
				end

				class Bar < ::Std::Error
					SECOND_CONSTANT := "całe te"
				end
			`,
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"superclass mismatch in Bar, expected: Foo, got: Std::Error",
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("Bar")
			},
		},
		"incorrect superclass": {
			source: `
				A := 3
				class Foo < ::A; end
			`,
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`3` cannot be used as a superclass",
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("A")
			},
		},
		"redefined constant": {
			source: `
				Foo := 3
				class Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.RedefinedConstantErrorClass,
				"module Root already has a constant named `:Foo`",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_DefineModule(t *testing.T) {
	tests := sourceTestTable{
		"module without a body with a relative name": {
			source: "module Foo; end",
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"module without a body with an absolute name": {
			source: "module ::Foo; end",
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"module with a body": {
			source: `
				module Foo
					a := 5
					Bar := a - 2
				end
			`,
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithName("&Foo"),
						value.ClassWithParent(value.ModuleClass),
					),
				),
				value.ModuleWithName("Foo"),
				value.ModuleWithConstants(
					value.SymbolMap{
						value.ToSymbol("Bar"): value.SmallInt(3),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"nested modules": {
			source: `
				module Gdańsk
					module Gdynia
						module Sopot
							Trójmiasto := "jest super"
							::Gdańsk::Warszawa := "to stolica"
						end
					end
				end
			`,
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithName("Gdańsk"),
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithName("&Gdańsk"),
						value.ClassWithParent(value.ModuleClass),
					),
				),
				value.ModuleWithConstants(
					value.SymbolMap{
						value.ToSymbol("Gdynia"): value.NewModuleWithOptions(
							value.ModuleWithName("Gdańsk::Gdynia"),
							value.ModuleWithClass(
								value.NewClassWithOptions(
									value.ClassWithSingleton(),
									value.ClassWithName("&Gdańsk::Gdynia"),
									value.ClassWithParent(value.ModuleClass),
								),
							),
							value.ModuleWithConstants(
								value.SymbolMap{
									value.ToSymbol("Sopot"): value.NewModuleWithOptions(
										value.ModuleWithName("Gdańsk::Gdynia::Sopot"),
										value.ModuleWithClass(
											value.NewClassWithOptions(
												value.ClassWithSingleton(),
												value.ClassWithName("&Gdańsk::Gdynia::Sopot"),
												value.ClassWithParent(value.ModuleClass),
											),
										),
										value.ModuleWithConstants(
											value.SymbolMap{
												value.ToSymbol("Trójmiasto"): value.String("jest super"),
											},
										),
									),
								},
							),
						),
						value.ToSymbol("Warszawa"): value.String("to stolica"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Gdańsk") },
		},
		"open an existing module": {
			source: `
				module Foo
					FIRST_CONSTANT := "oguem"
				end

				module Foo
					SECOND_CONSTANT := "całe te"
				end
			`,
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithName("&Foo"),
						value.ClassWithParent(value.ModuleClass),
					),
				),
				value.ModuleWithName("Foo"),
				value.ModuleWithConstants(
					value.SymbolMap{
						value.ToSymbol("FIRST_CONSTANT"):  value.String("oguem"),
						value.ToSymbol("SECOND_CONSTANT"): value.String("całe te"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"redefined constant": {
			source: `
				Foo := 3
				module Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.RedefinedConstantErrorClass,
				"module Root already has a constant named `:Foo`",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"redefined class as module": {
			source: `
				class Foo; end
				module Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.RedefinedConstantErrorClass,
				"module Root already has a constant named `:Foo`",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_GetModuleConstant(t *testing.T) {
	tests := simpleSourceTestTable{
		"::Std":                     value.StdModule,
		"::Std::Int":                value.IntClass,
		"::Std::Float::INF":         value.FloatInf(),
		"a := ::Std::Float; a::INF": value.FloatInf(),
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}

func TestVMSource_DefineModuleConstant(t *testing.T) {
	tests := sourceTestTable{
		"Set constant under Root": {
			source:       "::Foo := 3i64",
			wantStackTop: value.Int64(3),
			teardown:     func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"Set constant under Root and read it": {
			source: `
				::Foo := 3i64
				::Foo
			`,
			wantStackTop: value.Int64(3),
			teardown:     func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"Set constant under nested modules": {
			source:       `::Std::Int::Foo := 3i64`,
			wantStackTop: value.Int64(3),
			teardown:     func() { value.IntClass.Constants.DeleteString("Foo") },
		},
		"Set constant under a variable": {
			source: `
				a := ::Std::Int
				a::Bar := "baz"
			`,
			wantStackTop: value.String("baz"),
			teardown:     func() { value.IntClass.Constants.DeleteString("Bar") },
		},
		"Set constant under a variable and read it": {
			source: `
				a := ::Std::Int
				a::Bar := "baz"
				::Std::Int::Bar
			`,
			wantStackTop: value.String("baz"),
			teardown:     func() { value.IntClass.Constants.DeleteString("Bar") },
		},
		"Set a constant under Int": {
			source: `
				a := 3
				a::Foo := 10
			`,
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`3` is not a module",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
