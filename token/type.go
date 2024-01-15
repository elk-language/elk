package token

// Represents the type of token
type Type uint8

// Name of the token.
func (t Type) String() string {
	if int(t) > len(tokenNames) {
		return "UNKNOWN"
	}

	return tokenNames[t]
}

// Returns `true` if the token can be a beginning of
// an argument to a function call without parentheses
// eg. `foo 2`
func (t Type) IsValidAsArgumentToNoParenFunctionCall() bool {
	switch t {
	case BANG, TILDE, LBRACE, PUBLIC_IDENTIFIER, PRIVATE_IDENTIFIER,
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
// a range value eg. `..2`
func (t Type) IsValidAsEndInRangeLiteral() bool {
	switch t {
	case SCOPE_RES_OP, BANG, TILDE, LBRACE, LPAREN, LBRACKET, PUBLIC_IDENTIFIER, PRIVATE_IDENTIFIER,
		PUBLIC_CONSTANT, PRIVATE_CONSTANT, INSTANCE_VARIABLE,
		RAW_STRING, STRING_BEG, FLOAT, NIL, FALSE, TRUE, LOOP, ENUM,
		VAR, VAL, CONST, DO, SELF, SUPER, SWITCH:
		return true
	}

	if t.IsIntLiteral() || t.IsSpecialCollectionLiteralBeg() {
		return true
	}

	return false
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
	return OpBegToken < t && t < LABEL_OP_END
}

// Check whether the token is a valid simple symbol content.
func (t Type) IsValidSimpleSymbolContent() bool {
	if t.IsIdentifier() || t == RAW_STRING || t.IsKeyword() || t.IsOverridableOperator() {
		return true
	}

	return false
}

// Check whether the token is a valid method name (without operators).
func (t Type) IsValidRegularMethodName() bool {
	if t == PUBLIC_IDENTIFIER || t == PRIVATE_IDENTIFIER || t.IsKeyword() {
		return true
	}

	return false
}

// Check whether the token is a valid method name (including operators).
func (t Type) IsValidMethodName() bool {
	if t.IsValidRegularMethodName() || t.IsOverridableOperator() {
		return true
	}

	return false
}

// Check whether the token is a valid method name in method
// call expressions.
func (t Type) IsValidPublicMethodName() bool {
	if t == PUBLIC_IDENTIFIER || t.IsKeyword() || t.IsOverridableOperator() {
		return true
	}

	return false
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
		STRICT_NOT_EQUAL, REF_EQUAL, REF_NOT_EQUAL:
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
	ZERO_VALUE        Type = iota // Zero value for Type
	ERROR                         // ERROR Token with a message
	END_OF_FILE                   // End Of File has been reached
	NEWLINE                       // Newline `\n`, `\r\n`
	SEMICOLON                     // SEMICOLON `;`
	THICK_ARROW                   // Thick arrow `=>`
	THIN_ARROW                    // Thin arrow `->` (closure arrow)
	WIGGLY_ARROW                  // Wiggly arrow `~>` (lambda arrow)
	LPAREN                        // Left parenthesis `(`
	RPAREN                        // Right parenthesis `)`
	LBRACE                        // Left brace `{`
	RBRACE                        // Right brace `}`
	LBRACKET                      // Left bracket `[`
	QUESTION_LBRACKET             // Safe access `?[`
	RBRACKET                      // Right bracket `]`
	COMMA                         // Comma `,`
	DOT                           // Dot `.`
	QUESTION_DOT                  // Safe method call operator `?.`
	COLON                         // Colon `:`
	QUESTION                      // Question mark `?`

	// Operators start here
	OpBegToken

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

	PLUS_PLUS          // Post increment operator `++`
	MINUS_MINUS        // Post decrement operator `--`
	SCOPE_RES_OP       // Scope resolution operator `::`
	RANGE_OP           // Inclusive range operator `..`
	EXCLUSIVE_RANGE_OP // Exclusive range operator `...`
	PIPE_OP            // Pipe operator `|>`
	AND_AND            // Logical and `&&`
	AND_BANG           // Logical expression sequencing operator `&!` with the precedence of `&&`
	OR_OR              // Logical or `||`
	OR_BANG            // Logical expression sequencing operator `|!` with the precedence of `||`
	NOT_EQUAL          // Not equal `!=`
	REF_NOT_EQUAL      // Reference not equal operator `=!=`
	STRICT_NOT_EQUAL   // Strict not equal `!==`

	// Overridable operators start here
	LABEL_OVERRIDABLE_OP_BEG
	MINUS                  // Minus `-`
	PLUS                   // Plus `+`
	STAR                   // Star `*`
	SLASH                  // Slash `/`
	STAR_STAR              // Two stars `**`
	LESS                   // Less than `<`
	LESS_EQUAL             // Less than or equal `<=`
	GREATER                // Greater than `>`
	GREATER_EQUAL          // Greater than or equal `>=`
	SPACESHIP_OP           // Spaceship operator `<=>`
	EQUAL_EQUAL            // Equal (comparison) `==`
	REF_EQUAL              // Reference equality operator `=:=`
	STRICT_EQUAL           // Strict equal `===`
	TILDE                  // Tilde `~`
	MATCH_OP               // Match operator `=~`
	AND                    // Bitwise and `&`
	OR                     // Bitwise or `|`
	XOR                    // Bitwise xor `^`
	QUESTION_QUESTION      // Nil coalescing operator `??`
	BANG                   // Logical not `!`
	ISA_OP                 // "is a" operator `<:`
	REVERSE_ISA_OP         // Reverse "is a" operator `:>`
	INSTANCE_OF_OP         // Instance of operator `<<:`
	REVERSE_INSTANCE_OF_OP // Reverse instance of operator `:>>`
	LBITSHIFT              // Left bitwise shift `<<`
	LTRIPLE_BITSHIFT       // Triple left bitwise shift `<<<`
	RBITSHIFT              // Right bitwise shift `>>`
	RTRIPLE_BITSHIFT       // Triple right bitwise shift `>>>`
	PERCENT                // Percent `%`
	LABEL_OP_END           // Operators end here

	// Identifiers start here
	LABEL_IDENTIFIER_BEG
	PUBLIC_IDENTIFIER    // Identifier
	PRIVATE_IDENTIFIER   // Identifier with a initial underscore
	PUBLIC_CONSTANT      // Constant (identifier with an initial capital letter)
	PRIVATE_CONSTANT     // Constant with an initial underscore
	LABEL_IDENTIFIER_END // Identifiers end here

	INSTANCE_VARIABLE  // Instance variable token eg. `@foo`
	SPECIAL_IDENTIFIER // Special identifier token eg. `$foo`

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

	HASH_SET_LITERAL_BEG // HashHashSet literal beginning `^[`
	TUPLE_LITERAL_BEG    // ArrayTuple literal beginning `%[`
	RECORD_LITERAL_BEG   // Record literal beginning `%{`
	DOC_COMMENT          // Documentation comment `##[` ... `]##`
	RAW_STRING           // Raw String literal delimited by single quotes `'` ... `'`
	CHAR_LITERAL         // Character literal delimited by double quotes
	RAW_CHAR_LITERAL     // Raw Character literal delimited by single quotes
	STRING_BEG           // Beginning delimiter of String literals `"`
	STRING_CONTENT       // String literal content
	STRING_INTERP_BEG    // Beginning of string interpolation `${`
	STRING_INTERP_END    // End of string interpolation `}`
	STRING_END           // Ending delimiter of String literals `"`

	LABEL_SPECIAL_COLLECTION_LITERAL_END // Special collection literals end here

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
	IN                // Keyword `in`
	HALT              // Keyword `halt`
	BREAK             // Keyword `break`
	CONTINUE          // Keyword `continue`
	RETURN            // Keyword `return`
	YIELD             // Keyword `yield`
	ASYNC             // Keyword `async`
	AWAIT             // Keyword `await`
	GO                // Keyword `go`
	DEF               // Keyword `def`
	SIG               // Keyword `sig`
	END               // Keyword `end`
	THEN              // Keyword `then`
	INIT              // Keyword `init`
	CLASS             // Keyword `class`
	STRUCT            // Keyword `struct`
	MODULE            // Keyword `module`
	MIXIN             // Keyword `mixin`
	INTERFACE         // Keyword `interface`
	INCLUDE           // Keyword `include`
	EXTEND            // Keyword `extend`
	ENHANCE           // Keyword `enhance`
	ENUM              // Keyword `enum`
	TYPE              // Keyword `type`
	TYPEDEF           // Keyword `typedef`
	VAR               // Keyword `var`
	VAL               // Keyword `val`
	CONST             // Keyword `const`
	THROW             // Keyword `throw`
	TRY               // Keyword `try`
	CATCH             // Keyword `catch`
	DO                // Keyword `do`
	ENSURE            // Keyword `ensure`
	ALIAS             // Keyword `alias`
	AS                // Keyword `as`
	IS                // Keyword `is`
	SELF              // Keyword `self`
	SUPER             // Keyword `super`
	SWITCH            // Keyword `switch`
	CASE              // Keyword `case`
	USING             // Keyword `using`
	BREAKPOINT        // Keyword `breakpoint`
	GETTER            // Keyword `getter`
	SETTER            // Keyword `setter`
	ACCESSOR          // Keyword `accessor`
	MUST              // Keyword `must`
	SINGLETON         // Keyword `singleton`
	ABSTRACT          // Keyword `abstract`
	SEALED            // Keyword `sealed`
	LABEL_KEYWORD_END // Keywords end here
)

// Maps keywords to their Token Type.
var Keywords = map[string]Type{
	"nil":        NIL,
	"false":      FALSE,
	"true":       TRUE,
	"if":         IF,
	"else":       ELSE,
	"elsif":      ELSIF,
	"unless":     UNLESS,
	"while":      WHILE,
	"until":      UNTIL,
	"loop":       LOOP,
	"for":        FOR,
	"in":         IN,
	"halt":       HALT,
	"break":      BREAK,
	"continue":   CONTINUE,
	"return":     RETURN,
	"yield":      YIELD,
	"async":      ASYNC,
	"await":      AWAIT,
	"go":         GO,
	"def":        DEF,
	"sig":        SIG,
	"end":        END,
	"then":       THEN,
	"init":       INIT,
	"class":      CLASS,
	"struct":     STRUCT,
	"module":     MODULE,
	"mixin":      MIXIN,
	"interface":  INTERFACE,
	"include":    INCLUDE,
	"extend":     EXTEND,
	"enhance":    ENHANCE,
	"enum":       ENUM,
	"type":       TYPE,
	"typedef":    TYPEDEF,
	"var":        VAR,
	"val":        VAL,
	"const":      CONST,
	"throw":      THROW,
	"try":        TRY,
	"catch":      CATCH,
	"do":         DO,
	"ensure":     ENSURE,
	"alias":      ALIAS,
	"as":         AS,
	"is":         IS,
	"self":       SELF,
	"super":      SUPER,
	"switch":     SWITCH,
	"case":       CASE,
	"using":      USING,
	"breakpoint": BREAKPOINT,
	"getter":     GETTER,
	"setter":     SETTER,
	"accessor":   ACCESSOR,
	"must":       MUST,
	"singleton":  SINGLETON,
	"abstract":   ABSTRACT,
	"sealed":     SEALED,
}

var tokenNames = [...]string{
	ERROR:              "ERROR",
	END_OF_FILE:        "END_OF_FILE",
	NEWLINE:            "NEWLINE",
	SEMICOLON:          ";",
	THICK_ARROW:        "=>",
	THIN_ARROW:         "->",
	WIGGLY_ARROW:       "~>",
	LPAREN:             "(",
	RPAREN:             ")",
	LBRACE:             "{",
	RBRACE:             "}",
	LBRACKET:           "[",
	QUESTION_LBRACKET:  "?[",
	RBRACKET:           "]",
	COMMA:              ",",
	DOT:                ".",
	QUESTION_DOT:       "?.",
	COLON:              ":",
	QUESTION:           "?",
	PLUS_PLUS:          "++",
	MINUS_MINUS:        "--",
	SCOPE_RES_OP:       "::",
	RANGE_OP:           "..",
	EXCLUSIVE_RANGE_OP: "...",
	PIPE_OP:            "|>",

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
	REF_NOT_EQUAL:           "=:=",
	STRICT_NOT_EQUAL:        "!==",

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
	REF_EQUAL:              "=:=",
	STRICT_EQUAL:           "===",
	TILDE:                  "~",
	MATCH_OP:               "=~",
	AND:                    "&",
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

	PUBLIC_IDENTIFIER:  "PUBLIC_IDENTIFIER",
	PRIVATE_IDENTIFIER: "PRIVATE_IDENTIFIER",
	PUBLIC_CONSTANT:    "PUBLIC_CONSTANT",
	PRIVATE_CONSTANT:   "PRIVATE_CONSTANT",

	INSTANCE_VARIABLE:  "INSTANCE_VARIABLE",
	SPECIAL_IDENTIFIER: "SPECIAL_IDENTIFIER",

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

	HASH_SET_LITERAL_BEG: "^[",
	TUPLE_LITERAL_BEG:    "%[",
	RECORD_LITERAL_BEG:   "%{",
	DOC_COMMENT:          "DOC_COMMENT",
	RAW_STRING:           "RAW_STRING",
	CHAR_LITERAL:         "CHAR_LITERAL",
	RAW_CHAR_LITERAL:     "RAW_CHAR_LITERAL",
	STRING_BEG:           "\" (STRING_BEG)",
	STRING_CONTENT:       "STRING_CONTENT",
	STRING_INTERP_BEG:    "${",
	STRING_INTERP_END:    "} (STRING_INTERP_END)",
	STRING_END:           "\" (STRING_END)",
	INT:                  "INT",
	INT64:                "INT64",
	UINT64:               "UINT64",
	INT32:                "INT32",
	UINT32:               "UINT32",
	INT16:                "INT16",
	UINT16:               "UINT16",
	INT8:                 "INT8",
	UINT8:                "UINT8",
	FLOAT:                "FLOAT",
	BIG_FLOAT:            "BIG_FLOAT",
	FLOAT64:              "FLOAT64",
	FLOAT32:              "FLOAT32",

	NIL:        "nil",
	FALSE:      "false",
	TRUE:       "true",
	IF:         "if",
	ELSE:       "else",
	ELSIF:      "elsif",
	UNLESS:     "unless",
	WHILE:      "while",
	UNTIL:      "until",
	LOOP:       "loop",
	FOR:        "for",
	IN:         "in",
	HALT:       "halt",
	BREAK:      "break",
	CONTINUE:   "continue",
	RETURN:     "return",
	YIELD:      "yield",
	ASYNC:      "async",
	AWAIT:      "await",
	GO:         "go",
	DEF:        "def",
	SIG:        "sig",
	END:        "end",
	THEN:       "then",
	INIT:       "init",
	CLASS:      "class",
	STRUCT:     "struct",
	MODULE:     "module",
	MIXIN:      "mixin",
	INTERFACE:  "interface",
	INCLUDE:    "include",
	EXTEND:     "extend",
	ENHANCE:    "enhance",
	ENUM:       "enum",
	TYPE:       "type",
	TYPEDEF:    "typedef",
	VAR:        "var",
	VAL:        "val",
	CONST:      "const",
	THROW:      "throw",
	TRY:        "try",
	CATCH:      "catch",
	DO:         "do",
	ENSURE:     "ensure",
	ALIAS:      "alias",
	AS:         "as",
	IS:         "is",
	SELF:       "self",
	SUPER:      "super",
	SWITCH:     "switch",
	CASE:       "case",
	USING:      "using",
	BREAKPOINT: "breakpoint",
	GETTER:     "getter",
	SETTER:     "setter",
	ACCESSOR:   "accessor",
	MUST:       "must",
	SINGLETON:  "singleton",
	ABSTRACT:   "abstract",
	SEALED:     "sealed",
}
