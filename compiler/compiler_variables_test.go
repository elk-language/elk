package compiler

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestInstanceVariables(t *testing.T) {
	tests := testTable{
		"initialise when declared": {
			input: "var @a: Int = 3",
			err: error.ErrorList{
				error.NewFailure(
					L(P(14, 1, 15), P(14, 1, 15)),
					"instance variables cannot be initialised when declared",
				),
			},
		},
		"declare in the top level": {
			input: "var @a: Float",
			err: error.ErrorList{
				error.NewFailure(
					L(P(0, 1, 1), P(12, 1, 13)),
					"instance variables can only be declared in class, module, mixin bodies",
				),
			},
		},
		"declare in a method": {
			input: "def foo; var @a: Float; end",
			err: error.ErrorList{
				error.NewFailure(
					L(P(9, 1, 10), P(21, 1, 22)),
					"instance variables can only be declared in class, module, mixin bodies",
				),
			},
		},
		"declare in a class": {
			input: "class Foo; var @a: Float; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 1, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(0, 1, 1), P(28, 1, 29)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 3),
						},
						nil,
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"declare in a mixin": {
			input: "mixin Foo; var @a: Float; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 1, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(0, 1, 1), P(28, 1, 29)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 3),
						},
						nil,
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"declare in a module": {
			input: "module Foo; var @a: Float; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 1, 30)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(0, 1, 1), P(29, 1, 30)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 3),
						},
						nil,
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"set instance variable in top level": {
			input: "@a = 2",
			err: error.ErrorList{
				error.NewFailure(
					L(P(0, 1, 1), P(5, 1, 6)),
					"instance variables cannot be set in the top level",
				),
			},
		},
		"set instance variable in a method": {
			input: `
				def foo
				  @foo = 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("foo"),
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_IVAR8), 1,
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(34, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 1),
						},
						[]value.Value{
							value.SmallInt(2),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("foo"),
				},
			),
		},
		"set instance variable in a class": {
			input: `
				class Foo
					@a = 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_IVAR8), 1,
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(33, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.SmallInt(2),
							value.ToSymbol("a"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"set instance variable in a mixin": {
			input: `
				mixin Foo
					@a = 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_IVAR8), 1,
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(33, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.SmallInt(2),
							value.ToSymbol("a"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"set instance variable in a module": {
			input: `
				module Foo
					@a = 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_IVAR8), 1,
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(34, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.SmallInt(2),
							value.ToSymbol("a"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"read instance variable in top level": {
			input: "@a",
			err: error.ErrorList{
				error.NewFailure(
					L(P(0, 1, 1), P(1, 1, 2)),
					"cannot read instance variables in the top level",
				),
			},
		},
		"read instance variable in a method": {
			input: `
				def foo then @a
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 2, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("foo"),
						[]byte{
							byte(bytecode.GET_IVAR8), 0,
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(19, 2, 19)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(2, 3),
						},
						[]value.Value{
							value.ToSymbol("a"),
						},
					),
					value.ToSymbol("foo"),
				},
			),
		},
		"read instance variable in a mixin": {
			input: `mixin Foo then @a`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.GET_IVAR8), 0,
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(0, 1, 1), P(16, 1, 17)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
						},
						[]value.Value{
							value.ToSymbol("a"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"read instance variable in a module": {
			input: `module Foo then @a`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.GET_IVAR8), 0,
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(0, 1, 1), P(17, 1, 18)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
						},
						[]value.Value{
							value.ToSymbol("a"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"read instance variable in a class": {
			input: `class Foo then @a`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.GET_IVAR8), 0,
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(0, 1, 1), P(16, 1, 17)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
						},
						[]value.Value{
							value.ToSymbol("a"),
						},
					),
					value.ToSymbol("Foo"),
				},
			),
		},
		"define with required parameters in a class": {
			input: `
				class Bar
					init then @a
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
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
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.RETURN_SELF),
						},
						L(P(5, 2, 5), P(39, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(20, 3, 6), P(31, 3, 17)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 4),
								},
								nil,
								0,
								-1,
								false,
								false,
								[]value.Value{
									value.ToSymbol("a"),
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

func TestLocalVariables(t *testing.T) {
	tests := testTable{
		"declare": {
			input: "var a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				nil,
			),
		},
		"declare with a type": {
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
				nil,
			),
		},
		"declare and initialise": {
			input: "var a = 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.SmallInt(3),
				},
			),
		},
		"declare with a pattern": {
			input: "var [1, a] = foo()",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CALL_SELF8), 0,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 37,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 27,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.EQUAL),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 14,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_IF), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.THROW),
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 57),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.ListMixin,
					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
					value.SmallInt(2),
					value.SmallInt(0),
					value.SmallInt(1),
					value.NewError(value.PatternNotMatchedErrorClass, "assigned value does not match the pattern defined in variable declaration"),
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
			err: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(0, 1, 1)), "undeclared variable: a"),
			},
		},
		"assign undeclared": {
			input: "a = 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.SmallInt(3),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(4, 1, 5)), "undeclared variable: a"),
			},
		},
		"assign uninitialised": {
			input: `
				var a
				a = 'foo'
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					value.String("foo"),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					value.String("foo"),
					value.String("bar"),
				},
			),
		},
		"read uninitialised": {
			input: `
				var a
				a + 2
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					value.SmallInt(2),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L(P(15, 3, 5), P(15, 3, 5)), "cannot access an uninitialised local: a"),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 6),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
					value.SmallInt(10),
					value.SmallInt(12),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.CLOSURE), 2, 3, 0xff,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.SmallInt(5),
					vm.NewBytecodeFunctionWithUpvalues(
						functionSymbol,
						[]byte{
							byte(bytecode.GET_UPVALUE8), 0,
							byte(bytecode.CALL_SELF8), 0,
							byte(bytecode.RETURN),
						},
						L(P(22, 3, 5), P(39, 3, 22)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
						},
						nil,
						0,
						-1,
						false,
						false,
						[]value.Value{
							value.NewCallSiteInfo(value.ToSymbol("println"), 1, nil),
						},
						1,
					),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.CLOSURE), 2, 3, 0xff,
					byte(bytecode.CLOSE_UPVALUE8), 3,
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(57, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 6),
					bytecode.NewLineInfo(5, 6),
				},
				[]value.Value{
					value.SmallInt(5),
					vm.NewBytecodeFunctionWithUpvalues(
						functionSymbol,
						[]byte{
							byte(bytecode.GET_UPVALUE8), 0,
							byte(bytecode.CALL_SELF8), 0,
							byte(bytecode.RETURN),
						},
						L(P(31, 4, 6), P(48, 4, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(4, 5),
						},
						nil,
						0,
						-1,
						false,
						false,
						[]value.Value{
							value.NewCallSiteInfo(value.ToSymbol("println"), 1, nil),
						},
						1,
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

func TestLocalValues(t *testing.T) {
	tests := testTable{
		"declare": {
			input: "val a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				nil,
			),
		},
		"declare with a type": {
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
				nil,
			),
		},
		"declare and initialise": {
			input: "val a = 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.SmallInt(3),
				},
			),
		},
		"declare with a pattern": {
			input: "val [1, a] = foo()",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CALL_SELF8), 0,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 37,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 27,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.EQUAL),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 14,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_IF), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.THROW),
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 57),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.ListMixin,
					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
					value.SmallInt(2),
					value.SmallInt(0),
					value.SmallInt(1),
					value.NewError(value.PatternNotMatchedErrorClass, "assigned value does not match the pattern defined in value declaration"),
				},
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
							byte(bytecode.PREP_LOCALS16),
							1,
							0,
						},
						declareNVariables(253)...,
					),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL16), 1, 0,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL16), 1, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL16), 1, 2,
					byte(bytecode.LEAVE_SCOPE32), 1, 2, 1, 0,
					byte(bytecode.RETURN),
				),
				L(P(0, 1, 1), P(1958, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 1285),
					bytecode.NewLineInfo(4, 6),
				},
				[]value.Value{
					value.SmallInt(1),
				},
			),
		},
		"assign uninitialised": {
			input: `
				val a
				a = 'foo'
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					value.String("foo"),
				},
			),
		},
		"assign initialised": {
			input: `
				val a = 'foo'
				a = 'bar'
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					value.String("foo"),
					value.String("bar"),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L(P(23, 3, 5), P(31, 3, 13)), "cannot reassign a val: a"),
			},
		},
		"read uninitialised": {
			input: `
				val a
				a + 2
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					value.SmallInt(2),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L(P(15, 3, 5), P(15, 3, 5)), "cannot access an uninitialised local: a"),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 6),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
					value.SmallInt(10),
					value.SmallInt(12),
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

func declareNVariables(n int) []byte {
	b := make([]byte, 0, n*4)
	for i := 0; i < n; i++ {
		b = append(
			b,
			byte(bytecode.LOAD_VALUE8), 0,
			byte(bytecode.SET_LOCAL8), byte(i+3),
			byte(bytecode.POP),
		)
	}

	return b
}
