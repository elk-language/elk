package compiler_test

import (
	"testing"
)

func TestGoGetConstant(t *testing.T) {
	tests := goTestTable{
		"builtin absolute path ::Std": {
			input: "a := ::Std",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.StdModule).ToValue()
}
`,
		},
		"custom absolute path": {
			input: `
				module Foo; end
				a := ::Foo
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (const0).ToValue()
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewModule()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}
`,
		},
		// "absolute nested path ::Std::Float::INF": {
		// 	input: "::Std::Float::INF",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(16, 1, 17)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 3),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Float::INF").ToValue(),
		// 		},
		// 	),
		// },
		// "relative path from using": {
		// 	input: `
		// 		using Std::Float::INF as I
		// 		I
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(37, 3, 6)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(3, 3),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Float::INF").ToValue(),
		// 		},
		// 	),
		// },
		// "relative path": {
		// 	input: `
		// 		module Foo
		// 			const BAR = 3
		// 			module Baz
		// 				println BAR
		// 			end
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.LOAD_VALUE_0),
		// 			byte(bytecode.EXEC),
		// 			byte(bytecode.POP),
		// 			byte(bytecode.GET_CONST8), 1,
		// 			byte(bytecode.LOAD_VALUE_2),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.DEF_CONST),
		// 			byte(bytecode.GET_CONST8), 1,
		// 			byte(bytecode.LOAD_VALUE_3),
		// 			byte(bytecode.INIT_NAMESPACE),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(85, 7, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 3),
		// 			bytecode.NewLineInfo(3, 5),
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(7, 1),
		// 		},
		// 		[]value.Value{
		// 			value.Ref(vm.NewBytecodeFunctionNoParams(
		// 				nil,
		// 				namespaceDefinitionsSymbol,
		// 				[]byte{
		// 					byte(bytecode.GET_CONST8), 0,
		// 					byte(bytecode.LOAD_VALUE_1),
		// 					byte(bytecode.DEF_NAMESPACE), 0,
		// 					byte(bytecode.GET_CONST8), 1,
		// 					byte(bytecode.LOAD_VALUE_2),
		// 					byte(bytecode.DEF_NAMESPACE), 0,
		// 					byte(bytecode.NIL),
		// 					byte(bytecode.RETURN),
		// 				},
		// 				L(P(0, 1, 1), P(85, 7, 8)),
		// 				bytecode.LineInfoList{
		// 					bytecode.NewLineInfo(1, 10),
		// 					bytecode.NewLineInfo(7, 2),
		// 				},
		// 				[]value.Value{
		// 					value.ToSymbol("Root").ToValue(),
		// 					value.ToSymbol("Foo").ToValue(),
		// 					value.ToSymbol("Baz").ToValue(),
		// 				},
		// 			)),
		// 			value.ToSymbol("Foo").ToValue(),
		// 			value.ToSymbol("BAR").ToValue(),
		// 			value.Ref(vm.NewBytecodeFunctionNoParams(
		// 				nil,
		// 				value.ToSymbol("<module: Foo>"),
		// 				[]byte{
		// 					byte(bytecode.GET_CONST8), 0,
		// 					byte(bytecode.LOAD_VALUE_1),
		// 					byte(bytecode.INIT_NAMESPACE),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.NIL),
		// 					byte(bytecode.RETURN),
		// 				},
		// 				L(P(5, 2, 5), P(84, 7, 7)),
		// 				bytecode.LineInfoList{
		// 					bytecode.NewLineInfo(4, 4),
		// 					bytecode.NewLineInfo(7, 3),
		// 				},
		// 				[]value.Value{
		// 					value.ToSymbol("Foo::Baz").ToValue(),
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						nil,
		// 						value.ToSymbol("<module: Foo::Baz>"),
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.GET_CONST8), 1,
		// 							byte(bytecode.CALL_METHOD_NT8), 2,
		// 							byte(bytecode.POP),
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(40, 4, 6), P(76, 6, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(5, 6),
		// 							bytecode.NewLineInfo(6, 3),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Std::Kernel").ToValue(),
		// 							value.ToSymbol("Foo::BAR").ToValue(),
		// 							value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "println@1"), 1)),
		// 						},
		// 					)),
		// 				},
		// 			)),
		// 		},
		// 	),
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

// func TestBytecodeDefConstant(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"relative path Foo": {
// 			input: "const Foo = 3",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				nil,
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.GET_CONST8), 0,
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.INT_3),
// 					byte(bytecode.DEF_CONST),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(12, 1, 13)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 7),
// 				},
// 				[]value.Value{
// 					value.ToSymbol("Root").ToValue(),
// 					value.ToSymbol("Foo").ToValue(),
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
