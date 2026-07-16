package compiler_test

import (
	"testing"
)

func TestGoSubscript(t *testing.T) {
	tests := goTestTable{
		"static": {
			input: `a := [5, 3][0]`,
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
	l0 = (value.SmallInt(5)).ToValue()
}
`,
		},
		"arraylist of value smallint arg": {
			input: `
				val a = [5, 3.2]
				b := a[1]
			`,
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
	var l0 *value.ArrayListOfValue // var a: Std::ArrayList[Std::Int | Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
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
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.Float(3.2)).ToValue())
	callFrame.SetNativeLineNumber(3)
	t1, err = (l0).Get(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"arraytuple of value smallint arg": {
			input: `
				val a = %[5, 3.2]
				b := a[1]
			`,
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
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.Float(3.2)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.ArrayTupleOfValue // var a: Std::ArrayTuple[Std::Int | Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
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
	l0 = arrtuple0
	callFrame.SetNativeLineNumber(3)
	t1, err = (l0).Get(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"arraylist of value any int arg": {
			input: `
				val a = [5, 3.2]
				var b: Int64 | Int32 = 1i64
				c := a[b]
			`,
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
	var l0 *value.ArrayListOfValue // var a: Std::ArrayList[Std::Int | Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int64 | Std::Int32
	_ = l1
	var l2 value.Value // var c: Std::Int | Std::Float
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.Float(3.2)).ToValue())
	l1 = (value.Int64(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Subscript(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"arraytuple of value any int arg": {
			input: `
				val a = %[5, 3.2]
				var b: Int64 | Int32 = 1i64
				c := a[b]
			`,
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
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.Float(3.2)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.ArrayTupleOfValue // var a: Std::ArrayTuple[Std::Int | Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int64 | Std::Int32
	_ = l1
	var l2 value.Value // var c: Std::Int | Std::Float
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	l1 = (value.Int64(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Subscript(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"arraylist of value int arg": {
			input: `
				val a = [5, 3.2]
				b := 0
				c := a[b]
			`,
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
	var l0 *value.ArrayListOfValue // var a: Std::ArrayList[Std::Int | Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int | Std::Float
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.Float(3.2)).ToValue())
	l1 = (value.SmallInt(0)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Get((l1).AsInt())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"arraytuple of value int arg": {
			input: `
				val a = %[5, 3.2]
				b := 0
				c := a[b]
			`,
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
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.Float(3.2)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.ArrayTupleOfValue // var a: Std::ArrayTuple[Std::Int | Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int | Std::Float
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	l1 = (value.SmallInt(0)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Get((l1).AsInt())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"arraylist of value value arg": {
			input: `
				val a = [5, 3.2]
				var b: Int64 | Int32 = 1i64
				c := a[b]
			`,
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
	var l0 *value.ArrayListOfValue // var a: Std::ArrayList[Std::Int | Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int64 | Std::Int32
	_ = l1
	var l2 value.Value // var c: Std::Int | Std::Float
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.Float(3.2)).ToValue())
	l1 = (value.Int64(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Subscript(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"arraytuple of value value arg": {
			input: `
				val a = %[5, 3.2]
				var b: Int64 | Int32 = 1i64
				c := a[b]
			`,
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
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.Float(3.2)).ToValue())

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.ArrayTupleOfValue // var a: Std::ArrayTuple[Std::Int | Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int64 | Std::Int32
	_ = l1
	var l2 value.Value // var c: Std::Int | Std::Float
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	l1 = (value.Int64(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Subscript(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"native arraylist smallint arg": {
			input: `
				val a = [3.2]
				val b = 2
				c := a[b]
			`,
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
	var l0 *value.NativeArrayList[value.Float] // var a: Std::ArrayList[Std::Float]
	_ = l0
	var l1 value.SmallInt // var b: 2
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var t1 value.Float
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Float](0, value.Float(3.2))
	l1 = value.SmallInt(2)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Get(int(l1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"native arraytuple smallint arg": {
			input: `
				val a = %[3.2]
				val b = 2
				c := a[b]
			`,
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
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Float](0, value.Float(3.2))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayTuple[value.Float] // var a: Std::ArrayTuple[Std::Float]
	_ = l0
	var l1 value.SmallInt // var b: 2
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var t1 value.Float
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	l1 = value.SmallInt(2)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Get(int(l1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"native arraylist bigint arg": {
			input: `
				val a = [3.2]
				val b = 9223372036854775808
				c := a[b]
			`,
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
var bi0 = value.ParseBigIntPanic("9223372036854775808", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Float] // var a: Std::ArrayList[Std::Float]
	_ = l0
	var l1 *value.BigInt // var b: 9223372036854775808
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var t1 value.Float
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Float](0, value.Float(3.2))
	l1 = bi0
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Get(int((l1).ToSmallInt()))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"native arraytuple bigint arg": {
			input: `
				val a = %[3.2]
				val b = 9223372036854775808
				c := a[b]
			`,
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
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Float](0, value.Float(3.2))
var bi0 = value.ParseBigIntPanic("9223372036854775808", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayTuple[value.Float] // var a: Std::ArrayTuple[Std::Float]
	_ = l0
	var l1 *value.BigInt // var b: 9223372036854775808
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var t1 value.Float
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	l1 = bi0
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Get(int((l1).ToSmallInt()))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"native arraylist int arg": {
			input: `
				val a = [3.2]
				b := 5
				c := a[b]
			`,
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
	var l0 *value.NativeArrayList[value.Float] // var a: Std::ArrayList[Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var t1 value.Float
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Float](0, value.Float(3.2))
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Get((l1).AsInt())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"native arraytuple int arg": {
			input: `
				val a = %[3.2]
				b := 5
				c := a[b]
			`,
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
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Float](0, value.Float(3.2))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayTuple[value.Float] // var a: Std::ArrayTuple[Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var t1 value.Float
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Get((l1).AsInt())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"native arraylist strict int arg": {
			input: `
				val a = [3.2]
				b := 0u64
				c := a[b]
			`,
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
	var l0 *value.NativeArrayList[value.Float] // var a: Std::ArrayList[Std::Float]
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var t1 value.Float
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Float](0, value.Float(3.2))
	l1 = value.UInt64(0)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Get(int(l1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"native arraytuple strict int arg": {
			input: `
				val a = %[3.2]
				b := 0u64
				c := a[b]
			`,
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
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Float](0, value.Float(3.2))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayTuple[value.Float] // var a: Std::ArrayTuple[Std::Float]
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var t1 value.Float
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	l1 = value.UInt64(0)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Get(int(l1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"native arraylist value arg": {
			input: `
				val a = [3.2]
				var b: Int64 | UInt64 = 5i64
				c := a[b]
			`,
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
	var l0 *value.NativeArrayList[value.Float] // var a: Std::ArrayList[Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int64 | Std::UInt64
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Float](0, value.Float(3.2))
	l1 = (value.Int64(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Subscript(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = (t1).AsFloat()
}
`,
		},
		"native arraytuple value arg": {
			input: `
				val a = %[3.2]
				var b: Int64 | UInt64 = 5i64
				c := a[b]
			`,
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
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Float](0, value.Float(3.2))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayTuple[value.Float] // var a: Std::ArrayTuple[Std::Float]
	_ = l0
	var l1 value.Value // var b: Std::Int64 | Std::UInt64
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	l1 = (value.Int64(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).Subscript(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = (t1).AsFloat()
}
`,
		},

		"hashmap of value": {
			input: `
				a := { 1 => 5 }
				b := 1
				c := a[b]
			`,
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
	var l0 vm.HashMap // var a: Std::HashMap[Std::Int, Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int?
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.MustNewHashMapOfValueWithCapacityAndElements(nil, 0, value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue()))
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).GetValNil(thread, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"hashrecord of value": {
			input: `
				a := %{ 1 => 5 }
				b := 1
				c := a[b]
			`,
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
var hshrec0 = vm.MustNewHashRecordOfValueWithElements(nil, value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue()))

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.HashRecord // var a: Std::HashRecord[Std::Int, Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int?
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = hshrec0
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).GetValNil(thread, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"native key hashmap": {
			input: `
				val a = { "foo" => 2 }
				b := "bar"
				c := a[b]
			`,
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
	var l0 *vm.NativeKeyHashMap[value.String] // var a: Std::HashMap[Std::String, Std::Int]
	_ = l0
	var l1 value.String // var b: Std::String
	_ = l1
	var l2 value.Value // var c: Std::Int?
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.NewNativeKeyHashMapFromMap(map[value.String]value.Value{value.String("foo"): (value.SmallInt(2)).ToValue()})
	l1 = value.String("bar")
	l2 = (l0).Get(l1)
}
`,
		},
		"native key hashrecord": {
			input: `
				val a = %{ "foo" => 2 }
				b := "bar"
				c := a[b]
			`,
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
var hshrec0 = vm.MakeNativeKeyHashRecordFromMap(map[value.String]value.Value{value.String("foo"): (value.SmallInt(2)).ToValue()})

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.NativeKeyHashRecord[value.String] // var a: Std::HashRecord[Std::String, Std::Int]
	_ = l0
	var l1 value.String // var b: Std::String
	_ = l1
	var l2 value.Value // var c: Std::Int?
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = hshrec0
	l1 = value.String("bar")
	l2 = (l0).Get(l1)
}
`,
		},

		"native hashmap": {
			input: `
				val a = { "foo" => 2.5 }
				b := "bar"
				c := a[b]
			`,
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
	var l0 *vm.NativeHashMap[value.String, value.Float] // var a: Std::HashMap[Std::String, Std::Float]
	_ = l0
	var l1 value.String // var b: Std::String
	_ = l1
	var l2 value.Value // var c: Std::Float?
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.NewNativeHashMapFromMap(map[value.String]value.Float{value.String("foo"): value.Float(2.5)})
	l1 = value.String("bar")
	l2 = ((l0).Get(l1)).ToValue()
}
`,
		},
		"native hashrecord": {
			input: `
				val a = %{ "foo" => 2.5 }
				b := "bar"
				c := a[b]
			`,
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
var hshrec0 = vm.MakeNativeHashRecordFromMap(map[value.String]value.Float{value.String("foo"): value.Float(2.5)})

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.NativeHashRecord[value.String, value.Float] // var a: Std::HashRecord[Std::String, Std::Float]
	_ = l0
	var l1 value.String // var b: Std::String
	_ = l1
	var l2 value.Value // var c: Std::Float?
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = hshrec0
	l1 = value.String("bar")
	l2 = ((l0).Get(l1)).ToValue()
}
`,
		},

		"custom overload": {
			input: `
				module Foo
					overload def [](a: String): String then a
					overload def [](a: Int): Float then a.to_float
				end

				a := Foo["lol"]
				b := Foo[1]
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym5 = value.ToSymbol("main")
var sym6 = value.ToSymbol("[]@1")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym4 = value.ToSymbol("Foo::[]")

