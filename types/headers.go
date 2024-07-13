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
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("ArrayList"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("ArrayTuple"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("BigFloat"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Bool"), objectClass, env)
		namespace.TryDefineClass("Represents a single Unicode code point.", false, true, true, value.ToSymbol("Char"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Class"), objectClass, env)
		namespace.TryDefineClass("A base class for most errors in Elk stdlib.", false, false, false, value.ToSymbol("Error"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("False"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Float"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Float32"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Float64"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("HashMap"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("HashRecord"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("HashSet"), objectClass, env)
		namespace.TryDefineClass("Represents an integer (a whole number like `1`, `2`, `3`, `-5`, `0`).\n\nThis integer type is automatically resized so\nit can hold an arbitrarily large/small number.", false, true, true, value.ToSymbol("Int"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Int16"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Int32"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Int64"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Int8"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Interface"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Method"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Mixin"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Module"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Nil"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Object"), objectClass, env)
		namespace.TryDefineClass("Thrown when a numeric value is too large or too small to be used in a particular setting.", false, false, false, value.ToSymbol("OutOfRangeError"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("Regex"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("", false, true, true, value.ToSymbol("String"), objectClass, env)
			namespace.TryDefineClass("Iterates over all bytes of a `String`.", false, false, false, value.ToSymbol("ByteIterator"), objectClass, env)
			namespace.TryDefineClass("Iterates over all unicode code points of a `String`.", false, false, false, value.ToSymbol("CharIterator"), objectClass, env)
		}
		namespace.TryDefineClass("Represents an interned string.\n\nA symbol is an integer ID that is associated\nwith a particular name (string).\n\nA few symbols with the same name refer to the same ID.\n\nComparing symbols happens in constant time, so it's\nusually faster than comparing strings.", false, true, true, value.ToSymbol("Symbol"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("True"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("UInt16"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("UInt32"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("UInt64"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, value.ToSymbol("UInt8"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, value.ToSymbol("Value"), objectClass, env)
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
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1\n\t+a #=> 1\n```", false, true, true, value.ToSymbol("+@"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Subtract `other` from this integer.", false, true, true, value.ToSymbol("-"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", false, true, true, value.ToSymbol("--"), nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1\n\t-a #=> -1\n```", false, true, true, value.ToSymbol("-@"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", false, true, true, value.ToSymbol("/"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", false, true, true, value.ToSymbol("<<"), []*Parameter{NewParameter(value.ToSymbol("other"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("<="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", false, true, true, value.ToSymbol("<=>"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol("=="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, value.ToSymbol(">="), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
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
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", false, false, true, value.ToSymbol("to_uint16"), nil, NameToType("Std::UInt16", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", false, false, true, value.ToSymbol("to_uint32"), nil, NameToType("Std::UInt32", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", false, false, true, value.ToSymbol("to_uint64"), nil, NameToType("Std::UInt64", env), nil)
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", false, false, true, value.ToSymbol("to_uint8"), nil, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("Performs bitwise OR.", false, true, true, value.ToSymbol("|"), []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)

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
				namespace.DefineMethod("Return a human readable `String`\nrepresentation of this object\nfor debugging etc.", false, false, true, value.ToSymbol("inspect"), nil, NameToType("Std::String", env), nil)
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
		}
	}
}
