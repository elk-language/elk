package test

import (
	"fmt"

	"github.com/elk-language/elk/ext"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func Init() {
	ext.Register("test", runtimeInit, typecheckerInit)
}

func runtimeInit() {
	testModule := value.NewModule()
	value.StdModule.AddConstantString("Test", value.Ref(testModule))

	assertionsMixin := value.NewMixin()
	testModule.AddConstantString("Assertions", value.Ref(assertionsMixin))

	assertionErrorClass := value.NewClassWithOptions(
		value.ClassWithSuperclass(value.ErrorClass),
	)

	c := &assertionsMixin.MethodContainer
	vm.Def(
		c,
		"assert_truthy",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argGot := args[1]

			if value.Truthy(argGot) {
				return value.Nil, value.Undefined
			}

			gotInspect, err := vm.Inspect(v, argGot)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			var message string
			if args[2].IsUndefined() {
				message = fmt.Sprintf("value `%s` is not truthy", gotInspect)
			} else {
				message = args[2].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"assert_falsy",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argGot := args[1]

			if value.Falsy(argGot) {
				return value.Nil, value.Undefined
			}

			gotInspect, err := vm.Inspect(v, argGot)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			var message string
			if args[2].IsUndefined() {
				message = fmt.Sprintf("value `%s` is not falsy", gotInspect)
			} else {
				message = args[2].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(2),
	)
}

// TODO: test macros
func typecheckerInit(checker types.Checker) {
	env := checker.Env()
	testModule := env.Root.DefineModule("", symbol.Test, env)
	assertionsMixin := testModule.DefineMixin("", false, value.ToSymbol("Assertions"), env)

	types.DefMacro(
		assertionsMixin,
		"",
		"assert!",
		[]*types.Parameter{
			types.NewParameter(
				value.ToSymbol("expression"),
			),
		},
		env.ExpressionNode(),
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {

		},
	)
}
