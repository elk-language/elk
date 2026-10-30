package value

var SyncOrderedHashSetClass *Class         // ::Std::Sync::OrderedHashSet
var SyncOrderedHashSetIteratorClass *Class // ::Std::Sync::OrderedHashSet::Iterator

func initSyncOrderedHashSet() {
	SyncOrderedHashSetClass = NewClass()
	SyncOrderedHashSetClass.IncludeMixin(SetMixin)
	SyncModule.AddConstantString("OrderedHashSet", Ref(SyncOrderedHashSetClass))
	RegisterNativeClass("Std::Sync::OrderedHashSet", "value.SyncOrderedHashSetClass")

	SyncOrderedHashSetIteratorClass = NewClass()
	SyncOrderedHashSetClass.AddConstantString("Iterator", Ref(SyncOrderedHashSetIteratorClass))
	RegisterNativeClass("Std::Sync::OrderedHashSet::Iterator", "value.SyncOrderedHashSetIteratorClass")
}
