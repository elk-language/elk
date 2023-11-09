package value

type GlobalObjectType struct{}

var GlobalObjectSingletonClass = NewClass()
var GlobalObject = GlobalObjectType{}

func (GlobalObjectType) Class() *Class {
	return ObjectClass
}

func (GlobalObjectType) DirectClass() *Class {
	return GlobalObjectSingletonClass
}

func (GlobalObjectType) SingletonClass() *Class {
	return GlobalObjectSingletonClass
}

func (GlobalObjectType) IsFrozen() bool {
	return true
}

func (GlobalObjectType) SetFrozen() {}

func (GlobalObjectType) Inspect() string {
	return "<GlobalObject>"
}

func (GlobalObjectType) InstanceVariables() SimpleSymbolMap {
	return nil
}
