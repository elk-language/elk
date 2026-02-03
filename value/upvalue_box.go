package value

var UpvalueBoxClass *Class // ::Std::UpvalueBox

func initUpvalueBox() {
	UpvalueBoxClass = NewClassWithOptions(
		ClassWithHidden("Std::UpvalueBox"),
		ClassWithSuperclass(BoxClass),
		ClassWithConstructor(UndefinedConstructor),
	)
}
