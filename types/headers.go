package types

// This file is auto-generated, please do not edit it manually

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

func setupGlobalEnvironmentFromHeaders(env *GlobalEnvironment) {
	objectClass := env.StdSubtypeClass(symbol.Object)
	namespace := env.Root
	var mixin *Mixin
	mixin.IsLiteral() // noop - avoid unused variable error

	// Define all namespaces
	namespace.DefineSubtype(value.ToSymbol("Byte"), NewNamedType("Byte", NameToType("Std::UInt8", env)))
	{
		namespace := namespace.TryDefineModule("", value.ToSymbol("Std"), env)
		namespace.DefineSubtype(value.ToSymbol("AnyFloat"), NewNamedType("Std::AnyFloat", NewUnion(NameToType("Std::Float", env), NameToType("Std::Float64", env), NameToType("Std::Float32", env), NameToType("Std::BigFloat", env))))
		namespace.DefineSubtype(value.ToSymbol("AnyInt"), NewNamedType("Std::AnyInt", NewUnion(NameToType("Std::Int", env), NameToType("Std::Int64", env), NameToType("Std::Int32", env), NameToType("Std::Int16", env), NameToType("Std::Int8", env), NameToType("Std::UInt64", env), NameToType("Std::UInt32", env), NameToType("Std::UInt16", env), NameToType("Std::UInt8", env))))
		{
			namespace := namespace.TryDefineClass("A dynamically resizable list data structure backed\nby an array.\n\nIt is an ordered collection of integer indexed values.", false, true, true, false, value.ToSymbol("ArrayList"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("A tuple data structure backed by an array.\n\nIt is an ordered, immutable collection of integer indexed values.\nA tuple is an immutable list.", false, true, true, false, value.ToSymbol("ArrayTuple"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents a closed range from -∞ to a given value *(-∞, end]*", false, true, true, false, value.ToSymbol("BeginlessClosedRange"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents an open range from -∞ to a given value *(-∞, end)*", false, true, true, false, value.ToSymbol("BeginlessOpenRange"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Represents a multi-precision floating point number (a fraction like `1.2`, `0.1`).\n\n```\nsign × mantissa × 2**exponent\n```\n\nwith 0.5 <= mantissa < 1.0, and MinExp <= exponent <= MaxExp.\nA `BigFloat` may also be zero (+0, -0) or infinite (+Inf, -Inf).\nAll BigFloats are ordered.\n\nBy setting the desired precision to 24 or 53,\n`BigFloat` operations produce the same results as the corresponding float32 or float64 IEEE-754 arithmetic for operands that\ncorrespond to normal (i.e., not denormal) `Float`, `Float32` and `Float64` numbers.\nExponent underflow and overflow lead to a `0` or an Infinity for different values than IEEE-754 because `BigFloat` exponents have a much larger range.", false, true, true, true, value.ToSymbol("BigFloat"), objectClass, env)
		namespace.TryDefineClass("Represents boolean values.\nThere are only two instances of `Bool`, `true` and `false`.", false, true, true, true, value.ToSymbol("Bool"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("Box wraps another value, it's a pointer to another `Value`.", false, true, true, false, value.ToSymbol("Box"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Represents a single function call in a stack trace.", false, true, true, true, value.ToSymbol("CallFrame"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("A `Channel` is an object tha can be used to send and receive values.\nIts useful for communicating between multiple threads of execution.\n\n## Instantiation\n\nYou can specify the capacity of the channel.\nA channel with `0` capacity is called an unbuffered channel.\nChannels with positive capacity are called buffered channel.\n\n```\n# instantiate an unbuffered channel of `String` values\nunbuffered_channel := Channel::[String]()\n\n# instantiate a buffered channel of `Int` values, that can hold up to 5 integers\nbuffered_channel := Channel::[Int](5)\n```\n\n## Pushing values\n\nYou can send values to the channel using the `<<` operator.\nUnbuffered channels will block the current thread until the pushed value\nis popped by another thread.\nBuffered channels will not block the current thread if there is enough capacity for another value.\n\n```\nch := Channel::[Int]() # instantiate a channel of `Int` values\nch << 5 # send `5` to the channel\n```\n\nPushing values to a closed channel will result in an unchecked error being thrown.\n\n## Popping values\n\nYou can receive values from the channel using the `pop` method.\nUnbuffered channels will block the current thread until a value is available.\nBuffered channels will not block the current thread if there is a value in the channel's buffer.\n\n```\nch := Channel::[Int](3) # instantiate a buffered channel of `Int` values\nch << 5 # send `5` to the channel\nv := try ch.pop # pop `5` from the channel\n```\n\nif the channel is closed `pop` will throw `:channel_closed`\n\n## Closing channels\n\nYou can close a channel using the `close` method when you no longer wish to send values to it.\nChannels should only be closed by the producer (the thread that pushes values to the channel).\nClosing a closed channel will result in an unchecked error being thrown.", false, true, true, false, value.ToSymbol("Channel"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Represents a single Unicode code point.", false, true, true, true, value.ToSymbol("Char"), objectClass, env)
		namespace.TryDefineClass("`Class` is a metaclass, it's the class of all classes.", false, false, false, true, value.ToSymbol("Class"), objectClass, env)
		namespace.TryDefineInterface("Represents a resource that can be closed.", value.ToSymbol("Closable"), env)
		{
			namespace := namespace.TryDefineClass("Represents a closed range from `start` to `end` *[start, end]*", false, true, true, false, value.ToSymbol("ClosedRange"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.DefineSubtype(value.ToSymbol("CoercibleNumeric"), NewNamedType("Std::CoercibleNumeric", NewUnion(NameToType("Std::Int", env), NameToType("Std::Float", env), NameToType("Std::BigFloat", env))))
		{
			namespace := namespace.TryDefineInterface("An interface that represents a finite, mutable collection\nof elements.", value.ToSymbol("Collection"), env)
			{
				namespace := namespace.TryDefineMixin("Provides default implementations of most collection methods.", true, value.ToSymbol("Base"), env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineInterface("Represents an object that can colorize\na string.", value.ToSymbol("Colorizer"), env)
		{
			namespace := namespace.TryDefineInterface("Represents a value that can be compared\nusing relational operators like `>`, `>=`, `<`, `<=`, `<=>`", value.ToSymbol("Comparable"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineInterface("Represents a data structure that\ncan be used to check if it contains\na value.", value.ToSymbol("Container"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineModule("Contains various debugging utilities.", value.ToSymbol("Debug"), env)
		{
			namespace := namespace.TryDefineInterface("Represents a value that can be decremented using\nthe `--` operator like `a--`", value.ToSymbol("Decrementable"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Represents an info message, a warning or a failure\nat a particular `Location`", false, true, true, false, value.ToSymbol("Diagnostic"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("A builtin list of diagnostics.", false, true, true, false, value.ToSymbol("DiagnosticList"), objectClass, env)
			namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Represents the elapsed time between two Times as an int64 nanosecond count.\n The representation limits the largest representable duration to approximately 290 years.", false, true, true, false, value.ToSymbol("Duration"), objectClass, env)
		{
			namespace := namespace.TryDefineModule("Contains utilities for dealing with Elk\nsource code.", value.ToSymbol("Elk"), env)
			{
				namespace := namespace.TryDefineModule("Contains the definitions of all\nElk AST (Abstract Syntax Tree) node types.", value.ToSymbol("AST"), env)
				namespace.TryDefineClass("A single alias entry eg. `new_name old_name`", false, true, true, false, value.ToSymbol("AliasDeclarationEntry"), objectClass, env)
				namespace.TryDefineClass("Represents a new alias declaration eg. `alias push append, add plus`", false, true, true, false, value.ToSymbol("AliasDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("`any` type.", false, true, true, false, value.ToSymbol("AnyTypeNode"), objectClass, env)
				namespace.TryDefineClass("Represents a ArrayList literal eg. `[1, 5, -6]`", false, true, true, false, value.ToSymbol("ArrayListLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a ArrayTuple literal eg. `%[1, 5, -6]`", false, true, true, false, value.ToSymbol("ArrayTupleLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents an as type downcast eg. `foo as String`", false, true, true, false, value.ToSymbol("AsExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents an as pattern eg. `> 5 && < 20 as foo`", false, true, true, false, value.ToSymbol("AsPatternNode"), objectClass, env)
				namespace.TryDefineClass("Assignment with the specified operator.", false, true, true, false, value.ToSymbol("AssignmentExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a new setter declaration eg. `attr foo: String`", false, true, true, false, value.ToSymbol("AttrDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents attribute access eg. `foo.bar`", false, true, true, false, value.ToSymbol("AttributeAccessNode"), objectClass, env)
				namespace.TryDefineClass("Represents an attribute declaration in getters, setters and accessors eg. `foo: String`", false, true, true, false, value.ToSymbol("AttributeParameterNode"), objectClass, env)
				namespace.TryDefineClass("Represents an `await` expression eg. `await foo()`", false, true, true, false, value.ToSymbol("AwaitExpressionNode"), objectClass, env)
				namespace.TryDefineClass("BigFloat literal eg. `5.2bf`, `.5bf`, `45e20bf`", false, true, true, false, value.ToSymbol("BigFloatLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a bin ArrayList literal eg. `\\b[11 10]`", false, true, true, false, value.ToSymbol("BinArrayListLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a bin ArrayTuple literal eg. `%b[11 10]`", false, true, true, false, value.ToSymbol("BinArrayTupleLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a bin HashSet literal eg. `^b[11 10]`", false, true, true, false, value.ToSymbol("BinHashSetLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Expression of an operator with two operands eg. `2 + 5`, `foo > bar`", false, true, true, false, value.ToSymbol("BinaryExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Pattern with two operands eg. `> 10 && < 50`", false, true, true, false, value.ToSymbol("BinaryPatternNode"), objectClass, env)
				namespace.TryDefineClass("Type expression of an operator with two operands eg. `String | Int`", false, true, true, false, value.ToSymbol("BinaryTypeNode"), objectClass, env)
				namespace.TryDefineClass("`bool` literal.", false, true, true, false, value.ToSymbol("BoolLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `break` expression eg. `break`, `break false`", false, true, true, false, value.ToSymbol("BreakExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a method call eg. `'123'.()`", false, true, true, false, value.ToSymbol("CallNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `case` node eg. `case 3 then println(\"eureka!\")`", false, true, true, false, value.ToSymbol("CaseNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `catch` eg.\n\n```\ncatch SomeError(message)\n  print(\"awesome!\")\nend\n```", false, true, true, false, value.ToSymbol("CatchNode"), objectClass, env)
				namespace.TryDefineClass("Char literal eg. `c\"a\"`", false, true, true, false, value.ToSymbol("CharLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a class declaration eg. `class Foo; end`", false, true, true, false, value.ToSymbol("ClassDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents a closure eg. `|i| -> println(i)`", false, true, true, false, value.ToSymbol("ClosureLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a closure type eg. `|i: Int|: String`", false, true, true, false, value.ToSymbol("ClosureTypeNode"), objectClass, env)
				namespace.TryDefineMixin("All nodes that should be valid in constant lookups\nshould implement this interface.", false, value.ToSymbol("ComplexConstantNode"), env)
				namespace.TryDefineClass("Represents a constant with as in using declarations\neg. `Foo::Bar as Bar`.", false, true, true, false, value.ToSymbol("ConstantAsNode"), objectClass, env)
				namespace.TryDefineClass("Represents a constant declaration eg. `const Foo: ArrayList[String] = [\"foo\", \"bar\"]`", false, true, true, false, value.ToSymbol("ConstantDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents a constant lookup expressions eg. `Foo::Bar`", false, true, true, false, value.ToSymbol("ConstantLookupNode"), objectClass, env)
				namespace.TryDefineMixin("All nodes that should be valid constants\nshould implement this interface.", false, value.ToSymbol("ConstantNode"), env)
				namespace.TryDefineClass("Represents a constructor call eg. `String(123)`", false, true, true, false, value.ToSymbol("ConstructorCallNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `continue` expression eg. `continue`, `continue \"foo\"`", false, true, true, false, value.ToSymbol("ContinueExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `do` expression eg.\n\n```\ndo\n  print(\"awesome!\")\nend\n```", false, true, true, false, value.ToSymbol("DoExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a simple double quoted string literal eg. `\"foo baz\"`", false, true, true, false, value.ToSymbol("DoubleQuotedStringLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a double splat expression eg. `**foo`", false, true, true, false, value.ToSymbol("DoubleSplatExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents an empty statement eg. a statement with only a semicolon or a newline.", false, true, true, false, value.ToSymbol("EmptyStatementNode"), objectClass, env)
				{
					namespace := namespace.TryDefineMixin("A mixin included in all Elk AST nodes\nthat can be treated as expressions.", false, value.ToSymbol("ExpressionNode"), env)
					namespace.TryDefineInterface("Represents a value that can be converted to an Elk AST Expression Node.", value.ToSymbol("Convertible"), env)
					namespace.Name() // noop - avoid unused variable error
				}
				namespace.TryDefineClass("Expression optionally terminated with a newline or a semicolon.", false, true, true, false, value.ToSymbol("ExpressionStatementNode"), objectClass, env)
				namespace.TryDefineClass("Represents an `extend where` block expression eg.\n\n```\nextend where T < Foo\n  def hello then println(\"awesome!\")\nend\n```", false, true, true, false, value.ToSymbol("ExtendWhereBlockExpressionNode"), objectClass, env)
				namespace.TryDefineClass("`false` literal.", false, true, true, false, value.ToSymbol("FalseLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Float32 literal eg. `5.2f32`, `.5f32`, `45e20f32`", false, true, true, false, value.ToSymbol("Float32LiteralNode"), objectClass, env)
				namespace.TryDefineClass("Float64 literal eg. `5.2f64`, `.5f64`, `45e20f64`", false, true, true, false, value.ToSymbol("Float64LiteralNode"), objectClass, env)
				namespace.TryDefineClass("Float literal eg. `5.2`, `.5`, `45e20`", false, true, true, false, value.ToSymbol("FloatLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `for in` expression eg. `for i in 5..15 then println(i)`", false, true, true, false, value.ToSymbol("ForInExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a formal parameter in function or struct declarations eg. `foo: String = 'bar'`", false, true, true, false, value.ToSymbol("FormalParameterNode"), objectClass, env)
				namespace.TryDefineClass("Represents a generic constant in type annotations eg. `ArrayList[String]`", false, true, true, false, value.ToSymbol("GenericConstantNode"), objectClass, env)
				namespace.TryDefineClass("Represents a constructor call eg. `ArrayList::[Int](1, 2, 3)`", false, true, true, false, value.ToSymbol("GenericConstructorCallNode"), objectClass, env)
				namespace.TryDefineClass("Represents a method call eg. `foo.bar::[String](a)`", false, true, true, false, value.ToSymbol("GenericMethodCallNode"), objectClass, env)
				namespace.TryDefineClass("Represents a generic function-like call eg. `foo::[Int](123)`", false, true, true, false, value.ToSymbol("GenericReceiverlessMethodCallNode"), objectClass, env)
				namespace.TryDefineClass("Represents a new generic type definition eg. `typedef Nilable[T] = T | nil`", false, true, true, false, value.ToSymbol("GenericTypeDefinitionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a new getter declaration eg. `getter foo: String`", false, true, true, false, value.ToSymbol("GetterDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `go` expression eg. `go foo()`, `go; foo(); end`", false, true, true, false, value.ToSymbol("GoExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a HashMap literal eg. `{ foo: 1, 'bar' => 5, baz }`", false, true, true, false, value.ToSymbol("HashMapLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a Record literal eg. `%{ foo: 1, 'bar' => 5, baz }`", false, true, true, false, value.ToSymbol("HashRecordLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a HashSet literal eg. `^[1, 5, -6]`", false, true, true, false, value.ToSymbol("HashSetLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a hex ArrayList literal eg. `\\x[ff ee]`", false, true, true, false, value.ToSymbol("HexArrayListLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a hex ArrayTuple literal eg. `%x[ff ee]`", false, true, true, false, value.ToSymbol("HexArrayTupleLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a hex HashSet literal eg. `^x[ff ee]`", false, true, true, false, value.ToSymbol("HexHashSetLiteralNode"), objectClass, env)
				namespace.TryDefineMixin("All nodes that should be valid identifiers\nshould implement this interface.", false, value.ToSymbol("IdentifierNode"), env)
				namespace.TryDefineClass("Represents an `if` expression eg. `if foo then println(\"bar\")`", false, true, true, false, value.ToSymbol("IfExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents an enhance expression eg. `implement Enumerable[V]`", false, true, true, false, value.ToSymbol("ImplementExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents an import statement eg. `import \"./foo/bar.elk\"`", false, true, true, false, value.ToSymbol("ImportStatementNode"), objectClass, env)
				namespace.TryDefineClass("Represents an include expression eg. `include Enumerable[V]`", false, true, true, false, value.ToSymbol("IncludeExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a constructor definition eg. `init then 'hello world'`", false, true, true, false, value.ToSymbol("InitDefinitionNode"), objectClass, env)
				namespace.TryDefineClass("Represents an instance type eg. `%self`", false, true, true, false, value.ToSymbol("InstanceOfTypeNode"), objectClass, env)
				namespace.TryDefineClass("Represents an instance variable declaration eg. `var @foo: String`", false, true, true, false, value.ToSymbol("InstanceVariableDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents an instance variable eg. `@foo`", false, true, true, false, value.ToSymbol("InstanceVariableNode"), objectClass, env)
				namespace.TryDefineClass("Int16 literal eg. `5i16`, `1_20i16`, `0xffi16`", false, true, true, false, value.ToSymbol("Int16LiteralNode"), objectClass, env)
				namespace.TryDefineClass("Int32 literal eg. `5i32`, `1_20i32`, `0xffi32`", false, true, true, false, value.ToSymbol("Int32LiteralNode"), objectClass, env)
				namespace.TryDefineClass("Int64 literal eg. `5i64`, `125_355i64`, `0xffi64`", false, true, true, false, value.ToSymbol("Int64LiteralNode"), objectClass, env)
				namespace.TryDefineClass("Int8 literal eg. `5i8`, `1_20i8`, `0xffi8`", false, true, true, false, value.ToSymbol("Int8LiteralNode"), objectClass, env)
				namespace.TryDefineMixin("All nodes that should be able to appear as\nelements of Int collection literals should\nimplement this interface.", false, value.ToSymbol("IntCollectionContentNode"), env)
				namespace.TryDefineClass("Int literal eg. `5`, `125_355`, `0xff`", false, true, true, false, value.ToSymbol("IntLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents an interface declaration eg. `interface Foo; end`", false, true, true, false, value.ToSymbol("InterfaceDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents an Interpolated regex literal eg. `%/foo${1 + 2}bar/`", false, true, true, false, value.ToSymbol("InterpolatedRegexLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents an interpolated string literal eg. `\"foo ${bar} baz\"`", false, true, true, false, value.ToSymbol("InterpolatedStringLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents an interpolated symbol eg. `:\"foo ${bar + 2}\"`", false, true, true, false, value.ToSymbol("InterpolatedSymbolLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Union type eg. `String & Int & Float`", false, true, true, false, value.ToSymbol("IntersectionTypeNode"), objectClass, env)
				namespace.TryDefineClass("Represents a syntax error.", false, true, true, false, value.ToSymbol("InvalidNode"), objectClass, env)
				namespace.TryDefineClass("Represents a key value expression eg. `foo => bar`", false, true, true, false, value.ToSymbol("KeyValueExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a key value pattern eg. `foo => bar`", false, true, true, false, value.ToSymbol("KeyValuePatternNode"), objectClass, env)
				namespace.TryDefineClass("Represents a labeled expression eg. `$foo: 1 + 2`", false, true, true, false, value.ToSymbol("LabeledExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a List pattern eg. `[1, a, >= 10]`", false, true, true, false, value.ToSymbol("ListPatternNode"), objectClass, env)
				namespace.TryDefineClass("Expression of a logical operator with two operands eg. `foo && bar`", false, true, true, false, value.ToSymbol("LogicalExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `loop` expression.", false, true, true, false, value.ToSymbol("LoopExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Create a new macro boundary node.\nIt's a block of code with an isolated scope\neg.\n\n```\ndo macro 'Dupa'\n  print(\"awesome!\")\nend\n```", false, true, true, false, value.ToSymbol("MacroBoundaryNode"), objectClass, env)
				namespace.TryDefineClass("Represents a macro call eg. `'123'.foo!()`", false, true, true, false, value.ToSymbol("MacroCallNode"), objectClass, env)
				namespace.TryDefineClass("Represents a macro definition eg. `macro foo(a: Elk::AST::StringLiteralNode); a; end`", false, true, true, false, value.ToSymbol("MacroDefinitionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a Map pattern eg. `{ foo: 5, bar: a, 5 => >= 10 }`", false, true, true, false, value.ToSymbol("MapPatternNode"), objectClass, env)
				namespace.TryDefineClass("Represents a method call eg. `'123'.to_int()`", false, true, true, false, value.ToSymbol("MethodCallNode"), objectClass, env)
				namespace.TryDefineClass("Represents a method definition eg. `def foo: String then 'hello world'`", false, true, true, false, value.ToSymbol("MethodDefinitionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a method lookup with as in using declarations\neg. `Foo::bar as Bar`.", false, true, true, false, value.ToSymbol("MethodLookupAsNode"), objectClass, env)
				namespace.TryDefineClass("Represents a method lookup expression eg. `Foo::bar`, `a::c`", false, true, true, false, value.ToSymbol("MethodLookupNode"), objectClass, env)
				namespace.TryDefineClass("Represents a formal parameter in method declarations eg. `foo: String = 'bar'`", false, true, true, false, value.ToSymbol("MethodParameterNode"), objectClass, env)
				namespace.TryDefineClass("Represents a method signature definition eg. `sig to_string(val: Int): String`", false, true, true, false, value.ToSymbol("MethodSignatureDefinitionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a mixin declaration eg. `mixin Foo; end`", false, true, true, false, value.ToSymbol("MixinDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents an `for .. in` modifier expression eg. `println(i) for i in 10..30`", false, true, true, false, value.ToSymbol("ModifierForInNode"), objectClass, env)
				namespace.TryDefineClass("Represents an `if .. else` modifier expression eg. `foo = 1 if bar else foo = 2`", false, true, true, false, value.ToSymbol("ModifierIfElseNode"), objectClass, env)
				namespace.TryDefineClass("Represents an `if`, `unless`, `while` or `until` modifier expression eg. `return true if foo`.", false, true, true, false, value.ToSymbol("ModifierNode"), objectClass, env)
				namespace.TryDefineClass("Represents a module declaration eg. `module Foo; end`", false, true, true, false, value.ToSymbol("ModuleDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `must` expression eg. `must foo()`", false, true, true, false, value.ToSymbol("MustExpressionNode"), objectClass, env)
				namespace.TryDefineMixin("Nodes that implement this interface represent\nnamed arguments in method calls.", false, value.ToSymbol("NamedArgumentNode"), env)
				namespace.TryDefineClass("Represents a named argument in a function call eg. `foo: 123`", false, true, true, false, value.ToSymbol("NamedCallArgumentNode"), objectClass, env)
				namespace.TryDefineClass("`never` type.", false, true, true, false, value.ToSymbol("NeverTypeNode"), objectClass, env)
				namespace.TryDefineClass("Represents a new expression eg. `new(123)`", false, true, true, false, value.ToSymbol("NewExpressionNode"), objectClass, env)
				namespace.TryDefineClass("`nil` literal.", false, true, true, false, value.ToSymbol("NilLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents nil-safe subscript access eg. `arr?[5]`", false, true, true, false, value.ToSymbol("NilSafeSubscriptExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents an optional or nilable type eg. `String?`", false, true, true, false, value.ToSymbol("NilableTypeNode"), objectClass, env)
				{
					namespace := namespace.TryDefineMixin("A base for all Elk AST (Abstract Syntax Tree) nodes.", false, value.ToSymbol("Node"), env)
					namespace.TryDefineInterface("Represents a value that can be converted to an Elk AST Node.", value.ToSymbol("Convertible"), env)
					namespace.Name() // noop - avoid unused variable error
				}
				namespace.TryDefineClass("Represents a not type eg. `~String`", false, true, true, false, value.ToSymbol("NotTypeNode"), objectClass, env)
				namespace.TryDefineClass("Represents a numeric `for` expression eg. `fornum i := 0; i < 10; i += 1 then println(i)`", false, true, true, false, value.ToSymbol("NumericForExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents an Object pattern eg. `Foo(foo: 5, bar: a, c)`", false, true, true, false, value.ToSymbol("ObjectPatternNode"), objectClass, env)
				namespace.TryDefineMixin("All nodes that should be valid in parameter declaration lists\nof methods or functions should implement this interface.", true, value.ToSymbol("ParameterNode"), env)
				namespace.TryDefineClass("Formal parameter optionally terminated with a newline or a semicolon.", false, true, true, false, value.ToSymbol("ParameterStatementNode"), objectClass, env)
				namespace.TryDefineMixin("Represents AST nodes that are valid expressions and patterns.", false, value.ToSymbol("PatternExpressionNode"), env)
				{
					namespace := namespace.TryDefineMixin("All nodes that should be valid in pattern matching should\nimplement this interface", false, value.ToSymbol("PatternNode"), env)
					namespace.TryDefineInterface("Represents a value that can be converted to an Elk AST Pattern Node.", value.ToSymbol("Convertible"), env)
					namespace.Name() // noop - avoid unused variable error
				}
				namespace.TryDefineClass("Postfix expression eg. `foo++`, `bar--`", false, true, true, false, value.ToSymbol("PostfixExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a private constant eg. `_Foo`", false, true, true, false, value.ToSymbol("PrivateConstantNode"), objectClass, env)
				namespace.TryDefineClass("Represents a private identifier eg. `_foo`", false, true, true, false, value.ToSymbol("PrivateIdentifierNode"), objectClass, env)
				namespace.TryDefineClass("Represents a single Elk program (usually a single file).", false, true, true, false, value.ToSymbol("ProgramNode"), objectClass, env)
				namespace.TryDefineClass("Represents a constant with as in using declarations\neg. `Foo as Bar`.", false, true, true, false, value.ToSymbol("PublicConstantAsNode"), objectClass, env)
				namespace.TryDefineClass("Represents a public constant eg. `Foo`.", false, true, true, false, value.ToSymbol("PublicConstantNode"), objectClass, env)
				namespace.TryDefineClass("Represents an identifier with as in using declarations\neg. `foo as bar`.", false, true, true, false, value.ToSymbol("PublicIdentifierAsNode"), objectClass, env)
				namespace.TryDefineClass("Represents a public identifier eg. `foo`.", false, true, true, false, value.ToSymbol("PublicIdentifierNode"), objectClass, env)
				namespace.TryDefineClass("Represents a quoted block of AST eg.\n\n```\nquote\n  print(\"awesome!\")\nend\n```", false, true, true, false, value.ToSymbol("QuoteExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a Range literal eg. `1...5`", false, true, true, false, value.ToSymbol("RangeLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Raw Char literal eg. `a`", false, true, true, false, value.ToSymbol("RawCharLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Raw string literal enclosed with single quotes eg. `'foo'`.", false, true, true, false, value.ToSymbol("RawStringLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a function-like macro call eg. `foo!(123)`", false, true, true, false, value.ToSymbol("ReceiverlessMacroCallNode"), objectClass, env)
				namespace.TryDefineClass("Represents a function-like call eg. `to_string(123)`", false, true, true, false, value.ToSymbol("ReceiverlessMethodCallNode"), objectClass, env)
				namespace.TryDefineClass("Represents a Record pattern eg. `%{ foo: 5, bar: a, 5 => >= 10 }`", false, true, true, false, value.ToSymbol("RecordPatternNode"), objectClass, env)
				namespace.TryDefineClass("Represents a single interpolated section of a regex literal eg. `bar + 2` in `%/foo${bar + 2}/`", false, true, true, false, value.ToSymbol("RegexInterpolationNode"), objectClass, env)
				namespace.TryDefineMixin("Nodes that implement this interface can appear\ninside of a Regex literal.", false, value.ToSymbol("RegexLiteralContentNode"), env)
				namespace.TryDefineClass("Represents a single section of characters of a regex literal eg. `foo` in `%/foo${bar}/`.", false, true, true, false, value.ToSymbol("RegexLiteralContentSectionNode"), objectClass, env)
				namespace.TryDefineMixin("All nodes that represent regexes should\nimplement this interface.", false, value.ToSymbol("RegexLiteralNode"), env)
				namespace.TryDefineClass("Represents a rest element in a list pattern eg. `*a`", false, true, true, false, value.ToSymbol("RestPatternNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `return` expression eg. `return`, `return true`", false, true, true, false, value.ToSymbol("ReturnExpressionNode"), objectClass, env)
				namespace.TryDefineClass("`self` literal.", false, true, true, false, value.ToSymbol("SelfLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a Set pattern eg. `^[1, \"foo\"]`", false, true, true, false, value.ToSymbol("SetPatternNode"), objectClass, env)
				namespace.TryDefineClass("Represents a new setter declaration eg. `setter foo: String`", false, true, true, false, value.ToSymbol("SetterDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents a signature parameter in method and function signatures eg. `foo?: String`", false, true, true, false, value.ToSymbol("SignatureParameterNode"), objectClass, env)
				namespace.TryDefineMixin("All nodes that represent simple strings (without interpolation)\nshould implement this interface.", false, value.ToSymbol("SimpleStringLiteralNode"), env)
				namespace.TryDefineClass("Represents a symbol literal with simple content eg. `:foo`, `:'foo bar`, `:\"lol\"`", false, true, true, false, value.ToSymbol("SimpleSymbolLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `singleton` block expression eg.\n\nsingleton\n  def hello then println(\"awesome!\")\nend", false, true, true, false, value.ToSymbol("SingletonBlockExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a singleton type eg. `&String`", false, true, true, false, value.ToSymbol("SingletonTypeNode"), objectClass, env)
				namespace.TryDefineClass("Represents a splat expression eg. `*foo`", false, true, true, false, value.ToSymbol("SplatExpressionNode"), objectClass, env)
				namespace.TryDefineMixin("Represents a single statement, so for example\na single valid \"line\" of Elk code.\nUsually its an expression optionally terminated with a newline or a semicolon.", false, value.ToSymbol("StatementNode"), env)
				namespace.TryDefineClass("Represents a single inspect interpolated section of a string literal eg. `bar + 2` in `\"foo#{bar + 2}\"`", false, true, true, false, value.ToSymbol("StringInspectInterpolationNode"), objectClass, env)
				namespace.TryDefineClass("Represents a single interpolated section of a string literal eg. `bar + 2` in `\"foo${bar + 2}\"`", false, true, true, false, value.ToSymbol("StringInterpolationNode"), objectClass, env)
				namespace.TryDefineMixin("Nodes that implement this interface can appear\ninside of a String literal.", false, value.ToSymbol("StringLiteralContentNode"), env)
				namespace.TryDefineClass("Represents a single section of characters of a string literal eg. `foo` in `\"foo${bar}\"`.", false, true, true, false, value.ToSymbol("StringLiteralContentSectionNode"), objectClass, env)
				namespace.TryDefineMixin("All nodes that represent strings should\nimplement this interface.", false, value.ToSymbol("StringLiteralNode"), env)
				namespace.TryDefineMixin("Represents a single statement of a struct body\noptionally terminated with a newline or semicolon.", false, value.ToSymbol("StructBodyStatementNode"), env)
				namespace.TryDefineClass("Represents a struct declaration eg. `struct Foo; end`", false, true, true, false, value.ToSymbol("StructDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents subscript access eg. `arr[5]`", false, true, true, false, value.ToSymbol("SubscriptExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `switch` expression eg.\n\nswitch a\ncase 3\n  println(\"eureka!\")\ncase nil\n  println(\"boo\")\nelse\n  println(\"nothing\")\nend", false, true, true, false, value.ToSymbol("SwitchExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a symbol ArrayList literal eg. `\\s[foo bar]`", false, true, true, false, value.ToSymbol("SymbolArrayListLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a symbol ArrayTuple literal eg. `%s[foo bar]`", false, true, true, false, value.ToSymbol("SymbolArrayTupleLiteralNode"), objectClass, env)
				namespace.TryDefineMixin("All nodes that should be able to appear as\nelements of symbol collection literals should\nimplement this interface.", false, value.ToSymbol("SymbolCollectionContentNode"), env)
				namespace.TryDefineClass("Represents a symbol HashSet literal eg. `^s[foo bar]`", false, true, true, false, value.ToSymbol("SymbolHashSetLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a symbol value expression eg. `foo: bar`", false, true, true, false, value.ToSymbol("SymbolKeyValueExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a symbol value pattern eg. `foo: bar`", false, true, true, false, value.ToSymbol("SymbolKeyValuePatternNode"), objectClass, env)
				namespace.TryDefineMixin("A mixin included in all Elk AST nodes\nthat can be treated as symbol literals.", false, value.ToSymbol("SymbolLiteralNode"), env)
				namespace.TryDefineClass("Represents a `throw` expression eg. `throw ArgumentError(\"foo\")`", false, true, true, false, value.ToSymbol("ThrowExpressionNode"), objectClass, env)
				namespace.TryDefineClass("`true` literal.", false, true, true, false, value.ToSymbol("TrueLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `try` expression eg. `try foo()`", false, true, true, false, value.ToSymbol("TryExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a Tuple pattern eg. `%[1, a, >= 10]`", false, true, true, false, value.ToSymbol("TuplePatternNode"), objectClass, env)
				namespace.TryDefineClass("Represents a new type definition eg. `typedef StringList = ArrayList[String]`", false, true, true, false, value.ToSymbol("TypeDefinitionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a type expression `type String?`", false, true, true, false, value.ToSymbol("TypeExpressionNode"), objectClass, env)
				{
					namespace := namespace.TryDefineMixin("All nodes that should be valid in type annotations should\nimplement this interface", false, value.ToSymbol("TypeNode"), env)
					namespace.TryDefineInterface("Represents a value that can be converted to an Elk AST Type Node.", value.ToSymbol("Convertible"), env)
					namespace.Name() // noop - avoid unused variable error
				}
				namespace.TryDefineMixin("Represents a type variable in generics like `class Foo[+V]; end`", false, value.ToSymbol("TypeParameterNode"), env)
				namespace.TryDefineClass("Represents a `typeof` expression eg. `typeof foo()`", false, true, true, false, value.ToSymbol("TypeofExpressionNode"), objectClass, env)
				namespace.TryDefineClass("UInt16 literal eg. `5u16`, `1_20u16`, `0xffu16`", false, true, true, false, value.ToSymbol("UInt16LiteralNode"), objectClass, env)
				namespace.TryDefineClass("UInt32 literal eg. `5u32`, `1_20u32`, `0xffu32`", false, true, true, false, value.ToSymbol("UInt32LiteralNode"), objectClass, env)
				namespace.TryDefineClass("UInt64 literal eg. `5u64`, `125_355u64`, `0xffu64`", false, true, true, false, value.ToSymbol("UInt64LiteralNode"), objectClass, env)
				namespace.TryDefineClass("UInt8 literal eg. `5u8`, `1_20u8`, `0xffu8`", false, true, true, false, value.ToSymbol("UInt8LiteralNode"), objectClass, env)
				namespace.TryDefineClass("Expression of an operator with one operand eg. `!foo`, `-bar`", false, true, true, false, value.ToSymbol("UnaryExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Type of an operator with one operand eg. `-2`, `+3`", false, true, true, false, value.ToSymbol("UnaryTypeNode"), objectClass, env)
				namespace.TryDefineClass("`undefined` literal.", false, true, true, false, value.ToSymbol("UndefinedLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents an uninterpolated regex literal eg. `%/foo/`", false, true, true, false, value.ToSymbol("UninterpolatedRegexLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Union type eg. `String | Int | Float`", false, true, true, false, value.ToSymbol("UnionTypeNode"), objectClass, env)
				namespace.TryDefineClass("Represents an `unless` expression eg. `unless foo then println(\"bar\")`", false, true, true, false, value.ToSymbol("UnlessExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents an unquoted block of AST in a quote eg.\n\n```\nunquote(x)\n```", false, true, true, false, value.ToSymbol("UnquoteExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `until` expression eg. `until i >= 5 then i += 5`", false, true, true, false, value.ToSymbol("UntilExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a using all entry node eg. `Foo::*`, `A::B::C::*`", false, true, true, false, value.ToSymbol("UsingAllEntryNode"), objectClass, env)
				namespace.TryDefineMixin("Represents all nodes that are valid in using declarations", false, value.ToSymbol("UsingEntryNode"), env)
				namespace.TryDefineClass("Represents a using entry node with subentries eg. `Foo::{Bar, baz}`, `A::B::C::{lol, foo as epic, Gro as Moe}`", false, true, true, false, value.ToSymbol("UsingEntryWithSubentriesNode"), objectClass, env)
				namespace.TryDefineClass("Represents a using expression eg. `using Foo`", false, true, true, false, value.ToSymbol("UsingExpressionNode"), objectClass, env)
				namespace.TryDefineMixin("Represents all nodes that are valid in using subentries\nin `UsingEntryWithSubentriesNode`", false, value.ToSymbol("UsingSubentryNode"), env)
				namespace.TryDefineClass("Represents a value declaration eg. `val foo: String`", false, true, true, false, value.ToSymbol("ValueDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents a value pattern declaration eg. `val [foo, { bar }] = baz()`", false, true, true, false, value.ToSymbol("ValuePatternDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents a variable declaration eg. `var foo: String`", false, true, true, false, value.ToSymbol("VariableDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents a variable declaration with patterns eg. `var [foo, { bar }] = baz()`", false, true, true, false, value.ToSymbol("VariablePatternDeclarationNode"), objectClass, env)
				namespace.TryDefineClass("Represents a type parameter eg. `+V`", false, true, true, false, value.ToSymbol("VariantTypeParameterNode"), objectClass, env)
				namespace.TryDefineClass("`void` type.", false, true, true, false, value.ToSymbol("VoidTypeNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `while` expression eg. `while i < 5 then i += 5`", false, true, true, false, value.ToSymbol("WhileExpressionNode"), objectClass, env)
				namespace.TryDefineClass("Represents a word ArrayList literal eg. `\\w[foo bar]`", false, true, true, false, value.ToSymbol("WordArrayListLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a word ArrayTuple literal eg. `%w[foo bar]`", false, true, true, false, value.ToSymbol("WordArrayTupleLiteralNode"), objectClass, env)
				namespace.TryDefineMixin("All nodes that should be able to appear as\nelements of word collection literals should\nimplement this interface.", false, value.ToSymbol("WordCollectionContentNode"), env)
				namespace.TryDefineClass("Represents a word HashSet literal eg. `^w[foo bar]`", false, true, true, false, value.ToSymbol("WordHashSetLiteralNode"), objectClass, env)
				namespace.TryDefineClass("Represents a `yield` expression eg. `yield`, `yield true`, `yield* foo()`", false, true, true, false, value.ToSymbol("YieldExpressionNode"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.TryDefineClass("Implements a lexer/tokenizer of the Elk programming language.", false, true, true, false, value.ToSymbol("Lexer"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("Implements a parser of the Elk programming language.", false, true, true, true, value.ToSymbol("Parser"), objectClass, env)
				namespace.TryDefineClass("Represents the result of parsing Elk source code.", false, true, true, true, value.ToSymbol("Result"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.TryDefineClass("Represents a token produced by the Elk lexer.\n\nA token is a single lexical unit of text\nwith a particular meaning.", false, true, true, false, value.ToSymbol("Token"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents a closed range from a given value to +∞ *[start, +∞)*", false, true, true, false, value.ToSymbol("EndlessClosedRange"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents an open range from a given value to +∞ *(start, +∞)*", false, true, true, false, value.ToSymbol("EndlessOpenRange"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("A base class for most errors in Elk stdlib.", false, false, false, false, value.ToSymbol("Error"), objectClass, env)
		{
			namespace := namespace.TryDefineModule("Contains types used for interacting\n with the file system.", value.ToSymbol("FS"), env)
			namespace.TryDefineClass("Represents the position of a piece of text in a file.\n\nIt is made up of a path and a span.", false, true, true, false, value.ToSymbol("Location"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("Represents a file system path.\nPaths are immutable and safe to use by multiple threads.", false, true, true, false, value.ToSymbol("Path"), objectClass, env)
				namespace.TryDefineClass("", false, false, false, false, value.ToSymbol("Error"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("", false, true, true, true, value.ToSymbol("False"), objectClass, env)
		namespace.DefineSubtype(value.ToSymbol("Falsy"), NewNamedType("Std::Falsy", NewUnion(Nil{}, False{})))
		namespace.TryDefineClass("", false, false, false, false, value.ToSymbol("FileSystemError"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("Represents a floating point number (a fraction like `1.2`, `0.1`).\n\nThis float type has 64 bits on 64 bit platforms\nand 32 bit on 32 bit platforms.", false, true, true, true, value.ToSymbol("Float"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a float.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents a floating point number (a fraction like `1.2`, `0.1`).\n\nThis float type has 64 bits.", false, true, true, true, value.ToSymbol("Float32"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a float32.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents a floating point number (a fraction like `1.2`, `0.1`).\n\nThis float type has 64 bits.", false, true, true, true, value.ToSymbol("Float64"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a float64.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Thrown when a literal or interpreted string has an incorrect format.", false, false, false, false, value.ToSymbol("FormatError"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("Implements a generator object that is iterable.", false, true, true, true, value.ToSymbol("Generator"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Thrown when an invalid glob pattern has been encountered.", false, false, false, false, value.ToSymbol("GlobError"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("A dynamically resizable map data structure backed\nby an array with a hashing algorithm.\n\nIt is an unordered collection of key-value pairs.", false, true, true, false, value.ToSymbol("HashMap"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("A record data structure backed by an array with a hashing algorithm.\n\nIt is an unordered immutable collection of key-value pairs.\nA record is an immutable map.", false, true, true, false, value.ToSymbol("HashRecord"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("A dynamically resizable set data structure backed\nby an array with a hashing algorithm.\n\nIt is an unordered collection of unique values.", false, true, true, false, value.ToSymbol("HashSet"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineInterface("Represents a value that can compute its own hash for use in\ndata structures like hashmaps, hashrecords, hashsets.", value.ToSymbol("Hashable"), env)
		{
			namespace := namespace.TryDefineInterface("An interface that represents a finite, immutable collection\nof elements.", value.ToSymbol("ImmutableCollection"), env)
			{
				namespace := namespace.TryDefineMixin("Provides default implementations of most immutable collection methods.", true, value.ToSymbol("Base"), env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineMixin("Represents an unordered, immutable collection\nof unique elements.", true, value.ToSymbol("ImmutableSet"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineInterface("Represents a value that can be incremented using\nthe `++` operator like `a++`", value.ToSymbol("Incrementable"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a human readable string\nthat represents the structure of the value.", value.ToSymbol("Inspectable"), env)
		{
			namespace := namespace.TryDefineClass("Represents an integer (a whole number like `1`, `2`, `3`, `-5`, `0`).\n\nThis integer type is automatically resized so\nit can hold an arbitrarily large/small number.", false, true, true, true, value.ToSymbol("Int"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to an integer.", value.ToSymbol("Convertible"), env)
			namespace.TryDefineClass("", false, false, false, false, value.ToSymbol("Iterator"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents a signed 16 bit integer (a whole number like `1i16`, `2i16`, `-3i16`, `0i16`).", false, true, true, true, value.ToSymbol("Int16"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a int16.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents a signed 32 bit integer (a whole number like `1i32`, `2i32`, `-3i32`, `0i32`).", false, true, true, true, value.ToSymbol("Int32"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a int32.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents a signed 64 bit integer (a whole number like `1i64`, `2i64`, `-3i64`, `0i64`).", false, true, true, true, value.ToSymbol("Int64"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a int64.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents a signed 8 bit integer (a whole number like `1i8`, `2i8`, `-3i8`, `0i8`).", false, true, true, true, value.ToSymbol("Int8"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a int8.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("`Interface` is the class of all interfaces.", false, false, false, true, value.ToSymbol("Interface"), objectClass, env)
		{
			namespace := namespace.TryDefineInterface("Represents a value that can be iterated over in a `for` loop and implement\nmany useful methods.", value.ToSymbol("Iterable"), env)
			{
				namespace := namespace.TryDefineMixin("Provides default implementations of most iterable methods.", true, value.ToSymbol("Base"), env)
				namespace.Name() // noop - avoid unused variable error
			}
			{
				namespace := namespace.TryDefineMixin("Provides default implementations of most iterable methods\nfor finite iterables.", true, value.ToSymbol("FiniteBase"), env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.TryDefineClass("", false, false, false, false, value.ToSymbol("NotFoundError"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineInterface("Represents a `Range` that can be iterated over.", value.ToSymbol("IterableRange"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineInterface("An interface that represents objects\nthat allow for external iteration.", value.ToSymbol("Iterator"), env)
			{
				namespace := namespace.TryDefineMixin("Provides default implementations of most iterator methods.", true, value.ToSymbol("Base"), env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineModule("Contains builtin global functions like `println` etc.", value.ToSymbol("Kernel"), env)
		{
			namespace := namespace.TryDefineClass("Represents a left-open range from `start` to `end` *(start, end]*", false, true, true, false, value.ToSymbol("LeftOpenRange"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineMixin("Represents an ordered, mutable collection\nof elements indexed by integers starting at `0`.", true, value.ToSymbol("List"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineInterface("Represents a resource that can be locked and unlocked.", value.ToSymbol("Lockable"), env)
		{
			namespace := namespace.TryDefineMixin("Represents an unordered mutable collection of key-value pairs.", true, value.ToSymbol("Map"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("", false, true, true, true, value.ToSymbol("Method"), objectClass, env)
		namespace.TryDefineClass("`Mixin` is the class of all mixins.", false, false, false, true, value.ToSymbol("Mixin"), objectClass, env)
		namespace.TryDefineClass("`Module` is the class of all modules.", false, false, false, true, value.ToSymbol("Module"), objectClass, env)
		namespace.TryDefineClass("Represents an empty value.", false, true, true, true, value.ToSymbol("Nil"), objectClass, env)
		namespace.TryDefineClass("", false, false, false, false, value.ToSymbol("Object"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("Represents an open range from `start` to `end` *(start, end)*", false, true, true, false, value.ToSymbol("OpenRange"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Thrown when a numeric value is too large or too small to be used in a particular setting.", false, false, false, false, value.ToSymbol("OutOfRangeError"), objectClass, env)
		{
			namespace := namespace.TryDefineClass("A `Pair` represents a 2-element tuple,\nor a key-value pair.", false, true, true, false, value.ToSymbol("Pair"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineInterface("Represents a value that can be iterated over in a `for` loop.", value.ToSymbol("PrimitiveIterable"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("A promise is the return type of a asynchronous function.\nIt is a placeholder for a value that will be available at some point\nin the future.", false, true, true, true, value.ToSymbol("Promise"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineMixin("Represents a range of values, an interval.\n\nThe default implementation of `Range` is `ClosedRange`.", true, value.ToSymbol("Range"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineMixin("Represents an unordered immutable collection of key-value pairs.\nA record is an immutable map.", true, value.ToSymbol("Record"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("A `Regex` represents regular expression that can be used\nto match a pattern against strings.", false, true, true, false, value.ToSymbol("Regex"), objectClass, env)
		{
			namespace := namespace.TryDefineInterface("An interface that represents iterators that can be reset.", value.ToSymbol("ResettableIterator"), env)
			{
				namespace := namespace.TryDefineMixin("Provides default implementations of most resettable iterator methods.", true, value.ToSymbol("Base"), env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents a right-open range from `start` to `end` *[start, end)*", false, true, true, false, value.ToSymbol("RightOpenRange"), objectClass, env)
			{
				namespace := namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineMixin("Represents an unordered, mutable collection\nof unique elements.", true, value.ToSymbol("Set"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents the state of the call stack at some point in time.", false, true, true, true, value.ToSymbol("StackTrace"), objectClass, env)
			namespace.TryDefineClass("", false, true, true, true, value.ToSymbol("Iterator"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("", false, true, true, true, value.ToSymbol("String"), objectClass, env)
			namespace.TryDefineClass("Iterates over all bytes of a `String`.", false, true, true, false, value.ToSymbol("ByteIterator"), objectClass, env)
			namespace.TryDefineClass("Iterates over all unicode code points of a `String`.", false, true, true, false, value.ToSymbol("CharIterator"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a string.", value.ToSymbol("Convertible"), env)
			namespace.TryDefineClass("Iterates over all grapheme clusters of a `String`.", false, true, true, false, value.ToSymbol("GraphemeIterator"), objectClass, env)
			namespace.TryDefineClass("Represents a position of a character in a string.", false, true, true, false, value.ToSymbol("Position"), objectClass, env)
			namespace.TryDefineClass("Represents the position of a piece of text in a string.\n\nIt is made up of two positions, start position\n(position of the first character) and end position (position of the last character).", false, true, true, false, value.ToSymbol("Span"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Represents an interned string.\n\nA symbol is an integer ID that is associated\nwith a particular name (string).\n\nA few symbols with the same name refer to the same ID.\n\nComparing symbols happens in constant time, so it's\nusually faster than comparing strings.", false, true, true, true, value.ToSymbol("Symbol"), objectClass, env)
		{
			namespace := namespace.TryDefineModule("`Sync` provides synchronisation utilities like mutexes.", value.ToSymbol("Sync"), env)
			{
				namespace := namespace.TryDefineClass("A thread safe `DiagnosticList`, synchronized with a Mutex.", false, true, true, false, value.ToSymbol("DiagnosticList"), objectClass, env)
				namespace.TryDefineClass("", false, true, true, false, value.ToSymbol("Iterator"), objectClass, env)
				namespace.Name() // noop - avoid unused variable error
			}
			namespace.TryDefineClass("A `Mutex` is a mutual exclusion lock.\nIt can be used to synchronise operations in multiple threads.", false, true, true, false, value.ToSymbol("Mutex"), objectClass, env)
			namespace.TryDefineClass("`Once` is a kind of lock ensuring that a piece of\ncode will be executed exactly one time.", false, true, true, false, value.ToSymbol("Once"), objectClass, env)
			namespace.TryDefineClass("Wraps a `RWMutex` and exposes its `read_lock` and `read_unlock`\nmethods as `lock` and `unlock` respectively.", false, true, true, false, value.ToSymbol("ROMutex"), objectClass, env)
			namespace.TryDefineClass("A `Mutex` is a mutual exclusion lock that allows many readers or a single writer\nto hold the lock.", false, true, true, false, value.ToSymbol("RWMutex"), objectClass, env)
			namespace.TryDefineClass("A `WaitGroup` waits for threads to finish.\n\nYou can use the `add` method to specify the amount of threads to wait for.\nAfterwards each thread should call `end` when finished\nThe `wait` method can be used to block until all threads have finished.", false, true, true, false, value.ToSymbol("WaitGroup"), objectClass, env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("Represents a single Elk thread of execution.", false, true, true, true, value.ToSymbol("Thread"), objectClass, env)
		namespace.TryDefineClass("A pool of thread workers with a task queue.", false, true, true, true, value.ToSymbol("ThreadPool"), objectClass, env)
		namespace.TryDefineClass("Represents a moment in time with nanosecond precision.", false, true, true, false, value.ToSymbol("Time"), objectClass, env)
		namespace.TryDefineClass("Represents a timezone from the IANA Timezone database.", false, true, true, false, value.ToSymbol("Timezone"), objectClass, env)
		namespace.TryDefineClass("", false, true, true, true, value.ToSymbol("True"), objectClass, env)
		namespace.DefineSubtype(value.ToSymbol("Truthy"), NewNamedType("Std::Truthy", NewNot(NewNamedType("Std::Falsy", NewUnion(Nil{}, False{})))))
		{
			namespace := namespace.TryDefineMixin("Represents an ordered, immutable collection\nof elements indexed by integers starting at `0`.", true, value.ToSymbol("Tuple"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents an unsigned 16 bit integer (a positive whole number like `1u16`, `2u16`, `3u16`, `0u16`).", false, true, true, true, value.ToSymbol("UInt16"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a uint16.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents an unsigned 32 bit integer (a positive whole number like `1u32`, `2u32`, `3u32`, `0u32`).", false, true, true, true, value.ToSymbol("UInt32"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a uint32.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents an unsigned 64 bit integer (a positive whole number like `1u64`, `2u64`, `3u64`, `0u64`).", false, true, true, true, value.ToSymbol("UInt64"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a int8.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		{
			namespace := namespace.TryDefineClass("Represents an unsigned 8 bit integer (a positive whole number like `1u8`, `2u8`, `3u8`, `0u8`).", false, true, true, true, value.ToSymbol("UInt8"), objectClass, env)
			namespace.TryDefineInterface("Values that conform to this interface\ncan be converted to a uint8.", value.ToSymbol("Convertible"), env)
			namespace.Name() // noop - avoid unused variable error
		}
		namespace.TryDefineClass("`Value` is the superclass class of all\nElk classes.", false, false, true, false, value.ToSymbol("Value"), nil, env)
		namespace.Name() // noop - avoid unused variable error
	}

	// Define methods, constants

	{
		namespace := env.Root

		namespace.Name() // noop - avoid unused variable error

		// Include mixins and implement interfaces

		// Define methods

		// Define constants

		// Define instance variables

		{
			namespace := namespace.MustSubtype("Std").(*Module)

			namespace.Name() // noop - avoid unused variable error

			// Include mixins and implement interfaces

			// Define methods

			// Define constants

			// Define instance variables

			{
				namespace := namespace.MustSubtype("ArrayList").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::List", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ArrayList::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

				// Define methods
				namespace.DefineMethod("Create a new `ArrayList` containing the elements of `self`\nrepeated `n` times.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Create a new `ArrayList` containing the elements of `self`\nand another given `Tuple`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::ArrayList::Val", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Adds the given value to the list.\n\nReallocates the underlying array if it is\ntoo small to hold it.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::ArrayList::Val", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Check whether the given value is an `ArrayList`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Check whether the given value is an `ArrayList` or `ArrayTuple`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Get the element under the given index.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::ArrayList::Val", env), Never{})
				namespace.DefineMethod("Set the element under the given index to the given value.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]="), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::ArrayList::Val", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Adds the given values to the list.\n\nReallocates the underlying array if it is\ntoo small to hold them.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("append"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::ArrayList::Val", env), PositionalRestParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Returns the number of elements that can be\nheld by the underlying array.\n\nThis value will change when the list gets resized,\nand the underlying array gets reallocated", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Removes all elements from the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("clear"), nil, nil, Void{}, Never{})
				namespace.DefineMethod("Check whether the given `value` is present in this list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::ArrayList::Val", env), NameToType("Std::ArrayList::Val", env), NameToType("Std::ArrayList::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::ArrayList::Val", env), NameToType("Std::ArrayList::Val", env), NameToType("Std::ArrayList::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Mutates the list.\n\nReallocates the underlying array to hold\nthe given number of new elements.\n\nExpands the `capacity` of the list by `new_slots`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("grow"), nil, []*Parameter{NewParameter(value.ToSymbol("new_slots"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::ArrayList::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ArrayList::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Returns the number of left slots for new elements\nin the underlying array.\nIt tells you how many more elements can be\nadded to the list before the underlying array gets\nreallocated.\n\nIt is always equal to `capacity - length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left_capacity"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of elements present in the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over the elements of this list,\nyielding them to the given closure.\n\nReturns a new List that consists of the elements returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::ArrayList::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the elements of this list,\nyielding them to the given closure.\n\nMutates the list in place replacing the elements with the ones\nreturned by the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_mut"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::ArrayList::Val", env), NormalParameterKind, false)}, NameToType("Std::ArrayList::Val", env), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Adds the given value to the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("push"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::ArrayList::Val", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Removes the element from the list.\n\nReturns `true` if the element has been removed,\notherwise returns `false.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("remove"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::ArrayList::Val", env), NormalParameterKind, false)}, Bool{}, Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ArrayList::Iterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("list"), NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ArrayList::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::ArrayList::Iterator::Val", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("ArrayTuple").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::ArrayTuple::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))

				// Define methods
				namespace.DefineMethod("Create a new `ArrayTuple` containing the elements of `self`\nrepeated `n` times.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ArrayTuple::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Create a new `ArrayTuple` containing the elements of `self`\nand another given `Tuple`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::ArrayTuple::Val", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Check whether the given value is an `ArrayTuple`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Check whether the given value is an `ArrayTuple` or `ArrayList`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Get the element under the given index.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::ArrayTuple::Val", env), Never{})
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the tuple.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ArrayTuple::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Returns the number of elements present in the tuple.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over the elements of this tuple,\nyielding them to the given closure.\n\nReturns a new Tuple that consists of the elements returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::ArrayTuple::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ArrayTuple::Iterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("tuple"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ArrayTuple::Iterator::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the tuple.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::ArrayTuple::Iterator::Val", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("BeginlessClosedRange").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
				typeParam.UpperBound = NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::BeginlessClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Range", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::BeginlessClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))

				// Define methods
				namespace.DefineMethod("Check whether the given `value` is present in this range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::BeginlessClosedRange::Val", env), NameToType("Std::BeginlessClosedRange::Val", env), NameToType("Std::BeginlessClosedRange::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::BeginlessClosedRange::Val", env), NameToType("Std::BeginlessClosedRange::Val", env), NameToType("Std::BeginlessClosedRange::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Returns the upper bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("end"), nil, nil, NameToType("Std::BeginlessClosedRange::Val", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_closed"), nil, nil, False{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_closed"), nil, nil, True{}, Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("BeginlessOpenRange").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
				typeParam.UpperBound = NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::BeginlessOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Range", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::BeginlessOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))

				// Define methods
				namespace.DefineMethod("Check whether the given `value` is present in this range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::BeginlessOpenRange::Val", env), NameToType("Std::BeginlessOpenRange::Val", env), NameToType("Std::BeginlessOpenRange::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::BeginlessOpenRange::Val", env), NameToType("Std::BeginlessOpenRange::Val", env), NameToType("Std::BeginlessOpenRange::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Returns the upper bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("end"), nil, nil, NameToType("Std::BeginlessOpenRange::Val", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_closed"), nil, nil, False{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_closed"), nil, nil, False{}, Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("BigFloat").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10bf\n\tvar b = 3bf\n\ta % b #=> 1bf\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Multiply this float by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Exponentiate this float, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Add `other` to this bigfloat.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1.2bf\n\t+a #=> 1.2bf\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Subtract `other` from this bigfloat.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Returns the result of negating the number.\n\n```\n\tvar a = 1.2bf\n\t-a #=> -1.2bf\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Divide this float by another float.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this float with another float.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Calculates a hash of the float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Sets the precision to the given integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("p"), nil, []*Parameter{NewParameter(value.ToSymbol("precision"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("returns the mantissa precision of `self` in bits.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("precision"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Sets the precision to the given integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("set_precision"), nil, []*Parameter{NewParameter(value.ToSymbol("precision"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::BigFloatLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::BigFloatLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::BigFloatLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::BigFloatLiteralNode", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_bigfloat"), nil, nil, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Converts to a fixed-precision floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to an automatically resized integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the bigfloat to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Bool").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces

				// Define methods

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Box").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("v"), NameToType("Std::Box::Val", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Retrieves the value stored in the `Box`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("get"), nil, nil, NameToType("Std::Box::Val", env), Never{})
				namespace.DefineMethod("Stores the given value in the `Box`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("set"), nil, []*Parameter{NewParameter(value.ToSymbol("v"), NameToType("Std::Box::Val", env), NormalParameterKind, false)}, NameToType("Std::Box::Val", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("CallFrame").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Name of the source file where the called function\nis defined.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("file_name"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Name of the called function.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("func_name"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Number of the line in the source file\nwhere the definition of the called function starts.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("line_number"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Number of optimised tail calls before this call.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("tail_calls"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the string representation of the call frame.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Channel").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("V"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("V"), typeParam)
				namespace.DefineConstant(value.ToSymbol("V"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterable::FiniteBase", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Channel::V", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))
				ImplementInterface(namespace, NameToType("Std::Closable", env).(*Interface))

				// Define methods
				namespace.DefineMethod("Create a new `Channel` with the given capacity.\nDefault capacity is `0`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("capacity"), NameToType("Std::Int", env), DefaultValueParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Pushes a new element to the channel.\nBlocks the current thread until another thread pops the element if the channel does not have any empty slots in the buffer.\n\nPushing to a closed channel throws an unchecked error.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Channel::V", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Returns the size of the buffer that can hold elements\nuntil they're popped.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Closes the channel, preventing any more values from being pushed or popped.\n\nClosing a closed channel results in an unchecked error.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("close"), nil, nil, Void{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, Self{}, Never{})
				namespace.DefineMethod("Returns the amount of slots in the buffer that are available for new elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left_capacity"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the amount of elements present in the buffer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::Channel::V", env), NewSymbolLiteral("stop_iteration"))
				namespace.DefineMethod("Pops an element from the channel.\nBlocks the current thread until another thread pushes an element if the channel does not have any values in the buffer.\n\nPopping from a closed channel throws `:closed_channel`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("pop"), nil, nil, NameToType("Std::Channel::V", env), NewSymbolLiteral("closed_channel"))

				// Define constants

				// Define instance variables

				{
					namespace := namespace.Singleton()

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Create a new `Channel` that is closed.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("closed"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :closed", true), Never{}, Any{}, nil, INVARIANT)}, nil, NewGeneric(NameToType("Std::Channel", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("V"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :closed", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("V")})), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Char").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Concatenate this `Char`\nwith another `Char` or `String`.\n\nCreates a new `String` containing the content\nof both operands.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Get the next Unicode codepoint by incrementing by 1.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::Char", env), Never{})
				namespace.DefineMethod("Get the previous Unicode codepoint by decrementing by 1.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::Char", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Compare this char with another char or string.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Get the number of bytes that this\ncharacter contains.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("byte_count"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Always returns 1.\nFor better compatibility with `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("char_count"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Concatenate this `Char`\nwith another `Char` or `String`.\n\nCreates a new `String` containing the content\nof both operands.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("concat"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Always returns 1.\nFor better compatibility with `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("grapheme_count"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Calculates a hash of the char.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Always returns false.\nFor better compatibility with `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_empty"), nil, nil, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Always returns 1.\nFor better compatibility with `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Return the lowercase version of this character.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("lowercase"), nil, nil, NameToType("Std::Char", env), Never{})
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("repeat"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::CharLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::CharLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::CharLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::CharLiteralNode", env), Never{})
				namespace.DefineMethod("Creates a new string with this character.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Converts the `Char` to a `Symbol`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_symbol"), nil, nil, NameToType("Std::Symbol", env), Never{})
				namespace.DefineMethod("Return the uppercase version of this character.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("uppercase"), nil, nil, NameToType("Std::Char", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Class").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the name of the class.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the superclass (parent class) of this class.\nReturns `nil` when the class does not inherit from any class.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("superclass"), nil, nil, NewNilable(NameToType("Std::Class", env)), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Closable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("close"), nil, nil, Void{}, Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("ClosedRange").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
				typeParam.UpperBound = NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::ClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Range", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::ClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))

				// extend where Val < Std::Incrementable[Val] & Std::Comparable[Val]
				mixin = NewMixin("", false, "", env)
				{
					namespace := mixin
					namespace.Name() // noop - avoid unused variable error
					namespace.DefineSubtype(value.ToSymbol("Val"), NewTypeParameter(value.ToSymbol("Val"), NameToType("Std::ClosedRange", env).(*Class), Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::ClosedRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::ClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT))

					// Define methods
					namespace.DefineMethod("Returns the iterator for this range.\nOnly ranges of incrementable values can be iterated over.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::ClosedRange::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				}
				IncludeMixinWithWhere(namespace, mixin, []*TypeParameter{NewTypeParameter(value.ToSymbol("Val"), mixin, Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::ClosedRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::ClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT)})

				// Define methods
				namespace.DefineMethod("Check whether the given `value` is present in this range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::ClosedRange::Val", env), NameToType("Std::ClosedRange::Val", env), NameToType("Std::ClosedRange::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::ClosedRange::Val", env), NameToType("Std::ClosedRange::Val", env), NameToType("Std::ClosedRange::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Returns the upper bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("end"), nil, nil, NameToType("Std::ClosedRange::Val", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_closed"), nil, nil, True{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_closed"), nil, nil, True{}, Never{})
				namespace.DefineMethod("Returns the lower bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start"), nil, nil, NameToType("Std::ClosedRange::Val", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
					typeParam.UpperBound = NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::ClosedRange::Iterator::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::ClosedRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ClosedRange::Iterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("range"), NewGeneric(NameToType("Std::ClosedRange", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ClosedRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::ClosedRange::Iterator::Val", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Collection").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NewGeneric(NameToType("Std::ImmutableCollection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Collection::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

				// Define methods
				namespace.DefineMethod("Adds the given value to the collection.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Collection::Val", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Adds the given values to the collection.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("append"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::Collection::Val", env), PositionalRestParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Removes all elements from the collection.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("clear"), nil, nil, Void{}, Never{})
				namespace.DefineMethod("Returns the number of elements present in the collection.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over the elements of this collection,\nyielding them to the given closure.\n\nReturns a new collection that consists of the elements returned\nby the given closure.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Collection::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Collection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Adds the given value to the collection.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("push"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Collection::Val", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Removes the element from the collection.\n\nReturns `true` if the element has been removed,\notherwise returns `false.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("remove"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Collection::Val", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Removes all given elements from the collection.\n\nReturns `true` if any elements have been removed,\notherwise returns `false.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("remove_all"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::Collection::Val", env), PositionalRestParameterKind, false)}, Bool{}, Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Base").(*Mixin)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					ImplementInterface(namespace, NewGeneric(NameToType("Std::Collection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Collection::Base::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ImmutableCollection::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Collection::Base::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Collection::Base::Val", env), NormalParameterKind, false)}, Self{}, Never{})
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("append"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::Collection::Base::Val", env), PositionalRestParameterKind, false)}, Self{}, Never{})
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Collection::Base::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Collection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("remove_all"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::Collection::Base::Val", env), PositionalRestParameterKind, false)}, Bool{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Colorizer").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("colorize"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Comparable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("T"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("T"), typeParam)
				namespace.DefineConstant(value.ToSymbol("T"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Check if `self` is less than `other`", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Comparable::T", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check if `self` is less than or equal to `other`", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Comparable::T", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Returns:\n\n- `1` if `self` is greater than `other`\n- `0` if both are equal.\n- `-1` if `self` is less than `other`.\n– `nil` if the comparison was impossible (NaN)", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Comparable::T", env), NormalParameterKind, false)}, NewNilable(NameToType("Std::Int", env)), Never{})
				namespace.DefineMethod("Check if `self` is greater than `other`", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Comparable::T", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check if `self` is greater than or equal to `other`", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Comparable::T", env), NormalParameterKind, false)}, Bool{}, Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Container").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Check whether the given `value` is present in this container.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Container::Val", env), NameToType("Std::Container::Val", env), NameToType("Std::Container::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Container::Val", env), NameToType("Std::Container::Val", env), NameToType("Std::Container::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Debug").(*Module)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect_call_stack"), nil, nil, Void{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect_value_stack"), nil, nil, Void{}, Never{})
				namespace.DefineMethod("Returns the current stack trace.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("stack_trace"), nil, nil, NameToType("Std::StackTrace", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start_cpu_profile"), nil, []*Parameter{NewParameter(value.ToSymbol("file_path"), NameToType("Std::String", env), NormalParameterKind, false)}, Void{}, NameToType("Std::FileSystemError", env))
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("stop_cpu_profile"), nil, nil, Void{}, Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Decrementable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("T"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("T"), typeParam)
				namespace.DefineConstant(value.ToSymbol("T"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Get the previous value of this type, the predecessor of `self`", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::Decrementable::T", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Diagnostic").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("message"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), NormalParameterKind, false), NewParameter(value.ToSymbol("severity"), NameToType("Std::UInt8", env), DefaultValueParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("message"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Indicates whether this Diagnostic\nrepresents an info message, a warning or a failure.\n\nValid values are defined as constants:\n- `INFO`\n- `WARN`\n- `FAIL`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("severity"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("severity_name"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Get the human readable\nstring representation of this diagnostic\nwith a snippet of source code.\n\nIt will attempt to read the source fragment from the file\nif no source string is given.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_human_string"), nil, []*Parameter{NewParameter(value.ToSymbol("style"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("colorizer"), NewNilable(NameToType("Std::Colorizer", env)), DefaultValueParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Get the string representation of this diagnostic.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, Void{}, Never{})

				// Define constants
				namespace.DefineConstant(value.ToSymbol("FAIL"), NameToType("Std::UInt8", env))
				namespace.DefineConstant(value.ToSymbol("INFO"), NameToType("Std::UInt8", env))
				namespace.DefineConstant(value.ToSymbol("WARN"), NameToType("Std::UInt8", env))

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("DiagnosticList").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::List", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Diagnostic", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

				// Define methods
				namespace.DefineMethod("Create a new list containing the elements of `self`\nrepeated `n` times.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Create a new `ArrayList` containing the elements of `self`\nand another given `Tuple`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::Diagnostic", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Adds the given value to the list.\n\nReallocates the underlying array if it is\ntoo small to hold it.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Check whether the given value is an `ArrayList`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Check whether the given value is an `ArrayList` or `ArrayTuple`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Get the element under the given index.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Diagnostic", env), Never{})
				namespace.DefineMethod("Set the element under the given index to the given value.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]="), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Adds the given values to the list.\n\nReallocates the underlying array if it is\ntoo small to hold them.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("append"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::Diagnostic", env), PositionalRestParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Returns the number of elements that can be\nheld by the underlying array.\n\nThis value will change when the list gets resized,\nand the underlying array gets reallocated", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Removes all elements from the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("clear"), nil, nil, Void{}, Never{})
				namespace.DefineMethod("Check whether the given `value` is present in this list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Diagnostic", env), NameToType("Std::Diagnostic", env), NameToType("Std::Diagnostic", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Diagnostic", env), NameToType("Std::Diagnostic", env), NameToType("Std::Diagnostic", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Mutates the list.\n\nReallocates the underlying array to hold\nthe given number of new elements.\n\nExpands the `capacity` of the list by `new_slots`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("grow"), nil, []*Parameter{NewParameter(value.ToSymbol("new_slots"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Check if this diagnostic list contains a fail.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_failure"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NameToType("Std::DiagnosticList::Iterator", env), Never{})
				namespace.DefineMethod("Returns the number of left slots for new elements\nin the underlying array.\nIt tells you how many more elements can be\nadded to the list before the underlying array gets\nreallocated.\n\nIt is always equal to `capacity - length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left_capacity"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of elements present in the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over the elements of this list,\nyielding them to the given closure.\n\nReturns a new List that consists of the elements returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the elements of this list,\nyielding them to the given closure.\n\nMutates the list in place replacing the elements with the ones\nreturned by the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_mut"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, NameToType("Std::Diagnostic", env), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Adds the given value to the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("push"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Removes the element from the list.\n\nReturns `true` if the element has been removed,\notherwise returns `false.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("remove"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Get the human readable\nstring representation of this diagnostic list\nwith a snippet of source code.\n\nIt will attempt to read the source fragment from the file\nif no source string is given.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_human_string"), nil, []*Parameter{NewParameter(value.ToSymbol("style"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("colorizer"), NewNilable(NameToType("Std::Colorizer", env)), DefaultValueParameterKind, false)}, Void{}, Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Diagnostic", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("list"), NameToType("Std::DiagnosticList", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::Diagnostic", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Duration").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Duration", env), NormalParameterKind, false)}, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Duration", env), NormalParameterKind, false)}, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the count of days in this duration as an Int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("days"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the count of hours in this duration as an Int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hours"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the count of days in this duration as a Float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_days"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns the count of hours in this duration as a Float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_hours"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns the count of microseconds in this duration as an Int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_microseconds"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns the count of milliseconds in this duration as a Float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_milliseconds"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns the count of minutes in this duration as a Float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_minutes"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns the count of nanoseconds in this duration as a Float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_nanoseconds"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns the count of seconds in this duration as a Float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_seconds"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns the count of weeks in this duration as a Float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_weeks"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns the count of years in this duration as a Float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_years"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns the count of microseconds in this duration as a Float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("microseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the count of milliseconds in this duration as an Int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("milliseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the count of minutes in this duration as an Int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("minutes"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the count of nanoseconds in this duration as an Int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("nanoseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the count of seconds in this duration as an Int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("seconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the string representation of the duration in the format \"51h15m0.12s\".", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the count of weeks in this duration as an Int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("weeks"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the count of years in this duration as an Int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("years"), nil, nil, NameToType("Std::Int", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.Singleton()

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Parses a duration string and creates a Duration value.\nA duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \"300ms\", \"-1.5h\" or \"2h45m\".\nValid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".", 0|METHOD_NATIVE_FLAG, value.ToSymbol("parse"), nil, []*Parameter{NewParameter(value.ToSymbol("s"), NameToType("Std::String", env), NormalParameterKind, false)}, NameToType("Std::Duration", env), Never{})
					namespace.DefineMethod("Returns the amount of elapsed since the given `Time`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("since"), nil, []*Parameter{NewParameter(value.ToSymbol("time"), NameToType("Std::Time", env), NormalParameterKind, false)}, NameToType("Std::Duration", env), Never{})
					namespace.DefineMethod("Returns the amount of time that is left until the given `Time`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("until"), nil, []*Parameter{NewParameter(value.ToSymbol("time"), NameToType("Std::Time", env), NormalParameterKind, false)}, NameToType("Std::Duration", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Elk").(*Module)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("AST").(*Module)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods

					// Define constants

					// Define instance variables

					{
						namespace := namespace.MustSubtype("AliasDeclarationEntry").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("new_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("old_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("new_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("old_name"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("AliasDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("entries"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::AliasDeclarationEntry", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("entries"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::AliasDeclarationEntry", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("AnyTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ArrayListLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ArrayTupleLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("AsExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("runtime_type"), NameToType("Std::Elk::AST::ComplexConstantNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("runtime_type"), nil, nil, NameToType("Std::Elk::AST::ComplexConstantNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("AsPatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("pattern"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("name"), NameToType("Std::Elk::AST::IdentifierNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::Elk::AST::IdentifierNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("pattern"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("AssignmentExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("left"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("right"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("right"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("AttrDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("entries"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("entries"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("AttributeAccessNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("receiver"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("attribute_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("attribute_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("receiver"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("AttributeParameterNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ParameterNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("initialiser"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("initialiser"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_named_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_normal"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the parameter is optional.\neg. `foo?: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_optional"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_positional_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("AwaitExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("BigFloatLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("BinArrayListLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("BinArrayTupleLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("BinHashSetLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("BinaryExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("left"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("right"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("right"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("BinaryPatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("left"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("right"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("right"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("BinaryTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("left"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("right"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("right"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("BoolLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("BreakExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("label"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("value"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("label"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("CallNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("receiver"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("nil_safe"), Bool{}, NormalParameterKind, false), NewParameter(value.ToSymbol("positional_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("named_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("named_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("nil_safe"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("positional_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("receiver"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("CaseNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("pattern"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("pattern"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("CatchNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("pattern"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("stack_trace_var"), NewNilable(NameToType("Std::Elk::AST::IdentifierNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("pattern"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("stack_trace_var"), nil, nil, NewNilable(NameToType("Std::Elk::AST::IdentifierNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("CharLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Char", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Char", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ClassDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constant"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("type_parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("is_abstract"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("is_sealed"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("is_primitive"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("is_no_init"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("superclass"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constant"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::Bool", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_abstract"), nil, nil, NameToType("Std::Bool", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_no_init"), nil, nil, NameToType("Std::Bool", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_primitive"), nil, nil, NameToType("Std::Bool", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_sealed"), nil, nil, NameToType("Std::Bool", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("superclass"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ClosureLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("return_type"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("throw_type"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("return_type"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("throw_type"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ClosureTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("return_type"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("throw_type"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("return_type"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("throw_type"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ComplexConstantNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ConstantAsNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constant"), NameToType("Std::Elk::AST::ComplexConstantNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("as_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("as_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constant"), nil, nil, NameToType("Std::Elk::AST::ComplexConstantNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ConstantDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constant"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("initialiser"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constant"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("initialiser"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ConstantLookupNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ComplexConstantNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("left"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), NormalParameterKind, false), NewParameter(value.ToSymbol("right"), NameToType("Std::Elk::AST::ComplexConstantNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("right"), nil, nil, NameToType("Std::Elk::AST::ComplexConstantNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ConstantNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ComplexConstantNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ConstructorCallNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("class_node"), NameToType("Std::Elk::AST::ComplexConstantNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("positional_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("named_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("class_node"), nil, nil, NameToType("Std::Elk::AST::ComplexConstantNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("named_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("positional_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ContinueExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("label"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("label"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("DoExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("catches"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::CatchNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("finally_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("catches"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::CatchNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("finally_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("DoubleQuotedStringLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StringLiteralNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::SimpleStringLiteralNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("DoubleSplatExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::NamedArgumentNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("EmptyStatementNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StatementNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StructBodyStatementNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ExpressionNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables

						{
							namespace := namespace.MustSubtype("Convertible").(*Interface)

							namespace.Name() // noop - avoid unused variable error

							// Include mixins and implement interfaces

							// Define methods
							namespace.DefineMethod("Converts the value to an Elk AST Expression Node.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

							// Define constants

							// Define instance variables
						}
					}
					{
						namespace := namespace.MustSubtype("ExpressionStatementNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StatementNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ExtendWhereBlockExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("where_params"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("where_params"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("FalseLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("Float32LiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("Float64LiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("FloatLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ForInExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("pattern"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("in_expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("then_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("pattern"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("then_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("FormalParameterNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ParameterNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("kind"), NameToType("Std::UInt8", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("initialiser"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("initialiser"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("Whether it's a named rest parameter eg `**foo: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_named_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the parameter is normal\n(not a positional rest parameter, not a named rest parameter).\n\neg. `foo: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_normal"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the parameter is optional (has a default value).\neg. `foo: String = \"lol\"`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_optional"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether it's a positional rest parameter eg `*foo: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_positional_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Returns an integer that determines whether\nthe parameter is a normal, positional rest or named rest parameter.\n\nValid values are available as constants under `Std::Elk::AST::ParameterNode`:\n- `NORMAL_KIND`\n- `POSITIONAL_REST_KIND`\n- `NAMED_REST_KIND`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("kind"), nil, nil, NameToType("Std::UInt8", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("GenericConstantNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ComplexConstantNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constant"), NameToType("Std::Elk::AST::ComplexConstantNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constant"), nil, nil, NameToType("Std::Elk::AST::ComplexConstantNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("GenericConstructorCallNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("class_node"), NameToType("Std::Elk::AST::ComplexConstantNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("positional_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("named_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("class_node"), nil, nil, NameToType("Std::Elk::AST::ComplexConstantNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("named_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("positional_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("GenericMethodCallNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("receiver"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("method_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("positional_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("named_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("method_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("named_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("positional_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("receiver"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("GenericReceiverlessMethodCallNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("method_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("positional_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("named_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("method_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("named_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("positional_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("GenericTypeDefinitionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constant"), NameToType("Std::Elk::AST::ComplexConstantNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constant"), nil, nil, NameToType("Std::Elk::AST::ComplexConstantNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("GetterDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("entries"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("entries"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("GoExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("HashMapLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("HashRecordLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("HashSetLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("HexArrayListLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("HexArrayTupleLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("HexHashSetLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::IntCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("IdentifierNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("IfExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("condition"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("then_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("else_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("condition"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("else_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("then_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ImplementExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constants"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ComplexConstantNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constants"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ComplexConstantNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ImportStatementNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StatementNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("path"), NameToType("Std::Elk::AST::StringLiteralNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("path"), nil, nil, NameToType("Std::Elk::AST::StringLiteralNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("IncludeExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constants"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ComplexConstantNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constants"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ComplexConstantNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("InitDefinitionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("throw_type"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("throw_type"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("InstanceOfTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("type_node"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("InstanceVariableDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("InstanceVariableNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("Int16LiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("Int32LiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("Int64LiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("Int8LiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("IntCollectionContentNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("IntLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::IntCollectionContentNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("InterfaceDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constant"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constant"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("InterpolatedRegexLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::RegexLiteralNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("content"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::RegexLiteralContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("flags"), NameToType("Std::UInt8", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("content"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::RegexLiteralContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("Flags of the regex literal.\n\nThey're defined as constants\nunder the `RegexLiteralNode` namespace:\n\n- `CASE_INSENSITIVE_FLAG`\n- `MULTILINE_FLAG`\n- `DOT_ALL_FLAG`\n- `UNGREEDY_FLAG`\n- `EXTENDED_FLAG`\n- `ASCII_FLAG`\n\nExample usage:\n\n```\nusing Elk::AST::RegexLiteralNode::*\n\n# check if a flag is set\nnode.flags & MULTILINE_FLAG != 0\n\n# combine flags\nnode.flags = UNGREEDY_FLAG | MULTILINE_FLAG\n```", 0|METHOD_NATIVE_FLAG, value.ToSymbol("flags"), nil, nil, NameToType("Std::UInt8", env), Never{})
						namespace.DefineMethod("Whether the regex has the `a` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_ascii"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the regex has the `i` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_case_insensitive"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the regex has the `s` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_dot_all"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the regex has the `x` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_extended"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the regex has the `m` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_multiline"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the regex has the `U` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_ungreedy"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("InterpolatedStringLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StringLiteralNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("content"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StringLiteralContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("content"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StringLiteralContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("InterpolatedSymbolLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::SymbolLiteralNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("content"), NameToType("Std::Elk::AST::InterpolatedStringLiteralNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("content"), nil, nil, NameToType("Std::Elk::AST::InterpolatedStringLiteralNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("IntersectionTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("InvalidNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StatementNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StructBodyStatementNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ParameterNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingSubentryNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::NamedArgumentNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ComplexConstantNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ConstantNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::WordCollectionContentNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::SymbolCollectionContentNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::IntCollectionContentNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::IdentifierNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::RegexLiteralContentNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::RegexLiteralNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StringLiteralContentNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StringLiteralNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::SimpleStringLiteralNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::SymbolLiteralNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeParameterNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("token"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_named_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_normal"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_optional"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_positional_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("token"), nil, nil, NameToType("Std::Elk::Token", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("KeyValueExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("key"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("KeyValuePatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::Elk::AST::PatternExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("key"), nil, nil, NameToType("Std::Elk::AST::PatternExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("LabeledExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("label"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("label"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ListPatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("LogicalExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("left"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("right"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("right"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("LoopExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("then_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("then_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MacroBoundaryNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MacroCallNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("receiver"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("macro_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("positional_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("named_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("macro_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("named_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("positional_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("receiver"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MacroDefinitionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("is_sealed"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("Whether the method is sealed eg. `sealed def foo; end`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_sealed"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MapPatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MethodCallNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("receiver"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("method_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("positional_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("named_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("method_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("named_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("positional_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("receiver"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MethodDefinitionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("return_type"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("throw_type"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("flags"), NameToType("Std::UInt8", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("flags"), nil, nil, NameToType("Std::UInt8", env), Never{})
						namespace.DefineMethod("Whether the method is abstract eg. `abstract def foo; end`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_abstract"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the method is asynchronous eg. `async def foo; end`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_async"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the method is a generator eg. `def* foo; end`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_generator"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the method is sealed eg. `sealed def foo; end`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_sealed"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("return_type"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("throw_type"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants
						namespace.DefineConstant(value.ToSymbol("ABSTRACT_FLAG"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("ASYNC_FLAG"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("GENERATOR_FLAG"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("SEALED_FLAG"), NameToType("Std::UInt8", env))

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MethodLookupAsNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("method_lookup"), NameToType("Std::Elk::AST::MethodLookupNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("as_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("as_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("method_lookup"), nil, nil, NameToType("Std::Elk::AST::MethodLookupNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MethodLookupNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("receiver"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("receiver"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MethodParameterNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ParameterNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("initialiser"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("set_instance_variable"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("kind"), NameToType("Std::UInt8", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("initialiser"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("Whether it's a named rest parameter eg `**foo: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_named_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the parameter is normal\n(not a positional rest parameter, not a named rest parameter).\n\neg. `foo: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_normal"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the parameter is optional (has a default value).\neg. `foo: String = \"lol\"`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_optional"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether it's a positional rest parameter eg `*foo: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_positional_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Returns an integer that determines whether\nthe parameter is a normal, positional rest or named rest parameter.\n\nValid values are available as constants under `Std::Elk::AST::ParameterNode`:\n- `NORMAL_KIND`\n- `POSITIONAL_REST_KIND`\n- `NAMED_REST_KIND`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("kind"), nil, nil, NameToType("Std::UInt8", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("set_instance_variable"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MethodSignatureDefinitionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("return_type"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("throw_type"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("return_type"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("throw_type"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MixinDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constant"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("is_abstract"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("type_parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constant"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_abstract"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ModifierForInNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("then_expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("pattern"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("in_expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("in_expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("pattern"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("then_expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ModifierIfElseNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("then_expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("condition"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("else_expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("condition"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("else_expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("then_expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ModifierNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("modifier"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("left"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("right"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("modifier"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("right"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ModuleDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constant"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constant"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("MustExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("NamedArgumentNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("NamedCallArgumentNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::NamedArgumentNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("NeverTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("NewExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("positional_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("named_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("named_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("positional_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("NilLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ComplexConstantNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("NilSafeSubscriptExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("receiver"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("key"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("key"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("receiver"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("NilableTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("type_node"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("Node").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("Create a new iterator by eagerly traversing the entire AST\nand collecting its nodes to a collection.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::Iterator", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::Node", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})), Never{})
						namespace.DefineMethod("Returns the span that represents\nthe position of this node in a source file/string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("Returns self.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast"), nil, nil, Self{}, Never{})
						namespace.DefineMethod("Convert the AST back to a `String` of source code.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("Visit all nodes of this AST\npassing them as arguments to the given closure.\n\nThe closure should return `true` to continue\ntraversing the tree, `false` will be treated as a `break` and will\nstop the traversal.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("traverse"), nil, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("node"), NameToType("Std::Elk::AST::Node", env), NormalParameterKind, false)}, Bool{}, Never{}), NormalParameterKind, false)}, Bool{}, Never{})

						// Define constants

						// Define instance variables

						{
							namespace := namespace.MustSubtype("Convertible").(*Interface)

							namespace.Name() // noop - avoid unused variable error

							// Include mixins and implement interfaces

							// Define methods
							namespace.DefineMethod("Converts the value to an Elk AST Node.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::Node", env), Never{})

							// Define constants

							// Define instance variables
						}
					}
					{
						namespace := namespace.MustSubtype("NotTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("type_node"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("NumericForExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("initialiser"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("condition"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("increment"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("then_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("condition"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("increment"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("initialiser"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("then_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ObjectPatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("object_type"), NameToType("Std::Elk::AST::ComplexConstantNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("attributes"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("attributes"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("object_type"), nil, nil, NameToType("Std::Elk::AST::ComplexConstantNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ParameterNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("Whether it's a named rest parameter eg `**foo: String`", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_named_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the parameter is normal\n(not a positional rest parameter, not a named rest parameter).\n\neg. `foo: String`", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_normal"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the parameter is optional (has a default value).\neg. `foo: String = \"lol\"`", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_optional"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether it's a positional rest parameter eg `*foo: String`", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_positional_rest"), nil, nil, Bool{}, Never{})

						// Define constants
						namespace.DefineConstant(value.ToSymbol("NAMED_REST_KIND"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("NORMAL_KIND"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("POSITIONAL_REST_KIND"), NameToType("Std::UInt8", env))

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ParameterStatementNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StructBodyStatementNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("parameter"), NameToType("Std::Elk::AST::ParameterNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("parameter"), nil, nil, NameToType("Std::Elk::AST::ParameterNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("PatternExpressionNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("PatternNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables

						{
							namespace := namespace.MustSubtype("Convertible").(*Interface)

							namespace.Name() // noop - avoid unused variable error

							// Include mixins and implement interfaces

							// Define methods
							namespace.DefineMethod("Converts the value to an Elk AST Pattern Node.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})

							// Define constants

							// Define instance variables
						}
					}
					{
						namespace := namespace.MustSubtype("PostfixExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("PrivateConstantNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ComplexConstantNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ConstantNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("PrivateIdentifierNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::IdentifierNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ProgramNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("PublicConstantAsNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingSubentryNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("target"), NameToType("Std::Elk::AST::PublicConstantNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("as_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("as_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("target"), nil, nil, NameToType("Std::Elk::AST::PublicConstantNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("PublicConstantNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingSubentryNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ComplexConstantNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ConstantNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("PublicIdentifierAsNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingSubentryNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("target"), NameToType("Std::Elk::AST::PublicIdentifierNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("as_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("as_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("target"), nil, nil, NameToType("Std::Elk::AST::PublicIdentifierNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("PublicIdentifierNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingSubentryNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::IdentifierNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("QuoteExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("RangeLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("start"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("end_node"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("end"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("RawCharLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Char", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Char", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("RawStringLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::WordCollectionContentNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StringLiteralNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::SimpleStringLiteralNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ReceiverlessMacroCallNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("macro_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("positional_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("named_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("macro_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("named_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("positional_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ReceiverlessMethodCallNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("method_name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("positional_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("named_arguments"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("method_name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("named_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::NamedArgumentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("positional_arguments"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ExpressionNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("RecordPatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("RegexInterpolationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::RegexLiteralContentNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("RegexLiteralContentNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("RegexLiteralContentSectionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::RegexLiteralContentNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("RegexLiteralNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods

						// Define constants
						namespace.DefineConstant(value.ToSymbol("ASCII_FLAG"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("CASE_INSENSITIVE_FLAG"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("DOT_ALL_FLAG"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("EXTENDED_FLAG"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("MULTILINE_FLAG"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("UNGREEDY_FLAG"), NameToType("Std::UInt8", env))

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("RestPatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("identifier"), NameToType("Std::Elk::AST::IdentifierNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("identifier"), nil, nil, NameToType("Std::Elk::AST::IdentifierNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ReturnExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SelfLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SetPatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SetterDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("entries"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("entries"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::ParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SignatureParameterNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ParameterNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("is_optional"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("kind"), NameToType("Std::UInt8", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("Whether it's a named rest parameter eg `**foo: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_named_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the parameter is normal\n(not a positional rest parameter, not a named rest parameter).\n\neg. `foo: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_normal"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the parameter is optional.\neg. `foo?: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_optional"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether it's a positional rest parameter eg `*foo: String`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_positional_rest"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Returns an integer that determines whether\nthe parameter is a normal, positional rest or named rest parameter.\n\nValid values are available as constants under `Std::Elk::AST::ParameterNode`:\n- `NORMAL_KIND`\n- `POSITIONAL_REST_KIND`\n- `NAMED_REST_KIND`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("kind"), nil, nil, NameToType("Std::UInt8", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SimpleStringLiteralNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StringLiteralNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SimpleSymbolLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::SymbolCollectionContentNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::SymbolLiteralNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("content"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("content"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SingletonBlockExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SingletonTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("type_node"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SplatExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("StatementNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("StringInspectInterpolationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StringLiteralContentNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("StringInterpolationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StringLiteralContentNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("StringLiteralContentNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("StringLiteralContentSectionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::StringLiteralContentNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("StringLiteralNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("StructBodyStatementNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("StructDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constant"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_parameters"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StructBodyStatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StructBodyStatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constant"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_parameters"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeParameterNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SubscriptExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("receiver"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("key"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("key"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("receiver"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SwitchExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("cases"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::CaseNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("else_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("cases"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::CaseNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("else_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SymbolArrayListLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::SymbolCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::SymbolCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SymbolArrayTupleLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::SymbolCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::SymbolCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SymbolCollectionContentNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SymbolHashSetLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::SymbolCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::SymbolCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SymbolKeyValueExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("key"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SymbolKeyValuePatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("key"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("SymbolLiteralNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ThrowExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("is_unchecked"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_unchecked"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("TrueLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("TryExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("TuplePatternNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::PatternNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("TypeDefinitionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("constant"), NameToType("Std::Elk::AST::ComplexConstantNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("constant"), nil, nil, NameToType("Std::Elk::AST::ComplexConstantNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("TypeExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("type_node"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("TypeNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables

						{
							namespace := namespace.MustSubtype("Convertible").(*Interface)

							namespace.Name() // noop - avoid unused variable error

							// Include mixins and implement interfaces

							// Define methods
							namespace.DefineMethod("Converts the value to an Elk AST Type Node.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

							// Define constants

							// Define instance variables
						}
					}
					{
						namespace := namespace.MustSubtype("TypeParameterNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::Node", env).(*Mixin))

						// Define methods

						// Define constants
						namespace.DefineConstant(value.ToSymbol("CONTRAVARIANT"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("COVARIANT"), NameToType("Std::UInt8", env))
						namespace.DefineConstant(value.ToSymbol("INVARIANT"), NameToType("Std::UInt8", env))

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("TypeofExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UInt16LiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UInt32LiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UInt64LiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UInt8LiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UnaryExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("right"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("right"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UnaryTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("op"), NameToType("Std::Elk::Token", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NameToType("Std::Elk::AST::TypeNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("op"), nil, nil, NameToType("Std::Elk::Token", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NameToType("Std::Elk::AST::TypeNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UndefinedLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UninterpolatedRegexLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::RegexLiteralNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("content"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("flags"), NameToType("Std::UInt8", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("content"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("Flags of the regex literal.\n\nThey're defined as constants\nunder the `RegexLiteralNode` namespace:\n\n- `CASE_INSENSITIVE_FLAG`\n- `MULTILINE_FLAG`\n- `DOT_ALL_FLAG`\n- `UNGREEDY_FLAG`\n- `EXTENDED_FLAG`\n- `ASCII_FLAG`\n\nExample usage:\n\n```\nusing Elk::AST::RegexLiteralNode::*\n\n# check if a flag is set\nnode.flags & MULTILINE_FLAG != 0\n\n# combine flags\nnode.flags = UNGREEDY_FLAG | MULTILINE_FLAG\n```", 0|METHOD_NATIVE_FLAG, value.ToSymbol("flags"), nil, nil, NameToType("Std::UInt8", env), Never{})
						namespace.DefineMethod("Whether the regex has the `a` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_ascii"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the regex has the `i` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_case_insensitive"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the regex has the `s` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_dot_all"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the regex has the `x` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_extended"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the regex has the `m` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_multiline"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the regex has the `U` flag.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_ungreedy"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UnionTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::TypeNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UnlessExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("condition"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("then_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("else_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("condition"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("else_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("then_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UnquoteExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("expression"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("expression"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UntilExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("condition"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("then_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("condition"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("then_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UsingAllEntryNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("namespace"), NameToType("Std::Elk::AST::UsingEntryNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("namespace"), nil, nil, NameToType("Std::Elk::AST::UsingEntryNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UsingEntryNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UsingEntryWithSubentriesNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::UsingEntryNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("namespace"), NameToType("Std::Elk::AST::UsingEntryNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("subentries"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::UsingSubentryNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("namespace"), nil, nil, NameToType("Std::Elk::AST::UsingEntryNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("subentries"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::UsingSubentryNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UsingExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("entries"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::UsingEntryNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("entries"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::UsingEntryNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("UsingSubentryNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ValueDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("initialiser"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("initialiser"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("ValuePatternDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("pattern"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("initialiser"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("initialiser"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("pattern"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("VariableDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("type_node"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("initialiser"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("doc_comment"), NameToType("Std::String", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("doc_comment"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("initialiser"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_node"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("VariablePatternDeclarationNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("pattern"), NameToType("Std::Elk::AST::PatternNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("initialiser"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("initialiser"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("pattern"), nil, nil, NameToType("Std::Elk::AST::PatternNode", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("VariantTypeParameterNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeParameterNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("lower_bound"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("upper_bound"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("default_node"), NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("variance"), NameToType("Std::UInt8", env), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("default_node"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})
						namespace.DefineMethod("Whether the type parameter is contravariant\neg. `-V` in `class Foo[-V]; end`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_contravariant"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the type parameter is covariant\neg. `+V` in `class Foo[+V]; end`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_covariant"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("Whether the type parameter is invariant\neg. `V` in `class Foo[V]; end`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_invariant"), nil, nil, Bool{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("lower_bound"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("upper_bound"), nil, nil, NewNilable(NameToType("Std::Elk::AST::TypeNode", env)), Never{})
						namespace.DefineMethod("Returns an integer that determines whether\nthe variance of the type parameter.\n\nValid values are available as constants under `Std::Elk::AST::TypeParameterNode`:\n- `INVARIANT`\n- `COVARIANT`\n- `CONTRAVARIANT`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("variance"), nil, nil, NameToType("Std::UInt8", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("VoidTypeNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::TypeNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("WhileExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("condition"), NameToType("Std::Elk::AST::ExpressionNode", env), NormalParameterKind, false), NewParameter(value.ToSymbol("then_body"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("condition"), nil, nil, NameToType("Std::Elk::AST::ExpressionNode", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("then_body"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::StatementNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("WordArrayListLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::WordCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::WordCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("WordArrayTupleLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::WordCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::WordCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("WordCollectionContentNode").(*Mixin)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("WordHashSetLiteralNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))
						IncludeMixin(namespace, NameToType("Std::Elk::AST::PatternNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::WordCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("capacity"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("elements"), nil, nil, NewGeneric(NameToType("Std::ArrayTuple", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::AST::WordCollectionContentNode", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("YieldExpressionNode").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NameToType("Std::Elk::AST::ExpressionNode", env).(*Mixin))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), DefaultValueParameterKind, false), NewParameter(value.ToSymbol("forward"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), DefaultValueParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("forward"), nil, nil, NameToType("Std::Bool", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NewNilable(NameToType("Std::Elk::AST::ExpressionNode", env)), Never{})

						// Define constants

						// Define instance variables
					}
				}
				{
					namespace := namespace.MustSubtype("Lexer").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("source"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("source_name"), NameToType("Std::String", env), DefaultValueParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next token.\nToken will be of type `END_OF_FILE` when\nthe end of source is reached.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::Elk::Token", env), Never{})

					// Define constants

					// Define instance variables

					{
						namespace := namespace.Singleton()

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces

						// Define methods
						namespace.DefineMethod("Colorize the given string of Elk source code\nwith ANSI escape codes.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("colorize"), nil, []*Parameter{NewParameter(value.ToSymbol("source"), NameToType("Std::String", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
						namespace.DefineMethod("Tokenize the given string of Elk source code\nand return an array of tokens.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("lex"), nil, []*Parameter{NewParameter(value.ToSymbol("source"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("source_name"), NameToType("Std::String", env), DefaultValueParameterKind, false)}, NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Elk::Token", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

						// Define constants

						// Define instance variables
					}
				}
				{
					namespace := namespace.MustSubtype("Parser").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods

					// Define constants

					// Define instance variables

					{
						namespace := namespace.Singleton()

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("parse"), nil, []*Parameter{NewParameter(value.ToSymbol("source"), NameToType("Std::String", env), NormalParameterKind, false), NewParameter(value.ToSymbol("source_name"), NameToType("Std::String", env), DefaultValueParameterKind, false)}, NameToType("Std::Elk::Parser::Result", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("Result").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces

						// Define methods
						namespace.DefineMethod("Returns the Abstract Syntax Tree produced by the parser.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("ast"), nil, nil, NameToType("Std::Elk::AST::ProgramNode", env), Never{})
						namespace.DefineMethod("Returns a list of diagnostics (warnings, errors, info messages)", 0|METHOD_NATIVE_FLAG, value.ToSymbol("diagnostics"), nil, nil, NameToType("Std::DiagnosticList", env), Never{})

						// Define constants

						// Define instance variables
					}
				}
				{
					namespace := namespace.MustSubtype("Token").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("typ"), NameToType("Std::UInt16", env), NormalParameterKind, false), NewParameter(value.ToSymbol("location"), NameToType("Std::FS::Location", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::String", env), DefaultValueParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Returns the lexeme, a piece of source text\nthat this token represents.\n\nFetches the lexeme from a global map of lexemes\nfor Elk tokens if the value is an empty string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("fetch_value"), nil, nil, NameToType("Std::String", env), Never{})
					namespace.DefineMethod("Returns the span that represents the position of\nthe token in the source string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("location"), nil, nil, NameToType("Std::FS::Location", env), Never{})
					namespace.DefineMethod("Returns an integer that represents the type of the token.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("typ"), nil, nil, NameToType("Std::UInt16", env), Never{})
					namespace.DefineMethod("Get the token type name like `\"LBRACE\"`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_name"), nil, nil, NameToType("Std::String", env), Never{})
					namespace.DefineMethod("Returns the lexeme, a piece of source text\nthat this token represents.\n\nMay return an empty string for simple tokens\nlike operators where the lexeme is obvious (it will always be the same).", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::String", env), Never{})

					// Define constants
					namespace.DefineConstant(value.ToSymbol("NAMES_TO_TYPES"), NewGeneric(NameToType("Std::HashMap", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::String", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::UInt16", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})))
					namespace.DefineConstant(value.ToSymbol("TYPES_TO_NAMES"), NewGeneric(NameToType("Std::HashMap", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::UInt16", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::String", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})))

					// Define instance variables

					{
						namespace := namespace.Singleton()

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces

						// Define methods
						namespace.DefineMethod("Get the token type name of the given token type.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("type_name"), nil, []*Parameter{NewParameter(value.ToSymbol("typ"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})

						// Define constants

						// Define instance variables
					}
				}
			}
			{
				namespace := namespace.MustSubtype("EndlessClosedRange").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
				typeParam.UpperBound = NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Range", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))

				// extend where Val < Std::Incrementable[Val] & Std::Comparable[Val]
				mixin = NewMixin("", false, "", env)
				{
					namespace := mixin
					namespace.Name() // noop - avoid unused variable error
					namespace.DefineSubtype(value.ToSymbol("Val"), NewTypeParameter(value.ToSymbol("Val"), NameToType("Std::EndlessClosedRange", env).(*Class), Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT))

					// Define methods
					namespace.DefineMethod("Returns the iterator for this range.\nOnly ranges of incrementable values can be iterated over.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::EndlessClosedRange::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				}
				IncludeMixinWithWhere(namespace, mixin, []*TypeParameter{NewTypeParameter(value.ToSymbol("Val"), mixin, Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT)})

				// Define methods
				namespace.DefineMethod("Check whether the given `value` is present in this range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::EndlessClosedRange::Val", env), NameToType("Std::EndlessClosedRange::Val", env), NameToType("Std::EndlessClosedRange::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::EndlessClosedRange::Val", env), NameToType("Std::EndlessClosedRange::Val", env), NameToType("Std::EndlessClosedRange::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_closed"), nil, nil, True{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_closed"), nil, nil, False{}, Never{})
				namespace.DefineMethod("Returns the lower bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start"), nil, nil, NameToType("Std::EndlessClosedRange::Val", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
					typeParam.UpperBound = NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Iterator::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Iterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("range"), NewGeneric(NameToType("Std::EndlessClosedRange", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::EndlessClosedRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::EndlessClosedRange::Iterator::Val", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("EndlessOpenRange").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
				typeParam.UpperBound = NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Range", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))

				// extend where Val < Std::Incrementable[Val] & Std::Comparable[Val]
				mixin = NewMixin("", false, "", env)
				{
					namespace := mixin
					namespace.Name() // noop - avoid unused variable error
					namespace.DefineSubtype(value.ToSymbol("Val"), NewTypeParameter(value.ToSymbol("Val"), NameToType("Std::EndlessOpenRange", env).(*Class), Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT))

					// Define methods
					namespace.DefineMethod("Returns the iterator for this range.\nOnly ranges of incrementable values can be iterated over.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::EndlessOpenRange::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				}
				IncludeMixinWithWhere(namespace, mixin, []*TypeParameter{NewTypeParameter(value.ToSymbol("Val"), mixin, Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT)})

				// Define methods
				namespace.DefineMethod("Check whether the given `value` is present in this range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::EndlessOpenRange::Val", env), NameToType("Std::EndlessOpenRange::Val", env), NameToType("Std::EndlessOpenRange::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::EndlessOpenRange::Val", env), NameToType("Std::EndlessOpenRange::Val", env), NameToType("Std::EndlessOpenRange::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_closed"), nil, nil, False{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_closed"), nil, nil, False{}, Never{})
				namespace.DefineMethod("Returns the lower bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start"), nil, nil, NameToType("Std::EndlessOpenRange::Val", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
					typeParam.UpperBound = NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Iterator::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Iterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("range"), NewGeneric(NameToType("Std::EndlessOpenRange", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::EndlessOpenRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::EndlessOpenRange::Iterator::Val", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("FS").(*Module)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Location").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("path"), NameToType("Std::FS::Path", env), NormalParameterKind, false), NewParameter(value.ToSymbol("span"), NameToType("Std::String::Span", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Returns the path of the file\nwhere the piece of text is located.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("file_path"), nil, nil, NameToType("Std::FS::Path", env), Never{})
					namespace.DefineMethod("Returns the position of the file fragment.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("span"), nil, nil, NameToType("Std::String::Span", env), Never{})

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("Path").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("path"), NameToType("Std::String", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Returns the last element of the path.\nTypically this is the name of the file.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("base"), nil, nil, NameToType("Std::String", env), Never{})
					namespace.DefineMethod("Returns a path based on `self` omitting the last element.\nTypically this would result in the path to the parent directory.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("dir"), nil, nil, NameToType("Std::FS::Path", env), Never{})
					namespace.DefineMethod("Returns the extension of the file.\nThe extension is the suffix beginning at the final dot in the final element of path; it is empty if there is no dot.\n\n```\nFS::Path(\"index\").extension #=> \"\"\nFS::Path(\"index.js\").extension #=> \".js\"\nFS::Path(\"index.html.erb\").extension #=> \".erb\"\n```", 0|METHOD_NATIVE_FLAG, value.ToSymbol("extension"), nil, nil, NameToType("Std::String", env), Never{})
					namespace.DefineMethod("Reports whether the path is absolute.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_absolute"), nil, nil, Bool{}, Never{})
					namespace.DefineMethod("Reports whether the path is local.\nIt is a primitive lexical check it does not take into account\nsymbolic links etc", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_local"), nil, nil, Bool{}, Never{})
					namespace.DefineMethod("Checks whether the path matches the given glob pattern.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("matches_glob"), nil, []*Parameter{NewParameter(value.ToSymbol("patten"), NameToType("Std::String", env), NormalParameterKind, false)}, Bool{}, Never{})
					namespace.DefineMethod("Returns a path based on `self` that is\nthe shortest possible version of it based on lexical analysis.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("normalize"), nil, nil, NameToType("Std::FS::Path", env), Never{})
					namespace.DefineMethod("Split the path into individual elements\nseparated by the OS separator (`/` or `\\`)", 0|METHOD_NATIVE_FLAG, value.ToSymbol("split"), nil, nil, NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::String", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
					namespace.DefineMethod("Returns the absolute version of this path.\nIf the path is not absolute it will be joined with the current working directory to turn it into an absolute path.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_absolute"), nil, nil, NameToType("Std::FS::Path", env), Never{})
					namespace.DefineMethod("Returns the path as a backslash `\"\\\\\"` separated `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_backslash_string"), nil, nil, NameToType("Std::String", env), Never{})
					namespace.DefineMethod("Returns a path based on `target` that is relative to `self`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_relative"), nil, []*Parameter{NewParameter(value.ToSymbol("target"), NameToType("Std::FS::Path", env), NormalParameterKind, false)}, NameToType("Std::FS::Path", env), Never{})
					namespace.DefineMethod("Returns the path as a slash `\"/\"` separated `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_slash_string"), nil, nil, NameToType("Std::String", env), Never{})
					namespace.DefineMethod("Returns the path as a `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})
					namespace.DefineMethod("VolumeName returns leading volume name.\nGiven \"C:\\foo\\bar\" it returns \"C:\" on Windows.\nGiven \"\\\\host\\share\\foo\" it returns \"\\\\host\\share\". On other platforms it returns \"\".", 0|METHOD_NATIVE_FLAG, value.ToSymbol("volume_name"), nil, nil, NameToType("Std::String", env), Never{})

					// Define constants

					// Define instance variables

					{
						namespace := namespace.Singleton()

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces

						// Define methods
						namespace.DefineMethod("Create a new path from a list of elements.\nUses the OS-specific separator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("build"), nil, []*Parameter{NewParameter(value.ToSymbol("elements"), NameToType("Std::String", env), PositionalRestParameterKind, false)}, NameToType("Std::FS::Path", env), Never{})
						namespace.DefineMethod("Create a new path from a slash-separated string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("from_slash"), nil, []*Parameter{NewParameter(value.ToSymbol("path"), NameToType("Std::String", env), NormalParameterKind, false)}, NameToType("Std::FS::Path", env), Never{})

						// Define constants

						// Define instance variables
					}
					{
						namespace := namespace.MustSubtype("Error").(*Class)

						namespace.Name() // noop - avoid unused variable error
						namespace.SetParent(NameToNamespace("Std::Error", env))

						// Include mixins and implement interfaces

						// Define methods

						// Define constants

						// Define instance variables
					}
				}
			}
			{
				namespace := namespace.MustSubtype("False").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Bool", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))

				// Define methods
				namespace.DefineMethod("Calculates a hash of the value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::FalseLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::FalseLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::FalseLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::FalseLiteralNode", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("FileSystemError").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Error", env))

				// Include mixins and implement interfaces

				// Define methods

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Float").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10\n\tvar b = 3\n\ta % b #=> 1\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Multiply this float by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Exponentiate this float, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Add `other` to this float.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1.2\n\t+a #=> 1.2\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Subtract `other` from this float.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns the result of negating the number.\n\n```\n\tvar a = 1.2\n\t-a #=> -1.2\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Divide this float by another float.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float", env), NormalParameterKind, false)}, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this float with another number.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` days.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("day"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` days.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("days"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Calculates a hash of the float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` hours.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hour"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` hours.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hours"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` microseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("microsecond"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` microseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("microseconds"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` milliseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("millisecond"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` milliseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("milliseconds"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` minutes.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("minute"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` minutes.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("minutes"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` nanoseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("nanosecond"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` nanoseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("nanoseconds"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` seconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("second"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` seconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("seconds"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::FloatLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::FloatLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::FloatLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::FloatLiteralNode", env), Never{})
				namespace.DefineMethod("Converts the float to a multi-precision floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_bigfloat"), nil, nil, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the float to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the float to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts the float to an automatically resized integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the float to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the float to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the float to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the float to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the Float to a String.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` weeks.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("week"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` weeks.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("weeks"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` years.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("year"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` years.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("years"), nil, nil, NameToType("Std::Duration", env), Never{})

				// Define constants
				namespace.DefineConstant(value.ToSymbol("INF"), NameToType("Std::Float", env))
				namespace.DefineConstant(value.ToSymbol("NAN"), NameToType("Std::Float", env))
				namespace.DefineConstant(value.ToSymbol("NEG_INF"), NameToType("Std::Float", env))

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a float.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Float32").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10\n\tvar b = 3\n\ta % b #=> 1\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Multiply this float by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Exponentiate this float, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Add `other` to this float.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1.2f32\n\t+a #=> 1.2f32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Subtract `other` from this float.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Returns the result of negating the number.\n\n```\n\tvar a = 1.2f32\n\t-a #=> -1.2f32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Divide this float by another float.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this float with another float.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Calculates a hash of the float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::Float32LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::Float32LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::Float32LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::Float32LiteralNode", env), Never{})
				namespace.DefineMethod("Converts the float to a multi-precision floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_bigfloat"), nil, nil, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Converts the float to a coercible floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the float to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts the float to an automatically resized integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the float to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the float to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the float to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the float to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a float32.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Float64").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10\n\tvar b = 3\n\ta % b #=> 1\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Multiply this float by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Exponentiate this float, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Add `other` to this float.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1.2f64\n\t+a #=> 1.2f64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Subtract `other` from this float.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Returns the result of negating the number.\n\n```\n\tvar a = 1.2f64\n\t-a #=> -1.2f64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Divide this float by another float.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this float with another float.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Float64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Calculates a hash of the float.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::Float64LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::Float64LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::Float64LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::Float64LiteralNode", env), Never{})
				namespace.DefineMethod("Converts the float to a multi-precision floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_bigfloat"), nil, nil, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Converts the float to a coercible floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the float to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts the float to an automatically resized integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the float to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the float to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the float to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the float to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the float to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a float64.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("FormatError").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Error", env))

				// Include mixins and implement interfaces

				// Define methods

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Generator").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Err"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Err"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Err"), NoValue{})
				typeParam.Default = Never{}

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Generator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::Generator::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

				// Define methods
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::Generator::Val", env), NewUnion(NewSymbolLiteral("stop_iteration"), NameToType("Std::Generator::Err", env)))
				namespace.DefineMethod("Resets the state of the generator.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("GlobError").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Error", env))

				// Include mixins and implement interfaces

				// Define methods

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("HashMap").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Key"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Key"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Key"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Value"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Value"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Value"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Map", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Value", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})))

				// Define methods
				namespace.DefineMethod("Create a new `HashMap` containing the pairs of `self`\nand another given record/map.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::Record", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::HashMap", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewUnion(NameToType("Std::HashMap::Key", env), NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewUnion(NameToType("Std::HashMap::Value", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), Never{})
				namespace.DefineMethod("Check whether the given value is a `HashMap`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check whether the given value is an `HashMap` or `HashRecord`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Get the element under the given key.\nReturns `nil` when the key is not present.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::HashMap::Key", env), NormalParameterKind, false)}, NewNilable(NameToType("Std::HashMap::Value", env)), Never{})
				namespace.DefineMethod("Set the element under the given index to the given value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]="), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::HashMap::Key", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::HashMap::Value", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Returns the number of key-value pairs that can be\nheld by the underlying array.\n\nThis value will change when the map gets resized,\nand the underlying array gets reallocated.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Check whether the given `pair` is present in this map.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check whether the given `key` is present in this map.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains_key"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::HashMap::Key", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check whether the given `value` is present in this map.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains_value"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::HashMap::Value", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Mutates the map.\n\nReallocates the underlying array to hold\nthe given number of new elements.\n\nExpands the `capacity` of the list by `new_slots`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("grow"), nil, []*Parameter{NewParameter(value.ToSymbol("new_slots"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the map.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::HashMap::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Value", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), Never{})
				namespace.DefineMethod("Returns the number of left slots for new key-value pairs\nin the underlying array.\nIt tells you how many more elements can be\nadded to the map before the underlying array gets\nreallocated.\n\nIt is always equal to `capacity - length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left_capacity"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of key-value pairs present in the map.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over the key value pairs of this map,\nyielding them to the given closure.\n\nReturns a new ArrayList that consists of the values returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the key value pairs of this map,\nyielding them to the given closure.\n\nReturns a new HashMap that consists of the key value pairs returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_pairs"), []*TypeParameter{NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("pair"), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::HashMap", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the values of this map,\nyielding them to the given closure.\n\nReturns a new HashMap that consists of the values returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_values"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::HashMap::Value", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::HashMap", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Key", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the values of this map,\nyielding them to the given closure.\n\nMutates the map in place replacing the values with the ones\nreturned by the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_values_mut"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values_mut", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::HashMap::Value", env), NormalParameterKind, false)}, NameToType("Std::HashMap::Value", env), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values_mut", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values_mut", true), Never{}, Any{}, nil, INVARIANT))

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 2)

					typeParam = NewTypeParameter(value.ToSymbol("Key"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Key"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Key"), NoValue{})

					typeParam = NewTypeParameter(value.ToSymbol("Value"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[1] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Value"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Value"), NoValue{})

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Iterator::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Iterator::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("map"), NewGeneric(NameToType("Std::HashMap", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Iterator::Key", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Iterator::Value", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next pair of the map.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashMap::Iterator::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashMap::Iterator::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("HashRecord").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Key"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Key"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Key"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Value"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Value"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Value"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Record", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Value", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})))

				// Define methods
				namespace.DefineMethod("Create a new `HashRecord` containing the pairs of `self`\nand another given record.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::Record", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::HashRecord", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewUnion(NameToType("Std::HashRecord::Key", env), NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewUnion(NameToType("Std::HashRecord::Value", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), Never{})
				namespace.DefineMethod("Check whether the given value is a `HashRecord`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check whether the given value is an `HashRecord` or `HashMap`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Get the element under the given key.\nReturns `nil` when the key is not present.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::HashRecord::Key", env), NormalParameterKind, false)}, NewNilable(NameToType("Std::HashRecord::Value", env)), Never{})
				namespace.DefineMethod("Check whether the given `pair` is present in this record.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check whether the given `key` is present in this record.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains_key"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::HashRecord::Key", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check whether the given `value` is present in this record.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains_value"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::HashRecord::Value", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the record.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::HashRecord::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Value", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), Never{})
				namespace.DefineMethod("Returns the number of key-value pairs present in the record.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over the key value pairs of this record,\nyielding them to the given closure.\n\nReturns a new ArrayList that consists of the values returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the key value pairs of this record,\nyielding them to the given closure.\n\nReturns a new HashRecord that consists of the key value pairs returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_pairs"), []*TypeParameter{NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("pair"), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::HashRecord", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the values of this record,\nyielding them to the given closure.\n\nReturns a new HashRecord that consists of the values returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_values"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::HashRecord::Value", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::HashRecord", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Key", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT))

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 2)

					typeParam = NewTypeParameter(value.ToSymbol("Key"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Key"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Key"), NoValue{})

					typeParam = NewTypeParameter(value.ToSymbol("Value"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[1] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Value"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Value"), NoValue{})

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Iterator::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Iterator::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("record"), NewGeneric(NameToType("Std::HashRecord", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Iterator::Key", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Iterator::Value", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next pair of the record.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::HashRecord::Iterator::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::HashRecord::Iterator::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("HashSet").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Set", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::HashSet::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

				// Define methods
				namespace.DefineMethod("Return the intersection of both sets.\nCreate a new `HashSet` containing only the elements\npresent both in `self` and `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"&\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"&\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::HashSet", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(Never{}, INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Return the union of both sets.\n\nCreate a new `HashSet` containing all the elements\npresent in `self` and `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::HashSet", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::HashSet::Val", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Adds the given value to the set.\n\nDoes nothing if the value is already present in the set.\n\nReallocates the underlying array if it is\ntoo small to hold it.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::HashSet::Val", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Check whether the given value is a `HashSet`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check whether the given value is a `HashSet`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Adds the given values to the set.\n\nSkips a value if it is already present in the set.\n\nReallocates the underlying array if it is\ntoo small to hold them.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("append"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::HashSet::Val", env), PositionalRestParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Returns the number of elements that can be\nheld by the underlying array.\n\nThis value will change when the set gets resized,\nand the underlying array gets reallocated.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Removes all elements from the set.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("clear"), nil, nil, Void{}, Never{})
				namespace.DefineMethod("Check whether the given `value` is present in this set.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::HashSet::Val", env), NameToType("Std::HashSet::Val", env), NameToType("Std::HashSet::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::HashSet::Val", env), NameToType("Std::HashSet::Val", env), NameToType("Std::HashSet::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Mutates the set.\n\nReallocates the underlying array to hold\nthe given number of new elements.\n\nExpands the `capacity` of the list by `new_slots`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("grow"), nil, []*Parameter{NewParameter(value.ToSymbol("new_slots"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the set.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::HashSet::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::HashSet::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Returns the number of left slots for new elements\nin the underlying array.\nIt tells you how many more elements can be\nadded to the set before the underlying array gets\nreallocated.\n\nIt is always equal to `capacity - length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left_capacity"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of elements present in the set.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over the elements of this set,\nyielding them to the given closure.\n\nReturns a new HashSet that consists of the elements returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::HashSet::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::HashSet", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Adds the given value to the set.\n\nReturns `false` if the value is already present in the set.\nOtherwise returns `true`.\n\nReallocates the underlying array if it is\ntoo small to hold it.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("push"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::HashSet::Val", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Removes the element from the set.\n\nReturns `true` if the element has been removed,\notherwise returns `false.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("remove"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::HashSet::Val", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Return the union of both sets.\n\nCreate a new `HashSet` containing all the elements\npresent in `self` and `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("|"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::HashSet", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::HashSet::Val", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::HashSet::Iterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("set"), NewGeneric(NameToType("Std::HashSet", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::HashSet::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the set.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::HashSet::Iterator::Val", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Hashable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Returns a hash.", 0|METHOD_ABSTRACT_FLAG|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("ImmutableCollection").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Iterable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ImmutableCollection::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

				// Define methods
				namespace.DefineMethod("Returns the number of elements present in the collection.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over the elements of this collection,\nyielding them to the given closure.\n\nReturns a new collection that consists of the elements returned\nby the given closure.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::ImmutableCollection::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ImmutableCollection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Base").(*Mixin)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					ImplementInterface(namespace, NewGeneric(NameToType("Std::ImmutableCollection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ImmutableCollection::Base::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))
					IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterable::FiniteBase", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ImmutableCollection::Base::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::ImmutableCollection::Base::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ImmutableCollection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("ImmutableSet").(*Mixin)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::ImmutableCollection::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ImmutableSet::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

				// Define methods
				namespace.DefineMethod("Return the intersection of both sets.\n\nCreate a new set containing only the elements\npresent both in `self` and `other`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("&"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"&\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"&\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Return the union of both sets.\n\nCreate a new set containing all the elements\npresent in `self` and `other`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::ImmutableSet::Val", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Iterates over the elements of this set,\nyielding them to the given closure.\n\nReturns a new ImmutableSet that consists of the elements returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::ImmutableSet::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Return the union of both sets.\n\nCreate a new set containing all the elements\npresent in `self` and `other`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("|"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::ImmutableSet::Val", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Incrementable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("T"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("T"), typeParam)
				namespace.DefineConstant(value.ToSymbol("T"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Get the next value of this type, the successor of `self`", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::Incrementable::T", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Inspectable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Returns a human readable `String`\nrepresentation of this value\nfor debugging etc.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Int").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10\n\tvar b = 3\n\ta % b #=> 1\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Performs bitwise AND.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Multiply this integer by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Add `other` to this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1\n\t+a #=> 1\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Subtract `other` from this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1\n\t-a #=> -1\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::CoercibleNumeric", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Performs bitwise XOR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("^"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` days.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("day"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` days.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("days"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Calculates a hash of the int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` hours.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hour"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` hours.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hours"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns an iterator that\niterates over every integer from `0` to `self`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NameToType("Std::Int::Iterator", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` microseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("microsecond"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` microseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("microseconds"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` milliseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("millisecond"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` milliseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("milliseconds"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` minutes.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("minute"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` minutes.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("minutes"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` nanoseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("nanosecond"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` nanoseconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("nanoseconds"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` seconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("second"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` seconds.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("seconds"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("times"), nil, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("i"), NameToType("Std::Int", env), NormalParameterKind, false)}, Void{}, Never{}), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::IntLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::IntLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::IntLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::IntLiteralNode", env), Never{})
				namespace.DefineMethod("Converts the integer to a floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the integer to a string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` weeks.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("week"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` weeks.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("weeks"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` years.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("year"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the duration equivalent to `self` years.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("years"), nil, nil, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Performs bitwise OR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("|"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4 #=> -5\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("~"), nil, nil, NameToType("Std::Int", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to an integer.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Int", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::Int", env), NewSymbolLiteral("stop_iteration"))

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Int16").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10i16\n\tvar b = 3i16\n\ta % b #=> 1i16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Performs bitwise AND.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Multiply this integer by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Add `other` to this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1i16\n\t+a #=> 1i16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Subtract `other` from this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1i16\n\t-a #=> -1i16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted arithmetically to the left by `other` positions, or to the right if `other` is negative.\n\nPreserves the integer's sign bit.\n\n4i16  << 1  #=> 8i16\n4i16  << -1 #=> 2i16\n-4i16 << 1  #=> -8i16\n-4i16 << -1 #=> -2i16", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Returns an integer shifted logically to the left by `other` positions, or to the right if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i16  <<< 1  #=> 8i16\n4i16  <<< -1 #=> 2i16\n-4i16 <<< 1  #=> -8i16\n-4i16 <<< -1 #=> 32766i16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted arithmetically to the right by `other` positions, or to the left if `other` is negative.\n\nPreserves the integer's sign bit.\n\n```\n4i16  >> 1  #=> 2i16\n4i16  >> -1 #=> 8i16\n-4i16 >> 1  #=> -2i16\n-4i16 >> -1 #=> -8i16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Returns an integer shifted logically the the right by `other` positions, or to the left if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i16  >>> 1  #=> 2i16\n4i16  >>> -1 #=> 8i16\n-4i16 >>> 1  #=> 32766i16\n-4i16 >>> -1 #=> -8i16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Performs bitwise XOR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("^"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Calculates a hash of the value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::Int16LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::Int16LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::Int16LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::Int16LiteralNode", env), Never{})
				namespace.DefineMethod("Converts the integer to a floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts to an automatically resizable integer type.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise OR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("|"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int16", env), NormalParameterKind, false)}, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4i16 #=> -5i16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("~"), nil, nil, NameToType("Std::Int16", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a int16.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Int32").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10i32\n\tvar b = 3i32\n\ta % b #=> 1i32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Performs bitwise AND.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Multiply this integer by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Add `other` to this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1i32\n\t+a #=> 1i32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Subtract `other` from this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1i32\n\t-a #=> -1i32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted arithmetically to the left by `other` positions, or to the right if `other` is negative.\n\nPreserves the integer's sign bit.\n\n4i32  << 1  #=> 8i32\n4i32  << -1 #=> 2i32\n-4i32 << 1  #=> -8i32\n-4i32 << -1 #=> -2i32", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Returns an integer shifted logically to the left by `other` positions, or to the right if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i32  <<< 1  #=> 8i32\n4i32  <<< -1 #=> 2i32\n-4i32 <<< 1  #=> -8i32\n-4i32 <<< -1 #=> 2147483646i32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted arithmetically to the right by `other` positions, or to the left if `other` is negative.\n\nPreserves the integer's sign bit.\n\n```\n4i32  >> 1  #=> 2i32\n4i32  >> -1 #=> 8i32\n-4i32 >> 1  #=> -2i32\n-4i32 >> -1 #=> -8i32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Returns an integer shifted logically the the right by `other` positions, or to the left if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i32  >>> 1  #=> 2i32\n4i32  >>> -1 #=> 8i32\n-4i32 >>> 1  #=> 2147483646i32\n-4i32 >>> -1 #=> -8i32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Performs bitwise XOR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("^"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Calculates a hash of the value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::Int32LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::Int32LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::Int32LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::Int32LiteralNode", env), Never{})
				namespace.DefineMethod("Converts the integer to a floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts to an automatically resizable integer type.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise OR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("|"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int32", env), NormalParameterKind, false)}, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4i32 #=> -5i32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("~"), nil, nil, NameToType("Std::Int32", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a int32.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Int64").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10i64\n\tvar b = 3i64\n\ta % b #=> 1i64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Performs bitwise AND.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Multiply this integer by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Add `other` to this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1i64\n\t+a #=> 1i64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Subtract `other` from this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1i64\n\t-a #=> -1i64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted arithmetically to the left by `other` positions, or to the right if `other` is negative.\n\nPreserves the integer's sign bit.\n\n4i64  << 1  #=> 8i64\n4i64  << -1 #=> 2i64\n-4i64 << 1  #=> -8i64\n-4i64 << -1 #=> -2i64", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Returns an integer shifted logically to the left by `other` positions, or to the right if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i64  <<< 1  #=> 8i64\n4i64  <<< -1 #=> 2i64\n-4i64 <<< 1  #=> -8i64\n-4i64 <<< -1 #=> 9223372036854775806i64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted arithmetically to the right by `other` positions, or to the left if `other` is negative.\n\nPreserves the integer's sign bit.\n\n```\n4i64  >> 1  #=> 2i64\n4i64  >> -1 #=> 8i64\n-4i64 >> 1  #=> -2i64\n-4i64 >> -1 #=> -8i64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Returns an integer shifted logically the the right by `other` positions, or to the left if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i64  >>> 1  #=> 2i64\n4i64  >>> -1 #=> 8i64\n-4i64 >>> 1  #=> 9223372036854775806i64\n-4i64 >>> -1 #=> -8i64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Performs bitwise XOR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("^"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Calculates a hash of the value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::Int64LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::Int64LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::Int64LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::Int64LiteralNode", env), Never{})
				namespace.DefineMethod("Converts the integer to a floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts to an automatically resizable integer type.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise OR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("|"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int64", env), NormalParameterKind, false)}, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4i64 #=> -5i64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("~"), nil, nil, NameToType("Std::Int64", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a int64.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Int8").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10i8\n\tvar b = 3i8\n\ta % b #=> 1i8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Performs bitwise AND.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Multiply this integer by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Add `other` to this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1i8\n\t+a #=> 1i8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Subtract `other` from this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1i8\n\t-a #=> -1i8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted arithmetically to the left by `other` positions, or to the right if `other` is negative.\n\nPreserves the integer's sign bit.\n\n4i8  << 1  #=> 8i8\n4i8  << -1 #=> 2i8\n-4i8 << 1  #=> -8i8\n-4i8 << -1 #=> -2i8", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Returns an integer shifted logically to the left by `other` positions, or to the right if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i8  <<< 1  #=> 8i8\n4i8  <<< -1 #=> 2i8\n-4i8 <<< 1  #=> -8i8\n-4i8 <<< -1 #=> 126i8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted arithmetically to the right by `other` positions, or to the left if `other` is negative.\n\nPreserves the integer's sign bit.\n\n```\n4i8  >> 1  #=> 2i8\n4i8  >> -1 #=> 8i8\n-4i8 >> 1  #=> -2i8\n-4i8 >> -1 #=> -8i8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Returns an integer shifted logically the the right by `other` positions, or to the left if `other` is negative.\n\nUnlike an arithmetic shift, a logical shift does not preserve the integer's sign bit.\n\n```\n4i8  >>> 1  #=> 2i8\n4i8  >>> -1 #=> 8i8\n-4i8 >>> 1  #=> 126i8\n-4i8 >>> -1 #=> -8i8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Performs bitwise XOR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("^"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Calculates a hash of the value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::Int8LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::Int8LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::Int8LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::Int8LiteralNode", env), Never{})
				namespace.DefineMethod("Converts the integer to a floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts to an automatically resizable integer type.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise OR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("|"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Int8", env), NormalParameterKind, false)}, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4i8 #=> -5i8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("~"), nil, nil, NameToType("Std::Int8", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a int8.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Interface").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the name of the interface.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Iterable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Err"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Err"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Err"), NoValue{})
				typeParam.Default = Never{}

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NewGeneric(NameToType("Std::PrimitiveIterable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Iterable::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::Iterable::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))
				ImplementInterface(namespace, NameToType("Std::Inspectable", env).(*Interface))

				// Define methods
				namespace.DefineMethod("Checks whether any element of this iterable satisfies the given predicate.\n\nMay never return if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("any"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :any", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :any", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Bool{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :any", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Check whether the given `value` is present in this iterable.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Iterable::Val", env), NameToType("Std::Iterable::Val", env), NameToType("Std::Iterable::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Iterable::Val", env), NameToType("Std::Iterable::Val", env), NameToType("Std::Iterable::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Returns the number of elements matching the given predicate.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("count"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :count", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :count", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NameToType("Std::Int", env), NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :count", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Returns a new iterable containing all elements except first `n` elements.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("drop"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Returns a new iterable containing all elements except first elements that satisfy the given predicate.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("drop_while"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :drop_while", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :drop_while", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :drop_while", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Checks whether every element of this iterable satisfies the given predicate.\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("every"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :every", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :every", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Bool{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :every", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Returns a new iterable containing only elements matching the given predicate.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("filter"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :filter", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :filter", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :filter", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Returns the first element matching the given predicate.\nReturns `nil` otherwise.\n\nMay never return if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("find"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewNilable(NameToType("Std::Iterable::Val", env)), NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Returns the first element matching the given predicate.\nThrows an error otherwise.\n\nMay never return if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("find_err"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find_err", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find_err", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NameToType("Std::Iterable::Val", env), NewUnion(NameToType("Std::Iterable::NotFoundError", env), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find_err", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Returns the first element.\nThrows an unchecked error when the iterable is empty.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("first"), nil, nil, NameToType("Std::Iterable::Val", env), NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Returns the first element.\nThrows an error when the iterable is empty.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("first_err"), nil, nil, NameToType("Std::Iterable::Val", env), NewUnion(NameToType("Std::Iterable::NotFoundError", env), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Reduces the elements of this iterable to a single value by\niteratively combining each element with an initial value using the provided function.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("fold"), []*TypeParameter{NewTypeParameter(value.ToSymbol("I"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("initial"), NewTypeParameter(value.ToSymbol("I"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NormalParameterKind, false), NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("accum"), NewTypeParameter(value.ToSymbol("I"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NormalParameterKind, false), NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("I"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("I"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Returns the first index of element, or -1 if it could not be found.\n\nMay never return if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("index_of"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :index_of", true), NameToType("Std::Iterable::Val", env), NameToType("Std::Iterable::Val", env), NameToType("Std::Iterable::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("element"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :index_of", true), NameToType("Std::Iterable::Val", env), NameToType("Std::Iterable::Val", env), NameToType("Std::Iterable::Val", env), INVARIANT), NormalParameterKind, false)}, NameToType("Std::Int", env), NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Checks whether the iterable is empty.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_empty"), nil, nil, Bool{}, NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Returns the first element.\nThrows an unchecked error when the iterable is empty.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("last"), nil, nil, NameToType("Std::Iterable::Val", env), NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Returns the last element.\nThrows an error when the iterable is empty.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("last_err"), nil, nil, NameToType("Std::Iterable::Val", env), NewUnion(NameToType("Std::Iterable::NotFoundError", env), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Returns the number of elements present in the iterable.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Iterates over the elements of this iterable,\nyielding them to the given closure.\n\nReturns a new iterable that consists of the elements returned\nby the given closure.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Iterable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})), NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Reduces the elements of this iterable to a single value by\niteratively combining them using the provided function.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("reduce"), []*TypeParameter{NewTypeParameter(value.ToSymbol("A"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("accum"), NewTypeParameter(value.ToSymbol("A"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT), NormalParameterKind, false), NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("A"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("A"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT), NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Returns a new iterable containing only elements not matching the given predicate.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("reject"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reject", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reject", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reject", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Returns a new iterable containing only the first `n` elements.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("take"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Returns a new iterable containing first elements satisfying the given predicate.\n\nMay never return if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("take_while"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :take_while", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :take_while", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :take_while", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::Err", env)))
				namespace.DefineMethod("Creates a new collection that contains the elements of this iterable.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_collection"), []*TypeParameter{NewTypeParameter(value.ToSymbol("T"), NewTypeParamNamespace("Type Parameter Container of :to_collection", true), NameToType("Std::Iterable::Val", env), Any{}, nil, INVARIANT)}, nil, NewGeneric(NameToType("Std::Collection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("T"), NewTypeParamNamespace("Type Parameter Container of :to_collection", true), NameToType("Std::Iterable::Val", env), Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Creates a new immutable collection that contains the elements of this iterable.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_immutable_collection"), nil, nil, NewGeneric(NameToType("Std::ImmutableCollection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Iterable::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Creates a new list that contains the elements of this iterable.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_list"), []*TypeParameter{NewTypeParameter(value.ToSymbol("T"), NewTypeParamNamespace("Type Parameter Container of :to_list", true), NameToType("Std::Iterable::Val", env), Any{}, nil, INVARIANT)}, nil, NewGeneric(NameToType("Std::List", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("T"), NewTypeParamNamespace("Type Parameter Container of :to_list", true), NameToType("Std::Iterable::Val", env), Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Creates a new tuple that contains the elements of this iterable.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_tuple"), nil, nil, NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::Iterable::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})), NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Returns the first element.\nReturns `nil` when the iterable is empty.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("try_first"), nil, nil, NewNilable(NameToType("Std::Iterable::Val", env)), NameToType("Std::Iterable::Err", env))
				namespace.DefineMethod("Returns the first element.\nReturns `nil` when the collection is empty.\n\nNever returns if the iterable is infinite.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("try_last"), nil, nil, NewNilable(NameToType("Std::Iterable::Val", env)), NameToType("Std::Iterable::Err", env))

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Base").(*Mixin)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 2)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

					typeParam = NewTypeParameter(value.ToSymbol("Err"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[1] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Err"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Err"), NoValue{})
					typeParam.Default = Never{}

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterable::FiniteBase", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Iterable::Base::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::Iterable::Base::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), NameToType("Std::Iterable::Base::Err", env))

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("FiniteBase").(*Mixin)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 2)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

					typeParam = NewTypeParameter(value.ToSymbol("Err"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[1] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Err"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Err"), NoValue{})
					typeParam.Default = Never{}

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					ImplementInterface(namespace, NewGeneric(NameToType("Std::Iterable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Iterable::FiniteBase::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::Iterable::FiniteBase::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("any"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :any", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :any", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Bool{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :any", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Iterable::FiniteBase::Val", env), NameToType("Std::Iterable::FiniteBase::Val", env), NameToType("Std::Iterable::FiniteBase::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Iterable::FiniteBase::Val", env), NameToType("Std::Iterable::FiniteBase::Val", env), NameToType("Std::Iterable::FiniteBase::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("count"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :count", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :count", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NameToType("Std::Int", env), NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :count", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("drop"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("drop_while"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :drop_while", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :drop_while", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :drop_while", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("every"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :every", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :every", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Bool{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :every", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("filter"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :filter", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :filter", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :filter", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("find"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewNilable(NameToType("Std::Iterable::FiniteBase::Val", env)), NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("find_err"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find_err", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find_err", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NameToType("Std::Iterable::FiniteBase::Val", env), NewUnion(NameToType("Std::Iterable::NotFoundError", env), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :find_err", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("first"), nil, nil, NameToType("Std::Iterable::FiniteBase::Val", env), NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("first_err"), nil, nil, NameToType("Std::Iterable::FiniteBase::Val", env), NewUnion(NameToType("Std::Iterable::NotFoundError", env), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("fold"), []*TypeParameter{NewTypeParameter(value.ToSymbol("I"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("initial"), NewTypeParameter(value.ToSymbol("I"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NormalParameterKind, false), NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("accum"), NewTypeParameter(value.ToSymbol("I"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NormalParameterKind, false), NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("I"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("I"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :fold", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("index_of"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :index_of", true), NameToType("Std::Iterable::FiniteBase::Val", env), NameToType("Std::Iterable::FiniteBase::Val", env), NameToType("Std::Iterable::FiniteBase::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("element"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :index_of", true), NameToType("Std::Iterable::FiniteBase::Val", env), NameToType("Std::Iterable::FiniteBase::Val", env), NameToType("Std::Iterable::FiniteBase::Val", env), INVARIANT), NormalParameterKind, false)}, NameToType("Std::Int", env), NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_empty"), nil, nil, Bool{}, NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("last"), nil, nil, NameToType("Std::Iterable::FiniteBase::Val", env), NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("last_err"), nil, nil, NameToType("Std::Iterable::FiniteBase::Val", env), NewUnion(NameToType("Std::Iterable::NotFoundError", env), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Iterable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})), NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reduce"), []*TypeParameter{NewTypeParameter(value.ToSymbol("A"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("accum"), NewTypeParameter(value.ToSymbol("A"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT), NormalParameterKind, false), NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("A"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("A"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT), NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reduce", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reject"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reject", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reject", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :reject", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("take"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("take_while"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :take_while", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Iterable::FiniteBase::Val", env), NormalParameterKind, false)}, Bool{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :take_while", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewUnion(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :take_while", true), Never{}, Any{}, nil, INVARIANT), NameToType("Std::Iterable::FiniteBase::Err", env)))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_collection"), []*TypeParameter{NewTypeParameter(value.ToSymbol("T"), NewTypeParamNamespace("Type Parameter Container of :to_collection", true), NameToType("Std::Iterable::FiniteBase::Val", env), Any{}, nil, INVARIANT)}, nil, NewGeneric(NameToType("Std::Collection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("T"), NewTypeParamNamespace("Type Parameter Container of :to_collection", true), NameToType("Std::Iterable::FiniteBase::Val", env), Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_immutable_collection"), nil, nil, NewGeneric(NameToType("Std::ImmutableCollection", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Iterable::FiniteBase::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_list"), []*TypeParameter{NewTypeParameter(value.ToSymbol("T"), NewTypeParamNamespace("Type Parameter Container of :to_list", true), NameToType("Std::Iterable::FiniteBase::Val", env), Any{}, nil, INVARIANT)}, nil, NewGeneric(NameToType("Std::List", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("T"), NewTypeParamNamespace("Type Parameter Container of :to_list", true), NameToType("Std::Iterable::FiniteBase::Val", env), Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_tuple"), nil, nil, NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::Iterable::FiniteBase::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})), NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("try_first"), nil, nil, NewNilable(NameToType("Std::Iterable::FiniteBase::Val", env)), NameToType("Std::Iterable::FiniteBase::Err", env))
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("try_last"), nil, nil, NewNilable(NameToType("Std::Iterable::FiniteBase::Val", env)), NameToType("Std::Iterable::FiniteBase::Err", env))

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("NotFoundError").(*Class)

					namespace.Name() // noop - avoid unused variable error
					namespace.SetParent(NameToNamespace("Std::Error", env))

					// Include mixins and implement interfaces

					// Define methods

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("IterableRange").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
				typeParam.UpperBound = NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::IterableRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::IterableRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Container", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::IterableRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::PrimitiveIterable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::IterableRange::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

				// Define methods
				namespace.DefineMethod("Returns the upper bound of the range.\nReturns `nil` if the range is endless.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("end"), nil, nil, NewNilable(NameToType("Std::IterableRange::Val", env)), Never{})
				namespace.DefineMethod("Returns `true` when the range is left-closed.\nOtherwise, the range is left-open.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_closed"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns `true` when the range is left-open.\nOtherwise, the range is left-closed.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_open"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns `true` when the range is right-closed.\nOtherwise, the range is right-open.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_closed"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns `true` when the range is right-open.\nOtherwise, the range is right-closed.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_open"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns the lower bound of the range.\nReturns `nil` if the Range is beginless.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("start"), nil, nil, NewNilable(NameToType("Std::IterableRange::Val", env)), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Iterator").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Err"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Err"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Err"), NoValue{})
				typeParam.Default = Never{}

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Iterable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Iterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::Iterator::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

				// Define methods
				namespace.DefineMethod("", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, Self{}, Never{})
				namespace.DefineMethod("Returns the next element.\nThrows `:stop_iteration` when no more elements are available.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::Iterator::Val", env), NewUnion(NewSymbolLiteral("stop_iteration"), NameToType("Std::Iterator::Err", env)))

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Base").(*Mixin)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 2)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

					typeParam = NewTypeParameter(value.ToSymbol("Err"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[1] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Err"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Err"), NoValue{})
					typeParam.Default = Never{}

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					ImplementInterface(namespace, NewGeneric(NameToType("Std::Iterator", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Iterator::Base::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::Iterator::Base::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))
					IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterable::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Iterator::Base::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::Iterator::Base::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, Self{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Kernel").(*Module)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Converts the values to `String`\nand prints them to stdout.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("print"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::String::Convertible", env), PositionalRestParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Converts the values to `String`\nand prints them to stdout with a newline.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("println"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::String::Convertible", env), PositionalRestParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Converts the values to `String`\nand prints them to stdout with a newline.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("puts"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::String::Convertible", env), PositionalRestParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Pauses the execution of the current thread for the amount\nof time represented by the passed `Duration`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("sleep"), nil, []*Parameter{NewParameter(value.ToSymbol("dur"), NameToType("Std::Duration", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("The asynchronous version of `sleep`.\nReturns a promise that gets resolved after the given amount of time.", 0|METHOD_NATIVE_FLAG|METHOD_ASYNC_FLAG, value.ToSymbol("timeout"), nil, []*Parameter{NewParameter(value.ToSymbol("dur"), NameToType("Std::Duration", env), NormalParameterKind, false), NewParameter(value.ToSymbol("_pool"), NameToType("Std::ThreadPool", env), DefaultValueParameterKind, false)}, NewGeneric(NameToType("Std::Promise", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(Void{}, COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("LeftOpenRange").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
				typeParam.UpperBound = NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::LeftOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Range", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::LeftOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))

				// extend where Val < Std::Incrementable[Val] & Std::Comparable[Val]
				mixin = NewMixin("", false, "", env)
				{
					namespace := mixin
					namespace.Name() // noop - avoid unused variable error
					namespace.DefineSubtype(value.ToSymbol("Val"), NewTypeParameter(value.ToSymbol("Val"), NameToType("Std::LeftOpenRange", env).(*Class), Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::LeftOpenRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::LeftOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT))

					// Define methods
					namespace.DefineMethod("Returns the iterator for this range.\nOnly ranges of incrementable values can be iterated over.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::LeftOpenRange::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::LeftOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				}
				IncludeMixinWithWhere(namespace, mixin, []*TypeParameter{NewTypeParameter(value.ToSymbol("Val"), mixin, Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::LeftOpenRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::LeftOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT)})

				// Define methods
				namespace.DefineMethod("Check whether the given `value` is present in this range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::LeftOpenRange::Val", env), NameToType("Std::LeftOpenRange::Val", env), NameToType("Std::LeftOpenRange::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::LeftOpenRange::Val", env), NameToType("Std::LeftOpenRange::Val", env), NameToType("Std::LeftOpenRange::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Returns the upper bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("end"), nil, nil, NameToType("Std::LeftOpenRange::Val", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_closed"), nil, nil, False{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_closed"), nil, nil, True{}, Never{})
				namespace.DefineMethod("Returns the lower bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start"), nil, nil, NameToType("Std::LeftOpenRange::Val", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
					typeParam.UpperBound = NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::LeftOpenRange::Iterator::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::LeftOpenRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::LeftOpenRange::Iterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("range"), NewGeneric(NameToType("Std::LeftOpenRange", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::LeftOpenRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::LeftOpenRange::Iterator::Val", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("List").(*Mixin)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::List::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Collection::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::List::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

				// Define methods
				namespace.DefineMethod("Create a new `List` containing the elements of `self`\nand another given `Tuple`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::List", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::List::Val", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Set the element under the given index to the given value.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("[]="), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::List::Val", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Iterates over the elements of this list,\nyielding them to the given closure.\n\nReturns a new list that consists of the elements returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::List::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::List", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the elements of this collection,\nyielding them to the given closure.\n\nMutates the collection in place replacing the elements with the ones\nreturned by the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_mut"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::List::Val", env), NormalParameterKind, false)}, NameToType("Std::List::Val", env), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT))

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Lockable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("lock"), nil, nil, Void{}, Never{})
				namespace.DefineMethod("", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("unlock"), nil, nil, Void{}, Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Map").(*Mixin)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Key"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Key"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Key"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Value"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Value"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Value"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Record", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Map::Key", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Map::Value", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})))

				// Define methods
				namespace.DefineMethod("Create a new map containing the pairs of `self`\nand another given record/map.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::Record", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Map", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewUnion(NameToType("Std::Map::Key", env), NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewUnion(NameToType("Std::Map::Value", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), Never{})
				namespace.DefineMethod("Set the element under the given index to the given value.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("[]="), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::Map::Key", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::Map::Value", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Iterates over the key value pairs of this map,\nyielding them to the given closure.\n\nReturns a new map that consists of the key value pairs returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_pairs"), []*TypeParameter{NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("pair"), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Map::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Map::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Map", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the values of this map,\nyielding them to the given closure.\n\nMutates the map in place replacing the values with the ones\nreturned by the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_values_mut"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values_mut", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Map::Value", env), NormalParameterKind, false)}, NameToType("Std::Map::Value", env), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values_mut", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values_mut", true), Never{}, Any{}, nil, INVARIANT))

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Mixin").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the name of the mixin.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Module").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the name of the module.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Nil").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))

				// Define methods
				namespace.DefineMethod("Calculates a hash of the value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::NilLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::NilLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::NilLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::NilLiteralNode", env), Never{})
				namespace.DefineMethod("Converts `nil` to `BigFloat`.\nAlways returns `0.0bf`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_big_float"), nil, nil, NameToType("Std::BigFloat", env), Never{})
				namespace.DefineMethod("Converts `nil` to `Char`.\nAlways returns a null char \\x00.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_char"), nil, nil, NameToType("Std::Char", env), Never{})
				namespace.DefineMethod("Converts `nil` to `Float`.\nAlways returns `0.0`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts `nil` to `Float32`.\nAlways returns `0.0f32`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts `nil` to `Float64`.\nAlways returns `0.0f64`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts `nil` to `Int`.\nAlways returns `0`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts `nil` to `Int16`.\nAlways returns `0i16`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts `nil` to `Int32`.\nAlways returns `0i32`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts `nil` to `Int64`.\nAlways returns `0i64`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts `nil` to `Int8`.\nAlways returns `0i8`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts `nil` to `String`.\nAlways returns an empty string `\"\"`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Converts `nil` to `UInt16`.\nAlways returns `0u16`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts `nil` to `UInt32`.\nAlways returns `0u32`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts `nil` to `UInt64`.\nAlways returns `0u64`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts `nil` to `UInt8`.\nAlways returns `0u8`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Object").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces

				// Define methods

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("OpenRange").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
				typeParam.UpperBound = NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::OpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Range", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::OpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))

				// extend where Val < Std::Incrementable[Val] & Std::Comparable[Val]
				mixin = NewMixin("", false, "", env)
				{
					namespace := mixin
					namespace.Name() // noop - avoid unused variable error
					namespace.DefineSubtype(value.ToSymbol("Val"), NewTypeParameter(value.ToSymbol("Val"), NameToType("Std::OpenRange", env).(*Class), Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::OpenRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::OpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT))

					// Define methods
					namespace.DefineMethod("Returns the iterator for this range.\nOnly ranges of incrementable values can be iterated over.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::OpenRange::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::OpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				}
				IncludeMixinWithWhere(namespace, mixin, []*TypeParameter{NewTypeParameter(value.ToSymbol("Val"), mixin, Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::OpenRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::OpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT)})

				// Define methods
				namespace.DefineMethod("Check whether the given `value` is present in this range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::OpenRange::Val", env), NameToType("Std::OpenRange::Val", env), NameToType("Std::OpenRange::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::OpenRange::Val", env), NameToType("Std::OpenRange::Val", env), NameToType("Std::OpenRange::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Returns the upper bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("end"), nil, nil, NameToType("Std::OpenRange::Val", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_closed"), nil, nil, False{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_closed"), nil, nil, False{}, Never{})
				namespace.DefineMethod("Returns the lower bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start"), nil, nil, NameToType("Std::OpenRange::Val", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
					typeParam.UpperBound = NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::OpenRange::Iterator::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::OpenRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::OpenRange::Iterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("range"), NewGeneric(NameToType("Std::OpenRange", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::OpenRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::OpenRange::Iterator::Val", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("OutOfRangeError").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Error", env))

				// Include mixins and implement interfaces

				// Define methods

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Pair").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Key"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Key"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Key"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Value"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Value"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Value"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NewUnion(NameToType("Std::Pair::Key", env), NameToType("Std::Pair::Value", env)), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))

				// Define methods
				namespace.DefineMethod("Instantiate the `Pair` with the given key and value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::Pair::Key", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::Pair::Value", env), NormalParameterKind, false)}, Void{}, Never{})
				namespace.DefineMethod("Check whether the given value\nis a `Pair` that is equal to this `Pair`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Get the element with the given index.\nThe key is `0`, value is `1`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NewUnion(NameToType("Std::Pair::Key", env), NameToType("Std::Pair::Value", env)), Never{})
				namespace.DefineMethod("Returns an iterator that iterates\nover each element of the pair.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::Pair::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Pair::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Pair::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), Never{})
				namespace.DefineMethod("Returns the key, the first element of the tuple.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("key"), nil, nil, NameToType("Std::Pair::Key", env), Never{})
				namespace.DefineMethod("Always returns `2`.\nFor compatibility with `Tuple`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the value, the second element of the tuple.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("value"), nil, nil, NameToType("Std::Pair::Value", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 2)

					typeParam = NewTypeParameter(value.ToSymbol("Key"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Key"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Key"), NoValue{})

					typeParam = NewTypeParameter(value.ToSymbol("Value"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[1] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Value"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Value"), NoValue{})

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::Pair::Iterator::Key", env), NameToType("Std::Pair::Iterator::Value", env)), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("pair"), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Pair::Iterator::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Pair::Iterator::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the pair.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NewUnion(NameToType("Std::Pair::Iterator::Key", env), NameToType("Std::Pair::Iterator::Value", env)), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("PrimitiveIterable").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Err"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Err"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Err"), NoValue{})
				typeParam.Default = Never{}

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Returns an iterator for this structure.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::Iterator", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::PrimitiveIterable::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::PrimitiveIterable::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Promise").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Err"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Err"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Err"), NoValue{})
				typeParam.Default = Never{}

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Blocks the current thread until the promise\ngets resolved.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("await_sync"), nil, nil, NameToType("Std::Promise::Val", env), NameToType("Std::Promise::Err", env))
				namespace.DefineMethod("Check whether the promise is done.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_resolved"), nil, nil, Bool{}, Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.Singleton()

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Creates a new promise that is immediately rejected with the given error.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("rejected"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :rejected", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("err"), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :rejected", true), Never{}, Any{}, nil, INVARIANT), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Promise", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(Never{}, COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :rejected", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})), Never{})
					namespace.DefineMethod("Creates a new promise that is immediately resolved with the given result.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("resolved"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :resolved", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("result"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :resolved", true), Never{}, Any{}, nil, INVARIANT), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Promise", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :resolved", true), Never{}, Any{}, nil, INVARIANT), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})), Never{})
					namespace.DefineMethod("Returns a new promise that gets resolved when all given promises are resolved.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("wait"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :wait", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :wait", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("promises"), NewGeneric(NameToType("Std::Promise", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :wait", true), Never{}, Any{}, nil, INVARIANT), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :wait", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})), PositionalRestParameterKind, false)}, NewGeneric(NameToType("Std::Promise", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :wait", true), Never{}, Any{}, nil, INVARIANT), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :wait", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Range").(*Mixin)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Element"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Element"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Element"), NoValue{})
				typeParam.UpperBound = NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::Range::Element", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Container", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Range::Element", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

				// Define methods
				namespace.DefineMethod("Returns the upper bound of the range.\nReturns `nil` if the range is endless.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("end"), nil, nil, NewNilable(NameToType("Std::Range::Element", env)), Never{})
				namespace.DefineMethod("Returns `true` when the range is left-closed.\nOtherwise, the range is left-open.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_closed"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns `true` when the range is left-open.\nOtherwise, the range is left-closed.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_open"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns `true` when the range is right-closed.\nOtherwise, the range is right-open.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_closed"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns `true` when the range is right-open.\nOtherwise, the range is right-closed.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_open"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns the lower bound of the range.\nReturns `nil` if the Range is beginless.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start"), nil, nil, NewNilable(NameToType("Std::Range::Element", env)), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Record").(*Mixin)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Key"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Key"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Key"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Value"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Value"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Value"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterable::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Record::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Record::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

				// Define methods
				namespace.DefineMethod("Create a new record containing the pairs of `self`\nand another given record.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::Record", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Record", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewUnion(NameToType("Std::Record::Key", env), NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewUnion(NameToType("Std::Record::Value", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), Never{})
				namespace.DefineMethod("Check whether the given value is the same type of record\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check whether the given value is a record\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Get the element under the given key.\nReturns `nil` when the key is not present.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::Record::Key", env), NormalParameterKind, false)}, NewNilable(NameToType("Std::Record::Value", env)), Never{})
				namespace.DefineMethod("Check whether the given `pair` is present in this record.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Record::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Record::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Record::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Record::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Record::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Record::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Record::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Record::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Record::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Record::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Record::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Record::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check whether the given `key` is present in this record.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains_key"), nil, []*Parameter{NewParameter(value.ToSymbol("key"), NameToType("Std::Record::Key", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Check whether the given `value` is present in this record.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains_value"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Record::Value", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Iterates over the key value pairs of this record,\nyielding them to the given closure.\n\nReturns a new list that consists of the values returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Record::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Record::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::List", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the key value pairs of this record,\nyielding them to the given closure.\n\nReturns a new record that consists of the key value pairs returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_pairs"), []*TypeParameter{NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("pair"), NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Record::Key", env), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NameToType("Std::Record::Value", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Pair", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), COVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Record", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NewTypeParameter(value.ToSymbol("K"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_pairs", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Iterates over the values of this record,\nyielding them to the given closure.\n\nReturns a new record that consists of the values returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_values"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Record::Value", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Record", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Key"): NewTypeArgument(NameToType("Std::Record::Key", env), INVARIANT), value.ToSymbol("Value"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Key"), value.ToSymbol("Value")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_values", true), Never{}, Any{}, nil, INVARIANT))

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Regex").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Creates a new `Regex` that contains the\npattern of `self` repeated `n` times.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Create a new regex that contains\nthe patterns present in both operands.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Regex", env), NormalParameterKind, false)}, NameToType("Std::Regex", env), Never{})
				namespace.DefineMethod("Check whether the pattern matches\nthe given string.\n\nReturns `true` if it matches, otherwise `false`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("matches"), nil, []*Parameter{NewParameter(value.ToSymbol("str"), NameToType("Std::String", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::UninterpolatedRegexLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::UninterpolatedRegexLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::UninterpolatedRegexLiteralNode", env), Never{})
				namespace.DefineMethod("Creates a new string with this character.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("ResettableIterator").(*Interface)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 2)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				typeParam = NewTypeParameter(value.ToSymbol("Err"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[1] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Err"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Err"), NoValue{})
				typeParam.Default = Never{}

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Iterator", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ResettableIterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::ResettableIterator::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

				// Define methods
				namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Base").(*Mixin)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 2)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

					typeParam = NewTypeParameter(value.ToSymbol("Err"), namespace, Never{}, Any{}, nil, COVARIANT)
					typeParams[1] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Err"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Err"), NoValue{})
					typeParam.Default = Never{}

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					ImplementInterface(namespace, NewGeneric(NameToType("Std::ResettableIterator", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ResettableIterator::Base::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::ResettableIterator::Base::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))
					IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::ResettableIterator::Base::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(NameToType("Std::ResettableIterator::Base::Err", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("RightOpenRange").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
				typeParam.UpperBound = NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::RightOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Range", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NameToType("Std::RightOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Element")})))

				// extend where Val < Std::Incrementable[Val] & Std::Comparable[Val]
				mixin = NewMixin("", false, "", env)
				{
					namespace := mixin
					namespace.Name() // noop - avoid unused variable error
					namespace.DefineSubtype(value.ToSymbol("Val"), NewTypeParameter(value.ToSymbol("Val"), NameToType("Std::RightOpenRange", env).(*Class), Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::RightOpenRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::RightOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT))

					// Define methods
					namespace.DefineMethod("Returns the iterator for this range.\nOnly ranges of incrementable values can be iterated over.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NewGeneric(NameToType("Std::RightOpenRange::Iterator", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::RightOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				}
				IncludeMixinWithWhere(namespace, mixin, []*TypeParameter{NewTypeParameter(value.ToSymbol("Val"), mixin, Never{}, NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::RightOpenRange::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::RightOpenRange::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")}))), nil, INVARIANT)})

				// Define methods
				namespace.DefineMethod("Check whether the given `value` is present in this range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::RightOpenRange::Val", env), NameToType("Std::RightOpenRange::Val", env), NameToType("Std::RightOpenRange::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::RightOpenRange::Val", env), NameToType("Std::RightOpenRange::Val", env), NameToType("Std::RightOpenRange::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Returns the upper bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("end"), nil, nil, NameToType("Std::RightOpenRange::Val", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_left_closed"), nil, nil, True{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_right_closed"), nil, nil, False{}, Never{})
				namespace.DefineMethod("Returns the lower bound of the range.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start"), nil, nil, NameToType("Std::RightOpenRange::Val", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Set up type parameters
					var typeParam *TypeParameter
					typeParams := make([]*TypeParameter, 1)

					typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
					typeParams[0] = typeParam
					namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
					namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})
					typeParam.UpperBound = NewIntersection(NewGeneric(NameToType("Std::Incrementable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::RightOpenRange::Iterator::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("T")})), NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(NameToType("Std::RightOpenRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

					namespace.SetTypeParameters(typeParams)

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::RightOpenRange::Iterator::Val", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("range"), NewGeneric(NameToType("Std::RightOpenRange", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::RightOpenRange::Iterator::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::RightOpenRange::Iterator::Val", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Set").(*Mixin)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Val"), namespace, Never{}, Any{}, nil, INVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Val"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Val"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Set::Val", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))
				IncludeMixin(namespace, NewGeneric(NameToType("Std::Collection::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Set::Val", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

				// Define methods
				namespace.DefineMethod("Return the intersection of both sets.\n\nCreate a new set containing only the elements\npresent both in `self` and `other`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("&"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"&\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"&\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Set", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(Never{}, INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Return the union of both sets.\n\nCreate a new set containing all the elements\npresent in `self` and `other`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Set", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::Set::Val", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
				namespace.DefineMethod("Adds the given value to the set.\n\nDoes nothing if the value is already present in the set.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Set::Val", env), NormalParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Adds the given values to the set.\n\nSkips a value if it is already present in the set.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("append"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::Set::Val", env), PositionalRestParameterKind, false)}, Self{}, Never{})
				namespace.DefineMethod("Check whether the given `value` is present in this set.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Set::Val", env), NameToType("Std::Set::Val", env), NameToType("Std::Set::Val", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Set::Val", env), NameToType("Std::Set::Val", env), NameToType("Std::Set::Val", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Iterates over the elements of this set,\nyielding them to the given closure.\n\nReturns a new Set that consists of the elements returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Set::Val", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Set", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Adds the given value to the set.\n\nReturns `false` if the value is already present in the set.\nOtherwise returns `true`.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("push"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Set::Val", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Return the union of both sets.\n\nCreate a new set containing all the elements\npresent in `self` and `other`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("|"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::ImmutableSet", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Set", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::Set::Val", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("StackTrace").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::ImmutableCollection::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::CallFrame", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

				// Define methods
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::CallFrame", env), Never{})
				namespace.DefineMethod("Returns an iterator that iterates\nover each call frame of the stack trace.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NameToType("Std::StackTrace::Iterator", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the string representation of the stack trace.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Iterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::CallFrame", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("stack_trace"), NameToType("Std::StackTrace", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::CallFrame", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("String").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Concatenate this `String`\nwith another `String` or `Char`.\n\nCreates a new `String` containing the content\nof both operands.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Remove the given suffix from the `String`.\n\nDoes nothing if the `String` doesn't end\nwith `suffix` and returns `self`.\n\nIf the `String` ends with the given suffix\na new `String` gets created and returned that doesn't contain\nthe suffix.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("suffix"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Get the byte with the given index.\nIndices start at 0.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("byte_at"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Get the number of bytes that this\nstring contains.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("byte_count"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over all bytes of a `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("byte_iter"), nil, nil, NameToType("Std::String::ByteIterator", env), Never{})
				namespace.DefineMethod("Get the number of Unicode code points\nthat this `String` contains.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("char_count"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over all unicode code points of a `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("char_iter"), nil, nil, NameToType("Std::String::CharIterator", env), Never{})
				namespace.DefineMethod("Get the Unicode code point with the given index.\nIndices start at 0.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("chat_at"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Char", env), Never{})
				namespace.DefineMethod("Concatenate this `String`\nwith another `String` or `Char`.\n\nCreates a new `String` containing the content\nof both operands.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("concat"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Get the Unicode grapheme cluster with the given index.\nIndices start at 0.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("grapheme_at"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Get the number of unicode grapheme clusters\npresent in this string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("grapheme_count"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Iterates over all grapheme clusters of a `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("grapheme_iter"), nil, nil, NameToType("Std::String::GraphemeIterator", env), Never{})
				namespace.DefineMethod("Calculates a hash of the string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Check whether the `String` is empty.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_empty"), nil, nil, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Iterates over all unicode code points of a `String`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NameToType("Std::String::CharIterator", env), Never{})
				namespace.DefineMethod("Get the number of Unicode code points\nthat this `String` contains.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Create a new string with all of the characters\nof this one turned into lowercase.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("lowercase"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Remove the given suffix from the `String`.\n\nDoes nothing if the `String` doesn't end\nwith `suffix` and returns `self`.\n\nIf the `String` ends with the given suffix\na new `String` gets created and returned that doesn't contain\nthe suffix.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("remove_suffix"), nil, []*Parameter{NewParameter(value.ToSymbol("suffix"), NewUnion(NameToType("Std::String", env), NameToType("Std::Char", env)), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Creates a new `String` that contains the\ncontent of `self` repeated `n` times.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("repeat"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::DoubleQuotedStringLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::DoubleQuotedStringLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::DoubleQuotedStringLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::DoubleQuotedStringLiteralNode", env), Never{})
				namespace.DefineMethod("Convert the `String` to an `Int`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), NameToType("Std::FormatError", env))
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Convert the `String` to a `Symbol`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_symbol"), nil, nil, NameToType("Std::Symbol", env), Never{})
				namespace.DefineMethod("Create a new string with all of the characters\nof this one turned into uppercase.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("uppercase"), nil, nil, NameToType("Std::String", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("ByteIterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::UInt8", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("string"), NameToType("Std::String", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NameToType("Std::String::ByteIterator", env), Never{})
					namespace.DefineMethod("Get the next byte.\nThrows `:stop_iteration` when no more bytes are available.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::UInt8", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("CharIterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Char", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("string"), NameToType("Std::String", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NameToType("Std::String::CharIterator", env), Never{})
					namespace.DefineMethod("Get the next character.\nThrows `:stop_iteration` when no more characters are available.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::Char", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a string.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("GraphemeIterator").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::Iterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::String", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("string"), NameToType("Std::String", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NameToType("Std::String::GraphemeIterator", env), Never{})
					namespace.DefineMethod("Get the next grapheme.\nThrows `:stop_iteration` when no more graphemes are available.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::String", env), NewSymbolLiteral("stop_iteration"))
					namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("Position").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("byte_offset"), NameToType("Std::Int", env), NormalParameterKind, false), NewParameter(value.ToSymbol("line"), NameToType("Std::Int", env), NormalParameterKind, false), NewParameter(value.ToSymbol("column"), NameToType("Std::Int", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Returns the offset of the first byte of the character in a string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("byte_offset"), nil, nil, NameToType("Std::Int", env), Never{})
					namespace.DefineMethod("Returns the number of the\ncolumn where the character is located.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("column"), nil, nil, NameToType("Std::Int", env), Never{})
					namespace.DefineMethod("Returns the number of the line\nwhere the character is located in a string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("line"), nil, nil, NameToType("Std::Int", env), Never{})

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("Span").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("start_pos"), NameToType("Std::String::Position", env), NormalParameterKind, false), NewParameter(value.ToSymbol("end_pos"), NameToType("Std::String::Position", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Returns the position of the last character\nof a piece of text in a string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("end_pos"), nil, nil, NameToType("Std::String::Position", env), Never{})
					namespace.DefineMethod("Returns the position of the first character\nof a piece of text in a string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start_pos"), nil, nil, NameToType("Std::String::Position", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Symbol").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))

				// Define methods
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Calculates a hash of the symbol.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the string associated with this symbol.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::SimpleSymbolLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::SimpleSymbolLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::SimpleSymbolLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::SimpleSymbolLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the string associated with this symbol.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_symbol"), nil, nil, NameToType("Std::Symbol", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Sync").(*Module)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("DiagnosticList").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces
					IncludeMixin(namespace, NewGeneric(NameToType("Std::List", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Diagnostic", env), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

					// Define methods
					namespace.DefineMethod("Create a new `ArrayList` containing the elements of `self`\nrepeated `n` times.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, Never{})
					namespace.DefineMethod("Create a new `ArrayList` containing the elements of `self`\nand another given `Tuple`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewUnion(NameToType("Std::Diagnostic", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), Never{})
					namespace.DefineMethod("Adds the given value to the list.\n\nReallocates the underlying array if it is\ntoo small to hold it.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, Self{}, Never{})
					namespace.DefineMethod("Check whether the given value is an `ArrayList`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
					namespace.DefineMethod("Check whether the given value is an `ArrayList` or `ArrayTuple`\nwith the same elements.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("=~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
					namespace.DefineMethod("Get the element under the given index.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Diagnostic", env), Never{})
					namespace.DefineMethod("Set the element under the given index to the given value.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]="), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false), NewParameter(value.ToSymbol("value"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Adds the given values to the list.\n\nReallocates the underlying array if it is\ntoo small to hold them.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("append"), nil, []*Parameter{NewParameter(value.ToSymbol("values"), NameToType("Std::Diagnostic", env), PositionalRestParameterKind, false)}, Self{}, Never{})
					namespace.DefineMethod("Returns the number of elements that can be\nheld by the underlying array.\n\nThis value will change when the list gets resized,\nand the underlying array gets reallocated", 0|METHOD_NATIVE_FLAG, value.ToSymbol("capacity"), nil, nil, NameToType("Std::Int", env), Never{})
					namespace.DefineMethod("Removes all elements from the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("clear"), nil, nil, Void{}, Never{})
					namespace.DefineMethod("Check whether the given `value` is present in this list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("contains"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Diagnostic", env), NameToType("Std::Diagnostic", env), NameToType("Std::Diagnostic", env), INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("value"), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :contains", true), NameToType("Std::Diagnostic", env), NameToType("Std::Diagnostic", env), NameToType("Std::Diagnostic", env), INVARIANT), NormalParameterKind, false)}, Bool{}, Never{})
					namespace.DefineMethod("Returns the internal `DiagnosticList`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("diagnostic_list"), nil, nil, NameToType("Std::DiagnosticList", env), Never{})
					namespace.DefineMethod("Mutates the list.\n\nReallocates the underlying array to hold\nthe given number of new elements.\n\nExpands the `capacity` of the list by `new_slots`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("grow"), nil, []*Parameter{NewParameter(value.ToSymbol("new_slots"), NameToType("Std::Int", env), NormalParameterKind, false)}, Self{}, Never{})
					namespace.DefineMethod("Check if this diagnostic list contains a fail.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_failure"), nil, nil, Bool{}, Never{})
					namespace.DefineMethod("Returns an iterator that iterates\nover each element of the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iter"), nil, nil, NameToType("Std::Sync::DiagnosticList::Iterator", env), Never{})
					namespace.DefineMethod("Returns the number of left slots for new elements\nin the underlying array.\nIt tells you how many more elements can be\nadded to the list before the underlying array gets\nreallocated.\n\nIt is always equal to `capacity - length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("left_capacity"), nil, nil, NameToType("Std::Int", env), Never{})
					namespace.DefineMethod("Returns the number of elements present in the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("length"), nil, nil, NameToType("Std::Int", env), Never{})
					namespace.DefineMethod("Iterates over the elements of this list,\nyielding them to the given closure.\n\nReturns a new List that consists of the elements returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::ArrayList", env).(*Class), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), INVARIANT)}, []value.Symbol{value.ToSymbol("Val")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
					namespace.DefineMethod("Iterates over the elements of this list,\nyielding them to the given closure.\n\nMutates the list in place replacing the elements with the ones\nreturned by the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map_mut"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, NameToType("Std::Diagnostic", env), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Self{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map_mut", true), Never{}, Any{}, nil, INVARIANT))
					namespace.DefineMethod("Adds the given value to the list.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("push"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Removes the element from the list.\n\nReturns `true` if the element has been removed,\notherwise returns `false.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("remove"), nil, []*Parameter{NewParameter(value.ToSymbol("value"), NameToType("Std::Diagnostic", env), NormalParameterKind, false)}, Bool{}, Never{})
					namespace.DefineMethod("Get the human readable\nstring representation of this diagnostic list\nwith a snippet of source code.\n\nIt will attempt to read the source fragment from the file\nif no source string is given.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_human_string"), nil, []*Parameter{NewParameter(value.ToSymbol("style"), Bool{}, DefaultValueParameterKind, false), NewParameter(value.ToSymbol("colorizer"), NewNilable(NameToType("Std::Colorizer", env)), DefaultValueParameterKind, false)}, Void{}, Never{})

					// Define constants

					// Define instance variables

					{
						namespace := namespace.MustSubtype("Iterator").(*Class)

						namespace.Name() // noop - avoid unused variable error

						// Include mixins and implement interfaces
						IncludeMixin(namespace, NewGeneric(NameToType("Std::ResettableIterator::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Diagnostic", env), COVARIANT), value.ToSymbol("Err"): NewTypeArgument(Never{}, COVARIANT)}, []value.Symbol{value.ToSymbol("Val"), value.ToSymbol("Err")})))

						// Define methods
						namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("list"), NameToType("Std::DiagnosticList", env), NormalParameterKind, false)}, Void{}, Never{})
						namespace.DefineMethod("Get the next element of the list.\nThrows `:stop_iteration` when there are no more elements.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("next"), nil, nil, NameToType("Std::Diagnostic", env), NewSymbolLiteral("stop_iteration"))
						namespace.DefineMethod("Resets the state of the iterator.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("reset"), nil, nil, Void{}, Never{})

						// Define constants

						// Define instance variables
					}
				}
				{
					namespace := namespace.MustSubtype("Mutex").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Locks the mutex.\nIf the mutex is already locked it blocks the current thread\nuntil the mutex becomes available.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("lock"), nil, nil, Void{}, Never{})
					namespace.DefineMethod("Unlocks the mutex.\nIf the mutex is already unlocked an unchecked error gets thrown.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("unlock"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("Once").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Executes the given function if it is the first call\nfor this instance of `Once`.\nOtherwise it's a noop.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("call"), []*TypeParameter{NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :call", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, nil, Void{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :call", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, Void{}, NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :call", true), Never{}, Any{}, nil, INVARIANT))

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("ROMutex").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("rwmutex"), NameToType("Std::Sync::RWMutex", env), DefaultValueParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Locks the mutex for reading.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("lock"), nil, nil, Void{}, Never{})
					namespace.DefineMethod("Returns the underlying RWMutex.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("rwmutex"), nil, nil, NameToType("Std::Sync::RWMutex", env), Never{})
					namespace.DefineMethod("Unlocks the mutex for reading.\nIf the mutex is already unlocked for reading an unchecked error gets thrown.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("unlock"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("RWMutex").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Locks the mutex for writing.\nIf the mutex is already locked for writing it blocks the current thread\nuntil the mutex becomes available.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("lock"), nil, nil, Void{}, Never{})
					namespace.DefineMethod("Locks the mutex for reading.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("read_lock"), nil, nil, Void{}, Never{})
					namespace.DefineMethod("Unlocks the mutex for reading.\nIf the mutex is already unlocked for reading an unchecked error gets thrown.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("read_unlock"), nil, nil, Void{}, Never{})
					namespace.DefineMethod("Creates a read only wrapper around this mutex\nthat exposes `read_lock` and `read_unlock` methods as `lock` and `unlock`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_read_only"), nil, nil, NameToType("Std::Sync::ROMutex", env), Never{})
					namespace.DefineMethod("Unlocks the mutex for writing.\nIf the mutex is already unlocked for writing an unchecked error gets thrown.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("unlock"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
				{
					namespace := namespace.MustSubtype("WaitGroup").(*Class)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Initialises the counter with `n` elements.\n`0` is the default value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("#init"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), DefaultValueParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Adds n elements to the counter, which may be negative.\nIf the counter becomes zero, all threads blocked on `wait` are released.\nIf the counter goes negative an unchecked error gets thrown.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("add"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Decrements the counter by one.\nIf the counter becomes zero, all threads blocked on `wait` are released.\nIf the counter goes negative an unchecked error gets thrown.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("end"), nil, nil, Void{}, Never{})
					namespace.DefineMethod("Decrements the counter by `n`.\nIf the counter becomes zero, all threads blocked on `wait` are released.\nIf the counter goes negative an unchecked error gets thrown.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("remove"), nil, []*Parameter{NewParameter(value.ToSymbol("n"), NameToType("Std::Int", env), NormalParameterKind, false)}, Void{}, Never{})
					namespace.DefineMethod("Increments the counter by one.\nIf the counter becomes zero, all threads blocked on `wait` are released.\nIf the counter goes negative an unchecked error gets thrown.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("start"), nil, nil, Void{}, Never{})
					namespace.DefineMethod("Blocks the current thread until the internal counter of the `WaitGroup` reaches zero.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("wait"), nil, nil, Void{}, Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Thread").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Returns the current state of the thread.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("state"), nil, nil, NameToType("Std::Symbol", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("ThreadPool").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Closes the thread pool, shuts down all thread workers.\n\nA thread pool has to be closed when its no longer needed\notherwise it will never get garbage collected and the threads\nwill keep on waiting for work indefinitely.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("close"), nil, nil, Void{}, Never{})
				namespace.DefineMethod("Returns the number of available slots in the task\nqueue.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("task_queue_size"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the count of thread workers available in the pool.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("thread_count"), nil, nil, NameToType("Std::Int", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Time").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Adds the given duration to the time.\nReturns a new time object.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Duration", env), NormalParameterKind, false)}, NameToType("Std::Time", env), Never{})
				namespace.DefineMethod("Subtracts the given duration from the time.\nReturns a new time object.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Duration", env), NormalParameterKind, false)}, NameToType("Std::Time", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Time", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Time", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Time", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Time", env), NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Returns the attosecond offset within the second specified by `self` in the range `0...999999999999999999`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("attosecond"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the day of the month.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("day"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Calculates the difference between two time objects.\nReturns a duration.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("diff"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::Time", env), NormalParameterKind, false)}, NameToType("Std::Duration", env), Never{})
				namespace.DefineMethod("Returns the femtosecond offset within the second specified by `self` in the range `0...999999999999999`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("femtosecond"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Create a string formatted according to the given format string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("format"), nil, []*Parameter{NewParameter(value.ToSymbol("fmt"), NameToType("Std::String", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the hour offset within the day specified by `self` in the range `0...23`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hour"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the hour of the day in a twelve hour clock.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hour12"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Whether the current hour is AM.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_am"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Checks whether the day of the week is friday.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_friday"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Checks whether the timezone it local (the same as the system timezone).", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_local"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Checks whether the day of the week is monday.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_monday"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Whether the current hour is PM.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_pm"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Checks whether the day of the week is saturday.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_saturday"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Checks whether the day of the week is sunday.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_sunday"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Checks whether the day of the week is thursday.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_thursday"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Checks whether the day of the week is tuesday.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_tuesday"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Checks whether the timezone it UTC.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_utc"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Checks whether the day of the week is wednesday.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_wednesday"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns the ISO 8601 week number in which `self` occurs.\nWeek ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1 of year n+1.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iso_week"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the ISO 8601 year in which `self` occurs.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("iso_year"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Convert the time to the local timezone.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("local"), nil, nil, NameToType("Std::Time", env), Never{})
				namespace.DefineMethod("Returns `\"AM\"` or `\"PM\"` based on the hour.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("meridiem"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the microsecond offset within the second specified by `self` in the range `0...999999`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("microsecond"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the millisecond offset within the second specified by `self` in the range `0...999`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("millisecond"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the minute offset within the hour specified by `self` in the range `0...59`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("minute"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the month in which `self` occurs.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("month"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the day of the month.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("month_day"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the nanosecond offset within the second specified by `self` in the range `0...999999999`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("nanosecond"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the picosecond offset within the second specified by `self` in the range `0...999999999999`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("picosecond"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the second offset within the minute specified by `self` in the range `0...59`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("second"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Create a string formatted according to the given format string.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("strftime"), nil, []*Parameter{NewParameter(value.ToSymbol("fmt"), NameToType("Std::String", env), NormalParameterKind, false)}, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Return the timezone associated with this Time object.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("timezone"), nil, nil, NameToType("Std::Timezone", env), Never{})
				namespace.DefineMethod("Return the name of the timezone associated with this Time object.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("timezone_name"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the offset of the timezone in hours east of UTC.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("timezone_offset_hours"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the offset of the timezone in seconds east of UTC.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("timezone_offset_seconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Convert the time to the local timezone.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_local"), nil, nil, NameToType("Std::Time", env), Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_string"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Convert the time to the UTC zone.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_utc"), nil, nil, NameToType("Std::Time", env), Never{})
				namespace.DefineMethod("Returns the number of attoseconds elapsed since January 1, 1970 UTC", 0|METHOD_NATIVE_FLAG, value.ToSymbol("unix_attoseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of femtoseconds elapsed since January 1, 1970 UTC", 0|METHOD_NATIVE_FLAG, value.ToSymbol("unix_femtoseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of microseconds elapsed since January 1, 1970 UTC", 0|METHOD_NATIVE_FLAG, value.ToSymbol("unix_microseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of milliseconds elapsed since January 1, 1970 UTC", 0|METHOD_NATIVE_FLAG, value.ToSymbol("unix_milliseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of nanoseconds elapsed since January 1, 1970 UTC", 0|METHOD_NATIVE_FLAG, value.ToSymbol("unix_nanoseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of picoseconds elapsed since January 1, 1970 UTC", 0|METHOD_NATIVE_FLAG, value.ToSymbol("unix_picoseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of seconds elapsed since January 1, 1970 UTC", 0|METHOD_NATIVE_FLAG, value.ToSymbol("unix_seconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of yoctoseconds elapsed since January 1, 1970 UTC", 0|METHOD_NATIVE_FLAG, value.ToSymbol("unix_yoctoseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of zeptoseconds elapsed since January 1, 1970 UTC", 0|METHOD_NATIVE_FLAG, value.ToSymbol("unix_zeptoseconds"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Convert the time to the UTC zone.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("utc"), nil, nil, NameToType("Std::Time", env), Never{})
				namespace.DefineMethod("The week number of the current year as a decimal number,\nrange 0 to 53, starting with the first Monday\nas the first day of week 1.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("week"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("The week number of the current year as a decimal number,\nrange 0 to 53, starting with the first Monday\nas the first day of week 1.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("week_from_monday"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("The week number of the current year as a decimal number,\nrange 0 to 53, starting with the first Sunday\nas the first day of week 01.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("week_from_sunday"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of the day of the week, where 1 is Monday, 7 is Sunday", 0|METHOD_NATIVE_FLAG, value.ToSymbol("weekday"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of the day of the week, where 1 is Monday, 7 is Sunday", 0|METHOD_NATIVE_FLAG, value.ToSymbol("weekday_from_monday"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the number of the day of the week, where 0 is Sunday, 6 is Saturday", 0|METHOD_NATIVE_FLAG, value.ToSymbol("weekday_from_sunday"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the name of the day of the week.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("weekday_name"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the year in which `self` occurs.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("year"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the day of the year.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("year_day"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the yoctosecond offset within the second specified by `self` in the range `0...999999999999999999999999`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("yoctosecond"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the zeptosecond offset within the second specified by `self` in the range `0...999999999999999999999`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("zeptosecond"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Return the timezone associated with this Time object.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("zone"), nil, nil, NameToType("Std::Timezone", env), Never{})
				namespace.DefineMethod("Return the name of the timezone associated with this Time object.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("zone_name"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the offset of the timezone in hours east of UTC.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("zone_offset_hours"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Returns the offset of the timezone in seconds east of UTC.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("zone_offset_seconds"), nil, nil, NameToType("Std::Int", env), Never{})

				// Define constants
				namespace.DefineConstant(value.ToSymbol("DEFAULT_FORMAT"), NameToType("Std::String", env))

				// Define instance variables

				{
					namespace := namespace.Singleton()

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Returns the current time.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("now"), nil, nil, NameToType("Std::Time", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Timezone").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_local"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("is_utc"), nil, nil, Bool{}, Never{})
				namespace.DefineMethod("Returns the name of the Timezone eg. `\"Local\"`, `\"UTC\"`, `\"Europe/Warsaw\"`", 0|METHOD_NATIVE_FLAG, value.ToSymbol("name"), nil, nil, NameToType("Std::String", env), Never{})

				// Define constants
				namespace.DefineConstant(value.ToSymbol("LOCAL"), NameToType("Std::Timezone", env))
				namespace.DefineConstant(value.ToSymbol("UTC"), NameToType("Std::Timezone", env))

				// Define instance variables

				{
					namespace := namespace.Singleton()

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Returns the Timezone for the given name.\n\nIf the name is \"\" or \"UTC\" the UTC timezone gets returned. If the name is \"Local\", the local (system) timezone gets returned.\n\nOtherwise, the name is taken to be a location name corresponding to a file in the IANA Time Zone database, such as `\"Europe/Warsaw\"`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false)}, NameToType("Std::Timezone", env), Never{})
					namespace.DefineMethod("Returns the Timezone for the given name.\n\nIf the name is \"\" or \"UTC\" the UTC timezone gets returned. If the name is \"Local\", the local (system) timezone gets returned.\n\nOtherwise, the name is taken to be a location name corresponding to a file in the IANA Time Zone database, such as `\"Europe/Warsaw\"`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("get"), nil, []*Parameter{NewParameter(value.ToSymbol("name"), NameToType("Std::String", env), NormalParameterKind, false)}, NameToType("Std::Timezone", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("True").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Bool", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))

				// Define methods
				namespace.DefineMethod("Calculates a hash of the value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::TrueLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::TrueLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::TrueLiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::TrueLiteralNode", env), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("Tuple").(*Mixin)

				namespace.Name() // noop - avoid unused variable error

				// Set up type parameters
				var typeParam *TypeParameter
				typeParams := make([]*TypeParameter, 1)

				typeParam = NewTypeParameter(value.ToSymbol("Element"), namespace, Never{}, Any{}, nil, COVARIANT)
				typeParams[0] = typeParam
				namespace.DefineSubtype(value.ToSymbol("Element"), typeParam)
				namespace.DefineConstant(value.ToSymbol("Element"), NoValue{})

				namespace.SetTypeParameters(typeParams)

				// Include mixins and implement interfaces
				IncludeMixin(namespace, NewGeneric(NameToType("Std::ImmutableCollection::Base", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Val"): NewTypeArgument(NameToType("Std::Tuple::Element", env), COVARIANT)}, []value.Symbol{value.ToSymbol("Val")})))

				// Define methods
				namespace.DefineMethod("Create a new `Tuple` containing the elements of `self`\nand another given `Tuple`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("other"), NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NewUnion(NameToType("Std::Tuple::Element", env), NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :\"+\"", true), Never{}, Any{}, nil, INVARIANT)), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})), Never{})
				namespace.DefineMethod("Get the element under the given index.\n\nThrows an unchecked error if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("[]"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Tuple::Element", env), Never{})
				namespace.DefineMethod("Get the element under the given index.\n\nThrows an error if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("at"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::Tuple::Element", env), NameToType("Std::OutOfRangeError", env))
				namespace.DefineMethod("Iterates over the elements of this tuple,\nyielding them to the given closure.\n\nReturns a new tuple that consists of the elements returned\nby the given closure.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("map"), []*TypeParameter{NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)}, []*Parameter{NewParameter(value.ToSymbol("fn"), NewClosureWithMethod("", 0|METHOD_NATIVE_FLAG, value.ToSymbol("call"), nil, []*Parameter{NewParameter(value.ToSymbol("element"), NameToType("Std::Tuple::Element", env), NormalParameterKind, false)}, NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT)), NormalParameterKind, false)}, NewGeneric(NameToType("Std::Tuple", env).(*Mixin), NewTypeArguments(TypeArgumentMap{value.ToSymbol("Element"): NewTypeArgument(NewTypeParameter(value.ToSymbol("V"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT), COVARIANT)}, []value.Symbol{value.ToSymbol("Element")})), NewTypeParameter(value.ToSymbol("E"), NewTypeParamNamespace("Type Parameter Container of :map", true), Never{}, Any{}, nil, INVARIANT))
				namespace.DefineMethod("Get the element under the given index.\n\nReturns `nil` if the index is a negative number\nor is greater or equal to `length`.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("try_at"), nil, []*Parameter{NewParameter(value.ToSymbol("index"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NewNilable(NameToType("Std::Tuple::Element", env)), Never{})

				// Define constants

				// Define instance variables
			}
			{
				namespace := namespace.MustSubtype("UInt16").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10u16\n\tvar b = 3u16\n\ta % b #=> 1u16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Performs bitwise AND.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Multiply this integer by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Add `other` to this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1u16\n\t+a #=> 1u16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Subtract `other` from this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1u16\n\t-a #=> 65535u16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Performs bitwise XOR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("^"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Calculates a hash of the int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::UInt16LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::UInt16LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::UInt16LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::UInt16LiteralNode", env), Never{})
				namespace.DefineMethod("Converts the integer to a floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts to an automatically resizable integer type.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise OR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("|"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt16", env), NormalParameterKind, false)}, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4u16 #=> 65531u16\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("~"), nil, nil, NameToType("Std::UInt16", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a uint16.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("UInt32").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10u32\n\tvar b = 3u32\n\ta % b #=> 1u32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Performs bitwise AND.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Multiply this integer by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Add `other` to this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1u32\n\t+a #=> 1u32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Subtract `other` from this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1u32\n\t-a #=> 4294967295u32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Performs bitwise XOR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("^"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Calculates a hash of the int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::UInt32LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::UInt32LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::UInt32LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::UInt32LiteralNode", env), Never{})
				namespace.DefineMethod("Converts the integer to a floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts to an automatically resizable integer type.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise OR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("|"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt32", env), NormalParameterKind, false)}, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4u32 #=> 4294967291u32\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("~"), nil, nil, NameToType("Std::UInt32", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a uint32.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("UInt64").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10u64\n\tvar b = 3u64\n\ta % b #=> 1u64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Performs bitwise AND.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Multiply this integer by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Add `other` to this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1u64\n\t+a #=> 1u64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Subtract `other` from this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1u64\n\t-a #=> 18446744073709551615u64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Performs bitwise XOR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("^"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Calculates a hash of the int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::UInt64LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::UInt64LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::UInt64LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::UInt64LiteralNode", env), Never{})
				namespace.DefineMethod("Converts the integer to a floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts to an automatically resizable integer type.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise OR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("|"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt64", env), NormalParameterKind, false)}, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4u64 #=> 18446744073709551611u64\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("~"), nil, nil, NameToType("Std::UInt64", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a uint64.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("UInt8").(*Class)

				namespace.Name() // noop - avoid unused variable error
				namespace.SetParent(NameToNamespace("Std::Value", env))

				// Include mixins and implement interfaces
				ImplementInterface(namespace, NameToType("Std::Hashable", env).(*Interface))
				ImplementInterface(namespace, NewGeneric(NameToType("Std::Comparable", env).(*Interface), NewTypeArguments(TypeArgumentMap{value.ToSymbol("T"): NewTypeArgument(Self{}, INVARIANT)}, []value.Symbol{value.ToSymbol("T")})))

				// Define methods
				namespace.DefineMethod("Returns the remainder of dividing by `other`.\n\n```\n\tvar a = 10u8\n\tvar b = 3u8\n\ta % b #=> 1u8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("%"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise AND.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise AND NOT (bit clear).", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("&~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Multiply this integer by `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("*"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Exponentiate this integer, raise it to the power of `other`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("**"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Add `other` to this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Get the next integer by incrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("++"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Returns itself.\n\n```\n\tvar a = 1u8\n\t+a #=> 1u8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("+@"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Subtract `other` from this integer.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Get the previous integer by decrementing by `1`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("--"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Returns the result of negating the integer.\n\n```\n\tvar a = 1u8\n\t-a #=> 255u8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("-@"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Divide this integer by another integer.\nThrows an unchecked runtime error when dividing by `0`.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("/"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Returns an integer shifted left by `other` positions, or right if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<<<"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Compare this integer with another integer.\nReturns `1` if it is greater, `0` if they're equal, `-1` if it's less than the other.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("<=>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::Bool", env), Never{})
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Returns an integer shifted right by `other` positions, or left if `other` is negative.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol(">>>"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::AnyInt", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise XOR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("^"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Calculates a hash of the int.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Return a human readable string\nrepresentation of this object\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_expr_node"), nil, nil, NameToType("Std::Elk::AST::UInt8LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_node"), nil, nil, NameToType("Std::Elk::AST::UInt8LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_pattern_node"), nil, nil, NameToType("Std::Elk::AST::UInt8LiteralNode", env), Never{})
				namespace.DefineMethod("Returns the AST Node that represents the same value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_ast_type_node"), nil, nil, NameToType("Std::Elk::AST::UInt8LiteralNode", env), Never{})
				namespace.DefineMethod("Converts the integer to a floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float"), nil, nil, NameToType("Std::Float", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float32"), nil, nil, NameToType("Std::Float32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit floating point number.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_float64"), nil, nil, NameToType("Std::Float64", env), Never{})
				namespace.DefineMethod("Converts to an automatically resizable integer type.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int"), nil, nil, NameToType("Std::Int", env), Never{})
				namespace.DefineMethod("Converts the integer to a 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int16"), nil, nil, NameToType("Std::Int16", env), Never{})
				namespace.DefineMethod("Converts the integer to a 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int32"), nil, nil, NameToType("Std::Int32", env), Never{})
				namespace.DefineMethod("Converts the integer to a 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int64"), nil, nil, NameToType("Std::Int64", env), Never{})
				namespace.DefineMethod("Converts the integer to a 8-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_int8"), nil, nil, NameToType("Std::Int8", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 16-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint16"), nil, nil, NameToType("Std::UInt16", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 32-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint32"), nil, nil, NameToType("Std::UInt32", env), Never{})
				namespace.DefineMethod("Converts the integer to an unsigned 64-bit integer.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint64"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns itself.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Performs bitwise OR.", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("|"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), NameToType("Std::UInt8", env), NormalParameterKind, false)}, NameToType("Std::UInt8", env), Never{})
				namespace.DefineMethod("Returns the result of applying bitwise NOT on the bits\nof this integer.\n\n```\n\t~4u8 #=> 251u8\n```", 0|METHOD_SEALED_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("~"), nil, nil, NameToType("Std::UInt8", env), Never{})

				// Define constants

				// Define instance variables

				{
					namespace := namespace.MustSubtype("Convertible").(*Interface)

					namespace.Name() // noop - avoid unused variable error

					// Include mixins and implement interfaces

					// Define methods
					namespace.DefineMethod("Convert the value to a uint8.", 0|METHOD_ABSTRACT_FLAG|METHOD_NATIVE_FLAG, value.ToSymbol("to_uint8"), nil, nil, NameToType("Std::UInt8", env), Never{})

					// Define constants

					// Define instance variables
				}
			}
			{
				namespace := namespace.MustSubtype("Value").(*Class)

				namespace.Name() // noop - avoid unused variable error

				// Include mixins and implement interfaces

				// Define methods
				namespace.DefineMethod("Compares this value with another value.\n\nReturns `true` when they are instances of the same class,\nand are equal.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("=="), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Compares this value with another value.\nReturns `true` when they are equal.\n\nInstances of different (but similar) classes\nmay be treated as equal.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("=~"), nil, []*Parameter{NewParameter(value.ToSymbol("other"), Any{}, NormalParameterKind, false)}, Bool{}, Never{})
				namespace.DefineMethod("Returns the class of the value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("class"), nil, nil, NewSingletonOf(Self{}), Never{})
				namespace.DefineMethod("Returns a shallow copy of the value.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("copy"), nil, nil, Self{}, Never{})
				namespace.DefineMethod("Returns a hash of the value,\n  that is used to calculate the slot\n  in a HashMap, HashRecord or HashSet\n  where the value will be stored.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("hash"), nil, nil, NameToType("Std::UInt64", env), Never{})
				namespace.DefineMethod("Returns a human readable `String`\nrepresentation of this value\nfor debugging etc.", 0|METHOD_NATIVE_FLAG, value.ToSymbol("inspect"), nil, nil, NameToType("Std::String", env), Never{})

				// Define constants

				// Define instance variables
			}
		}
	}
}
