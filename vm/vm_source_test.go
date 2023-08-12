package vm

import (
	"strings"
	"testing"

	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/object"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Represents a single VM source code test case.
type sourceTestCase struct {
	source         string
	wantStackTop   object.Value
	wantStdout     string
	wantRuntimeErr object.Value
	wantCompileErr errors.ErrorList
}

// Type of the compiler test table.
type sourceTestTable map[string]sourceTestCase

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
		cmp.AllowUnexported(object.Error{}, object.BigFloat{}, object.BigInt{}),
		cmpopts.IgnoreUnexported(object.Class{}),
		cmpopts.IgnoreFields(object.Class{}, "ConstructorFunc"),
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
		t.Fatalf(diff)
	}
}

func TestVMSourceLocals(t *testing.T) {
	tests := sourceTestTable{
		"define and initialise a variable": {
			source:       "var a = 'foo'",
			wantStackTop: object.String("foo"),
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
			wantStackTop: object.SmallInt(18),
		},
		"define and set a variable": {
			source: `
				var a = 'foo'
				a = a + ' bar'
				a
			`,
			wantStackTop: object.String("foo bar"),
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

func TestVMSourceIfExpressions(t *testing.T) {
	tests := sourceTestTable{
		"empty then truthy": {
			source:       "if true; end",
			wantStackTop: object.Nil,
		},
		"empty then falsy": {
			source:       "if false; end",
			wantStackTop: object.Nil,
		},
		"execute the then branch": {
			source: `
				a := 5
				if a
					a = a + 2
				end
			`,
			wantStackTop: object.SmallInt(7),
		},
		"execute the empty else branch": {
			source: `
				a := 5
				if false
					a = a * 2
				end
			`,
			wantStackTop: object.Nil,
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
			wantStackTop: object.SmallInt(7),
		},
		"execute the else branch instead of then": {
			source: `
				a := 5
				if false
					a = a + 2
				else
					a = 30
				end
			`,
			wantStackTop: object.SmallInt(30),
		},
		"if is an expression": {
			source: `
				a := 5
				b := if a
					"foo"
				else
					5
				end
				b
			`,
			wantStackTop: object.String("foo"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
