package value

var SyncOrderedHashMapClass *Class         // ::Std::Sync::OrderedHashMap
var SyncOrderedHashMapIteratorClass *Class // ::Std::Sync::OrderedHashMap::Iterator

func initSyncOrderedHashMap() {
	SyncOrderedHashMapClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SyncOrderedHashMapClass.IncludeMixin(RecordMixin)
	SyncModule.AddConstantString("OrderedHashMap", SyncOrderedHashMapClass.ToValue())
	RegisterNativeClass("Std::Sync::OrderedHashMap", "value.SyncOrderedHashMapClass")

	SyncOrderedHashMapIteratorClass = NewClass()
	SyncOrderedHashMapClass.AddConstantString("Iterator", SyncOrderedHashMapIteratorClass.ToValue())
	RegisterNativeClass("Std::OrderedHashMap::Iterator", "value.SyncOrderedHashMapIteratorClass")
}
