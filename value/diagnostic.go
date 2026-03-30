package value

import (
	"fmt"

	"github.com/elk-language/elk/position/diagnostic"
)

var DiagnosticClass *Class // ::Std::Diagnostic

type Diagnostic diagnostic.Diagnostic

// Creates a new Diagnostic.
func DiagnosticConstructor(class *Class) Value {
	return Ref(&Diagnostic{})
}

func (*Diagnostic) Class() *Class {
	return DiagnosticClass
}

func (*Diagnostic) DirectClass() *Class {
	return DiagnosticClass
}

func (*Diagnostic) SingletonClass() *Class {
	return nil
}

func (d *Diagnostic) Copy() Reference {
	return d
}

func (d *Diagnostic) ToValue() Value {
	return Ref(d)
}

func (*Diagnostic) InstanceVariables() *InstanceVariables {
	return nil
}

func (d *Diagnostic) Inspect() string {
	return fmt.Sprintf(
		"Std::Diagnostic{message: %s, location: %s, severity: %s}",
		String(d.Message).Inspect(),
		(*Location)(d.Location).Inspect(),
		d.Severity.String(),
	)
}

func (d *Diagnostic) Error() string {
	return d.Inspect()
}

func DiagnosticInfo() UInt8 {
	return UInt8(diagnostic.INFO)
}

func DiagnosticWarn() UInt8 {
	return UInt8(diagnostic.WARN)
}

func DiagnosticFail() UInt8 {
	return UInt8(diagnostic.FAIL)
}

func initDiagnostic() {
	DiagnosticClass = NewClassWithOptions(ClassWithConstructor(DiagnosticConstructor))
	StdModule.AddConstantString("Diagnostic", Ref(DiagnosticClass))
	RegisterNativeInterface("Std::Diagnostic", "value.DiagnosticClass")

	DiagnosticClass.AddConstantString("INFO", DiagnosticInfo().ToValue())
	RegisterNativeConstant("Std::Diagnostic", "value.DiagnosticInfo()", FetchGoType("value.UInt8"))

	DiagnosticClass.AddConstantString("WARN", DiagnosticWarn().ToValue())
	RegisterNativeConstant("Std::Diagnostic", "value.DiagnosticWarn()", FetchGoType("value.UInt8"))

	DiagnosticClass.AddConstantString("FAIL", DiagnosticFail().ToValue())
	RegisterNativeConstant("Std::Diagnostic", "value.DiagnosticFail()", FetchGoType("value.UInt8"))

}
