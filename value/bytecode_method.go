package value

// Method with Elk bytecode
type BytecodeMethod struct {
	BytecodeFunction
	Name  Symbol
	Owner *ModulelikeObject
}

func (*BytecodeMethod) method() {}
