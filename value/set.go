package value

var SetMixin *Mixin // ::Std::Set

func initSet() {
	SetMixin = NewMixin()
	StdModule.AddConstantString("Set", SetMixin)
}
