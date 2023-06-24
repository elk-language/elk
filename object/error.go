package object

import (
	"fmt"
)

var ExceptionClass *Class // ::Std::Exception
var ErrorClass *Class     // ::Std::Error
var TypeErrorClass *Class // ::Std::TypeError

type Error Object

// Create a new Elk error.
func NewError(class *Class, message string) *Error {
	return &Error{
		class: class,
		instanceVariables: SimpleSymbolMap{
			SymbolTable.Add("message").Id: String(message),
		},
	}
}

// Mimics fmt.Errorf but creates an Elk error.
func Errorf(class *Class, format string, a ...any) *Error {
	return NewError(class, fmt.Sprintf(format, a...))
}

// Implement the error interface.
func (e *Error) Error() string {
	switch msg := e.Message().(type) {
	case String:
		return fmt.Sprintf("%s: %s", e.class.PrintableName(), msg)
	default:
		return fmt.Sprintf("%s: %s", e.class.PrintableName(), msg.Inspect())
	}
}

// Set the error message.
func (e *Error) SetMessage(message string) {
	e.instanceVariables.SetString("message", String(message))
}

// Get the error message.
func (e *Error) Message() Value {
	return e.instanceVariables.GetString("message")
}

func (e *Error) InstanceVariables() SimpleSymbolMap {
	return e.instanceVariables
}

func (e *Error) Class() *Class {
	return e.class
}

func (e *Error) Inspect() string {
	return fmt.Sprintf(
		"%s%s",
		e.class.PrintableName(),
		e.instanceVariables.Inspect(),
	)
}

func (e *Error) IsFrozen() bool {
	return e.frozen
}

func (e *Error) SetFrozen() {
	e.frozen = true
}

func initException() {
	ExceptionClass = NewClass()
	StdModule.AddConstant("Exception", ExceptionClass)

	ErrorClass = NewClass(ClassWithParent(ExceptionClass))
	StdModule.AddConstant("Error", ErrorClass)

	TypeErrorClass = NewClass(ClassWithParent(ErrorClass))
	StdModule.AddConstant("TypeError", TypeErrorClass)
}
