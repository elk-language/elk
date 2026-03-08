package compiler_test

import (
	"bytes"
	"go/format"
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/compiler/colorize"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp/v3"
)

// Represents a single go compiler test case.
type goTestCase struct {
	input string
	want  string
	err   diagnostic.DiagnosticList
}

// Type of the compiler test table.
type goTestTable map[string]goTestCase

func goCompilerTest(tc goTestCase, t *testing.T) {
	t.Helper()

	pp.Default.SetColoringEnabled(false)

	var buff bytes.Buffer
	compiler, errDiag := checker.CheckSourceNative("<main>", tc.input, nil, &buff, vm.DefaultThreadPool)
	opts := comparer.Options()
	if diff := cmp.Diff(tc.err, errDiag, opts...); diff != "" {
		t.Log(pp.Sprint(errDiag))
		t.Log(diff)
		t.Fail()
	}
	if errDiag.IsFailure() {
		return
	}

	compiler.Flush()
	result, err := format.Source(buff.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(tc.want, string(result), opts...); diff != "" {
		t.Log(string(colorize.ColorizeWhen(result, true)))
		t.Log(diff)
		t.Fail()
	}
}
