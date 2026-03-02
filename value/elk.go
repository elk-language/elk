package value

import "github.com/elk-language/elk/info"

var ElkModule *Module // Std::Elk

func ElkVersion() String {
	return String(info.Version)
}

func initElk() {
	ElkModule = NewModule()
	StdModule.AddConstantString("Elk", Ref(ElkModule))
	RegisterNativeModule("Std::Elk", "value.ElkModule")

	ElkModule.AddConstantString("VERSION", Ref(ElkVersion()))
	RegisterNativeConstant("Std::Elk::VERSION", "value.ElkVersion()", "value.String")
}
