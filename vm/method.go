package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// Std::Method
func initMethod() {
	// Instance methods
	c := &value.MethodClass.MethodContainer
	Def(
		c,
		"doc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			var docValue value.Value

			switch s := self.(type) {
			case *NativeMethod:
				docValue = s.Doc
			case *BytecodeFunction:
				docValue = s.Doc
			case *GetterMethod:
				docValue = s.Doc
			case *SetterMethod:
				docValue = s.Doc
			default:
				panic(fmt.Sprintf("invalid method type: %#v", self))
			}

			if docValue == nil {
				return value.Nil, nil
			}

			return docValue, nil
		},
	)
	Def(
		c,
		"doc=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			val := args[1]
			switch val.(type) {
			case value.NilType, value.String:
			default:
				return nil, value.NewArgumentTypeError("val", val.Inspect(), "Std::String")
			}

			switch s := self.(type) {
			case *NativeMethod:
				s.Doc = val
			case *BytecodeFunction:
				s.Doc = val
			case *GetterMethod:
				s.Doc = val
			case *SetterMethod:
				s.Doc = val
			default:
				panic(fmt.Sprintf("invalid method type: %#v", self))
			}

			return val, nil
		},
		DefWithParameters(1),
	)
	nativeSymbol := value.ToSymbol("native")
	bytecodeSymbol := value.ToSymbol("bytecode")
	getterSymbol := value.ToSymbol("getter")
	setterSymbol := value.ToSymbol("setter")
	Def(
		c,
		"type",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]

			switch self.(type) {
			case *NativeMethod:
				return nativeSymbol, nil
			case *BytecodeFunction:
				return bytecodeSymbol, nil
			case *GetterMethod:
				return getterSymbol, nil
			case *SetterMethod:
				return setterSymbol, nil
			default:
				panic(fmt.Sprintf("invalid method type: %#v", self))
			}
		},
	)
}
