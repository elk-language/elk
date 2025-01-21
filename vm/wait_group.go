package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::WaitGroup
func initWaitGroup() {
	// Instance methods
	c := &value.WaitGroupClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.WaitGroup)(args[0].Pointer())
			nVal := args[1]
			if nVal.IsUndefined() {
				return args[0], value.Undefined
			}
			if nVal.IsReference() {
				return value.Undefined, value.Ref(value.NewError(value.OutOfRangeErrorClass, "n is too large"))
			}
			self.Add(int(nVal.AsSmallInt()))
			return args[0], value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"add",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.WaitGroup)(args[0].Pointer())
			nVal := args[1]
			if nVal.IsReference() {
				return value.Undefined, value.Ref(value.NewError(value.OutOfRangeErrorClass, "n is too large"))
			}
			self.Add(int(nVal.AsSmallInt()))
			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"remove",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.WaitGroup)(args[0].Pointer())
			nVal := args[1]
			if nVal.IsReference() {
				return value.Undefined, value.Ref(value.NewError(value.OutOfRangeErrorClass, "n is too large"))
			}
			self.Remove(int(nVal.AsSmallInt()))
			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"start",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.WaitGroup)(args[0].Pointer())
			self.Start()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"end",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.WaitGroup)(args[0].Pointer())
			self.End()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"wait",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.WaitGroup)(args[0].Pointer())
			self.Wait()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.WaitGroup)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

}
