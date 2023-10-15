package vm

import (
	"strings"
	"testing"

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

	opts := []cmp.Option{
		cmp.AllowUnexported(value.Error{}, value.BigInt{}),
		cmpopts.IgnoreUnexported(value.Class{}),
		cmpopts.IgnoreFields(value.Class{}, "ConstructorFunc"),
	}

	chunk, gotCompileErr := compiler.CompileSource(testFileName, tc.source)
	if gotCompileErr != nil {
		if diff := cmp.Diff(tc.wantCompileErr, gotCompileErr, opts...); diff != "" {
			t.Fatalf(diff)
		}
		return
	}
	var stdout strings.Builder
	vm := New(WithStdout(&stdout))
	gotStackTop, gotRuntimeErr := vm.InterpretBytecode(chunk)
	gotStdout := stdout.String()
	if diff := cmp.Diff(tc.wantRuntimeErr, gotRuntimeErr, opts...); diff != "" {
		t.Fatalf(diff)
	}
	if diff := cmp.Diff(tc.wantStdout, gotStdout, opts...); diff != "" {
		t.Fatalf(diff)
	}
	if diff := cmp.Diff(tc.wantStackTop, gotStackTop, opts...); diff != "" {
		t.Logf("got: %s, want: %s", gotStackTop.Inspect(), tc.wantStackTop.Inspect())
		t.Fatalf(diff)
	}
}

func vmSimpleSourceTest(source string, want value.Value, t *testing.T) {
	t.Helper()

	opts := []cmp.Option{
		cmp.AllowUnexported(value.Error{}, value.BigFloat{}, value.BigInt{}),
		cmpopts.IgnoreUnexported(value.Class{}),
		cmpopts.IgnoreFields(value.Class{}, "ConstructorFunc"),
	}

	chunk, gotCompileErr := compiler.CompileSource(testFileName, source)
	if gotCompileErr != nil {
		t.Fatalf("Compile Error: %s", gotCompileErr.Error())
		return
	}
	var stdout strings.Builder
	vm := New(WithStdout(&stdout))
	got, gotRuntimeErr := vm.InterpretBytecode(chunk)
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

func TestVMSource_Equal(t *testing.T) {
	tests := simpleSourceTestTable{
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

func TestVMSource_StrictEqual(t *testing.T) {
	tests := simpleSourceTestTable{
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
