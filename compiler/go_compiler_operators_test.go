package compiler_test

import (
	"testing"
)

func TestGoComplexAssignmentLocals(t *testing.T) {
	tests := goTestTable{
		"increment int": {
			input: "a := 1; a++",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.IncrementInt(l0)
}
`,
		},

		"increment int64": {
			input: "a := 1i64; a++",
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
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(1)
	l0 = (l0) + 1
}
`,
		},

		"increment int32": {
			input: "a := 1i32; a++",
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
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(1)
	l0 = (l0) + 1
}
`,
		},

		"increment int16": {
			input: "a := 1i16; a++",
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
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(1)
	l0 = (l0) + 1
}
`,
		},

		"increment int8": {
			input: "a := 1i8; a++",
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
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(1)
	l0 = (l0) + 1
}
`,
		},

		"increment uint64": {
			input: "a := 1u64; a++",
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
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(1)
	l0 = (l0) + 1
}
`,
		},

		"increment uint32": {
			input: "a := 1u32; a++",
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
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(1)
	l0 = (l0) + 1
}
`,
		},

		"increment uint16": {
			input: "a := 1u16; a++",
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
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(1)
	l0 = (l0) + 1
}
`,
		},

		"increment uint8": {
			input: "a := 1u8; a++",
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
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(1)
	l0 = (l0) + 1
}
`,
		},

		"increment uint": {
			input: "a := 1u; a++",
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
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(1)
	l0 = (l0) + 1
}
`,
		},

		"increment union": {
			input: "var a: Int | Int64 = 1; a++",
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
	var l0 value.Value // var a: Std::Int | Std::Int64
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.IncrementVal(l0)
}
`,
		},

		"increment char": {
			input: "a := `a`; a++",
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
	var l0 value.Char // var a: Std::Char
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Char('a')
	l0 = (l0) + 1
}
`,
		},

		"decrement int": {
			input: "a := 1; a--",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.DecrementInt(l0)
}
`,
		},

		"decrement int64": {
			input: "a := 1i64; a--",
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
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(1)
	l0 = (l0) - 1
}
`,
		},

		"decrement int32": {
			input: "a := 1i32; a--",
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
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(1)
	l0 = (l0) - 1
}
`,
		},

		"decrement int16": {
			input: "a := 1i16; a--",
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
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(1)
	l0 = (l0) - 1
}
`,
		},

		"decrement int8": {
			input: "a := 1i8; a--",
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
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(1)
	l0 = (l0) - 1
}
`,
		},

		"decrement uint64": {
			input: "a := 1u64; a--",
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
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(1)
	l0 = (l0) - 1
}
`,
		},

		"decrement uint32": {
			input: "a := 1u32; a--",
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
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(1)
	l0 = (l0) - 1
}
`,
		},

		"decrement uint16": {
			input: "a := 1u16; a--",
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
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(1)
	l0 = (l0) - 1
}
`,
		},

		"decrement uint8": {
			input: "a := 1u8; a--",
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
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(1)
	l0 = (l0) - 1
}
`,
		},

		"decrement uint": {
			input: "a := 1u; a--",
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
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(1)
	l0 = (l0) - 1
}
`,
		},

		"decrement char": {
			input: "a := `a`; a--",
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
	var l0 value.Char // var a: Std::Char
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Char('a')
	l0 = (l0) - 1
}
`,
		},

		"decrement union": {
			input: "var a: Int | Int64 = 1; a--",
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
	var l0 value.Value // var a: Std::Int | Std::Int64
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.DecrementVal(l0)
}
`,
		},

		"add": {
			input: "a := 1; a += 3",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.AddInts(l0, (value.SmallInt(3)).ToValue())
}
`,
		},

		"subtract": {
			input: "a := 1; a -= 3",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.SubtractInts(l0, (value.SmallInt(3)).ToValue())
}
`,
		},

		"multiply": {
			input: "a := 1; a *= 3",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.MultiplyInts(l0, (value.SmallInt(3)).ToValue())
}
`,
		},

		"divide": {
			input: "a := 1; a /= 3",
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
	t1, err = value.DivideInts(l0, (value.SmallInt(3)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
}
`,
		},

		"exponentiate": {
			input: "a := 1; a **= 3",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.ExponentiateInts(l0, (value.SmallInt(3)).ToValue())
}
`,
		},

		"modulo": {
			input: "a := 1; a %= 3",
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
	t1, err = value.ModuloInts(l0, (value.SmallInt(3)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
}
`,
		},

		"bitwise AND": {
			input: "a := 1; a &= 3",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.BitwiseAndInts(l0, (value.SmallInt(3)).ToValue())
}
`,
		},

		"bitwise OR": {
			input: "a := 1; a |= 3",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.BitwiseOrInts(l0, (value.SmallInt(3)).ToValue())
}
`,
		},

		"bitwise XOR": {
			input: "a := 1; a ^= 3",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.BitwiseXorInts(l0, (value.SmallInt(3)).ToValue())
}
`,
		},

		"left bitshift": {
			input: "a := 1; a <<= 3",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.LeftBitshiftInts(l0, (value.SmallInt(3)).ToValue())
}
`,
		},

		"left logical bitshift": {
			input: "a := 1u64; a <<<= 3",
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
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var t1 value.UInt64
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(1)
	t1, err = value.StrictIntLeftBitshift(l0, (value.SmallInt(3)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
}
`,
		},

		"right bitshift": {
			input: "a := 1; a >>= 3",
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
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.RightBitshiftInts(l0, (value.SmallInt(3)).ToValue())
}
`,
		},

		"right logical bitshift": {
			input: "a := 1u64; a >>>= 3",
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
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var t1 value.UInt64
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(1)
	t1, err = value.StrictIntRightBitshift(l0, (value.SmallInt(3)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
}
`,
		},

		"logic OR": {
			input: "var a: Int? = 1; a ||= 3",
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
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
	t1 = l0
	if value.Falsy(t1) {
		t1 = (value.SmallInt(3)).ToValue()
	}
	l0 = t1
}
`,
		},

		"logic AND": {
			input: "var a: Int? = 1; a &&= 3",
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
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
	t1 = l0
	if value.Truthy(t1) {
		t1 = (value.SmallInt(3)).ToValue()
	}
	l0 = t1
}
`,
		},

		"nil coalesce": {
			input: "var a: Int? = 1; a ??= 3",
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
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
	t1 = l0
	if value.IsNil(t1) {
		t1 = (value.SmallInt(3)).ToValue()
	}
	l0 = t1
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

func TestGoComplexAssignmentInstanceVariables(t *testing.T) {
	tests := goTestTable{
		"increment int": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.IncrementInt(value.GetInstanceVariable(self, 0))
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment int64": {
			input: `
				class Foo
					var @a: Int64
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Int64) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Int64
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsInt64()) + 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsInt64())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment int32": {
			input: `
				class Foo
					var @a: Int32
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Int32) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Int32
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsInt32()) + 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsInt32())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment int16": {
			input: `
				class Foo
					var @a: Int16
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Int16) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Int16
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsInt16()) + 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsInt16())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment int8": {
			input: `
				class Foo
					var @a: Int8
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Int8) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Int8
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsInt8()) + 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsInt8())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment uint64": {
			input: `
				class Foo
					var @a: UInt64
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt64) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt64
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsUInt64()) + 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt64())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment uint32": {
			input: `
				class Foo
					var @a: UInt32
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt32) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt32
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsUInt32()) + 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt32())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment uint16": {
			input: `
				class Foo
					var @a: UInt16
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt16) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt16
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsUInt16()) + 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt16())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment uint8": {
			input: `
				class Foo
					var @a: UInt8
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt8) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt8
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsUInt8()) + 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt8())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment uint": {
			input: `
				class Foo
					var @a: UInt
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsUInt()) + 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment char": {
			input: `
				class Foo
					var @a: Char
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Char) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Char
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsChar()) + 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsChar())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"increment union": {
			input: `
				class Foo
					var @a: Int | Int64
					init(@a); end

					def foo then @a++
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.IncrementVal(value.GetInstanceVariable(self, 0))
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement int": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.DecrementInt(value.GetInstanceVariable(self, 0))
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement int64": {
			input: `
				class Foo
					var @a: Int64
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Int64) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Int64
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsInt64()) - 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsInt64())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement int32": {
			input: `
				class Foo
					var @a: Int32
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Int32) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Int32
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsInt32()) - 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsInt32())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement int16": {
			input: `
				class Foo
					var @a: Int16
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Int16) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Int16
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsInt16()) - 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsInt16())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement int8": {
			input: `
				class Foo
					var @a: Int8
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Int8) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Int8
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsInt8()) - 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsInt8())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement uint64": {
			input: `
				class Foo
					var @a: UInt64
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt64) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt64
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsUInt64()) - 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt64())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement uint32": {
			input: `
				class Foo
					var @a: UInt32
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt32) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt32
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsUInt32()) - 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt32())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement uint16": {
			input: `
				class Foo
					var @a: UInt16
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt16) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt16
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsUInt16()) - 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt16())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement uint8": {
			input: `
				class Foo
					var @a: UInt8
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt8) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt8
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsUInt8()) - 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt8())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement uint": {
			input: `
				class Foo
					var @a: UInt
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsUInt()) - 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement char": {
			input: `
				class Foo
					var @a: Char
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Char) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Char
	_ = t1

	t1 = ((value.GetInstanceVariable(self, 0)).AsChar()) - 1
	value.SetInstanceVariable(self, 0, (t1).ToValue())
	return (t1).ToValue(), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsChar())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"decrement union": {
			input: `
				class Foo
					var @a: Int | Int64
					init(@a); end

					def foo then @a--
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.DecrementVal(value.GetInstanceVariable(self, 0))
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"add": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a += 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.AddInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"subtract": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a -= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.SubtractInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"multiply": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a *= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.MultiplyInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"divide": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a /= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2

	t1, err = value.DivideInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	t2 = t1
	value.SetInstanceVariable(self, 0, t2)
	return t2, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"exponentiate": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a **= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.ExponentiateInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"modulo": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a %= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2

	t1, err = value.ModuloInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	t2 = t1
	value.SetInstanceVariable(self, 0, t2)
	return t2, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"bitwise AND": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a &= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.BitwiseAndInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"bitwise OR": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a |= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.BitwiseOrInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"bitwise XOR": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a ^= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.BitwiseXorInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"left bitshift": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a <<= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.LeftBitshiftInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"left logical bitshift": {
			input: `
				class Foo
					var @a: UInt64
					init(@a); end

					def foo then @a <<<= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt64) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.UInt64
	_ = t1
	var t2 value.Value
	_ = t2

	t1, err = value.StrictIntLeftBitshift((value.GetInstanceVariable(self, 0)).AsUInt64(), (value.SmallInt(3)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	t2 = (t1).ToValue()
	value.SetInstanceVariable(self, 0, t2)
	return t2, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt64())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"right bitshift": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a >>= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	t1 = value.RightBitshiftInts(value.GetInstanceVariable(self, 0), (value.SmallInt(3)).ToValue())
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"right logical bitshift": {
			input: `
				class Foo
					var @a: Int64
					init(@a); end

					def foo then @a >>>= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Int64) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, (l0).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:6:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Int64
	_ = t1
	var t2 value.Value
	_ = t2

	t1, err = value.StrictIntLogicalRightBitshift((value.GetInstanceVariable(self, 0)).AsInt64(), (value.SmallInt(3)).ToValue(), value.LogicalRightShift64)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	t2 = (t1).ToValue()
	value.SetInstanceVariable(self, 0, t2)
	return t2, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsInt64())
		return result, err
	}, vm.DefWithParameters(1))
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},

		"logic OR": {
			input: `
				class Foo
					var @a: Int?

					def foo then @a ||= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2

	t1 = value.GetInstanceVariable(self, 0)
	if value.Falsy(t1) {
		t1 = (value.SmallInt(3)).ToValue()
	}
	t2 = t1
	value.SetInstanceVariable(self, 0, t2)
	return t2, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"nil coalesce": {
			input: `
				class Foo
					var @a: Int?

					def foo then @a ??= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2

	t1 = value.GetInstanceVariable(self, 0)
	if value.IsNil(t1) {
		t1 = (value.SmallInt(3)).ToValue()
	}
	t2 = t1
	value.SetInstanceVariable(self, 0, t2)
	return t2, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"logic AND": {
			input: `
				class Foo
					var @a: Int?

					def foo then @a &&= 3
				end
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

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2

	t1 = value.GetInstanceVariable(self, 0)
	if value.Truthy(t1) {
		t1 = (value.SmallInt(3)).ToValue()
	}
	t2 = t1
	value.SetInstanceVariable(self, 0, t2)
	return t2, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
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
