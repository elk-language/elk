package compiler_test

import (
	"testing"
)

func TestGoSwitch(t *testing.T) {
	tests := goTestTable{
		"with a few literal cases": {
			input: `
			  var a: any = 0
				b := switch a
				case true then "a"
				case false then "b"
				case 0 then "c"
				case 1 then "d"
				end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
var cc_main_2 = &value.CallCache{}
var cc_main_3 = &value.CallCache{}
var cc_main_4 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.True).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.False).ToValue()
	callFrame.SetNativeLineNumber(5)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_2, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("b")).ToValue()
		goto lbl1
	}
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.SmallInt(0)).ToValue()
	callFrame.SetNativeLineNumber(6)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_3, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("c")).ToValue()
		goto lbl1
	}
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(7)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_4, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("d")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"with a few literal cases with ignored value": {
			input: `
			  var a: any = 0
				switch a
				case true then "a"
				case false then "b"
				case 0 then "c"
				case 1 then "d"
				end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
var cc_main_2 = &value.CallCache{}
var cc_main_3 = &value.CallCache{}
var cc_main_4 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
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
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = (value.True).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t2...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t1) {
		goto lbl1
	}
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = (value.False).ToValue()
	callFrame.SetNativeLineNumber(5)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_2, t2...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t1) {
		goto lbl1
	}
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = (value.SmallInt(0)).ToValue()
	callFrame.SetNativeLineNumber(6)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_3, t2...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t1) {
		goto lbl1
	}
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(7)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_4, t2...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t1) {
		goto lbl1
	}
lbl1:
}
`,
		},
		"with else": {
			input: `
				var a: any = 0
				b := switch a
				case true then "a"
				case false then "b"
				else "c"
				end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
var cc_main_2 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.String // var b: Std::String
	_ = l1
	var t1 value.String
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.True).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = value.String("a")
		goto lbl1
	}
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.False).ToValue()
	callFrame.SetNativeLineNumber(5)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_2, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = value.String("b")
		goto lbl1
	}
	t1 = value.String("c")
lbl1:
	l1 = t1
}
`,
		},
		"literal true": {
			input: `
			  var a: any = 0
				b := switch a
				case true then "a"
				end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.True).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal false": {
			input: `
			  var a: any = 0
				b := switch a
				case false then "a"
				end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.False).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},

		"literal nil": {
			input: `
			  var a: any = 0
				b := switch a
				case nil then "a"
				end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = value.Nil
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal string": {
			input: `
			  var a: any = 0
				b := switch a
				case "foo" then "a"
				end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.String("foo")).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal raw string": {
			input: `
			  var a: any = 0
				b := switch a
				case 'foo' then "a"
				end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.String("foo")).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal interpolated string": {
			input: `
			  var a: any = 0
				b := switch a
				case "f${a}" then "a"
				end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
var cc_main_2 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2
	var t3 []value.Value
	_ = t3
	var err value.Value
	_ = err
	var t4 value.Value
	_ = t4
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 2)
	t3[0] = l0
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.L_to_string, &cc_main_1, t3...) // receiver: any, name: to_string
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.String("f") + (t2).AsString()).ToValue()
	t4, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_2, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t4) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal symbol": {
			input: `
	  var a: any = 0
		b := switch a
		case :foo then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("foo")
