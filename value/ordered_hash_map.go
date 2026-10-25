package value

var OrderedHashMapClass *Class         // ::Std::OrderedHashMap
var OrderedHashMapIteratorClass *Class // ::Std::OrderedHashMap::Iterator

func initOrderedHashMap() {
	OrderedHashMapClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	OrderedHashMapClass.IncludeMixin(RecordMixin)
	StdModule.AddConstantString("OrderedHashMap", OrderedHashMapClass.ToValue())
	RegisterNativeClass("Std::OrderedHashMap", "value.OrderedHashMapClass")

	OrderedHashMapIteratorClass = NewClass()
	OrderedHashMapClass.AddConstantString("Iterator", OrderedHashMapIteratorClass.ToValue())
	RegisterNativeClass("Std::OrderedHashMap::Iterator", "value.OrderedHashMapIteratorClass")
}
