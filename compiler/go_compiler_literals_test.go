package compiler_test

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

// func TestClosureLiteral(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"recursive closure": {
// 			input: `
// 				var calc_fib: |n: Int|: Int = |n| ->
// 					return 1 if n < 3

// 					calc_fib(n - 2) + calc_fib(n - 1)
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.CLOSURE), 2, 1, 0xff,
// 					byte(bytecode.DUP),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(112, 6, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 2),
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(6, 1),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionWithUpvalues(
// 						functionSymbol,
// 						[]byte{
// 							byte(bytecode.GET_LOCAL_1),
// 							byte(bytecode.INT_3),
// 							byte(bytecode.JUMP_UNLESS_ILT), 0, 2,
// 							byte(bytecode.INT_1),
// 							byte(bytecode.RETURN),
// 							byte(bytecode.GET_UPVALUE_0),
// 							byte(bytecode.GET_LOCAL_1),
// 							byte(bytecode.INT_2),
// 							byte(bytecode.SUBTRACT_INT),
// 							byte(bytecode.CALL8), 0,
// 							byte(bytecode.GET_UPVALUE_0),
// 							byte(bytecode.GET_LOCAL_1),
// 							byte(bytecode.INT_1),
// 							byte(bytecode.SUBTRACT_INT),
// 							byte(bytecode.CALL8), 1,
// 							byte(bytecode.ADD_INT),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(35, 2, 35), P(111, 6, 7)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(3, 7),
// 							bytecode.NewLineInfo(5, 13),
// 							bytecode.NewLineInfo(6, 1),
// 						},
// 						1,
// 						0,
// 						[]value.Value{
// 							value.Ref(value.NewCallSiteInfo(value.ToSymbol("call"), 1)),
// 							value.Ref(value.NewCallSiteInfo(value.ToSymbol("call"), 1)),
// 						},
// 						1,
// 					)),
// 				},
// 			),
// 		},
// 		"lambda": {
// 			input: `
// 				a := 5
// 				calc := |n: Int|: Int ~>
// 					return 1 if n < 3

// 					n * a
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 2,
// 					byte(bytecode.INT_5),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.CLOSED_CLOSURE),
// 					2, 1,
// 					0xff,
// 					byte(bytecode.DUP),
// 					byte(bytecode.SET_LOCAL_2),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(83, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 2),
// 					bytecode.NewLineInfo(2, 2),
// 					bytecode.NewLineInfo(3, 7),
// 					bytecode.NewLineInfo(7, 1),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionWithUpvalues(
// 						functionSymbol,
// 						[]byte{
// 							byte(bytecode.GET_LOCAL_1),
// 							byte(bytecode.INT_3),
// 							byte(bytecode.JUMP_UNLESS_ILT), 0, 2,
// 							byte(bytecode.INT_1),
// 							byte(bytecode.RETURN),
// 							byte(bytecode.GET_LOCAL_1),
// 							byte(bytecode.GET_UPVALUE_0),
// 							byte(bytecode.MULTIPLY_INT),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(24, 3, 13), P(82, 7, 7)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(4, 7),
// 							bytecode.NewLineInfo(6, 3),
// 							bytecode.NewLineInfo(7, 1),
// 						},
// 						1,
// 						0,
// 						nil,
// 						1,
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

func TestGoStringLiteral(t *testing.T) {
	tests := goTestTable{
		"static string": {
			input: `a := "foo bar"`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.String // var a: Std::String
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.String("foo bar")
}
`,
		},
		"interpolated string with builtin types": {
			input: `
				bar := 15.2
				baz := "bazzy"
				foo := 1
				a := "foo: ${foo + 2}, bar: $bar, baz: $baz"
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("to_string")
var Std_ns_Int_im_to_string vm.NativeFunction // Std::Int.:to_string
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float // var bar: Std::Float
	_ = l0
	var l1 value.String // var baz: Std::String
	_ = l1
	var l2 value.Value // var foo: Std::Int
	_ = l2
	var l3 value.String // var a: Std::String
	_ = l3
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var t2 value.Value
	_ = t2
	var t3 []value.Value
	_ = t3
	var t4 value.String
	_ = t4
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	Std_ns_Int_im_to_string = vm.MethodToFunc((value.IntClass).LookupMethod(sym0))

	l0 = value.Float(15.2)
	l1 = value.String("bazzy")
	l2 = (value.SmallInt(1)).ToValue()
	t1, err = value.AddVal(l2, (value.SmallInt(2)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t3 = value.ResizeNativeArgs(t3, 2)
	t3[0] = t1
	thread.AddNativeCallFrame(sym0, sym1, 5)
	t2, err = Std_ns_Int_im_to_string(thread, t3) // receiver: Std::Int, name: to_string
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t4 = (t2).AsString()
	l3 = value.String("foo: ") + t4 + value.String(", bar: ") + (l0).ToString() + value.String(", baz: ") + l1
}
`,
		},
		"interpolated string with complex types": {
			input: `
				bar := 15.2
				baz := Time.now
				foo := 1
				a := "foo: ${foo + 2}, bar: $bar, baz: $baz"
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("now")
var Std_ns_Time_ns_now vm.NativeFunction // Std::Time::now
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("to_string")
var Std_ns_Int_im_to_string vm.NativeFunction  // Std::Int.:to_string
var Std_ns_Time_im_to_string vm.NativeFunction // Std::Time.:to_string

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float // var bar: Std::Float
	_ = l0
	var l1 value.Time // var baz: Std::Time
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Time
	_ = t3
	var l2 value.Value // var foo: Std::Int
	_ = l2
	var l3 value.String // var a: Std::String
	_ = l3
	var t4 value.Value
	_ = t4
	var t5 []value.Value
	_ = t5
	var t6 value.String
	_ = t6
	var t7 []value.Value
	_ = t7
	var t8 value.String
	_ = t8
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	Std_ns_Time_ns_now = vm.MethodToFunc(((value.TimeClass).SingletonClass()).LookupMethod(sym0))
	Std_ns_Int_im_to_string = vm.MethodToFunc((value.IntClass).LookupMethod(sym2))
	Std_ns_Time_im_to_string = vm.MethodToFunc((value.TimeClass).LookupMethod(sym2))

	l0 = value.Float(15.2)
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (value.TimeClass).ToValue()
	thread.AddNativeCallFrame(sym0, sym1, 3)
	t1, err = Std_ns_Time_ns_now(thread, t2) // receiver: &Std::Time, name: now
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t3 = (t1).AsTime()
	l1 = t3
	l2 = (value.SmallInt(1)).ToValue()
	t1, err = value.AddVal(l2, (value.SmallInt(2)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t5 = value.ResizeNativeArgs(t5, 2)
	t5[0] = t1
	thread.AddNativeCallFrame(sym2, sym1, 5)
	t4, err = Std_ns_Int_im_to_string(thread, t5) // receiver: Std::Int, name: to_string
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t6 = (t4).AsString()
	t7 = value.ResizeNativeArgs(t7, 2)
	t7[0] = (l1).ToValue()
	thread.AddNativeCallFrame(sym2, sym1, 5)
	t1, err = Std_ns_Time_im_to_string(thread, t7) // receiver: Std::Time, name: to_string
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t8 = (t1).AsString()
	l3 = value.String("foo: ") + t6 + value.String(", bar: ") + (l0).ToString() + value.String(", baz: ") + t8
}
`,
		},
		"inspect interpolated string": {
			input: `
				bar := 15.2
				foo := 1
				baz := "bazzy"
				a := "foo: #{foo + 2}, bar: #bar, baz: #baz"
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("inspect")
var Std_ns_Int_im_inspect vm.NativeFunction // Std::Int.:inspect
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float // var bar: Std::Float
	_ = l0
	var l1 value.Value // var foo: Std::Int
	_ = l1
	var l2 value.String // var baz: Std::String
	_ = l2
	var l3 value.String // var a: Std::String
	_ = l3
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var t2 value.Value
	_ = t2
	var t3 []value.Value
	_ = t3
	var t4 value.String
	_ = t4
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	Std_ns_Int_im_inspect = vm.MethodToFunc((value.IntClass).LookupMethod(sym0))

	l0 = value.Float(15.2)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.String("bazzy")
	t1, err = value.AddVal(l1, (value.SmallInt(2)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t3 = value.ResizeNativeArgs(t3, 2)
	t3[0] = t1
	thread.AddNativeCallFrame(sym0, sym1, 5)
	t2, err = Std_ns_Int_im_inspect(thread, t3) // receiver: Std::Int, name: inspect
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t4 = (t2).AsString()
	l3 = value.String("foo: ") + t4 + value.String(", bar: ") + value.String((l0).Inspect()) + value.String(", baz: ") + value.String((l2).Inspect())
}
`,
		},
	}

	for name, tc := range tests {
		noop(name, tc)
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoSymbolLiteral(t *testing.T) {
	tests := goTestTable{
		"static symbol": {
			input: `a := :foo`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Symbol // var a: Std::Symbol
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = sym0
}
`,
		},
		"interpolated symbol with builtin types": {
			input: `
				bar := 15.2
				baz := "bazzy"
				foo := 1
				a := :"foo: ${foo + 2}, bar: $bar, baz: $baz"
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("to_string")
var Std_ns_Int_im_to_string vm.NativeFunction // Std::Int.:to_string
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float // var bar: Std::Float
	_ = l0
	var l1 value.String // var baz: Std::String
	_ = l1
	var l2 value.Value // var foo: Std::Int
	_ = l2
	var l3 value.Symbol // var a: Std::Symbol
	_ = l3
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var t2 value.Value
	_ = t2
	var t3 []value.Value
	_ = t3
	var t4 value.String
	_ = t4
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	Std_ns_Int_im_to_string = vm.MethodToFunc((value.IntClass).LookupMethod(sym0))

	l0 = value.Float(15.2)
	l1 = value.String("bazzy")
	l2 = (value.SmallInt(1)).ToValue()
	t1, err = value.AddVal(l2, (value.SmallInt(2)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t3 = value.ResizeNativeArgs(t3, 2)
	t3[0] = t1
	thread.AddNativeCallFrame(sym0, sym1, 5)
	t2, err = Std_ns_Int_im_to_string(thread, t3) // receiver: Std::Int, name: to_string
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t4 = (t2).AsString()
	l3 = (value.String("foo: ") + t4 + value.String(", bar: ") + (l0).ToString() + value.String(", baz: ") + l1).ToSymbol()
}
`,
		},
		"interpolated symbol with complex types": {
			input: `
				bar := 15.2
				baz := Time.now
				foo := 1
				a := :"foo: ${foo + 2}, bar: $bar, baz: $baz"
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("now")
var Std_ns_Time_ns_now vm.NativeFunction // Std::Time::now
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("to_string")
var Std_ns_Int_im_to_string vm.NativeFunction  // Std::Int.:to_string
var Std_ns_Time_im_to_string vm.NativeFunction // Std::Time.:to_string

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float // var bar: Std::Float
	_ = l0
	var l1 value.Time // var baz: Std::Time
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Time
	_ = t3
	var l2 value.Value // var foo: Std::Int
	_ = l2
	var l3 value.Symbol // var a: Std::Symbol
	_ = l3
	var t4 value.Value
	_ = t4
	var t5 []value.Value
	_ = t5
	var t6 value.String
	_ = t6
	var t7 []value.Value
	_ = t7
	var t8 value.String
	_ = t8
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	Std_ns_Time_ns_now = vm.MethodToFunc(((value.TimeClass).SingletonClass()).LookupMethod(sym0))
	Std_ns_Int_im_to_string = vm.MethodToFunc((value.IntClass).LookupMethod(sym2))
	Std_ns_Time_im_to_string = vm.MethodToFunc((value.TimeClass).LookupMethod(sym2))

	l0 = value.Float(15.2)
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (value.TimeClass).ToValue()
	thread.AddNativeCallFrame(sym0, sym1, 3)
	t1, err = Std_ns_Time_ns_now(thread, t2) // receiver: &Std::Time, name: now
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t3 = (t1).AsTime()
	l1 = t3
	l2 = (value.SmallInt(1)).ToValue()
	t1, err = value.AddVal(l2, (value.SmallInt(2)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t5 = value.ResizeNativeArgs(t5, 2)
	t5[0] = t1
	thread.AddNativeCallFrame(sym2, sym1, 5)
	t4, err = Std_ns_Int_im_to_string(thread, t5) // receiver: Std::Int, name: to_string
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t6 = (t4).AsString()
	t7 = value.ResizeNativeArgs(t7, 2)
	t7[0] = (l1).ToValue()
	thread.AddNativeCallFrame(sym2, sym1, 5)
	t1, err = Std_ns_Time_im_to_string(thread, t7) // receiver: Std::Time, name: to_string
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t8 = (t1).AsString()
	l3 = (value.String("foo: ") + t6 + value.String(", bar: ") + (l0).ToString() + value.String(", baz: ") + t8).ToSymbol()
}
`,
		},
		"inspect interpolated symbol": {
			input: `
				bar := 15.2
				foo := 1
				baz := "bazzy"
				a := :"foo: #{foo + 2}, bar: #bar, baz: #baz"
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("inspect")
var Std_ns_Int_im_inspect vm.NativeFunction // Std::Int.:inspect
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float // var bar: Std::Float
	_ = l0
	var l1 value.Value // var foo: Std::Int
	_ = l1
	var l2 value.String // var baz: Std::String
	_ = l2
	var l3 value.Symbol // var a: Std::Symbol
	_ = l3
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var t2 value.Value
	_ = t2
	var t3 []value.Value
	_ = t3
	var t4 value.String
	_ = t4
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	Std_ns_Int_im_inspect = vm.MethodToFunc((value.IntClass).LookupMethod(sym0))

	l0 = value.Float(15.2)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.String("bazzy")
	t1, err = value.AddVal(l1, (value.SmallInt(2)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t3 = value.ResizeNativeArgs(t3, 2)
	t3[0] = t1
	thread.AddNativeCallFrame(sym0, sym1, 5)
	t2, err = Std_ns_Int_im_inspect(thread, t3) // receiver: Std::Int, name: inspect
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t4 = (t2).AsString()
	l3 = (value.String("foo: ") + t4 + value.String(", bar: ") + value.String((l0).Inspect()) + value.String(", baz: ") + value.String((l2).Inspect())).ToSymbol()
}
`,
		},
	}

	for name, tc := range tests {
		noop(name, tc)
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoRangeLiteral(t *testing.T) {
	tests := goTestTable{
		"static closed range": {
			input: `a := 2...5`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var range0 = value.NewClosedRange((value.SmallInt(2)).ToValue(), (value.SmallInt(5)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.ClosedRange // var a: Std::ClosedRange[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = range0
}
`,
		},
		"static open range": {
			input: `a := 2<.<5`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var range0 = value.NewOpenRange((value.SmallInt(2)).ToValue(), (value.SmallInt(5)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.OpenRange // var a: Std::OpenRange[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = range0
}
`,
		},
		"static left open range": {
			input: `a := 2<..5`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var range0 = value.NewLeftOpenRange((value.SmallInt(2)).ToValue(), (value.SmallInt(5)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.LeftOpenRange // var a: Std::LeftOpenRange[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = range0
}
`,
		},
		"static right open range": {
			input: `a := 2..<5`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var range0 = value.NewRightOpenRange((value.SmallInt(2)).ToValue(), (value.SmallInt(5)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.RightOpenRange // var a: Std::RightOpenRange[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = range0
}
`,
		},
		"static beginless closed range": {
			input: `a := ...5`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var range0 = value.NewBeginlessClosedRange((value.SmallInt(5)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.BeginlessClosedRange // var a: Std::BeginlessClosedRange[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = range0
}
`,
		},
		"static beginless open range": {
			input: `a := ..<5`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var range0 = value.NewBeginlessOpenRange((value.SmallInt(5)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.BeginlessOpenRange // var a: Std::BeginlessOpenRange[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = range0
}
`,
		},
		"static endless closed range": {
			input: `a := 2...`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var range0 = value.NewEndlessClosedRange((value.SmallInt(2)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.EndlessClosedRange // var a: Std::EndlessClosedRange[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = range0
}
`,
		},
		"static endless open range": {
			input: `a := 2<..`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var range0 = value.NewEndlessOpenRange((value.SmallInt(2)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.EndlessOpenRange // var a: Std::EndlessOpenRange[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = range0
}
`,
		},
		"closed range": {
			input: `
			  a := 2
				b := a...5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 *value.ClosedRange // var b: Std::ClosedRange[Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	l1 = value.NewClosedRange(l0, (value.SmallInt(5)).ToValue())
}
`,
		},
		"open range": {
			input: `
			  a := 2
				b := a<.<5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 *value.OpenRange // var b: Std::OpenRange[Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	l1 = value.NewOpenRange(l0, (value.SmallInt(5)).ToValue())
}
`,
		},
		"left open range": {
			input: `
			  a := 2
				b := a<..5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 *value.LeftOpenRange // var b: Std::LeftOpenRange[Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	l1 = value.NewLeftOpenRange(l0, (value.SmallInt(5)).ToValue())
}
`,
		},
		"right open range": {
			input: `
			  a := 2
				b := a..<5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 *value.RightOpenRange // var b: Std::RightOpenRange[Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	l1 = value.NewRightOpenRange(l0, (value.SmallInt(5)).ToValue())
}
`,
		},
		"beginless closed range": {
			input: `
			  a := 2
				b := ...a
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 *value.BeginlessClosedRange // var b: Std::BeginlessClosedRange[Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	l1 = value.NewBeginlessClosedRange(l0)
}
`,
		},
		"beginless open range": {
			input: `
			  a := 2
				b := ..<a
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 *value.BeginlessOpenRange // var b: Std::BeginlessOpenRange[Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	l1 = value.NewBeginlessOpenRange(l0)
}
`,
		},
		"endless closed range": {
			input: `
			  a := 2
				b := a...
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 *value.EndlessClosedRange // var b: Std::EndlessClosedRange[Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	l1 = value.NewEndlessClosedRange(l0)
}
`,
		},
		"endless open range": {
			input: `
			  a := 2
				b := a<..
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 *value.EndlessOpenRange // var b: Std::EndlessOpenRange[Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	l1 = value.NewEndlessOpenRange(l0)
}
`,
		},
	}

	for name, tc := range tests {
		noop(name, tc)
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoSimpleLiterals(t *testing.T) {
	tests := goTestTable{
		"put UInt8": {
			input: "a := 1u8",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.UInt8(1)
}
`,
		},
		"put UInt16": {
			input: "a := 25u16",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.UInt16(25)
}
`,
		},
		"put UInt32": {
			input: "a := 450_200u32",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.UInt32(450200)
}
`,
		},
		"put UInt64": {
			input: "a := 450_200u64",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.UInt64(450200)
}
`,
		},
		"put UInt": {
			input: "a := 450_200u",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.UInt(450200)
}
`,
		},
		"put Int8": {
			input: "a := 1i8",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Int8(1)
}
`,
		},
		"put Int16": {
			input: "a := 25i16",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Int16(25)
}
`,
		},
		"put Int32": {
			input: "a := 450_200i32",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Int32(450200)
}
`,
		},
		"put Int64": {
			input: "a := 450_200i64",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Int64(450200)
}
`,
		},
		"put SmallInt": {
			input: "a := 450_200",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(450200)).ToValue()
}
`,
		},
		"put BigInt": {
			input: fmt.Sprintf("a := %s", (&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)).String()),
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var bi0 = value.ParseBigIntPanic("9223372036854775812", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (bi0).ToValue()
}
`,
		},
		"put Float64": {
			input: "a := 45.5f64",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float64 // var a: Std::Float64
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Float64(45.5)
}
`,
		},
		"put Float32": {
			input: "a := 45.5f32",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float32 // var a: Std::Float32
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Float32(45.5)
}
`,
		},
		"put Float": {
			input: "a := 45.5",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float // var a: Std::Float
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Float(45.5)
}
`,
		},
		"put precise Float": {
			input: "a := 0.5827489723984",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float // var a: Std::Float
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Float(0.5827489723984)
}
`,
		},
		"put Raw String": {
			input: `a := 'foo\n'`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.String // var a: Std::String
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.String("foo\\n")
}
`,
		},
		"put String": {
			input: `a := "foo\n"`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.String // var a: Std::String
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.String("foo\n")
}
`,
		},
		"put raw Char": {
			input: "a := `I`",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Char // var a: Std::Char
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Char('I')
}
`,
		},
		"put Char": {
			input: "a := `\\n`",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Char // var a: Std::Char
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Char('\n')
}
`,
		},
		"put nil": {
			input: `a :=nil`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: nil
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
}
`,
		},
		"put true": {
			input: `a := true`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Bool // var a: bool
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = true
}
`,
		},
		"put false": {
			input: `a := false`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Bool // var a: bool
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = false
}
`,
		},
		"put simple Symbol": {
			input: `a := :foo`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Symbol // var a: Std::Symbol
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = sym0
}
`,
		},
		"put self": {
			input: `a := self`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Object
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = self
}
`,
		},
	}

	for name, tc := range tests {
		noop(name, tc)
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoArrayTuples(t *testing.T) {
	tests := goTestTable{
		"empty arrayTuple": {
			input: "a := %[]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewArrayTupleOfValueWithElements(0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[any]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"empty arrayTuple with native type": {
			input: "var a: ArrayTuple[Float] = %[]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Float](0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"with static elements": {
			input: "a := %[1, 'foo', 5, 5.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.SmallInt(5)).ToValue(), (value.Float(5.6)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int | Std::String | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"with static elements in immutable local": {
			input: "val a = %[1, 'foo', 5, 5.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.SmallInt(5)).ToValue(), (value.Float(5.6)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.ArrayTupleOfValue // var a: Std::ArrayTuple[Std::Int | Std::String | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"with static native elements": {
			input: "a := %[`b`, `c`, `d`]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Char](0, value.Char('b'), value.Char('c'), value.Char('d'))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Char]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"with static native elements in immutable local": {
			input: "val a = %[`b`, `c`, `d`]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Char](0, value.Char('b'), value.Char('c'), value.Char('d'))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.NativeArrayTuple[value.Char] // var a: Std::ArrayTuple[Std::Char]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"with static keyed elements": {
			input: "a := %[1, 'foo', 5 => 5,  3 => 5.6, :lol]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("lol")
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), value.Nil, (value.Float(5.6)).ToValue(), value.Nil, (value.SmallInt(5)).ToValue(), (sym0).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int | Std::String | Std::Symbol | nil | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"with static native keyed elements": {
			input: "a := %[1.2, 2.4, 5 => 5.0,  3 => 5.6, 420.69]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.Float(1.2)).ToValue(), (value.Float(2.4)).ToValue(), value.Nil, (value.Float(5.6)).ToValue(), value.Nil, (value.Float(5)).ToValue(), (value.Float(420.69)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Float | nil]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"nested static arrayTuples": {
			input: "a := %[1, %['bar', %[7.2]]]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.Float(7.2)).ToValue())
var arrtuple1 = value.NewArrayTupleOfValueWithElements(0, (value.String("bar")).ToValue(), (arrtuple0).ToValue())
var arrtuple2 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (arrtuple1).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::String | Std::ArrayTuple[Std::Float]]]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple2
}
`,
		},
		"nested static native arrayTuples": {
			input: "a := %[%['foo', 'bar'], %['baz', 'buzz']]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewNativeArrayTupleWithElements[value.String](0, value.String("foo"), value.String("bar"))
var arrtuple1 = value.NewNativeArrayTupleWithElements[value.String](0, value.String("baz"), value.String("buzz"))
var arrtuple2 = value.NewNativeArrayTupleWithElements[*value.NativeArrayTuple[value.String]](0, arrtuple0, arrtuple1)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::ArrayTuple[Std::String]]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple2
}
`,
		},
		"nested static with mutable elements": {
			input: "a := %[1, %['bar', [7.2]]]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::String | Std::ArrayList[Std::Float]]]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.NewArrayTupleOfValueWithElements(0, (value.String("bar")).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.Float(7.2)).ToValue())).ToValue())).ToValue())
}
`,
		},
		"static keyed elements": {
			input: "a := %[1, 'foo', 5i64 => 5,  3 => 5.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), value.Nil, (value.Float(5.6)).ToValue(), value.Nil, (value.SmallInt(5)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int | Std::String | nil | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"with static keyed and dynamic elements": {
			input: `
