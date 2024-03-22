package value

var RangeMixin *Mixin // ::Std::Range

func initRange() {
	RangeMixin = NewMixin()
	StdModule.AddConstantString("Range", RangeMixin)
}