var cc_main_1 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (sym2).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal interpolated symbol": {
			input: `
	  var a: any = 0
		b := switch a
		case :"f${a}" then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
var cc_main_2 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2
	var t3 []value.Value
	_ = t3
	var err value.Value
	_ = err
	var t4 value.Value
	_ = t4
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 2)
	t3[0] = l0
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.L_to_string, &cc_main_1, t3...) // receiver: any, name: to_string
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = ((value.String("f") + (t2).AsString()).ToSymbol()).ToValue()
	t4, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_2, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t4) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal int": {
			input: `
	  var a: any = 0
		b := switch a
		case 5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal int64": {
			input: `
	  var a: any = 0
		b := switch a
		case 5i64 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.Int64(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal uint64": {
			input: `
	  var a: any = 0
		b := switch a
		case 5u64 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.UInt64(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal int32": {
			input: `
	  var a: any = 0
		b := switch a
		case 5i32 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.Int32(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal uint32": {
			input: `
	  var a: any = 0
		b := switch a
		case 5u32 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.UInt32(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal int16": {
			input: `
	  var a: any = 0
		b := switch a
		case 5i16 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.Int16(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal uint16": {
			input: `
	  var a: any = 0
		b := switch a
		case 5u16 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.UInt16(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal int8": {
			input: `
	  var a: any = 0
		b := switch a
		case 5i8 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.Int8(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal uint8": {
			input: `
	  var a: any = 0
		b := switch a
		case 5u8 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.UInt8(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal float": {
			input: `
	  var a: any = 0
		b := switch a
		case 5.8 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.Float(5.8)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal float64": {
			input: `
	  var a: any = 0
		b := switch a
		case 5.8f64 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.Float64(5.8)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal float32": {
			input: `
	  var a: any = 0
		b := switch a
		case 5.8f32 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.Float32(5.800000190734863)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal negative float32": {
			input: `
	  var a: any = 0
		b := switch a
		case -5.8f32 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.Float32(-5.8)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"literal big float": {
			input: `
	  var a: any = 0
		b := switch a
		case 5.8bf then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("5.8")
var cc_main_1 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (bf0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"root constant": {
			input: `
		const Foo = 3
	  var a: any = 0
		b := switch a
		case ::Foo then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("Foo")
var const0 value.SmallInt // constant: Foo, loc: <main>:2:3
var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")
var cc_main_1 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var namespace value.Value
	_ = namespace
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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

	namespace = (value.RootModule).ToValue()
	value.AddConstant(namespace, sym0, (value.SmallInt(3)).ToValue())
	const0 = value.SmallInt(3)
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (const0).ToValue()
	callFrame.SetNativeLineNumber(5)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"constant lookup": {
			input: `
		module Foo
			const Bar = 3
		end
	  var a: any = 0
		b := switch a
		case ::Foo::Bar then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("Bar")
var const1 value.SmallInt // constant: Foo::Bar, loc: <main>:3:4
var sym2 = value.ToSymbol("main")
var sym3 = value.ToSymbol("<main>")
var cc_main_1 = &value.CallCache{}

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var namespace value.Value
	_ = namespace
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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

	initGlobalEnv()

	namespace = (const0).ToValue()
	value.AddConstant(namespace, sym1, (value.SmallInt(3)).ToValue())
	const1 = value.SmallInt(3)
	callFrame = thread.AddNativeCallFrame(sym2, sym3, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (const1).ToValue()
	callFrame.SetNativeLineNumber(7)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
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
`,
		},
		"negative constant lookup": {
			input: `
		module Foo
			const Bar = 3
		end
	  var a: any = 0
		b := switch a
		case -::Foo::Bar then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("Bar")
var const1 value.SmallInt // constant: Foo::Bar, loc: <main>:3:4
var sym2 = value.ToSymbol("main")
var sym3 = value.ToSymbol("<main>")
var cc_main_1 = &value.CallCache{}

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var namespace value.Value
	_ = namespace
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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

	initGlobalEnv()

	namespace = (const0).ToValue()
	value.AddConstant(namespace, sym1, (value.SmallInt(3)).ToValue())
	const1 = value.SmallInt(3)
	callFrame = thread.AddNativeCallFrame(sym2, sym3, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (const1).NegateVal()
	callFrame.SetNativeLineNumber(7)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
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
`,
		},
		"less pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case < 5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	if value.Bool(value.IsA((value.SmallInt(5)).ToValue(), (l0).Class())) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = l0
		t4[1] = (value.SmallInt(5)).ToValue()
		callFrame.SetNativeLineNumber(4)
		t3, err = thread.CallMethodByNameWithCache(symbol.OpLessThan, &cc_main_1, t4...) // receiver: any, name: <
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t2 = value.ToBool(t3)
	} else {
		t2 = value.False
	}
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"less than root constant": {
			input: `
		const Foo = 5
	  var a: any = 0
		b := switch a
		case < ::Foo then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("Foo")
var const0 value.SmallInt // constant: Foo, loc: <main>:2:3
var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")
var cc_main_1 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var namespace value.Value
	_ = namespace
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	namespace = (value.RootModule).ToValue()
	value.AddConstant(namespace, sym0, (value.SmallInt(5)).ToValue())
	const0 = value.SmallInt(5)
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	if value.Bool(value.IsA((const0).ToValue(), (l0).Class())) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = l0
		t4[1] = (const0).ToValue()
		callFrame.SetNativeLineNumber(5)
		t3, err = thread.CallMethodByNameWithCache(symbol.OpLessThan, &cc_main_1, t4...) // receiver: any, name: <
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t2 = value.ToBool(t3)
	} else {
		t2 = value.False
	}
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"less than negative root constant": {
			input: `
		const Foo = 2
	  var a: any = 0
		b := switch a
		case < -::Foo then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("Foo")
var const0 value.SmallInt // constant: Foo, loc: <main>:2:3
var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")
var cc_main_1 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var namespace value.Value
	_ = namespace
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	namespace = (value.RootModule).ToValue()
	value.AddConstant(namespace, sym0, (value.SmallInt(2)).ToValue())
	const0 = value.SmallInt(2)
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	if value.Bool(value.IsA((const0).NegateVal(), (l0).Class())) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = l0
		t4[1] = (const0).NegateVal()
		callFrame.SetNativeLineNumber(5)
		t3, err = thread.CallMethodByNameWithCache(symbol.OpLessThan, &cc_main_1, t4...) // receiver: any, name: <
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t2 = value.ToBool(t3)
	} else {
		t2 = value.False
	}
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"less equal pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case <= 5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	if value.Bool(value.IsA((value.SmallInt(5)).ToValue(), (l0).Class())) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = l0
		t4[1] = (value.SmallInt(5)).ToValue()
		callFrame.SetNativeLineNumber(4)
		t3, err = thread.CallMethodByNameWithCache(symbol.OpLessThanEqual, &cc_main_1, t4...) // receiver: any, name: <=
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t2 = value.ToBool(t3)
	} else {
		t2 = value.False
	}
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"greater pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case > 5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	if value.Bool(value.IsA((value.SmallInt(5)).ToValue(), (l0).Class())) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = l0
		t4[1] = (value.SmallInt(5)).ToValue()
		callFrame.SetNativeLineNumber(4)
		t3, err = thread.CallMethodByNameWithCache(symbol.OpGreaterThan, &cc_main_1, t4...) // receiver: any, name: >
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t2 = value.ToBool(t3)
	} else {
		t2 = value.False
	}
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"greater equal pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case >= 5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	if value.Bool(value.IsA((value.SmallInt(5)).ToValue(), (l0).Class())) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = l0
		t4[1] = (value.SmallInt(5)).ToValue()
		callFrame.SetNativeLineNumber(4)
		t3, err = thread.CallMethodByNameWithCache(symbol.OpGreaterThanEqual, &cc_main_1, t4...) // receiver: any, name: >=
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t2 = value.ToBool(t3)
	} else {
		t2 = value.False
	}
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"equal pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case == 5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"equal regex pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case == %/fo+/ then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var regex0 = value.MustCompileRegex("fo+", bitfield.BitField8FromBitFlag(0))
var cc_main_1 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (regex0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"equal local pattern": {
			input: `
	  var a: any = 0
		b := 2
		c := switch a
		case == b then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::String | nil
	_ = l2
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
	l0 = (value.SmallInt(0)).ToValue()
	l1 = (value.SmallInt(2)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = l1
	callFrame.SetNativeLineNumber(5)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l2 = t1
}
`,
		},
		"not equal pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case != 5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t2)) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"lax equal pattern": {
			input: `
	  var a: any = 0
		b := 2
		c := switch a
		case =~ b then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::String | nil
	_ = l2
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
	l0 = (value.SmallInt(0)).ToValue()
	l1 = (value.SmallInt(2)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = l1
	callFrame.SetNativeLineNumber(5)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpLaxEqual, &cc_main_1, t3...) // receiver: any, name: =~
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l2 = t1
}
`,
		},
		"lax not equal pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case !~ 5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = l0
	t3[1] = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpLaxEqual, &cc_main_1, t3...) // receiver: any, name: =~
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t2)) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"strict equal pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case === 5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	if value.Bool(value.StrictEqual(l0, (value.SmallInt(5)).ToValue())) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"strict not equal pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case !== 5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	if !(value.Bool(value.StrictEqual(l0, (value.SmallInt(5)).ToValue()))) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"strict not equal negative pattern": {
			input: `
	  var a: any = 0
		b := switch a
		case !== -5 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	if !(value.Bool(value.StrictEqual(l0, (value.SmallInt(-5)).ToValue()))) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"regex pattern": {
			input: `
	  a := "foo"
		b := switch a
		case %/fo+/ then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var regex0 = value.MustCompileRegex("fo+", bitfield.BitField8FromBitFlag(0))
var sym2 = value.ToSymbol("matches")
var fn_method0 vm.NativeFunction // Std::Regex.:matches

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.String // var a: Std::String
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	fn_method0 = vm.MethodToFunc((value.RegexClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.String("foo")
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = (regex0).ToValue()
	t3[1] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = fn_method0(thread, t3) // receiver: Std::Regex, name: matches
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"variable pattern": {
			input: `
	  a := 0
		b := switch a
		case n then n + 2
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
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
	var l1 value.Value // var b: Std::Int | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var l2 value.Value // var n: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	l2 = l0
	if value.True {
		t1 = value.AddInts(l2, (value.SmallInt(2)).ToValue())
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"range": {
			input: `
	  var a: any = 0
		b := switch a
		case -2...9 then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var range0 = value.NewClosedRange((value.SmallInt(-2)).ToValue(), (value.SmallInt(9)).ToValue())
var cc_main_1 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = (range0).ToValue()
	t3[1] = l0
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(symbol.S_contains, &cc_main_1, t3...) // receiver: Std::Int, name: #contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"range with constants": {
			input: `
		const Foo = 3
		const Bar = 10
	  var a: any = 0
		b := switch a
		case ::Foo...-::Bar then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("Foo")
var const0 value.SmallInt // constant: Foo, loc: <main>:2:3
var sym1 = value.ToSymbol("Bar")
var const1 value.SmallInt // constant: Bar, loc: <main>:3:3
var sym2 = value.ToSymbol("main")
var sym3 = value.ToSymbol("<main>")
var cc_main_1 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var namespace value.Value
	_ = namespace
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: any
	_ = l0
	var l1 value.Value // var b: Std::String | nil
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

	namespace = (value.RootModule).ToValue()
	value.AddConstant(namespace, sym0, (value.SmallInt(3)).ToValue())
	const0 = value.SmallInt(3)

	namespace = (value.RootModule).ToValue()
	value.AddConstant(namespace, sym1, (value.SmallInt(10)).ToValue())
	const1 = value.SmallInt(10)
	callFrame = thread.AddNativeCallFrame(sym2, sym3, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = (value.NewClosedRange((const0).ToValue(), (const1).NegateVal())).ToValue()
	t3[1] = l0
	callFrame.SetNativeLineNumber(6)
	t2, err = thread.CallMethodByNameWithCache(symbol.S_contains, &cc_main_1, t3...) // receiver: Std::Int, name: #contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"set pattern": {
			input: `
	  a := ^[1, 5, -4]
		b := switch a
		case ^[1, _, -4] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::HashSet.:length
var sym3 = value.ToSymbol("contains")
var fn_method1 vm.NativeFunction // Std::HashSet.:contains

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.HashSet // var a: Std::HashSet[Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym2))
	fn_method1 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, (value.SmallInt(-4)).ToValue(), (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue())
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.SetMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::HashSet[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(3)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (value.SmallInt(1)).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::Int], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (value.SmallInt(-4)).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::Int], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"set pattern with rest elements": {
			input: `
	  a := ^[1, 5, -4]
		b := switch a
		case ^[1, *, -4] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::HashSet.:length
var sym3 = value.ToSymbol("contains")
var fn_method1 vm.NativeFunction // Std::HashSet.:contains

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.HashSet // var a: Std::HashSet[Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym2))
	fn_method1 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, (value.SmallInt(-4)).ToValue(), (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue())
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.SetMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::HashSet[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.GreaterThanEqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (value.SmallInt(1)).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::Int], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (value.SmallInt(-4)).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::Int], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"word set pattern": {
			input: `
	  a := ^['foo', 'bar']
		b := switch a
		case ^w[foo bar] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::HashSet.:length
var sym3 = value.ToSymbol("contains")
var fn_method1 vm.NativeFunction // Std::HashSet.:contains

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.HashSet // var a: Std::HashSet[Std::String]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym2))
	fn_method1 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.NewNativeHashSetWithElements[value.String](value.String("bar"), value.String("foo"))
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.SetMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::HashSet[Std::String], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (value.String("foo")).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::String], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (value.String("bar")).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::String], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"symbol set pattern": {
			input: `
	  a := ^[:foo, :bar]
		b := switch a
		case ^s[foo bar] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("bar")
var sym3 = value.ToSymbol("foo")
var sym4 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::HashSet.:length
var sym5 = value.ToSymbol("contains")
var fn_method1 vm.NativeFunction // Std::HashSet.:contains

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.HashSet // var a: Std::HashSet[Std::Symbol]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym4))
	fn_method1 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym5))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.NewNativeHashSetWithElements[value.Symbol](sym2, sym3)
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.SetMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::HashSet[Std::Symbol], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (sym3).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::Symbol], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (sym2).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::Symbol], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"hex set pattern": {
			input: `
	  a := ^[0xff, 0x26]
		b := switch a
		case ^x[ff 26] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::HashSet.:length
var sym3 = value.ToSymbol("contains")
var fn_method1 vm.NativeFunction // Std::HashSet.:contains

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.HashSet // var a: Std::HashSet[Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym2))
	fn_method1 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, (value.SmallInt(255)).ToValue(), (value.SmallInt(38)).ToValue())
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.SetMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::HashSet[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (value.SmallInt(255)).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::Int], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (value.SmallInt(38)).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::Int], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"bin set pattern": {
			input: `
	  a := ^[0b11, 0b10]
		b := switch a
		case ^b[11 10] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::HashSet.:length
var sym3 = value.ToSymbol("contains")
var fn_method1 vm.NativeFunction // Std::HashSet.:contains

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.HashSet // var a: Std::HashSet[Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym2))
	fn_method1 = vm.MethodToFunc((value.HashSetClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue())
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.SetMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::HashSet[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (value.SmallInt(3)).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::Int], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = (l0).ToValue()
	t4[1] = (value.SmallInt(2)).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::HashSet[Std::Int], name: contains
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t3)) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"list pattern": {
			input: `
	  a := [1, 5, [8, 3]]
		b := switch a
		case [1, < 8, [a, > 1 && < 5]] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayList.:length
var cc_main_1 = &value.CallCache{}
var cc_main_2 = &value.CallCache{}
var cc_main_3 = &value.CallCache{}
var cc_main_4 = &value.CallCache{}
var cc_main_5 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int | Std::ArrayList[Std::Int]]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var t6 value.Value
	_ = t6
	var t7 value.Bool
	_ = t7
	var l2 value.Value // var a: Std::Int
	_ = l2
	var t8 value.Bool
	_ = t8
	var t9 value.Bool
	_ = t9
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayListClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.SmallInt(8)).ToValue(), (value.SmallInt(3)).ToValue())).ToValue())
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.ListMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayList[Std::Int | Std::ArrayList[Std::Int]], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(3)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(0)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t5
	t4[1] = (value.SmallInt(1)).ToValue()
	t6, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t4...) // receiver: Std::Int | Std::ArrayList[Std::Int], name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t6)) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(8)).ToValue(), (t5).Class())) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = t5
		t4[1] = (value.SmallInt(8)).ToValue()
		t5, err = thread.CallMethodByNameWithCache(symbol.OpLessThan, &cc_main_2, t4...) // receiver: Std::Int | Std::ArrayList[Std::Int], name: <
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t7 = value.ToBool(t5)
	} else {
		t7 = value.False
	}
	if !(t7) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(2)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t7 = value.True
	if !(value.Bool(value.IsA(t5, value.ListMixin))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = t5
	t5, err = thread.CallMethodByNameWithCache(symbol.L_length, &cc_main_3, t4...) // receiver: Std::Int | Std::ArrayList[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t5, (value.SmallInt(2)).ToValue()))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t5
	t4[1] = (value.SmallInt(0)).ToValue()
	t6, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_4, t4...) // receiver: Std::Int | Std::ArrayList[Std::Int], name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t6
	if !(value.True) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t5
	t4[1] = (value.SmallInt(1)).ToValue()
	t5, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_5, t4...) // receiver: Std::Int | Std::ArrayList[Std::Int], name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(1)).ToValue(), (t5).Class())) {
		t9 = value.Bool(value.GreaterThanInts(t5, (value.SmallInt(1)).ToValue()))
	} else {
		t9 = value.False
	}
	t8 = t9
	if t8 {
		if value.Bool(value.IsA((value.SmallInt(5)).ToValue(), (t5).Class())) {
			t9 = value.Bool(value.LessThanInts(t5, (value.SmallInt(5)).ToValue()))
		} else {
			t9 = value.False
		}
		t8 = t9
	}
	if !(t8) {
		t7 = value.False
		goto lbl3
	}
