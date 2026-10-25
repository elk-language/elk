package value

var OrderedHashRecordClass *Class         // ::Std::OrderedHashRecord
var OrderedHashRecordIteratorClass *Class // ::Std::OrderedHashRecord::Iterator

func initOrderedHashRecord() {
	OrderedHashRecordClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	OrderedHashRecordClass.IncludeMixin(RecordMixin)
	StdModule.AddConstantString("OrderedHashRecord", OrderedHashRecordClass.ToValue())
	RegisterNativeClass("Std::OrderedHashRecord", "value.OrderedHashRecordClass")

	OrderedHashRecordIteratorClass = NewClass()
	OrderedHashRecordClass.AddConstantString("Iterator", OrderedHashRecordIteratorClass.ToValue())
	RegisterNativeClass("Std::OrderedHashRecord::Iterator", "value.OrderedHashRecordIteratorClass")
}
