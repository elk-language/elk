package value

import (
	"fmt"

	"github.com/elk-language/elk/position"
)

var SpanClass *Class // Std::String::Span

type Span position.Span

// Creates a new Span.
func SpanConstructor(class *Class) Value {
	return Ref(&Span{})
}

func (*Span) Class() *Class {
	return SpanClass
}

func (*Span) DirectClass() *Class {
	return SpanClass
}

func (*Span) SingletonClass() *Class {
	return nil
}

func (s *Span) Copy() Reference {
	return s
}

func (*Span) InstanceVariables() SymbolMap {
	return nil
}

func (s *Span) Inspect() string {
	return fmt.Sprintf(
		"Std::String::Span(%s, %s)",
		s.StartPosition().SimpleInspect(),
		s.EndPosition().SimpleInspect(),
	)
}

func (s *Span) Error() string {
	return s.Inspect()
}

func (s *Span) StartPosition() *Position {
	return (*Position)(s.StartPos)
}

func (s *Span) EndPosition() *Position {
	return (*Position)(s.EndPos)
}

func (s *Span) Equal(other *Span) bool {
	return (*position.Span)(s).Equal((*position.Span)(other))
}

func initSpan() {
	SpanClass = NewClassWithOptions(ClassWithConstructor(SpanConstructor))
	StringClass.AddConstantString("Span", Ref(SpanClass))
}
