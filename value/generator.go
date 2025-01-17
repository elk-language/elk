package value

var GeneratorClass *Class // ::Std::Generator

func initGenerator() {
	GeneratorClass = NewClass()
	StdModule.AddConstantString("Generator", Ref(GeneratorClass))
}
