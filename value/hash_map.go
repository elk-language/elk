package value

var HashMapClass *Class         // ::Std::HashMap
var HashMapIteratorClass *Class // ::Std::HashMap::Iterator

func initHashMap() {
	HashMapClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HashMapClass.IncludeMixin(MapMixin)
	StdModule.AddConstantString("HashMap", Ref(HashMapClass))
	RegisterNativeClass("Std::HashMap", "value.HashMapClass")

	HashMapIteratorClass = NewClass()
	HashMapClass.AddConstantString("Iterator", Ref(HashMapIteratorClass))
	RegisterNativeClass("Std::HashMap::Iterator", "value.HashMapIteratorClass")
}
