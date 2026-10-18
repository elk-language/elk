package vm

import (
	"github.com/elk-language/elk/value"
)

var macroSTORE = NewHashMapOfValue(5)

func MacroSTORE() *HashMapOfValue {
	return macroSTORE
}

func SetMacroSTORE(m *HashMapOfValue) {
	macroSTORE = m
	value.MacroModule.AddConstantString("STORE", m.ToValue())
}

// Std::Macro
func initMacro() {
	value.MacroModule.AddConstantString("STORE", MacroSTORE().ToValue())
	value.RegisterNativeConstant("Std::Macro::STORE", "vm.MacroSTORE()", value.FetchGoType("*vm.HashMapOfValue"))
}
