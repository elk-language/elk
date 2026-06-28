package compiler_test

import (
	"testing"
)

func TestGoGreater(t *testing.T) {
	tests := goTestTable{
		"resolve static greater": {
			input: "a := 23 > 10",
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
	var l0 value.Bool // var a: Std::Bool
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.True
}
`,
		},
		"greater ints": {
			input: `
				a := 23
				b := a > 10
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = value.Bool(value.GreaterThanInts(l0, (value.SmallInt(10)).ToValue()))
}
`,
		},
		"greater int": {
			input: `
				a := 23
				b := a > 10.5
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.GreaterThanInt(l0, (value.Float(10.5)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.ToBool(t1)
}
`,
		},

		"greater smallint smallint": {
			input: `
						val a = 23
						b := a > 10
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).GreaterThanSmallInt(value.SmallInt(10)))
}
`,
		},
		"greater smallint bigint": {
			input: `
				val a = 23
				b := a > 18446744073709551616
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).GreaterThanBigInt(bi0))
}
`,
		},
		"greater smallint float": {
			input: `
						val a = 23
						b := a > 2.5
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).GreaterThanFloat(value.Float(2.5)))
}
`,
		},
		"greater smallint bigfloat": {
			input: `
				val a = 23
				b := a > 2.5bf
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).GreaterThanBigFloat(bf0))
}
`,
		},
		"greater smallint int": {
			input: `
				val a = 23
				b := 5
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).GreaterThanInt(l1))
}
`,
		},
		"greater smallint value": {
			input: `
				val a = 23
				var b: Int | Float = 5
				c := a > b
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
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).GreaterThanInt(l1))
}
`,
		},

		"greater bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 5
				c := a > b
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
	var l1 value.SmallInt // var b: 5
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(5)
	l2 = value.Bool((l0).GreaterThanSmallInt(l1))
}
`,
		},
		"greater bigint float": {
			input: `
				val a = 18446744073709551616
				val b = 5.5
				c := a > b
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
	var l1 value.Float // var b: 5.5
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.Float(5.5)
	l2 = value.Bool((l0).GreaterThanFloat(l1))
}
`,
		},
		"greater bigint bigfloat": {
			input: `
				val a = 18446744073709551616
				val b = 5.5bf
				c := a > b
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
var bf0 = value.ParseBigFloatPanic("5.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigFloat // var b: 5.5bf
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bf0
	l2 = value.Bool((l0).GreaterThanBigFloat(l1))
}
`,
		},
		"greater bigint bigint": {
			input: `
				val a = 18446744073709551616
				val b = 18446744073709551616
				c := a > b
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
	var l1 *value.BigInt // var b: 18446744073709551616
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bi0
	l2 = value.Bool((l0).GreaterThanBigInt(l1))
}
`,
		},
		"greater bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).GreaterThanInt(l1))
}
`,
		},
		"greater bigint value": {
			input: `
				var a: Int | Float = 10
				b := 5
				c := a > b
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
	var l0 value.Value // var a: Std::Int | Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.GreaterThan(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = value.Bool(t1)
}
`,
		},

		"greater int64": {
			input: `
				a := 6i64
				b := 5i64
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(6)
	l1 = value.Int64(5)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater int32": {
			input: `
				a := 6i32
				b := 5i32
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(6)
	l1 = value.Int32(5)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater int16": {
			input: `
				a := 6i16
				b := 5i16
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(6)
	l1 = value.Int16(5)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater int8": {
			input: `
				a := 6i8
				b := 5i8
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(6)
	l1 = value.Int8(5)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater uint": {
			input: `
				a := 6u
				b := 5u
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(6)
	l1 = value.UInt(5)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater uint64": {
			input: `
				a := 6u64
				b := 5u64
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(6)
	l1 = value.UInt64(5)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater uint32": {
			input: `
				a := 6u32
				b := 5u32
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(6)
	l1 = value.UInt32(5)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater uint16": {
			input: `
				a := 6u16
				b := 5u16
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(6)
	l1 = value.UInt16(5)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater uint8": {
			input: `
				a := 6u8
				b := 5u8
				c := a > b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(6)
	l1 = value.UInt8(5)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater float smallint": {
			input: `
				a := 2.5
				val b = 20
				c := a > b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.SmallInt // var b: 20
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.SmallInt(20)
	l2 = value.Bool((l0).GreaterThanSmallInt(l1))
}
`,
		},
		"greater float bigint": {
			input: `
				a := 2.5
				c := a > 18446744073709551616
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Bool // var c: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Bool((l0).GreaterThanBigInt(bi0))
}
`,
		},
		"greater float float": {
			input: `
				a := 2.5
				b := 0.1
				c := a > b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Float(0.1)
	l2 = value.Bool((l0).GreaterThanFloat(l1))
}
`,
		},
		"greater float bigfloat": {
			input: `
				a := 2.5
				b := 0.1bf
				c := a > b
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
var bf0 = value.ParseBigFloatPanic("0.1")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = bf0
	l2 = value.Bool((l0).GreaterThanBigFloat(l1))
}
`,
		},
		"greater float int": {
			input: `
				a := 2.5
				b := 1
				c := a > b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).GreaterThanInt(l1))
}
`,
		},
		"greater float value": {
			input: `
				var a: Float | Int = 2.5
				b := 1
				c := a > b
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
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.GreaterThan(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = value.Bool(t1)
}
`,
		},

		"greater float64": {
			input: `
				a := 2.5f64
				b := 1f64
				c := a > b
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
	var l0 value.Float64 // var a: Std::Float64
	_ = l0
	var l1 value.Float64 // var b: Std::Float64
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float64(2.5)
	l1 = value.Float64(1)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater float32": {
			input: `
				a := 2.5f32
				b := 1f32
				c := a > b
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
	var l0 value.Float32 // var a: Std::Float32
	_ = l0
	var l1 value.Float32 // var b: Std::Float32
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float32(2.5)
	l1 = value.Float32(1)
	l2 = value.Bool((l0) > (l1))
}
`,
		},

		"greater bigfloat smallint": {
			input: `
				a := 2.5bf
				val b = 1
				c := a > b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.SmallInt // var b: 1
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.SmallInt(1)
	l2 = value.Bool((l0).GreaterThanSmallInt(l1))
}
`,
		},

		"greater bigfloat bigint": {
			input: `
				a := 2.5bf
				c := a > 18446744073709551616
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
var bf0 = value.ParseBigFloatPanic("2.5")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Bool // var c: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Bool((l0).GreaterThanBigInt(bi0))
}
`,
		},
		"greater bigfloat float": {
			input: `
				a := 2.5bf
				val b = 1.0
				c := a > b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Float // var b: 1.0
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Float(1)
	l2 = value.Bool((l0).GreaterThanFloat(l1))
}
`,
		},
		"greater bigfloat bigfloat": {
			input: `
				a := 2.5bf
				val b = 1.0bf
				c := a > b
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
var bf0 = value.ParseBigFloatPanic("2.5")
var bf1 = value.ParseBigFloatPanic("1.0")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var b: 1.0bf
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = bf1
	l2 = value.Bool((l0).GreaterThanBigFloat(l1))
}
`,
		},
		"greater bigfloat int": {
			input: `
				a := 2.5bf
				b := 1
				c := a > b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).GreaterThanInt(l1))
}
`,
		},
		"greater bigfloat value": {
			input: `
				a := 2.5bf
				var b: Int | Float = 1
				c := a > b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).GreaterThanInt(l1))
}
`,
		},

		"greater builtin values": {
			input: `
				var a: Float | Int = 2.5
				b := a > 0.1
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
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.GreaterThan(l0, (value.Float(0.1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.Bool(t1)
}
`,
		},

		"greater optimised value": {
			input: `
				module Foo
					def >(other: Foo | Int): bool
						true
					end
				end
				a := Foo
				b := a > 5
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

var sym1 = value.ToSymbol("Foo::>")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__gt_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Bool, err value.Value) { // method: Foo::>, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.True, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
	var t1 value.Bool
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
	callFrame.SetNativeLineNumber(8)
	t1, err = Foo_ns__gt_(thread, l0, (value.SmallInt(5)).ToValue()) // receiver: Foo, name: >
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
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
	vm.Def(&class.MethodContainer, ">", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__gt_(thread, args[0], args[1])
		return (result).ToValue(), err
	}, vm.DefWithParameters(1))
}
`,
		},

		"greater unoptimised value": {
			input: `
				module Foo
					def >(other: Foo | CoercibleNumeric): bool
						true
					end
				end
				var a: Int | Foo = Foo
				b := a > 5
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

