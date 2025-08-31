package token

import "iter"

//go:generate stringer -type=Type

// Represents the type of token
type Type uint16

// Name of the token.
func (t Type) Name() string {
	if int(t) > len(tokenNames) {
		return "UNKNOWN"
	}

	return tokenNames[t]
}

func (t Type) TypeName() string {
	return t.String()
}

func Length() int {
	return len(tokenNames)
}

func Types() iter.Seq2[uint16, string] {
	return func(yield func(uint16, string) bool) {
		for i := range uint16(Length()) {
			if !yield(i, Type(i).String()) {
				return
			}
		}
	}
}

// Returns `true` if the token can be a beginning of
// an argument to a function call without parentheses
// eg. `foo 2`
func (t Type) IsValidAsArgumentToNoParenFunctionCall() bool {
	switch t {
	case BANG, TILDE, LBRACE, DOLLAR_IDENTIFIER, PUBLIC_IDENTIFIER, PRIVATE_IDENTIFIER,
		PUBLIC_CONSTANT, PRIVATE_CONSTANT, INSTANCE_VARIABLE, COLON, CHAR_LITERAL, RAW_CHAR_LITERAL,
		RAW_STRING, STRING_BEG, NIL, FALSE, TRUE, LOOP, DEF, SIG,
		INIT, CLASS, STRUCT, MODULE, MIXIN, INTERFACE, ENUM, TYPE, TYPEDEF,
		VAR, VAL, CONST, DO, ALIAS, SELF, SUPER, SWITCH,
		INT, INT64, UINT64, INT32, UINT32, INT16, UINT16, INT8, UINT8,
		FLOAT, BIG_FLOAT, FLOAT64, FLOAT32:
		return true
	}

	return t.IsSpecialCollectionLiteralBeg()
}

// Returns `true` if the token can be an end of
// a range value eg. `...2`
func (t Type) IsValidAsEndInRangeLiteral() bool {
	switch t {
	case SCOPE_RES_OP, BANG, TILDE, LBRACE, LPAREN, LBRACKET, DOLLAR_IDENTIFIER, PUBLIC_IDENTIFIER, PRIVATE_IDENTIFIER,
		PUBLIC_CONSTANT, PRIVATE_CONSTANT, INSTANCE_VARIABLE,
		RAW_STRING, STRING_BEG, CHAR_LITERAL, RAW_CHAR_LITERAL, FLOAT, FLOAT32, FLOAT64,
		NIL, FALSE, TRUE, LOOP, ENUM,
		VAR, VAL, CONST, DO, SELF, SUPER, SWITCH, MINUS, PLUS:
		return true
	}

	if t.IsIntLiteral() || t.IsSpecialCollectionLiteralBeg() {
		return true
	}

	return false
}

// Returns `true` if the token can be an end of
// a range pattern eg. `...2`
func (t Type) IsValidAsEndInRangePattern() bool {
	switch t {
	case RAW_STRING, STRING_BEG, CHAR_LITERAL, RAW_CHAR_LITERAL,
		FLOAT, FLOAT32, FLOAT64, NIL, FALSE, TRUE, MINUS, PLUS,
		PUBLIC_CONSTANT, PRIVATE_CONSTANT, SCOPE_RES_OP:
		return true
	}

	return t.IsIntLiteral()
}

func (t Type) IsCollectionLiteralBeg() bool {
	switch t {
	case LBRACKET, LBRACE, TUPLE_LITERAL_BEG, RECORD_LITERAL_BEG, HASH_SET_LITERAL_BEG:
		return true
	default:
		return t.IsSpecialCollectionLiteral()
	}
}

// return `true` if the token is the beginning of a special
// collection literal
func (t Type) IsSpecialCollectionLiteralBeg() bool {
	switch t {
	case WORD_ARRAY_LIST_BEG, SYMBOL_ARRAY_LIST_BEG, HEX_ARRAY_LIST_BEG, BIN_ARRAY_LIST_BEG, WORD_HASH_SET_BEG,
		SYMBOL_HASH_SET_BEG, HEX_HASH_SET_BEG, BIN_HASH_SET_BEG, WORD_ARRAY_TUPLE_BEG, SYMBOL_ARRAY_TUPLE_BEG,
		HEX_ARRAY_TUPLE_BEG, BIN_ARRAY_TUPLE_BEG, HASH_SET_LITERAL_BEG, TUPLE_LITERAL_BEG:
		return true
	default:
		return false
	}
}

// Check whether the token represents a special collection literal like `%w[]`
func (t Type) IsSpecialCollectionLiteral() bool {
	return LABEL_SPECIAL_COLLECTION_LITERAL_BEG < t && t < LABEL_SPECIAL_COLLECTION_LITERAL_END
}

// Check whether the token represents a regex flag
func (t Type) IsRegexFlag() bool {
	return LABEL_REGEX_FLAG_BEG < t && t < LABEL_REGEX_FLAG_END
}

// Check whether the token marks the end of the file.
func (t Type) IsEndOfFile() bool {
	return t == END_OF_FILE
}

// Check whether the token is a keyword.
func (t Type) IsKeyword() bool {
	return LABEL_KEYWORD_BEG < t && t < LABEL_KEYWORD_END
}