k := 10
a := %[1, 'foo', 5 => 5,  String.name, 3 => 5.6, k => 12, 8.2]
`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("name")
var Std_ns_Class_im_name vm.NativeFunction // Std::Class.:name
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var k: Std::Int
	_ = l0
	var l1 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int | Std::String | Std::Float | nil]
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 value.String
	_ = t3
	var t4 *value.ArrayTupleOfValue
	_ = t4
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	Std_ns_Class_im_name = vm.MethodToFunc((value.ClassClass).LookupMethod(sym0))

	l0 = (value.SmallInt(10)).ToValue()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (value.StringClass).ToValue()
	thread.AddNativeCallFrame(sym0, sym1, 3)
	t1, err = Std_ns_Class_im_name(thread, t2) // receiver: &Std::String, name: name
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t3 = (t1).AsString()
	t4 = value.NewArrayTupleOfValueWithElementsAndTotalCapacity(7, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), value.Nil, (value.Float(5.6)).ToValue(), value.Nil, (value.SmallInt(5)).ToValue(), (t3).ToValue())
	err = t4.AppendAt(l0, (value.SmallInt(12)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t4.Append((value.Float(8.2)).ToValue())
	l1 = t4
}
`,
		},
		"with static elements and if modifiers": {
			input: `
				var a: String? = "bar"
				b := %[1, 5 if a, %[:foo]]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Symbol](0, sym0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::String?
	_ = l0
	var l1 value.ArrayTuple // var b: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::Symbol]]
	_ = l1
	var t1 *value.ArrayTupleOfValue
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.String("bar")).ToValue()
	t1 = value.NewArrayTupleOfValueWithElementsAndTotalCapacity(3, (value.SmallInt(1)).ToValue())
	if value.Truthy(l0) {
		t1.Append((value.SmallInt(5)).ToValue())
	}
	t1.Append((arrtuple0).ToValue())
	l1 = t1
}
`,
		},
		"with static elements and if else modifiers": {
			input: `
				var a: String? = "bar"
				b := %[1, 5 if a else 2, %[:foo]]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Symbol](0, sym0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::String?
	_ = l0
	var l1 value.ArrayTuple // var b: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::Symbol]]
	_ = l1
	var t1 *value.ArrayTupleOfValue
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.String("bar")).ToValue()
	t1 = value.NewArrayTupleOfValueWithElementsAndTotalCapacity(3, (value.SmallInt(1)).ToValue())
	if value.Truthy(l0) {
		t1.Append((value.SmallInt(5)).ToValue())
	} else {
		t1.Append((value.SmallInt(2)).ToValue())
	}
	t1.Append((arrtuple0).ToValue())
	l1 = t1
}
`,
		},
		"with static elements and unless modifiers": {
			input: `
				var a: String? = nil
				b := %[1, 5 unless a, %[:foo]]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Symbol](0, sym0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::String?
	_ = l0
	var l1 value.ArrayTuple // var b: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::Symbol]]
	_ = l1
	var t1 *value.ArrayTupleOfValue
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1 = value.NewArrayTupleOfValueWithElementsAndTotalCapacity(3, (value.SmallInt(1)).ToValue())
	if value.Falsy(l0) {
		t1.Append((value.SmallInt(5)).ToValue())
	}
	t1.Append((arrtuple0).ToValue())
	l1 = t1
}
`,
		},
		// TODO: for in loops
		// 		"with static elements and for in loops": {
		// 			input: `
		// 				%[1.8, i * 2.0 for i in [1.0, 2.0, 3.0]]
		// 			`,
		// 			want: `package main

		// import "github.com/elk-language/elk/value"
		// import "github.com/elk-language/elk/vm"

		// import "github.com/elk-language/elk/value/symbol"

		// var _ = symbol.Value
		// var _ = vm.New
		// var _ = value.Truthy

		// func main() { // loc: <main>
		// 	thread := vm.New()
		// 	_ = thread
		// 	var t1 *value.NativeArrayTuple[value.Float]
		// 	_ = t1
		// 	var self value.Value
		// 	_ = self

		// 	self = value.Ref(value.GlobalObject)
		// 	t1 = value.NewNativeArrayTupleWithElements[value.Float](0, value.Float(1.800000))
		// }
		// `,
		// 		},
		"with dynamic elements and if modifiers": {
			input: `
				var a: Object? = nil
				b := %[String.name, 5 if a, %[:foo]]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("name")
var Std_ns_Class_im_name vm.NativeFunction // Std::Class.:name
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("foo")
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Symbol](0, sym2)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Object?
	_ = l0
	var l1 value.ArrayTuple // var b: Std::ArrayTuple[Std::String | Std::Int | Std::ArrayTuple[Std::Symbol]]
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 value.String
	_ = t3
	var t4 *value.ArrayTupleOfValue
	_ = t4
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	Std_ns_Class_im_name = vm.MethodToFunc((value.ClassClass).LookupMethod(sym0))

	l0 = value.Nil
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (value.StringClass).ToValue()
	thread.AddNativeCallFrame(sym0, sym1, 3)
	t1, err = Std_ns_Class_im_name(thread, t2) // receiver: &Std::String, name: name
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t3 = (t1).AsString()
	t4 = value.NewArrayTupleOfValueWithElementsAndTotalCapacity(3, (t3).ToValue())
	if value.Truthy(l0) {
		t4.Append((value.SmallInt(5)).ToValue())
	}
	t4.Append((arrtuple0).ToValue())
	l1 = t4
}
`,
		},
		// TODO: constructor
		// "with dynamic and keyed elements": {
		// 	input: "%[Object(), 1, 'foo', 5 => 5,  3 => 5.6]",
		// 	want: `
		// 	`,
		// },
		"with keyed and if elements": {
			input: `
				var a: String? = nil
				b := %[3 => 5 if a]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::String?
	_ = l0
	var l1 value.ArrayTuple // var b: Std::ArrayTuple[Std::Int | nil]
	_ = l1
	var t1 *value.ArrayTupleOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1 = value.NewArrayTupleOfValueWithElementsAndTotalCapacity(1)
	if value.Truthy(l0) {
		err = t1.AppendAt((value.SmallInt(3)).ToValue(), (value.SmallInt(5)).ToValue())
		if err.IsNotUndefined() {
			thread.Panic(err)
		}
	}
	l1 = t1
}
`,
		},
		"with static concat": {
			input: "a := %[1, 2, 3] + %[4, 5, 6] + %[10]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue(), (value.SmallInt(6)).ToValue(), (value.SmallInt(10)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"with static concat with list": {
			input: "a := %[1, 2, 3] + [4, 5, 6] + %[10]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue(), (value.SmallInt(6)).ToValue(), (value.SmallInt(10)).ToValue())
}
`,
		},
		"with static repeat": {
			input: "a := %[1, 2, 3] * 3",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"with static concat and nested tuples": {
			input: "a := %[1, 2, 3] + %[4, 5, 6, %[7, 8]] + %[10]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(7)).ToValue(), (value.SmallInt(8)).ToValue())
var arrtuple1 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue(), (value.SmallInt(6)).ToValue(), (arrtuple0).ToValue(), (value.SmallInt(10)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::Int]]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple1
}
`,
		},
		"word arrayTuple": {
			input: `a := %w[foo bar baz]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewNativeArrayTupleWithElements[value.String](0, value.String("foo"), value.String("bar"), value.String("baz"))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::String]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"symbol arrayTuple": {
			input: `a := %s[foo bar baz]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")
var sym1 = value.ToSymbol("bar")
var sym2 = value.ToSymbol("baz")
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Symbol](0, sym0, sym1, sym2)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Symbol]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"hex arrayTuple uint8": {
			input: `a := %x[ab cd 5f]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewNativeArrayTupleWithElements[value.UInt8](0, value.UInt8(171), value.UInt8(205), value.UInt8(95))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::UInt8]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"hex arrayTuple uint16": {
			input: `a := %x[ab_cd cd 5f]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewNativeArrayTupleWithElements[value.UInt16](0, value.UInt16(43981), value.UInt16(205), value.UInt16(95))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::UInt16]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"hex arrayTuple uint32": {
			input: `a := %x[ab_cd_ab_cd cd 5f]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewNativeArrayTupleWithElements[value.UInt32](0, value.UInt32(2882382797), value.UInt32(205), value.UInt32(95))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::UInt32]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"hex arrayTuple uint64": {
			input: `a := %x[ab_cd_ab_cd_ab_cd_ab_cd cd 5f]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewNativeArrayTupleWithElements[value.UInt64](0, value.UInt64(12379739850550389709), value.UInt64(205), value.UInt64(95))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::UInt64]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"hex arrayTuple int": {
			input: `a := %x[ab_cd_ab_cd_ab_cd_ab_cd_ab_cd_ab_cd_ab_cd_ab_cd cd 5f]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var bi0 = value.ParseBigIntPanic("228365892722206371581333312115001109453", 0)
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (bi0).ToValue(), (value.SmallInt(205)).ToValue(), (value.SmallInt(95)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
		"bin arrayTuple": {
			input: `a := %b[101 11 10]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var arrtuple0 = value.NewNativeArrayTupleWithElements[value.UInt8](0, value.UInt8(5), value.UInt8(3), value.UInt8(2))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::UInt8]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = arrtuple0
}
`,
		},
	}

	for name, tc := range tests {
		noop(name, tc)
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func noop(args ...any) {}

func TestGoArrayLists(t *testing.T) {
	tests := goTestTable{
		"empty list": {
			input: "a := []",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[any]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElements(0)
}
`,
		},
		"with static elements": {
			input: "a := [1, 'foo', 5, 5.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int | Std::String | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.SmallInt(5)).ToValue(), (value.Float(5.6)).ToValue())
}
`,
		},
		"with static elements in immutable local": {
			input: "val a = [1, 'foo', 5, 5.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.ArrayListOfValue // var a: Std::ArrayList[Std::Int | Std::String | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.SmallInt(5)).ToValue(), (value.Float(5.6)).ToValue())
}
`,
		},
		"with native static elements": {
			input: "a := [1.2, 5.0, 10.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewNativeArrayListWithElements[value.Float](0, value.Float(1.2), value.Float(5), value.Float(10.6))
}
`,
		},
		"with native static elements in immutable local": {
			input: "val a = [1.2, 5.0, 10.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.NativeArrayList[value.Float] // var a: Std::ArrayList[Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewNativeArrayListWithElements[value.Float](0, value.Float(1.2), value.Float(5), value.Float(10.6))
}
`,
		},
		"with static elements and static capacity": {
			input: "a := [1, 'foo', 5, 5.6]:10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int | Std::String | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElementsAndTotalCapacity(4+int(value.SmallInt(10)), (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.SmallInt(5)).ToValue(), (value.Float(5.6)).ToValue())
}
`,
		},

		"with static elements and dynamic capacity": {
			input: `
				cap := 2
				a := [1, 'foo', 5, 5.6]:cap
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var cap: Std::Int
	_ = l0
	var l1 value.ArrayList // var a: Std::ArrayList[Std::Int | Std::String | Std::Float]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	l1 = value.NewArrayListOfValueWithElementsAndTotalCapacity(4+(l0).AsAnyInt(), (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.SmallInt(5)).ToValue(), (value.Float(5.6)).ToValue())
}
`,
		},
		"word list": {
			input: `a := \w[foo bar baz]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::String]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewNativeArrayListWithElements[value.String](0, value.String("foo"), value.String("bar"), value.String("baz"))
}
`,
		},
		"word list with capacity": {
			input: `a := \w[foo bar baz]:15`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::String]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewNativeArrayListWithElements[value.String](15, value.String("foo"), value.String("bar"), value.String("baz"))
}
`,
		},
		"symbol list": {
			input: `a := \s[foo bar baz]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")
var sym1 = value.ToSymbol("bar")
var sym2 = value.ToSymbol("baz")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Symbol]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewNativeArrayListWithElements[value.Symbol](0, sym0, sym1, sym2)
}
`,
		},

		"symbol list with capacity": {
			input: `a := \s[foo bar baz]:15`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")
var sym1 = value.ToSymbol("bar")
var sym2 = value.ToSymbol("baz")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Symbol]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewNativeArrayListWithElements[value.Symbol](15, sym0, sym1, sym2)
}
`,
		},
		"hex list": {
			input: `a := \x[ab cd 5f]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::UInt8]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewNativeArrayListWithElements[value.UInt8](0, value.UInt8(171), value.UInt8(205), value.UInt8(95))
}
`,
		},
		"hex list with capacity": {
			input: `a := \x[ab cd 5f]:2`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::UInt8]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewNativeArrayListWithElements[value.UInt8](2, value.UInt8(171), value.UInt8(205), value.UInt8(95))
}
`,
		},
		"bin list": {
			input: `a := \b[101 11 10]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::UInt8]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewNativeArrayListWithElements[value.UInt8](0, value.UInt8(5), value.UInt8(3), value.UInt8(2))
}
`,
		},
		"bin list with capacity": {
			input: `a := \b[101 11 10]:3`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::UInt8]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewNativeArrayListWithElements[value.UInt8](3, value.UInt8(5), value.UInt8(3), value.UInt8(2))
}
`,
		},

		"with static keyed elements": {
			input: "a := [1, 'foo', 5 => 5,  3 => 5.6, :lol]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("lol")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int | Std::String | Std::Symbol | nil | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElements(2, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), value.Nil, (value.Float(5.6)).ToValue(), value.Nil, (value.SmallInt(5)).ToValue(), (sym0).ToValue())
}
`,
		},
		"with static keyed elements and static capacity": {
			input: "a := [1, 'foo', 5 => 5,  3 => 5.6, :lol]:6",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("lol")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int | Std::String | Std::Symbol | nil | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElementsAndTotalCapacity(5+int(value.SmallInt(6)), (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), value.Nil, (value.Float(5.6)).ToValue(), value.Nil, (value.SmallInt(5)).ToValue(), (sym0).ToValue())
}
`,
		},
		"with static concat": {
			input: "a := [1, 2, 3] + [4, 5, 6] + [10]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue(), (value.SmallInt(6)).ToValue(), (value.SmallInt(10)).ToValue())
}
`,
		},
		"with static repeat": {
			input: "a := [1, 2, 3] * 3",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue())
}
`,
		},
		"with static concat and nested lists": {
			input: "a := [1, 2, 3] + [4, 5, 6, [7, 8]] + [10]",

			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int | Std::ArrayList[Std::Int]]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue(), (value.SmallInt(6)).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.SmallInt(7)).ToValue(), (value.SmallInt(8)).ToValue())).ToValue(), (value.SmallInt(10)).ToValue())
}
`,
		},

		"nested static lists": {
			input: "a := [1, ['bar', [7.2]]]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int | Std::ArrayList[Std::String | Std::ArrayList[Std::Float]]]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.String("bar")).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.Float(7.2)).ToValue())).ToValue())).ToValue())
}
`,
		},
		"with static keyed and dynamic elements": {
			input: `
				a := 5
				b := [1, 'foo', 5 => 5,  3 => 5.6, a]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.ArrayList // var b: Std::ArrayList[Std::Int | Std::String | nil | Std::Float]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(5)).ToValue()
	l1 = value.NewArrayListOfValueWithElementsAndTotalCapacity(5+0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), value.Nil, (value.Float(5.6)).ToValue(), value.Nil, (value.SmallInt(5)).ToValue(), l0)
}
`,
		},
		"with static keyed, dynamic elements and capacity": {
			input: `
				a := 5
				b := [1, 'foo', 5 => 5,  3 => 5.6, a]:15
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.ArrayList // var b: Std::ArrayList[Std::Int | Std::String | nil | Std::Float]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(5)).ToValue()
	l1 = value.NewArrayListOfValueWithElementsAndTotalCapacity(5+int(value.SmallInt(15)), (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), value.Nil, (value.Float(5.6)).ToValue(), value.Nil, (value.SmallInt(5)).ToValue(), l0)
}
`,
		},
		"with static and dynamic elements": {
			input: `
				var a: Int? = 3
				b := [1, 'foo', 5, a, 5, %[:foo]]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Symbol](0, sym0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.ArrayList // var b: Std::ArrayList[Std::Int | Std::String | Std::ArrayTuple[Std::Symbol] | nil]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(3)).ToValue()
	l1 = value.NewArrayListOfValueWithElementsAndTotalCapacity(6+0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.SmallInt(5)).ToValue(), l0, (value.SmallInt(5)).ToValue(), (arrtuple0).ToValue())
}
`,
		},
		"with dynamic elements": {
			input: `
				a := 3
				b := [a, 5, [:foo]]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.ArrayList // var b: Std::ArrayList[Std::Int | Std::ArrayList[Std::Symbol]]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(3)).ToValue()
	l1 = value.NewArrayListOfValueWithElementsAndTotalCapacity(3+0, l0, (value.SmallInt(5)).ToValue(), (value.NewNativeArrayListWithElements[value.Symbol](0, sym0)).ToValue())
}
`,
		},

		"with static elements and if modifiers": {
			input: `
				var a: Int? = nil
				b := [1, 5 if a, [:foo]]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.ArrayList // var b: Std::ArrayList[Std::Int | Std::ArrayList[Std::Symbol]]
	_ = l1
	var t1 *value.ArrayListOfValue
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1 = value.NewArrayListOfValueWithElementsAndTotalCapacity(3+0, (value.SmallInt(1)).ToValue())
	if value.Truthy(l0) {
		t1.Append((value.SmallInt(5)).ToValue())
	}
	t1.Append((value.NewNativeArrayListWithElements[value.Symbol](0, sym0)).ToValue())
	l1 = t1
}
`,
		},
		"with static elements, if modifiers and capacity": {
			input: `
				var a: Int? = nil
				[1, 5 if a, [:foo]]:45
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(47, 3, 25), P(48, 3, 26)),
					"capacity cannot be specified in collection literals with conditional elements or loops",
				),
			},
		},
		"with static elements and unless modifiers": {
			input: `
				var a: Int? = nil
				b := [1, 5 unless a, [:foo]]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.ArrayList // var b: Std::ArrayList[Std::Int | Std::ArrayList[Std::Symbol]]
	_ = l1
	var t1 *value.ArrayListOfValue
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1 = value.NewArrayListOfValueWithElementsAndTotalCapacity(3+0, (value.SmallInt(1)).ToValue())
	if value.Falsy(l0) {
		t1.Append((value.SmallInt(5)).ToValue())
	}
	t1.Append((value.NewNativeArrayListWithElements[value.Symbol](0, sym0)).ToValue())
	l1 = t1
}
`,
		},
		// TODO: for in
		// 		"with static elements and for in loops": {
		// 			input: `
		// 				b := [1, i * 2 for i in [1, 2, 3], %[:foo]]
		// 			`,
		// 			want: `
		// `,
		// 		},

		// TODO: Constructor
		// "with dynamic elements and if modifiers": {
		// 	input: `
		// 		var a: Int? = nil
		// 		[Object(), 5 if a, [:foo]]
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.NIL),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.INSTANTIATE8), 0,
		// 			byte(bytecode.NEW_ARRAY_LIST8), 1,
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.JUMP_UNLESS), 0, 5,
		// 			byte(bytecode.INT_5),
		// 			byte(bytecode.APPEND),
		// 			byte(bytecode.JUMP), 0, 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.COPY),
		// 			byte(bytecode.APPEND),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(53, 3, 31)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(3, 21),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Object").ToValue(),
		// 			value.Ref(&value.ArrayList{
		// 				value.ToSymbol("foo").ToValue(),
		// 			}),
		// 		},
		// 	),
		// },

		"with dynamic and keyed elements": {
			input: `
				var a: Int? = nil
				b := [a, 1, 'foo', 5 => 5,  3 => 5.6]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.ArrayList // var b: Std::ArrayList[Std::Int | Std::String | nil | Std::Float]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	l1 = value.NewArrayListOfValueWithElementsAndTotalCapacity(5+0, l0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.Float(5.6)).ToValue(), value.Nil, (value.SmallInt(5)).ToValue())
}
`,
		},
		"with dynamic, keyed elements and capacity": {
			input: `
				a := 3
				b := [a, 1, 'foo', 5 => 5,  3 => 5.6]:7
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.ArrayList // var b: Std::ArrayList[Std::Int | Std::String | nil | Std::Float]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(3)).ToValue()
	l1 = value.NewArrayListOfValueWithElementsAndTotalCapacity(5+int(value.SmallInt(7)), l0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.Float(5.6)).ToValue(), value.Nil, (value.SmallInt(5)).ToValue())
}
`,
		},
		"with keyed and if elements": {
			input: `
				var a: Int? = nil
				b := [3 => 5 if a]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.ArrayList // var b: Std::ArrayList[Std::Int | nil]
	_ = l1
	var t1 *value.ArrayListOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1 = value.NewArrayListOfValueWithElementsAndTotalCapacity(1 + 0)
	if value.Truthy(l0) {
		err = t1.AppendAt((value.SmallInt(3)).ToValue(), (value.SmallInt(5)).ToValue())
		if err.IsNotUndefined() {
			thread.Panic(err)
		}
	}
	l1 = t1
}
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoHashSet(t *testing.T) {
	tests := goTestTable{
		"empty list": {
			input: "a := ^[]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[any]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0)
}
`,
		},
		"with static elements": {
			input: "a := ^[1, 'foo', 5, 5.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::Int | Std::String | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, (value.String("foo")).ToValue(), (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue(), (value.Float(5.6)).ToValue())
}
`,
		},
		"with static elements in immutable local": {
			input: "val a = ^[1, 'foo', 5, 5.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *vm.HashSetOfValue // var a: Std::HashSet[Std::Int | Std::String | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, (value.String("foo")).ToValue(), (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue(), (value.Float(5.6)).ToValue())
}
`,
		},
		"with native static elements": {
			input: "a := ^[1.2, 10.0, 5.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.NewNativeHashSetWithElements[value.Float](value.Float(1.2), value.Float(10), value.Float(5.6))
}
`,
		},
		"with native static elements in immutable local": {
			input: "val a = ^[1.2, 10.0, 5.6]",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *vm.NativeHashSet[value.Float] // var a: Std::HashSet[Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.NewNativeHashSetWithElements[value.Float](value.Float(1.2), value.Float(10), value.Float(5.6))
}
`,
		},
		"with static elements and static capacity": {
			input: "a := ^[1, 'foo', 5, 5.6]:10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::Int | Std::String | Std::Float]
	_ = l0
	var t1 *vm.HashSetOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	t1, err = vm.NewHashSetOfValueWithCapacityAndElements(thread, 4+int(value.SmallInt(10)), (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.SmallInt(5)).ToValue(), (value.Float(5.6)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l0 = t1
}
`,
		},
		"with static elements and dynamic capacity": {
			input: `
				cap := 2
				a := ^[1, 'foo', 5, 5.6]:cap
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var cap: Std::Int
	_ = l0
	var l1 vm.HashSet // var a: Std::HashSet[Std::Int | Std::String | Std::Float]
	_ = l1
	var t1 *vm.HashSetOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	t1, err = vm.NewHashSetOfValueWithCapacityAndElements(thread, 4+(l0).AsAnyInt(), (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.SmallInt(5)).ToValue(), (value.Float(5.6)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},

		"word set": {
			input: `a := ^w[foo bar baz]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::String]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.NewNativeHashSetWithElements[value.String](value.String("bar"), value.String("baz"), value.String("foo"))
}
`,
		},
		"word set with capacity": {
			input: `a := ^w[foo bar baz]:15`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::String]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.NewNativeHashSetWithElements[value.String](value.String("bar"), value.String("baz"), value.String("foo"))
}
`,
		},
		"symbol set": {
			input: `a := ^s[foo bar baz]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("bar")
var sym1 = value.ToSymbol("baz")
var sym2 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::Symbol]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.NewNativeHashSetWithElements[value.Symbol](sym0, sym1, sym2)
}
`,
		},
		"symbol set with capacity": {
			input: `a := ^s[foo bar baz]:15`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("bar")
var sym1 = value.ToSymbol("baz")
var sym2 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::Symbol]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.NewNativeHashSetWithElements[value.Symbol](sym0, sym1, sym2)
}
`,
		},
		"hex set": {
			input: `a := ^x[ab cd 5f]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, (value.SmallInt(171)).ToValue(), (value.SmallInt(205)).ToValue(), (value.SmallInt(95)).ToValue())
}
`,
		},
		"hex set with capacity": {
			input: `a := ^x[ab cd 5f]:2`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, (value.SmallInt(171)).ToValue(), (value.SmallInt(205)).ToValue(), (value.SmallInt(95)).ToValue())
}
`,
		},

		"bin set": {
			input: `a := ^b[101 11 10]`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(5)).ToValue())
}
`,
		},
		"bin set with capacity": {
			input: `a := ^b[101 11 10]:3`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashSet // var a: Std::HashSet[Std::Int]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(5)).ToValue())
}
`,
		},
		"with static and dynamic elements": {
			input: `
				var a: Int? = nil
				b := ^[1, 'foo', 5, a, 5]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashSet // var b: Std::HashSet[Std::Int | Std::String | nil]
	_ = l1
	var t1 *vm.HashSetOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1, err = vm.NewHashSetOfValueWithCapacityAndElements(thread, 5+0, (value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue(), (value.SmallInt(5)).ToValue(), l0, (value.SmallInt(5)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"with dynamic elements": {
			input: `
				var a: Int? = nil
				b := ^[a, 5]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashSet // var b: Std::HashSet[Std::Int | nil]
	_ = l1
	var t1 *vm.HashSetOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1, err = vm.NewHashSetOfValueWithCapacityAndElements(thread, 2+0, l0, (value.SmallInt(5)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"with static elements and if modifiers": {
			input: `
				var a: Int? = nil
				b := ^[1, 5 if a]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashSet // var b: Std::HashSet[Std::Int]
	_ = l1
	var t1 *vm.HashSetOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1, err = vm.NewHashSetOfValueWithCapacityAndElements(thread, 2+0, (value.SmallInt(1)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	if value.Truthy(l0) {
		_, err = t1.AppendVal(thread, (value.SmallInt(5)).ToValue())
		if err.IsNotUndefined() {
			thread.Panic(err)
		}
	}
	l1 = t1
}
`,
		},
		"with native static elements and if modifiers": {
			input: `
				var a: Int? = nil
				b := ^[1.2, 5.0 if a]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashSet // var b: Std::HashSet[Std::Float]
	_ = l1
	var t1 *vm.NativeHashSet[value.Float]
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1 = vm.NewNativeHashSetWithElementsAndTotalCapacity[value.Float](2+0, value.Float(1.2))
	if value.Truthy(l0) {
		t1.Append(value.Float(5))
	}
	l1 = t1
}
`,
		},
		"with static elements, if modifiers and capacity": {
			input: `
				var a: Int? = nil
				^[1, 5 if a]:45
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(40, 3, 18), P(41, 3, 19)),
					"capacity cannot be specified in collection literals with conditional elements or loops",
				),
			},
		},
		"with static elements and unless modifiers": {
			input: `
				var a: Int? = nil
				b := ^[1, 5 unless a]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashSet // var b: Std::HashSet[Std::Int]
	_ = l1
	var t1 *vm.HashSetOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1, err = vm.NewHashSetOfValueWithCapacityAndElements(thread, 2+0, (value.SmallInt(1)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	if value.Falsy(l0) {
		_, err = t1.AppendVal(thread, (value.SmallInt(5)).ToValue())
		if err.IsNotUndefined() {
			thread.Panic(err)
		}
	}
	l1 = t1
}
`,
		},
		// TODO: For in loops
		// 		"with static elements and for in loops": {
		// 			input: `
		// 				a := ^[1, i * 2 for i in [1, 2, 3], 2]
		// 			`,
		// 			want: `
		// `,
		// 		},

		// TODO: constructors
		// "with dynamic elements and if modifiers": {
		// 	input: `
		// 		var a: Int? = nil
		// 		^[Object(), 5 if a]
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.NIL),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.LOAD_VALUE_0),
		// 			byte(bytecode.GET_CONST8), 1,
		// 			byte(bytecode.INSTANTIATE8), 0,
		// 			byte(bytecode.NEW_HASH_SET8), 1,
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.JUMP_UNLESS), 0, 5,
		// 			byte(bytecode.INT_5),
		// 			byte(bytecode.APPEND),
		// 			byte(bytecode.JUMP), 0, 0,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(46, 3, 24)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(3, 18),
		// 		},
		// 		[]value.Value{
		// 			value.Ref(vm.MustNewHashSetWithCapacityAndElements(
		// 				nil,
		// 				2,
		// 			)),
		// 			value.ToSymbol("Std::Object").ToValue(),
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

