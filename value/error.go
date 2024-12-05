package value

import (
	"fmt"
	"math"
)

// ::Std::Error
//
// Parent class for all exceptions.
var ErrorClass *Class

// ::Std::UnexpectedNilError
//
// Thrown when a `nil` value is encountered in a `must` expression.
var UnexpectedNilErrorClass *Class

// ::Std::TypeError
//
// Thrown when an argument given to a method
// has an incorrect type.
var TypeErrorClass *Class

// ::Std::PatternNotMatchedError
//
// Thrown when a pattern was not matched
// in destructuring etc
var PatternNotMatchedErrorClass *Class

// ::Std::ModifierMismatchError
//
// Thrown when a class was originally
// defined with some modifiers
// and later on reopened with different ones.
var ModifierMismatchErrorClass *Class

// ::Std::PrimitiveValueError
//
// Thrown when trying to access or set
// instance variables on a primitive object
// that cannot have instance variables.
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

// ::Std::IndexError
//
// Thrown when the index is invalid.
var IndexErrorClass *Class

// ::Std::ZeroDivisionError
//
// Thrown when an integer is divided by zero.
var ZeroDivisionErrorClass *Class

// ::Std::FormatError
//
// Thrown when a literal or interpreted string
// has an incorrect format.
var FormatErrorClass *Class

// ::Std::RegexCompileError
//
// Thrown when a Regex could not be compiled.
var RegexCompileErrorClass *Class

// ::Std::SealedClassError
//
// Thrown when trying to inherit
// from a sealed class.
var SealedClassErrorClass *Class

// ::Std::NoMethodError
//
// Thrown after attempting to call a method
// that is not available to the value.
var NoMethodErrorClass *Class

// ::Std::InvalidTimezoneError
//
// Thrown when a timezone wasn't found.
var InvalidTimezoneErrorClass *Class

// ::Std::NotBuiltinError
//
// Thrown when the called method is not builtin.
var NotBuiltinErrorClass *Class

var NotBuiltinError *Object

// Create a new Elk error.
func NewError(class *Class, message string) *Object {
	return &Object{
		class: class,
		instanceVariables: SymbolMap{
			SymbolTable.Add("message"): Ref(String(message)),
		},
	}
}

// Create a new error that signals that
// the given index is out of range.
func NewIndexOutOfRangeError(index string, length int) *Object {
	return Errorf(
		IndexErrorClass,
		"index %s out of range: %d...%d",
		index,
		-length,
		length,
	)
}

// Create a new error that signals that
// the big float precision is out of range (negative or too large).
func NewBigFloatPrecisionError(precision string) *Object {
	return Errorf(
		OutOfRangeErrorClass,
		"BigFloat precision cannot be negative or larger than %d, got: %s",
		uint64(math.MaxUint),
		precision,
	)
}

// Create a new error that signals that
// negative indices cannot be used in collection literals.
func NewNegativeIndicesInCollectionLiteralsError(index string) *Object {
	return Errorf(
		OutOfRangeErrorClass,
		"cannot use negative indices in collection literals: %s",
		index,
	)
}

// Create a new error that signals that
// the given capacity is too large.
func NewTooLargeCapacityError(capacity string) *Object {
	return Errorf(
		OutOfRangeErrorClass,
		"too large collection literal capacity: %s",
		capacity,
	)
}

// Create a new error that signals that
// the given capacity should not be negative.
func NewNegativeCapacityError(capacity string) *Object {
	return Errorf(
		OutOfRangeErrorClass,
		"capacity cannot be negative: %s",
		capacity,
	)
}

// Create a new error that signals that
// the given object cannot have a singleton class.
func NewSingletonError(given string) *Object {
	return Errorf(
		TypeErrorClass,
		"cannot get the singleton class of a primitive: `%s`",
		given,
	)
}

// Create a new error that signals that
// the given superclass is not a valid class object.
func NewInvalidSuperclassError(superclass string) *Object {
	return Errorf(
		TypeErrorClass,
		"`%s` cannot be used as a superclass",
		superclass,
	)
}

// Create a new error that signals that
// the given superclass doesn't match the original one.
func NewSuperclassMismatchError(class, wantSuperclass, gotSuperclass string) *Object {
	return Errorf(
		TypeErrorClass,
		"superclass mismatch in %s, expected: %s, got: %s",
		class,
		wantSuperclass,
		gotSuperclass,
	)
}

