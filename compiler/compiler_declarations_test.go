package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/error"
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
			err: error.ErrorList{
				error.NewFailure(L(P(5, 2, 5), P(44, 4, 7)), "singleton definitions cannot appear in this context"),
				error.NewFailure(L(P(20, 3, 6), P(36, 3, 22)), "method definitions cannot appear in this context"),
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
			err: error.ErrorList{
				error.NewFailure(L(P(18, 3, 6), P(59, 5, 8)), "singleton definitions cannot appear in this context"),
				error.NewFailure(L(P(34, 4, 7), P(50, 4, 23)), "method definitions cannot appear in this context"),
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
			err: error.ErrorList{
				error.NewFailure(L(P(14, 2, 14), P(16, 2, 16)), "cannot declare parameter `arg` without a type"),
				error.NewFailure(L(P(24, 3, 6), P(65, 5, 8)), "singleton definitions cannot appear in this context"),
				error.NewFailure(L(P(40, 4, 7), P(56, 4, 23)), "method definitions cannot appear in this context"),
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Baz"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Baz"),
							vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(36, 4, 7), P(52, 4, 23)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
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
		},
		"define in a mixin": {
			input: `
				mixin Baz
					singleton
						def foo then :bar
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Baz"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Baz"),
							vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(36, 4, 7), P(52, 4, 23)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
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
		},
		"define in an interface": {
			input: `
				interface Baz
					singleton
						def foo then :bar
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(74, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 3,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(74, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Baz"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(74, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Baz"),
							vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(40, 4, 7), P(56, 4, 23)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
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
		},
		"define in a module": {
			input: `
				module Baz
					singleton
						def foo then :bar
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L(P(21, 3, 6), P(62, 5, 8)), "singleton definitions cannot appear in this context"),
				error.NewFailure(L(P(37, 4, 7), P(53, 4, 23)), "method definitions cannot appear in this context"),
			},
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
			input: `
				class Foo
					getter foo
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(38, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(38, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_GETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(38, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							value.ToSymbol("foo"),
						},
					),
				},
			),
		},
		"define three getters": {
			input: `
				class Foo
					getter foo: Foo?, bar: Int?, baz: String?
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(69, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(69, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.DEF_GETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(69, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 12),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							value.ToSymbol("bar"),
							value.ToSymbol("baz"),
							value.ToSymbol("foo"),
						},
					),
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
			input: `
				class Foo
					setter foo: String?
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(47, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_SETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(47, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							value.ToSymbol("foo"),
						},
					),
				},
			),
		},
		"define three setters": {
			input: `
				class Foo
					setter foo: Foo?, baz: String?
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(58, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(58, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_SETTER),
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_SETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(58, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							value.ToSymbol("baz"),
							value.ToSymbol("foo"),
						},
					),
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

func TestAttr(t *testing.T) {
	tests := testTable{
		"define single attr": {
			input: `
				class Foo
					attr foo: String?
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(45, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_SETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(45, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							value.ToSymbol("foo"),
						},
					),
				},
			),
		},
		"define three attrs": {
			input: `
				class Foo
					attr foo: Foo?, baz: String?
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(56, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(56, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_SETTER),
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_SETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(56, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 15),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							value.ToSymbol("baz"),
							value.ToSymbol("foo"),
						},
					),
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
			input: `
				class Foo
					def bar; end
					alias foo bar
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(59, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(59, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.DEF_METHOD_ALIAS),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(59, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("bar"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(31, 3, 17)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								nil,
							),
							value.ToSymbol("bar"),
							value.ToSymbol("foo"),
						},
					),
				},
			),
		},
		"define three aliases": {
			input: `
				class Foo
					def bar; end
					def delete; end
					def plus; end
					alias foo bar, remove delete, add plus
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(124, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(124, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 5,
							byte(bytecode.LOAD_VALUE8), 6,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 7,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD_ALIAS),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.LOAD_VALUE8), 8,
							byte(bytecode.DEF_METHOD_ALIAS),
							byte(bytecode.LOAD_VALUE8), 6,
							byte(bytecode.LOAD_VALUE8), 9,
							byte(bytecode.DEF_METHOD_ALIAS),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(124, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 33),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("plus"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(59, 5, 6), P(71, 5, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 2),
								},
								nil,
							),
							value.ToSymbol("add"),
							vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("bar"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(31, 3, 17)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								nil,
							),
							value.ToSymbol("bar"),
							vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("delete"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(38, 4, 6), P(52, 4, 20)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 2),
								},
								nil,
							),
							value.ToSymbol("delete"),
							value.ToSymbol("plus"),
							value.ToSymbol("foo"),
							value.ToSymbol("remove"),
						},
					),
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
		"class with a relative name without a body": {
			input: "class Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(13, 1, 14)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					nil,
				},
			),
		},
		"named class inside of a method": {
			input: `
				def foo
				  class ::Bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L(P(19, 3, 7), P(34, 3, 22)), "class definitions cannot appear in this context"),
			},
		},
		"class with an absolute parent": {
			input: `
				class Bar; end
				class Foo < ::Bar; end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 3, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(3, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(46, 3, 27)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 22),
							bytecode.NewLineInfo(3, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					nil,
				},
			),
		},
		"class with an absolute nested parent": {
			input: `
				module Baz
					class Bar; end
				end
				class Foo < Baz::Bar; end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 5, 30)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 4,
							byte(bytecode.GET_CONST8), 5,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.GET_CONST8), 4,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(73, 5, 30)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 28),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Baz"),
							value.ToSymbol("Bar"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Baz::Bar"),
							value.ToSymbol("Std::Object"),
						},
					),
					nil,
				},
			),
		},
		"class with an absolute name without a body": {
			input: "class ::Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(15, 1, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					nil,
				},
			),
		},
		"class with an absolute nested name without a body": {
			input: "class ::Std::Int::Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 1, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(25, 1, 26)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
						},
						[]value.Value{
							value.ToSymbol("Std::Int"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Int::Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					nil,
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(45, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					nil,
					value.ToSymbol("Foo"),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<class: Foo>"),
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 1,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(44, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 7),
							bytecode.NewLineInfo(4, 5),
							bytecode.NewLineInfo(5, 3),
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.GET_CONST8), 4,
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(71, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 22),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Bar"),
							value.ToSymbol("Std::Object"),
							value.ToSymbol("Foo::Bar"),
						},
					),
					nil,
					value.ToSymbol("Foo"),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<class: Foo>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.INIT_NAMESPACE),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(70, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							value.ToSymbol("Foo::Bar"),
							vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("<class: Foo::Bar>"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.SET_LOCAL8), 1,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(62, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
									bytecode.NewLineInfo(5, 5),
									bytecode.NewLineInfo(6, 3),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
						},
					),
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
		"module with a relative name without a body": {
			input: "module Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(14, 1, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					nil,
				},
			),
		},
		"named module inside of a method": {
			input: `
				def foo
					module Bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L(P(18, 3, 6), P(32, 3, 20)), "module definitions cannot appear in this context"),
			},
		},
		"class with an absolute name without a body": {
			input: "module ::Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(16, 1, 17)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					nil,
				},
			),
		},
		"class with an absolute nested name without a body": {
			input: "module ::Std::Int::Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 1, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(26, 1, 27)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
						},
						[]value.Value{
							value.ToSymbol("Std::Int"),
							value.ToSymbol("Foo"),
						},
					),
					nil,
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(46, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					nil,
					value.ToSymbol("Foo"),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<module: Foo>"),
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 1,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(45, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 7),
							bytecode.NewLineInfo(4, 5),
							bytecode.NewLineInfo(5, 3),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(73, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 12),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Bar"),
						},
					),
					nil,
					value.ToSymbol("Foo"),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<module: Foo>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.INIT_NAMESPACE),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(72, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							value.ToSymbol("Foo::Bar"),
							vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("<module: Foo::Bar>"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.SET_LOCAL8), 1,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(64, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
									bytecode.NewLineInfo(5, 5),
									bytecode.NewLineInfo(6, 3),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
						},
					),
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 2, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(2, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(22, 2, 22)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(2, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel"),
							vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(21, 2, 21)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 3),
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
		},
		"define a setter": {
			input: `
				def foo=(a: Int)
					println(a + 2)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(49, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(49, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.GET_CONST8), 0,
									byte(bytecode.UNDEFINED),
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.NEW_ARRAY_TUPLE8), 1,
									byte(bytecode.CALL_METHOD8), 2,
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(5, 2, 5), P(48, 4, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 12),
									bytecode.NewLineInfo(4, 2),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("Std::Kernel"),
									value.SmallInt(2),
									value.NewCallSiteInfo(
										value.ToSymbol("println"),
										1,
									),
								},
							),
							value.ToSymbol("foo="),
						},
					),
				},
			),
		},
		"define a setter with return": {
			input: `
				def foo=(a: Int)
					println(a + 2)
					return "siema"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(69, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(69, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.GET_CONST8), 0,
									byte(bytecode.UNDEFINED),
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.NEW_ARRAY_TUPLE8), 1,
									byte(bytecode.CALL_METHOD8), 2,
									byte(bytecode.POP),
									byte(bytecode.LOAD_VALUE8), 3,
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(5, 2, 5), P(68, 5, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 13),
									bytecode.NewLineInfo(4, 4),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("Std::Kernel"),
									value.SmallInt(2),
									value.NewCallSiteInfo(
										value.ToSymbol("println"),
										1,
									),
									value.String("siema"),
								},
							),
							value.ToSymbol("foo="),
						},
					),
				},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(54, 4, 13), P(60, 4, 19)), "values returned in void context will be ignored"),
			},
		},
		"define method with required parameters in top level": {
			input: `
				def foo(a: Int, b: Int)
					c := 5
					a + b + c
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(63, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel"),
							vm.NewBytecodeFunction(
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
								L(P(5, 2, 5), P(62, 5, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 7),
									bytecode.NewLineInfo(4, 8),
									bytecode.NewLineInfo(5, 1),
								},
								2,
								0,
								[]value.Value{
									value.SmallInt(5),
								},
							),
							value.ToSymbol("foo"),
						},
					),
				},
			),
		},
		"define method with ivar parameters": {
			input: `
				class Bar
					init(@a: Int, @b: Int); end

					def foo(@a, @b)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(115, 9, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(9, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(115, 9, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(9, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(115, 9, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(9, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 2,
									byte(bytecode.SET_IVAR8), 1,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(20, 3, 6), P(46, 3, 32)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 13),
								},
								2,
								0,
								[]value.Value{
									value.ToSymbol("a"),
									value.ToSymbol("b"),
								},
							),
							value.ToSymbol("#init"),
							vm.NewBytecodeFunction(
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
								L(P(54, 5, 6), P(106, 8, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 12),
									bytecode.NewLineInfo(6, 5),
									bytecode.NewLineInfo(7, 8),
									bytecode.NewLineInfo(8, 1),
								},
								2,
								0,
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
			),
		},
		"define method with default ivar parameters": {
			input: `
				class Bar
					var @a: Int?
					var @b: Int?
					def foo(@a: Int = 5, @b: Int = 21)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(136, 9, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(9, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(136, 9, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(9, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(136, 9, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(9, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
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
								L(P(56, 5, 6), P(127, 8, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 34),
									bytecode.NewLineInfo(6, 5),
									bytecode.NewLineInfo(7, 8),
									bytecode.NewLineInfo(8, 1),
								},
								2,
								2,
								[]value.Value{
									value.SmallInt(5),
									value.ToSymbol("a"),
									value.SmallInt(21),
									value.ToSymbol("b"),
								},
							),
							value.ToSymbol("foo"),
						},
					),
				},
			),
		},
		"define method with optional parameters in top level": {
			input: `
				def foo(a: Int, b: Float = 5.2, c: Int = 10)
					d := 5
					a + b + c + d
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(88, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(88, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel"),
							vm.NewBytecodeFunction(
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
								L(P(5, 2, 5), P(87, 5, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 24),
									bytecode.NewLineInfo(3, 5),
									bytecode.NewLineInfo(4, 11),
									bytecode.NewLineInfo(5, 1),
								},
								3,
								2,
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
			),
		},
		"define method with required parameters in a class": {
			input: `
				class Bar
					def foo(a: Int, b: Int)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(89, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(89, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(89, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
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
								L(P(20, 3, 6), P(80, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
									bytecode.NewLineInfo(5, 8),
									bytecode.NewLineInfo(6, 1),
								},
								2,
								0,
								[]value.Value{
									value.SmallInt(5),
								},
							),
							value.ToSymbol("foo"),
						},
					),
				},
			),
		},
		"define method with required parameters in a module": {
			input: `
				module Bar
					def foo(a: Int, b: Int)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(90, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(90, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(90, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
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
								L(P(21, 3, 6), P(81, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
									bytecode.NewLineInfo(5, 8),
									bytecode.NewLineInfo(6, 1),
								},
								2,
								0,
								[]value.Value{
									value.SmallInt(5),
								},
							),
							value.ToSymbol("foo"),
						},
					),
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

// func TestDefInit(t *testing.T) {
// 	tests := testTable{
// 		"define init in top level": {
// 			input: `
// 				init then :bar
// 			`,
// 			err: error.ErrorList{
// 				error.NewFailure(L(P(5, 2, 5), P(18, 2, 18)), "init cannot be defined in the top level"),
// 			},
// 		},
// 		"define init in a module": {
// 			input: `
// 				module Foo
// 					init then :bar
// 				end
// 			`,
// 			err: error.ErrorList{
// 				error.NewFailure(L(P(21, 3, 6), P(34, 3, 19)), "modules cannot have initializers"),
// 			},
// 		},
// 		"define init in a method": {
// 			input: `
// 				def foo
// 					init then :bar
// 				end
// 			`,
// 			err: error.ErrorList{
// 				error.NewFailure(L(P(18, 3, 6), P(31, 3, 19)), "methods cannot be nested: #init"),
// 			},
// 		},
// 		"define init in init": {
// 			input: `
// 				class Foo
// 				  init
// 					  init then :bar
// 				  end
// 				end
// 			`,
// 			err: error.ErrorList{
// 				error.NewFailure(L(P(33, 4, 8), P(46, 4, 21)), "methods cannot be nested: #init"),
// 			},
// 		},
// 		"define with required parameters in a class": {
// 			input: `
// 				class Bar
// 					init(a, b)
// 						c := 5
// 						a + b + c
// 					end
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.CONSTANT_CONTAINER),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.INIT_CLASS), 0,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(76, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(7, 1),
// 				},
// 				[]value.Value{
// 					vm.NewBytecodeFunctionNoParams(
// 						classSymbol,
// 						[]byte{
// 							byte(bytecode.LOAD_VALUE8), 0,
// 							byte(bytecode.LOAD_VALUE8), 1,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.RETURN_SELF),
// 						},
// 						L(P(5, 2, 5), P(75, 7, 7)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(3, 5),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							vm.NewBytecodeFunction(
// 								value.ToSymbol("#init"),
// 								[]byte{
// 									byte(bytecode.PREP_LOCALS8), 1,
// 									byte(bytecode.LOAD_VALUE8), 0,
// 									byte(bytecode.SET_LOCAL8), 3,
// 									byte(bytecode.POP),
// 									byte(bytecode.GET_LOCAL8), 1,
// 									byte(bytecode.GET_LOCAL8), 2,
// 									byte(bytecode.ADD),
// 									byte(bytecode.GET_LOCAL8), 3,
// 									byte(bytecode.ADD),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(20, 3, 6), P(67, 6, 8)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 									bytecode.NewLineInfo(5, 8),
// 									bytecode.NewLineInfo(6, 2),
// 								},
// 								[]value.Symbol{
// 									value.ToSymbol("a"),
// 									value.ToSymbol("b"),
// 								},
// 								0,
// 								-1,
// 								false,
// 								false,
// 								[]value.Value{
// 									value.SmallInt(5),
// 								},
// 							),
// 							value.ToSymbol("#init"),
// 						},
// 					),
// 					value.ToSymbol("Bar"),
// 				},
// 			),
// 		},
// 		"define with required parameters in a mixin": {
// 			input: `
// 				mixin Bar
// 					init(a, b)
// 						c := 5
// 						a + b + c
// 					end
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.CONSTANT_CONTAINER),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.INIT_MIXIN),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(76, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 6),
// 					bytecode.NewLineInfo(7, 1),
// 				},
// 				[]value.Value{
// 					vm.NewBytecodeFunctionNoParams(
// 						mixinSymbol,
// 						[]byte{
// 							byte(bytecode.LOAD_VALUE8), 0,
// 							byte(bytecode.LOAD_VALUE8), 1,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.RETURN_SELF),
// 						},
// 						L(P(5, 2, 5), P(75, 7, 7)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(3, 5),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							vm.NewBytecodeFunction(
// 								value.ToSymbol("#init"),
// 								[]byte{
// 									byte(bytecode.PREP_LOCALS8), 1,
// 									byte(bytecode.LOAD_VALUE8), 0,
// 									byte(bytecode.SET_LOCAL8), 3,
// 									byte(bytecode.POP),
// 									byte(bytecode.GET_LOCAL8), 1,
// 									byte(bytecode.GET_LOCAL8), 2,
// 									byte(bytecode.ADD),
// 									byte(bytecode.GET_LOCAL8), 3,
// 									byte(bytecode.ADD),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(20, 3, 6), P(67, 6, 8)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 									bytecode.NewLineInfo(5, 8),
// 									bytecode.NewLineInfo(6, 2),
// 								},
// 								[]value.Symbol{
// 									value.ToSymbol("a"),
// 									value.ToSymbol("b"),
// 								},
// 								0,
// 								-1,
// 								false,
// 								false,
// 								[]value.Value{
// 									value.SmallInt(5),
// 								},
// 							),
// 							value.ToSymbol("#init"),
// 						},
// 					),
// 					value.ToSymbol("Bar"),
// 				},
// 			),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			compilerTest(tc, t)
// 		})
// 	}
// }

