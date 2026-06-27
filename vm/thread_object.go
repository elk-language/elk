package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

func (*Thread) Class() *value.Class {
	return value.ThreadClass
}

func (*Thread) DirectClass() *value.Class {
	return value.ThreadClass
}

func (*Thread) SingletonClass() *value.Class {
	return nil
}

func (vm *Thread) Inspect() string {
	return fmt.Sprintf(`Std::Thread{state: %s}`, stateSymbols[vm.state].Inspect())
}

func (vm *Thread) Error() string {
	return vm.Inspect()
}

func (vm *Thread) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (vm *Thread) Copy() value.Reference {
	return vm
}

func (vm *Thread) ToValue() value.Value {
	return value.Ref(vm)
}

func (vm *Thread) StateSymbol() value.Symbol {
	return stateSymbols[vm.state]
}

// Std::Thread
func initThread() {
	// Instance methods
	c := &value.ThreadClass.MethodContainer

	Def(
		c,
		"==",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return value.BoolVal(args[0] == args[1]), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"state",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*Thread)(args[0].Pointer())
			return self.StateSymbol().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*Thread)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)
	Def(
		c,
		"copy",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*Thread)(args[0].Pointer())
			return value.Ref(self.Copy()), value.Undefined
		},
	)

}
