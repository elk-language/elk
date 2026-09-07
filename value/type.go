package value

import (
	"github.com/elk-language/elk/value/symbol"
)

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
			symbol.ToSymbol("message"):     0,
			symbol.ToSymbol("diagnostics"): 1,
			symbol.ToSymbol("source_map"):  2,
		}),
	)
	ElkTypeCheckerClass.AddConstantString("Error", Ref(ElkTypeCheckerErrorClass))
	RegisterNativeClass("Std::Elk::Type::Checker::Error", "value.ElkTypeCheckerErrorClass")
}
