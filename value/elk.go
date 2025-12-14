package value

import "github.com/elk-language/elk/info"

var ElkModule *Module // Std::Elk

func initElk() {
	ElkModule = NewModule()
	StdModule.AddConstantString("Elk", Ref(ElkModule))
	ElkModule.AddConstantString("VERSION", Ref(String(info.Version)))
}
