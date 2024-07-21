package types

// This file is auto-generated, please do not edit it manually

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

func setupGlobalEnvironmentFromHeaders(env *GlobalEnvironment) {
	objectClass := env.StdSubtypeClass(symbol.Object)
	namespace := env.Root

	// Define all namespaces
	{
		namespace := namespace.TryDefineModule("", value.ToSymbol("Std"))
		{
			namespace := namespace.TryDefineClass("A dynamically resizable list data structure backed\nby an array.\n\nIt is an ordered collection of integer indexed values.", false, true, true, value.ToSymbol("ArrayList"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("A tuple data structure backed by an array.\n\nIt is an ordered, immutable collection of integer indexed values.\nA tuple is an immutable list.", false, true, true, value.ToSymbol("ArrayTuple"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Represents a multi-precision floating point number (a fraction like `1.2`, `0.1`).\n\n```\nsign × mantissa × 2**exponent\n```\n\nwith 0.5 <= mantissa < 1.0, and MinExp <= exponent <= MaxExp.\nA `BigFloat` may also be zero (+0, -0) or infinite (+Inf, -Inf).\nAll BigFloats are ordered.\n\nBy setting the desired precision to 24 or 53,\n`BigFloat` operations produce the same results as the corresponding float32 or float64 IEEE-754 arithmetic for operands that\ncorrespond to normal (i.e., not denormal) `Float`, `Float32` and `Float64` numbers.\nExponent underflow and overflow lead to a `0` or an Infinity for different values than IEEE-754 because `BigFloat` exponents have a much larger range.", false, true, true, value.ToSymbol("BigFloat"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Bool"), objectClass, env)
		namespace.TryDefineClass("Represents a single Unicode code point.", false, true, true, value.ToSymbol("Char"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Class"), objectClass, env)
		namespace.TryDefineClass("A base class for most errors in Elk stdlib.", false, false, false, value.ToSymbol("Error"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("False"), objectClass, env)
		namespace.TryDefineClass("Represents a floating point number (a fraction like `1.2`, `0.1`).\n\nThis float type has 64 bits on 64 bit platforms\nand 32 bit on 32 bit platforms.", false, true, true, value.ToSymbol("Float"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Float32"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Float64"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("A dynamically resizable map data structure backed\nby an array with a hashing algorithm.\n\nIt is an unordered collection of key-value pairs.", false, true, true, value.ToSymbol("HashMap"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("A record data structure backed by an array with a hashing algorithm.\n\nIt is an unordered immutable collection of key-value pairs.\nA record is an immutable map.", false, true, true, value.ToSymbol("HashRecord"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("A dynamically resizable set data structure backed\nby an array with a hashing algorithm.\n\nIt is an unordered collection of unique values.", false, true, true, value.ToSymbol("HashSet"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a human readable string\nthat represents the structure of the value.", value.ToSymbol("Inspectable"), env)
		namespace.TryDefineClass("Represents an integer (a whole number like `1`, `2`, `3`, `-5`, `0`).\n\nThis integer type is automatically resized so\nit can hold an arbitrarily large/small number.", false, true, true, value.ToSymbol("Int"), objectClass, env)
		namespace.TryDefineClass("Represents a signed 16 bit integer (a whole number like `1i16`, `2i16`, `-3i16`, `0i16`).", false, true, true, value.ToSymbol("Int16"), objectClass, env)
		namespace.TryDefineClass("Represents a signed 32 bit integer (a whole number like `1i32`, `2i32`, `-3i32`, `0i32`).", false, true, true, value.ToSymbol("Int32"), objectClass, env)
		namespace.TryDefineClass("Represents a signed 64 bit integer (a whole number like `1i64`, `2i64`, `-3i64`, `0i64`).", false, true, true, value.ToSymbol("Int64"), objectClass, env)
		namespace.TryDefineClass("Represents a signed 8 bit integer (a whole number like `1i8`, `2i8`, `-3i8`, `0i8`).", false, true, true, value.ToSymbol("Int8"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Interface"), objectClass, env)
		namespace.TryDefineInterface("Represents a value that can be iterated over in a `for` loop.", value.ToSymbol("Iterable"), env)
		{
			namespace := namespace.TryDefineInterface("An interface that represents objects\nthat allow for external iteration.", value.ToSymbol("Iterator"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineInterface("An interface that represents an ordered, mutable collection\nof elements indexed by integers starting at `0`.", value.ToSymbol("List"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Method"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Mixin"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Module"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Nil"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Object"), objectClass, env)
		namespace.TryDefineClass("Thrown when a numeric value is too large or too small to be used in a particular setting.", false, false, false, value.ToSymbol("OutOfRangeError"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("A `Pair` represents a 2-element tuple,\nor a key-value pair.", false, true, true, value.ToSymbol("Pair"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("A `Regex` represents regular expression that can be used\nto match a pattern against strings.", false, true, true, value.ToSymbol("Regex"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("", false, true, true, value.ToSymbol("String"), objectClass, env)
			namespace.TryDefineClass("Iterates over all bytes of a `String`.", false, true, true, value.ToSymbol("ByteIterator"), objectClass, env)
			namespace.TryDefineClass("Iterates over all unicode code points of a `String`.", false, true, true, value.ToSymbol("CharIterator"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a string.", value.ToSymbol("StringConvertible"), env)
		namespace.TryDefineClass("Represents an interned string.\n\nA symbol is an integer ID that is associated\nwith a particular name (string).\n\nA few symbols with the same name refer to the same ID.\n\nComparing symbols happens in constant time, so it's\nusually faster than comparing strings.", false, true, true, value.ToSymbol("Symbol"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("True"), objectClass, env)
		namespace.TryDefineClass("Represents an unsigned 16 bit integer (a positive whole number like `1u16`, `2u16`, `3u16`, `0u16`).", false, true, true, value.ToSymbol("UInt16"), objectClass, env)
		namespace.TryDefineClass("Represents an unsigned 32 bit integer (a positive whole number like `1u32`, `2u32`, `3u32`, `0u32`).", false, true, true, value.ToSymbol("UInt32"), objectClass, env)
		namespace.TryDefineClass("Represents an unsigned 64 bit integer (a positive whole number like `1u64`, `2u64`, `3u64`, `0u64`).", false, true, true, value.ToSymbol("UInt64"), objectClass, env)
		namespace.TryDefineClass("Represents an unsigned 8 bit integer (a positive whole number like `1u8`, `2u8`, `3u8`, `0u8`).", false, true, true, value.ToSymbol("UInt8"), objectClass, env)
		namespace.TryDefineClass("`Value` is the superclass class of all\nElk classes.", false, false, false, value.ToSymbol("Value"), nil, env)
		namespace.Name() // noop - avoid unused variable error
	}

	// Define methods, constants

	{
		namespace := env.Root

		namespace.Name() // noop - avoid unused variable error

		// Include mixins

		// Implement interfaces

		// Define methods

		// Define constants

		// Define instance variables

		{
			namespace := namespace.SubtypeString("Std").(*Module)

			namespace.Name() // noop - avoid unused variable error

			// Include mixins

			// Implement interfaces

			// Define methods

			// Define constants

			// Define instance variables

			{
				namespace := namespace.SubtypeString("ArrayList").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Create a new `ArrayList` containing the elements of `self`\nrepeated `n` times.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::ArrayList", env), nil)
				namespace.DefineMethod("Create a new `ArrayList` containing the elements of `self`\nand another given `ArrayList` or `ArrayTuple`.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::ArrayList", env), NameToType("Std::ArrayTuple", env)), NormalParameterKind, false)}, NameToType("Std::ArrayList", env), nil)
				namespace.DefineMethod("Adds the given value to the list.\n\nReallocates the underlying array if it is\ntoo small to hold it.", false, false, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("value"), NewNamedType("Std::ArrayList::Element", Any{}), NormalParameterKind, false)}, NameToType("Std::ArrayList", env), nil)
				namespace.DefineMethod("Check whether the given value is an `ArrayList`\nwith the same elements.", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Check whether the given value is an `ArrayList` or `ArrayTuple`\nwith the same elements.", false, true, true, value.ToSymbol("=~"), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the element under the given index.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", false, false, true, value.ToSymbol("[]"), []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NewNamedType("Std::ArrayList::Element", Any{}), nil)
				namespace.DefineMethod("Set the element under the given index to the given value.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", false, false, true, value.ToSymbol("[]="), []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NewNamedType("Std::ArrayList::Element", Any{}), NormalParameterKind, false)}, Void{}, nil)
				namespace.DefineMethod("Adds the given values to the list.\n\nReallocates the underlying array if it is\ntoo small to hold them.", false, false, true, value.ToSymbol("append"), []*Parameter{NewParameter(value.ToSymbol("values"), NewNamedType("Std::ArrayList::Element", Any{}), PositionalRestParameterKind, false)}, NameToType("Std::ArrayList", env), nil)
				namespace.DefineMethod("Returns the number of elements that can be\nheld by the underlying array.\n\nThis value will change when the list gets resized,\nand the underlying array gets reallocated", false, false, true, value.ToSymbol("capacity"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Check whether the given `value` is present in this list.", false, false, true, value.ToSymbol("contains"), []*Parameter{NewParameter(value.ToSymbol("value"), NewNamedType("Std::ArrayList::Element", Any{}), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Mutates the list.\n\nReallocates the underlying array to hold\nthe given number of new elements.\n\nExpands the `capacity` of the list by `new_slots`", false, false, true, value.ToSymbol("grow"), []*Parameter{NewParameter(value.ToSymbol("new_slots"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::ArrayList", env), nil)
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the list.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::ArrayList::Iterator", env), nil)
				namespace.DefineMethod("Returns the number of left slots for new elements\nin the underlying array.\nIt tells you how many more elements can be\nadded to the list before the underlying array gets\nreallocated.\n\nIt is always equal to `capacity - length`.", false, false, true, value.ToSymbol("left_capacity"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns the number of elements present in the list.", false, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, false, true, value.ToSymbol("map"), nil, Void{}, nil)
				namespace.DefineMethod("", false, false, true, value.ToSymbol("map_mut"), nil, Void{}, nil)

				// Define constants

				// Define instance variables

				{
					namespace := namespace.SubtypeString("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins

					// Implement interfaces

					// Define methods
					namespace.DefineMethod("Returns `self`.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::ArrayList::Iterator", env), nil)
					namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", false, false, true, value.ToSymbol("next"), nil, NewNamedType("Std::ArrayList::Iterator::Element", Any{}), NewSymbolLiteral("stop_iteration"))

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.SubtypeString("ArrayTuple").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Create a new `ArrayTuple` containing the elements of `self`\nrepeated `n` times.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::ArrayList", env), nil)
				namespace.DefineMethod("Create a new `ArrayTuple` containing the elements of `self`\nand another given `ArrayTuple`.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::ArrayTuple", env), NormalParameterKind, false)}, NameToType("Std::ArrayTuple", env), nil)
				namespace.DefineMethod("Check whether the given value is an `ArrayTuple`\nwith the same elements.", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Check whether the given value is an `ArrayTuple` or `ArrayList`\nwith the same elements.", false, true, true, value.ToSymbol("=~"), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the element under the given index.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", false, false, true, value.ToSymbol("[]"), []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NewNamedType("Std::ArrayTuple::Element", Any{}), nil)
				namespace.DefineMethod("Check whether the given `value` is present in this tuple.", false, false, true, value.ToSymbol("contains"), []*Parameter{NewParameter(value.ToSymbol("value"), NewNamedType("Std::ArrayTuple::Element", Any{}), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the tuple.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::ArrayTuple::Iterator", env), nil)
				namespace.DefineMethod("Returns the number of elements present in the tuple.", false, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, false, true, value.ToSymbol("map"), nil, Void{}, nil)

				// Define constants

				// Define instance variables

				{
					namespace := namespace.SubtypeString("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins

					// Implement interfaces

					// Define methods
					namespace.DefineMethod("Returns `self`.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::ArrayTuple::Iterator", env), nil)
					namespace.DefineMethod("Get the next element of the tuple.\nThrows `:stop_iteration` when there are no more elements.", false, false, true, value.ToSymbol("next"), nil, NewNamedType("Std::ArrayTuple::Iterator::Element", Any{}), NewSymbolLiteral("stop_iteration"))

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.SubtypeString("BigFloat").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10bf\n\tvar b = 3bf\n\ta % b #=> 1bf\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("Multiply this float by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("Exponentiate this float, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("Add `other` to this bigfloat.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1.2bf\n\t+a #=> 1.2bf\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("Subtract `other` from this bigfloat.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("Returns the result of negating the number.\n\n```\n\tvar a = 1.2bf\n\t-a #=> -1.2bf\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("Divide this float by another float.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this float with another float.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Sets the precision to the given integer.", false, false, true, value.ToSymbol("set_precision"), []*Parameter{NewParameter(value.ToSymbol("precision"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("returns the mantissa precision of `self` in bits.", false, false, true, value.ToSymbol("precision"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Sets the precision to the given integer.", false, false, true, value.ToSymbol("set_precision"), []*Parameter{NewParameter(value.ToSymbol("precision"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_bigfloat"), nil, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("Converts to a fixed-precision floating point number.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the bigfloat to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the bigfloat to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Converts the bigfloat to an automatically resized integer.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Converts the bigfloat to a 16-bit integer.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Converts the bigfloat to a 32-bit integer.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Converts the bigfloat to a 64-bit integer.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Converts the bigfloat to a 8-bit integer.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Converts the bigfloat to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the bigfloat to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the bigfloat to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the bigfloat to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Char").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Concatenate this `Char`\nwith another `Char` or `String`.\n\nCreates a new `String` containing the content\nof both operands.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Get the next Unicode codepoint by incrementing by 1.", false, true, true, value.ToSymbol("++"), nil, NameToType("Std::Char", env), nil)
				namespace.DefineMethod("Get the previous Unicode codepoint by decrementing by 1.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::Char", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this char with another char or string.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Object", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the number of bytes that this\ncharacter contains.", false, false, true, value.ToSymbol("byte_count"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Always returns 1.\nFor better compatibility with `String`.", false, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Concatenate this `Char`\nwith another `Char` or `String`.\n\nCreates a new `String` containing the content\nof both operands.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Always returns 1.\nFor better compatibility with `String`.", false, false, true, value.ToSymbol("grapheme_count"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Always returns false.\nFor better compatibility with `String`.", false, false, true, value.ToSymbol("is_empty"), nil, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Always returns 1.\nFor better compatibility with `String`.", false, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Return the lowercase version of this character.", false, false, true, value.ToSymbol("lowercase"), nil, NameToType("Std::Char", env), nil)
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Creates a new string with this character.", false, false, true, value.ToSymbol("to_string"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the `Char` to a `Symbol`.", false, false, true, value.ToSymbol("to_symbol"), nil, NameToType("Std::Symbol", env), nil)
				namespace.DefineMethod("Return the uppercase version of this character.", false, false, true, value.ToSymbol("uppercase"), nil, NameToType("Std::Char", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("False").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Bool", env))

				// Include mixins

				// Implement interfaces

				// Define methods

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Float").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10\n\tvar b = 3\n\ta % b #=> 1\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Multiply this float by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Exponentiate this float, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Add `other` to this float.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1.2\n\t+a #=> 1.2\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Subtract `other` from this float.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Returns the result of negating the number.\n\n```\n\tvar a = 1.2\n\t-a #=> -1.2\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Divide this float by another float.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this float with another number.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Converts the float to a multi-precision floating point number.", false, false, true, value.ToSymbol("to_bigfloat"), nil, NameToType("Std::BigFloat", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the float to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the float to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Converts the float to an automatically resized integer.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Converts the float to a 16-bit integer.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Converts the float to a 32-bit integer.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Converts the float to a 64-bit integer.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Converts the float to a 8-bit integer.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Converts the float to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the float to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the float to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the float to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("HashMap").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Create a new `HashMap` containing the pairs of `self`\nand another given `HashMap` or `HashRecord`.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::HashMap", env), NameToType("Std::HashRecord", env)), NormalParameterKind, false)}, NameToType("Std::HashMap", env), nil)
				namespace.DefineMethod("Check whether the given value is a `HashMap`\nwith the same elements.", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Check whether the given value is an `HashMap` or `HashRecord`\nwith the same elements.", false, true, true, value.ToSymbol("=~"), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the element under the given key.\nReturns `nil` when the key is not present.", false, false, true, value.ToSymbol("[]"), []*Parameter{NewParameter(value.ToSymbol("key"), NewNamedType("Std::HashMap::Key", Any{}), NormalParameterKind, false)}, NewNilable(NewNamedType("Std::HashMap::Value", Any{})), nil)
				namespace.DefineMethod("Set the element under the given index to the given value.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", false, false, true, value.ToSymbol("[]="), []*Parameter{NewParameter(value.ToSymbol("key"), NewNamedType("Std::HashMap::Key", Any{}), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NewNamedType("Std::HashMap::Value", Any{}), NormalParameterKind, false)}, Void{}, nil)
				namespace.DefineMethod("Returns the number of key-value pairs that can be\nheld by the underlying array.\n\nThis value will change when the map gets resized,\nand the underlying array gets reallocated.", false, false, true, value.ToSymbol("capacity"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Check whether the given `pair` is present in this map.", false, false, true, value.ToSymbol("contains"), []*Parameter{NewParameter(value.ToSymbol("pair"), NameToType("Std::Pair", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Check whether the given `key` is present in this map.", false, false, true, value.ToSymbol("contains_key"), []*Parameter{NewParameter(value.ToSymbol("key"), NewNamedType("Std::HashMap::Key", Any{}), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Check whether the given `value` is present in this map.", false, false, true, value.ToSymbol("contains_value"), []*Parameter{NewParameter(value.ToSymbol("value"), NewNamedType("Std::HashMap::Value", Any{}), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Mutates the map.\n\nReallocates the underlying array to hold\nthe given number of new elements.\n\nExpands the `capacity` of the list by `new_slots`", false, false, true, value.ToSymbol("grow"), []*Parameter{NewParameter(value.ToSymbol("new_slots"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::HashMap", env), nil)
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the map.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::HashMap::Iterator", env), nil)
				namespace.DefineMethod("Returns the number of left slots for new key-value pairs\nin the underlying array.\nIt tells you how many more elements can be\nadded to the map before the underlying array gets\nreallocated.\n\nIt is always equal to `capacity - length`.", false, false, true, value.ToSymbol("left_capacity"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns the number of key-value pairs present in the map.", false, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, false, true, value.ToSymbol("map"), nil, Void{}, nil)
				namespace.DefineMethod("", false, false, true, value.ToSymbol("map_values"), nil, Void{}, nil)
				namespace.DefineMethod("", false, false, true, value.ToSymbol("map_values_mut"), nil, Void{}, nil)

				// Define constants

				// Define instance variables

				{
					namespace := namespace.SubtypeString("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins

					// Implement interfaces

					// Define methods
					namespace.DefineMethod("Returns `self`.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::HashMap::Iterator", env), nil)
					namespace.DefineMethod("Get the next pair of the map.\nThrows `:stop_iteration` when there are no more elements.", false, false, true, value.ToSymbol("next"), nil, NewNamedType("Std::HashMap::Iterator::Element", Any{}), NewSymbolLiteral("stop_iteration"))

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.SubtypeString("HashRecord").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Create a new `HashRecord` containing the pairs of `self`\nand another given `HashRecord`.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::HashRecord", env), NormalParameterKind, false)}, NameToType("Std::HashRecord", env), nil)
				namespace.DefineMethod("Check whether the given value is a `HashRecord`\nwith the same elements.", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Check whether the given value is an `HashRecord` or `HashMap`\nwith the same elements.", false, true, true, value.ToSymbol("=~"), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the element under the given key.\nReturns `nil` when the key is not present.", false, false, true, value.ToSymbol("[]"), []*Parameter{NewParameter(value.ToSymbol("key"), NewNamedType("Std::HashRecord::Key", Any{}), NormalParameterKind, false)}, NewNilable(NewNamedType("Std::HashRecord::Value", Any{})), nil)
				namespace.DefineMethod("Check whether the given `pair` is present in this record.", false, false, true, value.ToSymbol("contains"), []*Parameter{NewParameter(value.ToSymbol("pair"), NameToType("Std::Pair", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Check whether the given `key` is present in this record.", false, false, true, value.ToSymbol("contains_key"), []*Parameter{NewParameter(value.ToSymbol("key"), NewNamedType("Std::HashRecord::Key", Any{}), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Check whether the given `value` is present in this record.", false, false, true, value.ToSymbol("contains_value"), []*Parameter{NewParameter(value.ToSymbol("value"), NewNamedType("Std::HashRecord::Value", Any{}), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the record.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::HashRecord::Iterator", env), nil)
				namespace.DefineMethod("Returns the number of key-value pairs present in the record.", false, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, false, true, value.ToSymbol("map"), nil, Void{}, nil)
				namespace.DefineMethod("", false, false, true, value.ToSymbol("map_values"), nil, Void{}, nil)

				// Define constants

				// Define instance variables

				{
					namespace := namespace.SubtypeString("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins

					// Implement interfaces

					// Define methods
					namespace.DefineMethod("Returns `self`.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::HashRecord::Iterator", env), nil)
					namespace.DefineMethod("Get the next pair of the record.\nThrows `:stop_iteration` when there are no more elements.", false, false, true, value.ToSymbol("next"), nil, NewNamedType("Std::HashRecord::Iterator::Element", Any{}), NewSymbolLiteral("stop_iteration"))

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.SubtypeString("HashSet").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Return the intersection of both sets.\n\nCreate a new `HashSet` containing only the elements\npresent both in `self` and `other`.", false, true, true, value.ToSymbol("&"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::HashSet", env), NormalParameterKind, false)}, NameToType("Std::HashSet", env), nil)
				namespace.DefineMethod("Return the union of both sets.\n\nCreate a new `HashSet` containing all the elements\npresent in `self` and `other`.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::HashSet", env), NormalParameterKind, false)}, NameToType("Std::HashSet", env), nil)
				namespace.DefineMethod("Adds the given value to the set.\n\nDoes nothing if the value is already present in the set.\n\nReallocates the underlying array if it is\ntoo small to hold it.", false, false, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("value"), NewNamedType("Std::HashSet::Element", Any{}), NormalParameterKind, false)}, NameToType("Std::HashSet", env), nil)
				namespace.DefineMethod("Check whether the given value is a `HashSet`\nwith the same elements.", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Check whether the given value is a `HashSet`\nwith the same elements.", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the element under the given index.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", false, false, true, value.ToSymbol("[]"), []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NewNamedType("Std::HashSet::Element", Any{}), nil)
				namespace.DefineMethod("Set the element under the given index to the given value.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", false, false, true, value.ToSymbol("[]="), []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NewNamedType("Std::HashSet::Element", Any{}), NormalParameterKind, false)}, Void{}, nil)
				namespace.DefineMethod("Adds the given values to the set.\n\nSkips a value if it is already present in the set.\n\nReallocates the underlying array if it is\ntoo small to hold them.", false, false, true, value.ToSymbol("append"), []*Parameter{NewParameter(value.ToSymbol("values"), NewNamedType("Std::HashSet::Element", Any{}), PositionalRestParameterKind, false)}, NameToType("Std::HashSet", env), nil)
				namespace.DefineMethod("Returns the number of elements that can be\nheld by the underlying array.\n\nThis value will change when the set gets resized,\nand the underlying array gets reallocated.", false, false, true, value.ToSymbol("capacity"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Check whether the given `value` is present in this set.", false, false, true, value.ToSymbol("contains"), []*Parameter{NewParameter(value.ToSymbol("value"), NewNamedType("Std::HashSet::Element", Any{}), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Mutates the set.\n\nReallocates the underlying array to hold\nthe given number of new elements.\n\nExpands the `capacity` of the list by `new_slots`", false, false, true, value.ToSymbol("grow"), []*Parameter{NewParameter(value.ToSymbol("new_slots"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::HashSet", env), nil)
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the set.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::HashSet::Iterator", env), nil)
				namespace.DefineMethod("Returns the number of left slots for new elements\nin the underlying array.\nIt tells you how many more elements can be\nadded to the set before the underlying array gets\nreallocated.\n\nIt is always equal to `capacity - length`.", false, false, true, value.ToSymbol("left_capacity"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns the number of elements present in the set.", false, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, false, true, value.ToSymbol("map"), nil, Void{}, nil)
				namespace.DefineMethod("Return the union of both sets.\n\nCreate a new `HashSet` containing all the elements\npresent in `self` and `other`.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::HashSet", env), NormalParameterKind, false)}, NameToType("Std::HashSet", env), nil)
				namespace.DefineMethod("Return the union of both sets.\n\nCreate a new `HashSet` containing all the elements\npresent in `self` and `other`.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::HashSet", env), NormalParameterKind, false)}, NameToType("Std::HashSet", env), nil)

				// Define constants

				// Define instance variables

				{
					namespace := namespace.SubtypeString("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins

					// Implement interfaces

					// Define methods
					namespace.DefineMethod("Returns `self`.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::HashSet::Iterator", env), nil)
					namespace.DefineMethod("Get the next element of the set.\nThrows `:stop_iteration` when there are no more elements.", false, false, true, value.ToSymbol("next"), nil, NewNamedType("Std::HashSet::Iterator::Element", Any{}), NewSymbolLiteral("stop_iteration"))

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.SubtypeString("Inspectable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns a human readable `String`\nrepresentation of this value\nfor debugging etc.", true, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Int").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10\n\tvar b = 3\n\ta % b #=> 1\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Performs bitwise AND.", false, true, true, value.ToSymbol("&"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", false, true, true, value.ToSymbol("&~"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Multiply this integer by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Add `other` to this integer.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", false, true, true, value.ToSymbol("++"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1\n\t+a #=> 1\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Subtract `other` from this integer.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1\n\t-a #=> -1\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Performs bitwise XOR.", false, true, true, value.ToSymbol("^"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the integer to a floating point number.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Converts the integer to a string.", false, false, true, value.ToSymbol("to_string"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise OR.", false, true, true, value.ToSymbol("|"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4 #=> -5\n```", false, true, true, value.ToSymbol("~"), nil, NameToType("Std::Int", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Int16").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10i16\n\tvar b = 3i16\n\ta % b #=> 1i16\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Performs bitwise AND.", false, true, true, value.ToSymbol("&"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", false, true, true, value.ToSymbol("&~"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Multiply this integer by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Add `other` to this integer.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", false, true, true, value.ToSymbol("++"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1i16\n\t+a #=> 1i16\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Subtract `other` from this integer.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1i16\n\t-a #=> -1i16\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted arithmetically to the left by `other` positions, or to the right if `other` is negative.\n\nPreserves the integer's sign bit.\n\n4i16  << 1  #=> 8i16\n4i16  << -1 #=> 2i16\n-4i16 << 1  #=> -8i16\n-4i16 << -1 #=> -2i16", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Returns an integer shifted logically to the left by `other` positions, or to the right if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i16  <<< 1  #=> 8i16\n4i16  <<< -1 #=> 2i16\n-4i16 <<< 1  #=> -8i16\n-4i16 <<< -1 #=> 32766i16\n```", false, true, true, value.ToSymbol("<<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted arithmetically to the right by `other` positions, or to the left if `other` is negative.\n\nPreserves the integer's sign bit.\n\n```\n4i16  >> 1  #=> 2i16\n4i16  >> -1 #=> 8i16\n-4i16 >> 1  #=> -2i16\n-4i16 >> -1 #=> -8i16\n```", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Returns an integer shifted logically the the right by `other` positions, or to the left if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i16  >>> 1  #=> 2i16\n4i16  >>> -1 #=> 8i16\n-4i16 >>> 1  #=> 32766i16\n-4i16 >>> -1 #=> -8i16\n```", false, true, true, value.ToSymbol(">>>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Performs bitwise XOR.", false, true, true, value.ToSymbol("^"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the integer to a floating point number.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Converts to an automatically resizable integer type.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise OR.", false, true, true, value.ToSymbol("|"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4i16 #=> -5i16\n```", false, true, true, value.ToSymbol("~"), nil, NameToType("Std::Int16", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Int32").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10i32\n\tvar b = 3i32\n\ta % b #=> 1i32\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Performs bitwise AND.", false, true, true, value.ToSymbol("&"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", false, true, true, value.ToSymbol("&~"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Multiply this integer by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Add `other` to this integer.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", false, true, true, value.ToSymbol("++"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1i32\n\t+a #=> 1i32\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Subtract `other` from this integer.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1i32\n\t-a #=> -1i32\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted arithmetically to the left by `other` positions, or to the right if `other` is negative.\n\nPreserves the integer's sign bit.\n\n4i32  << 1  #=> 8i32\n4i32  << -1 #=> 2i32\n-4i32 << 1  #=> -8i32\n-4i32 << -1 #=> -2i32", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Returns an integer shifted logically to the left by `other` positions, or to the right if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i32  <<< 1  #=> 8i32\n4i32  <<< -1 #=> 2i32\n-4i32 <<< 1  #=> -8i32\n-4i32 <<< -1 #=> 2147483646i32\n```", false, true, true, value.ToSymbol("<<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted arithmetically to the right by `other` positions, or to the left if `other` is negative.\n\nPreserves the integer's sign bit.\n\n```\n4i32  >> 1  #=> 2i32\n4i32  >> -1 #=> 8i32\n-4i32 >> 1  #=> -2i32\n-4i32 >> -1 #=> -8i32\n```", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Returns an integer shifted logically the the right by `other` positions, or to the left if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i32  >>> 1  #=> 2i32\n4i32  >>> -1 #=> 8i32\n-4i32 >>> 1  #=> 2147483646i32\n-4i32 >>> -1 #=> -8i32\n```", false, true, true, value.ToSymbol(">>>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Performs bitwise XOR.", false, true, true, value.ToSymbol("^"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the integer to a floating point number.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Converts to an automatically resizable integer type.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise OR.", false, true, true, value.ToSymbol("|"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4i32 #=> -5i32\n```", false, true, true, value.ToSymbol("~"), nil, NameToType("Std::Int32", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Int64").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10i64\n\tvar b = 3i64\n\ta % b #=> 1i64\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Performs bitwise AND.", false, true, true, value.ToSymbol("&"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", false, true, true, value.ToSymbol("&~"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Multiply this integer by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Add `other` to this integer.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", false, true, true, value.ToSymbol("++"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1i64\n\t+a #=> 1i64\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Subtract `other` from this integer.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1i64\n\t-a #=> -1i64\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted arithmetically to the left by `other` positions, or to the right if `other` is negative.\n\nPreserves the integer's sign bit.\n\n4i64  << 1  #=> 8i64\n4i64  << -1 #=> 2i64\n-4i64 << 1  #=> -8i64\n-4i64 << -1 #=> -2i64", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Returns an integer shifted logically to the left by `other` positions, or to the right if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i64  <<< 1  #=> 8i64\n4i64  <<< -1 #=> 2i64\n-4i64 <<< 1  #=> -8i64\n-4i64 <<< -1 #=> 9223372036854775806i64\n```", false, true, true, value.ToSymbol("<<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted arithmetically to the right by `other` positions, or to the left if `other` is negative.\n\nPreserves the integer's sign bit.\n\n```\n4i64  >> 1  #=> 2i64\n4i64  >> -1 #=> 8i64\n-4i64 >> 1  #=> -2i64\n-4i64 >> -1 #=> -8i64\n```", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Returns an integer shifted logically the the right by `other` positions, or to the left if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i64  >>> 1  #=> 2i64\n4i64  >>> -1 #=> 8i64\n-4i64 >>> 1  #=> 9223372036854775806i64\n-4i64 >>> -1 #=> -8i64\n```", false, true, true, value.ToSymbol(">>>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Performs bitwise XOR.", false, true, true, value.ToSymbol("^"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the integer to a floating point number.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Converts to an automatically resizable integer type.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise OR.", false, true, true, value.ToSymbol("|"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4i64 #=> -5i64\n```", false, true, true, value.ToSymbol("~"), nil, NameToType("Std::Int64", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Int8").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10i8\n\tvar b = 3i8\n\ta % b #=> 1i8\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Performs bitwise AND.", false, true, true, value.ToSymbol("&"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", false, true, true, value.ToSymbol("&~"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Multiply this integer by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Add `other` to this integer.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", false, true, true, value.ToSymbol("++"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1i8\n\t+a #=> 1i8\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Subtract `other` from this integer.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1i8\n\t-a #=> -1i8\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted arithmetically to the left by `other` positions, or to the right if `other` is negative.\n\nPreserves the integer's sign bit.\n\n4i8  << 1  #=> 8i8\n4i8  << -1 #=> 2i8\n-4i8 << 1  #=> -8i8\n-4i8 << -1 #=> -2i8", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Returns an integer shifted logically to the left by `other` positions, or to the right if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i8  <<< 1  #=> 8i8\n4i8  <<< -1 #=> 2i8\n-4i8 <<< 1  #=> -8i8\n-4i8 <<< -1 #=> 126i8\n```", false, true, true, value.ToSymbol("<<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted arithmetically to the right by `other` positions, or to the left if `other` is negative.\n\nPreserves the integer's sign bit.\n\n```\n4i8  >> 1  #=> 2i8\n4i8  >> -1 #=> 8i8\n-4i8 >> 1  #=> -2i8\n-4i8 >> -1 #=> -8i8\n```", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Returns an integer shifted logically the the right by `other` positions, or to the left if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i8  >>> 1  #=> 2i8\n4i8  >>> -1 #=> 8i8\n-4i8 >>> 1  #=> 126i8\n-4i8 >>> -1 #=> -8i8\n```", false, true, true, value.ToSymbol(">>>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Performs bitwise XOR.", false, true, true, value.ToSymbol("^"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the integer to a floating point number.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Converts to an automatically resizable integer type.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise OR.", false, true, true, value.ToSymbol("|"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4i8 #=> -5i8\n```", false, true, true, value.ToSymbol("~"), nil, NameToType("Std::Int8", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Iterable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns an iterator for this structure.", true, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::Iterator", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Iterator").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces
				namespace.ImplementInterface(NameToType("Std::Iterable", env).(*Interface))

				// Define methods
				namespace.DefineMethod("", true, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::Iterator", env), nil)
				namespace.DefineMethod("Returns the next element.\nThrows `:stop_iteration` when no more elements are available.", true, false, true, value.ToSymbol("next"), nil, NewNamedType("Std::Iterator::Element", Any{}), NewSymbolLiteral("stop_iteration"))

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("List").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Adds the given value to the list.", true, false, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("value"), NewNamedType("Std::List::Element", Any{}), NormalParameterKind, false)}, NameToType("Std::ArrayList", env), nil)
				namespace.DefineMethod("Check whether the given value is a list\nwith the same elements.", true, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Check whether the given value is a list\nwith the same elements.", true, true, true, value.ToSymbol("=~"), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the element under the given index.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", true, false, true, value.ToSymbol("[]"), []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NewNamedType("Std::List::Element", Any{}), nil)
				namespace.DefineMethod("Set the element under the given index to the given value.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", true, false, true, value.ToSymbol("[]="), []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NewNamedType("Std::List::Element", Any{}), NormalParameterKind, false)}, Void{}, nil)
				namespace.DefineMethod("Adds the given values to the list.", true, false, true, value.ToSymbol("append"), []*Parameter{NewParameter(value.ToSymbol("values"), NewNamedType("Std::List::Element", Any{}), PositionalRestParameterKind, false)}, NameToType("Std::ArrayList", env), nil)
				namespace.DefineMethod("Check whether the given `value` is present in this list.", true, false, true, value.ToSymbol("contains"), []*Parameter{NewParameter(value.ToSymbol("value"), NewNamedType("Std::List::Element", Any{}), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the list.", true, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::Iterator", env), nil)
				namespace.DefineMethod("Returns the number of elements present in the list.", true, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", true, false, true, value.ToSymbol("map"), nil, Void{}, nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Object").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins

				// Implement interfaces

				// Define methods

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("OutOfRangeError").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Error", env))

				// Include mixins

				// Implement interfaces

				// Define methods

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Pair").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Instantiate the `Pair` with the given key and value.", false, false, true, value.ToSymbol("#init"), []*Parameter{NewParameter(value.ToSymbol("key"), NewNamedType("Std::Pair::Key", Any{}), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NewNamedType("Std::Pair::Value", Any{}), NormalParameterKind, false)}, Void{}, nil)
				namespace.DefineMethod("Check whether the given value\nis a `Pair` that is equal to this `Pair`.", false, false, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the element with the given index.\nThe key is `0`, value is `1`.", false, false, true, value.ToSymbol("[]"), []*Parameter{NewParameter(value.ToSymbol("key"), NewNamedType("Std::Pair::Key", Any{}), NormalParameterKind, false)}, NewNamedType("Std::Pair::Value", Any{}), nil)
				namespace.DefineMethod("Set the element with the given index to the given value.\nThe key is `0`, value is `1`.", false, false, true, value.ToSymbol("[]="), []*Parameter{NewParameter(value.ToSymbol("key"), NewNamedType("Std::Pair::Key", Any{}), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NewNamedType("Std::Pair::Value", Any{}), NormalParameterKind, false)}, Void{}, nil)
				namespace.DefineMethod("Returns the key, the first element of the tuple.", false, false, true, value.ToSymbol("key"), nil, NewNamedType("Std::Pair::Key", Any{}), nil)
				namespace.DefineMethod("Always returns `2`.\nFor compatibility with `Tuple`.", false, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns the value, the second element of the tuple.", false, false, true, value.ToSymbol("value"), nil, NewNamedType("Std::Pair::Value", Any{}), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Regex").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Creates a new `Regex` that contains the\npattern of `self` repeated `n` times.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Create a new regex that contains\nthe patterns present in both operands.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Regex", env), NormalParameterKind, false)}, NameToType("Std::Regex", env), nil)
				namespace.DefineMethod("Check whether the pattern matches\nthe given string.\n\nReturns `true` if it matches, otherwise `false`.", false, false, true, value.ToSymbol("matches"), []*Parameter{NewParameter(value.ToSymbol("str"), NameToType("Std::String", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Creates a new string with this character.", false, false, true, value.ToSymbol("to_string"), nil, NameToType("Std::String", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("String").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Concatenate this `String`\nwith another `String` or `Char`.\n\nCreates a new `String` containing the content\nof both operands.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Remove the given suffix from the `String`.\n\nDoes nothing if the `String` doesn't end\nwith `suffix` and returns `self`.\n\nIf the `String` ends with the given suffix\na new `String` gets created and returned that doesn't contain\nthe suffix.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("suffix"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Object", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the byte with the given index.\nIndices start at 0.", false, false, true, value.ToSymbol("byte_at"), []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Get the number of bytes that this\nstring contains.", false, false, true, value.ToSymbol("byte_count"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Iterates over all bytes of a `String`.", false, false, true, value.ToSymbol("byte_iterator"), nil, NameToType("Std::String::ByteIterator", env), nil)
				namespace.DefineMethod("Get the number of Unicode code points\nthat this `String` contains.", false, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Iterates over all unicode code points of a `String`.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::String::CharIterator", env), nil)
				namespace.DefineMethod("Get the Unicode code point with the given index.\nIndices start at 0.", false, false, true, value.ToSymbol("chat_at"), []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Char", env), nil)
				namespace.DefineMethod("Concatenate this `String`\nwith another `String` or `Char`.\n\nCreates a new `String` containing the content\nof both operands.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Get the Unicode grapheme cluster with the given index.\nIndices start at 0.", false, false, true, value.ToSymbol("grapheme_at"), []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Get the number of unicode grapheme clusters\npresent in this string.", false, false, true, value.ToSymbol("grapheme_count"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Check whether the `String` is empty.", false, false, true, value.ToSymbol("is_empty"), nil, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Iterates over all unicode code points of a `String`.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::String::CharIterator", env), nil)
				namespace.DefineMethod("Get the number of Unicode code points\nthat this `String` contains.", false, false, true, value.ToSymbol("length"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Create a new string with all of the characters\nof this one turned into lowercase.", false, false, true, value.ToSymbol("lowercase"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Remove the given suffix from the `String`.\n\nDoes nothing if the `String` doesn't end\nwith `suffix` and returns `self`.\n\nIf the `String` ends with the given suffix\na new `String` gets created and returned that doesn't contain\nthe suffix.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("suffix"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_string"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Convert the `String` to a `Symbol`.", false, false, true, value.ToSymbol("to_symbol"), nil, NameToType("Std::Symbol", env), nil)
				namespace.DefineMethod("Create a new string with all of the characters\nof this one turned into uppercase.", false, false, true, value.ToSymbol("uppercase"), nil, NameToType("Std::String", env), nil)

				// Define constants

				// Define instance variables

				{
					namespace := namespace.SubtypeString("ByteIterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins

					// Implement interfaces

					// Define methods
					namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::String::ByteIterator", env), nil)
					namespace.DefineMethod("Get the next byte.\nThrows `:stop_iteration` when no more bytes are available.", false, false, true, value.ToSymbol("next"), nil, NameToType("Std::UInt8", env), NewSymbolLiteral("stop_iteration"))

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.SubtypeString("CharIterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins

					// Implement interfaces

					// Define methods
					namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("iterator"), nil, NameToType("Std::String::CharIterator", env), nil)
					namespace.DefineMethod("Get the next character.\nThrows `:stop_iteration` when no more characters are available.", false, false, true, value.ToSymbol("next"), nil, NameToType("Std::Char", env), NewSymbolLiteral("stop_iteration"))

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.SubtypeString("StringConvertible").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Convert the value to a string.", true, false, true, value.ToSymbol("to_string"), nil, NameToType("Std::String", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Symbol").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Object", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Returns the string associated with this symbol.", false, false, true, value.ToSymbol("to_string"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Returns the string associated with this symbol.", false, false, true, value.ToSymbol("to_string"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_symbol"), nil, NameToType("Std::Symbol", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("True").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Bool", env))

				// Include mixins

				// Implement interfaces

				// Define methods

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("UInt16").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10u16\n\tvar b = 3u16\n\ta % b #=> 1u16\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Performs bitwise AND.", false, true, true, value.ToSymbol("&"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", false, true, true, value.ToSymbol("&~"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Multiply this integer by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Add `other` to this integer.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", false, true, true, value.ToSymbol("++"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1u16\n\t+a #=> 1u16\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Subtract `other` from this integer.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1u16\n\t-a #=> 65535u16\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Performs bitwise XOR.", false, true, true, value.ToSymbol("^"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the integer to a floating point number.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Converts to an automatically resizable integer type.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise OR.", false, true, true, value.ToSymbol("|"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4u16 #=> 65531u16\n```", false, true, true, value.ToSymbol("~"), nil, NameToType("Std::UInt16", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("UInt32").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10u32\n\tvar b = 3u32\n\ta % b #=> 1u32\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Performs bitwise AND.", false, true, true, value.ToSymbol("&"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", false, true, true, value.ToSymbol("&~"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Multiply this integer by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Add `other` to this integer.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", false, true, true, value.ToSymbol("++"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1u32\n\t+a #=> 1u32\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Subtract `other` from this integer.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1u32\n\t-a #=> 4294967295u32\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Performs bitwise XOR.", false, true, true, value.ToSymbol("^"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the integer to a floating point number.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Converts to an automatically resizable integer type.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise OR.", false, true, true, value.ToSymbol("|"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4u32 #=> 4294967291u32\n```", false, true, true, value.ToSymbol("~"), nil, NameToType("Std::UInt32", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("UInt64").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10u64\n\tvar b = 3u64\n\ta % b #=> 1u64\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Performs bitwise AND.", false, true, true, value.ToSymbol("&"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", false, true, true, value.ToSymbol("&~"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Multiply this integer by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Add `other` to this integer.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", false, true, true, value.ToSymbol("++"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1u64\n\t+a #=> 1u64\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Subtract `other` from this integer.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1u64\n\t-a #=> 18446744073709551615u64\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Performs bitwise XOR.", false, true, true, value.ToSymbol("^"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the integer to a floating point number.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Converts to an automatically resizable integer type.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise OR.", false, true, true, value.ToSymbol("|"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4u64 #=> 18446744073709551611u64\n```", false, true, true, value.ToSymbol("~"), nil, NameToType("Std::UInt64", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("UInt8").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10u8\n\tvar b = 3u8\n\ta % b #=> 1u8\n```", false, true, true, value.ToSymbol("%"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise AND.", false, true, true, value.ToSymbol("&"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", false, true, true, value.ToSymbol("&~"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Multiply this integer by `other`.", false, true, true, value.ToSymbol("*"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", false, true, true, value.ToSymbol("**"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Add `other` to this integer.", false, true, true, value.ToSymbol("+"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", false, true, true, value.ToSymbol("++"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1u8\n\t+a #=> 1u8\n```", false, true, true, value.ToSymbol("+@"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Subtract `other` from this integer.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1u8\n\t-a #=> 255u8\n```", false, true, true, value.ToSymbol("-@"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", false, true, true, value.ToSymbol(">>"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise XOR.", false, true, true, value.ToSymbol("^"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Converts the integer to a floating point number.", false, false, true, value.ToSymbol("to_float"), nil, NameToType("Std::Float", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", false, false, true, value.ToSymbol("to_float32"), nil, NameToType("Std::Float32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", false, false, true, value.ToSymbol("to_float64"), nil, NameToType("Std::Float64", env), nil)
				namespace.DefineMethod("Converts to an automatically resizable integer type.", false, false, true, value.ToSymbol("to_int"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", false, false, true, value.ToSymbol("to_int16"), nil, NameToType("Std::Int16", env), nil)
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", false, false, true, value.ToSymbol("to_int32"), nil, NameToType("Std::Int32", env), nil)
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", false, false, true, value.ToSymbol("to_int64"), nil, NameToType("Std::Int64", env), nil)
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", false, false, true, value.ToSymbol("to_int8"), nil, NameToType("Std::Int8", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise OR.", false, true, true, value.ToSymbol("|"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4u8 #=> 251u8\n```", false, true, true, value.ToSymbol("~"), nil, NameToType("Std::UInt8", env), nil)

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.SubtypeString("Value").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Compares this value with another value.\n\nReturns `true` when they are instances of the same class,\nand are equal.", false, false, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Compares this value with another value.\nReturns `true` when they are equal.\n\nInstances of different (but similar) classes\nmay be treated as equal.", false, false, true, value.ToSymbol("=~"), []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Returns a human readable `String`\nrepresentation of this value\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)

				// Define constants

				// Define instance variables
			}
		}
	}
}