// Check whether the token is a literal.
func (t Type) IsLiteral() bool {
	return LABEL_LITERAL_BEG < t && t < LABEL_LITERAL_END
}

// Check whether the token is an int literal.
func (t Type) IsIntLiteral() bool {
	return LABEL_INT_LITERAL_BEG < t && t < LABEL_INT_LITERAL_END
}

// Check whether the token is a float literal.
func (t Type) IsFloatLiteral() bool {
	return LABEL_FLOAT_LITERAL_BEG < t && t < LABEL_FLOAT_LITERAL_END
}

// Check whether the token is a an operator.
func (t Type) IsOperator() bool {
	return LABEL_OP_BEG < t && t < LABEL_OP_END
}

// Check whether the token is a valid simple symbol content.
func (t Type) IsValidSimpleSymbolContent() bool {
	return t.IsIdentifier() || t == RAW_STRING || t.IsKeyword() || t.IsOverridableOperator()
}

// Check whether the token is a valid function name.
func (t Type) IsValidFunctionName() bool {
	return t == DOLLAR_IDENTIFIER || t == PUBLIC_IDENTIFIER || t == PRIVATE_IDENTIFIER
}

// Check whether the token is a valid method name (without operators).
func (t Type) IsValidRegularMethodName() bool {
	return t == DOLLAR_IDENTIFIER || t == PUBLIC_IDENTIFIER || t == PRIVATE_IDENTIFIER || t.IsKeyword()
}

// Check whether the token is a valid macro name.
func (t Type) IsValidMacroName() bool {
	return t == DOLLAR_IDENTIFIER || t == PUBLIC_IDENTIFIER || t.IsKeyword()
}

// Check whether the token is a valid method name (including operators).
func (t Type) IsValidMethodName() bool {
	return t == SHORT_UNQUOTE_BEG || t.IsValidRegularMethodName() || t.IsOverridableOperator()
}

// Check whether the token is a valid method name in method
// call expressions.
func (t Type) IsValidPublicMethodName() bool {
	return t == DOLLAR_IDENTIFIER || t == PUBLIC_IDENTIFIER || t.IsKeyword() || t.IsOverridableOperator()
}

// Check whether the token is an overridable operator.
// Overridable operators can be used as method names.
func (t Type) IsOverridableOperator() bool {
	return LABEL_OVERRIDABLE_OP_BEG < t && t < LABEL_OP_END
}

// Check whether the token is a non overridable operator.
func (t Type) IsNonOverridableOperator() bool {
	return t.IsOperator() && !t.IsOverridableOperator()
}

// Check whether the token is an identifier.
func (t Type) IsIdentifier() bool {
	return LABEL_IDENTIFIER_BEG < t && t < LABEL_IDENTIFIER_END
}

// Check whether the token can separate statements.
func (t Type) IsStatementSeparator() bool {
	return t == NEWLINE || t == SEMICOLON
}

// Check whether the token is an assignment operator.
func (t Type) IsAssignmentOperator() bool {
	return LABEL_ASSIGN_OP_BEG < t && t < LABEL_ASSIGN_OP_END
}

// Check whether the token is an equality operator.
func (t Type) IsEqualityOperator() bool {
	switch t {
	case EQUAL_EQUAL, NOT_EQUAL, STRICT_EQUAL,
		STRICT_NOT_EQUAL, LAX_EQUAL, LAX_NOT_EQUAL:
		return true
	default:
		return false
	}
}

// Check whether the token is a comparison operator.
func (t Type) IsComparisonOperator() bool {
	switch t {
	case LESS, LESS_EQUAL, GREATER,
		GREATER_EQUAL, ISA_OP, REVERSE_ISA_OP,
		INSTANCE_OF_OP, REVERSE_INSTANCE_OF_OP, SPACESHIP_OP:
		return true
	default:
		return false
	}
}

