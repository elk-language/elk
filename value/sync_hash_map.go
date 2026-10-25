package value

var SyncHashMapClass *Class         // ::Std::Sync::HashMap
var SyncHashMapIteratorClass *Class // ::Std::Sync::HashMap::Iterator

func initSyncHashMap() {
	SyncHashMapClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SyncHashMapClass.IncludeMixin(MapMixin)
	SyncModule.AddConstantString("HashMap", Ref(SyncHashMapClass))
	RegisterNativeClass("Std::Sync::HashMap", "value.SyncHashMapClass")

	SyncHashMapIteratorClass = NewClass()
	SyncHashMapClass.AddConstantString("Iterator", Ref(SyncHashMapIteratorClass))
	RegisterNativeClass("Std::Sync::HashMap::Iterator", "value.SyncHashMapIteratorClass")
}
