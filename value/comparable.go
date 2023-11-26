package value

var ComparableMixin *Mixin // ::Std::Comparable

func initComparable() {
	ComparableMixin = NewMixin()
	StdModule.AddConstantString("Comparable", ComparableMixin)
}
