package compiler_test

import (
	"testing"
)

func TestGoBitwiseAnd(t *testing.T) {
	tests := goTestTable{
		"resolve static AND": {
			input: "a := 23 & 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(2)).ToValue()
}
`,
		},
		"resolve static nested AND": {
			input: "a := 23 & 15 & 46",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(6)).ToValue()
}
`,
		},

		"and smallint smallint": {
			input: `
				val a = 23
				val b = 10
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.SmallInt // var b: 10
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.SmallInt(10)
	l2 = ((l0).BitwiseAndSmallInt(l1)).ToValue()
}
`,
		},
		"and smallint bigint": {
			input: `
				val a = 23
				c := a & 18446744073709551616
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var c: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).BitwiseAndBigInt(bi0)
}
`,
		},
		"and smallint int": {
			input: `
				val a = 23
				b := 5
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).BitwiseAndInt(l1)
}
`,
		},

		"and bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 10
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.SmallInt // var b: 10
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(10)
	l2 = (l0).BitwiseAndSmallInt(l1)
}
`,
		},
		"and bigint bigint": {
			input: `
				val a = 18446744073709551616
				c := a & 18446744073709551616
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var c: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (l0).BitwiseAndBigInt(bi0)
}
`,
		},
		"and bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).BitwiseAndInt(l1)
}
`,
		},

		"and int64": {
			input: `
				a := 23i64
				b := 10i64
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(23)
	l1 = value.Int64(10)
	l2 = (l0) & (l1)
}
`,
		},
		"and int32": {
			input: `
				a := 23i32
				b := 10i32
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(23)
	l1 = value.Int32(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and int16": {
			input: `
				a := 23i16
				b := 10i16
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(23)
	l1 = value.Int16(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and int8": {
			input: `
				a := 23i8
				b := 10i8
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(23)
	l1 = value.Int8(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and uint64": {
			input: `
				a := 23u64
				b := 10u64
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(23)
	l1 = value.UInt64(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and uint32": {
			input: `
				a := 23u32
				b := 10u32
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(23)
	l1 = value.UInt32(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and uint16": {
			input: `
				a := 23u16
				b := 10u16
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(23)
	l1 = value.UInt16(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and uint8": {
			input: `
				a := 23u8
				b := 10u8
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(23)
	l1 = value.UInt8(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and uint": {
			input: `
				a := 23u
				b := 10u
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(23)
	l1 = value.UInt(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and ints": {
			input: `
				a := 23
				b := 10
				c := a & b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	l2 = value.BitwiseAndInts(l0, l1)
}
`,
		},
		"and optimised value": {
			input: `
				module Foo
					def &(other: Int): Int
						5 & other
					end
				end

				a := Foo
				b := 10
				c := a & b
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

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::&")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__and_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::&, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).BitwiseAndInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	callFrame.SetNativeLineNumber(10)
	t1, err = Foo_ns__and_(thread, l0, l1) // receiver: Foo, name: &
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "&", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__and_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},
		"and unoptimised value": {
			input: `
				module Foo
					def &(other: Int): Int
						5 & other
					end
				end

				var a: Foo | Int = Foo
				b := 10
				c := a & b
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

var sym3 = value.ToSymbol("main")
var cc_main_1 = &value.CallCache{}

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::&")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__and_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::&, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).BitwiseAndInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = l1
	callFrame.SetNativeLineNumber(10)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpAnd, &cc_main_1, t2...) // receiver: Foo | Std::Int, name: &
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "&", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__and_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
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

func TestGoBitwiseAndNot(t *testing.T) {
	tests := goTestTable{
		"resolve static AND NOT": {
			input: "a := 23 &~ 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(21)).ToValue()
}
`,
		},
		"resolve static nested AND NOT": {
			input: "a := 23 &~ 15 &~ 46",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(16)).ToValue()
}
`,
		},

		"and not smallint smallint": {
			input: `
				val a = 23
				val b = 10
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.SmallInt // var b: 10
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.SmallInt(10)
	l2 = ((l0).BitwiseAndNotSmallInt(l1)).ToValue()
}
`,
		},
		"and not smallint bigint": {
			input: `
				val a = 23
				c := a &~ 18446744073709551616
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var c: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).BitwiseAndNotBigInt(bi0)
}
`,
		},
		"and not smallint int": {
			input: `
				val a = 23
				b := 5
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).BitwiseAndNotInt(l1)
}
`,
		},

		"and not bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 10
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.SmallInt // var b: 10
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(10)
	l2 = (l0).BitwiseAndNotSmallInt(l1)
}
`,
		},
		"and not bigint bigint": {
			input: `
				val a = 18446744073709551616
				c := a &~ 18446744073709551616
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var c: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (l0).BitwiseAndNotBigInt(bi0)
}
`,
		},
		"and not bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).BitwiseAndNotInt(l1)
}
`,
		},

		"and not int64": {
			input: `
				a := 23i64
				b := 10i64
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(23)
	l1 = value.Int64(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not int32": {
			input: `
				a := 23i32
				b := 10i32
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(23)
	l1 = value.Int32(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not int16": {
			input: `
				a := 23i16
				b := 10i16
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(23)
	l1 = value.Int16(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not int8": {
			input: `
				a := 23i8
				b := 10i8
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(23)
	l1 = value.Int8(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not uint64": {
			input: `
				a := 23u64
				b := 10u64
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(23)
	l1 = value.UInt64(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not uint32": {
			input: `
				a := 23u32
				b := 10u32
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(23)
	l1 = value.UInt32(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not uint16": {
			input: `
				a := 23u16
				b := 10u16
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(23)
	l1 = value.UInt16(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not uint8": {
			input: `
				a := 23u8
				b := 10u8
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(23)
	l1 = value.UInt8(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not uint": {
			input: `
				a := 23u
				b := 10u
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(23)
	l1 = value.UInt(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not ints": {
			input: `
				a := 23
				b := 10
				c := a &~ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	l2 = value.BitwiseAndNotInts(l0, l1)
}
`,
		},
		"and not optimised value": {
			input: `
				module Foo
					def &~(other: Int): Int
						5 &~ other
					end
				end

				a := Foo
				b := 10
				c := a &~ b
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

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::&~")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__andnot_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::&~, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).BitwiseAndNotInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	callFrame.SetNativeLineNumber(10)
	t1, err = Foo_ns__andnot_(thread, l0, l1) // receiver: Foo, name: &~
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "&~", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__andnot_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},
		"and not unoptimised value": {
			input: `
				module Foo
					def &~(other: Int): Int
						5 &~ other
					end
				end

				var a: Foo | Int = Foo
				b := 10
				c := a &~ b
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

var sym3 = value.ToSymbol("main")
var cc_main_1 = &value.CallCache{}

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::&~")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__andnot_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::&~, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).BitwiseAndNotInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = l1
	callFrame.SetNativeLineNumber(10)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpAndNot, &cc_main_1, t2...) // receiver: Foo | Std::Int, name: &~
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "&~", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__andnot_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
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

func TestGoBitwiseOr(t *testing.T) {
	tests := goTestTable{
		"resolve static OR": {
			input: "a := 23 | 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(31)).ToValue()
}
`,
		},
		"resolve static nested OR": {
			input: "a := 23 | 15 | 46",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(63)).ToValue()
}
`,
		},

		"or smallint smallint": {
			input: `
				val a = 23
				val b = 10
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.SmallInt // var b: 10
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.SmallInt(10)
	l2 = ((l0).BitwiseOrSmallInt(l1)).ToValue()
}
`,
		},
		"or smallint bigint": {
			input: `
				val a = 23
				c := a | 18446744073709551616
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var c: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).BitwiseOrBigInt(bi0)
}
`,
		},
		"or smallint int": {
			input: `
				val a = 23
				b := 5
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).BitwiseOrInt(l1)
}
`,
		},

		"or bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 10
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.SmallInt // var b: 10
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(10)
	l2 = (l0).BitwiseOrSmallInt(l1)
}
`,
		},
		"or bigint bigint": {
			input: `
				val a = 18446744073709551616
				c := a | 18446744073709551616
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var c: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (l0).BitwiseOrBigInt(bi0)
}
`,
		},
		"or bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).BitwiseOrInt(l1)
}
`,
		},

		"or int64": {
			input: `
				a := 23i64
				b := 10i64
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(23)
	l1 = value.Int64(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or int32": {
			input: `
				a := 23i32
				b := 10i32
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(23)
	l1 = value.Int32(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or int16": {
			input: `
				a := 23i16
				b := 10i16
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(23)
	l1 = value.Int16(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or int8": {
			input: `
				a := 23i8
				b := 10i8
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(23)
	l1 = value.Int8(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or uint64": {
			input: `
				a := 23u64
				b := 10u64
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(23)
	l1 = value.UInt64(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or uint32": {
			input: `
				a := 23u32
				b := 10u32
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(23)
	l1 = value.UInt32(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or uint16": {
			input: `
				a := 23u16
				b := 10u16
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(23)
	l1 = value.UInt16(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or uint8": {
			input: `
				a := 23u8
				b := 10u8
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(23)
	l1 = value.UInt8(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or uint": {
			input: `
				a := 23u
				b := 10u
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(23)
	l1 = value.UInt(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or ints": {
			input: `
				a := 23
				b := 10
				c := a | b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	l2 = value.BitwiseOrInts(l0, l1)
}
`,
		},
		"or optimised value": {
			input: `
				module Foo
					def |(other: Int): Int
						5 | other
					end
				end

				a := Foo
				b := 10
				c := a | b
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

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::|")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__or_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::|, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).BitwiseOrInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	callFrame.SetNativeLineNumber(10)
	t1, err = Foo_ns__or_(thread, l0, l1) // receiver: Foo, name: |
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "|", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__or_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},
		"or unoptimised value": {
			input: `
				module Foo
					def |(other: Int): Int
						5 | other
					end
				end

				var a: Foo | Int = Foo
				b := 10
				c := a | b
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

var sym3 = value.ToSymbol("main")
var cc_main_1 = &value.CallCache{}

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::|")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__or_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::|, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).BitwiseOrInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = l1
	callFrame.SetNativeLineNumber(10)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpOr, &cc_main_1, t2...) // receiver: Foo | Std::Int, name: |
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "|", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__or_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
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

func TestGoBitwiseXor(t *testing.T) {
	tests := goTestTable{
		"resolve static XOR": {
			input: "a := 23 ^ 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(29)).ToValue()
}
`,
		},
		"resolve static nested XOR": {
			input: "a := 23 ^ 15 ^ 46",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(54)).ToValue()
}
`,
		},

		"xor smallint smallint": {
			input: `
				val a = 23
				val b = 10
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.SmallInt // var b: 10
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.SmallInt(10)
	l2 = ((l0).BitwiseXorSmallInt(l1)).ToValue()
}
`,
		},
		"xor smallint bigint": {
			input: `
				val a = 23
				c := a ^ 18446744073709551616
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var c: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).BitwiseXorBigInt(bi0)
}
`,
		},
		"xor smallint int": {
			input: `
				val a = 23
				b := 5
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).BitwiseXorInt(l1)
}
`,
		},

		"xor bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 10
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.SmallInt // var b: 10
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(10)
	l2 = (l0).BitwiseXorSmallInt(l1)
}
`,
		},
		"xor bigint bigint": {
			input: `
				val a = 18446744073709551616
				c := a ^ 18446744073709551616
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var c: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (l0).BitwiseXorBigInt(bi0)
}
`,
		},
		"xor bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).BitwiseXorInt(l1)
}
`,
		},

		"xor int64": {
			input: `
				a := 23i64
				b := 10i64
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(23)
	l1 = value.Int64(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor int32": {
			input: `
				a := 23i32
				b := 10i32
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(23)
	l1 = value.Int32(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor int16": {
			input: `
				a := 23i16
				b := 10i16
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(23)
	l1 = value.Int16(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor int8": {
			input: `
				a := 23i8
				b := 10i8
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(23)
	l1 = value.Int8(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor uint64": {
			input: `
				a := 23u64
				b := 10u64
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(23)
	l1 = value.UInt64(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor uint32": {
			input: `
				a := 23u32
				b := 10u32
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(23)
	l1 = value.UInt32(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor uint16": {
			input: `
				a := 23u16
				b := 10u16
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(23)
	l1 = value.UInt16(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor uint8": {
			input: `
				a := 23u8
				b := 10u8
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(23)
	l1 = value.UInt8(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor uint": {
			input: `
				a := 23u
				b := 10u
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(23)
	l1 = value.UInt(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor ints": {
			input: `
				a := 23
				b := 10
				c := a ^ b
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

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	l2 = value.BitwiseXorInts(l0, l1)
}
`,
		},
		"xor optimised value": {
			input: `
				module Foo
					def ^(other: Int): Int
						5 ^ other
					end
				end

				a := Foo
				b := 10
				c := a ^ b
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

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::^")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__xor_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::^, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).BitwiseXorInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	callFrame.SetNativeLineNumber(10)
	t1, err = Foo_ns__xor_(thread, l0, l1) // receiver: Foo, name: ^
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "^", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__xor_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},
		"xor unoptimised value": {
			input: `
				module Foo
					def ^(other: Int): Int
						5 ^ other
					end
				end

				var a: Int | Foo = Foo
				b := 10
				c := a ^ b
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

var sym3 = value.ToSymbol("main")
var cc_main_1 = &value.CallCache{}

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::^")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__xor_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::^, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).BitwiseXorInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int | Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = l1
	callFrame.SetNativeLineNumber(10)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpXor, &cc_main_1, t2...) // receiver: Std::Int | Foo, name: ^
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "^", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__xor_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
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
