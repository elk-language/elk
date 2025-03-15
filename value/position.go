package value

import (
	"fmt"

	"github.com/elk-language/elk/position"
)

var PositionClass *Class // Std::String::Position

type Position position.Position

// Creates a new Position.
func PositionConstructor(class *Class) Value {
	return Ref(&Position{})
}

func (*Position) Class() *Class {
	return PositionClass
}

func (*Position) DirectClass() *Class {
	return PositionClass
}

func (*Position) SingletonClass() *Class {
	return nil
}

func (p *Position) Copy() Reference {
	return p
}

func (*Position) InstanceVariables() SymbolMap {
	return nil
}

func (p *Position) Inspect() string {
	return fmt.Sprintf(
		"Std::String::Position(%d, %d, %d)",
		p.ByteOffset,
		p.Line,
		p.Column,
	)
}

func (p *Position) SimpleInspect() string {
	return fmt.Sprintf(
		"(%d, %d, %d)",
		p.ByteOffset,
		p.Line,
		p.Column,
	)
}

func (p *Position) Error() string {
	return p.Inspect()
}

func (p *Position) Equal(other *Position) bool {
	return (*position.Position)(p).Equal((*position.Position)(other))
}

func initPosition() {
	PositionClass = NewClassWithOptions(ClassWithConstructor(PositionConstructor))
	StringClass.AddConstantString("Position", Ref(PositionClass))
}
