package vm_test

import (
	"strings"
	"testing"

	"github.com/elk-language/elk"
	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestBreakpointInGoReadsCapturedVariables(t *testing.T) {
	prev := vm.BREAKPOINT_HANDLER
	t.Cleanup(func() { vm.BREAKPOINT_HANDLER = prev })

	tests := []struct {
		name   string
		source string
		expr   string
		want   string
		stdout string
	}{
		{
			name: "local in the main thread",
			source: `
				a := 5
				breakpoint
				println a
			`,
			expr:   "a",
			want:   "5",
			stdout: "5\n",
		},
		{
			name: "captured variable in a go thread",
			source: `
				wg := Sync::WaitGroup(1)
				a := 5
				go
					c := a + 1
					breakpoint
					println c
					wg.end
				end
				wg.wait
			`,
			expr:   "a",
			want:   "5",
			stdout: "6\n",
		},
		{
			name: "variable captured after the breakpoint",
			source: `
				wg := Sync::WaitGroup(1)
				a := 5
				go
					breakpoint
					println a
					wg.end
				end
				wg.wait
			`,
			expr:   "a",
			want:   "5",
			stdout: "5\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			elk.InitGlobalEnvironment()
			var got string
			var gotErr string
			vm.BREAKPOINT_HANDLER = breakpointFunc(func(thread *vm.Thread, ctx value.Value) {
				context := ctx.AsReference().(*compiler.BytecodeBreakpointContext)
				typechecker := checker.NewBreakpointChecker(context)
				fn, dl := typechecker.CheckBreakpointSource("<bp>", tc.expr)
				if dl != nil && dl.IsFailure() {
					gotErr = dl.Error()
					return
				}
				val, runtimeErr := thread.InterpretBreakpoint(fn)
				if !runtimeErr.IsUndefined() {
					gotErr = runtimeErr.Inspect()
					return
				}
				got = val.Inspect()
			})

			stdout := &strings.Builder{}
			typechecker := checker.New()
			chunk, compileErr := typechecker.CheckSourceBytecode("sourceName", tc.source, nil)
			if compileErr.IsFailure() {
				t.Fatalf("compile: %v", compileErr)
			}
			v := vm.New(vm.WithStdout(stdout))
			_, runtimeErr := v.InterpretTopLevel(chunk)
			if !runtimeErr.IsUndefined() {
				t.Fatalf("runtime: %s", runtimeErr.Inspect())
			}
			if gotErr != "" {
				t.Fatalf("breakpoint eval: %s", gotErr)
			}
			if got != tc.want {
				t.Fatalf("breakpoint result: got %#v, want %#v", got, tc.want)
			}
			if stdout.String() != tc.stdout {
				t.Fatalf("stdout: got %#v, want %#v", stdout.String(), tc.stdout)
			}
		})
	}
}

type breakpointFunc func(thread *vm.Thread, ctx value.Value)

func (f breakpointFunc) RunBreakpoint(thread *vm.Thread, ctx value.Value) {
	f(thread, ctx)
}
