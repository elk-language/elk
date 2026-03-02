package value

var HashRecordClass *Class         // ::Std::HashRecord
var HashRecordIteratorClass *Class // ::Std::HashRecord::Iterator

func initHashRecord() {
	HashRecordClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HashRecordClass.IncludeMixin(RecordMixin)
	StdModule.AddConstantString("HashRecord", HashRecordClass.ToValue())
	RegisterNativeClass("Std::HashRecord", "value.HashRecordClass")

	HashRecordIteratorClass = NewClass()
	HashRecordClass.AddConstantString("Iterator", HashRecordIteratorClass.ToValue())
	RegisterNativeClass("Std::HashRecord::Iterator", "value.HashRecordIteratorClass")
}
