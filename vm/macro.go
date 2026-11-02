package vm

import (
	"github.com/elk-language/elk/value"
)

var macroSTORE = NewSyncNativeKeyHashMap[value.String](5)

func MacroSTORE() *SyncNativeKeyHashMap[value.String] {
	return macroSTORE
}

func SetMacroSTORE(m *SyncNativeKeyHashMap[value.String]) {
	macroSTORE = m
	value.MacroModule.AddConstantString("STORE", m.ToValue())
}

// Std::Macro
func initMacro() {
	value.MacroModule.AddConstantString("STORE", MacroSTORE().ToValue())
	value.RegisterNativeConstant("Std::Macro::STORE", "vm.MacroSTORE()", value.FetchGoType("*vm.HashMapOfValue"))
}
