package compiler_test

import (
	"testing"
)

func TestGoComplexAssignmentLocals(t *testing.T) {
	tests := goTestTable{
		// TODO:
		// 		"increment": {
		// 			input: "a := 1; a++",
		// 			want: `
		// `,
		// 		},
		// "decrement": {
		// 	input: "a := 1; a--",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.DECREMENT_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(10, 1, 11)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 9),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "add": {
		// 	input: "a := 1; a += 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.ADD_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "subtract": {
		// 	input: "a := 1; a -= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.SUBTRACT_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "multiply": {
		// 	input: "a := 1; a *= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.MULTIPLY_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "divide": {
		// 	input: "a := 1; a /= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.DIVIDE_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "exponentiate": {
		// 	input: "a := 1; a **= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.EXPONENTIATE_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(14, 1, 15)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "modulo": {
		// 	input: "a := 1; a %= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.MODULO_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "bitwise AND": {
		// 	input: "a := 1; a &= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.BITWISE_AND_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "bitwise OR": {
		// 	input: "a := 1; a |= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.BITWISE_OR_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "bitwise XOR": {
		// 	input: "a := 1; a ^= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.BITWISE_XOR_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "left bitshift": {
		// 	input: "a := 1; a <<= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.LBITSHIFT_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(14, 1, 15)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "left logical bitshift": {
		// 	input: "a := 1u64; a <<<= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.LOAD_UINT64_8), 1,
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.LOGIC_LBITSHIFT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(18, 1, 19)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 11),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "right bitshift": {
		// 	input: "a := 1; a >>= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.RBITSHIFT_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(14, 1, 15)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "right logical bitshift": {
		// 	input: "a := 1u64; a >>>= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.LOAD_UINT64_8), 1,
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.LOGIC_RBITSHIFT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(18, 1, 19)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 11),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "logic OR": {
		// 	input: "var a: Int? = 1; a ||= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.JUMP_IF_NP), 0, 2,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(23, 1, 24)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 13),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "logic AND": {
		// 	input: "var a: Int? = 1; a &&= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.JUMP_UNLESS_NP), 0, 2,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(23, 1, 24)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 13),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "nil coalesce": {
		// 	input: "var a: Int? = 1; a ??= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.JUMP_UNLESS_NNP), 0, 2,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(23, 1, 24)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 13),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

// func TestBytecodeComplexAssignmentInstanceVariables(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"increment": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a++
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(82, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_NAME16), 0, 0,
// 									byte(bytecode.INCREMENT_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(73, 6, 22)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 7),
// 								},
// 								[]value.Value{
// 									value.ToSymbol("a").ToValue(),
// 								},
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"decrement": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a--
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(82, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_NAME16), 0, 0,
// 									byte(bytecode.DECREMENT_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(73, 6, 22)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 7),
// 								},
// 								[]value.Value{
// 									value.ToSymbol("a").ToValue(),
// 								},
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"add": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a += 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.ADD_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"subtract": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a -= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.SUBTRACT_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"multiply": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a *= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.MULTIPLY_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"divide": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a /= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.DIVIDE_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"exponentiate": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a **= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(86, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.EXPONENTIATE_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(77, 6, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"modulo": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a %= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.MODULO_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"bitwise AND": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a &= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.BITWISE_AND_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"bitwise OR": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a |= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.BITWISE_OR_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"bitwise XOR": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a ^= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.BITWISE_XOR_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"left bitshift": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a <<= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(86, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.LBITSHIFT_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(77, 6, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"left logical bitshift": {
// 			input: `
// 				class Foo
// 					var @a: UInt64
// 					init(@a); end

// 					def foo then @a <<<= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(90, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(90, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(90, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(90, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(40, 4, 6), P(52, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.LOGIC_LBITSHIFT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(60, 6, 6), P(81, 6, 27)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"right bitshift": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a >>= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(86, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.RBITSHIFT_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(77, 6, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"right logical bitshift": {
// 			input: `
// 				class Foo
// 					var @a: Int64
// 					init(@a); end

// 					def foo then @a >>>= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(89, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(89, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(89, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(89, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(39, 4, 6), P(51, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.LOGIC_RBITSHIFT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(59, 6, 6), P(80, 6, 27)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"logic OR": {
// 			input: `
// 				class Foo
// 					var @a: Int?

// 					def foo then @a ||= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(68, 6, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(6, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 6),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.JUMP_IF_NP), 0, 2,
// 									byte(bytecode.POP),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(39, 5, 6), P(59, 5, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(5, 9),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"nil coalesce": {
// 			input: `
// 				class Foo
// 					var @a: Int?

// 					def foo then @a ??= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(68, 6, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(6, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 6),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.JUMP_UNLESS_NNP), 0, 2,
// 									byte(bytecode.POP),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(39, 5, 6), P(59, 5, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(5, 9),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"logic AND": {
// 			input: `
// 				class Foo
// 					var @a: Int?

// 					def foo then @a &&= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(68, 6, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(6, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 6),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.JUMP_UNLESS_NP), 0, 2,
// 									byte(bytecode.POP),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(39, 5, 6), P(59, 5, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(5, 9),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			bytecodeCompilerTest(tc, t)
// 		})
// 	}
// }
