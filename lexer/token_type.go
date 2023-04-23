package lexer

// Represents the type of token
type TokenType uint8

// Name of the token.
func (t TokenType) String() string {
	if int(t) > len(tokenNames) {
		return "UNKNOWN"
	}

	return tokenNames[t]
}

// Check whether the token marks the end of the file.
func (t TokenType) IsEndOfFile() bool {
	return t == EndOfFileToken
}

// Check whether the token is a keyword.
func (t TokenType) IsKeyword() bool {
	return KeywordBegToken < t && t < KeywordEndToken
}

// Check whether the token is a literal.
func (t TokenType) IsLiteral() bool {
	return LiteralBegToken < t && t < LiteralEndToken
}

// Check whether the token is an Int literal.
func (t TokenType) IsIntLiteral() bool {
	return IntLiteralBegToken < t && t < IntLiteralEndToken
}

// Check whether the token is a an operator.
func (t TokenType) IsOperator() bool {
	return OpBegToken < t && t < OpEndToken
}

// Check whether the token is an overridable operator.
func (t TokenType) IsOverridableOperator() bool {
	return OverridableOpBegToken < t && t < OpEndToken
}

// Check whether the token is an identifier.
func (t TokenType) IsIdentifier() bool {
	return IdentifierBegToken < t && t < IdentifierEndToken
}

// Check whether the token can separate statements.
func (t TokenType) IsStatementSeparator() bool {
	return t == EndLineToken || t == SemicolonToken
}

// Check whether the token is an assignment operator.
func (t TokenType) IsAssignmentOperator() bool {
	return AssignOpBegToken < t && t < AssignOpEndToken
}

// Check whether the token is an equality operator.
func (t TokenType) IsEqualityOperator() bool {
	switch t {
	case EqualEqualToken, NotEqualToken, StrictEqualToken, StrictNotEqualToken, RefEqualToken, RefNotEqualToken:
		return true
	default:
		return false
	}
}

// Check whether the token is a comparison operator.
func (t TokenType) IsComparisonOperator() bool {
	switch t {
	case LessToken, LessEqualToken, GreaterToken, GreaterEqualToken, SubtypeToken, ReverseSubtypeToken, InstanceOfToken, ReverseInstanceOfToken:
		return true
	default:
		return false
	}
}