const (
	ZERO_VALUE       Type = iota // Zero value for Type
	ERROR                        // ERROR Token with a message
	END_OF_FILE                  // End Of File has been reached
	TEXT                         // Plain text
	NEWLINE                      // Newline `\n`, `\r\n`
	SEMICOLON                    // SEMICOLON `;`
	COMMA                        // Comma `,`
	DOT                          // Dot `.`
	QUESTION_DOT                 // Safe method call operator `?.`
	DOT_DOT                      // Cascade call operator `..`
	QUESTION_DOT_DOT             // Safe cascade call operator `?..`
	COLON                        // Colon `:`
	QUESTION                     // Question mark `?`

	// Operators start here
	LABEL_OP_BEG

	THIN_ARROW   // Thin arrow `->` (closure arrow)
	WIGGLY_ARROW // Wiggly arrow `~>` (lambda arrow)
	THICK_ARROW  // Thick arrow `=>`

	// Assignment operators start here
	LABEL_ASSIGN_OP_BEG
	EQUAL_OP                // Equal (assignment) `=`
	MINUS_EQUAL             // Minus equal `-=`
	PLUS_EQUAL              // Plus equal `+=`
	STAR_EQUAL              // Star equal `*=`
	SLASH_EQUAL             // Slash equal `/=`
	STAR_STAR_EQUAL         // Two stars equal `**=`
	TILDE_EQUAL             // Tilde equal `~=`
	AND_AND_EQUAL           // Logical and equal `&&=`
	AND_EQUAL               // Bitwise and equal `&=`
	OR_OR_EQUAL             // Logical or `||=`
	OR_EQUAL                // Bitwise or equal `|=`
	XOR_EQUAL               // Bitwise xor equal `^=`
	QUESTION_QUESTION_EQUAL // Nil coalescing equal operator `??=`
	LBITSHIFT_EQUAL         // Left bitwise shift equal `<<=`
	LTRIPLE_BITSHIFT_EQUAL  // Triple left bitwise shift equal `<<<=`
	RBITSHIFT_EQUAL         // Right bitwise shift equal `>>=`
	RTRIPLE_BITSHIFT_EQUAL  // Triple right bitwise shift equal `>>>=`
	PERCENT_EQUAL           // Percent equal `%=`
	COLON_EQUAL             // Colon equal `:=`
	LABEL_ASSIGN_OP_END     // Assignment operators end here

	SHORT_UNQUOTE_BEG      // Short unquote beginning `!{`
	LPAREN                 // Left parenthesis `(`
	RPAREN                 // Right parenthesis `)`
	LBRACE                 // Left brace `{`
	RBRACE                 // Right brace `}`
	LBRACKET               // Left bracket `[`
	QUESTION_LBRACKET      // Safe access `?[`
	RBRACKET               // Right bracket `]`
	SCOPE_RES_OP           // Scope resolution operator `::`
	COLON_COLON_LBRACKET   // Colon, colon, left bracket `::[`
	DOT_COLON              // Dot colon `.:`
	CLOSED_RANGE_OP        // Closed range operator `...`
	OPEN_RANGE_OP          // Open range operator `<.<`
	RIGHT_OPEN_RANGE_OP    // Right open range operator `..<`
	LEFT_OPEN_RANGE_OP     // Left open range operator `<..`
	PIPE_OP                // Pipe operator `|>`
	AND_AND                // Logical and `&&`
	AND_BANG               // Logical expression sequencing operator `&!` with the precedence of `&&`
	OR_OR                  // Logical or `||`
	OR_BANG                // Logical expression sequencing operator `|!` with the precedence of `||`
	NOT_EQUAL              // Not equal `!=`
	LAX_NOT_EQUAL          // Lax not equal operator `!~`
	STRICT_EQUAL           // Strict equal `===`
	STRICT_NOT_EQUAL       // Strict not equal `!==`
	QUESTION_QUESTION      // Nil coalescing operator `??`
	BANG                   // Logical not `!`
	ISA_OP                 // "is a" operator `<:`
	REVERSE_ISA_OP         // Reverse "is a" operator `:>`
	INSTANCE_OF_OP         // Instance of operator `<<:`
	REVERSE_INSTANCE_OF_OP // Reverse instance of operator `:>>`

	// Overridable operators start here
	LABEL_OVERRIDABLE_OP_BEG
	PLUS_PLUS        // Increment operator `++`
	MINUS_MINUS      // Decrement operator `--`
	PLUS_AT          // Negate `+@`
	MINUS_AT         // Negate `-@`
	MINUS            // Minus `-`
	PLUS             // Plus `+`
	STAR             // Star `*`
	SLASH            // Slash `/`
	STAR_STAR        // Two stars `**`
	LESS             // Less than `<`
	LESS_EQUAL       // Less than or equal `<=`
	GREATER          // Greater than `>`
	GREATER_EQUAL    // Greater than or equal `>=`
	SPACESHIP_OP     // Spaceship operator `<=>`
	EQUAL_EQUAL      // Equal (comparison) `==`
	LAX_EQUAL        // Lax equality operator `=~`
	TILDE            // Tilde `~`
	AND              // Bitwise and `&`
	AND_TILDE        // Bitwise and not `&~`
	OR               // Bitwise or `|`
	XOR              // Bitwise xor `^`
	LBITSHIFT        // Left bitwise shift `<<`
	LTRIPLE_BITSHIFT // Triple left bitwise shift `<<<`
	RBITSHIFT        // Right bitwise shift `>>`
	RTRIPLE_BITSHIFT // Triple right bitwise shift `>>>`
	PERCENT          // Percent `%`
	LABEL_OP_END     // Operators end here

	// Identifiers start here
	LABEL_IDENTIFIER_BEG
	DOLLAR_IDENTIFIER    // Dollar Identifier eg. `$foo`
	PUBLIC_IDENTIFIER    // Identifier eg. `foo`
	PRIVATE_IDENTIFIER   // Identifier with a initial underscore eg. `_foo`
	PUBLIC_CONSTANT      // Constant (identifier with an initial capital letter) eg. `Foo`
	PRIVATE_CONSTANT     // Constant with an initial underscore eg. `_Foo`
	LABEL_IDENTIFIER_END // Identifiers end here

	INSTANCE_VARIABLE // Instance variable token eg. `@foo`

	// Literals start here
	LABEL_LITERAL_BEG

	// Special collection literals start here
	LABEL_SPECIAL_COLLECTION_LITERAL_BEG

	WORD_ARRAY_LIST_BEG   // Word array literal beginning `\w[`
	WORD_ARRAY_LIST_END   // Word array literal end `]`
	SYMBOL_ARRAY_LIST_BEG // Symbol array literal beginning `\s[`
	SYMBOL_ARRAY_LIST_END // Symbol array literal end `]`
	HEX_ARRAY_LIST_BEG    // Hexadecimal integer array literal beginning `\x[`
	HEX_ARRAY_LIST_END    // Hexadecimal integer array literal end `]`
	BIN_ARRAY_LIST_BEG    // Binary integer array literal beginning `\b[`
	BIN_ARRAY_LIST_END    // Binary integer array literal end `]`

	WORD_HASH_SET_BEG   // Word set literal beginning `^w[`
	WORD_HASH_SET_END   // Word set literal end `]`
	SYMBOL_HASH_SET_BEG // Symbol set literal beginning `^s[`
	SYMBOL_HASH_SET_END // Symbol set literal end `]`
	HEX_HASH_SET_BEG    // Hexadecimal integer set literal beginning `^x[`
	HEX_HASH_SET_END    // Hexadecimal integer set literal end `]`
	BIN_HASH_SET_BEG    // Binary integer set literal beginning `^b[`
	BIN_HASH_SET_END    // Binary integer set literal end `]`

	WORD_ARRAY_TUPLE_BEG   // Word arrayTuple literal beginning `%w[`
	WORD_ARRAY_TUPLE_END   // Word arrayTuple literal end `]`
	SYMBOL_ARRAY_TUPLE_BEG // Symbol arrayTuple literal beginning `%s[`
	SYMBOL_ARRAY_TUPLE_END // Symbol arrayTuple literal end `]`
	HEX_ARRAY_TUPLE_BEG    // Hexadecimal integer arrayTuple literal beginning `%x[`
	HEX_ARRAY_TUPLE_END    // Hexadecimal integer arrayTuple literal end `]`
	BIN_ARRAY_TUPLE_BEG    // Binary integer arrayTuple literal beginning `%b[`
	BIN_ARRAY_TUPLE_END    // Binary integer arrayTuple literal end `]`

	LABEL_SPECIAL_COLLECTION_LITERAL_END // Special collection literals end here

	CLOSURE_TYPE_BEG     // Beginning of a closure type `%|`
	HASH_SET_LITERAL_BEG // HashHashSet literal beginning `^[`
	TUPLE_LITERAL_BEG    // ArrayTuple literal beginning `%[`
	RECORD_LITERAL_BEG   // Record literal beginning `%{`
	DOC_COMMENT          // Documentation comment `##[` ... `]##`
	RAW_STRING           // Raw String literal delimited by single quotes `'` ... `'`
	CHAR_LITERAL         // Character literal delimited by backticks eg. `f`
	RAW_CHAR_LITERAL     // Raw Character literal delimited by r` eg. r`f`
	REGEX_BEG            // Beginning delimiter of Regex literals `%/`
	REGEX_CONTENT        // Regex literal content
	REGEX_INTERP_BEG     // Beginning of regex interpolation `${`
	REGEX_INTERP_END     // End of regex interpolation `}`
	REGEX_END            // Ending delimiter of Regex literals `/`

	// Regex flags start here
	LABEL_REGEX_FLAG_BEG

	REGEX_FLAG_i // Regex flag i
	REGEX_FLAG_m // Regex flag m
	REGEX_FLAG_U // Regex flag U
	REGEX_FLAG_a // Regex flag a
	REGEX_FLAG_x // Regex flag x
	REGEX_FLAG_s // Regex flag s

	LABEL_REGEX_FLAG_END // Regex flags end here

	STRING_BEG                     // Beginning delimiter of String literals `"`
	STRING_CONTENT                 // String literal content
	STRING_INTERP_LOCAL            // A local embedded in string interpolation eg. `$foo`
	STRING_INTERP_CONSTANT         // A constant embedded in string interpolation eg. `$Foo`
	STRING_INTERP_BEG              // Beginning of string interpolation `${`
	STRING_INSPECT_INTERP_BEG      // Beginning of string inspect interpolation `#{`
	STRING_INSPECT_INTERP_CONSTANT // A constant embedded in inspect string interpolation eg. `#Foo`
	STRING_INSPECT_INTERP_LOCAL    // A local embedded in inspect string interpolation eg. `#foo`
	STRING_INTERP_END              // End of string interpolation `}`
	STRING_END                     // Ending delimiter of String literals `"`

	// Int literals start here
	LABEL_INT_LITERAL_BEG
	INT                   // Int literal eg. `23`
	INT64                 // Int64 literal eg. `23i64`
	UINT64                // UInt64 literal eg. `23u64`
	INT32                 // Int32 literal eg. `23i32`
	UINT32                // UInt32 literal eg. `23u32`
	INT16                 // Int16 literal eg. `23i16`
	UINT16                // UInt16 literal eg. `23u16`
	INT8                  // Int8 literal eg. `23i8`
	UINT8                 // UInt8 literal eg. `23u8`
	LABEL_INT_LITERAL_END // Int literals end here

	// Float literals start here
	LABEL_FLOAT_LITERAL_BEG
	FLOAT                   // Float literal eg. `2.5`
	BIG_FLOAT               // BigFloat literal eg. `2.5bf`
	FLOAT64                 // Float64 literal eg. `2.5f64`
	FLOAT32                 // Float32 literal eg. `2.5f32`
	LABEL_FLOAT_LITERAL_END // Float literals end here

	LABEL_LITERAL_END // Literals end here

	// Keywords start here
	LABEL_KEYWORD_BEG
	NIL               // Keyword `nil`
	FALSE             // Keyword `false`
	TRUE              // Keyword `true`
	IF                // Keyword `if`
	ELSE              // Keyword `else`
	ELSIF             // Keyword `elsif`
	UNLESS            // Keyword `unless`
	WHILE             // Keyword `while`
	UNTIL             // Keyword `until`
	LOOP              // Keyword `loop`
	FOR               // Keyword `for`
	FORNUM            // Keyword `fornum`
	IN                // Keyword `in`
	OF                // Keyword `of`
	BREAK             // Keyword `break`
	CONTINUE          // Keyword `continue`
	RETURN            // Keyword `return`
	YIELD             // Keyword `yield`
	ASYNC             // Keyword `async`
	AWAIT             // Keyword `await`
	AWAIT_SYNC        // Keyword `await_sync`
	GO                // Keyword `go`
	DEF               // Keyword `def`
	SIG               // Keyword `sig`
	END               // Keyword `end`
	THEN              // Keyword `then`
	INIT              // Keyword `init`
	NOINIT            // Keyword `noinit`
	CLASS             // Keyword `class`
	STRUCT            // Keyword `struct`
	MODULE            // Keyword `module`
	MIXIN             // Keyword `mixin`
	INTERFACE         // Keyword `interface`
	INCLUDE           // Keyword `include`
	IMPLEMENT         // Keyword `implement`
	EXTEND            // Keyword `extend`
	ENUM              // Keyword `enum`
	TYPE              // Keyword `type`
	TYPEDEF           // Keyword `typedef`
	TYPEOF            // Keyword `typeof`
	VAR               // Keyword `var`
	VAL               // Keyword `val`
	CONST             // Keyword `const`
	THROW             // Keyword `throw`
	TRY               // Keyword `try`
	CATCH             // Keyword `catch`
	DO                // Keyword `do`
	FINALLY           // Keyword `finally`
	DEFER             // Keyword `defer`
	ALIAS             // Keyword `alias`
	AS                // Keyword `as`
	IS                // Keyword `is`
	SELF              // Keyword `self`
	SUPER             // Keyword `super`
	SWITCH            // Keyword `switch`
	CASE              // Keyword `case`
	MATCH             // Keyword `match`
	WITH              // Keyword `with`
	USING             // Keyword `using`
	BREAKPOINT        // Keyword `breakpoint`
	GETTER            // Keyword `getter`
	SETTER            // Keyword `setter`
	ATTR              // Keyword `attr`
	MUST              // Keyword `must`
	SINGLETON         // Keyword `singleton`
	ABSTRACT          // Keyword `abstract`
	SEALED            // Keyword `sealed`
	VOID              // Keyword `void`
	NEVER             // Keyword `never`
	NOTHING           // Keyword `nothing`
	ANY               // Keyword `any`
	PRIMITIVE         // Keyword `primitive`
	PUBLIC            // Keyword `public`
	PRIVATE           // Keyword `private`
	PROTECTED         // Keyword `protected`
	NATIVE            // Keyword `native`
	DEFAULT           // Keyword `default`
	MACRO             // Keyword `macro`
	BOOL              // Keyword `bool`
	NEW               // Keyword `new`
	EXTERN            // Keyword `extern`
	IMPORT            // Keyword `import`
	EXPORT            // Keyword `export`
	WHERE             // Keyword `where`
	UNTYPED           // Keyword `untyped`
	UNCHECKED         // Keyword `unchecked`
	GOTO              // Keyword `goto`
	QUOTE             // Keyword `quote`
	QUOTE_EXPR        // Keyword `quote_expr`
	QUOTE_TYPE        // Keyword `quote_type`
	QUOTE_PATTERN     // Keyword `quote_pattern`
	UNQUOTE           // Keyword `unquote`
	UNQUOTE_EXPR      // Keyword `unquote_expr`
	UNQUOTE_TYPE      // Keyword `unquote_type`
	UNQUOTE_IDENT     // Keyword `unquote_ident`
	UNQUOTE_CONST     // Keyword `unquote_const`
	UNQUOTE_IVAR      // Keyword `unquote_ivar`
	UNQUOTE_PATTERN   // Keyword `unquote_pattern`
	UNDEFINED         // Keyword `undefined`
	FUNC              // Keyword `func`
	OVERLOAD          // Keyword `overload`
	LABEL_KEYWORD_END // Keywords end here
)

