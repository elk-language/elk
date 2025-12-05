package value

// ::Std::Collection
var CollectionInterface *Interface

// ::Std::Collection::Base
var CollectionBaseMixin *Mixin

func initCollection() {
	CollectionInterface = NewInterface()
	StdModule.AddConstantString("Collection", Ref(CollectionInterface))

	CollectionBaseMixin = NewMixin()
	CollectionBaseMixin.IncludeMixin(ImmutableCollectionBaseMixin)
	CollectionInterface.AddConstantString("Base", Ref(CollectionBaseMixin))
}
