package value

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
)

var MacroModule *Module      // Std::Macro
var MacroContextClass *Class // Std::Macro::Context

// Macro execution context created by the typechecker
// and passed to the macro during expansion.
type MacroContext struct {
	TypeChecker types.Checker
	Location    *position.Location // location of the macro call
	Namespace   types.Namespace    // namespace the macro call is nested under
}

func NewMacroContext(checker types.Checker, loc *position.Location, namespace types.Namespace) *MacroContext {
	return &MacroContext{
		TypeChecker: checker,
		Location:    loc,
		Namespace:   namespace,
	}
}

func initMacro() {
	MacroModule = NewModule()
	StdModule.AddConstantString("Macro", Ref(MacroModule))
	RegisterNativeModule("Std::Macro", "value.MacroModule")

	MacroContextClass = NewClass()
	MacroModule.AddConstantString("Context", Ref(MacroContextClass))
	RegisterNativeClass("Std::Macro::Context", "value.MacroContextClass")
}