const (
	ZeroToken         TokenType = iota // Zero value for TokenType
	ErrorToken                         // Error Token with a message
	EndOfFileToken                     // End Of File has been reached
	EndLineToken                       // End-line `\n`, `\r\n`
	SemicolonToken                     // Semicolon `;`
	ThickArrowToken                    // Thick arrow `=>`
	ThinArrowToken                     // Thin arrow `->` (closure arrow)
	WigglyArrowToken                   // Wiggly arrow `~>` (lambda arrow)
	LParenToken                        // Left parenthesis `(`
	RParenToken                        // Right parenthesis `)`
	LBraceToken                        // Left brace `{`
	RBraceToken                        // Right brace `}`
	LBracketToken                      // Left bracket `[`
	RBracketToken                      // Right bracket `]`
	CommaToken                         // Comma `,`
	DotToken                           // Dot `.`
	ColonToken                         // Colon `:`
	QuestionMarkToken                  // Question mark `?`

	// Operators start here
	OpBegToken

	// Assignment operators start here
	AssignOpBegToken
	EqualToken            // Equal (assignment) `=`
	MinusEqualToken       // Minus equal `-=`
	PlusEqualToken        // Plus equal `+=`
	StarEqualToken        // Star equal `*=`
	SlashEqualToken       // Slash equal `/=`
	StarStarEqualToken    // Two stars equal `**=`
	TildeEqualToken       // Tilde equal `~=`
	AndAndEqualToken      // Logical and equal `&&=`
	AndEqualToken         // Bitwise and equal `&=`
	OrOrEqualToken        // Logical or `||=`
	OrEqualToken          // Bitwise or equal `|=`
	XorEqualToken         // Bitwise xor equal `^=`
	NilCoalesceEqualToken // Nil coalescing equal operator `??=`
	LBitShiftEqualToken   // Left bitwise shift equal `<<=`
	RBitShiftEqualToken   // Right bitwise shift equal `>>=`
	PercentEqualToken     // Percent equal `%=`
	AssignOpEndToken      // Assignment operators end here

	ColonEqualToken       // Colon equal `:=`
	ScopeResOpToken       // Scope resolution operator `::`
	RangeOpToken          // Inclusive range operator `..`
	ExclusiveRangeOpToken // Exclusive range operator `...`
	PipeOpToken           // Pipe operator `|>`
	AndAndToken           // Logical and `&&`
	OrOrToken             // Logical or `||`
	NotEqualToken         // Not equal `!=`
	RefNotEqualToken      // Reference not equal operator `=!=`
	StrictNotEqualToken   // Strict not equal `!==`

	// Overridable operators start here
	OverridableOpBegToken
	MinusToken             // Minus `-`
	PlusToken              // Plus `+`
	StarToken              // Star `*`
	SlashToken             // Slash `/`
	StarStarToken          // Two stars `**`
	LessToken              // Less than `<`
	LessEqualToken         // Less than or equal `<=`
	GreaterToken           // Greater than `>`
	GreaterEqualToken      // Greater than or equal `>=`
	EqualEqualToken        // Equal (comparison) `==`
	RefEqualToken          // Reference equality operator `=:=`
	StrictEqualToken       // Strict equal `===`
	TildeToken             // Tilde `~`
	MatchOpToken           // Match operator `=~`
	AndToken               // Bitwise and `&`
	OrToken                // Bitwise or `|`
	XorToken               // Bitwise xor `^`
	NilCoalesceToken       // Nil coalescing operator `??`
	BangToken              // Logical not `!`
	SubtypeToken           // Subtype operator `<:`
	ReverseSubtypeToken    // Reverse subtype operator `:>`
	InstanceOfToken        // Instance of operator `<<:`
	ReverseInstanceOfToken // Reverse instance of operator `:>>`
	LBitShiftToken         // Left bitwise shift `<<`
	RBitShiftToken         // Right bitwise shift `>>`
	PercentToken           // Percent `%`
	OpEndToken             // Operators end here

	// Identifiers start here
	IdentifierBegToken
	IdentifierToken        // Identifier
	PrivateIdentifierToken // Identifier with a initial underscore
	ConstantToken          // Constant (identifier with an initial capital letter)
	PrivateConstantToken   // Constant with an initial underscore
	IdentifierEndToken     // Identifiers end here

	InstanceVariableToken // Instance variable token eg. `@foo`

	// Literals start here
	LiteralBegToken
	SymbolBegToken // Beginning of a Symbol literal `:`

	WordArrayBegToken   // Word array literal beginning `%w[`
	WordArrayEndToken   // Word array literal end `]`
	SymbolArrayBegToken // Symbol array literal beginning `%s[`
	SymbolArrayEndToken // Symbol array literal end `]`
	HexArrayBegToken    // Hexadecimal integer array literal beginning `%x[`
	HexArrayEndToken    // Hexadecimal integer array literal end `]`
	BinArrayBegToken    // Binary integer array literal beginning `%b[`
	BinArrayEndToken    // Binary integer array literal end `]`

	WordSetBegToken   // Word set literal beginning `%w{`
	WordSetEndToken   // Word set literal end `}`
	SymbolSetBegToken // Symbol set literal beginning `%s{`
	SymbolSetEndToken // Symbol set literal end `}`
	HexSetBegToken    // Hexadecimal integer set literal beginning `%x{`
	HexSetEndToken    // Hexadecimal integer set literal end `}`
	BinSetBegToken    // Binary integer set literal beginning `%b{`
	BinSetEndToken    // Binary integer set literal end `}`

	WordTupleBegToken   // Word tuple literal beginning `%w(`
	WordTupleEndToken   // Word tuple literal end `)`
	SymbolTupleBegToken // Symbol tuple literal beginning `%s(`
	SymbolTupleEndToken // Symbol tuple literal end `)`
	HexTupleBegToken    // Hexadecimal integer tuple literal beginning `%x(`
	HexTupleEndToken    // Hexadecimal integer tuple literal end `)`
	BinTupleBegToken    // Binary integer tuple literal beginning `%b(`
	BinTupleEndToken    // Binary integer tuple literal end `)`

	SetLiteralBegToken   // Set literal beginning `%{`
	TupleLiteralBegToken // Tuple literal beginning `%(`
	DocCommentToken      // Documentation comment `##[` ... `]##`
	RawStringToken       // Raw String literal delimited by single quotes `'` ... `'`
	StringBegToken       // Beginning delimiter of String literals `"`
	StringContentToken   // String literal content
	StringInterpBegToken // Beginning of string interpolation `${`
	StringInterpEndToken // End of string interpolation `}`
	StringEndToken       // Ending delimiter of String literals `"`

	// Int literals start here
	IntLiteralBegToken
	HexIntToken        // Hexadecimal (base-16) Int literal eg. `0x5f`
	DuoIntToken        // Duodecimal (base-12) Int literal eg. `0d5b`
	DecIntToken        // Decimal (base-10) Int literal
	OctIntToken        // Octal (base-8) Int literal eg. `0o34`
	QuatIntToken       // Quaternary (base-4) Int literal eg. `0q31`
	BinIntToken        // Binary (base-2) Int literal eg. `0b1010`
	IntLiteralEndToken // Int literals end here

	FloatToken      // Float literal
	LiteralEndToken // Literals end here

	// Keywords start here
	KeywordBegToken
	NilToken        // Keyword `nil`
	FalseToken      // Keyword `false`
	TrueToken       // Keyword `true`
	IfToken         // Keyword `if`
	ElseToken       // Keyword `else`
	ElsifToken      // Keyword `elsif`
	UnlessToken     // Keyword `unless`
	WhileToken      // Keyword `while`
	UntilToken      // Keyword `until`
	LoopToken       // Keyword `loop`
	BreakToken      // Keyword `break`
	ReturnToken     // Keyword `return`
	DefToken        // Keyword `def`
	EndToken        // Keyword `end`
	ThenToken       // Keyword `then`
	ClassToken      // Keyword `class`
	ModuleToken     // Keyword `module`
	MixinToken      // Keyword `mixin`
	InterfaceToken  // Keyword `interface`
	TypeToken       // Keyword `type`
	VarToken        // Keyword `var`
	ConstToken      // Keyword `const`
	ThrowToken      // Keyword `throw`
	CatchToken      // Keyword `catch`
	DoToken         // Keyword `do`
	EnsureToken     // Keyword `ensure`
	AliasToken      // Keyword `alias`
	SelfToken       // Keyword `self`
	SuperToken      // Keyword `super`
	SwitchToken     // Keyword `switch`
	CaseToken       // Keyword `case`
	UsingToken      // Keyword `using`
	KeywordEndToken // Keywords end here
)

