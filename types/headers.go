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
		namespace.TryDefineClass("", false, true, true, "Bool", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "True", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Int", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Int64", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "HashMap", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Module", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "False", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Symbol", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Float", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "ArrayList", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "HashSet", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Class", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Mixin", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Interface", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "BigFloat", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "UInt16", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Float64", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Int32", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Int8", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "ArrayTuple", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "HashRecord", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Nil", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Char", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Regex", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Float32", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "UInt64", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "String", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Int16", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Value", objectClass, env)
		namespace.TryDefineClass("", false, false, false, "Object", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "UInt32", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "UInt8", objectClass, env)
		namespace.TryDefineClass("", false, true, true, "Method", objectClass, env)
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
				namespace := namespace.SubtypeString("True").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Bool", env))

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
				namespace := namespace.SubtypeString("String").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins

				// Implement interfaces

				// Define methods
				namespace.DefineMethod("Concatenate this `String`\nwith another `String` or `Char`.\n\nCreates a new `String` containing the content\nof both operands.", false, true, true, "+", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Concatenate this `String`\nwith another `String` or `Char`.\n\nCreates a new `String` containing the content\nof both operands.", false, true, true, "+", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Get the number of unicode code points\nthat this `String` contains.", false, false, true, "length", nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Get the number of unicode code points\nthat this `String` contains.", false, false, true, "length", nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Check whether the `String` is empty.", false, false, true, "is_empty", nil, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Remove the given suffix from the `String`.\n\nDoes nothing if the `String` doesn't end\nwith `suffix` and returns `self`.\n\nIf the `String` ends with the given suffix\na new `String` gets created and returned that doesn't contain\nthe suffix.", false, true, true, "-", []*Parameter{NewParameter(value.ToSymbol("suffix"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", false, true, true, "*", []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("", false, true, true, "<", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, ">=", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", false, true, true, "*", []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("", false, true, true, ">", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("", false, true, true, "==", []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Object", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the number of unicode grapheme clusters\npresent in the `String`.", false, false, true, "grapheme_count", nil, NameToType("Std::Int", env), nil)
				namespace.DefineMethod("Convert the `String` to a `Symbol`.", false, false, true, "to_symbol", nil, NameToType("Std::Symbol", env), nil)
				namespace.DefineMethod("Return a human readable `String`\nrepresentation of this object\nfor debugging etc.", false, false, true, "inspect", nil, NameToType("Std::String", env), nil)
				namespace.DefineMethod("Remove the given suffix from the `String`.\n\nDoes nothing if the `String` doesn't end\nwith `suffix` and returns `self`.\n\nIf the `String` ends with the given suffix\na new `String` gets created and returned that doesn't contain\nthe suffix.", false, true, true, "-", []*Parameter{NewParameter(value.ToSymbol("suffix"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), nil)
				namespace.DefineMethod("", false, true, true, "<=", []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Bool", env), nil)
				namespace.DefineMethod("Get the number of bytes that this\n`String` contains.", false, false, true, "byte_count", nil, NameToType("Std::Int", env), nil)

				// Define constants
			}
			{
				namespace := namespace.SubtypeString("Object").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins

				// Implement interfaces

				// Define methods

				// Define constants
			}
		}
	}
}
