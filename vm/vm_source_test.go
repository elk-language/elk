package vm

import (
	"strings"
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Represents a single VM source code test case.
type sourceTestCase struct {
	source         string
	wantStackTop   value.Value
	wantStdout     string
	wantRuntimeErr value.Value
	wantCompileErr errors.ErrorList
	teardown       func()
}

// Type of the compiler test table.
type sourceTestTable map[string]sourceTestCase

// Type of the simple compiler test table.
type simpleSourceTestTable map[string]value.Value

const testFileName = "sourceName"

// Create a new position in tests
var P = position.New

// Create a new span in tests
var S = position.NewSpan

// Create a new location in tests
func L(startPos, endPos *position.Position) *position.Location {
	return position.NewLocation(testFileName, startPos, endPos)
}

func vmSourceTest(tc sourceTestCase, t *testing.T) {
	t.Helper()

	chunk, gotCompileErr := compiler.CompileSource(testFileName, tc.source)
	if gotCompileErr != nil {
		if diff := cmp.Diff(tc.wantCompileErr, gotCompileErr, value.ValueComparerOptions...); diff != "" {
			t.Fatalf(diff)
		}
		return
	}
	var stdout strings.Builder
	vm := New(WithStdout(&stdout))
	gotStackTop, gotRuntimeErr := vm.InterpretTopLevel(chunk)
	gotStdout := stdout.String()
	if tc.teardown != nil {
		tc.teardown()
	}
	if diff := cmp.Diff(tc.wantStdout, gotStdout, value.ValueComparerOptions...); diff != "" {
		t.Fatalf(diff)
	}
	if diff := cmp.Diff(tc.wantRuntimeErr, gotRuntimeErr, value.ValueComparerOptions...); diff != "" {
		t.Fatalf(diff)
	}
	if tc.wantRuntimeErr != nil {
		return
	}
	if diff := cmp.Diff(tc.wantStackTop, gotStackTop, value.ValueComparerOptions...); diff != "" {
		t.Log(gotRuntimeErr)
		if gotStackTop != nil && tc.wantStackTop != nil {
			t.Logf("got: %s, want: %s", gotStackTop.Inspect(), tc.wantStackTop.Inspect())
		}
		t.Fatalf(diff)
	}
}

func vmSimpleSourceTest(source string, want value.Value, t *testing.T) {
	t.Helper()

	opts := []cmp.Option{
		cmp.AllowUnexported(value.Error{}, value.BigFloat{}, value.BigInt{}),
		cmpopts.IgnoreUnexported(value.Class{}, value.Module{}),
		cmpopts.IgnoreFields(value.Class{}, "ConstructorFunc"),
		value.FloatComparer,
		value.Float32Comparer,
		value.Float64Comparer,
		value.BigFloatComparer,
	}

	chunk, gotCompileErr := compiler.CompileSource(testFileName, source)
	if gotCompileErr != nil {
		t.Fatalf("Compile Error: %s", gotCompileErr.Error())
		return
	}
	var stdout strings.Builder
	vm := New(WithStdout(&stdout))
	got, gotRuntimeErr := vm.InterpretTopLevel(chunk)
	if gotRuntimeErr != nil {
		t.Fatalf("Runtime Error: %s", gotRuntimeErr.Inspect())
	}
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Logf("got: %s, want: %s", got.Inspect(), want.Inspect())
		t.Fatalf(diff)
	}
}

