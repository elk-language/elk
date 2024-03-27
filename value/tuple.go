package value

var TupleMixin *Mixin // ::Std::Tuple

func initTuple() {
	TupleMixin = NewMixin()
	StdModule.AddConstantString("Tuple", TupleMixin)
}
