package value

// An empty struct that implements ValueInterface.
// Can be used to quickly make a struct a valid elk value by embedding.
type ValueBase struct{}

func (ValueBase) Class() *Class {
	return UndefinedClass
}

func (ValueBase) DirectClass() *Class {
	return UndefinedClass
}

func (ValueBase) SingletonClass() *Class {
	return nil
}

func (ValueBase) Inspect() string {
	return "<value base>"
}

func (u ValueBase) Error() string {
	return u.Inspect()
}

func (ValueBase) InstanceVariables() *InstanceVariables {
	return nil
}
