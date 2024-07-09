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
		namespace := namespace.TryDefineModule("", "Std")
		namespace.TryDefineClass("", false, false, false, "Object", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Float64", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "UInt64", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Method", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Value", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Int64", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Int32", objectClass, env)
		namespace.TryDefineClass("A base class for most errors in Elk stdlib.", false, false, false, "Error", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "True", objectClass, env)
		{
			namespace := namespace.TryDefineClass("", false, true, true, "String", objectClass, env)
			namespace.TryDefineClass("Iterates over all unicode code points of a `String`.", false, false, false, "CharIterator", objectClass, env)
			namespace.TryDefineClass("Iterates over all bytes of a `String`.", false, false, false, "ByteIterator", objectClass, env)
		}
		namespace.TryDefineClass("", false, true, true, "Int", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Int16", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Int8", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "ArrayTuple", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Interface", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "BigFloat", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "UInt32", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "HashSet", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Regex", objectClass, env)
		namespace.TryDefineClass("Thrown when a numeric value is too large or too small to be used in a particular setting.", false, false, false, "OutOfRangeError", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "False", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Symbol", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "UInt8", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Module", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Nil", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Bool", objectClass, env)
		namespace.TryDefineClass("Represents a single Unicode code point.", false, true, true, "Char", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Float32", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "UInt16", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "ArrayList", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Class", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Float", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "HashMap", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "HashRecord", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Mixin", objectClass, env)
	}

	// Define methods, constants

	{
		namespace := env.Root

		namespace.Name() // noop - avoid unused variable error

		// Include mixins

		// Implement interfaces

		// Define methods

		// Define constants

		{
			namespace := namespace.SubtypeString("Std").(*Module)

			namespace.Name() // noop - avoid unused variable error

			// Include mixins

			// Implement interfaces

			// Define methods

			// Define constants

			{
				namespace := namespace.SubtypeString("Object").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins

				// Implement interfaces

				// Define methods

				// Define constants
			}
			{
				namespace := namespace.SubtypeString("True").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Bool", env))

				// Include mixins

				// Implement interfaces

				// Define methods

				// Define constants
			}
			{
				namespace := namespace.SubtypeString("String").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Check whether the `String` is empty.", false, false, true, "is_empty", nil, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, "<=", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the number of unicode grapheme clusters\npresent in this string.", false, false, true, "grapheme_count", nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Returns itself.", false, false, true, "to_string", nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Convert the `String` to a `Symbol`.", false, false, true, "to_symbol", nil, NameToType("Std::Symbol", env), nil)
				namespace.DefineMethod("Iterates over all unicode code points of a `String`.", false, false, true, "iterator", nil, NameToType("Std::String::CharIterator", env), nil)
				namespace.DefineMethod("Get the number of bytes that this\nstring contains.", false, false, true, "byte_count", nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Iterates over all bytes of a `String`.", false, false, true, "byte_iterator", nil, NameToType("Std::String::ByteIterator", env), nil)
				namespace.DefineMethod("Concatenate this `String`\nwith another `String` or `Char`.\n\nCreates a new `String` containing the content\nof both operands.", false, true, true, "+", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Concatenate this `String`\nwith another `String` or `Char`.\n\nCreates a new `String` containing the content\nof both operands.", false, true, true, "+", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Remove the given suffix from the `String`.\n\nDoes nothing if the `String` doesn't end\nwith `suffix` and returns `self`.\n\nIf the `String` ends with the given suffix\na new `String` gets created and returned that doesn't contain\nthe suffix.", false, true, true, "-", []*Parameter{NewParameter(value.ToSymbol("suffix"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("", false, true, true, ">", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the byte with the given index.\nIndices start at 0.", false, false, true, "byte_at", []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::UInt8", env), nil)
				namespace.DefineMethod("", false, true, true, "==", []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Object", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Return a human readable `String`\nrepresentation of this object\nfor debugging etc.", false, false, true, "inspect", nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Remove the given suffix from the `String`.\n\nDoes nothing if the `String` doesn't end\nwith `suffix` and returns `self`.\n\nIf the `String` ends with the given suffix\na new `String` gets created and returned that doesn't contain\nthe suffix.", false, true, true, "-", []*Parameter{NewParameter(value.ToSymbol("suffix"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("", false, false, true, "<=>", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Iterates over all unicode code points of a `String`.", false, false, true, "iterator", nil, NameToType("Std::String::CharIterator", env), nil)
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", false, true, true, "*", []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Get the Unicode grapheme cluster with the given index.\nIndices start at 0.", false, false, true, "grapheme_at", []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Create a new string with all of the characters\nof this one turned into uppercase.", false, false, true, "uppercase", nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Create a new string with all of the characters\nof this one turned into lowercase.", false, false, true, "lowercase", nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Get the number of Unicode code points\nthat this `String` contains.", false, false, true, "length", nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, "<", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the Unicode code point with the given index.\nIndices start at 0.", false, false, true, "chat_at", []*Parameter{NewParameter(value.ToSymbol("index"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))), NormalParameterKind, false)}, NameToType("Std::Char", env), nil)
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", false, true, true, "*", []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("", false, true, true, ">=", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the number of Unicode code points\nthat this `String` contains.", false, false, true, "length", nil, NameToType("Std::Int", env), nil)

				// Define constants

				{
					namespace := namespace.SubtypeString("ByteIterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins

					// Implement interfaces

					// Define methods
					namespace.DefineMethod("Get the next byte.\nThrows `:stop_iteration` when no more bytes are available.", false, false, true, "next", nil, NameToType("Std::UInt8", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Returns itself.", false, false, true, "iterator", nil, NameToType("Std::String::ByteIterator", env), nil)

					// Define constants
				}
				{
					namespace := namespace.SubtypeString("CharIterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins

					// Implement interfaces

					// Define methods
					namespace.DefineMethod("Get the next character.\nThrows `:stop_iteration` when no more characters are available.", false, false, true, "next", nil, NameToType("Std::Char", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Returns itself.", false, false, true, "iterator", nil, NameToType("Std::String::CharIterator", env), nil)

					// Define constants
				}
			}
			{
				namespace := namespace.SubtypeString("OutOfRangeError").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Error", env))

				// Include mixins

				// Implement interfaces

				// Define methods

				// Define constants
			}
			{
				namespace := namespace.SubtypeString("False").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Bool", env))

				// Include mixins

				// Implement interfaces

				// Define methods

				// Define constants
			}
			{
				namespace := namespace.SubtypeString("Char").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("", false, true, true, "<=", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, ">", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Return the lowercase version of this character.", false, false, true, "lowercase", nil, NameToType("Std::Char", env), nil)
				namespace.DefineMethod("Converts the `Char` to a `Symbol`.", false, false, true, "to_symbol", nil, NameToType("Std::Symbol", env), nil)
				namespace.DefineMethod("Concatenate this `Char`\nwith another `Char` or `String`.\n\nCreates a new `String` containing the content\nof both operands.", false, true, true, "+", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("", false, false, true, "<=>", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("", false, true, true, ">=", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, "==", []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Object", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the number of bytes that this\ncharacter contains.", false, false, true, "byte_count", nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", false, false, true, "inspect", nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Concatenate this `Char`\nwith another `Char` or `String`.\n\nCreates a new `String` containing the content\nof both operands.", false, true, true, "+", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", false, true, true, "*", []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Always returns 1.", false, false, true, "length", nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Always returns 1.", false, false, true, "grapheme_count", nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Return the uppercase version of this character.", false, false, true, "uppercase", nil, NameToType("Std::Char", env), nil)
				namespace.DefineMethod("Always returns 1.", false, false, true, "length", nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Creates a new string with this character.", false, false, true, "to_string", nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Always returns false.", false, false, true, "is_empty", nil, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", false, true, true, "*", []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("", false, true, true, "<", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)

				// Define constants
			}
		}
	}
}
