package value

// ::Std::PrimitiveIterable
var PrimitiveIterableInterface *Interface

// ::Std::Iterable
var IterableInterface *Interface

// ::Std::Iterable::FiniteBase
var IterableFiniteBaseMixin *Mixin

// ::Std::Iterable::Base
var IterableBaseMixin *Mixin

// ::Std::Iterable::NotFoundError
var IterableNotFoundErrorClass *Class

func initIterable() {
	PrimitiveIterableInterface = NewInterface()
	StdModule.AddConstantString("PrimitiveIterable", Ref(PrimitiveIterableInterface))
	RegisterNativeInterface("Std::PrimitiveIterable", "value.PrimitiveIterableInterface")

	IterableInterface = NewInterface()
	StdModule.AddConstantString("Iterable", Ref(IterableInterface))
	RegisterNativeInterface("Std::Iterable", "value.IterableInterface")

	IterableFiniteBaseMixin = NewMixin()
	IterableInterface.AddConstantString("FiniteBase", Ref(IterableFiniteBaseMixin))
	RegisterNativeMixin("Std::Iterable::FiniteBase", "value.IterableFiniteBaseMixin")

	IterableBaseMixin = NewMixin()
	IterableBaseMixin.IncludeMixin(IterableFiniteBaseMixin)
	IterableInterface.AddConstantString("Base", Ref(IterableBaseMixin))
	RegisterNativeMixin("Std::Iterable::Base", "value.IterableBaseMixin")

	IterableNotFoundErrorClass = NewClassWithOptions(ClassWithSuperclass(ErrorClass))
	IterableInterface.AddConstantString("NotFoundError", Ref(IterableNotFoundErrorClass))
	RegisterNativeClass("Std::Iterable::NotFoundError", "value.IterableNotFoundErrorClass")
}
