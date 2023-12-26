package compiler

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestSingletonBlock(t *testing.T) {
	tests := testTable{
		"define in top-level": {
			input: `
				singleton
					def foo then :bar
				end
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(5, 2, 5), P(44, 4, 7)), "cannot open a singleton class in the top level"),
			},
		},
		"define in a method": {
			input: `
				def baz
					singleton
						def foo then :bar
					end
				end
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(18, 3, 6), P(59, 5, 8)), "cannot open a singleton class in a method"),
			},
		},
		"define in a setter": {
			input: `
				def baz=(arg)
					singleton
						def foo then :bar
					end
				end
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(24, 3, 6), P(65, 5, 8)), "cannot open a singleton class in a method"),
			},
		},
		"define in a class": {
			input: `
				class Baz
					singleton
						def foo then :bar
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SELF),
							byte(bytecode.DEF_SINGLETON),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(69, 6, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							vm.NewBytecodeMethodNoParams(
								singletonClassSymbol,
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.DEF_METHOD),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(20, 3, 6), P(61, 5, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
									bytecode.NewLineInfo(5, 2),
								},
								[]value.Value{
									vm.NewBytecodeMethodNoParams(
										value.ToSymbol("foo"),
										[]byte{
											byte(bytecode.LOAD_VALUE8), 0,
											byte(bytecode.RETURN),
										},
										L(P(36, 4, 7), P(52, 4, 23)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 2),
										},
										[]value.Value{
											value.ToSymbol("bar"),
										},
									),
									value.ToSymbol("foo"),
								},
							),
						},
					),
					value.ToSymbol("Baz"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestDocComment(t *testing.T) {
	tests := testTable{
		"document a module": {
			input: `
				##[
					Foo is awesome
				]##
				module Foo; end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.DOC_COMMENT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(56, 5, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(5, 6),
				},
				[]value.Value{
					value.String("Foo is awesome"),
					value.ToSymbol("Foo"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestGetter(t *testing.T) {
	tests := testTable{
		"define single getter": {
			input: "getter foo",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_GETTER),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.ToSymbol("foo"),
				},
			),
		},
		"define three getters": {
			input: "getter foo: Foo, bar, baz: String",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_GETTER),

					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_GETTER),

					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DEF_GETTER),

					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 1, 33)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.ToSymbol("foo"),
					value.ToSymbol("bar"),
					value.ToSymbol("baz"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestSetter(t *testing.T) {
	tests := testTable{
		"define single setter": {
			input: "setter foo",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_SETTER),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.ToSymbol("foo"),
				},
			),
		},
		"define three setters": {
			input: "setter foo: Foo, bar, baz: String",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_SETTER),

					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_SETTER),

					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DEF_SETTER),

					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 1, 33)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.ToSymbol("foo"),
					value.ToSymbol("bar"),
					value.ToSymbol("baz"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestAccessor(t *testing.T) {
	tests := testTable{
		"define single accessor": {
			input: "accessor foo",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_GETTER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_SETTER),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.ToSymbol("foo"),
				},
			),
		},
		"define three accessors": {
			input: "accessor foo: Foo, bar, baz: String",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_GETTER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_SETTER),

					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_GETTER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_SETTER),

					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DEF_GETTER),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DEF_SETTER),

					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 1, 35)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 14),
				},
				[]value.Value{
					value.ToSymbol("foo"),
					value.ToSymbol("bar"),
					value.ToSymbol("baz"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestAlias(t *testing.T) {
	tests := testTable{
		"define single alias": {
			input: "alias foo bar",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_ALIAS),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.ToSymbol("bar"),
					value.ToSymbol("foo"),
				},
			),
		},
		"define three aliases": {
			input: "alias foo bar, remove delete, add plus",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_ALIAS),

					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.DEF_ALIAS),

					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.DEF_ALIAS),

					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(37, 1, 38)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{
					value.ToSymbol("bar"),
					value.ToSymbol("foo"),

					value.ToSymbol("delete"),
					value.ToSymbol("remove"),

					value.ToSymbol("plus"),
					value.ToSymbol("add"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestDefClass(t *testing.T) {
	tests := testTable{
		"anonymous class without a body": {
			input: "class; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				nil,
			),
		},
		"class with a relative name without a body": {
			input: "class Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
				},
			),
		},
		"abstract class": {
			input: "abstract class Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), byte(value.CLASS_ABSTRACT_FLAG),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 1, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
				},
			),
		},
		"named class inside of a method": {
			input: `
				def foo
				  class ::Bar; end
				end
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(19, 3, 7), P(34, 3, 22)), "cannot define named classes inside of a method: foo"),
			},
		},
		"anonymous class with an absolute parent": {
			input: "class < ::Bar; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.ToSymbol("Bar"),
				},
			),
		},
		"class with an absolute parent": {
			input: "class Foo < ::Bar; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 1, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
					value.ToSymbol("Bar"),
				},
			),
		},
		"anonymous class with a nested parent": {
			input: "class < ::Baz::Bar; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 1, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.ToSymbol("Baz"),
					value.ToSymbol("Bar"),
				},
			),
		},
		"class with an absolute nested parent": {
			input: "class Foo < ::Baz::Bar; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 1, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
					value.ToSymbol("Baz"),
					value.ToSymbol("Bar"),
				},
			),
		},
		"class with an absolute name without a body": {
			input: "class ::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
				},
			),
		},
		"class with an absolute nested name without a body": {
			input: "class ::Std::Int::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 1, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.ToSymbol("Std"),
					value.ToSymbol("Int"),
					value.ToSymbol("Foo"),
				},
			),
		},
		"class with a body": {
			input: `
				class Foo
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(44, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"class with a an error": {
			input: "class A then def a then a",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 1, 25)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(0, 1, 1), P(24, 1, 25)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
						},
						[]value.Value{
							vm.NewBytecodeMethodNoParams(
								value.ToSymbol("a"),
								[]byte{
									byte(bytecode.RETURN),
								},
								L(P(13, 1, 14), P(24, 1, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(1, 1),
								},
								nil,
							),
							value.ToSymbol("a"),
						},
					),
					value.ToSymbol("A"),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L(P(24, 1, 25), P(24, 1, 25)), "undeclared variable: a"),
			},
		},
		"anonymous class with a body": {
			input: `
				class
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(40, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
				},
			),
		},
		"nested classes": {
			input: `
				class Foo
					class Bar
						a := 1
						a + 2
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.UNDEFINED),
							byte(bytecode.DEF_CLASS), 0,
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(70, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							vm.NewBytecodeMethodNoParams(
								classSymbol,
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.SET_LOCAL8), 3,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 3,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(20, 3, 6), P(62, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 2),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
							value.ToSymbol("Bar"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestDefModule(t *testing.T) {
	tests := testTable{
		"anonymous module without a body": {
			input: "module; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				nil,
			),
		},
		"module with a relative name without a body": {
			input: "module Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
				},
			),
		},
		"named module inside of a method": {
			input: `
				def foo
					module Bar; end
				end
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(18, 3, 6), P(32, 3, 20)), "cannot define named modules inside of a method: foo"),
			},
		},
		"class with an absolute name without a body": {
			input: "module ::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
				},
			),
		},
		"class with an absolute nested name without a body": {
			input: "module ::Std::Int::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 1, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.ToSymbol("Std"),
					value.ToSymbol("Int"),
					value.ToSymbol("Foo"),
				},
			),
		},
		"anonymous module with a body": {
			input: `
				module
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_ANON_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(41, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
				},
			),
		},
		"module with a body": {
			input: `
				module Foo
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(45, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"nested modules": {
			input: `
				module Foo
					module Bar
						a := 1
						a + 2
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_MODULE),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(72, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							vm.NewBytecodeMethodNoParams(
								moduleSymbol,
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.SET_LOCAL8), 3,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 3,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(21, 3, 6), P(64, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 2),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
							value.ToSymbol("Bar"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestDefMethod(t *testing.T) {
	tests := testTable{
		"define method in top level": {
			input: `
				def foo then :bar
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 2, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						value.ToSymbol("foo"),
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(21, 2, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(2, 2),
						},
						[]value.Value{
							value.ToSymbol("bar"),
						},
					),
					value.ToSymbol("foo"),
				},
			),
		},
		"define a setter": {
			input: `
				def foo=(a)
					println(a + 2)
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(44, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethod(
						value.ToSymbol("foo="),
						[]byte{
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.ADD),
							byte(bytecode.CALL_FUNCTION8), 1,
							byte(bytecode.POP),
							byte(bytecode.RETURN_FIRST_ARG),
						},
						L(P(5, 2, 5), P(43, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Symbol{
							value.ToSymbol("a"),
						},
						0,
						-1,
						false,
						false,
						[]value.Value{
							value.SmallInt(2),
							value.NewCallSiteInfo(
								value.ToSymbol("println"),
								1,
								nil,
							),
						},
					),
					value.ToSymbol("foo="),
				},
			),
		},
		"define a setter with return": {
			input: `
				def foo=(a)
					println(a + 2)
					return "siema"
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethod(
						value.ToSymbol("foo="),
						[]byte{
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.ADD),
							byte(bytecode.CALL_FUNCTION8), 1,
							byte(bytecode.POP),
							byte(bytecode.RETURN_FIRST_ARG),
						},
						L(P(5, 2, 5), P(63, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
							bytecode.NewLineInfo(4, 1),
						},
						[]value.Symbol{
							value.ToSymbol("a"),
						},
						0,
						-1,
						false,
						false,
						[]value.Value{
							value.SmallInt(2),
							value.NewCallSiteInfo(
								value.ToSymbol("println"),
								1,
								nil,
							),
						},
					),
					value.ToSymbol("foo="),
				},
			),
		},
		"define method with required parameters in top level": {
			input: `
				def foo(a, b)
					c := 5
					a + b + c
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethod(
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
						L(P(5, 2, 5), P(52, 5, 7)),
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
					value.ToSymbol("foo"),
				},
			),
		},
		"define method with ivar parameters": {
			input: `
				def foo(@a, @b)
					c := 5
					a + b + c
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethod(
						value.ToSymbol("foo"),
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.SET_IVAR8), 0,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 2,
							byte(bytecode.SET_IVAR8), 1,
							byte(bytecode.POP),
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.GET_LOCAL8), 2,
							byte(bytecode.ADD),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.ADD),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(54, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(2, 7),
							bytecode.NewLineInfo(3, 3),
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
							value.ToSymbol("a"),
							value.ToSymbol("b"),
							value.SmallInt(5),
						},
					),
					value.ToSymbol("foo"),
				},
			),
		},
		"define method with default ivar parameters": {
			input: `
				def foo(@a = 5, @b = "b")
					c := 5
					a + b + c
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(65, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethod(
						value.ToSymbol("foo"),
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.JUMP_UNLESS_UNDEF), 0, 5,
							byte(bytecode.POP),
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 1,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.SET_IVAR8), 1,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 2,
							byte(bytecode.JUMP_UNLESS_UNDEF), 0, 5,
							byte(bytecode.POP),
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.SET_LOCAL8), 2,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 2,
							byte(bytecode.SET_IVAR8), 3,
							byte(bytecode.POP),
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
						L(P(5, 2, 5), P(64, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(2, 19),
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(4, 5),
							bytecode.NewLineInfo(5, 1),
						},
						[]value.Symbol{
							value.ToSymbol("a"),
							value.ToSymbol("b"),
						},
						2,
						-1,
						false,
						false,
						[]value.Value{
							value.SmallInt(5),
							value.ToSymbol("a"),
							value.String("b"),
							value.ToSymbol("b"),
						},
					),
					value.ToSymbol("foo"),
				},
			),
		},
		"define method with optional parameters in top level": {
			input: `
				def foo(a, b = 5.2, c = 10)
					d := 5
					a + b + c + d
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethod(
						value.ToSymbol("foo"),
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,

							byte(bytecode.GET_LOCAL8), 2,
							byte(bytecode.JUMP_UNLESS_UNDEF), 0, 5,
							byte(bytecode.POP),
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 2,
							byte(bytecode.POP),

							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.JUMP_UNLESS_UNDEF), 0, 5,
							byte(bytecode.POP),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),

							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.SET_LOCAL8), 4,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.GET_LOCAL8), 2,
							byte(bytecode.ADD),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.ADD),
							byte(bytecode.GET_LOCAL8), 4,
							byte(bytecode.ADD),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(70, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(2, 13),
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(4, 7),
							bytecode.NewLineInfo(5, 1),
						},
						[]value.Symbol{
							value.ToSymbol("a"),
							value.ToSymbol("b"),
							value.ToSymbol("c"),
						},
						2,
						-1,
						false,
						false,
						[]value.Value{
							value.Float(5.2),
							value.SmallInt(10),
							value.SmallInt(5),
						},
					),
					value.ToSymbol("foo"),
				},
			),
		},
		"define method with required parameters in a class": {
			input: `
				class Bar
					def foo(a, b)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(79, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(78, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							vm.NewBytecodeMethod(
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
								L(P(20, 3, 6), P(70, 6, 8)),
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
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Bar"),
				},
			),
		},
		"define method with required parameters in a module": {
			input: `
				module Bar
					def foo(a, b)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(80, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(79, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							vm.NewBytecodeMethod(
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
								L(P(21, 3, 6), P(71, 6, 8)),
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
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Bar"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestDefInit(t *testing.T) {
	tests := testTable{
		"define init in top level": {
			input: `
				init then :bar
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(5, 2, 5), P(18, 2, 18)), "init cannot be defined in the top level"),
			},
		},
		"define init in a module": {
			input: `
				module Foo
					init then :bar
				end
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(21, 3, 6), P(34, 3, 19)), "modules cannot have initializers"),
			},
		},
		"define init in a method": {
			input: `
				def foo
					init then :bar
				end
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(18, 3, 6), P(31, 3, 19)), "methods cannot be nested: #init"),
			},
		},
		"define init in init": {
			input: `
				class Foo
				  init
					  init then :bar
				  end
				end
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(33, 4, 8), P(46, 4, 21)), "methods cannot be nested: #init"),
			},
		},
		"define with required parameters in a class": {
			input: `
				class Bar
					init(a, b)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(75, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							vm.NewBytecodeMethod(
								value.ToSymbol("#init"),
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
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(20, 3, 6), P(67, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 5),
									bytecode.NewLineInfo(6, 2),
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
							value.ToSymbol("#init"),
						},
					),
					value.ToSymbol("Bar"),
				},
			),
		},
		"define with required parameters in a mixin": {
			input: `
				mixin Bar
					init(a, b)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(75, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							vm.NewBytecodeMethod(
								value.ToSymbol("#init"),
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
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(20, 3, 6), P(67, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 5),
									bytecode.NewLineInfo(6, 2),
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
							value.ToSymbol("#init"),
						},
					),
					value.ToSymbol("Bar"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestDefMixin(t *testing.T) {
	tests := testTable{
		"anonymous mixin without a body": {
			input: "mixin; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				nil,
			),
		},
		"mixin with a relative name without a body": {
			input: "mixin Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
				},
			),
		},
		"mixin with an absolute name without a body": {
			input: "mixin ::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
				},
			),
		},
		"named mixin inside of a method": {
			input: `
				def foo
					mixin Bar; end
				end
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(18, 3, 6), P(31, 3, 19)), "cannot define named mixins inside of a method: foo"),
			},
		},
		"mixin with an absolute nested name without a body": {
			input: "mixin ::Std::Int::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 1, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.ToSymbol("Std"),
					value.ToSymbol("Int"),
					value.ToSymbol("Foo"),
				},
			),
		},
		"anonymous mixin with a body": {
			input: `
				mixin
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_ANON_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(40, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
				},
			),
		},
		"mixin with a body": {
			input: `
				mixin Foo
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(44, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"nested mixins": {
			input: `
				mixin Foo
					mixin Bar
						a := 1
						a + 2
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_MIXIN),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(70, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							vm.NewBytecodeMethodNoParams(
								mixinSymbol,
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.SET_LOCAL8), 3,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 3,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(20, 3, 6), P(62, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 2),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
							value.ToSymbol("Bar"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestInclude(t *testing.T) {
	tests := testTable{
		"include a global constant in a class": {
			input: `
				class Foo
					include ::Bar
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(40, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"include two constants in a class": {
			input: `
				class Foo
					include ::Bar, ::Baz
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(48, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.INCLUDE),

							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 1,
							byte(bytecode.SELF),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),

							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(47, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 9),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							value.ToSymbol("Baz"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"include in top level": {
			input: `include ::Bar`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(0, 1, 1), P(12, 1, 13)),
					"cannot include mixins in the top level",
				),
			},
		},
		"include in a module": {
			input: `
				module Foo
					include ::Bar
				end
			`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(21, 3, 6), P(33, 3, 18)),
					"cannot include mixins in a module",
				),
			},
		},
		"include in a method": {
			input: `
				def foo
					include ::Bar
				end
			`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(18, 3, 6), P(30, 3, 18)),
					"cannot include mixins in a method",
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestExtend(t *testing.T) {
	tests := testTable{
		"extend a class with a global constant": {
			input: `
				class Foo
					extend ::Bar
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(39, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 6),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"extend a module with a global constant": {
			input: `
				module Foo
					extend ::Bar
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(40, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 6),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"extend a mixin with a global constant": {
			input: `
				mixin Foo
					extend ::Bar
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(39, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 6),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"extend a class with two global constants": {
			input: `
				class Foo
					extend ::Bar, ::Baz
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.INCLUDE),

							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 1,
							byte(bytecode.SELF),
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),

							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(46, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 11),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							value.ToSymbol("Baz"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"extend in top level": {
			input: `extend ::Bar`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(0, 1, 1), P(11, 1, 12)),
					"cannot extend mixins in the top level",
				),
			},
		},
		"extend in a method": {
			input: `
				def foo
					extend ::Bar
				end
			`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(18, 3, 6), P(29, 3, 17)),
					"cannot extend mixins in a method",
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}
