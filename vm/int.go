package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

func init() {
	DefineMethodWithOptions(
		value.IntClass.Methods,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			other := args[1]
			switch s := self.(type) {
			case value.SmallInt:
				return value.ToValueErr(s.Add(other))
			case *value.BigInt:
				return value.ToValueErr(s.Add(other))
			}

			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
		NativeMethodWithStringParameters("other"),
		NativeMethodWithFrozen(),
	)
	DefineMethodWithOptions(
		value.IntClass.Methods,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			switch s := self.(type) {
			case value.SmallInt:
				return value.String(s.Inspect()), nil
			case *value.BigInt:
				return value.String(s.Inspect()), nil
			}
			panic(fmt.Sprintf("expected SmallInt or BigInt, got: %#v", self))
		},
	)

}
