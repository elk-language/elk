package value

var SyncHashSetClass *Class         // ::Std::Sync::HashSet
var SyncHashSetIteratorClass *Class // ::Std::Sync::HashSet::Iterator

func initSyncHashSet() {
	SyncHashSetClass = NewClass()
	SyncHashSetClass.IncludeMixin(SetMixin)
	SyncModule.AddConstantString("HashSet", Ref(SyncHashSetClass))
	RegisterNativeClass("Std::Sync::HashSet", "value.SyncHashSetClass")

	SyncHashSetIteratorClass = NewClass()
	SyncHashSetClass.AddConstantString("Iterator", Ref(SyncHashSetIteratorClass))
	RegisterNativeClass("Std::Sync::HashSet::Iterator", "value.SyncHashSetIteratorClass")
}
