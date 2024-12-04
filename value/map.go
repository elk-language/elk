package value

var MapMixin *Mixin // ::Std::Map

func initMap() {
	MapMixin = NewMixin()
	MapMixin.IncludeMixin(RecordMixin)
	StdModule.AddConstantString("Map", Ref(MapMixin))
}
