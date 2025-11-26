package value

// ::Std::PrimitiveIterable
var PrimitiveIterableInterface *Interface

// ::Std::Iterable
var IterableInterface *Interface

// ::Std::Iterable::FiniteBase
var IterableFiniteBase *Mixin

// ::Std::Iterable::Base
var IterableBase *Mixin

// ::Std::Iterable::NotFoundError
var IterableNotFoundError *Class

func initIterable() {
	PrimitiveIterableInterface = NewInterface()
	StdModule.AddConstantString("PrimitiveIterable", Ref(PrimitiveIterableInterface))

	IterableInterface = NewInterface()
	StdModule.AddConstantString("Iterable", Ref(IterableInterface))

	IterableFiniteBase = NewMixin()
	IterableInterface.AddConstantString("FiniteBase", Ref(IterableFiniteBase))

	IterableBase = NewMixinWithOptions(MixinWithParent(IterableFiniteBase))
	IterableInterface.AddConstantString("Base", Ref(IterableBase))

	IterableNotFoundError = NewClassWithOptions(ClassWithSuperclass(ErrorClass))
	IterableInterface.AddConstantString("NotFoundError", Ref(IterableNotFoundError))
}
