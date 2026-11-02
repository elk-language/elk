package value

var OrderedHashSetClass *Class         // ::Std::OrderedHashSet
var OrderedHashSetIteratorClass *Class // ::Std::OrderedHashSet::Iterator

func initOrderedHashSet() {
	OrderedHashSetClass = NewClass()
	OrderedHashSetClass.IncludeMixin(SetMixin)
	StdModule.AddConstantString("OrderedHashSet", Ref(OrderedHashSetClass))
	RegisterNativeClass("Std::OrderedHashSet", "value.OrderedHashSetClass")

	OrderedHashSetIteratorClass = NewClass()
	OrderedHashSetClass.AddConstantString("Iterator", Ref(OrderedHashSetIteratorClass))
	RegisterNativeClass("Std::OrderedHashSet::Iterator", "value.OrderedHashSetIteratorClass")
}
