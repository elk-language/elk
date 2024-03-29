package value

type GlobalObjectType struct{}

var GlobalObjectSingletonClass *Class
var GlobalObject = GlobalObjectType{}

func initGlobalObject() {
	GlobalObjectSingletonClass = NewClassWithOptions(
		ClassWithSingleton(),
		ClassWithNoInstanceVariables(),
	)
}

func (GlobalObjectType) Class() *Class {
	return ObjectClass
}

func (GlobalObjectType) DirectClass() *Class {
	return GlobalObjectSingletonClass
}

func (GlobalObjectType) SingletonClass() *Class {
	return GlobalObjectSingletonClass
}

func (g GlobalObjectType) Copy() Value {
	return g
}

func (GlobalObjectType) Inspect() string {
	return "<GlobalObject>"
}

func (GlobalObjectType) InstanceVariables() SymbolMap {
	return nil
}
