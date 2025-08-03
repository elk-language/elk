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

func initDiagnostic() {
	DiagnosticClass = NewClassWithOptions(ClassWithConstructor(DiagnosticConstructor))
	DiagnosticClass.AddConstantString("INFO", UInt8(diagnostic.INFO).ToValue())
	DiagnosticClass.AddConstantString("WARN", UInt8(diagnostic.WARN).ToValue())
	DiagnosticClass.AddConstantString("FAIL", UInt8(diagnostic.FAIL).ToValue())
	StdModule.AddConstantString("Diagnostic", Ref(DiagnosticClass))
}
