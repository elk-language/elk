package value

import (
	"fmt"
)

// ::Std::Exception
//
// Parent class for all exceptions.
var ExceptionClass *Class

// ::Std::Error
//
// Parent class for all errors
// that are automatically caught
// by a `catch` expression without
// type constraints.
var ErrorClass *Class

// ::Std::TypeError
//
// Thrown when an argument given to a method
// has an incorrect type.
var TypeErrorClass *Class

// ::Std::ArgumentError
//
// Thrown when the arguments don't match
// the defined parameters for a given method.
var ArgumentErrorClass *Class

// ::Std::NoConstantError
//
// Thrown after trying to read a nonexistent constant.
var NoConstantErrorClass *Class

// ::Std::RedefinedConstantError
//
// Thrown after trying to redefine a constant.
var RedefinedConstantErrorClass *Class

// ::Std::OutOfRangeError
//
// Thrown when a numeric value is too large or too small
// to be used in a particular setting.
var OutOfRangeErrorClass *Class

// ::Std::ZeroDivisionError
//
// Thrown when an integer is divided by zero.
var ZeroDivisionErrorClass *Class

// ::Std::FormatError
//
// Thrown when a literal or interpreted string
// has an incorrect format.
var FormatErrorClass *Class

// ::Std::NoMethodError
//
// Thrown after attempting to call a method
// that is not available to the value.
var NoMethodErrorClass *Class

type Error Object

// Create a new Elk error.
func NewError(class *Class, message string) *Error {
	return &Error{
		class: class,
		instanceVariables: SimpleSymbolMap{
			SymbolTable.Add("message"): String(message),
		},
	}
}

// Create a new error that signals that
// the number of given arguments is wrong.
func NewWrongArgumentCountError(given, expected int) *Error {
	return Errorf(
		ArgumentErrorClass,
		"wrong number of arguments, given: %d, expected: %d",
		given,
		expected,
	)
}

// Create a new error that signals that the
// given value is not a module, even though it should be.
func NewIsNotModuleError(value string) *Error {
	return Errorf(
		TypeErrorClass,
		"`%s` is not a module",
		value,
	)
}

// Create a new Std::RedefinedConstantError
func NewRedefinedConstantError(module, symbol string) *Error {
	return Errorf(
		RedefinedConstantErrorClass,
		"%s already has a constant named `%s`",
		module,
		symbol,
	)
}

// Create a new Std::NoMethodError.
func NewNoMethodError(methodName string, receiver Value) *Error {
	return Errorf(
		NoMethodErrorClass,
		"method `%s` is not available to value of class `%s`: %s",
		methodName,
		receiver.Class().PrintableName(),
		receiver.Inspect(),
	)
}

// Create a new error which signals
// that a value of one type can't be coerced
// into the other type.
func NewCoerceError(receiver, other Value) *Error {
	return Errorf(
		TypeErrorClass,
		"`%s` can't be coerced into `%s`",
		other.Class().PrintableName(),
		receiver.Class().PrintableName(),
	)
}

// Create a new error which signals
// that the given operand is not suitable for bit shifting
func NewBitshiftOperandError(other Value) *Error {
	return Errorf(
		TypeErrorClass,
		"`%s` can't be used as a bitshift operand",
		other.Class().PrintableName(),
	)
}

// Create a new error which signals
// that a the program tried to divide by zero.
func NewZeroDivisionError() *Error {
	return NewError(
		ZeroDivisionErrorClass,
		"can't divide by zero",
	)
}

// Mimics fmt.Errorf but creates an Elk error value.
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
	val, _ := e.instanceVariables.GetString("message")
	return val
}

func (e *Error) InstanceVariables() SimpleSymbolMap {
	return e.instanceVariables
}

func (e *Error) Class() *Class {
	if !e.class.IsSingleton() {
		return e.class
	}

	return e.class.Class()
}

func (e *Error) DirectClass() *Class {
	return e.class
}

func (e *Error) SingletonClass() *Class {
	if e.class.IsSingleton() {
		return e.class
	}

	singletonClass := NewClass()
	singletonClass.SetSingleton()
	singletonClass.Parent = e.class
	e.class = singletonClass
	return singletonClass
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
	StdModule.AddConstantString("Exception", ExceptionClass)

	ErrorClass = NewClassWithOptions(ClassWithParent(ExceptionClass))
	StdModule.AddConstantString("Error", ErrorClass)

	TypeErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("TypeError", TypeErrorClass)

	NoConstantErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("NoConstantError", NoConstantErrorClass)

	RedefinedConstantErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("RedefinedConstantError", RedefinedConstantErrorClass)

	FormatErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("FormatError", FormatErrorClass)

	NoMethodErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("NoMethodError", NoMethodErrorClass)

	ZeroDivisionErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("ZeroDivisionError", ZeroDivisionErrorClass)

	OutOfRangeErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("OutOfRangeError", OutOfRangeErrorClass)

	ArgumentErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("ArgumentError", ArgumentErrorClass)
}
