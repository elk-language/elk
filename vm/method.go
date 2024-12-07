package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// Std::Method
func initMethod() {
	// Instance methods
	c := &value.MethodClass.MethodContainer
	nativeSymbol := value.ToSymbol("native")
	bytecodeSymbol := value.ToSymbol("bytecode")
	getterSymbol := value.ToSymbol("getter")
	setterSymbol := value.ToSymbol("setter")
	Def(
		c,
		"type",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]

			switch self.SafeAsReference().(type) {
			case *NativeMethod:
				return nativeSymbol.ToValue(), value.Nil
			case *BytecodeFunction:
				return bytecodeSymbol.ToValue(), value.Nil
			case *GetterMethod:
				return getterSymbol.ToValue(), value.Nil
			case *SetterMethod:
				return setterSymbol.ToValue(), value.Nil
			default:
				panic(fmt.Sprintf("invalid method type: %#v", self))
			}
		},
	)
}