var sym1 = value.ToSymbol("Foo::>")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__gt_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Bool, err value.Value) { // method: Foo::>, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.True, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int | Foo
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
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
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(8)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpGreaterThan, &cc_main_1, t2...) // receiver: Std::Int | Foo, name: >
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.ToBool(t1)
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
	vm.Def(&class.MethodContainer, ">", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__gt_(thread, args[0], args[1])
		return (result).ToValue(), err
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

func TestGoGreaterEqual(t *testing.T) {
	tests := goTestTable{
		"resolve static greater equal": {
			input: "a := 23 >= 10",
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
	var l0 value.Bool // var a: Std::Bool
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.True
}
`,
		},
		"greater equal ints": {
			input: `
				a := 23
				b := a >= 10
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(10)).ToValue()))
}
`,
		},
		"greater equal int": {
			input: `
				a := 23
				b := a >= 10.5
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.GreaterThanEqualInt(l0, (value.Float(10.5)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.ToBool(t1)
}
`,
		},

		"greater equal smallint smallint": {
			input: `
				val a = 23
				b := a >= 10
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).GreaterThanEqualSmallInt(value.SmallInt(10)))
}
`,
		},
		"greater equal smallint bigint": {
			input: `
				val a = 23
				b := a >= 18446744073709551616
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).GreaterThanEqualBigInt(bi0))
}
`,
		},
		"greater equal smallint float": {
			input: `
				val a = 23
				b := a >= 2.5
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).GreaterThanEqualFloat(value.Float(2.5)))
}
`,
		},
		"greater equal smallint bigfloat": {
			input: `
				val a = 23
				b := a >= 2.5bf
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).GreaterThanEqualBigFloat(bf0))
}
`,
		},
		"greater equal smallint int": {
			input: `
				val a = 23
				b := 5
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).GreaterThanEqualInt(l1))
}
`,
		},
		"greater equal smallint value": {
			input: `
				val a = 23
				var b: Int | Float = 5
				c := a >= b
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
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).GreaterThanEqualInt(l1))
}
`,
		},

		"greater equal bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 5
				c := a >= b
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
	var l1 value.SmallInt // var b: 5
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(5)
	l2 = value.Bool((l0).GreaterThanEqualSmallInt(l1))
}
`,
		},
		"greater equal bigint float": {
			input: `
				val a = 18446744073709551616
				val b = 5.5
				c := a >= b
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
	var l1 value.Float // var b: 5.5
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.Float(5.5)
	l2 = value.Bool((l0).GreaterThanEqualFloat(l1))
}
`,
		},
		"greater equal bigint bigfloat": {
			input: `
				val a = 18446744073709551616
				val b = 5.5bf
				c := a >= b
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
var bf0 = value.ParseBigFloatPanic("5.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigFloat // var b: 5.5bf
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bf0
	l2 = value.Bool((l0).GreaterThanEqualBigFloat(l1))
}
`,
		},
		"greater equal bigint bigint": {
			input: `
				val a = 18446744073709551616
				val b = 18446744073709551616
				c := a >= b
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
	var l1 *value.BigInt // var b: 18446744073709551616
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bi0
	l2 = value.Bool((l0).GreaterThanEqualBigInt(l1))
}
`,
		},
		"greater equal bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).GreaterThanEqualInt(l1))
}
`,
		},
		"greater equal bigint value": {
			input: `
				var a: Int | Float = 10
				b := 5
				c := a >= b
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
	var l0 value.Value // var a: Std::Int | Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.GreaterThanEqual(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = value.Bool(t1)
}
`,
		},

		"greater equal int64": {
			input: `
				a := 6i64
				b := 5i64
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(6)
	l1 = value.Int64(5)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal int32": {
			input: `
				a := 6i32
				b := 5i32
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(6)
	l1 = value.Int32(5)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal int16": {
			input: `
				a := 6i16
				b := 5i16
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(6)
	l1 = value.Int16(5)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal int8": {
			input: `
				a := 6i8
				b := 5i8
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(6)
	l1 = value.Int8(5)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal uint": {
			input: `
				a := 6u
				b := 5u
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(6)
	l1 = value.UInt(5)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal uint64": {
			input: `
				a := 6u64
				b := 5u64
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(6)
	l1 = value.UInt64(5)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal uint32": {
			input: `
				a := 6u32
				b := 5u32
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(6)
	l1 = value.UInt32(5)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal uint16": {
			input: `
				a := 6u16
				b := 5u16
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(6)
	l1 = value.UInt16(5)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal uint8": {
			input: `
				a := 6u8
				b := 5u8
				c := a >= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(6)
	l1 = value.UInt8(5)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal float smallint": {
			input: `
				a := 2.5
				val b = 20
				c := a >= b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.SmallInt // var b: 20
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.SmallInt(20)
	l2 = value.Bool((l0).GreaterThanEqualSmallInt(l1))
}
`,
		},
		"greater equal float bigint": {
			input: `
				a := 2.5
				c := a >= 18446744073709551616
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Bool // var c: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Bool((l0).GreaterThanEqualBigInt(bi0))
}
`,
		},
		"greater equal float float": {
			input: `
				a := 2.5
				b := 0.1
				c := a >= b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Float(0.1)
	l2 = value.Bool((l0).GreaterThanEqualFloat(l1))
}
`,
		},
		"greater equal float bigfloat": {
			input: `
				a := 2.5
				b := 0.1bf
				c := a >= b
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
var bf0 = value.ParseBigFloatPanic("0.1")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = bf0
	l2 = value.Bool((l0).GreaterThanEqualBigFloat(l1))
}
`,
		},
		"greater equal float int": {
			input: `
				a := 2.5
				b := 1
				c := a >= b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).GreaterThanEqualInt(l1))
}
`,
		},
		"greater equal float value": {
			input: `
				var a: Float | Int = 2.5
				b := 1
				c := a >= b
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
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.GreaterThanEqual(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = value.Bool(t1)
}
`,
		},

		"greater equal float64": {
			input: `
				a := 2.5f64
				b := 1f64
				c := a >= b
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
	var l0 value.Float64 // var a: Std::Float64
	_ = l0
	var l1 value.Float64 // var b: Std::Float64
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float64(2.5)
	l1 = value.Float64(1)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal float32": {
			input: `
				a := 2.5f32
				b := 1f32
				c := a >= b
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
	var l0 value.Float32 // var a: Std::Float32
	_ = l0
	var l1 value.Float32 // var b: Std::Float32
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float32(2.5)
	l1 = value.Float32(1)
	l2 = value.Bool((l0) >= (l1))
}
`,
		},

		"greater equal bigfloat smallint": {
			input: `
				a := 2.5bf
				val b = 1
				c := a >= b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.SmallInt // var b: 1
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.SmallInt(1)
	l2 = value.Bool((l0).GreaterThanEqualSmallInt(l1))
}
`,
		},

		"greater equal bigfloat bigint": {
			input: `
				a := 2.5bf
				c := a >= 18446744073709551616
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
var bf0 = value.ParseBigFloatPanic("2.5")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Bool // var c: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Bool((l0).GreaterThanEqualBigInt(bi0))
}
`,
		},
		"greater equal bigfloat float": {
			input: `
				a := 2.5bf
				val b = 1.0
				c := a >= b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Float // var b: 1.0
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Float(1)
	l2 = value.Bool((l0).GreaterThanEqualFloat(l1))
}
`,
		},
		"greater equal bigfloat bigfloat": {
			input: `
				a := 2.5bf
				val b = 1.0bf
				c := a >= b
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
var bf0 = value.ParseBigFloatPanic("2.5")
var bf1 = value.ParseBigFloatPanic("1.0")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var b: 1.0bf
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = bf1
	l2 = value.Bool((l0).GreaterThanEqualBigFloat(l1))
}
`,
		},
		"greater equal bigfloat int": {
			input: `
				a := 2.5bf
				b := 1
				c := a >= b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).GreaterThanEqualInt(l1))
}
`,
		},
		"greater equal bigfloat value": {
			input: `
				a := 2.5bf
				var b: Int | Float = 1
				c := a >= b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).GreaterThanEqualInt(l1))
}
`,
		},

		"greater equal builtin values": {
			input: `
				var a: Float | Int = 2.5
				b := a >= 0.1
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
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.GreaterThanEqual(l0, (value.Float(0.1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.Bool(t1)
}
`,
		},

		"greater equal optimised value": {
			input: `
				module Foo
					def >=(other: Foo | Int): bool
						true
					end
				end
				a := Foo
				b := a >= 5
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

var sym1 = value.ToSymbol("Foo::>=")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__gte_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Bool, err value.Value) { // method: Foo::>=, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.True, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
	var t1 value.Bool
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
	callFrame.SetNativeLineNumber(8)
	t1, err = Foo_ns__gte_(thread, l0, (value.SmallInt(5)).ToValue()) // receiver: Foo, name: >=
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
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
	vm.Def(&class.MethodContainer, ">=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__gte_(thread, args[0], args[1])
		return (result).ToValue(), err
	}, vm.DefWithParameters(1))
}
`,
		},

		"greater equal unoptimised value": {
			input: `
				module Foo
					def >=(other: Foo | CoercibleNumeric): bool
						true
					end
				end
				var a: Int | Foo = Foo
				b := a >= 5
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

