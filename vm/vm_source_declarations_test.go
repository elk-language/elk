package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
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
			wantStackTop: value.ToSymbol("boo").ToValue(),
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
			wantStackTop: value.ToSymbol("boo").ToValue(),
		},
		"define singleton methods on a module": {
			source: `
				module Foo
					singleton
						def bar then :boo
					end
				end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(21, 3, 6), P(62, 5, 8)), "singleton definitions cannot appear in this context"),
				diagnostic.NewFailure(L(P(37, 4, 7), P(53, 4, 23)), "method definitions cannot appear in this context"),
			},
		},
		"define singleton methods on an interface": {
			source: `
				interface Foo
					singleton
						def bar then :boo
					end
				end

				::Foo.bar
			`,
			wantStackTop: value.ToSymbol("boo").ToValue(),
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
			source:       "mixin Foo; end",
			wantStackTop: value.Nil,
		},
		"mixin without a body with an absolute name": {
			source:       "mixin ::Foo; end",
			wantStackTop: value.Nil,
		},
		"mixin with a body": {
			source: `
				mixin Foo
					a := 5
					println a
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "5\n",
		},
		"nested mixins": {
			source: `
				println Gdańsk::Gdynia::Sopot::Trójmiasto

				mixin Gdańsk
					mixin Gdynia
						mixin Sopot
							const Trójmiasto = "jest super"
						end
					end
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "jest super\n",
		},
		"open an existing mixin": {
			source: `
				println Foo::FIRST_CONSTANT
				println Foo::SECOND_CONSTANT

				mixin Foo
					const FIRST_CONSTANT = "oguem"
				end

				mixin Foo
					const SECOND_CONSTANT = "całe te"
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "oguem\ncałe te\n",
		},
		"redefined constant": {
			source: `
				const Foo = 3
				mixin Foo; end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 2, 5), P(17, 2, 17)), "cannot redeclare constant `Foo`"),
			},
		},
		"redefined class as mixin": {
			source: `
				class Foo; end
				mixin Foo; end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(24, 3, 5), P(37, 3, 18)), "cannot redeclare constant `Foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_DefineInterface(t *testing.T) {
	tests := sourceTestTable{
		"without a body with a relative name": {
			source:       "interface Foo; end",
			wantStackTop: value.Nil,
		},
		"without a body with an absolute name": {
			source:       "interface ::Foo; end",
			wantStackTop: value.Nil,
		},
		"with a body": {
			source: `
				interface Foo
					a := 5
					println a
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "5\n",
		},
		"nested": {
			source: `
				println Gdańsk::Gdynia::Sopot::Trójmiasto

				interface Gdańsk
					interface Gdynia
						interface Sopot
							const Trójmiasto = "jest super"
						end
					end
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "jest super\n",
		},
		"open an existing interface": {
			source: `
				println Foo::FIRST_CONSTANT
				println Foo::SECOND_CONSTANT

				interface Foo
					const FIRST_CONSTANT = "oguem"
				end

				interface Foo
					const SECOND_CONSTANT = "całe te"
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "oguem\ncałe te\n",
		},
		"redefined constant": {
			source: `
				const Foo = 3
				interface Foo; end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 2, 5), P(17, 2, 17)), "cannot redeclare constant `Foo`"),
			},
		},
		"redefined class as interface": {
			source: `
				class Foo; end
				interface Foo; end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(24, 3, 5), P(41, 3, 22)), "cannot redeclare constant `Foo`"),
			},
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

				class ::Std::Object < Value
					include ::Foo
				end

				self.foo
			`,
			wantStackTop: value.Ref(value.String("hey, it's foo")),
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

				class ::Std::Object < Value
					include ::Foo, ::Bar
				end

				self.foo + "; " + self.bar
			`,
			wantStackTop: value.Ref(value.String("hey, it's foo; hey, it's bar")),
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

				sealed primitive noinit class ::Std::Int < Value
					include ::Bar
				end

				1.foo + "; " + 1.bar
			`,
			wantStackTop: value.Ref(value.String("hey, it's foo; hey, it's bar")),
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
			source:       "class Foo; end",
			wantStackTop: value.Nil,
		},
		"class without a body with an absolute name": {
			source:       "class ::Foo; end",
			wantStackTop: value.Nil,
		},
		"class without a body with a parent": {
			source: `
				class Foo < ::Std::Error; end
				println Foo.superclass?.name
			`,
			wantStackTop: value.Nil,
			wantStdout:   "Std::Error\n",
		},
		"class with a body": {
			source: `
				class Foo
					a := 5
					println a
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "5\n",
		},
		"nested classes": {
			source: `
				println Gdańsk::Gdynia::Sopot::Trójmiasto

				class Gdańsk
					class Gdynia
						class Sopot
							const Trójmiasto = "jest super"
						end
					end
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "jest super\n",
		},
		"open an existing class": {
			source: `
				println Foo::FIRST_CONSTANT
				println Foo::SECOND_CONSTANT

				class Foo
					const FIRST_CONSTANT = "oguem"
				end

				class Foo
					const SECOND_CONSTANT = "całe te"
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "oguem\ncałe te\n",
		},
		"superclass mismatch": {
			source: `
				class Foo; end

				class Bar < ::Foo
					const FIRST_CONSTANT = "oguem"
				end

				class Bar < ::Std::Error
					const SECOND_CONSTANT = "całe te"
				end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(104, 8, 17), P(115, 8, 28)), "superclass mismatch in `Bar`, got `Std::Error`, expected `Foo`"),
			},
		},
		"incorrect superclass": {
			source: `
				const A = 3
				class Foo < ::A; end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(35, 3, 19), P(35, 3, 19)), "undefined type `A`"),
			},
		},
		"redefined constant": {
			source: `
				const Foo = 3
				class Foo; end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 2, 5), P(17, 2, 17)), "cannot redeclare constant `Foo`"),
			},
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
			source:       "module Foo; end",
			wantStackTop: value.Nil,
		},
		"module without a body with an absolute name": {
			source:       "module ::Foo; end",
			wantStackTop: value.Nil,
		},
		"module with a body": {
			source: `
				module Foo
					const B = 2
					a := 5
					println a + B
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "7\n",
		},
		"nested modules": {
			source: `
				println Gdańsk::Gdynia::Sopot::Trójmiasto

				module Gdańsk
					module Gdynia
						module Sopot
							const Trójmiasto = "jest super"
						end
					end
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "jest super\n",
		},
		"open an existing module": {
			source: `
				println Foo::FIRST_CONSTANT
				println Foo::SECOND_CONSTANT

				module Foo
					const FIRST_CONSTANT = "oguem"
				end

				module Foo
					const SECOND_CONSTANT = "całe te"
				end
				nil
			`,
			wantStackTop: value.Nil,
			wantStdout:   "oguem\ncałe te\n",
		},
		"redefined constant": {
			source: `
				const Foo = 3
				module Foo; end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 2, 5), P(17, 2, 17)), "cannot redeclare constant `Foo`"),
			},
		},
		"redefined class as module": {
			source: `
				class Foo; end
				module Foo; end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(24, 3, 5), P(38, 3, 19)), "cannot redeclare constant `Foo`"),
			},
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
		"::Std::Float::INF": value.FloatInf().ToValue(),
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
			source: `
				const Foo = 3i64
				Foo
			`,
			wantStackTop: value.Int64(3).ToValue(),
		},
		"Set constant under nested modules": {
			source: `
				module ::Std
					sealed primitive noinit class Int < Value
						const Foo = 3i64
					end
				end

				Int::Foo
			`,
			wantStackTop: value.Int64(3).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
