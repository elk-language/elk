package value

import "fmt"

// Wraps a native Go value in an Elk value
type Native struct {
	Value any
}

// Create a new native value
func NewNative(value any) *Native {
	return &Native{
		Value: value,
	}
}

func (n *Native) Copy() Reference {
	return n
}

func (*Native) Class() *Class {
	return nil
}

func (*Native) DirectClass() *Class {
	return nil
}

func (*Native) SingletonClass() *Class {
	return nil
}

func (n *Native) Inspect() string {
	return fmt.Sprintf("Native{Value: %#v}", n.Value)
}

func (n *Native) Error() string {
	return n.Inspect()
}

func (*Native) InstanceVariables() SymbolMap {
	return nil
}