var sym1 = value.ToSymbol("Foo::>=")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__gte_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Bool, err value.Value) { // method: Foo::>=, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.True, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int | Foo
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
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
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(8)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpGreaterThanEqual, &cc_main_1, t2...) // receiver: Std::Int | Foo, name: >=
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.ToBool(t1)
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
	vm.Def(&class.MethodContainer, ">=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__gte_(thread, args[0], args[1])
		return (result).ToValue(), err
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

func TestGoLess(t *testing.T) {
	tests := goTestTable{
		"resolve static less": {
			input: "a := 23 < 10",
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
	var l0 value.Bool // var a: Std::Bool
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.False
}
`,
		},
		"less ints": {
			input: `
				a := 23
				b := a < 10
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = value.Bool(value.LessThanInts(l0, (value.SmallInt(10)).ToValue()))
}
`,
		},
		"less int": {
			input: `
				a := 23
				b := a < 10.5
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.LessThanInt(l0, (value.Float(10.5)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.ToBool(t1)
}
`,
		},

		"less smallint smallint": {
			input: `
				val a = 23
				b := a < 10
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).LessThanSmallInt(value.SmallInt(10)))
}
`,
		},
		"less smallint bigint": {
			input: `
				val a = 23
				b := a < 18446744073709551616
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).LessThanBigInt(bi0))
}
`,
		},
		"less smallint float": {
			input: `
				val a = 23
				b := a < 2.5
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).LessThanFloat(value.Float(2.5)))
}
`,
		},
		"less smallint bigfloat": {
			input: `
				val a = 23
				b := a < 2.5bf
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).LessThanBigFloat(bf0))
}
`,
		},
		"less smallint int": {
			input: `
				val a = 23
				b := 5
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).LessThanInt(l1))
}
`,
		},
		"less smallint value": {
			input: `
				val a = 23
				var b: Int | Float = 5
				c := a < b
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
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).LessThanInt(l1))
}
`,
		},

		"less bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 5
				c := a < b
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
	var l1 value.SmallInt // var b: 5
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(5)
	l2 = value.Bool((l0).LessThanSmallInt(l1))
}
`,
		},
		"less bigint float": {
			input: `
				val a = 18446744073709551616
				val b = 5.5
				c := a < b
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
	var l1 value.Float // var b: 5.5
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.Float(5.5)
	l2 = value.Bool((l0).LessThanFloat(l1))
}
`,
		},
		"less bigint bigfloat": {
			input: `
				val a = 18446744073709551616
				val b = 5.5bf
				c := a < b
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
var bf0 = value.ParseBigFloatPanic("5.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigFloat // var b: 5.5bf
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bf0
	l2 = value.Bool((l0).LessThanBigFloat(l1))
}
`,
		},
		"less bigint bigint": {
			input: `
				val a = 18446744073709551616
				val b = 18446744073709551616
				c := a < b
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
	var l1 *value.BigInt // var b: 18446744073709551616
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bi0
	l2 = value.Bool((l0).LessThanBigInt(l1))
}
`,
		},
		"less bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).LessThanInt(l1))
}
`,
		},
		"less bigint value": {
			input: `
				var a: Int | Float = 10
				b := 5
				c := a < b
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
	var l0 value.Value // var a: Std::Int | Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.LessThan(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = value.Bool(t1)
}
`,
		},

		"less int64": {
			input: `
				a := 6i64
				b := 5i64
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(6)
	l1 = value.Int64(5)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less int32": {
			input: `
				a := 6i32
				b := 5i32
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(6)
	l1 = value.Int32(5)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less int16": {
			input: `
				a := 6i16
				b := 5i16
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(6)
	l1 = value.Int16(5)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less int8": {
			input: `
				a := 6i8
				b := 5i8
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(6)
	l1 = value.Int8(5)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less uint": {
			input: `
				a := 6u
				b := 5u
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(6)
	l1 = value.UInt(5)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less uint64": {
			input: `
				a := 6u64
				b := 5u64
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(6)
	l1 = value.UInt64(5)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less uint32": {
			input: `
				a := 6u32
				b := 5u32
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(6)
	l1 = value.UInt32(5)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less uint16": {
			input: `
				a := 6u16
				b := 5u16
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(6)
	l1 = value.UInt16(5)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less uint8": {
			input: `
				a := 6u8
				b := 5u8
				c := a < b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(6)
	l1 = value.UInt8(5)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less float smallint": {
			input: `
				a := 2.5
				val b = 20
				c := a < b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.SmallInt // var b: 20
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.SmallInt(20)
	l2 = value.Bool((l0).LessThanSmallInt(l1))
}
`,
		},
		"less float bigint": {
			input: `
				a := 2.5
				c := a < 18446744073709551616
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Bool // var c: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Bool((l0).LessThanBigInt(bi0))
}
`,
		},
		"less float float": {
			input: `
				a := 2.5
				b := 0.1
				c := a < b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Float(0.1)
	l2 = value.Bool((l0).LessThanFloat(l1))
}
`,
		},
		"less float bigfloat": {
			input: `
				a := 2.5
				b := 0.1bf
				c := a < b
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
var bf0 = value.ParseBigFloatPanic("0.1")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = bf0
	l2 = value.Bool((l0).LessThanBigFloat(l1))
}
`,
		},
		"less float int": {
			input: `
				a := 2.5
				b := 1
				c := a < b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).LessThanInt(l1))
}
`,
		},
		"less float value": {
			input: `
				var a: Float | Int = 2.5
				b := 1
				c := a < b
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
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.LessThan(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = value.Bool(t1)
}
`,
		},

		"less float64": {
			input: `
				a := 2.5f64
				b := 1f64
				c := a < b
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
	var l0 value.Float64 // var a: Std::Float64
	_ = l0
	var l1 value.Float64 // var b: Std::Float64
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float64(2.5)
	l1 = value.Float64(1)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less float32": {
			input: `
				a := 2.5f32
				b := 1f32
				c := a < b
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
	var l0 value.Float32 // var a: Std::Float32
	_ = l0
	var l1 value.Float32 // var b: Std::Float32
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float32(2.5)
	l1 = value.Float32(1)
	l2 = value.Bool((l0) < (l1))
}
`,
		},

		"less bigfloat smallint": {
			input: `
				a := 2.5bf
				val b = 1
				c := a < b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.SmallInt // var b: 1
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.SmallInt(1)
	l2 = value.Bool((l0).LessThanSmallInt(l1))
}
`,
		},

		"less bigfloat bigint": {
			input: `
				a := 2.5bf
				c := a < 18446744073709551616
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
var bf0 = value.ParseBigFloatPanic("2.5")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Bool // var c: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Bool((l0).LessThanBigInt(bi0))
}
`,
		},
		"less bigfloat float": {
			input: `
				a := 2.5bf
				val b = 1.0
				c := a < b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Float // var b: 1.0
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Float(1)
	l2 = value.Bool((l0).LessThanFloat(l1))
}
`,
		},
		"less bigfloat bigfloat": {
			input: `
				a := 2.5bf
				val b = 1.0bf
				c := a < b
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
var bf0 = value.ParseBigFloatPanic("2.5")
var bf1 = value.ParseBigFloatPanic("1.0")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var b: 1.0bf
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = bf1
	l2 = value.Bool((l0).LessThanBigFloat(l1))
}
`,
		},
		"less bigfloat int": {
			input: `
				a := 2.5bf
				b := 1
				c := a < b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).LessThanInt(l1))
}
`,
		},
		"less bigfloat value": {
			input: `
				a := 2.5bf
				var b: Int | Float = 1
				c := a < b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).LessThanInt(l1))
}
`,
		},

		"less builtin values": {
			input: `
				var a: Float | Int = 2.5
				b := a < 0.1
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
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.LessThan(l0, (value.Float(0.1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.Bool(t1)
}
`,
		},

		"less optimised value": {
			input: `
				module Foo
					def <(other: Foo | Int): bool
						true
					end
				end
				a := Foo
				b := a < 5
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

var sym1 = value.ToSymbol("Foo::<")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__lt_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Bool, err value.Value) { // method: Foo::<, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.True, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
	var t1 value.Bool
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
	callFrame.SetNativeLineNumber(8)
	t1, err = Foo_ns__lt_(thread, l0, (value.SmallInt(5)).ToValue()) // receiver: Foo, name: <
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
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
	vm.Def(&class.MethodContainer, "<", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__lt_(thread, args[0], args[1])
		return (result).ToValue(), err
	}, vm.DefWithParameters(1))
}
`,
		},

		"less unoptimised value": {
			input: `
				module Foo
					def <(other: Foo | CoercibleNumeric): bool
						true
					end
				end
				var a: Int | Foo = Foo
				b := a < 5
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

