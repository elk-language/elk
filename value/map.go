package value

var MapMixin *Mixin // ::Std::Map

func initMap() {
	MapMixin = NewMixin()
	MapMixin.IncludeMixin(RecordMixin)
	StdModule.AddConstantString("Map", Ref(MapMixin))
	RegisterNativeMixin("Std::Map", "value.MapMixin")
}