// Maps keywords to their Token Type.
var Keywords = map[string]Type{
	"nil":             NIL,
	"false":           FALSE,
	"true":            TRUE,
	"if":              IF,
	"else":            ELSE,
	"elsif":           ELSIF,
	"unless":          UNLESS,
	"while":           WHILE,
	"until":           UNTIL,
	"loop":            LOOP,
	"for":             FOR,
	"fornum":          FORNUM,
	"in":              IN,
	"of":              OF,
	"break":           BREAK,
	"continue":        CONTINUE,
	"return":          RETURN,
	"yield":           YIELD,
	"async":           ASYNC,
	"await":           AWAIT,
	"await_sync":      AWAIT_SYNC,
	"go":              GO,
	"def":             DEF,
	"sig":             SIG,
	"end":             END,
	"then":            THEN,
	"init":            INIT,
	"noinit":          NOINIT,
	"class":           CLASS,
	"struct":          STRUCT,
	"module":          MODULE,
	"mixin":           MIXIN,
	"interface":       INTERFACE,
	"include":         INCLUDE,
	"implement":       IMPLEMENT,
	"extend":          EXTEND,
	"enum":            ENUM,
	"type":            TYPE,
	"typedef":         TYPEDEF,
	"typeof":          TYPEOF,
	"var":             VAR,
	"val":             VAL,
	"const":           CONST,
	"throw":           THROW,
	"try":             TRY,
	"catch":           CATCH,
	"do":              DO,
	"finally":         FINALLY,
	"defer":           DEFER,
	"alias":           ALIAS,
	"as":              AS,
	"is":              IS,
	"self":            SELF,
	"super":           SUPER,
	"switch":          SWITCH,
	"case":            CASE,
	"match":           MATCH,
	"with":            WITH,
	"using":           USING,
	"breakpoint":      BREAKPOINT,
	"getter":          GETTER,
	"setter":          SETTER,
	"attr":            ATTR,
	"must":            MUST,
	"singleton":       SINGLETON,
	"abstract":        ABSTRACT,
	"sealed":          SEALED,
	"void":            VOID,
	"never":           NEVER,
	"nothing":         NOTHING,
	"any":             ANY,
	"native":          NATIVE,
	"primitive":       PRIMITIVE,
	"public":          PUBLIC,
	"private":         PRIVATE,
	"protected":       PROTECTED,
	"default":         DEFAULT,
	"macro":           MACRO,
	"bool":            BOOL,
	"new":             NEW,
	"extern":          EXTERN,
	"import":          IMPORT,
	"export":          EXPORT,
	"where":           WHERE,
	"untyped":         UNTYPED,
	"unchecked":       UNCHECKED,
	"goto":            GOTO,
	"quote":           QUOTE,
	"quote_expr":      QUOTE_EXPR,
	"quote_type":      QUOTE_TYPE,
	"quote_pattern":   QUOTE_PATTERN,
	"unquote":         UNQUOTE,
	"unquote_expr":    UNQUOTE_EXPR,
	"unquote_type":    UNQUOTE_TYPE,
	"unquote_ident":   UNQUOTE_IDENT,
	"unquote_const":   UNQUOTE_CONST,
	"unquote_ivar":    UNQUOTE_IVAR,
	"unquote_pattern": UNQUOTE_PATTERN,
	"undefined":       UNDEFINED,
	"func":            FUNC,
	"overload":        OVERLOAD,
}