// func TestDefMixin(t *testing.T) {
// 	tests := testTable{
// 		"mixin with a relative name without a body": {
// 			input: "mixin Foo; end",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.CONSTANT_CONTAINER),
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.INIT_MIXIN),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(13, 1, 14)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 6),
// 				},
// 				[]value.Value{
// 					value.ToSymbol("Foo"),
// 				},
// 			),
// 		},
// 		"mixin with an absolute name without a body": {
// 			input: "mixin ::Foo; end",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.ROOT),
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.INIT_MIXIN),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(15, 1, 16)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 6),
// 				},
// 				[]value.Value{
// 					value.ToSymbol("Foo"),
// 				},
// 			),
// 		},
// 		"named mixin inside of a method": {
// 			input: `
// 				def foo
// 					mixin Bar; end
// 				end
// 			`,
// 			err: error.ErrorList{
// 				error.NewFailure(L(P(18, 3, 6), P(31, 3, 19)), "cannot define named mixins inside of a method: foo"),
// 			},
// 		},
// 		"mixin with an absolute nested name without a body": {
// 			input: "mixin ::Std::Int::Foo; end",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 0,
// 					byte(bytecode.GET_MOD_CONST8), 1,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.INIT_MIXIN),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(25, 1, 26)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 10),
// 				},
// 				[]value.Value{
// 					value.ToSymbol("Std"),
// 					value.ToSymbol("Int"),
// 					value.ToSymbol("Foo"),
// 				},
// 			),
// 		},
// 		"mixin with a body": {
// 			input: `
// 				mixin Foo
// 					a := 1
// 					a + 2
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.CONSTANT_CONTAINER),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.INIT_MIXIN),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(45, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 6),
// 					bytecode.NewLineInfo(5, 1),
// 				},
// 				[]value.Value{
// 					vm.NewBytecodeFunctionNoParams(
// 						mixinSymbol,
// 						[]byte{
// 							byte(bytecode.PREP_LOCALS8), 1,
// 							byte(bytecode.LOAD_VALUE8), 0,
// 							byte(bytecode.SET_LOCAL8), 3,
// 							byte(bytecode.POP),
// 							byte(bytecode.GET_LOCAL8), 3,
// 							byte(bytecode.LOAD_VALUE8), 1,
// 							byte(bytecode.ADD),
// 							byte(bytecode.POP),
// 							byte(bytecode.RETURN_SELF),
// 						},
// 						L(P(5, 2, 5), P(44, 5, 7)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(3, 7),
// 							bytecode.NewLineInfo(4, 5),
// 							bytecode.NewLineInfo(5, 2),
// 						},
// 						[]value.Value{
// 							value.SmallInt(1),
// 							value.SmallInt(2),
// 						},
// 					),
// 					value.ToSymbol("Foo"),
// 				},
// 			),
// 		},
// 		"nested mixins": {
// 			input: `
// 				mixin Foo
// 					mixin Bar
// 						a := 1
// 						a + 2
// 					end
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.CONSTANT_CONTAINER),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.INIT_MIXIN),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(71, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 6),
// 					bytecode.NewLineInfo(7, 1),
// 				},
// 				[]value.Value{
// 					vm.NewBytecodeFunctionNoParams(
// 						mixinSymbol,
// 						[]byte{
// 							byte(bytecode.LOAD_VALUE8), 0,
// 							byte(bytecode.CONSTANT_CONTAINER),
// 							byte(bytecode.LOAD_VALUE8), 1,
// 							byte(bytecode.INIT_MIXIN),
// 							byte(bytecode.POP),
// 							byte(bytecode.RETURN_SELF),
// 						},
// 						L(P(5, 2, 5), P(70, 7, 7)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(3, 6),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							vm.NewBytecodeFunctionNoParams(
// 								mixinSymbol,
// 								[]byte{
// 									byte(bytecode.PREP_LOCALS8), 1,
// 									byte(bytecode.LOAD_VALUE8), 0,
// 									byte(bytecode.SET_LOCAL8), 3,
// 									byte(bytecode.POP),
// 									byte(bytecode.GET_LOCAL8), 3,
// 									byte(bytecode.LOAD_VALUE8), 1,
// 									byte(bytecode.ADD),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(20, 3, 6), P(62, 6, 8)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 									bytecode.NewLineInfo(5, 5),
// 									bytecode.NewLineInfo(6, 2),
// 								},
// 								[]value.Value{
// 									value.SmallInt(1),
// 									value.SmallInt(2),
// 								},
// 							),
// 							value.ToSymbol("Bar"),
// 						},
// 					),
// 					value.ToSymbol("Foo"),
// 				},
// 			),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			compilerTest(tc, t)
// 		})
// 	}
// }

