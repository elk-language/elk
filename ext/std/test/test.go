package test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/fatih/color"
)

type TestStatus uint8

const (
	TEST_PENDING TestStatus = iota
	TEST_FAILED
	TEST_ERROR
	TEST_SKIPPED
	TEST_RUNNING
	TEST_SUCCESS
)

func Run() *SuiteReport {
	v := vm.New()
	events := make(chan *ReportEvent, 50)

	var reporter Reporter
	if color.NoColor {
		reporter = NewPlainReporter()
	} else {
		reporter = NewRichReporter()
	}

	seed := uint64(time.Now().UnixNano())
	return RunWith(v, reporter, events, seed)
}

// Run the all tests under the root suite
func RunWith(v *vm.VM, reporter Reporter, events chan *ReportEvent, seed uint64) *SuiteReport {
	var wg sync.WaitGroup
	wg.Add(1)
	shutdownCtx, shutdown := context.WithCancel(context.Background())

	go func() {
		reporter.Report(events, shutdown)
		wg.Done()
	}()

	rng := rand.New(rand.NewPCG(seed, ^seed+1))
	report := RootSuite.Run(v, events, rng, shutdownCtx)
	close(events)
	wg.Wait()

	return report
}

var RootSuite = NewSuite("", nil)
var CurrentSuite = RootSuite

func initTest() *value.Module {
	testModule := value.NewModule()
	value.StdModule.AddConstantString("Test", value.Ref(testModule))

	c := &testModule.SingletonClass().MethodContainer
	vm.Def(
		c,
		"describe",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argName := args[1].AsReference().(value.String)
			argFn := args[2].AsReference().(*vm.Closure)

			prevSuite := CurrentSuite
			CurrentSuite = CurrentSuite.NewSubSuite(string(argName))

			_, err = v.CallClosure(argFn)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			CurrentSuite = prevSuite
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(2),
	)
	vm.Alias(c, "context", "describe")

	vm.Def(
		c,
		"test",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argName := args[1].AsReference().(value.String)
			argFn := args[2].AsReference().(*vm.Closure)

			CurrentSuite.NewCase(string(argName), argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"it",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argName := args[1].AsReference().(value.String)
			argFn := args[2].AsReference().(*vm.Closure)

			CurrentSuite.NewCase(fmt.Sprintf("it %s", string(argName)), argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"should",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argName := args[1].AsReference().(value.String)
			argFn := args[2].AsReference().(*vm.Closure)

			CurrentSuite.NewCase(fmt.Sprintf("should %s", string(argName)), argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"before_each",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argFn := args[1].AsReference().(*vm.Closure)
			CurrentSuite.RegisterBeforeEach(argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"before_all",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argFn := args[1].AsReference().(*vm.Closure)
			CurrentSuite.RegisterBeforeAll(argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"after_each",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argFn := args[1].AsReference().(*vm.Closure)
			CurrentSuite.RegisterAfterEach(argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"after_all",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argFn := args[1].AsReference().(*vm.Closure)
			CurrentSuite.RegisterAfterAll(argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(1),
	)

	return testModule
}