// Maps keywords to their Token Type.
var keywords = map[string]TokenType{
	"nil":       NilToken,
	"false":     FalseToken,
	"true":      TrueToken,
	"if":        IfToken,
	"else":      ElseToken,
	"elsif":     ElsifToken,
	"unless":    UnlessToken,
	"while":     WhileToken,
	"until":     UntilToken,
	"loop":      LoopToken,
	"break":     BreakToken,
	"return":    ReturnToken,
	"def":       DefToken,
	"end":       EndToken,
	"then":      ThenToken,
	"class":     ClassToken,
	"module":    ModuleToken,
	"mixin":     MixinToken,
	"interface": InterfaceToken,
	"type":      TypeToken,
	"var":       VarToken,
	"const":     ConstToken,
	"throw":     ThrowToken,
	"catch":     CatchToken,
	"do":        DoToken,
	"ensure":    EnsureToken,
	"alias":     AliasToken,
	"self":      SelfToken,
	"super":     SuperToken,
	"switch":    SwitchToken,
	"case":      CaseToken,
	"using":     UsingToken,
}

var tokenNames = [...]string{
	ErrorToken:            "Error",
	EndOfFileToken:        "EndOfFile",
	EndLineToken:          "EndLine",
	SemicolonToken:        ";",
	ThickArrowToken:       "=>",
	ThinArrowToken:        "->",
	WigglyArrowToken:      "~>",
	LParenToken:           "(",
	RParenToken:           ")",
	LBraceToken:           "{",
	RBraceToken:           "}",
	LBracketToken:         "[",
	RBracketToken:         "]",
	CommaToken:            ",",
	DotToken:              ".",
	ColonToken:            ":",
	QuestionMarkToken:     "?",
	ScopeResOpToken:       "::",
	RangeOpToken:          "..",
	ExclusiveRangeOpToken: "...",
	PipeOpToken:           "|>",

	EqualToken:            "=",
	MinusEqualToken:       "-=",
	PlusEqualToken:        "+=",
	StarEqualToken:        "*=",
	SlashEqualToken:       "/=",
	StarStarEqualToken:    "**=",
	ColonEqualToken:       ":=",
	TildeEqualToken:       "~=",
	AndAndToken:           "&&",
	AndAndEqualToken:      "&&=",
	OrOrToken:             "||",
	OrOrEqualToken:        "||=",
	OrEqualToken:          "|=",
	XorEqualToken:         "^=",
	NilCoalesceEqualToken: "??=",
	LBitShiftEqualToken:   "<<=",
	RBitShiftEqualToken:   ">>=",
	PercentEqualToken:     "%=",
	NotEqualToken:         "!=",
	RefNotEqualToken:      "=:=",
	StrictNotEqualToken:   "!==",

	MinusToken:             "-",
	PlusToken:              "+",
	StarToken:              "*",
	SlashToken:             "/",
	StarStarToken:          "**",
	LessToken:              "<",
	LessEqualToken:         "<=",
	GreaterToken:           ">",
	GreaterEqualToken:      ">=",
	EqualEqualToken:        "==",
	RefEqualToken:          "=:=",
	StrictEqualToken:       "===",
	TildeToken:             "~",
	MatchOpToken:           "=~",
	AndToken:               "&",
	OrToken:                "|",
	XorToken:               "^",
	NilCoalesceToken:       "??",
	BangToken:              "!",
	SubtypeToken:           "<:",
	ReverseSubtypeToken:    ":>",
	InstanceOfToken:        "<<:",
	ReverseInstanceOfToken: ":>>",
	LBitShiftToken:         "<<",
	RBitShiftToken:         ">>",
	PercentToken:           "%",

	IdentifierToken:        "Identifier",
	PrivateIdentifierToken: "PrivateIdentifier",
	ConstantToken:          "Constant",
	PrivateConstantToken:   "PrivateConstant",

	InstanceVariableToken: "InstanceVariable",

	SymbolBegToken: "SymbolBeg",

	WordArrayBegToken:   "%w[",
	WordArrayEndToken:   "] (WordArrayEnd)",
	SymbolArrayBegToken: "%s[",
	SymbolArrayEndToken: "] (SymbolArrayEnd)",
	HexArrayBegToken:    "%x[",
	HexArrayEndToken:    "] (HexArrayEnd)",
	BinArrayBegToken:    "%b[",
	BinArrayEndToken:    "] (BinArrayEnd)",

	WordSetBegToken:   "%w{",
	WordSetEndToken:   "} (WordSetEnd)",
	SymbolSetBegToken: "%s{",
	SymbolSetEndToken: "} (SymbolSetEnd)",
	HexSetBegToken:    "%x{",
	HexSetEndToken:    "} (HexSetEnd)",
	BinSetBegToken:    "%b{",
	BinSetEndToken:    "} (BinSetEnd)",

	WordTupleBegToken:   "%w(",
	WordTupleEndToken:   ") (WordTupleEnd)",
	SymbolTupleBegToken: "%s(",
	SymbolTupleEndToken: ") (SymbolTupleEnd)",
	HexTupleBegToken:    "%x(",
	HexTupleEndToken:    ") (HexTupleEnd)",
	BinTupleBegToken:    "%b(",
	BinTupleEndToken:    ") (BinTupleEnd)",

	SetLiteralBegToken:   "%{",
	TupleLiteralBegToken: "%(",
	DocCommentToken:      "DocComment",
	RawStringToken:       "RawString",
	StringBegToken:       "\" (StringBeg)",
	StringContentToken:   "StringContent",
	StringInterpBegToken: "${",
	StringInterpEndToken: "} (StringInterpEnd)",
	StringEndToken:       "\" (StringEnd)",
	HexIntToken:          "HexInt",
	DuoIntToken:          "DuoInt",
	DecIntToken:          "DecInt",
	OctIntToken:          "OctInt",
	QuatIntToken:         "QuatInt",
	BinIntToken:          "BinInt",
	FloatToken:           "Float",

	NilToken:       "nil",
	FalseToken:     "false",
	TrueToken:      "true",
	IfToken:        "if",
	ElseToken:      "else",
	ElsifToken:     "elsif",
	UnlessToken:    "unless",
	WhileToken:     "while",
	UntilToken:     "until",
	LoopToken:      "loop",
	BreakToken:     "break",
	ReturnToken:    "return",
	DefToken:       "def",
	EndToken:       "end",
	ThenToken:      "then",
	ClassToken:     "class",
	ModuleToken:    "module",
	MixinToken:     "mixin",
	InterfaceToken: "interface",
	TypeToken:      "type",
	VarToken:       "var",
	ConstToken:     "const",
	ThrowToken:     "throw",
	CatchToken:     "catch",
	DoToken:        "do",
	EnsureToken:    "ensure",
	AliasToken:     "alias",
	SelfToken:      "self",
	SuperToken:     "super",
	SwitchToken:    "switch",
	CaseToken:      "case",
	UsingToken:     "using",
}
