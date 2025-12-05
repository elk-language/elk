package value

var TupleMixin *Mixin // ::Std::Tuple

func initTuple() {
	TupleMixin = NewMixin()
	TupleMixin.IncludeMixin(ImmutableCollectionBaseMixin)
	StdModule.AddConstantString("Tuple", Ref(TupleMixin))
}