var tokenNames = [...]string{
	ERROR:                "ERROR",
	END_OF_FILE:          "END_OF_FILE",
	NEWLINE:              "NEWLINE",
	SEMICOLON:            ";",
	THICK_ARROW:          "=>",
	THIN_ARROW:           "->",
	WIGGLY_ARROW:         "~>",
	SHORT_UNQUOTE_BEG:    "!{",
	LPAREN:               "(",
	RPAREN:               ")",
	LBRACE:               "{",
	RBRACE:               "}",
	LBRACKET:             "[",
	QUESTION_LBRACKET:    "?[",
	RBRACKET:             "]",
	COMMA:                ",",
	DOT:                  ".",
	QUESTION_DOT:         "?.",
	DOT_DOT:              "..",
	QUESTION_DOT_DOT:     "?..",
	COLON:                ":",
	QUESTION:             "?",
	PLUS_PLUS:            "++",
	MINUS_MINUS:          "--",
	SCOPE_RES_OP:         "::",
	COLON_COLON_LBRACKET: "::[",
	DOT_COLON:            ".:",
	CLOSED_RANGE_OP:      "...",
	OPEN_RANGE_OP:        "<.<",
	RIGHT_OPEN_RANGE_OP:  "..<",
	LEFT_OPEN_RANGE_OP:   "<..",
	PIPE_OP:              "|>",

	EQUAL_OP:                "=",
	MINUS_EQUAL:             "-=",
	PLUS_EQUAL:              "+=",
	STAR_EQUAL:              "*=",
	SLASH_EQUAL:             "/=",
	STAR_STAR_EQUAL:         "**=",
	COLON_EQUAL:             ":=",
	TILDE_EQUAL:             "~=",
	AND_AND:                 "&&",
	AND_BANG:                "&!",
	AND_AND_EQUAL:           "&&=",
	OR_OR:                   "||",
	OR_BANG:                 "|!",
	OR_OR_EQUAL:             "||=",
	OR_EQUAL:                "|=",
	XOR_EQUAL:               "^=",
	QUESTION_QUESTION_EQUAL: "??=",
	LBITSHIFT_EQUAL:         "<<=",
	LTRIPLE_BITSHIFT_EQUAL:  "<<<=",
	RTRIPLE_BITSHIFT_EQUAL:  ">>>=",
	RBITSHIFT_EQUAL:         ">>=",
	PERCENT_EQUAL:           "%=",
	NOT_EQUAL:               "!=",
	LAX_NOT_EQUAL:           "!~",
	STRICT_NOT_EQUAL:        "!==",

	MINUS_AT:               "-@",
	PLUS_AT:                "+@",
	MINUS:                  "-",
	PLUS:                   "+",
	STAR:                   "*",
	SLASH:                  "/",
	STAR_STAR:              "**",
	LESS:                   "<",
	LESS_EQUAL:             "<=",
	GREATER:                ">",
	GREATER_EQUAL:          ">=",
	SPACESHIP_OP:           "<=>",
	EQUAL_EQUAL:            "==",
	LAX_EQUAL:              "=~",
	STRICT_EQUAL:           "===",
	TILDE:                  "~",
	AND:                    "&",
	AND_TILDE:              "&~",
	OR:                     "|",
	XOR:                    "^",
	QUESTION_QUESTION:      "??",
	BANG:                   "!",
	ISA_OP:                 "<:",
	REVERSE_ISA_OP:         ":>",
	INSTANCE_OF_OP:         "<<:",
	REVERSE_INSTANCE_OF_OP: ":>>",
	LBITSHIFT:              "<<",
	LTRIPLE_BITSHIFT:       "<<<",
	RBITSHIFT:              ">>",
	RTRIPLE_BITSHIFT:       ">>>",
	PERCENT:                "%",

	DOLLAR_IDENTIFIER:  "DOLLAR_IDENTIFIER",
	PUBLIC_IDENTIFIER:  "PUBLIC_IDENTIFIER",
	PRIVATE_IDENTIFIER: "PRIVATE_IDENTIFIER",
	PUBLIC_CONSTANT:    "PUBLIC_CONSTANT",
	PRIVATE_CONSTANT:   "PRIVATE_CONSTANT",

	INSTANCE_VARIABLE: "INSTANCE_VARIABLE",

	WORD_ARRAY_LIST_BEG:   "\\w[",
	WORD_ARRAY_LIST_END:   "] (WORD_ARRAY_LIST_END)",
	SYMBOL_ARRAY_LIST_BEG: "\\s[",
	SYMBOL_ARRAY_LIST_END: "] (SYMBOL_ARRAY_LIST_END)",
	HEX_ARRAY_LIST_BEG:    "\\x[",
	HEX_ARRAY_LIST_END:    "] (HEX_ARRAY_LIST_END)",
	BIN_ARRAY_LIST_BEG:    "\\b[",
	BIN_ARRAY_LIST_END:    "] (BIN_ARRAY_LIST_END)",

	WORD_HASH_SET_BEG:   "^w[",
	WORD_HASH_SET_END:   "] (WORD_HASH_SET_END)",
	SYMBOL_HASH_SET_BEG: "^s[",
	SYMBOL_HASH_SET_END: "] (SYMBOL_HASH_SET_END)",
	HEX_HASH_SET_BEG:    "^x[",
	HEX_HASH_SET_END:    "] (HEX_HASH_SET_END)",
	BIN_HASH_SET_BEG:    "^b[",
	BIN_HASH_SET_END:    "] (BIN_HASH_SET_END)",

	WORD_ARRAY_TUPLE_BEG:   "%w[",
	WORD_ARRAY_TUPLE_END:   "] (WORD_ARRAY_TUPLE_END)",
	SYMBOL_ARRAY_TUPLE_BEG: "%s[",
	SYMBOL_ARRAY_TUPLE_END: "] (SYMBOL_ARRAY_TUPLE_END)",
	HEX_ARRAY_TUPLE_BEG:    "%x[",
	HEX_ARRAY_TUPLE_END:    "] (HEX_ARRAY_TUPLE_END)",
	BIN_ARRAY_TUPLE_BEG:    "%b[",
	BIN_ARRAY_TUPLE_END:    "] (BIN_ARRAY_TUPLE_END)",

	HASH_SET_LITERAL_BEG:           "^[",
	TUPLE_LITERAL_BEG:              "%[",
	RECORD_LITERAL_BEG:             "%{",
	DOC_COMMENT:                    "DOC_COMMENT",
	RAW_STRING:                     "RAW_STRING",
	CHAR_LITERAL:                   "CHAR_LITERAL",
	RAW_CHAR_LITERAL:               "RAW_CHAR_LITERAL",
	CLOSURE_TYPE_BEG:               "%|",
	REGEX_BEG:                      "%/",
	REGEX_CONTENT:                  "REGEX_CONTENT",
	REGEX_INTERP_BEG:               "${ (REGEX_INTERP_BEG)",
	REGEX_INTERP_END:               "} (REGEX_INTERP_END)",
	REGEX_END:                      "/ (REGEX_END)",
	REGEX_FLAG_i:                   "i (REGEX_FLAG)",
	REGEX_FLAG_m:                   "m (REGEX_FLAG)",
	REGEX_FLAG_U:                   "U (REGEX_FLAG)",
	REGEX_FLAG_a:                   "a (REGEX_FLAG)",
	REGEX_FLAG_x:                   "x (REGEX_FLAG)",
	REGEX_FLAG_s:                   "s (REGEX_FLAG)",
	STRING_CONTENT:                 "STRING_CONTENT",
	STRING_INTERP_LOCAL:            "STRING_INTERP_LOCAL",
	STRING_INTERP_CONSTANT:         "STRING_INTERP_CONSTANT",
	STRING_INTERP_BEG:              "${ (STRING_INTERP_BEG)",
	STRING_INSPECT_INTERP_BEG:      "#{ (STRING_INSPECT_INTERP_BEG)",
	STRING_INSPECT_INTERP_CONSTANT: "STRING_INSPECT_INTERP_CONSTANT",
	STRING_INSPECT_INTERP_LOCAL:    "STRING_INSPECT_INTERP_LOCAL",
	STRING_INTERP_END:              "} (STRING_INTERP_END)",
	STRING_END:                     "\" (STRING_END)",
	INT:                            "INT",
	INT64:                          "INT64",
	UINT64:                         "UINT64",
	INT32:                          "INT32",
	UINT32:                         "UINT32",
	INT16:                          "INT16",
	UINT16:                         "UINT16",
	INT8:                           "INT8",
	UINT8:                          "UINT8",
	FLOAT:                          "FLOAT",
	BIG_FLOAT:                      "BIG_FLOAT",
	FLOAT64:                        "FLOAT64",
	FLOAT32:                        "FLOAT32",

	NIL:             "nil",
	FALSE:           "false",
	TRUE:            "true",
	IF:              "if",
	ELSE:            "else",
	ELSIF:           "elsif",
	UNLESS:          "unless",
	WHILE:           "while",
	UNTIL:           "until",
	LOOP:            "loop",
	FOR:             "for",
	FORNUM:          "fornum",
	IN:              "in",
	OF:              "of",
	BREAK:           "break",
	CONTINUE:        "continue",
	RETURN:          "return",
	YIELD:           "yield",
	ASYNC:           "async",
	AWAIT:           "await",
	AWAIT_SYNC:      "await_sync",
	GO:              "go",
	DEF:             "def",
	SIG:             "sig",
	END:             "end",
	THEN:            "then",
	INIT:            "init",
	NOINIT:          "noinit",
	CLASS:           "class",
	STRUCT:          "struct",
	MODULE:          "module",
	MIXIN:           "mixin",
	INTERFACE:       "interface",
	INCLUDE:         "include",
	IMPLEMENT:       "implement",
	EXTEND:          "extend",
	ENUM:            "enum",
	TYPE:            "type",
	TYPEDEF:         "typedef",
	TYPEOF:          "typeof",
	VAR:             "var",
	VAL:             "val",
	CONST:           "const",
	THROW:           "throw",
	TRY:             "try",
	CATCH:           "catch",
	DO:              "do",
	FINALLY:         "finally",
	DEFER:           "defer",
	ALIAS:           "alias",
	AS:              "as",
	IS:              "is",
	SELF:            "self",
	SUPER:           "super",
	SWITCH:          "switch",
	CASE:            "case",
	MATCH:           "match",
	WITH:            "with",
	USING:           "using",
	BREAKPOINT:      "breakpoint",
	GETTER:          "getter",
	SETTER:          "setter",
	ATTR:            "attr",
	MUST:            "must",
	SINGLETON:       "singleton",
	ABSTRACT:        "abstract",
	SEALED:          "sealed",
	VOID:            "void",
	NEVER:           "never",
	NOTHING:         "nothing",
	ANY:             "any",
	PRIMITIVE:       "primitive",
	PUBLIC:          "public",
	PRIVATE:         "private",
	PROTECTED:       "protected",
	NATIVE:          "native",
	DEFAULT:         "default",
	MACRO:           "macro",
	BOOL:            "bool",
	NEW:             "new",
	EXTERN:          "extern",
	IMPORT:          "import",
	EXPORT:          "export",
	WHERE:           "where",
	UNTYPED:         "untyped",
	UNCHECKED:       "unchecked",
	GOTO:            "goto",
	QUOTE:           "quote",
	QUOTE_EXPR:      "quote_expr",
	QUOTE_TYPE:      "quote_type",
	QUOTE_PATTERN:   "quote_pattern",
	UNQUOTE:         "unquote",
	UNQUOTE_EXPR:    "unquote_expr",
	UNQUOTE_TYPE:    "unquote_type",
	UNQUOTE_IDENT:   "unquote_ident",
	UNQUOTE_CONST:   "unquote_const",
	UNQUOTE_IVAR:    "unquote_ivar",
	UNQUOTE_PATTERN: "unquote_pattern",
	UNDEFINED:       "undefined",
	FUNC:            "func",
	OVERLOAD:        "overload",
}
