package value

var ComparableMixin *Mixin // ::Std::Comparable

func initComparable() {
	ComparableMixin = NewMixin()
	StdModule.AddConstantString("Comparable", Ref(ComparableMixin))
	RegisterNativeMixin("Std::Comparable", "value.ComparableMixin")
}
