package compiler_test

import (
	"testing"

	"github.com/elk-language/elk"
	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp/v3"
)

func init() {
	elk.InitGlobalEnvironment()
}

var namespaceDefinitionsSymbol value.Symbol = value.ToSymbol("<namespaceDefinitions>")
var methodDefinitionsSymbol value.Symbol = value.ToSymbol("<methodDefinitions>")
var ivarIndicesSymbol value.Symbol = value.ToSymbol("<ivarIndices>")
var mainSymbol value.Symbol = value.ToSymbol("<main>")
var functionSymbol value.Symbol = value.ToSymbol("<closure>")

// Represents a single compiler test case.
type bytecodeTestCase struct {
	input  string
	want   *vm.BytecodeFunction
	wantFn func(bytecodeTestCase) *vm.BytecodeFunction
	err    diagnostic.DiagnosticList
}

// Type of the compiler test table.
type bytecodeTestTable map[string]bytecodeTestCase

func bytecodeCompilerTest(tc bytecodeTestCase, t *testing.T) {
	t.Helper()

	pp.Default.SetColoringEnabled(false)

	typechecker := checker.New()
	typechecker.SetBuiltinImportsProcessed(true)
	got, err := typechecker.CheckSourceBytecode("<main>", tc.input)
	opts := comparer.Options()
	if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
		t.Log(pp.Sprint(err))
		t.Log(diff)
		t.Fail()
	}
	if err.IsFailure() {
		return
	}

	var want *vm.BytecodeFunction
	if tc.wantFn != nil {
		want = tc.wantFn(tc)
	} else {
		want = tc.want
	}

	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Log(got.DisassembleString())
		t.Log(diff)
		t.Fail()
	}
}

const testFileName = "<main>"