func TestVMSource_Locals(t *testing.T) {
	tests := sourceTestTable{
		"define and initialise a variable": {
			source:       "var a = 'foo'",
			wantStackTop: value.String("foo"),
		},
		"shadow a variable": {
			source: `
				var a = 10
				var b = do
					var a = 5
					a + 3
				end
				a + b
			`,
			wantStackTop: value.SmallInt(18),
		},
		"define and set a variable": {
			source: `
				var a = 'foo'
				a = a + ' bar'
				a
			`,
			wantStackTop: value.String("foo bar"),
		},
		"try to read an uninitialised variable": {
			source: `
				var a
				a
			`,
			wantCompileErr: errors.ErrorList{
				errors.NewError(L(P(15, 3, 5), P(15, 3, 5)), "can't access an uninitialised local: a"),
			},
		},
		"try to read a nonexistent variable": {
			source: `
				a
			`,
			wantCompileErr: errors.ErrorList{
				errors.NewError(L(P(5, 2, 5), P(5, 2, 5)), "undeclared variable: a"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_IfExpressions(t *testing.T) {
	tests := sourceTestTable{
		"return nil when condition is truthy and then is empty": {
			source:       "if true; end",
			wantStackTop: value.Nil,
		},
		"return nil when condition is falsy and then is empty": {
			source:       "if false; end",
			wantStackTop: value.Nil,
		},
		"execute the then branch": {
			source: `
				a := 5
				if a
					a = a + 2
				end
			`,
			wantStackTop: value.SmallInt(7),
		},
		"execute the empty else branch": {
			source: `
				a := 5
				if false
					a = a * 2
				end
			`,
			wantStackTop: value.Nil,
		},
		"execute the then branch instead of else": {
			source: `
				a := 5
				if a
					a = a + 2
				else
					a = 30
				end
			`,
			wantStackTop: value.SmallInt(7),
		},
		"execute the else branch instead of then": {
			source: `
				a := 5
				if nil
					a = a + 2
				else
					a = 30
				end
			`,
			wantStackTop: value.SmallInt(30),
		},
		"is an expression": {
			source: `
				a := 5
				b := if a
					"foo"
				else
					5
				end
				b
			`,
			wantStackTop: value.String("foo"),
		},
		"modifier binds more strongly than assignment": {
			source: `
				a := 5
				b := "foo" if a else 5
				b
			`,
			wantCompileErr: errors.ErrorList{
				errors.NewError(L(P(43, 4, 5), P(43, 4, 5)), "undeclared variable: b"),
			},
		},
		"modifier returns the left side if the condition is satisfied": {
			source: `
				a := 5
				"foo" if a else 5
			`,
			wantStackTop: value.String("foo"),
		},
		"modifier returns the right side if the condition is not satisfied": {
			source: `
				a := nil
				"foo" if a else 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"modifier returns nil when condition is not satisfied": {
			source: `
				a := nil
				"foo" if a
			`,
			wantStackTop: value.Nil,
		},
		"can access variables defined in the condition": {
			source: `
				a + " bar" if a := "foo"
			`,
			wantStackTop: value.String("foo bar"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_UnlessExpressions(t *testing.T) {
	tests := sourceTestTable{
		"return nil when condition is falsy and then is empty": {
			source:       "unless false; end",
			wantStackTop: value.Nil,
		},
		"return nil when condition is truthy and then is empty": {
			source:       "unless true; end",
			wantStackTop: value.Nil,
		},
		"execute the then branch": {
			source: `
				a := nil
				unless a
					a = 7
				end
			`,
			wantStackTop: value.SmallInt(7),
		},
		"execute the empty else branch": {
			source: `
				a := 5
				unless true
					a = a * 2
				end
			`,
			wantStackTop: value.Nil,
		},
		"execute the then branch instead of else": {
			source: `
				a := false
				unless a
					a = 10
				else
					a = a + 2
				end
			`,
			wantStackTop: value.SmallInt(10),
		},
		"execute the else branch instead of then": {
			source: `
				a := 5
				unless a
					a = 30
				else
					a = a + 2
				end
			`,
			wantStackTop: value.SmallInt(7),
		},
		"is an expression": {
			source: `
				a := 5
				b := unless a
					"foo"
				else
					5
				end
				b
			`,
			wantStackTop: value.SmallInt(5),
		},
		"modifier binds more strongly than assignment": {
			source: `
				a := 5
				b := "foo" unless a
				b
			`,
			wantCompileErr: errors.ErrorList{
				errors.NewError(L(P(40, 4, 5), P(40, 4, 5)), "undeclared variable: b"),
			},
		},
		"modifier returns the left side if the condition is satisfied": {
			source: `
				a := nil
				"foo" unless a
			`,
			wantStackTop: value.String("foo"),
		},
		"modifier returns nil if the condition is not satisfied": {
			source: `
				a := 5
				"foo" unless a
			`,
			wantStackTop: value.Nil,
		},
		"can access variables defined in the condition": {
			source: `
				a unless a := false
			`,
			wantStackTop: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LogicalOrOperator(t *testing.T) {
	tests := sourceTestTable{
		"return right operand if left is nil": {
			source:       "nil || 4",
			wantStackTop: value.SmallInt(4),
		},
		"return right operand (nil) if left is nil": {
			source:       "nil || nil",
			wantStackTop: value.Nil,
		},
		"return right operand (false) if left is nil": {
			source:       "nil || false",
			wantStackTop: value.False,
		},
		"return right operand if left is false": {
			source:       "false || 'foo'",
			wantStackTop: value.String("foo"),
		},
		"return left operand if it's truthy": {
			source:       "3 || 'foo'",
			wantStackTop: value.SmallInt(3),
		},
		"return right nested operand if left are falsy": {
			source:       "false || nil || 4",
			wantStackTop: value.SmallInt(4),
		},
		"return middle nested operand if left is falsy": {
			source:       "false || 2 || 5",
			wantStackTop: value.SmallInt(2),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LogicalAndOperator(t *testing.T) {
	tests := sourceTestTable{
		"return left operand if left is nil": {
			source:       "nil && 4",
			wantStackTop: value.Nil,
		},
		"return left operand if left is false": {
			source:       "false && 'foo'",
			wantStackTop: value.False,
		},
		"return right operand if left is truthy": {
			source:       "3 && 'foo'",
			wantStackTop: value.String("foo"),
		},
		"return right operand (false) if left is truthy": {
			source:       "3 && false",
			wantStackTop: value.False,
		},
		"return right nested operand if left are truthy": {
			source:       "4 && 'bar' && 16",
			wantStackTop: value.SmallInt(16),
		},
		"return middle nested operand if left is truthy": {
			source:       "4 && nil && 5",
			wantStackTop: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_NilCoalescingOperator(t *testing.T) {
	tests := sourceTestTable{
		"return right operand if left is nil": {
			source:       "nil ?? 4",
			wantStackTop: value.SmallInt(4),
		},
		"return right operand (nil) if left is nil": {
			source:       "nil ?? nil",
			wantStackTop: value.Nil,
		},
		"return right operand (false) if left is nil": {
			source:       "nil ?? false",
			wantStackTop: value.False,
		},
		"return left operand if left is false": {
			source:       "false ?? 'foo'",
			wantStackTop: value.False,
		},
		"return left operand if it's not nil": {
			source:       "3 ?? 'foo'",
			wantStackTop: value.SmallInt(3),
		},
		"return right nested operand if left are nil": {
			source:       "nil ?? nil ?? 4",
			wantStackTop: value.SmallInt(4),
		},
		"return middle nested operand if left is nil": {
			source:       "nil ?? false ?? 5",
			wantStackTop: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Exponentiate(t *testing.T) {
	tests := sourceTestTable{
		"Int64 ** Int64": {
			source:       "2i64 ** 10i64",
			wantStackTop: value.Int64(1024),
		},
		"Int64 ** Int32": {
			source: "2i64 ** 10i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int64`",
			),
			wantStackTop: value.Int64(2),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Modulo(t *testing.T) {
	tests := sourceTestTable{
		"Int64 % Int64": {
			source:       "29i64 % 3i64",
			wantStackTop: value.Int64(2),
		},
		"SmallInt % Float": {
			source:       "250 % 4.5",
			wantStackTop: value.Float(2.5),
		},
		"Int64 % Int32": {
			source: "11i64 % 2i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int64`",
			),
			wantStackTop: value.Int64(11),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_GreaterThan(t *testing.T) {
	tests := sourceTestTable{
		// String
		"'25' > '25'": {
			source:       "'25' > '25'",
			wantStackTop: value.False,
		},
		"'7' > '10'": {
			source:       "'7' > '10'",
			wantStackTop: value.True,
		},
		"'10' > '7'": {
			source:       "'10' > '7'",
			wantStackTop: value.False,
		},
		"'25' > '22'": {
			source:       "'25' > '22'",
			wantStackTop: value.True,
		},
		"'22' > '25'": {
			source:       "'22' > '25'",
			wantStackTop: value.False,
		},
		"'foo' > 'foo'": {
			source:       "'foo' > 'foo'",
			wantStackTop: value.False,
		},
		"'foo' > 'foa'": {
			source:       "'foo' > 'foa'",
			wantStackTop: value.True,
		},
		"'foa' > 'foo'": {
			source:       "'foa' > 'foo'",
			wantStackTop: value.False,
		},
		"'foo' > 'foo bar'": {
			source:       "'foo' > 'foo bar'",
			wantStackTop: value.False,
		},
		"'foo bar' > 'foo'": {
			source:       "'foo bar' > 'foo'",
			wantStackTop: value.True,
		},

		"'2' > c'2'": {
			source:       "'2' > c'2'",
			wantStackTop: value.False,
		},
		"'72' > c'7'": {
			source:       "'72' > c'7'",
			wantStackTop: value.True,
		},
		"'8' > c'7'": {
			source:       "'8' > c'7'",
			wantStackTop: value.True,
		},
		"'7' > c'8'": {
			source:       "'7' > c'8'",
			wantStackTop: value.False,
		},
		"'ba' > c'b'": {
			source:       "'ba' > c'b'",
			wantStackTop: value.True,
		},
		"'b' > c'a'": {
			source:       "'b' > c'a'",
			wantStackTop: value.True,
		},
		"'a' > c'b'": {
			source:       "'a' > c'b'",
			wantStackTop: value.False,
		},

		"'2' > 2.0": {
			source: "'2' > 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::String`",
			),
		},

		"'28' > 25.2bf": {
			source: "'28' > 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::String`",
			),
		},

		"'28.8' > 12.9f64": {
			source: "'28.8' > 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::String`",
			),
		},

		"'28.8' > 12.9f32": {
			source: "'28.8' > 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::String`",
			),
		},

		"'93' > 19i64": {
			source: "'93' > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::String`",
			),
		},

		"'93' > 19i32": {
			source: "'93' > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::String`",
			),
		},

		"'93' > 19i16": {
			source: "'93' > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::String`",
			),
		},

		"'93' > 19i8": {
			source: "'93' > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::String`",
			),
		},

		"'93' > 19u64": {
			source: "'93' > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::String`",
			),
		},

		"'93' > 19u32": {
			source: "'93' > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::String`",
			),
		},

		"'93' > 19u16": {
			source: "'93' > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::String`",
			),
		},

		"'93' > 19u8": {
			source: "'93' > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::String`",
			),
		},

		// Char
		"c'2' > c'2'": {
			source:       "c'2' > c'2'",
			wantStackTop: value.False,
		},
		"c'8' > c'7'": {
			source:       "c'8' > c'7'",
			wantStackTop: value.True,
		},
		"c'7' > c'8'": {
			source:       "c'7' > c'8'",
			wantStackTop: value.False,
		},
		"c'b' > c'a'": {
			source:       "c'b' > c'a'",
			wantStackTop: value.True,
		},
		"c'a' > c'b'": {
			source:       "c'a' > c'b'",
			wantStackTop: value.False,
		},

		"c'2' > '2'": {
			source:       "c'2' > '2'",
			wantStackTop: value.False,
		},
		"c'7' > '72'": {
			source:       "c'7' > '72'",
			wantStackTop: value.False,
		},
		"c'8' > '7'": {
			source:       "c'8' > '7'",
			wantStackTop: value.True,
		},
		"c'7' > '8'": {
			source:       "c'7' > '8'",
			wantStackTop: value.False,
		},
		"c'b' > 'a'": {
			source:       "c'b' > 'a'",
			wantStackTop: value.True,
		},
		"c'b' > 'ba'": {
			source:       "c'b' > 'ba'",
			wantStackTop: value.False,
		},
		"c'a' > 'b'": {
			source:       "c'a' > 'b'",
			wantStackTop: value.False,
		},

		"c'2' > 2.0": {
			source: "c'2' > 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Char`",
			),
		},
		"c'i' > 25.2bf": {
			source: "c'i' > 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Char`",
			),
		},
		"c'f' > 12.9f64": {
			source: "c'f' > 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Char`",
			),
		},
		"c'0' > 12.9f32": {
			source: "c'0' > 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Char`",
			),
		},
		"c'9' > 19i64": {
			source: "c'9' > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Char`",
			),
		},
		"c'u' > 19i32": {
			source: "c'u' > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Char`",
			),
		},
		"c'4' > 19i16": {
			source: "c'4' > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Char`",
			),
		},
		"c'6' > 19i8": {
			source: "c'6' > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Char`",
			),
		},
		"c'9' > 19u64": {
			source: "c'9' > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Char`",
			),
		},
		"c'u' > 19u32": {
			source: "c'u' > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Char`",
			),
		},
		"c'4' > 19u16": {
			source: "c'4' > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Char`",
			),
		},
		"c'6' > 19u8": {
			source: "c'6' > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Char`",
			),
		},

		// Int
		"25 > 25": {
			source:       "25 > 25",
			wantStackTop: value.False,
		},
		"25 > -25": {
			source:       "25 > -25",
			wantStackTop: value.True,
		},
		"-25 > 25": {
			source:       "-25 > 25",
			wantStackTop: value.False,
		},
		"13 > 7": {
			source:       "13 > 7",
			wantStackTop: value.True,
		},
		"7 > 13": {
			source:       "7 > 13",
			wantStackTop: value.False,
		},

		"25 > 25.0": {
			source:       "25 > 25.0",
			wantStackTop: value.False,
		},
		"25 > -25.0": {
			source:       "25 > -25.0",
			wantStackTop: value.True,
		},
		"-25 > 25.0": {
			source:       "-25 > 25.0",
			wantStackTop: value.False,
		},
		"13 > 7.0": {
			source:       "13 > 7.0",
			wantStackTop: value.True,
		},
		"7 > 13.0": {
			source:       "7 > 13.0",
			wantStackTop: value.False,
		},
		"7 > 7.5": {
			source:       "7 > 7.5",
			wantStackTop: value.False,
		},
		"7 > 6.9": {
			source:       "7 > 6.9",
			wantStackTop: value.True,
		},

		"25 > 25bf": {
			source:       "25 > 25bf",
			wantStackTop: value.False,
		},
		"25 > -25bf": {
			source:       "25 > -25bf",
			wantStackTop: value.True,
		},
		"-25 > 25bf": {
			source:       "-25 > 25bf",
			wantStackTop: value.False,
		},
		"13 > 7bf": {
			source:       "13 > 7bf",
			wantStackTop: value.True,
		},
		"7 > 13bf": {
			source:       "7 > 13bf",
			wantStackTop: value.False,
		},
		"7 > 7.5bf": {
			source:       "7 > 7.5bf",
			wantStackTop: value.False,
		},
		"7 > 6.9bf": {
			source:       "7 > 6.9bf",
			wantStackTop: value.True,
		},

		"6 > 19f64": {
			source: "6 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 > 19f32": {
			source: "6 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 > 19i64": {
			source: "6 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 > 19i32": {
			source: "6 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 > 19i16": {
			source: "6 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 > 19i8": {
			source: "6 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 > 19u64": {
			source: "6 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 > 19u32": {
			source: "6 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 > 19u16": {
			source: "6 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 > 19u8": {
			source: "6 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::SmallInt`",
			),
		},

		// Float
		"25.0 > 25.0": {
			source:       "25.0 > 25.0",
			wantStackTop: value.False,
		},
		"25.0 > -25.0": {
			source:       "25.0 > -25.0",
			wantStackTop: value.True,
		},
		"-25.0 > 25.0": {
			source:       "-25.0 > 25.0",
			wantStackTop: value.False,
		},
		"13.0 > 7.0": {
			source:       "13.0 > 7.0",
			wantStackTop: value.True,
		},
		"7.0 > 13.0": {
			source:       "7.0 > 13.0",
			wantStackTop: value.False,
		},
		"7.0 > 7.5": {
			source:       "7.0 > 7.5",
			wantStackTop: value.False,
		},
		"7.5 > 7.0": {
			source:       "7.5 > 7.0",
			wantStackTop: value.True,
		},
		"7.0 > 6.9": {
			source:       "7.0 > 6.9",
			wantStackTop: value.True,
		},

		"25.0 > 25": {
			source:       "25.0 > 25",
			wantStackTop: value.False,
		},
		"25.0 > -25": {
			source:       "25.0 > -25",
			wantStackTop: value.True,
		},
		"-25.0 > 25": {
			source:       "-25.0 > 25",
			wantStackTop: value.False,
		},
		"13.0 > 7": {
			source:       "13.0 > 7",
			wantStackTop: value.True,
		},
		"7.0 > 13": {
			source:       "7.0 > 13",
			wantStackTop: value.False,
		},
		"7.5 > 7": {
			source:       "7.5 > 7",
			wantStackTop: value.True,
		},

		"25.0 > 25bf": {
			source:       "25.0 > 25bf",
			wantStackTop: value.False,
		},
		"25.0 > -25bf": {
			source:       "25.0 > -25bf",
			wantStackTop: value.True,
		},
		"-25.0 > 25bf": {
			source:       "-25.0 > 25bf",
			wantStackTop: value.False,
		},
		"13.0 > 7bf": {
			source:       "13.0 > 7bf",
			wantStackTop: value.True,
		},
		"7.0 > 13bf": {
			source:       "7.0 > 13bf",
			wantStackTop: value.False,
		},
		"7.0 > 7.5bf": {
			source:       "7.0 > 7.5bf",
			wantStackTop: value.False,
		},
		"7.5 > 7bf": {
			source:       "7.5 > 7bf",
			wantStackTop: value.True,
		},
		"7.0 > 6.9bf": {
			source:       "7.0 > 6.9bf",
			wantStackTop: value.True,
		},

		"6.0 > 19f64": {
			source: "6.0 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 > 19f32": {
			source: "6.0 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 > 19i64": {
			source: "6.0 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 > 19i32": {
			source: "6.0 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 > 19i16": {
			source: "6.0 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float`",
			),
		},
		"6.0 > 19i8": {
			source: "6.0 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float`",
			),
		},
		"6.0 > 19u64": {
			source: "6.0 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 > 19u32": {
			source: "6.0 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 > 19u16": {
			source: "6.0 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float`",
			),
		},
		"6.0 > 19u8": {
			source: "6.0 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float`",
			),
		},

		// BigFloat
		"25bf > 25.0": {
			source:       "25bf > 25.0",
			wantStackTop: value.False,
		},
		"25bf > -25.0": {
			source:       "25bf > -25.0",
			wantStackTop: value.True,
		},
		"-25bf > 25.0": {
			source:       "-25bf > 25.0",
			wantStackTop: value.False,
		},
		"13bf > 7.0": {
			source:       "13bf > 7.0",
			wantStackTop: value.True,
		},
		"7bf > 13.0": {
			source:       "7bf > 13.0",
			wantStackTop: value.False,
		},
		"7bf > 7.5": {
			source:       "7bf > 7.5",
			wantStackTop: value.False,
		},
		"7.5bf > 7.0": {
			source:       "7.5bf > 7.0",
			wantStackTop: value.True,
		},
		"7bf > 6.9": {
			source:       "7bf > 6.9",
			wantStackTop: value.True,
		},

		"25bf > 25": {
			source:       "25bf > 25",
			wantStackTop: value.False,
		},
		"25bf > -25": {
			source:       "25bf > -25",
			wantStackTop: value.True,
		},
		"-25bf > 25": {
			source:       "-25bf > 25",
			wantStackTop: value.False,
		},
		"13bf > 7": {
			source:       "13bf > 7",
			wantStackTop: value.True,
		},
		"7bf > 13": {
			source:       "7bf > 13",
			wantStackTop: value.False,
		},
		"7.5bf > 7": {
			source:       "7.5bf > 7",
			wantStackTop: value.True,
		},

		"25bf > 25bf": {
			source:       "25bf > 25bf",
			wantStackTop: value.False,
		},
		"25bf > -25bf": {
			source:       "25bf > -25bf",
			wantStackTop: value.True,
		},
		"-25bf > 25bf": {
			source:       "-25bf > 25bf",
			wantStackTop: value.False,
		},
		"13bf > 7bf": {
			source:       "13bf > 7bf",
			wantStackTop: value.True,
		},
		"7bf > 13bf": {
			source:       "7bf > 13bf",
			wantStackTop: value.False,
		},
		"7bf > 7.5bf": {
			source:       "7bf > 7.5bf",
			wantStackTop: value.False,
		},
		"7.5bf > 7bf": {
			source:       "7.5bf > 7bf",
			wantStackTop: value.True,
		},
		"7bf > 6.9bf": {
			source:       "7bf > 6.9bf",
			wantStackTop: value.True,
		},

		"6bf > 19f64": {
			source: "6bf > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19f32": {
			source: "6bf > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19i64": {
			source: "6bf > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19i32": {
			source: "6bf > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19i16": {
			source: "6bf > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19i8": {
			source: "6bf > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19u64": {
			source: "6bf > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19u32": {
			source: "6bf > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19u16": {
			source: "6bf > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf > 19u8": {
			source: "6bf > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::BigFloat`",
			),
		},

		// Float64
		"25f64 > 25f64": {
			source:       "25f64 > 25f64",
			wantStackTop: value.False,
		},
		"25f64 > -25f64": {
			source:       "25f64 > -25f64",
			wantStackTop: value.True,
		},
		"-25f64 > 25f64": {
			source:       "-25f64 > 25f64",
			wantStackTop: value.False,
		},
		"13f64 > 7f64": {
			source:       "13f64 > 7f64",
			wantStackTop: value.True,
		},
		"7f64 > 13f64": {
			source:       "7f64 > 13f64",
			wantStackTop: value.False,
		},
		"7f64 > 7.5f64": {
			source:       "7f64 > 7.5f64",
			wantStackTop: value.False,
		},
		"7.5f64 > 7f64": {
			source:       "7.5f64 > 7f64",
			wantStackTop: value.True,
		},
		"7f64 > 6.9f64": {
			source:       "7f64 > 6.9f64",
			wantStackTop: value.True,
		},

		"6f64 > 19.0": {
			source: "6f64 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Float64`",
			),
		},

		"6f64 > 19": {
			source: "6f64 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19bf": {
			source: "6f64 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19f32": {
			source: "6f64 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19i64": {
			source: "6f64 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19i32": {
			source: "6f64 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19i16": {
			source: "6f64 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19i8": {
			source: "6f64 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19u64": {
			source: "6f64 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19u32": {
			source: "6f64 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19u16": {
			source: "6f64 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 > 19u8": {
			source: "6f64 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float64`",
			),
		},

		// Float32
		"25f32 > 25f32": {
			source:       "25f32 > 25f32",
			wantStackTop: value.False,
		},
		"25f32 > -25f32": {
			source:       "25f32 > -25f32",
			wantStackTop: value.True,
		},
		"-25f32 > 25f32": {
			source:       "-25f32 > 25f32",
			wantStackTop: value.False,
		},
		"13f32 > 7f32": {
			source:       "13f32 > 7f32",
			wantStackTop: value.True,
		},
		"7f32 > 13f32": {
			source:       "7f32 > 13f32",
			wantStackTop: value.False,
		},
		"7f32 > 7.5f32": {
			source:       "7f32 > 7.5f32",
			wantStackTop: value.False,
		},
		"7.5f32 > 7f32": {
			source:       "7.5f32 > 7f32",
			wantStackTop: value.True,
		},
		"7f32 > 6.9f32": {
			source:       "7f32 > 6.9f32",
			wantStackTop: value.True,
		},

		"6f32 > 19.0": {
			source: "6f32 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Float32`",
			),
		},

		"6f32 > 19": {
			source: "6f32 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19bf": {
			source: "6f32 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19f64": {
			source: "6f32 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19i64": {
			source: "6f32 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19i32": {
			source: "6f32 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19i16": {
			source: "6f32 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19i8": {
			source: "6f32 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19u64": {
			source: "6f32 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19u32": {
			source: "6f32 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19u16": {
			source: "6f32 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 > 19u8": {
			source: "6f32 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float32`",
			),
		},

		// Int64
		"25i64 > 25i64": {
			source:       "25i64 > 25i64",
			wantStackTop: value.False,
		},
		"25i64 > -25i64": {
			source:       "25i64 > -25i64",
			wantStackTop: value.True,
		},
		"-25i64 > 25i64": {
			source:       "-25i64 > 25i64",
			wantStackTop: value.False,
		},
		"13i64 > 7i64": {
			source:       "13i64 > 7i64",
			wantStackTop: value.True,
		},
		"7i64 > 13i64": {
			source:       "7i64 > 13i64",
			wantStackTop: value.False,
		},

		"6i64 > 19": {
			source: "6i64 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19.0": {
			source: "6i64 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19bf": {
			source: "6i64 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19f64": {
			source: "6i64 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19f32": {
			source: "6i64 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19i32": {
			source: "6i64 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19i16": {
			source: "6i64 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19i8": {
			source: "6i64 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19u64": {
			source: "6i64 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19u32": {
			source: "6i64 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19u16": {
			source: "6i64 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 > 19u8": {
			source: "6i64 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int64`",
			),
		},

		// Int32
		"25i32 > 25i32": {
			source:       "25i32 > 25i32",
			wantStackTop: value.False,
		},
		"25i32 > -25i32": {
			source:       "25i32 > -25i32",
			wantStackTop: value.True,
		},
		"-25i32 > 25i32": {
			source:       "-25i32 > 25i32",
			wantStackTop: value.False,
		},
		"13i32 > 7i32": {
			source:       "13i32 > 7i32",
			wantStackTop: value.True,
		},
		"7i32 > 13i32": {
			source:       "7i32 > 13i32",
			wantStackTop: value.False,
		},

		"6i32 > 19": {
			source: "6i32 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19.0": {
			source: "6i32 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19bf": {
			source: "6i32 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19f64": {
			source: "6i32 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19f32": {
			source: "6i32 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19i64": {
			source: "6i32 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19i16": {
			source: "6i32 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19i8": {
			source: "6i32 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19u64": {
			source: "6i32 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19u32": {
			source: "6i32 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19u16": {
			source: "6i32 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 > 19u8": {
			source: "6i32 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int32`",
			),
		},

		// Int16
		"25i16 > 25i16": {
			source:       "25i16 > 25i16",
			wantStackTop: value.False,
		},
		"25i16 > -25i16": {
			source:       "25i16 > -25i16",
			wantStackTop: value.True,
		},
		"-25i16 > 25i16": {
			source:       "-25i16 > 25i16",
			wantStackTop: value.False,
		},
		"13i16 > 7i16": {
			source:       "13i16 > 7i16",
			wantStackTop: value.True,
		},
		"7i16 > 13i16": {
			source:       "7i16 > 13i16",
			wantStackTop: value.False,
		},

		"6i16 > 19": {
			source: "6i16 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19.0": {
			source: "6i16 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19bf": {
			source: "6i16 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19f64": {
			source: "6i16 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19f32": {
			source: "6i16 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19i64": {
			source: "6i16 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19i32": {
			source: "6i16 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19i8": {
			source: "6i16 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19u64": {
			source: "6i16 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19u32": {
			source: "6i16 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19u16": {
			source: "6i16 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 > 19u8": {
			source: "6i16 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int16`",
			),
		},

		// Int8
		"25i8 > 25i8": {
			source:       "25i8 > 25i8",
			wantStackTop: value.False,
		},
		"25i8 > -25i8": {
			source:       "25i8 > -25i8",
			wantStackTop: value.True,
		},
		"-25i8 > 25i8": {
			source:       "-25i8 > 25i8",
			wantStackTop: value.False,
		},
		"13i8 > 7i8": {
			source:       "13i8 > 7i8",
			wantStackTop: value.True,
		},
		"7i8 > 13i8": {
			source:       "7i8 > 13i8",
			wantStackTop: value.False,
		},

		"6i8 > 19": {
			source: "6i8 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19.0": {
			source: "6i8 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19bf": {
			source: "6i8 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19f64": {
			source: "6i8 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19f32": {
			source: "6i8 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19i64": {
			source: "6i8 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19i32": {
			source: "6i8 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19i16": {
			source: "6i8 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19u64": {
			source: "6i8 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19u32": {
			source: "6i8 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19u16": {
			source: "6i8 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 > 19u8": {
			source: "6i8 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int8`",
			),
		},

		// UInt64
		"25u64 > 25u64": {
			source:       "25u64 > 25u64",
			wantStackTop: value.False,
		},
		"13u64 > 7u64": {
			source:       "13u64 > 7u64",
			wantStackTop: value.True,
		},
		"7u64 > 13u64": {
			source:       "7u64 > 13u64",
			wantStackTop: value.False,
		},

		"6u64 > 19": {
			source: "6u64 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19.0": {
			source: "6u64 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19bf": {
			source: "6u64 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19f64": {
			source: "6u64 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19f32": {
			source: "6u64 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19i64": {
			source: "6u64 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19i32": {
			source: "6u64 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19i16": {
			source: "6u64 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19i8": {
			source: "6u64 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19u32": {
			source: "6u64 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19u16": {
			source: "6u64 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 > 19u8": {
			source: "6u64 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt64`",
			),
		},

		// UInt32
		"25u32 > 25u32": {
			source:       "25u32 > 25u32",
			wantStackTop: value.False,
		},
		"13u32 > 7u32": {
			source:       "13u32 > 7u32",
			wantStackTop: value.True,
		},
		"7u32 > 13u32": {
			source:       "7u32 > 13u32",
			wantStackTop: value.False,
		},

		"6u32 > 19": {
			source: "6u32 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19.0": {
			source: "6u32 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19bf": {
			source: "6u32 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19f64": {
			source: "6u32 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19f32": {
			source: "6u32 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19i64": {
			source: "6u32 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19i32": {
			source: "6u32 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19i16": {
			source: "6u32 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19i8": {
			source: "6u32 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19u64": {
			source: "6u32 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19u16": {
			source: "6u32 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 > 19u8": {
			source: "6u32 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt32`",
			),
		},

		// Int16
		"25u16 > 25u16": {
			source:       "25u16 > 25u16",
			wantStackTop: value.False,
		},
		"13u16 > 7u16": {
			source:       "13u16 > 7u16",
			wantStackTop: value.True,
		},
		"7u16 > 13u16": {
			source:       "7u16 > 13u16",
			wantStackTop: value.False,
		},

		"6u16 > 19": {
			source: "6u16 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19.0": {
			source: "6u16 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19bf": {
			source: "6u16 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19f64": {
			source: "6u16 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19f32": {
			source: "6u16 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19i64": {
			source: "6u16 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19i32": {
			source: "6u16 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19i16": {
			source: "6u16 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19i8": {
			source: "6u16 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19u64": {
			source: "6u16 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19u32": {
			source: "6u16 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 > 19u8": {
			source: "6u16 > 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt16`",
			),
		},

		// Int8
		"25u8 > 25u8": {
			source:       "25u8 > 25u8",
			wantStackTop: value.False,
		},
		"13u8 > 7u8": {
			source:       "13u8 > 7u8",
			wantStackTop: value.True,
		},
		"7u8 > 13u8": {
			source:       "7u8 > 13u8",
			wantStackTop: value.False,
		},

		"6u8 > 19": {
			source: "6u8 > 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19.0": {
			source: "6u8 > 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19bf": {
			source: "6u8 > 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19f64": {
			source: "6u8 > 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19f32": {
			source: "6u8 > 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19i64": {
			source: "6u8 > 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19i32": {
			source: "6u8 > 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19i16": {
			source: "6u8 > 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19i8": {
			source: "6u8 > 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19u64": {
			source: "6u8 > 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19u32": {
			source: "6u8 > 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 > 19u16": {
			source: "6u8 > 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt8`",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_GreaterThanEqual(t *testing.T) {
	tests := sourceTestTable{
		// String
		"'25' >= '25'": {
			source:       "'25' >= '25'",
			wantStackTop: value.True,
		},
		"'7' >= '10'": {
			source:       "'7' >= '10'",
			wantStackTop: value.True,
		},
		"'10' >= '7'": {
			source:       "'10' >= '7'",
			wantStackTop: value.False,
		},
		"'25' >= '22'": {
			source:       "'25' >= '22'",
			wantStackTop: value.True,
		},
		"'22' >= '25'": {
			source:       "'22' >= '25'",
			wantStackTop: value.False,
		},
		"'foo' >= 'foo'": {
			source:       "'foo' >= 'foo'",
			wantStackTop: value.True,
		},
		"'foo' >= 'foa'": {
			source:       "'foo' >= 'foa'",
			wantStackTop: value.True,
		},
		"'foa' >= 'foo'": {
			source:       "'foa' >= 'foo'",
			wantStackTop: value.False,
		},
		"'foo' >= 'foo bar'": {
			source:       "'foo' >= 'foo bar'",
			wantStackTop: value.False,
		},
		"'foo bar' >= 'foo'": {
			source:       "'foo bar' >= 'foo'",
			wantStackTop: value.True,
		},

		"'2' >= c'2'": {
			source:       "'2' >= c'2'",
			wantStackTop: value.True,
		},
		"'72' >= c'7'": {
			source:       "'72' >= c'7'",
			wantStackTop: value.True,
		},
		"'8' >= c'7'": {
			source:       "'8' >= c'7'",
			wantStackTop: value.True,
		},
		"'7' >= c'8'": {
			source:       "'7' >= c'8'",
			wantStackTop: value.False,
		},
		"'ba' >= c'b'": {
			source:       "'ba' >= c'b'",
			wantStackTop: value.True,
		},
		"'b' >= c'a'": {
			source:       "'b' >= c'a'",
			wantStackTop: value.True,
		},
		"'a' >= c'b'": {
			source:       "'a' >= c'b'",
			wantStackTop: value.False,
		},

		"'2' >= 2.0": {
			source: "'2' >= 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::String`",
			),
		},

		"'28' >= 25.2bf": {
			source: "'28' >= 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::String`",
			),
		},

		"'28.8' >= 12.9f64": {
			source: "'28.8' >= 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::String`",
			),
		},

		"'28.8' >= 12.9f32": {
			source: "'28.8' >= 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::String`",
			),
		},

		"'93' >= 19i64": {
			source: "'93' >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::String`",
			),
		},

		"'93' >= 19i32": {
			source: "'93' >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::String`",
			),
		},

		"'93' >= 19i16": {
			source: "'93' >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::String`",
			),
		},

		"'93' >= 19i8": {
			source: "'93' >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::String`",
			),
		},

		"'93' >= 19u64": {
			source: "'93' >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::String`",
			),
		},

		"'93' >= 19u32": {
			source: "'93' >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::String`",
			),
		},

		"'93' >= 19u16": {
			source: "'93' >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::String`",
			),
		},

		"'93' >= 19u8": {
			source: "'93' >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::String`",
			),
		},

		// Char
		"c'2' >= c'2'": {
			source:       "c'2' >= c'2'",
			wantStackTop: value.True,
		},
		"c'8' >= c'7'": {
			source:       "c'8' >= c'7'",
			wantStackTop: value.True,
		},
		"c'7' >= c'8'": {
			source:       "c'7' >= c'8'",
			wantStackTop: value.False,
		},
		"c'b' >= c'a'": {
			source:       "c'b' >= c'a'",
			wantStackTop: value.True,
		},
		"c'a' >= c'b'": {
			source:       "c'a' >= c'b'",
			wantStackTop: value.False,
		},

		"c'2' >= '2'": {
			source:       "c'2' >= '2'",
			wantStackTop: value.True,
		},
		"c'7' >= '72'": {
			source:       "c'7' >= '72'",
			wantStackTop: value.False,
		},
		"c'8' >= '7'": {
			source:       "c'8' >= '7'",
			wantStackTop: value.True,
		},
		"c'7' >= '8'": {
			source:       "c'7' >= '8'",
			wantStackTop: value.False,
		},
		"c'b' >= 'a'": {
			source:       "c'b' >= 'a'",
			wantStackTop: value.True,
		},
		"c'b' >= 'ba'": {
			source:       "c'b' >= 'ba'",
			wantStackTop: value.False,
		},
		"c'a' >= 'b'": {
			source:       "c'a' >= 'b'",
			wantStackTop: value.False,
		},

		"c'2' >= 2.0": {
			source: "c'2' >= 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Char`",
			),
		},
		"c'i' >= 25.2bf": {
			source: "c'i' >= 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Char`",
			),
		},
		"c'f' >= 12.9f64": {
			source: "c'f' >= 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Char`",
			),
		},
		"c'0' >= 12.9f32": {
			source: "c'0' >= 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Char`",
			),
		},
		"c'9' >= 19i64": {
			source: "c'9' >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Char`",
			),
		},
		"c'u' >= 19i32": {
			source: "c'u' >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Char`",
			),
		},
		"c'4' >= 19i16": {
			source: "c'4' >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Char`",
			),
		},
		"c'6' >= 19i8": {
			source: "c'6' >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Char`",
			),
		},
		"c'9' >= 19u64": {
			source: "c'9' >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Char`",
			),
		},
		"c'u' >= 19u32": {
			source: "c'u' >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Char`",
			),
		},
		"c'4' >= 19u16": {
			source: "c'4' >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Char`",
			),
		},
		"c'6' >= 19u8": {
			source: "c'6' >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Char`",
			),
		},

		// Int
		"25 >= 25": {
			source:       "25 >= 25",
			wantStackTop: value.True,
		},
		"25 >= -25": {
			source:       "25 >= -25",
			wantStackTop: value.True,
		},
		"-25 >= 25": {
			source:       "-25 >= 25",
			wantStackTop: value.False,
		},
		"13 >= 7": {
			source:       "13 >= 7",
			wantStackTop: value.True,
		},
		"7 >= 13": {
			source:       "7 >= 13",
			wantStackTop: value.False,
		},

		"25 >= 25.0": {
			source:       "25 >= 25.0",
			wantStackTop: value.True,
		},
		"25 >= -25.0": {
			source:       "25 >= -25.0",
			wantStackTop: value.True,
		},
		"-25 >= 25.0": {
			source:       "-25 >= 25.0",
			wantStackTop: value.False,
		},
		"13 >= 7.0": {
			source:       "13 >= 7.0",
			wantStackTop: value.True,
		},
		"7 >= 13.0": {
			source:       "7 >= 13.0",
			wantStackTop: value.False,
		},
		"7 >= 7.5": {
			source:       "7 >= 7.5",
			wantStackTop: value.False,
		},
		"7 >= 6.9": {
			source:       "7 >= 6.9",
			wantStackTop: value.True,
		},

		"25 >= 25bf": {
			source:       "25 >= 25bf",
			wantStackTop: value.True,
		},
		"25 >= -25bf": {
			source:       "25 >= -25bf",
			wantStackTop: value.True,
		},
		"-25 >= 25bf": {
			source:       "-25 >= 25bf",
			wantStackTop: value.False,
		},
		"13 >= 7bf": {
			source:       "13 >= 7bf",
			wantStackTop: value.True,
		},
		"7 >= 13bf": {
			source:       "7 >= 13bf",
			wantStackTop: value.False,
		},
		"7 >= 7.5bf": {
			source:       "7 >= 7.5bf",
			wantStackTop: value.False,
		},
		"7 >= 6.9bf": {
			source:       "7 >= 6.9bf",
			wantStackTop: value.True,
		},

		"6 >= 19f64": {
			source: "6 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 >= 19f32": {
			source: "6 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 >= 19i64": {
			source: "6 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 >= 19i32": {
			source: "6 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 >= 19i16": {
			source: "6 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 >= 19i8": {
			source: "6 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 >= 19u64": {
			source: "6 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 >= 19u32": {
			source: "6 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 >= 19u16": {
			source: "6 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 >= 19u8": {
			source: "6 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::SmallInt`",
			),
		},

		// Float
		"25.0 >= 25.0": {
			source:       "25.0 >= 25.0",
			wantStackTop: value.True,
		},
		"25.0 >= -25.0": {
			source:       "25.0 >= -25.0",
			wantStackTop: value.True,
		},
		"-25.0 >= 25.0": {
			source:       "-25.0 >= 25.0",
			wantStackTop: value.False,
		},
		"13.0 >= 7.0": {
			source:       "13.0 >= 7.0",
			wantStackTop: value.True,
		},
		"7.0 >= 13.0": {
			source:       "7.0 >= 13.0",
			wantStackTop: value.False,
		},
		"7.0 >= 7.5": {
			source:       "7.0 >= 7.5",
			wantStackTop: value.False,
		},
		"7.5 >= 7.0": {
			source:       "7.5 >= 7.0",
			wantStackTop: value.True,
		},
		"7.0 >= 6.9": {
			source:       "7.0 >= 6.9",
			wantStackTop: value.True,
		},

		"25.0 >= 25": {
			source:       "25.0 >= 25",
			wantStackTop: value.True,
		},
		"25.0 >= -25": {
			source:       "25.0 >= -25",
			wantStackTop: value.True,
		},
		"-25.0 >= 25": {
			source:       "-25.0 >= 25",
			wantStackTop: value.False,
		},
		"13.0 >= 7": {
			source:       "13.0 >= 7",
			wantStackTop: value.True,
		},
		"7.0 >= 13": {
			source:       "7.0 >= 13",
			wantStackTop: value.False,
		},
		"7.5 >= 7": {
			source:       "7.5 >= 7",
			wantStackTop: value.True,
		},

		"25.0 >= 25bf": {
			source:       "25.0 >= 25bf",
			wantStackTop: value.True,
		},
		"25.0 >= -25bf": {
			source:       "25.0 >= -25bf",
			wantStackTop: value.True,
		},
		"-25.0 >= 25bf": {
			source:       "-25.0 >= 25bf",
			wantStackTop: value.False,
		},
		"13.0 >= 7bf": {
			source:       "13.0 >= 7bf",
			wantStackTop: value.True,
		},
		"7.0 >= 13bf": {
			source:       "7.0 >= 13bf",
			wantStackTop: value.False,
		},
		"7.0 >= 7.5bf": {
			source:       "7.0 >= 7.5bf",
			wantStackTop: value.False,
		},
		"7.5 >= 7bf": {
			source:       "7.5 >= 7bf",
			wantStackTop: value.True,
		},
		"7.0 >= 6.9bf": {
			source:       "7.0 >= 6.9bf",
			wantStackTop: value.True,
		},

		"6.0 >= 19f64": {
			source: "6.0 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19f32": {
			source: "6.0 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19i64": {
			source: "6.0 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19i32": {
			source: "6.0 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19i16": {
			source: "6.0 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19i8": {
			source: "6.0 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19u64": {
			source: "6.0 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19u32": {
			source: "6.0 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19u16": {
			source: "6.0 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float`",
			),
		},
		"6.0 >= 19u8": {
			source: "6.0 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float`",
			),
		},

		// BigFloat
		"25bf >= 25.0": {
			source:       "25bf >= 25.0",
			wantStackTop: value.True,
		},
		"25bf >= -25.0": {
			source:       "25bf >= -25.0",
			wantStackTop: value.True,
		},
		"-25bf >= 25.0": {
			source:       "-25bf >= 25.0",
			wantStackTop: value.False,
		},
		"13bf >= 7.0": {
			source:       "13bf >= 7.0",
			wantStackTop: value.True,
		},
		"7bf >= 13.0": {
			source:       "7bf >= 13.0",
			wantStackTop: value.False,
		},
		"7bf >= 7.5": {
			source:       "7bf >= 7.5",
			wantStackTop: value.False,
		},
		"7.5bf >= 7.0": {
			source:       "7.5bf >= 7.0",
			wantStackTop: value.True,
		},
		"7bf >= 6.9": {
			source:       "7bf >= 6.9",
			wantStackTop: value.True,
		},

		"25bf >= 25": {
			source:       "25bf >= 25",
			wantStackTop: value.True,
		},
		"25bf >= -25": {
			source:       "25bf >= -25",
			wantStackTop: value.True,
		},
		"-25bf >= 25": {
			source:       "-25bf >= 25",
			wantStackTop: value.False,
		},
		"13bf >= 7": {
			source:       "13bf >= 7",
			wantStackTop: value.True,
		},
		"7bf >= 13": {
			source:       "7bf >= 13",
			wantStackTop: value.False,
		},
		"7.5bf >= 7": {
			source:       "7.5bf >= 7",
			wantStackTop: value.True,
		},

		"25bf >= 25bf": {
			source:       "25bf >= 25bf",
			wantStackTop: value.True,
		},
		"25bf >= -25bf": {
			source:       "25bf >= -25bf",
			wantStackTop: value.True,
		},
		"-25bf >= 25bf": {
			source:       "-25bf >= 25bf",
			wantStackTop: value.False,
		},
		"13bf >= 7bf": {
			source:       "13bf >= 7bf",
			wantStackTop: value.True,
		},
		"7bf >= 13bf": {
			source:       "7bf >= 13bf",
			wantStackTop: value.False,
		},
		"7bf >= 7.5bf": {
			source:       "7bf >= 7.5bf",
			wantStackTop: value.False,
		},
		"7.5bf >= 7bf": {
			source:       "7.5bf >= 7bf",
			wantStackTop: value.True,
		},
		"7bf >= 6.9bf": {
			source:       "7bf >= 6.9bf",
			wantStackTop: value.True,
		},

		"6bf >= 19f64": {
			source: "6bf >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19f32": {
			source: "6bf >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19i64": {
			source: "6bf >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19i32": {
			source: "6bf >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19i16": {
			source: "6bf >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19i8": {
			source: "6bf >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19u64": {
			source: "6bf >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19u32": {
			source: "6bf >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19u16": {
			source: "6bf >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf >= 19u8": {
			source: "6bf >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::BigFloat`",
			),
		},

		// Float64
		"25f64 >= 25f64": {
			source:       "25f64 >= 25f64",
			wantStackTop: value.True,
		},
		"25f64 >= -25f64": {
			source:       "25f64 >= -25f64",
			wantStackTop: value.True,
		},
		"-25f64 >= 25f64": {
			source:       "-25f64 >= 25f64",
			wantStackTop: value.False,
		},
		"13f64 >= 7f64": {
			source:       "13f64 >= 7f64",
			wantStackTop: value.True,
		},
		"7f64 >= 13f64": {
			source:       "7f64 >= 13f64",
			wantStackTop: value.False,
		},
		"7f64 >= 7.5f64": {
			source:       "7f64 >= 7.5f64",
			wantStackTop: value.False,
		},
		"7.5f64 >= 7f64": {
			source:       "7.5f64 >= 7f64",
			wantStackTop: value.True,
		},
		"7f64 >= 6.9f64": {
			source:       "7f64 >= 6.9f64",
			wantStackTop: value.True,
		},

		"6f64 >= 19.0": {
			source: "6f64 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Float64`",
			),
		},

		"6f64 >= 19": {
			source: "6f64 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19bf": {
			source: "6f64 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19f32": {
			source: "6f64 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19i64": {
			source: "6f64 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19i32": {
			source: "6f64 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19i16": {
			source: "6f64 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19i8": {
			source: "6f64 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19u64": {
			source: "6f64 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19u32": {
			source: "6f64 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19u16": {
			source: "6f64 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 >= 19u8": {
			source: "6f64 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float64`",
			),
		},

		// Float32
		"25f32 >= 25f32": {
			source:       "25f32 >= 25f32",
			wantStackTop: value.True,
		},
		"25f32 >= -25f32": {
			source:       "25f32 >= -25f32",
			wantStackTop: value.True,
		},
		"-25f32 >= 25f32": {
			source:       "-25f32 >= 25f32",
			wantStackTop: value.False,
		},
		"13f32 >= 7f32": {
			source:       "13f32 >= 7f32",
			wantStackTop: value.True,
		},
		"7f32 >= 13f32": {
			source:       "7f32 >= 13f32",
			wantStackTop: value.False,
		},
		"7f32 >= 7.5f32": {
			source:       "7f32 >= 7.5f32",
			wantStackTop: value.False,
		},
		"7.5f32 >= 7f32": {
			source:       "7.5f32 >= 7f32",
			wantStackTop: value.True,
		},
		"7f32 >= 6.9f32": {
			source:       "7f32 >= 6.9f32",
			wantStackTop: value.True,
		},

		"6f32 >= 19.0": {
			source: "6f32 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Float32`",
			),
		},

		"6f32 >= 19": {
			source: "6f32 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19bf": {
			source: "6f32 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19f64": {
			source: "6f32 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19i64": {
			source: "6f32 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19i32": {
			source: "6f32 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19i16": {
			source: "6f32 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19i8": {
			source: "6f32 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19u64": {
			source: "6f32 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19u32": {
			source: "6f32 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19u16": {
			source: "6f32 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 >= 19u8": {
			source: "6f32 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float32`",
			),
		},

		// Int64
		"25i64 >= 25i64": {
			source:       "25i64 >= 25i64",
			wantStackTop: value.True,
		},
		"25i64 >= -25i64": {
			source:       "25i64 >= -25i64",
			wantStackTop: value.True,
		},
		"-25i64 >= 25i64": {
			source:       "-25i64 >= 25i64",
			wantStackTop: value.False,
		},
		"13i64 >= 7i64": {
			source:       "13i64 >= 7i64",
			wantStackTop: value.True,
		},
		"7i64 >= 13i64": {
			source:       "7i64 >= 13i64",
			wantStackTop: value.False,
		},

		"6i64 >= 19": {
			source: "6i64 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19.0": {
			source: "6i64 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19bf": {
			source: "6i64 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19f64": {
			source: "6i64 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19f32": {
			source: "6i64 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19i32": {
			source: "6i64 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19i16": {
			source: "6i64 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19i8": {
			source: "6i64 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19u64": {
			source: "6i64 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19u32": {
			source: "6i64 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19u16": {
			source: "6i64 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 >= 19u8": {
			source: "6i64 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int64`",
			),
		},

		// Int32
		"25i32 >= 25i32": {
			source:       "25i32 >= 25i32",
			wantStackTop: value.True,
		},
		"25i32 >= -25i32": {
			source:       "25i32 >= -25i32",
			wantStackTop: value.True,
		},
		"-25i32 >= 25i32": {
			source:       "-25i32 >= 25i32",
			wantStackTop: value.False,
		},
		"13i32 >= 7i32": {
			source:       "13i32 >= 7i32",
			wantStackTop: value.True,
		},
		"7i32 >= 13i32": {
			source:       "7i32 >= 13i32",
			wantStackTop: value.False,
		},

		"6i32 >= 19": {
			source: "6i32 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19.0": {
			source: "6i32 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19bf": {
			source: "6i32 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19f64": {
			source: "6i32 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19f32": {
			source: "6i32 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19i64": {
			source: "6i32 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19i16": {
			source: "6i32 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19i8": {
			source: "6i32 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19u64": {
			source: "6i32 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19u32": {
			source: "6i32 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19u16": {
			source: "6i32 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 >= 19u8": {
			source: "6i32 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int32`",
			),
		},

		// Int16
		"25i16 >= 25i16": {
			source:       "25i16 >= 25i16",
			wantStackTop: value.True,
		},
		"25i16 >= -25i16": {
			source:       "25i16 >= -25i16",
			wantStackTop: value.True,
		},
		"-25i16 >= 25i16": {
			source:       "-25i16 >= 25i16",
			wantStackTop: value.False,
		},
		"13i16 >= 7i16": {
			source:       "13i16 >= 7i16",
			wantStackTop: value.True,
		},
		"7i16 >= 13i16": {
			source:       "7i16 >= 13i16",
			wantStackTop: value.False,
		},

		"6i16 >= 19": {
			source: "6i16 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19.0": {
			source: "6i16 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19bf": {
			source: "6i16 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19f64": {
			source: "6i16 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19f32": {
			source: "6i16 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19i64": {
			source: "6i16 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19i32": {
			source: "6i16 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19i8": {
			source: "6i16 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19u64": {
			source: "6i16 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19u32": {
			source: "6i16 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19u16": {
			source: "6i16 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 >= 19u8": {
			source: "6i16 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int16`",
			),
		},

		// Int8
		"25i8 >= 25i8": {
			source:       "25i8 >= 25i8",
			wantStackTop: value.True,
		},
		"25i8 >= -25i8": {
			source:       "25i8 >= -25i8",
			wantStackTop: value.True,
		},
		"-25i8 >= 25i8": {
			source:       "-25i8 >= 25i8",
			wantStackTop: value.False,
		},
		"13i8 >= 7i8": {
			source:       "13i8 >= 7i8",
			wantStackTop: value.True,
		},
		"7i8 >= 13i8": {
			source:       "7i8 >= 13i8",
			wantStackTop: value.False,
		},

		"6i8 >= 19": {
			source: "6i8 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19.0": {
			source: "6i8 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19bf": {
			source: "6i8 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19f64": {
			source: "6i8 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19f32": {
			source: "6i8 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19i64": {
			source: "6i8 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19i32": {
			source: "6i8 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19i16": {
			source: "6i8 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19u64": {
			source: "6i8 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19u32": {
			source: "6i8 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19u16": {
			source: "6i8 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 >= 19u8": {
			source: "6i8 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int8`",
			),
		},

		// UInt64
		"25u64 >= 25u64": {
			source:       "25u64 >= 25u64",
			wantStackTop: value.True,
		},
		"13u64 >= 7u64": {
			source:       "13u64 >= 7u64",
			wantStackTop: value.True,
		},
		"7u64 >= 13u64": {
			source:       "7u64 >= 13u64",
			wantStackTop: value.False,
		},

		"6u64 >= 19": {
			source: "6u64 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19.0": {
			source: "6u64 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19bf": {
			source: "6u64 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19f64": {
			source: "6u64 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19f32": {
			source: "6u64 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19i64": {
			source: "6u64 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19i32": {
			source: "6u64 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19i16": {
			source: "6u64 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19i8": {
			source: "6u64 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19u32": {
			source: "6u64 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19u16": {
			source: "6u64 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 >= 19u8": {
			source: "6u64 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt64`",
			),
		},

		// UInt32
		"25u32 >= 25u32": {
			source:       "25u32 >= 25u32",
			wantStackTop: value.True,
		},
		"13u32 >= 7u32": {
			source:       "13u32 >= 7u32",
			wantStackTop: value.True,
		},
		"7u32 >= 13u32": {
			source:       "7u32 >= 13u32",
			wantStackTop: value.False,
		},

		"6u32 >= 19": {
			source: "6u32 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19.0": {
			source: "6u32 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19bf": {
			source: "6u32 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19f64": {
			source: "6u32 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19f32": {
			source: "6u32 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19i64": {
			source: "6u32 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19i32": {
			source: "6u32 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19i16": {
			source: "6u32 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19i8": {
			source: "6u32 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19u64": {
			source: "6u32 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19u16": {
			source: "6u32 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 >= 19u8": {
			source: "6u32 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt32`",
			),
		},

		// Int16
		"25u16 >= 25u16": {
			source:       "25u16 >= 25u16",
			wantStackTop: value.True,
		},
		"13u16 >= 7u16": {
			source:       "13u16 >= 7u16",
			wantStackTop: value.True,
		},
		"7u16 >= 13u16": {
			source:       "7u16 >= 13u16",
			wantStackTop: value.False,
		},

		"6u16 >= 19": {
			source: "6u16 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19.0": {
			source: "6u16 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19bf": {
			source: "6u16 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19f64": {
			source: "6u16 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19f32": {
			source: "6u16 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19i64": {
			source: "6u16 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19i32": {
			source: "6u16 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19i16": {
			source: "6u16 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19i8": {
			source: "6u16 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19u64": {
			source: "6u16 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19u32": {
			source: "6u16 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 >= 19u8": {
			source: "6u16 >= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt16`",
			),
		},

		// Int8
		"25u8 >= 25u8": {
			source:       "25u8 >= 25u8",
			wantStackTop: value.True,
		},
		"13u8 >= 7u8": {
			source:       "13u8 >= 7u8",
			wantStackTop: value.True,
		},
		"7u8 >= 13u8": {
			source:       "7u8 >= 13u8",
			wantStackTop: value.False,
		},

		"6u8 >= 19": {
			source: "6u8 >= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19.0": {
			source: "6u8 >= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19bf": {
			source: "6u8 >= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19f64": {
			source: "6u8 >= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19f32": {
			source: "6u8 >= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19i64": {
			source: "6u8 >= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19i32": {
			source: "6u8 >= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19i16": {
			source: "6u8 >= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19i8": {
			source: "6u8 >= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19u64": {
			source: "6u8 >= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19u32": {
			source: "6u8 >= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 >= 19u16": {
			source: "6u8 >= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt8`",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LessThan(t *testing.T) {
	tests := sourceTestTable{
		// String
		"'25' < '25'": {
			source:       "'25' < '25'",
			wantStackTop: value.False,
		},
		"'7' < '10'": {
			source:       "'7' < '10'",
			wantStackTop: value.False,
		},
		"'10' < '7'": {
			source:       "'10' < '7'",
			wantStackTop: value.True,
		},
		"'25' < '22'": {
			source:       "'25' < '22'",
			wantStackTop: value.False,
		},
		"'22' < '25'": {
			source:       "'22' < '25'",
			wantStackTop: value.True,
		},
		"'foo' < 'foo'": {
			source:       "'foo' < 'foo'",
			wantStackTop: value.False,
		},
		"'foo' < 'foa'": {
			source:       "'foo' < 'foa'",
			wantStackTop: value.False,
		},
		"'foa' < 'foo'": {
			source:       "'foa' < 'foo'",
			wantStackTop: value.True,
		},
		"'foo' < 'foo bar'": {
			source:       "'foo' < 'foo bar'",
			wantStackTop: value.True,
		},
		"'foo bar' < 'foo'": {
			source:       "'foo bar' < 'foo'",
			wantStackTop: value.False,
		},

		"'2' < c'2'": {
			source:       "'2' < c'2'",
			wantStackTop: value.False,
		},
		"'72' < c'7'": {
			source:       "'72' < c'7'",
			wantStackTop: value.False,
		},
		"'8' < c'7'": {
			source:       "'8' < c'7'",
			wantStackTop: value.False,
		},
		"'7' < c'8'": {
			source:       "'7' < c'8'",
			wantStackTop: value.True,
		},
		"'ba' < c'b'": {
			source:       "'ba' < c'b'",
			wantStackTop: value.False,
		},
		"'b' < c'a'": {
			source:       "'b' < c'a'",
			wantStackTop: value.False,
		},
		"'a' < c'b'": {
			source:       "'a' < c'b'",
			wantStackTop: value.True,
		},

		"'2' < 2.0": {
			source: "'2' < 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::String`",
			),
		},

		"'28' < 25.2bf": {
			source: "'28' < 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::String`",
			),
		},

		"'28.8' < 12.9f64": {
			source: "'28.8' < 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::String`",
			),
		},

		"'28.8' < 12.9f32": {
			source: "'28.8' < 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::String`",
			),
		},

		"'93' < 19i64": {
			source: "'93' < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::String`",
			),
		},

		"'93' < 19i32": {
			source: "'93' < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::String`",
			),
		},

		"'93' < 19i16": {
			source: "'93' < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::String`",
			),
		},

		"'93' < 19i8": {
			source: "'93' < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::String`",
			),
		},

		"'93' < 19u64": {
			source: "'93' < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::String`",
			),
		},

		"'93' < 19u32": {
			source: "'93' < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::String`",
			),
		},

		"'93' < 19u16": {
			source: "'93' < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::String`",
			),
		},

		"'93' < 19u8": {
			source: "'93' < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::String`",
			),
		},

		// Char
		"c'2' < c'2'": {
			source:       "c'2' < c'2'",
			wantStackTop: value.False,
		},
		"c'8' < c'7'": {
			source:       "c'8' < c'7'",
			wantStackTop: value.False,
		},
		"c'7' < c'8'": {
			source:       "c'7' < c'8'",
			wantStackTop: value.True,
		},
		"c'b' < c'a'": {
			source:       "c'b' < c'a'",
			wantStackTop: value.False,
		},
		"c'a' < c'b'": {
			source:       "c'a' < c'b'",
			wantStackTop: value.True,
		},

		"c'2' < '2'": {
			source:       "c'2' < '2'",
			wantStackTop: value.False,
		},
		"c'7' < '72'": {
			source:       "c'7' < '72'",
			wantStackTop: value.True,
		},
		"c'8' < '7'": {
			source:       "c'8' < '7'",
			wantStackTop: value.False,
		},
		"c'7' < '8'": {
			source:       "c'7' < '8'",
			wantStackTop: value.True,
		},
		"c'b' < 'a'": {
			source:       "c'b' < 'a'",
			wantStackTop: value.False,
		},
		"c'b' < 'ba'": {
			source:       "c'b' < 'ba'",
			wantStackTop: value.True,
		},
		"c'a' < 'b'": {
			source:       "c'a' < 'b'",
			wantStackTop: value.True,
		},

		"c'2' < 2.0": {
			source: "c'2' < 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Char`",
			),
		},
		"c'i' < 25.2bf": {
			source: "c'i' < 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Char`",
			),
		},
		"c'f' < 12.9f64": {
			source: "c'f' < 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Char`",
			),
		},
		"c'0' < 12.9f32": {
			source: "c'0' < 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Char`",
			),
		},
		"c'9' < 19i64": {
			source: "c'9' < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Char`",
			),
		},
		"c'u' < 19i32": {
			source: "c'u' < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Char`",
			),
		},
		"c'4' < 19i16": {
			source: "c'4' < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Char`",
			),
		},
		"c'6' < 19i8": {
			source: "c'6' < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Char`",
			),
		},
		"c'9' < 19u64": {
			source: "c'9' < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Char`",
			),
		},
		"c'u' < 19u32": {
			source: "c'u' < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Char`",
			),
		},
		"c'4' < 19u16": {
			source: "c'4' < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Char`",
			),
		},
		"c'6' < 19u8": {
			source: "c'6' < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Char`",
			),
		},

		// Int
		"25 < 25": {
			source:       "25 < 25",
			wantStackTop: value.False,
		},
		"25 < -25": {
			source:       "25 < -25",
			wantStackTop: value.False,
		},
		"-25 < 25": {
			source:       "-25 < 25",
			wantStackTop: value.True,
		},
		"13 < 7": {
			source:       "13 < 7",
			wantStackTop: value.False,
		},
		"7 < 13": {
			source:       "7 < 13",
			wantStackTop: value.True,
		},

		"25 < 25.0": {
			source:       "25 < 25.0",
			wantStackTop: value.False,
		},
		"25 < -25.0": {
			source:       "25 < -25.0",
			wantStackTop: value.False,
		},
		"-25 < 25.0": {
			source:       "-25 < 25.0",
			wantStackTop: value.True,
		},
		"13 < 7.0": {
			source:       "13 < 7.0",
			wantStackTop: value.False,
		},
		"7 < 13.0": {
			source:       "7 < 13.0",
			wantStackTop: value.True,
		},
		"7 < 7.5": {
			source:       "7 < 7.5",
			wantStackTop: value.True,
		},
		"7 < 6.9": {
			source:       "7 < 6.9",
			wantStackTop: value.False,
		},

		"25 < 25bf": {
			source:       "25 < 25bf",
			wantStackTop: value.False,
		},
		"25 < -25bf": {
			source:       "25 < -25bf",
			wantStackTop: value.False,
		},
		"-25 < 25bf": {
			source:       "-25 < 25bf",
			wantStackTop: value.True,
		},
		"13 < 7bf": {
			source:       "13 < 7bf",
			wantStackTop: value.False,
		},
		"7 < 13bf": {
			source:       "7 < 13bf",
			wantStackTop: value.True,
		},
		"7 < 7.5bf": {
			source:       "7 < 7.5bf",
			wantStackTop: value.True,
		},
		"7 < 6.9bf": {
			source:       "7 < 6.9bf",
			wantStackTop: value.False,
		},

		"6 < 19f64": {
			source: "6 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 < 19f32": {
			source: "6 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 < 19i64": {
			source: "6 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 < 19i32": {
			source: "6 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 < 19i16": {
			source: "6 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 < 19i8": {
			source: "6 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 < 19u64": {
			source: "6 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 < 19u32": {
			source: "6 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 < 19u16": {
			source: "6 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 < 19u8": {
			source: "6 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::SmallInt`",
			),
		},

		// Float
		"25.0 < 25.0": {
			source:       "25.0 < 25.0",
			wantStackTop: value.False,
		},
		"25.0 < -25.0": {
			source:       "25.0 < -25.0",
			wantStackTop: value.False,
		},
		"-25.0 < 25.0": {
			source:       "-25.0 < 25.0",
			wantStackTop: value.True,
		},
		"13.0 < 7.0": {
			source:       "13.0 < 7.0",
			wantStackTop: value.False,
		},
		"7.0 < 13.0": {
			source:       "7.0 < 13.0",
			wantStackTop: value.True,
		},
		"7.0 < 7.5": {
			source:       "7.0 < 7.5",
			wantStackTop: value.True,
		},
		"7.5 < 7.0": {
			source:       "7.5 < 7.0",
			wantStackTop: value.False,
		},
		"7.0 < 6.9": {
			source:       "7.0 < 6.9",
			wantStackTop: value.False,
		},

		"25.0 < 25": {
			source:       "25.0 < 25",
			wantStackTop: value.False,
		},
		"25.0 < -25": {
			source:       "25.0 < -25",
			wantStackTop: value.False,
		},
		"-25.0 < 25": {
			source:       "-25.0 < 25",
			wantStackTop: value.True,
		},
		"13.0 < 7": {
			source:       "13.0 < 7",
			wantStackTop: value.False,
		},
		"7.0 < 13": {
			source:       "7.0 < 13",
			wantStackTop: value.True,
		},
		"7.5 < 7": {
			source:       "7.5 < 7",
			wantStackTop: value.False,
		},

		"25.0 < 25bf": {
			source:       "25.0 < 25bf",
			wantStackTop: value.False,
		},
		"25.0 < -25bf": {
			source:       "25.0 < -25bf",
			wantStackTop: value.False,
		},
		"-25.0 < 25bf": {
			source:       "-25.0 < 25bf",
			wantStackTop: value.True,
		},
		"13.0 < 7bf": {
			source:       "13.0 < 7bf",
			wantStackTop: value.False,
		},
		"7.0 < 13bf": {
			source:       "7.0 < 13bf",
			wantStackTop: value.True,
		},
		"7.0 < 7.5bf": {
			source:       "7.0 < 7.5bf",
			wantStackTop: value.True,
		},
		"7.5 < 7bf": {
			source:       "7.5 < 7bf",
			wantStackTop: value.False,
		},
		"7.0 < 6.9bf": {
			source:       "7.0 < 6.9bf",
			wantStackTop: value.False,
		},

		"6.0 < 19f64": {
			source: "6.0 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 < 19f32": {
			source: "6.0 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 < 19i64": {
			source: "6.0 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 < 19i32": {
			source: "6.0 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 < 19i16": {
			source: "6.0 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float`",
			),
		},
		"6.0 < 19i8": {
			source: "6.0 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float`",
			),
		},
		"6.0 < 19u64": {
			source: "6.0 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 < 19u32": {
			source: "6.0 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 < 19u16": {
			source: "6.0 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float`",
			),
		},
		"6.0 < 19u8": {
			source: "6.0 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float`",
			),
		},

		// BigFloat
		"25bf < 25.0": {
			source:       "25bf < 25.0",
			wantStackTop: value.False,
		},
		"25bf < -25.0": {
			source:       "25bf < -25.0",
			wantStackTop: value.False,
		},
		"-25bf < 25.0": {
			source:       "-25bf < 25.0",
			wantStackTop: value.True,
		},
		"13bf < 7.0": {
			source:       "13bf < 7.0",
			wantStackTop: value.False,
		},
		"7bf < 13.0": {
			source:       "7bf < 13.0",
			wantStackTop: value.True,
		},
		"7bf < 7.5": {
			source:       "7bf < 7.5",
			wantStackTop: value.True,
		},
		"7.5bf < 7.0": {
			source:       "7.5bf < 7.0",
			wantStackTop: value.False,
		},
		"7bf < 6.9": {
			source:       "7bf < 6.9",
			wantStackTop: value.False,
		},

		"25bf < 25": {
			source:       "25bf < 25",
			wantStackTop: value.False,
		},
		"25bf < -25": {
			source:       "25bf < -25",
			wantStackTop: value.False,
		},
		"-25bf < 25": {
			source:       "-25bf < 25",
			wantStackTop: value.True,
		},
		"13bf < 7": {
			source:       "13bf < 7",
			wantStackTop: value.False,
		},
		"7bf < 13": {
			source:       "7bf < 13",
			wantStackTop: value.True,
		},
		"7.5bf < 7": {
			source:       "7.5bf < 7",
			wantStackTop: value.False,
		},

		"25bf < 25bf": {
			source:       "25bf < 25bf",
			wantStackTop: value.False,
		},
		"25bf < -25bf": {
			source:       "25bf < -25bf",
			wantStackTop: value.False,
		},
		"-25bf < 25bf": {
			source:       "-25bf < 25bf",
			wantStackTop: value.True,
		},
		"13bf < 7bf": {
			source:       "13bf < 7bf",
			wantStackTop: value.False,
		},
		"7bf < 13bf": {
			source:       "7bf < 13bf",
			wantStackTop: value.True,
		},
		"7bf < 7.5bf": {
			source:       "7bf < 7.5bf",
			wantStackTop: value.True,
		},
		"7.5bf < 7bf": {
			source:       "7.5bf < 7bf",
			wantStackTop: value.False,
		},
		"7bf < 6.9bf": {
			source:       "7bf < 6.9bf",
			wantStackTop: value.False,
		},

		"6bf < 19f64": {
			source: "6bf < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19f32": {
			source: "6bf < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19i64": {
			source: "6bf < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19i32": {
			source: "6bf < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19i16": {
			source: "6bf < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19i8": {
			source: "6bf < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19u64": {
			source: "6bf < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19u32": {
			source: "6bf < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19u16": {
			source: "6bf < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf < 19u8": {
			source: "6bf < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::BigFloat`",
			),
		},

		// Float64
		"25f64 < 25f64": {
			source:       "25f64 < 25f64",
			wantStackTop: value.False,
		},
		"25f64 < -25f64": {
			source:       "25f64 < -25f64",
			wantStackTop: value.False,
		},
		"-25f64 < 25f64": {
			source:       "-25f64 < 25f64",
			wantStackTop: value.True,
		},
		"13f64 < 7f64": {
			source:       "13f64 < 7f64",
			wantStackTop: value.False,
		},
		"7f64 < 13f64": {
			source:       "7f64 < 13f64",
			wantStackTop: value.True,
		},
		"7f64 < 7.5f64": {
			source:       "7f64 < 7.5f64",
			wantStackTop: value.True,
		},
		"7.5f64 < 7f64": {
			source:       "7.5f64 < 7f64",
			wantStackTop: value.False,
		},
		"7f64 < 6.9f64": {
			source:       "7f64 < 6.9f64",
			wantStackTop: value.False,
		},

		"6f64 < 19.0": {
			source: "6f64 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Float64`",
			),
		},

		"6f64 < 19": {
			source: "6f64 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19bf": {
			source: "6f64 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19f32": {
			source: "6f64 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19i64": {
			source: "6f64 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19i32": {
			source: "6f64 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19i16": {
			source: "6f64 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19i8": {
			source: "6f64 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19u64": {
			source: "6f64 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19u32": {
			source: "6f64 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19u16": {
			source: "6f64 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 < 19u8": {
			source: "6f64 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float64`",
			),
		},

		// Float32
		"25f32 < 25f32": {
			source:       "25f32 < 25f32",
			wantStackTop: value.False,
		},
		"25f32 < -25f32": {
			source:       "25f32 < -25f32",
			wantStackTop: value.False,
		},
		"-25f32 < 25f32": {
			source:       "-25f32 < 25f32",
			wantStackTop: value.True,
		},
		"13f32 < 7f32": {
			source:       "13f32 < 7f32",
			wantStackTop: value.False,
		},
		"7f32 < 13f32": {
			source:       "7f32 < 13f32",
			wantStackTop: value.True,
		},
		"7f32 < 7.5f32": {
			source:       "7f32 < 7.5f32",
			wantStackTop: value.True,
		},
		"7.5f32 < 7f32": {
			source:       "7.5f32 < 7f32",
			wantStackTop: value.False,
		},
		"7f32 < 6.9f32": {
			source:       "7f32 < 6.9f32",
			wantStackTop: value.False,
		},

		"6f32 < 19.0": {
			source: "6f32 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Float32`",
			),
		},

		"6f32 < 19": {
			source: "6f32 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19bf": {
			source: "6f32 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19f64": {
			source: "6f32 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19i64": {
			source: "6f32 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19i32": {
			source: "6f32 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19i16": {
			source: "6f32 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19i8": {
			source: "6f32 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19u64": {
			source: "6f32 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19u32": {
			source: "6f32 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19u16": {
			source: "6f32 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 < 19u8": {
			source: "6f32 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float32`",
			),
		},

		// Int64
		"25i64 < 25i64": {
			source:       "25i64 < 25i64",
			wantStackTop: value.False,
		},
		"25i64 < -25i64": {
			source:       "25i64 < -25i64",
			wantStackTop: value.False,
		},
		"-25i64 < 25i64": {
			source:       "-25i64 < 25i64",
			wantStackTop: value.True,
		},
		"13i64 < 7i64": {
			source:       "13i64 < 7i64",
			wantStackTop: value.False,
		},
		"7i64 < 13i64": {
			source:       "7i64 < 13i64",
			wantStackTop: value.True,
		},

		"6i64 < 19": {
			source: "6i64 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19.0": {
			source: "6i64 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19bf": {
			source: "6i64 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19f64": {
			source: "6i64 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19f32": {
			source: "6i64 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19i32": {
			source: "6i64 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19i16": {
			source: "6i64 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19i8": {
			source: "6i64 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19u64": {
			source: "6i64 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19u32": {
			source: "6i64 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19u16": {
			source: "6i64 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 < 19u8": {
			source: "6i64 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int64`",
			),
		},

		// Int32
		"25i32 < 25i32": {
			source:       "25i32 < 25i32",
			wantStackTop: value.False,
		},
		"25i32 < -25i32": {
			source:       "25i32 < -25i32",
			wantStackTop: value.False,
		},
		"-25i32 < 25i32": {
			source:       "-25i32 < 25i32",
			wantStackTop: value.True,
		},
		"13i32 < 7i32": {
			source:       "13i32 < 7i32",
			wantStackTop: value.False,
		},
		"7i32 < 13i32": {
			source:       "7i32 < 13i32",
			wantStackTop: value.True,
		},

		"6i32 < 19": {
			source: "6i32 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19.0": {
			source: "6i32 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19bf": {
			source: "6i32 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19f64": {
			source: "6i32 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19f32": {
			source: "6i32 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19i64": {
			source: "6i32 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19i16": {
			source: "6i32 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19i8": {
			source: "6i32 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19u64": {
			source: "6i32 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19u32": {
			source: "6i32 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19u16": {
			source: "6i32 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 < 19u8": {
			source: "6i32 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int32`",
			),
		},

		// Int16
		"25i16 < 25i16": {
			source:       "25i16 < 25i16",
			wantStackTop: value.False,
		},
		"25i16 < -25i16": {
			source:       "25i16 < -25i16",
			wantStackTop: value.False,
		},
		"-25i16 < 25i16": {
			source:       "-25i16 < 25i16",
			wantStackTop: value.True,
		},
		"13i16 < 7i16": {
			source:       "13i16 < 7i16",
			wantStackTop: value.False,
		},
		"7i16 < 13i16": {
			source:       "7i16 < 13i16",
			wantStackTop: value.True,
		},

		"6i16 < 19": {
			source: "6i16 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19.0": {
			source: "6i16 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19bf": {
			source: "6i16 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19f64": {
			source: "6i16 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19f32": {
			source: "6i16 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19i64": {
			source: "6i16 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19i32": {
			source: "6i16 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19i8": {
			source: "6i16 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19u64": {
			source: "6i16 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19u32": {
			source: "6i16 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19u16": {
			source: "6i16 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 < 19u8": {
			source: "6i16 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int16`",
			),
		},

		// Int8
		"25i8 < 25i8": {
			source:       "25i8 < 25i8",
			wantStackTop: value.False,
		},
		"25i8 < -25i8": {
			source:       "25i8 < -25i8",
			wantStackTop: value.False,
		},
		"-25i8 < 25i8": {
			source:       "-25i8 < 25i8",
			wantStackTop: value.True,
		},
		"13i8 < 7i8": {
			source:       "13i8 < 7i8",
			wantStackTop: value.False,
		},
		"7i8 < 13i8": {
			source:       "7i8 < 13i8",
			wantStackTop: value.True,
		},

		"6i8 < 19": {
			source: "6i8 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19.0": {
			source: "6i8 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19bf": {
			source: "6i8 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19f64": {
			source: "6i8 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19f32": {
			source: "6i8 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19i64": {
			source: "6i8 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19i32": {
			source: "6i8 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19i16": {
			source: "6i8 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19u64": {
			source: "6i8 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19u32": {
			source: "6i8 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19u16": {
			source: "6i8 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 < 19u8": {
			source: "6i8 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int8`",
			),
		},

		// UInt64
		"25u64 < 25u64": {
			source:       "25u64 < 25u64",
			wantStackTop: value.False,
		},
		"13u64 < 7u64": {
			source:       "13u64 < 7u64",
			wantStackTop: value.False,
		},
		"7u64 < 13u64": {
			source:       "7u64 < 13u64",
			wantStackTop: value.True,
		},

		"6u64 < 19": {
			source: "6u64 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19.0": {
			source: "6u64 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19bf": {
			source: "6u64 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19f64": {
			source: "6u64 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19f32": {
			source: "6u64 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19i64": {
			source: "6u64 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19i32": {
			source: "6u64 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19i16": {
			source: "6u64 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19i8": {
			source: "6u64 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19u32": {
			source: "6u64 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19u16": {
			source: "6u64 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 < 19u8": {
			source: "6u64 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt64`",
			),
		},

		// UInt32
		"25u32 < 25u32": {
			source:       "25u32 < 25u32",
			wantStackTop: value.False,
		},
		"13u32 < 7u32": {
			source:       "13u32 < 7u32",
			wantStackTop: value.False,
		},
		"7u32 < 13u32": {
			source:       "7u32 < 13u32",
			wantStackTop: value.True,
		},

		"6u32 < 19": {
			source: "6u32 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19.0": {
			source: "6u32 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19bf": {
			source: "6u32 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19f64": {
			source: "6u32 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19f32": {
			source: "6u32 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19i64": {
			source: "6u32 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19i32": {
			source: "6u32 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19i16": {
			source: "6u32 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19i8": {
			source: "6u32 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19u64": {
			source: "6u32 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19u16": {
			source: "6u32 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 < 19u8": {
			source: "6u32 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt32`",
			),
		},

		// Int16
		"25u16 < 25u16": {
			source:       "25u16 < 25u16",
			wantStackTop: value.False,
		},
		"13u16 < 7u16": {
			source:       "13u16 < 7u16",
			wantStackTop: value.False,
		},
		"7u16 < 13u16": {
			source:       "7u16 < 13u16",
			wantStackTop: value.True,
		},

		"6u16 < 19": {
			source: "6u16 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19.0": {
			source: "6u16 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19bf": {
			source: "6u16 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19f64": {
			source: "6u16 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19f32": {
			source: "6u16 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19i64": {
			source: "6u16 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19i32": {
			source: "6u16 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19i16": {
			source: "6u16 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19i8": {
			source: "6u16 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19u64": {
			source: "6u16 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19u32": {
			source: "6u16 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 < 19u8": {
			source: "6u16 < 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt16`",
			),
		},

		// Int8
		"25u8 < 25u8": {
			source:       "25u8 < 25u8",
			wantStackTop: value.False,
		},
		"13u8 < 7u8": {
			source:       "13u8 < 7u8",
			wantStackTop: value.False,
		},
		"7u8 < 13u8": {
			source:       "7u8 < 13u8",
			wantStackTop: value.True,
		},

		"6u8 < 19": {
			source: "6u8 < 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19.0": {
			source: "6u8 < 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19bf": {
			source: "6u8 < 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19f64": {
			source: "6u8 < 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19f32": {
			source: "6u8 < 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19i64": {
			source: "6u8 < 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19i32": {
			source: "6u8 < 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19i16": {
			source: "6u8 < 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19i8": {
			source: "6u8 < 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19u64": {
			source: "6u8 < 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19u32": {
			source: "6u8 < 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 < 19u16": {
			source: "6u8 < 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt8`",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LessThanEqual(t *testing.T) {
	tests := sourceTestTable{
		// String
		"'25' <= '25'": {
			source:       "'25' <= '25'",
			wantStackTop: value.True,
		},
		"'7' <= '10'": {
			source:       "'7' <= '10'",
			wantStackTop: value.False,
		},
		"'10' <= '7'": {
			source:       "'10' <= '7'",
			wantStackTop: value.True,
		},
		"'25' <= '22'": {
			source:       "'25' <= '22'",
			wantStackTop: value.False,
		},
		"'22' <= '25'": {
			source:       "'22' <= '25'",
			wantStackTop: value.True,
		},
		"'foo' <= 'foo'": {
			source:       "'foo' <= 'foo'",
			wantStackTop: value.True,
		},
		"'foo' <= 'foa'": {
			source:       "'foo' <= 'foa'",
			wantStackTop: value.False,
		},
		"'foa' <= 'foo'": {
			source:       "'foa' <= 'foo'",
			wantStackTop: value.True,
		},
		"'foo' <= 'foo bar'": {
			source:       "'foo' <= 'foo bar'",
			wantStackTop: value.True,
		},
		"'foo bar' <= 'foo'": {
			source:       "'foo bar' <= 'foo'",
			wantStackTop: value.False,
		},

		"'2' <= c'2'": {
			source:       "'2' <= c'2'",
			wantStackTop: value.True,
		},
		"'72' <= c'7'": {
			source:       "'72' <= c'7'",
			wantStackTop: value.False,
		},
		"'8' <= c'7'": {
			source:       "'8' <= c'7'",
			wantStackTop: value.False,
		},
		"'7' <= c'8'": {
			source:       "'7' <= c'8'",
			wantStackTop: value.True,
		},
		"'ba' <= c'b'": {
			source:       "'ba' <= c'b'",
			wantStackTop: value.False,
		},
		"'b' <= c'a'": {
			source:       "'b' <= c'a'",
			wantStackTop: value.False,
		},
		"'a' <= c'b'": {
			source:       "'a' <= c'b'",
			wantStackTop: value.True,
		},

		"'2' <= 2.0": {
			source: "'2' <= 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::String`",
			),
		},

		"'28' <= 25.2bf": {
			source: "'28' <= 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::String`",
			),
		},

		"'28.8' <= 12.9f64": {
			source: "'28.8' <= 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::String`",
			),
		},

		"'28.8' <= 12.9f32": {
			source: "'28.8' <= 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::String`",
			),
		},

		"'93' <= 19i64": {
			source: "'93' <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::String`",
			),
		},

		"'93' <= 19i32": {
			source: "'93' <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::String`",
			),
		},

		"'93' <= 19i16": {
			source: "'93' <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::String`",
			),
		},

		"'93' <= 19i8": {
			source: "'93' <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::String`",
			),
		},

		"'93' <= 19u64": {
			source: "'93' <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::String`",
			),
		},

		"'93' <= 19u32": {
			source: "'93' <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::String`",
			),
		},

		"'93' <= 19u16": {
			source: "'93' <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::String`",
			),
		},

		"'93' <= 19u8": {
			source: "'93' <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::String`",
			),
		},

		// Char
		"c'2' <= c'2'": {
			source:       "c'2' <= c'2'",
			wantStackTop: value.True,
		},
		"c'8' <= c'7'": {
			source:       "c'8' <= c'7'",
			wantStackTop: value.False,
		},
		"c'7' <= c'8'": {
			source:       "c'7' <= c'8'",
			wantStackTop: value.True,
		},
		"c'b' <= c'a'": {
			source:       "c'b' <= c'a'",
			wantStackTop: value.False,
		},
		"c'a' <= c'b'": {
			source:       "c'a' <= c'b'",
			wantStackTop: value.True,
		},

		"c'2' <= '2'": {
			source:       "c'2' <= '2'",
			wantStackTop: value.True,
		},
		"c'7' <= '72'": {
			source:       "c'7' <= '72'",
			wantStackTop: value.True,
		},
		"c'8' <= '7'": {
			source:       "c'8' <= '7'",
			wantStackTop: value.False,
		},
		"c'7' <= '8'": {
			source:       "c'7' <= '8'",
			wantStackTop: value.True,
		},
		"c'b' <= 'a'": {
			source:       "c'b' <= 'a'",
			wantStackTop: value.False,
		},
		"c'b' <= 'ba'": {
			source:       "c'b' <= 'ba'",
			wantStackTop: value.True,
		},
		"c'a' <= 'b'": {
			source:       "c'a' <= 'b'",
			wantStackTop: value.True,
		},

		"c'2' <= 2.0": {
			source: "c'2' <= 2.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Char`",
			),
		},
		"c'i' <= 25.2bf": {
			source: "c'i' <= 25.2bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Char`",
			),
		},
		"c'f' <= 12.9f64": {
			source: "c'f' <= 12.9f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Char`",
			),
		},
		"c'0' <= 12.9f32": {
			source: "c'0' <= 12.9f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Char`",
			),
		},
		"c'9' <= 19i64": {
			source: "c'9' <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Char`",
			),
		},
		"c'u' <= 19i32": {
			source: "c'u' <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Char`",
			),
		},
		"c'4' <= 19i16": {
			source: "c'4' <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Char`",
			),
		},
		"c'6' <= 19i8": {
			source: "c'6' <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Char`",
			),
		},
		"c'9' <= 19u64": {
			source: "c'9' <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Char`",
			),
		},
		"c'u' <= 19u32": {
			source: "c'u' <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Char`",
			),
		},
		"c'4' <= 19u16": {
			source: "c'4' <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Char`",
			),
		},
		"c'6' <= 19u8": {
			source: "c'6' <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Char`",
			),
		},

		// Int
		"25 <= 25": {
			source:       "25 <= 25",
			wantStackTop: value.True,
		},
		"25 <= -25": {
			source:       "25 <= -25",
			wantStackTop: value.False,
		},
		"-25 <= 25": {
			source:       "-25 <= 25",
			wantStackTop: value.True,
		},
		"13 <= 7": {
			source:       "13 <= 7",
			wantStackTop: value.False,
		},
		"7 <= 13": {
			source:       "7 <= 13",
			wantStackTop: value.True,
		},

		"25 <= 25.0": {
			source:       "25 <= 25.0",
			wantStackTop: value.True,
		},
		"25 <= -25.0": {
			source:       "25 <= -25.0",
			wantStackTop: value.False,
		},
		"-25 <= 25.0": {
			source:       "-25 <= 25.0",
			wantStackTop: value.True,
		},
		"13 <= 7.0": {
			source:       "13 <= 7.0",
			wantStackTop: value.False,
		},
		"7 <= 13.0": {
			source:       "7 <= 13.0",
			wantStackTop: value.True,
		},
		"7 <= 7.5": {
			source:       "7 <= 7.5",
			wantStackTop: value.True,
		},
		"7 <= 6.9": {
			source:       "7 <= 6.9",
			wantStackTop: value.False,
		},

		"25 <= 25bf": {
			source:       "25 <= 25bf",
			wantStackTop: value.True,
		},
		"25 <= -25bf": {
			source:       "25 <= -25bf",
			wantStackTop: value.False,
		},
		"-25 <= 25bf": {
			source:       "-25 <= 25bf",
			wantStackTop: value.True,
		},
		"13 <= 7bf": {
			source:       "13 <= 7bf",
			wantStackTop: value.False,
		},
		"7 <= 13bf": {
			source:       "7 <= 13bf",
			wantStackTop: value.True,
		},
		"7 <= 7.5bf": {
			source:       "7 <= 7.5bf",
			wantStackTop: value.True,
		},
		"7 <= 6.9bf": {
			source:       "7 <= 6.9bf",
			wantStackTop: value.False,
		},

		"6 <= 19f64": {
			source: "6 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 <= 19f32": {
			source: "6 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 <= 19i64": {
			source: "6 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 <= 19i32": {
			source: "6 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 <= 19i16": {
			source: "6 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 <= 19i8": {
			source: "6 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 <= 19u64": {
			source: "6 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 <= 19u32": {
			source: "6 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 <= 19u16": {
			source: "6 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::SmallInt`",
			),
		},
		"6 <= 19u8": {
			source: "6 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::SmallInt`",
			),
		},

		// Float
		"25.0 <= 25.0": {
			source:       "25.0 <= 25.0",
			wantStackTop: value.True,
		},
		"25.0 <= -25.0": {
			source:       "25.0 <= -25.0",
			wantStackTop: value.False,
		},
		"-25.0 <= 25.0": {
			source:       "-25.0 <= 25.0",
			wantStackTop: value.True,
		},
		"13.0 <= 7.0": {
			source:       "13.0 <= 7.0",
			wantStackTop: value.False,
		},
		"7.0 <= 13.0": {
			source:       "7.0 <= 13.0",
			wantStackTop: value.True,
		},
		"7.0 <= 7.5": {
			source:       "7.0 <= 7.5",
			wantStackTop: value.True,
		},
		"7.5 <= 7.0": {
			source:       "7.5 <= 7.0",
			wantStackTop: value.False,
		},
		"7.0 <= 6.9": {
			source:       "7.0 <= 6.9",
			wantStackTop: value.False,
		},

		"25.0 <= 25": {
			source:       "25.0 <= 25",
			wantStackTop: value.True,
		},
		"25.0 <= -25": {
			source:       "25.0 <= -25",
			wantStackTop: value.False,
		},
		"-25.0 <= 25": {
			source:       "-25.0 <= 25",
			wantStackTop: value.True,
		},
		"13.0 <= 7": {
			source:       "13.0 <= 7",
			wantStackTop: value.False,
		},
		"7.0 <= 13": {
			source:       "7.0 <= 13",
			wantStackTop: value.True,
		},
		"7.5 <= 7": {
			source:       "7.5 <= 7",
			wantStackTop: value.False,
		},

		"25.0 <= 25bf": {
			source:       "25.0 <= 25bf",
			wantStackTop: value.True,
		},
		"25.0 <= -25bf": {
			source:       "25.0 <= -25bf",
			wantStackTop: value.False,
		},
		"-25.0 <= 25bf": {
			source:       "-25.0 <= 25bf",
			wantStackTop: value.True,
		},
		"13.0 <= 7bf": {
			source:       "13.0 <= 7bf",
			wantStackTop: value.False,
		},
		"7.0 <= 13bf": {
			source:       "7.0 <= 13bf",
			wantStackTop: value.True,
		},
		"7.0 <= 7.5bf": {
			source:       "7.0 <= 7.5bf",
			wantStackTop: value.True,
		},
		"7.5 <= 7bf": {
			source:       "7.5 <= 7bf",
			wantStackTop: value.False,
		},
		"7.0 <= 6.9bf": {
			source:       "7.0 <= 6.9bf",
			wantStackTop: value.False,
		},

		"6.0 <= 19f64": {
			source: "6.0 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19f32": {
			source: "6.0 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19i64": {
			source: "6.0 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19i32": {
			source: "6.0 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19i16": {
			source: "6.0 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19i8": {
			source: "6.0 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19u64": {
			source: "6.0 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19u32": {
			source: "6.0 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19u16": {
			source: "6.0 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float`",
			),
		},
		"6.0 <= 19u8": {
			source: "6.0 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float`",
			),
		},

		// BigFloat
		"25bf <= 25.0": {
			source:       "25bf <= 25.0",
			wantStackTop: value.True,
		},
		"25bf <= -25.0": {
			source:       "25bf <= -25.0",
			wantStackTop: value.False,
		},
		"-25bf <= 25.0": {
			source:       "-25bf <= 25.0",
			wantStackTop: value.True,
		},
		"13bf <= 7.0": {
			source:       "13bf <= 7.0",
			wantStackTop: value.False,
		},
		"7bf <= 13.0": {
			source:       "7bf <= 13.0",
			wantStackTop: value.True,
		},
		"7bf <= 7.5": {
			source:       "7bf <= 7.5",
			wantStackTop: value.True,
		},
		"7.5bf <= 7.0": {
			source:       "7.5bf <= 7.0",
			wantStackTop: value.False,
		},
		"7bf <= 6.9": {
			source:       "7bf <= 6.9",
			wantStackTop: value.False,
		},

		"25bf <= 25": {
			source:       "25bf <= 25",
			wantStackTop: value.True,
		},
		"25bf <= -25": {
			source:       "25bf <= -25",
			wantStackTop: value.False,
		},
		"-25bf <= 25": {
			source:       "-25bf <= 25",
			wantStackTop: value.True,
		},
		"13bf <= 7": {
			source:       "13bf <= 7",
			wantStackTop: value.False,
		},
		"7bf <= 13": {
			source:       "7bf <= 13",
			wantStackTop: value.True,
		},
		"7.5bf <= 7": {
			source:       "7.5bf <= 7",
			wantStackTop: value.False,
		},

		"25bf <= 25bf": {
			source:       "25bf <= 25bf",
			wantStackTop: value.True,
		},
		"25bf <= -25bf": {
			source:       "25bf <= -25bf",
			wantStackTop: value.False,
		},
		"-25bf <= 25bf": {
			source:       "-25bf <= 25bf",
			wantStackTop: value.True,
		},
		"13bf <= 7bf": {
			source:       "13bf <= 7bf",
			wantStackTop: value.False,
		},
		"7bf <= 13bf": {
			source:       "7bf <= 13bf",
			wantStackTop: value.True,
		},
		"7bf <= 7.5bf": {
			source:       "7bf <= 7.5bf",
			wantStackTop: value.True,
		},
		"7.5bf <= 7bf": {
			source:       "7.5bf <= 7bf",
			wantStackTop: value.False,
		},
		"7bf <= 6.9bf": {
			source:       "7bf <= 6.9bf",
			wantStackTop: value.False,
		},

		"6bf <= 19f64": {
			source: "6bf <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19f32": {
			source: "6bf <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19i64": {
			source: "6bf <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19i32": {
			source: "6bf <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19i16": {
			source: "6bf <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19i8": {
			source: "6bf <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19u64": {
			source: "6bf <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19u32": {
			source: "6bf <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19u16": {
			source: "6bf <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::BigFloat`",
			),
		},
		"6bf <= 19u8": {
			source: "6bf <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::BigFloat`",
			),
		},

		// Float64
		"25f64 <= 25f64": {
			source:       "25f64 <= 25f64",
			wantStackTop: value.True,
		},
		"25f64 <= -25f64": {
			source:       "25f64 <= -25f64",
			wantStackTop: value.False,
		},
		"-25f64 <= 25f64": {
			source:       "-25f64 <= 25f64",
			wantStackTop: value.True,
		},
		"13f64 <= 7f64": {
			source:       "13f64 <= 7f64",
			wantStackTop: value.False,
		},
		"7f64 <= 13f64": {
			source:       "7f64 <= 13f64",
			wantStackTop: value.True,
		},
		"7f64 <= 7.5f64": {
			source:       "7f64 <= 7.5f64",
			wantStackTop: value.True,
		},
		"7.5f64 <= 7f64": {
			source:       "7.5f64 <= 7f64",
			wantStackTop: value.False,
		},
		"7f64 <= 6.9f64": {
			source:       "7f64 <= 6.9f64",
			wantStackTop: value.False,
		},

		"6f64 <= 19.0": {
			source: "6f64 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Float64`",
			),
		},

		"6f64 <= 19": {
			source: "6f64 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19bf": {
			source: "6f64 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19f32": {
			source: "6f64 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19i64": {
			source: "6f64 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19i32": {
			source: "6f64 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19i16": {
			source: "6f64 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19i8": {
			source: "6f64 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19u64": {
			source: "6f64 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19u32": {
			source: "6f64 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19u16": {
			source: "6f64 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float64`",
			),
		},
		"6f64 <= 19u8": {
			source: "6f64 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float64`",
			),
		},

		// Float32
		"25f32 <= 25f32": {
			source:       "25f32 <= 25f32",
			wantStackTop: value.True,
		},
		"25f32 <= -25f32": {
			source:       "25f32 <= -25f32",
			wantStackTop: value.False,
		},
		"-25f32 <= 25f32": {
			source:       "-25f32 <= 25f32",
			wantStackTop: value.True,
		},
		"13f32 <= 7f32": {
			source:       "13f32 <= 7f32",
			wantStackTop: value.False,
		},
		"7f32 <= 13f32": {
			source:       "7f32 <= 13f32",
			wantStackTop: value.True,
		},
		"7f32 <= 7.5f32": {
			source:       "7f32 <= 7.5f32",
			wantStackTop: value.True,
		},
		"7.5f32 <= 7f32": {
			source:       "7.5f32 <= 7f32",
			wantStackTop: value.False,
		},
		"7f32 <= 6.9f32": {
			source:       "7f32 <= 6.9f32",
			wantStackTop: value.False,
		},

		"6f32 <= 19.0": {
			source: "6f32 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Float32`",
			),
		},

		"6f32 <= 19": {
			source: "6f32 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19bf": {
			source: "6f32 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19f64": {
			source: "6f32 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19i64": {
			source: "6f32 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19i32": {
			source: "6f32 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19i16": {
			source: "6f32 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19i8": {
			source: "6f32 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19u64": {
			source: "6f32 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19u32": {
			source: "6f32 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19u16": {
			source: "6f32 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Float32`",
			),
		},
		"6f32 <= 19u8": {
			source: "6f32 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Float32`",
			),
		},

		// Int64
		"25i64 <= 25i64": {
			source:       "25i64 <= 25i64",
			wantStackTop: value.True,
		},
		"25i64 <= -25i64": {
			source:       "25i64 <= -25i64",
			wantStackTop: value.False,
		},
		"-25i64 <= 25i64": {
			source:       "-25i64 <= 25i64",
			wantStackTop: value.True,
		},
		"13i64 <= 7i64": {
			source:       "13i64 <= 7i64",
			wantStackTop: value.False,
		},
		"7i64 <= 13i64": {
			source:       "7i64 <= 13i64",
			wantStackTop: value.True,
		},

		"6i64 <= 19": {
			source: "6i64 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19.0": {
			source: "6i64 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19bf": {
			source: "6i64 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19f64": {
			source: "6i64 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19f32": {
			source: "6i64 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19i32": {
			source: "6i64 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19i16": {
			source: "6i64 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19i8": {
			source: "6i64 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19u64": {
			source: "6i64 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19u32": {
			source: "6i64 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19u16": {
			source: "6i64 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int64`",
			),
		},
		"6i64 <= 19u8": {
			source: "6i64 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int64`",
			),
		},

		// Int32
		"25i32 <= 25i32": {
			source:       "25i32 <= 25i32",
			wantStackTop: value.True,
		},
		"25i32 <= -25i32": {
			source:       "25i32 <= -25i32",
			wantStackTop: value.False,
		},
		"-25i32 <= 25i32": {
			source:       "-25i32 <= 25i32",
			wantStackTop: value.True,
		},
		"13i32 <= 7i32": {
			source:       "13i32 <= 7i32",
			wantStackTop: value.False,
		},
		"7i32 <= 13i32": {
			source:       "7i32 <= 13i32",
			wantStackTop: value.True,
		},

		"6i32 <= 19": {
			source: "6i32 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19.0": {
			source: "6i32 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19bf": {
			source: "6i32 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19f64": {
			source: "6i32 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19f32": {
			source: "6i32 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19i64": {
			source: "6i32 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19i16": {
			source: "6i32 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19i8": {
			source: "6i32 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19u64": {
			source: "6i32 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19u32": {
			source: "6i32 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19u16": {
			source: "6i32 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int32`",
			),
		},
		"6i32 <= 19u8": {
			source: "6i32 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int32`",
			),
		},

		// Int16
		"25i16 <= 25i16": {
			source:       "25i16 <= 25i16",
			wantStackTop: value.True,
		},
		"25i16 <= -25i16": {
			source:       "25i16 <= -25i16",
			wantStackTop: value.False,
		},
		"-25i16 <= 25i16": {
			source:       "-25i16 <= 25i16",
			wantStackTop: value.True,
		},
		"13i16 <= 7i16": {
			source:       "13i16 <= 7i16",
			wantStackTop: value.False,
		},
		"7i16 <= 13i16": {
			source:       "7i16 <= 13i16",
			wantStackTop: value.True,
		},

		"6i16 <= 19": {
			source: "6i16 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19.0": {
			source: "6i16 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19bf": {
			source: "6i16 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19f64": {
			source: "6i16 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19f32": {
			source: "6i16 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19i64": {
			source: "6i16 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19i32": {
			source: "6i16 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19i8": {
			source: "6i16 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19u64": {
			source: "6i16 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19u32": {
			source: "6i16 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19u16": {
			source: "6i16 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int16`",
			),
		},
		"6i16 <= 19u8": {
			source: "6i16 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int16`",
			),
		},

		// Int8
		"25i8 <= 25i8": {
			source:       "25i8 <= 25i8",
			wantStackTop: value.True,
		},
		"25i8 <= -25i8": {
			source:       "25i8 <= -25i8",
			wantStackTop: value.False,
		},
		"-25i8 <= 25i8": {
			source:       "-25i8 <= 25i8",
			wantStackTop: value.True,
		},
		"13i8 <= 7i8": {
			source:       "13i8 <= 7i8",
			wantStackTop: value.False,
		},
		"7i8 <= 13i8": {
			source:       "7i8 <= 13i8",
			wantStackTop: value.True,
		},

		"6i8 <= 19": {
			source: "6i8 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19.0": {
			source: "6i8 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19bf": {
			source: "6i8 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19f64": {
			source: "6i8 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19f32": {
			source: "6i8 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19i64": {
			source: "6i8 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19i32": {
			source: "6i8 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19i16": {
			source: "6i8 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19u64": {
			source: "6i8 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19u32": {
			source: "6i8 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19u16": {
			source: "6i8 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::Int8`",
			),
		},
		"6i8 <= 19u8": {
			source: "6i8 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::Int8`",
			),
		},

		// UInt64
		"25u64 <= 25u64": {
			source:       "25u64 <= 25u64",
			wantStackTop: value.True,
		},
		"13u64 <= 7u64": {
			source:       "13u64 <= 7u64",
			wantStackTop: value.False,
		},
		"7u64 <= 13u64": {
			source:       "7u64 <= 13u64",
			wantStackTop: value.True,
		},

		"6u64 <= 19": {
			source: "6u64 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19.0": {
			source: "6u64 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19bf": {
			source: "6u64 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19f64": {
			source: "6u64 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19f32": {
			source: "6u64 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19i64": {
			source: "6u64 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19i32": {
			source: "6u64 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19i16": {
			source: "6u64 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19i8": {
			source: "6u64 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19u32": {
			source: "6u64 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19u16": {
			source: "6u64 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt64`",
			),
		},
		"6u64 <= 19u8": {
			source: "6u64 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt64`",
			),
		},

		// UInt32
		"25u32 <= 25u32": {
			source:       "25u32 <= 25u32",
			wantStackTop: value.True,
		},
		"13u32 <= 7u32": {
			source:       "13u32 <= 7u32",
			wantStackTop: value.False,
		},
		"7u32 <= 13u32": {
			source:       "7u32 <= 13u32",
			wantStackTop: value.True,
		},

		"6u32 <= 19": {
			source: "6u32 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19.0": {
			source: "6u32 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19bf": {
			source: "6u32 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19f64": {
			source: "6u32 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19f32": {
			source: "6u32 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19i64": {
			source: "6u32 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19i32": {
			source: "6u32 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19i16": {
			source: "6u32 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19i8": {
			source: "6u32 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19u64": {
			source: "6u32 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19u16": {
			source: "6u32 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt32`",
			),
		},
		"6u32 <= 19u8": {
			source: "6u32 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt32`",
			),
		},

		// Int16
		"25u16 <= 25u16": {
			source:       "25u16 <= 25u16",
			wantStackTop: value.True,
		},
		"13u16 <= 7u16": {
			source:       "13u16 <= 7u16",
			wantStackTop: value.False,
		},
		"7u16 <= 13u16": {
			source:       "7u16 <= 13u16",
			wantStackTop: value.True,
		},

		"6u16 <= 19": {
			source: "6u16 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19.0": {
			source: "6u16 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19bf": {
			source: "6u16 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19f64": {
			source: "6u16 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19f32": {
			source: "6u16 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19i64": {
			source: "6u16 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19i32": {
			source: "6u16 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19i16": {
			source: "6u16 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19i8": {
			source: "6u16 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19u64": {
			source: "6u16 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19u32": {
			source: "6u16 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt16`",
			),
		},
		"6u16 <= 19u8": {
			source: "6u16 <= 19u8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt8` can't be coerced into `Std::UInt16`",
			),
		},

		// Int8
		"25u8 <= 25u8": {
			source:       "25u8 <= 25u8",
			wantStackTop: value.True,
		},
		"13u8 <= 7u8": {
			source:       "13u8 <= 7u8",
			wantStackTop: value.False,
		},
		"7u8 <= 13u8": {
			source:       "7u8 <= 13u8",
			wantStackTop: value.True,
		},

		"6u8 <= 19": {
			source: "6u8 <= 19",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19.0": {
			source: "6u8 <= 19.0",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19bf": {
			source: "6u8 <= 19bf",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::BigFloat` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19f64": {
			source: "6u8 <= 19f64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19f32": {
			source: "6u8 <= 19f32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19i64": {
			source: "6u8 <= 19i64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19i32": {
			source: "6u8 <= 19i32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19i16": {
			source: "6u8 <= 19i16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int16` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19i8": {
			source: "6u8 <= 19i8",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Int8` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19u64": {
			source: "6u8 <= 19u64",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt64` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19u32": {
			source: "6u8 <= 19u32",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt32` can't be coerced into `Std::UInt8`",
			),
		},
		"6u8 <= 19u16": {
			source: "6u8 <= 19u16",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::UInt16` can't be coerced into `Std::UInt8`",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Equal(t *testing.T) {
	tests := simpleSourceTestTable{
		// String
		"'25' == '25'":   value.True,
		"'25' == '25.0'": value.False,
		"'25' == '7'":    value.False,

		"'7' == c'7'":  value.True,
		"'a' == c'a'":  value.True,
		"'7' == c'5'":  value.False,
		"'ab' == c'a'": value.False,

		"'25' == 25.0":   value.False,
		"'13.3' == 13.3": value.False,

		"'25' == 25bf":     value.False,
		"'13.3' == 13.3bf": value.False,

		"'25' == 25f64": value.False,

		"'25' == 25f32": value.False,

		"'1' == 1i64": value.False,

		"'5' == 5i32": value.False,

		"'5' == 5i16": value.False,

		"'5' == 5i8": value.False,

		"'1' == 1u64": value.False,

		"'5' == 5u32": value.False,

		"'5' == 5u16": value.False,

		"'5' == 5u8": value.False,

		// Char
		"c'2' == '2'":   value.True,
		"c'a' == 'a'":   value.True,
		"c'a' == 'ab'":  value.False,
		"c'2' == '2.0'": value.False,

		"c'7' == c'7'": value.True,
		"c'a' == c'a'": value.True,
		"c'7' == c'5'": value.False,
		"c'a' == c'b'": value.False,

		"c'2' == 2.0": value.False,

		"c'9' == 9bf": value.False,

		"c'3' == 3f64": value.False,

		"c'7' == 7f32": value.False,

		"c'1' == 1i64": value.False,

		"c'5' == 5i32": value.False,

		"c'5' == 5i16": value.False,

		"c'5' == 5i8": value.False,

		"c'1' == 1u64": value.False,

		"c'5' == 5u32": value.False,

		"c'5' == 5u16": value.False,

		"c'5' == 5u8": value.False,

		// Int
		"25 == 25":  value.True,
		"-25 == 25": value.False,
		"25 == -25": value.False,
		"25 == 28":  value.False,
		"28 == 25":  value.False,

		"25 == '25'": value.False,

		"7 == c'7'": value.False,

		"-73 == 73.0": value.False,
		"73 == -73.0": value.False,
		"25 == 25.0":  value.True,
		"1 == 1.2":    value.False,

		"-73 == 73bf": value.False,
		"73 == -73bf": value.False,
		"25 == 25bf":  value.True,
		"1 == 1.2bf":  value.False,

		"-73 == 73f64": value.False,
		"73 == -73f64": value.False,
		"25 == 25f64":  value.True,
		"1 == 1.2f64":  value.False,

		"-73 == 73f32": value.False,
		"73 == -73f32": value.False,
		"25 == 25f32":  value.True,
		"1 == 1.2f32":  value.False,

		"1 == 1i64":   value.True,
		"4 == -4i64":  value.False,
		"-8 == 8i64":  value.False,
		"-8 == -8i64": value.True,
		"91 == 27i64": value.False,

		"5 == 5i32":  value.True,
		"4 == -4i32": value.False,
		"-8 == 8i32": value.False,
		"3 == 71i32": value.False,

		"5 == 5i16":  value.True,
		"4 == -4i16": value.False,
		"-8 == 8i16": value.False,
		"3 == 71i16": value.False,

		"5 == 5i8":  value.True,
		"4 == -4i8": value.False,
		"-8 == 8i8": value.False,
		"3 == 71i8": value.False,

		"1 == 1u64":   value.True,
		"-8 == 8u64":  value.False,
		"91 == 27u64": value.False,

		"5 == 5u32":  value.True,
		"-8 == 8u32": value.False,
		"3 == 71u32": value.False,

		"53000 == 32767u16": value.False,
		"5 == 5u16":         value.True,
		"-8 == 8u16":        value.False,
		"3 == 71u16":        value.False,

		"256 == 127u8": value.False,
		"5 == 5u8":     value.True,
		"-8 == 8u8":    value.False,
		"3 == 71u8":    value.False,

		// Int64
		"25i64 == 25":  value.True,
		"-25i64 == 25": value.False,
		"25i64 == -25": value.False,
		"25i64 == 28":  value.False,
		"28i64 == 25":  value.False,

		"25i64 == '25'": value.False,

		"7i64 == c'7'": value.False,

		"-73i64 == 73.0": value.False,
		"73i64 == -73.0": value.False,
		"25i64 == 25.0":  value.True,
		"1i64 == 1.2":    value.False,

		"-73i64 == 73bf": value.False,
		"73i64 == -73bf": value.False,
		"25i64 == 25bf":  value.True,
		"1i64 == 1.2bf":  value.False,

		"-73i64 == 73f64": value.False,
		"73i64 == -73f64": value.False,
		"25i64 == 25f64":  value.True,
		"1i64 == 1.2f64":  value.False,

		"-73i64 == 73f32": value.False,
		"73i64 == -73f32": value.False,
		"25i64 == 25f32":  value.True,
		"1i64 == 1.2f32":  value.False,

		"1i64 == 1i64":   value.True,
		"4i64 == -4i64":  value.False,
		"-8i64 == 8i64":  value.False,
		"-8i64 == -8i64": value.True,
		"91i64 == 27i64": value.False,

		"5i64 == 5i32":  value.True,
		"4i64 == -4i32": value.False,
		"-8i64 == 8i32": value.False,
		"3i64 == 71i32": value.False,

		"5i64 == 5i16":  value.True,
		"4i64 == -4i16": value.False,
		"-8i64 == 8i16": value.False,
		"3i64 == 71i16": value.False,

		"5i64 == 5i8":  value.True,
		"4i64 == -4i8": value.False,
		"-8i64 == 8i8": value.False,
		"3i64 == 71i8": value.False,

		"1i64 == 1u64":   value.True,
		"-8i64 == 8u64":  value.False,
		"91i64 == 27u64": value.False,

		"5i64 == 5u32":  value.True,
		"-8i64 == 8u32": value.False,
		"3i64 == 71u32": value.False,

		"53000i64 == 32767u16": value.False,
		"5i64 == 5u16":         value.True,
		"-8i64 == 8u16":        value.False,
		"3i64 == 71u16":        value.False,

		"256i64 == 127u8": value.False,
		"5i64 == 5u8":     value.True,
		"-8i64 == 8u8":    value.False,
		"3i64 == 71u8":    value.False,

		// Int32
		"25i32 == 25":  value.True,
		"-25i32 == 25": value.False,
		"25i32 == -25": value.False,
		"25i32 == 28":  value.False,
		"28i32 == 25":  value.False,

		"25i32 == '25'": value.False,

		"7i32 == c'7'": value.False,

		"-73i32 == 73.0": value.False,
		"73i32 == -73.0": value.False,
		"25i32 == 25.0":  value.True,
		"1i32 == 1.2":    value.False,

		"-73i32 == 73bf": value.False,
		"73i32 == -73bf": value.False,
		"25i32 == 25bf":  value.True,
		"1i32 == 1.2bf":  value.False,

		"-73i32 == 73f64": value.False,
		"73i32 == -73f64": value.False,
		"25i32 == 25f64":  value.True,
		"1i32 == 1.2f64":  value.False,

		"-73i32 == 73f32": value.False,
		"73i32 == -73f32": value.False,
		"25i32 == 25f32":  value.True,
		"1i32 == 1.2f32":  value.False,

		"1i32 == 1i64":   value.True,
		"4i32 == -4i64":  value.False,
		"-8i32 == 8i64":  value.False,
		"-8i32 == -8i64": value.True,
		"91i32 == 27i64": value.False,

		"5i32 == 5i32":  value.True,
		"4i32 == -4i32": value.False,
		"-8i32 == 8i32": value.False,
		"3i32 == 71i32": value.False,

		"5i32 == 5i16":  value.True,
		"4i32 == -4i16": value.False,
		"-8i32 == 8i16": value.False,
		"3i32 == 71i16": value.False,

		"5i32 == 5i8":  value.True,
		"4i32 == -4i8": value.False,
		"-8i32 == 8i8": value.False,
		"3i32 == 71i8": value.False,

		"1i32 == 1u64":   value.True,
		"-8i32 == 8u64":  value.False,
		"91i32 == 27u64": value.False,

		"5i32 == 5u32":  value.True,
		"-8i32 == 8u32": value.False,
		"3i32 == 71u32": value.False,

		"53000i32 == 32767u16": value.False,
		"5i32 == 5u16":         value.True,
		"-8i32 == 8u16":        value.False,
		"3i32 == 71u16":        value.False,

		"256i32 == 127u8": value.False,
		"5i32 == 5u8":     value.True,
		"-8i32 == 8u8":    value.False,
		"3i32 == 71u8":    value.False,

		// Int16
		"25i16 == 25":  value.True,
		"-25i16 == 25": value.False,
		"25i16 == -25": value.False,
		"25i16 == 28":  value.False,
		"28i16 == 25":  value.False,

		"25i16 == '25'": value.False,

		"7i16 == c'7'": value.False,

		"-73i16 == 73.0": value.False,
		"73i16 == -73.0": value.False,
		"25i16 == 25.0":  value.True,
		"1i16 == 1.2":    value.False,

		"-73i16 == 73bf": value.False,
		"73i16 == -73bf": value.False,
		"25i16 == 25bf":  value.True,
		"1i16 == 1.2bf":  value.False,

		"-73i16 == 73f64": value.False,
		"73i16 == -73f64": value.False,
		"25i16 == 25f64":  value.True,
		"1i16 == 1.2f64":  value.False,

		"-73i16 == 73f32": value.False,
		"73i16 == -73f32": value.False,
		"25i16 == 25f32":  value.True,
		"1i16 == 1.2f32":  value.False,

		"1i16 == 1i64":   value.True,
		"4i16 == -4i64":  value.False,
		"-8i16 == 8i64":  value.False,
		"-8i16 == -8i64": value.True,
		"91i16 == 27i64": value.False,

		"5i16 == 5i32":  value.True,
		"4i16 == -4i32": value.False,
		"-8i16 == 8i32": value.False,
		"3i16 == 71i32": value.False,

		"5i16 == 5i16":  value.True,
		"4i16 == -4i16": value.False,
		"-8i16 == 8i16": value.False,
		"3i16 == 71i16": value.False,

		"5i16 == 5i8":  value.True,
		"4i16 == -4i8": value.False,
		"-8i16 == 8i8": value.False,
		"3i16 == 71i8": value.False,

		"1i16 == 1u64":   value.True,
		"-8i16 == 8u64":  value.False,
		"91i16 == 27u64": value.False,

		"5i16 == 5u32":  value.True,
		"-8i16 == 8u32": value.False,
		"3i16 == 71u32": value.False,

		"5i16 == 5u16":  value.True,
		"-8i16 == 8u16": value.False,
		"3i16 == 71u16": value.False,

		"256i16 == 127u8": value.False,
		"5i16 == 5u8":     value.True,
		"-8i16 == 8u8":    value.False,
		"3i16 == 71u8":    value.False,

		// Int8
		"25i8 == 25":  value.True,
		"-25i8 == 25": value.False,
		"25i8 == -25": value.False,
		"25i8 == 28":  value.False,
		"28i8 == 25":  value.False,

		"25i8 == '25'": value.False,

		"7i8 == c'7'": value.False,

		"-73i8 == 73.0": value.False,
		"73i8 == -73.0": value.False,
		"25i8 == 25.0":  value.True,
		"1i8 == 1.2":    value.False,

		"-73i8 == 73bf": value.False,
		"73i8 == -73bf": value.False,
		"25i8 == 25bf":  value.True,
		"1i8 == 1.2bf":  value.False,

		"-73i8 == 73f64": value.False,
		"73i8 == -73f64": value.False,
		"25i8 == 25f64":  value.True,
		"1i8 == 1.2f64":  value.False,

		"-73i8 == 73f32": value.False,
		"73i8 == -73f32": value.False,
		"25i8 == 25f32":  value.True,
		"1i8 == 1.2f32":  value.False,

		"1i8 == 1i64":   value.True,
		"4i8 == -4i64":  value.False,
		"-8i8 == 8i64":  value.False,
		"-8i8 == -8i64": value.True,
		"91i8 == 27i64": value.False,

		"5i8 == 5i32":  value.True,
		"4i8 == -4i32": value.False,
		"-8i8 == 8i32": value.False,
		"3i8 == 71i32": value.False,

		"5i8 == 5i16":  value.True,
		"4i8 == -4i16": value.False,
		"-8i8 == 8i16": value.False,
		"3i8 == 71i16": value.False,

		"5i8 == 5i8":  value.True,
		"4i8 == -4i8": value.False,
		"-8i8 == 8i8": value.False,
		"3i8 == 71i8": value.False,

		"1i8 == 1u64":   value.True,
		"-8i8 == 8u64":  value.False,
		"91i8 == 27u64": value.False,

		"5i8 == 5u32":  value.True,
		"-8i8 == 8u32": value.False,
		"3i8 == 71u32": value.False,

		"5i8 == 5u16":  value.True,
		"-8i8 == 8u16": value.False,
		"3i8 == 71u16": value.False,

		"5i8 == 5u8":  value.True,
		"-8i8 == 8u8": value.False,
		"3i8 == 71u8": value.False,

		// UInt64
		"25u64 == 25":  value.True,
		"25u64 == -25": value.False,
		"25u64 == 28":  value.False,
		"28u64 == 25":  value.False,

		"25u64 == '25'": value.False,

		"7u64 == c'7'": value.False,

		"73u64 == -73.0": value.False,
		"25u64 == 25.0":  value.True,
		"1u64 == 1.2":    value.False,

		"73u64 == -73bf": value.False,
		"25u64 == 25bf":  value.True,
		"1u64 == 1.2bf":  value.False,

		"73u64 == -73f64": value.False,
		"25u64 == 25f64":  value.True,
		"1u64 == 1.2f64":  value.False,

		"73u64 == -73f32": value.False,
		"25u64 == 25f32":  value.True,
		"1u64 == 1.2f32":  value.False,

		"1u64 == 1i64":   value.True,
		"4u64 == -4i64":  value.False,
		"91u64 == 27i64": value.False,

		"5u64 == 5i32":  value.True,
		"4u64 == -4i32": value.False,
		"3u64 == 71i32": value.False,

		"5u64 == 5i16":  value.True,
		"4u64 == -4i16": value.False,
		"3u64 == 71i16": value.False,

		"5u64 == 5i8":  value.True,
		"4u64 == -4i8": value.False,
		"3u64 == 71i8": value.False,

		"1u64 == 1u64":   value.True,
		"91u64 == 27u64": value.False,

		"5u64 == 5u32":  value.True,
		"3u64 == 71u32": value.False,

		"53000u64 == 32767u16": value.False,
		"5u64 == 5u16":         value.True,
		"3u64 == 71u16":        value.False,

		"256u64 == 127u8": value.False,
		"5u64 == 5u8":     value.True,
		"3u64 == 71u8":    value.False,

		// UInt32
		"25u32 == 25":  value.True,
		"25u32 == -25": value.False,
		"25u32 == 28":  value.False,
		"28u32 == 25":  value.False,

		"25u32 == '25'": value.False,

		"7u32 == c'7'": value.False,

		"73u32 == -73.0": value.False,
		"25u32 == 25.0":  value.True,
		"1u32 == 1.2":    value.False,

		"73u32 == -73bf": value.False,
		"25u32 == 25bf":  value.True,
		"1u32 == 1.2bf":  value.False,

		"73u32 == -73f64": value.False,
		"25u32 == 25f64":  value.True,
		"1u32 == 1.2f64":  value.False,

		"73u32 == -73f32": value.False,
		"25u32 == 25f32":  value.True,
		"1u32 == 1.2f32":  value.False,

		"1u32 == 1i64":   value.True,
		"4u32 == -4i64":  value.False,
		"91u32 == 27i64": value.False,

		"5u32 == 5i32":  value.True,
		"4u32 == -4i32": value.False,
		"3u32 == 71i32": value.False,

		"5u32 == 5i16":  value.True,
		"4u32 == -4i16": value.False,
		"3u32 == 71i16": value.False,

		"5u32 == 5i8":  value.True,
		"4u32 == -4i8": value.False,
		"3u32 == 71i8": value.False,

		"1u32 == 1u64":   value.True,
		"91u32 == 27u64": value.False,

		"5u32 == 5u32":  value.True,
		"3u32 == 71u32": value.False,

		"53000u32 == 32767u16": value.False,
		"5u32 == 5u16":         value.True,
		"3u32 == 71u16":        value.False,

		"256u32 == 127u8": value.False,
		"5u32 == 5u8":     value.True,
		"3u32 == 71u8":    value.False,

		// UInt16
		"25u16 == 25":  value.True,
		"25u16 == -25": value.False,
		"25u16 == 28":  value.False,
		"28u16 == 25":  value.False,

		"25u16 == '25'": value.False,

		"7u16 == c'7'": value.False,

		"73u16 == -73.0": value.False,
		"25u16 == 25.0":  value.True,
		"1u16 == 1.2":    value.False,

		"73u16 == -73bf": value.False,
		"25u16 == 25bf":  value.True,
		"1u16 == 1.2bf":  value.False,

		"73u16 == -73f64": value.False,
		"25u16 == 25f64":  value.True,
		"1u16 == 1.2f64":  value.False,

		"73u16 == -73f32": value.False,
		"25u16 == 25f32":  value.True,
		"1u16 == 1.2f32":  value.False,

		"1u16 == 1i64":   value.True,
		"4u16 == -4i64":  value.False,
		"91u16 == 27i64": value.False,

		"5u16 == 5i32":  value.True,
		"4u16 == -4i32": value.False,
		"3u16 == 71i32": value.False,

		"5u16 == 5i16":  value.True,
		"4u16 == -4i16": value.False,
		"3u16 == 71i16": value.False,

		"5u16 == 5i8":  value.True,
		"4u16 == -4i8": value.False,
		"3u16 == 71i8": value.False,

		"1u16 == 1u64":   value.True,
		"91u16 == 27u64": value.False,

		"5u16 == 5u32":  value.True,
		"3u16 == 71u32": value.False,

		"53000u16 == 32767u16": value.False,
		"5u16 == 5u16":         value.True,
		"3u16 == 71u16":        value.False,

		"256u16 == 127u8": value.False,
		"5u16 == 5u8":     value.True,
		"3u16 == 71u8":    value.False,

		// UInt8
		"25u8 == 25":  value.True,
		"25u8 == -25": value.False,
		"25u8 == 28":  value.False,
		"28u8 == 25":  value.False,

		"25u8 == '25'": value.False,

		"7u8 == c'7'": value.False,

		"73u8 == -73.0": value.False,
		"25u8 == 25.0":  value.True,
		"1u8 == 1.2":    value.False,

		"73u8 == -73bf": value.False,
		"25u8 == 25bf":  value.True,
		"1u8 == 1.2bf":  value.False,

		"73u8 == -73f64": value.False,
		"25u8 == 25f64":  value.True,
		"1u8 == 1.2f64":  value.False,

		"73u8 == -73f32": value.False,
		"25u8 == 25f32":  value.True,
		"1u8 == 1.2f32":  value.False,

		"1u8 == 1i64":   value.True,
		"4u8 == -4i64":  value.False,
		"91u8 == 27i64": value.False,

		"5u8 == 5i32":  value.True,
		"4u8 == -4i32": value.False,
		"3u8 == 71i32": value.False,

		"5u8 == 5i16":  value.True,
		"4u8 == -4i16": value.False,
		"3u8 == 71i16": value.False,

		"5u8 == 5i8":  value.True,
		"4u8 == -4i8": value.False,
		"3u8 == 71i8": value.False,

		"1u8 == 1u64":   value.True,
		"91u8 == 27u64": value.False,

		"5u8 == 5u32":  value.True,
		"3u8 == 71u32": value.False,

		"5u8 == 5u16":  value.True,
		"3u8 == 71u16": value.False,

		"5u8 == 5u8":  value.True,
		"3u8 == 71u8": value.False,

		// Float
		"-73.0 == 73.0": value.False,
		"73.0 == -73.0": value.False,
		"25.0 == 25.0":  value.True,
		"1.0 == 1.2":    value.False,
		"1.2 == 1.0":    value.False,
		"78.5 == 78.5":  value.True,

		"8.25 == '8.25'": value.False,

		"4.0 == c'4'": value.False,

		"25.0 == 25":  value.True,
		"32.3 == 32":  value.False,
		"-25.0 == 25": value.False,
		"25.0 == -25": value.False,
		"25.0 == 28":  value.False,
		"28.0 == 25":  value.False,

		"-73.0 == 73bf":  value.False,
		"73.0 == -73bf":  value.False,
		"25.0 == 25bf":   value.True,
		"1.0 == 1.2bf":   value.False,
		"15.5 == 15.5bf": value.True,

		"-73.0 == 73f64":    value.False,
		"73.0 == -73f64":    value.False,
		"25.0 == 25f64":     value.True,
		"1.0 == 1.2f64":     value.False,
		"15.26 == 15.26f64": value.True,

		"-73.0 == 73f32":  value.False,
		"73.0 == -73f32":  value.False,
		"25.0 == 25f32":   value.True,
		"1.0 == 1.2f32":   value.False,
		"15.5 == 15.5f32": value.True,

		"1.0 == 1i64":   value.True,
		"1.5 == 1i64":   value.False,
		"4.0 == -4i64":  value.False,
		"-8.0 == 8i64":  value.False,
		"-8.0 == -8i64": value.True,
		"91.0 == 27i64": value.False,

		"1.0 == 1i32":   value.True,
		"1.5 == 1i32":   value.False,
		"4.0 == -4i32":  value.False,
		"-8.0 == 8i32":  value.False,
		"-8.0 == -8i32": value.True,
		"91.0 == 27i32": value.False,

		"1.0 == 1i16":   value.True,
		"1.5 == 1i16":   value.False,
		"4.0 == -4i16":  value.False,
		"-8.0 == 8i16":  value.False,
		"-8.0 == -8i16": value.True,
		"91.0 == 27i16": value.False,

		"1.0 == 1i8":   value.True,
		"1.5 == 1i8":   value.False,
		"4.0 == -4i8":  value.False,
		"-8.0 == 8i8":  value.False,
		"-8.0 == -8i8": value.True,
		"91.0 == 27i8": value.False,

		"1.0 == 1u64":   value.True,
		"1.5 == 1u64":   value.False,
		"-8.0 == 8u64":  value.False,
		"91.0 == 27u64": value.False,

		"1.0 == 1u32":   value.True,
		"1.5 == 1u32":   value.False,
		"-8.0 == 8u32":  value.False,
		"91.0 == 27u32": value.False,

		"53000.0 == 32767u16": value.False,
		"1.0 == 1u16":         value.True,
		"1.5 == 1u16":         value.False,
		"-8.0 == 8u16":        value.False,
		"91.0 == 27u16":       value.False,

		"256.0 == 127u8": value.False,
		"1.0 == 1u8":     value.True,
		"1.5 == 1u8":     value.False,
		"-8.0 == 8u8":    value.False,
		"91.0 == 27u8":   value.False,

		// Float64
		"-73f64 == 73.0":  value.False,
		"73f64 == -73.0":  value.False,
		"25f64 == 25.0":   value.True,
		"1f64 == 1.2":     value.False,
		"1.2f64 == 1.0":   value.False,
		"78.5f64 == 78.5": value.True,

		"8.25f64 == '8.25'": value.False,

		"4f64 == c'4'": value.False,

		"25f64 == 25":   value.True,
		"32.3f64 == 32": value.False,
		"-25f64 == 25":  value.False,
		"25f64 == -25":  value.False,
		"25f64 == 28":   value.False,
		"28f64 == 25":   value.False,

		"-73f64 == 73bf":    value.False,
		"73f64 == -73bf":    value.False,
		"25f64 == 25bf":     value.True,
		"1f64 == 1.2bf":     value.False,
		"15.5f64 == 15.5bf": value.True,

		"-73f64 == 73f64":      value.False,
		"73f64 == -73f64":      value.False,
		"25f64 == 25f64":       value.True,
		"1f64 == 1.2f64":       value.False,
		"15.26f64 == 15.26f64": value.True,

		"-73f64 == 73f32":    value.False,
		"73f64 == -73f32":    value.False,
		"25f64 == 25f32":     value.True,
		"1f64 == 1.2f32":     value.False,
		"15.5f64 == 15.5f32": value.True,

		"1f64 == 1i64":   value.True,
		"1.5f64 == 1i64": value.False,
		"4f64 == -4i64":  value.False,
		"-8f64 == 8i64":  value.False,
		"-8f64 == -8i64": value.True,
		"91f64 == 27i64": value.False,

		"1f64 == 1i32":   value.True,
		"1.5f64 == 1i32": value.False,
		"4f64 == -4i32":  value.False,
		"-8f64 == 8i32":  value.False,
		"-8f64 == -8i32": value.True,
		"91f64 == 27i32": value.False,

		"1f64 == 1i16":   value.True,
		"1.5f64 == 1i16": value.False,
		"4f64 == -4i16":  value.False,
		"-8f64 == 8i16":  value.False,
		"-8f64 == -8i16": value.True,
		"91f64 == 27i16": value.False,

		"1f64 == 1i8":   value.True,
		"1.5f64 == 1i8": value.False,
		"4f64 == -4i8":  value.False,
		"-8f64 == 8i8":  value.False,
		"-8f64 == -8i8": value.True,
		"91f64 == 27i8": value.False,

		"1f64 == 1u64":   value.True,
		"1.5f64 == 1u64": value.False,
		"-8f64 == 8u64":  value.False,
		"91f64 == 27u64": value.False,

		"1f64 == 1u32":   value.True,
		"1.5f64 == 1u32": value.False,
		"-8f64 == 8u32":  value.False,
		"91f64 == 27u32": value.False,

		"53000f64 == 32767u16": value.False,
		"1f64 == 1u16":         value.True,
		"1.5f64 == 1u16":       value.False,
		"-8f64 == 8u16":        value.False,
		"91f64 == 27u16":       value.False,

		"256f64 == 127u8": value.False,
		"1f64 == 1u8":     value.True,
		"1.5f64 == 1u8":   value.False,
		"-8f64 == 8u8":    value.False,
		"91f64 == 27u8":   value.False,

		// Float32
		"-73f32 == 73.0":  value.False,
		"73f32 == -73.0":  value.False,
		"25f32 == 25.0":   value.True,
		"1f32 == 1.2":     value.False,
		"1.2f32 == 1.0":   value.False,
		"78.5f32 == 78.5": value.True,

		"8.25f32 == '8.25'": value.False,

		"4f32 == c'4'": value.False,

		"25f32 == 25":   value.True,
		"32.3f32 == 32": value.False,
		"-25f32 == 25":  value.False,
		"25f32 == -25":  value.False,
		"25f32 == 28":   value.False,
		"28f32 == 25":   value.False,

		"-73f32 == 73bf":    value.False,
		"73f32 == -73bf":    value.False,
		"25f32 == 25bf":     value.True,
		"1f32 == 1.2bf":     value.False,
		"15.5f32 == 15.5bf": value.True,

		"-73f32 == 73f64":    value.False,
		"73f32 == -73f64":    value.False,
		"25f32 == 25f64":     value.True,
		"1f32 == 1.2f64":     value.False,
		"15.5f32 == 15.5f64": value.True,

		"-73f32 == 73f32":    value.False,
		"73f32 == -73f32":    value.False,
		"25f32 == 25f32":     value.True,
		"1f32 == 1.2f32":     value.False,
		"15.5f32 == 15.5f32": value.True,

		"1f32 == 1i64":   value.True,
		"1.5f32 == 1i64": value.False,
		"4f32 == -4i64":  value.False,
		"-8f32 == 8i64":  value.False,
		"-8f32 == -8i64": value.True,
		"91f32 == 27i64": value.False,

		"1f32 == 1i32":   value.True,
		"1.5f32 == 1i32": value.False,
		"4f32 == -4i32":  value.False,
		"-8f32 == 8i32":  value.False,
		"-8f32 == -8i32": value.True,
		"91f32 == 27i32": value.False,

		"1f32 == 1i16":   value.True,
		"1.5f32 == 1i16": value.False,
		"4f32 == -4i16":  value.False,
		"-8f32 == 8i16":  value.False,
		"-8f32 == -8i16": value.True,
		"91f32 == 27i16": value.False,

		"1f32 == 1i8":   value.True,
		"1.5f32 == 1i8": value.False,
		"4f32 == -4i8":  value.False,
		"-8f32 == 8i8":  value.False,
		"-8f32 == -8i8": value.True,
		"91f32 == 27i8": value.False,

		"1f32 == 1u64":   value.True,
		"1.5f32 == 1u64": value.False,
		"-8f32 == 8u64":  value.False,
		"91f32 == 27u64": value.False,

		"1f32 == 1u32":   value.True,
		"1.5f32 == 1u32": value.False,
		"-8f32 == 8u32":  value.False,
		"91f32 == 27u32": value.False,

		"53000f32 == 32767u16": value.False,
		"1f32 == 1u16":         value.True,
		"1.5f32 == 1u16":       value.False,
		"-8f32 == 8u16":        value.False,
		"91f32 == 27u16":       value.False,

		"256f32 == 127u8": value.False,
		"1f32 == 1u8":     value.True,
		"1.5f32 == 1u8":   value.False,
		"-8f32 == 8u8":    value.False,
		"91f32 == 27u8":   value.False,
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}

func TestVMSource_NotEqual(t *testing.T) {
	tests := simpleSourceTestTable{
		// String
		"'25' != '25'":   value.False,
		"'25' != '25.0'": value.True,
		"'25' != '7'":    value.True,

		"'7' != c'7'":  value.False,
		"'a' != c'a'":  value.False,
		"'7' != c'5'":  value.True,
		"'ab' != c'a'": value.True,

		"'25' != 25.0":   value.True,
		"'13.3' != 13.3": value.True,

		"'25' != 25bf":     value.True,
		"'13.3' != 13.3bf": value.True,

		"'25' != 25f64": value.True,

		"'25' != 25f32": value.True,

		"'1' != 1i64": value.True,

		"'5' != 5i32": value.True,

		"'5' != 5i16": value.True,

		"'5' != 5i8": value.True,

		"'1' != 1u64": value.True,

		"'5' != 5u32": value.True,

		"'5' != 5u16": value.True,

		"'5' != 5u8": value.True,

		// Char
		"c'2' != '2'":   value.False,
		"c'a' != 'a'":   value.False,
		"c'a' != 'ab'":  value.True,
		"c'2' != '2.0'": value.True,

		"c'7' != c'7'": value.False,
		"c'a' != c'a'": value.False,
		"c'7' != c'5'": value.True,
		"c'a' != c'b'": value.True,

		"c'2' != 2.0": value.True,

		"c'9' != 9bf": value.True,

		"c'3' != 3f64": value.True,

		"c'7' != 7f32": value.True,

		"c'1' != 1i64": value.True,

		"c'5' != 5i32": value.True,

		"c'5' != 5i16": value.True,

		"c'5' != 5i8": value.True,

		"c'1' != 1u64": value.True,

		"c'5' != 5u32": value.True,

		"c'5' != 5u16": value.True,

		"c'5' != 5u8": value.True,

		// Int
		"25 != 25":  value.False,
		"-25 != 25": value.True,
		"25 != -25": value.True,
		"25 != 28":  value.True,
		"28 != 25":  value.True,

		"25 != '25'": value.True,

		"7 != c'7'": value.True,

		"-73 != 73.0": value.True,
		"73 != -73.0": value.True,
		"25 != 25.0":  value.False,
		"1 != 1.2":    value.True,

		"-73 != 73bf": value.True,
		"73 != -73bf": value.True,
		"25 != 25bf":  value.False,
		"1 != 1.2bf":  value.True,

		"-73 != 73f64": value.True,
		"73 != -73f64": value.True,
		"25 != 25f64":  value.False,
		"1 != 1.2f64":  value.True,

		"-73 != 73f32": value.True,
		"73 != -73f32": value.True,
		"25 != 25f32":  value.False,
		"1 != 1.2f32":  value.True,

		"1 != 1i64":   value.False,
		"4 != -4i64":  value.True,
		"-8 != 8i64":  value.True,
		"-8 != -8i64": value.False,
		"91 != 27i64": value.True,

		"5 != 5i32":  value.False,
		"4 != -4i32": value.True,
		"-8 != 8i32": value.True,
		"3 != 71i32": value.True,

		"5 != 5i16":  value.False,
		"4 != -4i16": value.True,
		"-8 != 8i16": value.True,
		"3 != 71i16": value.True,

		"5 != 5i8":  value.False,
		"4 != -4i8": value.True,
		"-8 != 8i8": value.True,
		"3 != 71i8": value.True,

		"1 != 1u64":   value.False,
		"-8 != 8u64":  value.True,
		"91 != 27u64": value.True,

		"5 != 5u32":  value.False,
		"-8 != 8u32": value.True,
		"3 != 71u32": value.True,

		"53000 != 32767u16": value.True,
		"5 != 5u16":         value.False,
		"-8 != 8u16":        value.True,
		"3 != 71u16":        value.True,

		"256 != 127u8": value.True,
		"5 != 5u8":     value.False,
		"-8 != 8u8":    value.True,
		"3 != 71u8":    value.True,

		// Int64
		"25i64 != 25":  value.False,
		"-25i64 != 25": value.True,
		"25i64 != -25": value.True,
		"25i64 != 28":  value.True,
		"28i64 != 25":  value.True,

		"25i64 != '25'": value.True,

		"7i64 != c'7'": value.True,

		"-73i64 != 73.0": value.True,
		"73i64 != -73.0": value.True,
		"25i64 != 25.0":  value.False,
		"1i64 != 1.2":    value.True,

		"-73i64 != 73bf": value.True,
		"73i64 != -73bf": value.True,
		"25i64 != 25bf":  value.False,
		"1i64 != 1.2bf":  value.True,

		"-73i64 != 73f64": value.True,
		"73i64 != -73f64": value.True,
		"25i64 != 25f64":  value.False,
		"1i64 != 1.2f64":  value.True,

		"-73i64 != 73f32": value.True,
		"73i64 != -73f32": value.True,
		"25i64 != 25f32":  value.False,
		"1i64 != 1.2f32":  value.True,

		"1i64 != 1i64":   value.False,
		"4i64 != -4i64":  value.True,
		"-8i64 != 8i64":  value.True,
		"-8i64 != -8i64": value.False,
		"91i64 != 27i64": value.True,

		"5i64 != 5i32":  value.False,
		"4i64 != -4i32": value.True,
		"-8i64 != 8i32": value.True,
		"3i64 != 71i32": value.True,

		"5i64 != 5i16":  value.False,
		"4i64 != -4i16": value.True,
		"-8i64 != 8i16": value.True,
		"3i64 != 71i16": value.True,

		"5i64 != 5i8":  value.False,
		"4i64 != -4i8": value.True,
		"-8i64 != 8i8": value.True,
		"3i64 != 71i8": value.True,

		"1i64 != 1u64":   value.False,
		"-8i64 != 8u64":  value.True,
		"91i64 != 27u64": value.True,

		"5i64 != 5u32":  value.False,
		"-8i64 != 8u32": value.True,
		"3i64 != 71u32": value.True,

		"53000i64 != 32767u16": value.True,
		"5i64 != 5u16":         value.False,
		"-8i64 != 8u16":        value.True,
		"3i64 != 71u16":        value.True,

		"256i64 != 127u8": value.True,
		"5i64 != 5u8":     value.False,
		"-8i64 != 8u8":    value.True,
		"3i64 != 71u8":    value.True,

		// Int32
		"25i32 != 25":  value.False,
		"-25i32 != 25": value.True,
		"25i32 != -25": value.True,
		"25i32 != 28":  value.True,
		"28i32 != 25":  value.True,

		"25i32 != '25'": value.True,

		"7i32 != c'7'": value.True,

		"-73i32 != 73.0": value.True,
		"73i32 != -73.0": value.True,
		"25i32 != 25.0":  value.False,
		"1i32 != 1.2":    value.True,

		"-73i32 != 73bf": value.True,
		"73i32 != -73bf": value.True,
		"25i32 != 25bf":  value.False,
		"1i32 != 1.2bf":  value.True,

		"-73i32 != 73f64": value.True,
		"73i32 != -73f64": value.True,
		"25i32 != 25f64":  value.False,
		"1i32 != 1.2f64":  value.True,

		"-73i32 != 73f32": value.True,
		"73i32 != -73f32": value.True,
		"25i32 != 25f32":  value.False,
		"1i32 != 1.2f32":  value.True,

		"1i32 != 1i64":   value.False,
		"4i32 != -4i64":  value.True,
		"-8i32 != 8i64":  value.True,
		"-8i32 != -8i64": value.False,
		"91i32 != 27i64": value.True,

		"5i32 != 5i32":  value.False,
		"4i32 != -4i32": value.True,
		"-8i32 != 8i32": value.True,
		"3i32 != 71i32": value.True,

		"5i32 != 5i16":  value.False,
		"4i32 != -4i16": value.True,
		"-8i32 != 8i16": value.True,
		"3i32 != 71i16": value.True,

		"5i32 != 5i8":  value.False,
		"4i32 != -4i8": value.True,
		"-8i32 != 8i8": value.True,
		"3i32 != 71i8": value.True,

		"1i32 != 1u64":   value.False,
		"-8i32 != 8u64":  value.True,
		"91i32 != 27u64": value.True,

		"5i32 != 5u32":  value.False,
		"-8i32 != 8u32": value.True,
		"3i32 != 71u32": value.True,

		"53000i32 != 32767u16": value.True,
		"5i32 != 5u16":         value.False,
		"-8i32 != 8u16":        value.True,
		"3i32 != 71u16":        value.True,

		"256i32 != 127u8": value.True,
		"5i32 != 5u8":     value.False,
		"-8i32 != 8u8":    value.True,
		"3i32 != 71u8":    value.True,

		// Int16
		"25i16 != 25":  value.False,
		"-25i16 != 25": value.True,
		"25i16 != -25": value.True,
		"25i16 != 28":  value.True,
		"28i16 != 25":  value.True,

		"25i16 != '25'": value.True,

		"7i16 != c'7'": value.True,

		"-73i16 != 73.0": value.True,
		"73i16 != -73.0": value.True,
		"25i16 != 25.0":  value.False,
		"1i16 != 1.2":    value.True,

		"-73i16 != 73bf": value.True,
		"73i16 != -73bf": value.True,
		"25i16 != 25bf":  value.False,
		"1i16 != 1.2bf":  value.True,

		"-73i16 != 73f64": value.True,
		"73i16 != -73f64": value.True,
		"25i16 != 25f64":  value.False,
		"1i16 != 1.2f64":  value.True,

		"-73i16 != 73f32": value.True,
		"73i16 != -73f32": value.True,
		"25i16 != 25f32":  value.False,
		"1i16 != 1.2f32":  value.True,

		"1i16 != 1i64":   value.False,
		"4i16 != -4i64":  value.True,
		"-8i16 != 8i64":  value.True,
		"-8i16 != -8i64": value.False,
		"91i16 != 27i64": value.True,

		"5i16 != 5i32":  value.False,
		"4i16 != -4i32": value.True,
		"-8i16 != 8i32": value.True,
		"3i16 != 71i32": value.True,

		"5i16 != 5i16":  value.False,
		"4i16 != -4i16": value.True,
		"-8i16 != 8i16": value.True,
		"3i16 != 71i16": value.True,

		"5i16 != 5i8":  value.False,
		"4i16 != -4i8": value.True,
		"-8i16 != 8i8": value.True,
		"3i16 != 71i8": value.True,

		"1i16 != 1u64":   value.False,
		"-8i16 != 8u64":  value.True,
		"91i16 != 27u64": value.True,

		"5i16 != 5u32":  value.False,
		"-8i16 != 8u32": value.True,
		"3i16 != 71u32": value.True,

		"5i16 != 5u16":  value.False,
		"-8i16 != 8u16": value.True,
		"3i16 != 71u16": value.True,

		"256i16 != 127u8": value.True,
		"5i16 != 5u8":     value.False,
		"-8i16 != 8u8":    value.True,
		"3i16 != 71u8":    value.True,

		// Int8
		"25i8 != 25":  value.False,
		"-25i8 != 25": value.True,
		"25i8 != -25": value.True,
		"25i8 != 28":  value.True,
		"28i8 != 25":  value.True,

		"25i8 != '25'": value.True,

		"7i8 != c'7'": value.True,

		"-73i8 != 73.0": value.True,
		"73i8 != -73.0": value.True,
		"25i8 != 25.0":  value.False,
		"1i8 != 1.2":    value.True,

		"-73i8 != 73bf": value.True,
		"73i8 != -73bf": value.True,
		"25i8 != 25bf":  value.False,
		"1i8 != 1.2bf":  value.True,

		"-73i8 != 73f64": value.True,
		"73i8 != -73f64": value.True,
		"25i8 != 25f64":  value.False,
		"1i8 != 1.2f64":  value.True,

		"-73i8 != 73f32": value.True,
		"73i8 != -73f32": value.True,
		"25i8 != 25f32":  value.False,
		"1i8 != 1.2f32":  value.True,

		"1i8 != 1i64":   value.False,
		"4i8 != -4i64":  value.True,
		"-8i8 != 8i64":  value.True,
		"-8i8 != -8i64": value.False,
		"91i8 != 27i64": value.True,

		"5i8 != 5i32":  value.False,
		"4i8 != -4i32": value.True,
		"-8i8 != 8i32": value.True,
		"3i8 != 71i32": value.True,

		"5i8 != 5i16":  value.False,
		"4i8 != -4i16": value.True,
		"-8i8 != 8i16": value.True,
		"3i8 != 71i16": value.True,

		"5i8 != 5i8":  value.False,
		"4i8 != -4i8": value.True,
		"-8i8 != 8i8": value.True,
		"3i8 != 71i8": value.True,

		"1i8 != 1u64":   value.False,
		"-8i8 != 8u64":  value.True,
		"91i8 != 27u64": value.True,

		"5i8 != 5u32":  value.False,
		"-8i8 != 8u32": value.True,
		"3i8 != 71u32": value.True,

		"5i8 != 5u16":  value.False,
		"-8i8 != 8u16": value.True,
		"3i8 != 71u16": value.True,

		"5i8 != 5u8":  value.False,
		"-8i8 != 8u8": value.True,
		"3i8 != 71u8": value.True,

		// UInt64
		"25u64 != 25":  value.False,
		"25u64 != -25": value.True,
		"25u64 != 28":  value.True,
		"28u64 != 25":  value.True,

		"25u64 != '25'": value.True,

		"7u64 != c'7'": value.True,

		"73u64 != -73.0": value.True,
		"25u64 != 25.0":  value.False,
		"1u64 != 1.2":    value.True,

		"73u64 != -73bf": value.True,
		"25u64 != 25bf":  value.False,
		"1u64 != 1.2bf":  value.True,

		"73u64 != -73f64": value.True,
		"25u64 != 25f64":  value.False,
		"1u64 != 1.2f64":  value.True,

		"73u64 != -73f32": value.True,
		"25u64 != 25f32":  value.False,
		"1u64 != 1.2f32":  value.True,

		"1u64 != 1i64":   value.False,
		"4u64 != -4i64":  value.True,
		"91u64 != 27i64": value.True,

		"5u64 != 5i32":  value.False,
		"4u64 != -4i32": value.True,
		"3u64 != 71i32": value.True,

		"5u64 != 5i16":  value.False,
		"4u64 != -4i16": value.True,
		"3u64 != 71i16": value.True,

		"5u64 != 5i8":  value.False,
		"4u64 != -4i8": value.True,
		"3u64 != 71i8": value.True,

		"1u64 != 1u64":   value.False,
		"91u64 != 27u64": value.True,

		"5u64 != 5u32":  value.False,
		"3u64 != 71u32": value.True,

		"53000u64 != 32767u16": value.True,
		"5u64 != 5u16":         value.False,
		"3u64 != 71u16":        value.True,

		"256u64 != 127u8": value.True,
		"5u64 != 5u8":     value.False,
		"3u64 != 71u8":    value.True,

		// UInt32
		"25u32 != 25":  value.False,
		"25u32 != -25": value.True,
		"25u32 != 28":  value.True,
		"28u32 != 25":  value.True,

		"25u32 != '25'": value.True,

		"7u32 != c'7'": value.True,

		"73u32 != -73.0": value.True,
		"25u32 != 25.0":  value.False,
		"1u32 != 1.2":    value.True,

		"73u32 != -73bf": value.True,
		"25u32 != 25bf":  value.False,
		"1u32 != 1.2bf":  value.True,

		"73u32 != -73f64": value.True,
		"25u32 != 25f64":  value.False,
		"1u32 != 1.2f64":  value.True,

		"73u32 != -73f32": value.True,
		"25u32 != 25f32":  value.False,
		"1u32 != 1.2f32":  value.True,

		"1u32 != 1i64":   value.False,
		"4u32 != -4i64":  value.True,
		"91u32 != 27i64": value.True,

		"5u32 != 5i32":  value.False,
		"4u32 != -4i32": value.True,
		"3u32 != 71i32": value.True,

		"5u32 != 5i16":  value.False,
		"4u32 != -4i16": value.True,
		"3u32 != 71i16": value.True,

		"5u32 != 5i8":  value.False,
		"4u32 != -4i8": value.True,
		"3u32 != 71i8": value.True,

		"1u32 != 1u64":   value.False,
		"91u32 != 27u64": value.True,

		"5u32 != 5u32":  value.False,
		"3u32 != 71u32": value.True,

		"53000u32 != 32767u16": value.True,
		"5u32 != 5u16":         value.False,
		"3u32 != 71u16":        value.True,

		"256u32 != 127u8": value.True,
		"5u32 != 5u8":     value.False,
		"3u32 != 71u8":    value.True,

		// UInt16
		"25u16 != 25":  value.False,
		"25u16 != -25": value.True,
		"25u16 != 28":  value.True,
		"28u16 != 25":  value.True,

		"25u16 != '25'": value.True,

		"7u16 != c'7'": value.True,

		"73u16 != -73.0": value.True,
		"25u16 != 25.0":  value.False,
		"1u16 != 1.2":    value.True,

		"73u16 != -73bf": value.True,
		"25u16 != 25bf":  value.False,
		"1u16 != 1.2bf":  value.True,

		"73u16 != -73f64": value.True,
		"25u16 != 25f64":  value.False,
		"1u16 != 1.2f64":  value.True,

		"73u16 != -73f32": value.True,
		"25u16 != 25f32":  value.False,
		"1u16 != 1.2f32":  value.True,

		"1u16 != 1i64":   value.False,
		"4u16 != -4i64":  value.True,
		"91u16 != 27i64": value.True,

		"5u16 != 5i32":  value.False,
		"4u16 != -4i32": value.True,
		"3u16 != 71i32": value.True,

		"5u16 != 5i16":  value.False,
		"4u16 != -4i16": value.True,
		"3u16 != 71i16": value.True,

		"5u16 != 5i8":  value.False,
		"4u16 != -4i8": value.True,
		"3u16 != 71i8": value.True,

		"1u16 != 1u64":   value.False,
		"91u16 != 27u64": value.True,

		"5u16 != 5u32":  value.False,
		"3u16 != 71u32": value.True,

		"53000u16 != 32767u16": value.True,
		"5u16 != 5u16":         value.False,
		"3u16 != 71u16":        value.True,

		"256u16 != 127u8": value.True,
		"5u16 != 5u8":     value.False,
		"3u16 != 71u8":    value.True,

		// UInt8
		"25u8 != 25":  value.False,
		"25u8 != -25": value.True,
		"25u8 != 28":  value.True,
		"28u8 != 25":  value.True,

		"25u8 != '25'": value.True,

		"7u8 != c'7'": value.True,

		"73u8 != -73.0": value.True,
		"25u8 != 25.0":  value.False,
		"1u8 != 1.2":    value.True,

		"73u8 != -73bf": value.True,
		"25u8 != 25bf":  value.False,
		"1u8 != 1.2bf":  value.True,

		"73u8 != -73f64": value.True,
		"25u8 != 25f64":  value.False,
		"1u8 != 1.2f64":  value.True,

		"73u8 != -73f32": value.True,
		"25u8 != 25f32":  value.False,
		"1u8 != 1.2f32":  value.True,

		"1u8 != 1i64":   value.False,
		"4u8 != -4i64":  value.True,
		"91u8 != 27i64": value.True,

		"5u8 != 5i32":  value.False,
		"4u8 != -4i32": value.True,
		"3u8 != 71i32": value.True,

		"5u8 != 5i16":  value.False,
		"4u8 != -4i16": value.True,
		"3u8 != 71i16": value.True,

		"5u8 != 5i8":  value.False,
		"4u8 != -4i8": value.True,
		"3u8 != 71i8": value.True,

		"1u8 != 1u64":   value.False,
		"91u8 != 27u64": value.True,

		"5u8 != 5u32":  value.False,
		"3u8 != 71u32": value.True,

		"5u8 != 5u16":  value.False,
		"3u8 != 71u16": value.True,

		"5u8 != 5u8":  value.False,
		"3u8 != 71u8": value.True,

		// Float
		"-73.0 != 73.0": value.True,
		"73.0 != -73.0": value.True,
		"25.0 != 25.0":  value.False,
		"1.0 != 1.2":    value.True,
		"1.2 != 1.0":    value.True,
		"78.5 != 78.5":  value.False,

		"8.25 != '8.25'": value.True,

		"4.0 != c'4'": value.True,

		"25.0 != 25":  value.False,
		"32.3 != 32":  value.True,
		"-25.0 != 25": value.True,
		"25.0 != -25": value.True,
		"25.0 != 28":  value.True,
		"28.0 != 25":  value.True,

		"-73.0 != 73bf":  value.True,
		"73.0 != -73bf":  value.True,
		"25.0 != 25bf":   value.False,
		"1.0 != 1.2bf":   value.True,
		"15.5 != 15.5bf": value.False,

		"-73.0 != 73f64":    value.True,
		"73.0 != -73f64":    value.True,
		"25.0 != 25f64":     value.False,
		"1.0 != 1.2f64":     value.True,
		"15.26 != 15.26f64": value.False,

		"-73.0 != 73f32":  value.True,
		"73.0 != -73f32":  value.True,
		"25.0 != 25f32":   value.False,
		"1.0 != 1.2f32":   value.True,
		"15.5 != 15.5f32": value.False,

		"1.0 != 1i64":   value.False,
		"1.5 != 1i64":   value.True,
		"4.0 != -4i64":  value.True,
		"-8.0 != 8i64":  value.True,
		"-8.0 != -8i64": value.False,
		"91.0 != 27i64": value.True,

		"1.0 != 1i32":   value.False,
		"1.5 != 1i32":   value.True,
		"4.0 != -4i32":  value.True,
		"-8.0 != 8i32":  value.True,
		"-8.0 != -8i32": value.False,
		"91.0 != 27i32": value.True,

		"1.0 != 1i16":   value.False,
		"1.5 != 1i16":   value.True,
		"4.0 != -4i16":  value.True,
		"-8.0 != 8i16":  value.True,
		"-8.0 != -8i16": value.False,
		"91.0 != 27i16": value.True,

		"1.0 != 1i8":   value.False,
		"1.5 != 1i8":   value.True,
		"4.0 != -4i8":  value.True,
		"-8.0 != 8i8":  value.True,
		"-8.0 != -8i8": value.False,
		"91.0 != 27i8": value.True,

		"1.0 != 1u64":   value.False,
		"1.5 != 1u64":   value.True,
		"-8.0 != 8u64":  value.True,
		"91.0 != 27u64": value.True,

		"1.0 != 1u32":   value.False,
		"1.5 != 1u32":   value.True,
		"-8.0 != 8u32":  value.True,
		"91.0 != 27u32": value.True,

		"53000.0 != 32767u16": value.True,
		"1.0 != 1u16":         value.False,
		"1.5 != 1u16":         value.True,
		"-8.0 != 8u16":        value.True,
		"91.0 != 27u16":       value.True,

		"256.0 != 127u8": value.True,
		"1.0 != 1u8":     value.False,
		"1.5 != 1u8":     value.True,
		"-8.0 != 8u8":    value.True,
		"91.0 != 27u8":   value.True,

		// Float64
		"-73f64 != 73.0":  value.True,
		"73f64 != -73.0":  value.True,
		"25f64 != 25.0":   value.False,
		"1f64 != 1.2":     value.True,
		"1.2f64 != 1.0":   value.True,
		"78.5f64 != 78.5": value.False,

		"8.25f64 != '8.25'": value.True,

		"4f64 != c'4'": value.True,

		"25f64 != 25":   value.False,
		"32.3f64 != 32": value.True,
		"-25f64 != 25":  value.True,
		"25f64 != -25":  value.True,
		"25f64 != 28":   value.True,
		"28f64 != 25":   value.True,

		"-73f64 != 73bf":    value.True,
		"73f64 != -73bf":    value.True,
		"25f64 != 25bf":     value.False,
		"1f64 != 1.2bf":     value.True,
		"15.5f64 != 15.5bf": value.False,

		"-73f64 != 73f64":      value.True,
		"73f64 != -73f64":      value.True,
		"25f64 != 25f64":       value.False,
		"1f64 != 1.2f64":       value.True,
		"15.26f64 != 15.26f64": value.False,

		"-73f64 != 73f32":    value.True,
		"73f64 != -73f32":    value.True,
		"25f64 != 25f32":     value.False,
		"1f64 != 1.2f32":     value.True,
		"15.5f64 != 15.5f32": value.False,

		"1f64 != 1i64":   value.False,
		"1.5f64 != 1i64": value.True,
		"4f64 != -4i64":  value.True,
		"-8f64 != 8i64":  value.True,
		"-8f64 != -8i64": value.False,
		"91f64 != 27i64": value.True,

		"1f64 != 1i32":   value.False,
		"1.5f64 != 1i32": value.True,
		"4f64 != -4i32":  value.True,
		"-8f64 != 8i32":  value.True,
		"-8f64 != -8i32": value.False,
		"91f64 != 27i32": value.True,

		"1f64 != 1i16":   value.False,
		"1.5f64 != 1i16": value.True,
		"4f64 != -4i16":  value.True,
		"-8f64 != 8i16":  value.True,
		"-8f64 != -8i16": value.False,
		"91f64 != 27i16": value.True,

		"1f64 != 1i8":   value.False,
		"1.5f64 != 1i8": value.True,
		"4f64 != -4i8":  value.True,
		"-8f64 != 8i8":  value.True,
		"-8f64 != -8i8": value.False,
		"91f64 != 27i8": value.True,

		"1f64 != 1u64":   value.False,
		"1.5f64 != 1u64": value.True,
		"-8f64 != 8u64":  value.True,
		"91f64 != 27u64": value.True,

		"1f64 != 1u32":   value.False,
		"1.5f64 != 1u32": value.True,
		"-8f64 != 8u32":  value.True,
		"91f64 != 27u32": value.True,

		"53000f64 != 32767u16": value.True,
		"1f64 != 1u16":         value.False,
		"1.5f64 != 1u16":       value.True,
		"-8f64 != 8u16":        value.True,
		"91f64 != 27u16":       value.True,

		"256f64 != 127u8": value.True,
		"1f64 != 1u8":     value.False,
		"1.5f64 != 1u8":   value.True,
		"-8f64 != 8u8":    value.True,
		"91f64 != 27u8":   value.True,

		// Float32
		"-73f32 != 73.0":  value.True,
		"73f32 != -73.0":  value.True,
		"25f32 != 25.0":   value.False,
		"1f32 != 1.2":     value.True,
		"1.2f32 != 1.0":   value.True,
		"78.5f32 != 78.5": value.False,

		"8.25f32 != '8.25'": value.True,

		"4f32 != c'4'": value.True,

		"25f32 != 25":   value.False,
		"32.3f32 != 32": value.True,
		"-25f32 != 25":  value.True,
		"25f32 != -25":  value.True,
		"25f32 != 28":   value.True,
		"28f32 != 25":   value.True,

		"-73f32 != 73bf":    value.True,
		"73f32 != -73bf":    value.True,
		"25f32 != 25bf":     value.False,
		"1f32 != 1.2bf":     value.True,
		"15.5f32 != 15.5bf": value.False,

		"-73f32 != 73f64":    value.True,
		"73f32 != -73f64":    value.True,
		"25f32 != 25f64":     value.False,
		"1f32 != 1.2f64":     value.True,
		"15.5f32 != 15.5f64": value.False,

		"-73f32 != 73f32":    value.True,
		"73f32 != -73f32":    value.True,
		"25f32 != 25f32":     value.False,
		"1f32 != 1.2f32":     value.True,
		"15.5f32 != 15.5f32": value.False,

		"1f32 != 1i64":   value.False,
		"1.5f32 != 1i64": value.True,
		"4f32 != -4i64":  value.True,
		"-8f32 != 8i64":  value.True,
		"-8f32 != -8i64": value.False,
		"91f32 != 27i64": value.True,

		"1f32 != 1i32":   value.False,
		"1.5f32 != 1i32": value.True,
		"4f32 != -4i32":  value.True,
		"-8f32 != 8i32":  value.True,
		"-8f32 != -8i32": value.False,
		"91f32 != 27i32": value.True,

		"1f32 != 1i16":   value.False,
		"1.5f32 != 1i16": value.True,
		"4f32 != -4i16":  value.True,
		"-8f32 != 8i16":  value.True,
		"-8f32 != -8i16": value.False,
		"91f32 != 27i16": value.True,

		"1f32 != 1i8":   value.False,
		"1.5f32 != 1i8": value.True,
		"4f32 != -4i8":  value.True,
		"-8f32 != 8i8":  value.True,
		"-8f32 != -8i8": value.False,
		"91f32 != 27i8": value.True,

		"1f32 != 1u64":   value.False,
		"1.5f32 != 1u64": value.True,
		"-8f32 != 8u64":  value.True,
		"91f32 != 27u64": value.True,

		"1f32 != 1u32":   value.False,
		"1.5f32 != 1u32": value.True,
		"-8f32 != 8u32":  value.True,
		"91f32 != 27u32": value.True,

		"53000f32 != 32767u16": value.True,
		"1f32 != 1u16":         value.False,
		"1.5f32 != 1u16":       value.True,
		"-8f32 != 8u16":        value.True,
		"91f32 != 27u16":       value.True,

		"256f32 != 127u8": value.True,
		"1f32 != 1u8":     value.False,
		"1.5f32 != 1u8":   value.True,
		"-8f32 != 8u8":    value.True,
		"91f32 != 27u8":   value.True,
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}

func TestVMSource_StrictEqual(t *testing.T) {
	tests := simpleSourceTestTable{
		// String
		"'25' === '25'":   value.True,
		"'25' === '25.0'": value.False,
		"'25' === '7'":    value.False,

		"'7' === c'7'": value.False,

		"'25' === 25.0": value.False,

		"'25' === 25bf": value.False,

		"'25' === 25f64": value.False,

		"'25' === 25f32": value.False,

		"'1' === 1i64": value.False,

		"'5' === 5i32": value.False,

		"'5' === 5i16": value.False,

		"'5' === 5i8": value.False,

		"'1' === 1u64": value.False,

		"'5' === 5u32": value.False,

		"'5' === 5u16": value.False,

		"'5' === 5u8": value.False,

		// Char
		"c'2' === 25": value.False,

		"c'2' === '2'": value.False,

		"c'7' === c'7'": value.True,
		"c'b' === c'b'": value.True,
		"c'c' === c'g'": value.False,
		"c'7' === c'8'": value.False,

		"c'2' === 2.0": value.False,

		"c'3' === 3bf": value.False,

		"c'9' === 9f64": value.False,

		"c'1' === 1f32": value.False,

		"c'1' === 1i64": value.False,

		"c'5' === 5i32": value.False,

		"c'5' === 5i16": value.False,

		"c'5' === 5i8": value.False,

		"c'1' === 1u64": value.False,

		"c'5' === 5u32": value.False,

		"c'5' === 5u16": value.False,

		"c'5' === 5u8": value.False,

		// Int
		"25 === 25":  value.True,
		"-25 === 25": value.False,
		"25 === -25": value.False,
		"25 === 28":  value.False,
		"28 === 25":  value.False,

		"25 === '25'": value.False,

		"7 === c'7'": value.False,

		"25 === 25.0": value.False,

		"25 === 25bf": value.False,

		"25 === 25f64": value.False,

		"25 === 25f32": value.False,

		"1 === 1i64": value.False,

		"5 === 5i32": value.False,

		"5 === 5i16": value.False,

		"5 === 5i8": value.False,

		"1 === 1u64": value.False,

		"5 === 5u32": value.False,

		"5 === 5u16": value.False,

		"5 === 5u8": value.False,

		// Int64
		"25i64 === 25": value.False,

		"25i64 === '25'": value.False,

		"7i64 === c'7'": value.False,

		"25i64 === 25.0": value.False,

		"25i64 === 25bf": value.False,

		"25i64 === 25f64": value.False,

		"25i64 === 25f32": value.False,

		"1i64 === 1i64":   value.True,
		"4i64 === -4i64":  value.False,
		"-8i64 === 8i64":  value.False,
		"-8i64 === -8i64": value.True,
		"91i64 === 27i64": value.False,

		"5i64 === 5i32": value.False,

		"5i64 === 5i16": value.False,

		"5i64 === 5i8": value.False,

		"1i64 === 1u64": value.False,

		"5i64 === 5u32": value.False,

		"5i64 === 5u16": value.False,

		"5i64 === 5u8": value.False,

		// Int32
		"25i32 === 25": value.False,

		"25i32 === '25'": value.False,

		"7i32 === c'7'": value.False,

		"25i32 === 25.0": value.False,

		"25i32 === 25bf": value.False,

		"25i32 === 25f64": value.False,

		"25i32 === 25f32": value.False,

		"1i32 === 1i64": value.False,

		"5i32 === 5i32":  value.True,
		"4i32 === -4i32": value.False,
		"-8i32 === 8i32": value.False,
		"3i32 === 71i32": value.False,

		"5i32 === 5i16": value.False,

		"5i32 === 5i8": value.False,

		"1i32 === 1u64": value.False,

		"5i32 === 5u32": value.False,

		"5i32 === 5u16": value.False,

		"5i32 === 5u8": value.False,

		// Int16
		"25i16 === 25": value.False,

		"25i16 === '25'": value.False,

		"7i16 === c'7'": value.False,

		"25i16 === 25.0": value.False,

		"25i16 === 25bf": value.False,

		"25i16 === 25f64": value.False,

		"25i16 === 25f32": value.False,

		"1i16 === 1i64": value.False,

		"5i16 === 5i32": value.False,

		"5i16 === 5i16":  value.True,
		"4i16 === -4i16": value.False,
		"-8i16 === 8i16": value.False,
		"3i16 === 71i16": value.False,

		"5i16 === 5i8": value.False,

		"1i16 === 1u64": value.False,

		"5i16 === 5u32": value.False,

		"5i16 === 5u16": value.False,

		"5i16 === 5u8": value.False,

		// Int8
		"25i8 === 25": value.False,

		"25i8 === '25'": value.False,

		"7i8 === c'7'": value.False,

		"25i8 === 25.0": value.False,

		"25i8 === 25bf": value.False,

		"25i8 === 25f64": value.False,

		"25i8 === 25f32": value.False,

		"1i8 === 1i64": value.False,

		"5i8 === 5i32": value.False,

		"5i8 === 5i16": value.False,

		"5i8 === 5i8":  value.True,
		"4i8 === -4i8": value.False,
		"-8i8 === 8i8": value.False,
		"3i8 === 71i8": value.False,

		"1i8 === 1u64": value.False,

		"5i8 === 5u32": value.False,

		"5i8 === 5u16": value.False,

		"5i8 === 5u8": value.False,

		// UInt64
		"25u64 === 25": value.False,

		"25u64 === '25'": value.False,

		"7u64 === c'7'": value.False,

		"25u64 === 25.0": value.False,

		"25u64 === 25bf": value.False,

		"25u64 === 25f64": value.False,

		"25u64 === 25f32": value.False,

		"1u64 === 1i64": value.False,

		"5u64 === 5i32": value.False,

		"5u64 === 5i16": value.False,

		"5u64 === 5i8": value.False,

		"1u64 === 1u64":   value.True,
		"91u64 === 27u64": value.False,

		"5u64 === 5u32": value.False,

		"5u64 === 5u16": value.False,

		"5u64 === 5u8": value.False,

		// UInt32
		"25u32 === 25": value.False,

		"25u32 === '25'": value.False,

		"7u32 === c'7'": value.False,

		"25u32 === 25.0": value.False,

		"25u32 === 25bf": value.False,

		"25u32 === 25f64": value.False,

		"25u32 === 25f32": value.False,

		"1u32 === 1i64": value.False,

		"5u32 === 5i32": value.False,

		"5u32 === 5i16": value.False,

		"5u32 === 5i8": value.False,

		"1u32 === 1u64": value.False,

		"5u32 === 5u32":  value.True,
		"3u32 === 71u32": value.False,

		"5u32 === 5u16": value.False,

		"5u32 === 5u8": value.False,

		// UInt16
		"25u16 === 25": value.False,

		"25u16 === '25'": value.False,

		"7u16 === c'7'": value.False,

		"25u16 === 25.0": value.False,

		"25u16 === 25bf": value.False,

		"25u16 === 25f64": value.False,

		"25u16 === 25f32": value.False,

		"1u16 === 1i64": value.False,

		"5u16 === 5i32": value.False,

		"5u16 === 5i16": value.False,

		"5u16 === 5i8": value.False,

		"1u16 === 1u64": value.False,

		"5u16 === 5u32": value.False,

		"53000u16 === 32767u16": value.False,
		"5u16 === 5u16":         value.True,
		"3u16 === 71u16":        value.False,

		"5u16 === 5u8": value.False,

		// UInt8
		"25u8 === 25": value.False,

		"25u8 === '25'": value.False,

		"7u8 === c'7'": value.False,

		"25u8 === 25.0": value.False,

		"25u8 === 25bf": value.False,

		"25u8 === 25f64": value.False,

		"25u8 === 25f32": value.False,

		"1u8 === 1i64": value.False,

		"5u8 === 5i32": value.False,

		"5u8 === 5i16": value.False,

		"5u8 === 5i8": value.False,

		"1u8 === 1u64": value.False,

		"5u8 === 5u32": value.False,

		"5u8 === 5u16": value.False,

		"5u8 === 5u8":  value.True,
		"3u8 === 71u8": value.False,

		// Float
		"-73.0 === 73.0": value.False,
		"73.0 === -73.0": value.False,
		"25.0 === 25.0":  value.True,
		"1.0 === 1.2":    value.False,
		"1.2 === 1.0":    value.False,
		"78.5 === 78.5":  value.True,

		"8.25 === '8.25'": value.False,

		"4.0 === c'4'": value.False,

		"25.0 === 25": value.False,

		"25.0 === 25bf":   value.False,
		"15.5 === 15.5bf": value.False,

		"25.0 === 25f64":     value.False,
		"15.26 === 15.26f64": value.False,

		"25.0 === 25f32":   value.False,
		"15.5 === 15.5f32": value.False,

		"1.0 === 1i64":   value.False,
		"-8.0 === -8i64": value.False,

		"1.0 === 1i32":   value.False,
		"-8.0 === -8i32": value.False,

		"1.0 === 1i16":   value.False,
		"-8.0 === -8i16": value.False,

		"1.0 === 1i8":   value.False,
		"-8.0 === -8i8": value.False,

		"1.0 === 1u64": value.False,

		"1.0 === 1u32": value.False,

		"1.0 === 1u16": value.False,

		"1.0 === 1u8": value.False,

		// Float64
		"25f64 === 25.0":   value.False,
		"78.5f64 === 78.5": value.False,

		"8.25f64 === '8.25'": value.False,

		"4f64 === c'4'": value.False,

		"25f64 === 25": value.False,

		"25f64 === 25bf":     value.False,
		"15.5f64 === 15.5bf": value.False,

		"-73f64 === 73f64":      value.False,
		"73f64 === -73f64":      value.False,
		"25f64 === 25f64":       value.True,
		"1f64 === 1.2f64":       value.False,
		"15.26f64 === 15.26f64": value.True,

		"25f64 === 25f32":     value.False,
		"15.5f64 === 15.5f32": value.False,

		"1f64 === 1i64":   value.False,
		"-8f64 === -8i64": value.False,

		"1f64 === 1i32":   value.False,
		"-8f64 === -8i32": value.False,

		"1f64 === 1i16":   value.False,
		"-8f64 === -8i16": value.False,

		"1f64 === 1i8":   value.False,
		"-8f64 === -8i8": value.False,

		"1f64 === 1u64": value.False,

		"1f64 === 1u32": value.False,

		"1f64 === 1u16": value.False,

		"1f64 === 1u8": value.False,

		// Float32
		"25f32 === 25.0":   value.False,
		"78.5f32 === 78.5": value.False,

		"8.25f32 === '8.25'": value.False,

		"4f32 === c'4'": value.False,

		"25f32 === 25": value.False,

		"25f32 === 25bf":     value.False,
		"15.5f32 === 15.5bf": value.False,

		"25f32 === 25f64":     value.False,
		"15.5f32 === 15.5f64": value.False,

		"-73f32 === 73f32":    value.False,
		"73f32 === -73f32":    value.False,
		"25f32 === 25f32":     value.True,
		"1f32 === 1.2f32":     value.False,
		"15.5f32 === 15.5f32": value.True,

		"1f32 === 1i64":   value.False,
		"-8f32 === -8i64": value.False,

		"1f32 === 1i32":   value.False,
		"-8f32 === -8i32": value.False,

		"1f32 === 1i16":   value.False,
		"-8f32 === -8i16": value.False,

		"1f32 === 1i8":   value.False,
		"-8f32 === -8i8": value.False,

		"1f32 === 1u64": value.False,

		"1f32 === 1u32": value.False,

		"1f32 === 1u16": value.False,

		"1f32 === 1u8": value.False,
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}

func TestVMSource_StrictNotEqual(t *testing.T) {
	tests := simpleSourceTestTable{
		// String
		"'25' !== '25'":   value.False,
		"'25' !== '25.0'": value.True,
		"'25' !== '7'":    value.True,

		"'7' !== c'7'": value.True,

		"'25' !== 25.0": value.True,

		"'25' !== 25bf": value.True,

		"'25' !== 25f64": value.True,

		"'25' !== 25f32": value.True,

		"'1' !== 1i64": value.True,

		"'5' !== 5i32": value.True,

		"'5' !== 5i16": value.True,

		"'5' !== 5i8": value.True,

		"'1' !== 1u64": value.True,

		"'5' !== 5u32": value.True,

		"'5' !== 5u16": value.True,

		"'5' !== 5u8": value.True,

		// Char
		"c'2' !== 25": value.True,

		"c'2' !== '2'": value.True,

		"c'7' !== c'7'": value.False,
		"c'b' !== c'b'": value.False,
		"c'c' !== c'g'": value.True,
		"c'7' !== c'8'": value.True,

		"c'2' !== 2.0": value.True,

		"c'3' !== 3bf": value.True,

		"c'9' !== 9f64": value.True,

		"c'1' !== 1f32": value.True,

		"c'1' !== 1i64": value.True,

		"c'5' !== 5i32": value.True,

		"c'5' !== 5i16": value.True,

		"c'5' !== 5i8": value.True,

		"c'1' !== 1u64": value.True,

		"c'5' !== 5u32": value.True,

		"c'5' !== 5u16": value.True,

		"c'5' !== 5u8": value.True,

		// Int
		"25 !== 25":  value.False,
		"-25 !== 25": value.True,
		"25 !== -25": value.True,
		"25 !== 28":  value.True,
		"28 !== 25":  value.True,

		"25 !== '25'": value.True,

		"7 !== c'7'": value.True,

		"25 !== 25.0": value.True,

		"25 !== 25bf": value.True,

		"25 !== 25f64": value.True,

		"25 !== 25f32": value.True,

		"1 !== 1i64": value.True,

		"5 !== 5i32": value.True,

		"5 !== 5i16": value.True,

		"5 !== 5i8": value.True,

		"1 !== 1u64": value.True,

		"5 !== 5u32": value.True,

		"5 !== 5u16": value.True,

		"5 !== 5u8": value.True,

		// Int64
		"25i64 !== 25": value.True,

		"25i64 !== '25'": value.True,

		"7i64 !== c'7'": value.True,

		"25i64 !== 25.0": value.True,

		"25i64 !== 25bf": value.True,

		"25i64 !== 25f64": value.True,

		"25i64 !== 25f32": value.True,

		"1i64 !== 1i64":   value.False,
		"4i64 !== -4i64":  value.True,
		"-8i64 !== 8i64":  value.True,
		"-8i64 !== -8i64": value.False,
		"91i64 !== 27i64": value.True,

		"5i64 !== 5i32": value.True,

		"5i64 !== 5i16": value.True,

		"5i64 !== 5i8": value.True,

		"1i64 !== 1u64": value.True,

		"5i64 !== 5u32": value.True,

		"5i64 !== 5u16": value.True,

		"5i64 !== 5u8": value.True,

		// Int32
		"25i32 !== 25": value.True,

		"25i32 !== '25'": value.True,

		"7i32 !== c'7'": value.True,

		"25i32 !== 25.0": value.True,

		"25i32 !== 25bf": value.True,

		"25i32 !== 25f64": value.True,

		"25i32 !== 25f32": value.True,

		"1i32 !== 1i64": value.True,

		"5i32 !== 5i32":  value.False,
		"4i32 !== -4i32": value.True,
		"-8i32 !== 8i32": value.True,
		"3i32 !== 71i32": value.True,

		"5i32 !== 5i16": value.True,

		"5i32 !== 5i8": value.True,

		"1i32 !== 1u64": value.True,

		"5i32 !== 5u32": value.True,

		"5i32 !== 5u16": value.True,

		"5i32 !== 5u8": value.True,

		// Int16
		"25i16 !== 25": value.True,

		"25i16 !== '25'": value.True,

		"7i16 !== c'7'": value.True,

		"25i16 !== 25.0": value.True,

		"25i16 !== 25bf": value.True,

		"25i16 !== 25f64": value.True,

		"25i16 !== 25f32": value.True,

		"1i16 !== 1i64": value.True,

		"5i16 !== 5i32": value.True,

		"5i16 !== 5i16":  value.False,
		"4i16 !== -4i16": value.True,
		"-8i16 !== 8i16": value.True,
		"3i16 !== 71i16": value.True,

		"5i16 !== 5i8": value.True,

		"1i16 !== 1u64": value.True,

		"5i16 !== 5u32": value.True,

		"5i16 !== 5u16": value.True,

		"5i16 !== 5u8": value.True,

		// Int8
		"25i8 !== 25": value.True,

		"25i8 !== '25'": value.True,

		"7i8 !== c'7'": value.True,

		"25i8 !== 25.0": value.True,

		"25i8 !== 25bf": value.True,

		"25i8 !== 25f64": value.True,

		"25i8 !== 25f32": value.True,

		"1i8 !== 1i64": value.True,

		"5i8 !== 5i32": value.True,

		"5i8 !== 5i16": value.True,

		"5i8 !== 5i8":  value.False,
		"4i8 !== -4i8": value.True,
		"-8i8 !== 8i8": value.True,
		"3i8 !== 71i8": value.True,

		"1i8 !== 1u64": value.True,

		"5i8 !== 5u32": value.True,

		"5i8 !== 5u16": value.True,

		"5i8 !== 5u8": value.True,

		// UInt64
		"25u64 !== 25": value.True,

		"25u64 !== '25'": value.True,

		"7u64 !== c'7'": value.True,

		"25u64 !== 25.0": value.True,

		"25u64 !== 25bf": value.True,

		"25u64 !== 25f64": value.True,

		"25u64 !== 25f32": value.True,

		"1u64 !== 1i64": value.True,

		"5u64 !== 5i32": value.True,

		"5u64 !== 5i16": value.True,

		"5u64 !== 5i8": value.True,

		"1u64 !== 1u64":   value.False,
		"91u64 !== 27u64": value.True,

		"5u64 !== 5u32": value.True,

		"5u64 !== 5u16": value.True,

		"5u64 !== 5u8": value.True,

		// UInt32
		"25u32 !== 25": value.True,

		"25u32 !== '25'": value.True,

		"7u32 !== c'7'": value.True,

		"25u32 !== 25.0": value.True,

		"25u32 !== 25bf": value.True,

		"25u32 !== 25f64": value.True,

		"25u32 !== 25f32": value.True,

		"1u32 !== 1i64": value.True,

		"5u32 !== 5i32": value.True,

		"5u32 !== 5i16": value.True,

		"5u32 !== 5i8": value.True,

		"1u32 !== 1u64": value.True,

		"5u32 !== 5u32":  value.False,
		"3u32 !== 71u32": value.True,

		"5u32 !== 5u16": value.True,

		"5u32 !== 5u8": value.True,

		// UInt16
		"25u16 !== 25": value.True,

		"25u16 !== '25'": value.True,

		"7u16 !== c'7'": value.True,

		"25u16 !== 25.0": value.True,

		"25u16 !== 25bf": value.True,

		"25u16 !== 25f64": value.True,

		"25u16 !== 25f32": value.True,

		"1u16 !== 1i64": value.True,

		"5u16 !== 5i32": value.True,

		"5u16 !== 5i16": value.True,

		"5u16 !== 5i8": value.True,

		"1u16 !== 1u64": value.True,

		"5u16 !== 5u32": value.True,

		"53000u16 !== 32767u16": value.True,
		"5u16 !== 5u16":         value.False,
		"3u16 !== 71u16":        value.True,

		"5u16 !== 5u8": value.True,

		// UInt8
		"25u8 !== 25": value.True,

		"25u8 !== '25'": value.True,

		"7u8 !== c'7'": value.True,

		"25u8 !== 25.0": value.True,

		"25u8 !== 25bf": value.True,

		"25u8 !== 25f64": value.True,

		"25u8 !== 25f32": value.True,

		"1u8 !== 1i64": value.True,

		"5u8 !== 5i32": value.True,

		"5u8 !== 5i16": value.True,

		"5u8 !== 5i8": value.True,

		"1u8 !== 1u64": value.True,

		"5u8 !== 5u32": value.True,

		"5u8 !== 5u16": value.True,

		"5u8 !== 5u8":  value.False,
		"3u8 !== 71u8": value.True,

		// Float
		"-73.0 !== 73.0": value.True,
		"73.0 !== -73.0": value.True,
		"25.0 !== 25.0":  value.False,
		"1.0 !== 1.2":    value.True,
		"1.2 !== 1.0":    value.True,
		"78.5 !== 78.5":  value.False,

		"8.25 !== '8.25'": value.True,

		"4.0 !== c'4'": value.True,

		"25.0 !== 25": value.True,

		"25.0 !== 25bf":   value.True,
		"15.5 !== 15.5bf": value.True,

		"25.0 !== 25f64":     value.True,
		"15.26 !== 15.26f64": value.True,

		"25.0 !== 25f32":   value.True,
		"15.5 !== 15.5f32": value.True,

		"1.0 !== 1i64":   value.True,
		"-8.0 !== -8i64": value.True,

		"1.0 !== 1i32":   value.True,
		"-8.0 !== -8i32": value.True,

		"1.0 !== 1i16":   value.True,
		"-8.0 !== -8i16": value.True,

		"1.0 !== 1i8":   value.True,
		"-8.0 !== -8i8": value.True,

		"1.0 !== 1u64": value.True,

		"1.0 !== 1u32": value.True,

		"1.0 !== 1u16": value.True,

		"1.0 !== 1u8": value.True,

		// Float64
		"25f64 !== 25.0":   value.True,
		"78.5f64 !== 78.5": value.True,

		"8.25f64 !== '8.25'": value.True,

		"4f64 !== c'4'": value.True,

		"25f64 !== 25": value.True,

		"25f64 !== 25bf":     value.True,
		"15.5f64 !== 15.5bf": value.True,

		"-73f64 !== 73f64":      value.True,
		"73f64 !== -73f64":      value.True,
		"25f64 !== 25f64":       value.False,
		"1f64 !== 1.2f64":       value.True,
		"15.26f64 !== 15.26f64": value.False,

		"25f64 !== 25f32":     value.True,
		"15.5f64 !== 15.5f32": value.True,

		"1f64 !== 1i64":   value.True,
		"-8f64 !== -8i64": value.True,

		"1f64 !== 1i32":   value.True,
		"-8f64 !== -8i32": value.True,

		"1f64 !== 1i16":   value.True,
		"-8f64 !== -8i16": value.True,

		"1f64 !== 1i8":   value.True,
		"-8f64 !== -8i8": value.True,

		"1f64 !== 1u64": value.True,

		"1f64 !== 1u32": value.True,

		"1f64 !== 1u16": value.True,

		"1f64 !== 1u8": value.True,

		// Float32
		"25f32 !== 25.0":   value.True,
		"78.5f32 !== 78.5": value.True,

		"8.25f32 !== '8.25'": value.True,

		"4f32 !== c'4'": value.True,

		"25f32 !== 25": value.True,

		"25f32 !== 25bf":     value.True,
		"15.5f32 !== 15.5bf": value.True,

		"25f32 !== 25f64":     value.True,
		"15.5f32 !== 15.5f64": value.True,

		"-73f32 !== 73f32":    value.True,
		"73f32 !== -73f32":    value.True,
		"25f32 !== 25f32":     value.False,
		"1f32 !== 1.2f32":     value.True,
		"15.5f32 !== 15.5f32": value.False,

		"1f32 !== 1i64":   value.True,
		"-8f32 !== -8i64": value.True,

		"1f32 !== 1i32":   value.True,
		"-8f32 !== -8i32": value.True,

		"1f32 !== 1i16":   value.True,
		"-8f32 !== -8i16": value.True,

		"1f32 !== 1i8":   value.True,
		"-8f32 !== -8i8": value.True,

		"1f32 !== 1u64": value.True,

		"1f32 !== 1u32": value.True,

		"1f32 !== 1u16": value.True,

		"1f32 !== 1u8": value.True,
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}

func TestVMSource_RightBitshift(t *testing.T) {
	tests := sourceTestTable{
		"Int >> String": {
			source: "3 >> 'foo'",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::String` can't be used as a bitshift operand",
			),
			wantStackTop: value.SmallInt(3),
		},
		"UInt16 >> Float": {
			source: "3u16 >> 5.2",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be used as a bitshift operand",
			),
			wantStackTop: value.UInt16(3),
		},
		"String >> Int": {
			source: "'36' >> 5",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `>>` is not available to value of class `Std::String`: \"36\"",
			),
			wantStackTop: value.String("36"),
		},

		"Int >> Int": {
			source:       "16 >> 2",
			wantStackTop: value.SmallInt(4),
		},
		"-Int >> Int": {
			source:       "-16 >> 2",
			wantStackTop: value.SmallInt(-4),
		},
		"Int >> -Int": {
			source:       "16 >> -2",
			wantStackTop: value.SmallInt(64),
		},
		"Int >> Int32": {
			source:       "39 >> 1i32",
			wantStackTop: value.SmallInt(19),
		},

		"Int64 >> Int64": {
			source:       "16i64 >> 2i64",
			wantStackTop: value.Int64(4),
		},
		"-Int64 >> Int64": {
			source:       "-16i64 >> 2i64",
			wantStackTop: value.Int64(-4),
		},
		"Int64 >> -Int64": {
			source:       "16i64 >> -2i64",
			wantStackTop: value.Int64(64),
		},
		"Int64 >> Int32": {
			source:       "39i64 >> 1i32",
			wantStackTop: value.Int64(19),
		},
		"Int64 >> UInt8": {
			source:       "120i64 >> 5u8",
			wantStackTop: value.Int64(3),
		},
		"Int64 >> Int": {
			source:       "54i64 >> 3",
			wantStackTop: value.Int64(6),
		},

		"Int32 >> Int32": {
			source:       "16i32 >> 2i32",
			wantStackTop: value.Int32(4),
		},
		"-Int32 >> Int32": {
			source:       "-16i32 >> 2i32",
			wantStackTop: value.Int32(-4),
		},
		"Int32 >> -Int32": {
			source:       "16i32 >> -2i32",
			wantStackTop: value.Int32(64),
		},
		"Int32 >> Int16": {
			source:       "39i32 >> 1i16",
			wantStackTop: value.Int32(19),
		},
		"Int32 >> UInt8": {
			source:       "120i32 >> 5u8",
			wantStackTop: value.Int32(3),
		},
		"Int32 >> Int": {
			source:       "54i32 >> 3",
			wantStackTop: value.Int32(6),
		},

		"Int16 >> Int16": {
			source:       "16i16 >> 2i16",
			wantStackTop: value.Int16(4),
		},
		"-Int16 >> Int16": {
			source:       "-16i16 >> 2i16",
			wantStackTop: value.Int16(-4),
		},
		"Int16 >> -Int16": {
			source:       "16i16 >> -2i16",
			wantStackTop: value.Int16(64),
		},
		"Int16 >> Int32": {
			source:       "39i16 >> 1i32",
			wantStackTop: value.Int16(19),
		},
		"Int16 >> UInt8": {
			source:       "120i16 >> 5u8",
			wantStackTop: value.Int16(3),
		},
		"Int16 >> Int": {
			source:       "54i16 >> 3",
			wantStackTop: value.Int16(6),
		},

		"Int8 >> Int8": {
			source:       "16i8 >> 2i8",
			wantStackTop: value.Int8(4),
		},
		"-Int8 >> Int8": {
			source:       "-16i8 >> 2i8",
			wantStackTop: value.Int8(-4),
		},
		"Int8 >> -Int8": {
			source:       "16i8 >> -2i8",
			wantStackTop: value.Int8(64),
		},
		"Int8 >> Int16": {
			source:       "39i8 >> 1i16",
			wantStackTop: value.Int8(19),
		},
		"Int8 >> UInt8": {
			source:       "120i8 >> 5u8",
			wantStackTop: value.Int8(3),
		},
		"Int8 >> Int": {
			source:       "54i8 >> 3",
			wantStackTop: value.Int8(6),
		},

		"UInt64 >> UInt64": {
			source:       "16u64 >> 2u64",
			wantStackTop: value.UInt64(4),
		},
		"UInt64 >> -Int": {
			source:       "16u64 >> -2",
			wantStackTop: value.UInt64(64),
		},
		"UInt64 >> Int32": {
			source:       "39u64 >> 1i32",
			wantStackTop: value.UInt64(19),
		},

		"UInt32 >> UInt32": {
			source:       "16u32 >> 2u32",
			wantStackTop: value.UInt32(4),
		},
		"UInt32 >> -Int": {
			source:       "16u32 >> -2",
			wantStackTop: value.UInt32(64),
		},
		"UInt32 >> Int32": {
			source:       "39u32 >> 1i32",
			wantStackTop: value.UInt32(19),
		},

		"UInt16 >> UInt16": {
			source:       "16u16 >> 2u16",
			wantStackTop: value.UInt16(4),
		},
		"UInt16 >> -Int": {
			source:       "16u16 >> -2",
			wantStackTop: value.UInt16(64),
		},
		"UInt16 >> Int32": {
			source:       "39u16 >> 1i32",
			wantStackTop: value.UInt16(19),
		},

		"UInt8 >> UInt8": {
			source:       "16u8 >> 2u8",
			wantStackTop: value.UInt8(4),
		},
		"UInt8 >> -Int": {
			source:       "16u8 >> -2",
			wantStackTop: value.UInt8(64),
		},
		"UInt8 >> Int32": {
			source:       "39u8 >> 1i32",
			wantStackTop: value.UInt8(19),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LogicalRightBitshift(t *testing.T) {
	tests := sourceTestTable{
		"Int >>> String": {
			source: "3 >>> 'foo'",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `>>>` is not available to value of class `Std::SmallInt`: 3",
			),
			wantStackTop: value.SmallInt(3),
		},
		"Int64 >>> String": {
			source: "3i64 >>> 'foo'",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::String` can't be used as a bitshift operand",
			),
			wantStackTop: value.Int64(3),
		},
		"UInt16 >>> Float": {
			source: "3u16 >>> 5.2",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be used as a bitshift operand",
			),
			wantStackTop: value.UInt16(3),
		},
		"String >>> Int": {
			source: "'36' >>> 5",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `>>>` is not available to value of class `Std::String`: \"36\"",
			),
			wantStackTop: value.String("36"),
		},

		"Int64 >>> Int64": {
			source:       "16i64 >>> 2i64",
			wantStackTop: value.Int64(4),
		},
		"-Int64 >>> Int64": {
			source:       "-16i64 >>> 2i64",
			wantStackTop: value.Int64(4611686018427387900),
		},
		"Int64 >>> -Int64": {
			source:       "16i64 >>> -2i64",
			wantStackTop: value.Int64(64),
		},
		"Int64 >>> Int32": {
			source:       "39i64 >>> 1i32",
			wantStackTop: value.Int64(19),
		},
		"Int64 >>> UInt8": {
			source:       "120i64 >>> 5u8",
			wantStackTop: value.Int64(3),
		},
		"Int64 >>> Int": {
			source:       "54i64 >>> 3",
			wantStackTop: value.Int64(6),
		},

		"Int32 >>> Int32": {
			source:       "16i32 >>> 2i32",
			wantStackTop: value.Int32(4),
		},
		"-Int32 >>> Int32": {
			source:       "-16i32 >>> 2i32",
			wantStackTop: value.Int32(1073741820),
		},
		"Int32 >>> -Int32": {
			source:       "16i32 >>> -2i32",
			wantStackTop: value.Int32(64),
		},
		"Int32 >>> Int16": {
			source:       "39i32 >>> 1i16",
			wantStackTop: value.Int32(19),
		},
		"Int32 >>> UInt8": {
			source:       "120i32 >>> 5u8",
			wantStackTop: value.Int32(3),
		},
		"Int32 >>> Int": {
			source:       "54i32 >>> 3",
			wantStackTop: value.Int32(6),
		},

		"Int16 >>> Int16": {
			source:       "16i16 >>> 2i16",
			wantStackTop: value.Int16(4),
		},
		"-Int16 >>> Int16": {
			source:       "-16i16 >>> 2i16",
			wantStackTop: value.Int16(16380),
		},
		"Int16 >>> -Int16": {
			source:       "16i16 >>> -2i16",
			wantStackTop: value.Int16(64),
		},
		"Int16 >>> Int32": {
			source:       "39i16 >>> 1i32",
			wantStackTop: value.Int16(19),
		},
		"Int16 >>> UInt8": {
			source:       "120i16 >>> 5u8",
			wantStackTop: value.Int16(3),
		},
		"Int16 >>> Int": {
			source:       "54i16 >>> 3",
			wantStackTop: value.Int16(6),
		},

		"Int8 >>> Int8": {
			source:       "16i8 >>> 2i8",
			wantStackTop: value.Int8(4),
		},
		"-Int8 >>> Int8": {
			source:       "-16i8 >>> 2i8",
			wantStackTop: value.Int8(60),
		},
		"Int8 >>> -Int8": {
			source:       "16i8 >>> -2i8",
			wantStackTop: value.Int8(64),
		},
		"Int8 >>> Int16": {
			source:       "39i8 >>> 1i16",
			wantStackTop: value.Int8(19),
		},
		"Int8 >>> UInt8": {
			source:       "120i8 >>> 5u8",
			wantStackTop: value.Int8(3),
		},
		"Int8 >>> Int": {
			source:       "54i8 >>> 3",
			wantStackTop: value.Int8(6),
		},

		"UInt64 >>> UInt64": {
			source:       "16u64 >>> 2u64",
			wantStackTop: value.UInt64(4),
		},
		"UInt64 >>> -Int": {
			source:       "16u64 >>> -2",
			wantStackTop: value.UInt64(64),
		},
		"UInt64 >>> Int32": {
			source:       "39u64 >>> 1i32",
			wantStackTop: value.UInt64(19),
		},

		"UInt32 >>> UInt32": {
			source:       "16u32 >>> 2u32",
			wantStackTop: value.UInt32(4),
		},
		"UInt32 >>> -Int": {
			source:       "16u32 >>> -2",
			wantStackTop: value.UInt32(64),
		},
		"UInt32 >>> Int32": {
			source:       "39u32 >>> 1i32",
			wantStackTop: value.UInt32(19),
		},

		"UInt16 >>> UInt16": {
			source:       "16u16 >>> 2u16",
			wantStackTop: value.UInt16(4),
		},
		"UInt16 >>> -Int": {
			source:       "16u16 >>> -2",
			wantStackTop: value.UInt16(64),
		},
		"UInt16 >>> Int32": {
			source:       "39u16 >>> 1i32",
			wantStackTop: value.UInt16(19),
		},

		"UInt8 >>> UInt8": {
			source:       "16u8 >>> 2u8",
			wantStackTop: value.UInt8(4),
		},
		"UInt8 >>> -Int": {
			source:       "16u8 >>> -2",
			wantStackTop: value.UInt8(64),
		},
		"UInt8 >>> Int32": {
			source:       "39u8 >>> 1i32",
			wantStackTop: value.UInt8(19),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LeftBitshift(t *testing.T) {
	tests := sourceTestTable{
		"Int << String": {
			source: "3 << 'foo'",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::String` can't be used as a bitshift operand",
			),
			wantStackTop: value.SmallInt(3),
		},
		"UInt16 << Float": {
			source: "3u16 << 5.2",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be used as a bitshift operand",
			),
			wantStackTop: value.UInt16(3),
		},
		"String << Int": {
			source: "'36' << 5",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `<<` is not available to value of class `Std::String`: \"36\"",
			),
			wantStackTop: value.String("36"),
		},

		"Int << Int": {
			source:       "16 << 2",
			wantStackTop: value.SmallInt(64),
		},
		"-Int << Int": {
			source:       "-16 << 2",
			wantStackTop: value.SmallInt(-64),
		},
		"Int << -Int": {
			source:       "16 << -2",
			wantStackTop: value.SmallInt(4),
		},
		"Int << Int32": {
			source:       "39 << 1i32",
			wantStackTop: value.SmallInt(78),
		},

		"Int64 << Int64": {
			source:       "16i64 << 2i64",
			wantStackTop: value.Int64(64),
		},
		"-Int64 << Int64": {
			source:       "-16i64 << 2i64",
			wantStackTop: value.Int64(-64),
		},
		"Int64 << -Int64": {
			source:       "16i64 << -2i64",
			wantStackTop: value.Int64(4),
		},
		"Int64 << Int32": {
			source:       "39i64 << 1i32",
			wantStackTop: value.Int64(78),
		},
		"Int64 << UInt8": {
			source:       "120i64 << 5u8",
			wantStackTop: value.Int64(3840),
		},
		"Int64 << Int": {
			source:       "54i64 << 3",
			wantStackTop: value.Int64(432),
		},

		"Int32 << Int32": {
			source:       "16i32 << 2i32",
			wantStackTop: value.Int32(64),
		},
		"-Int32 << Int32": {
			source:       "-16i32 << 2i32",
			wantStackTop: value.Int32(-64),
		},
		"Int32 << -Int32": {
			source:       "16i32 << -2i32",
			wantStackTop: value.Int32(4),
		},
		"Int32 << Int16": {
			source:       "39i32 << 1i16",
			wantStackTop: value.Int32(78),
		},
		"Int32 << UInt8": {
			source:       "120i32 << 5u8",
			wantStackTop: value.Int32(3840),
		},
		"Int32 << Int": {
			source:       "54i32 << 3",
			wantStackTop: value.Int32(432),
		},

		"Int16 << Int16": {
			source:       "16i16 << 2i16",
			wantStackTop: value.Int16(64),
		},
		"-Int16 << Int16": {
			source:       "-16i16 << 2i16",
			wantStackTop: value.Int16(-64),
		},
		"Int16 << -Int16": {
			source:       "16i16 << -2i16",
			wantStackTop: value.Int16(4),
		},
		"Int16 << Int32": {
			source:       "39i16 << 1i32",
			wantStackTop: value.Int16(78),
		},
		"Int16 << UInt8": {
			source:       "120i16 << 5u8",
			wantStackTop: value.Int16(3840),
		},
		"Int16 << Int": {
			source:       "54i16 << 3",
			wantStackTop: value.Int16(432),
		},

		"Int8 << Int8": {
			source:       "16i8 << 2i8",
			wantStackTop: value.Int8(64),
		},
		"-Int8 << Int8": {
			source:       "-16i8 << 2i8",
			wantStackTop: value.Int8(-64),
		},
		"Int8 << -Int8": {
			source:       "16i8 << -2i8",
			wantStackTop: value.Int8(4),
		},
		"Int8 << Int16": {
			source:       "39i8 << 1i16",
			wantStackTop: value.Int8(78),
		},
		"Int8 << UInt8": {
			source:       "120i8 << 5u8",
			wantStackTop: value.Int8(0),
		},
		"Int8 << Int": {
			source:       "54i8 << 3",
			wantStackTop: value.Int8(-80),
		},

		"UInt64 << UInt64": {
			source:       "16u64 << 2u64",
			wantStackTop: value.UInt64(64),
		},
		"UInt64 << -Int": {
			source:       "16u64 << -2",
			wantStackTop: value.UInt64(4),
		},
		"UInt64 << Int32": {
			source:       "39u64 << 1i32",
			wantStackTop: value.UInt64(78),
		},

		"UInt32 << UInt32": {
			source:       "16u32 << 2u32",
			wantStackTop: value.UInt32(64),
		},
		"UInt32 << -Int": {
			source:       "16u32 << -2",
			wantStackTop: value.UInt32(4),
		},
		"UInt32 << Int32": {
			source:       "39u32 << 1i32",
			wantStackTop: value.UInt32(78),
		},

		"UInt16 << UInt16": {
			source:       "16u16 << 2u16",
			wantStackTop: value.UInt16(64),
		},
		"UInt16 << -Int": {
			source:       "16u16 << -2",
			wantStackTop: value.UInt16(4),
		},
		"UInt16 << Int32": {
			source:       "39u16 << 1i32",
			wantStackTop: value.UInt16(78),
		},

		"UInt8 << UInt8": {
			source:       "16u8 << 2u8",
			wantStackTop: value.UInt8(64),
		},
		"UInt8 << -Int": {
			source:       "16u8 << -2",
			wantStackTop: value.UInt8(4),
		},
		"UInt8 << Int32": {
			source:       "39u8 << 1i32",
			wantStackTop: value.UInt8(78),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LogicalLeftBitshift(t *testing.T) {
	tests := sourceTestTable{
		"Int64 <<< String": {
			source: "3i64 <<< 'foo'",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::String` can't be used as a bitshift operand",
			),
			wantStackTop: value.Int64(3),
		},
		"UInt16 <<< Float": {
			source: "3u16 <<< 5.2",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be used as a bitshift operand",
			),
			wantStackTop: value.UInt16(3),
		},
		"String <<< Int": {
			source: "'36' <<< 5",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `<<<` is not available to value of class `Std::String`: \"36\"",
			),
			wantStackTop: value.String("36"),
		},
		"Int <<< Int": {
			source: "16 <<< 2",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `<<<` is not available to value of class `Std::SmallInt`: 16",
			),
			wantStackTop: value.SmallInt(16),
		},

		"Int64 <<< Int64": {
			source:       "16i64 <<< 2i64",
			wantStackTop: value.Int64(64),
		},
		"-Int64 <<< Int64": {
			source:       "-16i64 <<< 2i64",
			wantStackTop: value.Int64(-64),
		},
		"Int64 <<< -Int64": {
			source:       "16i64 <<< -2i64",
			wantStackTop: value.Int64(4),
		},
		"Int64 <<< Int32": {
			source:       "39i64 <<< 1i32",
			wantStackTop: value.Int64(78),
		},
		"Int64 <<< UInt8": {
			source:       "120i64 <<< 5u8",
			wantStackTop: value.Int64(3840),
		},
		"Int64 <<< Int": {
			source:       "54i64 <<< 3",
			wantStackTop: value.Int64(432),
		},

		"Int32 <<< Int32": {
			source:       "16i32 <<< 2i32",
			wantStackTop: value.Int32(64),
		},
		"-Int32 <<< Int32": {
			source:       "-16i32 <<< 2i32",
			wantStackTop: value.Int32(-64),
		},
		"Int32 <<< -Int32": {
			source:       "16i32 <<< -2i32",
			wantStackTop: value.Int32(4),
		},
		"Int32 <<< Int16": {
			source:       "39i32 <<< 1i16",
			wantStackTop: value.Int32(78),
		},
		"Int32 <<< UInt8": {
			source:       "120i32 <<< 5u8",
			wantStackTop: value.Int32(3840),
		},
		"Int32 <<< Int": {
			source:       "54i32 <<< 3",
			wantStackTop: value.Int32(432),
		},

		"Int16 <<< Int16": {
			source:       "16i16 <<< 2i16",
			wantStackTop: value.Int16(64),
		},
		"-Int16 <<< Int16": {
			source:       "-16i16 <<< 2i16",
			wantStackTop: value.Int16(-64),
		},
		"Int16 <<< -Int16": {
			source:       "16i16 <<< -2i16",
			wantStackTop: value.Int16(4),
		},
		"Int16 <<< Int32": {
			source:       "39i16 <<< 1i32",
			wantStackTop: value.Int16(78),
		},
		"Int16 <<< UInt8": {
			source:       "120i16 <<< 5u8",
			wantStackTop: value.Int16(3840),
		},
		"Int16 <<< Int": {
			source:       "54i16 <<< 3",
			wantStackTop: value.Int16(432),
		},

		"Int8 <<< Int8": {
			source:       "16i8 <<< 2i8",
			wantStackTop: value.Int8(64),
		},
		"-Int8 <<< Int8": {
			source:       "-16i8 <<< 2i8",
			wantStackTop: value.Int8(-64),
		},
		"Int8 <<< -Int8": {
			source:       "16i8 <<< -2i8",
			wantStackTop: value.Int8(4),
		},
		"Int8 <<< Int16": {
			source:       "39i8 <<< 1i16",
			wantStackTop: value.Int8(78),
		},
		"Int8 <<< UInt8": {
			source:       "120i8 <<< 5u8",
			wantStackTop: value.Int8(0),
		},
		"Int8 <<< Int": {
			source:       "54i8 <<< 3",
			wantStackTop: value.Int8(-80),
		},

		"UInt64 <<< UInt64": {
			source:       "16u64 <<< 2u64",
			wantStackTop: value.UInt64(64),
		},
		"UInt64 <<< -Int": {
			source:       "16u64 <<< -2",
			wantStackTop: value.UInt64(4),
		},
		"UInt64 <<< Int32": {
			source:       "39u64 <<< 1i32",
			wantStackTop: value.UInt64(78),
		},

		"UInt32 <<< UInt32": {
			source:       "16u32 <<< 2u32",
			wantStackTop: value.UInt32(64),
		},
		"UInt32 <<< -Int": {
			source:       "16u32 <<< -2",
			wantStackTop: value.UInt32(4),
		},
		"UInt32 <<< Int32": {
			source:       "39u32 <<< 1i32",
			wantStackTop: value.UInt32(78),
		},

		"UInt16 <<< UInt16": {
			source:       "16u16 <<< 2u16",
			wantStackTop: value.UInt16(64),
		},
		"UInt16 <<< -Int": {
			source:       "16u16 <<< -2",
			wantStackTop: value.UInt16(4),
		},
		"UInt16 <<< Int32": {
			source:       "39u16 <<< 1i32",
			wantStackTop: value.UInt16(78),
		},

		"UInt8 <<< UInt8": {
			source:       "16u8 <<< 2u8",
			wantStackTop: value.UInt8(64),
		},
		"UInt8 <<< -Int": {
			source:       "16u8 <<< -2",
			wantStackTop: value.UInt8(4),
		},
		"UInt8 <<< Int32": {
			source:       "39u8 <<< 1i32",
			wantStackTop: value.UInt8(78),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_BitwiseAnd(t *testing.T) {
	tests := sourceTestTable{
		"Int64 & String": {
			source: "3i64 & 'foo'",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::String` can't be coerced into `Std::Int64`",
			),
			wantStackTop: value.Int64(3),
		},
		"Int64 & SmallInt": {
			source: "3i64 & 5",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int64`",
			),
			wantStackTop: value.Int64(3),
		},
		"UInt16 & Float": {
			source: "3u16 & 5.2",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt16`",
			),
			wantStackTop: value.UInt16(3),
		},
		"String & Int": {
			source: "'36' & 5",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `&` is not available to value of class `Std::String`: \"36\"",
			),
			wantStackTop: value.String("36"),
		},
		"Float & Int": {
			source: "3.6 & 5",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `&` is not available to value of class `Std::Float`: 3.6",
			),
			wantStackTop: value.Float(3.6),
		},
		"Int & Int": {
			source:       "25 & 14",
			wantStackTop: value.SmallInt(8),
		},
		"Int & BigInt": {
			source:       "255 & 9223372036857247042",
			wantStackTop: value.SmallInt(66),
		},
		"Int8 & Int8": {
			source:       "59i8 & 122i8",
			wantStackTop: value.Int8(58),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_BitwiseOr(t *testing.T) {
	tests := sourceTestTable{
		"Int64 | String": {
			source: "3i64 | 'foo'",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::String` can't be coerced into `Std::Int64`",
			),
			wantStackTop: value.Int64(3),
		},
		"Int64 | SmallInt": {
			source: "3i64 | 5",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int64`",
			),
			wantStackTop: value.Int64(3),
		},
		"UInt16 | Float": {
			source: "3u16 | 5.2",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt16`",
			),
			wantStackTop: value.UInt16(3),
		},
		"String | Int": {
			source: "'36' | 5",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `|` is not available to value of class `Std::String`: \"36\"",
			),
			wantStackTop: value.String("36"),
		},
		"Float | Int": {
			source: "3.6 | 5",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `|` is not available to value of class `Std::Float`: 3.6",
			),
			wantStackTop: value.Float(3.6),
		},
		"Int | Int": {
			source:       "25 | 14",
			wantStackTop: value.SmallInt(31),
		},
		"Int | BigInt": {
			source:       "255 | 9223372036857247042",
			wantStackTop: value.ParseBigIntPanic("9223372036857247231", 10),
		},
		"Int8 | Int8": {
			source:       "59i8 | 122i8",
			wantStackTop: value.Int8(123),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_BitwiseXor(t *testing.T) {
	tests := sourceTestTable{
		"Int64 ^ String": {
			source: "3i64 ^ 'foo'",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::String` can't be coerced into `Std::Int64`",
			),
			wantStackTop: value.Int64(3),
		},
		"Int64 ^ SmallInt": {
			source: "3i64 ^ 5",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::SmallInt` can't be coerced into `Std::Int64`",
			),
			wantStackTop: value.Int64(3),
		},
		"UInt16 ^ Float": {
			source: "3u16 ^ 5.2",
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`Std::Float` can't be coerced into `Std::UInt16`",
			),
			wantStackTop: value.UInt16(3),
		},
		"String ^ Int": {
			source: "'36' ^ 5",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `^` is not available to value of class `Std::String`: \"36\"",
			),
			wantStackTop: value.String("36"),
		},
		"Float ^ Int": {
			source: "3.6 ^ 5",
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `^` is not available to value of class `Std::Float`: 3.6",
			),
			wantStackTop: value.Float(3.6),
		},
		"Int ^ Int": {
			source:       "25 ^ 14",
			wantStackTop: value.SmallInt(23),
		},
		"Int ^ BigInt": {
			source:       "255 ^ 9223372036857247042",
			wantStackTop: value.ParseBigIntPanic("9223372036857247165", 10),
		},
		"Int8 ^ Int8": {
			source:       "59i8 ^ 122i8",
			wantStackTop: value.Int8(65),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_NumericFor(t *testing.T) {
	tests := sourceTestTable{
		"calculate the sum of consecutive natural numbers": {
			source: `
				a := 0
				for i := 1; i < 6; i += 1
					a += i
				end
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"create a repeated string": {
			source: `
				a := ""
				for i := 20; i > 0; i -= 2
					a += "-"
				end
				a
			`,
			wantStackTop: value.String("----------"),
		},
		"calculate the factorial of 10": {
			source: `
				a := 1
				for i := 2; i <= 10; i += 1
					a *= i
				end
				a
			`,
			wantStackTop: value.SmallInt(3628800),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_GetModuleConstant(t *testing.T) {
	tests := simpleSourceTestTable{
		"::Std":                     value.StdModule,
		"::Std::Int":                value.IntClass,
		"::Std::Float::INF":         value.FloatInf(),
		"a := ::Std::Float; a::INF": value.FloatInf(),
	}

	for source, want := range tests {
		t.Run(source, func(t *testing.T) {
			vmSimpleSourceTest(source, want, t)
		})
	}
}

func TestVMSource_DefineModuleConstant(t *testing.T) {
	tests := sourceTestTable{
		"Set constant under Root": {
			source:       "::Foo := 3i64",
			wantStackTop: value.Int64(3),
			teardown:     func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"Set constant under Root and read it": {
			source: `
				::Foo := 3i64
				::Foo
			`,
			wantStackTop: value.Int64(3),
			teardown:     func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"Set constant under nested modules": {
			source:       `::Std::Int::Foo := 3i64`,
			wantStackTop: value.Int64(3),
			teardown:     func() { value.IntClass.Constants.DeleteString("Foo") },
		},
		"Set constant under a variable": {
			source: `
				a := ::Std::Int
				a::Bar := "baz"
			`,
			wantStackTop: value.String("baz"),
			teardown:     func() { value.IntClass.Constants.DeleteString("Bar") },
		},
		"Set constant under a variable and read it": {
			source: `
				a := ::Std::Int
				a::Bar := "baz"
				::Std::Int::Bar
			`,
			wantStackTop: value.String("baz"),
			teardown:     func() { value.IntClass.Constants.DeleteString("Bar") },
		},
		"Set a constant under Int": {
			source: `
				a := 3
				a::Foo := 10
			`,
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`3` is not a module",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_DefineClass(t *testing.T) {
	tests := sourceTestTable{
		"anonymous class without a body": {
			source:       "class; end",
			wantStackTop: value.NewClass(),
		},
		"class without a body with a relative name": {
			source: "class Foo; end",
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"class without a body with an absolute name": {
			source: "class ::Foo; end",
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"class without a body with a parent": {
			source: "class Foo < ::Std::Error; end",
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithParent(value.ErrorClass),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"anonymous class without a body with a parent": {
			source: "class < ::Std::Error; end",
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithParent(value.ErrorClass),
			),
		},
		"class with a body": {
			source: `
				class Foo
					a := 5
					Bar := a - 2
				end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("Bar"): value.SmallInt(3),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"anonymous class with a body": {
			source: `
				class
					a := 5
					Bar := a - 2
				end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("Bar"): value.SmallInt(3),
					},
				),
			),
		},
		"nested classes": {
			source: `
				class Gdańsk
					class Gdynia
						class Sopot
							Trójmiasto := "jest super"
							::Gdańsk::Warszawa := "to stolica"
						end
					end
				end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Gdańsk"),
				value.ClassWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("Gdynia"): value.NewClassWithOptions(
							value.ClassWithName("Gdańsk::Gdynia"),
							value.ClassWithConstants(
								value.SimpleSymbolMap{
									value.SymbolTable.Add("Sopot"): value.NewClassWithOptions(
										value.ClassWithName("Gdańsk::Gdynia::Sopot"),
										value.ClassWithConstants(
											value.SimpleSymbolMap{
												value.SymbolTable.Add("Trójmiasto"): value.String("jest super"),
											},
										),
									),
								},
							),
						),
						value.SymbolTable.Add("Warszawa"): value.String("to stolica"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Gdańsk") },
		},
		"open an existing class": {
			source: `
				class Foo
					FIRST_CONSTANT := "oguem"
				end

				class Foo
					SECOND_CONSTANT := "całe te"
				end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("FIRST_CONSTANT"):  value.String("oguem"),
						value.SymbolTable.Add("SECOND_CONSTANT"): value.String("całe te"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"superclass mismatch": {
			source: `
				class Foo; end

				class Bar < ::Foo
					FIRST_CONSTANT := "oguem"
				end

				class Bar < ::Std::Error
					SECOND_CONSTANT := "całe te"
				end
			`,
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"superclass mismatch in Bar, expected: Foo, got: Std::Error",
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("Bar")
			},
		},
		"incorrect superclass": {
			source: `
				A := 3
				class Foo < ::A; end
			`,
			wantRuntimeErr: value.NewError(
				value.TypeErrorClass,
				"`3` can't be used as a superclass",
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("A")
			},
		},
		"redefined constant": {
			source: `
				Foo := 3
				class Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.RedefinedConstantErrorClass,
				"module Root already has a constant named `:Foo`",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_DefineModule(t *testing.T) {
	tests := sourceTestTable{
		"anonymous module without a body": {
			source:       "module; end",
			wantStackTop: value.NewModule(),
		},
		"module without a body with a relative name": {
			source: "module Foo; end",
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"module without a body with an absolute name": {
			source: "module ::Foo; end",
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"module with a body": {
			source: `
				module Foo
					a := 5
					Bar := a - 2
				end
			`,
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithParent(value.ModuleClass),
					),
				),
				value.ModuleWithName("Foo"),
				value.ModuleWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("Bar"): value.SmallInt(3),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"anonymous module with a body": {
			source: `
				module
					a := 5
					Bar := a - 2
				end
			`,
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithParent(value.ModuleClass),
					),
				),
				value.ModuleWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("Bar"): value.SmallInt(3),
					},
				),
			),
		},
		"nested modules": {
			source: `
				module Gdańsk
					module Gdynia
						module Sopot
							Trójmiasto := "jest super"
							::Gdańsk::Warszawa := "to stolica"
						end
					end
				end
			`,
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithName("Gdańsk"),
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithParent(value.ModuleClass),
					),
				),
				value.ModuleWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("Gdynia"): value.NewModuleWithOptions(
							value.ModuleWithName("Gdańsk::Gdynia"),
							value.ModuleWithClass(
								value.NewClassWithOptions(
									value.ClassWithSingleton(),
									value.ClassWithParent(value.ModuleClass),
								),
							),
							value.ModuleWithConstants(
								value.SimpleSymbolMap{
									value.SymbolTable.Add("Sopot"): value.NewModuleWithOptions(
										value.ModuleWithName("Gdańsk::Gdynia::Sopot"),
										value.ModuleWithClass(
											value.NewClassWithOptions(
												value.ClassWithSingleton(),
												value.ClassWithParent(value.ModuleClass),
											),
										),
										value.ModuleWithConstants(
											value.SimpleSymbolMap{
												value.SymbolTable.Add("Trójmiasto"): value.String("jest super"),
											},
										),
									),
								},
							),
						),
						value.SymbolTable.Add("Warszawa"): value.String("to stolica"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Gdańsk") },
		},
		"open an existing module": {
			source: `
				module Foo
					FIRST_CONSTANT := "oguem"
				end

				module Foo
					SECOND_CONSTANT := "całe te"
				end
			`,
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithParent(value.ModuleClass),
					),
				),
				value.ModuleWithName("Foo"),
				value.ModuleWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("FIRST_CONSTANT"):  value.String("oguem"),
						value.SymbolTable.Add("SECOND_CONSTANT"): value.String("całe te"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"redefined constant": {
			source: `
				Foo := 3
				module Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.RedefinedConstantErrorClass,
				"module Root already has a constant named `:Foo`",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"redefined class as module": {
			source: `
				class Foo; end
				module Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.RedefinedConstantErrorClass,
				"module Root already has a constant named `:Foo`",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_DefineMethod(t *testing.T) {
	tests := sourceTestTable{
		"define a method in top level": {
			source: `
				def foo: Symbol then :bar
			`,
			wantStackTop: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 2),
				},
				Location: L(P(5, 2, 5), P(29, 2, 29)),
				Name:     value.SymbolTable.Add("foo"),
				Values: []value.Value{
					value.SymbolTable.Add("bar"),
				},
			},
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("bar"))
			},
		},
		"define a method with positional arguments in top level": {
			source: `
				def foo(a: Int, b: Int): Int
					c := 5
					a + b + c
				end
			`,
			wantStackTop: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.ADD),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 1),
				},
				Location: L(P(5, 2, 5), P(67, 5, 7)),
				Name:     value.SymbolTable.Add("foo"),
				Parameters: []value.Symbol{
					value.SymbolTable.Add("a"),
					value.SymbolTable.Add("b"),
				},
				Values: []value.Value{
					value.SmallInt(5),
				},
			},
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("bar"))
			},
		},
		"define a method with positional arguments in a class": {
			source: `
				class Bar
					def foo(a: Int, b: Int): Int
						c := 5
						a + b + c
					end
				end
			`,
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Bar"),
				value.ClassWithMethods(
					value.MethodMap{
						value.SymbolTable.Add("foo"): &value.BytecodeFunction{
							Instructions: []byte{
								byte(bytecode.PREP_LOCALS8), 1,
								byte(bytecode.LOAD_VALUE8), 0,
								byte(bytecode.SET_LOCAL8), 3,
								byte(bytecode.POP),
								byte(bytecode.GET_LOCAL8), 1,
								byte(bytecode.GET_LOCAL8), 2,
								byte(bytecode.ADD),
								byte(bytecode.GET_LOCAL8), 3,
								byte(bytecode.ADD),
								byte(bytecode.RETURN),
							},
							LineInfoList: bytecode.LineInfoList{
								bytecode.NewLineInfo(4, 4),
								bytecode.NewLineInfo(5, 5),
								bytecode.NewLineInfo(6, 1),
							},
							Location: L(P(20, 3, 6), P(85, 6, 8)),
							Name:     value.SymbolTable.Add("foo"),
							Parameters: []value.Symbol{
								value.SymbolTable.Add("a"),
								value.SymbolTable.Add("b"),
							},
							Values: []value.Value{
								value.SmallInt(5),
							},
						},
					},
				),
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Bar")
			},
		},
		"define a method with positional arguments in a module": {
			source: `
				module Bar
					def foo(a: Int, b: Int): Int
						c := 5
						a + b + c
					end
				end
			`,
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithName("Bar"),
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithParent(value.ModuleClass),
						value.ClassWithMethods(
							value.MethodMap{
								value.SymbolTable.Add("foo"): &value.BytecodeFunction{
									Instructions: []byte{
										byte(bytecode.PREP_LOCALS8), 1,
										byte(bytecode.LOAD_VALUE8), 0,
										byte(bytecode.SET_LOCAL8), 3,
										byte(bytecode.POP),
										byte(bytecode.GET_LOCAL8), 1,
										byte(bytecode.GET_LOCAL8), 2,
										byte(bytecode.ADD),
										byte(bytecode.GET_LOCAL8), 3,
										byte(bytecode.ADD),
										byte(bytecode.RETURN),
									},
									LineInfoList: bytecode.LineInfoList{
										bytecode.NewLineInfo(4, 4),
										bytecode.NewLineInfo(5, 5),
										bytecode.NewLineInfo(6, 1),
									},
									Location: L(P(21, 3, 6), P(86, 6, 8)),
									Name:     value.SymbolTable.Add("foo"),
									Parameters: []value.Symbol{
										value.SymbolTable.Add("a"),
										value.SymbolTable.Add("b"),
									},
									Values: []value.Value{
										value.SmallInt(5),
									},
								},
							},
						),
					),
				),
			),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Bar")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_CallMethod(t *testing.T) {
	tests := sourceTestTable{
		"call a global method without arguments": {
			source: `
				def foo: Symbol
					:bar
				end

				self.foo
			`,
			wantStackTop: value.SymbolTable.Add("bar"),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("foo"))
			},
		},
		"call a global method with positional arguments": {
			source: `
				def add(a: Int, b: Int): Int
					a + b
				end

				self.add(5, 9)
			`,
			wantStackTop: value.SmallInt(14),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.SymbolTable.Add("add"))
			},
		},
		"call a module method without arguments": {
			source: `
				module Foo
					def bar: Symbol
						:baz
					end
				end

				::Foo.bar
			`,
			wantStackTop: value.SymbolTable.Add("baz"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
		},
		"call a module method with positional arguments": {
			source: `
				module Foo
					def add(a: Int, b: Int): Int
						a + b
					end
				end

				::Foo.add 4, 12
			`,
			wantStackTop: value.SmallInt(16),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
		},
		"call an instance method without arguments": {
			source: `
				class ::Std::Object
					def bar: Symbol
						:baz
					end
				end

				self.bar
			`,
			wantStackTop: value.SymbolTable.Add("baz"),
			teardown: func() {
				delete(value.ObjectClass.Methods, value.SymbolTable.Add("bar"))
			},
		},
		"call an instance method with positional arguments": {
			source: `
				class ::Std::Object
					def add(a: Int, b: Int): Int
						a + b
					end
				end

				self.add 1, 8
			`,
			wantStackTop: value.SmallInt(9),
			teardown: func() {
				delete(value.ObjectClass.Methods, value.SymbolTable.Add("add"))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_DefineMixin(t *testing.T) {
	tests := sourceTestTable{
		"anonymous mixin without a body": {
			source:       "mixin; end",
			wantStackTop: value.NewMixin(),
		},
		"mixin without a body with a relative name": {
			source: "mixin Foo; end",
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"mixin without a body with an absolute name": {
			source: "mixin ::Foo; end",
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"mixin with a body": {
			source: `
				mixin Foo
					a := 5
					Bar := a - 2
				end
			`,
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithParent(value.MixinClass),
					),
				),
				value.MixinWithName("Foo"),
				value.MixinWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("Bar"): value.SmallInt(3),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"anonymous mixin with a body": {
			source: `
				mixin
					a := 5
					Bar := a - 2
				end
			`,
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithParent(value.MixinClass),
					),
				),
				value.MixinWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("Bar"): value.SmallInt(3),
					},
				),
			),
		},
		"nested mixins": {
			source: `
				mixin Gdańsk
					mixin Gdynia
						mixin Sopot
							Trójmiasto := "jest super"
							::Gdańsk::Warszawa := "to stolica"
						end
					end
				end
			`,
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Gdańsk"),
				value.MixinWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("Gdynia"): value.NewMixinWithOptions(
							value.MixinWithName("Gdańsk::Gdynia"),
							value.MixinWithConstants(
								value.SimpleSymbolMap{
									value.SymbolTable.Add("Sopot"): value.NewMixinWithOptions(
										value.MixinWithName("Gdańsk::Gdynia::Sopot"),
										value.MixinWithConstants(
											value.SimpleSymbolMap{
												value.SymbolTable.Add("Trójmiasto"): value.String("jest super"),
											},
										),
									),
								},
							),
						),
						value.SymbolTable.Add("Warszawa"): value.String("to stolica"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Gdańsk") },
		},
		"open an existing mixin": {
			source: `
				mixin Foo
					FIRST_CONSTANT := "oguem"
				end

				mixin Foo
					SECOND_CONSTANT := "całe te"
				end
			`,
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("FIRST_CONSTANT"):  value.String("oguem"),
						value.SymbolTable.Add("SECOND_CONSTANT"): value.String("całe te"),
					},
				),
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"redefined constant": {
			source: `
				Foo := 3
				mixin Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.RedefinedConstantErrorClass,
				"module Root already has a constant named `:Foo`",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
		"redefined class as mixin": {
			source: `
				class Foo; end
				mixin Foo; end
			`,
			wantRuntimeErr: value.NewError(
				value.RedefinedConstantErrorClass,
				"module Root already has a constant named `:Foo`",
			),
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Include(t *testing.T) {
	tests := sourceTestTable{
		"include a mixin to a class": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				class ::Std::Object
					include ::Foo
				end

				self.foo
			`,
			wantStackTop: value.String("hey, it's foo"),
			teardown:     func() { value.ObjectClass.Parent = value.PrimitiveObjectClass },
		},
		"include two mixins to a class": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				mixin Bar
					def bar: String
						"hey, it's bar"
					end
				end

				class ::Std::Object
					include ::Foo, ::Bar
				end

				self.foo + "; " + self.bar
			`,
			wantStackTop: value.String("hey, it's foo; hey, it's bar"),
			teardown:     func() { value.ObjectClass.Parent = value.PrimitiveObjectClass },
		},
		"include a complex mixin in a class": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				mixin Bar
					include ::Foo

					def bar: String
						"hey, it's bar"
					end
				end

				class ::Std::Int
					include ::Bar
				end

				1.foo + "; " + 1.bar
			`,
			wantStackTop: value.String("hey, it's foo; hey, it's bar"),
			teardown:     func() { value.ObjectClass.Parent = value.PrimitiveObjectClass },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Extend(t *testing.T) {
	tests := sourceTestTable{
		"extend a class with a mixin": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				class ::Std::String
					extend ::Foo
				end

				::Std::String.foo
			`,
			wantStackTop: value.String("hey, it's foo"),
			teardown:     func() { value.StringClass.SetDirectClass(value.ObjectClass) },
		},
		"extend a module with a mixin": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				module Std
					extend ::Foo
				end

				::Std.foo
			`,
			wantStackTop: value.String("hey, it's foo"),
			teardown:     func() { value.StdModule.SetDirectClass(value.ModuleClass) },
		},
		"extend a mixin with a mixin": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				mixin Bar
					extend ::Foo
				end

				::Bar.foo
			`,
			wantStackTop: value.String("hey, it's foo"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("Bar")
			},
		},
		"extend a class with two mixins": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				mixin Bar
					def bar: String
						"hey, it's bar"
					end
				end

				class ::Std::String
					extend ::Foo, ::Bar
				end

				::Std::String.foo + "; " + ::Std::String.bar
			`,
			wantStackTop: value.String("hey, it's foo; hey, it's bar"),
			teardown:     func() { value.StringClass.SetDirectClass(value.ObjectClass) },
		},
		"extend a class with a complex mixin": {
			source: `
				mixin Foo
					def foo: String
						"hey, it's foo"
					end
				end

				mixin Bar
					include ::Foo

					def bar: String
						"hey, it's bar"
					end
				end

				class Baz
					extend ::Bar
				end

				::Baz.foo + "; " + ::Baz.bar
			`,
			wantStackTop: value.String("hey, it's foo; hey, it's bar"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
				value.RootModule.Constants.DeleteString("Bar")
				value.RootModule.Constants.DeleteString("Baz")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
