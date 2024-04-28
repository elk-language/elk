package vm

// Wraps a bytecode method with associated local variables
// from the outer context
type Closure struct {
	Bytecode *BytecodeFunction
}

// Create a new closure
func NewClosure(bytecode *BytecodeFunction) *Closure {
	return &Closure{
		Bytecode: bytecode,
	}
}
