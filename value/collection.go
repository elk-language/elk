package value

// ::Std::Collection
var CollectionInterface *Interface

// ::Std::Collection::Base
var CollectionBaseMixin *Mixin

func initCollection() {
	CollectionInterface = NewInterface()
	StdModule.AddConstantString("Collection", Ref(CollectionInterface))
	RegisterNativeInterface("Std::Collection", "value.CollectionInterface")

	CollectionBaseMixin = NewMixin()
	CollectionBaseMixin.IncludeMixin(ImmutableCollectionBaseMixin)
	CollectionInterface.AddConstantString("Base", Ref(CollectionBaseMixin))
	RegisterNativeMixin("Std::Collection::Base", "value.CollectionBaseMixin")
}