lbl3:
	if !(t7) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"word list pattern": {
			input: `
	  a := ['foo', 'bar']
		b := switch a
		case \w[foo bar] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayList.:length

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayList // var a: Std::ArrayList[Std::String]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayListClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.String](0, value.String("foo"), value.String("bar"))
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.ListMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayList[Std::String], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(0)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(((t5).AsString()) == (value.String("foo")))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(((t5).AsString()) == (value.String("bar")))) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"symbol list pattern": {
			input: `
	  a := [:foo, :bar]
		b := switch a
		case \s[foo bar] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("foo")
var sym3 = value.ToSymbol("bar")
var sym4 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayList.:length
var sym5 = value.ToSymbol("==")
var fn_method1 vm.NativeFunction // Std::Symbol.:==

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Symbol]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var t6 value.Value
	_ = t6
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayListClass).LookupMethod(sym4))
	fn_method1 = vm.MethodToFunc((value.SymbolClass).LookupMethod(sym5))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewNativeArrayListWithElements[value.Symbol](0, sym2, sym3)
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.ListMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayList[Std::Symbol], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(0)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t5
	t4[1] = (sym2).ToValue()
	t6, err = fn_method1(thread, t4) // receiver: Std::Symbol, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t6)) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t5
	t4[1] = (sym3).ToValue()
	t6, err = fn_method1(thread, t4) // receiver: Std::Symbol, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t6)) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"hex list pattern": {
			input: `
	  a := [0xff, 0x26]
		b := switch a
		case \x[ff 26] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayList.:length

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayListClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(255)).ToValue(), (value.SmallInt(38)).ToValue())
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.ListMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayList[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(0)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t5, (value.SmallInt(255)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t5, (value.SmallInt(38)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"bin list pattern": {
			input: `
	  a := [0b11, 0b10]
		b := switch a
		case \b[11 10] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayList.:length

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayListClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(3)).ToValue(), (value.SmallInt(2)).ToValue())
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.ListMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayList[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(0)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t5, (value.SmallInt(3)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t5, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"list pattern with rest elements": {
			input: `
	  a := [1, 5, [-2, 8, 3, 6]]
		b := switch a
		case [*b, [< 0, *c]] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayList.:length
var cc_main_1 = &value.CallCache{}
var cc_main_2 = &value.CallCache{}
var cc_main_3 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int | Std::ArrayList[Std::Int]]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var l2 *value.ArrayListOfValue // var b: Std::ArrayList[Std::Int | Std::ArrayList[Std::Int]]
	_ = l2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var t6 value.Value
	_ = t6
	var t7 value.Bool
	_ = t7
	var l3 *value.ArrayListOfValue // var c: Std::ArrayList[Std::Int]
	_ = l3
	var t8 value.Bool
	_ = t8
	var t9 value.Value
	_ = t9
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayListClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.SmallInt(-2)).ToValue(), (value.SmallInt(8)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(6)).ToValue())).ToValue())
	t2 = value.True
	l2 = value.NewArrayListOfValueWithElements(0)
	if !(value.Bool(value.IsA((l0).ToValue(), value.ListMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayList[Std::Int | Std::ArrayList[Std::Int]], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.GreaterThanEqualInts(t3, (value.SmallInt(1)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t3 = value.SubtractInts(t3, (value.SmallInt(1)).ToValue())
	t5 = (value.SmallInt(0)).ToValue()
	for value.Bool(value.LessThanInts(t5, t3)) {
		t6, err = (l0).SubscriptInt((t5).AsInt())
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		l2.Append(t6)
		t5 = value.IncrementInt(t5)
	}
	t6, err = (l0).SubscriptInt((t5).AsInt())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t7 = value.True
	l3 = value.NewArrayListOfValueWithElements(0)
	if !(value.Bool(value.IsA(t6, value.ListMixin))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = t6
	t5, err = thread.CallMethodByNameWithCache(symbol.L_length, &cc_main_1, t4...) // receiver: Std::Int | Std::ArrayList[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.GreaterThanEqualInts(t5, (value.SmallInt(1)).ToValue()))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t6
	t4[1] = (value.SmallInt(0)).ToValue()
	t6, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_2, t4...) // receiver: Std::Int | Std::ArrayList[Std::Int], name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(0)).ToValue(), (t6).Class())) {
		t8 = value.Bool(value.LessThanInts(t6, (value.SmallInt(0)).ToValue()))
	} else {
		t8 = value.False
	}
	if !(t8) {
		t7 = value.False
		goto lbl3
	}
	t6 = (value.SmallInt(1)).ToValue()
	for value.Bool(value.LessThanInts(t6, t5)) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = t6
		t4[1] = t6
		t9, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_3, t4...) // receiver: Std::Int | Std::ArrayList[Std::Int], name: []
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		l3.Append(t9)
		t6 = value.IncrementInt(t6)
	}
lbl3:
	if !(t7) {
		t2 = value.False
		goto lbl2
	}
	t5 = value.IncrementInt(t5)
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"list pattern with unnamed rest elements": {
			input: `
	  a := [1, 5, [-2, 8, 3, 6]]
		b := switch a
		case [*, [< 0, *]] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayList.:length
var cc_main_1 = &value.CallCache{}
var cc_main_2 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int | Std::ArrayList[Std::Int]]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var t6 value.Value
	_ = t6
	var t7 value.Bool
	_ = t7
	var t8 value.Bool
	_ = t8
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayListClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.SmallInt(-2)).ToValue(), (value.SmallInt(8)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(6)).ToValue())).ToValue())
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.ListMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayList[Std::Int | Std::ArrayList[Std::Int]], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.GreaterThanEqualInts(t3, (value.SmallInt(1)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5 = value.SubtractInts(t3, (value.SmallInt(1)).ToValue())
	t6, err = (l0).SubscriptInt((t5).AsInt())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t7 = value.True
	if !(value.Bool(value.IsA(t6, value.ListMixin))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = t6
	t5, err = thread.CallMethodByNameWithCache(symbol.L_length, &cc_main_1, t4...) // receiver: Std::Int | Std::ArrayList[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.GreaterThanEqualInts(t5, (value.SmallInt(1)).ToValue()))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t6
	t4[1] = (value.SmallInt(0)).ToValue()
	t6, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_2, t4...) // receiver: Std::Int | Std::ArrayList[Std::Int], name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(0)).ToValue(), (t6).Class())) {
		t8 = value.Bool(value.LessThanInts(t6, (value.SmallInt(0)).ToValue()))
	} else {
		t8 = value.False
	}
	if !(t8) {
		t7 = value.False
		goto lbl3
	}
	t6 = t5
lbl3:
	if !(t7) {
		t2 = value.False
		goto lbl2
	}
	t5 = value.IncrementInt(t5)
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"tuple pattern": {
			input: `
	  a := %[1, 5, %[8, 3]]
		b := switch a
		case %[1, < 8, %[a, > 1 && < 5]] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(8)).ToValue(), (value.SmallInt(3)).ToValue())
var arrtuple1 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue(), (arrtuple0).ToValue())
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayTuple.:length
var cc_main_1 = &value.CallCache{}
var cc_main_2 = &value.CallCache{}
var cc_main_3 = &value.CallCache{}
var cc_main_4 = &value.CallCache{}
var cc_main_5 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::Int]]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var t6 value.Value
	_ = t6
	var t7 value.Bool
	_ = t7
	var l2 value.Value // var a: Std::Int
	_ = l2
	var t8 value.Bool
	_ = t8
	var t9 value.Bool
	_ = t9
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayTupleClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple1
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.TupleMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::Int]], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(3)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(0)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t5
	t4[1] = (value.SmallInt(1)).ToValue()
	t6, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t4...) // receiver: Std::Int | Std::ArrayTuple[Std::Int], name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t6)) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(8)).ToValue(), (t5).Class())) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = t5
		t4[1] = (value.SmallInt(8)).ToValue()
		t5, err = thread.CallMethodByNameWithCache(symbol.OpLessThan, &cc_main_2, t4...) // receiver: Std::Int | Std::ArrayTuple[Std::Int], name: <
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t7 = value.ToBool(t5)
	} else {
		t7 = value.False
	}
	if !(t7) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(2)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t7 = value.True
	if !(value.Bool(value.IsA(t5, value.TupleMixin))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = t5
	t5, err = thread.CallMethodByNameWithCache(symbol.L_length, &cc_main_3, t4...) // receiver: Std::Int | Std::ArrayTuple[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t5, (value.SmallInt(2)).ToValue()))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t5
	t4[1] = (value.SmallInt(0)).ToValue()
	t6, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_4, t4...) // receiver: Std::Int | Std::ArrayTuple[Std::Int], name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t6
	if !(value.True) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t5
	t4[1] = (value.SmallInt(1)).ToValue()
	t5, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_5, t4...) // receiver: Std::Int | Std::ArrayTuple[Std::Int], name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(1)).ToValue(), (t5).Class())) {
		t9 = value.Bool(value.GreaterThanInts(t5, (value.SmallInt(1)).ToValue()))
	} else {
		t9 = value.False
	}
	t8 = t9
	if t8 {
		if value.Bool(value.IsA((value.SmallInt(5)).ToValue(), (t5).Class())) {
			t9 = value.Bool(value.LessThanInts(t5, (value.SmallInt(5)).ToValue()))
		} else {
			t9 = value.False
		}
		t8 = t9
	}
	if !(t8) {
		t7 = value.False
		goto lbl3
	}
lbl3:
	if !(t7) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"word tuple pattern": {
			input: `
	  a := %['foo', 'bar']
		b := switch a
		case %w[foo bar] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.String](0, value.String("foo"), value.String("bar"))
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayTuple.:length

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::String]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayTupleClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.TupleMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayTuple[Std::String], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(0)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(((t5).AsString()) == (value.String("foo")))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(((t5).AsString()) == (value.String("bar")))) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"symbol tuple pattern": {
			input: `
	  a := %[:foo, :bar]
		b := switch a
		case %s[foo bar] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("foo")
var sym3 = value.ToSymbol("bar")
var arrtuple0 = value.NewNativeArrayTupleWithElements[value.Symbol](0, sym2, sym3)
var sym4 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayTuple.:length
var sym5 = value.ToSymbol("==")
var fn_method1 vm.NativeFunction // Std::Symbol.:==

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Symbol]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var t6 value.Value
	_ = t6
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayTupleClass).LookupMethod(sym4))
	fn_method1 = vm.MethodToFunc((value.SymbolClass).LookupMethod(sym5))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.TupleMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayTuple[Std::Symbol], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(0)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t5
	t4[1] = (sym2).ToValue()
	t6, err = fn_method1(thread, t4) // receiver: Std::Symbol, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t6)) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t5
	t4[1] = (sym3).ToValue()
	t6, err = fn_method1(thread, t4) // receiver: Std::Symbol, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.ToBool(t6)) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"hex tuple pattern": {
			input: `
	  a := %[0xff, 0x26]
		b := switch a
		case %x[ff 26] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(255)).ToValue(), (value.SmallInt(38)).ToValue())
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayTuple.:length

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayTupleClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.TupleMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayTuple[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(0)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t5, (value.SmallInt(255)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t5, (value.SmallInt(38)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"bin tuple pattern": {
			input: `
	  a := %[0b11, 0b10]
		b := switch a
		case %b[11 10] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(3)).ToValue(), (value.SmallInt(2)).ToValue())
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayTuple.:length

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayTupleClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple0
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.TupleMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayTuple[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(0)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t5, (value.SmallInt(3)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5, err = (l0).SubscriptInt(int(value.SmallInt(1)))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t5, (value.SmallInt(2)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"tuple pattern with rest elements": {
			input: `
	  a := %[1, 5, %[-2, 8, 3, 6]]
		b := switch a
		case %[*b, %[< 0, *c]] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(-2)).ToValue(), (value.SmallInt(8)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(6)).ToValue())
var arrtuple1 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue(), (arrtuple0).ToValue())
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayTuple.:length
var cc_main_1 = &value.CallCache{}
var cc_main_2 = &value.CallCache{}
var cc_main_3 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::Int]]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var l2 *value.ArrayListOfValue // var b: Std::ArrayList[Std::Int | Std::ArrayTuple[Std::Int]]
	_ = l2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var t6 value.Value
	_ = t6
	var t7 value.Bool
	_ = t7
	var l3 *value.ArrayListOfValue // var c: Std::ArrayList[Std::Int]
	_ = l3
	var t8 value.Bool
	_ = t8
	var t9 value.Value
	_ = t9
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayTupleClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple1
	t2 = value.True
	l2 = value.NewArrayListOfValueWithElements(0)
	if !(value.Bool(value.IsA((l0).ToValue(), value.TupleMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::Int]], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.GreaterThanEqualInts(t3, (value.SmallInt(1)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t3 = value.SubtractInts(t3, (value.SmallInt(1)).ToValue())
	t5 = (value.SmallInt(0)).ToValue()
	for value.Bool(value.LessThanInts(t5, t3)) {
		t6, err = (l0).SubscriptInt((t5).AsInt())
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		l2.Append(t6)
		t5 = value.IncrementInt(t5)
	}
	t6, err = (l0).SubscriptInt((t5).AsInt())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t7 = value.True
	l3 = value.NewArrayListOfValueWithElements(0)
	if !(value.Bool(value.IsA(t6, value.TupleMixin))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = t6
	t5, err = thread.CallMethodByNameWithCache(symbol.L_length, &cc_main_1, t4...) // receiver: Std::Int | Std::ArrayTuple[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.GreaterThanEqualInts(t5, (value.SmallInt(1)).ToValue()))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t6
	t4[1] = (value.SmallInt(0)).ToValue()
	t6, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_2, t4...) // receiver: Std::Int | Std::ArrayTuple[Std::Int], name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(0)).ToValue(), (t6).Class())) {
		t8 = value.Bool(value.LessThanInts(t6, (value.SmallInt(0)).ToValue()))
	} else {
		t8 = value.False
	}
	if !(t8) {
		t7 = value.False
		goto lbl3
	}
	t6 = (value.SmallInt(1)).ToValue()
	for value.Bool(value.LessThanInts(t6, t5)) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = t6
		t4[1] = t6
		t9, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_3, t4...) // receiver: Std::Int | Std::ArrayTuple[Std::Int], name: []
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		l3.Append(t9)
		t6 = value.IncrementInt(t6)
	}
