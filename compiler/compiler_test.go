package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp/v3"
)

var classSymbol value.Symbol = value.ToSymbol("<class>")
var singletonClassSymbol value.Symbol = value.ToSymbol("<singleton_class>")
var moduleSymbol value.Symbol = value.ToSymbol("<module>")
var mixinSymbol value.Symbol = value.ToSymbol("<mixin>")
var mainSymbol value.Symbol = value.ToSymbol("<main>")
var functionSymbol value.Symbol = value.ToSymbol("<closure>")

// Represents a single compiler test case.
type testCase struct {
	input string
	want  *vm.BytecodeFunction
	err   error.ErrorList
}

// Type of the compiler test table.
type testTable map[string]testCase

func compilerTest(tc testCase, t *testing.T) {
	t.Helper()

	pp.Default.SetColoringEnabled(false)

	typechecker := checker.New()
	got, err := typechecker.CheckSource("<main>", tc.input)
	opts := comparer.Options()
	if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
		t.Log(pp.Sprint(err))
		t.Fatal(diff)
	}
	if err.IsFailure() {
		return
	}
	if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
		t.Log(got.DisassembleString())
		t.Fatal(diff)
	}
}

// Create a new position in tests
var P = position.New

// Create a new span in tests
var S = position.NewSpan

const testFileName = "<main>"

// Create a new source location in tests.
// Create a new location in tests
func L(startPos, endPos *position.Position) *position.Location {

	return position.NewLocation(testFileName, startPos, endPos)
}
