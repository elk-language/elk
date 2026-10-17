package value

import "fmt"

// Elk Method object
type Method interface {
	Reference
	Function
	// Name of the method
	Name() Symbol
	Namespace() *Class
	SetNamespace(*Class)
	MethodBody()
}

func MethodNamespacedName(m Method) string {
	namespace := m.Namespace()
	name := m.Name()
	if namespace.IsSingleton() {
		namespaceName := namespace.Name()[1:]
		return fmt.Sprintf("%s::%s", namespaceName, name.String())
	}

	namespaceName := namespace.Name()
	return fmt.Sprintf("%s.:%s", namespaceName, name.String())
}

var MethodClass *Class // ::Std::Method

func initMethod() {
	MethodClass = NewClassWithOptions(
		ClassWithSuperclass(FunctionClass),
	)
	StdModule.AddConstantString("Method", Ref(MethodClass))
	RegisterNativeClass("Std::Method", "value.MethodClass")
}
