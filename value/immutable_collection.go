package value

// ::Std::ImmutableCollection
var ImmutableCollectionInterface *Interface

// ::Std::ImmutableCollection::Base
var ImmutableCollectionBaseMixin *Mixin

func initImmutableCollection() {
	ImmutableCollectionInterface = NewInterface()
	StdModule.AddConstantString("ImmutableCollection", Ref(ImmutableCollectionInterface))

	ImmutableCollectionBaseMixin = NewMixin()
	ImmutableCollectionBaseMixin.IncludeMixin(IterableFiniteBase)
	ImmutableCollectionInterface.AddConstantString("Base", Ref(ImmutableCollectionBaseMixin))
}
