package test

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

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
			argFn := args[2]

			prevSuite := CurrentSuite
			CurrentSuite = CurrentSuite.NewSubSuite(string(argName))

			_, err = v.CallCallable(argFn)
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
			argFn := args[2]

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
			argFn := args[2]

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
			argFn := args[2]

			CurrentSuite.NewCase(fmt.Sprintf("should %s", string(argName)), argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"before_each",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argFn := args[2]
			CurrentSuite.RegisterBeforeEach(argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"before_all",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argFn := args[2]
			CurrentSuite.RegisterBeforeAll(argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"after_each",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argFn := args[2]
			CurrentSuite.RegisterAfterEach(argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"after_all",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argFn := args[2]
			CurrentSuite.RegisterAfterAll(argFn)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(1),
	)

	return testModule
}