lbl3:
	if !(t7) {
		t2 = value.False
		goto lbl2
	}
	t5 = value.IncrementInt(t5)
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"tuple pattern with unnamed rest elements": {
			input: `
	  a := %[1, 5, %[-2, 8, 3, 6]]
		b := switch a
		case %[*, %[< 0, *]] then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(-2)).ToValue(), (value.SmallInt(8)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(6)).ToValue())
var arrtuple1 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(5)).ToValue(), (arrtuple0).ToValue())
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayTuple.:length
var cc_main_1 = &value.CallCache{}
var cc_main_2 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayTuple // var a: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::Int]]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var t5 value.Value
	_ = t5
	var t6 value.Value
	_ = t6
	var t7 value.Bool
	_ = t7
	var t8 value.Bool
	_ = t8
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayTupleClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = arrtuple1
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.TupleMixin))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayTuple[Std::Int | Std::ArrayTuple[Std::Int]], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.GreaterThanEqualInts(t3, (value.SmallInt(1)).ToValue()))) {
		t2 = value.False
		goto lbl2
	}
	t5 = value.SubtractInts(t3, (value.SmallInt(1)).ToValue())
	t6, err = (l0).SubscriptInt((t5).AsInt())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t7 = value.True
	if !(value.Bool(value.IsA(t6, value.TupleMixin))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = t6
	t5, err = thread.CallMethodByNameWithCache(symbol.L_length, &cc_main_1, t4...) // receiver: Std::Int | Std::ArrayTuple[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.GreaterThanEqualInts(t5, (value.SmallInt(1)).ToValue()))) {
		t7 = value.False
		goto lbl3
	}
	t4 = value.ResizeNativeArgs(t4, 3)
	t4[0] = t6
	t4[1] = (value.SmallInt(0)).ToValue()
	t6, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_2, t4...) // receiver: Std::Int | Std::ArrayTuple[Std::Int], name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(0)).ToValue(), (t6).Class())) {
		t8 = value.Bool(value.LessThanInts(t6, (value.SmallInt(0)).ToValue()))
	} else {
		t8 = value.False
	}
	if !(t8) {
		t7 = value.False
		goto lbl3
	}
	t6 = t5
lbl3:
	if !(t7) {
		t2 = value.False
		goto lbl2
	}
	t5 = value.IncrementInt(t5)
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"map pattern": {
			input: `
	  a := { 1 => 2, foo: :bar, "baz" => { dupa: [8, 3] } }
		b := switch a
		case { 1 => < 8, foo, "baz" => { dupa: [a, > 1 && < 5] } } then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("dupa")
var sym3 = value.ToSymbol("foo")
var sym4 = value.ToSymbol("bar")
var cc_main_1 = &value.CallCache{}
var cc_main_2 = &value.CallCache{}
var cc_main_3 = &value.CallCache{}
var cc_main_4 = &value.CallCache{}
var cc_main_5 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.HashMap // var a: Std::HashMap[Std::Int | Std::Symbol | Std::String, Std::Int | Std::Symbol | Std::HashMap[Std::Symbol, Std::ArrayList[Std::Int]]]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var err value.Value
	_ = err
	var t4 value.Bool
	_ = t4
	var t5 []value.Value
	_ = t5
	var l2 value.Value // var foo: Std::HashMap[Std::Int | Std::Symbol | Std::String, Std::Int | Std::Symbol | Std::HashMap[Std::Symbol, Std::ArrayList[Std::Int]]]
	_ = l2
	var t6 value.Bool
	_ = t6
	var t7 value.Value
	_ = t7
	var l3 value.Value // var a: Std::Int
	_ = l3
	var t8 value.Bool
	_ = t8
	var t9 value.Bool
	_ = t9
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.MustNewHashMapOfValueWithCapacityAndElements(nil, 0, value.MakePairOfValue((value.String("baz")).ToValue(), (vm.MustNewHashMapOfValueWithCapacityAndElements(nil, 0, value.MakePairOfValue((sym2).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.SmallInt(8)).ToValue(), (value.SmallInt(3)).ToValue())).ToValue()))).ToValue()), value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue()), value.MakePairOfValue((sym3).ToValue(), (sym4).ToValue()))
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.MapMixin))) {
		t2 = value.False
		goto lbl2
	}
	callFrame.SetNativeLineNumber(4)
	t3, err = (l0).GetValNil(thread, (value.SmallInt(1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(8)).ToValue(), (t3).Class())) {
		t5 = value.ResizeNativeArgs(t5, 3)
		t5[0] = t3
		t5[1] = (value.SmallInt(8)).ToValue()
		t3, err = thread.CallMethodByNameWithCache(symbol.OpLessThan, &cc_main_1, t5...) // receiver: void, name: <
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t4 = value.ToBool(t3)
	} else {
		t4 = value.False
	}
	if !(t4) {
		t2 = value.False
		goto lbl2
	}
	t3, err = (l0).GetValNil(thread, (sym3).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t3
	t3, err = (l0).GetValNil(thread, (value.String("baz")).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t4 = value.True
	if !(value.Bool(value.IsA(t3, value.MapMixin))) {
		t4 = value.False
		goto lbl3
	}
	t5 = value.ResizeNativeArgs(t5, 3)
	t5[0] = t3
	t5[1] = (sym2).ToValue()
	t3, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_2, t5...) // receiver: void, name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t6 = value.True
	if !(value.Bool(value.IsA(t3, value.ListMixin))) {
		t6 = value.False
		goto lbl4
	}
	t5 = value.ResizeNativeArgs(t5, 2)
	t5[0] = t3
	t3, err = thread.CallMethodByNameWithCache(symbol.L_length, &cc_main_3, t5...) // receiver: void, name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t6 = value.False
		goto lbl4
	}
	t5 = value.ResizeNativeArgs(t5, 3)
	t5[0] = t3
	t5[1] = (value.SmallInt(0)).ToValue()
	t7, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_4, t5...) // receiver: void, name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l3 = t7
	if !(value.True) {
		t6 = value.False
		goto lbl4
	}
	t5 = value.ResizeNativeArgs(t5, 3)
	t5[0] = t3
	t5[1] = (value.SmallInt(1)).ToValue()
	t3, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_5, t5...) // receiver: void, name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(1)).ToValue(), (t3).Class())) {
		t9 = value.Bool(value.GreaterThanInts(t3, (value.SmallInt(1)).ToValue()))
	} else {
		t9 = value.False
	}
	t8 = t9
	if t8 {
		if value.Bool(value.IsA((value.SmallInt(5)).ToValue(), (t3).Class())) {
			t9 = value.Bool(value.LessThanInts(t3, (value.SmallInt(5)).ToValue()))
		} else {
			t9 = value.False
		}
		t8 = t9
	}
	if !(t8) {
		t6 = value.False
		goto lbl4
	}
