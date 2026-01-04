package value

var DurationMixin *Class // ::Std::Duration

func initDuration() {
	DurationMixin = NewMixin()
	StdModule.AddConstantString("Duration", Ref(DurationMixin))
	RegisterNativeMixin("Std::Duration", "value.DurationMixin")
}
