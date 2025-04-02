package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/diagnostic"
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 2, 5), P(44, 4, 7)), "singleton definitions cannot appear in this context"),
				diagnostic.NewFailure(L(P(20, 3, 6), P(36, 3, 22)), "method definitions cannot appear in this context"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(18, 3, 6), P(59, 5, 8)), "singleton definitions cannot appear in this context"),
				diagnostic.NewFailure(L(P(34, 4, 7), P(50, 4, 23)), "method definitions cannot appear in this context"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(14, 2, 14), P(16, 2, 16)), "cannot declare parameter `arg` without a type"),
				diagnostic.NewFailure(L(P(24, 3, 6), P(65, 5, 8)), "singleton definitions cannot appear in this context"),
				diagnostic.NewFailure(L(P(40, 4, 7), P(56, 4, 23)), "method definitions cannot appear in this context"),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Baz").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Baz").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE_0),
									byte(bytecode.RETURN),
								},
								L(P(36, 4, 7), P(52, 4, 23)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 2),
								},
								[]value.Value{
									value.ToSymbol("bar").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Baz").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Baz").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE_0),
									byte(bytecode.RETURN),
								},
								L(P(36, 4, 7), P(52, 4, 23)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 2),
								},
								[]value.Value{
									value.ToSymbol("bar").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(74, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 3,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(74, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Baz").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(74, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Baz").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE_0),
									byte(bytecode.RETURN),
								},
								L(P(40, 4, 7), P(56, 4, 23)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 2),
								},
								[]value.Value{
									value.ToSymbol("bar").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(21, 3, 6), P(62, 5, 8)), "singleton definitions cannot appear in this context"),
				diagnostic.NewFailure(L(P(37, 4, 7), P(53, 4, 23)), "method definitions cannot appear in this context"),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
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
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(38, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(38, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(69, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(69, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(69, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("bar").ToValue(),
							value.ToSymbol("baz").ToValue(),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
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
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(47, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_SETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(47, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(58, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(58, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_SETTER),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_SETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(58, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("baz").ToValue(),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(45, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_SETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(45, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(56, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(56, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_SETTER),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_SETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(56, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("baz").ToValue(),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(59, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(59, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.DEF_METHOD_ALIAS),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(59, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
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
							)),
							value.ToSymbol("bar").ToValue(),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(124, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(124, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 5,
							byte(bytecode.LOAD_VALUE8), 6,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 7,
							byte(bytecode.LOAD_VALUE_2),
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
							bytecode.NewLineInfo(1, 29),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
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
							)),
							value.ToSymbol("add").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
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
							)),
							value.ToSymbol("bar").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
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
							)),
							value.ToSymbol("delete").ToValue(),
							value.ToSymbol("plus").ToValue(),
							value.ToSymbol("foo").ToValue(),
							value.ToSymbol("remove").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(13, 1, 14)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 12),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"named class inside of a method": {
			input: `
				def foo
				  class ::Bar; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(19, 3, 7), P(34, 3, 22)), "class definitions cannot appear in this context"),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 3, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(3, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_2),
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
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(3, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Undefined,
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 5, 30)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_3),
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
							bytecode.NewLineInfo(1, 25),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Baz").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Baz::Bar").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"class with an absolute name without a body": {
			input: "class ::Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(15, 1, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 12),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"class with an absolute nested name without a body": {
			input: "class ::Std::Int::Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 1, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(25, 1, 26)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 12),
						},
						[]value.Value{
							value.ToSymbol("Std::Int").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Int::Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Undefined,
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(45, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Undefined,
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<class: Foo>"),
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.INT_1),
							byte(bytecode.SET_LOCAL_1),
							byte(bytecode.GET_LOCAL_1),
							byte(bytecode.INT_2),
							byte(bytecode.ADD_INT),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(44, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						nil,
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE_2),
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
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
							value.ToSymbol("Foo::Bar").ToValue(),
						},
					)),
					value.Undefined,
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<class: Foo>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.INIT_NAMESPACE),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(70, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							value.ToSymbol("Foo::Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("<class: Foo::Bar>"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.INT_1),
									byte(bytecode.SET_LOCAL_1),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.INT_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(62, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 3),
								},
								nil,
							)),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(14, 1, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"named module inside of a method": {
			input: `
				def foo
					module Bar; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(18, 3, 6), P(32, 3, 20)), "module definitions cannot appear in this context"),
			},
		},
		"class with an absolute name without a body": {
			input: "module ::Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(16, 1, 17)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"class with an absolute nested name without a body": {
			input: "module ::Std::Int::Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 1, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(26, 1, 27)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
						},
						[]value.Value{
							value.ToSymbol("Std::Int").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(46, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<module: Foo>"),
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.INT_1),
							byte(bytecode.SET_LOCAL_1),
							byte(bytecode.GET_LOCAL_1),
							byte(bytecode.INT_2),
							byte(bytecode.ADD_INT),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(45, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						nil,
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(73, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Undefined,
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<module: Foo>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.INIT_NAMESPACE),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(72, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							value.ToSymbol("Foo::Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("<module: Foo::Bar>"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.INT_1),
									byte(bytecode.SET_LOCAL_1),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.INT_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(64, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 3),
								},
								nil,
							)),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 2, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(2, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(22, 2, 22)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(2, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE_0),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(21, 2, 21)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 2),
								},
								[]value.Value{
									value.ToSymbol("bar").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(49, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(49, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.UNDEFINED),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.INT_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.NEW_ARRAY_TUPLE8), 1,
									byte(bytecode.CALL_SELF_TCO8), 0,
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(5, 2, 5), P(48, 4, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 8),
									bytecode.NewLineInfo(4, 2),
								},
								1,
								0,
								[]value.Value{
									value.Ref(value.NewCallSiteInfo(
										value.ToSymbol("println"),
										1,
									)),
								},
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(69, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(69, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.UNDEFINED),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.INT_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.NEW_ARRAY_TUPLE8), 1,
									byte(bytecode.CALL_SELF8), 0,
									byte(bytecode.POP),
									byte(bytecode.LOAD_VALUE_1),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(5, 2, 5), P(68, 5, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 9),
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								[]value.Value{
									value.Ref(value.NewCallSiteInfo(
										value.ToSymbol("println"),
										1,
									)),
									value.Ref(value.String("siema")),
								},
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(54, 4, 13), P(60, 4, 19)), "values returned in void context will be ignored"),
			},
		},
		"define generator": {
			input: `
				def *foo(a: Int, b: Int = 2): Int ! String
					yield a + b
					throw "lol"
					return 10
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(104, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(104, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunctionWithCatchEntries(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.JUMP_UNLESS_UNDEF), 0, 2,
									byte(bytecode.INT_2),
									byte(bytecode.SET_LOCAL_2),
									byte(bytecode.GENERATOR),
									byte(bytecode.RETURN),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.YIELD),
									byte(bytecode.LOAD_VALUE_0),
									byte(bytecode.THROW),
									byte(bytecode.POP),
									byte(bytecode.LOAD_INT_8), 10,
									byte(bytecode.YIELD),
									byte(bytecode.STOP_ITERATION),
									byte(bytecode.YIELD),
									byte(bytecode.STOP_ITERATION),
									byte(bytecode.LOOP), 0, 4,
								},
								L(P(5, 2, 5), P(103, 6, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 7),
									bytecode.NewLineInfo(6, 1),
									bytecode.NewLineInfo(3, 4),
									bytecode.NewLineInfo(4, 3),
									bytecode.NewLineInfo(5, 4),
									bytecode.NewLineInfo(6, 2),
									bytecode.NewLineInfo(2, 3),
								},
								2,
								1,
								[]value.Value{
									value.Ref(value.String("lol")),
								},
								[]*vm.CatchEntry{
									{From: -1, To: -1, JumpAddress: 8},
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(87, 5, 6), P(95, 5, 14)), "unreachable code"),
			},
		},
		"define an async method": {
			input: `
				async def foo(a: Int, b: Int = 2): Int ! String
					await timeout(5.seconds)
					return a + b
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(108, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(108, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunctionWithCatchEntries(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.JUMP_UNLESS_UNDEF), 0, 2,
									byte(bytecode.INT_2),
									byte(bytecode.SET_LOCAL_2),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.PROMISE),
									byte(bytecode.RETURN),
									byte(bytecode.INT_5),
									byte(bytecode.CALL_METHOD8), 0,
									byte(bytecode.UNDEFINED),
									byte(bytecode.CALL_SELF8), 1,
									byte(bytecode.AWAIT),
									byte(bytecode.AWAIT_RESULT),
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(107, 5, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 8),
									bytecode.NewLineInfo(5, 1),
									bytecode.NewLineInfo(3, 9),
									bytecode.NewLineInfo(4, 4),
								},
								3,
								2,
								[]value.Value{
									value.Ref(value.NewCallSiteInfo(value.ToSymbol("seconds"), 0)),
									value.Ref(value.NewCallSiteInfo(value.ToSymbol("timeout"), 2)),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(63, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.INT_5),
									byte(bytecode.SET_LOCAL_3),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.ADD_INT),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(62, 5, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 4),
									bytecode.NewLineInfo(4, 5),
									bytecode.NewLineInfo(5, 1),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(115, 9, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(9, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(115, 9, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(9, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(115, 9, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(9, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.DUP),
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
									value.ToSymbol("a").ToValue(),
									value.ToSymbol("b").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR8), 1,
									byte(bytecode.POP),
									byte(bytecode.INT_5),
									byte(bytecode.SET_LOCAL_3),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.ADD_INT),
									byte(bytecode.RETURN),
								},
								L(P(54, 5, 6), P(106, 8, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 12),
									bytecode.NewLineInfo(6, 2),
									bytecode.NewLineInfo(7, 5),
									bytecode.NewLineInfo(8, 1),
								},
								2,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.ToSymbol("b").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(136, 9, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(9, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(136, 9, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(9, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(136, 9, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(9, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.JUMP_UNLESS_UNDEF), 0, 2,
									byte(bytecode.INT_5),
									byte(bytecode.SET_LOCAL_1),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.JUMP_UNLESS_UNDEF), 0, 3,
									byte(bytecode.LOAD_INT_8), 21,
									byte(bytecode.SET_LOCAL_2),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR8), 1,
									byte(bytecode.POP),
									byte(bytecode.INT_5),
									byte(bytecode.SET_LOCAL_3),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.ADD_INT),
									byte(bytecode.RETURN),
								},
								L(P(56, 5, 6), P(127, 8, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 25),
									bytecode.NewLineInfo(6, 2),
									bytecode.NewLineInfo(7, 5),
									bytecode.NewLineInfo(8, 1),
								},
								2,
								2,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.ToSymbol("b").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(88, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(88, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.JUMP_UNLESS_UNDEF), 0, 2,
									byte(bytecode.LOAD_VALUE_0),
									byte(bytecode.SET_LOCAL_2),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.JUMP_UNLESS_UNDEF), 0, 3,
									byte(bytecode.LOAD_INT_8), 10,
									byte(bytecode.SET_LOCAL_3),
									byte(bytecode.INT_5),
									byte(bytecode.SET_LOCAL_4),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.ADD_FLOAT),
									byte(bytecode.GET_LOCAL_4),
									byte(bytecode.ADD_FLOAT),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(87, 5, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 15),
									bytecode.NewLineInfo(3, 2),
									bytecode.NewLineInfo(4, 7),
									bytecode.NewLineInfo(5, 1),
								},
								3,
								2,
								[]value.Value{
									value.Float(5.2).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(89, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(89, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(89, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.INT_5),
									byte(bytecode.SET_LOCAL_3),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.ADD_INT),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(80, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 5),
									bytecode.NewLineInfo(6, 1),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(90, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(90, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(90, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.INT_5),
									byte(bytecode.SET_LOCAL_3),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.ADD_INT),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(81, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 5),
									bytecode.NewLineInfo(6, 1),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 2, 5), P(18, 2, 18)), "init definitions cannot appear outside of classes"),
			},
		},
		"define init in a module": {
			input: `
				module Foo
					init then :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(21, 3, 6), P(34, 3, 19)), "init definitions cannot appear outside of classes"),
			},
		},
		"define init in a mixin": {
			input: `
				mixin Foo
					init then :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(20, 3, 6), P(33, 3, 19)), "init definitions cannot appear outside of classes"),
			},
		},
		"define init in an interface": {
			input: `
				interface Foo
					init then :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(24, 3, 6), P(37, 3, 19)), "init definitions cannot appear outside of classes"),
				diagnostic.NewFailure(L(P(24, 3, 6), P(37, 3, 19)), "method `#init` cannot have a body because it is abstract"),
			},
		},
		"define init in a method": {
			input: `
				def foo
					init then :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(18, 3, 6), P(31, 3, 19)), "method definitions cannot appear in this context"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(33, 4, 8), P(46, 4, 21)), "method definitions cannot appear in this context"),
			},
		},
		"define with required parameters in a class": {
			input: `
				class Bar
					init(a: Int, b: Int)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(86, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.INT_5),
									byte(bytecode.SET_LOCAL_3),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.ADD_INT),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(20, 3, 6), P(77, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 5),
									bytecode.NewLineInfo(6, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
						},
					)),
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
		"mixin with a relative name without a body": {
			input: "mixin Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(13, 1, 14)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"mixin with an absolute name without a body": {
			input: "mixin ::Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(15, 1, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"named mixin inside of a method": {
			input: `
				def foo
					mixin Bar; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(18, 3, 6), P(31, 3, 19)), "mixin definitions cannot appear in this context"),
			},
		},
		"mixin with an absolute nested name without a body": {
			input: "mixin ::Std::Int::Foo; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 1, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(25, 1, 26)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
						},
						[]value.Value{
							value.ToSymbol("Std::Int").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(45, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<mixin: Foo>"),
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.INT_1),
							byte(bytecode.SET_LOCAL_1),
							byte(bytecode.GET_LOCAL_1),
							byte(bytecode.INT_2),
							byte(bytecode.ADD_INT),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(44, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						nil,
					)),
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(71, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Undefined,
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<mixin: Foo>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.INIT_NAMESPACE),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(70, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							value.ToSymbol("Foo::Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("<mixin: Foo::Bar>"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.INT_1),
									byte(bytecode.SET_LOCAL_1),
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.INT_2),
									byte(bytecode.ADD_INT),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(62, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 3),
								},
								nil,
							)),
						},
					)),
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
				mixin Bar; end
				class Foo
					include ::Bar
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(60, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(60, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"include a global constant in a mixin": {
			input: `
				mixin Bar; end
				mixin Foo
					include ::Bar
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(60, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(60, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 15),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"include two constants in a class": {
			input: `
				mixin Bar; end
				mixin Baz; end
				class Foo
					include ::Bar, ::Baz
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(86, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_NAMESPACE), 2,
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.GET_CONST8), 4,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.INCLUDE),
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 30),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Baz").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"include in top level": {
			input: `
				mixin Bar; end
				include ::Bar
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(24, 3, 5), P(36, 3, 17)),
					"cannot include mixins in this context",
				),
			},
		},
		"include in a module": {
			input: `
				mixin Bar; end
				module Foo
					include ::Bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(40, 4, 6), P(52, 4, 18)),
					"cannot include mixins in this context",
				),
			},
		},
		"include in an interface": {
			input: `
				mixin Bar; end
				interface Foo
					include ::Bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(43, 4, 6), P(55, 4, 18)),
					"cannot include mixins in this context",
				),
			},
		},
		"include in a method": {
			input: `
				mixin Bar; end
				def foo
					include ::Bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(37, 4, 6), P(49, 4, 18)),
					"cannot include mixins in this context",
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
