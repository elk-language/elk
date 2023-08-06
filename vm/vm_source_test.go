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

// Create a new source position in tests.
var P = position.New

// Create a new source location in tests.
func L(start, length, line, column int) *position.Location {
	return position.NewLocation(testFileName, start, length, line, column)
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

func TestVMSource(t *testing.T) {
	tests := sourceTestTable{
		"add two ints": {
			source:       "1 + 2",
			wantStackTop: object.SmallInt(3),
		},
		"empty source": {
			source:       "",
			wantStackTop: object.Nil,
		},
		"define and initialise a variable": {
			source:       "var a = 'foo'",
			wantStackTop: object.String("foo"),
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
				errors.NewError(L(15, 1, 3, 5), "can't access an uninitialised local: a"),
			},
		},
		"try to read a nonexistent variable": {
			source: `
				a
			`,
			wantCompileErr: errors.ErrorList{
				errors.NewError(L(5, 1, 2, 5), "undeclared variable: a"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
