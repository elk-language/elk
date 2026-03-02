package value

var HashSetClass *Class         // ::Std::HashSet
var HashSetIteratorClass *Class // ::Std::HashSet::Iterator

func initHashSet() {
	HashSetClass = NewClass()
	HashSetClass.IncludeMixin(SetMixin)
	StdModule.AddConstantString("HashSet", Ref(HashSetClass))
	RegisterNativeClass("Std::HashSet", "value.HashSetClass")

	HashSetIteratorClass = NewClass()
	HashSetClass.AddConstantString("Iterator", Ref(HashSetIteratorClass))
	RegisterNativeClass("Std::HashSet::Iterator", "value.HashSetIteratorClass")
}
