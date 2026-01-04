package value

// ::Std::ImmutableCollection
var ImmutableCollectionInterface *Interface

// ::Std::ImmutableCollection::Base
var ImmutableCollectionBaseMixin *Mixin

func initImmutableCollection() {
	ImmutableCollectionInterface = NewInterface()
	StdModule.AddConstantString("ImmutableCollection", Ref(ImmutableCollectionInterface))
	RegisterNativeInterface("Std::ImmutableCollection", "value.ImmutableCollectionInterface")

	ImmutableCollectionBaseMixin = NewMixin()
	ImmutableCollectionBaseMixin.IncludeMixin(IterableFiniteBaseMixin)
	ImmutableCollectionInterface.AddConstantString("Base", Ref(ImmutableCollectionBaseMixin))
	RegisterNativeMixin("Std::ImmutableCollection::Base", "value.ImmutableCollectionBaseMixin")
}