lbl4:
	if !(t6) {
		t4 = value.False
		goto lbl3
	}
lbl3:
	if !(t4) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"record pattern": {
			input: `
	  a := %{ 1 => 2, foo: :bar, "baz" => %{ dupa: [8, 3] } }
		b := switch a
		case %{ 1 => < 8, foo, "baz" => %{ dupa: [a, > 1 && < 5] } } then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("dupa")
var sym3 = value.ToSymbol("foo")
var sym4 = value.ToSymbol("bar")
var cc_main_1 = &value.CallCache{}
var cc_main_2 = &value.CallCache{}
var cc_main_3 = &value.CallCache{}
var cc_main_4 = &value.CallCache{}
var cc_main_5 = &value.CallCache{}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 vm.HashRecord // var a: Std::HashRecord[Std::Int | Std::Symbol | Std::String, Std::Int | Std::Symbol | Std::HashRecord[Std::Symbol, Std::ArrayList[Std::Int]]]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var err value.Value
	_ = err
	var t4 value.Bool
	_ = t4
	var t5 []value.Value
	_ = t5
	var l2 value.Value // var foo: void
	_ = l2
	var t6 value.Bool
	_ = t6
	var t7 value.Value
	_ = t7
	var l3 value.Value // var a: Std::Int
	_ = l3
	var t8 value.Bool
	_ = t8
	var t9 value.Bool
	_ = t9
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = vm.MustNewHashRecordOfValueWithElements(nil, value.MakePairOfValue((value.String("baz")).ToValue(), (vm.MustNewHashRecordOfValueWithElements(nil, value.MakePairOfValue((sym2).ToValue(), (value.NewArrayListOfValueWithElements(0, (value.SmallInt(8)).ToValue(), (value.SmallInt(3)).ToValue())).ToValue()))).ToValue()), value.MakePairOfValue((value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue()), value.MakePairOfValue((sym3).ToValue(), (sym4).ToValue()))
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.RecordMixin))) {
		t2 = value.False
		goto lbl2
	}
	callFrame.SetNativeLineNumber(4)
	t3, err = (l0).GetValNil(thread, (value.SmallInt(1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(8)).ToValue(), (t3).Class())) {
		t5 = value.ResizeNativeArgs(t5, 3)
		t5[0] = t3
		t5[1] = (value.SmallInt(8)).ToValue()
		t3, err = thread.CallMethodByNameWithCache(symbol.OpLessThan, &cc_main_1, t5...) // receiver: void, name: <
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t4 = value.ToBool(t3)
	} else {
		t4 = value.False
	}
	if !(t4) {
		t2 = value.False
		goto lbl2
	}
	t3, err = (l0).GetValNil(thread, (sym3).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t3
	t3, err = (l0).GetValNil(thread, (value.String("baz")).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t4 = value.True
	if !(value.Bool(value.IsA(t3, value.RecordMixin))) {
		t4 = value.False
		goto lbl3
	}
	t5 = value.ResizeNativeArgs(t5, 3)
	t5[0] = t3
	t5[1] = (sym2).ToValue()
	t3, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_2, t5...) // receiver: void, name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t6 = value.True
	if !(value.Bool(value.IsA(t3, value.ListMixin))) {
		t6 = value.False
		goto lbl4
	}
	t5 = value.ResizeNativeArgs(t5, 2)
	t5[0] = t3
	t3, err = thread.CallMethodByNameWithCache(symbol.L_length, &cc_main_3, t5...) // receiver: void, name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if !(value.Bool(value.EqualInts(t3, (value.SmallInt(2)).ToValue()))) {
		t6 = value.False
		goto lbl4
	}
	t5 = value.ResizeNativeArgs(t5, 3)
	t5[0] = t3
	t5[1] = (value.SmallInt(0)).ToValue()
	t7, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_4, t5...) // receiver: void, name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l3 = t7
	if !(value.True) {
		t6 = value.False
		goto lbl4
	}
	t5 = value.ResizeNativeArgs(t5, 3)
	t5[0] = t3
	t5[1] = (value.SmallInt(1)).ToValue()
	t3, err = thread.CallMethodByNameWithCache(symbol.OpSubscript, &cc_main_5, t5...) // receiver: void, name: []
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.Bool(value.IsA((value.SmallInt(1)).ToValue(), (t3).Class())) {
		t9 = value.Bool(value.GreaterThanInts(t3, (value.SmallInt(1)).ToValue()))
	} else {
		t9 = value.False
	}
	t8 = t9
	if t8 {
		if value.Bool(value.IsA((value.SmallInt(5)).ToValue(), (t3).Class())) {
			t9 = value.Bool(value.LessThanInts(t3, (value.SmallInt(5)).ToValue()))
		} else {
			t9 = value.False
		}
		t8 = t9
	}
	if !(t8) {
		t6 = value.False
		goto lbl4
	}
lbl4:
	if !(t6) {
		t4 = value.False
		goto lbl3
	}
lbl3:
	if !(t4) {
		t2 = value.False
		goto lbl2
	}
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
	l1 = t1
}
`,
		},
		"object pattern": {
			input: `
	  a := [0b11, 0b10]
		b := switch a
		case ::Std::ArrayList(length: > 1 && < 5 as l, first) then "a"
		end
			`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayList.:length
var cc_main_1 = &value.CallCache{}
var cc_main_2 = &value.CallCache{}
var sym3 = value.ToSymbol("first")
var fn_method1 vm.NativeFunction // Std::Iterable::FiniteBase.:first

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int]
	_ = l0
	var l1 value.Value // var b: Std::String | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Bool
	_ = t2
	var t3 value.Value
	_ = t3
	var t4 []value.Value
	_ = t4
	var err value.Value
	_ = err
	var l2 value.Value // var l: void
	_ = l2
	var t5 value.Bool
	_ = t5
	var t6 value.Bool
	_ = t6
	var l3 value.Value // var first: Std::Int
	_ = l3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayListClass).LookupMethod(sym2))
	fn_method1 = vm.MethodToFunc((value.IterableFiniteBaseMixin).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(3)).ToValue(), (value.SmallInt(2)).ToValue())
	t2 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.ArrayListClass))) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(4)
	t3, err = fn_method0(thread, t4) // receiver: Std::ArrayList[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t3
	if value.Bool(value.IsA((value.SmallInt(1)).ToValue(), (t3).Class())) {
		t4 = value.ResizeNativeArgs(t4, 3)
		t4[0] = t3
		t4[1] = (value.SmallInt(1)).ToValue()
		t3, err = thread.CallMethodByNameWithCache(symbol.OpGreaterThan, &cc_main_1, t4...) // receiver: void, name: >
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t6 = value.ToBool(t3)
	} else {
		t6 = value.False
	}
	t5 = t6
	if t5 {
		if value.Bool(value.IsA((value.SmallInt(5)).ToValue(), (t3).Class())) {
			t4 = value.ResizeNativeArgs(t4, 3)
			t4[0] = t3
			t4[1] = (value.SmallInt(5)).ToValue()
			t3, err = thread.CallMethodByNameWithCache(symbol.OpLessThan, &cc_main_2, t4...) // receiver: void, name: <
			if err.IsNotUndefined() {
				thread.CaptureStackTrace()
				thread.Panic(err)
			}
			t6 = value.ToBool(t3)
		} else {
			t6 = value.False
		}
		t5 = t6
	}
	if !(t5) {
		t2 = value.False
		goto lbl2
	}
	t4 = value.ResizeNativeArgs(t4, 2)
	t4[0] = (l0).ToValue()
	t3, err = fn_method1(thread, t4) // receiver: Std::ArrayList[Std::Int], name: first
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l3 = t3
lbl2:
	if t2 {
		t1 = (value.String("a")).ToValue()
		goto lbl1
	}
	t1 = value.Nil
