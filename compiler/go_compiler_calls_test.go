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
var sym7 = value.ToSymbol("[]@1")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::[]")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.String) (result value.String, err value.Value) { // method: Foo::[], loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return l0, value.Undefined

}

var sym3 = value.ToSymbol("Foo::[]@1")
var sym4 = value.ToSymbol("to_float")
var fn_method2 vm.NativeFunction // Std::Int.:to_float
func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Float, err value.Value) { // method: Foo::[]@1, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2

	callFrame = thread.AddNativeCallFrame(sym3, sym2, 4)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = l0
	t1, err = fn_method2(thread, t2) // receiver: Std::Int, name: to_float
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
	fn_method2 = vm.MethodToFunc((value.IntClass).LookupMethod(sym4))

	callFrame = thread.AddNativeCallFrame(sym5, sym2, 1)
	defer thread.PopNativeCallFrame()
	callFrame.SetNativeLineNumber(7)
	t1, err = fn_method0(thread, (const0).ToValue(), value.String("lol")) // receiver: Foo, name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
	callFrame.SetNativeLineNumber(8)
	t2, err = fn_method1(thread, (const0).ToValue(), (value.SmallInt(1)).ToValue()) // receiver: Foo, name: []@1
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
	const0 = value.NewModule()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "[]", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsString())
		return (result).ToValue(), err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "[]@1", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], args[1])
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

