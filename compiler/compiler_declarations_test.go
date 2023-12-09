package compiler

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

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
					value.SymbolTable.Add("Foo"),
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
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"named class inside osf a method": {
			input: `
				def foo
				  class ::Bar; end
				end
			`,
			err: errors.ErrorList{
				errors.NewError(L(P(19, 3, 7), P(34, 3, 22)), "can't define named classes inside of a method: foo"),
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
					value.SymbolTable.Add("Bar"),
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
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 1, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
					value.SymbolTable.Add("Bar"),
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
					value.SymbolTable.Add("Baz"),
					value.SymbolTable.Add("Bar"),
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
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 1, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
					value.SymbolTable.Add("Baz"),
					value.SymbolTable.Add("Bar"),
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
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
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
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 1, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SymbolTable.Add("Std"),
					value.SymbolTable.Add("Int"),
					value.SymbolTable.Add("Foo"),
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
					byte(bytecode.DEF_CLASS),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(44, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
					value.SymbolTable.Add("Foo"),
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
					byte(bytecode.DEF_CLASS),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(24, 1, 25)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
						},
						[]value.Value{
							vm.NewBytecodeMethodNoParams(
								value.SymbolTable.Add("a"),
								[]byte{
									byte(bytecode.RETURN),
								},
								L(P(13, 1, 14), P(24, 1, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(1, 1),
								},
								nil,
							),
							value.SymbolTable.Add("a"),
						},
					),
					value.SymbolTable.Add("A"),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(40, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
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
							byte(bytecode.DEF_CLASS),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(70, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
							bytecode.NewLineInfo(7, 3),
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
									byte(bytecode.SELF),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(62, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 3),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
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
					value.SymbolTable.Add("Foo"),
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
				errors.NewError(L(P(18, 3, 6), P(32, 3, 20)), "can't define named modules inside of a method: foo"),
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
					value.SymbolTable.Add("Foo"),
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
					value.SymbolTable.Add("Std"),
					value.SymbolTable.Add("Int"),
					value.SymbolTable.Add("Foo"),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(41, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(45, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
					value.SymbolTable.Add("Foo"),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(72, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(7, 3),
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
									byte(bytecode.SELF),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(64, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 3),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
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
						value.SymbolTable.Add("foo"),
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(21, 2, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(2, 2),
						},
						[]value.Value{
							value.SymbolTable.Add("bar"),
						},
					),
					value.SymbolTable.Add("foo"),
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
						L(P(5, 2, 5), P(52, 5, 7)),
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
					value.SymbolTable.Add("foo"),
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
						value.SymbolTable.Add("foo"),
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
							value.SymbolTable.Add("a"),
							value.SymbolTable.Add("b"),
							value.SymbolTable.Add("c"),
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
					value.SymbolTable.Add("foo"),
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
					byte(bytecode.DEF_CLASS),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(78, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							vm.NewBytecodeMethod(
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
								L(P(20, 3, 6), P(70, 6, 8)),
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
							value.SymbolTable.Add("foo"),
						},
					),
					value.SymbolTable.Add("Bar"),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(79, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							vm.NewBytecodeMethod(
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
								L(P(21, 3, 6), P(71, 6, 8)),
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
							value.SymbolTable.Add("foo"),
						},
					),
					value.SymbolTable.Add("Bar"),
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
					value.SymbolTable.Add("Foo"),
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
					value.SymbolTable.Add("Foo"),
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
				errors.NewError(L(P(18, 3, 6), P(31, 3, 19)), "can't define named mixins inside of a method: foo"),
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
					value.SymbolTable.Add("Std"),
					value.SymbolTable.Add("Int"),
					value.SymbolTable.Add("Foo"),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(40, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(44, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
					value.SymbolTable.Add("Foo"),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(70, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(7, 3),
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
									byte(bytecode.SELF),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(62, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 3),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
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
					byte(bytecode.DEF_CLASS),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(40, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
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
					byte(bytecode.DEF_CLASS),
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

							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(47, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 9),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
							value.SymbolTable.Add("Baz"),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"include in top level": {
			input: `include ::Bar`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(0, 1, 1), P(12, 1, 13)),
					"can't include mixins in the top level",
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
					"can't include mixins in a module",
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
					"can't include mixins in a method",
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
					byte(bytecode.DEF_CLASS),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(39, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 6),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(40, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 6),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
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
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(39, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 6),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
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
					byte(bytecode.DEF_CLASS),
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

							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(46, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 11),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
							value.SymbolTable.Add("Baz"),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"extend in top level": {
			input: `extend ::Bar`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(0, 1, 1), P(11, 1, 12)),
					"can't extend mixins in the top level",
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
					"can't extend mixins in a method",
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