lbl1:
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

func TestGoMatch(t *testing.T) {
	tests := goTestTable{
		"contains a pattern": {
			input: `
			a := [0b11, 0b10]
			if a match ::Std::ArrayList(length: > 1 && < 5 as l, first)
				puts "hooray"
			end
		`,
			want: `package main

import (
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayList.:length
var cc_main_1 = &value.CallCache{}
var cc_main_2 = &value.CallCache{}
var sym3 = value.ToSymbol("first")
var fn_method1 vm.NativeFunction // Std::Iterable::FiniteBase.:first
var sym4 = value.ToSymbol("puts@1")
var fn_method2 vm.NativeFunction // Std::Kernel::puts@1

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.ArrayList // var a: Std::ArrayList[Std::Int]
	_ = l0
	var t1 value.Bool
	_ = t1
	var t2 value.Value
	_ = t2
	var t3 []value.Value
	_ = t3
	var err value.Value
	_ = err
	var l1 value.Value // var l: void
	_ = l1
	var t4 value.Bool
	_ = t4
	var t5 value.Bool
	_ = t5
	var l2 value.Value // var first: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayListClass).LookupMethod(sym2))
	fn_method1 = vm.MethodToFunc((value.IterableFiniteBaseMixin).LookupMethod(sym3))
	fn_method2 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym4))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.NewArrayListOfValueWithElements(0, (value.SmallInt(3)).ToValue(), (value.SmallInt(2)).ToValue())
	t1 = value.True
	if !(value.Bool(value.IsA((l0).ToValue(), value.ArrayListClass))) {
		t1 = value.False
		goto lbl1
	}
	t3 = value.ResizeNativeArgs(t3, 2)
	t3[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(3)
	t2, err = fn_method0(thread, t3) // receiver: Std::ArrayList[Std::Int], name: length
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t2
	if value.Bool(value.IsA((value.SmallInt(1)).ToValue(), (t2).Class())) {
		t3 = value.ResizeNativeArgs(t3, 3)
		t3[0] = t2
		t3[1] = (value.SmallInt(1)).ToValue()
		t2, err = thread.CallMethodByNameWithCache(symbol.OpGreaterThan, &cc_main_1, t3...) // receiver: void, name: >
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t5 = value.ToBool(t2)
	} else {
		t5 = value.False
	}
	t4 = t5
	if t4 {
		if value.Bool(value.IsA((value.SmallInt(5)).ToValue(), (t2).Class())) {
			t3 = value.ResizeNativeArgs(t3, 3)
			t3[0] = t2
			t3[1] = (value.SmallInt(5)).ToValue()
			t2, err = thread.CallMethodByNameWithCache(symbol.OpLessThan, &cc_main_2, t3...) // receiver: void, name: <
			if err.IsNotUndefined() {
				thread.CaptureStackTrace()
				thread.Panic(err)
			}
			t5 = value.ToBool(t2)
		} else {
			t5 = value.False
		}
		t4 = t5
	}
	if !(t4) {
		t1 = value.False
		goto lbl1
	}
	t3 = value.ResizeNativeArgs(t3, 2)
	t3[0] = (l0).ToValue()
	t2, err = fn_method1(thread, t3) // receiver: Std::ArrayList[Std::Int], name: first
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t2
lbl1:
	if t1 {
		t3 = value.ResizeNativeArgs(t3, 3)
		t3[0] = (value.KernelModule).ToValue()
		t3[1] = (value.String("hooray")).ToValue()
		callFrame.SetNativeLineNumber(4)
		_, err = fn_method2(thread, t3) // receiver: Std::Kernel, name: puts@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}
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