func Foo_ns__sub_(thread *vm.Thread, self value.Value, l0 value.String) (result value.String, err value.Value) { // method: Foo::[], loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym4, sym2, 3)
	defer thread.PopNativeCallFrame()
	return l0, value.Undefined

}

var sym1 = value.ToSymbol("Foo::[]@1")
var sym2 = value.ToSymbol("<main>")
var sym3 = value.ToSymbol("to_float")
var Std_ns_Int_im_to_float vm.NativeFunction // Std::Int.:to_float
func Foo_ns__sub__at_1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Float, err value.Value) { // method: Foo::[]@1, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 4)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = l0
	t1, err = Std_ns_Int_im_to_float(thread, t2) // receiver: Std::Int, name: to_float
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	return (t1).AsFloat(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.String // var a: Std::String
	_ = l0
	var t1 value.String
	_ = t1
	var err value.Value
	_ = err
	var l1 value.Float // var b: Std::Float
	_ = l1
	var t2 value.Float
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	Std_ns_Int_im_to_float = vm.MethodToFunc((value.IntClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym5, sym2, 1)
	defer thread.PopNativeCallFrame()
	callFrame.SetNativeLineNumber(7)
	t1, err = Foo_ns__sub_(thread, (Foo).ToValue(), value.String("lol")) // receiver: Foo, name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
	callFrame.SetNativeLineNumber(8)
	t2, err = Foo_ns__sub__at_1(thread, (Foo).ToValue(), (value.SmallInt(1)).ToValue()) // receiver: Foo, name: []@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t2
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
	vm.Def(&class.MethodContainer, "[]", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__sub_(thread, args[0], (args[1]).AsString())
		return (result).ToValue(), err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "[]@1", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__sub__at_1(thread, args[0], args[1])
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