// Create a new error that signals that
// the given class should have different modifiers.
func NewModifierMismatchError(object, modifier string, with bool) *Object {
	var withStr string
	if with {
		withStr = "with"
	} else {
		withStr = "without"
	}
	return Errorf(
		ModifierMismatchErrorClass,
		"%s should be reopened %s the `%s` modifier",
		object,
		withStr,
		modifier,
	)
}

// Create a new error that signals that
// the type of the given argument is wrong.
func NewArgumentTypeError(argName, given, expected string) *Object {
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
func NewWrongArgumentCountError(method string, given, expected int) *Object {
	return Errorf(
		ArgumentErrorClass,
		"`%s` wrong number of arguments, given: %d, expected: %d",
		method,
		given,
		expected,
	)
}

// Create a new error that signals that
// some given arguments are not defined in the method.
func NewUnknownArgumentsError(method string, names []Symbol) *Object {
	return Errorf(
		ArgumentErrorClass,
		"`%s` unknown arguments: %s",
		method,
		InspectSlice(names),
	)
}

// Create a new error that signals that
// accessing instance variables of primitive values
// is impossible.
func NewCantAccessInstanceVariablesOnPrimitiveError(value string) *Object {
	return Errorf(
		PrimitiveValueErrorClass,
		"cannot access instance variables of a primitive value `%s`",
		value,
	)
}

// Create a new error that signals that
// setting instance variables of primitive values
// is impossible.
func NewCantSetInstanceVariablesOnPrimitiveError(value string) *Object {
	return Errorf(
		PrimitiveValueErrorClass,
		"cannot set instance variables of a primitive value `%s`",
		value,
	)
}

// Create a new error that signals that
// a required argument was not given.
func NewRequiredArgumentMissingError(methodName, paramName string) *Object {
	return Errorf(
		ArgumentErrorClass,
		"`%s` missing required argument `%s`",
		methodName,
		paramName,
	)
}

// Create a new error that signals that
// an argument is duplicated
func NewDuplicatedArgumentError(methodName, paramName string) *Object {
	return Errorf(
		ArgumentErrorClass,
		"`%s` duplicated argument `%s`",
		methodName,
		paramName,
	)
}

// Create a new error that signals that
// the number of given arguments is not within the accepted range.
func NewWrongArgumentCountRangeError(method string, given, expectedFrom, expectedTo int) *Object {
	return Errorf(
		ArgumentErrorClass,
		"`%s` wrong number of arguments, given: %d, expected: %d..%d",
		method,
		given,
		expectedFrom,
		expectedTo,
	)
}

// Create a new error that signals that
// the number of given arguments is not within the accepted range.
// For methods with rest parameters.
func NewWrongArgumentCountRestError(method string, given, expectedFrom int) *Object {
	return Errorf(
		ArgumentErrorClass,
		"`%s` wrong number of arguments, given: %d, expected: %d..",
		method,
		given,
		expectedFrom,
	)
}

// Create a new error that signals that
// the number of given arguments is not within the accepted range.
// For methods with rest parameters.
func NewWrongPositionalArgumentCountError(method string, given, expectedFrom int) *Object {
	return Errorf(
		ArgumentErrorClass,
		"`%s` wrong number of positional arguments, given: %d, expected: %d..",
		method,
		given,
		expectedFrom,
	)
}

// Create a new error that signals that the
// given value is not a class, even though it should be.
func NewIsNotClassError(value string) *Object {
	return Errorf(
		TypeErrorClass,
		"`%s` is not a class",
		value,
	)
}

// Create a new error that signals that the
// given value is not a class or mixin, even though it should be.
func NewIsNotClassOrMixinError(value string) *Object {
	return Errorf(
		TypeErrorClass,
		"`%s` is not a class or mixin",
		value,
	)
}

// Create a new error that signals that the
// given value is not a mixin, even though it should be.
func NewIsNotMixinError(value string) *Object {
	return Errorf(
		TypeErrorClass,
		"`%s` is not a mixin",
		value,
	)
}

// Create a new Std::RedefinedConstantError
func NewRedefinedConstantError(module, symbol string) *Object {
	return Errorf(
		RedefinedConstantErrorClass,
		"%s already has a constant named `%s`",
		module,
		symbol,
	)
}

// Create a new Std::NoMethodError.
func NewNoMethodError(methodName string, receiver Value) *Object {
	return Errorf(
		NoMethodErrorClass,
		"method `%s` is not available to value of class `%s`: %s",
		methodName,
		receiver.Class().PrintableName(),
		receiver.Inspect(),
	)
}

// Create a new error which signals
// that a value of one type cannot be coerced
// into the other type.
func NewCoerceError(target, other *Class) *Object {
	return Errorf(
		TypeErrorClass,
		"`%s` cannot be coerced into `%s`",
		other.PrintableName(),
		target.PrintableName(),
	)
}

// Create a new error which signals
// that a `nil` value has been encountered in a `must` expression
func NewUnexpectedNilError() *Object {
	return NewError(
		TypeErrorClass,
		"unexpected nil value in a must expression",
	)
}

// Create a new error which signals
// that the value can't be used as capacity.
func NewCapacityTypeError(val string) *Object {
	return Errorf(
		TypeErrorClass,
		"`%s` cannot be used as a collection literal's capacity",
		val,
	)
}

// Create a new error which signals
// that the given operand is not suitable for bit shifting
func NewBitshiftOperandError(other Value) *Object {
	return Errorf(
		TypeErrorClass,
		"`%s` cannot be used as a bitshift operand",
		other.Class().PrintableName(),
	)
}

// Create a new error which signals
// that a the program tried to divide by zero.
func NewZeroDivisionError() *Object {
	return NewError(
		ZeroDivisionErrorClass,
		"cannot divide by zero",
	)
}

// Mimics fmt.Errorf but creates an Elk error value.
func Errorf(class *Class, format string, a ...any) *Object {
	return NewError(class, fmt.Sprintf(format, a...))
}

// Implement the error interface.
func (e *Object) Error() string {
	msg := e.Message()
	if !msg.IsReference() {
		return fmt.Sprintf("%s: %s", e.class.PrintableName(), msg.Inspect())
	}

	switch msg := msg.AsReference().(type) {
	case String:
		return fmt.Sprintf("%s: %s", e.class.PrintableName(), msg.String())
	default:
		return fmt.Sprintf("%s: %s", e.class.PrintableName(), msg.Inspect())
	}
}

// Set the error message.
func (e *Object) SetMessage(message string) {
	e.instanceVariables.SetString("message", Ref(String(message)))
}

// Get the error message.
func (e *Object) Message() Value {
	return e.instanceVariables.GetString("message")
}

func initException() {
	ErrorClass = NewClass()
	StdModule.AddConstantString("Error", Ref(ErrorClass))

	UnexpectedNilErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("UnexpectedNilError", Ref(UnexpectedNilErrorClass))

	TypeErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("TypeError", Ref(TypeErrorClass))

	ModifierMismatchErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("ModifierMismatchError", Ref(ModifierMismatchErrorClass))

	NoConstantErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("NoConstantError", Ref(NoConstantErrorClass))

	RedefinedConstantErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("RedefinedConstantError", Ref(RedefinedConstantErrorClass))

	FormatErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("FormatError", Ref(FormatErrorClass))

	RegexCompileErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("RegexCompileError", Ref(RegexCompileErrorClass))

	NoMethodErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("NoMethodError", Ref(NoMethodErrorClass))

	ZeroDivisionErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("ZeroDivisionError", Ref(ZeroDivisionErrorClass))

	OutOfRangeErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("OutOfRangeError", Ref(OutOfRangeErrorClass))

	ArgumentErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("ArgumentError", Ref(ArgumentErrorClass))

	SealedClassErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("SealedClassError", Ref(SealedClassErrorClass))

	InvalidTimezoneErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("InvalidTimezoneError", Ref(InvalidTimezoneErrorClass))

	PrimitiveValueErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("PrimitiveValueError", Ref(PrimitiveValueErrorClass))

	PatternNotMatchedErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("PatternNotMatchedError", Ref(PatternNotMatchedErrorClass))

	IndexErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("IndexError", Ref(IndexErrorClass))

	NotBuiltinErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	StdModule.AddConstantString("NotBuiltinError", Ref(NotBuiltinErrorClass))
	NotBuiltinError = NewError(NotBuiltinErrorClass, "")
}
