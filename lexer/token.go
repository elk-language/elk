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
	ErrorToken             TokenType = iota // Error Token with a message
	EOFToken                                // End Of File has been reached
	SeparatorToken                          // Statement separator `\n`, `\r\n` or `;`
	ThickArrowToken                         // Thick arrow `=>`
	ThinArrowToken                          // Thin arrow `->` (closure arrow)
	WigglyArrowToken                        // Wiggly arrow `~>` (lambda arrow)
	LParenToken                             // Left parenthesis `(`
	RParenToken                             // Right parenthesis `)`
	LBraceToken                             // Left brace `{`
	RBraceToken                             // Right brace `}`
	LBracketToken                           // Left bracket `[`
	RBracketToken                           // Right bracket `]`
	CommaToken                              // Comma `,`
	DotToken                                // Dot `.`
	ColonToken                              // Colon `:`
	QuestionMarkToken                       // Question mark `?`
	IdentifierToken                         // Identifier
	PrivateIdentifierToken                  // Identifier with a initial underscore
	ConstantToken                           // Constant (identifier with an initial capital letter)
	PrivateConstantToken                    // Constant with an initial underscore
	InstanceVariableToken                   // Instance variable token eg. `@foo`

	// Literals start here
	LiteralBegToken
	PercentWToken        // Word collection literal prefix `%w`
	PercentSToken        // Symbol collection literal prefix `%s`
	PercentIToken        // Integer collection literal prefix `%i`
	PercentFToken        // Float collection literal prefix `%f`
	SetLiteralBegToken   // Set literal beginning `%{`
	TupleLiteralBegToken // Tuple literal beginning `%(`
	DocCommentToken      // Documentation comment `##[` ... `]##`
	RawStringToken       // Raw String literal delimited by single quotes `'` ... `'`
	StringBegToken       // Beginning delimiter of String literals `"`
	StringContentToken   // String literal content
	StringInterpBegToken // Beginning of string interpolation `${`
	StringInterpEndToken // End of string interpolation `}`
	StringEndToken       // Ending delimiter of String literals `"`
	IntToken             // Int literal
	FloatToken           // Float literal
	LiteralEndToken      // Literals end here

	// Operators start here
	OpBegToken
	AssignToken           // Assign `=`
	ScopeResOperatorToken // Scope resolution operator `::`
	PipeOperatorToken     // Pipe operator `|>`
	MinusEqualToken       // Minus equal `-=`
	TheAnswerToken        // The answer to the great question of life, the universe, and everything.
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
	MatchOperatorToken     // Match operator `=~`
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
}

// Check whether the token type is a keyword.
func (t TokenType) isKeyword() bool {
	return KeywordBegToken < t && t < KeywordEndToken
}
