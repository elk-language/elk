package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestGoSubscript(t *testing.T) {
	tests := goTestTable{
		"static": {
			input: `a := [5, 3][0]`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
		"dynamic nil safe": {
			input: `
				var arr: List[Int]? = [5, 3]
				b := arr?[1]
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var cc_main_1 = &vm.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var arr: Std::List[Std::Int]?
	_ = l0
	var l1 value.Value // var b: Std::Int?
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2
	var t3 []value.Value
	_ = t3
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.NewArrayListOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.SmallInt(3)).ToValue())).ToValue()
	if value.IsNil(l0) {
		t3 = value.ResizeNativeArgs(t3, 3)
		t3[0] = l0
		t3[1] = (value.SmallInt(1)).ToValue()
		callFrame.SetNativeLineNumber(3)
		t2, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_1, t3...) // receiver: Std::List[Std::Int]?, name: []
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t1 = t2
	} else {
		t1 = value.Nil
	}
	l1 = t1
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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym5 = value.ToSymbol("main")
var sym7 = value.ToSymbol("[]@1")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::[]")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.String) (result value.String, err value.Value) { // method: Foo::[], loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

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
	t1, err = fn_method0(thread, (const0).ToValue(), value.String("lol")) // receiver: Foo, name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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

		"increment int": {
			input: `
				val a = [5, 3]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.ArrayListOfValue // var a: Std::ArrayList[Std::Int]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Value
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.SmallInt(3)).ToValue())
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = value.IncrementInt(t2)
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment int64": {
			input: `
				val a = [5i64, 3i64]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Int64] // var a: Std::ArrayList[Std::Int64]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Int64
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Int64
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Int64](0, value.Int64(5), value.Int64(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) + 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment int32": {
			input: `
				val a = [5i32, 3i32]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Int32] // var a: Std::ArrayList[Std::Int32]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Int32
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Int32
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Int32](0, value.Int32(5), value.Int32(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) + 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment int16": {
			input: `
				val a = [5i16, 3i16]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Int16] // var a: Std::ArrayList[Std::Int16]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Int16
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Int16
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Int16](0, value.Int16(5), value.Int16(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) + 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment int8": {
			input: `
				val a = [5i8, 3i8]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Int8] // var a: Std::ArrayList[Std::Int8]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Int8
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Int8
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Int8](0, value.Int8(5), value.Int8(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) + 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment uint64": {
			input: `
				val a = [5u64, 3u64]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.UInt64] // var a: Std::ArrayList[Std::UInt64]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.UInt64
	_ = t2
	var err value.Value
	_ = err
	var t3 value.UInt64
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.UInt64](0, value.UInt64(5), value.UInt64(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) + 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment uint32": {
			input: `
				val a = [5u32, 3u32]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.UInt32] // var a: Std::ArrayList[Std::UInt32]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.UInt32
	_ = t2
	var err value.Value
	_ = err
	var t3 value.UInt32
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.UInt32](0, value.UInt32(5), value.UInt32(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) + 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment uint16": {
			input: `
				val a = [5u16, 3u16]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.UInt16] // var a: Std::ArrayList[Std::UInt16]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.UInt16
	_ = t2
	var err value.Value
	_ = err
	var t3 value.UInt16
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.UInt16](0, value.UInt16(5), value.UInt16(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) + 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment uint8": {
			input: `
				val a = [5u8, 3u8]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.UInt8] // var a: Std::ArrayList[Std::UInt8]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.UInt8
	_ = t2
	var err value.Value
	_ = err
	var t3 value.UInt8
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.UInt8](0, value.UInt8(5), value.UInt8(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) + 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment uint": {
			input: `
				val a = [5u, 3u]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.UInt] // var a: Std::ArrayList[Std::UInt]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.UInt
	_ = t2
	var err value.Value
	_ = err
	var t3 value.UInt
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElementsAndTotalCapacity[value.UInt](2+0, value.UInt(5), value.UInt(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) + 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment char": {
			input: "val a = [`a`, `b`]; a[1]++",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Char] // var a: Std::ArrayList[Std::Char]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Char
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Char
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Char](0, value.Char('a'), value.Char('b'))
	t1 = value.SmallInt(1)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) + 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"increment union": {
			input: `
				val a: ArrayList[Int | Int64] = [5, 3]
				a[1]++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.ArrayListOfValue // var a: Std::ArrayList[Std::Int | Std::Int64]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Value
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.SmallInt(3)).ToValue())
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = value.IncrementVal(t2)
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement int": {
			input: `
				val a = [5, 3]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.ArrayListOfValue // var a: Std::ArrayList[Std::Int]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Value
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.SmallInt(3)).ToValue())
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = value.DecrementInt(t2)
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement int64": {
			input: `
				val a = [5i64, 3i64]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Int64] // var a: Std::ArrayList[Std::Int64]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Int64
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Int64
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Int64](0, value.Int64(5), value.Int64(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) - 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement int32": {
			input: `
				val a = [5i32, 3i32]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Int32] // var a: Std::ArrayList[Std::Int32]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Int32
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Int32
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Int32](0, value.Int32(5), value.Int32(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) - 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement int16": {
			input: `
				val a = [5i16, 3i16]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Int16] // var a: Std::ArrayList[Std::Int16]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Int16
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Int16
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Int16](0, value.Int16(5), value.Int16(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) - 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement int8": {
			input: `
				val a = [5i8, 3i8]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Int8] // var a: Std::ArrayList[Std::Int8]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Int8
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Int8
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Int8](0, value.Int8(5), value.Int8(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) - 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement uint64": {
			input: `
				val a = [5u64, 3u64]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.UInt64] // var a: Std::ArrayList[Std::UInt64]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.UInt64
	_ = t2
	var err value.Value
	_ = err
	var t3 value.UInt64
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.UInt64](0, value.UInt64(5), value.UInt64(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) - 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement uint32": {
			input: `
				val a = [5u32, 3u32]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.UInt32] // var a: Std::ArrayList[Std::UInt32]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.UInt32
	_ = t2
	var err value.Value
	_ = err
	var t3 value.UInt32
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.UInt32](0, value.UInt32(5), value.UInt32(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) - 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement uint16": {
			input: `
				val a = [5u16, 3u16]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.UInt16] // var a: Std::ArrayList[Std::UInt16]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.UInt16
	_ = t2
	var err value.Value
	_ = err
	var t3 value.UInt16
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.UInt16](0, value.UInt16(5), value.UInt16(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) - 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement uint8": {
			input: `
				val a = [5u8, 3u8]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.UInt8] // var a: Std::ArrayList[Std::UInt8]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.UInt8
	_ = t2
	var err value.Value
	_ = err
	var t3 value.UInt8
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.UInt8](0, value.UInt8(5), value.UInt8(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) - 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement uint": {
			input: `
				val a = [5u, 3u]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.UInt] // var a: Std::ArrayList[Std::UInt]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.UInt
	_ = t2
	var err value.Value
	_ = err
	var t3 value.UInt
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElementsAndTotalCapacity[value.UInt](2+0, value.UInt(5), value.UInt(3))
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) - 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement char": {
			input: "val a = [`a`, `b`]; a[1]--",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.NativeArrayList[value.Char] // var a: Std::ArrayList[Std::Char]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Char
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Char
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Char](0, value.Char('a'), value.Char('b'))
	t1 = value.SmallInt(1)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = (t2) - 1
	err = (l0).Set(int(t1), t3)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},

		"decrement union": {
			input: `
				val a: ArrayList[Int | Int64] = [5, 3]
				a[1]--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.ArrayListOfValue // var a: Std::ArrayList[Std::Int | Std::Int64]
	_ = l0
	var t1 value.SmallInt
	_ = t1
	var t2 value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Value
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.SmallInt(3)).ToValue())
	t1 = value.SmallInt(1)
	callFrame.SetNativeLineNumber(3)
	t2, err = (l0).Get(int(t1))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = value.DecrementVal(t2)
	err = (l0).Set(int(t1), t3)
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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var cc_main_1 = &vm.CallCache{}

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("[]=@1")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::[]=")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.String, l1 value.String) (result value.Value, err value.Value) { // method: Foo::[]=, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Nil, value.Undefined

}

var sym3 = value.ToSymbol("Foo::[]=@1")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.Value) (result value.Value, err value.Value) { // method: Foo::[]=@1, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

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
	_, err = fn_method1(thread, (const0).ToValue(), (value.SmallInt(1)).ToValue(), (value.SmallInt(15)).ToValue()) // receiver: Foo, name: []=@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (value.SmallInt(3)).ToValue(), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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

	return (value.SmallInt(3)).ToValue(), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

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
		"increment int": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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

	return (value.SmallInt(3)).ToValue(), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, value.IncrementInt(t1)) // receiver: Bar, name: foo=
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

		"increment union": {
			input: `
				module Bar
					def foo: Int | Int64 then 3
					def foo=(value: Int | Int64); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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

	return (value.SmallInt(3)).ToValue(), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, value.IncrementVal(t1)) // receiver: Bar, name: foo=
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

		"increment int64": {
			input: `
				module Bar
					def foo: Int64 then 3i64
					def foo=(value: Int64); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Int64, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Int64(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Int64) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Int64
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)+1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsInt64())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"increment int32": {
			input: `
				module Bar
					def foo: Int32 then 3i32
					def foo=(value: Int32); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Int32, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Int32(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Int32) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Int32
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)+1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsInt32())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"increment int16": {
			input: `
				module Bar
					def foo: Int16 then 3i16
					def foo=(value: Int16); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Int16, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Int16(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Int16) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Int16
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)+1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsInt16())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"increment int8": {
			input: `
				module Bar
					def foo: Int8 then 3i8
					def foo=(value: Int8); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Int8, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Int8(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Int8) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Int8
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)+1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsInt8())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"increment uint64": {
			input: `
				module Bar
					def foo: UInt64 then 3u64
					def foo=(value: UInt64); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.UInt64, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.UInt64(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.UInt64) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.UInt64
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)+1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsUInt64())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"increment uint32": {
			input: `
				module Bar
					def foo: UInt32 then 3u32
					def foo=(value: UInt32); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.UInt32, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.UInt32(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.UInt32) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.UInt32
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)+1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsUInt32())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"increment uint16": {
			input: `
				module Bar
					def foo: UInt16 then 3u16
					def foo=(value: UInt16); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.UInt16, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.UInt16(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.UInt16) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.UInt16
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)+1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsUInt16())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"increment uint8": {
			input: `
				module Bar
					def foo: UInt8 then 3u8
					def foo=(value: UInt8); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.UInt8, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.UInt8(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.UInt8) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.UInt8
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)+1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsUInt8())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"increment uint": {
			input: `
				module Bar
					def foo: UInt then 3u
					def foo=(value: UInt); end
				end
				a := Bar
				a.foo++
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.UInt, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.UInt(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.UInt) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.UInt
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)+1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsUInt())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"increment char": {
			input: "module Bar; def foo: Char then `a`; def foo=(value: Char); end; end; a := Bar; a.foo++",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Char, err value.Value) { // method: Bar::foo, loc: <main>:1:13
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Char('a'), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Char) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:1:37
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Char
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)+1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsChar())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"decrement int": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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

	return (value.SmallInt(3)).ToValue(), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, value.DecrementInt(t1)) // receiver: Bar, name: foo=
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

		"decrement union": {
			input: `
				module Bar
					def foo: Int | Int64 then 3
					def foo=(value: Int | Int64); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

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

	return (value.SmallInt(3)).ToValue(), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, value.DecrementVal(t1)) // receiver: Bar, name: foo=
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

		"decrement int64": {
			input: `
				module Bar
					def foo: Int64 then 3i64
					def foo=(value: Int64); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Int64, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Int64(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Int64) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Int64
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)-1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsInt64())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"decrement int32": {
			input: `
				module Bar
					def foo: Int32 then 3i32
					def foo=(value: Int32); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Int32, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Int32(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Int32) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Int32
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)-1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsInt32())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"decrement int16": {
			input: `
				module Bar
					def foo: Int16 then 3i16
					def foo=(value: Int16); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Int16, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Int16(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Int16) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Int16
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)-1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsInt16())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"decrement int8": {
			input: `
				module Bar
					def foo: Int8 then 3i8
					def foo=(value: Int8); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Int8, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Int8(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Int8) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Int8
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)-1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsInt8())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"decrement uint64": {
			input: `
				module Bar
					def foo: UInt64 then 3u64
					def foo=(value: UInt64); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.UInt64, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.UInt64(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.UInt64) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.UInt64
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)-1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsUInt64())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"decrement uint32": {
			input: `
				module Bar
					def foo: UInt32 then 3u32
					def foo=(value: UInt32); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.UInt32, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.UInt32(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.UInt32) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.UInt32
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)-1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsUInt32())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"decrement uint16": {
			input: `
				module Bar
					def foo: UInt16 then 3u16
					def foo=(value: UInt16); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.UInt16, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.UInt16(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.UInt16) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.UInt16
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)-1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsUInt16())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"decrement uint8": {
			input: `
				module Bar
					def foo: UInt8 then 3u8
					def foo=(value: UInt8); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.UInt8, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.UInt8(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.UInt8) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.UInt8
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)-1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsUInt8())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"decrement uint": {
			input: `
				module Bar
					def foo: UInt then 3u
					def foo=(value: UInt); end
				end
				a := Bar
				a.foo--
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.UInt, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.UInt(3), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.UInt) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.UInt
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)-1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsUInt())
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"decrement char": {
			input: "module Bar; def foo: Char then `a`; def foo=(value: Char); end; end; a := Bar; a.foo--",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("foo")
var sym7 = value.ToSymbol("foo=")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Char, err value.Value) { // method: Bar::foo, loc: <main>:1:13
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Char('a'), value.Undefined

}

var sym3 = value.ToSymbol("Bar::foo=")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Char) (result value.Value, err value.Value) { // method: Bar::foo=, loc: <main>:1:37
	var callFrame *vm.CallFrame
	_ = callFrame

	return (l0).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Bar
	_ = l0
	var t1 value.Char
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
	t1, err = fn_method0(thread, l0) // receiver: Bar, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0, (t1)-1) // receiver: Bar, name: foo=
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
		return (result).ToValue(), err
	})
	vm.Def(&class.MethodContainer, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], (args[1]).AsChar())
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
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym3 = value.ToSymbol("main")
var sym5 = value.ToSymbol("foo")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

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

func TestGoInstantiate(t *testing.T) {
	tests := goTestTable{
		"without arguments and unused": {
			input: `
				class Foo; end
				::Foo()
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
	defer thread.PopNativeCallFrame()
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
`,
		},
		"without arguments": {
			input: `
				class Foo; end
				a := ::Foo()
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	l0 = const0.CreateInstance()
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
`,
		},
		"without arguments native": {
			input: `
				a := Elk::AST::BreakpointNode()
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("Std::Elk::AST::BreakpointNode")
var const0 *value.Class // Std::Elk::AST::BreakpointNode
var sym3 = value.ToSymbol("#init")
var fn_method0 vm.NativeFunction // Std::Elk::AST::BreakpointNode.:#init

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Elk::AST::BreakpointNode
	_ = l0
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	const0 = (*value.Class)((value.GetConstant(sym2)).Pointer())

	fn_method0 = vm.MethodToFunc((const0).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = const0.CreateInstance()
	t2[1] = value.Undefined
	callFrame.SetNativeLineNumber(2)
	t1, err = fn_method0(thread, t2) // receiver: Std::Elk::AST::BreakpointNode, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
}
`,
		},
		"complex constant": {
			input: `
				module Foo
					module Bar
						class Baz; end
					end
				end
				a := ::Foo::Bar::Baz()
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym5 = value.ToSymbol("main")
var sym6 = value.ToSymbol("<main>")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")
var const1 *value.Module // Bar
var sym1 = value.ToSymbol("Bar")
var sym2 = value.ToSymbol("Foo::Bar")
var const2 value.Value  // Foo::Bar
var const3 *value.Class // Baz
var sym3 = value.ToSymbol("Baz")
var sym4 = value.ToSymbol("Foo::Bar::Baz")
var const4 *value.Class // Foo::Bar::Baz

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo::Bar::Baz
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	const2 = value.GetConstant(sym2)
	const4 = (*value.Class)((value.GetConstant(sym4)).Pointer())

	callFrame = thread.AddNativeCallFrame(sym5, sym6, 1)
	defer thread.PopNativeCallFrame()
	l0 = const4.CreateInstance()
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

	parentNamespace = (const0).ToValue()
	const1 = value.NewModule()
	namespace = value.Ref(const1)
	value.AddConstant(parentNamespace, sym1, namespace)

	parentNamespace = const2
	const3 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const3)
	value.AddConstant(parentNamespace, sym3, namespace)

	class = const4
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
`,
		},
		"with positional arguments and unused": {
			input: `
				class Foo
					init(a: Int, b: String); end
				end
				::Foo(1, 'lol')
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym3 = value.ToSymbol("main")
var sym5 = value.ToSymbol("#init")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.String) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return self, value.Undefined

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
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	_, err = fn_method0(thread, const0.CreateInstance(), (value.SmallInt(1)).ToValue(), value.String("lol")) // receiver: Foo, name: #init
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1], (args[2]).AsString())
		return result, err
	}, vm.DefWithParameters(2))
}
`,
		},
		"with positional arguments": {
			input: `
				class Foo
					init(a: Int, b: String); end
				end
				a := ::Foo(1, 'lol')
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym3 = value.ToSymbol("main")
var sym5 = value.ToSymbol("#init")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.String) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return self, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
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
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	t1, err = fn_method0(thread, const0.CreateInstance(), (value.SmallInt(1)).ToValue(), value.String("lol")) // receiver: Foo, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1], (args[2]).AsString())
		return result, err
	}, vm.DefWithParameters(2))
}
`,
		},
		"with named args": {
			input: `
				class Foo
					init(a: Int, b: String); end
				end
				a := ::Foo(1, b: 'lol')
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym3 = value.ToSymbol("main")
var sym5 = value.ToSymbol("#init")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.String) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return self, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
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
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	t1, err = fn_method0(thread, const0.CreateInstance(), (value.SmallInt(1)).ToValue(), value.String("lol")) // receiver: Foo, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1], (args[2]).AsString())
		return result, err
	}, vm.DefWithParameters(2))
}
`,
		},
		"with duplicated named args": {
			input: `
				class Foo
					init(a: String, b: Int); end
				end
				::Foo(b: 1, a: 'lol', b: 2)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(83, 5, 27), P(86, 5, 30)),
					"duplicated argument `b` in call to `Foo.:#init`",
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
