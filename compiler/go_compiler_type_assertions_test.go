package compiler_test

import (
	"testing"
)

func TestGoAs(t *testing.T) {
	tests := goTestTable{
		"assert": {
			input: `
				var a: Int | Float = 1
				a as ::Std::Int
				nil
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
	var l0 value.Value // var a: Std::Int | Std::Float
	_ = l0
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(3)
	err = value.As(l0, value.IntClass)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"cast": {
			input: `
				var a: Int | Float = 1
				b := a as ::Std::Int
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
	var l0 value.Value // var a: Std::Int | Std::Float
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
	l0 = (value.SmallInt(1)).ToValue()
	t1 = l0
	callFrame.SetNativeLineNumber(3)
	err = value.As(t1, value.IntClass)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
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

func TestGoMust(t *testing.T) {
	tests := goTestTable{
		"assert": {
			input: `
				var a: Int? = nil
				must a
				nil
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Nil
	callFrame.SetNativeLineNumber(3)
	err = value.Must(l0)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
		"cast": {
			input: `
				var a: Int? = nil
				b := must a
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
	var l0 value.Value // var a: Std::Int?
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
	l0 = value.Nil
	t1 = l0
	callFrame.SetNativeLineNumber(3)
	err = value.Must(t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
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

func TestGoInstanceOf(t *testing.T) {
	tests := goTestTable{
		"optimised instance of": {
			input: `
				var a: Int? = 5
				result := a <<: Int
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.Bool // var result: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	l1 = value.Bool(value.InstanceOf(l0, value.IntClass))
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

func TestGoReverseInstanceOf(t *testing.T) {
	tests := goTestTable{
		"optimised reverse instance of": {
			input: `
				var a: Int? = 5
				result := Int :>> a
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.Bool // var result: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	l1 = value.Bool(value.InstanceOf(l0, value.IntClass))
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

func TestGoIsA(t *testing.T) {
	tests := goTestTable{
		"optimised is a class": {
			input: `
				var a: Int? = 5
				result := a <: Int
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.Bool // var result: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	l1 = value.Bool(value.IsA(l0, value.IntClass))
}
`,
		},
		"optimised is a mixin": {
			input: `
				var a: HashMap[String, String]? = { 'foo' => 'bar' }
				result := a <: Map
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
	var l0 value.Value // var a: Std::HashMap[Std::String, Std::String]?
	_ = l0
	var l1 value.Bool // var result: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (vm.NewNativeHashMapFromMap(map[value.String]value.String{value.String("foo"): value.String("bar")})).ToValue()
	l1 = value.Bool(value.IsA(l0, value.MapMixin))
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

func TestGoReverseIsA(t *testing.T) {
	tests := goTestTable{
		"optimised reverse is a class": {
			input: `
				var a: Int? = 5
				result := Int :> a
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.Bool // var result: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	l1 = value.Bool(value.IsA(l0, value.IntClass))
}
`,
		},
		"optimised reverse is a mixin": {
			input: `
				var a: HashMap[String, String]? = { 'foo' => 'bar' }
				result := Map :> a
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
	var l0 value.Value // var a: Std::HashMap[Std::String, Std::String]?
	_ = l0
	var l1 value.Bool // var result: Std::Bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (vm.NewNativeHashMapFromMap(map[value.String]value.String{value.String("foo"): value.String("bar")})).ToValue()
	l1 = value.Bool(value.IsA(l0, value.MapMixin))
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
