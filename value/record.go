package value

var RecordMixin *Mixin // ::Std::Record

func initRecord() {
	RecordMixin = NewMixin()
	StdModule.AddConstantString("Record", Ref(RecordMixin))
}
