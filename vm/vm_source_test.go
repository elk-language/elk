package vm_test

import (
	"strings"
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp"
)

// Represents a single VM source code test case.
type sourceTestCase struct {
	source         string
	wantStackTop   value.Value
	wantStdout     string
	wantRuntimeErr value.Value
	wantCompileErr error.ErrorList
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

	typechecker := checker.New()
	chunk, gotCompileErr := typechecker.CheckSource(testFileName, tc.source)
	if diff := cmp.Diff(tc.wantCompileErr, gotCompileErr, comparer.Options()...); diff != "" {
		t.Log(pp.Sprint(gotCompileErr))
		t.Fatal(diff)
	}

	var stdout strings.Builder
	v := vm.New(vm.WithStdout(&stdout))
	gotStackTop, gotRuntimeErr := v.InterpretTopLevel(chunk)
	gotStdout := stdout.String()
	if tc.teardown != nil {
		tc.teardown()
	}
	if diff := cmp.Diff(tc.wantRuntimeErr, gotRuntimeErr, comparer.Options()...); diff != "" {
		t.Log(pp.Sprint(gotRuntimeErr))
		t.Fatal(diff)
	}
	if tc.wantRuntimeErr != nil {
		return
	}
	if diff := cmp.Diff(tc.wantStdout, gotStdout, comparer.Options()...); diff != "" {
		t.Fatal(diff)
	}
	if diff := cmp.Diff(tc.wantStackTop, gotStackTop, comparer.Options()...); diff != "" {
		t.Log(gotRuntimeErr)
		if gotStackTop != nil && tc.wantStackTop != nil {
			t.Logf("got: %s, want: %s", gotStackTop.Inspect(), tc.wantStackTop.Inspect())
		}
		t.Fatal(diff)
	}
}

func vmSimpleSourceTest(source string, want value.Value, t *testing.T) {
	t.Helper()

	opts := comparer.Options()

	chunk, gotCompileErr := compiler.CompileSource(testFileName, source)
	if gotCompileErr != nil {
		t.Fatalf("Compile Error: %s", gotCompileErr.Error())
		return
	}
	var stdout strings.Builder
	vm := vm.New(vm.WithStdout(&stdout))
	got, gotRuntimeErr := vm.InterpretTopLevel(chunk)
	if gotRuntimeErr != nil {
		t.Fatalf("Runtime Error: %s", gotRuntimeErr.Inspect())
	}
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Logf("got: %s, want: %s", got.Inspect(), want.Inspect())
		t.Fatalf(diff)
	}
}