// func TestInclude(t *testing.T) {
// 	tests := testTable{
// 		"include a global constant in a class": {
// 			input: `
// 				class Foo
// 					include ::Bar
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.CONSTANT_CONTAINER),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.INIT_CLASS), 0,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(41, 4, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(4, 1),
// 				},
// 				[]value.Value{
// 					vm.NewBytecodeFunctionNoParams(
// 						classSymbol,
// 						[]byte{
// 							byte(bytecode.ROOT),
// 							byte(bytecode.GET_MOD_CONST8), 0,
// 							byte(bytecode.SELF),
// 							byte(bytecode.INCLUDE),
// 							byte(bytecode.NIL),
// 							byte(bytecode.POP),
// 							byte(bytecode.RETURN_SELF),
// 						},
// 						L(P(5, 2, 5), P(40, 4, 7)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(3, 6),
// 							bytecode.NewLineInfo(4, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Bar"),
// 						},
// 					),
// 					value.ToSymbol("Foo"),
// 				},
// 			),
// 		},
// 		"include two constants in a class": {
// 			input: `
// 				class Foo
// 					include ::Bar, ::Baz
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.CONSTANT_CONTAINER),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.INIT_CLASS), 0,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(48, 4, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(4, 1),
// 				},
// 				[]value.Value{
// 					vm.NewBytecodeFunctionNoParams(
// 						classSymbol,
// 						[]byte{
// 							byte(bytecode.ROOT),
// 							byte(bytecode.GET_MOD_CONST8), 0,
// 							byte(bytecode.SELF),
// 							byte(bytecode.INCLUDE),

// 							byte(bytecode.ROOT),
// 							byte(bytecode.GET_MOD_CONST8), 1,
// 							byte(bytecode.SELF),
// 							byte(bytecode.INCLUDE),
// 							byte(bytecode.NIL),
// 							byte(bytecode.POP),

// 							byte(bytecode.RETURN_SELF),
// 						},
// 						L(P(5, 2, 5), P(47, 4, 7)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(3, 11),
// 							bytecode.NewLineInfo(4, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Bar"),
// 							value.ToSymbol("Baz"),
// 						},
// 					),
// 					value.ToSymbol("Foo"),
// 				},
// 			),
// 		},
// 		"include in top level": {
// 			input: `include ::Bar`,
// 			err: error.ErrorList{
// 				error.NewFailure(
// 					L(P(0, 1, 1), P(12, 1, 13)),
// 					"cannot include mixins in the top level",
// 				),
// 			},
// 		},
// 		"include in a module": {
// 			input: `
// 				module Foo
// 					include ::Bar
// 				end
// 			`,
// 			err: error.ErrorList{
// 				error.NewFailure(
// 					L(P(21, 3, 6), P(33, 3, 18)),
// 					"cannot include mixins in a module",
// 				),
// 			},
// 		},
// 		"include in a method": {
// 			input: `
// 				def foo
// 					include ::Bar
// 				end
// 			`,
// 			err: error.ErrorList{
// 				error.NewFailure(
// 					L(P(18, 3, 6), P(30, 3, 18)),
// 					"cannot include mixins in a method",
// 				),
// 			},
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			compilerTest(tc, t)
// 		})
// 	}
// }