func TestGoSubscriptSet(t *testing.T) {
	tests := goTestTable{
		"short add set": {
			input: `
				val a = [5.5, 3.2]
				a[1] += 9.1
			`,
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
	var t1 value.Float
	_ = t1
	var err value.Value
	_ = err
	var t2 value.Float
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Float](0, value.Float(5.5), value.Float(3.2))
	callFrame.SetNativeLineNumber(3)
	t1, err = (l0).Get(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t2 = (t1).AddFloat(value.Float(9.1))
	err = (l0).Set(int(value.SmallInt(1)), t2)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"arraylist of value smallint arg": {
			input: `
				val a = [5, 3.2]
				a[1] = 9.1
			`,
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
	t1 = (value.Float(9.1)).ToValue()
	callFrame.SetNativeLineNumber(3)
	err = (l0).Set(int(value.SmallInt(1)), t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"arraylist of value int arg": {
			input: `
				val a = [5, 3.2]
				b := 0
				a[b] = 1
			`,
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
	t1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	err = (l0).Set((l1).AsInt(), t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"arraylist of value bigint arg": {
			input: `
				val a = [5, 3.2]
				val b = 9223372036854775808
				a[b] = 1
			`,
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
	var l0 *value.ArrayListOfValue // var a: Std::ArrayList[Std::Int | Std::Float]
	_ = l0
	var l1 *value.BigInt // var b: 9223372036854775808
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
	l1 = bi0
	t1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	err = (l0).Set(int((l1).ToSmallInt()), t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"arraylist of value strict int arg": {
			input: `
				val a = [5, 3.2]
				b := 0u64
				a[b] = 1
			`,
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
	var l1 value.UInt64 // var b: Std::UInt64
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
	l1 = value.UInt64(0)
	t1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	err = (l0).Set(int(l1), t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"arraylist of value value arg": {
			input: `
				val a = [5, 3.2]
				var b: Int64 | Int32 = 1i64
				a[b] = 1
			`,
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
	t1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	err = (l0).SubscriptSet(l1, t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"native arraylist smallint arg": {
			input: `
				val a = [3.2]
				val b = 0
				a[b] = 1.5
			`,
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
	var l1 value.SmallInt // var b: 0
	_ = l1
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
	l1 = value.SmallInt(0)
	t1 = value.Float(1.5)
	callFrame.SetNativeLineNumber(4)
	err = (l0).Set(int(l1), t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"native arraylist bigint arg": {
			input: `
				val a = [3.2]
				val b = 9223372036854775808
				a[b] = 1.5
			`,
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
	t1 = value.Float(1.5)
	callFrame.SetNativeLineNumber(4)
	err = (l0).Set(int((l1).ToSmallInt()), t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"native arraylist int arg": {
			input: `
				val a = [3.2]
				b := 0
				a[b] = 1.5
			`,
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
	l1 = (value.SmallInt(0)).ToValue()
	t1 = value.Float(1.5)
	callFrame.SetNativeLineNumber(4)
	err = (l0).Set((l1).AsInt(), t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"native arraylist strict int arg": {
			input: `
				val a = [3.2]
				b := 0u64
				a[b] = 1.5
			`,
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
	t1 = value.Float(1.5)
	callFrame.SetNativeLineNumber(4)
	err = (l0).Set(int(l1), t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"native arraylist value arg": {
			input: `
				val a = [3.2]
				var b: Int64 | UInt64 = 0i64
				a[b] = 1.5
			`,
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
	l1 = (value.Int64(0)).ToValue()
	t1 = (value.Float(1.5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	err = (l0).SubscriptSet(l1, t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"list interface int arg": {
			input: `
				var a: List[Int] = [5, 3]
				a[1] = 15
			`,
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
var cc_main_1 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::List[Std::Int]
	_ = l0
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.NewArrayListOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.SmallInt(3)).ToValue())).ToValue()
	t1 = value.ResizeNativeArgs(t1, 4)
	t1[0] = l0
	t1[1] = (value.SmallInt(1)).ToValue()
	t1[2] = (value.SmallInt(15)).ToValue()
	callFrame.SetNativeLineNumber(3)
	_, err = thread.CallMethodByNameWithCache(symbol.OpSubscriptSet, &cc_main_1, t1...) // receiver: Std::List[Std::Int], name: []=
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"hashmap of value": {
			input: `
				val a = { 1 => 5 }
				b := 1
				a[b] = 10
			`,
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
	var l0 *vm.HashMapOfValue // var a: Std::HashMap[Std::Int, Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::Int
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
	l0 = vm.MustNewHashMapOfValueWithCapacityAndElements(nil, 0, value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue()))
	l1 = (value.SmallInt(1)).ToValue()
	t1 = (value.SmallInt(10)).ToValue()
	callFrame.SetNativeLineNumber(4)
	err = (l0).SetVal(thread, l1, t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"native key hashmap": {
			input: `
				val a = { "foo" => 2 }
				b := "bar"
				a[b] = 3
			`,
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
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.NewNativeKeyHashMapFromMap(map[value.String]value.Value{value.String("foo"): (value.SmallInt(2)).ToValue()})
	l1 = value.String("bar")
	t1 = (value.SmallInt(3)).ToValue()
	(l0).Set(l1, t1)
}
`,
		},
		"native hashmap": {
			input: `
				val a = { "foo" => 2.5 }
				b := "bar"
				a[b] = 3.5
			`,
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
	var t1 value.Float
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.NewNativeHashMapFromMap(map[value.String]value.Float{value.String("foo"): value.Float(2.5)})
	l1 = value.String("bar")
	t1 = value.Float(3.5)
	(l0).Set(l1, t1)
}
`,
		},
		"used result": {
			input: `
				val a = [5, 3.2]
				b := a[0] = 9
			`,
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
	t1 = (value.SmallInt(9)).ToValue()
	callFrame.SetNativeLineNumber(3)
	err = (l0).Set(int(value.SmallInt(0)), t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"custom overload": {
			input: `
				module Foo
					overload def []=(a: String, b: String); end
					overload def []=(a: Int, b: Int); end
				end
				Foo[1] = 15
				Foo["lol"] = "foo"
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("[]=@1")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::[]=")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.String, l1 value.String) (result value.Value, err value.Value) { // method: Foo::[]=, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.Nil, value.Undefined

}

var sym3 = value.ToSymbol("Foo::[]=@1")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.Value) (result value.Value, err value.Value) { // method: Foo::[]=@1, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym3, sym2, 4)
	defer thread.PopNativeCallFrame()
	return value.Nil, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
	defer thread.PopNativeCallFrame()
	callFrame.SetNativeLineNumber(6)
	_, err = fn_method1(thread, (const0).ToValue(), (value.SmallInt(1)).ToValue(), (value.SmallInt(15)).ToValue()) // receiver: Foo, name: []=@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	callFrame.SetNativeLineNumber(7)
	_, err = fn_method0(thread, (const0).ToValue(), value.String("lol"), value.String("foo")) // receiver: Foo, name: []=
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
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

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "[]=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsString(), (args[2]).AsString())
		return result, err
	}, vm.DefWithParameters(2))
	vm.Def(&class.MethodContainer, "[]=@1", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], args[1], args[2])
		return result, err
	}, vm.DefWithParameters(2))
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

func TestGoAttributeSet(t *testing.T) {
	tests := goTestTable{
		"set attribute": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo = 3
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(3)).ToValue(), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym3, sym2, 4)
	defer thread.PopNativeCallFrame()
	return l0, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (const0).ToValue()
	callFrame.SetNativeLineNumber(7)
	_, err = fn_method1(thread, l0, (value.SmallInt(3)).ToValue()) // receiver: Bar, name: foo=
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
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

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Bar
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},
		"short add set": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo += 3
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(3)).ToValue(), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym3, sym2, 4)
	defer thread.PopNativeCallFrame()
	return l0, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (const0).ToValue()
	callFrame.SetNativeLineNumber(7)
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, value.AddInts(t1, (value.SmallInt(3)).ToValue())) // receiver: Bar, name: foo=
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
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

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Bar
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], args[1])
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

func TestGoAttributeGet(t *testing.T) {
	tests := goTestTable{
		"get attribute": {
			input: `
				module Bar
					def foo: Int then 3
				end
				a := Bar
				b := a.foo
			`,
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
var sym5 = value.ToSymbol("foo")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(3)).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
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
	l0 = (const0).ToValue()
	callFrame.SetNativeLineNumber(6)
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
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
	const0 = value.NewModule()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Bar
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
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
