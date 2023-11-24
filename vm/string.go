package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

func init() {
	DefineMethodRestParam(
		value.ObjectClass.Methods,
		"print",
		[]string{"values"},
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].(value.List)
			for _, val := range values {
				fmt.Print(val)
			}

			return value.Nil, nil
		},
	)
	DefineMethodRestParam(
		value.ObjectClass.Methods,
		"println",
		[]string{"values"},
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			values := args[1].(value.List)
			for _, val := range values {
				fmt.Println(val)
			}

			return value.Nil, nil
		},
	)
	DefineMethodNoParams(
		value.ObjectClass.Methods,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.String(self.Inspect()), nil
		},
	)
	DefineMethodNoParams(
		value.ObjectClass.Methods,
		"class",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return self.Class(), nil
		},
	)

	DefineMethodReqParams(
		value.ObjectClass.Methods,
		"==",
		[]string{"other"},
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			return value.ToElkBool(self == other), nil
		},
	)
	value.ObjectClass.DefineAliasString("===", "==")

}
