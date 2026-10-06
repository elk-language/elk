package compiler_test

import (
	"testing"
)

func TestGoQuote(t *testing.T) {
	tests := goTestTable{
		"without unquote": {
			input: `
				quote
					1 + 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
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
	var t1 ast.Node
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t1 = ast.NewBinaryExpressionNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(20, 3, 10))), token.New(position.NewLocation("<main>", position.NewSpan(position.New(18, 3, 8), position.New(18, 3, 8))), token.PLUS), ast.NewIntLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(16, 3, 6))), "1"), ast.NewIntLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(20, 3, 10), position.New(20, 3, 10))), "2"))
}
`,
		},

		"unquote expression": {
			input: `
				quote
					1 + unquote(5)
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
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
var sym2 = value.ToSymbol("to_ast_expr_node")
var fn_method0 vm.NativeFunction // Std::Int.:to_ast_expr_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.IntClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Int, name: to_ast_expr_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewBinaryExpressionNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(29, 3, 19))), token.New(position.NewLocation("<main>", position.NewSpan(position.New(18, 3, 8), position.New(18, 3, 8))), token.PLUS), ast.NewIntLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(16, 3, 6))), "1"), (t1).AsReference().(ast.ExpressionNode))
}
`,
		},
		"short unquote expression": {
			input: `
				quote
					1 + !{5}
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
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
var sym2 = value.ToSymbol("to_ast_expr_node")
var fn_method0 vm.NativeFunction // Std::Int.:to_ast_expr_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.IntClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Int, name: to_ast_expr_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewBinaryExpressionNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(23, 3, 13))), token.New(position.NewLocation("<main>", position.NewSpan(position.New(18, 3, 8), position.New(18, 3, 8))), token.PLUS), ast.NewIntLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(16, 3, 6))), "1"), (t1).AsReference().(ast.ExpressionNode))
}
`,
		},

		"unquote identifier": {
			input: `
				quote
					var unquote(:foo): String
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("foo")
var sym3 = value.ToSymbol("to_ast_ident_node")
var fn_method0 vm.NativeFunction // Std::Symbol.:to_ast_ident_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.SymbolClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (sym2).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Symbol, name: to_ast_ident_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewVariableDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(40, 3, 30))), "", (t1).AsReference().(ast.IdentifierNode), ast.NewPublicConstantNode(position.NewLocation("<main>", position.NewSpan(position.New(35, 3, 25), position.New(40, 3, 30))), "String"), nil)
}
`,
		},
		"short unquote identifier": {
			input: `
				quote
					var !{:foo}: String
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("foo")
var sym3 = value.ToSymbol("to_ast_ident_node")
var fn_method0 vm.NativeFunction // Std::Symbol.:to_ast_ident_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.SymbolClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (sym2).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Symbol, name: to_ast_ident_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewVariableDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(34, 3, 24))), "", (t1).AsReference().(ast.IdentifierNode), ast.NewPublicConstantNode(position.NewLocation("<main>", position.NewSpan(position.New(29, 3, 19), position.New(34, 3, 24))), "String"), nil)
}
`,
		},
		"unquote_ident": {
			input: `
				quote
					var unquote_ident(:foo): String
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("foo")
var sym3 = value.ToSymbol("to_ast_ident_node")
var fn_method0 vm.NativeFunction // Std::Symbol.:to_ast_ident_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.SymbolClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (sym2).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Symbol, name: to_ast_ident_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewVariableDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(46, 3, 36))), "", (t1).AsReference().(ast.IdentifierNode), ast.NewPublicConstantNode(position.NewLocation("<main>", position.NewSpan(position.New(41, 3, 31), position.New(46, 3, 36))), "String"), nil)
}
`,
		},

		"unquote pattern expression": {
			input: `
				quote
					var ^[unquote(1)] = "foo"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("to_ast_pattern_expr_node")
var fn_method0 vm.NativeFunction // Std::Int.:to_ast_pattern_expr_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.IntClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Int, name: to_ast_pattern_expr_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewVariablePatternDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(40, 3, 30))), ast.NewSetPatternNode(position.NewLocation("<main>", position.NewSpan(position.New(20, 3, 10), position.New(32, 3, 22))), []ast.PatternNode{(t1).AsReference().(ast.LiteralPatternNode)}), ast.NewDoubleQuotedStringLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(36, 3, 26), position.New(40, 3, 30))), "foo"))
}
`,
		},
		"short unquote pattern expression": {
			input: `
				quote
					var ^[!{1}] = "foo"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("to_ast_pattern_expr_node")
var fn_method0 vm.NativeFunction // Std::Int.:to_ast_pattern_expr_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.IntClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Int, name: to_ast_pattern_expr_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewVariablePatternDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(34, 3, 24))), ast.NewSetPatternNode(position.NewLocation("<main>", position.NewSpan(position.New(20, 3, 10), position.New(26, 3, 16))), []ast.PatternNode{(t1).AsReference().(ast.LiteralPatternNode)}), ast.NewDoubleQuotedStringLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(30, 3, 20), position.New(34, 3, 24))), "foo"))
}
`,
		},

		"unquote pattern": {
			input: `
				quote
					var [unquote(Elk::AST::ListPatternNode())] = "foo"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("Std::Elk::AST::ListPatternNode")
var const0 *value.Class // Std::Elk::AST::ListPatternNode
var sym3 = value.ToSymbol("#init")
var fn_method0 vm.NativeFunction // Std::Elk::AST::ListPatternNode.:#init

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	const0 = (*value.Class)((value.GetConstant(sym2)).Pointer())

	fn_method0 = vm.MethodToFunc((const0).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 4)
	t2[0] = const0.CreateInstance()
	t2[1] = value.Undefined
	t2[2] = value.Undefined
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Elk::AST::ListPatternNode, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewVariablePatternDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(65, 3, 55))), ast.NewListPatternNode(position.NewLocation("<main>", position.NewSpan(position.New(20, 3, 10), position.New(57, 3, 47))), []ast.PatternNode{(t1).AsReference().(ast.PatternNode)}), ast.NewDoubleQuotedStringLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(61, 3, 51), position.New(65, 3, 55))), "foo"))
}
`,
		},
		"unquote_pattern": {
			input: `
				quote
					var [unquote_pattern(Elk::AST::ListPatternNode())] = "foo"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("Std::Elk::AST::ListPatternNode")
var const0 *value.Class // Std::Elk::AST::ListPatternNode
var sym3 = value.ToSymbol("#init")
var fn_method0 vm.NativeFunction // Std::Elk::AST::ListPatternNode.:#init

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	const0 = (*value.Class)((value.GetConstant(sym2)).Pointer())

	fn_method0 = vm.MethodToFunc((const0).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 4)
	t2[0] = const0.CreateInstance()
	t2[1] = value.Undefined
	t2[2] = value.Undefined
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Elk::AST::ListPatternNode, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewVariablePatternDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(73, 3, 63))), ast.NewListPatternNode(position.NewLocation("<main>", position.NewSpan(position.New(20, 3, 10), position.New(65, 3, 55))), []ast.PatternNode{(t1).AsReference().(ast.PatternNode)}), ast.NewDoubleQuotedStringLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(69, 3, 59), position.New(73, 3, 63))), "foo"))
}
`,
		},
		"short unquote pattern": {
			input: `
				quote
					var [!{Elk::AST::ListPatternNode()}] = "foo"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("Std::Elk::AST::ListPatternNode")
var const0 *value.Class // Std::Elk::AST::ListPatternNode
var sym3 = value.ToSymbol("#init")
var fn_method0 vm.NativeFunction // Std::Elk::AST::ListPatternNode.:#init

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	const0 = (*value.Class)((value.GetConstant(sym2)).Pointer())

	fn_method0 = vm.MethodToFunc((const0).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 4)
	t2[0] = const0.CreateInstance()
	t2[1] = value.Undefined
	t2[2] = value.Undefined
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Elk::AST::ListPatternNode, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewVariablePatternDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(59, 3, 49))), ast.NewListPatternNode(position.NewLocation("<main>", position.NewSpan(position.New(20, 3, 10), position.New(51, 3, 41))), []ast.PatternNode{(t1).AsReference().(ast.PatternNode)}), ast.NewDoubleQuotedStringLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(55, 3, 45), position.New(59, 3, 49))), "foo"))
}
`,
		},

		"unquote constant": {
			input: `
				quote
					const unquote(:Bar) = "foo"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("Bar")
var sym3 = value.ToSymbol("to_ast_const_node")
var fn_method0 vm.NativeFunction // Std::Symbol.:to_ast_const_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.SymbolClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (sym2).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Symbol, name: to_ast_const_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewConstantDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(42, 3, 32))), "", (t1).AsReference().(ast.LiteralPatternNode), nil, ast.NewDoubleQuotedStringLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(38, 3, 28), position.New(42, 3, 32))), "foo"))
}
`,
		},
		"unquote_const": {
			input: `
				quote
					const unquote_const(:Bar) = "foo"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("Bar")
var sym3 = value.ToSymbol("to_ast_const_node")
var fn_method0 vm.NativeFunction // Std::Symbol.:to_ast_const_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.SymbolClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (sym2).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Symbol, name: to_ast_const_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewConstantDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(48, 3, 38))), "", (t1).AsReference().(ast.LiteralPatternNode), nil, ast.NewDoubleQuotedStringLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(44, 3, 34), position.New(48, 3, 38))), "foo"))
}
`,
		},
		"short unquote constant": {
			input: `
				quote
					const !{:Bar} = "foo"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("Bar")
var sym3 = value.ToSymbol("to_ast_const_node")
var fn_method0 vm.NativeFunction // Std::Symbol.:to_ast_const_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.SymbolClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (sym2).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Symbol, name: to_ast_const_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewConstantDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(36, 3, 26))), "", (t1).AsReference().(ast.LiteralPatternNode), nil, ast.NewDoubleQuotedStringLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(32, 3, 22), position.New(36, 3, 26))), "foo"))
}
`,
		},

		"unquote_ivar": {
			input: `
				quote
					var unquote_ivar(:foo): String?
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
var sym2 = value.ToSymbol("foo")
var sym3 = value.ToSymbol("to_ast_ivar_node")
var fn_method0 vm.NativeFunction // Std::Symbol.:to_ast_ivar_node

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 ast.Node
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.SymbolClass).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (sym2).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::Symbol, name: to_ast_ivar_node
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t3 = ast.NewInstanceVariableDeclarationNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(46, 3, 36))), "", (t1).AsReference().(ast.InstanceVariableNode), ast.NewNilableTypeNode(position.NewLocation("<main>", position.NewSpan(position.New(40, 3, 30), position.New(46, 3, 36))), ast.NewPublicConstantNode(position.NewLocation("<main>", position.NewSpan(position.New(40, 3, 30), position.New(45, 3, 35))), "String")))
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

func TestGoMacroExpansion(t *testing.T) {
	tests := goTestTable{
		"compile-time fibonacci": {
			input: `
				using Std::Elk::AST::*

				macro fib(i: IntLiteralNode)
					calc_fib := |n: Int|: Int ->
						return 1 if n < 3

						calc_fib(n - 2) + calc_fib(n - 1)
					end

					calc_fib(i.to_int).to_ast_node
				end

				a := fib!(10) * 2
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(110)).ToValue()
}
`,
		},
		"call a scoped macro": {
			input: `
				using Std::Elk::AST::*

				module Math
					macro fib(i: IntLiteralNode)
						calc_fib := |n: Int|: Int ->
							return 1 if n < 3

							calc_fib(n - 2) + calc_fib(n - 1)
						end

						calc_fib(i.to_int).to_ast_node
					end
				end

				a := Math::fib!(10) * 2
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

var const0 *value.Module // Math
var sym0 = value.ToSymbol("Math")

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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(110)).ToValue()
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
		"recursive fibonacci macro": {
			input: `
				using Std::Elk::AST::*

				macro fib(i: IntLiteralNode)
					int := i.to_int
					return try IntLiteralNode('1') if int < 3

					quote
						fib!(!{int - 1}) + fib!(!{int - 2})
					end
				end

				a := fib!(10) * 2
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(110)).ToValue()
}
`,
		},
		"define a class": {
			input: `
				using Std::Elk::AST::*

				macro box(name: ConstantNode, typ: TypeExpressionNode)
					quote
						class !{name}
							attr value: !{typ.type_node}
							init(@value: !{typ.type_node}); end
						end
					end
				end

				box!(BoxString, type String?)

				b := BoxString("foo")
				b.value
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

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("#init")
var fn_method1 vm.NativeFunction // BoxString.:value

var const0 *value.Class // BoxString
var sym0 = value.ToSymbol("BoxString")

var sym1 = value.ToSymbol("BoxString.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: BoxString.:#init, loc: <main>:13:5
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("value")

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
	var l0 value.Value // var b: BoxString
	_ = l0
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var t2 []value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	fn_method1 = vm.MethodToFunc((const0).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
	defer thread.PopNativeCallFrame()
	t1, err = fn_method0(thread, const0.CreateInstance(), (value.String("foo")).ToValue()) // receiver: BoxString, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = l0
	callFrame.SetNativeLineNumber(16)
	_, err = fn_method1(thread, t2) // receiver: BoxString, name: value
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("value"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // BoxString
	vm.Def(class, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.DefineGetter(class, sym3, 0)
	vm.DefineSetter(class, sym3, 0)
}
`,
		},
		"define a method": {
			input: `
				using Std::Elk::AST::*

				macro reader(name: IdentifierNode, typ: TypeExpressionNode)
					ivar := PublicInstanceVariableNode(name.value)

					quote
						var unquote_ivar(ivar): !{typ.type_node}
						def !{name}: !{typ.type_node}
							!{ivar}
						end
					end
				end

				class Foo
					reader!(bar, type Int | Float?)
				end

				b := Foo()
				r := b.bar
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

var sym3 = value.ToSymbol("main")
var sym5 = value.ToSymbol("bar")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:bar")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:bar, loc: <main>:16:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.GetInstanceVariable(self, 0), value.Undefined

}

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
	var l0 value.Value // var b: Foo
	_ = l0
	var l1 value.Value // var r: Std::Int | Std::Float | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = const0.CreateInstance()
	t1, err = fn_method0(thread, l0) // receiver: Foo, name: bar
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
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("bar"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(class, "bar", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
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
