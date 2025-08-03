package value

import (
	"fmt"

	"github.com/elk-language/elk/position"
)

var LocationClass *Class // Std::FS::Location

type Location position.Location

// Creates a new Location.
func LocationConstructor(class *Class) Value {
	return Ref(&Location{})
}

func (*Location) Class() *Class {
	return LocationClass
}

func (*Location) DirectClass() *Class {
	return LocationClass
}

func (*Location) SingletonClass() *Class {
	return nil
}

func (l *Location) Copy() Reference {
	return l
}

func (*Location) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *Location) Inspect() string {
	return fmt.Sprintf(
		"Std::Location{&: %p, span: %s, file_path: %s}",
		l,
		l.SpanValue().Inspect(),
		l.FilePath,
	)
}

func (l *Location) Error() string {
	return l.Inspect()
}

func (l *Location) SpanValue() *Span {
	return (*Span)(l.Span)
}

func (l *Location) Equal(other *Location) bool {
	return (*position.Location)(l).Equal((*position.Location)(other))
}

func initLocation() {
	LocationClass = NewClassWithOptions(ClassWithConstructor(LocationConstructor))
	FSModule.AddConstantString("Location", Ref(LocationClass))
}
