package value

var ElkTypeMixin *Mixin             // ::Std::Elk::Type
var ElkTypeCheckerClass *Mixin      // ::Std::Elk::Type::Checker
var ElkTypeCheckerErrorClass *Class // ::Std::Elk::Type::Checker::Error

func initElkType() {
	ElkTypeMixin = NewMixin()
	ElkModule.AddConstantString("Type", Ref(ElkTypeMixin))
	RegisterNativeMixin("Std::Elk::Type", "value.ElkTypeMixin")

	ElkTypeCheckerClass = NewClass()
	ElkTypeMixin.AddConstantString("Checker", Ref(ElkTypeCheckerClass))
	RegisterNativeClass("Std::Elk::Type::Checker", "value.ElkTypeCheckerClass")

	ElkTypeCheckerErrorClass = NewClassWithOptions(
		ClassWithSuperclass(ErrorClass),
		ClassWithIvarIndices(IvarIndices{
			ToSymbol("message"):     0,
			ToSymbol("diagnostics"): 1,
			ToSymbol("source_map"):  2,
		}),
	)
	ElkTypeMixin.AddConstantString("Error", Ref(ElkTypeCheckerErrorClass))
	RegisterNativeClass("Std::Elk::Type::Checker::Error", "value.ElkTypeCheckerErrorClass")
}
