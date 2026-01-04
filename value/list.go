package value

var ListMixin *Mixin // ::Std::List

func initList() {
	ListMixin = NewMixin()
	ListMixin.IncludeMixin(TupleMixin)
	StdModule.AddConstantString("List", Ref(ListMixin))
	RegisterNativeMixin("Std::List", "value.ListMixin")
}