var sym1 = value.ToSymbol("Foo::<")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__lt_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Bool, err value.Value) { // method: Foo::<, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.True, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int | Foo
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
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
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(8)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpLessThan, &cc_main_1, t2...) // receiver: Std::Int | Foo, name: <
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.ToBool(t1)
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
	vm.Def(&class.MethodContainer, "<", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__lt_(thread, args[0], args[1])
		return (result).ToValue(), err
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

func TestGoLessEqual(t *testing.T) {
	tests := goTestTable{
		"resolve static less equal": {
			input: "a := 23 <= 10",
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
	var l0 value.Bool // var a: Std::Bool
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.False
}
`,
		},
		"less equal ints": {
			input: `
				a := 23
				b := a <= 10
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = value.Bool(value.LessThanEqualInts(l0, (value.SmallInt(10)).ToValue()))
}
`,
		},
		"less equal int": {
			input: `
				a := 23
				b := a <= 10.5
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.LessThanEqualInt(l0, (value.Float(10.5)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.ToBool(t1)
}
`,
		},

		"less equal smallint smallint": {
			input: `
				val a = 23
				b := a <= 10
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).LessThanEqualSmallInt(value.SmallInt(10)))
}
`,
		},
		"less equal smallint bigint": {
			input: `
				val a = 23
				b := a <= 18446744073709551616
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).LessThanEqualBigInt(bi0))
}
`,
		},
		"less equal smallint float": {
			input: `
				val a = 23
				b := a <= 2.5
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
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).LessThanEqualFloat(value.Float(2.5)))
}
`,
		},
		"less equal smallint bigfloat": {
			input: `
				val a = 23
				b := a <= 2.5bf
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.Bool((l0).LessThanEqualBigFloat(bf0))
}
`,
		},
		"less equal smallint int": {
			input: `
				val a = 23
				b := 5
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).LessThanEqualInt(l1))
}
`,
		},
		"less equal smallint value": {
			input: `
				val a = 23
				var b: Int | Float = 5
				c := a <= b
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
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).LessThanEqualInt(l1))
}
`,
		},

		"less equal bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 5
				c := a <= b
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
	var l1 value.SmallInt // var b: 5
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(5)
	l2 = value.Bool((l0).LessThanEqualSmallInt(l1))
}
`,
		},
		"less equal bigint float": {
			input: `
				val a = 18446744073709551616
				val b = 5.5
				c := a <= b
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
	var l1 value.Float // var b: 5.5
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.Float(5.5)
	l2 = value.Bool((l0).LessThanEqualFloat(l1))
}
`,
		},
		"less equal bigint bigfloat": {
			input: `
				val a = 18446744073709551616
				val b = 5.5bf
				c := a <= b
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
var bf0 = value.ParseBigFloatPanic("5.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigFloat // var b: 5.5bf
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bf0
	l2 = value.Bool((l0).LessThanEqualBigFloat(l1))
}
`,
		},
		"less equal bigint bigint": {
			input: `
				val a = 18446744073709551616
				val b = 18446744073709551616
				c := a <= b
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
	var l1 *value.BigInt // var b: 18446744073709551616
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bi0
	l2 = value.Bool((l0).LessThanEqualBigInt(l1))
}
`,
		},
		"less equal bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = value.Bool((l0).LessThanEqualInt(l1))
}
`,
		},
		"less equal bigint value": {
			input: `
				var a: Int | Float = 10
				b := 5
				c := a <= b
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
	var l0 value.Value // var a: Std::Int | Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.LessThanEqual(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = value.Bool(t1)
}
`,
		},

		"greater equal int64": {
			input: `
				a := 6i64
				b := 5i64
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(6)
	l1 = value.Int64(5)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal int32": {
			input: `
				a := 6i32
				b := 5i32
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(6)
	l1 = value.Int32(5)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal int16": {
			input: `
				a := 6i16
				b := 5i16
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(6)
	l1 = value.Int16(5)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal int8": {
			input: `
				a := 6i8
				b := 5i8
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(6)
	l1 = value.Int8(5)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal uint": {
			input: `
				a := 6u
				b := 5u
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(6)
	l1 = value.UInt(5)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal uint64": {
			input: `
				a := 6u64
				b := 5u64
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(6)
	l1 = value.UInt64(5)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal uint32": {
			input: `
				a := 6u32
				b := 5u32
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(6)
	l1 = value.UInt32(5)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal uint16": {
			input: `
				a := 6u16
				b := 5u16
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(6)
	l1 = value.UInt16(5)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal uint8": {
			input: `
				a := 6u8
				b := 5u8
				c := a <= b
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
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(6)
	l1 = value.UInt8(5)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal float smallint": {
			input: `
				a := 2.5
				val b = 20
				c := a <= b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.SmallInt // var b: 20
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.SmallInt(20)
	l2 = value.Bool((l0).LessThanEqualSmallInt(l1))
}
`,
		},
		"less equal float bigint": {
			input: `
				a := 2.5
				c := a <= 18446744073709551616
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Bool // var c: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Bool((l0).LessThanEqualBigInt(bi0))
}
`,
		},
		"less equal float float": {
			input: `
				a := 2.5
				b := 0.1
				c := a <= b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Float(0.1)
	l2 = value.Bool((l0).LessThanEqualFloat(l1))
}
`,
		},
		"less equal float bigfloat": {
			input: `
				a := 2.5
				b := 0.1bf
				c := a <= b
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
var bf0 = value.ParseBigFloatPanic("0.1")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = bf0
	l2 = value.Bool((l0).LessThanEqualBigFloat(l1))
}
`,
		},
		"less equal float int": {
			input: `
				a := 2.5
				b := 1
				c := a <= b
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
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).LessThanEqualInt(l1))
}
`,
		},
		"less equal float value": {
			input: `
				var a: Float | Int = 2.5
				b := 1
				c := a <= b
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
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.LessThanEqual(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = value.Bool(t1)
}
`,
		},

		"less equal float64": {
			input: `
				a := 2.5f64
				b := 1f64
				c := a <= b
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
	var l0 value.Float64 // var a: Std::Float64
	_ = l0
	var l1 value.Float64 // var b: Std::Float64
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float64(2.5)
	l1 = value.Float64(1)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal float32": {
			input: `
				a := 2.5f32
				b := 1f32
				c := a <= b
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
	var l0 value.Float32 // var a: Std::Float32
	_ = l0
	var l1 value.Float32 // var b: Std::Float32
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float32(2.5)
	l1 = value.Float32(1)
	l2 = value.Bool((l0) <= (l1))
}
`,
		},

		"less equal bigfloat smallint": {
			input: `
				a := 2.5bf
				val b = 1
				c := a <= b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.SmallInt // var b: 1
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.SmallInt(1)
	l2 = value.Bool((l0).LessThanEqualSmallInt(l1))
}
`,
		},

		"less equal bigfloat bigint": {
			input: `
				a := 2.5bf
				c := a <= 18446744073709551616
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
var bf0 = value.ParseBigFloatPanic("2.5")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Bool // var c: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Bool((l0).LessThanEqualBigInt(bi0))
}
`,
		},
		"less equal bigfloat float": {
			input: `
				a := 2.5bf
				val b = 1.0
				c := a <= b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Float // var b: 1.0
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Float(1)
	l2 = value.Bool((l0).LessThanEqualFloat(l1))
}
`,
		},
		"less equal bigfloat bigfloat": {
			input: `
				a := 2.5bf
				val b = 1.0bf
				c := a <= b
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
var bf0 = value.ParseBigFloatPanic("2.5")
var bf1 = value.ParseBigFloatPanic("1.0")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var b: 1.0bf
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = bf1
	l2 = value.Bool((l0).LessThanEqualBigFloat(l1))
}
`,
		},
		"less equal bigfloat int": {
			input: `
				a := 2.5bf
				b := 1
				c := a <= b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).LessThanEqualInt(l1))
}
`,
		},
		"less equal bigfloat value": {
			input: `
				a := 2.5bf
				var b: Int | Float = 1
				c := a <= b
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
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Bool // var c: Std::Bool
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = value.Bool((l0).LessThanEqualInt(l1))
}
`,
		},

		"less equal builtin values": {
			input: `
				var a: Float | Int = 2.5
				b := a <= 0.1
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
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Bool // var b: Std::Bool
	_ = l1
	var t1 bool
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.LessThanEqual(l0, (value.Float(0.1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.Bool(t1)
}
`,
		},

		"less equal optimised value": {
			input: `
				module Foo
					def <=(other: Foo | Int): bool
						true
					end
				end
				a := Foo
				b := a <= 5
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

var sym1 = value.ToSymbol("Foo::<=")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__lte_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Bool, err value.Value) { // method: Foo::<=, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.True, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
	var t1 value.Bool
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
	callFrame.SetNativeLineNumber(8)
	t1, err = Foo_ns__lte_(thread, l0, (value.SmallInt(5)).ToValue()) // receiver: Foo, name: <=
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
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
	vm.Def(&class.MethodContainer, "<=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__lte_(thread, args[0], args[1])
		return (result).ToValue(), err
	}, vm.DefWithParameters(1))
}
`,
		},

		"less equal unoptimised value": {
			input: `
				module Foo
					def <=(other: Foo | CoercibleNumeric): bool
						true
					end
				end
				var a: Int | Foo = Foo
				b := a <= 5
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

var sym1 = value.ToSymbol("Foo::<=")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__lte_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Bool, err value.Value) { // method: Foo::<=, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.True, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int | Foo
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
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
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(8)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpLessThanEqual, &cc_main_1, t2...) // receiver: Std::Int | Foo, name: <=
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = value.ToBool(t1)
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
	vm.Def(&class.MethodContainer, "<=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__lte_(thread, args[0], args[1])
		return (result).ToValue(), err
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
