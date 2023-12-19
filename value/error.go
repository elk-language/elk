package value

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
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

// ::Std::PrimitiveValueError
//
// Thrown when trying to access or set
// instance variables on a primitive object
// that can't have instance variables.
var PrimitiveValueErrorClass *Class

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

// ::Std::FrozenMethodError
//
// Thrown when trying to override
// a frozen method.
var FrozenMethodErrorClass *Class

// ::Std::NoMethodError
//
// Thrown after attempting to call a method
// that is not available to the value.
var NoMethodErrorClass *Class

// ::Std::InvalidTimezoneError
//
// Thrown when a timezone wasn't found.
var InvalidTimezoneErrorClass *Class

type Error struct {
	Object
}

func NewErrorComparer(opts cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *Error) bool {
		if x == nil && y == nil {
			return true
		}

		if x == nil || y == nil {
			return false
		}

		return x.class == y.class &&
			cmp.Equal(x.instanceVariables, y.instanceVariables, opts...)
	})
}

// Create a new Elk error.
func NewError(class *Class, message string) *Error {
	return &Error{
		Object: Object{
			class: class,
			instanceVariables: SymbolMap{
				SymbolTable.Add("message"): String(message),
			},
		},
	}
}

// Create a new error that signals that
// the number of given arguments is wrong.
func NewArgumentTypeError(argName, given, expected string) *Error {
	return Errorf(
		TypeErrorClass,
		"wrong argument type for `%s`, given: `%s`, expected: `%s`",
		argName,
		given,
		expected,
	)
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

// Create a new error that signals that
// some given arguments are not defined in the method.
func NewUnknownArgumentsError(names []Symbol) *Error {
	return Errorf(
		ArgumentErrorClass,
		"unknown arguments: %s",
		InspectSlice(names),
	)
}

// Create a new error that signals that
// the user tried to create an alias for a nonexistent method.
func NewCantCreateAnAliasForNonexistentMethod(methodName string) *Error {
	return Errorf(
		NoMethodErrorClass,
		"can't create an alias for a nonexistent method: %s",
		methodName,
	)
}

// Create a new error that signals that
// some given arguments are not defined in the method.
func NewCantOverrideAFrozenMethod(methodName string) *Error {
	return Errorf(
		FrozenMethodErrorClass,
		"can't override a frozen method: %s",
		methodName,
	)
}

// Create a new error that signals that
// accessing instance variables of primitive values
// is impossible.
func NewCantAccessInstanceVariablesOnPrimitiveError(value string) *Error {
	return Errorf(
		PrimitiveValueErrorClass,
		"can't access instance variables of a primitive value `%s`",
		value,
	)
}

// Create a new error that signals that
// a required argument was not given.
func NewRequiredArgumentMissingError(methodName, paramName string) *Error {
	return Errorf(
		ArgumentErrorClass,
		"missing required argument `%s` in call to `%s`",
		paramName,
		methodName,
	)
}

// Create a new error that signals that
// an argument is duplicated
func NewDuplicatedArgumentError(methodName, paramName string) *Error {
	return Errorf(
		ArgumentErrorClass,
		"duplicated argument `%s` in call to `%s`",
		paramName,
		methodName,
	)
}

// Create a new error that signals that
// the number of given arguments is not within the accepted range.
func NewWrongArgumentCountRangeError(given, expectedFrom, expectedTo int) *Error {
	return Errorf(
		ArgumentErrorClass,
		"wrong number of arguments, given: %d, expected: %d..%d",
		given,
		expectedFrom,
		expectedTo,
	)
}

// Create a new error that signals that
// the number of given arguments is not within the accepted range.
// For methods with rest parameters.
func NewWrongArgumentCountRestError(given, expectedFrom int) *Error {
	return Errorf(
		ArgumentErrorClass,
		"wrong number of arguments, given: %d, expected: %d..",
		given,
		expectedFrom,
	)
}

// Create a new error that signals that
// the number of given arguments is not within the accepted range.
// For methods with rest parameters.
func NewWrongPositionalArgumentCountError(given, expectedFrom int) *Error {
	return Errorf(
		ArgumentErrorClass,
		"wrong number of positional arguments, given: %d, expected: %d..",
		given,
		expectedFrom,
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

// Create a new error that signals that the
// given value is not a mixin, even though it should be.
func NewIsNotMixinError(value string) *Error {
	return Errorf(
		TypeErrorClass,
		"`%s` is not a mixin",
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
	return e.instanceVariables.GetString("message")
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

	FrozenMethodErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("FrozenMethodError", FrozenMethodErrorClass)

	InvalidTimezoneErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("InvalidTimezoneError", InvalidTimezoneErrorClass)

	PrimitiveValueErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("PrimitiveValueError", InvalidTimezoneErrorClass)
}
