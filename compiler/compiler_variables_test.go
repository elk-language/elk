package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestInstanceVariables(t *testing.T) {
	tests := testTable{
		"initialise when declared": {
			input: `
				class Foo
					var @a: Int = 3
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(34, 3, 20), P(34, 3, 20)),
					"instance variables cannot be initialised when declared",
				),
			},
		},
		"declare in the top level": {
			input: "var @a: Float",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(0, 1, 1), P(12, 1, 13)),
					"cannot declare instance variable `@a` in this context",
				),
			},
		},
		"declare in a method": {
			input: "def foo; var @a: Float; end",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(9, 1, 10), P(21, 1, 22)),
					"instance variable definitions cannot appear in this context",
				),
			},
		},
		"declare in a class": {
			input: "class Foo; var @a: Float?; end",
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
				L(P(0, 1, 1), P(29, 1, 30)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
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
						L(P(0, 1, 1), P(29, 1, 30)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 12),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(29, 1, 30)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
				},
			),
		},
		"declare in a mixin": {
			input: "mixin Foo; var @a: Float; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 1, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
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
						L(P(0, 1, 1), P(28, 1, 29)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
				},
			),
		},
		"declare in a module": {
			input: "module Foo; var @a: Float?; end",
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
				L(P(0, 1, 1), P(30, 1, 31)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(30, 1, 31)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(30, 1, 31)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
				},
			),
		},
		"set instance variable in a class instance method": {
			input: `
				class Foo
					var @foo: Int?

					def foo
				  	@foo = 2
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(81, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(8, 2),
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
						L(P(0, 1, 1), P(81, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(81, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(81, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(41, 5, 6), P(72, 7, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 3),
									bytecode.NewLineInfo(7, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"set instance variable from parent in a class instance method": {
			input: `
				class Bar
					var @foo: Int?
				end

				class Foo < Bar
					def foo
				  	@foo = 2
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(109, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(10, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
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
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(69, 7, 6), P(100, 9, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(8, 3),
									bytecode.NewLineInfo(9, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"set instance variable from parent and self in a class instance method": {
			input: `
				class Bar
					var @bar: Int?
				end

				class Foo < Bar
					var @foo: Int?

					def foo
				  	@foo = 2
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(130, 12, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(12, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
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
						L(P(0, 1, 1), P(130, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(130, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("bar"): 0,
							}),
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("bar"): 0,
								value.ToSymbol("foo"): 1,
							}),
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
						L(P(0, 1, 1), P(130, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_1),
									byte(bytecode.RETURN),
								},
								L(P(90, 9, 6), P(121, 11, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(10, 3),
									bytecode.NewLineInfo(11, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"set instance variable from mixin in a class instance method": {
			input: `
				mixin Bar
					var @foo: Int?
				end

				class Foo
					include Bar

					def foo
				  	@foo = 2
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(121, 12, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(12, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
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
						L(P(0, 1, 1), P(121, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(121, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(121, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(81, 9, 6), P(112, 11, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(10, 3),
									bytecode.NewLineInfo(11, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"set instance variable in a mixin instance method": {
			input: `
				mixin Foo
					var @foo: Int?

					def foo
				  	@foo = 2
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
				L(P(0, 1, 1), P(81, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(8, 2),
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
						L(P(0, 1, 1), P(81, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
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
						L(P(0, 1, 1), P(81, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_NAME16), 0, 0,
									byte(bytecode.RETURN),
								},
								L(P(41, 5, 6), P(72, 7, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 5),
									bytecode.NewLineInfo(7, 1),
								},
								[]value.Value{
									value.ToSymbol("foo").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"set instance variable in a module method": {
			input: `
				module Foo
					var @foo: Int?

					def foo
				  	@foo = 2
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(82, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(8, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(82, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(82, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(82, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(42, 5, 6), P(73, 7, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 3),
									bytecode.NewLineInfo(7, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"set instance variable in a class method": {
			input: `
				class Foo
					singleton
						var @foo: Int?

						def foo
				  		@foo = 2
						end
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(109, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(10, 2),
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
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(58, 6, 7), P(91, 8, 9)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(7, 3),
									bytecode.NewLineInfo(8, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"set instance variable in a mixin singleton method": {
			input: `
				mixin Foo
					singleton
						var @foo: Int?

						def foo
				  		@foo = 2
						end
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(109, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(10, 2),
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
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(58, 6, 7), P(91, 8, 9)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(7, 3),
									bytecode.NewLineInfo(8, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"set instance variable in an interface singleton method": {
			input: `
				interface Foo
					singleton
						var @foo: Int?

						def foo
				  		@foo = 2
						end
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(113, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(10, 2),
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
						L(P(0, 1, 1), P(113, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(113, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(113, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(62, 6, 7), P(95, 8, 9)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(7, 3),
									bytecode.NewLineInfo(8, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"set instance variable in a class": {
			input: `
				class Foo
					singleton
						var @a: Int?
					end
					@a = 2
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
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
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
						L(P(0, 1, 1), P(77, 7, 8)),
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
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(77, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<class: Foo>"),
						[]byte{
							byte(bytecode.INT_2),
							byte(bytecode.DUP),
							byte(bytecode.SET_IVAR_0),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(76, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(6, 3),
							bytecode.NewLineInfo(7, 3),
						},
						nil,
					)),
				},
			),
		},
		"set instance variable in a mixin": {
			input: `
				mixin Foo
					singleton
						var @a: Int?
					end
					@a = 2
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
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
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
						L(P(0, 1, 1), P(77, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(77, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<mixin: Foo>"),
						[]byte{
							byte(bytecode.INT_2),
							byte(bytecode.DUP),
							byte(bytecode.SET_IVAR_0),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(76, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(6, 3),
							bytecode.NewLineInfo(7, 3),
						},
						nil,
					)),
				},
			),
		},
		"set instance variable in an interface": {
			input: `
				interface Foo
					singleton
						var @a: Int?
					end
					@a = 2
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
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(81, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 3,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(81, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(81, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<interface: Foo>"),
						[]byte{
							byte(bytecode.INT_2),
							byte(bytecode.DUP),
							byte(bytecode.SET_IVAR_0),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(80, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(6, 3),
							bytecode.NewLineInfo(7, 3),
						},
						nil,
					)),
				},
			),
		},
		"set instance variable in a module": {
			input: `
				module Foo
					var @a: Int?
					@a = 2
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
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
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
						L(P(0, 1, 1), P(53, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(53, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<module: Foo>"),
						[]byte{
							byte(bytecode.INT_2),
							byte(bytecode.DUP),
							byte(bytecode.SET_IVAR_0),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(52, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						nil,
					)),
				},
			),
		},

		"get box of instance variable in a class instance method": {
			input: `
				class Foo
					var @foo: Int?

					def foo: ^Int?
				  	&@foo
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(8, 2),
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
						L(P(0, 1, 1), P(85, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(85, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_0),
									byte(bytecode.CALL_METHOD8), 0,
									byte(bytecode.RETURN),
								},
								L(P(41, 5, 6), P(76, 7, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 3),
									bytecode.NewLineInfo(7, 1),
								},
								[]value.Value{
									value.Ref(&value.CallSiteInfo{
										Name:          value.ToSymbol("#box_of_ivar_index"),
										ArgumentCount: 1,
									}),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get box of instance variable from parent in a class instance method": {
			input: `
				class Bar
					var @foo: Int?
				end

				class Foo < Bar
					def foo: ^Int?
				  	&@foo
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(113, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(10, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
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
						L(P(0, 1, 1), P(113, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(113, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(113, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_0),
									byte(bytecode.CALL_METHOD8), 0,
									byte(bytecode.RETURN),
								},
								L(P(69, 7, 6), P(104, 9, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(8, 3),
									bytecode.NewLineInfo(9, 1),
								},
								[]value.Value{
									value.Ref(&value.CallSiteInfo{
										Name:          value.ToSymbol("#box_of_ivar_index"),
										ArgumentCount: 1,
									}),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get box of instance variable from mixin in a class instance method": {
			input: `
				mixin Bar
					var @foo: Int?
				end

				class Foo
					include Bar

					def foo: ^Int?
				  	&@foo
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(125, 12, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(12, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
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
						L(P(0, 1, 1), P(125, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(125, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(125, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_0),
									byte(bytecode.CALL_METHOD8), 0,
									byte(bytecode.RETURN),
								},
								L(P(81, 9, 6), P(116, 11, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(10, 3),
									bytecode.NewLineInfo(11, 1),
								},
								[]value.Value{
									value.Ref(&value.CallSiteInfo{
										Name:          value.ToSymbol("#box_of_ivar_index"),
										ArgumentCount: 1,
									}),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get box of instance variable in a mixin instance method": {
			input: `
				mixin Foo
					var @foo: Int?

					def foo: ^Int?
				  	&@foo
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
				L(P(0, 1, 1), P(85, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(8, 2),
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
						L(P(0, 1, 1), P(85, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
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
						L(P(0, 1, 1), P(85, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE_0),
									byte(bytecode.CALL_METHOD8), 1,
									byte(bytecode.RETURN),
								},
								L(P(41, 5, 6), P(76, 7, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 3),
									bytecode.NewLineInfo(7, 1),
								},
								[]value.Value{
									value.ToSymbol("foo").ToValue(),
									value.Ref(&value.CallSiteInfo{
										Name:          value.ToSymbol("#box_of_ivar_name"),
										ArgumentCount: 1,
									}),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},

		"get instance variable in a class instance method": {
			input: `
				class Foo
					var @foo: Int?

					def foo
				  	@foo
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(8, 2),
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
						L(P(0, 1, 1), P(77, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(77, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(77, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(41, 5, 6), P(68, 7, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 1),
									bytecode.NewLineInfo(7, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get instance variable from parent in a class instance method": {
			input: `
				class Bar
					var @foo: Int?
				end

				class Foo < Bar
					def foo
				  	@foo
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(10, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
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
						L(P(0, 1, 1), P(105, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(105, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(105, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(69, 7, 6), P(96, 9, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(8, 1),
									bytecode.NewLineInfo(9, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get instance variable from parent and self in a class instance method": {
			input: `
				class Bar
					var @bar: Int?
				end

				class Foo < Bar
					var @foo: Int?

					def foo
				  	@bar.must + @foo.must
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(143, 12, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(12, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
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
						L(P(0, 1, 1), P(143, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(143, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("bar"): 0,
							}),
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("bar"): 0,
								value.ToSymbol("foo"): 1,
							}),
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
						L(P(0, 1, 1), P(143, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.MUST),
									byte(bytecode.GET_IVAR_1),
									byte(bytecode.MUST),
									byte(bytecode.ADD_INT),
									byte(bytecode.RETURN),
								},
								L(P(90, 9, 6), P(134, 11, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(10, 5),
									bytecode.NewLineInfo(11, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get instance variable from mixin in a class instance method": {
			input: `
				mixin Bar
					var @foo: Int?
				end

				class Foo
					include Bar

					def foo
				  	@foo
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(117, 12, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(12, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
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
						L(P(0, 1, 1), P(117, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 20),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(117, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(117, 12, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(12, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(81, 9, 6), P(108, 11, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(10, 1),
									bytecode.NewLineInfo(11, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get instance variable in a mixin instance method": {
			input: `
				mixin Foo
					var @foo: Int?

					def foo
				  	@foo
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
				L(P(0, 1, 1), P(77, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(8, 2),
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
						L(P(0, 1, 1), P(77, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
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
						L(P(0, 1, 1), P(77, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_NAME16), 0, 0,
									byte(bytecode.RETURN),
								},
								L(P(41, 5, 6), P(68, 7, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 3),
									bytecode.NewLineInfo(7, 1),
								},
								[]value.Value{
									value.ToSymbol("foo").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get instance variable in a module method": {
			input: `
				module Foo
					var @foo: Int?

					def foo
				  	@foo
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(78, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(8, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(78, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(78, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(78, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(42, 5, 6), P(69, 7, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 1),
									bytecode.NewLineInfo(7, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get instance variable in a class method": {
			input: `
				class Foo
					singleton
						var @foo: Int?

						def foo
				  		@foo
						end
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(10, 2),
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
						L(P(0, 1, 1), P(105, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(105, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(105, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(58, 6, 7), P(87, 8, 9)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(7, 1),
									bytecode.NewLineInfo(8, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get instance variable in a mixin singleton method": {
			input: `
				mixin Foo
					singleton
						var @foo: Int?

						def foo
				  		@foo
						end
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(10, 2),
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
						L(P(0, 1, 1), P(105, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(105, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(105, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(58, 6, 7), P(87, 8, 9)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(7, 1),
									bytecode.NewLineInfo(8, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get instance variable in an interface singleton method": {
			input: `
				interface Foo
					singleton
						var @foo: Int?

						def foo
				  		@foo
						end
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(109, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(10, 2),
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
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("foo"): 0,
							}),
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
						L(P(0, 1, 1), P(109, 10, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(10, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(62, 6, 7), P(91, 8, 9)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(7, 1),
									bytecode.NewLineInfo(8, 1),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"get instance variable in a class": {
			input: `
				class Foo
					singleton
						var @a: Int?
					end
					@a
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
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
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
						L(P(0, 1, 1), P(73, 7, 8)),
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
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(73, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<class: Foo>"),
						[]byte{
							byte(bytecode.GET_IVAR_0),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(72, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(6, 1),
							bytecode.NewLineInfo(7, 3),
						},
						nil,
					)),
				},
			),
		},
		"get instance variable in a mixin": {
			input: `
				mixin Foo
					singleton
						var @a: Int?
					end
					@a
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
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
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
						L(P(0, 1, 1), P(73, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(73, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<mixin: Foo>"),
						[]byte{
							byte(bytecode.GET_IVAR_0),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(72, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(6, 1),
							bytecode.NewLineInfo(7, 3),
						},
						nil,
					)),
				},
			),
		},
		"get instance variable in an interface": {
			input: `
				interface Foo
					singleton
						var @a: Int?
					end
					@a
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
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 3,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(77, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(77, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<interface: Foo>"),
						[]byte{
							byte(bytecode.GET_IVAR_0),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(76, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(6, 1),
							bytecode.NewLineInfo(7, 3),
						},
						nil,
					)),
				},
			),
		},
		"get instance variable in a module": {
			input: `
				module Foo
					var @a: Int?
					@a
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
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(49, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
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
						L(P(0, 1, 1), P(49, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(49, 5, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<module: Foo>"),
						[]byte{
							byte(bytecode.GET_IVAR_0),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(48, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(4, 1),
							bytecode.NewLineInfo(5, 3),
						},
						nil,
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

func TestLocalVariables(t *testing.T) {
	tests := testTable{
		"declare": {
			input: "var a: Int",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{},
			),
		},
		"declare and initialise": {
			input: "var a = 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_3),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{},
			),
		},
		"declare with a pattern": {
			input: "var [1, a] = [1, 2]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS_NP), 0, 33,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.INT_2),
					byte(bytecode.EQUAL_INT),
					byte(bytecode.JUMP_UNLESS_NP), 0, 24,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.INT_0),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.INT_1),
					byte(bytecode.EQUAL),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 13,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_IF), 0, 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.THROW),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 49),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{value.SmallInt(1).ToValue(), value.SmallInt(2).ToValue()}),
					value.Ref(value.ListMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.Ref(value.NewError(value.PatternNotMatchedErrorClass, "assigned value does not match the pattern defined in variable declaration")),
				},
			),
		},

		"read undeclared": {
			input: "a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(0, 1, 1)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				nil,
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(0, 1, 1)), "undefined local `a`"),
			},
		},
		"assign undeclared": {
			input: "a = 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(0, 1, 1)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.SmallInt(3).ToValue(),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(0, 1, 1)), "undefined local `a`"),
			},
		},
		"assign uninitialised": {
			input: `
				var a: String
				a = 'foo'
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
				},
			),
		},
		"assign initialised": {
			input: `
				var a = 'foo'
				a = 'bar'
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
					value.Ref(value.String("bar")),
				},
			),
		},
		"read uninitialised": {
			input: `
				var a: Int
				a + 2
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 3, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 6),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(2).ToValue(),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"read initialised": {
			input: `
				var a = 5
				a + 2
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_2),
					byte(bytecode.ADD_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{},
			),
		},
		"read initialised in child scope": {
			input: `
				var a = 5
				do
					a + 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_2),
					byte(bytecode.ADD_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{},
			),
		},
		"shadow in child scope": {
			input: `
				var a = 5
				2 + do
					var a = 10
					a + 12
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_2),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.LOAD_INT_8), 12,
					byte(bytecode.ADD_INT),
					byte(bytecode.ADD_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestUpvalues(t *testing.T) {
	tests := testTable{
		"in child scope": {
			input: `
				upvalue := 5
				-> println upvalue
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.CLOSURE), 2, 1, 0xff,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 6),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionWithUpvalues(
						functionSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.UNDEFINED),
							byte(bytecode.GET_UPVALUE_0),
							byte(bytecode.NEW_ARRAY_TUPLE8), 1,
							byte(bytecode.CALL_METHOD_TCO8), 1,
							byte(bytecode.RETURN),
						},
						L(P(22, 3, 5), P(39, 3, 22)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 9),
						},
						0,
						0,
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
						},
						1,
					)),
				},
			),
		},
		"close upvalue": {
			input: `
				do
					upvalue := 5
					-> println upvalue
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.CLOSURE), 2, 1, 0xff,
					byte(bytecode.CLOSE_UPVALUE_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(57, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionWithUpvalues(
						functionSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.UNDEFINED),
							byte(bytecode.GET_UPVALUE_0),
							byte(bytecode.NEW_ARRAY_TUPLE8), 1,
							byte(bytecode.CALL_METHOD_TCO8), 1,
							byte(bytecode.RETURN),
						},
						L(P(31, 4, 6), P(48, 4, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(4, 9),
						},
						0,
						0,
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
						},
						1,
					)),
				},
			),
		},
		"close upvalues in multiple scopes": {
			input: `
				a := 2
				b := 10

				if b == 40
					g := 5
					e := 9
					f := 20
					h := -> g + f
					if a == 20
						a := 10
						b := 50
						c := 30
						d := -> a + b + c
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 10,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.LOAD_INT_8), 40,
					byte(bytecode.JUMP_UNLESS_IEQ), 0, 58,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.LOAD_INT_8), 9,
					byte(bytecode.SET_LOCAL_4),
					byte(bytecode.LOAD_INT_8), 20,
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.CLOSURE),
					2, 3,
					2, 5,
					0xff,
					byte(bytecode.SET_LOCAL8), 6,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 20,
					byte(bytecode.JUMP_UNLESS_IEQ), 0, 29,
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.SET_LOCAL8), 7,
					byte(bytecode.LOAD_INT_8), 50,
					byte(bytecode.SET_LOCAL8), 8,
					byte(bytecode.LOAD_INT_8), 30,
					byte(bytecode.SET_LOCAL8), 9,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.CLOSURE),
					2, 7,
					2, 8,
					2, 9,
					0xff,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL8), 10,
					byte(bytecode.CLOSE_UPVALUE8), 7,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.CLOSE_UPVALUE_3),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(194, 16, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(5, 6),
					bytecode.NewLineInfo(6, 2),
					bytecode.NewLineInfo(7, 3),
					bytecode.NewLineInfo(8, 4),
					bytecode.NewLineInfo(9, 9),
					bytecode.NewLineInfo(10, 6),
					bytecode.NewLineInfo(11, 4),
					bytecode.NewLineInfo(12, 4),
					bytecode.NewLineInfo(13, 4),
					bytecode.NewLineInfo(14, 12),
					bytecode.NewLineInfo(10, 6),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(16, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionWithUpvalues(
						functionSymbol,
						[]byte{
							byte(bytecode.GET_UPVALUE_0),
							byte(bytecode.GET_UPVALUE_1),
							byte(bytecode.ADD_INT),
							byte(bytecode.RETURN),
						},
						L(P(87, 9, 11), P(94, 9, 18)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(9, 4),
						},
						0,
						0,
						nil,
						2,
					)),
					value.Ref(vm.NewBytecodeFunctionWithUpvalues(
						functionSymbol,
						[]byte{
							byte(bytecode.GET_UPVALUE_0),
							byte(bytecode.GET_UPVALUE_1),
							byte(bytecode.ADD_INT),
							byte(bytecode.GET_UPVALUE8), 2,
							byte(bytecode.ADD_INT),
							byte(bytecode.RETURN),
						},
						L(P(165, 14, 12), P(176, 14, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(14, 7),
						},
						0,
						0,
						nil,
						3,
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

func TestLocalValues(t *testing.T) {
	tests := testTable{
		"declare": {
			input: "val a: Int",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{},
			),
		},
		"declare and initialise": {
			input: "val a = 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_3),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{},
			),
		},
		"declare with a pattern": {
			input: "val [1, a] = [1, 2]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS_NP), 0, 33,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.INT_2),
					byte(bytecode.EQUAL_INT),
					byte(bytecode.JUMP_UNLESS_NP), 0, 24,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.INT_0),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.INT_1),
					byte(bytecode.EQUAL),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 13,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_IF), 0, 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.THROW),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 49),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{value.SmallInt(1).ToValue(), value.SmallInt(2).ToValue()}),
					value.Ref(value.ListMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.Ref(value.NewError(value.PatternNotMatchedErrorClass, "assigned value does not match the pattern defined in value declaration")),
				},
			),
		},
		"assign uninitialised": {
			input: `
				val a: String
				a = 'foo'
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
				},
			),
		},
		"assign initialised": {
			input: `
				val a = 'foo'
				a = 'bar'
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(23, 3, 5), P(23, 3, 5)), "local value `a` cannot be reassigned"),
				diagnostic.NewFailure(L(P(27, 3, 9), P(31, 3, 13)), "type `\"bar\"` cannot be assigned to type `\"foo\"`"),
			},
		},
		"read uninitialised": {
			input: `
				val a: Int
				a + 2
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"read initialised": {
			input: `
				val a = 5
				a + 2
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_2),
					byte(bytecode.ADD_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{},
			),
		},
		"read initialised in child scope": {
			input: `
				val a = 5
				do
					a + 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_2),
					byte(bytecode.ADD_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{},
			),
		},
		"shadow in child scope": {
			input: `
				val a = 5
				2 + do
					val a = 10
					a + 12
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_2),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.LOAD_INT_8), 12,
					byte(bytecode.ADD_INT),
					byte(bytecode.ADD_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{},
			),
		},
		"declare and initialise 255 variables": {
			input: `
				do
					a0:=1;a1:=1;a2:=1;a3:=1;a4:=1;a5:=1;a6:=1;a7:=1;a8:=1;a9:=1;a10:=1;a11:=1;a12:=1;a13:=1;a14:=1;a15:=1;a16:=1;a17:=1;a18:=1;a19:=1;a20:=1;a21:=1;a22:=1;a23:=1;a24:=1;a25:=1;a26:=1;a27:=1;a28:=1;a29:=1;a30:=1;a31:=1;a32:=1;a33:=1;a34:=1;a35:=1;a36:=1;a37:=1;a38:=1;a39:=1;a40:=1;a41:=1;a42:=1;a43:=1;a44:=1;a45:=1;a46:=1;a47:=1;a48:=1;a49:=1;a50:=1;a51:=1;a52:=1;a53:=1;a54:=1;a55:=1;a56:=1;a57:=1;a58:=1;a59:=1;a60:=1;a61:=1;a62:=1;a63:=1;a64:=1;a65:=1;a66:=1;a67:=1;a68:=1;a69:=1;a70:=1;a71:=1;a72:=1;a73:=1;a74:=1;a75:=1;a76:=1;a77:=1;a78:=1;a79:=1;a80:=1;a81:=1;a82:=1;a83:=1;a84:=1;a85:=1;a86:=1;a87:=1;a88:=1;a89:=1;a90:=1;a91:=1;a92:=1;a93:=1;a94:=1;a95:=1;a96:=1;a97:=1;a98:=1;a99:=1;a100:=1;a101:=1;a102:=1;a103:=1;a104:=1;a105:=1;a106:=1;a107:=1;a108:=1;a109:=1;a110:=1;a111:=1;a112:=1;a113:=1;a114:=1;a115:=1;a116:=1;a117:=1;a118:=1;a119:=1;a120:=1;a121:=1;a122:=1;a123:=1;a124:=1;a125:=1;a126:=1;a127:=1;a128:=1;a129:=1;a130:=1;a131:=1;a132:=1;a133:=1;a134:=1;a135:=1;a136:=1;a137:=1;a138:=1;a139:=1;a140:=1;a141:=1;a142:=1;a143:=1;a144:=1;a145:=1;a146:=1;a147:=1;a148:=1;a149:=1;a150:=1;a151:=1;a152:=1;a153:=1;a154:=1;a155:=1;a156:=1;a157:=1;a158:=1;a159:=1;a160:=1;a161:=1;a162:=1;a163:=1;a164:=1;a165:=1;a166:=1;a167:=1;a168:=1;a169:=1;a170:=1;a171:=1;a172:=1;a173:=1;a174:=1;a175:=1;a176:=1;a177:=1;a178:=1;a179:=1;a180:=1;a181:=1;a182:=1;a183:=1;a184:=1;a185:=1;a186:=1;a187:=1;a188:=1;a189:=1;a190:=1;a191:=1;a192:=1;a193:=1;a194:=1;a195:=1;a196:=1;a197:=1;a198:=1;a199:=1;a200:=1;a201:=1;a202:=1;a203:=1;a204:=1;a205:=1;a206:=1;a207:=1;a208:=1;a209:=1;a210:=1;a211:=1;a212:=1;a213:=1;a214:=1;a215:=1;a216:=1;a217:=1;a218:=1;a219:=1;a220:=1;a221:=1;a222:=1;a223:=1;a224:=1;a225:=1;a226:=1;a227:=1;a228:=1;a229:=1;a230:=1;a231:=1;a232:=1;a233:=1;a234:=1;a235:=1;a236:=1;a237:=1;a238:=1;a239:=1;a240:=1;a241:=1;a242:=1;a243:=1;a244:=1;a245:=1;a246:=1;a247:=1;a248:=1;a249:=1;a250:=1;a251:=1;a252:=1;a253:=1;a254:=1;a255:=1
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				append(
					append(
						[]byte{
							byte(bytecode.PREP_LOCALS16), 1, 0,
							byte(bytecode.INT_1),
							byte(bytecode.SET_LOCAL_1),
							byte(bytecode.INT_1),
							byte(bytecode.SET_LOCAL_2),
							byte(bytecode.INT_1),
							byte(bytecode.SET_LOCAL_3),
							byte(bytecode.INT_1),
							byte(bytecode.SET_LOCAL_4),
						},
						declareNVariables(253)...,
					),
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL8), 254,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL8), 255,
					byte(bytecode.INT_1),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL16), 1, 0,
					byte(bytecode.RETURN),
				),
				L(P(0, 1, 1), P(1958, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(3, 766),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func declareNVariables(n int) []byte {
	b := make([]byte, 0, n*4)
	for i := 5; i <= n; i++ {
		b = append(
			b,
			byte(bytecode.INT_1),
			byte(bytecode.SET_LOCAL8), byte(i),
		)
	}

	return b
}
