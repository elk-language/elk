package lexer

// Represents the type of token
type TokenType int

type Token struct {
	TokenType
	Value      string // Literal value of the token
	StartByte  int    // Index of the first byte of this token
	ByteLength int    // Number of bytes of the token
	Line       int    // Source line number where the token starts
	Column     int    // Source column number where the token starts
}

// Allocate a new End of File token.
func newEOF() *Token {
	return &Token{TokenType: EOFToken}
}

const (
	ErrorToken        TokenType = iota // Error Token with a message
	EOFToken                           // End Of File has been reached
	EndLineToken                       // Statement separator `\n`, `\r\n`
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
	AssignToken           // Assign `=`
	ScopeResOpToken       // Scope resolution operator `::`
	RangeOpToken          // Inclusive range operator `..`
	ExclusiveRangeOpToken // Exclusive range operator `...`
	PipeOpToken           // Pipe operator `|>`
	MinusEqualToken       // Minus equal `-=`
	PlusEqualToken        // Plus equal `+=`
	StarEqualToken        // Star equal `*=`
	PowerEqualToken       // Power equal `**=`
	ColonEqualToken       // Colon equal `:=`
	TildeEqualToken       // Tilde equal `~=`
	AndAndToken           // Logical and `&&`
	AndAndEqualToken      // Logical and equal `&&=`
	AndEqualToken         // Bitwise and equal `&=`
	OrOrToken             // Logical or `||`
	OrOrEqualToken        // Logical or `||=`
	OrEqualToken          // Bitwise or equal `|=`
	XorEqualToken         // Bitwise xor equal `^=`
	NilCoalesceEqualToken // Nil coalescing equal operator `??=`
	LBitShiftEqualToken   // Left bitwise shift equal `<<=`
	RBitShiftEqualToken   // Right bitwise shift equal `>>=`
	PercentEqualToken     // Percent equal `%=`
	NotEqualToken         // Not equal `!=`
	RefNotEqualToken      // Reference not equal operator `=!=`
	StrictNotEqualToken   // Strict not equal `!==`

	// Overridable operators start here
	OverridableOpBegToken
	MinusToken             // Minus `-`
	PlusToken              // Plus `+`
	StarToken              // Star `*`
	PowerToken             // Power `**`
	LessToken              // Less than `<`
	LessEqualToken         // Less than or equal `<=`
	GreaterToken           // Greater than `>`
	GreaterEqualToken      // Greater than or equal `>=`
	EqualToken             // Equal `==`
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
	SymbolBegToken       // Beginning of a Symbol literal `:`
	WordArrayBegToken    // Word array literal beginning `%w[`
	WordArrayEndToken    // Word array literal end `]`
	SymbolArrayBegToken  // Word array literal beginning `%s[`
	SymbolArrayEndToken  // Word array literal end `]`
	WordSetBegToken      // Word array literal beginning `%w{`
	WordSetEndToken      // Word array literal end `}`
	SymbolSetBegToken    // Word array literal beginning `%s{`
	SymbolSetEndToken    // Word array literal end `}`
	WordTupleBegToken    // Word array literal beginning `%w(`
	WordTupleEndToken    // Word array literal end `)`
	SymbolTupleBegToken  // Word array literal beginning `%s(`
	SymbolTupleEndToken  // Word array literal end `)`
	SetLiteralBegToken   // Set literal beginning `%{`
	TupleLiteralBegToken // Tuple literal beginning `%(`
	DocCommentToken      // Documentation comment `##[` ... `]##`
	RawStringToken       // Raw String literal delimited by single quotes `'` ... `'`
	StringBegToken       // Beginning delimiter of String literals `"`
	StringContentToken   // String literal content
	StringInterpBegToken // Beginning of string interpolation `${`
	StringInterpEndToken // End of string interpolation `}`
	StringEndToken       // Ending delimiter of String literals `"`
	HexIntToken          // Hexadecimal (base-16) Int literal eg. `0x5f`
	DuoIntToken          // Duodecimal (base-12) Int literal eg. `0d5b`
	DecIntToken          // Decimal (base-10) Int literal
	OctIntToken          // Octal (base-8) Int literal eg. `0o34`
	QuatIntToken         // Quaternary (base-4) Int literal eg. `0q31`
	BinIntToken          // Binary (base-2) Int literal eg. `0b1010`
	FloatToken           // Float literal
	LiteralEndToken      // Literals end here

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
	KeywordEndToken // any types lesser than this value can be considered keywords
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

// Check whether the token marks the end of the file.
func (t TokenType) IsEOF() bool {
	return t == EOFToken
}

// Check whether the token is a keyword.
func (t TokenType) IsKeyword() bool {
	return KeywordBegToken < t && t < KeywordEndToken
}

// Check whether the token is a literal.
func (t TokenType) IsLiteral() bool {
	return LiteralBegToken < t && t < LiteralEndToken
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
