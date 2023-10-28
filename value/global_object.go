package value

type GlobalObjectType struct{}

var GlobalObject = GlobalObjectType{}

func (GlobalObjectType) Class() *Class {
	return ObjectClass
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