func TestGoHashMap(t *testing.T) {
	tests := goTestTable{
		"empty": {
			input: "a := {}",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashMap // var a: Std::HashMap[any, any]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashMapOfValueWithCapacityAndElements(nil, 0)
}
`,
		},
		"shorthand local": {
			input: `
				foo := 3
				a := { foo }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var foo: Std::Int
	_ = l0
	var l1 vm.HashMap // var a: Std::HashMap[Std::Symbol, Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(3)).ToValue()
	l1 = vm.NewNativeKeyHashMapWithElementsAndTotalCapacity[value.Symbol](1+0, value.MakeNativePair(sym0, l0))
}
`,
		},
		"native map in immutable local": {
			input: `
				foo := 3.2
				val a = { foo }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float // var foo: Std::Float
	_ = l0
	var l1 *vm.NativeHashMap[value.Symbol, value.Float] // var a: Std::HashMap[Std::Symbol, Std::Float]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Float(3.2)
	l1 = vm.NewNativeHashMapWithElementsAndTotalCapacity[value.Symbol, value.Float](1+0, value.MakeNativePair(sym0, l0))
}
`,
		},
		"shorthand private local": {
			input: `
				_foo := 3
				a := { _foo }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("_foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var _foo: Std::Int
	_ = l0
	var l1 vm.HashMap // var a: Std::HashMap[Std::Symbol, Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(3)).ToValue()
	l1 = vm.NewNativeKeyHashMapWithElementsAndTotalCapacity[value.Symbol](1+0, value.MakeNativePair(sym0, l0))
}
`,
		},
		"with static elements": {
			input: `a := { 1 => 'foo', foo: 5, "bar" => 5.6 }`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashMap // var a: Std::HashMap[Std::Int | Std::Symbol | Std::String, Std::String | Std::Int | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashMapOfValueWithCapacityAndElements(nil, 0, value.MakePairOfValue((value.String("bar")).ToValue(), (value.Float(5.6)).ToValue()), value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue()), value.MakePairOfValue((sym0).ToValue(), (value.SmallInt(5)).ToValue()))
}
`,
		},
		// TODO: for in loops
		// 		"with static elements and for loops": {
		// 			input: `{ 1 => 'foo', i => i ** 2 for i in [1, 2, 3], 2 => 5.6 }`,
		// 			want: vm.NewBytecodeFunctionNoParams(
		// 				mainSymbol,
		// 				[]byte{
		// 					byte(bytecode.PREP_LOCALS8), 2,
		// 					byte(bytecode.UNDEFINED),
		// 					byte(bytecode.LOAD_VALUE_0),
		// 					byte(bytecode.NEW_HASH_MAP8), 0,
		// 					byte(bytecode.LOAD_VALUE_1),
		// 					byte(bytecode.COPY),
		// 					byte(bytecode.GET_ITERATOR),
		// 					byte(bytecode.SET_LOCAL_1),
		// 					byte(bytecode.GET_LOCAL_1),
		// 					byte(bytecode.FOR_IN_BUILTIN), 0, 9,
		// 					byte(bytecode.SET_LOCAL_2),
		// 					byte(bytecode.GET_LOCAL_2),
		// 					byte(bytecode.GET_LOCAL_2),
		// 					byte(bytecode.INT_2),
		// 					byte(bytecode.EXPONENTIATE_INT),
		// 					byte(bytecode.MAP_SET),
		// 					byte(bytecode.LOOP), 0, 13,
		// 					byte(bytecode.INT_2),
		// 					byte(bytecode.LOAD_VALUE_2),
		// 					byte(bytecode.MAP_SET),
		// 					byte(bytecode.RETURN),
		// 				},
		// 				L(P(0, 1, 1), P(55, 1, 56)),
		// 				bytecode.LineInfoList{
		// 					bytecode.NewLineInfo(1, 27),
		// 				},
		// 				[]value.Value{
		// 					value.Ref(vm.MustNewHashMapWithCapacityAndElements(
		// 						nil,
		// 						3,
		// 						value.Pair{
		// 							Key:   value.SmallInt(1).ToValue(),
		// 							Value: value.Ref(value.String("foo")),
		// 						},
		// 					)),
		// 					value.Ref(&value.ArrayList{
		// 						value.SmallInt(1).ToValue(),
		// 						value.SmallInt(2).ToValue(),
		// 						value.SmallInt(3).ToValue(),
		// 					}),
		// 					value.Float(5.6).ToValue(),
		// 				},
		// 			),
		// 		},
		"with static elements and static capacity": {
			input: `a := { 1 => 'foo', foo: 5, "bar" => 5.6 }:10`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashMap // var a: Std::HashMap[Std::Int | Std::Symbol | Std::String, Std::String | Std::Int | Std::Float]
	_ = l0
	var t1 *vm.HashMapOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	t1, err = vm.NewHashMapOfValueWithCapacityAndElements(thread, 3+int(value.SmallInt(10)), value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue()), value.MakePairOfValue((sym0).ToValue(), (value.SmallInt(5)).ToValue()), value.MakePairOfValue((value.String("bar")).ToValue(), (value.Float(5.6)).ToValue()))
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l0 = t1
}
`,
		},
		"with static elements and dynamic capacity": {
			input: `
				cap := 2
				a := { 1 => 'foo', foo: 5, "bar" => 5.6 }:cap
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var cap: Std::Int
	_ = l0
	var l1 vm.HashMap // var a: Std::HashMap[Std::Int | Std::Symbol | Std::String, Std::String | Std::Int | Std::Float]
	_ = l1
	var t1 *vm.HashMapOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(2)).ToValue()
	t1, err = vm.NewHashMapOfValueWithCapacityAndElements(thread, 3+(l0).AsAnyInt(), value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue()), value.MakePairOfValue((sym0).ToValue(), (value.SmallInt(5)).ToValue()), value.MakePairOfValue((value.String("bar")).ToValue(), (value.Float(5.6)).ToValue()))
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"nested static": {
			input: "a := { 1 => { 'bar' => [7.2] } }",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashMap // var a: Std::HashMap[Std::Int, Std::HashMap[Std::String, Std::ArrayList[Std::Float]]]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashMapOfValueWithCapacityAndElements(nil, 0, value.MakePairOfValue((value.SmallInt(1)).ToValue(), (vm.MustNewHashMapOfValueWithCapacityAndElements(nil, 0, value.MakePairOfValue((value.String("bar")).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.Float(7.2)).ToValue())).ToValue()))).ToValue()))
}
`,
		},
		"with static and dynamic elements": {
			input: `
				a := 5
				b := { 1 => 'foo', 5 => a, 5 => %[:foo] }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Symbol](0, sym0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 vm.HashMap // var b: Std::HashMap[Std::Int, Std::String | Std::Int | Std::ArrayTuple[Std::Symbol]]
	_ = l1
	var t1 *vm.HashMapOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(5)).ToValue()
	t1, err = vm.NewHashMapOfValueWithCapacityAndElements(thread, 3+0, value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue()), value.MakePairOfValue((value.SmallInt(5)).ToValue(), l0), value.MakePairOfValue((value.SmallInt(5)).ToValue(), (arrtuple0).ToValue()))
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"with static elements and if modifiers": {
			input: `
				var a: Int? = nil
				b := { 2 => 5, 1 => 5 if a, a: [:foo] }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("a")
var sym1 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashMap // var b: Std::HashMap[Std::Int | Std::Symbol, Std::Int | Std::ArrayList[Std::Symbol]]
	_ = l1
	var t1 *vm.HashMapOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1, err = vm.NewHashMapOfValueWithCapacityAndElements(thread, 3+0, value.MakePairOfValue((value.SmallInt(2)).ToValue(), (value.SmallInt(5)).ToValue()))
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	if value.Truthy(l0) {
		err = t1.SetVal(thread, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue())
		if err.IsNotUndefined() {
			thread.Panic(err)
		}
	}
	err = t1.SetVal(thread, (sym0).ToValue(), (value.NewNativeArrayListWithElements[value.Symbol](0, sym1)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"with static elements and if modifiers in immutable local": {
			input: `
				var a: Int? = nil
				val b = { 2 => 5, 1 => 5 if a, a: [:foo] }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("a")
var sym1 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var t1 *vm.HashMapOfValue
	_ = t1
	var err value.Value
	_ = err
	var l1 *vm.HashMapOfValue // var b: Std::HashMap[Std::Int | Std::Symbol, Std::Int | Std::ArrayList[Std::Symbol]]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1, err = vm.NewHashMapOfValueWithCapacityAndElements(thread, 3+0, value.MakePairOfValue((value.SmallInt(2)).ToValue(), (value.SmallInt(5)).ToValue()))
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	if value.Truthy(l0) {
		err = t1.SetVal(thread, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue())
		if err.IsNotUndefined() {
			thread.Panic(err)
		}
	}
	err = t1.SetVal(thread, (sym0).ToValue(), (value.NewNativeArrayListWithElements[value.Symbol](0, sym1)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"with native static elements and if modifiers": {
			input: `
				var a: Int? = nil
				b := { "foo" => 5.2, "bar" => 124.99 if a, "baz" => 0.01 }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashMap // var b: Std::HashMap[Std::String, Std::Float]
	_ = l1
	var t1 *vm.NativeHashMap[value.String, value.Float]
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1 = vm.NewNativeHashMapWithElementsAndTotalCapacity[value.String, value.Float](3+0, value.MakeNativePair(value.String("foo"), value.Float(5.2)))
	if value.Truthy(l0) {
		t1.Set(value.String("bar"), value.Float(124.99))
	}
	t1.Set(value.String("baz"), value.Float(0.01))
	l1 = t1
}
`,
		},
		"with native static elements and if modifiers in immutable local": {
			input: `
				var a: Int? = nil
				val b = { "foo" => 5.2, "bar" => 124.99 if a, "baz" => 0.01 }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var t1 *vm.NativeHashMap[value.String, value.Float]
	_ = t1
	var l1 *vm.NativeHashMap[value.String, value.Float] // var b: Std::HashMap[Std::String, Std::Float]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1 = vm.NewNativeHashMapWithElementsAndTotalCapacity[value.String, value.Float](3+0, value.MakeNativePair(value.String("foo"), value.Float(5.2)))
	if value.Truthy(l0) {
		t1.Set(value.String("bar"), value.Float(124.99))
	}
	t1.Set(value.String("baz"), value.Float(0.01))
	l1 = t1
}
`,
		},
		"with static elements, if modifiers and capacity": {
			input: `
				var a: Int? = nil
				b := { 1 => 5 if a, 6 => [:foo] }:45
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(61, 3, 39), P(62, 3, 40)),
					"capacity cannot be specified in collection literals with conditional elements or loops",
				),
			},
		},
		"with static elements and unless modifiers": {
			input: `
				var a: Int? = nil
				b := { 1 => 5 unless a, 9 => [:foo] }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashMap // var b: Std::HashMap[Std::Int, Std::Int | Std::ArrayList[Std::Symbol]]
	_ = l1
	var t1 *vm.HashMapOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1, err = vm.NewHashMapOfValueWithCapacityAndElements(thread, 2+0)
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	if value.Falsy(l0) {
		err = t1.SetVal(thread, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue())
		if err.IsNotUndefined() {
			thread.Panic(err)
		}
	}
	err = t1.SetVal(thread, (value.SmallInt(9)).ToValue(), (value.NewNativeArrayListWithElements[value.Symbol](0, sym0)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		// TODO: constructors
		// 		"with dynamic elements and if modifiers": {
		// 			input: `
		// 				var a: Int? = nil
		// 				{ Object() => 5 if a, 0 => [:foo] }
		// 			`,
		// 			want: vm.NewBytecodeFunctionNoParams(
		// 				mainSymbol,
		// 				[]byte{
		// 					byte(bytecode.PREP_LOCALS8), 1,
		// 					byte(bytecode.NIL),
		// 					byte(bytecode.SET_LOCAL_1),
		// 					byte(bytecode.UNDEFINED),
		// 					byte(bytecode.LOAD_VALUE_0),
		// 					byte(bytecode.NEW_HASH_MAP8), 0,
		// 					byte(bytecode.GET_LOCAL_1),
		// 					byte(bytecode.JUMP_UNLESS), 0, 9,
		// 					byte(bytecode.GET_CONST8), 1,
		// 					byte(bytecode.INSTANTIATE8), 0,
		// 					byte(bytecode.INT_5),
		// 					byte(bytecode.MAP_SET),
		// 					byte(bytecode.JUMP), 0, 0,
		// 					byte(bytecode.INT_0),
		// 					byte(bytecode.LOAD_VALUE_2),
		// 					byte(bytecode.COPY),
		// 					byte(bytecode.MAP_SET),
		// 					byte(bytecode.RETURN),
		// 				},
		// 				L(P(0, 1, 1), P(62, 3, 40)),
		// 				bytecode.LineInfoList{
		// 					bytecode.NewLineInfo(1, 2),
		// 					bytecode.NewLineInfo(2, 2),
		// 					bytecode.NewLineInfo(3, 22),
		// 				},
		// 				[]value.Value{
		// 					value.Ref(value.NewHashMap(2)),
		// 					value.ToSymbol("Std::Object").ToValue(),
		// 					value.Ref(&value.ArrayList{
		// 						value.ToSymbol("foo").ToValue(),
		// 					}),
		// 				},
		// 			),
		// 		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoHashRecord(t *testing.T) {
	tests := goTestTable{
		"empty": {
			input: "a := %{}",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var hshrec0 = vm.MustNewHashRecordOfValueWithElements(nil)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashRecord // var a: Std::HashRecord[any, any]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = hshrec0
}
`,
		},
		"shorthand local": {
			input: `
				foo := 3
				a := %{ foo }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var foo: Std::Int
	_ = l0
	var l1 vm.HashRecord // var a: Std::HashRecord[Std::Symbol, Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(3)).ToValue()
	l1 = vm.MakeNativeKeyHashRecordFromMap(map[value.Symbol]value.Value{sym0: l0})
}
`,
		},
		"shorthand private local": {
			input: `
				_foo := 3
				a := %{ _foo }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("_foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var _foo: Std::Int
	_ = l0
	var l1 vm.HashRecord // var a: Std::HashRecord[Std::Symbol, Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = (value.SmallInt(3)).ToValue()
	l1 = vm.MakeNativeKeyHashRecordFromMap(map[value.Symbol]value.Value{sym0: l0})
}
`,
		},
		"with static elements": {
			input: `a := %{ 1 => 'foo', foo: 5, "bar" => 5.6 }`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")
var hshrec0 = vm.MustNewHashRecordOfValueWithElements(nil, value.MakePairOfValue((value.String("bar")).ToValue(), (value.Float(5.6)).ToValue()), value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue()), value.MakePairOfValue((sym0).ToValue(), (value.SmallInt(5)).ToValue()))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashRecord // var a: Std::HashRecord[Std::Int | Std::Symbol | Std::String, Std::String | Std::Int | Std::Float]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = hshrec0
}
`,
		},
		"with native static elements": {
			input: `a := %{ foo: 5u8, bar: 120u8, baz: 62u8 }`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("bar")
var sym1 = value.ToSymbol("baz")
var sym2 = value.ToSymbol("foo")
var hshrec0 = vm.MakeNativeHashRecordFromMap(map[value.Symbol]value.UInt8{sym0: value.UInt8(120), sym1: value.UInt8(62), sym2: value.UInt8(5)})

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashRecord // var a: Std::HashRecord[Std::Symbol, Std::UInt8]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = hshrec0
}
`,
		},
		"with native static keys": {
			input: `a := %{ 'foo' => 5, 'bar' => 120.7, 'baz' => 'lol' }`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var hshrec0 = vm.MakeNativeKeyHashRecordFromMap(map[value.String]value.Value{value.String("bar"): (value.Float(120.7)).ToValue(), value.String("baz"): (value.String("lol")).ToValue(), value.String("foo"): (value.SmallInt(5)).ToValue()})

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashRecord // var a: Std::HashRecord[Std::String, Std::Int | Std::Float | Std::String]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = hshrec0
}
`,
		},
		// TODO: for in loops
		// 		"with static elements and for loops": {
		// 			input: `%{ 1 => 'foo', i => i ** 2 for i in [1, 2, 3], 2 => 5.6 }`,
		// 			want: vm.NewBytecodeFunctionNoParams(
		// 				mainSymbol,
		// 				[]byte{
		// 					byte(bytecode.PREP_LOCALS8), 2,
		// 					byte(bytecode.LOAD_VALUE_0),
		// 					byte(bytecode.NEW_HASH_RECORD8), 0,
		// 					byte(bytecode.LOAD_VALUE_1),
		// 					byte(bytecode.COPY),
		// 					byte(bytecode.GET_ITERATOR),
		// 					byte(bytecode.SET_LOCAL_1),
		// 					byte(bytecode.GET_LOCAL_1),
		// 					byte(bytecode.FOR_IN_BUILTIN), 0, 9,
		// 					byte(bytecode.SET_LOCAL_2),
		// 					byte(bytecode.GET_LOCAL_2),
		// 					byte(bytecode.GET_LOCAL_2),
		// 					byte(bytecode.INT_2),
		// 					byte(bytecode.EXPONENTIATE_INT),
		// 					byte(bytecode.MAP_SET),
		// 					byte(bytecode.LOOP), 0, 13,
		// 					byte(bytecode.INT_2),
		// 					byte(bytecode.LOAD_VALUE_2),
		// 					byte(bytecode.MAP_SET),
		// 					byte(bytecode.RETURN),
		// 				},
		// 				L(P(0, 1, 1), P(56, 1, 57)),
		// 				bytecode.LineInfoList{
		// 					bytecode.NewLineInfo(1, 26),
		// 				},
		// 				[]value.Value{
		// 					value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
		// 						nil,
		// 						3,
		// 						value.Pair{
		// 							Key:   value.SmallInt(1).ToValue(),
		// 							Value: value.Ref(value.String("foo")),
		// 						},
		// 					)),
		// 					value.Ref(&value.ArrayList{
		// 						value.SmallInt(1).ToValue(),
		// 						value.SmallInt(2).ToValue(),
		// 						value.SmallInt(3).ToValue(),
		// 					}),
		// 					value.Float(5.6).ToValue(),
		// 				},
		// 			),
		// 		},
		"nested static": {
			input: "a := %{ 'foo' => 9, 1 => %{ 'bar' => [7.2] } }",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 vm.HashRecord // var a: Std::HashRecord[Std::String | Std::Int, Std::Int | Std::HashRecord[Std::String, Std::ArrayList[Std::Float]]]
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = vm.MustNewHashRecordOfValueWithElements(nil, value.MakePairOfValue((value.String("foo")).ToValue(), (value.SmallInt(9)).ToValue()), value.MakePairOfValue((value.SmallInt(1)).ToValue(), (vm.MustNewHashRecordOfValueWithElements(nil, value.MakePairOfValue((value.String("bar")).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.Float(7.2)).ToValue())).ToValue()))).ToValue()))
}
`,
		},
		"with static and dynamic elements": {
			input: `
				var a: Int? = nil
				b := %{ 1 => 'foo', 5 => a, 5 => %[:foo] }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Symbol](0, sym0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashRecord // var b: Std::HashRecord[Std::Int, Std::String | Std::ArrayTuple[Std::Symbol] | Std::Int | nil]
	_ = l1
	var t1 *vm.HashRecordOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1, err = vm.NewHashRecordOfValueWithElements(thread, value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.String("foo")).ToValue()), value.MakePairOfValue((value.SmallInt(5)).ToValue(), l0), value.MakePairOfValue((value.SmallInt(5)).ToValue(), (arrtuple0).ToValue()))
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"with static elements and if modifiers": {
			input: `
				var a: Int? = nil
				b := %{ 2 => 5, 1 => 5 if a, a: [:foo] }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("a")
var sym1 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashRecord // var b: Std::HashRecord[Std::Int | Std::Symbol, Std::Int | Std::ArrayList[Std::Symbol]]
	_ = l1
	var t1 *vm.HashRecordOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1, err = vm.NewHashRecordOfValueWithElements(thread, value.MakePairOfValue((value.SmallInt(2)).ToValue(), (value.SmallInt(5)).ToValue()))
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	if value.Truthy(l0) {
		err = t1.SetVal(thread, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue())
		if err.IsNotUndefined() {
			thread.Panic(err)
		}
	}
	err = t1.SetVal(thread, (sym0).ToValue(), (value.NewNativeArrayListWithElements[value.Symbol](0, sym1)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"with native static elements and if modifiers": {
			input: `
				var a: Int? = nil
				b := %{ 2.5 => "foo", 1.0 => "bar" if a, 4.92 => "baz" }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashRecord // var b: Std::HashRecord[Std::Float, Std::String]
	_ = l1
	var t1 vm.NativeHashRecord[value.Float, value.String]
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1 = vm.MakeNativeHashRecordFromMap(map[value.Float]value.String{value.Float(2.5): value.String("foo")})
	if value.Truthy(l0) {
		t1.Set(value.Float(1), value.String("bar"))
	}
	t1.Set(value.Float(4.92), value.String("baz"))
	l1 = t1
}
`,
		},
		"with native static keys and if modifiers": {
			input: `
				var a: Int? = nil
				b := %{ 2.5 => "foo", 1.0 => 190 if a, 4.92 => nil }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashRecord // var b: Std::HashRecord[Std::Float, Std::String | Std::Int | nil]
	_ = l1
	var t1 vm.NativeKeyHashRecord[value.Float]
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1 = vm.MakeNativeKeyHashRecordFromMap(map[value.Float]value.Value{value.Float(2.5): (value.String("foo")).ToValue()})
	if value.Truthy(l0) {
		t1.Set(value.Float(1), (value.SmallInt(190)).ToValue())
	}
	t1.Set(value.Float(4.92), value.Nil)
	l1 = t1
}
`,
		},
		"with static elements and unless modifiers": {
			input: `
				var a: Int? = nil
				b := %{ 1 => 5 unless a, 9 => [:foo] }
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 vm.HashRecord // var b: Std::HashRecord[Std::Int, Std::Int | Std::ArrayList[Std::Symbol]]
	_ = l1
	var t1 *vm.HashRecordOfValue
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Nil
	t1, err = vm.NewHashRecordOfValueWithElements(thread)
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	if value.Falsy(l0) {
		err = t1.SetVal(thread, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue())
		if err.IsNotUndefined() {
			thread.Panic(err)
		}
	}
	err = t1.SetVal(thread, (value.SmallInt(9)).ToValue(), (value.NewNativeArrayListWithElements[value.Symbol](0, sym0)).ToValue())
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		// TODO: constructors
		// 		"with dynamic elements and if modifiers": {
		// 			input: `
		// 				var a: Int? = nil
		// 				%{ Object() => 5 if a, 0 => [:foo] }
		// 			`,
		// 			want: vm.NewBytecodeFunctionNoParams(
		// 				mainSymbol,
		// 				[]byte{
		// 					byte(bytecode.PREP_LOCALS8), 1,
		// 					byte(bytecode.NIL),
		// 					byte(bytecode.SET_LOCAL_1),
		// 					byte(bytecode.UNDEFINED),
		// 					byte(bytecode.NEW_HASH_RECORD8), 0,
		// 					byte(bytecode.GET_LOCAL_1),
		// 					byte(bytecode.JUMP_UNLESS), 0, 9,
		// 					byte(bytecode.GET_CONST8), 0,
		// 					byte(bytecode.INSTANTIATE8), 0,
		// 					byte(bytecode.INT_5),
		// 					byte(bytecode.MAP_SET),
		// 					byte(bytecode.JUMP), 0, 0,
		// 					byte(bytecode.INT_0),
		// 					byte(bytecode.LOAD_VALUE_1),
		// 					byte(bytecode.COPY),
		// 					byte(bytecode.MAP_SET),
		// 					byte(bytecode.RETURN),
		// 				},
		// 				L(P(0, 1, 1), P(63, 3, 41)),
		// 				bytecode.LineInfoList{
		// 					bytecode.NewLineInfo(1, 2),
		// 					bytecode.NewLineInfo(2, 2),
		// 					bytecode.NewLineInfo(3, 21),
		// 				},
		// 				[]value.Value{
		// 					value.ToSymbol("Std::Object").ToValue(),
		// 					value.Ref(&value.ArrayList{
		// 						value.ToSymbol("foo").ToValue(),
		// 					}),
		// 				},
		// 			),
		// 		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoRegex(t *testing.T) {
	tests := goTestTable{
		"empty": {
			input: "a := %//",
			want: `package main

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var regex0 = value.MustCompileRegex("", bitfield.BitField8FromBitFlag(0))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.Regex // var a: Std::Regex
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = regex0
}
`,
		},
		"empty with flags": {
			input: "a := %//imx",
			want: `package main

import (
	"github.com/elk-language/elk/bitfield"
	reflag "github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var regex0 = value.MustCompileRegex("", bitfield.BitField8FromBitFlag(reflag.CaseInsensitiveFlag|reflag.MultilineFlag|reflag.ExtendedFlag))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.Regex // var a: Std::Regex
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = regex0
}
`,
		},
		"with content": {
			input: `a := %/foo \w+ bar/i`,
			want: `package main

import (
	"github.com/elk-language/elk/bitfield"
	reflag "github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var regex0 = value.MustCompileRegex("foo \\w+ bar", bitfield.BitField8FromBitFlag(reflag.CaseInsensitiveFlag))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.Regex // var a: Std::Regex
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = regex0
}
`,
		},
		"reuse the same regex": {
			input: `
				a := %/foo \w+ bar/i
				b := %/foo \w+ bar/i
			`,
			want: `package main

import (
	"github.com/elk-language/elk/bitfield"
	reflag "github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var regex0 = value.MustCompileRegex("foo \\w+ bar", bitfield.BitField8FromBitFlag(reflag.CaseInsensitiveFlag))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 *value.Regex // var a: Std::Regex
	_ = l0
	var l1 *value.Regex // var b: Std::Regex
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = regex0
	l1 = regex0
}
`,
		},
		"with interpolation": {
			input: `
				a := "baz"
				b := %/foo \w+ ${a} bar/i
			`,
			want: `package main

import (
	"github.com/elk-language/elk/bitfield"
	reflag "github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.String // var a: Std::String
	_ = l0
	var l1 *value.Regex // var b: Std::Regex
	_ = l1
	var t1 *value.Regex
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.String("baz")
	t1, err = value.CompileRegexVal(value.String("foo \\w+ ")+l0+value.String(" bar"), bitfield.BitField8FromBitFlag(reflag.CaseInsensitiveFlag))
	l1 = t1
}
`,
		},
		"with primitive interpolation": {
			input: `
				a := 2.5
				b := %/foo \w+ ${a} bar/i
			`,
			want: `package main

import (
	"github.com/elk-language/elk/bitfield"
	reflag "github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 *value.Regex // var b: Std::Regex
	_ = l1
	var t1 *value.Regex
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	l0 = value.Float(2.5)
	t1, err = value.CompileRegexVal(value.String("foo \\w+ ")+(l0).ToString()+value.String(" bar"), bitfield.BitField8FromBitFlag(reflag.CaseInsensitiveFlag))
	l1 = t1
}
`,
		},
		"with complex interpolation": {
			input: `
				a := Time.now
				b := %/foo \w+ ${a} bar/i
			`,
			want: `package main

import (
	"github.com/elk-language/elk/bitfield"
	reflag "github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("now")
var Std_ns_Time_ns_now vm.NativeFunction // Std::Time::now
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("to_string")
var Std_ns_Time_im_to_string vm.NativeFunction // Std::Time.:to_string

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var l0 value.Time // var a: Std::Time
	_ = l0
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Time
	_ = t3
	var l1 *value.Regex // var b: Std::Regex
	_ = l1
	var t4 *value.Regex
	_ = t4
	var t5 []value.Value
	_ = t5
	var t6 value.String
	_ = t6
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	Std_ns_Time_ns_now = vm.MethodToFunc(((value.TimeClass).SingletonClass()).LookupMethod(sym0))
	Std_ns_Time_im_to_string = vm.MethodToFunc((value.TimeClass).LookupMethod(sym2))

	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (value.TimeClass).ToValue()
	thread.AddNativeCallFrame(sym0, sym1, 2)
	t1, err = Std_ns_Time_ns_now(thread, t2) // receiver: &Std::Time, name: now
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t3 = (t1).AsTime()
	l0 = t3
	t5 = value.ResizeNativeArgs(t5, 2)
	t5[0] = (l0).ToValue()
	thread.AddNativeCallFrame(sym2, sym1, 3)
	t1, err = Std_ns_Time_im_to_string(thread, t5) // receiver: Std::Time, name: to_string
	thread.PopNativeCallFrame()
	if err.IsNotUndefined() {
		thread.Panic(err)
	}
	t6 = (t1).AsString()
	t4, err = value.CompileRegexVal(value.String("foo \\w+ ")+t6+value.String(" bar"), bitfield.BitField8FromBitFlag(reflag.CaseInsensitiveFlag))
	l1 = t4
}
`,
		},
		"with compile error": {
			input: `%/foo\y/i`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(5, 1, 6), P(6, 1, 7)),
					`invalid escape sequence: \y`,
				),
			},
		},
		"with compile error from Go": {
			input: ` %/foo{1000000}/i`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(1, 1, 2), P(16, 1, 17)),
					"error parsing regexp: invalid repeat count: `{1000000}`",
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}
